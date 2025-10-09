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
	_ resource.Resource                = &verityPodResource{}
	_ resource.ResourceWithConfigure   = &verityPodResource{}
	_ resource.ResourceWithImportState = &verityPodResource{}
)

func NewVerityPodResource() resource.Resource {
	return &verityPodResource{}
}

type verityPodResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *utils.BulkOperationManager
	notifyOperationAdded func()
}

type verityPodResourceModel struct {
	Name               types.String                     `tfsdk:"name"`
	Enable             types.Bool                       `tfsdk:"enable"`
	ExpectedSpineCount types.Int64                      `tfsdk:"expected_spine_count"`
	ObjectProperties   []verityPodObjectPropertiesModel `tfsdk:"object_properties"`
}

type verityPodObjectPropertiesModel struct {
	Notes types.String `tfsdk:"notes"`
}

func (r *verityPodResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_pod"
}

func (r *verityPodResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityPodResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Pod resource",
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
			"expected_spine_count": schema.Int64Attribute{
				Description: "Number of spine switches expected in this pod",
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"object_properties": schema.ListNestedBlock{
				Description: "Object properties for the pod",
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

func (r *verityPodResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityPodResourceModel
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
	podReq := &openapi.PodsPutRequestPodValue{
		Name: openapi.PtrString(name),
	}

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &podReq.Enable, TFValue: plan.Enable},
	})

	// Handle int64 fields
	utils.SetInt64Fields([]utils.Int64FieldMapping{
		{FieldName: "ExpectedSpineCount", APIField: &podReq.ExpectedSpineCount, TFValue: plan.ExpectedSpineCount},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
		if !op.Notes.IsNull() {
			objProps.Notes = openapi.PtrString(op.Notes.ValueString())
		} else {
			objProps.Notes = nil
		}
		podReq.ObjectProperties = &objProps
	}

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "pod", name, *podReq, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Pod %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "pods")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
}

func (r *verityPodResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityPodResourceModel
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

	podName := state.Name.ValueString()

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("pod") {
		tflog.Info(ctx, fmt.Sprintf("Skipping Pod %s verification - trusting recent successful API operation", podName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("No recent Pod operations found, performing normal verification for %s", podName))

	type PodsResponse struct {
		Pod map[string]map[string]interface{} `json:"pod"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "pods", podName,
		func() (PodsResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch Pods")
			respAPI, err := r.client.PodsAPI.PodsGet(ctx).Execute()
			if err != nil {
				return PodsResponse{}, fmt.Errorf("error reading Pod: %v", err)
			}
			defer respAPI.Body.Close()

			var res PodsResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return PodsResponse{}, fmt.Errorf("failed to decode Pods response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d Pods from API", len(res.Pod)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Pod %s", podName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for Pod with name: %s", podName))

	podData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.Pod,
		podName,
		func(data map[string]interface{}) (string, bool) {
			if name, ok := data["name"].(string); ok {
				return name, true
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Pod with name '%s' not found in API response", podName))
		resp.State.RemoveResource(ctx)
		return
	}

	podMap, ok := (interface{}(podData)).(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Pod Data",
			fmt.Sprintf("Pod data is not in expected format for %s", podName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found Pod '%s' under API key '%s'", podName, actualAPIName))

	state.Name = utils.MapStringFromAPI(podMap["name"])

	// Handle object properties
	if objProps, ok := podMap["object_properties"].(map[string]interface{}); ok {
		state.ObjectProperties = []verityPodObjectPropertiesModel{
			{Notes: utils.MapStringFromAPI(objProps["notes"])},
		}
	} else {
		state.ObjectProperties = nil
	}

	// Map boolean fields
	boolFieldMappings := map[string]*types.Bool{
		"enable": &state.Enable,
	}

	for apiKey, stateField := range boolFieldMappings {
		*stateField = utils.MapBoolFromAPI(podMap[apiKey])
	}

	// Map int64 fields
	int64FieldMappings := map[string]*types.Int64{
		"expected_spine_count": &state.ExpectedSpineCount,
	}

	for apiKey, stateField := range int64FieldMappings {
		*stateField = utils.MapInt64FromAPI(podMap[apiKey])
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityPodResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityPodResourceModel

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
	podReq := openapi.PodsPutRequestPodValue{}
	hasChanges := false

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { podReq.Name = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { podReq.Enable = v }, &hasChanges)

	// Handle int64 field changes
	utils.CompareAndSetInt64Field(plan.ExpectedSpineCount, state.ExpectedSpineCount, func(v *int32) { podReq.ExpectedSpineCount = v }, &hasChanges)

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		if len(state.ObjectProperties) == 0 || !plan.ObjectProperties[0].Notes.Equal(state.ObjectProperties[0].Notes) {
			objProps := openapi.AclsPutRequestIpFilterValueObjectProperties{}
			if !plan.ObjectProperties[0].Notes.IsNull() {
				objProps.Notes = openapi.PtrString(plan.ObjectProperties[0].Notes.ValueString())
			} else {
				objProps.Notes = nil
			}
			podReq.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "pod", name, podReq, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Pod %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "pods")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *verityPodResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityPodResourceModel
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

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "pod", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Pod %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "pods")
	resp.State.RemoveResource(ctx)
}

func (r *verityPodResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
