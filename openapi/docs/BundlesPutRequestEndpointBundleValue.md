# BundlesPutRequestEndpointBundleValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. It&#39;s highly recommended to set this value to true so that validation on the object will be ran. | [optional] [default to false]
**Protocol** | Pointer to **string** | Voice Protocol: MGCP or SIP | [optional] [default to "SIP"]
**DeviceSettings** | Pointer to **string** | Device Settings for device | [optional] [default to "eth_device_profile|(Device Settings)|"]
**DeviceSettingsRefType** | Pointer to **string** | Object type for device_settings field | [optional] 
**CliCommands** | Pointer to **string** | CLI Commands | [optional] [default to ""]
**DiagnosticsProfile** | Pointer to **string** | Diagnostics Profile for device | [optional] [default to ""]
**DiagnosticsProfileRefType** | Pointer to **string** | Object type for diagnostics_profile field | [optional] 
**EthPortPaths** | Pointer to [**[]BundlesPutRequestEndpointBundleValueEthPortPathsInner**](BundlesPutRequestEndpointBundleValueEthPortPathsInner.md) |  | [optional] 
**UserServices** | Pointer to [**[]BundlesPutRequestEndpointBundleValueUserServicesInner**](BundlesPutRequestEndpointBundleValueUserServicesInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**BundlesPutRequestEndpointBundleValueObjectProperties**](BundlesPutRequestEndpointBundleValueObjectProperties.md) |  | [optional] 
**DeviceVoiceSettings** | Pointer to **string** | Device Voice Settings for device | [optional] [default to "voice_device_profile|(SIP Voice Device)|"]
**DeviceVoiceSettingsRefType** | Pointer to **string** | Object type for device_voice_settings field | [optional] 
**VoicePortProfilePaths** | Pointer to [**[]BundlesPutRequestEndpointBundleValueVoicePortProfilePathsInner**](BundlesPutRequestEndpointBundleValueVoicePortProfilePathsInner.md) |  | [optional] 

## Methods

### NewBundlesPutRequestEndpointBundleValue

`func NewBundlesPutRequestEndpointBundleValue() *BundlesPutRequestEndpointBundleValue`

NewBundlesPutRequestEndpointBundleValue instantiates a new BundlesPutRequestEndpointBundleValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBundlesPutRequestEndpointBundleValueWithDefaults

`func NewBundlesPutRequestEndpointBundleValueWithDefaults() *BundlesPutRequestEndpointBundleValue`

NewBundlesPutRequestEndpointBundleValueWithDefaults instantiates a new BundlesPutRequestEndpointBundleValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *BundlesPutRequestEndpointBundleValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *BundlesPutRequestEndpointBundleValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *BundlesPutRequestEndpointBundleValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *BundlesPutRequestEndpointBundleValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *BundlesPutRequestEndpointBundleValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *BundlesPutRequestEndpointBundleValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *BundlesPutRequestEndpointBundleValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *BundlesPutRequestEndpointBundleValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetProtocol

`func (o *BundlesPutRequestEndpointBundleValue) GetProtocol() string`

GetProtocol returns the Protocol field if non-nil, zero value otherwise.

### GetProtocolOk

`func (o *BundlesPutRequestEndpointBundleValue) GetProtocolOk() (*string, bool)`

GetProtocolOk returns a tuple with the Protocol field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProtocol

`func (o *BundlesPutRequestEndpointBundleValue) SetProtocol(v string)`

SetProtocol sets Protocol field to given value.

### HasProtocol

`func (o *BundlesPutRequestEndpointBundleValue) HasProtocol() bool`

HasProtocol returns a boolean if a field has been set.

### GetDeviceSettings

`func (o *BundlesPutRequestEndpointBundleValue) GetDeviceSettings() string`

GetDeviceSettings returns the DeviceSettings field if non-nil, zero value otherwise.

### GetDeviceSettingsOk

`func (o *BundlesPutRequestEndpointBundleValue) GetDeviceSettingsOk() (*string, bool)`

GetDeviceSettingsOk returns a tuple with the DeviceSettings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceSettings

`func (o *BundlesPutRequestEndpointBundleValue) SetDeviceSettings(v string)`

SetDeviceSettings sets DeviceSettings field to given value.

### HasDeviceSettings

`func (o *BundlesPutRequestEndpointBundleValue) HasDeviceSettings() bool`

HasDeviceSettings returns a boolean if a field has been set.

### GetDeviceSettingsRefType

