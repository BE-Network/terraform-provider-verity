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
	_ resource.Resource                = &verityEthPortProfileResource{}
	_ resource.ResourceWithConfigure   = &verityEthPortProfileResource{}
	_ resource.ResourceWithImportState = &verityEthPortProfileResource{}
	_ resource.ResourceWithModifyPlan  = &verityEthPortProfileResource{}
)

const ethPortProfileResourceType = "ethportprofiles"
const ethPortProfileTerraformType = "verity_eth_port_profile"

func NewVerityEthPortProfileResource() resource.Resource {
	return &verityEthPortProfileResource{}
}

type verityEthPortProfileResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
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
				Computed:    true,
			},
			"ingress_acl": schema.StringAttribute{
				Description: "Choose an ingress access control list",
				Optional:    true,
				Computed:    true,
			},
			"ingress_acl_ref_type_": schema.StringAttribute{
				Description: "Object type for ingress_acl field",
				Optional:    true,
				Computed:    true,
			},
			"egress_acl": schema.StringAttribute{
				Description: "Choose an egress access control list",
				Optional:    true,
				Computed:    true,
			},
			"egress_acl_ref_type_": schema.StringAttribute{
				Description: "Object type for egress_acl field",
				Optional:    true,
				Computed:    true,
			},
			"tls": schema.BoolAttribute{
				Description: "Transparent LAN Service Trunk",
				Optional:    true,
				Computed:    true,
			},
			"tls_service": schema.StringAttribute{
				Description: "Choose a Service supporting Transparent LAN Service",
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
		},
		Blocks: map[string]schema.Block{
			"object_properties": schema.ListNestedBlock{
				Description: "Object properties for the profile",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"group": schema.StringAttribute{
							Description: "Group",
							Optional:    true,
							Computed:    true,
						},
						"port_monitoring": schema.StringAttribute{
							Description: "Defines importance of Link Down on this port",
							Optional:    true,
							Computed:    true,
						},
						"sort_by_name": schema.BoolAttribute{
							Description: "Choose to sort by service name or by order of creation",
							Optional:    true,
							Computed:    true,
						},
						"label": schema.StringAttribute{
							Description: "Port Label displayed ports provisioned with this Eth Port Profile but with no Port Label defined in the endpoint",
							Optional:    true,
							Computed:    true,
						},
						"icon": schema.StringAttribute{
							Description: "Port Icon displayed ports provisioned with this Eth Port Profile but with no Port Icon defined in the endpoint",
							Optional:    true,
							Computed:    true,
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
							Computed:    true,
						},
						"row_num_service": schema.StringAttribute{
							Description: "Choose a Service to connect",
							Optional:    true,
							Computed:    true,
						},
						"row_num_service_ref_type_": schema.StringAttribute{
							Description: "Object type for row_num_service field",
							Optional:    true,
							Computed:    true,
						},
						"row_num_external_vlan": schema.Int64Attribute{
							Description: "Choose an external vlan. A value of 0 will make the VLAN untagged, while null will use service VLAN.",
							Optional:    true,
							Computed:    true,
						},
						"row_num_ingress_acl": schema.StringAttribute{
							Description: "Choose an ingress access control list",
							Optional:    true,
							Computed:    true,
						},
						"row_num_ingress_acl_ref_type_": schema.StringAttribute{
							Description: "Object type for row_num_ingress_acl field",
							Optional:    true,
							Computed:    true,
						},
						"row_num_egress_acl": schema.StringAttribute{
							Description: "Choose an egress access control list",
							Optional:    true,
							Computed:    true,
						},
						"row_num_egress_acl_ref_type_": schema.StringAttribute{
							Description: "Object type for row_num_egress_acl field",
							Optional:    true,
							Computed:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object. Zero if you want to add an object to the list.",
							Optional:    true,
							Computed:    true,
						},
						"row_num_mac_filter": schema.StringAttribute{
							Description: "Choose an access control list",
							Optional:    true,
							Computed:    true,
						},
						"row_num_mac_filter_ref_type_": schema.StringAttribute{
							Description: "Object type for row_num_mac_filter field",
							Optional:    true,
							Computed:    true,
						},
						"row_num_lan_iptv": schema.StringAttribute{
							Description: "Denotes a LAN or IPTV service",
							Optional:    true,
							Computed:    true,
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

	var config verityEthPortProfileResourceModel
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
		utils.SetObjectPropertiesFields([]utils.ObjectPropertiesField{
			{Name: "Group", TFValue: op.Group, APIValue: &objProps.Group},
			{Name: "PortMonitoring", TFValue: op.PortMonitoring, APIValue: &objProps.PortMonitoring},
			{Name: "SortByName", TFValue: op.SortByName, APIValue: &objProps.SortByName},
			{Name: "Label", TFValue: op.Label, APIValue: &objProps.Label},
			{Name: "Icon", TFValue: op.Icon, APIValue: &objProps.Icon},
		})
		ethPortProfileProps.ObjectProperties = &objProps
	}

	// Parse HCL to detect explicitly configured attributes
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, ethPortProfileTerraformType, name)

	// Handle Services
	if len(plan.Services) > 0 {
		servicesItems := make([]openapi.EthportprofilesPutRequestEthPortProfileValueServicesInner, len(plan.Services))
		servicesConfigMap := utils.BuildIndexedConfigMap(config.Services)
		for i, item := range plan.Services {
			serviceItem := openapi.EthportprofilesPutRequestEthPortProfileValueServicesInner{}

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "RowNumEnable", APIField: &serviceItem.RowNumEnable, TFValue: item.RowNumEnable},
			})
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "RowNumService", APIField: &serviceItem.RowNumService, TFValue: item.RowNumService},
				{FieldName: "RowNumServiceRefType", APIField: &serviceItem.RowNumServiceRefType, TFValue: item.RowNumServiceRefType},
				{FieldName: "RowNumIngressAcl", APIField: &serviceItem.RowNumIngressAcl, TFValue: item.RowNumIngressAcl},
				{FieldName: "RowNumIngressAclRefType", APIField: &serviceItem.RowNumIngressAclRefType, TFValue: item.RowNumIngressAclRefType},
				{FieldName: "RowNumEgressAcl", APIField: &serviceItem.RowNumEgressAcl, TFValue: item.RowNumEgressAcl},
				{FieldName: "RowNumEgressAclRefType", APIField: &serviceItem.RowNumEgressAclRefType, TFValue: item.RowNumEgressAclRefType},
				{FieldName: "RowNumMacFilter", APIField: &serviceItem.RowNumMacFilter, TFValue: item.RowNumMacFilter},
				{FieldName: "RowNumMacFilterRefType", APIField: &serviceItem.RowNumMacFilterRefType, TFValue: item.RowNumMacFilterRefType},
				{FieldName: "RowNumLanIptv", APIField: &serviceItem.RowNumLanIptv, TFValue: item.RowNumLanIptv},
			})

			// Get per-block configured info for nullable Int64 fields
			configItem, cfg := utils.GetIndexedBlockConfig(item, servicesConfigMap, "services", configuredAttrs)
			utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
				{FieldName: "RowNumExternalVlan", APIField: &serviceItem.RowNumExternalVlan, TFValue: configItem.RowNumExternalVlan, IsConfigured: cfg.IsFieldConfigured("row_num_external_vlan")},
			})
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &serviceItem.Index, TFValue: item.Index},
			})

			servicesItems[i] = serviceItem
		}
		ethPortProfileProps.Services = servicesItems
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "eth_port_profile", name, *ethPortProfileProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Eth Port Profile %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "eth_port_profiles")

	var minState verityEthPortProfileResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if ethPortProfileData, exists := bulkMgr.GetResourceResponse("eth_port_profile", name); exists {
			state := populateEthPortProfileState(ctx, minState, ethPortProfileData, r.provCtx.mode)
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

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if ethPortProfileData, exists := r.bulkOpsMgr.GetResourceResponse("eth_port_profile", ethPortProfileName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached eth port profile data for %s from recent operation", ethPortProfileName))
			state = populateEthPortProfileState(ctx, state, ethPortProfileData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

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

	state = populateEthPortProfileState(ctx, state, ethPortProfileMap, r.provCtx.mode)
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
	if len(plan.ObjectProperties) > 0 && len(state.ObjectProperties) > 0 {
		objProps := openapi.EthportprofilesPutRequestEthPortProfileValueObjectProperties{}
		op := plan.ObjectProperties[0]
		st := state.ObjectProperties[0]
		objPropsChanged := false

		utils.CompareAndSetObjectPropertiesFields([]utils.ObjectPropertiesFieldWithComparison{
			{Name: "Group", PlanValue: op.Group, StateValue: st.Group, APIValue: &objProps.Group},
			{Name: "PortMonitoring", PlanValue: op.PortMonitoring, StateValue: st.PortMonitoring, APIValue: &objProps.PortMonitoring},
			{Name: "SortByName", PlanValue: op.SortByName, StateValue: st.SortByName, APIValue: &objProps.SortByName},
			{Name: "Label", PlanValue: op.Label, StateValue: st.Label, APIValue: &objProps.Label},
			{Name: "Icon", PlanValue: op.Icon, StateValue: st.Icon, APIValue: &objProps.Icon},
		}, &objPropsChanged)

		if objPropsChanged {
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
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, ethPortProfileTerraformType, name)
	var config verityEthPortProfileResourceModel
	req.Config.Get(ctx, &config)
	servicesConfigMap := utils.BuildIndexedConfigMap(config.Services)

	servicesHandler := utils.IndexedItemHandler[servicesModel, openapi.EthportprofilesPutRequestEthPortProfileValueServicesInner]{
		CreateNew: func(planItem servicesModel) openapi.EthportprofilesPutRequestEthPortProfileValueServicesInner {
			service := openapi.EthportprofilesPutRequestEthPortProfileValueServicesInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &service.Index, TFValue: planItem.Index},
			})

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "RowNumEnable", APIField: &service.RowNumEnable, TFValue: planItem.RowNumEnable},
			})

			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "RowNumService", APIField: &service.RowNumService, TFValue: planItem.RowNumService},
				{FieldName: "RowNumServiceRefType", APIField: &service.RowNumServiceRefType, TFValue: planItem.RowNumServiceRefType},
				{FieldName: "RowNumIngressAcl", APIField: &service.RowNumIngressAcl, TFValue: planItem.RowNumIngressAcl},
				{FieldName: "RowNumIngressAclRefType", APIField: &service.RowNumIngressAclRefType, TFValue: planItem.RowNumIngressAclRefType},
				{FieldName: "RowNumEgressAcl", APIField: &service.RowNumEgressAcl, TFValue: planItem.RowNumEgressAcl},
				{FieldName: "RowNumEgressAclRefType", APIField: &service.RowNumEgressAclRefType, TFValue: planItem.RowNumEgressAclRefType},
				{FieldName: "RowNumMacFilter", APIField: &service.RowNumMacFilter, TFValue: planItem.RowNumMacFilter},
				{FieldName: "RowNumMacFilterRefType", APIField: &service.RowNumMacFilterRefType, TFValue: planItem.RowNumMacFilterRefType},
				{FieldName: "RowNumLanIptv", APIField: &service.RowNumLanIptv, TFValue: planItem.RowNumLanIptv},
			})

			// Get per-block configured info for nullable Int64 fields
			configItem, cfg := utils.GetIndexedBlockConfig(planItem, servicesConfigMap, "services", configuredAttrs)
			utils.SetNullableInt64Fields([]utils.NullableInt64FieldMapping{
				{FieldName: "RowNumExternalVlan", APIField: &service.RowNumExternalVlan, TFValue: configItem.RowNumExternalVlan, IsConfigured: cfg.IsFieldConfigured("row_num_external_vlan")},
			})

			return service
		},
		UpdateExisting: func(planItem servicesModel, stateItem servicesModel) (openapi.EthportprofilesPutRequestEthPortProfileValueServicesInner, bool) {
			service := openapi.EthportprofilesPutRequestEthPortProfileValueServicesInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &service.Index, TFValue: planItem.Index},
			})

			fieldChanged := false

			// Handle boolean fields
			utils.CompareAndSetBoolField(planItem.RowNumEnable, stateItem.RowNumEnable, func(v *bool) { service.RowNumEnable = v }, &fieldChanged)

			// Handle nullable int64 fields
			configItem, cfg := utils.GetIndexedBlockConfig(planItem, servicesConfigMap, "services", configuredAttrs)
			utils.CompareAndSetNullableInt64Field(configItem.RowNumExternalVlan, stateItem.RowNumExternalVlan, cfg.IsFieldConfigured("row_num_external_vlan"), func(v *openapi.NullableInt32) { service.RowNumExternalVlan = *v }, &fieldChanged)

			// Handle row_num_service and row_num_service_ref_type_ using "One ref type supported" pattern
			if !utils.HandleOneRefTypeSupported(
				planItem.RowNumService, stateItem.RowNumService, planItem.RowNumServiceRefType, stateItem.RowNumServiceRefType,
				func(v *string) { service.RowNumService = v },
				func(v *string) { service.RowNumServiceRefType = v },
				"row_num_service", "row_num_service_ref_type_",
				&fieldChanged,
				&resp.Diagnostics,
			) {
				return service, false
			}

			// Handle non-ref-type string fields
			utils.CompareAndSetStringField(planItem.RowNumLanIptv, stateItem.RowNumLanIptv, func(v *string) { service.RowNumLanIptv = v }, &fieldChanged)

			// Handle row_num_ingress_acl and row_num_ingress_acl_ref_type_ using "One ref type supported" pattern
			if !utils.HandleOneRefTypeSupported(
				planItem.RowNumIngressAcl, stateItem.RowNumIngressAcl, planItem.RowNumIngressAclRefType, stateItem.RowNumIngressAclRefType,
				func(v *string) { service.RowNumIngressAcl = v },
				func(v *string) { service.RowNumIngressAclRefType = v },
				"row_num_ingress_acl", "row_num_ingress_acl_ref_type_",
				&fieldChanged,
				&resp.Diagnostics,
			) {
				return service, false
			}

			// Handle row_num_egress_acl and row_num_egress_acl_ref_type_ using "One ref type supported" pattern
			if !utils.HandleOneRefTypeSupported(
				planItem.RowNumEgressAcl, stateItem.RowNumEgressAcl, planItem.RowNumEgressAclRefType, stateItem.RowNumEgressAclRefType,
				func(v *string) { service.RowNumEgressAcl = v },
				func(v *string) { service.RowNumEgressAclRefType = v },
				"row_num_egress_acl", "row_num_egress_acl_ref_type_",
				&fieldChanged,
				&resp.Diagnostics,
			) {
				return service, false
			}

			// Handle row_num_mac_filter and row_num_mac_filter_ref_type_ using "One ref type supported" pattern
			if !utils.HandleOneRefTypeSupported(
				planItem.RowNumMacFilter, stateItem.RowNumMacFilter, planItem.RowNumMacFilterRefType, stateItem.RowNumMacFilterRefType,
				func(v *string) { service.RowNumMacFilter = v },
				func(v *string) { service.RowNumMacFilterRefType = v },
				"row_num_mac_filter", "row_num_mac_filter_ref_type_",
				&fieldChanged,
				&resp.Diagnostics,
			) {
				return service, false
			}

			return service, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.EthportprofilesPutRequestEthPortProfileValueServicesInner {
			service := openapi.EthportprofilesPutRequestEthPortProfileValueServicesInner{}
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &service.Index, TFValue: types.Int64Value(index)},
			})
			return service
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "eth_port_profile", name, ethPortProfileProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Eth Port Profile %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "eth_port_profiles")

	var minState verityEthPortProfileResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Try to use cached response from bulk operation to populate state with API values
	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if ethPortProfileData, exists := bulkMgr.GetResourceResponse("eth_port_profile", name); exists {
			newState := populateEthPortProfileState(ctx, minState, ethPortProfileData, r.provCtx.mode)
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "eth_port_profile", name, nil, &resp.Diagnostics)
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

