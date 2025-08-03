package provider

import (
	"context"
	"encoding/json"
	"fmt"
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
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"vlan": schema.Int64Attribute{
				Description: "A Value between 1 and 4096",
				Optional:    true,
			},
			"vni": schema.Int64Attribute{
				Description: "Indication of the outgoing VLAN layer 2 service",
				Optional:    true,
				Computed:    true,
			},
			"vni_auto_assigned_": schema.BoolAttribute{
				Description: "Whether or not the value in vni field has been automatically assigned or not.",
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

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Authenticate",
			fmt.Sprintf("Error authenticating with API: %s", err),
		)
		return
	}

	name := plan.Name.ValueString()

	serviceReq := openapi.NewServicesPutRequestServiceValue()
	serviceReq.Name = openapi.PtrString(name)
	if !plan.Enable.IsNull() {
		serviceReq.Enable = openapi.PtrBool(plan.Enable.ValueBool())
	}

	if len(plan.ObjectProperties) > 0 {
		objProps := openapi.ServicesPutRequestServiceValueObjectProperties{}
		op := plan.ObjectProperties[0]
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
	} else {
		serviceReq.ObjectProperties = nil
	}

	if !plan.Vlan.IsNull() {
		vlanVal := int32(plan.Vlan.ValueInt64())
		serviceReq.Vlan = *openapi.NewNullableInt32(&vlanVal)
	} else {
		serviceReq.Vlan = *openapi.NewNullableInt32(nil)
	}
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
	if !plan.Tenant.IsNull() {
		serviceReq.Tenant = openapi.PtrString(plan.Tenant.ValueString())
	}
	if !plan.TenantRefType.IsNull() {
		serviceReq.TenantRefType = openapi.PtrString(plan.TenantRefType.ValueString())
	}
	if !plan.DhcpServerIpv4.IsNull() {
		serviceReq.DhcpServerIpv4 = openapi.PtrString(plan.DhcpServerIpv4.ValueString())
	}
	if !plan.DhcpServerIpv6.IsNull() {
		serviceReq.DhcpServerIpv6 = openapi.PtrString(plan.DhcpServerIpv6.ValueString())
	}
	if !plan.Mtu.IsNull() {
		mtuVal := int32(plan.Mtu.ValueInt64())
		serviceReq.Mtu = *openapi.NewNullableInt32(&mtuVal)
	} else {
		serviceReq.Mtu = *openapi.NewNullableInt32(nil)
	}
	if !plan.AnycastIpv4Mask.IsNull() {
		serviceReq.AnycastIpv4Mask = openapi.PtrString(plan.AnycastIpv4Mask.ValueString())
	}
	if !plan.AnycastIpv6Mask.IsNull() {
		serviceReq.AnycastIpv6Mask = openapi.PtrString(plan.AnycastIpv6Mask.ValueString())
	}
	if !plan.MaxUpstreamRateMbps.IsNull() {
		serviceReq.MaxUpstreamRateMbps = openapi.PtrInt32(int32(plan.MaxUpstreamRateMbps.ValueInt64()))
	}
	if !plan.MaxDownstreamRateMbps.IsNull() {
		serviceReq.MaxDownstreamRateMbps = openapi.PtrInt32(int32(plan.MaxDownstreamRateMbps.ValueInt64()))
	}
	if !plan.PacketPriority.IsNull() {
		serviceReq.PacketPriority = openapi.PtrString(plan.PacketPriority.ValueString())
	}
	if !plan.MulticastManagementMode.IsNull() {
		serviceReq.MulticastManagementMode = openapi.PtrString(plan.MulticastManagementMode.ValueString())
	}
	if !plan.TaggedPackets.IsNull() {
		serviceReq.TaggedPackets = openapi.PtrBool(plan.TaggedPackets.ValueBool())
	}
	if !plan.Tls.IsNull() {
		serviceReq.Tls = openapi.PtrBool(plan.Tls.ValueBool())
	}
	if !plan.AllowLocalSwitching.IsNull() {
		serviceReq.AllowLocalSwitching = openapi.PtrBool(plan.AllowLocalSwitching.ValueBool())
	}
	if !plan.ActAsMulticastQuerier.IsNull() {
		serviceReq.ActAsMulticastQuerier = openapi.PtrBool(plan.ActAsMulticastQuerier.ValueBool())
	}
	if !plan.BlockUnknownUnicastFlood.IsNull() {
		serviceReq.BlockUnknownUnicastFlood = openapi.PtrBool(plan.BlockUnknownUnicastFlood.ValueBool())
	}
	if !plan.BlockDownstreamDhcpServer.IsNull() {
		serviceReq.BlockDownstreamDhcpServer = openapi.PtrBool(plan.BlockDownstreamDhcpServer.ValueBool())
	}
	if !plan.IsManagementService.IsNull() {
		serviceReq.IsManagementService = openapi.PtrBool(plan.IsManagementService.ValueBool())
	}
	if !plan.UseDscpToPBitMappingForL3PacketsIfAvailable.IsNull() {
		serviceReq.UseDscpToPBitMappingForL3PacketsIfAvailable = openapi.PtrBool(plan.UseDscpToPBitMappingForL3PacketsIfAvailable.ValueBool())
	}
	if !plan.AllowFastLeave.IsNull() {
		serviceReq.AllowFastLeave = openapi.PtrBool(plan.AllowFastLeave.ValueBool())
	}
	if !plan.MstInstance.IsNull() {
		serviceReq.MstInstance = openapi.PtrInt32(int32(plan.MstInstance.ValueInt64()))
	}

	provCtx := r.provCtx
	bulkMgr := provCtx.bulkOpsMgr
	operationID := bulkMgr.AddPut(ctx, "service", name, *serviceReq)

	provCtx.NotifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for service creation operation %s to complete", operationID))
	if err := bulkMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Create Service %s", name))...,
		)
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
			fmt.Sprintf("Error authenticating with API: %v", err),
		)
		return
	}

	tflog.Debug(ctx, "Reading service resource")

	provCtx := r.provCtx
	bulkOpsMgr := provCtx.bulkOpsMgr
	serviceName := state.Name.ValueString()

	var serviceData map[string]interface{}
	var exists bool

	if bulkOpsMgr != nil {
		serviceData, exists = bulkOpsMgr.GetResourceResponse("service", serviceName)
		if exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached service data for %s from recent operation", serviceName))
			state = populateServiceState(ctx, state, serviceData, nil)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if bulkOpsMgr != nil && bulkOpsMgr.HasPendingOrRecentOperations("service") {
		tflog.Info(ctx, fmt.Sprintf("Skipping service %s verification - trusting recent successful API operation", serviceName))
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("No recent service operations found, performing normal verification for %s", serviceName))

	type ServicesResponse struct {
		Service map[string]map[string]interface{} `json:"service"`
	}

	var result ServicesResponse
	var err error
	maxRetries := 3

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch services on attempt %d, retrying in %v", attempt, sleepTime))
			time.Sleep(sleepTime)
		}

		servicesData, fetchErr := getCachedResponse(ctx, provCtx, "services", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch services")
			apiReq := provCtx.client.ServicesAPI.ServicesGet(ctx)
			apiResp, err := apiReq.Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading service: %v", err)
			}
			defer apiResp.Body.Close()

			var res ServicesResponse
			if err := json.NewDecoder(apiResp.Body).Decode(&res); err != nil {
				return nil, fmt.Errorf("failed to decode services response: %v", err)
			}
			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched services data with %d services", len(res.Service)))
			return res, nil
		})

		if fetchErr == nil {
			result = servicesData.(ServicesResponse)
			break
		}
		err = fetchErr
	}

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Service %s", serviceName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for service with ID: %s", serviceName))
	if data, ok := result.Service[serviceName]; ok {
		serviceData = data
		exists = true
		tflog.Debug(ctx, fmt.Sprintf("Found service directly by ID: %s", serviceName))
	} else {
		for apiName, s := range result.Service {
			if name, ok := s["name"].(string); ok && name == serviceName {
				serviceData = s
				serviceName = apiName
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found service with name '%s' under API key '%s'", name, apiName))
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Service with ID '%s' not found in API response", serviceName))
		resp.State.RemoveResource(ctx)
		return
	}

	state = populateServiceState(ctx, state, serviceData, nil)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityServiceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityServiceResourceModel
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
	serviceReq := &openapi.ServicesPutRequestServiceValue{}
	hasChanges := false

	objPropsChanged := false
	if len(plan.ObjectProperties) > 0 && len(state.ObjectProperties) > 0 {
		op_plan := plan.ObjectProperties[0]
		op_state := state.ObjectProperties[0]
		if !op_plan.Group.Equal(op_state.Group) || !op_plan.OnSummary.Equal(op_state.OnSummary) || !op_plan.WarnOnNoExternalSource.Equal(op_state.WarnOnNoExternalSource) {
			objPropsChanged = true
		}
	} else if len(plan.ObjectProperties) != len(state.ObjectProperties) {
		objPropsChanged = true
	}

	if objPropsChanged {
		if len(plan.ObjectProperties) > 0 {
			objProps := openapi.ServicesPutRequestServiceValueObjectProperties{}
			op := plan.ObjectProperties[0]
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
		} else {
			serviceReq.ObjectProperties = nil
		}
		hasChanges = true
	}

	if !plan.Enable.Equal(state.Enable) {
		serviceReq.Enable = openapi.PtrBool(plan.Enable.ValueBool())
		hasChanges = true
	}

	if !plan.Vlan.Equal(state.Vlan) {
		if !plan.Vlan.IsNull() {
			vlanVal := int32(plan.Vlan.ValueInt64())
			serviceReq.Vlan = *openapi.NewNullableInt32(&vlanVal)
		} else {
			serviceReq.Vlan = *openapi.NewNullableInt32(nil)
		}
		hasChanges = true
	}

	if !plan.Vni.Equal(state.Vni) {
		if !plan.Vni.IsNull() {
			vniVal := int32(plan.Vni.ValueInt64())
			serviceReq.Vni = *openapi.NewNullableInt32(&vniVal)
		} else {
			serviceReq.Vni = *openapi.NewNullableInt32(nil)
		}
		hasChanges = true
	}

	if !plan.VniAutoAssigned.Equal(state.VniAutoAssigned) {
		serviceReq.VniAutoAssigned = openapi.PtrBool(plan.VniAutoAssigned.ValueBool())
		hasChanges = true
	}

	tenantChanged := !plan.Tenant.Equal(state.Tenant)
	tenantRefTypeChanged := !plan.TenantRefType.Equal(state.TenantRefType)

	if tenantChanged || tenantRefTypeChanged {
		// Validate using one ref type supported rules
		if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
			plan.Tenant, plan.TenantRefType,
			"tenant", "tenant_ref_type_",
			tenantChanged, tenantRefTypeChanged) {
			return
		}

		// Only send the base field if only it changed
		if tenantChanged && !tenantRefTypeChanged {
			// Just send the base field
			if !plan.Tenant.IsNull() && plan.Tenant.ValueString() != "" {
				serviceReq.Tenant = openapi.PtrString(plan.Tenant.ValueString())
			} else {
				serviceReq.Tenant = openapi.PtrString("")
			}
			hasChanges = true
		} else if tenantRefTypeChanged {
			// Send both fields
			if !plan.Tenant.IsNull() && plan.Tenant.ValueString() != "" {
				serviceReq.Tenant = openapi.PtrString(plan.Tenant.ValueString())
			} else {
				serviceReq.Tenant = openapi.PtrString("")
			}

			if !plan.TenantRefType.IsNull() && plan.TenantRefType.ValueString() != "" {
				serviceReq.TenantRefType = openapi.PtrString(plan.TenantRefType.ValueString())
			} else {
				serviceReq.TenantRefType = openapi.PtrString("")
			}
			hasChanges = true
		}
	}

	if !plan.DhcpServerIpv4.Equal(state.DhcpServerIpv4) {
		serviceReq.DhcpServerIpv4 = openapi.PtrString(plan.DhcpServerIpv4.ValueString())
		hasChanges = true
	}

	if !plan.DhcpServerIpv6.Equal(state.DhcpServerIpv6) {
		serviceReq.DhcpServerIpv6 = openapi.PtrString(plan.DhcpServerIpv6.ValueString())
		hasChanges = true
	}

	if !plan.Mtu.Equal(state.Mtu) {
		if !plan.Mtu.IsNull() {
			mtuVal := int32(plan.Mtu.ValueInt64())
			serviceReq.Mtu = *openapi.NewNullableInt32(&mtuVal)
		} else {
			serviceReq.Mtu = *openapi.NewNullableInt32(nil)
		}
		hasChanges = true
	}

	if !plan.AnycastIpv4Mask.Equal(state.AnycastIpv4Mask) {
		serviceReq.AnycastIpv4Mask = openapi.PtrString(plan.AnycastIpv4Mask.ValueString())
		hasChanges = true
	}

	if !plan.AnycastIpv6Mask.Equal(state.AnycastIpv6Mask) {
		serviceReq.AnycastIpv6Mask = openapi.PtrString(plan.AnycastIpv6Mask.ValueString())
		hasChanges = true
	}

	if !plan.MaxUpstreamRateMbps.Equal(state.MaxUpstreamRateMbps) {
		if !plan.MaxUpstreamRateMbps.IsNull() {
			serviceReq.MaxUpstreamRateMbps = openapi.PtrInt32(int32(plan.MaxUpstreamRateMbps.ValueInt64()))
		}
		hasChanges = true
	}

	if !plan.MaxDownstreamRateMbps.Equal(state.MaxDownstreamRateMbps) {
		if !plan.MaxDownstreamRateMbps.IsNull() {
			serviceReq.MaxDownstreamRateMbps = openapi.PtrInt32(int32(plan.MaxDownstreamRateMbps.ValueInt64()))
		}
		hasChanges = true
	}

	if !plan.PacketPriority.Equal(state.PacketPriority) {
		serviceReq.PacketPriority = openapi.PtrString(plan.PacketPriority.ValueString())
		hasChanges = true
	}

	if !plan.MulticastManagementMode.Equal(state.MulticastManagementMode) {
		serviceReq.MulticastManagementMode = openapi.PtrString(plan.MulticastManagementMode.ValueString())
		hasChanges = true
	}

	if !plan.TaggedPackets.Equal(state.TaggedPackets) {
		serviceReq.TaggedPackets = openapi.PtrBool(plan.TaggedPackets.ValueBool())
		hasChanges = true
	}

	if !plan.Tls.Equal(state.Tls) {
		serviceReq.Tls = openapi.PtrBool(plan.Tls.ValueBool())
		hasChanges = true
	}

	if !plan.AllowLocalSwitching.Equal(state.AllowLocalSwitching) {
		serviceReq.AllowLocalSwitching = openapi.PtrBool(plan.AllowLocalSwitching.ValueBool())
		hasChanges = true
	}

	if !plan.ActAsMulticastQuerier.Equal(state.ActAsMulticastQuerier) {
		serviceReq.ActAsMulticastQuerier = openapi.PtrBool(plan.ActAsMulticastQuerier.ValueBool())
		hasChanges = true
	}

	if !plan.BlockUnknownUnicastFlood.Equal(state.BlockUnknownUnicastFlood) {
		serviceReq.BlockUnknownUnicastFlood = openapi.PtrBool(plan.BlockUnknownUnicastFlood.ValueBool())
		hasChanges = true
	}

	if !plan.BlockDownstreamDhcpServer.Equal(state.BlockDownstreamDhcpServer) {
		serviceReq.BlockDownstreamDhcpServer = openapi.PtrBool(plan.BlockDownstreamDhcpServer.ValueBool())
		hasChanges = true
	}

	if !plan.IsManagementService.Equal(state.IsManagementService) {
		serviceReq.IsManagementService = openapi.PtrBool(plan.IsManagementService.ValueBool())
		hasChanges = true
	}

	if !plan.UseDscpToPBitMappingForL3PacketsIfAvailable.Equal(state.UseDscpToPBitMappingForL3PacketsIfAvailable) {
		serviceReq.UseDscpToPBitMappingForL3PacketsIfAvailable = openapi.PtrBool(plan.UseDscpToPBitMappingForL3PacketsIfAvailable.ValueBool())
		hasChanges = true
	}

	if !plan.AllowFastLeave.Equal(state.AllowFastLeave) {
		serviceReq.AllowFastLeave = openapi.PtrBool(plan.AllowFastLeave.ValueBool())
		hasChanges = true
	}

	if !plan.MstInstance.Equal(state.MstInstance) {
		serviceReq.MstInstance = openapi.PtrInt32(int32(plan.MstInstance.ValueInt64()))
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	bulkOpsMgr := r.provCtx.bulkOpsMgr
	operationID := bulkOpsMgr.AddPatch(ctx, "service", name, *serviceReq)
	r.provCtx.NotifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for service update operation %s to complete", operationID))
	if err := bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update Service %s", name))...,
		)
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
			fmt.Sprintf("Error authenticating with API: %v", err),
		)
		return
	}

	name := state.Name.ValueString()
	bulkOpsMgr := r.provCtx.bulkOpsMgr
	operationID := bulkOpsMgr.AddDelete(ctx, "service", name)
	r.provCtx.NotifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for service deletion operation %s to complete", operationID))
	if err := bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Delete Service %s", name))...,
		)
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
			if plan != nil && !plan.Vlan.IsNull() {
				state.Vlan = plan.Vlan
			}
		}
	} else if plan != nil && !plan.Vlan.IsNull() {
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
				if plan != nil && !plan.Vni.IsNull() {
					state.Vni = plan.Vni
				}
			}
		}
	} else if plan != nil && !plan.Vni.IsNull() {
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

	// For new resources (where state is null)
	if req.State.Raw.IsNull() {
		if !plan.VniAutoAssigned.IsNull() && plan.VniAutoAssigned.ValueBool() {
			resp.Diagnostics.AddWarning(
				"VNI will be assigned by the API",
				"The 'vni' field value in your configuration will be ignored because 'vni_auto_assigned_' is set to true. The API will assign this value automatically.",
			)
		}
		return
	}

	var state verityServiceResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Only show warning and suppress diff if auto-assignment is enabled AND the field is actually changing
	if !plan.VniAutoAssigned.IsNull() && plan.VniAutoAssigned.ValueBool() && !plan.Vni.Equal(state.Vni) {
		resp.Diagnostics.AddWarning(
			"Ignoring vni changes with auto-assignment enabled",
			"The 'vni' field changes will be ignored because 'vni_auto_assigned_' is set to true. The API will assign this value automatically.",
		)

		// Use current state value to suppress the diff
		if !state.Vni.IsNull() {
			resp.Plan.SetAttribute(ctx, path.Root("vni"), state.Vni)
		}
	}
}
