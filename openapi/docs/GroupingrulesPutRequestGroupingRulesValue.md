# GroupingrulesPutRequestGroupingRulesValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**Type** | Pointer to **string** | Type of elements to group | [optional] [default to "device"]
**Operation** | Pointer to **string** | How to combine rules | [optional] [default to "and"]
**Rules** | Pointer to [**[]GroupingrulesPutRequestGroupingRulesValueRulesInner**](GroupingrulesPutRequestGroupingRulesValueRulesInner.md) |  | [optional] 

## Methods

### NewGroupingrulesPutRequestGroupingRulesValue

`func NewGroupingrulesPutRequestGroupingRulesValue() *GroupingrulesPutRequestGroupingRulesValue`

NewGroupingrulesPutRequestGroupingRulesValue instantiates a new GroupingrulesPutRequestGroupingRulesValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGroupingrulesPutRequestGroupingRulesValueWithDefaults

`func NewGroupingrulesPutRequestGroupingRulesValueWithDefaults() *GroupingrulesPutRequestGroupingRulesValue`

NewGroupingrulesPutRequestGroupingRulesValueWithDefaults instantiates a new GroupingrulesPutRequestGroupingRulesValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *GroupingrulesPutRequestGroupingRulesValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *GroupingrulesPutRequestGroupingRulesValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *GroupingrulesPutRequestGroupingRulesValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *GroupingrulesPutRequestGroupingRulesValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *GroupingrulesPutRequestGroupingRulesValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *GroupingrulesPutRequestGroupingRulesValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *GroupingrulesPutRequestGroupingRulesValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *GroupingrulesPutRequestGroupingRulesValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetType

`func (o *GroupingrulesPutRequestGroupingRulesValue) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *GroupingrulesPutRequestGroupingRulesValue) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *GroupingrulesPutRequestGroupingRulesValue) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *GroupingrulesPutRequestGroupingRulesValue) HasType() bool`

HasType returns a boolean if a field has been set.

### GetOperation

`func (o *GroupingrulesPutRequestGroupingRulesValue) GetOperation() string`

GetOperation returns the Operation field if non-nil, zero value otherwise.

### GetOperationOk

`func (o *GroupingrulesPutRequestGroupingRulesValue) GetOperationOk() (*string, bool)`

GetOperationOk returns a tuple with the Operation field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOperation

`func (o *GroupingrulesPutRequestGroupingRulesValue) SetOperation(v string)`

SetOperation sets Operation field to given value.

### HasOperation

`func (o *GroupingrulesPutRequestGroupingRulesValue) HasOperation() bool`

HasOperation returns a boolean if a field has been set.

### GetRules

`func (o *GroupingrulesPutRequestGroupingRulesValue) GetRules() []GroupingrulesPutRequestGroupingRulesValueRulesInner`

GetRules returns the Rules field if non-nil, zero value otherwise.

### GetRulesOk

`func (o *GroupingrulesPutRequestGroupingRulesValue) GetRulesOk() (*[]GroupingrulesPutRequestGroupingRulesValueRulesInner, bool)`

GetRulesOk returns a tuple with the Rules field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRules

`func (o *GroupingrulesPutRequestGroupingRulesValue) SetRules(v []GroupingrulesPutRequestGroupingRulesValueRulesInner)`

SetRules sets Rules field to given value.

### HasRules

`func (o *GroupingrulesPutRequestGroupingRulesValue) HasRules() bool`

HasRules returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


