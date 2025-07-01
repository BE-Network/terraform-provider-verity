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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-verity/internal/utils"
	"terraform-provider-verity/openapi"
)

var (
	_ resource.Resource                = &verityServicePortProfileResource{}
	_ resource.ResourceWithConfigure   = &verityServicePortProfileResource{}
	_ resource.ResourceWithImportState = &verityServicePortProfileResource{}
)

func NewVerityServicePortProfileResource() resource.Resource {
	return &verityServicePortProfileResource{}
}

type verityServicePortProfileResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *utils.BulkOperationManager
	notifyOperationAdded func()
}

type verityServicePortProfileResourceModel struct {
	Name              types.String                                    `tfsdk:"name"`
	Enable            types.Bool                                      `tfsdk:"enable"`
	PortType          types.String                                    `tfsdk:"port_type"`
	TlsLimitIn        types.Int64                                     `tfsdk:"tls_limit_in"`
	TlsService        types.String                                    `tfsdk:"tls_service"`
	TlsServiceRefType types.String                                    `tfsdk:"tls_service_ref_type_"`
	TrustedPort       types.Bool                                      `tfsdk:"trusted_port"`
	IpMask            types.String                                    `tfsdk:"ip_mask"`
	Services          []verityServicePortProfileServiceModel          `tfsdk:"services"`
	ObjectProperties  []verityServicePortProfileObjectPropertiesModel `tfsdk:"object_properties"`
}

type verityServicePortProfileServiceModel struct {
	RowNumEnable         types.Bool   `tfsdk:"row_num_enable"`
	RowNumService        types.String `tfsdk:"row_num_service"`
	RowNumServiceRefType types.String `tfsdk:"row_num_service_ref_type_"`
	RowNumExternalVlan   types.Int64  `tfsdk:"row_num_external_vlan"`
	RowNumLimitIn        types.Int64  `tfsdk:"row_num_limit_in"`
	RowNumLimitOut       types.Int64  `tfsdk:"row_num_limit_out"`
	Index                types.Int64  `tfsdk:"index"`
}

type verityServicePortProfileObjectPropertiesModel struct {
	OnSummary      types.Bool   `tfsdk:"on_summary"`
	PortMonitoring types.String `tfsdk:"port_monitoring"`
	Group          types.String `tfsdk:"group"`
}

func (r *verityServicePortProfileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_service_port_profile"
}

