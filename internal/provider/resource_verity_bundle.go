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
	Name                  types.String                        `tfsdk:"name"`
	DeviceSettings        types.String                        `tfsdk:"device_settings"`
	DeviceSettingsRefType types.String                        `tfsdk:"device_settings_ref_type_"`
	CliCommands           types.String                        `tfsdk:"cli_commands"`
	ObjectProperties      []verityBundleObjectPropertiesModel `tfsdk:"object_properties"`
	EthPortPaths          []ethPortPathsModel                 `tfsdk:"eth_port_paths"`
	UserServices          []userServicesModel                 `tfsdk:"user_services"`
}

type verityBundleObjectPropertiesModel struct {
	IsForSwitch types.Bool `tfsdk:"is_for_switch"`
}

type ethPortPathsModel struct {
	EthPortNumEthPortProfile         types.String `tfsdk:"eth_port_num_eth_port_profile"`
	EthPortNumEthPortProfileRefType  types.String `tfsdk:"eth_port_num_eth_port_profile_ref_type_"`
	EthPortNumEthPortSettings        types.String `tfsdk:"eth_port_num_eth_port_settings"`
	EthPortNumEthPortSettingsRefType types.String `tfsdk:"eth_port_num_eth_port_settings_ref_type_"`
	EthPortNumGatewayProfile         types.String `tfsdk:"eth_port_num_gateway_profile"`
	EthPortNumGatewayProfileRefType  types.String `tfsdk:"eth_port_num_gateway_profile_ref_type_"`
	PortName                         types.String `tfsdk:"port_name"`
	Index                            types.Int64  `tfsdk:"index"`
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
		},
	}
}

