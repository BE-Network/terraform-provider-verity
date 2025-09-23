package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
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
	bulkOpsMgr           *utils.BulkOperationManager
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
	IsDefault      types.Bool   `tfsdk:"isdefault"`
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
						"isdefault": schema.BoolAttribute{
							Description: "Default object.",
							Optional:    true,
						},
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

	if !plan.Enable.IsNull() {
		vppProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
	}
	if !plan.Protocol.IsNull() {
		vppProps.Protocol = openapi.PtrString(plan.Protocol.ValueString())
	}
	if !plan.DigitMap.IsNull() {
		vppProps.DigitMap = openapi.PtrString(plan.DigitMap.ValueString())
	}
	if !plan.SignalingCode.IsNull() {
		vppProps.SignalingCode = openapi.PtrString(plan.SignalingCode.ValueString())
	}
	if !plan.CidNumPresentationStatus.IsNull() {
		vppProps.CidNumPresentationStatus = openapi.PtrString(plan.CidNumPresentationStatus.ValueString())
	}
	if !plan.CidNamePresentationStatus.IsNull() {
		vppProps.CidNamePresentationStatus = openapi.PtrString(plan.CidNamePresentationStatus.ValueString())
	}

	if !plan.CallThreeWayEnable.IsNull() {
		vppProps.CallThreeWayEnable = openapi.PtrBool(plan.CallThreeWayEnable.ValueBool())
	}
	if !plan.CallerIdEnable.IsNull() {
		vppProps.CallerIdEnable = openapi.PtrBool(plan.CallerIdEnable.ValueBool())
	}
	if !plan.CallerIdNameEnable.IsNull() {
		vppProps.CallerIdNameEnable = openapi.PtrBool(plan.CallerIdNameEnable.ValueBool())
	}
	if !plan.CallWaitingEnable.IsNull() {
		vppProps.CallWaitingEnable = openapi.PtrBool(plan.CallWaitingEnable.ValueBool())
	}
	if !plan.CallForwardUnconditionalEnable.IsNull() {
		vppProps.CallForwardUnconditionalEnable = openapi.PtrBool(plan.CallForwardUnconditionalEnable.ValueBool())
	}
	if !plan.CallForwardOnBusyEnable.IsNull() {
		vppProps.CallForwardOnBusyEnable = openapi.PtrBool(plan.CallForwardOnBusyEnable.ValueBool())
	}
	if !plan.CallTransferEnable.IsNull() {
		vppProps.CallTransferEnable = openapi.PtrBool(plan.CallTransferEnable.ValueBool())
	}
	if !plan.AudioMwiEnable.IsNull() {
		vppProps.AudioMwiEnable = openapi.PtrBool(plan.AudioMwiEnable.ValueBool())
	}
	if !plan.AnonymousCallBlockEnable.IsNull() {
		vppProps.AnonymousCallBlockEnable = openapi.PtrBool(plan.AnonymousCallBlockEnable.ValueBool())
	}
	if !plan.DoNotDisturbEnable.IsNull() {
		vppProps.DoNotDisturbEnable = openapi.PtrBool(plan.DoNotDisturbEnable.ValueBool())
	}
	if !plan.CidBlockingEnable.IsNull() {
		vppProps.CidBlockingEnable = openapi.PtrBool(plan.CidBlockingEnable.ValueBool())
	}
	if !plan.CallWaitingCallerIdEnable.IsNull() {
		vppProps.CallWaitingCallerIdEnable = openapi.PtrBool(plan.CallWaitingCallerIdEnable.ValueBool())
	}
	if !plan.CallHoldEnable.IsNull() {
		vppProps.CallHoldEnable = openapi.PtrBool(plan.CallHoldEnable.ValueBool())
	}
	if !plan.VisualMwiEnable.IsNull() {
		vppProps.VisualMwiEnable = openapi.PtrBool(plan.VisualMwiEnable.ValueBool())
	}
	if !plan.HotlineEnable.IsNull() {
		vppProps.HotlineEnable = openapi.PtrBool(plan.HotlineEnable.ValueBool())
	}
	if !plan.IntercomEnable.IsNull() {
		vppProps.IntercomEnable = openapi.PtrBool(plan.IntercomEnable.ValueBool())
	}
	if !plan.IntercomTransferEnable.IsNull() {
		vppProps.IntercomTransferEnable = openapi.PtrBool(plan.IntercomTransferEnable.ValueBool())
	}
	if !plan.EchoCancellationEnable.IsNull() {
		vppProps.EchoCancellationEnable = openapi.PtrBool(plan.EchoCancellationEnable.ValueBool())
	}

	if !plan.CallForwardOnNoAnswerRingCount.IsNull() {
		val := int32(plan.CallForwardOnNoAnswerRingCount.ValueInt64())
		vppProps.CallForwardOnNoAnswerRingCount = *openapi.NewNullableInt32(&val)
	} else {
		vppProps.CallForwardOnNoAnswerRingCount = *openapi.NewNullableInt32(nil)
	}
	if !plan.MwiRefreshTimer.IsNull() {
		val := int32(plan.MwiRefreshTimer.ValueInt64())
		vppProps.MwiRefreshTimer = *openapi.NewNullableInt32(&val)
	} else {
		vppProps.MwiRefreshTimer = *openapi.NewNullableInt32(nil)
	}
	if !plan.DialToneFeatureDelay.IsNull() {
		val := int32(plan.DialToneFeatureDelay.ValueInt64())
		vppProps.DialToneFeatureDelay = *openapi.NewNullableInt32(&val)
	} else {
		vppProps.DialToneFeatureDelay = *openapi.NewNullableInt32(nil)
	}
	if !plan.TransmitGain.IsNull() {
		val := int32(plan.TransmitGain.ValueInt64())
		vppProps.TransmitGain = *openapi.NewNullableInt32(&val)
	} else {
		vppProps.TransmitGain = *openapi.NewNullableInt32(nil)
	}
	if !plan.ReceiveGain.IsNull() {
		val := int32(plan.ReceiveGain.ValueInt64())
		vppProps.ReceiveGain = *openapi.NewNullableInt32(&val)
	} else {
		vppProps.ReceiveGain = *openapi.NewNullableInt32(nil)
	}
	if !plan.JitterTarget.IsNull() {
		val := int32(plan.JitterTarget.ValueInt64())
		vppProps.JitterTarget = *openapi.NewNullableInt32(&val)
	} else {
		vppProps.JitterTarget = *openapi.NewNullableInt32(nil)
	}
	if !plan.JitterBufferMax.IsNull() {
		val := int32(plan.JitterBufferMax.ValueInt64())
		vppProps.JitterBufferMax = *openapi.NewNullableInt32(&val)
	} else {
		vppProps.JitterBufferMax = *openapi.NewNullableInt32(nil)
	}
	if !plan.ReleaseTimer.IsNull() {
		val := int32(plan.ReleaseTimer.ValueInt64())
		vppProps.ReleaseTimer = *openapi.NewNullableInt32(&val)
	} else {
		vppProps.ReleaseTimer = *openapi.NewNullableInt32(nil)
	}
	if !plan.RohTimer.IsNull() {
		val := int32(plan.RohTimer.ValueInt64())
		vppProps.RohTimer = *openapi.NewNullableInt32(&val)
	} else {
		vppProps.RohTimer = *openapi.NewNullableInt32(nil)
	}

	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.VoiceportprofilesPutRequestVoicePortProfilesValueObjectProperties{}
		if !op.IsDefault.IsNull() {
			objProps.Isdefault = openapi.PtrBool(op.IsDefault.ValueBool())
		}
		if !op.PortMonitoring.IsNull() {
			objProps.PortMonitoring = openapi.PtrString(op.PortMonitoring.ValueString())
		}
		if !op.Group.IsNull() {
			objProps.Group = openapi.PtrString(op.Group.ValueString())
		}
		if !op.FormatDialPlan.IsNull() {
			objProps.FormatDialPlan = openapi.PtrBool(op.FormatDialPlan.ValueBool())
		}
		vppProps.ObjectProperties = &objProps
	}

	operationID := r.bulkOpsMgr.AddPut(ctx, "voice_port_profile", name, *vppProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for voice port profile creation operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Create Voice Port Profile %s", name))...,
		)
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

	var result VoicePortProfileResponse
	var err error
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		vppData, fetchErr := getCachedResponse(ctx, r.provCtx, "voice_port_profiles", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch Voice Port Profiles")
			respAPI, err := r.client.VoicePortProfilesAPI.VoiceportprofilesGet(ctx).Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading Voice Port Profiles: %v", err)
			}
			defer respAPI.Body.Close()

			var res VoicePortProfileResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return nil, fmt.Errorf("failed to decode Voice Port Profiles response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d Voice Port Profiles", len(res.VoicePortProfiles)))
			return res, nil
		})
		if fetchErr != nil {
			err = fetchErr
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch Voice Port Profiles on attempt %d, retrying in %v", attempt+1, sleepTime))
			time.Sleep(sleepTime)
			continue
		}
		result = vppData.(VoicePortProfileResponse)
		break
	}
	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Voice Port Profile %s", vppName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for Voice Port Profile with ID: %s", vppName))
	var vppData map[string]interface{}
	exists := false

	if data, ok := result.VoicePortProfiles[vppName].(map[string]interface{}); ok {
		vppData = data
		exists = true
		tflog.Debug(ctx, fmt.Sprintf("Found Voice Port Profile directly by ID: %s", vppName))
	} else {
		for apiName, v := range result.VoicePortProfiles {
			voicePortProfile, ok := v.(map[string]interface{})
			if !ok {
				continue
			}

			if name, ok := voicePortProfile["name"].(string); ok && name == vppName {
				vppData = voicePortProfile
				vppName = apiName
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found Voice Port Profile with name '%s' under API key '%s'", name, apiName))
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Voice Port Profile with ID '%s' not found in API response", vppName))
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(fmt.Sprintf("%v", vppData["name"]))

	boolFields := map[string]*types.Bool{
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

	for apiKey, stateField := range boolFields {
		if value, ok := vppData[apiKey].(bool); ok {
			*stateField = types.BoolValue(value)
		} else {
			*stateField = types.BoolNull()
		}
	}

	stringFields := map[string]*types.String{
		"protocol":                     &state.Protocol,
		"digit_map":                    &state.DigitMap,
		"cid_num_presentation_status":  &state.CidNumPresentationStatus,
		"cid_name_presentation_status": &state.CidNamePresentationStatus,
		"signaling_code":               &state.SignalingCode,
	}

	for apiKey, stateField := range stringFields {
		if value, ok := vppData[apiKey].(string); ok {
			*stateField = types.StringValue(value)
		} else {
			*stateField = types.StringNull()
		}
	}

	intFields := map[string]*types.Int64{
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

	for apiKey, stateField := range intFields {
		if value, ok := vppData[apiKey]; ok && value != nil {
			switch v := value.(type) {
			case int:
				*stateField = types.Int64Value(int64(v))
			case int32:
				*stateField = types.Int64Value(int64(v))
			case int64:
				*stateField = types.Int64Value(v)
			case float64:
				*stateField = types.Int64Value(int64(v))
			case string:
				if intVal, err := strconv.ParseInt(v, 10, 64); err == nil {
					*stateField = types.Int64Value(intVal)
				} else {
					*stateField = types.Int64Null()
				}
			default:
				*stateField = types.Int64Null()
			}
		} else {
			*stateField = types.Int64Null()
		}
	}

	if objProps, ok := vppData["object_properties"].(map[string]interface{}); ok {
		op := verityVoicePortProfileObjectPropertiesModel{}
		if isDefault, ok := objProps["isdefault"].(bool); ok {
			op.IsDefault = types.BoolValue(isDefault)
		} else {
			op.IsDefault = types.BoolNull()
		}
		if portMonitoring, ok := objProps["port_monitoring"].(string); ok {
			op.PortMonitoring = types.StringValue(portMonitoring)
		} else {
			op.PortMonitoring = types.StringNull()
		}
		if group, ok := objProps["group"].(string); ok {
			op.Group = types.StringValue(group)
		} else {
			op.Group = types.StringNull()
		}
		if formatDialPlan, ok := objProps["format_dial_plan"].(bool); ok {
			op.FormatDialPlan = types.BoolValue(formatDialPlan)
		} else {
			op.FormatDialPlan = types.BoolNull()
		}
		state.ObjectProperties = []verityVoicePortProfileObjectPropertiesModel{op}
	} else {
		state.ObjectProperties = nil
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

	if !plan.Name.Equal(state.Name) {
		vppProps.Name = openapi.PtrString(name)
		hasChanges = true
	}

	if !plan.Enable.Equal(state.Enable) {
		vppProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
		hasChanges = true
	}
	if !plan.Protocol.Equal(state.Protocol) {
		vppProps.Protocol = openapi.PtrString(plan.Protocol.ValueString())
		hasChanges = true
	}
	if !plan.DigitMap.Equal(state.DigitMap) {
		vppProps.DigitMap = openapi.PtrString(plan.DigitMap.ValueString())
		hasChanges = true
	}
	if !plan.SignalingCode.Equal(state.SignalingCode) {
		vppProps.SignalingCode = openapi.PtrString(plan.SignalingCode.ValueString())
		hasChanges = true
	}
	if !plan.CidNumPresentationStatus.Equal(state.CidNumPresentationStatus) {
		vppProps.CidNumPresentationStatus = openapi.PtrString(plan.CidNumPresentationStatus.ValueString())
		hasChanges = true
	}
	if !plan.CidNamePresentationStatus.Equal(state.CidNamePresentationStatus) {
		vppProps.CidNamePresentationStatus = openapi.PtrString(plan.CidNamePresentationStatus.ValueString())
		hasChanges = true
	}

	boolFields := map[string]struct {
		planField  types.Bool
		stateField types.Bool
		setter     func(bool)
	}{
		"call_three_way_enable":             {plan.CallThreeWayEnable, state.CallThreeWayEnable, func(v bool) { vppProps.CallThreeWayEnable = openapi.PtrBool(v) }},
		"caller_id_enable":                  {plan.CallerIdEnable, state.CallerIdEnable, func(v bool) { vppProps.CallerIdEnable = openapi.PtrBool(v) }},
		"caller_id_name_enable":             {plan.CallerIdNameEnable, state.CallerIdNameEnable, func(v bool) { vppProps.CallerIdNameEnable = openapi.PtrBool(v) }},
		"call_waiting_enable":               {plan.CallWaitingEnable, state.CallWaitingEnable, func(v bool) { vppProps.CallWaitingEnable = openapi.PtrBool(v) }},
		"call_forward_unconditional_enable": {plan.CallForwardUnconditionalEnable, state.CallForwardUnconditionalEnable, func(v bool) { vppProps.CallForwardUnconditionalEnable = openapi.PtrBool(v) }},
		"call_forward_on_busy_enable":       {plan.CallForwardOnBusyEnable, state.CallForwardOnBusyEnable, func(v bool) { vppProps.CallForwardOnBusyEnable = openapi.PtrBool(v) }},
		"call_transfer_enable":              {plan.CallTransferEnable, state.CallTransferEnable, func(v bool) { vppProps.CallTransferEnable = openapi.PtrBool(v) }},
		"audio_mwi_enable":                  {plan.AudioMwiEnable, state.AudioMwiEnable, func(v bool) { vppProps.AudioMwiEnable = openapi.PtrBool(v) }},
		"anonymous_call_block_enable":       {plan.AnonymousCallBlockEnable, state.AnonymousCallBlockEnable, func(v bool) { vppProps.AnonymousCallBlockEnable = openapi.PtrBool(v) }},
		"do_not_disturb_enable":             {plan.DoNotDisturbEnable, state.DoNotDisturbEnable, func(v bool) { vppProps.DoNotDisturbEnable = openapi.PtrBool(v) }},
		"cid_blocking_enable":               {plan.CidBlockingEnable, state.CidBlockingEnable, func(v bool) { vppProps.CidBlockingEnable = openapi.PtrBool(v) }},
		"call_waiting_caller_id_enable":     {plan.CallWaitingCallerIdEnable, state.CallWaitingCallerIdEnable, func(v bool) { vppProps.CallWaitingCallerIdEnable = openapi.PtrBool(v) }},
		"call_hold_enable":                  {plan.CallHoldEnable, state.CallHoldEnable, func(v bool) { vppProps.CallHoldEnable = openapi.PtrBool(v) }},
		"visual_mwi_enable":                 {plan.VisualMwiEnable, state.VisualMwiEnable, func(v bool) { vppProps.VisualMwiEnable = openapi.PtrBool(v) }},
		"hotline_enable":                    {plan.HotlineEnable, state.HotlineEnable, func(v bool) { vppProps.HotlineEnable = openapi.PtrBool(v) }},
		"intercom_enable":                   {plan.IntercomEnable, state.IntercomEnable, func(v bool) { vppProps.IntercomEnable = openapi.PtrBool(v) }},
		"intercom_transfer_enable":          {plan.IntercomTransferEnable, state.IntercomTransferEnable, func(v bool) { vppProps.IntercomTransferEnable = openapi.PtrBool(v) }},
		"echo_cancellation_enable":          {plan.EchoCancellationEnable, state.EchoCancellationEnable, func(v bool) { vppProps.EchoCancellationEnable = openapi.PtrBool(v) }},
	}

	for _, field := range boolFields {
		if !field.planField.Equal(field.stateField) {
			field.setter(field.planField.ValueBool())
			hasChanges = true
		}
	}

	intFields := map[string]struct {
		planField  types.Int64
		stateField types.Int64
		setter     func(types.Int64)
	}{
		"call_forward_on_no_answer_ring_count": {plan.CallForwardOnNoAnswerRingCount, state.CallForwardOnNoAnswerRingCount, func(v types.Int64) {
			if !v.IsNull() {
				val := int32(v.ValueInt64())
				vppProps.CallForwardOnNoAnswerRingCount = *openapi.NewNullableInt32(&val)
			} else {
				vppProps.CallForwardOnNoAnswerRingCount = *openapi.NewNullableInt32(nil)
			}
		}},
		"mwi_refresh_timer": {plan.MwiRefreshTimer, state.MwiRefreshTimer, func(v types.Int64) {
			if !v.IsNull() {
				val := int32(v.ValueInt64())
				vppProps.MwiRefreshTimer = *openapi.NewNullableInt32(&val)
			} else {
				vppProps.MwiRefreshTimer = *openapi.NewNullableInt32(nil)
			}
		}},
		"dial_tone_feature_delay": {plan.DialToneFeatureDelay, state.DialToneFeatureDelay, func(v types.Int64) {
			if !v.IsNull() {
				val := int32(v.ValueInt64())
				vppProps.DialToneFeatureDelay = *openapi.NewNullableInt32(&val)
			} else {
				vppProps.DialToneFeatureDelay = *openapi.NewNullableInt32(nil)
			}
		}},
		"transmit_gain": {plan.TransmitGain, state.TransmitGain, func(v types.Int64) {
			if !v.IsNull() {
				val := int32(v.ValueInt64())
				vppProps.TransmitGain = *openapi.NewNullableInt32(&val)
			} else {
				vppProps.TransmitGain = *openapi.NewNullableInt32(nil)
			}
		}},
		"receive_gain": {plan.ReceiveGain, state.ReceiveGain, func(v types.Int64) {
			if !v.IsNull() {
				val := int32(v.ValueInt64())
				vppProps.ReceiveGain = *openapi.NewNullableInt32(&val)
			} else {
				vppProps.ReceiveGain = *openapi.NewNullableInt32(nil)
			}
		}},
		"jitter_target": {plan.JitterTarget, state.JitterTarget, func(v types.Int64) {
			if !v.IsNull() {
				val := int32(v.ValueInt64())
				vppProps.JitterTarget = *openapi.NewNullableInt32(&val)
			} else {
				vppProps.JitterTarget = *openapi.NewNullableInt32(nil)
			}
		}},
		"jitter_buffer_max": {plan.JitterBufferMax, state.JitterBufferMax, func(v types.Int64) {
			if !v.IsNull() {
				val := int32(v.ValueInt64())
				vppProps.JitterBufferMax = *openapi.NewNullableInt32(&val)
			} else {
				vppProps.JitterBufferMax = *openapi.NewNullableInt32(nil)
			}
		}},
		"release_timer": {plan.ReleaseTimer, state.ReleaseTimer, func(v types.Int64) {
			if !v.IsNull() {
				val := int32(v.ValueInt64())
				vppProps.ReleaseTimer = *openapi.NewNullableInt32(&val)
			} else {
				vppProps.ReleaseTimer = *openapi.NewNullableInt32(nil)
			}
		}},
		"roh_timer": {plan.RohTimer, state.RohTimer, func(v types.Int64) {
			if !v.IsNull() {
				val := int32(v.ValueInt64())
				vppProps.RohTimer = *openapi.NewNullableInt32(&val)
			} else {
				vppProps.RohTimer = *openapi.NewNullableInt32(nil)
			}
		}},
	}

	for _, field := range intFields {
		if !field.planField.Equal(field.stateField) {
			field.setter(field.planField)
			hasChanges = true
		}
	}

	if len(plan.ObjectProperties) > 0 {
		if len(state.ObjectProperties) == 0 ||
			!plan.ObjectProperties[0].IsDefault.Equal(state.ObjectProperties[0].IsDefault) ||
			!plan.ObjectProperties[0].PortMonitoring.Equal(state.ObjectProperties[0].PortMonitoring) ||
			!plan.ObjectProperties[0].Group.Equal(state.ObjectProperties[0].Group) ||
			!plan.ObjectProperties[0].FormatDialPlan.Equal(state.ObjectProperties[0].FormatDialPlan) {
			op := plan.ObjectProperties[0]
			objProps := openapi.VoiceportprofilesPutRequestVoicePortProfilesValueObjectProperties{}
			if !op.IsDefault.IsNull() {
				objProps.Isdefault = openapi.PtrBool(op.IsDefault.ValueBool())
			}
			if !op.PortMonitoring.IsNull() {
				objProps.PortMonitoring = openapi.PtrString(op.PortMonitoring.ValueString())
			}
			if !op.Group.IsNull() {
				objProps.Group = openapi.PtrString(op.Group.ValueString())
			}
			if !op.FormatDialPlan.IsNull() {
				objProps.FormatDialPlan = openapi.PtrBool(op.FormatDialPlan.ValueBool())
			}
			vppProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	operationID := r.bulkOpsMgr.AddPatch(ctx, "voice_port_profile", name, vppProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for Voice Port Profile update operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update Voice Port Profile %s", name))...,
		)
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
	operationID := r.bulkOpsMgr.AddDelete(ctx, "voice_port_profile", name)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for Voice Port Profile deletion operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Delete Voice Port Profile %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Voice Port Profile %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "voice_port_profiles")
	resp.State.RemoveResource(ctx)
}

func (r *verityVoicePortProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
