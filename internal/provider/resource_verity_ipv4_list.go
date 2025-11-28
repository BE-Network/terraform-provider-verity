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
	_ resource.Resource                = &verityIpv4ListResource{}
	_ resource.ResourceWithConfigure   = &verityIpv4ListResource{}
	_ resource.ResourceWithImportState = &verityIpv4ListResource{}
)

func NewVerityIpv4ListResource() resource.Resource {
	return &verityIpv4ListResource{}
}

type verityIpv4ListResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type verityIpv4ListResourceModel struct {
	Name     types.String `tfsdk:"name"`
	Enable   types.Bool   `tfsdk:"enable"`
	Ipv4List types.String `tfsdk:"ipv4_list"`
}

func (r *verityIpv4ListResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_ipv4_list"
}

func (r *verityIpv4ListResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityIpv4ListResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity IPv4 List Filter",
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
			"ipv4_list": schema.StringAttribute{
				Description: "Comma separated list of IPv4 addresses",
				Optional:    true,
			},
		},
	}
}

func (r *verityIpv4ListResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityIpv4ListResourceModel
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
	ipv4ListProps := &openapi.Ipv4listsPutRequestIpv4ListFilterValue{
		Name: openapi.PtrString(name),
	}

	// Handle string fields
	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "Ipv4List", APIField: &ipv4ListProps.Ipv4List, TFValue: plan.Ipv4List},
	})

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &ipv4ListProps.Enable, TFValue: plan.Enable},
	})

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "ipv4_list", name, *ipv4ListProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("IPv4 List %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "ipv4_lists")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
}

func (r *verityIpv4ListResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityIpv4ListResourceModel
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

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("ipv4_list") {
		tflog.Info(ctx, fmt.Sprintf("Skipping IPv4 List %s verification – trusting recent successful API operation", name))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching IPv4 List Filters for verification of %s", name))

	type Ipv4ListResponse struct {
		Ipv4ListFilter map[string]interface{} `json:"ipv4_list_filter"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "ipv4_lists", name,
		func() (Ipv4ListResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch IPv4 List Filters")
			respAPI, err := r.client.IPv4ListFiltersAPI.Ipv4listsGet(ctx).Execute()
			if err != nil {
				return Ipv4ListResponse{}, fmt.Errorf("error reading IPv4 List Filters: %v", err)
			}
			defer respAPI.Body.Close()

			var res Ipv4ListResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return Ipv4ListResponse{}, fmt.Errorf("failed to decode IPv4 List Filters response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d IPv4 List Filters", len(res.Ipv4ListFilter)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read IPv4 List %s", name))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for IPv4 List with name: %s", name))

	ipv4ListData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.Ipv4ListFilter,
		name,
		func(data interface{}) (string, bool) {
			if ipv4List, ok := data.(map[string]interface{}); ok {
				if resourceName, ok := ipv4List["name"].(string); ok {
					return resourceName, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("IPv4 List with name '%s' not found in API response", name))
		resp.State.RemoveResource(ctx)
		return
	}

	ipv4ListMap, ok := ipv4ListData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid IPv4 List Data",
			fmt.Sprintf("IPv4 List data is not in expected format for %s", name),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found IPv4 List '%s' under API key '%s'", name, actualAPIName))

	state.Name = utils.MapStringFromAPI(ipv4ListMap["name"])

	// Map boolean fields
	boolFieldMappings := map[string]*types.Bool{
		"enable": &state.Enable,
	}

	for apiKey, stateField := range boolFieldMappings {
		*stateField = utils.MapBoolFromAPI(ipv4ListMap[apiKey])
	}

	// Map string fields
	stringFieldMappings := map[string]*types.String{
		"ipv4_list": &state.Ipv4List,
	}

	for apiKey, stateField := range stringFieldMappings {
		*stateField = utils.MapStringFromAPI(ipv4ListMap[apiKey])
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityIpv4ListResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityIpv4ListResourceModel

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
	ipv4ListProps := openapi.Ipv4listsPutRequestIpv4ListFilterValue{}
	hasChanges := false

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { ipv4ListProps.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Ipv4List, state.Ipv4List, func(v *string) { ipv4ListProps.Ipv4List = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { ipv4ListProps.Enable = v }, &hasChanges)

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "ipv4_list", name, ipv4ListProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("IPv4 List %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "ipv4_lists")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *verityIpv4ListResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityIpv4ListResourceModel
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "ipv4_list", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("IPv4 List %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "ipv4_lists")
	resp.State.RemoveResource(ctx)
}

func (r *verityIpv4ListResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
