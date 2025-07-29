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
	_ resource.Resource                = &verityEthPortProfileResource{}
	_ resource.ResourceWithConfigure   = &verityEthPortProfileResource{}
	_ resource.ResourceWithImportState = &verityEthPortProfileResource{}
)

func NewVerityEthPortProfileResource() resource.Resource {
	return &verityEthPortProfileResource{}
}

type verityEthPortProfileResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *utils.BulkOperationManager
	notifyOperationAdded func()
}

type verityEthPortProfileResourceModel struct {
	Name               types.String                                `tfsdk:"name"`
	Enable             types.Bool                                  `tfsdk:"enable"`
	TenantSliceManaged types.Bool                                  `tfsdk:"tenant_slice_managed"`
	ObjectProperties   []verityEthPortProfileObjectPropertiesModel `tfsdk:"object_properties"`
	Services           []servicesModel                             `tfsdk:"services"`
	Tls                types.Bool                                  `tfsdk:"tls"`
	TlsService         types.String                                `tfsdk:"tls_service"`
	TlsServiceRefType  types.String                                `tfsdk:"tls_service_ref_type_"`
	TrustedPort        types.Bool                                  `tfsdk:"trusted_port"`
}

type verityEthPortProfileObjectPropertiesModel struct {
	Group          types.String `tfsdk:"group"`
	PortMonitoring types.String `tfsdk:"port_monitoring"`
	SortByName     types.Bool   `tfsdk:"sort_by_name"`
	Label          types.String `tfsdk:"label"`
	Icon           types.String `tfsdk:"icon"`
}

type servicesModel struct {
	RowNumEnable           types.Bool   `tfsdk:"row_num_enable"`
	RowNumService          types.String `tfsdk:"row_num_service"`
	RowNumServiceRefType   types.String `tfsdk:"row_num_service_ref_type_"`
	RowNumExternalVlan     types.Int64  `tfsdk:"row_num_external_vlan"`
	Index                  types.Int64  `tfsdk:"index"`
	RowNumMacFilter        types.String `tfsdk:"row_num_mac_filter"`
	RowNumMacFilterRefType types.String `tfsdk:"row_num_mac_filter_ref_type_"`
	RowNumLanIptv          types.String `tfsdk:"row_num_lan_iptv"`
}

func (r *verityEthPortProfileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_eth_port_profile"
}

