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
	bulkOpsMgr           *utils.BulkOperationManager
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

	// Handle int64 fields
	utils.SetInt64Fields([]utils.Int64FieldMapping{
		{FieldName: "MaxUpstreamRateMbps", APIField: &serviceReq.MaxUpstreamRateMbps, TFValue: plan.MaxUpstreamRateMbps},
		{FieldName: "MaxDownstreamRateMbps", APIField: &serviceReq.MaxDownstreamRateMbps, TFValue: plan.MaxDownstreamRateMbps},
		{FieldName: "MstInstance", APIField: &serviceReq.MstInstance, TFValue: plan.MstInstance},
	})

	// Handle nullable int64 fields
	utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
		{FieldName: "Vlan", APIField: &serviceReq.Vlan, TFValue: plan.Vlan},
		{FieldName: "Mtu", APIField: &serviceReq.Mtu, TFValue: plan.Mtu},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.ServicesPutRequestServiceValueObjectProperties{}
		if !op.Group.IsNull() {
			objProps.Group = openapi.PtrString(op.Group.ValueString())
		} else {
			objProps.Group = nil
		}
		if !op.OnSummary.IsNull() {
			objProps.OnSummary = openapi.PtrBool(op.OnSummary.ValueBool())
		} else {
			objProps.OnSummary = nil
		}
		if !op.WarnOnNoExternalSource.IsNull() {
			objProps.WarnOnNoExternalSource = openapi.PtrBool(op.WarnOnNoExternalSource.ValueBool())
		} else {
			objProps.WarnOnNoExternalSource = nil
		}
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

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "service", name, *serviceReq, &resp.Diagnostics)
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
			// Use the cached data with plan values as fallback
			state := populateServiceState(ctx, minState, serviceData, &plan)
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
			state = populateServiceState(ctx, state, serviceData, nil)
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

	state = populateServiceState(ctx, state, serviceMap, nil)
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

	// Handle int64 field changes
	utils.CompareAndSetInt64Field(plan.MaxUpstreamRateMbps, state.MaxUpstreamRateMbps, func(v *int32) { serviceReq.MaxUpstreamRateMbps = v }, &hasChanges)
	utils.CompareAndSetInt64Field(plan.MaxDownstreamRateMbps, state.MaxDownstreamRateMbps, func(v *int32) { serviceReq.MaxDownstreamRateMbps = v }, &hasChanges)
	utils.CompareAndSetInt64Field(plan.MstInstance, state.MstInstance, func(v *int32) { serviceReq.MstInstance = v }, &hasChanges)

	// Handle nullable int64 field changes
	utils.CompareAndSetNullableInt64Field(plan.Mtu, state.Mtu, func(v *openapi.NullableInt32) { serviceReq.Mtu = *v }, &hasChanges)

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		if len(state.ObjectProperties) == 0 ||
			!plan.ObjectProperties[0].Group.Equal(state.ObjectProperties[0].Group) ||
			!plan.ObjectProperties[0].OnSummary.Equal(state.ObjectProperties[0].OnSummary) ||
			!plan.ObjectProperties[0].WarnOnNoExternalSource.Equal(state.ObjectProperties[0].WarnOnNoExternalSource) {
			objProps := openapi.ServicesPutRequestServiceValueObjectProperties{}
			if !plan.ObjectProperties[0].Group.IsNull() {
				objProps.Group = openapi.PtrString(plan.ObjectProperties[0].Group.ValueString())
			} else {
				objProps.Group = nil
			}
			if !plan.ObjectProperties[0].OnSummary.IsNull() {
				objProps.OnSummary = openapi.PtrBool(plan.ObjectProperties[0].OnSummary.ValueBool())
			} else {
				objProps.OnSummary = nil
			}
			if !plan.ObjectProperties[0].WarnOnNoExternalSource.IsNull() {
				objProps.WarnOnNoExternalSource = openapi.PtrBool(plan.ObjectProperties[0].WarnOnNoExternalSource.ValueBool())
			} else {
				objProps.WarnOnNoExternalSource = nil
			}
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

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "service", name, serviceReq, &resp.Diagnostics)
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
			// Use the cached data from the API response with plan values as fallback
			state := populateServiceState(ctx, minState, serviceData, &plan)
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

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "service", name, nil, &resp.Diagnostics)
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

