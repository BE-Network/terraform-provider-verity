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
	_ resource.Resource                = &verityDeviceControllerResource{}
	_ resource.ResourceWithConfigure   = &verityDeviceControllerResource{}
	_ resource.ResourceWithImportState = &verityDeviceControllerResource{}
)

func NewVerityDeviceControllerResource() resource.Resource {
	return &verityDeviceControllerResource{}
}

type verityDeviceControllerResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *utils.BulkOperationManager
	notifyOperationAdded func()
}

type verityDeviceControllerResourceModel struct {
	Name                      types.String `tfsdk:"name"`
	Enable                    types.Bool   `tfsdk:"enable"`
	IpSource                  types.String `tfsdk:"ip_source"`
	ControllerIpAndMask       types.String `tfsdk:"controller_ip_and_mask"`
	Gateway                   types.String `tfsdk:"gateway"`
	SwitchIpAndMask           types.String `tfsdk:"switch_ip_and_mask"`
	SwitchGateway             types.String `tfsdk:"switch_gateway"`
	CommType                  types.String `tfsdk:"comm_type"`
	SnmpCommunityString       types.String `tfsdk:"snmp_community_string"`
	UplinkPort                types.String `tfsdk:"uplink_port"`
	LldpSearchString          types.String `tfsdk:"lldp_search_string"`
	ZtpIdentification         types.String `tfsdk:"ztp_identification"`
	LocatedBy                 types.String `tfsdk:"located_by"`
	PowerState                types.String `tfsdk:"power_state"`
	CommunicationMode         types.String `tfsdk:"communication_mode"`
	CliAccessMode             types.String `tfsdk:"cli_access_mode"`
	Username                  types.String `tfsdk:"username"`
	Password                  types.String `tfsdk:"password"`
	EnablePassword            types.String `tfsdk:"enable_password"`
	SshKeyOrPassword          types.String `tfsdk:"ssh_key_or_password"`
	ManagedOnNativeVlan       types.Bool   `tfsdk:"managed_on_native_vlan"`
	Sdlc                      types.String `tfsdk:"sdlc"`
	Switchpoint               types.String `tfsdk:"switchpoint"`
	SwitchpointRefType        types.String `tfsdk:"switchpoint_ref_type_"`
	SecurityType              types.String `tfsdk:"security_type"`
	Snmpv3Username            types.String `tfsdk:"snmpv3_username"`
	AuthenticationProtocol    types.String `tfsdk:"authentication_protocol"`
	Passphrase                types.String `tfsdk:"passphrase"`
	PrivateProtocol           types.String `tfsdk:"private_protocol"`
	PrivatePassword           types.String `tfsdk:"private_password"`
	PasswordEncrypted         types.String `tfsdk:"password_encrypted"`
	EnablePasswordEncrypted   types.String `tfsdk:"enable_password_encrypted"`
	SshKeyOrPasswordEncrypted types.String `tfsdk:"ssh_key_or_password_encrypted"`
	PassphraseEncrypted       types.String `tfsdk:"passphrase_encrypted"`
	PrivatePasswordEncrypted  types.String `tfsdk:"private_password_encrypted"`
	DeviceManagedAs           types.String `tfsdk:"device_managed_as"`
	Switch                    types.String `tfsdk:"switch"`
	SwitchRefType             types.String `tfsdk:"switch_ref_type_"`
	ConnectionService         types.String `tfsdk:"connection_service"`
	ConnectionServiceRefType  types.String `tfsdk:"connection_service_ref_type_"`
	Port                      types.String `tfsdk:"port"`
	SfpMacAddressOrSn         types.String `tfsdk:"sfp_mac_address_or_sn"`
	UsesTaggedPackets         types.Bool   `tfsdk:"uses_tagged_packets"`
}

func (r *verityDeviceControllerResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device_controller"
}

