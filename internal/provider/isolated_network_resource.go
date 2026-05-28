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
	_ resource.Resource                = &isolatedNetworkResource{}
	_ resource.ResourceWithConfigure   = &isolatedNetworkResource{}
	_ resource.ResourceWithImportState = &isolatedNetworkResource{}
)

func NewIsolatedNetworkResource() resource.Resource {
	return &isolatedNetworkResource{}
}

type isolatedNetworkResource struct {
	client *client.Client
}

type isolatedNetworkResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Region      types.String `tfsdk:"region"`
	Netmask     types.String `tfsdk:"netmask"`
	Gateway     types.String `tfsdk:"gateway"`
	Status      types.String `tfsdk:"status"`
	Created     types.String `tfsdk:"created"`
}

func (r *isolatedNetworkResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_isolated_network"
}

func (r *isolatedNetworkResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an American Cloud isolated network.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the isolated network.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Description: "A description of the network.",
			},
			"region": schema.StringAttribute{
				Required:    true,
				Description: "The region for the network.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"netmask": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The netmask for the network.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"gateway": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The gateway for the network.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"status": schema.StringAttribute{
				Computed:    true,
				Description: "Current status of the network.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"created": schema.StringAttribute{
				Computed:    true,
				Description: "Timestamp when the network was created.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *isolatedNetworkResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *isolatedNetworkResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan isolatedNetworkResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &client.CreateIsolatedNetworkRequest{
		Name:   plan.Name.ValueString(),
		Region: plan.Region.ValueString(),
	}
	if !plan.Description.IsNull() {
		createReq.Description = plan.Description.ValueString()
	}
	if !plan.Netmask.IsNull() {
		createReq.Netmask = plan.Netmask.ValueString()
	}
	if !plan.Gateway.IsNull() {
		createReq.Gateway = plan.Gateway.ValueString()
	}

	network, err := r.client.CreateIsolatedNetwork(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating isolated network", err.Error())
		return
	}

	plan.ID = types.StringValue(network.ID)
	plan.Name = types.StringValue(network.Name)
	plan.Description = types.StringValue(network.Description)
	plan.Region = types.StringValue(network.Region)
	plan.Netmask = types.StringValue(network.CIDR)
	plan.Gateway = types.StringValue(network.Gateway)
	plan.Status = types.StringValue(network.Status)
	plan.Created = types.StringValue(network.Created)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *isolatedNetworkResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state isolatedNetworkResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	network, err := r.client.GetIsolatedNetwork(ctx, state.ID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading isolated network", err.Error())
		return
	}

	state.Name = types.StringValue(network.Name)
	state.Description = types.StringValue(network.Description)
	state.Region = types.StringValue(network.Region)
	state.Netmask = types.StringValue(network.CIDR)
	state.Gateway = types.StringValue(network.Gateway)
	state.Status = types.StringValue(network.Status)
	state.Created = types.StringValue(network.Created)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *isolatedNetworkResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan isolatedNetworkResourceModel
	var state isolatedNetworkResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &client.UpdateIsolatedNetworkRequest{}
	if !plan.Name.Equal(state.Name) {
		updateReq.Name = plan.Name.ValueString()
	}
	if !plan.Description.Equal(state.Description) {
		updateReq.Description = plan.Description.ValueString()
	}

	network, err := r.client.UpdateIsolatedNetwork(ctx, state.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating isolated network", err.Error())
		return
	}

	plan.ID = state.ID
	plan.Name = types.StringValue(network.Name)
	plan.Description = types.StringValue(network.Description)
	plan.Region = types.StringValue(network.Region)
	plan.Netmask = types.StringValue(network.CIDR)
	plan.Gateway = types.StringValue(network.Gateway)
	plan.Status = types.StringValue(network.Status)
	plan.Created = types.StringValue(network.Created)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *isolatedNetworkResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state isolatedNetworkResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteIsolatedNetwork(ctx, state.ID.ValueString())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Error deleting isolated network", err.Error())
	}
}

func (r *isolatedNetworkResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
