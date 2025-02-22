package provider

import (
	"cmp"
	"context"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/uberspace-community/terraform-provider-uberspace/ssh"
	"github.com/uberspace-community/terraform-provider-uberspace/uberspace"
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
	Host       types.String `tfsdk:"host"`
	User       types.String `tfsdk:"user"`
	Password   types.String `tfsdk:"password"`
	PrivateKey types.String `tfsdk:"private_key"`
}

func (p *UberspaceProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "uberspace"
	resp.Version = p.version
}

func (p *UberspaceProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Description: "The hostname of the SSH server",
				Optional:    true,
			},
			"user": schema.StringAttribute{
				Description: "The user to authenticate with",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "The password to authenticate with, either this or private_key must be set",
				Optional:    true,
				Sensitive:   true,
			},
			"private_key": schema.StringAttribute{
				Description: "The private key to authenticate with, either this or password must be set",
				Optional:    true,
				Sensitive:   true,
			},
		},
	}
}

func (p *UberspaceProvider) ValidateConfig(ctx context.Context, req provider.ValidateConfigRequest, resp *provider.ValidateConfigResponse) { //nolint:cyclop
	var data UberspaceProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Host.ValueString() == "" && os.Getenv("UBERSPACE_HOST") == "" {
		resp.Diagnostics.AddError("Invalid configuration", "host or UBERSPACE_HOST must be set")
	}

	if data.User.ValueString() == "" && os.Getenv("UBERSPACE_USER") == "" {
		resp.Diagnostics.AddError("Invalid configuration", "user or UBERSPACE_USER must be set")
	}

	if data.Password.ValueString() == "" && data.PrivateKey.ValueString() == "" && os.Getenv("UBERSPACE_PASSWORD") == "" && os.Getenv("UBERSPACE_PRIVATE_KEY") == "" {
		resp.Diagnostics.AddError("Invalid configuration", "password, private_key, UBERSPACE_PASSWORD or UBERSPACE_PRIVATE_KEY must be set")
	}

	if data.Password.ValueString() != "" && data.PrivateKey.ValueString() != "" {
		resp.Diagnostics.AddError("Invalid configuration", "only one of password or private_key must be set")
	}
}

func (p *UberspaceProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data UberspaceProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	user := cmp.Or(data.User.ValueString(), os.Getenv("UBERSPACE_USER"))

	sshClient, err := ssh.NewClient(&ssh.Config{
		Host:       cmp.Or(data.Host.ValueString(), os.Getenv("UBERSPACE_HOST")),
		User:       user,
		Password:   cmp.Or(data.Password.ValueString(), os.Getenv("UBERSPACE_PASSWORD")),
		PrivateKey: cmp.Or(data.PrivateKey.ValueString(), os.Getenv("UBERSPACE_PRIVATE_KEY")),
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to create SSH client", err.Error())
		return
	}

	client := &uberspace.Client{User: user, SSHClient: sshClient}
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *UberspaceProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewCronTabEntryResource,
		NewWebDomainResource,
		NewWebBackendResource,
		NewMySQLDatabaseResource,
		NewSupervisorServiceResource,
		NewRemoteFileResource,
	}
}

func (p *UberspaceProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewUserDataSource,
		NewMyCnfDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &UberspaceProvider{
			version: version,
		}
	}
}
