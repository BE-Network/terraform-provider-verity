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
	_ resource.Resource                = &verityRouteMapClauseResource{}
	_ resource.ResourceWithConfigure   = &verityRouteMapClauseResource{}
	_ resource.ResourceWithImportState = &verityRouteMapClauseResource{}
	_ resource.ResourceWithModifyPlan  = &verityRouteMapClauseResource{}
)

const routeMapClauseResourceType = "routemapclauses"
const routeMapClauseTerraformType = "verity_route_map_clause"

func NewVerityRouteMapClauseResource() resource.Resource {
	return &verityRouteMapClauseResource{}
}

type verityRouteMapClauseResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
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
				Computed:    true,
			},
			"permit_deny": schema.StringAttribute{
				Description: "Action upon match of Community Strings.",
				Optional:    true,
				Computed:    true,
			},
			"match_as_path_access_list": schema.StringAttribute{
				Description: "Match AS Path Access List",
				Optional:    true,
				Computed:    true,
			},
			"match_as_path_access_list_ref_type_": schema.StringAttribute{
				Description: "Object type for match_as_path_access_list field",
				Optional:    true,
				Computed:    true,
			},
			"match_community_list": schema.StringAttribute{
				Description: "Match Community List",
				Optional:    true,
				Computed:    true,
			},
			"match_community_list_ref_type_": schema.StringAttribute{
				Description: "Object type for match_community_list field",
				Optional:    true,
				Computed:    true,
			},
			"match_extended_community_list": schema.StringAttribute{
				Description: "Match Extended Community List",
				Optional:    true,
				Computed:    true,
			},
			"match_extended_community_list_ref_type_": schema.StringAttribute{
				Description: "Object type for match_extended_community_list field",
				Optional:    true,
				Computed:    true,
			},
			"match_interface_number": schema.Int64Attribute{
				Description: "Match Interface Number (minimum: 1, maximum: 256)",
				Optional:    true,
				Computed:    true,
			},
			"match_interface_vlan": schema.Int64Attribute{
				Description: "Match Interface VLAN (minimum: 1, maximum: 4094)",
				Optional:    true,
				Computed:    true,
			},
			"match_ipv4_address_ip_prefix_list": schema.StringAttribute{
				Description: "Match IPv4 Address IPv4 Prefix List",
				Optional:    true,
				Computed:    true,
			},
			"match_ipv4_address_ip_prefix_list_ref_type_": schema.StringAttribute{
				Description: "Object type for match_ipv4_address_ip_prefix_list field",
				Optional:    true,
				Computed:    true,
			},
			"match_ipv4_next_hop_ip_prefix_list": schema.StringAttribute{
				Description: "Match IPv4 Next Hop IPv4 Prefix List",
				Optional:    true,
				Computed:    true,
			},
			"match_ipv4_next_hop_ip_prefix_list_ref_type_": schema.StringAttribute{
				Description: "Object type for match_ipv4_next_hop_ip_prefix_list field",
				Optional:    true,
				Computed:    true,
			},
			"match_local_preference": schema.Int64Attribute{
				Description: "Match BGP Local Preference value on the route (maximum: 4294967295)",
				Optional:    true,
				Computed:    true,
			},
			"match_metric": schema.Int64Attribute{
				Description: "Match Metric of the IP route entry (minimum: 1, maximum: 4294967295)",
				Optional:    true,
				Computed:    true,
			},
			"match_origin": schema.StringAttribute{
				Description: "Match routes based on the value of the BGP Origin attribute",
				Optional:    true,
				Computed:    true,
			},
			"match_peer_ip_address": schema.StringAttribute{
				Description: "Match BGP Peer IP Address the route was learned from",
				Optional:    true,
				Computed:    true,
			},
			"match_peer_interface": schema.Int64Attribute{
				Description: "Match BGP Peer port the route was learned from (minimum: 1, maximum: 256)",
				Optional:    true,
				Computed:    true,
			},
			"match_peer_vlan": schema.Int64Attribute{
				Description: "Match BGP Peer VLAN over which the route was learned (minimum: 1, maximum: 4094)",
				Optional:    true,
				Computed:    true,
			},
			"match_source_protocol": schema.StringAttribute{
				Description: "Match Routing Protocol the route originated from",
				Optional:    true,
				Computed:    true,
			},
			"match_vrf": schema.StringAttribute{
				Description: "Match VRF the route is associated with",
				Optional:    true,
				Computed:    true,
			},
			"match_vrf_ref_type_": schema.StringAttribute{
				Description: "Object type for match_vrf field",
				Optional:    true,
				Computed:    true,
			},
			"match_tag": schema.Int64Attribute{
				Description: "Match routes that have this value for a Tag attribute (minimum: 1, maximum: 4294967295)",
				Optional:    true,
				Computed:    true,
			},
			"match_evpn_route_type_default": schema.BoolAttribute{
				Description: "Match based on the type of EVPN Route Type being Default",
				Optional:    true,
				Computed:    true,
			},
			"match_evpn_route_type": schema.StringAttribute{
				Description: "Match based on the indicated EVPN Route Type",
				Optional:    true,
				Computed:    true,
			},
			"match_vni": schema.Int64Attribute{
				Description: "Match based on the VNI value (minimum: 1, maximum: 16777215)",
				Optional:    true,
				Computed:    true,
			},
			"match_ipv6_address_ipv6_prefix_list": schema.StringAttribute{
				Description: "Match IPv4 Address IPv6 Prefix List",
				Optional:    true,
				Computed:    true,
			},
			"match_ipv6_address_ipv6_prefix_list_ref_type_": schema.StringAttribute{
				Description: "Object type for match_ipv6_address_ipv6_prefix_list field",
				Optional:    true,
				Computed:    true,
			},
			"match_ipv6_next_hop_ipv6_prefix_list": schema.StringAttribute{
				Description: "Match IPv6 Next Hop IPv6 Prefix List",
				Optional:    true,
				Computed:    true,
			},
			"match_ipv6_next_hop_ipv6_prefix_list_ref_type_": schema.StringAttribute{
				Description: "Object type for match_ipv6_next_hop_ipv6_prefix_list field",
				Optional:    true,
				Computed:    true,
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
							Computed:    true,
						},
						"match_fields_shown": schema.StringAttribute{
							Description: "Match fields shown",
							Optional:    true,
							Computed:    true,
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

	var config verityRouteMapClauseResourceModel
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

	// Handle nullable int64 fields - parse HCL to detect explicit config
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, routeMapClauseTerraformType, name)

	utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
		{FieldName: "MatchInterfaceNumber", APIField: &routeMapClauseProps.MatchInterfaceNumber, TFValue: config.MatchInterfaceNumber, IsConfigured: configuredAttrs.IsConfigured("match_interface_number")},
		{FieldName: "MatchInterfaceVlan", APIField: &routeMapClauseProps.MatchInterfaceVlan, TFValue: config.MatchInterfaceVlan, IsConfigured: configuredAttrs.IsConfigured("match_interface_vlan")},
		{FieldName: "MatchLocalPreference", APIField: &routeMapClauseProps.MatchLocalPreference, TFValue: config.MatchLocalPreference, IsConfigured: configuredAttrs.IsConfigured("match_local_preference")},
		{FieldName: "MatchMetric", APIField: &routeMapClauseProps.MatchMetric, TFValue: config.MatchMetric, IsConfigured: configuredAttrs.IsConfigured("match_metric")},
		{FieldName: "MatchPeerInterface", APIField: &routeMapClauseProps.MatchPeerInterface, TFValue: config.MatchPeerInterface, IsConfigured: configuredAttrs.IsConfigured("match_peer_interface")},
		{FieldName: "MatchPeerVlan", APIField: &routeMapClauseProps.MatchPeerVlan, TFValue: config.MatchPeerVlan, IsConfigured: configuredAttrs.IsConfigured("match_peer_vlan")},
		{FieldName: "MatchTag", APIField: &routeMapClauseProps.MatchTag, TFValue: config.MatchTag, IsConfigured: configuredAttrs.IsConfigured("match_tag")},
		{FieldName: "MatchVni", APIField: &routeMapClauseProps.MatchVni, TFValue: config.MatchVni, IsConfigured: configuredAttrs.IsConfigured("match_vni")},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.RoutemapclausesPutRequestRouteMapClauseValueObjectProperties{}
		utils.SetObjectPropertiesFields([]utils.ObjectPropertiesField{
			{Name: "Notes", TFValue: op.Notes, APIValue: &objProps.Notes},
			{Name: "MatchFieldsShown", TFValue: op.MatchFieldsShown, APIValue: &objProps.MatchFieldsShown},
		})
		routeMapClauseProps.ObjectProperties = &objProps
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "route_map_clause", name, *routeMapClauseProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Route Map Clause %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "route_map_clauses")

	var minState verityRouteMapClauseResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if routeMapClauseData, exists := bulkMgr.GetResourceResponse("route_map_clause", name); exists {
			state := populateRouteMapClauseState(ctx, minState, routeMapClauseData, r.provCtx.mode)
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

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if routeMapClauseData, exists := r.bulkOpsMgr.GetResourceResponse("route_map_clause", name); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached route_map_clause data for %s from recent operation", name))
			state = populateRouteMapClauseState(ctx, state, routeMapClauseData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("route_map_clause") {
		tflog.Info(ctx, fmt.Sprintf("Skipping Route Map Clause %s verification – trusting recent successful API operation", name))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching Route Map Clause for verification of %s", name))

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

	state = populateRouteMapClauseState(ctx, state, routeMapClauseMap, r.provCtx.mode)
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

	// Get config for nullable field handling
	var config verityRouteMapClauseResourceModel
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
	routeMapClauseProps := openapi.RoutemapclausesPutRequestRouteMapClauseValue{}
	hasChanges := false

	// Parse HCL to detect which fields are explicitly configured
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, routeMapClauseTerraformType, name)

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

	// Handle nullable int64 field changes - parse HCL to detect explicit config
	utils.CompareAndSetNullableInt64Field(config.MatchInterfaceNumber, state.MatchInterfaceNumber, configuredAttrs.IsConfigured("match_interface_number"), func(v *openapi.NullableInt32) { routeMapClauseProps.MatchInterfaceNumber = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.MatchInterfaceVlan, state.MatchInterfaceVlan, configuredAttrs.IsConfigured("match_interface_vlan"), func(v *openapi.NullableInt32) { routeMapClauseProps.MatchInterfaceVlan = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.MatchLocalPreference, state.MatchLocalPreference, configuredAttrs.IsConfigured("match_local_preference"), func(v *openapi.NullableInt32) { routeMapClauseProps.MatchLocalPreference = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.MatchMetric, state.MatchMetric, configuredAttrs.IsConfigured("match_metric"), func(v *openapi.NullableInt32) { routeMapClauseProps.MatchMetric = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.MatchPeerInterface, state.MatchPeerInterface, configuredAttrs.IsConfigured("match_peer_interface"), func(v *openapi.NullableInt32) { routeMapClauseProps.MatchPeerInterface = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.MatchPeerVlan, state.MatchPeerVlan, configuredAttrs.IsConfigured("match_peer_vlan"), func(v *openapi.NullableInt32) { routeMapClauseProps.MatchPeerVlan = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.MatchTag, state.MatchTag, configuredAttrs.IsConfigured("match_tag"), func(v *openapi.NullableInt32) { routeMapClauseProps.MatchTag = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(config.MatchVni, state.MatchVni, configuredAttrs.IsConfigured("match_vni"), func(v *openapi.NullableInt32) { routeMapClauseProps.MatchVni = *v }, &hasChanges)

	// Handle object properties
	if len(plan.ObjectProperties) > 0 && len(state.ObjectProperties) > 0 {
		objProps := openapi.RoutemapclausesPutRequestRouteMapClauseValueObjectProperties{}
		op := plan.ObjectProperties[0]
		st := state.ObjectProperties[0]
		objPropsChanged := false

		utils.CompareAndSetObjectPropertiesFields([]utils.ObjectPropertiesFieldWithComparison{
			{Name: "Notes", PlanValue: op.Notes, StateValue: st.Notes, APIValue: &objProps.Notes},
			{Name: "MatchFieldsShown", PlanValue: op.MatchFieldsShown, StateValue: st.MatchFieldsShown, APIValue: &objProps.MatchFieldsShown},
		}, &objPropsChanged)

		if objPropsChanged {
			routeMapClauseProps.ObjectProperties = &objProps
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "route_map_clause", name, routeMapClauseProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Route Map Clause %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "route_map_clauses")

	var minState verityRouteMapClauseResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Try to use cached response from bulk operation to populate state with API values
	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if routeMapClauseData, exists := bulkMgr.GetResourceResponse("route_map_clause", name); exists {
			newState := populateRouteMapClauseState(ctx, minState, routeMapClauseData, r.provCtx.mode)
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "route_map_clause", name, nil, &resp.Diagnostics)
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

func populateRouteMapClauseState(ctx context.Context, state verityRouteMapClauseResourceModel, data map[string]interface{}, mode string) verityRouteMapClauseResourceModel {
	const resourceType = routeMapClauseResourceType

	state.Name = utils.MapStringFromAPI(data["name"])

	// String fields
	state.PermitDeny = utils.MapStringWithMode(data, "permit_deny", resourceType, mode)
	state.MatchAsPathAccessList = utils.MapStringWithMode(data, "match_as_path_access_list", resourceType, mode)
	state.MatchAsPathAccessListRefType = utils.MapStringWithMode(data, "match_as_path_access_list_ref_type_", resourceType, mode)
	state.MatchCommunityList = utils.MapStringWithMode(data, "match_community_list", resourceType, mode)
	state.MatchCommunityListRefType = utils.MapStringWithMode(data, "match_community_list_ref_type_", resourceType, mode)
	state.MatchExtendedCommunityList = utils.MapStringWithMode(data, "match_extended_community_list", resourceType, mode)
	state.MatchExtendedCommunityListRefType = utils.MapStringWithMode(data, "match_extended_community_list_ref_type_", resourceType, mode)
	state.MatchIpv4AddressIpPrefixList = utils.MapStringWithMode(data, "match_ipv4_address_ip_prefix_list", resourceType, mode)
	state.MatchIpv4AddressIpPrefixListRefType = utils.MapStringWithMode(data, "match_ipv4_address_ip_prefix_list_ref_type_", resourceType, mode)
	state.MatchIpv4NextHopIpPrefixList = utils.MapStringWithMode(data, "match_ipv4_next_hop_ip_prefix_list", resourceType, mode)
	state.MatchIpv4NextHopIpPrefixListRefType = utils.MapStringWithMode(data, "match_ipv4_next_hop_ip_prefix_list_ref_type_", resourceType, mode)
	state.MatchOrigin = utils.MapStringWithMode(data, "match_origin", resourceType, mode)
	state.MatchPeerIpAddress = utils.MapStringWithMode(data, "match_peer_ip_address", resourceType, mode)
	state.MatchSourceProtocol = utils.MapStringWithMode(data, "match_source_protocol", resourceType, mode)
	state.MatchVrf = utils.MapStringWithMode(data, "match_vrf", resourceType, mode)
	state.MatchVrfRefType = utils.MapStringWithMode(data, "match_vrf_ref_type_", resourceType, mode)
	state.MatchEvpnRouteType = utils.MapStringWithMode(data, "match_evpn_route_type", resourceType, mode)
	state.MatchIpv6AddressIpv6PrefixList = utils.MapStringWithMode(data, "match_ipv6_address_ipv6_prefix_list", resourceType, mode)
	state.MatchIpv6AddressIpv6PrefixListRefType = utils.MapStringWithMode(data, "match_ipv6_address_ipv6_prefix_list_ref_type_", resourceType, mode)
	state.MatchIpv6NextHopIpv6PrefixList = utils.MapStringWithMode(data, "match_ipv6_next_hop_ipv6_prefix_list", resourceType, mode)
	state.MatchIpv6NextHopIpv6PrefixListRefType = utils.MapStringWithMode(data, "match_ipv6_next_hop_ipv6_prefix_list_ref_type_", resourceType, mode)

	// Boolean fields
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)
	state.MatchEvpnRouteTypeDefault = utils.MapBoolWithMode(data, "match_evpn_route_type_default", resourceType, mode)

	// Int64 fields
	state.MatchInterfaceNumber = utils.MapInt64WithMode(data, "match_interface_number", resourceType, mode)
	state.MatchInterfaceVlan = utils.MapInt64WithMode(data, "match_interface_vlan", resourceType, mode)
	state.MatchLocalPreference = utils.MapInt64WithMode(data, "match_local_preference", resourceType, mode)
	state.MatchMetric = utils.MapInt64WithMode(data, "match_metric", resourceType, mode)
	state.MatchPeerInterface = utils.MapInt64WithMode(data, "match_peer_interface", resourceType, mode)
	state.MatchPeerVlan = utils.MapInt64WithMode(data, "match_peer_vlan", resourceType, mode)
	state.MatchTag = utils.MapInt64WithMode(data, "match_tag", resourceType, mode)
	state.MatchVni = utils.MapInt64WithMode(data, "match_vni", resourceType, mode)

	// Handle object_properties block
	if utils.FieldAppliesToMode(resourceType, "object_properties", mode) {
		if objProps, ok := data["object_properties"].(map[string]interface{}); ok {
			objPropsModel := verityRouteMapClauseObjectPropertiesModel{
				Notes:            utils.MapStringWithModeNested(objProps, "notes", resourceType, "object_properties.notes", mode),
				MatchFieldsShown: utils.MapStringWithModeNested(objProps, "match_fields_shown", resourceType, "object_properties.match_fields_shown", mode),
			}
			state.ObjectProperties = []verityRouteMapClauseObjectPropertiesModel{objPropsModel}
		} else {
			state.ObjectProperties = nil
		}
	} else {
		state.ObjectProperties = nil
	}

	return state
}

func (r *verityRouteMapClauseResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// =========================================================================
	// Skip if deleting
	// =========================================================================
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityRouteMapClauseResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Mode-aware field nullification
	// Set fields that don't apply to current mode to null to prevent
	// "known after apply" messages for irrelevant fields.
	// =========================================================================
	const resourceType = routeMapClauseResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyBools(
		"enable",
		"match_evpn_route_type_default",
	)

	nullifier.NullifyStrings(
		"permit_deny",
		"match_as_path_access_list",
		"match_as_path_access_list_ref_type_",
		"match_community_list",
		"match_community_list_ref_type_",
		"match_extended_community_list",
		"match_extended_community_list_ref_type_",
		"match_ipv4_address_ip_prefix_list",
		"match_ipv4_address_ip_prefix_list_ref_type_",
		"match_ipv4_next_hop_ip_prefix_list",
		"match_ipv4_next_hop_ip_prefix_list_ref_type_",
		"match_origin",
		"match_peer_ip_address",
		"match_source_protocol",
		"match_vrf",
		"match_vrf_ref_type_",
		"match_evpn_route_type",
		"match_ipv6_address_ipv6_prefix_list",
		"match_ipv6_address_ipv6_prefix_list_ref_type_",
		"match_ipv6_next_hop_ipv6_prefix_list",
		"match_ipv6_next_hop_ipv6_prefix_list_ref_type_",
	)

	nullifier.NullifyInt64s(
		"match_interface_number",
		"match_interface_vlan",
		"match_local_preference",
		"match_metric",
		"match_peer_interface",
		"match_peer_vlan",
		"match_tag",
		"match_vni",
	)

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "object_properties",
		ItemCount:    len(plan.ObjectProperties),
		StringFields: []string{"notes", "match_fields_shown"},
	})

	// =========================================================================
	// Skip UPDATE-specific logic during CREATE
	// =========================================================================
	if req.State.Raw.IsNull() {
		return
	}

	// =========================================================================
	// UPDATE operation - get state and config
	// =========================================================================
	var state verityRouteMapClauseResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config verityRouteMapClauseResourceModel
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
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, routeMapClauseTerraformType, name)

	utils.HandleNullableFields(utils.NullableFieldsConfig{
		Ctx:             ctx,
		Plan:            &resp.Plan,
		ConfiguredAttrs: configuredAttrs,
		Int64Fields: []utils.NullableInt64Field{
			{AttrName: "match_interface_number", ConfigVal: config.MatchInterfaceNumber, StateVal: state.MatchInterfaceNumber},
			{AttrName: "match_interface_vlan", ConfigVal: config.MatchInterfaceVlan, StateVal: state.MatchInterfaceVlan},
			{AttrName: "match_local_preference", ConfigVal: config.MatchLocalPreference, StateVal: state.MatchLocalPreference},
			{AttrName: "match_metric", ConfigVal: config.MatchMetric, StateVal: state.MatchMetric},
			{AttrName: "match_peer_interface", ConfigVal: config.MatchPeerInterface, StateVal: state.MatchPeerInterface},
			{AttrName: "match_peer_vlan", ConfigVal: config.MatchPeerVlan, StateVal: state.MatchPeerVlan},
			{AttrName: "match_tag", ConfigVal: config.MatchTag, StateVal: state.MatchTag},
			{AttrName: "match_vni", ConfigVal: config.MatchVni, StateVal: state.MatchVni},
		},
	})
}
