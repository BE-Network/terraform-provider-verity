# GatewaysPutRequestGatewayValueStaticRoutesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Enable** | Pointer to **bool** | Enable of this static route | [optional] [default to false]
**Ipv4RoutePrefix** | Pointer to **string** | IPv4 unicast IP address followed by a subnet mask length | [optional] [default to ""]
**NextHopIpAddress** | Pointer to **string** | Next Hop IP Address. Must be a unicast IP address | [optional] [default to ""]
**AdValue** | Pointer to **NullableInt32** | Administrative distancing value, also known as route preference - values from 0-255 | [optional] 
**Index** | Pointer to **int32** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 

## Methods

### NewGatewaysPutRequestGatewayValueStaticRoutesInner

`func NewGatewaysPutRequestGatewayValueStaticRoutesInner() *GatewaysPutRequestGatewayValueStaticRoutesInner`

NewGatewaysPutRequestGatewayValueStaticRoutesInner instantiates a new GatewaysPutRequestGatewayValueStaticRoutesInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGatewaysPutRequestGatewayValueStaticRoutesInnerWithDefaults

`func NewGatewaysPutRequestGatewayValueStaticRoutesInnerWithDefaults() *GatewaysPutRequestGatewayValueStaticRoutesInner`

NewGatewaysPutRequestGatewayValueStaticRoutesInnerWithDefaults instantiates a new GatewaysPutRequestGatewayValueStaticRoutesInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEnable

`func (o *GatewaysPutRequestGatewayValueStaticRoutesInner) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *GatewaysPutRequestGatewayValueStaticRoutesInner) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *GatewaysPutRequestGatewayValueStaticRoutesInner) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *GatewaysPutRequestGatewayValueStaticRoutesInner) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetIpv4RoutePrefix

`func (o *GatewaysPutRequestGatewayValueStaticRoutesInner) GetIpv4RoutePrefix() string`

GetIpv4RoutePrefix returns the Ipv4RoutePrefix field if non-nil, zero value otherwise.

### GetIpv4RoutePrefixOk

`func (o *GatewaysPutRequestGatewayValueStaticRoutesInner) GetIpv4RoutePrefixOk() (*string, bool)`

GetIpv4RoutePrefixOk returns a tuple with the Ipv4RoutePrefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpv4RoutePrefix

`func (o *GatewaysPutRequestGatewayValueStaticRoutesInner) SetIpv4RoutePrefix(v string)`

SetIpv4RoutePrefix sets Ipv4RoutePrefix field to given value.

### HasIpv4RoutePrefix

`func (o *GatewaysPutRequestGatewayValueStaticRoutesInner) HasIpv4RoutePrefix() bool`

HasIpv4RoutePrefix returns a boolean if a field has been set.

### GetNextHopIpAddress

`func (o *GatewaysPutRequestGatewayValueStaticRoutesInner) GetNextHopIpAddress() string`

GetNextHopIpAddress returns the NextHopIpAddress field if non-nil, zero value otherwise.

### GetNextHopIpAddressOk

`func (o *GatewaysPutRequestGatewayValueStaticRoutesInner) GetNextHopIpAddressOk() (*string, bool)`

GetNextHopIpAddressOk returns a tuple with the NextHopIpAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNextHopIpAddress

`func (o *GatewaysPutRequestGatewayValueStaticRoutesInner) SetNextHopIpAddress(v string)`

SetNextHopIpAddress sets NextHopIpAddress field to given value.

### HasNextHopIpAddress

`func (o *GatewaysPutRequestGatewayValueStaticRoutesInner) HasNextHopIpAddress() bool`

HasNextHopIpAddress returns a boolean if a field has been set.

### GetAdValue

`func (o *GatewaysPutRequestGatewayValueStaticRoutesInner) GetAdValue() int32`

GetAdValue returns the AdValue field if non-nil, zero value otherwise.

### GetAdValueOk

`func (o *GatewaysPutRequestGatewayValueStaticRoutesInner) GetAdValueOk() (*int32, bool)`

GetAdValueOk returns a tuple with the AdValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAdValue

`func (o *GatewaysPutRequestGatewayValueStaticRoutesInner) SetAdValue(v int32)`

SetAdValue sets AdValue field to given value.

### HasAdValue

`func (o *GatewaysPutRequestGatewayValueStaticRoutesInner) HasAdValue() bool`

HasAdValue returns a boolean if a field has been set.

### SetAdValueNil

`func (o *GatewaysPutRequestGatewayValueStaticRoutesInner) SetAdValueNil(b bool)`

 SetAdValueNil sets the value for AdValue to be an explicit nil

### UnsetAdValue
`func (o *GatewaysPutRequestGatewayValueStaticRoutesInner) UnsetAdValue()`

UnsetAdValue ensures that no value is present for AdValue, not even an explicit nil
### GetIndex

`func (o *GatewaysPutRequestGatewayValueStaticRoutesInner) GetIndex() int32`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *GatewaysPutRequestGatewayValueStaticRoutesInner) GetIndexOk() (*int32, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *GatewaysPutRequestGatewayValueStaticRoutesInner) SetIndex(v int32)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *GatewaysPutRequestGatewayValueStaticRoutesInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


