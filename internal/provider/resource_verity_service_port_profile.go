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

func (m verityServicePortProfileServiceModel) GetIndex() types.Int64 {
	return m.Index
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
	sppProps := &openapi.ServiceportprofilesPutRequestServicePortProfileValue{
		Name: openapi.PtrString(name),
	}

	// Handle string fields
	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "PortType", APIField: &sppProps.PortType, TFValue: plan.PortType},
		{FieldName: "TlsService", APIField: &sppProps.TlsService, TFValue: plan.TlsService},
		{FieldName: "TlsServiceRefType", APIField: &sppProps.TlsServiceRefType, TFValue: plan.TlsServiceRefType},
		{FieldName: "IpMask", APIField: &sppProps.IpMask, TFValue: plan.IpMask},
	})

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &sppProps.Enable, TFValue: plan.Enable},
		{FieldName: "TrustedPort", APIField: &sppProps.TrustedPort, TFValue: plan.TrustedPort},
	})

	// Handle nullable int64 fields
	utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
		{FieldName: "TlsLimitIn", APIField: &sppProps.TlsLimitIn, TFValue: plan.TlsLimitIn},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.ServiceportprofilesPutRequestServicePortProfileValueObjectProperties{}
		if !op.OnSummary.IsNull() {
			objProps.OnSummary = openapi.PtrBool(op.OnSummary.ValueBool())
		} else {
			objProps.OnSummary = nil
		}
		if !op.PortMonitoring.IsNull() {
			objProps.PortMonitoring = openapi.PtrString(op.PortMonitoring.ValueString())
		} else {
			objProps.PortMonitoring = nil
		}
		if !op.Group.IsNull() {
			objProps.Group = openapi.PtrString(op.Group.ValueString())
		} else {
			objProps.Group = nil
		}
		sppProps.ObjectProperties = &objProps
	}

	// Handle services
	if len(plan.Services) > 0 {
		services := make([]openapi.ServiceportprofilesPutRequestServicePortProfileValueServicesInner, len(plan.Services))
		for i, service := range plan.Services {
			serviceItem := openapi.ServiceportprofilesPutRequestServicePortProfileValueServicesInner{}

			if !service.RowNumService.IsNull() {
				serviceItem.RowNumService = openapi.PtrString(service.RowNumService.ValueString())
			}
			if !service.RowNumServiceRefType.IsNull() {
				serviceItem.RowNumServiceRefType = openapi.PtrString(service.RowNumServiceRefType.ValueString())
			}
			if !service.RowNumEnable.IsNull() {
				serviceItem.RowNumEnable = openapi.PtrBool(service.RowNumEnable.ValueBool())
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

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "service_port_profile", name, *sppProps, &resp.Diagnostics)
	if !success {
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

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("service_port_profile") {
		tflog.Info(ctx, fmt.Sprintf("Skipping Service Port Profile %s verification – trusting recent successful API operation", sppName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching Service Port Profiles for verification of %s", sppName))

	type ServicePortProfileResponse struct {
		ServicePortProfile map[string]interface{} `json:"service_port_profile"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "service_port_profiles", sppName,
		func() (ServicePortProfileResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch Service Port Profiles")
			respAPI, err := r.client.ServicePortProfilesAPI.ServiceportprofilesGet(ctx).Execute()
			if err != nil {
				return ServicePortProfileResponse{}, fmt.Errorf("error reading Service Port Profiles: %v", err)
			}
			defer respAPI.Body.Close()

			var res ServicePortProfileResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return ServicePortProfileResponse{}, fmt.Errorf("failed to decode Service Port Profiles response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d Service Port Profiles", len(res.ServicePortProfile)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Service Port Profile %s", sppName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for Service Port Profile with name: %s", sppName))

	sppData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.ServicePortProfile,
		sppName,
		func(data interface{}) (string, bool) {
			if spp, ok := data.(map[string]interface{}); ok {
				if name, ok := spp["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Service Port Profile with name '%s' not found in API response", sppName))
		resp.State.RemoveResource(ctx)
		return
	}

	sppMap, ok := sppData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Service Port Profile Data",
			fmt.Sprintf("Service Port Profile data is not in expected format for %s", sppName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found Service Port Profile '%s' under API key '%s'", sppName, actualAPIName))

	state.Name = utils.MapStringFromAPI(sppMap["name"])

	// Handle object properties
	if objProps, ok := sppMap["object_properties"].(map[string]interface{}); ok {
		onSummary := utils.MapBoolFromAPI(objProps["on_summary"])
		if onSummary.IsNull() {
			onSummary = types.BoolValue(false)
		}
		portMonitoring := utils.MapStringFromAPI(objProps["port_monitoring"])
		if portMonitoring.IsNull() {
			portMonitoring = types.StringValue("")
		}
		group := utils.MapStringFromAPI(objProps["group"])
		if group.IsNull() {
			group = types.StringValue("")
		}
		state.ObjectProperties = []verityServicePortProfileObjectPropertiesModel{
			{
				OnSummary:      onSummary,
				PortMonitoring: portMonitoring,
				Group:          group,
			},
		}
	} else {
		state.ObjectProperties = nil
	}

	// Map string fields
	stringFieldMappings := map[string]*types.String{
		"port_type":             &state.PortType,
		"tls_service":           &state.TlsService,
		"tls_service_ref_type_": &state.TlsServiceRefType,
		"ip_mask":               &state.IpMask,
	}

	for apiKey, stateField := range stringFieldMappings {
		*stateField = utils.MapStringFromAPI(sppMap[apiKey])
	}

	// Map boolean fields
	boolFieldMappings := map[string]*types.Bool{
		"enable":       &state.Enable,
		"trusted_port": &state.TrustedPort,
	}

	for apiKey, stateField := range boolFieldMappings {
		*stateField = utils.MapBoolFromAPI(sppMap[apiKey])
	}

	// Map int64 fields
	int64FieldMappings := map[string]*types.Int64{
		"tls_limit_in": &state.TlsLimitIn,
	}

	for apiKey, stateField := range int64FieldMappings {
		*stateField = utils.MapInt64FromAPI(sppMap[apiKey])
	}

	// Handle services
	if servicesArray, ok := sppMap["services"].([]interface{}); ok && len(servicesArray) > 0 {
		var services []verityServicePortProfileServiceModel
		for _, s := range servicesArray {
			service, ok := s.(map[string]interface{})
			if !ok {
				continue
			}

			srModel := verityServicePortProfileServiceModel{
				RowNumEnable:         utils.MapBoolFromAPI(service["row_num_enable"]),
				RowNumService:        utils.MapStringFromAPI(service["row_num_service"]),
				RowNumServiceRefType: utils.MapStringFromAPI(service["row_num_service_ref_type_"]),
				RowNumExternalVlan:   utils.MapInt64FromAPI(service["row_num_external_vlan"]),
				RowNumLimitIn:        utils.MapInt64FromAPI(service["row_num_limit_in"]),
				RowNumLimitOut:       utils.MapInt64FromAPI(service["row_num_limit_out"]),
				Index:                utils.MapInt64FromAPI(service["index"]),
			}

			services = append(services, srModel)
		}
		state.Services = services
	} else {
		state.Services = nil
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
	sppProps := openapi.ServiceportprofilesPutRequestServicePortProfileValue{}
	hasChanges := false

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { sppProps.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.PortType, state.PortType, func(v *string) { sppProps.PortType = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.TlsService, state.TlsService, func(v *string) { sppProps.TlsService = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.TlsServiceRefType, state.TlsServiceRefType, func(v *string) { sppProps.TlsServiceRefType = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.IpMask, state.IpMask, func(v *string) { sppProps.IpMask = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { sppProps.Enable = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.TrustedPort, state.TrustedPort, func(v *bool) { sppProps.TrustedPort = v }, &hasChanges)

	// Handle nullable int64 field changes
	utils.CompareAndSetNullableInt64Field(plan.TlsLimitIn, state.TlsLimitIn, func(v *openapi.NullableInt32) { sppProps.TlsLimitIn = *v }, &hasChanges)

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		if len(state.ObjectProperties) == 0 ||
			!plan.ObjectProperties[0].OnSummary.Equal(state.ObjectProperties[0].OnSummary) ||
			!plan.ObjectProperties[0].PortMonitoring.Equal(state.ObjectProperties[0].PortMonitoring) ||
			!plan.ObjectProperties[0].Group.Equal(state.ObjectProperties[0].Group) {
			op := plan.ObjectProperties[0]
			objProps := openapi.ServiceportprofilesPutRequestServicePortProfileValueObjectProperties{}
			if !op.OnSummary.IsNull() {
				objProps.OnSummary = openapi.PtrBool(op.OnSummary.ValueBool())
			} else {
				objProps.OnSummary = nil
			}
			if !op.PortMonitoring.IsNull() {
				objProps.PortMonitoring = openapi.PtrString(op.PortMonitoring.ValueString())
			} else {
				objProps.PortMonitoring = nil
			}
			if !op.Group.IsNull() {
				objProps.Group = openapi.PtrString(op.Group.ValueString())
			} else {
				objProps.Group = nil
			}
			sppProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	// Handle TlsService and TlsServiceRefType using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.TlsService, state.TlsService, plan.TlsServiceRefType, state.TlsServiceRefType,
		func(v *string) { sppProps.TlsService = v },
		func(v *string) { sppProps.TlsServiceRefType = v },
		"tls_service", "tls_service_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	// Handle services
	servicesHandler := utils.IndexedItemHandler[verityServicePortProfileServiceModel, openapi.ServiceportprofilesPutRequestServicePortProfileValueServicesInner]{
		CreateNew: func(planItem verityServicePortProfileServiceModel) openapi.ServiceportprofilesPutRequestServicePortProfileValueServicesInner {
			service := openapi.ServiceportprofilesPutRequestServicePortProfileValueServicesInner{
				Index: openapi.PtrInt32(int32(planItem.Index.ValueInt64())),
			}

			if !planItem.RowNumEnable.IsNull() {
				service.RowNumEnable = openapi.PtrBool(planItem.RowNumEnable.ValueBool())
			} else {
				service.RowNumEnable = openapi.PtrBool(false)
			}

			if !planItem.RowNumService.IsNull() {
				service.RowNumService = openapi.PtrString(planItem.RowNumService.ValueString())
			} else {
				service.RowNumService = openapi.PtrString("")
			}

			if !planItem.RowNumServiceRefType.IsNull() {
				service.RowNumServiceRefType = openapi.PtrString(planItem.RowNumServiceRefType.ValueString())
			} else {
				service.RowNumServiceRefType = openapi.PtrString("")
			}

			if !planItem.RowNumExternalVlan.IsNull() {
				val := int32(planItem.RowNumExternalVlan.ValueInt64())
				service.RowNumExternalVlan = *openapi.NewNullableInt32(&val)
			} else {
				service.RowNumExternalVlan = *openapi.NewNullableInt32(nil)
			}

			if !planItem.RowNumLimitIn.IsNull() {
				val := int32(planItem.RowNumLimitIn.ValueInt64())
				service.RowNumLimitIn = *openapi.NewNullableInt32(&val)
			} else {
				service.RowNumLimitIn = *openapi.NewNullableInt32(nil)
			}

			if !planItem.RowNumLimitOut.IsNull() {
				val := int32(planItem.RowNumLimitOut.ValueInt64())
				service.RowNumLimitOut = *openapi.NewNullableInt32(&val)
			} else {
				service.RowNumLimitOut = *openapi.NewNullableInt32(nil)
			}

			return service
		},
		UpdateExisting: func(planItem verityServicePortProfileServiceModel, stateItem verityServicePortProfileServiceModel) (openapi.ServiceportprofilesPutRequestServicePortProfileValueServicesInner, bool) {
			service := openapi.ServiceportprofilesPutRequestServicePortProfileValueServicesInner{
				Index: openapi.PtrInt32(int32(planItem.Index.ValueInt64())),
			}

			fieldChanged := false

			if !planItem.RowNumEnable.Equal(stateItem.RowNumEnable) {
				service.RowNumEnable = openapi.PtrBool(planItem.RowNumEnable.ValueBool())
				fieldChanged = true
			}

			if !planItem.RowNumService.Equal(stateItem.RowNumService) {
				if !planItem.RowNumService.IsNull() {
					service.RowNumService = openapi.PtrString(planItem.RowNumService.ValueString())
				} else {
					service.RowNumService = openapi.PtrString("")
				}
				fieldChanged = true
			}

			if !planItem.RowNumServiceRefType.Equal(stateItem.RowNumServiceRefType) {
				if !planItem.RowNumServiceRefType.IsNull() {
					service.RowNumServiceRefType = openapi.PtrString(planItem.RowNumServiceRefType.ValueString())
				} else {
					service.RowNumServiceRefType = openapi.PtrString("")
				}
				fieldChanged = true
			}

			if !planItem.RowNumExternalVlan.Equal(stateItem.RowNumExternalVlan) {
				if !planItem.RowNumExternalVlan.IsNull() {
					val := int32(planItem.RowNumExternalVlan.ValueInt64())
					service.RowNumExternalVlan = *openapi.NewNullableInt32(&val)
				} else {
					service.RowNumExternalVlan = *openapi.NewNullableInt32(nil)
				}
				fieldChanged = true
			}

			if !planItem.RowNumLimitIn.Equal(stateItem.RowNumLimitIn) {
				if !planItem.RowNumLimitIn.IsNull() {
					val := int32(planItem.RowNumLimitIn.ValueInt64())
					service.RowNumLimitIn = *openapi.NewNullableInt32(&val)
				} else {
					service.RowNumLimitIn = *openapi.NewNullableInt32(nil)
				}
				fieldChanged = true
			}

			if !planItem.RowNumLimitOut.Equal(stateItem.RowNumLimitOut) {
				if !planItem.RowNumLimitOut.IsNull() {
					val := int32(planItem.RowNumLimitOut.ValueInt64())
					service.RowNumLimitOut = *openapi.NewNullableInt32(&val)
				} else {
					service.RowNumLimitOut = *openapi.NewNullableInt32(nil)
				}
				fieldChanged = true
			}

			return service, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.ServiceportprofilesPutRequestServicePortProfileValueServicesInner {
			return openapi.ServiceportprofilesPutRequestServicePortProfileValueServicesInner{
				Index: openapi.PtrInt32(int32(index)),
			}
		},
	}

	changedServices, servicesChanged := utils.ProcessIndexedArrayUpdates(plan.Services, state.Services, servicesHandler)
	if servicesChanged {
		sppProps.Services = changedServices
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "service_port_profile", name, sppProps, &resp.Diagnostics)
	if !success {
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

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "service_port_profile", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Service Port Profile %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "service_port_profiles")
	resp.State.RemoveResource(ctx)
}

func (r *verityServicePortProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
