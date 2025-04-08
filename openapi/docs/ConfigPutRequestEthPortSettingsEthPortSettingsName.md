# ConfigPutRequestEthPortSettingsEthPortSettingsName

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
**MaxAllowedValue** | Pointer to **int32** | Max Percentage of the ports bandwidth allowed for broadcast/multicast/unknown-unicast traffic before invoking the protective action | [optional] 
**MaxAllowedUnit** | Pointer to **string** | Max Percentage of the ports bandwidth allowed for broadcast/multicast/unknown-unicast traffic before invoking the protective action &lt;br&gt;                                                 &lt;div class&#x3D;\&quot;tab\&quot;&gt;                                                     %: Percentage.&lt;br&gt;                                                                                                                                                     kbps: kilobits per second &lt;br&gt;                                                     mbps: megabits per second &lt;br&gt;                                                     gbps: gigabits per second &lt;br&gt;                                                     pps: packet per second &lt;br&gt;                                                     kpps: kilopacket per second &lt;br&gt;                                                 &lt;/div&gt;                                                  | [optional] [default to "pps"]
**Action** | Pointer to **string** | Action taken if broadcast/multicast/unknown-unicast traffic excedes the Max. One of: &lt;br&gt;                                                 &lt;div class&#x3D;\&quot;tab\&quot;&gt;                                                     Protect: Broadcast/Multicast packets beyond the percent rate are silently dropped. QOS drop counters should indicate the drops.&lt;br&gt;&lt;br&gt;                                                     Restrict: Broadcast/Multicast packets beyond the percent rate are dropped. QOS drop counters should indicate the drops.                                                      Alarm is raised . Alarm automatically clears when rate is below configured threshold. &lt;br&gt;&lt;br&gt;                                                     Shutdown: Alarm is raised and port is taken out of service. User must administratively Disable and Enable the port to restore service. &lt;br&gt;                                                 &lt;/div&gt;                                              | [optional] [default to "Protect"]
**Fec** | Pointer to **string** | FEC is Forward Error Correction which is error correction on the fiber link.                                                 &lt;div class&#x3D;\&quot;tab\&quot;&gt;                                                     Any: Allows switch Negotiation between FC and RS &lt;br&gt;                                                     None: Disables FEC on an interface.&lt;br&gt;                                                     FC: Enables FEC on supported interfaces. FC stands for fire code.&lt;br&gt;                                                     RS: Enables FEC on supported interfaces. RS stands for Reed-Solomon code. &lt;br&gt;                                                     None: VnetC doesn&#39;t alter the Switch Value.&lt;br&gt;                                                 &lt;/div&gt;                                              | [optional] [default to "unaltered"]
**SingleLink** | Pointer to **bool** | Ports with this setting will be disabled when link state tracking takes effect | [optional] [default to false]
**ObjectProperties** | Pointer to [**ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties**](ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties.md) |  | [optional] 

## Methods

### NewConfigPutRequestEthPortSettingsEthPortSettingsName

`func NewConfigPutRequestEthPortSettingsEthPortSettingsName() *ConfigPutRequestEthPortSettingsEthPortSettingsName`

NewConfigPutRequestEthPortSettingsEthPortSettingsName instantiates a new ConfigPutRequestEthPortSettingsEthPortSettingsName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestEthPortSettingsEthPortSettingsNameWithDefaults

`func NewConfigPutRequestEthPortSettingsEthPortSettingsNameWithDefaults() *ConfigPutRequestEthPortSettingsEthPortSettingsName`

NewConfigPutRequestEthPortSettingsEthPortSettingsNameWithDefaults instantiates a new ConfigPutRequestEthPortSettingsEthPortSettingsName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetAutoNegotiation

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetAutoNegotiation() bool`

GetAutoNegotiation returns the AutoNegotiation field if non-nil, zero value otherwise.

### GetAutoNegotiationOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetAutoNegotiationOk() (*bool, bool)`

GetAutoNegotiationOk returns a tuple with the AutoNegotiation field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAutoNegotiation

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetAutoNegotiation(v bool)`

SetAutoNegotiation sets AutoNegotiation field to given value.

### HasAutoNegotiation

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasAutoNegotiation() bool`

HasAutoNegotiation returns a boolean if a field has been set.

### GetMaxBitRate

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetMaxBitRate() string`

GetMaxBitRate returns the MaxBitRate field if non-nil, zero value otherwise.

### GetMaxBitRateOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetMaxBitRateOk() (*string, bool)`

