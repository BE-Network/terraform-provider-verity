# EthportsettingsPutRequestEthPortSettingsValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. It&#39;s highly recommended to set this value to true so that validation on the object will be ran. | [optional] [default to false]
**AutoNegotiation** | Pointer to **bool** | Indicates if port speed and duplex mode should be auto negotiated | [optional] [default to true]
**MaxBitRate** | Pointer to **string** | Maximum Bit Rate allowed | [optional] [default to "-1"]
**DuplexMode** | Pointer to **string** | Duplex Mode | [optional] [default to "Auto"]
**StpEnable** | Pointer to **bool** | Enable Spanning Tree on the port.  Note: the Spanning Tree Type (VLAN, Port, MST) is controlled in the Site Settings | [optional] [default to false]
**FastLearningMode** | Pointer to **bool** | Enable Immediate Transition to Forwarding | [optional] [default to true]
**BpduGuard** | Pointer to **bool** | Block port on BPDU Receive | [optional] [default to false]
**BpduFilter** | Pointer to **bool** | Drop all Rx and Tx BPDUs | [optional] [default to false]
**GuardLoop** | Pointer to **bool** | Enable Cisco Guard Loop | [optional] [default to false]
**PoeEnable** | Pointer to **bool** | PoE Enable | [optional] [default to false]
**Priority** | Pointer to **string** | Priority given when assigning power in a limited power situation | [optional] [default to "High"]
**AllocatedPower** | Pointer to **string** | Power the PoE system will attempt to allocate on this port | [optional] [default to "0.0"]
**BspEnable** | Pointer to **bool** | Enable Traffic Storm Protection which prevents excessive broadcast/multicast/unknown-unicast traffic from overwhelming the Switch CPU | [optional] [default to false]
**Broadcast** | Pointer to **bool** | Broadcast | [optional] [default to true]
**Multicast** | Pointer to **bool** | Multicast | [optional] [default to true]
**MaxAllowedValue** | Pointer to **NullableInt32** | Max Percentage of the ports bandwidth allowed for broadcast/multicast/unknown-unicast traffic before invoking the protective action | [optional] [default to 1000]
**MaxAllowedUnit** | Pointer to **string** | Max Percentage of the ports bandwidth allowed for broadcast/multicast/unknown-unicast traffic before invoking the protective action &lt;br&gt;                                                 &lt;div class&#x3D;\&quot;tab\&quot;&gt;                                                     %: Percentage.&lt;br&gt;                                                                                                                                                     kbps: kilobits per second &lt;br&gt;                                                     mbps: megabits per second &lt;br&gt;                                                     gbps: gigabits per second &lt;br&gt;                                                     pps: packet per second &lt;br&gt;                                                     kpps: kilopacket per second &lt;br&gt;                                                 &lt;/div&gt;                                                  | [optional] [default to "pps"]
**Action** | Pointer to **string** | Action taken if broadcast/multicast/unknown-unicast traffic excedes the Max. One of: &lt;br&gt;                                                 &lt;div class&#x3D;\&quot;tab\&quot;&gt;                                                     Protect: Broadcast/Multicast packets beyond the percent rate are silently dropped. QOS drop counters should indicate the drops.&lt;br&gt;&lt;br&gt;                                                     Restrict: Broadcast/Multicast packets beyond the percent rate are dropped. QOS drop counters should indicate the drops.                                                      Alarm is raised . Alarm automatically clears when rate is below configured threshold. &lt;br&gt;&lt;br&gt;                                                     Shutdown: Alarm is raised and port is taken out of service. User must administratively Disable and Enable the port to restore service. &lt;br&gt;                                                 &lt;/div&gt;                                              | [optional] [default to "Protect"]
**Fec** | Pointer to **string** | FEC is Forward Error Correction which is error correction on the fiber link.                                                 &lt;div class&#x3D;\&quot;tab\&quot;&gt;                                                     Any: Allows switch Negotiation between FC and RS &lt;br&gt;                                                     None: Disables FEC on an interface.&lt;br&gt;                                                     FC: Enables FEC on supported interfaces. FC stands for fire code.&lt;br&gt;                                                     RS: Enables FEC on supported interfaces. RS stands for Reed-Solomon code. &lt;br&gt;                                                     None: VnetC doesn&#39;t alter the Switch Value.&lt;br&gt;                                                 &lt;/div&gt;                                              | [optional] [default to "unaltered"]
**SingleLink** | Pointer to **bool** | Ports with this setting will be disabled when link state tracking takes effect | [optional] [default to false]
**ObjectProperties** | Pointer to [**EthportsettingsPutRequestEthPortSettingsValueObjectProperties**](EthportsettingsPutRequestEthPortSettingsValueObjectProperties.md) |  | [optional] 

