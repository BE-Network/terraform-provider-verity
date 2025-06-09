# ConfigPutRequestFeatureFlagFeatureFlagName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**Feature** | Pointer to [**[]ConfigPutRequestFeatureFlagFeatureFlagNameFeatureInner**](ConfigPutRequestFeatureFlagFeatureFlagNameFeatureInner.md) |  | [optional] 

## Methods

### NewConfigPutRequestFeatureFlagFeatureFlagName

`func NewConfigPutRequestFeatureFlagFeatureFlagName() *ConfigPutRequestFeatureFlagFeatureFlagName`

NewConfigPutRequestFeatureFlagFeatureFlagName instantiates a new ConfigPutRequestFeatureFlagFeatureFlagName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestFeatureFlagFeatureFlagNameWithDefaults

`func NewConfigPutRequestFeatureFlagFeatureFlagNameWithDefaults() *ConfigPutRequestFeatureFlagFeatureFlagName`

NewConfigPutRequestFeatureFlagFeatureFlagNameWithDefaults instantiates a new ConfigPutRequestFeatureFlagFeatureFlagName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestFeatureFlagFeatureFlagName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestFeatureFlagFeatureFlagName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestFeatureFlagFeatureFlagName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestFeatureFlagFeatureFlagName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ConfigPutRequestFeatureFlagFeatureFlagName) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestFeatureFlagFeatureFlagName) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestFeatureFlagFeatureFlagName) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestFeatureFlagFeatureFlagName) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetFeature

`func (o *ConfigPutRequestFeatureFlagFeatureFlagName) GetFeature() []ConfigPutRequestFeatureFlagFeatureFlagNameFeatureInner`

GetFeature returns the Feature field if non-nil, zero value otherwise.

### GetFeatureOk

`func (o *ConfigPutRequestFeatureFlagFeatureFlagName) GetFeatureOk() (*[]ConfigPutRequestFeatureFlagFeatureFlagNameFeatureInner, bool)`

GetFeatureOk returns a tuple with the Feature field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFeature

`func (o *ConfigPutRequestFeatureFlagFeatureFlagName) SetFeature(v []ConfigPutRequestFeatureFlagFeatureFlagNameFeatureInner)`

SetFeature sets Feature field to given value.

### HasFeature

`func (o *ConfigPutRequestFeatureFlagFeatureFlagName) HasFeature() bool`

HasFeature returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


