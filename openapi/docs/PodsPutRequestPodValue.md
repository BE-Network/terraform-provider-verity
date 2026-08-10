# PodsPutRequestPodValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Template Name. Must be unique within type. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to true]
**ExpectedSpineCount** | Pointer to **NullableInt64** | Number of spine switches expected in this pod | [optional] [default to 1]
**Site** | Pointer to **string** | Fabric this Pod is assigned to | [optional] [default to ""]
**SiteRefType** | Pointer to **string** | Object type for site field | [optional] 
**Position** | Pointer to **NullableFloat64** | Position of the Switch | [optional] 
**ObjectProperties** | Pointer to [**AclsPutRequestIpFilterValueObjectProperties**](AclsPutRequestIpFilterValueObjectProperties.md) |  | [optional] 

## Methods

### NewPodsPutRequestPodValue

`func NewPodsPutRequestPodValue() *PodsPutRequestPodValue`

NewPodsPutRequestPodValue instantiates a new PodsPutRequestPodValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPodsPutRequestPodValueWithDefaults

`func NewPodsPutRequestPodValueWithDefaults() *PodsPutRequestPodValue`

NewPodsPutRequestPodValueWithDefaults instantiates a new PodsPutRequestPodValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *PodsPutRequestPodValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *PodsPutRequestPodValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *PodsPutRequestPodValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *PodsPutRequestPodValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *PodsPutRequestPodValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *PodsPutRequestPodValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *PodsPutRequestPodValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *PodsPutRequestPodValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetExpectedSpineCount

`func (o *PodsPutRequestPodValue) GetExpectedSpineCount() int64`

GetExpectedSpineCount returns the ExpectedSpineCount field if non-nil, zero value otherwise.

### GetExpectedSpineCountOk

`func (o *PodsPutRequestPodValue) GetExpectedSpineCountOk() (*int64, bool)`

GetExpectedSpineCountOk returns a tuple with the ExpectedSpineCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpectedSpineCount

`func (o *PodsPutRequestPodValue) SetExpectedSpineCount(v int64)`

SetExpectedSpineCount sets ExpectedSpineCount field to given value.

### HasExpectedSpineCount

`func (o *PodsPutRequestPodValue) HasExpectedSpineCount() bool`

HasExpectedSpineCount returns a boolean if a field has been set.

### SetExpectedSpineCountNil

`func (o *PodsPutRequestPodValue) SetExpectedSpineCountNil(b bool)`

 SetExpectedSpineCountNil sets the value for ExpectedSpineCount to be an explicit nil

### UnsetExpectedSpineCount
`func (o *PodsPutRequestPodValue) UnsetExpectedSpineCount()`

UnsetExpectedSpineCount ensures that no value is present for ExpectedSpineCount, not even an explicit nil
### GetSite

`func (o *PodsPutRequestPodValue) GetSite() string`

GetSite returns the Site field if non-nil, zero value otherwise.

### GetSiteOk

`func (o *PodsPutRequestPodValue) GetSiteOk() (*string, bool)`

GetSiteOk returns a tuple with the Site field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSite

`func (o *PodsPutRequestPodValue) SetSite(v string)`

SetSite sets Site field to given value.

### HasSite

`func (o *PodsPutRequestPodValue) HasSite() bool`

HasSite returns a boolean if a field has been set.

### GetSiteRefType

`func (o *PodsPutRequestPodValue) GetSiteRefType() string`

GetSiteRefType returns the SiteRefType field if non-nil, zero value otherwise.

### GetSiteRefTypeOk

`func (o *PodsPutRequestPodValue) GetSiteRefTypeOk() (*string, bool)`

GetSiteRefTypeOk returns a tuple with the SiteRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteRefType

`func (o *PodsPutRequestPodValue) SetSiteRefType(v string)`

SetSiteRefType sets SiteRefType field to given value.

### HasSiteRefType

`func (o *PodsPutRequestPodValue) HasSiteRefType() bool`

HasSiteRefType returns a boolean if a field has been set.

### GetPosition

`func (o *PodsPutRequestPodValue) GetPosition() float64`

GetPosition returns the Position field if non-nil, zero value otherwise.

### GetPositionOk

`func (o *PodsPutRequestPodValue) GetPositionOk() (*float64, bool)`

GetPositionOk returns a tuple with the Position field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPosition

`func (o *PodsPutRequestPodValue) SetPosition(v float64)`

SetPosition sets Position field to given value.

### HasPosition

`func (o *PodsPutRequestPodValue) HasPosition() bool`

HasPosition returns a boolean if a field has been set.

### SetPositionNil

`func (o *PodsPutRequestPodValue) SetPositionNil(b bool)`

 SetPositionNil sets the value for Position to be an explicit nil

### UnsetPosition
`func (o *PodsPutRequestPodValue) UnsetPosition()`

UnsetPosition ensures that no value is present for Position, not even an explicit nil
### GetObjectProperties

`func (o *PodsPutRequestPodValue) GetObjectProperties() AclsPutRequestIpFilterValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *PodsPutRequestPodValue) GetObjectPropertiesOk() (*AclsPutRequestIpFilterValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *PodsPutRequestPodValue) SetObjectProperties(v AclsPutRequestIpFilterValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *PodsPutRequestPodValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


