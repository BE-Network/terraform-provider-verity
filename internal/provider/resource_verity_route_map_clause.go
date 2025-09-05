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
	_ resource.Resource                = &verityRouteMapClauseResource{}
	_ resource.ResourceWithConfigure   = &verityRouteMapClauseResource{}
	_ resource.ResourceWithImportState = &verityRouteMapClauseResource{}
)

func NewVerityRouteMapClauseResource() resource.Resource {
	return &verityRouteMapClauseResource{}
}

type verityRouteMapClauseResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *utils.BulkOperationManager
	notifyOperationAdded func()
}

type verityRouteMapClauseResourceModel struct {
	Name                                  types.String                                `tfsdk:"name"`
	Enable                                types.Bool                                  `tfsdk:"enable"`
	PermitDeny                            types.String                                `tfsdk:"permit_deny"`
	MatchAsPathAccessList                 types.String                                `tfsdk:"match_as_path_access_list"`
	MatchAsPathAccessListRefType          types.String                                `tfsdk:"match_as_path_access_list_ref_type_"`
	MatchCommunityList                    types.String                                `tfsdk:"match_community_list"`
	MatchCommunityListRefType             types.String                                `tfsdk:"match_community_list_ref_type_"`
	MatchExtendedCommunityList            types.String                                `tfsdk:"match_extended_community_list"`
	MatchExtendedCommunityListRefType     types.String                                `tfsdk:"match_extended_community_list_ref_type_"`
	MatchInterfaceNumber                  types.Int64                                 `tfsdk:"match_interface_number"`
	MatchInterfaceVlan                    types.Int64                                 `tfsdk:"match_interface_vlan"`
	MatchIpv4AddressIpPrefixList          types.String                                `tfsdk:"match_ipv4_address_ip_prefix_list"`
	MatchIpv4AddressIpPrefixListRefType   types.String                                `tfsdk:"match_ipv4_address_ip_prefix_list_ref_type_"`
	MatchIpv4NextHopIpPrefixList          types.String                                `tfsdk:"match_ipv4_next_hop_ip_prefix_list"`
	MatchIpv4NextHopIpPrefixListRefType   types.String                                `tfsdk:"match_ipv4_next_hop_ip_prefix_list_ref_type_"`
	MatchLocalPreference                  types.Int64                                 `tfsdk:"match_local_preference"`
	MatchMetric                           types.Int64                                 `tfsdk:"match_metric"`
	MatchOrigin                           types.String                                `tfsdk:"match_origin"`
	MatchPeerIpAddress                    types.String                                `tfsdk:"match_peer_ip_address"`
	MatchPeerInterface                    types.Int64                                 `tfsdk:"match_peer_interface"`
	MatchPeerVlan                         types.Int64                                 `tfsdk:"match_peer_vlan"`
	MatchSourceProtocol                   types.String                                `tfsdk:"match_source_protocol"`
	MatchVrf                              types.String                                `tfsdk:"match_vrf"`
	MatchVrfRefType                       types.String                                `tfsdk:"match_vrf_ref_type_"`
	MatchTag                              types.Int64                                 `tfsdk:"match_tag"`
	MatchEvpnRouteTypeDefault             types.Bool                                  `tfsdk:"match_evpn_route_type_default"`
	MatchEvpnRouteType                    types.String                                `tfsdk:"match_evpn_route_type"`
	MatchVni                              types.Int64                                 `tfsdk:"match_vni"`
	MatchIpv6AddressIpv6PrefixList        types.String                                `tfsdk:"match_ipv6_address_ipv6_prefix_list"`
	MatchIpv6AddressIpv6PrefixListRefType types.String                                `tfsdk:"match_ipv6_address_ipv6_prefix_list_ref_type_"`
	MatchIpv6NextHopIpv6PrefixList        types.String                                `tfsdk:"match_ipv6_next_hop_ipv6_prefix_list"`
	MatchIpv6NextHopIpv6PrefixListRefType types.String                                `tfsdk:"match_ipv6_next_hop_ipv6_prefix_list_ref_type_"`
	ObjectProperties                      []verityRouteMapClauseObjectPropertiesModel `tfsdk:"object_properties"`
}

type verityRouteMapClauseObjectPropertiesModel struct {
	Notes            types.String `tfsdk:"notes"`
	MatchFieldsShown types.String `tfsdk:"match_fields_shown"`
}

func (r *verityRouteMapClauseResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_route_map_clause"
}

