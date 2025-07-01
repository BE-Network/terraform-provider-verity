# DevicecontrollersPutRequestDeviceControllerValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**IpSource** | Pointer to **string** | IP Source | [optional] [default to "dhcp"]
**ControllerIpAndMask** | Pointer to **string** | Controller IP and Mask | [optional] [default to ""]
**Gateway** | Pointer to **string** | Gateway | [optional] [default to ""]
**SwitchIpAndMask** | Pointer to **string** | Switch IP and Mask | [optional] [default to ""]
**SwitchGateway** | Pointer to **string** | Gateway of Managed Device | [optional] [default to ""]
**CommType** | Pointer to **string** | Comm Type | [optional] [default to "snmpv2"]
**SnmpCommunityString** | Pointer to **string** | Comm Credentials | [optional] [default to ""]
**UplinkPort** | Pointer to **string** | Uplink Port of Managed Device | [optional] [default to ""]
**LldpSearchString** | Pointer to **string** | Optional unless Located By is \&quot;LLDP\&quot; or Device managed as \&quot;Active SFP\&quot;. Must be either the chassis-id or the hostname of the LLDP from the managed device. Used to detect connections between managed devices. If blank, the chassis-id detected by the Device Controller via SNMP/CLI is used | [optional] [default to ""]
**ZtpIdentification** | Pointer to **string** | Service Tag or Serial Number to identify device for Zero Touch Provisioning | [optional] [default to ""]
**LocatedBy** | Pointer to **string** | Controls how the system locates this Device within its LAN | [optional] [default to "LLDP"]
**PowerState** | Pointer to **string** | Power state of Switch Controller | [optional] [default to "on"]
**CommunicationMode** | Pointer to **string** | Communication Mode | [optional] [default to "generic_snmp"]
**CliAccessMode** | Pointer to **string** | CLI Access Mode | [optional] [default to "SSH"]
**Username** | Pointer to **string** | Username | [optional] [default to ""]
**Password** | Pointer to **string** | Password | [optional] [default to ""]
**EnablePassword** | Pointer to **string** | Enable Password - to enable privileged CLI operations | [optional] [default to ""]
**SshKeyOrPassword** | Pointer to **string** | SSH Key or Password | [optional] [default to ""]
**ManagedOnNativeVlan** | Pointer to **bool** | Managed on native VLAN | [optional] [default to false]
**Sdlc** | Pointer to **string** | SDLC that Device Controller belongs to | [optional] [default to ""]
**Switchpoint** | Pointer to **string** | Endpoint reference | [optional] [default to ""]
**SwitchpointRefType** | Pointer to **string** | Object type for switchpoint field | [optional] 
**SecurityType** | Pointer to **string** | Security level | [optional] [default to "noAuthNoPriv"]
**Snmpv3Username** | Pointer to **string** | Username | [optional] [default to ""]
**AuthenticationProtocol** | Pointer to **string** | Protocol | [optional] [default to "MD5"]
**Passphrase** | Pointer to **string** | Passphrase | [optional] [default to ""]
**PrivateProtocol** | Pointer to **string** | Protocol | [optional] [default to "DES"]
**PrivatePassword** | Pointer to **string** | Password | [optional] [default to ""]
**PasswordEncrypted** | Pointer to **string** | Password | [optional] [default to ""]
**EnablePasswordEncrypted** | Pointer to **string** | Enable Password - to enable privileged CLI operations | [optional] [default to ""]
**SshKeyOrPasswordEncrypted** | Pointer to **string** | SSH Key or Password | [optional] [default to ""]
**PassphraseEncrypted** | Pointer to **string** | Passphrase | [optional] [default to ""]
**PrivatePasswordEncrypted** | Pointer to **string** | Password | [optional] [default to ""]
**DeviceManagedAs** | Pointer to **string** | Device managed as | [optional] [default to "switch"]
**Switch** | Pointer to **string** | Endpoint locating the Switch to be controlled | [optional] [default to ""]
**SwitchRefType** | Pointer to **string** | Object type for switch field | [optional] 
**ConnectionService** | Pointer to **string** | Connect a Service | [optional] [default to ""]
**ConnectionServiceRefType** | Pointer to **string** | Object type for connection_service field | [optional] 
**Port** | Pointer to **string** | Port locating the Switch to be controlled | [optional] [default to ""]
**SfpMacAddressOrSn** | Pointer to **string** | SFP MAC Address or SN | [optional] [default to ""]
**UsesTaggedPackets** | Pointer to **bool** | Indicates if the direct interface expects tagged or untagged packets | [optional] [default to true]