`func (o *BundlesPutRequestEndpointBundleValue) GetDeviceSettingsRefType() string`

GetDeviceSettingsRefType returns the DeviceSettingsRefType field if non-nil, zero value otherwise.

### GetDeviceSettingsRefTypeOk

`func (o *BundlesPutRequestEndpointBundleValue) GetDeviceSettingsRefTypeOk() (*string, bool)`

GetDeviceSettingsRefTypeOk returns a tuple with the DeviceSettingsRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceSettingsRefType

`func (o *BundlesPutRequestEndpointBundleValue) SetDeviceSettingsRefType(v string)`

SetDeviceSettingsRefType sets DeviceSettingsRefType field to given value.

### HasDeviceSettingsRefType

`func (o *BundlesPutRequestEndpointBundleValue) HasDeviceSettingsRefType() bool`

HasDeviceSettingsRefType returns a boolean if a field has been set.

### GetCliCommands

`func (o *BundlesPutRequestEndpointBundleValue) GetCliCommands() string`

GetCliCommands returns the CliCommands field if non-nil, zero value otherwise.

### GetCliCommandsOk

`func (o *BundlesPutRequestEndpointBundleValue) GetCliCommandsOk() (*string, bool)`

GetCliCommandsOk returns a tuple with the CliCommands field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCliCommands

`func (o *BundlesPutRequestEndpointBundleValue) SetCliCommands(v string)`

SetCliCommands sets CliCommands field to given value.

### HasCliCommands

`func (o *BundlesPutRequestEndpointBundleValue) HasCliCommands() bool`

HasCliCommands returns a boolean if a field has been set.

### GetDiagnosticsProfile

`func (o *BundlesPutRequestEndpointBundleValue) GetDiagnosticsProfile() string`

GetDiagnosticsProfile returns the DiagnosticsProfile field if non-nil, zero value otherwise.

### GetDiagnosticsProfileOk

`func (o *BundlesPutRequestEndpointBundleValue) GetDiagnosticsProfileOk() (*string, bool)`

GetDiagnosticsProfileOk returns a tuple with the DiagnosticsProfile field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDiagnosticsProfile

`func (o *BundlesPutRequestEndpointBundleValue) SetDiagnosticsProfile(v string)`

SetDiagnosticsProfile sets DiagnosticsProfile field to given value.

### HasDiagnosticsProfile

`func (o *BundlesPutRequestEndpointBundleValue) HasDiagnosticsProfile() bool`

HasDiagnosticsProfile returns a boolean if a field has been set.

### GetDiagnosticsProfileRefType

`func (o *BundlesPutRequestEndpointBundleValue) GetDiagnosticsProfileRefType() string`

GetDiagnosticsProfileRefType returns the DiagnosticsProfileRefType field if non-nil, zero value otherwise.

### GetDiagnosticsProfileRefTypeOk

`func (o *BundlesPutRequestEndpointBundleValue) GetDiagnosticsProfileRefTypeOk() (*string, bool)`

GetDiagnosticsProfileRefTypeOk returns a tuple with the DiagnosticsProfileRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDiagnosticsProfileRefType

`func (o *BundlesPutRequestEndpointBundleValue) SetDiagnosticsProfileRefType(v string)`

SetDiagnosticsProfileRefType sets DiagnosticsProfileRefType field to given value.

### HasDiagnosticsProfileRefType

`func (o *BundlesPutRequestEndpointBundleValue) HasDiagnosticsProfileRefType() bool`

HasDiagnosticsProfileRefType returns a boolean if a field has been set.

### GetEthPortPaths

`func (o *BundlesPutRequestEndpointBundleValue) GetEthPortPaths() []BundlesPutRequestEndpointBundleValueEthPortPathsInner`

GetEthPortPaths returns the EthPortPaths field if non-nil, zero value otherwise.

### GetEthPortPathsOk

`func (o *BundlesPutRequestEndpointBundleValue) GetEthPortPathsOk() (*[]BundlesPutRequestEndpointBundleValueEthPortPathsInner, bool)`

GetEthPortPathsOk returns a tuple with the EthPortPaths field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEthPortPaths

`func (o *BundlesPutRequestEndpointBundleValue) SetEthPortPaths(v []BundlesPutRequestEndpointBundleValueEthPortPathsInner)`

SetEthPortPaths sets EthPortPaths field to given value.

### HasEthPortPaths

`func (o *BundlesPutRequestEndpointBundleValue) HasEthPortPaths() bool`

HasEthPortPaths returns a boolean if a field has been set.

