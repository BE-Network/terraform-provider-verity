# ConfigPutRequestVoicePortProfilesVoicePortProfilesName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**Protocol** | Pointer to **string** | Voice Protocol: MGCP or SIP | [optional] [default to "SIP"]
**DigitMap** | Pointer to **string** | Dial Plan | [optional] [default to "(T)"]
**CallThreeWayEnable** | Pointer to **bool** | Enable three way calling | [optional] [default to false]
**CallerIdEnable** | Pointer to **bool** | Caller ID | [optional] [default to false]
**CallerIdNameEnable** | Pointer to **bool** | Caller ID Name | [optional] [default to false]
**CallWaitingEnable** | Pointer to **bool** | Call Waiting | [optional] [default to false]
**CallForwardUnconditionalEnable** | Pointer to **bool** | Call Forward Unconditional | [optional] [default to false]
**CallForwardOnBusyEnable** | Pointer to **bool** | Call Forward On Busy | [optional] [default to false]
**CallForwardOnNoAnswerRingCount** | Pointer to **NullableInt32** | Call Forward on number of rings | [optional] [default to 4]
**CallTransferEnable** | Pointer to **bool** | Call Transfer | [optional] [default to false]
**AudioMwiEnable** | Pointer to **bool** | Audio Message Waiting Indicator | [optional] [default to false]
**AnonymousCallBlockEnable** | Pointer to **bool** | Block all anonymous calls | [optional] [default to false]
**DoNotDisturbEnable** | Pointer to **bool** | Do not disturb | [optional] [default to false]
**CidBlockingEnable** | Pointer to **bool** | CID Blocking | [optional] [default to false]
**CidNumPresentationStatus** | Pointer to **string** | CID Number Presentation | [optional] [default to "Public"]
**CidNamePresentationStatus** | Pointer to **string** | CID Name Presentation | [optional] [default to "Public"]
**CallWaitingCallerIdEnable** | Pointer to **bool** | Call Waiting Caller ID | [optional] [default to false]
**CallHoldEnable** | Pointer to **bool** | Call Hold | [optional] [default to false]
**VisualMwiEnable** | Pointer to **bool** | Visual Message Waiting Indicator | [optional] [default to false]
**MwiRefreshTimer** | Pointer to **NullableInt32** | Message Waiting Indicator Refresh | [optional] [default to 30]
**HotlineEnable** | Pointer to **bool** | Direct Connect | [optional] [default to false]
**DialToneFeatureDelay** | Pointer to **NullableInt32** | Dial Tone Feature Delay | [optional] [default to 4]
**IntercomEnable** | Pointer to **bool** | Intercom | [optional] [default to false]
**IntercomTransferEnable** | Pointer to **bool** | Intercom Transfer | [optional] [default to false]
**TransmitGain** | Pointer to **NullableInt32** | Transmit Gain in tenths of a dB.Example -30 would equal -3.0db | [optional] [default to -30]
**ReceiveGain** | Pointer to **NullableInt32** | Receive Gainin tenths of a dB. Example -30 would equal -3.0db | [optional] [default to -30]
**EchoCancellationEnable** | Pointer to **bool** | Echo Cancellation Enable | [optional] [default to true]
**JitterTarget** | Pointer to **NullableInt32** | The target value of the jitter buffer in milliseconds | [optional] [default to 40]
**JitterBufferMax** | Pointer to **NullableInt32** | The maximum depth of the jitter buffer in milliseconds | [optional] [default to 180]
**SignalingCode** | Pointer to **string** | Signaling Code | [optional] [default to "LoopStart"]
**ReleaseTimer** | Pointer to **NullableInt32** | Release timer defined in seconds. The default value of this attribute is 10 seconds | [optional] [default to 10]
**RohTimer** | Pointer to **NullableInt32** | Time in seconds for the receiver is off-hook before ROH tone is applied. The value 0 disables ROH timing. The default value is 15 seconds | [optional] [default to 15]
**ObjectProperties** | Pointer to [**ConfigPutRequestVoicePortProfilesVoicePortProfilesNameObjectProperties**](ConfigPutRequestVoicePortProfilesVoicePortProfilesNameObjectProperties.md) |  | [optional] 

