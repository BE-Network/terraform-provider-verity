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
	_ resource.Resource                = &verityVoicePortProfileResource{}
	_ resource.ResourceWithConfigure   = &verityVoicePortProfileResource{}
	_ resource.ResourceWithImportState = &verityVoicePortProfileResource{}
)

func NewVerityVoicePortProfileResource() resource.Resource {
	return &verityVoicePortProfileResource{}
}

type verityVoicePortProfileResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type verityVoicePortProfileResourceModel struct {
	Name                           types.String                                  `tfsdk:"name"`
	Enable                         types.Bool                                    `tfsdk:"enable"`
	Protocol                       types.String                                  `tfsdk:"protocol"`
	DigitMap                       types.String                                  `tfsdk:"digit_map"`
	CallThreeWayEnable             types.Bool                                    `tfsdk:"call_three_way_enable"`
	CallerIdEnable                 types.Bool                                    `tfsdk:"caller_id_enable"`
	CallerIdNameEnable             types.Bool                                    `tfsdk:"caller_id_name_enable"`
	CallWaitingEnable              types.Bool                                    `tfsdk:"call_waiting_enable"`
	CallForwardUnconditionalEnable types.Bool                                    `tfsdk:"call_forward_unconditional_enable"`
	CallForwardOnBusyEnable        types.Bool                                    `tfsdk:"call_forward_on_busy_enable"`
	CallForwardOnNoAnswerRingCount types.Int64                                   `tfsdk:"call_forward_on_no_answer_ring_count"`
	CallTransferEnable             types.Bool                                    `tfsdk:"call_transfer_enable"`
	AudioMwiEnable                 types.Bool                                    `tfsdk:"audio_mwi_enable"`
	AnonymousCallBlockEnable       types.Bool                                    `tfsdk:"anonymous_call_block_enable"`
	DoNotDisturbEnable             types.Bool                                    `tfsdk:"do_not_disturb_enable"`
	CidBlockingEnable              types.Bool                                    `tfsdk:"cid_blocking_enable"`
	CidNumPresentationStatus       types.String                                  `tfsdk:"cid_num_presentation_status"`
	CidNamePresentationStatus      types.String                                  `tfsdk:"cid_name_presentation_status"`
	CallWaitingCallerIdEnable      types.Bool                                    `tfsdk:"call_waiting_caller_id_enable"`
	CallHoldEnable                 types.Bool                                    `tfsdk:"call_hold_enable"`
	VisualMwiEnable                types.Bool                                    `tfsdk:"visual_mwi_enable"`
	MwiRefreshTimer                types.Int64                                   `tfsdk:"mwi_refresh_timer"`
	HotlineEnable                  types.Bool                                    `tfsdk:"hotline_enable"`
	DialToneFeatureDelay           types.Int64                                   `tfsdk:"dial_tone_feature_delay"`
	IntercomEnable                 types.Bool                                    `tfsdk:"intercom_enable"`
	IntercomTransferEnable         types.Bool                                    `tfsdk:"intercom_transfer_enable"`
	TransmitGain                   types.Int64                                   `tfsdk:"transmit_gain"`
	ReceiveGain                    types.Int64                                   `tfsdk:"receive_gain"`
	EchoCancellationEnable         types.Bool                                    `tfsdk:"echo_cancellation_enable"`
	JitterTarget                   types.Int64                                   `tfsdk:"jitter_target"`
	JitterBufferMax                types.Int64                                   `tfsdk:"jitter_buffer_max"`
	SignalingCode                  types.String                                  `tfsdk:"signaling_code"`
	ReleaseTimer                   types.Int64                                   `tfsdk:"release_timer"`
	RohTimer                       types.Int64                                   `tfsdk:"roh_timer"`
	ObjectProperties               []verityVoicePortProfileObjectPropertiesModel `tfsdk:"object_properties"`
}

