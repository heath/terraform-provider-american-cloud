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
	_ resource.Resource              = &portForwardingResource{}
	_ resource.ResourceWithConfigure = &portForwardingResource{}
)

func NewPortForwardingResource() resource.Resource {
	return &portForwardingResource{}
}

type portForwardingResource struct {
	client *client.Client
}

type portForwardingResourceModel struct {
	ID            types.String `tfsdk:"id"`
	IPAddressID   types.String `tfsdk:"ip_address_id"`
	VMID          types.String `tfsdk:"vm_id"`
	Protocol      types.String `tfsdk:"protocol"`
	PublicPort    types.Int64  `tfsdk:"public_port"`
	PrivatePort   types.Int64  `tfsdk:"private_port"`
	OpenFirewall  types.Bool   `tfsdk:"open_firewall"`
	State         types.String `tfsdk:"state"`
	Created       types.String `tfsdk:"created"`
}

func (r *portForwardingResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_port_forwarding_rule"
}

func (r *portForwardingResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a port forwarding rule on a public IP. All attributes are ForceNew.",
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
			"vm_id": schema.StringAttribute{
				Required:    true,
				Description: "The VM to forward traffic to.",
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
			"open_firewall": schema.BoolAttribute{
				Optional:    true,
				Description: "If true, also creates a firewall rule for the public port.",
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

func (r *portForwardingResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *portForwardingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan portForwardingResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	createReq := &client.CreatePortForwardingRuleRequest{
		PrivatePort: int(plan.PrivatePort.ValueInt64()),
		PublicPort:  int(plan.PublicPort.ValueInt64()),
		Protocol:    plan.Protocol.ValueString(),
		VMID:        plan.VMID.ValueString(),
	}
	if !plan.OpenFirewall.IsNull() {
		createReq.OpenFirewall = plan.OpenFirewall.ValueBool()
	}

	rule, err := r.client.CreatePortForwardingRule(ctx, plan.IPAddressID.ValueString(), createReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating port forwarding rule", err.Error())
		return
	}

	plan.ID = types.StringValue(rule.ID)
	plan.State = types.StringValue(rule.State)
	plan.Created = types.StringValue(rule.Created)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *portForwardingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state portForwardingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	rules, err := r.client.ListPortForwardingRules(ctx, state.IPAddressID.ValueString())
	if err != nil {
		if client.IsNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Error reading port forwarding rules", err.Error())
		return
	}

	var found *client.PortForwardingRule
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

	state.Protocol = types.StringValue(found.Protocol)
	state.PublicPort = types.Int64Value(int64(found.PublicPort))
	state.PrivatePort = types.Int64Value(int64(found.PrivatePort))
	state.VMID = types.StringValue(found.VMID)
	state.State = types.StringValue(found.State)
	state.Created = types.StringValue(found.Created)

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *portForwardingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddError("Update not supported", "Port forwarding rules are immutable — all changes require replacement.")
}

func (r *portForwardingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state portForwardingResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeletePortForwardingRule(ctx, state.ID.ValueString())
	if err != nil && !client.IsNotFound(err) {
		resp.Diagnostics.AddError("Error deleting port forwarding rule", err.Error())
	}
}
