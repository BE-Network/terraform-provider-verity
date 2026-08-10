# SusPutRequestSuValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Template Name. Must be unique within type. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to true]
**Pod** | Pointer to **string** | Pod this SU is assigned to | [optional] [default to ""]
**PodRefType** | Pointer to **string** | Object type for pod field | [optional] 
**Position** | Pointer to **NullableFloat64** | Position of the Switch | [optional] 
**ObjectProperties** | Pointer to [**AclsPutRequestIpFilterValueObjectProperties**](AclsPutRequestIpFilterValueObjectProperties.md) |  | [optional] 

## Methods

### NewSusPutRequestSuValue

`func NewSusPutRequestSuValue() *SusPutRequestSuValue`

NewSusPutRequestSuValue instantiates a new SusPutRequestSuValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSusPutRequestSuValueWithDefaults

`func NewSusPutRequestSuValueWithDefaults() *SusPutRequestSuValue`

NewSusPutRequestSuValueWithDefaults instantiates a new SusPutRequestSuValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *SusPutRequestSuValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *SusPutRequestSuValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *SusPutRequestSuValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *SusPutRequestSuValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *SusPutRequestSuValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *SusPutRequestSuValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *SusPutRequestSuValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *SusPutRequestSuValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetPod

`func (o *SusPutRequestSuValue) GetPod() string`

GetPod returns the Pod field if non-nil, zero value otherwise.

### GetPodOk

`func (o *SusPutRequestSuValue) GetPodOk() (*string, bool)`

GetPodOk returns a tuple with the Pod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPod

`func (o *SusPutRequestSuValue) SetPod(v string)`

SetPod sets Pod field to given value.

### HasPod

`func (o *SusPutRequestSuValue) HasPod() bool`

HasPod returns a boolean if a field has been set.

### GetPodRefType

`func (o *SusPutRequestSuValue) GetPodRefType() string`

GetPodRefType returns the PodRefType field if non-nil, zero value otherwise.

### GetPodRefTypeOk

`func (o *SusPutRequestSuValue) GetPodRefTypeOk() (*string, bool)`

GetPodRefTypeOk returns a tuple with the PodRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPodRefType

`func (o *SusPutRequestSuValue) SetPodRefType(v string)`

SetPodRefType sets PodRefType field to given value.

### HasPodRefType

`func (o *SusPutRequestSuValue) HasPodRefType() bool`

HasPodRefType returns a boolean if a field has been set.

### GetPosition

`func (o *SusPutRequestSuValue) GetPosition() float64`

GetPosition returns the Position field if non-nil, zero value otherwise.

### GetPositionOk

`func (o *SusPutRequestSuValue) GetPositionOk() (*float64, bool)`

GetPositionOk returns a tuple with the Position field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPosition

`func (o *SusPutRequestSuValue) SetPosition(v float64)`

SetPosition sets Position field to given value.

### HasPosition

`func (o *SusPutRequestSuValue) HasPosition() bool`

HasPosition returns a boolean if a field has been set.

### SetPositionNil

`func (o *SusPutRequestSuValue) SetPositionNil(b bool)`

 SetPositionNil sets the value for Position to be an explicit nil

### UnsetPosition
`func (o *SusPutRequestSuValue) UnsetPosition()`

UnsetPosition ensures that no value is present for Position, not even an explicit nil
### GetObjectProperties

`func (o *SusPutRequestSuValue) GetObjectProperties() AclsPutRequestIpFilterValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *SusPutRequestSuValue) GetObjectPropertiesOk() (*AclsPutRequestIpFilterValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *SusPutRequestSuValue) SetObjectProperties(v AclsPutRequestIpFilterValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *SusPutRequestSuValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


