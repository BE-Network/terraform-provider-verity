# ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Enable** | Pointer to **bool** | Enable of this IPv4 Prefix List | [optional] [default to false]
**PermitDeny** | Pointer to **string** | Action upon match of Community Strings. | [optional] [default to "permit"]
**Ipv4Prefix** | Pointer to **string** | IPv4 address and subnet to match against  | [optional] [default to ""]
**GreaterThanEqualValue** | Pointer to **NullableInt32** | Match IP routes with a subnet mask greater than or equal to the value indicated  | [optional] 
**LessThanEqualValue** | Pointer to **NullableInt32** | Match IP routes with a subnet mask less than or equal to the value indicated | [optional] 
**Index** | Pointer to **int32** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 

## Methods

### NewConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner

`func NewConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner() *ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner`

NewConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner instantiates a new ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInnerWithDefaults

`func NewConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInnerWithDefaults() *ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner`

NewConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInnerWithDefaults instantiates a new ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEnable

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetPermitDeny

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner) GetPermitDeny() string`

GetPermitDeny returns the PermitDeny field if non-nil, zero value otherwise.

### GetPermitDenyOk

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner) GetPermitDenyOk() (*string, bool)`

GetPermitDenyOk returns a tuple with the PermitDeny field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPermitDeny

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner) SetPermitDeny(v string)`

SetPermitDeny sets PermitDeny field to given value.

### HasPermitDeny

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner) HasPermitDeny() bool`

HasPermitDeny returns a boolean if a field has been set.

### GetIpv4Prefix

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner) GetIpv4Prefix() string`

GetIpv4Prefix returns the Ipv4Prefix field if non-nil, zero value otherwise.

### GetIpv4PrefixOk

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner) GetIpv4PrefixOk() (*string, bool)`

GetIpv4PrefixOk returns a tuple with the Ipv4Prefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpv4Prefix

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner) SetIpv4Prefix(v string)`

SetIpv4Prefix sets Ipv4Prefix field to given value.

### HasIpv4Prefix

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner) HasIpv4Prefix() bool`

HasIpv4Prefix returns a boolean if a field has been set.

### GetGreaterThanEqualValue

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner) GetGreaterThanEqualValue() int32`

GetGreaterThanEqualValue returns the GreaterThanEqualValue field if non-nil, zero value otherwise.

### GetGreaterThanEqualValueOk

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner) GetGreaterThanEqualValueOk() (*int32, bool)`

GetGreaterThanEqualValueOk returns a tuple with the GreaterThanEqualValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGreaterThanEqualValue

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner) SetGreaterThanEqualValue(v int32)`

SetGreaterThanEqualValue sets GreaterThanEqualValue field to given value.

### HasGreaterThanEqualValue

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner) HasGreaterThanEqualValue() bool`

HasGreaterThanEqualValue returns a boolean if a field has been set.

### SetGreaterThanEqualValueNil

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner) SetGreaterThanEqualValueNil(b bool)`

 SetGreaterThanEqualValueNil sets the value for GreaterThanEqualValue to be an explicit nil

### UnsetGreaterThanEqualValue
`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner) UnsetGreaterThanEqualValue()`

UnsetGreaterThanEqualValue ensures that no value is present for GreaterThanEqualValue, not even an explicit nil
### GetLessThanEqualValue

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner) GetLessThanEqualValue() int32`

GetLessThanEqualValue returns the LessThanEqualValue field if non-nil, zero value otherwise.

### GetLessThanEqualValueOk

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner) GetLessThanEqualValueOk() (*int32, bool)`

GetLessThanEqualValueOk returns a tuple with the LessThanEqualValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLessThanEqualValue

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner) SetLessThanEqualValue(v int32)`

SetLessThanEqualValue sets LessThanEqualValue field to given value.

### HasLessThanEqualValue

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner) HasLessThanEqualValue() bool`

HasLessThanEqualValue returns a boolean if a field has been set.

### SetLessThanEqualValueNil

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner) SetLessThanEqualValueNil(b bool)`

 SetLessThanEqualValueNil sets the value for LessThanEqualValue to be an explicit nil

### UnsetLessThanEqualValue
`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner) UnsetLessThanEqualValue()`

UnsetLessThanEqualValue ensures that no value is present for LessThanEqualValue, not even an explicit nil
### GetIndex

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner) GetIndex() int32`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner) GetIndexOk() (*int32, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner) SetIndex(v int32)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


