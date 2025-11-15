package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/uberspace-community/terraform-provider-uberspace/gen/client"
	"github.com/uberspace-community/terraform-provider-uberspace/gen/provider/resource_mailuser"
)

// Ensure provider defined types fully satisfy framework interfaces.
var (
	_ resource.Resource = &MailuserResource{}
	// _ resource.ResourceWithImportState = &MailuserResource{}.
)

func NewMailuserResource() resource.Resource {
	return &MailuserResource{}
}

// MailuserResource defines the resource implementation.
type MailuserResource struct {
	client *client.Client
}

func (r *MailuserResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mailuser"
}

func (r *MailuserResource) Schema(ctx context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = resource_mailuser.MailuserResourceSchema(ctx)
}

func (r *MailuserResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// convertForwards converts a slice of client.NestedMailForward into a Terraform types.List
// using the generated resource_mailuser.ForwardsValue helpers and returns any diagnostics.
func convertForwards(ctx context.Context, forwards []client.NestedMailForward) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	attrTypes := resource_mailuser.ForwardsValue{}.AttributeTypes(ctx)
	vals := make([]attr.Value, 0, len(forwards))

	for _, f := range forwards {
		attrs := map[string]attr.Value{
			"destination": types.StringValue(f.Destination),
			"keep":        types.BoolValue(f.Keep),
		}

		fv, d := resource_mailuser.NewForwardsValue(attrTypes, attrs) //nolint:contextcheck
		if d.HasError() {
			diags = append(diags, d...)
			return types.ListNull(resource_mailuser.ForwardsValue{}.Type(ctx)), diags
		}

		obj, d := fv.ToObjectValue(ctx)
		if d.HasError() {
			diags = append(diags, d...)
			return types.ListNull(resource_mailuser.ForwardsValue{}.Type(ctx)), diags
		}

		vals = append(vals, obj)
	}

	listVal, d := types.ListValue(resource_mailuser.ForwardsValue{}.Type(ctx), vals)
	if d.HasError() {
		diags = append(diags, d...)
		return types.ListNull(resource_mailuser.ForwardsValue{}.Type(ctx)), diags
	}

	return listVal, diags
}

