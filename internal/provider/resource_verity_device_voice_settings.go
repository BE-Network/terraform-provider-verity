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
	_ resource.Resource                = &verityDeviceVoiceSettingsResource{}
	_ resource.ResourceWithConfigure   = &verityDeviceVoiceSettingsResource{}
	_ resource.ResourceWithImportState = &verityDeviceVoiceSettingsResource{}
	_ resource.ResourceWithModifyPlan  = &verityDeviceVoiceSettingsResource{}
)

const deviceVoiceSettingsResourceType = "devicevoicesettings"
const deviceVoiceSettingsTerraformType = "verity_device_voice_settings"

func NewVerityDeviceVoiceSettingsResource() resource.Resource {
	return &verityDeviceVoiceSettingsResource{}
}

type verityDeviceVoiceSettingsResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type verityDeviceVoiceSettingsResourceModel struct {
	Name                               types.String                                     `tfsdk:"name"`
	Enable                             types.Bool                                       `tfsdk:"enable"`
	DtmfMethod                         types.String                                     `tfsdk:"dtmf_method"`
	Region                             types.String                                     `tfsdk:"region"`
	Protocol                           types.String                                     `tfsdk:"protocol"`
	ProxyServer                        types.String                                     `tfsdk:"proxy_server"`
	ProxyServerPort                    types.Int64                                      `tfsdk:"proxy_server_port"`
	ProxyServerSecondary               types.String                                     `tfsdk:"proxy_server_secondary"`
	ProxyServerSecondaryPort           types.Int64                                      `tfsdk:"proxy_server_secondary_port"`
	RegistrarServer                    types.String                                     `tfsdk:"registrar_server"`
	RegistrarServerPort                types.Int64                                      `tfsdk:"registrar_server_port"`
	RegistrarServerSecondary           types.String                                     `tfsdk:"registrar_server_secondary"`
	RegistrarServerSecondaryPort       types.Int64                                      `tfsdk:"registrar_server_secondary_port"`
	UserAgentDomain                    types.String                                     `tfsdk:"user_agent_domain"`
	UserAgentTransport                 types.String                                     `tfsdk:"user_agent_transport"`
	UserAgentPort                      types.Int64                                      `tfsdk:"user_agent_port"`
	OutboundProxy                      types.String                                     `tfsdk:"outbound_proxy"`
	OutboundProxyPort                  types.Int64                                      `tfsdk:"outbound_proxy_port"`
	OutboundProxySecondary             types.String                                     `tfsdk:"outbound_proxy_secondary"`
	OutboundProxySecondaryPort         types.Int64                                      `tfsdk:"outbound_proxy_secondary_port"`
	RegistrationPeriod                 types.Int64                                      `tfsdk:"registration_period"`
	RegisterExpires                    types.Int64                                      `tfsdk:"register_expires"`
	VoicemailServer                    types.String                                     `tfsdk:"voicemail_server"`
	VoicemailServerPort                types.Int64                                      `tfsdk:"voicemail_server_port"`
	VoicemailServerExpires             types.Int64                                      `tfsdk:"voicemail_server_expires"`
	SipDscpMark                        types.Int64                                      `tfsdk:"sip_dscp_mark"`
	CallAgent1                         types.String                                     `tfsdk:"call_agent_1"`
	CallAgentPort1                     types.Int64                                      `tfsdk:"call_agent_port_1"`
	CallAgent2                         types.String                                     `tfsdk:"call_agent_2"`
	CallAgentPort2                     types.Int64                                      `tfsdk:"call_agent_port_2"`
	Domain                             types.String                                     `tfsdk:"domain"`
	MgcpDscpMark                       types.Int64                                      `tfsdk:"mgcp_dscp_mark"`
	TerminationBase                    types.String                                     `tfsdk:"termination_base"`
	LocalPortMin                       types.Int64                                      `tfsdk:"local_port_min"`
	LocalPortMax                       types.Int64                                      `tfsdk:"local_port_max"`
	EventPayloadType                   types.Int64                                      `tfsdk:"event_payload_type"`
	CasEvents                          types.Int64                                      `tfsdk:"cas_events"`
	DscpMark                           types.Int64                                      `tfsdk:"dscp_mark"`
	Rtcp                               types.Bool                                       `tfsdk:"rtcp"`
	FaxT38                             types.Bool                                       `tfsdk:"fax_t38"`
	BitRate                            types.String                                     `tfsdk:"bit_rate"`
	CancelCallWaiting                  types.String                                     `tfsdk:"cancel_call_waiting"`
	CallHold                           types.String                                     `tfsdk:"call_hold"`
	CidsActivate                       types.String                                     `tfsdk:"cids_activate"`
	CidsDeactivate                     types.String                                     `tfsdk:"cids_deactivate"`
	DoNotDisturbActivate               types.String                                     `tfsdk:"do_not_disturb_activate"`
	DoNotDisturbDeactivate             types.String                                     `tfsdk:"do_not_disturb_deactivate"`
	DoNotDisturbPinChange              types.String                                     `tfsdk:"do_not_disturb_pin_change"`
	EmergencyServiceNumber             types.String                                     `tfsdk:"emergency_service_number"`
	AnonCidBlockActivate               types.String                                     `tfsdk:"anon_cid_block_activate"`
	AnonCidBlockDeactivate             types.String                                     `tfsdk:"anon_cid_block_deactivate"`
	CallForwardUnconditionalActivate   types.String                                     `tfsdk:"call_forward_unconditional_activate"`
	CallForwardUnconditionalDeactivate types.String                                     `tfsdk:"call_forward_unconditional_deactivate"`
	CallForwardOnBusyActivate          types.String                                     `tfsdk:"call_forward_on_busy_activate"`
	CallForwardOnBusyDeactivate        types.String                                     `tfsdk:"call_forward_on_busy_deactivate"`
	CallForwardOnNoAnswerActivate      types.String                                     `tfsdk:"call_forward_on_no_answer_activate"`
	CallForwardOnNoAnswerDeactivate    types.String                                     `tfsdk:"call_forward_on_no_answer_deactivate"`
	Intercom1                          types.String                                     `tfsdk:"intercom_1"`
	Intercom2                          types.String                                     `tfsdk:"intercom_2"`
	Intercom3                          types.String                                     `tfsdk:"intercom_3"`
	Codecs                             []verityDeviceVoiceSettingsCodecModel            `tfsdk:"codecs"`
	ObjectProperties                   []verityDeviceVoiceSettingsObjectPropertiesModel `tfsdk:"object_properties"`
}

type verityDeviceVoiceSettingsCodecModel struct {
	CodecNumName                types.String `tfsdk:"codec_num_name"`
	CodecNumEnable              types.Bool   `tfsdk:"codec_num_enable"`
	CodecNumPacketizationPeriod types.String `tfsdk:"codec_num_packetization_period"`
	CodecNumSilenceSuppression  types.Bool   `tfsdk:"codec_num_silence_suppression"`
	Index                       types.Int64  `tfsdk:"index"`
}

func (c verityDeviceVoiceSettingsCodecModel) GetIndex() types.Int64 {
	return c.Index
}

type verityDeviceVoiceSettingsObjectPropertiesModel struct {
	Group types.String `tfsdk:"group"`
}

