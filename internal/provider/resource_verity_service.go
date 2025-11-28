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
	_ resource.Resource                = &verityServiceResource{}
	_ resource.ResourceWithConfigure   = &verityServiceResource{}
	_ resource.ResourceWithImportState = &verityServiceResource{}
	_ resource.ResourceWithModifyPlan  = &verityServiceResource{}
)

func NewVerityServiceResource() resource.Resource {
	return &verityServiceResource{}
}

type verityServiceResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type verityServiceResourceModel struct {
	Name                                        types.String                         `tfsdk:"name"`
	Enable                                      types.Bool                           `tfsdk:"enable"`
	ObjectProperties                            []verityServiceObjectPropertiesModel `tfsdk:"object_properties"`
	Vlan                                        types.Int64                          `tfsdk:"vlan"`
	Vni                                         types.Int64                          `tfsdk:"vni"`
	VniAutoAssigned                             types.Bool                           `tfsdk:"vni_auto_assigned_"`
	Tenant                                      types.String                         `tfsdk:"tenant"`
	TenantRefType                               types.String                         `tfsdk:"tenant_ref_type_"`
	DhcpServerIpv4                              types.String                         `tfsdk:"dhcp_server_ipv4"`
	DhcpServerIpv6                              types.String                         `tfsdk:"dhcp_server_ipv6"`
	Mtu                                         types.Int64                          `tfsdk:"mtu"`
	AnycastIpv4Mask                             types.String                         `tfsdk:"anycast_ipv4_mask"`
	AnycastIpv6Mask                             types.String                         `tfsdk:"anycast_ipv6_mask"`
	MaxUpstreamRateMbps                         types.Int64                          `tfsdk:"max_upstream_rate_mbps"`
	MaxDownstreamRateMbps                       types.Int64                          `tfsdk:"max_downstream_rate_mbps"`
	PacketPriority                              types.String                         `tfsdk:"packet_priority"`
	MulticastManagementMode                     types.String                         `tfsdk:"multicast_management_mode"`
	TaggedPackets                               types.Bool                           `tfsdk:"tagged_packets"`
	Tls                                         types.Bool                           `tfsdk:"tls"`
	AllowLocalSwitching                         types.Bool                           `tfsdk:"allow_local_switching"`
	ActAsMulticastQuerier                       types.Bool                           `tfsdk:"act_as_multicast_querier"`
	BlockUnknownUnicastFlood                    types.Bool                           `tfsdk:"block_unknown_unicast_flood"`
	BlockDownstreamDhcpServer                   types.Bool                           `tfsdk:"block_downstream_dhcp_server"`
	IsManagementService                         types.Bool                           `tfsdk:"is_management_service"`
	UseDscpToPBitMappingForL3PacketsIfAvailable types.Bool                           `tfsdk:"use_dscp_to_p_bit_mapping_for_l3_packets_if_available"`
	AllowFastLeave                              types.Bool                           `tfsdk:"allow_fast_leave"`
	MstInstance                                 types.Int64                          `tfsdk:"mst_instance"`
	PolicyBasedRouting                          types.String                         `tfsdk:"policy_based_routing"`
	PolicyBasedRoutingRefType                   types.String                         `tfsdk:"policy_based_routing_ref_type_"`
}

type verityServiceObjectPropertiesModel struct {
	Group                  types.String `tfsdk:"group"`
	OnSummary              types.Bool   `tfsdk:"on_summary"`
	WarnOnNoExternalSource types.Bool   `tfsdk:"warn_on_no_external_source"`
}

func (r *verityServiceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service"
}

func (r *verityServiceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	provCtx, ok := req.ProviderData.(*providerContext)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *providerContext, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.provCtx = provCtx
	r.client = provCtx.client
	r.bulkOpsMgr = provCtx.bulkOpsMgr
	r.notifyOperationAdded = provCtx.NotifyOperationAdded
}