## Methods

### NewEthportsettingsPutRequestEthPortSettingsValue

`func NewEthportsettingsPutRequestEthPortSettingsValue() *EthportsettingsPutRequestEthPortSettingsValue`

NewEthportsettingsPutRequestEthPortSettingsValue instantiates a new EthportsettingsPutRequestEthPortSettingsValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewEthportsettingsPutRequestEthPortSettingsValueWithDefaults

`func NewEthportsettingsPutRequestEthPortSettingsValueWithDefaults() *EthportsettingsPutRequestEthPortSettingsValue`

NewEthportsettingsPutRequestEthPortSettingsValueWithDefaults instantiates a new EthportsettingsPutRequestEthPortSettingsValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *EthportsettingsPutRequestEthPortSettingsValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *EthportsettingsPutRequestEthPortSettingsValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *EthportsettingsPutRequestEthPortSettingsValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *EthportsettingsPutRequestEthPortSettingsValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetAutoNegotiation

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetAutoNegotiation() bool`

GetAutoNegotiation returns the AutoNegotiation field if non-nil, zero value otherwise.

### GetAutoNegotiationOk

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetAutoNegotiationOk() (*bool, bool)`

GetAutoNegotiationOk returns a tuple with the AutoNegotiation field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAutoNegotiation

`func (o *EthportsettingsPutRequestEthPortSettingsValue) SetAutoNegotiation(v bool)`

SetAutoNegotiation sets AutoNegotiation field to given value.

### HasAutoNegotiation

`func (o *EthportsettingsPutRequestEthPortSettingsValue) HasAutoNegotiation() bool`

HasAutoNegotiation returns a boolean if a field has been set.

### GetMaxBitRate

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetMaxBitRate() string`

GetMaxBitRate returns the MaxBitRate field if non-nil, zero value otherwise.

### GetMaxBitRateOk

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetMaxBitRateOk() (*string, bool)`

GetMaxBitRateOk returns a tuple with the MaxBitRate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxBitRate

`func (o *EthportsettingsPutRequestEthPortSettingsValue) SetMaxBitRate(v string)`

SetMaxBitRate sets MaxBitRate field to given value.

### HasMaxBitRate

`func (o *EthportsettingsPutRequestEthPortSettingsValue) HasMaxBitRate() bool`

HasMaxBitRate returns a boolean if a field has been set.

### GetDuplexMode

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetDuplexMode() string`

GetDuplexMode returns the DuplexMode field if non-nil, zero value otherwise.

### GetDuplexModeOk

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetDuplexModeOk() (*string, bool)`

GetDuplexModeOk returns a tuple with the DuplexMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDuplexMode

`func (o *EthportsettingsPutRequestEthPortSettingsValue) SetDuplexMode(v string)`

SetDuplexMode sets DuplexMode field to given value.

### HasDuplexMode

`func (o *EthportsettingsPutRequestEthPortSettingsValue) HasDuplexMode() bool`

HasDuplexMode returns a boolean if a field has been set.

### GetStpEnable

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetStpEnable() bool`

GetStpEnable returns the StpEnable field if non-nil, zero value otherwise.

### GetStpEnableOk

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetStpEnableOk() (*bool, bool)`

GetStpEnableOk returns a tuple with the StpEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStpEnable

`func (o *EthportsettingsPutRequestEthPortSettingsValue) SetStpEnable(v bool)`

SetStpEnable sets StpEnable field to given value.

### HasStpEnable

`func (o *EthportsettingsPutRequestEthPortSettingsValue) HasStpEnable() bool`

HasStpEnable returns a boolean if a field has been set.

### GetFastLearningMode

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetFastLearningMode() bool`

GetFastLearningMode returns the FastLearningMode field if non-nil, zero value otherwise.

### GetFastLearningModeOk

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetFastLearningModeOk() (*bool, bool)`

GetFastLearningModeOk returns a tuple with the FastLearningMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFastLearningMode

`func (o *EthportsettingsPutRequestEthPortSettingsValue) SetFastLearningMode(v bool)`

SetFastLearningMode sets FastLearningMode field to given value.

### HasFastLearningMode

`func (o *EthportsettingsPutRequestEthPortSettingsValue) HasFastLearningMode() bool`

HasFastLearningMode returns a boolean if a field has been set.

### GetBpduGuard

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetBpduGuard() bool`

GetBpduGuard returns the BpduGuard field if non-nil, zero value otherwise.

### GetBpduGuardOk

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetBpduGuardOk() (*bool, bool)`

GetBpduGuardOk returns a tuple with the BpduGuard field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBpduGuard

`func (o *EthportsettingsPutRequestEthPortSettingsValue) SetBpduGuard(v bool)`

