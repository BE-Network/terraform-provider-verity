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

	if !plan.Enable.IsNull() {
		dvsProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
	}
	if !plan.DtmfMethod.IsNull() {
		dvsProps.DtmfMethod = openapi.PtrString(plan.DtmfMethod.ValueString())
	}
	if !plan.Region.IsNull() {
		dvsProps.Region = openapi.PtrString(plan.Region.ValueString())
	}
	if !plan.Protocol.IsNull() {
		dvsProps.Protocol = openapi.PtrString(plan.Protocol.ValueString())
	}

	if !plan.ProxyServer.IsNull() {
		dvsProps.ProxyServer = openapi.PtrString(plan.ProxyServer.ValueString())
	}
	if !plan.ProxyServerSecondary.IsNull() {
		dvsProps.ProxyServerSecondary = openapi.PtrString(plan.ProxyServerSecondary.ValueString())
	}
	if !plan.RegistrarServer.IsNull() {
		dvsProps.RegistrarServer = openapi.PtrString(plan.RegistrarServer.ValueString())
	}
	if !plan.RegistrarServerSecondary.IsNull() {
		dvsProps.RegistrarServerSecondary = openapi.PtrString(plan.RegistrarServerSecondary.ValueString())
	}
	if !plan.UserAgentDomain.IsNull() {
		dvsProps.UserAgentDomain = openapi.PtrString(plan.UserAgentDomain.ValueString())
	}
	if !plan.UserAgentTransport.IsNull() {
		dvsProps.UserAgentTransport = openapi.PtrString(plan.UserAgentTransport.ValueString())
	}
	if !plan.OutboundProxy.IsNull() {
		dvsProps.OutboundProxy = openapi.PtrString(plan.OutboundProxy.ValueString())
	}
	if !plan.OutboundProxySecondary.IsNull() {
		dvsProps.OutboundProxySecondary = openapi.PtrString(plan.OutboundProxySecondary.ValueString())
	}
	if !plan.VoicemailServer.IsNull() {
		dvsProps.VoicemailServer = openapi.PtrString(plan.VoicemailServer.ValueString())
	}
	if !plan.CallAgent1.IsNull() {
		dvsProps.CallAgent1 = openapi.PtrString(plan.CallAgent1.ValueString())
	}
	if !plan.CallAgent2.IsNull() {
		dvsProps.CallAgent2 = openapi.PtrString(plan.CallAgent2.ValueString())
	}
	if !plan.Domain.IsNull() {
		dvsProps.Domain = openapi.PtrString(plan.Domain.ValueString())
	}
	if !plan.TerminationBase.IsNull() {
		dvsProps.TerminationBase = openapi.PtrString(plan.TerminationBase.ValueString())
	}
	if !plan.BitRate.IsNull() {
		dvsProps.BitRate = openapi.PtrString(plan.BitRate.ValueString())
	}
	if !plan.CancelCallWaiting.IsNull() {
		dvsProps.CancelCallWaiting = openapi.PtrString(plan.CancelCallWaiting.ValueString())
	}
	if !plan.CallHold.IsNull() {
		dvsProps.CallHold = openapi.PtrString(plan.CallHold.ValueString())
	}
	if !plan.CidsActivate.IsNull() {
		dvsProps.CidsActivate = openapi.PtrString(plan.CidsActivate.ValueString())
	}
	if !plan.CidsDeactivate.IsNull() {
		dvsProps.CidsDeactivate = openapi.PtrString(plan.CidsDeactivate.ValueString())
	}
	if !plan.DoNotDisturbActivate.IsNull() {
		dvsProps.DoNotDisturbActivate = openapi.PtrString(plan.DoNotDisturbActivate.ValueString())
	}
	if !plan.DoNotDisturbDeactivate.IsNull() {
		dvsProps.DoNotDisturbDeactivate = openapi.PtrString(plan.DoNotDisturbDeactivate.ValueString())
	}
	if !plan.DoNotDisturbPinChange.IsNull() {
		dvsProps.DoNotDisturbPinChange = openapi.PtrString(plan.DoNotDisturbPinChange.ValueString())
	}
	if !plan.EmergencyServiceNumber.IsNull() {
		dvsProps.EmergencyServiceNumber = openapi.PtrString(plan.EmergencyServiceNumber.ValueString())
	}
	if !plan.AnonCidBlockActivate.IsNull() {
		dvsProps.AnonCidBlockActivate = openapi.PtrString(plan.AnonCidBlockActivate.ValueString())
	}
	if !plan.AnonCidBlockDeactivate.IsNull() {
		dvsProps.AnonCidBlockDeactivate = openapi.PtrString(plan.AnonCidBlockDeactivate.ValueString())
	}
	if !plan.CallForwardUnconditionalActivate.IsNull() {
		dvsProps.CallForwardUnconditionalActivate = openapi.PtrString(plan.CallForwardUnconditionalActivate.ValueString())
	}
	if !plan.CallForwardUnconditionalDeactivate.IsNull() {
		dvsProps.CallForwardUnconditionalDeactivate = openapi.PtrString(plan.CallForwardUnconditionalDeactivate.ValueString())
	}
	if !plan.CallForwardOnBusyActivate.IsNull() {
		dvsProps.CallForwardOnBusyActivate = openapi.PtrString(plan.CallForwardOnBusyActivate.ValueString())
	}
	if !plan.CallForwardOnBusyDeactivate.IsNull() {
		dvsProps.CallForwardOnBusyDeactivate = openapi.PtrString(plan.CallForwardOnBusyDeactivate.ValueString())
	}
	if !plan.CallForwardOnNoAnswerActivate.IsNull() {
		dvsProps.CallForwardOnNoAnswerActivate = openapi.PtrString(plan.CallForwardOnNoAnswerActivate.ValueString())
	}
	if !plan.CallForwardOnNoAnswerDeactivate.IsNull() {
		dvsProps.CallForwardOnNoAnswerDeactivate = openapi.PtrString(plan.CallForwardOnNoAnswerDeactivate.ValueString())
	}
	if !plan.Intercom1.IsNull() {
		dvsProps.Intercom1 = openapi.PtrString(plan.Intercom1.ValueString())
	}
	if !plan.Intercom2.IsNull() {
		dvsProps.Intercom2 = openapi.PtrString(plan.Intercom2.ValueString())
	}
	if !plan.Intercom3.IsNull() {
		dvsProps.Intercom3 = openapi.PtrString(plan.Intercom3.ValueString())
	}

	if !plan.ProxyServerPort.IsNull() {
		val := int32(plan.ProxyServerPort.ValueInt64())
		dvsProps.ProxyServerPort = *openapi.NewNullableInt32(&val)
	} else {
		dvsProps.ProxyServerPort = *openapi.NewNullableInt32(nil)
	}
	if !plan.ProxyServerSecondaryPort.IsNull() {
		val := int32(plan.ProxyServerSecondaryPort.ValueInt64())
		dvsProps.ProxyServerSecondaryPort = *openapi.NewNullableInt32(&val)
	} else {
		dvsProps.ProxyServerSecondaryPort = *openapi.NewNullableInt32(nil)
	}
	if !plan.RegistrarServerPort.IsNull() {
		val := int32(plan.RegistrarServerPort.ValueInt64())
		dvsProps.RegistrarServerPort = *openapi.NewNullableInt32(&val)
	} else {
		dvsProps.RegistrarServerPort = *openapi.NewNullableInt32(nil)
	}
	if !plan.RegistrarServerSecondaryPort.IsNull() {
		val := int32(plan.RegistrarServerSecondaryPort.ValueInt64())
		dvsProps.RegistrarServerSecondaryPort = *openapi.NewNullableInt32(&val)
	} else {
		dvsProps.RegistrarServerSecondaryPort = *openapi.NewNullableInt32(nil)
	}
	if !plan.UserAgentPort.IsNull() {
		val := int32(plan.UserAgentPort.ValueInt64())
		dvsProps.UserAgentPort = *openapi.NewNullableInt32(&val)
	} else {
		dvsProps.UserAgentPort = *openapi.NewNullableInt32(nil)
	}
	if !plan.OutboundProxyPort.IsNull() {
		val := int32(plan.OutboundProxyPort.ValueInt64())
		dvsProps.OutboundProxyPort = *openapi.NewNullableInt32(&val)
	} else {
		dvsProps.OutboundProxyPort = *openapi.NewNullableInt32(nil)
	}
	if !plan.OutboundProxySecondaryPort.IsNull() {
		val := int32(plan.OutboundProxySecondaryPort.ValueInt64())
		dvsProps.OutboundProxySecondaryPort = *openapi.NewNullableInt32(&val)
	} else {
		dvsProps.OutboundProxySecondaryPort = *openapi.NewNullableInt32(nil)
	}
	if !plan.RegistrationPeriod.IsNull() {
		val := int32(plan.RegistrationPeriod.ValueInt64())
		dvsProps.RegistrationPeriod = *openapi.NewNullableInt32(&val)
	} else {
		dvsProps.RegistrationPeriod = *openapi.NewNullableInt32(nil)
	}
	if !plan.RegisterExpires.IsNull() {
		val := int32(plan.RegisterExpires.ValueInt64())
		dvsProps.RegisterExpires = *openapi.NewNullableInt32(&val)
	} else {
		dvsProps.RegisterExpires = *openapi.NewNullableInt32(nil)
	}
	if !plan.VoicemailServerPort.IsNull() {
		val := int32(plan.VoicemailServerPort.ValueInt64())
		dvsProps.VoicemailServerPort = *openapi.NewNullableInt32(&val)
	} else {
		dvsProps.VoicemailServerPort = *openapi.NewNullableInt32(nil)
	}
	if !plan.VoicemailServerExpires.IsNull() {
		val := int32(plan.VoicemailServerExpires.ValueInt64())
		dvsProps.VoicemailServerExpires = *openapi.NewNullableInt32(&val)
	} else {
		dvsProps.VoicemailServerExpires = *openapi.NewNullableInt32(nil)
	}
	if !plan.SipDscpMark.IsNull() {
		val := int32(plan.SipDscpMark.ValueInt64())
		dvsProps.SipDscpMark = *openapi.NewNullableInt32(&val)
	} else {
		dvsProps.SipDscpMark = *openapi.NewNullableInt32(nil)
	}
	if !plan.CallAgentPort1.IsNull() {
		val := int32(plan.CallAgentPort1.ValueInt64())
		dvsProps.CallAgentPort1 = *openapi.NewNullableInt32(&val)
	} else {
		dvsProps.CallAgentPort1 = *openapi.NewNullableInt32(nil)
	}
	if !plan.CallAgentPort2.IsNull() {
		val := int32(plan.CallAgentPort2.ValueInt64())
		dvsProps.CallAgentPort2 = *openapi.NewNullableInt32(&val)
	} else {
		dvsProps.CallAgentPort2 = *openapi.NewNullableInt32(nil)
	}
	if !plan.MgcpDscpMark.IsNull() {
		val := int32(plan.MgcpDscpMark.ValueInt64())
		dvsProps.MgcpDscpMark = *openapi.NewNullableInt32(&val)
	} else {
		dvsProps.MgcpDscpMark = *openapi.NewNullableInt32(nil)
	}
	if !plan.LocalPortMin.IsNull() {
		val := int32(plan.LocalPortMin.ValueInt64())
		dvsProps.LocalPortMin = *openapi.NewNullableInt32(&val)
	} else {
		dvsProps.LocalPortMin = *openapi.NewNullableInt32(nil)
	}
	if !plan.LocalPortMax.IsNull() {
		val := int32(plan.LocalPortMax.ValueInt64())
		dvsProps.LocalPortMax = *openapi.NewNullableInt32(&val)
	} else {
		dvsProps.LocalPortMax = *openapi.NewNullableInt32(nil)
	}
	if !plan.EventPayloadType.IsNull() {
		val := int32(plan.EventPayloadType.ValueInt64())
		dvsProps.EventPayloadType = *openapi.NewNullableInt32(&val)
	} else {
		dvsProps.EventPayloadType = *openapi.NewNullableInt32(nil)
	}
	if !plan.CasEvents.IsNull() {
		dvsProps.CasEvents = openapi.PtrInt32(int32(plan.CasEvents.ValueInt64()))
	}
	if !plan.DscpMark.IsNull() {
		val := int32(plan.DscpMark.ValueInt64())
		dvsProps.DscpMark = *openapi.NewNullableInt32(&val)
	} else {
		dvsProps.DscpMark = *openapi.NewNullableInt32(nil)
	}

	if !plan.Rtcp.IsNull() {
		dvsProps.Rtcp = openapi.PtrBool(plan.Rtcp.ValueBool())
	}
	if !plan.FaxT38.IsNull() {
		dvsProps.FaxT38 = openapi.PtrBool(plan.FaxT38.ValueBool())
	}

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

	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.PacketqueuesPutRequestPacketQueueValueObjectProperties{}
		if !op.IsDefault.IsNull() {
			objProps.Isdefault = openapi.PtrBool(op.IsDefault.ValueBool())
		}
		if !op.Group.IsNull() {
			objProps.Group = openapi.PtrString(op.Group.ValueString())
		}
		dvsProps.ObjectProperties = &objProps
	}

	operationID := r.bulkOpsMgr.AddPut(ctx, "device_voice_settings", name, *dvsProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for device voice settings creation operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Create Device Voice Settings %s", name))...,
		)
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

	var result DeviceVoiceSettingsResponse
	var err error
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		dvsData, fetchErr := getCachedResponse(ctx, r.provCtx, "device_voice_settings", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch Device Voice Settings")
			respAPI, err := r.client.DeviceVoiceSettingsAPI.DevicevoicesettingsGet(ctx).Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading Device Voice Settings: %v", err)
			}
			defer respAPI.Body.Close()

			var res DeviceVoiceSettingsResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return nil, fmt.Errorf("failed to decode Device Voice Settings response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d Device Voice Settings", len(res.DeviceVoiceSettings)))
			return res, nil
		})
		if fetchErr != nil {
			err = fetchErr
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch Device Voice Settings on attempt %d, retrying in %v", attempt+1, sleepTime))
			time.Sleep(sleepTime)
			continue
		}
		result = dvsData.(DeviceVoiceSettingsResponse)
		break
	}
	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Device Voice Settings %s", dvsName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for Device Voice Settings with ID: %s", dvsName))
	var dvsData map[string]interface{}
	exists := false

	if data, ok := result.DeviceVoiceSettings[dvsName].(map[string]interface{}); ok {
		dvsData = data
		exists = true
		tflog.Debug(ctx, fmt.Sprintf("Found Device Voice Settings directly by ID: %s", dvsName))
	} else {
		for apiName, d := range result.DeviceVoiceSettings {
			deviceVoice, ok := d.(map[string]interface{})
			if !ok {
				continue
			}

			if name, ok := deviceVoice["name"].(string); ok && name == dvsName {
				dvsData = deviceVoice
				dvsName = apiName
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found Device Voice Settings with name '%s' under API key '%s'", name, apiName))
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Device Voice Settings with ID '%s' not found in API response", dvsName))
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(fmt.Sprintf("%v", dvsData["name"]))

	if enable, ok := dvsData["enable"].(bool); ok {
		state.Enable = types.BoolValue(enable)
	} else {
		state.Enable = types.BoolNull()
	}

	if rtcp, ok := dvsData["rtcp"].(bool); ok {
		state.Rtcp = types.BoolValue(rtcp)
	} else {
		state.Rtcp = types.BoolNull()
	}

	if faxT38, ok := dvsData["fax_t38"].(bool); ok {
		state.FaxT38 = types.BoolValue(faxT38)
	} else {
		state.FaxT38 = types.BoolNull()
	}

	stringFields := map[string]*types.String{
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

	for apiKey, stateField := range stringFields {
		if value, ok := dvsData[apiKey].(string); ok {
			*stateField = types.StringValue(value)
		} else {
			*stateField = types.StringNull()
		}
	}

	intFields := map[string]*types.Int64{
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
		"cas_events":                      &state.CasEvents,
		"dscp_mark":                       &state.DscpMark,
	}

	for apiKey, stateField := range intFields {
		if value, ok := dvsData[apiKey]; ok && value != nil {
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

	if codecsArray, ok := dvsData["Codecs"].([]interface{}); ok && len(codecsArray) > 0 {
		var codecs []verityDeviceVoiceSettingsCodecModel
		for _, c := range codecsArray {
			codec, ok := c.(map[string]interface{})
			if !ok {
				continue
			}
			codecModel := verityDeviceVoiceSettingsCodecModel{}

			if name, ok := codec["codec_num_name"].(string); ok {
				codecModel.CodecNumName = types.StringValue(name)
			} else {
				codecModel.CodecNumName = types.StringNull()
			}

			if enable, ok := codec["codec_num_enable"].(bool); ok {
				codecModel.CodecNumEnable = types.BoolValue(enable)
			} else {
				codecModel.CodecNumEnable = types.BoolNull()
			}

			if period, ok := codec["codec_num_packetization_period"].(string); ok {
				codecModel.CodecNumPacketizationPeriod = types.StringValue(period)
			} else {
				codecModel.CodecNumPacketizationPeriod = types.StringNull()
			}

			if suppression, ok := codec["codec_num_silence_suppression"].(bool); ok {
				codecModel.CodecNumSilenceSuppression = types.BoolValue(suppression)
			} else {
				codecModel.CodecNumSilenceSuppression = types.BoolNull()
			}

			if index, ok := codec["index"]; ok && index != nil {
				if intVal, ok := index.(float64); ok {
					codecModel.Index = types.Int64Value(int64(intVal))
				} else if intVal, ok := index.(int); ok {
					codecModel.Index = types.Int64Value(int64(intVal))
				} else {
					codecModel.Index = types.Int64Null()
				}
			} else {
				codecModel.Index = types.Int64Null()
			}

			codecs = append(codecs, codecModel)
		}
		state.Codecs = codecs
	} else {
		state.Codecs = nil
	}

	if objProps, ok := dvsData["object_properties"].(map[string]interface{}); ok {
		op := verityDeviceVoiceSettingsObjectPropertiesModel{}
		if isDefault, ok := objProps["isdefault"].(bool); ok {
			op.IsDefault = types.BoolValue(isDefault)
		} else {
			op.IsDefault = types.BoolNull()
		}
		if group, ok := objProps["group"].(string); ok {
			op.Group = types.StringValue(group)
		} else {
			op.Group = types.StringNull()
		}
		state.ObjectProperties = []verityDeviceVoiceSettingsObjectPropertiesModel{op}
	} else {
		state.ObjectProperties = nil
	}

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

	if !plan.Name.Equal(state.Name) {
		dvsProps.Name = openapi.PtrString(name)
		hasChanges = true
	}

	if !plan.Enable.Equal(state.Enable) {
		dvsProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
		hasChanges = true
	}
	if !plan.DtmfMethod.Equal(state.DtmfMethod) {
		dvsProps.DtmfMethod = openapi.PtrString(plan.DtmfMethod.ValueString())
		hasChanges = true
	}
	if !plan.Region.Equal(state.Region) {
		dvsProps.Region = openapi.PtrString(plan.Region.ValueString())
		hasChanges = true
	}
	if !plan.Protocol.Equal(state.Protocol) {
		dvsProps.Protocol = openapi.PtrString(plan.Protocol.ValueString())
		hasChanges = true
	}

	stringFields := map[string]struct {
		planField  types.String
		stateField types.String
		setter     func(string)
	}{
		"proxy_server":                          {plan.ProxyServer, state.ProxyServer, func(v string) { dvsProps.ProxyServer = openapi.PtrString(v) }},
		"proxy_server_secondary":                {plan.ProxyServerSecondary, state.ProxyServerSecondary, func(v string) { dvsProps.ProxyServerSecondary = openapi.PtrString(v) }},
		"registrar_server":                      {plan.RegistrarServer, state.RegistrarServer, func(v string) { dvsProps.RegistrarServer = openapi.PtrString(v) }},
		"registrar_server_secondary":            {plan.RegistrarServerSecondary, state.RegistrarServerSecondary, func(v string) { dvsProps.RegistrarServerSecondary = openapi.PtrString(v) }},
		"user_agent_domain":                     {plan.UserAgentDomain, state.UserAgentDomain, func(v string) { dvsProps.UserAgentDomain = openapi.PtrString(v) }},
		"user_agent_transport":                  {plan.UserAgentTransport, state.UserAgentTransport, func(v string) { dvsProps.UserAgentTransport = openapi.PtrString(v) }},
		"outbound_proxy":                        {plan.OutboundProxy, state.OutboundProxy, func(v string) { dvsProps.OutboundProxy = openapi.PtrString(v) }},
		"outbound_proxy_secondary":              {plan.OutboundProxySecondary, state.OutboundProxySecondary, func(v string) { dvsProps.OutboundProxySecondary = openapi.PtrString(v) }},
		"voicemail_server":                      {plan.VoicemailServer, state.VoicemailServer, func(v string) { dvsProps.VoicemailServer = openapi.PtrString(v) }},
		"call_agent_1":                          {plan.CallAgent1, state.CallAgent1, func(v string) { dvsProps.CallAgent1 = openapi.PtrString(v) }},
		"call_agent_2":                          {plan.CallAgent2, state.CallAgent2, func(v string) { dvsProps.CallAgent2 = openapi.PtrString(v) }},
		"domain":                                {plan.Domain, state.Domain, func(v string) { dvsProps.Domain = openapi.PtrString(v) }},
		"termination_base":                      {plan.TerminationBase, state.TerminationBase, func(v string) { dvsProps.TerminationBase = openapi.PtrString(v) }},
		"bit_rate":                              {plan.BitRate, state.BitRate, func(v string) { dvsProps.BitRate = openapi.PtrString(v) }},
		"cancel_call_waiting":                   {plan.CancelCallWaiting, state.CancelCallWaiting, func(v string) { dvsProps.CancelCallWaiting = openapi.PtrString(v) }},
		"call_hold":                             {plan.CallHold, state.CallHold, func(v string) { dvsProps.CallHold = openapi.PtrString(v) }},
		"cids_activate":                         {plan.CidsActivate, state.CidsActivate, func(v string) { dvsProps.CidsActivate = openapi.PtrString(v) }},
		"cids_deactivate":                       {plan.CidsDeactivate, state.CidsDeactivate, func(v string) { dvsProps.CidsDeactivate = openapi.PtrString(v) }},
		"do_not_disturb_activate":               {plan.DoNotDisturbActivate, state.DoNotDisturbActivate, func(v string) { dvsProps.DoNotDisturbActivate = openapi.PtrString(v) }},
		"do_not_disturb_deactivate":             {plan.DoNotDisturbDeactivate, state.DoNotDisturbDeactivate, func(v string) { dvsProps.DoNotDisturbDeactivate = openapi.PtrString(v) }},
		"do_not_disturb_pin_change":             {plan.DoNotDisturbPinChange, state.DoNotDisturbPinChange, func(v string) { dvsProps.DoNotDisturbPinChange = openapi.PtrString(v) }},
		"emergency_service_number":              {plan.EmergencyServiceNumber, state.EmergencyServiceNumber, func(v string) { dvsProps.EmergencyServiceNumber = openapi.PtrString(v) }},
		"anon_cid_block_activate":               {plan.AnonCidBlockActivate, state.AnonCidBlockActivate, func(v string) { dvsProps.AnonCidBlockActivate = openapi.PtrString(v) }},
		"anon_cid_block_deactivate":             {plan.AnonCidBlockDeactivate, state.AnonCidBlockDeactivate, func(v string) { dvsProps.AnonCidBlockDeactivate = openapi.PtrString(v) }},
		"call_forward_unconditional_activate":   {plan.CallForwardUnconditionalActivate, state.CallForwardUnconditionalActivate, func(v string) { dvsProps.CallForwardUnconditionalActivate = openapi.PtrString(v) }},
		"call_forward_unconditional_deactivate": {plan.CallForwardUnconditionalDeactivate, state.CallForwardUnconditionalDeactivate, func(v string) { dvsProps.CallForwardUnconditionalDeactivate = openapi.PtrString(v) }},
		"call_forward_on_busy_activate":         {plan.CallForwardOnBusyActivate, state.CallForwardOnBusyActivate, func(v string) { dvsProps.CallForwardOnBusyActivate = openapi.PtrString(v) }},
		"call_forward_on_busy_deactivate":       {plan.CallForwardOnBusyDeactivate, state.CallForwardOnBusyDeactivate, func(v string) { dvsProps.CallForwardOnBusyDeactivate = openapi.PtrString(v) }},
		"call_forward_on_no_answer_activate":    {plan.CallForwardOnNoAnswerActivate, state.CallForwardOnNoAnswerActivate, func(v string) { dvsProps.CallForwardOnNoAnswerActivate = openapi.PtrString(v) }},
		"call_forward_on_no_answer_deactivate":  {plan.CallForwardOnNoAnswerDeactivate, state.CallForwardOnNoAnswerDeactivate, func(v string) { dvsProps.CallForwardOnNoAnswerDeactivate = openapi.PtrString(v) }},
		"intercom_1":                            {plan.Intercom1, state.Intercom1, func(v string) { dvsProps.Intercom1 = openapi.PtrString(v) }},
		"intercom_2":                            {plan.Intercom2, state.Intercom2, func(v string) { dvsProps.Intercom2 = openapi.PtrString(v) }},
		"intercom_3":                            {plan.Intercom3, state.Intercom3, func(v string) { dvsProps.Intercom3 = openapi.PtrString(v) }},
	}

	for _, field := range stringFields {
		if !field.planField.Equal(field.stateField) {
			field.setter(field.planField.ValueString())
			hasChanges = true
		}
	}

	if !plan.Rtcp.Equal(state.Rtcp) {
		dvsProps.Rtcp = openapi.PtrBool(plan.Rtcp.ValueBool())
		hasChanges = true
	}
	if !plan.FaxT38.Equal(state.FaxT38) {
		dvsProps.FaxT38 = openapi.PtrBool(plan.FaxT38.ValueBool())
		hasChanges = true
	}

	intFields := map[string]struct {
		planField  types.Int64
		stateField types.Int64
		setter     func(types.Int64)
	}{
		"proxy_server_port":               {plan.ProxyServerPort, state.ProxyServerPort, func(v types.Int64) { setNullableInt32(&dvsProps.ProxyServerPort, v) }},
		"proxy_server_secondary_port":     {plan.ProxyServerSecondaryPort, state.ProxyServerSecondaryPort, func(v types.Int64) { setNullableInt32(&dvsProps.ProxyServerSecondaryPort, v) }},
		"registrar_server_port":           {plan.RegistrarServerPort, state.RegistrarServerPort, func(v types.Int64) { setNullableInt32(&dvsProps.RegistrarServerPort, v) }},
		"registrar_server_secondary_port": {plan.RegistrarServerSecondaryPort, state.RegistrarServerSecondaryPort, func(v types.Int64) { setNullableInt32(&dvsProps.RegistrarServerSecondaryPort, v) }},
		"user_agent_port":                 {plan.UserAgentPort, state.UserAgentPort, func(v types.Int64) { setNullableInt32(&dvsProps.UserAgentPort, v) }},
		"outbound_proxy_port":             {plan.OutboundProxyPort, state.OutboundProxyPort, func(v types.Int64) { setNullableInt32(&dvsProps.OutboundProxyPort, v) }},
		"outbound_proxy_secondary_port":   {plan.OutboundProxySecondaryPort, state.OutboundProxySecondaryPort, func(v types.Int64) { setNullableInt32(&dvsProps.OutboundProxySecondaryPort, v) }},
		"registration_period":             {plan.RegistrationPeriod, state.RegistrationPeriod, func(v types.Int64) { setNullableInt32(&dvsProps.RegistrationPeriod, v) }},
		"register_expires":                {plan.RegisterExpires, state.RegisterExpires, func(v types.Int64) { setNullableInt32(&dvsProps.RegisterExpires, v) }},
		"voicemail_server_port":           {plan.VoicemailServerPort, state.VoicemailServerPort, func(v types.Int64) { setNullableInt32(&dvsProps.VoicemailServerPort, v) }},
		"voicemail_server_expires":        {plan.VoicemailServerExpires, state.VoicemailServerExpires, func(v types.Int64) { setNullableInt32(&dvsProps.VoicemailServerExpires, v) }},
		"sip_dscp_mark":                   {plan.SipDscpMark, state.SipDscpMark, func(v types.Int64) { setNullableInt32(&dvsProps.SipDscpMark, v) }},
		"call_agent_port_1":               {plan.CallAgentPort1, state.CallAgentPort1, func(v types.Int64) { setNullableInt32(&dvsProps.CallAgentPort1, v) }},
		"call_agent_port_2":               {plan.CallAgentPort2, state.CallAgentPort2, func(v types.Int64) { setNullableInt32(&dvsProps.CallAgentPort2, v) }},
		"mgcp_dscp_mark":                  {plan.MgcpDscpMark, state.MgcpDscpMark, func(v types.Int64) { setNullableInt32(&dvsProps.MgcpDscpMark, v) }},
		"local_port_min":                  {plan.LocalPortMin, state.LocalPortMin, func(v types.Int64) { setNullableInt32(&dvsProps.LocalPortMin, v) }},
		"local_port_max":                  {plan.LocalPortMax, state.LocalPortMax, func(v types.Int64) { setNullableInt32(&dvsProps.LocalPortMax, v) }},
		"event_payload_type":              {plan.EventPayloadType, state.EventPayloadType, func(v types.Int64) { setNullableInt32(&dvsProps.EventPayloadType, v) }},
		"dscp_mark":                       {plan.DscpMark, state.DscpMark, func(v types.Int64) { setNullableInt32(&dvsProps.DscpMark, v) }},
	}

	for _, field := range intFields {
		if !field.planField.Equal(field.stateField) {
			field.setter(field.planField)
			hasChanges = true
		}
	}

	if !plan.CasEvents.Equal(state.CasEvents) {
		dvsProps.CasEvents = openapi.PtrInt32(int32(plan.CasEvents.ValueInt64()))
		hasChanges = true
	}

	oldCodecsByIndex := make(map[int64]verityDeviceVoiceSettingsCodecModel)
	for _, codec := range state.Codecs {
		if !codec.Index.IsNull() {
			idx := codec.Index.ValueInt64()
			oldCodecsByIndex[idx] = codec
		}
	}

	var changedCodecs []openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner
	codecsChanged := false

	for _, planCodec := range plan.Codecs {
		if planCodec.Index.IsNull() {
			continue // Skip items without identifier
		}

		idx := planCodec.Index.ValueInt64()
		stateCodec, exists := oldCodecsByIndex[idx]

		if !exists {
			// CREATE: new codec, include all fields
			newCodec := openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner{
				Index: openapi.PtrInt32(int32(idx)),
			}

			if !planCodec.CodecNumName.IsNull() && planCodec.CodecNumName.ValueString() != "" {
				newCodec.CodecNumName = openapi.PtrString(planCodec.CodecNumName.ValueString())
			} else {
				newCodec.CodecNumName = openapi.PtrString("")
			}

			if !planCodec.CodecNumEnable.IsNull() {
				newCodec.CodecNumEnable = openapi.PtrBool(planCodec.CodecNumEnable.ValueBool())
			} else {
				newCodec.CodecNumEnable = openapi.PtrBool(false)
			}

			if !planCodec.CodecNumPacketizationPeriod.IsNull() && planCodec.CodecNumPacketizationPeriod.ValueString() != "" {
				newCodec.CodecNumPacketizationPeriod = openapi.PtrString(planCodec.CodecNumPacketizationPeriod.ValueString())
			} else {
				newCodec.CodecNumPacketizationPeriod = openapi.PtrString("")
			}

			if !planCodec.CodecNumSilenceSuppression.IsNull() {
				newCodec.CodecNumSilenceSuppression = openapi.PtrBool(planCodec.CodecNumSilenceSuppression.ValueBool())
			} else {
				newCodec.CodecNumSilenceSuppression = openapi.PtrBool(false)
			}

			changedCodecs = append(changedCodecs, newCodec)
			codecsChanged = true
			continue
		}

		// UPDATE: existing codec, check which fields changed
		updateCodec := openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner{
			Index: openapi.PtrInt32(int32(idx)),
		}

		fieldChanged := false

		if !planCodec.CodecNumName.Equal(stateCodec.CodecNumName) {
			if !planCodec.CodecNumName.IsNull() && planCodec.CodecNumName.ValueString() != "" {
				updateCodec.CodecNumName = openapi.PtrString(planCodec.CodecNumName.ValueString())
			} else {
				updateCodec.CodecNumName = openapi.PtrString("")
			}
			fieldChanged = true
		}

		if !planCodec.CodecNumEnable.Equal(stateCodec.CodecNumEnable) {
			updateCodec.CodecNumEnable = openapi.PtrBool(planCodec.CodecNumEnable.ValueBool())
			fieldChanged = true
		}

		if !planCodec.CodecNumPacketizationPeriod.Equal(stateCodec.CodecNumPacketizationPeriod) {
			if !planCodec.CodecNumPacketizationPeriod.IsNull() && planCodec.CodecNumPacketizationPeriod.ValueString() != "" {
				updateCodec.CodecNumPacketizationPeriod = openapi.PtrString(planCodec.CodecNumPacketizationPeriod.ValueString())
			} else {
				updateCodec.CodecNumPacketizationPeriod = openapi.PtrString("")
			}
			fieldChanged = true
		}

		if !planCodec.CodecNumSilenceSuppression.Equal(stateCodec.CodecNumSilenceSuppression) {
			updateCodec.CodecNumSilenceSuppression = openapi.PtrBool(planCodec.CodecNumSilenceSuppression.ValueBool())
			fieldChanged = true
		}

		if fieldChanged {
			changedCodecs = append(changedCodecs, updateCodec)
			codecsChanged = true
		}
	}

	// DELETE: Check for deleted items
	for stateIdx := range oldCodecsByIndex {
		found := false
		for _, planCodec := range plan.Codecs {
			if !planCodec.Index.IsNull() && planCodec.Index.ValueInt64() == stateIdx {
				found = true
				break
			}
		}

		if !found {
			// codec removed - include only the index for deletion
			deletedCodec := openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner{
				Index: openapi.PtrInt32(int32(stateIdx)),
			}
			changedCodecs = append(changedCodecs, deletedCodec)
			codecsChanged = true
		}
	}

	if codecsChanged && len(changedCodecs) > 0 {
		dvsProps.Codecs = changedCodecs
		hasChanges = true
	}

	if len(plan.ObjectProperties) > 0 {
		if len(state.ObjectProperties) == 0 || !r.equalObjectProperties(plan.ObjectProperties[0], state.ObjectProperties[0]) {
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

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	operationID := r.bulkOpsMgr.AddPatch(ctx, "device_voice_settings", name, dvsProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for Device Voice Settings update operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update Device Voice Settings %s", name))...,
		)
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
	operationID := r.bulkOpsMgr.AddDelete(ctx, "device_voice_settings", name)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for Device Voice Settings deletion operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Delete Device Voice Settings %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Device Voice Settings %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "device_voice_settings")
	resp.State.RemoveResource(ctx)
}

func (r *verityDeviceVoiceSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func (r *verityDeviceVoiceSettingsResource) equalObjectProperties(a, b verityDeviceVoiceSettingsObjectPropertiesModel) bool {
	return a.IsDefault.Equal(b.IsDefault) && a.Group.Equal(b.Group)
}

func setNullableInt32(field *openapi.NullableInt32, value types.Int64) {
	if !value.IsNull() {
		val := int32(value.ValueInt64())
		*field = *openapi.NewNullableInt32(&val)
	} else {
		*field = *openapi.NewNullableInt32(nil)
	}
}
