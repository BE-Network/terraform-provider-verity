# SitesPatchRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Site** | Pointer to [**map[string]SitesPatchRequestSiteValue**](SitesPatchRequestSiteValue.md) |  | [optional] 

## Methods

### NewSitesPatchRequest

`func NewSitesPatchRequest() *SitesPatchRequest`

NewSitesPatchRequest instantiates a new SitesPatchRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSitesPatchRequestWithDefaults

`func NewSitesPatchRequestWithDefaults() *SitesPatchRequest`

NewSitesPatchRequestWithDefaults instantiates a new SitesPatchRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSite

`func (o *SitesPatchRequest) GetSite() map[string]SitesPatchRequestSiteValue`

GetSite returns the Site field if non-nil, zero value otherwise.

### GetSiteOk

`func (o *SitesPatchRequest) GetSiteOk() (*map[string]SitesPatchRequestSiteValue, bool)`

GetSiteOk returns a tuple with the Site field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSite

`func (o *SitesPatchRequest) SetSite(v map[string]SitesPatchRequestSiteValue)`

SetSite sets Site field to given value.

### HasSite

`func (o *SitesPatchRequest) HasSite() bool`

HasSite returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


