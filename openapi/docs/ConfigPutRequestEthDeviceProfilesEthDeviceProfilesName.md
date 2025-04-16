# ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**Mode** | Pointer to **string** | Mode | [optional] [default to "IEEE 802.3af"]
**UsageThreshold** | Pointer to **float32** | Usage Threshold | [optional] 
**ExternalBatteryPowerAvailable** | Pointer to **int32** | External Battery Power Available | [optional] [default to 40]
**ExternalPowerAvailable** | Pointer to **int32** | External Power Available | [optional] [default to 75]
**SecurityAuditInterval** | Pointer to **NullableInt32** | Frequency in minutes of rereading this Switch running configuration and comparing it to expected values.                                                 &lt;br&gt;if the value is blank, audit will use default switch settings.                                                 &lt;br&gt;if the value is 0, audit will be turned off.                                                  | [optional] [default to 60]
**CommitToFlashInterval** | Pointer to **NullableInt32** | Frequency in minutes to write the Switch configuration to flash.                                                 &lt;br&gt;if the value is blank, commit will use default switch settings.                                                 &lt;br&gt;if the value is 0, commit will be turned off. | [optional] [default to 60]
**Rocev2** | Pointer to **bool** | Enable RDMA over Converged Ethernet version 2 network protocol. Switches that are set to ROCE mode should already have their port breakouts set up and should not have any ports configured with LAGs. | [optional] [default to false]
**CutThroughSwitching** | Pointer to **bool** | Enable Cut-through Switching on all Switches | [optional] [default to false]
**ObjectProperties** | Pointer to [**ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties**](ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties.md) |  | [optional] 

## Methods

### NewConfigPutRequestEthDeviceProfilesEthDeviceProfilesName

`func NewConfigPutRequestEthDeviceProfilesEthDeviceProfilesName() *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName`

NewConfigPutRequestEthDeviceProfilesEthDeviceProfilesName instantiates a new ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameWithDefaults

`func NewConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameWithDefaults() *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName`

NewConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameWithDefaults instantiates a new ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetMode

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) GetMode() string`

GetMode returns the Mode field if non-nil, zero value otherwise.

### GetModeOk

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) GetModeOk() (*string, bool)`

GetModeOk returns a tuple with the Mode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMode

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) SetMode(v string)`

SetMode sets Mode field to given value.

### HasMode

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) HasMode() bool`

HasMode returns a boolean if a field has been set.

### GetUsageThreshold

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) GetUsageThreshold() float32`

GetUsageThreshold returns the UsageThreshold field if non-nil, zero value otherwise.

### GetUsageThresholdOk

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) GetUsageThresholdOk() (*float32, bool)`

GetUsageThresholdOk returns a tuple with the UsageThreshold field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsageThreshold

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) SetUsageThreshold(v float32)`

SetUsageThreshold sets UsageThreshold field to given value.

### HasUsageThreshold

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) HasUsageThreshold() bool`

HasUsageThreshold returns a boolean if a field has been set.

### GetExternalBatteryPowerAvailable

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) GetExternalBatteryPowerAvailable() int32`

GetExternalBatteryPowerAvailable returns the ExternalBatteryPowerAvailable field if non-nil, zero value otherwise.

### GetExternalBatteryPowerAvailableOk

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) GetExternalBatteryPowerAvailableOk() (*int32, bool)`

GetExternalBatteryPowerAvailableOk returns a tuple with the ExternalBatteryPowerAvailable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalBatteryPowerAvailable

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) SetExternalBatteryPowerAvailable(v int32)`

SetExternalBatteryPowerAvailable sets ExternalBatteryPowerAvailable field to given value.

### HasExternalBatteryPowerAvailable

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) HasExternalBatteryPowerAvailable() bool`

HasExternalBatteryPowerAvailable returns a boolean if a field has been set.

### GetExternalPowerAvailable

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) GetExternalPowerAvailable() int32`

GetExternalPowerAvailable returns the ExternalPowerAvailable field if non-nil, zero value otherwise.

### GetExternalPowerAvailableOk

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) GetExternalPowerAvailableOk() (*int32, bool)`

