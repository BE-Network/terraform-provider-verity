# RacksPutRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Rack** | Pointer to [**map[string]RacksPutRequestRackValue**](RacksPutRequestRackValue.md) |  | [optional] 

## Methods

### NewRacksPutRequest

`func NewRacksPutRequest() *RacksPutRequest`

NewRacksPutRequest instantiates a new RacksPutRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRacksPutRequestWithDefaults

`func NewRacksPutRequestWithDefaults() *RacksPutRequest`

NewRacksPutRequestWithDefaults instantiates a new RacksPutRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRack

`func (o *RacksPutRequest) GetRack() map[string]RacksPutRequestRackValue`

GetRack returns the Rack field if non-nil, zero value otherwise.

### GetRackOk

`func (o *RacksPutRequest) GetRackOk() (*map[string]RacksPutRequestRackValue, bool)`

GetRackOk returns a tuple with the Rack field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRack

`func (o *RacksPutRequest) SetRack(v map[string]RacksPutRequestRackValue)`

SetRack sets Rack field to given value.

### HasRack

`func (o *RacksPutRequest) HasRack() bool`

HasRack returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


