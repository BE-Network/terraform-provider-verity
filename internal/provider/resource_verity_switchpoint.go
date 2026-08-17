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
	_ resource.Resource                = &veritySwitchpointResource{}
	_ resource.ResourceWithConfigure   = &veritySwitchpointResource{}
	_ resource.ResourceWithImportState = &veritySwitchpointResource{}
	_ resource.ResourceWithModifyPlan  = &veritySwitchpointResource{}
)

const switchpointResourceType = "switchpoints"
const switchpointTerraformType = "verity_switchpoint"

func NewVeritySwitchpointResource() resource.Resource {
	return &veritySwitchpointResource{}
}

type veritySwitchpointResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type veritySwitchpointResourceModel struct {
	Name                             types.String                             `tfsdk:"name"`
	Enable                           types.Bool                               `tfsdk:"enable"`
	Tenant                           types.String                             `tfsdk:"tenant"`
	TenantRefType                    types.String                             `tfsdk:"tenant_ref_type_"`
	DeviceSerialNumber               types.String                             `tfsdk:"device_serial_number"`
	ConnectedBundle                  types.String                             `tfsdk:"connected_bundle"`
	ConnectedBundleRefType           types.String                             `tfsdk:"connected_bundle_ref_type_"`
	IsTopOfIsland                    types.Bool                               `tfsdk:"is_top_of_island"`
	ReadOnlyMode                     types.Bool                               `tfsdk:"read_only_mode"`
	Locked                           types.Bool                               `tfsdk:"locked"`
	ExpectedSite                     types.String                             `tfsdk:"expected_site"`
	ExpectedSiteRefType              types.String                             `tfsdk:"expected_site_ref_type_"`
	OutOfBandManagement              types.Bool                               `tfsdk:"out_of_band_management"`
	Type                             types.String                             `tfsdk:"type"`
	Plane                            types.String                             `tfsdk:"plane"`
	PlaneRefType                     types.String                             `tfsdk:"plane_ref_type_"`
	SpinePlane                       types.String                             `tfsdk:"spine_plane"`
	SpinePlaneRefType                types.String                             `tfsdk:"spine_plane_ref_type_"`
	Pod                              types.String                             `tfsdk:"pod"`
	PodRefType                       types.String                             `tfsdk:"pod_ref_type_"`
	Su                               types.String                             `tfsdk:"su"`
	SuRefType                        types.String                             `tfsdk:"su_ref_type_"`
	SspGroup                         types.String                             `tfsdk:"ssp_group"`
	SspGroupRefType                  types.String                             `tfsdk:"ssp_group_ref_type_"`
	RackInfo                         types.String                             `tfsdk:"rack_info"`
	Rack                             types.String                             `tfsdk:"rack"`
	RackRefType                      types.String                             `tfsdk:"rack_ref_type_"`
	Position                         types.Number                             `tfsdk:"position"`
	RailGroup                        types.Number                             `tfsdk:"rail_group"`
	SwitchRouterIdIpMask             types.String                             `tfsdk:"switch_router_id_ip_mask"`
	SwitchRouterIdIpMaskAutoAssigned types.Bool                               `tfsdk:"switch_router_id_ip_mask_auto_assigned_"`
	SwitchVtepIdIpMask               types.String                             `tfsdk:"switch_vtep_id_ip_mask"`
	SwitchVtepIdIpMaskAutoAssigned   types.Bool                               `tfsdk:"switch_vtep_id_ip_mask_auto_assigned_"`
	BgpAsNumber                      types.Int64                              `tfsdk:"bgp_as_number"`
	BgpAsNumberAutoAssigned          types.Bool                               `tfsdk:"bgp_as_number_auto_assigned_"`
	BbSwitch                         types.Bool                               `tfsdk:"bb_switch"`
	PasswordEncrypted                types.String                             `tfsdk:"password_encrypted"`
	EnablePasswordEncrypted          types.String                             `tfsdk:"enable_password_encrypted"`
	SshKeyOrPasswordEncrypted        types.String                             `tfsdk:"ssh_key_or_password_encrypted"`
	PassphraseEncrypted              types.String                             `tfsdk:"passphrase_encrypted"`
	PrivatePasswordEncrypted         types.String                             `tfsdk:"private_password_encrypted"`
	IpSource                         types.String                             `tfsdk:"ip_source"`
	ControllerIpAndMask              types.String                             `tfsdk:"controller_ip_and_mask"`
	ControllerIpAndMaskAutoAssigned  types.Bool                               `tfsdk:"controller_ip_and_mask_auto_assigned_"`
	Gateway                          types.String                             `tfsdk:"gateway"`
	SwitchIpAndMask                  types.String                             `tfsdk:"switch_ip_and_mask"`
	SwitchIpAndMaskAutoAssigned      types.Bool                               `tfsdk:"switch_ip_and_mask_auto_assigned_"`
	SwitchGateway                    types.String                             `tfsdk:"switch_gateway"`
	CommType                         types.String                             `tfsdk:"comm_type"`
	SnmpCommunityString              types.String                             `tfsdk:"snmp_community_string"`
	UplinkPort                       types.String                             `tfsdk:"uplink_port"`
	ExpectedUplinkPort               types.Int64                              `tfsdk:"expected_uplink_port"`
	ExpectedBreakout                 types.String                             `tfsdk:"expected_breakout"`
	LldpSearchString                 types.String                             `tfsdk:"lldp_search_string"`
	ZtpIdentification                types.String                             `tfsdk:"ztp_identification"`
	LocatedBy                        types.String                             `tfsdk:"located_by"`
	PowerState                       types.String                             `tfsdk:"power_state"`
	CommunicationMode                types.String                             `tfsdk:"communication_mode"`
	CliAccessMode                    types.String                             `tfsdk:"cli_access_mode"`
	Username                         types.String                             `tfsdk:"username"`
	Password                         types.String                             `tfsdk:"password"`
	EnablePassword                   types.String                             `tfsdk:"enable_password"`
	SshKeyOrPassword                 types.String                             `tfsdk:"ssh_key_or_password"`
	ManagedOnNativeVlan              types.Bool                               `tfsdk:"managed_on_native_vlan"`
	Sdlc                             types.String                             `tfsdk:"sdlc"`
	SecurityType                     types.String                             `tfsdk:"security_type"`
	Snmpv3Username                   types.String                             `tfsdk:"snmpv3_username"`
	AuthenticationProtocol           types.String                             `tfsdk:"authentication_protocol"`
	Passphrase                       types.String                             `tfsdk:"passphrase"`
	PrivateProtocol                  types.String                             `tfsdk:"private_protocol"`
	PrivatePassword                  types.String                             `tfsdk:"private_password"`
	IsFabric                         types.Bool                               `tfsdk:"is_fabric"`
	DeviceManagedAs                  types.String                             `tfsdk:"device_managed_as"`
	Switch                           types.String                             `tfsdk:"switch"`
	SwitchRefType                    types.String                             `tfsdk:"switch_ref_type_"`
	ConnectionService                types.String                             `tfsdk:"connection_service"`
	ConnectionServiceRefType         types.String                             `tfsdk:"connection_service_ref_type_"`
	Port                             types.String                             `tfsdk:"port"`
	UsesTaggedPackets                types.Bool                               `tfsdk:"uses_tagged_packets"`
	Badges                           []veritySwitchpointBadgeModel            `tfsdk:"badges"`
	Children                         []veritySwitchpointChildModel            `tfsdk:"children"`
	TrafficMirrors                   []veritySwitchpointTrafficMirrorModel    `tfsdk:"traffic_mirrors"`
	Eths                             []veritySwitchpointEthModel              `tfsdk:"eths"`
	Pots                             []veritySwitchpointPotsModel             `tfsdk:"pots"`
	ObjectProperties                 []veritySwitchpointObjectPropertiesModel `tfsdk:"object_properties"`
}

type veritySwitchpointBadgeModel struct {
	Badge        types.String `tfsdk:"badge"`
	BadgeRefType types.String `tfsdk:"badge_ref_type_"`
	Index        types.Int64  `tfsdk:"index"`
}

func (b veritySwitchpointBadgeModel) GetIndex() types.Int64 {
	return b.Index
}

type veritySwitchpointChildModel struct {
	ChildNumEndpoint        types.String `tfsdk:"child_num_endpoint"`
	ChildNumEndpointRefType types.String `tfsdk:"child_num_endpoint_ref_type_"`
	ChildNumDevice          types.String `tfsdk:"child_num_device"`
	Index                   types.Int64  `tfsdk:"index"`
}

func (c veritySwitchpointChildModel) GetIndex() types.Int64 {
	return c.Index
}

type veritySwitchpointTrafficMirrorModel struct {
	TrafficMirrorNumEnable             types.Bool   `tfsdk:"traffic_mirror_num_enable"`
	TrafficMirrorNumSourcePort         types.String `tfsdk:"traffic_mirror_num_source_port"`
	TrafficMirrorNumSourceLagIndicator types.Bool   `tfsdk:"traffic_mirror_num_source_lag_indicator"`
	TrafficMirrorNumDestinationPort    types.String `tfsdk:"traffic_mirror_num_destination_port"`
	TrafficMirrorNumInboundTraffic     types.Bool   `tfsdk:"traffic_mirror_num_inbound_traffic"`
	TrafficMirrorNumOutboundTraffic    types.Bool   `tfsdk:"traffic_mirror_num_outbound_traffic"`
	Index                              types.Int64  `tfsdk:"index"`
}

func (tm veritySwitchpointTrafficMirrorModel) GetIndex() types.Int64 {
	return tm.Index
}

type veritySwitchpointEthModel struct {
	Breakout     types.String `tfsdk:"breakout"`
	CustomerVlan types.String `tfsdk:"customer_vlan"`
	Index        types.Int64  `tfsdk:"index"`
	EthNumIcon   types.String `tfsdk:"eth_num_icon"`
	EthNumLabel  types.String `tfsdk:"eth_num_label"`
	Enable       types.Bool   `tfsdk:"enable"`
	PortName     types.String `tfsdk:"port_name"`
}

func (e veritySwitchpointEthModel) GetIndex() types.Int64 {
	return e.Index
}

type veritySwitchpointPotsModel struct {
	PotsNumEnable            types.Bool   `tfsdk:"pots_num_enable"`
	PotsNumUri               types.String `tfsdk:"pots_num_uri"`
	PotsNumUsername          types.String `tfsdk:"pots_num_username"`
	PotsNumPassword          types.String `tfsdk:"pots_num_password"`
	PotsNumCallerId          types.String `tfsdk:"pots_num_caller_id"`
	PotsNumHotLine           types.String `tfsdk:"pots_num_hot_line"`
	PotsNumPasswordEncrypted types.String `tfsdk:"pots_num_password_encrypted"`
	Index                    types.Int64  `tfsdk:"index"`
}

func (p veritySwitchpointPotsModel) GetIndex() types.Int64 {
	return p.Index
}

type veritySwitchpointObjectPropertiesModel struct {
	UserNotes                     types.String `tfsdk:"user_notes"`
	ExpectedParentEndpoint        types.String `tfsdk:"expected_parent_endpoint"`
	ExpectedParentEndpointRefType types.String `tfsdk:"expected_parent_endpoint_ref_type_"`
	NumberOfMultipoints           types.Int64  `tfsdk:"number_of_multipoints"`
	Aggregate                     types.Bool   `tfsdk:"aggregate"`
	IsHost                        types.Bool   `tfsdk:"is_host"`
	EmulateRfVideoPort            types.Bool   `tfsdk:"emulate_rf_video_port"`
	DrawAsEdgeDevice              types.Bool   `tfsdk:"draw_as_edge_device"`
}

func (r *veritySwitchpointResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_switchpoint"
}

