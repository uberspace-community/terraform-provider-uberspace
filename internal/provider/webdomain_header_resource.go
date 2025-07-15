package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uberspace-community/terraform-provider-uberspace/gen/client"
	"github.com/uberspace-community/terraform-provider-uberspace/gen/provider/resource_webdomain_header"
)

var _ resource.Resource = &WebdomainHeaderResource{}

func NewWebdomainHeaderResource() resource.Resource {
	return &WebdomainHeaderResource{}
}

// WebdomainHeaderResource defines the resource implementation.
type WebdomainHeaderResource struct {
	client *client.Client
}

func (r *WebdomainHeaderResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webdomain_header"
}

func (r *WebdomainHeaderResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_webdomain_header.WebdomainHeaderResourceSchema(ctx)
}

func (r *WebdomainHeaderResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *WebdomainHeaderResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resource_webdomain_header.WebdomainHeaderModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var value client.OptNilString
	if !plan.Value.IsNull() {
		value = client.NewOptNilString(plan.Value.ValueString())
	}

	reqBody := client.WebHeaderRequest{
		Asteroid: plan.Asteroid.ValueString(),
		Domain:   client.NewNilString(plan.Domain.ValueString()),
		Path:     plan.Path.ValueString(),
		Name:     plan.Name.ValueString(),
		Value:    value,
	}

	apiReq := client.CreateAsteroidsWebdomainsHeadersApplicationJSON(reqBody)

	header, err := r.client.CreateAsteroidsWebdomainsHeaders(ctx, &apiReq, client.CreateAsteroidsWebdomainsHeadersParams{
		AsteroidName:    plan.AsteroidName.ValueString(),
		WebdomainDomain: plan.WebdomainDomain.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create web domain header, got error: %s", err))
		return
	}

	plan.CreatedAt = types.StringValue(header.CreatedAt.Format(time.RFC3339))
	plan.UpdatedAt = types.StringValue(header.UpdatedAt.Format(time.RFC3339))
	plan.Asteroid = types.StringValue(header.Asteroid)
	plan.Domain = types.StringValue(header.Domain.Or(""))
	plan.Path = types.StringValue(header.Path)

	plan.Name = types.StringValue(header.Name)
	if v, ok := header.Value.Get(); ok {
		plan.Value = types.StringValue(v)
	} else if header.Value.IsNull() {
		plan.Value = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WebdomainHeaderResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resource_webdomain_header.WebdomainHeaderModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	header, err := r.client.GetAsteroidWebdomainHeader(ctx, client.GetAsteroidWebdomainHeaderParams{
		AsteroidName:    state.AsteroidName.ValueString(),
		WebdomainDomain: state.WebdomainDomain.ValueString(),
		Path:            state.Path.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read web domain header, got error: %s", err))
		return
	}

	state.CreatedAt = types.StringValue(header.CreatedAt.Format(time.RFC3339))
	state.UpdatedAt = types.StringValue(header.UpdatedAt.Format(time.RFC3339))
	state.Asteroid = types.StringValue(header.Asteroid)
	state.Domain = types.StringValue(header.Domain.Or(""))
	state.Path = types.StringValue(header.Path)

	state.Name = types.StringValue(header.Name)
	if v, ok := header.Value.Get(); ok {
		state.Value = types.StringValue(v)
	} else if header.Value.IsNull() {
		state.Value = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WebdomainHeaderResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan resource_webdomain_header.WebdomainHeaderModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteAsteroidWebdomainHeader(ctx, client.DeleteAsteroidWebdomainHeaderParams{
		AsteroidName:    state.AsteroidName.ValueString(),
		WebdomainDomain: state.WebdomainDomain.ValueString(),
		Path:            state.Path.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete web domain header, got error: %s", err))
		return
	}

	var value client.OptNilString
	if !plan.Value.IsNull() {
		value = client.NewOptNilString(plan.Value.ValueString())
	}

	reqBody := client.WebHeaderRequest{
		Asteroid: plan.Asteroid.ValueString(),
		Domain:   client.NewNilString(plan.Domain.ValueString()),
		Path:     plan.Path.ValueString(),
		Name:     plan.Name.ValueString(),
		Value:    value,
	}

	apiReq := client.CreateAsteroidsWebdomainsHeadersApplicationJSON(reqBody)

	header, err := r.client.CreateAsteroidsWebdomainsHeaders(ctx, &apiReq, client.CreateAsteroidsWebdomainsHeadersParams{
		AsteroidName:    plan.AsteroidName.ValueString(),
		WebdomainDomain: plan.WebdomainDomain.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create web domain header, got error: %s", err))
		return
	}

	plan.CreatedAt = types.StringValue(header.CreatedAt.Format(time.RFC3339))
	plan.UpdatedAt = types.StringValue(header.UpdatedAt.Format(time.RFC3339))
	plan.Asteroid = types.StringValue(header.Asteroid)
	plan.Domain = types.StringValue(header.Domain.Or(""))
	plan.Path = types.StringValue(header.Path)

	plan.Name = types.StringValue(header.Name)
	if v, ok := header.Value.Get(); ok {
		plan.Value = types.StringValue(v)
	} else if header.Value.IsNull() {
		plan.Value = types.StringNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WebdomainHeaderResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resource_webdomain_header.WebdomainHeaderModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteAsteroidWebdomainHeader(ctx, client.DeleteAsteroidWebdomainHeaderParams{
		AsteroidName:    state.AsteroidName.ValueString(),
		WebdomainDomain: state.WebdomainDomain.ValueString(),
		Path:            state.Path.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete web domain header, got error: %s", err))
		return
	}
}
