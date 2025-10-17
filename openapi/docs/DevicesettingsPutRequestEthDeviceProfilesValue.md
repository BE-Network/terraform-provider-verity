# DevicesettingsPutRequestEthDeviceProfilesValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**Mode** | Pointer to **string** | Mode | [optional] [default to "IEEE 802.3af"]
**UsageThreshold** | Pointer to **float32** | Usage Threshold | [optional] 
**ExternalBatteryPowerAvailable** | Pointer to **int32** | External Battery Power Available | [optional] [default to 40]
**ExternalPowerAvailable** | Pointer to **int32** | External Power Available | [optional] [default to 75]
**DisableTcpUdpLearnedPacketAcceleration** | Pointer to **bool** | Required for AVB, PTP and Cobranet Support for ONT Devices | [optional] [default to false]
**PacketQueueId** | Pointer to **string** | Packet Queue for device | [optional] [default to "packet_queue|(Packet Queue)|"]
**PacketQueueIdRefType** | Pointer to **string** | Object type for packet_queue_id field | [optional] 
**SecurityAuditInterval** | Pointer to **NullableInt32** | Frequency in minutes of rereading this Switch running configuration and comparing it to expected values.                                                 &lt;br&gt;if the value is blank, audit will use default switch settings.                                                 &lt;br&gt;if the value is 0, audit will be turned off.                                                  | [optional] [default to 60]
**CommitToFlashInterval** | Pointer to **NullableInt32** | Time delay in minutes to write the Switch configuration to flash after a change is made.                                                 &lt;br&gt;if the value is blank, commit will use default switch settings of 12 hours.                                                 &lt;br&gt;if the value is 0, commit will be turned off. | [optional] [default to 60]
**Rocev2** | Pointer to **bool** | Enable RDMA over Converged Ethernet version 2 network protocol. Switches that are set to ROCE mode should already have their port breakouts set up and should not have any ports configured with LAGs. | [optional] [default to false]
**CutThroughSwitching** | Pointer to **bool** | Enable Cut-through Switching on all Switches | [optional] [default to false]
**ObjectProperties** | Pointer to [**DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties**](DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties.md) |  | [optional] 
**HoldTimer** | Pointer to **NullableInt32** | Hold Timer | [optional] [default to 0]
**MacAgingTimerOverride** | Pointer to **NullableInt32** | Blank uses the Device&#39;s default; otherwise an integer between 1 to 1,000,000 seconds | [optional] 
**SpanningTreePriority** | Pointer to **string** | STP per switch, priority are in 4096 increments, the lower the number, the higher the priority. | [optional] [default to "byLevel"]

## Methods

### NewDevicesettingsPutRequestEthDeviceProfilesValue

`func NewDevicesettingsPutRequestEthDeviceProfilesValue() *DevicesettingsPutRequestEthDeviceProfilesValue`

NewDevicesettingsPutRequestEthDeviceProfilesValue instantiates a new DevicesettingsPutRequestEthDeviceProfilesValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDevicesettingsPutRequestEthDeviceProfilesValueWithDefaults

`func NewDevicesettingsPutRequestEthDeviceProfilesValueWithDefaults() *DevicesettingsPutRequestEthDeviceProfilesValue`

NewDevicesettingsPutRequestEthDeviceProfilesValueWithDefaults instantiates a new DevicesettingsPutRequestEthDeviceProfilesValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetMode

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetMode() string`

GetMode returns the Mode field if non-nil, zero value otherwise.

### GetModeOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetModeOk() (*string, bool)`

GetModeOk returns a tuple with the Mode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMode

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetMode(v string)`

SetMode sets Mode field to given value.

### HasMode

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) HasMode() bool`

HasMode returns a boolean if a field has been set.

### GetUsageThreshold

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetUsageThreshold() float32`

GetUsageThreshold returns the UsageThreshold field if non-nil, zero value otherwise.

### GetUsageThresholdOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetUsageThresholdOk() (*float32, bool)`

GetUsageThresholdOk returns a tuple with the UsageThreshold field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsageThreshold

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetUsageThreshold(v float32)`

SetUsageThreshold sets UsageThreshold field to given value.

### HasUsageThreshold

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) HasUsageThreshold() bool`

HasUsageThreshold returns a boolean if a field has been set.

### GetExternalBatteryPowerAvailable

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetExternalBatteryPowerAvailable() int32`

GetExternalBatteryPowerAvailable returns the ExternalBatteryPowerAvailable field if non-nil, zero value otherwise.

### GetExternalBatteryPowerAvailableOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetExternalBatteryPowerAvailableOk() (*int32, bool)`

