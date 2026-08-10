# DevicesettingsPutRequestEthDeviceProfilesValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Template Name. Must be unique within type. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**CliCommands** | Pointer to **string** | CLI Commands | [optional] [default to ""]
**Mode** | Pointer to **string** | Mode | [optional] [default to "IEEE 802.3af"]
**UsageThreshold** | Pointer to **NullableFloat64** | Usage Threshold | [optional] 
**ExternalBatteryPowerAvailable** | Pointer to **NullableInt64** | External Battery Power Available | [optional] [default to 40]
**ExternalPowerAvailable** | Pointer to **NullableInt64** | External Power Available | [optional] [default to 75]
**DisableTcpUdpLearnedPacketAcceleration** | Pointer to **bool** | Required for AVB, PTP and Cobranet Support for ONT Devices | [optional] [default to false]
**PacketQueue** | Pointer to **string** | Packet Queue for device | [optional] [default to ""]
**PacketQueueRefType** | Pointer to **string** | Object type for packet_queue field | [optional] 
**DeviceAaaProfile** | Pointer to **string** | Device AAA Profile for authentication settings | [optional] [default to ""]
**DeviceAaaProfileRefType** | Pointer to **string** | Object type for device_aaa_profile field | [optional] 
**SecurityAuditInterval** | Pointer to **NullableInt64** | Frequency in minutes of rereading this Switch running configuration and comparing it to expected values.                                                 &lt;br&gt;if the value is blank, audit will use default switch settings.                                                 &lt;br&gt;if the value is 0, audit will be turned off.                                                  | [optional] [default to 60]
**CommitToFlashInterval** | Pointer to **NullableInt64** | Time delay in minutes to write the Switch configuration to flash after a change is made.                                                 &lt;br&gt;if the value is blank, commit will use default switch settings of 12 hours.                                                 &lt;br&gt;if the value is 0, commit will be turned off. | [optional] [default to 60]
**Rocev2** | Pointer to **bool** | Enable RDMA over Converged Ethernet version 2 network protocol. Switches that are set to ROCE mode should already have their port breakouts set up and should not have any ports configured with LAGs. | [optional] [default to false]
**CutThroughSwitching** | Pointer to **bool** | Enable Cut-through Switching on all Switches | [optional] [default to false]
**LoginBanner** | Pointer to **string** | Banner message displayed at login | [optional] [default to ""]
**DnsServers** | Pointer to [**[]DevicesettingsPutRequestEthDeviceProfilesValueDnsServersInner**](DevicesettingsPutRequestEthDeviceProfilesValueDnsServersInner.md) |  | [optional] 
**NtpServers** | Pointer to [**[]DevicesettingsPutRequestEthDeviceProfilesValueNtpServersInner**](DevicesettingsPutRequestEthDeviceProfilesValueNtpServersInner.md) |  | [optional] 
**SyslogServers** | Pointer to [**[]DevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInner**](DevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInner.md) |  | [optional] 
**ObjectProperties** | Pointer to **map[string]interface{}** |  | [optional] 
**HoldTimer** | Pointer to **NullableInt64** | Hold Timer | [optional] [default to 0]
**MacAgingTimerOverride** | Pointer to **NullableInt64** | Blank uses the Device&#39;s default; otherwise an integer between 1 to 1,000,000 seconds | [optional] 
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

### GetCliCommands

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetCliCommands() string`

GetCliCommands returns the CliCommands field if non-nil, zero value otherwise.

### GetCliCommandsOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetCliCommandsOk() (*string, bool)`

GetCliCommandsOk returns a tuple with the CliCommands field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCliCommands

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetCliCommands(v string)`

SetCliCommands sets CliCommands field to given value.

### HasCliCommands

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) HasCliCommands() bool`

HasCliCommands returns a boolean if a field has been set.

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

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetUsageThreshold() float64`

GetUsageThreshold returns the UsageThreshold field if non-nil, zero value otherwise.

### GetUsageThresholdOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetUsageThresholdOk() (*float64, bool)`

GetUsageThresholdOk returns a tuple with the UsageThreshold field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsageThreshold

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetUsageThreshold(v float64)`

SetUsageThreshold sets UsageThreshold field to given value.

### HasUsageThreshold

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) HasUsageThreshold() bool`

HasUsageThreshold returns a boolean if a field has been set.

### SetUsageThresholdNil

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetUsageThresholdNil(b bool)`

 SetUsageThresholdNil sets the value for UsageThreshold to be an explicit nil

### UnsetUsageThreshold
`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) UnsetUsageThreshold()`

UnsetUsageThreshold ensures that no value is present for UsageThreshold, not even an explicit nil
### GetExternalBatteryPowerAvailable

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetExternalBatteryPowerAvailable() int64`

