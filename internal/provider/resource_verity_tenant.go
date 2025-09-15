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
	_ resource.Resource                = &verityTenantResource{}
	_ resource.ResourceWithConfigure   = &verityTenantResource{}
	_ resource.ResourceWithImportState = &verityTenantResource{}
	_ resource.ResourceWithModifyPlan  = &verityTenantResource{}
)

func NewVerityTenantResource() resource.Resource {
	return &verityTenantResource{}
}

type verityTenantResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *utils.BulkOperationManager
	notifyOperationAdded func()
}

type verityTenantResourceModel struct {
	Name                       types.String                        `tfsdk:"name"`
	Enable                     types.Bool                          `tfsdk:"enable"`
	ObjectProperties           []verityTenantObjectPropertiesModel `tfsdk:"object_properties"`
	Layer3Vni                  types.Int64                         `tfsdk:"layer_3_vni"`
	Layer3VniAutoAssigned      types.Bool                          `tfsdk:"layer_3_vni_auto_assigned_"`
	Layer3Vlan                 types.Int64                         `tfsdk:"layer_3_vlan"`
	Layer3VlanAutoAssigned     types.Bool                          `tfsdk:"layer_3_vlan_auto_assigned_"`
	DhcpRelaySourceIpv4sSubnet types.String                        `tfsdk:"dhcp_relay_source_ipv4s_subnet"`
	DhcpRelaySourceIpv6sSubnet types.String                        `tfsdk:"dhcp_relay_source_ipv6s_subnet"`
	RouteDistinguisher         types.String                        `tfsdk:"route_distinguisher"`
	RouteTargetImport          types.String                        `tfsdk:"route_target_import"`
	RouteTargetExport          types.String                        `tfsdk:"route_target_export"`
	ImportRouteMap             types.String                        `tfsdk:"import_route_map"`
	ImportRouteMapRefType      types.String                        `tfsdk:"import_route_map_ref_type_"`
	ExportRouteMap             types.String                        `tfsdk:"export_route_map"`
	ExportRouteMapRefType      types.String                        `tfsdk:"export_route_map_ref_type_"`
	VrfName                    types.String                        `tfsdk:"vrf_name"`
	VrfNameAutoAssigned        types.Bool                          `tfsdk:"vrf_name_auto_assigned_"`
	RouteTenants               []verityTenantRouteTenantModel      `tfsdk:"route_tenants"`
	DefaultOriginate           types.Bool                          `tfsdk:"default_originate"`
}

type verityTenantObjectPropertiesModel struct {
	Group types.String `tfsdk:"group"`
}

type verityTenantRouteTenantModel struct {
	Enable types.Bool   `tfsdk:"enable"`
	Tenant types.String `tfsdk:"tenant"`
	Index  types.Int64  `tfsdk:"index"`
}

func (r *verityTenantResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tenant"
}