func (r *verityDeviceControllerResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *verityDeviceControllerResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a Verity Device Controller",
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
			"ip_source": schema.StringAttribute{
				Description: "IP Source",
				Optional:    true,
			},
			"controller_ip_and_mask": schema.StringAttribute{
				Description: "Controller IP and Mask",
				Optional:    true,
			},
			"gateway": schema.StringAttribute{
				Description: "Gateway",
				Optional:    true,
			},
			"switch_ip_and_mask": schema.StringAttribute{
				Description: "Switch IP and Mask",
				Optional:    true,
			},
			"switch_gateway": schema.StringAttribute{
				Description: "Gateway of Managed Device",
				Optional:    true,
			},
			"comm_type": schema.StringAttribute{
				Description: "Comm Type",
				Optional:    true,
			},
			"snmp_community_string": schema.StringAttribute{
				Description: "Comm Credentials",
				Optional:    true,
			},
			"uplink_port": schema.StringAttribute{
				Description: "Uplink Port of Managed Device",
				Optional:    true,
			},
			"lldp_search_string": schema.StringAttribute{
				Description: "Optional unless Located By is \"LLDP\" or Device managed as \"Active SFP\". Must be either the chassis-id or the hostname of the LLDP from the managed device. Used to detect connections between managed devices. If blank, the chassis-id detected by the Device Controller via SNMP/CLI is used",
				Optional:    true,
			},
			"ztp_identification": schema.StringAttribute{
				Description: "Service Tag or Serial Number to identify device for Zero Touch Provisioning",
				Optional:    true,
			},
			"located_by": schema.StringAttribute{
				Description: "Controls how the system locates this Device within its LAN",
				Optional:    true,
			},
			"power_state": schema.StringAttribute{
				Description: "Power state of Switch Controller",
				Optional:    true,
			},
			"communication_mode": schema.StringAttribute{
				Description: "Communication Mode",
				Optional:    true,
			},
			"cli_access_mode": schema.StringAttribute{
				Description: "CLI Access Mode",
				Optional:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username",
				Optional:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password",
				Optional:    true,
			},
			"enable_password": schema.StringAttribute{
				Description: "Enable Password - to enable privileged CLI operations",
				Optional:    true,
			},
			"ssh_key_or_password": schema.StringAttribute{
				Description: "SSH Key or Password",
				Optional:    true,
			},
			"managed_on_native_vlan": schema.BoolAttribute{
				Description: "Managed on native VLAN",
				Optional:    true,
			},
			"sdlc": schema.StringAttribute{
				Description: "SDLC that Device Controller belongs to",
				Optional:    true,
			},
			"switchpoint": schema.StringAttribute{
				Description: "Endpoint reference",
				Optional:    true,
			},
			"switchpoint_ref_type_": schema.StringAttribute{
				Description: "Object type for switchpoint field",
				Optional:    true,
			},
			"security_type": schema.StringAttribute{
				Description: "Security level",
				Optional:    true,
			},
			"snmpv3_username": schema.StringAttribute{
				Description: "SNMPv3 Username",
				Optional:    true,
			},
			"authentication_protocol": schema.StringAttribute{
				Description: "Authentication Protocol",
				Optional:    true,
			},
			"passphrase": schema.StringAttribute{
				Description: "Passphrase",
				Optional:    true,
			},
			"private_protocol": schema.StringAttribute{
				Description: "Private Protocol",
				Optional:    true,
			},
			"private_password": schema.StringAttribute{
				Description: "Private Password",
				Optional:    true,
			},
			"password_encrypted": schema.StringAttribute{
				Description: "Encrypted Password",
				Optional:    true,
			},
			"enable_password_encrypted": schema.StringAttribute{
				Description: "Encrypted Enable Password - to enable privileged CLI operations",
				Optional:    true,
			},
			"ssh_key_or_password_encrypted": schema.StringAttribute{
				Description: "Encrypted SSH Key or Password",
				Optional:    true,
			},
			"passphrase_encrypted": schema.StringAttribute{
				Description: "Encrypted Passphrase",
				Optional:    true,
			},
			"private_password_encrypted": schema.StringAttribute{
				Description: "Encrypted Private Password",
				Optional:    true,
			},
			"device_managed_as": schema.StringAttribute{
				Description: "Device managed as",
				Optional:    true,
			},
			"switch": schema.StringAttribute{
				Description: "Endpoint locating the Switch to be controlled",
				Optional:    true,
			},
			"switch_ref_type_": schema.StringAttribute{
				Description: "Object type for switch field",
				Optional:    true,
			},
			"connection_service": schema.StringAttribute{
				Description: "Connect a Service",
				Optional:    true,
			},
			"connection_service_ref_type_": schema.StringAttribute{
				Description: "Object type for connection_service field",
				Optional:    true,
			},
			"port": schema.StringAttribute{
				Description: "Port locating the Switch to be controlled",
				Optional:    true,
			},
			"sfp_mac_address_or_sn": schema.StringAttribute{
				Description: "SFP MAC Address or SN",
				Optional:    true,
			},
			"uses_tagged_packets": schema.BoolAttribute{
				Description: "Indicates if the direct interface expects tagged or untagged packets",
				Optional:    true,
			},
		},
	}
}

