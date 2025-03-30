# ConfigPutRequestThresholdRulesThresholdRulesName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**Rules** | Pointer to [**[]ConfigPutRequestThresholdRulesThresholdRulesNameRulesInner**](ConfigPutRequestThresholdRulesThresholdRulesNameRulesInner.md) |  | [optional] 

## Methods

### NewConfigPutRequestThresholdRulesThresholdRulesName

`func NewConfigPutRequestThresholdRulesThresholdRulesName() *ConfigPutRequestThresholdRulesThresholdRulesName`

NewConfigPutRequestThresholdRulesThresholdRulesName instantiates a new ConfigPutRequestThresholdRulesThresholdRulesName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestThresholdRulesThresholdRulesNameWithDefaults

`func NewConfigPutRequestThresholdRulesThresholdRulesNameWithDefaults() *ConfigPutRequestThresholdRulesThresholdRulesName`

NewConfigPutRequestThresholdRulesThresholdRulesNameWithDefaults instantiates a new ConfigPutRequestThresholdRulesThresholdRulesName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestThresholdRulesThresholdRulesName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestThresholdRulesThresholdRulesName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestThresholdRulesThresholdRulesName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestThresholdRulesThresholdRulesName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ConfigPutRequestThresholdRulesThresholdRulesName) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestThresholdRulesThresholdRulesName) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestThresholdRulesThresholdRulesName) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestThresholdRulesThresholdRulesName) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetRules

`func (o *ConfigPutRequestThresholdRulesThresholdRulesName) GetRules() []ConfigPutRequestThresholdRulesThresholdRulesNameRulesInner`

GetRules returns the Rules field if non-nil, zero value otherwise.

### GetRulesOk

`func (o *ConfigPutRequestThresholdRulesThresholdRulesName) GetRulesOk() (*[]ConfigPutRequestThresholdRulesThresholdRulesNameRulesInner, bool)`

GetRulesOk returns a tuple with the Rules field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRules

`func (o *ConfigPutRequestThresholdRulesThresholdRulesName) SetRules(v []ConfigPutRequestThresholdRulesThresholdRulesNameRulesInner)`

SetRules sets Rules field to given value.

### HasRules

`func (o *ConfigPutRequestThresholdRulesThresholdRulesName) HasRules() bool`

HasRules returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


