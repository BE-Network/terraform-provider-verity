# PodsPutRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Pod** | Pointer to [**map[string]PodsPutRequestPodValue**](PodsPutRequestPodValue.md) |  | [optional] 

## Methods

### NewPodsPutRequest

`func NewPodsPutRequest() *PodsPutRequest`

NewPodsPutRequest instantiates a new PodsPutRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPodsPutRequestWithDefaults

`func NewPodsPutRequestWithDefaults() *PodsPutRequest`

NewPodsPutRequestWithDefaults instantiates a new PodsPutRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPod

`func (o *PodsPutRequest) GetPod() map[string]PodsPutRequestPodValue`

GetPod returns the Pod field if non-nil, zero value otherwise.

### GetPodOk

`func (o *PodsPutRequest) GetPodOk() (*map[string]PodsPutRequestPodValue, bool)`

GetPodOk returns a tuple with the Pod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPod

`func (o *PodsPutRequest) SetPod(v map[string]PodsPutRequestPodValue)`

SetPod sets Pod field to given value.

### HasPod

`func (o *PodsPutRequest) HasPod() bool`

HasPod returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


