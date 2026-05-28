package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-americancloud/internal/client"
)

var (
	_ resource.Resource                = &egressRuleResource{}
	_ resource.ResourceWithConfigure   = &egressRuleResource{}
	_ resource.ResourceWithImportState = &egressRuleResource{}
)

func NewEgressRuleResource() resource.Resource {
	return &egressRuleResource{}
}

type egressRuleResource struct {
	client *client.Client
}

type egressRuleResourceModel struct {
	ID             types.String `tfsdk:"id"`
	NetworkID      types.String `tfsdk:"network_id"`
	Protocol       types.String `tfsdk:"protocol"`
	StartPort      types.Int64  `tfsdk:"start_port"`
	EndPort        types.Int64  `tfsdk:"end_port"`
	SourceCIDRList types.String `tfsdk:"source_cidr_list"`
	DestCIDRList   types.String `tfsdk:"dest_cidr_list"`
	State          types.String `tfsdk:"state"`
	Created        types.String `tfsdk:"created"`
}

func (r *egressRuleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_egress_rule"
}

func (r *egressRuleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an egress firewall rule. All attributes are ForceNew.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"network_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"protocol": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"start_port": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"end_port": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"source_cidr_list": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"dest_cidr_list": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"state": schema.StringAttribute{
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

func (r *egressRuleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *egressRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan egressRuleResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &client.CreateEgressRuleRequest{
		Protocol:  plan.Protocol.ValueString(),
		StartPort: strconv.FormatInt(plan.StartPort.ValueInt64(), 10),
		EndPort:   strconv.FormatInt(plan.EndPort.ValueInt64(), 10),
		NetworkID: plan.NetworkID.ValueString(),
	}
	if !plan.SourceCIDRList.IsNull() {
		createReq.SourceCIDRList = plan.SourceCIDRList.ValueString()
	}
	if !plan.DestCIDRList.IsNull() {
		createReq.DestCIDRList = plan.DestCIDRList.ValueString()
	}

	rule, err := r.client.CreateEgressRule(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating egress rule", err.Error())
		return
	}

	plan.ID = types.StringValue(rule.ID)
	plan.State = types.StringValue(rule.State)
	plan.Created = types.StringValue(rule.Created)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *egressRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state egressRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	rule, err := r.client.GetEgressRule(ctx, state.ID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading egress rule", err.Error())
		return
	}

	state.Protocol = types.StringValue(rule.Protocol)
	if v, err := strconv.ParseInt(rule.StartPort, 10, 64); err == nil {
		state.StartPort = types.Int64Value(v)
	}
	if v, err := strconv.ParseInt(rule.EndPort, 10, 64); err == nil {
		state.EndPort = types.Int64Value(v)
	}
	state.NetworkID = types.StringValue(rule.NetworkID)
	state.State = types.StringValue(rule.State)
	state.Created = types.StringValue(rule.Created)
	if rule.SourceCIDRList != "" {
		state.SourceCIDRList = types.StringValue(rule.SourceCIDRList)
	}
	if rule.DestCIDRList != "" {
		state.DestCIDRList = types.StringValue(rule.DestCIDRList)
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *egressRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "Egress rules are immutable — all changes require replacement.")
}

func (r *egressRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state egressRuleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteEgressRule(ctx, state.ID.ValueString())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Error deleting egress rule", err.Error())
	}
}

func (r *egressRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
