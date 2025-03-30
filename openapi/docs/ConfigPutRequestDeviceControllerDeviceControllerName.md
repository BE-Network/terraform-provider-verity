# ConfigPutRequestDeviceControllerDeviceControllerName

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
**CommType** | Pointer to **string** | Comm Type | [optional] [default to "gnmi"]
**SnmpCommunityString** | Pointer to **string** | Comm Credentials | [optional] [default to "private"]
**UplinkPort** | Pointer to **string** | Uplink Port of Managed Device | [optional] [default to ""]
**LldpSearchString** | Pointer to **string** | The unique identifier associated with the managed device. | [optional] [default to ""]
**ZtpIdentification** | Pointer to **string** | Service Tag or Serial Number to identify device for Zero Touch Provisioning | [optional] [default to ""]
**LocatedBy** | Pointer to **string** | Controls how the system locates this Device within its LAN | [optional] [default to "LLDP"]
**PowerState** | Pointer to **string** | Power state of Switch Controller | [optional] [default to "on"]
**CommunicationMode** | Pointer to **string** | Communication Mode | [optional] [default to "sonic"]
**CliAccessMode** | Pointer to **string** | CLI Access Mode | [optional] [default to "SSH"]
**Username** | Pointer to **string** | Username | [optional] [default to ""]
**Password** | Pointer to **string** | Password | [optional] [default to ""]
**EnablePassword** | Pointer to **string** | Enable Password - to enable privileged CLI operations | [optional] [default to ""]
**SshKeyOrPassword** | Pointer to **string** | SSH Key or Password | [optional] [default to ""]
**ManagedOnNativeVlan** | Pointer to **bool** | Managed on native VLAN | [optional] [default to true]
**Sdlc** | Pointer to **string** | SDLC that Device Controller belongs to | [optional] [default to ""]
**Switchpoint** | Pointer to **string** | Switchpoint reference | [optional] [default to ""]
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

## Methods

### NewConfigPutRequestDeviceControllerDeviceControllerName

`func NewConfigPutRequestDeviceControllerDeviceControllerName() *ConfigPutRequestDeviceControllerDeviceControllerName`

NewConfigPutRequestDeviceControllerDeviceControllerName instantiates a new ConfigPutRequestDeviceControllerDeviceControllerName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestDeviceControllerDeviceControllerNameWithDefaults

`func NewConfigPutRequestDeviceControllerDeviceControllerNameWithDefaults() *ConfigPutRequestDeviceControllerDeviceControllerName`

NewConfigPutRequestDeviceControllerDeviceControllerNameWithDefaults instantiates a new ConfigPutRequestDeviceControllerDeviceControllerName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetIpSource

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetIpSource() string`

GetIpSource returns the IpSource field if non-nil, zero value otherwise.

### GetIpSourceOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetIpSourceOk() (*string, bool)`

GetIpSourceOk returns a tuple with the IpSource field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpSource

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetIpSource(v string)`

SetIpSource sets IpSource field to given value.

### HasIpSource

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasIpSource() bool`

HasIpSource returns a boolean if a field has been set.

### GetControllerIpAndMask

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetControllerIpAndMask() string`

GetControllerIpAndMask returns the ControllerIpAndMask field if non-nil, zero value otherwise.

### GetControllerIpAndMaskOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetControllerIpAndMaskOk() (*string, bool)`

GetControllerIpAndMaskOk returns a tuple with the ControllerIpAndMask field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetControllerIpAndMask

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetControllerIpAndMask(v string)`

SetControllerIpAndMask sets ControllerIpAndMask field to given value.

### HasControllerIpAndMask

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasControllerIpAndMask() bool`

HasControllerIpAndMask returns a boolean if a field has been set.

### GetGateway

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetGateway() string`

GetGateway returns the Gateway field if non-nil, zero value otherwise.

### GetGatewayOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetGatewayOk() (*string, bool)`

GetGatewayOk returns a tuple with the Gateway field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGateway

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetGateway(v string)`

SetGateway sets Gateway field to given value.

### HasGateway

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasGateway() bool`

HasGateway returns a boolean if a field has been set.

