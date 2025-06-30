# BundlesPutRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**EndpointBundle** | Pointer to [**map[string]BundlesPutRequestEndpointBundleValue**](BundlesPutRequestEndpointBundleValue.md) |  | [optional] 

## Methods

### NewBundlesPutRequest

`func NewBundlesPutRequest() *BundlesPutRequest`

NewBundlesPutRequest instantiates a new BundlesPutRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBundlesPutRequestWithDefaults

`func NewBundlesPutRequestWithDefaults() *BundlesPutRequest`

NewBundlesPutRequestWithDefaults instantiates a new BundlesPutRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEndpointBundle

`func (o *BundlesPutRequest) GetEndpointBundle() map[string]BundlesPutRequestEndpointBundleValue`

GetEndpointBundle returns the EndpointBundle field if non-nil, zero value otherwise.

### GetEndpointBundleOk

`func (o *BundlesPutRequest) GetEndpointBundleOk() (*map[string]BundlesPutRequestEndpointBundleValue, bool)`

GetEndpointBundleOk returns a tuple with the EndpointBundle field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpointBundle

`func (o *BundlesPutRequest) SetEndpointBundle(v map[string]BundlesPutRequestEndpointBundleValue)`

SetEndpointBundle sets EndpointBundle field to given value.

### HasEndpointBundle

`func (o *BundlesPutRequest) HasEndpointBundle() bool`

HasEndpointBundle returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


