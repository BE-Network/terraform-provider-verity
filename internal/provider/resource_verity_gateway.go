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
	_ resource.ResourceWithModifyPlan  = &verityGatewayResource{}
)

const gatewayResourceType = "gateways"

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
	Name                       types.String                         `tfsdk:"name"`
	Enable                     types.Bool                           `tfsdk:"enable"`
	ObjectProperties           []verityGatewayObjectPropertiesModel `tfsdk:"object_properties"`
	Tenant                     types.String                         `tfsdk:"tenant"`
	TenantRefType              types.String                         `tfsdk:"tenant_ref_type_"`
	NeighborIpAddress          types.String                         `tfsdk:"neighbor_ip_address"`
	NeighborAsNumber           types.Int64                          `tfsdk:"neighbor_as_number"`
	FabricInterconnect         types.Bool                           `tfsdk:"fabric_interconnect"`
	KeepaliveTimer             types.Int64                          `tfsdk:"keepalive_timer"`
	HoldTimer                  types.Int64                          `tfsdk:"hold_timer"`
	ConnectTimer               types.Int64                          `tfsdk:"connect_timer"`
	AdvertisementInterval      types.Int64                          `tfsdk:"advertisement_interval"`
	EbgpMultihop               types.Int64                          `tfsdk:"ebgp_multihop"`
	EgressVlan                 types.Int64                          `tfsdk:"egress_vlan"`
	SourceIpAddress            types.String                         `tfsdk:"source_ip_address"`
	AnycastIpMask              types.String                         `tfsdk:"anycast_ip_mask"`
	Md5Password                types.String                         `tfsdk:"md5_password"`
	Md5PasswordEncrypted       types.String                         `tfsdk:"md5_password_encrypted"`
	SwitchEncryptedMd5Password types.Bool                           `tfsdk:"switch_encrypted_md5_password"`
	ImportRouteMap             types.String                         `tfsdk:"import_route_map"`
	StaticRoutes               []verityGatewayStaticRoutesModel     `tfsdk:"static_routes"`
	ExportRouteMap             types.String                         `tfsdk:"export_route_map"`
	GatewayMode                types.String                         `tfsdk:"gateway_mode"`
	LocalAsNumber              types.Int64                          `tfsdk:"local_as_number"`
	LocalAsNoPrepend           types.Bool                           `tfsdk:"local_as_no_prepend"`
	ReplaceAs                  types.Bool                           `tfsdk:"replace_as"`
	MaxLocalAsOccurrences      types.Int64                          `tfsdk:"max_local_as_occurrences"`
	DynamicBgpSubnet           types.String                         `tfsdk:"dynamic_bgp_subnet"`
	DynamicBgpLimits           types.Int64                          `tfsdk:"dynamic_bgp_limits"`
	HelperHopIpAddress         types.String                         `tfsdk:"helper_hop_ip_address"`
	EnableBfd                  types.Bool                           `tfsdk:"enable_bfd"`
	BfdReceiveInterval         types.Int64                          `tfsdk:"bfd_receive_interval"`
	BfdTransmissionInterval    types.Int64                          `tfsdk:"bfd_transmission_interval"`
	BfdDetectMultiplier        types.Int64                          `tfsdk:"bfd_detect_multiplier"`
	BfdMultihop                types.Bool                           `tfsdk:"bfd_multihop"`
	NextHopSelf                types.Bool                           `tfsdk:"next_hop_self"`
	DefaultOriginate           types.Bool                           `tfsdk:"default_originate"`
	ExportRouteMapRefType      types.String                         `tfsdk:"export_route_map_ref_type_"`
	ImportRouteMapRefType      types.String                         `tfsdk:"import_route_map_ref_type_"`
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
				Computed:    true,
			},
			"tenant": schema.StringAttribute{
				Description: "Tenant",
				Optional:    true,
				Computed:    true,
			},
			"tenant_ref_type_": schema.StringAttribute{
				Description: "Object type for tenant field",
				Optional:    true,
				Computed:    true,
			},
			"neighbor_ip_address": schema.StringAttribute{
				Description: "IP address of remote BGP peer",
				Optional:    true,
				Computed:    true,
			},
			"neighbor_as_number": schema.Int64Attribute{
				Description: "Autonomous System Number of remote BGP peer",
				Optional:    true,
				Computed:    true,
			},
			"fabric_interconnect": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"keepalive_timer": schema.Int64Attribute{
				Description: "Interval in seconds between Keepalive messages",
				Optional:    true,
				Computed:    true,
			},
			"hold_timer": schema.Int64Attribute{
				Optional: true,
				Computed: true,
			},
			"connect_timer": schema.Int64Attribute{
				Optional: true,
				Computed: true,
			},
			"advertisement_interval": schema.Int64Attribute{
				Optional: true,
				Computed: true,
			},
			"ebgp_multihop": schema.Int64Attribute{
				Optional: true,
				Computed: true,
			},
			"egress_vlan": schema.Int64Attribute{
				Optional: true,
				Computed: true,
			},
			"source_ip_address": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"anycast_ip_mask": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"md5_password": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"md5_password_encrypted": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"switch_encrypted_md5_password": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"import_route_map": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"export_route_map": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"gateway_mode": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"local_as_number": schema.Int64Attribute{
				Optional: true,
				Computed: true,
			},
			"local_as_no_prepend": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"replace_as": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"max_local_as_occurrences": schema.Int64Attribute{
				Optional: true,
				Computed: true,
			},
			"dynamic_bgp_subnet": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"dynamic_bgp_limits": schema.Int64Attribute{
				Optional: true,
				Computed: true,
			},
			"helper_hop_ip_address": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"enable_bfd": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"bfd_receive_interval": schema.Int64Attribute{
				Optional: true,
				Computed: true,
			},
			"bfd_transmission_interval": schema.Int64Attribute{
				Optional: true,
				Computed: true,
			},
			"bfd_detect_multiplier": schema.Int64Attribute{
				Optional: true,
				Computed: true,
			},
			"bfd_multihop": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"next_hop_self": schema.BoolAttribute{
				Optional: true,
				Computed: true,
			},
			"default_originate": schema.BoolAttribute{
				Description: "Instructs BGP to generate and send a default route 0.0.0.0/0 to the specified neighbor.",
				Optional:    true,
				Computed:    true,
			},
			"export_route_map_ref_type_": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"import_route_map_ref_type_": schema.StringAttribute{
				Optional: true,
				Computed: true,
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
							Computed:    true,
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
							Computed:    true,
						},
						"ipv4_route_prefix": schema.StringAttribute{
							Description: "IPv4 unicast IP address with subnet mask",
							Optional:    true,
							Computed:    true,
						},
						"next_hop_ip_address": schema.StringAttribute{
							Description: "Next Hop IP Address",
							Optional:    true,
							Computed:    true,
						},
						"ad_value": schema.Int64Attribute{
							Description: "Administrative distance value (0-255)",
							Optional:    true,
							Computed:    true,
						},
						"index": schema.Int64Attribute{
							Description: "Index identifying the object",
							Optional:    true,
							Computed:    true,
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

	var config verityGatewayResourceModel
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
		{FieldName: "Md5PasswordEncrypted", APIField: &gatewayProps.Md5PasswordEncrypted, TFValue: plan.Md5PasswordEncrypted},
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
		{FieldName: "SwitchEncryptedMd5Password", APIField: &gatewayProps.SwitchEncryptedMd5Password, TFValue: plan.SwitchEncryptedMd5Password},
	})

	// Handle nullable int64 fields - parse HCL to detect explicit config
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, "verity_gateway", name)

	utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
		{FieldName: "NeighborAsNumber", APIField: &gatewayProps.NeighborAsNumber, TFValue: config.NeighborAsNumber, IsConfigured: configuredAttrs.IsConfigured("neighbor_as_number")},
		{FieldName: "KeepaliveTimer", APIField: &gatewayProps.KeepaliveTimer, TFValue: config.KeepaliveTimer, IsConfigured: configuredAttrs.IsConfigured("keepalive_timer")},
		{FieldName: "HoldTimer", APIField: &gatewayProps.HoldTimer, TFValue: config.HoldTimer, IsConfigured: configuredAttrs.IsConfigured("hold_timer")},
		{FieldName: "ConnectTimer", APIField: &gatewayProps.ConnectTimer, TFValue: config.ConnectTimer, IsConfigured: configuredAttrs.IsConfigured("connect_timer")},
		{FieldName: "AdvertisementInterval", APIField: &gatewayProps.AdvertisementInterval, TFValue: config.AdvertisementInterval, IsConfigured: configuredAttrs.IsConfigured("advertisement_interval")},
		{FieldName: "EbgpMultihop", APIField: &gatewayProps.EbgpMultihop, TFValue: config.EbgpMultihop, IsConfigured: configuredAttrs.IsConfigured("ebgp_multihop")},
		{FieldName: "EgressVlan", APIField: &gatewayProps.EgressVlan, TFValue: config.EgressVlan, IsConfigured: configuredAttrs.IsConfigured("egress_vlan")},
		{FieldName: "LocalAsNumber", APIField: &gatewayProps.LocalAsNumber, TFValue: config.LocalAsNumber, IsConfigured: configuredAttrs.IsConfigured("local_as_number")},
		{FieldName: "MaxLocalAsOccurrences", APIField: &gatewayProps.MaxLocalAsOccurrences, TFValue: config.MaxLocalAsOccurrences, IsConfigured: configuredAttrs.IsConfigured("max_local_as_occurrences")},
		{FieldName: "DynamicBgpLimits", APIField: &gatewayProps.DynamicBgpLimits, TFValue: config.DynamicBgpLimits, IsConfigured: configuredAttrs.IsConfigured("dynamic_bgp_limits")},
		{FieldName: "BfdReceiveInterval", APIField: &gatewayProps.BfdReceiveInterval, TFValue: config.BfdReceiveInterval, IsConfigured: configuredAttrs.IsConfigured("bfd_receive_interval")},
		{FieldName: "BfdTransmissionInterval", APIField: &gatewayProps.BfdTransmissionInterval, TFValue: config.BfdTransmissionInterval, IsConfigured: configuredAttrs.IsConfigured("bfd_transmission_interval")},
		{FieldName: "BfdDetectMultiplier", APIField: &gatewayProps.BfdDetectMultiplier, TFValue: config.BfdDetectMultiplier, IsConfigured: configuredAttrs.IsConfigured("bfd_detect_multiplier")},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties{}
		utils.SetObjectPropertiesFields([]utils.ObjectPropertiesField{
			{Name: "Group", TFValue: op.Group, APIValue: &objProps.Group},
		})
		gatewayProps.ObjectProperties = &objProps
	}

	// Handle static routes
	if len(plan.StaticRoutes) > 0 {
		staticRoutesConfigMap := utils.BuildIndexedConfigMap(config.StaticRoutes)
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

			// Get per-block configured info for nullable Int64 fields
			itemIndex := item.Index.ValueInt64()
			configItem := item // fallback to plan item
			if cfgItem, ok := staticRoutesConfigMap[itemIndex]; ok {
				configItem = cfgItem
			}
			cfg := &utils.IndexedBlockNullableFieldConfig{
				BlockType:       "static_routes",
				BlockIndex:      itemIndex,
				ConfiguredAttrs: configuredAttrs,
			}
			utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
				{FieldName: "AdValue", APIField: &rItem.AdValue, TFValue: configItem.AdValue, IsConfigured: cfg.IsFieldConfigured("ad_value")},
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

	var minState verityGatewayResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if gatewayData, exists := bulkMgr.GetResourceResponse("gateway", name); exists {
			state := populateGatewayState(ctx, minState, gatewayData, r.provCtx.mode)
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

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if gatewayData, exists := r.bulkOpsMgr.GetResourceResponse("gateway", gatewayName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached gateway data for %s from recent operation", gatewayName))
			state = populateGatewayState(ctx, state, gatewayData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

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

	state = populateGatewayState(ctx, state, gatewayMap, r.provCtx.mode)
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
	utils.CompareAndSetStringField(plan.Md5PasswordEncrypted, state.Md5PasswordEncrypted, func(v *string) { gatewayProps.Md5PasswordEncrypted = v }, &hasChanges)
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
	utils.CompareAndSetBoolField(plan.SwitchEncryptedMd5Password, state.SwitchEncryptedMd5Password, func(v *bool) { gatewayProps.SwitchEncryptedMd5Password = v }, &hasChanges)

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
		objProps := openapi.DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties{}
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
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, "verity_gateway", name)
	var config verityGatewayResourceModel
	req.Config.Get(ctx, &config)
	staticRoutesConfigMap := utils.BuildIndexedConfigMap(config.StaticRoutes)

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

			// Get per-block configured info for nullable Int64 fields
			itemIndex := planItem.Index.ValueInt64()
			configItem := planItem // fallback to plan item
			if cfgItem, ok := staticRoutesConfigMap[itemIndex]; ok {
				configItem = cfgItem
			}
			cfg := &utils.IndexedBlockNullableFieldConfig{
				BlockType:       "static_routes",
				BlockIndex:      itemIndex,
				ConfiguredAttrs: configuredAttrs,
			}
			utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
				{FieldName: "AdValue", APIField: &route.AdValue, TFValue: configItem.AdValue, IsConfigured: cfg.IsFieldConfigured("ad_value")},
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

	var minState verityGatewayResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Try to use cached response from bulk operation to populate state with API values
	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if gatewayData, exists := bulkMgr.GetResourceResponse("gateway", name); exists {
			newState := populateGatewayState(ctx, minState, gatewayData, r.provCtx.mode)
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

func populateGatewayState(ctx context.Context, state verityGatewayResourceModel, data map[string]interface{}, mode string) verityGatewayResourceModel {
	const resourceType = gatewayResourceType

	state.Name = utils.MapStringFromAPI(data["name"])

	// Int fields
	state.NeighborAsNumber = utils.MapInt64WithMode(data, "neighbor_as_number", resourceType, mode)
	state.KeepaliveTimer = utils.MapInt64WithMode(data, "keepalive_timer", resourceType, mode)
	state.HoldTimer = utils.MapInt64WithMode(data, "hold_timer", resourceType, mode)
	state.ConnectTimer = utils.MapInt64WithMode(data, "connect_timer", resourceType, mode)
	state.AdvertisementInterval = utils.MapInt64WithMode(data, "advertisement_interval", resourceType, mode)
	state.EbgpMultihop = utils.MapInt64WithMode(data, "ebgp_multihop", resourceType, mode)
	state.EgressVlan = utils.MapInt64WithMode(data, "egress_vlan", resourceType, mode)
	state.LocalAsNumber = utils.MapInt64WithMode(data, "local_as_number", resourceType, mode)
	state.MaxLocalAsOccurrences = utils.MapInt64WithMode(data, "max_local_as_occurrences", resourceType, mode)
	state.DynamicBgpLimits = utils.MapInt64WithMode(data, "dynamic_bgp_limits", resourceType, mode)
	state.BfdReceiveInterval = utils.MapInt64WithMode(data, "bfd_receive_interval", resourceType, mode)
	state.BfdTransmissionInterval = utils.MapInt64WithMode(data, "bfd_transmission_interval", resourceType, mode)
	state.BfdDetectMultiplier = utils.MapInt64WithMode(data, "bfd_detect_multiplier", resourceType, mode)

	// Boolean fields
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)
	state.FabricInterconnect = utils.MapBoolWithMode(data, "fabric_interconnect", resourceType, mode)
	state.LocalAsNoPrepend = utils.MapBoolWithMode(data, "local_as_no_prepend", resourceType, mode)
	state.ReplaceAs = utils.MapBoolWithMode(data, "replace_as", resourceType, mode)
	state.EnableBfd = utils.MapBoolWithMode(data, "enable_bfd", resourceType, mode)
	state.BfdMultihop = utils.MapBoolWithMode(data, "bfd_multihop", resourceType, mode)
	state.NextHopSelf = utils.MapBoolWithMode(data, "next_hop_self", resourceType, mode)
	state.DefaultOriginate = utils.MapBoolWithMode(data, "default_originate", resourceType, mode)
	state.SwitchEncryptedMd5Password = utils.MapBoolWithMode(data, "switch_encrypted_md5_password", resourceType, mode)

	// String fields
	state.Tenant = utils.MapStringWithMode(data, "tenant", resourceType, mode)
	state.TenantRefType = utils.MapStringWithMode(data, "tenant_ref_type_", resourceType, mode)
	state.NeighborIpAddress = utils.MapStringWithMode(data, "neighbor_ip_address", resourceType, mode)
	state.SourceIpAddress = utils.MapStringWithMode(data, "source_ip_address", resourceType, mode)
	state.AnycastIpMask = utils.MapStringWithMode(data, "anycast_ip_mask", resourceType, mode)
	state.Md5Password = utils.MapStringWithMode(data, "md5_password", resourceType, mode)
	state.Md5PasswordEncrypted = utils.MapStringWithMode(data, "md5_password_encrypted", resourceType, mode)
	state.ImportRouteMap = utils.MapStringWithMode(data, "import_route_map", resourceType, mode)
	state.ExportRouteMap = utils.MapStringWithMode(data, "export_route_map", resourceType, mode)
	state.GatewayMode = utils.MapStringWithMode(data, "gateway_mode", resourceType, mode)
	state.DynamicBgpSubnet = utils.MapStringWithMode(data, "dynamic_bgp_subnet", resourceType, mode)
	state.HelperHopIpAddress = utils.MapStringWithMode(data, "helper_hop_ip_address", resourceType, mode)
	state.ExportRouteMapRefType = utils.MapStringWithMode(data, "export_route_map_ref_type_", resourceType, mode)
	state.ImportRouteMapRefType = utils.MapStringWithMode(data, "import_route_map_ref_type_", resourceType, mode)

	// Handle object_properties block
	if utils.FieldAppliesToMode(resourceType, "object_properties", mode) {
		if objProps, ok := data["object_properties"].(map[string]interface{}); ok {
			objPropsModel := verityGatewayObjectPropertiesModel{
				Group: utils.MapStringWithModeNested(objProps, "group", resourceType, "object_properties.group", mode),
			}
			state.ObjectProperties = []verityGatewayObjectPropertiesModel{objPropsModel}
		} else {
			state.ObjectProperties = nil
		}
	} else {
		state.ObjectProperties = nil
	}

	// Handle static_routes list block
	if utils.FieldAppliesToMode(resourceType, "static_routes", mode) {
		if routesData, ok := data["static_routes"].([]interface{}); ok && len(routesData) > 0 {
			var routesList []verityGatewayStaticRoutesModel

			for _, item := range routesData {
				itemMap, ok := item.(map[string]interface{})
				if !ok {
					continue
				}

				routeItem := verityGatewayStaticRoutesModel{
					Enable:           utils.MapBoolWithModeNested(itemMap, "enable", resourceType, "static_routes.enable", mode),
					Ipv4RoutePrefix:  utils.MapStringWithModeNested(itemMap, "ipv4_route_prefix", resourceType, "static_routes.ipv4_route_prefix", mode),
					NextHopIpAddress: utils.MapStringWithModeNested(itemMap, "next_hop_ip_address", resourceType, "static_routes.next_hop_ip_address", mode),
					AdValue:          utils.MapInt64WithModeNested(itemMap, "ad_value", resourceType, "static_routes.ad_value", mode),
					Index:            utils.MapInt64WithModeNested(itemMap, "index", resourceType, "static_routes.index", mode),
				}

				routesList = append(routesList, routeItem)
			}

			state.StaticRoutes = routesList
		} else {
			state.StaticRoutes = nil
		}
	} else {
		state.StaticRoutes = nil
	}

	return state
}

func (r *verityGatewayResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// =========================================================================
	// Skip if deleting
	// =========================================================================
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityGatewayResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Mode-aware field nullification
	// Set fields that don't apply to current mode to null to prevent
	// "known after apply" messages for irrelevant fields.
	// =========================================================================
	const resourceType = gatewayResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyStrings(
		"tenant", "tenant_ref_type_", "neighbor_ip_address", "source_ip_address",
		"anycast_ip_mask", "md5_password", "md5_password_encrypted",
		"import_route_map", "export_route_map", "gateway_mode",
		"dynamic_bgp_subnet", "helper_hop_ip_address",
		"export_route_map_ref_type_", "import_route_map_ref_type_",
	)

	nullifier.NullifyBools(
		"enable", "fabric_interconnect", "local_as_no_prepend", "replace_as",
		"enable_bfd", "bfd_multihop", "next_hop_self", "default_originate",
		"switch_encrypted_md5_password",
	)

	nullifier.NullifyInt64s(
		"neighbor_as_number", "keepalive_timer", "hold_timer", "connect_timer",
		"advertisement_interval", "ebgp_multihop", "egress_vlan",
		"local_as_number", "max_local_as_occurrences", "dynamic_bgp_limits",
		"bfd_receive_interval", "bfd_transmission_interval", "bfd_detect_multiplier",
	)

	nullifier.NullifyNestedBlocks(
		"static_routes", "object_properties",
	)

	// =========================================================================
	// Skip UPDATE-specific logic during CREATE
	// =========================================================================
	if req.State.Raw.IsNull() {
		return
	}

	// =========================================================================
	// UPDATE operation - get state and config
	// =========================================================================
	var state verityGatewayResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config verityGatewayResourceModel
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
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, "verity_gateway", name)

	utils.HandleNullableFields(utils.NullableFieldsConfig{
		Ctx:             ctx,
		Plan:            &resp.Plan,
		ConfiguredAttrs: configuredAttrs,
		Int64Fields: []utils.NullableInt64Field{
			{AttrName: "neighbor_as_number", ConfigVal: config.NeighborAsNumber, StateVal: state.NeighborAsNumber},
			{AttrName: "keepalive_timer", ConfigVal: config.KeepaliveTimer, StateVal: state.KeepaliveTimer},
			{AttrName: "hold_timer", ConfigVal: config.HoldTimer, StateVal: state.HoldTimer},
			{AttrName: "connect_timer", ConfigVal: config.ConnectTimer, StateVal: state.ConnectTimer},
			{AttrName: "advertisement_interval", ConfigVal: config.AdvertisementInterval, StateVal: state.AdvertisementInterval},
			{AttrName: "ebgp_multihop", ConfigVal: config.EbgpMultihop, StateVal: state.EbgpMultihop},
			{AttrName: "egress_vlan", ConfigVal: config.EgressVlan, StateVal: state.EgressVlan},
			{AttrName: "local_as_number", ConfigVal: config.LocalAsNumber, StateVal: state.LocalAsNumber},
			{AttrName: "max_local_as_occurrences", ConfigVal: config.MaxLocalAsOccurrences, StateVal: state.MaxLocalAsOccurrences},
			{AttrName: "dynamic_bgp_limits", ConfigVal: config.DynamicBgpLimits, StateVal: state.DynamicBgpLimits},
			{AttrName: "bfd_receive_interval", ConfigVal: config.BfdReceiveInterval, StateVal: state.BfdReceiveInterval},
			{AttrName: "bfd_transmission_interval", ConfigVal: config.BfdTransmissionInterval, StateVal: state.BfdTransmissionInterval},
			{AttrName: "bfd_detect_multiplier", ConfigVal: config.BfdDetectMultiplier, StateVal: state.BfdDetectMultiplier},
		},
	})

	// =========================================================================
	// Handle nullable fields in nested blocks
	// =========================================================================
	for i, configRoute := range config.StaticRoutes {
		routeIndex := configRoute.Index.ValueInt64()
		var stateRoute *verityGatewayStaticRoutesModel
		for j := range state.StaticRoutes {
			if state.StaticRoutes[j].Index.ValueInt64() == routeIndex {
				stateRoute = &state.StaticRoutes[j]
				break
			}
		}

		if stateRoute != nil {
			utils.HandleNullableNestedFields(utils.NullableNestedFieldsConfig{
				Ctx:             ctx,
				Plan:            &resp.Plan,
				ConfiguredAttrs: configuredAttrs,
				BlockType:       "static_routes",
				BlockListPath:   "static_routes",
				BlockListIndex:  i,
				Int64Fields: []utils.NullableNestedInt64Field{
					{BlockIndex: routeIndex, AttrName: "ad_value", ConfigVal: configRoute.AdValue, StateVal: stateRoute.AdValue},
				},
			})
		}
	}
}