### GetSwitchIpAndMask

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetSwitchIpAndMask() string`

GetSwitchIpAndMask returns the SwitchIpAndMask field if non-nil, zero value otherwise.

### GetSwitchIpAndMaskOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetSwitchIpAndMaskOk() (*string, bool)`

GetSwitchIpAndMaskOk returns a tuple with the SwitchIpAndMask field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchIpAndMask

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetSwitchIpAndMask(v string)`

SetSwitchIpAndMask sets SwitchIpAndMask field to given value.

### HasSwitchIpAndMask

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasSwitchIpAndMask() bool`

HasSwitchIpAndMask returns a boolean if a field has been set.

### GetSwitchGateway

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetSwitchGateway() string`

GetSwitchGateway returns the SwitchGateway field if non-nil, zero value otherwise.

### GetSwitchGatewayOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetSwitchGatewayOk() (*string, bool)`

GetSwitchGatewayOk returns a tuple with the SwitchGateway field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchGateway

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetSwitchGateway(v string)`

SetSwitchGateway sets SwitchGateway field to given value.

### HasSwitchGateway

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasSwitchGateway() bool`

HasSwitchGateway returns a boolean if a field has been set.

### GetCommType

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetCommType() string`

GetCommType returns the CommType field if non-nil, zero value otherwise.

### GetCommTypeOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetCommTypeOk() (*string, bool)`

GetCommTypeOk returns a tuple with the CommType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommType

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetCommType(v string)`

SetCommType sets CommType field to given value.

### HasCommType

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasCommType() bool`

HasCommType returns a boolean if a field has been set.

### GetSnmpCommunityString

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetSnmpCommunityString() string`

GetSnmpCommunityString returns the SnmpCommunityString field if non-nil, zero value otherwise.

### GetSnmpCommunityStringOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetSnmpCommunityStringOk() (*string, bool)`

GetSnmpCommunityStringOk returns a tuple with the SnmpCommunityString field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSnmpCommunityString

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetSnmpCommunityString(v string)`

SetSnmpCommunityString sets SnmpCommunityString field to given value.

### HasSnmpCommunityString

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasSnmpCommunityString() bool`

HasSnmpCommunityString returns a boolean if a field has been set.

### GetUplinkPort

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetUplinkPort() string`

GetUplinkPort returns the UplinkPort field if non-nil, zero value otherwise.

### GetUplinkPortOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetUplinkPortOk() (*string, bool)`

GetUplinkPortOk returns a tuple with the UplinkPort field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUplinkPort

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetUplinkPort(v string)`

SetUplinkPort sets UplinkPort field to given value.

### HasUplinkPort

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasUplinkPort() bool`

HasUplinkPort returns a boolean if a field has been set.

### GetLldpSearchString

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetLldpSearchString() string`

GetLldpSearchString returns the LldpSearchString field if non-nil, zero value otherwise.

### GetLldpSearchStringOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetLldpSearchStringOk() (*string, bool)`

GetLldpSearchStringOk returns a tuple with the LldpSearchString field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLldpSearchString

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetLldpSearchString(v string)`

SetLldpSearchString sets LldpSearchString field to given value.

### HasLldpSearchString

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasLldpSearchString() bool`

HasLldpSearchString returns a boolean if a field has been set.

### GetZtpIdentification

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetZtpIdentification() string`

GetZtpIdentification returns the ZtpIdentification field if non-nil, zero value otherwise.

### GetZtpIdentificationOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetZtpIdentificationOk() (*string, bool)`

GetZtpIdentificationOk returns a tuple with the ZtpIdentification field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetZtpIdentification

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetZtpIdentification(v string)`

SetZtpIdentification sets ZtpIdentification field to given value.

### HasZtpIdentification

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasZtpIdentification() bool`

HasZtpIdentification returns a boolean if a field has been set.

### GetLocatedBy

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetLocatedBy() string`

GetLocatedBy returns the LocatedBy field if non-nil, zero value otherwise.

### GetLocatedByOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetLocatedByOk() (*string, bool)`

GetLocatedByOk returns a tuple with the LocatedBy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocatedBy

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetLocatedBy(v string)`

SetLocatedBy sets LocatedBy field to given value.

### HasLocatedBy

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasLocatedBy() bool`

