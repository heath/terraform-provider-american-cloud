package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-americancloud/internal/client"
)

var (
	_ resource.Resource                = &vmResource{}
	_ resource.ResourceWithConfigure   = &vmResource{}
	_ resource.ResourceWithImportState = &vmResource{}
)

func NewVMResource() resource.Resource {
	return &vmResource{}
}

type vmResource struct {
	client *client.Client
}

type vmResourceModel struct {
	ID                 types.String  `tfsdk:"id"`
	Name               types.String  `tfsdk:"name"`
	Region             types.String  `tfsdk:"region"`
	Image              types.String  `tfsdk:"image"`
	VMPackage          types.String  `tfsdk:"vm_package"`
	CPU                types.Int64   `tfsdk:"cpu"`
	MemoryMB           types.Int64   `tfsdk:"memory_mb"`
	RootDiskGB         types.Int64   `tfsdk:"root_disk_gb"`
	Network            types.String  `tfsdk:"network"`
	SubscriptionPeriod types.String  `tfsdk:"subscription_period"`
	Tags               types.List    `tfsdk:"tags"`
	Keypairs           types.List    `tfsdk:"keypairs"`
	Userdata           types.String  `tfsdk:"userdata"`
	NetworkAccess      types.String  `tfsdk:"network_access"`
	Status             types.String  `tfsdk:"status"`
	IPAddress          types.String  `tfsdk:"ip_address"`
	Hostname           types.String  `tfsdk:"hostname"`
	Created            types.String  `tfsdk:"created"`
	Password           types.String  `tfsdk:"password"`
	Timeouts           *timeoutModel `tfsdk:"timeouts"`
}

type timeoutModel struct {
	Create types.String `tfsdk:"create"`
	Update types.String `tfsdk:"update"`
	Delete types.String `tfsdk:"delete"`
}

func (r *vmResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vm"
}

func (r *vmResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an American Cloud virtual machine.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier of the VM.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the virtual machine.",
			},
			"region": schema.StringAttribute{
				Required:    true,
				Description: "The region to deploy the VM in (e.g. \"us-central-1\").",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"image": schema.StringAttribute{
				Required:    true,
				Description: "The image label to use for the VM (e.g. \"ubuntu-24.04\").",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vm_package": schema.StringAttribute{
				Optional:    true,
				Description: "The VM package label. Mutually exclusive with cpu/memory_mb/root_disk_gb.",
			},
			"cpu": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Number of virtual CPUs. Used with vm_specs (custom sizing).",
			},
			"memory_mb": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Memory in MB. Used with vm_specs (custom sizing).",
			},
			"root_disk_gb": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Description: "Root disk size in GB. Used with vm_specs (custom sizing).",
			},
			"network": schema.StringAttribute{
				Optional:    true,
				Description: "Network ID to attach the VM to.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"subscription_period": schema.StringAttribute{
				Optional:    true,
				Description: "Billing period (e.g. \"hourly\", \"monthly\").",
			},
			"tags": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "Tags to apply to the VM.",
			},
			"keypairs": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
				Description: "SSH keypair names to inject into the VM.",
			},
			"userdata": schema.StringAttribute{
				Optional:    true,
				Description: "Cloud-init userdata (base64 encoded).",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"network_access": schema.StringAttribute{
				Optional:    true,
				Description: "Network access mode.",
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Current status of the VM (e.g. Running, Stopped).",
			},
			"ip_address": schema.StringAttribute{
				Computed:    true,
				Description: "The IP address assigned to the VM.",
			},
			"hostname": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The hostname of the VM.",
			},
			"created": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the VM was created.",
			},
			"password": schema.StringAttribute{
				Computed:    true,
				Sensitive:   true,
				Description: "The initial password for the VM (only available at creation time).",
			},
			"timeouts": schema.SingleNestedAttribute{
				Optional:    true,
				Description: "Timeout configuration.",
				Attributes: map[string]schema.Attribute{
					"create": schema.StringAttribute{
						Optional:    true,
						Description: "Timeout for create operations. Defaults to 30m.",
					},
					"update": schema.StringAttribute{
						Optional:    true,
						Description: "Timeout for update operations. Defaults to 15m.",
					},
					"delete": schema.StringAttribute{
						Optional:    true,
						Description: "Timeout for delete operations. Defaults to 15m.",
					},
				},
			},
		},
	}
}

func (r *vmResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T.", req.ProviderData),
		)
		return
	}
	r.client = c
}

