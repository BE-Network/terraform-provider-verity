# GatewayprofilesPutRequestGatewayProfileValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. It&#39;s highly recommended to set this value to true so that validation on the object will be ran. | [optional] [default to false]
**TenantSliceManaged** | Pointer to **bool** | Profiles that Tenant Slice creates and manages | [optional] [default to false]
**ExternalGateways** | Pointer to [**[]GatewayprofilesPutRequestGatewayProfileValueExternalGatewaysInner**](GatewayprofilesPutRequestGatewayProfileValueExternalGatewaysInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**EthportsettingsPutRequestEthPortSettingsValueObjectProperties**](EthportsettingsPutRequestEthPortSettingsValueObjectProperties.md) |  | [optional] 

## Methods

### NewGatewayprofilesPutRequestGatewayProfileValue

`func NewGatewayprofilesPutRequestGatewayProfileValue() *GatewayprofilesPutRequestGatewayProfileValue`

NewGatewayprofilesPutRequestGatewayProfileValue instantiates a new GatewayprofilesPutRequestGatewayProfileValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGatewayprofilesPutRequestGatewayProfileValueWithDefaults

`func NewGatewayprofilesPutRequestGatewayProfileValueWithDefaults() *GatewayprofilesPutRequestGatewayProfileValue`

NewGatewayprofilesPutRequestGatewayProfileValueWithDefaults instantiates a new GatewayprofilesPutRequestGatewayProfileValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *GatewayprofilesPutRequestGatewayProfileValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *GatewayprofilesPutRequestGatewayProfileValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *GatewayprofilesPutRequestGatewayProfileValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *GatewayprofilesPutRequestGatewayProfileValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *GatewayprofilesPutRequestGatewayProfileValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *GatewayprofilesPutRequestGatewayProfileValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *GatewayprofilesPutRequestGatewayProfileValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *GatewayprofilesPutRequestGatewayProfileValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetTenantSliceManaged

`func (o *GatewayprofilesPutRequestGatewayProfileValue) GetTenantSliceManaged() bool`

GetTenantSliceManaged returns the TenantSliceManaged field if non-nil, zero value otherwise.

### GetTenantSliceManagedOk

`func (o *GatewayprofilesPutRequestGatewayProfileValue) GetTenantSliceManagedOk() (*bool, bool)`

GetTenantSliceManagedOk returns a tuple with the TenantSliceManaged field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantSliceManaged

`func (o *GatewayprofilesPutRequestGatewayProfileValue) SetTenantSliceManaged(v bool)`

SetTenantSliceManaged sets TenantSliceManaged field to given value.

### HasTenantSliceManaged

`func (o *GatewayprofilesPutRequestGatewayProfileValue) HasTenantSliceManaged() bool`

HasTenantSliceManaged returns a boolean if a field has been set.

### GetExternalGateways

`func (o *GatewayprofilesPutRequestGatewayProfileValue) GetExternalGateways() []GatewayprofilesPutRequestGatewayProfileValueExternalGatewaysInner`

GetExternalGateways returns the ExternalGateways field if non-nil, zero value otherwise.

### GetExternalGatewaysOk

`func (o *GatewayprofilesPutRequestGatewayProfileValue) GetExternalGatewaysOk() (*[]GatewayprofilesPutRequestGatewayProfileValueExternalGatewaysInner, bool)`

GetExternalGatewaysOk returns a tuple with the ExternalGateways field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalGateways

`func (o *GatewayprofilesPutRequestGatewayProfileValue) SetExternalGateways(v []GatewayprofilesPutRequestGatewayProfileValueExternalGatewaysInner)`

SetExternalGateways sets ExternalGateways field to given value.

### HasExternalGateways

`func (o *GatewayprofilesPutRequestGatewayProfileValue) HasExternalGateways() bool`

HasExternalGateways returns a boolean if a field has been set.

### GetObjectProperties

`func (o *GatewayprofilesPutRequestGatewayProfileValue) GetObjectProperties() EthportsettingsPutRequestEthPortSettingsValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *GatewayprofilesPutRequestGatewayProfileValue) GetObjectPropertiesOk() (*EthportsettingsPutRequestEthPortSettingsValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *GatewayprofilesPutRequestGatewayProfileValue) SetObjectProperties(v EthportsettingsPutRequestEthPortSettingsValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *GatewayprofilesPutRequestGatewayProfileValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


