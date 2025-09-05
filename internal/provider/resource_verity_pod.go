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
	Name   types.String `tfsdk:"name"`
	Enable types.Bool   `tfsdk:"enable"`
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
	podProps := &openapi.PodsPutRequestPodValue{
		Name: openapi.PtrString(name),
	}

	if !plan.Enable.IsNull() {
		podProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
	}

	operationID := r.bulkOpsMgr.AddPut(ctx, "pod", name, *podProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for pod creation operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Create Pod %s", name))...,
		)
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
		tflog.Info(ctx, fmt.Sprintf("Skipping Pod %s verification – trusting recent successful API operation", podName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching Pods for verification of %s", podName))

	type PodResponse struct {
		Pod map[string]map[string]interface{} `json:"pod"`
	}

	var result PodResponse
	var err error
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		podsData, fetchErr := getCachedResponse(ctx, r.provCtx, "pods", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch pods")
			respAPI, err := r.client.PodsAPI.PodsGet(ctx).Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading pods: %v", err)
			}
			defer respAPI.Body.Close()

			var result PodResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&result); err != nil {
				return nil, fmt.Errorf("failed to decode pods response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d pods", len(result.Pod)))
			return result, nil
		})
		if fetchErr != nil {
			err = fetchErr
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch pods on attempt %d, retrying in %v", attempt+1, sleepTime))
			time.Sleep(sleepTime)
			continue
		}
		result = podsData.(PodResponse)
		break
	}

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Pod %s", podName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for pod with ID: %s", podName))
	var podData map[string]interface{}
	exists := false

	if data, ok := result.Pod[podName]; ok {
		podData = data
		exists = true
		tflog.Debug(ctx, fmt.Sprintf("Found pod directly by ID: %s", podName))
	} else {
		for apiName, pod := range result.Pod {
			if resourceName, ok := pod["name"].(string); ok && resourceName == podName {
				podData = pod
				podName = apiName
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found pod with name '%s' under API key '%s'", resourceName, apiName))
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Pod with ID '%s' not found in API response", podName))
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(fmt.Sprintf("%v", podData["name"]))

	if enable, ok := podData["enable"].(bool); ok {
		state.Enable = types.BoolValue(enable)
	} else {
		state.Enable = types.BoolNull()
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
	podProps := openapi.PodsPutRequestPodValue{}
	hasChanges := false

	if !plan.Name.Equal(state.Name) {
		podProps.Name = openapi.PtrString(name)
		hasChanges = true
	}

	if !plan.Enable.Equal(state.Enable) {
		podProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	operationID := r.bulkOpsMgr.AddPatch(ctx, "pod", name, podProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for pod update operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update Pod %s", name))...,
		)
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
	operationID := r.bulkOpsMgr.AddDelete(ctx, "pod", name)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for pod deletion operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Delete Pod %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Pod %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "pods")
	resp.State.RemoveResource(ctx)
}

func (r *verityPodResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