## Methods

### NewDevicecontrollersPutRequestDeviceControllerValue

`func NewDevicecontrollersPutRequestDeviceControllerValue() *DevicecontrollersPutRequestDeviceControllerValue`

NewDevicecontrollersPutRequestDeviceControllerValue instantiates a new DevicecontrollersPutRequestDeviceControllerValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDevicecontrollersPutRequestDeviceControllerValueWithDefaults

`func NewDevicecontrollersPutRequestDeviceControllerValueWithDefaults() *DevicecontrollersPutRequestDeviceControllerValue`

NewDevicecontrollersPutRequestDeviceControllerValueWithDefaults instantiates a new DevicecontrollersPutRequestDeviceControllerValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetIpSource

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetIpSource() string`

GetIpSource returns the IpSource field if non-nil, zero value otherwise.

### GetIpSourceOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetIpSourceOk() (*string, bool)`

GetIpSourceOk returns a tuple with the IpSource field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpSource

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetIpSource(v string)`

SetIpSource sets IpSource field to given value.

### HasIpSource

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasIpSource() bool`

HasIpSource returns a boolean if a field has been set.

### GetControllerIpAndMask

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetControllerIpAndMask() string`

GetControllerIpAndMask returns the ControllerIpAndMask field if non-nil, zero value otherwise.

### GetControllerIpAndMaskOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetControllerIpAndMaskOk() (*string, bool)`

GetControllerIpAndMaskOk returns a tuple with the ControllerIpAndMask field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetControllerIpAndMask

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetControllerIpAndMask(v string)`

SetControllerIpAndMask sets ControllerIpAndMask field to given value.

### HasControllerIpAndMask

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasControllerIpAndMask() bool`

HasControllerIpAndMask returns a boolean if a field has been set.

### GetGateway

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetGateway() string`

GetGateway returns the Gateway field if non-nil, zero value otherwise.

### GetGatewayOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetGatewayOk() (*string, bool)`

GetGatewayOk returns a tuple with the Gateway field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGateway

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetGateway(v string)`

SetGateway sets Gateway field to given value.

### HasGateway

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasGateway() bool`

HasGateway returns a boolean if a field has been set.

### GetSwitchIpAndMask

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetSwitchIpAndMask() string`

GetSwitchIpAndMask returns the SwitchIpAndMask field if non-nil, zero value otherwise.

### GetSwitchIpAndMaskOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetSwitchIpAndMaskOk() (*string, bool)`

GetSwitchIpAndMaskOk returns a tuple with the SwitchIpAndMask field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchIpAndMask

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetSwitchIpAndMask(v string)`

SetSwitchIpAndMask sets SwitchIpAndMask field to given value.

### HasSwitchIpAndMask

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasSwitchIpAndMask() bool`

HasSwitchIpAndMask returns a boolean if a field has been set.

### GetSwitchGateway

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetSwitchGateway() string`

GetSwitchGateway returns the SwitchGateway field if non-nil, zero value otherwise.

### GetSwitchGatewayOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetSwitchGatewayOk() (*string, bool)`

GetSwitchGatewayOk returns a tuple with the SwitchGateway field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchGateway

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetSwitchGateway(v string)`

SetSwitchGateway sets SwitchGateway field to given value.

### HasSwitchGateway

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasSwitchGateway() bool`

HasSwitchGateway returns a boolean if a field has been set.

