# SwitchpointsMarkoutofservicePutRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Mos** | Pointer to **bool** | Mark all devices out of service or back in service | [optional] [default to true]
**DeviceNames** | **[]string** |  | 

## Methods

### NewSwitchpointsMarkoutofservicePutRequest

`func NewSwitchpointsMarkoutofservicePutRequest(deviceNames []string, ) *SwitchpointsMarkoutofservicePutRequest`

NewSwitchpointsMarkoutofservicePutRequest instantiates a new SwitchpointsMarkoutofservicePutRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSwitchpointsMarkoutofservicePutRequestWithDefaults

`func NewSwitchpointsMarkoutofservicePutRequestWithDefaults() *SwitchpointsMarkoutofservicePutRequest`

NewSwitchpointsMarkoutofservicePutRequestWithDefaults instantiates a new SwitchpointsMarkoutofservicePutRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMos

`func (o *SwitchpointsMarkoutofservicePutRequest) GetMos() bool`

GetMos returns the Mos field if non-nil, zero value otherwise.

### GetMosOk

`func (o *SwitchpointsMarkoutofservicePutRequest) GetMosOk() (*bool, bool)`

GetMosOk returns a tuple with the Mos field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMos

`func (o *SwitchpointsMarkoutofservicePutRequest) SetMos(v bool)`

SetMos sets Mos field to given value.

### HasMos

`func (o *SwitchpointsMarkoutofservicePutRequest) HasMos() bool`

HasMos returns a boolean if a field has been set.

### GetDeviceNames

`func (o *SwitchpointsMarkoutofservicePutRequest) GetDeviceNames() []string`

GetDeviceNames returns the DeviceNames field if non-nil, zero value otherwise.

### GetDeviceNamesOk

`func (o *SwitchpointsMarkoutofservicePutRequest) GetDeviceNamesOk() (*[]string, bool)`

GetDeviceNamesOk returns a tuple with the DeviceNames field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceNames

`func (o *SwitchpointsMarkoutofservicePutRequest) SetDeviceNames(v []string)`

SetDeviceNames sets DeviceNames field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


