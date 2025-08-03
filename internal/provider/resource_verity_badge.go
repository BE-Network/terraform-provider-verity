package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

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

	if !plan.Color.IsNull() {
		badgeProps.Color = openapi.PtrString(plan.Color.ValueString())
	}

	if !plan.Number.IsNull() {
		badgeProps.Number = openapi.PtrInt32(int32(plan.Number.ValueInt64()))
	}

	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
		if !op.Notes.IsNull() {
			objProps.Notes = openapi.PtrString(op.Notes.ValueString())
		} else {
			objProps.Notes = nil
		}
		badgeProps.ObjectProperties = &objProps
	}

	operationID := r.bulkOpsMgr.AddPut(ctx, "badge", name, *badgeProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for badge creation operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Create Badge %s", name))...,
		)
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

	var result BadgesResponse
	var err error
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		badgesData, fetchErr := getCachedResponse(ctx, r.provCtx, "badges", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch badges")
			respAPI, err := r.client.BadgesAPI.BadgesGet(ctx).Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading badges: %v", err)
			}
			defer respAPI.Body.Close()

			var res BadgesResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return nil, fmt.Errorf("failed to decode badges response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d badges", len(res.Badge)))
			return res, nil
		})
		if fetchErr != nil {
			err = fetchErr
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch badges on attempt %d, retrying in %v", attempt+1, sleepTime))
			time.Sleep(sleepTime)
			continue
		}
		result = badgesData.(BadgesResponse)
		break
	}
	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Badge %s", badgeName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for badge with ID: %s", badgeName))
	var badgeData map[string]interface{}
	exists := false

	if data, ok := result.Badge[badgeName].(map[string]interface{}); ok {
		badgeData = data
		exists = true
		tflog.Debug(ctx, fmt.Sprintf("Found badge directly by ID: %s", badgeName))
	} else {
		for apiName, b := range result.Badge {
			badge, ok := b.(map[string]interface{})
			if !ok {
				continue
			}

			if name, ok := badge["name"].(string); ok && name == badgeName {
				badgeData = badge
				badgeName = apiName
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found badge with name '%s' under API key '%s'", name, apiName))
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Badge with ID '%s' not found in API response", badgeName))
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(fmt.Sprintf("%v", badgeData["name"]))

	if color, ok := badgeData["color"].(string); ok && color != "" {
		state.Color = types.StringValue(color)
	} else {
		state.Color = types.StringNull()
	}

	if number, ok := badgeData["number"].(float64); ok {
		state.Number = types.Int64Value(int64(number))
	} else if number, ok := badgeData["number"].(int); ok {
		state.Number = types.Int64Value(int64(number))
	} else if number, ok := badgeData["number"].(int64); ok {
		state.Number = types.Int64Value(number)
	} else {
		state.Number = types.Int64Null()
	}

	if objProps, ok := badgeData["object_properties"].(map[string]interface{}); ok {
		if notes, ok := objProps["notes"].(string); ok {
			state.ObjectProperties = []verityBadgeObjectPropertiesModel{
				{Notes: types.StringValue(notes)},
			}
		} else {
			state.ObjectProperties = []verityBadgeObjectPropertiesModel{
				{Notes: types.StringNull()},
			}
		}
	} else {
		state.ObjectProperties = nil
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

	if !plan.Name.Equal(state.Name) {
		badgeProps.Name = openapi.PtrString(name)
		hasChanges = true
	}

	if !plan.Color.Equal(state.Color) {
		if !plan.Color.IsNull() {
			badgeProps.Color = openapi.PtrString(plan.Color.ValueString())
		} else {
			badgeProps.Color = nil
		}
		hasChanges = true
	}

	if !plan.Number.Equal(state.Number) {
		if !plan.Number.IsNull() {
			badgeProps.Number = openapi.PtrInt32(int32(plan.Number.ValueInt64()))
		} else {
			badgeProps.Number = nil
		}
		hasChanges = true
	}

	if len(plan.ObjectProperties) > 0 {
		if len(state.ObjectProperties) == 0 || !plan.ObjectProperties[0].Notes.Equal(state.ObjectProperties[0].Notes) {
			objProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
			if !plan.ObjectProperties[0].Notes.IsNull() {
				objProps.Notes = openapi.PtrString(plan.ObjectProperties[0].Notes.ValueString())
			} else {
				objProps.Notes = nil
			}
			badgeProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	operationID := r.bulkOpsMgr.AddPatch(ctx, "badge", name, badgeProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for badge update operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update Badge %s", name))...,
		)
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
	operationID := r.bulkOpsMgr.AddDelete(ctx, "badge", name)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for badge deletion operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Delete Badge %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Badge %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "badges")
	resp.State.RemoveResource(ctx)
}

func (r *verityBadgeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