## Methods

### NewConfigPutRequestVoicePortProfilesVoicePortProfilesName

`func NewConfigPutRequestVoicePortProfilesVoicePortProfilesName() *ConfigPutRequestVoicePortProfilesVoicePortProfilesName`

NewConfigPutRequestVoicePortProfilesVoicePortProfilesName instantiates a new ConfigPutRequestVoicePortProfilesVoicePortProfilesName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestVoicePortProfilesVoicePortProfilesNameWithDefaults

`func NewConfigPutRequestVoicePortProfilesVoicePortProfilesNameWithDefaults() *ConfigPutRequestVoicePortProfilesVoicePortProfilesName`

NewConfigPutRequestVoicePortProfilesVoicePortProfilesNameWithDefaults instantiates a new ConfigPutRequestVoicePortProfilesVoicePortProfilesName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetProtocol

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetProtocol() string`

GetProtocol returns the Protocol field if non-nil, zero value otherwise.

### GetProtocolOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetProtocolOk() (*string, bool)`

GetProtocolOk returns a tuple with the Protocol field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProtocol

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetProtocol(v string)`

SetProtocol sets Protocol field to given value.

### HasProtocol

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasProtocol() bool`

HasProtocol returns a boolean if a field has been set.

### GetDigitMap

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetDigitMap() string`

GetDigitMap returns the DigitMap field if non-nil, zero value otherwise.

### GetDigitMapOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetDigitMapOk() (*string, bool)`

GetDigitMapOk returns a tuple with the DigitMap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDigitMap

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetDigitMap(v string)`

SetDigitMap sets DigitMap field to given value.

### HasDigitMap

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasDigitMap() bool`

HasDigitMap returns a boolean if a field has been set.

### GetCallThreeWayEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetCallThreeWayEnable() bool`

GetCallThreeWayEnable returns the CallThreeWayEnable field if non-nil, zero value otherwise.

### GetCallThreeWayEnableOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetCallThreeWayEnableOk() (*bool, bool)`

GetCallThreeWayEnableOk returns a tuple with the CallThreeWayEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallThreeWayEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetCallThreeWayEnable(v bool)`

SetCallThreeWayEnable sets CallThreeWayEnable field to given value.

### HasCallThreeWayEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasCallThreeWayEnable() bool`

HasCallThreeWayEnable returns a boolean if a field has been set.

### GetCallerIdEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetCallerIdEnable() bool`

GetCallerIdEnable returns the CallerIdEnable field if non-nil, zero value otherwise.

### GetCallerIdEnableOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetCallerIdEnableOk() (*bool, bool)`

GetCallerIdEnableOk returns a tuple with the CallerIdEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallerIdEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetCallerIdEnable(v bool)`

SetCallerIdEnable sets CallerIdEnable field to given value.

### HasCallerIdEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasCallerIdEnable() bool`

HasCallerIdEnable returns a boolean if a field has been set.

### GetCallerIdNameEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetCallerIdNameEnable() bool`

GetCallerIdNameEnable returns the CallerIdNameEnable field if non-nil, zero value otherwise.

### GetCallerIdNameEnableOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetCallerIdNameEnableOk() (*bool, bool)`

GetCallerIdNameEnableOk returns a tuple with the CallerIdNameEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallerIdNameEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetCallerIdNameEnable(v bool)`

SetCallerIdNameEnable sets CallerIdNameEnable field to given value.

### HasCallerIdNameEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasCallerIdNameEnable() bool`

HasCallerIdNameEnable returns a boolean if a field has been set.

### GetCallWaitingEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetCallWaitingEnable() bool`

GetCallWaitingEnable returns the CallWaitingEnable field if non-nil, zero value otherwise.

### GetCallWaitingEnableOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetCallWaitingEnableOk() (*bool, bool)`