HasLocatedBy returns a boolean if a field has been set.

### GetPowerState

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetPowerState() string`

GetPowerState returns the PowerState field if non-nil, zero value otherwise.

### GetPowerStateOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetPowerStateOk() (*string, bool)`

GetPowerStateOk returns a tuple with the PowerState field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPowerState

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetPowerState(v string)`

SetPowerState sets PowerState field to given value.

### HasPowerState

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasPowerState() bool`

HasPowerState returns a boolean if a field has been set.

### GetCommunicationMode

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetCommunicationMode() string`

GetCommunicationMode returns the CommunicationMode field if non-nil, zero value otherwise.

### GetCommunicationModeOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetCommunicationModeOk() (*string, bool)`

GetCommunicationModeOk returns a tuple with the CommunicationMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommunicationMode

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetCommunicationMode(v string)`

SetCommunicationMode sets CommunicationMode field to given value.

### HasCommunicationMode

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasCommunicationMode() bool`

HasCommunicationMode returns a boolean if a field has been set.

### GetCliAccessMode

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetCliAccessMode() string`

GetCliAccessMode returns the CliAccessMode field if non-nil, zero value otherwise.

### GetCliAccessModeOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetCliAccessModeOk() (*string, bool)`

GetCliAccessModeOk returns a tuple with the CliAccessMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCliAccessMode

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetCliAccessMode(v string)`

SetCliAccessMode sets CliAccessMode field to given value.

### HasCliAccessMode

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasCliAccessMode() bool`

HasCliAccessMode returns a boolean if a field has been set.

### GetUsername

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsername

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetUsername(v string)`

SetUsername sets Username field to given value.

### HasUsername

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasUsername() bool`

HasUsername returns a boolean if a field has been set.

### GetPassword

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetPassword(v string)`

SetPassword sets Password field to given value.

### HasPassword

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasPassword() bool`

HasPassword returns a boolean if a field has been set.

### GetEnablePassword

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetEnablePassword() string`

GetEnablePassword returns the EnablePassword field if non-nil, zero value otherwise.

### GetEnablePasswordOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetEnablePasswordOk() (*string, bool)`

GetEnablePasswordOk returns a tuple with the EnablePassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnablePassword

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetEnablePassword(v string)`

SetEnablePassword sets EnablePassword field to given value.

### HasEnablePassword

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasEnablePassword() bool`

HasEnablePassword returns a boolean if a field has been set.

### GetSshKeyOrPassword

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetSshKeyOrPassword() string`

GetSshKeyOrPassword returns the SshKeyOrPassword field if non-nil, zero value otherwise.

### GetSshKeyOrPasswordOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetSshKeyOrPasswordOk() (*string, bool)`

GetSshKeyOrPasswordOk returns a tuple with the SshKeyOrPassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSshKeyOrPassword

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetSshKeyOrPassword(v string)`

SetSshKeyOrPassword sets SshKeyOrPassword field to given value.

### HasSshKeyOrPassword

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasSshKeyOrPassword() bool`

HasSshKeyOrPassword returns a boolean if a field has been set.

### GetManagedOnNativeVlan

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetManagedOnNativeVlan() bool`

GetManagedOnNativeVlan returns the ManagedOnNativeVlan field if non-nil, zero value otherwise.

### GetManagedOnNativeVlanOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetManagedOnNativeVlanOk() (*bool, bool)`

GetManagedOnNativeVlanOk returns a tuple with the ManagedOnNativeVlan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetManagedOnNativeVlan

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetManagedOnNativeVlan(v bool)`

SetManagedOnNativeVlan sets ManagedOnNativeVlan field to given value.

### HasManagedOnNativeVlan

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasManagedOnNativeVlan() bool`

HasManagedOnNativeVlan returns a boolean if a field has been set.

### GetSdlc

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetSdlc() string`

GetSdlc returns the Sdlc field if non-nil, zero value otherwise.

### GetSdlcOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetSdlcOk() (*string, bool)`

