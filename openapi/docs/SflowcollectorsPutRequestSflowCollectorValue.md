# SflowcollectorsPutRequestSflowCollectorValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**Ip** | Pointer to **string** | IP address of the sFlow Collector  | [optional] [default to ""]
**Port** | Pointer to **NullableInt32** | Port | [optional] [default to 6343]

## Methods

### NewSflowcollectorsPutRequestSflowCollectorValue

`func NewSflowcollectorsPutRequestSflowCollectorValue() *SflowcollectorsPutRequestSflowCollectorValue`

NewSflowcollectorsPutRequestSflowCollectorValue instantiates a new SflowcollectorsPutRequestSflowCollectorValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSflowcollectorsPutRequestSflowCollectorValueWithDefaults

`func NewSflowcollectorsPutRequestSflowCollectorValueWithDefaults() *SflowcollectorsPutRequestSflowCollectorValue`

NewSflowcollectorsPutRequestSflowCollectorValueWithDefaults instantiates a new SflowcollectorsPutRequestSflowCollectorValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *SflowcollectorsPutRequestSflowCollectorValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *SflowcollectorsPutRequestSflowCollectorValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *SflowcollectorsPutRequestSflowCollectorValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *SflowcollectorsPutRequestSflowCollectorValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *SflowcollectorsPutRequestSflowCollectorValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *SflowcollectorsPutRequestSflowCollectorValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *SflowcollectorsPutRequestSflowCollectorValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *SflowcollectorsPutRequestSflowCollectorValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetIp

`func (o *SflowcollectorsPutRequestSflowCollectorValue) GetIp() string`

GetIp returns the Ip field if non-nil, zero value otherwise.

### GetIpOk

`func (o *SflowcollectorsPutRequestSflowCollectorValue) GetIpOk() (*string, bool)`

GetIpOk returns a tuple with the Ip field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIp

`func (o *SflowcollectorsPutRequestSflowCollectorValue) SetIp(v string)`

SetIp sets Ip field to given value.

### HasIp

`func (o *SflowcollectorsPutRequestSflowCollectorValue) HasIp() bool`

HasIp returns a boolean if a field has been set.

### GetPort

`func (o *SflowcollectorsPutRequestSflowCollectorValue) GetPort() int32`

GetPort returns the Port field if non-nil, zero value otherwise.

### GetPortOk

`func (o *SflowcollectorsPutRequestSflowCollectorValue) GetPortOk() (*int32, bool)`

GetPortOk returns a tuple with the Port field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPort

`func (o *SflowcollectorsPutRequestSflowCollectorValue) SetPort(v int32)`

SetPort sets Port field to given value.

### HasPort

`func (o *SflowcollectorsPutRequestSflowCollectorValue) HasPort() bool`

HasPort returns a boolean if a field has been set.

### SetPortNil

`func (o *SflowcollectorsPutRequestSflowCollectorValue) SetPortNil(b bool)`

 SetPortNil sets the value for Port to be an explicit nil

### UnsetPort
`func (o *SflowcollectorsPutRequestSflowCollectorValue) UnsetPort()`

UnsetPort ensures that no value is present for Port, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


