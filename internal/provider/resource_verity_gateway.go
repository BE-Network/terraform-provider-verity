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
	_ resource.Resource                = &verityGatewayResource{}
	_ resource.ResourceWithConfigure   = &verityGatewayResource{}
	_ resource.ResourceWithImportState = &verityGatewayResource{}
)

func NewVerityGatewayResource() resource.Resource {
	return &verityGatewayResource{}
}

type verityGatewayResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type verityGatewayResourceModel struct {
	Name                    types.String                         `tfsdk:"name"`
	Enable                  types.Bool                           `tfsdk:"enable"`
	ObjectProperties        []verityGatewayObjectPropertiesModel `tfsdk:"object_properties"`
	Tenant                  types.String                         `tfsdk:"tenant"`
	TenantRefType           types.String                         `tfsdk:"tenant_ref_type_"`
	NeighborIpAddress       types.String                         `tfsdk:"neighbor_ip_address"`
	NeighborAsNumber        types.Int64                          `tfsdk:"neighbor_as_number"`
	FabricInterconnect      types.Bool                           `tfsdk:"fabric_interconnect"`
	KeepaliveTimer          types.Int64                          `tfsdk:"keepalive_timer"`
	HoldTimer               types.Int64                          `tfsdk:"hold_timer"`
	ConnectTimer            types.Int64                          `tfsdk:"connect_timer"`
	AdvertisementInterval   types.Int64                          `tfsdk:"advertisement_interval"`
	EbgpMultihop            types.Int64                          `tfsdk:"ebgp_multihop"`
	EgressVlan              types.Int64                          `tfsdk:"egress_vlan"`
	SourceIpAddress         types.String                         `tfsdk:"source_ip_address"`
	AnycastIpMask           types.String                         `tfsdk:"anycast_ip_mask"`
	Md5Password             types.String                         `tfsdk:"md5_password"`
	ImportRouteMap          types.String                         `tfsdk:"import_route_map"`
	StaticRoutes            []verityGatewayStaticRoutesModel     `tfsdk:"static_routes"`
	ExportRouteMap          types.String                         `tfsdk:"export_route_map"`
	GatewayMode             types.String                         `tfsdk:"gateway_mode"`
	LocalAsNumber           types.Int64                          `tfsdk:"local_as_number"`
	LocalAsNoPrepend        types.Bool                           `tfsdk:"local_as_no_prepend"`
	ReplaceAs               types.Bool                           `tfsdk:"replace_as"`
	MaxLocalAsOccurrences   types.Int64                          `tfsdk:"max_local_as_occurrences"`
	DynamicBgpSubnet        types.String                         `tfsdk:"dynamic_bgp_subnet"`
	DynamicBgpLimits        types.Int64                          `tfsdk:"dynamic_bgp_limits"`
	HelperHopIpAddress      types.String                         `tfsdk:"helper_hop_ip_address"`
	EnableBfd               types.Bool                           `tfsdk:"enable_bfd"`
	BfdReceiveInterval      types.Int64                          `tfsdk:"bfd_receive_interval"`
	BfdTransmissionInterval types.Int64                          `tfsdk:"bfd_transmission_interval"`
	BfdDetectMultiplier     types.Int64                          `tfsdk:"bfd_detect_multiplier"`
	BfdMultihop             types.Bool                           `tfsdk:"bfd_multihop"`
	NextHopSelf             types.Bool                           `tfsdk:"next_hop_self"`
	DefaultOriginate        types.Bool                           `tfsdk:"default_originate"`
	ExportRouteMapRefType   types.String                         `tfsdk:"export_route_map_ref_type_"`
	ImportRouteMapRefType   types.String                         `tfsdk:"import_route_map_ref_type_"`
}

type verityGatewayObjectPropertiesModel struct {
	Group types.String `tfsdk:"group"`
}

type verityGatewayStaticRoutesModel struct {
	Enable           types.Bool   `tfsdk:"enable"`
	Ipv4RoutePrefix  types.String `tfsdk:"ipv4_route_prefix"`
	NextHopIpAddress types.String `tfsdk:"next_hop_ip_address"`
	AdValue          types.Int64  `tfsdk:"ad_value"`
	Index            types.Int64  `tfsdk:"index"`
}

