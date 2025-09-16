package provider

import (
	"context"
	"encoding/json"
	"fmt"
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
	Name                           types.String                      `tfsdk:"name"`
	Enable                         types.Bool                        `tfsdk:"enable"`
	ServiceForSite                 types.String                      `tfsdk:"service_for_site"`
	ServiceForSiteRefType          types.String                      `tfsdk:"service_for_site_ref_type_"`
	SpanningTreeType               types.String                      `tfsdk:"spanning_tree_type"`
	RegionName                     types.String                      `tfsdk:"region_name"`
	Revision                       types.Int64                       `tfsdk:"revision"`
	ForceSpanningTreeOnFabricPorts types.Bool                        `tfsdk:"force_spanning_tree_on_fabric_ports"`
	ReadOnlyMode                   types.Bool                        `tfsdk:"read_only_mode"`
	DscpToPBitMap                  types.String                      `tfsdk:"dscp_to_p_bit_map"`
	AnycastMacAddress              types.String                      `tfsdk:"anycast_mac_address"`
	AnycastMacAddressAutoAssigned  types.Bool                        `tfsdk:"anycast_mac_address_auto_assigned_"`
	MacAddressAgingTime            types.Int64                       `tfsdk:"mac_address_aging_time"`
	MlagDelayRestoreTimer          types.Int64                       `tfsdk:"mlag_delay_restore_timer"`
	BgpKeepaliveTimer              types.Int64                       `tfsdk:"bgp_keepalive_timer"`
	BgpHoldDownTimer               types.Int64                       `tfsdk:"bgp_hold_down_timer"`
	SpineBgpAdvertisementInterval  types.Int64                       `tfsdk:"spine_bgp_advertisement_interval"`
	SpineBgpConnectTimer           types.Int64                       `tfsdk:"spine_bgp_connect_timer"`
	LeafBgpKeepAliveTimer          types.Int64                       `tfsdk:"leaf_bgp_keep_alive_timer"`
	LeafBgpHoldDownTimer           types.Int64                       `tfsdk:"leaf_bgp_hold_down_timer"`
	LeafBgpAdvertisementInterval   types.Int64                       `tfsdk:"leaf_bgp_advertisement_interval"`
	LeafBgpConnectTimer            types.Int64                       `tfsdk:"leaf_bgp_connect_timer"`
	LinkStateTimeoutValue          types.Int64                       `tfsdk:"link_state_timeout_value"`
	EvpnMultihomingStartupDelay    types.Int64                       `tfsdk:"evpn_multihoming_startup_delay"`
	EvpnMacHoldtime                types.Int64                       `tfsdk:"evpn_mac_holdtime"`
	AggressiveReporting            types.Bool                        `tfsdk:"aggressive_reporting"`
	CrcFailureThreshold            types.Int64                       `tfsdk:"crc_failure_threshold"`
	EnableDhcpSnooping             types.Bool                        `tfsdk:"enable_dhcp_snooping"`
	IpSourceGuard                  types.Bool                        `tfsdk:"ip_source_guard"`
	Islands                        []veritySiteIslandsModel          `tfsdk:"islands"`
	Pairs                          []veritySitePairsModel            `tfsdk:"pairs"`
	ObjectProperties               []veritySiteObjectPropertiesModel `tfsdk:"object_properties"`
}