func (r *verityRouteMapClauseResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityRouteMapClauseResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity Route Map Clause",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Object Name. Must be unique.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enable": schema.BoolAttribute{
				Description: "Enable flag of this provisioning object",
				Optional:    true,
			},
			"permit_deny": schema.StringAttribute{
				Description: "Action upon match of Community Strings.",
				Optional:    true,
			},
			"match_as_path_access_list": schema.StringAttribute{
				Description: "Match AS Path Access List",
				Optional:    true,
			},
			"match_as_path_access_list_ref_type_": schema.StringAttribute{
				Description: "Object type for match_as_path_access_list field",
				Optional:    true,
			},
			"match_community_list": schema.StringAttribute{
				Description: "Match Community List",
				Optional:    true,
			},
			"match_community_list_ref_type_": schema.StringAttribute{
				Description: "Object type for match_community_list field",
				Optional:    true,
			},
			"match_extended_community_list": schema.StringAttribute{
				Description: "Match Extended Community List",
				Optional:    true,
			},
			"match_extended_community_list_ref_type_": schema.StringAttribute{
				Description: "Object type for match_extended_community_list field",
				Optional:    true,
			},
			"match_interface_number": schema.Int64Attribute{
				Description: "Match Interface Number (minimum: 1, maximum: 256)",
				Optional:    true,
			},
			"match_interface_vlan": schema.Int64Attribute{
				Description: "Match Interface VLAN (minimum: 1, maximum: 4094)",
				Optional:    true,
			},
			"match_ipv4_address_ip_prefix_list": schema.StringAttribute{
				Description: "Match IPv4 Address IPv4 Prefix List",
				Optional:    true,
			},
			"match_ipv4_address_ip_prefix_list_ref_type_": schema.StringAttribute{
				Description: "Object type for match_ipv4_address_ip_prefix_list field",
				Optional:    true,
			},
			"match_ipv4_next_hop_ip_prefix_list": schema.StringAttribute{
				Description: "Match IPv4 Next Hop IPv4 Prefix List",
				Optional:    true,
			},
			"match_ipv4_next_hop_ip_prefix_list_ref_type_": schema.StringAttribute{
				Description: "Object type for match_ipv4_next_hop_ip_prefix_list field",
				Optional:    true,
			},
			"match_local_preference": schema.Int64Attribute{
				Description: "Match BGP Local Preference value on the route (maximum: 4294967295)",
				Optional:    true,
			},
			"match_metric": schema.Int64Attribute{
				Description: "Match Metric of the IP route entry (minimum: 1, maximum: 4294967295)",
				Optional:    true,
			},
			"match_origin": schema.StringAttribute{
				Description: "Match routes based on the value of the BGP Origin attribute",
				Optional:    true,
			},
			"match_peer_ip_address": schema.StringAttribute{
				Description: "Match BGP Peer IP Address the route was learned from",
				Optional:    true,
			},
			"match_peer_interface": schema.Int64Attribute{
				Description: "Match BGP Peer port the route was learned from (minimum: 1, maximum: 256)",
				Optional:    true,
			},
			"match_peer_vlan": schema.Int64Attribute{
				Description: "Match BGP Peer VLAN over which the route was learned (minimum: 1, maximum: 4094)",
				Optional:    true,
			},
			"match_source_protocol": schema.StringAttribute{
				Description: "Match Routing Protocol the route originated from",
				Optional:    true,
			},
			"match_vrf": schema.StringAttribute{
				Description: "Match VRF the route is associated with",
				Optional:    true,
			},
			"match_vrf_ref_type_": schema.StringAttribute{
				Description: "Object type for match_vrf field",
				Optional:    true,
			},
			"match_tag": schema.Int64Attribute{
				Description: "Match routes that have this value for a Tag attribute (minimum: 1, maximum: 4294967295)",
				Optional:    true,
			},
			"match_evpn_route_type_default": schema.BoolAttribute{
				Description: "Match based on the type of EVPN Route Type being Default",
				Optional:    true,
			},
			"match_evpn_route_type": schema.StringAttribute{
				Description: "Match based on the indicated EVPN Route Type",
				Optional:    true,
			},
			"match_vni": schema.Int64Attribute{
				Description: "Match based on the VNI value (minimum: 1, maximum: 16777215)",
				Optional:    true,
			},
			"match_ipv6_address_ipv6_prefix_list": schema.StringAttribute{
				Description: "Match IPv4 Address IPv6 Prefix List",
				Optional:    true,
			},
			"match_ipv6_address_ipv6_prefix_list_ref_type_": schema.StringAttribute{
				Description: "Object type for match_ipv6_address_ipv6_prefix_list field",
				Optional:    true,
			},
			"match_ipv6_next_hop_ipv6_prefix_list": schema.StringAttribute{
				Description: "Match IPv6 Next Hop IPv6 Prefix List",
				Optional:    true,
			},
			"match_ipv6_next_hop_ipv6_prefix_list_ref_type_": schema.StringAttribute{
				Description: "Object type for match_ipv6_next_hop_ipv6_prefix_list field",
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"object_properties": schema.ListNestedBlock{
				Description: "Object properties for the Route Map Clause",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"notes": schema.StringAttribute{
							Description: "User Notes.",
							Optional:    true,
						},
						"match_fields_shown": schema.StringAttribute{
							Description: "Match fields shown",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func (r *verityRouteMapClauseResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityRouteMapClauseResourceModel
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
	routeMapClauseProps := &openapi.RoutemapclausesPutRequestRouteMapClauseValue{
		Name: openapi.PtrString(name),
	}

	if !plan.Enable.IsNull() {
		routeMapClauseProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
	}

	if !plan.PermitDeny.IsNull() {
		routeMapClauseProps.PermitDeny = openapi.PtrString(plan.PermitDeny.ValueString())
	}

	if !plan.MatchAsPathAccessList.IsNull() {
		routeMapClauseProps.MatchAsPathAccessList = openapi.PtrString(plan.MatchAsPathAccessList.ValueString())
	}

	if !plan.MatchAsPathAccessListRefType.IsNull() {
		routeMapClauseProps.MatchAsPathAccessListRefType = openapi.PtrString(plan.MatchAsPathAccessListRefType.ValueString())
	}

	if !plan.MatchCommunityList.IsNull() {
		routeMapClauseProps.MatchCommunityList = openapi.PtrString(plan.MatchCommunityList.ValueString())
	}

	if !plan.MatchCommunityListRefType.IsNull() {
		routeMapClauseProps.MatchCommunityListRefType = openapi.PtrString(plan.MatchCommunityListRefType.ValueString())
	}

	if !plan.MatchExtendedCommunityList.IsNull() {
		routeMapClauseProps.MatchExtendedCommunityList = openapi.PtrString(plan.MatchExtendedCommunityList.ValueString())
	}

	if !plan.MatchExtendedCommunityListRefType.IsNull() {
		routeMapClauseProps.MatchExtendedCommunityListRefType = openapi.PtrString(plan.MatchExtendedCommunityListRefType.ValueString())
	}

	if !plan.MatchInterfaceNumber.IsNull() {
		val := int32(plan.MatchInterfaceNumber.ValueInt64())
		routeMapClauseProps.MatchInterfaceNumber = *openapi.NewNullableInt32(&val)
	} else {
		routeMapClauseProps.MatchInterfaceNumber = *openapi.NewNullableInt32(nil)
	}

	if !plan.MatchInterfaceVlan.IsNull() {
		val := int32(plan.MatchInterfaceVlan.ValueInt64())
		routeMapClauseProps.MatchInterfaceVlan = *openapi.NewNullableInt32(&val)
	} else {
		routeMapClauseProps.MatchInterfaceVlan = *openapi.NewNullableInt32(nil)
	}

	if !plan.MatchIpv4AddressIpPrefixList.IsNull() {
		routeMapClauseProps.MatchIpv4AddressIpPrefixList = openapi.PtrString(plan.MatchIpv4AddressIpPrefixList.ValueString())
	}

	if !plan.MatchIpv4AddressIpPrefixListRefType.IsNull() {
		routeMapClauseProps.MatchIpv4AddressIpPrefixListRefType = openapi.PtrString(plan.MatchIpv4AddressIpPrefixListRefType.ValueString())
	}

	if !plan.MatchIpv4NextHopIpPrefixList.IsNull() {
		routeMapClauseProps.MatchIpv4NextHopIpPrefixList = openapi.PtrString(plan.MatchIpv4NextHopIpPrefixList.ValueString())
	}

	if !plan.MatchIpv4NextHopIpPrefixListRefType.IsNull() {
		routeMapClauseProps.MatchIpv4NextHopIpPrefixListRefType = openapi.PtrString(plan.MatchIpv4NextHopIpPrefixListRefType.ValueString())
	}

	if !plan.MatchLocalPreference.IsNull() {
		val := int32(plan.MatchLocalPreference.ValueInt64())
		routeMapClauseProps.MatchLocalPreference = *openapi.NewNullableInt32(&val)
	} else {
		routeMapClauseProps.MatchLocalPreference = *openapi.NewNullableInt32(nil)
	}

	if !plan.MatchMetric.IsNull() {
		val := int32(plan.MatchMetric.ValueInt64())
		routeMapClauseProps.MatchMetric = *openapi.NewNullableInt32(&val)
	} else {
		routeMapClauseProps.MatchMetric = *openapi.NewNullableInt32(nil)
	}

	if !plan.MatchOrigin.IsNull() {
		routeMapClauseProps.MatchOrigin = openapi.PtrString(plan.MatchOrigin.ValueString())
	}

	if !plan.MatchPeerIpAddress.IsNull() {
		routeMapClauseProps.MatchPeerIpAddress = openapi.PtrString(plan.MatchPeerIpAddress.ValueString())
	}

	if !plan.MatchPeerInterface.IsNull() {
		val := int32(plan.MatchPeerInterface.ValueInt64())
		routeMapClauseProps.MatchPeerInterface = *openapi.NewNullableInt32(&val)
	} else {
		routeMapClauseProps.MatchPeerInterface = *openapi.NewNullableInt32(nil)
	}

	if !plan.MatchPeerVlan.IsNull() {
		val := int32(plan.MatchPeerVlan.ValueInt64())
		routeMapClauseProps.MatchPeerVlan = *openapi.NewNullableInt32(&val)
	} else {
		routeMapClauseProps.MatchPeerVlan = *openapi.NewNullableInt32(nil)
	}

	if !plan.MatchSourceProtocol.IsNull() {
		routeMapClauseProps.MatchSourceProtocol = openapi.PtrString(plan.MatchSourceProtocol.ValueString())
	}

	if !plan.MatchVrf.IsNull() {
		routeMapClauseProps.MatchVrf = openapi.PtrString(plan.MatchVrf.ValueString())
	}

	if !plan.MatchVrfRefType.IsNull() {
		routeMapClauseProps.MatchVrfRefType = openapi.PtrString(plan.MatchVrfRefType.ValueString())
	}

	if !plan.MatchTag.IsNull() {
		val := int32(plan.MatchTag.ValueInt64())
		routeMapClauseProps.MatchTag = *openapi.NewNullableInt32(&val)
	} else {
		routeMapClauseProps.MatchTag = *openapi.NewNullableInt32(nil)
	}

	if !plan.MatchEvpnRouteTypeDefault.IsNull() {
		routeMapClauseProps.MatchEvpnRouteTypeDefault = openapi.PtrBool(plan.MatchEvpnRouteTypeDefault.ValueBool())
	}

	if !plan.MatchEvpnRouteType.IsNull() {
		routeMapClauseProps.MatchEvpnRouteType = openapi.PtrString(plan.MatchEvpnRouteType.ValueString())
	}

	if !plan.MatchVni.IsNull() {
		val := int32(plan.MatchVni.ValueInt64())
		routeMapClauseProps.MatchVni = *openapi.NewNullableInt32(&val)
	} else {
		routeMapClauseProps.MatchVni = *openapi.NewNullableInt32(nil)
	}

	if !plan.MatchIpv6AddressIpv6PrefixList.IsNull() {
		routeMapClauseProps.MatchIpv6AddressIpv6PrefixList = openapi.PtrString(plan.MatchIpv6AddressIpv6PrefixList.ValueString())
	}

	if !plan.MatchIpv6AddressIpv6PrefixListRefType.IsNull() {
		routeMapClauseProps.MatchIpv6AddressIpv6PrefixListRefType = openapi.PtrString(plan.MatchIpv6AddressIpv6PrefixListRefType.ValueString())
	}

	if !plan.MatchIpv6NextHopIpv6PrefixList.IsNull() {
		routeMapClauseProps.MatchIpv6NextHopIpv6PrefixList = openapi.PtrString(plan.MatchIpv6NextHopIpv6PrefixList.ValueString())
	}

	if !plan.MatchIpv6NextHopIpv6PrefixListRefType.IsNull() {
		routeMapClauseProps.MatchIpv6NextHopIpv6PrefixListRefType = openapi.PtrString(plan.MatchIpv6NextHopIpv6PrefixListRefType.ValueString())
	}

	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objectProps := openapi.RoutemapclausesPutRequestRouteMapClauseValueObjectProperties{}
		if !op.Notes.IsNull() {
			objectProps.Notes = openapi.PtrString(op.Notes.ValueString())
		} else {
			objectProps.Notes = nil
		}
		if !op.MatchFieldsShown.IsNull() {
			objectProps.MatchFieldsShown = openapi.PtrString(op.MatchFieldsShown.ValueString())
		} else {
			objectProps.MatchFieldsShown = nil
		}
		routeMapClauseProps.ObjectProperties = &objectProps
	}

	operationID := r.bulkOpsMgr.AddPut(ctx, "route_map_clause", name, *routeMapClauseProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for Route Map Clause create operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Create Route Map Clause %s", name))...,
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Route Map Clause %s create operation completed successfully", name))
	clearCache(ctx, r.provCtx, "route_map_clauses")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
}

func (r *verityRouteMapClauseResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityRouteMapClauseResourceModel
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

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("route_map_clause") {
		tflog.Info(ctx, fmt.Sprintf("Skipping Route Map Clause %s verification – trusting recent successful API operation", name))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching Route Map Clauses for verification of %s", name))

	type RouteMapClauseResponse struct {
		RouteMapClause map[string]map[string]interface{} `json:"route_map_clause"`
	}

	var result RouteMapClauseResponse
	var err error
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		routeMapClausesData, fetchErr := getCachedResponse(ctx, r.provCtx, "route_map_clauses", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch route map clauses")
			respAPI, err := r.client.RouteMapClausesAPI.RoutemapclausesGet(ctx).Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading route map clauses: %v", err)
			}
			defer respAPI.Body.Close()

			var res RouteMapClauseResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return nil, fmt.Errorf("failed to decode route map clauses response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d route map clauses", len(res.RouteMapClause)))
			return res, nil
		})
		if fetchErr != nil {
			err = fetchErr
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch route map clauses on attempt %d, retrying in %v", attempt+1, sleepTime))
			time.Sleep(sleepTime)
			continue
		}
		result = routeMapClausesData.(RouteMapClauseResponse)
		break
	}
	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Route Map Clause %s", name))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for route map clause with ID: %s", name))
	var routeMapClauseData map[string]interface{}
	exists := false

	if data, ok := result.RouteMapClause[name]; ok {
		routeMapClauseData = data
		exists = true
		tflog.Debug(ctx, fmt.Sprintf("Found route map clause directly by ID: %s", name))
	} else {
		for apiName, p := range result.RouteMapClause {
			if nameVal, ok := p["name"].(string); ok && nameVal == name {
				routeMapClauseData = p
				name = apiName
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found route map clause with name '%s' under API key '%s'", nameVal, apiName))
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Route Map Clause with ID '%s' not found in API response", name))
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(fmt.Sprintf("%v", routeMapClauseData["name"]))

	if enable, ok := routeMapClauseData["enable"].(bool); ok {
		state.Enable = types.BoolValue(enable)
	} else {
		state.Enable = types.BoolNull()
	}

	// Only set object_properties if it exists in the API response
	if objProps, ok := routeMapClauseData["object_properties"].(map[string]interface{}); ok {
		if notes, ok := objProps["notes"].(string); ok {
			state.ObjectProperties = []verityRouteMapClauseObjectPropertiesModel{
				{Notes: types.StringValue(notes)},
			}
		} else {
			state.ObjectProperties = []verityRouteMapClauseObjectPropertiesModel{
				{Notes: types.StringValue("")},
			}
		}
		// Update match_fields_shown if it exists
		if len(state.ObjectProperties) > 0 {
			if matchFieldsShown, ok := objProps["match_fields_shown"].(string); ok {
				state.ObjectProperties[0].MatchFieldsShown = types.StringValue(matchFieldsShown)
			} else {
				state.ObjectProperties[0].MatchFieldsShown = types.StringValue("")
			}
		}
	} else {
		state.ObjectProperties = nil
	}

	stringAttrs := map[string]*types.String{
		"permit_deny":                                    &state.PermitDeny,
		"match_as_path_access_list":                      &state.MatchAsPathAccessList,
		"match_as_path_access_list_ref_type_":            &state.MatchAsPathAccessListRefType,
		"match_community_list":                           &state.MatchCommunityList,
		"match_community_list_ref_type_":                 &state.MatchCommunityListRefType,
		"match_extended_community_list":                  &state.MatchExtendedCommunityList,
		"match_extended_community_list_ref_type_":        &state.MatchExtendedCommunityListRefType,
		"match_ipv4_address_ip_prefix_list":              &state.MatchIpv4AddressIpPrefixList,
		"match_ipv4_address_ip_prefix_list_ref_type_":    &state.MatchIpv4AddressIpPrefixListRefType,
		"match_ipv4_next_hop_ip_prefix_list":             &state.MatchIpv4NextHopIpPrefixList,
		"match_ipv4_next_hop_ip_prefix_list_ref_type_":   &state.MatchIpv4NextHopIpPrefixListRefType,
		"match_origin":                                   &state.MatchOrigin,
		"match_peer_ip_address":                          &state.MatchPeerIpAddress,
		"match_source_protocol":                          &state.MatchSourceProtocol,
		"match_vrf":                                      &state.MatchVrf,
		"match_vrf_ref_type_":                            &state.MatchVrfRefType,
		"match_evpn_route_type":                          &state.MatchEvpnRouteType,
		"match_ipv6_address_ipv6_prefix_list":            &state.MatchIpv6AddressIpv6PrefixList,
		"match_ipv6_address_ipv6_prefix_list_ref_type_":  &state.MatchIpv6AddressIpv6PrefixListRefType,
		"match_ipv6_next_hop_ipv6_prefix_list":           &state.MatchIpv6NextHopIpv6PrefixList,
		"match_ipv6_next_hop_ipv6_prefix_list_ref_type_": &state.MatchIpv6NextHopIpv6PrefixListRefType,
	}

	for apiKey, stateField := range stringAttrs {
		if value, ok := routeMapClauseData[apiKey].(string); ok {
			*stateField = types.StringValue(value)
		} else {
			*stateField = types.StringNull()
		}
	}

	boolAttrs := map[string]*types.Bool{
		"match_evpn_route_type_default": &state.MatchEvpnRouteTypeDefault,
	}

	for apiKey, stateField := range boolAttrs {
		if value, ok := routeMapClauseData[apiKey].(bool); ok {
			*stateField = types.BoolValue(value)
		} else {
			*stateField = types.BoolNull()
		}
	}

	intAttrs := map[string]*types.Int64{
		"match_interface_number": &state.MatchInterfaceNumber,
		"match_interface_vlan":   &state.MatchInterfaceVlan,
		"match_local_preference": &state.MatchLocalPreference,
		"match_metric":           &state.MatchMetric,
		"match_peer_interface":   &state.MatchPeerInterface,
		"match_peer_vlan":        &state.MatchPeerVlan,
		"match_tag":              &state.MatchTag,
		"match_vni":              &state.MatchVni,
	}

	for apiKey, stateField := range intAttrs {
		if value, ok := routeMapClauseData[apiKey]; ok && value != nil {
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityRouteMapClauseResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityRouteMapClauseResourceModel

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
	routeMapClauseProps := openapi.RoutemapclausesPutRequestRouteMapClauseValue{}
	hasChanges := false

	if !plan.Enable.Equal(state.Enable) {
		routeMapClauseProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
		hasChanges = true
	}

	if !plan.PermitDeny.Equal(state.PermitDeny) {
		routeMapClauseProps.PermitDeny = openapi.PtrString(plan.PermitDeny.ValueString())
		hasChanges = true
	}

	// Handle MatchAsPathAccessList and MatchAsPathAccessListRefType according to "One ref type supported" rules
	matchAsPathAccessListChanged := !plan.MatchAsPathAccessList.Equal(state.MatchAsPathAccessList)
	matchAsPathAccessListRefTypeChanged := !plan.MatchAsPathAccessListRefType.Equal(state.MatchAsPathAccessListRefType)

	if matchAsPathAccessListChanged || matchAsPathAccessListRefTypeChanged {
		// Validate using "one ref type supported" rules
		if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
			plan.MatchAsPathAccessList, plan.MatchAsPathAccessListRefType,
			"match_as_path_access_list", "match_as_path_access_list_ref_type_",
			matchAsPathAccessListChanged, matchAsPathAccessListRefTypeChanged) {
			return
		}

		if matchAsPathAccessListChanged && !matchAsPathAccessListRefTypeChanged {
			// Only the base field changes, only the base field is sent
			if !plan.MatchAsPathAccessList.IsNull() && plan.MatchAsPathAccessList.ValueString() != "" {
				routeMapClauseProps.MatchAsPathAccessList = openapi.PtrString(plan.MatchAsPathAccessList.ValueString())
			} else {
				routeMapClauseProps.MatchAsPathAccessList = openapi.PtrString("")
			}
			hasChanges = true
		} else if matchAsPathAccessListRefTypeChanged {
			// ref_type changes (or both change), both fields are sent
			if !plan.MatchAsPathAccessList.IsNull() && plan.MatchAsPathAccessList.ValueString() != "" {
				routeMapClauseProps.MatchAsPathAccessList = openapi.PtrString(plan.MatchAsPathAccessList.ValueString())
			} else {
				routeMapClauseProps.MatchAsPathAccessList = openapi.PtrString("")
			}

			if !plan.MatchAsPathAccessListRefType.IsNull() && plan.MatchAsPathAccessListRefType.ValueString() != "" {
				routeMapClauseProps.MatchAsPathAccessListRefType = openapi.PtrString(plan.MatchAsPathAccessListRefType.ValueString())
			} else {
				routeMapClauseProps.MatchAsPathAccessListRefType = openapi.PtrString("")
			}
			hasChanges = true
		}
	}

	// Handle MatchCommunityList and MatchCommunityListRefType according to "One ref type supported" rules
	matchCommunityListChanged := !plan.MatchCommunityList.Equal(state.MatchCommunityList)
	matchCommunityListRefTypeChanged := !plan.MatchCommunityListRefType.Equal(state.MatchCommunityListRefType)

	if matchCommunityListChanged || matchCommunityListRefTypeChanged {
		if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
			plan.MatchCommunityList, plan.MatchCommunityListRefType,
			"match_community_list", "match_community_list_ref_type_",
			matchCommunityListChanged, matchCommunityListRefTypeChanged) {
			return
		}

		if matchCommunityListChanged && !matchCommunityListRefTypeChanged {
			if !plan.MatchCommunityList.IsNull() && plan.MatchCommunityList.ValueString() != "" {
				routeMapClauseProps.MatchCommunityList = openapi.PtrString(plan.MatchCommunityList.ValueString())
			} else {
				routeMapClauseProps.MatchCommunityList = openapi.PtrString("")
			}
			hasChanges = true
		} else if matchCommunityListRefTypeChanged {
			if !plan.MatchCommunityList.IsNull() && plan.MatchCommunityList.ValueString() != "" {
				routeMapClauseProps.MatchCommunityList = openapi.PtrString(plan.MatchCommunityList.ValueString())
			} else {
				routeMapClauseProps.MatchCommunityList = openapi.PtrString("")
			}

			if !plan.MatchCommunityListRefType.IsNull() && plan.MatchCommunityListRefType.ValueString() != "" {
				routeMapClauseProps.MatchCommunityListRefType = openapi.PtrString(plan.MatchCommunityListRefType.ValueString())
			} else {
				routeMapClauseProps.MatchCommunityListRefType = openapi.PtrString("")
			}
			hasChanges = true
		}
	}

	// Handle MatchExtendedCommunityList and MatchExtendedCommunityListRefType according to "One ref type supported" rules
	matchExtendedCommunityListChanged := !plan.MatchExtendedCommunityList.Equal(state.MatchExtendedCommunityList)
	matchExtendedCommunityListRefTypeChanged := !plan.MatchExtendedCommunityListRefType.Equal(state.MatchExtendedCommunityListRefType)

	if matchExtendedCommunityListChanged || matchExtendedCommunityListRefTypeChanged {
		if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
			plan.MatchExtendedCommunityList, plan.MatchExtendedCommunityListRefType,
			"match_extended_community_list", "match_extended_community_list_ref_type_",
			matchExtendedCommunityListChanged, matchExtendedCommunityListRefTypeChanged) {
			return
		}

		if matchExtendedCommunityListChanged && !matchExtendedCommunityListRefTypeChanged {
			if !plan.MatchExtendedCommunityList.IsNull() && plan.MatchExtendedCommunityList.ValueString() != "" {
				routeMapClauseProps.MatchExtendedCommunityList = openapi.PtrString(plan.MatchExtendedCommunityList.ValueString())
			} else {
				routeMapClauseProps.MatchExtendedCommunityList = openapi.PtrString("")
			}
			hasChanges = true
		} else if matchExtendedCommunityListRefTypeChanged {
			if !plan.MatchExtendedCommunityList.IsNull() && plan.MatchExtendedCommunityList.ValueString() != "" {
				routeMapClauseProps.MatchExtendedCommunityList = openapi.PtrString(plan.MatchExtendedCommunityList.ValueString())
			} else {
				routeMapClauseProps.MatchExtendedCommunityList = openapi.PtrString("")
			}

			if !plan.MatchExtendedCommunityListRefType.IsNull() && plan.MatchExtendedCommunityListRefType.ValueString() != "" {
				routeMapClauseProps.MatchExtendedCommunityListRefType = openapi.PtrString(plan.MatchExtendedCommunityListRefType.ValueString())
			} else {
				routeMapClauseProps.MatchExtendedCommunityListRefType = openapi.PtrString("")
			}
			hasChanges = true
		}
	}

	if !plan.MatchInterfaceNumber.Equal(state.MatchInterfaceNumber) {
		if !plan.MatchInterfaceNumber.IsNull() {
			val := int32(plan.MatchInterfaceNumber.ValueInt64())
			routeMapClauseProps.MatchInterfaceNumber = *openapi.NewNullableInt32(&val)
		} else {
			routeMapClauseProps.MatchInterfaceNumber = *openapi.NewNullableInt32(nil)
		}
		hasChanges = true
	}

	if !plan.MatchInterfaceVlan.Equal(state.MatchInterfaceVlan) {
		if !plan.MatchInterfaceVlan.IsNull() {
			val := int32(plan.MatchInterfaceVlan.ValueInt64())
			routeMapClauseProps.MatchInterfaceVlan = *openapi.NewNullableInt32(&val)
		} else {
			routeMapClauseProps.MatchInterfaceVlan = *openapi.NewNullableInt32(nil)
		}
		hasChanges = true
	}

	// Handle MatchIpv4AddressIpPrefixList and MatchIpv4AddressIpPrefixListRefType according to "One ref type supported" rules
	matchIpv4AddressIpPrefixListChanged := !plan.MatchIpv4AddressIpPrefixList.Equal(state.MatchIpv4AddressIpPrefixList)
	matchIpv4AddressIpPrefixListRefTypeChanged := !plan.MatchIpv4AddressIpPrefixListRefType.Equal(state.MatchIpv4AddressIpPrefixListRefType)

	if matchIpv4AddressIpPrefixListChanged || matchIpv4AddressIpPrefixListRefTypeChanged {
		if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
			plan.MatchIpv4AddressIpPrefixList, plan.MatchIpv4AddressIpPrefixListRefType,
			"match_ipv4_address_ip_prefix_list", "match_ipv4_address_ip_prefix_list_ref_type_",
			matchIpv4AddressIpPrefixListChanged, matchIpv4AddressIpPrefixListRefTypeChanged) {
			return
		}

		if matchIpv4AddressIpPrefixListChanged && !matchIpv4AddressIpPrefixListRefTypeChanged {
			if !plan.MatchIpv4AddressIpPrefixList.IsNull() && plan.MatchIpv4AddressIpPrefixList.ValueString() != "" {
				routeMapClauseProps.MatchIpv4AddressIpPrefixList = openapi.PtrString(plan.MatchIpv4AddressIpPrefixList.ValueString())
			} else {
				routeMapClauseProps.MatchIpv4AddressIpPrefixList = openapi.PtrString("")
			}
			hasChanges = true
		} else if matchIpv4AddressIpPrefixListRefTypeChanged {
			if !plan.MatchIpv4AddressIpPrefixList.IsNull() && plan.MatchIpv4AddressIpPrefixList.ValueString() != "" {
				routeMapClauseProps.MatchIpv4AddressIpPrefixList = openapi.PtrString(plan.MatchIpv4AddressIpPrefixList.ValueString())
			} else {
				routeMapClauseProps.MatchIpv4AddressIpPrefixList = openapi.PtrString("")
			}

			if !plan.MatchIpv4AddressIpPrefixListRefType.IsNull() && plan.MatchIpv4AddressIpPrefixListRefType.ValueString() != "" {
				routeMapClauseProps.MatchIpv4AddressIpPrefixListRefType = openapi.PtrString(plan.MatchIpv4AddressIpPrefixListRefType.ValueString())
			} else {
				routeMapClauseProps.MatchIpv4AddressIpPrefixListRefType = openapi.PtrString("")
			}
			hasChanges = true
		}
	}

	// Handle MatchIpv4NextHopIpPrefixList and MatchIpv4NextHopIpPrefixListRefType according to "One ref type supported" rules
	matchIpv4NextHopIpPrefixListChanged := !plan.MatchIpv4NextHopIpPrefixList.Equal(state.MatchIpv4NextHopIpPrefixList)
	matchIpv4NextHopIpPrefixListRefTypeChanged := !plan.MatchIpv4NextHopIpPrefixListRefType.Equal(state.MatchIpv4NextHopIpPrefixListRefType)

	if matchIpv4NextHopIpPrefixListChanged || matchIpv4NextHopIpPrefixListRefTypeChanged {
		if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
			plan.MatchIpv4NextHopIpPrefixList, plan.MatchIpv4NextHopIpPrefixListRefType,
			"match_ipv4_next_hop_ip_prefix_list", "match_ipv4_next_hop_ip_prefix_list_ref_type_",
			matchIpv4NextHopIpPrefixListChanged, matchIpv4NextHopIpPrefixListRefTypeChanged) {
			return
		}

		if matchIpv4NextHopIpPrefixListChanged && !matchIpv4NextHopIpPrefixListRefTypeChanged {
			if !plan.MatchIpv4NextHopIpPrefixList.IsNull() && plan.MatchIpv4NextHopIpPrefixList.ValueString() != "" {
				routeMapClauseProps.MatchIpv4NextHopIpPrefixList = openapi.PtrString(plan.MatchIpv4NextHopIpPrefixList.ValueString())
			} else {
				routeMapClauseProps.MatchIpv4NextHopIpPrefixList = openapi.PtrString("")
			}
			hasChanges = true
		} else if matchIpv4NextHopIpPrefixListRefTypeChanged {
			if !plan.MatchIpv4NextHopIpPrefixList.IsNull() && plan.MatchIpv4NextHopIpPrefixList.ValueString() != "" {
				routeMapClauseProps.MatchIpv4NextHopIpPrefixList = openapi.PtrString(plan.MatchIpv4NextHopIpPrefixList.ValueString())
			} else {
				routeMapClauseProps.MatchIpv4NextHopIpPrefixList = openapi.PtrString("")
			}

			if !plan.MatchIpv4NextHopIpPrefixListRefType.IsNull() && plan.MatchIpv4NextHopIpPrefixListRefType.ValueString() != "" {
				routeMapClauseProps.MatchIpv4NextHopIpPrefixListRefType = openapi.PtrString(plan.MatchIpv4NextHopIpPrefixListRefType.ValueString())
			} else {
				routeMapClauseProps.MatchIpv4NextHopIpPrefixListRefType = openapi.PtrString("")
			}
			hasChanges = true
		}
	}

	if !plan.MatchLocalPreference.Equal(state.MatchLocalPreference) {
		if !plan.MatchLocalPreference.IsNull() {
			val := int32(plan.MatchLocalPreference.ValueInt64())
			routeMapClauseProps.MatchLocalPreference = *openapi.NewNullableInt32(&val)
		} else {
			routeMapClauseProps.MatchLocalPreference = *openapi.NewNullableInt32(nil)
		}
		hasChanges = true
	}

	if !plan.MatchMetric.Equal(state.MatchMetric) {
		if !plan.MatchMetric.IsNull() {
			val := int32(plan.MatchMetric.ValueInt64())
			routeMapClauseProps.MatchMetric = *openapi.NewNullableInt32(&val)
		} else {
			routeMapClauseProps.MatchMetric = *openapi.NewNullableInt32(nil)
		}
		hasChanges = true
	}

	if !plan.MatchTag.Equal(state.MatchTag) {
		if !plan.MatchTag.IsNull() {
			val := int32(plan.MatchTag.ValueInt64())
			routeMapClauseProps.MatchTag = *openapi.NewNullableInt32(&val)
		} else {
			routeMapClauseProps.MatchTag = *openapi.NewNullableInt32(nil)
		}
		hasChanges = true
	}

	if len(plan.ObjectProperties) > 0 {
		if len(state.ObjectProperties) == 0 ||
			!plan.ObjectProperties[0].Notes.Equal(state.ObjectProperties[0].Notes) ||
			!plan.ObjectProperties[0].MatchFieldsShown.Equal(state.ObjectProperties[0].MatchFieldsShown) {
			routeMapClauseObjProps := openapi.RoutemapclausesPutRequestRouteMapClauseValueObjectProperties{}

			if !plan.ObjectProperties[0].Notes.IsNull() {
				routeMapClauseObjProps.Notes = openapi.PtrString(plan.ObjectProperties[0].Notes.ValueString())
			} else {
				routeMapClauseObjProps.Notes = nil
			}

			if !plan.ObjectProperties[0].MatchFieldsShown.IsNull() {
				routeMapClauseObjProps.MatchFieldsShown = openapi.PtrString(plan.ObjectProperties[0].MatchFieldsShown.ValueString())
			} else {
				routeMapClauseObjProps.MatchFieldsShown = nil
			}

			routeMapClauseProps.ObjectProperties = &routeMapClauseObjProps
			hasChanges = true
		}
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	operationID := r.bulkOpsMgr.AddPatch(ctx, "route_map_clause", name, routeMapClauseProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for Route Map Clause update operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update Route Map Clause %s", name))...,
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Route Map Clause %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "route_map_clauses")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *verityRouteMapClauseResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityRouteMapClauseResourceModel
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
	operationID := r.bulkOpsMgr.AddDelete(ctx, "route_map_clause", name)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for Route Map Clause deletion operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Delete Route Map Clause %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Route Map Clause %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "route_map_clauses")
	resp.State.RemoveResource(ctx)
}

func (r *verityRouteMapClauseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
