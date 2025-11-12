package provider

import (
	"cmp"
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uberspace-community/terraform-provider-uberspace/gen/client"
)

// Ensure UberspaceProvider satisfies various provider interfaces.
var _ provider.Provider = &UberspaceProvider{}

// UberspaceProvider defines the provider implementation.
type UberspaceProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// UberspaceProviderModel describes the provider data model.
type UberspaceProviderModel struct {
	APIKey types.String `tfsdk:"apikey"`
}

func (p *UberspaceProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "uberspace"
	resp.Version = p.version
}

func (p *UberspaceProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"apikey": schema.StringAttribute{
				Description: "The API key for the Uberspace API. If not set, the environment variable UBERSPACE_APIKEY will be used.",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *UberspaceProvider) ValidateConfig(ctx context.Context, req provider.ValidateConfigRequest, resp *provider.ValidateConfigResponse) {
	var data UberspaceProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !data.APIKey.IsUnknown() && data.APIKey.ValueString() == "" && os.Getenv("UBERSPACE_APIKEY") == "" {
		resp.Diagnostics.AddError("Invalid configuration", "apikey or UBERSPACE_APIKEY must be set")
	}
}

func (p *UberspaceProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data UberspaceProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apikey := cmp.Or(data.APIKey.ValueString(), os.Getenv("UBERSPACE_APIKEY"))

	if apikey == "" {
		resp.Diagnostics.AddError(
			"Invalid configuration",
			"apikey or UBERSPACE_APIKEY must be set",
		)

		return
	}

	client, err := client.NewClient("https://marvin.uberspace.is", client.WithClient(newAuthClient(apikey)))
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to create Uberspace client",
			"An error occurred while creating the Uberspace client: "+err.Error(),
		)

		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *UberspaceProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewWebdomainResource,
		NewSshkeyResource,
		NewWebdomainBackendResource,
		NewWebdomainHeaderResource,
		NewMaildomainResource,
		NewMailuserResource,
	}
}

func (p *UberspaceProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &UberspaceProvider{
			version: version,
		}
	}
}
