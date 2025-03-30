# ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Enable** | Pointer to **bool** | Enable row | [optional] [default to false]
**Gateway** | Pointer to **string** | BGP Gateway referenced for port profile | [optional] [default to ""]
**GatewayRefType** | Pointer to **string** | Object type for gateway field | [optional] 
**SourceIpMask** | Pointer to **string** | Source address on the port if untagged or on the VLAN if tagged used for the outgoing BGP session  | [optional] [default to ""]
**PeerGw** | Pointer to **bool** | Setting for paired switches only. Flag indicating that this gateway is a peer gateway. For each gateway profile referencing a BGP session on a member of a leaf pair, the peer should have a gateway profile entry indicating the IP address for the peers gateway. | [optional] [default to false]
**Index** | Pointer to **int32** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 

## Methods

### NewConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner

`func NewConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner() *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner`

NewConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner instantiates a new ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInnerWithDefaults

`func NewConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInnerWithDefaults() *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner`

NewConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInnerWithDefaults instantiates a new ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEnable

`func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetGateway

`func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) GetGateway() string`

GetGateway returns the Gateway field if non-nil, zero value otherwise.

### GetGatewayOk

`func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) GetGatewayOk() (*string, bool)`

GetGatewayOk returns a tuple with the Gateway field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGateway

`func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) SetGateway(v string)`

SetGateway sets Gateway field to given value.

### HasGateway

`func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) HasGateway() bool`

HasGateway returns a boolean if a field has been set.

### GetGatewayRefType

`func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) GetGatewayRefType() string`

GetGatewayRefType returns the GatewayRefType field if non-nil, zero value otherwise.

### GetGatewayRefTypeOk

`func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) GetGatewayRefTypeOk() (*string, bool)`

GetGatewayRefTypeOk returns a tuple with the GatewayRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGatewayRefType

`func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) SetGatewayRefType(v string)`

SetGatewayRefType sets GatewayRefType field to given value.

### HasGatewayRefType

`func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) HasGatewayRefType() bool`

HasGatewayRefType returns a boolean if a field has been set.

### GetSourceIpMask

`func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) GetSourceIpMask() string`

GetSourceIpMask returns the SourceIpMask field if non-nil, zero value otherwise.

### GetSourceIpMaskOk

`func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) GetSourceIpMaskOk() (*string, bool)`

GetSourceIpMaskOk returns a tuple with the SourceIpMask field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSourceIpMask

`func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) SetSourceIpMask(v string)`

SetSourceIpMask sets SourceIpMask field to given value.

### HasSourceIpMask

`func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) HasSourceIpMask() bool`

HasSourceIpMask returns a boolean if a field has been set.

### GetPeerGw

`func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) GetPeerGw() bool`

GetPeerGw returns the PeerGw field if non-nil, zero value otherwise.

### GetPeerGwOk

`func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) GetPeerGwOk() (*bool, bool)`

GetPeerGwOk returns a tuple with the PeerGw field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPeerGw

`func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) SetPeerGw(v bool)`

SetPeerGw sets PeerGw field to given value.

### HasPeerGw

`func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) HasPeerGw() bool`

HasPeerGw returns a boolean if a field has been set.

### GetIndex

`func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) GetIndex() int32`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) GetIndexOk() (*int32, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) SetIndex(v int32)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


