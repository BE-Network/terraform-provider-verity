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
	Name              types.String                                `tfsdk:"name"`
	Enable            types.Bool                                  `tfsdk:"enable"`
	IngressAcl        types.String                                `tfsdk:"ingress_acl"`
	IngressAclRefType types.String                                `tfsdk:"ingress_acl_ref_type_"`
	EgressAcl         types.String                                `tfsdk:"egress_acl"`
	EgressAclRefType  types.String                                `tfsdk:"egress_acl_ref_type_"`
	ObjectProperties  []verityEthPortProfileObjectPropertiesModel `tfsdk:"object_properties"`
	Services          []servicesModel                             `tfsdk:"services"`
	Tls               types.Bool                                  `tfsdk:"tls"`
	TlsService        types.String                                `tfsdk:"tls_service"`
	TlsServiceRefType types.String                                `tfsdk:"tls_service_ref_type_"`
	TrustedPort       types.Bool                                  `tfsdk:"trusted_port"`
}

type verityEthPortProfileObjectPropertiesModel struct {
	Group          types.String `tfsdk:"group"`
	PortMonitoring types.String `tfsdk:"port_monitoring"`
	SortByName     types.Bool   `tfsdk:"sort_by_name"`
	Label          types.String `tfsdk:"label"`
	Icon           types.String `tfsdk:"icon"`
}

