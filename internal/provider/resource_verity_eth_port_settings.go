package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-verity/internal/utils"
	"terraform-provider-verity/openapi"
)

var (
	_ resource.Resource                = &verityEthPortSettingsResource{}
	_ resource.ResourceWithConfigure   = &verityEthPortSettingsResource{}
	_ resource.ResourceWithImportState = &verityEthPortSettingsResource{}
)

func NewVerityEthPortSettingsResource() resource.Resource {
	return &verityEthPortSettingsResource{}
}

type verityEthPortSettingsResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *utils.BulkOperationManager
	notifyOperationAdded func()
}

type verityEthPortSettingsResourceModel struct {
	Name             types.String                                 `tfsdk:"name"`
	Enable           types.Bool                                   `tfsdk:"enable"`
	ObjectProperties []verityEthPortSettingsObjectPropertiesModel `tfsdk:"object_properties"`
	AutoNegotiation  types.Bool                                   `tfsdk:"auto_negotiation"`
	MaxBitRate       types.String                                 `tfsdk:"max_bit_rate"`
	DuplexMode       types.String                                 `tfsdk:"duplex_mode"`
	StpEnable        types.Bool                                   `tfsdk:"stp_enable"`
	FastLearningMode types.Bool                                   `tfsdk:"fast_learning_mode"`
	BpduGuard        types.Bool                                   `tfsdk:"bpdu_guard"`
	BpduFilter       types.Bool                                   `tfsdk:"bpdu_filter"`
	GuardLoop        types.Bool                                   `tfsdk:"guard_loop"`
	PoeEnable        types.Bool                                   `tfsdk:"poe_enable"`
	Priority         types.String                                 `tfsdk:"priority"`
	AllocatedPower   types.String                                 `tfsdk:"allocated_power"`
	BspEnable        types.Bool                                   `tfsdk:"bsp_enable"`
	Broadcast        types.Bool                                   `tfsdk:"broadcast"`
	Multicast        types.Bool                                   `tfsdk:"multicast"`
	MaxAllowedValue  types.Int64                                  `tfsdk:"max_allowed_value"`
	MaxAllowedUnit   types.String                                 `tfsdk:"max_allowed_unit"`
	Action           types.String                                 `tfsdk:"action"`
	Fec              types.String                                 `tfsdk:"fec"`
	SingleLink       types.Bool                                   `tfsdk:"single_link"`
}

type verityEthPortSettingsObjectPropertiesModel struct {
	Group types.String `tfsdk:"group"`
}

func (r *verityEthPortSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_eth_port_settings"
}

