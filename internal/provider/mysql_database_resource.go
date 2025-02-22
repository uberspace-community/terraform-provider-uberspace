package provider

import (
	"context"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uberspace-community/terraform-provider-uberspace/uberspace"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource = &MySQLDatabaseResource{}
	// _ resource.ResourceWithImportState    = &MySQLDatabaseResource{}.
	_ resource.ResourceWithValidateConfig = &MySQLDatabaseResource{}
)

func NewMySQLDatabaseResource() resource.Resource {
	return &MySQLDatabaseResource{}
}

// MySQLDatabaseResource defines the resource implementation.
type MySQLDatabaseResource struct {
	client *uberspace.Client
}

// MySQLDatabaseResourceModel describes the resource data model.
type MySQLDatabaseResourceModel struct {
	Suffix types.String `tfsdk:"suffix"`
}

func (r *MySQLDatabaseResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mysql_database"
}

func (r *MySQLDatabaseResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage MySQL databases.",

		Attributes: map[string]schema.Attribute{
			"suffix": schema.StringAttribute{
				Description: "The suffix of the MySQL database, all databases will be prefixed with the user name.",
				Required:    true,
			},
		},
	}
}

var suffixRegex = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)

func (r *MySQLDatabaseResource) ValidateConfig(ctx context.Context, request resource.ValidateConfigRequest, response *resource.ValidateConfigResponse) {
	var model MySQLDatabaseResourceModel

	response.Diagnostics.Append(request.Config.Get(ctx, &model)...)

	if response.Diagnostics.HasError() {
		return
	}

	if !model.Suffix.IsUnknown() && !suffixRegex.MatchString(model.Suffix.ValueString()) {
		response.Diagnostics.AddAttributeError(path.Root("suffix"), "Invalid Suffix", "Suffix must match the regex: "+suffixRegex.String())
	}
}

func (r *MySQLDatabaseResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*uberspace.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *uberspace.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *MySQLDatabaseResource) databaseName(suffix string) string {
	return r.client.User + "_" + suffix
}

func (r *MySQLDatabaseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state MySQLDatabaseResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.MySQLDatabaseCreate(r.databaseName(state.Suffix.ValueString())); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create mysql database, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *MySQLDatabaseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state MySQLDatabaseResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	stateDatabaseName := r.databaseName(state.Suffix.ValueString())

	found, err := r.client.MySQLDatabaseExists(stateDatabaseName)
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read mysql database, got error: %s", err))
		return
	}

	if !found {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Mysql database for %q not found", stateDatabaseName))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *MySQLDatabaseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, planning MySQLDatabaseResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planning)...)

	if resp.Diagnostics.HasError() {
		return
	}

	stateDatabaseName := r.databaseName(state.Suffix.ValueString())
	planningDatabaseName := r.databaseName(planning.Suffix.ValueString())

	if err := r.client.MySQLDatabaseDrop(stateDatabaseName); err != nil {
		resp.Diagnostics.AddWarning("Client Error", fmt.Sprintf("Unable to update mysql database, got error: %s", err))
	}

	if err := r.client.MySQLDatabaseCreate(planningDatabaseName); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update mysql database, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &planning)...)
}

func (r *MySQLDatabaseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state MySQLDatabaseResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	stateDatabaseName := r.databaseName(state.Suffix.ValueString())

	if err := r.client.MySQLDatabaseDrop(stateDatabaseName); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete mysql database, got error: %s", err))
		return
	}
}

// func (r *MySQLDatabaseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
// 	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
// }