SetBpduGuard sets BpduGuard field to given value.

### HasBpduGuard

`func (o *EthportsettingsPutRequestEthPortSettingsValue) HasBpduGuard() bool`

HasBpduGuard returns a boolean if a field has been set.

### GetBpduFilter

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetBpduFilter() bool`

GetBpduFilter returns the BpduFilter field if non-nil, zero value otherwise.

### GetBpduFilterOk

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetBpduFilterOk() (*bool, bool)`

GetBpduFilterOk returns a tuple with the BpduFilter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBpduFilter

`func (o *EthportsettingsPutRequestEthPortSettingsValue) SetBpduFilter(v bool)`

SetBpduFilter sets BpduFilter field to given value.

### HasBpduFilter

`func (o *EthportsettingsPutRequestEthPortSettingsValue) HasBpduFilter() bool`

HasBpduFilter returns a boolean if a field has been set.

### GetGuardLoop

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetGuardLoop() bool`

GetGuardLoop returns the GuardLoop field if non-nil, zero value otherwise.

### GetGuardLoopOk

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetGuardLoopOk() (*bool, bool)`

GetGuardLoopOk returns a tuple with the GuardLoop field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGuardLoop

`func (o *EthportsettingsPutRequestEthPortSettingsValue) SetGuardLoop(v bool)`

SetGuardLoop sets GuardLoop field to given value.

### HasGuardLoop

`func (o *EthportsettingsPutRequestEthPortSettingsValue) HasGuardLoop() bool`

HasGuardLoop returns a boolean if a field has been set.

### GetPoeEnable

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetPoeEnable() bool`

GetPoeEnable returns the PoeEnable field if non-nil, zero value otherwise.

### GetPoeEnableOk

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetPoeEnableOk() (*bool, bool)`

GetPoeEnableOk returns a tuple with the PoeEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoeEnable

`func (o *EthportsettingsPutRequestEthPortSettingsValue) SetPoeEnable(v bool)`

SetPoeEnable sets PoeEnable field to given value.

### HasPoeEnable

`func (o *EthportsettingsPutRequestEthPortSettingsValue) HasPoeEnable() bool`

HasPoeEnable returns a boolean if a field has been set.

### GetPriority

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetPriority() string`

GetPriority returns the Priority field if non-nil, zero value otherwise.

### GetPriorityOk

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetPriorityOk() (*string, bool)`

GetPriorityOk returns a tuple with the Priority field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPriority

`func (o *EthportsettingsPutRequestEthPortSettingsValue) SetPriority(v string)`

SetPriority sets Priority field to given value.

### HasPriority

`func (o *EthportsettingsPutRequestEthPortSettingsValue) HasPriority() bool`

HasPriority returns a boolean if a field has been set.

### GetAllocatedPower

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetAllocatedPower() string`

GetAllocatedPower returns the AllocatedPower field if non-nil, zero value otherwise.

### GetAllocatedPowerOk

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetAllocatedPowerOk() (*string, bool)`

GetAllocatedPowerOk returns a tuple with the AllocatedPower field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllocatedPower

`func (o *EthportsettingsPutRequestEthPortSettingsValue) SetAllocatedPower(v string)`

SetAllocatedPower sets AllocatedPower field to given value.

### HasAllocatedPower

`func (o *EthportsettingsPutRequestEthPortSettingsValue) HasAllocatedPower() bool`

HasAllocatedPower returns a boolean if a field has been set.

### GetBspEnable

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetBspEnable() bool`

GetBspEnable returns the BspEnable field if non-nil, zero value otherwise.

### GetBspEnableOk

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetBspEnableOk() (*bool, bool)`

GetBspEnableOk returns a tuple with the BspEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBspEnable

`func (o *EthportsettingsPutRequestEthPortSettingsValue) SetBspEnable(v bool)`

SetBspEnable sets BspEnable field to given value.

### HasBspEnable

`func (o *EthportsettingsPutRequestEthPortSettingsValue) HasBspEnable() bool`

HasBspEnable returns a boolean if a field has been set.

### GetBroadcast

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetBroadcast() bool`

GetBroadcast returns the Broadcast field if non-nil, zero value otherwise.

### GetBroadcastOk

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetBroadcastOk() (*bool, bool)`

GetBroadcastOk returns a tuple with the Broadcast field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBroadcast

`func (o *EthportsettingsPutRequestEthPortSettingsValue) SetBroadcast(v bool)`

SetBroadcast sets Broadcast field to given value.

### HasBroadcast

`func (o *EthportsettingsPutRequestEthPortSettingsValue) HasBroadcast() bool`

HasBroadcast returns a boolean if a field has been set.

