# ConfigPutRequestEthPortSettingsEthPortSettingsName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. It&#39;s highly recommended to set this value to true so that validation on the object will be ran. | [optional] [default to false]
**AutoNegotiation** | Pointer to **bool** | Indicates if duplex mode should be auto negotiated | [optional] [default to true]
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
**MaxAllowedValue** | Pointer to **int32** | Max Percentage of the ports bandwidth allowed for broadcast/multicast/unknown-unicast traffic before invoking the protective action | [optional] [default to 1000]
**MaxAllowedUnit** | Pointer to **string** | Max Percentage of the ports bandwidth allowed for broadcast/multicast/unknown-unicast traffic before invoking the protective action &lt;br&gt;                                                 &lt;div class&#x3D;\&quot;tab\&quot;&gt;                                                     %: Percentage.&lt;br&gt;                                                     kbps: kilobits per second &lt;br&gt;                                                     mbps: megabits per second &lt;br&gt;                                                     gbps: gigabits per second &lt;br&gt;                                                     pps: packet per second &lt;br&gt;                                                     kpps: kilopacket per second &lt;br&gt;                                                 &lt;/div&gt;                                                  | [optional] [default to "pps"]
**Action** | Pointer to **string** | Action taken if broadcast/multicast/unknown-unicast traffic excedes the Max. One of: &lt;br&gt;                                                 &lt;div class&#x3D;\&quot;tab\&quot;&gt;                                                     Protect: Broadcast/Multicast packets beyond the percent rate are silently dropped. QOS drop counters should indicate the drops.&lt;br&gt;&lt;br&gt;                                                     Restrict: Broadcast/Multicast packets beyond the percent rate are dropped. QOS drop counters should indicate the drops.                                                     Alarm is raised . Alarm automatically clears when rate is below configured threshold. &lt;br&gt;&lt;br&gt;                                                     Shutdown: Alarm is raised and port is taken out of service. User must administratively Disable and Enable the port to restore service. &lt;br&gt;                                                 &lt;/div&gt;                                              | [optional] [default to "Protect"]
**Fec** | Pointer to **string** | FEC is Forward Error Correction which is error correction on the fiber link.                                                 &lt;div class&#x3D;\&quot;tab\&quot;&gt;                                                     Any: Allows switch Negotiation between FC and RS &lt;br&gt;                                                     None: Disables FEC on an interface.&lt;br&gt;                                                     FC: Enables FEC on supported interfaces. FC stands for fire code.&lt;br&gt;                                                     RS: Enables FEC on supported interfaces. RS stands for Reed-Solomon code. &lt;br&gt;                                                     None: VnetC doesn&#39;t alter the Switch Value.&lt;br&gt;                                                 &lt;/div&gt;                                              | [optional] [default to "unaltered"]
**SingleLink** | Pointer to **bool** | Ports with this setting will be disabled when link state tracking takes effect | [optional] [default to false]
**ObjectProperties** | Pointer to [**ConfigPutRequestEthPortSettingsEthPortSettingsNameObjectProperties**](ConfigPutRequestEthPortSettingsEthPortSettingsNameObjectProperties.md) |  | [optional] 
**MinimumWredThreshold** | Pointer to **int32** | A value between 1 to 12480(in KiloBytes) | [optional] [default to 1]
**MaximumWredThreshold** | Pointer to **int32** | A value between 1 to 12480(in KiloBytes) | [optional] [default to 1]
**WredDropProbability** | Pointer to **int32** | A value between 0 to 100 | [optional] [default to 0]
**PriorityFlowControlWatchdogAction** | Pointer to **string** | Ports with this setting will be disabled when link state tracking takes effect | [optional] [default to "DROP"]
**PriorityFlowControlWatchdogDetectTime** | Pointer to **int32** | A value between 100 to 5000 | [optional] [default to 100]
**PriorityFlowControlWatchdogRestoreTime** | Pointer to **int32** | A value between 100 to 60000 | [optional] [default to 100]
**PacketQueue** | Pointer to **string** | Packet Queue | [optional] [default to ""]
**PacketQueueRefType** | Pointer to **string** | Object type for packet_queue field | [optional] 
**EnableWredTuning** | Pointer to **bool** | Enables custom tuning of WRED values. Uncheck to use Switch default values. | [optional] [default to false]
**EnableEcn** | Pointer to **bool** | Enables Explicit Congestion Notification for WRED. | [optional] [default to true]
**EnableWatchdogTuning** | Pointer to **bool** | Enables custom tuning of Watchdog values. Uncheck to use Switch default values. | [optional] [default to false]
**CliCommands** | Pointer to **string** | CLI Commands | [optional] [default to ""]
**DetectBridgingLoops** | Pointer to **bool** | Enable Detection of Bridging Loops | [optional] [default to false]
**UnidirectionalLinkDetection** | Pointer to **bool** | Enable Detection of Unidirectional Link | [optional] [default to false]
**MacSecurityMode** | Pointer to **string** | Dynamic - MACs are learned and aged normally up to the limit. &lt;br&gt;                                 &lt;div class&#x3D;\&quot;tab\&quot;&gt;                                     Packets will be dropped from clients exceeding the limit. &lt;br&gt;                                     Once a client ages out, a new client can take its slot. &lt;br&gt;                                     When the port goes operationally down (disconnecting or disabling), the MACs will be flushed.&lt;br&gt;                                 &lt;/div&gt;                             Sticky - Semi permenant learning. &lt;br&gt;                                 &lt;div class&#x3D;\&quot;tab\&quot;&gt;                                     Packets will be dropped from clients exceeding the limit. &lt;br&gt;                                     Addresses do not age out or move within the same switch. &lt;br&gt;                                     Operationally downing a port (disconnecting) does NOT flush the entries. &lt;br&gt;                                     Learned MACs can only be flushed by administratively taking the port down or rebooting the switch.                                 &lt;/div&gt; | [optional] [default to "disabled"]
**MacLimit** | Pointer to **int32** | Between 1-1000 | [optional] [default to 1000]
**SecurityViolationAction** | Pointer to **string** | Protect - All packets are dropped from clients above the MAC Limit. &lt;br&gt;                                 &lt;div class&#x3D;\&quot;tab\&quot;&gt;                                     Exceeding the limit is not alarmed. &lt;br&gt;                                 &lt;/div&gt;                             Restrict - All packets are dropped from clients above the MAC Limit. &lt;br&gt;                                 &lt;div class&#x3D;\&quot;tab\&quot;&gt;                                     Alarm is raised while attempts to exceed limit are active (MAC has not aged). Alarm automatically clears. &lt;br&gt;                                 &lt;/div&gt;                             Shutdown - Alarm is raised and port is taken down if attempt to exceed MAC limit is made. &lt;br&gt;                                 &lt;div class&#x3D;\&quot;tab\&quot;&gt;                                     User must administratively Disable and Enable the port to restore service.                                 &lt;/div&gt; | [optional] [default to "protect"]
**AgingType** | Pointer to **string** | Limit MAC authentication based on inactivity or on absolute time. See Also Aging Time | [optional] [default to "absolute"]
**AgingTime** | Pointer to **int32** | In minutes, how long the client will stay authenticated. See Also Aging Type | [optional] [default to 0]
**LldpEnable** | Pointer to **bool** | LLDP enable | [optional] [default to true]
**LldpMode** | Pointer to **string** | LLDP mode.  Enables LLDP Rx and/or LLDP Tx | [optional] [default to "RxAndTx"]
**LldpMedEnable** | Pointer to **bool** | LLDP med enable | [optional] [default to false]
**LldpMed** | Pointer to [**[]ConfigPutRequestEthPortSettingsEthPortSettingsNameLldpMedInner**](ConfigPutRequestEthPortSettingsEthPortSettingsNameLldpMedInner.md) |  | [optional] 

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

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetObjectProperties() ConfigPutRequestEthPortSettingsEthPortSettingsNameObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetObjectPropertiesOk() (*ConfigPutRequestEthPortSettingsEthPortSettingsNameObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetObjectProperties(v ConfigPutRequestEthPortSettingsEthPortSettingsNameObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.

### GetMinimumWredThreshold

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetMinimumWredThreshold() int32`

GetMinimumWredThreshold returns the MinimumWredThreshold field if non-nil, zero value otherwise.

### GetMinimumWredThresholdOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetMinimumWredThresholdOk() (*int32, bool)`

GetMinimumWredThresholdOk returns a tuple with the MinimumWredThreshold field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMinimumWredThreshold

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetMinimumWredThreshold(v int32)`

SetMinimumWredThreshold sets MinimumWredThreshold field to given value.

### HasMinimumWredThreshold

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasMinimumWredThreshold() bool`

HasMinimumWredThreshold returns a boolean if a field has been set.

### GetMaximumWredThreshold

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetMaximumWredThreshold() int32`

GetMaximumWredThreshold returns the MaximumWredThreshold field if non-nil, zero value otherwise.

### GetMaximumWredThresholdOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetMaximumWredThresholdOk() (*int32, bool)`

GetMaximumWredThresholdOk returns a tuple with the MaximumWredThreshold field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaximumWredThreshold

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetMaximumWredThreshold(v int32)`

SetMaximumWredThreshold sets MaximumWredThreshold field to given value.

### HasMaximumWredThreshold

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasMaximumWredThreshold() bool`

HasMaximumWredThreshold returns a boolean if a field has been set.

### GetWredDropProbability

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetWredDropProbability() int32`

GetWredDropProbability returns the WredDropProbability field if non-nil, zero value otherwise.

### GetWredDropProbabilityOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetWredDropProbabilityOk() (*int32, bool)`

GetWredDropProbabilityOk returns a tuple with the WredDropProbability field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWredDropProbability

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetWredDropProbability(v int32)`

SetWredDropProbability sets WredDropProbability field to given value.

### HasWredDropProbability

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasWredDropProbability() bool`

HasWredDropProbability returns a boolean if a field has been set.

### GetPriorityFlowControlWatchdogAction

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetPriorityFlowControlWatchdogAction() string`

GetPriorityFlowControlWatchdogAction returns the PriorityFlowControlWatchdogAction field if non-nil, zero value otherwise.

### GetPriorityFlowControlWatchdogActionOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetPriorityFlowControlWatchdogActionOk() (*string, bool)`

GetPriorityFlowControlWatchdogActionOk returns a tuple with the PriorityFlowControlWatchdogAction field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPriorityFlowControlWatchdogAction

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetPriorityFlowControlWatchdogAction(v string)`

SetPriorityFlowControlWatchdogAction sets PriorityFlowControlWatchdogAction field to given value.

### HasPriorityFlowControlWatchdogAction

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasPriorityFlowControlWatchdogAction() bool`

HasPriorityFlowControlWatchdogAction returns a boolean if a field has been set.

### GetPriorityFlowControlWatchdogDetectTime

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetPriorityFlowControlWatchdogDetectTime() int32`

GetPriorityFlowControlWatchdogDetectTime returns the PriorityFlowControlWatchdogDetectTime field if non-nil, zero value otherwise.

### GetPriorityFlowControlWatchdogDetectTimeOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetPriorityFlowControlWatchdogDetectTimeOk() (*int32, bool)`

GetPriorityFlowControlWatchdogDetectTimeOk returns a tuple with the PriorityFlowControlWatchdogDetectTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPriorityFlowControlWatchdogDetectTime

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetPriorityFlowControlWatchdogDetectTime(v int32)`

SetPriorityFlowControlWatchdogDetectTime sets PriorityFlowControlWatchdogDetectTime field to given value.

### HasPriorityFlowControlWatchdogDetectTime

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasPriorityFlowControlWatchdogDetectTime() bool`

HasPriorityFlowControlWatchdogDetectTime returns a boolean if a field has been set.

### GetPriorityFlowControlWatchdogRestoreTime

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetPriorityFlowControlWatchdogRestoreTime() int32`

GetPriorityFlowControlWatchdogRestoreTime returns the PriorityFlowControlWatchdogRestoreTime field if non-nil, zero value otherwise.

### GetPriorityFlowControlWatchdogRestoreTimeOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetPriorityFlowControlWatchdogRestoreTimeOk() (*int32, bool)`

GetPriorityFlowControlWatchdogRestoreTimeOk returns a tuple with the PriorityFlowControlWatchdogRestoreTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPriorityFlowControlWatchdogRestoreTime

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetPriorityFlowControlWatchdogRestoreTime(v int32)`

SetPriorityFlowControlWatchdogRestoreTime sets PriorityFlowControlWatchdogRestoreTime field to given value.

### HasPriorityFlowControlWatchdogRestoreTime

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasPriorityFlowControlWatchdogRestoreTime() bool`

HasPriorityFlowControlWatchdogRestoreTime returns a boolean if a field has been set.

### GetPacketQueue

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetPacketQueue() string`

GetPacketQueue returns the PacketQueue field if non-nil, zero value otherwise.

### GetPacketQueueOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetPacketQueueOk() (*string, bool)`

GetPacketQueueOk returns a tuple with the PacketQueue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPacketQueue

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetPacketQueue(v string)`

SetPacketQueue sets PacketQueue field to given value.

### HasPacketQueue

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasPacketQueue() bool`

HasPacketQueue returns a boolean if a field has been set.

### GetPacketQueueRefType

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetPacketQueueRefType() string`

GetPacketQueueRefType returns the PacketQueueRefType field if non-nil, zero value otherwise.

### GetPacketQueueRefTypeOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetPacketQueueRefTypeOk() (*string, bool)`

GetPacketQueueRefTypeOk returns a tuple with the PacketQueueRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPacketQueueRefType

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetPacketQueueRefType(v string)`

SetPacketQueueRefType sets PacketQueueRefType field to given value.

### HasPacketQueueRefType

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasPacketQueueRefType() bool`

HasPacketQueueRefType returns a boolean if a field has been set.

### GetEnableWredTuning

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetEnableWredTuning() bool`

GetEnableWredTuning returns the EnableWredTuning field if non-nil, zero value otherwise.

### GetEnableWredTuningOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetEnableWredTuningOk() (*bool, bool)`

GetEnableWredTuningOk returns a tuple with the EnableWredTuning field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnableWredTuning

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetEnableWredTuning(v bool)`

SetEnableWredTuning sets EnableWredTuning field to given value.

### HasEnableWredTuning

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasEnableWredTuning() bool`

HasEnableWredTuning returns a boolean if a field has been set.

### GetEnableEcn

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetEnableEcn() bool`

GetEnableEcn returns the EnableEcn field if non-nil, zero value otherwise.

### GetEnableEcnOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetEnableEcnOk() (*bool, bool)`

GetEnableEcnOk returns a tuple with the EnableEcn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnableEcn

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetEnableEcn(v bool)`

SetEnableEcn sets EnableEcn field to given value.

### HasEnableEcn

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasEnableEcn() bool`

HasEnableEcn returns a boolean if a field has been set.

### GetEnableWatchdogTuning

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetEnableWatchdogTuning() bool`

GetEnableWatchdogTuning returns the EnableWatchdogTuning field if non-nil, zero value otherwise.

### GetEnableWatchdogTuningOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetEnableWatchdogTuningOk() (*bool, bool)`

GetEnableWatchdogTuningOk returns a tuple with the EnableWatchdogTuning field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnableWatchdogTuning

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetEnableWatchdogTuning(v bool)`

SetEnableWatchdogTuning sets EnableWatchdogTuning field to given value.

### HasEnableWatchdogTuning

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasEnableWatchdogTuning() bool`

HasEnableWatchdogTuning returns a boolean if a field has been set.

### GetCliCommands

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetCliCommands() string`

GetCliCommands returns the CliCommands field if non-nil, zero value otherwise.

### GetCliCommandsOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetCliCommandsOk() (*string, bool)`

GetCliCommandsOk returns a tuple with the CliCommands field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCliCommands

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetCliCommands(v string)`

SetCliCommands sets CliCommands field to given value.

### HasCliCommands

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasCliCommands() bool`

HasCliCommands returns a boolean if a field has been set.

### GetDetectBridgingLoops

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetDetectBridgingLoops() bool`

GetDetectBridgingLoops returns the DetectBridgingLoops field if non-nil, zero value otherwise.

### GetDetectBridgingLoopsOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetDetectBridgingLoopsOk() (*bool, bool)`

GetDetectBridgingLoopsOk returns a tuple with the DetectBridgingLoops field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDetectBridgingLoops

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetDetectBridgingLoops(v bool)`

SetDetectBridgingLoops sets DetectBridgingLoops field to given value.

### HasDetectBridgingLoops

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasDetectBridgingLoops() bool`

HasDetectBridgingLoops returns a boolean if a field has been set.

### GetUnidirectionalLinkDetection

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetUnidirectionalLinkDetection() bool`

GetUnidirectionalLinkDetection returns the UnidirectionalLinkDetection field if non-nil, zero value otherwise.

### GetUnidirectionalLinkDetectionOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetUnidirectionalLinkDetectionOk() (*bool, bool)`

GetUnidirectionalLinkDetectionOk returns a tuple with the UnidirectionalLinkDetection field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUnidirectionalLinkDetection

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetUnidirectionalLinkDetection(v bool)`

SetUnidirectionalLinkDetection sets UnidirectionalLinkDetection field to given value.

### HasUnidirectionalLinkDetection

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasUnidirectionalLinkDetection() bool`

HasUnidirectionalLinkDetection returns a boolean if a field has been set.

### GetMacSecurityMode

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetMacSecurityMode() string`

GetMacSecurityMode returns the MacSecurityMode field if non-nil, zero value otherwise.

### GetMacSecurityModeOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetMacSecurityModeOk() (*string, bool)`

GetMacSecurityModeOk returns a tuple with the MacSecurityMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMacSecurityMode

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetMacSecurityMode(v string)`

SetMacSecurityMode sets MacSecurityMode field to given value.

### HasMacSecurityMode

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasMacSecurityMode() bool`

HasMacSecurityMode returns a boolean if a field has been set.

### GetMacLimit

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetMacLimit() int32`

GetMacLimit returns the MacLimit field if non-nil, zero value otherwise.

### GetMacLimitOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetMacLimitOk() (*int32, bool)`

GetMacLimitOk returns a tuple with the MacLimit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMacLimit

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetMacLimit(v int32)`

SetMacLimit sets MacLimit field to given value.

### HasMacLimit

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasMacLimit() bool`

HasMacLimit returns a boolean if a field has been set.

### GetSecurityViolationAction

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetSecurityViolationAction() string`

GetSecurityViolationAction returns the SecurityViolationAction field if non-nil, zero value otherwise.

### GetSecurityViolationActionOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetSecurityViolationActionOk() (*string, bool)`

GetSecurityViolationActionOk returns a tuple with the SecurityViolationAction field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSecurityViolationAction

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetSecurityViolationAction(v string)`

SetSecurityViolationAction sets SecurityViolationAction field to given value.

### HasSecurityViolationAction

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasSecurityViolationAction() bool`

HasSecurityViolationAction returns a boolean if a field has been set.

### GetAgingType

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetAgingType() string`

GetAgingType returns the AgingType field if non-nil, zero value otherwise.

### GetAgingTypeOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetAgingTypeOk() (*string, bool)`

GetAgingTypeOk returns a tuple with the AgingType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAgingType

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetAgingType(v string)`

SetAgingType sets AgingType field to given value.

### HasAgingType

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasAgingType() bool`

HasAgingType returns a boolean if a field has been set.

### GetAgingTime

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetAgingTime() int32`

GetAgingTime returns the AgingTime field if non-nil, zero value otherwise.

### GetAgingTimeOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetAgingTimeOk() (*int32, bool)`

GetAgingTimeOk returns a tuple with the AgingTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAgingTime

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetAgingTime(v int32)`

SetAgingTime sets AgingTime field to given value.

### HasAgingTime

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasAgingTime() bool`

HasAgingTime returns a boolean if a field has been set.

### GetLldpEnable

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetLldpEnable() bool`

GetLldpEnable returns the LldpEnable field if non-nil, zero value otherwise.

### GetLldpEnableOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetLldpEnableOk() (*bool, bool)`

GetLldpEnableOk returns a tuple with the LldpEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLldpEnable

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetLldpEnable(v bool)`

SetLldpEnable sets LldpEnable field to given value.

### HasLldpEnable

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasLldpEnable() bool`

HasLldpEnable returns a boolean if a field has been set.

### GetLldpMode

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetLldpMode() string`

GetLldpMode returns the LldpMode field if non-nil, zero value otherwise.

### GetLldpModeOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetLldpModeOk() (*string, bool)`

GetLldpModeOk returns a tuple with the LldpMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLldpMode

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetLldpMode(v string)`

SetLldpMode sets LldpMode field to given value.

### HasLldpMode

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasLldpMode() bool`

HasLldpMode returns a boolean if a field has been set.

### GetLldpMedEnable

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetLldpMedEnable() bool`

GetLldpMedEnable returns the LldpMedEnable field if non-nil, zero value otherwise.

### GetLldpMedEnableOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetLldpMedEnableOk() (*bool, bool)`

GetLldpMedEnableOk returns a tuple with the LldpMedEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLldpMedEnable

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetLldpMedEnable(v bool)`

SetLldpMedEnable sets LldpMedEnable field to given value.

### HasLldpMedEnable

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasLldpMedEnable() bool`

HasLldpMedEnable returns a boolean if a field has been set.

### GetLldpMed

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetLldpMed() []ConfigPutRequestEthPortSettingsEthPortSettingsNameLldpMedInner`

GetLldpMed returns the LldpMed field if non-nil, zero value otherwise.

### GetLldpMedOk

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetLldpMedOk() (*[]ConfigPutRequestEthPortSettingsEthPortSettingsNameLldpMedInner, bool)`

GetLldpMedOk returns a tuple with the LldpMed field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLldpMed

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetLldpMed(v []ConfigPutRequestEthPortSettingsEthPortSettingsNameLldpMedInner)`

SetLldpMed sets LldpMed field to given value.

### HasLldpMed

`func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasLldpMed() bool`

HasLldpMed returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


