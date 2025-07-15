package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uberspace-community/terraform-provider-uberspace/gen/client"
	"github.com/uberspace-community/terraform-provider-uberspace/gen/provider/resource_webdomain_backend"
)

var _ resource.Resource = &WebdomainBackendResource{}

func NewWebdomainBackendResource() resource.Resource {
	return &WebdomainBackendResource{}
}

// WebdomainBackendResource defines the resource implementation.
type WebdomainBackendResource struct {
	client *client.Client
}

func (r *WebdomainBackendResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webdomain_backend"
}

func (r *WebdomainBackendResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_webdomain_backend.WebdomainBackendResourceSchema(ctx)
}

func (r *WebdomainBackendResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	c, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *uberspace.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = c
}

func (r *WebdomainBackendResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resource_webdomain_backend.WebdomainBackendModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var port client.OptNilInt
	if !plan.Port.IsNull() {
		port = client.NewOptNilInt(int(plan.Port.ValueInt64()))
	}

	reqBody := client.WebBackendRequest{
		Asteroid:     plan.Asteroid.ValueString(),
		Domain:       client.NewNilString(plan.Domain.ValueString()),
		Path:         plan.Path.ValueString(),
		RemovePrefix: client.NewOptBool(plan.RemovePrefix.ValueBool()),
		Destination:  client.DestinationEnum(plan.Destination.ValueString()),
		Port:         port,
	}

	apiReq := client.CreateAsteroidsWebdomainsBackendsApplicationJSON(reqBody)

	backend, err := r.client.CreateAsteroidsWebdomainsBackends(ctx, &apiReq, client.CreateAsteroidsWebdomainsBackendsParams{
		AsteroidName:    plan.AsteroidName.ValueString(),
		WebdomainDomain: plan.WebdomainDomain.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create web domain backend, got error: %s", err))
		return
	}

	plan.Pk = types.Int64Value(int64(backend.Pk))
	plan.CreatedAt = types.StringValue(backend.CreatedAt.Format(time.RFC3339))
	plan.UpdatedAt = types.StringValue(backend.UpdatedAt.Format(time.RFC3339))
	plan.Asteroid = types.StringValue(backend.Asteroid)
	plan.Domain = types.StringValue(backend.Domain.Or(""))
	plan.Path = types.StringValue(backend.Path)
	plan.RemovePrefix = types.BoolValue(backend.RemovePrefix.Or(false))

	plan.Destination = types.StringValue(string(backend.Destination))
	if v, ok := backend.Port.Get(); ok {
		plan.Port = types.Int64Value(int64(v))
	} else if backend.Port.IsNull() {
		plan.Port = types.Int64Null()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WebdomainBackendResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resource_webdomain_backend.WebdomainBackendModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	backend, err := r.client.GetAsteroidWebdomainBackend(ctx, client.GetAsteroidWebdomainBackendParams{
		AsteroidName:    state.AsteroidName.ValueString(),
		WebdomainDomain: state.WebdomainDomain.ValueString(),
		Path:            state.Path.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read web domain backend, got error: %s", err))
		return
	}

	state.Pk = types.Int64Value(int64(backend.Pk))
	state.CreatedAt = types.StringValue(backend.CreatedAt.Format(time.RFC3339))
	state.UpdatedAt = types.StringValue(backend.UpdatedAt.Format(time.RFC3339))
	state.Asteroid = types.StringValue(backend.Asteroid)
	state.Domain = types.StringValue(backend.Domain.Or(""))
	state.Path = types.StringValue(backend.Path)
	state.RemovePrefix = types.BoolValue(backend.RemovePrefix.Or(false))

	state.Destination = types.StringValue(string(backend.Destination))
	if v, ok := backend.Port.Get(); ok {
		state.Port = types.Int64Value(int64(v))
	} else if backend.Port.IsNull() {
		state.Port = types.Int64Null()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WebdomainBackendResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan resource_webdomain_backend.WebdomainBackendModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteAsteroidWebdomainBackend(ctx, client.DeleteAsteroidWebdomainBackendParams{
		AsteroidName:    state.AsteroidName.ValueString(),
		WebdomainDomain: state.WebdomainDomain.ValueString(),
		Path:            state.Path.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete web domain backend, got error: %s", err))
		return
	}

	var port client.OptNilInt
	if !plan.Port.IsNull() {
		port = client.NewOptNilInt(int(plan.Port.ValueInt64()))
	}

	reqBody := client.WebBackendRequest{
		Asteroid:     plan.Asteroid.ValueString(),
		Domain:       client.NewNilString(plan.Domain.ValueString()),
		Path:         plan.Path.ValueString(),
		RemovePrefix: client.NewOptBool(plan.RemovePrefix.ValueBool()),
		Destination:  client.DestinationEnum(plan.Destination.ValueString()),
		Port:         port,
	}

	apiReq := client.CreateAsteroidsWebdomainsBackendsApplicationJSON(reqBody)

	backend, err := r.client.CreateAsteroidsWebdomainsBackends(ctx, &apiReq, client.CreateAsteroidsWebdomainsBackendsParams{
		AsteroidName:    plan.AsteroidName.ValueString(),
		WebdomainDomain: plan.WebdomainDomain.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create web domain backend, got error: %s", err))
		return
	}

	plan.Pk = types.Int64Value(int64(backend.Pk))
	plan.CreatedAt = types.StringValue(backend.CreatedAt.Format(time.RFC3339))
	plan.UpdatedAt = types.StringValue(backend.UpdatedAt.Format(time.RFC3339))
	plan.Asteroid = types.StringValue(backend.Asteroid)
	plan.Domain = types.StringValue(backend.Domain.Or(""))
	plan.Path = types.StringValue(backend.Path)
	plan.RemovePrefix = types.BoolValue(backend.RemovePrefix.Or(false))

	plan.Destination = types.StringValue(string(backend.Destination))
	if v, ok := backend.Port.Get(); ok {
		plan.Port = types.Int64Value(int64(v))
	} else if backend.Port.IsNull() {
		plan.Port = types.Int64Null()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WebdomainBackendResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resource_webdomain_backend.WebdomainBackendModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteAsteroidWebdomainBackend(ctx, client.DeleteAsteroidWebdomainBackendParams{
		AsteroidName:    state.AsteroidName.ValueString(),
		WebdomainDomain: state.WebdomainDomain.ValueString(),
		Path:            state.Path.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete web domain backend, got error: %s", err))
		return
	}
}
