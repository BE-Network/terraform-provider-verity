# AclsPutRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**IpFilter** | Pointer to [**map[string]ConfigPutRequestIpv4FilterIpv4FilterName**](ConfigPutRequestIpv4FilterIpv4FilterName.md) |  | [optional] 

## Methods

### NewAclsPutRequest

`func NewAclsPutRequest() *AclsPutRequest`

NewAclsPutRequest instantiates a new AclsPutRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAclsPutRequestWithDefaults

`func NewAclsPutRequestWithDefaults() *AclsPutRequest`

NewAclsPutRequestWithDefaults instantiates a new AclsPutRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetIpFilter

`func (o *AclsPutRequest) GetIpFilter() map[string]ConfigPutRequestIpv4FilterIpv4FilterName`

GetIpFilter returns the IpFilter field if non-nil, zero value otherwise.

### GetIpFilterOk

`func (o *AclsPutRequest) GetIpFilterOk() (*map[string]ConfigPutRequestIpv4FilterIpv4FilterName, bool)`

GetIpFilterOk returns a tuple with the IpFilter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpFilter

`func (o *AclsPutRequest) SetIpFilter(v map[string]ConfigPutRequestIpv4FilterIpv4FilterName)`

SetIpFilter sets IpFilter field to given value.

### HasIpFilter

`func (o *AclsPutRequest) HasIpFilter() bool`

HasIpFilter returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


