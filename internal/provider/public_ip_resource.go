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
	_ resource.Resource                = &publicIPResource{}
	_ resource.ResourceWithConfigure   = &publicIPResource{}
	_ resource.ResourceWithImportState = &publicIPResource{}
)

func NewPublicIPResource() resource.Resource {
	return &publicIPResource{}
}

type publicIPResource struct {
	client *client.Client
}

type publicIPResourceModel struct {
	ID              types.String `tfsdk:"id"`
	IPAddress       types.String `tfsdk:"ip_address"`
	NetworkID       types.String `tfsdk:"network_id"`
	VPCID           types.String `tfsdk:"vpc_id"`
	Region          types.String `tfsdk:"region"`
	StaticNatVMID   types.String `tfsdk:"static_nat_vm_id"`
	IsSourceNat     types.Bool   `tfsdk:"is_source_nat"`
	IsStaticNat     types.Bool   `tfsdk:"is_static_nat"`
	Status          types.String `tfsdk:"status"`
}

func (r *publicIPResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_public_ip"
}

func (r *publicIPResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an American Cloud public IP address with optional static NAT.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ip_address": schema.StringAttribute{
				Computed:    true,
				Description: "The allocated public IP address.",
			},
			"network_id": schema.StringAttribute{
				Optional:    true,
				Description: "Isolated network ID to associate the IP with.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vpc_id": schema.StringAttribute{
				Optional:    true,
				Description: "VPC ID to associate the IP with.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"region": schema.StringAttribute{
				Optional:    true,
				Description: "Region for the IP reservation.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"static_nat_vm_id": schema.StringAttribute{
				Optional:    true,
				Description: "VM ID to enable static NAT to. Set to enable, remove to disable.",
			},
			"is_source_nat": schema.BoolAttribute{
				Computed: true,
			},
			"is_static_nat": schema.BoolAttribute{
				Computed: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *publicIPResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *publicIPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan publicIPResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	reserveReq := &client.ReservePublicIPRequest{}
	if !plan.NetworkID.IsNull() {
		reserveReq.NetworkID = plan.NetworkID.ValueString()
	}
	if !plan.VPCID.IsNull() {
		reserveReq.VPCID = plan.VPCID.ValueString()
	}
	if !plan.Region.IsNull() {
		reserveReq.Region = plan.Region.ValueString()
	}

	ip, err := r.client.ReservePublicIP(ctx, reserveReq)
	if err != nil {
		resp.Diagnostics.AddError("Error reserving public IP", err.Error())
		return
	}

	plan.ID = types.StringValue(ip.ID)
	plan.IPAddress = types.StringValue(ip.IPAddress)
	plan.IsSourceNat = types.BoolValue(ip.IsSourceNat)
	plan.IsStaticNat = types.BoolValue(ip.IsStaticNat)
	plan.Status = types.StringValue(ip.Status)

	// Enable static NAT if requested
	if !plan.StaticNatVMID.IsNull() && !plan.StaticNatVMID.IsUnknown() {
		err := r.client.EnableStaticNat(ctx, ip.ID, plan.StaticNatVMID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError("Error enabling static NAT", err.Error())
			return
		}
		plan.IsStaticNat = types.BoolValue(true)
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *publicIPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state publicIPResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	ip, err := r.client.GetPublicIP(ctx, state.ID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading public IP", err.Error())
		return
	}

	state.IPAddress = types.StringValue(ip.IPAddress)
	state.IsSourceNat = types.BoolValue(ip.IsSourceNat)
	state.IsStaticNat = types.BoolValue(ip.IsStaticNat)
	state.Status = types.StringValue(ip.Status)
	if ip.NetworkID != "" {
		state.NetworkID = types.StringValue(ip.NetworkID)
	}
	if ip.VPCID != "" {
		state.VPCID = types.StringValue(ip.VPCID)
	}
	if ip.Region != "" {
		state.Region = types.StringValue(ip.Region)
	}
	if ip.IsStaticNat && ip.VMID != "" {
		state.StaticNatVMID = types.StringValue(ip.VMID)
	} else {
		state.StaticNatVMID = types.StringNull()
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *publicIPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan publicIPResourceModel
	var state publicIPResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id := state.ID.ValueString()

	// Handle static NAT changes
	oldVMID := ""
	if !state.StaticNatVMID.IsNull() {
		oldVMID = state.StaticNatVMID.ValueString()
	}
	newVMID := ""
	if !plan.StaticNatVMID.IsNull() {
		newVMID = plan.StaticNatVMID.ValueString()
	}

	if oldVMID != newVMID {
		// Disable existing static NAT first
		if oldVMID != "" {
			err := r.client.DisableStaticNat(ctx, id)
			if err != nil {
				resp.Diagnostics.AddError("Error disabling static NAT", err.Error())
				return
			}
		}
		// Enable new static NAT
		if newVMID != "" {
			err := r.client.EnableStaticNat(ctx, id, newVMID)
			if err != nil {
				resp.Diagnostics.AddError("Error enabling static NAT", err.Error())
				return
			}
		}
	}

	// Re-read state
	ip, err := r.client.GetPublicIP(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Error reading public IP after update", err.Error())
		return
	}

	plan.ID = state.ID
	plan.IPAddress = types.StringValue(ip.IPAddress)
	plan.IsSourceNat = types.BoolValue(ip.IsSourceNat)
	plan.IsStaticNat = types.BoolValue(ip.IsStaticNat)
	plan.Status = types.StringValue(ip.Status)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *publicIPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state publicIPResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Disable static NAT before releasing
	if !state.StaticNatVMID.IsNull() && state.StaticNatVMID.ValueString() != "" {
		_ = r.client.DisableStaticNat(ctx, state.ID.ValueString())
	}

	err := r.client.ReleasePublicIP(ctx, state.ID.ValueString())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Error releasing public IP", err.Error())
	}
}

func (r *publicIPResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
