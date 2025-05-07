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
	_ resource.Resource                = &verityTenantResource{}
	_ resource.ResourceWithConfigure   = &verityTenantResource{}
	_ resource.ResourceWithImportState = &verityTenantResource{}
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
	Name                     types.String                        `tfsdk:"name"`
	Enable                   types.Bool                          `tfsdk:"enable"`
	ObjectProperties         []verityTenantObjectPropertiesModel `tfsdk:"object_properties"`
	Layer3Vni                types.Int64                         `tfsdk:"layer_3_vni"`
	Layer3VniAutoAssigned    types.Bool                          `tfsdk:"layer_3_vni_auto_assigned_"`
	Layer3Vlan               types.Int64                         `tfsdk:"layer_3_vlan"`
	Layer3VlanAutoAssigned   types.Bool                          `tfsdk:"layer_3_vlan_auto_assigned_"`
	DhcpRelaySourceIpsSubnet types.String                        `tfsdk:"dhcp_relay_source_ips_subnet"`
	RouteDistinguisher       types.String                        `tfsdk:"route_distinguisher"`
	RouteTargetImport        types.String                        `tfsdk:"route_target_import"`
	RouteTargetExport        types.String                        `tfsdk:"route_target_export"`
	ImportRouteMap           types.String                        `tfsdk:"import_route_map"`
	ImportRouteMapRefType    types.String                        `tfsdk:"import_route_map_ref_type_"`
	ExportRouteMap           types.String                        `tfsdk:"export_route_map"`
	ExportRouteMapRefType    types.String                        `tfsdk:"export_route_map_ref_type_"`
	VrfName                  types.String                        `tfsdk:"vrf_name"`
	VrfNameAutoAssigned      types.Bool                          `tfsdk:"vrf_name_auto_assigned_"`
	RouteTenants             []verityTenantRouteTenantModel      `tfsdk:"route_tenants"`
	DefaultOriginate         types.Bool                          `tfsdk:"default_originate"`
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
				Computed:    true,
				Default:     booldefault.StaticBool(false),
			},
			"layer_3_vni": schema.Int64Attribute{
				Description: "VNI value used to transport traffic between services of a Tenant",
				Optional:    true,
				Computed:    true,
			},
			"layer_3_vni_auto_assigned_": schema.BoolAttribute{
				Description: "Whether or not the value in layer_3_vni field has been automatically assigned",
				Optional:    true,
			},
			"layer_3_vlan": schema.Int64Attribute{
				Description: "VLAN value used to transport traffic between services of a Tenant",
				Optional:    true,
				Computed:    true,
			},
			"layer_3_vlan_auto_assigned_": schema.BoolAttribute{
				Description: "Whether or not the value in layer_3_vlan field has been automatically assigned",
				Optional:    true,
			},
			"dhcp_relay_source_ips_subnet": schema.StringAttribute{
				Description: "Range of IP addresses used to configure the source IP of each DHCP Relay",
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
				Description: "Virtual Routing and Forwarding instance name",
				Optional:    true,
				Computed:    true,
			},
			"vrf_name_auto_assigned_": schema.BoolAttribute{
				Description: "Whether or not the value in vrf_name field has been automatically assigned",
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

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Authenticate",
			fmt.Sprintf("Error authenticating with API: %s", err),
		)
		return
	}

	name := plan.Name.ValueString()

	tenantReq := openapi.ConfigPutRequestTenantTenantName{
		Name:                     openapi.PtrString(name),
		Enable:                   openapi.PtrBool(plan.Enable.ValueBool()),
		DhcpRelaySourceIpsSubnet: openapi.PtrString(plan.DhcpRelaySourceIpsSubnet.ValueString()),
		RouteDistinguisher:       openapi.PtrString(plan.RouteDistinguisher.ValueString()),
		RouteTargetImport:        openapi.PtrString(plan.RouteTargetImport.ValueString()),
		RouteTargetExport:        openapi.PtrString(plan.RouteTargetExport.ValueString()),
		ImportRouteMap:           openapi.PtrString(plan.ImportRouteMap.ValueString()),
		ExportRouteMap:           openapi.PtrString(plan.ExportRouteMap.ValueString()),
		VrfName:                  openapi.PtrString(plan.VrfName.ValueString()),
		RouteTenants:             []openapi.ConfigPutRequestTenantTenantNameRouteTenantsInner{},
	}

	if len(plan.ObjectProperties) > 0 {
		objProps := openapi.ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties{}
		if !plan.ObjectProperties[0].Group.IsNull() {
			objProps.Group = openapi.PtrString(plan.ObjectProperties[0].Group.ValueString())
		} else {
			objProps.Group = nil
		}
		tenantReq.ObjectProperties = &objProps
	} else {
		tenantReq.ObjectProperties = nil
	}

	if !plan.Layer3Vni.IsNull() {
		val := int32(plan.Layer3Vni.ValueInt64())
		tenantReq.Layer3Vni = *openapi.NewNullableInt32(&val)
	} else {
		tenantReq.Layer3Vni = *openapi.NewNullableInt32(nil)
	}
	if !plan.Layer3Vlan.IsNull() {
		val := int32(plan.Layer3Vlan.ValueInt64())
		tenantReq.Layer3Vlan = *openapi.NewNullableInt32(&val)
	} else {
		tenantReq.Layer3Vlan = *openapi.NewNullableInt32(nil)
	}
	if !plan.Layer3VniAutoAssigned.IsNull() {
		tenantReq.Layer3VniAutoAssigned = openapi.PtrBool(plan.Layer3VniAutoAssigned.ValueBool())
	}
	if !plan.Layer3VlanAutoAssigned.IsNull() {
		tenantReq.Layer3VlanAutoAssigned = openapi.PtrBool(plan.Layer3VlanAutoAssigned.ValueBool())
	}
	if !plan.VrfNameAutoAssigned.IsNull() {
		tenantReq.VrfNameAutoAssigned = openapi.PtrBool(plan.VrfNameAutoAssigned.ValueBool())
	}
	if !plan.DefaultOriginate.IsNull() {
		tenantReq.DefaultOriginate = openapi.PtrBool(plan.DefaultOriginate.ValueBool())
	}

	if len(plan.RouteTenants) > 0 {
		for _, rt := range plan.RouteTenants {
			tenant := openapi.ConfigPutRequestTenantTenantNameRouteTenantsInner{
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
	operationID := bulkOpsMgr.AddTenantPut(ctx, name, tenantReq)

	provCtx.NotifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for tenant creation operation %s to complete", operationID))
	if err := bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Create Tenant",
			fmt.Sprintf("Error creating tenant %s: %v", name, err),
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
		if tenantData, exists := bulkMgr.GetTenantResponse(name); exists {
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
		tenantData, exists = bulkOpsMgr.GetTenantResponse(tenantName)
		if exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached tenant data for %s from recent operation", tenantName))
			state = populateTenantState(ctx, state, tenantData, nil)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if bulkOpsMgr != nil && bulkOpsMgr.HasPendingOrRecentTenantOperations() {
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
		resp.Diagnostics.AddError("Failed to Read Tenant", err.Error())
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

	if !plan.Layer3Vni.Equal(state.Layer3Vni) && !plan.Layer3VniAutoAssigned.IsNull() && plan.Layer3VniAutoAssigned.ValueBool() {
		resp.Diagnostics.AddError(
			"Cannot modify auto-assigned field",
			"The 'layer_3_vni' field cannot be modified because 'layer_3_vni_auto_assigned_' is set to true.",
		)
		return
	}

	if !plan.Layer3Vlan.Equal(state.Layer3Vlan) && !plan.Layer3VlanAutoAssigned.IsNull() && plan.Layer3VlanAutoAssigned.ValueBool() {
		resp.Diagnostics.AddError(
			"Cannot modify auto-assigned field",
			"The 'layer_3_vlan' field cannot be modified because 'layer_3_vlan_auto_assigned_' is set to true.",
		)
		return
	}

	if !plan.VrfName.Equal(state.VrfName) && !plan.VrfNameAutoAssigned.IsNull() && plan.VrfNameAutoAssigned.ValueBool() {
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

	tenantReq := openapi.ConfigPutRequestTenantTenantName{}
	hasChanges := false

	if len(plan.ObjectProperties) > 0 {
		if len(state.ObjectProperties) == 0 || !plan.ObjectProperties[0].Group.Equal(state.ObjectProperties[0].Group) {
			objProps := openapi.ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties{}
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

	if !plan.Layer3Vni.Equal(state.Layer3Vni) {
		if !plan.Layer3Vni.IsNull() {
			val := int32(plan.Layer3Vni.ValueInt64())
			tenantReq.Layer3Vni = *openapi.NewNullableInt32(&val)
		} else {
			tenantReq.Layer3Vni = *openapi.NewNullableInt32(nil)
		}
		if plan.Layer3VniAutoAssigned.IsNull() {
			if !state.Layer3VniAutoAssigned.IsNull() {
				tenantReq.Layer3VniAutoAssigned = openapi.PtrBool(state.Layer3VniAutoAssigned.ValueBool())
			} else {
				tenantReq.Layer3VniAutoAssigned = openapi.PtrBool(false)
			}
		}
		hasChanges = true
	} else if !plan.Layer3VniAutoAssigned.IsNull() {
		if !plan.Layer3VniAutoAssigned.Equal(state.Layer3VniAutoAssigned) {
			tenantReq.Layer3VniAutoAssigned = openapi.PtrBool(plan.Layer3VniAutoAssigned.ValueBool())
			hasChanges = true
		}
	}

	if !plan.Layer3Vlan.Equal(state.Layer3Vlan) {
		if !plan.Layer3Vlan.IsNull() {
			val := int32(plan.Layer3Vlan.ValueInt64())
			tenantReq.Layer3Vlan = *openapi.NewNullableInt32(&val)
		} else {
			tenantReq.Layer3Vlan = *openapi.NewNullableInt32(nil)
		}
		if plan.Layer3VlanAutoAssigned.IsNull() {
			if !state.Layer3VlanAutoAssigned.IsNull() {
				tenantReq.Layer3VlanAutoAssigned = openapi.PtrBool(state.Layer3VlanAutoAssigned.ValueBool())
			} else {
				tenantReq.Layer3VlanAutoAssigned = openapi.PtrBool(false)
			}
		}
		hasChanges = true
	} else if !plan.Layer3VlanAutoAssigned.IsNull() {
		if !plan.Layer3VlanAutoAssigned.Equal(state.Layer3VlanAutoAssigned) {
			tenantReq.Layer3VlanAutoAssigned = openapi.PtrBool(plan.Layer3VlanAutoAssigned.ValueBool())
			hasChanges = true
		}
	}

	if !plan.VrfName.Equal(state.VrfName) {
		if !plan.VrfName.IsNull() && plan.VrfName.ValueString() != "" {
			tenantReq.VrfName = openapi.PtrString(plan.VrfName.ValueString())
		} else {
			tenantReq.VrfName = openapi.PtrString("")
		}
		if plan.VrfNameAutoAssigned.IsNull() {
			if !state.VrfNameAutoAssigned.IsNull() {
				tenantReq.VrfNameAutoAssigned = openapi.PtrBool(state.VrfNameAutoAssigned.ValueBool())
			} else {
				tenantReq.VrfNameAutoAssigned = openapi.PtrBool(false)
			}
		}
		hasChanges = true
	} else if !plan.VrfNameAutoAssigned.IsNull() {
		if !plan.VrfNameAutoAssigned.Equal(state.VrfNameAutoAssigned) {
			tenantReq.VrfNameAutoAssigned = openapi.PtrBool(plan.VrfNameAutoAssigned.ValueBool())
			hasChanges = true
		}
	}

	if !plan.DhcpRelaySourceIpsSubnet.Equal(state.DhcpRelaySourceIpsSubnet) {
		if !plan.DhcpRelaySourceIpsSubnet.IsNull() && plan.DhcpRelaySourceIpsSubnet.ValueString() != "" {
			tenantReq.DhcpRelaySourceIpsSubnet = openapi.PtrString(plan.DhcpRelaySourceIpsSubnet.ValueString())
		} else {
			tenantReq.DhcpRelaySourceIpsSubnet = openapi.PtrString("")
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

	if !plan.ImportRouteMap.Equal(state.ImportRouteMap) {
		if !plan.ImportRouteMap.IsNull() && plan.ImportRouteMap.ValueString() != "" {
			tenantReq.ImportRouteMap = openapi.PtrString(plan.ImportRouteMap.ValueString())
		} else {
			tenantReq.ImportRouteMap = openapi.PtrString("")
		}
		hasChanges = true
	}

	if !plan.ExportRouteMap.Equal(state.ExportRouteMap) {
		if !plan.ExportRouteMap.IsNull() && plan.ExportRouteMap.ValueString() != "" {
			tenantReq.ExportRouteMap = openapi.PtrString(plan.ExportRouteMap.ValueString())
		} else {
			tenantReq.ExportRouteMap = openapi.PtrString("")
		}
		hasChanges = true
	}

	if !plan.ImportRouteMapRefType.Equal(state.ImportRouteMapRefType) {
		if !plan.ImportRouteMapRefType.IsNull() && plan.ImportRouteMapRefType.ValueString() != "" {
			tenantReq.ImportRouteMapRefType = openapi.PtrString(plan.ImportRouteMapRefType.ValueString())
		} else {
			tenantReq.ImportRouteMapRefType = openapi.PtrString("")
		}
		hasChanges = true
	}

	if !plan.ExportRouteMapRefType.Equal(state.ExportRouteMapRefType) {
		if !plan.ExportRouteMapRefType.IsNull() && plan.ExportRouteMapRefType.ValueString() != "" {
			tenantReq.ExportRouteMapRefType = openapi.PtrString(plan.ExportRouteMapRefType.ValueString())
		} else {
			tenantReq.ExportRouteMapRefType = openapi.PtrString("")
		}
		hasChanges = true
	}

	if !plan.DefaultOriginate.Equal(state.DefaultOriginate) {
		tenantReq.DefaultOriginate = openapi.PtrBool(plan.DefaultOriginate.ValueBool())
		hasChanges = true
	}

	if len(plan.RouteTenants) != len(state.RouteTenants) {
		oldRouteTenantsByIndex := make(map[int64]verityTenantRouteTenantModel)
		for _, rt := range state.RouteTenants {
			if !rt.Index.IsNull() {
				idx := rt.Index.ValueInt64()
				oldRouteTenantsByIndex[idx] = rt
			}
		}

		var changedRouteTenants []openapi.ConfigPutRequestTenantTenantNameRouteTenantsInner
		for _, rt := range plan.RouteTenants {
			if rt.Index.IsNull() {
				continue
			}

			idx := rt.Index.ValueInt64()
			oldRouteTenant, exists := oldRouteTenantsByIndex[idx]

			routeTenantChanged := !exists
			if exists {
				if !rt.Enable.Equal(oldRouteTenant.Enable) ||
					!rt.Tenant.Equal(oldRouteTenant.Tenant) {
					routeTenantChanged = true
				}
			}

			if routeTenantChanged {
				tenantRoute := openapi.ConfigPutRequestTenantTenantNameRouteTenantsInner{
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
			}
		}

		if len(changedRouteTenants) > 0 {
			tenantReq.RouteTenants = changedRouteTenants
			hasChanges = true
		}
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	bulkOpsMgr := r.provCtx.bulkOpsMgr
	operationID := bulkOpsMgr.AddTenantPatch(ctx, state.Name.ValueString(), tenantReq)
	r.provCtx.NotifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for tenant update operation %s to complete", operationID))
	if err := bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Update Tenant",
			fmt.Sprintf("Error updating tenant %s: %v", state.Name.ValueString(), err),
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Tenant %s update operation completed successfully", state.Name.ValueString()))
	clearCache(ctx, r.provCtx, "tenants")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
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
	operationID := bulkOpsMgr.AddTenantDelete(ctx, tenantName)
	r.provCtx.NotifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for tenant deletion operation %s to complete", operationID))
	if err := bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Delete Tenant",
			fmt.Sprintf("Error deleting tenant %s: %v", tenantName, err),
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
	// if not and we have a plan value, use that instead of null

	if val, ok := tenantData["enable"].(bool); ok {
		state.Enable = types.BoolValue(val)
	} else if plan != nil && !plan.Enable.IsNull() {
		state.Enable = plan.Enable
	} else {
		state.Enable = types.BoolNull()
	}

	if op, ok := tenantData["object_properties"].(map[string]interface{}); ok {
		state.ObjectProperties = []verityTenantObjectPropertiesModel{
			{
				Group: types.StringValue(fmt.Sprintf("%v", op["group"])),
			},
		}
	} else if plan != nil && len(plan.ObjectProperties) > 0 {
		state.ObjectProperties = plan.ObjectProperties
	} else {
		state.ObjectProperties = nil
	}

	if val, ok := tenantData["layer_3_vni"]; ok {
		switch v := val.(type) {
		case float64:
			state.Layer3Vni = types.Int64Value(int64(v))
		case int:
			state.Layer3Vni = types.Int64Value(int64(v))
		default:
			if plan != nil && !plan.Layer3VniAutoAssigned.IsNull() && plan.Layer3VniAutoAssigned.ValueBool() {
				state.Layer3Vni = types.Int64Null()
			} else if plan != nil && !plan.Layer3Vni.IsNull() {
				state.Layer3Vni = plan.Layer3Vni
			} else {
				state.Layer3Vni = types.Int64Null()
			}
		}
	} else if plan != nil && !plan.Layer3Vni.IsNull() {
		state.Layer3Vni = plan.Layer3Vni
	} else {
		state.Layer3Vni = types.Int64Null()
	}

	if val, ok := tenantData["layer_3_vni_auto_assigned_"].(bool); ok {
		state.Layer3VniAutoAssigned = types.BoolValue(val)
	} else if plan != nil && !plan.Layer3VniAutoAssigned.IsNull() {
		state.Layer3VniAutoAssigned = plan.Layer3VniAutoAssigned
	} else {
		state.Layer3VniAutoAssigned = types.BoolNull()
	}

	if val, ok := tenantData["layer_3_vlan"]; ok {
		switch v := val.(type) {
		case float64:
			state.Layer3Vlan = types.Int64Value(int64(v))
		case int:
			state.Layer3Vlan = types.Int64Value(int64(v))
		default:
			if plan != nil && !plan.Layer3VlanAutoAssigned.IsNull() && plan.Layer3VlanAutoAssigned.ValueBool() {
				state.Layer3Vlan = types.Int64Null()
			} else if plan != nil && !plan.Layer3Vlan.IsNull() {
				state.Layer3Vlan = plan.Layer3Vlan
			} else {
				state.Layer3Vlan = types.Int64Null()
			}
		}
	} else if plan != nil && !plan.Layer3Vlan.IsNull() {
		state.Layer3Vlan = plan.Layer3Vlan
	} else {
		state.Layer3Vlan = types.Int64Null()
	}

	if val, ok := tenantData["layer_3_vlan_auto_assigned_"].(bool); ok {
		state.Layer3VlanAutoAssigned = types.BoolValue(val)
	} else if plan != nil && !plan.Layer3VlanAutoAssigned.IsNull() {
		state.Layer3VlanAutoAssigned = plan.Layer3VlanAutoAssigned
	} else {
		state.Layer3VlanAutoAssigned = types.BoolNull()
	}

	if val, ok := tenantData["vrf_name"].(string); ok {
		state.VrfName = types.StringValue(val)
	} else if plan != nil && !plan.VrfNameAutoAssigned.IsNull() && plan.VrfNameAutoAssigned.ValueBool() {
		state.VrfName = types.StringValue("")
	} else if plan != nil && !plan.VrfName.IsNull() {
		state.VrfName = plan.VrfName
	} else {
		state.VrfName = types.StringNull()
	}

	if val, ok := tenantData["vrf_name_auto_assigned_"].(bool); ok {
		state.VrfNameAutoAssigned = types.BoolValue(val)
	} else if plan != nil && !plan.VrfNameAutoAssigned.IsNull() {
		state.VrfNameAutoAssigned = plan.VrfNameAutoAssigned
	} else {
		state.VrfNameAutoAssigned = types.BoolNull()
	}

	if val, ok := tenantData["default_originate"].(bool); ok {
		state.DefaultOriginate = types.BoolValue(val)
	} else if plan != nil && !plan.DefaultOriginate.IsNull() {
		state.DefaultOriginate = plan.DefaultOriginate
	} else {
		state.DefaultOriginate = types.BoolNull()
	}

	stringFields := map[string]*types.String{
		"dhcp_relay_source_ips_subnet": &state.DhcpRelaySourceIpsSubnet,
		"route_distinguisher":          &state.RouteDistinguisher,
		"route_target_import":          &state.RouteTargetImport,
		"route_target_export":          &state.RouteTargetExport,
		"import_route_map":             &state.ImportRouteMap,
		"import_route_map_ref_type_":   &state.ImportRouteMapRefType,
		"export_route_map":             &state.ExportRouteMap,
		"export_route_map_ref_type_":   &state.ExportRouteMapRefType,
	}

	for apiKey, stateField := range stringFields {
		if val, ok := tenantData[apiKey].(string); ok {
			*stateField = types.StringValue(val)
		} else if plan != nil {
			switch apiKey {
			case "dhcp_relay_source_ips_subnet":
				if !plan.DhcpRelaySourceIpsSubnet.IsNull() {
					*stateField = plan.DhcpRelaySourceIpsSubnet
				} else {
					*stateField = types.StringNull()
				}
			case "route_distinguisher":
				if !plan.RouteDistinguisher.IsNull() {
					*stateField = plan.RouteDistinguisher
				} else {
					*stateField = types.StringNull()
				}
			case "route_target_import":
				if !plan.RouteTargetImport.IsNull() {
					*stateField = plan.RouteTargetImport
				} else {
					*stateField = types.StringNull()
				}
			case "route_target_export":
				if !plan.RouteTargetExport.IsNull() {
					*stateField = plan.RouteTargetExport
				} else {
					*stateField = types.StringNull()
				}
			case "import_route_map":
				if !plan.ImportRouteMap.IsNull() {
					*stateField = plan.ImportRouteMap
				} else {
					*stateField = types.StringNull()
				}
			case "import_route_map_ref_type_":
				if !plan.ImportRouteMapRefType.IsNull() {
					*stateField = plan.ImportRouteMapRefType
				} else {
					*stateField = types.StringNull()
				}
			case "export_route_map":
				if !plan.ExportRouteMap.IsNull() {
					*stateField = plan.ExportRouteMap
				} else {
					*stateField = types.StringNull()
				}
			case "export_route_map_ref_type_":
				if !plan.ExportRouteMapRefType.IsNull() {
					*stateField = plan.ExportRouteMapRefType
				} else {
					*stateField = types.StringNull()
				}
			}
		} else {
			*stateField = types.StringNull()
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
	} else {
		state.RouteTenants = nil
	}

	return state
}