func (r *verityTenantResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityTenantResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Tenant resource",
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
			"layer_3_vni": schema.Int64Attribute{
				Description: "VNI value used to transport traffic between services of a Tenant. This field should not be specified when 'layer_3_vni_auto_assigned_' is set to true, as the API will assign this value automatically.",
				Optional:    true,
				Computed:    true,
			},
			"layer_3_vni_auto_assigned_": schema.BoolAttribute{
				Description: "Whether the Layer 3 VNI value should be automatically assigned by the API. When set to true, do not specify the 'layer_3_vni' field in your configuration.",
				Optional:    true,
			},
			"layer_3_vlan": schema.Int64Attribute{
				Description: "VLAN value used to transport traffic between services of a Tenant. This field should not be specified when 'layer_3_vlan_auto_assigned_' is set to true, as the API will assign this value automatically.",
				Optional:    true,
				Computed:    true,
			},
			"layer_3_vlan_auto_assigned_": schema.BoolAttribute{
				Description: "Whether the Layer 3 VLAN value should be automatically assigned by the API. When set to true, do not specify the 'layer_3_vlan' field in your configuration.",
				Optional:    true,
			},
			"dhcp_relay_source_ipv4s_subnet": schema.StringAttribute{
				Description: "Range of IPv4 addresses (represented in IPv4 subnet format) used to configure the source IP of each DHCP Relay on each switch that this Tenant is provisioned on.",
				Optional:    true,
			},
			"dhcp_relay_source_ipv6s_subnet": schema.StringAttribute{
				Description: "Range of IPv6 addresses (represented in IPv6 subnet format) used to configure the source IP of each DHCP Relay on each switch that this Tenant is provisioned on.",
				Optional:    true,
			},
			"route_distinguisher": schema.StringAttribute{
				Description: "Route Distinguisher (BGP Community) for uniqueness among identical routes",
				Optional:    true,
			},
			"route_target_import": schema.StringAttribute{
				Description: "Route-target to attach while importing routes into the tenant",
				Optional:    true,
			},
			"route_target_export": schema.StringAttribute{
				Description: "Route-target to attach while exporting routes from the tenant",
				Optional:    true,
			},
			"import_route_map": schema.StringAttribute{
				Description: "Route-map applied to routes imported into the tenant",
				Optional:    true,
			},
			"import_route_map_ref_type_": schema.StringAttribute{
				Description: "Object type for import_route_map field",
				Optional:    true,
			},
			"export_route_map": schema.StringAttribute{
				Description: "Route-map applied to routes exported from the tenant",
				Optional:    true,
			},
			"export_route_map_ref_type_": schema.StringAttribute{
				Description: "Object type for export_route_map field",
				Optional:    true,
			},
			"vrf_name": schema.StringAttribute{
				Description: "Virtual Routing and Forwarding instance name. This field should not be specified when 'vrf_name_auto_assigned_' is set to true, as the API will assign this value automatically.",
				Optional:    true,
				Computed:    true,
			},
			"vrf_name_auto_assigned_": schema.BoolAttribute{
				Description: "Whether the VRF name should be automatically assigned by the API. When set to true, do not specify the 'vrf_name' field in your configuration.",
				Optional:    true,
			},
			"default_originate": schema.BoolAttribute{
				Description: "Enables a leaf switch to originate IPv4 default type-5 EVPN routes across the switching fabric.",
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"object_properties": schema.ListNestedBlock{
				Description: "Object properties for the tenant",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"group": schema.StringAttribute{
							Description: "Group",
							Optional:    true,
						},
					},
				},
			},
			"route_tenants": schema.ListNestedBlock{
				Description: "Route tenants configuration",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"enable": schema.BoolAttribute{
							Description: "Enable",
							Optional:    true,
						},
						"tenant": schema.StringAttribute{
							Description: "Tenant",
							Optional:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func (r *verityTenantResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityTenantResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate auto-assigned field specifications
	if !plan.Layer3VniAutoAssigned.IsNull() && plan.Layer3VniAutoAssigned.ValueBool() {
		if !plan.Layer3Vni.IsNull() && !plan.Layer3Vni.IsUnknown() {
			resp.Diagnostics.AddError(
				"Layer 3 VNI cannot be specified when auto-assigned",
				"The 'layer_3_vni' field cannot be specified in the configuration when 'layer_3_vni_auto_assigned_' is set to true. The API will assign this value automatically.",
			)
			return
		}
	}

	if !plan.Layer3VlanAutoAssigned.IsNull() && plan.Layer3VlanAutoAssigned.ValueBool() {
		if !plan.Layer3Vlan.IsNull() && !plan.Layer3Vlan.IsUnknown() {
			resp.Diagnostics.AddError(
				"Layer 3 VLAN cannot be specified when auto-assigned",
				"The 'layer_3_vlan' field cannot be specified in the configuration when 'layer_3_vlan_auto_assigned_' is set to true. The API will assign this value automatically.",
			)
			return
		}
	}

	if !plan.VrfNameAutoAssigned.IsNull() && plan.VrfNameAutoAssigned.ValueBool() {
		if !plan.VrfName.IsNull() && !plan.VrfName.IsUnknown() && plan.VrfName.ValueString() != "" {
			resp.Diagnostics.AddError(
				"VRF name cannot be specified when auto-assigned",
				"The 'vrf_name' field cannot be specified in the configuration when 'vrf_name_auto_assigned_' is set to true. The API will assign this value automatically.",
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
	tenantReq := &openapi.TenantsPutRequestTenantValue{
		Name:         openapi.PtrString(name),
		RouteTenants: []openapi.TenantsPutRequestTenantValueRouteTenantsInner{},
	}

	if !plan.Enable.IsNull() {
		tenantReq.Enable = openapi.PtrBool(plan.Enable.ValueBool())
	}
	if !plan.DhcpRelaySourceIpv4sSubnet.IsNull() {
		tenantReq.DhcpRelaySourceIpv4sSubnet = openapi.PtrString(plan.DhcpRelaySourceIpv4sSubnet.ValueString())
	}
	if !plan.DhcpRelaySourceIpv6sSubnet.IsNull() {
		tenantReq.DhcpRelaySourceIpv6sSubnet = openapi.PtrString(plan.DhcpRelaySourceIpv6sSubnet.ValueString())
	}
	if !plan.RouteDistinguisher.IsNull() {
		tenantReq.RouteDistinguisher = openapi.PtrString(plan.RouteDistinguisher.ValueString())
	}
	if !plan.RouteTargetImport.IsNull() {
		tenantReq.RouteTargetImport = openapi.PtrString(plan.RouteTargetImport.ValueString())
	}
	if !plan.RouteTargetExport.IsNull() {
		tenantReq.RouteTargetExport = openapi.PtrString(plan.RouteTargetExport.ValueString())
	}
	if !plan.ImportRouteMap.IsNull() {
		tenantReq.ImportRouteMap = openapi.PtrString(plan.ImportRouteMap.ValueString())
	}
	if !plan.ImportRouteMapRefType.IsNull() {
		tenantReq.ImportRouteMapRefType = openapi.PtrString(plan.ImportRouteMapRefType.ValueString())
	}
	if !plan.ExportRouteMap.IsNull() {
		tenantReq.ExportRouteMap = openapi.PtrString(plan.ExportRouteMap.ValueString())
	}
	if !plan.ExportRouteMapRefType.IsNull() {
		tenantReq.ExportRouteMapRefType = openapi.PtrString(plan.ExportRouteMapRefType.ValueString())
	}

	if len(plan.ObjectProperties) > 0 {
		objProps := openapi.GatewayprofilesPutRequestGatewayProfileValueObjectProperties{}
		if !plan.ObjectProperties[0].Group.IsNull() {
			objProps.Group = openapi.PtrString(plan.ObjectProperties[0].Group.ValueString())
		} else {
			objProps.Group = nil
		}
		tenantReq.ObjectProperties = &objProps
	} else {
		tenantReq.ObjectProperties = nil
	}

	if !plan.Layer3VniAutoAssigned.IsNull() && plan.Layer3VniAutoAssigned.ValueBool() {
		tenantReq.Layer3VniAutoAssigned = openapi.PtrBool(true)
		// Don't include the specific VNI in the request
	} else if !plan.Layer3Vni.IsNull() {
		// User explicitly specified a value
		val := int32(plan.Layer3Vni.ValueInt64())
		tenantReq.Layer3Vni = *openapi.NewNullableInt32(&val)
		if !plan.Layer3VniAutoAssigned.IsNull() {
			tenantReq.Layer3VniAutoAssigned = openapi.PtrBool(plan.Layer3VniAutoAssigned.ValueBool())
		}
	} else {
		tenantReq.Layer3Vni = *openapi.NewNullableInt32(nil)
		if !plan.Layer3VniAutoAssigned.IsNull() {
			tenantReq.Layer3VniAutoAssigned = openapi.PtrBool(plan.Layer3VniAutoAssigned.ValueBool())
		}
	}
	if !plan.Layer3VlanAutoAssigned.IsNull() && plan.Layer3VlanAutoAssigned.ValueBool() {
		tenantReq.Layer3VlanAutoAssigned = openapi.PtrBool(true)
		// Don't include the specific VLAN in the request
	} else if !plan.Layer3Vlan.IsNull() {
		// User explicitly specified a value
		val := int32(plan.Layer3Vlan.ValueInt64())
		tenantReq.Layer3Vlan = *openapi.NewNullableInt32(&val)
		if !plan.Layer3VlanAutoAssigned.IsNull() {
			tenantReq.Layer3VlanAutoAssigned = openapi.PtrBool(plan.Layer3VlanAutoAssigned.ValueBool())
		}
	} else {
		tenantReq.Layer3Vlan = *openapi.NewNullableInt32(nil)
		if !plan.Layer3VlanAutoAssigned.IsNull() {
			tenantReq.Layer3VlanAutoAssigned = openapi.PtrBool(plan.Layer3VlanAutoAssigned.ValueBool())
		}
	}
	if !plan.VrfNameAutoAssigned.IsNull() && plan.VrfNameAutoAssigned.ValueBool() {
		tenantReq.VrfNameAutoAssigned = openapi.PtrBool(true)
		// Don't include the specific VRF name in the request
	} else if !plan.VrfName.IsNull() {
		// User explicitly specified a value
		tenantReq.VrfName = openapi.PtrString(plan.VrfName.ValueString())
		if !plan.VrfNameAutoAssigned.IsNull() {
			tenantReq.VrfNameAutoAssigned = openapi.PtrBool(plan.VrfNameAutoAssigned.ValueBool())
		}
	} else {
		if !plan.VrfNameAutoAssigned.IsNull() {
			tenantReq.VrfNameAutoAssigned = openapi.PtrBool(plan.VrfNameAutoAssigned.ValueBool())
		}
	}
	if !plan.DefaultOriginate.IsNull() {
		tenantReq.DefaultOriginate = openapi.PtrBool(plan.DefaultOriginate.ValueBool())
	}

	if len(plan.RouteTenants) > 0 {
		for _, rt := range plan.RouteTenants {
			tenant := openapi.TenantsPutRequestTenantValueRouteTenantsInner{
				Enable: openapi.PtrBool(rt.Enable.ValueBool()),
				Tenant: openapi.PtrString(rt.Tenant.ValueString()),
			}
			if !rt.Index.IsNull() {
				tenant.Index = openapi.PtrInt32(int32(rt.Index.ValueInt64()))
			}
			tenantReq.RouteTenants = append(tenantReq.RouteTenants, tenant)
		}
	}

	provCtx := r.provCtx
	bulkOpsMgr := provCtx.bulkOpsMgr
	operationID := bulkOpsMgr.AddPut(ctx, "tenant", name, *tenantReq)

	provCtx.NotifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for tenant creation operation %s to complete", operationID))
	if err := bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Create Tenant %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Tenant %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "tenants")

	var minState verityTenantResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if tenantData, exists := bulkMgr.GetResourceResponse("tenant", name); exists {
			// Use the cached data with plan values as fallback
			state := populateTenantState(ctx, minState, tenantData, &plan)
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

func (r *verityTenantResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityTenantResourceModel
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

	tflog.Debug(ctx, "Reading tenant resource")

	provCtx := r.provCtx
	bulkOpsMgr := provCtx.bulkOpsMgr
	tenantName := state.Name.ValueString()

	var tenantData map[string]interface{}
	var exists bool

	if bulkOpsMgr != nil {
		tenantData, exists = bulkOpsMgr.GetResourceResponse("tenant", tenantName)
		if exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached tenant data for %s from recent operation", tenantName))
			state = populateTenantState(ctx, state, tenantData, nil)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if bulkOpsMgr != nil && bulkOpsMgr.HasPendingOrRecentOperations("tenant") {
		tflog.Info(ctx, fmt.Sprintf("Skipping tenant %s verification - trusting recent successful API operation", tenantName))
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("No recent tenant operations found, performing normal verification for %s", tenantName))

	type TenantsResponse struct {
		Tenant map[string]map[string]interface{} `json:"tenant"`
	}

	var result TenantsResponse
	var err error
	maxRetries := 3

	for attempt := 0; attempt < maxRetries; attempt++ {
		tenantsData, fetchErr := getCachedResponse(ctx, provCtx, "tenants", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch tenants")
			apiResp, err := provCtx.client.TenantsAPI.TenantsGet(ctx).Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading tenant: %v", err)
			}
			defer apiResp.Body.Close()

			var res TenantsResponse
			if err := json.NewDecoder(apiResp.Body).Decode(&res); err != nil {
				return nil, fmt.Errorf("failed to decode tenants response: %v", err)
			}
			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d tenants from API", len(res.Tenant)))
			return res, nil
		})
		if fetchErr != nil {
			err = fetchErr
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch tenants on attempt %d, retrying in %v", attempt+1, sleepTime))
			time.Sleep(sleepTime)
			continue
		}
		result = tenantsData.(TenantsResponse)
		break
	}

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Tenant %s", tenantName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for tenant with ID: %s", tenantName))
	exists = false

	if data, ok := result.Tenant[tenantName]; ok {
		tenantData = data
		exists = true
		tflog.Debug(ctx, fmt.Sprintf("Found tenant directly by ID: %s", tenantName))
	} else {
		for apiName, t := range result.Tenant {
			if nameVal, ok := t["name"].(string); ok && nameVal == tenantName {
				tenantData = t
				tenantName = apiName
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found tenant with name '%s' under API key '%s'", nameVal, apiName))
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Tenant with ID '%s' not found in API response", tenantName))
		resp.State.RemoveResource(ctx)
		return
	}

	state = populateTenantState(ctx, state, tenantData, nil)
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityTenantResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityTenantResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate auto-assigned fields - these checks prevent ineffective API calls
	// Only error if the auto-assigned flag is enabled AND the user is explicitly setting a value
	// AND the auto-assigned flag itself is not changing (which would be a valid operation)
	// Don't error if the field is unknown (computed during plan recalculation)
	if !plan.Layer3Vni.IsUnknown() && !plan.Layer3Vni.Equal(state.Layer3Vni) && !plan.Layer3VniAutoAssigned.IsNull() && plan.Layer3VniAutoAssigned.ValueBool() {
		resp.Diagnostics.AddError(
			"Cannot modify auto-assigned field",
			"The 'layer_3_vni' field cannot be modified because 'layer_3_vni_auto_assigned_' is set to true.",
		)
		return
	}

	if !plan.Layer3Vlan.IsUnknown() && !plan.Layer3Vlan.Equal(state.Layer3Vlan) && !plan.Layer3VlanAutoAssigned.IsNull() && plan.Layer3VlanAutoAssigned.ValueBool() {
		resp.Diagnostics.AddError(
			"Cannot modify auto-assigned field",
			"The 'layer_3_vlan' field cannot be modified because 'layer_3_vlan_auto_assigned_' is set to true.",
		)
		return
	}

	if !plan.VrfName.IsUnknown() && !plan.VrfName.Equal(state.VrfName) && !plan.VrfNameAutoAssigned.IsNull() && plan.VrfNameAutoAssigned.ValueBool() {
		resp.Diagnostics.AddError(
			"Cannot modify auto-assigned field",
			"The 'vrf_name' field cannot be modified because 'vrf_name_auto_assigned_' is set to true.",
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

	tenantReq := &openapi.TenantsPutRequestTenantValue{}
	hasChanges := false

	if len(plan.ObjectProperties) > 0 {
		if len(state.ObjectProperties) == 0 || !plan.ObjectProperties[0].Group.Equal(state.ObjectProperties[0].Group) {
			objProps := openapi.GatewayprofilesPutRequestGatewayProfileValueObjectProperties{}
			if !plan.ObjectProperties[0].Group.IsNull() {
				objProps.Group = openapi.PtrString(plan.ObjectProperties[0].Group.ValueString())
			} else {
				objProps.Group = nil
			}
			tenantReq.ObjectProperties = &objProps
			hasChanges = true
		}
	} else {
		tenantReq.ObjectProperties = nil
	}

	if !plan.Enable.Equal(state.Enable) {
		tenantReq.Enable = openapi.PtrBool(plan.Enable.ValueBool())
		hasChanges = true
	}

	// Handle Layer3Vni and Layer3VniAutoAssigned changes in a coordinated way
	layer3VniChanged := !plan.Layer3Vni.IsUnknown() && !plan.Layer3Vni.Equal(state.Layer3Vni)
	layer3VniAutoAssignedChanged := !plan.Layer3VniAutoAssigned.Equal(state.Layer3VniAutoAssigned)

	if layer3VniChanged || layer3VniAutoAssignedChanged {
		// Handle Layer3Vni field changes
		if layer3VniChanged {
			if !plan.Layer3Vni.IsNull() {
				val := int32(plan.Layer3Vni.ValueInt64())
				tenantReq.Layer3Vni = *openapi.NewNullableInt32(&val)
			} else {
				tenantReq.Layer3Vni = *openapi.NewNullableInt32(nil)
			}
		}

		// Handle Layer3VniAutoAssigned field changes
		if layer3VniAutoAssignedChanged {
			// Only send layer_3_vni_auto_assigned_ if the user has explicitly specified it in their configuration
			var config verityTenantResourceModel
			userSpecifiedLayer3VniAutoAssigned := false
			if !req.Config.Raw.IsNull() {
				if err := req.Config.Get(ctx, &config); err == nil {
					userSpecifiedLayer3VniAutoAssigned = !config.Layer3VniAutoAssigned.IsNull()
				}
			}

			if userSpecifiedLayer3VniAutoAssigned {
				tenantReq.Layer3VniAutoAssigned = openapi.PtrBool(plan.Layer3VniAutoAssigned.ValueBool())

				// Special case: When changing from auto-assigned (true) to manual (false),
				// the API requires both layer_3_vni_auto_assigned_ and layer_3_vni fields to be sent.
				if !state.Layer3VniAutoAssigned.IsNull() && state.Layer3VniAutoAssigned.ValueBool() &&
					!plan.Layer3VniAutoAssigned.ValueBool() {
					// Changing from auto-assigned=true to auto-assigned=false
					// Must include Layer3Vni value in the request for the change to take effect
					if !plan.Layer3Vni.IsNull() {
						val := int32(plan.Layer3Vni.ValueInt64())
						tenantReq.Layer3Vni = *openapi.NewNullableInt32(&val)
					} else if !state.Layer3Vni.IsNull() {
						// Use current state Layer3Vni if plan doesn't specify one
						val := int32(state.Layer3Vni.ValueInt64())
						tenantReq.Layer3Vni = *openapi.NewNullableInt32(&val)
					}
				}
			}
		} else if layer3VniChanged {
			// Layer3Vni changed but Layer3VniAutoAssigned didn't change
			// Send the auto-assigned flag to maintain consistency with API
			if !plan.Layer3VniAutoAssigned.IsNull() {
				tenantReq.Layer3VniAutoAssigned = openapi.PtrBool(plan.Layer3VniAutoAssigned.ValueBool())
			} else if !state.Layer3VniAutoAssigned.IsNull() {
				tenantReq.Layer3VniAutoAssigned = openapi.PtrBool(state.Layer3VniAutoAssigned.ValueBool())
			} else {
				tenantReq.Layer3VniAutoAssigned = openapi.PtrBool(false)
			}
		}

		hasChanges = true
	}

	// Handle Layer3Vlan and Layer3VlanAutoAssigned changes in a coordinated way
	layer3VlanChanged := !plan.Layer3Vlan.IsUnknown() && !plan.Layer3Vlan.Equal(state.Layer3Vlan)
	layer3VlanAutoAssignedChanged := !plan.Layer3VlanAutoAssigned.Equal(state.Layer3VlanAutoAssigned)

	if layer3VlanChanged || layer3VlanAutoAssignedChanged {
		// Handle Layer3Vlan field changes
		if layer3VlanChanged {
			if !plan.Layer3Vlan.IsNull() {
				val := int32(plan.Layer3Vlan.ValueInt64())
				tenantReq.Layer3Vlan = *openapi.NewNullableInt32(&val)
			} else {
				tenantReq.Layer3Vlan = *openapi.NewNullableInt32(nil)
			}
		}

		// Handle Layer3VlanAutoAssigned field changes
		if layer3VlanAutoAssignedChanged {
			// Only send layer_3_vlan_auto_assigned_ if the user has explicitly specified it in their configuration
			var config verityTenantResourceModel
			userSpecifiedLayer3VlanAutoAssigned := false
			if !req.Config.Raw.IsNull() {
				if err := req.Config.Get(ctx, &config); err == nil {
					userSpecifiedLayer3VlanAutoAssigned = !config.Layer3VlanAutoAssigned.IsNull()
				}
			}

			if userSpecifiedLayer3VlanAutoAssigned {
				tenantReq.Layer3VlanAutoAssigned = openapi.PtrBool(plan.Layer3VlanAutoAssigned.ValueBool())

				// Special case: When changing from auto-assigned (true) to manual (false),
				// the API requires both layer_3_vlan_auto_assigned_ and layer_3_vlan fields to be sent.
				if !state.Layer3VlanAutoAssigned.IsNull() && state.Layer3VlanAutoAssigned.ValueBool() &&
					!plan.Layer3VlanAutoAssigned.ValueBool() {
					// Changing from auto-assigned=true to auto-assigned=false
					// Must include Layer3Vlan value in the request for the change to take effect
					if !plan.Layer3Vlan.IsNull() {
						val := int32(plan.Layer3Vlan.ValueInt64())
						tenantReq.Layer3Vlan = *openapi.NewNullableInt32(&val)
					} else if !state.Layer3Vlan.IsNull() {
						// Use current state Layer3Vlan if plan doesn't specify one
						val := int32(state.Layer3Vlan.ValueInt64())
						tenantReq.Layer3Vlan = *openapi.NewNullableInt32(&val)
					}
				}
			}
		} else if layer3VlanChanged {
			// Layer3Vlan changed but Layer3VlanAutoAssigned didn't change
			// Send the auto-assigned flag to maintain consistency with API
			if !plan.Layer3VlanAutoAssigned.IsNull() {
				tenantReq.Layer3VlanAutoAssigned = openapi.PtrBool(plan.Layer3VlanAutoAssigned.ValueBool())
			} else if !state.Layer3VlanAutoAssigned.IsNull() {
				tenantReq.Layer3VlanAutoAssigned = openapi.PtrBool(state.Layer3VlanAutoAssigned.ValueBool())
			} else {
				tenantReq.Layer3VlanAutoAssigned = openapi.PtrBool(false)
			}
		}

		hasChanges = true
	}

	// Handle VrfName and VrfNameAutoAssigned changes in a coordinated way
	vrfNameChanged := !plan.VrfName.IsUnknown() && !plan.VrfName.Equal(state.VrfName)
	vrfNameAutoAssignedChanged := !plan.VrfNameAutoAssigned.Equal(state.VrfNameAutoAssigned)

	if vrfNameChanged || vrfNameAutoAssignedChanged {
		// Handle VrfName field changes
		if vrfNameChanged {
			if !plan.VrfName.IsNull() && plan.VrfName.ValueString() != "" {
				tenantReq.VrfName = openapi.PtrString(plan.VrfName.ValueString())
			} else {
				tenantReq.VrfName = openapi.PtrString("")
			}
		}

		// Handle VrfNameAutoAssigned field changes
		if vrfNameAutoAssignedChanged {
			// Only send vrf_name_auto_assigned_ if the user has explicitly specified it in their configuration
			var config verityTenantResourceModel
			userSpecifiedVrfNameAutoAssigned := false
			if !req.Config.Raw.IsNull() {
				if err := req.Config.Get(ctx, &config); err == nil {
					userSpecifiedVrfNameAutoAssigned = !config.VrfNameAutoAssigned.IsNull()
				}
			}

			if userSpecifiedVrfNameAutoAssigned {
				tenantReq.VrfNameAutoAssigned = openapi.PtrBool(plan.VrfNameAutoAssigned.ValueBool())

				// Special case: When changing from auto-assigned (true) to manual (false),
				// the API requires both vrf_name_auto_assigned_ and vrf_name fields to be sent.
				if !state.VrfNameAutoAssigned.IsNull() && state.VrfNameAutoAssigned.ValueBool() &&
					!plan.VrfNameAutoAssigned.ValueBool() {
					// Changing from auto-assigned=true to auto-assigned=false
					// Must include VrfName value in the request for the change to take effect
					if !plan.VrfName.IsNull() && plan.VrfName.ValueString() != "" {
						tenantReq.VrfName = openapi.PtrString(plan.VrfName.ValueString())
					} else if !state.VrfName.IsNull() && state.VrfName.ValueString() != "" {
						// Use current state VrfName if plan doesn't specify one
						tenantReq.VrfName = openapi.PtrString(state.VrfName.ValueString())
					}
				}
			}
		} else if vrfNameChanged {
			// VrfName changed but VrfNameAutoAssigned didn't change
			// Send the auto-assigned flag to maintain consistency with API
			if !plan.VrfNameAutoAssigned.IsNull() {
				tenantReq.VrfNameAutoAssigned = openapi.PtrBool(plan.VrfNameAutoAssigned.ValueBool())
			} else if !state.VrfNameAutoAssigned.IsNull() {
				tenantReq.VrfNameAutoAssigned = openapi.PtrBool(state.VrfNameAutoAssigned.ValueBool())
			} else {
				tenantReq.VrfNameAutoAssigned = openapi.PtrBool(false)
			}
		}

		hasChanges = true
	}

	if !plan.DhcpRelaySourceIpv4sSubnet.Equal(state.DhcpRelaySourceIpv4sSubnet) {
		if !plan.DhcpRelaySourceIpv4sSubnet.IsNull() && plan.DhcpRelaySourceIpv4sSubnet.ValueString() != "" {
			tenantReq.DhcpRelaySourceIpv4sSubnet = openapi.PtrString(plan.DhcpRelaySourceIpv4sSubnet.ValueString())
		} else {
			tenantReq.DhcpRelaySourceIpv4sSubnet = openapi.PtrString("")
		}
		hasChanges = true
	}

	if !plan.DhcpRelaySourceIpv6sSubnet.Equal(state.DhcpRelaySourceIpv6sSubnet) {
		if !plan.DhcpRelaySourceIpv6sSubnet.IsNull() && plan.DhcpRelaySourceIpv6sSubnet.ValueString() != "" {
			tenantReq.DhcpRelaySourceIpv6sSubnet = openapi.PtrString(plan.DhcpRelaySourceIpv6sSubnet.ValueString())
		} else {
			tenantReq.DhcpRelaySourceIpv6sSubnet = openapi.PtrString("")
		}
		hasChanges = true
	}

	if !plan.RouteDistinguisher.Equal(state.RouteDistinguisher) {
		if !plan.RouteDistinguisher.IsNull() && plan.RouteDistinguisher.ValueString() != "" {
			tenantReq.RouteDistinguisher = openapi.PtrString(plan.RouteDistinguisher.ValueString())
		} else {
			tenantReq.RouteDistinguisher = openapi.PtrString("")
		}
		hasChanges = true
	}

	if !plan.RouteTargetImport.Equal(state.RouteTargetImport) {
		if !plan.RouteTargetImport.IsNull() && plan.RouteTargetImport.ValueString() != "" {
			tenantReq.RouteTargetImport = openapi.PtrString(plan.RouteTargetImport.ValueString())
		} else {
			tenantReq.RouteTargetImport = openapi.PtrString("")
		}
		hasChanges = true
	}

	if !plan.RouteTargetExport.Equal(state.RouteTargetExport) {
		if !plan.RouteTargetExport.IsNull() && plan.RouteTargetExport.ValueString() != "" {
			tenantReq.RouteTargetExport = openapi.PtrString(plan.RouteTargetExport.ValueString())
		} else {
			tenantReq.RouteTargetExport = openapi.PtrString("")
		}
		hasChanges = true
	}

	importRouteMapChanged := !plan.ImportRouteMap.Equal(state.ImportRouteMap)
	importRouteMapRefTypeChanged := !plan.ImportRouteMapRefType.Equal(state.ImportRouteMapRefType)

	if importRouteMapChanged || importRouteMapRefTypeChanged {
		// Validate using one ref type supported rules
		if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
			plan.ImportRouteMap, plan.ImportRouteMapRefType,
			"import_route_map", "import_route_map_ref_type_",
			importRouteMapChanged, importRouteMapRefTypeChanged) {
			return
		}

		// Only send the base field if only it changed
		if importRouteMapChanged && !importRouteMapRefTypeChanged {
			// Just send the base field
			if !plan.ImportRouteMap.IsNull() && plan.ImportRouteMap.ValueString() != "" {
				tenantReq.ImportRouteMap = openapi.PtrString(plan.ImportRouteMap.ValueString())
			} else {
				tenantReq.ImportRouteMap = openapi.PtrString("")
			}
			hasChanges = true
		} else if importRouteMapRefTypeChanged {
			// Send both fields
			if !plan.ImportRouteMap.IsNull() && plan.ImportRouteMap.ValueString() != "" {
				tenantReq.ImportRouteMap = openapi.PtrString(plan.ImportRouteMap.ValueString())
			} else {
				tenantReq.ImportRouteMap = openapi.PtrString("")
			}

			if !plan.ImportRouteMapRefType.IsNull() && plan.ImportRouteMapRefType.ValueString() != "" {
				tenantReq.ImportRouteMapRefType = openapi.PtrString(plan.ImportRouteMapRefType.ValueString())
			} else {
				tenantReq.ImportRouteMapRefType = openapi.PtrString("")
			}
			hasChanges = true
		}
	}

	// Handle ExportRouteMap and ExportRouteMapRefType
	exportRouteMapChanged := !plan.ExportRouteMap.Equal(state.ExportRouteMap)
	exportRouteMapRefTypeChanged := !plan.ExportRouteMapRefType.Equal(state.ExportRouteMapRefType)

	if exportRouteMapChanged || exportRouteMapRefTypeChanged {
		// Validate using one ref type supported rules
		if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
			plan.ExportRouteMap, plan.ExportRouteMapRefType,
			"export_route_map", "export_route_map_ref_type_",
			exportRouteMapChanged, exportRouteMapRefTypeChanged) {
			return
		}

		// Only send the base field if only it changed
		if exportRouteMapChanged && !exportRouteMapRefTypeChanged {
			// Just send the base field
			if !plan.ExportRouteMap.IsNull() && plan.ExportRouteMap.ValueString() != "" {
				tenantReq.ExportRouteMap = openapi.PtrString(plan.ExportRouteMap.ValueString())
			} else {
				tenantReq.ExportRouteMap = openapi.PtrString("")
			}
			hasChanges = true
		} else if exportRouteMapRefTypeChanged {
			// Send both fields
			if !plan.ExportRouteMap.IsNull() && plan.ExportRouteMap.ValueString() != "" {
				tenantReq.ExportRouteMap = openapi.PtrString(plan.ExportRouteMap.ValueString())
			} else {
				tenantReq.ExportRouteMap = openapi.PtrString("")
			}

			if !plan.ExportRouteMapRefType.IsNull() && plan.ExportRouteMapRefType.ValueString() != "" {
				tenantReq.ExportRouteMapRefType = openapi.PtrString(plan.ExportRouteMapRefType.ValueString())
			} else {
				tenantReq.ExportRouteMapRefType = openapi.PtrString("")
			}
			hasChanges = true
		}
	}

	if !plan.DefaultOriginate.Equal(state.DefaultOriginate) {
		tenantReq.DefaultOriginate = openapi.PtrBool(plan.DefaultOriginate.ValueBool())
		hasChanges = true
	}

	oldRouteTenantsByIndex := make(map[int64]verityTenantRouteTenantModel)
	for _, rt := range state.RouteTenants {
		if !rt.Index.IsNull() {
			idx := rt.Index.ValueInt64()
			oldRouteTenantsByIndex[idx] = rt
		}
	}

	var changedRouteTenants []openapi.TenantsPutRequestTenantValueRouteTenantsInner
	routeTenantsChanged := false

	for _, rt := range plan.RouteTenants {
		if rt.Index.IsNull() {
			continue
		}

		idx := rt.Index.ValueInt64()
		oldRouteTenant, exists := oldRouteTenantsByIndex[idx]

		if !exists {
			// new route tenant, include all fields
			tenantRoute := openapi.TenantsPutRequestTenantValueRouteTenantsInner{
				Index: openapi.PtrInt32(int32(idx)),
			}

			if !rt.Enable.IsNull() {
				tenantRoute.Enable = openapi.PtrBool(rt.Enable.ValueBool())
			} else {
				tenantRoute.Enable = openapi.PtrBool(false)
			}

			if !rt.Tenant.IsNull() && rt.Tenant.ValueString() != "" {
				tenantRoute.Tenant = openapi.PtrString(rt.Tenant.ValueString())
			} else {
				tenantRoute.Tenant = openapi.PtrString("")
			}

			changedRouteTenants = append(changedRouteTenants, tenantRoute)
			routeTenantsChanged = true
			continue
		}

		// existing route tenant, check which fields changed
		tenantRoute := openapi.TenantsPutRequestTenantValueRouteTenantsInner{
			Index: openapi.PtrInt32(int32(idx)),
		}

		fieldChanged := false

		if !rt.Enable.Equal(oldRouteTenant.Enable) {
			tenantRoute.Enable = openapi.PtrBool(rt.Enable.ValueBool())
			fieldChanged = true
		}

		if !rt.Tenant.Equal(oldRouteTenant.Tenant) {
			if !rt.Tenant.IsNull() && rt.Tenant.ValueString() != "" {
				tenantRoute.Tenant = openapi.PtrString(rt.Tenant.ValueString())
			} else {
				tenantRoute.Tenant = openapi.PtrString("")
			}
			fieldChanged = true
		}

		if fieldChanged {
			changedRouteTenants = append(changedRouteTenants, tenantRoute)
			routeTenantsChanged = true
		}
	}

	for idx := range oldRouteTenantsByIndex {
		found := false
		for _, rt := range plan.RouteTenants {
			if !rt.Index.IsNull() && rt.Index.ValueInt64() == idx {
				found = true
				break
			}
		}

		if !found {
			// route tenant removed - include only the index for deletion
			deletedTenant := openapi.TenantsPutRequestTenantValueRouteTenantsInner{
				Index: openapi.PtrInt32(int32(idx)),
			}
			changedRouteTenants = append(changedRouteTenants, deletedTenant)
			routeTenantsChanged = true
		}
	}

	if routeTenantsChanged && len(changedRouteTenants) > 0 {
		tenantReq.RouteTenants = changedRouteTenants
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	bulkOpsMgr := r.provCtx.bulkOpsMgr
	operationID := bulkOpsMgr.AddPatch(ctx, "tenant", state.Name.ValueString(), tenantReq)
	r.provCtx.NotifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for tenant update operation %s to complete", operationID))
	if err := bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update Tenant %s", state.Name.ValueString()))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Tenant %s update operation completed successfully", state.Name.ValueString()))
	clearCache(ctx, r.provCtx, "tenants")

	var minState verityTenantResourceModel
	minState.Name = types.StringValue(state.Name.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if tenantData, exists := bulkMgr.GetResourceResponse("tenant", state.Name.ValueString()); exists {
			// Use the cached data from the API response with plan values as fallback
			updatedState := populateTenantState(ctx, minState, tenantData, &plan)
			resp.Diagnostics.Append(resp.State.Set(ctx, &updatedState)...)
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

func (r *verityTenantResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityTenantResourceModel
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

	tenantName := state.Name.ValueString()
	bulkOpsMgr := r.provCtx.bulkOpsMgr
	operationID := bulkOpsMgr.AddDelete(ctx, "tenant", tenantName)
	r.provCtx.NotifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for tenant deletion operation %s to complete", operationID))
	if err := bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Delete Tenant %s", tenantName))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Tenant %s deletion operation completed successfully", tenantName))
	clearCache(ctx, r.provCtx, "tenants")
	resp.State.RemoveResource(ctx)
}

func (r *verityTenantResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func populateTenantState(ctx context.Context, state verityTenantResourceModel, tenantData map[string]interface{}, plan *verityTenantResourceModel) verityTenantResourceModel {
	state.Name = types.StringValue(fmt.Sprintf("%v", tenantData["name"]))

	// For each field, check if it's in the API response first,
	// if not: use plan value (if plan provided and not null), otherwise preserve current state value

	if val, ok := tenantData["enable"].(bool); ok {
		state.Enable = types.BoolValue(val)
	} else if plan != nil && !plan.Enable.IsNull() {
		state.Enable = plan.Enable
	}

	if op, ok := tenantData["object_properties"].(map[string]interface{}); ok {
		state.ObjectProperties = []verityTenantObjectPropertiesModel{
			{
				Group: types.StringValue(fmt.Sprintf("%v", op["group"])),
			},
		}
	} else if plan != nil && len(plan.ObjectProperties) > 0 {
		state.ObjectProperties = plan.ObjectProperties
	}

	if val, ok := tenantData["layer_3_vni"]; ok {
		if val == nil {
			state.Layer3Vni = types.Int64Null()
		} else {
			switch v := val.(type) {
			case float64:
				state.Layer3Vni = types.Int64Value(int64(v))
			case int:
				state.Layer3Vni = types.Int64Value(int64(v))
			case int32:
				state.Layer3Vni = types.Int64Value(int64(v))
			default:
				if plan != nil && !plan.Layer3Vni.IsNull() {
					state.Layer3Vni = plan.Layer3Vni
				}
			}
		}
	} else if plan != nil && !plan.Layer3Vni.IsNull() {
		state.Layer3Vni = plan.Layer3Vni
	}

	if val, ok := tenantData["layer_3_vni_auto_assigned_"].(bool); ok {
		state.Layer3VniAutoAssigned = types.BoolValue(val)
	} else if plan != nil && !plan.Layer3VniAutoAssigned.IsNull() {
		state.Layer3VniAutoAssigned = plan.Layer3VniAutoAssigned
	}

	if val, ok := tenantData["layer_3_vlan"]; ok {
		if val == nil {
			state.Layer3Vlan = types.Int64Null()
		} else {
			switch v := val.(type) {
			case float64:
				state.Layer3Vlan = types.Int64Value(int64(v))
			case int:
				state.Layer3Vlan = types.Int64Value(int64(v))
			case int32:
				state.Layer3Vlan = types.Int64Value(int64(v))
			default:
				if plan != nil && !plan.Layer3Vlan.IsNull() {
					state.Layer3Vlan = plan.Layer3Vlan
				}
			}
		}
	}

	if val, ok := tenantData["layer_3_vlan_auto_assigned_"].(bool); ok {
		state.Layer3VlanAutoAssigned = types.BoolValue(val)
	} else if plan != nil && !plan.Layer3VlanAutoAssigned.IsNull() {
		state.Layer3VlanAutoAssigned = plan.Layer3VlanAutoAssigned
	}

	if val, ok := tenantData["vrf_name"].(string); ok {
		state.VrfName = types.StringValue(val)
	} else if plan != nil && !plan.VrfNameAutoAssigned.IsNull() && plan.VrfNameAutoAssigned.ValueBool() {
		state.VrfName = types.StringValue("")
	} else if plan != nil && !plan.VrfName.IsNull() {
		state.VrfName = plan.VrfName
	}

	if val, ok := tenantData["vrf_name_auto_assigned_"].(bool); ok {
		state.VrfNameAutoAssigned = types.BoolValue(val)
	} else if plan != nil && !plan.VrfNameAutoAssigned.IsNull() {
		state.VrfNameAutoAssigned = plan.VrfNameAutoAssigned
	}

	if val, ok := tenantData["default_originate"].(bool); ok {
		state.DefaultOriginate = types.BoolValue(val)
	} else if plan != nil && !plan.DefaultOriginate.IsNull() {
		state.DefaultOriginate = plan.DefaultOriginate
	}

	stringFields := map[string]*types.String{
		"dhcp_relay_source_ipv4s_subnet": &state.DhcpRelaySourceIpv4sSubnet,
		"dhcp_relay_source_ipv6s_subnet": &state.DhcpRelaySourceIpv6sSubnet,
		"route_distinguisher":            &state.RouteDistinguisher,
		"route_target_import":            &state.RouteTargetImport,
		"route_target_export":            &state.RouteTargetExport,
		"import_route_map":               &state.ImportRouteMap,
		"import_route_map_ref_type_":     &state.ImportRouteMapRefType,
		"export_route_map":               &state.ExportRouteMap,
		"export_route_map_ref_type_":     &state.ExportRouteMapRefType,
	}

	for apiKey, stateField := range stringFields {
		if val, ok := tenantData[apiKey].(string); ok {
			*stateField = types.StringValue(val)
		} else if plan != nil {
			switch apiKey {
			case "dhcp_relay_source_ipv4s_subnet":
				if !plan.DhcpRelaySourceIpv4sSubnet.IsNull() {
					*stateField = plan.DhcpRelaySourceIpv4sSubnet
				}
			case "dhcp_relay_source_ipv6s_subnet":
				if !plan.DhcpRelaySourceIpv6sSubnet.IsNull() {
					*stateField = plan.DhcpRelaySourceIpv6sSubnet
				}
			case "route_distinguisher":
				if !plan.RouteDistinguisher.IsNull() {
					*stateField = plan.RouteDistinguisher
				}
			case "route_target_import":
				if !plan.RouteTargetImport.IsNull() {
					*stateField = plan.RouteTargetImport
				}
			case "route_target_export":
				if !plan.RouteTargetExport.IsNull() {
					*stateField = plan.RouteTargetExport
				}
			case "import_route_map":
				if !plan.ImportRouteMap.IsNull() {
					*stateField = plan.ImportRouteMap
				}
			case "import_route_map_ref_type_":
				if !plan.ImportRouteMapRefType.IsNull() {
					*stateField = plan.ImportRouteMapRefType
				}
			case "export_route_map":
				if !plan.ExportRouteMap.IsNull() {
					*stateField = plan.ExportRouteMap
				}
			case "export_route_map_ref_type_":
				if !plan.ExportRouteMapRefType.IsNull() {
					*stateField = plan.ExportRouteMapRefType
				}
			}
		}
	}

	if rtVal, ok := tenantData["route_tenants"].([]interface{}); ok {
		var routeTenants []verityTenantRouteTenantModel
		for _, rt := range rtVal {
			rtMap, ok := rt.(map[string]interface{})
			if !ok {
				continue
			}
			routeTenant := verityTenantRouteTenantModel{}
			if val, ok := rtMap["enable"].(bool); ok {
				routeTenant.Enable = types.BoolValue(val)
			} else {
				routeTenant.Enable = types.BoolNull()
			}
			if val, ok := rtMap["tenant"].(string); ok {
				routeTenant.Tenant = types.StringValue(val)
			} else {
				routeTenant.Tenant = types.StringNull()
			}
			if val, ok := rtMap["index"].(float64); ok {
				routeTenant.Index = types.Int64Value(int64(val))
			} else {
				routeTenant.Index = types.Int64Null()
			}
			routeTenants = append(routeTenants, routeTenant)
		}
		state.RouteTenants = routeTenants
	} else if plan != nil && len(plan.RouteTenants) > 0 {
		state.RouteTenants = plan.RouteTenants
	}

	return state
}

func (r *verityTenantResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Skip modification if we're deleting the resource
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityTenantResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Validate auto-assigned field specifications in configuration
	// Check the actual configuration, not the plan
	var config verityTenantResourceModel
	if !req.Config.Raw.IsNull() {
		resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
		if resp.Diagnostics.HasError() {
			return
		}

		if !config.Layer3VniAutoAssigned.IsNull() && config.Layer3VniAutoAssigned.ValueBool() {
			if !config.Layer3Vni.IsNull() && !config.Layer3Vni.IsUnknown() {
				resp.Diagnostics.AddError(
					"Layer 3 VNI cannot be specified when auto-assigned",
					"The 'layer_3_vni' field cannot be specified in the configuration when 'layer_3_vni_auto_assigned_' is set to true. The API will assign this value automatically.",
				)
				return
			}
		}

		if !config.Layer3VlanAutoAssigned.IsNull() && config.Layer3VlanAutoAssigned.ValueBool() {
			if !config.Layer3Vlan.IsNull() && !config.Layer3Vlan.IsUnknown() {
				resp.Diagnostics.AddError(
					"Layer 3 VLAN cannot be specified when auto-assigned",
					"The 'layer_3_vlan' field cannot be specified in the configuration when 'layer_3_vlan_auto_assigned_' is set to true. The API will assign this value automatically.",
				)
				return
			}
		}

		if !config.VrfNameAutoAssigned.IsNull() && config.VrfNameAutoAssigned.ValueBool() {
			if !config.VrfName.IsNull() && !config.VrfName.IsUnknown() && config.VrfName.ValueString() != "" {
				resp.Diagnostics.AddError(
					"VRF name cannot be specified when auto-assigned",
					"The 'vrf_name' field cannot be specified in the configuration when 'vrf_name_auto_assigned_' is set to true. The API will assign this value automatically.",
				)
				return
			}
		}
	}

	// For new resources (where state is null)
	if req.State.Raw.IsNull() {
		if !plan.Layer3VniAutoAssigned.IsNull() && plan.Layer3VniAutoAssigned.ValueBool() {
			resp.Plan.SetAttribute(ctx, path.Root("layer_3_vni"), types.Int64Unknown())
		}

		if !plan.Layer3VlanAutoAssigned.IsNull() && plan.Layer3VlanAutoAssigned.ValueBool() {
			resp.Plan.SetAttribute(ctx, path.Root("layer_3_vlan"), types.Int64Unknown())
		}

		if !plan.VrfNameAutoAssigned.IsNull() && plan.VrfNameAutoAssigned.ValueBool() {
			resp.Plan.SetAttribute(ctx, path.Root("vrf_name"), types.StringUnknown())
		}
		return
	}

	var state verityTenantResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Handle auto-assigned field behavior for existing resources
	if !plan.Layer3VniAutoAssigned.IsNull() && plan.Layer3VniAutoAssigned.ValueBool() {
		if !plan.Layer3VniAutoAssigned.Equal(state.Layer3VniAutoAssigned) {
			// layer_3_vni_auto_assigned_ is changing to true, API will assign value
			resp.Plan.SetAttribute(ctx, path.Root("layer_3_vni"), types.Int64Unknown())
			resp.Diagnostics.AddWarning(
				"Layer 3 VNI will be assigned by the API",
				"The 'layer_3_vni' field will be automatically assigned by the API because 'layer_3_vni_auto_assigned_' is being set to true.",
			)
		} else if !plan.Layer3Vni.Equal(state.Layer3Vni) {
			// User tried to change Layer3Vni but it's auto-assigned
			resp.Diagnostics.AddWarning(
				"Ignoring layer_3_vni changes with auto-assignment enabled",
				"The 'layer_3_vni' field changes will be ignored because 'layer_3_vni_auto_assigned_' is set to true. The API will assign this value automatically.",
			)
			// Keep the current state value to suppress the diff
			if !state.Layer3Vni.IsNull() {
				resp.Plan.SetAttribute(ctx, path.Root("layer_3_vni"), state.Layer3Vni)
			}
		}
	}

	if !plan.Layer3VlanAutoAssigned.IsNull() && plan.Layer3VlanAutoAssigned.ValueBool() {
		if !plan.Layer3VlanAutoAssigned.Equal(state.Layer3VlanAutoAssigned) {
			// layer_3_vlan_auto_assigned_ is changing to true, API will assign value
			resp.Plan.SetAttribute(ctx, path.Root("layer_3_vlan"), types.Int64Unknown())
			resp.Diagnostics.AddWarning(
				"Layer 3 VLAN will be assigned by the API",
				"The 'layer_3_vlan' field will be automatically assigned by the API because 'layer_3_vlan_auto_assigned_' is being set to true.",
			)
		} else if !plan.Layer3Vlan.Equal(state.Layer3Vlan) {
			// User tried to change Layer3Vlan but it's auto-assigned
			resp.Diagnostics.AddWarning(
				"Ignoring layer_3_vlan changes with auto-assignment enabled",
				"The 'layer_3_vlan' field changes will be ignored because 'layer_3_vlan_auto_assigned_' is set to true. The API will assign this value automatically.",
			)
			// Keep the current state value to suppress the diff
			if !state.Layer3Vlan.IsNull() {
				resp.Plan.SetAttribute(ctx, path.Root("layer_3_vlan"), state.Layer3Vlan)
			}
		}
	}

	if !plan.VrfNameAutoAssigned.IsNull() && plan.VrfNameAutoAssigned.ValueBool() {
		if !plan.VrfNameAutoAssigned.Equal(state.VrfNameAutoAssigned) {
			// vrf_name_auto_assigned_ is changing to true, API will assign value
			resp.Plan.SetAttribute(ctx, path.Root("vrf_name"), types.StringUnknown())
			resp.Diagnostics.AddWarning(
				"VRF name will be assigned by the API",
				"The 'vrf_name' field will be automatically assigned by the API because 'vrf_name_auto_assigned_' is being set to true.",
			)
		} else if !plan.VrfName.Equal(state.VrfName) {
			// User tried to change VrfName but it's auto-assigned
			resp.Diagnostics.AddWarning(
				"Ignoring vrf_name changes with auto-assignment enabled",
				"The 'vrf_name' field changes will be ignored because 'vrf_name_auto_assigned_' is set to true. The API will assign this value automatically.",
			)
			// Keep the current state value to suppress the diff
			if !state.VrfName.IsNull() {
				resp.Plan.SetAttribute(ctx, path.Root("vrf_name"), state.VrfName)
			}
		}
	}

	// Check for ineffective changes to import_route_map_ref_type_
	// Warn if import_route_map_ref_type_ is changing BUT import_route_map is empty and NOT changing.
	if !plan.ImportRouteMapRefType.Equal(state.ImportRouteMapRefType) &&
		(plan.ImportRouteMap.IsNull() || plan.ImportRouteMap.ValueString() == "") &&
		plan.ImportRouteMap.Equal(state.ImportRouteMap) {
		resp.Diagnostics.AddWarning(
			"Ineffective change to import_route_map_ref_type_",
			"The change to 'import_route_map_ref_type_' will likely be ignored by the API because 'import_route_map' is empty and not being changed. The API may require 'import_route_map' to have a value for 'import_route_map_ref_type_' to be effective.",
		)
	}

	// Check for ineffective changes to export_route_map_ref_type_
	// Warn if export_route_map_ref_type_ is changing BUT export_route_map is empty and NOT changing.
	if !plan.ExportRouteMapRefType.Equal(state.ExportRouteMapRefType) &&
		(plan.ExportRouteMap.IsNull() || plan.ExportRouteMap.ValueString() == "") &&
		plan.ExportRouteMap.Equal(state.ExportRouteMap) {
		resp.Diagnostics.AddWarning(
			"Ineffective change to export_route_map_ref_type_",
			"The change to 'export_route_map_ref_type_' will likely be ignored by the API because 'export_route_map' is empty and not being changed. The API may require 'export_route_map' to have a value for 'export_route_map_ref_type_' to be effective.",
		)
	}
}