### GetCommType

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetCommType() string`

GetCommType returns the CommType field if non-nil, zero value otherwise.

### GetCommTypeOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetCommTypeOk() (*string, bool)`

GetCommTypeOk returns a tuple with the CommType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommType

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetCommType(v string)`

SetCommType sets CommType field to given value.

### HasCommType

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasCommType() bool`

HasCommType returns a boolean if a field has been set.

### GetSnmpCommunityString

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetSnmpCommunityString() string`

GetSnmpCommunityString returns the SnmpCommunityString field if non-nil, zero value otherwise.

### GetSnmpCommunityStringOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetSnmpCommunityStringOk() (*string, bool)`

GetSnmpCommunityStringOk returns a tuple with the SnmpCommunityString field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSnmpCommunityString

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetSnmpCommunityString(v string)`

SetSnmpCommunityString sets SnmpCommunityString field to given value.

### HasSnmpCommunityString

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasSnmpCommunityString() bool`

HasSnmpCommunityString returns a boolean if a field has been set.

### GetUplinkPort

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetUplinkPort() string`

GetUplinkPort returns the UplinkPort field if non-nil, zero value otherwise.

### GetUplinkPortOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetUplinkPortOk() (*string, bool)`

GetUplinkPortOk returns a tuple with the UplinkPort field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUplinkPort

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetUplinkPort(v string)`

SetUplinkPort sets UplinkPort field to given value.

### HasUplinkPort

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasUplinkPort() bool`

HasUplinkPort returns a boolean if a field has been set.

### GetLldpSearchString

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetLldpSearchString() string`

GetLldpSearchString returns the LldpSearchString field if non-nil, zero value otherwise.

### GetLldpSearchStringOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetLldpSearchStringOk() (*string, bool)`

GetLldpSearchStringOk returns a tuple with the LldpSearchString field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLldpSearchString

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetLldpSearchString(v string)`

SetLldpSearchString sets LldpSearchString field to given value.

### HasLldpSearchString

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasLldpSearchString() bool`

HasLldpSearchString returns a boolean if a field has been set.

### GetZtpIdentification

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetZtpIdentification() string`

GetZtpIdentification returns the ZtpIdentification field if non-nil, zero value otherwise.

### GetZtpIdentificationOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetZtpIdentificationOk() (*string, bool)`

GetZtpIdentificationOk returns a tuple with the ZtpIdentification field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetZtpIdentification

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetZtpIdentification(v string)`

SetZtpIdentification sets ZtpIdentification field to given value.

### HasZtpIdentification

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasZtpIdentification() bool`

HasZtpIdentification returns a boolean if a field has been set.

### GetLocatedBy

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetLocatedBy() string`

GetLocatedBy returns the LocatedBy field if non-nil, zero value otherwise.

### GetLocatedByOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetLocatedByOk() (*string, bool)`

GetLocatedByOk returns a tuple with the LocatedBy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocatedBy

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetLocatedBy(v string)`

SetLocatedBy sets LocatedBy field to given value.

### HasLocatedBy

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasLocatedBy() bool`

HasLocatedBy returns a boolean if a field has been set.

### GetPowerState

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetPowerState() string`

GetPowerState returns the PowerState field if non-nil, zero value otherwise.

### GetPowerStateOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetPowerStateOk() (*string, bool)`

GetPowerStateOk returns a tuple with the PowerState field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPowerState

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetPowerState(v string)`

SetPowerState sets PowerState field to given value.

### HasPowerState

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasPowerState() bool`

HasPowerState returns a boolean if a field has been set.

### GetCommunicationMode

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetCommunicationMode() string`

GetCommunicationMode returns the CommunicationMode field if non-nil, zero value otherwise.

### GetCommunicationModeOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetCommunicationModeOk() (*string, bool)`

GetCommunicationModeOk returns a tuple with the CommunicationMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommunicationMode

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetCommunicationMode(v string)`

SetCommunicationMode sets CommunicationMode field to given value.

### HasCommunicationMode

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasCommunicationMode() bool`

