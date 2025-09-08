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
	_ resource.Resource                = &verityDeviceSettingsResource{}
	_ resource.ResourceWithConfigure   = &verityDeviceSettingsResource{}
	_ resource.ResourceWithImportState = &verityDeviceSettingsResource{}
)

func NewVerityDeviceSettingsResource() resource.Resource {
	return &verityDeviceSettingsResource{}
}

type verityDeviceSettingsResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *utils.BulkOperationManager
	notifyOperationAdded func()
}

type verityDeviceSettingsResourceModel struct {
	Name                                   types.String                                `tfsdk:"name"`
	Enable                                 types.Bool                                  `tfsdk:"enable"`
	Mode                                   types.String                                `tfsdk:"mode"`
	UsageThreshold                         types.Float64                               `tfsdk:"usage_threshold"`
	ExternalBatteryPowerAvailable          types.Int64                                 `tfsdk:"external_battery_power_available"`
	ExternalPowerAvailable                 types.Int64                                 `tfsdk:"external_power_available"`
	SecurityAuditInterval                  types.Int64                                 `tfsdk:"security_audit_interval"`
	CommitToFlashInterval                  types.Int64                                 `tfsdk:"commit_to_flash_interval"`
	Rocev2                                 types.Bool                                  `tfsdk:"rocev2"`
	CutThroughSwitching                    types.Bool                                  `tfsdk:"cut_through_switching"`
	HoldTimer                              types.Int64                                 `tfsdk:"hold_timer"`
	DisableTcpUdpLearnedPacketAcceleration types.Bool                                  `tfsdk:"disable_tcp_udp_learned_packet_acceleration"`
	MacAgingTimerOverride                  types.Int64                                 `tfsdk:"mac_aging_timer_override"`
	SpanningTreePriority                   types.String                                `tfsdk:"spanning_tree_priority"`
	PacketQueueId                          types.String                                `tfsdk:"packet_queue_id"`
	PacketQueueIdRefType                   types.String                                `tfsdk:"packet_queue_id_ref_type_"`
	ObjectProperties                       []verityDeviceSettingsObjectPropertiesModel `tfsdk:"object_properties"`
}

type verityDeviceSettingsObjectPropertiesModel struct {
	Group     types.String `tfsdk:"group"`
	IsDefault types.Bool   `tfsdk:"isdefault"`
}

func (r *verityDeviceSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_settings"
}

