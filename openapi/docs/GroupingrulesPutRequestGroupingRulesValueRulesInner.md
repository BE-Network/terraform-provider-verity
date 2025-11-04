# GroupingrulesPutRequestGroupingRulesValueRulesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Enable** | Pointer to **bool** | Enable | [optional] [default to false]
**RuleInvert** | Pointer to **bool** | Invert the rule | [optional] [default to false]
**RuleType** | Pointer to **string** | Which type of rule to apply | [optional] [default to ""]
**RuleValue** | Pointer to **string** | Value to compare | [optional] [default to ""]
**RuleValuePath** | Pointer to **string** | Object to compare | [optional] [default to ""]
**RuleValuePathRefType** | Pointer to **string** | Object type for rule_value_path field | [optional] 
**Index** | Pointer to **int32** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 

## Methods

### NewGroupingrulesPutRequestGroupingRulesValueRulesInner

`func NewGroupingrulesPutRequestGroupingRulesValueRulesInner() *GroupingrulesPutRequestGroupingRulesValueRulesInner`

NewGroupingrulesPutRequestGroupingRulesValueRulesInner instantiates a new GroupingrulesPutRequestGroupingRulesValueRulesInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGroupingrulesPutRequestGroupingRulesValueRulesInnerWithDefaults

`func NewGroupingrulesPutRequestGroupingRulesValueRulesInnerWithDefaults() *GroupingrulesPutRequestGroupingRulesValueRulesInner`

NewGroupingrulesPutRequestGroupingRulesValueRulesInnerWithDefaults instantiates a new GroupingrulesPutRequestGroupingRulesValueRulesInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEnable

`func (o *GroupingrulesPutRequestGroupingRulesValueRulesInner) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *GroupingrulesPutRequestGroupingRulesValueRulesInner) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *GroupingrulesPutRequestGroupingRulesValueRulesInner) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *GroupingrulesPutRequestGroupingRulesValueRulesInner) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetRuleInvert

`func (o *GroupingrulesPutRequestGroupingRulesValueRulesInner) GetRuleInvert() bool`

GetRuleInvert returns the RuleInvert field if non-nil, zero value otherwise.

### GetRuleInvertOk

`func (o *GroupingrulesPutRequestGroupingRulesValueRulesInner) GetRuleInvertOk() (*bool, bool)`

GetRuleInvertOk returns a tuple with the RuleInvert field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRuleInvert

`func (o *GroupingrulesPutRequestGroupingRulesValueRulesInner) SetRuleInvert(v bool)`

SetRuleInvert sets RuleInvert field to given value.

### HasRuleInvert

`func (o *GroupingrulesPutRequestGroupingRulesValueRulesInner) HasRuleInvert() bool`

HasRuleInvert returns a boolean if a field has been set.

### GetRuleType

`func (o *GroupingrulesPutRequestGroupingRulesValueRulesInner) GetRuleType() string`

GetRuleType returns the RuleType field if non-nil, zero value otherwise.

### GetRuleTypeOk

`func (o *GroupingrulesPutRequestGroupingRulesValueRulesInner) GetRuleTypeOk() (*string, bool)`

GetRuleTypeOk returns a tuple with the RuleType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRuleType

`func (o *GroupingrulesPutRequestGroupingRulesValueRulesInner) SetRuleType(v string)`

SetRuleType sets RuleType field to given value.

### HasRuleType

`func (o *GroupingrulesPutRequestGroupingRulesValueRulesInner) HasRuleType() bool`

HasRuleType returns a boolean if a field has been set.

### GetRuleValue

`func (o *GroupingrulesPutRequestGroupingRulesValueRulesInner) GetRuleValue() string`

GetRuleValue returns the RuleValue field if non-nil, zero value otherwise.

### GetRuleValueOk

`func (o *GroupingrulesPutRequestGroupingRulesValueRulesInner) GetRuleValueOk() (*string, bool)`

GetRuleValueOk returns a tuple with the RuleValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRuleValue

`func (o *GroupingrulesPutRequestGroupingRulesValueRulesInner) SetRuleValue(v string)`

SetRuleValue sets RuleValue field to given value.

### HasRuleValue

`func (o *GroupingrulesPutRequestGroupingRulesValueRulesInner) HasRuleValue() bool`

HasRuleValue returns a boolean if a field has been set.

### GetRuleValuePath

`func (o *GroupingrulesPutRequestGroupingRulesValueRulesInner) GetRuleValuePath() string`

GetRuleValuePath returns the RuleValuePath field if non-nil, zero value otherwise.

### GetRuleValuePathOk

`func (o *GroupingrulesPutRequestGroupingRulesValueRulesInner) GetRuleValuePathOk() (*string, bool)`

GetRuleValuePathOk returns a tuple with the RuleValuePath field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRuleValuePath

`func (o *GroupingrulesPutRequestGroupingRulesValueRulesInner) SetRuleValuePath(v string)`

SetRuleValuePath sets RuleValuePath field to given value.

### HasRuleValuePath

`func (o *GroupingrulesPutRequestGroupingRulesValueRulesInner) HasRuleValuePath() bool`

HasRuleValuePath returns a boolean if a field has been set.

### GetRuleValuePathRefType

`func (o *GroupingrulesPutRequestGroupingRulesValueRulesInner) GetRuleValuePathRefType() string`

GetRuleValuePathRefType returns the RuleValuePathRefType field if non-nil, zero value otherwise.

### GetRuleValuePathRefTypeOk

`func (o *GroupingrulesPutRequestGroupingRulesValueRulesInner) GetRuleValuePathRefTypeOk() (*string, bool)`

GetRuleValuePathRefTypeOk returns a tuple with the RuleValuePathRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRuleValuePathRefType

`func (o *GroupingrulesPutRequestGroupingRulesValueRulesInner) SetRuleValuePathRefType(v string)`

SetRuleValuePathRefType sets RuleValuePathRefType field to given value.

### HasRuleValuePathRefType

`func (o *GroupingrulesPutRequestGroupingRulesValueRulesInner) HasRuleValuePathRefType() bool`

HasRuleValuePathRefType returns a boolean if a field has been set.

### GetIndex

`func (o *GroupingrulesPutRequestGroupingRulesValueRulesInner) GetIndex() int32`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *GroupingrulesPutRequestGroupingRulesValueRulesInner) GetIndexOk() (*int32, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *GroupingrulesPutRequestGroupingRulesValueRulesInner) SetIndex(v int32)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *GroupingrulesPutRequestGroupingRulesValueRulesInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