GetSdlcOk returns a tuple with the Sdlc field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSdlc

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetSdlc(v string)`

SetSdlc sets Sdlc field to given value.

### HasSdlc

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasSdlc() bool`

HasSdlc returns a boolean if a field has been set.

### GetSwitchpoint

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetSwitchpoint() string`

GetSwitchpoint returns the Switchpoint field if non-nil, zero value otherwise.

### GetSwitchpointOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetSwitchpointOk() (*string, bool)`

GetSwitchpointOk returns a tuple with the Switchpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchpoint

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetSwitchpoint(v string)`

SetSwitchpoint sets Switchpoint field to given value.

### HasSwitchpoint

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasSwitchpoint() bool`

HasSwitchpoint returns a boolean if a field has been set.

### GetSwitchpointRefType

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetSwitchpointRefType() string`

GetSwitchpointRefType returns the SwitchpointRefType field if non-nil, zero value otherwise.

### GetSwitchpointRefTypeOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetSwitchpointRefTypeOk() (*string, bool)`

GetSwitchpointRefTypeOk returns a tuple with the SwitchpointRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchpointRefType

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetSwitchpointRefType(v string)`

SetSwitchpointRefType sets SwitchpointRefType field to given value.

### HasSwitchpointRefType

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasSwitchpointRefType() bool`

HasSwitchpointRefType returns a boolean if a field has been set.

### GetSecurityType

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetSecurityType() string`

GetSecurityType returns the SecurityType field if non-nil, zero value otherwise.

### GetSecurityTypeOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetSecurityTypeOk() (*string, bool)`

GetSecurityTypeOk returns a tuple with the SecurityType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSecurityType

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetSecurityType(v string)`

SetSecurityType sets SecurityType field to given value.

### HasSecurityType

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasSecurityType() bool`

HasSecurityType returns a boolean if a field has been set.

### GetSnmpv3Username

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetSnmpv3Username() string`

GetSnmpv3Username returns the Snmpv3Username field if non-nil, zero value otherwise.

### GetSnmpv3UsernameOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetSnmpv3UsernameOk() (*string, bool)`

GetSnmpv3UsernameOk returns a tuple with the Snmpv3Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSnmpv3Username

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetSnmpv3Username(v string)`

SetSnmpv3Username sets Snmpv3Username field to given value.

### HasSnmpv3Username

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasSnmpv3Username() bool`

HasSnmpv3Username returns a boolean if a field has been set.

### GetAuthenticationProtocol

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetAuthenticationProtocol() string`

GetAuthenticationProtocol returns the AuthenticationProtocol field if non-nil, zero value otherwise.

### GetAuthenticationProtocolOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetAuthenticationProtocolOk() (*string, bool)`

GetAuthenticationProtocolOk returns a tuple with the AuthenticationProtocol field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthenticationProtocol

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetAuthenticationProtocol(v string)`

SetAuthenticationProtocol sets AuthenticationProtocol field to given value.

### HasAuthenticationProtocol

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasAuthenticationProtocol() bool`

HasAuthenticationProtocol returns a boolean if a field has been set.

### GetPassphrase

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetPassphrase() string`

GetPassphrase returns the Passphrase field if non-nil, zero value otherwise.

### GetPassphraseOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetPassphraseOk() (*string, bool)`

GetPassphraseOk returns a tuple with the Passphrase field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassphrase

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetPassphrase(v string)`

SetPassphrase sets Passphrase field to given value.

### HasPassphrase

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasPassphrase() bool`

HasPassphrase returns a boolean if a field has been set.

### GetPrivateProtocol

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetPrivateProtocol() string`

GetPrivateProtocol returns the PrivateProtocol field if non-nil, zero value otherwise.

### GetPrivateProtocolOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetPrivateProtocolOk() (*string, bool)`

GetPrivateProtocolOk returns a tuple with the PrivateProtocol field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrivateProtocol

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetPrivateProtocol(v string)`

SetPrivateProtocol sets PrivateProtocol field to given value.

### HasPrivateProtocol

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasPrivateProtocol() bool`

HasPrivateProtocol returns a boolean if a field has been set.