type verityVoicePortProfileObjectPropertiesModel struct {
	PortMonitoring types.String `tfsdk:"port_monitoring"`
	Group          types.String `tfsdk:"group"`
	FormatDialPlan types.Bool   `tfsdk:"format_dial_plan"`
}

func (r *verityVoicePortProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_voice_port_profile"
}

func (r *verityVoicePortProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityVoicePortProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity Voice Port Profile",
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
			"protocol": schema.StringAttribute{
				Description: "Voice Protocol: MGCP or SIP",
				Optional:    true,
			},
			"digit_map": schema.StringAttribute{
				Description: "Dial Plan",
				Optional:    true,
			},
			"call_three_way_enable": schema.BoolAttribute{
				Description: "Enable three way calling",
				Optional:    true,
			},
			"caller_id_enable": schema.BoolAttribute{
				Description: "Caller ID",
				Optional:    true,
			},
			"caller_id_name_enable": schema.BoolAttribute{
				Description: "Caller ID Name",
				Optional:    true,
			},
			"call_waiting_enable": schema.BoolAttribute{
				Description: "Call Waiting",
				Optional:    true,
			},
			"call_forward_unconditional_enable": schema.BoolAttribute{
				Description: "Call Forward Unconditional",
				Optional:    true,
			},
			"call_forward_on_busy_enable": schema.BoolAttribute{
				Description: "Call Forward On Busy",
				Optional:    true,
			},
			"call_forward_on_no_answer_ring_count": schema.Int64Attribute{
				Description: "Call Forward on number of rings",
				Optional:    true,
			},
			"call_transfer_enable": schema.BoolAttribute{
				Description: "Call Transfer",
				Optional:    true,
			},
			"audio_mwi_enable": schema.BoolAttribute{
				Description: "Audio Message Waiting Indicator",
				Optional:    true,
			},
			"anonymous_call_block_enable": schema.BoolAttribute{
				Description: "Block all anonymous calls",
				Optional:    true,
			},
			"do_not_disturb_enable": schema.BoolAttribute{
				Description: "Do not disturb",
				Optional:    true,
			},
			"cid_blocking_enable": schema.BoolAttribute{
				Description: "CID Blocking",
				Optional:    true,
			},
			"cid_num_presentation_status": schema.StringAttribute{
				Description: "CID Number Presentation",
				Optional:    true,
			},
			"cid_name_presentation_status": schema.StringAttribute{
				Description: "CID Name Presentation",
				Optional:    true,
			},
			"call_waiting_caller_id_enable": schema.BoolAttribute{
				Description: "Call Waiting Caller ID",
				Optional:    true,
			},
			"call_hold_enable": schema.BoolAttribute{
				Description: "Call Hold",
				Optional:    true,
			},
			"visual_mwi_enable": schema.BoolAttribute{
				Description: "Visual Message Waiting Indicator",
				Optional:    true,
			},
			"mwi_refresh_timer": schema.Int64Attribute{
				Description: "Message Waiting Indicator Refresh",
				Optional:    true,
			},
			"hotline_enable": schema.BoolAttribute{
				Description: "Direct Connect",
				Optional:    true,
			},
			"dial_tone_feature_delay": schema.Int64Attribute{
				Description: "Dial Tone Feature Delay",
				Optional:    true,
			},
			"intercom_enable": schema.BoolAttribute{
				Description: "Intercom",
				Optional:    true,
			},
			"intercom_transfer_enable": schema.BoolAttribute{
				Description: "Intercom Transfer",
				Optional:    true,
			},
			"transmit_gain": schema.Int64Attribute{
				Description: "Transmit Gain in tenths of a dB. Example -30 would equal -3.0db",
				Optional:    true,
			},
			"receive_gain": schema.Int64Attribute{
				Description: "Receive Gain in tenths of a dB. Example -30 would equal -3.0db",
				Optional:    true,
			},
			"echo_cancellation_enable": schema.BoolAttribute{
				Description: "Echo Cancellation Enable",
				Optional:    true,
			},
			"jitter_target": schema.Int64Attribute{
				Description: "The target value of the jitter buffer in milliseconds",
				Optional:    true,
			},
			"jitter_buffer_max": schema.Int64Attribute{
				Description: "The maximum depth of the jitter buffer in milliseconds",
				Optional:    true,
			},
			"signaling_code": schema.StringAttribute{
				Description: "Signaling Code",
				Optional:    true,
			},
			"release_timer": schema.Int64Attribute{
				Description: "Release timer defined in seconds. The default value of this attribute is 10 seconds",
				Optional:    true,
			},
			"roh_timer": schema.Int64Attribute{
				Description: "Time in seconds for the receiver is off-hook before ROH tone is applied. The value 0 disables ROH timing. The default value is 15 seconds",
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"object_properties": schema.ListNestedBlock{
				Description: "Object properties for the voice port profile",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"port_monitoring": schema.StringAttribute{
							Description: "Defines importance of Link Down on this port",
							Optional:    true,
						},
						"group": schema.StringAttribute{
							Description: "Group",
							Optional:    true,
						},
						"format_dial_plan": schema.BoolAttribute{
							Description: "Format dial plan for easier viewing",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func (r *verityVoicePortProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityVoicePortProfileResourceModel
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
	vppProps := &openapi.VoiceportprofilesPutRequestVoicePortProfilesValue{
		Name: openapi.PtrString(name),
	}

	// Handle string fields
	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "Protocol", APIField: &vppProps.Protocol, TFValue: plan.Protocol},
		{FieldName: "DigitMap", APIField: &vppProps.DigitMap, TFValue: plan.DigitMap},
		{FieldName: "SignalingCode", APIField: &vppProps.SignalingCode, TFValue: plan.SignalingCode},
		{FieldName: "CidNumPresentationStatus", APIField: &vppProps.CidNumPresentationStatus, TFValue: plan.CidNumPresentationStatus},
		{FieldName: "CidNamePresentationStatus", APIField: &vppProps.CidNamePresentationStatus, TFValue: plan.CidNamePresentationStatus},
	})

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &vppProps.Enable, TFValue: plan.Enable},
		{FieldName: "CallThreeWayEnable", APIField: &vppProps.CallThreeWayEnable, TFValue: plan.CallThreeWayEnable},
		{FieldName: "CallerIdEnable", APIField: &vppProps.CallerIdEnable, TFValue: plan.CallerIdEnable},
		{FieldName: "CallerIdNameEnable", APIField: &vppProps.CallerIdNameEnable, TFValue: plan.CallerIdNameEnable},
		{FieldName: "CallWaitingEnable", APIField: &vppProps.CallWaitingEnable, TFValue: plan.CallWaitingEnable},
		{FieldName: "CallForwardUnconditionalEnable", APIField: &vppProps.CallForwardUnconditionalEnable, TFValue: plan.CallForwardUnconditionalEnable},
		{FieldName: "CallForwardOnBusyEnable", APIField: &vppProps.CallForwardOnBusyEnable, TFValue: plan.CallForwardOnBusyEnable},
		{FieldName: "CallTransferEnable", APIField: &vppProps.CallTransferEnable, TFValue: plan.CallTransferEnable},
		{FieldName: "AudioMwiEnable", APIField: &vppProps.AudioMwiEnable, TFValue: plan.AudioMwiEnable},
		{FieldName: "AnonymousCallBlockEnable", APIField: &vppProps.AnonymousCallBlockEnable, TFValue: plan.AnonymousCallBlockEnable},
		{FieldName: "DoNotDisturbEnable", APIField: &vppProps.DoNotDisturbEnable, TFValue: plan.DoNotDisturbEnable},
		{FieldName: "CidBlockingEnable", APIField: &vppProps.CidBlockingEnable, TFValue: plan.CidBlockingEnable},
		{FieldName: "CallWaitingCallerIdEnable", APIField: &vppProps.CallWaitingCallerIdEnable, TFValue: plan.CallWaitingCallerIdEnable},
		{FieldName: "CallHoldEnable", APIField: &vppProps.CallHoldEnable, TFValue: plan.CallHoldEnable},
		{FieldName: "VisualMwiEnable", APIField: &vppProps.VisualMwiEnable, TFValue: plan.VisualMwiEnable},
		{FieldName: "HotlineEnable", APIField: &vppProps.HotlineEnable, TFValue: plan.HotlineEnable},
		{FieldName: "IntercomEnable", APIField: &vppProps.IntercomEnable, TFValue: plan.IntercomEnable},
		{FieldName: "IntercomTransferEnable", APIField: &vppProps.IntercomTransferEnable, TFValue: plan.IntercomTransferEnable},
		{FieldName: "EchoCancellationEnable", APIField: &vppProps.EchoCancellationEnable, TFValue: plan.EchoCancellationEnable},
	})

	// Handle nullable int64 fields
	utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
		{FieldName: "CallForwardOnNoAnswerRingCount", APIField: &vppProps.CallForwardOnNoAnswerRingCount, TFValue: plan.CallForwardOnNoAnswerRingCount},
		{FieldName: "MwiRefreshTimer", APIField: &vppProps.MwiRefreshTimer, TFValue: plan.MwiRefreshTimer},
		{FieldName: "DialToneFeatureDelay", APIField: &vppProps.DialToneFeatureDelay, TFValue: plan.DialToneFeatureDelay},
		{FieldName: "TransmitGain", APIField: &vppProps.TransmitGain, TFValue: plan.TransmitGain},
		{FieldName: "ReceiveGain", APIField: &vppProps.ReceiveGain, TFValue: plan.ReceiveGain},
		{FieldName: "JitterTarget", APIField: &vppProps.JitterTarget, TFValue: plan.JitterTarget},
		{FieldName: "JitterBufferMax", APIField: &vppProps.JitterBufferMax, TFValue: plan.JitterBufferMax},
		{FieldName: "ReleaseTimer", APIField: &vppProps.ReleaseTimer, TFValue: plan.ReleaseTimer},
		{FieldName: "RohTimer", APIField: &vppProps.RohTimer, TFValue: plan.RohTimer},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.VoiceportprofilesPutRequestVoicePortProfilesValueObjectProperties{}
		utils.SetObjectPropertiesFields([]utils.ObjectPropertiesField{
			{Name: "PortMonitoring", TFValue: op.PortMonitoring, APIValue: &objProps.PortMonitoring},
			{Name: "Group", TFValue: op.Group, APIValue: &objProps.Group},
			{Name: "FormatDialPlan", TFValue: op.FormatDialPlan, APIValue: &objProps.FormatDialPlan},
		})
		vppProps.ObjectProperties = &objProps
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "voice_port_profile", name, *vppProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Voice Port Profile %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "voice_port_profiles")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
}

