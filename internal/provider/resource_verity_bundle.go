package provider

import (
	"context"
	"encoding/json"
	"fmt"
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
	bulkOpsMgr           *utils.BulkOperationManager
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

type voicePortProfilePathsModel struct {
	VoicePortNumVoicePortProfiles        types.String `tfsdk:"voice_port_num_voice_port_profiles"`
	VoicePortNumVoicePortProfilesRefType types.String `tfsdk:"voice_port_num_voice_port_profiles_ref_type_"`
	Index                                types.Int64  `tfsdk:"index"`
}

type userServicesModel struct {
	RowAppEnable                  types.Bool   `tfsdk:"row_app_enable"`
	RowAppConnectedService        types.String `tfsdk:"row_app_connected_service"`
	RowAppConnectedServiceRefType types.String `tfsdk:"row_app_connected_service_ref_type_"`
	RowAppCliCommands             types.String `tfsdk:"row_app_cli_commands"`
	RowIpMask                     types.String `tfsdk:"row_ip_mask"`
	Index                         types.Int64  `tfsdk:"index"`
}

func (r *verityBundleResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bundle"
}

func (r *verityBundleResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityBundleResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// API version check: Only allow on 6.5+
	apiVersion, err := getApiVersion(ctx, r.provCtx)
	if err != nil {
		resp.Diagnostics.AddError("API Version Error", fmt.Sprintf("Unable to determine API version: %s", err))
		return
	}
	if apiVersion < "6.5" {
		resp.Diagnostics.AddError("Unsupported API Version", "Bundle resource creation is only supported on API version 6.5 and above.")
		return
	}

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError("Authentication Error", fmt.Sprintf("Unable to authenticate client: %s", err))
		return
	}

	name := plan.Name.ValueString()
	bundleProps := &openapi.BundlesPutRequestEndpointBundleValue{
		Name: openapi.PtrString(name),
	}

	if !plan.Enable.IsNull() {
		bundleProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
	}
	if !plan.DeviceSettings.IsNull() {
		bundleProps.DeviceSettings = openapi.PtrString(plan.DeviceSettings.ValueString())
	}
	if !plan.DeviceSettingsRefType.IsNull() {
		bundleProps.DeviceSettingsRefType = openapi.PtrString(plan.DeviceSettingsRefType.ValueString())
	}
	if !plan.CliCommands.IsNull() {
		bundleProps.CliCommands = openapi.PtrString(plan.CliCommands.ValueString())
	}
	if !plan.Protocol.IsNull() {
		bundleProps.Protocol = openapi.PtrString(plan.Protocol.ValueString())
	}
	if !plan.DiagnosticsProfile.IsNull() {
		bundleProps.DiagnosticsProfile = openapi.PtrString(plan.DiagnosticsProfile.ValueString())
	}
	if !plan.DiagnosticsProfileRefType.IsNull() {
		bundleProps.DiagnosticsProfileRefType = openapi.PtrString(plan.DiagnosticsProfileRefType.ValueString())
	}
	if !plan.DeviceVoiceSettings.IsNull() {
		bundleProps.DeviceVoiceSettings = openapi.PtrString(plan.DeviceVoiceSettings.ValueString())
	}
	if !plan.DeviceVoiceSettingsRefType.IsNull() {
		bundleProps.DeviceVoiceSettingsRefType = openapi.PtrString(plan.DeviceVoiceSettingsRefType.ValueString())
	}

	if len(plan.ObjectProperties) > 0 {
		op := plan.ObjectProperties[0]
		objProps := openapi.BundlesPutRequestEndpointBundleValueObjectProperties{}
		if !op.IsForSwitch.IsNull() {
			objProps.IsForSwitch = openapi.PtrBool(op.IsForSwitch.ValueBool())
		}
		if !op.Group.IsNull() {
			objProps.Group = openapi.PtrString(op.Group.ValueString())
		}
		if !op.IsPublic.IsNull() {
			objProps.IsPublic = openapi.PtrBool(op.IsPublic.ValueBool())
		}
		bundleProps.ObjectProperties = &objProps
	}

	if len(plan.EthPortPaths) > 0 {
		ethPortPaths := make([]openapi.BundlesPutRequestEndpointBundleValueEthPortPathsInner, len(plan.EthPortPaths))
		for i, path := range plan.EthPortPaths {
			pathItem := openapi.BundlesPutRequestEndpointBundleValueEthPortPathsInner{}
			if !path.EthPortNumEthPortProfile.IsNull() {
				pathItem.EthPortNumEthPortProfile = openapi.PtrString(path.EthPortNumEthPortProfile.ValueString())
			}
			if !path.EthPortNumEthPortProfileRefType.IsNull() {
				pathItem.EthPortNumEthPortProfileRefType = openapi.PtrString(path.EthPortNumEthPortProfileRefType.ValueString())
			}
			if !path.EthPortNumEthPortSettings.IsNull() {
				pathItem.EthPortNumEthPortSettings = openapi.PtrString(path.EthPortNumEthPortSettings.ValueString())
			}
			if !path.EthPortNumEthPortSettingsRefType.IsNull() {
				pathItem.EthPortNumEthPortSettingsRefType = openapi.PtrString(path.EthPortNumEthPortSettingsRefType.ValueString())
			}
			if !path.EthPortNumGatewayProfile.IsNull() {
				pathItem.EthPortNumGatewayProfile = openapi.PtrString(path.EthPortNumGatewayProfile.ValueString())
			}
			if !path.EthPortNumGatewayProfileRefType.IsNull() {
				pathItem.EthPortNumGatewayProfileRefType = openapi.PtrString(path.EthPortNumGatewayProfileRefType.ValueString())
			}
			if !path.DiagnosticsPortProfileNumDiagnosticsPortProfile.IsNull() {
				pathItem.DiagnosticsPortProfileNumDiagnosticsPortProfile = openapi.PtrString(path.DiagnosticsPortProfileNumDiagnosticsPortProfile.ValueString())
			}
			if !path.DiagnosticsPortProfileNumDiagnosticsPortProfileRefType.IsNull() {
				pathItem.DiagnosticsPortProfileNumDiagnosticsPortProfileRefType = openapi.PtrString(path.DiagnosticsPortProfileNumDiagnosticsPortProfileRefType.ValueString())
			}
			if !path.PortName.IsNull() {
				pathItem.PortName = openapi.PtrString(path.PortName.ValueString())
			}
			if !path.Index.IsNull() {
				pathItem.Index = openapi.PtrInt32(int32(path.Index.ValueInt64()))
			}
			ethPortPaths[i] = pathItem
		}
		bundleProps.EthPortPaths = ethPortPaths
	}

	if len(plan.UserServices) > 0 {
		userServices := make([]openapi.BundlesPutRequestEndpointBundleValueUserServicesInner, len(plan.UserServices))
		for i, service := range plan.UserServices {
			serviceItem := openapi.BundlesPutRequestEndpointBundleValueUserServicesInner{}
			if !service.RowAppEnable.IsNull() {
				serviceItem.RowAppEnable = openapi.PtrBool(service.RowAppEnable.ValueBool())
			}
			if !service.RowAppConnectedService.IsNull() {
				serviceItem.RowAppConnectedService = openapi.PtrString(service.RowAppConnectedService.ValueString())
			}
			if !service.RowAppConnectedServiceRefType.IsNull() {
				serviceItem.RowAppConnectedServiceRefType = openapi.PtrString(service.RowAppConnectedServiceRefType.ValueString())
			}
			if !service.RowAppCliCommands.IsNull() {
				serviceItem.RowAppCliCommands = openapi.PtrString(service.RowAppCliCommands.ValueString())
			}
			if !service.RowIpMask.IsNull() {
				serviceItem.RowIpMask = openapi.PtrString(service.RowIpMask.ValueString())
			}
			if !service.Index.IsNull() {
				serviceItem.Index = openapi.PtrInt32(int32(service.Index.ValueInt64()))
			}
			userServices[i] = serviceItem
		}
		bundleProps.UserServices = userServices
	}

	if len(plan.VoicePortProfilePaths) > 0 {
		voicePortProfilePaths := make([]openapi.BundlesPutRequestEndpointBundleValueVoicePortProfilePathsInner, len(plan.VoicePortProfilePaths))
		for i, path := range plan.VoicePortProfilePaths {
			pathItem := openapi.BundlesPutRequestEndpointBundleValueVoicePortProfilePathsInner{}
			if !path.VoicePortNumVoicePortProfiles.IsNull() {
				pathItem.VoicePortNumVoicePortProfiles = openapi.PtrString(path.VoicePortNumVoicePortProfiles.ValueString())
			}
			if !path.VoicePortNumVoicePortProfilesRefType.IsNull() {
				pathItem.VoicePortNumVoicePortProfilesRefType = openapi.PtrString(path.VoicePortNumVoicePortProfilesRefType.ValueString())
			}
			if !path.Index.IsNull() {
				pathItem.Index = openapi.PtrInt32(int32(path.Index.ValueInt64()))
			}
			voicePortProfilePaths[i] = pathItem
		}
		bundleProps.VoicePortProfilePaths = voicePortProfilePaths
	}

	operationID := r.bulkOpsMgr.AddPut(ctx, "bundle", name, *bundleProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for bundle creation operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Create Bundle %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Bundle %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "bundles")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
}

