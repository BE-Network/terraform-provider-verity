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
	Name             types.String                         `tfsdk:"name"`
	Enable           types.Bool                           `tfsdk:"enable"`
	ObjectProperties []verityServiceObjectPropertiesModel `tfsdk:"object_properties"`
	Vlan             types.Int64                          `tfsdk:"vlan"`
	Vni              types.Int64                          `tfsdk:"vni"`
	VniAutoAssigned  types.Bool                           `tfsdk:"vni_auto_assigned_"`
	Tenant           types.String                         `tfsdk:"tenant"`
	TenantRefType    types.String                         `tfsdk:"tenant_ref_type_"`
	AnycastIpMask    types.String                         `tfsdk:"anycast_ip_mask"`
	DhcpServerIp     types.String                         `tfsdk:"dhcp_server_ip"`
	Mtu              types.Int64                          `tfsdk:"mtu"`
}

type verityServiceObjectPropertiesModel struct {
	Group types.String `tfsdk:"group"`
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
			"anycast_ip_mask": schema.StringAttribute{
				Description: "Static anycast gateway address for service",
				Optional:    true,
			},
			"dhcp_server_ip": schema.StringAttribute{
				Description: "IP address(s) of the DHCP server for service. May have up to four separated by commas.",
				Optional:    true,
			},
			"mtu": schema.Int64Attribute{
				Description: "MTU (Maximum Transmission Unit) - the size used by a switch to determine when large packets must be broken up for delivery.",
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

	serviceReq := openapi.NewConfigPutRequestServiceServiceName()
	serviceReq.Name = openapi.PtrString(name)
	if !plan.Enable.IsNull() {
		serviceReq.Enable = openapi.PtrBool(plan.Enable.ValueBool())
	}

	if len(plan.ObjectProperties) > 0 {
		objProps := openapi.ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties{}
		op := plan.ObjectProperties[0]
		if !op.Group.IsNull() {
			objProps.Group = openapi.PtrString(op.Group.ValueString())
		} else {
			objProps.Group = nil
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
	if !plan.Vni.IsNull() {
		vniVal := int32(plan.Vni.ValueInt64())
		serviceReq.Vni = *openapi.NewNullableInt32(&vniVal)
	} else {
		serviceReq.Vni = *openapi.NewNullableInt32(nil)
	}
	if !plan.VniAutoAssigned.IsNull() {
		serviceReq.VniAutoAssigned = openapi.PtrBool(plan.VniAutoAssigned.ValueBool())
	}
	if !plan.Tenant.IsNull() {
		serviceReq.Tenant = openapi.PtrString(plan.Tenant.ValueString())
	}
	if !plan.TenantRefType.IsNull() {
		serviceReq.TenantRefType = openapi.PtrString(plan.TenantRefType.ValueString())
	}
	if !plan.AnycastIpMask.IsNull() {
		serviceReq.AnycastIpMask = openapi.PtrString(plan.AnycastIpMask.ValueString())
	}
	if !plan.DhcpServerIp.IsNull() {
		serviceReq.DhcpServerIp = openapi.PtrString(plan.DhcpServerIp.ValueString())
	}
	if !plan.Mtu.IsNull() {
		mtuVal := int32(plan.Mtu.ValueInt64())
		serviceReq.Mtu = *openapi.NewNullableInt32(&mtuVal)
	} else {
		serviceReq.Mtu = *openapi.NewNullableInt32(nil)
	}

	provCtx := r.provCtx
	bulkMgr := provCtx.bulkOpsMgr
	operationID := bulkMgr.AddServicePut(ctx, name, *serviceReq)

	provCtx.NotifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for service creation operation %s to complete", operationID))
	if err := bulkMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Create Service",
			fmt.Sprintf("Error creating service %s: %v", name, err),
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Service %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "services")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
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

	if bulkOpsMgr != nil && bulkOpsMgr.HasPendingOrRecentServiceOperations() {
		tflog.Info(ctx, fmt.Sprintf("Skipping service %s verification - trusting recent successful API operation", serviceName))
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
		resp.Diagnostics.AddError(
			"Failed to Read Service",
			fmt.Sprintf("Error reading service %s: %v", serviceName, err),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for service with ID: %s", serviceName))
	var serviceData map[string]interface{}
	exists := false

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

	state.Name = types.StringValue(fmt.Sprintf("%v", serviceData["name"]))

	if val, ok := serviceData["enable"].(bool); ok {
		state.Enable = types.BoolValue(val)
	} else {
		state.Enable = types.BoolNull()
	}

	if op, ok := serviceData["object_properties"].(map[string]interface{}); ok {
		state.ObjectProperties = []verityServiceObjectPropertiesModel{
			{
				Group: types.StringValue(fmt.Sprintf("%v", op["group"])),
			},
		}
	} else {
		state.ObjectProperties = nil
	}

	if val, ok := serviceData["vlan"]; ok {
		switch v := val.(type) {
		case float64:
			state.Vlan = types.Int64Value(int64(v))
		case int:
			state.Vlan = types.Int64Value(int64(v))
		default:
			state.Vlan = types.Int64Null()
		}
	} else {
		state.Vlan = types.Int64Null()
	}

	if val, ok := serviceData["vni"]; ok {
		switch v := val.(type) {
		case float64:
			state.Vni = types.Int64Value(int64(v))
		case int:
			state.Vni = types.Int64Value(int64(v))
		default:
			state.Vni = types.Int64Null()
		}
	} else {
		state.Vni = types.Int64Null()
	}

	if val, ok := serviceData["vni_auto_assigned_"].(bool); ok {
		state.VniAutoAssigned = types.BoolValue(val)
	} else {
		state.VniAutoAssigned = types.BoolNull()
	}

	if val, ok := serviceData["tenant"].(string); ok {
		state.Tenant = types.StringValue(val)
	} else {
		state.Tenant = types.StringNull()
	}

	if val, ok := serviceData["tenant_ref_type_"].(string); ok {
		state.TenantRefType = types.StringValue(val)
	} else {
		state.TenantRefType = types.StringNull()
	}

	if val, ok := serviceData["anycast_ip_mask"].(string); ok {
		state.AnycastIpMask = types.StringValue(val)
	} else {
		state.AnycastIpMask = types.StringNull()
	}

	if val, ok := serviceData["dhcp_server_ip"].(string); ok {
		state.DhcpServerIp = types.StringValue(val)
	} else {
		state.DhcpServerIp = types.StringNull()
	}

	if val, ok := serviceData["mtu"]; ok {
		switch v := val.(type) {
		case float64:
			state.Mtu = types.Int64Value(int64(v))
		case int:
			state.Mtu = types.Int64Value(int64(v))
		default:
			state.Mtu = types.Int64Null()
		}
	} else {
		state.Mtu = types.Int64Null()
	}

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

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Authenticate",
			fmt.Sprintf("Error authenticating with API: %v", err),
		)
		return
	}

	name := plan.Name.ValueString()
	serviceReq := openapi.NewConfigPutRequestServiceServiceName()
	serviceReq.Name = openapi.PtrString(name)
	hasChanges := false

	if len(plan.ObjectProperties) > 0 {
		if len(state.ObjectProperties) == 0 || !plan.ObjectProperties[0].Group.Equal(state.ObjectProperties[0].Group) {
			objProps := openapi.ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties{}
			if !plan.ObjectProperties[0].Group.IsNull() {
				objProps.Group = openapi.PtrString(plan.ObjectProperties[0].Group.ValueString())
			} else {
				objProps.Group = nil
			}
			serviceReq.ObjectProperties = &objProps
			hasChanges = true
		}
	} else {
		serviceReq.ObjectProperties = nil
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
	} else if !state.Vlan.IsNull() {
		vlanVal := int32(state.Vlan.ValueInt64())
		serviceReq.Vlan = *openapi.NewNullableInt32(&vlanVal)
	} else {
		serviceReq.Vlan = *openapi.NewNullableInt32(nil)
	}
	if !plan.Vni.Equal(state.Vni) {
		if !plan.Vni.IsNull() {
			vniVal := int32(plan.Vni.ValueInt64())
			serviceReq.Vni = *openapi.NewNullableInt32(&vniVal)
		} else {
			serviceReq.Vni = *openapi.NewNullableInt32(nil)
		}
		hasChanges = true
	} else if !state.Vni.IsNull() {
		vniVal := int32(state.Vni.ValueInt64())
		serviceReq.Vni = *openapi.NewNullableInt32(&vniVal)
	} else {
		serviceReq.Vni = *openapi.NewNullableInt32(nil)
	}

	if !plan.VniAutoAssigned.Equal(state.VniAutoAssigned) {
		serviceReq.VniAutoAssigned = openapi.PtrBool(plan.VniAutoAssigned.ValueBool())
		hasChanges = true
	}

	if !plan.Tenant.Equal(state.Tenant) {
		serviceReq.Tenant = openapi.PtrString(plan.Tenant.ValueString())
		hasChanges = true
	}
	if !plan.TenantRefType.Equal(state.TenantRefType) {
		serviceReq.TenantRefType = openapi.PtrString(plan.TenantRefType.ValueString())
		hasChanges = true
	}
	if !plan.AnycastIpMask.Equal(state.AnycastIpMask) {
		serviceReq.AnycastIpMask = openapi.PtrString(plan.AnycastIpMask.ValueString())
		hasChanges = true
	}
	if !plan.DhcpServerIp.Equal(state.DhcpServerIp) {
		serviceReq.DhcpServerIp = openapi.PtrString(plan.DhcpServerIp.ValueString())
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
	} else if !state.Mtu.IsNull() {
		mtuVal := int32(state.Mtu.ValueInt64())
		serviceReq.Mtu = *openapi.NewNullableInt32(&mtuVal)
	} else {
		serviceReq.Mtu = *openapi.NewNullableInt32(nil)
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	bulkOpsMgr := r.provCtx.bulkOpsMgr
	operationID := bulkOpsMgr.AddServicePatch(ctx, name, *serviceReq)
	r.provCtx.NotifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for service update operation %s to complete", operationID))
	if err := bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Update Service",
			fmt.Sprintf("Error updating service %s: %v", name, err),
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Service %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "services")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
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
	operationID := bulkOpsMgr.AddServiceDelete(ctx, name)
	r.provCtx.NotifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for service deletion operation %s to complete", operationID))
	if err := bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Delete Service",
			fmt.Sprintf("Error deleting service %s: %v", name, err),
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
