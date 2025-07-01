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
	deviceControllerProps := openapi.ConfigPutRequestDeviceControllerDeviceControllerName{}
	deviceControllerProps.Name = openapi.PtrString(name)

	if !plan.Enable.IsNull() {
		deviceControllerProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
	}
	if !plan.IpSource.IsNull() {
		deviceControllerProps.IpSource = openapi.PtrString(plan.IpSource.ValueString())
	}
	if !plan.ControllerIpAndMask.IsNull() {
		deviceControllerProps.ControllerIpAndMask = openapi.PtrString(plan.ControllerIpAndMask.ValueString())
	}
	if !plan.Gateway.IsNull() {
		deviceControllerProps.Gateway = openapi.PtrString(plan.Gateway.ValueString())
	}
	if !plan.SwitchIpAndMask.IsNull() {
		deviceControllerProps.SwitchIpAndMask = openapi.PtrString(plan.SwitchIpAndMask.ValueString())
	}
	if !plan.SwitchGateway.IsNull() {
		deviceControllerProps.SwitchGateway = openapi.PtrString(plan.SwitchGateway.ValueString())
	}
	if !plan.CommType.IsNull() {
		deviceControllerProps.CommType = openapi.PtrString(plan.CommType.ValueString())
	}
	if !plan.SnmpCommunityString.IsNull() {
		deviceControllerProps.SnmpCommunityString = openapi.PtrString(plan.SnmpCommunityString.ValueString())
	}
	if !plan.UplinkPort.IsNull() {
		deviceControllerProps.UplinkPort = openapi.PtrString(plan.UplinkPort.ValueString())
	}
	if !plan.LldpSearchString.IsNull() {
		deviceControllerProps.LldpSearchString = openapi.PtrString(plan.LldpSearchString.ValueString())
	}
	if !plan.ZtpIdentification.IsNull() {
		deviceControllerProps.ZtpIdentification = openapi.PtrString(plan.ZtpIdentification.ValueString())
	}
	if !plan.LocatedBy.IsNull() {
		deviceControllerProps.LocatedBy = openapi.PtrString(plan.LocatedBy.ValueString())
	}
	if !plan.PowerState.IsNull() {
		deviceControllerProps.PowerState = openapi.PtrString(plan.PowerState.ValueString())
	}
	if !plan.CommunicationMode.IsNull() {
		deviceControllerProps.CommunicationMode = openapi.PtrString(plan.CommunicationMode.ValueString())
	}
	if !plan.CliAccessMode.IsNull() {
		deviceControllerProps.CliAccessMode = openapi.PtrString(plan.CliAccessMode.ValueString())
	}
	if !plan.Username.IsNull() {
		deviceControllerProps.Username = openapi.PtrString(plan.Username.ValueString())
	}
	if !plan.Password.IsNull() {
		deviceControllerProps.Password = openapi.PtrString(plan.Password.ValueString())
	}
	if !plan.EnablePassword.IsNull() {
		deviceControllerProps.EnablePassword = openapi.PtrString(plan.EnablePassword.ValueString())
	}
	if !plan.SshKeyOrPassword.IsNull() {
		deviceControllerProps.SshKeyOrPassword = openapi.PtrString(plan.SshKeyOrPassword.ValueString())
	}
	if !plan.ManagedOnNativeVlan.IsNull() {
		deviceControllerProps.ManagedOnNativeVlan = openapi.PtrBool(plan.ManagedOnNativeVlan.ValueBool())
	}
	if !plan.Sdlc.IsNull() {
		deviceControllerProps.Sdlc = openapi.PtrString(plan.Sdlc.ValueString())
	}
	if !plan.Switchpoint.IsNull() {
		deviceControllerProps.Switchpoint = openapi.PtrString(plan.Switchpoint.ValueString())
	}
	if !plan.SwitchpointRefType.IsNull() {
		deviceControllerProps.SwitchpointRefType = openapi.PtrString(plan.SwitchpointRefType.ValueString())
	}
	if !plan.SecurityType.IsNull() {
		deviceControllerProps.SecurityType = openapi.PtrString(plan.SecurityType.ValueString())
	}
	if !plan.Snmpv3Username.IsNull() {
		deviceControllerProps.Snmpv3Username = openapi.PtrString(plan.Snmpv3Username.ValueString())
	}
	if !plan.AuthenticationProtocol.IsNull() {
		deviceControllerProps.AuthenticationProtocol = openapi.PtrString(plan.AuthenticationProtocol.ValueString())
	}
	if !plan.Passphrase.IsNull() {
		deviceControllerProps.Passphrase = openapi.PtrString(plan.Passphrase.ValueString())
	}
	if !plan.PrivateProtocol.IsNull() {
		deviceControllerProps.PrivateProtocol = openapi.PtrString(plan.PrivateProtocol.ValueString())
	}
	if !plan.PrivatePassword.IsNull() {
		deviceControllerProps.PrivatePassword = openapi.PtrString(plan.PrivatePassword.ValueString())
	}
	if !plan.PasswordEncrypted.IsNull() {
		deviceControllerProps.PasswordEncrypted = openapi.PtrString(plan.PasswordEncrypted.ValueString())
	}
	if !plan.EnablePasswordEncrypted.IsNull() {
		deviceControllerProps.EnablePasswordEncrypted = openapi.PtrString(plan.EnablePasswordEncrypted.ValueString())
	}
	if !plan.SshKeyOrPasswordEncrypted.IsNull() {
		deviceControllerProps.SshKeyOrPasswordEncrypted = openapi.PtrString(plan.SshKeyOrPasswordEncrypted.ValueString())
	}
	if !plan.PassphraseEncrypted.IsNull() {
		deviceControllerProps.PassphraseEncrypted = openapi.PtrString(plan.PassphraseEncrypted.ValueString())
	}
	if !plan.PrivatePasswordEncrypted.IsNull() {
		deviceControllerProps.PrivatePasswordEncrypted = openapi.PtrString(plan.PrivatePasswordEncrypted.ValueString())
	}
	if !plan.DeviceManagedAs.IsNull() {
		deviceControllerProps.DeviceManagedAs = openapi.PtrString(plan.DeviceManagedAs.ValueString())
	}
	if !plan.Switch.IsNull() {
		deviceControllerProps.Switch = openapi.PtrString(plan.Switch.ValueString())
	}
	if !plan.SwitchRefType.IsNull() {
		deviceControllerProps.SwitchRefType = openapi.PtrString(plan.SwitchRefType.ValueString())
	}
	if !plan.ConnectionService.IsNull() {
		deviceControllerProps.ConnectionService = openapi.PtrString(plan.ConnectionService.ValueString())
	}
	if !plan.ConnectionServiceRefType.IsNull() {
		deviceControllerProps.ConnectionServiceRefType = openapi.PtrString(plan.ConnectionServiceRefType.ValueString())
	}
	if !plan.Port.IsNull() {
		deviceControllerProps.Port = openapi.PtrString(plan.Port.ValueString())
	}
	if !plan.SfpMacAddressOrSn.IsNull() {
		deviceControllerProps.SfpMacAddressOrSn = openapi.PtrString(plan.SfpMacAddressOrSn.ValueString())
	}
	if !plan.UsesTaggedPackets.IsNull() {
		deviceControllerProps.UsesTaggedPackets = openapi.PtrBool(plan.UsesTaggedPackets.ValueBool())
	}

	operationID := r.bulkOpsMgr.AddDeviceControllerPut(ctx, name, deviceControllerProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for device controller creation operation %s to complete", operationID))
	if err := r.bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Create Device Controller %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Device Controller %s creation operation completed successfully", name))
	clearCache(ctx, r.provCtx, "devicecontrollers")

	var minState verityDeviceControllerResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if deviceControllerData, exists := bulkMgr.GetResourceResponse("device_controller", name); exists {
			tflog.Debug(ctx, fmt.Sprintf("Using cached device controller data for %s after creation", name))
			state := populateDeviceControllerState(ctx, minState, deviceControllerData, &plan)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

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

	tflog.Debug(ctx, "Reading device controller resource")

	provCtx := r.provCtx
	bulkOpsMgr := provCtx.bulkOpsMgr
	deviceControllerName := state.Name.ValueString()

	var deviceControllerData map[string]interface{}
	var exists bool

	if bulkOpsMgr != nil {
		deviceControllerData, exists = bulkOpsMgr.GetResourceResponse("device_controller", deviceControllerName)
		if exists {
			tflog.Debug(ctx, fmt.Sprintf("Using cached device controller data for %s", deviceControllerName))
			state = populateDeviceControllerState(ctx, state, deviceControllerData, nil)
			resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
			return
		}
	}

	if bulkOpsMgr != nil && bulkOpsMgr.HasPendingOrRecentDeviceControllerOperations() {
		tflog.Info(ctx, fmt.Sprintf("Skipping device controller %s verification - trusting recent successful API operation", deviceControllerName))
		resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("No recent device controller operations found, performing normal verification for %s", deviceControllerName))

	type DeviceControllersResponse struct {
		DeviceController map[string]interface{} `json:"device_controller"`
	}

	var result DeviceControllersResponse
	var err error
	maxRetries := 3

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			sleepTime := time.Duration(100*(attempt)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Retrying device controller fetch in %v", sleepTime))
			time.Sleep(sleepTime)
		}

		deviceControllersData, fetchErr := getCachedResponse(ctx, r.provCtx, "devicecontrollers", func() (interface{}, error) {
			tflog.Debug(ctx, "Making API call to fetch device controllers")
			respAPI, err := r.client.DeviceControllersAPI.DevicecontrollersGet(ctx).Execute()
			if err != nil {
				return nil, fmt.Errorf("error reading device controllers: %v", err)
			}
			defer respAPI.Body.Close()

			var res DeviceControllersResponse
			if err := json.NewDecoder(respAPI.Body).Decode(&res); err != nil {
				return nil, fmt.Errorf("failed to decode device controllers response: %v", err)
			}

			tflog.Debug(ctx, fmt.Sprintf("Successfully fetched %d device controllers", len(res.DeviceController)))
			return res, nil
		})
		if fetchErr != nil {
			err = fetchErr
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch device controllers on attempt %d, retrying in %v", attempt+1, sleepTime))
			time.Sleep(sleepTime)
			continue
		}
		result = deviceControllersData.(DeviceControllersResponse)
		break
	}
	if err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read Device Controller %s", deviceControllerName))...,
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Looking for device controller with ID: %s", deviceControllerName))
	if data, ok := result.DeviceController[deviceControllerName].(map[string]interface{}); ok {
		deviceControllerData = data
		exists = true
		tflog.Debug(ctx, fmt.Sprintf("Found device controller directly by ID: %s", deviceControllerName))
	} else {
		for apiName, dc := range result.DeviceController {
			deviceController, ok := dc.(map[string]interface{})
			if !ok {
				continue
			}

			if name, ok := deviceController["name"].(string); ok && name == deviceControllerName {
				deviceControllerData = deviceController
				deviceControllerName = apiName
				exists = true
				tflog.Debug(ctx, fmt.Sprintf("Found device controller with name '%s' under API key '%s'", name, apiName))
				break
			}
		}
	}

	if !exists {
		tflog.Debug(ctx, fmt.Sprintf("Device Controller with ID '%s' not found in API response", deviceControllerName))
		resp.State.RemoveResource(ctx)
		return
	}

	state = populateDeviceControllerState(ctx, state, deviceControllerData, nil)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *verityDeviceControllerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state verityDeviceControllerResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
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
	deviceControllerProps := openapi.ConfigPutRequestDeviceControllerDeviceControllerName{}
	hasChanges := false

	if !plan.Enable.Equal(state.Enable) {
		deviceControllerProps.Enable = openapi.PtrBool(plan.Enable.ValueBool())
		hasChanges = true
	}
	if !plan.IpSource.Equal(state.IpSource) {
		deviceControllerProps.IpSource = openapi.PtrString(plan.IpSource.ValueString())
		hasChanges = true
	}
	if !plan.ControllerIpAndMask.Equal(state.ControllerIpAndMask) {
		deviceControllerProps.ControllerIpAndMask = openapi.PtrString(plan.ControllerIpAndMask.ValueString())
		hasChanges = true
	}
	if !plan.Gateway.Equal(state.Gateway) {
		deviceControllerProps.Gateway = openapi.PtrString(plan.Gateway.ValueString())
		hasChanges = true
	}
	if !plan.SwitchIpAndMask.Equal(state.SwitchIpAndMask) {
		deviceControllerProps.SwitchIpAndMask = openapi.PtrString(plan.SwitchIpAndMask.ValueString())
		hasChanges = true
	}
	if !plan.SwitchGateway.Equal(state.SwitchGateway) {
		deviceControllerProps.SwitchGateway = openapi.PtrString(plan.SwitchGateway.ValueString())
		hasChanges = true
	}
	if !plan.CommType.Equal(state.CommType) {
		deviceControllerProps.CommType = openapi.PtrString(plan.CommType.ValueString())
		hasChanges = true
	}
	if !plan.SnmpCommunityString.Equal(state.SnmpCommunityString) {
		deviceControllerProps.SnmpCommunityString = openapi.PtrString(plan.SnmpCommunityString.ValueString())
		hasChanges = true
	}
	if !plan.UplinkPort.Equal(state.UplinkPort) {
		deviceControllerProps.UplinkPort = openapi.PtrString(plan.UplinkPort.ValueString())
		hasChanges = true
	}
	if !plan.LldpSearchString.Equal(state.LldpSearchString) {
		deviceControllerProps.LldpSearchString = openapi.PtrString(plan.LldpSearchString.ValueString())
		hasChanges = true
	}
	if !plan.ZtpIdentification.Equal(state.ZtpIdentification) {
		deviceControllerProps.ZtpIdentification = openapi.PtrString(plan.ZtpIdentification.ValueString())
		hasChanges = true
	}
	if !plan.LocatedBy.Equal(state.LocatedBy) {
		deviceControllerProps.LocatedBy = openapi.PtrString(plan.LocatedBy.ValueString())
		hasChanges = true
	}
	if !plan.PowerState.Equal(state.PowerState) {
		deviceControllerProps.PowerState = openapi.PtrString(plan.PowerState.ValueString())
		hasChanges = true
	}
	if !plan.CommunicationMode.Equal(state.CommunicationMode) {
		deviceControllerProps.CommunicationMode = openapi.PtrString(plan.CommunicationMode.ValueString())
		hasChanges = true
	}
	if !plan.CliAccessMode.Equal(state.CliAccessMode) {
		deviceControllerProps.CliAccessMode = openapi.PtrString(plan.CliAccessMode.ValueString())
		hasChanges = true
	}
	if !plan.Username.Equal(state.Username) {
		deviceControllerProps.Username = openapi.PtrString(plan.Username.ValueString())
		hasChanges = true
	}
	if !plan.Password.Equal(state.Password) {
		deviceControllerProps.Password = openapi.PtrString(plan.Password.ValueString())
		hasChanges = true
	}
	if !plan.EnablePassword.Equal(state.EnablePassword) {
		deviceControllerProps.EnablePassword = openapi.PtrString(plan.EnablePassword.ValueString())
		hasChanges = true
	}
	if !plan.SshKeyOrPassword.Equal(state.SshKeyOrPassword) {
		deviceControllerProps.SshKeyOrPassword = openapi.PtrString(plan.SshKeyOrPassword.ValueString())
		hasChanges = true
	}
	if !plan.ManagedOnNativeVlan.Equal(state.ManagedOnNativeVlan) {
		deviceControllerProps.ManagedOnNativeVlan = openapi.PtrBool(plan.ManagedOnNativeVlan.ValueBool())
		hasChanges = true
	}
	if !plan.Sdlc.Equal(state.Sdlc) {
		deviceControllerProps.Sdlc = openapi.PtrString(plan.Sdlc.ValueString())
		hasChanges = true
	}
	if !plan.SecurityType.Equal(state.SecurityType) {
		deviceControllerProps.SecurityType = openapi.PtrString(plan.SecurityType.ValueString())
		hasChanges = true
	}
	if !plan.Snmpv3Username.Equal(state.Snmpv3Username) {
		deviceControllerProps.Snmpv3Username = openapi.PtrString(plan.Snmpv3Username.ValueString())
		hasChanges = true
	}
	if !plan.AuthenticationProtocol.Equal(state.AuthenticationProtocol) {
		deviceControllerProps.AuthenticationProtocol = openapi.PtrString(plan.AuthenticationProtocol.ValueString())
		hasChanges = true
	}
	if !plan.Passphrase.Equal(state.Passphrase) {
		deviceControllerProps.Passphrase = openapi.PtrString(plan.Passphrase.ValueString())
		hasChanges = true
	}
	if !plan.PrivateProtocol.Equal(state.PrivateProtocol) {
		deviceControllerProps.PrivateProtocol = openapi.PtrString(plan.PrivateProtocol.ValueString())
		hasChanges = true
	}
	if !plan.PrivatePassword.Equal(state.PrivatePassword) {
		deviceControllerProps.PrivatePassword = openapi.PtrString(plan.PrivatePassword.ValueString())
		hasChanges = true
	}
	if !plan.PasswordEncrypted.Equal(state.PasswordEncrypted) {
		deviceControllerProps.PasswordEncrypted = openapi.PtrString(plan.PasswordEncrypted.ValueString())
		hasChanges = true
	}
	if !plan.EnablePasswordEncrypted.Equal(state.EnablePasswordEncrypted) {
		deviceControllerProps.EnablePasswordEncrypted = openapi.PtrString(plan.EnablePasswordEncrypted.ValueString())
		hasChanges = true
	}
	if !plan.SshKeyOrPasswordEncrypted.Equal(state.SshKeyOrPasswordEncrypted) {
		deviceControllerProps.SshKeyOrPasswordEncrypted = openapi.PtrString(plan.SshKeyOrPasswordEncrypted.ValueString())
		hasChanges = true
	}
	if !plan.PassphraseEncrypted.Equal(state.PassphraseEncrypted) {
		deviceControllerProps.PassphraseEncrypted = openapi.PtrString(plan.PassphraseEncrypted.ValueString())
		hasChanges = true
	}
	if !plan.PrivatePasswordEncrypted.Equal(state.PrivatePasswordEncrypted) {
		deviceControllerProps.PrivatePasswordEncrypted = openapi.PtrString(plan.PrivatePasswordEncrypted.ValueString())
		hasChanges = true
	}
	if !plan.DeviceManagedAs.Equal(state.DeviceManagedAs) {
		deviceControllerProps.DeviceManagedAs = openapi.PtrString(plan.DeviceManagedAs.ValueString())
		hasChanges = true
	}
	if !plan.Port.Equal(state.Port) {
		deviceControllerProps.Port = openapi.PtrString(plan.Port.ValueString())
		hasChanges = true
	}
	if !plan.SfpMacAddressOrSn.Equal(state.SfpMacAddressOrSn) {
		deviceControllerProps.SfpMacAddressOrSn = openapi.PtrString(plan.SfpMacAddressOrSn.ValueString())
		hasChanges = true
	}
	if !plan.UsesTaggedPackets.Equal(state.UsesTaggedPackets) {
		deviceControllerProps.UsesTaggedPackets = openapi.PtrBool(plan.UsesTaggedPackets.ValueBool())
		hasChanges = true
	}

	// Handle switchpoint and switchpoint_ref_type_ according to "One ref type supported" rules
	switchpointChanged := !plan.Switchpoint.Equal(state.Switchpoint)
	switchpointRefTypeChanged := !plan.SwitchpointRefType.Equal(state.SwitchpointRefType)

	if switchpointChanged || switchpointRefTypeChanged {
		// Validate using "one ref type supported" rules
		if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
			plan.Switchpoint, plan.SwitchpointRefType,
			"switchpoint", "switchpoint_ref_type_",
			switchpointChanged, switchpointRefTypeChanged) {
			return
		}

		// Only send the base field if only it changed
		if switchpointChanged && !switchpointRefTypeChanged {
			// Just send the base field
			if !plan.Switchpoint.IsNull() && plan.Switchpoint.ValueString() != "" {
				deviceControllerProps.Switchpoint = openapi.PtrString(plan.Switchpoint.ValueString())
			} else {
				deviceControllerProps.Switchpoint = openapi.PtrString("")
			}
			hasChanges = true
		} else if switchpointRefTypeChanged {
			// Send both fields
			if !plan.Switchpoint.IsNull() && plan.Switchpoint.ValueString() != "" {
				deviceControllerProps.Switchpoint = openapi.PtrString(plan.Switchpoint.ValueString())
			} else {
				deviceControllerProps.Switchpoint = openapi.PtrString("")
			}

			if !plan.SwitchpointRefType.IsNull() && plan.SwitchpointRefType.ValueString() != "" {
				deviceControllerProps.SwitchpointRefType = openapi.PtrString(plan.SwitchpointRefType.ValueString())
			} else {
				deviceControllerProps.SwitchpointRefType = openapi.PtrString("")
			}
			hasChanges = true
		}
	}

	// Handle switch and switch_ref_type_ according to "One ref type supported" rules
	switchChanged := !plan.Switch.Equal(state.Switch)
	switchRefTypeChanged := !plan.SwitchRefType.Equal(state.SwitchRefType)

	if switchChanged || switchRefTypeChanged {
		// Validate using "one ref type supported" rules
		if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
			plan.Switch, plan.SwitchRefType,
			"switch", "switch_ref_type_",
			switchChanged, switchRefTypeChanged) {
			return
		}

		// Only send the base field if only it changed
		if switchChanged && !switchRefTypeChanged {
			// Just send the base field
			if !plan.Switch.IsNull() && plan.Switch.ValueString() != "" {
				deviceControllerProps.Switch = openapi.PtrString(plan.Switch.ValueString())
			} else {
				deviceControllerProps.Switch = openapi.PtrString("")
			}
			hasChanges = true
		} else if switchRefTypeChanged {
			// Send both fields
			if !plan.Switch.IsNull() && plan.Switch.ValueString() != "" {
				deviceControllerProps.Switch = openapi.PtrString(plan.Switch.ValueString())
			} else {
				deviceControllerProps.Switch = openapi.PtrString("")
			}

			if !plan.SwitchRefType.IsNull() && plan.SwitchRefType.ValueString() != "" {
				deviceControllerProps.SwitchRefType = openapi.PtrString(plan.SwitchRefType.ValueString())
			} else {
				deviceControllerProps.SwitchRefType = openapi.PtrString("")
			}
			hasChanges = true
		}
	}

	// Handle connection_service and connection_service_ref_type_ according to "One ref type supported" rules
	connectionServiceChanged := !plan.ConnectionService.Equal(state.ConnectionService)
	connectionServiceRefTypeChanged := !plan.ConnectionServiceRefType.Equal(state.ConnectionServiceRefType)

	if connectionServiceChanged || connectionServiceRefTypeChanged {
		// Validate using "one ref type supported" rules
		if !utils.ValidateOneRefTypeSupported(&resp.Diagnostics,
			plan.ConnectionService, plan.ConnectionServiceRefType,
			"connection_service", "connection_service_ref_type_",
			connectionServiceChanged, connectionServiceRefTypeChanged) {
			return
		}

		// Only send the base field if only it changed
		if connectionServiceChanged && !connectionServiceRefTypeChanged {
			// Just send the base field
			if !plan.ConnectionService.IsNull() && plan.ConnectionService.ValueString() != "" {
				deviceControllerProps.ConnectionService = openapi.PtrString(plan.ConnectionService.ValueString())
			} else {
				deviceControllerProps.ConnectionService = openapi.PtrString("")
			}
			hasChanges = true
		} else if connectionServiceRefTypeChanged {
			// Send both fields
			if !plan.ConnectionService.IsNull() && plan.ConnectionService.ValueString() != "" {
				deviceControllerProps.ConnectionService = openapi.PtrString(plan.ConnectionService.ValueString())
			} else {
				deviceControllerProps.ConnectionService = openapi.PtrString("")
			}

			if !plan.ConnectionServiceRefType.IsNull() && plan.ConnectionServiceRefType.ValueString() != "" {
				deviceControllerProps.ConnectionServiceRefType = openapi.PtrString(plan.ConnectionServiceRefType.ValueString())
			} else {
				deviceControllerProps.ConnectionServiceRefType = openapi.PtrString("")
			}
			hasChanges = true
		}
	}

	if !hasChanges {
		resp.Diagnostics.Append(resp.State.Set(ctx, plan)...)
		return
	}

	bulkOpsMgr := r.bulkOpsMgr
	operationID := bulkOpsMgr.AddDeviceControllerPatch(ctx, name, deviceControllerProps)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for device controller update operation %s to complete", operationID))
	if err := bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Update Device Controller %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Device Controller %s update operation completed successfully", name))
	clearCache(ctx, r.provCtx, "devicecontrollers")

	var minState verityDeviceControllerResourceModel
	minState.Name = types.StringValue(name)
	resp.Diagnostics.Append(resp.State.Set(ctx, &minState)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if bulkMgr := r.provCtx.bulkOpsMgr; bulkMgr != nil {
		if deviceControllerData, exists := bulkMgr.GetResourceResponse("device_controller", name); exists {
			tflog.Debug(ctx, fmt.Sprintf("Using cached device controller data for %s after update", name))
			state = populateDeviceControllerState(ctx, state, deviceControllerData, &plan)
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
	bulkOpsMgr := r.bulkOpsMgr
	operationID := bulkOpsMgr.AddDeviceControllerDelete(ctx, name)
	r.notifyOperationAdded()

	tflog.Debug(ctx, fmt.Sprintf("Waiting for device controller deletion operation %s to complete", operationID))
	if err := bulkOpsMgr.WaitForOperation(ctx, operationID, utils.OperationTimeout); err != nil {
		resp.Diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Delete Device Controller %s", name))...,
		)
		return
	}

	tflog.Info(ctx, fmt.Sprintf("Device Controller %s deletion operation completed successfully", name))
	clearCache(ctx, r.provCtx, "devicecontrollers")
	resp.State.RemoveResource(ctx)
}

func (r *verityDeviceControllerResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("name"), req, resp)
}

