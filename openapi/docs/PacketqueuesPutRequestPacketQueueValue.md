# PacketqueuesPutRequestPacketQueueValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**Pbit** | Pointer to [**[]PacketqueuesPutRequestPacketQueueValuePbitInner**](PacketqueuesPutRequestPacketQueueValuePbitInner.md) |  | [optional] 
**Queue** | Pointer to [**[]PacketqueuesPutRequestPacketQueueValueQueueInner**](PacketqueuesPutRequestPacketQueueValueQueueInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**PacketqueuesPutRequestPacketQueueValueObjectProperties**](PacketqueuesPutRequestPacketQueueValueObjectProperties.md) |  | [optional] 

## Methods

### NewPacketqueuesPutRequestPacketQueueValue

`func NewPacketqueuesPutRequestPacketQueueValue() *PacketqueuesPutRequestPacketQueueValue`

NewPacketqueuesPutRequestPacketQueueValue instantiates a new PacketqueuesPutRequestPacketQueueValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPacketqueuesPutRequestPacketQueueValueWithDefaults

`func NewPacketqueuesPutRequestPacketQueueValueWithDefaults() *PacketqueuesPutRequestPacketQueueValue`

NewPacketqueuesPutRequestPacketQueueValueWithDefaults instantiates a new PacketqueuesPutRequestPacketQueueValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *PacketqueuesPutRequestPacketQueueValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *PacketqueuesPutRequestPacketQueueValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *PacketqueuesPutRequestPacketQueueValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *PacketqueuesPutRequestPacketQueueValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *PacketqueuesPutRequestPacketQueueValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *PacketqueuesPutRequestPacketQueueValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *PacketqueuesPutRequestPacketQueueValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *PacketqueuesPutRequestPacketQueueValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetPbit

`func (o *PacketqueuesPutRequestPacketQueueValue) GetPbit() []PacketqueuesPutRequestPacketQueueValuePbitInner`

GetPbit returns the Pbit field if non-nil, zero value otherwise.

### GetPbitOk

`func (o *PacketqueuesPutRequestPacketQueueValue) GetPbitOk() (*[]PacketqueuesPutRequestPacketQueueValuePbitInner, bool)`

GetPbitOk returns a tuple with the Pbit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPbit

`func (o *PacketqueuesPutRequestPacketQueueValue) SetPbit(v []PacketqueuesPutRequestPacketQueueValuePbitInner)`

SetPbit sets Pbit field to given value.

### HasPbit

`func (o *PacketqueuesPutRequestPacketQueueValue) HasPbit() bool`

HasPbit returns a boolean if a field has been set.

### GetQueue

`func (o *PacketqueuesPutRequestPacketQueueValue) GetQueue() []PacketqueuesPutRequestPacketQueueValueQueueInner`

GetQueue returns the Queue field if non-nil, zero value otherwise.

### GetQueueOk

`func (o *PacketqueuesPutRequestPacketQueueValue) GetQueueOk() (*[]PacketqueuesPutRequestPacketQueueValueQueueInner, bool)`

GetQueueOk returns a tuple with the Queue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetQueue

`func (o *PacketqueuesPutRequestPacketQueueValue) SetQueue(v []PacketqueuesPutRequestPacketQueueValueQueueInner)`

SetQueue sets Queue field to given value.

### HasQueue

`func (o *PacketqueuesPutRequestPacketQueueValue) HasQueue() bool`

HasQueue returns a boolean if a field has been set.

### GetObjectProperties

`func (o *PacketqueuesPutRequestPacketQueueValue) GetObjectProperties() PacketqueuesPutRequestPacketQueueValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *PacketqueuesPutRequestPacketQueueValue) GetObjectPropertiesOk() (*PacketqueuesPutRequestPacketQueueValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *PacketqueuesPutRequestPacketQueueValue) SetObjectProperties(v PacketqueuesPutRequestPacketQueueValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *PacketqueuesPutRequestPacketQueueValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