GetExternalBatteryPowerAvailableOk returns a tuple with the ExternalBatteryPowerAvailable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalBatteryPowerAvailable

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetExternalBatteryPowerAvailable(v int32)`

SetExternalBatteryPowerAvailable sets ExternalBatteryPowerAvailable field to given value.

### HasExternalBatteryPowerAvailable

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) HasExternalBatteryPowerAvailable() bool`

HasExternalBatteryPowerAvailable returns a boolean if a field has been set.

### GetExternalPowerAvailable

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetExternalPowerAvailable() int32`

GetExternalPowerAvailable returns the ExternalPowerAvailable field if non-nil, zero value otherwise.

### GetExternalPowerAvailableOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetExternalPowerAvailableOk() (*int32, bool)`

GetExternalPowerAvailableOk returns a tuple with the ExternalPowerAvailable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalPowerAvailable

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetExternalPowerAvailable(v int32)`

SetExternalPowerAvailable sets ExternalPowerAvailable field to given value.

### HasExternalPowerAvailable

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) HasExternalPowerAvailable() bool`

HasExternalPowerAvailable returns a boolean if a field has been set.

### GetDisableTcpUdpLearnedPacketAcceleration

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetDisableTcpUdpLearnedPacketAcceleration() bool`

GetDisableTcpUdpLearnedPacketAcceleration returns the DisableTcpUdpLearnedPacketAcceleration field if non-nil, zero value otherwise.

### GetDisableTcpUdpLearnedPacketAccelerationOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetDisableTcpUdpLearnedPacketAccelerationOk() (*bool, bool)`

GetDisableTcpUdpLearnedPacketAccelerationOk returns a tuple with the DisableTcpUdpLearnedPacketAcceleration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDisableTcpUdpLearnedPacketAcceleration

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetDisableTcpUdpLearnedPacketAcceleration(v bool)`

SetDisableTcpUdpLearnedPacketAcceleration sets DisableTcpUdpLearnedPacketAcceleration field to given value.

### HasDisableTcpUdpLearnedPacketAcceleration

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) HasDisableTcpUdpLearnedPacketAcceleration() bool`

HasDisableTcpUdpLearnedPacketAcceleration returns a boolean if a field has been set.

### GetPacketQueueId

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetPacketQueueId() string`

GetPacketQueueId returns the PacketQueueId field if non-nil, zero value otherwise.

### GetPacketQueueIdOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetPacketQueueIdOk() (*string, bool)`

GetPacketQueueIdOk returns a tuple with the PacketQueueId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPacketQueueId

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetPacketQueueId(v string)`

SetPacketQueueId sets PacketQueueId field to given value.

### HasPacketQueueId

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) HasPacketQueueId() bool`

HasPacketQueueId returns a boolean if a field has been set.

### GetPacketQueueIdRefType

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetPacketQueueIdRefType() string`

GetPacketQueueIdRefType returns the PacketQueueIdRefType field if non-nil, zero value otherwise.

### GetPacketQueueIdRefTypeOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetPacketQueueIdRefTypeOk() (*string, bool)`

GetPacketQueueIdRefTypeOk returns a tuple with the PacketQueueIdRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPacketQueueIdRefType

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetPacketQueueIdRefType(v string)`

SetPacketQueueIdRefType sets PacketQueueIdRefType field to given value.

### HasPacketQueueIdRefType

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) HasPacketQueueIdRefType() bool`

HasPacketQueueIdRefType returns a boolean if a field has been set.

### GetSecurityAuditInterval

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetSecurityAuditInterval() int32`

GetSecurityAuditInterval returns the SecurityAuditInterval field if non-nil, zero value otherwise.

### GetSecurityAuditIntervalOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetSecurityAuditIntervalOk() (*int32, bool)`

GetSecurityAuditIntervalOk returns a tuple with the SecurityAuditInterval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSecurityAuditInterval

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetSecurityAuditInterval(v int32)`

SetSecurityAuditInterval sets SecurityAuditInterval field to given value.

### HasSecurityAuditInterval

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) HasSecurityAuditInterval() bool`

HasSecurityAuditInterval returns a boolean if a field has been set.

### SetSecurityAuditIntervalNil

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetSecurityAuditIntervalNil(b bool)`

 SetSecurityAuditIntervalNil sets the value for SecurityAuditInterval to be an explicit nil

### UnsetSecurityAuditInterval
`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) UnsetSecurityAuditInterval()`

UnsetSecurityAuditInterval ensures that no value is present for SecurityAuditInterval, not even an explicit nil
### GetCommitToFlashInterval

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetCommitToFlashInterval() int32`

GetCommitToFlashInterval returns the CommitToFlashInterval field if non-nil, zero value otherwise.

### GetCommitToFlashIntervalOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetCommitToFlashIntervalOk() (*int32, bool)`

GetCommitToFlashIntervalOk returns a tuple with the CommitToFlashInterval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommitToFlashInterval

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetCommitToFlashInterval(v int32)`

SetCommitToFlashInterval sets CommitToFlashInterval field to given value.

### HasCommitToFlashInterval

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) HasCommitToFlashInterval() bool`

HasCommitToFlashInterval returns a boolean if a field has been set.

### SetCommitToFlashIntervalNil

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetCommitToFlashIntervalNil(b bool)`

 SetCommitToFlashIntervalNil sets the value for CommitToFlashInterval to be an explicit nil

### UnsetCommitToFlashInterval
`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) UnsetCommitToFlashInterval()`

UnsetCommitToFlashInterval ensures that no value is present for CommitToFlashInterval, not even an explicit nil
### GetRocev2

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetRocev2() bool`

GetRocev2 returns the Rocev2 field if non-nil, zero value otherwise.

### GetRocev2Ok

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetRocev2Ok() (*bool, bool)`

GetRocev2Ok returns a tuple with the Rocev2 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRocev2

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetRocev2(v bool)`

SetRocev2 sets Rocev2 field to given value.

### HasRocev2

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) HasRocev2() bool`

HasRocev2 returns a boolean if a field has been set.

### GetCutThroughSwitching

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetCutThroughSwitching() bool`

GetCutThroughSwitching returns the CutThroughSwitching field if non-nil, zero value otherwise.

### GetCutThroughSwitchingOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetCutThroughSwitchingOk() (*bool, bool)`

GetCutThroughSwitchingOk returns a tuple with the CutThroughSwitching field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCutThroughSwitching

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetCutThroughSwitching(v bool)`

SetCutThroughSwitching sets CutThroughSwitching field to given value.

### HasCutThroughSwitching

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) HasCutThroughSwitching() bool`

HasCutThroughSwitching returns a boolean if a field has been set.

### GetObjectProperties

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetObjectProperties() DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetObjectPropertiesOk() (*DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetObjectProperties(v DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.

### GetHoldTimer

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetHoldTimer() int32`

GetHoldTimer returns the HoldTimer field if non-nil, zero value otherwise.

### GetHoldTimerOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetHoldTimerOk() (*int32, bool)`

GetHoldTimerOk returns a tuple with the HoldTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHoldTimer

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetHoldTimer(v int32)`

SetHoldTimer sets HoldTimer field to given value.

### HasHoldTimer

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) HasHoldTimer() bool`

HasHoldTimer returns a boolean if a field has been set.

### SetHoldTimerNil

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetHoldTimerNil(b bool)`

 SetHoldTimerNil sets the value for HoldTimer to be an explicit nil

### UnsetHoldTimer
`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) UnsetHoldTimer()`

UnsetHoldTimer ensures that no value is present for HoldTimer, not even an explicit nil
### GetMacAgingTimerOverride

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetMacAgingTimerOverride() int32`

GetMacAgingTimerOverride returns the MacAgingTimerOverride field if non-nil, zero value otherwise.

### GetMacAgingTimerOverrideOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetMacAgingTimerOverrideOk() (*int32, bool)`

GetMacAgingTimerOverrideOk returns a tuple with the MacAgingTimerOverride field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMacAgingTimerOverride

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetMacAgingTimerOverride(v int32)`

SetMacAgingTimerOverride sets MacAgingTimerOverride field to given value.

### HasMacAgingTimerOverride

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) HasMacAgingTimerOverride() bool`

HasMacAgingTimerOverride returns a boolean if a field has been set.

### SetMacAgingTimerOverrideNil

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetMacAgingTimerOverrideNil(b bool)`

 SetMacAgingTimerOverrideNil sets the value for MacAgingTimerOverride to be an explicit nil

### UnsetMacAgingTimerOverride
`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) UnsetMacAgingTimerOverride()`

UnsetMacAgingTimerOverride ensures that no value is present for MacAgingTimerOverride, not even an explicit nil
### GetSpanningTreePriority

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetSpanningTreePriority() string`

GetSpanningTreePriority returns the SpanningTreePriority field if non-nil, zero value otherwise.

### GetSpanningTreePriorityOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetSpanningTreePriorityOk() (*string, bool)`

GetSpanningTreePriorityOk returns a tuple with the SpanningTreePriority field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpanningTreePriority

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetSpanningTreePriority(v string)`

SetSpanningTreePriority sets SpanningTreePriority field to given value.

### HasSpanningTreePriority

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) HasSpanningTreePriority() bool`

HasSpanningTreePriority returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


