# ConfigPutRequestGatewayProfileGatewayProfileName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**TenantSliceManaged** | Pointer to **bool** | Profiles that Tenant Slice creates and manages | [optional] [default to false]
**ExternalGateways** | Pointer to [**[]ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner**](ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties**](ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties.md) |  | [optional] 

## Methods

### NewConfigPutRequestGatewayProfileGatewayProfileName

`func NewConfigPutRequestGatewayProfileGatewayProfileName() *ConfigPutRequestGatewayProfileGatewayProfileName`

NewConfigPutRequestGatewayProfileGatewayProfileName instantiates a new ConfigPutRequestGatewayProfileGatewayProfileName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestGatewayProfileGatewayProfileNameWithDefaults

`func NewConfigPutRequestGatewayProfileGatewayProfileNameWithDefaults() *ConfigPutRequestGatewayProfileGatewayProfileName`

NewConfigPutRequestGatewayProfileGatewayProfileNameWithDefaults instantiates a new ConfigPutRequestGatewayProfileGatewayProfileName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestGatewayProfileGatewayProfileName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestGatewayProfileGatewayProfileName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestGatewayProfileGatewayProfileName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestGatewayProfileGatewayProfileName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ConfigPutRequestGatewayProfileGatewayProfileName) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestGatewayProfileGatewayProfileName) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestGatewayProfileGatewayProfileName) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestGatewayProfileGatewayProfileName) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetTenantSliceManaged

`func (o *ConfigPutRequestGatewayProfileGatewayProfileName) GetTenantSliceManaged() bool`

GetTenantSliceManaged returns the TenantSliceManaged field if non-nil, zero value otherwise.

### GetTenantSliceManagedOk

`func (o *ConfigPutRequestGatewayProfileGatewayProfileName) GetTenantSliceManagedOk() (*bool, bool)`

GetTenantSliceManagedOk returns a tuple with the TenantSliceManaged field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantSliceManaged

`func (o *ConfigPutRequestGatewayProfileGatewayProfileName) SetTenantSliceManaged(v bool)`

SetTenantSliceManaged sets TenantSliceManaged field to given value.

### HasTenantSliceManaged

`func (o *ConfigPutRequestGatewayProfileGatewayProfileName) HasTenantSliceManaged() bool`

HasTenantSliceManaged returns a boolean if a field has been set.

### GetExternalGateways

`func (o *ConfigPutRequestGatewayProfileGatewayProfileName) GetExternalGateways() []ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner`

GetExternalGateways returns the ExternalGateways field if non-nil, zero value otherwise.

### GetExternalGatewaysOk

`func (o *ConfigPutRequestGatewayProfileGatewayProfileName) GetExternalGatewaysOk() (*[]ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner, bool)`

GetExternalGatewaysOk returns a tuple with the ExternalGateways field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalGateways

`func (o *ConfigPutRequestGatewayProfileGatewayProfileName) SetExternalGateways(v []ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner)`

SetExternalGateways sets ExternalGateways field to given value.

### HasExternalGateways

`func (o *ConfigPutRequestGatewayProfileGatewayProfileName) HasExternalGateways() bool`

HasExternalGateways returns a boolean if a field has been set.

### GetObjectProperties

`func (o *ConfigPutRequestGatewayProfileGatewayProfileName) GetObjectProperties() ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ConfigPutRequestGatewayProfileGatewayProfileName) GetObjectPropertiesOk() (*ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ConfigPutRequestGatewayProfileGatewayProfileName) SetObjectProperties(v ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ConfigPutRequestGatewayProfileGatewayProfileName) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


