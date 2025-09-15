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

	if bulkOpsMgr != nil && bulkOpsMgr.HasPendingOrRecentServiceOperations() {
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

	// Check for auto-assigned VNI modifications, but allow when VNI is unknown (expecting API to set it)
	if !plan.Vni.IsUnknown() && !plan.Vni.Equal(state.Vni) && !plan.VniAutoAssigned.IsNull() && plan.VniAutoAssigned.ValueBool() {
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
	serviceReq := &openapi.ConfigPutRequestServiceServiceName{}
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
	} else if len(state.ObjectProperties) > 0 {
		serviceReq.ObjectProperties = nil
		hasChanges = true
	}

	if !plan.Enable.Equal(state.Enable) {
		serviceReq.Enable = openapi.PtrBool(plan.Enable.ValueBool())
		hasChanges = true
	}

	if !plan.Vlan.IsUnknown() && !plan.Vlan.Equal(state.Vlan) {
		if !plan.Vlan.IsNull() {
			vlanVal := int32(plan.Vlan.ValueInt64())
			serviceReq.Vlan = *openapi.NewNullableInt32(&vlanVal)
		} else {
			serviceReq.Vlan = *openapi.NewNullableInt32(nil)
		}
		hasChanges = true
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
	operationID := bulkOpsMgr.AddServiceDelete(ctx, name)
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
		state.ObjectProperties = []verityServiceObjectPropertiesModel{
			{
				Group: types.StringValue(fmt.Sprintf("%v", op["group"])),
			},
		}
	} else if plan != nil && len(plan.ObjectProperties) > 0 {
		state.ObjectProperties = plan.ObjectProperties
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

	if val, ok := serviceData["anycast_ip_mask"].(string); ok {
		state.AnycastIpMask = types.StringValue(val)
	} else if plan != nil && !plan.AnycastIpMask.IsNull() {
		state.AnycastIpMask = plan.AnycastIpMask
	} else {
		state.AnycastIpMask = types.StringNull()
	}

	if val, ok := serviceData["dhcp_server_ip"].(string); ok {
		state.DhcpServerIp = types.StringValue(val)
	} else if plan != nil && !plan.DhcpServerIp.IsNull() {
		state.DhcpServerIp = plan.DhcpServerIp
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