func (r *verityBundleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data verityBundleResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := ensureAuthenticated(ctx, r.provCtx); err != nil {
		resp.Diagnostics.AddError("Authentication Error", fmt.Sprintf("Unable to authenticate client: %s", err))
		return
	}

	clearCache(ctx, r.provCtx, "bundles")

	name := data.Name.ValueString()
	data.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	if !resp.Diagnostics.HasError() {
		readReq := resource.ReadRequest{
			State: resp.State,
		}
		readResp := resource.ReadResponse{
			State:       resp.State,
			Diagnostics: resp.Diagnostics,
		}
		r.Read(ctx, readReq, &readResp)
		resp.Diagnostics = readResp.Diagnostics
	}
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

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentBundleOperations() {
		tflog.Info(ctx, fmt.Sprintf("Skipping bundle %s verification - trusting recent successful API operation", bundleName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("No recent bundle operations found, performing normal verification for %s", bundleName))

	type BundleResponse struct {
		EndpointBundle map[string]map[string]interface{} `json:"endpoint_bundle"`
	}

	bundlesData, err := getCachedResponse(ctx, r.provCtx, "bundles", func() (interface{}, error) {
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
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Bundle",
			fmt.Sprintf("Could not read bundle %s: %s", bundleName, err),
		)
		return
	}

	result := bundlesData.(BundleResponse)
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

	objectProps := verityBundleObjectPropertiesModel{
		IsForSwitch: types.BoolValue(false),
	}

	if objProps, ok := bundleData["object_properties"].(map[string]interface{}); ok {
		if isForSwitch, ok := objProps["is_for_switch"].(bool); ok {
			objectProps.IsForSwitch = types.BoolValue(isForSwitch)
		}
	}

	data.ObjectProperties = []verityBundleObjectPropertiesModel{objectProps}

	var ethPortPaths []ethPortPathsModel
	if paths, ok := bundleData["eth_port_paths"].([]interface{}); ok {
		for _, p := range paths {
			path, ok := p.(map[string]interface{})
			if !ok {
				continue
			}

			ethPortPath := ethPortPathsModel{}

			stringFields := map[string]*types.String{
				"eth_port_num_eth_port_profile":            &ethPortPath.EthPortNumEthPortProfile,
				"eth_port_num_eth_port_settings":           &ethPortPath.EthPortNumEthPortSettings,
				"eth_port_num_eth_port_settings_ref_type_": &ethPortPath.EthPortNumEthPortSettingsRefType,
				"eth_port_num_eth_port_profile_ref_type_":  &ethPortPath.EthPortNumEthPortProfileRefType,
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
	bundleValue := openapi.BundlesPatchRequestEndpointBundleValue{}
	name := data.Name.ValueString()

	if !data.CliCommands.Equal(state.CliCommands) {
		if !data.CliCommands.IsNull() {
			bundleValue.CliCommands = openapi.PtrString(data.CliCommands.ValueString())
		} else {
			bundleValue.CliCommands = openapi.PtrString("")
		}
		hasChanges = true
	}

	if !data.DeviceSettings.Equal(state.DeviceSettings) {
		if !data.DeviceSettings.IsNull() {
			bundleValue.DeviceSettings = openapi.PtrString(data.DeviceSettings.ValueString())
		} else {
			bundleValue.DeviceSettings = openapi.PtrString("")
		}
		hasChanges = true
	}

	if !data.DeviceSettingsRefType.Equal(state.DeviceSettingsRefType) {
		if !data.DeviceSettingsRefType.IsNull() {
			bundleValue.DeviceSettingsRefType = openapi.PtrString(data.DeviceSettingsRefType.ValueString())
		} else {
			bundleValue.DeviceSettingsRefType = openapi.PtrString("")
		}
		hasChanges = true
	}

	if len(data.ObjectProperties) > 0 {
		objPropsChanged := false
		objectProperties := openapi.BundlesPatchRequestEndpointBundleValueObjectProperties{}

		if len(state.ObjectProperties) == 0 ||
			!data.ObjectProperties[0].IsForSwitch.Equal(state.ObjectProperties[0].IsForSwitch) {
			objPropsChanged = true

			if !data.ObjectProperties[0].IsForSwitch.IsNull() {
				objectProperties.IsForSwitch = openapi.PtrBool(data.ObjectProperties[0].IsForSwitch.ValueBool())
			} else {
				objectProperties.IsForSwitch = openapi.PtrBool(false)
			}
		}

		if objPropsChanged {
			bundleValue.ObjectProperties = &objectProperties
			hasChanges = true
		}
	}

	if len(data.EthPortPaths) > 0 {
		statePathsByIndex := make(map[int64]ethPortPathsModel)
		for _, path := range state.EthPortPaths {
			if !path.Index.IsNull() {
				statePathsByIndex[path.Index.ValueInt64()] = path
			}
		}

		var changedPaths []openapi.BundlesPatchRequestEndpointBundleValueEthPortPathsInner

		for _, path := range data.EthPortPaths {
			if path.Index.IsNull() {
				continue
			}

			index := path.Index.ValueInt64()
			statePath, exists := statePathsByIndex[index]
			pathChanged := false

			if !exists ||
				!path.PortName.Equal(statePath.PortName) ||
				!path.EthPortNumEthPortSettings.Equal(statePath.EthPortNumEthPortSettings) ||
				!path.EthPortNumEthPortProfile.Equal(statePath.EthPortNumEthPortProfile) ||
				!path.EthPortNumGatewayProfile.Equal(statePath.EthPortNumGatewayProfile) ||
				!path.EthPortNumEthPortProfileRefType.Equal(statePath.EthPortNumEthPortProfileRefType) ||
				!path.EthPortNumEthPortSettingsRefType.Equal(statePath.EthPortNumEthPortSettingsRefType) ||
				!path.EthPortNumGatewayProfileRefType.Equal(statePath.EthPortNumGatewayProfileRefType) {
				pathChanged = true
			}

			if pathChanged {
				ethPortPath := openapi.BundlesPatchRequestEndpointBundleValueEthPortPathsInner{
					Index: openapi.PtrInt32(int32(index)),
				}

				if !path.PortName.IsNull() {
					ethPortPath.PortName = openapi.PtrString(path.PortName.ValueString())
				} else {
					ethPortPath.PortName = openapi.PtrString("")
				}

				if !path.EthPortNumEthPortSettings.IsNull() {
					ethPortPath.EthPortNumEthPortSettings = openapi.PtrString(path.EthPortNumEthPortSettings.ValueString())
				} else {
					ethPortPath.EthPortNumEthPortSettings = openapi.PtrString("")
				}

				if !path.EthPortNumEthPortProfile.IsNull() {
					ethPortPath.EthPortNumEthPortProfile = openapi.PtrString(path.EthPortNumEthPortProfile.ValueString())
				} else {
					ethPortPath.EthPortNumEthPortProfile = openapi.PtrString("")
				}

				if !path.EthPortNumGatewayProfile.IsNull() {
					ethPortPath.EthPortNumGatewayProfile = openapi.PtrString(path.EthPortNumGatewayProfile.ValueString())
				} else {
					ethPortPath.EthPortNumGatewayProfile = openapi.PtrString("")
				}

				if !path.EthPortNumEthPortProfileRefType.IsNull() {
					ethPortPath.EthPortNumEthPortProfileRefType = openapi.PtrString(path.EthPortNumEthPortProfileRefType.ValueString())
				} else {
					ethPortPath.EthPortNumEthPortProfileRefType = openapi.PtrString("")
				}

				if !path.EthPortNumEthPortSettingsRefType.IsNull() {
					ethPortPath.EthPortNumEthPortSettingsRefType = openapi.PtrString(path.EthPortNumEthPortSettingsRefType.ValueString())
				} else {
					ethPortPath.EthPortNumEthPortSettingsRefType = openapi.PtrString("")
				}

				if !path.EthPortNumGatewayProfileRefType.IsNull() {
					ethPortPath.EthPortNumGatewayProfileRefType = openapi.PtrString(path.EthPortNumGatewayProfileRefType.ValueString())
				} else {
					ethPortPath.EthPortNumGatewayProfileRefType = openapi.PtrString("")
				}

				changedPaths = append(changedPaths, ethPortPath)
			}
		}

		if len(changedPaths) > 0 {
			bundleValue.EthPortPaths = changedPaths
			hasChanges = true
		}
	}

	if len(data.UserServices) > 0 {
		stateServicesByIndex := make(map[int64]userServicesModel)
		for _, service := range state.UserServices {
			if !service.Index.IsNull() {
				stateServicesByIndex[service.Index.ValueInt64()] = service
			}
		}

		var changedServices []openapi.BundlesPatchRequestEndpointBundleValueUserServicesInner

		for _, service := range data.UserServices {
			if service.Index.IsNull() {
				continue
			}

			index := service.Index.ValueInt64()
			stateService, exists := stateServicesByIndex[index]
			serviceChanged := false

			if !exists ||
				!service.RowAppEnable.Equal(stateService.RowAppEnable) ||
				!service.RowAppConnectedService.Equal(stateService.RowAppConnectedService) ||
				!service.RowAppCliCommands.Equal(stateService.RowAppCliCommands) ||
				!service.RowIpMask.Equal(stateService.RowIpMask) ||
				!service.RowAppConnectedServiceRefType.Equal(stateService.RowAppConnectedServiceRefType) {
				serviceChanged = true
			}

			if serviceChanged {
				userService := openapi.BundlesPatchRequestEndpointBundleValueUserServicesInner{
					Index: openapi.PtrInt32(int32(index)),
				}

				if !service.RowAppEnable.IsNull() {
					userService.RowAppEnable = openapi.PtrBool(service.RowAppEnable.ValueBool())
				} else {
					userService.RowAppEnable = openapi.PtrBool(false)
				}

				if !service.RowAppConnectedService.IsNull() {
					userService.RowAppConnectedService = openapi.PtrString(service.RowAppConnectedService.ValueString())
				} else {
					userService.RowAppConnectedService = openapi.PtrString("")
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

				if !service.RowAppConnectedServiceRefType.IsNull() {
					userService.RowAppConnectedServiceRefType = openapi.PtrString(service.RowAppConnectedServiceRefType.ValueString())
				} else {
					userService.RowAppConnectedServiceRefType = openapi.PtrString("")
				}

				changedServices = append(changedServices, userService)
			}
		}

		if len(changedServices) > 0 {
			bundleValue.UserServices = changedServices
			hasChanges = true
		}
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
		return
	}

	operationID := r.bulkOpsMgr.AddBundlePatch(ctx, name, bundleValue)
	r.notifyOperationAdded()

	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Bundle",
			fmt.Sprintf("Could not update bundle %s: %s", name, err),
		)
		return
	}

	clearCache(ctx, r.provCtx, "bundles")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *verityBundleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddError(
		"Operation Not Supported",
		"Deletion of bundles is not supported by the API",
	)
}

func (r *verityBundleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
