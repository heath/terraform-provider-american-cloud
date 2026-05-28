package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-americancloud/internal/client"
)

var (
	_ resource.Resource                = &blockStorageResource{}
	_ resource.ResourceWithConfigure   = &blockStorageResource{}
	_ resource.ResourceWithImportState = &blockStorageResource{}
)

func NewBlockStorageResource() resource.Resource {
	return &blockStorageResource{}
}

type blockStorageResource struct {
	client *client.Client
}

type blockStorageResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	SizeGB  types.Int64  `tfsdk:"size_gb"`
	Region  types.String `tfsdk:"region"`
	VMID    types.String `tfsdk:"vm_id"`
	Status  types.String `tfsdk:"status"`
	Created types.String `tfsdk:"created"`
}

func (r *blockStorageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_block_storage"
}

func (r *blockStorageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an American Cloud block storage volume. Supports resize (grow only) and attach/detach via vm_id.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"size_gb": schema.Int64Attribute{
				Required:    true,
				Description: "Size in GB. Can only be increased (grow only).",
			},
			"region": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vm_id": schema.StringAttribute{
				Optional:    true,
				Description: "VM ID to attach the volume to. Set to attach, remove to detach.",
			},
			"status": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *blockStorageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", fmt.Sprintf("Expected *client.Client, got: %T.", req.ProviderData))
		return
	}
	r.client = c
}

func (r *blockStorageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan blockStorageResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &client.CreateVolumeRequest{
		Name:   plan.Name.ValueString(),
		SizeGB: int(plan.SizeGB.ValueInt64()),
		Region: plan.Region.ValueString(),
	}

	volume, err := r.client.CreateVolume(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating block storage volume", err.Error())
		return
	}

	plan.ID = types.StringValue(volume.ID)
	plan.Status = types.StringValue(volume.Status)
	plan.Created = types.StringValue(volume.Created)

	// Attach if vm_id specified
	if !plan.VMID.IsNull() && !plan.VMID.IsUnknown() {
		err := r.client.AttachVolume(ctx, volume.ID, plan.VMID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error attaching volume", err.Error())
			return
		}
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *blockStorageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state blockStorageResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	volume, err := r.client.GetVolume(ctx, state.ID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading block storage volume", err.Error())
		return
	}

	state.Name = types.StringValue(volume.Name)
	state.SizeGB = types.Int64Value(int64(volume.SizeGB))
	state.Region = types.StringValue(volume.Region)
	state.Status = types.StringValue(volume.Status)
	state.Created = types.StringValue(volume.Created)
	if volume.VMID != "" {
		state.VMID = types.StringValue(volume.VMID)
	} else {
		state.VMID = types.StringNull()
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *blockStorageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan blockStorageResourceModel
	var state blockStorageResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	// Resize (grow only)
	if !plan.SizeGB.Equal(state.SizeGB) {
		err := r.client.ResizeVolume(ctx, id, int(plan.SizeGB.ValueInt64()))
		if err != nil {
			resp.Diagnostics.AddError("Error resizing volume", err.Error())
			return
		}
	}

	// Attach/detach
	oldVMID := ""
	if !state.VMID.IsNull() {
		oldVMID = state.VMID.ValueString()
	}
	newVMID := ""
	if !plan.VMID.IsNull() {
		newVMID = plan.VMID.ValueString()
	}

	if oldVMID != newVMID {
		if oldVMID != "" {
			err := r.client.DetachVolume(ctx, id)
			if err != nil {
				resp.Diagnostics.AddError("Error detaching volume", err.Error())
				return
			}
		}
		if newVMID != "" {
			err := r.client.AttachVolume(ctx, id, newVMID)
			if err != nil {
				resp.Diagnostics.AddError("Error attaching volume", err.Error())
				return
			}
		}
	}

	// Re-read
	volume, err := r.client.GetVolume(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Error reading volume after update", err.Error())
		return
	}

	plan.ID = state.ID
	plan.Status = types.StringValue(volume.Status)
	plan.Created = types.StringValue(volume.Created)
	plan.SizeGB = types.Int64Value(int64(volume.SizeGB))
	if volume.VMID != "" {
		plan.VMID = types.StringValue(volume.VMID)
	} else {
		plan.VMID = types.StringNull()
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *blockStorageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state blockStorageResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Detach first if attached
	if !state.VMID.IsNull() && state.VMID.ValueString() != "" {
		_ = r.client.DetachVolume(ctx, state.ID.ValueString())
	}

	err := r.client.DeleteVolume(ctx, state.ID.ValueString())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Error deleting block storage volume", err.Error())
	}
}

func (r *blockStorageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