### GetPrivatePassword

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetPrivatePassword() string`

GetPrivatePassword returns the PrivatePassword field if non-nil, zero value otherwise.

### GetPrivatePasswordOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetPrivatePasswordOk() (*string, bool)`

GetPrivatePasswordOk returns a tuple with the PrivatePassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrivatePassword

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetPrivatePassword(v string)`

SetPrivatePassword sets PrivatePassword field to given value.

### HasPrivatePassword

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasPrivatePassword() bool`

HasPrivatePassword returns a boolean if a field has been set.

### GetPasswordEncrypted

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetPasswordEncrypted() string`

GetPasswordEncrypted returns the PasswordEncrypted field if non-nil, zero value otherwise.

### GetPasswordEncryptedOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetPasswordEncryptedOk() (*string, bool)`

GetPasswordEncryptedOk returns a tuple with the PasswordEncrypted field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPasswordEncrypted

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetPasswordEncrypted(v string)`

SetPasswordEncrypted sets PasswordEncrypted field to given value.

### HasPasswordEncrypted

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasPasswordEncrypted() bool`

HasPasswordEncrypted returns a boolean if a field has been set.

### GetEnablePasswordEncrypted

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetEnablePasswordEncrypted() string`

GetEnablePasswordEncrypted returns the EnablePasswordEncrypted field if non-nil, zero value otherwise.

### GetEnablePasswordEncryptedOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetEnablePasswordEncryptedOk() (*string, bool)`

GetEnablePasswordEncryptedOk returns a tuple with the EnablePasswordEncrypted field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnablePasswordEncrypted

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetEnablePasswordEncrypted(v string)`

SetEnablePasswordEncrypted sets EnablePasswordEncrypted field to given value.

### HasEnablePasswordEncrypted

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasEnablePasswordEncrypted() bool`

HasEnablePasswordEncrypted returns a boolean if a field has been set.

### GetSshKeyOrPasswordEncrypted

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetSshKeyOrPasswordEncrypted() string`

GetSshKeyOrPasswordEncrypted returns the SshKeyOrPasswordEncrypted field if non-nil, zero value otherwise.

### GetSshKeyOrPasswordEncryptedOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetSshKeyOrPasswordEncryptedOk() (*string, bool)`

GetSshKeyOrPasswordEncryptedOk returns a tuple with the SshKeyOrPasswordEncrypted field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSshKeyOrPasswordEncrypted

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetSshKeyOrPasswordEncrypted(v string)`

SetSshKeyOrPasswordEncrypted sets SshKeyOrPasswordEncrypted field to given value.

### HasSshKeyOrPasswordEncrypted

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasSshKeyOrPasswordEncrypted() bool`

HasSshKeyOrPasswordEncrypted returns a boolean if a field has been set.

### GetPassphraseEncrypted

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetPassphraseEncrypted() string`

GetPassphraseEncrypted returns the PassphraseEncrypted field if non-nil, zero value otherwise.

### GetPassphraseEncryptedOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetPassphraseEncryptedOk() (*string, bool)`

GetPassphraseEncryptedOk returns a tuple with the PassphraseEncrypted field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassphraseEncrypted

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetPassphraseEncrypted(v string)`

SetPassphraseEncrypted sets PassphraseEncrypted field to given value.

### HasPassphraseEncrypted

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasPassphraseEncrypted() bool`

HasPassphraseEncrypted returns a boolean if a field has been set.

### GetPrivatePasswordEncrypted

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetPrivatePasswordEncrypted() string`

GetPrivatePasswordEncrypted returns the PrivatePasswordEncrypted field if non-nil, zero value otherwise.

### GetPrivatePasswordEncryptedOk

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) GetPrivatePasswordEncryptedOk() (*string, bool)`

GetPrivatePasswordEncryptedOk returns a tuple with the PrivatePasswordEncrypted field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrivatePasswordEncrypted

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) SetPrivatePasswordEncrypted(v string)`

SetPrivatePasswordEncrypted sets PrivatePasswordEncrypted field to given value.

### HasPrivatePasswordEncrypted

`func (o *ConfigPutRequestDeviceControllerDeviceControllerName) HasPrivatePasswordEncrypted() bool`

HasPrivatePasswordEncrypted returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