type veritySiteIslandsModel struct {
	ToiSwitchpoint        types.String `tfsdk:"toi_switchpoint"`
	ToiSwitchpointRefType types.String `tfsdk:"toi_switchpoint_ref_type_"`
	Index                 types.Int64  `tfsdk:"index"`
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

	tflog.Debug(ctx, "Reading site resource")

	provCtx := r.provCtx
	bulkOpsMgr := provCtx.bulkOpsMgr
	siteName := state.Name.ValueString()

	var siteData map[string]interface{}
	var exists bool

	if bulkOpsMgr != nil {
		siteData, exists = bulkOpsMgr.GetResourceResponse("site", siteName)
		if exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached site data for %s from recent operation", siteName))
			state = populateSiteState(ctx, state, siteData, nil)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if bulkOpsMgr != nil && bulkOpsMgr.HasPendingOrRecentOperations("site") {
		tflog.Info(ctx, fmt.Sprintf("Skipping site %s verification - trusting recent successful API operation", siteName))
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("No recent site operations found, performing normal verification for %s", siteName))

	type SitesResponse struct {
		Site map[string]interface{} `json:"site"`
	}

	var result SitesResponse
	var err error
	maxRetries := 3

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch sites on attempt %d, retrying in %v", attempt, sleepTime))
			time.Sleep(sleepTime)
		}

		sitesData, fetchErr := getCachedResponse(ctx, provCtx, "sites", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch sites")
			respAPI, err := r.client.SitesAPI.SitesGet(ctx).Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading sites: %v", err)
			}
			defer respAPI.Body.Close()

			var res SitesResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return nil, fmt.Errorf("failed to decode sites response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched sites data with %d sites", len(res.Site)))
			return res, nil
		})

		if fetchErr == nil {
			result = sitesData.(SitesResponse)
			break
		}
		err = fetchErr
	}

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Site %s", siteName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for site with ID: %s", siteName))
	if data, ok := result.Site[siteName].(map[string]interface{}); ok {
		siteData = data
		exists = true
		tflog.Debug(ctx, fmt.Sprintf("Found site directly by ID: %s", siteName))
	} else {
		for apiName, s := range result.Site {
			site, ok := s.(map[string]interface{})
			if !ok {
				continue
			}

			if name, ok := site["name"].(string); ok && name == siteName {
				siteData = site
				siteName = apiName
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found site with name '%s' under API key '%s'", name, apiName))
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Site with ID '%s' not found in API response", siteName))
		resp.State.RemoveResource(ctx)
		return
	}

	state = populateSiteState(ctx, state, siteData, nil)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *veritySiteResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state veritySiteResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
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
	siteReq := &openapi.SitesPatchRequestSiteValue{}
	hasChanges := false

	objPropsChanged := false
	if len(plan.ObjectProperties) > 0 && len(state.ObjectProperties) > 0 {
		if len(plan.ObjectProperties[0].SystemGraphs) != len(state.ObjectProperties[0].SystemGraphs) {
			objPropsChanged = true
		} else {
			for i, planGraph := range plan.ObjectProperties[0].SystemGraphs {
				if i >= len(state.ObjectProperties[0].SystemGraphs) ||
					!planGraph.GraphNumData.Equal(state.ObjectProperties[0].SystemGraphs[i].GraphNumData) ||
					!planGraph.Index.Equal(state.ObjectProperties[0].SystemGraphs[i].Index) {
					objPropsChanged = true
					break
				}
			}
		}
	} else if len(plan.ObjectProperties) != len(state.ObjectProperties) {
		objPropsChanged = true
	}

	if objPropsChanged {
		if len(plan.ObjectProperties) > 0 {
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
		} else {
			siteReq.ObjectProperties = nil
		}
		hasChanges = true
	}

	if !plan.Enable.Equal(state.Enable) {
		siteReq.Enable = openapi.PtrBool(plan.Enable.ValueBool())
		hasChanges = true
	}

	serviceForSiteChanged := !plan.ServiceForSite.Equal(state.ServiceForSite)
	serviceForSiteRefTypeChanged := !plan.ServiceForSiteRefType.Equal(state.ServiceForSiteRefType)

	if serviceForSiteChanged || serviceForSiteRefTypeChanged {
		// Validate using one ref type supported rules
		if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
			plan.ServiceForSite, plan.ServiceForSiteRefType,
			"service_for_site", "service_for_site_ref_type_",
			serviceForSiteChanged, serviceForSiteRefTypeChanged) {
			return
		}

		// Only send the base field if only it changed
		if serviceForSiteChanged && !serviceForSiteRefTypeChanged {
			// Just send the base field
			if !plan.ServiceForSite.IsNull() && plan.ServiceForSite.ValueString() != "" {
				siteReq.ServiceForSite = openapi.PtrString(plan.ServiceForSite.ValueString())
			} else {
				siteReq.ServiceForSite = openapi.PtrString("")
			}
			hasChanges = true
		} else if serviceForSiteRefTypeChanged {
			// Send both fields
			if !plan.ServiceForSite.IsNull() && plan.ServiceForSite.ValueString() != "" {
				siteReq.ServiceForSite = openapi.PtrString(plan.ServiceForSite.ValueString())
			} else {
				siteReq.ServiceForSite = openapi.PtrString("")
			}

			if !plan.ServiceForSiteRefType.IsNull() && plan.ServiceForSiteRefType.ValueString() != "" {
				siteReq.ServiceForSiteRefType = openapi.PtrString(plan.ServiceForSiteRefType.ValueString())
			} else {
				siteReq.ServiceForSiteRefType = openapi.PtrString("")
			}
			hasChanges = true
		}
	}

	if !plan.SpanningTreeType.Equal(state.SpanningTreeType) {
		siteReq.SpanningTreeType = openapi.PtrString(plan.SpanningTreeType.ValueString())
		hasChanges = true
	}

	if !plan.RegionName.Equal(state.RegionName) {
		siteReq.RegionName = openapi.PtrString(plan.RegionName.ValueString())
		hasChanges = true
	}

	if !plan.Revision.Equal(state.Revision) {
		if !plan.Revision.IsNull() {
			revisionVal := int32(plan.Revision.ValueInt64())
			siteReq.Revision = *openapi.NewNullableInt32(&revisionVal)
		} else {
			siteReq.Revision = *openapi.NewNullableInt32(nil)
		}
		hasChanges = true
	}

	if !plan.ForceSpanningTreeOnFabricPorts.Equal(state.ForceSpanningTreeOnFabricPorts) {
		siteReq.ForceSpanningTreeOnFabricPorts = openapi.PtrBool(plan.ForceSpanningTreeOnFabricPorts.ValueBool())
		hasChanges = true
	}

	if !plan.ReadOnlyMode.Equal(state.ReadOnlyMode) {
		siteReq.ReadOnlyMode = openapi.PtrBool(plan.ReadOnlyMode.ValueBool())
		hasChanges = true
	}

	if !plan.DscpToPBitMap.Equal(state.DscpToPBitMap) {
		siteReq.DscpToPBitMap = openapi.PtrString(plan.DscpToPBitMap.ValueString())
		hasChanges = true
	}

	// Handle AnycastMacAddress and AnycastMacAddressAutoAssigned changes in a coordinated way
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

	if !plan.MacAddressAgingTime.Equal(state.MacAddressAgingTime) {
		if !plan.MacAddressAgingTime.IsNull() {
			siteReq.MacAddressAgingTime = openapi.PtrInt32(int32(plan.MacAddressAgingTime.ValueInt64()))
		}
		hasChanges = true
	}

	if !plan.MlagDelayRestoreTimer.Equal(state.MlagDelayRestoreTimer) {
		if !plan.MlagDelayRestoreTimer.IsNull() {
			siteReq.MlagDelayRestoreTimer = openapi.PtrInt32(int32(plan.MlagDelayRestoreTimer.ValueInt64()))
		}
		hasChanges = true
	}

	if !plan.BgpKeepaliveTimer.Equal(state.BgpKeepaliveTimer) {
		if !plan.BgpKeepaliveTimer.IsNull() {
			siteReq.BgpKeepaliveTimer = openapi.PtrInt32(int32(plan.BgpKeepaliveTimer.ValueInt64()))
		}
		hasChanges = true
	}

	if !plan.BgpHoldDownTimer.Equal(state.BgpHoldDownTimer) {
		if !plan.BgpHoldDownTimer.IsNull() {
			siteReq.BgpHoldDownTimer = openapi.PtrInt32(int32(plan.BgpHoldDownTimer.ValueInt64()))
		}
		hasChanges = true
	}

	if !plan.SpineBgpAdvertisementInterval.Equal(state.SpineBgpAdvertisementInterval) {
		if !plan.SpineBgpAdvertisementInterval.IsNull() {
			siteReq.SpineBgpAdvertisementInterval = openapi.PtrInt32(int32(plan.SpineBgpAdvertisementInterval.ValueInt64()))
		}
		hasChanges = true
	}

	if !plan.SpineBgpConnectTimer.Equal(state.SpineBgpConnectTimer) {
		if !plan.SpineBgpConnectTimer.IsNull() {
			siteReq.SpineBgpConnectTimer = openapi.PtrInt32(int32(plan.SpineBgpConnectTimer.ValueInt64()))
		}
		hasChanges = true
	}

	if !plan.LeafBgpKeepAliveTimer.Equal(state.LeafBgpKeepAliveTimer) {
		if !plan.LeafBgpKeepAliveTimer.IsNull() {
			siteReq.LeafBgpKeepAliveTimer = openapi.PtrInt32(int32(plan.LeafBgpKeepAliveTimer.ValueInt64()))
		}
		hasChanges = true
	}

	if !plan.LeafBgpHoldDownTimer.Equal(state.LeafBgpHoldDownTimer) {
		if !plan.LeafBgpHoldDownTimer.IsNull() {
			siteReq.LeafBgpHoldDownTimer = openapi.PtrInt32(int32(plan.LeafBgpHoldDownTimer.ValueInt64()))
		}
		hasChanges = true
	}

	if !plan.LeafBgpAdvertisementInterval.Equal(state.LeafBgpAdvertisementInterval) {
		if !plan.LeafBgpAdvertisementInterval.IsNull() {
			siteReq.LeafBgpAdvertisementInterval = openapi.PtrInt32(int32(plan.LeafBgpAdvertisementInterval.ValueInt64()))
		}
		hasChanges = true
	}

	if !plan.LeafBgpConnectTimer.Equal(state.LeafBgpConnectTimer) {
		if !plan.LeafBgpConnectTimer.IsNull() {
			siteReq.LeafBgpConnectTimer = openapi.PtrInt32(int32(plan.LeafBgpConnectTimer.ValueInt64()))
		}
		hasChanges = true
	}

	if !plan.LinkStateTimeoutValue.Equal(state.LinkStateTimeoutValue) {
		if !plan.LinkStateTimeoutValue.IsNull() {
			val := int32(plan.LinkStateTimeoutValue.ValueInt64())
			siteReq.LinkStateTimeoutValue = *openapi.NewNullableInt32(&val)
		} else {
			siteReq.LinkStateTimeoutValue = *openapi.NewNullableInt32(nil)
		}
		hasChanges = true
	}

	if !plan.EvpnMultihomingStartupDelay.Equal(state.EvpnMultihomingStartupDelay) {
		if !plan.EvpnMultihomingStartupDelay.IsNull() {
			val := int32(plan.EvpnMultihomingStartupDelay.ValueInt64())
			siteReq.EvpnMultihomingStartupDelay = *openapi.NewNullableInt32(&val)
		} else {
			siteReq.EvpnMultihomingStartupDelay = *openapi.NewNullableInt32(nil)
		}
		hasChanges = true
	}

	if !plan.EvpnMacHoldtime.Equal(state.EvpnMacHoldtime) {
		if !plan.EvpnMacHoldtime.IsNull() {
			val := int32(plan.EvpnMacHoldtime.ValueInt64())
			siteReq.EvpnMacHoldtime = *openapi.NewNullableInt32(&val)
		} else {
			siteReq.EvpnMacHoldtime = *openapi.NewNullableInt32(nil)
		}
		hasChanges = true
	}

	if !plan.AggressiveReporting.Equal(state.AggressiveReporting) {
		siteReq.AggressiveReporting = openapi.PtrBool(plan.AggressiveReporting.ValueBool())
		hasChanges = true
	}

	if !plan.CrcFailureThreshold.Equal(state.CrcFailureThreshold) {
		if !plan.CrcFailureThreshold.IsNull() {
			val := int32(plan.CrcFailureThreshold.ValueInt64())
			siteReq.CrcFailureThreshold = *openapi.NewNullableInt32(&val)
		} else {
			siteReq.CrcFailureThreshold = *openapi.NewNullableInt32(nil)
		}
		hasChanges = true
	}

	if !plan.EnableDhcpSnooping.Equal(state.EnableDhcpSnooping) {
		siteReq.EnableDhcpSnooping = openapi.PtrBool(plan.EnableDhcpSnooping.ValueBool())
		hasChanges = true
	}

	if !plan.IpSourceGuard.Equal(state.IpSourceGuard) {
		siteReq.IpSourceGuard = openapi.PtrBool(plan.IpSourceGuard.ValueBool())
		hasChanges = true
	}

	islandsChanged := len(plan.Islands) != len(state.Islands)
	if !islandsChanged {
		for i, planIsland := range plan.Islands {
			if i >= len(state.Islands) ||
				!planIsland.ToiSwitchpoint.Equal(state.Islands[i].ToiSwitchpoint) ||
				!planIsland.ToiSwitchpointRefType.Equal(state.Islands[i].ToiSwitchpointRefType) ||
				!planIsland.Index.Equal(state.Islands[i].Index) {
				islandsChanged = true
				break
			}
		}
	}

	if islandsChanged {
		if len(plan.Islands) > 0 {
			islandsList := make([]openapi.SitesPatchRequestSiteValueIslandsInner, len(plan.Islands))
			for i, island := range plan.Islands {
				islandProps := openapi.SitesPatchRequestSiteValueIslandsInner{}
				if !island.ToiSwitchpoint.IsNull() {
					islandProps.ToiSwitchpoint = openapi.PtrString(island.ToiSwitchpoint.ValueString())
				}
				if !island.ToiSwitchpointRefType.IsNull() {
					islandProps.ToiSwitchpointRefType = openapi.PtrString(island.ToiSwitchpointRefType.ValueString())
				}
				if !island.Index.IsNull() {
					islandProps.Index = openapi.PtrInt32(int32(island.Index.ValueInt64()))
				}
				islandsList[i] = islandProps
			}
			siteReq.Islands = islandsList
		} else {
			siteReq.Islands = []openapi.SitesPatchRequestSiteValueIslandsInner{}
		}
		hasChanges = true
	}

	pairsChanged := len(plan.Pairs) != len(state.Pairs)
	if !pairsChanged {
		for i, planPair := range plan.Pairs {
			if i >= len(state.Pairs) ||
				!planPair.Name.Equal(state.Pairs[i].Name) ||
				!planPair.Switchpoint1.Equal(state.Pairs[i].Switchpoint1) ||
				!planPair.Switchpoint1RefType.Equal(state.Pairs[i].Switchpoint1RefType) ||
				!planPair.Switchpoint2.Equal(state.Pairs[i].Switchpoint2) ||
				!planPair.Switchpoint2RefType.Equal(state.Pairs[i].Switchpoint2RefType) ||
				!planPair.LagGroup.Equal(state.Pairs[i].LagGroup) ||
				!planPair.LagGroupRefType.Equal(state.Pairs[i].LagGroupRefType) ||
				!planPair.IsWhiteboxPair.Equal(state.Pairs[i].IsWhiteboxPair) ||
				!planPair.Index.Equal(state.Pairs[i].Index) {
				pairsChanged = true
				break
			}
		}
	}

	if pairsChanged {
		if len(plan.Pairs) > 0 {
			pairsList := make([]openapi.SitesPatchRequestSiteValuePairsInner, len(plan.Pairs))
			for i, pair := range plan.Pairs {
				pairProps := openapi.SitesPatchRequestSiteValuePairsInner{}
				if !pair.Name.IsNull() {
					pairProps.Name = openapi.PtrString(pair.Name.ValueString())
				}
				if !pair.Switchpoint1.IsNull() {
					pairProps.Switchpoint1 = openapi.PtrString(pair.Switchpoint1.ValueString())
				}
				if !pair.Switchpoint1RefType.IsNull() {
					pairProps.Switchpoint1RefType = openapi.PtrString(pair.Switchpoint1RefType.ValueString())
				}
				if !pair.Switchpoint2.IsNull() {
					pairProps.Switchpoint2 = openapi.PtrString(pair.Switchpoint2.ValueString())
				}
				if !pair.Switchpoint2RefType.IsNull() {
					pairProps.Switchpoint2RefType = openapi.PtrString(pair.Switchpoint2RefType.ValueString())
				}
				if !pair.LagGroup.IsNull() {
					pairProps.LagGroup = openapi.PtrString(pair.LagGroup.ValueString())
				}
				if !pair.LagGroupRefType.IsNull() {
					pairProps.LagGroupRefType = openapi.PtrString(pair.LagGroupRefType.ValueString())
				}
				if !pair.IsWhiteboxPair.IsNull() {
					pairProps.IsWhiteboxPair = openapi.PtrBool(pair.IsWhiteboxPair.ValueBool())
				}
				if !pair.Index.IsNull() {
					pairProps.Index = openapi.PtrInt32(int32(pair.Index.ValueInt64()))
				}
				pairsList[i] = pairProps
			}
			siteReq.Pairs = pairsList
		} else {
			siteReq.Pairs = []openapi.SitesPatchRequestSiteValuePairsInner{}
		}
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	bulkOpsMgr := r.provCtx.bulkOpsMgr
	operationID := bulkOpsMgr.AddPatch(ctx, "site", name, *siteReq)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for Site update operation %s to complete", operationID))
	if err := bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update Site %s", name))...,
		)
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

		// Following Gateway pattern - set as slice with single element
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
