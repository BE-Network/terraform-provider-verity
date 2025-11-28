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
	_ resource.Resource                = &verityBundleResource{}
	_ resource.ResourceWithConfigure   = &verityBundleResource{}
	_ resource.ResourceWithImportState = &verityBundleResource{}
)

func NewVerityBundleResource() resource.Resource {
	return &verityBundleResource{}
}

type verityBundleResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
	notifyOperationAdded func()
}

type verityBundleResourceModel struct {
	Name                       types.String                        `tfsdk:"name"`
	Enable                     types.Bool                          `tfsdk:"enable"`
	DeviceSettings             types.String                        `tfsdk:"device_settings"`
	DeviceSettingsRefType      types.String                        `tfsdk:"device_settings_ref_type_"`
	CliCommands                types.String                        `tfsdk:"cli_commands"`
	Protocol                   types.String                        `tfsdk:"protocol"`
	DiagnosticsProfile         types.String                        `tfsdk:"diagnostics_profile"`
	DiagnosticsProfileRefType  types.String                        `tfsdk:"diagnostics_profile_ref_type_"`
	DeviceVoiceSettings        types.String                        `tfsdk:"device_voice_settings"`
	DeviceVoiceSettingsRefType types.String                        `tfsdk:"device_voice_settings_ref_type_"`
	ObjectProperties           []verityBundleObjectPropertiesModel `tfsdk:"object_properties"`
	EthPortPaths               []ethPortPathsModel                 `tfsdk:"eth_port_paths"`
	UserServices               []userServicesModel                 `tfsdk:"user_services"`
	VoicePortProfilePaths      []voicePortProfilePathsModel        `tfsdk:"voice_port_profile_paths"`
}

type verityBundleObjectPropertiesModel struct {
	IsForSwitch types.Bool   `tfsdk:"is_for_switch"`
	Group       types.String `tfsdk:"group"`
	IsPublic    types.Bool   `tfsdk:"is_public"`
}

type ethPortPathsModel struct {
	EthPortNumEthPortProfile                               types.String `tfsdk:"eth_port_num_eth_port_profile"`
	EthPortNumEthPortProfileRefType                        types.String `tfsdk:"eth_port_num_eth_port_profile_ref_type_"`
	EthPortNumEthPortSettings                              types.String `tfsdk:"eth_port_num_eth_port_settings"`
	EthPortNumEthPortSettingsRefType                       types.String `tfsdk:"eth_port_num_eth_port_settings_ref_type_"`
	EthPortNumGatewayProfile                               types.String `tfsdk:"eth_port_num_gateway_profile"`
	EthPortNumGatewayProfileRefType                        types.String `tfsdk:"eth_port_num_gateway_profile_ref_type_"`
	DiagnosticsPortProfileNumDiagnosticsPortProfile        types.String `tfsdk:"diagnostics_port_profile_num_diagnostics_port_profile"`
	DiagnosticsPortProfileNumDiagnosticsPortProfileRefType types.String `tfsdk:"diagnostics_port_profile_num_diagnostics_port_profile_ref_type_"`
	PortName                                               types.String `tfsdk:"port_name"`
	Index                                                  types.Int64  `tfsdk:"index"`
}

func (epp ethPortPathsModel) GetIndex() types.Int64 {
	return epp.Index
}

type voicePortProfilePathsModel struct {
	VoicePortNumVoicePortProfiles        types.String `tfsdk:"voice_port_num_voice_port_profiles"`
	VoicePortNumVoicePortProfilesRefType types.String `tfsdk:"voice_port_num_voice_port_profiles_ref_type_"`
	Index                                types.Int64  `tfsdk:"index"`
}

func (vppp voicePortProfilePathsModel) GetIndex() types.Int64 {
	return vppp.Index
}

type userServicesModel struct {
	RowAppEnable                  types.Bool   `tfsdk:"row_app_enable"`
	RowAppConnectedService        types.String `tfsdk:"row_app_connected_service"`
	RowAppConnectedServiceRefType types.String `tfsdk:"row_app_connected_service_ref_type_"`
	RowAppCliCommands             types.String `tfsdk:"row_app_cli_commands"`
	RowIpMask                     types.String `tfsdk:"row_ip_mask"`
	Index                         types.Int64  `tfsdk:"index"`
}

func (us userServicesModel) GetIndex() types.Int64 {
	return us.Index
}

func (r *verityBundleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bundle"
}

