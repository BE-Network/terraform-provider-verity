# BundlesPatchRequestEndpointBundleValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**DeviceSettings** | Pointer to **string** | Device Settings for device | [optional] [default to "eth_device_profile|(Device Settings)|"]
**DeviceSettingsRefType** | Pointer to **string** | Object type for device_settings field | [optional] 
**CliCommands** | Pointer to **string** | CLI Commands | [optional] [default to ""]
**EthPortPaths** | Pointer to [**[]BundlesPatchRequestEndpointBundleValueEthPortPathsInner**](BundlesPatchRequestEndpointBundleValueEthPortPathsInner.md) |  | [optional] 
**UserServices** | Pointer to [**[]BundlesPatchRequestEndpointBundleValueUserServicesInner**](BundlesPatchRequestEndpointBundleValueUserServicesInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**BundlesPatchRequestEndpointBundleValueObjectProperties**](BundlesPatchRequestEndpointBundleValueObjectProperties.md) |  | [optional] 
**RgServices** | Pointer to [**[]BundlesPatchRequestEndpointBundleValueRgServicesInner**](BundlesPatchRequestEndpointBundleValueRgServicesInner.md) |  | [optional] 
**Enable** | Pointer to **bool** | Enable object. It&#39;s highly recommended to set this value to true so that validation on the object will be ran. | [optional] [default to false]
**Protocol** | Pointer to **string** | Voice Protocol: MGCP or SIP | [optional] [default to "SIP"]
**DeviceVoiceSettings** | Pointer to **string** | Device Voice Settings for device | [optional] [default to "voice_device_profile|(SIP Voice Device)|"]
**DeviceVoiceSettingsRefType** | Pointer to **string** | Object type for device_voice_settings field | [optional] 
**VoicePortProfilePaths** | Pointer to [**[]BundlesPatchRequestEndpointBundleValueVoicePortProfilePathsInner**](BundlesPatchRequestEndpointBundleValueVoicePortProfilePathsInner.md) |  | [optional] 

## Methods

### NewBundlesPatchRequestEndpointBundleValue

`func NewBundlesPatchRequestEndpointBundleValue() *BundlesPatchRequestEndpointBundleValue`

NewBundlesPatchRequestEndpointBundleValue instantiates a new BundlesPatchRequestEndpointBundleValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBundlesPatchRequestEndpointBundleValueWithDefaults

`func NewBundlesPatchRequestEndpointBundleValueWithDefaults() *BundlesPatchRequestEndpointBundleValue`

NewBundlesPatchRequestEndpointBundleValueWithDefaults instantiates a new BundlesPatchRequestEndpointBundleValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *BundlesPatchRequestEndpointBundleValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *BundlesPatchRequestEndpointBundleValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *BundlesPatchRequestEndpointBundleValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *BundlesPatchRequestEndpointBundleValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetDeviceSettings

`func (o *BundlesPatchRequestEndpointBundleValue) GetDeviceSettings() string`

GetDeviceSettings returns the DeviceSettings field if non-nil, zero value otherwise.

### GetDeviceSettingsOk

`func (o *BundlesPatchRequestEndpointBundleValue) GetDeviceSettingsOk() (*string, bool)`

GetDeviceSettingsOk returns a tuple with the DeviceSettings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceSettings

`func (o *BundlesPatchRequestEndpointBundleValue) SetDeviceSettings(v string)`

SetDeviceSettings sets DeviceSettings field to given value.

### HasDeviceSettings

`func (o *BundlesPatchRequestEndpointBundleValue) HasDeviceSettings() bool`

HasDeviceSettings returns a boolean if a field has been set.

### GetDeviceSettingsRefType

`func (o *BundlesPatchRequestEndpointBundleValue) GetDeviceSettingsRefType() string`

GetDeviceSettingsRefType returns the DeviceSettingsRefType field if non-nil, zero value otherwise.

### GetDeviceSettingsRefTypeOk

`func (o *BundlesPatchRequestEndpointBundleValue) GetDeviceSettingsRefTypeOk() (*string, bool)`

GetDeviceSettingsRefTypeOk returns a tuple with the DeviceSettingsRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceSettingsRefType

`func (o *BundlesPatchRequestEndpointBundleValue) SetDeviceSettingsRefType(v string)`

