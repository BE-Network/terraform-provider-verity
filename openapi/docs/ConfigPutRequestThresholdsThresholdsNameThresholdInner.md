# ConfigPutRequestThresholdsThresholdsNameThresholdInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ThresholdNumEnabled** | Pointer to **bool** | Enable/Disable the threshold | [optional] [default to false]
**ThresholdNumStatType** | Pointer to **string** | Stat Type | [optional] [default to ""]
**ThresholdNumMethod** | Pointer to **string** | Method | [optional] [default to ""]
**ThresholdNumValue** | Pointer to **string** | Value | [optional] [default to ""]
**ThresholdNumAction** | Pointer to **string** | Action | [optional] [default to ""]
**ThresholdNumAlarmClearMethod** | Pointer to **string** | Alarm clear method | [optional] [default to ""]
**ThresholdNumAlarmClearValue** | Pointer to **string** | Alarm clear value | [optional] [default to ""]
**Index** | Pointer to **int32** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 

## Methods

### NewConfigPutRequestThresholdsThresholdsNameThresholdInner

`func NewConfigPutRequestThresholdsThresholdsNameThresholdInner() *ConfigPutRequestThresholdsThresholdsNameThresholdInner`

NewConfigPutRequestThresholdsThresholdsNameThresholdInner instantiates a new ConfigPutRequestThresholdsThresholdsNameThresholdInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestThresholdsThresholdsNameThresholdInnerWithDefaults

`func NewConfigPutRequestThresholdsThresholdsNameThresholdInnerWithDefaults() *ConfigPutRequestThresholdsThresholdsNameThresholdInner`

NewConfigPutRequestThresholdsThresholdsNameThresholdInnerWithDefaults instantiates a new ConfigPutRequestThresholdsThresholdsNameThresholdInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetThresholdNumEnabled

`func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) GetThresholdNumEnabled() bool`

GetThresholdNumEnabled returns the ThresholdNumEnabled field if non-nil, zero value otherwise.

### GetThresholdNumEnabledOk

`func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) GetThresholdNumEnabledOk() (*bool, bool)`

GetThresholdNumEnabledOk returns a tuple with the ThresholdNumEnabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetThresholdNumEnabled

`func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) SetThresholdNumEnabled(v bool)`

SetThresholdNumEnabled sets ThresholdNumEnabled field to given value.

### HasThresholdNumEnabled

`func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) HasThresholdNumEnabled() bool`

HasThresholdNumEnabled returns a boolean if a field has been set.

### GetThresholdNumStatType

`func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) GetThresholdNumStatType() string`

GetThresholdNumStatType returns the ThresholdNumStatType field if non-nil, zero value otherwise.

### GetThresholdNumStatTypeOk

`func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) GetThresholdNumStatTypeOk() (*string, bool)`

GetThresholdNumStatTypeOk returns a tuple with the ThresholdNumStatType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetThresholdNumStatType

`func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) SetThresholdNumStatType(v string)`

SetThresholdNumStatType sets ThresholdNumStatType field to given value.

### HasThresholdNumStatType

`func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) HasThresholdNumStatType() bool`

HasThresholdNumStatType returns a boolean if a field has been set.

### GetThresholdNumMethod

`func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) GetThresholdNumMethod() string`

GetThresholdNumMethod returns the ThresholdNumMethod field if non-nil, zero value otherwise.

### GetThresholdNumMethodOk

`func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) GetThresholdNumMethodOk() (*string, bool)`

GetThresholdNumMethodOk returns a tuple with the ThresholdNumMethod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetThresholdNumMethod

`func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) SetThresholdNumMethod(v string)`

SetThresholdNumMethod sets ThresholdNumMethod field to given value.

### HasThresholdNumMethod

`func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) HasThresholdNumMethod() bool`

HasThresholdNumMethod returns a boolean if a field has been set.

### GetThresholdNumValue

`func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) GetThresholdNumValue() string`

GetThresholdNumValue returns the ThresholdNumValue field if non-nil, zero value otherwise.

### GetThresholdNumValueOk

`func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) GetThresholdNumValueOk() (*string, bool)`

GetThresholdNumValueOk returns a tuple with the ThresholdNumValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetThresholdNumValue

`func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) SetThresholdNumValue(v string)`

SetThresholdNumValue sets ThresholdNumValue field to given value.

### HasThresholdNumValue

`func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) HasThresholdNumValue() bool`

HasThresholdNumValue returns a boolean if a field has been set.

### GetThresholdNumAction

`func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) GetThresholdNumAction() string`

GetThresholdNumAction returns the ThresholdNumAction field if non-nil, zero value otherwise.

### GetThresholdNumActionOk

`func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) GetThresholdNumActionOk() (*string, bool)`

GetThresholdNumActionOk returns a tuple with the ThresholdNumAction field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetThresholdNumAction

`func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) SetThresholdNumAction(v string)`

SetThresholdNumAction sets ThresholdNumAction field to given value.

### HasThresholdNumAction

`func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) HasThresholdNumAction() bool`

HasThresholdNumAction returns a boolean if a field has been set.

### GetThresholdNumAlarmClearMethod

`func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) GetThresholdNumAlarmClearMethod() string`

GetThresholdNumAlarmClearMethod returns the ThresholdNumAlarmClearMethod field if non-nil, zero value otherwise.

### GetThresholdNumAlarmClearMethodOk

`func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) GetThresholdNumAlarmClearMethodOk() (*string, bool)`

GetThresholdNumAlarmClearMethodOk returns a tuple with the ThresholdNumAlarmClearMethod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetThresholdNumAlarmClearMethod

`func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) SetThresholdNumAlarmClearMethod(v string)`

SetThresholdNumAlarmClearMethod sets ThresholdNumAlarmClearMethod field to given value.

### HasThresholdNumAlarmClearMethod

`func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) HasThresholdNumAlarmClearMethod() bool`

HasThresholdNumAlarmClearMethod returns a boolean if a field has been set.

### GetThresholdNumAlarmClearValue

`func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) GetThresholdNumAlarmClearValue() string`

GetThresholdNumAlarmClearValue returns the ThresholdNumAlarmClearValue field if non-nil, zero value otherwise.

### GetThresholdNumAlarmClearValueOk

`func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) GetThresholdNumAlarmClearValueOk() (*string, bool)`

GetThresholdNumAlarmClearValueOk returns a tuple with the ThresholdNumAlarmClearValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetThresholdNumAlarmClearValue

`func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) SetThresholdNumAlarmClearValue(v string)`

SetThresholdNumAlarmClearValue sets ThresholdNumAlarmClearValue field to given value.

### HasThresholdNumAlarmClearValue

`func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) HasThresholdNumAlarmClearValue() bool`

HasThresholdNumAlarmClearValue returns a boolean if a field has been set.

### GetIndex

`func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) GetIndex() int32`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) GetIndexOk() (*int32, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) SetIndex(v int32)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