func (r *verityEthPortProfileResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityEthPortProfileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an Ethernet Port Profile",
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
			"tenant_slice_managed": schema.BoolAttribute{
				Description: "Profiles that Tenant Slice creates and manages",
				Optional:    true,
				Computed:    true,
			},
			"tls": schema.BoolAttribute{
				Description: "Transparent LAN Service Trunk",
				Optional:    true,
			},
			"tls_service": schema.StringAttribute{
				Description: "Choose a Service supporting Transparent LAN Service",
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
		},
		Blocks: map[string]schema.Block{
			"object_properties": schema.ListNestedBlock{
				Description: "Object properties for the profile",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"group": schema.StringAttribute{
							Description: "Group",
							Optional:    true,
						},
						"port_monitoring": schema.StringAttribute{
							Description: "Defines importance of Link Down on this port",
							Optional:    true,
						},
						"sort_by_name": schema.BoolAttribute{
							Description: "Choose to sort by service name or by order of creation",
							Optional:    true,
						},
						"label": schema.StringAttribute{
							Description: "Port Label displayed ports provisioned with this Eth Port Profile but with no Port Label defined in the endpoint",
							Optional:    true,
						},
						"icon": schema.StringAttribute{
							Description: "Port Icon displayed ports provisioned with this Eth Port Profile but with no Port Icon defined in the endpoint",
							Optional:    true,
						},
					},
				},
			},
			"services": schema.ListNestedBlock{
				Description: "List of service configurations",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"row_num_enable": schema.BoolAttribute{
							Description: "Enable row",
							Optional:    true,
						},
						"row_num_service": schema.StringAttribute{
							Description: "Choose a Service to connect",
							Optional:    true,
						},
						"row_num_service_ref_type_": schema.StringAttribute{
							Description: "Object type for row_num_service field",
							Optional:    true,
						},
						"row_num_external_vlan": schema.Int64Attribute{
							Description: "Choose an external vlan. A value of 0 will make the VLAN untagged, while null will use service VLAN.",
							Optional:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object. Zero if you want to add an object to the list.",
							Optional:    true,
						},
						"row_num_mac_filter": schema.StringAttribute{
							Description: "Choose an access control list",
							Optional:    true,
						},
						"row_num_mac_filter_ref_type_": schema.StringAttribute{
							Description: "Object type for row_num_mac_filter field",
							Optional:    true,
						},
						"row_num_lan_iptv": schema.StringAttribute{
							Description: "Denotes a LAN or IPTV service",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func (r *verityEthPortProfileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityEthPortProfileResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := ensureAuthenticated(ctx, r.provCtx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Authenticate",
			fmt.Sprintf("Error authenticating with API: %s", err),
		)
		return
	}

	name := plan.Name.ValueString()

	ethPortName := openapi.ConfigPutRequestEthPortProfileEthPortProfileName{}
	ethPortName.Name = openapi.PtrString(name)
	ethPortName.Enable = openapi.PtrBool(plan.Enable.ValueBool())
	ethPortName.TenantSliceManaged = openapi.PtrBool(plan.TenantSliceManaged.ValueBool())

	if !plan.Tls.IsNull() {
		ethPortName.Tls = openapi.PtrBool(plan.Tls.ValueBool())
	}

	if !plan.TlsService.IsNull() {
		ethPortName.TlsService = openapi.PtrString(plan.TlsService.ValueString())
	}

	if !plan.TlsServiceRefType.IsNull() {
		ethPortName.TlsServiceRefType = openapi.PtrString(plan.TlsServiceRefType.ValueString())
	}

	if !plan.TrustedPort.IsNull() {
		ethPortName.TrustedPort = openapi.PtrBool(plan.TrustedPort.ValueBool())
	}

	if len(plan.ObjectProperties) > 0 {
		objProps := openapi.ConfigPutRequestEthPortProfileEthPortProfileNameObjectProperties{}
		objProp := plan.ObjectProperties[0]

		if !objProp.Group.IsNull() {
			objProps.Group = openapi.PtrString(objProp.Group.ValueString())
		} else {
			objProps.Group = nil
		}

		if !objProp.PortMonitoring.IsNull() {
			objProps.PortMonitoring = openapi.PtrString(objProp.PortMonitoring.ValueString())
		} else {
			objProps.PortMonitoring = nil
		}

		if !objProp.SortByName.IsNull() {
			objProps.SortByName = openapi.PtrBool(objProp.SortByName.ValueBool())
		}

		if !objProp.Label.IsNull() {
			objProps.Label = openapi.PtrString(objProp.Label.ValueString())
		}

		if !objProp.Icon.IsNull() {
			objProps.Icon = openapi.PtrString(objProp.Icon.ValueString())
		}

		ethPortName.ObjectProperties = &objProps
	}

	if len(plan.Services) > 0 {
		services := make([]openapi.ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner, len(plan.Services))

		for i, service := range plan.Services {
			s := openapi.ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner{}

			if !service.RowNumEnable.IsNull() {
				s.RowNumEnable = openapi.PtrBool(service.RowNumEnable.ValueBool())
			}

			if !service.RowNumService.IsNull() {
				s.RowNumService = openapi.PtrString(service.RowNumService.ValueString())
			}

			if !service.RowNumServiceRefType.IsNull() {
				s.RowNumServiceRefType = openapi.PtrString(service.RowNumServiceRefType.ValueString())
			}

			if !service.RowNumExternalVlan.IsNull() {
				intVal := int32(service.RowNumExternalVlan.ValueInt64())
				s.RowNumExternalVlan = *openapi.NewNullableInt32(&intVal)
			} else {
				s.RowNumExternalVlan = *openapi.NewNullableInt32(nil)
			}

			if !service.Index.IsNull() {
				s.Index = openapi.PtrInt32(int32(service.Index.ValueInt64()))
			}

			if !service.RowNumMacFilter.IsNull() {
				s.RowNumMacFilter = openapi.PtrString(service.RowNumMacFilter.ValueString())
			}

			if !service.RowNumMacFilterRefType.IsNull() {
				s.RowNumMacFilterRefType = openapi.PtrString(service.RowNumMacFilterRefType.ValueString())
			}

			if !service.RowNumLanIptv.IsNull() {
				s.RowNumLanIptv = openapi.PtrString(service.RowNumLanIptv.ValueString())
			}

			services[i] = s
		}

		ethPortName.Services = services
	}

	operationID := r.bulkOpsMgr.AddPut(ctx, "ethportprofile", name, ethPortName)

	r.notifyOperationAdded()
	tflog.Debug(ctx, fmt.Sprintf("Waiting for eth port profile creation operation %s to complete", operationID))

	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Create Eth Port Profile %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Eth port profile %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "eth_port_profiles")
	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
}

