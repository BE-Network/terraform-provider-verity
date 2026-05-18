# SpineplanesPutRequestSpinePlaneValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Template Name. Must be unique within type. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to true]
**Site** | Pointer to **string** | Fabric this Spine Plane is assigned to | [optional] [default to ""]
**SiteRefType** | Pointer to **string** | Object type for site field | [optional] 
**ObjectProperties** | Pointer to [**AclsPutRequestIpFilterValueObjectProperties**](AclsPutRequestIpFilterValueObjectProperties.md) |  | [optional] 

## Methods

### NewSpineplanesPutRequestSpinePlaneValue

`func NewSpineplanesPutRequestSpinePlaneValue() *SpineplanesPutRequestSpinePlaneValue`

NewSpineplanesPutRequestSpinePlaneValue instantiates a new SpineplanesPutRequestSpinePlaneValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSpineplanesPutRequestSpinePlaneValueWithDefaults

`func NewSpineplanesPutRequestSpinePlaneValueWithDefaults() *SpineplanesPutRequestSpinePlaneValue`

NewSpineplanesPutRequestSpinePlaneValueWithDefaults instantiates a new SpineplanesPutRequestSpinePlaneValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *SpineplanesPutRequestSpinePlaneValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *SpineplanesPutRequestSpinePlaneValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *SpineplanesPutRequestSpinePlaneValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *SpineplanesPutRequestSpinePlaneValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *SpineplanesPutRequestSpinePlaneValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *SpineplanesPutRequestSpinePlaneValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *SpineplanesPutRequestSpinePlaneValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *SpineplanesPutRequestSpinePlaneValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetSite

`func (o *SpineplanesPutRequestSpinePlaneValue) GetSite() string`

GetSite returns the Site field if non-nil, zero value otherwise.

### GetSiteOk

`func (o *SpineplanesPutRequestSpinePlaneValue) GetSiteOk() (*string, bool)`

GetSiteOk returns a tuple with the Site field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSite

`func (o *SpineplanesPutRequestSpinePlaneValue) SetSite(v string)`

SetSite sets Site field to given value.

### HasSite

`func (o *SpineplanesPutRequestSpinePlaneValue) HasSite() bool`

HasSite returns a boolean if a field has been set.

### GetSiteRefType

`func (o *SpineplanesPutRequestSpinePlaneValue) GetSiteRefType() string`

GetSiteRefType returns the SiteRefType field if non-nil, zero value otherwise.

### GetSiteRefTypeOk

`func (o *SpineplanesPutRequestSpinePlaneValue) GetSiteRefTypeOk() (*string, bool)`

GetSiteRefTypeOk returns a tuple with the SiteRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteRefType

`func (o *SpineplanesPutRequestSpinePlaneValue) SetSiteRefType(v string)`

SetSiteRefType sets SiteRefType field to given value.

### HasSiteRefType

`func (o *SpineplanesPutRequestSpinePlaneValue) HasSiteRefType() bool`

HasSiteRefType returns a boolean if a field has been set.

### GetObjectProperties

`func (o *SpineplanesPutRequestSpinePlaneValue) GetObjectProperties() AclsPutRequestIpFilterValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *SpineplanesPutRequestSpinePlaneValue) GetObjectPropertiesOk() (*AclsPutRequestIpFilterValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *SpineplanesPutRequestSpinePlaneValue) SetObjectProperties(v AclsPutRequestIpFilterValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *SpineplanesPutRequestSpinePlaneValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


