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
	bulkOpsMgr           *utils.BulkOperationManager
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

	if !plan.Enable.IsNull() {
		ipv4ListProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
	}

	if !plan.Ipv4List.IsNull() {
		ipv4ListProps.Ipv4List = openapi.PtrString(plan.Ipv4List.ValueString())
	}

	operationID := r.bulkOpsMgr.AddPut(ctx, "ipv4_list_filter", name, *ipv4ListProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for IPv4 List create operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Create IPv4 List %s", name))...,
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("IPv4 List %s create operation completed successfully", name))
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

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("ipv4_list_filter") {
		tflog.Info(ctx, fmt.Sprintf("Skipping IPv4 List %s verification – trusting recent successful API operation", name))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching IPv4 List Filters for verification of %s", name))

	type Ipv4ListResponse struct {
		Ipv4ListFilter map[string]interface{} `json:"ipv4_list_filter"`
	}

	var result Ipv4ListResponse
	var err error
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		ipv4Data, fetchErr := getCachedResponse(ctx, r.provCtx, "ipv4_list_filters", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch IPv4 List Filters")
			respAPI, err := r.client.IPv4ListFiltersAPI.Ipv4listsGet(ctx).Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading IPv4 List Filters: %v", err)
			}
			defer respAPI.Body.Close()

			var res Ipv4ListResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return nil, fmt.Errorf("failed to decode IPv4 List Filters response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d IPv4 List Filters", len(res.Ipv4ListFilter)))
			return res, nil
		})
		if fetchErr != nil {
			err = fetchErr
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch IPv4 List Filters on attempt %d, retrying in %v", attempt+1, sleepTime))
			time.Sleep(sleepTime)
			continue
		}
		result = ipv4Data.(Ipv4ListResponse)
		break
	}
	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read IPv4 List %s", name))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for IPv4 List with ID: %s", name))
	var ipv4ListData map[string]interface{}
	exists := false

	if data, ok := result.Ipv4ListFilter[name].(map[string]interface{}); ok {
		ipv4ListData = data
		exists = true
		tflog.Debug(ctx, fmt.Sprintf("Found IPv4 List directly by ID: %s", name))
	} else {
		for apiName, ipv4 := range result.Ipv4ListFilter {
			ipv4List, ok := ipv4.(map[string]interface{})
			if !ok {
				continue
			}

			if resourceName, ok := ipv4List["name"].(string); ok && resourceName == name {
				ipv4ListData = ipv4List
				name = apiName
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found IPv4 List with name '%s' under API key '%s'", resourceName, apiName))
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("IPv4 List with ID '%s' not found in API response", name))
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(fmt.Sprintf("%v", ipv4ListData["name"]))

	if enable, ok := ipv4ListData["enable"].(bool); ok {
		state.Enable = types.BoolValue(enable)
	} else {
		state.Enable = types.BoolNull()
	}

	if ipv4List, ok := ipv4ListData["ipv4_list"].(string); ok {
		state.Ipv4List = types.StringValue(ipv4List)
	} else {
		state.Ipv4List = types.StringNull()
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

	if !plan.Enable.Equal(state.Enable) {
		ipv4ListProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
		hasChanges = true
	}

	if !plan.Ipv4List.Equal(state.Ipv4List) {
		ipv4ListProps.Ipv4List = openapi.PtrString(plan.Ipv4List.ValueString())
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	operationID := r.bulkOpsMgr.AddPatch(ctx, "ipv4_list_filter", name, ipv4ListProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for IPv4 List update operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update IPv4 List %s", name))...,
		)
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
	operationID := r.bulkOpsMgr.AddDelete(ctx, "ipv4_list_filter", name)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for IPv4 List delete operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Delete IPv4 List %s", name))...,
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("IPv4 List %s delete operation completed successfully", name))
	clearCache(ctx, r.provCtx, "ipv4_lists")
	resp.State.RemoveResource(ctx)
}

func (r *verityIpv4ListResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