HasCommunicationMode returns a boolean if a field has been set.

### GetCliAccessMode

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetCliAccessMode() string`

GetCliAccessMode returns the CliAccessMode field if non-nil, zero value otherwise.

### GetCliAccessModeOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetCliAccessModeOk() (*string, bool)`

GetCliAccessModeOk returns a tuple with the CliAccessMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCliAccessMode

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetCliAccessMode(v string)`

SetCliAccessMode sets CliAccessMode field to given value.

### HasCliAccessMode

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasCliAccessMode() bool`

HasCliAccessMode returns a boolean if a field has been set.

### GetUsername

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsername

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetUsername(v string)`

SetUsername sets Username field to given value.

### HasUsername

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasUsername() bool`

HasUsername returns a boolean if a field has been set.

### GetPassword

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetPassword(v string)`

SetPassword sets Password field to given value.

### HasPassword

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasPassword() bool`

HasPassword returns a boolean if a field has been set.

### GetEnablePassword

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetEnablePassword() string`

GetEnablePassword returns the EnablePassword field if non-nil, zero value otherwise.

### GetEnablePasswordOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetEnablePasswordOk() (*string, bool)`

GetEnablePasswordOk returns a tuple with the EnablePassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnablePassword

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetEnablePassword(v string)`

SetEnablePassword sets EnablePassword field to given value.

### HasEnablePassword

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasEnablePassword() bool`

HasEnablePassword returns a boolean if a field has been set.

### GetSshKeyOrPassword

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetSshKeyOrPassword() string`

GetSshKeyOrPassword returns the SshKeyOrPassword field if non-nil, zero value otherwise.

### GetSshKeyOrPasswordOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetSshKeyOrPasswordOk() (*string, bool)`

GetSshKeyOrPasswordOk returns a tuple with the SshKeyOrPassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSshKeyOrPassword

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetSshKeyOrPassword(v string)`

SetSshKeyOrPassword sets SshKeyOrPassword field to given value.

### HasSshKeyOrPassword

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasSshKeyOrPassword() bool`

HasSshKeyOrPassword returns a boolean if a field has been set.

### GetManagedOnNativeVlan

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetManagedOnNativeVlan() bool`

GetManagedOnNativeVlan returns the ManagedOnNativeVlan field if non-nil, zero value otherwise.

### GetManagedOnNativeVlanOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetManagedOnNativeVlanOk() (*bool, bool)`

GetManagedOnNativeVlanOk returns a tuple with the ManagedOnNativeVlan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetManagedOnNativeVlan

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetManagedOnNativeVlan(v bool)`

SetManagedOnNativeVlan sets ManagedOnNativeVlan field to given value.

### HasManagedOnNativeVlan

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasManagedOnNativeVlan() bool`

HasManagedOnNativeVlan returns a boolean if a field has been set.

### GetSdlc

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetSdlc() string`

GetSdlc returns the Sdlc field if non-nil, zero value otherwise.

### GetSdlcOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetSdlcOk() (*string, bool)`

GetSdlcOk returns a tuple with the Sdlc field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSdlc

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetSdlc(v string)`

SetSdlc sets Sdlc field to given value.

### HasSdlc

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasSdlc() bool`

HasSdlc returns a boolean if a field has been set.

### GetSwitchpoint

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetSwitchpoint() string`

GetSwitchpoint returns the Switchpoint field if non-nil, zero value otherwise.

### GetSwitchpointOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetSwitchpointOk() (*string, bool)`

GetSwitchpointOk returns a tuple with the Switchpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchpoint

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetSwitchpoint(v string)`

SetSwitchpoint sets Switchpoint field to given value.

### HasSwitchpoint

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasSwitchpoint() bool`

HasSwitchpoint returns a boolean if a field has been set.

### GetSwitchpointRefType

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetSwitchpointRefType() string`

GetSwitchpointRefType returns the SwitchpointRefType field if non-nil, zero value otherwise.

### GetSwitchpointRefTypeOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetSwitchpointRefTypeOk() (*string, bool)`