### GetUserServices

`func (o *BundlesPutRequestEndpointBundleValue) GetUserServices() []BundlesPutRequestEndpointBundleValueUserServicesInner`

GetUserServices returns the UserServices field if non-nil, zero value otherwise.

### GetUserServicesOk

`func (o *BundlesPutRequestEndpointBundleValue) GetUserServicesOk() (*[]BundlesPutRequestEndpointBundleValueUserServicesInner, bool)`

GetUserServicesOk returns a tuple with the UserServices field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserServices

`func (o *BundlesPutRequestEndpointBundleValue) SetUserServices(v []BundlesPutRequestEndpointBundleValueUserServicesInner)`

SetUserServices sets UserServices field to given value.

### HasUserServices

`func (o *BundlesPutRequestEndpointBundleValue) HasUserServices() bool`

HasUserServices returns a boolean if a field has been set.

### GetObjectProperties

`func (o *BundlesPutRequestEndpointBundleValue) GetObjectProperties() BundlesPutRequestEndpointBundleValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *BundlesPutRequestEndpointBundleValue) GetObjectPropertiesOk() (*BundlesPutRequestEndpointBundleValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *BundlesPutRequestEndpointBundleValue) SetObjectProperties(v BundlesPutRequestEndpointBundleValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *BundlesPutRequestEndpointBundleValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.

### GetDeviceVoiceSettings

`func (o *BundlesPutRequestEndpointBundleValue) GetDeviceVoiceSettings() string`

GetDeviceVoiceSettings returns the DeviceVoiceSettings field if non-nil, zero value otherwise.

### GetDeviceVoiceSettingsOk

`func (o *BundlesPutRequestEndpointBundleValue) GetDeviceVoiceSettingsOk() (*string, bool)`

GetDeviceVoiceSettingsOk returns a tuple with the DeviceVoiceSettings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceVoiceSettings

`func (o *BundlesPutRequestEndpointBundleValue) SetDeviceVoiceSettings(v string)`

SetDeviceVoiceSettings sets DeviceVoiceSettings field to given value.

### HasDeviceVoiceSettings

`func (o *BundlesPutRequestEndpointBundleValue) HasDeviceVoiceSettings() bool`

HasDeviceVoiceSettings returns a boolean if a field has been set.

### GetDeviceVoiceSettingsRefType

`func (o *BundlesPutRequestEndpointBundleValue) GetDeviceVoiceSettingsRefType() string`

GetDeviceVoiceSettingsRefType returns the DeviceVoiceSettingsRefType field if non-nil, zero value otherwise.

### GetDeviceVoiceSettingsRefTypeOk

`func (o *BundlesPutRequestEndpointBundleValue) GetDeviceVoiceSettingsRefTypeOk() (*string, bool)`

GetDeviceVoiceSettingsRefTypeOk returns a tuple with the DeviceVoiceSettingsRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceVoiceSettingsRefType

`func (o *BundlesPutRequestEndpointBundleValue) SetDeviceVoiceSettingsRefType(v string)`

SetDeviceVoiceSettingsRefType sets DeviceVoiceSettingsRefType field to given value.

### HasDeviceVoiceSettingsRefType

`func (o *BundlesPutRequestEndpointBundleValue) HasDeviceVoiceSettingsRefType() bool`

HasDeviceVoiceSettingsRefType returns a boolean if a field has been set.

### GetVoicePortProfilePaths

`func (o *BundlesPutRequestEndpointBundleValue) GetVoicePortProfilePaths() []BundlesPutRequestEndpointBundleValueVoicePortProfilePathsInner`

GetVoicePortProfilePaths returns the VoicePortProfilePaths field if non-nil, zero value otherwise.

### GetVoicePortProfilePathsOk

`func (o *BundlesPutRequestEndpointBundleValue) GetVoicePortProfilePathsOk() (*[]BundlesPutRequestEndpointBundleValueVoicePortProfilePathsInner, bool)`

GetVoicePortProfilePathsOk returns a tuple with the VoicePortProfilePaths field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVoicePortProfilePaths

`func (o *BundlesPutRequestEndpointBundleValue) SetVoicePortProfilePaths(v []BundlesPutRequestEndpointBundleValueVoicePortProfilePathsInner)`

SetVoicePortProfilePaths sets VoicePortProfilePaths field to given value.

### HasVoicePortProfilePaths

`func (o *BundlesPutRequestEndpointBundleValue) HasVoicePortProfilePaths() bool`

HasVoicePortProfilePaths returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


