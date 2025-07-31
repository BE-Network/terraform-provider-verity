# VoiceportprofilesPutRequestVoicePortProfilesValue

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
**ObjectProperties** | Pointer to [**VoiceportprofilesPutRequestVoicePortProfilesValueObjectProperties**](VoiceportprofilesPutRequestVoicePortProfilesValueObjectProperties.md) |  | [optional] 

## Methods

### NewVoiceportprofilesPutRequestVoicePortProfilesValue

`func NewVoiceportprofilesPutRequestVoicePortProfilesValue() *VoiceportprofilesPutRequestVoicePortProfilesValue`

NewVoiceportprofilesPutRequestVoicePortProfilesValue instantiates a new VoiceportprofilesPutRequestVoicePortProfilesValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewVoiceportprofilesPutRequestVoicePortProfilesValueWithDefaults

`func NewVoiceportprofilesPutRequestVoicePortProfilesValueWithDefaults() *VoiceportprofilesPutRequestVoicePortProfilesValue`

NewVoiceportprofilesPutRequestVoicePortProfilesValueWithDefaults instantiates a new VoiceportprofilesPutRequestVoicePortProfilesValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetProtocol

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetProtocol() string`

GetProtocol returns the Protocol field if non-nil, zero value otherwise.

### GetProtocolOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetProtocolOk() (*string, bool)`

GetProtocolOk returns a tuple with the Protocol field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProtocol

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetProtocol(v string)`

SetProtocol sets Protocol field to given value.

### HasProtocol

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasProtocol() bool`

HasProtocol returns a boolean if a field has been set.

### GetDigitMap

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetDigitMap() string`

GetDigitMap returns the DigitMap field if non-nil, zero value otherwise.

### GetDigitMapOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetDigitMapOk() (*string, bool)`

GetDigitMapOk returns a tuple with the DigitMap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDigitMap

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetDigitMap(v string)`

SetDigitMap sets DigitMap field to given value.

### HasDigitMap

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasDigitMap() bool`

HasDigitMap returns a boolean if a field has been set.

### GetCallThreeWayEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetCallThreeWayEnable() bool`

GetCallThreeWayEnable returns the CallThreeWayEnable field if non-nil, zero value otherwise.

### GetCallThreeWayEnableOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetCallThreeWayEnableOk() (*bool, bool)`

GetCallThreeWayEnableOk returns a tuple with the CallThreeWayEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallThreeWayEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetCallThreeWayEnable(v bool)`

SetCallThreeWayEnable sets CallThreeWayEnable field to given value.

### HasCallThreeWayEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasCallThreeWayEnable() bool`

HasCallThreeWayEnable returns a boolean if a field has been set.

### GetCallerIdEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetCallerIdEnable() bool`

GetCallerIdEnable returns the CallerIdEnable field if non-nil, zero value otherwise.

### GetCallerIdEnableOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetCallerIdEnableOk() (*bool, bool)`

GetCallerIdEnableOk returns a tuple with the CallerIdEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallerIdEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetCallerIdEnable(v bool)`

SetCallerIdEnable sets CallerIdEnable field to given value.

### HasCallerIdEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasCallerIdEnable() bool`

HasCallerIdEnable returns a boolean if a field has been set.

### GetCallerIdNameEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetCallerIdNameEnable() bool`

GetCallerIdNameEnable returns the CallerIdNameEnable field if non-nil, zero value otherwise.

### GetCallerIdNameEnableOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetCallerIdNameEnableOk() (*bool, bool)`

GetCallerIdNameEnableOk returns a tuple with the CallerIdNameEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallerIdNameEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetCallerIdNameEnable(v bool)`

SetCallerIdNameEnable sets CallerIdNameEnable field to given value.

### HasCallerIdNameEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasCallerIdNameEnable() bool`

HasCallerIdNameEnable returns a boolean if a field has been set.

### GetCallWaitingEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetCallWaitingEnable() bool`

GetCallWaitingEnable returns the CallWaitingEnable field if non-nil, zero value otherwise.

### GetCallWaitingEnableOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetCallWaitingEnableOk() (*bool, bool)`

GetCallWaitingEnableOk returns a tuple with the CallWaitingEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallWaitingEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetCallWaitingEnable(v bool)`

SetCallWaitingEnable sets CallWaitingEnable field to given value.

### HasCallWaitingEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasCallWaitingEnable() bool`

HasCallWaitingEnable returns a boolean if a field has been set.

### GetCallForwardUnconditionalEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetCallForwardUnconditionalEnable() bool`

GetCallForwardUnconditionalEnable returns the CallForwardUnconditionalEnable field if non-nil, zero value otherwise.

### GetCallForwardUnconditionalEnableOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetCallForwardUnconditionalEnableOk() (*bool, bool)`

GetCallForwardUnconditionalEnableOk returns a tuple with the CallForwardUnconditionalEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallForwardUnconditionalEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetCallForwardUnconditionalEnable(v bool)`

SetCallForwardUnconditionalEnable sets CallForwardUnconditionalEnable field to given value.

### HasCallForwardUnconditionalEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasCallForwardUnconditionalEnable() bool`

HasCallForwardUnconditionalEnable returns a boolean if a field has been set.

### GetCallForwardOnBusyEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetCallForwardOnBusyEnable() bool`

GetCallForwardOnBusyEnable returns the CallForwardOnBusyEnable field if non-nil, zero value otherwise.

### GetCallForwardOnBusyEnableOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetCallForwardOnBusyEnableOk() (*bool, bool)`

GetCallForwardOnBusyEnableOk returns a tuple with the CallForwardOnBusyEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallForwardOnBusyEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetCallForwardOnBusyEnable(v bool)`

SetCallForwardOnBusyEnable sets CallForwardOnBusyEnable field to given value.

### HasCallForwardOnBusyEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasCallForwardOnBusyEnable() bool`

HasCallForwardOnBusyEnable returns a boolean if a field has been set.

### GetCallForwardOnNoAnswerRingCount

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetCallForwardOnNoAnswerRingCount() int32`

GetCallForwardOnNoAnswerRingCount returns the CallForwardOnNoAnswerRingCount field if non-nil, zero value otherwise.

### GetCallForwardOnNoAnswerRingCountOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetCallForwardOnNoAnswerRingCountOk() (*int32, bool)`

GetCallForwardOnNoAnswerRingCountOk returns a tuple with the CallForwardOnNoAnswerRingCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallForwardOnNoAnswerRingCount

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetCallForwardOnNoAnswerRingCount(v int32)`

SetCallForwardOnNoAnswerRingCount sets CallForwardOnNoAnswerRingCount field to given value.

### HasCallForwardOnNoAnswerRingCount

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasCallForwardOnNoAnswerRingCount() bool`

HasCallForwardOnNoAnswerRingCount returns a boolean if a field has been set.

### SetCallForwardOnNoAnswerRingCountNil

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetCallForwardOnNoAnswerRingCountNil(b bool)`

 SetCallForwardOnNoAnswerRingCountNil sets the value for CallForwardOnNoAnswerRingCount to be an explicit nil

### UnsetCallForwardOnNoAnswerRingCount
`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) UnsetCallForwardOnNoAnswerRingCount()`

UnsetCallForwardOnNoAnswerRingCount ensures that no value is present for CallForwardOnNoAnswerRingCount, not even an explicit nil
### GetCallTransferEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetCallTransferEnable() bool`

GetCallTransferEnable returns the CallTransferEnable field if non-nil, zero value otherwise.

### GetCallTransferEnableOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetCallTransferEnableOk() (*bool, bool)`

GetCallTransferEnableOk returns a tuple with the CallTransferEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallTransferEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetCallTransferEnable(v bool)`

SetCallTransferEnable sets CallTransferEnable field to given value.

### HasCallTransferEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasCallTransferEnable() bool`

HasCallTransferEnable returns a boolean if a field has been set.

### GetAudioMwiEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetAudioMwiEnable() bool`

GetAudioMwiEnable returns the AudioMwiEnable field if non-nil, zero value otherwise.

### GetAudioMwiEnableOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetAudioMwiEnableOk() (*bool, bool)`

GetAudioMwiEnableOk returns a tuple with the AudioMwiEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAudioMwiEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetAudioMwiEnable(v bool)`

SetAudioMwiEnable sets AudioMwiEnable field to given value.

### HasAudioMwiEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasAudioMwiEnable() bool`

HasAudioMwiEnable returns a boolean if a field has been set.

### GetAnonymousCallBlockEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetAnonymousCallBlockEnable() bool`

GetAnonymousCallBlockEnable returns the AnonymousCallBlockEnable field if non-nil, zero value otherwise.

### GetAnonymousCallBlockEnableOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetAnonymousCallBlockEnableOk() (*bool, bool)`

GetAnonymousCallBlockEnableOk returns a tuple with the AnonymousCallBlockEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnonymousCallBlockEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetAnonymousCallBlockEnable(v bool)`

SetAnonymousCallBlockEnable sets AnonymousCallBlockEnable field to given value.

### HasAnonymousCallBlockEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasAnonymousCallBlockEnable() bool`

HasAnonymousCallBlockEnable returns a boolean if a field has been set.

### GetDoNotDisturbEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetDoNotDisturbEnable() bool`

GetDoNotDisturbEnable returns the DoNotDisturbEnable field if non-nil, zero value otherwise.

### GetDoNotDisturbEnableOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetDoNotDisturbEnableOk() (*bool, bool)`

GetDoNotDisturbEnableOk returns a tuple with the DoNotDisturbEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDoNotDisturbEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetDoNotDisturbEnable(v bool)`

SetDoNotDisturbEnable sets DoNotDisturbEnable field to given value.

### HasDoNotDisturbEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasDoNotDisturbEnable() bool`

HasDoNotDisturbEnable returns a boolean if a field has been set.

### GetCidBlockingEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetCidBlockingEnable() bool`

GetCidBlockingEnable returns the CidBlockingEnable field if non-nil, zero value otherwise.

### GetCidBlockingEnableOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetCidBlockingEnableOk() (*bool, bool)`

GetCidBlockingEnableOk returns a tuple with the CidBlockingEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCidBlockingEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetCidBlockingEnable(v bool)`

SetCidBlockingEnable sets CidBlockingEnable field to given value.

### HasCidBlockingEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasCidBlockingEnable() bool`

HasCidBlockingEnable returns a boolean if a field has been set.

### GetCidNumPresentationStatus

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetCidNumPresentationStatus() string`

GetCidNumPresentationStatus returns the CidNumPresentationStatus field if non-nil, zero value otherwise.

### GetCidNumPresentationStatusOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetCidNumPresentationStatusOk() (*string, bool)`

GetCidNumPresentationStatusOk returns a tuple with the CidNumPresentationStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCidNumPresentationStatus

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetCidNumPresentationStatus(v string)`

SetCidNumPresentationStatus sets CidNumPresentationStatus field to given value.

### HasCidNumPresentationStatus

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasCidNumPresentationStatus() bool`

HasCidNumPresentationStatus returns a boolean if a field has been set.

### GetCidNamePresentationStatus

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetCidNamePresentationStatus() string`

GetCidNamePresentationStatus returns the CidNamePresentationStatus field if non-nil, zero value otherwise.

### GetCidNamePresentationStatusOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetCidNamePresentationStatusOk() (*string, bool)`

GetCidNamePresentationStatusOk returns a tuple with the CidNamePresentationStatus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCidNamePresentationStatus

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetCidNamePresentationStatus(v string)`

SetCidNamePresentationStatus sets CidNamePresentationStatus field to given value.

### HasCidNamePresentationStatus

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasCidNamePresentationStatus() bool`

HasCidNamePresentationStatus returns a boolean if a field has been set.

### GetCallWaitingCallerIdEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetCallWaitingCallerIdEnable() bool`

GetCallWaitingCallerIdEnable returns the CallWaitingCallerIdEnable field if non-nil, zero value otherwise.

### GetCallWaitingCallerIdEnableOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetCallWaitingCallerIdEnableOk() (*bool, bool)`

GetCallWaitingCallerIdEnableOk returns a tuple with the CallWaitingCallerIdEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallWaitingCallerIdEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetCallWaitingCallerIdEnable(v bool)`

SetCallWaitingCallerIdEnable sets CallWaitingCallerIdEnable field to given value.

### HasCallWaitingCallerIdEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasCallWaitingCallerIdEnable() bool`

HasCallWaitingCallerIdEnable returns a boolean if a field has been set.

### GetCallHoldEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetCallHoldEnable() bool`

GetCallHoldEnable returns the CallHoldEnable field if non-nil, zero value otherwise.

### GetCallHoldEnableOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetCallHoldEnableOk() (*bool, bool)`

GetCallHoldEnableOk returns a tuple with the CallHoldEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallHoldEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetCallHoldEnable(v bool)`

SetCallHoldEnable sets CallHoldEnable field to given value.

### HasCallHoldEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasCallHoldEnable() bool`

HasCallHoldEnable returns a boolean if a field has been set.

### GetVisualMwiEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetVisualMwiEnable() bool`

GetVisualMwiEnable returns the VisualMwiEnable field if non-nil, zero value otherwise.

### GetVisualMwiEnableOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetVisualMwiEnableOk() (*bool, bool)`

GetVisualMwiEnableOk returns a tuple with the VisualMwiEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVisualMwiEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetVisualMwiEnable(v bool)`

SetVisualMwiEnable sets VisualMwiEnable field to given value.

### HasVisualMwiEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasVisualMwiEnable() bool`

HasVisualMwiEnable returns a boolean if a field has been set.

### GetMwiRefreshTimer

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetMwiRefreshTimer() int32`

GetMwiRefreshTimer returns the MwiRefreshTimer field if non-nil, zero value otherwise.

### GetMwiRefreshTimerOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetMwiRefreshTimerOk() (*int32, bool)`

GetMwiRefreshTimerOk returns a tuple with the MwiRefreshTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMwiRefreshTimer

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetMwiRefreshTimer(v int32)`

SetMwiRefreshTimer sets MwiRefreshTimer field to given value.

### HasMwiRefreshTimer

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasMwiRefreshTimer() bool`

HasMwiRefreshTimer returns a boolean if a field has been set.

### SetMwiRefreshTimerNil

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetMwiRefreshTimerNil(b bool)`

 SetMwiRefreshTimerNil sets the value for MwiRefreshTimer to be an explicit nil

### UnsetMwiRefreshTimer
`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) UnsetMwiRefreshTimer()`

UnsetMwiRefreshTimer ensures that no value is present for MwiRefreshTimer, not even an explicit nil
### GetHotlineEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetHotlineEnable() bool`

GetHotlineEnable returns the HotlineEnable field if non-nil, zero value otherwise.

### GetHotlineEnableOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetHotlineEnableOk() (*bool, bool)`

GetHotlineEnableOk returns a tuple with the HotlineEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHotlineEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetHotlineEnable(v bool)`

SetHotlineEnable sets HotlineEnable field to given value.

### HasHotlineEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasHotlineEnable() bool`

HasHotlineEnable returns a boolean if a field has been set.

### GetDialToneFeatureDelay

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetDialToneFeatureDelay() int32`

GetDialToneFeatureDelay returns the DialToneFeatureDelay field if non-nil, zero value otherwise.

### GetDialToneFeatureDelayOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetDialToneFeatureDelayOk() (*int32, bool)`

GetDialToneFeatureDelayOk returns a tuple with the DialToneFeatureDelay field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDialToneFeatureDelay

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetDialToneFeatureDelay(v int32)`

SetDialToneFeatureDelay sets DialToneFeatureDelay field to given value.

### HasDialToneFeatureDelay

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasDialToneFeatureDelay() bool`

HasDialToneFeatureDelay returns a boolean if a field has been set.

