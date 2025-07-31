# Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Enable** | Pointer to **bool** | Enable of this IPv6 Prefix List | [optional] [default to false]
**PermitDeny** | Pointer to **string** | Action upon match of Community Strings. | [optional] [default to "permit"]
**Ipv6Prefix** | Pointer to **string** | IPv6 address and subnet to match against  | [optional] [default to ""]
**GreaterThanEqualValue** | Pointer to **NullableInt32** | Match IP routes with a subnet mask greater than or equal to the value indicated  | [optional] 
**LessThanEqualValue** | Pointer to **NullableInt32** | Match IP routes with a subnet mask less than or equal to the value indicated | [optional] 

## Methods

### NewIpv6prefixlistsPutRequestIpv6PrefixListValueListsInner

`func NewIpv6prefixlistsPutRequestIpv6PrefixListValueListsInner() *Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner`

NewIpv6prefixlistsPutRequestIpv6PrefixListValueListsInner instantiates a new Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewIpv6prefixlistsPutRequestIpv6PrefixListValueListsInnerWithDefaults

`func NewIpv6prefixlistsPutRequestIpv6PrefixListValueListsInnerWithDefaults() *Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner`

NewIpv6prefixlistsPutRequestIpv6PrefixListValueListsInnerWithDefaults instantiates a new Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEnable

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetPermitDeny

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner) GetPermitDeny() string`

GetPermitDeny returns the PermitDeny field if non-nil, zero value otherwise.

### GetPermitDenyOk

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner) GetPermitDenyOk() (*string, bool)`

GetPermitDenyOk returns a tuple with the PermitDeny field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPermitDeny

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner) SetPermitDeny(v string)`

SetPermitDeny sets PermitDeny field to given value.

### HasPermitDeny

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner) HasPermitDeny() bool`

HasPermitDeny returns a boolean if a field has been set.

### GetIpv6Prefix

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner) GetIpv6Prefix() string`

GetIpv6Prefix returns the Ipv6Prefix field if non-nil, zero value otherwise.

### GetIpv6PrefixOk

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner) GetIpv6PrefixOk() (*string, bool)`

GetIpv6PrefixOk returns a tuple with the Ipv6Prefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpv6Prefix

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner) SetIpv6Prefix(v string)`

SetIpv6Prefix sets Ipv6Prefix field to given value.

### HasIpv6Prefix

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner) HasIpv6Prefix() bool`

HasIpv6Prefix returns a boolean if a field has been set.

### GetGreaterThanEqualValue

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner) GetGreaterThanEqualValue() int32`

GetGreaterThanEqualValue returns the GreaterThanEqualValue field if non-nil, zero value otherwise.

### GetGreaterThanEqualValueOk

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner) GetGreaterThanEqualValueOk() (*int32, bool)`

GetGreaterThanEqualValueOk returns a tuple with the GreaterThanEqualValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGreaterThanEqualValue

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner) SetGreaterThanEqualValue(v int32)`

SetGreaterThanEqualValue sets GreaterThanEqualValue field to given value.

### HasGreaterThanEqualValue

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner) HasGreaterThanEqualValue() bool`

HasGreaterThanEqualValue returns a boolean if a field has been set.

### SetGreaterThanEqualValueNil

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner) SetGreaterThanEqualValueNil(b bool)`

 SetGreaterThanEqualValueNil sets the value for GreaterThanEqualValue to be an explicit nil

### UnsetGreaterThanEqualValue
`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner) UnsetGreaterThanEqualValue()`

UnsetGreaterThanEqualValue ensures that no value is present for GreaterThanEqualValue, not even an explicit nil
### GetLessThanEqualValue

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner) GetLessThanEqualValue() int32`

GetLessThanEqualValue returns the LessThanEqualValue field if non-nil, zero value otherwise.

### GetLessThanEqualValueOk

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner) GetLessThanEqualValueOk() (*int32, bool)`

GetLessThanEqualValueOk returns a tuple with the LessThanEqualValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLessThanEqualValue

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner) SetLessThanEqualValue(v int32)`

SetLessThanEqualValue sets LessThanEqualValue field to given value.

### HasLessThanEqualValue

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner) HasLessThanEqualValue() bool`

HasLessThanEqualValue returns a boolean if a field has been set.

### SetLessThanEqualValueNil

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner) SetLessThanEqualValueNil(b bool)`

 SetLessThanEqualValueNil sets the value for LessThanEqualValue to be an explicit nil

### UnsetLessThanEqualValue
`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner) UnsetLessThanEqualValue()`

UnsetLessThanEqualValue ensures that no value is present for LessThanEqualValue, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


