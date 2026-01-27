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
	_ resource.Resource                = &verityPacketQueueResource{}
	_ resource.ResourceWithConfigure   = &verityPacketQueueResource{}
	_ resource.ResourceWithImportState = &verityPacketQueueResource{}
	_ resource.ResourceWithModifyPlan  = &verityPacketQueueResource{}
)

const packetQueueResourceType = "packetqueues"

func NewVerityPacketQueueResource() resource.Resource {
	return &verityPacketQueueResource{}
}

type verityPacketQueueResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type verityPacketQueueResourceModel struct {
	Name             types.String                             `tfsdk:"name"`
	Enable           types.Bool                               `tfsdk:"enable"`
	Pbit             []verityPacketQueuePbitModel             `tfsdk:"pbit"`
	Queue            []verityPacketQueueQueueModel            `tfsdk:"queue"`
	ObjectProperties []verityPacketQueueObjectPropertiesModel `tfsdk:"object_properties"`
}

type verityPacketQueuePbitModel struct {
	PacketQueueForPBit types.Int64 `tfsdk:"packet_queue_for_p_bit"`
	Index              types.Int64 `tfsdk:"index"`
}

func (p verityPacketQueuePbitModel) GetIndex() types.Int64 {
	return p.Index
}

type verityPacketQueueQueueModel struct {
	BandwidthForQueue types.Int64  `tfsdk:"bandwidth_for_queue"`
	SchedulerType     types.String `tfsdk:"scheduler_type"`
	SchedulerWeight   types.Int64  `tfsdk:"scheduler_weight"`
	Index             types.Int64  `tfsdk:"index"`
}

func (q verityPacketQueueQueueModel) GetIndex() types.Int64 {
	return q.Index
}

type verityPacketQueueObjectPropertiesModel struct {
	Group types.String `tfsdk:"group"`
}

func (r *verityPacketQueueResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_packet_queue"
}