SetDeviceSettingsRefType sets DeviceSettingsRefType field to given value.

### HasDeviceSettingsRefType

`func (o *BundlesPatchRequestEndpointBundleValue) HasDeviceSettingsRefType() bool`

HasDeviceSettingsRefType returns a boolean if a field has been set.

### GetCliCommands

`func (o *BundlesPatchRequestEndpointBundleValue) GetCliCommands() string`

GetCliCommands returns the CliCommands field if non-nil, zero value otherwise.

### GetCliCommandsOk

`func (o *BundlesPatchRequestEndpointBundleValue) GetCliCommandsOk() (*string, bool)`

GetCliCommandsOk returns a tuple with the CliCommands field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCliCommands

`func (o *BundlesPatchRequestEndpointBundleValue) SetCliCommands(v string)`

SetCliCommands sets CliCommands field to given value.

### HasCliCommands

`func (o *BundlesPatchRequestEndpointBundleValue) HasCliCommands() bool`

HasCliCommands returns a boolean if a field has been set.

### GetEthPortPaths

`func (o *BundlesPatchRequestEndpointBundleValue) GetEthPortPaths() []BundlesPatchRequestEndpointBundleValueEthPortPathsInner`

GetEthPortPaths returns the EthPortPaths field if non-nil, zero value otherwise.

### GetEthPortPathsOk

`func (o *BundlesPatchRequestEndpointBundleValue) GetEthPortPathsOk() (*[]BundlesPatchRequestEndpointBundleValueEthPortPathsInner, bool)`

GetEthPortPathsOk returns a tuple with the EthPortPaths field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEthPortPaths

`func (o *BundlesPatchRequestEndpointBundleValue) SetEthPortPaths(v []BundlesPatchRequestEndpointBundleValueEthPortPathsInner)`

SetEthPortPaths sets EthPortPaths field to given value.

### HasEthPortPaths

`func (o *BundlesPatchRequestEndpointBundleValue) HasEthPortPaths() bool`

HasEthPortPaths returns a boolean if a field has been set.

### GetUserServices

`func (o *BundlesPatchRequestEndpointBundleValue) GetUserServices() []BundlesPatchRequestEndpointBundleValueUserServicesInner`

GetUserServices returns the UserServices field if non-nil, zero value otherwise.

### GetUserServicesOk

`func (o *BundlesPatchRequestEndpointBundleValue) GetUserServicesOk() (*[]BundlesPatchRequestEndpointBundleValueUserServicesInner, bool)`

GetUserServicesOk returns a tuple with the UserServices field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserServices

`func (o *BundlesPatchRequestEndpointBundleValue) SetUserServices(v []BundlesPatchRequestEndpointBundleValueUserServicesInner)`

SetUserServices sets UserServices field to given value.

### HasUserServices

`func (o *BundlesPatchRequestEndpointBundleValue) HasUserServices() bool`

HasUserServices returns a boolean if a field has been set.

### GetObjectProperties