func populateDeviceControllerState(ctx context.Context, state verityDeviceControllerResourceModel, deviceControllerData map[string]interface{}, plan *verityDeviceControllerResourceModel) verityDeviceControllerResourceModel {
	state.Name = types.StringValue(fmt.Sprintf("%v", deviceControllerData["name"]))

	if val, ok := deviceControllerData["enable"].(bool); ok {
		state.Enable = types.BoolValue(val)
	} else if plan != nil && !plan.Enable.IsNull() {
		state.Enable = plan.Enable
	} else {
		state.Enable = types.BoolNull()
	}

	stringAttrs := map[string]*types.String{
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

	for apiKey, stateField := range stringAttrs {
		if val, ok := deviceControllerData[apiKey].(string); ok {
			*stateField = types.StringValue(val)
		} else if plan != nil {
			planField := getDeviceControllerPlanStringField(plan, apiKey)
			if planField != nil && !planField.IsNull() {
				*stateField = *planField
			} else {
				*stateField = types.StringNull()
			}
		} else {
			*stateField = types.StringNull()
		}
	}

	if val, ok := deviceControllerData["managed_on_native_vlan"].(bool); ok {
		state.ManagedOnNativeVlan = types.BoolValue(val)
	} else if plan != nil && !plan.ManagedOnNativeVlan.IsNull() {
		state.ManagedOnNativeVlan = plan.ManagedOnNativeVlan
	} else {
		state.ManagedOnNativeVlan = types.BoolNull()
	}

	if val, ok := deviceControllerData["uses_tagged_packets"].(bool); ok {
		state.UsesTaggedPackets = types.BoolValue(val)
	} else if plan != nil && !plan.UsesTaggedPackets.IsNull() {
		state.UsesTaggedPackets = plan.UsesTaggedPackets
	} else {
		state.UsesTaggedPackets = types.BoolNull()
	}

	return state
}

func getDeviceControllerPlanStringField(plan *verityDeviceControllerResourceModel, apiKey string) *types.String {
	switch apiKey {
	case "ip_source":
		return &plan.IpSource
	case "controller_ip_and_mask":
		return &plan.ControllerIpAndMask
	case "gateway":
		return &plan.Gateway
	case "switch_ip_and_mask":
		return &plan.SwitchIpAndMask
	case "switch_gateway":
		return &plan.SwitchGateway
	case "comm_type":
		return &plan.CommType
	case "snmp_community_string":
		return &plan.SnmpCommunityString
	case "uplink_port":
		return &plan.UplinkPort
	case "lldp_search_string":
		return &plan.LldpSearchString
	case "ztp_identification":
		return &plan.ZtpIdentification
	case "located_by":
		return &plan.LocatedBy
	case "power_state":
		return &plan.PowerState
	case "communication_mode":
		return &plan.CommunicationMode
	case "cli_access_mode":
		return &plan.CliAccessMode
	case "username":
		return &plan.Username
	case "password":
		return &plan.Password
	case "enable_password":
		return &plan.EnablePassword
	case "ssh_key_or_password":
		return &plan.SshKeyOrPassword
	case "sdlc":
		return &plan.Sdlc
	case "switchpoint":
		return &plan.Switchpoint
	case "switchpoint_ref_type_":
		return &plan.SwitchpointRefType
	case "security_type":
		return &plan.SecurityType
	case "snmpv3_username":
		return &plan.Snmpv3Username
	case "authentication_protocol":
		return &plan.AuthenticationProtocol
	case "passphrase":
		return &plan.Passphrase
	case "private_protocol":
		return &plan.PrivateProtocol
	case "private_password":
		return &plan.PrivatePassword
	case "password_encrypted":
		return &plan.PasswordEncrypted
	case "enable_password_encrypted":
		return &plan.EnablePasswordEncrypted
	case "ssh_key_or_password_encrypted":
		return &plan.SshKeyOrPasswordEncrypted
	case "passphrase_encrypted":
		return &plan.PassphraseEncrypted
	case "private_password_encrypted":
		return &plan.PrivatePasswordEncrypted
	case "device_managed_as":
		return &plan.DeviceManagedAs
	case "switch":
		return &plan.Switch
	case "switch_ref_type_":
		return &plan.SwitchRefType
	case "connection_service":
		return &plan.ConnectionService
	case "connection_service_ref_type_":
		return &plan.ConnectionServiceRefType
	case "port":
		return &plan.Port
	case "sfp_mac_address_or_sn":
		return &plan.SfpMacAddressOrSn
	default:
		return nil
	}
}