GetMaxBitRateOk returns a tuple with the MaxBitRate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxBitRate

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetMaxBitRate(v string)`

SetMaxBitRate sets MaxBitRate field to given value.

### HasMaxBitRate

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasMaxBitRate() bool`

HasMaxBitRate returns a boolean if a field has been set.

### GetDuplexMode

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetDuplexMode() string`

GetDuplexMode returns the DuplexMode field if non-nil, zero value otherwise.

### GetDuplexModeOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetDuplexModeOk() (*string, bool)`

GetDuplexModeOk returns a tuple with the DuplexMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDuplexMode

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetDuplexMode(v string)`

SetDuplexMode sets DuplexMode field to given value.

### HasDuplexMode

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasDuplexMode() bool`

HasDuplexMode returns a boolean if a field has been set.

### GetStpEnable

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetStpEnable() bool`

GetStpEnable returns the StpEnable field if non-nil, zero value otherwise.

### GetStpEnableOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetStpEnableOk() (*bool, bool)`

GetStpEnableOk returns a tuple with the StpEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStpEnable

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetStpEnable(v bool)`

SetStpEnable sets StpEnable field to given value.

### HasStpEnable

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasStpEnable() bool`

HasStpEnable returns a boolean if a field has been set.

### GetFastLearningMode

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetFastLearningMode() bool`

GetFastLearningMode returns the FastLearningMode field if non-nil, zero value otherwise.

### GetFastLearningModeOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetFastLearningModeOk() (*bool, bool)`

GetFastLearningModeOk returns a tuple with the FastLearningMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFastLearningMode

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetFastLearningMode(v bool)`

SetFastLearningMode sets FastLearningMode field to given value.

### HasFastLearningMode

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasFastLearningMode() bool`

HasFastLearningMode returns a boolean if a field has been set.

### GetBpduGuard

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetBpduGuard() bool`

GetBpduGuard returns the BpduGuard field if non-nil, zero value otherwise.

### GetBpduGuardOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetBpduGuardOk() (*bool, bool)`

GetBpduGuardOk returns a tuple with the BpduGuard field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBpduGuard

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetBpduGuard(v bool)`

SetBpduGuard sets BpduGuard field to given value.

### HasBpduGuard

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasBpduGuard() bool`

HasBpduGuard returns a boolean if a field has been set.

### GetBpduFilter

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetBpduFilter() bool`

GetBpduFilter returns the BpduFilter field if non-nil, zero value otherwise.

### GetBpduFilterOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetBpduFilterOk() (*bool, bool)`

GetBpduFilterOk returns a tuple with the BpduFilter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBpduFilter

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetBpduFilter(v bool)`

SetBpduFilter sets BpduFilter field to given value.

### HasBpduFilter

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasBpduFilter() bool`

HasBpduFilter returns a boolean if a field has been set.

### GetGuardLoop

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetGuardLoop() bool`

GetGuardLoop returns the GuardLoop field if non-nil, zero value otherwise.

### GetGuardLoopOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetGuardLoopOk() (*bool, bool)`

GetGuardLoopOk returns a tuple with the GuardLoop field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGuardLoop

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetGuardLoop(v bool)`

SetGuardLoop sets GuardLoop field to given value.

### HasGuardLoop

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasGuardLoop() bool`

HasGuardLoop returns a boolean if a field has been set.

### GetPoeEnable

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetPoeEnable() bool`

GetPoeEnable returns the PoeEnable field if non-nil, zero value otherwise.

### GetPoeEnableOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetPoeEnableOk() (*bool, bool)`

GetPoeEnableOk returns a tuple with the PoeEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPoeEnable

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetPoeEnable(v bool)`

SetPoeEnable sets PoeEnable field to given value.

### HasPoeEnable

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasPoeEnable() bool`

HasPoeEnable returns a boolean if a field has been set.

### GetPriority

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetPriority() string`

GetPriority returns the Priority field if non-nil, zero value otherwise.

### GetPriorityOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetPriorityOk() (*string, bool)`

GetPriorityOk returns a tuple with the Priority field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPriority

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetPriority(v string)`

SetPriority sets Priority field to given value.

### HasPriority

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasPriority() bool`

HasPriority returns a boolean if a field has been set.

