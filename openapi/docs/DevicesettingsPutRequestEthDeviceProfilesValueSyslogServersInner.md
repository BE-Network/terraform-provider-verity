# DevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Enabled** | Pointer to **bool** | Enable syslog server | [optional] [default to false]
**Scheme** | Pointer to **string** | Syslog connection scheme | [optional] [default to "udp_bsd"]
**Server** | Pointer to **string** | IPv4, IPv6, or DNS name for syslog server | [optional] [default to ""]
**Port** | Pointer to **string** | Syslog server port | [optional] [default to ""]
**Index** | Pointer to **int32** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 

## Methods

### NewDevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInner

`func NewDevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInner() *DevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInner`

NewDevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInner instantiates a new DevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInnerWithDefaults

`func NewDevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInnerWithDefaults() *DevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInner`

NewDevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInnerWithDefaults instantiates a new DevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEnabled

`func (o *DevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInner) GetEnabled() bool`

GetEnabled returns the Enabled field if non-nil, zero value otherwise.

### GetEnabledOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInner) GetEnabledOk() (*bool, bool)`

GetEnabledOk returns a tuple with the Enabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnabled

`func (o *DevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInner) SetEnabled(v bool)`

SetEnabled sets Enabled field to given value.

### HasEnabled

`func (o *DevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInner) HasEnabled() bool`

HasEnabled returns a boolean if a field has been set.

### GetScheme

`func (o *DevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInner) GetScheme() string`

GetScheme returns the Scheme field if non-nil, zero value otherwise.

### GetSchemeOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInner) GetSchemeOk() (*string, bool)`

GetSchemeOk returns a tuple with the Scheme field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScheme

`func (o *DevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInner) SetScheme(v string)`

SetScheme sets Scheme field to given value.

### HasScheme

`func (o *DevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInner) HasScheme() bool`

HasScheme returns a boolean if a field has been set.

### GetServer

`func (o *DevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInner) GetServer() string`

GetServer returns the Server field if non-nil, zero value otherwise.

### GetServerOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInner) GetServerOk() (*string, bool)`

GetServerOk returns a tuple with the Server field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServer

`func (o *DevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInner) SetServer(v string)`

SetServer sets Server field to given value.

### HasServer

`func (o *DevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInner) HasServer() bool`

HasServer returns a boolean if a field has been set.

### GetPort

`func (o *DevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInner) GetPort() string`

GetPort returns the Port field if non-nil, zero value otherwise.

### GetPortOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInner) GetPortOk() (*string, bool)`

GetPortOk returns a tuple with the Port field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPort

`func (o *DevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInner) SetPort(v string)`

SetPort sets Port field to given value.

### HasPort

`func (o *DevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInner) HasPort() bool`

HasPort returns a boolean if a field has been set.

### GetIndex

`func (o *DevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInner) GetIndex() int32`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInner) GetIndexOk() (*int32, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *DevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInner) SetIndex(v int32)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *DevicesettingsPutRequestEthDeviceProfilesValueSyslogServersInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