func populateEthPortProfileState(ctx context.Context, state verityEthPortProfileResourceModel, data map[string]interface{}, mode string) verityEthPortProfileResourceModel {
	const resourceType = ethPortProfileResourceType

	state.Name = utils.MapStringFromAPI(data["name"])

	// Boolean fields
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)
	state.Tls = utils.MapBoolWithMode(data, "tls", resourceType, mode)
	state.TrustedPort = utils.MapBoolWithMode(data, "trusted_port", resourceType, mode)

	// String fields
	state.IngressAcl = utils.MapStringWithMode(data, "ingress_acl", resourceType, mode)
	state.IngressAclRefType = utils.MapStringWithMode(data, "ingress_acl_ref_type_", resourceType, mode)
	state.EgressAcl = utils.MapStringWithMode(data, "egress_acl", resourceType, mode)
	state.EgressAclRefType = utils.MapStringWithMode(data, "egress_acl_ref_type_", resourceType, mode)
	state.TlsService = utils.MapStringWithMode(data, "tls_service", resourceType, mode)
	state.TlsServiceRefType = utils.MapStringWithMode(data, "tls_service_ref_type_", resourceType, mode)

	// Handle object_properties block
	if utils.FieldAppliesToMode(resourceType, "object_properties", mode) {
		if objProps, ok := data["object_properties"].(map[string]interface{}); ok {
			objPropsModel := verityEthPortProfileObjectPropertiesModel{
				Group:          utils.MapStringWithModeNested(objProps, "group", resourceType, "object_properties.group", mode),
				PortMonitoring: utils.MapStringWithModeNested(objProps, "port_monitoring", resourceType, "object_properties.port_monitoring", mode),
				SortByName:     utils.MapBoolWithModeNested(objProps, "sort_by_name", resourceType, "object_properties.sort_by_name", mode),
				Label:          utils.MapStringWithModeNested(objProps, "label", resourceType, "object_properties.label", mode),
				Icon:           utils.MapStringWithModeNested(objProps, "icon", resourceType, "object_properties.icon", mode),
			}
			state.ObjectProperties = []verityEthPortProfileObjectPropertiesModel{objPropsModel}
		} else {
			state.ObjectProperties = nil
		}
	} else {
		state.ObjectProperties = nil
	}

	// Handle services list block
	if utils.FieldAppliesToMode(resourceType, "services", mode) {
		if servicesData, ok := data["services"].([]interface{}); ok && len(servicesData) > 0 {
			var servicesList []servicesModel

			for _, item := range servicesData {
				itemMap, ok := item.(map[string]interface{})
				if !ok {
					continue
				}

				serviceItem := servicesModel{
					RowNumEnable:            utils.MapBoolWithModeNested(itemMap, "row_num_enable", resourceType, "services.row_num_enable", mode),
					RowNumService:           utils.MapStringWithModeNested(itemMap, "row_num_service", resourceType, "services.row_num_service", mode),
					RowNumServiceRefType:    utils.MapStringWithModeNested(itemMap, "row_num_service_ref_type_", resourceType, "services.row_num_service_ref_type_", mode),
					RowNumExternalVlan:      utils.MapInt64WithModeNested(itemMap, "row_num_external_vlan", resourceType, "services.row_num_external_vlan", mode),
					RowNumIngressAcl:        utils.MapStringWithModeNested(itemMap, "row_num_ingress_acl", resourceType, "services.row_num_ingress_acl", mode),
					RowNumIngressAclRefType: utils.MapStringWithModeNested(itemMap, "row_num_ingress_acl_ref_type_", resourceType, "services.row_num_ingress_acl_ref_type_", mode),
					RowNumEgressAcl:         utils.MapStringWithModeNested(itemMap, "row_num_egress_acl", resourceType, "services.row_num_egress_acl", mode),
					RowNumEgressAclRefType:  utils.MapStringWithModeNested(itemMap, "row_num_egress_acl_ref_type_", resourceType, "services.row_num_egress_acl_ref_type_", mode),
					Index:                   utils.MapInt64WithModeNested(itemMap, "index", resourceType, "services.index", mode),
					RowNumMacFilter:         utils.MapStringWithModeNested(itemMap, "row_num_mac_filter", resourceType, "services.row_num_mac_filter", mode),
					RowNumMacFilterRefType:  utils.MapStringWithModeNested(itemMap, "row_num_mac_filter_ref_type_", resourceType, "services.row_num_mac_filter_ref_type_", mode),
					RowNumLanIptv:           utils.MapStringWithModeNested(itemMap, "row_num_lan_iptv", resourceType, "services.row_num_lan_iptv", mode),
				}

				servicesList = append(servicesList, serviceItem)
			}

			state.Services = servicesList
		} else {
			state.Services = nil
		}
	} else {
		state.Services = nil
	}

	return state
}

