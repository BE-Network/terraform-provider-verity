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
	_ resource.Resource                = &verityPacketQueueResource{}
	_ resource.ResourceWithConfigure   = &verityPacketQueueResource{}
	_ resource.ResourceWithImportState = &verityPacketQueueResource{}
)

func NewVerityPacketQueueResource() resource.Resource {
	return &verityPacketQueueResource{}
}

type verityPacketQueueResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *utils.BulkOperationManager
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
	IsDefault types.Bool   `tfsdk:"isdefault"`
	Group     types.String `tfsdk:"group"`
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
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object. Zero if you want to add an object to the list.",
							Optional:    true,
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
						},
						"scheduler_type": schema.StringAttribute{
							Description: "Scheduler Type for Queue",
							Optional:    true,
						},
						"scheduler_weight": schema.Int64Attribute{
							Description: "Weight associated with WRR or DWRR scheduler",
							Optional:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object. Zero if you want to add an object to the list.",
							Optional:    true,
						},
					},
				},
			},
			"object_properties": schema.ListNestedBlock{
				Description: "Object properties for the packet queue",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"isdefault": schema.BoolAttribute{
							Description: "Default object.",
							Optional:    true,
						},
						"group": schema.StringAttribute{
							Description: "Group",
							Optional:    true,
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
		objProps := openapi.PacketqueuesPutRequestPacketQueueValueObjectProperties{}
		if !op.IsDefault.IsNull() {
			objProps.Isdefault = openapi.PtrBool(op.IsDefault.ValueBool())
		} else {
			objProps.Isdefault = nil
		}
		if !op.Group.IsNull() {
			objProps.Group = openapi.PtrString(op.Group.ValueString())
		} else {
			objProps.Group = nil
		}
		pqProps.ObjectProperties = &objProps
	}

	// Handle Pbit
	if len(plan.Pbit) > 0 {
		pbitArray := make([]openapi.PacketqueuesPutRequestPacketQueueValuePbitInner, len(plan.Pbit))
		for i, pbit := range plan.Pbit {
			pbitItem := openapi.PacketqueuesPutRequestPacketQueueValuePbitInner{}

			// Handle int64 fields
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &pbitItem.Index, TFValue: pbit.Index},
				{FieldName: "PacketQueueForPBit", APIField: &pbitItem.PacketQueueForPBit, TFValue: pbit.PacketQueueForPBit},
			})

			pbitArray[i] = pbitItem
		}
		pqProps.Pbit = pbitArray
	}

	// Handle Queue
	if len(plan.Queue) > 0 {
		queueArray := make([]openapi.PacketqueuesPutRequestPacketQueueValueQueueInner, len(plan.Queue))
		for i, queue := range plan.Queue {
			queueItem := openapi.PacketqueuesPutRequestPacketQueueValueQueueInner{}

			// Handle int64 fields
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &queueItem.Index, TFValue: queue.Index},
			})

			// Handle nullable int64 fields
			utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
				{FieldName: "BandwidthForQueue", APIField: &queueItem.BandwidthForQueue, TFValue: queue.BandwidthForQueue},
				{FieldName: "SchedulerWeight", APIField: &queueItem.SchedulerWeight, TFValue: queue.SchedulerWeight},
			})

			// Handle string fields
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "SchedulerType", APIField: &queueItem.SchedulerType, TFValue: queue.SchedulerType},
			})

			queueArray[i] = queueItem
		}
		pqProps.Queue = queueArray
	}

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "packet_queue", name, *pqProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Packet Queue %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "packet_queues")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
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

	state.Name = utils.MapStringFromAPI(pqMap["name"])

	// Handle object properties
	if objProps, ok := pqMap["object_properties"].(map[string]interface{}); ok {
		state.ObjectProperties = []verityPacketQueueObjectPropertiesModel{
			{
				IsDefault: utils.MapBoolFromAPI(objProps["isdefault"]),
				Group:     utils.MapStringFromAPI(objProps["group"]),
			},
		}
	} else {
		state.ObjectProperties = nil
	}

	// Map boolean fields
	boolFieldMappings := map[string]*types.Bool{
		"enable": &state.Enable,
	}

	for apiKey, stateField := range boolFieldMappings {
		*stateField = utils.MapBoolFromAPI(pqMap[apiKey])
	}

	// Handle pbit array
	if pbitArray, ok := pqMap["pbit"].([]interface{}); ok && len(pbitArray) > 0 {
		var pbits []verityPacketQueuePbitModel
		for _, p := range pbitArray {
			pbit, ok := p.(map[string]interface{})
			if !ok {
				continue
			}
			pbitModel := verityPacketQueuePbitModel{
				PacketQueueForPBit: utils.MapInt64FromAPI(pbit["packet_queue_for_p_bit"]),
				Index:              utils.MapInt64FromAPI(pbit["index"]),
			}
			pbits = append(pbits, pbitModel)
		}
		state.Pbit = pbits
	} else {
		state.Pbit = nil
	}

	// Handle queue array
	if queueArray, ok := pqMap["queue"].([]interface{}); ok && len(queueArray) > 0 {
		var queues []verityPacketQueueQueueModel
		for _, q := range queueArray {
			queue, ok := q.(map[string]interface{})
			if !ok {
				continue
			}
			queueModel := verityPacketQueueQueueModel{
				BandwidthForQueue: utils.MapInt64FromAPI(queue["bandwidth_for_queue"]),
				SchedulerType:     utils.MapStringFromAPI(queue["scheduler_type"]),
				SchedulerWeight:   utils.MapInt64FromAPI(queue["scheduler_weight"]),
				Index:             utils.MapInt64FromAPI(queue["index"]),
			}
			queues = append(queues, queueModel)
		}
		state.Queue = queues
	} else {
		state.Queue = nil
	}

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
	if len(plan.ObjectProperties) > 0 {
		if len(state.ObjectProperties) == 0 ||
			!plan.ObjectProperties[0].IsDefault.Equal(state.ObjectProperties[0].IsDefault) ||
			!plan.ObjectProperties[0].Group.Equal(state.ObjectProperties[0].Group) {
			objProps := openapi.PacketqueuesPutRequestPacketQueueValueObjectProperties{}
			if !plan.ObjectProperties[0].IsDefault.IsNull() {
				objProps.Isdefault = openapi.PtrBool(plan.ObjectProperties[0].IsDefault.ValueBool())
			} else {
				objProps.Isdefault = nil
			}
			if !plan.ObjectProperties[0].Group.IsNull() {
				objProps.Group = openapi.PtrString(plan.ObjectProperties[0].Group.ValueString())
			} else {
				objProps.Group = nil
			}
			pqProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	// Handle Pbit
	changedPbits, pbitsChanged := utils.ProcessIndexedArrayUpdates(plan.Pbit, state.Pbit,
		utils.IndexedItemHandler[verityPacketQueuePbitModel, openapi.PacketqueuesPutRequestPacketQueueValuePbitInner]{
			CreateNew: func(planItem verityPacketQueuePbitModel) openapi.PacketqueuesPutRequestPacketQueueValuePbitInner {
				newPbit := openapi.PacketqueuesPutRequestPacketQueueValuePbitInner{}

				// Handle int64 fields
				utils.SetInt64Fields([]utils.Int64FieldMapping{
					{FieldName: "Index", APIField: &newPbit.Index, TFValue: planItem.Index},
					{FieldName: "PacketQueueForPBit", APIField: &newPbit.PacketQueueForPBit, TFValue: planItem.PacketQueueForPBit},
				})

				return newPbit
			},
			UpdateExisting: func(planItem verityPacketQueuePbitModel, stateItem verityPacketQueuePbitModel) (openapi.PacketqueuesPutRequestPacketQueueValuePbitInner, bool) {
				updatePbit := openapi.PacketqueuesPutRequestPacketQueueValuePbitInner{}
				fieldChanged := false

				// Handle int64 field changes
				utils.CompareAndSetInt64Field(planItem.Index, stateItem.Index, func(v *int32) { updatePbit.Index = v }, &fieldChanged)
				utils.CompareAndSetInt64Field(planItem.PacketQueueForPBit, stateItem.PacketQueueForPBit, func(v *int32) { updatePbit.PacketQueueForPBit = v }, &fieldChanged)

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

				// Handle nullable int64 fields
				utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
					{FieldName: "BandwidthForQueue", APIField: &newQueue.BandwidthForQueue, TFValue: planItem.BandwidthForQueue},
					{FieldName: "SchedulerWeight", APIField: &newQueue.SchedulerWeight, TFValue: planItem.SchedulerWeight},
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

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "packet_queue", name, pqProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Packet Queue %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "packet_queues")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
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

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "packet_queue", name, nil, &resp.Diagnostics)
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
