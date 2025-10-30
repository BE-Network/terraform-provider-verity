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
	_ resource.Resource                = &veritySiteResource{}
	_ resource.ResourceWithConfigure   = &veritySiteResource{}
	_ resource.ResourceWithImportState = &veritySiteResource{}
	_ resource.ResourceWithModifyPlan  = &veritySiteResource{}
)

func NewVeritySiteResource() resource.Resource {
	return &veritySiteResource{}
}

type veritySiteResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *utils.BulkOperationManager
	notifyOperationAdded func()
}

type veritySiteResourceModel struct {
	Name                                      types.String                      `tfsdk:"name"`
	Enable                                    types.Bool                        `tfsdk:"enable"`
	ServiceForSite                            types.String                      `tfsdk:"service_for_site"`
	ServiceForSiteRefType                     types.String                      `tfsdk:"service_for_site_ref_type_"`
	SpanningTreeType                          types.String                      `tfsdk:"spanning_tree_type"`
	RegionName                                types.String                      `tfsdk:"region_name"`
	Revision                                  types.Int64                       `tfsdk:"revision"`
	ForceSpanningTreeOnFabricPorts            types.Bool                        `tfsdk:"force_spanning_tree_on_fabric_ports"`
	ReadOnlyMode                              types.Bool                        `tfsdk:"read_only_mode"`
	DscpToPBitMap                             types.String                      `tfsdk:"dscp_to_p_bit_map"`
	AnycastMacAddress                         types.String                      `tfsdk:"anycast_mac_address"`
	AnycastMacAddressAutoAssigned             types.Bool                        `tfsdk:"anycast_mac_address_auto_assigned_"`
	MacAddressAgingTime                       types.Int64                       `tfsdk:"mac_address_aging_time"`
	MlagDelayRestoreTimer                     types.Int64                       `tfsdk:"mlag_delay_restore_timer"`
	BgpKeepaliveTimer                         types.Int64                       `tfsdk:"bgp_keepalive_timer"`
	BgpHoldDownTimer                          types.Int64                       `tfsdk:"bgp_hold_down_timer"`
	SpineBgpAdvertisementInterval             types.Int64                       `tfsdk:"spine_bgp_advertisement_interval"`
	SpineBgpConnectTimer                      types.Int64                       `tfsdk:"spine_bgp_connect_timer"`
	LeafBgpKeepAliveTimer                     types.Int64                       `tfsdk:"leaf_bgp_keep_alive_timer"`
	LeafBgpHoldDownTimer                      types.Int64                       `tfsdk:"leaf_bgp_hold_down_timer"`
	LeafBgpAdvertisementInterval              types.Int64                       `tfsdk:"leaf_bgp_advertisement_interval"`
	LeafBgpConnectTimer                       types.Int64                       `tfsdk:"leaf_bgp_connect_timer"`
	LinkStateTimeoutValue                     types.Int64                       `tfsdk:"link_state_timeout_value"`
	EvpnMultihomingStartupDelay               types.Int64                       `tfsdk:"evpn_multihoming_startup_delay"`
	EvpnMacHoldtime                           types.Int64                       `tfsdk:"evpn_mac_holdtime"`
	AggressiveReporting                       types.Bool                        `tfsdk:"aggressive_reporting"`
	CrcFailureThreshold                       types.Int64                       `tfsdk:"crc_failure_threshold"`
	EnableDhcpSnooping                        types.Bool                        `tfsdk:"enable_dhcp_snooping"`
	IpSourceGuard                             types.Bool                        `tfsdk:"ip_source_guard"`
	DuplicateAddressDetectionMaxNumberOfMoves types.Int64                       `tfsdk:"duplicate_address_detection_max_number_of_moves"`
	DuplicateAddressDetectionTime             types.Int64                       `tfsdk:"duplicate_address_detection_time"`
	Islands                                   []veritySiteIslandsModel          `tfsdk:"islands"`
	Pairs                                     []veritySitePairsModel            `tfsdk:"pairs"`
	ObjectProperties                          []veritySiteObjectPropertiesModel `tfsdk:"object_properties"`
}

type veritySiteIslandsModel struct {
	ToiSwitchpoint        types.String `tfsdk:"toi_switchpoint"`
	ToiSwitchpointRefType types.String `tfsdk:"toi_switchpoint_ref_type_"`
	Index                 types.Int64  `tfsdk:"index"`
}

func (m veritySiteIslandsModel) GetIndex() types.Int64 {
	return m.Index
}

type veritySitePairsModel struct {
	Name                types.String `tfsdk:"name"`
	Switchpoint1        types.String `tfsdk:"switchpoint_1"`
	Switchpoint1RefType types.String `tfsdk:"switchpoint_1_ref_type_"`
	Switchpoint2        types.String `tfsdk:"switchpoint_2"`
	Switchpoint2RefType types.String `tfsdk:"switchpoint_2_ref_type_"`
	LagGroup            types.String `tfsdk:"lag_group"`
	LagGroupRefType     types.String `tfsdk:"lag_group_ref_type_"`
	IsWhiteboxPair      types.Bool   `tfsdk:"is_whitebox_pair"`
	Index               types.Int64  `tfsdk:"index"`
}

func (m veritySitePairsModel) GetIndex() types.Int64 {
	return m.Index
}

type veritySiteObjectPropertiesModel struct {
	SystemGraphs []veritySiteSystemGraphsModel `tfsdk:"system_graphs"`
}

type veritySiteSystemGraphsModel struct {
	GraphNumData types.String `tfsdk:"graph_num_data"`
	Index        types.Int64  `tfsdk:"index"`
}

func (r *veritySiteResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_site"
}