func (r *verityEthPortProfileResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// =========================================================================
	// Skip if deleting
	// =========================================================================
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityEthPortProfileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Mode-aware field nullification
	// Set fields that don't apply to current mode to null to prevent
	// "known after apply" messages for irrelevant fields.
	// =========================================================================
	const resourceType = ethPortProfileResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyStrings(
		"ingress_acl", "ingress_acl_ref_type_",
		"egress_acl", "egress_acl_ref_type_",
		"tls_service", "tls_service_ref_type_",
	)

	nullifier.NullifyBools(
		"enable", "tls", "trusted_port",
	)

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "object_properties",
		ItemCount:    len(plan.ObjectProperties),
		StringFields: []string{"group", "port_monitoring", "label", "icon"},
		BoolFields:   []string{"sort_by_name"},
	})

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName: "services",
		ItemCount: len(plan.Services),
		StringFields: []string{"row_num_service", "row_num_service_ref_type_", "row_num_ingress_acl", "row_num_ingress_acl_ref_type_",
			"row_num_egress_acl", "row_num_egress_acl_ref_type_", "row_num_mac_filter", "row_num_mac_filter_ref_type_",
			"row_num_lan_iptv"},
		BoolFields:  []string{"row_num_enable"},
		Int64Fields: []string{"index", "row_num_external_vlan"},
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
	var state verityEthPortProfileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var config verityEthPortProfileResourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Handle nullable fields in services nested blocks
	// =========================================================================
	name := plan.Name.ValueString()
	workDir := utils.GetWorkingDirectory()
	configuredAttrs := utils.ParseResourceConfiguredAttributes(ctx, workDir, ethPortProfileTerraformType, name)

	for i, configItem := range config.Services {
		itemIndex := configItem.Index.ValueInt64()
		var stateItem *servicesModel
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
				},
			})
		}
	}
}
