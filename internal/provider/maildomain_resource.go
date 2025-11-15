package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uberspace-community/terraform-provider-uberspace/gen/client"
	"github.com/uberspace-community/terraform-provider-uberspace/gen/provider/resource_maildomain"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource = &MaildomainResource{}
	// _ resource.ResourceWithImportState = &MaildomainResource{}.
)

func NewMaildomainResource() resource.Resource {
	return &MaildomainResource{}
}

// MaildomainResource defines the resource implementation.
type MaildomainResource struct {
	client *client.Client
}

func (r *MaildomainResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_maildomain"
}

func (r *MaildomainResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_maildomain.MaildomainResourceSchema(ctx)
}

func (r *MaildomainResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *MaildomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resource_maildomain.MaildomainModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := client.AsteroidsMaildomainsCreateApplicationJSON(client.MailDomainRequest{
		Name:     plan.Name.ValueString(),
		Asteroid: plan.Asteroid.ValueString(),
	})

	Maildomain, err := r.client.AsteroidsMaildomainsCreate(ctx, &apiReq, client.AsteroidsMaildomainsCreateParams{
		AsteroidName: plan.Asteroid.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create mail domain, got error: %s", err))
		return
	}

	plan.Asteroid = types.StringValue(Maildomain.Asteroid)
	plan.AsteroidName = types.StringValue(Maildomain.Asteroid)
	plan.CreatedAt = types.StringValue(Maildomain.CreatedAt.Format(time.RFC3339))
	plan.DnsState = types.StringValue(string(Maildomain.DNSState))
	plan.DnsLastCheck = types.StringValue(Maildomain.DNSLastCheck.Or(time.Now()).Format(time.RFC3339))
	plan.DnsError = types.StringValue(Maildomain.DNSError.Or(""))
	plan.Domain = types.StringValue(Maildomain.Domain)
	plan.DomainDisplay = types.StringValue(Maildomain.DomainDisplay)
	plan.DomainIdn = types.StringValue(Maildomain.DomainIdn)
	plan.Format = types.StringValue("json")
	plan.NameIdn = types.StringValue(Maildomain.Name)
	plan.NameDisplay = types.StringValue(Maildomain.NameDisplay)
	plan.NameIdn = types.StringValue(Maildomain.NameIdn)
	plan.UpdatedAt = types.StringValue(Maildomain.UpdatedAt.Format(time.RFC3339))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MaildomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resource_maildomain.MaildomainModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	Maildomain, err := r.client.AsteroidsMaildomainsGet(ctx, client.AsteroidsMaildomainsGetParams{
		AsteroidName: state.Asteroid.ValueString(),
		Name:         state.Name.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read mail domain, got error: %s", err))
		return
	}

	state.Asteroid = types.StringValue(Maildomain.Asteroid)
	state.AsteroidName = types.StringValue(Maildomain.Asteroid)
	state.CreatedAt = types.StringValue(Maildomain.CreatedAt.Format(time.RFC3339))
	state.DnsState = types.StringValue(string(Maildomain.DNSState))
	state.DnsLastCheck = types.StringValue(Maildomain.DNSLastCheck.Or(time.Now()).Format(time.RFC3339))
	state.DnsError = types.StringValue(Maildomain.DNSError.Or(""))
	state.Domain = types.StringValue(Maildomain.Domain)
	state.DomainDisplay = types.StringValue(Maildomain.DomainDisplay)
	state.DomainIdn = types.StringValue(Maildomain.DomainIdn)
	state.Format = types.StringValue("json")
	state.NameIdn = types.StringValue(Maildomain.Name)
	state.NameDisplay = types.StringValue(Maildomain.NameDisplay)
	state.NameIdn = types.StringValue(Maildomain.NameIdn)
	state.UpdatedAt = types.StringValue(Maildomain.UpdatedAt.Format(time.RFC3339))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *MaildomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan resource_maildomain.MaildomainModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.AsteroidsMaildomainsDelete(ctx, client.AsteroidsMaildomainsDeleteParams{
		AsteroidName: state.Asteroid.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete mail domain, got error: %s", err))
		return
	}

	apiReq := client.AsteroidsMaildomainsCreateApplicationJSON(client.MailDomainRequest{
		Name:     plan.Name.ValueString(),
		Asteroid: plan.Asteroid.ValueString(),
	})

	Maildomain, err := r.client.AsteroidsMaildomainsCreate(ctx, &apiReq, client.AsteroidsMaildomainsCreateParams{
		AsteroidName: plan.Asteroid.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create mail domain, got error: %s", err))
		return
	}

	plan.Asteroid = types.StringValue(Maildomain.Asteroid)
	plan.AsteroidName = types.StringValue(Maildomain.Asteroid)
	plan.CreatedAt = types.StringValue(Maildomain.CreatedAt.Format(time.RFC3339))
	plan.DnsState = types.StringValue(string(Maildomain.DNSState))
	plan.DnsLastCheck = types.StringValue(Maildomain.DNSLastCheck.Or(time.Now()).Format(time.RFC3339))
	plan.DnsError = types.StringValue(Maildomain.DNSError.Or(""))
	plan.Domain = types.StringValue(Maildomain.Domain)
	plan.DomainDisplay = types.StringValue(Maildomain.DomainDisplay)
	plan.DomainIdn = types.StringValue(Maildomain.DomainIdn)
	plan.Format = types.StringValue("json")
	plan.NameIdn = types.StringValue(Maildomain.Name)
	plan.NameDisplay = types.StringValue(Maildomain.NameDisplay)
	plan.NameIdn = types.StringValue(Maildomain.NameIdn)
	plan.UpdatedAt = types.StringValue(Maildomain.UpdatedAt.Format(time.RFC3339))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MaildomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resource_maildomain.MaildomainModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.AsteroidsMaildomainsDelete(ctx, client.AsteroidsMaildomainsDeleteParams{
		AsteroidName: state.Asteroid.ValueString(),
		Name:         state.Name.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete mail domain, got error: %s", err))
		return
	}
}
