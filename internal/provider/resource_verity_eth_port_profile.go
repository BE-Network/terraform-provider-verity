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

func (s servicesModel) GetIndex() types.Int64 {
	return s.Index
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

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Authenticate",
			fmt.Sprintf("Error authenticating with API: %s", err),
		)
		return
	}

	name := plan.Name.ValueString()
	ethPortProfileProps := &openapi.EthportprofilesPutRequestEthPortProfileValue{
		Name: openapi.PtrString(name),
	}

	// Handle string fields
	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "IngressAcl", APIField: &ethPortProfileProps.IngressAcl, TFValue: plan.IngressAcl},
		{FieldName: "IngressAclRefType", APIField: &ethPortProfileProps.IngressAclRefType, TFValue: plan.IngressAclRefType},
		{FieldName: "EgressAcl", APIField: &ethPortProfileProps.EgressAcl, TFValue: plan.EgressAcl},
		{FieldName: "EgressAclRefType", APIField: &ethPortProfileProps.EgressAclRefType, TFValue: plan.EgressAclRefType},
		{FieldName: "TlsService", APIField: &ethPortProfileProps.TlsService, TFValue: plan.TlsService},
		{FieldName: "TlsServiceRefType", APIField: &ethPortProfileProps.TlsServiceRefType, TFValue: plan.TlsServiceRefType},
	})

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &ethPortProfileProps.Enable, TFValue: plan.Enable},
		{FieldName: "Tls", APIField: &ethPortProfileProps.Tls, TFValue: plan.Tls},
		{FieldName: "TrustedPort", APIField: &ethPortProfileProps.TrustedPort, TFValue: plan.TrustedPort},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.EthportprofilesPutRequestEthPortProfileValueObjectProperties{}
		if !op.Group.IsNull() {
			objProps.Group = openapi.PtrString(op.Group.ValueString())
		} else {
			objProps.Group = nil
		}
		if !op.PortMonitoring.IsNull() {
			objProps.PortMonitoring = openapi.PtrString(op.PortMonitoring.ValueString())
		} else {
			objProps.PortMonitoring = nil
		}
		if !op.SortByName.IsNull() {
			objProps.SortByName = openapi.PtrBool(op.SortByName.ValueBool())
		}
		if !op.Label.IsNull() {
			objProps.Label = openapi.PtrString(op.Label.ValueString())
		}
		if !op.Icon.IsNull() {
			objProps.Icon = openapi.PtrString(op.Icon.ValueString())
		}
		ethPortProfileProps.ObjectProperties = &objProps
	}

	// Handle Services
	if len(plan.Services) > 0 {
		servicesItems := make([]openapi.EthportprofilesPutRequestEthPortProfileValueServicesInner, len(plan.Services))
		for i, service := range plan.Services {
			serviceItem := openapi.EthportprofilesPutRequestEthPortProfileValueServicesInner{}

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
				intVal := int32(service.RowNumExternalVlan.ValueInt64())
				serviceItem.RowNumExternalVlan = *openapi.NewNullableInt32(&intVal)
			} else {
				serviceItem.RowNumExternalVlan = *openapi.NewNullableInt32(nil)
			}
			if !service.RowNumIngressAcl.IsNull() {
				serviceItem.RowNumIngressAcl = openapi.PtrString(service.RowNumIngressAcl.ValueString())
			}
			if !service.RowNumIngressAclRefType.IsNull() {
				serviceItem.RowNumIngressAclRefType = openapi.PtrString(service.RowNumIngressAclRefType.ValueString())
			}
			if !service.RowNumEgressAcl.IsNull() {
				serviceItem.RowNumEgressAcl = openapi.PtrString(service.RowNumEgressAcl.ValueString())
			}
			if !service.RowNumEgressAclRefType.IsNull() {
				serviceItem.RowNumEgressAclRefType = openapi.PtrString(service.RowNumEgressAclRefType.ValueString())
			}
			if !service.Index.IsNull() {
				serviceItem.Index = openapi.PtrInt32(int32(service.Index.ValueInt64()))
			}
			if !service.RowNumMacFilter.IsNull() {
				serviceItem.RowNumMacFilter = openapi.PtrString(service.RowNumMacFilter.ValueString())
			}
			if !service.RowNumMacFilterRefType.IsNull() {
				serviceItem.RowNumMacFilterRefType = openapi.PtrString(service.RowNumMacFilterRefType.ValueString())
			}
			if !service.RowNumLanIptv.IsNull() {
				serviceItem.RowNumLanIptv = openapi.PtrString(service.RowNumLanIptv.ValueString())
			}

			servicesItems[i] = serviceItem
		}
		ethPortProfileProps.Services = servicesItems
	}

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "eth_port_profile", name, *ethPortProfileProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Eth Port Profile %s creation operation completed successfully", name))
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
			fmt.Sprintf("Error authenticating with API: %s", err),
		)
		return
	}

	ethPortProfileName := state.Name.ValueString()

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("eth_port_profile") {
		tflog.Info(ctx, fmt.Sprintf("Skipping eth port profile %s verification – trusting recent successful API operation", ethPortProfileName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching eth port profiles for verification of %s", ethPortProfileName))

	type EthPortProfileResponse struct {
		EthPortProfile map[string]interface{} `json:"eth_port_profile_"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "eth_port_profiles", ethPortProfileName,
		func() (EthPortProfileResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch eth port profiles")
			respAPI, err := r.client.EthPortProfilesAPI.EthportprofilesGet(ctx).Execute()
			if err != nil {
				return EthPortProfileResponse{}, fmt.Errorf("error reading eth port profiles: %v", err)
			}
			defer respAPI.Body.Close()

			var res EthPortProfileResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return EthPortProfileResponse{}, fmt.Errorf("failed to decode eth port profiles response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d eth port profiles", len(res.EthPortProfile)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Eth Port Profile %s", ethPortProfileName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for eth port profile with name: %s", ethPortProfileName))

	ethPortProfileData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.EthPortProfile,
		ethPortProfileName,
		func(data interface{}) (string, bool) {
			if ethPortProfile, ok := data.(map[string]interface{}); ok {
				if name, ok := ethPortProfile["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Eth Port Profile with name '%s' not found in API response", ethPortProfileName))
		resp.State.RemoveResource(ctx)
		return
	}

	ethPortProfileMap, ok := ethPortProfileData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Eth Port Profile Data",
			fmt.Sprintf("Eth Port Profile data is not in expected format for %s", ethPortProfileName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found eth port profile '%s' under API key '%s'", ethPortProfileName, actualAPIName))

	state.Name = utils.MapStringFromAPI(ethPortProfileMap["name"])

	// Handle object properties
	if objProps, ok := ethPortProfileMap["object_properties"].(map[string]interface{}); ok {
		state.ObjectProperties = []verityEthPortProfileObjectPropertiesModel{
			{
				Group:          utils.MapStringFromAPI(objProps["group"]),
				PortMonitoring: utils.MapStringFromAPI(objProps["port_monitoring"]),
				SortByName:     utils.MapBoolFromAPI(objProps["sort_by_name"]),
				Label:          utils.MapStringFromAPI(objProps["label"]),
				Icon:           utils.MapStringFromAPI(objProps["icon"]),
			},
		}
	} else {
		state.ObjectProperties = nil
	}

	// Map string fields
	stringFieldMappings := map[string]*types.String{
		"ingress_acl":           &state.IngressAcl,
		"ingress_acl_ref_type_": &state.IngressAclRefType,
		"egress_acl":            &state.EgressAcl,
		"egress_acl_ref_type_":  &state.EgressAclRefType,
		"tls_service":           &state.TlsService,
		"tls_service_ref_type_": &state.TlsServiceRefType,
	}

	for apiKey, stateField := range stringFieldMappings {
		*stateField = utils.MapStringFromAPI(ethPortProfileMap[apiKey])
	}

	// Map boolean fields
	boolFieldMappings := map[string]*types.Bool{
		"enable":       &state.Enable,
		"tls":          &state.Tls,
		"trusted_port": &state.TrustedPort,
	}

	for apiKey, stateField := range boolFieldMappings {
		*stateField = utils.MapBoolFromAPI(ethPortProfileMap[apiKey])
	}

	// Handle services
	if services, ok := ethPortProfileMap["services"].([]interface{}); ok && len(services) > 0 {
		var servicesList []servicesModel

		for _, s := range services {
			service, ok := s.(map[string]interface{})
			if !ok {
				continue
			}

			serviceModel := servicesModel{
				RowNumEnable:            utils.MapBoolFromAPI(service["row_num_enable"]),
				RowNumService:           utils.MapStringFromAPI(service["row_num_service"]),
				RowNumServiceRefType:    utils.MapStringFromAPI(service["row_num_service_ref_type_"]),
				RowNumExternalVlan:      utils.MapInt64FromAPI(service["row_num_external_vlan"]),
				RowNumIngressAcl:        utils.MapStringFromAPI(service["row_num_ingress_acl"]),
				RowNumIngressAclRefType: utils.MapStringFromAPI(service["row_num_ingress_acl_ref_type_"]),
				RowNumEgressAcl:         utils.MapStringFromAPI(service["row_num_egress_acl"]),
				RowNumEgressAclRefType:  utils.MapStringFromAPI(service["row_num_egress_acl_ref_type_"]),
				Index:                   utils.MapInt64FromAPI(service["index"]),
				RowNumMacFilter:         utils.MapStringFromAPI(service["row_num_mac_filter"]),
				RowNumMacFilterRefType:  utils.MapStringFromAPI(service["row_num_mac_filter_ref_type_"]),
				RowNumLanIptv:           utils.MapStringFromAPI(service["row_num_lan_iptv"]),
			}

			servicesList = append(servicesList, serviceModel)
		}

		state.Services = servicesList
	} else {
		state.Services = nil
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

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Authenticate",
			fmt.Sprintf("Error authenticating with API: %s", err),
		)
		return
	}

	name := plan.Name.ValueString()
	ethPortProfileProps := openapi.EthportprofilesPutRequestEthPortProfileValue{}
	hasChanges := false

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { ethPortProfileProps.Name = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { ethPortProfileProps.Enable = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.Tls, state.Tls, func(v *bool) { ethPortProfileProps.Tls = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.TrustedPort, state.TrustedPort, func(v *bool) { ethPortProfileProps.TrustedPort = v }, &hasChanges)

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		if len(state.ObjectProperties) == 0 ||
			!plan.ObjectProperties[0].Group.Equal(state.ObjectProperties[0].Group) ||
			!plan.ObjectProperties[0].PortMonitoring.Equal(state.ObjectProperties[0].PortMonitoring) ||
			!plan.ObjectProperties[0].SortByName.Equal(state.ObjectProperties[0].SortByName) ||
			!plan.ObjectProperties[0].Label.Equal(state.ObjectProperties[0].Label) ||
			!plan.ObjectProperties[0].Icon.Equal(state.ObjectProperties[0].Icon) {

			objProps := openapi.EthportprofilesPutRequestEthPortProfileValueObjectProperties{}
			op := plan.ObjectProperties[0]

			if !op.Group.IsNull() {
				objProps.Group = openapi.PtrString(op.Group.ValueString())
			} else {
				objProps.Group = nil
			}
			if !op.PortMonitoring.IsNull() {
				objProps.PortMonitoring = openapi.PtrString(op.PortMonitoring.ValueString())
			} else {
				objProps.PortMonitoring = nil
			}
			if !op.SortByName.IsNull() {
				objProps.SortByName = openapi.PtrBool(op.SortByName.ValueBool())
			} else {
				objProps.SortByName = nil
			}
			if !op.Label.IsNull() {
				objProps.Label = openapi.PtrString(op.Label.ValueString())
			} else {
				objProps.Label = nil
			}
			if !op.Icon.IsNull() {
				objProps.Icon = openapi.PtrString(op.Icon.ValueString())
			} else {
				objProps.Icon = nil
			}
			ethPortProfileProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	// Handle ingress_acl and ingress_acl_ref_type_ using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.IngressAcl, state.IngressAcl, plan.IngressAclRefType, state.IngressAclRefType,
		func(v *string) { ethPortProfileProps.IngressAcl = v },
		func(v *string) { ethPortProfileProps.IngressAclRefType = v },
		"ingress_acl", "ingress_acl_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	// Handle egress_acl and egress_acl_ref_type_ using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.EgressAcl, state.EgressAcl, plan.EgressAclRefType, state.EgressAclRefType,
		func(v *string) { ethPortProfileProps.EgressAcl = v },
		func(v *string) { ethPortProfileProps.EgressAclRefType = v },
		"egress_acl", "egress_acl_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	// Handle tls_service and tls_service_ref_type_ using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.TlsService, state.TlsService, plan.TlsServiceRefType, state.TlsServiceRefType,
		func(v *string) { ethPortProfileProps.TlsService = v },
		func(v *string) { ethPortProfileProps.TlsServiceRefType = v },
		"tls_service", "tls_service_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	// Handle services
	servicesHandler := utils.IndexedItemHandler[servicesModel, openapi.EthportprofilesPutRequestEthPortProfileValueServicesInner]{
		CreateNew: func(planItem servicesModel) openapi.EthportprofilesPutRequestEthPortProfileValueServicesInner {
			service := openapi.EthportprofilesPutRequestEthPortProfileValueServicesInner{
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
				intVal := int32(planItem.RowNumExternalVlan.ValueInt64())
				service.RowNumExternalVlan = *openapi.NewNullableInt32(&intVal)
			} else {
				service.RowNumExternalVlan = *openapi.NewNullableInt32(nil)
			}

			if !planItem.RowNumIngressAcl.IsNull() {
				service.RowNumIngressAcl = openapi.PtrString(planItem.RowNumIngressAcl.ValueString())
			} else {
				service.RowNumIngressAcl = openapi.PtrString("")
			}

			if !planItem.RowNumIngressAclRefType.IsNull() {
				service.RowNumIngressAclRefType = openapi.PtrString(planItem.RowNumIngressAclRefType.ValueString())
			} else {
				service.RowNumIngressAclRefType = openapi.PtrString("")
			}

			if !planItem.RowNumEgressAcl.IsNull() {
				service.RowNumEgressAcl = openapi.PtrString(planItem.RowNumEgressAcl.ValueString())
			} else {
				service.RowNumEgressAcl = openapi.PtrString("")
			}

			if !planItem.RowNumEgressAclRefType.IsNull() {
				service.RowNumEgressAclRefType = openapi.PtrString(planItem.RowNumEgressAclRefType.ValueString())
			} else {
				service.RowNumEgressAclRefType = openapi.PtrString("")
			}

			if !planItem.RowNumMacFilter.IsNull() {
				service.RowNumMacFilter = openapi.PtrString(planItem.RowNumMacFilter.ValueString())
			} else {
				service.RowNumMacFilter = openapi.PtrString("")
			}

			if !planItem.RowNumMacFilterRefType.IsNull() {
				service.RowNumMacFilterRefType = openapi.PtrString(planItem.RowNumMacFilterRefType.ValueString())
			} else {
				service.RowNumMacFilterRefType = openapi.PtrString("")
			}

			if !planItem.RowNumLanIptv.IsNull() {
				service.RowNumLanIptv = openapi.PtrString(planItem.RowNumLanIptv.ValueString())
			} else {
				service.RowNumLanIptv = openapi.PtrString("")
			}

			return service
		},
		UpdateExisting: func(planItem servicesModel, stateItem servicesModel) (openapi.EthportprofilesPutRequestEthPortProfileValueServicesInner, bool) {
			service := openapi.EthportprofilesPutRequestEthPortProfileValueServicesInner{
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
					intVal := int32(planItem.RowNumExternalVlan.ValueInt64())
					service.RowNumExternalVlan = *openapi.NewNullableInt32(&intVal)
				} else {
					service.RowNumExternalVlan = *openapi.NewNullableInt32(nil)
				}
				fieldChanged = true
			}

			if !planItem.RowNumIngressAcl.Equal(stateItem.RowNumIngressAcl) {
				if !planItem.RowNumIngressAcl.IsNull() {
					service.RowNumIngressAcl = openapi.PtrString(planItem.RowNumIngressAcl.ValueString())
				} else {
					service.RowNumIngressAcl = openapi.PtrString("")
				}
				fieldChanged = true
			}

			if !planItem.RowNumIngressAclRefType.Equal(stateItem.RowNumIngressAclRefType) {
				if !planItem.RowNumIngressAclRefType.IsNull() {
					service.RowNumIngressAclRefType = openapi.PtrString(planItem.RowNumIngressAclRefType.ValueString())
				} else {
					service.RowNumIngressAclRefType = openapi.PtrString("")
				}
				fieldChanged = true
			}

			if !planItem.RowNumEgressAcl.Equal(stateItem.RowNumEgressAcl) {
				if !planItem.RowNumEgressAcl.IsNull() {
					service.RowNumEgressAcl = openapi.PtrString(planItem.RowNumEgressAcl.ValueString())
				} else {
					service.RowNumEgressAcl = openapi.PtrString("")
				}
				fieldChanged = true
			}

			if !planItem.RowNumEgressAclRefType.Equal(stateItem.RowNumEgressAclRefType) {
				if !planItem.RowNumEgressAclRefType.IsNull() {
					service.RowNumEgressAclRefType = openapi.PtrString(planItem.RowNumEgressAclRefType.ValueString())
				} else {
					service.RowNumEgressAclRefType = openapi.PtrString("")
				}
				fieldChanged = true
			}

			if !planItem.RowNumMacFilter.Equal(stateItem.RowNumMacFilter) {
				if !planItem.RowNumMacFilter.IsNull() {
					service.RowNumMacFilter = openapi.PtrString(planItem.RowNumMacFilter.ValueString())
				} else {
					service.RowNumMacFilter = openapi.PtrString("")
				}
				fieldChanged = true
			}

			if !planItem.RowNumMacFilterRefType.Equal(stateItem.RowNumMacFilterRefType) {
				if !planItem.RowNumMacFilterRefType.IsNull() {
					service.RowNumMacFilterRefType = openapi.PtrString(planItem.RowNumMacFilterRefType.ValueString())
				} else {
					service.RowNumMacFilterRefType = openapi.PtrString("")
				}
				fieldChanged = true
			}

			if !planItem.RowNumLanIptv.Equal(stateItem.RowNumLanIptv) {
				if !planItem.RowNumLanIptv.IsNull() {
					service.RowNumLanIptv = openapi.PtrString(planItem.RowNumLanIptv.ValueString())
				} else {
					service.RowNumLanIptv = openapi.PtrString("")
				}
				fieldChanged = true
			}

			return service, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.EthportprofilesPutRequestEthPortProfileValueServicesInner {
			return openapi.EthportprofilesPutRequestEthPortProfileValueServicesInner{
				Index: openapi.PtrInt32(int32(index)),
			}
		},
	}

	changedServices, servicesChanged := utils.ProcessIndexedArrayUpdates(plan.Services, state.Services, servicesHandler)
	if servicesChanged {
		ethPortProfileProps.Services = changedServices
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "eth_port_profile", name, ethPortProfileProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Eth Port Profile %s update operation completed successfully", name))
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

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError(
			"Failed to Authenticate",
			fmt.Sprintf("Error authenticating with API: %s", err),
		)
		return
	}

	name := state.Name.ValueString()

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "eth_port_profile", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Eth Port Profile %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "eth_port_profiles")
	resp.State.RemoveResource(ctx)
}

func (r *verityEthPortProfileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