GetCallWaitingEnableOk returns a tuple with the CallWaitingEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallWaitingEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetCallWaitingEnable(v bool)`

SetCallWaitingEnable sets CallWaitingEnable field to given value.

### HasCallWaitingEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasCallWaitingEnable() bool`

HasCallWaitingEnable returns a boolean if a field has been set.

### GetCallForwardUnconditionalEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetCallForwardUnconditionalEnable() bool`

GetCallForwardUnconditionalEnable returns the CallForwardUnconditionalEnable field if non-nil, zero value otherwise.

### GetCallForwardUnconditionalEnableOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetCallForwardUnconditionalEnableOk() (*bool, bool)`

GetCallForwardUnconditionalEnableOk returns a tuple with the CallForwardUnconditionalEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallForwardUnconditionalEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetCallForwardUnconditionalEnable(v bool)`

SetCallForwardUnconditionalEnable sets CallForwardUnconditionalEnable field to given value.

### HasCallForwardUnconditionalEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasCallForwardUnconditionalEnable() bool`

HasCallForwardUnconditionalEnable returns a boolean if a field has been set.

### GetCallForwardOnBusyEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetCallForwardOnBusyEnable() bool`

GetCallForwardOnBusyEnable returns the CallForwardOnBusyEnable field if non-nil, zero value otherwise.

### GetCallForwardOnBusyEnableOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetCallForwardOnBusyEnableOk() (*bool, bool)`

GetCallForwardOnBusyEnableOk returns a tuple with the CallForwardOnBusyEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallForwardOnBusyEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetCallForwardOnBusyEnable(v bool)`

SetCallForwardOnBusyEnable sets CallForwardOnBusyEnable field to given value.

### HasCallForwardOnBusyEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasCallForwardOnBusyEnable() bool`

HasCallForwardOnBusyEnable returns a boolean if a field has been set.

### GetCallForwardOnNoAnswerRingCount

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetCallForwardOnNoAnswerRingCount() int32`

GetCallForwardOnNoAnswerRingCount returns the CallForwardOnNoAnswerRingCount field if non-nil, zero value otherwise.

### GetCallForwardOnNoAnswerRingCountOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetCallForwardOnNoAnswerRingCountOk() (*int32, bool)`

GetCallForwardOnNoAnswerRingCountOk returns a tuple with the CallForwardOnNoAnswerRingCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallForwardOnNoAnswerRingCount

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetCallForwardOnNoAnswerRingCount(v int32)`

SetCallForwardOnNoAnswerRingCount sets CallForwardOnNoAnswerRingCount field to given value.

### HasCallForwardOnNoAnswerRingCount

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasCallForwardOnNoAnswerRingCount() bool`

HasCallForwardOnNoAnswerRingCount returns a boolean if a field has been set.

### SetCallForwardOnNoAnswerRingCountNil

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetCallForwardOnNoAnswerRingCountNil(b bool)`

 SetCallForwardOnNoAnswerRingCountNil sets the value for CallForwardOnNoAnswerRingCount to be an explicit nil

### UnsetCallForwardOnNoAnswerRingCount
`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) UnsetCallForwardOnNoAnswerRingCount()`

UnsetCallForwardOnNoAnswerRingCount ensures that no value is present for CallForwardOnNoAnswerRingCount, not even an explicit nil
### GetCallTransferEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetCallTransferEnable() bool`

GetCallTransferEnable returns the CallTransferEnable field if non-nil, zero value otherwise.

### GetCallTransferEnableOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetCallTransferEnableOk() (*bool, bool)`

GetCallTransferEnableOk returns a tuple with the CallTransferEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallTransferEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetCallTransferEnable(v bool)`

SetCallTransferEnable sets CallTransferEnable field to given value.

### HasCallTransferEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasCallTransferEnable() bool`

HasCallTransferEnable returns a boolean if a field has been set.

### GetAudioMwiEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetAudioMwiEnable() bool`

GetAudioMwiEnable returns the AudioMwiEnable field if non-nil, zero value otherwise.

### GetAudioMwiEnableOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetAudioMwiEnableOk() (*bool, bool)`