func (r *verityEthPortSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityEthPortSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages Ethernet port settings",
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
				Default:     booldefault.StaticBool(false),
			},
			"auto_negotiation": schema.BoolAttribute{
				Description: "Indicates if port speed and duplex mode should be auto negotiated",
				Optional:    true,
			},
			"max_bit_rate": schema.StringAttribute{
				Description: "Maximum Bit Rate allowed",
				Optional:    true,
			},
			"duplex_mode": schema.StringAttribute{
				Description: "Duplex Mode",
				Optional:    true,
			},
			"stp_enable": schema.BoolAttribute{
				Description: "Enable Spanning Tree on the port. Note: the Spanning Tree Type (VLAN, Port, MST) is controlled in the Site Settings",
				Optional:    true,
			},
			"fast_learning_mode": schema.BoolAttribute{
				Description: "Enable Immediate Transition to Forwarding",
				Optional:    true,
			},
			"bpdu_guard": schema.BoolAttribute{
				Description: "Block port on BPDU Receive",
				Optional:    true,
			},
			"bpdu_filter": schema.BoolAttribute{
				Description: "Drop all Rx and Tx BPDUs",
				Optional:    true,
			},
			"guard_loop": schema.BoolAttribute{
				Description: "Enable Cisco Guard Loop",
				Optional:    true,
			},
			"poe_enable": schema.BoolAttribute{
				Description: "PoE Enable",
				Optional:    true,
			},
			"priority": schema.StringAttribute{
				Description: "Priority given when assigning power in a limited power situation",
				Optional:    true,
			},
			"allocated_power": schema.StringAttribute{
				Description: "Power the PoE system will attempt to allocate on this port",
				Optional:    true,
			},
			"bsp_enable": schema.BoolAttribute{
				Description: "Enable Traffic Storm Protection which prevents excessive broadcast/multicast/unknown-unicast traffic from overwhelming the Switch CPU",
				Optional:    true,
			},
			"broadcast": schema.BoolAttribute{
				Description: "Broadcast",
				Optional:    true,
			},
			"multicast": schema.BoolAttribute{
				Description: "Multicast",
				Optional:    true,
			},
			"max_allowed_value": schema.Int64Attribute{
				Description: "Max Percentage of the port's bandwidth allowed for broadcast/multicast/unknown-unicast traffic before invoking the protective action",
				Optional:    true,
			},
			"max_allowed_unit": schema.StringAttribute{
				Description: "Max Percentage unit for broadcast/multicast/unknown-unicast traffic",
				Optional:    true,
			},
			"action": schema.StringAttribute{
				Description: "Action taken if broadcast/multicast/unknown-unicast traffic exceeds the Max",
				Optional:    true,
			},
			"fec": schema.StringAttribute{
				Description: "FEC is Forward Error Correction which is error correction on the fiber link",
				Optional:    true,
			},
			"single_link": schema.BoolAttribute{
				Description: "Ports with this setting will be disabled when link state tracking takes effect",
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"object_properties": schema.ListNestedBlock{
				Description: "Object properties for the eth port settings",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
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

func (r *verityEthPortSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityEthPortSettingsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Authenticate",
			fmt.Sprintf("Error authenticating with API: %v", err),
		)
		return
	}

	name := plan.Name.ValueString()

	ethPortSettingsReq := openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName{}
	ethPortSettingsReq.Name = openapi.PtrString(name)
	if !plan.Enable.IsNull() {
		ethPortSettingsReq.Enable = openapi.PtrBool(plan.Enable.ValueBool())
	}
	if !plan.AutoNegotiation.IsNull() {
		ethPortSettingsReq.AutoNegotiation = openapi.PtrBool(plan.AutoNegotiation.ValueBool())
	}
	if !plan.MaxBitRate.IsNull() {
		ethPortSettingsReq.MaxBitRate = openapi.PtrString(plan.MaxBitRate.ValueString())
	}
	if !plan.DuplexMode.IsNull() {
		ethPortSettingsReq.DuplexMode = openapi.PtrString(plan.DuplexMode.ValueString())
	}
	if !plan.StpEnable.IsNull() {
		ethPortSettingsReq.StpEnable = openapi.PtrBool(plan.StpEnable.ValueBool())
	}
	if !plan.FastLearningMode.IsNull() {
		ethPortSettingsReq.FastLearningMode = openapi.PtrBool(plan.FastLearningMode.ValueBool())
	}
	if !plan.BpduGuard.IsNull() {
		ethPortSettingsReq.BpduGuard = openapi.PtrBool(plan.BpduGuard.ValueBool())
	}
	if !plan.BpduFilter.IsNull() {
		ethPortSettingsReq.BpduFilter = openapi.PtrBool(plan.BpduFilter.ValueBool())
	}
	if !plan.GuardLoop.IsNull() {
		ethPortSettingsReq.GuardLoop = openapi.PtrBool(plan.GuardLoop.ValueBool())
	}
	if !plan.PoeEnable.IsNull() {
		ethPortSettingsReq.PoeEnable = openapi.PtrBool(plan.PoeEnable.ValueBool())
	}
	if !plan.Priority.IsNull() {
		ethPortSettingsReq.Priority = openapi.PtrString(plan.Priority.ValueString())
	}
	if !plan.AllocatedPower.IsNull() {
		ethPortSettingsReq.AllocatedPower = openapi.PtrString(plan.AllocatedPower.ValueString())
	}
	if !plan.BspEnable.IsNull() {
		ethPortSettingsReq.BspEnable = openapi.PtrBool(plan.BspEnable.ValueBool())
	}
	if !plan.Broadcast.IsNull() {
		ethPortSettingsReq.Broadcast = openapi.PtrBool(plan.Broadcast.ValueBool())
	}
	if !plan.Multicast.IsNull() {
		ethPortSettingsReq.Multicast = openapi.PtrBool(plan.Multicast.ValueBool())
	}
	if !plan.MaxAllowedValue.IsNull() {
		ethPortSettingsReq.MaxAllowedValue = openapi.PtrInt32(int32(plan.MaxAllowedValue.ValueInt64()))
	}
	if !plan.MaxAllowedUnit.IsNull() {
		ethPortSettingsReq.MaxAllowedUnit = openapi.PtrString(plan.MaxAllowedUnit.ValueString())
	}
	if !plan.Action.IsNull() {
		ethPortSettingsReq.Action = openapi.PtrString(plan.Action.ValueString())
	}
	if !plan.Fec.IsNull() {
		ethPortSettingsReq.Fec = openapi.PtrString(plan.Fec.ValueString())
	}
	if !plan.SingleLink.IsNull() {
		ethPortSettingsReq.SingleLink = openapi.PtrBool(plan.SingleLink.ValueBool())
	}

	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties{}
		if !op.Group.IsNull() {
			objProps.Group = openapi.PtrString(op.Group.ValueString())
		} else {
			objProps.Group = nil
		}
		ethPortSettingsReq.ObjectProperties = &objProps
	} else {
		ethPortSettingsReq.ObjectProperties = nil
	}

	provCtx := r.provCtx
	bulkOpsMgr := provCtx.bulkOpsMgr
	operationID := bulkOpsMgr.AddEthPortSettingsPut(ctx, name, ethPortSettingsReq)

	provCtx.NotifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for Ethernet port settings creation operation %s to complete", operationID))
	if err := bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Create Ethernet Port Settings",
			fmt.Sprintf("Error creating Ethernet port settings %s: %v", name, err),
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Ethernet port settings %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "eth_port_settings")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
}