func (r *veritySwitchpointResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *veritySwitchpointResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity Switchpoint",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Object Name. Must be unique.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enable": schema.BoolAttribute{
				Description: "Enable object. It's highly recommended to set this value to true so that validation on the object will be ran.",
				Optional:    true,
				Computed:    true,
			},
			"tenant": schema.StringAttribute{
				Description: "The Tenant of this Device",
				Optional:    true,
				Computed:    true,
			},
			"tenant_ref_type_": schema.StringAttribute{
				Description: "Object type for tenant field",
				Optional:    true,
				Computed:    true,
			},
			"device_serial_number": schema.StringAttribute{
				Description: "Device Serial Number",
				Optional:    true,
				Computed:    true,
			},
			"connected_bundle": schema.StringAttribute{
				Description: "Connected Bundle",
				Optional:    true,
				Computed:    true,
			},
			"connected_bundle_ref_type_": schema.StringAttribute{
				Description: "Object type for connected_bundle field",
				Optional:    true,
				Computed:    true,
			},
			"is_top_of_island": schema.BoolAttribute{
				Description: "Mark this Switchpoint as Top of Island",
				Optional:    true,
				Computed:    true,
			},
			"read_only_mode": schema.BoolAttribute{
				Description: "When Read Only Mode is checked, vNetC will perform all functions except writing database updates to the target hardware",
				Optional:    true,
				Computed:    true,
			},
			"locked": schema.BoolAttribute{
				Description: "Permission lock",
				Optional:    true,
				Computed:    true,
			},
			"expected_site": schema.StringAttribute{
				Description: "Expected Fabric",
				Optional:    true,
				Computed:    true,
			},
			"expected_site_ref_type_": schema.StringAttribute{
				Description: "Object type for expected_site field",
				Optional:    true,
				Computed:    true,
			},
			"out_of_band_management": schema.BoolAttribute{
				Description: "For Switch Endpoints. Denotes a Switch is managed out of band via the management port",
				Optional:    true,
				Computed:    true,
			},
			"type": schema.StringAttribute{
				Description: "Type of Switchpoint",
				Optional:    true,
				Computed:    true,
			},
			"plane": schema.StringAttribute{
				Description: "Plane",
				Optional:    true,
				Computed:    true,
			},
			"plane_ref_type_": schema.StringAttribute{
				Description: "Object type for plane field",
				Optional:    true,
				Computed:    true,
			},
			"spine_plane": schema.StringAttribute{
				Description: "Spine Plane - subgrouping of super spine and spine",
				Optional:    true,
				Computed:    true,
			},
			"spine_plane_ref_type_": schema.StringAttribute{
				Description: "Object type for spine_plane field",
				Optional:    true,
				Computed:    true,
			},
			"pod": schema.StringAttribute{
				Description: "Pod - subgrouping of spine and leaf switches",
				Optional:    true,
				Computed:    true,
			},
			"pod_ref_type_": schema.StringAttribute{
				Description: "Object type for pod field",
				Optional:    true,
				Computed:    true,
			},
			"su": schema.StringAttribute{
				Description: "SU",
				Optional:    true,
				Computed:    true,
			},
			"su_ref_type_": schema.StringAttribute{
				Description: "Object type for su field",
				Optional:    true,
				Computed:    true,
			},
			"ssp_group": schema.StringAttribute{
				Description: "SuperSpine Group - grouping of superspines in 3-tier config",
				Optional:    true,
				Computed:    true,
			},
			"ssp_group_ref_type_": schema.StringAttribute{
				Description: "Object type for ssp_group field",
				Optional:    true,
				Computed:    true,
			},
			"rack_info": schema.StringAttribute{
				Description: "Physical Rack location of the Switch",
				Optional:    true,
				Computed:    true,
			},
			"rack": schema.StringAttribute{
				Description: "Rack",
				Optional:    true,
				Computed:    true,
			},
			"rack_ref_type_": schema.StringAttribute{
				Description: "Object type for rack field",
				Optional:    true,
				Computed:    true,
			},
			"position": schema.NumberAttribute{
				Description: "Position of the Switch",
				Optional:    true,
				Computed:    true,
			},
			"rail_group": schema.NumberAttribute{
				Description: "Rail Group the Switch is part of",
				Optional:    true,
				Computed:    true,
			},
			"switch_router_id_ip_mask": schema.StringAttribute{
				Description: "Switch BGP Router Identifier. This field should not be specified when 'switch_router_id_ip_mask_auto_assigned_' is set to true, as the API will assign this value automatically.",
				Optional:    true,
				Computed:    true,
			},
			"switch_router_id_ip_mask_auto_assigned_": schema.BoolAttribute{
				Description: "Whether the Switch BGP Router Identifier should be automatically assigned by the API. When set to true, do not specify the 'switch_router_id_ip_mask' field in your configuration.",
				Optional:    true,
				Computed:    true,
			},
			"switch_vtep_id_ip_mask": schema.StringAttribute{
				Description: "Switch VTEP Identifier. This field should not be specified when 'switch_vtep_id_ip_mask_auto_assigned_' is set to true, as the API will assign this value automatically.",
				Optional:    true,
				Computed:    true,
			},
			"switch_vtep_id_ip_mask_auto_assigned_": schema.BoolAttribute{
				Description: "Whether the Switch VTEP Identifier should be automatically assigned by the API. When set to true, do not specify the 'switch_vtep_id_ip_mask' field in your configuration.",
				Optional:    true,
				Computed:    true,
			},
			"bgp_as_number": schema.Int64Attribute{
				Description: "BGP Autonomous System Number for the site underlay. This field should not be specified when 'bgp_as_number_auto_assigned_' is set to true, as the API will assign this value automatically.",
				Optional:    true,
				Computed:    true,
			},
			"bgp_as_number_auto_assigned_": schema.BoolAttribute{
				Description: "Whether the BGP AS Number should be automatically assigned by the API. When set to true, do not specify the 'bgp_as_number' field in your configuration.",
				Optional:    true,
				Computed:    true,
			},
			"bb_switch": schema.BoolAttribute{
				Description: "Expose fields for Device Management",
				Optional:    true,
				Computed:    true,
			},
			"password_encrypted": schema.StringAttribute{
				Description: "Password",
				Optional:    true,
				Computed:    true,
			},
			"enable_password_encrypted": schema.StringAttribute{
				Description: "Enable Password - to enable privileged CLI operations",
				Optional:    true,
				Computed:    true,
			},
			"ssh_key_or_password_encrypted": schema.StringAttribute{
				Description: "SSH Key or Password",
				Optional:    true,
				Computed:    true,
			},
			"passphrase_encrypted": schema.StringAttribute{
				Description: "Passphrase",
				Optional:    true,
				Computed:    true,
			},
			"private_password_encrypted": schema.StringAttribute{
				Description: "Password",
				Optional:    true,
				Computed:    true,
			},
			"ip_source": schema.StringAttribute{
				Description: "IP Source",
				Optional:    true,
				Computed:    true,
			},
			"controller_ip_and_mask": schema.StringAttribute{
				Description: "Controller IP and Mask. This field should not be specified when 'controller_ip_and_mask_auto_assigned_' is set to true, as the API will assign this value automatically.",
				Optional:    true,
				Computed:    true,
			},
			"controller_ip_and_mask_auto_assigned_": schema.BoolAttribute{
				Description: "Whether the Controller IP and Mask should be automatically assigned by the API. When set to true, do not specify the 'controller_ip_and_mask' field in your configuration.",
				Optional:    true,
				Computed:    true,
			},
			"gateway": schema.StringAttribute{
				Description: "Gateway",
				Optional:    true,
				Computed:    true,
			},
			"switch_ip_and_mask": schema.StringAttribute{
				Description: "Switch IP and Mask. This field should not be specified when 'switch_ip_and_mask_auto_assigned_' is set to true, as the API will assign this value automatically.",
				Optional:    true,
				Computed:    true,
			},
			"switch_ip_and_mask_auto_assigned_": schema.BoolAttribute{
				Description: "Whether the Switch IP and Mask should be automatically assigned by the API. When set to true, do not specify the 'switch_ip_and_mask' field in your configuration.",
				Optional:    true,
				Computed:    true,
			},
			"switch_gateway": schema.StringAttribute{
				Description: "Gateway of Managed Device",
				Optional:    true,
				Computed:    true,
			},
			"comm_type": schema.StringAttribute{
				Description: "Comm Type",
				Optional:    true,
				Computed:    true,
			},
			"snmp_community_string": schema.StringAttribute{
				Description: "Comm Credentials",
				Optional:    true,
				Computed:    true,
			},
			"uplink_port": schema.StringAttribute{
				Description: "Uplink Port of Managed Device",
				Optional:    true,
				Computed:    true,
			},
			"expected_uplink_port": schema.Int64Attribute{
				Description: "First breakout port, using a one-based index, that ZTP configures as an uplink for SFP-based ports",
				Optional:    true,
				Computed:    true,
			},
			"expected_breakout": schema.StringAttribute{
				Description: "Full breakout configuration for the SFP used as the uplink",
				Optional:    true,
				Computed:    true,
			},
			"lldp_search_string": schema.StringAttribute{
				Description: "Optional unless Located By is LLDP or Device managed as Active SFP",
				Optional:    true,
				Computed:    true,
			},
			"ztp_identification": schema.StringAttribute{
				Description: "Service Tag or Serial Number to identify device for Zero Touch Provisioning",
				Optional:    true,
				Computed:    true,
			},
			"located_by": schema.StringAttribute{
				Description: "Controls how the system locates this Device within its LAN",
				Optional:    true,
				Computed:    true,
			},
			"power_state": schema.StringAttribute{
				Description: "Power state of Switch Controller",
				Optional:    true,
				Computed:    true,
			},
			"communication_mode": schema.StringAttribute{
				Description: "Communication Mode",
				Optional:    true,
				Computed:    true,
			},
			"cli_access_mode": schema.StringAttribute{
				Description: "CLI Access Mode",
				Optional:    true,
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username",
				Optional:    true,
				Computed:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password",
				Optional:    true,
				Computed:    true,
			},
			"enable_password": schema.StringAttribute{
				Description: "Enable Password - to enable privileged CLI operations",
				Optional:    true,
				Computed:    true,
			},
			"ssh_key_or_password": schema.StringAttribute{
				Description: "SSH Key or Password",
				Optional:    true,
				Computed:    true,
			},
			"managed_on_native_vlan": schema.BoolAttribute{
				Description: "Managed on native VLAN",
				Optional:    true,
				Computed:    true,
			},
			"sdlc": schema.StringAttribute{
				Description: "SDLC that Device Controller belongs to",
				Optional:    true,
				Computed:    true,
			},
			"security_type": schema.StringAttribute{
				Description: "Security level",
				Optional:    true,
				Computed:    true,
			},
			"snmpv3_username": schema.StringAttribute{
				Description: "Username",
				Optional:    true,
				Computed:    true,
			},
			"authentication_protocol": schema.StringAttribute{
				Description: "Protocol",
				Optional:    true,
				Computed:    true,
			},
			"passphrase": schema.StringAttribute{
				Description: "Passphrase",
				Optional:    true,
				Computed:    true,
			},
			"private_protocol": schema.StringAttribute{
				Description: "Protocol",
				Optional:    true,
				Computed:    true,
			},
			"private_password": schema.StringAttribute{
				Description: "Password",
				Optional:    true,
				Computed:    true,
			},
			"is_fabric": schema.BoolAttribute{
				Description: "For Switch Endpoints. Denotes a Switch that is Fabric rather than an Edge Device",
				Optional:    true,
				Computed:    true,
			},
			"device_managed_as": schema.StringAttribute{
				Description: "Device managed as",
				Optional:    true,
				Computed:    true,
			},
			"switch": schema.StringAttribute{
				Description: "Switchpoint locating the Switch to be controlled",
				Optional:    true,
				Computed:    true,
			},
			"switch_ref_type_": schema.StringAttribute{
				Description: "Object type for switch field",
				Optional:    true,
				Computed:    true,
			},
			"connection_service": schema.StringAttribute{
				Description: "Connect a Service",
				Optional:    true,
				Computed:    true,
			},
			"connection_service_ref_type_": schema.StringAttribute{
				Description: "Object type for connection_service field",
				Optional:    true,
				Computed:    true,
			},
			"port": schema.StringAttribute{
				Description: "Port locating the Switch to be controlled",
				Optional:    true,
				Computed:    true,
			},
			"uses_tagged_packets": schema.BoolAttribute{
				Description: "Indicates if the direct interface expects tagged or untagged packets",
				Optional:    true,
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"badges": schema.ListNestedBlock{
				Description: "Badge configurations",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"badge": schema.StringAttribute{
							Description: "Badge name",
							Optional:    true,
							Computed:    true,
						},
						"badge_ref_type_": schema.StringAttribute{
							Description: "Object type for badge field",
							Optional:    true,
							Computed:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"children": schema.ListNestedBlock{
				Description: "Child configurations",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"child_num_endpoint": schema.StringAttribute{
							Description: "Switchpoint associated with the Child",
							Optional:    true,
							Computed:    true,
						},
						"child_num_endpoint_ref_type_": schema.StringAttribute{
							Description: "Object type for child_num_endpoint field",
							Optional:    true,
							Computed:    true,
						},
						"child_num_device": schema.StringAttribute{
							Description: "Device associated with the Child",
							Optional:    true,
							Computed:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"traffic_mirrors": schema.ListNestedBlock{
				Description: "Traffic mirror configurations",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"traffic_mirror_num_enable": schema.BoolAttribute{
							Description: "Enable Traffic Mirror",
							Optional:    true,
							Computed:    true,
						},
						"traffic_mirror_num_source_port": schema.StringAttribute{
							Description: "Source Port for Traffic Mirror",
							Optional:    true,
							Computed:    true,
						},
						"traffic_mirror_num_source_lag_indicator": schema.BoolAttribute{
							Description: "Source LAG Indicator for Traffic Mirror",
							Optional:    true,
							Computed:    true,
						},
						"traffic_mirror_num_destination_port": schema.StringAttribute{
							Description: "Destination Port for Traffic Mirror",
							Optional:    true,
							Computed:    true,
						},
						"traffic_mirror_num_inbound_traffic": schema.BoolAttribute{
							Description: "Boolean value indicating if the mirror is for inbound traffic",
							Optional:    true,
							Computed:    true,
						},
						"traffic_mirror_num_outbound_traffic": schema.BoolAttribute{
							Description: "Boolean value indicating if the mirror is for outbound traffic",
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
			"eths": schema.ListNestedBlock{
				Description: "Ethernet port configurations",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"breakout": schema.StringAttribute{
							Description: "Breakout Port Override",
							Optional:    true,
							Computed:    true,
						},
						"customer_vlan": schema.StringAttribute{
							Description: "A Value between 1 and 4096",
							Optional:    true,
							Computed:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object",
							Optional:    true,
							Computed:    true,
						},
						"eth_num_icon": schema.StringAttribute{
							Description: "Icon of this Eth Port",
							Optional:    true,
							Computed:    true,
						},
						"eth_num_label": schema.StringAttribute{
							Description: "Label of this Eth Port",
							Optional:    true,
							Computed:    true,
						},
						"enable": schema.BoolAttribute{
							Description: "Enable port",
							Optional:    true,
							Computed:    true,
						},
						"port_name": schema.StringAttribute{
							Description: "The name identifying the port. Used for reference only, it won't actually change the port name.",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"pots": schema.ListNestedBlock{
				Description: "POTS configurations",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"pots_num_enable": schema.BoolAttribute{
							Description: "Enable POTS port",
							Optional:    true,
							Computed:    true,
						},
						"pots_num_uri": schema.StringAttribute{
							Description: "Specific telephone extension for SIP for POTS port",
							Optional:    true,
							Computed:    true,
						},
						"pots_num_username": schema.StringAttribute{
							Description: "SIP username used for authentication for POTS port",
							Optional:    true,
							Computed:    true,
						},
						"pots_num_password": schema.StringAttribute{
							Description: "SIP password used for authentication for POTS port",
							Optional:    true,
							Computed:    true,
						},
						"pots_num_caller_id": schema.StringAttribute{
							Description: "ASCII string defining the user for the Caller ID display for POTS port",
							Optional:    true,
							Computed:    true,
						},
						"pots_num_hot_line": schema.StringAttribute{
							Description: "URI of line to autodial upon off-hook for POTS port",
							Optional:    true,
							Computed:    true,
						},
						"pots_num_password_encrypted": schema.StringAttribute{
							Description: "SIP password used for authentication for POTS port",
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
				Description: "Object properties for the switchpoint",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"user_notes": schema.StringAttribute{
							Description: "Notes written by User about the site",
							Optional:    true,
							Computed:    true,
						},
						"expected_parent_endpoint": schema.StringAttribute{
							Description: "Expected Parent Endpoint",
							Optional:    true,
							Computed:    true,
						},
						"expected_parent_endpoint_ref_type_": schema.StringAttribute{
							Description: "Object type for expected_parent_endpoint field",
							Optional:    true,
							Computed:    true,
						},
						"number_of_multipoints": schema.Int64Attribute{
							Description: "Number of Multipoints",
							Optional:    true,
							Computed:    true,
						},
						"aggregate": schema.BoolAttribute{
							Description: "For Switch Endpoints. Denotes switch aggregated with all of its sub switches",
							Optional:    true,
							Computed:    true,
						},
						"is_host": schema.BoolAttribute{
							Description: "For Switch Endpoints. Denotes the Host Switch",
							Optional:    true,
							Computed:    true,
						},
						"emulate_rf_video_port": schema.BoolAttribute{
							Description: "Emulate RF Video Port",
							Optional:    true,
							Computed:    true,
						},
						"draw_as_edge_device": schema.BoolAttribute{
							Description: "Turn on to display the switch as an edge device instead of as a switch",
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func (r *veritySwitchpointResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan veritySwitchpointResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config veritySwitchpointResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate auto-assigned field specifications
	if !plan.BgpAsNumberAutoAssigned.IsNull() && plan.BgpAsNumberAutoAssigned.ValueBool() {
		if !plan.BgpAsNumber.IsNull() && !plan.BgpAsNumber.IsUnknown() {
			resp.Diagnostics.AddError(
				"BGP AS Number cannot be specified when auto-assigned",
				"The 'bgp_as_number' field cannot be specified in the configuration when 'bgp_as_number_auto_assigned_' is set to true. The API will assign this value automatically.",
			)
			return
		}
	}

	if !plan.SwitchRouterIdIpMaskAutoAssigned.IsNull() && plan.SwitchRouterIdIpMaskAutoAssigned.ValueBool() {
		if !plan.SwitchRouterIdIpMask.IsNull() && !plan.SwitchRouterIdIpMask.IsUnknown() && plan.SwitchRouterIdIpMask.ValueString() != "" {
			resp.Diagnostics.AddError(
				"Switch Router ID IP Mask cannot be specified when auto-assigned",
				"The 'switch_router_id_ip_mask' field cannot be specified in the configuration when 'switch_router_id_ip_mask_auto_assigned_' is set to true. The API will assign this value automatically.",
			)
			return
		}
	}

	if !plan.SwitchVtepIdIpMaskAutoAssigned.IsNull() && plan.SwitchVtepIdIpMaskAutoAssigned.ValueBool() {
		if !plan.SwitchVtepIdIpMask.IsNull() && !plan.SwitchVtepIdIpMask.IsUnknown() && plan.SwitchVtepIdIpMask.ValueString() != "" {
			resp.Diagnostics.AddError(
				"Switch VTEP ID IP Mask cannot be specified when auto-assigned",
				"The 'switch_vtep_id_ip_mask' field cannot be specified in the configuration when 'switch_vtep_id_ip_mask_auto_assigned_' is set to true. The API will assign this value automatically.",
			)
			return
		}
	}

	if !plan.ControllerIpAndMaskAutoAssigned.IsNull() && plan.ControllerIpAndMaskAutoAssigned.ValueBool() {
		if !plan.ControllerIpAndMask.IsNull() && !plan.ControllerIpAndMask.IsUnknown() && plan.ControllerIpAndMask.ValueString() != "" {
			resp.Diagnostics.AddError(
				"Controller IP and Mask cannot be specified when auto-assigned",
				"The 'controller_ip_and_mask' field cannot be specified in the configuration when 'controller_ip_and_mask_auto_assigned_' is set to true. The API will assign this value automatically.",
			)
			return
		}
	}

	if !plan.SwitchIpAndMaskAutoAssigned.IsNull() && plan.SwitchIpAndMaskAutoAssigned.ValueBool() {
		if !plan.SwitchIpAndMask.IsNull() && !plan.SwitchIpAndMask.IsUnknown() && plan.SwitchIpAndMask.ValueString() != "" {
			resp.Diagnostics.AddError(
				"Switch IP and Mask cannot be specified when auto-assigned",
				"The 'switch_ip_and_mask' field cannot be specified in the configuration when 'switch_ip_and_mask_auto_assigned_' is set to true. The API will assign this value automatically.",
			)
			return
		}
	}

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Authenticate",
			fmt.Sprintf("Error authenticating with API: %s", err),
		)
		return
	}

	name := plan.Name.ValueString()
	spProps := &openapi.SwitchpointsPutRequestSwitchpointValue{
		Name: openapi.PtrString(name),
	}

	// Handle string fields
	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "Tenant", APIField: &spProps.Tenant, TFValue: plan.Tenant},
		{FieldName: "TenantRefType", APIField: &spProps.TenantRefType, TFValue: plan.TenantRefType},
		{FieldName: "DeviceSerialNumber", APIField: &spProps.DeviceSerialNumber, TFValue: plan.DeviceSerialNumber},
		{FieldName: "ConnectedBundle", APIField: &spProps.ConnectedBundle, TFValue: plan.ConnectedBundle},
		{FieldName: "ConnectedBundleRefType", APIField: &spProps.ConnectedBundleRefType, TFValue: plan.ConnectedBundleRefType},
		{FieldName: "ExpectedSite", APIField: &spProps.ExpectedSite, TFValue: plan.ExpectedSite},
		{FieldName: "ExpectedSiteRefType", APIField: &spProps.ExpectedSiteRefType, TFValue: plan.ExpectedSiteRefType},
		{FieldName: "Type", APIField: &spProps.Type, TFValue: plan.Type},
		{FieldName: "Plane", APIField: &spProps.Plane, TFValue: plan.Plane},
		{FieldName: "PlaneRefType", APIField: &spProps.PlaneRefType, TFValue: plan.PlaneRefType},
		{FieldName: "SpinePlane", APIField: &spProps.SpinePlane, TFValue: plan.SpinePlane},
		{FieldName: "SpinePlaneRefType", APIField: &spProps.SpinePlaneRefType, TFValue: plan.SpinePlaneRefType},
		{FieldName: "Pod", APIField: &spProps.Pod, TFValue: plan.Pod},
		{FieldName: "PodRefType", APIField: &spProps.PodRefType, TFValue: plan.PodRefType},
		{FieldName: "Su", APIField: &spProps.Su, TFValue: plan.Su},
		{FieldName: "SuRefType", APIField: &spProps.SuRefType, TFValue: plan.SuRefType},
		{FieldName: "SspGroup", APIField: &spProps.SspGroup, TFValue: plan.SspGroup},
		{FieldName: "SspGroupRefType", APIField: &spProps.SspGroupRefType, TFValue: plan.SspGroupRefType},
		{FieldName: "RackInfo", APIField: &spProps.RackInfo, TFValue: plan.RackInfo},
		{FieldName: "Rack", APIField: &spProps.Rack, TFValue: plan.Rack},
		{FieldName: "RackRefType", APIField: &spProps.RackRefType, TFValue: plan.RackRefType},
		{FieldName: "SwitchRouterIdIpMask", APIField: &spProps.SwitchRouterIdIpMask, TFValue: plan.SwitchRouterIdIpMask},
		{FieldName: "SwitchVtepIdIpMask", APIField: &spProps.SwitchVtepIdIpMask, TFValue: plan.SwitchVtepIdIpMask},
		{FieldName: "PasswordEncrypted", APIField: &spProps.PasswordEncrypted, TFValue: plan.PasswordEncrypted},
		{FieldName: "EnablePasswordEncrypted", APIField: &spProps.EnablePasswordEncrypted, TFValue: plan.EnablePasswordEncrypted},
		{FieldName: "SshKeyOrPasswordEncrypted", APIField: &spProps.SshKeyOrPasswordEncrypted, TFValue: plan.SshKeyOrPasswordEncrypted},
		{FieldName: "PassphraseEncrypted", APIField: &spProps.PassphraseEncrypted, TFValue: plan.PassphraseEncrypted},
		{FieldName: "PrivatePasswordEncrypted", APIField: &spProps.PrivatePasswordEncrypted, TFValue: plan.PrivatePasswordEncrypted},
		{FieldName: "IpSource", APIField: &spProps.IpSource, TFValue: plan.IpSource},
		{FieldName: "ControllerIpAndMask", APIField: &spProps.ControllerIpAndMask, TFValue: plan.ControllerIpAndMask},
		{FieldName: "Gateway", APIField: &spProps.Gateway, TFValue: plan.Gateway},
		{FieldName: "SwitchIpAndMask", APIField: &spProps.SwitchIpAndMask, TFValue: plan.SwitchIpAndMask},
		{FieldName: "SwitchGateway", APIField: &spProps.SwitchGateway, TFValue: plan.SwitchGateway},
		{FieldName: "CommType", APIField: &spProps.CommType, TFValue: plan.CommType},
		{FieldName: "SnmpCommunityString", APIField: &spProps.SnmpCommunityString, TFValue: plan.SnmpCommunityString},
		{FieldName: "UplinkPort", APIField: &spProps.UplinkPort, TFValue: plan.UplinkPort},
		{FieldName: "ExpectedBreakout", APIField: &spProps.ExpectedBreakout, TFValue: plan.ExpectedBreakout},
		{FieldName: "LldpSearchString", APIField: &spProps.LldpSearchString, TFValue: plan.LldpSearchString},
		{FieldName: "ZtpIdentification", APIField: &spProps.ZtpIdentification, TFValue: plan.ZtpIdentification},
		{FieldName: "LocatedBy", APIField: &spProps.LocatedBy, TFValue: plan.LocatedBy},
		{FieldName: "PowerState", APIField: &spProps.PowerState, TFValue: plan.PowerState},
		{FieldName: "CommunicationMode", APIField: &spProps.CommunicationMode, TFValue: plan.CommunicationMode},
		{FieldName: "CliAccessMode", APIField: &spProps.CliAccessMode, TFValue: plan.CliAccessMode},
		{FieldName: "Username", APIField: &spProps.Username, TFValue: plan.Username},
		{FieldName: "Password", APIField: &spProps.Password, TFValue: plan.Password},
		{FieldName: "EnablePassword", APIField: &spProps.EnablePassword, TFValue: plan.EnablePassword},
		{FieldName: "SshKeyOrPassword", APIField: &spProps.SshKeyOrPassword, TFValue: plan.SshKeyOrPassword},
		{FieldName: "Sdlc", APIField: &spProps.Sdlc, TFValue: plan.Sdlc},
		{FieldName: "SecurityType", APIField: &spProps.SecurityType, TFValue: plan.SecurityType},
		{FieldName: "Snmpv3Username", APIField: &spProps.Snmpv3Username, TFValue: plan.Snmpv3Username},
		{FieldName: "AuthenticationProtocol", APIField: &spProps.AuthenticationProtocol, TFValue: plan.AuthenticationProtocol},
		{FieldName: "Passphrase", APIField: &spProps.Passphrase, TFValue: plan.Passphrase},
		{FieldName: "PrivateProtocol", APIField: &spProps.PrivateProtocol, TFValue: plan.PrivateProtocol},
		{FieldName: "PrivatePassword", APIField: &spProps.PrivatePassword, TFValue: plan.PrivatePassword},
		{FieldName: "DeviceManagedAs", APIField: &spProps.DeviceManagedAs, TFValue: plan.DeviceManagedAs},
		{FieldName: "Switch", APIField: &spProps.Switch, TFValue: plan.Switch},
		{FieldName: "SwitchRefType", APIField: &spProps.SwitchRefType, TFValue: plan.SwitchRefType},
		{FieldName: "ConnectionService", APIField: &spProps.ConnectionService, TFValue: plan.ConnectionService},
		{FieldName: "ConnectionServiceRefType", APIField: &spProps.ConnectionServiceRefType, TFValue: plan.ConnectionServiceRefType},
		{FieldName: "Port", APIField: &spProps.Port, TFValue: plan.Port},
	})

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &spProps.Enable, TFValue: plan.Enable},
		{FieldName: "IsTopOfIsland", APIField: &spProps.IsTopOfIsland, TFValue: plan.IsTopOfIsland},
		{FieldName: "ReadOnlyMode", APIField: &spProps.ReadOnlyMode, TFValue: plan.ReadOnlyMode},
		{FieldName: "Locked", APIField: &spProps.Locked, TFValue: plan.Locked},
		{FieldName: "OutOfBandManagement", APIField: &spProps.OutOfBandManagement, TFValue: plan.OutOfBandManagement},
		{FieldName: "SwitchRouterIdIpMaskAutoAssigned", APIField: &spProps.SwitchRouterIdIpMaskAutoAssigned, TFValue: plan.SwitchRouterIdIpMaskAutoAssigned},
		{FieldName: "SwitchVtepIdIpMaskAutoAssigned", APIField: &spProps.SwitchVtepIdIpMaskAutoAssigned, TFValue: plan.SwitchVtepIdIpMaskAutoAssigned},
		{FieldName: "BgpAsNumberAutoAssigned", APIField: &spProps.BgpAsNumberAutoAssigned, TFValue: plan.BgpAsNumberAutoAssigned},
		{FieldName: "ControllerIpAndMaskAutoAssigned", APIField: &spProps.ControllerIpAndMaskAutoAssigned, TFValue: plan.ControllerIpAndMaskAutoAssigned},
		{FieldName: "SwitchIpAndMaskAutoAssigned", APIField: &spProps.SwitchIpAndMaskAutoAssigned, TFValue: plan.SwitchIpAndMaskAutoAssigned},
		{FieldName: "BbSwitch", APIField: &spProps.BbSwitch, TFValue: plan.BbSwitch},
		{FieldName: "ManagedOnNativeVlan", APIField: &spProps.ManagedOnNativeVlan, TFValue: plan.ManagedOnNativeVlan},
		{FieldName: "IsFabric", APIField: &spProps.IsFabric, TFValue: plan.IsFabric},
		{FieldName: "UsesTaggedPackets", APIField: &spProps.UsesTaggedPackets, TFValue: plan.UsesTaggedPackets},
	})

	// Handle nullable int64 fields - parse HCL to detect explicit config
	workDir := r.provCtx.workDir
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, switchpointTerraformType, name)

	utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
		{FieldName: "BgpAsNumber", APIField: &spProps.BgpAsNumber, TFValue: config.BgpAsNumber, IsConfigured: configuredAttrs.IsConfigured("bgp_as_number")},
		{FieldName: "ExpectedUplinkPort", APIField: &spProps.ExpectedUplinkPort, TFValue: config.ExpectedUplinkPort, IsConfigured: configuredAttrs.IsConfigured("expected_uplink_port")},
	})
	utils.SetNullableNumberFields([]utils.NullableNumberFieldMapping{
		{FieldName: "Position", APIField: &spProps.Position, TFValue: config.Position, IsConfigured: configuredAttrs.IsConfigured("position")},
		{FieldName: "RailGroup", APIField: &spProps.RailGroup, TFValue: config.RailGroup, IsConfigured: configuredAttrs.IsConfigured("rail_group")},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		configOp, objPropsCfg := utils.GetObjectPropertiesConfig(op, config.ObjectProperties, configuredAttrs)
		objProps := openapi.SwitchpointsPutRequestSwitchpointValueObjectProperties{}
		utils.SetObjectPropertiesFields([]utils.ObjectPropertiesField{
			{Name: "UserNotes", TFValue: op.UserNotes, APIValue: &objProps.UserNotes},
			{Name: "ExpectedParentEndpoint", TFValue: op.ExpectedParentEndpoint, APIValue: &objProps.ExpectedParentEndpoint},
			{Name: "Aggregate", TFValue: op.Aggregate, APIValue: &objProps.Aggregate},
			{Name: "IsHost", TFValue: op.IsHost, APIValue: &objProps.IsHost},
			{Name: "EmulateRfVideoPort", TFValue: op.EmulateRfVideoPort, APIValue: &objProps.EmulateRfVideoPort},
			{Name: "DrawAsEdgeDevice", TFValue: op.DrawAsEdgeDevice, APIValue: &objProps.DrawAsEdgeDevice},
		})
		utils.SetRefTypeFields([]utils.RefTypeFieldMapping{
			{
				FieldName:        "expected_parent_endpoint",
				RefTypeFieldName: "expected_parent_endpoint_ref_type_",
				APIField:         &objProps.ExpectedParentEndpoint,
				RefTypeAPIField:  &objProps.ExpectedParentEndpointRefType,
				TFValue:          op.ExpectedParentEndpoint,
				RefTypeTFValue:   op.ExpectedParentEndpointRefType,
			},
		})
		utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
			{FieldName: "NumberOfMultipoints", APIField: &objProps.NumberOfMultipoints, TFValue: configOp.NumberOfMultipoints, IsConfigured: objPropsCfg.IsFieldConfigured("number_of_multipoints")},
		})
		spProps.ObjectProperties = &objProps
	}

	// Handle badges
	if len(plan.Badges) > 0 {
		badges := make([]openapi.SwitchpointsPutRequestSwitchpointValueBadgesInner, len(plan.Badges))
		for i, badge := range plan.Badges {
			badgeItem := openapi.SwitchpointsPutRequestSwitchpointValueBadgesInner{}

			// Handle string fields
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Badge", APIField: &badgeItem.Badge, TFValue: badge.Badge},
				{FieldName: "BadgeRefType", APIField: &badgeItem.BadgeRefType, TFValue: badge.BadgeRefType},
			})

			// Handle int64 fields
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &badgeItem.Index, TFValue: badge.Index},
			})

			badges[i] = badgeItem
		}
		spProps.Badges = badges
	}

	// Handle children
	if len(plan.Children) > 0 {
		children := make([]openapi.SwitchpointsPutRequestSwitchpointValueChildrenInner, len(plan.Children))
		for i, child := range plan.Children {
			childItem := openapi.SwitchpointsPutRequestSwitchpointValueChildrenInner{}

			// Handle string fields
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "ChildNumEndpoint", APIField: &childItem.ChildNumEndpoint, TFValue: child.ChildNumEndpoint},
				{FieldName: "ChildNumEndpointRefType", APIField: &childItem.ChildNumEndpointRefType, TFValue: child.ChildNumEndpointRefType},
				{FieldName: "ChildNumDevice", APIField: &childItem.ChildNumDevice, TFValue: child.ChildNumDevice},
			})

			// Handle int64 fields
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &childItem.Index, TFValue: child.Index},
			})

			children[i] = childItem
		}
		spProps.Children = children
	}

	// Handle traffic mirrors
	if len(plan.TrafficMirrors) > 0 {
		mirrors := make([]openapi.SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner, len(plan.TrafficMirrors))
		for i, mirror := range plan.TrafficMirrors {
			mirrorItem := openapi.SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner{}

			// Handle boolean fields
			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "TrafficMirrorNumEnable", APIField: &mirrorItem.TrafficMirrorNumEnable, TFValue: mirror.TrafficMirrorNumEnable},
				{FieldName: "TrafficMirrorNumSourceLagIndicator", APIField: &mirrorItem.TrafficMirrorNumSourceLagIndicator, TFValue: mirror.TrafficMirrorNumSourceLagIndicator},
				{FieldName: "TrafficMirrorNumInboundTraffic", APIField: &mirrorItem.TrafficMirrorNumInboundTraffic, TFValue: mirror.TrafficMirrorNumInboundTraffic},
				{FieldName: "TrafficMirrorNumOutboundTraffic", APIField: &mirrorItem.TrafficMirrorNumOutboundTraffic, TFValue: mirror.TrafficMirrorNumOutboundTraffic},
			})

			// Handle string fields
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "TrafficMirrorNumSourcePort", APIField: &mirrorItem.TrafficMirrorNumSourcePort, TFValue: mirror.TrafficMirrorNumSourcePort},
				{FieldName: "TrafficMirrorNumDestinationPort", APIField: &mirrorItem.TrafficMirrorNumDestinationPort, TFValue: mirror.TrafficMirrorNumDestinationPort},
			})

			// Handle int64 fields
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &mirrorItem.Index, TFValue: mirror.Index},
			})

			mirrors[i] = mirrorItem
		}
		spProps.TrafficMirrors = mirrors
	}

	// Handle eths
	if len(plan.Eths) > 0 {
		eths := make([]openapi.SwitchpointsPutRequestSwitchpointValueEthsInner, len(plan.Eths))
		for i, eth := range plan.Eths {
			ethItem := openapi.SwitchpointsPutRequestSwitchpointValueEthsInner{}

			// Handle string fields
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Breakout", APIField: &ethItem.Breakout, TFValue: eth.Breakout},
				{FieldName: "CustomerVlan", APIField: &ethItem.CustomerVlan, TFValue: eth.CustomerVlan},
				{FieldName: "EthNumIcon", APIField: &ethItem.EthNumIcon, TFValue: eth.EthNumIcon},
				{FieldName: "EthNumLabel", APIField: &ethItem.EthNumLabel, TFValue: eth.EthNumLabel},
				{FieldName: "PortName", APIField: &ethItem.PortName, TFValue: eth.PortName},
			})

			// Handle boolean fields
			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &ethItem.Enable, TFValue: eth.Enable},
			})

			// Handle int64 fields
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &ethItem.Index, TFValue: eth.Index},
			})

			eths[i] = ethItem
		}
		spProps.Eths = eths
	}

	if len(plan.Pots) > 0 {
		pots := make([]openapi.SwitchpointsPutRequestSwitchpointValuePotsInner, len(plan.Pots))
		for i, pot := range plan.Pots {
			potItem := openapi.SwitchpointsPutRequestSwitchpointValuePotsInner{}

			// Handle string fields
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "PotsNumUri", APIField: &potItem.PotsNumUri, TFValue: pot.PotsNumUri},
				{FieldName: "PotsNumUsername", APIField: &potItem.PotsNumUsername, TFValue: pot.PotsNumUsername},
				{FieldName: "PotsNumPassword", APIField: &potItem.PotsNumPassword, TFValue: pot.PotsNumPassword},
				{FieldName: "PotsNumCallerId", APIField: &potItem.PotsNumCallerId, TFValue: pot.PotsNumCallerId},
				{FieldName: "PotsNumHotLine", APIField: &potItem.PotsNumHotLine, TFValue: pot.PotsNumHotLine},
				{FieldName: "PotsNumPasswordEncrypted", APIField: &potItem.PotsNumPasswordEncrypted, TFValue: pot.PotsNumPasswordEncrypted},
			})

			// Handle boolean fields
			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "PotsNumEnable", APIField: &potItem.PotsNumEnable, TFValue: pot.PotsNumEnable},
			})

			// Handle int64 fields
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &potItem.Index, TFValue: pot.Index},
			})
			pots[i] = potItem
		}
		spProps.Pots = pots
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "switchpoint", name, *spProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Switchpoint %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "switchpoints")

	var minState veritySwitchpointResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if switchpointData, exists := bulkMgr.GetResourceResponse("switchpoint", name); exists {
			effectiveData := utils.MergeMissingPlanScalars(switchpointData, plan, switchpointResourceType, r.provCtx.mode)
			state := populateSwitchpointState(ctx, minState, effectiveData, r.provCtx.mode)
			preserveSwitchpointPortNames(&state, &plan)
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

	postOpCtx := utils.WithPostOperationFallback(ctx, plan, switchpointResourceType, r.provCtx.mode)
	r.Read(postOpCtx, readReq, &readResp)
	if readResp.State.Raw.IsNull() {
		_, diags := utils.SetPostOperationFallbackState(postOpCtx, &readResp.State)
		readResp.Diagnostics.Append(diags...)
	}
	resp.State = readResp.State
	resp.Diagnostics = readResp.Diagnostics

	if resp.Diagnostics.HasError() {
		return
	}
	var readState veritySwitchpointResourceModel
	resp.Diagnostics.Append(readResp.State.Get(ctx, &readState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	preserveSwitchpointPortNames(&readState, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &readState)...)
}

func (r *veritySwitchpointResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state veritySwitchpointResourceModel
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

	spName := state.Name.ValueString()
	priorState := state // save prior state to preserve reference-only fields

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if switchpointData, exists := r.bulkOpsMgr.GetResourceResponse("switchpoint", spName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached switchpoint data for %s from recent operation", spName))
			state = populateSwitchpointState(ctx, state, utils.ApplyPostOperationFallback(ctx, switchpointData), r.provCtx.mode)
			preserveSwitchpointPortNames(&state, &priorState)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("switchpoint") {
		tflog.Info(ctx, fmt.Sprintf("Skipping switchpoint %s verification – trusting recent successful API operation", spName))
		if handled, diags := utils.SetPostOperationFallbackState(ctx, &resp.State); handled {
			resp.Diagnostics.Append(diags...)
		}
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching switchpoints for verification of %s", spName))

	type SwitchpointResponse struct {
		Switchpoint map[string]interface{} `json:"switchpoint"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "switchpoints", spName,
		func() (SwitchpointResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch switchpoints")
			respAPI, err := r.client.SwitchpointsAPI.SwitchpointsGet(ctx).Execute()
			if err != nil {
				return SwitchpointResponse{}, fmt.Errorf("error reading switchpoints: %v", err)
			}
			defer respAPI.Body.Close()

			var res SwitchpointResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return SwitchpointResponse{}, fmt.Errorf("failed to decode switchpoints response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d switchpoints", len(res.Switchpoint)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Switchpoint %s", spName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for switchpoint with name: %s", spName))

	spData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.Switchpoint,
		spName,
		func(data interface{}) (string, bool) {
			if switchpoint, ok := data.(map[string]interface{}); ok {
				if name, ok := switchpoint["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Switchpoint with name '%s' not found in API response", spName))
		resp.State.RemoveResource(ctx)
		return
	}

	switchpointMap, ok := spData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Switchpoint Data",
			fmt.Sprintf("Switchpoint data is not in expected format for %s", spName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found switchpoint '%s' under API key '%s'", spName, actualAPIName))

	state = populateSwitchpointState(ctx, state, utils.ApplyPostOperationFallback(ctx, switchpointMap), r.provCtx.mode)
	preserveSwitchpointPortNames(&state, &priorState)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *veritySwitchpointResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state veritySwitchpointResourceModel

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

	// Validate auto-assigned fields - these checks prevent ineffective API calls
	// Only error if the auto-assigned flag is enabled AND the user is explicitly setting a value
	// AND the auto-assigned flag itself is not changing (which would be a valid operation)
	// Don't error if the field is unknown (computed during plan recalculation)
	if !plan.BgpAsNumber.Equal(state.BgpAsNumber) &&
		!plan.BgpAsNumber.IsNull() && !plan.BgpAsNumber.IsUnknown() && // User is explicitly setting a value
		!plan.BgpAsNumberAutoAssigned.IsNull() && plan.BgpAsNumberAutoAssigned.ValueBool() &&
		plan.BgpAsNumberAutoAssigned.Equal(state.BgpAsNumberAutoAssigned) {
		resp.Diagnostics.AddError(
			"Cannot modify auto-assigned field",
			"The 'bgp_as_number' field cannot be modified because 'bgp_as_number_auto_assigned_' is set to true.",
		)
		return
	}

	if !plan.SwitchRouterIdIpMask.Equal(state.SwitchRouterIdIpMask) &&
		!plan.SwitchRouterIdIpMask.IsNull() && !plan.SwitchRouterIdIpMask.IsUnknown() && // User is explicitly setting a value
		!plan.SwitchRouterIdIpMaskAutoAssigned.IsNull() && plan.SwitchRouterIdIpMaskAutoAssigned.ValueBool() &&
		plan.SwitchRouterIdIpMaskAutoAssigned.Equal(state.SwitchRouterIdIpMaskAutoAssigned) {
		resp.Diagnostics.AddError(
			"Cannot modify auto-assigned field",
			"The 'switch_router_id_ip_mask' field cannot be modified because 'switch_router_id_ip_mask_auto_assigned_' is set to true.",
		)
		return
	}

	if !plan.SwitchVtepIdIpMask.Equal(state.SwitchVtepIdIpMask) &&
		!plan.SwitchVtepIdIpMask.IsNull() && !plan.SwitchVtepIdIpMask.IsUnknown() && // User is explicitly setting a value
		!plan.SwitchVtepIdIpMaskAutoAssigned.IsNull() && plan.SwitchVtepIdIpMaskAutoAssigned.ValueBool() &&
		plan.SwitchVtepIdIpMaskAutoAssigned.Equal(state.SwitchVtepIdIpMaskAutoAssigned) {
		resp.Diagnostics.AddError(
			"Cannot modify auto-assigned field",
			"The 'switch_vtep_id_ip_mask' field cannot be modified because 'switch_vtep_id_ip_mask_auto_assigned_' is set to true.",
		)
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
	spProps := openapi.SwitchpointsPutRequestSwitchpointValue{}
	hasChanges := false

	// Get config for nullable field handling
	var config veritySwitchpointResourceModel
	req.Config.Get(ctx, &config)
	workDir := r.provCtx.workDir
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, switchpointTerraformType, name)

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { spProps.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.DeviceSerialNumber, state.DeviceSerialNumber, func(v *string) { spProps.DeviceSerialNumber = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.ExpectedSite, state.ExpectedSite, func(v *string) { spProps.ExpectedSite = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Type, state.Type, func(v *string) { spProps.Type = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.RackInfo, state.RackInfo, func(v *string) { spProps.RackInfo = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.PasswordEncrypted, state.PasswordEncrypted, func(v *string) { spProps.PasswordEncrypted = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.EnablePasswordEncrypted, state.EnablePasswordEncrypted, func(v *string) { spProps.EnablePasswordEncrypted = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.SshKeyOrPasswordEncrypted, state.SshKeyOrPasswordEncrypted, func(v *string) { spProps.SshKeyOrPasswordEncrypted = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.PassphraseEncrypted, state.PassphraseEncrypted, func(v *string) { spProps.PassphraseEncrypted = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.PrivatePasswordEncrypted, state.PrivatePasswordEncrypted, func(v *string) { spProps.PrivatePasswordEncrypted = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.IpSource, state.IpSource, func(v *string) { spProps.IpSource = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Gateway, state.Gateway, func(v *string) { spProps.Gateway = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.SwitchGateway, state.SwitchGateway, func(v *string) { spProps.SwitchGateway = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.CommType, state.CommType, func(v *string) { spProps.CommType = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.SnmpCommunityString, state.SnmpCommunityString, func(v *string) { spProps.SnmpCommunityString = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.UplinkPort, state.UplinkPort, func(v *string) { spProps.UplinkPort = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.ExpectedBreakout, state.ExpectedBreakout, func(v *string) { spProps.ExpectedBreakout = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.LldpSearchString, state.LldpSearchString, func(v *string) { spProps.LldpSearchString = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.ZtpIdentification, state.ZtpIdentification, func(v *string) { spProps.ZtpIdentification = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.LocatedBy, state.LocatedBy, func(v *string) { spProps.LocatedBy = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.PowerState, state.PowerState, func(v *string) { spProps.PowerState = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.CommunicationMode, state.CommunicationMode, func(v *string) { spProps.CommunicationMode = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.CliAccessMode, state.CliAccessMode, func(v *string) { spProps.CliAccessMode = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Username, state.Username, func(v *string) { spProps.Username = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Password, state.Password, func(v *string) { spProps.Password = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.EnablePassword, state.EnablePassword, func(v *string) { spProps.EnablePassword = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.SshKeyOrPassword, state.SshKeyOrPassword, func(v *string) { spProps.SshKeyOrPassword = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Sdlc, state.Sdlc, func(v *string) { spProps.Sdlc = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.SecurityType, state.SecurityType, func(v *string) { spProps.SecurityType = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Snmpv3Username, state.Snmpv3Username, func(v *string) { spProps.Snmpv3Username = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.AuthenticationProtocol, state.AuthenticationProtocol, func(v *string) { spProps.AuthenticationProtocol = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Passphrase, state.Passphrase, func(v *string) { spProps.Passphrase = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.PrivateProtocol, state.PrivateProtocol, func(v *string) { spProps.PrivateProtocol = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.PrivatePassword, state.PrivatePassword, func(v *string) { spProps.PrivatePassword = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.DeviceManagedAs, state.DeviceManagedAs, func(v *string) { spProps.DeviceManagedAs = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Port, state.Port, func(v *string) { spProps.Port = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { spProps.Enable = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.IsTopOfIsland, state.IsTopOfIsland, func(v *bool) { spProps.IsTopOfIsland = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.ReadOnlyMode, state.ReadOnlyMode, func(v *bool) { spProps.ReadOnlyMode = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.Locked, state.Locked, func(v *bool) { spProps.Locked = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.OutOfBandManagement, state.OutOfBandManagement, func(v *bool) { spProps.OutOfBandManagement = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.BbSwitch, state.BbSwitch, func(v *bool) { spProps.BbSwitch = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.ManagedOnNativeVlan, state.ManagedOnNativeVlan, func(v *bool) { spProps.ManagedOnNativeVlan = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.IsFabric, state.IsFabric, func(v *bool) { spProps.IsFabric = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.UsesTaggedPackets, state.UsesTaggedPackets, func(v *bool) { spProps.UsesTaggedPackets = v }, &hasChanges)

	// Handle ConnectedBundle and ConnectedBundleRefType using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.ConnectedBundle, state.ConnectedBundle, plan.ConnectedBundleRefType, state.ConnectedBundleRefType,
		func(val *string) { spProps.ConnectedBundle = val },
		func(val *string) { spProps.ConnectedBundleRefType = val },
		"connected_bundle", "connected_bundle_ref_type_",
		&hasChanges, &resp.Diagnostics,
	) {
		return
	}

	// Handle Plane and PlaneRefType using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.Plane, state.Plane, plan.PlaneRefType, state.PlaneRefType,
		func(v *string) { spProps.Plane = v },
		func(v *string) { spProps.PlaneRefType = v },
		"plane", "plane_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	// Handle SpinePlane and SpinePlaneRefType using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.SpinePlane, state.SpinePlane, plan.SpinePlaneRefType, state.SpinePlaneRefType,
		func(v *string) { spProps.SpinePlane = v },
		func(v *string) { spProps.SpinePlaneRefType = v },
		"spine_plane", "spine_plane_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	// Handle Pod and PodRefType using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.Pod, state.Pod, plan.PodRefType, state.PodRefType,
		func(v *string) { spProps.Pod = v },
		func(v *string) { spProps.PodRefType = v },
		"pod", "pod_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	if !utils.HandleOneRefTypeSupported(
		plan.Tenant, state.Tenant, plan.TenantRefType, state.TenantRefType,
		func(v *string) { spProps.Tenant = v },
		func(v *string) { spProps.TenantRefType = v },
		"tenant", "tenant_ref_type_",
		&hasChanges, &resp.Diagnostics,
	) {
		return
	}

	if !utils.HandleOneRefTypeSupported(
		plan.ExpectedSite, state.ExpectedSite, plan.ExpectedSiteRefType, state.ExpectedSiteRefType,
		func(v *string) { spProps.ExpectedSite = v },
		func(v *string) { spProps.ExpectedSiteRefType = v },
		"expected_site", "expected_site_ref_type_",
		&hasChanges, &resp.Diagnostics,
	) {
		return
	}

	if !utils.HandleOneRefTypeSupported(
		plan.Su, state.Su, plan.SuRefType, state.SuRefType,
		func(v *string) { spProps.Su = v },
		func(v *string) { spProps.SuRefType = v },
		"su", "su_ref_type_",
		&hasChanges, &resp.Diagnostics,
	) {
		return
	}

	if !utils.HandleOneRefTypeSupported(
		plan.SspGroup, state.SspGroup, plan.SspGroupRefType, state.SspGroupRefType,
		func(v *string) { spProps.SspGroup = v },
		func(v *string) { spProps.SspGroupRefType = v },
		"ssp_group", "ssp_group_ref_type_",
		&hasChanges, &resp.Diagnostics,
	) {
		return
	}

	if !utils.HandleOneRefTypeSupported(
		plan.Rack, state.Rack, plan.RackRefType, state.RackRefType,
		func(v *string) { spProps.Rack = v },
		func(v *string) { spProps.RackRefType = v },
		"rack", "rack_ref_type_",
		&hasChanges, &resp.Diagnostics,
	) {
		return
	}

	if !utils.HandleOneRefTypeSupported(
		plan.Switch, state.Switch, plan.SwitchRefType, state.SwitchRefType,
		func(v *string) { spProps.Switch = v },
		func(v *string) { spProps.SwitchRefType = v },
		"switch", "switch_ref_type_",
		&hasChanges, &resp.Diagnostics,
	) {
		return
	}

	if !utils.HandleOneRefTypeSupported(
		plan.ConnectionService, state.ConnectionService, plan.ConnectionServiceRefType, state.ConnectionServiceRefType,
		func(v *string) { spProps.ConnectionService = v },
		func(v *string) { spProps.ConnectionServiceRefType = v },
		"connection_service", "connection_service_ref_type_",
		&hasChanges, &resp.Diagnostics,
	) {
		return
	}

	utils.CompareAndSetNullableNumberField(config.Position, state.Position, configuredAttrs.IsConfigured("position"), func(v *openapi.NullableFloat64) { spProps.Position = *v }, &hasChanges)
	utils.CompareAndSetNullableNumberField(config.RailGroup, state.RailGroup, configuredAttrs.IsConfigured("rail_group"), func(v *openapi.NullableFloat64) { spProps.RailGroup = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.ExpectedUplinkPort, state.ExpectedUplinkPort, configuredAttrs.IsConfigured("expected_uplink_port"), func(v *openapi.NullableInt64) { spProps.ExpectedUplinkPort = *v }, &hasChanges)

	// Handle BgpAsNumber and BgpAsNumberAutoAssigned changes
	bgpAsNumberChanged := !plan.BgpAsNumber.IsUnknown() && !plan.BgpAsNumber.Equal(state.BgpAsNumber)
	bgpAsNumberAutoAssignedChanged := !plan.BgpAsNumberAutoAssigned.Equal(state.BgpAsNumberAutoAssigned)

	if bgpAsNumberChanged || bgpAsNumberAutoAssignedChanged {
		// Handle BgpAsNumber field changes
		if bgpAsNumberChanged {
			if !plan.BgpAsNumber.IsNull() {
				val := plan.BgpAsNumber.ValueInt64()
				spProps.BgpAsNumber = *openapi.NewNullableInt64(&val)
			} else {
				spProps.BgpAsNumber = *openapi.NewNullableInt64(nil)
			}
		}

		// Handle BgpAsNumberAutoAssigned field changes
		if bgpAsNumberAutoAssignedChanged {
			// Only send bgp_as_number_auto_assigned_ if the user has explicitly specified it in their configuration
			var config veritySwitchpointResourceModel
			userSpecifiedBgpAsNumberAutoAssigned := false
			if !req.Config.Raw.IsNull() {
				if err := req.Config.Get(ctx, &config); err == nil {
					userSpecifiedBgpAsNumberAutoAssigned = !config.BgpAsNumberAutoAssigned.IsNull()
				}
			}

			if userSpecifiedBgpAsNumberAutoAssigned {
				spProps.BgpAsNumberAutoAssigned = openapi.PtrBool(plan.BgpAsNumberAutoAssigned.ValueBool())

				// Special case: When changing from auto-assigned (true) to manual (false),
				// the API requires both bgp_as_number_auto_assigned_ and bgp_as_number fields to be sent.
				if !state.BgpAsNumberAutoAssigned.IsNull() && state.BgpAsNumberAutoAssigned.ValueBool() &&
					!plan.BgpAsNumberAutoAssigned.ValueBool() {
					// Changing from auto-assigned=true to auto-assigned=false
					// Must include BgpAsNumber value in the request for the change to take effect
					if !plan.BgpAsNumber.IsNull() {
						val := plan.BgpAsNumber.ValueInt64()
						spProps.BgpAsNumber = *openapi.NewNullableInt64(&val)
					} else if !state.BgpAsNumber.IsNull() {
						// Use current state BgpAsNumber if plan doesn't specify one
						val := state.BgpAsNumber.ValueInt64()
						spProps.BgpAsNumber = *openapi.NewNullableInt64(&val)
					}
				}
			}
		} else if bgpAsNumberChanged {
			// BgpAsNumber changed but BgpAsNumberAutoAssigned didn't change
			// Send the auto-assigned flag to maintain consistency with API
			if !plan.BgpAsNumberAutoAssigned.IsNull() {
				spProps.BgpAsNumberAutoAssigned = openapi.PtrBool(plan.BgpAsNumberAutoAssigned.ValueBool())
			} else if !state.BgpAsNumberAutoAssigned.IsNull() {
				spProps.BgpAsNumberAutoAssigned = openapi.PtrBool(state.BgpAsNumberAutoAssigned.ValueBool())
			} else {
				spProps.BgpAsNumberAutoAssigned = openapi.PtrBool(false)
			}
		}

		hasChanges = true
	}

	// Handle SwitchRouterIdIpMask and SwitchRouterIdIpMaskAutoAssigned changes
	switchRouterIdIpMaskChanged := !plan.SwitchRouterIdIpMask.IsUnknown() && !plan.SwitchRouterIdIpMask.Equal(state.SwitchRouterIdIpMask)
	switchRouterIdIpMaskAutoAssignedChanged := !plan.SwitchRouterIdIpMaskAutoAssigned.Equal(state.SwitchRouterIdIpMaskAutoAssigned)

	if switchRouterIdIpMaskChanged || switchRouterIdIpMaskAutoAssignedChanged {
		// Handle SwitchRouterIdIpMask field changes
		if switchRouterIdIpMaskChanged {
			spProps.SwitchRouterIdIpMask = openapi.PtrString(plan.SwitchRouterIdIpMask.ValueString())
		}

		// Handle SwitchRouterIdIpMaskAutoAssigned field changes
		if switchRouterIdIpMaskAutoAssignedChanged {
			// Only send switch_router_id_ip_mask_auto_assigned_ if the user has explicitly specified it in their configuration
			var config veritySwitchpointResourceModel
			userSpecifiedSwitchRouterIdIpMaskAutoAssigned := false
			if !req.Config.Raw.IsNull() {
				if err := req.Config.Get(ctx, &config); err == nil {
					userSpecifiedSwitchRouterIdIpMaskAutoAssigned = !config.SwitchRouterIdIpMaskAutoAssigned.IsNull()
				}
			}

			if userSpecifiedSwitchRouterIdIpMaskAutoAssigned {
				spProps.SwitchRouterIdIpMaskAutoAssigned = openapi.PtrBool(plan.SwitchRouterIdIpMaskAutoAssigned.ValueBool())

				// Special case: When changing from auto-assigned (true) to manual (false),
				// the API requires both switch_router_id_ip_mask_auto_assigned_ and switch_router_id_ip_mask fields to be sent.
				if !state.SwitchRouterIdIpMaskAutoAssigned.IsNull() && state.SwitchRouterIdIpMaskAutoAssigned.ValueBool() &&
					!plan.SwitchRouterIdIpMaskAutoAssigned.ValueBool() {
					// Changing from auto-assigned=true to auto-assigned=false
					// Must include SwitchRouterIdIpMask value in the request for the change to take effect
					if !plan.SwitchRouterIdIpMask.IsNull() {
						spProps.SwitchRouterIdIpMask = openapi.PtrString(plan.SwitchRouterIdIpMask.ValueString())
					} else if !state.SwitchRouterIdIpMask.IsNull() {
						// Use current state SwitchRouterIdIpMask if plan doesn't specify one
						spProps.SwitchRouterIdIpMask = openapi.PtrString(state.SwitchRouterIdIpMask.ValueString())
					}
				}
			}
		} else if switchRouterIdIpMaskChanged {
			// SwitchRouterIdIpMask changed but SwitchRouterIdIpMaskAutoAssigned didn't change
			// Send the auto-assigned flag to maintain consistency with API
			if !plan.SwitchRouterIdIpMaskAutoAssigned.IsNull() {
				spProps.SwitchRouterIdIpMaskAutoAssigned = openapi.PtrBool(plan.SwitchRouterIdIpMaskAutoAssigned.ValueBool())
			} else if !state.SwitchRouterIdIpMaskAutoAssigned.IsNull() {
				spProps.SwitchRouterIdIpMaskAutoAssigned = openapi.PtrBool(state.SwitchRouterIdIpMaskAutoAssigned.ValueBool())
			} else {
				spProps.SwitchRouterIdIpMaskAutoAssigned = openapi.PtrBool(false)
			}
		}

		hasChanges = true
	}

	// Handle SwitchVtepIdIpMask and SwitchVtepIdIpMaskAutoAssigned changes
	switchVtepIdIpMaskChanged := !plan.SwitchVtepIdIpMask.IsUnknown() && !plan.SwitchVtepIdIpMask.Equal(state.SwitchVtepIdIpMask)
	switchVtepIdIpMaskAutoAssignedChanged := !plan.SwitchVtepIdIpMaskAutoAssigned.Equal(state.SwitchVtepIdIpMaskAutoAssigned)

	if switchVtepIdIpMaskChanged || switchVtepIdIpMaskAutoAssignedChanged {
		// Handle SwitchVtepIdIpMask field changes
		if switchVtepIdIpMaskChanged {
			spProps.SwitchVtepIdIpMask = openapi.PtrString(plan.SwitchVtepIdIpMask.ValueString())
		}

		// Handle SwitchVtepIdIpMaskAutoAssigned field changes
		if switchVtepIdIpMaskAutoAssignedChanged {
			// Only send switch_vtep_id_ip_mask_auto_assigned_ if the user has explicitly specified it in their configuration
			var config veritySwitchpointResourceModel
			userSpecifiedSwitchVtepIdIpMaskAutoAssigned := false
			if !req.Config.Raw.IsNull() {
				if err := req.Config.Get(ctx, &config); err == nil {
					userSpecifiedSwitchVtepIdIpMaskAutoAssigned = !config.SwitchVtepIdIpMaskAutoAssigned.IsNull()
				}
			}

			if userSpecifiedSwitchVtepIdIpMaskAutoAssigned {
				spProps.SwitchVtepIdIpMaskAutoAssigned = openapi.PtrBool(plan.SwitchVtepIdIpMaskAutoAssigned.ValueBool())

				// Special case: When changing from auto-assigned (true) to manual (false),
				// the API requires both switch_vtep_id_ip_mask_auto_assigned_ and switch_vtep_id_ip_mask fields to be sent.
				if !state.SwitchVtepIdIpMaskAutoAssigned.IsNull() && state.SwitchVtepIdIpMaskAutoAssigned.ValueBool() &&
					!plan.SwitchVtepIdIpMaskAutoAssigned.ValueBool() {
					// Changing from auto-assigned=true to auto-assigned=false
					// Must include SwitchVtepIdIpMask value in the request for the change to take effect
					if !plan.SwitchVtepIdIpMask.IsNull() {
						spProps.SwitchVtepIdIpMask = openapi.PtrString(plan.SwitchVtepIdIpMask.ValueString())
					} else if !state.SwitchVtepIdIpMask.IsNull() {
						// Use current state SwitchVtepIdIpMask if plan doesn't specify one
						spProps.SwitchVtepIdIpMask = openapi.PtrString(state.SwitchVtepIdIpMask.ValueString())
					}
				}
			}
		} else if switchVtepIdIpMaskChanged {
			// SwitchVtepIdIpMask changed but SwitchVtepIdIpMaskAutoAssigned didn't change
			// Send the auto-assigned flag to maintain consistency with API
			if !plan.SwitchVtepIdIpMaskAutoAssigned.IsNull() {
				spProps.SwitchVtepIdIpMaskAutoAssigned = openapi.PtrBool(plan.SwitchVtepIdIpMaskAutoAssigned.ValueBool())
			} else if !state.SwitchVtepIdIpMaskAutoAssigned.IsNull() {
				spProps.SwitchVtepIdIpMaskAutoAssigned = openapi.PtrBool(state.SwitchVtepIdIpMaskAutoAssigned.ValueBool())
			} else {
				spProps.SwitchVtepIdIpMaskAutoAssigned = openapi.PtrBool(false)
			}
		}

		hasChanges = true
	}

	// Handle ControllerIpAndMask and ControllerIpAndMaskAutoAssigned changes
	controllerIpAndMaskChanged := !plan.ControllerIpAndMask.IsUnknown() && !plan.ControllerIpAndMask.Equal(state.ControllerIpAndMask)
	controllerIpAndMaskAutoAssignedChanged := !plan.ControllerIpAndMaskAutoAssigned.Equal(state.ControllerIpAndMaskAutoAssigned)

	if controllerIpAndMaskChanged || controllerIpAndMaskAutoAssignedChanged {
		if controllerIpAndMaskChanged {
			spProps.ControllerIpAndMask = openapi.PtrString(plan.ControllerIpAndMask.ValueString())
		}

		if controllerIpAndMaskAutoAssignedChanged {
			var config veritySwitchpointResourceModel
			userSpecifiedControllerIpAndMaskAutoAssigned := false
			if !req.Config.Raw.IsNull() {
				if err := req.Config.Get(ctx, &config); err == nil {
					userSpecifiedControllerIpAndMaskAutoAssigned = !config.ControllerIpAndMaskAutoAssigned.IsNull()
				}
			}

			if userSpecifiedControllerIpAndMaskAutoAssigned {
				spProps.ControllerIpAndMaskAutoAssigned = openapi.PtrBool(plan.ControllerIpAndMaskAutoAssigned.ValueBool())

				if !state.ControllerIpAndMaskAutoAssigned.IsNull() && state.ControllerIpAndMaskAutoAssigned.ValueBool() &&
					!plan.ControllerIpAndMaskAutoAssigned.ValueBool() {
					if !plan.ControllerIpAndMask.IsNull() {
						spProps.ControllerIpAndMask = openapi.PtrString(plan.ControllerIpAndMask.ValueString())
					} else if !state.ControllerIpAndMask.IsNull() {
						spProps.ControllerIpAndMask = openapi.PtrString(state.ControllerIpAndMask.ValueString())
					}
				}
			}
		} else if controllerIpAndMaskChanged {
			if !plan.ControllerIpAndMaskAutoAssigned.IsNull() {
				spProps.ControllerIpAndMaskAutoAssigned = openapi.PtrBool(plan.ControllerIpAndMaskAutoAssigned.ValueBool())
			} else if !state.ControllerIpAndMaskAutoAssigned.IsNull() {
				spProps.ControllerIpAndMaskAutoAssigned = openapi.PtrBool(state.ControllerIpAndMaskAutoAssigned.ValueBool())
			} else {
				spProps.ControllerIpAndMaskAutoAssigned = openapi.PtrBool(false)
			}
		}

		hasChanges = true
	}

	// Handle SwitchIpAndMask and SwitchIpAndMaskAutoAssigned changes
	switchIpAndMaskChanged := !plan.SwitchIpAndMask.IsUnknown() && !plan.SwitchIpAndMask.Equal(state.SwitchIpAndMask)
	switchIpAndMaskAutoAssignedChanged := !plan.SwitchIpAndMaskAutoAssigned.Equal(state.SwitchIpAndMaskAutoAssigned)

	if switchIpAndMaskChanged || switchIpAndMaskAutoAssignedChanged {
		if switchIpAndMaskChanged {
			spProps.SwitchIpAndMask = openapi.PtrString(plan.SwitchIpAndMask.ValueString())
		}

		if switchIpAndMaskAutoAssignedChanged {
			var config veritySwitchpointResourceModel
			userSpecifiedSwitchIpAndMaskAutoAssigned := false
			if !req.Config.Raw.IsNull() {
				if err := req.Config.Get(ctx, &config); err == nil {
					userSpecifiedSwitchIpAndMaskAutoAssigned = !config.SwitchIpAndMaskAutoAssigned.IsNull()
				}
			}

			if userSpecifiedSwitchIpAndMaskAutoAssigned {
				spProps.SwitchIpAndMaskAutoAssigned = openapi.PtrBool(plan.SwitchIpAndMaskAutoAssigned.ValueBool())

				if !state.SwitchIpAndMaskAutoAssigned.IsNull() && state.SwitchIpAndMaskAutoAssigned.ValueBool() &&
					!plan.SwitchIpAndMaskAutoAssigned.ValueBool() {
					if !plan.SwitchIpAndMask.IsNull() {
						spProps.SwitchIpAndMask = openapi.PtrString(plan.SwitchIpAndMask.ValueString())
					} else if !state.SwitchIpAndMask.IsNull() {
						spProps.SwitchIpAndMask = openapi.PtrString(state.SwitchIpAndMask.ValueString())
					}
				}
			}
		} else if switchIpAndMaskChanged {
			if !plan.SwitchIpAndMaskAutoAssigned.IsNull() {
				spProps.SwitchIpAndMaskAutoAssigned = openapi.PtrBool(plan.SwitchIpAndMaskAutoAssigned.ValueBool())
			} else if !state.SwitchIpAndMaskAutoAssigned.IsNull() {
				spProps.SwitchIpAndMaskAutoAssigned = openapi.PtrBool(state.SwitchIpAndMaskAutoAssigned.ValueBool())
			} else {
				spProps.SwitchIpAndMaskAutoAssigned = openapi.PtrBool(false)
			}
		}

		hasChanges = true
	}

	// Handle badges
	changedBadges, badgesChanged := utils.ProcessIndexedArrayUpdates(plan.Badges, state.Badges,
		utils.IndexedItemHandler[veritySwitchpointBadgeModel, openapi.SwitchpointsPutRequestSwitchpointValueBadgesInner]{
			CreateNew: func(planItem veritySwitchpointBadgeModel) openapi.SwitchpointsPutRequestSwitchpointValueBadgesInner {
				badge := openapi.SwitchpointsPutRequestSwitchpointValueBadgesInner{}

				// Handle string fields
				utils.SetStringFields([]utils.StringFieldMapping{
					{FieldName: "Badge", APIField: &badge.Badge, TFValue: planItem.Badge},
					{FieldName: "BadgeRefType", APIField: &badge.BadgeRefType, TFValue: planItem.BadgeRefType},
				})

				// Handle int64 fields
				utils.SetInt64Fields([]utils.Int64FieldMapping{
					{FieldName: "Index", APIField: &badge.Index, TFValue: planItem.Index},
				})

				return badge
			},
			UpdateExisting: func(planItem veritySwitchpointBadgeModel, stateItem veritySwitchpointBadgeModel) (openapi.SwitchpointsPutRequestSwitchpointValueBadgesInner, bool) {
				badge := openapi.SwitchpointsPutRequestSwitchpointValueBadgesInner{}
				fieldChanged := false

				// Handle badge and badge_ref_type_ using "One ref type supported" pattern
				if !utils.HandleOneRefTypeSupported(
					planItem.Badge, stateItem.Badge, planItem.BadgeRefType, stateItem.BadgeRefType,
					func(v *string) { badge.Badge = v },
					func(v *string) { badge.BadgeRefType = v },
					"badge", "badge_ref_type_",
					&fieldChanged,
					&resp.Diagnostics,
				) {
					return badge, false
				}

				// Always include index — API requires it to identify which array element to modify
				utils.SetInt64Fields([]utils.Int64FieldMapping{
					{FieldName: "Index", APIField: &badge.Index, TFValue: planItem.Index},
				})

				return badge, fieldChanged
			},
			CreateDeleted: func(index int64) openapi.SwitchpointsPutRequestSwitchpointValueBadgesInner {
				return openapi.SwitchpointsPutRequestSwitchpointValueBadgesInner{
					Index: openapi.PtrInt64(int64(index)),
				}
			},
		})
	if badgesChanged {
		spProps.Badges = changedBadges
		hasChanges = true
	}

	// Handle children
	changedChildren, childrenChanged := utils.ProcessIndexedArrayUpdates(plan.Children, state.Children,
		utils.IndexedItemHandler[veritySwitchpointChildModel, openapi.SwitchpointsPutRequestSwitchpointValueChildrenInner]{
			CreateNew: func(planItem veritySwitchpointChildModel) openapi.SwitchpointsPutRequestSwitchpointValueChildrenInner {
				child := openapi.SwitchpointsPutRequestSwitchpointValueChildrenInner{}

				// Handle string fields
				utils.SetStringFields([]utils.StringFieldMapping{
					{FieldName: "ChildNumEndpoint", APIField: &child.ChildNumEndpoint, TFValue: planItem.ChildNumEndpoint},
					{FieldName: "ChildNumEndpointRefType", APIField: &child.ChildNumEndpointRefType, TFValue: planItem.ChildNumEndpointRefType},
					{FieldName: "ChildNumDevice", APIField: &child.ChildNumDevice, TFValue: planItem.ChildNumDevice},
				})

				// Handle int64 fields
				utils.SetInt64Fields([]utils.Int64FieldMapping{
					{FieldName: "Index", APIField: &child.Index, TFValue: planItem.Index},
				})

				return child
			},
			UpdateExisting: func(planItem veritySwitchpointChildModel, stateItem veritySwitchpointChildModel) (openapi.SwitchpointsPutRequestSwitchpointValueChildrenInner, bool) {
				child := openapi.SwitchpointsPutRequestSwitchpointValueChildrenInner{}
				fieldChanged := false

				// Handle child_num_endpoint and child_num_endpoint_ref_type_ using "One ref type supported" pattern
				if !utils.HandleOneRefTypeSupported(
					planItem.ChildNumEndpoint, stateItem.ChildNumEndpoint, planItem.ChildNumEndpointRefType, stateItem.ChildNumEndpointRefType,
					func(v *string) { child.ChildNumEndpoint = v },
					func(v *string) { child.ChildNumEndpointRefType = v },
					"child_num_endpoint", "child_num_endpoint_ref_type_",
					&fieldChanged,
					&resp.Diagnostics,
				) {
					return child, false
				}

				// Handle other string field changes
				utils.CompareAndSetStringField(planItem.ChildNumDevice, stateItem.ChildNumDevice, func(v *string) { child.ChildNumDevice = v }, &fieldChanged)

				// Always include index — API requires it to identify which array element to modify
				utils.SetInt64Fields([]utils.Int64FieldMapping{
					{FieldName: "Index", APIField: &child.Index, TFValue: planItem.Index},
				})

				return child, fieldChanged
			},
			CreateDeleted: func(index int64) openapi.SwitchpointsPutRequestSwitchpointValueChildrenInner {
				return openapi.SwitchpointsPutRequestSwitchpointValueChildrenInner{
					Index: openapi.PtrInt64(int64(index)),
				}
			},
		})
	if childrenChanged {
		spProps.Children = changedChildren
		hasChanges = true
	}

	// Handle traffic mirrors
	changedTrafficMirrors, trafficMirrorsChanged := utils.ProcessIndexedArrayUpdates(plan.TrafficMirrors, state.TrafficMirrors,
		utils.IndexedItemHandler[veritySwitchpointTrafficMirrorModel, openapi.SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner]{
			CreateNew: func(planItem veritySwitchpointTrafficMirrorModel) openapi.SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner {
				mirror := openapi.SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner{}

				// Handle boolean fields
				utils.SetBoolFields([]utils.BoolFieldMapping{
					{FieldName: "TrafficMirrorNumEnable", APIField: &mirror.TrafficMirrorNumEnable, TFValue: planItem.TrafficMirrorNumEnable},
					{FieldName: "TrafficMirrorNumSourceLagIndicator", APIField: &mirror.TrafficMirrorNumSourceLagIndicator, TFValue: planItem.TrafficMirrorNumSourceLagIndicator},
					{FieldName: "TrafficMirrorNumInboundTraffic", APIField: &mirror.TrafficMirrorNumInboundTraffic, TFValue: planItem.TrafficMirrorNumInboundTraffic},
					{FieldName: "TrafficMirrorNumOutboundTraffic", APIField: &mirror.TrafficMirrorNumOutboundTraffic, TFValue: planItem.TrafficMirrorNumOutboundTraffic},
				})

				// Handle string fields
				utils.SetStringFields([]utils.StringFieldMapping{
					{FieldName: "TrafficMirrorNumSourcePort", APIField: &mirror.TrafficMirrorNumSourcePort, TFValue: planItem.TrafficMirrorNumSourcePort},
					{FieldName: "TrafficMirrorNumDestinationPort", APIField: &mirror.TrafficMirrorNumDestinationPort, TFValue: planItem.TrafficMirrorNumDestinationPort},
				})

				// Handle int64 fields
				utils.SetInt64Fields([]utils.Int64FieldMapping{
					{FieldName: "Index", APIField: &mirror.Index, TFValue: planItem.Index},
				})

				return mirror
			},
			UpdateExisting: func(planItem veritySwitchpointTrafficMirrorModel, stateItem veritySwitchpointTrafficMirrorModel) (openapi.SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner, bool) {
				mirror := openapi.SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner{}
				fieldChanged := false

				// Handle boolean field changes
				utils.CompareAndSetBoolField(planItem.TrafficMirrorNumEnable, stateItem.TrafficMirrorNumEnable, func(v *bool) { mirror.TrafficMirrorNumEnable = v }, &fieldChanged)
				utils.CompareAndSetBoolField(planItem.TrafficMirrorNumSourceLagIndicator, stateItem.TrafficMirrorNumSourceLagIndicator, func(v *bool) { mirror.TrafficMirrorNumSourceLagIndicator = v }, &fieldChanged)
				utils.CompareAndSetBoolField(planItem.TrafficMirrorNumInboundTraffic, stateItem.TrafficMirrorNumInboundTraffic, func(v *bool) { mirror.TrafficMirrorNumInboundTraffic = v }, &fieldChanged)
				utils.CompareAndSetBoolField(planItem.TrafficMirrorNumOutboundTraffic, stateItem.TrafficMirrorNumOutboundTraffic, func(v *bool) { mirror.TrafficMirrorNumOutboundTraffic = v }, &fieldChanged)

				// Handle string field changes
				utils.CompareAndSetStringField(planItem.TrafficMirrorNumSourcePort, stateItem.TrafficMirrorNumSourcePort, func(v *string) { mirror.TrafficMirrorNumSourcePort = v }, &fieldChanged)
				utils.CompareAndSetStringField(planItem.TrafficMirrorNumDestinationPort, stateItem.TrafficMirrorNumDestinationPort, func(v *string) { mirror.TrafficMirrorNumDestinationPort = v }, &fieldChanged)

				// Always include index — API requires it to identify which array element to modify
				utils.SetInt64Fields([]utils.Int64FieldMapping{
					{FieldName: "Index", APIField: &mirror.Index, TFValue: planItem.Index},
				})

				return mirror, fieldChanged
			},
			CreateDeleted: func(index int64) openapi.SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner {
				return openapi.SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner{
					Index: openapi.PtrInt64(int64(index)),
				}
			},
		})
	if trafficMirrorsChanged {
		spProps.TrafficMirrors = changedTrafficMirrors
		hasChanges = true
	}

	// Handle eths
	changedEths, ethsChanged := utils.ProcessIndexedArrayUpdates(plan.Eths, state.Eths,
		utils.IndexedItemHandler[veritySwitchpointEthModel, openapi.SwitchpointsPutRequestSwitchpointValueEthsInner]{
			CreateNew: func(planItem veritySwitchpointEthModel) openapi.SwitchpointsPutRequestSwitchpointValueEthsInner {
				eth := openapi.SwitchpointsPutRequestSwitchpointValueEthsInner{}

				// Handle string fields
				utils.SetStringFields([]utils.StringFieldMapping{
					{FieldName: "Breakout", APIField: &eth.Breakout, TFValue: planItem.Breakout},
					{FieldName: "CustomerVlan", APIField: &eth.CustomerVlan, TFValue: planItem.CustomerVlan},
					{FieldName: "EthNumIcon", APIField: &eth.EthNumIcon, TFValue: planItem.EthNumIcon},
					{FieldName: "EthNumLabel", APIField: &eth.EthNumLabel, TFValue: planItem.EthNumLabel},
					{FieldName: "PortName", APIField: &eth.PortName, TFValue: planItem.PortName},
				})

				// Handle boolean fields
				utils.SetBoolFields([]utils.BoolFieldMapping{
					{FieldName: "Enable", APIField: &eth.Enable, TFValue: planItem.Enable},
				})

				// Handle int64 fields
				utils.SetInt64Fields([]utils.Int64FieldMapping{
					{FieldName: "Index", APIField: &eth.Index, TFValue: planItem.Index},
				})

				return eth
			},
			UpdateExisting: func(planItem veritySwitchpointEthModel, stateItem veritySwitchpointEthModel) (openapi.SwitchpointsPutRequestSwitchpointValueEthsInner, bool) {
				eth := openapi.SwitchpointsPutRequestSwitchpointValueEthsInner{}
				fieldChanged := false

				// Handle string field changes
				utils.CompareAndSetStringField(planItem.Breakout, stateItem.Breakout, func(v *string) { eth.Breakout = v }, &fieldChanged)
				utils.CompareAndSetStringField(planItem.CustomerVlan, stateItem.CustomerVlan, func(v *string) { eth.CustomerVlan = v }, &fieldChanged)
				utils.CompareAndSetStringField(planItem.EthNumIcon, stateItem.EthNumIcon, func(v *string) { eth.EthNumIcon = v }, &fieldChanged)
				utils.CompareAndSetStringField(planItem.EthNumLabel, stateItem.EthNumLabel, func(v *string) { eth.EthNumLabel = v }, &fieldChanged)
				utils.CompareAndSetStringField(planItem.PortName, stateItem.PortName, func(v *string) { eth.PortName = v }, &fieldChanged)

				// Handle boolean field changes
				utils.CompareAndSetBoolField(planItem.Enable, stateItem.Enable, func(v *bool) { eth.Enable = v }, &fieldChanged)

				// Always include index — API requires it to identify which array element to modify
				utils.SetInt64Fields([]utils.Int64FieldMapping{
					{FieldName: "Index", APIField: &eth.Index, TFValue: planItem.Index},
				})

				return eth, fieldChanged
			},
			CreateDeleted: func(index int64) openapi.SwitchpointsPutRequestSwitchpointValueEthsInner {
				return openapi.SwitchpointsPutRequestSwitchpointValueEthsInner{
					Index: openapi.PtrInt64(int64(index)),
				}
			},
		})
	if ethsChanged {
		spProps.Eths = changedEths
		hasChanges = true
	}

	changedPots, potsChanged := utils.ProcessIndexedArrayUpdates(plan.Pots, state.Pots,
		utils.IndexedItemHandler[veritySwitchpointPotsModel, openapi.SwitchpointsPutRequestSwitchpointValuePotsInner]{
			CreateNew: func(planItem veritySwitchpointPotsModel) openapi.SwitchpointsPutRequestSwitchpointValuePotsInner {
				pot := openapi.SwitchpointsPutRequestSwitchpointValuePotsInner{}
				utils.SetBoolFields([]utils.BoolFieldMapping{
					{FieldName: "PotsNumEnable", APIField: &pot.PotsNumEnable, TFValue: planItem.PotsNumEnable},
				})
				utils.SetStringFields([]utils.StringFieldMapping{
					{FieldName: "PotsNumUri", APIField: &pot.PotsNumUri, TFValue: planItem.PotsNumUri},
					{FieldName: "PotsNumUsername", APIField: &pot.PotsNumUsername, TFValue: planItem.PotsNumUsername},
					{FieldName: "PotsNumPassword", APIField: &pot.PotsNumPassword, TFValue: planItem.PotsNumPassword},
					{FieldName: "PotsNumCallerId", APIField: &pot.PotsNumCallerId, TFValue: planItem.PotsNumCallerId},
					{FieldName: "PotsNumHotLine", APIField: &pot.PotsNumHotLine, TFValue: planItem.PotsNumHotLine},
					{FieldName: "PotsNumPasswordEncrypted", APIField: &pot.PotsNumPasswordEncrypted, TFValue: planItem.PotsNumPasswordEncrypted},
				})
				utils.SetInt64Fields([]utils.Int64FieldMapping{
					{FieldName: "Index", APIField: &pot.Index, TFValue: planItem.Index},
				})
				return pot
			},
			UpdateExisting: func(planItem veritySwitchpointPotsModel, stateItem veritySwitchpointPotsModel) (openapi.SwitchpointsPutRequestSwitchpointValuePotsInner, bool) {
				pot := openapi.SwitchpointsPutRequestSwitchpointValuePotsInner{}
				fieldChanged := false
				utils.CompareAndSetBoolField(planItem.PotsNumEnable, stateItem.PotsNumEnable, func(v *bool) { pot.PotsNumEnable = v }, &fieldChanged)
				utils.CompareAndSetStringField(planItem.PotsNumUri, stateItem.PotsNumUri, func(v *string) { pot.PotsNumUri = v }, &fieldChanged)
				utils.CompareAndSetStringField(planItem.PotsNumUsername, stateItem.PotsNumUsername, func(v *string) { pot.PotsNumUsername = v }, &fieldChanged)
				utils.CompareAndSetStringField(planItem.PotsNumPassword, stateItem.PotsNumPassword, func(v *string) { pot.PotsNumPassword = v }, &fieldChanged)
				utils.CompareAndSetStringField(planItem.PotsNumCallerId, stateItem.PotsNumCallerId, func(v *string) { pot.PotsNumCallerId = v }, &fieldChanged)
				utils.CompareAndSetStringField(planItem.PotsNumHotLine, stateItem.PotsNumHotLine, func(v *string) { pot.PotsNumHotLine = v }, &fieldChanged)
				utils.CompareAndSetStringField(planItem.PotsNumPasswordEncrypted, stateItem.PotsNumPasswordEncrypted, func(v *string) { pot.PotsNumPasswordEncrypted = v }, &fieldChanged)
				utils.SetInt64Fields([]utils.Int64FieldMapping{
					{FieldName: "Index", APIField: &pot.Index, TFValue: planItem.Index},
				})
				return pot, fieldChanged
			},
			CreateDeleted: func(index int64) openapi.SwitchpointsPutRequestSwitchpointValuePotsInner {
				return openapi.SwitchpointsPutRequestSwitchpointValuePotsInner{
					Index: openapi.PtrInt64(int64(index)),
				}
			},
		})
	if potsChanged {
		spProps.Pots = changedPots
		hasChanges = true
	}

	// Handle object_properties
	if len(plan.ObjectProperties) > 0 && len(state.ObjectProperties) > 0 {
		objProps := openapi.SwitchpointsPutRequestSwitchpointValueObjectProperties{}
		op := plan.ObjectProperties[0]
		st := state.ObjectProperties[0]
		objPropsChanged := false

		// Get config for nullable field handling in object_properties
		configOp, objPropsCfg := utils.GetObjectPropertiesConfig(op, config.ObjectProperties, configuredAttrs)

		utils.CompareAndSetObjectPropertiesFields([]utils.ObjectPropertiesFieldWithComparison{
			{Name: "UserNotes", PlanValue: op.UserNotes, StateValue: st.UserNotes, APIValue: &objProps.UserNotes},
			{Name: "Aggregate", PlanValue: op.Aggregate, StateValue: st.Aggregate, APIValue: &objProps.Aggregate},
			{Name: "IsHost", PlanValue: op.IsHost, StateValue: st.IsHost, APIValue: &objProps.IsHost},
			{Name: "EmulateRfVideoPort", PlanValue: op.EmulateRfVideoPort, StateValue: st.EmulateRfVideoPort, APIValue: &objProps.EmulateRfVideoPort},
			{Name: "DrawAsEdgeDevice", PlanValue: op.DrawAsEdgeDevice, StateValue: st.DrawAsEdgeDevice, APIValue: &objProps.DrawAsEdgeDevice},
		}, &objPropsChanged)

		if !utils.CompareAndSetRefTypeFields([]utils.RefTypeFieldWithComparison{{
			FieldName:         "expected_parent_endpoint",
			RefTypeFieldName:  "expected_parent_endpoint_ref_type_",
			APIField:          func(v *string) { objProps.ExpectedParentEndpoint = v },
			RefTypeAPIField:   func(v *string) { objProps.ExpectedParentEndpointRefType = v },
			PlanValue:         op.ExpectedParentEndpoint,
			StateValue:        st.ExpectedParentEndpoint,
			PlanRefTypeValue:  op.ExpectedParentEndpointRefType,
			StateRefTypeValue: st.ExpectedParentEndpointRefType,
			SupportMode:       utils.RefTypeSupportOne,
		}}, &objPropsChanged, &resp.Diagnostics) {
			return
		}

		// Handle nullable field in object_properties
		utils.CompareAndSetNullableInt64Field(configOp.NumberOfMultipoints, st.NumberOfMultipoints, objPropsCfg.IsFieldConfigured("number_of_multipoints"), func(v *openapi.NullableInt64) { objProps.NumberOfMultipoints = *v }, &objPropsChanged)

		if objPropsChanged {
			spProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "switchpoint", name, spProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Switchpoint %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "switchpoints")

	var minState veritySwitchpointResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if switchpointData, exists := bulkMgr.GetResourceResponse("switchpoint", name); exists {
			effectiveData := utils.MergeMissingPlanScalars(switchpointData, plan, switchpointResourceType, r.provCtx.mode)
			state := populateSwitchpointState(ctx, minState, effectiveData, r.provCtx.mode)
			preserveSwitchpointPortNames(&state, &plan)
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

	postOpCtx := utils.WithPostOperationFallback(ctx, plan, switchpointResourceType, r.provCtx.mode)
	r.Read(postOpCtx, readReq, &readResp)
	if readResp.State.Raw.IsNull() {
		_, diags := utils.SetPostOperationFallbackState(postOpCtx, &readResp.State)
		readResp.Diagnostics.Append(diags...)
	}
	resp.State = readResp.State
	resp.Diagnostics = readResp.Diagnostics

	if resp.Diagnostics.HasError() {
		return
	}
	var readState veritySwitchpointResourceModel
	resp.Diagnostics.Append(readResp.State.Get(ctx, &readState)...)
	if resp.Diagnostics.HasError() {
		return
	}
	preserveSwitchpointPortNames(&readState, &plan)
	resp.Diagnostics.Append(resp.State.Set(ctx, &readState)...)
}

func (r *veritySwitchpointResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state veritySwitchpointResourceModel
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "switchpoint", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Switchpoint %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "switchpoints")
	resp.State.RemoveResource(ctx)
}

func (r *veritySwitchpointResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func populateSwitchpointState(ctx context.Context, state veritySwitchpointResourceModel, switchpointData map[string]interface{}, mode string) veritySwitchpointResourceModel {
	const resourceType = switchpointResourceType

	state.Name = utils.MapStringFromAPI(switchpointData["name"])

	// Int fields
	state.BgpAsNumber = utils.MapInt64WithMode(switchpointData, "bgp_as_number", resourceType, mode)
	state.ExpectedUplinkPort = utils.MapInt64WithMode(switchpointData, "expected_uplink_port", resourceType, mode)

	// Number fields
	state.Position = utils.MapNumberWithMode(switchpointData, "position", resourceType, mode)
	state.RailGroup = utils.MapNumberWithMode(switchpointData, "rail_group", resourceType, mode)

	// Bool fields
	state.Enable = utils.MapBoolWithMode(switchpointData, "enable", resourceType, mode)
	state.IsTopOfIsland = utils.MapBoolWithMode(switchpointData, "is_top_of_island", resourceType, mode)
	state.ReadOnlyMode = utils.MapBoolWithMode(switchpointData, "read_only_mode", resourceType, mode)
	state.Locked = utils.MapBoolWithMode(switchpointData, "locked", resourceType, mode)
	state.OutOfBandManagement = utils.MapBoolWithMode(switchpointData, "out_of_band_management", resourceType, mode)
	state.BbSwitch = utils.MapBoolWithMode(switchpointData, "bb_switch", resourceType, mode)
	state.ManagedOnNativeVlan = utils.MapBoolWithMode(switchpointData, "managed_on_native_vlan", resourceType, mode)
	state.IsFabric = utils.MapBoolWithMode(switchpointData, "is_fabric", resourceType, mode)
	state.UsesTaggedPackets = utils.MapBoolWithMode(switchpointData, "uses_tagged_packets", resourceType, mode)
	state.BgpAsNumberAutoAssigned = utils.MapBoolWithMode(switchpointData, "bgp_as_number_auto_assigned_", resourceType, mode)
	state.SwitchVtepIdIpMaskAutoAssigned = utils.MapBoolWithMode(switchpointData, "switch_vtep_id_ip_mask_auto_assigned_", resourceType, mode)
	state.SwitchRouterIdIpMaskAutoAssigned = utils.MapBoolWithMode(switchpointData, "switch_router_id_ip_mask_auto_assigned_", resourceType, mode)
	state.ControllerIpAndMaskAutoAssigned = utils.MapBoolWithMode(switchpointData, "controller_ip_and_mask_auto_assigned_", resourceType, mode)
	state.SwitchIpAndMaskAutoAssigned = utils.MapBoolWithMode(switchpointData, "switch_ip_and_mask_auto_assigned_", resourceType, mode)
	// String fields
	state.Tenant = utils.MapStringWithMode(switchpointData, "tenant", resourceType, mode)
	state.TenantRefType = utils.MapStringWithMode(switchpointData, "tenant_ref_type_", resourceType, mode)
	state.DeviceSerialNumber = utils.MapStringWithMode(switchpointData, "device_serial_number", resourceType, mode)
	state.ConnectedBundle = utils.MapStringWithMode(switchpointData, "connected_bundle", resourceType, mode)
	state.ConnectedBundleRefType = utils.MapStringWithMode(switchpointData, "connected_bundle_ref_type_", resourceType, mode)
	state.ExpectedSite = utils.MapStringWithMode(switchpointData, "expected_site", resourceType, mode)
	state.ExpectedSiteRefType = utils.MapStringWithMode(switchpointData, "expected_site_ref_type_", resourceType, mode)
	state.Type = utils.MapStringWithMode(switchpointData, "type", resourceType, mode)
	state.Plane = utils.MapStringWithMode(switchpointData, "plane", resourceType, mode)
	state.PlaneRefType = utils.MapStringWithMode(switchpointData, "plane_ref_type_", resourceType, mode)
	state.SpinePlane = utils.MapStringWithMode(switchpointData, "spine_plane", resourceType, mode)
	state.SpinePlaneRefType = utils.MapStringWithMode(switchpointData, "spine_plane_ref_type_", resourceType, mode)
	state.Pod = utils.MapStringWithMode(switchpointData, "pod", resourceType, mode)
	state.PodRefType = utils.MapStringWithMode(switchpointData, "pod_ref_type_", resourceType, mode)
	state.Su = utils.MapStringWithMode(switchpointData, "su", resourceType, mode)
	state.SuRefType = utils.MapStringWithMode(switchpointData, "su_ref_type_", resourceType, mode)
	state.SspGroup = utils.MapStringWithMode(switchpointData, "ssp_group", resourceType, mode)
	state.SspGroupRefType = utils.MapStringWithMode(switchpointData, "ssp_group_ref_type_", resourceType, mode)
	state.RackInfo = utils.MapStringWithMode(switchpointData, "rack_info", resourceType, mode)
	state.Rack = utils.MapStringWithMode(switchpointData, "rack", resourceType, mode)
	state.RackRefType = utils.MapStringWithMode(switchpointData, "rack_ref_type_", resourceType, mode)
	state.SwitchRouterIdIpMask = utils.MapStringWithMode(switchpointData, "switch_router_id_ip_mask", resourceType, mode)
	state.SwitchVtepIdIpMask = utils.MapStringWithMode(switchpointData, "switch_vtep_id_ip_mask", resourceType, mode)
	state.PasswordEncrypted = utils.MapStringWithMode(switchpointData, "password_encrypted", resourceType, mode)
	state.EnablePasswordEncrypted = utils.MapStringWithMode(switchpointData, "enable_password_encrypted", resourceType, mode)
	state.SshKeyOrPasswordEncrypted = utils.MapStringWithMode(switchpointData, "ssh_key_or_password_encrypted", resourceType, mode)
	state.PassphraseEncrypted = utils.MapStringWithMode(switchpointData, "passphrase_encrypted", resourceType, mode)
	state.PrivatePasswordEncrypted = utils.MapStringWithMode(switchpointData, "private_password_encrypted", resourceType, mode)
	state.IpSource = utils.MapStringWithMode(switchpointData, "ip_source", resourceType, mode)
	state.ControllerIpAndMask = utils.MapStringWithMode(switchpointData, "controller_ip_and_mask", resourceType, mode)
	state.Gateway = utils.MapStringWithMode(switchpointData, "gateway", resourceType, mode)
	state.SwitchIpAndMask = utils.MapStringWithMode(switchpointData, "switch_ip_and_mask", resourceType, mode)
	state.SwitchGateway = utils.MapStringWithMode(switchpointData, "switch_gateway", resourceType, mode)
	state.CommType = utils.MapStringWithMode(switchpointData, "comm_type", resourceType, mode)
	state.SnmpCommunityString = utils.MapStringWithMode(switchpointData, "snmp_community_string", resourceType, mode)
	state.UplinkPort = utils.MapStringWithMode(switchpointData, "uplink_port", resourceType, mode)
	state.ExpectedBreakout = utils.MapStringWithMode(switchpointData, "expected_breakout", resourceType, mode)
	state.LldpSearchString = utils.MapStringWithMode(switchpointData, "lldp_search_string", resourceType, mode)
	state.ZtpIdentification = utils.MapStringWithMode(switchpointData, "ztp_identification", resourceType, mode)
	state.LocatedBy = utils.MapStringWithMode(switchpointData, "located_by", resourceType, mode)
	state.PowerState = utils.MapStringWithMode(switchpointData, "power_state", resourceType, mode)
	state.CommunicationMode = utils.MapStringWithMode(switchpointData, "communication_mode", resourceType, mode)
	state.CliAccessMode = utils.MapStringWithMode(switchpointData, "cli_access_mode", resourceType, mode)
	state.Username = utils.MapStringWithMode(switchpointData, "username", resourceType, mode)
	state.Password = utils.MapStringWithMode(switchpointData, "password", resourceType, mode)
	state.EnablePassword = utils.MapStringWithMode(switchpointData, "enable_password", resourceType, mode)
	state.SshKeyOrPassword = utils.MapStringWithMode(switchpointData, "ssh_key_or_password", resourceType, mode)
	state.Sdlc = utils.MapStringWithMode(switchpointData, "sdlc", resourceType, mode)
	state.SecurityType = utils.MapStringWithMode(switchpointData, "security_type", resourceType, mode)
	state.Snmpv3Username = utils.MapStringWithMode(switchpointData, "snmpv3_username", resourceType, mode)
	state.AuthenticationProtocol = utils.MapStringWithMode(switchpointData, "authentication_protocol", resourceType, mode)
	state.Passphrase = utils.MapStringWithMode(switchpointData, "passphrase", resourceType, mode)
	state.PrivateProtocol = utils.MapStringWithMode(switchpointData, "private_protocol", resourceType, mode)
	state.PrivatePassword = utils.MapStringWithMode(switchpointData, "private_password", resourceType, mode)
	state.DeviceManagedAs = utils.MapStringWithMode(switchpointData, "device_managed_as", resourceType, mode)
	state.Switch = utils.MapStringWithMode(switchpointData, "switch", resourceType, mode)
	state.SwitchRefType = utils.MapStringWithMode(switchpointData, "switch_ref_type_", resourceType, mode)
	state.ConnectionService = utils.MapStringWithMode(switchpointData, "connection_service", resourceType, mode)
	state.ConnectionServiceRefType = utils.MapStringWithMode(switchpointData, "connection_service_ref_type_", resourceType, mode)
	state.Port = utils.MapStringWithMode(switchpointData, "port", resourceType, mode)

	// Handle object_properties block
	if utils.FieldAppliesToMode(resourceType, "object_properties", mode) {
		if objProps, ok := switchpointData["object_properties"].(map[string]interface{}); ok {
			op := veritySwitchpointObjectPropertiesModel{
				UserNotes:                     utils.MapStringWithModeNested(objProps, "user_notes", resourceType, "object_properties.user_notes", mode),
				ExpectedParentEndpoint:        utils.MapStringWithModeNested(objProps, "expected_parent_endpoint", resourceType, "object_properties.expected_parent_endpoint", mode),
				ExpectedParentEndpointRefType: utils.MapStringWithModeNested(objProps, "expected_parent_endpoint_ref_type_", resourceType, "object_properties.expected_parent_endpoint_ref_type_", mode),
				NumberOfMultipoints:           utils.MapInt64WithModeNested(objProps, "number_of_multipoints", resourceType, "object_properties.number_of_multipoints", mode),
				Aggregate:                     utils.MapBoolWithModeNested(objProps, "aggregate", resourceType, "object_properties.aggregate", mode),
				IsHost:                        utils.MapBoolWithModeNested(objProps, "is_host", resourceType, "object_properties.is_host", mode),
				EmulateRfVideoPort:            utils.MapBoolWithModeNested(objProps, "emulate_rf_video_port", resourceType, "object_properties.emulate_rf_video_port", mode),
				DrawAsEdgeDevice:              utils.MapBoolWithModeNested(objProps, "draw_as_edge_device", resourceType, "object_properties.draw_as_edge_device", mode),
			}
			state.ObjectProperties = []veritySwitchpointObjectPropertiesModel{op}
		} else {
			state.ObjectProperties = nil
		}
	} else {
		state.ObjectProperties = nil
	}

	// Handle badges block
	if utils.FieldAppliesToMode(resourceType, "badges", mode) {
		if badgesArray, ok := switchpointData["badges"].([]interface{}); ok && len(badgesArray) > 0 {
			var badges []veritySwitchpointBadgeModel
			for _, b := range badgesArray {
				badge, ok := b.(map[string]interface{})
				if !ok {
					continue
				}
				badgeModel := veritySwitchpointBadgeModel{
					Badge:        utils.MapStringWithModeNested(badge, "badge", resourceType, "badges.badge", mode),
					BadgeRefType: utils.MapStringWithModeNested(badge, "badge_ref_type_", resourceType, "badges.badge_ref_type_", mode),
					Index:        utils.MapInt64WithModeNested(badge, "index", resourceType, "badges.index", mode),
				}
				badges = append(badges, badgeModel)
			}
			state.Badges = badges
		} else {
			state.Badges = nil
		}
	} else {
		state.Badges = nil
	}

	// Handle children block
	if utils.FieldAppliesToMode(resourceType, "children", mode) {
		if childrenArray, ok := switchpointData["children"].([]interface{}); ok && len(childrenArray) > 0 {
			var children []veritySwitchpointChildModel
			for _, c := range childrenArray {
				child, ok := c.(map[string]interface{})
				if !ok {
					continue
				}
				childModel := veritySwitchpointChildModel{
					ChildNumEndpoint:        utils.MapStringWithModeNested(child, "child_num_endpoint", resourceType, "children.child_num_endpoint", mode),
					ChildNumEndpointRefType: utils.MapStringWithModeNested(child, "child_num_endpoint_ref_type_", resourceType, "children.child_num_endpoint_ref_type_", mode),
					ChildNumDevice:          utils.MapStringWithModeNested(child, "child_num_device", resourceType, "children.child_num_device", mode),
					Index:                   utils.MapInt64WithModeNested(child, "index", resourceType, "children.index", mode),
				}
				children = append(children, childModel)
			}
			state.Children = children
		} else {
			state.Children = nil
		}
	} else {
		state.Children = nil
	}

	// Handle traffic_mirrors block
	if utils.FieldAppliesToMode(resourceType, "traffic_mirrors", mode) {
		if mirrorsArray, ok := switchpointData["traffic_mirrors"].([]interface{}); ok && len(mirrorsArray) > 0 {
			var mirrors []veritySwitchpointTrafficMirrorModel
			for _, m := range mirrorsArray {
				mirror, ok := m.(map[string]interface{})
				if !ok {
					continue
				}
				mirrorModel := veritySwitchpointTrafficMirrorModel{
					TrafficMirrorNumEnable:             utils.MapBoolWithModeNested(mirror, "traffic_mirror_num_enable", resourceType, "traffic_mirrors.traffic_mirror_num_enable", mode),
					TrafficMirrorNumSourcePort:         utils.MapStringWithModeNested(mirror, "traffic_mirror_num_source_port", resourceType, "traffic_mirrors.traffic_mirror_num_source_port", mode),
					TrafficMirrorNumSourceLagIndicator: utils.MapBoolWithModeNested(mirror, "traffic_mirror_num_source_lag_indicator", resourceType, "traffic_mirrors.traffic_mirror_num_source_lag_indicator", mode),
					TrafficMirrorNumDestinationPort:    utils.MapStringWithModeNested(mirror, "traffic_mirror_num_destination_port", resourceType, "traffic_mirrors.traffic_mirror_num_destination_port", mode),
					TrafficMirrorNumInboundTraffic:     utils.MapBoolWithModeNested(mirror, "traffic_mirror_num_inbound_traffic", resourceType, "traffic_mirrors.traffic_mirror_num_inbound_traffic", mode),
					TrafficMirrorNumOutboundTraffic:    utils.MapBoolWithModeNested(mirror, "traffic_mirror_num_outbound_traffic", resourceType, "traffic_mirrors.traffic_mirror_num_outbound_traffic", mode),
					Index:                              utils.MapInt64WithModeNested(mirror, "index", resourceType, "traffic_mirrors.index", mode),
				}
				mirrors = append(mirrors, mirrorModel)
			}
			state.TrafficMirrors = mirrors
		} else {
			state.TrafficMirrors = nil
		}
	} else {
		state.TrafficMirrors = nil
	}

	// Handle eths block
	if utils.FieldAppliesToMode(resourceType, "eths", mode) {
		if ethsArray, ok := switchpointData["eths"].([]interface{}); ok && len(ethsArray) > 0 {
			var eths []veritySwitchpointEthModel
			for _, e := range ethsArray {
				eth, ok := e.(map[string]interface{})
				if !ok {
					continue
				}
				ethModel := veritySwitchpointEthModel{
					Breakout:     utils.MapStringWithModeNested(eth, "breakout", resourceType, "eths.breakout", mode),
					CustomerVlan: utils.MapStringWithModeNested(eth, "customer_vlan", resourceType, "eths.customer_vlan", mode),
					Index:        utils.MapInt64WithModeNested(eth, "index", resourceType, "eths.index", mode),
					EthNumIcon:   utils.MapStringWithModeNested(eth, "eth_num_icon", resourceType, "eths.eth_num_icon", mode),
					EthNumLabel:  utils.MapStringWithModeNested(eth, "eth_num_label", resourceType, "eths.eth_num_label", mode),
					Enable:       utils.MapBoolWithModeNested(eth, "enable", resourceType, "eths.enable", mode),
					PortName:     utils.MapStringWithModeNested(eth, "port_name", resourceType, "eths.port_name", mode),
				}
				eths = append(eths, ethModel)
			}
			state.Eths = eths
		} else {
			state.Eths = nil
		}
	} else {
		state.Eths = nil
	}

	if utils.FieldAppliesToMode(resourceType, "pots", mode) {
		if potsArray, ok := switchpointData["pots"].([]interface{}); ok && len(potsArray) > 0 {
			var pots []veritySwitchpointPotsModel
			for _, p := range potsArray {
				pot, ok := p.(map[string]interface{})
				if !ok {
					continue
				}
				potModel := veritySwitchpointPotsModel{
					PotsNumEnable:            utils.MapBoolWithModeNested(pot, "pots_num_enable", resourceType, "pots.pots_num_enable", mode),
					PotsNumUri:               utils.MapStringWithModeNested(pot, "pots_num_uri", resourceType, "pots.pots_num_uri", mode),
					PotsNumUsername:          utils.MapStringWithModeNested(pot, "pots_num_username", resourceType, "pots.pots_num_username", mode),
					PotsNumPassword:          utils.MapStringWithModeNested(pot, "pots_num_password", resourceType, "pots.pots_num_password", mode),
					PotsNumCallerId:          utils.MapStringWithModeNested(pot, "pots_num_caller_id", resourceType, "pots.pots_num_caller_id", mode),
					PotsNumHotLine:           utils.MapStringWithModeNested(pot, "pots_num_hot_line", resourceType, "pots.pots_num_hot_line", mode),
					PotsNumPasswordEncrypted: utils.MapStringWithModeNested(pot, "pots_num_password_encrypted", resourceType, "pots.pots_num_password_encrypted", mode),
					Index:                    utils.MapInt64WithModeNested(pot, "index", resourceType, "pots.index", mode),
				}
				pots = append(pots, potModel)
			}
			state.Pots = pots
		} else {
			state.Pots = nil
		}
	} else {
		state.Pots = nil
	}

	return state
}

func (r *veritySwitchpointResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// =========================================================================
	// Skip if deleting
	// =========================================================================
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan veritySwitchpointResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Mode-aware field nullification
	// Set fields that don't apply to current mode to null to prevent
	// "known after apply" messages for irrelevant fields.
	// =========================================================================
	const resourceType = switchpointResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyStrings(
		"tenant", "tenant_ref_type_", "device_serial_number",
		"connected_bundle", "connected_bundle_ref_type_",
		"expected_site", "expected_site_ref_type_",
		"type", "plane", "plane_ref_type_", "spine_plane", "spine_plane_ref_type_",
		"pod", "pod_ref_type_", "su", "su_ref_type_",
		"ssp_group", "ssp_group_ref_type_", "rack_info", "rack", "rack_ref_type_",
		"switch_router_id_ip_mask", "switch_vtep_id_ip_mask",
		"password_encrypted", "enable_password_encrypted",
		"ssh_key_or_password_encrypted", "passphrase_encrypted",
		"private_password_encrypted", "ip_source",
		"controller_ip_and_mask", "gateway", "switch_ip_and_mask",
		"switch_gateway", "comm_type", "snmp_community_string",
		"uplink_port", "expected_breakout", "lldp_search_string", "ztp_identification",
		"located_by", "power_state", "communication_mode",
		"cli_access_mode", "username", "password", "enable_password",
		"ssh_key_or_password", "sdlc", "security_type",
		"snmpv3_username", "authentication_protocol", "passphrase",
		"private_protocol", "private_password", "device_managed_as",
		"switch", "switch_ref_type_", "connection_service",
		"connection_service_ref_type_", "port",
	)

	nullifier.NullifyBools(
		"enable", "is_top_of_island", "read_only_mode", "locked",
		"out_of_band_management", "bb_switch", "managed_on_native_vlan",
		"is_fabric", "uses_tagged_packets",
		"switch_router_id_ip_mask_auto_assigned_",
		"switch_vtep_id_ip_mask_auto_assigned_",
		"bgp_as_number_auto_assigned_",
		"controller_ip_and_mask_auto_assigned_",
		"switch_ip_and_mask_auto_assigned_",
	)

	nullifier.NullifyInt64s(
		"bgp_as_number", "expected_uplink_port",
	)
	nullifier.NullifyNumbers(
		"position", "rail_group",
	)

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "badges",
		ItemCount:    len(plan.Badges),
		StringFields: []string{"badge", "badge_ref_type_"},
		Int64Fields:  []string{"index"},
	})

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "children",
		ItemCount:    len(plan.Children),
		StringFields: []string{"child_num_endpoint", "child_num_endpoint_ref_type_", "child_num_device"},
		Int64Fields:  []string{"index"},
	})

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "traffic_mirrors",
		ItemCount:    len(plan.TrafficMirrors),
		StringFields: []string{"traffic_mirror_num_source_port", "traffic_mirror_num_destination_port"},
		BoolFields:   []string{"traffic_mirror_num_enable", "traffic_mirror_num_source_lag_indicator", "traffic_mirror_num_inbound_traffic", "traffic_mirror_num_outbound_traffic"},
		Int64Fields:  []string{"index"},
	})

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "eths",
		ItemCount:    len(plan.Eths),
		StringFields: []string{"breakout", "customer_vlan", "eth_num_icon", "eth_num_label", "port_name"},
		BoolFields:   []string{"enable"},
		Int64Fields:  []string{"index"},
	})

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "pots",
		ItemCount:    len(plan.Pots),
		StringFields: []string{"pots_num_uri", "pots_num_username", "pots_num_password", "pots_num_caller_id", "pots_num_hot_line", "pots_num_password_encrypted"},
		BoolFields:   []string{"pots_num_enable"},
		Int64Fields:  []string{"index"},
	})

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "object_properties",
		ItemCount:    len(plan.ObjectProperties),
		StringFields: []string{"user_notes", "expected_parent_endpoint", "expected_parent_endpoint_ref_type_"},
		BoolFields:   []string{"aggregate", "is_host", "emulate_rf_video_port", "draw_as_edge_device"},
		Int64Fields:  []string{"number_of_multipoints"},
	})

	// =========================================================================
	// CREATE operation - handle auto-assigned fields
	// =========================================================================
	if req.State.Raw.IsNull() {
		// Switchpoint-specific: auto-assignment on create
		if !plan.BgpAsNumberAutoAssigned.IsNull() && plan.BgpAsNumberAutoAssigned.ValueBool() {
			resp.Plan.SetAttribute(ctx, path.Root("bgp_as_number"), types.Int64Unknown())
		}
		if !plan.SwitchRouterIdIpMaskAutoAssigned.IsNull() && plan.SwitchRouterIdIpMaskAutoAssigned.ValueBool() {
			resp.Plan.SetAttribute(ctx, path.Root("switch_router_id_ip_mask"), types.StringUnknown())
		}
		if !plan.SwitchVtepIdIpMaskAutoAssigned.IsNull() && plan.SwitchVtepIdIpMaskAutoAssigned.ValueBool() {
			resp.Plan.SetAttribute(ctx, path.Root("switch_vtep_id_ip_mask"), types.StringUnknown())
		}
		if !plan.ControllerIpAndMaskAutoAssigned.IsNull() && plan.ControllerIpAndMaskAutoAssigned.ValueBool() {
			resp.Plan.SetAttribute(ctx, path.Root("controller_ip_and_mask"), types.StringUnknown())
		}
		if !plan.SwitchIpAndMaskAutoAssigned.IsNull() && plan.SwitchIpAndMaskAutoAssigned.ValueBool() {
			resp.Plan.SetAttribute(ctx, path.Root("switch_ip_and_mask"), types.StringUnknown())
		}
		return
	}

	// =========================================================================
	// UPDATE operation - get state and config
	// =========================================================================
	var state veritySwitchpointResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config veritySwitchpointResourceModel
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
	workDir := r.provCtx.workDir
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, switchpointTerraformType, name)

	utils.HandleNullableFields(utils.NullableFieldsConfig{
		Ctx:             ctx,
		Plan:            &resp.Plan,
		ConfiguredAttrs: configuredAttrs,
		Int64Fields: []utils.NullableInt64Field{
			{AttrName: "bgp_as_number", ConfigVal: config.BgpAsNumber, StateVal: state.BgpAsNumber},
		},
		NumberFields: []utils.NullableNumberField{
			{AttrName: "position", ConfigVal: config.Position, StateVal: state.Position},
			{AttrName: "rail_group", ConfigVal: config.RailGroup, StateVal: state.RailGroup},
		},
	})

	// =========================================================================
	// Handle nullable fields in nested blocks
	// =========================================================================
	if len(config.ObjectProperties) > 0 && len(state.ObjectProperties) > 0 {
		configOp := config.ObjectProperties[0]
		stateOp := state.ObjectProperties[0]

		if configuredAttrs.IsBlockAttributeConfigured("object_properties.number_of_multipoints") &&
			configOp.NumberOfMultipoints.IsNull() && !stateOp.NumberOfMultipoints.IsNull() {
			resp.Plan.SetAttribute(ctx, path.Root("object_properties").AtListIndex(0).AtName("number_of_multipoints"), types.Int64Null())
		}
	}

	// =========================================================================
	// Validate auto-assigned field specifications
	// =========================================================================
	if !config.BgpAsNumberAutoAssigned.IsNull() && config.BgpAsNumberAutoAssigned.ValueBool() {
		if !config.BgpAsNumber.IsNull() && !config.BgpAsNumber.IsUnknown() {
			resp.Diagnostics.AddError(
				"BGP AS Number cannot be specified when auto-assigned",
				"The 'bgp_as_number' field cannot be specified in the configuration when 'bgp_as_number_auto_assigned_' is set to true. The API will assign this value automatically.",
			)
			return
		}
	}

	if !config.SwitchRouterIdIpMaskAutoAssigned.IsNull() && config.SwitchRouterIdIpMaskAutoAssigned.ValueBool() {
		if !config.SwitchRouterIdIpMask.IsNull() && !config.SwitchRouterIdIpMask.IsUnknown() && config.SwitchRouterIdIpMask.ValueString() != "" {
			resp.Diagnostics.AddError(
				"Switch Router ID IP Mask cannot be specified when auto-assigned",
				"The 'switch_router_id_ip_mask' field cannot be specified in the configuration when 'switch_router_id_ip_mask_auto_assigned_' is set to true. The API will assign this value automatically.",
			)
			return
		}
	}

	if !config.SwitchVtepIdIpMaskAutoAssigned.IsNull() && config.SwitchVtepIdIpMaskAutoAssigned.ValueBool() {
		if !config.SwitchVtepIdIpMask.IsNull() && !config.SwitchVtepIdIpMask.IsUnknown() && config.SwitchVtepIdIpMask.ValueString() != "" {
			resp.Diagnostics.AddError(
				"Switch VTEP ID IP Mask cannot be specified when auto-assigned",
				"The 'switch_vtep_id_ip_mask' field cannot be specified in the configuration when 'switch_vtep_id_ip_mask_auto_assigned_' is set to true. The API will assign this value automatically.",
			)
			return
		}
	}

	if !config.ControllerIpAndMaskAutoAssigned.IsNull() && config.ControllerIpAndMaskAutoAssigned.ValueBool() {
		if !config.ControllerIpAndMask.IsNull() && !config.ControllerIpAndMask.IsUnknown() && config.ControllerIpAndMask.ValueString() != "" {
			resp.Diagnostics.AddError(
				"Controller IP and Mask cannot be specified when auto-assigned",
				"The 'controller_ip_and_mask' field cannot be specified in the configuration when 'controller_ip_and_mask_auto_assigned_' is set to true. The API will assign this value automatically.",
			)
			return
		}
	}

	if !config.SwitchIpAndMaskAutoAssigned.IsNull() && config.SwitchIpAndMaskAutoAssigned.ValueBool() {
		if !config.SwitchIpAndMask.IsNull() && !config.SwitchIpAndMask.IsUnknown() && config.SwitchIpAndMask.ValueString() != "" {
			resp.Diagnostics.AddError(
				"Switch IP and Mask cannot be specified when auto-assigned",
				"The 'switch_ip_and_mask' field cannot be specified in the configuration when 'switch_ip_and_mask_auto_assigned_' is set to true. The API will assign this value automatically.",
			)
			return
		}
	}

	// =========================================================================
	// Resource-specific auto-assigned field logic (BgpAsNumber)
	// =========================================================================
	if !plan.BgpAsNumberAutoAssigned.IsNull() && plan.BgpAsNumberAutoAssigned.ValueBool() {
		if !plan.BgpAsNumberAutoAssigned.Equal(state.BgpAsNumberAutoAssigned) {
			// bgp_as_number_auto_assigned_ is changing to true - API will assign value
			resp.Plan.SetAttribute(ctx, path.Root("bgp_as_number"), types.Int64Unknown())
			resp.Diagnostics.AddWarning(
				"BGP AS Number will be assigned by the API",
				"The 'bgp_as_number' field will be automatically assigned by the API because 'bgp_as_number_auto_assigned_' is being set to true.",
			)
		} else if !plan.BgpAsNumber.Equal(state.BgpAsNumber) {
			// User tried to change BgpAsNumber but it's auto-assigned - suppress diff
			resp.Diagnostics.AddWarning(
				"Ignoring bgp_as_number changes with auto-assignment enabled",
				"The 'bgp_as_number' field changes will be ignored because 'bgp_as_number_auto_assigned_' is set to true.",
			)
			if !state.BgpAsNumber.IsNull() {
				resp.Plan.SetAttribute(ctx, path.Root("bgp_as_number"), state.BgpAsNumber)
			}
		}
	}

	// =========================================================================
	// Resource-specific auto-assigned field logic (SwitchRouterIdIpMask)
	// =========================================================================
	if !plan.SwitchRouterIdIpMaskAutoAssigned.IsNull() && plan.SwitchRouterIdIpMaskAutoAssigned.ValueBool() {
		if !plan.SwitchRouterIdIpMaskAutoAssigned.Equal(state.SwitchRouterIdIpMaskAutoAssigned) {
			// switch_router_id_ip_mask_auto_assigned_ is changing to true - API will assign value
			resp.Plan.SetAttribute(ctx, path.Root("switch_router_id_ip_mask"), types.StringUnknown())
			resp.Diagnostics.AddWarning(
				"Switch Router ID IP Mask will be assigned by the API",
				"The 'switch_router_id_ip_mask' field will be automatically assigned by the API because 'switch_router_id_ip_mask_auto_assigned_' is being set to true.",
			)
		} else if !plan.SwitchRouterIdIpMask.Equal(state.SwitchRouterIdIpMask) {
			// User tried to change SwitchRouterIdIpMask but it's auto-assigned - suppress diff
			resp.Diagnostics.AddWarning(
				"Ignoring switch_router_id_ip_mask changes with auto-assignment enabled",
				"The 'switch_router_id_ip_mask' field changes will be ignored because 'switch_router_id_ip_mask_auto_assigned_' is set to true.",
			)
			if !state.SwitchRouterIdIpMask.IsNull() {
				resp.Plan.SetAttribute(ctx, path.Root("switch_router_id_ip_mask"), state.SwitchRouterIdIpMask)
			}
		}
	}

	// =========================================================================
	// Resource-specific auto-assigned field logic (SwitchVtepIdIpMask)
	// =========================================================================
	if !plan.SwitchVtepIdIpMaskAutoAssigned.IsNull() && plan.SwitchVtepIdIpMaskAutoAssigned.ValueBool() {
		if !plan.SwitchVtepIdIpMaskAutoAssigned.Equal(state.SwitchVtepIdIpMaskAutoAssigned) {
			// switch_vtep_id_ip_mask_auto_assigned_ is changing to true - API will assign value
			resp.Plan.SetAttribute(ctx, path.Root("switch_vtep_id_ip_mask"), types.StringUnknown())
			resp.Diagnostics.AddWarning(
				"Switch VTEP ID IP Mask will be assigned by the API",
				"The 'switch_vtep_id_ip_mask' field will be automatically assigned by the API because 'switch_vtep_id_ip_mask_auto_assigned_' is being set to true.",
			)
		} else if !plan.SwitchVtepIdIpMask.Equal(state.SwitchVtepIdIpMask) {
			// User tried to change SwitchVtepIdIpMask but it's auto-assigned - suppress diff
			resp.Diagnostics.AddWarning(
				"Ignoring switch_vtep_id_ip_mask changes with auto-assignment enabled",
				"The 'switch_vtep_id_ip_mask' field changes will be ignored because 'switch_vtep_id_ip_mask_auto_assigned_' is set to true.",
			)
			if !state.SwitchVtepIdIpMask.IsNull() {
				resp.Plan.SetAttribute(ctx, path.Root("switch_vtep_id_ip_mask"), state.SwitchVtepIdIpMask)
			}
		}
	}

	// =========================================================================
	// Resource-specific auto-assigned field logic (ControllerIpAndMask)
	// =========================================================================
	if !plan.ControllerIpAndMaskAutoAssigned.IsNull() && plan.ControllerIpAndMaskAutoAssigned.ValueBool() {
		if !plan.ControllerIpAndMaskAutoAssigned.Equal(state.ControllerIpAndMaskAutoAssigned) {
			resp.Plan.SetAttribute(ctx, path.Root("controller_ip_and_mask"), types.StringUnknown())
			resp.Diagnostics.AddWarning(
				"Controller IP and Mask will be assigned by the API",
				"The 'controller_ip_and_mask' field will be automatically assigned by the API because 'controller_ip_and_mask_auto_assigned_' is being set to true.",
			)
		} else if !plan.ControllerIpAndMask.Equal(state.ControllerIpAndMask) {
			resp.Diagnostics.AddWarning(
				"Ignoring controller_ip_and_mask changes with auto-assignment enabled",
				"The 'controller_ip_and_mask' field changes will be ignored because 'controller_ip_and_mask_auto_assigned_' is set to true.",
			)
			if !state.ControllerIpAndMask.IsNull() {
				resp.Plan.SetAttribute(ctx, path.Root("controller_ip_and_mask"), state.ControllerIpAndMask)
			}
		}
	}

	// =========================================================================
	// Resource-specific auto-assigned field logic (SwitchIpAndMask)
	// =========================================================================
	if !plan.SwitchIpAndMaskAutoAssigned.IsNull() && plan.SwitchIpAndMaskAutoAssigned.ValueBool() {
		if !plan.SwitchIpAndMaskAutoAssigned.Equal(state.SwitchIpAndMaskAutoAssigned) {
			resp.Plan.SetAttribute(ctx, path.Root("switch_ip_and_mask"), types.StringUnknown())
			resp.Diagnostics.AddWarning(
				"Switch IP and Mask will be assigned by the API",
				"The 'switch_ip_and_mask' field will be automatically assigned by the API because 'switch_ip_and_mask_auto_assigned_' is being set to true.",
			)
		} else if !plan.SwitchIpAndMask.Equal(state.SwitchIpAndMask) {
			resp.Diagnostics.AddWarning(
				"Ignoring switch_ip_and_mask changes with auto-assignment enabled",
				"The 'switch_ip_and_mask' field changes will be ignored because 'switch_ip_and_mask_auto_assigned_' is set to true.",
			)
			if !state.SwitchIpAndMask.IsNull() {
				resp.Plan.SetAttribute(ctx, path.Root("switch_ip_and_mask"), state.SwitchIpAndMask)
			}
		}
	}
}

// preserveSwitchpointPortNames copies port_name values from a reference source (plan or prior state)
// into the populated state. The API documents port_name as "reference only" – it accepts the value on PUT
// but never persists or returns it, so GET always gives back "".
func preserveSwitchpointPortNames(state *veritySwitchpointResourceModel, ref *veritySwitchpointResourceModel) {
	if ref == nil || len(ref.Eths) == 0 || len(state.Eths) == 0 {
		return
	}

	for i := range state.Eths {
		if i >= len(ref.Eths) {
			break
		}

		refPortName := ref.Eths[i].PortName

		if !refPortName.IsNull() && !refPortName.IsUnknown() && state.Eths[i].PortName.ValueString() == "" {
			state.Eths[i].PortName = refPortName
		}
	}
}