GetSwitchpointRefTypeOk returns a tuple with the SwitchpointRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchpointRefType

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetSwitchpointRefType(v string)`

SetSwitchpointRefType sets SwitchpointRefType field to given value.

### HasSwitchpointRefType

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasSwitchpointRefType() bool`

HasSwitchpointRefType returns a boolean if a field has been set.

### GetSecurityType

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetSecurityType() string`

GetSecurityType returns the SecurityType field if non-nil, zero value otherwise.

### GetSecurityTypeOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetSecurityTypeOk() (*string, bool)`

GetSecurityTypeOk returns a tuple with the SecurityType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSecurityType

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetSecurityType(v string)`

SetSecurityType sets SecurityType field to given value.

### HasSecurityType

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasSecurityType() bool`

HasSecurityType returns a boolean if a field has been set.

### GetSnmpv3Username

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetSnmpv3Username() string`

GetSnmpv3Username returns the Snmpv3Username field if non-nil, zero value otherwise.

### GetSnmpv3UsernameOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetSnmpv3UsernameOk() (*string, bool)`

GetSnmpv3UsernameOk returns a tuple with the Snmpv3Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSnmpv3Username

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetSnmpv3Username(v string)`

SetSnmpv3Username sets Snmpv3Username field to given value.

### HasSnmpv3Username

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasSnmpv3Username() bool`

HasSnmpv3Username returns a boolean if a field has been set.

### GetAuthenticationProtocol

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetAuthenticationProtocol() string`

GetAuthenticationProtocol returns the AuthenticationProtocol field if non-nil, zero value otherwise.

### GetAuthenticationProtocolOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetAuthenticationProtocolOk() (*string, bool)`

GetAuthenticationProtocolOk returns a tuple with the AuthenticationProtocol field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthenticationProtocol

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetAuthenticationProtocol(v string)`

SetAuthenticationProtocol sets AuthenticationProtocol field to given value.

### HasAuthenticationProtocol

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasAuthenticationProtocol() bool`

HasAuthenticationProtocol returns a boolean if a field has been set.

### GetPassphrase

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetPassphrase() string`

GetPassphrase returns the Passphrase field if non-nil, zero value otherwise.

### GetPassphraseOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetPassphraseOk() (*string, bool)`

GetPassphraseOk returns a tuple with the Passphrase field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassphrase

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetPassphrase(v string)`

SetPassphrase sets Passphrase field to given value.

### HasPassphrase

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasPassphrase() bool`

HasPassphrase returns a boolean if a field has been set.

### GetPrivateProtocol

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetPrivateProtocol() string`

GetPrivateProtocol returns the PrivateProtocol field if non-nil, zero value otherwise.

### GetPrivateProtocolOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetPrivateProtocolOk() (*string, bool)`

GetPrivateProtocolOk returns a tuple with the PrivateProtocol field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrivateProtocol

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetPrivateProtocol(v string)`

SetPrivateProtocol sets PrivateProtocol field to given value.

### HasPrivateProtocol

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasPrivateProtocol() bool`

HasPrivateProtocol returns a boolean if a field has been set.

### GetPrivatePassword

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetPrivatePassword() string`

GetPrivatePassword returns the PrivatePassword field if non-nil, zero value otherwise.

### GetPrivatePasswordOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetPrivatePasswordOk() (*string, bool)`

GetPrivatePasswordOk returns a tuple with the PrivatePassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrivatePassword

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetPrivatePassword(v string)`

SetPrivatePassword sets PrivatePassword field to given value.

### HasPrivatePassword

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasPrivatePassword() bool`

HasPrivatePassword returns a boolean if a field has been set.

### GetPasswordEncrypted

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetPasswordEncrypted() string`

GetPasswordEncrypted returns the PasswordEncrypted field if non-nil, zero value otherwise.

### GetPasswordEncryptedOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetPasswordEncryptedOk() (*string, bool)`

