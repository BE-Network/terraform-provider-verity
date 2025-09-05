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
	bulkOpsMgr           *utils.BulkOperationManager
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

	if !plan.Enable.IsNull() {
		ipv6ListProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
	}

	if !plan.Ipv6List.IsNull() {
		ipv6ListProps.Ipv6List = openapi.PtrString(plan.Ipv6List.ValueString())
	}

	operationID := r.bulkOpsMgr.AddPut(ctx, "ipv6_list_filter", name, *ipv6ListProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for IPv6 List create operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Create IPv6 List %s", name))...,
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("IPv6 List %s create operation completed successfully", name))
	clearCache(ctx, r.provCtx, "ipv6_list_filters")

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

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("ipv6_list_filter") {
		tflog.Info(ctx, fmt.Sprintf("Skipping IPv6 List %s verification – trusting recent successful API operation", ipv6ListName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching IPv6 List Filters for verification of %s", ipv6ListName))

	type Ipv6ListsResponse struct {
		Ipv6ListFilter map[string]interface{} `json:"ipv6_list_filter"`
	}

	var result Ipv6ListsResponse
	var err error
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		ipv6ListsData, fetchErr := getCachedResponse(ctx, r.provCtx, "ipv6_list_filters", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch IPv6 List Filters")
			respAPI, err := r.client.IPv6ListFiltersAPI.Ipv6listsGet(ctx).Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading IPv6 List Filters: %v", err)
			}
			defer respAPI.Body.Close()

			var res Ipv6ListsResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return nil, fmt.Errorf("failed to decode IPv6 List Filters response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d IPv6 List Filters", len(res.Ipv6ListFilter)))
			return res, nil
		})
		if fetchErr != nil {
			err = fetchErr
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch IPv6 List Filters on attempt %d, retrying in %v", attempt+1, sleepTime))
			time.Sleep(sleepTime)
			continue
		}
		result = ipv6ListsData.(Ipv6ListsResponse)
		break
	}
	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read IPv6 List %s", ipv6ListName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for IPv6 List Filter with ID: %s", ipv6ListName))
	var ipv6ListData map[string]interface{}
	exists := false

	if data, ok := result.Ipv6ListFilter[ipv6ListName].(map[string]interface{}); ok {
		ipv6ListData = data
		exists = true
		tflog.Debug(ctx, fmt.Sprintf("Found IPv6 List Filter directly by ID: %s", ipv6ListName))
	} else {
		for apiName, i := range result.Ipv6ListFilter {
			ipv6List, ok := i.(map[string]interface{})
			if !ok {
				continue
			}

			if name, ok := ipv6List["name"].(string); ok && name == ipv6ListName {
				ipv6ListData = ipv6List
				ipv6ListName = apiName
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found IPv6 List Filter with name '%s' under API key '%s'", name, apiName))
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("IPv6 List Filter with ID '%s' not found in API response", ipv6ListName))
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(fmt.Sprintf("%v", ipv6ListData["name"]))

	if enable, ok := ipv6ListData["enable"].(bool); ok {
		state.Enable = types.BoolValue(enable)
	} else {
		state.Enable = types.BoolNull()
	}

	if ipv6List, ok := ipv6ListData["ipv6_list"].(string); ok {
		state.Ipv6List = types.StringValue(ipv6List)
	} else {
		state.Ipv6List = types.StringNull()
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

	if !plan.Enable.Equal(state.Enable) {
		ipv6ListProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
		hasChanges = true
	}

	if !plan.Ipv6List.Equal(state.Ipv6List) {
		ipv6ListProps.Ipv6List = openapi.PtrString(plan.Ipv6List.ValueString())
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	operationID := r.bulkOpsMgr.AddPatch(ctx, "ipv6_list_filter", name, ipv6ListProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for IPv6 List update operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update IPv6 List %s", name))...,
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("IPv6 List %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "ipv6_list_filters")
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
	operationID := r.bulkOpsMgr.AddDelete(ctx, "ipv6_list_filter", name)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for IPv6 List delete operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Delete IPv6 List %s", name))...,
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("IPv6 List %s delete operation completed successfully", name))
	clearCache(ctx, r.provCtx, "ipv6_list_filters")
	resp.State.RemoveResource(ctx)
}

func (r *verityIpv6ListResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