func (r *verityEthPortSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityEthPortSettingsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Authenticate",
			fmt.Sprintf("Error authenticating with API: %v", err),
		)
		return
	}

	provCtx := r.provCtx
	bulkOpsMgr := provCtx.bulkOpsMgr
	settingsName := state.Name.ValueString()

	if bulkOpsMgr != nil && bulkOpsMgr.HasPendingOrRecentEthPortSettingsOperations() {
		tflog.Info(ctx, fmt.Sprintf("Skipping eth port settings %s verification - trusting recent successful API operation", settingsName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("No recent eth port settings operations found, performing normal verification for %s", settingsName))

	type EthPortSettingsResponse struct {
		EthPortSettings map[string]map[string]interface{} `json:"eth_port_settings"`
	}

	var result EthPortSettingsResponse
	var err error
	maxRetries := 3

	for attempt := 0; attempt < maxRetries; attempt++ {
		settingsData, fetchErr := getCachedResponse(ctx, provCtx, "eth_port_settings", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch Ethernet port settings")
			apiResp, err := provCtx.client.EthPortSettingsAPI.EthportsettingsGet(ctx).Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading eth port settings: %v", err)
			}
			defer apiResp.Body.Close()

			var res EthPortSettingsResponse
			if err := json.NewDecoder(apiResp.Body).Decode(&res); err != nil {
				return nil, fmt.Errorf("error decoding eth port settings response: %v", err)
			}
			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d Ethernet port settings from API", len(res.EthPortSettings)))
			return res, nil
		})

		if fetchErr != nil {
			err = fetchErr
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch eth port settings on attempt %d, retrying in %v", attempt+1, sleepTime))
			time.Sleep(sleepTime)
			continue
		}
		result = settingsData.(EthPortSettingsResponse)
		break
	}

	if err != nil {
		resp.Diagnostics.AddError("Failed to Read Ethernet Port Settings", err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for Ethernet port settings with ID: %s", settingsName))
	var settings map[string]interface{}
	exists := false

	if data, ok := result.EthPortSettings[settingsName]; ok {
		settings = data
		exists = true
		tflog.Debug(ctx, fmt.Sprintf("Found Ethernet port settings directly by ID: %s", settingsName))
	} else {
		var nameStr types.String
		diags := req.State.GetAttribute(ctx, path.Root("name"), &nameStr)
		if !diags.HasError() && !nameStr.IsNull() {
			settingsNameFromAttr := nameStr.ValueString()
			tflog.Debug(ctx, fmt.Sprintf("Settings not found by ID, trying name attribute: %s", settingsNameFromAttr))
			if data, ok := result.EthPortSettings[settingsNameFromAttr]; ok {
				settings = data
				settingsName = settingsNameFromAttr
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found Ethernet port settings by name attribute: %s", settingsNameFromAttr))
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Settings not found directly by ID '%s', searching through all settings", settingsName))
		for apiName, s := range result.EthPortSettings {
			if name, ok := s["name"].(string); ok && name == settingsName {
				settings = s
				settingsName = apiName
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found Ethernet port settings with name '%s' under API key '%s'", name, apiName))
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Ethernet port settings with ID '%s' not found in API response", settingsName))
		resp.State.RemoveResource(ctx)
		return
	}

	if op, ok := settings["object_properties"].(map[string]interface{}); ok {
		state.ObjectProperties = []verityEthPortSettingsObjectPropertiesModel{
			{
				Group: types.StringValue(fmt.Sprintf("%v", op["group"])),
			},
		}
	} else {
		state.ObjectProperties = nil
	}

	if val, ok := settings["name"]; ok {
		state.Name = types.StringValue(fmt.Sprintf("%v", val))
	}
	if val, ok := settings["enable"].(bool); ok {
		state.Enable = types.BoolValue(val)
	} else {
		state.Enable = types.BoolNull()
	}
	if val, ok := settings["auto_negotiation"].(bool); ok {
		state.AutoNegotiation = types.BoolValue(val)
	} else {
		state.AutoNegotiation = types.BoolNull()
	}
	if val, ok := settings["max_bit_rate"].(string); ok {
		state.MaxBitRate = types.StringValue(val)
	} else {
		state.MaxBitRate = types.StringNull()
	}
	if val, ok := settings["duplex_mode"].(string); ok {
		state.DuplexMode = types.StringValue(val)
	} else {
		state.DuplexMode = types.StringNull()
	}
	if val, ok := settings["stp_enable"].(bool); ok {
		state.StpEnable = types.BoolValue(val)
	} else {
		state.StpEnable = types.BoolNull()
	}
	if val, ok := settings["fast_learning_mode"].(bool); ok {
		state.FastLearningMode = types.BoolValue(val)
	} else {
		state.FastLearningMode = types.BoolNull()
	}
	if val, ok := settings["bpdu_guard"].(bool); ok {
		state.BpduGuard = types.BoolValue(val)
	} else {
		state.BpduGuard = types.BoolNull()
	}
	if val, ok := settings["bpdu_filter"].(bool); ok {
		state.BpduFilter = types.BoolValue(val)
	} else {
		state.BpduFilter = types.BoolNull()
	}
	if val, ok := settings["guard_loop"].(bool); ok {
		state.GuardLoop = types.BoolValue(val)
	} else {
		state.GuardLoop = types.BoolNull()
	}
	if val, ok := settings["poe_enable"].(bool); ok {
		state.PoeEnable = types.BoolValue(val)
	} else {
		state.PoeEnable = types.BoolNull()
	}
	if val, ok := settings["priority"].(string); ok {
		state.Priority = types.StringValue(val)
	} else {
		state.Priority = types.StringNull()
	}
	if val, ok := settings["allocated_power"].(string); ok {
		state.AllocatedPower = types.StringValue(val)
	} else {
		state.AllocatedPower = types.StringNull()
	}
	if val, ok := settings["bsp_enable"].(bool); ok {
		state.BspEnable = types.BoolValue(val)
	} else {
		state.BspEnable = types.BoolNull()
	}
	if val, ok := settings["broadcast"].(bool); ok {
		state.Broadcast = types.BoolValue(val)
	} else {
		state.Broadcast = types.BoolNull()
	}
	if val, ok := settings["multicast"].(bool); ok {
		state.Multicast = types.BoolValue(val)
	} else {
		state.Multicast = types.BoolNull()
	}
	if val, ok := settings["max_allowed_value"]; ok {
		switch v := val.(type) {
		case float64:
			state.MaxAllowedValue = types.Int64Value(int64(v))
		case int:
			state.MaxAllowedValue = types.Int64Value(int64(v))
		default:
			state.MaxAllowedValue = types.Int64Null()
		}
	} else {
		state.MaxAllowedValue = types.Int64Null()
	}
	if val, ok := settings["max_allowed_unit"].(string); ok {
		state.MaxAllowedUnit = types.StringValue(val)
	} else {
		state.MaxAllowedUnit = types.StringNull()
	}
	if val, ok := settings["action"].(string); ok {
		state.Action = types.StringValue(val)
	} else {
		state.Action = types.StringNull()
	}
	if val, ok := settings["fec"].(string); ok {
		state.Fec = types.StringValue(val)
	} else {
		state.Fec = types.StringNull()
	}
	if val, ok := settings["single_link"].(bool); ok {
		state.SingleLink = types.BoolValue(val)
	} else {
		state.SingleLink = types.BoolNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityEthPortSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan verityEthPortSettingsResourceModel
	var state verityEthPortSettingsResourceModel

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
			fmt.Sprintf("Error authenticating with API: %v", err),
		)
		return
	}

	ethPortSettingsReq := openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName{}
	hasChanges := false

	objectPropertiesChanged := false

	if (len(plan.ObjectProperties) > 0 && len(state.ObjectProperties) == 0) ||
		(len(plan.ObjectProperties) == 0 && len(state.ObjectProperties) > 0) {
		objectPropertiesChanged = true
	} else if len(plan.ObjectProperties) > 0 && len(state.ObjectProperties) > 0 {
		if !plan.ObjectProperties[0].Group.Equal(state.ObjectProperties[0].Group) {
			objectPropertiesChanged = true
		}
	}

	if objectPropertiesChanged {
		if len(plan.ObjectProperties) > 0 {
			objProps := openapi.ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties{}
			if !plan.ObjectProperties[0].Group.IsNull() {
				objProps.Group = openapi.PtrString(plan.ObjectProperties[0].Group.ValueString())
			} else {
				objProps.Group = nil
			}
			ethPortSettingsReq.ObjectProperties = &objProps
		} else {
			ethPortSettingsReq.ObjectProperties = nil
		}
		hasChanges = true
	}

	boolFields := []struct {
		planValue  types.Bool
		stateValue types.Bool
		setter     func(val bool)
	}{
		{
			planValue:  plan.Enable,
			stateValue: state.Enable,
			setter: func(val bool) {
				ethPortSettingsReq.Enable = openapi.PtrBool(val)
			},
		},
		{
			planValue:  plan.AutoNegotiation,
			stateValue: state.AutoNegotiation,
			setter: func(val bool) {
				ethPortSettingsReq.AutoNegotiation = openapi.PtrBool(val)
			},
		},
		{
			planValue:  plan.StpEnable,
			stateValue: state.StpEnable,
			setter: func(val bool) {
				ethPortSettingsReq.StpEnable = openapi.PtrBool(val)
			},
		},
		{
			planValue:  plan.FastLearningMode,
			stateValue: state.FastLearningMode,
			setter: func(val bool) {
				ethPortSettingsReq.FastLearningMode = openapi.PtrBool(val)
			},
		},
		{
			planValue:  plan.BpduGuard,
			stateValue: state.BpduGuard,
			setter: func(val bool) {
				ethPortSettingsReq.BpduGuard = openapi.PtrBool(val)
			},
		},
		{
			planValue:  plan.BpduFilter,
			stateValue: state.BpduFilter,
			setter: func(val bool) {
				ethPortSettingsReq.BpduFilter = openapi.PtrBool(val)
			},
		},
		{
			planValue:  plan.GuardLoop,
			stateValue: state.GuardLoop,
			setter: func(val bool) {
				ethPortSettingsReq.GuardLoop = openapi.PtrBool(val)
			},
		},
		{
			planValue:  plan.PoeEnable,
			stateValue: state.PoeEnable,
			setter: func(val bool) {
				ethPortSettingsReq.PoeEnable = openapi.PtrBool(val)
			},
		},
		{
			planValue:  plan.BspEnable,
			stateValue: state.BspEnable,
			setter: func(val bool) {
				ethPortSettingsReq.BspEnable = openapi.PtrBool(val)
			},
		},
		{
			planValue:  plan.Broadcast,
			stateValue: state.Broadcast,
			setter: func(val bool) {
				ethPortSettingsReq.Broadcast = openapi.PtrBool(val)
			},
		},
		{
			planValue:  plan.Multicast,
			stateValue: state.Multicast,
			setter: func(val bool) {
				ethPortSettingsReq.Multicast = openapi.PtrBool(val)
			},
		},
		{
			planValue:  plan.SingleLink,
			stateValue: state.SingleLink,
			setter: func(val bool) {
				ethPortSettingsReq.SingleLink = openapi.PtrBool(val)
			},
		},
	}

	for _, field := range boolFields {
		if !field.planValue.Equal(field.stateValue) {
			field.setter(field.planValue.ValueBool())
			hasChanges = true
		}
	}

	stringFields := []struct {
		planValue  types.String
		stateValue types.String
		setter     func(val string)
	}{
		{
			planValue:  plan.MaxBitRate,
			stateValue: state.MaxBitRate,
			setter: func(val string) {
				ethPortSettingsReq.MaxBitRate = openapi.PtrString(val)
			},
		},
		{
			planValue:  plan.DuplexMode,
			stateValue: state.DuplexMode,
			setter: func(val string) {
				ethPortSettingsReq.DuplexMode = openapi.PtrString(val)
			},
		},
		{
			planValue:  plan.Priority,
			stateValue: state.Priority,
			setter: func(val string) {
				ethPortSettingsReq.Priority = openapi.PtrString(val)
			},
		},
		{
			planValue:  plan.AllocatedPower,
			stateValue: state.AllocatedPower,
			setter: func(val string) {
				ethPortSettingsReq.AllocatedPower = openapi.PtrString(val)
			},
		},
		{
			planValue:  plan.MaxAllowedUnit,
			stateValue: state.MaxAllowedUnit,
			setter: func(val string) {
				ethPortSettingsReq.MaxAllowedUnit = openapi.PtrString(val)
			},
		},
		{
			planValue:  plan.Action,
			stateValue: state.Action,
			setter: func(val string) {
				ethPortSettingsReq.Action = openapi.PtrString(val)
			},
		},
		{
			planValue:  plan.Fec,
			stateValue: state.Fec,
			setter: func(val string) {
				ethPortSettingsReq.Fec = openapi.PtrString(val)
			},
		},
	}

	for _, field := range stringFields {
		if !field.planValue.Equal(field.stateValue) {
			field.setter(field.planValue.ValueString())
			hasChanges = true
		}
	}

	intFields := []struct {
		planValue  types.Int64
		stateValue types.Int64
		setter     func(val int32)
	}{
		{
			planValue:  plan.MaxAllowedValue,
			stateValue: state.MaxAllowedValue,
			setter: func(val int32) {
				ethPortSettingsReq.MaxAllowedValue = openapi.PtrInt32(val)
			},
		},
	}

	for _, field := range intFields {
		if !field.planValue.Equal(field.stateValue) {
			field.setter(int32(field.planValue.ValueInt64()))
			hasChanges = true
		}
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	provCtx := r.provCtx
	bulkMgr := provCtx.bulkOpsMgr
	resourceName := plan.Name.ValueString()

	operationID := bulkMgr.AddEthPortSettingsPatch(ctx, resourceName, ethPortSettingsReq)
	provCtx.NotifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for Ethernet port settings update operation %s to complete", operationID))
	if err := bulkMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Update Ethernet Port Settings",
			fmt.Sprintf("Error updating Ethernet port settings %s: %v", resourceName, err),
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Ethernet port settings %s update operation completed successfully", resourceName))
	clearCache(ctx, provCtx, "eth_port_settings")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *verityEthPortSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityEthPortSettingsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Authenticate",
			fmt.Sprintf("Error authenticating with API: %v", err),
		)
		return
	}

	name := state.Name.ValueString()
	bulkMgr := r.provCtx.bulkOpsMgr
	operationID := bulkMgr.AddEthPortSettingsDelete(ctx, name)
	r.provCtx.NotifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for Ethernet port settings deletion operation %s to complete", operationID))
	if err := bulkMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Delete Ethernet Port Settings",
			fmt.Sprintf("Error deleting Ethernet port settings %s: %v", name, err),
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Ethernet port settings %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "eth_port_settings")
	resp.State.RemoveResource(ctx)
}

func (r *verityEthPortSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
