package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uberspace-community/terraform-provider-uberspace/gen/client"
	"github.com/uberspace-community/terraform-provider-uberspace/gen/provider/resource_sshkey"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &SshkeyResource{}

// NewSshkeyResource returns a new resource instance.
func NewSshkeyResource() resource.Resource {
	return &SshkeyResource{}
}

// SshkeyResource defines the resource implementation.
type SshkeyResource struct {
	client *client.Client
}

func (r *SshkeyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sshkey"
}

func (r *SshkeyResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_sshkey.SshkeyResourceSchema(ctx)
}

func (r *SshkeyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SshkeyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resource_sshkey.SshkeyModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	reqBody := client.SshKeyRequest{
		Asteroid: plan.Asteroid.ValueString(),
		Key:      plan.Key.ValueString(),
		KeyType:  client.KeyTypeEnum(plan.KeyType.ValueString()),
	}
	if !plan.KeyComment.IsNull() {
		reqBody.KeyComment.SetTo(plan.KeyComment.ValueString())
	}

	apiReq := client.AsteroidsSshkeysCreateApplicationJSON(reqBody)

	sshKey, err := r.client.AsteroidsSshkeysCreate(ctx, &apiReq, client.AsteroidsSshkeysCreateParams{
		AsteroidName: plan.Asteroid.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create ssh key, got error: %s", err))
		return
	}

	plan.Pk = types.Int64Value(int64(sshKey.Pk))
	plan.Id = types.Int64Value(int64(sshKey.Pk))
	plan.FormattedKey = types.StringValue(sshKey.FormattedKey)
	plan.ShortenedKey = types.StringValue(sshKey.ShortenedKey)
	plan.CreatedAt = types.StringValue(sshKey.CreatedAt.Format(time.RFC3339))
	plan.UpdatedAt = types.StringValue(sshKey.UpdatedAt.Format(time.RFC3339))

	plan.Key = types.StringValue(sshKey.Key)
	if v, ok := sshKey.KeyComment.Get(); ok {
		plan.KeyComment = types.StringValue(v)
	}

	plan.KeyType = types.StringValue(string(sshKey.KeyType))
	plan.Asteroid = types.StringValue(sshKey.Asteroid)
	plan.AsteroidName = types.StringValue(sshKey.Asteroid)
	plan.Format = types.StringValue("json")

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SshkeyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resource_sshkey.SshkeyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	sshKey, err := r.client.AsteroidsSshkeysGet(ctx, client.AsteroidsSshkeysGetParams{
		AsteroidName: state.Asteroid.ValueString(),
		ID:           int(state.Id.ValueInt64()),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read ssh key, got error: %s", err))
		return
	}

	state.Pk = types.Int64Value(int64(sshKey.Pk))
	state.FormattedKey = types.StringValue(sshKey.FormattedKey)
	state.ShortenedKey = types.StringValue(sshKey.ShortenedKey)
	state.CreatedAt = types.StringValue(sshKey.CreatedAt.Format(time.RFC3339))
	state.UpdatedAt = types.StringValue(sshKey.UpdatedAt.Format(time.RFC3339))

	state.Key = types.StringValue(sshKey.Key)
	if v, ok := sshKey.KeyComment.Get(); ok {
		state.KeyComment = types.StringValue(v)
	} else {
		state.KeyComment = types.StringNull()
	}

	state.KeyType = types.StringValue(string(sshKey.KeyType))
	state.Asteroid = types.StringValue(sshKey.Asteroid)
	state.AsteroidName = types.StringValue(sshKey.Asteroid)
	state.Id = types.Int64Value(int64(sshKey.Pk))
	state.Format = types.StringValue("json")

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *SshkeyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan resource_sshkey.SshkeyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.AsteroidsSshkeysDelete(ctx, client.AsteroidsSshkeysDeleteParams{
		AsteroidName: state.Asteroid.ValueString(),
		ID:           int(state.Id.ValueInt64()),
	}); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete ssh key, got error: %s", err))
		return
	}

	reqBody := client.SshKeyRequest{
		Asteroid: plan.Asteroid.ValueString(),
		Key:      plan.Key.ValueString(),
		KeyType:  client.KeyTypeEnum(plan.KeyType.ValueString()),
	}
	if !plan.KeyComment.IsNull() {
		reqBody.KeyComment.SetTo(plan.KeyComment.ValueString())
	}

	apiReq := client.AsteroidsSshkeysCreateApplicationJSON(reqBody)

	sshKey, err := r.client.AsteroidsSshkeysCreate(ctx, &apiReq, client.AsteroidsSshkeysCreateParams{
		AsteroidName: plan.Asteroid.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create ssh key, got error: %s", err))
		return
	}

	plan.Pk = types.Int64Value(int64(sshKey.Pk))
	plan.Id = types.Int64Value(int64(sshKey.Pk))
	plan.FormattedKey = types.StringValue(sshKey.FormattedKey)
	plan.ShortenedKey = types.StringValue(sshKey.ShortenedKey)
	plan.CreatedAt = types.StringValue(sshKey.CreatedAt.Format(time.RFC3339))
	plan.UpdatedAt = types.StringValue(sshKey.UpdatedAt.Format(time.RFC3339))

	plan.Key = types.StringValue(sshKey.Key)
	if v, ok := sshKey.KeyComment.Get(); ok {
		plan.KeyComment = types.StringValue(v)
	}

	plan.KeyType = types.StringValue(string(sshKey.KeyType))
	plan.Asteroid = types.StringValue(sshKey.Asteroid)
	plan.AsteroidName = types.StringValue(sshKey.Asteroid)
	plan.Format = types.StringValue("json")

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *SshkeyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resource_sshkey.SshkeyModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.AsteroidsSshkeysDelete(ctx, client.AsteroidsSshkeysDeleteParams{
		AsteroidName: state.Asteroid.ValueString(),
		ID:           int(state.Id.ValueInt64()),
	}); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete ssh key, got error: %s", err))
		return
	}
}
