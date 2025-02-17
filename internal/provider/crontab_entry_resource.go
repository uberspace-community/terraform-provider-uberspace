package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/cugu/terraform-provider-uberspace/uberspace"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource                = &CronTabEntryResource{}
	_ resource.ResourceWithImportState = &CronTabEntryResource{}
)

func NewCronTabEntryResource() resource.Resource {
	return &CronTabEntryResource{}
}

// CronTabEntryResource defines the resource implementation.
type CronTabEntryResource struct {
	client *uberspace.Client
}

// CronTabEntryResourceModel describes the resource data model.
type CronTabEntryResourceModel struct {
	Entry types.String `tfsdk:"entry"`
}

func (r *CronTabEntryResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_crontab_entry"
}

func (r *CronTabEntryResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Crontab entry",

		Attributes: map[string]schema.Attribute{
			"entry": schema.StringAttribute{
				Description: "Crontab entry, e.g. `* * * * * echo hello`",
				Required:    true,
			},
		},
	}
}

func (r *CronTabEntryResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *CronTabEntryResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var planning CronTabEntryResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &planning)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.CrontabEntryAdd(ctx, planning.Entry.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create crontab entry, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &planning)...)
}

func (r *CronTabEntryResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state CronTabEntryResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	found, err := r.client.CrontabEntryExists(ctx, state.Entry.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read crontab entry, got error: %s", err))
		return
	}

	if !found {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Crontab entry for %q not found", state.Entry.ValueString()))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *CronTabEntryResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, planning CronTabEntryResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planning)...)

	if resp.Diagnostics.HasError() {
		return
	}

	removed, err := r.client.CrontabEntryRemove(ctx, state.Entry.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update crontab entry, got error: %s", err))
		return
	}

	if !removed {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Unable to update crontab entry, got error: %s", err))
		return
	}

	err = r.client.CrontabEntryAdd(ctx, planning.Entry.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update crontab entry, got error: %s", err))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &planning)...)
}

func (r *CronTabEntryResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state CronTabEntryResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	removed, err := r.client.CrontabEntryRemove(ctx, state.Entry.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete crontab entry, got error: %s", err))
		return
	}

	if !removed {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Unable to delete crontab entry, got error: %s", err))
		return
	}
}

func (r *CronTabEntryResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("entry"), req, resp)
}
