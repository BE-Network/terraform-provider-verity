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

	// Handle string fields
	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "Mode", APIField: &deviceSettingsProps.Mode, TFValue: plan.Mode},
		{FieldName: "SpanningTreePriority", APIField: &deviceSettingsProps.SpanningTreePriority, TFValue: plan.SpanningTreePriority},
		{FieldName: "PacketQueueId", APIField: &deviceSettingsProps.PacketQueueId, TFValue: plan.PacketQueueId},
		{FieldName: "PacketQueueIdRefType", APIField: &deviceSettingsProps.PacketQueueIdRefType, TFValue: plan.PacketQueueIdRefType},
	})

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &deviceSettingsProps.Enable, TFValue: plan.Enable},
		{FieldName: "Rocev2", APIField: &deviceSettingsProps.Rocev2, TFValue: plan.Rocev2},
		{FieldName: "CutThroughSwitching", APIField: &deviceSettingsProps.CutThroughSwitching, TFValue: plan.CutThroughSwitching},
		{FieldName: "DisableTcpUdpLearnedPacketAcceleration", APIField: &deviceSettingsProps.DisableTcpUdpLearnedPacketAcceleration, TFValue: plan.DisableTcpUdpLearnedPacketAcceleration},
	})

	// Handle int64 fields
	utils.SetInt64Fields([]utils.Int64FieldMapping{
		{FieldName: "ExternalBatteryPowerAvailable", APIField: &deviceSettingsProps.ExternalBatteryPowerAvailable, TFValue: plan.ExternalBatteryPowerAvailable},
		{FieldName: "ExternalPowerAvailable", APIField: &deviceSettingsProps.ExternalPowerAvailable, TFValue: plan.ExternalPowerAvailable},
	})

	// Handle float64 fields
	utils.SetFloat64Fields([]utils.Float64FieldMapping{
		{FieldName: "UsageThreshold", APIField: &deviceSettingsProps.UsageThreshold, TFValue: plan.UsageThreshold},
	})

	// Handle nullable int64 fields
	utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
		{FieldName: "SecurityAuditInterval", APIField: &deviceSettingsProps.SecurityAuditInterval, TFValue: plan.SecurityAuditInterval},
		{FieldName: "CommitToFlashInterval", APIField: &deviceSettingsProps.CommitToFlashInterval, TFValue: plan.CommitToFlashInterval},
		{FieldName: "HoldTimer", APIField: &deviceSettingsProps.HoldTimer, TFValue: plan.HoldTimer},
		{FieldName: "MacAgingTimerOverride", APIField: &deviceSettingsProps.MacAgingTimerOverride, TFValue: plan.MacAgingTimerOverride},
	})

	// Handle object properties
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

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "device_settings", name, *deviceSettingsProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Device Settings %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "device_settings")

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

	tflog.Debug(ctx, fmt.Sprintf("Fetching device settings for verification of %s", deviceSettingsName))

	type DeviceSettingsResponse struct {
		EthDeviceProfiles map[string]interface{} `json:"eth_device_profiles"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "device_settings", deviceSettingsName,
		func() (DeviceSettingsResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch device settings")
			respAPI, err := r.client.DeviceSettingsAPI.DevicesettingsGet(ctx).Execute()
			if err != nil {
				return DeviceSettingsResponse{}, fmt.Errorf("error reading Device Settings: %v", err)
			}
			defer respAPI.Body.Close()

			var res DeviceSettingsResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return DeviceSettingsResponse{}, fmt.Errorf("failed to decode Device Settings response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d device settings", len(res.EthDeviceProfiles)))
			return res, nil
		},
		getCachedResponse,
	)
	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Device Settings %s", deviceSettingsName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for device settings with name: %s", deviceSettingsName))

	deviceData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.EthDeviceProfiles,
		deviceSettingsName,
		func(data interface{}) (string, bool) {
			if device, ok := data.(map[string]interface{}); ok {
				if name, ok := device["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Device Settings with name '%s' not found in API response", deviceSettingsName))
		resp.State.RemoveResource(ctx)
		return
	}

	deviceMap, ok := deviceData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Device Settings Data",
			fmt.Sprintf("Device Settings data is not in expected format for %s", deviceSettingsName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found device settings '%s' under API key '%s'", deviceSettingsName, actualAPIName))

	state.Name = utils.MapStringFromAPI(deviceMap["name"])

	// Handle object properties
	if objProps, ok := deviceMap["object_properties"].(map[string]interface{}); ok {
		group := utils.MapStringFromAPI(objProps["group"])
		isDefault := utils.MapBoolFromAPI(objProps["isdefault"])
		if group.IsNull() {
			group = types.StringValue("")
		}
		state.ObjectProperties = []verityDeviceSettingsObjectPropertiesModel{
			{Group: group, IsDefault: isDefault},
		}
	} else {
		state.ObjectProperties = nil
	}

	// Map string fields
	stringFieldMappings := map[string]*types.String{
		"mode":                      &state.Mode,
		"spanning_tree_priority":    &state.SpanningTreePriority,
		"packet_queue_id":           &state.PacketQueueId,
		"packet_queue_id_ref_type_": &state.PacketQueueIdRefType,
	}

	for apiKey, stateField := range stringFieldMappings {
		*stateField = utils.MapStringFromAPI(deviceMap[apiKey])
	}

	// Map boolean fields
	boolFieldMappings := map[string]*types.Bool{
		"enable":                &state.Enable,
		"rocev2":                &state.Rocev2,
		"cut_through_switching": &state.CutThroughSwitching,
		"disable_tcp_udp_learned_packet_acceleration": &state.DisableTcpUdpLearnedPacketAcceleration,
	}

	for apiKey, stateField := range boolFieldMappings {
		*stateField = utils.MapBoolFromAPI(deviceMap[apiKey])
	}

	// Map int64 fields
	int64FieldMappings := map[string]*types.Int64{
		"external_battery_power_available": &state.ExternalBatteryPowerAvailable,
		"external_power_available":         &state.ExternalPowerAvailable,
	}

	for apiKey, stateField := range int64FieldMappings {
		*stateField = utils.MapInt64FromAPI(deviceMap[apiKey])
	}

	// Map nullable int64 fields
	nullableInt64FieldMappings := map[string]*types.Int64{
		"security_audit_interval":  &state.SecurityAuditInterval,
		"commit_to_flash_interval": &state.CommitToFlashInterval,
		"hold_timer":               &state.HoldTimer,
		"mac_aging_timer_override": &state.MacAgingTimerOverride,
	}

	for apiKey, stateField := range nullableInt64FieldMappings {
		*stateField = utils.MapNullableInt64FromAPI(deviceMap[apiKey])
	}

	// Map float64 fields
	float64FieldMappings := map[string]*types.Float64{
		"usage_threshold": &state.UsageThreshold,
	}

	for apiKey, stateField := range float64FieldMappings {
		*stateField = utils.MapFloat64FromAPI(deviceMap[apiKey])
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

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { deviceSettingsProps.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Mode, state.Mode, func(v *string) { deviceSettingsProps.Mode = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.SpanningTreePriority, state.SpanningTreePriority, func(v *string) { deviceSettingsProps.SpanningTreePriority = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { deviceSettingsProps.Enable = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.Rocev2, state.Rocev2, func(v *bool) { deviceSettingsProps.Rocev2 = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.CutThroughSwitching, state.CutThroughSwitching, func(v *bool) { deviceSettingsProps.CutThroughSwitching = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.DisableTcpUdpLearnedPacketAcceleration, state.DisableTcpUdpLearnedPacketAcceleration, func(v *bool) { deviceSettingsProps.DisableTcpUdpLearnedPacketAcceleration = v }, &hasChanges)

	// Handle int64 field changes
	utils.CompareAndSetInt64Field(plan.ExternalBatteryPowerAvailable, state.ExternalBatteryPowerAvailable, func(v *int32) { deviceSettingsProps.ExternalBatteryPowerAvailable = v }, &hasChanges)
	utils.CompareAndSetInt64Field(plan.ExternalPowerAvailable, state.ExternalPowerAvailable, func(v *int32) { deviceSettingsProps.ExternalPowerAvailable = v }, &hasChanges)

	// Handle float64 field changes
	utils.CompareAndSetFloat64Field(plan.UsageThreshold, state.UsageThreshold, func(v *float32) { deviceSettingsProps.UsageThreshold = v }, &hasChanges)

	// Handle nullable int64 field changes
	utils.CompareAndSetNullableInt64Field(plan.SecurityAuditInterval, state.SecurityAuditInterval, func(v *openapi.NullableInt32) { deviceSettingsProps.SecurityAuditInterval = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.CommitToFlashInterval, state.CommitToFlashInterval, func(v *openapi.NullableInt32) { deviceSettingsProps.CommitToFlashInterval = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.HoldTimer, state.HoldTimer, func(v *openapi.NullableInt32) { deviceSettingsProps.HoldTimer = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.MacAgingTimerOverride, state.MacAgingTimerOverride, func(v *openapi.NullableInt32) { deviceSettingsProps.MacAgingTimerOverride = *v }, &hasChanges)

	// Handle object properties
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

	// Handle PacketQueueId and PacketQueueIdRefType using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.PacketQueueId, state.PacketQueueId, plan.PacketQueueIdRefType, state.PacketQueueIdRefType,
		func(v *string) { deviceSettingsProps.PacketQueueId = v },
		func(v *string) { deviceSettingsProps.PacketQueueIdRefType = v },
		"packet_queue_id", "packet_queue_id_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "device_settings", name, deviceSettingsProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Device Settings %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "device_settings")
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

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "device_settings", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Device Settings %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "device_settings")
	resp.State.RemoveResource(ctx)
}

func (r *verityDeviceSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
