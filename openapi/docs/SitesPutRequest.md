# SitesPutRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Site** | Pointer to [**map[string]SitesPutRequestSiteValue**](SitesPutRequestSiteValue.md) |  | [optional] 

## Methods

### NewSitesPutRequest

`func NewSitesPutRequest() *SitesPutRequest`

NewSitesPutRequest instantiates a new SitesPutRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSitesPutRequestWithDefaults

`func NewSitesPutRequestWithDefaults() *SitesPutRequest`

NewSitesPutRequestWithDefaults instantiates a new SitesPutRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSite

`func (o *SitesPutRequest) GetSite() map[string]SitesPutRequestSiteValue`

GetSite returns the Site field if non-nil, zero value otherwise.

### GetSiteOk

`func (o *SitesPutRequest) GetSiteOk() (*map[string]SitesPutRequestSiteValue, bool)`

GetSiteOk returns a tuple with the Site field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSite

`func (o *SitesPutRequest) SetSite(v map[string]SitesPutRequestSiteValue)`

SetSite sets Site field to given value.

### HasSite

`func (o *SitesPutRequest) HasSite() bool`

HasSite returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


