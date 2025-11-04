# ThresholdsPutRequestThresholdValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**Type** | Pointer to **string** | Type of elements threshold applies to | [optional] [default to "device"]
**Operation** | Pointer to **string** | How to combine rules | [optional] [default to "and"]
**Severity** | Pointer to **string** | Severity of the alarm when the threshold is met | [optional] [default to "notice"]
**For** | Pointer to **string** | Duration in minutes the threshold must be met before firing the alarm | [optional] [default to "5"]
**KeepFiringFor** | Pointer to **string** | Duration in minutes to keep firing the alarm after the threshold is no longer met | [optional] [default to "5"]
**Rules** | Pointer to [**[]ThresholdsPutRequestThresholdValueRulesInner**](ThresholdsPutRequestThresholdValueRulesInner.md) |  | [optional] 

## Methods

### NewThresholdsPutRequestThresholdValue

`func NewThresholdsPutRequestThresholdValue() *ThresholdsPutRequestThresholdValue`

NewThresholdsPutRequestThresholdValue instantiates a new ThresholdsPutRequestThresholdValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewThresholdsPutRequestThresholdValueWithDefaults

`func NewThresholdsPutRequestThresholdValueWithDefaults() *ThresholdsPutRequestThresholdValue`

NewThresholdsPutRequestThresholdValueWithDefaults instantiates a new ThresholdsPutRequestThresholdValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ThresholdsPutRequestThresholdValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ThresholdsPutRequestThresholdValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ThresholdsPutRequestThresholdValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ThresholdsPutRequestThresholdValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ThresholdsPutRequestThresholdValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ThresholdsPutRequestThresholdValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ThresholdsPutRequestThresholdValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ThresholdsPutRequestThresholdValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetType

`func (o *ThresholdsPutRequestThresholdValue) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ThresholdsPutRequestThresholdValue) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ThresholdsPutRequestThresholdValue) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *ThresholdsPutRequestThresholdValue) HasType() bool`

HasType returns a boolean if a field has been set.

### GetOperation

`func (o *ThresholdsPutRequestThresholdValue) GetOperation() string`

GetOperation returns the Operation field if non-nil, zero value otherwise.

### GetOperationOk

`func (o *ThresholdsPutRequestThresholdValue) GetOperationOk() (*string, bool)`

GetOperationOk returns a tuple with the Operation field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOperation

`func (o *ThresholdsPutRequestThresholdValue) SetOperation(v string)`

SetOperation sets Operation field to given value.

### HasOperation

`func (o *ThresholdsPutRequestThresholdValue) HasOperation() bool`

HasOperation returns a boolean if a field has been set.

### GetSeverity

`func (o *ThresholdsPutRequestThresholdValue) GetSeverity() string`

GetSeverity returns the Severity field if non-nil, zero value otherwise.

### GetSeverityOk

`func (o *ThresholdsPutRequestThresholdValue) GetSeverityOk() (*string, bool)`

GetSeverityOk returns a tuple with the Severity field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSeverity

`func (o *ThresholdsPutRequestThresholdValue) SetSeverity(v string)`

SetSeverity sets Severity field to given value.

### HasSeverity

`func (o *ThresholdsPutRequestThresholdValue) HasSeverity() bool`

HasSeverity returns a boolean if a field has been set.

### GetFor

`func (o *ThresholdsPutRequestThresholdValue) GetFor() string`

GetFor returns the For field if non-nil, zero value otherwise.

### GetForOk

`func (o *ThresholdsPutRequestThresholdValue) GetForOk() (*string, bool)`

GetForOk returns a tuple with the For field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFor

`func (o *ThresholdsPutRequestThresholdValue) SetFor(v string)`

SetFor sets For field to given value.

### HasFor

`func (o *ThresholdsPutRequestThresholdValue) HasFor() bool`

HasFor returns a boolean if a field has been set.

### GetKeepFiringFor

`func (o *ThresholdsPutRequestThresholdValue) GetKeepFiringFor() string`

GetKeepFiringFor returns the KeepFiringFor field if non-nil, zero value otherwise.

### GetKeepFiringForOk

`func (o *ThresholdsPutRequestThresholdValue) GetKeepFiringForOk() (*string, bool)`

GetKeepFiringForOk returns a tuple with the KeepFiringFor field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeepFiringFor

`func (o *ThresholdsPutRequestThresholdValue) SetKeepFiringFor(v string)`

SetKeepFiringFor sets KeepFiringFor field to given value.

### HasKeepFiringFor

`func (o *ThresholdsPutRequestThresholdValue) HasKeepFiringFor() bool`

HasKeepFiringFor returns a boolean if a field has been set.

### GetRules

`func (o *ThresholdsPutRequestThresholdValue) GetRules() []ThresholdsPutRequestThresholdValueRulesInner`

GetRules returns the Rules field if non-nil, zero value otherwise.

### GetRulesOk

`func (o *ThresholdsPutRequestThresholdValue) GetRulesOk() (*[]ThresholdsPutRequestThresholdValueRulesInner, bool)`

GetRulesOk returns a tuple with the Rules field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRules

`func (o *ThresholdsPutRequestThresholdValue) SetRules(v []ThresholdsPutRequestThresholdValueRulesInner)`

SetRules sets Rules field to given value.

### HasRules

`func (o *ThresholdsPutRequestThresholdValue) HasRules() bool`

HasRules returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


