# EthportprofilesPutRequestEthPortProfileValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. It&#39;s highly recommended to set this value to true so that validation on the object will be ran. | [optional] [default to false]
**TenantSliceManaged** | Pointer to **bool** | Profiles that Tenant Slice creates and manages | [optional] [default to false]
**Services** | Pointer to [**[]EthportprofilesPutRequestEthPortProfileValueServicesInner**](EthportprofilesPutRequestEthPortProfileValueServicesInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**EthportprofilesPutRequestEthPortProfileValueObjectProperties**](EthportprofilesPutRequestEthPortProfileValueObjectProperties.md) |  | [optional] 

## Methods

### NewEthportprofilesPutRequestEthPortProfileValue

`func NewEthportprofilesPutRequestEthPortProfileValue() *EthportprofilesPutRequestEthPortProfileValue`

NewEthportprofilesPutRequestEthPortProfileValue instantiates a new EthportprofilesPutRequestEthPortProfileValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewEthportprofilesPutRequestEthPortProfileValueWithDefaults

`func NewEthportprofilesPutRequestEthPortProfileValueWithDefaults() *EthportprofilesPutRequestEthPortProfileValue`

NewEthportprofilesPutRequestEthPortProfileValueWithDefaults instantiates a new EthportprofilesPutRequestEthPortProfileValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *EthportprofilesPutRequestEthPortProfileValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *EthportprofilesPutRequestEthPortProfileValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *EthportprofilesPutRequestEthPortProfileValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *EthportprofilesPutRequestEthPortProfileValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetTenantSliceManaged

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetTenantSliceManaged() bool`

GetTenantSliceManaged returns the TenantSliceManaged field if non-nil, zero value otherwise.

### GetTenantSliceManagedOk

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetTenantSliceManagedOk() (*bool, bool)`

GetTenantSliceManagedOk returns a tuple with the TenantSliceManaged field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantSliceManaged

`func (o *EthportprofilesPutRequestEthPortProfileValue) SetTenantSliceManaged(v bool)`

SetTenantSliceManaged sets TenantSliceManaged field to given value.

### HasTenantSliceManaged

`func (o *EthportprofilesPutRequestEthPortProfileValue) HasTenantSliceManaged() bool`

HasTenantSliceManaged returns a boolean if a field has been set.

### GetServices

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetServices() []EthportprofilesPutRequestEthPortProfileValueServicesInner`

GetServices returns the Services field if non-nil, zero value otherwise.

### GetServicesOk

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetServicesOk() (*[]EthportprofilesPutRequestEthPortProfileValueServicesInner, bool)`

GetServicesOk returns a tuple with the Services field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServices

`func (o *EthportprofilesPutRequestEthPortProfileValue) SetServices(v []EthportprofilesPutRequestEthPortProfileValueServicesInner)`

SetServices sets Services field to given value.

### HasServices

`func (o *EthportprofilesPutRequestEthPortProfileValue) HasServices() bool`

HasServices returns a boolean if a field has been set.

### GetObjectProperties

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetObjectProperties() EthportprofilesPutRequestEthPortProfileValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetObjectPropertiesOk() (*EthportprofilesPutRequestEthPortProfileValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *EthportprofilesPutRequestEthPortProfileValue) SetObjectProperties(v EthportprofilesPutRequestEthPortProfileValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *EthportprofilesPutRequestEthPortProfileValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