func (sr verityGatewayStaticRoutesModel) GetIndex() types.Int64 {
	return sr.Index
}

func (r *verityGatewayResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_gateway"
}

func (r *verityGatewayResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityGatewayResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity Gateway",
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
			"tenant": schema.StringAttribute{
				Description: "Tenant",
				Optional:    true,
			},
			"tenant_ref_type_": schema.StringAttribute{
				Description: "Object type for tenant field",
				Optional:    true,
			},
			"neighbor_ip_address": schema.StringAttribute{
				Description: "IP address of remote BGP peer",
				Optional:    true,
			},
			"neighbor_as_number": schema.Int64Attribute{
				Description: "Autonomous System Number of remote BGP peer",
				Optional:    true,
			},
			"fabric_interconnect": schema.BoolAttribute{
				Optional: true,
			},
			"keepalive_timer": schema.Int64Attribute{
				Description: "Interval in seconds between Keepalive messages",
				Optional:    true,
			},
			"hold_timer": schema.Int64Attribute{
				Optional: true,
			},
			"connect_timer": schema.Int64Attribute{
				Optional: true,
			},
			"advertisement_interval": schema.Int64Attribute{
				Optional: true,
			},
			"ebgp_multihop": schema.Int64Attribute{
				Optional: true,
			},
			"egress_vlan": schema.Int64Attribute{
				Optional: true,
			},
			"source_ip_address": schema.StringAttribute{
				Optional: true,
			},
			"anycast_ip_mask": schema.StringAttribute{
				Optional: true,
			},
			"md5_password": schema.StringAttribute{
				Optional: true,
			},
			"import_route_map": schema.StringAttribute{
				Optional: true,
			},
			"export_route_map": schema.StringAttribute{
				Optional: true,
			},
			"gateway_mode": schema.StringAttribute{
				Optional: true,
			},
			"local_as_number": schema.Int64Attribute{
				Optional: true,
			},
			"local_as_no_prepend": schema.BoolAttribute{
				Optional: true,
			},
			"replace_as": schema.BoolAttribute{
				Optional: true,
			},
			"max_local_as_occurrences": schema.Int64Attribute{
				Optional: true,
			},
			"dynamic_bgp_subnet": schema.StringAttribute{
				Optional: true,
			},
			"dynamic_bgp_limits": schema.Int64Attribute{
				Optional: true,
			},
			"helper_hop_ip_address": schema.StringAttribute{
				Optional: true,
			},
			"enable_bfd": schema.BoolAttribute{
				Optional: true,
			},
			"bfd_receive_interval": schema.Int64Attribute{
				Optional: true,
			},
			"bfd_transmission_interval": schema.Int64Attribute{
				Optional: true,
			},
			"bfd_detect_multiplier": schema.Int64Attribute{
				Optional: true,
			},
			"bfd_multihop": schema.BoolAttribute{
				Optional: true,
			},
			"next_hop_self": schema.BoolAttribute{
				Optional: true,
			},
			"default_originate": schema.BoolAttribute{
				Description: "Instructs BGP to generate and send a default route 0.0.0.0/0 to the specified neighbor.",
				Optional:    true,
			},
			"export_route_map_ref_type_": schema.StringAttribute{
				Optional: true,
			},
			"import_route_map_ref_type_": schema.StringAttribute{
				Optional: true,
			},
		},
		Blocks: map[string]schema.Block{
			"object_properties": schema.ListNestedBlock{
				Description: "Object properties for the gateway",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"group": schema.StringAttribute{
							Description: "Group",
							Optional:    true,
						},
					},
				},
			},
			"static_routes": schema.ListNestedBlock{
				Description: "List of static routes",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enable": schema.BoolAttribute{
							Description: "Enable of this static route",
							Optional:    true,
						},
						"ipv4_route_prefix": schema.StringAttribute{
							Description: "IPv4 unicast IP address with subnet mask",
							Optional:    true,
						},
						"next_hop_ip_address": schema.StringAttribute{
							Description: "Next Hop IP Address",
							Optional:    true,
						},
						"ad_value": schema.Int64Attribute{
							Description: "Administrative distance value (0-255)",
							Optional:    true,
						},
						"index": schema.Int64Attribute{
							Description: "Index identifying the object",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func (r *verityGatewayResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityGatewayResourceModel
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
	gatewayProps := &openapi.GatewaysPutRequestGatewayValue{
		Name: openapi.PtrString(name),
	}

	// Handle string fields
	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "Tenant", APIField: &gatewayProps.Tenant, TFValue: plan.Tenant},
		{FieldName: "TenantRefType", APIField: &gatewayProps.TenantRefType, TFValue: plan.TenantRefType},
		{FieldName: "NeighborIpAddress", APIField: &gatewayProps.NeighborIpAddress, TFValue: plan.NeighborIpAddress},
		{FieldName: "SourceIpAddress", APIField: &gatewayProps.SourceIpAddress, TFValue: plan.SourceIpAddress},
		{FieldName: "AnycastIpMask", APIField: &gatewayProps.AnycastIpMask, TFValue: plan.AnycastIpMask},
		{FieldName: "Md5Password", APIField: &gatewayProps.Md5Password, TFValue: plan.Md5Password},
		{FieldName: "ImportRouteMap", APIField: &gatewayProps.ImportRouteMap, TFValue: plan.ImportRouteMap},
		{FieldName: "ExportRouteMap", APIField: &gatewayProps.ExportRouteMap, TFValue: plan.ExportRouteMap},
		{FieldName: "GatewayMode", APIField: &gatewayProps.GatewayMode, TFValue: plan.GatewayMode},
		{FieldName: "DynamicBgpSubnet", APIField: &gatewayProps.DynamicBgpSubnet, TFValue: plan.DynamicBgpSubnet},
		{FieldName: "HelperHopIpAddress", APIField: &gatewayProps.HelperHopIpAddress, TFValue: plan.HelperHopIpAddress},
		{FieldName: "ExportRouteMapRefType", APIField: &gatewayProps.ExportRouteMapRefType, TFValue: plan.ExportRouteMapRefType},
		{FieldName: "ImportRouteMapRefType", APIField: &gatewayProps.ImportRouteMapRefType, TFValue: plan.ImportRouteMapRefType},
	})

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &gatewayProps.Enable, TFValue: plan.Enable},
		{FieldName: "FabricInterconnect", APIField: &gatewayProps.FabricInterconnect, TFValue: plan.FabricInterconnect},
		{FieldName: "LocalAsNoPrepend", APIField: &gatewayProps.LocalAsNoPrepend, TFValue: plan.LocalAsNoPrepend},
		{FieldName: "ReplaceAs", APIField: &gatewayProps.ReplaceAs, TFValue: plan.ReplaceAs},
		{FieldName: "EnableBfd", APIField: &gatewayProps.EnableBfd, TFValue: plan.EnableBfd},
		{FieldName: "BfdMultihop", APIField: &gatewayProps.BfdMultihop, TFValue: plan.BfdMultihop},
		{FieldName: "NextHopSelf", APIField: &gatewayProps.NextHopSelf, TFValue: plan.NextHopSelf},
		{FieldName: "DefaultOriginate", APIField: &gatewayProps.DefaultOriginate, TFValue: plan.DefaultOriginate},
	})

	// Handle nullable int64 fields
	utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
		{FieldName: "NeighborAsNumber", APIField: &gatewayProps.NeighborAsNumber, TFValue: plan.NeighborAsNumber},
		{FieldName: "KeepaliveTimer", APIField: &gatewayProps.KeepaliveTimer, TFValue: plan.KeepaliveTimer},
		{FieldName: "HoldTimer", APIField: &gatewayProps.HoldTimer, TFValue: plan.HoldTimer},
		{FieldName: "ConnectTimer", APIField: &gatewayProps.ConnectTimer, TFValue: plan.ConnectTimer},
		{FieldName: "AdvertisementInterval", APIField: &gatewayProps.AdvertisementInterval, TFValue: plan.AdvertisementInterval},
		{FieldName: "EbgpMultihop", APIField: &gatewayProps.EbgpMultihop, TFValue: plan.EbgpMultihop},
		{FieldName: "EgressVlan", APIField: &gatewayProps.EgressVlan, TFValue: plan.EgressVlan},
		{FieldName: "LocalAsNumber", APIField: &gatewayProps.LocalAsNumber, TFValue: plan.LocalAsNumber},
		{FieldName: "MaxLocalAsOccurrences", APIField: &gatewayProps.MaxLocalAsOccurrences, TFValue: plan.MaxLocalAsOccurrences},
		{FieldName: "DynamicBgpLimits", APIField: &gatewayProps.DynamicBgpLimits, TFValue: plan.DynamicBgpLimits},
		{FieldName: "BfdReceiveInterval", APIField: &gatewayProps.BfdReceiveInterval, TFValue: plan.BfdReceiveInterval},
		{FieldName: "BfdTransmissionInterval", APIField: &gatewayProps.BfdTransmissionInterval, TFValue: plan.BfdTransmissionInterval},
		{FieldName: "BfdDetectMultiplier", APIField: &gatewayProps.BfdDetectMultiplier, TFValue: plan.BfdDetectMultiplier},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.EthportsettingsPutRequestEthPortSettingsValueObjectProperties{}
		utils.SetObjectPropertiesFields([]utils.ObjectPropertiesField{
			{Name: "Group", TFValue: op.Group, APIValue: &objProps.Group},
		})
		gatewayProps.ObjectProperties = &objProps
	}

	// Handle static routes
	if len(plan.StaticRoutes) > 0 {
		routes := make([]openapi.GatewaysPutRequestGatewayValueStaticRoutesInner, len(plan.StaticRoutes))
		for i, item := range plan.StaticRoutes {
			rItem := openapi.GatewaysPutRequestGatewayValueStaticRoutesInner{}
			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &rItem.Enable, TFValue: item.Enable},
			})
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Ipv4RoutePrefix", APIField: &rItem.Ipv4RoutePrefix, TFValue: item.Ipv4RoutePrefix},
				{FieldName: "NextHopIpAddress", APIField: &rItem.NextHopIpAddress, TFValue: item.NextHopIpAddress},
			})
			utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
				{FieldName: "AdValue", APIField: &rItem.AdValue, TFValue: item.AdValue},
			})
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &rItem.Index, TFValue: item.Index},
			})
			routes[i] = rItem
		}
		gatewayProps.StaticRoutes = routes
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "gateway", name, *gatewayProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Gateway %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "gateways")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
}

