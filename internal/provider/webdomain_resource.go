package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uberspace-community/terraform-provider-uberspace/gen/client"
	"github.com/uberspace-community/terraform-provider-uberspace/gen/provider/resource_webdomain"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource = &WebdomainResource{}
	// _ resource.ResourceWithImportState = &WebdomainResource{}.
)

func NewWebdomainResource() resource.Resource {
	return &WebdomainResource{}
}

// WebdomainResource defines the resource implementation.
type WebdomainResource struct {
	client *client.Client
}

func (r *WebdomainResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_webdomain"
}

func (r *WebdomainResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_webdomain.WebdomainResourceSchema(ctx)
}

func (r *WebdomainResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *uberspace.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *WebdomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resource_webdomain.WebdomainModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := client.CreateAsteroidsWebdomainsApplicationJSON(client.WebDomainRequest{
		Domain:   plan.Domain.ValueString(),
		Asteroid: plan.Asteroid.ValueString(),
	})

	Webdomain, err := r.client.CreateAsteroidsWebdomains(ctx, &apiReq, client.CreateAsteroidsWebdomainsParams{
		AsteroidName: plan.Asteroid.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create web domain, got error: %s", err))
		return
	}

	plan.AsteroidName = types.StringValue(Webdomain.Asteroid)
	plan.Asteroid = types.StringValue(Webdomain.Asteroid)
	plan.CreatedAt = types.StringValue(Webdomain.CreatedAt.Format(time.RFC3339))
	plan.Domain = types.StringValue(Webdomain.Domain)
	plan.DomainIdn = types.StringValue(Webdomain.DomainIdn)
	plan.UpdatedAt = types.StringValue(Webdomain.UpdatedAt.Format(time.RFC3339))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WebdomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resource_webdomain.WebdomainModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	Webdomain, err := r.client.GetAsteroidWebdomain(ctx, client.GetAsteroidWebdomainParams{
		AsteroidName: state.Asteroid.ValueString(),
		Domain:       state.Domain.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read web domain, got error: %s", err))
		return
	}

	state.AsteroidName = types.StringValue(Webdomain.Asteroid)
	state.Asteroid = types.StringValue(Webdomain.Asteroid)
	state.CreatedAt = types.StringValue(Webdomain.CreatedAt.Format(time.RFC3339))
	state.Domain = types.StringValue(Webdomain.Domain)
	state.DomainIdn = types.StringValue(Webdomain.DomainIdn)
	state.UpdatedAt = types.StringValue(Webdomain.UpdatedAt.Format(time.RFC3339))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WebdomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan resource_webdomain.WebdomainModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteAsteroidWebdomain(ctx, client.DeleteAsteroidWebdomainParams{
		AsteroidName: state.Asteroid.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete web domain, got error: %s", err))
		return
	}

	apiReq := client.CreateAsteroidsWebdomainsApplicationJSON(client.WebDomainRequest{
		Domain:   plan.Domain.ValueString(),
		Asteroid: plan.Asteroid.ValueString(),
	})

	Webdomain, err := r.client.CreateAsteroidsWebdomains(ctx, &apiReq, client.CreateAsteroidsWebdomainsParams{
		AsteroidName: plan.Asteroid.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create web domain, got error: %s", err))
		return
	}

	plan.AsteroidName = types.StringValue(Webdomain.Asteroid)
	plan.Asteroid = types.StringValue(Webdomain.Asteroid)
	plan.CreatedAt = types.StringValue(Webdomain.CreatedAt.Format(time.RFC3339))
	plan.Domain = types.StringValue(Webdomain.Domain)
	plan.DomainIdn = types.StringValue(Webdomain.DomainIdn)
	plan.UpdatedAt = types.StringValue(Webdomain.UpdatedAt.Format(time.RFC3339))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *WebdomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resource_webdomain.WebdomainModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteAsteroidWebdomain(ctx, client.DeleteAsteroidWebdomainParams{
		AsteroidName: state.Asteroid.ValueString(),
		Domain:       state.Domain.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete web domain, got error: %s", err))
		return
	}
}
