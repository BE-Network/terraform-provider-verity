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
	_ resource.Resource                = &verityFabricResource{}
	_ resource.ResourceWithConfigure   = &verityFabricResource{}
	_ resource.ResourceWithImportState = &verityFabricResource{}
	_ resource.ResourceWithModifyPlan  = &verityFabricResource{}
)

const fabricResourceType = "fabrics"
const fabricTerraformType = "verity_fabric"

func NewVerityFabricResource() resource.Resource {
	return &verityFabricResource{}
}

type verityFabricResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type verityFabricResourceModel struct {
	Name                                      types.String                        `tfsdk:"name"`
	Enable                                    types.Bool                          `tfsdk:"enable"`
	PlaneCount                                types.String                        `tfsdk:"plane_count"`
	SuSize                                    types.String                        `tfsdk:"su_size"`
	SuSupport                                 types.Bool                          `tfsdk:"su_support"`
	ServerManagement                          types.Bool                          `tfsdk:"server_management"`
	AllowAllUnderlayConnections               types.Bool                          `tfsdk:"allow_all_underlay_connections"`
	SiteType                                  types.String                        `tfsdk:"site_type"`
	ServiceForSite                            types.String                        `tfsdk:"service_for_site"`
	ServiceForSiteRefType                     types.String                        `tfsdk:"service_for_site_ref_type_"`
	PortAdminPollingInterval                  types.Int64                         `tfsdk:"port_admin_polling_interval"`
	PortStatusPollingInterval                 types.Int64                         `tfsdk:"port_status_polling_interval"`
	SpanningTreeType                          types.String                        `tfsdk:"spanning_tree_type"`
	RegionName                                types.String                        `tfsdk:"region_name"`
	Revision                                  types.Int64                         `tfsdk:"revision"`
	ForceSpanningTreeOnFabricPorts            types.Bool                          `tfsdk:"force_spanning_tree_on_fabric_ports"`
	ReadOnlyMode                              types.Bool                          `tfsdk:"read_only_mode"`
	DomainForSite                             types.String                        `tfsdk:"domain_for_site"`
	DomainForSiteRefType                      types.String                        `tfsdk:"domain_for_site_ref_type_"`
	EnableDscp                                types.Bool                          `tfsdk:"enable_dscp"`
	DscpToPBitMap                             types.String                        `tfsdk:"dscp_to_p_bit_map"`
	AnycastMacAddress                         types.String                        `tfsdk:"anycast_mac_address"`
	AnycastMacAddressAutoAssigned             types.Bool                          `tfsdk:"anycast_mac_address_auto_assigned_"`
	MacAddressAgingTime                       types.Int64                         `tfsdk:"mac_address_aging_time"`
	MlagDelayRestoreTimer                     types.Int64                         `tfsdk:"mlag_delay_restore_timer"`
	BgpKeepaliveTimer                         types.Int64                         `tfsdk:"bgp_keepalive_timer"`
	BgpHoldDownTimer                          types.Int64                         `tfsdk:"bgp_hold_down_timer"`
	SpineBgpAdvertisementInterval             types.Int64                         `tfsdk:"spine_bgp_advertisement_interval"`
	SpineBgpConnectTimer                      types.Int64                         `tfsdk:"spine_bgp_connect_timer"`
	SpineAsNumber                             types.Int64                         `tfsdk:"spine_as_number"`
	LeafBgpKeepAliveTimer                     types.Int64                         `tfsdk:"leaf_bgp_keep_alive_timer"`
	LeafBgpHoldDownTimer                      types.Int64                         `tfsdk:"leaf_bgp_hold_down_timer"`
	LeafBgpAdvertisementInterval              types.Int64                         `tfsdk:"leaf_bgp_advertisement_interval"`
	LeafBgpConnectTimer                       types.Int64                         `tfsdk:"leaf_bgp_connect_timer"`
	LinkStateTimeoutValue                     types.Int64                         `tfsdk:"link_state_timeout_value"`
	EvpnMultihomingStartupDelay               types.Int64                         `tfsdk:"evpn_multihoming_startup_delay"`
	EvpnMacHoldtime                           types.Int64                         `tfsdk:"evpn_mac_holdtime"`
	AggressiveReporting                       types.Bool                          `tfsdk:"aggressive_reporting"`
	SwitchIpBase                              types.String                        `tfsdk:"switch_ip_base"`
	ControllerIpBase                          types.String                        `tfsdk:"controller_ip_base"`
	SwitchUsername                            types.String                        `tfsdk:"switch_username"`
	SwitchPassword                            types.String                        `tfsdk:"switch_password"`
	SwitchPasswordEncrypted                   types.String                        `tfsdk:"switch_password_encrypted"`
	HgxUsername                               types.String                        `tfsdk:"hgx_username"`
	HgxPassword                               types.String                        `tfsdk:"hgx_password"`
	HgxPasswordEncrypted                      types.String                        `tfsdk:"hgx_password_encrypted"`
	SwitchGateway                             types.String                        `tfsdk:"switch_gateway"`
	ControllerGateway                         types.String                        `tfsdk:"controller_gateway"`
	HgxGateway                                types.String                        `tfsdk:"hgx_gateway"`
	GpuArchitecture                           types.String                        `tfsdk:"gpu_architecture"`
	MultiTenant                               types.Bool                          `tfsdk:"multi_tenant"`
	BaseBgpAsNumber                           types.String                        `tfsdk:"base_bgp_as_number"`
	RouterIdBasePrefix                        types.String                        `tfsdk:"router_id_base_prefix"`
	VtepIdBasePrefix                          types.String                        `tfsdk:"vtep_id_base_prefix"`
	PairedIpSubnet                            types.String                        `tfsdk:"paired_ip_subnet"`
	MaxSwitches                               types.String                        `tfsdk:"max_switches"`
	PauseValidationAlarms                     types.Bool                          `tfsdk:"pause_validation_alarms"`
	StartingOctet                             types.Int64                         `tfsdk:"starting_octet"`
	MaxSus                                    types.Int64                         `tfsdk:"max_sus"`
	MaxPods                                   types.Int64                         `tfsdk:"max_pods"`
	EnableDhcpSnooping                        types.Bool                          `tfsdk:"enable_dhcp_snooping"`
	IpSourceGuard                             types.Bool                          `tfsdk:"ip_source_guard"`
	DuplicateAddressDetectionMaxNumberOfMoves types.Int64                         `tfsdk:"duplicate_address_detection_max_number_of_moves"`
	DuplicateAddressDetectionTime             types.Int64                         `tfsdk:"duplicate_address_detection_time"`
	ObjectProperties                          []verityFabricObjectPropertiesModel `tfsdk:"object_properties"`
}

type verityFabricObjectPropertiesModel struct {
	SystemGraphs []verityFabricSystemGraphsModel `tfsdk:"system_graphs"`
}

type verityFabricSystemGraphsModel struct {
	Index types.Int64 `tfsdk:"index"`
}

func (m verityFabricSystemGraphsModel) GetIndex() types.Int64 {
	return m.Index
}

func (r *verityFabricResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fabric"
}

