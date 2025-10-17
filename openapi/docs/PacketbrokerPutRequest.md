# PacketbrokerPutRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**PortAcl** | Pointer to [**map[string]PacketbrokerPutRequestPortAclValue**](PacketbrokerPutRequestPortAclValue.md) |  | [optional] 

## Methods

### NewPacketbrokerPutRequest

`func NewPacketbrokerPutRequest() *PacketbrokerPutRequest`

NewPacketbrokerPutRequest instantiates a new PacketbrokerPutRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPacketbrokerPutRequestWithDefaults

`func NewPacketbrokerPutRequestWithDefaults() *PacketbrokerPutRequest`

NewPacketbrokerPutRequestWithDefaults instantiates a new PacketbrokerPutRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPortAcl

`func (o *PacketbrokerPutRequest) GetPortAcl() map[string]PacketbrokerPutRequestPortAclValue`

GetPortAcl returns the PortAcl field if non-nil, zero value otherwise.

### GetPortAclOk

`func (o *PacketbrokerPutRequest) GetPortAclOk() (*map[string]PacketbrokerPutRequestPortAclValue, bool)`

GetPortAclOk returns a tuple with the PortAcl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPortAcl

`func (o *PacketbrokerPutRequest) SetPortAcl(v map[string]PacketbrokerPutRequestPortAclValue)`

SetPortAcl sets PortAcl field to given value.

### HasPortAcl

`func (o *PacketbrokerPutRequest) HasPortAcl() bool`

HasPortAcl returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


