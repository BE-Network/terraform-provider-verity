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
	_ resource.Resource                = &verityDeviceSettingsResource{}
	_ resource.ResourceWithConfigure   = &verityDeviceSettingsResource{}
	_ resource.ResourceWithImportState = &verityDeviceSettingsResource{}
	_ resource.ResourceWithModifyPlan  = &verityDeviceSettingsResource{}
)

const deviceSettingsResourceType = "devicesettings"

func NewVerityDeviceSettingsResource() resource.Resource {
	return &verityDeviceSettingsResource{}
}

type verityDeviceSettingsResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type verityDeviceSettingsResourceModel struct {
	Name                                   types.String                                `tfsdk:"name"`
	Enable                                 types.Bool                                  `tfsdk:"enable"`
	Mode                                   types.String                                `tfsdk:"mode"`
	UsageThreshold                         types.Number                                `tfsdk:"usage_threshold"`
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
	PacketQueue                            types.String                                `tfsdk:"packet_queue"`
	PacketQueueRefType                     types.String                                `tfsdk:"packet_queue_ref_type_"`
	ObjectProperties                       []verityDeviceSettingsObjectPropertiesModel `tfsdk:"object_properties"`
}

type verityDeviceSettingsObjectPropertiesModel struct {
	Group types.String `tfsdk:"group"`
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
				Computed:    true,
			},
			"mode": schema.StringAttribute{
				Description: "Mode",
				Optional:    true,
				Computed:    true,
			},
			"usage_threshold": schema.NumberAttribute{
				Description: "Usage Threshold",
				Optional:    true,
				Computed:    true,
			},
			"external_battery_power_available": schema.Int64Attribute{
				Description: "External Battery Power Available (maximum: 2000)",
				Optional:    true,
				Computed:    true,
			},
			"external_power_available": schema.Int64Attribute{
				Description: "External Power Available (maximum: 2000)",
				Optional:    true,
				Computed:    true,
			},
			"security_audit_interval": schema.Int64Attribute{
				Description: "Frequency in minutes of rereading this Switch running configuration and comparing it to expected values. If the value is blank, audit will use default switch settings. If the value is 0, audit will be turned off. (maximum: 1440)",
				Optional:    true,
				Computed:    true,
			},
			"commit_to_flash_interval": schema.Int64Attribute{
				Description: "Frequency in minutes to write the Switch configuration to flash. If the value is blank, commit will use default switch settings. If the value is 0, commit will be turned off. (maximum: 1440)",
				Optional:    true,
				Computed:    true,
			},
			"rocev2": schema.BoolAttribute{
				Description: "Enable RDMA over Converged Ethernet version 2 network protocol. Switches that are set to ROCE mode should already have their port breakouts set up and should not have any ports configured with LAGs.",
				Optional:    true,
				Computed:    true,
			},
			"cut_through_switching": schema.BoolAttribute{
				Description: "Enable Cut-through Switching on all Switches",
				Optional:    true,
				Computed:    true,
			},
			"hold_timer": schema.Int64Attribute{
				Description: "Hold Timer (maximum: 86400)",
				Optional:    true,
				Computed:    true,
			},
			"disable_tcp_udp_learned_packet_acceleration": schema.BoolAttribute{
				Description: "Required for AVB, PTP and Cobranet Support",
				Optional:    true,
				Computed:    true,
			},
			"mac_aging_timer_override": schema.Int64Attribute{
				Description: "Blank uses the Device's default; otherwise an integer between 1 to 1,000,000 seconds (minimum: 1, maximum: 1000000)",
				Optional:    true,
				Computed:    true,
			},
			"spanning_tree_priority": schema.StringAttribute{
				Description: "STP per switch, priority are in 4096 increments, the lower the number, the higher the priority.",
				Optional:    true,
				Computed:    true,
			},
			"packet_queue": schema.StringAttribute{
				Description: "Packet Queue for device",
				Optional:    true,
				Computed:    true,
			},
			"packet_queue_ref_type_": schema.StringAttribute{
				Description: "Object type for packet_queue field",
				Optional:    true,
				Computed:    true,
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
							Computed:    true,
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

	var config verityDeviceSettingsResourceModel
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
	deviceSettingsProps := &openapi.DevicesettingsPutRequestEthDeviceProfilesValue{
		Name: openapi.PtrString(name),
	}

	// Handle string fields
	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "Mode", APIField: &deviceSettingsProps.Mode, TFValue: plan.Mode},
		{FieldName: "SpanningTreePriority", APIField: &deviceSettingsProps.SpanningTreePriority, TFValue: plan.SpanningTreePriority},
		{FieldName: "PacketQueue", APIField: &deviceSettingsProps.PacketQueue, TFValue: plan.PacketQueue},
		{FieldName: "PacketQueueRefType", APIField: &deviceSettingsProps.PacketQueueRefType, TFValue: plan.PacketQueueRefType},
	})

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &deviceSettingsProps.Enable, TFValue: plan.Enable},
		{FieldName: "Rocev2", APIField: &deviceSettingsProps.Rocev2, TFValue: plan.Rocev2},
		{FieldName: "CutThroughSwitching", APIField: &deviceSettingsProps.CutThroughSwitching, TFValue: plan.CutThroughSwitching},
		{FieldName: "DisableTcpUdpLearnedPacketAcceleration", APIField: &deviceSettingsProps.DisableTcpUdpLearnedPacketAcceleration, TFValue: plan.DisableTcpUdpLearnedPacketAcceleration},
	})

	// Handle nullable int64 fields - parse HCL to detect explicit config
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, "verity_device_settings", name)

	utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
		{FieldName: "ExternalBatteryPowerAvailable", APIField: &deviceSettingsProps.ExternalBatteryPowerAvailable, TFValue: config.ExternalBatteryPowerAvailable, IsConfigured: configuredAttrs.IsConfigured("external_battery_power_available")},
		{FieldName: "ExternalPowerAvailable", APIField: &deviceSettingsProps.ExternalPowerAvailable, TFValue: config.ExternalPowerAvailable, IsConfigured: configuredAttrs.IsConfigured("external_power_available")},
		{FieldName: "SecurityAuditInterval", APIField: &deviceSettingsProps.SecurityAuditInterval, TFValue: config.SecurityAuditInterval, IsConfigured: configuredAttrs.IsConfigured("security_audit_interval")},
		{FieldName: "CommitToFlashInterval", APIField: &deviceSettingsProps.CommitToFlashInterval, TFValue: config.CommitToFlashInterval, IsConfigured: configuredAttrs.IsConfigured("commit_to_flash_interval")},
		{FieldName: "HoldTimer", APIField: &deviceSettingsProps.HoldTimer, TFValue: config.HoldTimer, IsConfigured: configuredAttrs.IsConfigured("hold_timer")},
		{FieldName: "MacAgingTimerOverride", APIField: &deviceSettingsProps.MacAgingTimerOverride, TFValue: config.MacAgingTimerOverride, IsConfigured: configuredAttrs.IsConfigured("mac_aging_timer_override")},
	})

	// Handle nullable float fields - parse HCL to detect explicit config
	utils.SetNullableNumberFields([]utils.NullableNumberFieldMapping{
		{FieldName: "UsageThreshold", APIField: &deviceSettingsProps.UsageThreshold, TFValue: config.UsageThreshold, IsConfigured: configuredAttrs.IsConfigured("usage_threshold")},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties{}
		utils.SetObjectPropertiesFields([]utils.ObjectPropertiesField{
			{Name: "Group", TFValue: op.Group, APIValue: &objProps.Group},
		})
		deviceSettingsProps.ObjectProperties = &objProps
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "device_settings", name, *deviceSettingsProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Device Settings %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "device_settings")

	var minState verityDeviceSettingsResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if deviceSettingsData, exists := bulkMgr.GetResourceResponse("device_settings", name); exists {
			state := populateDeviceSettingsState(ctx, minState, deviceSettingsData, r.provCtx.mode)
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

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if deviceSettingsData, exists := r.bulkOpsMgr.GetResourceResponse("device_settings", deviceSettingsName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached device settings data for %s from recent operation", deviceSettingsName))
			state = populateDeviceSettingsState(ctx, state, deviceSettingsData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

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

	state = populateDeviceSettingsState(ctx, state, deviceMap, r.provCtx.mode)
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

	// Handle nullable int64 field changes
	utils.CompareAndSetNullableInt64Field(plan.ExternalBatteryPowerAvailable, state.ExternalBatteryPowerAvailable, func(v *openapi.NullableInt32) { deviceSettingsProps.ExternalBatteryPowerAvailable = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.ExternalPowerAvailable, state.ExternalPowerAvailable, func(v *openapi.NullableInt32) { deviceSettingsProps.ExternalPowerAvailable = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.SecurityAuditInterval, state.SecurityAuditInterval, func(v *openapi.NullableInt32) { deviceSettingsProps.SecurityAuditInterval = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.CommitToFlashInterval, state.CommitToFlashInterval, func(v *openapi.NullableInt32) { deviceSettingsProps.CommitToFlashInterval = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.HoldTimer, state.HoldTimer, func(v *openapi.NullableInt32) { deviceSettingsProps.HoldTimer = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.MacAgingTimerOverride, state.MacAgingTimerOverride, func(v *openapi.NullableInt32) { deviceSettingsProps.MacAgingTimerOverride = *v }, &hasChanges)

	// Handle nullable float field changes
	utils.CompareAndSetNullableNumberField(plan.UsageThreshold, state.UsageThreshold, func(v *openapi.NullableFloat32) { deviceSettingsProps.UsageThreshold = *v }, &hasChanges)

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
			deviceSettingsProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	// Handle PacketQueue and PacketQueueRefType using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.PacketQueue, state.PacketQueue, plan.PacketQueueRefType, state.PacketQueueRefType,
		func(v *string) { deviceSettingsProps.PacketQueue = v },
		func(v *string) { deviceSettingsProps.PacketQueueRefType = v },
		"packet_queue", "packet_queue_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "device_settings", name, deviceSettingsProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Device Settings %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "device_settings")

	var minState verityDeviceSettingsResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Try to use cached response from bulk operation to populate state with API values
	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if deviceSettingsData, exists := bulkMgr.GetResourceResponse("device_settings", name); exists {
			newState := populateDeviceSettingsState(ctx, minState, deviceSettingsData, r.provCtx.mode)
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "device_settings", name, nil, &resp.Diagnostics)
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

func populateDeviceSettingsState(ctx context.Context, state verityDeviceSettingsResourceModel, data map[string]interface{}, mode string) verityDeviceSettingsResourceModel {
	const resourceType = deviceSettingsResourceType

	state.Name = utils.MapStringFromAPI(data["name"])

	// Boolean fields
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)
	state.Rocev2 = utils.MapBoolWithMode(data, "rocev2", resourceType, mode)
	state.CutThroughSwitching = utils.MapBoolWithMode(data, "cut_through_switching", resourceType, mode)
	state.DisableTcpUdpLearnedPacketAcceleration = utils.MapBoolWithMode(data, "disable_tcp_udp_learned_packet_acceleration", resourceType, mode)

	// String fields
	state.Mode = utils.MapStringWithMode(data, "mode", resourceType, mode)
	state.SpanningTreePriority = utils.MapStringWithMode(data, "spanning_tree_priority", resourceType, mode)
	state.PacketQueue = utils.MapStringWithMode(data, "packet_queue", resourceType, mode)
	state.PacketQueueRefType = utils.MapStringWithMode(data, "packet_queue_ref_type_", resourceType, mode)

	// Int64 fields
	state.ExternalBatteryPowerAvailable = utils.MapInt64WithMode(data, "external_battery_power_available", resourceType, mode)
	state.ExternalPowerAvailable = utils.MapInt64WithMode(data, "external_power_available", resourceType, mode)
	state.SecurityAuditInterval = utils.MapInt64WithMode(data, "security_audit_interval", resourceType, mode)
	state.CommitToFlashInterval = utils.MapInt64WithMode(data, "commit_to_flash_interval", resourceType, mode)
	state.HoldTimer = utils.MapInt64WithMode(data, "hold_timer", resourceType, mode)
	state.MacAgingTimerOverride = utils.MapInt64WithMode(data, "mac_aging_timer_override", resourceType, mode)

	// Float fields
	state.UsageThreshold = utils.MapNumberWithMode(data, "usage_threshold", resourceType, mode)

	// Handle object_properties block
	if utils.FieldAppliesToMode(resourceType, "object_properties", mode) {
		if objProps, ok := data["object_properties"].(map[string]interface{}); ok {
			objPropsModel := verityDeviceSettingsObjectPropertiesModel{
				Group: utils.MapStringWithModeNested(objProps, "group", resourceType, "object_properties.group", mode),
			}
			state.ObjectProperties = []verityDeviceSettingsObjectPropertiesModel{objPropsModel}
		} else {
			state.ObjectProperties = nil
		}
	} else {
		state.ObjectProperties = nil
	}

	return state
}

func (r *verityDeviceSettingsResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// =========================================================================
	// Skip if deleting
	// =========================================================================
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityDeviceSettingsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Mode-aware field nullification
	// Set fields that don't apply to current mode to null to prevent
	// "known after apply" messages for irrelevant fields.
	// =========================================================================
	const resourceType = deviceSettingsResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyStrings(
		"mode", "spanning_tree_priority", "packet_queue", "packet_queue_ref_type_",
	)

	nullifier.NullifyBools(
		"enable", "rocev2", "cut_through_switching", "disable_tcp_udp_learned_packet_acceleration",
	)

	nullifier.NullifyInt64s(
		"external_battery_power_available", "external_power_available",
		"security_audit_interval", "commit_to_flash_interval",
		"hold_timer", "mac_aging_timer_override",
	)

	// Float fields
	nullifier.NullifyNumbers(
		"usage_threshold",
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
	var state verityDeviceSettingsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config verityDeviceSettingsResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Handle nullable fields (explicit null detection)
	// For Optional+Computed fields, Terraform copies state to plan when config
	// is null. We detect explicit null in HCL and force plan to null.
	// =========================================================================
	name := plan.Name.ValueString()
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, "verity_device_settings", name)

	utils.HandleNullableFields(utils.NullableFieldsConfig{
		Ctx:             ctx,
		Plan:            &resp.Plan,
		ConfiguredAttrs: configuredAttrs,
		Int64Fields: []utils.NullableInt64Field{
			{AttrName: "external_battery_power_available", ConfigVal: config.ExternalBatteryPowerAvailable, StateVal: state.ExternalBatteryPowerAvailable},
			{AttrName: "external_power_available", ConfigVal: config.ExternalPowerAvailable, StateVal: state.ExternalPowerAvailable},
			{AttrName: "security_audit_interval", ConfigVal: config.SecurityAuditInterval, StateVal: state.SecurityAuditInterval},
			{AttrName: "commit_to_flash_interval", ConfigVal: config.CommitToFlashInterval, StateVal: state.CommitToFlashInterval},
			{AttrName: "hold_timer", ConfigVal: config.HoldTimer, StateVal: state.HoldTimer},
			{AttrName: "mac_aging_timer_override", ConfigVal: config.MacAgingTimerOverride, StateVal: state.MacAgingTimerOverride},
		},
		NumberFields: []utils.NullableNumberField{
			{AttrName: "usage_threshold", ConfigVal: config.UsageThreshold, StateVal: state.UsageThreshold},
		},
	})
}