func (r *verityFabricResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityFabricResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity Fabric",
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
			"plane_count": schema.StringAttribute{
				Description: "Number of planes in this Fabric",
				Optional:    true,
				Computed:    true,
			},
			"su_size": schema.StringAttribute{
				Description: "Number of HGXs per SU. Valid values are \"32\" and \"64\".",
				Optional:    true,
				Computed:    true,
			},
			"su_support": schema.BoolAttribute{
				Description: "Support grouping leaf switches in SUs",
				Optional:    true,
				Computed:    true,
			},
			"server_management": schema.BoolAttribute{
				Description: "Support managing servers",
				Optional:    true,
				Computed:    true,
			},
			"allow_all_underlay_connections": schema.BoolAttribute{
				Description: "Allows underlay connections between PODs",
				Optional:    true,
				Computed:    true,
			},
			"site_type": schema.StringAttribute{
				Description: "Type of Fabric",
				Optional:    true,
				Computed:    true,
			},
			"service_for_site": schema.StringAttribute{
				Description: "Service for Fabric",
				Optional:    true,
				Computed:    true,
			},
			"service_for_site_ref_type_": schema.StringAttribute{
				Description: "Object type for service_for_site field",
				Optional:    true,
				Computed:    true,
			},
			"port_admin_polling_interval": schema.Int64Attribute{
				Description: "Polling interval value in seconds used when aggressive reporting is disabled",
				Optional:    true,
				Computed:    true,
			},
			"port_status_polling_interval": schema.Int64Attribute{
				Description: "Polling interval value in seconds used when aggressive reporting is disabled",
				Optional:    true,
				Computed:    true,
			},
			"spanning_tree_type": schema.StringAttribute{
				Description: "Sets the spanning tree type for all Ports in this Fabric with Spanning Tree enabled",
				Optional:    true,
				Computed:    true,
			},
			"region_name": schema.StringAttribute{
				Description: "Defines the logical boundary of the network. All switches in an MSTP region must have the same configured region name",
				Optional:    true,
				Computed:    true,
			},
			"revision": schema.Int64Attribute{
				Description: "A logical number that signifies a revision for the MSTP configuration. All switches in an MSTP region must have the same revision number (maximum: 65535)",
				Optional:    true,
				Computed:    true,
			},
			"force_spanning_tree_on_fabric_ports": schema.BoolAttribute{
				Description: "Enable spanning tree on all fabric connections. This overrides the Eth Port Settings for Fabric ports",
				Optional:    true,
				Computed:    true,
			},
			"read_only_mode": schema.BoolAttribute{
				Description: "When Read Only Mode is checked, vNetC will perform all functions except writing database updates to the target hardware",
				Optional:    true,
				Computed:    true,
			},
			"domain_for_site": schema.StringAttribute{
				Description: "Fabric Collection for Fabric",
				Optional:    true,
				Computed:    true,
			},
			"domain_for_site_ref_type_": schema.StringAttribute{
				Description: "Object type for domain_for_site field",
				Optional:    true,
				Computed:    true,
			},
			"enable_dscp": schema.BoolAttribute{
				Description: "Enable DSCP to p-bit/TC configuration. When enabled, DSCP to p-bit/TC mappings are applied.",
				Optional:    true,
				Computed:    true,
			},
			"dscp_to_p_bit_map": schema.StringAttribute{
				Description: "For any Service that is using DSCP to TC map packet prioritization. A string of length 64 with a 0-7 in each position (maxLength: 64)",
				Optional:    true,
				Computed:    true,
			},
			"anycast_mac_address": schema.StringAttribute{
				Description: "Anycast MAC address to use. This field should not be specified when 'anycast_mac_address_auto_assigned_' is set to true, as the API will assign this value automatically. Used for MAC VRRP.",
				Optional:    true,
				Computed:    true,
			},
			"anycast_mac_address_auto_assigned_": schema.BoolAttribute{
				Description: "Whether the anycast MAC address should be automatically assigned by the API. When set to true, do not specify the 'anycast_mac_address' field in your configuration.",
				Optional:    true,
				Computed:    true,
			},
			"mac_address_aging_time": schema.Int64Attribute{
				Description: "MAC Address Aging Time (minimum: 1, maximum: 100000)",
				Optional:    true,
				Computed:    true,
			},
			"mlag_delay_restore_timer": schema.Int64Attribute{
				Description: "MLAG Delay Restore Timer (minimum: 1, maximum: 3600)",
				Optional:    true,
				Computed:    true,
			},
			"bgp_keepalive_timer": schema.Int64Attribute{
				Description: "Spine BGP Keepalive Timer (minimum: 1, maximum: 3600)",
				Optional:    true,
				Computed:    true,
			},
			"bgp_hold_down_timer": schema.Int64Attribute{
				Description: "Spine BGP Hold Down Timer (minimum: 1, maximum: 3600)",
				Optional:    true,
				Computed:    true,
			},
			"spine_bgp_advertisement_interval": schema.Int64Attribute{
				Description: "BGP Advertisement Interval for spines/superspines. Use \"0\" for immediate updates (maximum: 3600)",
				Optional:    true,
				Computed:    true,
			},
			"spine_bgp_connect_timer": schema.Int64Attribute{
				Description: "BGP Connect Timer (minimum: 1, maximum: 3600)",
				Optional:    true,
				Computed:    true,
			},
			"spine_as_number": schema.Int64Attribute{
				Description: "BGP AS number applied uniformly to all spine endpoints in this CLOS fabric on save. Leave blank to manage spine AS numbers individually.",
				Optional:    true,
				Computed:    true,
			},
			"leaf_bgp_keep_alive_timer": schema.Int64Attribute{
				Description: "Leaf BGP Keep Alive Timer (minimum: 1, maximum: 3600)",
				Optional:    true,
				Computed:    true,
			},
			"leaf_bgp_hold_down_timer": schema.Int64Attribute{
				Description: "Leaf BGP Hold Down Timer (minimum: 1, maximum: 3600)",
				Optional:    true,
				Computed:    true,
			},
			"leaf_bgp_advertisement_interval": schema.Int64Attribute{
				Description: "BGP Advertisement Interval for leafs. Use \"0\" for immediate updates (maximum: 3600)",
				Optional:    true,
				Computed:    true,
			},
			"leaf_bgp_connect_timer": schema.Int64Attribute{
				Description: "BGP Connect Timer (minimum: 1, maximum: 3600)",
				Optional:    true,
				Computed:    true,
			},
			"link_state_timeout_value": schema.Int64Attribute{
				Description: "Link State Timeout Value",
				Optional:    true,
				Computed:    true,
			},
			"evpn_multihoming_startup_delay": schema.Int64Attribute{
				Description: "Startup Delay",
				Optional:    true,
				Computed:    true,
			},
			"evpn_mac_holdtime": schema.Int64Attribute{
				Description: "MAC Holdtime",
				Optional:    true,
				Computed:    true,
			},
			"aggressive_reporting": schema.BoolAttribute{
				Description: "Fast Reporting of Switch Communications, Link Up/Down, and BGP Status",
				Optional:    true,
				Computed:    true,
			},
			"switch_ip_base": schema.StringAttribute{
				Description: "Base IPv4 address for switch IPs in this Fabric",
				Optional:    true,
				Computed:    true,
			},
			"controller_ip_base": schema.StringAttribute{
				Description: "Base IPv4 address for controller IPs in this Fabric",
				Optional:    true,
				Computed:    true,
			},
			"switch_username": schema.StringAttribute{
				Description: "Default username for managed switches in this Fabric",
				Optional:    true,
				Computed:    true,
			},
			"switch_password": schema.StringAttribute{
				Description: "Default password for managed switches in this Fabric",
				Optional:    true,
				Computed:    true,
			},
			"switch_password_encrypted": schema.StringAttribute{
				Description: "Default password for managed switches in this Fabric",
				Optional:    true,
				Computed:    true,
			},
			"hgx_username": schema.StringAttribute{
				Description: "Default username for HGX devices in this Fabric",
				Optional:    true,
				Computed:    true,
			},
			"hgx_password": schema.StringAttribute{
				Description: "Default password for HGX devices in this Fabric",
				Optional:    true,
				Computed:    true,
			},
			"hgx_password_encrypted": schema.StringAttribute{
				Description: "Default password for HGX devices in this Fabric",
				Optional:    true,
				Computed:    true,
			},
			"switch_gateway": schema.StringAttribute{
				Description: "Default switch management gateway IP for devices in this Fabric",
				Optional:    true,
				Computed:    true,
			},
			"controller_gateway": schema.StringAttribute{
				Description: "Default Device Management VM gateway IP for devices in this Fabric",
				Optional:    true,
				Computed:    true,
			},
			"hgx_gateway": schema.StringAttribute{
				Description: "Default HGX management gateway IP for devices in this Fabric",
				Optional:    true,
				Computed:    true,
			},
			"gpu_architecture": schema.StringAttribute{
				Description: "GPU Architecture used within this Fabric",
				Optional:    true,
				Computed:    true,
			},
			"multi_tenant": schema.BoolAttribute{
				Description: "Allow multiple tenants to HGX endpoints on this fabric.",
				Optional:    true,
				Computed:    true,
			},
			"base_bgp_as_number": schema.StringAttribute{
				Description: "Base BGP Autonomous System Number used for switches in the fabric",
				Optional:    true,
				Computed:    true,
			},
			"router_id_base_prefix": schema.StringAttribute{
				Description: "Router ID starting IP address",
				Optional:    true,
				Computed:    true,
			},
			"vtep_id_base_prefix": schema.StringAttribute{
				Description: "Vtep ID starting IP address",
				Optional:    true,
				Computed:    true,
			},
			"paired_ip_subnet": schema.StringAttribute{
				Description: "IP address range reserved for communication between paired switches",
				Optional:    true,
				Computed:    true,
			},
			"max_switches": schema.StringAttribute{
				Description: "Max number Switches to support in this site",
				Optional:    true,
				Computed:    true,
			},
			"pause_validation_alarms": schema.BoolAttribute{
				Description: "Validation still runs, but validation alarms are not raised for this Fabric while enabled.",
				Optional:    true,
				Computed:    true,
			},
			"starting_octet": schema.Int64Attribute{
				Description: "Starting Octet for HGX Port IPs",
				Optional:    true,
				Computed:    true,
			},
			"max_sus": schema.Int64Attribute{
				Description: "Maximum number of SUs allowed per POD",
				Optional:    true,
				Computed:    true,
			},
			"max_pods": schema.Int64Attribute{
				Description: "Maximum number of PODs allowed in the Fabric",
				Optional:    true,
				Computed:    true,
			},
			"enable_dhcp_snooping": schema.BoolAttribute{
				Description: "Enables the switches to monitor DHCP traffic and collect assigned IP addresses which are then placed in the DHCP assigned IPs report.",
				Optional:    true,
				Computed:    true,
			},
			"ip_source_guard": schema.BoolAttribute{
				Description: "On untrusted ports, only allow known traffic from known IP addresses. IP addresses are discovered via DHCP snooping or with static IP settings",
				Optional:    true,
				Computed:    true,
			},
			"duplicate_address_detection_max_number_of_moves": schema.Int64Attribute{
				Description: "Duplicate Address Detection Max Number of Moves",
				Optional:    true,
				Computed:    true,
			},
			"duplicate_address_detection_time": schema.Int64Attribute{
				Description: "Duplicate Address Detection Time",
				Optional:    true,
				Computed:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"object_properties": schema.ListNestedBlock{
				Description: "Object properties for the Fabric",
				NestedObject: schema.NestedBlockObject{
					Blocks: map[string]schema.Block{
						"system_graphs": schema.ListNestedBlock{
							Description: "System graphs",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"index": schema.Int64Attribute{
										Description: "The index identifying the object. Zero if you want to add an object to the list.",
										Optional:    true,
										Computed:    true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *verityFabricResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityFabricResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config verityFabricResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !plan.AnycastMacAddressAutoAssigned.IsNull() && plan.AnycastMacAddressAutoAssigned.ValueBool() {
		if !plan.AnycastMacAddress.IsNull() && !plan.AnycastMacAddress.IsUnknown() && plan.AnycastMacAddress.ValueString() != "" {
			resp.Diagnostics.AddError(
				"Anycast MAC Address cannot be specified when auto-assigned",
				"The 'anycast_mac_address' field cannot be specified in the configuration when 'anycast_mac_address_auto_assigned_' is set to true. The API will assign this value automatically.",
			)
			return
		}
	}

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Authenticate",
			fmt.Sprintf("Error authenticating with API: %v", err),
		)
		return
	}

	name := plan.Name.ValueString()
	fabricReq := &openapi.FabricsPutRequestFabricValue{
		Name: openapi.PtrString(name),
	}

	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "PlaneCount", APIField: &fabricReq.PlaneCount, TFValue: plan.PlaneCount},
		{FieldName: "SuSize", APIField: &fabricReq.SuSize, TFValue: plan.SuSize},
		{FieldName: "SiteType", APIField: &fabricReq.SiteType, TFValue: plan.SiteType},
		{FieldName: "ServiceForSite", APIField: &fabricReq.ServiceForSite, TFValue: plan.ServiceForSite},
		{FieldName: "ServiceForSiteRefType", APIField: &fabricReq.ServiceForSiteRefType, TFValue: plan.ServiceForSiteRefType},
		{FieldName: "SpanningTreeType", APIField: &fabricReq.SpanningTreeType, TFValue: plan.SpanningTreeType},
		{FieldName: "RegionName", APIField: &fabricReq.RegionName, TFValue: plan.RegionName},
		{FieldName: "DomainForSite", APIField: &fabricReq.DomainForSite, TFValue: plan.DomainForSite},
		{FieldName: "DomainForSiteRefType", APIField: &fabricReq.DomainForSiteRefType, TFValue: plan.DomainForSiteRefType},
		{FieldName: "DscpToPBitMap", APIField: &fabricReq.DscpToPBitMap, TFValue: plan.DscpToPBitMap},
		{FieldName: "SwitchIpBase", APIField: &fabricReq.SwitchIpBase, TFValue: plan.SwitchIpBase},
		{FieldName: "ControllerIpBase", APIField: &fabricReq.ControllerIpBase, TFValue: plan.ControllerIpBase},
		{FieldName: "SwitchUsername", APIField: &fabricReq.SwitchUsername, TFValue: plan.SwitchUsername},
		{FieldName: "SwitchPassword", APIField: &fabricReq.SwitchPassword, TFValue: plan.SwitchPassword},
		{FieldName: "SwitchPasswordEncrypted", APIField: &fabricReq.SwitchPasswordEncrypted, TFValue: plan.SwitchPasswordEncrypted},
		{FieldName: "HgxUsername", APIField: &fabricReq.HgxUsername, TFValue: plan.HgxUsername},
		{FieldName: "HgxPassword", APIField: &fabricReq.HgxPassword, TFValue: plan.HgxPassword},
		{FieldName: "HgxPasswordEncrypted", APIField: &fabricReq.HgxPasswordEncrypted, TFValue: plan.HgxPasswordEncrypted},
		{FieldName: "SwitchGateway", APIField: &fabricReq.SwitchGateway, TFValue: plan.SwitchGateway},
		{FieldName: "ControllerGateway", APIField: &fabricReq.ControllerGateway, TFValue: plan.ControllerGateway},
		{FieldName: "HgxGateway", APIField: &fabricReq.HgxGateway, TFValue: plan.HgxGateway},
		{FieldName: "GpuArchitecture", APIField: &fabricReq.GpuArchitecture, TFValue: plan.GpuArchitecture},
		{FieldName: "BaseBgpAsNumber", APIField: &fabricReq.BaseBgpAsNumber, TFValue: plan.BaseBgpAsNumber},
		{FieldName: "RouterIdBasePrefix", APIField: &fabricReq.RouterIdBasePrefix, TFValue: plan.RouterIdBasePrefix},
		{FieldName: "VtepIdBasePrefix", APIField: &fabricReq.VtepIdBasePrefix, TFValue: plan.VtepIdBasePrefix},
		{FieldName: "PairedIpSubnet", APIField: &fabricReq.PairedIpSubnet, TFValue: plan.PairedIpSubnet},
		{FieldName: "MaxSwitches", APIField: &fabricReq.MaxSwitches, TFValue: plan.MaxSwitches},
	})

	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &fabricReq.Enable, TFValue: plan.Enable},
		{FieldName: "ServerManagement", APIField: &fabricReq.ServerManagement, TFValue: plan.ServerManagement},
		{FieldName: "SuSupport", APIField: &fabricReq.SuSupport, TFValue: plan.SuSupport},
		{FieldName: "AllowAllUnderlayConnections", APIField: &fabricReq.AllowAllUnderlayConnections, TFValue: plan.AllowAllUnderlayConnections},
		{FieldName: "ForceSpanningTreeOnFabricPorts", APIField: &fabricReq.ForceSpanningTreeOnFabricPorts, TFValue: plan.ForceSpanningTreeOnFabricPorts},
		{FieldName: "ReadOnlyMode", APIField: &fabricReq.ReadOnlyMode, TFValue: plan.ReadOnlyMode},
		{FieldName: "EnableDscp", APIField: &fabricReq.EnableDscp, TFValue: plan.EnableDscp},
		{FieldName: "AggressiveReporting", APIField: &fabricReq.AggressiveReporting, TFValue: plan.AggressiveReporting},
		{FieldName: "MultiTenant", APIField: &fabricReq.MultiTenant, TFValue: plan.MultiTenant},
		{FieldName: "PauseValidationAlarms", APIField: &fabricReq.PauseValidationAlarms, TFValue: plan.PauseValidationAlarms},
		{FieldName: "EnableDhcpSnooping", APIField: &fabricReq.EnableDhcpSnooping, TFValue: plan.EnableDhcpSnooping},
		{FieldName: "IpSourceGuard", APIField: &fabricReq.IpSourceGuard, TFValue: plan.IpSourceGuard},
	})

	workDir := r.provCtx.workDir
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, fabricTerraformType, name)

	utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
		{FieldName: "PortAdminPollingInterval", APIField: &fabricReq.PortAdminPollingInterval, TFValue: config.PortAdminPollingInterval, IsConfigured: configuredAttrs.IsConfigured("port_admin_polling_interval")},
		{FieldName: "PortStatusPollingInterval", APIField: &fabricReq.PortStatusPollingInterval, TFValue: config.PortStatusPollingInterval, IsConfigured: configuredAttrs.IsConfigured("port_status_polling_interval")},
		{FieldName: "MacAddressAgingTime", APIField: &fabricReq.MacAddressAgingTime, TFValue: config.MacAddressAgingTime, IsConfigured: configuredAttrs.IsConfigured("mac_address_aging_time")},
		{FieldName: "MlagDelayRestoreTimer", APIField: &fabricReq.MlagDelayRestoreTimer, TFValue: config.MlagDelayRestoreTimer, IsConfigured: configuredAttrs.IsConfigured("mlag_delay_restore_timer")},
		{FieldName: "BgpKeepaliveTimer", APIField: &fabricReq.BgpKeepaliveTimer, TFValue: config.BgpKeepaliveTimer, IsConfigured: configuredAttrs.IsConfigured("bgp_keepalive_timer")},
		{FieldName: "BgpHoldDownTimer", APIField: &fabricReq.BgpHoldDownTimer, TFValue: config.BgpHoldDownTimer, IsConfigured: configuredAttrs.IsConfigured("bgp_hold_down_timer")},
		{FieldName: "SpineBgpAdvertisementInterval", APIField: &fabricReq.SpineBgpAdvertisementInterval, TFValue: config.SpineBgpAdvertisementInterval, IsConfigured: configuredAttrs.IsConfigured("spine_bgp_advertisement_interval")},
		{FieldName: "SpineBgpConnectTimer", APIField: &fabricReq.SpineBgpConnectTimer, TFValue: config.SpineBgpConnectTimer, IsConfigured: configuredAttrs.IsConfigured("spine_bgp_connect_timer")},
		{FieldName: "SpineAsNumber", APIField: &fabricReq.SpineAsNumber, TFValue: config.SpineAsNumber, IsConfigured: configuredAttrs.IsConfigured("spine_as_number")},
		{FieldName: "LeafBgpKeepAliveTimer", APIField: &fabricReq.LeafBgpKeepAliveTimer, TFValue: config.LeafBgpKeepAliveTimer, IsConfigured: configuredAttrs.IsConfigured("leaf_bgp_keep_alive_timer")},
		{FieldName: "LeafBgpHoldDownTimer", APIField: &fabricReq.LeafBgpHoldDownTimer, TFValue: config.LeafBgpHoldDownTimer, IsConfigured: configuredAttrs.IsConfigured("leaf_bgp_hold_down_timer")},
		{FieldName: "LeafBgpAdvertisementInterval", APIField: &fabricReq.LeafBgpAdvertisementInterval, TFValue: config.LeafBgpAdvertisementInterval, IsConfigured: configuredAttrs.IsConfigured("leaf_bgp_advertisement_interval")},
		{FieldName: "LeafBgpConnectTimer", APIField: &fabricReq.LeafBgpConnectTimer, TFValue: config.LeafBgpConnectTimer, IsConfigured: configuredAttrs.IsConfigured("leaf_bgp_connect_timer")},
		{FieldName: "Revision", APIField: &fabricReq.Revision, TFValue: config.Revision, IsConfigured: configuredAttrs.IsConfigured("revision")},
		{FieldName: "LinkStateTimeoutValue", APIField: &fabricReq.LinkStateTimeoutValue, TFValue: config.LinkStateTimeoutValue, IsConfigured: configuredAttrs.IsConfigured("link_state_timeout_value")},
		{FieldName: "EvpnMultihomingStartupDelay", APIField: &fabricReq.EvpnMultihomingStartupDelay, TFValue: config.EvpnMultihomingStartupDelay, IsConfigured: configuredAttrs.IsConfigured("evpn_multihoming_startup_delay")},
		{FieldName: "EvpnMacHoldtime", APIField: &fabricReq.EvpnMacHoldtime, TFValue: config.EvpnMacHoldtime, IsConfigured: configuredAttrs.IsConfigured("evpn_mac_holdtime")},
		{FieldName: "StartingOctet", APIField: &fabricReq.StartingOctet, TFValue: config.StartingOctet, IsConfigured: configuredAttrs.IsConfigured("starting_octet")},
		{FieldName: "MaxSus", APIField: &fabricReq.MaxSus, TFValue: config.MaxSus, IsConfigured: configuredAttrs.IsConfigured("max_sus")},
		{FieldName: "MaxPods", APIField: &fabricReq.MaxPods, TFValue: config.MaxPods, IsConfigured: configuredAttrs.IsConfigured("max_pods")},
		{FieldName: "DuplicateAddressDetectionMaxNumberOfMoves", APIField: &fabricReq.DuplicateAddressDetectionMaxNumberOfMoves, TFValue: config.DuplicateAddressDetectionMaxNumberOfMoves, IsConfigured: configuredAttrs.IsConfigured("duplicate_address_detection_max_number_of_moves")},
		{FieldName: "DuplicateAddressDetectionTime", APIField: &fabricReq.DuplicateAddressDetectionTime, TFValue: config.DuplicateAddressDetectionTime, IsConfigured: configuredAttrs.IsConfigured("duplicate_address_detection_time")},
	})

	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		systemGraphs := make([]openapi.FabricsPutRequestFabricValueObjectPropertiesSystemGraphsInner, len(op.SystemGraphs))
		for i, graph := range op.SystemGraphs {
			graphProps := openapi.FabricsPutRequestFabricValueObjectPropertiesSystemGraphsInner{}
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &graphProps.Index, TFValue: graph.Index},
			})
			systemGraphs[i] = graphProps
		}
		fabricReq.SetObjectProperties(openapi.FabricsPutRequestFabricValueObjectProperties{
			SystemGraphs: systemGraphs,
		})
	}

	if !plan.AnycastMacAddressAutoAssigned.IsNull() && plan.AnycastMacAddressAutoAssigned.ValueBool() {
		fabricReq.AnycastMacAddressAutoAssigned = openapi.PtrBool(true)
		// Don't include the specific MAC in the request
	} else if configuredAttrs.IsConfigured("anycast_mac_address") && !plan.AnycastMacAddress.IsNull() && !plan.AnycastMacAddress.IsUnknown() {
		// Preserve an explicitly configured value, including an explicit empty string.
		fabricReq.AnycastMacAddress = openapi.PtrString(plan.AnycastMacAddress.ValueString())
		if !plan.AnycastMacAddressAutoAssigned.IsNull() {
			fabricReq.AnycastMacAddressAutoAssigned = openapi.PtrBool(plan.AnycastMacAddressAutoAssigned.ValueBool())
		}
	} else if !plan.AnycastMacAddressAutoAssigned.IsNull() {
		// No MAC value was provided, but preserve the user's explicit flag setting
		fabricReq.AnycastMacAddressAutoAssigned = openapi.PtrBool(plan.AnycastMacAddressAutoAssigned.ValueBool())
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "fabric", name, *fabricReq, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Fabric %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "fabrics")

	var minState verityFabricResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if fabricData, exists := bulkMgr.GetResourceResponse("fabric", name); exists {
			state := populateFabricState(ctx, minState, utils.MergeMissingPlanScalars(fabricData, plan, fabricResourceType, r.provCtx.mode), r.provCtx.mode)
			filterFabricIndexedEntries(&state, &plan)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	readReq := resource.ReadRequest{
		State: resp.State,
	}
	readResp := resource.ReadResponse{
		State:       resp.State,
		Diagnostics: resp.Diagnostics,
	}

	postOpCtx := utils.WithPostOperationFallback(ctx, plan, fabricResourceType, r.provCtx.mode)
	r.Read(postOpCtx, readReq, &readResp)
	if readResp.State.Raw.IsNull() {
		_, diags := utils.SetPostOperationFallbackState(postOpCtx, &readResp.State)
		readResp.Diagnostics.Append(diags...)
	}
	resp.State = readResp.State
	resp.Diagnostics = readResp.Diagnostics
}