func (r *verityDeviceControllerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan verityDeviceControllerResourceModel
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
	deviceControllerProps := &openapi.DevicecontrollersPutRequestDeviceControllerValue{
		Name: openapi.PtrString(name),
	}

	// Handle string fields
	utils.SetStringFields([]utils.StringFieldMapping{
		{FieldName: "IpSource", APIField: &deviceControllerProps.IpSource, TFValue: plan.IpSource},
		{FieldName: "ControllerIpAndMask", APIField: &deviceControllerProps.ControllerIpAndMask, TFValue: plan.ControllerIpAndMask},
		{FieldName: "Gateway", APIField: &deviceControllerProps.Gateway, TFValue: plan.Gateway},
		{FieldName: "SwitchIpAndMask", APIField: &deviceControllerProps.SwitchIpAndMask, TFValue: plan.SwitchIpAndMask},
		{FieldName: "SwitchGateway", APIField: &deviceControllerProps.SwitchGateway, TFValue: plan.SwitchGateway},
		{FieldName: "CommType", APIField: &deviceControllerProps.CommType, TFValue: plan.CommType},
		{FieldName: "SnmpCommunityString", APIField: &deviceControllerProps.SnmpCommunityString, TFValue: plan.SnmpCommunityString},
		{FieldName: "UplinkPort", APIField: &deviceControllerProps.UplinkPort, TFValue: plan.UplinkPort},
		{FieldName: "LldpSearchString", APIField: &deviceControllerProps.LldpSearchString, TFValue: plan.LldpSearchString},
		{FieldName: "ZtpIdentification", APIField: &deviceControllerProps.ZtpIdentification, TFValue: plan.ZtpIdentification},
		{FieldName: "LocatedBy", APIField: &deviceControllerProps.LocatedBy, TFValue: plan.LocatedBy},
		{FieldName: "PowerState", APIField: &deviceControllerProps.PowerState, TFValue: plan.PowerState},
		{FieldName: "CommunicationMode", APIField: &deviceControllerProps.CommunicationMode, TFValue: plan.CommunicationMode},
		{FieldName: "CliAccessMode", APIField: &deviceControllerProps.CliAccessMode, TFValue: plan.CliAccessMode},
		{FieldName: "Username", APIField: &deviceControllerProps.Username, TFValue: plan.Username},
		{FieldName: "Password", APIField: &deviceControllerProps.Password, TFValue: plan.Password},
		{FieldName: "EnablePassword", APIField: &deviceControllerProps.EnablePassword, TFValue: plan.EnablePassword},
		{FieldName: "SshKeyOrPassword", APIField: &deviceControllerProps.SshKeyOrPassword, TFValue: plan.SshKeyOrPassword},
		{FieldName: "Sdlc", APIField: &deviceControllerProps.Sdlc, TFValue: plan.Sdlc},
		{FieldName: "Switchpoint", APIField: &deviceControllerProps.Switchpoint, TFValue: plan.Switchpoint},
		{FieldName: "SwitchpointRefType", APIField: &deviceControllerProps.SwitchpointRefType, TFValue: plan.SwitchpointRefType},
		{FieldName: "SecurityType", APIField: &deviceControllerProps.SecurityType, TFValue: plan.SecurityType},
		{FieldName: "Snmpv3Username", APIField: &deviceControllerProps.Snmpv3Username, TFValue: plan.Snmpv3Username},
		{FieldName: "AuthenticationProtocol", APIField: &deviceControllerProps.AuthenticationProtocol, TFValue: plan.AuthenticationProtocol},
		{FieldName: "Passphrase", APIField: &deviceControllerProps.Passphrase, TFValue: plan.Passphrase},
		{FieldName: "PrivateProtocol", APIField: &deviceControllerProps.PrivateProtocol, TFValue: plan.PrivateProtocol},
		{FieldName: "PrivatePassword", APIField: &deviceControllerProps.PrivatePassword, TFValue: plan.PrivatePassword},
		{FieldName: "PasswordEncrypted", APIField: &deviceControllerProps.PasswordEncrypted, TFValue: plan.PasswordEncrypted},
		{FieldName: "EnablePasswordEncrypted", APIField: &deviceControllerProps.EnablePasswordEncrypted, TFValue: plan.EnablePasswordEncrypted},
		{FieldName: "SshKeyOrPasswordEncrypted", APIField: &deviceControllerProps.SshKeyOrPasswordEncrypted, TFValue: plan.SshKeyOrPasswordEncrypted},
		{FieldName: "PassphraseEncrypted", APIField: &deviceControllerProps.PassphraseEncrypted, TFValue: plan.PassphraseEncrypted},
		{FieldName: "PrivatePasswordEncrypted", APIField: &deviceControllerProps.PrivatePasswordEncrypted, TFValue: plan.PrivatePasswordEncrypted},
		{FieldName: "DeviceManagedAs", APIField: &deviceControllerProps.DeviceManagedAs, TFValue: plan.DeviceManagedAs},
		{FieldName: "Switch", APIField: &deviceControllerProps.Switch, TFValue: plan.Switch},
		{FieldName: "SwitchRefType", APIField: &deviceControllerProps.SwitchRefType, TFValue: plan.SwitchRefType},
		{FieldName: "ConnectionService", APIField: &deviceControllerProps.ConnectionService, TFValue: plan.ConnectionService},
		{FieldName: "ConnectionServiceRefType", APIField: &deviceControllerProps.ConnectionServiceRefType, TFValue: plan.ConnectionServiceRefType},
		{FieldName: "Port", APIField: &deviceControllerProps.Port, TFValue: plan.Port},
		{FieldName: "SfpMacAddressOrSn", APIField: &deviceControllerProps.SfpMacAddressOrSn, TFValue: plan.SfpMacAddressOrSn},
	})

	// Handle boolean fields
	utils.SetBoolFields([]utils.BoolFieldMapping{
		{FieldName: "Enable", APIField: &deviceControllerProps.Enable, TFValue: plan.Enable},
		{FieldName: "ManagedOnNativeVlan", APIField: &deviceControllerProps.ManagedOnNativeVlan, TFValue: plan.ManagedOnNativeVlan},
		{FieldName: "UsesTaggedPackets", APIField: &deviceControllerProps.UsesTaggedPackets, TFValue: plan.UsesTaggedPackets},
	})

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "device_controller", name, *deviceControllerProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Device Controller %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "device_controllers")

	plan.Name = types.StringValue(name)
	resp.State.Set(ctx, plan)
}

