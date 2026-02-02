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

	"terraform-provider-verity/internal/bulkops"
	"terraform-provider-verity/internal/utils"
	"terraform-provider-verity/openapi"
)

var (
	_ resource.Resource                = &verityServicePortProfileResource{}
	_ resource.ResourceWithConfigure   = &verityServicePortProfileResource{}
	_ resource.ResourceWithImportState = &verityServicePortProfileResource{}
	_ resource.ResourceWithModifyPlan  = &verityServicePortProfileResource{}
)

const servicePortProfileResourceType = "serviceportprofiles"
const servicePortProfileTerraformType = "verity_service_port_profile"

func NewVerityServicePortProfileResource() resource.Resource {
	return &verityServicePortProfileResource{}
}

type verityServicePortProfileResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
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
				Computed:    true,
			},
			"port_type": schema.StringAttribute{
				Description: "Determines what Service are provisioned on the port and if those Services are propagated upstream",
				Optional:    true,
				Computed:    true,
			},
			"tls_limit_in": schema.Int64Attribute{
				Description: "Speed of ingress (Mbps) for TLS (Transparent LAN Service)",
				Optional:    true,
				Computed:    true,
			},
			"tls_service": schema.StringAttribute{
				Description: "Service used for TLS (Transparent LAN Service)",
				Optional:    true,
				Computed:    true,
			},
			"tls_service_ref_type_": schema.StringAttribute{
				Description: "Object type for tls_service field",
				Optional:    true,
				Computed:    true,
			},
			"trusted_port": schema.BoolAttribute{
				Description: "Trusted Ports do not participate in IP Source Guard, Dynamic ARP Inspection, nor DHCP Snooping, meaning all packets are forwarded without any checks.",
				Optional:    true,
				Computed:    true,
			},
			"ip_mask": schema.StringAttribute{
				Description: "IP/Mask",
				Optional:    true,
				Computed:    true,
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
							Computed:    true,
						},
						"row_num_service": schema.StringAttribute{
							Description: "Connect a Service",
							Optional:    true,
							Computed:    true,
						},
						"row_num_service_ref_type_": schema.StringAttribute{
							Description: "Object type for row_num_service field",
							Optional:    true,
							Computed:    true,
						},
						"row_num_external_vlan": schema.Int64Attribute{
							Description: "Choose an external vlan",
							Optional:    true,
							Computed:    true,
						},
						"row_num_limit_in": schema.Int64Attribute{
							Description: "Speed of ingress (Mbps)",
							Optional:    true,
							Computed:    true,
						},
						"row_num_limit_out": schema.Int64Attribute{
							Description: "Speed of egress (Mbps)",
							Optional:    true,
							Computed:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object. Zero if you want to add an object to the list.",
							Optional:    true,
							Computed:    true,
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
							Computed:    true,
						},
						"port_monitoring": schema.StringAttribute{
							Description: "Defines importance of Link Down on this port",
							Optional:    true,
							Computed:    true,
						},
						"group": schema.StringAttribute{
							Description: "Group",
							Optional:    true,
							Computed:    true,
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

	var config verityServicePortProfileResourceModel
	diags = req.Config.Get(ctx, &config)
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

	// Handle nullable int64 fields - parse HCL to detect explicit config
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, servicePortProfileTerraformType, name)

	utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
		{FieldName: "TlsLimitIn", APIField: &sppProps.TlsLimitIn, TFValue: config.TlsLimitIn, IsConfigured: configuredAttrs.IsConfigured("tls_limit_in")},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.ServiceportprofilesPutRequestServicePortProfileValueObjectProperties{}
		utils.SetObjectPropertiesFields([]utils.ObjectPropertiesField{
			{Name: "OnSummary", TFValue: op.OnSummary, APIValue: &objProps.OnSummary},
			{Name: "PortMonitoring", TFValue: op.PortMonitoring, APIValue: &objProps.PortMonitoring},
			{Name: "Group", TFValue: op.Group, APIValue: &objProps.Group},
		})
		sppProps.ObjectProperties = &objProps
	}

	// Handle services
	if len(plan.Services) > 0 {
		services := make([]openapi.ServiceportprofilesPutRequestServicePortProfileValueServicesInner, len(plan.Services))
		servicesConfigMap := utils.BuildIndexedConfigMap(config.Services)
		for i, service := range plan.Services {
			serviceItem := openapi.ServiceportprofilesPutRequestServicePortProfileValueServicesInner{}

			// Handle boolean fields
			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "RowNumEnable", APIField: &serviceItem.RowNumEnable, TFValue: service.RowNumEnable},
			})

			// Handle string fields
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "RowNumService", APIField: &serviceItem.RowNumService, TFValue: service.RowNumService},
				{FieldName: "RowNumServiceRefType", APIField: &serviceItem.RowNumServiceRefType, TFValue: service.RowNumServiceRefType},
			})

			// Handle int64 fields
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &serviceItem.Index, TFValue: service.Index},
			})

			// Get per-block configured info for nullable Int64 fields
			configItem, cfg := utils.GetIndexedBlockConfig(service, servicesConfigMap, "services", configuredAttrs)
			utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
				{FieldName: "RowNumExternalVlan", APIField: &serviceItem.RowNumExternalVlan, TFValue: configItem.RowNumExternalVlan, IsConfigured: cfg.IsFieldConfigured("row_num_external_vlan")},
				{FieldName: "RowNumLimitIn", APIField: &serviceItem.RowNumLimitIn, TFValue: configItem.RowNumLimitIn, IsConfigured: cfg.IsFieldConfigured("row_num_limit_in")},
				{FieldName: "RowNumLimitOut", APIField: &serviceItem.RowNumLimitOut, TFValue: configItem.RowNumLimitOut, IsConfigured: cfg.IsFieldConfigured("row_num_limit_out")},
			})

			services[i] = serviceItem
		}
		sppProps.Services = services
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "service_port_profile", name, *sppProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Service Port Profile %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "service_port_profiles")

	var minState verityServicePortProfileResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if sppData, exists := bulkMgr.GetResourceResponse("service_port_profile", name); exists {
			state := populateServicePortProfileState(ctx, minState, sppData, r.provCtx.mode)
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

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if sppData, exists := r.bulkOpsMgr.GetResourceResponse("service_port_profile", sppName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached service_port_profile data for %s from recent operation", sppName))
			state = populateServicePortProfileState(ctx, state, sppData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

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

	state = populateServicePortProfileState(ctx, state, sppMap, r.provCtx.mode)
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

	// Get config for nullable field handling
	var config verityServicePortProfileResourceModel
	diags = req.Config.Get(ctx, &config)
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

	// Parse HCL to detect which fields are explicitly configured
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, servicePortProfileTerraformType, name)

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { sppProps.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.PortType, state.PortType, func(v *string) { sppProps.PortType = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.TlsService, state.TlsService, func(v *string) { sppProps.TlsService = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.TlsServiceRefType, state.TlsServiceRefType, func(v *string) { sppProps.TlsServiceRefType = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.IpMask, state.IpMask, func(v *string) { sppProps.IpMask = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { sppProps.Enable = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.TrustedPort, state.TrustedPort, func(v *bool) { sppProps.TrustedPort = v }, &hasChanges)

	// Handle nullable int64 field changes - parse HCL to detect explicit config
	utils.CompareAndSetNullableInt64Field(config.TlsLimitIn, state.TlsLimitIn, configuredAttrs.IsConfigured("tls_limit_in"), func(v *openapi.NullableInt32) { sppProps.TlsLimitIn = *v }, &hasChanges)

	// Handle object properties
	if len(plan.ObjectProperties) > 0 && len(state.ObjectProperties) > 0 {
		objProps := openapi.ServiceportprofilesPutRequestServicePortProfileValueObjectProperties{}
		op := plan.ObjectProperties[0]
		st := state.ObjectProperties[0]
		objPropsChanged := false

		utils.CompareAndSetObjectPropertiesFields([]utils.ObjectPropertiesFieldWithComparison{
			{Name: "OnSummary", PlanValue: op.OnSummary, StateValue: st.OnSummary, APIValue: &objProps.OnSummary},
			{Name: "PortMonitoring", PlanValue: op.PortMonitoring, StateValue: st.PortMonitoring, APIValue: &objProps.PortMonitoring},
			{Name: "Group", PlanValue: op.Group, StateValue: st.Group, APIValue: &objProps.Group},
		}, &objPropsChanged)

		if objPropsChanged {
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
	servicesConfigMap := utils.BuildIndexedConfigMap(config.Services)

	changedServices, servicesChanged := utils.ProcessIndexedArrayUpdates(plan.Services, state.Services,
		utils.IndexedItemHandler[verityServicePortProfileServiceModel, openapi.ServiceportprofilesPutRequestServicePortProfileValueServicesInner]{
			CreateNew: func(planItem verityServicePortProfileServiceModel) openapi.ServiceportprofilesPutRequestServicePortProfileValueServicesInner {
				newService := openapi.ServiceportprofilesPutRequestServicePortProfileValueServicesInner{}

				// Handle boolean fields
				utils.SetBoolFields([]utils.BoolFieldMapping{
					{FieldName: "RowNumEnable", APIField: &newService.RowNumEnable, TFValue: planItem.RowNumEnable},
				})

				// Handle string fields
				utils.SetStringFields([]utils.StringFieldMapping{
					{FieldName: "RowNumService", APIField: &newService.RowNumService, TFValue: planItem.RowNumService},
					{FieldName: "RowNumServiceRefType", APIField: &newService.RowNumServiceRefType, TFValue: planItem.RowNumServiceRefType},
				})

				// Handle int64 fields
				utils.SetInt64Fields([]utils.Int64FieldMapping{
					{FieldName: "Index", APIField: &newService.Index, TFValue: planItem.Index},
				})

				// Get per-block configured info for nullable Int64 fields
				configItem, cfg := utils.GetIndexedBlockConfig(planItem, servicesConfigMap, "services", configuredAttrs)
				utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
					{FieldName: "RowNumExternalVlan", APIField: &newService.RowNumExternalVlan, TFValue: configItem.RowNumExternalVlan, IsConfigured: cfg.IsFieldConfigured("row_num_external_vlan")},
					{FieldName: "RowNumLimitIn", APIField: &newService.RowNumLimitIn, TFValue: configItem.RowNumLimitIn, IsConfigured: cfg.IsFieldConfigured("row_num_limit_in")},
					{FieldName: "RowNumLimitOut", APIField: &newService.RowNumLimitOut, TFValue: configItem.RowNumLimitOut, IsConfigured: cfg.IsFieldConfigured("row_num_limit_out")},
				})

				return newService
			},
			UpdateExisting: func(planItem verityServicePortProfileServiceModel, stateItem verityServicePortProfileServiceModel) (openapi.ServiceportprofilesPutRequestServicePortProfileValueServicesInner, bool) {
				updateService := openapi.ServiceportprofilesPutRequestServicePortProfileValueServicesInner{}
				fieldChanged := false

				// Handle boolean field changes
				utils.CompareAndSetBoolField(planItem.RowNumEnable, stateItem.RowNumEnable, func(v *bool) { updateService.RowNumEnable = v }, &fieldChanged)

				// Handle row_num_service and row_num_service_ref_type_ using one ref type supported pattern
				if !utils.HandleOneRefTypeSupported(
					planItem.RowNumService, stateItem.RowNumService, planItem.RowNumServiceRefType, stateItem.RowNumServiceRefType,
					func(v *string) { updateService.RowNumService = v },
					func(v *string) { updateService.RowNumServiceRefType = v },
					"row_num_service", "row_num_service_ref_type_",
					&fieldChanged, &resp.Diagnostics,
				) {
					return updateService, false
				}

				// Handle index field change
				utils.CompareAndSetInt64Field(planItem.Index, stateItem.Index, func(v *int32) { updateService.Index = v }, &fieldChanged)

				// Handle nullable int64 field changes
				configItem, cfg := utils.GetIndexedBlockConfig(planItem, servicesConfigMap, "services", configuredAttrs)
				utils.CompareAndSetNullableInt64Field(configItem.RowNumExternalVlan, stateItem.RowNumExternalVlan, cfg.IsFieldConfigured("row_num_external_vlan"), func(v *openapi.NullableInt32) { updateService.RowNumExternalVlan = *v }, &fieldChanged)
				utils.CompareAndSetNullableInt64Field(configItem.RowNumLimitIn, stateItem.RowNumLimitIn, cfg.IsFieldConfigured("row_num_limit_in"), func(v *openapi.NullableInt32) { updateService.RowNumLimitIn = *v }, &fieldChanged)
				utils.CompareAndSetNullableInt64Field(configItem.RowNumLimitOut, stateItem.RowNumLimitOut, cfg.IsFieldConfigured("row_num_limit_out"), func(v *openapi.NullableInt32) { updateService.RowNumLimitOut = *v }, &fieldChanged)

				return updateService, fieldChanged
			},
			CreateDeleted: func(index int64) openapi.ServiceportprofilesPutRequestServicePortProfileValueServicesInner {
				return openapi.ServiceportprofilesPutRequestServicePortProfileValueServicesInner{
					Index: openapi.PtrInt32(int32(index)),
				}
			},
		})
	if servicesChanged {
		sppProps.Services = changedServices
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "service_port_profile", name, sppProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Service Port Profile %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "service_port_profiles")

	var minState verityServicePortProfileResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Try to use cached response from bulk operation to populate state with API values
	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if sppData, exists := bulkMgr.GetResourceResponse("service_port_profile", name); exists {
			newState := populateServicePortProfileState(ctx, minState, sppData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "service_port_profile", name, nil, &resp.Diagnostics)
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

func populateServicePortProfileState(ctx context.Context, state verityServicePortProfileResourceModel, data map[string]interface{}, mode string) verityServicePortProfileResourceModel {
	const resourceType = servicePortProfileResourceType

	state.Name = utils.MapStringFromAPI(data["name"])

	// String fields
	state.PortType = utils.MapStringWithMode(data, "port_type", resourceType, mode)
	state.TlsService = utils.MapStringWithMode(data, "tls_service", resourceType, mode)
	state.TlsServiceRefType = utils.MapStringWithMode(data, "tls_service_ref_type_", resourceType, mode)
	state.IpMask = utils.MapStringWithMode(data, "ip_mask", resourceType, mode)

	// Boolean fields
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)
	state.TrustedPort = utils.MapBoolWithMode(data, "trusted_port", resourceType, mode)

	// Int64 fields
	state.TlsLimitIn = utils.MapInt64WithMode(data, "tls_limit_in", resourceType, mode)

	// Handle object_properties block
	if utils.FieldAppliesToMode(resourceType, "object_properties", mode) {
		if objProps, ok := data["object_properties"].(map[string]interface{}); ok {
			objPropsModel := verityServicePortProfileObjectPropertiesModel{
				OnSummary:      utils.MapBoolWithModeNested(objProps, "on_summary", resourceType, "object_properties.on_summary", mode),
				PortMonitoring: utils.MapStringWithModeNested(objProps, "port_monitoring", resourceType, "object_properties.port_monitoring", mode),
				Group:          utils.MapStringWithModeNested(objProps, "group", resourceType, "object_properties.group", mode),
			}
			state.ObjectProperties = []verityServicePortProfileObjectPropertiesModel{objPropsModel}
		} else {
			state.ObjectProperties = nil
		}
	} else {
		state.ObjectProperties = nil
	}

	// Handle services block
	if utils.FieldAppliesToMode(resourceType, "services", mode) {
		if servicesArray, ok := data["services"].([]interface{}); ok && len(servicesArray) > 0 {
			var services []verityServicePortProfileServiceModel
			for _, s := range servicesArray {
				service, ok := s.(map[string]interface{})
				if !ok {
					continue
				}

				srModel := verityServicePortProfileServiceModel{
					RowNumEnable:         utils.MapBoolWithModeNested(service, "row_num_enable", resourceType, "services.row_num_enable", mode),
					RowNumService:        utils.MapStringWithModeNested(service, "row_num_service", resourceType, "services.row_num_service", mode),
					RowNumServiceRefType: utils.MapStringWithModeNested(service, "row_num_service_ref_type_", resourceType, "services.row_num_service_ref_type_", mode),
					RowNumExternalVlan:   utils.MapInt64WithModeNested(service, "row_num_external_vlan", resourceType, "services.row_num_external_vlan", mode),
					RowNumLimitIn:        utils.MapInt64WithModeNested(service, "row_num_limit_in", resourceType, "services.row_num_limit_in", mode),
					RowNumLimitOut:       utils.MapInt64WithModeNested(service, "row_num_limit_out", resourceType, "services.row_num_limit_out", mode),
					Index:                utils.MapInt64WithModeNested(service, "index", resourceType, "services.index", mode),
				}

				services = append(services, srModel)
			}
			state.Services = services
		} else {
			state.Services = nil
		}
	} else {
		state.Services = nil
	}

	return state
}

func (r *verityServicePortProfileResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// =========================================================================
	// Skip if deleting
	// =========================================================================
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityServicePortProfileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Mode-aware field nullification
	// Set fields that don't apply to current mode to null to prevent
	// "known after apply" messages for irrelevant fields.
	// =========================================================================
	const resourceType = servicePortProfileResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyBools(
		"enable",
		"trusted_port",
	)

	nullifier.NullifyStrings(
		"port_type",
		"tls_service",
		"tls_service_ref_type_",
		"ip_mask",
	)

	nullifier.NullifyInt64s(
		"tls_limit_in",
	)

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "services",
		ItemCount:    len(plan.Services),
		StringFields: []string{"row_num_service", "row_num_service_ref_type_"},
		BoolFields:   []string{"row_num_enable"},
		Int64Fields:  []string{"index", "row_num_external_vlan", "row_num_limit_in", "row_num_limit_out"},
	})

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "object_properties",
		ItemCount:    len(plan.ObjectProperties),
		StringFields: []string{"group", "port_monitoring"},
		BoolFields:   []string{"on_summary"},
	})

	// =========================================================================
	// Skip UPDATE-specific logic during CREATE
	// =========================================================================
	if req.State.Raw.IsNull() {
		return
	}

	// =========================================================================
	// UPDATE operation - get state and config
	// =========================================================================
	var state verityServicePortProfileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config verityServicePortProfileResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Handle nullable Int64 fields (explicit null detection)
	// For Optional+Computed fields, Terraform copies state to plan when config
	// is null. We detect explicit null in HCL and force plan to null.
	// =========================================================================
	name := plan.Name.ValueString()
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, servicePortProfileTerraformType, name)

	utils.HandleNullableFields(utils.NullableFieldsConfig{
		Ctx:             ctx,
		Plan:            &resp.Plan,
		ConfiguredAttrs: configuredAttrs,
		Int64Fields: []utils.NullableInt64Field{
			{AttrName: "tls_limit_in", ConfigVal: config.TlsLimitIn, StateVal: state.TlsLimitIn},
		},
	})

	// =========================================================================
	// Handle nullable fields in nested blocks
	// =========================================================================
	for i, configItem := range config.Services {
		itemIndex := configItem.Index.ValueInt64()
		var stateItem *verityServicePortProfileServiceModel
		for j := range state.Services {
			if state.Services[j].Index.ValueInt64() == itemIndex {
				stateItem = &state.Services[j]
				break
			}
		}

		if stateItem != nil {
			utils.HandleNullableNestedFields(utils.NullableNestedFieldsConfig{
				Ctx:             ctx,
				Plan:            &resp.Plan,
				ConfiguredAttrs: configuredAttrs,
				BlockType:       "services",
				BlockListPath:   "services",
				BlockListIndex:  i,
				Int64Fields: []utils.NullableNestedInt64Field{
					{BlockIndex: itemIndex, AttrName: "row_num_external_vlan", ConfigVal: configItem.RowNumExternalVlan, StateVal: stateItem.RowNumExternalVlan},
					{BlockIndex: itemIndex, AttrName: "row_num_limit_in", ConfigVal: configItem.RowNumLimitIn, StateVal: stateItem.RowNumLimitIn},
					{BlockIndex: itemIndex, AttrName: "row_num_limit_out", ConfigVal: configItem.RowNumLimitOut, StateVal: stateItem.RowNumLimitOut},
				},
			})
		}
	}
}