func (r *verityBundleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityBundleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity Bundle",
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "Object Name. Must be unique.",
				Required:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"enable": schema.BoolAttribute{
				Description: "Enable object. It's highly recommended to set this value to true so that validation on the object will be ran.",
				Optional:    true,
			},
			"device_settings": schema.StringAttribute{
				Description: "Device Settings for device",
				Optional:    true,
			},
			"device_settings_ref_type_": schema.StringAttribute{
				Description: "Object type for device_settings field",
				Optional:    true,
			},
			"cli_commands": schema.StringAttribute{
				Description: "CLI Commands",
				Optional:    true,
			},
			"protocol": schema.StringAttribute{
				Description: "Voice Protocol: MGCP or SIP",
				Optional:    true,
			},
			"diagnostics_profile": schema.StringAttribute{
				Description: "Diagnostics Profile for device",
				Optional:    true,
			},
			"diagnostics_profile_ref_type_": schema.StringAttribute{
				Description: "Object type for diagnostics_profile field",
				Optional:    true,
			},
			"device_voice_settings": schema.StringAttribute{
				Description: "Device Voice Settings for device",
				Optional:    true,
			},
			"device_voice_settings_ref_type_": schema.StringAttribute{
				Description: "Object type for device_voice_settings field",
				Optional:    true,
			},
		},
		Blocks: map[string]schema.Block{
			"object_properties": schema.ListNestedBlock{
				Description: "Object properties for the bundle",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"is_for_switch": schema.BoolAttribute{
							Description: "Denotes a Switch Bundle",
							Optional:    true,
						},
						"group": schema.StringAttribute{
							Description: "Group",
							Optional:    true,
						},
						"is_public": schema.BoolAttribute{
							Description: "Denotes a shared Switch Bundle",
							Optional:    true,
						},
					},
				},
			},
			"eth_port_paths": schema.ListNestedBlock{
				Description: "List of ethernet port configurations",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"eth_port_num_eth_port_profile": schema.StringAttribute{
							Description: "Eth Port Profile or LAG for Eth Port",
							Optional:    true,
						},
						"eth_port_num_eth_port_profile_ref_type_": schema.StringAttribute{
							Description: "Object type for eth_port_num_eth_port_profile field",
							Optional:    true,
						},
						"eth_port_num_eth_port_settings": schema.StringAttribute{
							Description: "Choose an Eth Port Settings",
							Optional:    true,
						},
						"eth_port_num_eth_port_settings_ref_type_": schema.StringAttribute{
							Description: "Object type for eth_port_num_eth_port_settings field",
							Optional:    true,
						},
						"eth_port_num_gateway_profile": schema.StringAttribute{
							Description: "Gateway Profile or LAG for Eth Port",
							Optional:    true,
						},
						"eth_port_num_gateway_profile_ref_type_": schema.StringAttribute{
							Description: "Object type for eth_port_num_gateway_profile field",
							Optional:    true,
						},
						"diagnostics_port_profile_num_diagnostics_port_profile": schema.StringAttribute{
							Description: "Diagnostics Port Profile for port",
							Optional:    true,
						},
						"diagnostics_port_profile_num_diagnostics_port_profile_ref_type_": schema.StringAttribute{
							Description: "Object type for diagnostics_port_profile_num_diagnostics_port_profile field",
							Optional:    true,
						},
						"port_name": schema.StringAttribute{
							Description: "The name identifying the port",
							Optional:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object. Zero if you want to add an object to the list.",
							Optional:    true,
						},
					},
				},
			},
			"user_services": schema.ListNestedBlock{
				Description: "List of user services configurations",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"row_app_enable": schema.BoolAttribute{
							Description: "Enable of this User application",
							Optional:    true,
						},
						"row_app_connected_service": schema.StringAttribute{
							Description: "Service connected to this User application",
							Optional:    true,
						},
						"row_app_connected_service_ref_type_": schema.StringAttribute{
							Description: "Object type for row_app_connected_service field",
							Optional:    true,
						},
						"row_app_cli_commands": schema.StringAttribute{
							Description: "CLI Commands of this User application",
							Optional:    true,
						},
						"row_ip_mask": schema.StringAttribute{
							Description: "IP/Mask in IPv4 format",
							Optional:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object. Zero if you want to add an object to the list.",
							Optional:    true,
						},
					},
				},
			},
			"voice_port_profile_paths": schema.ListNestedBlock{
				Description: "List of voice port profile configurations",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"voice_port_num_voice_port_profiles": schema.StringAttribute{
							Description: "Voice Port Profile for Voice Port",
							Optional:    true,
						},
						"voice_port_num_voice_port_profiles_ref_type_": schema.StringAttribute{
							Description: "Object type for voice_port_num_voice_port_profiles field",
							Optional:    true,
						},
						"index": schema.Int64Attribute{
							Description: "The index identifying the object. Zero if you want to add an object to the list.",
							Optional:    true,
						},
					},
				},
			},
		},
	}
}