GetAudioMwiEnableOk returns a tuple with the AudioMwiEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAudioMwiEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetAudioMwiEnable(v bool)`

SetAudioMwiEnable sets AudioMwiEnable field to given value.

### HasAudioMwiEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasAudioMwiEnable() bool`

HasAudioMwiEnable returns a boolean if a field has been set.

### GetAnonymousCallBlockEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetAnonymousCallBlockEnable() bool`

GetAnonymousCallBlockEnable returns the AnonymousCallBlockEnable field if non-nil, zero value otherwise.

### GetAnonymousCallBlockEnableOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetAnonymousCallBlockEnableOk() (*bool, bool)`

GetAnonymousCallBlockEnableOk returns a tuple with the AnonymousCallBlockEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnonymousCallBlockEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetAnonymousCallBlockEnable(v bool)`

SetAnonymousCallBlockEnable sets AnonymousCallBlockEnable field to given value.

### HasAnonymousCallBlockEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasAnonymousCallBlockEnable() bool`

HasAnonymousCallBlockEnable returns a boolean if a field has been set.

### GetDoNotDisturbEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetDoNotDisturbEnable() bool`

GetDoNotDisturbEnable returns the DoNotDisturbEnable field if non-nil, zero value otherwise.

### GetDoNotDisturbEnableOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetDoNotDisturbEnableOk() (*bool, bool)`

GetDoNotDisturbEnableOk returns a tuple with the DoNotDisturbEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDoNotDisturbEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetDoNotDisturbEnable(v bool)`

SetDoNotDisturbEnable sets DoNotDisturbEnable field to given value.

### HasDoNotDisturbEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasDoNotDisturbEnable() bool`

HasDoNotDisturbEnable returns a boolean if a field has been set.

### GetCidBlockingEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetCidBlockingEnable() bool`

GetCidBlockingEnable returns the CidBlockingEnable field if non-nil, zero value otherwise.

### GetCidBlockingEnableOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetCidBlockingEnableOk() (*bool, bool)`

GetCidBlockingEnableOk returns a tuple with the CidBlockingEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCidBlockingEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetCidBlockingEnable(v bool)`

SetCidBlockingEnable sets CidBlockingEnable field to given value.

### HasCidBlockingEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasCidBlockingEnable() bool`

HasCidBlockingEnable returns a boolean if a field has been set.

### GetCidNumPresentationStatus

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetCidNumPresentationStatus() string`

GetCidNumPresentationStatus returns the CidNumPresentationStatus field if non-nil, zero value otherwise.

### GetCidNumPresentationStatusOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetCidNumPresentationStatusOk() (*string, bool)`

GetCidNumPresentationStatusOk returns a tuple with the CidNumPresentationStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCidNumPresentationStatus

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetCidNumPresentationStatus(v string)`

SetCidNumPresentationStatus sets CidNumPresentationStatus field to given value.

### HasCidNumPresentationStatus

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasCidNumPresentationStatus() bool`

HasCidNumPresentationStatus returns a boolean if a field has been set.

### GetCidNamePresentationStatus

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetCidNamePresentationStatus() string`

GetCidNamePresentationStatus returns the CidNamePresentationStatus field if non-nil, zero value otherwise.

### GetCidNamePresentationStatusOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetCidNamePresentationStatusOk() (*string, bool)`

GetCidNamePresentationStatusOk returns a tuple with the CidNamePresentationStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCidNamePresentationStatus

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetCidNamePresentationStatus(v string)`

SetCidNamePresentationStatus sets CidNamePresentationStatus field to given value.

### HasCidNamePresentationStatus

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasCidNamePresentationStatus() bool`

HasCidNamePresentationStatus returns a boolean if a field has been set.

### GetCallWaitingCallerIdEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetCallWaitingCallerIdEnable() bool`

GetCallWaitingCallerIdEnable returns the CallWaitingCallerIdEnable field if non-nil, zero value otherwise.

### GetCallWaitingCallerIdEnableOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetCallWaitingCallerIdEnableOk() (*bool, bool)`