GetExternalBatteryPowerAvailable returns the ExternalBatteryPowerAvailable field if non-nil, zero value otherwise.

### GetExternalBatteryPowerAvailableOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetExternalBatteryPowerAvailableOk() (*int64, bool)`

GetExternalBatteryPowerAvailableOk returns a tuple with the ExternalBatteryPowerAvailable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalBatteryPowerAvailable

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetExternalBatteryPowerAvailable(v int64)`

SetExternalBatteryPowerAvailable sets ExternalBatteryPowerAvailable field to given value.

### HasExternalBatteryPowerAvailable

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) HasExternalBatteryPowerAvailable() bool`

HasExternalBatteryPowerAvailable returns a boolean if a field has been set.

### SetExternalBatteryPowerAvailableNil

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetExternalBatteryPowerAvailableNil(b bool)`

 SetExternalBatteryPowerAvailableNil sets the value for ExternalBatteryPowerAvailable to be an explicit nil

### UnsetExternalBatteryPowerAvailable
`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) UnsetExternalBatteryPowerAvailable()`

UnsetExternalBatteryPowerAvailable ensures that no value is present for ExternalBatteryPowerAvailable, not even an explicit nil
### GetExternalPowerAvailable

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetExternalPowerAvailable() int64`

GetExternalPowerAvailable returns the ExternalPowerAvailable field if non-nil, zero value otherwise.

### GetExternalPowerAvailableOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetExternalPowerAvailableOk() (*int64, bool)`

GetExternalPowerAvailableOk returns a tuple with the ExternalPowerAvailable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalPowerAvailable

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetExternalPowerAvailable(v int64)`

SetExternalPowerAvailable sets ExternalPowerAvailable field to given value.

### HasExternalPowerAvailable

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) HasExternalPowerAvailable() bool`

HasExternalPowerAvailable returns a boolean if a field has been set.

### SetExternalPowerAvailableNil

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetExternalPowerAvailableNil(b bool)`

 SetExternalPowerAvailableNil sets the value for ExternalPowerAvailable to be an explicit nil

### UnsetExternalPowerAvailable
`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) UnsetExternalPowerAvailable()`

UnsetExternalPowerAvailable ensures that no value is present for ExternalPowerAvailable, not even an explicit nil
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

### GetPacketQueue

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetPacketQueue() string`

GetPacketQueue returns the PacketQueue field if non-nil, zero value otherwise.

### GetPacketQueueOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetPacketQueueOk() (*string, bool)`

GetPacketQueueOk returns a tuple with the PacketQueue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPacketQueue

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetPacketQueue(v string)`

SetPacketQueue sets PacketQueue field to given value.

### HasPacketQueue

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) HasPacketQueue() bool`

HasPacketQueue returns a boolean if a field has been set.

### GetPacketQueueRefType

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetPacketQueueRefType() string`

GetPacketQueueRefType returns the PacketQueueRefType field if non-nil, zero value otherwise.

### GetPacketQueueRefTypeOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetPacketQueueRefTypeOk() (*string, bool)`

GetPacketQueueRefTypeOk returns a tuple with the PacketQueueRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPacketQueueRefType

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetPacketQueueRefType(v string)`

SetPacketQueueRefType sets PacketQueueRefType field to given value.

### HasPacketQueueRefType

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) HasPacketQueueRefType() bool`

HasPacketQueueRefType returns a boolean if a field has been set.

### GetDeviceAaaProfile

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetDeviceAaaProfile() string`

GetDeviceAaaProfile returns the DeviceAaaProfile field if non-nil, zero value otherwise.

### GetDeviceAaaProfileOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetDeviceAaaProfileOk() (*string, bool)`

GetDeviceAaaProfileOk returns a tuple with the DeviceAaaProfile field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceAaaProfile

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetDeviceAaaProfile(v string)`

SetDeviceAaaProfile sets DeviceAaaProfile field to given value.

### HasDeviceAaaProfile

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) HasDeviceAaaProfile() bool`

HasDeviceAaaProfile returns a boolean if a field has been set.

### GetDeviceAaaProfileRefType

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetDeviceAaaProfileRefType() string`

GetDeviceAaaProfileRefType returns the DeviceAaaProfileRefType field if non-nil, zero value otherwise.

### GetDeviceAaaProfileRefTypeOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetDeviceAaaProfileRefTypeOk() (*string, bool)`

