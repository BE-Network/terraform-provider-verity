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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

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
	bulkOpsMgr           *utils.BulkOperationManager
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
	Md5PasswordEncrypted    types.String                         `tfsdk:"md5_password_encrypted"`
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
				Default:     booldefault.StaticBool(false),
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
			"md5_password_encrypted": schema.StringAttribute{
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
				Computed: true,
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
	gatewayProps := openapi.ConfigPutRequestGatewayGatewayName{}
	gatewayProps.Name = openapi.PtrString(name)

	if !plan.Enable.IsNull() {
		gatewayProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
	}
	if !plan.Tenant.IsNull() {
		gatewayProps.Tenant = openapi.PtrString(plan.Tenant.ValueString())
	}
	if !plan.TenantRefType.IsNull() {
		gatewayProps.TenantRefType = openapi.PtrString(plan.TenantRefType.ValueString())
	}
	if !plan.NeighborIpAddress.IsNull() {
		gatewayProps.NeighborIpAddress = openapi.PtrString(plan.NeighborIpAddress.ValueString())
	}
	if !plan.NeighborAsNumber.IsNull() {
		val := int32(plan.NeighborAsNumber.ValueInt64())
		gatewayProps.NeighborAsNumber = *openapi.NewNullableInt32(&val)
	} else {
		gatewayProps.NeighborAsNumber = *openapi.NewNullableInt32(nil)
	}
	if !plan.FabricInterconnect.IsNull() {
		gatewayProps.FabricInterconnect = openapi.PtrBool(plan.FabricInterconnect.ValueBool())
	}
	if !plan.KeepaliveTimer.IsNull() {
		gatewayProps.KeepaliveTimer = openapi.PtrInt32(int32(plan.KeepaliveTimer.ValueInt64()))
	}
	if !plan.HoldTimer.IsNull() {
		gatewayProps.HoldTimer = openapi.PtrInt32(int32(plan.HoldTimer.ValueInt64()))
	}
	if !plan.ConnectTimer.IsNull() {
		gatewayProps.ConnectTimer = openapi.PtrInt32(int32(plan.ConnectTimer.ValueInt64()))
	}
	if !plan.AdvertisementInterval.IsNull() {
		gatewayProps.AdvertisementInterval = openapi.PtrInt32(int32(plan.AdvertisementInterval.ValueInt64()))
	}
	if !plan.EbgpMultihop.IsNull() {
		gatewayProps.SetEbgpMultihop(int32(plan.EbgpMultihop.ValueInt64()))
	}
	if !plan.EgressVlan.IsNull() {
		val := int32(plan.EgressVlan.ValueInt64())
		gatewayProps.EgressVlan = *openapi.NewNullableInt32(&val)
	} else {
		gatewayProps.EgressVlan = *openapi.NewNullableInt32(nil)
	}
	if !plan.SourceIpAddress.IsNull() {
		gatewayProps.SourceIpAddress = openapi.PtrString(plan.SourceIpAddress.ValueString())
	}
	if !plan.AnycastIpMask.IsNull() {
		gatewayProps.AnycastIpMask = openapi.PtrString(plan.AnycastIpMask.ValueString())
	}
	if !plan.Md5Password.IsNull() {
		gatewayProps.Md5Password = openapi.PtrString(plan.Md5Password.ValueString())
	}
	if !plan.ImportRouteMap.IsNull() {
		gatewayProps.ImportRouteMap = openapi.PtrString(plan.ImportRouteMap.ValueString())
	}
	if !plan.ExportRouteMap.IsNull() {
		gatewayProps.ExportRouteMap = openapi.PtrString(plan.ExportRouteMap.ValueString())
	}
	if !plan.GatewayMode.IsNull() {
		gatewayProps.GatewayMode = openapi.PtrString(plan.GatewayMode.ValueString())
	}
	if !plan.LocalAsNumber.IsNull() {
		val := int32(plan.LocalAsNumber.ValueInt64())
		gatewayProps.LocalAsNumber = *openapi.NewNullableInt32(&val)
	} else {
		gatewayProps.LocalAsNumber = *openapi.NewNullableInt32(nil)
	}
	if !plan.LocalAsNoPrepend.IsNull() {
		gatewayProps.LocalAsNoPrepend = openapi.PtrBool(plan.LocalAsNoPrepend.ValueBool())
	}
	if !plan.ReplaceAs.IsNull() {
		gatewayProps.ReplaceAs = openapi.PtrBool(plan.ReplaceAs.ValueBool())
	}
	if !plan.MaxLocalAsOccurrences.IsNull() {
		val := int32(plan.MaxLocalAsOccurrences.ValueInt64())
		gatewayProps.MaxLocalAsOccurrences = *openapi.NewNullableInt32(&val)
	} else {
		gatewayProps.MaxLocalAsOccurrences = *openapi.NewNullableInt32(nil)
	}
	if !plan.DynamicBgpSubnet.IsNull() {
		gatewayProps.DynamicBgpSubnet = openapi.PtrString(plan.DynamicBgpSubnet.ValueString())
	}
	if !plan.DynamicBgpLimits.IsNull() {
		val := int32(plan.DynamicBgpLimits.ValueInt64())
		gatewayProps.DynamicBgpLimits = *openapi.NewNullableInt32(&val)
	} else {
		gatewayProps.DynamicBgpLimits = *openapi.NewNullableInt32(nil)
	}
	if !plan.HelperHopIpAddress.IsNull() {
		gatewayProps.HelperHopIpAddress = openapi.PtrString(plan.HelperHopIpAddress.ValueString())
	}
	if !plan.EnableBfd.IsNull() {
		gatewayProps.EnableBfd = openapi.PtrBool(plan.EnableBfd.ValueBool())
	}
	if !plan.BfdReceiveInterval.IsNull() {
		val := int32(plan.BfdReceiveInterval.ValueInt64())
		gatewayProps.BfdReceiveInterval = *openapi.NewNullableInt32(&val)
	} else {
		gatewayProps.BfdReceiveInterval = *openapi.NewNullableInt32(nil)
	}
	if !plan.BfdTransmissionInterval.IsNull() {
		val := int32(plan.BfdTransmissionInterval.ValueInt64())
		gatewayProps.BfdTransmissionInterval = *openapi.NewNullableInt32(&val)
	} else {
		gatewayProps.BfdTransmissionInterval = *openapi.NewNullableInt32(nil)
	}
	if !plan.BfdDetectMultiplier.IsNull() {
		val := int32(plan.BfdDetectMultiplier.ValueInt64())
		gatewayProps.BfdDetectMultiplier = *openapi.NewNullableInt32(&val)
	} else {
		gatewayProps.BfdDetectMultiplier = *openapi.NewNullableInt32(nil)
	}
	if !plan.BfdMultihop.IsNull() {
		gatewayProps.BfdMultihop = openapi.PtrBool(plan.BfdMultihop.ValueBool())
	}
	if !plan.NextHopSelf.IsNull() {
		gatewayProps.NextHopSelf = openapi.PtrBool(plan.NextHopSelf.ValueBool())
	}
	if !plan.DefaultOriginate.IsNull() {
		gatewayProps.DefaultOriginate = openapi.PtrBool(plan.DefaultOriginate.ValueBool())
	}
	if !plan.ExportRouteMapRefType.IsNull() {
		gatewayProps.ExportRouteMapRefType = openapi.PtrString(plan.ExportRouteMapRefType.ValueString())
	}
	if !plan.ImportRouteMapRefType.IsNull() {
		gatewayProps.ImportRouteMapRefType = openapi.PtrString(plan.ImportRouteMapRefType.ValueString())
	}

	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties{}
		if !op.Group.IsNull() {
			objProps.Group = openapi.PtrString(op.Group.ValueString())
		} else {
			objProps.Group = nil
		}
		gatewayProps.ObjectProperties = &objProps
	}

	if len(plan.StaticRoutes) > 0 {
		routes := make([]openapi.ConfigPutRequestGatewayGatewayNameStaticRoutesInner, len(plan.StaticRoutes))
		for i, route := range plan.StaticRoutes {
			rItem := openapi.ConfigPutRequestGatewayGatewayNameStaticRoutesInner{}
			if !route.Enable.IsNull() {
				rItem.Enable = openapi.PtrBool(route.Enable.ValueBool())
			}
			if !route.Ipv4RoutePrefix.IsNull() {
				rItem.Ipv4RoutePrefix = openapi.PtrString(route.Ipv4RoutePrefix.ValueString())
			}
			if !route.NextHopIpAddress.IsNull() {
				rItem.NextHopIpAddress = openapi.PtrString(route.NextHopIpAddress.ValueString())
			}
			if !route.AdValue.IsNull() {
				adVal := int32(route.AdValue.ValueInt64())
				rItem.AdValue = *openapi.NewNullableInt32(&adVal)
			} else {
				rItem.AdValue = *openapi.NewNullableInt32(nil)
			}
			if !route.Index.IsNull() {
				rItem.Index = openapi.PtrInt32(int32(route.Index.ValueInt64()))
			}
			routes[i] = rItem
		}
		gatewayProps.StaticRoutes = routes
	}

	operationID := r.bulkOpsMgr.AddGatewayPut(ctx, name, gatewayProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for gateway creation operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Create Gateway",
			fmt.Sprintf("Error creating gateway %s: %v", name, err),
		)
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

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentGatewayOperations() {
		tflog.Info(ctx, fmt.Sprintf("Skipping gateway %s verification – trusting recent successful API operation", gatewayName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching gateways for verification of %s", gatewayName))

	type GatewaysResponse struct {
		Gateway map[string]interface{} `json:"gateway"`
	}

	var result GatewaysResponse
	var err error
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		gatewaysData, fetchErr := getCachedResponse(ctx, r.provCtx, "gateways", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch gateways")
			respAPI, err := r.client.GatewaysAPI.GatewaysGet(ctx).Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading gateways: %v", err)
			}
			defer respAPI.Body.Close()

			var res GatewaysResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return nil, fmt.Errorf("failed to decode gateways response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d gateways", len(res.Gateway)))
			return res, nil
		})
		if fetchErr != nil {
			err = fetchErr
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch gateways on attempt %d, retrying in %v", attempt+1, sleepTime))
			time.Sleep(sleepTime)
			continue
		}
		result = gatewaysData.(GatewaysResponse)
		break
	}
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Read Gateways",
			fmt.Sprintf("Error reading gateways: %v", err),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for gateway with ID: %s", gatewayName))
	var gatewayData map[string]interface{}
	exists := false

	if data, ok := result.Gateway[gatewayName].(map[string]interface{}); ok {
		gatewayData = data
		exists = true
		tflog.Debug(ctx, fmt.Sprintf("Found gateway directly by ID: %s", gatewayName))
	} else {
		for apiName, g := range result.Gateway {
			gateway, ok := g.(map[string]interface{})
			if !ok {
				continue
			}

			if name, ok := gateway["name"].(string); ok && name == gatewayName {
				gatewayData = gateway
				gatewayName = apiName
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found gateway with name '%s' under API key '%s'", name, apiName))
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Gateway with ID '%s' not found in API response", gatewayName))
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(fmt.Sprintf("%v", gatewayData["name"]))

	if enable, ok := gatewayData["enable"].(bool); ok {
		state.Enable = types.BoolValue(enable)
	} else {
		state.Enable = types.BoolNull()
	}

	if objProps, ok := gatewayData["object_properties"].(map[string]interface{}); ok {
		if group, ok := objProps["group"].(string); ok {
			state.ObjectProperties = []verityGatewayObjectPropertiesModel{
				{Group: types.StringValue(group)},
			}
		} else {
			state.ObjectProperties = []verityGatewayObjectPropertiesModel{
				{Group: types.StringValue("")},
			}
		}
	} else {
		state.ObjectProperties = []verityGatewayObjectPropertiesModel{
			{Group: types.StringValue("")},
		}
	}

	stringAttrs := map[string]*types.String{
		"tenant":                     &state.Tenant,
		"tenant_ref_type_":           &state.TenantRefType,
		"neighbor_ip_address":        &state.NeighborIpAddress,
		"source_ip_address":          &state.SourceIpAddress,
		"anycast_ip_mask":            &state.AnycastIpMask,
		"md5_password":               &state.Md5Password,
		"md5_password_encrypted":     &state.Md5PasswordEncrypted,
		"import_route_map":           &state.ImportRouteMap,
		"export_route_map":           &state.ExportRouteMap,
		"gateway_mode":               &state.GatewayMode,
		"dynamic_bgp_subnet":         &state.DynamicBgpSubnet,
		"helper_hop_ip_address":      &state.HelperHopIpAddress,
		"export_route_map_ref_type_": &state.ExportRouteMapRefType,
		"import_route_map_ref_type_": &state.ImportRouteMapRefType,
	}

	for apiKey, stateField := range stringAttrs {
		if value, ok := gatewayData[apiKey].(string); ok {
			*stateField = types.StringValue(value)
		} else {
			*stateField = types.StringNull()
		}
	}

	boolAttrs := map[string]*types.Bool{
		"fabric_interconnect": &state.FabricInterconnect,
		"local_as_no_prepend": &state.LocalAsNoPrepend,
		"replace_as":          &state.ReplaceAs,
		"enable_bfd":          &state.EnableBfd,
		"bfd_multihop":        &state.BfdMultihop,
		"next_hop_self":       &state.NextHopSelf,
		"default_originate":   &state.DefaultOriginate,
	}

	for apiKey, stateField := range boolAttrs {
		if value, ok := gatewayData[apiKey].(bool); ok {
			*stateField = types.BoolValue(value)
		} else {
			*stateField = types.BoolNull()
		}
	}

	intAttrs := map[string]*types.Int64{
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

	for apiKey, stateField := range intAttrs {
		if value, ok := gatewayData[apiKey]; ok && value != nil {
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

	if routes, ok := gatewayData["static_routes"].([]interface{}); ok && len(routes) > 0 {
		var staticRoutes []verityGatewayStaticRoutesModel

		for _, r := range routes {
			route, ok := r.(map[string]interface{})
			if !ok {
				continue
			}

			srModel := verityGatewayStaticRoutesModel{}

			if enable, ok := route["enable"].(bool); ok {
				srModel.Enable = types.BoolValue(enable)
			} else {
				srModel.Enable = types.BoolNull()
			}

			if prefix, ok := route["ipv4_route_prefix"].(string); ok {
				srModel.Ipv4RoutePrefix = types.StringValue(prefix)
			} else {
				srModel.Ipv4RoutePrefix = types.StringNull()
			}

			if nextHop, ok := route["next_hop_ip_address"].(string); ok {
				srModel.NextHopIpAddress = types.StringValue(nextHop)
			} else {
				srModel.NextHopIpAddress = types.StringNull()
			}

			if adValue, exists := route["ad_value"]; exists && adValue != nil {
				switch v := adValue.(type) {
				case int:
					srModel.AdValue = types.Int64Value(int64(v))
				case int32:
					srModel.AdValue = types.Int64Value(int64(v))
				case int64:
					srModel.AdValue = types.Int64Value(v)
				case float64:
					srModel.AdValue = types.Int64Value(int64(v))
				case string:
					if intVal, err := strconv.ParseInt(v, 10, 64); err == nil {
						srModel.AdValue = types.Int64Value(intVal)
					} else {
						srModel.AdValue = types.Int64Null()
					}
				default:
					srModel.AdValue = types.Int64Null()
				}
			} else {
				srModel.AdValue = types.Int64Null()
			}

			if index, exists := route["index"]; exists && index != nil {
				switch v := index.(type) {
				case int:
					srModel.Index = types.Int64Value(int64(v))
				case float64:
					srModel.Index = types.Int64Value(int64(v))
				default:
					srModel.Index = types.Int64Null()
				}
			} else {
				srModel.Index = types.Int64Null()
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
	gatewayProps := openapi.ConfigPutRequestGatewayGatewayName{}
	hasChanges := false

	if !plan.Name.Equal(state.Name) {
		gatewayProps.Name = openapi.PtrString(name)
		hasChanges = true
	}

	if len(plan.ObjectProperties) > 0 {
		if len(state.ObjectProperties) == 0 || !plan.ObjectProperties[0].Group.Equal(state.ObjectProperties[0].Group) {
			objProps := openapi.ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties{}
			if !plan.ObjectProperties[0].Group.IsNull() {
				objProps.Group = openapi.PtrString(plan.ObjectProperties[0].Group.ValueString())
			} else {
				objProps.Group = nil
			}
			gatewayProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	if !plan.Enable.Equal(state.Enable) {
		gatewayProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
		hasChanges = true
	}
	if !plan.FabricInterconnect.Equal(state.FabricInterconnect) {
		gatewayProps.FabricInterconnect = openapi.PtrBool(plan.FabricInterconnect.ValueBool())
		hasChanges = true
	}
	if !plan.LocalAsNoPrepend.Equal(state.LocalAsNoPrepend) {
		gatewayProps.LocalAsNoPrepend = openapi.PtrBool(plan.LocalAsNoPrepend.ValueBool())
		hasChanges = true
	}
	if !plan.ReplaceAs.Equal(state.ReplaceAs) {
		gatewayProps.ReplaceAs = openapi.PtrBool(plan.ReplaceAs.ValueBool())
		hasChanges = true
	}
	if !plan.EnableBfd.Equal(state.EnableBfd) {
		gatewayProps.EnableBfd = openapi.PtrBool(plan.EnableBfd.ValueBool())
		hasChanges = true
	}
	if !plan.BfdMultihop.Equal(state.BfdMultihop) {
		gatewayProps.BfdMultihop = openapi.PtrBool(plan.BfdMultihop.ValueBool())
		hasChanges = true
	}
	if !plan.NextHopSelf.Equal(state.NextHopSelf) {
		gatewayProps.NextHopSelf = openapi.PtrBool(plan.NextHopSelf.ValueBool())
		hasChanges = true
	}
	if !plan.Tenant.Equal(state.Tenant) {
		gatewayProps.Tenant = openapi.PtrString(plan.Tenant.ValueString())
		hasChanges = true
	}
	if !plan.TenantRefType.Equal(state.TenantRefType) {
		gatewayProps.TenantRefType = openapi.PtrString(plan.TenantRefType.ValueString())
		hasChanges = true
	}
	if !plan.NeighborIpAddress.Equal(state.NeighborIpAddress) {
		gatewayProps.NeighborIpAddress = openapi.PtrString(plan.NeighborIpAddress.ValueString())
		hasChanges = true
	}
	if !plan.SourceIpAddress.Equal(state.SourceIpAddress) {
		gatewayProps.SourceIpAddress = openapi.PtrString(plan.SourceIpAddress.ValueString())
		hasChanges = true
	}
	if !plan.AnycastIpMask.Equal(state.AnycastIpMask) {
		gatewayProps.AnycastIpMask = openapi.PtrString(plan.AnycastIpMask.ValueString())
		hasChanges = true
	}
	if !plan.Md5Password.Equal(state.Md5Password) {
		gatewayProps.Md5Password = openapi.PtrString(plan.Md5Password.ValueString())
		hasChanges = true
	}
	if !plan.ImportRouteMap.Equal(state.ImportRouteMap) {
		gatewayProps.ImportRouteMap = openapi.PtrString(plan.ImportRouteMap.ValueString())
		hasChanges = true
	}
	if !plan.ExportRouteMap.Equal(state.ExportRouteMap) {
		gatewayProps.ExportRouteMap = openapi.PtrString(plan.ExportRouteMap.ValueString())
		hasChanges = true
	}
	if !plan.GatewayMode.Equal(state.GatewayMode) {
		gatewayProps.GatewayMode = openapi.PtrString(plan.GatewayMode.ValueString())
		hasChanges = true
	}
	if !plan.DynamicBgpSubnet.Equal(state.DynamicBgpSubnet) {
		gatewayProps.DynamicBgpSubnet = openapi.PtrString(plan.DynamicBgpSubnet.ValueString())
		hasChanges = true
	}
	if !plan.HelperHopIpAddress.Equal(state.HelperHopIpAddress) {
		gatewayProps.HelperHopIpAddress = openapi.PtrString(plan.HelperHopIpAddress.ValueString())
		hasChanges = true
	}
	if !plan.DefaultOriginate.Equal(state.DefaultOriginate) {
		gatewayProps.DefaultOriginate = openapi.PtrBool(plan.DefaultOriginate.ValueBool())
		hasChanges = true
	}
	if !plan.ExportRouteMapRefType.Equal(state.ExportRouteMapRefType) {
		gatewayProps.ExportRouteMapRefType = openapi.PtrString(plan.ExportRouteMapRefType.ValueString())
		hasChanges = true
	}
	if !plan.ImportRouteMapRefType.Equal(state.ImportRouteMapRefType) {
		gatewayProps.ImportRouteMapRefType = openapi.PtrString(plan.ImportRouteMapRefType.ValueString())
		hasChanges = true
	}
	if !plan.NeighborAsNumber.Equal(state.NeighborAsNumber) {
		if !plan.NeighborAsNumber.IsNull() {
			val := int32(plan.NeighborAsNumber.ValueInt64())
			gatewayProps.NeighborAsNumber = *openapi.NewNullableInt32(&val)
		} else {
			gatewayProps.NeighborAsNumber = *openapi.NewNullableInt32(nil)
		}
		hasChanges = true
	}
	if !plan.KeepaliveTimer.Equal(state.KeepaliveTimer) {
		gatewayProps.KeepaliveTimer = openapi.PtrInt32(int32(plan.KeepaliveTimer.ValueInt64()))
		hasChanges = true
	}
	if !plan.HoldTimer.Equal(state.HoldTimer) {
		gatewayProps.HoldTimer = openapi.PtrInt32(int32(plan.HoldTimer.ValueInt64()))
		hasChanges = true
	}
	if !plan.ConnectTimer.Equal(state.ConnectTimer) {
		gatewayProps.ConnectTimer = openapi.PtrInt32(int32(plan.ConnectTimer.ValueInt64()))
		hasChanges = true
	}
	if !plan.AdvertisementInterval.Equal(state.AdvertisementInterval) {
		gatewayProps.AdvertisementInterval = openapi.PtrInt32(int32(plan.AdvertisementInterval.ValueInt64()))
		hasChanges = true
	}
	if !plan.EbgpMultihop.Equal(state.EbgpMultihop) {
		gatewayProps.SetEbgpMultihop(int32(plan.EbgpMultihop.ValueInt64()))
		hasChanges = true
	}
	if !plan.EgressVlan.Equal(state.EgressVlan) {
		if !plan.EgressVlan.IsNull() {
			val := int32(plan.EgressVlan.ValueInt64())
			gatewayProps.EgressVlan = *openapi.NewNullableInt32(&val)
		} else {
			gatewayProps.EgressVlan = *openapi.NewNullableInt32(nil)
		}
		hasChanges = true
	}
	if !plan.LocalAsNumber.Equal(state.LocalAsNumber) {
		if !plan.LocalAsNumber.IsNull() {
			val := int32(plan.LocalAsNumber.ValueInt64())
			gatewayProps.LocalAsNumber = *openapi.NewNullableInt32(&val)
		} else {
			gatewayProps.LocalAsNumber = *openapi.NewNullableInt32(nil)
		}
		hasChanges = true
	}
	if !plan.MaxLocalAsOccurrences.Equal(state.MaxLocalAsOccurrences) {
		if !plan.MaxLocalAsOccurrences.IsNull() {
			val := int32(plan.MaxLocalAsOccurrences.ValueInt64())
			gatewayProps.MaxLocalAsOccurrences = *openapi.NewNullableInt32(&val)
		} else {
			gatewayProps.MaxLocalAsOccurrences = *openapi.NewNullableInt32(nil)
		}
		hasChanges = true
	}
	if !plan.DynamicBgpLimits.Equal(state.DynamicBgpLimits) {
		if !plan.DynamicBgpLimits.IsNull() {
			val := int32(plan.DynamicBgpLimits.ValueInt64())
			gatewayProps.DynamicBgpLimits = *openapi.NewNullableInt32(&val)
		} else {
			gatewayProps.DynamicBgpLimits = *openapi.NewNullableInt32(nil)
		}
		hasChanges = true
	}
	if !plan.BfdReceiveInterval.Equal(state.BfdReceiveInterval) {
		if !plan.BfdReceiveInterval.IsNull() {
			val := int32(plan.BfdReceiveInterval.ValueInt64())
			gatewayProps.BfdReceiveInterval = *openapi.NewNullableInt32(&val)
		} else {
			gatewayProps.BfdReceiveInterval = *openapi.NewNullableInt32(nil)
		}
		hasChanges = true
	}
	if !plan.BfdTransmissionInterval.Equal(state.BfdTransmissionInterval) {
		if !plan.BfdTransmissionInterval.IsNull() {
			val := int32(plan.BfdTransmissionInterval.ValueInt64())
			gatewayProps.BfdTransmissionInterval = *openapi.NewNullableInt32(&val)
		} else {
			gatewayProps.BfdTransmissionInterval = *openapi.NewNullableInt32(nil)
		}
		hasChanges = true
	}
	if !plan.BfdDetectMultiplier.Equal(state.BfdDetectMultiplier) {
		if !plan.BfdDetectMultiplier.IsNull() {
			val := int32(plan.BfdDetectMultiplier.ValueInt64())
			gatewayProps.BfdDetectMultiplier = *openapi.NewNullableInt32(&val)
		} else {
			gatewayProps.BfdDetectMultiplier = *openapi.NewNullableInt32(nil)
		}
		hasChanges = true
	}

	oldStaticRoutesByIndex := make(map[int64]verityGatewayStaticRoutesModel)
	for _, sr := range state.StaticRoutes {
		if !sr.Index.IsNull() {
			oldStaticRoutesByIndex[sr.Index.ValueInt64()] = sr
		}
	}

	var changedStaticRoutes []openapi.ConfigPutRequestGatewayGatewayNameStaticRoutesInner
	for _, sr := range plan.StaticRoutes {
		if sr.Index.IsNull() {
			continue
		}
		index := sr.Index.ValueInt64()
		oldRoute, exists := oldStaticRoutesByIndex[index]

		routeChanged := !exists ||
			!sr.Enable.Equal(oldRoute.Enable) ||
			!sr.Ipv4RoutePrefix.Equal(oldRoute.Ipv4RoutePrefix) ||
			!sr.NextHopIpAddress.Equal(oldRoute.NextHopIpAddress) ||
			!sr.AdValue.Equal(oldRoute.AdValue)

		if routeChanged {
			route := openapi.ConfigPutRequestGatewayGatewayNameStaticRoutesInner{
				Index: openapi.PtrInt32(int32(index)),
			}
			if !sr.Enable.IsNull() {
				route.Enable = openapi.PtrBool(sr.Enable.ValueBool())
			} else {
				route.Enable = openapi.PtrBool(false)
			}
			if !sr.Ipv4RoutePrefix.IsNull() {
				route.Ipv4RoutePrefix = openapi.PtrString(sr.Ipv4RoutePrefix.ValueString())
			} else {
				route.Ipv4RoutePrefix = openapi.PtrString("")
			}
			if !sr.NextHopIpAddress.IsNull() {
				route.NextHopIpAddress = openapi.PtrString(sr.NextHopIpAddress.ValueString())
			} else {
				route.NextHopIpAddress = openapi.PtrString("")
			}
			if !sr.AdValue.IsNull() {
				adVal := int32(sr.AdValue.ValueInt64())
				route.AdValue = *openapi.NewNullableInt32(&adVal)
			} else {
				route.AdValue = *openapi.NewNullableInt32(nil)
			}
			changedStaticRoutes = append(changedStaticRoutes, route)
		}
	}
	if len(changedStaticRoutes) > 0 {
		gatewayProps.StaticRoutes = changedStaticRoutes
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	operationID := r.bulkOpsMgr.AddGatewayPatch(ctx, name, gatewayProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for gateway update operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Update Gateway",
			fmt.Sprintf("Error updating gateway %s: %v", name, err),
		)
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
	operationID := r.bulkOpsMgr.AddGatewayDelete(ctx, name)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for gateway deletion operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Delete Gateway",
			fmt.Sprintf("Error deleting gateway %s: %v", name, err),
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Gateway %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "gateways")
	resp.State.RemoveResource(ctx)
}

func (r *verityGatewayResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
