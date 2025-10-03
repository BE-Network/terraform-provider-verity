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

	// Handle string fields
	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "PermitDeny", APIField: &routeMapClauseProps.PermitDeny, TFValue: plan.PermitDeny},
		{FieldName: "MatchAsPathAccessList", APIField: &routeMapClauseProps.MatchAsPathAccessList, TFValue: plan.MatchAsPathAccessList},
		{FieldName: "MatchAsPathAccessListRefType", APIField: &routeMapClauseProps.MatchAsPathAccessListRefType, TFValue: plan.MatchAsPathAccessListRefType},
		{FieldName: "MatchCommunityList", APIField: &routeMapClauseProps.MatchCommunityList, TFValue: plan.MatchCommunityList},
		{FieldName: "MatchCommunityListRefType", APIField: &routeMapClauseProps.MatchCommunityListRefType, TFValue: plan.MatchCommunityListRefType},
		{FieldName: "MatchExtendedCommunityList", APIField: &routeMapClauseProps.MatchExtendedCommunityList, TFValue: plan.MatchExtendedCommunityList},
		{FieldName: "MatchExtendedCommunityListRefType", APIField: &routeMapClauseProps.MatchExtendedCommunityListRefType, TFValue: plan.MatchExtendedCommunityListRefType},
		{FieldName: "MatchIpv4AddressIpPrefixList", APIField: &routeMapClauseProps.MatchIpv4AddressIpPrefixList, TFValue: plan.MatchIpv4AddressIpPrefixList},
		{FieldName: "MatchIpv4AddressIpPrefixListRefType", APIField: &routeMapClauseProps.MatchIpv4AddressIpPrefixListRefType, TFValue: plan.MatchIpv4AddressIpPrefixListRefType},
		{FieldName: "MatchIpv4NextHopIpPrefixList", APIField: &routeMapClauseProps.MatchIpv4NextHopIpPrefixList, TFValue: plan.MatchIpv4NextHopIpPrefixList},
		{FieldName: "MatchIpv4NextHopIpPrefixListRefType", APIField: &routeMapClauseProps.MatchIpv4NextHopIpPrefixListRefType, TFValue: plan.MatchIpv4NextHopIpPrefixListRefType},
		{FieldName: "MatchOrigin", APIField: &routeMapClauseProps.MatchOrigin, TFValue: plan.MatchOrigin},
		{FieldName: "MatchPeerIpAddress", APIField: &routeMapClauseProps.MatchPeerIpAddress, TFValue: plan.MatchPeerIpAddress},
		{FieldName: "MatchSourceProtocol", APIField: &routeMapClauseProps.MatchSourceProtocol, TFValue: plan.MatchSourceProtocol},
		{FieldName: "MatchVrf", APIField: &routeMapClauseProps.MatchVrf, TFValue: plan.MatchVrf},
		{FieldName: "MatchVrfRefType", APIField: &routeMapClauseProps.MatchVrfRefType, TFValue: plan.MatchVrfRefType},
		{FieldName: "MatchEvpnRouteType", APIField: &routeMapClauseProps.MatchEvpnRouteType, TFValue: plan.MatchEvpnRouteType},
		{FieldName: "MatchIpv6AddressIpv6PrefixList", APIField: &routeMapClauseProps.MatchIpv6AddressIpv6PrefixList, TFValue: plan.MatchIpv6AddressIpv6PrefixList},
		{FieldName: "MatchIpv6AddressIpv6PrefixListRefType", APIField: &routeMapClauseProps.MatchIpv6AddressIpv6PrefixListRefType, TFValue: plan.MatchIpv6AddressIpv6PrefixListRefType},
		{FieldName: "MatchIpv6NextHopIpv6PrefixList", APIField: &routeMapClauseProps.MatchIpv6NextHopIpv6PrefixList, TFValue: plan.MatchIpv6NextHopIpv6PrefixList},
		{FieldName: "MatchIpv6NextHopIpv6PrefixListRefType", APIField: &routeMapClauseProps.MatchIpv6NextHopIpv6PrefixListRefType, TFValue: plan.MatchIpv6NextHopIpv6PrefixListRefType},
	})

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &routeMapClauseProps.Enable, TFValue: plan.Enable},
		{FieldName: "MatchEvpnRouteTypeDefault", APIField: &routeMapClauseProps.MatchEvpnRouteTypeDefault, TFValue: plan.MatchEvpnRouteTypeDefault},
	})

	// Handle nullable int64 fields
	utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
		{FieldName: "MatchInterfaceNumber", APIField: &routeMapClauseProps.MatchInterfaceNumber, TFValue: plan.MatchInterfaceNumber},
		{FieldName: "MatchInterfaceVlan", APIField: &routeMapClauseProps.MatchInterfaceVlan, TFValue: plan.MatchInterfaceVlan},
		{FieldName: "MatchLocalPreference", APIField: &routeMapClauseProps.MatchLocalPreference, TFValue: plan.MatchLocalPreference},
		{FieldName: "MatchMetric", APIField: &routeMapClauseProps.MatchMetric, TFValue: plan.MatchMetric},
		{FieldName: "MatchPeerInterface", APIField: &routeMapClauseProps.MatchPeerInterface, TFValue: plan.MatchPeerInterface},
		{FieldName: "MatchPeerVlan", APIField: &routeMapClauseProps.MatchPeerVlan, TFValue: plan.MatchPeerVlan},
		{FieldName: "MatchTag", APIField: &routeMapClauseProps.MatchTag, TFValue: plan.MatchTag},
		{FieldName: "MatchVni", APIField: &routeMapClauseProps.MatchVni, TFValue: plan.MatchVni},
	})

	// Handle object properties
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

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "route_map_clause", name, *routeMapClauseProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Route Map Clause %s creation operation completed successfully", name))
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

	tflog.Debug(ctx, fmt.Sprintf("No recent Route Map Clause operations found, performing normal verification for %s", name))

	type RouteMapClauseResponse struct {
		RouteMapClause map[string]map[string]interface{} `json:"route_map_clause"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "route_map_clauses", name,
		func() (RouteMapClauseResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch Route Map Clauses")
			respAPI, err := r.client.RouteMapClausesAPI.RoutemapclausesGet(ctx).Execute()
			if err != nil {
				return RouteMapClauseResponse{}, fmt.Errorf("error reading Route Map Clause: %v", err)
			}
			defer respAPI.Body.Close()

			var res RouteMapClauseResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return RouteMapClauseResponse{}, fmt.Errorf("failed to decode Route Map Clauses response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d Route Map Clauses from API", len(res.RouteMapClause)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Route Map Clause %s", name))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for Route Map Clause with name: %s", name))

	routeMapClauseData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.RouteMapClause,
		name,
		func(data map[string]interface{}) (string, bool) {
			if name, ok := data["name"].(string); ok {
				return name, true
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Route Map Clause with name '%s' not found in API response", name))
		resp.State.RemoveResource(ctx)
		return
	}

	routeMapClauseMap, ok := (interface{}(routeMapClauseData)).(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Route Map Clause Data",
			fmt.Sprintf("Route Map Clause data is not in expected format for %s", name),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found Route Map Clause '%s' under API key '%s'", name, actualAPIName))

	state.Name = utils.MapStringFromAPI(routeMapClauseMap["name"])

	// Handle object properties
	if objProps, ok := routeMapClauseMap["object_properties"].(map[string]interface{}); ok {
		notes := utils.MapStringFromAPI(objProps["notes"])
		if notes.IsNull() {
			notes = types.StringValue("")
		}
		matchFieldsShown := utils.MapStringFromAPI(objProps["match_fields_shown"])
		if matchFieldsShown.IsNull() {
			matchFieldsShown = types.StringValue("")
		}
		state.ObjectProperties = []verityRouteMapClauseObjectPropertiesModel{
			{
				Notes:            notes,
				MatchFieldsShown: matchFieldsShown,
			},
		}
	} else {
		state.ObjectProperties = nil
	}

	// Map string fields
	stringFieldMappings := map[string]*types.String{
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

	for apiKey, stateField := range stringFieldMappings {
		*stateField = utils.MapStringFromAPI(routeMapClauseMap[apiKey])
	}

	// Map boolean fields
	boolFieldMappings := map[string]*types.Bool{
		"enable":                        &state.Enable,
		"match_evpn_route_type_default": &state.MatchEvpnRouteTypeDefault,
	}

	for apiKey, stateField := range boolFieldMappings {
		*stateField = utils.MapBoolFromAPI(routeMapClauseMap[apiKey])
	}

	// Map int64 fields
	int64FieldMappings := map[string]*types.Int64{
		"match_interface_number": &state.MatchInterfaceNumber,
		"match_interface_vlan":   &state.MatchInterfaceVlan,
		"match_local_preference": &state.MatchLocalPreference,
		"match_metric":           &state.MatchMetric,
		"match_peer_interface":   &state.MatchPeerInterface,
		"match_peer_vlan":        &state.MatchPeerVlan,
		"match_tag":              &state.MatchTag,
		"match_vni":              &state.MatchVni,
	}

	for apiKey, stateField := range int64FieldMappings {
		*stateField = utils.MapInt64FromAPI(routeMapClauseMap[apiKey])
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

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { routeMapClauseProps.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.PermitDeny, state.PermitDeny, func(v *string) { routeMapClauseProps.PermitDeny = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.MatchOrigin, state.MatchOrigin, func(v *string) { routeMapClauseProps.MatchOrigin = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.MatchPeerIpAddress, state.MatchPeerIpAddress, func(v *string) { routeMapClauseProps.MatchPeerIpAddress = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.MatchSourceProtocol, state.MatchSourceProtocol, func(v *string) { routeMapClauseProps.MatchSourceProtocol = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.MatchEvpnRouteType, state.MatchEvpnRouteType, func(v *string) { routeMapClauseProps.MatchEvpnRouteType = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { routeMapClauseProps.Enable = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.MatchEvpnRouteTypeDefault, state.MatchEvpnRouteTypeDefault, func(v *bool) { routeMapClauseProps.MatchEvpnRouteTypeDefault = v }, &hasChanges)

	// Handle nullable int64 field changes
	utils.CompareAndSetNullableInt64Field(plan.MatchInterfaceNumber, state.MatchInterfaceNumber, func(v *openapi.NullableInt32) { routeMapClauseProps.MatchInterfaceNumber = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.MatchInterfaceVlan, state.MatchInterfaceVlan, func(v *openapi.NullableInt32) { routeMapClauseProps.MatchInterfaceVlan = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.MatchLocalPreference, state.MatchLocalPreference, func(v *openapi.NullableInt32) { routeMapClauseProps.MatchLocalPreference = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.MatchMetric, state.MatchMetric, func(v *openapi.NullableInt32) { routeMapClauseProps.MatchMetric = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.MatchPeerInterface, state.MatchPeerInterface, func(v *openapi.NullableInt32) { routeMapClauseProps.MatchPeerInterface = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.MatchPeerVlan, state.MatchPeerVlan, func(v *openapi.NullableInt32) { routeMapClauseProps.MatchPeerVlan = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.MatchTag, state.MatchTag, func(v *openapi.NullableInt32) { routeMapClauseProps.MatchTag = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.MatchVni, state.MatchVni, func(v *openapi.NullableInt32) { routeMapClauseProps.MatchVni = *v }, &hasChanges)

	// Handle object properties
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

	// Handle MatchAsPathAccessList and MatchAsPathAccessListRefType using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.MatchAsPathAccessList, state.MatchAsPathAccessList, plan.MatchAsPathAccessListRefType, state.MatchAsPathAccessListRefType,
		func(v *string) { routeMapClauseProps.MatchAsPathAccessList = v },
		func(v *string) { routeMapClauseProps.MatchAsPathAccessListRefType = v },
		"match_as_path_access_list", "match_as_path_access_list_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	// Handle MatchCommunityList and MatchCommunityListRefType using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.MatchCommunityList, state.MatchCommunityList, plan.MatchCommunityListRefType, state.MatchCommunityListRefType,
		func(v *string) { routeMapClauseProps.MatchCommunityList = v },
		func(v *string) { routeMapClauseProps.MatchCommunityListRefType = v },
		"match_community_list", "match_community_list_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	// Handle MatchExtendedCommunityList and MatchExtendedCommunityListRefType using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.MatchExtendedCommunityList, state.MatchExtendedCommunityList, plan.MatchExtendedCommunityListRefType, state.MatchExtendedCommunityListRefType,
		func(v *string) { routeMapClauseProps.MatchExtendedCommunityList = v },
		func(v *string) { routeMapClauseProps.MatchExtendedCommunityListRefType = v },
		"match_extended_community_list", "match_extended_community_list_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	// Handle MatchIpv4AddressIpPrefixList and MatchIpv4AddressIpPrefixListRefType using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.MatchIpv4AddressIpPrefixList, state.MatchIpv4AddressIpPrefixList, plan.MatchIpv4AddressIpPrefixListRefType, state.MatchIpv4AddressIpPrefixListRefType,
		func(v *string) { routeMapClauseProps.MatchIpv4AddressIpPrefixList = v },
		func(v *string) { routeMapClauseProps.MatchIpv4AddressIpPrefixListRefType = v },
		"match_ipv4_address_ip_prefix_list", "match_ipv4_address_ip_prefix_list_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	// Handle MatchIpv4NextHopIpPrefixList and MatchIpv4NextHopIpPrefixListRefType using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.MatchIpv4NextHopIpPrefixList, state.MatchIpv4NextHopIpPrefixList, plan.MatchIpv4NextHopIpPrefixListRefType, state.MatchIpv4NextHopIpPrefixListRefType,
		func(v *string) { routeMapClauseProps.MatchIpv4NextHopIpPrefixList = v },
		func(v *string) { routeMapClauseProps.MatchIpv4NextHopIpPrefixListRefType = v },
		"match_ipv4_next_hop_ip_prefix_list", "match_ipv4_next_hop_ip_prefix_list_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	// Handle MatchVrf and MatchVrfRefType using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.MatchVrf, state.MatchVrf, plan.MatchVrfRefType, state.MatchVrfRefType,
		func(v *string) { routeMapClauseProps.MatchVrf = v },
		func(v *string) { routeMapClauseProps.MatchVrfRefType = v },
		"match_vrf", "match_vrf_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	// Handle MatchIpv6AddressIpv6PrefixList and MatchIpv6AddressIpv6PrefixListRefType using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.MatchIpv6AddressIpv6PrefixList, state.MatchIpv6AddressIpv6PrefixList, plan.MatchIpv6AddressIpv6PrefixListRefType, state.MatchIpv6AddressIpv6PrefixListRefType,
		func(v *string) { routeMapClauseProps.MatchIpv6AddressIpv6PrefixList = v },
		func(v *string) { routeMapClauseProps.MatchIpv6AddressIpv6PrefixListRefType = v },
		"match_ipv6_address_ipv6_prefix_list", "match_ipv6_address_ipv6_prefix_list_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	// Handle MatchIpv6NextHopIpv6PrefixList and MatchIpv6NextHopIpv6PrefixListRefType using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.MatchIpv6NextHopIpv6PrefixList, state.MatchIpv6NextHopIpv6PrefixList, plan.MatchIpv6NextHopIpv6PrefixListRefType, state.MatchIpv6NextHopIpv6PrefixListRefType,
		func(v *string) { routeMapClauseProps.MatchIpv6NextHopIpv6PrefixList = v },
		func(v *string) { routeMapClauseProps.MatchIpv6NextHopIpv6PrefixListRefType = v },
		"match_ipv6_next_hop_ipv6_prefix_list", "match_ipv6_next_hop_ipv6_prefix_list_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "route_map_clause", name, routeMapClauseProps, &resp.Diagnostics)
	if !success {
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

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "route_map_clause", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Route Map Clause %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "route_map_clauses")
	resp.State.RemoveResource(ctx)
}

func (r *verityRouteMapClauseResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
