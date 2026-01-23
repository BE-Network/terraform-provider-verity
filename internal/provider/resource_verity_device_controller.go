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
	_ resource.Resource                = &verityDeviceControllerResource{}
	_ resource.ResourceWithConfigure   = &verityDeviceControllerResource{}
	_ resource.ResourceWithImportState = &verityDeviceControllerResource{}
	_ resource.ResourceWithModifyPlan  = &verityDeviceControllerResource{}
)

const deviceControllerResourceType = "devicecontrollers"

func NewVerityDeviceControllerResource() resource.Resource {
	return &verityDeviceControllerResource{}
}

type verityDeviceControllerResource struct {
	provCtx              *providerContext
	client               *openapi.APIClient
	bulkOpsMgr           *bulkops.Manager
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
				Computed:    true,
			},
			"ip_source": schema.StringAttribute{
				Description: "IP Source",
				Optional:    true,
				Computed:    true,
			},
			"controller_ip_and_mask": schema.StringAttribute{
				Description: "Controller IP and Mask",
				Optional:    true,
				Computed:    true,
			},
			"gateway": schema.StringAttribute{
				Description: "Gateway",
				Optional:    true,
				Computed:    true,
			},
			"switch_ip_and_mask": schema.StringAttribute{
				Description: "Switch IP and Mask",
				Optional:    true,
				Computed:    true,
			},
			"switch_gateway": schema.StringAttribute{
				Description: "Gateway of Managed Device",
				Optional:    true,
				Computed:    true,
			},
			"comm_type": schema.StringAttribute{
				Description: "Comm Type",
				Optional:    true,
				Computed:    true,
			},
			"snmp_community_string": schema.StringAttribute{
				Description: "Comm Credentials",
				Optional:    true,
				Computed:    true,
			},
			"uplink_port": schema.StringAttribute{
				Description: "Uplink Port of Managed Device",
				Optional:    true,
				Computed:    true,
			},
			"lldp_search_string": schema.StringAttribute{
				Description: "Optional unless Located By is \"LLDP\" or Device managed as \"Active SFP\". Must be either the chassis-id or the hostname of the LLDP from the managed device. Used to detect connections between managed devices. If blank, the chassis-id detected by the Device Controller via SNMP/CLI is used",
				Optional:    true,
				Computed:    true,
			},
			"ztp_identification": schema.StringAttribute{
				Description: "Service Tag or Serial Number to identify device for Zero Touch Provisioning",
				Optional:    true,
				Computed:    true,
			},
			"located_by": schema.StringAttribute{
				Description: "Controls how the system locates this Device within its LAN",
				Optional:    true,
				Computed:    true,
			},
			"power_state": schema.StringAttribute{
				Description: "Power state of Switch Controller",
				Optional:    true,
				Computed:    true,
			},
			"communication_mode": schema.StringAttribute{
				Description: "Communication Mode",
				Optional:    true,
				Computed:    true,
			},
			"cli_access_mode": schema.StringAttribute{
				Description: "CLI Access Mode",
				Optional:    true,
				Computed:    true,
			},
			"username": schema.StringAttribute{
				Description: "Username",
				Optional:    true,
				Computed:    true,
			},
			"password": schema.StringAttribute{
				Description: "Password",
				Optional:    true,
				Computed:    true,
			},
			"enable_password": schema.StringAttribute{
				Description: "Enable Password - to enable privileged CLI operations",
				Optional:    true,
				Computed:    true,
			},
			"ssh_key_or_password": schema.StringAttribute{
				Description: "SSH Key or Password",
				Optional:    true,
				Computed:    true,
			},
			"managed_on_native_vlan": schema.BoolAttribute{
				Description: "Managed on native VLAN",
				Optional:    true,
				Computed:    true,
			},
			"sdlc": schema.StringAttribute{
				Description: "SDLC that Device Controller belongs to",
				Optional:    true,
				Computed:    true,
			},
			"switchpoint": schema.StringAttribute{
				Description: "Endpoint reference",
				Optional:    true,
				Computed:    true,
			},
			"switchpoint_ref_type_": schema.StringAttribute{
				Description: "Object type for switchpoint field",
				Optional:    true,
				Computed:    true,
			},
			"security_type": schema.StringAttribute{
				Description: "Security level",
				Optional:    true,
				Computed:    true,
			},
			"snmpv3_username": schema.StringAttribute{
				Description: "SNMPv3 Username",
				Optional:    true,
				Computed:    true,
			},
			"authentication_protocol": schema.StringAttribute{
				Description: "Authentication Protocol",
				Optional:    true,
				Computed:    true,
			},
			"passphrase": schema.StringAttribute{
				Description: "Passphrase",
				Optional:    true,
				Computed:    true,
			},
			"private_protocol": schema.StringAttribute{
				Description: "Private Protocol",
				Optional:    true,
				Computed:    true,
			},
			"private_password": schema.StringAttribute{
				Description: "Private Password",
				Optional:    true,
				Computed:    true,
			},
			"password_encrypted": schema.StringAttribute{
				Description: "Encrypted Password",
				Optional:    true,
				Computed:    true,
			},
			"enable_password_encrypted": schema.StringAttribute{
				Description: "Encrypted Enable Password - to enable privileged CLI operations",
				Optional:    true,
				Computed:    true,
			},
			"ssh_key_or_password_encrypted": schema.StringAttribute{
				Description: "Encrypted SSH Key or Password",
				Optional:    true,
				Computed:    true,
			},
			"passphrase_encrypted": schema.StringAttribute{
				Description: "Encrypted Passphrase",
				Optional:    true,
				Computed:    true,
			},
			"private_password_encrypted": schema.StringAttribute{
				Description: "Encrypted Private Password",
				Optional:    true,
				Computed:    true,
			},
			"device_managed_as": schema.StringAttribute{
				Description: "Device managed as",
				Optional:    true,
				Computed:    true,
			},
			"switch": schema.StringAttribute{
				Description: "Endpoint locating the Switch to be controlled",
				Optional:    true,
				Computed:    true,
			},
			"switch_ref_type_": schema.StringAttribute{
				Description: "Object type for switch field",
				Optional:    true,
				Computed:    true,
			},
			"connection_service": schema.StringAttribute{
				Description: "Connect a Service",
				Optional:    true,
				Computed:    true,
			},
			"connection_service_ref_type_": schema.StringAttribute{
				Description: "Object type for connection_service field",
				Optional:    true,
				Computed:    true,
			},
			"port": schema.StringAttribute{
				Description: "Port locating the Switch to be controlled",
				Optional:    true,
				Computed:    true,
			},
			"sfp_mac_address_or_sn": schema.StringAttribute{
				Description: "SFP MAC Address or SN",
				Optional:    true,
				Computed:    true,
			},
			"uses_tagged_packets": schema.BoolAttribute{
				Description: "Indicates if the direct interface expects tagged or untagged packets",
				Optional:    true,
				Computed:    true,
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "create", "device_controller", name, *deviceControllerProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Device Controller %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "device_controllers")

	var minState verityDeviceControllerResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if deviceControllerData, exists := bulkMgr.GetResourceResponse("device_controller", name); exists {
			state := populateDeviceControllerState(ctx, minState, deviceControllerData, r.provCtx.mode)
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

	// Check for cached data from recent operations first
	if r.bulkOpsMgr != nil {
		if deviceControllerData, exists := r.bulkOpsMgr.GetResourceResponse("device_controller", deviceControllerName); exists {
			tflog.Info(ctx, fmt.Sprintf("Using cached device controller data for %s from recent operation", deviceControllerName))
			state = populateDeviceControllerState(ctx, state, deviceControllerData, r.provCtx.mode)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

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

	state = populateDeviceControllerState(ctx, state, deviceControllerMap, r.provCtx.mode)
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "update", "device_controller", name, deviceControllerProps, &resp.Diagnostics)
	if !success {
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Device Controller %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "device_controllers")

	var minState verityDeviceControllerResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Try to use cached response from bulk operation to populate state with API values
	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if deviceControllerData, exists := bulkMgr.GetResourceResponse("device_controller", name); exists {
			newState := populateDeviceControllerState(ctx, minState, deviceControllerData, r.provCtx.mode)
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

	success := bulkops.ExecuteResourceOperation(ctx, r.bulkOpsMgr, r.notifyOperationAdded, "delete", "device_controller", name, nil, &resp.Diagnostics)
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

func populateDeviceControllerState(ctx context.Context, state verityDeviceControllerResourceModel, data map[string]interface{}, mode string) verityDeviceControllerResourceModel {
	const resourceType = deviceControllerResourceType

	state.Name = utils.MapStringFromAPI(data["name"])

	// Boolean fields
	state.Enable = utils.MapBoolWithMode(data, "enable", resourceType, mode)
	state.ManagedOnNativeVlan = utils.MapBoolWithMode(data, "managed_on_native_vlan", resourceType, mode)
	state.UsesTaggedPackets = utils.MapBoolWithMode(data, "uses_tagged_packets", resourceType, mode)

	// String fields
	state.IpSource = utils.MapStringWithMode(data, "ip_source", resourceType, mode)
	state.ControllerIpAndMask = utils.MapStringWithMode(data, "controller_ip_and_mask", resourceType, mode)
	state.Gateway = utils.MapStringWithMode(data, "gateway", resourceType, mode)
	state.SwitchIpAndMask = utils.MapStringWithMode(data, "switch_ip_and_mask", resourceType, mode)
	state.SwitchGateway = utils.MapStringWithMode(data, "switch_gateway", resourceType, mode)
	state.CommType = utils.MapStringWithMode(data, "comm_type", resourceType, mode)
	state.SnmpCommunityString = utils.MapStringWithMode(data, "snmp_community_string", resourceType, mode)
	state.UplinkPort = utils.MapStringWithMode(data, "uplink_port", resourceType, mode)
	state.LldpSearchString = utils.MapStringWithMode(data, "lldp_search_string", resourceType, mode)
	state.ZtpIdentification = utils.MapStringWithMode(data, "ztp_identification", resourceType, mode)
	state.LocatedBy = utils.MapStringWithMode(data, "located_by", resourceType, mode)
	state.PowerState = utils.MapStringWithMode(data, "power_state", resourceType, mode)
	state.CommunicationMode = utils.MapStringWithMode(data, "communication_mode", resourceType, mode)
	state.CliAccessMode = utils.MapStringWithMode(data, "cli_access_mode", resourceType, mode)
	state.Username = utils.MapStringWithMode(data, "username", resourceType, mode)
	state.Password = utils.MapStringWithMode(data, "password", resourceType, mode)
	state.EnablePassword = utils.MapStringWithMode(data, "enable_password", resourceType, mode)
	state.SshKeyOrPassword = utils.MapStringWithMode(data, "ssh_key_or_password", resourceType, mode)
	state.Sdlc = utils.MapStringWithMode(data, "sdlc", resourceType, mode)
	state.Switchpoint = utils.MapStringWithMode(data, "switchpoint", resourceType, mode)
	state.SwitchpointRefType = utils.MapStringWithMode(data, "switchpoint_ref_type_", resourceType, mode)
	state.SecurityType = utils.MapStringWithMode(data, "security_type", resourceType, mode)
	state.Snmpv3Username = utils.MapStringWithMode(data, "snmpv3_username", resourceType, mode)
	state.AuthenticationProtocol = utils.MapStringWithMode(data, "authentication_protocol", resourceType, mode)
	state.Passphrase = utils.MapStringWithMode(data, "passphrase", resourceType, mode)
	state.PrivateProtocol = utils.MapStringWithMode(data, "private_protocol", resourceType, mode)
	state.PrivatePassword = utils.MapStringWithMode(data, "private_password", resourceType, mode)
	state.PasswordEncrypted = utils.MapStringWithMode(data, "password_encrypted", resourceType, mode)
	state.EnablePasswordEncrypted = utils.MapStringWithMode(data, "enable_password_encrypted", resourceType, mode)
	state.SshKeyOrPasswordEncrypted = utils.MapStringWithMode(data, "ssh_key_or_password_encrypted", resourceType, mode)
	state.PassphraseEncrypted = utils.MapStringWithMode(data, "passphrase_encrypted", resourceType, mode)
	state.PrivatePasswordEncrypted = utils.MapStringWithMode(data, "private_password_encrypted", resourceType, mode)
	state.DeviceManagedAs = utils.MapStringWithMode(data, "device_managed_as", resourceType, mode)
	state.Switch = utils.MapStringWithMode(data, "switch", resourceType, mode)
	state.SwitchRefType = utils.MapStringWithMode(data, "switch_ref_type_", resourceType, mode)
	state.ConnectionService = utils.MapStringWithMode(data, "connection_service", resourceType, mode)
	state.ConnectionServiceRefType = utils.MapStringWithMode(data, "connection_service_ref_type_", resourceType, mode)
	state.Port = utils.MapStringWithMode(data, "port", resourceType, mode)
	state.SfpMacAddressOrSn = utils.MapStringWithMode(data, "sfp_mac_address_or_sn", resourceType, mode)

	return state
}

func (r *verityDeviceControllerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// =========================================================================
	// Skip if deleting
	// =========================================================================
	if req.Plan.Raw.IsNull() {
		return
	}

	var plan verityDeviceControllerResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// =========================================================================
	// Mode-aware field nullification
	// Set fields that don't apply to current mode to null to prevent
	// "known after apply" messages for irrelevant fields.
	// =========================================================================
	const resourceType = deviceControllerResourceType
	mode := r.provCtx.mode

	nullifier := &utils.ModeFieldNullifier{
		Ctx:          ctx,
		ResourceType: resourceType,
		Mode:         mode,
		Plan:         &resp.Plan,
	}

	nullifier.NullifyStrings(
		"ip_source", "controller_ip_and_mask", "gateway", "switch_ip_and_mask", "switch_gateway",
		"comm_type", "snmp_community_string", "uplink_port", "lldp_search_string", "ztp_identification",
		"located_by", "power_state", "communication_mode", "cli_access_mode",
		"username", "password", "enable_password", "ssh_key_or_password",
		"sdlc", "switchpoint", "switchpoint_ref_type_", "security_type",
		"snmpv3_username", "authentication_protocol", "passphrase", "private_protocol", "private_password",
		"password_encrypted", "enable_password_encrypted", "ssh_key_or_password_encrypted",
		"passphrase_encrypted", "private_password_encrypted",
		"device_managed_as", "switch", "switch_ref_type_",
		"connection_service", "connection_service_ref_type_", "port", "sfp_mac_address_or_sn",
	)

	nullifier.NullifyBools(
		"enable", "managed_on_native_vlan", "uses_tagged_packets",
	)
}
