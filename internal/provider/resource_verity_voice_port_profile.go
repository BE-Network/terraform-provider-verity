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
	_ resource.ResourceWithModifyPlan  = &verityVoicePortProfileResource{}
)

const voicePortProfileResourceType = "voiceportprofiles"

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
				Computed:    true,
			},
			"protocol": schema.StringAttribute{
				Description: "Voice Protocol: MGCP or SIP",
				Optional:    true,
				Computed:    true,
			},
			"digit_map": schema.StringAttribute{
				Description: "Dial Plan",
				Optional:    true,
				Computed:    true,
			},
			"call_three_way_enable": schema.BoolAttribute{
				Description: "Enable three way calling",
				Optional:    true,
				Computed:    true,
			},
			"caller_id_enable": schema.BoolAttribute{
				Description: "Caller ID",
				Optional:    true,
				Computed:    true,
			},
			"caller_id_name_enable": schema.BoolAttribute{
				Description: "Caller ID Name",
				Optional:    true,
				Computed:    true,
			},
			"call_waiting_enable": schema.BoolAttribute{
				Description: "Call Waiting",
				Optional:    true,
				Computed:    true,
			},
			"call_forward_unconditional_enable": schema.BoolAttribute{
				Description: "Call Forward Unconditional",
				Optional:    true,
				Computed:    true,
			},
			"call_forward_on_busy_enable": schema.BoolAttribute{
				Description: "Call Forward On Busy",
				Optional:    true,
				Computed:    true,
			},
			"call_forward_on_no_answer_ring_count": schema.Int64Attribute{
				Description: "Call Forward on number of rings",
				Optional:    true,
				Computed:    true,
			},
			"call_transfer_enable": schema.BoolAttribute{
				Description: "Call Transfer",
				Optional:    true,
				Computed:    true,
			},
			"audio_mwi_enable": schema.BoolAttribute{
				Description: "Audio Message Waiting Indicator",
				Optional:    true,
				Computed:    true,
			},
			"anonymous_call_block_enable": schema.BoolAttribute{
				Description: "Block all anonymous calls",
				Optional:    true,
				Computed:    true,
			},
			"do_not_disturb_enable": schema.BoolAttribute{
				Description: "Do not disturb",
				Optional:    true,
				Computed:    true,
			},
			"cid_blocking_enable": schema.BoolAttribute{
				Description: "CID Blocking",
				Optional:    true,
				Computed:    true,
			},
			"cid_num_presentation_status": schema.StringAttribute{
				Description: "CID Number Presentation",
				Optional:    true,
				Computed:    true,
			},
			"cid_name_presentation_status": schema.StringAttribute{
				Description: "CID Name Presentation",
				Optional:    true,
				Computed:    true,
			},
			"call_waiting_caller_id_enable": schema.BoolAttribute{
				Description: "Call Waiting Caller ID",
				Optional:    true,
				Computed:    true,
			},
			"call_hold_enable": schema.BoolAttribute{
				Description: "Call Hold",
				Optional:    true,
				Computed:    true,
			},
			"visual_mwi_enable": schema.BoolAttribute{
				Description: "Visual Message Waiting Indicator",
				Optional:    true,
				Computed:    true,
			},
			"mwi_refresh_timer": schema.Int64Attribute{
				Description: "Message Waiting Indicator Refresh",
				Optional:    true,
				Computed:    true,
			},
			"hotline_enable": schema.BoolAttribute{
				Description: "Direct Connect",
				Optional:    true,
				Computed:    true,
			},
			"dial_tone_feature_delay": schema.Int64Attribute{
				Description: "Dial Tone Feature Delay",
				Optional:    true,
				Computed:    true,
			},
			"intercom_enable": schema.BoolAttribute{
				Description: "Intercom",
				Optional:    true,
				Computed:    true,
			},
			"intercom_transfer_enable": schema.BoolAttribute{
				Description: "Intercom Transfer",
				Optional:    true,
				Computed:    true,
			},
			"transmit_gain": schema.Int64Attribute{
				Description: "Transmit Gain in tenths of a dB. Example -30 would equal -3.0db",
				Optional:    true,
				Computed:    true,
			},
			"receive_gain": schema.Int64Attribute{
				Description: "Receive Gain in tenths of a dB. Example -30 would equal -3.0db",
				Optional:    true,
				Computed:    true,
			},
			"echo_cancellation_enable": schema.BoolAttribute{
				Description: "Echo Cancellation Enable",
				Optional:    true,
				Computed:    true,
			},
			"jitter_target": schema.Int64Attribute{
				Description: "The target value of the jitter buffer in milliseconds",
				Optional:    true,
				Computed:    true,
			},
			"jitter_buffer_max": schema.Int64Attribute{
				Description: "The maximum depth of the jitter buffer in milliseconds",
				Optional:    true,
				Computed:    true,
			},
			"signaling_code": schema.StringAttribute{
				Description: "Signaling Code",
				Optional:    true,
				Computed:    true,
			},
			"release_timer": schema.Int64Attribute{
				Description: "Release timer defined in seconds. The default value of this attribute is 10 seconds",
				Optional:    true,
				Computed:    true,
			},
			"roh_timer": schema.Int64Attribute{
				Description: "Time in seconds for the receiver is off-hook before ROH tone is applied. The value 0 disables ROH timing. The default value is 15 seconds",
				Optional:    true,
				Computed:    true,
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
							Computed:    true,
						},
						"group": schema.StringAttribute{
							Description: "Group",
							Optional:    true,
							Computed:    true,
						},
						"format_dial_plan": schema.BoolAttribute{
							Description: "Format dial plan for easier viewing",
							Optional:    true,
							Computed:    true,
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

	var config verityVoicePortProfileResourceModel
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

	// Handle nullable int64 fields - parse HCL to detect explicit config
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, "verity_voice_port_profile", name)

	utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
		{FieldName: "CallForwardOnNoAnswerRingCount", APIField: &vppProps.CallForwardOnNoAnswerRingCount, TFValue: config.CallForwardOnNoAnswerRingCount, IsConfigured: configuredAttrs.IsConfigured("call_forward_on_no_answer_ring_count")},
		{FieldName: "MwiRefreshTimer", APIField: &vppProps.MwiRefreshTimer, TFValue: config.MwiRefreshTimer, IsConfigured: configuredAttrs.IsConfigured("mwi_refresh_timer")},
		{FieldName: "DialToneFeatureDelay", APIField: &vppProps.DialToneFeatureDelay, TFValue: config.DialToneFeatureDelay, IsConfigured: configuredAttrs.IsConfigured("dial_tone_feature_delay")},
		{FieldName: "TransmitGain", APIField: &vppProps.TransmitGain, TFValue: config.TransmitGain, IsConfigured: configuredAttrs.IsConfigured("transmit_gain")},
		{FieldName: "ReceiveGain", APIField: &vppProps.ReceiveGain, TFValue: config.ReceiveGain, IsConfigured: configuredAttrs.IsConfigured("receive_gain")},
		{FieldName: "JitterTarget", APIField: &vppProps.JitterTarget, TFValue: config.JitterTarget, IsConfigured: configuredAttrs.IsConfigured("jitter_target")},
		{FieldName: "JitterBufferMax", APIField: &vppProps.JitterBufferMax, TFValue: config.JitterBufferMax, IsConfigured: configuredAttrs.IsConfigured("jitter_buffer_max")},
		{FieldName: "ReleaseTimer", APIField: &vppProps.ReleaseTimer, TFValue: config.ReleaseTimer, IsConfigured: configuredAttrs.IsConfigured("release_timer")},
		{FieldName: "RohTimer", APIField: &vppProps.RohTimer, TFValue: config.RohTimer, IsConfigured: configuredAttrs.IsConfigured("roh_timer")},
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

	var minState verityVoicePortProfileResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if voicePortProfileData, exists := bulkMgr.GetResourceResponse("voice_port_profile", name); exists {
			state := populateVoicePortProfileState(ctx, minState, voicePortProfileData, r.provCtx.mode)
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

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if voicePortProfileData, exists := r.bulkOpsMgr.GetResourceResponse("voice_port_profile", vppName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached voice port profile data for %s from recent operation", vppName))
			state = populateVoicePortProfileState(ctx, state, voicePortProfileData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

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

	state = populateVoicePortProfileState(ctx, state, vppMap, r.provCtx.mode)
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

	// Get config for nullable field handling
	var config verityVoicePortProfileResourceModel
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
	vppProps := openapi.VoiceportprofilesPutRequestVoicePortProfilesValue{}
	hasChanges := false

	// Parse HCL to detect which fields are explicitly configured
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, "verity_voice_port_profile", name)

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

	// Handle nullable int64 field changes - parse HCL to detect explicit config
	utils.CompareAndSetNullableInt64Field(config.CallForwardOnNoAnswerRingCount, state.CallForwardOnNoAnswerRingCount, configuredAttrs.IsConfigured("call_forward_on_no_answer_ring_count"), func(val *openapi.NullableInt32) { vppProps.CallForwardOnNoAnswerRingCount = *val }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.MwiRefreshTimer, state.MwiRefreshTimer, configuredAttrs.IsConfigured("mwi_refresh_timer"), func(val *openapi.NullableInt32) { vppProps.MwiRefreshTimer = *val }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.DialToneFeatureDelay, state.DialToneFeatureDelay, configuredAttrs.IsConfigured("dial_tone_feature_delay"), func(val *openapi.NullableInt32) { vppProps.DialToneFeatureDelay = *val }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.TransmitGain, state.TransmitGain, configuredAttrs.IsConfigured("transmit_gain"), func(val *openapi.NullableInt32) { vppProps.TransmitGain = *val }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.ReceiveGain, state.ReceiveGain, configuredAttrs.IsConfigured("receive_gain"), func(val *openapi.NullableInt32) { vppProps.ReceiveGain = *val }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.JitterTarget, state.JitterTarget, configuredAttrs.IsConfigured("jitter_target"), func(val *openapi.NullableInt32) { vppProps.JitterTarget = *val }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.JitterBufferMax, state.JitterBufferMax, configuredAttrs.IsConfigured("jitter_buffer_max"), func(val *openapi.NullableInt32) { vppProps.JitterBufferMax = *val }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.ReleaseTimer, state.ReleaseTimer, configuredAttrs.IsConfigured("release_timer"), func(val *openapi.NullableInt32) { vppProps.ReleaseTimer = *val }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.RohTimer, state.RohTimer, configuredAttrs.IsConfigured("roh_timer"), func(val *openapi.NullableInt32) { vppProps.RohTimer = *val }, &hasChanges)

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

	var minState verityVoicePortProfileResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Try to use cached response from bulk operation to populate state with API values
	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if voicePortProfileData, exists := bulkMgr.GetResourceResponse("voice_port_profile", name); exists {
			updatedState := populateVoicePortProfileState(ctx, minState, voicePortProfileData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
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

func populateVoicePortProfileState(ctx context.Context, state verityVoicePortProfileResourceModel, data map[string]interface{}, mode string) verityVoicePortProfileResourceModel {
	const resourceType = voicePortProfileResourceType

	state.Name = utils.MapStringFromAPI(data["name"])

	// String fields
	state.Protocol = utils.MapStringWithMode(data, "protocol", resourceType, mode)
	state.DigitMap = utils.MapStringWithMode(data, "digit_map", resourceType, mode)
	state.SignalingCode = utils.MapStringWithMode(data, "signaling_code", resourceType, mode)
	state.CidNumPresentationStatus = utils.MapStringWithMode(data, "cid_num_presentation_status", resourceType, mode)
	state.CidNamePresentationStatus = utils.MapStringWithMode(data, "cid_name_presentation_status", resourceType, mode)

	// Boolean fields
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)
	state.CallThreeWayEnable = utils.MapBoolWithMode(data, "call_three_way_enable", resourceType, mode)
	state.CallerIdEnable = utils.MapBoolWithMode(data, "caller_id_enable", resourceType, mode)
	state.CallerIdNameEnable = utils.MapBoolWithMode(data, "caller_id_name_enable", resourceType, mode)
	state.CallWaitingEnable = utils.MapBoolWithMode(data, "call_waiting_enable", resourceType, mode)
	state.CallForwardUnconditionalEnable = utils.MapBoolWithMode(data, "call_forward_unconditional_enable", resourceType, mode)
	state.CallForwardOnBusyEnable = utils.MapBoolWithMode(data, "call_forward_on_busy_enable", resourceType, mode)
	state.CallTransferEnable = utils.MapBoolWithMode(data, "call_transfer_enable", resourceType, mode)
	state.AudioMwiEnable = utils.MapBoolWithMode(data, "audio_mwi_enable", resourceType, mode)
	state.AnonymousCallBlockEnable = utils.MapBoolWithMode(data, "anonymous_call_block_enable", resourceType, mode)
	state.DoNotDisturbEnable = utils.MapBoolWithMode(data, "do_not_disturb_enable", resourceType, mode)
	state.CidBlockingEnable = utils.MapBoolWithMode(data, "cid_blocking_enable", resourceType, mode)
	state.CallWaitingCallerIdEnable = utils.MapBoolWithMode(data, "call_waiting_caller_id_enable", resourceType, mode)
	state.CallHoldEnable = utils.MapBoolWithMode(data, "call_hold_enable", resourceType, mode)
	state.VisualMwiEnable = utils.MapBoolWithMode(data, "visual_mwi_enable", resourceType, mode)
	state.HotlineEnable = utils.MapBoolWithMode(data, "hotline_enable", resourceType, mode)
	state.IntercomEnable = utils.MapBoolWithMode(data, "intercom_enable", resourceType, mode)
	state.IntercomTransferEnable = utils.MapBoolWithMode(data, "intercom_transfer_enable", resourceType, mode)
	state.EchoCancellationEnable = utils.MapBoolWithMode(data, "echo_cancellation_enable", resourceType, mode)

	// Int fields
	state.CallForwardOnNoAnswerRingCount = utils.MapInt64WithMode(data, "call_forward_on_no_answer_ring_count", resourceType, mode)
	state.MwiRefreshTimer = utils.MapInt64WithMode(data, "mwi_refresh_timer", resourceType, mode)
	state.DialToneFeatureDelay = utils.MapInt64WithMode(data, "dial_tone_feature_delay", resourceType, mode)
	state.TransmitGain = utils.MapInt64WithMode(data, "transmit_gain", resourceType, mode)
	state.ReceiveGain = utils.MapInt64WithMode(data, "receive_gain", resourceType, mode)
	state.JitterTarget = utils.MapInt64WithMode(data, "jitter_target", resourceType, mode)
	state.JitterBufferMax = utils.MapInt64WithMode(data, "jitter_buffer_max", resourceType, mode)
	state.ReleaseTimer = utils.MapInt64WithMode(data, "release_timer", resourceType, mode)
	state.RohTimer = utils.MapInt64WithMode(data, "roh_timer", resourceType, mode)

	// Handle object_properties block
	if utils.FieldAppliesToMode(resourceType, "object_properties", mode) {
		if objProps, ok := data["object_properties"].(map[string]interface{}); ok {
			objPropsModel := verityVoicePortProfileObjectPropertiesModel{
				PortMonitoring: utils.MapStringWithModeNested(objProps, "port_monitoring", resourceType, "object_properties.port_monitoring", mode),
				Group:          utils.MapStringWithModeNested(objProps, "group", resourceType, "object_properties.group", mode),
				FormatDialPlan: utils.MapBoolWithModeNested(objProps, "format_dial_plan", resourceType, "object_properties.format_dial_plan", mode),
			}
			state.ObjectProperties = []verityVoicePortProfileObjectPropertiesModel{objPropsModel}
		} else {
			state.ObjectProperties = nil
		}
	} else {
		state.ObjectProperties = nil
	}

	return state
}

func (r *verityVoicePortProfileResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// =========================================================================
	// Skip if deleting
	// =========================================================================
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityVoicePortProfileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Mode-aware field nullification
	// Set fields that don't apply to current mode to null to prevent
	// "known after apply" messages for irrelevant fields.
	// =========================================================================
	const resourceType = voicePortProfileResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyStrings(
		"protocol", "digit_map", "signaling_code",
		"cid_num_presentation_status", "cid_name_presentation_status",
	)

	nullifier.NullifyBools(
		"enable", "call_three_way_enable", "caller_id_enable", "caller_id_name_enable",
		"call_waiting_enable", "call_forward_unconditional_enable", "call_forward_on_busy_enable",
		"call_transfer_enable", "audio_mwi_enable", "anonymous_call_block_enable",
		"do_not_disturb_enable", "cid_blocking_enable", "call_waiting_caller_id_enable",
		"call_hold_enable", "visual_mwi_enable", "hotline_enable",
		"intercom_enable", "intercom_transfer_enable", "echo_cancellation_enable",
	)

	nullifier.NullifyInt64s(
		"call_forward_on_no_answer_ring_count", "mwi_refresh_timer", "dial_tone_feature_delay",
		"transmit_gain", "receive_gain", "jitter_target",
		"jitter_buffer_max", "release_timer", "roh_timer",
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
	var state verityVoicePortProfileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config verityVoicePortProfileResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Handle nullable Int64 fields (explicit null detection)
	// For Optional+Computed fields, Terraform copies state to plan when config
	// is null. We detect explicit null in HCL and force plan to null.
	// =========================================================================
	name := plan.Name.ValueString()
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, "verity_voice_port_profile", name)

	utils.HandleNullableFields(utils.NullableFieldsConfig{
		Ctx:             ctx,
		Plan:            &resp.Plan,
		ConfiguredAttrs: configuredAttrs,
		Int64Fields: []utils.NullableInt64Field{
			{AttrName: "call_forward_on_no_answer_ring_count", ConfigVal: config.CallForwardOnNoAnswerRingCount, StateVal: state.CallForwardOnNoAnswerRingCount},
			{AttrName: "mwi_refresh_timer", ConfigVal: config.MwiRefreshTimer, StateVal: state.MwiRefreshTimer},
			{AttrName: "dial_tone_feature_delay", ConfigVal: config.DialToneFeatureDelay, StateVal: state.DialToneFeatureDelay},
			{AttrName: "transmit_gain", ConfigVal: config.TransmitGain, StateVal: state.TransmitGain},
			{AttrName: "receive_gain", ConfigVal: config.ReceiveGain, StateVal: state.ReceiveGain},
			{AttrName: "jitter_target", ConfigVal: config.JitterTarget, StateVal: state.JitterTarget},
			{AttrName: "jitter_buffer_max", ConfigVal: config.JitterBufferMax, StateVal: state.JitterBufferMax},
			{AttrName: "release_timer", ConfigVal: config.ReleaseTimer, StateVal: state.ReleaseTimer},
			{AttrName: "roh_timer", ConfigVal: config.RohTimer, StateVal: state.RohTimer},
		},
	})
}
