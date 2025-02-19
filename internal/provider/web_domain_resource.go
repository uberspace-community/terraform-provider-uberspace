package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uberspace-community/terraform-provider-uberspace/uberspace"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource = &WebDomainResource{}
	// _ resource.ResourceWithImportState = &WebDomainResource{}.
)

func NewWebDomainResource() resource.Resource {
	return &WebDomainResource{}
}

// WebDomainResource defines the resource implementation.
type WebDomainResource struct {
	client *uberspace.Client
}

// WebDomainResourceModel describes the resource data model.
type WebDomainResourceModel struct {
	Domain types.String `tfsdk:"domain"`
}

func (r *WebDomainResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_web_domain"
}

func (r *WebDomainResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage domains in web server configuration.",

		Attributes: map[string]schema.Attribute{
			"domain": schema.StringAttribute{
				Description: "The domain name to manage.",
				Required:    true,
			},
		},
	}
}

func (r *WebDomainResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *WebDomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planning WebDomainResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planning)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.WebDomainAdd(planning.Domain.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create web domain, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &planning)...)
}

func (r *WebDomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WebDomainResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	found, err := r.client.WebDomainExists(state.Domain.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read web domain, got error: %s", err))
		return
	}

	if !found {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Web domain for %q not found", state.Domain.ValueString()))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WebDomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, planning WebDomainResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planning)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.WebDomainDelete(state.Domain.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update web domain, got error: %s", err))
		return
	}

	if err := r.client.WebDomainAdd(planning.Domain.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update web domain, got error: %s", err))
		return
	}

	planning.Domain = types.StringValue(planning.Domain.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &planning)...)
}

func (r *WebDomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state WebDomainResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.WebDomainDelete(state.Domain.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete web domain, got error: %s", err))
		return
	}
}

// func (r *WebDomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
// 	resource.ImportStatePassthroughID(ctx, path.Root("domain"), req, resp)
// }