func (r *verityServicePortProfileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityServicePortProfileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity Service Port Profile",
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
			"port_type": schema.StringAttribute{
				Description: "Determines what Service are provisioned on the port and if those Services are propagated upstream",
				Optional:    true,
			},
			"tls_limit_in": schema.Int64Attribute{
				Description: "Speed of ingress (Mbps) for TLS (Transparent LAN Service)",
				Optional:    true,
			},
			"tls_service": schema.StringAttribute{
				Description: "Service used for TLS (Transparent LAN Service)",
				Optional:    true,
			},
			"tls_service_ref_type_": schema.StringAttribute{
				Description: "Object type for tls_service field",
				Optional:    true,
			},
			"trusted_port": schema.BoolAttribute{
				Description: "Trusted Ports do not participate in IP Source Guard, Dynamic ARP Inspection, nor DHCP Snooping, meaning all packets are forwarded without any checks.",
				Optional:    true,
			},
			"ip_mask": schema.StringAttribute{
				Description: "IP/Mask",
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"services": schema.ListNestedBlock{
				Description: "Service configurations",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"row_num_enable": schema.BoolAttribute{
							Description: "Enable row",
							Optional:    true,
						},
						"row_num_service": schema.StringAttribute{
							Description: "Connect a Service",
							Optional:    true,
						},
						"row_num_service_ref_type_": schema.StringAttribute{
							Description: "Object type for row_num_service field",
							Optional:    true,
						},
						"row_num_external_vlan": schema.Int64Attribute{
							Description: "Choose an external vlan",
							Optional:    true,
						},
						"row_num_limit_in": schema.Int64Attribute{
							Description: "Speed of ingress (Mbps)",
							Optional:    true,
						},
						"row_num_limit_out": schema.Int64Attribute{
							Description: "Speed of egress (Mbps)",
							Optional:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object. Zero if you want to add an object to the list.",
							Optional:    true,
						},
					},
				},
			},
			"object_properties": schema.ListNestedBlock{
				Description: "Object properties for the service port profile",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"on_summary": schema.BoolAttribute{
							Description: "Show on the summary view",
							Optional:    true,
						},
						"port_monitoring": schema.StringAttribute{
							Description: "Defines importance of Link Down on this port",
							Optional:    true,
						},
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

func (r *verityServicePortProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityServicePortProfileResourceModel
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
	sppProps := openapi.ConfigPutRequestServicePortProfileServicePortProfileName{}
	sppProps.Name = openapi.PtrString(name)

	if !plan.Enable.IsNull() {
		sppProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
	}
	if !plan.PortType.IsNull() {
		sppProps.PortType = openapi.PtrString(plan.PortType.ValueString())
	}
	if !plan.TlsLimitIn.IsNull() {
		val := int32(plan.TlsLimitIn.ValueInt64())
		sppProps.TlsLimitIn = *openapi.NewNullableInt32(&val)
	} else {
		sppProps.TlsLimitIn = *openapi.NewNullableInt32(nil)
	}
	if !plan.TlsService.IsNull() {
		sppProps.TlsService = openapi.PtrString(plan.TlsService.ValueString())
	}
	if !plan.TlsServiceRefType.IsNull() {
		sppProps.TlsServiceRefType = openapi.PtrString(plan.TlsServiceRefType.ValueString())
	}
	if !plan.TrustedPort.IsNull() {
		sppProps.TrustedPort = openapi.PtrBool(plan.TrustedPort.ValueBool())
	}
	if !plan.IpMask.IsNull() {
		sppProps.IpMask = openapi.PtrString(plan.IpMask.ValueString())
	}

	if len(plan.Services) > 0 {
		services := make([]openapi.ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner, len(plan.Services))
		for i, service := range plan.Services {
			serviceItem := openapi.ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner{}
			if !service.RowNumEnable.IsNull() {
				serviceItem.RowNumEnable = openapi.PtrBool(service.RowNumEnable.ValueBool())
			}
			if !service.RowNumService.IsNull() {
				serviceItem.RowNumService = openapi.PtrString(service.RowNumService.ValueString())
			}
			if !service.RowNumServiceRefType.IsNull() {
				serviceItem.RowNumServiceRefType = openapi.PtrString(service.RowNumServiceRefType.ValueString())
			}
			if !service.RowNumExternalVlan.IsNull() {
				val := int32(service.RowNumExternalVlan.ValueInt64())
				serviceItem.RowNumExternalVlan = *openapi.NewNullableInt32(&val)
			} else {
				serviceItem.RowNumExternalVlan = *openapi.NewNullableInt32(nil)
			}
			if !service.RowNumLimitIn.IsNull() {
				val := int32(service.RowNumLimitIn.ValueInt64())
				serviceItem.RowNumLimitIn = *openapi.NewNullableInt32(&val)
			} else {
				serviceItem.RowNumLimitIn = *openapi.NewNullableInt32(nil)
			}
			if !service.RowNumLimitOut.IsNull() {
				val := int32(service.RowNumLimitOut.ValueInt64())
				serviceItem.RowNumLimitOut = *openapi.NewNullableInt32(&val)
			} else {
				serviceItem.RowNumLimitOut = *openapi.NewNullableInt32(nil)
			}
			if !service.Index.IsNull() {
				serviceItem.Index = openapi.PtrInt32(int32(service.Index.ValueInt64()))
			}
			services[i] = serviceItem
		}
		sppProps.Services = services
	}

	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.ConfigPutRequestServicePortProfileServicePortProfileNameObjectProperties{}
		if !op.OnSummary.IsNull() {
			objProps.OnSummary = openapi.PtrBool(op.OnSummary.ValueBool())
		}
		if !op.PortMonitoring.IsNull() {
			objProps.PortMonitoring = openapi.PtrString(op.PortMonitoring.ValueString())
		}
		if !op.Group.IsNull() {
			objProps.Group = openapi.PtrString(op.Group.ValueString())
		}
		sppProps.ObjectProperties = &objProps
	}

	operationID := r.bulkOpsMgr.AddServicePortProfilePut(ctx, name, sppProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for service port profile creation operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Create Service Port Profile %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Service Port Profile %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "service_port_profiles")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
}

