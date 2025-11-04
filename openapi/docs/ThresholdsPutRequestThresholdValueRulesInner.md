# ThresholdsPutRequestThresholdValueRulesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Enable** | Pointer to **bool** | Enable | [optional] [default to false]
**Type** | Pointer to **string** | Use a metric or a nested threshold | [optional] [default to "metric"]
**Metric** | Pointer to **string** | Metric threshold is on | [optional] [default to ""]
**Operation** | Pointer to **string** | How to compare the metric to the value | [optional] [default to "=="]
**Value** | Pointer to **string** | Value to compare the metric to | [optional] [default to ""]
**Threshold** | Pointer to **string** | How to compare the metric to the value | [optional] [default to ""]
**ThresholdRefType** | Pointer to **string** | Object type for threshold field | [optional] 
**Index** | Pointer to **int32** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 

## Methods

### NewThresholdsPutRequestThresholdValueRulesInner

`func NewThresholdsPutRequestThresholdValueRulesInner() *ThresholdsPutRequestThresholdValueRulesInner`

NewThresholdsPutRequestThresholdValueRulesInner instantiates a new ThresholdsPutRequestThresholdValueRulesInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewThresholdsPutRequestThresholdValueRulesInnerWithDefaults

`func NewThresholdsPutRequestThresholdValueRulesInnerWithDefaults() *ThresholdsPutRequestThresholdValueRulesInner`

NewThresholdsPutRequestThresholdValueRulesInnerWithDefaults instantiates a new ThresholdsPutRequestThresholdValueRulesInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEnable

`func (o *ThresholdsPutRequestThresholdValueRulesInner) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ThresholdsPutRequestThresholdValueRulesInner) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ThresholdsPutRequestThresholdValueRulesInner) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ThresholdsPutRequestThresholdValueRulesInner) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetType

`func (o *ThresholdsPutRequestThresholdValueRulesInner) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ThresholdsPutRequestThresholdValueRulesInner) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ThresholdsPutRequestThresholdValueRulesInner) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *ThresholdsPutRequestThresholdValueRulesInner) HasType() bool`

HasType returns a boolean if a field has been set.

### GetMetric

`func (o *ThresholdsPutRequestThresholdValueRulesInner) GetMetric() string`

GetMetric returns the Metric field if non-nil, zero value otherwise.

### GetMetricOk

`func (o *ThresholdsPutRequestThresholdValueRulesInner) GetMetricOk() (*string, bool)`

GetMetricOk returns a tuple with the Metric field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMetric

`func (o *ThresholdsPutRequestThresholdValueRulesInner) SetMetric(v string)`

SetMetric sets Metric field to given value.

### HasMetric

`func (o *ThresholdsPutRequestThresholdValueRulesInner) HasMetric() bool`

HasMetric returns a boolean if a field has been set.

### GetOperation

`func (o *ThresholdsPutRequestThresholdValueRulesInner) GetOperation() string`

GetOperation returns the Operation field if non-nil, zero value otherwise.

### GetOperationOk

`func (o *ThresholdsPutRequestThresholdValueRulesInner) GetOperationOk() (*string, bool)`

GetOperationOk returns a tuple with the Operation field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOperation

`func (o *ThresholdsPutRequestThresholdValueRulesInner) SetOperation(v string)`

SetOperation sets Operation field to given value.

### HasOperation

`func (o *ThresholdsPutRequestThresholdValueRulesInner) HasOperation() bool`

HasOperation returns a boolean if a field has been set.

### GetValue

`func (o *ThresholdsPutRequestThresholdValueRulesInner) GetValue() string`

GetValue returns the Value field if non-nil, zero value otherwise.

### GetValueOk

`func (o *ThresholdsPutRequestThresholdValueRulesInner) GetValueOk() (*string, bool)`

GetValueOk returns a tuple with the Value field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetValue

`func (o *ThresholdsPutRequestThresholdValueRulesInner) SetValue(v string)`

SetValue sets Value field to given value.

### HasValue

`func (o *ThresholdsPutRequestThresholdValueRulesInner) HasValue() bool`

HasValue returns a boolean if a field has been set.

### GetThreshold

`func (o *ThresholdsPutRequestThresholdValueRulesInner) GetThreshold() string`

GetThreshold returns the Threshold field if non-nil, zero value otherwise.

### GetThresholdOk

`func (o *ThresholdsPutRequestThresholdValueRulesInner) GetThresholdOk() (*string, bool)`

GetThresholdOk returns a tuple with the Threshold field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetThreshold

`func (o *ThresholdsPutRequestThresholdValueRulesInner) SetThreshold(v string)`

SetThreshold sets Threshold field to given value.

### HasThreshold

`func (o *ThresholdsPutRequestThresholdValueRulesInner) HasThreshold() bool`

HasThreshold returns a boolean if a field has been set.

### GetThresholdRefType

`func (o *ThresholdsPutRequestThresholdValueRulesInner) GetThresholdRefType() string`

GetThresholdRefType returns the ThresholdRefType field if non-nil, zero value otherwise.

### GetThresholdRefTypeOk

`func (o *ThresholdsPutRequestThresholdValueRulesInner) GetThresholdRefTypeOk() (*string, bool)`

GetThresholdRefTypeOk returns a tuple with the ThresholdRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetThresholdRefType

`func (o *ThresholdsPutRequestThresholdValueRulesInner) SetThresholdRefType(v string)`

SetThresholdRefType sets ThresholdRefType field to given value.

### HasThresholdRefType

`func (o *ThresholdsPutRequestThresholdValueRulesInner) HasThresholdRefType() bool`

HasThresholdRefType returns a boolean if a field has been set.

### GetIndex

`func (o *ThresholdsPutRequestThresholdValueRulesInner) GetIndex() int32`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *ThresholdsPutRequestThresholdValueRulesInner) GetIndexOk() (*int32, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *ThresholdsPutRequestThresholdValueRulesInner) SetIndex(v int32)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *ThresholdsPutRequestThresholdValueRulesInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