func (r *verityGatewayResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityGatewayResourceModel
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

	gatewayName := state.Name.ValueString()

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("gateway") {
		tflog.Info(ctx, fmt.Sprintf("Skipping gateway %s verification – trusting recent successful API operation", gatewayName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching gateways for verification of %s", gatewayName))

	type GatewaysResponse struct {
		Gateway map[string]interface{} `json:"gateway"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "gateways", gatewayName,
		func() (GatewaysResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch gateways")
			respAPI, err := r.client.GatewaysAPI.GatewaysGet(ctx).Execute()
			if err != nil {
				return GatewaysResponse{}, fmt.Errorf("error reading gateways: %v", err)
			}
			defer respAPI.Body.Close()

			var res GatewaysResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return GatewaysResponse{}, fmt.Errorf("failed to decode gateways response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d gateways", len(res.Gateway)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Gateway %s", gatewayName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for gateway with name: %s", gatewayName))

	gatewayData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.Gateway,
		gatewayName,
		func(data interface{}) (string, bool) {
			if gateway, ok := data.(map[string]interface{}); ok {
				if name, ok := gateway["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Gateway with name '%s' not found in API response", gatewayName))
		resp.State.RemoveResource(ctx)
		return
	}

	gatewayMap, ok := gatewayData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Gateway Data",
			fmt.Sprintf("Gateway data is not in expected format for %s", gatewayName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found gateway '%s' under API key '%s'", gatewayName, actualAPIName))

	state.Name = utils.MapStringFromAPI(gatewayMap["name"])

	// Handle object properties
	if objProps, ok := gatewayMap["object_properties"].(map[string]interface{}); ok {
		state.ObjectProperties = []verityGatewayObjectPropertiesModel{
			{Group: utils.MapStringFromAPI(objProps["group"])},
		}
	} else {
		state.ObjectProperties = nil
	}

	// Map string fields
	stringFieldMappings := map[string]*types.String{
		"tenant":                     &state.Tenant,
		"tenant_ref_type_":           &state.TenantRefType,
		"neighbor_ip_address":        &state.NeighborIpAddress,
		"source_ip_address":          &state.SourceIpAddress,
		"anycast_ip_mask":            &state.AnycastIpMask,
		"md5_password":               &state.Md5Password,
		"import_route_map":           &state.ImportRouteMap,
		"export_route_map":           &state.ExportRouteMap,
		"gateway_mode":               &state.GatewayMode,
		"dynamic_bgp_subnet":         &state.DynamicBgpSubnet,
		"helper_hop_ip_address":      &state.HelperHopIpAddress,
		"export_route_map_ref_type_": &state.ExportRouteMapRefType,
		"import_route_map_ref_type_": &state.ImportRouteMapRefType,
	}

	for apiKey, stateField := range stringFieldMappings {
		*stateField = utils.MapStringFromAPI(gatewayMap[apiKey])
	}

	// Map boolean fields
	boolFieldMappings := map[string]*types.Bool{
		"enable":              &state.Enable,
		"fabric_interconnect": &state.FabricInterconnect,
		"local_as_no_prepend": &state.LocalAsNoPrepend,
		"replace_as":          &state.ReplaceAs,
		"enable_bfd":          &state.EnableBfd,
		"bfd_multihop":        &state.BfdMultihop,
		"next_hop_self":       &state.NextHopSelf,
		"default_originate":   &state.DefaultOriginate,
	}

	for apiKey, stateField := range boolFieldMappings {
		*stateField = utils.MapBoolFromAPI(gatewayMap[apiKey])
	}

	// Map int64 fields
	int64FieldMappings := map[string]*types.Int64{
		"neighbor_as_number":        &state.NeighborAsNumber,
		"keepalive_timer":           &state.KeepaliveTimer,
		"hold_timer":                &state.HoldTimer,
		"connect_timer":             &state.ConnectTimer,
		"advertisement_interval":    &state.AdvertisementInterval,
		"ebgp_multihop":             &state.EbgpMultihop,
		"egress_vlan":               &state.EgressVlan,
		"local_as_number":           &state.LocalAsNumber,
		"max_local_as_occurrences":  &state.MaxLocalAsOccurrences,
		"dynamic_bgp_limits":        &state.DynamicBgpLimits,
		"bfd_receive_interval":      &state.BfdReceiveInterval,
		"bfd_transmission_interval": &state.BfdTransmissionInterval,
		"bfd_detect_multiplier":     &state.BfdDetectMultiplier,
	}

	for apiKey, stateField := range int64FieldMappings {
		*stateField = utils.MapInt64FromAPI(gatewayMap[apiKey])
	}

	// Handle static routes
	if routes, ok := gatewayMap["static_routes"].([]interface{}); ok && len(routes) > 0 {
		var staticRoutes []verityGatewayStaticRoutesModel

		for _, r := range routes {
			route, ok := r.(map[string]interface{})
			if !ok {
				continue
			}

			srModel := verityGatewayStaticRoutesModel{
				Enable:           utils.MapBoolFromAPI(route["enable"]),
				Ipv4RoutePrefix:  utils.MapStringFromAPI(route["ipv4_route_prefix"]),
				NextHopIpAddress: utils.MapStringFromAPI(route["next_hop_ip_address"]),
				AdValue:          utils.MapInt64FromAPI(route["ad_value"]),
				Index:            utils.MapInt64FromAPI(route["index"]),
			}

			staticRoutes = append(staticRoutes, srModel)
		}

		state.StaticRoutes = staticRoutes
	} else {
		state.StaticRoutes = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityGatewayResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityGatewayResourceModel

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
	gatewayProps := openapi.GatewaysPutRequestGatewayValue{}
	hasChanges := false

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { gatewayProps.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.NeighborIpAddress, state.NeighborIpAddress, func(v *string) { gatewayProps.NeighborIpAddress = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.SourceIpAddress, state.SourceIpAddress, func(v *string) { gatewayProps.SourceIpAddress = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.AnycastIpMask, state.AnycastIpMask, func(v *string) { gatewayProps.AnycastIpMask = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Md5Password, state.Md5Password, func(v *string) { gatewayProps.Md5Password = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.GatewayMode, state.GatewayMode, func(v *string) { gatewayProps.GatewayMode = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.DynamicBgpSubnet, state.DynamicBgpSubnet, func(v *string) { gatewayProps.DynamicBgpSubnet = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.HelperHopIpAddress, state.HelperHopIpAddress, func(v *string) { gatewayProps.HelperHopIpAddress = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { gatewayProps.Enable = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.FabricInterconnect, state.FabricInterconnect, func(v *bool) { gatewayProps.FabricInterconnect = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.LocalAsNoPrepend, state.LocalAsNoPrepend, func(v *bool) { gatewayProps.LocalAsNoPrepend = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.ReplaceAs, state.ReplaceAs, func(v *bool) { gatewayProps.ReplaceAs = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.EnableBfd, state.EnableBfd, func(v *bool) { gatewayProps.EnableBfd = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.BfdMultihop, state.BfdMultihop, func(v *bool) { gatewayProps.BfdMultihop = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.NextHopSelf, state.NextHopSelf, func(v *bool) { gatewayProps.NextHopSelf = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.DefaultOriginate, state.DefaultOriginate, func(v *bool) { gatewayProps.DefaultOriginate = v }, &hasChanges)

	// Handle nullable int64 field changes
	utils.CompareAndSetNullableInt64Field(plan.NeighborAsNumber, state.NeighborAsNumber, func(v *openapi.NullableInt32) { gatewayProps.NeighborAsNumber = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.KeepaliveTimer, state.KeepaliveTimer, func(v *openapi.NullableInt32) { gatewayProps.KeepaliveTimer = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.HoldTimer, state.HoldTimer, func(v *openapi.NullableInt32) { gatewayProps.HoldTimer = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.ConnectTimer, state.ConnectTimer, func(v *openapi.NullableInt32) { gatewayProps.ConnectTimer = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.AdvertisementInterval, state.AdvertisementInterval, func(v *openapi.NullableInt32) { gatewayProps.AdvertisementInterval = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.EbgpMultihop, state.EbgpMultihop, func(v *openapi.NullableInt32) { gatewayProps.EbgpMultihop = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.EgressVlan, state.EgressVlan, func(v *openapi.NullableInt32) { gatewayProps.EgressVlan = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.LocalAsNumber, state.LocalAsNumber, func(v *openapi.NullableInt32) { gatewayProps.LocalAsNumber = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.MaxLocalAsOccurrences, state.MaxLocalAsOccurrences, func(v *openapi.NullableInt32) { gatewayProps.MaxLocalAsOccurrences = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.DynamicBgpLimits, state.DynamicBgpLimits, func(v *openapi.NullableInt32) { gatewayProps.DynamicBgpLimits = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.BfdReceiveInterval, state.BfdReceiveInterval, func(v *openapi.NullableInt32) { gatewayProps.BfdReceiveInterval = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.BfdTransmissionInterval, state.BfdTransmissionInterval, func(v *openapi.NullableInt32) { gatewayProps.BfdTransmissionInterval = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.BfdDetectMultiplier, state.BfdDetectMultiplier, func(v *openapi.NullableInt32) { gatewayProps.BfdDetectMultiplier = *v }, &hasChanges)

	// Handle object properties
	if len(plan.ObjectProperties) > 0 && len(state.ObjectProperties) > 0 {
		objProps := openapi.EthportsettingsPutRequestEthPortSettingsValueObjectProperties{}
		op := plan.ObjectProperties[0]
		st := state.ObjectProperties[0]
		objPropsChanged := false

		utils.CompareAndSetObjectPropertiesFields([]utils.ObjectPropertiesFieldWithComparison{
			{Name: "Group", PlanValue: op.Group, StateValue: st.Group, APIValue: &objProps.Group},
		}, &objPropsChanged)

		if objPropsChanged {
			gatewayProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	// Handle tenant and tenant_ref_type_ fields using "Many ref types supported" pattern
	if !utils.HandleMultipleRefTypesSupported(
		plan.Tenant, state.Tenant, plan.TenantRefType, state.TenantRefType,
		func(v *string) { gatewayProps.Tenant = v },
		func(v *string) { gatewayProps.TenantRefType = v },
		"tenant", "tenant_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	// Handle ImportRouteMap and ImportRouteMapRefType using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.ImportRouteMap, state.ImportRouteMap, plan.ImportRouteMapRefType, state.ImportRouteMapRefType,
		func(v *string) { gatewayProps.ImportRouteMap = v },
		func(v *string) { gatewayProps.ImportRouteMapRefType = v },
		"import_route_map", "import_route_map_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	// Handle ExportRouteMap and ExportRouteMapRefType using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.ExportRouteMap, state.ExportRouteMap, plan.ExportRouteMapRefType, state.ExportRouteMapRefType,
		func(v *string) { gatewayProps.ExportRouteMap = v },
		func(v *string) { gatewayProps.ExportRouteMapRefType = v },
		"export_route_map", "export_route_map_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	// Handle static routes
	staticRoutesHandler := utils.IndexedItemHandler[verityGatewayStaticRoutesModel, openapi.GatewaysPutRequestGatewayValueStaticRoutesInner]{
		CreateNew: func(planItem verityGatewayStaticRoutesModel) openapi.GatewaysPutRequestGatewayValueStaticRoutesInner {
			route := openapi.GatewaysPutRequestGatewayValueStaticRoutesInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &route.Index, TFValue: planItem.Index},
			})

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "Enable", APIField: &route.Enable, TFValue: planItem.Enable},
			})

			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "Ipv4RoutePrefix", APIField: &route.Ipv4RoutePrefix, TFValue: planItem.Ipv4RoutePrefix},
				{FieldName: "NextHopIpAddress", APIField: &route.NextHopIpAddress, TFValue: planItem.NextHopIpAddress},
			})

			utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
				{FieldName: "AdValue", APIField: &route.AdValue, TFValue: planItem.AdValue},
			})

			return route
		},
		UpdateExisting: func(planItem verityGatewayStaticRoutesModel, stateItem verityGatewayStaticRoutesModel) (openapi.GatewaysPutRequestGatewayValueStaticRoutesInner, bool) {
			route := openapi.GatewaysPutRequestGatewayValueStaticRoutesInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &route.Index, TFValue: planItem.Index},
			})

			fieldChanged := false

			// Handle boolean fields
			utils.CompareAndSetBoolField(planItem.Enable, stateItem.Enable, func(v *bool) { route.Enable = v }, &fieldChanged)

			// Handle string fields
			utils.CompareAndSetStringField(planItem.Ipv4RoutePrefix, stateItem.Ipv4RoutePrefix, func(v *string) { route.Ipv4RoutePrefix = v }, &fieldChanged)
			utils.CompareAndSetStringField(planItem.NextHopIpAddress, stateItem.NextHopIpAddress, func(v *string) { route.NextHopIpAddress = v }, &fieldChanged)

			// Handle nullable int64 fields
			utils.CompareAndSetNullableInt64Field(planItem.AdValue, stateItem.AdValue, func(v *openapi.NullableInt32) { route.AdValue = *v }, &fieldChanged)

			return route, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.GatewaysPutRequestGatewayValueStaticRoutesInner {
			route := openapi.GatewaysPutRequestGatewayValueStaticRoutesInner{}
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &route.Index, TFValue: types.Int64Value(index)},
			})
			return route
		},
	}

	changedStaticRoutes, staticRoutesChanged := utils.ProcessIndexedArrayUpdates(plan.StaticRoutes, state.StaticRoutes, staticRoutesHandler)
	if staticRoutesChanged {
		gatewayProps.StaticRoutes = changedStaticRoutes
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "gateway", name, gatewayProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Gateway %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "gateways")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *verityGatewayResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityGatewayResourceModel
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "gateway", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Gateway %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "gateways")
	resp.State.RemoveResource(ctx)
}

func (r *verityGatewayResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
