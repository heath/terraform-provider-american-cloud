package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-americancloud/internal/client"
)

var (
	_ resource.Resource              = &loadBalancerRuleResource{}
	_ resource.ResourceWithConfigure = &loadBalancerRuleResource{}
)

func NewLoadBalancerRuleResource() resource.Resource {
	return &loadBalancerRuleResource{}
}

type loadBalancerRuleResource struct {
	client *client.Client
}

type loadBalancerRuleResourceModel struct {
	ID             types.String `tfsdk:"id"`
	IPAddressID    types.String `tfsdk:"ip_address_id"`
	Name           types.String `tfsdk:"name"`
	Algorithm      types.String `tfsdk:"algorithm"`
	PublicPort     types.Int64  `tfsdk:"public_port"`
	PrivatePort    types.Int64  `tfsdk:"private_port"`
	Protocol       types.String `tfsdk:"protocol"`
	SourceCIDRList types.String `tfsdk:"source_cidr_list"`
	Description    types.String `tfsdk:"description"`
	State          types.String `tfsdk:"state"`
	Created        types.String `tfsdk:"created"`
}

func (r *loadBalancerRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_load_balancer_rule"
}

func (r *loadBalancerRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a load balancer rule on a public IP. Name, algorithm, and description can be updated.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ip_address_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"algorithm": schema.StringAttribute{
				Required:    true,
				Description: "Load balancing algorithm (e.g. \"roundrobin\", \"leastconn\", \"source\").",
			},
			"public_port": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"private_port": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"protocol": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"source_cidr_list": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"description": schema.StringAttribute{
				Optional: true,
			},
			"state": schema.StringAttribute{
				Computed: true,
			},
			"created": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *loadBalancerRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *loadBalancerRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan loadBalancerRuleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &client.CreateLoadBalancerRuleRequest{
		Name:        plan.Name.ValueString(),
		Algorithm:   plan.Algorithm.ValueString(),
		PublicPort:  int(plan.PublicPort.ValueInt64()),
		PrivatePort: int(plan.PrivatePort.ValueInt64()),
	}
	if !plan.Protocol.IsNull() {
		createReq.Protocol = plan.Protocol.ValueString()
	}
	if !plan.SourceCIDRList.IsNull() {
		createReq.SourceCIDRList = plan.SourceCIDRList.ValueString()
	}
	if !plan.Description.IsNull() {
		createReq.Description = plan.Description.ValueString()
	}

	rule, err := r.client.CreateLoadBalancerRule(ctx, plan.IPAddressID.ValueString(), createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating load balancer rule", err.Error())
		return
	}

	plan.ID = types.StringValue(rule.ID)
	plan.State = types.StringValue(rule.State)
	plan.Created = types.StringValue(rule.Created)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *loadBalancerRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state loadBalancerRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	rules, err := r.client.ListLoadBalancerRules(ctx, state.IPAddressID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading load balancer rules", err.Error())
		return
	}

	var found *client.LoadBalancerRule
	for i := range rules {
		if rules[i].ID == state.ID.ValueString() {
			found = &rules[i]
			break
		}
	}
	if found == nil {
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(found.Name)
	state.Algorithm = types.StringValue(found.Algorithm)
	state.PublicPort = types.Int64Value(int64(found.PublicPort))
	state.PrivatePort = types.Int64Value(int64(found.PrivatePort))
	state.State = types.StringValue(found.State)
	state.Created = types.StringValue(found.Created)
	if found.Protocol != "" {
		state.Protocol = types.StringValue(found.Protocol)
	}
	if found.SourceCIDRList != "" {
		state.SourceCIDRList = types.StringValue(found.SourceCIDRList)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *loadBalancerRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan loadBalancerRuleResourceModel
	var state loadBalancerRuleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateReq := &client.UpdateLoadBalancerRuleRequest{}
	if !plan.Name.Equal(state.Name) {
		updateReq.Name = plan.Name.ValueString()
	}
	if !plan.Algorithm.Equal(state.Algorithm) {
		updateReq.Algorithm = plan.Algorithm.ValueString()
	}
	if !plan.Description.Equal(state.Description) && !plan.Description.IsNull() {
		updateReq.Description = plan.Description.ValueString()
	}

	rule, err := r.client.UpdateLoadBalancerRule(ctx, state.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError("Error updating load balancer rule", err.Error())
		return
	}

	plan.ID = state.ID
	plan.State = types.StringValue(rule.State)
	plan.Created = types.StringValue(rule.Created)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *loadBalancerRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state loadBalancerRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteLoadBalancerRule(ctx, state.ID.ValueString())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Error deleting load balancer rule", err.Error())
	}
}