### GetMulticast

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetMulticast() bool`

GetMulticast returns the Multicast field if non-nil, zero value otherwise.

### GetMulticastOk

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetMulticastOk() (*bool, bool)`

GetMulticastOk returns a tuple with the Multicast field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMulticast

`func (o *EthportsettingsPutRequestEthPortSettingsValue) SetMulticast(v bool)`

SetMulticast sets Multicast field to given value.

### HasMulticast

`func (o *EthportsettingsPutRequestEthPortSettingsValue) HasMulticast() bool`

HasMulticast returns a boolean if a field has been set.

### GetMaxAllowedValue

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetMaxAllowedValue() int32`

GetMaxAllowedValue returns the MaxAllowedValue field if non-nil, zero value otherwise.

### GetMaxAllowedValueOk

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetMaxAllowedValueOk() (*int32, bool)`

GetMaxAllowedValueOk returns a tuple with the MaxAllowedValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxAllowedValue

`func (o *EthportsettingsPutRequestEthPortSettingsValue) SetMaxAllowedValue(v int32)`

SetMaxAllowedValue sets MaxAllowedValue field to given value.

### HasMaxAllowedValue

`func (o *EthportsettingsPutRequestEthPortSettingsValue) HasMaxAllowedValue() bool`

HasMaxAllowedValue returns a boolean if a field has been set.

### SetMaxAllowedValueNil

`func (o *EthportsettingsPutRequestEthPortSettingsValue) SetMaxAllowedValueNil(b bool)`

 SetMaxAllowedValueNil sets the value for MaxAllowedValue to be an explicit nil

### UnsetMaxAllowedValue
`func (o *EthportsettingsPutRequestEthPortSettingsValue) UnsetMaxAllowedValue()`

UnsetMaxAllowedValue ensures that no value is present for MaxAllowedValue, not even an explicit nil
### GetMaxAllowedUnit

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetMaxAllowedUnit() string`

GetMaxAllowedUnit returns the MaxAllowedUnit field if non-nil, zero value otherwise.

### GetMaxAllowedUnitOk

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetMaxAllowedUnitOk() (*string, bool)`

GetMaxAllowedUnitOk returns a tuple with the MaxAllowedUnit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxAllowedUnit

`func (o *EthportsettingsPutRequestEthPortSettingsValue) SetMaxAllowedUnit(v string)`

SetMaxAllowedUnit sets MaxAllowedUnit field to given value.

### HasMaxAllowedUnit

`func (o *EthportsettingsPutRequestEthPortSettingsValue) HasMaxAllowedUnit() bool`

HasMaxAllowedUnit returns a boolean if a field has been set.

### GetAction

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetAction() string`

GetAction returns the Action field if non-nil, zero value otherwise.

### GetActionOk

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetActionOk() (*string, bool)`

GetActionOk returns a tuple with the Action field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAction

`func (o *EthportsettingsPutRequestEthPortSettingsValue) SetAction(v string)`

SetAction sets Action field to given value.

### HasAction

`func (o *EthportsettingsPutRequestEthPortSettingsValue) HasAction() bool`

HasAction returns a boolean if a field has been set.

### GetFec

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetFec() string`

GetFec returns the Fec field if non-nil, zero value otherwise.

### GetFecOk

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetFecOk() (*string, bool)`

GetFecOk returns a tuple with the Fec field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFec

`func (o *EthportsettingsPutRequestEthPortSettingsValue) SetFec(v string)`

SetFec sets Fec field to given value.

### HasFec

`func (o *EthportsettingsPutRequestEthPortSettingsValue) HasFec() bool`

HasFec returns a boolean if a field has been set.

### GetSingleLink

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetSingleLink() bool`

GetSingleLink returns the SingleLink field if non-nil, zero value otherwise.

### GetSingleLinkOk

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetSingleLinkOk() (*bool, bool)`

GetSingleLinkOk returns a tuple with the SingleLink field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSingleLink

`func (o *EthportsettingsPutRequestEthPortSettingsValue) SetSingleLink(v bool)`

SetSingleLink sets SingleLink field to given value.

### HasSingleLink

`func (o *EthportsettingsPutRequestEthPortSettingsValue) HasSingleLink() bool`

HasSingleLink returns a boolean if a field has been set.

### GetObjectProperties

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetObjectProperties() EthportsettingsPutRequestEthPortSettingsValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *EthportsettingsPutRequestEthPortSettingsValue) GetObjectPropertiesOk() (*EthportsettingsPutRequestEthPortSettingsValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *EthportsettingsPutRequestEthPortSettingsValue) SetObjectProperties(v EthportsettingsPutRequestEthPortSettingsValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *EthportsettingsPutRequestEthPortSettingsValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


