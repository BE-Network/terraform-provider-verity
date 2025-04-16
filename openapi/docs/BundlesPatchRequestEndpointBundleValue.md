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


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