func (r *verityBundleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data verityBundleResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError("Authentication Error", fmt.Sprintf("Unable to authenticate client: %s", err))
		return
	}

	bundleName := data.Name.ValueString()

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("bundle") {
		tflog.Info(ctx, fmt.Sprintf("Skipping bundle %s verification - trusting recent successful API operation", bundleName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("No recent bundle operations found, performing normal verification for %s", bundleName))

	type BundleResponse struct {
		EndpointBundle map[string]map[string]interface{} `json:"endpoint_bundle"`
	}

	var result BundleResponse
	var err error
	maxRetries := 3

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch bundles on attempt %d, retrying in %v", attempt, sleepTime))
			time.Sleep(sleepTime)
		}

		bundlesData, fetchErr := getCachedResponse(ctx, r.provCtx, "bundles", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch bundles")
			resp, err := r.client.BundlesAPI.BundlesGet(ctx).Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading bundle: %v", err)
			}
			defer resp.Body.Close()

			var result BundleResponse
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				return nil, fmt.Errorf("failed to decode bundles response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d bundles from API", len(result.EndpointBundle)))
			return result, nil
		})

		if fetchErr == nil {
			result = bundlesData.(BundleResponse)
			break
		}
		err = fetchErr
	}

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Bundle %s", bundleName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for bundle with ID: %s", bundleName))

	bundleData, exists := result.EndpointBundle[bundleName]
	if exists {
		tflog.Debug(ctx, fmt.Sprintf("Found bundle directly by ID: %s", bundleName))
	}

	if !exists {
		for apiName, b := range result.EndpointBundle {
			if name, ok := b["name"].(string); ok && name == bundleName {
				bundleData = b
				bundleName = apiName
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found bundle with name '%s' under API key '%s'", name, apiName))
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Bundle with ID '%s' not found in API response", bundleName))
		resp.State.RemoveResource(ctx)
		return
	}

	if name, ok := bundleData["name"].(string); ok {
		data.Name = types.StringValue(name)
	}

	if deviceSettings, ok := bundleData["device_settings"].(string); ok {
		data.DeviceSettings = types.StringValue(deviceSettings)
	} else {
		data.DeviceSettings = types.StringNull()
	}

	if deviceSettingsRefType, ok := bundleData["device_settings_ref_type_"].(string); ok {
		data.DeviceSettingsRefType = types.StringValue(deviceSettingsRefType)
	} else {
		data.DeviceSettingsRefType = types.StringNull()
	}

	if cliCmds, ok := bundleData["cli_commands"]; ok && cliCmds != nil {
		if cmds, ok := cliCmds.(string); ok {
			data.CliCommands = types.StringValue(cmds)
		} else {
			data.CliCommands = types.StringNull()
		}
	} else {
		data.CliCommands = types.StringNull()
	}

	if enableVal, ok := bundleData["enable"].(bool); ok {
		data.Enable = types.BoolValue(enableVal)
	} else {
		data.Enable = types.BoolNull()
	}

	if protocol, ok := bundleData["protocol"].(string); ok {
		data.Protocol = types.StringValue(protocol)
	} else {
		data.Protocol = types.StringNull()
	}

	if diagnosticsProfile, ok := bundleData["diagnostics_profile"].(string); ok {
		data.DiagnosticsProfile = types.StringValue(diagnosticsProfile)
	} else {
		data.DiagnosticsProfile = types.StringNull()
	}

	if diagnosticsProfileRefType, ok := bundleData["diagnostics_profile_ref_type_"].(string); ok {
		data.DiagnosticsProfileRefType = types.StringValue(diagnosticsProfileRefType)
	} else {
		data.DiagnosticsProfileRefType = types.StringNull()
	}

	if deviceVoiceSettings, ok := bundleData["device_voice_settings"].(string); ok {
		data.DeviceVoiceSettings = types.StringValue(deviceVoiceSettings)
	} else {
		data.DeviceVoiceSettings = types.StringNull()
	}

	if deviceVoiceSettingsRefType, ok := bundleData["device_voice_settings_ref_type_"].(string); ok {
		data.DeviceVoiceSettingsRefType = types.StringValue(deviceVoiceSettingsRefType)
	} else {
		data.DeviceVoiceSettingsRefType = types.StringNull()
	}

	// Only set object_properties if it exists in the API response
	if objProps, ok := bundleData["object_properties"].(map[string]interface{}); ok {
		objectProps := verityBundleObjectPropertiesModel{
			IsForSwitch: types.BoolValue(false),
			Group:       types.StringNull(),
			IsPublic:    types.BoolNull(),
		}

		if isForSwitch, ok := objProps["is_for_switch"].(bool); ok {
			objectProps.IsForSwitch = types.BoolValue(isForSwitch)
		}
		if group, ok := objProps["group"].(string); ok {
			objectProps.Group = types.StringValue(group)
		}
		if isPublic, ok := objProps["is_public"].(bool); ok {
			objectProps.IsPublic = types.BoolValue(isPublic)
		}

		data.ObjectProperties = []verityBundleObjectPropertiesModel{objectProps}
	} else {
		data.ObjectProperties = nil
	}

	var ethPortPaths []ethPortPathsModel
	if paths, ok := bundleData["eth_port_paths"].([]interface{}); ok {
		for _, p := range paths {
			path, ok := p.(map[string]interface{})
			if !ok {
				continue
			}

			ethPortPath := ethPortPathsModel{}

			stringFields := map[string]*types.String{
				"eth_port_num_eth_port_profile":                                   &ethPortPath.EthPortNumEthPortProfile,
				"eth_port_num_eth_port_settings":                                  &ethPortPath.EthPortNumEthPortSettings,
				"eth_port_num_eth_port_settings_ref_type_":                        &ethPortPath.EthPortNumEthPortSettingsRefType,
				"eth_port_num_eth_port_profile_ref_type_":                         &ethPortPath.EthPortNumEthPortProfileRefType,
				"diagnostics_port_profile_num_diagnostics_port_profile":           &ethPortPath.DiagnosticsPortProfileNumDiagnosticsPortProfile,
				"diagnostics_port_profile_num_diagnostics_port_profile_ref_type_": &ethPortPath.DiagnosticsPortProfileNumDiagnosticsPortProfileRefType,
				"port_name": &ethPortPath.PortName,
			}

			for apiKey, field := range stringFields {
				if val, ok := path[apiKey].(string); ok {
					*field = types.StringValue(val)
				} else {
					*field = types.StringNull()
				}
			}

			if gwProfile, ok := path["eth_port_num_gateway_profile"]; ok && gwProfile != nil {
				if val, ok := gwProfile.(string); ok {
					ethPortPath.EthPortNumGatewayProfile = types.StringValue(val)
				} else {
					ethPortPath.EthPortNumGatewayProfile = types.StringNull()
				}
			} else {
				ethPortPath.EthPortNumGatewayProfile = types.StringNull()
			}

			if gwProfileRefType, ok := path["eth_port_num_gateway_profile_ref_type_"]; ok && gwProfileRefType != nil {
				if val, ok := gwProfileRefType.(string); ok {
					ethPortPath.EthPortNumGatewayProfileRefType = types.StringValue(val)
				} else {
					ethPortPath.EthPortNumGatewayProfileRefType = types.StringNull()
				}
			} else {
				ethPortPath.EthPortNumGatewayProfileRefType = types.StringNull()
			}

			if index, ok := path["index"].(float64); ok {
				ethPortPath.Index = types.Int64Value(int64(index))
			} else if index, ok := path["index"].(int); ok {
				ethPortPath.Index = types.Int64Value(int64(index))
			} else {
				ethPortPath.Index = types.Int64Null()
			}

			ethPortPaths = append(ethPortPaths, ethPortPath)
		}
	}
	data.EthPortPaths = ethPortPaths

	var userServices []userServicesModel
	if services, ok := bundleData["user_services"].([]interface{}); ok {
		for _, s := range services {
			service, ok := s.(map[string]interface{})
			if !ok {
				continue
			}

			userService := userServicesModel{}

			if enable, ok := service["row_app_enable"].(bool); ok {
				userService.RowAppEnable = types.BoolValue(enable)
			} else {
				userService.RowAppEnable = types.BoolNull()
			}

			stringFields := map[string]struct {
				field  *types.String
				apiKey string
			}{
				"row_app_connected_service":           {&userService.RowAppConnectedService, "row_app_connected_service"},
				"row_app_cli_commands":                {&userService.RowAppCliCommands, "row_app_cli_commands"},
				"row_ip_mask":                         {&userService.RowIpMask, "row_ip_mask"},
				"row_app_connected_service_ref_type_": {&userService.RowAppConnectedServiceRefType, "row_app_connected_service_ref_type_"},
			}

			for _, item := range stringFields {
				if val, ok := service[item.apiKey]; ok && val != nil {
					if strVal, ok := val.(string); ok {
						*item.field = types.StringValue(strVal)
					} else {
						*item.field = types.StringNull()
					}
				} else {
					*item.field = types.StringNull()
				}
			}

			if index, ok := service["index"].(float64); ok {
				userService.Index = types.Int64Value(int64(index))
			} else if index, ok := service["index"].(int); ok {
				userService.Index = types.Int64Value(int64(index))
			} else {
				userService.Index = types.Int64Null()
			}

			userServices = append(userServices, userService)
		}
	}
	data.UserServices = userServices

	var voicePortProfilePaths []voicePortProfilePathsModel
	if paths, ok := bundleData["voice_port_profile_paths"].([]interface{}); ok {
		for _, p := range paths {
			path, ok := p.(map[string]interface{})
			if !ok {
				continue
			}

			voicePortPath := voicePortProfilePathsModel{}

			if voicePortProfiles, ok := path["voice_port_num_voice_port_profiles"].(string); ok {
				voicePortPath.VoicePortNumVoicePortProfiles = types.StringValue(voicePortProfiles)
			} else {
				voicePortPath.VoicePortNumVoicePortProfiles = types.StringNull()
			}

			if voicePortProfilesRefType, ok := path["voice_port_num_voice_port_profiles_ref_type_"].(string); ok {
				voicePortPath.VoicePortNumVoicePortProfilesRefType = types.StringValue(voicePortProfilesRefType)
			} else {
				voicePortPath.VoicePortNumVoicePortProfilesRefType = types.StringNull()
			}

			if index, ok := path["index"].(float64); ok {
				voicePortPath.Index = types.Int64Value(int64(index))
			} else if index, ok := path["index"].(int); ok {
				voicePortPath.Index = types.Int64Value(int64(index))
			} else {
				voicePortPath.Index = types.Int64Null()
			}

			voicePortProfilePaths = append(voicePortProfilePaths, voicePortPath)
		}
	}
	data.VoicePortProfilePaths = voicePortProfilePaths

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *verityBundleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data verityBundleResourceModel
	var state verityBundleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError("Authentication Error", fmt.Sprintf("Unable to authenticate client: %s", err))
		return
	}

	hasChanges := false
	bundleValue := openapi.BundlesPutRequestEndpointBundleValue{}
	name := data.Name.ValueString()

	if !data.CliCommands.Equal(state.CliCommands) {
		if !data.CliCommands.IsNull() {
			bundleValue.CliCommands = openapi.PtrString(data.CliCommands.ValueString())
		} else {
			bundleValue.CliCommands = openapi.PtrString("")
		}
		hasChanges = true
	}

	if !data.Enable.Equal(state.Enable) {
		if !data.Enable.IsNull() {
			bundleValue.Enable = openapi.PtrBool(data.Enable.ValueBool())
		} else {
			bundleValue.Enable = openapi.PtrBool(false)
		}
		hasChanges = true
	}

	if !data.Protocol.Equal(state.Protocol) {
		if !data.Protocol.IsNull() {
			bundleValue.Protocol = openapi.PtrString(data.Protocol.ValueString())
		} else {
			bundleValue.Protocol = openapi.PtrString("")
		}
		hasChanges = true
	}

	deviceSettingsChanged := !data.DeviceSettings.Equal(state.DeviceSettings)
	deviceSettingsRefTypeChanged := !data.DeviceSettingsRefType.Equal(state.DeviceSettingsRefType)

	if deviceSettingsChanged || deviceSettingsRefTypeChanged {
		if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
			data.DeviceSettings, data.DeviceSettingsRefType,
			"device_settings", "device_settings_ref_type_",
			deviceSettingsChanged, deviceSettingsRefTypeChanged) {
			return
		}

		// For fields with one reference type:
		// If only base field changes, send only base field
		// If ref type field changes (or both), send both fields
		if deviceSettingsChanged && !deviceSettingsRefTypeChanged {
			// Just send the base field
			if !data.DeviceSettings.IsNull() && data.DeviceSettings.ValueString() != "" {
				bundleValue.DeviceSettings = openapi.PtrString(data.DeviceSettings.ValueString())
			} else {
				bundleValue.DeviceSettings = openapi.PtrString("")
			}
			hasChanges = true
		} else if deviceSettingsRefTypeChanged {
			// Send both fields
			if !data.DeviceSettings.IsNull() && data.DeviceSettings.ValueString() != "" {
				bundleValue.DeviceSettings = openapi.PtrString(data.DeviceSettings.ValueString())
			} else {
				bundleValue.DeviceSettings = openapi.PtrString("")
			}

			if !data.DeviceSettingsRefType.IsNull() && data.DeviceSettingsRefType.ValueString() != "" {
				bundleValue.DeviceSettingsRefType = openapi.PtrString(data.DeviceSettingsRefType.ValueString())
			} else {
				bundleValue.DeviceSettingsRefType = openapi.PtrString("")
			}
			hasChanges = true
		}
	}

	diagnosticsProfileChanged := !data.DiagnosticsProfile.Equal(state.DiagnosticsProfile)
	diagnosticsProfileRefTypeChanged := !data.DiagnosticsProfileRefType.Equal(state.DiagnosticsProfileRefType)

	if diagnosticsProfileChanged || diagnosticsProfileRefTypeChanged {
		if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
			data.DiagnosticsProfile, data.DiagnosticsProfileRefType,
			"diagnostics_profile", "diagnostics_profile_ref_type_",
			diagnosticsProfileChanged, diagnosticsProfileRefTypeChanged) {
			return
		}

		// Only send the base field if only it changed
		if diagnosticsProfileChanged && !diagnosticsProfileRefTypeChanged {
			// Just send the base field
			if !data.DiagnosticsProfile.IsNull() && data.DiagnosticsProfile.ValueString() != "" {
				bundleValue.DiagnosticsProfile = openapi.PtrString(data.DiagnosticsProfile.ValueString())
			} else {
				bundleValue.DiagnosticsProfile = openapi.PtrString("")
			}
			hasChanges = true
		} else if diagnosticsProfileRefTypeChanged {
			// Send both fields
			if !data.DiagnosticsProfile.IsNull() && data.DiagnosticsProfile.ValueString() != "" {
				bundleValue.DiagnosticsProfile = openapi.PtrString(data.DiagnosticsProfile.ValueString())
			} else {
				bundleValue.DiagnosticsProfile = openapi.PtrString("")
			}

			if !data.DiagnosticsProfileRefType.IsNull() && data.DiagnosticsProfileRefType.ValueString() != "" {
				bundleValue.DiagnosticsProfileRefType = openapi.PtrString(data.DiagnosticsProfileRefType.ValueString())
			} else {
				bundleValue.DiagnosticsProfileRefType = openapi.PtrString("")
			}
			hasChanges = true
		}
	}

	deviceVoiceSettingsChanged := !data.DeviceVoiceSettings.Equal(state.DeviceVoiceSettings)
	deviceVoiceSettingsRefTypeChanged := !data.DeviceVoiceSettingsRefType.Equal(state.DeviceVoiceSettingsRefType)

	if deviceVoiceSettingsChanged || deviceVoiceSettingsRefTypeChanged {
		if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
			data.DeviceVoiceSettings, data.DeviceVoiceSettingsRefType,
			"device_voice_settings", "device_voice_settings_ref_type_",
			deviceVoiceSettingsChanged, deviceVoiceSettingsRefTypeChanged) {
			return
		}

		// For fields with one reference type:
		// If only base field changes, send only base field
		// If ref type field changes (or both), send both fields
		if deviceVoiceSettingsChanged && !deviceVoiceSettingsRefTypeChanged {
			// Just send the base field
			if !data.DeviceVoiceSettings.IsNull() && data.DeviceVoiceSettings.ValueString() != "" {
				bundleValue.DeviceVoiceSettings = openapi.PtrString(data.DeviceVoiceSettings.ValueString())
			} else {
				bundleValue.DeviceVoiceSettings = openapi.PtrString("")
			}
			hasChanges = true
		} else if deviceVoiceSettingsRefTypeChanged {
			// Send both fields
			if !data.DeviceVoiceSettings.IsNull() && data.DeviceVoiceSettings.ValueString() != "" {
				bundleValue.DeviceVoiceSettings = openapi.PtrString(data.DeviceVoiceSettings.ValueString())
			} else {
				bundleValue.DeviceVoiceSettings = openapi.PtrString("")
			}

			if !data.DeviceVoiceSettingsRefType.IsNull() && data.DeviceVoiceSettingsRefType.ValueString() != "" {
				bundleValue.DeviceVoiceSettingsRefType = openapi.PtrString(data.DeviceVoiceSettingsRefType.ValueString())
			} else {
				bundleValue.DeviceVoiceSettingsRefType = openapi.PtrString("")
			}
			hasChanges = true
		}
	}

	if len(data.ObjectProperties) > 0 {
		objPropsChanged := false
		objectProperties := openapi.BundlesPutRequestEndpointBundleValueObjectProperties{}

		if len(state.ObjectProperties) == 0 ||
			!data.ObjectProperties[0].IsForSwitch.Equal(state.ObjectProperties[0].IsForSwitch) {
			objPropsChanged = true

			if !data.ObjectProperties[0].IsForSwitch.IsNull() {
				objectProperties.IsForSwitch = openapi.PtrBool(data.ObjectProperties[0].IsForSwitch.ValueBool())
			} else {
				objectProperties.IsForSwitch = openapi.PtrBool(false)
			}
		}

		if len(state.ObjectProperties) == 0 ||
			!data.ObjectProperties[0].Group.Equal(state.ObjectProperties[0].Group) {
			objPropsChanged = true

			if !data.ObjectProperties[0].Group.IsNull() {
				objectProperties.Group = openapi.PtrString(data.ObjectProperties[0].Group.ValueString())
			} else {
				objectProperties.Group = openapi.PtrString("")
			}
		}

		if len(state.ObjectProperties) == 0 ||
			!data.ObjectProperties[0].IsPublic.Equal(state.ObjectProperties[0].IsPublic) {
			objPropsChanged = true

			if !data.ObjectProperties[0].IsPublic.IsNull() {
				objectProperties.IsPublic = openapi.PtrBool(data.ObjectProperties[0].IsPublic.ValueBool())
			} else {
				objectProperties.IsPublic = openapi.PtrBool(false)
			}
		}

		if objPropsChanged {
			bundleValue.ObjectProperties = &objectProperties
			hasChanges = true
		}
	}

	statePathsByIndex := make(map[int64]ethPortPathsModel)
	for _, path := range state.EthPortPaths {
		if !path.Index.IsNull() {
			statePathsByIndex[path.Index.ValueInt64()] = path
		}
	}

	var changedPaths []openapi.BundlesPutRequestEndpointBundleValueEthPortPathsInner
	ethPortPathsChanged := false

	for _, path := range data.EthPortPaths {
		if path.Index.IsNull() {
			continue
		}

		index := path.Index.ValueInt64()
		statePath, exists := statePathsByIndex[index]

		if !exists {
			// new eth port path, include all fields
			ethPortPath := openapi.BundlesPutRequestEndpointBundleValueEthPortPathsInner{
				Index: openapi.PtrInt32(int32(index)),
			}

			if !path.PortName.IsNull() {
				ethPortPath.PortName = openapi.PtrString(path.PortName.ValueString())
			} else {
				ethPortPath.PortName = openapi.PtrString("")
			}

			hasEthPortSettings := !path.EthPortNumEthPortSettings.IsNull() && path.EthPortNumEthPortSettings.ValueString() != ""
			hasEthPortSettingsRefType := !path.EthPortNumEthPortSettingsRefType.IsNull() && path.EthPortNumEthPortSettingsRefType.ValueString() != ""

			if hasEthPortSettings || hasEthPortSettingsRefType {
				if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
					path.EthPortNumEthPortSettings, path.EthPortNumEthPortSettingsRefType,
					"eth_port_num_eth_port_settings", "eth_port_num_eth_port_settings_ref_type_",
					hasEthPortSettings, hasEthPortSettingsRefType) {
					return
				}
			}

			if !path.EthPortNumEthPortSettings.IsNull() {
				ethPortPath.EthPortNumEthPortSettings = openapi.PtrString(path.EthPortNumEthPortSettings.ValueString())
			} else {
				ethPortPath.EthPortNumEthPortSettings = openapi.PtrString("")
			}

			if !path.EthPortNumEthPortSettingsRefType.IsNull() {
				ethPortPath.EthPortNumEthPortSettingsRefType = openapi.PtrString(path.EthPortNumEthPortSettingsRefType.ValueString())
			} else {
				ethPortPath.EthPortNumEthPortSettingsRefType = openapi.PtrString("")
			}

			hasEthPortProfile := !path.EthPortNumEthPortProfile.IsNull() && path.EthPortNumEthPortProfile.ValueString() != ""
			hasEthPortProfileRefType := !path.EthPortNumEthPortProfileRefType.IsNull() && path.EthPortNumEthPortProfileRefType.ValueString() != ""

			if hasEthPortProfile || hasEthPortProfileRefType {
				if !utils.ValidateMultipleRefTypesSupported(&resp.Diagnostics,
					path.EthPortNumEthPortProfile, path.EthPortNumEthPortProfileRefType,
					"eth_port_num_eth_port_profile", "eth_port_num_eth_port_profile_ref_type_") {
					return
				}
			}

			// Always send both for multiple ref types
			if !path.EthPortNumEthPortProfile.IsNull() {
				ethPortPath.EthPortNumEthPortProfile = openapi.PtrString(path.EthPortNumEthPortProfile.ValueString())
			} else {
				ethPortPath.EthPortNumEthPortProfile = openapi.PtrString("")
			}

			if !path.EthPortNumEthPortProfileRefType.IsNull() {
				ethPortPath.EthPortNumEthPortProfileRefType = openapi.PtrString(path.EthPortNumEthPortProfileRefType.ValueString())
			} else {
				ethPortPath.EthPortNumEthPortProfileRefType = openapi.PtrString("")
			}

			hasGatewayProfile := !path.EthPortNumGatewayProfile.IsNull() && path.EthPortNumGatewayProfile.ValueString() != ""
			hasGatewayProfileRefType := !path.EthPortNumGatewayProfileRefType.IsNull() && path.EthPortNumGatewayProfileRefType.ValueString() != ""

			if hasGatewayProfile || hasGatewayProfileRefType {
				if !utils.ValidateMultipleRefTypesSupported(&resp.Diagnostics,
					path.EthPortNumGatewayProfile, path.EthPortNumGatewayProfileRefType,
					"eth_port_num_gateway_profile", "eth_port_num_gateway_profile_ref_type_") {
					return
				}
			}

			// Always send both fields for multiple ref types
			if !path.EthPortNumGatewayProfile.IsNull() {
				ethPortPath.EthPortNumGatewayProfile = openapi.PtrString(path.EthPortNumGatewayProfile.ValueString())
			} else {
				ethPortPath.EthPortNumGatewayProfile = openapi.PtrString("")
			}

			if !path.EthPortNumGatewayProfileRefType.IsNull() {
				ethPortPath.EthPortNumGatewayProfileRefType = openapi.PtrString(path.EthPortNumGatewayProfileRefType.ValueString())
			} else {
				ethPortPath.EthPortNumGatewayProfileRefType = openapi.PtrString("")
			}

			// Handle diagnostics port profile fields
			hasDiagnosticsProfile := !path.DiagnosticsPortProfileNumDiagnosticsPortProfile.IsNull() && path.DiagnosticsPortProfileNumDiagnosticsPortProfile.ValueString() != ""
			hasDiagnosticsProfileRefType := !path.DiagnosticsPortProfileNumDiagnosticsPortProfileRefType.IsNull() && path.DiagnosticsPortProfileNumDiagnosticsPortProfileRefType.ValueString() != ""

			if hasDiagnosticsProfile || hasDiagnosticsProfileRefType {
				if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
					path.DiagnosticsPortProfileNumDiagnosticsPortProfile, path.DiagnosticsPortProfileNumDiagnosticsPortProfileRefType,
					"diagnostics_port_profile_num_diagnostics_port_profile", "diagnostics_port_profile_num_diagnostics_port_profile_ref_type_",
					hasDiagnosticsProfile, hasDiagnosticsProfileRefType) {
					return
				}
			}

			if !path.DiagnosticsPortProfileNumDiagnosticsPortProfile.IsNull() {
				ethPortPath.DiagnosticsPortProfileNumDiagnosticsPortProfile = openapi.PtrString(path.DiagnosticsPortProfileNumDiagnosticsPortProfile.ValueString())
			} else {
				ethPortPath.DiagnosticsPortProfileNumDiagnosticsPortProfile = openapi.PtrString("")
			}

			if !path.DiagnosticsPortProfileNumDiagnosticsPortProfileRefType.IsNull() {
				ethPortPath.DiagnosticsPortProfileNumDiagnosticsPortProfileRefType = openapi.PtrString(path.DiagnosticsPortProfileNumDiagnosticsPortProfileRefType.ValueString())
			} else {
				ethPortPath.DiagnosticsPortProfileNumDiagnosticsPortProfileRefType = openapi.PtrString("")
			}

			changedPaths = append(changedPaths, ethPortPath)
			ethPortPathsChanged = true
			continue
		}

		// existing eth port path, check which fields changed
		ethPortPath := openapi.BundlesPutRequestEndpointBundleValueEthPortPathsInner{
			Index: openapi.PtrInt32(int32(index)),
		}

		fieldChanged := false

		if !path.PortName.Equal(statePath.PortName) {
			if !path.PortName.IsNull() {
				ethPortPath.PortName = openapi.PtrString(path.PortName.ValueString())
			} else {
				ethPortPath.PortName = openapi.PtrString("")
			}
			fieldChanged = true
		}

		ethPortSettingsChanged := !path.EthPortNumEthPortSettings.Equal(statePath.EthPortNumEthPortSettings)
		ethPortSettingsRefTypeChanged := !path.EthPortNumEthPortSettingsRefType.Equal(statePath.EthPortNumEthPortSettingsRefType)

		if ethPortSettingsChanged || ethPortSettingsRefTypeChanged {
			if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
				path.EthPortNumEthPortSettings, path.EthPortNumEthPortSettingsRefType,
				"eth_port_num_eth_port_settings", "eth_port_num_eth_port_settings_ref_type_",
				ethPortSettingsChanged, ethPortSettingsRefTypeChanged) {
				return
			}

			// For fields with one reference type:
			// If only base field changes, send only base field
			// If ref type field changes (or both), send both fields
			if ethPortSettingsChanged {
				if !path.EthPortNumEthPortSettings.IsNull() {
					ethPortPath.EthPortNumEthPortSettings = openapi.PtrString(path.EthPortNumEthPortSettings.ValueString())
				} else {
					ethPortPath.EthPortNumEthPortSettings = openapi.PtrString("")
				}
			}

			if ethPortSettingsRefTypeChanged {
				if !path.EthPortNumEthPortSettingsRefType.IsNull() {
					ethPortPath.EthPortNumEthPortSettingsRefType = openapi.PtrString(path.EthPortNumEthPortSettingsRefType.ValueString())
				} else {
					ethPortPath.EthPortNumEthPortSettingsRefType = openapi.PtrString("")
				}

				// If ref type changes, also send base field
				if !ethPortSettingsChanged {
					if !path.EthPortNumEthPortSettings.IsNull() {
						ethPortPath.EthPortNumEthPortSettings = openapi.PtrString(path.EthPortNumEthPortSettings.ValueString())
					} else {
						ethPortPath.EthPortNumEthPortSettings = openapi.PtrString("")
					}
				}
			}

			fieldChanged = true
		}

		ethPortProfileChanged := !path.EthPortNumEthPortProfile.Equal(statePath.EthPortNumEthPortProfile)
		ethPortProfileRefTypeChanged := !path.EthPortNumEthPortProfileRefType.Equal(statePath.EthPortNumEthPortProfileRefType)

		if ethPortProfileChanged || ethPortProfileRefTypeChanged {
			if !utils.ValidateMultipleRefTypesSupported(&resp.Diagnostics,
				path.EthPortNumEthPortProfile, path.EthPortNumEthPortProfileRefType,
				"eth_port_num_eth_port_profile", "eth_port_num_eth_port_profile_ref_type_") {
				return
			}

			// For fields with multiple reference types:
			// Always send both fields when either changes
			if !path.EthPortNumEthPortProfile.IsNull() {
				ethPortPath.EthPortNumEthPortProfile = openapi.PtrString(path.EthPortNumEthPortProfile.ValueString())
			} else {
				ethPortPath.EthPortNumEthPortProfile = openapi.PtrString("")
			}

			if !path.EthPortNumEthPortProfileRefType.IsNull() {
				ethPortPath.EthPortNumEthPortProfileRefType = openapi.PtrString(path.EthPortNumEthPortProfileRefType.ValueString())
			} else {
				ethPortPath.EthPortNumEthPortProfileRefType = openapi.PtrString("")
			}

			fieldChanged = true
		}

		gatewayProfileChanged := !path.EthPortNumGatewayProfile.Equal(statePath.EthPortNumGatewayProfile)
		gatewayProfileRefTypeChanged := !path.EthPortNumGatewayProfileRefType.Equal(statePath.EthPortNumGatewayProfileRefType)

		if gatewayProfileChanged || gatewayProfileRefTypeChanged {
			if !utils.ValidateMultipleRefTypesSupported(&resp.Diagnostics,
				path.EthPortNumGatewayProfile, path.EthPortNumGatewayProfileRefType,
				"eth_port_num_gateway_profile", "eth_port_num_gateway_profile_ref_type_") {
				return
			}

			// For fields with multiple reference types:
			// Always send both fields when either changes
			if !path.EthPortNumGatewayProfile.IsNull() {
				ethPortPath.EthPortNumGatewayProfile = openapi.PtrString(path.EthPortNumGatewayProfile.ValueString())
			} else {
				ethPortPath.EthPortNumGatewayProfile = openapi.PtrString("")
			}

			if !path.EthPortNumGatewayProfileRefType.IsNull() {
				ethPortPath.EthPortNumGatewayProfileRefType = openapi.PtrString(path.EthPortNumGatewayProfileRefType.ValueString())
			} else {
				ethPortPath.EthPortNumGatewayProfileRefType = openapi.PtrString("")
			}

			fieldChanged = true
		}

		// Handle diagnostics port profile fields changes
		diagnosticsProfileChanged := !path.DiagnosticsPortProfileNumDiagnosticsPortProfile.Equal(statePath.DiagnosticsPortProfileNumDiagnosticsPortProfile)
		diagnosticsProfileRefTypeChanged := !path.DiagnosticsPortProfileNumDiagnosticsPortProfileRefType.Equal(statePath.DiagnosticsPortProfileNumDiagnosticsPortProfileRefType)

		if diagnosticsProfileChanged || diagnosticsProfileRefTypeChanged {
			if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
				path.DiagnosticsPortProfileNumDiagnosticsPortProfile, path.DiagnosticsPortProfileNumDiagnosticsPortProfileRefType,
				"diagnostics_port_profile_num_diagnostics_port_profile", "diagnostics_port_profile_num_diagnostics_port_profile_ref_type_",
				diagnosticsProfileChanged, diagnosticsProfileRefTypeChanged) {
				return
			}

			// Only send the base field if only it changed
			if diagnosticsProfileChanged && !diagnosticsProfileRefTypeChanged {
				// Just send the base field
				if !path.DiagnosticsPortProfileNumDiagnosticsPortProfile.IsNull() && path.DiagnosticsPortProfileNumDiagnosticsPortProfile.ValueString() != "" {
					ethPortPath.DiagnosticsPortProfileNumDiagnosticsPortProfile = openapi.PtrString(path.DiagnosticsPortProfileNumDiagnosticsPortProfile.ValueString())
				} else {
					ethPortPath.DiagnosticsPortProfileNumDiagnosticsPortProfile = openapi.PtrString("")
				}
			} else if diagnosticsProfileRefTypeChanged {
				// Send both fields
				if !path.DiagnosticsPortProfileNumDiagnosticsPortProfile.IsNull() && path.DiagnosticsPortProfileNumDiagnosticsPortProfile.ValueString() != "" {
					ethPortPath.DiagnosticsPortProfileNumDiagnosticsPortProfile = openapi.PtrString(path.DiagnosticsPortProfileNumDiagnosticsPortProfile.ValueString())
				} else {
					ethPortPath.DiagnosticsPortProfileNumDiagnosticsPortProfile = openapi.PtrString("")
				}

				if !path.DiagnosticsPortProfileNumDiagnosticsPortProfileRefType.IsNull() && path.DiagnosticsPortProfileNumDiagnosticsPortProfileRefType.ValueString() != "" {
					ethPortPath.DiagnosticsPortProfileNumDiagnosticsPortProfileRefType = openapi.PtrString(path.DiagnosticsPortProfileNumDiagnosticsPortProfileRefType.ValueString())
				} else {
					ethPortPath.DiagnosticsPortProfileNumDiagnosticsPortProfileRefType = openapi.PtrString("")
				}
			}

			fieldChanged = true
		}

		if fieldChanged {
			changedPaths = append(changedPaths, ethPortPath)
			ethPortPathsChanged = true
		}
	}

	for idx := range statePathsByIndex {
		found := false
		for _, path := range data.EthPortPaths {
			if !path.Index.IsNull() && path.Index.ValueInt64() == idx {
				found = true
				break
			}
		}

		if !found {
			// Path removed - include only the index for deletion
			deletedPath := openapi.BundlesPutRequestEndpointBundleValueEthPortPathsInner{
				Index: openapi.PtrInt32(int32(idx)),
			}
			changedPaths = append(changedPaths, deletedPath)
			ethPortPathsChanged = true
		}
	}

	if ethPortPathsChanged && len(changedPaths) > 0 {
		bundleValue.EthPortPaths = changedPaths
		hasChanges = true
	}

	stateServicesByIndex := make(map[int64]userServicesModel)
	for _, service := range state.UserServices {
		if !service.Index.IsNull() {
			stateServicesByIndex[service.Index.ValueInt64()] = service
		}
	}

	var changedServices []openapi.BundlesPutRequestEndpointBundleValueUserServicesInner
	userServicesChanged := false

	for _, service := range data.UserServices {
		if service.Index.IsNull() {
			continue
		}

		index := service.Index.ValueInt64()
		stateService, exists := stateServicesByIndex[index]

		if !exists {
			// new user service, include all fields
			userService := openapi.BundlesPutRequestEndpointBundleValueUserServicesInner{
				Index: openapi.PtrInt32(int32(index)),
			}

			if !service.RowAppEnable.IsNull() {
				userService.RowAppEnable = openapi.PtrBool(service.RowAppEnable.ValueBool())
			} else {
				userService.RowAppEnable = openapi.PtrBool(false)
			}

			hasConnectedService := !service.RowAppConnectedService.IsNull() && service.RowAppConnectedService.ValueString() != ""
			hasConnectedServiceRefType := !service.RowAppConnectedServiceRefType.IsNull() && service.RowAppConnectedServiceRefType.ValueString() != ""

			if hasConnectedService || hasConnectedServiceRefType {
				if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
					service.RowAppConnectedService, service.RowAppConnectedServiceRefType,
					"row_app_connected_service", "row_app_connected_service_ref_type_",
					hasConnectedService, hasConnectedServiceRefType) {
					return
				}
			}

			if !service.RowAppConnectedService.IsNull() {
				userService.RowAppConnectedService = openapi.PtrString(service.RowAppConnectedService.ValueString())
			} else {
				userService.RowAppConnectedService = openapi.PtrString("")
			}

			if !service.RowAppConnectedServiceRefType.IsNull() {
				userService.RowAppConnectedServiceRefType = openapi.PtrString(service.RowAppConnectedServiceRefType.ValueString())
			} else {
				userService.RowAppConnectedServiceRefType = openapi.PtrString("")
			}

			if !service.RowAppCliCommands.IsNull() {
				userService.RowAppCliCommands = openapi.PtrString(service.RowAppCliCommands.ValueString())
			} else {
				userService.RowAppCliCommands = openapi.PtrString("")
			}

			if !service.RowIpMask.IsNull() {
				userService.RowIpMask = openapi.PtrString(service.RowIpMask.ValueString())
			} else {
				userService.RowIpMask = openapi.PtrString("")
			}

			changedServices = append(changedServices, userService)
			userServicesChanged = true
			continue
		}

		// existing user service, check which fields changed
		userService := openapi.BundlesPutRequestEndpointBundleValueUserServicesInner{
			Index: openapi.PtrInt32(int32(index)),
		}

		fieldChanged := false

		if !service.RowAppEnable.Equal(stateService.RowAppEnable) {
			userService.RowAppEnable = openapi.PtrBool(service.RowAppEnable.ValueBool())
			fieldChanged = true
		}

		connectedServiceChanged := !service.RowAppConnectedService.Equal(stateService.RowAppConnectedService)
		connectedServiceRefTypeChanged := !service.RowAppConnectedServiceRefType.Equal(stateService.RowAppConnectedServiceRefType)

		if connectedServiceChanged || connectedServiceRefTypeChanged {
			if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
				service.RowAppConnectedService, service.RowAppConnectedServiceRefType,
				"row_app_connected_service", "row_app_connected_service_ref_type_",
				connectedServiceChanged, connectedServiceRefTypeChanged) {
				return
			}

			// For fields with one reference type:
			// If only base field changes, send only base field
			// If ref type field changes (or both), send both fields
			if connectedServiceChanged && !connectedServiceRefTypeChanged {
				// Just send the base field
				if !service.RowAppConnectedService.IsNull() && service.RowAppConnectedService.ValueString() != "" {
					userService.RowAppConnectedService = openapi.PtrString(service.RowAppConnectedService.ValueString())
				} else {
					userService.RowAppConnectedService = openapi.PtrString("")
				}
			} else if connectedServiceRefTypeChanged {
				// Send both fields
				if !service.RowAppConnectedService.IsNull() && service.RowAppConnectedService.ValueString() != "" {
					userService.RowAppConnectedService = openapi.PtrString(service.RowAppConnectedService.ValueString())
				} else {
					userService.RowAppConnectedService = openapi.PtrString("")
				}

				if !service.RowAppConnectedServiceRefType.IsNull() && service.RowAppConnectedServiceRefType.ValueString() != "" {
					userService.RowAppConnectedServiceRefType = openapi.PtrString(service.RowAppConnectedServiceRefType.ValueString())
				} else {
					userService.RowAppConnectedServiceRefType = openapi.PtrString("")
				}
			}

			fieldChanged = true
		}

		if !service.RowAppCliCommands.Equal(stateService.RowAppCliCommands) {
			if !service.RowAppCliCommands.IsNull() {
				userService.RowAppCliCommands = openapi.PtrString(service.RowAppCliCommands.ValueString())
			} else {
				userService.RowAppCliCommands = openapi.PtrString("")
			}
			fieldChanged = true
		}

		if !service.RowIpMask.Equal(stateService.RowIpMask) {
			if !service.RowIpMask.IsNull() {
				userService.RowIpMask = openapi.PtrString(service.RowIpMask.ValueString())
			} else {
				userService.RowIpMask = openapi.PtrString("")
			}
			fieldChanged = true
		}

		if fieldChanged {
			changedServices = append(changedServices, userService)
			userServicesChanged = true
		}
	}

	for idx := range stateServicesByIndex {
		found := false
		for _, service := range data.UserServices {
			if !service.Index.IsNull() && service.Index.ValueInt64() == idx {
				found = true
				break
			}
		}

		if !found {
			// service removed - include only the index for deletion
			deletedService := openapi.BundlesPutRequestEndpointBundleValueUserServicesInner{
				Index: openapi.PtrInt32(int32(idx)),
			}
			changedServices = append(changedServices, deletedService)
			userServicesChanged = true
		}
	}

	if userServicesChanged && len(changedServices) > 0 {
		bundleValue.UserServices = changedServices
		hasChanges = true
	}

	stateVoicePortProfilePathsByIndex := make(map[int64]voicePortProfilePathsModel)
	for _, path := range state.VoicePortProfilePaths {
		if !path.Index.IsNull() {
			stateVoicePortProfilePathsByIndex[path.Index.ValueInt64()] = path
		}
	}

	var changedVoicePortProfilePaths []openapi.BundlesPutRequestEndpointBundleValueVoicePortProfilePathsInner
	voicePortProfilePathsChanged := false

	for _, path := range data.VoicePortProfilePaths {
		if path.Index.IsNull() {
			continue
		}

		index := path.Index.ValueInt64()
		statePath, exists := stateVoicePortProfilePathsByIndex[index]

		if !exists {
			// new voice port profile path, include all fields
			voicePortPath := openapi.BundlesPutRequestEndpointBundleValueVoicePortProfilePathsInner{
				Index: openapi.PtrInt32(int32(index)),
			}

			hasVoicePortProfiles := !path.VoicePortNumVoicePortProfiles.IsNull() && path.VoicePortNumVoicePortProfiles.ValueString() != ""
			hasVoicePortProfilesRefType := !path.VoicePortNumVoicePortProfilesRefType.IsNull() && path.VoicePortNumVoicePortProfilesRefType.ValueString() != ""

			if hasVoicePortProfiles || hasVoicePortProfilesRefType {
				if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
					path.VoicePortNumVoicePortProfiles, path.VoicePortNumVoicePortProfilesRefType,
					"voice_port_num_voice_port_profiles", "voice_port_num_voice_port_profiles_ref_type_",
					hasVoicePortProfiles, hasVoicePortProfilesRefType) {
					return
				}
			}

			if !path.VoicePortNumVoicePortProfiles.IsNull() {
				voicePortPath.VoicePortNumVoicePortProfiles = openapi.PtrString(path.VoicePortNumVoicePortProfiles.ValueString())
			} else {
				voicePortPath.VoicePortNumVoicePortProfiles = openapi.PtrString("")
			}

			if !path.VoicePortNumVoicePortProfilesRefType.IsNull() {
				voicePortPath.VoicePortNumVoicePortProfilesRefType = openapi.PtrString(path.VoicePortNumVoicePortProfilesRefType.ValueString())
			} else {
				voicePortPath.VoicePortNumVoicePortProfilesRefType = openapi.PtrString("")
			}

			changedVoicePortProfilePaths = append(changedVoicePortProfilePaths, voicePortPath)
			voicePortProfilePathsChanged = true
			continue
		}

		// existing voice port profile path, check which fields changed
		voicePortPath := openapi.BundlesPutRequestEndpointBundleValueVoicePortProfilePathsInner{
			Index: openapi.PtrInt32(int32(index)),
		}

		fieldChanged := false

		voicePortProfilesChanged := !path.VoicePortNumVoicePortProfiles.Equal(statePath.VoicePortNumVoicePortProfiles)
		voicePortProfilesRefTypeChanged := !path.VoicePortNumVoicePortProfilesRefType.Equal(statePath.VoicePortNumVoicePortProfilesRefType)

		if voicePortProfilesChanged || voicePortProfilesRefTypeChanged {
			if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
				path.VoicePortNumVoicePortProfiles, path.VoicePortNumVoicePortProfilesRefType,
				"voice_port_num_voice_port_profiles", "voice_port_num_voice_port_profiles_ref_type_",
				voicePortProfilesChanged, voicePortProfilesRefTypeChanged) {
				return
			}

			// For fields with one reference type:
			// If only base field changes, send only base field
			// If ref type field changes (or both), send both fields
			if voicePortProfilesChanged {
				if !path.VoicePortNumVoicePortProfiles.IsNull() {
					voicePortPath.VoicePortNumVoicePortProfiles = openapi.PtrString(path.VoicePortNumVoicePortProfiles.ValueString())
				} else {
					voicePortPath.VoicePortNumVoicePortProfiles = openapi.PtrString("")
				}
			}

			if voicePortProfilesRefTypeChanged {
				if !path.VoicePortNumVoicePortProfilesRefType.IsNull() {
					voicePortPath.VoicePortNumVoicePortProfilesRefType = openapi.PtrString(path.VoicePortNumVoicePortProfilesRefType.ValueString())
				} else {
					voicePortPath.VoicePortNumVoicePortProfilesRefType = openapi.PtrString("")
				}

				// If ref type changes, also send base field
				if !voicePortProfilesChanged {
					if !path.VoicePortNumVoicePortProfiles.IsNull() {
						voicePortPath.VoicePortNumVoicePortProfiles = openapi.PtrString(path.VoicePortNumVoicePortProfiles.ValueString())
					} else {
						voicePortPath.VoicePortNumVoicePortProfiles = openapi.PtrString("")
					}
				}
			}

			fieldChanged = true
		}

		if fieldChanged {
			changedVoicePortProfilePaths = append(changedVoicePortProfilePaths, voicePortPath)
			voicePortProfilePathsChanged = true
		}
	}

	for idx := range stateVoicePortProfilePathsByIndex {
		found := false
		for _, path := range data.VoicePortProfilePaths {
			if !path.Index.IsNull() && path.Index.ValueInt64() == idx {
				found = true
				break
			}
		}

		if !found {
			// Path removed - include only the index for deletion
			deletedPath := openapi.BundlesPutRequestEndpointBundleValueVoicePortProfilePathsInner{
				Index: openapi.PtrInt32(int32(idx)),
			}
			changedVoicePortProfilePaths = append(changedVoicePortProfilePaths, deletedPath)
			voicePortProfilePathsChanged = true
		}
	}

	if voicePortProfilePathsChanged && len(changedVoicePortProfilePaths) > 0 {
		bundleValue.VoicePortProfilePaths = changedVoicePortProfilePaths
		hasChanges = true
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	operationID := r.bulkOpsMgr.AddPatch(ctx, "bundle", name, bundleValue)
	r.notifyOperationAdded()

	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update Bundle %s", name))...,
		)
		return
	}

	clearCache(ctx, r.provCtx, "bundles")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *verityBundleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityBundleResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// API version check: Only allow on 6.5+
	apiVersion, err := getApiVersion(ctx, r.provCtx)
	if err != nil {
		resp.Diagnostics.AddError("API Version Error", fmt.Sprintf("Unable to determine API version: %s", err))
		return
	}
	if apiVersion < "6.5" {
		resp.Diagnostics.AddError("Unsupported API Version", "Bundle resource deletion is only supported on API version 6.5 and above.")
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
	operationID := r.bulkOpsMgr.AddDelete(ctx, "bundle", name)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for bundle deletion operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Delete Bundle %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Bundle %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "bundles")
	resp.State.RemoveResource(ctx)
}

func (r *verityBundleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
