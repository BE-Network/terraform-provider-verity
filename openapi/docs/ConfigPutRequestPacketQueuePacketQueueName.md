# ConfigPutRequestPacketQueuePacketQueueName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**Pbit** | Pointer to [**[]ConfigPutRequestPacketQueuePacketQueueNamePbitInner**](ConfigPutRequestPacketQueuePacketQueueNamePbitInner.md) |  | [optional] 
**Queue** | Pointer to [**[]ConfigPutRequestPacketQueuePacketQueueNameQueueInner**](ConfigPutRequestPacketQueuePacketQueueNameQueueInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**ConfigPutRequestPacketQueuePacketQueueNameObjectProperties**](ConfigPutRequestPacketQueuePacketQueueNameObjectProperties.md) |  | [optional] 

## Methods

### NewConfigPutRequestPacketQueuePacketQueueName

`func NewConfigPutRequestPacketQueuePacketQueueName() *ConfigPutRequestPacketQueuePacketQueueName`

NewConfigPutRequestPacketQueuePacketQueueName instantiates a new ConfigPutRequestPacketQueuePacketQueueName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestPacketQueuePacketQueueNameWithDefaults

`func NewConfigPutRequestPacketQueuePacketQueueNameWithDefaults() *ConfigPutRequestPacketQueuePacketQueueName`

NewConfigPutRequestPacketQueuePacketQueueNameWithDefaults instantiates a new ConfigPutRequestPacketQueuePacketQueueName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestPacketQueuePacketQueueName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestPacketQueuePacketQueueName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestPacketQueuePacketQueueName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestPacketQueuePacketQueueName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ConfigPutRequestPacketQueuePacketQueueName) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestPacketQueuePacketQueueName) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestPacketQueuePacketQueueName) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestPacketQueuePacketQueueName) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetPbit

`func (o *ConfigPutRequestPacketQueuePacketQueueName) GetPbit() []ConfigPutRequestPacketQueuePacketQueueNamePbitInner`

GetPbit returns the Pbit field if non-nil, zero value otherwise.

### GetPbitOk

`func (o *ConfigPutRequestPacketQueuePacketQueueName) GetPbitOk() (*[]ConfigPutRequestPacketQueuePacketQueueNamePbitInner, bool)`

GetPbitOk returns a tuple with the Pbit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPbit

`func (o *ConfigPutRequestPacketQueuePacketQueueName) SetPbit(v []ConfigPutRequestPacketQueuePacketQueueNamePbitInner)`

SetPbit sets Pbit field to given value.

### HasPbit

`func (o *ConfigPutRequestPacketQueuePacketQueueName) HasPbit() bool`

HasPbit returns a boolean if a field has been set.

### GetQueue

`func (o *ConfigPutRequestPacketQueuePacketQueueName) GetQueue() []ConfigPutRequestPacketQueuePacketQueueNameQueueInner`

GetQueue returns the Queue field if non-nil, zero value otherwise.

### GetQueueOk

`func (o *ConfigPutRequestPacketQueuePacketQueueName) GetQueueOk() (*[]ConfigPutRequestPacketQueuePacketQueueNameQueueInner, bool)`

GetQueueOk returns a tuple with the Queue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQueue

`func (o *ConfigPutRequestPacketQueuePacketQueueName) SetQueue(v []ConfigPutRequestPacketQueuePacketQueueNameQueueInner)`

SetQueue sets Queue field to given value.

### HasQueue

`func (o *ConfigPutRequestPacketQueuePacketQueueName) HasQueue() bool`

HasQueue returns a boolean if a field has been set.

### GetObjectProperties

`func (o *ConfigPutRequestPacketQueuePacketQueueName) GetObjectProperties() ConfigPutRequestPacketQueuePacketQueueNameObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ConfigPutRequestPacketQueuePacketQueueName) GetObjectPropertiesOk() (*ConfigPutRequestPacketQueuePacketQueueNameObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ConfigPutRequestPacketQueuePacketQueueName) SetObjectProperties(v ConfigPutRequestPacketQueuePacketQueueNameObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ConfigPutRequestPacketQueuePacketQueueName) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


