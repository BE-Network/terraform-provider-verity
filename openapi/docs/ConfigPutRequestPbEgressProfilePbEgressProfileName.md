# ConfigPutRequestPbEgressProfilePbEgressProfileName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**Ipv4Permit** | Pointer to [**[]ConfigPutRequestPbEgressProfilePbEgressProfileNameIpv4PermitInner**](ConfigPutRequestPbEgressProfilePbEgressProfileNameIpv4PermitInner.md) |  | [optional] 
**Ipv4Deny** | Pointer to [**[]ConfigPutRequestPbEgressProfilePbEgressProfileNameIpv4PermitInner**](ConfigPutRequestPbEgressProfilePbEgressProfileNameIpv4PermitInner.md) |  | [optional] 
**Ipv6Permit** | Pointer to [**[]ConfigPutRequestPbEgressProfilePbEgressProfileNameIpv6PermitInner**](ConfigPutRequestPbEgressProfilePbEgressProfileNameIpv6PermitInner.md) |  | [optional] 
**Ipv6Deny** | Pointer to [**[]ConfigPutRequestPbEgressProfilePbEgressProfileNameIpv6PermitInner**](ConfigPutRequestPbEgressProfilePbEgressProfileNameIpv6PermitInner.md) |  | [optional] 

## Methods

### NewConfigPutRequestPbEgressProfilePbEgressProfileName

`func NewConfigPutRequestPbEgressProfilePbEgressProfileName() *ConfigPutRequestPbEgressProfilePbEgressProfileName`

NewConfigPutRequestPbEgressProfilePbEgressProfileName instantiates a new ConfigPutRequestPbEgressProfilePbEgressProfileName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestPbEgressProfilePbEgressProfileNameWithDefaults

`func NewConfigPutRequestPbEgressProfilePbEgressProfileNameWithDefaults() *ConfigPutRequestPbEgressProfilePbEgressProfileName`

NewConfigPutRequestPbEgressProfilePbEgressProfileNameWithDefaults instantiates a new ConfigPutRequestPbEgressProfilePbEgressProfileName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestPbEgressProfilePbEgressProfileName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestPbEgressProfilePbEgressProfileName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestPbEgressProfilePbEgressProfileName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestPbEgressProfilePbEgressProfileName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ConfigPutRequestPbEgressProfilePbEgressProfileName) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestPbEgressProfilePbEgressProfileName) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestPbEgressProfilePbEgressProfileName) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestPbEgressProfilePbEgressProfileName) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetIpv4Permit

`func (o *ConfigPutRequestPbEgressProfilePbEgressProfileName) GetIpv4Permit() []ConfigPutRequestPbEgressProfilePbEgressProfileNameIpv4PermitInner`

GetIpv4Permit returns the Ipv4Permit field if non-nil, zero value otherwise.

### GetIpv4PermitOk

`func (o *ConfigPutRequestPbEgressProfilePbEgressProfileName) GetIpv4PermitOk() (*[]ConfigPutRequestPbEgressProfilePbEgressProfileNameIpv4PermitInner, bool)`

GetIpv4PermitOk returns a tuple with the Ipv4Permit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpv4Permit

`func (o *ConfigPutRequestPbEgressProfilePbEgressProfileName) SetIpv4Permit(v []ConfigPutRequestPbEgressProfilePbEgressProfileNameIpv4PermitInner)`

SetIpv4Permit sets Ipv4Permit field to given value.

### HasIpv4Permit

`func (o *ConfigPutRequestPbEgressProfilePbEgressProfileName) HasIpv4Permit() bool`

HasIpv4Permit returns a boolean if a field has been set.

### GetIpv4Deny

`func (o *ConfigPutRequestPbEgressProfilePbEgressProfileName) GetIpv4Deny() []ConfigPutRequestPbEgressProfilePbEgressProfileNameIpv4PermitInner`

GetIpv4Deny returns the Ipv4Deny field if non-nil, zero value otherwise.

### GetIpv4DenyOk

`func (o *ConfigPutRequestPbEgressProfilePbEgressProfileName) GetIpv4DenyOk() (*[]ConfigPutRequestPbEgressProfilePbEgressProfileNameIpv4PermitInner, bool)`

GetIpv4DenyOk returns a tuple with the Ipv4Deny field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpv4Deny

`func (o *ConfigPutRequestPbEgressProfilePbEgressProfileName) SetIpv4Deny(v []ConfigPutRequestPbEgressProfilePbEgressProfileNameIpv4PermitInner)`

SetIpv4Deny sets Ipv4Deny field to given value.

### HasIpv4Deny

`func (o *ConfigPutRequestPbEgressProfilePbEgressProfileName) HasIpv4Deny() bool`

HasIpv4Deny returns a boolean if a field has been set.

### GetIpv6Permit

`func (o *ConfigPutRequestPbEgressProfilePbEgressProfileName) GetIpv6Permit() []ConfigPutRequestPbEgressProfilePbEgressProfileNameIpv6PermitInner`

GetIpv6Permit returns the Ipv6Permit field if non-nil, zero value otherwise.

### GetIpv6PermitOk

`func (o *ConfigPutRequestPbEgressProfilePbEgressProfileName) GetIpv6PermitOk() (*[]ConfigPutRequestPbEgressProfilePbEgressProfileNameIpv6PermitInner, bool)`

GetIpv6PermitOk returns a tuple with the Ipv6Permit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpv6Permit

`func (o *ConfigPutRequestPbEgressProfilePbEgressProfileName) SetIpv6Permit(v []ConfigPutRequestPbEgressProfilePbEgressProfileNameIpv6PermitInner)`

SetIpv6Permit sets Ipv6Permit field to given value.

### HasIpv6Permit

`func (o *ConfigPutRequestPbEgressProfilePbEgressProfileName) HasIpv6Permit() bool`

HasIpv6Permit returns a boolean if a field has been set.

### GetIpv6Deny

`func (o *ConfigPutRequestPbEgressProfilePbEgressProfileName) GetIpv6Deny() []ConfigPutRequestPbEgressProfilePbEgressProfileNameIpv6PermitInner`

GetIpv6Deny returns the Ipv6Deny field if non-nil, zero value otherwise.

### GetIpv6DenyOk

`func (o *ConfigPutRequestPbEgressProfilePbEgressProfileName) GetIpv6DenyOk() (*[]ConfigPutRequestPbEgressProfilePbEgressProfileNameIpv6PermitInner, bool)`

GetIpv6DenyOk returns a tuple with the Ipv6Deny field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpv6Deny

`func (o *ConfigPutRequestPbEgressProfilePbEgressProfileName) SetIpv6Deny(v []ConfigPutRequestPbEgressProfilePbEgressProfileNameIpv6PermitInner)`

SetIpv6Deny sets Ipv6Deny field to given value.

### HasIpv6Deny

`func (o *ConfigPutRequestPbEgressProfilePbEgressProfileName) HasIpv6Deny() bool`

HasIpv6Deny returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


