# SfpbreakoutsPatchRequestSfpBreakoutsValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**Breakout** | Pointer to [**[]SfpbreakoutsPatchRequestSfpBreakoutsValueBreakoutInner**](SfpbreakoutsPatchRequestSfpBreakoutsValueBreakoutInner.md) |  | [optional] 
**ObjectProperties** | Pointer to **map[string]interface{}** |  | [optional] 

## Methods

### NewSfpbreakoutsPatchRequestSfpBreakoutsValue

`func NewSfpbreakoutsPatchRequestSfpBreakoutsValue() *SfpbreakoutsPatchRequestSfpBreakoutsValue`

NewSfpbreakoutsPatchRequestSfpBreakoutsValue instantiates a new SfpbreakoutsPatchRequestSfpBreakoutsValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSfpbreakoutsPatchRequestSfpBreakoutsValueWithDefaults

`func NewSfpbreakoutsPatchRequestSfpBreakoutsValueWithDefaults() *SfpbreakoutsPatchRequestSfpBreakoutsValue`

NewSfpbreakoutsPatchRequestSfpBreakoutsValueWithDefaults instantiates a new SfpbreakoutsPatchRequestSfpBreakoutsValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *SfpbreakoutsPatchRequestSfpBreakoutsValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *SfpbreakoutsPatchRequestSfpBreakoutsValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *SfpbreakoutsPatchRequestSfpBreakoutsValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *SfpbreakoutsPatchRequestSfpBreakoutsValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *SfpbreakoutsPatchRequestSfpBreakoutsValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *SfpbreakoutsPatchRequestSfpBreakoutsValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *SfpbreakoutsPatchRequestSfpBreakoutsValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *SfpbreakoutsPatchRequestSfpBreakoutsValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetBreakout

`func (o *SfpbreakoutsPatchRequestSfpBreakoutsValue) GetBreakout() []SfpbreakoutsPatchRequestSfpBreakoutsValueBreakoutInner`

GetBreakout returns the Breakout field if non-nil, zero value otherwise.

### GetBreakoutOk

`func (o *SfpbreakoutsPatchRequestSfpBreakoutsValue) GetBreakoutOk() (*[]SfpbreakoutsPatchRequestSfpBreakoutsValueBreakoutInner, bool)`

GetBreakoutOk returns a tuple with the Breakout field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBreakout

`func (o *SfpbreakoutsPatchRequestSfpBreakoutsValue) SetBreakout(v []SfpbreakoutsPatchRequestSfpBreakoutsValueBreakoutInner)`

SetBreakout sets Breakout field to given value.

### HasBreakout

`func (o *SfpbreakoutsPatchRequestSfpBreakoutsValue) HasBreakout() bool`

HasBreakout returns a boolean if a field has been set.

### GetObjectProperties

`func (o *SfpbreakoutsPatchRequestSfpBreakoutsValue) GetObjectProperties() map[string]interface{}`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *SfpbreakoutsPatchRequestSfpBreakoutsValue) GetObjectPropertiesOk() (*map[string]interface{}, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *SfpbreakoutsPatchRequestSfpBreakoutsValue) SetObjectProperties(v map[string]interface{})`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *SfpbreakoutsPatchRequestSfpBreakoutsValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


