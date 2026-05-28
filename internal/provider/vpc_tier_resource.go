package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-americancloud/internal/client"
)

var (
	_ resource.Resource              = &vpcTierResource{}
	_ resource.ResourceWithConfigure = &vpcTierResource{}
)

func NewVPCTierResource() resource.Resource {
	return &vpcTierResource{}
}

type vpcTierResource struct {
	client *client.Client
}

type vpcTierResourceModel struct {
	ID          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	VPCID       types.String `tfsdk:"vpc_id"`
	Gateway     types.String `tfsdk:"gateway"`
	Netmask     types.String `tfsdk:"netmask"`
	ACLID       types.String `tfsdk:"acl_id"`
	Status      types.String `tfsdk:"status"`
	Created     types.String `tfsdk:"created"`
}

func (r *vpcTierResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vpc_tier"
}

func (r *vpcTierResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a tier (subnet) within an American Cloud VPC. All attributes are ForceNew (no update endpoint).",
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
			"description": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"vpc_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the parent VPC.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"gateway": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"netmask": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"acl_id": schema.StringAttribute{
				Optional:    true,
				Description: "Network ACL list ID.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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

func (r *vpcTierResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *vpcTierResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan vpcTierResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &client.CreateVPCTierRequest{
		Name:    plan.Name.ValueString(),
		VPCID:   plan.VPCID.ValueString(),
		Gateway: plan.Gateway.ValueString(),
		Netmask: plan.Netmask.ValueString(),
	}
	if !plan.Description.IsNull() {
		createReq.Description = plan.Description.ValueString()
	}
	if !plan.ACLID.IsNull() {
		createReq.ACLID = plan.ACLID.ValueString()
	}

	tier, err := r.client.CreateVPCTier(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating VPC tier", err.Error())
		return
	}

	plan.ID = types.StringValue(tier.ID)
	plan.Status = types.StringValue(tier.Status)
	plan.Created = types.StringValue(tier.Created)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *vpcTierResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state vpcTierResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// VPC tiers are read via the parent VPC detail endpoint
	vpc, err := r.client.GetVPC(ctx, state.VPCID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading VPC for tier", err.Error())
		return
	}

	// Find our tier in the VPC's tier list
	var found *client.VPCTier
	for i := range vpc.Tiers {
		if vpc.Tiers[i].ID == state.ID.ValueString() {
			found = &vpc.Tiers[i]
			break
		}
	}
	if found == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(found.Name)
	state.Description = types.StringValue(found.Description)
	state.Gateway = types.StringValue(found.Gateway)
	state.Netmask = types.StringValue(found.Netmask)
	if found.ACLID != "" {
		state.ACLID = types.StringValue(found.ACLID)
	}
	state.Status = types.StringValue(found.Status)
	state.Created = types.StringValue(found.Created)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *vpcTierResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// All attributes are ForceNew — this should never be called.
	resp.Diagnostics.AddError("Update not supported", "VPC tier attributes all require replacement.")
}

func (r *vpcTierResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// The AC API doesn't expose a standalone delete-tier endpoint.
	// Tiers are removed when the parent VPC is deleted.
	// We simply remove from state.
}
