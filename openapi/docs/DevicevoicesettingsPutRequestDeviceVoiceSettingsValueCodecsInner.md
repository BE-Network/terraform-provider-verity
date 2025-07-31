# DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CodecNumName** | Pointer to **string** | Name of this Codec | [optional] [default to "G.711MuLaw"]
**CodecNumEnable** | Pointer to **bool** | Enable Codec | [optional] [default to true]
**CodecNumPacketizationPeriod** | Pointer to **string** | Packet period selection interval in milliseconds | [optional] [default to "20"]
**CodecNumSilenceSuppression** | Pointer to **bool** | Specifies whether silence suppression is on or off. Valid values are 0 &#x3D; off and 1 &#x3D; on | [optional] [default to false]
**Index** | Pointer to **int32** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 

## Methods

### NewDevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner

`func NewDevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner() *DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner`

NewDevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner instantiates a new DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInnerWithDefaults

`func NewDevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInnerWithDefaults() *DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner`

NewDevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInnerWithDefaults instantiates a new DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCodecNumName

`func (o *DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner) GetCodecNumName() string`

GetCodecNumName returns the CodecNumName field if non-nil, zero value otherwise.

### GetCodecNumNameOk

`func (o *DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner) GetCodecNumNameOk() (*string, bool)`

GetCodecNumNameOk returns a tuple with the CodecNumName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCodecNumName

`func (o *DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner) SetCodecNumName(v string)`

SetCodecNumName sets CodecNumName field to given value.

### HasCodecNumName

`func (o *DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner) HasCodecNumName() bool`

HasCodecNumName returns a boolean if a field has been set.

### GetCodecNumEnable

`func (o *DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner) GetCodecNumEnable() bool`

GetCodecNumEnable returns the CodecNumEnable field if non-nil, zero value otherwise.

### GetCodecNumEnableOk

`func (o *DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner) GetCodecNumEnableOk() (*bool, bool)`

GetCodecNumEnableOk returns a tuple with the CodecNumEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCodecNumEnable

`func (o *DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner) SetCodecNumEnable(v bool)`

SetCodecNumEnable sets CodecNumEnable field to given value.

### HasCodecNumEnable

`func (o *DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner) HasCodecNumEnable() bool`

HasCodecNumEnable returns a boolean if a field has been set.

### GetCodecNumPacketizationPeriod

`func (o *DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner) GetCodecNumPacketizationPeriod() string`

GetCodecNumPacketizationPeriod returns the CodecNumPacketizationPeriod field if non-nil, zero value otherwise.

### GetCodecNumPacketizationPeriodOk

`func (o *DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner) GetCodecNumPacketizationPeriodOk() (*string, bool)`

GetCodecNumPacketizationPeriodOk returns a tuple with the CodecNumPacketizationPeriod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCodecNumPacketizationPeriod

`func (o *DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner) SetCodecNumPacketizationPeriod(v string)`

SetCodecNumPacketizationPeriod sets CodecNumPacketizationPeriod field to given value.

### HasCodecNumPacketizationPeriod

`func (o *DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner) HasCodecNumPacketizationPeriod() bool`

HasCodecNumPacketizationPeriod returns a boolean if a field has been set.

### GetCodecNumSilenceSuppression

`func (o *DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner) GetCodecNumSilenceSuppression() bool`

GetCodecNumSilenceSuppression returns the CodecNumSilenceSuppression field if non-nil, zero value otherwise.

### GetCodecNumSilenceSuppressionOk

`func (o *DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner) GetCodecNumSilenceSuppressionOk() (*bool, bool)`

GetCodecNumSilenceSuppressionOk returns a tuple with the CodecNumSilenceSuppression field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCodecNumSilenceSuppression

`func (o *DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner) SetCodecNumSilenceSuppression(v bool)`

SetCodecNumSilenceSuppression sets CodecNumSilenceSuppression field to given value.

### HasCodecNumSilenceSuppression

`func (o *DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner) HasCodecNumSilenceSuppression() bool`

HasCodecNumSilenceSuppression returns a boolean if a field has been set.

### GetIndex

`func (o *DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner) GetIndex() int32`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner) GetIndexOk() (*int32, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner) SetIndex(v int32)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *DevicevoicesettingsPutRequestDeviceVoiceSettingsValueCodecsInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


