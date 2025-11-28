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
)

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
			},
			"ipv6_list": schema.StringAttribute{
				Description: "Comma separated list of IPv6 addresses",
				Optional:    true,
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

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
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

	state.Name = utils.MapStringFromAPI(ipv6ListMap["name"])

	// Map boolean fields
	boolFieldMappings := map[string]*types.Bool{
		"enable": &state.Enable,
	}

	for apiKey, stateField := range boolFieldMappings {
		*stateField = utils.MapBoolFromAPI(ipv6ListMap[apiKey])
	}

	// Map string fields
	stringFieldMappings := map[string]*types.String{
		"ipv6_list": &state.Ipv6List,
	}

	for apiKey, stateField := range stringFieldMappings {
		*stateField = utils.MapStringFromAPI(ipv6ListMap[apiKey])
	}

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
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
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