GetExternalPowerAvailableOk returns a tuple with the ExternalPowerAvailable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExternalPowerAvailable

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) SetExternalPowerAvailable(v int32)`

SetExternalPowerAvailable sets ExternalPowerAvailable field to given value.

### HasExternalPowerAvailable

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) HasExternalPowerAvailable() bool`

HasExternalPowerAvailable returns a boolean if a field has been set.

### GetSecurityAuditInterval

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) GetSecurityAuditInterval() int32`

GetSecurityAuditInterval returns the SecurityAuditInterval field if non-nil, zero value otherwise.

### GetSecurityAuditIntervalOk

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) GetSecurityAuditIntervalOk() (*int32, bool)`

GetSecurityAuditIntervalOk returns a tuple with the SecurityAuditInterval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSecurityAuditInterval

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) SetSecurityAuditInterval(v int32)`

SetSecurityAuditInterval sets SecurityAuditInterval field to given value.

### HasSecurityAuditInterval

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) HasSecurityAuditInterval() bool`

HasSecurityAuditInterval returns a boolean if a field has been set.

### SetSecurityAuditIntervalNil

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) SetSecurityAuditIntervalNil(b bool)`

 SetSecurityAuditIntervalNil sets the value for SecurityAuditInterval to be an explicit nil

### UnsetSecurityAuditInterval
`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) UnsetSecurityAuditInterval()`

UnsetSecurityAuditInterval ensures that no value is present for SecurityAuditInterval, not even an explicit nil
### GetCommitToFlashInterval

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) GetCommitToFlashInterval() int32`

GetCommitToFlashInterval returns the CommitToFlashInterval field if non-nil, zero value otherwise.

### GetCommitToFlashIntervalOk

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) GetCommitToFlashIntervalOk() (*int32, bool)`

GetCommitToFlashIntervalOk returns a tuple with the CommitToFlashInterval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommitToFlashInterval

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) SetCommitToFlashInterval(v int32)`

SetCommitToFlashInterval sets CommitToFlashInterval field to given value.

### HasCommitToFlashInterval

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) HasCommitToFlashInterval() bool`

HasCommitToFlashInterval returns a boolean if a field has been set.

### SetCommitToFlashIntervalNil

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) SetCommitToFlashIntervalNil(b bool)`

 SetCommitToFlashIntervalNil sets the value for CommitToFlashInterval to be an explicit nil

### UnsetCommitToFlashInterval
`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) UnsetCommitToFlashInterval()`

UnsetCommitToFlashInterval ensures that no value is present for CommitToFlashInterval, not even an explicit nil
### GetRocev2

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) GetRocev2() bool`

GetRocev2 returns the Rocev2 field if non-nil, zero value otherwise.

### GetRocev2Ok

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) GetRocev2Ok() (*bool, bool)`

GetRocev2Ok returns a tuple with the Rocev2 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRocev2

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) SetRocev2(v bool)`

SetRocev2 sets Rocev2 field to given value.

### HasRocev2

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) HasRocev2() bool`

HasRocev2 returns a boolean if a field has been set.

### GetCutThroughSwitching

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) GetCutThroughSwitching() bool`

GetCutThroughSwitching returns the CutThroughSwitching field if non-nil, zero value otherwise.

### GetCutThroughSwitchingOk

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) GetCutThroughSwitchingOk() (*bool, bool)`

GetCutThroughSwitchingOk returns a tuple with the CutThroughSwitching field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCutThroughSwitching

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) SetCutThroughSwitching(v bool)`

SetCutThroughSwitching sets CutThroughSwitching field to given value.

### HasCutThroughSwitching

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) HasCutThroughSwitching() bool`

HasCutThroughSwitching returns a boolean if a field has been set.

### GetObjectProperties

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) GetObjectProperties() ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) GetObjectPropertiesOk() (*ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) SetObjectProperties(v ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesName) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


