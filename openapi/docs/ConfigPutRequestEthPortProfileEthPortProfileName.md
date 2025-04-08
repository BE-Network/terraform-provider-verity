# ConfigPutRequestEthPortProfileEthPortProfileName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. It&#39;s highly recommended to set this value to true so that validation on the object will be ran. | [optional] [default to false]
**TenantSliceManaged** | Pointer to **bool** | Profiles that Tenant Slice creates and manages | [optional] [default to false]
**Services** | Pointer to [**[]ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner**](ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**ConfigPutRequestEthPortProfileEthPortProfileNameObjectProperties**](ConfigPutRequestEthPortProfileEthPortProfileNameObjectProperties.md) |  | [optional] 

## Methods

### NewConfigPutRequestEthPortProfileEthPortProfileName

`func NewConfigPutRequestEthPortProfileEthPortProfileName() *ConfigPutRequestEthPortProfileEthPortProfileName`

NewConfigPutRequestEthPortProfileEthPortProfileName instantiates a new ConfigPutRequestEthPortProfileEthPortProfileName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestEthPortProfileEthPortProfileNameWithDefaults

`func NewConfigPutRequestEthPortProfileEthPortProfileNameWithDefaults() *ConfigPutRequestEthPortProfileEthPortProfileName`

NewConfigPutRequestEthPortProfileEthPortProfileNameWithDefaults instantiates a new ConfigPutRequestEthPortProfileEthPortProfileName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestEthPortProfileEthPortProfileName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestEthPortProfileEthPortProfileName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestEthPortProfileEthPortProfileName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestEthPortProfileEthPortProfileName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ConfigPutRequestEthPortProfileEthPortProfileName) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestEthPortProfileEthPortProfileName) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestEthPortProfileEthPortProfileName) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestEthPortProfileEthPortProfileName) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetTenantSliceManaged

`func (o *ConfigPutRequestEthPortProfileEthPortProfileName) GetTenantSliceManaged() bool`

GetTenantSliceManaged returns the TenantSliceManaged field if non-nil, zero value otherwise.

### GetTenantSliceManagedOk

`func (o *ConfigPutRequestEthPortProfileEthPortProfileName) GetTenantSliceManagedOk() (*bool, bool)`

GetTenantSliceManagedOk returns a tuple with the TenantSliceManaged field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantSliceManaged

`func (o *ConfigPutRequestEthPortProfileEthPortProfileName) SetTenantSliceManaged(v bool)`

SetTenantSliceManaged sets TenantSliceManaged field to given value.

### HasTenantSliceManaged

`func (o *ConfigPutRequestEthPortProfileEthPortProfileName) HasTenantSliceManaged() bool`

HasTenantSliceManaged returns a boolean if a field has been set.

### GetServices

`func (o *ConfigPutRequestEthPortProfileEthPortProfileName) GetServices() []ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner`

GetServices returns the Services field if non-nil, zero value otherwise.

### GetServicesOk

`func (o *ConfigPutRequestEthPortProfileEthPortProfileName) GetServicesOk() (*[]ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner, bool)`

GetServicesOk returns a tuple with the Services field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServices

`func (o *ConfigPutRequestEthPortProfileEthPortProfileName) SetServices(v []ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner)`

SetServices sets Services field to given value.

### HasServices

`func (o *ConfigPutRequestEthPortProfileEthPortProfileName) HasServices() bool`

HasServices returns a boolean if a field has been set.

### GetObjectProperties

`func (o *ConfigPutRequestEthPortProfileEthPortProfileName) GetObjectProperties() ConfigPutRequestEthPortProfileEthPortProfileNameObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ConfigPutRequestEthPortProfileEthPortProfileName) GetObjectPropertiesOk() (*ConfigPutRequestEthPortProfileEthPortProfileNameObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ConfigPutRequestEthPortProfileEthPortProfileName) SetObjectProperties(v ConfigPutRequestEthPortProfileEthPortProfileNameObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ConfigPutRequestEthPortProfileEthPortProfileName) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


