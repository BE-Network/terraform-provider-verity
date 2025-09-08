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
	_ resource.Resource                = &veritySflowCollectorResource{}
	_ resource.ResourceWithConfigure   = &veritySflowCollectorResource{}
	_ resource.ResourceWithImportState = &veritySflowCollectorResource{}
)

func NewVeritySflowCollectorResource() resource.Resource {
	return &veritySflowCollectorResource{}
}

type veritySflowCollectorResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *utils.BulkOperationManager
	notifyOperationAdded func()
}

type veritySflowCollectorResourceModel struct {
	Name   types.String `tfsdk:"name"`
	Enable types.Bool   `tfsdk:"enable"`
	Ip     types.String `tfsdk:"ip"`
	Port   types.Int64  `tfsdk:"port"`
}

func (r *veritySflowCollectorResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sflow_collector"
}

func (r *veritySflowCollectorResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *veritySflowCollectorResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity SFlow Collector",
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
			"ip": schema.StringAttribute{
				Description: "IP address of the sFlow Collector",
				Optional:    true,
			},
			"port": schema.Int64Attribute{
				Description: "Port (maximum 65535)",
				Optional:    true,
			},
		},
	}
}

func (r *veritySflowCollectorResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan veritySflowCollectorResourceModel
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
	sflowCollectorProps := &openapi.SflowcollectorsPutRequestSflowCollectorValue{
		Name: openapi.PtrString(name),
	}

	if !plan.Enable.IsNull() {
		sflowCollectorProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
	}

	if !plan.Ip.IsNull() {
		sflowCollectorProps.Ip = openapi.PtrString(plan.Ip.ValueString())
	}

	if !plan.Port.IsNull() {
		sflowCollectorProps.Port = openapi.PtrInt32(int32(plan.Port.ValueInt64()))
	}

	operationID := r.bulkOpsMgr.AddPut(ctx, "sflow_collector", name, *sflowCollectorProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for sflow collector creation operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Create SFlow Collector %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("SFlow Collector %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "sflowcollectors")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
}

func (r *veritySflowCollectorResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state veritySflowCollectorResourceModel
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

	sflowCollectorName := state.Name.ValueString()

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("sflow_collector") {
		tflog.Info(ctx, fmt.Sprintf("Skipping sflow collector %s verification – trusting recent successful API operation", sflowCollectorName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching sflow collectors for verification of %s", sflowCollectorName))

	type SflowCollectorsResponse struct {
		SflowCollector map[string]interface{} `json:"sflow_collector"`
	}

	var result SflowCollectorsResponse
	var err error
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		sflowCollectorsData, fetchErr := getCachedResponse(ctx, r.provCtx, "sflowcollectors", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch sflow collectors")
			respAPI, err := r.client.SFlowCollectorsAPI.SflowcollectorsGet(ctx).Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading sflow collectors: %v", err)
			}
			defer respAPI.Body.Close()

			var res SflowCollectorsResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return nil, fmt.Errorf("failed to decode sflow collectors response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d sflow collectors", len(res.SflowCollector)))
			return res, nil
		})
		if fetchErr != nil {
			err = fetchErr
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch sflow collectors on attempt %d, retrying in %v", attempt+1, sleepTime))
			time.Sleep(sleepTime)
			continue
		}
		result = sflowCollectorsData.(SflowCollectorsResponse)
		break
	}
	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read SFlow Collector %s", sflowCollectorName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for sflow collector with ID: %s", sflowCollectorName))
	var sflowCollectorData map[string]interface{}
	exists := false

	if data, ok := result.SflowCollector[sflowCollectorName].(map[string]interface{}); ok {
		sflowCollectorData = data
		exists = true
		tflog.Debug(ctx, fmt.Sprintf("Found sflow collector directly by ID: %s", sflowCollectorName))
	} else {
		for apiName, s := range result.SflowCollector {
			sflowCollector, ok := s.(map[string]interface{})
			if !ok {
				continue
			}

			if name, ok := sflowCollector["name"].(string); ok && name == sflowCollectorName {
				sflowCollectorData = sflowCollector
				sflowCollectorName = apiName
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found sflow collector with name '%s' under API key '%s'", name, apiName))
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("SFlow Collector with ID '%s' not found in API response", sflowCollectorName))
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(fmt.Sprintf("%v", sflowCollectorData["name"]))

	if ip, ok := sflowCollectorData["ip"].(string); ok {
		state.Ip = types.StringValue(ip)
	} else {
		state.Ip = types.StringNull()
	}

	if enable, ok := sflowCollectorData["enable"].(bool); ok {
		state.Enable = types.BoolValue(enable)
	} else {
		state.Enable = types.BoolNull()
	}

	if value, ok := sflowCollectorData["port"]; ok && value != nil {
		switch v := value.(type) {
		case int:
			state.Port = types.Int64Value(int64(v))
		case int32:
			state.Port = types.Int64Value(int64(v))
		case int64:
			state.Port = types.Int64Value(v)
		case float64:
			state.Port = types.Int64Value(int64(v))
		default:
			state.Port = types.Int64Null()
		}
	} else {
		state.Port = types.Int64Null()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *veritySflowCollectorResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state veritySflowCollectorResourceModel

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
	sflowCollectorProps := openapi.SflowcollectorsPutRequestSflowCollectorValue{}
	hasChanges := false

	if !plan.Name.Equal(state.Name) {
		sflowCollectorProps.Name = openapi.PtrString(name)
		hasChanges = true
	}

	if !plan.Enable.Equal(state.Enable) {
		sflowCollectorProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
		hasChanges = true
	}

	if !plan.Ip.Equal(state.Ip) {
		sflowCollectorProps.Ip = openapi.PtrString(plan.Ip.ValueString())
		hasChanges = true
	}

	if !plan.Port.Equal(state.Port) {
		sflowCollectorProps.Port = openapi.PtrInt32(int32(plan.Port.ValueInt64()))
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	operationID := r.bulkOpsMgr.AddPatch(ctx, "sflow_collector", name, sflowCollectorProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for sflow collector update operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update SFlow Collector %s", name))...,
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("SFlow Collector %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "sflowcollectors")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *veritySflowCollectorResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state veritySflowCollectorResourceModel
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
	operationID := r.bulkOpsMgr.AddDelete(ctx, "sflow_collector", name)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for sflow collector deletion operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Delete SFlow Collector %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("SFlow Collector %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "sflowcollectors")
	resp.State.RemoveResource(ctx)
}

func (r *veritySflowCollectorResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