func (r *verityServicePortProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityServicePortProfileResourceModel
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

	sppName := state.Name.ValueString()

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentServicePortProfileOperations() {
		tflog.Info(ctx, fmt.Sprintf("Skipping Service Port Profile %s verification – trusting recent successful API operation", sppName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching Service Port Profiles for verification of %s", sppName))

	type ServicePortProfileResponse struct {
		ServicePortProfile map[string]interface{} `json:"service_port_profile"`
	}

	var result ServicePortProfileResponse
	var err error
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		sppData, fetchErr := getCachedResponse(ctx, r.provCtx, "service_port_profiles", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch Service Port Profiles")
			respAPI, err := r.client.ServicePortProfilesAPI.ServiceportprofilesGet(ctx).Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading Service Port Profiles: %v", err)
			}
			defer respAPI.Body.Close()

			var res ServicePortProfileResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return nil, fmt.Errorf("failed to decode Service Port Profiles response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d Service Port Profiles", len(res.ServicePortProfile)))
			return res, nil
		})
		if fetchErr != nil {
			err = fetchErr
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch Service Port Profiles on attempt %d, retrying in %v", attempt+1, sleepTime))
			time.Sleep(sleepTime)
			continue
		}
		result = sppData.(ServicePortProfileResponse)
		break
	}
	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Service Port Profile %s", sppName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for Service Port Profile with ID: %s", sppName))
	var sppData map[string]interface{}
	exists := false

	if data, ok := result.ServicePortProfile[sppName].(map[string]interface{}); ok {
		sppData = data
		exists = true
		tflog.Debug(ctx, fmt.Sprintf("Found Service Port Profile directly by ID: %s", sppName))
	} else {
		for apiName, s := range result.ServicePortProfile {
			servicePortProfile, ok := s.(map[string]interface{})
			if !ok {
				continue
			}

			if name, ok := servicePortProfile["name"].(string); ok && name == sppName {
				sppData = servicePortProfile
				sppName = apiName
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found Service Port Profile with name '%s' under API key '%s'", name, apiName))
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Service Port Profile with ID '%s' not found in API response", sppName))
		resp.State.RemoveResource(ctx)
		return
	}

	state.Name = types.StringValue(fmt.Sprintf("%v", sppData["name"]))

	if enable, ok := sppData["enable"].(bool); ok {
		state.Enable = types.BoolValue(enable)
	} else {
		state.Enable = types.BoolNull()
	}

	if trustedPort, ok := sppData["trusted_port"].(bool); ok {
		state.TrustedPort = types.BoolValue(trustedPort)
	} else {
		state.TrustedPort = types.BoolNull()
	}

	if portType, ok := sppData["port_type"].(string); ok {
		state.PortType = types.StringValue(portType)
	} else {
		state.PortType = types.StringNull()
	}

	if tlsService, ok := sppData["tls_service"].(string); ok {
		state.TlsService = types.StringValue(tlsService)
	} else {
		state.TlsService = types.StringNull()
	}

	if tlsServiceRefType, ok := sppData["tls_service_ref_type_"].(string); ok {
		state.TlsServiceRefType = types.StringValue(tlsServiceRefType)
	} else {
		state.TlsServiceRefType = types.StringNull()
	}

	if ipMask, ok := sppData["ip_mask"].(string); ok {
		state.IpMask = types.StringValue(ipMask)
	} else {
		state.IpMask = types.StringNull()
	}

	if tlsLimitIn, ok := sppData["tls_limit_in"]; ok && tlsLimitIn != nil {
		switch v := tlsLimitIn.(type) {
		case int:
			state.TlsLimitIn = types.Int64Value(int64(v))
		case int32:
			state.TlsLimitIn = types.Int64Value(int64(v))
		case int64:
			state.TlsLimitIn = types.Int64Value(v)
		case float64:
			state.TlsLimitIn = types.Int64Value(int64(v))
		case string:
			if intVal, err := strconv.ParseInt(v, 10, 64); err == nil {
				state.TlsLimitIn = types.Int64Value(intVal)
			} else {
				state.TlsLimitIn = types.Int64Null()
			}
		default:
			state.TlsLimitIn = types.Int64Null()
		}
	} else {
		state.TlsLimitIn = types.Int64Null()
	}

	if servicesArray, ok := sppData["services"].([]interface{}); ok && len(servicesArray) > 0 {
		var services []verityServicePortProfileServiceModel
		for _, s := range servicesArray {
			service, ok := s.(map[string]interface{})
			if !ok {
				continue
			}
			serviceModel := verityServicePortProfileServiceModel{}

			if enable, ok := service["row_num_enable"].(bool); ok {
				serviceModel.RowNumEnable = types.BoolValue(enable)
			} else {
				serviceModel.RowNumEnable = types.BoolNull()
			}

			if serviceName, ok := service["row_num_service"].(string); ok {
				serviceModel.RowNumService = types.StringValue(serviceName)
			} else {
				serviceModel.RowNumService = types.StringNull()
			}

			if serviceRefType, ok := service["row_num_service_ref_type_"].(string); ok {
				serviceModel.RowNumServiceRefType = types.StringValue(serviceRefType)
			} else {
				serviceModel.RowNumServiceRefType = types.StringNull()
			}

			intFields := map[string]*types.Int64{
				"row_num_external_vlan": &serviceModel.RowNumExternalVlan,
				"row_num_limit_in":      &serviceModel.RowNumLimitIn,
				"row_num_limit_out":     &serviceModel.RowNumLimitOut,
			}

			for apiKey, stateField := range intFields {
				if value, ok := service[apiKey]; ok && value != nil {
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

			if index, ok := service["index"]; ok && index != nil {
				if intVal, ok := index.(float64); ok {
					serviceModel.Index = types.Int64Value(int64(intVal))
				} else if intVal, ok := index.(int); ok {
					serviceModel.Index = types.Int64Value(int64(intVal))
				} else {
					serviceModel.Index = types.Int64Null()
				}
			} else {
				serviceModel.Index = types.Int64Null()
			}

			services = append(services, serviceModel)
		}
		state.Services = services
	} else {
		state.Services = nil
	}

	if objProps, ok := sppData["object_properties"].(map[string]interface{}); ok {
		op := verityServicePortProfileObjectPropertiesModel{}
		if onSummary, ok := objProps["on_summary"].(bool); ok {
			op.OnSummary = types.BoolValue(onSummary)
		} else {
			op.OnSummary = types.BoolNull()
		}
		if portMonitoring, ok := objProps["port_monitoring"].(string); ok {
			op.PortMonitoring = types.StringValue(portMonitoring)
		} else {
			op.PortMonitoring = types.StringNull()
		}
		if group, ok := objProps["group"].(string); ok {
			op.Group = types.StringValue(group)
		} else {
			op.Group = types.StringNull()
		}
		state.ObjectProperties = []verityServicePortProfileObjectPropertiesModel{op}
	} else {
		state.ObjectProperties = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityServicePortProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityServicePortProfileResourceModel

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
	sppProps := openapi.ConfigPutRequestServicePortProfileServicePortProfileName{}
	hasChanges := false

	if !plan.Name.Equal(state.Name) {
		sppProps.Name = openapi.PtrString(name)
		hasChanges = true
	}

	if !plan.Enable.Equal(state.Enable) {
		sppProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
		hasChanges = true
	}
	if !plan.PortType.Equal(state.PortType) {
		sppProps.PortType = openapi.PtrString(plan.PortType.ValueString())
		hasChanges = true
	}
	if !plan.TlsService.Equal(state.TlsService) {
		sppProps.TlsService = openapi.PtrString(plan.TlsService.ValueString())
		hasChanges = true
	}
	if !plan.TlsServiceRefType.Equal(state.TlsServiceRefType) {
		sppProps.TlsServiceRefType = openapi.PtrString(plan.TlsServiceRefType.ValueString())
		hasChanges = true
	}
	if !plan.TrustedPort.Equal(state.TrustedPort) {
		sppProps.TrustedPort = openapi.PtrBool(plan.TrustedPort.ValueBool())
		hasChanges = true
	}
	if !plan.IpMask.Equal(state.IpMask) {
		sppProps.IpMask = openapi.PtrString(plan.IpMask.ValueString())
		hasChanges = true
	}

	if !plan.TlsLimitIn.Equal(state.TlsLimitIn) {
		if !plan.TlsLimitIn.IsNull() {
			val := int32(plan.TlsLimitIn.ValueInt64())
			sppProps.TlsLimitIn = *openapi.NewNullableInt32(&val)
		} else {
			sppProps.TlsLimitIn = *openapi.NewNullableInt32(nil)
		}
		hasChanges = true
	}

	if !r.equalServiceArrays(plan.Services, state.Services) {
		services := make([]openapi.ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner, len(plan.Services))
		for i, service := range plan.Services {
			serviceItem := openapi.ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner{}
			if !service.RowNumEnable.IsNull() {
				serviceItem.RowNumEnable = openapi.PtrBool(service.RowNumEnable.ValueBool())
			}
			if !service.RowNumService.IsNull() {
				serviceItem.RowNumService = openapi.PtrString(service.RowNumService.ValueString())
			}
			if !service.RowNumServiceRefType.IsNull() {
				serviceItem.RowNumServiceRefType = openapi.PtrString(service.RowNumServiceRefType.ValueString())
			}
			if !service.RowNumExternalVlan.IsNull() {
				val := int32(service.RowNumExternalVlan.ValueInt64())
				serviceItem.RowNumExternalVlan = *openapi.NewNullableInt32(&val)
			} else {
				serviceItem.RowNumExternalVlan = *openapi.NewNullableInt32(nil)
			}
			if !service.RowNumLimitIn.IsNull() {
				val := int32(service.RowNumLimitIn.ValueInt64())
				serviceItem.RowNumLimitIn = *openapi.NewNullableInt32(&val)
			} else {
				serviceItem.RowNumLimitIn = *openapi.NewNullableInt32(nil)
			}
			if !service.RowNumLimitOut.IsNull() {
				val := int32(service.RowNumLimitOut.ValueInt64())
				serviceItem.RowNumLimitOut = *openapi.NewNullableInt32(&val)
			} else {
				serviceItem.RowNumLimitOut = *openapi.NewNullableInt32(nil)
			}
			if !service.Index.IsNull() {
				serviceItem.Index = openapi.PtrInt32(int32(service.Index.ValueInt64()))
			}
			services[i] = serviceItem
		}
		sppProps.Services = services
		hasChanges = true
	}

	if len(plan.ObjectProperties) > 0 {
		if len(state.ObjectProperties) == 0 || !r.equalObjectProperties(plan.ObjectProperties[0], state.ObjectProperties[0]) {
			op := plan.ObjectProperties[0]
			objProps := openapi.ConfigPutRequestServicePortProfileServicePortProfileNameObjectProperties{}
			if !op.OnSummary.IsNull() {
				objProps.OnSummary = openapi.PtrBool(op.OnSummary.ValueBool())
			}
			if !op.PortMonitoring.IsNull() {
				objProps.PortMonitoring = openapi.PtrString(op.PortMonitoring.ValueString())
			}
			if !op.Group.IsNull() {
				objProps.Group = openapi.PtrString(op.Group.ValueString())
			}
			sppProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	operationID := r.bulkOpsMgr.AddServicePortProfilePatch(ctx, name, sppProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for Service Port Profile update operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update Service Port Profile %s", name))...,
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Service Port Profile %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "service_port_profiles")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *verityServicePortProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityServicePortProfileResourceModel
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
	operationID := r.bulkOpsMgr.AddServicePortProfileDelete(ctx, name)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for Service Port Profile deletion operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Delete Service Port Profile %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Service Port Profile %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "service_port_profiles")
	resp.State.RemoveResource(ctx)
}

func (r *verityServicePortProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func (r *verityServicePortProfileResource) equalServiceArrays(a, b []verityServicePortProfileServiceModel) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !a[i].RowNumEnable.Equal(b[i].RowNumEnable) ||
			!a[i].RowNumService.Equal(b[i].RowNumService) ||
			!a[i].RowNumServiceRefType.Equal(b[i].RowNumServiceRefType) ||
			!a[i].RowNumExternalVlan.Equal(b[i].RowNumExternalVlan) ||
			!a[i].RowNumLimitIn.Equal(b[i].RowNumLimitIn) ||
			!a[i].RowNumLimitOut.Equal(b[i].RowNumLimitOut) ||
			!a[i].Index.Equal(b[i].Index) {
			return false
		}
	}
	return true
}

func (r *verityServicePortProfileResource) equalObjectProperties(a, b verityServicePortProfileObjectPropertiesModel) bool {
	return a.OnSummary.Equal(b.OnSummary) &&
		a.PortMonitoring.Equal(b.PortMonitoring) &&
		a.Group.Equal(b.Group)
}
