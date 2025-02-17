package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/cugu/terraform-provider-uberspace/ssh"
	"github.com/cugu/terraform-provider-uberspace/uberspace"
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
	Host       string  `tfsdk:"host"`
	User       string  `tfsdk:"user"`
	Password   *string `tfsdk:"password"`
	PrivateKey *string `tfsdk:"private_key"`
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
				Required:    true,
			},
			"user": schema.StringAttribute{
				Description: "The user to authenticate with",
				Required:    true,
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

func (p *UberspaceProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) { //nolint:cyclop
	var data UberspaceProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	host := os.Getenv("UBERSPACE_HOST")
	if host != "" {
		data.Host = host
	}

	user := os.Getenv("UBERSPACE_USER")
	if user != "" {
		data.User = user
	}

	password := os.Getenv("UBERSPACE_PASSWORD")
	if password != "" {
		data.Password = &password
	}

	privateKey := os.Getenv("UBERSPACE_PRIVATE_KEY")
	if privateKey != "" {
		data.PrivateKey = &privateKey
	}

	if data.Password == nil && data.PrivateKey == nil {
		resp.Diagnostics.AddError("Invalid configuration", "either password or private_key must be set")
		return
	}

	config := &ssh.Config{
		Host: data.Host,
		User: data.User,
	}

	if data.PrivateKey != nil {
		config.PrivateKey = *data.PrivateKey
	}

	if data.Password != nil {
		config.Password = *data.Password
	}

	sshClient, err := ssh.NewClient(config)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create SSH client", err.Error())
		return
	}

	client := &uberspace.Client{
		User:   data.User,
		Runner: sshClient,
	}
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