GetCallWaitingCallerIdEnableOk returns a tuple with the CallWaitingCallerIdEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallWaitingCallerIdEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetCallWaitingCallerIdEnable(v bool)`

SetCallWaitingCallerIdEnable sets CallWaitingCallerIdEnable field to given value.

### HasCallWaitingCallerIdEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasCallWaitingCallerIdEnable() bool`

HasCallWaitingCallerIdEnable returns a boolean if a field has been set.

### GetCallHoldEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetCallHoldEnable() bool`

GetCallHoldEnable returns the CallHoldEnable field if non-nil, zero value otherwise.

### GetCallHoldEnableOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetCallHoldEnableOk() (*bool, bool)`

GetCallHoldEnableOk returns a tuple with the CallHoldEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallHoldEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetCallHoldEnable(v bool)`

SetCallHoldEnable sets CallHoldEnable field to given value.

### HasCallHoldEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasCallHoldEnable() bool`

HasCallHoldEnable returns a boolean if a field has been set.

### GetVisualMwiEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetVisualMwiEnable() bool`

GetVisualMwiEnable returns the VisualMwiEnable field if non-nil, zero value otherwise.

### GetVisualMwiEnableOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetVisualMwiEnableOk() (*bool, bool)`

GetVisualMwiEnableOk returns a tuple with the VisualMwiEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVisualMwiEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetVisualMwiEnable(v bool)`

SetVisualMwiEnable sets VisualMwiEnable field to given value.

### HasVisualMwiEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasVisualMwiEnable() bool`

HasVisualMwiEnable returns a boolean if a field has been set.

### GetMwiRefreshTimer

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetMwiRefreshTimer() int32`

GetMwiRefreshTimer returns the MwiRefreshTimer field if non-nil, zero value otherwise.

### GetMwiRefreshTimerOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetMwiRefreshTimerOk() (*int32, bool)`

GetMwiRefreshTimerOk returns a tuple with the MwiRefreshTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMwiRefreshTimer

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetMwiRefreshTimer(v int32)`

SetMwiRefreshTimer sets MwiRefreshTimer field to given value.

### HasMwiRefreshTimer

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasMwiRefreshTimer() bool`

HasMwiRefreshTimer returns a boolean if a field has been set.

### SetMwiRefreshTimerNil

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetMwiRefreshTimerNil(b bool)`

 SetMwiRefreshTimerNil sets the value for MwiRefreshTimer to be an explicit nil

### UnsetMwiRefreshTimer
`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) UnsetMwiRefreshTimer()`

UnsetMwiRefreshTimer ensures that no value is present for MwiRefreshTimer, not even an explicit nil
### GetHotlineEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetHotlineEnable() bool`

GetHotlineEnable returns the HotlineEnable field if non-nil, zero value otherwise.

### GetHotlineEnableOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetHotlineEnableOk() (*bool, bool)`

GetHotlineEnableOk returns a tuple with the HotlineEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHotlineEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetHotlineEnable(v bool)`

SetHotlineEnable sets HotlineEnable field to given value.

### HasHotlineEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasHotlineEnable() bool`

HasHotlineEnable returns a boolean if a field has been set.

### GetDialToneFeatureDelay

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetDialToneFeatureDelay() int32`

GetDialToneFeatureDelay returns the DialToneFeatureDelay field if non-nil, zero value otherwise.

### GetDialToneFeatureDelayOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetDialToneFeatureDelayOk() (*int32, bool)`

GetDialToneFeatureDelayOk returns a tuple with the DialToneFeatureDelay field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDialToneFeatureDelay

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetDialToneFeatureDelay(v int32)`

SetDialToneFeatureDelay sets DialToneFeatureDelay field to given value.

### HasDialToneFeatureDelay

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasDialToneFeatureDelay() bool`

HasDialToneFeatureDelay returns a boolean if a field has been set.

