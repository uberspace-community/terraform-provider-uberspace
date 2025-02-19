package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uberspace-community/terraform-provider-uberspace/uberspace"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &SupervisorServiceResource{}
	_ resource.ResourceWithImportState = &SupervisorServiceResource{}
)

func NewSupervisorServiceResource() resource.Resource {
	return &SupervisorServiceResource{}
}

// SupervisorServiceResource defines the resource implementation.
type SupervisorServiceResource struct {
	client *uberspace.Client
}

// SupervisorServiceResourceModel describes the resource data model.
type SupervisorServiceResourceModel struct {
	Name        types.String `tfsdk:"name"`
	Command     types.String `tfsdk:"command"`
	Environment types.Map    `tfsdk:"environment"`
}

func (r *SupervisorServiceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_supervisor_service"
}

func (r *SupervisorServiceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage supervisor services.",

		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the supervisor service.",
				Required:    true,
			},
			"command": schema.StringAttribute{
				Description: "The command to run.",
				Required:    true,
			},
			"environment": schema.MapAttribute{
				Description: "The environment variables to set.",
				Optional:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (r *SupervisorServiceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SupervisorServiceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state SupervisorServiceResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.SupervisorServiceCreate(
		ctx,
		r.client.User,
		state.Name.ValueString(),
		state.Command.ValueString(),
		stringMap(state.Environment),
	); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create supervisor service, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SupervisorServiceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state SupervisorServiceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	found, err := r.client.SupervisorServiceExists(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read supervisor service, got error: %s", err))
		return
	}

	if !found {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Supervisor service for %q not found", state.Name.ValueString()))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SupervisorServiceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, planning SupervisorServiceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planning)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.SupervisorServiceRemove(r.client.User, state.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update supervisor service, got error: %s", err))
		return
	}

	if err := r.client.SupervisorServiceCreate(
		ctx,
		r.client.User,
		planning.Name.ValueString(),
		planning.Command.ValueString(),
		stringMap(planning.Environment),
	); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update supervisor service, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &planning)...)
}

func (r *SupervisorServiceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state SupervisorServiceResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.SupervisorServiceRemove(r.client.User, state.Name.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete supervisor service, got error: %s", err))
		return
	}
}

func (r *SupervisorServiceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func stringMap(environment types.Map) map[string]string {
	env := make(map[string]string)
	for k, v := range environment.Elements() {
		env[k] = v.String()
	}

	return env
}
