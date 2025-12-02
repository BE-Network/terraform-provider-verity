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

func (epp ethPortPathsModel) GetIndex() types.Int64 {
	return epp.Index
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
	resp.Diagnostics.AddError(
		"Create Not Supported",
		"Bundle resources cannot be created. They represent existing bundle configurations that can only be read and updated.",
	)
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
			},
		}
	} else {
		state.ObjectProperties = nil
	}

	// Map string fields
	stringFieldMappings := map[string]*types.String{
		"device_settings":           &state.DeviceSettings,
		"device_settings_ref_type_": &state.DeviceSettingsRefType,
		"cli_commands":              &state.CliCommands,
	}

	for apiKey, stateField := range stringFieldMappings {
		*stateField = utils.MapStringFromAPI(bundleMap[apiKey])
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
				EthPortNumEthPortProfile:         utils.MapStringFromAPI(path["eth_port_num_eth_port_profile"]),
				EthPortNumEthPortSettings:        utils.MapStringFromAPI(path["eth_port_num_eth_port_settings"]),
				EthPortNumEthPortSettingsRefType: utils.MapStringFromAPI(path["eth_port_num_eth_port_settings_ref_type_"]),
				EthPortNumEthPortProfileRefType:  utils.MapStringFromAPI(path["eth_port_num_eth_port_profile_ref_type_"]),
				PortName:                         utils.MapStringFromAPI(path["port_name"]),
				EthPortNumGatewayProfile:         utils.MapStringFromAPI(path["eth_port_num_gateway_profile"]),
				EthPortNumGatewayProfileRefType:  utils.MapStringFromAPI(path["eth_port_num_gateway_profile_ref_type_"]),
				Index:                            utils.MapInt64FromAPI(path["index"]),
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
	bundleProps := openapi.BundlesPatchRequestEndpointBundleValue{}
	hasChanges := false

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { bundleProps.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.CliCommands, state.CliCommands, func(v *string) { bundleProps.CliCommands = v }, &hasChanges)

	// Handle object properties
	if len(plan.ObjectProperties) > 0 && len(state.ObjectProperties) > 0 {
		objProps := openapi.BundlesPatchRequestEndpointBundleValueObjectProperties{}
		op := plan.ObjectProperties[0]
		st := state.ObjectProperties[0]
		objPropsChanged := false

		utils.CompareAndSetObjectPropertiesFields([]utils.ObjectPropertiesFieldWithComparison{
			{Name: "IsForSwitch", PlanValue: op.IsForSwitch, StateValue: st.IsForSwitch, APIValue: &objProps.IsForSwitch},
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

	// Handle eth port paths
	ethPortPathsHandler := utils.IndexedItemHandler[ethPortPathsModel, openapi.BundlesPatchRequestEndpointBundleValueEthPortPathsInner]{
		CreateNew: func(planItem ethPortPathsModel) openapi.BundlesPatchRequestEndpointBundleValueEthPortPathsInner {
			ethPortPath := openapi.BundlesPatchRequestEndpointBundleValueEthPortPathsInner{}

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
			})

			return ethPortPath
		},
		UpdateExisting: func(planItem ethPortPathsModel, stateItem ethPortPathsModel) (openapi.BundlesPatchRequestEndpointBundleValueEthPortPathsInner, bool) {
			ethPortPath := openapi.BundlesPatchRequestEndpointBundleValueEthPortPathsInner{}

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

			return ethPortPath, fieldChanged
		},
		CreateDeleted: func(index int64) openapi.BundlesPatchRequestEndpointBundleValueEthPortPathsInner {
			item := openapi.BundlesPatchRequestEndpointBundleValueEthPortPathsInner{}
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
	userServicesHandler := utils.IndexedItemHandler[userServicesModel, openapi.BundlesPatchRequestEndpointBundleValueUserServicesInner]{
		CreateNew: func(planItem userServicesModel) openapi.BundlesPatchRequestEndpointBundleValueUserServicesInner {
			userService := openapi.BundlesPatchRequestEndpointBundleValueUserServicesInner{}

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
		UpdateExisting: func(planItem userServicesModel, stateItem userServicesModel) (openapi.BundlesPatchRequestEndpointBundleValueUserServicesInner, bool) {
			userService := openapi.BundlesPatchRequestEndpointBundleValueUserServicesInner{}

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
		CreateDeleted: func(index int64) openapi.BundlesPatchRequestEndpointBundleValueUserServicesInner {
			item := openapi.BundlesPatchRequestEndpointBundleValueUserServicesInner{}
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
	resp.Diagnostics.AddError(
		"Delete Not Supported",
		"Bundle resources cannot be deleted. They represent existing bundle configurations that can only be read and updated.",
	)
}

func (r *verityBundleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