func (r *verityBundleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityBundleResourceModel
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
	bundleProps := &openapi.BundlesPutRequestEndpointBundleValue{
		Name: openapi.PtrString(name),
	}

	// Handle string fields
	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "DeviceSettings", APIField: &bundleProps.DeviceSettings, TFValue: plan.DeviceSettings},
		{FieldName: "DeviceSettingsRefType", APIField: &bundleProps.DeviceSettingsRefType, TFValue: plan.DeviceSettingsRefType},
		{FieldName: "CliCommands", APIField: &bundleProps.CliCommands, TFValue: plan.CliCommands},
		{FieldName: "Protocol", APIField: &bundleProps.Protocol, TFValue: plan.Protocol},
		{FieldName: "DiagnosticsProfile", APIField: &bundleProps.DiagnosticsProfile, TFValue: plan.DiagnosticsProfile},
		{FieldName: "DiagnosticsProfileRefType", APIField: &bundleProps.DiagnosticsProfileRefType, TFValue: plan.DiagnosticsProfileRefType},
		{FieldName: "DeviceVoiceSettings", APIField: &bundleProps.DeviceVoiceSettings, TFValue: plan.DeviceVoiceSettings},
		{FieldName: "DeviceVoiceSettingsRefType", APIField: &bundleProps.DeviceVoiceSettingsRefType, TFValue: plan.DeviceVoiceSettingsRefType},
	})

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &bundleProps.Enable, TFValue: plan.Enable},
	})

	// Handle object properties
	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.BundlesPutRequestEndpointBundleValueObjectProperties{}
		utils.SetObjectPropertiesFields([]utils.ObjectPropertiesField{
			{Name: "IsForSwitch", TFValue: op.IsForSwitch, APIValue: &objProps.IsForSwitch},
			{Name: "Group", TFValue: op.Group, APIValue: &objProps.Group},
			{Name: "IsPublic", TFValue: op.IsPublic, APIValue: &objProps.IsPublic},
		})
		bundleProps.ObjectProperties = &objProps
	}

	// Handle eth port paths
	if len(plan.EthPortPaths) > 0 {
		ethPortPaths := make([]openapi.BundlesPutRequestEndpointBundleValueEthPortPathsInner, len(plan.EthPortPaths))
		for i, item := range plan.EthPortPaths {
			pathItem := openapi.BundlesPutRequestEndpointBundleValueEthPortPathsInner{}
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "EthPortNumEthPortProfile", APIField: &pathItem.EthPortNumEthPortProfile, TFValue: item.EthPortNumEthPortProfile},
				{FieldName: "EthPortNumEthPortProfileRefType", APIField: &pathItem.EthPortNumEthPortProfileRefType, TFValue: item.EthPortNumEthPortProfileRefType},
				{FieldName: "EthPortNumEthPortSettings", APIField: &pathItem.EthPortNumEthPortSettings, TFValue: item.EthPortNumEthPortSettings},
				{FieldName: "EthPortNumEthPortSettingsRefType", APIField: &pathItem.EthPortNumEthPortSettingsRefType, TFValue: item.EthPortNumEthPortSettingsRefType},
				{FieldName: "EthPortNumGatewayProfile", APIField: &pathItem.EthPortNumGatewayProfile, TFValue: item.EthPortNumGatewayProfile},
				{FieldName: "EthPortNumGatewayProfileRefType", APIField: &pathItem.EthPortNumGatewayProfileRefType, TFValue: item.EthPortNumGatewayProfileRefType},
				{FieldName: "DiagnosticsPortProfileNumDiagnosticsPortProfile", APIField: &pathItem.DiagnosticsPortProfileNumDiagnosticsPortProfile, TFValue: item.DiagnosticsPortProfileNumDiagnosticsPortProfile},
				{FieldName: "DiagnosticsPortProfileNumDiagnosticsPortProfileRefType", APIField: &pathItem.DiagnosticsPortProfileNumDiagnosticsPortProfileRefType, TFValue: item.DiagnosticsPortProfileNumDiagnosticsPortProfileRefType},
				{FieldName: "PortName", APIField: &pathItem.PortName, TFValue: item.PortName},
			})
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &pathItem.Index, TFValue: item.Index},
			})
			ethPortPaths[i] = pathItem
		}
		bundleProps.EthPortPaths = ethPortPaths
	}

	// Handle user services
	if len(plan.UserServices) > 0 {
		userServices := make([]openapi.BundlesPutRequestEndpointBundleValueUserServicesInner, len(plan.UserServices))
		for i, item := range plan.UserServices {
			serviceItem := openapi.BundlesPutRequestEndpointBundleValueUserServicesInner{}
			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "RowAppEnable", APIField: &serviceItem.RowAppEnable, TFValue: item.RowAppEnable},
			})
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "RowAppConnectedService", APIField: &serviceItem.RowAppConnectedService, TFValue: item.RowAppConnectedService},
				{FieldName: "RowAppConnectedServiceRefType", APIField: &serviceItem.RowAppConnectedServiceRefType, TFValue: item.RowAppConnectedServiceRefType},
				{FieldName: "RowAppCliCommands", APIField: &serviceItem.RowAppCliCommands, TFValue: item.RowAppCliCommands},
				{FieldName: "RowIpMask", APIField: &serviceItem.RowIpMask, TFValue: item.RowIpMask},
			})
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &serviceItem.Index, TFValue: item.Index},
			})
			userServices[i] = serviceItem
		}
		bundleProps.UserServices = userServices
	}

	// Handle voice port profile paths
	if len(plan.VoicePortProfilePaths) > 0 {
		voicePortProfilePaths := make([]openapi.BundlesPutRequestEndpointBundleValueVoicePortProfilePathsInner, len(plan.VoicePortProfilePaths))
		for i, item := range plan.VoicePortProfilePaths {
			pathItem := openapi.BundlesPutRequestEndpointBundleValueVoicePortProfilePathsInner{}
			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "VoicePortNumVoicePortProfiles", APIField: &pathItem.VoicePortNumVoicePortProfiles, TFValue: item.VoicePortNumVoicePortProfiles},
				{FieldName: "VoicePortNumVoicePortProfilesRefType", APIField: &pathItem.VoicePortNumVoicePortProfilesRefType, TFValue: item.VoicePortNumVoicePortProfilesRefType},
			})
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &pathItem.Index, TFValue: item.Index},
			})
			voicePortProfilePaths[i] = pathItem
		}
		bundleProps.VoicePortProfilePaths = voicePortProfilePaths
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "bundle", name, *bundleProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Bundle %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "bundles")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
}

