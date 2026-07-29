# PlanesPutRequestPlaneValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Template Name. Must be unique within type. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to true]
**Site** | Pointer to **string** | Fabric this Plane is assigned to | [optional] [default to ""]
**SiteRefType** | Pointer to **string** | Object type for site field | [optional] 
**Position** | Pointer to **NullableFloat32** | Position of the Plane | [optional] 
**ObjectProperties** | Pointer to [**AclsPutRequestIpFilterValueObjectProperties**](AclsPutRequestIpFilterValueObjectProperties.md) |  | [optional] 

## Methods

### NewPlanesPutRequestPlaneValue

`func NewPlanesPutRequestPlaneValue() *PlanesPutRequestPlaneValue`

NewPlanesPutRequestPlaneValue instantiates a new PlanesPutRequestPlaneValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPlanesPutRequestPlaneValueWithDefaults

`func NewPlanesPutRequestPlaneValueWithDefaults() *PlanesPutRequestPlaneValue`

NewPlanesPutRequestPlaneValueWithDefaults instantiates a new PlanesPutRequestPlaneValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *PlanesPutRequestPlaneValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *PlanesPutRequestPlaneValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *PlanesPutRequestPlaneValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *PlanesPutRequestPlaneValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *PlanesPutRequestPlaneValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *PlanesPutRequestPlaneValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *PlanesPutRequestPlaneValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *PlanesPutRequestPlaneValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetSite

`func (o *PlanesPutRequestPlaneValue) GetSite() string`

GetSite returns the Site field if non-nil, zero value otherwise.

### GetSiteOk

`func (o *PlanesPutRequestPlaneValue) GetSiteOk() (*string, bool)`

GetSiteOk returns a tuple with the Site field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSite

`func (o *PlanesPutRequestPlaneValue) SetSite(v string)`

SetSite sets Site field to given value.

### HasSite

`func (o *PlanesPutRequestPlaneValue) HasSite() bool`

HasSite returns a boolean if a field has been set.

### GetSiteRefType

`func (o *PlanesPutRequestPlaneValue) GetSiteRefType() string`

GetSiteRefType returns the SiteRefType field if non-nil, zero value otherwise.

### GetSiteRefTypeOk

`func (o *PlanesPutRequestPlaneValue) GetSiteRefTypeOk() (*string, bool)`

GetSiteRefTypeOk returns a tuple with the SiteRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteRefType

`func (o *PlanesPutRequestPlaneValue) SetSiteRefType(v string)`

SetSiteRefType sets SiteRefType field to given value.

### HasSiteRefType

`func (o *PlanesPutRequestPlaneValue) HasSiteRefType() bool`

HasSiteRefType returns a boolean if a field has been set.

### GetPosition

`func (o *PlanesPutRequestPlaneValue) GetPosition() float32`

GetPosition returns the Position field if non-nil, zero value otherwise.

### GetPositionOk

`func (o *PlanesPutRequestPlaneValue) GetPositionOk() (*float32, bool)`

GetPositionOk returns a tuple with the Position field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPosition

`func (o *PlanesPutRequestPlaneValue) SetPosition(v float32)`

SetPosition sets Position field to given value.

### HasPosition

`func (o *PlanesPutRequestPlaneValue) HasPosition() bool`

HasPosition returns a boolean if a field has been set.

### SetPositionNil

`func (o *PlanesPutRequestPlaneValue) SetPositionNil(b bool)`

 SetPositionNil sets the value for Position to be an explicit nil

### UnsetPosition
`func (o *PlanesPutRequestPlaneValue) UnsetPosition()`

UnsetPosition ensures that no value is present for Position, not even an explicit nil
### GetObjectProperties

`func (o *PlanesPutRequestPlaneValue) GetObjectProperties() AclsPutRequestIpFilterValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *PlanesPutRequestPlaneValue) GetObjectPropertiesOk() (*AclsPutRequestIpFilterValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *PlanesPutRequestPlaneValue) SetObjectProperties(v AclsPutRequestIpFilterValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *PlanesPutRequestPlaneValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


