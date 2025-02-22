package provider

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uberspace-community/terraform-provider-uberspace/uberspace"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource = &RemoteFileResource{}
	// _ resource.ResourceWithImportState    = &RemoteFileResource{}.
	_ resource.ResourceWithValidateConfig = &RemoteFileResource{}
)

func NewRemoteFileResource() resource.Resource {
	return &RemoteFileResource{}
}

// RemoteFileResource defines the resource implementation.
type RemoteFileResource struct {
	client *uberspace.Client
}

// RemoteFileResourceModel describes the resource data model.
type RemoteFileResourceModel struct {
	Src        types.String `tfsdk:"src"`
	SrcHash    types.String `tfsdk:"src_hash"`
	Content    types.String `tfsdk:"content"`
	Dst        types.String `tfsdk:"dst"`
	Executable types.Bool   `tfsdk:"executable"`
}

func (r *RemoteFileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_remote_file"
}

func (r *RemoteFileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage remote files.",

		Attributes: map[string]schema.Attribute{
			"src": schema.StringAttribute{
				Description: "The local file path or a http(s) URL to fetch the file from. Either this or content must be set.",
				Optional:    true,
			},
			"src_hash": schema.StringAttribute{
				Description: "The hash of the source file, used to detect changes. Required if src is set.",
				Optional:    true,
			},
			"content": schema.StringAttribute{
				Description: "The content of the file to create. Either this or src must be set.",
				Optional:    true,
			},
			"dst": schema.StringAttribute{
				Description: "The destination file.",
				Required:    true,
			},
			"executable": schema.BoolAttribute{
				Description: "Whether the destination file should be executable.",
				Optional:    true,
			},
		},
	}
}

func (r *RemoteFileResource) ValidateConfig(ctx context.Context, request resource.ValidateConfigRequest, response *resource.ValidateConfigResponse) { //nolint:cyclop
	var model RemoteFileResourceModel

	response.Diagnostics.Append(request.Config.Get(ctx, &model)...)

	if !model.Src.IsUnknown() { //nolint:nestif
		hasSrc := model.Src.ValueString() != ""

		if hasSrc {
			src := model.Src.ValueString()

			if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
				response.Diagnostics.AddError("Not Implemented", "http(s) URLs are not yet supported")
			} else {
				info, err := os.Stat(src)
				if err != nil {
					response.Diagnostics.AddAttributeError(path.Root("src"), "Invalid Configuration", fmt.Sprintf("Unable to find src file, got error: %s", err))
				}

				if info.IsDir() {
					response.Diagnostics.AddAttributeError(path.Root("src"), "Invalid Configuration", "src must be a file, not a directory")
				}
			}
		}

		if !model.Content.IsUnknown() {
			hasContent := model.Content.ValueString() != ""

			if !hasSrc && !hasContent {
				response.Diagnostics.AddError("Invalid Configuration", "Either src or content must be set")
			}

			if hasSrc && hasContent {
				response.Diagnostics.AddError("Invalid Configuration", "src and content cannot be set at the same time")
			}
		}

		if !model.SrcHash.IsUnknown() && hasSrc && model.SrcHash.ValueString() == "" {
			response.Diagnostics.AddAttributeError(path.Root("src_hash"), "Invalid Configuration", "src_hash must be set if src is set")
		}
	}
}

func (r *RemoteFileResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *RemoteFileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state RemoteFileResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if state.Src.ValueString() != "" {
		if err := r.client.RemoteFileCopy(
			ctx,
			state.Src.ValueString(),
			state.Dst.ValueString(),
			state.Executable.ValueBool(),
		); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to copy remote file, got error: %s", err))
			return
		}
	}

	if state.Content.ValueString() != "" {
		if err := r.client.RemoteFileCreate(
			ctx,
			[]byte(state.Content.ValueString()),
			state.Dst.ValueString(),
			state.Executable.ValueBool(),
		); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create remote file, got error: %s", err))
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *RemoteFileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state RemoteFileResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	found, err := r.client.RemoteFileExists(state.Dst.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read remote file, got error: %s", err))
		return
	}

	if !found {
		resp.Diagnostics.AddError("Not Found", fmt.Sprintf("Remote file for %q not found", state.Dst.ValueString()))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *RemoteFileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, planning RemoteFileResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &planning)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// srcChanged := state.Src.ValueString() != planning.Src.ValueString()
	// srcHashChanged := state.SrcHash.ValueString() != planning.SrcHash.ValueString()
	// contentChanged := state.Content.ValueString() != planning.Content.ValueString()
	dstChanged := state.Dst.ValueString() != planning.Dst.ValueString()
	// executableChanged := state.Executable.ValueBool() != planning.Executable.ValueBool()

	if dstChanged {
		if err := r.client.RemoteFileDelete(state.Dst.ValueString()); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update remote file, got error: %s", err))
			return
		}
	}

	if state.Src.ValueString() != "" {
		if err := r.client.RemoteFileCopy(
			ctx,
			planning.Src.ValueString(),
			planning.Dst.ValueString(),
			planning.Executable.ValueBool(),
		); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update remote file, got error: %s", err))
			return
		}
	}

	if state.Content.ValueString() != "" {
		if err := r.client.RemoteFileCreate(
			ctx,
			[]byte(planning.Content.ValueString()),
			planning.Dst.ValueString(),
			planning.Executable.ValueBool(),
		); err != nil {
			resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to update remote file, got error: %s", err))
			return
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &planning)...)
}

func (r *RemoteFileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state RemoteFileResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.RemoteFileDelete(state.Dst.ValueString()); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete remote file, got error: %s", err))
		return
	}
}

// func (r *RemoteFileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
// 	resource.ImportStatePassthroughID(ctx, path.Root("dst"), req, resp)
// }