GetDeviceAaaProfileRefTypeOk returns a tuple with the DeviceAaaProfileRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceAaaProfileRefType

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetDeviceAaaProfileRefType(v string)`

SetDeviceAaaProfileRefType sets DeviceAaaProfileRefType field to given value.

### HasDeviceAaaProfileRefType

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) HasDeviceAaaProfileRefType() bool`

HasDeviceAaaProfileRefType returns a boolean if a field has been set.

### GetSecurityAuditInterval

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetSecurityAuditInterval() int64`

GetSecurityAuditInterval returns the SecurityAuditInterval field if non-nil, zero value otherwise.

### GetSecurityAuditIntervalOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetSecurityAuditIntervalOk() (*int64, bool)`

GetSecurityAuditIntervalOk returns a tuple with the SecurityAuditInterval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSecurityAuditInterval

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetSecurityAuditInterval(v int64)`

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

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetCommitToFlashInterval() int64`

GetCommitToFlashInterval returns the CommitToFlashInterval field if non-nil, zero value otherwise.

### GetCommitToFlashIntervalOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetCommitToFlashIntervalOk() (*int64, bool)`

GetCommitToFlashIntervalOk returns a tuple with the CommitToFlashInterval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommitToFlashInterval

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetCommitToFlashInterval(v int64)`

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

### GetLoginBanner

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetLoginBanner() string`

GetLoginBanner returns the LoginBanner field if non-nil, zero value otherwise.

### GetLoginBannerOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetLoginBannerOk() (*string, bool)`

GetLoginBannerOk returns a tuple with the LoginBanner field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLoginBanner

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetLoginBanner(v string)`

SetLoginBanner sets LoginBanner field to given value.

### HasLoginBanner

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) HasLoginBanner() bool`

HasLoginBanner returns a boolean if a field has been set.

### GetDnsServers

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetDnsServers() []DevicesettingsPutRequestEthDeviceProfilesValueDnsServersInner`

GetDnsServers returns the DnsServers field if non-nil, zero value otherwise.

### GetDnsServersOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetDnsServersOk() (*[]DevicesettingsPutRequestEthDeviceProfilesValueDnsServersInner, bool)`

GetDnsServersOk returns a tuple with the DnsServers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDnsServers

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetDnsServers(v []DevicesettingsPutRequestEthDeviceProfilesValueDnsServersInner)`

SetDnsServers sets DnsServers field to given value.

### HasDnsServers

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) HasDnsServers() bool`

HasDnsServers returns a boolean if a field has been set.

### GetNtpServers

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetNtpServers() []DevicesettingsPutRequestEthDeviceProfilesValueNtpServersInner`

GetNtpServers returns the NtpServers field if non-nil, zero value otherwise.

### GetNtpServersOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetNtpServersOk() (*[]DevicesettingsPutRequestEthDeviceProfilesValueNtpServersInner, bool)`

GetNtpServersOk returns a tuple with the NtpServers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNtpServers

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetNtpServers(v []DevicesettingsPutRequestEthDeviceProfilesValueNtpServersInner)`

SetNtpServers sets NtpServers field to given value.

### HasNtpServers

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) HasNtpServers() bool`

HasNtpServers returns a boolean if a field has been set.

### GetSyslogServers

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetSyslogServers() []DevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInner`

GetSyslogServers returns the SyslogServers field if non-nil, zero value otherwise.

### GetSyslogServersOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetSyslogServersOk() (*[]DevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInner, bool)`

GetSyslogServersOk returns a tuple with the SyslogServers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSyslogServers

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetSyslogServers(v []DevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInner)`

SetSyslogServers sets SyslogServers field to given value.

### HasSyslogServers

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) HasSyslogServers() bool`

HasSyslogServers returns a boolean if a field has been set.

### GetObjectProperties

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetObjectProperties() map[string]interface{}`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetObjectPropertiesOk() (*map[string]interface{}, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetObjectProperties(v map[string]interface{})`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.

### GetHoldTimer

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetHoldTimer() int64`

GetHoldTimer returns the HoldTimer field if non-nil, zero value otherwise.

### GetHoldTimerOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetHoldTimerOk() (*int64, bool)`

GetHoldTimerOk returns a tuple with the HoldTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHoldTimer

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetHoldTimer(v int64)`

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

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetMacAgingTimerOverride() int64`

GetMacAgingTimerOverride returns the MacAgingTimerOverride field if non-nil, zero value otherwise.

### GetMacAgingTimerOverrideOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) GetMacAgingTimerOverrideOk() (*int64, bool)`

GetMacAgingTimerOverrideOk returns a tuple with the MacAgingTimerOverride field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMacAgingTimerOverride

`func (o *DevicesettingsPutRequestEthDeviceProfilesValue) SetMacAgingTimerOverride(v int64)`

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


