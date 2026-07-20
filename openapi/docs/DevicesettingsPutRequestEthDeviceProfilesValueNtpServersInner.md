# DevicesettingsPutRequestEthDeviceProfilesValueNtpServersInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Enabled** | Pointer to **bool** | Enable NTP server | [optional] [default to false]
**Server** | Pointer to **string** | IPv4, IPv6, or DNS name for NTP server | [optional] [default to ""]
**Index** | Pointer to **int32** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 

## Methods

### NewDevicesettingsPutRequestEthDeviceProfilesValueNtpServersInner

`func NewDevicesettingsPutRequestEthDeviceProfilesValueNtpServersInner() *DevicesettingsPutRequestEthDeviceProfilesValueNtpServersInner`

NewDevicesettingsPutRequestEthDeviceProfilesValueNtpServersInner instantiates a new DevicesettingsPutRequestEthDeviceProfilesValueNtpServersInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDevicesettingsPutRequestEthDeviceProfilesValueNtpServersInnerWithDefaults

`func NewDevicesettingsPutRequestEthDeviceProfilesValueNtpServersInnerWithDefaults() *DevicesettingsPutRequestEthDeviceProfilesValueNtpServersInner`

NewDevicesettingsPutRequestEthDeviceProfilesValueNtpServersInnerWithDefaults instantiates a new DevicesettingsPutRequestEthDeviceProfilesValueNtpServersInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEnabled

`func (o *DevicesettingsPutRequestEthDeviceProfilesValueNtpServersInner) GetEnabled() bool`

GetEnabled returns the Enabled field if non-nil, zero value otherwise.

### GetEnabledOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValueNtpServersInner) GetEnabledOk() (*bool, bool)`

GetEnabledOk returns a tuple with the Enabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnabled

`func (o *DevicesettingsPutRequestEthDeviceProfilesValueNtpServersInner) SetEnabled(v bool)`

SetEnabled sets Enabled field to given value.

### HasEnabled

`func (o *DevicesettingsPutRequestEthDeviceProfilesValueNtpServersInner) HasEnabled() bool`

HasEnabled returns a boolean if a field has been set.

### GetServer

`func (o *DevicesettingsPutRequestEthDeviceProfilesValueNtpServersInner) GetServer() string`

GetServer returns the Server field if non-nil, zero value otherwise.

### GetServerOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValueNtpServersInner) GetServerOk() (*string, bool)`

GetServerOk returns a tuple with the Server field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServer

`func (o *DevicesettingsPutRequestEthDeviceProfilesValueNtpServersInner) SetServer(v string)`

SetServer sets Server field to given value.

### HasServer

`func (o *DevicesettingsPutRequestEthDeviceProfilesValueNtpServersInner) HasServer() bool`

HasServer returns a boolean if a field has been set.

### GetIndex

`func (o *DevicesettingsPutRequestEthDeviceProfilesValueNtpServersInner) GetIndex() int32`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *DevicesettingsPutRequestEthDeviceProfilesValueNtpServersInner) GetIndexOk() (*int32, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *DevicesettingsPutRequestEthDeviceProfilesValueNtpServersInner) SetIndex(v int32)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *DevicesettingsPutRequestEthDeviceProfilesValueNtpServersInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