func (r *verityBundleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityBundleResourceModel
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

	bundleName := state.Name.ValueString()

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("bundle") {
		tflog.Info(ctx, fmt.Sprintf("Skipping bundle %s verification – trusting recent successful API operation", bundleName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching bundles for verification of %s", bundleName))

	type BundleResponse struct {
		EndpointBundle map[string]interface{} `json:"endpoint_bundle"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "bundles", bundleName,
		func() (BundleResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch bundles")
			resp, err := r.client.BundlesAPI.BundlesGet(ctx).Execute()
			if err != nil {
				return BundleResponse{}, fmt.Errorf("error reading bundle: %v", err)
			}
			defer resp.Body.Close()

			var result BundleResponse
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				return BundleResponse{}, fmt.Errorf("failed to decode bundles response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d bundles from API", len(result.EndpointBundle)))
			return result, nil
		}, getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Bundle %s", bundleName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for bundle with name: %s", bundleName))

	bundleData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.EndpointBundle,
		bundleName,
		func(data interface{}) (string, bool) {
			if bundle, ok := data.(map[string]interface{}); ok {
				if name, ok := bundle["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Bundle with name '%s' not found in API response", bundleName))
		resp.State.RemoveResource(ctx)
		return
	}

	bundleMap, ok := bundleData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Bundle Data",
			fmt.Sprintf("Bundle data is not in expected format for %s", bundleName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found bundle '%s' under API key '%s'", bundleName, actualAPIName))

	state.Name = utils.MapStringFromAPI(bundleMap["name"])

	// Handle object properties
	if objProps, ok := bundleMap["object_properties"].(map[string]interface{}); ok {
		state.ObjectProperties = []verityBundleObjectPropertiesModel{
			{
				IsForSwitch: utils.MapBoolFromAPI(objProps["is_for_switch"]),
				Group:       utils.MapStringFromAPI(objProps["group"]),
				IsPublic:    utils.MapBoolFromAPI(objProps["is_public"]),
			},
		}
	} else {
		state.ObjectProperties = nil
	}

	// Map string fields
	stringFieldMappings := map[string]*types.String{
		"device_settings":                 &state.DeviceSettings,
		"device_settings_ref_type_":       &state.DeviceSettingsRefType,
		"cli_commands":                    &state.CliCommands,
		"protocol":                        &state.Protocol,
		"diagnostics_profile":             &state.DiagnosticsProfile,
		"diagnostics_profile_ref_type_":   &state.DiagnosticsProfileRefType,
		"device_voice_settings":           &state.DeviceVoiceSettings,
		"device_voice_settings_ref_type_": &state.DeviceVoiceSettingsRefType,
	}

	for apiKey, stateField := range stringFieldMappings {
		*stateField = utils.MapStringFromAPI(bundleMap[apiKey])
	}

	// Map boolean fields
	boolFieldMappings := map[string]*types.Bool{
		"enable": &state.Enable,
	}

	for apiKey, stateField := range boolFieldMappings {
		*stateField = utils.MapBoolFromAPI(bundleMap[apiKey])
	}

	// Handle eth port paths
	var ethPortPaths []ethPortPathsModel
	if paths, ok := bundleMap["eth_port_paths"].([]interface{}); ok && len(paths) > 0 {
		for _, p := range paths {
			path, ok := p.(map[string]interface{})
			if !ok {
				continue
			}

			ethPortPath := ethPortPathsModel{
				EthPortNumEthPortProfile:                               utils.MapStringFromAPI(path["eth_port_num_eth_port_profile"]),
				EthPortNumEthPortSettings:                              utils.MapStringFromAPI(path["eth_port_num_eth_port_settings"]),
				EthPortNumEthPortSettingsRefType:                       utils.MapStringFromAPI(path["eth_port_num_eth_port_settings_ref_type_"]),
				EthPortNumEthPortProfileRefType:                        utils.MapStringFromAPI(path["eth_port_num_eth_port_profile_ref_type_"]),
				DiagnosticsPortProfileNumDiagnosticsPortProfile:        utils.MapStringFromAPI(path["diagnostics_port_profile_num_diagnostics_port_profile"]),
				DiagnosticsPortProfileNumDiagnosticsPortProfileRefType: utils.MapStringFromAPI(path["diagnostics_port_profile_num_diagnostics_port_profile_ref_type_"]),
				PortName:                        utils.MapStringFromAPI(path["port_name"]),
				EthPortNumGatewayProfile:        utils.MapStringFromAPI(path["eth_port_num_gateway_profile"]),
				EthPortNumGatewayProfileRefType: utils.MapStringFromAPI(path["eth_port_num_gateway_profile_ref_type_"]),
				Index:                           utils.MapInt64FromAPI(path["index"]),
			}

			ethPortPaths = append(ethPortPaths, ethPortPath)
		}
		state.EthPortPaths = ethPortPaths
	} else {
		state.EthPortPaths = nil
	}

	// Handle user services
	var userServices []userServicesModel
	if services, ok := bundleMap["user_services"].([]interface{}); ok && len(services) > 0 {
		for _, s := range services {
			service, ok := s.(map[string]interface{})
			if !ok {
				continue
			}

			userService := userServicesModel{
				RowAppEnable:                  utils.MapBoolFromAPI(service["row_app_enable"]),
				RowAppConnectedService:        utils.MapStringFromAPI(service["row_app_connected_service"]),
				RowAppCliCommands:             utils.MapStringFromAPI(service["row_app_cli_commands"]),
				RowIpMask:                     utils.MapStringFromAPI(service["row_ip_mask"]),
				RowAppConnectedServiceRefType: utils.MapStringFromAPI(service["row_app_connected_service_ref_type_"]),
				Index:                         utils.MapInt64FromAPI(service["index"]),
			}

			userServices = append(userServices, userService)
		}
		state.UserServices = userServices
	} else {
		state.UserServices = nil
	}

	// Handle voice port profile paths
	var voicePortProfilePaths []voicePortProfilePathsModel
	if paths, ok := bundleMap["voice_port_profile_paths"].([]interface{}); ok && len(paths) > 0 {
		for _, p := range paths {
			path, ok := p.(map[string]interface{})
			if !ok {
				continue
			}

			voicePortPath := voicePortProfilePathsModel{
				VoicePortNumVoicePortProfiles:        utils.MapStringFromAPI(path["voice_port_num_voice_port_profiles"]),
				VoicePortNumVoicePortProfilesRefType: utils.MapStringFromAPI(path["voice_port_num_voice_port_profiles_ref_type_"]),
				Index:                                utils.MapInt64FromAPI(path["index"]),
			}

			voicePortProfilePaths = append(voicePortProfilePaths, voicePortPath)
		}
		state.VoicePortProfilePaths = voicePortProfilePaths
	} else {
		state.VoicePortProfilePaths = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityBundleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityBundleResourceModel

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
	bundleProps := openapi.BundlesPutRequestEndpointBundleValue{}
	hasChanges := false

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { bundleProps.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.CliCommands, state.CliCommands, func(v *string) { bundleProps.CliCommands = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Protocol, state.Protocol, func(v *string) { bundleProps.Protocol = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { bundleProps.Enable = v }, &hasChanges)

	// Handle object properties
	if len(plan.ObjectProperties) > 0 && len(state.ObjectProperties) > 0 {
		objProps := openapi.BundlesPutRequestEndpointBundleValueObjectProperties{}
		op := plan.ObjectProperties[0]
		st := state.ObjectProperties[0]
		objPropsChanged := false

		utils.CompareAndSetObjectPropertiesFields([]utils.ObjectPropertiesFieldWithComparison{
			{Name: "IsForSwitch", PlanValue: op.IsForSwitch, StateValue: st.IsForSwitch, APIValue: &objProps.IsForSwitch},
			{Name: "Group", PlanValue: op.Group, StateValue: st.Group, APIValue: &objProps.Group},
			{Name: "IsPublic", PlanValue: op.IsPublic, StateValue: st.IsPublic, APIValue: &objProps.IsPublic},
		}, &objPropsChanged)

		if objPropsChanged {
			bundleProps.ObjectProperties = &objProps
			hasChanges = true
		}
	}

	// Handle device settings reference type using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.DeviceSettings, state.DeviceSettings, plan.DeviceSettingsRefType, state.DeviceSettingsRefType,
		func(v *string) { bundleProps.DeviceSettings = v },
		func(v *string) { bundleProps.DeviceSettingsRefType = v },
		"device_settings", "device_settings_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	// Handle diagnostics profile reference type using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.DiagnosticsProfile, state.DiagnosticsProfile, plan.DiagnosticsProfileRefType, state.DiagnosticsProfileRefType,
		func(v *string) { bundleProps.DiagnosticsProfile = v },
		func(v *string) { bundleProps.DiagnosticsProfileRefType = v },
		"diagnostics_profile", "diagnostics_profile_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	// Handle device voice settings reference type using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.DeviceVoiceSettings, state.DeviceVoiceSettings, plan.DeviceVoiceSettingsRefType, state.DeviceVoiceSettingsRefType,
		func(v *string) { bundleProps.DeviceVoiceSettings = v },
		func(v *string) { bundleProps.DeviceVoiceSettingsRefType = v },
		"device_voice_settings", "device_voice_settings_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	// Handle eth port paths
	ethPortPathsHandler := utils.IndexedItemHandler[ethPortPathsModel, openapi.BundlesPutRequestEndpointBundleValueEthPortPathsInner]{
		CreateNew: func(planItem ethPortPathsModel) openapi.BundlesPutRequestEndpointBundleValueEthPortPathsInner {
			ethPortPath := openapi.BundlesPutRequestEndpointBundleValueEthPortPathsInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &ethPortPath.Index, TFValue: planItem.Index},
			})

			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "PortName", APIField: &ethPortPath.PortName, TFValue: planItem.PortName},
				{FieldName: "EthPortNumEthPortSettings", APIField: &ethPortPath.EthPortNumEthPortSettings, TFValue: planItem.EthPortNumEthPortSettings},
				{FieldName: "EthPortNumEthPortSettingsRefType", APIField: &ethPortPath.EthPortNumEthPortSettingsRefType, TFValue: planItem.EthPortNumEthPortSettingsRefType},
				{FieldName: "EthPortNumEthPortProfile", APIField: &ethPortPath.EthPortNumEthPortProfile, TFValue: planItem.EthPortNumEthPortProfile},
				{FieldName: "EthPortNumEthPortProfileRefType", APIField: &ethPortPath.EthPortNumEthPortProfileRefType, TFValue: planItem.EthPortNumEthPortProfileRefType},
				{FieldName: "EthPortNumGatewayProfile", APIField: &ethPortPath.EthPortNumGatewayProfile, TFValue: planItem.EthPortNumGatewayProfile},
				{FieldName: "EthPortNumGatewayProfileRefType", APIField: &ethPortPath.EthPortNumGatewayProfileRefType, TFValue: planItem.EthPortNumGatewayProfileRefType},
				{FieldName: "DiagnosticsPortProfileNumDiagnosticsPortProfile", APIField: &ethPortPath.DiagnosticsPortProfileNumDiagnosticsPortProfile, TFValue: planItem.DiagnosticsPortProfileNumDiagnosticsPortProfile},
				{FieldName: "DiagnosticsPortProfileNumDiagnosticsPortProfileRefType", APIField: &ethPortPath.DiagnosticsPortProfileNumDiagnosticsPortProfileRefType, TFValue: planItem.DiagnosticsPortProfileNumDiagnosticsPortProfileRefType},
			})

			return ethPortPath
		},
		UpdateExisting: func(planItem ethPortPathsModel, stateItem ethPortPathsModel) (openapi.BundlesPutRequestEndpointBundleValueEthPortPathsInner, bool) {
			ethPortPath := openapi.BundlesPutRequestEndpointBundleValueEthPortPathsInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &ethPortPath.Index, TFValue: planItem.Index},
			})

			fieldChanged := false

			// Handle non-ref-type string fields
			utils.CompareAndSetStringField(planItem.PortName, stateItem.PortName, func(v *string) { ethPortPath.PortName = v }, &fieldChanged)

			// Handle eth_port_num_eth_port_settings and eth_port_num_eth_port_settings_ref_type_ using "One ref type supported" pattern
			if !utils.HandleOneRefTypeSupported(
				planItem.EthPortNumEthPortSettings, stateItem.EthPortNumEthPortSettings, planItem.EthPortNumEthPortSettingsRefType, stateItem.EthPortNumEthPortSettingsRefType,
				func(v *string) { ethPortPath.EthPortNumEthPortSettings = v },
				func(v *string) { ethPortPath.EthPortNumEthPortSettingsRefType = v },
				"eth_port_num_eth_port_settings", "eth_port_num_eth_port_settings_ref_type_",
				&fieldChanged,
				&resp.Diagnostics,
			) {
				return ethPortPath, false
			}

			// Handle eth_port_num_eth_port_profile and eth_port_num_eth_port_profile_ref_type_ using "Many ref types supported" pattern
			if !utils.HandleMultipleRefTypesSupported(
				planItem.EthPortNumEthPortProfile, stateItem.EthPortNumEthPortProfile, planItem.EthPortNumEthPortProfileRefType, stateItem.EthPortNumEthPortProfileRefType,
				func(v *string) { ethPortPath.EthPortNumEthPortProfile = v },
				func(v *string) { ethPortPath.EthPortNumEthPortProfileRefType = v },
				"eth_port_num_eth_port_profile", "eth_port_num_eth_port_profile_ref_type_",
				&fieldChanged,
				&resp.Diagnostics,
			) {
				return ethPortPath, false
			}

			// Handle eth_port_num_gateway_profile and eth_port_num_gateway_profile_ref_type_ using "Many ref types supported" pattern
			if !utils.HandleMultipleRefTypesSupported(
				planItem.EthPortNumGatewayProfile, stateItem.EthPortNumGatewayProfile, planItem.EthPortNumGatewayProfileRefType, stateItem.EthPortNumGatewayProfileRefType,
				func(v *string) { ethPortPath.EthPortNumGatewayProfile = v },
				func(v *string) { ethPortPath.EthPortNumGatewayProfileRefType = v },
				"eth_port_num_gateway_profile", "eth_port_num_gateway_profile_ref_type_",
				&fieldChanged,
				&resp.Diagnostics,
			) {
				return ethPortPath, false
			}

			// Handle diagnostics_port_profile_num_diagnostics_port_profile and diagnostics_port_profile_num_diagnostics_port_profile_ref_type_ using "One ref type supported" pattern
			if !utils.HandleOneRefTypeSupported(
				planItem.DiagnosticsPortProfileNumDiagnosticsPortProfile, stateItem.DiagnosticsPortProfileNumDiagnosticsPortProfile, planItem.DiagnosticsPortProfileNumDiagnosticsPortProfileRefType, stateItem.DiagnosticsPortProfileNumDiagnosticsPortProfileRefType,
				func(v *string) { ethPortPath.DiagnosticsPortProfileNumDiagnosticsPortProfile = v },
				func(v *string) { ethPortPath.DiagnosticsPortProfileNumDiagnosticsPortProfileRefType = v },
				"diagnostics_port_profile_num_diagnostics_port_profile", "diagnostics_port_profile_num_diagnostics_port_profile_ref_type_",
				&fieldChanged,
				&resp.Diagnostics,
			) {
				return ethPortPath, false
			}

			return ethPortPath, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.BundlesPutRequestEndpointBundleValueEthPortPathsInner {
			item := openapi.BundlesPutRequestEndpointBundleValueEthPortPathsInner{}
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &item.Index, TFValue: types.Int64Value(index)},
			})
			return item
		},
	}

	changedEthPortPaths, ethPortPathsChanged := utils.ProcessIndexedArrayUpdates(plan.EthPortPaths, state.EthPortPaths, ethPortPathsHandler)
	if ethPortPathsChanged {
		bundleProps.EthPortPaths = changedEthPortPaths
		hasChanges = true
	}

	// Handle user services
	userServicesHandler := utils.IndexedItemHandler[userServicesModel, openapi.BundlesPutRequestEndpointBundleValueUserServicesInner]{
		CreateNew: func(planItem userServicesModel) openapi.BundlesPutRequestEndpointBundleValueUserServicesInner {
			userService := openapi.BundlesPutRequestEndpointBundleValueUserServicesInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &userService.Index, TFValue: planItem.Index},
			})

			utils.SetBoolFields([]utils.BoolFieldMapping{
				{FieldName: "RowAppEnable", APIField: &userService.RowAppEnable, TFValue: planItem.RowAppEnable},
			})

			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "RowAppConnectedService", APIField: &userService.RowAppConnectedService, TFValue: planItem.RowAppConnectedService},
				{FieldName: "RowAppConnectedServiceRefType", APIField: &userService.RowAppConnectedServiceRefType, TFValue: planItem.RowAppConnectedServiceRefType},
				{FieldName: "RowAppCliCommands", APIField: &userService.RowAppCliCommands, TFValue: planItem.RowAppCliCommands},
				{FieldName: "RowIpMask", APIField: &userService.RowIpMask, TFValue: planItem.RowIpMask},
			})

			return userService
		},
		UpdateExisting: func(planItem userServicesModel, stateItem userServicesModel) (openapi.BundlesPutRequestEndpointBundleValueUserServicesInner, bool) {
			userService := openapi.BundlesPutRequestEndpointBundleValueUserServicesInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &userService.Index, TFValue: planItem.Index},
			})

			fieldChanged := false

			// Handle boolean fields
			utils.CompareAndSetBoolField(planItem.RowAppEnable, stateItem.RowAppEnable, func(v *bool) { userService.RowAppEnable = v }, &fieldChanged)

			// Handle row_app_connected_service and row_app_connected_service_ref_type_ using "One ref type supported" pattern
			if !utils.HandleOneRefTypeSupported(
				planItem.RowAppConnectedService, stateItem.RowAppConnectedService, planItem.RowAppConnectedServiceRefType, stateItem.RowAppConnectedServiceRefType,
				func(v *string) { userService.RowAppConnectedService = v },
				func(v *string) { userService.RowAppConnectedServiceRefType = v },
				"row_app_connected_service", "row_app_connected_service_ref_type_",
				&fieldChanged,
				&resp.Diagnostics,
			) {
				return userService, false
			}

			// Handle non-ref-type string fields
			utils.CompareAndSetStringField(planItem.RowAppCliCommands, stateItem.RowAppCliCommands, func(v *string) { userService.RowAppCliCommands = v }, &fieldChanged)
			utils.CompareAndSetStringField(planItem.RowIpMask, stateItem.RowIpMask, func(v *string) { userService.RowIpMask = v }, &fieldChanged)

			return userService, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.BundlesPutRequestEndpointBundleValueUserServicesInner {
			item := openapi.BundlesPutRequestEndpointBundleValueUserServicesInner{}
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &item.Index, TFValue: types.Int64Value(index)},
			})
			return item
		},
	}

	changedUserServices, userServicesChanged := utils.ProcessIndexedArrayUpdates(plan.UserServices, state.UserServices, userServicesHandler)
	if userServicesChanged {
		bundleProps.UserServices = changedUserServices
		hasChanges = true
	}

	// Handle voice port profile paths
	voicePortProfilePathsHandler := utils.IndexedItemHandler[voicePortProfilePathsModel, openapi.BundlesPutRequestEndpointBundleValueVoicePortProfilePathsInner]{
		CreateNew: func(planItem voicePortProfilePathsModel) openapi.BundlesPutRequestEndpointBundleValueVoicePortProfilePathsInner {
			voicePortPath := openapi.BundlesPutRequestEndpointBundleValueVoicePortProfilePathsInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &voicePortPath.Index, TFValue: planItem.Index},
			})

			utils.SetStringFields([]utils.StringFieldMapping{
				{FieldName: "VoicePortNumVoicePortProfiles", APIField: &voicePortPath.VoicePortNumVoicePortProfiles, TFValue: planItem.VoicePortNumVoicePortProfiles},
				{FieldName: "VoicePortNumVoicePortProfilesRefType", APIField: &voicePortPath.VoicePortNumVoicePortProfilesRefType, TFValue: planItem.VoicePortNumVoicePortProfilesRefType},
			})

			return voicePortPath
		},
		UpdateExisting: func(planItem voicePortProfilePathsModel, stateItem voicePortProfilePathsModel) (openapi.BundlesPutRequestEndpointBundleValueVoicePortProfilePathsInner, bool) {
			voicePortPath := openapi.BundlesPutRequestEndpointBundleValueVoicePortProfilePathsInner{}

			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &voicePortPath.Index, TFValue: planItem.Index},
			})

			fieldChanged := false

			// Handle voice_port_num_voice_port_profiles and voice_port_num_voice_port_profiles_ref_type_ using "One ref type supported" pattern
			if !utils.HandleOneRefTypeSupported(
				planItem.VoicePortNumVoicePortProfiles, stateItem.VoicePortNumVoicePortProfiles, planItem.VoicePortNumVoicePortProfilesRefType, stateItem.VoicePortNumVoicePortProfilesRefType,
				func(v *string) { voicePortPath.VoicePortNumVoicePortProfiles = v },
				func(v *string) { voicePortPath.VoicePortNumVoicePortProfilesRefType = v },
				"voice_port_num_voice_port_profiles", "voice_port_num_voice_port_profiles_ref_type_",
				&fieldChanged,
				&resp.Diagnostics,
			) {
				return voicePortPath, false
			}

			return voicePortPath, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.BundlesPutRequestEndpointBundleValueVoicePortProfilePathsInner {
			item := openapi.BundlesPutRequestEndpointBundleValueVoicePortProfilePathsInner{}
			utils.SetInt64Fields([]utils.Int64FieldMapping{
				{FieldName: "Index", APIField: &item.Index, TFValue: types.Int64Value(index)},
			})
			return item
		},
	}

	changedVoicePortProfilePaths, voicePortProfilePathsChanged := utils.ProcessIndexedArrayUpdates(plan.VoicePortProfilePaths, state.VoicePortProfilePaths, voicePortProfilePathsHandler)
	if voicePortProfilePathsChanged {
		bundleProps.VoicePortProfilePaths = changedVoicePortProfilePaths
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "bundle", name, bundleProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Bundle %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "bundles")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *verityBundleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityBundleResourceModel
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "bundle", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Bundle %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "bundles")
	resp.State.RemoveResource(ctx)
}

func (r *verityBundleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
