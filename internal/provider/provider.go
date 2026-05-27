package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-americancloud/internal/client"
)

var _ provider.Provider = &americancloudProvider{}

// New returns a provider.Provider constructor for use with providerserver.Serve.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &americancloudProvider{
			version: version,
		}
	}
}

type americancloudProvider struct {
	version string
}

type americancloudProviderModel struct {
	APIURL       types.String `tfsdk:"api_url"`
	BearerToken  types.String `tfsdk:"bearer_token"`
	ClientID     types.String `tfsdk:"client_id"`
	ClientSecret types.String `tfsdk:"client_secret"`
}

func (p *americancloudProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "americancloud"
	resp.Version = p.version
}

func (p *americancloudProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage American Cloud infrastructure resources.",
		Attributes: map[string]schema.Attribute{
			"api_url": schema.StringAttribute{
				Optional:    true,
				Description: "Base URL for the American Cloud API. Defaults to AC_API_URL env var, then https://api.americancloud.com.",
			},
			"bearer_token": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Bearer token for API authentication. Can also be set via AC_BEARER_TOKEN env var.",
			},
			"client_id": schema.StringAttribute{
				Optional:    true,
				Description: "Client ID for API Key authentication. Can also be set via AC_CLIENT_ID env var.",
			},
			"client_secret": schema.StringAttribute{
				Optional:    true,
				Sensitive:   true,
				Description: "Client secret for API Key authentication. Can also be set via AC_CLIENT_SECRET env var.",
			},
		},
	}
}

func (p *americancloudProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config americancloudProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Resolve values: config > env var > default
	apiURL := stringValueOrEnv(config.APIURL, "AC_API_URL", "")
	bearerToken := stringValueOrEnv(config.BearerToken, "AC_BEARER_TOKEN", "")
	clientID := stringValueOrEnv(config.ClientID, "AC_CLIENT_ID", "")
	clientSecret := stringValueOrEnv(config.ClientSecret, "AC_CLIENT_SECRET", "")

	// Validate auth configuration
	if bearerToken == "" && (clientID == "" || clientSecret == "") {
		resp.Diagnostics.AddAttributeError(
			path.Root("bearer_token"),
			"Missing API credentials",
			"Either bearer_token or both client_id and client_secret must be configured. "+
				"Set them in the provider block or via AC_BEARER_TOKEN / AC_CLIENT_ID + AC_CLIENT_SECRET environment variables.",
		)
		return
	}

	c, err := client.NewClient(apiURL, bearerToken, clientID, clientSecret)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create American Cloud client",
			err.Error(),
		)
		return
	}

	resp.DataSourceData = c
	resp.ResourceData = c
}

func (p *americancloudProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewImagesDataSource,
		NewRegionsDataSource,
		NewPackagesDataSource,
	}
}

func (p *americancloudProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewVMResource,
		NewIsolatedNetworkResource,
		NewVPCResource,
		NewVPCTierResource,
		NewPublicIPResource,
		NewFirewallRuleResource,
		NewPortForwardingResource,
		NewEgressRuleResource,
		NewLoadBalancerRuleResource,
	}
}

// stringValueOrEnv returns the Terraform config value if set, otherwise
// falls back to the named environment variable, then the default.
func stringValueOrEnv(val types.String, envVar, defaultVal string) string {
	if !val.IsNull() && !val.IsUnknown() {
		return val.ValueString()
	}
	if v := os.Getenv(envVar); v != "" {
		return v
	}
	return defaultVal
}
