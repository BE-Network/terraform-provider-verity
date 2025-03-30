# ConfigPutRequestStaticConnectionsStaticConnectionsName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to true]
**Connections** | Pointer to [**[]ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner**](ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner.md) |  | [optional] 
**ObjectProperties** | Pointer to **map[string]interface{}** |  | [optional] 

## Methods

### NewConfigPutRequestStaticConnectionsStaticConnectionsName

`func NewConfigPutRequestStaticConnectionsStaticConnectionsName() *ConfigPutRequestStaticConnectionsStaticConnectionsName`

NewConfigPutRequestStaticConnectionsStaticConnectionsName instantiates a new ConfigPutRequestStaticConnectionsStaticConnectionsName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestStaticConnectionsStaticConnectionsNameWithDefaults

`func NewConfigPutRequestStaticConnectionsStaticConnectionsNameWithDefaults() *ConfigPutRequestStaticConnectionsStaticConnectionsName`

NewConfigPutRequestStaticConnectionsStaticConnectionsNameWithDefaults instantiates a new ConfigPutRequestStaticConnectionsStaticConnectionsName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEnable

`func (o *ConfigPutRequestStaticConnectionsStaticConnectionsName) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestStaticConnectionsStaticConnectionsName) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestStaticConnectionsStaticConnectionsName) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestStaticConnectionsStaticConnectionsName) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetConnections

`func (o *ConfigPutRequestStaticConnectionsStaticConnectionsName) GetConnections() []ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner`

GetConnections returns the Connections field if non-nil, zero value otherwise.

### GetConnectionsOk

`func (o *ConfigPutRequestStaticConnectionsStaticConnectionsName) GetConnectionsOk() (*[]ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner, bool)`

GetConnectionsOk returns a tuple with the Connections field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnections

`func (o *ConfigPutRequestStaticConnectionsStaticConnectionsName) SetConnections(v []ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner)`

SetConnections sets Connections field to given value.

### HasConnections

`func (o *ConfigPutRequestStaticConnectionsStaticConnectionsName) HasConnections() bool`

HasConnections returns a boolean if a field has been set.

### GetObjectProperties

`func (o *ConfigPutRequestStaticConnectionsStaticConnectionsName) GetObjectProperties() map[string]interface{}`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ConfigPutRequestStaticConnectionsStaticConnectionsName) GetObjectPropertiesOk() (*map[string]interface{}, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ConfigPutRequestStaticConnectionsStaticConnectionsName) SetObjectProperties(v map[string]interface{})`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ConfigPutRequestStaticConnectionsStaticConnectionsName) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


