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

	"terraform-provider-verity/internal/bulkops"
	"terraform-provider-verity/internal/utils"
	"terraform-provider-verity/openapi"
)

var (
	_ resource.Resource                = &verityIpv6ListResource{}
	_ resource.ResourceWithConfigure   = &verityIpv6ListResource{}
	_ resource.ResourceWithImportState = &verityIpv6ListResource{}
	_ resource.ResourceWithModifyPlan  = &verityIpv6ListResource{}
)

const ipv6ListResourceType = "ipv6lists"

func NewVerityIpv6ListResource() resource.Resource {
	return &verityIpv6ListResource{}
}

type verityIpv6ListResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type verityIpv6ListResourceModel struct {
	Name     types.String `tfsdk:"name"`
	Enable   types.Bool   `tfsdk:"enable"`
	Ipv6List types.String `tfsdk:"ipv6_list"`
}

func (r *verityIpv6ListResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ipv6_list"
}

func (r *verityIpv6ListResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityIpv6ListResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity IPv6 List Filter",
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
				Computed:    true,
			},
			"ipv6_list": schema.StringAttribute{
				Description: "Comma separated list of IPv6 addresses",
				Optional:    true,
				Computed:    true,
			},
		},
	}
}

func (r *verityIpv6ListResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityIpv6ListResourceModel
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
	ipv6ListProps := &openapi.Ipv6listsPutRequestIpv6ListFilterValue{
		Name: openapi.PtrString(name),
	}

	// Handle string fields
	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "Ipv6List", APIField: &ipv6ListProps.Ipv6List, TFValue: plan.Ipv6List},
	})

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &ipv6ListProps.Enable, TFValue: plan.Enable},
	})

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "ipv6_list", name, *ipv6ListProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("IPv6 List %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "ipv6_lists")

	var minState verityIpv6ListResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if ipv6ListData, exists := bulkMgr.GetResourceResponse("ipv6_list", name); exists {
			state := populateIpv6ListState(ctx, minState, ipv6ListData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	// If no cached data, fall back to normal Read
	readReq := resource.ReadRequest{
		State: resp.State,
	}
	readResp := resource.ReadResponse{
		State:       resp.State,
		Diagnostics: resp.Diagnostics,
	}

	r.Read(ctx, readReq, &readResp)
}

func (r *verityIpv6ListResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityIpv6ListResourceModel
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

	ipv6ListName := state.Name.ValueString()

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if ipv6ListData, exists := r.bulkOpsMgr.GetResourceResponse("ipv6_list", ipv6ListName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached IPv6 List data for %s from recent operation", ipv6ListName))
			state = populateIpv6ListState(ctx, state, ipv6ListData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("ipv6_list") {
		tflog.Info(ctx, fmt.Sprintf("Skipping IPv6 List %s verification – trusting recent successful API operation", ipv6ListName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching IPv6 List Filters for verification of %s", ipv6ListName))

	type Ipv6ListsResponse struct {
		Ipv6ListFilter map[string]interface{} `json:"ipv6_list_filter"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "ipv6_lists", ipv6ListName,
		func() (Ipv6ListsResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch IPv6 List Filters")
			respAPI, err := r.client.IPv6ListFiltersAPI.Ipv6listsGet(ctx).Execute()
			if err != nil {
				return Ipv6ListsResponse{}, fmt.Errorf("error reading IPv6 List Filters: %v", err)
			}
			defer respAPI.Body.Close()

			var res Ipv6ListsResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return Ipv6ListsResponse{}, fmt.Errorf("failed to decode IPv6 List Filters response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d IPv6 List Filters", len(res.Ipv6ListFilter)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read IPv6 List %s", ipv6ListName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for IPv6 List with name: %s", ipv6ListName))

	ipv6ListData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.Ipv6ListFilter,
		ipv6ListName,
		func(data interface{}) (string, bool) {
			if ipv6List, ok := data.(map[string]interface{}); ok {
				if resourceName, ok := ipv6List["name"].(string); ok {
					return resourceName, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("IPv6 List with name '%s' not found in API response", ipv6ListName))
		resp.State.RemoveResource(ctx)
		return
	}

	ipv6ListMap, ok := ipv6ListData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid IPv6 List Data",
			fmt.Sprintf("IPv6 List data is not in expected format for %s", ipv6ListName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found IPv6 List '%s' under API key '%s'", ipv6ListName, actualAPIName))

	state = populateIpv6ListState(ctx, state, ipv6ListMap, r.provCtx.mode)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityIpv6ListResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityIpv6ListResourceModel

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
	ipv6ListProps := openapi.Ipv6listsPutRequestIpv6ListFilterValue{}
	hasChanges := false

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { ipv6ListProps.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Ipv6List, state.Ipv6List, func(v *string) { ipv6ListProps.Ipv6List = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { ipv6ListProps.Enable = v }, &hasChanges)

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "ipv6_list", name, ipv6ListProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("IPv6 List %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "ipv6_lists")

	var minState verityIpv6ListResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Try to use cached response from bulk operation to populate state with API values
	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if ipv6ListData, exists := bulkMgr.GetResourceResponse("ipv6_list", name); exists {
			newState := populateIpv6ListState(ctx, minState, ipv6ListData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
			return
		}
	}

	// If no cached data, fall back to normal Read
	readReq := resource.ReadRequest{
		State: resp.State,
	}
	readResp := resource.ReadResponse{
		State:       resp.State,
		Diagnostics: resp.Diagnostics,
	}

	r.Read(ctx, readReq, &readResp)
}

func (r *verityIpv6ListResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityIpv6ListResourceModel
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "ipv6_list", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("IPv6 List %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "ipv6_lists")
	resp.State.RemoveResource(ctx)
}

func (r *verityIpv6ListResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func populateIpv6ListState(ctx context.Context, state verityIpv6ListResourceModel, data map[string]interface{}, mode string) verityIpv6ListResourceModel {
	const resourceType = ipv6ListResourceType

	state.Name = utils.MapStringFromAPI(data["name"])

	// Boolean fields
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)

	// String fields
	state.Ipv6List = utils.MapStringWithMode(data, "ipv6_list", resourceType, mode)

	return state
}

func (r *verityIpv6ListResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// =========================================================================
	// Skip if deleting
	// =========================================================================
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityIpv6ListResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Mode-aware field nullification
	// Set fields that don't apply to current mode to null to prevent
	// "known after apply" messages for irrelevant fields.
	// =========================================================================
	const resourceType = ipv6ListResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyStrings(
		"ipv6_list",
	)

	nullifier.NullifyBools(
		"enable",
	)
}