func populateServiceState(ctx context.Context, state verityServiceResourceModel, serviceData map[string]interface{}, plan *verityServiceResourceModel) verityServiceResourceModel {
	state.Name = types.StringValue(fmt.Sprintf("%v", serviceData["name"]))

	// For each field, check if it's in the API response first,
	// if not: use plan value (if plan provided and not null), otherwise preserve current state value

	if val, ok := serviceData["enable"].(bool); ok {
		state.Enable = types.BoolValue(val)
	} else if plan != nil && !plan.Enable.IsNull() {
		state.Enable = plan.Enable
	}

	if op, ok := serviceData["object_properties"].(map[string]interface{}); ok {
		objProps := verityServiceObjectPropertiesModel{}

		if group, exists := op["group"]; exists {
			objProps.Group = types.StringValue(fmt.Sprintf("%v", group))
		} else if len(state.ObjectProperties) > 0 {
			objProps.Group = state.ObjectProperties[0].Group
		} else {
			objProps.Group = types.StringNull()
		}

		if onSummary, exists := op["on_summary"].(bool); exists {
			objProps.OnSummary = types.BoolValue(onSummary)
		} else if plan != nil && len(plan.ObjectProperties) > 0 && !plan.ObjectProperties[0].OnSummary.IsNull() {
			objProps.OnSummary = plan.ObjectProperties[0].OnSummary
		} else if len(state.ObjectProperties) > 0 {
			objProps.OnSummary = state.ObjectProperties[0].OnSummary
		} else {
			objProps.OnSummary = types.BoolNull()
		}

		if warnOnNoExternal, exists := op["warn_on_no_external_source"].(bool); exists {
			objProps.WarnOnNoExternalSource = types.BoolValue(warnOnNoExternal)
		} else if plan != nil && len(plan.ObjectProperties) > 0 && !plan.ObjectProperties[0].WarnOnNoExternalSource.IsNull() {
			objProps.WarnOnNoExternalSource = plan.ObjectProperties[0].WarnOnNoExternalSource
		} else if len(state.ObjectProperties) > 0 {
			objProps.WarnOnNoExternalSource = state.ObjectProperties[0].WarnOnNoExternalSource
		} else {
			objProps.WarnOnNoExternalSource = types.BoolNull()
		}

		state.ObjectProperties = []verityServiceObjectPropertiesModel{objProps}
	} else if plan != nil && len(plan.ObjectProperties) > 0 {
		state.ObjectProperties = plan.ObjectProperties
	}

	if val, ok := serviceData["vlan"]; ok {
		switch v := val.(type) {
		case float64:
			state.Vlan = types.Int64Value(int64(v))
		case int:
			state.Vlan = types.Int64Value(int64(v))
		default:
			if plan != nil && !plan.Vlan.IsNull() && !plan.Vlan.IsUnknown() {
				state.Vlan = plan.Vlan
			}
		}
	} else if plan != nil && !plan.Vlan.IsNull() && !plan.Vlan.IsUnknown() {
		state.Vlan = plan.Vlan
	}

	if val, ok := serviceData["vni"]; ok {
		if val == nil {
			state.Vni = types.Int64Null()
		} else {
			switch v := val.(type) {
			case float64:
				state.Vni = types.Int64Value(int64(v))
			case int:
				state.Vni = types.Int64Value(int64(v))
			case int32:
				state.Vni = types.Int64Value(int64(v))
			default:
				if plan != nil && !plan.Vni.IsNull() && !plan.Vni.IsUnknown() {
					state.Vni = plan.Vni
				}
			}
		}
	} else if plan != nil && !plan.Vni.IsNull() && !plan.Vni.IsUnknown() {
		state.Vni = plan.Vni
	}

	if val, ok := serviceData["vni_auto_assigned_"].(bool); ok {
		state.VniAutoAssigned = types.BoolValue(val)
	} else if plan != nil && !plan.VniAutoAssigned.IsNull() {
		state.VniAutoAssigned = plan.VniAutoAssigned
	} else {
		state.VniAutoAssigned = types.BoolNull()
	}

	if val, ok := serviceData["tenant"].(string); ok {
		state.Tenant = types.StringValue(val)
	} else if plan != nil && !plan.Tenant.IsNull() {
		state.Tenant = plan.Tenant
	} else {
		state.Tenant = types.StringNull()
	}

	if val, ok := serviceData["tenant_ref_type_"].(string); ok {
		state.TenantRefType = types.StringValue(val)
	} else if plan != nil && !plan.TenantRefType.IsNull() {
		state.TenantRefType = plan.TenantRefType
	} else {
		state.TenantRefType = types.StringNull()
	}

	if val, ok := serviceData["dhcp_server_ipv4"].(string); ok {
		state.DhcpServerIpv4 = types.StringValue(val)
	} else if plan != nil && !plan.DhcpServerIpv4.IsNull() {
		state.DhcpServerIpv4 = plan.DhcpServerIpv4
	} else {
		state.DhcpServerIpv4 = types.StringNull()
	}

	if val, ok := serviceData["dhcp_server_ipv6"].(string); ok {
		state.DhcpServerIpv6 = types.StringValue(val)
	} else if plan != nil && !plan.DhcpServerIpv6.IsNull() {
		state.DhcpServerIpv6 = plan.DhcpServerIpv6
	} else {
		state.DhcpServerIpv6 = types.StringNull()
	}

	if val, ok := serviceData["mtu"]; ok {
		switch v := val.(type) {
		case float64:
			state.Mtu = types.Int64Value(int64(v))
		case int:
			state.Mtu = types.Int64Value(int64(v))
		default:
			if plan != nil && !plan.Mtu.IsNull() {
				state.Mtu = plan.Mtu
			} else {
				state.Mtu = types.Int64Null()
			}
		}
	} else if plan != nil && !plan.Mtu.IsNull() {
		state.Mtu = plan.Mtu
	} else {
		state.Mtu = types.Int64Null()
	}

	if val, ok := serviceData["anycast_ipv4_mask"].(string); ok {
		state.AnycastIpv4Mask = types.StringValue(val)
	} else if plan != nil && !plan.AnycastIpv4Mask.IsNull() {
		state.AnycastIpv4Mask = plan.AnycastIpv4Mask
	} else {
		state.AnycastIpv4Mask = types.StringNull()
	}

	if val, ok := serviceData["anycast_ipv6_mask"].(string); ok {
		state.AnycastIpv6Mask = types.StringValue(val)
	} else if plan != nil && !plan.AnycastIpv6Mask.IsNull() {
		state.AnycastIpv6Mask = plan.AnycastIpv6Mask
	} else {
		state.AnycastIpv6Mask = types.StringNull()
	}

	if val, ok := serviceData["max_upstream_rate_mbps"]; ok {
		switch v := val.(type) {
		case float64:
			state.MaxUpstreamRateMbps = types.Int64Value(int64(v))
		case int:
			state.MaxUpstreamRateMbps = types.Int64Value(int64(v))
		case int32:
			state.MaxUpstreamRateMbps = types.Int64Value(int64(v))
		default:
			if plan != nil && !plan.MaxUpstreamRateMbps.IsNull() {
				state.MaxUpstreamRateMbps = plan.MaxUpstreamRateMbps
			} else {
				state.MaxUpstreamRateMbps = types.Int64Null()
			}
		}
	} else if plan != nil && !plan.MaxUpstreamRateMbps.IsNull() {
		state.MaxUpstreamRateMbps = plan.MaxUpstreamRateMbps
	} else {
		state.MaxUpstreamRateMbps = types.Int64Null()
	}

	if val, ok := serviceData["max_downstream_rate_mbps"]; ok {
		switch v := val.(type) {
		case float64:
			state.MaxDownstreamRateMbps = types.Int64Value(int64(v))
		case int:
			state.MaxDownstreamRateMbps = types.Int64Value(int64(v))
		case int32:
			state.MaxDownstreamRateMbps = types.Int64Value(int64(v))
		default:
			if plan != nil && !plan.MaxDownstreamRateMbps.IsNull() {
				state.MaxDownstreamRateMbps = plan.MaxDownstreamRateMbps
			} else {
				state.MaxDownstreamRateMbps = types.Int64Null()
			}
		}
	} else if plan != nil && !plan.MaxDownstreamRateMbps.IsNull() {
		state.MaxDownstreamRateMbps = plan.MaxDownstreamRateMbps
	} else {
		state.MaxDownstreamRateMbps = types.Int64Null()
	}

	if val, ok := serviceData["packet_priority"].(string); ok {
		state.PacketPriority = types.StringValue(val)
	} else if plan != nil && !plan.PacketPriority.IsNull() {
		state.PacketPriority = plan.PacketPriority
	} else {
		state.PacketPriority = types.StringNull()
	}

	if val, ok := serviceData["multicast_management_mode"].(string); ok {
		state.MulticastManagementMode = types.StringValue(val)
	} else if plan != nil && !plan.MulticastManagementMode.IsNull() {
		state.MulticastManagementMode = plan.MulticastManagementMode
	} else {
		state.MulticastManagementMode = types.StringNull()
	}

	if val, ok := serviceData["tagged_packets"].(bool); ok {
		state.TaggedPackets = types.BoolValue(val)
	} else if plan != nil && !plan.TaggedPackets.IsNull() {
		state.TaggedPackets = plan.TaggedPackets
	} else {
		state.TaggedPackets = types.BoolNull()
	}

	if val, ok := serviceData["tls"].(bool); ok {
		state.Tls = types.BoolValue(val)
	} else if plan != nil && !plan.Tls.IsNull() {
		state.Tls = plan.Tls
	} else {
		state.Tls = types.BoolNull()
	}

	if val, ok := serviceData["allow_local_switching"].(bool); ok {
		state.AllowLocalSwitching = types.BoolValue(val)
	} else if plan != nil && !plan.AllowLocalSwitching.IsNull() {
		state.AllowLocalSwitching = plan.AllowLocalSwitching
	} else {
		state.AllowLocalSwitching = types.BoolNull()
	}

	if val, ok := serviceData["act_as_multicast_querier"].(bool); ok {
		state.ActAsMulticastQuerier = types.BoolValue(val)
	} else if plan != nil && !plan.ActAsMulticastQuerier.IsNull() {
		state.ActAsMulticastQuerier = plan.ActAsMulticastQuerier
	} else {
		state.ActAsMulticastQuerier = types.BoolNull()
	}

	if val, ok := serviceData["block_unknown_unicast_flood"].(bool); ok {
		state.BlockUnknownUnicastFlood = types.BoolValue(val)
	} else if plan != nil && !plan.BlockUnknownUnicastFlood.IsNull() {
		state.BlockUnknownUnicastFlood = plan.BlockUnknownUnicastFlood
	} else {
		state.BlockUnknownUnicastFlood = types.BoolNull()
	}

	if val, ok := serviceData["block_downstream_dhcp_server"].(bool); ok {
		state.BlockDownstreamDhcpServer = types.BoolValue(val)
	} else if plan != nil && !plan.BlockDownstreamDhcpServer.IsNull() {
		state.BlockDownstreamDhcpServer = plan.BlockDownstreamDhcpServer
	} else {
		state.BlockDownstreamDhcpServer = types.BoolNull()
	}

	if val, ok := serviceData["is_management_service"].(bool); ok {
		state.IsManagementService = types.BoolValue(val)
	} else if plan != nil && !plan.IsManagementService.IsNull() {
		state.IsManagementService = plan.IsManagementService
	} else {
		state.IsManagementService = types.BoolNull()
	}

	if val, ok := serviceData["use_dscp_to_p_bit_mapping_for_l3_packets_if_available"].(bool); ok {
		state.UseDscpToPBitMappingForL3PacketsIfAvailable = types.BoolValue(val)
	} else if plan != nil && !plan.UseDscpToPBitMappingForL3PacketsIfAvailable.IsNull() {
		state.UseDscpToPBitMappingForL3PacketsIfAvailable = plan.UseDscpToPBitMappingForL3PacketsIfAvailable
	} else {
		state.UseDscpToPBitMappingForL3PacketsIfAvailable = types.BoolNull()
	}

	if val, ok := serviceData["allow_fast_leave"].(bool); ok {
		state.AllowFastLeave = types.BoolValue(val)
	} else if plan != nil && !plan.AllowFastLeave.IsNull() {
		state.AllowFastLeave = plan.AllowFastLeave
	} else {
		state.AllowFastLeave = types.BoolNull()
	}

	if val, ok := serviceData["mst_instance"]; ok {
		switch v := val.(type) {
		case float64:
			state.MstInstance = types.Int64Value(int64(v))
		case int:
			state.MstInstance = types.Int64Value(int64(v))
		case int32:
			state.MstInstance = types.Int64Value(int64(v))
		default:
			if plan != nil && !plan.MstInstance.IsNull() {
				state.MstInstance = plan.MstInstance
			} else {
				state.MstInstance = types.Int64Null()
			}
		}
	} else if plan != nil && !plan.MstInstance.IsNull() {
		state.MstInstance = plan.MstInstance
	} else {
		state.MstInstance = types.Int64Null()
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