func (r *MailuserResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resource_mailuser.MailuserModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := client.AsteroidsMaildomainsUsersCreateApplicationJSON(client.MailUserRequest{
		Name:         plan.Name.ValueString(),
		PasswordHash: client.NewOptNilString(plan.PasswordHash.ValueString()),
	})

	Mailuser, err := r.client.AsteroidsMaildomainsUsersCreate(ctx, &apiReq, client.AsteroidsMaildomainsUsersCreateParams{
		AsteroidName:   plan.AsteroidName.ValueString(),
		MaildomainName: plan.MaildomainName.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create mail user, got error: %s", err))
		return
	}

	plan.Asteroid = types.StringValue(Mailuser.Asteroid)
	plan.AsteroidName = types.StringValue(Mailuser.Asteroid)
	plan.CreatedAt = types.StringValue(Mailuser.CreatedAt.Format(time.RFC3339))
	plan.Format = types.StringValue("json")

	lv, d := convertForwards(ctx, Mailuser.Forwards)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	plan.Forwards = lv

	plan.IsCatchall = types.BoolValue(Mailuser.IsCatchall.Or(false))
	plan.IsSysmail = types.BoolValue(Mailuser.IsSysmail.Or(false))
	plan.KeepForwards = types.BoolValue(Mailuser.KeepForwards.Or(false))
	plan.Local = types.StringValue(Mailuser.Name)
	plan.Mailaddr = types.StringValue(Mailuser.Mailaddr)
	plan.Name = types.StringValue(Mailuser.Name)
	plan.PasswordHash = types.StringValue(Mailuser.PasswordHash.Or(""))
	plan.Pk = types.StringValue(Mailuser.Pk)
	plan.UpdatedAt = types.StringValue(Mailuser.UpdatedAt.Format(time.RFC3339))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MailuserResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resource_mailuser.MailuserModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	Mailuser, err := r.client.AsteroidsMaildomainsUsersGet(ctx, client.AsteroidsMaildomainsUsersGetParams{
		AsteroidName:   state.AsteroidName.ValueString(),
		Local:          state.Local.ValueString(),
		MaildomainName: state.MaildomainName.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to read mail user, got error: %s", err))
		return
	}

	state.Asteroid = types.StringValue(Mailuser.Asteroid)
	state.AsteroidName = types.StringValue(Mailuser.Asteroid)
	state.CreatedAt = types.StringValue(Mailuser.CreatedAt.Format(time.RFC3339))
	state.Format = types.StringValue("json")

	lv, d := convertForwards(ctx, Mailuser.Forwards)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	state.Forwards = lv
	state.IsCatchall = types.BoolValue(Mailuser.IsCatchall.Or(false))
	state.IsSysmail = types.BoolValue(Mailuser.IsSysmail.Or(false))
	state.KeepForwards = types.BoolValue(Mailuser.KeepForwards.Or(false))
	state.Local = types.StringValue(Mailuser.Name)
	state.Mailaddr = types.StringValue(Mailuser.Mailaddr)
	state.Name = types.StringValue(Mailuser.Name)
	state.PasswordHash = types.StringValue(Mailuser.PasswordHash.Or(""))
	state.Pk = types.StringValue(Mailuser.Pk)
	state.UpdatedAt = types.StringValue(Mailuser.UpdatedAt.Format(time.RFC3339))

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *MailuserResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan resource_mailuser.MailuserModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.AsteroidsMaildomainsUsersDelete(ctx, client.AsteroidsMaildomainsUsersDeleteParams{
		AsteroidName: state.AsteroidName.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete mail user, got error: %s", err))
		return
	}

	apiReq := client.AsteroidsMaildomainsUsersCreateApplicationJSON(client.MailUserRequest{
		Name: plan.Name.ValueString(),
	})

	Mailuser, err := r.client.AsteroidsMaildomainsUsersCreate(ctx, &apiReq, client.AsteroidsMaildomainsUsersCreateParams{
		AsteroidName:   plan.AsteroidName.ValueString(),
		MaildomainName: plan.MaildomainName.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to create mail user, got error: %s", err))
		return
	}

	plan.Asteroid = types.StringValue(Mailuser.Asteroid)
	plan.AsteroidName = types.StringValue(Mailuser.Asteroid)
	plan.CreatedAt = types.StringValue(Mailuser.CreatedAt.Format(time.RFC3339))
	plan.Format = types.StringValue("json")

	lv, d := convertForwards(ctx, Mailuser.Forwards)
	if d.HasError() {
		resp.Diagnostics.Append(d...)
		return
	}

	plan.Forwards = lv

	plan.IsCatchall = types.BoolValue(Mailuser.IsCatchall.Or(false))
	plan.IsSysmail = types.BoolValue(Mailuser.IsSysmail.Or(false))
	plan.KeepForwards = types.BoolValue(Mailuser.KeepForwards.Or(false))
	plan.Local = types.StringValue(Mailuser.Name)
	plan.Mailaddr = types.StringValue(Mailuser.Mailaddr)
	plan.Name = types.StringValue(Mailuser.Name)
	plan.PasswordHash = types.StringValue(Mailuser.PasswordHash.Or(""))
	plan.Pk = types.StringValue(Mailuser.Pk)
	plan.UpdatedAt = types.StringValue(Mailuser.UpdatedAt.Format(time.RFC3339))

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *MailuserResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resource_mailuser.MailuserModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.AsteroidsMaildomainsUsersDelete(ctx, client.AsteroidsMaildomainsUsersDeleteParams{
		AsteroidName:   state.AsteroidName.ValueString(),
		MaildomainName: state.MaildomainName.ValueString(),
		Local:          state.Local.ValueString(),
	}); err != nil {
		resp.Diagnostics.AddError("Client Error", fmt.Sprintf("Unable to delete mail user, got error: %s", err))
		return
	}
}
