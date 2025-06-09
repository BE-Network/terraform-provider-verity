# ConfigPutRequestPacketQueuePacketQueueNameQueueInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BandwidthForQueue** | Pointer to **NullableInt32** | Percentage bandwidth allocated to Queue. 0 is no limit | [optional] [default to 0]
**SchedulerType** | Pointer to **string** | Scheduler Type for Queue | [optional] [default to "SP"]
**SchedulerWeight** | Pointer to **NullableInt32** | Weight associated with WRR or DWRR scheduler | [optional] [default to 0]
**Index** | Pointer to **int32** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 

## Methods

### NewConfigPutRequestPacketQueuePacketQueueNameQueueInner

`func NewConfigPutRequestPacketQueuePacketQueueNameQueueInner() *ConfigPutRequestPacketQueuePacketQueueNameQueueInner`

NewConfigPutRequestPacketQueuePacketQueueNameQueueInner instantiates a new ConfigPutRequestPacketQueuePacketQueueNameQueueInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestPacketQueuePacketQueueNameQueueInnerWithDefaults

`func NewConfigPutRequestPacketQueuePacketQueueNameQueueInnerWithDefaults() *ConfigPutRequestPacketQueuePacketQueueNameQueueInner`

NewConfigPutRequestPacketQueuePacketQueueNameQueueInnerWithDefaults instantiates a new ConfigPutRequestPacketQueuePacketQueueNameQueueInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBandwidthForQueue

`func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) GetBandwidthForQueue() int32`

GetBandwidthForQueue returns the BandwidthForQueue field if non-nil, zero value otherwise.

### GetBandwidthForQueueOk

`func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) GetBandwidthForQueueOk() (*int32, bool)`

GetBandwidthForQueueOk returns a tuple with the BandwidthForQueue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBandwidthForQueue

`func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) SetBandwidthForQueue(v int32)`

SetBandwidthForQueue sets BandwidthForQueue field to given value.

### HasBandwidthForQueue

`func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) HasBandwidthForQueue() bool`

HasBandwidthForQueue returns a boolean if a field has been set.

### SetBandwidthForQueueNil

`func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) SetBandwidthForQueueNil(b bool)`

 SetBandwidthForQueueNil sets the value for BandwidthForQueue to be an explicit nil

### UnsetBandwidthForQueue
`func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) UnsetBandwidthForQueue()`

UnsetBandwidthForQueue ensures that no value is present for BandwidthForQueue, not even an explicit nil
### GetSchedulerType

`func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) GetSchedulerType() string`

GetSchedulerType returns the SchedulerType field if non-nil, zero value otherwise.

### GetSchedulerTypeOk

`func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) GetSchedulerTypeOk() (*string, bool)`

GetSchedulerTypeOk returns a tuple with the SchedulerType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchedulerType

`func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) SetSchedulerType(v string)`

SetSchedulerType sets SchedulerType field to given value.

### HasSchedulerType

`func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) HasSchedulerType() bool`

HasSchedulerType returns a boolean if a field has been set.

### GetSchedulerWeight

`func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) GetSchedulerWeight() int32`

GetSchedulerWeight returns the SchedulerWeight field if non-nil, zero value otherwise.

### GetSchedulerWeightOk

`func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) GetSchedulerWeightOk() (*int32, bool)`

GetSchedulerWeightOk returns a tuple with the SchedulerWeight field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSchedulerWeight

`func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) SetSchedulerWeight(v int32)`

SetSchedulerWeight sets SchedulerWeight field to given value.

### HasSchedulerWeight

`func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) HasSchedulerWeight() bool`

HasSchedulerWeight returns a boolean if a field has been set.

### SetSchedulerWeightNil

`func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) SetSchedulerWeightNil(b bool)`

 SetSchedulerWeightNil sets the value for SchedulerWeight to be an explicit nil

### UnsetSchedulerWeight
`func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) UnsetSchedulerWeight()`

UnsetSchedulerWeight ensures that no value is present for SchedulerWeight, not even an explicit nil
### GetIndex

`func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) GetIndex() int32`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) GetIndexOk() (*int32, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) SetIndex(v int32)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


