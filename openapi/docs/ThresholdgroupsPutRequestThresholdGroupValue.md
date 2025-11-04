# ThresholdgroupsPutRequestThresholdGroupValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**Type** | Pointer to **string** | Type of elements to apply thresholds to | [optional] [default to "device"]
**Targets** | Pointer to [**[]ThresholdgroupsPutRequestThresholdGroupValueTargetsInner**](ThresholdgroupsPutRequestThresholdGroupValueTargetsInner.md) |  | [optional] 
**Thresholds** | Pointer to [**[]ThresholdgroupsPutRequestThresholdGroupValueThresholdsInner**](ThresholdgroupsPutRequestThresholdGroupValueThresholdsInner.md) |  | [optional] 

## Methods

### NewThresholdgroupsPutRequestThresholdGroupValue

`func NewThresholdgroupsPutRequestThresholdGroupValue() *ThresholdgroupsPutRequestThresholdGroupValue`

NewThresholdgroupsPutRequestThresholdGroupValue instantiates a new ThresholdgroupsPutRequestThresholdGroupValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewThresholdgroupsPutRequestThresholdGroupValueWithDefaults

`func NewThresholdgroupsPutRequestThresholdGroupValueWithDefaults() *ThresholdgroupsPutRequestThresholdGroupValue`

NewThresholdgroupsPutRequestThresholdGroupValueWithDefaults instantiates a new ThresholdgroupsPutRequestThresholdGroupValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ThresholdgroupsPutRequestThresholdGroupValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ThresholdgroupsPutRequestThresholdGroupValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ThresholdgroupsPutRequestThresholdGroupValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ThresholdgroupsPutRequestThresholdGroupValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ThresholdgroupsPutRequestThresholdGroupValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ThresholdgroupsPutRequestThresholdGroupValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ThresholdgroupsPutRequestThresholdGroupValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ThresholdgroupsPutRequestThresholdGroupValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetType

`func (o *ThresholdgroupsPutRequestThresholdGroupValue) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ThresholdgroupsPutRequestThresholdGroupValue) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ThresholdgroupsPutRequestThresholdGroupValue) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *ThresholdgroupsPutRequestThresholdGroupValue) HasType() bool`

HasType returns a boolean if a field has been set.

### GetTargets

`func (o *ThresholdgroupsPutRequestThresholdGroupValue) GetTargets() []ThresholdgroupsPutRequestThresholdGroupValueTargetsInner`

GetTargets returns the Targets field if non-nil, zero value otherwise.

### GetTargetsOk

`func (o *ThresholdgroupsPutRequestThresholdGroupValue) GetTargetsOk() (*[]ThresholdgroupsPutRequestThresholdGroupValueTargetsInner, bool)`

GetTargetsOk returns a tuple with the Targets field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTargets

`func (o *ThresholdgroupsPutRequestThresholdGroupValue) SetTargets(v []ThresholdgroupsPutRequestThresholdGroupValueTargetsInner)`

SetTargets sets Targets field to given value.

### HasTargets

`func (o *ThresholdgroupsPutRequestThresholdGroupValue) HasTargets() bool`

HasTargets returns a boolean if a field has been set.

### GetThresholds

`func (o *ThresholdgroupsPutRequestThresholdGroupValue) GetThresholds() []ThresholdgroupsPutRequestThresholdGroupValueThresholdsInner`

GetThresholds returns the Thresholds field if non-nil, zero value otherwise.

### GetThresholdsOk

`func (o *ThresholdgroupsPutRequestThresholdGroupValue) GetThresholdsOk() (*[]ThresholdgroupsPutRequestThresholdGroupValueThresholdsInner, bool)`

GetThresholdsOk returns a tuple with the Thresholds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetThresholds

`func (o *ThresholdgroupsPutRequestThresholdGroupValue) SetThresholds(v []ThresholdgroupsPutRequestThresholdGroupValueThresholdsInner)`

SetThresholds sets Thresholds field to given value.

### HasThresholds

`func (o *ThresholdgroupsPutRequestThresholdGroupValue) HasThresholds() bool`

HasThresholds returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


