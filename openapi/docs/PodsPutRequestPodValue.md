# PodsPutRequestPodValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to true]
**ObjectProperties** | Pointer to [**AclsPutRequestIpFilterValueObjectProperties**](AclsPutRequestIpFilterValueObjectProperties.md) |  | [optional] 
**ExpectedSpineCount** | Pointer to **NullableInt32** | Number of spine switches expected in this pod | [optional] [default to 1]

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

### GetExpectedSpineCount

`func (o *PodsPutRequestPodValue) GetExpectedSpineCount() int32`

GetExpectedSpineCount returns the ExpectedSpineCount field if non-nil, zero value otherwise.

### GetExpectedSpineCountOk

`func (o *PodsPutRequestPodValue) GetExpectedSpineCountOk() (*int32, bool)`

GetExpectedSpineCountOk returns a tuple with the ExpectedSpineCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpectedSpineCount

`func (o *PodsPutRequestPodValue) SetExpectedSpineCount(v int32)`

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

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