func (r *verityEthPortProfileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityEthPortProfileResourceModel
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

	profileName := state.Name.ValueString()

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("ethportprofile") {
		tflog.Info(ctx, fmt.Sprintf("Skipping eth port profile %s verification - trusting recent successful API operation", profileName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("No recent eth port profile operations found, performing normal verification for %s", profileName))

	type EthPortProfileResponse struct {
		EthPortProfile map[string]map[string]interface{} `json:"eth_port_profile_"`
	}

	var result EthPortProfileResponse
	var err error
	maxRetries := 3

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch eth port profiles on attempt %d, retrying in %v", attempt, sleepTime))
			time.Sleep(sleepTime)
		}

		profilesData, fetchErr := getCachedResponse(ctx, r.provCtx, "eth_port_profiles", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch Ethernet port profiles")
			resp, err := r.client.EthPortProfilesAPI.EthportprofilesGet(ctx).Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading EthPort profiles: %v", err)
			}
			defer resp.Body.Close()

			var result EthPortProfileResponse
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				return nil, fmt.Errorf("error decoding EthPort profiles response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d Ethernet port profiles from API", len(result.EthPortProfile)))
			return result, nil
		})

		if fetchErr == nil {
			result = profilesData.(EthPortProfileResponse)
			break
		}
		err = fetchErr
	}

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, "Failed to Read Eth Port Profiles")...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for Ethernet port profile with ID: %s", profileName))

	profile, exists := result.EthPortProfile[profileName]
	if exists {
		tflog.Debug(ctx, fmt.Sprintf("Found Ethernet port profile directly by ID: %s", profileName))
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Profile not found directly by ID '%s', searching through all profiles", profileName))
		for apiName, p := range result.EthPortProfile {
			if name, ok := p["name"].(string); ok && name == profileName {
				profile = p
				profileName = apiName
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found Ethernet port profile with name '%s' under API key '%s'", name, apiName))
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Ethernet port profile with ID '%s' not found in API response", profileName))
		resp.State.RemoveResource(ctx)
		return
	}

	if name, ok := profile["name"].(string); ok {
		state.Name = types.StringValue(name)
	}

	if enable, ok := profile["enable"].(bool); ok {
		state.Enable = types.BoolValue(enable)
	}

	if tenantSliceManaged, ok := profile["tenant_slice_managed"].(bool); ok {
		state.TenantSliceManaged = types.BoolValue(tenantSliceManaged)
	}

	if tls, ok := profile["tls"].(bool); ok {
		state.Tls = types.BoolValue(tls)
	} else {
		state.Tls = types.BoolNull()
	}

	if tlsService, ok := profile["tls_service"].(string); ok {
		state.TlsService = types.StringValue(tlsService)
	} else {
		state.TlsService = types.StringNull()
	}

	if tlsServiceRefType, ok := profile["tls_service_ref_type_"].(string); ok {
		state.TlsServiceRefType = types.StringValue(tlsServiceRefType)
	} else {
		state.TlsServiceRefType = types.StringNull()
	}

	if trustedPort, ok := profile["trusted_port"].(bool); ok {
		state.TrustedPort = types.BoolValue(trustedPort)
	} else {
		state.TrustedPort = types.BoolNull()
	}

	if objProps, ok := profile["object_properties"].(map[string]interface{}); ok {
		var objProperties []verityEthPortProfileObjectPropertiesModel
		objProperty := verityEthPortProfileObjectPropertiesModel{}

		if group, ok := objProps["group"].(string); ok {
			objProperty.Group = types.StringValue(group)
		} else {
			objProperty.Group = types.StringNull()
		}

		if portMonitoring, ok := objProps["port_monitoring"].(string); ok {
			objProperty.PortMonitoring = types.StringValue(portMonitoring)
		} else {
			objProperty.PortMonitoring = types.StringNull()
		}

		if sortByName, ok := objProps["sort_by_name"].(bool); ok {
			objProperty.SortByName = types.BoolValue(sortByName)
		} else {
			objProperty.SortByName = types.BoolNull()
		}

		if label, ok := objProps["label"].(string); ok {
			objProperty.Label = types.StringValue(label)
		} else {
			objProperty.Label = types.StringNull()
		}

		if icon, ok := objProps["icon"].(string); ok {
			objProperty.Icon = types.StringValue(icon)
		} else {
			objProperty.Icon = types.StringNull()
		}

		objProperties = append(objProperties, objProperty)
		state.ObjectProperties = objProperties
	} else {
		state.ObjectProperties = nil
	}

	if services, ok := profile["services"].([]interface{}); ok {
		var servicesList []servicesModel

		for _, service := range services {
			if serviceMap, ok := service.(map[string]interface{}); ok {
				serviceModel := servicesModel{}

				if enable, ok := serviceMap["row_num_enable"].(bool); ok {
					serviceModel.RowNumEnable = types.BoolValue(enable)
				} else {
					serviceModel.RowNumEnable = types.BoolNull()
				}

				if service, ok := serviceMap["row_num_service"].(string); ok {
					serviceModel.RowNumService = types.StringValue(service)
				} else {
					serviceModel.RowNumService = types.StringNull()
				}

				if refType, ok := serviceMap["row_num_service_ref_type_"].(string); ok {
					serviceModel.RowNumServiceRefType = types.StringValue(refType)
				} else {
					serviceModel.RowNumServiceRefType = types.StringNull()
				}

				if vlan, ok := serviceMap["row_num_external_vlan"]; ok {
					if vlan == nil {
						serviceModel.RowNumExternalVlan = types.Int64Null()
					} else {
						switch v := vlan.(type) {
						case float64:
							serviceModel.RowNumExternalVlan = types.Int64Value(int64(v))
						case int:
							serviceModel.RowNumExternalVlan = types.Int64Value(int64(v))
						case int32:
							serviceModel.RowNumExternalVlan = types.Int64Value(int64(v))
						case int64:
							serviceModel.RowNumExternalVlan = types.Int64Value(v)
						case string:
							if intVal, err := strconv.ParseInt(v, 10, 64); err == nil {
								serviceModel.RowNumExternalVlan = types.Int64Value(intVal)
							} else {
								tflog.Warn(ctx, fmt.Sprintf("Cannot convert row_num_external_vlan value '%s' to integer", v))
								serviceModel.RowNumExternalVlan = types.Int64Null()
							}
						default:
							strVal := fmt.Sprintf("%v", v)
							if intVal, err := strconv.ParseInt(strVal, 10, 64); err == nil {
								serviceModel.RowNumExternalVlan = types.Int64Value(intVal)
							} else {
								tflog.Warn(ctx, fmt.Sprintf("Cannot convert row_num_external_vlan value '%v' to integer", v))
								serviceModel.RowNumExternalVlan = types.Int64Null()
							}
						}
					}
				} else {
					serviceModel.RowNumExternalVlan = types.Int64Null()
				}
				if index, ok := serviceMap["index"].(float64); ok {
					serviceModel.Index = types.Int64Value(int64(index))
				} else {
					serviceModel.Index = types.Int64Null()
				}

				if macFilter, ok := serviceMap["row_num_mac_filter"].(string); ok {
					serviceModel.RowNumMacFilter = types.StringValue(macFilter)
				} else {
					serviceModel.RowNumMacFilter = types.StringNull()
				}

				if macFilterRefType, ok := serviceMap["row_num_mac_filter_ref_type_"].(string); ok {
					serviceModel.RowNumMacFilterRefType = types.StringValue(macFilterRefType)
				} else {
					serviceModel.RowNumMacFilterRefType = types.StringNull()
				}

				if lanIptv, ok := serviceMap["row_num_lan_iptv"].(string); ok {
					serviceModel.RowNumLanIptv = types.StringValue(lanIptv)
				} else {
					serviceModel.RowNumLanIptv = types.StringNull()
				}

				servicesList = append(servicesList, serviceModel)
			}
		}

		state.Services = servicesList
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityEthPortProfileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityEthPortProfileResourceModel

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

	err := ensureAuthenticated(ctx, r.provCtx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Authenticate",
			fmt.Sprintf("Error authenticating with API: %s", err),
		)
		return
	}

	name := plan.Name.ValueString()
	ethPortName := openapi.ConfigPutRequestEthPortProfileEthPortProfileName{}

	hasChanges := false

	if !plan.Enable.Equal(state.Enable) {
		ethPortName.Enable = openapi.PtrBool(plan.Enable.ValueBool())
		hasChanges = true
	}

	if !plan.TenantSliceManaged.Equal(state.TenantSliceManaged) {
		ethPortName.TenantSliceManaged = openapi.PtrBool(plan.TenantSliceManaged.ValueBool())
		hasChanges = true
	}

	if !plan.Tls.Equal(state.Tls) {
		ethPortName.Tls = openapi.PtrBool(plan.Tls.ValueBool())
		hasChanges = true
	}

	if !plan.TlsService.Equal(state.TlsService) {
		ethPortName.TlsService = openapi.PtrString(plan.TlsService.ValueString())
		hasChanges = true
	}

	if !plan.TlsServiceRefType.Equal(state.TlsServiceRefType) {
		ethPortName.TlsServiceRefType = openapi.PtrString(plan.TlsServiceRefType.ValueString())
		hasChanges = true
	}

	if !plan.TrustedPort.Equal(state.TrustedPort) {
		ethPortName.TrustedPort = openapi.PtrBool(plan.TrustedPort.ValueBool())
		hasChanges = true
	}

	if len(plan.ObjectProperties) > 0 && (len(state.ObjectProperties) == 0 ||
		!plan.ObjectProperties[0].Group.Equal(state.ObjectProperties[0].Group) ||
		!plan.ObjectProperties[0].PortMonitoring.Equal(state.ObjectProperties[0].PortMonitoring) ||
		!plan.ObjectProperties[0].SortByName.Equal(state.ObjectProperties[0].SortByName) ||
		!plan.ObjectProperties[0].Label.Equal(state.ObjectProperties[0].Label) ||
		!plan.ObjectProperties[0].Icon.Equal(state.ObjectProperties[0].Icon)) {

		objProps := openapi.ConfigPutRequestEthPortProfileEthPortProfileNameObjectProperties{}
		hasObjPropsChanges := false

		if len(plan.ObjectProperties) > 0 {
			objProp := plan.ObjectProperties[0]

			if len(state.ObjectProperties) == 0 || !objProp.Group.Equal(state.ObjectProperties[0].Group) {
				hasObjPropsChanges = true
				if !objProp.Group.IsNull() {
					objProps.Group = openapi.PtrString(objProp.Group.ValueString())
				} else {
					objProps.Group = nil
				}
			}

			if len(state.ObjectProperties) == 0 || !objProp.PortMonitoring.Equal(state.ObjectProperties[0].PortMonitoring) {
				hasObjPropsChanges = true
				if !objProp.PortMonitoring.IsNull() {
					objProps.PortMonitoring = openapi.PtrString(objProp.PortMonitoring.ValueString())
				} else {
					objProps.PortMonitoring = nil
				}
			}

			if len(state.ObjectProperties) == 0 || !objProp.SortByName.Equal(state.ObjectProperties[0].SortByName) {
				hasObjPropsChanges = true
				if !objProp.SortByName.IsNull() {
					objProps.SortByName = openapi.PtrBool(objProp.SortByName.ValueBool())
				}
			}

			if len(state.ObjectProperties) == 0 || !objProp.Label.Equal(state.ObjectProperties[0].Label) {
				hasObjPropsChanges = true
				if !objProp.Label.IsNull() {
					objProps.Label = openapi.PtrString(objProp.Label.ValueString())
				}
			}

			if len(state.ObjectProperties) == 0 || !objProp.Icon.Equal(state.ObjectProperties[0].Icon) {
				hasObjPropsChanges = true
				if !objProp.Icon.IsNull() {
					objProps.Icon = openapi.PtrString(objProp.Icon.ValueString())
				}
			}
		}

		if hasObjPropsChanges {
			ethPortName.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	oldServicesByIndex := make(map[int64]servicesModel)
	for _, service := range state.Services {
		if !service.Index.IsNull() {
			oldServicesByIndex[service.Index.ValueInt64()] = service
		}
	}

	var changedServices []openapi.ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner
	servicesChanged := false

	for _, service := range plan.Services {
		if service.Index.IsNull() {
			continue
		}

		index := service.Index.ValueInt64()
		oldService, exists := oldServicesByIndex[index]

		if !exists {
			// new service, include all fields
			s := openapi.ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner{
				Index: openapi.PtrInt32(int32(index)),
			}

			if !service.RowNumEnable.IsNull() {
				s.RowNumEnable = openapi.PtrBool(service.RowNumEnable.ValueBool())
			} else {
				s.RowNumEnable = openapi.PtrBool(false)
			}

			hasService := !service.RowNumService.IsNull() && service.RowNumService.ValueString() != ""
			hasRefType := !service.RowNumServiceRefType.IsNull() && service.RowNumServiceRefType.ValueString() != ""

			if hasService || hasRefType {
				if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
					service.RowNumService, service.RowNumServiceRefType,
					"row_num_service", "row_num_service_ref_type_",
					hasService, hasRefType) {
					return
				}

				// Set both fields for new entries that have at least one of the fields
				if !service.RowNumService.IsNull() {
					s.RowNumService = openapi.PtrString(service.RowNumService.ValueString())
				} else {
					s.RowNumService = openapi.PtrString("")
				}

				if !service.RowNumServiceRefType.IsNull() {
					s.RowNumServiceRefType = openapi.PtrString(service.RowNumServiceRefType.ValueString())
				} else {
					s.RowNumServiceRefType = openapi.PtrString("")
				}
			} else {
				// If neither field is set, set both to empty strings
				s.RowNumService = openapi.PtrString("")
				s.RowNumServiceRefType = openapi.PtrString("")
			}

			if !service.RowNumExternalVlan.IsNull() {
				intVal := int32(service.RowNumExternalVlan.ValueInt64())
				s.RowNumExternalVlan = *openapi.NewNullableInt32(&intVal)
			} else {
				s.RowNumExternalVlan = *openapi.NewNullableInt32(nil)
			}

			if !service.RowNumMacFilter.IsNull() {
				s.RowNumMacFilter = openapi.PtrString(service.RowNumMacFilter.ValueString())
			} else {
				s.RowNumMacFilter = openapi.PtrString("")
			}

			if !service.RowNumMacFilterRefType.IsNull() {
				s.RowNumMacFilterRefType = openapi.PtrString(service.RowNumMacFilterRefType.ValueString())
			} else {
				s.RowNumMacFilterRefType = openapi.PtrString("")
			}

			if !service.RowNumLanIptv.IsNull() {
				s.RowNumLanIptv = openapi.PtrString(service.RowNumLanIptv.ValueString())
			} else {
				s.RowNumLanIptv = openapi.PtrString("")
			}

			changedServices = append(changedServices, s)
			servicesChanged = true
			continue
		}

		// existing service, check which fields changed
		s := openapi.ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner{
			Index: openapi.PtrInt32(int32(index)),
		}

		fieldChanged := false

		if !service.RowNumEnable.Equal(oldService.RowNumEnable) {
			s.RowNumEnable = openapi.PtrBool(service.RowNumEnable.ValueBool())
			fieldChanged = true
		}

		rowNumServiceChanged := !service.RowNumService.Equal(oldService.RowNumService)
		rowNumServiceRefTypeChanged := !service.RowNumServiceRefType.Equal(oldService.RowNumServiceRefType)

		if rowNumServiceChanged || rowNumServiceRefTypeChanged {
			// Validate using one ref type supported rules
			if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
				service.RowNumService, service.RowNumServiceRefType,
				"row_num_service", "row_num_service_ref_type_",
				rowNumServiceChanged, rowNumServiceRefTypeChanged) {
				return
			}

			// For fields with one reference type:
			// If only base field changes, send only base field
			// If ref type field changes (or both), send both fields
			if rowNumServiceChanged {
				if !service.RowNumService.IsNull() {
					s.RowNumService = openapi.PtrString(service.RowNumService.ValueString())
				} else {
					s.RowNumService = openapi.PtrString("")
				}
			}

			if rowNumServiceRefTypeChanged {
				if !service.RowNumServiceRefType.IsNull() {
					s.RowNumServiceRefType = openapi.PtrString(service.RowNumServiceRefType.ValueString())
				} else {
					s.RowNumServiceRefType = openapi.PtrString("")
				}

				if !rowNumServiceChanged {
					if !service.RowNumService.IsNull() {
						s.RowNumService = openapi.PtrString(service.RowNumService.ValueString())
					} else {
						s.RowNumService = openapi.PtrString("")
					}
				}
			}
			fieldChanged = true
		}

		if !service.RowNumExternalVlan.Equal(oldService.RowNumExternalVlan) {
			if !service.RowNumExternalVlan.IsNull() {
				intVal := int32(service.RowNumExternalVlan.ValueInt64())
				s.RowNumExternalVlan = *openapi.NewNullableInt32(&intVal)
			} else {
				s.RowNumExternalVlan = *openapi.NewNullableInt32(nil)
			}
			fieldChanged = true
		}

		if !service.RowNumMacFilter.Equal(oldService.RowNumMacFilter) {
			if !service.RowNumMacFilter.IsNull() {
				s.RowNumMacFilter = openapi.PtrString(service.RowNumMacFilter.ValueString())
			} else {
				s.RowNumMacFilter = openapi.PtrString("")
			}
			fieldChanged = true
		}

		if !service.RowNumMacFilterRefType.Equal(oldService.RowNumMacFilterRefType) {
			if !service.RowNumMacFilterRefType.IsNull() {
				s.RowNumMacFilterRefType = openapi.PtrString(service.RowNumMacFilterRefType.ValueString())
			} else {
				s.RowNumMacFilterRefType = openapi.PtrString("")
			}
			fieldChanged = true
		}

		if !service.RowNumLanIptv.Equal(oldService.RowNumLanIptv) {
			if !service.RowNumLanIptv.IsNull() {
				s.RowNumLanIptv = openapi.PtrString(service.RowNumLanIptv.ValueString())
			} else {
				s.RowNumLanIptv = openapi.PtrString("")
			}
			fieldChanged = true
		}

		if fieldChanged {
			changedServices = append(changedServices, s)
			servicesChanged = true
		}
	}

	for idx := range oldServicesByIndex {
		found := false
		for _, service := range plan.Services {
			if !service.Index.IsNull() && service.Index.ValueInt64() == idx {
				found = true
				break
			}
		}

		if !found {
			// service removed - include only the index for deletion
			deletedService := openapi.ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner{
				Index: openapi.PtrInt32(int32(idx)),
			}
			changedServices = append(changedServices, deletedService)
			servicesChanged = true
		}
	}

	if servicesChanged && len(changedServices) > 0 {
		ethPortName.Services = changedServices
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	operationID := r.bulkOpsMgr.AddPatch(ctx, "ethportprofile", name, ethPortName)
	r.notifyOperationAdded()
	tflog.Debug(ctx, fmt.Sprintf("Waiting for eth port profile update operation %s to complete", operationID))

	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update Eth Port Profile %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Eth port profile %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "eth_port_profiles")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *verityEthPortProfileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityEthPortProfileResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := ensureAuthenticated(ctx, r.provCtx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Authenticate",
			fmt.Sprintf("Error authenticating with API: %s", err),
		)
		return
	}

	name := state.Name.ValueString()

	operationID := r.bulkOpsMgr.AddDelete(ctx, "ethportprofile", name)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for eth port profile deletion operation %s to complete", operationID))

	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Delete Eth Port Profile %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Eth port profile %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "eth_port_profiles")
	resp.State.RemoveResource(ctx)
}

func (r *verityEthPortProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
