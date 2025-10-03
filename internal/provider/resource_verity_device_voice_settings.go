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
	_ resource.Resource                = &verityDeviceVoiceSettingsResource{}
	_ resource.ResourceWithConfigure   = &verityDeviceVoiceSettingsResource{}
	_ resource.ResourceWithImportState = &verityDeviceVoiceSettingsResource{}
)

func NewVerityDeviceVoiceSettingsResource() resource.Resource {
	return &verityDeviceVoiceSettingsResource{}
}

type verityDeviceVoiceSettingsResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *utils.BulkOperationManager
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
	IsDefault types.Bool   `tfsdk:"isdefault"`
	Group     types.String `tfsdk:"group"`
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
			},
			"dtmf_method": schema.StringAttribute{
				Description: "Specifies how DTMF signals are carried",
				Optional:    true,
			},
			"region": schema.StringAttribute{
				Description: "Region",
				Optional:    true,
			},
			"protocol": schema.StringAttribute{
				Description: "Voice Protocol: MGCP or SIP",
				Optional:    true,
			},
			"proxy_server": schema.StringAttribute{
				Description: "IP address or URI of the SIP proxy server for SIP signalling messages",
				Optional:    true,
			},
			"proxy_server_port": schema.Int64Attribute{
				Description: "Proxy Server Port",
				Optional:    true,
			},
			"proxy_server_secondary": schema.StringAttribute{
				Description: "IP address or URI of the secondary SIP proxy server for SIP signalling messages",
				Optional:    true,
			},
			"proxy_server_secondary_port": schema.Int64Attribute{
				Description: "Secondary Proxy Server Port",
				Optional:    true,
			},
			"registrar_server": schema.StringAttribute{
				Description: "Name or IP address or resolved name of the registrar server for SIP signalling messages",
				Optional:    true,
			},
			"registrar_server_port": schema.Int64Attribute{
				Description: "Registrar Server Port",
				Optional:    true,
			},
			"registrar_server_secondary": schema.StringAttribute{
				Description: "Name or IP address or resolved name of the secondary registrar server for SIP signalling messages",
				Optional:    true,
			},
			"registrar_server_secondary_port": schema.Int64Attribute{
				Description: "Secondary Registrar Server Port",
				Optional:    true,
			},
			"user_agent_domain": schema.StringAttribute{
				Description: "User Agent Domain",
				Optional:    true,
			},
			"user_agent_transport": schema.StringAttribute{
				Description: "User Agent Transport",
				Optional:    true,
			},
			"user_agent_port": schema.Int64Attribute{
				Description: "User Agent Port",
				Optional:    true,
			},
			"outbound_proxy": schema.StringAttribute{
				Description: "IP address or URI of the outbound proxy server for SIP signalling messages",
				Optional:    true,
			},
			"outbound_proxy_port": schema.Int64Attribute{
				Description: "Outbound Proxy Port",
				Optional:    true,
			},
			"outbound_proxy_secondary": schema.StringAttribute{
				Description: "IP address or URI of the secondary outbound proxy server for SIP signalling messages",
				Optional:    true,
			},
			"outbound_proxy_secondary_port": schema.Int64Attribute{
				Description: "Secondary Outbound Proxy Port",
				Optional:    true,
			},
			"registration_period": schema.Int64Attribute{
				Description: "Specifies the time in seconds to start the re-registration process. The default value is 3240 seconds",
				Optional:    true,
			},
			"register_expires": schema.Int64Attribute{
				Description: "SIP registration expiration time in seconds. If value is 0, the SIP agent does not add an expiration time to the registration requests and does not perform re-registration. The default value is 3600 seconds",
				Optional:    true,
			},
			"voicemail_server": schema.StringAttribute{
				Description: "Name or IP address or resolved name of the external voicemail server if not provided by SIP server for MWI control",
				Optional:    true,
			},
			"voicemail_server_port": schema.Int64Attribute{
				Description: "Voicemail Server Port",
				Optional:    true,
			},
			"voicemail_server_expires": schema.Int64Attribute{
				Description: "Voicemail server expiration time in seconds. If value is 0, the Register Expires time is used instead. The default value is 3600 seconds",
				Optional:    true,
			},
			"sip_dscp_mark": schema.Int64Attribute{
				Description: "SIP Differentiated Services Code point (DSCP)",
				Optional:    true,
			},
			"call_agent_1": schema.StringAttribute{
				Description: "Call Agent 1",
				Optional:    true,
			},
			"call_agent_port_1": schema.Int64Attribute{
				Description: "Call Agent Port 1",
				Optional:    true,
			},
			"call_agent_2": schema.StringAttribute{
				Description: "Call Agent 2",
				Optional:    true,
			},
			"call_agent_port_2": schema.Int64Attribute{
				Description: "Call Agent Port 2",
				Optional:    true,
			},
			"domain": schema.StringAttribute{
				Description: "Domain",
				Optional:    true,
			},
			"mgcp_dscp_mark": schema.Int64Attribute{
				Description: "MGCP Differentiated Services Code point (DSCP)",
				Optional:    true,
			},
			"termination_base": schema.StringAttribute{
				Description: "Base string for the MGCP physical termination id(s)",
				Optional:    true,
			},
			"local_port_min": schema.Int64Attribute{
				Description: "Defines the base RTP port that should be used for voice traffic",
				Optional:    true,
			},
			"local_port_max": schema.Int64Attribute{
				Description: "Defines the highest RTP port used for voice traffic, must be greater than local Local Port Min",
				Optional:    true,
			},
			"event_payload_type": schema.Int64Attribute{
				Description: "Telephone Event Payload Type",
				Optional:    true,
			},
			"cas_events": schema.Int64Attribute{
				Description: "Enables or disables handling of CAS via RTP CAS events. Valid values are 0 = off and 1 = on",
				Optional:    true,
			},
			"dscp_mark": schema.Int64Attribute{
				Description: "Differentiated Services Code Point (DSCP) to be used for outgoing RTP packets",
				Optional:    true,
			},
			"rtcp": schema.BoolAttribute{
				Description: "RTCP Enable",
				Optional:    true,
			},
			"fax_t38": schema.BoolAttribute{
				Description: "Fax T.38 Enable",
				Optional:    true,
			},
			"bit_rate": schema.StringAttribute{
				Description: "T.38 Bit Rate in bps. Most available fax machines support up to 14,400bps",
				Optional:    true,
			},
			"cancel_call_waiting": schema.StringAttribute{
				Description: "Cancel Call waiting",
				Optional:    true,
			},
			"call_hold": schema.StringAttribute{
				Description: "Call hold",
				Optional:    true,
			},
			"cids_activate": schema.StringAttribute{
				Description: "Caller ID Delivery Blocking (single call) Activate",
				Optional:    true,
			},
			"cids_deactivate": schema.StringAttribute{
				Description: "Caller ID Delivery Blocking (single call) Deactivate",
				Optional:    true,
			},
			"do_not_disturb_activate": schema.StringAttribute{
				Description: "Do not Disturb Activate",
				Optional:    true,
			},
			"do_not_disturb_deactivate": schema.StringAttribute{
				Description: "Do not Disturb Deactivate",
				Optional:    true,
			},
			"do_not_disturb_pin_change": schema.StringAttribute{
				Description: "Do not Disturb PIN Change",
				Optional:    true,
			},
			"emergency_service_number": schema.StringAttribute{
				Description: "Emergency Service Number",
				Optional:    true,
			},
			"anon_cid_block_activate": schema.StringAttribute{
				Description: "Anonymous Caller ID Block Activate",
				Optional:    true,
			},
			"anon_cid_block_deactivate": schema.StringAttribute{
				Description: "Anonymous Caller ID Block Deactivate",
				Optional:    true,
			},
			"call_forward_unconditional_activate": schema.StringAttribute{
				Description: "Call Forward Unconditional Activate",
				Optional:    true,
			},
			"call_forward_unconditional_deactivate": schema.StringAttribute{
				Description: "Call Forward Unconditional Deactivate",
				Optional:    true,
			},
			"call_forward_on_busy_activate": schema.StringAttribute{
				Description: "Call Forward On Busy Activate",
				Optional:    true,
			},
			"call_forward_on_busy_deactivate": schema.StringAttribute{
				Description: "Call Forward On Busy Deactivate",
				Optional:    true,
			},
			"call_forward_on_no_answer_activate": schema.StringAttribute{
				Description: "Call Forward On No Answer Activate",
				Optional:    true,
			},
			"call_forward_on_no_answer_deactivate": schema.StringAttribute{
				Description: "Call Forward On No Answer Deactivate",
				Optional:    true,
			},
			"intercom_1": schema.StringAttribute{
				Description: "Intercom 1",
				Optional:    true,
			},
			"intercom_2": schema.StringAttribute{
				Description: "Intercom 2",
				Optional:    true,
			},
			"intercom_3": schema.StringAttribute{
				Description: "Intercom 3",
				Optional:    true,
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
						},
						"codec_num_enable": schema.BoolAttribute{
							Description: "Enable Codec",
							Optional:    true,
						},
						"codec_num_packetization_period": schema.StringAttribute{
							Description: "Packet period selection interval in milliseconds",
							Optional:    true,
						},
						"codec_num_silence_suppression": schema.BoolAttribute{
							Description: "Specifies whether silence suppression is on or off. Valid values are 0 = off and 1 = on",
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
				Description: "Object properties for the device voice settings",
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

func (r *verityDeviceVoiceSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityDeviceVoiceSettingsResourceModel
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

	// Handle int64 fields
	utils.SetInt64Fields([]utils.Int64FieldMapping{
		{FieldName: "CasEvents", APIField: &dvsProps.CasEvents, TFValue: plan.CasEvents},
	})

	// Handle nullable int64 fields
	utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
		{FieldName: "ProxyServerPort", APIField: &dvsProps.ProxyServerPort, TFValue: plan.ProxyServerPort},
		{FieldName: "ProxyServerSecondaryPort", APIField: &dvsProps.ProxyServerSecondaryPort, TFValue: plan.ProxyServerSecondaryPort},
		{FieldName: "RegistrarServerPort", APIField: &dvsProps.RegistrarServerPort, TFValue: plan.RegistrarServerPort},
		{FieldName: "RegistrarServerSecondaryPort", APIField: &dvsProps.RegistrarServerSecondaryPort, TFValue: plan.RegistrarServerSecondaryPort},
		{FieldName: "UserAgentPort", APIField: &dvsProps.UserAgentPort, TFValue: plan.UserAgentPort},
		{FieldName: "OutboundProxyPort", APIField: &dvsProps.OutboundProxyPort, TFValue: plan.OutboundProxyPort},
		{FieldName: "OutboundProxySecondaryPort", APIField: &dvsProps.OutboundProxySecondaryPort, TFValue: plan.OutboundProxySecondaryPort},
		{FieldName: "RegistrationPeriod", APIField: &dvsProps.RegistrationPeriod, TFValue: plan.RegistrationPeriod},
		{FieldName: "RegisterExpires", APIField: &dvsProps.RegisterExpires, TFValue: plan.RegisterExpires},
		{FieldName: "VoicemailServerPort", APIField: &dvsProps.VoicemailServerPort, TFValue: plan.VoicemailServerPort},
		{FieldName: "VoicemailServerExpires", APIField: &dvsProps.VoicemailServerExpires, TFValue: plan.VoicemailServerExpires},
		{FieldName: "SipDscpMark", APIField: &dvsProps.SipDscpMark, TFValue: plan.SipDscpMark},
		{FieldName: "CallAgentPort1", APIField: &dvsProps.CallAgentPort1, TFValue: plan.CallAgentPort1},
		{FieldName: "CallAgentPort2", APIField: &dvsProps.CallAgentPort2, TFValue: plan.CallAgentPort2},
		{FieldName: "MgcpDscpMark", APIField: &dvsProps.MgcpDscpMark, TFValue: plan.MgcpDscpMark},
		{FieldName: "LocalPortMin", APIField: &dvsProps.LocalPortMin, TFValue: plan.LocalPortMin},
		{FieldName: "LocalPortMax", APIField: &dvsProps.LocalPortMax, TFValue: plan.LocalPortMax},
		{FieldName: "EventPayloadType", APIField: &dvsProps.EventPayloadType, TFValue: plan.EventPayloadType},
		{FieldName: "DscpMark", APIField: &dvsProps.DscpMark, TFValue: plan.DscpMark},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.PacketqueuesPutRequestPacketQueueValueObjectProperties{}
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
		dvsProps.ObjectProperties = &objProps
	}

	// Handle codecs
	if len(plan.Codecs) > 0 {
		codecs := make([]openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner, len(plan.Codecs))
		for i, codec := range plan.Codecs {
			codecItem := openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner{}
			if !codec.CodecNumName.IsNull() {
				codecItem.CodecNumName = openapi.PtrString(codec.CodecNumName.ValueString())
			}
			if !codec.CodecNumEnable.IsNull() {
				codecItem.CodecNumEnable = openapi.PtrBool(codec.CodecNumEnable.ValueBool())
			}
			if !codec.CodecNumPacketizationPeriod.IsNull() {
				codecItem.CodecNumPacketizationPeriod = openapi.PtrString(codec.CodecNumPacketizationPeriod.ValueString())
			}
			if !codec.CodecNumSilenceSuppression.IsNull() {
				codecItem.CodecNumSilenceSuppression = openapi.PtrBool(codec.CodecNumSilenceSuppression.ValueBool())
			}
			if !codec.Index.IsNull() {
				codecItem.Index = openapi.PtrInt32(int32(codec.Index.ValueInt64()))
			}
			codecs[i] = codecItem
		}
		dvsProps.Codecs = codecs
	}

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "device_voice_settings", name, *dvsProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Device Voice Settings %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "device_voice_settings")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
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

	state.Name = utils.MapStringFromAPI(dvsMap["name"])

	// Handle object properties
	if objProps, ok := dvsMap["object_properties"].(map[string]interface{}); ok {
		group := utils.MapStringFromAPI(objProps["group"])
		if group.IsNull() {
			group = types.StringValue("")
		}
		isdefault := utils.MapBoolFromAPI(objProps["isdefault"])
		state.ObjectProperties = []verityDeviceVoiceSettingsObjectPropertiesModel{
			{Group: group, IsDefault: isdefault},
		}
	} else {
		state.ObjectProperties = nil
	}

	// Map string fields
	stringFieldMappings := map[string]*types.String{
		"dtmf_method":                           &state.DtmfMethod,
		"region":                                &state.Region,
		"protocol":                              &state.Protocol,
		"proxy_server":                          &state.ProxyServer,
		"proxy_server_secondary":                &state.ProxyServerSecondary,
		"registrar_server":                      &state.RegistrarServer,
		"registrar_server_secondary":            &state.RegistrarServerSecondary,
		"user_agent_domain":                     &state.UserAgentDomain,
		"user_agent_transport":                  &state.UserAgentTransport,
		"outbound_proxy":                        &state.OutboundProxy,
		"outbound_proxy_secondary":              &state.OutboundProxySecondary,
		"voicemail_server":                      &state.VoicemailServer,
		"call_agent_1":                          &state.CallAgent1,
		"call_agent_2":                          &state.CallAgent2,
		"domain":                                &state.Domain,
		"termination_base":                      &state.TerminationBase,
		"bit_rate":                              &state.BitRate,
		"cancel_call_waiting":                   &state.CancelCallWaiting,
		"call_hold":                             &state.CallHold,
		"cids_activate":                         &state.CidsActivate,
		"cids_deactivate":                       &state.CidsDeactivate,
		"do_not_disturb_activate":               &state.DoNotDisturbActivate,
		"do_not_disturb_deactivate":             &state.DoNotDisturbDeactivate,
		"do_not_disturb_pin_change":             &state.DoNotDisturbPinChange,
		"emergency_service_number":              &state.EmergencyServiceNumber,
		"anon_cid_block_activate":               &state.AnonCidBlockActivate,
		"anon_cid_block_deactivate":             &state.AnonCidBlockDeactivate,
		"call_forward_unconditional_activate":   &state.CallForwardUnconditionalActivate,
		"call_forward_unconditional_deactivate": &state.CallForwardUnconditionalDeactivate,
		"call_forward_on_busy_activate":         &state.CallForwardOnBusyActivate,
		"call_forward_on_busy_deactivate":       &state.CallForwardOnBusyDeactivate,
		"call_forward_on_no_answer_activate":    &state.CallForwardOnNoAnswerActivate,
		"call_forward_on_no_answer_deactivate":  &state.CallForwardOnNoAnswerDeactivate,
		"intercom_1":                            &state.Intercom1,
		"intercom_2":                            &state.Intercom2,
		"intercom_3":                            &state.Intercom3,
	}

	for apiKey, stateField := range stringFieldMappings {
		*stateField = utils.MapStringFromAPI(dvsMap[apiKey])
	}

	// Map boolean fields
	boolFieldMappings := map[string]*types.Bool{
		"enable":  &state.Enable,
		"rtcp":    &state.Rtcp,
		"fax_t38": &state.FaxT38,
	}

	for apiKey, stateField := range boolFieldMappings {
		*stateField = utils.MapBoolFromAPI(dvsMap[apiKey])
	}

	// Map nullable int64 fields
	nullableInt64FieldMappings := map[string]*types.Int64{
		"proxy_server_port":               &state.ProxyServerPort,
		"proxy_server_secondary_port":     &state.ProxyServerSecondaryPort,
		"registrar_server_port":           &state.RegistrarServerPort,
		"registrar_server_secondary_port": &state.RegistrarServerSecondaryPort,
		"user_agent_port":                 &state.UserAgentPort,
		"outbound_proxy_port":             &state.OutboundProxyPort,
		"outbound_proxy_secondary_port":   &state.OutboundProxySecondaryPort,
		"registration_period":             &state.RegistrationPeriod,
		"register_expires":                &state.RegisterExpires,
		"voicemail_server_port":           &state.VoicemailServerPort,
		"voicemail_server_expires":        &state.VoicemailServerExpires,
		"sip_dscp_mark":                   &state.SipDscpMark,
		"call_agent_port_1":               &state.CallAgentPort1,
		"call_agent_port_2":               &state.CallAgentPort2,
		"mgcp_dscp_mark":                  &state.MgcpDscpMark,
		"local_port_min":                  &state.LocalPortMin,
		"local_port_max":                  &state.LocalPortMax,
		"event_payload_type":              &state.EventPayloadType,
		"dscp_mark":                       &state.DscpMark,
		"cas_events":                      &state.CasEvents,
	}

	for apiKey, stateField := range nullableInt64FieldMappings {
		*stateField = utils.MapInt64FromAPI(dvsMap[apiKey])
	}

	// Handle codecs
	if codecsArray, ok := dvsMap["codecs"].([]interface{}); ok && len(codecsArray) > 0 {
		var codecs []verityDeviceVoiceSettingsCodecModel

		for _, c := range codecsArray {
			codec, ok := c.(map[string]interface{})
			if !ok {
				continue
			}

			codecModel := verityDeviceVoiceSettingsCodecModel{
				CodecNumName:                utils.MapStringFromAPI(codec["codec_num_name"]),
				CodecNumEnable:              utils.MapBoolFromAPI(codec["codec_num_enable"]),
				CodecNumPacketizationPeriod: utils.MapStringFromAPI(codec["codec_num_packetization_period"]),
				CodecNumSilenceSuppression:  utils.MapBoolFromAPI(codec["codec_num_silence_suppression"]),
				Index:                       utils.MapInt64FromAPI(codec["index"]),
			}

			codecs = append(codecs, codecModel)
		}

		state.Codecs = codecs
	} else {
		state.Codecs = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
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

	// Handle nullable int64 field changes
	utils.CompareAndSetNullableInt64Field(plan.ProxyServerPort, state.ProxyServerPort, func(v *openapi.NullableInt32) { dvsProps.ProxyServerPort = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.ProxyServerSecondaryPort, state.ProxyServerSecondaryPort, func(v *openapi.NullableInt32) { dvsProps.ProxyServerSecondaryPort = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.RegistrarServerPort, state.RegistrarServerPort, func(v *openapi.NullableInt32) { dvsProps.RegistrarServerPort = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.RegistrarServerSecondaryPort, state.RegistrarServerSecondaryPort, func(v *openapi.NullableInt32) { dvsProps.RegistrarServerSecondaryPort = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.UserAgentPort, state.UserAgentPort, func(v *openapi.NullableInt32) { dvsProps.UserAgentPort = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.OutboundProxyPort, state.OutboundProxyPort, func(v *openapi.NullableInt32) { dvsProps.OutboundProxyPort = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.OutboundProxySecondaryPort, state.OutboundProxySecondaryPort, func(v *openapi.NullableInt32) { dvsProps.OutboundProxySecondaryPort = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.RegistrationPeriod, state.RegistrationPeriod, func(v *openapi.NullableInt32) { dvsProps.RegistrationPeriod = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.RegisterExpires, state.RegisterExpires, func(v *openapi.NullableInt32) { dvsProps.RegisterExpires = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.VoicemailServerPort, state.VoicemailServerPort, func(v *openapi.NullableInt32) { dvsProps.VoicemailServerPort = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.VoicemailServerExpires, state.VoicemailServerExpires, func(v *openapi.NullableInt32) { dvsProps.VoicemailServerExpires = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.SipDscpMark, state.SipDscpMark, func(v *openapi.NullableInt32) { dvsProps.SipDscpMark = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.CallAgentPort1, state.CallAgentPort1, func(v *openapi.NullableInt32) { dvsProps.CallAgentPort1 = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.CallAgentPort2, state.CallAgentPort2, func(v *openapi.NullableInt32) { dvsProps.CallAgentPort2 = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.MgcpDscpMark, state.MgcpDscpMark, func(v *openapi.NullableInt32) { dvsProps.MgcpDscpMark = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.LocalPortMin, state.LocalPortMin, func(v *openapi.NullableInt32) { dvsProps.LocalPortMin = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.LocalPortMax, state.LocalPortMax, func(v *openapi.NullableInt32) { dvsProps.LocalPortMax = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.EventPayloadType, state.EventPayloadType, func(v *openapi.NullableInt32) { dvsProps.EventPayloadType = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.DscpMark, state.DscpMark, func(v *openapi.NullableInt32) { dvsProps.DscpMark = *v }, &hasChanges)

	// Handle non-nullable int64 field
	utils.CompareAndSetInt64Field(plan.CasEvents, state.CasEvents, func(v *int32) { dvsProps.CasEvents = v }, &hasChanges)

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		if len(state.ObjectProperties) == 0 ||
			!plan.ObjectProperties[0].IsDefault.Equal(state.ObjectProperties[0].IsDefault) ||
			!plan.ObjectProperties[0].Group.Equal(state.ObjectProperties[0].Group) {
			op := plan.ObjectProperties[0]
			objProps := openapi.PacketqueuesPutRequestPacketQueueValueObjectProperties{}
			if !op.IsDefault.IsNull() {
				objProps.Isdefault = openapi.PtrBool(op.IsDefault.ValueBool())
			}
			if !op.Group.IsNull() {
				objProps.Group = openapi.PtrString(op.Group.ValueString())
			}
			dvsProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	// Handle codecs
	codecsHandler := utils.IndexedItemHandler[verityDeviceVoiceSettingsCodecModel, openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner]{
		CreateNew: func(planItem verityDeviceVoiceSettingsCodecModel) openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner {
			codec := openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner{
				Index: openapi.PtrInt32(int32(planItem.Index.ValueInt64())),
			}

			if !planItem.CodecNumName.IsNull() {
				codec.CodecNumName = openapi.PtrString(planItem.CodecNumName.ValueString())
			} else {
				codec.CodecNumName = openapi.PtrString("")
			}

			if !planItem.CodecNumEnable.IsNull() {
				codec.CodecNumEnable = openapi.PtrBool(planItem.CodecNumEnable.ValueBool())
			} else {
				codec.CodecNumEnable = openapi.PtrBool(false)
			}

			if !planItem.CodecNumPacketizationPeriod.IsNull() {
				codec.CodecNumPacketizationPeriod = openapi.PtrString(planItem.CodecNumPacketizationPeriod.ValueString())
			} else {
				codec.CodecNumPacketizationPeriod = openapi.PtrString("")
			}

			if !planItem.CodecNumSilenceSuppression.IsNull() {
				codec.CodecNumSilenceSuppression = openapi.PtrBool(planItem.CodecNumSilenceSuppression.ValueBool())
			} else {
				codec.CodecNumSilenceSuppression = openapi.PtrBool(false)
			}

			return codec
		},
		UpdateExisting: func(planItem verityDeviceVoiceSettingsCodecModel, stateItem verityDeviceVoiceSettingsCodecModel) (openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner, bool) {
			codec := openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner{
				Index: openapi.PtrInt32(int32(planItem.Index.ValueInt64())),
			}

			fieldChanged := false

			if !planItem.CodecNumName.Equal(stateItem.CodecNumName) {
				if !planItem.CodecNumName.IsNull() {
					codec.CodecNumName = openapi.PtrString(planItem.CodecNumName.ValueString())
				} else {
					codec.CodecNumName = openapi.PtrString("")
				}
				fieldChanged = true
			}

			if !planItem.CodecNumEnable.Equal(stateItem.CodecNumEnable) {
				codec.CodecNumEnable = openapi.PtrBool(planItem.CodecNumEnable.ValueBool())
				fieldChanged = true
			}

			if !planItem.CodecNumPacketizationPeriod.Equal(stateItem.CodecNumPacketizationPeriod) {
				if !planItem.CodecNumPacketizationPeriod.IsNull() {
					codec.CodecNumPacketizationPeriod = openapi.PtrString(planItem.CodecNumPacketizationPeriod.ValueString())
				} else {
					codec.CodecNumPacketizationPeriod = openapi.PtrString("")
				}
				fieldChanged = true
			}

			if !planItem.CodecNumSilenceSuppression.Equal(stateItem.CodecNumSilenceSuppression) {
				codec.CodecNumSilenceSuppression = openapi.PtrBool(planItem.CodecNumSilenceSuppression.ValueBool())
				fieldChanged = true
			}

			return codec, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner {
			return openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner{
				Index: openapi.PtrInt32(int32(index)),
			}
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

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "device_voice_settings", name, dvsProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Device Voice Settings %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "device_voice_settings")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
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

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "device_voice_settings", name, nil, &resp.Diagnostics)
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