func (r *verityDeviceSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityDeviceSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity Device Settings (Eth Device Profile)",
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
			"mode": schema.StringAttribute{
				Description: "Mode",
				Optional:    true,
			},
			"usage_threshold": schema.Float64Attribute{
				Description: "Usage Threshold",
				Optional:    true,
			},
			"external_battery_power_available": schema.Int64Attribute{
				Description: "External Battery Power Available (maximum: 2000)",
				Optional:    true,
			},
			"external_power_available": schema.Int64Attribute{
				Description: "External Power Available (maximum: 2000)",
				Optional:    true,
			},
			"security_audit_interval": schema.Int64Attribute{
				Description: "Frequency in minutes of rereading this Switch running configuration and comparing it to expected values. If the value is blank, audit will use default switch settings. If the value is 0, audit will be turned off. (maximum: 1440)",
				Optional:    true,
			},
			"commit_to_flash_interval": schema.Int64Attribute{
				Description: "Frequency in minutes to write the Switch configuration to flash. If the value is blank, commit will use default switch settings. If the value is 0, commit will be turned off. (maximum: 1440)",
				Optional:    true,
			},
			"rocev2": schema.BoolAttribute{
				Description: "Enable RDMA over Converged Ethernet version 2 network protocol. Switches that are set to ROCE mode should already have their port breakouts set up and should not have any ports configured with LAGs.",
				Optional:    true,
			},
			"cut_through_switching": schema.BoolAttribute{
				Description: "Enable Cut-through Switching on all Switches",
				Optional:    true,
			},
			"hold_timer": schema.Int64Attribute{
				Description: "Hold Timer (maximum: 86400)",
				Optional:    true,
			},
			"disable_tcp_udp_learned_packet_acceleration": schema.BoolAttribute{
				Description: "Required for AVB, PTP and Cobranet Support",
				Optional:    true,
			},
			"mac_aging_timer_override": schema.Int64Attribute{
				Description: "Blank uses the Device's default; otherwise an integer between 1 to 1,000,000 seconds (minimum: 1, maximum: 1000000)",
				Optional:    true,
			},
			"spanning_tree_priority": schema.StringAttribute{
				Description: "STP per switch, priority are in 4096 increments, the lower the number, the higher the priority.",
				Optional:    true,
			},
			"packet_queue_id": schema.StringAttribute{
				Description: "Packet Queue for device",
				Optional:    true,
			},
			"packet_queue_id_ref_type_": schema.StringAttribute{
				Description: "Object type for packet_queue_id field",
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"object_properties": schema.ListNestedBlock{
				Description: "Object properties for the Device Settings",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"group": schema.StringAttribute{
							Description: "Group",
							Optional:    true,
						},
						"isdefault": schema.BoolAttribute{
							Description: "Default object.",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func (r *verityDeviceSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityDeviceSettingsResourceModel
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
	deviceSettingsProps := &openapi.DevicesettingsPutRequestEthDeviceProfilesValue{
		Name: openapi.PtrString(name),
	}

	if !plan.Enable.IsNull() {
		deviceSettingsProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
	}

	if !plan.Mode.IsNull() {
		deviceSettingsProps.Mode = openapi.PtrString(plan.Mode.ValueString())
	}

	if !plan.UsageThreshold.IsNull() {
		deviceSettingsProps.UsageThreshold = openapi.PtrFloat32(float32(plan.UsageThreshold.ValueFloat64()))
	}

	if !plan.ExternalBatteryPowerAvailable.IsNull() {
		deviceSettingsProps.ExternalBatteryPowerAvailable = openapi.PtrInt32(int32(plan.ExternalBatteryPowerAvailable.ValueInt64()))
	}

	if !plan.ExternalPowerAvailable.IsNull() {
		deviceSettingsProps.ExternalPowerAvailable = openapi.PtrInt32(int32(plan.ExternalPowerAvailable.ValueInt64()))
	}

	if !plan.SecurityAuditInterval.IsNull() {
		val := int32(plan.SecurityAuditInterval.ValueInt64())
		deviceSettingsProps.SecurityAuditInterval = *openapi.NewNullableInt32(&val)
	} else {
		deviceSettingsProps.SecurityAuditInterval = *openapi.NewNullableInt32(nil)
	}

	if !plan.CommitToFlashInterval.IsNull() {
		val := int32(plan.CommitToFlashInterval.ValueInt64())
		deviceSettingsProps.CommitToFlashInterval = *openapi.NewNullableInt32(&val)
	} else {
		deviceSettingsProps.CommitToFlashInterval = *openapi.NewNullableInt32(nil)
	}

	if !plan.Rocev2.IsNull() {
		deviceSettingsProps.Rocev2 = openapi.PtrBool(plan.Rocev2.ValueBool())
	}

	if !plan.CutThroughSwitching.IsNull() {
		deviceSettingsProps.CutThroughSwitching = openapi.PtrBool(plan.CutThroughSwitching.ValueBool())
	}

	if !plan.HoldTimer.IsNull() {
		val := int32(plan.HoldTimer.ValueInt64())
		deviceSettingsProps.HoldTimer = *openapi.NewNullableInt32(&val)
	} else {
		deviceSettingsProps.HoldTimer = *openapi.NewNullableInt32(nil)
	}

	if !plan.DisableTcpUdpLearnedPacketAcceleration.IsNull() {
		deviceSettingsProps.DisableTcpUdpLearnedPacketAcceleration = openapi.PtrBool(plan.DisableTcpUdpLearnedPacketAcceleration.ValueBool())
	}

	if !plan.MacAgingTimerOverride.IsNull() {
		val := int32(plan.MacAgingTimerOverride.ValueInt64())
		deviceSettingsProps.MacAgingTimerOverride = *openapi.NewNullableInt32(&val)
	} else {
		deviceSettingsProps.MacAgingTimerOverride = *openapi.NewNullableInt32(nil)
	}

	if !plan.SpanningTreePriority.IsNull() {
		deviceSettingsProps.SpanningTreePriority = openapi.PtrString(plan.SpanningTreePriority.ValueString())
	}

	if !plan.PacketQueueId.IsNull() {
		deviceSettingsProps.PacketQueueId = openapi.PtrString(plan.PacketQueueId.ValueString())
	}

	if !plan.PacketQueueIdRefType.IsNull() {
		deviceSettingsProps.PacketQueueIdRefType = openapi.PtrString(plan.PacketQueueIdRefType.ValueString())
	}

	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties{}
		if !op.Group.IsNull() {
			objProps.Group = openapi.PtrString(op.Group.ValueString())
		} else {
			objProps.Group = nil
		}
		if !op.IsDefault.IsNull() {
			objProps.Isdefault = openapi.PtrBool(op.IsDefault.ValueBool())
		} else {
			objProps.Isdefault = nil
		}
		deviceSettingsProps.ObjectProperties = &objProps
	}

	operationID := r.bulkOpsMgr.AddPut(ctx, "device_settings", name, *deviceSettingsProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for Device Settings creation operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Create Device Settings %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Device Settings %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "devicesettings")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
}

func (r *verityDeviceSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityDeviceSettingsResourceModel
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

	deviceSettingsName := state.Name.ValueString()

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("device_settings") {
		tflog.Info(ctx, fmt.Sprintf("Skipping Device Settings %s verification – trusting recent successful API operation", deviceSettingsName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching Device Settings for verification of %s", deviceSettingsName))

	type DeviceSettingsResponse struct {
		EthDeviceProfiles map[string]interface{} `json:"eth_device_profiles"`
	}

	var result DeviceSettingsResponse
	var err error
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		deviceSettingsData, fetchErr := getCachedResponse(ctx, r.provCtx, "devicesettings", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch Device Settings")
			respAPI, err := r.client.DeviceSettingsAPI.DevicesettingsGet(ctx).Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading Device Settings: %v", err)
			}
			defer respAPI.Body.Close()

			var res DeviceSettingsResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return nil, fmt.Errorf("failed to decode Device Settings response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d Device Settings", len(res.EthDeviceProfiles)))
			return res, nil
		})
		if fetchErr != nil {
			err = fetchErr
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch Device Settings on attempt %d, retrying in %v", attempt+1, sleepTime))
			time.Sleep(sleepTime)
			continue
		}
		result = deviceSettingsData.(DeviceSettingsResponse)
		break
	}
	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Device Settings %s", deviceSettingsName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for Device Settings with ID: %s", deviceSettingsName))
	var deviceData map[string]interface{}
	exists := false

	if data, ok := result.EthDeviceProfiles[deviceSettingsName].(map[string]interface{}); ok {
		deviceData = data
		exists = true
		tflog.Debug(ctx, fmt.Sprintf("Found Device Settings directly by ID: %s", deviceSettingsName))
	} else {
		for apiName, device := range result.EthDeviceProfiles {
			deviceProfile, ok := device.(map[string]interface{})
			if !ok {
				continue
			}

			if name, ok := deviceProfile["name"].(string); ok && name == deviceSettingsName {
				deviceData = deviceProfile
				deviceSettingsName = apiName
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found Device Settings with name '%s' under API key '%s'", name, apiName))
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Device Settings with ID '%s' not found in API response", deviceSettingsName))
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(fmt.Sprintf("%v", deviceData["name"]))

	if enable, ok := deviceData["enable"].(bool); ok {
		state.Enable = types.BoolValue(enable)
	} else {
		state.Enable = types.BoolNull()
	}

	if mode, ok := deviceData["mode"].(string); ok && mode != "" {
		state.Mode = types.StringValue(mode)
	} else {
		state.Mode = types.StringNull()
	}

	if usageThreshold, ok := deviceData["usage_threshold"].(float64); ok {
		state.UsageThreshold = types.Float64Value(usageThreshold)
	} else {
		state.UsageThreshold = types.Float64Null()
	}

	if extBatteryPower, ok := deviceData["external_battery_power_available"].(float64); ok {
		state.ExternalBatteryPowerAvailable = types.Int64Value(int64(extBatteryPower))
	} else {
		state.ExternalBatteryPowerAvailable = types.Int64Null()
	}

	if extPower, ok := deviceData["external_power_available"].(float64); ok {
		state.ExternalPowerAvailable = types.Int64Value(int64(extPower))
	} else {
		state.ExternalPowerAvailable = types.Int64Null()
	}

	if securityAudit, ok := deviceData["security_audit_interval"]; ok && securityAudit != nil {
		if val, ok := securityAudit.(float64); ok {
			state.SecurityAuditInterval = types.Int64Value(int64(val))
		}
	} else {
		state.SecurityAuditInterval = types.Int64Null()
	}

	if commitFlash, ok := deviceData["commit_to_flash_interval"]; ok && commitFlash != nil {
		if val, ok := commitFlash.(float64); ok {
			state.CommitToFlashInterval = types.Int64Value(int64(val))
		}
	} else {
		state.CommitToFlashInterval = types.Int64Null()
	}

	if rocev2, ok := deviceData["rocev2"].(bool); ok {
		state.Rocev2 = types.BoolValue(rocev2)
	} else {
		state.Rocev2 = types.BoolNull()
	}

	if cutThrough, ok := deviceData["cut_through_switching"].(bool); ok {
		state.CutThroughSwitching = types.BoolValue(cutThrough)
	} else {
		state.CutThroughSwitching = types.BoolNull()
	}

	if holdTimer, ok := deviceData["hold_timer"]; ok && holdTimer != nil {
		if val, ok := holdTimer.(float64); ok {
			state.HoldTimer = types.Int64Value(int64(val))
		}
	} else {
		state.HoldTimer = types.Int64Null()
	}

	if disableTcp, ok := deviceData["disable_tcp_udp_learned_packet_acceleration"].(bool); ok {
		state.DisableTcpUdpLearnedPacketAcceleration = types.BoolValue(disableTcp)
	} else {
		state.DisableTcpUdpLearnedPacketAcceleration = types.BoolNull()
	}

	if macAging, ok := deviceData["mac_aging_timer_override"]; ok && macAging != nil {
		if val, ok := macAging.(float64); ok {
			state.MacAgingTimerOverride = types.Int64Value(int64(val))
		}
	} else {
		state.MacAgingTimerOverride = types.Int64Null()
	}

	if spanningTree, ok := deviceData["spanning_tree_priority"].(string); ok && spanningTree != "" {
		state.SpanningTreePriority = types.StringValue(spanningTree)
	} else {
		state.SpanningTreePriority = types.StringNull()
	}

	if packetQueue, ok := deviceData["packet_queue_id"].(string); ok {
		state.PacketQueueId = types.StringValue(packetQueue)
	} else {
		state.PacketQueueId = types.StringNull()
	}

	if packetQueueRefType, ok := deviceData["packet_queue_id_ref_type_"].(string); ok {
		state.PacketQueueIdRefType = types.StringValue(packetQueueRefType)
	} else {
		state.PacketQueueIdRefType = types.StringNull()
	}

	if objProps, ok := deviceData["object_properties"].(map[string]interface{}); ok {
		if group, ok := objProps["group"].(string); ok {
			if isDefault, ok := objProps["isdefault"].(bool); ok {
				state.ObjectProperties = []verityDeviceSettingsObjectPropertiesModel{
					{Group: types.StringValue(group), IsDefault: types.BoolValue(isDefault)},
				}
			} else {
				state.ObjectProperties = []verityDeviceSettingsObjectPropertiesModel{
					{Group: types.StringValue(group), IsDefault: types.BoolNull()},
				}
			}
		} else {
			if isDefault, ok := objProps["isdefault"].(bool); ok {
				state.ObjectProperties = []verityDeviceSettingsObjectPropertiesModel{
					{Group: types.StringValue(""), IsDefault: types.BoolValue(isDefault)},
				}
			} else {
				state.ObjectProperties = []verityDeviceSettingsObjectPropertiesModel{
					{Group: types.StringValue(""), IsDefault: types.BoolNull()},
				}
			}
		}
	} else {
		state.ObjectProperties = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityDeviceSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityDeviceSettingsResourceModel

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
	deviceSettingsProps := openapi.DevicesettingsPutRequestEthDeviceProfilesValue{}
	hasChanges := false

	if !plan.Enable.Equal(state.Enable) {
		deviceSettingsProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
		hasChanges = true
	}

	if !plan.Mode.Equal(state.Mode) {
		deviceSettingsProps.Mode = openapi.PtrString(plan.Mode.ValueString())
		hasChanges = true
	}

	if !plan.UsageThreshold.Equal(state.UsageThreshold) {
		deviceSettingsProps.UsageThreshold = openapi.PtrFloat32(float32(plan.UsageThreshold.ValueFloat64()))
		hasChanges = true
	}

	if !plan.ExternalBatteryPowerAvailable.Equal(state.ExternalBatteryPowerAvailable) {
		deviceSettingsProps.ExternalBatteryPowerAvailable = openapi.PtrInt32(int32(plan.ExternalBatteryPowerAvailable.ValueInt64()))
		hasChanges = true
	}

	if !plan.ExternalPowerAvailable.Equal(state.ExternalPowerAvailable) {
		deviceSettingsProps.ExternalPowerAvailable = openapi.PtrInt32(int32(plan.ExternalPowerAvailable.ValueInt64()))
		hasChanges = true
	}

	if !plan.SecurityAuditInterval.Equal(state.SecurityAuditInterval) {
		if !plan.SecurityAuditInterval.IsNull() {
			val := int32(plan.SecurityAuditInterval.ValueInt64())
			deviceSettingsProps.SecurityAuditInterval = *openapi.NewNullableInt32(&val)
		} else {
			deviceSettingsProps.SecurityAuditInterval = *openapi.NewNullableInt32(nil)
		}
		hasChanges = true
	}

	if !plan.CommitToFlashInterval.Equal(state.CommitToFlashInterval) {
		if !plan.CommitToFlashInterval.IsNull() {
			val := int32(plan.CommitToFlashInterval.ValueInt64())
			deviceSettingsProps.CommitToFlashInterval = *openapi.NewNullableInt32(&val)
		} else {
			deviceSettingsProps.CommitToFlashInterval = *openapi.NewNullableInt32(nil)
		}
		hasChanges = true
	}

	if !plan.Rocev2.Equal(state.Rocev2) {
		deviceSettingsProps.Rocev2 = openapi.PtrBool(plan.Rocev2.ValueBool())
		hasChanges = true
	}

	if !plan.CutThroughSwitching.Equal(state.CutThroughSwitching) {
		deviceSettingsProps.CutThroughSwitching = openapi.PtrBool(plan.CutThroughSwitching.ValueBool())
		hasChanges = true
	}

	if !plan.HoldTimer.Equal(state.HoldTimer) {
		if !plan.HoldTimer.IsNull() {
			val := int32(plan.HoldTimer.ValueInt64())
			deviceSettingsProps.HoldTimer = *openapi.NewNullableInt32(&val)
		} else {
			deviceSettingsProps.HoldTimer = *openapi.NewNullableInt32(nil)
		}
		hasChanges = true
	}

	if !plan.DisableTcpUdpLearnedPacketAcceleration.Equal(state.DisableTcpUdpLearnedPacketAcceleration) {
		deviceSettingsProps.DisableTcpUdpLearnedPacketAcceleration = openapi.PtrBool(plan.DisableTcpUdpLearnedPacketAcceleration.ValueBool())
		hasChanges = true
	}

	if !plan.MacAgingTimerOverride.Equal(state.MacAgingTimerOverride) {
		if !plan.MacAgingTimerOverride.IsNull() {
			val := int32(plan.MacAgingTimerOverride.ValueInt64())
			deviceSettingsProps.MacAgingTimerOverride = *openapi.NewNullableInt32(&val)
		} else {
			deviceSettingsProps.MacAgingTimerOverride = *openapi.NewNullableInt32(nil)
		}
		hasChanges = true
	}

	if !plan.SpanningTreePriority.Equal(state.SpanningTreePriority) {
		deviceSettingsProps.SpanningTreePriority = openapi.PtrString(plan.SpanningTreePriority.ValueString())
		hasChanges = true
	}

	// Handle PacketQueueId and PacketQueueIdRefType according to "One ref type supported" rules
	packetQueueChanged := !plan.PacketQueueId.Equal(state.PacketQueueId)
	packetQueueRefTypeChanged := !plan.PacketQueueIdRefType.Equal(state.PacketQueueIdRefType)

	// Case: Validate using "one ref type supported" rules
	if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
		plan.PacketQueueId, plan.PacketQueueIdRefType,
		"packet_queue_id", "packet_queue_id_ref_type_",
		packetQueueChanged, packetQueueRefTypeChanged) {
		return
	}

	if packetQueueChanged && !packetQueueRefTypeChanged {
		// Case: Only the base field changes, only the base field is sent
		// Just send the base field
		if !plan.PacketQueueId.IsNull() && plan.PacketQueueId.ValueString() != "" {
			deviceSettingsProps.PacketQueueId = openapi.PtrString(plan.PacketQueueId.ValueString())
		} else {
			deviceSettingsProps.PacketQueueId = openapi.PtrString("")
		}
		hasChanges = true
	} else if packetQueueRefTypeChanged {
		// Case: ref_type changes (or both change), both fields are sent

		// Send both fields
		if !plan.PacketQueueId.IsNull() && plan.PacketQueueId.ValueString() != "" {
			deviceSettingsProps.PacketQueueId = openapi.PtrString(plan.PacketQueueId.ValueString())
		} else {
			deviceSettingsProps.PacketQueueId = openapi.PtrString("")
		}

		if !plan.PacketQueueIdRefType.IsNull() && plan.PacketQueueIdRefType.ValueString() != "" {
			deviceSettingsProps.PacketQueueIdRefType = openapi.PtrString(plan.PacketQueueIdRefType.ValueString())
		} else {
			deviceSettingsProps.PacketQueueIdRefType = openapi.PtrString("")
		}

		hasChanges = true
	}

	if len(plan.ObjectProperties) > 0 {
		if len(state.ObjectProperties) == 0 ||
			!plan.ObjectProperties[0].Group.Equal(state.ObjectProperties[0].Group) ||
			!plan.ObjectProperties[0].IsDefault.Equal(state.ObjectProperties[0].IsDefault) {
			objProps := openapi.DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties{}
			if !plan.ObjectProperties[0].Group.IsNull() {
				objProps.Group = openapi.PtrString(plan.ObjectProperties[0].Group.ValueString())
			} else {
				objProps.Group = nil
			}
			if !plan.ObjectProperties[0].IsDefault.IsNull() {
				objProps.Isdefault = openapi.PtrBool(plan.ObjectProperties[0].IsDefault.ValueBool())
			} else {
				objProps.Isdefault = nil
			}
			deviceSettingsProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	operationID := r.bulkOpsMgr.AddPatch(ctx, "device_settings", name, deviceSettingsProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for Device Settings update operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update Device Settings %s", name))...,
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Device Settings %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "devicesettings")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *verityDeviceSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityDeviceSettingsResourceModel
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
	operationID := r.bulkOpsMgr.AddDelete(ctx, "device_settings", name)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for Device Settings deletion operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Delete Device Settings %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Device Settings %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "devicesettings")
	resp.State.RemoveResource(ctx)
}

func (r *verityDeviceSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