GetPasswordEncryptedOk returns a tuple with the PasswordEncrypted field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPasswordEncrypted

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetPasswordEncrypted(v string)`

SetPasswordEncrypted sets PasswordEncrypted field to given value.

### HasPasswordEncrypted

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasPasswordEncrypted() bool`

HasPasswordEncrypted returns a boolean if a field has been set.

### GetEnablePasswordEncrypted

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetEnablePasswordEncrypted() string`

GetEnablePasswordEncrypted returns the EnablePasswordEncrypted field if non-nil, zero value otherwise.

### GetEnablePasswordEncryptedOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetEnablePasswordEncryptedOk() (*string, bool)`

GetEnablePasswordEncryptedOk returns a tuple with the EnablePasswordEncrypted field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnablePasswordEncrypted

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetEnablePasswordEncrypted(v string)`

SetEnablePasswordEncrypted sets EnablePasswordEncrypted field to given value.

### HasEnablePasswordEncrypted

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasEnablePasswordEncrypted() bool`

HasEnablePasswordEncrypted returns a boolean if a field has been set.

### GetSshKeyOrPasswordEncrypted

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetSshKeyOrPasswordEncrypted() string`

GetSshKeyOrPasswordEncrypted returns the SshKeyOrPasswordEncrypted field if non-nil, zero value otherwise.

### GetSshKeyOrPasswordEncryptedOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetSshKeyOrPasswordEncryptedOk() (*string, bool)`

GetSshKeyOrPasswordEncryptedOk returns a tuple with the SshKeyOrPasswordEncrypted field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSshKeyOrPasswordEncrypted

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetSshKeyOrPasswordEncrypted(v string)`

SetSshKeyOrPasswordEncrypted sets SshKeyOrPasswordEncrypted field to given value.

### HasSshKeyOrPasswordEncrypted

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasSshKeyOrPasswordEncrypted() bool`

HasSshKeyOrPasswordEncrypted returns a boolean if a field has been set.

### GetPassphraseEncrypted

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetPassphraseEncrypted() string`

GetPassphraseEncrypted returns the PassphraseEncrypted field if non-nil, zero value otherwise.

### GetPassphraseEncryptedOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetPassphraseEncryptedOk() (*string, bool)`

GetPassphraseEncryptedOk returns a tuple with the PassphraseEncrypted field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassphraseEncrypted

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetPassphraseEncrypted(v string)`

SetPassphraseEncrypted sets PassphraseEncrypted field to given value.

### HasPassphraseEncrypted

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasPassphraseEncrypted() bool`

HasPassphraseEncrypted returns a boolean if a field has been set.

### GetPrivatePasswordEncrypted

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetPrivatePasswordEncrypted() string`

GetPrivatePasswordEncrypted returns the PrivatePasswordEncrypted field if non-nil, zero value otherwise.

### GetPrivatePasswordEncryptedOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetPrivatePasswordEncryptedOk() (*string, bool)`

GetPrivatePasswordEncryptedOk returns a tuple with the PrivatePasswordEncrypted field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrivatePasswordEncrypted

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetPrivatePasswordEncrypted(v string)`

SetPrivatePasswordEncrypted sets PrivatePasswordEncrypted field to given value.

### HasPrivatePasswordEncrypted

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasPrivatePasswordEncrypted() bool`

HasPrivatePasswordEncrypted returns a boolean if a field has been set.

### GetDeviceManagedAs

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetDeviceManagedAs() string`

GetDeviceManagedAs returns the DeviceManagedAs field if non-nil, zero value otherwise.

### GetDeviceManagedAsOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetDeviceManagedAsOk() (*string, bool)`

GetDeviceManagedAsOk returns a tuple with the DeviceManagedAs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceManagedAs

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetDeviceManagedAs(v string)`

SetDeviceManagedAs sets DeviceManagedAs field to given value.

### HasDeviceManagedAs

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasDeviceManagedAs() bool`

HasDeviceManagedAs returns a boolean if a field has been set.