func (r *vmResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vmResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createTimeout := 30 * time.Minute
	if plan.Timeouts != nil && !plan.Timeouts.Create.IsNull() {
		d, err := time.ParseDuration(plan.Timeouts.Create.ValueString())
		if err == nil {
			createTimeout = d
		}
	}
	ctx, cancel := context.WithTimeout(ctx, createTimeout)
	defer cancel()

	createReq := &client.CreateVMRequest{
		Name:   plan.Name.ValueString(),
		Region: plan.Region.ValueString(),
		Image:  plan.Image.ValueString(),
	}

	if !plan.VMPackage.IsNull() && !plan.VMPackage.IsUnknown() {
		createReq.VMPackage = plan.VMPackage.ValueString()
	}

	if (!plan.CPU.IsNull() && !plan.CPU.IsUnknown()) ||
		(!plan.MemoryMB.IsNull() && !plan.MemoryMB.IsUnknown()) ||
		(!plan.RootDiskGB.IsNull() && !plan.RootDiskGB.IsUnknown()) {
		createReq.VMSpecs = &client.VMSpecs{}
		if !plan.CPU.IsNull() {
			createReq.VMSpecs.VCPU = int(plan.CPU.ValueInt64())
		}
		if !plan.MemoryMB.IsNull() {
			createReq.VMSpecs.MemoryMB = int(plan.MemoryMB.ValueInt64())
		}
		if !plan.RootDiskGB.IsNull() {
			createReq.VMSpecs.RootDiskGB = int(plan.RootDiskGB.ValueInt64())
		}
	}

	if !plan.Network.IsNull() {
		createReq.Network = plan.Network.ValueString()
	}
	if !plan.SubscriptionPeriod.IsNull() {
		createReq.SubscriptionPeriod = plan.SubscriptionPeriod.ValueString()
	}
	if !plan.Userdata.IsNull() {
		createReq.Userdata = plan.Userdata.ValueString()
	}
	if !plan.NetworkAccess.IsNull() {
		createReq.NetworkAccess = plan.NetworkAccess.ValueString()
	}

	if !plan.Tags.IsNull() {
		var tags []string
		diags = plan.Tags.ElementsAs(ctx, &tags, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		createReq.Tags = tags
	}

	if !plan.Keypairs.IsNull() {
		var keypairs []string
		diags = plan.Keypairs.ElementsAs(ctx, &keypairs, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
		createReq.Keypairs = keypairs
	}

	vm, err := r.client.CreateVM(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VM", err.Error())
		return
	}

	vm, err = r.client.WaitForVMReady(ctx, vm.ID, 10*time.Second)
	if err != nil {
		resp.Diagnostics.AddError("Error waiting for VM to be ready", err.Error())
		return
	}

	mapVMToState(ctx, vm, &plan, &resp.Diagnostics)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *vmResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vmResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	vm, err := r.client.GetVM(ctx, state.ID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError(
			"Error reading VM",
			fmt.Sprintf("Could not read VM %s: %s", state.ID.ValueString(), err.Error()),
		)
		return
	}

	mapVMToState(ctx, vm, &state, &resp.Diagnostics)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *vmResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not yet implemented", "VM update will be added in Phase 4.")
}

func (r *vmResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddError("Delete not yet implemented", "VM delete will be added in Phase 4.")
}

func (r *vmResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapVMToState(ctx context.Context, vm *client.VM, state *vmResourceModel, diagnostics *diag.Diagnostics) {
	state.ID = types.StringValue(vm.ID)
	state.Name = types.StringValue(vm.Name)
	state.Status = types.StringValue(vm.Status)
	state.Region = types.StringValue(vm.Region)
	state.Image = types.StringValue(vm.Image)
	state.CPU = types.Int64Value(int64(vm.CPU))
	state.MemoryMB = types.Int64Value(int64(vm.MemoryMB))
	state.RootDiskGB = types.Int64Value(int64(vm.RootDiskGB))
	state.Created = types.StringValue(vm.Created)

	if vm.VMPackage != "" {
		state.VMPackage = types.StringValue(vm.VMPackage)
	}
	if vm.IPAddress != "" {
		state.IPAddress = types.StringValue(vm.IPAddress)
	} else {
		state.IPAddress = types.StringNull()
	}
	if vm.NetworkID != "" {
		state.Network = types.StringValue(vm.NetworkID)
	}
	if vm.Hostname != "" {
		state.Hostname = types.StringValue(vm.Hostname)
	}
	if vm.Password != "" {
		state.Password = types.StringValue(vm.Password)
	}

	if len(vm.Tags) > 0 {
		tagList, d := types.ListValueFrom(ctx, types.StringType, vm.Tags)
		diagnostics.Append(d...)
		state.Tags = tagList
	}

	if len(vm.KeypairNames) > 0 {
		kpList, d := types.ListValueFrom(ctx, types.StringType, vm.KeypairNames)
		diagnostics.Append(d...)
		state.Keypairs = kpList
	}
}
