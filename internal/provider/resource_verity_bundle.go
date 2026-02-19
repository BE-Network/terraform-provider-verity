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
	_ resource.ResourceWithModifyPlan  = &verityBundleResource{}
)

const bundleResourceType = "bundles"

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
				Computed:    true,
			},
			"device_settings": schema.StringAttribute{
				Description: "Device Settings for device",
				Optional:    true,
				Computed:    true,
			},
			"device_settings_ref_type_": schema.StringAttribute{
				Description: "Object type for device_settings field",
				Optional:    true,
				Computed:    true,
			},
			"cli_commands": schema.StringAttribute{
				Description: "CLI Commands",
				Optional:    true,
				Computed:    true,
			},
			"protocol": schema.StringAttribute{
				Description: "Voice Protocol: MGCP or SIP",
				Optional:    true,
				Computed:    true,
			},
			"diagnostics_profile": schema.StringAttribute{
				Description: "Diagnostics Profile for device",
				Optional:    true,
				Computed:    true,
			},
			"diagnostics_profile_ref_type_": schema.StringAttribute{
				Description: "Object type for diagnostics_profile field",
				Optional:    true,
				Computed:    true,
			},
			"device_voice_settings": schema.StringAttribute{
				Description: "Device Voice Settings for device",
				Optional:    true,
				Computed:    true,
			},
			"device_voice_settings_ref_type_": schema.StringAttribute{
				Description: "Object type for device_voice_settings field",
				Optional:    true,
				Computed:    true,
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
							Computed:    true,
						},
						"group": schema.StringAttribute{
							Description: "Group",
							Optional:    true,
							Computed:    true,
						},
						"is_public": schema.BoolAttribute{
							Description: "Denotes a shared Switch Bundle",
							Optional:    true,
							Computed:    true,
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
							Computed:    true,
						},
						"eth_port_num_eth_port_profile_ref_type_": schema.StringAttribute{
							Description: "Object type for eth_port_num_eth_port_profile field",
							Optional:    true,
							Computed:    true,
						},
						"eth_port_num_eth_port_settings": schema.StringAttribute{
							Description: "Choose an Eth Port Settings",
							Optional:    true,
							Computed:    true,
						},
						"eth_port_num_eth_port_settings_ref_type_": schema.StringAttribute{
							Description: "Object type for eth_port_num_eth_port_settings field",
							Optional:    true,
							Computed:    true,
						},
						"eth_port_num_gateway_profile": schema.StringAttribute{
							Description: "Gateway Profile or LAG for Eth Port",
							Optional:    true,
							Computed:    true,
						},
						"eth_port_num_gateway_profile_ref_type_": schema.StringAttribute{
							Description: "Object type for eth_port_num_gateway_profile field",
							Optional:    true,
							Computed:    true,
						},
						"diagnostics_port_profile_num_diagnostics_port_profile": schema.StringAttribute{
							Description: "Diagnostics Port Profile for port",
							Optional:    true,
							Computed:    true,
						},
						"diagnostics_port_profile_num_diagnostics_port_profile_ref_type_": schema.StringAttribute{
							Description: "Object type for diagnostics_port_profile_num_diagnostics_port_profile field",
							Optional:    true,
							Computed:    true,
						},
						"port_name": schema.StringAttribute{
							Description: "The name identifying the port",
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
			"user_services": schema.ListNestedBlock{
				Description: "List of user services configurations",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"row_app_enable": schema.BoolAttribute{
							Description: "Enable of this User application",
							Optional:    true,
							Computed:    true,
						},
						"row_app_connected_service": schema.StringAttribute{
							Description: "Service connected to this User application",
							Optional:    true,
							Computed:    true,
						},
						"row_app_connected_service_ref_type_": schema.StringAttribute{
							Description: "Object type for row_app_connected_service field",
							Optional:    true,
							Computed:    true,
						},
						"row_app_cli_commands": schema.StringAttribute{
							Description: "CLI Commands of this User application",
							Optional:    true,
							Computed:    true,
						},
						"row_ip_mask": schema.StringAttribute{
							Description: "IP/Mask in IPv4 format",
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
			"voice_port_profile_paths": schema.ListNestedBlock{
				Description: "List of voice port profile configurations",
				NestedObject: schema.NestedBlockObject{
					Attributes: map[string]schema.Attribute{
						"voice_port_num_voice_port_profiles": schema.StringAttribute{
							Description: "Voice Port Profile for Voice Port",
							Optional:    true,
							Computed:    true,
						},
						"voice_port_num_voice_port_profiles_ref_type_": schema.StringAttribute{
							Description: "Object type for voice_port_num_voice_port_profiles field",
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

	var minState verityBundleResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if bundleData, exists := bulkMgr.GetResourceResponse("bundle", name); exists {
			state := populateBundleState(ctx, minState, bundleData, r.provCtx.mode)
			preserveBundlePortNames(&state, &plan)
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

	if !readResp.Diagnostics.HasError() {
		var readState verityBundleResourceModel
		readResp.State.Get(ctx, &readState)
		preserveBundlePortNames(&readState, &plan)
		resp.State.Set(ctx, &readState)
	}
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
	priorState := state // save prior state to preserve reference-only fields

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if bundleData, exists := r.bulkOpsMgr.GetResourceResponse("bundle", bundleName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached bundle data for %s from recent operation", bundleName))
			state = populateBundleState(ctx, state, bundleData, r.provCtx.mode)
			preserveBundlePortNames(&state, &priorState)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

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

	state = populateBundleState(ctx, state, bundleMap, r.provCtx.mode)
	preserveBundlePortNames(&state, &priorState)
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

	var minState verityBundleResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Try to use cached response from bulk operation to populate state with API values
	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if bundleData, exists := bulkMgr.GetResourceResponse("bundle", name); exists {
			newState := populateBundleState(ctx, minState, bundleData, r.provCtx.mode)
			preserveBundlePortNames(&newState, &plan)
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

	if !readResp.Diagnostics.HasError() {
		var readState verityBundleResourceModel
		readResp.State.Get(ctx, &readState)
		preserveBundlePortNames(&readState, &plan)
		resp.State.Set(ctx, &readState)
	}
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

func populateBundleState(ctx context.Context, state verityBundleResourceModel, data map[string]interface{}, mode string) verityBundleResourceModel {
	const resourceType = bundleResourceType

	state.Name = utils.MapStringFromAPI(data["name"])

	// Boolean fields
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)

	// String fields
	state.DeviceSettings = utils.MapStringWithMode(data, "device_settings", resourceType, mode)
	state.DeviceSettingsRefType = utils.MapStringWithMode(data, "device_settings_ref_type_", resourceType, mode)
	state.CliCommands = utils.MapStringWithMode(data, "cli_commands", resourceType, mode)
	state.Protocol = utils.MapStringWithMode(data, "protocol", resourceType, mode)
	state.DiagnosticsProfile = utils.MapStringWithMode(data, "diagnostics_profile", resourceType, mode)
	state.DiagnosticsProfileRefType = utils.MapStringWithMode(data, "diagnostics_profile_ref_type_", resourceType, mode)
	state.DeviceVoiceSettings = utils.MapStringWithMode(data, "device_voice_settings", resourceType, mode)
	state.DeviceVoiceSettingsRefType = utils.MapStringWithMode(data, "device_voice_settings_ref_type_", resourceType, mode)

	// Handle object_properties block
	if utils.FieldAppliesToMode(resourceType, "object_properties", mode) {
		if objProps, ok := data["object_properties"].(map[string]interface{}); ok {
			objPropsModel := verityBundleObjectPropertiesModel{
				IsForSwitch: utils.MapBoolWithModeNested(objProps, "is_for_switch", resourceType, "object_properties.is_for_switch", mode),
				Group:       utils.MapStringWithModeNested(objProps, "group", resourceType, "object_properties.group", mode),
				IsPublic:    utils.MapBoolWithModeNested(objProps, "is_public", resourceType, "object_properties.is_public", mode),
			}
			state.ObjectProperties = []verityBundleObjectPropertiesModel{objPropsModel}
		} else {
			state.ObjectProperties = nil
		}
	} else {
		state.ObjectProperties = nil
	}

	// Handle eth_port_paths array
	if utils.FieldAppliesToMode(resourceType, "eth_port_paths", mode) {
		if pathsData, ok := data["eth_port_paths"].([]interface{}); ok && len(pathsData) > 0 {
			var ethPortPaths []ethPortPathsModel
			for _, p := range pathsData {
				pathItem, ok := p.(map[string]interface{})
				if !ok {
					continue
				}
				pathModel := ethPortPathsModel{
					EthPortNumEthPortProfile:                               utils.MapStringWithModeNested(pathItem, "eth_port_num_eth_port_profile", resourceType, "eth_port_paths.eth_port_num_eth_port_profile", mode),
					EthPortNumEthPortProfileRefType:                        utils.MapStringWithModeNested(pathItem, "eth_port_num_eth_port_profile_ref_type_", resourceType, "eth_port_paths.eth_port_num_eth_port_profile_ref_type_", mode),
					EthPortNumEthPortSettings:                              utils.MapStringWithModeNested(pathItem, "eth_port_num_eth_port_settings", resourceType, "eth_port_paths.eth_port_num_eth_port_settings", mode),
					EthPortNumEthPortSettingsRefType:                       utils.MapStringWithModeNested(pathItem, "eth_port_num_eth_port_settings_ref_type_", resourceType, "eth_port_paths.eth_port_num_eth_port_settings_ref_type_", mode),
					EthPortNumGatewayProfile:                               utils.MapStringWithModeNested(pathItem, "eth_port_num_gateway_profile", resourceType, "eth_port_paths.eth_port_num_gateway_profile", mode),
					EthPortNumGatewayProfileRefType:                        utils.MapStringWithModeNested(pathItem, "eth_port_num_gateway_profile_ref_type_", resourceType, "eth_port_paths.eth_port_num_gateway_profile_ref_type_", mode),
					DiagnosticsPortProfileNumDiagnosticsPortProfile:        utils.MapStringWithModeNested(pathItem, "diagnostics_port_profile_num_diagnostics_port_profile", resourceType, "eth_port_paths.diagnostics_port_profile_num_diagnostics_port_profile", mode),
					DiagnosticsPortProfileNumDiagnosticsPortProfileRefType: utils.MapStringWithModeNested(pathItem, "diagnostics_port_profile_num_diagnostics_port_profile_ref_type_", resourceType, "eth_port_paths.diagnostics_port_profile_num_diagnostics_port_profile_ref_type_", mode),
					PortName: utils.MapStringWithModeNested(pathItem, "port_name", resourceType, "eth_port_paths.port_name", mode),
					Index:    utils.MapInt64WithModeNested(pathItem, "index", resourceType, "eth_port_paths.index", mode),
				}
				ethPortPaths = append(ethPortPaths, pathModel)
			}
			state.EthPortPaths = ethPortPaths
		} else {
			state.EthPortPaths = nil
		}
	} else {
		state.EthPortPaths = nil
	}

	// Handle user_services array
	if utils.FieldAppliesToMode(resourceType, "user_services", mode) {
		if servicesData, ok := data["user_services"].([]interface{}); ok && len(servicesData) > 0 {
			var userServices []userServicesModel
			for _, s := range servicesData {
				serviceItem, ok := s.(map[string]interface{})
				if !ok {
					continue
				}
				serviceModel := userServicesModel{
					RowAppEnable:                  utils.MapBoolWithModeNested(serviceItem, "row_app_enable", resourceType, "user_services.row_app_enable", mode),
					RowAppConnectedService:        utils.MapStringWithModeNested(serviceItem, "row_app_connected_service", resourceType, "user_services.row_app_connected_service", mode),
					RowAppConnectedServiceRefType: utils.MapStringWithModeNested(serviceItem, "row_app_connected_service_ref_type_", resourceType, "user_services.row_app_connected_service_ref_type_", mode),
					RowAppCliCommands:             utils.MapStringWithModeNested(serviceItem, "row_app_cli_commands", resourceType, "user_services.row_app_cli_commands", mode),
					RowIpMask:                     utils.MapStringWithModeNested(serviceItem, "row_ip_mask", resourceType, "user_services.row_ip_mask", mode),
					Index:                         utils.MapInt64WithModeNested(serviceItem, "index", resourceType, "user_services.index", mode),
				}
				userServices = append(userServices, serviceModel)
			}
			state.UserServices = userServices
		} else {
			state.UserServices = nil
		}
	} else {
		state.UserServices = nil
	}

	// Handle voice_port_profile_paths array
	if utils.FieldAppliesToMode(resourceType, "voice_port_profile_paths", mode) {
		if pathsData, ok := data["voice_port_profile_paths"].([]interface{}); ok && len(pathsData) > 0 {
			var voicePortProfilePaths []voicePortProfilePathsModel
			for _, p := range pathsData {
				pathItem, ok := p.(map[string]interface{})
				if !ok {
					continue
				}
				pathModel := voicePortProfilePathsModel{
					VoicePortNumVoicePortProfiles:        utils.MapStringWithModeNested(pathItem, "voice_port_num_voice_port_profiles", resourceType, "voice_port_profile_paths.voice_port_num_voice_port_profiles", mode),
					VoicePortNumVoicePortProfilesRefType: utils.MapStringWithModeNested(pathItem, "voice_port_num_voice_port_profiles_ref_type_", resourceType, "voice_port_profile_paths.voice_port_num_voice_port_profiles_ref_type_", mode),
					Index:                                utils.MapInt64WithModeNested(pathItem, "index", resourceType, "voice_port_profile_paths.index", mode),
				}
				voicePortProfilePaths = append(voicePortProfilePaths, pathModel)
			}
			state.VoicePortProfilePaths = voicePortProfilePaths
		} else {
			state.VoicePortProfilePaths = nil
		}
	} else {
		state.VoicePortProfilePaths = nil
	}

	return state
}

func (r *verityBundleResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// =========================================================================
	// Skip if deleting
	// =========================================================================
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityBundleResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Mode-aware field nullification
	// Set fields that don't apply to current mode to null to prevent
	// "known after apply" messages for irrelevant fields.
	// =========================================================================
	const resourceType = bundleResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyStrings(
		"device_settings", "device_settings_ref_type_",
		"cli_commands", "protocol",
		"diagnostics_profile", "diagnostics_profile_ref_type_",
		"device_voice_settings", "device_voice_settings_ref_type_",
	)

	nullifier.NullifyBools(
		"enable",
	)

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "object_properties",
		ItemCount:    len(plan.ObjectProperties),
		StringFields: []string{"group"},
		BoolFields:   []string{"is_for_switch", "is_public"},
	})

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName: "eth_port_paths",
		ItemCount: len(plan.EthPortPaths),
		StringFields: []string{
			"eth_port_num_eth_port_profile", "eth_port_num_eth_port_profile_ref_type_",
			"eth_port_num_eth_port_settings", "eth_port_num_eth_port_settings_ref_type_",
			"eth_port_num_gateway_profile", "eth_port_num_gateway_profile_ref_type_",
			"diagnostics_port_profile_num_diagnostics_port_profile", "diagnostics_port_profile_num_diagnostics_port_profile_ref_type_",
			"port_name",
		},
	})

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "user_services",
		ItemCount:    len(plan.UserServices),
		StringFields: []string{"row_app_connected_service", "row_app_connected_service_ref_type_", "row_app_cli_commands", "row_ip_mask"},
		BoolFields:   []string{"row_app_enable"},
	})

	nullifier.NullifyNestedBlockFields(utils.NestedBlockFieldConfig{
		BlockName:    "voice_port_profile_paths",
		ItemCount:    len(plan.VoicePortProfilePaths),
		StringFields: []string{"voice_port_num_voice_port_profiles", "voice_port_num_voice_port_profiles_ref_type_"},
	})
}

// preserveBundlePortNames copies port_name values from a reference source (plan or prior state)
// into the populated state. The API documents port_name as "reference only" – it accepts the value on PUT
// but never persists or returns it, so GET always gives back "".
func preserveBundlePortNames(state *verityBundleResourceModel, ref *verityBundleResourceModel) {
	if ref == nil || len(ref.EthPortPaths) == 0 || len(state.EthPortPaths) == 0 {
		return
	}

	for i := range state.EthPortPaths {
		if i >= len(ref.EthPortPaths) {
			break
		}

		refPortName := ref.EthPortPaths[i].PortName

		if !refPortName.IsNull() && !refPortName.IsUnknown() && state.EthPortPaths[i].PortName.ValueString() == "" {
			state.EthPortPaths[i].PortName = refPortName
		}
	}
}
