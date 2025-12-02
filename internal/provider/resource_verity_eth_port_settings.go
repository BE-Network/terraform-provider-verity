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
	bulkOpsMgr           *bulkops.Manager
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
			fmt.Sprintf("Error authenticating with API: %s", err),
		)
		return
	}

	name := plan.Name.ValueString()
	ethPortSettingsProps := &openapi.EthportsettingsPutRequestEthPortSettingsValue{
		Name: openapi.PtrString(name),
	}

	// Handle string fields
	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "MaxBitRate", APIField: &ethPortSettingsProps.MaxBitRate, TFValue: plan.MaxBitRate},
		{FieldName: "DuplexMode", APIField: &ethPortSettingsProps.DuplexMode, TFValue: plan.DuplexMode},
		{FieldName: "Priority", APIField: &ethPortSettingsProps.Priority, TFValue: plan.Priority},
		{FieldName: "AllocatedPower", APIField: &ethPortSettingsProps.AllocatedPower, TFValue: plan.AllocatedPower},
		{FieldName: "MaxAllowedUnit", APIField: &ethPortSettingsProps.MaxAllowedUnit, TFValue: plan.MaxAllowedUnit},
		{FieldName: "Action", APIField: &ethPortSettingsProps.Action, TFValue: plan.Action},
		{FieldName: "Fec", APIField: &ethPortSettingsProps.Fec, TFValue: plan.Fec},
	})

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &ethPortSettingsProps.Enable, TFValue: plan.Enable},
		{FieldName: "AutoNegotiation", APIField: &ethPortSettingsProps.AutoNegotiation, TFValue: plan.AutoNegotiation},
		{FieldName: "StpEnable", APIField: &ethPortSettingsProps.StpEnable, TFValue: plan.StpEnable},
		{FieldName: "FastLearningMode", APIField: &ethPortSettingsProps.FastLearningMode, TFValue: plan.FastLearningMode},
		{FieldName: "BpduGuard", APIField: &ethPortSettingsProps.BpduGuard, TFValue: plan.BpduGuard},
		{FieldName: "BpduFilter", APIField: &ethPortSettingsProps.BpduFilter, TFValue: plan.BpduFilter},
		{FieldName: "GuardLoop", APIField: &ethPortSettingsProps.GuardLoop, TFValue: plan.GuardLoop},
		{FieldName: "PoeEnable", APIField: &ethPortSettingsProps.PoeEnable, TFValue: plan.PoeEnable},
		{FieldName: "BspEnable", APIField: &ethPortSettingsProps.BspEnable, TFValue: plan.BspEnable},
		{FieldName: "Broadcast", APIField: &ethPortSettingsProps.Broadcast, TFValue: plan.Broadcast},
		{FieldName: "Multicast", APIField: &ethPortSettingsProps.Multicast, TFValue: plan.Multicast},
		{FieldName: "SingleLink", APIField: &ethPortSettingsProps.SingleLink, TFValue: plan.SingleLink},
	})

	// Handle nullable int64 fields
	utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
		{FieldName: "MaxAllowedValue", APIField: &ethPortSettingsProps.MaxAllowedValue, TFValue: plan.MaxAllowedValue},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.EthportsettingsPutRequestEthPortSettingsValueObjectProperties{}
		utils.SetObjectPropertiesFields([]utils.ObjectPropertiesField{
			{Name: "Group", TFValue: op.Group, APIValue: &objProps.Group},
		})
		ethPortSettingsProps.ObjectProperties = &objProps
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "eth_port_settings", name, *ethPortSettingsProps, &resp.Diagnostics)
	if !success {
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
			fmt.Sprintf("Error authenticating with API: %s", err),
		)
		return
	}

	ethPortSettingsName := state.Name.ValueString()

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("eth_port_settings") {
		tflog.Info(ctx, fmt.Sprintf("Skipping eth port settings %s verification – trusting recent successful API operation", ethPortSettingsName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching eth port settings for verification of %s", ethPortSettingsName))

	type EthPortSettingsResponse struct {
		EthPortSettings map[string]interface{} `json:"eth_port_settings"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "eth_port_settings", ethPortSettingsName,
		func() (EthPortSettingsResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch eth port settings")
			respAPI, err := r.client.EthPortSettingsAPI.EthportsettingsGet(ctx).Execute()
			if err != nil {
				return EthPortSettingsResponse{}, fmt.Errorf("error reading eth port settings: %v", err)
			}
			defer respAPI.Body.Close()

			var res EthPortSettingsResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return EthPortSettingsResponse{}, fmt.Errorf("failed to decode eth port settings response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d eth port settings", len(res.EthPortSettings)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Eth Port Settings %s", ethPortSettingsName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for eth port settings with name: %s", ethPortSettingsName))

	ethPortSettingsData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.EthPortSettings,
		ethPortSettingsName,
		func(data interface{}) (string, bool) {
			if ethPortSettings, ok := data.(map[string]interface{}); ok {
				if name, ok := ethPortSettings["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Eth Port Settings with name '%s' not found in API response", ethPortSettingsName))
		resp.State.RemoveResource(ctx)
		return
	}

	ethPortSettingsMap, ok := ethPortSettingsData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Eth Port Settings Data",
			fmt.Sprintf("Eth Port Settings data is not in expected format for %s", ethPortSettingsName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found eth port settings '%s' under API key '%s'", ethPortSettingsName, actualAPIName))

	state.Name = utils.MapStringFromAPI(ethPortSettingsMap["name"])

	// Handle object properties
	if objProps, ok := ethPortSettingsMap["object_properties"].(map[string]interface{}); ok {
		state.ObjectProperties = []verityEthPortSettingsObjectPropertiesModel{
			{
				Group: utils.MapStringFromAPI(objProps["group"]),
			},
		}
	} else {
		state.ObjectProperties = nil
	}

	// Map string fields
	stringFieldMappings := map[string]*types.String{
		"max_bit_rate":     &state.MaxBitRate,
		"duplex_mode":      &state.DuplexMode,
		"priority":         &state.Priority,
		"allocated_power":  &state.AllocatedPower,
		"max_allowed_unit": &state.MaxAllowedUnit,
		"action":           &state.Action,
		"fec":              &state.Fec,
	}

	for apiKey, stateField := range stringFieldMappings {
		*stateField = utils.MapStringFromAPI(ethPortSettingsMap[apiKey])
	}

	// Map boolean fields
	boolFieldMappings := map[string]*types.Bool{
		"enable":             &state.Enable,
		"auto_negotiation":   &state.AutoNegotiation,
		"stp_enable":         &state.StpEnable,
		"fast_learning_mode": &state.FastLearningMode,
		"bpdu_guard":         &state.BpduGuard,
		"bpdu_filter":        &state.BpduFilter,
		"guard_loop":         &state.GuardLoop,
		"poe_enable":         &state.PoeEnable,
		"bsp_enable":         &state.BspEnable,
		"broadcast":          &state.Broadcast,
		"multicast":          &state.Multicast,
		"single_link":        &state.SingleLink,
	}

	for apiKey, stateField := range boolFieldMappings {
		*stateField = utils.MapBoolFromAPI(ethPortSettingsMap[apiKey])
	}

	// Map int64 fields
	int64FieldMappings := map[string]*types.Int64{
		"max_allowed_value": &state.MaxAllowedValue,
	}

	for apiKey, stateField := range int64FieldMappings {
		*stateField = utils.MapInt64FromAPI(ethPortSettingsMap[apiKey])
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

	name := plan.Name.ValueString()
	ethPortSettingsProps := openapi.EthportsettingsPutRequestEthPortSettingsValue{}
	hasChanges := false

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { ethPortSettingsProps.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.MaxBitRate, state.MaxBitRate, func(v *string) { ethPortSettingsProps.MaxBitRate = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.DuplexMode, state.DuplexMode, func(v *string) { ethPortSettingsProps.DuplexMode = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Priority, state.Priority, func(v *string) { ethPortSettingsProps.Priority = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.AllocatedPower, state.AllocatedPower, func(v *string) { ethPortSettingsProps.AllocatedPower = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.MaxAllowedUnit, state.MaxAllowedUnit, func(v *string) { ethPortSettingsProps.MaxAllowedUnit = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Action, state.Action, func(v *string) { ethPortSettingsProps.Action = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Fec, state.Fec, func(v *string) { ethPortSettingsProps.Fec = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { ethPortSettingsProps.Enable = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.AutoNegotiation, state.AutoNegotiation, func(v *bool) { ethPortSettingsProps.AutoNegotiation = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.StpEnable, state.StpEnable, func(v *bool) { ethPortSettingsProps.StpEnable = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.FastLearningMode, state.FastLearningMode, func(v *bool) { ethPortSettingsProps.FastLearningMode = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.BpduGuard, state.BpduGuard, func(v *bool) { ethPortSettingsProps.BpduGuard = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.BpduFilter, state.BpduFilter, func(v *bool) { ethPortSettingsProps.BpduFilter = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.GuardLoop, state.GuardLoop, func(v *bool) { ethPortSettingsProps.GuardLoop = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.PoeEnable, state.PoeEnable, func(v *bool) { ethPortSettingsProps.PoeEnable = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.BspEnable, state.BspEnable, func(v *bool) { ethPortSettingsProps.BspEnable = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.Broadcast, state.Broadcast, func(v *bool) { ethPortSettingsProps.Broadcast = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.Multicast, state.Multicast, func(v *bool) { ethPortSettingsProps.Multicast = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.SingleLink, state.SingleLink, func(v *bool) { ethPortSettingsProps.SingleLink = v }, &hasChanges)

	// Handle nullable int64 field changes
	utils.CompareAndSetNullableInt64Field(plan.MaxAllowedValue, state.MaxAllowedValue, func(v *openapi.NullableInt32) { ethPortSettingsProps.MaxAllowedValue = *v }, &hasChanges)

	// Handle object properties
	if len(plan.ObjectProperties) > 0 && len(state.ObjectProperties) > 0 {
		objProps := openapi.EthportsettingsPutRequestEthPortSettingsValueObjectProperties{}
		op := plan.ObjectProperties[0]
		st := state.ObjectProperties[0]
		objPropsChanged := false

		utils.CompareAndSetObjectPropertiesFields([]utils.ObjectPropertiesFieldWithComparison{
			{Name: "Group", PlanValue: op.Group, StateValue: st.Group, APIValue: &objProps.Group},
		}, &objPropsChanged)

		if objPropsChanged {
			ethPortSettingsProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "eth_port_settings", name, ethPortSettingsProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Eth Port Settings %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "eth_port_settings")
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
			fmt.Sprintf("Error authenticating with API: %s", err),
		)
		return
	}

	name := state.Name.ValueString()

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "eth_port_settings", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Eth Port Settings %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "eth_port_settings")
	resp.State.RemoveResource(ctx)
}

func (r *verityEthPortSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