### SetDialToneFeatureDelayNil

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetDialToneFeatureDelayNil(b bool)`

 SetDialToneFeatureDelayNil sets the value for DialToneFeatureDelay to be an explicit nil

### UnsetDialToneFeatureDelay
`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) UnsetDialToneFeatureDelay()`

UnsetDialToneFeatureDelay ensures that no value is present for DialToneFeatureDelay, not even an explicit nil
### GetIntercomEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetIntercomEnable() bool`

GetIntercomEnable returns the IntercomEnable field if non-nil, zero value otherwise.

### GetIntercomEnableOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetIntercomEnableOk() (*bool, bool)`

GetIntercomEnableOk returns a tuple with the IntercomEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIntercomEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetIntercomEnable(v bool)`

SetIntercomEnable sets IntercomEnable field to given value.

### HasIntercomEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasIntercomEnable() bool`

HasIntercomEnable returns a boolean if a field has been set.

### GetIntercomTransferEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetIntercomTransferEnable() bool`

GetIntercomTransferEnable returns the IntercomTransferEnable field if non-nil, zero value otherwise.

### GetIntercomTransferEnableOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetIntercomTransferEnableOk() (*bool, bool)`

GetIntercomTransferEnableOk returns a tuple with the IntercomTransferEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIntercomTransferEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetIntercomTransferEnable(v bool)`

SetIntercomTransferEnable sets IntercomTransferEnable field to given value.

### HasIntercomTransferEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasIntercomTransferEnable() bool`

HasIntercomTransferEnable returns a boolean if a field has been set.

### GetTransmitGain

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetTransmitGain() int32`

GetTransmitGain returns the TransmitGain field if non-nil, zero value otherwise.

### GetTransmitGainOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetTransmitGainOk() (*int32, bool)`

GetTransmitGainOk returns a tuple with the TransmitGain field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTransmitGain

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetTransmitGain(v int32)`

SetTransmitGain sets TransmitGain field to given value.

### HasTransmitGain

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasTransmitGain() bool`

HasTransmitGain returns a boolean if a field has been set.

### SetTransmitGainNil

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetTransmitGainNil(b bool)`

 SetTransmitGainNil sets the value for TransmitGain to be an explicit nil

### UnsetTransmitGain
`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) UnsetTransmitGain()`

UnsetTransmitGain ensures that no value is present for TransmitGain, not even an explicit nil
### GetReceiveGain

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetReceiveGain() int32`

GetReceiveGain returns the ReceiveGain field if non-nil, zero value otherwise.

### GetReceiveGainOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetReceiveGainOk() (*int32, bool)`

GetReceiveGainOk returns a tuple with the ReceiveGain field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReceiveGain

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetReceiveGain(v int32)`

SetReceiveGain sets ReceiveGain field to given value.

### HasReceiveGain

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasReceiveGain() bool`

HasReceiveGain returns a boolean if a field has been set.

### SetReceiveGainNil

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetReceiveGainNil(b bool)`

 SetReceiveGainNil sets the value for ReceiveGain to be an explicit nil

### UnsetReceiveGain
`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) UnsetReceiveGain()`

UnsetReceiveGain ensures that no value is present for ReceiveGain, not even an explicit nil
### GetEchoCancellationEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetEchoCancellationEnable() bool`

GetEchoCancellationEnable returns the EchoCancellationEnable field if non-nil, zero value otherwise.

### GetEchoCancellationEnableOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetEchoCancellationEnableOk() (*bool, bool)`

GetEchoCancellationEnableOk returns a tuple with the EchoCancellationEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEchoCancellationEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetEchoCancellationEnable(v bool)`

SetEchoCancellationEnable sets EchoCancellationEnable field to given value.

### HasEchoCancellationEnable

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasEchoCancellationEnable() bool`

HasEchoCancellationEnable returns a boolean if a field has been set.

### GetJitterTarget

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetJitterTarget() int32`

GetJitterTarget returns the JitterTarget field if non-nil, zero value otherwise.

### GetJitterTargetOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetJitterTargetOk() (*int32, bool)`

