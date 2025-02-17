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
	_ resource.Resource                = &WebBackendResource{}
	_ resource.ResourceWithImportState = &WebBackendResource{}
)

func NewWebBackendResource() resource.Resource {
	return &WebBackendResource{}
}

// WebBackendResource defines the resource implementation.
type WebBackendResource struct {
	client *uberspace.Client
}

// WebBackendResourceModel describes the resource data model.
type WebBackendResourceModel struct {
	URI  types.String `tfsdk:"uri"`
	Port types.Int32  `tfsdk:"port"`
}

func (r *WebBackendResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_web_backend"
}

func (r *WebBackendResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage web backends in web server configuration.",
		Attributes: map[string]schema.Attribute{
			"uri": schema.StringAttribute{
				Description: "The URI of the web backend.",
				Required:    true,
			},
			"port": schema.Int32Attribute{
				Description: "The port of the web backend.",
				Required:    true,
			},
		},
	}
}

func (r *WebBackendResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *WebBackendResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state WebBackendResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.WebBackendSet(ctx, state.URI.ValueString(), state.Port.ValueInt32()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create web backend, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WebBackendResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state WebBackendResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	found, err := r.client.WebBackendRead(ctx, state.URI.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read web backend, got error: %s", err))
		return
	}

	if !found {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Web backend for %q not found", state.URI.ValueString()))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *WebBackendResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, planning WebBackendResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planning)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.WebBackendSet(ctx, state.URI.ValueString(), planning.Port.ValueInt32())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update web backend, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &planning)...)
}

func (r *WebBackendResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state WebBackendResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	removed, err := r.client.WebBackendDelete(ctx, state.URI.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete web backend, got error: %s", err))
		return
	}

	if !removed {
		resp.Diagnostics.AddWarning("Not Found", fmt.Sprintf("Web backend for %q not found", state.URI.ValueString()))
		return
	}
}

func (r *WebBackendResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("uri"), req, resp)
}