### GetSwitch

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetSwitch() string`

GetSwitch returns the Switch field if non-nil, zero value otherwise.

### GetSwitchOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetSwitchOk() (*string, bool)`

GetSwitchOk returns a tuple with the Switch field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitch

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetSwitch(v string)`

SetSwitch sets Switch field to given value.

### HasSwitch

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasSwitch() bool`

HasSwitch returns a boolean if a field has been set.

### GetSwitchRefType

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetSwitchRefType() string`

GetSwitchRefType returns the SwitchRefType field if non-nil, zero value otherwise.

### GetSwitchRefTypeOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetSwitchRefTypeOk() (*string, bool)`

GetSwitchRefTypeOk returns a tuple with the SwitchRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchRefType

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetSwitchRefType(v string)`

SetSwitchRefType sets SwitchRefType field to given value.

### HasSwitchRefType

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasSwitchRefType() bool`

HasSwitchRefType returns a boolean if a field has been set.

### GetConnectionService

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetConnectionService() string`

GetConnectionService returns the ConnectionService field if non-nil, zero value otherwise.

### GetConnectionServiceOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetConnectionServiceOk() (*string, bool)`

GetConnectionServiceOk returns a tuple with the ConnectionService field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnectionService

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetConnectionService(v string)`

SetConnectionService sets ConnectionService field to given value.

### HasConnectionService

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasConnectionService() bool`

HasConnectionService returns a boolean if a field has been set.

### GetConnectionServiceRefType

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetConnectionServiceRefType() string`

GetConnectionServiceRefType returns the ConnectionServiceRefType field if non-nil, zero value otherwise.

### GetConnectionServiceRefTypeOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetConnectionServiceRefTypeOk() (*string, bool)`

GetConnectionServiceRefTypeOk returns a tuple with the ConnectionServiceRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnectionServiceRefType

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetConnectionServiceRefType(v string)`

SetConnectionServiceRefType sets ConnectionServiceRefType field to given value.

### HasConnectionServiceRefType

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasConnectionServiceRefType() bool`

HasConnectionServiceRefType returns a boolean if a field has been set.

### GetPort

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetPort() string`

GetPort returns the Port field if non-nil, zero value otherwise.

### GetPortOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetPortOk() (*string, bool)`

GetPortOk returns a tuple with the Port field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPort

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetPort(v string)`

SetPort sets Port field to given value.

### HasPort

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasPort() bool`

HasPort returns a boolean if a field has been set.

### GetSfpMacAddressOrSn

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetSfpMacAddressOrSn() string`

GetSfpMacAddressOrSn returns the SfpMacAddressOrSn field if non-nil, zero value otherwise.

### GetSfpMacAddressOrSnOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetSfpMacAddressOrSnOk() (*string, bool)`

GetSfpMacAddressOrSnOk returns a tuple with the SfpMacAddressOrSn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSfpMacAddressOrSn

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetSfpMacAddressOrSn(v string)`

SetSfpMacAddressOrSn sets SfpMacAddressOrSn field to given value.

### HasSfpMacAddressOrSn

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasSfpMacAddressOrSn() bool`

HasSfpMacAddressOrSn returns a boolean if a field has been set.

### GetUsesTaggedPackets

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetUsesTaggedPackets() bool`

GetUsesTaggedPackets returns the UsesTaggedPackets field if non-nil, zero value otherwise.

### GetUsesTaggedPacketsOk

`func (o *DevicecontrollersPutRequestDeviceControllerValue) GetUsesTaggedPacketsOk() (*bool, bool)`

GetUsesTaggedPacketsOk returns a tuple with the UsesTaggedPackets field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsesTaggedPackets

`func (o *DevicecontrollersPutRequestDeviceControllerValue) SetUsesTaggedPackets(v bool)`

SetUsesTaggedPackets sets UsesTaggedPackets field to given value.

### HasUsesTaggedPackets

`func (o *DevicecontrollersPutRequestDeviceControllerValue) HasUsesTaggedPackets() bool`

HasUsesTaggedPackets returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


