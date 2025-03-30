# ConfigPutRequestSfpBreakoutsSfpBreakoutsName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**Breakout** | Pointer to [**[]ConfigPutRequestSfpBreakoutsSfpBreakoutsNameBreakoutInner**](ConfigPutRequestSfpBreakoutsSfpBreakoutsNameBreakoutInner.md) |  | [optional] 
**ObjectProperties** | Pointer to **map[string]interface{}** |  | [optional] 

## Methods

### NewConfigPutRequestSfpBreakoutsSfpBreakoutsName

`func NewConfigPutRequestSfpBreakoutsSfpBreakoutsName() *ConfigPutRequestSfpBreakoutsSfpBreakoutsName`

NewConfigPutRequestSfpBreakoutsSfpBreakoutsName instantiates a new ConfigPutRequestSfpBreakoutsSfpBreakoutsName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestSfpBreakoutsSfpBreakoutsNameWithDefaults

`func NewConfigPutRequestSfpBreakoutsSfpBreakoutsNameWithDefaults() *ConfigPutRequestSfpBreakoutsSfpBreakoutsName`

NewConfigPutRequestSfpBreakoutsSfpBreakoutsNameWithDefaults instantiates a new ConfigPutRequestSfpBreakoutsSfpBreakoutsName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestSfpBreakoutsSfpBreakoutsName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestSfpBreakoutsSfpBreakoutsName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestSfpBreakoutsSfpBreakoutsName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestSfpBreakoutsSfpBreakoutsName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ConfigPutRequestSfpBreakoutsSfpBreakoutsName) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestSfpBreakoutsSfpBreakoutsName) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestSfpBreakoutsSfpBreakoutsName) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestSfpBreakoutsSfpBreakoutsName) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetBreakout

`func (o *ConfigPutRequestSfpBreakoutsSfpBreakoutsName) GetBreakout() []ConfigPutRequestSfpBreakoutsSfpBreakoutsNameBreakoutInner`

GetBreakout returns the Breakout field if non-nil, zero value otherwise.

### GetBreakoutOk

`func (o *ConfigPutRequestSfpBreakoutsSfpBreakoutsName) GetBreakoutOk() (*[]ConfigPutRequestSfpBreakoutsSfpBreakoutsNameBreakoutInner, bool)`

GetBreakoutOk returns a tuple with the Breakout field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBreakout

`func (o *ConfigPutRequestSfpBreakoutsSfpBreakoutsName) SetBreakout(v []ConfigPutRequestSfpBreakoutsSfpBreakoutsNameBreakoutInner)`

SetBreakout sets Breakout field to given value.

### HasBreakout

`func (o *ConfigPutRequestSfpBreakoutsSfpBreakoutsName) HasBreakout() bool`

HasBreakout returns a boolean if a field has been set.

### GetObjectProperties

`func (o *ConfigPutRequestSfpBreakoutsSfpBreakoutsName) GetObjectProperties() map[string]interface{}`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ConfigPutRequestSfpBreakoutsSfpBreakoutsName) GetObjectPropertiesOk() (*map[string]interface{}, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ConfigPutRequestSfpBreakoutsSfpBreakoutsName) SetObjectProperties(v map[string]interface{})`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ConfigPutRequestSfpBreakoutsSfpBreakoutsName) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


