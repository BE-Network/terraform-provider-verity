# ConfigPutRequestMacFilterMacFilterName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**Type** | Pointer to **string** | Black vs White MAC Filter | [optional] [default to "White"]
**Acl** | Pointer to [**[]ConfigPutRequestMacFilterMacFilterNameAclInner**](ConfigPutRequestMacFilterMacFilterNameAclInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties**](ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties.md) |  | [optional] 

## Methods

### NewConfigPutRequestMacFilterMacFilterName

`func NewConfigPutRequestMacFilterMacFilterName() *ConfigPutRequestMacFilterMacFilterName`

NewConfigPutRequestMacFilterMacFilterName instantiates a new ConfigPutRequestMacFilterMacFilterName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestMacFilterMacFilterNameWithDefaults

`func NewConfigPutRequestMacFilterMacFilterNameWithDefaults() *ConfigPutRequestMacFilterMacFilterName`

NewConfigPutRequestMacFilterMacFilterNameWithDefaults instantiates a new ConfigPutRequestMacFilterMacFilterName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestMacFilterMacFilterName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestMacFilterMacFilterName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestMacFilterMacFilterName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestMacFilterMacFilterName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ConfigPutRequestMacFilterMacFilterName) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestMacFilterMacFilterName) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestMacFilterMacFilterName) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestMacFilterMacFilterName) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetType

`func (o *ConfigPutRequestMacFilterMacFilterName) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ConfigPutRequestMacFilterMacFilterName) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ConfigPutRequestMacFilterMacFilterName) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *ConfigPutRequestMacFilterMacFilterName) HasType() bool`

HasType returns a boolean if a field has been set.

### GetAcl

`func (o *ConfigPutRequestMacFilterMacFilterName) GetAcl() []ConfigPutRequestMacFilterMacFilterNameAclInner`

GetAcl returns the Acl field if non-nil, zero value otherwise.

### GetAclOk

`func (o *ConfigPutRequestMacFilterMacFilterName) GetAclOk() (*[]ConfigPutRequestMacFilterMacFilterNameAclInner, bool)`

GetAclOk returns a tuple with the Acl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAcl

`func (o *ConfigPutRequestMacFilterMacFilterName) SetAcl(v []ConfigPutRequestMacFilterMacFilterNameAclInner)`

SetAcl sets Acl field to given value.

### HasAcl

`func (o *ConfigPutRequestMacFilterMacFilterName) HasAcl() bool`

HasAcl returns a boolean if a field has been set.

### GetObjectProperties

`func (o *ConfigPutRequestMacFilterMacFilterName) GetObjectProperties() ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ConfigPutRequestMacFilterMacFilterName) GetObjectPropertiesOk() (*ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ConfigPutRequestMacFilterMacFilterName) SetObjectProperties(v ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ConfigPutRequestMacFilterMacFilterName) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