func (r *verityDeviceControllerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state verityDeviceControllerResourceModel
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

	deviceControllerName := state.Name.ValueString()

	if r.bulkOpsMgr != nil && r.bulkOpsMgr.HasPendingOrRecentOperations("device_controller") {
		tflog.Info(ctx, fmt.Sprintf("Skipping device controller %s verification – trusting recent successful API operation", deviceControllerName))
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Fetching device controllers for verification of %s", deviceControllerName))

	type DeviceControllersResponse struct {
		DeviceController map[string]interface{} `json:"device_controller"`
	}

	result, err := utils.FetchResourceWithRetry(ctx, r.provCtx, "device_controllers", deviceControllerName,
		func() (DeviceControllersResponse, error) {
			tflog.Debug(ctx, "Making API call to fetch device controllers")
			respAPI, err := r.client.DeviceControllersAPI.DevicecontrollersGet(ctx).Execute()
			if err != nil {
				return DeviceControllersResponse{}, fmt.Errorf("error reading device controllers: %v", err)
			}
			defer respAPI.Body.Close()

			var res DeviceControllersResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return DeviceControllersResponse{}, fmt.Errorf("failed to decode device controllers response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d device controllers", len(res.DeviceController)))
			return res, nil
		},
		getCachedResponse,
	)

	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Device Controller %s", deviceControllerName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for device controller with name: %s", deviceControllerName))

	deviceControllerData, actualAPIName, exists := utils.FindResourceByAPIName(
		result.DeviceController,
		deviceControllerName,
		func(data interface{}) (string, bool) {
			if deviceController, ok := data.(map[string]interface{}); ok {
				if name, ok := deviceController["name"].(string); ok {
					return name, true
				}
			}
			return "", false
		},
	)

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Device Controller with name '%s' not found in API response", deviceControllerName))
		resp.State.RemoveResource(ctx)
		return
	}

	deviceControllerMap, ok := deviceControllerData.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Invalid Device Controller Data",
			fmt.Sprintf("Device Controller data is not in expected format for %s", deviceControllerName),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Found device controller '%s' under API key '%s'", deviceControllerName, actualAPIName))

	state.Name = utils.MapStringFromAPI(deviceControllerMap["name"])

	// Map string fields
	stringFieldMappings := map[string]*types.String{
		"ip_source":                     &state.IpSource,
		"controller_ip_and_mask":        &state.ControllerIpAndMask,
		"gateway":                       &state.Gateway,
		"switch_ip_and_mask":            &state.SwitchIpAndMask,
		"switch_gateway":                &state.SwitchGateway,
		"comm_type":                     &state.CommType,
		"snmp_community_string":         &state.SnmpCommunityString,
		"uplink_port":                   &state.UplinkPort,
		"lldp_search_string":            &state.LldpSearchString,
		"ztp_identification":            &state.ZtpIdentification,
		"located_by":                    &state.LocatedBy,
		"power_state":                   &state.PowerState,
		"communication_mode":            &state.CommunicationMode,
		"cli_access_mode":               &state.CliAccessMode,
		"username":                      &state.Username,
		"password":                      &state.Password,
		"enable_password":               &state.EnablePassword,
		"ssh_key_or_password":           &state.SshKeyOrPassword,
		"sdlc":                          &state.Sdlc,
		"switchpoint":                   &state.Switchpoint,
		"switchpoint_ref_type_":         &state.SwitchpointRefType,
		"security_type":                 &state.SecurityType,
		"snmpv3_username":               &state.Snmpv3Username,
		"authentication_protocol":       &state.AuthenticationProtocol,
		"passphrase":                    &state.Passphrase,
		"private_protocol":              &state.PrivateProtocol,
		"private_password":              &state.PrivatePassword,
		"password_encrypted":            &state.PasswordEncrypted,
		"enable_password_encrypted":     &state.EnablePasswordEncrypted,
		"ssh_key_or_password_encrypted": &state.SshKeyOrPasswordEncrypted,
		"passphrase_encrypted":          &state.PassphraseEncrypted,
		"private_password_encrypted":    &state.PrivatePasswordEncrypted,
		"device_managed_as":             &state.DeviceManagedAs,
		"switch":                        &state.Switch,
		"switch_ref_type_":              &state.SwitchRefType,
		"connection_service":            &state.ConnectionService,
		"connection_service_ref_type_":  &state.ConnectionServiceRefType,
		"port":                          &state.Port,
		"sfp_mac_address_or_sn":         &state.SfpMacAddressOrSn,
	}

	for apiKey, stateField := range stringFieldMappings {
		*stateField = utils.MapStringFromAPI(deviceControllerMap[apiKey])
	}

	// Map boolean fields
	boolFieldMappings := map[string]*types.Bool{
		"enable":                 &state.Enable,
		"managed_on_native_vlan": &state.ManagedOnNativeVlan,
		"uses_tagged_packets":    &state.UsesTaggedPackets,
	}

	for apiKey, stateField := range boolFieldMappings {
		*stateField = utils.MapBoolFromAPI(deviceControllerMap[apiKey])
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityDeviceControllerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityDeviceControllerResourceModel

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
	deviceControllerProps := openapi.DevicecontrollersPutRequestDeviceControllerValue{}
	hasChanges := false

	// Handle string field changes
	utils.CompareAndSetStringField(plan.Name, state.Name, func(v *string) { deviceControllerProps.Name = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.IpSource, state.IpSource, func(v *string) { deviceControllerProps.IpSource = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.ControllerIpAndMask, state.ControllerIpAndMask, func(v *string) { deviceControllerProps.ControllerIpAndMask = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Gateway, state.Gateway, func(v *string) { deviceControllerProps.Gateway = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.SwitchIpAndMask, state.SwitchIpAndMask, func(v *string) { deviceControllerProps.SwitchIpAndMask = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.SwitchGateway, state.SwitchGateway, func(v *string) { deviceControllerProps.SwitchGateway = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.CommType, state.CommType, func(v *string) { deviceControllerProps.CommType = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.SnmpCommunityString, state.SnmpCommunityString, func(v *string) { deviceControllerProps.SnmpCommunityString = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.UplinkPort, state.UplinkPort, func(v *string) { deviceControllerProps.UplinkPort = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.LldpSearchString, state.LldpSearchString, func(v *string) { deviceControllerProps.LldpSearchString = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.ZtpIdentification, state.ZtpIdentification, func(v *string) { deviceControllerProps.ZtpIdentification = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.LocatedBy, state.LocatedBy, func(v *string) { deviceControllerProps.LocatedBy = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.PowerState, state.PowerState, func(v *string) { deviceControllerProps.PowerState = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.CommunicationMode, state.CommunicationMode, func(v *string) { deviceControllerProps.CommunicationMode = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.CliAccessMode, state.CliAccessMode, func(v *string) { deviceControllerProps.CliAccessMode = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Username, state.Username, func(v *string) { deviceControllerProps.Username = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Password, state.Password, func(v *string) { deviceControllerProps.Password = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.EnablePassword, state.EnablePassword, func(v *string) { deviceControllerProps.EnablePassword = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.SshKeyOrPassword, state.SshKeyOrPassword, func(v *string) { deviceControllerProps.SshKeyOrPassword = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Sdlc, state.Sdlc, func(v *string) { deviceControllerProps.Sdlc = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.SecurityType, state.SecurityType, func(v *string) { deviceControllerProps.SecurityType = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Snmpv3Username, state.Snmpv3Username, func(v *string) { deviceControllerProps.Snmpv3Username = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.AuthenticationProtocol, state.AuthenticationProtocol, func(v *string) { deviceControllerProps.AuthenticationProtocol = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Passphrase, state.Passphrase, func(v *string) { deviceControllerProps.Passphrase = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.PrivateProtocol, state.PrivateProtocol, func(v *string) { deviceControllerProps.PrivateProtocol = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.PrivatePassword, state.PrivatePassword, func(v *string) { deviceControllerProps.PrivatePassword = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.PasswordEncrypted, state.PasswordEncrypted, func(v *string) { deviceControllerProps.PasswordEncrypted = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.EnablePasswordEncrypted, state.EnablePasswordEncrypted, func(v *string) { deviceControllerProps.EnablePasswordEncrypted = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.SshKeyOrPasswordEncrypted, state.SshKeyOrPasswordEncrypted, func(v *string) { deviceControllerProps.SshKeyOrPasswordEncrypted = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.PassphraseEncrypted, state.PassphraseEncrypted, func(v *string) { deviceControllerProps.PassphraseEncrypted = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.PrivatePasswordEncrypted, state.PrivatePasswordEncrypted, func(v *string) { deviceControllerProps.PrivatePasswordEncrypted = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.DeviceManagedAs, state.DeviceManagedAs, func(v *string) { deviceControllerProps.DeviceManagedAs = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.Port, state.Port, func(v *string) { deviceControllerProps.Port = v }, &hasChanges)
	utils.CompareAndSetStringField(plan.SfpMacAddressOrSn, state.SfpMacAddressOrSn, func(v *string) { deviceControllerProps.SfpMacAddressOrSn = v }, &hasChanges)

	// Handle boolean field changes
	utils.CompareAndSetBoolField(plan.Enable, state.Enable, func(v *bool) { deviceControllerProps.Enable = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.ManagedOnNativeVlan, state.ManagedOnNativeVlan, func(v *bool) { deviceControllerProps.ManagedOnNativeVlan = v }, &hasChanges)
	utils.CompareAndSetBoolField(plan.UsesTaggedPackets, state.UsesTaggedPackets, func(v *bool) { deviceControllerProps.UsesTaggedPackets = v }, &hasChanges)

	// Handle Switchpoint and SwitchpointRefType using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.Switchpoint, state.Switchpoint, plan.SwitchpointRefType, state.SwitchpointRefType,
		func(v *string) { deviceControllerProps.Switchpoint = v },
		func(v *string) { deviceControllerProps.SwitchpointRefType = v },
		"switchpoint", "switchpoint_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	// Handle Switch and SwitchRefType using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.Switch, state.Switch, plan.SwitchRefType, state.SwitchRefType,
		func(v *string) { deviceControllerProps.Switch = v },
		func(v *string) { deviceControllerProps.SwitchRefType = v },
		"switch", "switch_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	// Handle ConnectionService and ConnectionServiceRefType using "One ref type supported" pattern
	if !utils.HandleOneRefTypeSupported(
		plan.ConnectionService, state.ConnectionService, plan.ConnectionServiceRefType, state.ConnectionServiceRefType,
		func(v *string) { deviceControllerProps.ConnectionService = v },
		func(v *string) { deviceControllerProps.ConnectionServiceRefType = v },
		"connection_service", "connection_service_ref_type_",
		&hasChanges,
		&resp.Diagnostics,
	) {
		return
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "device_controller", name, deviceControllerProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Device Controller %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "device_controllers")
	resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
}

func (r *verityDeviceControllerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state verityDeviceControllerResourceModel
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

	success := utils.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "device_controller", name, nil, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Device Controller %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "device_controllers")
	resp.State.RemoveResource(ctx)
}

func (r *verityDeviceControllerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}
