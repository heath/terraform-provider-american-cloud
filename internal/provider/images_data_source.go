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
	_ datasource.DataSource              = &imagesDataSource{}
	_ datasource.DataSourceWithConfigure = &imagesDataSource{}
)

func NewImagesDataSource() datasource.DataSource {
	return &imagesDataSource{}
}

type imagesDataSource struct {
	client *client.Client
}

type imagesDataSourceModel struct {
	OS      types.String `tfsdk:"os"`
	Version types.String `tfsdk:"version"`
	Images  []imageModel `tfsdk:"images"`
}

type imageModel struct {
	ID          types.String `tfsdk:"id"`
	Label       types.String `tfsdk:"label"`
	Description types.String `tfsdk:"description"`
	OS          types.String `tfsdk:"os"`
	Version     types.String `tfsdk:"version"`
}

func (d *imagesDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_images"
}

func (d *imagesDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List available OS images for VM provisioning.",
		Attributes: map[string]schema.Attribute{
			"os": schema.StringAttribute{
				Optional:    true,
				Description: "Filter images by operating system (e.g. \"ubuntu\", \"centos\").",
			},
			"version": schema.StringAttribute{
				Optional:    true,
				Description: "Filter images by OS version (used together with os).",
			},
			"images": schema.ListNestedAttribute{
				Computed:    true,
				Description: "The list of images.",
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
						"os": schema.StringAttribute{
							Computed: true,
						},
						"version": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *imagesDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *imagesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state imagesDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	osFilter := ""
	if !state.OS.IsNull() {
		osFilter = state.OS.ValueString()
	}
	versionFilter := ""
	if !state.Version.IsNull() {
		versionFilter = state.Version.ValueString()
	}

	images, err := d.client.ListImages(ctx, osFilter, versionFilter)
	if err != nil {
		resp.Diagnostics.AddError("Unable to list images", err.Error())
		return
	}

	state.Images = make([]imageModel, len(images))
	for i, img := range images {
		state.Images[i] = imageModel{
			ID:          types.StringValue(img.ID),
			Label:       types.StringValue(img.Label),
			Description: types.StringValue(img.Description),
			OS:          types.StringValue(img.OS),
			Version:     types.StringValue(img.Version),
		}
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