### SetDialToneFeatureDelayNil

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetDialToneFeatureDelayNil(b bool)`

 SetDialToneFeatureDelayNil sets the value for DialToneFeatureDelay to be an explicit nil

### UnsetDialToneFeatureDelay
`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) UnsetDialToneFeatureDelay()`

UnsetDialToneFeatureDelay ensures that no value is present for DialToneFeatureDelay, not even an explicit nil
### GetIntercomEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetIntercomEnable() bool`

GetIntercomEnable returns the IntercomEnable field if non-nil, zero value otherwise.

### GetIntercomEnableOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetIntercomEnableOk() (*bool, bool)`

GetIntercomEnableOk returns a tuple with the IntercomEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIntercomEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetIntercomEnable(v bool)`

SetIntercomEnable sets IntercomEnable field to given value.

### HasIntercomEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasIntercomEnable() bool`

HasIntercomEnable returns a boolean if a field has been set.

### GetIntercomTransferEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetIntercomTransferEnable() bool`

GetIntercomTransferEnable returns the IntercomTransferEnable field if non-nil, zero value otherwise.

### GetIntercomTransferEnableOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetIntercomTransferEnableOk() (*bool, bool)`

GetIntercomTransferEnableOk returns a tuple with the IntercomTransferEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIntercomTransferEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetIntercomTransferEnable(v bool)`

SetIntercomTransferEnable sets IntercomTransferEnable field to given value.

### HasIntercomTransferEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasIntercomTransferEnable() bool`

HasIntercomTransferEnable returns a boolean if a field has been set.

### GetTransmitGain

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetTransmitGain() int32`

GetTransmitGain returns the TransmitGain field if non-nil, zero value otherwise.

### GetTransmitGainOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetTransmitGainOk() (*int32, bool)`

GetTransmitGainOk returns a tuple with the TransmitGain field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTransmitGain

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetTransmitGain(v int32)`

SetTransmitGain sets TransmitGain field to given value.

### HasTransmitGain

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasTransmitGain() bool`

HasTransmitGain returns a boolean if a field has been set.

### SetTransmitGainNil

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetTransmitGainNil(b bool)`

 SetTransmitGainNil sets the value for TransmitGain to be an explicit nil

### UnsetTransmitGain
`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) UnsetTransmitGain()`

UnsetTransmitGain ensures that no value is present for TransmitGain, not even an explicit nil
### GetReceiveGain

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetReceiveGain() int32`

GetReceiveGain returns the ReceiveGain field if non-nil, zero value otherwise.

### GetReceiveGainOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetReceiveGainOk() (*int32, bool)`

GetReceiveGainOk returns a tuple with the ReceiveGain field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReceiveGain

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetReceiveGain(v int32)`

SetReceiveGain sets ReceiveGain field to given value.

### HasReceiveGain

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasReceiveGain() bool`

HasReceiveGain returns a boolean if a field has been set.

### SetReceiveGainNil

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetReceiveGainNil(b bool)`

 SetReceiveGainNil sets the value for ReceiveGain to be an explicit nil

### UnsetReceiveGain
`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) UnsetReceiveGain()`

UnsetReceiveGain ensures that no value is present for ReceiveGain, not even an explicit nil
### GetEchoCancellationEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetEchoCancellationEnable() bool`

GetEchoCancellationEnable returns the EchoCancellationEnable field if non-nil, zero value otherwise.

### GetEchoCancellationEnableOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetEchoCancellationEnableOk() (*bool, bool)`

GetEchoCancellationEnableOk returns a tuple with the EchoCancellationEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEchoCancellationEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetEchoCancellationEnable(v bool)`

SetEchoCancellationEnable sets EchoCancellationEnable field to given value.

### HasEchoCancellationEnable

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasEchoCancellationEnable() bool`

HasEchoCancellationEnable returns a boolean if a field has been set.

### GetJitterTarget

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetJitterTarget() int32`

GetJitterTarget returns the JitterTarget field if non-nil, zero value otherwise.

### GetJitterTargetOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetJitterTargetOk() (*int32, bool)`