### GetAllocatedPower

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetAllocatedPower() string`

GetAllocatedPower returns the AllocatedPower field if non-nil, zero value otherwise.

### GetAllocatedPowerOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetAllocatedPowerOk() (*string, bool)`

GetAllocatedPowerOk returns a tuple with the AllocatedPower field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllocatedPower

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetAllocatedPower(v string)`

SetAllocatedPower sets AllocatedPower field to given value.

### HasAllocatedPower

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasAllocatedPower() bool`

HasAllocatedPower returns a boolean if a field has been set.

### GetBspEnable

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetBspEnable() bool`

GetBspEnable returns the BspEnable field if non-nil, zero value otherwise.

### GetBspEnableOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetBspEnableOk() (*bool, bool)`

GetBspEnableOk returns a tuple with the BspEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBspEnable

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetBspEnable(v bool)`

SetBspEnable sets BspEnable field to given value.

### HasBspEnable

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasBspEnable() bool`

HasBspEnable returns a boolean if a field has been set.

### GetBroadcast

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetBroadcast() bool`

GetBroadcast returns the Broadcast field if non-nil, zero value otherwise.

### GetBroadcastOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetBroadcastOk() (*bool, bool)`

GetBroadcastOk returns a tuple with the Broadcast field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBroadcast

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetBroadcast(v bool)`

SetBroadcast sets Broadcast field to given value.

### HasBroadcast

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasBroadcast() bool`

HasBroadcast returns a boolean if a field has been set.

### GetMulticast

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetMulticast() bool`

GetMulticast returns the Multicast field if non-nil, zero value otherwise.

### GetMulticastOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetMulticastOk() (*bool, bool)`

GetMulticastOk returns a tuple with the Multicast field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMulticast

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetMulticast(v bool)`

SetMulticast sets Multicast field to given value.

### HasMulticast

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasMulticast() bool`

HasMulticast returns a boolean if a field has been set.

### GetMaxAllowedValue

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetMaxAllowedValue() int32`

GetMaxAllowedValue returns the MaxAllowedValue field if non-nil, zero value otherwise.

### GetMaxAllowedValueOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetMaxAllowedValueOk() (*int32, bool)`

GetMaxAllowedValueOk returns a tuple with the MaxAllowedValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxAllowedValue

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetMaxAllowedValue(v int32)`

SetMaxAllowedValue sets MaxAllowedValue field to given value.

### HasMaxAllowedValue

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasMaxAllowedValue() bool`

HasMaxAllowedValue returns a boolean if a field has been set.

### GetMaxAllowedUnit

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetMaxAllowedUnit() string`

GetMaxAllowedUnit returns the MaxAllowedUnit field if non-nil, zero value otherwise.

### GetMaxAllowedUnitOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetMaxAllowedUnitOk() (*string, bool)`

GetMaxAllowedUnitOk returns a tuple with the MaxAllowedUnit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxAllowedUnit

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetMaxAllowedUnit(v string)`

SetMaxAllowedUnit sets MaxAllowedUnit field to given value.

### HasMaxAllowedUnit

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasMaxAllowedUnit() bool`

HasMaxAllowedUnit returns a boolean if a field has been set.

### GetAction

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetAction() string`

GetAction returns the Action field if non-nil, zero value otherwise.

### GetActionOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetActionOk() (*string, bool)`

GetActionOk returns a tuple with the Action field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAction

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetAction(v string)`

SetAction sets Action field to given value.

### HasAction

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasAction() bool`

HasAction returns a boolean if a field has been set.

### GetFec

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetFec() string`

GetFec returns the Fec field if non-nil, zero value otherwise.

### GetFecOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetFecOk() (*string, bool)`

GetFecOk returns a tuple with the Fec field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFec

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetFec(v string)`

SetFec sets Fec field to given value.

### HasFec

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasFec() bool`

HasFec returns a boolean if a field has been set.

### GetSingleLink

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetSingleLink() bool`

GetSingleLink returns the SingleLink field if non-nil, zero value otherwise.

### GetSingleLinkOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetSingleLinkOk() (*bool, bool)`

GetSingleLinkOk returns a tuple with the SingleLink field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSingleLink

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetSingleLink(v bool)`

SetSingleLink sets SingleLink field to given value.

### HasSingleLink

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasSingleLink() bool`

HasSingleLink returns a boolean if a field has been set.

### GetObjectProperties

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetObjectProperties() ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetObjectPropertiesOk() (*ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetObjectProperties(v ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


