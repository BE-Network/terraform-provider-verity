# LagsPutRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Lag** | Pointer to [**map[string]LagsPutRequestLagValue**](LagsPutRequestLagValue.md) |  | [optional] 

## Methods

### NewLagsPutRequest

`func NewLagsPutRequest() *LagsPutRequest`

NewLagsPutRequest instantiates a new LagsPutRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewLagsPutRequestWithDefaults

`func NewLagsPutRequestWithDefaults() *LagsPutRequest`

NewLagsPutRequestWithDefaults instantiates a new LagsPutRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLag

`func (o *LagsPutRequest) GetLag() map[string]LagsPutRequestLagValue`

GetLag returns the Lag field if non-nil, zero value otherwise.

### GetLagOk

`func (o *LagsPutRequest) GetLagOk() (*map[string]LagsPutRequestLagValue, bool)`

GetLagOk returns a tuple with the Lag field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLag

`func (o *LagsPutRequest) SetLag(v map[string]LagsPutRequestLagValue)`

SetLag sets Lag field to given value.

### HasLag

`func (o *LagsPutRequest) HasLag() bool`

HasLag returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