GetJitterTargetOk returns a tuple with the JitterTarget field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetJitterTarget

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetJitterTarget(v int32)`

SetJitterTarget sets JitterTarget field to given value.

### HasJitterTarget

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasJitterTarget() bool`

HasJitterTarget returns a boolean if a field has been set.

### SetJitterTargetNil

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetJitterTargetNil(b bool)`

 SetJitterTargetNil sets the value for JitterTarget to be an explicit nil

### UnsetJitterTarget
`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) UnsetJitterTarget()`

UnsetJitterTarget ensures that no value is present for JitterTarget, not even an explicit nil
### GetJitterBufferMax

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetJitterBufferMax() int32`

GetJitterBufferMax returns the JitterBufferMax field if non-nil, zero value otherwise.

### GetJitterBufferMaxOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetJitterBufferMaxOk() (*int32, bool)`

GetJitterBufferMaxOk returns a tuple with the JitterBufferMax field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetJitterBufferMax

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetJitterBufferMax(v int32)`

SetJitterBufferMax sets JitterBufferMax field to given value.

### HasJitterBufferMax

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasJitterBufferMax() bool`

HasJitterBufferMax returns a boolean if a field has been set.

### SetJitterBufferMaxNil

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetJitterBufferMaxNil(b bool)`

 SetJitterBufferMaxNil sets the value for JitterBufferMax to be an explicit nil

### UnsetJitterBufferMax
`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) UnsetJitterBufferMax()`

UnsetJitterBufferMax ensures that no value is present for JitterBufferMax, not even an explicit nil
### GetSignalingCode

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetSignalingCode() string`

GetSignalingCode returns the SignalingCode field if non-nil, zero value otherwise.

### GetSignalingCodeOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetSignalingCodeOk() (*string, bool)`

GetSignalingCodeOk returns a tuple with the SignalingCode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSignalingCode

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetSignalingCode(v string)`

SetSignalingCode sets SignalingCode field to given value.

### HasSignalingCode

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasSignalingCode() bool`

HasSignalingCode returns a boolean if a field has been set.

### GetReleaseTimer

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetReleaseTimer() int32`

GetReleaseTimer returns the ReleaseTimer field if non-nil, zero value otherwise.

### GetReleaseTimerOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetReleaseTimerOk() (*int32, bool)`

GetReleaseTimerOk returns a tuple with the ReleaseTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReleaseTimer

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetReleaseTimer(v int32)`

SetReleaseTimer sets ReleaseTimer field to given value.

### HasReleaseTimer

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasReleaseTimer() bool`

HasReleaseTimer returns a boolean if a field has been set.

### SetReleaseTimerNil

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetReleaseTimerNil(b bool)`

 SetReleaseTimerNil sets the value for ReleaseTimer to be an explicit nil

### UnsetReleaseTimer
`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) UnsetReleaseTimer()`

UnsetReleaseTimer ensures that no value is present for ReleaseTimer, not even an explicit nil
### GetRohTimer

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetRohTimer() int32`

GetRohTimer returns the RohTimer field if non-nil, zero value otherwise.

### GetRohTimerOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetRohTimerOk() (*int32, bool)`

GetRohTimerOk returns a tuple with the RohTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRohTimer

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetRohTimer(v int32)`

SetRohTimer sets RohTimer field to given value.

### HasRohTimer

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasRohTimer() bool`

HasRohTimer returns a boolean if a field has been set.

### SetRohTimerNil

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetRohTimerNil(b bool)`

 SetRohTimerNil sets the value for RohTimer to be an explicit nil

### UnsetRohTimer
`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) UnsetRohTimer()`

UnsetRohTimer ensures that no value is present for RohTimer, not even an explicit nil
### GetObjectProperties

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetObjectProperties() VoiceportprofilesPutRequestVoicePortProfilesValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) GetObjectPropertiesOk() (*VoiceportprofilesPutRequestVoicePortProfilesValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) SetObjectProperties(v VoiceportprofilesPutRequestVoicePortProfilesValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *VoiceportprofilesPutRequestVoicePortProfilesValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


