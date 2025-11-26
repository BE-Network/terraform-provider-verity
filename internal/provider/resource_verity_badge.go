package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-verity/internal/utils"
	"terraform-provider-verity/openapi"
)

var (
	_ resource.Resource                = &verityBadgeResource{}
	_ resource.ResourceWithConfigure   = &verityBadgeResource{}
	_ resource.ResourceWithImportState = &verityBadgeResource{}
)

func NewVerityBadgeResource() resource.Resource {
	return &verityBadgeResource{}
}

type verityBadgeResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *utils.BulkOperationManager
	notifyOperationAdded func()
}

type verityBadgeResourceModel struct {
	Name             types.String                       `tfsdk:"name"`
	Enable           types.Bool                         `tfsdk:"enable"`
	Color            types.String                       `tfsdk:"color"`
	Number           types.Int64                        `tfsdk:"number"`
	ObjectProperties []verityBadgeObjectPropertiesModel `tfsdk:"object_properties"`
}

type verityBadgeObjectPropertiesModel struct {
	Notes types.String `tfsdk:"notes"`
}

func (r *verityBadgeResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_badge"
}

func (r *verityBadgeResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	provCtx, ok := req.ProviderData.(*providerContext)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *providerContext, got: %T", req.ProviderData),
		)
		return
	}

	r.provCtx = provCtx
	r.client = provCtx.client
	r.bulkOpsMgr = provCtx.bulkOpsMgr
	r.notifyOperationAdded = provCtx.NotifyOperationAdded
}

func (r *verityBadgeResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Badge resource",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Object Name. Must be unique.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enable": schema.BoolAttribute{
				Description: "Enable object.",
				Optional:    true,
			},
			"color": schema.StringAttribute{
				Description: "Badge color.",
				Optional:    true,
			},
			"number": schema.Int64Attribute{
				Description: "Badge number.",
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"object_properties": schema.ListNestedBlock{
				Description: "Object properties for the badge",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"notes": schema.StringAttribute{
							Description: "User Notes.",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func (r *verityBadgeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityBadgeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Authenticate",
			fmt.Sprintf("Error authenticating with API: %s", err),
		)
		return
	}

	name := plan.Name.ValueString()
	badgeProps := &openapi.BadgesPutRequestBadgeValue{
		Name: openapi.PtrString(name),
	}

	// Handle string fields
	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "Color", APIField: &badgeProps.Color, TFValue: plan.Color},
	})

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &badgeProps.Enable, TFValue: plan.Enable},
	})

	// Handle nullable int64 fields
	utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
		{FieldName: "Number", APIField: &badgeProps.Number, TFValue: plan.Number},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
		utils.SetObjectPropertiesFields([]utils.ObjectPropertiesField{
			{Name: "Notes", TFValue: op.Notes, APIValue: &objProps.Notes},
		})
		badgeProps.ObjectProperties = &objProps
	}

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "badge", name, *badgeProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Badge %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "badges")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
}

func (r *verityBadgeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityBadgeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Authenticate",
			fmt.Sprintf("Error authenticating with API: %s", err),
		)
		return
	}

	badgeName := state.Name.ValueString()

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("badge") {
		tflog.Info(ctx, fmt.Sprintf("Skipping badge %s verification – trusting recent successful API operation", badgeName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching badges for verification of %s", badgeName))

	type BadgesResponse struct {
		Badge map[string]interface{} `json:"badge"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "badges", badgeName,
		func() (BadgesResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch badges")
			respAPI, err := r.client.BadgesAPI.BadgesGet(ctx).Execute()
			if err != nil {
				return BadgesResponse{}, fmt.Errorf("error reading badges: %v", err)
			}
			defer respAPI.Body.Close()

			var res BadgesResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return BadgesResponse{}, fmt.Errorf("failed to decode badges response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d badges", len(res.Badge)))
			return res, nil
		},
		getCachedResponse,
	)
	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Badge %s", badgeName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for badge with name: %s", badgeName))

	badgeData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.Badge,
		badgeName,
		func(data interface{}) (string, bool) {
			if badge, ok := data.(map[string]interface{}); ok {
				if name, ok := badge["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Badge with name '%s' not found in API response", badgeName))
		resp.State.RemoveResource(ctx)
		return
	}

	badgeMap, ok := badgeData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Badge Data",
			fmt.Sprintf("Badge data is not in expected format for %s", badgeName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found badge '%s' under API key '%s'", badgeName, actualAPIName))

	state.Name = utils.MapStringFromAPI(badgeMap["name"])

	// Handle object properties
	if objProps, ok := badgeMap["object_properties"].(map[string]interface{}); ok {
		state.ObjectProperties = []verityBadgeObjectPropertiesModel{
			{Notes: utils.MapStringFromAPI(objProps["notes"])},
		}
	} else {
		state.ObjectProperties = nil
	}

	// Map string fields
	stringFieldMappings := map[string]*types.String{
		"color": &state.Color,
	}

	for apiKey, stateField := range stringFieldMappings {
		*stateField = utils.MapStringFromAPI(badgeMap[apiKey])
	}

	// Map boolean fields
	boolFieldMappings := map[string]*types.Bool{
		"enable": &state.Enable,
	}

	for apiKey, stateField := range boolFieldMappings {
		*stateField = utils.MapBoolFromAPI(badgeMap[apiKey])
	}

	// Map int64 fields
	int64FieldMappings := map[string]*types.Int64{
		"number": &state.Number,
	}

	for apiKey, stateField := range int64FieldMappings {
		*stateField = utils.MapInt64FromAPI(badgeMap[apiKey])
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityBadgeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityBadgeResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Authenticate",
			fmt.Sprintf("Error authenticating with API: %s", err),
		)
		return
	}

	name := plan.Name.ValueString()
	badgeProps := openapi.BadgesPutRequestBadgeValue{}
	hasChanges := false

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { badgeProps.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Color, state.Color, func(v *string) { badgeProps.Color = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { badgeProps.Enable = v }, &hasChanges)

	// Handle nullable int64 field changes
	utils.CompareAndSetNullableInt64Field(plan.Number, state.Number, func(v *openapi.NullableInt32) { badgeProps.Number = *v }, &hasChanges)

	// Handle object properties
	if len(plan.ObjectProperties) > 0 && len(state.ObjectProperties) > 0 {
		objProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
		op := plan.ObjectProperties[0]
		st := state.ObjectProperties[0]
		objPropsChanged := false

		utils.CompareAndSetObjectPropertiesFields([]utils.ObjectPropertiesFieldWithComparison{
			{Name: "Notes", PlanValue: op.Notes, StateValue: st.Notes, APIValue: &objProps.Notes},
		}, &objPropsChanged)

		if objPropsChanged {
			badgeProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "badge", name, badgeProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Badge %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "badges")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *verityBadgeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityBadgeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Authenticate",
			fmt.Sprintf("Error authenticating with API: %s", err),
		)
		return
	}

	name := state.Name.ValueString()

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "badge", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Badge %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "badges")
	resp.State.RemoveResource(ctx)
}

func (r *verityBadgeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