func (r *veritySiteResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *veritySiteResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity Site",
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
			"service_for_site": schema.StringAttribute{
				Description: "Service for Site",
				Optional:    true,
			},
			"service_for_site_ref_type_": schema.StringAttribute{
				Description: "Object type for service_for_site field",
				Optional:    true,
			},
			"spanning_tree_type": schema.StringAttribute{
				Description: "Sets the spanning tree type for all Ports in this Site with Spanning Tree enabled",
				Optional:    true,
			},
			"region_name": schema.StringAttribute{
				Description: "Defines the logical boundary of the network. All switches in an MSTP region must have the same configured region name",
				Optional:    true,
			},
			"revision": schema.Int64Attribute{
				Description: "A logical number that signifies a revision for the MSTP configuration. All switches in an MSTP region must have the same revision number (maximum: 65535)",
				Optional:    true,
			},
			"force_spanning_tree_on_fabric_ports": schema.BoolAttribute{
				Description: "Enable spanning tree on all fabric connections. This overrides the Eth Port Settings for Fabric ports",
				Optional:    true,
			},
			"read_only_mode": schema.BoolAttribute{
				Description: "When Read Only Mode is checked, vNetC will perform all functions except writing database updates to the target hardware",
				Optional:    true,
			},
			"dscp_to_p_bit_map": schema.StringAttribute{
				Description: "For any Service that is using DSCP to TC map packet prioritization. A string of length 64 with a 0-7 in each position (maxLength: 64)",
				Optional:    true,
			},
			"anycast_mac_address": schema.StringAttribute{
				Description: "Anycast MAC address to use. This field should not be specified when 'anycast_mac_address_auto_assigned_' is set to true, as the API will assign this value automatically. Used for MAC VRRP.",
				Optional:    true,
				Computed:    true,
			},
			"anycast_mac_address_auto_assigned_": schema.BoolAttribute{
				Description: "Whether the anycast MAC address should be automatically assigned by the API. When set to true, do not specify the 'anycast_mac_address' field in your configuration.",
				Optional:    true,
			},
			"mac_address_aging_time": schema.Int64Attribute{
				Description: "MAC Address Aging Time (minimum: 1, maximum: 100000)",
				Optional:    true,
			},
			"mlag_delay_restore_timer": schema.Int64Attribute{
				Description: "MLAG Delay Restore Timer (minimum: 1, maximum: 3600)",
				Optional:    true,
			},
			"bgp_keepalive_timer": schema.Int64Attribute{
				Description: "Spine BGP Keepalive Timer (minimum: 1, maximum: 3600)",
				Optional:    true,
			},
			"bgp_hold_down_timer": schema.Int64Attribute{
				Description: "Spine BGP Hold Down Timer (minimum: 1, maximum: 3600)",
				Optional:    true,
			},
			"spine_bgp_advertisement_interval": schema.Int64Attribute{
				Description: "BGP Advertisement Interval for spines/superspines. Use \"0\" for immediate updates (maximum: 3600)",
				Optional:    true,
			},
			"spine_bgp_connect_timer": schema.Int64Attribute{
				Description: "BGP Connect Timer (minimum: 1, maximum: 3600)",
				Optional:    true,
			},
			"leaf_bgp_keep_alive_timer": schema.Int64Attribute{
				Description: "Leaf BGP Keep Alive Timer (minimum: 1, maximum: 3600)",
				Optional:    true,
			},
			"leaf_bgp_hold_down_timer": schema.Int64Attribute{
				Description: "Leaf BGP Hold Down Timer (minimum: 1, maximum: 3600)",
				Optional:    true,
			},
			"leaf_bgp_advertisement_interval": schema.Int64Attribute{
				Description: "BGP Advertisement Interval for leafs. Use \"0\" for immediate updates (maximum: 3600)",
				Optional:    true,
			},
			"leaf_bgp_connect_timer": schema.Int64Attribute{
				Description: "BGP Connect Timer (minimum: 1, maximum: 3600)",
				Optional:    true,
			},
			"link_state_timeout_value": schema.Int64Attribute{
				Description: "Link State Timeout Value",
				Optional:    true,
			},
			"evpn_multihoming_startup_delay": schema.Int64Attribute{
				Description: "Startup Delay",
				Optional:    true,
			},
			"evpn_mac_holdtime": schema.Int64Attribute{
				Description: "MAC Holdtime",
				Optional:    true,
			},
			"aggressive_reporting": schema.BoolAttribute{
				Description: "Fast Reporting of Switch Communications, Link Up/Down, and BGP Status",
				Optional:    true,
			},
			"crc_failure_threshold": schema.Int64Attribute{
				Description: "Threshold in Errors per second that when met will disable the links as part of LAGs (minimum: 1, maximum: 4294967296)",
				Optional:    true,
			},
			"enable_dhcp_snooping": schema.BoolAttribute{
				Description: "Enables the switches to monitor DHCP traffic and collect assigned IP addresses which are then placed in the DHCP assigned IPs report.",
				Optional:    true,
			},
			"ip_source_guard": schema.BoolAttribute{
				Description: "On untrusted ports, only allow known traffic from known IP addresses. IP addresses are discovered via DHCP snooping or with static IP settings",
				Optional:    true,
			},
			"duplicate_address_detection_max_number_of_moves": schema.Int64Attribute{
				Description: "Duplicate Address Detection Max Number of Moves",
				Optional:    true,
			},
			"duplicate_address_detection_time": schema.Int64Attribute{
				Description: "Duplicate Address Detection Time",
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"islands": schema.ListNestedBlock{
				Description: "List of islands",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"toi_switchpoint": schema.StringAttribute{
							Description: "TOI Switchpoint",
							Optional:    true,
						},
						"toi_switchpoint_ref_type_": schema.StringAttribute{
							Description: "Object type for toi_switchpoint field",
							Optional:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object. Zero if you want to add an object to the list.",
							Optional:    true,
						},
					},
				},
			},
			"pairs": schema.ListNestedBlock{
				Description: "List of pairs",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							Description: "Object Name. Must be unique.",
							Optional:    true,
						},
						"switchpoint_1": schema.StringAttribute{
							Description: "Switchpoint",
							Optional:    true,
						},
						"switchpoint_1_ref_type_": schema.StringAttribute{
							Description: "Object type for switchpoint_1 field",
							Optional:    true,
						},
						"switchpoint_2": schema.StringAttribute{
							Description: "Switchpoint",
							Optional:    true,
						},
						"switchpoint_2_ref_type_": schema.StringAttribute{
							Description: "Object type for switchpoint_2 field",
							Optional:    true,
						},
						"lag_group": schema.StringAttribute{
							Description: "LAG Group",
							Optional:    true,
						},
						"lag_group_ref_type_": schema.StringAttribute{
							Description: "Object type for lag_group field",
							Optional:    true,
						},
						"is_whitebox_pair": schema.BoolAttribute{
							Description: "LAG Pair",
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
				Description: "Object properties for the Site",
				NestedObject: schema.NestedBlockObject{
					Blocks: map[string]schema.Block{
						"system_graphs": schema.ListNestedBlock{
							Description: "System graphs",
							NestedObject: schema.NestedBlockObject{
								Attributes: map[string]schema.Attribute{
									"graph_num_data": schema.StringAttribute{
										Description: "The graph data detailing this graph choice",
										Optional:    true,
									},
									"index": schema.Int64Attribute{
										Description: "The index identifying the object. Zero if you want to add an object to the list.",
										Optional:    true,
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

func (r *veritySiteResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	resp.Diagnostics.AddError(
		"Create Not Supported",
		"Site resources cannot be created. They represent existing site configurations that can only be read and updated.",
	)
}

func (r *veritySiteResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state veritySiteResourceModel
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

	siteName := state.Name.ValueString()

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if siteData, exists := r.bulkOpsMgr.GetResourceResponse("site", siteName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached site data for %s from recent operation", siteName))
			state = populateSiteState(ctx, state, siteData, nil)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("site") {
		tflog.Info(ctx, fmt.Sprintf("Skipping site %s verification – trusting recent successful API operation", siteName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching sites for verification of %s", siteName))

	type SitesResponse struct {
		Site map[string]interface{} `json:"site"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "sites", siteName,
		func() (SitesResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch sites")
			respAPI, err := r.client.SitesAPI.SitesGet(ctx).Execute()
			if err != nil {
				return SitesResponse{}, fmt.Errorf("error reading sites: %v", err)
			}
			defer respAPI.Body.Close()

			var res SitesResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return SitesResponse{}, fmt.Errorf("failed to decode sites response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched sites data with %d sites", len(res.Site)))
			return res, nil
		}, getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Site %s", siteName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for site with name: %s", siteName))

	siteData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.Site,
		siteName,
		func(data interface{}) (string, bool) {
			if site, ok := data.(map[string]interface{}); ok {
				if name, ok := site["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Site with name '%s' not found in API response", siteName))
		resp.State.RemoveResource(ctx)
		return
	}

	siteMap, ok := siteData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Site Data",
			fmt.Sprintf("Site data is not in expected format for %s", siteName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found site '%s' under API key '%s'", siteName, actualAPIName))

	state = populateSiteState(ctx, state, siteMap, nil)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *veritySiteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state veritySiteResourceModel

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
	siteReq := openapi.SitesPatchRequestSiteValue{}
	hasChanges := false

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { siteReq.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.SpanningTreeType, state.SpanningTreeType, func(v *string) { siteReq.SpanningTreeType = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.RegionName, state.RegionName, func(v *string) { siteReq.RegionName = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.DscpToPBitMap, state.DscpToPBitMap, func(v *string) { siteReq.DscpToPBitMap = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { siteReq.Enable = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.ForceSpanningTreeOnFabricPorts, state.ForceSpanningTreeOnFabricPorts, func(v *bool) { siteReq.ForceSpanningTreeOnFabricPorts = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.ReadOnlyMode, state.ReadOnlyMode, func(v *bool) { siteReq.ReadOnlyMode = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.AggressiveReporting, state.AggressiveReporting, func(v *bool) { siteReq.AggressiveReporting = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.EnableDhcpSnooping, state.EnableDhcpSnooping, func(v *bool) { siteReq.EnableDhcpSnooping = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.IpSourceGuard, state.IpSourceGuard, func(v *bool) { siteReq.IpSourceGuard = v }, &hasChanges)

	// Handle nullable int64 field changes
	utils.CompareAndSetNullableInt64Field(plan.MacAddressAgingTime, state.MacAddressAgingTime, func(v *openapi.NullableInt32) { siteReq.MacAddressAgingTime = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.MlagDelayRestoreTimer, state.MlagDelayRestoreTimer, func(v *openapi.NullableInt32) { siteReq.MlagDelayRestoreTimer = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.BgpKeepaliveTimer, state.BgpKeepaliveTimer, func(v *openapi.NullableInt32) { siteReq.BgpKeepaliveTimer = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.BgpHoldDownTimer, state.BgpHoldDownTimer, func(v *openapi.NullableInt32) { siteReq.BgpHoldDownTimer = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.SpineBgpAdvertisementInterval, state.SpineBgpAdvertisementInterval, func(v *openapi.NullableInt32) { siteReq.SpineBgpAdvertisementInterval = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.SpineBgpConnectTimer, state.SpineBgpConnectTimer, func(v *openapi.NullableInt32) { siteReq.SpineBgpConnectTimer = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.LeafBgpKeepAliveTimer, state.LeafBgpKeepAliveTimer, func(v *openapi.NullableInt32) { siteReq.LeafBgpKeepAliveTimer = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.LeafBgpHoldDownTimer, state.LeafBgpHoldDownTimer, func(v *openapi.NullableInt32) { siteReq.LeafBgpHoldDownTimer = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.LeafBgpAdvertisementInterval, state.LeafBgpAdvertisementInterval, func(v *openapi.NullableInt32) { siteReq.LeafBgpAdvertisementInterval = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.LeafBgpConnectTimer, state.LeafBgpConnectTimer, func(v *openapi.NullableInt32) { siteReq.LeafBgpConnectTimer = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.Revision, state.Revision, func(v *openapi.NullableInt32) { siteReq.Revision = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.LinkStateTimeoutValue, state.LinkStateTimeoutValue, func(v *openapi.NullableInt32) { siteReq.LinkStateTimeoutValue = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.EvpnMultihomingStartupDelay, state.EvpnMultihomingStartupDelay, func(v *openapi.NullableInt32) { siteReq.EvpnMultihomingStartupDelay = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.EvpnMacHoldtime, state.EvpnMacHoldtime, func(v *openapi.NullableInt32) { siteReq.EvpnMacHoldtime = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.CrcFailureThreshold, state.CrcFailureThreshold, func(v *openapi.NullableInt32) { siteReq.CrcFailureThreshold = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.DuplicateAddressDetectionMaxNumberOfMoves, state.DuplicateAddressDetectionMaxNumberOfMoves, func(v *openapi.NullableInt32) { siteReq.DuplicateAddressDetectionMaxNumberOfMoves = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.DuplicateAddressDetectionTime, state.DuplicateAddressDetectionTime, func(v *openapi.NullableInt32) { siteReq.DuplicateAddressDetectionTime = *v }, &hasChanges)

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		if len(state.ObjectProperties) == 0 ||
			len(plan.ObjectProperties[0].SystemGraphs) != len(state.ObjectProperties[0].SystemGraphs) {
			siteObjProps := openapi.SitesPatchRequestSiteValueObjectProperties{}
			if len(plan.ObjectProperties[0].SystemGraphs) > 0 {
				systemGraphsList := make([]openapi.SitesPatchRequestSiteValueObjectPropertiesSystemGraphsInner, len(plan.ObjectProperties[0].SystemGraphs))
				for i, graph := range plan.ObjectProperties[0].SystemGraphs {
					graphProps := openapi.SitesPatchRequestSiteValueObjectPropertiesSystemGraphsInner{}
					if !graph.GraphNumData.IsNull() {
						graphProps.GraphNumData = openapi.PtrString(graph.GraphNumData.ValueString())
					} else {
						graphProps.GraphNumData = nil
					}
					if !graph.Index.IsNull() {
						graphProps.Index = openapi.PtrInt32(int32(graph.Index.ValueInt64()))
					} else {
						graphProps.Index = nil
					}
					systemGraphsList[i] = graphProps
				}
				siteObjProps.SystemGraphs = systemGraphsList
			}
			siteReq.ObjectProperties = &siteObjProps
			hasChanges = true
		} else {
			// Check if any individual graph changed
			graphsChanged := false
			for i, planGraph := range plan.ObjectProperties[0].SystemGraphs {
				if i >= len(state.ObjectProperties[0].SystemGraphs) ||
					!planGraph.GraphNumData.Equal(state.ObjectProperties[0].SystemGraphs[i].GraphNumData) ||
					!planGraph.Index.Equal(state.ObjectProperties[0].SystemGraphs[i].Index) {
					graphsChanged = true
					break
				}
			}
			if graphsChanged {
				siteObjProps := openapi.SitesPatchRequestSiteValueObjectProperties{}
				if len(plan.ObjectProperties[0].SystemGraphs) > 0 {
					systemGraphsList := make([]openapi.SitesPatchRequestSiteValueObjectPropertiesSystemGraphsInner, len(plan.ObjectProperties[0].SystemGraphs))
					for i, graph := range plan.ObjectProperties[0].SystemGraphs {
						graphProps := openapi.SitesPatchRequestSiteValueObjectPropertiesSystemGraphsInner{}
						if !graph.GraphNumData.IsNull() {
							graphProps.GraphNumData = openapi.PtrString(graph.GraphNumData.ValueString())
						} else {
							graphProps.GraphNumData = nil
						}
						if !graph.Index.IsNull() {
							graphProps.Index = openapi.PtrInt32(int32(graph.Index.ValueInt64()))
						} else {
							graphProps.Index = nil
						}
						systemGraphsList[i] = graphProps
					}
					siteObjProps.SystemGraphs = systemGraphsList
				}
				siteReq.ObjectProperties = &siteObjProps
				hasChanges = true
			}
		}
	}

	// Handle service_for_site and service_for_site_ref_type_ fields using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.ServiceForSite, state.ServiceForSite, plan.ServiceForSiteRefType, state.ServiceForSiteRefType,
		func(v *string) { siteReq.ServiceForSite = v },
		func(v *string) { siteReq.ServiceForSiteRefType = v },
		"service_for_site", "service_for_site_ref_type_",
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
				siteReq.AnycastMacAddress = openapi.PtrString(plan.AnycastMacAddress.ValueString())
			} else {
				siteReq.AnycastMacAddress = openapi.PtrString("")
			}
		}

		// Handle AnycastMacAddressAutoAssigned field changes
		if anycastMacAddressAutoAssignedChanged {
			// Only send anycast_mac_address_auto_assigned_ if the user has explicitly specified it in their configuration
			var config veritySiteResourceModel
			userSpecifiedAnycastMacAddressAutoAssigned := false
			if !req.Config.Raw.IsNull() {
				if err := req.Config.Get(ctx, &config); err == nil {
					userSpecifiedAnycastMacAddressAutoAssigned = !config.AnycastMacAddressAutoAssigned.IsNull()
				}
			}

			if userSpecifiedAnycastMacAddressAutoAssigned {
				siteReq.AnycastMacAddressAutoAssigned = openapi.PtrBool(plan.AnycastMacAddressAutoAssigned.ValueBool())

				// Special case: When changing from auto-assigned (true) to manual (false),
				// the API requires both anycast_mac_address_auto_assigned_ and anycast_mac_address fields to be sent.
				if !state.AnycastMacAddressAutoAssigned.IsNull() && state.AnycastMacAddressAutoAssigned.ValueBool() &&
					!plan.AnycastMacAddressAutoAssigned.ValueBool() {
					// Changing from auto-assigned=true to auto-assigned=false
					// Must include AnycastMacAddress value in the request for the change to take effect
					if !plan.AnycastMacAddress.IsNull() && plan.AnycastMacAddress.ValueString() != "" {
						siteReq.AnycastMacAddress = openapi.PtrString(plan.AnycastMacAddress.ValueString())
					} else if !state.AnycastMacAddress.IsNull() && state.AnycastMacAddress.ValueString() != "" {
						// Use current state AnycastMacAddress if plan doesn't specify one
						siteReq.AnycastMacAddress = openapi.PtrString(state.AnycastMacAddress.ValueString())
					}
				}
			}
		} else if anycastMacAddressChanged {
			// AnycastMacAddress changed but AnycastMacAddressAutoAssigned didn't change
			// Send the auto-assigned flag to maintain consistency with API
			if !plan.AnycastMacAddressAutoAssigned.IsNull() {
				siteReq.AnycastMacAddressAutoAssigned = openapi.PtrBool(plan.AnycastMacAddressAutoAssigned.ValueBool())
			} else if !state.AnycastMacAddressAutoAssigned.IsNull() {
				siteReq.AnycastMacAddressAutoAssigned = openapi.PtrBool(state.AnycastMacAddressAutoAssigned.ValueBool())
			} else {
				siteReq.AnycastMacAddressAutoAssigned = openapi.PtrBool(false)
			}
		}

		hasChanges = true
	}

	// Handle islands
	changedIslands, islandsChanged := utils.ProcessIndexedArrayUpdates(plan.Islands, state.Islands,
		utils.IndexedItemHandler[veritySiteIslandsModel, openapi.SitesPatchRequestSiteValueIslandsInner]{
			CreateNew: func(planItem veritySiteIslandsModel) openapi.SitesPatchRequestSiteValueIslandsInner {
				newIsland := openapi.SitesPatchRequestSiteValueIslandsInner{}

				// Handle string fields
				utils.SetStringFields([]utils.StringFieldMapping{
					{FieldName: "ToiSwitchpoint", APIField: &newIsland.ToiSwitchpoint, TFValue: planItem.ToiSwitchpoint},
					{FieldName: "ToiSwitchpointRefType", APIField: &newIsland.ToiSwitchpointRefType, TFValue: planItem.ToiSwitchpointRefType},
				})

				// Handle int64 fields
				utils.SetInt64Fields([]utils.Int64FieldMapping{
					{FieldName: "Index", APIField: &newIsland.Index, TFValue: planItem.Index},
				})

				return newIsland
			},
			UpdateExisting: func(planItem veritySiteIslandsModel, stateItem veritySiteIslandsModel) (openapi.SitesPatchRequestSiteValueIslandsInner, bool) {
				updateIsland := openapi.SitesPatchRequestSiteValueIslandsInner{}
				fieldChanged := false

				// Handle toi_switchpoint and toi_switchpoint_ref_type_ using one ref type supported pattern
				if !utils.HandleOneRefTypeSupported(
					planItem.ToiSwitchpoint, stateItem.ToiSwitchpoint, planItem.ToiSwitchpointRefType, stateItem.ToiSwitchpointRefType,
					func(v *string) { updateIsland.ToiSwitchpoint = v },
					func(v *string) { updateIsland.ToiSwitchpointRefType = v },
					"toi_switchpoint", "toi_switchpoint_ref_type_",
					&fieldChanged, &resp.Diagnostics,
				) {
					return updateIsland, false
				}

				// Handle index field change
				utils.CompareAndSetInt64Field(planItem.Index, stateItem.Index, func(v *int32) { updateIsland.Index = v }, &fieldChanged)

				return updateIsland, fieldChanged
			},
			CreateDeleted: func(index int64) openapi.SitesPatchRequestSiteValueIslandsInner {
				return openapi.SitesPatchRequestSiteValueIslandsInner{
					Index: openapi.PtrInt32(int32(index)),
				}
			},
		})
	if islandsChanged {
		siteReq.Islands = changedIslands
		hasChanges = true
	}

	// Handle pairs list
	changedPairs, pairsChanged := utils.ProcessIndexedArrayUpdates(plan.Pairs, state.Pairs,
		utils.IndexedItemHandler[veritySitePairsModel, openapi.SitesPatchRequestSiteValuePairsInner]{
			CreateNew: func(planItem veritySitePairsModel) openapi.SitesPatchRequestSiteValuePairsInner {
				newPair := openapi.SitesPatchRequestSiteValuePairsInner{}

				// Handle boolean fields
				utils.SetBoolFields([]utils.BoolFieldMapping{
					{FieldName: "IsWhiteboxPair", APIField: &newPair.IsWhiteboxPair, TFValue: planItem.IsWhiteboxPair},
				})

				// Handle string fields
				utils.SetStringFields([]utils.StringFieldMapping{
					{FieldName: "Name", APIField: &newPair.Name, TFValue: planItem.Name},
					{FieldName: "Switchpoint1", APIField: &newPair.Switchpoint1, TFValue: planItem.Switchpoint1},
					{FieldName: "Switchpoint1RefType", APIField: &newPair.Switchpoint1RefType, TFValue: planItem.Switchpoint1RefType},
					{FieldName: "Switchpoint2", APIField: &newPair.Switchpoint2, TFValue: planItem.Switchpoint2},
					{FieldName: "Switchpoint2RefType", APIField: &newPair.Switchpoint2RefType, TFValue: planItem.Switchpoint2RefType},
					{FieldName: "LagGroup", APIField: &newPair.LagGroup, TFValue: planItem.LagGroup},
					{FieldName: "LagGroupRefType", APIField: &newPair.LagGroupRefType, TFValue: planItem.LagGroupRefType},
				})

				// Handle int64 fields
				utils.SetInt64Fields([]utils.Int64FieldMapping{
					{FieldName: "Index", APIField: &newPair.Index, TFValue: planItem.Index},
				})

				return newPair
			},
			UpdateExisting: func(planItem veritySitePairsModel, stateItem veritySitePairsModel) (openapi.SitesPatchRequestSiteValuePairsInner, bool) {
				updatePair := openapi.SitesPatchRequestSiteValuePairsInner{}
				fieldChanged := false

				// Handle boolean field changes
				utils.CompareAndSetBoolField(planItem.IsWhiteboxPair, stateItem.IsWhiteboxPair, func(v *bool) { updatePair.IsWhiteboxPair = v }, &fieldChanged)

				// Handle simple string field changes
				utils.CompareAndSetStringField(planItem.Name, stateItem.Name, func(v *string) { updatePair.Name = v }, &fieldChanged)

				// Handle switchpoint_1 and switchpoint_1_ref_type_ using one ref type supported pattern
				if !utils.HandleOneRefTypeSupported(
					planItem.Switchpoint1, stateItem.Switchpoint1, planItem.Switchpoint1RefType, stateItem.Switchpoint1RefType,
					func(v *string) { updatePair.Switchpoint1 = v },
					func(v *string) { updatePair.Switchpoint1RefType = v },
					"switchpoint_1", "switchpoint_1_ref_type_",
					&fieldChanged, &resp.Diagnostics,
				) {
					return updatePair, false
				}

				// Handle switchpoint_2 and switchpoint_2_ref_type_ using one ref type supported pattern
				if !utils.HandleOneRefTypeSupported(
					planItem.Switchpoint2, stateItem.Switchpoint2, planItem.Switchpoint2RefType, stateItem.Switchpoint2RefType,
					func(v *string) { updatePair.Switchpoint2 = v },
					func(v *string) { updatePair.Switchpoint2RefType = v },
					"switchpoint_2", "switchpoint_2_ref_type_",
					&fieldChanged, &resp.Diagnostics,
				) {
					return updatePair, false
				}

				// Handle lag_group and lag_group_ref_type_ using one ref type supported pattern
				if !utils.HandleOneRefTypeSupported(
					planItem.LagGroup, stateItem.LagGroup, planItem.LagGroupRefType, stateItem.LagGroupRefType,
					func(v *string) { updatePair.LagGroup = v },
					func(v *string) { updatePair.LagGroupRefType = v },
					"lag_group", "lag_group_ref_type_",
					&fieldChanged, &resp.Diagnostics,
				) {
					return updatePair, false
				}

				// Handle index field change
				utils.CompareAndSetInt64Field(planItem.Index, stateItem.Index, func(v *int32) { updatePair.Index = v }, &fieldChanged)

				return updatePair, fieldChanged
			},
			CreateDeleted: func(index int64) openapi.SitesPatchRequestSiteValuePairsInner {
				return openapi.SitesPatchRequestSiteValuePairsInner{
					Index: openapi.PtrInt32(int32(index)),
				}
			},
		})
	if pairsChanged {
		siteReq.Pairs = changedPairs
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "site", name, siteReq, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Site %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "sites")

	var minState veritySiteResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if siteData, exists := bulkMgr.GetResourceResponse("site", name); exists {
			// Use the cached data from the API response with plan values as fallback
			state := populateSiteState(ctx, minState, siteData, &plan)
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

func (r *veritySiteResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddError(
		"Delete Not Supported",
		"Site resources cannot be deleted. They represent existing site configurations that can only be read and updated.",
	)
}

func (r *veritySiteResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func populateSiteState(ctx context.Context, state veritySiteResourceModel, siteData map[string]interface{}, plan *veritySiteResourceModel) veritySiteResourceModel {
	state.Name = types.StringValue(fmt.Sprintf("%v", siteData["name"]))

	// For each field, check if it's in the API response first,
	// if not: use plan value (if plan provided and not null), otherwise preserve current state value

	if val, ok := siteData["enable"].(bool); ok {
		state.Enable = types.BoolValue(val)
	} else if plan != nil && !plan.Enable.IsNull() {
		state.Enable = plan.Enable
	} else {
		state.Enable = types.BoolNull()
	}

	if val, ok := siteData["service_for_site"].(string); ok {
		state.ServiceForSite = types.StringValue(val)
	} else if plan != nil && !plan.ServiceForSite.IsNull() {
		state.ServiceForSite = plan.ServiceForSite
	} else {
		state.ServiceForSite = types.StringNull()
	}

	if val, ok := siteData["service_for_site_ref_type_"].(string); ok {
		state.ServiceForSiteRefType = types.StringValue(val)
	} else if plan != nil && !plan.ServiceForSiteRefType.IsNull() {
		state.ServiceForSiteRefType = plan.ServiceForSiteRefType
	} else {
		state.ServiceForSiteRefType = types.StringNull()
	}

	if val, ok := siteData["spanning_tree_type"].(string); ok {
		state.SpanningTreeType = types.StringValue(val)
	} else if plan != nil && !plan.SpanningTreeType.IsNull() {
		state.SpanningTreeType = plan.SpanningTreeType
	} else {
		state.SpanningTreeType = types.StringNull()
	}

	if val, ok := siteData["region_name"].(string); ok {
		state.RegionName = types.StringValue(val)
	} else if plan != nil && !plan.RegionName.IsNull() {
		state.RegionName = plan.RegionName
	} else {
		state.RegionName = types.StringNull()
	}

	if val, ok := siteData["revision"]; ok {
		switch v := val.(type) {
		case float64:
			state.Revision = types.Int64Value(int64(v))
		case int:
			state.Revision = types.Int64Value(int64(v))
		case int32:
			state.Revision = types.Int64Value(int64(v))
		case nil:
			state.Revision = types.Int64Null()
		default:
			if plan != nil && !plan.Revision.IsNull() {
				state.Revision = plan.Revision
			} else {
				state.Revision = types.Int64Null()
			}
		}
	} else if plan != nil && !plan.Revision.IsNull() {
		state.Revision = plan.Revision
	} else {
		state.Revision = types.Int64Null()
	}

	if val, ok := siteData["force_spanning_tree_on_fabric_ports"].(bool); ok {
		state.ForceSpanningTreeOnFabricPorts = types.BoolValue(val)
	} else if plan != nil && !plan.ForceSpanningTreeOnFabricPorts.IsNull() {
		state.ForceSpanningTreeOnFabricPorts = plan.ForceSpanningTreeOnFabricPorts
	} else {
		state.ForceSpanningTreeOnFabricPorts = types.BoolNull()
	}

	if val, ok := siteData["read_only_mode"].(bool); ok {
		state.ReadOnlyMode = types.BoolValue(val)
	} else if plan != nil && !plan.ReadOnlyMode.IsNull() {
		state.ReadOnlyMode = plan.ReadOnlyMode
	} else {
		state.ReadOnlyMode = types.BoolNull()
	}

	if val, ok := siteData["dscp_to_p_bit_map"].(string); ok {
		state.DscpToPBitMap = types.StringValue(val)
	} else if plan != nil && !plan.DscpToPBitMap.IsNull() {
		state.DscpToPBitMap = plan.DscpToPBitMap
	} else {
		state.DscpToPBitMap = types.StringNull()
	}

	if val, ok := siteData["anycast_mac_address"].(string); ok {
		state.AnycastMacAddress = types.StringValue(val)
	} else if plan != nil && !plan.AnycastMacAddress.IsNull() && !plan.AnycastMacAddress.IsUnknown() {
		state.AnycastMacAddress = plan.AnycastMacAddress
	} else {
		state.AnycastMacAddress = types.StringNull()
	}

	if val, ok := siteData["anycast_mac_address_auto_assigned_"].(bool); ok {
		state.AnycastMacAddressAutoAssigned = types.BoolValue(val)
	} else if plan != nil && !plan.AnycastMacAddressAutoAssigned.IsNull() {
		state.AnycastMacAddressAutoAssigned = plan.AnycastMacAddressAutoAssigned
	} else {
		state.AnycastMacAddressAutoAssigned = types.BoolNull()
	}

	// Handle all the integer fields
	if val, ok := siteData["mac_address_aging_time"]; ok {
		switch v := val.(type) {
		case float64:
			state.MacAddressAgingTime = types.Int64Value(int64(v))
		case int:
			state.MacAddressAgingTime = types.Int64Value(int64(v))
		case int32:
			state.MacAddressAgingTime = types.Int64Value(int64(v))
		case nil:
			state.MacAddressAgingTime = types.Int64Null()
		default:
			if plan != nil && !plan.MacAddressAgingTime.IsNull() {
				state.MacAddressAgingTime = plan.MacAddressAgingTime
			} else {
				state.MacAddressAgingTime = types.Int64Null()
			}
		}
	} else if plan != nil && !plan.MacAddressAgingTime.IsNull() {
		state.MacAddressAgingTime = plan.MacAddressAgingTime
	} else {
		state.MacAddressAgingTime = types.Int64Null()
	}

	if val, ok := siteData["mlag_delay_restore_timer"]; ok {
		switch v := val.(type) {
		case float64:
			state.MlagDelayRestoreTimer = types.Int64Value(int64(v))
		case int:
			state.MlagDelayRestoreTimer = types.Int64Value(int64(v))
		case int32:
			state.MlagDelayRestoreTimer = types.Int64Value(int64(v))
		case nil:
			state.MlagDelayRestoreTimer = types.Int64Null()
		default:
			if plan != nil && !plan.MlagDelayRestoreTimer.IsNull() {
				state.MlagDelayRestoreTimer = plan.MlagDelayRestoreTimer
			} else {
				state.MlagDelayRestoreTimer = types.Int64Null()
			}
		}
	} else if plan != nil && !plan.MlagDelayRestoreTimer.IsNull() {
		state.MlagDelayRestoreTimer = plan.MlagDelayRestoreTimer
	} else {
		state.MlagDelayRestoreTimer = types.Int64Null()
	}

	if val, ok := siteData["bgp_keepalive_timer"]; ok {
		switch v := val.(type) {
		case float64:
			state.BgpKeepaliveTimer = types.Int64Value(int64(v))
		case int:
			state.BgpKeepaliveTimer = types.Int64Value(int64(v))
		case int32:
			state.BgpKeepaliveTimer = types.Int64Value(int64(v))
		case nil:
			state.BgpKeepaliveTimer = types.Int64Null()
		default:
			if plan != nil && !plan.BgpKeepaliveTimer.IsNull() {
				state.BgpKeepaliveTimer = plan.BgpKeepaliveTimer
			} else {
				state.BgpKeepaliveTimer = types.Int64Null()
			}
		}
	} else if plan != nil && !plan.BgpKeepaliveTimer.IsNull() {
		state.BgpKeepaliveTimer = plan.BgpKeepaliveTimer
	} else {
		state.BgpKeepaliveTimer = types.Int64Null()
	}

	if val, ok := siteData["bgp_hold_down_timer"]; ok {
		switch v := val.(type) {
		case float64:
			state.BgpHoldDownTimer = types.Int64Value(int64(v))
		case int:
			state.BgpHoldDownTimer = types.Int64Value(int64(v))
		case int32:
			state.BgpHoldDownTimer = types.Int64Value(int64(v))
		case nil:
			state.BgpHoldDownTimer = types.Int64Null()
		default:
			if plan != nil && !plan.BgpHoldDownTimer.IsNull() {
				state.BgpHoldDownTimer = plan.BgpHoldDownTimer
			} else {
				state.BgpHoldDownTimer = types.Int64Null()
			}
		}
	} else if plan != nil && !plan.BgpHoldDownTimer.IsNull() {
		state.BgpHoldDownTimer = plan.BgpHoldDownTimer
	} else {
		state.BgpHoldDownTimer = types.Int64Null()
	}

	if val, ok := siteData["spine_bgp_advertisement_interval"]; ok {
		switch v := val.(type) {
		case float64:
			state.SpineBgpAdvertisementInterval = types.Int64Value(int64(v))
		case int:
			state.SpineBgpAdvertisementInterval = types.Int64Value(int64(v))
		case int32:
			state.SpineBgpAdvertisementInterval = types.Int64Value(int64(v))
		case nil:
			state.SpineBgpAdvertisementInterval = types.Int64Null()
		default:
			if plan != nil && !plan.SpineBgpAdvertisementInterval.IsNull() {
				state.SpineBgpAdvertisementInterval = plan.SpineBgpAdvertisementInterval
			} else {
				state.SpineBgpAdvertisementInterval = types.Int64Null()
			}
		}
	} else if plan != nil && !plan.SpineBgpAdvertisementInterval.IsNull() {
		state.SpineBgpAdvertisementInterval = plan.SpineBgpAdvertisementInterval
	} else {
		state.SpineBgpAdvertisementInterval = types.Int64Null()
	}

	if val, ok := siteData["spine_bgp_connect_timer"]; ok {
		switch v := val.(type) {
		case float64:
			state.SpineBgpConnectTimer = types.Int64Value(int64(v))
		case int:
			state.SpineBgpConnectTimer = types.Int64Value(int64(v))
		case int32:
			state.SpineBgpConnectTimer = types.Int64Value(int64(v))
		case nil:
			state.SpineBgpConnectTimer = types.Int64Null()
		default:
			if plan != nil && !plan.SpineBgpConnectTimer.IsNull() {
				state.SpineBgpConnectTimer = plan.SpineBgpConnectTimer
			} else {
				state.SpineBgpConnectTimer = types.Int64Null()
			}
		}
	} else if plan != nil && !plan.SpineBgpConnectTimer.IsNull() {
		state.SpineBgpConnectTimer = plan.SpineBgpConnectTimer
	} else {
		state.SpineBgpConnectTimer = types.Int64Null()
	}

	if val, ok := siteData["leaf_bgp_keep_alive_timer"]; ok {
		switch v := val.(type) {
		case float64:
			state.LeafBgpKeepAliveTimer = types.Int64Value(int64(v))
		case int:
			state.LeafBgpKeepAliveTimer = types.Int64Value(int64(v))
		case int32:
			state.LeafBgpKeepAliveTimer = types.Int64Value(int64(v))
		case nil:
			state.LeafBgpKeepAliveTimer = types.Int64Null()
		default:
			if plan != nil && !plan.LeafBgpKeepAliveTimer.IsNull() {
				state.LeafBgpKeepAliveTimer = plan.LeafBgpKeepAliveTimer
			} else {
				state.LeafBgpKeepAliveTimer = types.Int64Null()
			}
		}
	} else if plan != nil && !plan.LeafBgpKeepAliveTimer.IsNull() {
		state.LeafBgpKeepAliveTimer = plan.LeafBgpKeepAliveTimer
	} else {
		state.LeafBgpKeepAliveTimer = types.Int64Null()
	}

	if val, ok := siteData["leaf_bgp_hold_down_timer"]; ok {
		switch v := val.(type) {
		case float64:
			state.LeafBgpHoldDownTimer = types.Int64Value(int64(v))
		case int:
			state.LeafBgpHoldDownTimer = types.Int64Value(int64(v))
		case int32:
			state.LeafBgpHoldDownTimer = types.Int64Value(int64(v))
		case nil:
			state.LeafBgpHoldDownTimer = types.Int64Null()
		default:
			if plan != nil && !plan.LeafBgpHoldDownTimer.IsNull() {
				state.LeafBgpHoldDownTimer = plan.LeafBgpHoldDownTimer
			} else {
				state.LeafBgpHoldDownTimer = types.Int64Null()
			}
		}
	} else if plan != nil && !plan.LeafBgpHoldDownTimer.IsNull() {
		state.LeafBgpHoldDownTimer = plan.LeafBgpHoldDownTimer
	} else {
		state.LeafBgpHoldDownTimer = types.Int64Null()
	}

	if val, ok := siteData["leaf_bgp_advertisement_interval"]; ok {
		switch v := val.(type) {
		case float64:
			state.LeafBgpAdvertisementInterval = types.Int64Value(int64(v))
		case int:
			state.LeafBgpAdvertisementInterval = types.Int64Value(int64(v))
		case int32:
			state.LeafBgpAdvertisementInterval = types.Int64Value(int64(v))
		case nil:
			state.LeafBgpAdvertisementInterval = types.Int64Null()
		default:
			if plan != nil && !plan.LeafBgpAdvertisementInterval.IsNull() {
				state.LeafBgpAdvertisementInterval = plan.LeafBgpAdvertisementInterval
			} else {
				state.LeafBgpAdvertisementInterval = types.Int64Null()
			}
		}
	} else if plan != nil && !plan.LeafBgpAdvertisementInterval.IsNull() {
		state.LeafBgpAdvertisementInterval = plan.LeafBgpAdvertisementInterval
	} else {
		state.LeafBgpAdvertisementInterval = types.Int64Null()
	}

	if val, ok := siteData["leaf_bgp_connect_timer"]; ok {
		switch v := val.(type) {
		case float64:
			state.LeafBgpConnectTimer = types.Int64Value(int64(v))
		case int:
			state.LeafBgpConnectTimer = types.Int64Value(int64(v))
		case int32:
			state.LeafBgpConnectTimer = types.Int64Value(int64(v))
		case nil:
			state.LeafBgpConnectTimer = types.Int64Null()
		default:
			if plan != nil && !plan.LeafBgpConnectTimer.IsNull() {
				state.LeafBgpConnectTimer = plan.LeafBgpConnectTimer
			} else {
				state.LeafBgpConnectTimer = types.Int64Null()
			}
		}
	} else if plan != nil && !plan.LeafBgpConnectTimer.IsNull() {
		state.LeafBgpConnectTimer = plan.LeafBgpConnectTimer
	} else {
		state.LeafBgpConnectTimer = types.Int64Null()
	}

	if val, ok := siteData["link_state_timeout_value"]; ok {
		switch v := val.(type) {
		case float64:
			state.LinkStateTimeoutValue = types.Int64Value(int64(v))
		case int:
			state.LinkStateTimeoutValue = types.Int64Value(int64(v))
		case int32:
			state.LinkStateTimeoutValue = types.Int64Value(int64(v))
		case nil:
			state.LinkStateTimeoutValue = types.Int64Null()
		default:
			if plan != nil && !plan.LinkStateTimeoutValue.IsNull() {
				state.LinkStateTimeoutValue = plan.LinkStateTimeoutValue
			} else {
				state.LinkStateTimeoutValue = types.Int64Null()
			}
		}
	} else if plan != nil && !plan.LinkStateTimeoutValue.IsNull() {
		state.LinkStateTimeoutValue = plan.LinkStateTimeoutValue
	} else {
		state.LinkStateTimeoutValue = types.Int64Null()
	}

	if val, ok := siteData["evpn_multihoming_startup_delay"]; ok {
		switch v := val.(type) {
		case float64:
			state.EvpnMultihomingStartupDelay = types.Int64Value(int64(v))
		case int:
			state.EvpnMultihomingStartupDelay = types.Int64Value(int64(v))
		case int32:
			state.EvpnMultihomingStartupDelay = types.Int64Value(int64(v))
		case nil:
			state.EvpnMultihomingStartupDelay = types.Int64Null()
		default:
			if plan != nil && !plan.EvpnMultihomingStartupDelay.IsNull() {
				state.EvpnMultihomingStartupDelay = plan.EvpnMultihomingStartupDelay
			} else {
				state.EvpnMultihomingStartupDelay = types.Int64Null()
			}
		}
	} else if plan != nil && !plan.EvpnMultihomingStartupDelay.IsNull() {
		state.EvpnMultihomingStartupDelay = plan.EvpnMultihomingStartupDelay
	} else {
		state.EvpnMultihomingStartupDelay = types.Int64Null()
	}

	if val, ok := siteData["evpn_mac_holdtime"]; ok {
		switch v := val.(type) {
		case float64:
			state.EvpnMacHoldtime = types.Int64Value(int64(v))
		case int:
			state.EvpnMacHoldtime = types.Int64Value(int64(v))
		case int32:
			state.EvpnMacHoldtime = types.Int64Value(int64(v))
		case nil:
			state.EvpnMacHoldtime = types.Int64Null()
		default:
			if plan != nil && !plan.EvpnMacHoldtime.IsNull() {
				state.EvpnMacHoldtime = plan.EvpnMacHoldtime
			} else {
				state.EvpnMacHoldtime = types.Int64Null()
			}
		}
	} else if plan != nil && !plan.EvpnMacHoldtime.IsNull() {
		state.EvpnMacHoldtime = plan.EvpnMacHoldtime
	} else {
		state.EvpnMacHoldtime = types.Int64Null()
	}

	if val, ok := siteData["aggressive_reporting"].(bool); ok {
		state.AggressiveReporting = types.BoolValue(val)
	} else if plan != nil && !plan.AggressiveReporting.IsNull() {
		state.AggressiveReporting = plan.AggressiveReporting
	} else {
		state.AggressiveReporting = types.BoolNull()
	}

	if val, ok := siteData["crc_failure_threshold"]; ok {
		switch v := val.(type) {
		case float64:
			state.CrcFailureThreshold = types.Int64Value(int64(v))
		case int:
			state.CrcFailureThreshold = types.Int64Value(int64(v))
		case int32:
			state.CrcFailureThreshold = types.Int64Value(int64(v))
		case nil:
			state.CrcFailureThreshold = types.Int64Null()
		default:
			if plan != nil && !plan.CrcFailureThreshold.IsNull() {
				state.CrcFailureThreshold = plan.CrcFailureThreshold
			} else {
				state.CrcFailureThreshold = types.Int64Null()
			}
		}
	} else if plan != nil && !plan.CrcFailureThreshold.IsNull() {
		state.CrcFailureThreshold = plan.CrcFailureThreshold
	} else {
		state.CrcFailureThreshold = types.Int64Null()
	}

	if val, ok := siteData["duplicate_address_detection_max_number_of_moves"]; ok {
		switch v := val.(type) {
		case float64:
			state.DuplicateAddressDetectionMaxNumberOfMoves = types.Int64Value(int64(v))
		case int:
			state.DuplicateAddressDetectionMaxNumberOfMoves = types.Int64Value(int64(v))
		case int32:
			state.DuplicateAddressDetectionMaxNumberOfMoves = types.Int64Value(int64(v))
		case nil:
			state.DuplicateAddressDetectionMaxNumberOfMoves = types.Int64Null()
		default:
			if plan != nil && !plan.DuplicateAddressDetectionMaxNumberOfMoves.IsNull() {
				state.DuplicateAddressDetectionMaxNumberOfMoves = plan.DuplicateAddressDetectionMaxNumberOfMoves
			} else {
				state.DuplicateAddressDetectionMaxNumberOfMoves = types.Int64Null()
			}
		}
	} else if plan != nil && !plan.DuplicateAddressDetectionMaxNumberOfMoves.IsNull() {
		state.DuplicateAddressDetectionMaxNumberOfMoves = plan.DuplicateAddressDetectionMaxNumberOfMoves
	} else {
		state.DuplicateAddressDetectionMaxNumberOfMoves = types.Int64Null()
	}

	if val, ok := siteData["duplicate_address_detection_time"]; ok {
		switch v := val.(type) {
		case float64:
			state.DuplicateAddressDetectionTime = types.Int64Value(int64(v))
		case int:
			state.DuplicateAddressDetectionTime = types.Int64Value(int64(v))
		case int32:
			state.DuplicateAddressDetectionTime = types.Int64Value(int64(v))
		case nil:
			state.DuplicateAddressDetectionTime = types.Int64Null()
		default:
			if plan != nil && !plan.DuplicateAddressDetectionTime.IsNull() {
				state.DuplicateAddressDetectionTime = plan.DuplicateAddressDetectionTime
			} else {
				state.DuplicateAddressDetectionTime = types.Int64Null()
			}
		}
	} else if plan != nil && !plan.DuplicateAddressDetectionTime.IsNull() {
		state.DuplicateAddressDetectionTime = plan.DuplicateAddressDetectionTime
	} else {
		state.DuplicateAddressDetectionTime = types.Int64Null()
	}

	if val, ok := siteData["enable_dhcp_snooping"].(bool); ok {
		state.EnableDhcpSnooping = types.BoolValue(val)
	} else if plan != nil && !plan.EnableDhcpSnooping.IsNull() {
		state.EnableDhcpSnooping = plan.EnableDhcpSnooping
	} else {
		state.EnableDhcpSnooping = types.BoolNull()
	}

	if val, ok := siteData["ip_source_guard"].(bool); ok {
		state.IpSourceGuard = types.BoolValue(val)
	} else if plan != nil && !plan.IpSourceGuard.IsNull() {
		state.IpSourceGuard = plan.IpSourceGuard
	} else {
		state.IpSourceGuard = types.BoolNull()
	}

	// Handle islands list
	if islands, ok := siteData["islands"].([]interface{}); ok {
		islandsList := make([]veritySiteIslandsModel, len(islands))
		for i, island := range islands {
			if islandMap, ok := island.(map[string]interface{}); ok {
				islandModel := veritySiteIslandsModel{}

				if val, exists := islandMap["toi_switchpoint"].(string); exists {
					islandModel.ToiSwitchpoint = types.StringValue(val)
				} else {
					islandModel.ToiSwitchpoint = types.StringNull()
				}

				if val, exists := islandMap["toi_switchpoint_ref_type_"].(string); exists {
					islandModel.ToiSwitchpointRefType = types.StringValue(val)
				} else {
					islandModel.ToiSwitchpointRefType = types.StringNull()
				}

				if val, exists := islandMap["index"]; exists {
					switch v := val.(type) {
					case float64:
						islandModel.Index = types.Int64Value(int64(v))
					case int:
						islandModel.Index = types.Int64Value(int64(v))
					case int32:
						islandModel.Index = types.Int64Value(int64(v))
					default:
						islandModel.Index = types.Int64Null()
					}
				} else {
					islandModel.Index = types.Int64Null()
				}

				islandsList[i] = islandModel
			}
		}
		state.Islands = islandsList
	} else if plan != nil {
		state.Islands = plan.Islands
	}

	// Handle pairs list
	if pairs, ok := siteData["pairs"].([]interface{}); ok {
		pairsList := make([]veritySitePairsModel, len(pairs))
		for i, pair := range pairs {
			if pairMap, ok := pair.(map[string]interface{}); ok {
				pairModel := veritySitePairsModel{}

				if val, exists := pairMap["name"].(string); exists {
					pairModel.Name = types.StringValue(val)
				} else {
					pairModel.Name = types.StringNull()
				}

				if val, exists := pairMap["switchpoint_1"].(string); exists {
					pairModel.Switchpoint1 = types.StringValue(val)
				} else {
					pairModel.Switchpoint1 = types.StringNull()
				}

				if val, exists := pairMap["switchpoint_1_ref_type_"].(string); exists {
					pairModel.Switchpoint1RefType = types.StringValue(val)
				} else {
					pairModel.Switchpoint1RefType = types.StringNull()
				}

				if val, exists := pairMap["switchpoint_2"].(string); exists {
					pairModel.Switchpoint2 = types.StringValue(val)
				} else {
					pairModel.Switchpoint2 = types.StringNull()
				}

				if val, exists := pairMap["switchpoint_2_ref_type_"].(string); exists {
					pairModel.Switchpoint2RefType = types.StringValue(val)
				} else {
					pairModel.Switchpoint2RefType = types.StringNull()
				}

				if val, exists := pairMap["lag_group"].(string); exists {
					pairModel.LagGroup = types.StringValue(val)
				} else {
					pairModel.LagGroup = types.StringNull()
				}

				if val, exists := pairMap["lag_group_ref_type_"].(string); exists {
					pairModel.LagGroupRefType = types.StringValue(val)
				} else {
					pairModel.LagGroupRefType = types.StringNull()
				}

				if val, exists := pairMap["is_whitebox_pair"].(bool); exists {
					pairModel.IsWhiteboxPair = types.BoolValue(val)
				} else {
					pairModel.IsWhiteboxPair = types.BoolNull()
				}

				if val, exists := pairMap["index"]; exists {
					switch v := val.(type) {
					case float64:
						pairModel.Index = types.Int64Value(int64(v))
					case int:
						pairModel.Index = types.Int64Value(int64(v))
					case int32:
						pairModel.Index = types.Int64Value(int64(v))
					default:
						pairModel.Index = types.Int64Null()
					}
				} else {
					pairModel.Index = types.Int64Null()
				}

				pairsList[i] = pairModel
			}
		}
		state.Pairs = pairsList
	} else if plan != nil {
		state.Pairs = plan.Pairs
	}

	// Handle object properties
	if op, ok := siteData["object_properties"].(map[string]interface{}); ok {
		objProps := &veritySiteObjectPropertiesModel{}

		if systemGraphs, exists := op["system_graphs"].([]interface{}); exists {
			graphsList := make([]veritySiteSystemGraphsModel, len(systemGraphs))
			for i, graph := range systemGraphs {
				if graphMap, ok := graph.(map[string]interface{}); ok {
					graphModel := veritySiteSystemGraphsModel{}

					if val, graphExists := graphMap["graph_num_data"].(string); graphExists {
						graphModel.GraphNumData = types.StringValue(val)
					} else {
						graphModel.GraphNumData = types.StringNull()
					}

					if val, graphExists := graphMap["index"]; graphExists {
						switch v := val.(type) {
						case float64:
							graphModel.Index = types.Int64Value(int64(v))
						case int:
							graphModel.Index = types.Int64Value(int64(v))
						case int32:
							graphModel.Index = types.Int64Value(int64(v))
						default:
							graphModel.Index = types.Int64Null()
						}
					} else {
						graphModel.Index = types.Int64Null()
					}

					graphsList[i] = graphModel
				}
			}
			objProps.SystemGraphs = graphsList
		} else if len(state.ObjectProperties) > 0 {
			objProps.SystemGraphs = state.ObjectProperties[0].SystemGraphs
		} else {
			objProps.SystemGraphs = []veritySiteSystemGraphsModel{}
		}

		state.ObjectProperties = []veritySiteObjectPropertiesModel{*objProps}
	} else if plan != nil && len(plan.ObjectProperties) > 0 {
		state.ObjectProperties = plan.ObjectProperties
	}

	return state
}

func (r *veritySiteResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Skip modification if we're deleting the resource
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan veritySiteResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate auto-assigned field specifications in configuration when auto-assigned
	// Check the actual configuration, not the plan
	var config veritySiteResourceModel
	if !req.Config.Raw.IsNull() {
		resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if !config.AnycastMacAddressAutoAssigned.IsNull() && config.AnycastMacAddressAutoAssigned.ValueBool() {
			if !config.AnycastMacAddress.IsNull() && !config.AnycastMacAddress.IsUnknown() && config.AnycastMacAddress.ValueString() != "" {
				resp.Diagnostics.AddError(
					"Anycast MAC Address cannot be specified when auto-assigned",
					"The 'anycast_mac_address' field cannot be specified in the configuration when 'anycast_mac_address_auto_assigned_' is set to true. The API will assign this value automatically.",
				)
				return
			}
		}
	}

	// For new resources (where state is null), mark auto-assigned fields as Unknown
	if req.State.Raw.IsNull() {
		if !plan.AnycastMacAddressAutoAssigned.IsNull() && plan.AnycastMacAddressAutoAssigned.ValueBool() {
			resp.Plan.SetAttribute(ctx, path.Root("anycast_mac_address"), types.StringUnknown())
		}
		return
	}

	var state veritySiteResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Handle auto-assigned field behavior
	if !plan.AnycastMacAddressAutoAssigned.IsNull() && plan.AnycastMacAddressAutoAssigned.ValueBool() {
		if !plan.AnycastMacAddressAutoAssigned.Equal(state.AnycastMacAddressAutoAssigned) {
			// anycast_mac_address_auto_assigned_ is changing to true, API will assign the value
			resp.Plan.SetAttribute(ctx, path.Root("anycast_mac_address"), types.StringUnknown())
			resp.Diagnostics.AddWarning(
				"Anycast MAC Address will be assigned by the API",
				"The 'anycast_mac_address' field will be automatically assigned by the API because 'anycast_mac_address_auto_assigned_' is being set to true.",
			)
		} else if !plan.AnycastMacAddress.Equal(state.AnycastMacAddress) {
			// User tried to change AnycastMacAddress but it's auto-assigned
			resp.Diagnostics.AddWarning(
				"Ignoring anycast_mac_address changes with auto-assignment enabled",
				"The 'anycast_mac_address' field changes will be ignored because 'anycast_mac_address_auto_assigned_' is set to true. The API will assign this value automatically.",
			)
			// Keep the current state value to suppress the diff
			if !state.AnycastMacAddress.IsNull() {
				resp.Plan.SetAttribute(ctx, path.Root("anycast_mac_address"), state.AnycastMacAddress)
			}
		}
	}
}