func (r *verityPacketQueueResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityPacketQueueResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity Packet Queue",
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
		},
		Blocks: map[string]schema.Block{
			"pbit": schema.ListNestedBlock{
				Description: "P-bit configurations",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"packet_queue_for_p_bit": schema.Int64Attribute{
							Description: "Flag indicating this p-bit's Queue",
							Optional:    true,
							Computed:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object. Zero if you want to add an object to the list.",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"queue": schema.ListNestedBlock{
				Description: "Queue configurations",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"bandwidth_for_queue": schema.Int64Attribute{
							Description: "Percentage bandwidth allocated to Queue. 0 is no limit",
							Optional:    true,
							Computed:    true,
						},
						"scheduler_type": schema.StringAttribute{
							Description: "Scheduler Type for Queue",
							Optional:    true,
							Computed:    true,
						},
						"scheduler_weight": schema.Int64Attribute{
							Description: "Weight associated with WRR or DWRR scheduler",
							Optional:    true,
							Computed:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object. Zero if you want to add an object to the list.",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"object_properties": schema.ListNestedBlock{
				Description: "Object properties for the packet queue",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"group": schema.StringAttribute{
							Description: "Group",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (r *verityPacketQueueResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityPacketQueueResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config verityPacketQueueResourceModel
	diags = req.Config.Get(ctx, &config)
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
	pqProps := &openapi.PacketqueuesPutRequestPacketQueueValue{
		Name: openapi.PtrString(name),
	}

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &pqProps.Enable, TFValue: plan.Enable},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties{}
		utils.SetObjectPropertiesFields([]utils.ObjectPropertiesField{
			{Name: "Group", TFValue: op.Group, APIValue: &objProps.Group},
		})
		pqProps.ObjectProperties = &objProps
	}

	// Parse HCL to detect explicitly configured attributes
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, "verity_packet_queue", name)

	// Handle Pbit
	if len(plan.Pbit) > 0 {
		pbitConfigMap := utils.BuildIndexedConfigMap(config.Pbit)

		pbitArray := make([]openapi.PacketqueuesPutRequestPacketQueueValuePbitInner, len(plan.Pbit))
		for i, pbit := range plan.Pbit {
			pbitItem := openapi.PacketqueuesPutRequestPacketQueueValuePbitInner{}

			// Handle int64 fields
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &pbitItem.Index, TFValue: pbit.Index},
			})

			// Get per-block configured info for nullable Int64 fields
			itemIndex := pbit.Index.ValueInt64()
			configItem := pbit // fallback to plan item
			if cfgItem, ok := pbitConfigMap[itemIndex]; ok {
				configItem = cfgItem
			}
			cfg := &utils.IndexedBlockNullableFieldConfig{
				BlockType:       "pbit",
				BlockIndex:      itemIndex,
				ConfiguredAttrs: configuredAttrs,
			}
			utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
				{FieldName: "PacketQueueForPBit", APIField: &pbitItem.PacketQueueForPBit, TFValue: configItem.PacketQueueForPBit, IsConfigured: cfg.IsFieldConfigured("packet_queue_for_p_bit")},
			})

			pbitArray[i] = pbitItem
		}
		pqProps.Pbit = pbitArray
	}

	// Handle Queue
	if len(plan.Queue) > 0 {
		queueConfigMap := utils.BuildIndexedConfigMap(config.Queue)

		queueArray := make([]openapi.PacketqueuesPutRequestPacketQueueValueQueueInner, len(plan.Queue))
		for i, queue := range plan.Queue {
			queueItem := openapi.PacketqueuesPutRequestPacketQueueValueQueueInner{}

			// Handle int64 fields
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &queueItem.Index, TFValue: queue.Index},
			})

			// Get per-block configured info for nullable Int64 fields
			itemIndex := queue.Index.ValueInt64()
			configItem := queue // fallback to plan item
			if cfgItem, ok := queueConfigMap[itemIndex]; ok {
				configItem = cfgItem
			}
			cfg := &utils.IndexedBlockNullableFieldConfig{
				BlockType:       "queue",
				BlockIndex:      itemIndex,
				ConfiguredAttrs: configuredAttrs,
			}
			utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
				{FieldName: "BandwidthForQueue", APIField: &queueItem.BandwidthForQueue, TFValue: configItem.BandwidthForQueue, IsConfigured: cfg.IsFieldConfigured("bandwidth_for_queue")},
				{FieldName: "SchedulerWeight", APIField: &queueItem.SchedulerWeight, TFValue: configItem.SchedulerWeight, IsConfigured: cfg.IsFieldConfigured("scheduler_weight")},
			})

			// Handle string fields
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "SchedulerType", APIField: &queueItem.SchedulerType, TFValue: queue.SchedulerType},
			})

			queueArray[i] = queueItem
		}
		pqProps.Queue = queueArray
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "packet_queue", name, *pqProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Packet Queue %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "packet_queues")

	var minState verityPacketQueueResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if pqData, exists := bulkMgr.GetResourceResponse("packet_queue", name); exists {
			state := populatePacketQueueState(ctx, minState, pqData, r.provCtx.mode)
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

func (r *verityPacketQueueResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityPacketQueueResourceModel
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

	pqName := state.Name.ValueString()

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if pqData, exists := r.bulkOpsMgr.GetResourceResponse("packet_queue", pqName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached packet queue data for %s from recent operation", pqName))
			state = populatePacketQueueState(ctx, state, pqData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("packet_queue") {
		tflog.Info(ctx, fmt.Sprintf("Skipping Packet Queue %s verification – trusting recent successful API operation", pqName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching Packet Queues for verification of %s", pqName))

	type PacketQueueResponse struct {
		PacketQueue map[string]interface{} `json:"packet_queue"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "packet_queues", pqName,
		func() (PacketQueueResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch Packet Queues")
			respAPI, err := r.client.PacketQueuesAPI.PacketqueuesGet(ctx).Execute()
			if err != nil {
				return PacketQueueResponse{}, fmt.Errorf("error reading Packet Queues: %v", err)
			}
			defer respAPI.Body.Close()

			var res PacketQueueResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return PacketQueueResponse{}, fmt.Errorf("failed to decode Packet Queues response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d Packet Queues", len(res.PacketQueue)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Packet Queue %s", pqName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for Packet Queue with name: %s", pqName))

	pqData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.PacketQueue,
		pqName,
		func(data interface{}) (string, bool) {
			if queue, ok := data.(map[string]interface{}); ok {
				if name, ok := queue["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Packet Queue with name '%s' not found in API response", pqName))
		resp.State.RemoveResource(ctx)
		return
	}

	pqMap, ok := pqData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Packet Queue Data",
			fmt.Sprintf("Packet Queue data is not in expected format for %s", pqName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found Packet Queue '%s' under API key '%s'", pqName, actualAPIName))

	state = populatePacketQueueState(ctx, state, pqMap, r.provCtx.mode)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityPacketQueueResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityPacketQueueResourceModel

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
	pqProps := openapi.PacketqueuesPutRequestPacketQueueValue{}
	hasChanges := false

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { pqProps.Name = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { pqProps.Enable = v }, &hasChanges)

	// Handle object properties
	if len(plan.ObjectProperties) > 0 && len(state.ObjectProperties) > 0 {
		objProps := openapi.DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties{}
		op := plan.ObjectProperties[0]
		st := state.ObjectProperties[0]
		objPropsChanged := false

		utils.CompareAndSetObjectPropertiesFields([]utils.ObjectPropertiesFieldWithComparison{
			{Name: "Group", PlanValue: op.Group, StateValue: st.Group, APIValue: &objProps.Group},
		}, &objPropsChanged)

		if objPropsChanged {
			pqProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, "verity_packet_queue", name)
	var config verityPacketQueueResourceModel
	req.Config.Get(ctx, &config)
	pbitConfigMap := utils.BuildIndexedConfigMap(config.Pbit)
	queueConfigMap := utils.BuildIndexedConfigMap(config.Queue)

	// Handle Pbit
	changedPbits, pbitsChanged := utils.ProcessIndexedArrayUpdates(plan.Pbit, state.Pbit,
		utils.IndexedItemHandler[verityPacketQueuePbitModel, openapi.PacketqueuesPutRequestPacketQueueValuePbitInner]{
			CreateNew: func(planItem verityPacketQueuePbitModel) openapi.PacketqueuesPutRequestPacketQueueValuePbitInner {
				newPbit := openapi.PacketqueuesPutRequestPacketQueueValuePbitInner{}

				// Handle int64 fields
				utils.SetInt64Fields([]utils.Int64FieldMapping{
					{FieldName: "Index", APIField: &newPbit.Index, TFValue: planItem.Index},
				})

				// Get per-block configured info for nullable Int64 fields
				itemIndex := planItem.Index.ValueInt64()
				configItem := planItem // fallback to plan item
				if cfgItem, ok := pbitConfigMap[itemIndex]; ok {
					configItem = cfgItem
				}
				cfg := &utils.IndexedBlockNullableFieldConfig{
					BlockType:       "pbit",
					BlockIndex:      itemIndex,
					ConfiguredAttrs: configuredAttrs,
				}
				utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
					{FieldName: "PacketQueueForPBit", APIField: &newPbit.PacketQueueForPBit, TFValue: configItem.PacketQueueForPBit, IsConfigured: cfg.IsFieldConfigured("packet_queue_for_p_bit")},
				})

				return newPbit
			},
			UpdateExisting: func(planItem verityPacketQueuePbitModel, stateItem verityPacketQueuePbitModel) (openapi.PacketqueuesPutRequestPacketQueueValuePbitInner, bool) {
				updatePbit := openapi.PacketqueuesPutRequestPacketQueueValuePbitInner{}
				fieldChanged := false

				// Handle int64 field changes
				utils.CompareAndSetInt64Field(planItem.Index, stateItem.Index, func(v *int32) { updatePbit.Index = v }, &fieldChanged)

				// Handle nullable int64 field changes
				utils.CompareAndSetNullableInt64Field(planItem.PacketQueueForPBit, stateItem.PacketQueueForPBit, func(v *openapi.NullableInt32) { updatePbit.PacketQueueForPBit = *v }, &fieldChanged)

				return updatePbit, fieldChanged
			},
			CreateDeleted: func(index int64) openapi.PacketqueuesPutRequestPacketQueueValuePbitInner {
				return openapi.PacketqueuesPutRequestPacketQueueValuePbitInner{
					Index: openapi.PtrInt32(int32(index)),
				}
			},
		})
	if pbitsChanged {
		pqProps.Pbit = changedPbits
		hasChanges = true
	}

	// Handle Queue
	updatedQueues, queuesChanged := utils.ProcessIndexedArrayUpdates(plan.Queue, state.Queue,
		utils.IndexedItemHandler[verityPacketQueueQueueModel, openapi.PacketqueuesPutRequestPacketQueueValueQueueInner]{
			CreateNew: func(planItem verityPacketQueueQueueModel) openapi.PacketqueuesPutRequestPacketQueueValueQueueInner {
				newQueue := openapi.PacketqueuesPutRequestPacketQueueValueQueueInner{}

				// Handle int64 fields
				utils.SetInt64Fields([]utils.Int64FieldMapping{
					{FieldName: "Index", APIField: &newQueue.Index, TFValue: planItem.Index},
				})

				// Get per-block configured info for nullable Int64 fields
				itemIndex := planItem.Index.ValueInt64()
				configItem := planItem // fallback to plan item
				if cfgItem, ok := queueConfigMap[itemIndex]; ok {
					configItem = cfgItem
				}
				cfg := &utils.IndexedBlockNullableFieldConfig{
					BlockType:       "queue",
					BlockIndex:      itemIndex,
					ConfiguredAttrs: configuredAttrs,
				}
				utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
					{FieldName: "BandwidthForQueue", APIField: &newQueue.BandwidthForQueue, TFValue: configItem.BandwidthForQueue, IsConfigured: cfg.IsFieldConfigured("bandwidth_for_queue")},
					{FieldName: "SchedulerWeight", APIField: &newQueue.SchedulerWeight, TFValue: configItem.SchedulerWeight, IsConfigured: cfg.IsFieldConfigured("scheduler_weight")},
				})

				// Handle string fields
				utils.SetStringFields([]utils.StringFieldMapping{
					{FieldName: "SchedulerType", APIField: &newQueue.SchedulerType, TFValue: planItem.SchedulerType},
				})

				return newQueue
			},
			UpdateExisting: func(planItem, stateItem verityPacketQueueQueueModel) (openapi.PacketqueuesPutRequestPacketQueueValueQueueInner, bool) {
				updateQueue := openapi.PacketqueuesPutRequestPacketQueueValueQueueInner{}
				fieldChanged := false

				// Handle int64 field changes
				utils.CompareAndSetInt64Field(planItem.Index, stateItem.Index, func(v *int32) { updateQueue.Index = v }, &fieldChanged)

				// Handle nullable int64 field changes
				utils.CompareAndSetNullableInt64Field(planItem.BandwidthForQueue, stateItem.BandwidthForQueue, func(v *openapi.NullableInt32) { updateQueue.BandwidthForQueue = *v }, &fieldChanged)
				utils.CompareAndSetNullableInt64Field(planItem.SchedulerWeight, stateItem.SchedulerWeight, func(v *openapi.NullableInt32) { updateQueue.SchedulerWeight = *v }, &fieldChanged)

				// Handle string field changes
				utils.CompareAndSetStringField(planItem.SchedulerType, stateItem.SchedulerType, func(v *string) { updateQueue.SchedulerType = v }, &fieldChanged)

				return updateQueue, fieldChanged
			},
			CreateDeleted: func(idx int64) openapi.PacketqueuesPutRequestPacketQueueValueQueueInner {
				return openapi.PacketqueuesPutRequestPacketQueueValueQueueInner{
					Index: openapi.PtrInt32(int32(idx)),
				}
			},
		})
	if queuesChanged {
		pqProps.Queue = updatedQueues
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "packet_queue", name, pqProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Packet Queue %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "packet_queues")

	var minState verityPacketQueueResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Try to use cached response from bulk operation to populate state with API values
	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if pqData, exists := bulkMgr.GetResourceResponse("packet_queue", name); exists {
			newState := populatePacketQueueState(ctx, minState, pqData, r.provCtx.mode)
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

func (r *verityPacketQueueResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityPacketQueueResourceModel
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "packet_queue", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Packet Queue %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "packet_queues")
	resp.State.RemoveResource(ctx)
}

func (r *verityPacketQueueResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func populatePacketQueueState(ctx context.Context, state verityPacketQueueResourceModel, data map[string]interface{}, mode string) verityPacketQueueResourceModel {
	const resourceType = packetQueueResourceType

	state.Name = utils.MapStringFromAPI(data["name"])

	// Boolean fields
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)

	// Handle object_properties block
	if utils.FieldAppliesToMode(resourceType, "object_properties", mode) {
		if objProps, ok := data["object_properties"].(map[string]interface{}); ok {
			objPropsModel := verityPacketQueueObjectPropertiesModel{
				Group: utils.MapStringWithModeNested(objProps, "group", resourceType, "object_properties.group", mode),
			}
			state.ObjectProperties = []verityPacketQueueObjectPropertiesModel{objPropsModel}
		} else {
			state.ObjectProperties = nil
		}
	} else {
		state.ObjectProperties = nil
	}

	// Handle pbit array with mode awareness
	if utils.FieldAppliesToMode(resourceType, "pbit", mode) {
		if pbitArray, ok := data["pbit"].([]interface{}); ok && len(pbitArray) > 0 {
			var pbits []verityPacketQueuePbitModel
			for _, p := range pbitArray {
				pbit, ok := p.(map[string]interface{})
				if !ok {
					continue
				}
				pbitModel := verityPacketQueuePbitModel{
					PacketQueueForPBit: utils.MapInt64WithModeNested(pbit, "packet_queue_for_p_bit", resourceType, "pbit.packet_queue_for_p_bit", mode),
					Index:              utils.MapInt64WithModeNested(pbit, "index", resourceType, "pbit.index", mode),
				}
				pbits = append(pbits, pbitModel)
			}
			state.Pbit = pbits
		} else {
			state.Pbit = nil
		}
	} else {
		state.Pbit = nil
	}

	// Handle queue array with mode awareness
	if utils.FieldAppliesToMode(resourceType, "queue", mode) {
		if queueArray, ok := data["queue"].([]interface{}); ok && len(queueArray) > 0 {
			var queues []verityPacketQueueQueueModel
			for _, q := range queueArray {
				queue, ok := q.(map[string]interface{})
				if !ok {
					continue
				}
				queueModel := verityPacketQueueQueueModel{
					BandwidthForQueue: utils.MapInt64WithModeNested(queue, "bandwidth_for_queue", resourceType, "queue.bandwidth_for_queue", mode),
					SchedulerType:     utils.MapStringWithModeNested(queue, "scheduler_type", resourceType, "queue.scheduler_type", mode),
					SchedulerWeight:   utils.MapInt64WithModeNested(queue, "scheduler_weight", resourceType, "queue.scheduler_weight", mode),
					Index:             utils.MapInt64WithModeNested(queue, "index", resourceType, "queue.index", mode),
				}
				queues = append(queues, queueModel)
			}
			state.Queue = queues
		} else {
			state.Queue = nil
		}
	} else {
		state.Queue = nil
	}

	return state
}

func (r *verityPacketQueueResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// =========================================================================
	// Skip if deleting
	// =========================================================================
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityPacketQueueResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Mode-aware field nullification
	// Set fields that don't apply to current mode to null to prevent
	// "known after apply" messages for irrelevant fields.
	// =========================================================================
	const resourceType = packetQueueResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyBools(
		"enable",
	)

	nullifier.NullifyNestedBlocks(
		"pbit", "queue", "object_properties",
	)

	// =========================================================================
	// Skip UPDATE-specific logic during CREATE
	// =========================================================================
	if req.State.Raw.IsNull() {
		return
	}

	// =========================================================================
	// UPDATE operation - get state and config
	// =========================================================================
	var state verityPacketQueueResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config verityPacketQueueResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Handle nullable fields in nested blocks
	// =========================================================================
	name := plan.Name.ValueString()
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, "verity_packet_queue", name)

	// Handle pbit block nullable fields
	for i, configItem := range config.Pbit {
		itemIndex := configItem.Index.ValueInt64()
		var stateItem *verityPacketQueuePbitModel
		for j := range state.Pbit {
			if state.Pbit[j].Index.ValueInt64() == itemIndex {
				stateItem = &state.Pbit[j]
				break
			}
		}

		if stateItem != nil {
			utils.HandleNullableNestedFields(utils.NullableNestedFieldsConfig{
				Ctx:             ctx,
				Plan:            &resp.Plan,
				ConfiguredAttrs: configuredAttrs,
				BlockType:       "pbit",
				BlockListPath:   "pbit",
				BlockListIndex:  i,
				Int64Fields: []utils.NullableNestedInt64Field{
					{BlockIndex: itemIndex, AttrName: "packet_queue_for_p_bit", ConfigVal: configItem.PacketQueueForPBit, StateVal: stateItem.PacketQueueForPBit},
				},
			})
		}
	}

	// Handle queue block nullable fields
	for i, configItem := range config.Queue {
		itemIndex := configItem.Index.ValueInt64()
		var stateItem *verityPacketQueueQueueModel
		for j := range state.Queue {
			if state.Queue[j].Index.ValueInt64() == itemIndex {
				stateItem = &state.Queue[j]
				break
			}
		}

		if stateItem != nil {
			utils.HandleNullableNestedFields(utils.NullableNestedFieldsConfig{
				Ctx:             ctx,
				Plan:            &resp.Plan,
				ConfiguredAttrs: configuredAttrs,
				BlockType:       "queue",
				BlockListPath:   "queue",
				BlockListIndex:  i,
				Int64Fields: []utils.NullableNestedInt64Field{
					{BlockIndex: itemIndex, AttrName: "bandwidth_for_queue", ConfigVal: configItem.BandwidthForQueue, StateVal: stateItem.BandwidthForQueue},
					{BlockIndex: itemIndex, AttrName: "scheduler_weight", ConfigVal: configItem.SchedulerWeight, StateVal: stateItem.SchedulerWeight},
				},
			})
		}
	}
}