type servicesModel struct {
	RowNumEnable            types.Bool   `tfsdk:"row_num_enable"`
	RowNumService           types.String `tfsdk:"row_num_service"`
	RowNumServiceRefType    types.String `tfsdk:"row_num_service_ref_type_"`
	RowNumExternalVlan      types.Int64  `tfsdk:"row_num_external_vlan"`
	RowNumIngressAcl        types.String `tfsdk:"row_num_ingress_acl"`
	RowNumIngressAclRefType types.String `tfsdk:"row_num_ingress_acl_ref_type_"`
	RowNumEgressAcl         types.String `tfsdk:"row_num_egress_acl"`
	RowNumEgressAclRefType  types.String `tfsdk:"row_num_egress_acl_ref_type_"`
	Index                   types.Int64  `tfsdk:"index"`
	RowNumMacFilter         types.String `tfsdk:"row_num_mac_filter"`
	RowNumMacFilterRefType  types.String `tfsdk:"row_num_mac_filter_ref_type_"`
	RowNumLanIptv           types.String `tfsdk:"row_num_lan_iptv"`
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
			},
			"ingress_acl": schema.StringAttribute{
				Description: "Choose an ingress access control list",
				Optional:    true,
			},
			"ingress_acl_ref_type_": schema.StringAttribute{
				Description: "Object type for ingress_acl field",
				Optional:    true,
			},
			"egress_acl": schema.StringAttribute{
				Description: "Choose an egress access control list",
				Optional:    true,
			},
			"egress_acl_ref_type_": schema.StringAttribute{
				Description: "Object type for egress_acl field",
				Optional:    true,
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
						"row_num_ingress_acl": schema.StringAttribute{
							Description: "Choose an ingress access control list",
							Optional:    true,
						},
						"row_num_ingress_acl_ref_type_": schema.StringAttribute{
							Description: "Object type for row_num_ingress_acl field",
							Optional:    true,
						},
						"row_num_egress_acl": schema.StringAttribute{
							Description: "Choose an egress access control list",
							Optional:    true,
						},
						"row_num_egress_acl_ref_type_": schema.StringAttribute{
							Description: "Object type for row_num_egress_acl field",
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
	ethPortName := &openapi.EthportprofilesPutRequestEthPortProfileValue{
		Name: openapi.PtrString(name),
	}

	if !plan.Enable.IsNull() {
		ethPortName.Enable = openapi.PtrBool(plan.Enable.ValueBool())
	}

	if !plan.IngressAcl.IsNull() {
		ethPortName.IngressAcl = openapi.PtrString(plan.IngressAcl.ValueString())
	}

	if !plan.IngressAclRefType.IsNull() {
		ethPortName.IngressAclRefType = openapi.PtrString(plan.IngressAclRefType.ValueString())
	}

	if !plan.EgressAcl.IsNull() {
		ethPortName.EgressAcl = openapi.PtrString(plan.EgressAcl.ValueString())
	}

	if !plan.EgressAclRefType.IsNull() {
		ethPortName.EgressAclRefType = openapi.PtrString(plan.EgressAclRefType.ValueString())
	}

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
		objProps := openapi.EthportprofilesPutRequestEthPortProfileValueObjectProperties{}
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
		services := make([]openapi.EthportprofilesPutRequestEthPortProfileValueServicesInner, len(plan.Services))

		for i, service := range plan.Services {
			s := openapi.EthportprofilesPutRequestEthPortProfileValueServicesInner{}

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

			if !service.RowNumIngressAcl.IsNull() {
				s.RowNumIngressAcl = openapi.PtrString(service.RowNumIngressAcl.ValueString())
			}

			if !service.RowNumIngressAclRefType.IsNull() {
				s.RowNumIngressAclRefType = openapi.PtrString(service.RowNumIngressAclRefType.ValueString())
			}

			if !service.RowNumEgressAcl.IsNull() {
				s.RowNumEgressAcl = openapi.PtrString(service.RowNumEgressAcl.ValueString())
			}

			if !service.RowNumEgressAclRefType.IsNull() {
				s.RowNumEgressAclRefType = openapi.PtrString(service.RowNumEgressAclRefType.ValueString())
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

	operationID := r.bulkOpsMgr.AddPut(ctx, "eth_port_profile", name, *ethPortName)

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

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("eth_port_profile") {
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

	if ingressAcl, ok := profile["ingress_acl"].(string); ok {
		state.IngressAcl = types.StringValue(ingressAcl)
	} else {
		state.IngressAcl = types.StringNull()
	}

	if ingressAclRefType, ok := profile["ingress_acl_ref_type_"].(string); ok {
		state.IngressAclRefType = types.StringValue(ingressAclRefType)
	} else {
		state.IngressAclRefType = types.StringNull()
	}

	if egressAcl, ok := profile["egress_acl"].(string); ok {
		state.EgressAcl = types.StringValue(egressAcl)
	} else {
		state.EgressAcl = types.StringNull()
	}

	if egressAclRefType, ok := profile["egress_acl_ref_type_"].(string); ok {
		state.EgressAclRefType = types.StringValue(egressAclRefType)
	} else {
		state.EgressAclRefType = types.StringNull()
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

				if ingressAcl, ok := serviceMap["row_num_ingress_acl"].(string); ok {
					serviceModel.RowNumIngressAcl = types.StringValue(ingressAcl)
				} else {
					serviceModel.RowNumIngressAcl = types.StringNull()
				}

				if ingressAclRefType, ok := serviceMap["row_num_ingress_acl_ref_type_"].(string); ok {
					serviceModel.RowNumIngressAclRefType = types.StringValue(ingressAclRefType)
				} else {
					serviceModel.RowNumIngressAclRefType = types.StringNull()
				}

				if egressAcl, ok := serviceMap["row_num_egress_acl"].(string); ok {
					serviceModel.RowNumEgressAcl = types.StringValue(egressAcl)
				} else {
					serviceModel.RowNumEgressAcl = types.StringNull()
				}

				if egressAclRefType, ok := serviceMap["row_num_egress_acl_ref_type_"].(string); ok {
					serviceModel.RowNumEgressAclRefType = types.StringValue(egressAclRefType)
				} else {
					serviceModel.RowNumEgressAclRefType = types.StringNull()
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
	ethPortName := openapi.EthportprofilesPutRequestEthPortProfileValue{}

	hasChanges := false

	if !plan.Enable.Equal(state.Enable) {
		ethPortName.Enable = openapi.PtrBool(plan.Enable.ValueBool())
		hasChanges = true
	}

	ingressAclChanged := !plan.IngressAcl.Equal(state.IngressAcl)
	ingressAclRefTypeChanged := !plan.IngressAclRefType.Equal(state.IngressAclRefType)

	if ingressAclChanged || ingressAclRefTypeChanged {
		if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
			plan.IngressAcl, plan.IngressAclRefType,
			"ingress_acl", "ingress_acl_ref_type_",
			ingressAclChanged, ingressAclRefTypeChanged) {
			return
		}

		// For fields with one reference type:
		// If only base field changes, send only base field
		// If ref type field changes (or both), send both fields
		if ingressAclChanged && !ingressAclRefTypeChanged {
			// Just send the base field
			if !plan.IngressAcl.IsNull() && plan.IngressAcl.ValueString() != "" {
				ethPortName.IngressAcl = openapi.PtrString(plan.IngressAcl.ValueString())
			} else {
				ethPortName.IngressAcl = openapi.PtrString("")
			}
			hasChanges = true
		} else if ingressAclRefTypeChanged {
			// Send both fields
			if !plan.IngressAcl.IsNull() && plan.IngressAcl.ValueString() != "" {
				ethPortName.IngressAcl = openapi.PtrString(plan.IngressAcl.ValueString())
			} else {
				ethPortName.IngressAcl = openapi.PtrString("")
			}

			if !plan.IngressAclRefType.IsNull() && plan.IngressAclRefType.ValueString() != "" {
				ethPortName.IngressAclRefType = openapi.PtrString(plan.IngressAclRefType.ValueString())
			} else {
				ethPortName.IngressAclRefType = openapi.PtrString("")
			}
			hasChanges = true
		}
	}

	egressAclChanged := !plan.EgressAcl.Equal(state.EgressAcl)
	egressAclRefTypeChanged := !plan.EgressAclRefType.Equal(state.EgressAclRefType)

	if egressAclChanged || egressAclRefTypeChanged {
		if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
			plan.EgressAcl, plan.EgressAclRefType,
			"egress_acl", "egress_acl_ref_type_",
			egressAclChanged, egressAclRefTypeChanged) {
			return
		}

		// For fields with one reference type:
		// If only base field changes, send only base field
		// If ref type field changes (or both), send both fields
		if egressAclChanged && !egressAclRefTypeChanged {
			// Just send the base field
			if !plan.EgressAcl.IsNull() && plan.EgressAcl.ValueString() != "" {
				ethPortName.EgressAcl = openapi.PtrString(plan.EgressAcl.ValueString())
			} else {
				ethPortName.EgressAcl = openapi.PtrString("")
			}
			hasChanges = true
		} else if egressAclRefTypeChanged {
			// Send both fields
			if !plan.EgressAcl.IsNull() && plan.EgressAcl.ValueString() != "" {
				ethPortName.EgressAcl = openapi.PtrString(plan.EgressAcl.ValueString())
			} else {
				ethPortName.EgressAcl = openapi.PtrString("")
			}

			if !plan.EgressAclRefType.IsNull() && plan.EgressAclRefType.ValueString() != "" {
				ethPortName.EgressAclRefType = openapi.PtrString(plan.EgressAclRefType.ValueString())
			} else {
				ethPortName.EgressAclRefType = openapi.PtrString("")
			}
			hasChanges = true
		}
	}

	if !plan.Tls.Equal(state.Tls) {
		ethPortName.Tls = openapi.PtrBool(plan.Tls.ValueBool())
		hasChanges = true
	}

	tlsServiceChanged := !plan.TlsService.Equal(state.TlsService)
	tlsServiceRefTypeChanged := !plan.TlsServiceRefType.Equal(state.TlsServiceRefType)

	if tlsServiceChanged || tlsServiceRefTypeChanged {
		if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
			plan.TlsService, plan.TlsServiceRefType,
			"tls_service", "tls_service_ref_type_",
			tlsServiceChanged, tlsServiceRefTypeChanged) {
			return
		}

		// For fields with one reference type:
		// If only base field changes, send only base field
		// If ref type field changes (or both), send both fields
		if tlsServiceChanged && !tlsServiceRefTypeChanged {
			// Just send the base field
			if !plan.TlsService.IsNull() && plan.TlsService.ValueString() != "" {
				ethPortName.TlsService = openapi.PtrString(plan.TlsService.ValueString())
			} else {
				ethPortName.TlsService = openapi.PtrString("")
			}
			hasChanges = true
		} else if tlsServiceRefTypeChanged {
			// Send both fields
			if !plan.TlsService.IsNull() && plan.TlsService.ValueString() != "" {
				ethPortName.TlsService = openapi.PtrString(plan.TlsService.ValueString())
			} else {
				ethPortName.TlsService = openapi.PtrString("")
			}

			if !plan.TlsServiceRefType.IsNull() && plan.TlsServiceRefType.ValueString() != "" {
				ethPortName.TlsServiceRefType = openapi.PtrString(plan.TlsServiceRefType.ValueString())
			} else {
				ethPortName.TlsServiceRefType = openapi.PtrString("")
			}
			hasChanges = true
		}
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

		objProps := openapi.EthportprofilesPutRequestEthPortProfileValueObjectProperties{}
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

	var changedServices []openapi.EthportprofilesPutRequestEthPortProfileValueServicesInner
	servicesChanged := false

	for _, service := range plan.Services {
		if service.Index.IsNull() {
			continue
		}

		index := service.Index.ValueInt64()
		oldService, exists := oldServicesByIndex[index]

		if !exists {
			// new service, include all fields
			s := openapi.EthportprofilesPutRequestEthPortProfileValueServicesInner{
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

			hasIngressAcl := !service.RowNumIngressAcl.IsNull() && service.RowNumIngressAcl.ValueString() != ""
			hasIngressAclRefType := !service.RowNumIngressAclRefType.IsNull() && service.RowNumIngressAclRefType.ValueString() != ""

			if hasIngressAcl || hasIngressAclRefType {
				if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
					service.RowNumIngressAcl, service.RowNumIngressAclRefType,
					"row_num_ingress_acl", "row_num_ingress_acl_ref_type_",
					hasIngressAcl, hasIngressAclRefType) {
					return
				}
			}

			if !service.RowNumIngressAcl.IsNull() {
				s.RowNumIngressAcl = openapi.PtrString(service.RowNumIngressAcl.ValueString())
			} else {
				s.RowNumIngressAcl = openapi.PtrString("")
			}

			if !service.RowNumIngressAclRefType.IsNull() {
				s.RowNumIngressAclRefType = openapi.PtrString(service.RowNumIngressAclRefType.ValueString())
			} else {
				s.RowNumIngressAclRefType = openapi.PtrString("")
			}

			hasEgressAcl := !service.RowNumEgressAcl.IsNull() && service.RowNumEgressAcl.ValueString() != ""
			hasEgressAclRefType := !service.RowNumEgressAclRefType.IsNull() && service.RowNumEgressAclRefType.ValueString() != ""

			if hasEgressAcl || hasEgressAclRefType {
				if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
					service.RowNumEgressAcl, service.RowNumEgressAclRefType,
					"row_num_egress_acl", "row_num_egress_acl_ref_type_",
					hasEgressAcl, hasEgressAclRefType) {
					return
				}
			}

			if !service.RowNumEgressAcl.IsNull() {
				s.RowNumEgressAcl = openapi.PtrString(service.RowNumEgressAcl.ValueString())
			} else {
				s.RowNumEgressAcl = openapi.PtrString("")
			}

			if !service.RowNumEgressAclRefType.IsNull() {
				s.RowNumEgressAclRefType = openapi.PtrString(service.RowNumEgressAclRefType.ValueString())
			} else {
				s.RowNumEgressAclRefType = openapi.PtrString("")
			}

			hasMacFilter := !service.RowNumMacFilter.IsNull() && service.RowNumMacFilter.ValueString() != ""
			hasMacFilterRefType := !service.RowNumMacFilterRefType.IsNull() && service.RowNumMacFilterRefType.ValueString() != ""

			if hasMacFilter || hasMacFilterRefType {
				if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
					service.RowNumMacFilter, service.RowNumMacFilterRefType,
					"row_num_mac_filter", "row_num_mac_filter_ref_type_",
					hasMacFilter, hasMacFilterRefType) {
					return
				}
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
		s := openapi.EthportprofilesPutRequestEthPortProfileValueServicesInner{
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

		rowNumIngressAclChanged := !service.RowNumIngressAcl.Equal(oldService.RowNumIngressAcl)
		rowNumIngressAclRefTypeChanged := !service.RowNumIngressAclRefType.Equal(oldService.RowNumIngressAclRefType)

		if rowNumIngressAclChanged || rowNumIngressAclRefTypeChanged {
			if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
				service.RowNumIngressAcl, service.RowNumIngressAclRefType,
				"row_num_ingress_acl", "row_num_ingress_acl_ref_type_",
				rowNumIngressAclChanged, rowNumIngressAclRefTypeChanged) {
				return
			}

			// For fields with one reference type:
			// If only base field changes, send only base field
			// If ref type field changes (or both), send both fields
			if rowNumIngressAclChanged && !rowNumIngressAclRefTypeChanged {
				// Just send the base field
				if !service.RowNumIngressAcl.IsNull() && service.RowNumIngressAcl.ValueString() != "" {
					s.RowNumIngressAcl = openapi.PtrString(service.RowNumIngressAcl.ValueString())
				} else {
					s.RowNumIngressAcl = openapi.PtrString("")
				}
				fieldChanged = true
			} else if rowNumIngressAclRefTypeChanged {
				// Send both fields
				if !service.RowNumIngressAcl.IsNull() && service.RowNumIngressAcl.ValueString() != "" {
					s.RowNumIngressAcl = openapi.PtrString(service.RowNumIngressAcl.ValueString())
				} else {
					s.RowNumIngressAcl = openapi.PtrString("")
				}

				if !service.RowNumIngressAclRefType.IsNull() && service.RowNumIngressAclRefType.ValueString() != "" {
					s.RowNumIngressAclRefType = openapi.PtrString(service.RowNumIngressAclRefType.ValueString())
				} else {
					s.RowNumIngressAclRefType = openapi.PtrString("")
				}
				fieldChanged = true
			}
		}

		rowNumEgressAclChanged := !service.RowNumEgressAcl.Equal(oldService.RowNumEgressAcl)
		rowNumEgressAclRefTypeChanged := !service.RowNumEgressAclRefType.Equal(oldService.RowNumEgressAclRefType)

		if rowNumEgressAclChanged || rowNumEgressAclRefTypeChanged {
			if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
				service.RowNumEgressAcl, service.RowNumEgressAclRefType,
				"row_num_egress_acl", "row_num_egress_acl_ref_type_",
				rowNumEgressAclChanged, rowNumEgressAclRefTypeChanged) {
				return
			}

			// For fields with one reference type:
			// If only base field changes, send only base field
			// If ref type field changes (or both), send both fields
			if rowNumEgressAclChanged && !rowNumEgressAclRefTypeChanged {
				// Just send the base field
				if !service.RowNumEgressAcl.IsNull() && service.RowNumEgressAcl.ValueString() != "" {
					s.RowNumEgressAcl = openapi.PtrString(service.RowNumEgressAcl.ValueString())
				} else {
					s.RowNumEgressAcl = openapi.PtrString("")
				}
				fieldChanged = true
			} else if rowNumEgressAclRefTypeChanged {
				// Send both fields
				if !service.RowNumEgressAcl.IsNull() && service.RowNumEgressAcl.ValueString() != "" {
					s.RowNumEgressAcl = openapi.PtrString(service.RowNumEgressAcl.ValueString())
				} else {
					s.RowNumEgressAcl = openapi.PtrString("")
				}

				if !service.RowNumEgressAclRefType.IsNull() && service.RowNumEgressAclRefType.ValueString() != "" {
					s.RowNumEgressAclRefType = openapi.PtrString(service.RowNumEgressAclRefType.ValueString())
				} else {
					s.RowNumEgressAclRefType = openapi.PtrString("")
				}
				fieldChanged = true
			}
		}

		rowNumMacFilterChanged := !service.RowNumMacFilter.Equal(oldService.RowNumMacFilter)
		rowNumMacFilterRefTypeChanged := !service.RowNumMacFilterRefType.Equal(oldService.RowNumMacFilterRefType)

		if rowNumMacFilterChanged || rowNumMacFilterRefTypeChanged {
			if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
				service.RowNumMacFilter, service.RowNumMacFilterRefType,
				"row_num_mac_filter", "row_num_mac_filter_ref_type_",
				rowNumMacFilterChanged, rowNumMacFilterRefTypeChanged) {
				return
			}

			// For fields with one reference type:
			// If only base field changes, send only base field
			// If ref type field changes (or both), send both fields
			if rowNumMacFilterChanged && !rowNumMacFilterRefTypeChanged {
				// Just send the base field
				if !service.RowNumMacFilter.IsNull() && service.RowNumMacFilter.ValueString() != "" {
					s.RowNumMacFilter = openapi.PtrString(service.RowNumMacFilter.ValueString())
				} else {
					s.RowNumMacFilter = openapi.PtrString("")
				}
				fieldChanged = true
			} else if rowNumMacFilterRefTypeChanged {
				// Send both fields
				if !service.RowNumMacFilter.IsNull() && service.RowNumMacFilter.ValueString() != "" {
					s.RowNumMacFilter = openapi.PtrString(service.RowNumMacFilter.ValueString())
				} else {
					s.RowNumMacFilter = openapi.PtrString("")
				}

				if !service.RowNumMacFilterRefType.IsNull() && service.RowNumMacFilterRefType.ValueString() != "" {
					s.RowNumMacFilterRefType = openapi.PtrString(service.RowNumMacFilterRefType.ValueString())
				} else {
					s.RowNumMacFilterRefType = openapi.PtrString("")
				}
				fieldChanged = true
			}
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
			deletedService := openapi.EthportprofilesPutRequestEthPortProfileValueServicesInner{
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

	operationID := r.bulkOpsMgr.AddPatch(ctx, "eth_port_profile", name, ethPortName)
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

	operationID := r.bulkOpsMgr.AddDelete(ctx, "eth_port_profile", name)
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