func (r *verityVoicePortProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityVoicePortProfileResourceModel
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

	vppName := state.Name.ValueString()

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("voice_port_profile") {
		tflog.Info(ctx, fmt.Sprintf("Skipping Voice Port Profile %s verification – trusting recent successful API operation", vppName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching Voice Port Profiles for verification of %s", vppName))

	type VoicePortProfileResponse struct {
		VoicePortProfiles map[string]interface{} `json:"voice_port_profiles"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "voice_port_profiles", vppName,
		func() (VoicePortProfileResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch Voice Port Profiles")
			respAPI, err := r.client.VoicePortProfilesAPI.VoiceportprofilesGet(ctx).Execute()
			if err != nil {
				return VoicePortProfileResponse{}, fmt.Errorf("error reading Voice Port Profiles: %v", err)
			}
			defer respAPI.Body.Close()

			var res VoicePortProfileResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return VoicePortProfileResponse{}, fmt.Errorf("failed to decode Voice Port Profiles response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d Voice Port Profiles", len(res.VoicePortProfiles)))
			return res, nil
		},
		getCachedResponse,
	)
	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Voice Port Profile %s", vppName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for Voice Port Profile with name: %s", vppName))

	vppData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.VoicePortProfiles,
		vppName,
		func(data interface{}) (string, bool) {
			if voicePortProfile, ok := data.(map[string]interface{}); ok {
				if name, ok := voicePortProfile["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Voice Port Profile with name '%s' not found in API response", vppName))
		resp.State.RemoveResource(ctx)
		return
	}

	vppMap, ok := vppData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Voice Port Profile Data",
			fmt.Sprintf("Voice Port Profile data is not in expected format for %s", vppName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found Voice Port Profile '%s' under API key '%s'", vppName, actualAPIName))

	state.Name = utils.MapStringFromAPI(vppMap["name"])

	// Handle object properties
	if objProps, ok := vppMap["object_properties"].(map[string]interface{}); ok {
		state.ObjectProperties = []verityVoicePortProfileObjectPropertiesModel{
			{
				PortMonitoring: utils.MapStringFromAPI(objProps["port_monitoring"]),
				Group:          utils.MapStringFromAPI(objProps["group"]),
				FormatDialPlan: utils.MapBoolFromAPI(objProps["format_dial_plan"]),
			},
		}
	} else {
		state.ObjectProperties = nil
	}

	// Map string fields
	stringFieldMappings := map[string]*types.String{
		"protocol":                     &state.Protocol,
		"digit_map":                    &state.DigitMap,
		"signaling_code":               &state.SignalingCode,
		"cid_num_presentation_status":  &state.CidNumPresentationStatus,
		"cid_name_presentation_status": &state.CidNamePresentationStatus,
	}

	for apiKey, stateField := range stringFieldMappings {
		*stateField = utils.MapStringFromAPI(vppMap[apiKey])
	}

	// Map boolean fields
	boolFieldMappings := map[string]*types.Bool{
		"enable":                            &state.Enable,
		"call_three_way_enable":             &state.CallThreeWayEnable,
		"caller_id_enable":                  &state.CallerIdEnable,
		"caller_id_name_enable":             &state.CallerIdNameEnable,
		"call_waiting_enable":               &state.CallWaitingEnable,
		"call_forward_unconditional_enable": &state.CallForwardUnconditionalEnable,
		"call_forward_on_busy_enable":       &state.CallForwardOnBusyEnable,
		"call_transfer_enable":              &state.CallTransferEnable,
		"audio_mwi_enable":                  &state.AudioMwiEnable,
		"anonymous_call_block_enable":       &state.AnonymousCallBlockEnable,
		"do_not_disturb_enable":             &state.DoNotDisturbEnable,
		"cid_blocking_enable":               &state.CidBlockingEnable,
		"call_waiting_caller_id_enable":     &state.CallWaitingCallerIdEnable,
		"call_hold_enable":                  &state.CallHoldEnable,
		"visual_mwi_enable":                 &state.VisualMwiEnable,
		"hotline_enable":                    &state.HotlineEnable,
		"intercom_enable":                   &state.IntercomEnable,
		"intercom_transfer_enable":          &state.IntercomTransferEnable,
		"echo_cancellation_enable":          &state.EchoCancellationEnable,
	}

	for apiKey, stateField := range boolFieldMappings {
		*stateField = utils.MapBoolFromAPI(vppMap[apiKey])
	}

	// Map int64 fields
	int64FieldMappings := map[string]*types.Int64{
		"call_forward_on_no_answer_ring_count": &state.CallForwardOnNoAnswerRingCount,
		"mwi_refresh_timer":                    &state.MwiRefreshTimer,
		"dial_tone_feature_delay":              &state.DialToneFeatureDelay,
		"transmit_gain":                        &state.TransmitGain,
		"receive_gain":                         &state.ReceiveGain,
		"jitter_target":                        &state.JitterTarget,
		"jitter_buffer_max":                    &state.JitterBufferMax,
		"release_timer":                        &state.ReleaseTimer,
		"roh_timer":                            &state.RohTimer,
	}

	for apiKey, stateField := range int64FieldMappings {
		*stateField = utils.MapNullableInt64FromAPI(vppMap[apiKey])
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityVoicePortProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityVoicePortProfileResourceModel

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
	vppProps := openapi.VoiceportprofilesPutRequestVoicePortProfilesValue{}
	hasChanges := false

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(val *string) { vppProps.Name = val }, &hasChanges)
	utils.CompareAndSetStringField(plan.Protocol, state.Protocol, func(val *string) { vppProps.Protocol = val }, &hasChanges)
	utils.CompareAndSetStringField(plan.DigitMap, state.DigitMap, func(val *string) { vppProps.DigitMap = val }, &hasChanges)
	utils.CompareAndSetStringField(plan.SignalingCode, state.SignalingCode, func(val *string) { vppProps.SignalingCode = val }, &hasChanges)
	utils.CompareAndSetStringField(plan.CidNumPresentationStatus, state.CidNumPresentationStatus, func(val *string) { vppProps.CidNumPresentationStatus = val }, &hasChanges)
	utils.CompareAndSetStringField(plan.CidNamePresentationStatus, state.CidNamePresentationStatus, func(val *string) { vppProps.CidNamePresentationStatus = val }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(val *bool) { vppProps.Enable = val }, &hasChanges)
	utils.CompareAndSetBoolField(plan.CallThreeWayEnable, state.CallThreeWayEnable, func(val *bool) { vppProps.CallThreeWayEnable = val }, &hasChanges)
	utils.CompareAndSetBoolField(plan.CallerIdEnable, state.CallerIdEnable, func(val *bool) { vppProps.CallerIdEnable = val }, &hasChanges)
	utils.CompareAndSetBoolField(plan.CallerIdNameEnable, state.CallerIdNameEnable, func(val *bool) { vppProps.CallerIdNameEnable = val }, &hasChanges)
	utils.CompareAndSetBoolField(plan.CallWaitingEnable, state.CallWaitingEnable, func(val *bool) { vppProps.CallWaitingEnable = val }, &hasChanges)
	utils.CompareAndSetBoolField(plan.CallForwardUnconditionalEnable, state.CallForwardUnconditionalEnable, func(val *bool) { vppProps.CallForwardUnconditionalEnable = val }, &hasChanges)
	utils.CompareAndSetBoolField(plan.CallForwardOnBusyEnable, state.CallForwardOnBusyEnable, func(val *bool) { vppProps.CallForwardOnBusyEnable = val }, &hasChanges)
	utils.CompareAndSetBoolField(plan.CallTransferEnable, state.CallTransferEnable, func(val *bool) { vppProps.CallTransferEnable = val }, &hasChanges)
	utils.CompareAndSetBoolField(plan.AudioMwiEnable, state.AudioMwiEnable, func(val *bool) { vppProps.AudioMwiEnable = val }, &hasChanges)
	utils.CompareAndSetBoolField(plan.AnonymousCallBlockEnable, state.AnonymousCallBlockEnable, func(val *bool) { vppProps.AnonymousCallBlockEnable = val }, &hasChanges)
	utils.CompareAndSetBoolField(plan.DoNotDisturbEnable, state.DoNotDisturbEnable, func(val *bool) { vppProps.DoNotDisturbEnable = val }, &hasChanges)
	utils.CompareAndSetBoolField(plan.CidBlockingEnable, state.CidBlockingEnable, func(val *bool) { vppProps.CidBlockingEnable = val }, &hasChanges)
	utils.CompareAndSetBoolField(plan.CallWaitingCallerIdEnable, state.CallWaitingCallerIdEnable, func(val *bool) { vppProps.CallWaitingCallerIdEnable = val }, &hasChanges)
	utils.CompareAndSetBoolField(plan.CallHoldEnable, state.CallHoldEnable, func(val *bool) { vppProps.CallHoldEnable = val }, &hasChanges)
	utils.CompareAndSetBoolField(plan.VisualMwiEnable, state.VisualMwiEnable, func(val *bool) { vppProps.VisualMwiEnable = val }, &hasChanges)
	utils.CompareAndSetBoolField(plan.HotlineEnable, state.HotlineEnable, func(val *bool) { vppProps.HotlineEnable = val }, &hasChanges)
	utils.CompareAndSetBoolField(plan.IntercomEnable, state.IntercomEnable, func(val *bool) { vppProps.IntercomEnable = val }, &hasChanges)
	utils.CompareAndSetBoolField(plan.IntercomTransferEnable, state.IntercomTransferEnable, func(val *bool) { vppProps.IntercomTransferEnable = val }, &hasChanges)
	utils.CompareAndSetBoolField(plan.EchoCancellationEnable, state.EchoCancellationEnable, func(val *bool) { vppProps.EchoCancellationEnable = val }, &hasChanges)

	// Handle nullable int64 field changes
	utils.CompareAndSetNullableInt64Field(plan.CallForwardOnNoAnswerRingCount, state.CallForwardOnNoAnswerRingCount, func(val *openapi.NullableInt32) { vppProps.CallForwardOnNoAnswerRingCount = *val }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.MwiRefreshTimer, state.MwiRefreshTimer, func(val *openapi.NullableInt32) { vppProps.MwiRefreshTimer = *val }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.DialToneFeatureDelay, state.DialToneFeatureDelay, func(val *openapi.NullableInt32) { vppProps.DialToneFeatureDelay = *val }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.TransmitGain, state.TransmitGain, func(val *openapi.NullableInt32) { vppProps.TransmitGain = *val }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.ReceiveGain, state.ReceiveGain, func(val *openapi.NullableInt32) { vppProps.ReceiveGain = *val }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.JitterTarget, state.JitterTarget, func(val *openapi.NullableInt32) { vppProps.JitterTarget = *val }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.JitterBufferMax, state.JitterBufferMax, func(val *openapi.NullableInt32) { vppProps.JitterBufferMax = *val }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.ReleaseTimer, state.ReleaseTimer, func(val *openapi.NullableInt32) { vppProps.ReleaseTimer = *val }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.RohTimer, state.RohTimer, func(val *openapi.NullableInt32) { vppProps.RohTimer = *val }, &hasChanges)

	// Handle object properties
	if len(plan.ObjectProperties) > 0 && len(state.ObjectProperties) > 0 {
		objProps := openapi.VoiceportprofilesPutRequestVoicePortProfilesValueObjectProperties{}
		op := plan.ObjectProperties[0]
		st := state.ObjectProperties[0]
		objPropsChanged := false

		utils.CompareAndSetObjectPropertiesFields([]utils.ObjectPropertiesFieldWithComparison{
			{Name: "PortMonitoring", PlanValue: op.PortMonitoring, StateValue: st.PortMonitoring, APIValue: &objProps.PortMonitoring},
			{Name: "Group", PlanValue: op.Group, StateValue: st.Group, APIValue: &objProps.Group},
			{Name: "FormatDialPlan", PlanValue: op.FormatDialPlan, StateValue: st.FormatDialPlan, APIValue: &objProps.FormatDialPlan},
		}, &objPropsChanged)

		if objPropsChanged {
			vppProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "voice_port_profile", name, vppProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Voice Port Profile %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "voice_port_profiles")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *verityVoicePortProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityVoicePortProfileResourceModel
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "voice_port_profile", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Voice Port Profile %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "voice_port_profiles")
	resp.State.RemoveResource(ctx)
}

func (r *verityVoicePortProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
