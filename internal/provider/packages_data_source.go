package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-americancloud/internal/client"
)

var (
	_ datasource.DataSource              = &packagesDataSource{}
	_ datasource.DataSourceWithConfigure = &packagesDataSource{}
)

func NewPackagesDataSource() datasource.DataSource {
	return &packagesDataSource{}
}

type packagesDataSource struct {
	client *client.Client
}

type packagesDataSourceModel struct {
	Packages []packageModel `tfsdk:"packages"`
}

type packageModel struct {
	ID            types.String `tfsdk:"id"`
	Label         types.String `tfsdk:"label"`
	Description   types.String `tfsdk:"description"`
	MinCPU        types.Int64  `tfsdk:"min_cpu"`
	MaxCPU        types.Int64  `tfsdk:"max_cpu"`
	MinMemoryMB   types.Int64  `tfsdk:"min_memory_mb"`
	MaxMemoryMB   types.Int64  `tfsdk:"max_memory_mb"`
	MinRootDiskGB types.Int64  `tfsdk:"min_root_disk_gb"`
	MaxRootDiskGB types.Int64  `tfsdk:"max_root_disk_gb"`
	Tier          types.String `tfsdk:"tier"`
}

func (d *packagesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_packages"
}

func (d *packagesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List available VM packages (sizing tiers).",
		Attributes: map[string]schema.Attribute{
			"packages": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The list of VM packages.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"label": schema.StringAttribute{
							Computed: true,
						},
						"description": schema.StringAttribute{
							Computed: true,
						},
						"min_cpu": schema.Int64Attribute{
							Computed: true,
						},
						"max_cpu": schema.Int64Attribute{
							Computed: true,
						},
						"min_memory_mb": schema.Int64Attribute{
							Computed: true,
						},
						"max_memory_mb": schema.Int64Attribute{
							Computed: true,
						},
						"min_root_disk_gb": schema.Int64Attribute{
							Computed: true,
						},
						"max_root_disk_gb": schema.Int64Attribute{
							Computed: true,
						},
						"tier": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *packagesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T.", req.ProviderData),
		)
		return
	}
	d.client = c
}

func (d *packagesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	packages, err := d.client.ListPackages(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Unable to list packages", err.Error())
		return
	}

	state := packagesDataSourceModel{
		Packages: make([]packageModel, len(packages)),
	}
	for i, p := range packages {
		state.Packages[i] = packageModel{
			ID:            types.StringValue(p.ID),
			Label:         types.StringValue(p.Label),
			Description:   types.StringValue(p.Description),
			MinCPU:        types.Int64Value(int64(p.MinCPU)),
			MaxCPU:        types.Int64Value(int64(p.MaxCPU)),
			MinMemoryMB:   types.Int64Value(int64(p.MinMemoryMB)),
			MaxMemoryMB:   types.Int64Value(int64(p.MaxMemoryMB)),
			MinRootDiskGB: types.Int64Value(int64(p.MinRootDiskGB)),
			MaxRootDiskGB: types.Int64Value(int64(p.MaxRootDiskGB)),
			Tier:          types.StringValue(p.Tier),
		}
	}

	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