GetJitterTargetOk returns a tuple with the JitterTarget field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetJitterTarget

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetJitterTarget(v int32)`

SetJitterTarget sets JitterTarget field to given value.

### HasJitterTarget

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasJitterTarget() bool`

HasJitterTarget returns a boolean if a field has been set.

### SetJitterTargetNil

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetJitterTargetNil(b bool)`

 SetJitterTargetNil sets the value for JitterTarget to be an explicit nil

### UnsetJitterTarget
`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) UnsetJitterTarget()`

UnsetJitterTarget ensures that no value is present for JitterTarget, not even an explicit nil
### GetJitterBufferMax

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetJitterBufferMax() int32`

GetJitterBufferMax returns the JitterBufferMax field if non-nil, zero value otherwise.

### GetJitterBufferMaxOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetJitterBufferMaxOk() (*int32, bool)`

GetJitterBufferMaxOk returns a tuple with the JitterBufferMax field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetJitterBufferMax

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetJitterBufferMax(v int32)`

SetJitterBufferMax sets JitterBufferMax field to given value.

### HasJitterBufferMax

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasJitterBufferMax() bool`

HasJitterBufferMax returns a boolean if a field has been set.

### SetJitterBufferMaxNil

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetJitterBufferMaxNil(b bool)`

 SetJitterBufferMaxNil sets the value for JitterBufferMax to be an explicit nil

### UnsetJitterBufferMax
`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) UnsetJitterBufferMax()`

UnsetJitterBufferMax ensures that no value is present for JitterBufferMax, not even an explicit nil
### GetSignalingCode

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetSignalingCode() string`

GetSignalingCode returns the SignalingCode field if non-nil, zero value otherwise.

### GetSignalingCodeOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetSignalingCodeOk() (*string, bool)`

GetSignalingCodeOk returns a tuple with the SignalingCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSignalingCode

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetSignalingCode(v string)`

SetSignalingCode sets SignalingCode field to given value.

### HasSignalingCode

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasSignalingCode() bool`

HasSignalingCode returns a boolean if a field has been set.

### GetReleaseTimer

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetReleaseTimer() int32`

GetReleaseTimer returns the ReleaseTimer field if non-nil, zero value otherwise.

### GetReleaseTimerOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetReleaseTimerOk() (*int32, bool)`

GetReleaseTimerOk returns a tuple with the ReleaseTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReleaseTimer

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetReleaseTimer(v int32)`

SetReleaseTimer sets ReleaseTimer field to given value.

### HasReleaseTimer

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasReleaseTimer() bool`

HasReleaseTimer returns a boolean if a field has been set.

### SetReleaseTimerNil

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetReleaseTimerNil(b bool)`

 SetReleaseTimerNil sets the value for ReleaseTimer to be an explicit nil

### UnsetReleaseTimer
`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) UnsetReleaseTimer()`

UnsetReleaseTimer ensures that no value is present for ReleaseTimer, not even an explicit nil
### GetRohTimer

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetRohTimer() int32`

GetRohTimer returns the RohTimer field if non-nil, zero value otherwise.

### GetRohTimerOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetRohTimerOk() (*int32, bool)`

GetRohTimerOk returns a tuple with the RohTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRohTimer

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetRohTimer(v int32)`

SetRohTimer sets RohTimer field to given value.

### HasRohTimer

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasRohTimer() bool`

HasRohTimer returns a boolean if a field has been set.

### SetRohTimerNil

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetRohTimerNil(b bool)`

 SetRohTimerNil sets the value for RohTimer to be an explicit nil

### UnsetRohTimer
`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) UnsetRohTimer()`

UnsetRohTimer ensures that no value is present for RohTimer, not even an explicit nil
### GetObjectProperties

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetObjectProperties() ConfigPutRequestVoicePortProfilesVoicePortProfilesNameObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) GetObjectPropertiesOk() (*ConfigPutRequestVoicePortProfilesVoicePortProfilesNameObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) SetObjectProperties(v ConfigPutRequestVoicePortProfilesVoicePortProfilesNameObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ConfigPutRequestVoicePortProfilesVoicePortProfilesName) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