`func (o *BundlesPatchRequestEndpointBundleValue) GetObjectProperties() BundlesPatchRequestEndpointBundleValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *BundlesPatchRequestEndpointBundleValue) GetObjectPropertiesOk() (*BundlesPatchRequestEndpointBundleValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *BundlesPatchRequestEndpointBundleValue) SetObjectProperties(v BundlesPatchRequestEndpointBundleValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *BundlesPatchRequestEndpointBundleValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.

### GetRgServices

`func (o *BundlesPatchRequestEndpointBundleValue) GetRgServices() []BundlesPatchRequestEndpointBundleValueRgServicesInner`

GetRgServices returns the RgServices field if non-nil, zero value otherwise.

### GetRgServicesOk

`func (o *BundlesPatchRequestEndpointBundleValue) GetRgServicesOk() (*[]BundlesPatchRequestEndpointBundleValueRgServicesInner, bool)`

GetRgServicesOk returns a tuple with the RgServices field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRgServices

`func (o *BundlesPatchRequestEndpointBundleValue) SetRgServices(v []BundlesPatchRequestEndpointBundleValueRgServicesInner)`

SetRgServices sets RgServices field to given value.

### HasRgServices

`func (o *BundlesPatchRequestEndpointBundleValue) HasRgServices() bool`

HasRgServices returns a boolean if a field has been set.

### GetEnable

`func (o *BundlesPatchRequestEndpointBundleValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *BundlesPatchRequestEndpointBundleValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *BundlesPatchRequestEndpointBundleValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *BundlesPatchRequestEndpointBundleValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetProtocol

`func (o *BundlesPatchRequestEndpointBundleValue) GetProtocol() string`

GetProtocol returns the Protocol field if non-nil, zero value otherwise.

### GetProtocolOk

`func (o *BundlesPatchRequestEndpointBundleValue) GetProtocolOk() (*string, bool)`

GetProtocolOk returns a tuple with the Protocol field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProtocol

`func (o *BundlesPatchRequestEndpointBundleValue) SetProtocol(v string)`

SetProtocol sets Protocol field to given value.

### HasProtocol

`func (o *BundlesPatchRequestEndpointBundleValue) HasProtocol() bool`

HasProtocol returns a boolean if a field has been set.

### GetDeviceVoiceSettings

`func (o *BundlesPatchRequestEndpointBundleValue) GetDeviceVoiceSettings() string`

GetDeviceVoiceSettings returns the DeviceVoiceSettings field if non-nil, zero value otherwise.

### GetDeviceVoiceSettingsOk

`func (o *BundlesPatchRequestEndpointBundleValue) GetDeviceVoiceSettingsOk() (*string, bool)`

GetDeviceVoiceSettingsOk returns a tuple with the DeviceVoiceSettings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceVoiceSettings

`func (o *BundlesPatchRequestEndpointBundleValue) SetDeviceVoiceSettings(v string)`

SetDeviceVoiceSettings sets DeviceVoiceSettings field to given value.

### HasDeviceVoiceSettings

`func (o *BundlesPatchRequestEndpointBundleValue) HasDeviceVoiceSettings() bool`

HasDeviceVoiceSettings returns a boolean if a field has been set.

### GetDeviceVoiceSettingsRefType

`func (o *BundlesPatchRequestEndpointBundleValue) GetDeviceVoiceSettingsRefType() string`

GetDeviceVoiceSettingsRefType returns the DeviceVoiceSettingsRefType field if non-nil, zero value otherwise.

### GetDeviceVoiceSettingsRefTypeOk

`func (o *BundlesPatchRequestEndpointBundleValue) GetDeviceVoiceSettingsRefTypeOk() (*string, bool)`

GetDeviceVoiceSettingsRefTypeOk returns a tuple with the DeviceVoiceSettingsRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceVoiceSettingsRefType

`func (o *BundlesPatchRequestEndpointBundleValue) SetDeviceVoiceSettingsRefType(v string)`

SetDeviceVoiceSettingsRefType sets DeviceVoiceSettingsRefType field to given value.

### HasDeviceVoiceSettingsRefType

`func (o *BundlesPatchRequestEndpointBundleValue) HasDeviceVoiceSettingsRefType() bool`

HasDeviceVoiceSettingsRefType returns a boolean if a field has been set.

### GetVoicePortProfilePaths

`func (o *BundlesPatchRequestEndpointBundleValue) GetVoicePortProfilePaths() []BundlesPatchRequestEndpointBundleValueVoicePortProfilePathsInner`

GetVoicePortProfilePaths returns the VoicePortProfilePaths field if non-nil, zero value otherwise.

### GetVoicePortProfilePathsOk

`func (o *BundlesPatchRequestEndpointBundleValue) GetVoicePortProfilePathsOk() (*[]BundlesPatchRequestEndpointBundleValueVoicePortProfilePathsInner, bool)`

GetVoicePortProfilePathsOk returns a tuple with the VoicePortProfilePaths field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVoicePortProfilePaths

`func (o *BundlesPatchRequestEndpointBundleValue) SetVoicePortProfilePaths(v []BundlesPatchRequestEndpointBundleValueVoicePortProfilePathsInner)`

SetVoicePortProfilePaths sets VoicePortProfilePaths field to given value.

### HasVoicePortProfilePaths

`func (o *BundlesPatchRequestEndpointBundleValue) HasVoicePortProfilePaths() bool`

HasVoicePortProfilePaths returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


