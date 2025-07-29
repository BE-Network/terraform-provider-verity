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

type verityPacketQueueQueueModel struct {
	BandwidthForQueue types.Int64  `tfsdk:"bandwidth_for_queue"`
	SchedulerType     types.String `tfsdk:"scheduler_type"`
	SchedulerWeight   types.Int64  `tfsdk:"scheduler_weight"`
	Index             types.Int64  `tfsdk:"index"`
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
	pqProps := openapi.ConfigPutRequestPacketQueuePacketQueueName{}
	pqProps.Name = openapi.PtrString(name)

	if !plan.Enable.IsNull() {
		pqProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
	}

	if len(plan.Pbit) > 0 {
		pbitArray := make([]openapi.ConfigPutRequestPacketQueuePacketQueueNamePbitInner, len(plan.Pbit))
		for i, pbit := range plan.Pbit {
			pbitItem := openapi.ConfigPutRequestPacketQueuePacketQueueNamePbitInner{}
			if !pbit.PacketQueueForPBit.IsNull() {
				pbitItem.PacketQueueForPBit = openapi.PtrInt32(int32(pbit.PacketQueueForPBit.ValueInt64()))
			}
			if !pbit.Index.IsNull() {
				pbitItem.Index = openapi.PtrInt32(int32(pbit.Index.ValueInt64()))
			}
			pbitArray[i] = pbitItem
		}
		pqProps.Pbit = pbitArray
	}

	if len(plan.Queue) > 0 {
		queueArray := make([]openapi.ConfigPutRequestPacketQueuePacketQueueNameQueueInner, len(plan.Queue))
		for i, queue := range plan.Queue {
			queueItem := openapi.ConfigPutRequestPacketQueuePacketQueueNameQueueInner{}
			if !queue.BandwidthForQueue.IsNull() {
				val := int32(queue.BandwidthForQueue.ValueInt64())
				queueItem.BandwidthForQueue = *openapi.NewNullableInt32(&val)
			} else {
				queueItem.BandwidthForQueue = *openapi.NewNullableInt32(nil)
			}
			if !queue.SchedulerType.IsNull() {
				queueItem.SchedulerType = openapi.PtrString(queue.SchedulerType.ValueString())
			}
			if !queue.SchedulerWeight.IsNull() {
				val := int32(queue.SchedulerWeight.ValueInt64())
				queueItem.SchedulerWeight = *openapi.NewNullableInt32(&val)
			} else {
				queueItem.SchedulerWeight = *openapi.NewNullableInt32(nil)
			}
			if !queue.Index.IsNull() {
				queueItem.Index = openapi.PtrInt32(int32(queue.Index.ValueInt64()))
			}
			queueArray[i] = queueItem
		}
		pqProps.Queue = queueArray
	}

	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.ConfigPutRequestPacketQueuePacketQueueNameObjectProperties{}
		if !op.IsDefault.IsNull() {
			objProps.Isdefault = openapi.PtrBool(op.IsDefault.ValueBool())
		}
		if !op.Group.IsNull() {
			objProps.Group = openapi.PtrString(op.Group.ValueString())
		}
		pqProps.ObjectProperties = &objProps
	}

	operationID := r.bulkOpsMgr.AddPut(ctx, "packetqueue", name, pqProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for packet queue creation operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Create Packet Queue %s", name))...,
		)
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

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("packetqueue") {
		tflog.Info(ctx, fmt.Sprintf("Skipping Packet Queue %s verification – trusting recent successful API operation", pqName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching Packet Queues for verification of %s", pqName))

	type PacketQueueResponse struct {
		PacketQueue map[string]interface{} `json:"packet_queue"`
	}

	var result PacketQueueResponse
	var err error
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		pqData, fetchErr := getCachedResponse(ctx, r.provCtx, "packet_queues", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch Packet Queues")
			respAPI, err := r.client.PacketQueuesAPI.PacketqueuesGet(ctx).Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading Packet Queues: %v", err)
			}
			defer respAPI.Body.Close()

			var res PacketQueueResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return nil, fmt.Errorf("failed to decode Packet Queues response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d Packet Queues", len(res.PacketQueue)))
			return res, nil
		})
		if fetchErr != nil {
			err = fetchErr
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch Packet Queues on attempt %d, retrying in %v", attempt+1, sleepTime))
			time.Sleep(sleepTime)
			continue
		}
		result = pqData.(PacketQueueResponse)
		break
	}
	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Packet Queue %s", pqName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for Packet Queue with ID: %s", pqName))
	var pqData map[string]interface{}
	exists := false

	if data, ok := result.PacketQueue[pqName].(map[string]interface{}); ok {
		pqData = data
		exists = true
		tflog.Debug(ctx, fmt.Sprintf("Found Packet Queue directly by ID: %s", pqName))
	} else {
		for apiName, p := range result.PacketQueue {
			queue, ok := p.(map[string]interface{})
			if !ok {
				continue
			}

			if name, ok := queue["name"].(string); ok && name == pqName {
				pqData = queue
				pqName = apiName
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found Packet Queue with name '%s' under API key '%s'", name, apiName))
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Packet Queue with ID '%s' not found in API response", pqName))
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(fmt.Sprintf("%v", pqData["name"]))

	if enable, ok := pqData["enable"].(bool); ok {
		state.Enable = types.BoolValue(enable)
	} else {
		state.Enable = types.BoolNull()
	}

	if pbitArray, ok := pqData["pbit"].([]interface{}); ok && len(pbitArray) > 0 {
		var pbits []verityPacketQueuePbitModel
		for _, p := range pbitArray {
			pbit, ok := p.(map[string]interface{})
			if !ok {
				continue
			}
			pbitModel := verityPacketQueuePbitModel{}
			if val, ok := pbit["packet_queue_for_p_bit"]; ok && val != nil {
				if intVal, ok := val.(float64); ok {
					pbitModel.PacketQueueForPBit = types.Int64Value(int64(intVal))
				} else if intVal, ok := val.(int); ok {
					pbitModel.PacketQueueForPBit = types.Int64Value(int64(intVal))
				} else {
					pbitModel.PacketQueueForPBit = types.Int64Null()
				}
			} else {
				pbitModel.PacketQueueForPBit = types.Int64Null()
			}
			if index, ok := pbit["index"]; ok && index != nil {
				if intVal, ok := index.(float64); ok {
					pbitModel.Index = types.Int64Value(int64(intVal))
				} else if intVal, ok := index.(int); ok {
					pbitModel.Index = types.Int64Value(int64(intVal))
				} else {
					pbitModel.Index = types.Int64Null()
				}
			} else {
				pbitModel.Index = types.Int64Null()
			}
			pbits = append(pbits, pbitModel)
		}
		state.Pbit = pbits
	} else {
		state.Pbit = nil
	}

	if queueArray, ok := pqData["queue"].([]interface{}); ok && len(queueArray) > 0 {
		var queues []verityPacketQueueQueueModel
		for _, q := range queueArray {
			queue, ok := q.(map[string]interface{})
			if !ok {
				continue
			}
			queueModel := verityPacketQueueQueueModel{}
			if val, ok := queue["bandwidth_for_queue"]; ok && val != nil {
				if intVal, ok := val.(float64); ok {
					queueModel.BandwidthForQueue = types.Int64Value(int64(intVal))
				} else if intVal, ok := val.(int); ok {
					queueModel.BandwidthForQueue = types.Int64Value(int64(intVal))
				} else {
					queueModel.BandwidthForQueue = types.Int64Null()
				}
			} else {
				queueModel.BandwidthForQueue = types.Int64Null()
			}
			if schedulerType, ok := queue["scheduler_type"].(string); ok {
				queueModel.SchedulerType = types.StringValue(schedulerType)
			} else {
				queueModel.SchedulerType = types.StringNull()
			}
			if val, ok := queue["scheduler_weight"]; ok && val != nil {
				if intVal, ok := val.(float64); ok {
					queueModel.SchedulerWeight = types.Int64Value(int64(intVal))
				} else if intVal, ok := val.(int); ok {
					queueModel.SchedulerWeight = types.Int64Value(int64(intVal))
				} else {
					queueModel.SchedulerWeight = types.Int64Null()
				}
			} else {
				queueModel.SchedulerWeight = types.Int64Null()
			}
			if index, ok := queue["index"]; ok && index != nil {
				if intVal, ok := index.(float64); ok {
					queueModel.Index = types.Int64Value(int64(intVal))
				} else if intVal, ok := index.(int); ok {
					queueModel.Index = types.Int64Value(int64(intVal))
				} else {
					queueModel.Index = types.Int64Null()
				}
			} else {
				queueModel.Index = types.Int64Null()
			}
			queues = append(queues, queueModel)
		}
		state.Queue = queues
	} else {
		state.Queue = nil
	}

	if objProps, ok := pqData["object_properties"].(map[string]interface{}); ok {
		op := verityPacketQueueObjectPropertiesModel{}
		if isDefault, ok := objProps["isdefault"].(bool); ok {
			op.IsDefault = types.BoolValue(isDefault)
		} else {
			op.IsDefault = types.BoolNull()
		}
		if group, ok := objProps["group"].(string); ok {
			op.Group = types.StringValue(group)
		} else {
			op.Group = types.StringNull()
		}
		state.ObjectProperties = []verityPacketQueueObjectPropertiesModel{op}
	} else {
		state.ObjectProperties = nil
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
	pqProps := openapi.ConfigPutRequestPacketQueuePacketQueueName{}
	hasChanges := false

	if !plan.Name.Equal(state.Name) {
		pqProps.Name = openapi.PtrString(name)
		hasChanges = true
	}

	if !plan.Enable.Equal(state.Enable) {
		pqProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
		hasChanges = true
	}

	if !r.equalPbitArrays(plan.Pbit, state.Pbit) {
		pbitArray := make([]openapi.ConfigPutRequestPacketQueuePacketQueueNamePbitInner, len(plan.Pbit))
		for i, pbit := range plan.Pbit {
			pbitItem := openapi.ConfigPutRequestPacketQueuePacketQueueNamePbitInner{}
			if !pbit.PacketQueueForPBit.IsNull() {
				pbitItem.PacketQueueForPBit = openapi.PtrInt32(int32(pbit.PacketQueueForPBit.ValueInt64()))
			}
			if !pbit.Index.IsNull() {
				pbitItem.Index = openapi.PtrInt32(int32(pbit.Index.ValueInt64()))
			}
			pbitArray[i] = pbitItem
		}
		pqProps.Pbit = pbitArray
		hasChanges = true
	}

	if !r.equalQueueArrays(plan.Queue, state.Queue) {
		queueArray := make([]openapi.ConfigPutRequestPacketQueuePacketQueueNameQueueInner, len(plan.Queue))
		for i, queue := range plan.Queue {
			queueItem := openapi.ConfigPutRequestPacketQueuePacketQueueNameQueueInner{}
			if !queue.BandwidthForQueue.IsNull() {
				val := int32(queue.BandwidthForQueue.ValueInt64())
				queueItem.BandwidthForQueue = *openapi.NewNullableInt32(&val)
			} else {
				queueItem.BandwidthForQueue = *openapi.NewNullableInt32(nil)
			}
			if !queue.SchedulerType.IsNull() {
				queueItem.SchedulerType = openapi.PtrString(queue.SchedulerType.ValueString())
			}
			if !queue.SchedulerWeight.IsNull() {
				val := int32(queue.SchedulerWeight.ValueInt64())
				queueItem.SchedulerWeight = *openapi.NewNullableInt32(&val)
			} else {
				queueItem.SchedulerWeight = *openapi.NewNullableInt32(nil)
			}
			if !queue.Index.IsNull() {
				queueItem.Index = openapi.PtrInt32(int32(queue.Index.ValueInt64()))
			}
			queueArray[i] = queueItem
		}
		pqProps.Queue = queueArray
		hasChanges = true
	}

	if len(plan.ObjectProperties) > 0 {
		if len(state.ObjectProperties) == 0 ||
			!plan.ObjectProperties[0].IsDefault.Equal(state.ObjectProperties[0].IsDefault) ||
			!plan.ObjectProperties[0].Group.Equal(state.ObjectProperties[0].Group) {
			objProps := openapi.ConfigPutRequestPacketQueuePacketQueueNameObjectProperties{}
			if !plan.ObjectProperties[0].IsDefault.IsNull() {
				objProps.Isdefault = openapi.PtrBool(plan.ObjectProperties[0].IsDefault.ValueBool())
			}
			if !plan.ObjectProperties[0].Group.IsNull() {
				objProps.Group = openapi.PtrString(plan.ObjectProperties[0].Group.ValueString())
			}
			pqProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	operationID := r.bulkOpsMgr.AddPatch(ctx, "packetqueue", name, pqProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for Packet Queue update operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update Packet Queue %s", name))...,
		)
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
	operationID := r.bulkOpsMgr.AddDelete(ctx, "packetqueue", name)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for Packet Queue deletion operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Delete Packet Queue %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Packet Queue %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "packet_queues")
	resp.State.RemoveResource(ctx)
}

func (r *verityPacketQueueResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func (r *verityPacketQueueResource) equalPbitArrays(a, b []verityPacketQueuePbitModel) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !a[i].PacketQueueForPBit.Equal(b[i].PacketQueueForPBit) ||
			!a[i].Index.Equal(b[i].Index) {
			return false
		}
	}
	return true
}

func (r *verityPacketQueueResource) equalQueueArrays(a, b []verityPacketQueueQueueModel) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !a[i].BandwidthForQueue.Equal(b[i].BandwidthForQueue) ||
			!a[i].SchedulerType.Equal(b[i].SchedulerType) ||
			!a[i].SchedulerWeight.Equal(b[i].SchedulerWeight) ||
			!a[i].Index.Equal(b[i].Index) {
			return false
		}
	}
	return true
}
