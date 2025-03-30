# ConfigPutRequestIpPrefixListIpPrefixListNameListsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Enable** | Pointer to **bool** | Enable of this IP Prefix List | [optional] [default to false]
**PermitDeny** | Pointer to **string** | Action upon match of Community Strings. | [optional] [default to "permit"]
**Ipv4Prefix** | Pointer to **string** | IP address and subnet to match against  | [optional] [default to ""]
**GreaterThanEqualValue** | Pointer to **NullableInt32** | Match IP routes with a subnet mask greater than or equal to the value indicated  | [optional] 
**LessThanEqualValue** | Pointer to **NullableInt32** | Match IP routes with a subnet mask less than or equal to the value indicated | [optional] 
**Index** | Pointer to **int32** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 

## Methods

### NewConfigPutRequestIpPrefixListIpPrefixListNameListsInner

`func NewConfigPutRequestIpPrefixListIpPrefixListNameListsInner() *ConfigPutRequestIpPrefixListIpPrefixListNameListsInner`

NewConfigPutRequestIpPrefixListIpPrefixListNameListsInner instantiates a new ConfigPutRequestIpPrefixListIpPrefixListNameListsInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestIpPrefixListIpPrefixListNameListsInnerWithDefaults

`func NewConfigPutRequestIpPrefixListIpPrefixListNameListsInnerWithDefaults() *ConfigPutRequestIpPrefixListIpPrefixListNameListsInner`

NewConfigPutRequestIpPrefixListIpPrefixListNameListsInnerWithDefaults instantiates a new ConfigPutRequestIpPrefixListIpPrefixListNameListsInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEnable

`func (o *ConfigPutRequestIpPrefixListIpPrefixListNameListsInner) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestIpPrefixListIpPrefixListNameListsInner) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestIpPrefixListIpPrefixListNameListsInner) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestIpPrefixListIpPrefixListNameListsInner) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetPermitDeny

`func (o *ConfigPutRequestIpPrefixListIpPrefixListNameListsInner) GetPermitDeny() string`

GetPermitDeny returns the PermitDeny field if non-nil, zero value otherwise.

### GetPermitDenyOk

`func (o *ConfigPutRequestIpPrefixListIpPrefixListNameListsInner) GetPermitDenyOk() (*string, bool)`

GetPermitDenyOk returns a tuple with the PermitDeny field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPermitDeny

`func (o *ConfigPutRequestIpPrefixListIpPrefixListNameListsInner) SetPermitDeny(v string)`

SetPermitDeny sets PermitDeny field to given value.

### HasPermitDeny

`func (o *ConfigPutRequestIpPrefixListIpPrefixListNameListsInner) HasPermitDeny() bool`

HasPermitDeny returns a boolean if a field has been set.

### GetIpv4Prefix

`func (o *ConfigPutRequestIpPrefixListIpPrefixListNameListsInner) GetIpv4Prefix() string`

GetIpv4Prefix returns the Ipv4Prefix field if non-nil, zero value otherwise.

### GetIpv4PrefixOk

`func (o *ConfigPutRequestIpPrefixListIpPrefixListNameListsInner) GetIpv4PrefixOk() (*string, bool)`

GetIpv4PrefixOk returns a tuple with the Ipv4Prefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpv4Prefix

`func (o *ConfigPutRequestIpPrefixListIpPrefixListNameListsInner) SetIpv4Prefix(v string)`

SetIpv4Prefix sets Ipv4Prefix field to given value.

### HasIpv4Prefix

`func (o *ConfigPutRequestIpPrefixListIpPrefixListNameListsInner) HasIpv4Prefix() bool`

HasIpv4Prefix returns a boolean if a field has been set.

### GetGreaterThanEqualValue

`func (o *ConfigPutRequestIpPrefixListIpPrefixListNameListsInner) GetGreaterThanEqualValue() int32`

GetGreaterThanEqualValue returns the GreaterThanEqualValue field if non-nil, zero value otherwise.

### GetGreaterThanEqualValueOk

`func (o *ConfigPutRequestIpPrefixListIpPrefixListNameListsInner) GetGreaterThanEqualValueOk() (*int32, bool)`

GetGreaterThanEqualValueOk returns a tuple with the GreaterThanEqualValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGreaterThanEqualValue

`func (o *ConfigPutRequestIpPrefixListIpPrefixListNameListsInner) SetGreaterThanEqualValue(v int32)`

SetGreaterThanEqualValue sets GreaterThanEqualValue field to given value.

### HasGreaterThanEqualValue

`func (o *ConfigPutRequestIpPrefixListIpPrefixListNameListsInner) HasGreaterThanEqualValue() bool`

HasGreaterThanEqualValue returns a boolean if a field has been set.

### SetGreaterThanEqualValueNil

`func (o *ConfigPutRequestIpPrefixListIpPrefixListNameListsInner) SetGreaterThanEqualValueNil(b bool)`

 SetGreaterThanEqualValueNil sets the value for GreaterThanEqualValue to be an explicit nil

### UnsetGreaterThanEqualValue
`func (o *ConfigPutRequestIpPrefixListIpPrefixListNameListsInner) UnsetGreaterThanEqualValue()`

UnsetGreaterThanEqualValue ensures that no value is present for GreaterThanEqualValue, not even an explicit nil
### GetLessThanEqualValue

`func (o *ConfigPutRequestIpPrefixListIpPrefixListNameListsInner) GetLessThanEqualValue() int32`

GetLessThanEqualValue returns the LessThanEqualValue field if non-nil, zero value otherwise.

### GetLessThanEqualValueOk

`func (o *ConfigPutRequestIpPrefixListIpPrefixListNameListsInner) GetLessThanEqualValueOk() (*int32, bool)`

GetLessThanEqualValueOk returns a tuple with the LessThanEqualValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLessThanEqualValue

`func (o *ConfigPutRequestIpPrefixListIpPrefixListNameListsInner) SetLessThanEqualValue(v int32)`

SetLessThanEqualValue sets LessThanEqualValue field to given value.

### HasLessThanEqualValue

`func (o *ConfigPutRequestIpPrefixListIpPrefixListNameListsInner) HasLessThanEqualValue() bool`

HasLessThanEqualValue returns a boolean if a field has been set.

### SetLessThanEqualValueNil

`func (o *ConfigPutRequestIpPrefixListIpPrefixListNameListsInner) SetLessThanEqualValueNil(b bool)`

 SetLessThanEqualValueNil sets the value for LessThanEqualValue to be an explicit nil

### UnsetLessThanEqualValue
`func (o *ConfigPutRequestIpPrefixListIpPrefixListNameListsInner) UnsetLessThanEqualValue()`

UnsetLessThanEqualValue ensures that no value is present for LessThanEqualValue, not even an explicit nil
### GetIndex

`func (o *ConfigPutRequestIpPrefixListIpPrefixListNameListsInner) GetIndex() int32`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *ConfigPutRequestIpPrefixListIpPrefixListNameListsInner) GetIndexOk() (*int32, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *ConfigPutRequestIpPrefixListIpPrefixListNameListsInner) SetIndex(v int32)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *ConfigPutRequestIpPrefixListIpPrefixListNameListsInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


