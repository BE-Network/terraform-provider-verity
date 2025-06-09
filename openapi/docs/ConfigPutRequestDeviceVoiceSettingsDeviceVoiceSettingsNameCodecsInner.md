# ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CodecNumName** | Pointer to **string** | Name of this Codec | [optional] [default to "G.711MuLaw"]
**CodecNumEnable** | Pointer to **bool** | Enable Codec | [optional] [default to true]
**CodecNumPacketizationPeriod** | Pointer to **string** | Packet period selection interval in milliseconds | [optional] [default to "20"]
**CodecNumSilenceSuppression** | Pointer to **bool** | Specifies whether silence suppression is on or off. Valid values are 0 &#x3D; off and 1 &#x3D; on | [optional] [default to false]
**Index** | Pointer to **int32** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 

## Methods

### NewConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInner

`func NewConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInner() *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInner`

NewConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInner instantiates a new ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInnerWithDefaults

`func NewConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInnerWithDefaults() *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInner`

NewConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInnerWithDefaults instantiates a new ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCodecNumName

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInner) GetCodecNumName() string`

GetCodecNumName returns the CodecNumName field if non-nil, zero value otherwise.

### GetCodecNumNameOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInner) GetCodecNumNameOk() (*string, bool)`

GetCodecNumNameOk returns a tuple with the CodecNumName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCodecNumName

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInner) SetCodecNumName(v string)`

SetCodecNumName sets CodecNumName field to given value.

### HasCodecNumName

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInner) HasCodecNumName() bool`

HasCodecNumName returns a boolean if a field has been set.

### GetCodecNumEnable

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInner) GetCodecNumEnable() bool`

GetCodecNumEnable returns the CodecNumEnable field if non-nil, zero value otherwise.

### GetCodecNumEnableOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInner) GetCodecNumEnableOk() (*bool, bool)`

GetCodecNumEnableOk returns a tuple with the CodecNumEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCodecNumEnable

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInner) SetCodecNumEnable(v bool)`

SetCodecNumEnable sets CodecNumEnable field to given value.

### HasCodecNumEnable

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInner) HasCodecNumEnable() bool`

HasCodecNumEnable returns a boolean if a field has been set.

### GetCodecNumPacketizationPeriod

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInner) GetCodecNumPacketizationPeriod() string`

GetCodecNumPacketizationPeriod returns the CodecNumPacketizationPeriod field if non-nil, zero value otherwise.

### GetCodecNumPacketizationPeriodOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInner) GetCodecNumPacketizationPeriodOk() (*string, bool)`

GetCodecNumPacketizationPeriodOk returns a tuple with the CodecNumPacketizationPeriod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCodecNumPacketizationPeriod

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInner) SetCodecNumPacketizationPeriod(v string)`

SetCodecNumPacketizationPeriod sets CodecNumPacketizationPeriod field to given value.

### HasCodecNumPacketizationPeriod

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInner) HasCodecNumPacketizationPeriod() bool`

HasCodecNumPacketizationPeriod returns a boolean if a field has been set.

### GetCodecNumSilenceSuppression

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInner) GetCodecNumSilenceSuppression() bool`

GetCodecNumSilenceSuppression returns the CodecNumSilenceSuppression field if non-nil, zero value otherwise.

### GetCodecNumSilenceSuppressionOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInner) GetCodecNumSilenceSuppressionOk() (*bool, bool)`

GetCodecNumSilenceSuppressionOk returns a tuple with the CodecNumSilenceSuppression field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCodecNumSilenceSuppression

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInner) SetCodecNumSilenceSuppression(v bool)`

SetCodecNumSilenceSuppression sets CodecNumSilenceSuppression field to given value.

### HasCodecNumSilenceSuppression

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInner) HasCodecNumSilenceSuppression() bool`

HasCodecNumSilenceSuppression returns a boolean if a field has been set.

### GetIndex

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInner) GetIndex() int32`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInner) GetIndexOk() (*int32, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInner) SetIndex(v int32)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


