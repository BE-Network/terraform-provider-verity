# BadgesPutRequestBadgeValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to true]
**Color** | Pointer to **string** | Color of Badge | [optional] [default to "next available color"]
**Number** | Pointer to **int32** | Number of Badge | [optional] 
**ObjectProperties** | Pointer to [**AclsPutRequestIpFilterValueObjectProperties**](AclsPutRequestIpFilterValueObjectProperties.md) |  | [optional] 

## Methods

### NewBadgesPutRequestBadgeValue

`func NewBadgesPutRequestBadgeValue() *BadgesPutRequestBadgeValue`

NewBadgesPutRequestBadgeValue instantiates a new BadgesPutRequestBadgeValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBadgesPutRequestBadgeValueWithDefaults

`func NewBadgesPutRequestBadgeValueWithDefaults() *BadgesPutRequestBadgeValue`

NewBadgesPutRequestBadgeValueWithDefaults instantiates a new BadgesPutRequestBadgeValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *BadgesPutRequestBadgeValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *BadgesPutRequestBadgeValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *BadgesPutRequestBadgeValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *BadgesPutRequestBadgeValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *BadgesPutRequestBadgeValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *BadgesPutRequestBadgeValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *BadgesPutRequestBadgeValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *BadgesPutRequestBadgeValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetColor

`func (o *BadgesPutRequestBadgeValue) GetColor() string`

GetColor returns the Color field if non-nil, zero value otherwise.

### GetColorOk

`func (o *BadgesPutRequestBadgeValue) GetColorOk() (*string, bool)`

GetColorOk returns a tuple with the Color field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetColor

`func (o *BadgesPutRequestBadgeValue) SetColor(v string)`

SetColor sets Color field to given value.

### HasColor

`func (o *BadgesPutRequestBadgeValue) HasColor() bool`

HasColor returns a boolean if a field has been set.

### GetNumber

`func (o *BadgesPutRequestBadgeValue) GetNumber() int32`

GetNumber returns the Number field if non-nil, zero value otherwise.

### GetNumberOk

`func (o *BadgesPutRequestBadgeValue) GetNumberOk() (*int32, bool)`

GetNumberOk returns a tuple with the Number field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNumber

`func (o *BadgesPutRequestBadgeValue) SetNumber(v int32)`

SetNumber sets Number field to given value.

### HasNumber

`func (o *BadgesPutRequestBadgeValue) HasNumber() bool`

HasNumber returns a boolean if a field has been set.

### GetObjectProperties

`func (o *BadgesPutRequestBadgeValue) GetObjectProperties() AclsPutRequestIpFilterValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *BadgesPutRequestBadgeValue) GetObjectPropertiesOk() (*AclsPutRequestIpFilterValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *BadgesPutRequestBadgeValue) SetObjectProperties(v AclsPutRequestIpFilterValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *BadgesPutRequestBadgeValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