func (r *verityFabricResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityFabricResourceModel
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

	fabricName := state.Name.ValueString()
	priorState := state

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if fabricData, exists := r.bulkOpsMgr.GetResourceResponse("fabric", fabricName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached fabric data for %s from recent operation", fabricName))
			state = populateFabricState(ctx, state, utils.ApplyPostOperationFallback(ctx, fabricData), r.provCtx.mode)
			filterFabricIndexedEntries(&state, &priorState)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("fabric") {
		tflog.Info(ctx, fmt.Sprintf("Skipping fabric %s verification – trusting recent successful API operation", fabricName))
		if handled, diags := utils.SetPostOperationFallbackState(ctx, &resp.State); handled {
			resp.Diagnostics.Append(diags...)
		}
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching fabrics for verification of %s", fabricName))

	type FabricsResponse struct {
		Fabric map[string]interface{} `json:"fabric"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "fabrics", fabricName,
		func() (FabricsResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch fabrics")
			respAPI, err := r.client.FabricsAPI.FabricsGet(ctx).Execute()
			if err != nil {
				return FabricsResponse{}, fmt.Errorf("error reading fabrics: %v", err)
			}
			defer respAPI.Body.Close()

			var res FabricsResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return FabricsResponse{}, fmt.Errorf("failed to decode fabrics response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched fabrics data with %d fabrics", len(res.Fabric)))
			return res, nil
		}, getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Fabric %s", fabricName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for fabric with name: %s", fabricName))

	fabricData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.Fabric,
		fabricName,
		func(data interface{}) (string, bool) {
			if fabric, ok := data.(map[string]interface{}); ok {
				if name, ok := fabric["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Fabric with name '%s' not found in API response", fabricName))
		resp.State.RemoveResource(ctx)
		return
	}

	fabricMap, ok := fabricData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Fabric Data",
			fmt.Sprintf("Fabric data is not in expected format for %s", fabricName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found fabric '%s' under API key '%s'", fabricName, actualAPIName))

	state = populateFabricState(ctx, state, utils.ApplyPostOperationFallback(ctx, fabricMap), r.provCtx.mode)
	filterFabricIndexedEntries(&state, &priorState)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityFabricResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityFabricResourceModel

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
	var config verityFabricResourceModel
	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate auto-assigned fields - this check prevents ineffective API calls
	// Only error if the auto-assigned flag is enabled AND the user is explicitly setting a value
	// AND the auto-assigned flag itself is not changing (which would be a valid operation)
	// Don't error if the field is unknown (computed during plan recalculation)
	if !plan.AnycastMacAddress.Equal(state.AnycastMacAddress) &&
		!plan.AnycastMacAddress.IsNull() && !plan.AnycastMacAddress.IsUnknown() &&
		!plan.AnycastMacAddressAutoAssigned.IsNull() && plan.AnycastMacAddressAutoAssigned.ValueBool() &&
		plan.AnycastMacAddressAutoAssigned.Equal(state.AnycastMacAddressAutoAssigned) {
		resp.Diagnostics.AddError(
			"Cannot modify auto-assigned field",
			"The 'anycast_mac_address' field cannot be modified because 'anycast_mac_address_auto_assigned_' is set to true.",
		)
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
	fabricReq := openapi.FabricsPutRequestFabricValue{}
	hasChanges := false

	// Parse HCL to detect which fields are explicitly configured
	workDir := r.provCtx.workDir
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, fabricTerraformType, name)

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { fabricReq.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.PlaneCount, state.PlaneCount, func(v *string) { fabricReq.PlaneCount = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.SuSize, state.SuSize, func(v *string) { fabricReq.SuSize = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.SiteType, state.SiteType, func(v *string) { fabricReq.SiteType = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.SpanningTreeType, state.SpanningTreeType, func(v *string) { fabricReq.SpanningTreeType = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.RegionName, state.RegionName, func(v *string) { fabricReq.RegionName = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.DomainForSite, state.DomainForSite, func(v *string) { fabricReq.DomainForSite = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.DscpToPBitMap, state.DscpToPBitMap, func(v *string) { fabricReq.DscpToPBitMap = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.SwitchIpBase, state.SwitchIpBase, func(v *string) { fabricReq.SwitchIpBase = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.ControllerIpBase, state.ControllerIpBase, func(v *string) { fabricReq.ControllerIpBase = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.SwitchUsername, state.SwitchUsername, func(v *string) { fabricReq.SwitchUsername = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.SwitchPassword, state.SwitchPassword, func(v *string) { fabricReq.SwitchPassword = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.SwitchPasswordEncrypted, state.SwitchPasswordEncrypted, func(v *string) { fabricReq.SwitchPasswordEncrypted = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.HgxUsername, state.HgxUsername, func(v *string) { fabricReq.HgxUsername = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.HgxPassword, state.HgxPassword, func(v *string) { fabricReq.HgxPassword = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.HgxPasswordEncrypted, state.HgxPasswordEncrypted, func(v *string) { fabricReq.HgxPasswordEncrypted = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.SwitchGateway, state.SwitchGateway, func(v *string) { fabricReq.SwitchGateway = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.ControllerGateway, state.ControllerGateway, func(v *string) { fabricReq.ControllerGateway = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.HgxGateway, state.HgxGateway, func(v *string) { fabricReq.HgxGateway = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.GpuArchitecture, state.GpuArchitecture, func(v *string) { fabricReq.GpuArchitecture = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.BaseBgpAsNumber, state.BaseBgpAsNumber, func(v *string) { fabricReq.BaseBgpAsNumber = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.RouterIdBasePrefix, state.RouterIdBasePrefix, func(v *string) { fabricReq.RouterIdBasePrefix = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.VtepIdBasePrefix, state.VtepIdBasePrefix, func(v *string) { fabricReq.VtepIdBasePrefix = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.PairedIpSubnet, state.PairedIpSubnet, func(v *string) { fabricReq.PairedIpSubnet = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.MaxSwitches, state.MaxSwitches, func(v *string) { fabricReq.MaxSwitches = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { fabricReq.Enable = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.ServerManagement, state.ServerManagement, func(v *bool) { fabricReq.ServerManagement = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.SuSupport, state.SuSupport, func(v *bool) { fabricReq.SuSupport = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.AllowAllUnderlayConnections, state.AllowAllUnderlayConnections, func(v *bool) { fabricReq.AllowAllUnderlayConnections = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.ForceSpanningTreeOnFabricPorts, state.ForceSpanningTreeOnFabricPorts, func(v *bool) { fabricReq.ForceSpanningTreeOnFabricPorts = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.ReadOnlyMode, state.ReadOnlyMode, func(v *bool) { fabricReq.ReadOnlyMode = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.EnableDscp, state.EnableDscp, func(v *bool) { fabricReq.EnableDscp = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.AggressiveReporting, state.AggressiveReporting, func(v *bool) { fabricReq.AggressiveReporting = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.MultiTenant, state.MultiTenant, func(v *bool) { fabricReq.MultiTenant = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.PauseValidationAlarms, state.PauseValidationAlarms, func(v *bool) { fabricReq.PauseValidationAlarms = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.EnableDhcpSnooping, state.EnableDhcpSnooping, func(v *bool) { fabricReq.EnableDhcpSnooping = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.IpSourceGuard, state.IpSourceGuard, func(v *bool) { fabricReq.IpSourceGuard = v }, &hasChanges)

	// Handle nullable int64 field changes - parse HCL to detect explicit config
	utils.CompareAndSetNullableInt64Field(config.PortAdminPollingInterval, state.PortAdminPollingInterval, configuredAttrs.IsConfigured("port_admin_polling_interval"), func(v *openapi.NullableInt64) { fabricReq.PortAdminPollingInterval = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.PortStatusPollingInterval, state.PortStatusPollingInterval, configuredAttrs.IsConfigured("port_status_polling_interval"), func(v *openapi.NullableInt64) { fabricReq.PortStatusPollingInterval = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.MacAddressAgingTime, state.MacAddressAgingTime, configuredAttrs.IsConfigured("mac_address_aging_time"), func(v *openapi.NullableInt64) { fabricReq.MacAddressAgingTime = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.MlagDelayRestoreTimer, state.MlagDelayRestoreTimer, configuredAttrs.IsConfigured("mlag_delay_restore_timer"), func(v *openapi.NullableInt64) { fabricReq.MlagDelayRestoreTimer = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.BgpKeepaliveTimer, state.BgpKeepaliveTimer, configuredAttrs.IsConfigured("bgp_keepalive_timer"), func(v *openapi.NullableInt64) { fabricReq.BgpKeepaliveTimer = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.BgpHoldDownTimer, state.BgpHoldDownTimer, configuredAttrs.IsConfigured("bgp_hold_down_timer"), func(v *openapi.NullableInt64) { fabricReq.BgpHoldDownTimer = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.SpineBgpAdvertisementInterval, state.SpineBgpAdvertisementInterval, configuredAttrs.IsConfigured("spine_bgp_advertisement_interval"), func(v *openapi.NullableInt64) { fabricReq.SpineBgpAdvertisementInterval = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.SpineBgpConnectTimer, state.SpineBgpConnectTimer, configuredAttrs.IsConfigured("spine_bgp_connect_timer"), func(v *openapi.NullableInt64) { fabricReq.SpineBgpConnectTimer = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.SpineAsNumber, state.SpineAsNumber, configuredAttrs.IsConfigured("spine_as_number"), func(v *openapi.NullableInt64) { fabricReq.SpineAsNumber = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.LeafBgpKeepAliveTimer, state.LeafBgpKeepAliveTimer, configuredAttrs.IsConfigured("leaf_bgp_keep_alive_timer"), func(v *openapi.NullableInt64) { fabricReq.LeafBgpKeepAliveTimer = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.LeafBgpHoldDownTimer, state.LeafBgpHoldDownTimer, configuredAttrs.IsConfigured("leaf_bgp_hold_down_timer"), func(v *openapi.NullableInt64) { fabricReq.LeafBgpHoldDownTimer = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.LeafBgpAdvertisementInterval, state.LeafBgpAdvertisementInterval, configuredAttrs.IsConfigured("leaf_bgp_advertisement_interval"), func(v *openapi.NullableInt64) { fabricReq.LeafBgpAdvertisementInterval = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.LeafBgpConnectTimer, state.LeafBgpConnectTimer, configuredAttrs.IsConfigured("leaf_bgp_connect_timer"), func(v *openapi.NullableInt64) { fabricReq.LeafBgpConnectTimer = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.Revision, state.Revision, configuredAttrs.IsConfigured("revision"), func(v *openapi.NullableInt64) { fabricReq.Revision = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.LinkStateTimeoutValue, state.LinkStateTimeoutValue, configuredAttrs.IsConfigured("link_state_timeout_value"), func(v *openapi.NullableInt64) { fabricReq.LinkStateTimeoutValue = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.EvpnMultihomingStartupDelay, state.EvpnMultihomingStartupDelay, configuredAttrs.IsConfigured("evpn_multihoming_startup_delay"), func(v *openapi.NullableInt64) { fabricReq.EvpnMultihomingStartupDelay = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.EvpnMacHoldtime, state.EvpnMacHoldtime, configuredAttrs.IsConfigured("evpn_mac_holdtime"), func(v *openapi.NullableInt64) { fabricReq.EvpnMacHoldtime = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.StartingOctet, state.StartingOctet, configuredAttrs.IsConfigured("starting_octet"), func(v *openapi.NullableInt64) { fabricReq.StartingOctet = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.MaxSus, state.MaxSus, configuredAttrs.IsConfigured("max_sus"), func(v *openapi.NullableInt64) { fabricReq.MaxSus = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.MaxPods, state.MaxPods, configuredAttrs.IsConfigured("max_pods"), func(v *openapi.NullableInt64) { fabricReq.MaxPods = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.DuplicateAddressDetectionMaxNumberOfMoves, state.DuplicateAddressDetectionMaxNumberOfMoves, configuredAttrs.IsConfigured("duplicate_address_detection_max_number_of_moves"), func(v *openapi.NullableInt64) { fabricReq.DuplicateAddressDetectionMaxNumberOfMoves = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.DuplicateAddressDetectionTime, state.DuplicateAddressDetectionTime, configuredAttrs.IsConfigured("duplicate_address_detection_time"), func(v *openapi.NullableInt64) { fabricReq.DuplicateAddressDetectionTime = *v }, &hasChanges)

	// Handle object properties with nested system_graphs
	if len(plan.ObjectProperties) > 0 || len(state.ObjectProperties) > 0 {
		var op verityFabricObjectPropertiesModel
		var st verityFabricObjectPropertiesModel
		if len(plan.ObjectProperties) > 0 {
			op = plan.ObjectProperties[0]
		}
		if len(state.ObjectProperties) > 0 {
			st = state.ObjectProperties[0]
		}

		changedSystemGraphs, systemGraphsChanged := utils.ProcessIndexedArrayUpdates(op.SystemGraphs, st.SystemGraphs,
			utils.IndexedItemHandler[verityFabricSystemGraphsModel, openapi.FabricsPutRequestFabricValueObjectPropertiesSystemGraphsInner]{
				CreateNew: func(planItem verityFabricSystemGraphsModel) openapi.FabricsPutRequestFabricValueObjectPropertiesSystemGraphsInner {
					graphProps := openapi.FabricsPutRequestFabricValueObjectPropertiesSystemGraphsInner{}
					utils.SetInt64Fields([]utils.Int64FieldMapping{
						{FieldName: "Index", APIField: &graphProps.Index, TFValue: planItem.Index},
					})
					return graphProps
				},
				UpdateExisting: func(planItem verityFabricSystemGraphsModel, stateItem verityFabricSystemGraphsModel) (openapi.FabricsPutRequestFabricValueObjectPropertiesSystemGraphsInner, bool) {
					graphProps := openapi.FabricsPutRequestFabricValueObjectPropertiesSystemGraphsInner{}
					// Always include index — API requires it to identify which array element to modify
					utils.SetInt64Fields([]utils.Int64FieldMapping{
						{FieldName: "Index", APIField: &graphProps.Index, TFValue: planItem.Index},
					})
					return graphProps, false
				},
				CreateDeleted: func(index int64) openapi.FabricsPutRequestFabricValueObjectPropertiesSystemGraphsInner {
					return openapi.FabricsPutRequestFabricValueObjectPropertiesSystemGraphsInner{
						Index: openapi.PtrInt64(int64(index)),
					}
				},
			})

		if systemGraphsChanged {
			fabricObjProps := openapi.FabricsPutRequestFabricValueObjectProperties{
				SystemGraphs: changedSystemGraphs,
			}
			fabricReq.SetObjectProperties(fabricObjProps)
			hasChanges = true
		}
	}

	// Handle service_for_site and service_for_site_ref_type_ fields using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.ServiceForSite, state.ServiceForSite, plan.ServiceForSiteRefType, state.ServiceForSiteRefType,
		func(v *string) { fabricReq.ServiceForSite = v },
		func(v *string) { fabricReq.ServiceForSiteRefType = v },
		"service_for_site", "service_for_site_ref_type_",
		&hasChanges, &resp.Diagnostics,
	) {
		return
	}

	if !utils.HandleOneRefTypeSupported(
		plan.DomainForSite, state.DomainForSite, plan.DomainForSiteRefType, state.DomainForSiteRefType,
		func(v *string) { fabricReq.DomainForSite = v },
		func(v *string) { fabricReq.DomainForSiteRefType = v },
		"domain_for_site", "domain_for_site_ref_type_",
		&hasChanges, &resp.Diagnostics,
	) {
		return
	}

	// Handle AnycastMacAddress and AnycastMacAddressAutoAssigned changes
	anycastMacAddressChanged := !plan.AnycastMacAddress.IsUnknown() && !plan.AnycastMacAddress.Equal(state.AnycastMacAddress)
	anycastMacAddressAutoAssignedChanged := !plan.AnycastMacAddressAutoAssigned.Equal(state.AnycastMacAddressAutoAssigned)

	if anycastMacAddressChanged || anycastMacAddressAutoAssignedChanged {
		// Handle AnycastMacAddress field changes
		if anycastMacAddressChanged {
			if !plan.AnycastMacAddress.IsNull() && plan.AnycastMacAddress.ValueString() != "" {
				fabricReq.AnycastMacAddress = openapi.PtrString(plan.AnycastMacAddress.ValueString())
			} else {
				fabricReq.AnycastMacAddress = openapi.PtrString("")
			}
		}

		// Handle AnycastMacAddressAutoAssigned field changes
		if anycastMacAddressAutoAssignedChanged {
			// Only send anycast_mac_address_auto_assigned_ if the user has explicitly specified it in their configuration
			var config verityFabricResourceModel
			userSpecifiedAnycastMacAddressAutoAssigned := false
			if !req.Config.Raw.IsNull() {
				if err := req.Config.Get(ctx, &config); err == nil {
					userSpecifiedAnycastMacAddressAutoAssigned = !config.AnycastMacAddressAutoAssigned.IsNull()
				}
			}

			if userSpecifiedAnycastMacAddressAutoAssigned {
				fabricReq.AnycastMacAddressAutoAssigned = openapi.PtrBool(plan.AnycastMacAddressAutoAssigned.ValueBool())

				// Special case: When changing from auto-assigned (true) to manual (false),
				// the API requires both anycast_mac_address_auto_assigned_ and anycast_mac_address fields to be sent.
				if !state.AnycastMacAddressAutoAssigned.IsNull() && state.AnycastMacAddressAutoAssigned.ValueBool() &&
					!plan.AnycastMacAddressAutoAssigned.ValueBool() {
					// Changing from auto-assigned=true to auto-assigned=false
					// Must include AnycastMacAddress value in the request for the change to take effect
					if !plan.AnycastMacAddress.IsNull() && plan.AnycastMacAddress.ValueString() != "" {
						fabricReq.AnycastMacAddress = openapi.PtrString(plan.AnycastMacAddress.ValueString())
					} else if !state.AnycastMacAddress.IsNull() && state.AnycastMacAddress.ValueString() != "" {
						// Use current state AnycastMacAddress if plan doesn't specify one
						fabricReq.AnycastMacAddress = openapi.PtrString(state.AnycastMacAddress.ValueString())
					}
				}
			}
		} else if anycastMacAddressChanged {
			// AnycastMacAddress changed but AnycastMacAddressAutoAssigned didn't change
			// Send the auto-assigned flag to maintain consistency with API
			if !plan.AnycastMacAddressAutoAssigned.IsNull() {
				fabricReq.AnycastMacAddressAutoAssigned = openapi.PtrBool(plan.AnycastMacAddressAutoAssigned.ValueBool())
			} else if !state.AnycastMacAddressAutoAssigned.IsNull() {
				fabricReq.AnycastMacAddressAutoAssigned = openapi.PtrBool(state.AnycastMacAddressAutoAssigned.ValueBool())
			} else {
				fabricReq.AnycastMacAddressAutoAssigned = openapi.PtrBool(false)
			}
		}

		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "fabric", name, fabricReq, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Fabric %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "fabrics")

	var minState verityFabricResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if fabricData, exists := bulkMgr.GetResourceResponse("fabric", name); exists {
			state := populateFabricState(ctx, minState, utils.MergeMissingPlanScalars(fabricData, plan, fabricResourceType, r.provCtx.mode), r.provCtx.mode)
			filterFabricIndexedEntries(&state, &plan)
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

	postOpCtx := utils.WithPostOperationFallback(ctx, plan, fabricResourceType, r.provCtx.mode)
	r.Read(postOpCtx, readReq, &readResp)
	if readResp.State.Raw.IsNull() {
		_, diags := utils.SetPostOperationFallbackState(postOpCtx, &readResp.State)
		readResp.Diagnostics.Append(diags...)
	}
	resp.State = readResp.State
	resp.Diagnostics = readResp.Diagnostics

	if !resp.Diagnostics.HasError() {
		var readState verityFabricResourceModel
		readResp.State.Get(ctx, &readState)
		filterFabricIndexedEntries(&readState, &plan)
		resp.State.Set(ctx, &readState)
	}
}

func (r *verityFabricResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityFabricResourceModel
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "fabric", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Fabric %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "fabrics")
	resp.State.RemoveResource(ctx)
}

func (r *verityFabricResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func populateFabricState(ctx context.Context, state verityFabricResourceModel, fabricData map[string]interface{}, mode string) verityFabricResourceModel {
	const resourceType = fabricResourceType

	state.Name = utils.MapStringFromAPI(fabricData["name"])

	// Int fields
	state.Revision = utils.MapInt64WithMode(fabricData, "revision", resourceType, mode)
	state.PortAdminPollingInterval = utils.MapInt64WithMode(fabricData, "port_admin_polling_interval", resourceType, mode)
	state.PortStatusPollingInterval = utils.MapInt64WithMode(fabricData, "port_status_polling_interval", resourceType, mode)
	state.MacAddressAgingTime = utils.MapInt64WithMode(fabricData, "mac_address_aging_time", resourceType, mode)
	state.MlagDelayRestoreTimer = utils.MapInt64WithMode(fabricData, "mlag_delay_restore_timer", resourceType, mode)
	state.BgpKeepaliveTimer = utils.MapInt64WithMode(fabricData, "bgp_keepalive_timer", resourceType, mode)
	state.BgpHoldDownTimer = utils.MapInt64WithMode(fabricData, "bgp_hold_down_timer", resourceType, mode)
	state.SpineBgpAdvertisementInterval = utils.MapInt64WithMode(fabricData, "spine_bgp_advertisement_interval", resourceType, mode)
	state.SpineBgpConnectTimer = utils.MapInt64WithMode(fabricData, "spine_bgp_connect_timer", resourceType, mode)
	state.SpineAsNumber = utils.MapInt64WithMode(fabricData, "spine_as_number", resourceType, mode)
	state.LeafBgpKeepAliveTimer = utils.MapInt64WithMode(fabricData, "leaf_bgp_keep_alive_timer", resourceType, mode)
	state.LeafBgpHoldDownTimer = utils.MapInt64WithMode(fabricData, "leaf_bgp_hold_down_timer", resourceType, mode)
	state.LeafBgpAdvertisementInterval = utils.MapInt64WithMode(fabricData, "leaf_bgp_advertisement_interval", resourceType, mode)
	state.LeafBgpConnectTimer = utils.MapInt64WithMode(fabricData, "leaf_bgp_connect_timer", resourceType, mode)
	state.LinkStateTimeoutValue = utils.MapInt64WithMode(fabricData, "link_state_timeout_value", resourceType, mode)
	state.EvpnMultihomingStartupDelay = utils.MapInt64WithMode(fabricData, "evpn_multihoming_startup_delay", resourceType, mode)
	state.EvpnMacHoldtime = utils.MapInt64WithMode(fabricData, "evpn_mac_holdtime", resourceType, mode)
	state.StartingOctet = utils.MapInt64WithMode(fabricData, "starting_octet", resourceType, mode)
	state.MaxSus = utils.MapInt64WithMode(fabricData, "max_sus", resourceType, mode)
	state.MaxPods = utils.MapInt64WithMode(fabricData, "max_pods", resourceType, mode)
	state.DuplicateAddressDetectionMaxNumberOfMoves = utils.MapInt64WithMode(fabricData, "duplicate_address_detection_max_number_of_moves", resourceType, mode)
	state.DuplicateAddressDetectionTime = utils.MapInt64WithMode(fabricData, "duplicate_address_detection_time", resourceType, mode)

	// Bool fields
	state.Enable = utils.MapBoolWithMode(fabricData, "enable", resourceType, mode)
	state.ServerManagement = utils.MapBoolWithMode(fabricData, "server_management", resourceType, mode)
	state.SuSupport = utils.MapBoolWithMode(fabricData, "su_support", resourceType, mode)
	state.AllowAllUnderlayConnections = utils.MapBoolWithMode(fabricData, "allow_all_underlay_connections", resourceType, mode)
	state.ForceSpanningTreeOnFabricPorts = utils.MapBoolWithMode(fabricData, "force_spanning_tree_on_fabric_ports", resourceType, mode)
	state.ReadOnlyMode = utils.MapBoolWithMode(fabricData, "read_only_mode", resourceType, mode)
	state.EnableDscp = utils.MapBoolWithMode(fabricData, "enable_dscp", resourceType, mode)
	state.AggressiveReporting = utils.MapBoolWithMode(fabricData, "aggressive_reporting", resourceType, mode)
	state.MultiTenant = utils.MapBoolWithMode(fabricData, "multi_tenant", resourceType, mode)
	state.PauseValidationAlarms = utils.MapBoolWithMode(fabricData, "pause_validation_alarms", resourceType, mode)
	state.EnableDhcpSnooping = utils.MapBoolWithMode(fabricData, "enable_dhcp_snooping", resourceType, mode)
	state.IpSourceGuard = utils.MapBoolWithMode(fabricData, "ip_source_guard", resourceType, mode)
	state.AnycastMacAddressAutoAssigned = utils.MapBoolWithMode(fabricData, "anycast_mac_address_auto_assigned_", resourceType, mode)

	// String fields
	state.SiteType = utils.MapStringWithMode(fabricData, "site_type", resourceType, mode)
	state.PlaneCount = utils.MapStringWithMode(fabricData, "plane_count", resourceType, mode)
	state.SuSize = utils.MapStringWithMode(fabricData, "su_size", resourceType, mode)
	state.ServiceForSite = utils.MapStringWithMode(fabricData, "service_for_site", resourceType, mode)
	state.ServiceForSiteRefType = utils.MapStringWithMode(fabricData, "service_for_site_ref_type_", resourceType, mode)
	state.SpanningTreeType = utils.MapStringWithMode(fabricData, "spanning_tree_type", resourceType, mode)
	state.RegionName = utils.MapStringWithMode(fabricData, "region_name", resourceType, mode)
	state.DomainForSite = utils.MapStringWithMode(fabricData, "domain_for_site", resourceType, mode)
	state.DomainForSiteRefType = utils.MapStringWithMode(fabricData, "domain_for_site_ref_type_", resourceType, mode)
	state.DscpToPBitMap = utils.MapStringWithMode(fabricData, "dscp_to_p_bit_map", resourceType, mode)
	state.AnycastMacAddress = utils.MapStringWithMode(fabricData, "anycast_mac_address", resourceType, mode)
	state.SwitchIpBase = utils.MapStringWithMode(fabricData, "switch_ip_base", resourceType, mode)
	state.ControllerIpBase = utils.MapStringWithMode(fabricData, "controller_ip_base", resourceType, mode)
	state.SwitchUsername = utils.MapStringWithMode(fabricData, "switch_username", resourceType, mode)
	state.SwitchPassword = utils.MapStringWithMode(fabricData, "switch_password", resourceType, mode)
	state.SwitchPasswordEncrypted = utils.MapStringWithMode(fabricData, "switch_password_encrypted", resourceType, mode)
	state.HgxUsername = utils.MapStringWithMode(fabricData, "hgx_username", resourceType, mode)
	state.HgxPassword = utils.MapStringWithMode(fabricData, "hgx_password", resourceType, mode)
	state.HgxPasswordEncrypted = utils.MapStringWithMode(fabricData, "hgx_password_encrypted", resourceType, mode)
	state.SwitchGateway = utils.MapStringWithMode(fabricData, "switch_gateway", resourceType, mode)
	state.ControllerGateway = utils.MapStringWithMode(fabricData, "controller_gateway", resourceType, mode)
	state.HgxGateway = utils.MapStringWithMode(fabricData, "hgx_gateway", resourceType, mode)
	state.GpuArchitecture = utils.MapStringWithMode(fabricData, "gpu_architecture", resourceType, mode)
	state.BaseBgpAsNumber = utils.MapStringWithMode(fabricData, "base_bgp_as_number", resourceType, mode)
	state.RouterIdBasePrefix = utils.MapStringWithMode(fabricData, "router_id_base_prefix", resourceType, mode)
	state.VtepIdBasePrefix = utils.MapStringWithMode(fabricData, "vtep_id_base_prefix", resourceType, mode)
	state.PairedIpSubnet = utils.MapStringWithMode(fabricData, "paired_ip_subnet", resourceType, mode)
	state.MaxSwitches = utils.MapStringWithMode(fabricData, "max_switches", resourceType, mode)

	// Handle object_properties block
	if utils.FieldAppliesToMode(resourceType, "object_properties", mode) {
		if op, ok := fabricData["object_properties"].(map[string]interface{}); ok {
			objProps := verityFabricObjectPropertiesModel{}

			// Handle nested system_graphs array
			if systemGraphs, exists := op["system_graphs"].([]interface{}); exists && len(systemGraphs) > 0 {
				var graphsList []verityFabricSystemGraphsModel
				for _, graph := range systemGraphs {
					graphMap, ok := graph.(map[string]interface{})
					if !ok {
						continue
					}
					graphModel := verityFabricSystemGraphsModel{
						Index: utils.MapInt64WithModeNested(graphMap, "index", resourceType, "object_properties.system_graphs.index", mode),
					}
					graphsList = append(graphsList, graphModel)
				}
				objProps.SystemGraphs = graphsList
			} else {
				objProps.SystemGraphs = []verityFabricSystemGraphsModel{}
			}

			state.ObjectProperties = []verityFabricObjectPropertiesModel{objProps}
		} else {
			state.ObjectProperties = nil
		}
	} else {
		state.ObjectProperties = nil
	}

	return state
}

func (r *verityFabricResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// =========================================================================
	// Skip if deleting
	// =========================================================================
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityFabricResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Mode-aware field nullification
	// Set fields that don't apply to current mode to null to prevent
	// "known after apply" messages for irrelevant fields.
	// =========================================================================
	const resourceType = fabricResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyStrings(
		"site_type", "plane_count", "su_size",
		"service_for_site", "service_for_site_ref_type_",
		"spanning_tree_type", "region_name",
		"domain_for_site", "domain_for_site_ref_type_",
		"dscp_to_p_bit_map", "anycast_mac_address",
		"switch_ip_base", "controller_ip_base",
		"switch_username", "switch_password", "switch_password_encrypted",
		"hgx_username", "hgx_password", "hgx_password_encrypted",
		"switch_gateway", "controller_gateway", "hgx_gateway", "gpu_architecture",
		"base_bgp_as_number", "router_id_base_prefix",
		"vtep_id_base_prefix", "paired_ip_subnet", "max_switches",
	)

	nullifier.NullifyBools(
		"su_support", "server_management", "allow_all_underlay_connections",
		"enable", "force_spanning_tree_on_fabric_ports",
		"read_only_mode", "enable_dscp", "aggressive_reporting",
		"multi_tenant", "pause_validation_alarms",
		"enable_dhcp_snooping", "ip_source_guard",
		"anycast_mac_address_auto_assigned_",
	)

	nullifier.NullifyInt64s(
		"revision", "port_admin_polling_interval",
		"port_status_polling_interval", "mac_address_aging_time",
		"mlag_delay_restore_timer", "bgp_keepalive_timer",
		"bgp_hold_down_timer", "spine_bgp_advertisement_interval",
		"spine_bgp_connect_timer", "spine_as_number",
		"leaf_bgp_keep_alive_timer",
		"leaf_bgp_hold_down_timer", "leaf_bgp_advertisement_interval",
		"leaf_bgp_connect_timer", "link_state_timeout_value",
		"evpn_multihoming_startup_delay", "evpn_mac_holdtime",
		"starting_octet", "max_sus", "max_pods",
		"duplicate_address_detection_max_number_of_moves",
		"duplicate_address_detection_time",
	)

	// Handle object_properties with nested system_graphs sub-block
	// Build item counts for system_graphs within each object_properties item
	systemGraphsCounts := make([]int, len(plan.ObjectProperties))
	for i, op := range plan.ObjectProperties {
		systemGraphsCounts[i] = len(op.SystemGraphs)
	}
	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName: "object_properties",
		ItemCount: len(plan.ObjectProperties),
		SubBlocks: []utils.SubBlockFieldConfig{
			{
				SubBlockName: "system_graphs",
				ItemCounts:   systemGraphsCounts,
			},
		},
	})

	// =========================================================================
	// CREATE operation - handle auto-assigned fields
	// =========================================================================
	if req.State.Raw.IsNull() {
		if !plan.AnycastMacAddressAutoAssigned.IsNull() && plan.AnycastMacAddressAutoAssigned.ValueBool() {
			if !plan.AnycastMacAddress.IsNull() && !plan.AnycastMacAddress.IsUnknown() && plan.AnycastMacAddress.ValueString() != "" {
				resp.Diagnostics.AddError(
					"Anycast MAC Address cannot be specified when auto-assigned",
					"The 'anycast_mac_address' field cannot be specified in the configuration when 'anycast_mac_address_auto_assigned_' is set to true. The API will assign this value automatically.",
				)
				return
			}
		}

		// Fabric-specific: AnycastMacAddress auto-assignment on create
		if !plan.AnycastMacAddressAutoAssigned.IsNull() && plan.AnycastMacAddressAutoAssigned.ValueBool() {
			resp.Plan.SetAttribute(ctx, path.Root("anycast_mac_address"), types.StringUnknown())
		}
		return
	}

	// =========================================================================
	// UPDATE operation - get state and config
	// =========================================================================
	var state verityFabricResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config verityFabricResourceModel
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
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, fabricTerraformType, name)

	utils.HandleNullableFields(utils.NullableFieldsConfig{
		Ctx:             ctx,
		Plan:            &resp.Plan,
		ConfiguredAttrs: configuredAttrs,
		Int64Fields: []utils.NullableInt64Field{
			{AttrName: "revision", ConfigVal: config.Revision, StateVal: state.Revision},
			{AttrName: "port_admin_polling_interval", ConfigVal: config.PortAdminPollingInterval, StateVal: state.PortAdminPollingInterval},
			{AttrName: "port_status_polling_interval", ConfigVal: config.PortStatusPollingInterval, StateVal: state.PortStatusPollingInterval},
			{AttrName: "mac_address_aging_time", ConfigVal: config.MacAddressAgingTime, StateVal: state.MacAddressAgingTime},
			{AttrName: "mlag_delay_restore_timer", ConfigVal: config.MlagDelayRestoreTimer, StateVal: state.MlagDelayRestoreTimer},
			{AttrName: "bgp_keepalive_timer", ConfigVal: config.BgpKeepaliveTimer, StateVal: state.BgpKeepaliveTimer},
			{AttrName: "bgp_hold_down_timer", ConfigVal: config.BgpHoldDownTimer, StateVal: state.BgpHoldDownTimer},
			{AttrName: "spine_bgp_advertisement_interval", ConfigVal: config.SpineBgpAdvertisementInterval, StateVal: state.SpineBgpAdvertisementInterval},
			{AttrName: "spine_bgp_connect_timer", ConfigVal: config.SpineBgpConnectTimer, StateVal: state.SpineBgpConnectTimer},
			{AttrName: "spine_as_number", ConfigVal: config.SpineAsNumber, StateVal: state.SpineAsNumber},
			{AttrName: "leaf_bgp_keep_alive_timer", ConfigVal: config.LeafBgpKeepAliveTimer, StateVal: state.LeafBgpKeepAliveTimer},
			{AttrName: "leaf_bgp_hold_down_timer", ConfigVal: config.LeafBgpHoldDownTimer, StateVal: state.LeafBgpHoldDownTimer},
			{AttrName: "leaf_bgp_advertisement_interval", ConfigVal: config.LeafBgpAdvertisementInterval, StateVal: state.LeafBgpAdvertisementInterval},
			{AttrName: "leaf_bgp_connect_timer", ConfigVal: config.LeafBgpConnectTimer, StateVal: state.LeafBgpConnectTimer},
			{AttrName: "link_state_timeout_value", ConfigVal: config.LinkStateTimeoutValue, StateVal: state.LinkStateTimeoutValue},
			{AttrName: "evpn_multihoming_startup_delay", ConfigVal: config.EvpnMultihomingStartupDelay, StateVal: state.EvpnMultihomingStartupDelay},
			{AttrName: "evpn_mac_holdtime", ConfigVal: config.EvpnMacHoldtime, StateVal: state.EvpnMacHoldtime},
			{AttrName: "starting_octet", ConfigVal: config.StartingOctet, StateVal: state.StartingOctet},
			{AttrName: "max_sus", ConfigVal: config.MaxSus, StateVal: state.MaxSus},
			{AttrName: "max_pods", ConfigVal: config.MaxPods, StateVal: state.MaxPods},
			{AttrName: "duplicate_address_detection_max_number_of_moves", ConfigVal: config.DuplicateAddressDetectionMaxNumberOfMoves, StateVal: state.DuplicateAddressDetectionMaxNumberOfMoves},
			{AttrName: "duplicate_address_detection_time", ConfigVal: config.DuplicateAddressDetectionTime, StateVal: state.DuplicateAddressDetectionTime},
		},
	})

	// =========================================================================
	// Validate auto-assigned field specifications
	// =========================================================================
	if !config.AnycastMacAddressAutoAssigned.IsNull() && config.AnycastMacAddressAutoAssigned.ValueBool() {
		if !config.AnycastMacAddress.IsNull() && !config.AnycastMacAddress.IsUnknown() && config.AnycastMacAddress.ValueString() != "" {
			resp.Diagnostics.AddError(
				"Anycast MAC Address cannot be specified when auto-assigned",
				"The 'anycast_mac_address' field cannot be specified in the configuration when 'anycast_mac_address_auto_assigned_' is set to true. The API will assign this value automatically.",
			)
			return
		}
	}

	// =========================================================================
	// Resource-specific auto-assigned field logic (AnycastMacAddress)
	// =========================================================================
	if !plan.AnycastMacAddressAutoAssigned.IsNull() && plan.AnycastMacAddressAutoAssigned.ValueBool() {
		if !plan.AnycastMacAddressAutoAssigned.Equal(state.AnycastMacAddressAutoAssigned) {
			// anycast_mac_address_auto_assigned_ is changing to true - API will assign value
			resp.Plan.SetAttribute(ctx, path.Root("anycast_mac_address"), types.StringUnknown())
			resp.Diagnostics.AddWarning(
				"Anycast MAC Address will be assigned by the API",
				"The 'anycast_mac_address' field will be automatically assigned by the API because 'anycast_mac_address_auto_assigned_' is being set to true.",
			)
		} else if !plan.AnycastMacAddress.Equal(state.AnycastMacAddress) {
			// User tried to change AnycastMacAddress but it's auto-assigned - suppress diff
			resp.Diagnostics.AddWarning(
				"Ignoring anycast_mac_address changes with auto-assignment enabled",
				"The 'anycast_mac_address' field changes will be ignored because 'anycast_mac_address_auto_assigned_' is set to true.",
			)
			if !state.AnycastMacAddress.IsNull() {
				resp.Plan.SetAttribute(ctx, path.Root("anycast_mac_address"), state.AnycastMacAddress)
			}
		}
	}
}

func filterFabricIndexedEntries(state *verityFabricResourceModel, ref *verityFabricResourceModel) {
	if ref == nil {
		return
	}
	if len(state.ObjectProperties) > 0 && len(ref.ObjectProperties) > 0 {
		state.ObjectProperties[0].SystemGraphs = utils.FilterIndexedEntries(
			state.ObjectProperties[0].SystemGraphs,
			ref.ObjectProperties[0].SystemGraphs,
		)
	}
}
