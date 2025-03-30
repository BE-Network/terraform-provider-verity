# ConfigPutRequestThresholdsThresholdsName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**InheritedThresholds** | Pointer to **string** | Choose a Thresholds | [optional] [default to ""]
**InheritedThresholdsRefType** | Pointer to **string** | Object type for inherited_thresholds field | [optional] 
**Threshold** | Pointer to [**[]ConfigPutRequestThresholdsThresholdsNameThresholdInner**](ConfigPutRequestThresholdsThresholdsNameThresholdInner.md) |  | [optional] 

## Methods

### NewConfigPutRequestThresholdsThresholdsName

`func NewConfigPutRequestThresholdsThresholdsName() *ConfigPutRequestThresholdsThresholdsName`

NewConfigPutRequestThresholdsThresholdsName instantiates a new ConfigPutRequestThresholdsThresholdsName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestThresholdsThresholdsNameWithDefaults

`func NewConfigPutRequestThresholdsThresholdsNameWithDefaults() *ConfigPutRequestThresholdsThresholdsName`

NewConfigPutRequestThresholdsThresholdsNameWithDefaults instantiates a new ConfigPutRequestThresholdsThresholdsName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestThresholdsThresholdsName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestThresholdsThresholdsName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestThresholdsThresholdsName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestThresholdsThresholdsName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ConfigPutRequestThresholdsThresholdsName) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestThresholdsThresholdsName) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestThresholdsThresholdsName) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestThresholdsThresholdsName) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetInheritedThresholds

`func (o *ConfigPutRequestThresholdsThresholdsName) GetInheritedThresholds() string`

GetInheritedThresholds returns the InheritedThresholds field if non-nil, zero value otherwise.

### GetInheritedThresholdsOk

`func (o *ConfigPutRequestThresholdsThresholdsName) GetInheritedThresholdsOk() (*string, bool)`

GetInheritedThresholdsOk returns a tuple with the InheritedThresholds field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInheritedThresholds

`func (o *ConfigPutRequestThresholdsThresholdsName) SetInheritedThresholds(v string)`

SetInheritedThresholds sets InheritedThresholds field to given value.

### HasInheritedThresholds

`func (o *ConfigPutRequestThresholdsThresholdsName) HasInheritedThresholds() bool`

HasInheritedThresholds returns a boolean if a field has been set.

### GetInheritedThresholdsRefType

`func (o *ConfigPutRequestThresholdsThresholdsName) GetInheritedThresholdsRefType() string`

GetInheritedThresholdsRefType returns the InheritedThresholdsRefType field if non-nil, zero value otherwise.

### GetInheritedThresholdsRefTypeOk

`func (o *ConfigPutRequestThresholdsThresholdsName) GetInheritedThresholdsRefTypeOk() (*string, bool)`

GetInheritedThresholdsRefTypeOk returns a tuple with the InheritedThresholdsRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInheritedThresholdsRefType

`func (o *ConfigPutRequestThresholdsThresholdsName) SetInheritedThresholdsRefType(v string)`

SetInheritedThresholdsRefType sets InheritedThresholdsRefType field to given value.

### HasInheritedThresholdsRefType

`func (o *ConfigPutRequestThresholdsThresholdsName) HasInheritedThresholdsRefType() bool`

HasInheritedThresholdsRefType returns a boolean if a field has been set.

### GetThreshold

`func (o *ConfigPutRequestThresholdsThresholdsName) GetThreshold() []ConfigPutRequestThresholdsThresholdsNameThresholdInner`

GetThreshold returns the Threshold field if non-nil, zero value otherwise.

### GetThresholdOk

`func (o *ConfigPutRequestThresholdsThresholdsName) GetThresholdOk() (*[]ConfigPutRequestThresholdsThresholdsNameThresholdInner, bool)`

GetThresholdOk returns a tuple with the Threshold field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetThreshold

`func (o *ConfigPutRequestThresholdsThresholdsName) SetThreshold(v []ConfigPutRequestThresholdsThresholdsNameThresholdInner)`

SetThreshold sets Threshold field to given value.

### HasThreshold

`func (o *ConfigPutRequestThresholdsThresholdsName) HasThreshold() bool`

HasThreshold returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


