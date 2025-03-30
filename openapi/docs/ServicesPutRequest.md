# ServicesPutRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Service** | Pointer to [**map[string]ConfigPutRequestServiceServiceName**](ConfigPutRequestServiceServiceName.md) |  | [optional] 

## Methods

### NewServicesPutRequest

`func NewServicesPutRequest() *ServicesPutRequest`

NewServicesPutRequest instantiates a new ServicesPutRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewServicesPutRequestWithDefaults

`func NewServicesPutRequestWithDefaults() *ServicesPutRequest`

NewServicesPutRequestWithDefaults instantiates a new ServicesPutRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetService

`func (o *ServicesPutRequest) GetService() map[string]ConfigPutRequestServiceServiceName`

GetService returns the Service field if non-nil, zero value otherwise.

### GetServiceOk

`func (o *ServicesPutRequest) GetServiceOk() (*map[string]ConfigPutRequestServiceServiceName, bool)`

GetServiceOk returns a tuple with the Service field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetService

`func (o *ServicesPutRequest) SetService(v map[string]ConfigPutRequestServiceServiceName)`

SetService sets Service field to given value.

### HasService

`func (o *ServicesPutRequest) HasService() bool`

HasService returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


