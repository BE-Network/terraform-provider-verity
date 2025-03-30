# ConfigPutRequestBizdConfigBizdConfigName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**GlobalReadOnly** | Pointer to **bool** | when true, switches system-wide do not apply the configuration writen by bizd | [optional] [default to true]
**GlobalImageUpdatesDisable** | Pointer to **bool** | when true, switches system-wide have their firmware upgrades disabled | [optional] [default to true]
**ObjectProperties** | Pointer to **map[string]interface{}** |  | [optional] 

## Methods

### NewConfigPutRequestBizdConfigBizdConfigName

`func NewConfigPutRequestBizdConfigBizdConfigName() *ConfigPutRequestBizdConfigBizdConfigName`

NewConfigPutRequestBizdConfigBizdConfigName instantiates a new ConfigPutRequestBizdConfigBizdConfigName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestBizdConfigBizdConfigNameWithDefaults

`func NewConfigPutRequestBizdConfigBizdConfigNameWithDefaults() *ConfigPutRequestBizdConfigBizdConfigName`

NewConfigPutRequestBizdConfigBizdConfigNameWithDefaults instantiates a new ConfigPutRequestBizdConfigBizdConfigName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestBizdConfigBizdConfigName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestBizdConfigBizdConfigName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestBizdConfigBizdConfigName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestBizdConfigBizdConfigName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ConfigPutRequestBizdConfigBizdConfigName) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestBizdConfigBizdConfigName) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestBizdConfigBizdConfigName) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestBizdConfigBizdConfigName) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetGlobalReadOnly

`func (o *ConfigPutRequestBizdConfigBizdConfigName) GetGlobalReadOnly() bool`

GetGlobalReadOnly returns the GlobalReadOnly field if non-nil, zero value otherwise.

### GetGlobalReadOnlyOk

`func (o *ConfigPutRequestBizdConfigBizdConfigName) GetGlobalReadOnlyOk() (*bool, bool)`

GetGlobalReadOnlyOk returns a tuple with the GlobalReadOnly field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGlobalReadOnly

`func (o *ConfigPutRequestBizdConfigBizdConfigName) SetGlobalReadOnly(v bool)`

SetGlobalReadOnly sets GlobalReadOnly field to given value.

### HasGlobalReadOnly

`func (o *ConfigPutRequestBizdConfigBizdConfigName) HasGlobalReadOnly() bool`

HasGlobalReadOnly returns a boolean if a field has been set.

### GetGlobalImageUpdatesDisable

`func (o *ConfigPutRequestBizdConfigBizdConfigName) GetGlobalImageUpdatesDisable() bool`

GetGlobalImageUpdatesDisable returns the GlobalImageUpdatesDisable field if non-nil, zero value otherwise.

### GetGlobalImageUpdatesDisableOk

`func (o *ConfigPutRequestBizdConfigBizdConfigName) GetGlobalImageUpdatesDisableOk() (*bool, bool)`

GetGlobalImageUpdatesDisableOk returns a tuple with the GlobalImageUpdatesDisable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGlobalImageUpdatesDisable

`func (o *ConfigPutRequestBizdConfigBizdConfigName) SetGlobalImageUpdatesDisable(v bool)`

SetGlobalImageUpdatesDisable sets GlobalImageUpdatesDisable field to given value.

### HasGlobalImageUpdatesDisable

`func (o *ConfigPutRequestBizdConfigBizdConfigName) HasGlobalImageUpdatesDisable() bool`

HasGlobalImageUpdatesDisable returns a boolean if a field has been set.

### GetObjectProperties

`func (o *ConfigPutRequestBizdConfigBizdConfigName) GetObjectProperties() map[string]interface{}`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ConfigPutRequestBizdConfigBizdConfigName) GetObjectPropertiesOk() (*map[string]interface{}, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ConfigPutRequestBizdConfigBizdConfigName) SetObjectProperties(v map[string]interface{})`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ConfigPutRequestBizdConfigBizdConfigName) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


