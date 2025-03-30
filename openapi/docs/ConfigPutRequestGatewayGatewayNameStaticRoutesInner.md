# ConfigPutRequestGatewayGatewayNameStaticRoutesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Enable** | Pointer to **bool** | Enable of this static route | [optional] [default to false]
**Ipv4RoutePrefix** | Pointer to **string** | IPv4 unicast IP address followed by a subnet mask length | [optional] [default to ""]
**NextHopIpAddress** | Pointer to **string** | Next Hop IP Address. Must be a unicast IP address | [optional] [default to ""]
**AdValue** | Pointer to **NullableInt32** | Administrative distancing value, also known as route preference - values from 0-255 | [optional] 
**Index** | Pointer to **int32** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 

## Methods

### NewConfigPutRequestGatewayGatewayNameStaticRoutesInner

`func NewConfigPutRequestGatewayGatewayNameStaticRoutesInner() *ConfigPutRequestGatewayGatewayNameStaticRoutesInner`

NewConfigPutRequestGatewayGatewayNameStaticRoutesInner instantiates a new ConfigPutRequestGatewayGatewayNameStaticRoutesInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestGatewayGatewayNameStaticRoutesInnerWithDefaults

`func NewConfigPutRequestGatewayGatewayNameStaticRoutesInnerWithDefaults() *ConfigPutRequestGatewayGatewayNameStaticRoutesInner`

NewConfigPutRequestGatewayGatewayNameStaticRoutesInnerWithDefaults instantiates a new ConfigPutRequestGatewayGatewayNameStaticRoutesInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEnable

`func (o *ConfigPutRequestGatewayGatewayNameStaticRoutesInner) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestGatewayGatewayNameStaticRoutesInner) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestGatewayGatewayNameStaticRoutesInner) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestGatewayGatewayNameStaticRoutesInner) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetIpv4RoutePrefix

`func (o *ConfigPutRequestGatewayGatewayNameStaticRoutesInner) GetIpv4RoutePrefix() string`

GetIpv4RoutePrefix returns the Ipv4RoutePrefix field if non-nil, zero value otherwise.

### GetIpv4RoutePrefixOk

`func (o *ConfigPutRequestGatewayGatewayNameStaticRoutesInner) GetIpv4RoutePrefixOk() (*string, bool)`

GetIpv4RoutePrefixOk returns a tuple with the Ipv4RoutePrefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpv4RoutePrefix

`func (o *ConfigPutRequestGatewayGatewayNameStaticRoutesInner) SetIpv4RoutePrefix(v string)`

SetIpv4RoutePrefix sets Ipv4RoutePrefix field to given value.

### HasIpv4RoutePrefix

`func (o *ConfigPutRequestGatewayGatewayNameStaticRoutesInner) HasIpv4RoutePrefix() bool`

HasIpv4RoutePrefix returns a boolean if a field has been set.

### GetNextHopIpAddress

`func (o *ConfigPutRequestGatewayGatewayNameStaticRoutesInner) GetNextHopIpAddress() string`

GetNextHopIpAddress returns the NextHopIpAddress field if non-nil, zero value otherwise.

### GetNextHopIpAddressOk

`func (o *ConfigPutRequestGatewayGatewayNameStaticRoutesInner) GetNextHopIpAddressOk() (*string, bool)`

GetNextHopIpAddressOk returns a tuple with the NextHopIpAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNextHopIpAddress

`func (o *ConfigPutRequestGatewayGatewayNameStaticRoutesInner) SetNextHopIpAddress(v string)`

SetNextHopIpAddress sets NextHopIpAddress field to given value.

### HasNextHopIpAddress

`func (o *ConfigPutRequestGatewayGatewayNameStaticRoutesInner) HasNextHopIpAddress() bool`

HasNextHopIpAddress returns a boolean if a field has been set.

### GetAdValue

`func (o *ConfigPutRequestGatewayGatewayNameStaticRoutesInner) GetAdValue() int32`

GetAdValue returns the AdValue field if non-nil, zero value otherwise.

### GetAdValueOk

`func (o *ConfigPutRequestGatewayGatewayNameStaticRoutesInner) GetAdValueOk() (*int32, bool)`

GetAdValueOk returns a tuple with the AdValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAdValue

`func (o *ConfigPutRequestGatewayGatewayNameStaticRoutesInner) SetAdValue(v int32)`

SetAdValue sets AdValue field to given value.

### HasAdValue

`func (o *ConfigPutRequestGatewayGatewayNameStaticRoutesInner) HasAdValue() bool`

HasAdValue returns a boolean if a field has been set.

### SetAdValueNil

`func (o *ConfigPutRequestGatewayGatewayNameStaticRoutesInner) SetAdValueNil(b bool)`

 SetAdValueNil sets the value for AdValue to be an explicit nil

### UnsetAdValue
`func (o *ConfigPutRequestGatewayGatewayNameStaticRoutesInner) UnsetAdValue()`

UnsetAdValue ensures that no value is present for AdValue, not even an explicit nil
### GetIndex

`func (o *ConfigPutRequestGatewayGatewayNameStaticRoutesInner) GetIndex() int32`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *ConfigPutRequestGatewayGatewayNameStaticRoutesInner) GetIndexOk() (*int32, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *ConfigPutRequestGatewayGatewayNameStaticRoutesInner) SetIndex(v int32)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *ConfigPutRequestGatewayGatewayNameStaticRoutesInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