func (r *verityDeviceVoiceSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_voice_settings"
}

func (r *verityDeviceVoiceSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityDeviceVoiceSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity Device Voice Settings",
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
			"dtmf_method": schema.StringAttribute{
				Description: "Specifies how DTMF signals are carried",
				Optional:    true,
				Computed:    true,
			},
			"region": schema.StringAttribute{
				Description: "Region",
				Optional:    true,
				Computed:    true,
			},
			"protocol": schema.StringAttribute{
				Description: "Voice Protocol: MGCP or SIP",
				Optional:    true,
				Computed:    true,
			},
			"proxy_server": schema.StringAttribute{
				Description: "IP address or URI of the SIP proxy server for SIP signalling messages",
				Optional:    true,
				Computed:    true,
			},
			"proxy_server_port": schema.Int64Attribute{
				Description: "Proxy Server Port",
				Optional:    true,
				Computed:    true,
			},
			"proxy_server_secondary": schema.StringAttribute{
				Description: "IP address or URI of the secondary SIP proxy server for SIP signalling messages",
				Optional:    true,
				Computed:    true,
			},
			"proxy_server_secondary_port": schema.Int64Attribute{
				Description: "Secondary Proxy Server Port",
				Optional:    true,
				Computed:    true,
			},
			"registrar_server": schema.StringAttribute{
				Description: "Name or IP address or resolved name of the registrar server for SIP signalling messages",
				Optional:    true,
				Computed:    true,
			},
			"registrar_server_port": schema.Int64Attribute{
				Description: "Registrar Server Port",
				Optional:    true,
				Computed:    true,
			},
			"registrar_server_secondary": schema.StringAttribute{
				Description: "Name or IP address or resolved name of the secondary registrar server for SIP signalling messages",
				Optional:    true,
				Computed:    true,
			},
			"registrar_server_secondary_port": schema.Int64Attribute{
				Description: "Secondary Registrar Server Port",
				Optional:    true,
				Computed:    true,
			},
			"user_agent_domain": schema.StringAttribute{
				Description: "User Agent Domain",
				Optional:    true,
				Computed:    true,
			},
			"user_agent_transport": schema.StringAttribute{
				Description: "User Agent Transport",
				Optional:    true,
				Computed:    true,
			},
			"user_agent_port": schema.Int64Attribute{
				Description: "User Agent Port",
				Optional:    true,
				Computed:    true,
			},
			"outbound_proxy": schema.StringAttribute{
				Description: "IP address or URI of the outbound proxy server for SIP signalling messages",
				Optional:    true,
				Computed:    true,
			},
			"outbound_proxy_port": schema.Int64Attribute{
				Description: "Outbound Proxy Port",
				Optional:    true,
				Computed:    true,
			},
			"outbound_proxy_secondary": schema.StringAttribute{
				Description: "IP address or URI of the secondary outbound proxy server for SIP signalling messages",
				Optional:    true,
				Computed:    true,
			},
			"outbound_proxy_secondary_port": schema.Int64Attribute{
				Description: "Secondary Outbound Proxy Port",
				Optional:    true,
				Computed:    true,
			},
			"registration_period": schema.Int64Attribute{
				Description: "Specifies the time in seconds to start the re-registration process. The default value is 3240 seconds",
				Optional:    true,
				Computed:    true,
			},
			"register_expires": schema.Int64Attribute{
				Description: "SIP registration expiration time in seconds. If value is 0, the SIP agent does not add an expiration time to the registration requests and does not perform re-registration. The default value is 3600 seconds",
				Optional:    true,
				Computed:    true,
			},
			"voicemail_server": schema.StringAttribute{
				Description: "Name or IP address or resolved name of the external voicemail server if not provided by SIP server for MWI control",
				Optional:    true,
				Computed:    true,
			},
			"voicemail_server_port": schema.Int64Attribute{
				Description: "Voicemail Server Port",
				Optional:    true,
				Computed:    true,
			},
			"voicemail_server_expires": schema.Int64Attribute{
				Description: "Voicemail server expiration time in seconds. If value is 0, the Register Expires time is used instead. The default value is 3600 seconds",
				Optional:    true,
				Computed:    true,
			},
			"sip_dscp_mark": schema.Int64Attribute{
				Description: "SIP Differentiated Services Code point (DSCP)",
				Optional:    true,
				Computed:    true,
			},
			"call_agent_1": schema.StringAttribute{
				Description: "Call Agent 1",
				Optional:    true,
				Computed:    true,
			},
			"call_agent_port_1": schema.Int64Attribute{
				Description: "Call Agent Port 1",
				Optional:    true,
				Computed:    true,
			},
			"call_agent_2": schema.StringAttribute{
				Description: "Call Agent 2",
				Optional:    true,
				Computed:    true,
			},
			"call_agent_port_2": schema.Int64Attribute{
				Description: "Call Agent Port 2",
				Optional:    true,
				Computed:    true,
			},
			"domain": schema.StringAttribute{
				Description: "Domain",
				Optional:    true,
				Computed:    true,
			},
			"mgcp_dscp_mark": schema.Int64Attribute{
				Description: "MGCP Differentiated Services Code point (DSCP)",
				Optional:    true,
				Computed:    true,
			},
			"termination_base": schema.StringAttribute{
				Description: "Base string for the MGCP physical termination id(s)",
				Optional:    true,
				Computed:    true,
			},
			"local_port_min": schema.Int64Attribute{
				Description: "Defines the base RTP port that should be used for voice traffic",
				Optional:    true,
				Computed:    true,
			},
			"local_port_max": schema.Int64Attribute{
				Description: "Defines the highest RTP port used for voice traffic, must be greater than local Local Port Min",
				Optional:    true,
				Computed:    true,
			},
			"event_payload_type": schema.Int64Attribute{
				Description: "Telephone Event Payload Type",
				Optional:    true,
				Computed:    true,
			},
			"cas_events": schema.Int64Attribute{
				Description: "Enables or disables handling of CAS via RTP CAS events. Valid values are 0 = off and 1 = on",
				Optional:    true,
				Computed:    true,
			},
			"dscp_mark": schema.Int64Attribute{
				Description: "Differentiated Services Code Point (DSCP) to be used for outgoing RTP packets",
				Optional:    true,
				Computed:    true,
			},
			"rtcp": schema.BoolAttribute{
				Description: "RTCP Enable",
				Optional:    true,
				Computed:    true,
			},
			"fax_t38": schema.BoolAttribute{
				Description: "Fax T.38 Enable",
				Optional:    true,
				Computed:    true,
			},
			"bit_rate": schema.StringAttribute{
				Description: "T.38 Bit Rate in bps. Most available fax machines support up to 14,400bps",
				Optional:    true,
				Computed:    true,
			},
			"cancel_call_waiting": schema.StringAttribute{
				Description: "Cancel Call waiting",
				Optional:    true,
				Computed:    true,
			},
			"call_hold": schema.StringAttribute{
				Description: "Call hold",
				Optional:    true,
				Computed:    true,
			},
			"cids_activate": schema.StringAttribute{
				Description: "Caller ID Delivery Blocking (single call) Activate",
				Optional:    true,
				Computed:    true,
			},
			"cids_deactivate": schema.StringAttribute{
				Description: "Caller ID Delivery Blocking (single call) Deactivate",
				Optional:    true,
				Computed:    true,
			},
			"do_not_disturb_activate": schema.StringAttribute{
				Description: "Do not Disturb Activate",
				Optional:    true,
				Computed:    true,
			},
			"do_not_disturb_deactivate": schema.StringAttribute{
				Description: "Do not Disturb Deactivate",
				Optional:    true,
				Computed:    true,
			},
			"do_not_disturb_pin_change": schema.StringAttribute{
				Description: "Do not Disturb PIN Change",
				Optional:    true,
				Computed:    true,
			},
			"emergency_service_number": schema.StringAttribute{
				Description: "Emergency Service Number",
				Optional:    true,
				Computed:    true,
			},
			"anon_cid_block_activate": schema.StringAttribute{
				Description: "Anonymous Caller ID Block Activate",
				Optional:    true,
				Computed:    true,
			},
			"anon_cid_block_deactivate": schema.StringAttribute{
				Description: "Anonymous Caller ID Block Deactivate",
				Optional:    true,
				Computed:    true,
			},
			"call_forward_unconditional_activate": schema.StringAttribute{
				Description: "Call Forward Unconditional Activate",
				Optional:    true,
				Computed:    true,
			},
			"call_forward_unconditional_deactivate": schema.StringAttribute{
				Description: "Call Forward Unconditional Deactivate",
				Optional:    true,
				Computed:    true,
			},
			"call_forward_on_busy_activate": schema.StringAttribute{
				Description: "Call Forward On Busy Activate",
				Optional:    true,
				Computed:    true,
			},
			"call_forward_on_busy_deactivate": schema.StringAttribute{
				Description: "Call Forward On Busy Deactivate",
				Optional:    true,
				Computed:    true,
			},
			"call_forward_on_no_answer_activate": schema.StringAttribute{
				Description: "Call Forward On No Answer Activate",
				Optional:    true,
				Computed:    true,
			},
			"call_forward_on_no_answer_deactivate": schema.StringAttribute{
				Description: "Call Forward On No Answer Deactivate",
				Optional:    true,
				Computed:    true,
			},
			"intercom_1": schema.StringAttribute{
				Description: "Intercom 1",
				Optional:    true,
				Computed:    true,
			},
			"intercom_2": schema.StringAttribute{
				Description: "Intercom 2",
				Optional:    true,
				Computed:    true,
			},
			"intercom_3": schema.StringAttribute{
				Description: "Intercom 3",
				Optional:    true,
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"codecs": schema.ListNestedBlock{
				Description: "Codec configurations",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"codec_num_name": schema.StringAttribute{
							Description: "Name of this Codec",
							Optional:    true,
							Computed:    true,
						},
						"codec_num_enable": schema.BoolAttribute{
							Description: "Enable Codec",
							Optional:    true,
							Computed:    true,
						},
						"codec_num_packetization_period": schema.StringAttribute{
							Description: "Packet period selection interval in milliseconds",
							Optional:    true,
							Computed:    true,
						},
						"codec_num_silence_suppression": schema.BoolAttribute{
							Description: "Specifies whether silence suppression is on or off. Valid values are 0 = off and 1 = on",
							Optional:    true,
							Computed:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object. Zero if you want to add an object to the list.",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"object_properties": schema.ListNestedBlock{
				Description: "Object properties for the device voice settings",
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

func (r *verityDeviceVoiceSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityDeviceVoiceSettingsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config verityDeviceVoiceSettingsResourceModel
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
	dvsProps := &openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValue{
		Name: openapi.PtrString(name),
	}

	// Handle string fields
	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "DtmfMethod", APIField: &dvsProps.DtmfMethod, TFValue: plan.DtmfMethod},
		{FieldName: "Region", APIField: &dvsProps.Region, TFValue: plan.Region},
		{FieldName: "Protocol", APIField: &dvsProps.Protocol, TFValue: plan.Protocol},
		{FieldName: "ProxyServer", APIField: &dvsProps.ProxyServer, TFValue: plan.ProxyServer},
		{FieldName: "ProxyServerSecondary", APIField: &dvsProps.ProxyServerSecondary, TFValue: plan.ProxyServerSecondary},
		{FieldName: "RegistrarServer", APIField: &dvsProps.RegistrarServer, TFValue: plan.RegistrarServer},
		{FieldName: "RegistrarServerSecondary", APIField: &dvsProps.RegistrarServerSecondary, TFValue: plan.RegistrarServerSecondary},
		{FieldName: "UserAgentDomain", APIField: &dvsProps.UserAgentDomain, TFValue: plan.UserAgentDomain},
		{FieldName: "UserAgentTransport", APIField: &dvsProps.UserAgentTransport, TFValue: plan.UserAgentTransport},
		{FieldName: "OutboundProxy", APIField: &dvsProps.OutboundProxy, TFValue: plan.OutboundProxy},
		{FieldName: "OutboundProxySecondary", APIField: &dvsProps.OutboundProxySecondary, TFValue: plan.OutboundProxySecondary},
		{FieldName: "VoicemailServer", APIField: &dvsProps.VoicemailServer, TFValue: plan.VoicemailServer},
		{FieldName: "CallAgent1", APIField: &dvsProps.CallAgent1, TFValue: plan.CallAgent1},
		{FieldName: "CallAgent2", APIField: &dvsProps.CallAgent2, TFValue: plan.CallAgent2},
		{FieldName: "Domain", APIField: &dvsProps.Domain, TFValue: plan.Domain},
		{FieldName: "TerminationBase", APIField: &dvsProps.TerminationBase, TFValue: plan.TerminationBase},
		{FieldName: "BitRate", APIField: &dvsProps.BitRate, TFValue: plan.BitRate},
		{FieldName: "CancelCallWaiting", APIField: &dvsProps.CancelCallWaiting, TFValue: plan.CancelCallWaiting},
		{FieldName: "CallHold", APIField: &dvsProps.CallHold, TFValue: plan.CallHold},
		{FieldName: "CidsActivate", APIField: &dvsProps.CidsActivate, TFValue: plan.CidsActivate},
		{FieldName: "CidsDeactivate", APIField: &dvsProps.CidsDeactivate, TFValue: plan.CidsDeactivate},
		{FieldName: "DoNotDisturbActivate", APIField: &dvsProps.DoNotDisturbActivate, TFValue: plan.DoNotDisturbActivate},
		{FieldName: "DoNotDisturbDeactivate", APIField: &dvsProps.DoNotDisturbDeactivate, TFValue: plan.DoNotDisturbDeactivate},
		{FieldName: "DoNotDisturbPinChange", APIField: &dvsProps.DoNotDisturbPinChange, TFValue: plan.DoNotDisturbPinChange},
		{FieldName: "EmergencyServiceNumber", APIField: &dvsProps.EmergencyServiceNumber, TFValue: plan.EmergencyServiceNumber},
		{FieldName: "AnonCidBlockActivate", APIField: &dvsProps.AnonCidBlockActivate, TFValue: plan.AnonCidBlockActivate},
		{FieldName: "AnonCidBlockDeactivate", APIField: &dvsProps.AnonCidBlockDeactivate, TFValue: plan.AnonCidBlockDeactivate},
		{FieldName: "CallForwardUnconditionalActivate", APIField: &dvsProps.CallForwardUnconditionalActivate, TFValue: plan.CallForwardUnconditionalActivate},
		{FieldName: "CallForwardUnconditionalDeactivate", APIField: &dvsProps.CallForwardUnconditionalDeactivate, TFValue: plan.CallForwardUnconditionalDeactivate},
		{FieldName: "CallForwardOnBusyActivate", APIField: &dvsProps.CallForwardOnBusyActivate, TFValue: plan.CallForwardOnBusyActivate},
		{FieldName: "CallForwardOnBusyDeactivate", APIField: &dvsProps.CallForwardOnBusyDeactivate, TFValue: plan.CallForwardOnBusyDeactivate},
		{FieldName: "CallForwardOnNoAnswerActivate", APIField: &dvsProps.CallForwardOnNoAnswerActivate, TFValue: plan.CallForwardOnNoAnswerActivate},
		{FieldName: "CallForwardOnNoAnswerDeactivate", APIField: &dvsProps.CallForwardOnNoAnswerDeactivate, TFValue: plan.CallForwardOnNoAnswerDeactivate},
		{FieldName: "Intercom1", APIField: &dvsProps.Intercom1, TFValue: plan.Intercom1},
		{FieldName: "Intercom2", APIField: &dvsProps.Intercom2, TFValue: plan.Intercom2},
		{FieldName: "Intercom3", APIField: &dvsProps.Intercom3, TFValue: plan.Intercom3},
	})

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &dvsProps.Enable, TFValue: plan.Enable},
		{FieldName: "Rtcp", APIField: &dvsProps.Rtcp, TFValue: plan.Rtcp},
		{FieldName: "FaxT38", APIField: &dvsProps.FaxT38, TFValue: plan.FaxT38},
	})

	// Handle nullable int64 fields - parse HCL to detect explicit config
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, deviceVoiceSettingsTerraformType, name)

	utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
		{FieldName: "ProxyServerPort", APIField: &dvsProps.ProxyServerPort, TFValue: config.ProxyServerPort, IsConfigured: configuredAttrs.IsConfigured("proxy_server_port")},
		{FieldName: "ProxyServerSecondaryPort", APIField: &dvsProps.ProxyServerSecondaryPort, TFValue: config.ProxyServerSecondaryPort, IsConfigured: configuredAttrs.IsConfigured("proxy_server_secondary_port")},
		{FieldName: "RegistrarServerPort", APIField: &dvsProps.RegistrarServerPort, TFValue: config.RegistrarServerPort, IsConfigured: configuredAttrs.IsConfigured("registrar_server_port")},
		{FieldName: "RegistrarServerSecondaryPort", APIField: &dvsProps.RegistrarServerSecondaryPort, TFValue: config.RegistrarServerSecondaryPort, IsConfigured: configuredAttrs.IsConfigured("registrar_server_secondary_port")},
		{FieldName: "UserAgentPort", APIField: &dvsProps.UserAgentPort, TFValue: config.UserAgentPort, IsConfigured: configuredAttrs.IsConfigured("user_agent_port")},
		{FieldName: "OutboundProxyPort", APIField: &dvsProps.OutboundProxyPort, TFValue: config.OutboundProxyPort, IsConfigured: configuredAttrs.IsConfigured("outbound_proxy_port")},
		{FieldName: "OutboundProxySecondaryPort", APIField: &dvsProps.OutboundProxySecondaryPort, TFValue: config.OutboundProxySecondaryPort, IsConfigured: configuredAttrs.IsConfigured("outbound_proxy_secondary_port")},
		{FieldName: "RegistrationPeriod", APIField: &dvsProps.RegistrationPeriod, TFValue: config.RegistrationPeriod, IsConfigured: configuredAttrs.IsConfigured("registration_period")},
		{FieldName: "RegisterExpires", APIField: &dvsProps.RegisterExpires, TFValue: config.RegisterExpires, IsConfigured: configuredAttrs.IsConfigured("register_expires")},
		{FieldName: "VoicemailServerPort", APIField: &dvsProps.VoicemailServerPort, TFValue: config.VoicemailServerPort, IsConfigured: configuredAttrs.IsConfigured("voicemail_server_port")},
		{FieldName: "VoicemailServerExpires", APIField: &dvsProps.VoicemailServerExpires, TFValue: config.VoicemailServerExpires, IsConfigured: configuredAttrs.IsConfigured("voicemail_server_expires")},
		{FieldName: "SipDscpMark", APIField: &dvsProps.SipDscpMark, TFValue: config.SipDscpMark, IsConfigured: configuredAttrs.IsConfigured("sip_dscp_mark")},
		{FieldName: "CallAgentPort1", APIField: &dvsProps.CallAgentPort1, TFValue: config.CallAgentPort1, IsConfigured: configuredAttrs.IsConfigured("call_agent_port_1")},
		{FieldName: "CallAgentPort2", APIField: &dvsProps.CallAgentPort2, TFValue: config.CallAgentPort2, IsConfigured: configuredAttrs.IsConfigured("call_agent_port_2")},
		{FieldName: "MgcpDscpMark", APIField: &dvsProps.MgcpDscpMark, TFValue: config.MgcpDscpMark, IsConfigured: configuredAttrs.IsConfigured("mgcp_dscp_mark")},
		{FieldName: "LocalPortMin", APIField: &dvsProps.LocalPortMin, TFValue: config.LocalPortMin, IsConfigured: configuredAttrs.IsConfigured("local_port_min")},
		{FieldName: "LocalPortMax", APIField: &dvsProps.LocalPortMax, TFValue: config.LocalPortMax, IsConfigured: configuredAttrs.IsConfigured("local_port_max")},
		{FieldName: "EventPayloadType", APIField: &dvsProps.EventPayloadType, TFValue: config.EventPayloadType, IsConfigured: configuredAttrs.IsConfigured("event_payload_type")},
		{FieldName: "CasEvents", APIField: &dvsProps.CasEvents, TFValue: config.CasEvents, IsConfigured: configuredAttrs.IsConfigured("cas_events")},
		{FieldName: "DscpMark", APIField: &dvsProps.DscpMark, TFValue: config.DscpMark, IsConfigured: configuredAttrs.IsConfigured("dscp_mark")},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties{}
		utils.SetObjectPropertiesFields([]utils.ObjectPropertiesField{
			{Name: "Group", TFValue: op.Group, APIValue: &objProps.Group},
		})
		dvsProps.ObjectProperties = &objProps
	}

	// Handle codecs
	if len(plan.Codecs) > 0 {
		codecs := make([]openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner, len(plan.Codecs))
		for i, item := range plan.Codecs {
			codecItem := openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner{}
			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "CodecNumEnable", APIField: &codecItem.CodecNumEnable, TFValue: item.CodecNumEnable},
				{FieldName: "CodecNumSilenceSuppression", APIField: &codecItem.CodecNumSilenceSuppression, TFValue: item.CodecNumSilenceSuppression},
			})
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "CodecNumName", APIField: &codecItem.CodecNumName, TFValue: item.CodecNumName},
				{FieldName: "CodecNumPacketizationPeriod", APIField: &codecItem.CodecNumPacketizationPeriod, TFValue: item.CodecNumPacketizationPeriod},
			})
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &codecItem.Index, TFValue: item.Index},
			})
			codecs[i] = codecItem
		}
		dvsProps.Codecs = codecs
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "device_voice_settings", name, *dvsProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Device Voice Settings %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "device_voice_settings")

	var minState verityDeviceVoiceSettingsResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if dvsData, exists := bulkMgr.GetResourceResponse("device_voice_settings", name); exists {
			state := populateDeviceVoiceSettingsState(ctx, minState, dvsData, r.provCtx.mode)
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

func (r *verityDeviceVoiceSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityDeviceVoiceSettingsResourceModel
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

	dvsName := state.Name.ValueString()

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if dvsData, exists := r.bulkOpsMgr.GetResourceResponse("device_voice_settings", dvsName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached device voice settings data for %s from recent operation", dvsName))
			state = populateDeviceVoiceSettingsState(ctx, state, dvsData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("device_voice_settings") {
		tflog.Info(ctx, fmt.Sprintf("Skipping Device Voice Settings %s verification – trusting recent successful API operation", dvsName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching Device Voice Settings for verification of %s", dvsName))

	type DeviceVoiceSettingsResponse struct {
		DeviceVoiceSettings map[string]interface{} `json:"device_voice_settings"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "device_voice_settings", dvsName,
		func() (DeviceVoiceSettingsResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch Device Voice Settings")
			respAPI, err := r.client.DeviceVoiceSettingsAPI.DevicevoicesettingsGet(ctx).Execute()
			if err != nil {
				return DeviceVoiceSettingsResponse{}, fmt.Errorf("error reading Device Voice Settings: %v", err)
			}
			defer respAPI.Body.Close()

			var res DeviceVoiceSettingsResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return DeviceVoiceSettingsResponse{}, fmt.Errorf("failed to decode Device Voice Settings response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d Device Voice Settings", len(res.DeviceVoiceSettings)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Device Voice Settings %s", dvsName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for Device Voice Settings with name: %s", dvsName))

	dvsData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.DeviceVoiceSettings,
		dvsName,
		func(data interface{}) (string, bool) {
			if dvs, ok := data.(map[string]interface{}); ok {
				if name, ok := dvs["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Device Voice Settings with name '%s' not found in API response", dvsName))
		resp.State.RemoveResource(ctx)
		return
	}

	dvsMap, ok := dvsData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Device Voice Settings Data",
			fmt.Sprintf("Device Voice Settings data is not in expected format for %s", dvsName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found Device Voice Settings '%s' under API key '%s'", dvsName, actualAPIName))

	state = populateDeviceVoiceSettingsState(ctx, state, dvsMap, r.provCtx.mode)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityDeviceVoiceSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityDeviceVoiceSettingsResourceModel

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
	var config verityDeviceVoiceSettingsResourceModel
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
	dvsProps := openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValue{}
	hasChanges := false

	// Parse HCL to detect which fields are explicitly configured
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, deviceVoiceSettingsTerraformType, name)

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { dvsProps.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.DtmfMethod, state.DtmfMethod, func(v *string) { dvsProps.DtmfMethod = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Region, state.Region, func(v *string) { dvsProps.Region = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Protocol, state.Protocol, func(v *string) { dvsProps.Protocol = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.ProxyServer, state.ProxyServer, func(v *string) { dvsProps.ProxyServer = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.ProxyServerSecondary, state.ProxyServerSecondary, func(v *string) { dvsProps.ProxyServerSecondary = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.RegistrarServer, state.RegistrarServer, func(v *string) { dvsProps.RegistrarServer = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.RegistrarServerSecondary, state.RegistrarServerSecondary, func(v *string) { dvsProps.RegistrarServerSecondary = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.UserAgentDomain, state.UserAgentDomain, func(v *string) { dvsProps.UserAgentDomain = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.UserAgentTransport, state.UserAgentTransport, func(v *string) { dvsProps.UserAgentTransport = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.OutboundProxy, state.OutboundProxy, func(v *string) { dvsProps.OutboundProxy = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.OutboundProxySecondary, state.OutboundProxySecondary, func(v *string) { dvsProps.OutboundProxySecondary = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.VoicemailServer, state.VoicemailServer, func(v *string) { dvsProps.VoicemailServer = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.CallAgent1, state.CallAgent1, func(v *string) { dvsProps.CallAgent1 = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.CallAgent2, state.CallAgent2, func(v *string) { dvsProps.CallAgent2 = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Domain, state.Domain, func(v *string) { dvsProps.Domain = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.TerminationBase, state.TerminationBase, func(v *string) { dvsProps.TerminationBase = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.BitRate, state.BitRate, func(v *string) { dvsProps.BitRate = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.CancelCallWaiting, state.CancelCallWaiting, func(v *string) { dvsProps.CancelCallWaiting = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.CallHold, state.CallHold, func(v *string) { dvsProps.CallHold = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.CidsActivate, state.CidsActivate, func(v *string) { dvsProps.CidsActivate = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.CidsDeactivate, state.CidsDeactivate, func(v *string) { dvsProps.CidsDeactivate = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.DoNotDisturbActivate, state.DoNotDisturbActivate, func(v *string) { dvsProps.DoNotDisturbActivate = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.DoNotDisturbDeactivate, state.DoNotDisturbDeactivate, func(v *string) { dvsProps.DoNotDisturbDeactivate = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.DoNotDisturbPinChange, state.DoNotDisturbPinChange, func(v *string) { dvsProps.DoNotDisturbPinChange = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.EmergencyServiceNumber, state.EmergencyServiceNumber, func(v *string) { dvsProps.EmergencyServiceNumber = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.AnonCidBlockActivate, state.AnonCidBlockActivate, func(v *string) { dvsProps.AnonCidBlockActivate = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.AnonCidBlockDeactivate, state.AnonCidBlockDeactivate, func(v *string) { dvsProps.AnonCidBlockDeactivate = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.CallForwardUnconditionalActivate, state.CallForwardUnconditionalActivate, func(v *string) { dvsProps.CallForwardUnconditionalActivate = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.CallForwardUnconditionalDeactivate, state.CallForwardUnconditionalDeactivate, func(v *string) { dvsProps.CallForwardUnconditionalDeactivate = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.CallForwardOnBusyActivate, state.CallForwardOnBusyActivate, func(v *string) { dvsProps.CallForwardOnBusyActivate = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.CallForwardOnBusyDeactivate, state.CallForwardOnBusyDeactivate, func(v *string) { dvsProps.CallForwardOnBusyDeactivate = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.CallForwardOnNoAnswerActivate, state.CallForwardOnNoAnswerActivate, func(v *string) { dvsProps.CallForwardOnNoAnswerActivate = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.CallForwardOnNoAnswerDeactivate, state.CallForwardOnNoAnswerDeactivate, func(v *string) { dvsProps.CallForwardOnNoAnswerDeactivate = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Intercom1, state.Intercom1, func(v *string) { dvsProps.Intercom1 = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Intercom2, state.Intercom2, func(v *string) { dvsProps.Intercom2 = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Intercom3, state.Intercom3, func(v *string) { dvsProps.Intercom3 = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { dvsProps.Enable = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.Rtcp, state.Rtcp, func(v *bool) { dvsProps.Rtcp = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.FaxT38, state.FaxT38, func(v *bool) { dvsProps.FaxT38 = v }, &hasChanges)

	// Handle nullable int64 field changes - parse HCL to detect explicit config
	utils.CompareAndSetNullableInt64Field(config.ProxyServerPort, state.ProxyServerPort, configuredAttrs.IsConfigured("proxy_server_port"), func(v *openapi.NullableInt32) { dvsProps.ProxyServerPort = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.ProxyServerSecondaryPort, state.ProxyServerSecondaryPort, configuredAttrs.IsConfigured("proxy_server_secondary_port"), func(v *openapi.NullableInt32) { dvsProps.ProxyServerSecondaryPort = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.RegistrarServerPort, state.RegistrarServerPort, configuredAttrs.IsConfigured("registrar_server_port"), func(v *openapi.NullableInt32) { dvsProps.RegistrarServerPort = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.RegistrarServerSecondaryPort, state.RegistrarServerSecondaryPort, configuredAttrs.IsConfigured("registrar_server_secondary_port"), func(v *openapi.NullableInt32) { dvsProps.RegistrarServerSecondaryPort = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.UserAgentPort, state.UserAgentPort, configuredAttrs.IsConfigured("user_agent_port"), func(v *openapi.NullableInt32) { dvsProps.UserAgentPort = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.OutboundProxyPort, state.OutboundProxyPort, configuredAttrs.IsConfigured("outbound_proxy_port"), func(v *openapi.NullableInt32) { dvsProps.OutboundProxyPort = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.OutboundProxySecondaryPort, state.OutboundProxySecondaryPort, configuredAttrs.IsConfigured("outbound_proxy_secondary_port"), func(v *openapi.NullableInt32) { dvsProps.OutboundProxySecondaryPort = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.RegistrationPeriod, state.RegistrationPeriod, configuredAttrs.IsConfigured("registration_period"), func(v *openapi.NullableInt32) { dvsProps.RegistrationPeriod = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.RegisterExpires, state.RegisterExpires, configuredAttrs.IsConfigured("register_expires"), func(v *openapi.NullableInt32) { dvsProps.RegisterExpires = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.VoicemailServerPort, state.VoicemailServerPort, configuredAttrs.IsConfigured("voicemail_server_port"), func(v *openapi.NullableInt32) { dvsProps.VoicemailServerPort = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.VoicemailServerExpires, state.VoicemailServerExpires, configuredAttrs.IsConfigured("voicemail_server_expires"), func(v *openapi.NullableInt32) { dvsProps.VoicemailServerExpires = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.SipDscpMark, state.SipDscpMark, configuredAttrs.IsConfigured("sip_dscp_mark"), func(v *openapi.NullableInt32) { dvsProps.SipDscpMark = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.CallAgentPort1, state.CallAgentPort1, configuredAttrs.IsConfigured("call_agent_port_1"), func(v *openapi.NullableInt32) { dvsProps.CallAgentPort1 = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.CallAgentPort2, state.CallAgentPort2, configuredAttrs.IsConfigured("call_agent_port_2"), func(v *openapi.NullableInt32) { dvsProps.CallAgentPort2 = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.MgcpDscpMark, state.MgcpDscpMark, configuredAttrs.IsConfigured("mgcp_dscp_mark"), func(v *openapi.NullableInt32) { dvsProps.MgcpDscpMark = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.LocalPortMin, state.LocalPortMin, configuredAttrs.IsConfigured("local_port_min"), func(v *openapi.NullableInt32) { dvsProps.LocalPortMin = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.LocalPortMax, state.LocalPortMax, configuredAttrs.IsConfigured("local_port_max"), func(v *openapi.NullableInt32) { dvsProps.LocalPortMax = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.EventPayloadType, state.EventPayloadType, configuredAttrs.IsConfigured("event_payload_type"), func(v *openapi.NullableInt32) { dvsProps.EventPayloadType = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.CasEvents, state.CasEvents, configuredAttrs.IsConfigured("cas_events"), func(v *openapi.NullableInt32) { dvsProps.CasEvents = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.DscpMark, state.DscpMark, configuredAttrs.IsConfigured("dscp_mark"), func(v *openapi.NullableInt32) { dvsProps.DscpMark = *v }, &hasChanges)

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
			dvsProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	// Handle codecs
	codecsHandler := utils.IndexedItemHandler[verityDeviceVoiceSettingsCodecModel, openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner]{
		CreateNew: func(planItem verityDeviceVoiceSettingsCodecModel) openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner {
			codec := openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &codec.Index, TFValue: planItem.Index},
			})

			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "CodecNumName", APIField: &codec.CodecNumName, TFValue: planItem.CodecNumName},
				{FieldName: "CodecNumPacketizationPeriod", APIField: &codec.CodecNumPacketizationPeriod, TFValue: planItem.CodecNumPacketizationPeriod},
			})

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "CodecNumEnable", APIField: &codec.CodecNumEnable, TFValue: planItem.CodecNumEnable},
				{FieldName: "CodecNumSilenceSuppression", APIField: &codec.CodecNumSilenceSuppression, TFValue: planItem.CodecNumSilenceSuppression},
			})

			return codec
		},
		UpdateExisting: func(planItem verityDeviceVoiceSettingsCodecModel, stateItem verityDeviceVoiceSettingsCodecModel) (openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner, bool) {
			codec := openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &codec.Index, TFValue: planItem.Index},
			})

			fieldChanged := false

			// Handle string fields
			utils.CompareAndSetStringField(planItem.CodecNumName, stateItem.CodecNumName, func(v *string) { codec.CodecNumName = v }, &fieldChanged)
			utils.CompareAndSetStringField(planItem.CodecNumPacketizationPeriod, stateItem.CodecNumPacketizationPeriod, func(v *string) { codec.CodecNumPacketizationPeriod = v }, &fieldChanged)

			// Handle boolean fields
			utils.CompareAndSetBoolField(planItem.CodecNumEnable, stateItem.CodecNumEnable, func(v *bool) { codec.CodecNumEnable = v }, &fieldChanged)
			utils.CompareAndSetBoolField(planItem.CodecNumSilenceSuppression, stateItem.CodecNumSilenceSuppression, func(v *bool) { codec.CodecNumSilenceSuppression = v }, &fieldChanged)

			return codec, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner {
			codec := openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner{}
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &codec.Index, TFValue: types.Int64Value(index)},
			})
			return codec
		},
	}

	changedCodecs, codecsChanged := utils.ProcessIndexedArrayUpdates(plan.Codecs, state.Codecs, codecsHandler)
	if codecsChanged {
		dvsProps.Codecs = changedCodecs
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "device_voice_settings", name, dvsProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Device Voice Settings %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "device_voice_settings")

	var minState verityDeviceVoiceSettingsResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Try to use cached response from bulk operation to populate state with API values
	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if dvsData, exists := bulkMgr.GetResourceResponse("device_voice_settings", name); exists {
			newState := populateDeviceVoiceSettingsState(ctx, minState, dvsData, r.provCtx.mode)
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

func (r *verityDeviceVoiceSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityDeviceVoiceSettingsResourceModel
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "device_voice_settings", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Device Voice Settings %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "device_voice_settings")
	resp.State.RemoveResource(ctx)
}

func (r *verityDeviceVoiceSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func populateDeviceVoiceSettingsState(ctx context.Context, state verityDeviceVoiceSettingsResourceModel, data map[string]interface{}, mode string) verityDeviceVoiceSettingsResourceModel {
	const resourceType = deviceVoiceSettingsResourceType

	state.Name = utils.MapStringFromAPI(data["name"])

	// Boolean fields
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)
	state.Rtcp = utils.MapBoolWithMode(data, "rtcp", resourceType, mode)
	state.FaxT38 = utils.MapBoolWithMode(data, "fax_t38", resourceType, mode)

	// String fields
	state.DtmfMethod = utils.MapStringWithMode(data, "dtmf_method", resourceType, mode)
	state.Region = utils.MapStringWithMode(data, "region", resourceType, mode)
	state.Protocol = utils.MapStringWithMode(data, "protocol", resourceType, mode)
	state.ProxyServer = utils.MapStringWithMode(data, "proxy_server", resourceType, mode)
	state.ProxyServerSecondary = utils.MapStringWithMode(data, "proxy_server_secondary", resourceType, mode)
	state.RegistrarServer = utils.MapStringWithMode(data, "registrar_server", resourceType, mode)
	state.RegistrarServerSecondary = utils.MapStringWithMode(data, "registrar_server_secondary", resourceType, mode)
	state.UserAgentDomain = utils.MapStringWithMode(data, "user_agent_domain", resourceType, mode)
	state.UserAgentTransport = utils.MapStringWithMode(data, "user_agent_transport", resourceType, mode)
	state.OutboundProxy = utils.MapStringWithMode(data, "outbound_proxy", resourceType, mode)
	state.OutboundProxySecondary = utils.MapStringWithMode(data, "outbound_proxy_secondary", resourceType, mode)
	state.VoicemailServer = utils.MapStringWithMode(data, "voicemail_server", resourceType, mode)
	state.CallAgent1 = utils.MapStringWithMode(data, "call_agent_1", resourceType, mode)
	state.CallAgent2 = utils.MapStringWithMode(data, "call_agent_2", resourceType, mode)
	state.Domain = utils.MapStringWithMode(data, "domain", resourceType, mode)
	state.TerminationBase = utils.MapStringWithMode(data, "termination_base", resourceType, mode)
	state.BitRate = utils.MapStringWithMode(data, "bit_rate", resourceType, mode)
	state.CancelCallWaiting = utils.MapStringWithMode(data, "cancel_call_waiting", resourceType, mode)
	state.CallHold = utils.MapStringWithMode(data, "call_hold", resourceType, mode)
	state.CidsActivate = utils.MapStringWithMode(data, "cids_activate", resourceType, mode)
	state.CidsDeactivate = utils.MapStringWithMode(data, "cids_deactivate", resourceType, mode)
	state.DoNotDisturbActivate = utils.MapStringWithMode(data, "do_not_disturb_activate", resourceType, mode)
	state.DoNotDisturbDeactivate = utils.MapStringWithMode(data, "do_not_disturb_deactivate", resourceType, mode)
	state.DoNotDisturbPinChange = utils.MapStringWithMode(data, "do_not_disturb_pin_change", resourceType, mode)
	state.EmergencyServiceNumber = utils.MapStringWithMode(data, "emergency_service_number", resourceType, mode)
	state.AnonCidBlockActivate = utils.MapStringWithMode(data, "anon_cid_block_activate", resourceType, mode)
	state.AnonCidBlockDeactivate = utils.MapStringWithMode(data, "anon_cid_block_deactivate", resourceType, mode)
	state.CallForwardUnconditionalActivate = utils.MapStringWithMode(data, "call_forward_unconditional_activate", resourceType, mode)
	state.CallForwardUnconditionalDeactivate = utils.MapStringWithMode(data, "call_forward_unconditional_deactivate", resourceType, mode)
	state.CallForwardOnBusyActivate = utils.MapStringWithMode(data, "call_forward_on_busy_activate", resourceType, mode)
	state.CallForwardOnBusyDeactivate = utils.MapStringWithMode(data, "call_forward_on_busy_deactivate", resourceType, mode)
	state.CallForwardOnNoAnswerActivate = utils.MapStringWithMode(data, "call_forward_on_no_answer_activate", resourceType, mode)
	state.CallForwardOnNoAnswerDeactivate = utils.MapStringWithMode(data, "call_forward_on_no_answer_deactivate", resourceType, mode)
	state.Intercom1 = utils.MapStringWithMode(data, "intercom_1", resourceType, mode)
	state.Intercom2 = utils.MapStringWithMode(data, "intercom_2", resourceType, mode)
	state.Intercom3 = utils.MapStringWithMode(data, "intercom_3", resourceType, mode)

	// Int64 fields
	state.ProxyServerPort = utils.MapInt64WithMode(data, "proxy_server_port", resourceType, mode)
	state.ProxyServerSecondaryPort = utils.MapInt64WithMode(data, "proxy_server_secondary_port", resourceType, mode)
	state.RegistrarServerPort = utils.MapInt64WithMode(data, "registrar_server_port", resourceType, mode)
	state.RegistrarServerSecondaryPort = utils.MapInt64WithMode(data, "registrar_server_secondary_port", resourceType, mode)
	state.UserAgentPort = utils.MapInt64WithMode(data, "user_agent_port", resourceType, mode)
	state.OutboundProxyPort = utils.MapInt64WithMode(data, "outbound_proxy_port", resourceType, mode)
	state.OutboundProxySecondaryPort = utils.MapInt64WithMode(data, "outbound_proxy_secondary_port", resourceType, mode)
	state.RegistrationPeriod = utils.MapInt64WithMode(data, "registration_period", resourceType, mode)
	state.RegisterExpires = utils.MapInt64WithMode(data, "register_expires", resourceType, mode)
	state.VoicemailServerPort = utils.MapInt64WithMode(data, "voicemail_server_port", resourceType, mode)
	state.VoicemailServerExpires = utils.MapInt64WithMode(data, "voicemail_server_expires", resourceType, mode)
	state.SipDscpMark = utils.MapInt64WithMode(data, "sip_dscp_mark", resourceType, mode)
	state.CallAgentPort1 = utils.MapInt64WithMode(data, "call_agent_port_1", resourceType, mode)
	state.CallAgentPort2 = utils.MapInt64WithMode(data, "call_agent_port_2", resourceType, mode)
	state.MgcpDscpMark = utils.MapInt64WithMode(data, "mgcp_dscp_mark", resourceType, mode)
	state.LocalPortMin = utils.MapInt64WithMode(data, "local_port_min", resourceType, mode)
	state.LocalPortMax = utils.MapInt64WithMode(data, "local_port_max", resourceType, mode)
	state.EventPayloadType = utils.MapInt64WithMode(data, "event_payload_type", resourceType, mode)
	state.CasEvents = utils.MapInt64WithMode(data, "cas_events", resourceType, mode)
	state.DscpMark = utils.MapInt64WithMode(data, "dscp_mark", resourceType, mode)

	// Handle codecs array
	if utils.FieldAppliesToMode(resourceType, "codecs", mode) {
		if codecsArray, ok := data["codecs"].([]interface{}); ok && len(codecsArray) > 0 {
			var codecs []verityDeviceVoiceSettingsCodecModel
			for _, c := range codecsArray {
				codec, ok := c.(map[string]interface{})
				if !ok {
					continue
				}
				codecModel := verityDeviceVoiceSettingsCodecModel{
					CodecNumName:                utils.MapStringWithModeNested(codec, "codec_num_name", resourceType, "codecs.codec_num_name", mode),
					CodecNumEnable:              utils.MapBoolWithModeNested(codec, "codec_num_enable", resourceType, "codecs.codec_num_enable", mode),
					CodecNumPacketizationPeriod: utils.MapStringWithModeNested(codec, "codec_num_packetization_period", resourceType, "codecs.codec_num_packetization_period", mode),
					CodecNumSilenceSuppression:  utils.MapBoolWithModeNested(codec, "codec_num_silence_suppression", resourceType, "codecs.codec_num_silence_suppression", mode),
					Index:                       utils.MapInt64WithModeNested(codec, "index", resourceType, "codecs.index", mode),
				}
				codecs = append(codecs, codecModel)
			}
			state.Codecs = codecs
		} else {
			state.Codecs = nil
		}
	} else {
		state.Codecs = nil
	}

	// Handle object_properties block
	if utils.FieldAppliesToMode(resourceType, "object_properties", mode) {
		if objProps, ok := data["object_properties"].(map[string]interface{}); ok {
			objPropsModel := verityDeviceVoiceSettingsObjectPropertiesModel{
				Group: utils.MapStringWithModeNested(objProps, "group", resourceType, "object_properties.group", mode),
			}
			state.ObjectProperties = []verityDeviceVoiceSettingsObjectPropertiesModel{objPropsModel}
		} else {
			state.ObjectProperties = nil
		}
	} else {
		state.ObjectProperties = nil
	}

	return state
}

func (r *verityDeviceVoiceSettingsResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// =========================================================================
	// Skip if deleting
	// =========================================================================
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityDeviceVoiceSettingsResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Mode-aware field nullification
	// Set fields that don't apply to current mode to null to prevent
	// "known after apply" messages for irrelevant fields.
	// =========================================================================
	const resourceType = deviceVoiceSettingsResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyStrings(
		"dtmf_method", "region", "protocol", "proxy_server", "proxy_server_secondary",
		"registrar_server", "registrar_server_secondary", "user_agent_domain", "user_agent_transport",
		"outbound_proxy", "outbound_proxy_secondary", "voicemail_server",
		"call_agent_1", "call_agent_2", "domain", "termination_base", "bit_rate",
		"cancel_call_waiting", "call_hold", "cids_activate", "cids_deactivate",
		"do_not_disturb_activate", "do_not_disturb_deactivate", "do_not_disturb_pin_change",
		"emergency_service_number", "anon_cid_block_activate", "anon_cid_block_deactivate",
		"call_forward_unconditional_activate", "call_forward_unconditional_deactivate",
		"call_forward_on_busy_activate", "call_forward_on_busy_deactivate",
		"call_forward_on_no_answer_activate", "call_forward_on_no_answer_deactivate",
		"intercom_1", "intercom_2", "intercom_3",
	)

	nullifier.NullifyBools(
		"enable", "rtcp", "fax_t38",
	)

	nullifier.NullifyInt64s(
		"proxy_server_port", "proxy_server_secondary_port",
		"registrar_server_port", "registrar_server_secondary_port",
		"user_agent_port", "outbound_proxy_port", "outbound_proxy_secondary_port",
		"registration_period", "register_expires",
		"voicemail_server_port", "voicemail_server_expires", "sip_dscp_mark",
		"call_agent_port_1", "call_agent_port_2", "mgcp_dscp_mark",
		"local_port_min", "local_port_max", "event_payload_type", "cas_events", "dscp_mark",
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
	var state verityDeviceVoiceSettingsResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config verityDeviceVoiceSettingsResourceModel
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
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, deviceVoiceSettingsTerraformType, name)

	utils.HandleNullableFields(utils.NullableFieldsConfig{
		Ctx:             ctx,
		Plan:            &resp.Plan,
		ConfiguredAttrs: configuredAttrs,
		Int64Fields: []utils.NullableInt64Field{
			{AttrName: "proxy_server_port", ConfigVal: config.ProxyServerPort, StateVal: state.ProxyServerPort},
			{AttrName: "proxy_server_secondary_port", ConfigVal: config.ProxyServerSecondaryPort, StateVal: state.ProxyServerSecondaryPort},
			{AttrName: "registrar_server_port", ConfigVal: config.RegistrarServerPort, StateVal: state.RegistrarServerPort},
			{AttrName: "registrar_server_secondary_port", ConfigVal: config.RegistrarServerSecondaryPort, StateVal: state.RegistrarServerSecondaryPort},
			{AttrName: "user_agent_port", ConfigVal: config.UserAgentPort, StateVal: state.UserAgentPort},
			{AttrName: "outbound_proxy_port", ConfigVal: config.OutboundProxyPort, StateVal: state.OutboundProxyPort},
			{AttrName: "outbound_proxy_secondary_port", ConfigVal: config.OutboundProxySecondaryPort, StateVal: state.OutboundProxySecondaryPort},
			{AttrName: "registration_period", ConfigVal: config.RegistrationPeriod, StateVal: state.RegistrationPeriod},
			{AttrName: "register_expires", ConfigVal: config.RegisterExpires, StateVal: state.RegisterExpires},
			{AttrName: "voicemail_server_port", ConfigVal: config.VoicemailServerPort, StateVal: state.VoicemailServerPort},
			{AttrName: "voicemail_server_expires", ConfigVal: config.VoicemailServerExpires, StateVal: state.VoicemailServerExpires},
			{AttrName: "sip_dscp_mark", ConfigVal: config.SipDscpMark, StateVal: state.SipDscpMark},
			{AttrName: "call_agent_port_1", ConfigVal: config.CallAgentPort1, StateVal: state.CallAgentPort1},
			{AttrName: "call_agent_port_2", ConfigVal: config.CallAgentPort2, StateVal: state.CallAgentPort2},
			{AttrName: "mgcp_dscp_mark", ConfigVal: config.MgcpDscpMark, StateVal: state.MgcpDscpMark},
			{AttrName: "local_port_min", ConfigVal: config.LocalPortMin, StateVal: state.LocalPortMin},
			{AttrName: "local_port_max", ConfigVal: config.LocalPortMax, StateVal: state.LocalPortMax},
			{AttrName: "event_payload_type", ConfigVal: config.EventPayloadType, StateVal: state.EventPayloadType},
			{AttrName: "cas_events", ConfigVal: config.CasEvents, StateVal: state.CasEvents},
			{AttrName: "dscp_mark", ConfigVal: config.DscpMark, StateVal: state.DscpMark},
		},
	})
}