func (r *verityServiceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Service resource",
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
			"vlan": schema.Int64Attribute{
				Description: "A Value between 1 and 4096",
				Optional:    true,
			},
			"vni": schema.Int64Attribute{
				Description: "Indication of the outgoing VLAN layer 2 service. This field should not be specified when 'vni_auto_assigned_' is set to true, as the API will assign this value automatically. When specified, it represents an explicit VNI value.",
				Optional:    true,
				Computed:    true,
			},
			"vni_auto_assigned_": schema.BoolAttribute{
				Description: "Whether the VNI value should be automatically assigned by the API. When set to true, do not specify the 'vni' field in your configuration. The API will assign the VNI value, typically as VLAN + 100000.",
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
			"dhcp_server_ipv4": schema.StringAttribute{
				Description: "IPv4 address(s) of the DHCP server for service. May have up to four separated by commas.",
				Optional:    true,
			},
			"dhcp_server_ipv6": schema.StringAttribute{
				Description: "IPv6 address(s) of the DHCP server for service. May have up to four separated by commas.",
				Optional:    true,
			},
			"mtu": schema.Int64Attribute{
				Description: "MTU (Maximum Transmission Unit) - the size used by a switch to determine when large packets must be broken up for delivery.",
				Optional:    true,
			},
			"anycast_ipv4_mask": schema.StringAttribute{
				Description: "Static anycast gateway addresses(IPv4) for service",
				Optional:    true,
			},
			"anycast_ipv6_mask": schema.StringAttribute{
				Description: "Static anycast gateway addresses(IPv6) for service",
				Optional:    true,
			},
			"max_upstream_rate_mbps": schema.Int64Attribute{
				Description: "Bandwidth allocated per port in the upstream direction. (Max 10000 Mbps)",
				Optional:    true,
			},
			"max_downstream_rate_mbps": schema.Int64Attribute{
				Description: "Bandwidth allocated per port in the downstream direction. (Max 10000 Mbps)",
				Optional:    true,
			},
			"packet_priority": schema.StringAttribute{
				Description: "Priority untagged packets will be tagged with on ingress to the network.",
				Optional:    true,
			},
			"multicast_management_mode": schema.StringAttribute{
				Description: "Determines how to handle multicast packets for Service",
				Optional:    true,
			},
			"tagged_packets": schema.BoolAttribute{
				Description: "Overrides priority bits on incoming tagged packets.",
				Optional:    true,
			},
			"tls": schema.BoolAttribute{
				Description: "Is a Transparent LAN Service?",
				Optional:    true,
			},
			"allow_local_switching": schema.BoolAttribute{
				Description: "Allow Edge Devices to communicate with each other.",
				Optional:    true,
			},
			"act_as_multicast_querier": schema.BoolAttribute{
				Description: "Multicast management through IGMP requires a multicast querier.",
				Optional:    true,
			},
			"block_unknown_unicast_flood": schema.BoolAttribute{
				Description: "Block unknown unicast traffic flooding.",
				Optional:    true,
			},
			"block_downstream_dhcp_server": schema.BoolAttribute{
				Description: "Block inbound packets sent by Downstream DHCP servers",
				Optional:    true,
			},
			"is_management_service": schema.BoolAttribute{
				Description: "Denotes a Management Service",
				Optional:    true,
			},
			"use_dscp_to_p_bit_mapping_for_l3_packets_if_available": schema.BoolAttribute{
				Description: "use DSCP to p-bit Mapping for L3 packets if available",
				Optional:    true,
			},
			"allow_fast_leave": schema.BoolAttribute{
				Description: "The Fast Leave feature causes the switch to immediately remove a port from the forwarding list.",
				Optional:    true,
			},
			"mst_instance": schema.Int64Attribute{
				Description: "MST Instance ID (0-4094)",
				Optional:    true,
			},
			"policy_based_routing": schema.StringAttribute{
				Description: "Policy Based Routing",
				Optional:    true,
			},
			"policy_based_routing_ref_type_": schema.StringAttribute{
				Description: "Object type for policy_based_routing field",
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"object_properties": schema.ListNestedBlock{
				Description: "Object properties for the service",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"group": schema.StringAttribute{
							Description: "Group",
							Optional:    true,
						},
						"on_summary": schema.BoolAttribute{
							Description: "Show on the summary view",
							Optional:    true,
						},
						"warn_on_no_external_source": schema.BoolAttribute{
							Description: "Warn if there is not outbound path for service in SD-Router or a Service Port Profile",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func (r *verityServiceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityServiceResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate VNI specification when auto-assigned
	if !plan.VniAutoAssigned.IsNull() && plan.VniAutoAssigned.ValueBool() {
		if !plan.Vni.IsNull() && !plan.Vni.IsUnknown() {
			resp.Diagnostics.AddError(
				"VNI cannot be specified when auto-assigned",
				"The 'vni' field cannot be specified in the configuration when 'vni_auto_assigned_' is set to true. The API will assign this value automatically.",
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
	serviceReq := &openapi.ServicesPutRequestServiceValue{
		Name: openapi.PtrString(name),
	}

	// Handle string fields
	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "Tenant", APIField: &serviceReq.Tenant, TFValue: plan.Tenant},
		{FieldName: "TenantRefType", APIField: &serviceReq.TenantRefType, TFValue: plan.TenantRefType},
		{FieldName: "DhcpServerIpv4", APIField: &serviceReq.DhcpServerIpv4, TFValue: plan.DhcpServerIpv4},
		{FieldName: "DhcpServerIpv6", APIField: &serviceReq.DhcpServerIpv6, TFValue: plan.DhcpServerIpv6},
		{FieldName: "AnycastIpv4Mask", APIField: &serviceReq.AnycastIpv4Mask, TFValue: plan.AnycastIpv4Mask},
		{FieldName: "AnycastIpv6Mask", APIField: &serviceReq.AnycastIpv6Mask, TFValue: plan.AnycastIpv6Mask},
		{FieldName: "PacketPriority", APIField: &serviceReq.PacketPriority, TFValue: plan.PacketPriority},
		{FieldName: "MulticastManagementMode", APIField: &serviceReq.MulticastManagementMode, TFValue: plan.MulticastManagementMode},
		{FieldName: "PolicyBasedRouting", APIField: &serviceReq.PolicyBasedRouting, TFValue: plan.PolicyBasedRouting},
		{FieldName: "PolicyBasedRoutingRefType", APIField: &serviceReq.PolicyBasedRoutingRefType, TFValue: plan.PolicyBasedRoutingRefType},
	})

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &serviceReq.Enable, TFValue: plan.Enable},
		{FieldName: "TaggedPackets", APIField: &serviceReq.TaggedPackets, TFValue: plan.TaggedPackets},
		{FieldName: "Tls", APIField: &serviceReq.Tls, TFValue: plan.Tls},
		{FieldName: "AllowLocalSwitching", APIField: &serviceReq.AllowLocalSwitching, TFValue: plan.AllowLocalSwitching},
		{FieldName: "ActAsMulticastQuerier", APIField: &serviceReq.ActAsMulticastQuerier, TFValue: plan.ActAsMulticastQuerier},
		{FieldName: "BlockUnknownUnicastFlood", APIField: &serviceReq.BlockUnknownUnicastFlood, TFValue: plan.BlockUnknownUnicastFlood},
		{FieldName: "BlockDownstreamDhcpServer", APIField: &serviceReq.BlockDownstreamDhcpServer, TFValue: plan.BlockDownstreamDhcpServer},
		{FieldName: "IsManagementService", APIField: &serviceReq.IsManagementService, TFValue: plan.IsManagementService},
		{FieldName: "UseDscpToPBitMappingForL3PacketsIfAvailable", APIField: &serviceReq.UseDscpToPBitMappingForL3PacketsIfAvailable, TFValue: plan.UseDscpToPBitMappingForL3PacketsIfAvailable},
		{FieldName: "AllowFastLeave", APIField: &serviceReq.AllowFastLeave, TFValue: plan.AllowFastLeave},
	})

	// Handle nullable int64 fields
	utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
		{FieldName: "Vlan", APIField: &serviceReq.Vlan, TFValue: plan.Vlan},
		{FieldName: "Mtu", APIField: &serviceReq.Mtu, TFValue: plan.Mtu},
		{FieldName: "MaxUpstreamRateMbps", APIField: &serviceReq.MaxUpstreamRateMbps, TFValue: plan.MaxUpstreamRateMbps},
		{FieldName: "MaxDownstreamRateMbps", APIField: &serviceReq.MaxDownstreamRateMbps, TFValue: plan.MaxDownstreamRateMbps},
		{FieldName: "MstInstance", APIField: &serviceReq.MstInstance, TFValue: plan.MstInstance},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.ServicesPutRequestServiceValueObjectProperties{}
		utils.SetObjectPropertiesFields([]utils.ObjectPropertiesField{
			{Name: "Group", TFValue: op.Group, APIValue: &objProps.Group},
			{Name: "OnSummary", TFValue: op.OnSummary, APIValue: &objProps.OnSummary},
			{Name: "WarnOnNoExternalSource", TFValue: op.WarnOnNoExternalSource, APIValue: &objProps.WarnOnNoExternalSource},
		})
		serviceReq.ObjectProperties = &objProps
	}

	// Handle auto-assigned VNI logic
	if !plan.VniAutoAssigned.IsNull() && plan.VniAutoAssigned.ValueBool() {
		serviceReq.VniAutoAssigned = openapi.PtrBool(true)
		// Don't include the specific VNI in the request
	} else if !plan.Vni.IsNull() {
		// User explicitly specified a value
		vniVal := int32(plan.Vni.ValueInt64())
		serviceReq.Vni = *openapi.NewNullableInt32(&vniVal)
		if !plan.VniAutoAssigned.IsNull() {
			serviceReq.VniAutoAssigned = openapi.PtrBool(plan.VniAutoAssigned.ValueBool())
		}
	} else {
		serviceReq.Vni = *openapi.NewNullableInt32(nil)
		if !plan.VniAutoAssigned.IsNull() {
			serviceReq.VniAutoAssigned = openapi.PtrBool(plan.VniAutoAssigned.ValueBool())
		}
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "service", name, *serviceReq, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Service %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "services")

	var minState verityServiceResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if serviceData, exists := bulkMgr.GetResourceResponse("service", name); exists {
			state := populateServiceState(ctx, minState, serviceData)
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

func (r *verityServiceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityServiceResourceModel
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

	serviceName := state.Name.ValueString()

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if serviceData, exists := r.bulkOpsMgr.GetResourceResponse("service", serviceName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached service data for %s from recent operation", serviceName))
			state = populateServiceState(ctx, state, serviceData)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("service") {
		tflog.Info(ctx, fmt.Sprintf("Skipping service %s verification – trusting recent successful API operation", serviceName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching services for verification of %s", serviceName))

	type ServicesResponse struct {
		Service map[string]interface{} `json:"service"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "services", serviceName,
		func() (ServicesResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch services")
			respAPI, err := r.client.ServicesAPI.ServicesGet(ctx).Execute()
			if err != nil {
				return ServicesResponse{}, fmt.Errorf("error reading services: %v", err)
			}
			defer respAPI.Body.Close()

			var res ServicesResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return ServicesResponse{}, fmt.Errorf("failed to decode services response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d services", len(res.Service)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Service %s", serviceName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for service with name: %s", serviceName))

	serviceData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.Service,
		serviceName,
		func(data interface{}) (string, bool) {
			if service, ok := data.(map[string]interface{}); ok {
				if name, ok := service["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Service with name '%s' not found in API response", serviceName))
		resp.State.RemoveResource(ctx)
		return
	}

	serviceMap, ok := serviceData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Service Data",
			fmt.Sprintf("Service data is not in expected format for %s", serviceName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found service '%s' under API key '%s'", serviceName, actualAPIName))

	state = populateServiceState(ctx, state, serviceMap)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityServiceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityServiceResourceModel

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
	if !plan.Vni.Equal(state.Vni) &&
		!plan.Vni.IsNull() && !plan.Vni.IsUnknown() && // User is explicitly setting a value
		!plan.VniAutoAssigned.IsNull() && plan.VniAutoAssigned.ValueBool() &&
		plan.VniAutoAssigned.Equal(state.VniAutoAssigned) {
		resp.Diagnostics.AddError(
			"Cannot modify auto-assigned field",
			"The 'vni' field cannot be modified because 'vni_auto_assigned_' is set to true.",
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
	serviceReq := openapi.ServicesPutRequestServiceValue{}
	hasChanges := false

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { serviceReq.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Tenant, state.Tenant, func(v *string) { serviceReq.Tenant = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.DhcpServerIpv4, state.DhcpServerIpv4, func(v *string) { serviceReq.DhcpServerIpv4 = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.DhcpServerIpv6, state.DhcpServerIpv6, func(v *string) { serviceReq.DhcpServerIpv6 = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.AnycastIpv4Mask, state.AnycastIpv4Mask, func(v *string) { serviceReq.AnycastIpv4Mask = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.AnycastIpv6Mask, state.AnycastIpv6Mask, func(v *string) { serviceReq.AnycastIpv6Mask = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.PacketPriority, state.PacketPriority, func(v *string) { serviceReq.PacketPriority = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.MulticastManagementMode, state.MulticastManagementMode, func(v *string) { serviceReq.MulticastManagementMode = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { serviceReq.Enable = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.TaggedPackets, state.TaggedPackets, func(v *bool) { serviceReq.TaggedPackets = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.Tls, state.Tls, func(v *bool) { serviceReq.Tls = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.AllowLocalSwitching, state.AllowLocalSwitching, func(v *bool) { serviceReq.AllowLocalSwitching = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.ActAsMulticastQuerier, state.ActAsMulticastQuerier, func(v *bool) { serviceReq.ActAsMulticastQuerier = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.BlockUnknownUnicastFlood, state.BlockUnknownUnicastFlood, func(v *bool) { serviceReq.BlockUnknownUnicastFlood = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.BlockDownstreamDhcpServer, state.BlockDownstreamDhcpServer, func(v *bool) { serviceReq.BlockDownstreamDhcpServer = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.IsManagementService, state.IsManagementService, func(v *bool) { serviceReq.IsManagementService = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.UseDscpToPBitMappingForL3PacketsIfAvailable, state.UseDscpToPBitMappingForL3PacketsIfAvailable, func(v *bool) { serviceReq.UseDscpToPBitMappingForL3PacketsIfAvailable = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.AllowFastLeave, state.AllowFastLeave, func(v *bool) { serviceReq.AllowFastLeave = v }, &hasChanges)

	// Handle nullable int64 field changes
	utils.CompareAndSetNullableInt64Field(plan.Mtu, state.Mtu, func(v *openapi.NullableInt32) { serviceReq.Mtu = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.MaxUpstreamRateMbps, state.MaxUpstreamRateMbps, func(v *openapi.NullableInt32) { serviceReq.MaxUpstreamRateMbps = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.MaxDownstreamRateMbps, state.MaxDownstreamRateMbps, func(v *openapi.NullableInt32) { serviceReq.MaxDownstreamRateMbps = *v }, &hasChanges)
	utils.CompareAndSetNullableInt64Field(plan.MstInstance, state.MstInstance, func(v *openapi.NullableInt32) { serviceReq.MstInstance = *v }, &hasChanges)

	// Handle object properties
	if len(plan.ObjectProperties) > 0 && len(state.ObjectProperties) > 0 {
		objProps := openapi.ServicesPutRequestServiceValueObjectProperties{}
		op := plan.ObjectProperties[0]
		st := state.ObjectProperties[0]
		objPropsChanged := false

		utils.CompareAndSetObjectPropertiesFields([]utils.ObjectPropertiesFieldWithComparison{
			{Name: "Group", PlanValue: op.Group, StateValue: st.Group, APIValue: &objProps.Group},
			{Name: "OnSummary", PlanValue: op.OnSummary, StateValue: st.OnSummary, APIValue: &objProps.OnSummary},
			{Name: "WarnOnNoExternalSource", PlanValue: op.WarnOnNoExternalSource, StateValue: st.WarnOnNoExternalSource, APIValue: &objProps.WarnOnNoExternalSource},
		}, &objPropsChanged)

		if objPropsChanged {
			serviceReq.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	// Handle VLAN changes (preserve special handling for Unknown state)
	if !plan.Vlan.IsUnknown() && !plan.Vlan.Equal(state.Vlan) {
		utils.CompareAndSetNullableInt64Field(plan.Vlan, state.Vlan, func(v *openapi.NullableInt32) { serviceReq.Vlan = *v }, &hasChanges)
	}

	// Handle VNI and VniAutoAssigned changes
	vniChanged := !plan.Vni.IsUnknown() && !plan.Vni.Equal(state.Vni)
	vniAutoAssignedChanged := !plan.VniAutoAssigned.Equal(state.VniAutoAssigned)

	if vniChanged || vniAutoAssignedChanged {
		if vniChanged {
			if !plan.Vni.IsNull() {
				vniVal := int32(plan.Vni.ValueInt64())
				serviceReq.Vni = *openapi.NewNullableInt32(&vniVal)
			} else {
				serviceReq.Vni = *openapi.NewNullableInt32(nil)
			}
		}

		if vniAutoAssignedChanged {
			// Only send vni_auto_assigned_ if the user has explicitly specified it in their configuration
			var config verityServiceResourceModel
			userSpecifiedVniAutoAssigned := false
			if !req.Config.Raw.IsNull() {
				if err := req.Config.Get(ctx, &config); err == nil {
					userSpecifiedVniAutoAssigned = !config.VniAutoAssigned.IsNull()
				}
			}

			if userSpecifiedVniAutoAssigned {
				serviceReq.VniAutoAssigned = openapi.PtrBool(plan.VniAutoAssigned.ValueBool())

				// Special case: When changing from auto-assigned (true) to manual (false),
				// the API requires both vni_auto_assigned_ and vni fields to be sent.
				// Otherwise, the vni_auto_assigned_ change will be ignored by the API.
				if !state.VniAutoAssigned.IsNull() && state.VniAutoAssigned.ValueBool() &&
					!plan.VniAutoAssigned.ValueBool() {
					// Changing from auto-assigned=true to auto-assigned=false
					// Must include VNI value in the request for the change to take effect
					if !plan.Vni.IsNull() {
						vniVal := int32(plan.Vni.ValueInt64())
						serviceReq.Vni = *openapi.NewNullableInt32(&vniVal)
					} else if !state.Vni.IsNull() {
						// Use current state VNI if plan doesn't specify one
						vniVal := int32(state.Vni.ValueInt64())
						serviceReq.Vni = *openapi.NewNullableInt32(&vniVal)
					}
				}
			}
		} else if vniChanged {
			// VNI changed but VniAutoAssigned didn't change
			// Send the auto-assigned flag to maintain consistency with API
			if !plan.VniAutoAssigned.IsNull() {
				serviceReq.VniAutoAssigned = openapi.PtrBool(plan.VniAutoAssigned.ValueBool())
			} else if !state.VniAutoAssigned.IsNull() {
				serviceReq.VniAutoAssigned = openapi.PtrBool(state.VniAutoAssigned.ValueBool())
			} else {
				serviceReq.VniAutoAssigned = openapi.PtrBool(false)
			}
		}

		hasChanges = true
	}

	// Handle tenant and tenant_ref_type_ fields using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.Tenant, state.Tenant, plan.TenantRefType, state.TenantRefType,
		func(v *string) { serviceReq.Tenant = v },
		func(v *string) { serviceReq.TenantRefType = v },
		"tenant", "tenant_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	if !utils.HandleOneRefTypeSupported(
		plan.PolicyBasedRouting, state.PolicyBasedRouting, plan.PolicyBasedRoutingRefType, state.PolicyBasedRoutingRefType,
		func(val *string) { serviceReq.PolicyBasedRouting = val },
		func(val *string) { serviceReq.PolicyBasedRoutingRefType = val },
		"policy_based_routing", "policy_based_routing_ref_type_",
		&hasChanges, &resp.Diagnostics,
	) {
		return
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "service", name, serviceReq, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Service %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "services")

	var minState verityServiceResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if serviceData, exists := bulkMgr.GetResourceResponse("service", name); exists {
			state := populateServiceState(ctx, minState, serviceData)
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

func (r *verityServiceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityServiceResourceModel
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "service", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Service %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "services")
	resp.State.RemoveResource(ctx)
}

func (r *verityServiceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func populateServiceState(ctx context.Context, state verityServiceResourceModel, serviceData map[string]interface{}) verityServiceResourceModel {
	state.Name = utils.MapStringFromAPI(serviceData["name"])

	// Int fields
	state.Vlan = utils.MapInt64FromAPI(serviceData["vlan"])
	state.Vni = utils.MapInt64FromAPI(serviceData["vni"])
	state.Mtu = utils.MapInt64FromAPI(serviceData["mtu"])
	state.MaxUpstreamRateMbps = utils.MapInt64FromAPI(serviceData["max_upstream_rate_mbps"])
	state.MaxDownstreamRateMbps = utils.MapInt64FromAPI(serviceData["max_downstream_rate_mbps"])
	state.MstInstance = utils.MapInt64FromAPI(serviceData["mst_instance"])

	// Bool fields
	state.Enable = utils.MapBoolFromAPI(serviceData["enable"])
	state.VniAutoAssigned = utils.MapBoolFromAPI(serviceData["vni_auto_assigned_"])
	state.TaggedPackets = utils.MapBoolFromAPI(serviceData["tagged_packets"])
	state.Tls = utils.MapBoolFromAPI(serviceData["tls"])
	state.AllowLocalSwitching = utils.MapBoolFromAPI(serviceData["allow_local_switching"])
	state.ActAsMulticastQuerier = utils.MapBoolFromAPI(serviceData["act_as_multicast_querier"])
	state.BlockUnknownUnicastFlood = utils.MapBoolFromAPI(serviceData["block_unknown_unicast_flood"])
	state.BlockDownstreamDhcpServer = utils.MapBoolFromAPI(serviceData["block_downstream_dhcp_server"])
	state.IsManagementService = utils.MapBoolFromAPI(serviceData["is_management_service"])
	state.UseDscpToPBitMappingForL3PacketsIfAvailable = utils.MapBoolFromAPI(serviceData["use_dscp_to_p_bit_mapping_for_l3_packets_if_available"])
	state.AllowFastLeave = utils.MapBoolFromAPI(serviceData["allow_fast_leave"])

	// String fields
	state.Tenant = utils.MapStringFromAPI(serviceData["tenant"])
	state.TenantRefType = utils.MapStringFromAPI(serviceData["tenant_ref_type_"])
	state.DhcpServerIpv4 = utils.MapStringFromAPI(serviceData["dhcp_server_ipv4"])
	state.DhcpServerIpv6 = utils.MapStringFromAPI(serviceData["dhcp_server_ipv6"])
	state.AnycastIpv4Mask = utils.MapStringFromAPI(serviceData["anycast_ipv4_mask"])
	state.AnycastIpv6Mask = utils.MapStringFromAPI(serviceData["anycast_ipv6_mask"])
	state.PacketPriority = utils.MapStringFromAPI(serviceData["packet_priority"])
	state.MulticastManagementMode = utils.MapStringFromAPI(serviceData["multicast_management_mode"])
	state.PolicyBasedRouting = utils.MapStringFromAPI(serviceData["policy_based_routing"])
	state.PolicyBasedRoutingRefType = utils.MapStringFromAPI(serviceData["policy_based_routing_ref_type_"])

	// Object properties
	if op, ok := serviceData["object_properties"].(map[string]interface{}); ok {
		objProps := verityServiceObjectPropertiesModel{
			Group:                  utils.MapStringFromAPI(op["group"]),
			OnSummary:              utils.MapBoolFromAPI(op["on_summary"]),
			WarnOnNoExternalSource: utils.MapBoolFromAPI(op["warn_on_no_external_source"]),
		}
		state.ObjectProperties = []verityServiceObjectPropertiesModel{objProps}
	} else {
		state.ObjectProperties = nil
	}

	return state
}

func (r *verityServiceResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Skip modification if we're deleting the resource
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityServiceResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate VNI specification in configuration when auto-assigned
	// Check the actual configuration, not the plan
	var config verityServiceResourceModel
	if !req.Config.Raw.IsNull() {
		resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if !config.VniAutoAssigned.IsNull() && config.VniAutoAssigned.ValueBool() {
			if !config.Vni.IsNull() && !config.Vni.IsUnknown() {
				resp.Diagnostics.AddError(
					"VNI cannot be specified when auto-assigned",
					"The 'vni' field cannot be specified in the configuration when 'vni_auto_assigned_' is set to true. The API will assign this value automatically.",
				)
				return
			}
		}
	}

	// For new resources (where state is null)
	if req.State.Raw.IsNull() {
		if !plan.VniAutoAssigned.IsNull() && plan.VniAutoAssigned.ValueBool() {
			resp.Plan.SetAttribute(ctx, path.Root("vni"), types.Int64Unknown())
		}
		return
	}

	var state verityServiceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Handle VNI behavior based on auto-assignment and VLAN changes
	if !plan.VniAutoAssigned.IsNull() && plan.VniAutoAssigned.ValueBool() {
		if !plan.VniAutoAssigned.Equal(state.VniAutoAssigned) {
			// vni_auto_assigned_ is changing to true, API will assign VNI based on VLAN
			// Only mark VNI as unknown, VLAN stays as configured
			resp.Plan.SetAttribute(ctx, path.Root("vni"), types.Int64Unknown())
			resp.Diagnostics.AddWarning(
				"VNI will be assigned by the API",
				"The 'vni' field will be automatically assigned by the API because 'vni_auto_assigned_' is being set to true. The API will assign VNI based on the VLAN value.",
			)
		} else if !plan.Vlan.Equal(state.Vlan) {
			// VLAN is changing, so VNI will be auto-updated by API
			resp.Plan.SetAttribute(ctx, path.Root("vni"), types.Int64Unknown())
			resp.Diagnostics.AddWarning(
				"VNI will be updated by the API",
				"The 'vni' field will be automatically updated by the API because 'vni_auto_assigned_' is set to true and VLAN is changing.",
			)
		} else if !plan.Vni.Equal(state.Vni) {
			// User tried to change VNI but it's auto-assigned
			resp.Diagnostics.AddWarning(
				"Ignoring vni changes with auto-assignment enabled",
				"The 'vni' field changes will be ignored because 'vni_auto_assigned_' is set to true. The API will assign this value automatically.",
			)
			// Keep the current state value to suppress the diff
			if !state.Vni.IsNull() {
				resp.Plan.SetAttribute(ctx, path.Root("vni"), state.Vni)
			}
		}
	} else if !plan.Vlan.Equal(state.Vlan) && plan.Vni.Equal(state.Vni) && (plan.VniAutoAssigned.IsNull() || !plan.VniAutoAssigned.ValueBool()) {
		// VLAN is changing and VNI hasn't been explicitly changed by the user
		// Only mark VNI as unknown when vni_auto_assigned_ is false
		// AND the user hasn't explicitly specified VNI in their configuration

		// Check if user explicitly specified VNI in config
		var config verityServiceResourceModel
		hasExplicitVni := false
		if !req.Config.Raw.IsNull() {
			if err := req.Config.Get(ctx, &config); err == nil {
				hasExplicitVni = !config.Vni.IsNull() && !config.Vni.IsUnknown()
			}
		}

		// Only mark VNI as unknown if user hasn't explicitly specified it
		if !hasExplicitVni {
			resp.Plan.SetAttribute(ctx, path.Root("vni"), types.Int64Unknown())
			tflog.Info(ctx, "Marking VNI as unknown due to VLAN change - API will determine the actual value")
		}
	}
}
