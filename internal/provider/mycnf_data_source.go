package provider

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"gopkg.in/ini.v1"

	"github.com/cugu/terraform-provider-uberspace/uberspace"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &MyCnfDataSource{}

func NewMyCnfDataSource() datasource.DataSource {
	return &MyCnfDataSource{}
}

// MyCnfDataSource defines the data source implementation.
type MyCnfDataSource struct {
	client *uberspace.Client
}

// MyCnfDataSourceModel describes the data source data model.
type MyCnfDataSourceModel struct {
	Client         *MyCnfClientDataSourceModel `tfsdk:"client"`
	ClientReadOnly *MyCnfClientDataSourceModel `tfsdk:"clientreadonly"`
}

type MyCnfClientDataSourceModel struct {
	User     types.String `tfsdk:"user"`
	Password types.String `tfsdk:"password"`
}

func (d *MyCnfDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mycnf"
}

func (d *MyCnfDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "MySQL configuration data source",
		Attributes: map[string]schema.Attribute{
			"client": schema.SingleNestedAttribute{
				Description: "Client configuration",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"user": schema.StringAttribute{
						Description: "MySQL user",
						Computed:    true,
					},
					"password": schema.StringAttribute{
						Description: "MySQL password",
						Computed:    true,
						Sensitive:   true,
					},
				},
			},
			"clientreadonly": schema.SingleNestedAttribute{
				Description: "Client configuration",
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"user": schema.StringAttribute{
						Description: "MySQL user",
						Computed:    true,
					},
					"password": schema.StringAttribute{
						Description: "MySQL password",
						Computed:    true,
						Sensitive:   true,
					},
				},
			},
		},
	}
}

func (d *MyCnfDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*uberspace.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *uberspace.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.client = client
}

func (d *MyCnfDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var config MyCnfDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	cmd := exec.CommandContext(ctx, "cat", "~/.my.cnf")

	out, err := d.client.Runner.Run(cmd)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read MySQL configuration", err.Error())
		return
	}

	cfg, err := ini.Load(out)
	if err != nil {
		resp.Diagnostics.AddError("Failed to parse MySQL configuration", err.Error())
		return
	}

	config.Client = &MyCnfClientDataSourceModel{
		User:     types.StringValue(cfg.Section("client").Key("user").String()),
		Password: types.StringValue(cfg.Section("client").Key("password").String()),
	}
	config.ClientReadOnly = &MyCnfClientDataSourceModel{
		User:     types.StringValue(cfg.Section("clientreadonly").Key("user").String()),
		Password: types.StringValue(cfg.Section("clientreadonly").Key("password").String()),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &config)...)
}
