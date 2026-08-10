# RacksPutRequestRackValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Template Name. Must be unique within type. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to true]
**Position** | Pointer to **NullableFloat64** | Position of the Rack | [optional] 
**Su** | Pointer to **string** | SU this Rack is assigned to | [optional] [default to ""]
**SuRefType** | Pointer to **string** | Object type for su field | [optional] 
**ObjectProperties** | Pointer to [**AclsPutRequestIpFilterValueObjectProperties**](AclsPutRequestIpFilterValueObjectProperties.md) |  | [optional] 

## Methods

### NewRacksPutRequestRackValue

`func NewRacksPutRequestRackValue() *RacksPutRequestRackValue`

NewRacksPutRequestRackValue instantiates a new RacksPutRequestRackValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRacksPutRequestRackValueWithDefaults

`func NewRacksPutRequestRackValueWithDefaults() *RacksPutRequestRackValue`

NewRacksPutRequestRackValueWithDefaults instantiates a new RacksPutRequestRackValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *RacksPutRequestRackValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *RacksPutRequestRackValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *RacksPutRequestRackValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *RacksPutRequestRackValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *RacksPutRequestRackValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *RacksPutRequestRackValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *RacksPutRequestRackValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *RacksPutRequestRackValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetPosition

`func (o *RacksPutRequestRackValue) GetPosition() float64`

GetPosition returns the Position field if non-nil, zero value otherwise.

### GetPositionOk

`func (o *RacksPutRequestRackValue) GetPositionOk() (*float64, bool)`

GetPositionOk returns a tuple with the Position field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPosition

`func (o *RacksPutRequestRackValue) SetPosition(v float64)`

SetPosition sets Position field to given value.

### HasPosition

`func (o *RacksPutRequestRackValue) HasPosition() bool`

HasPosition returns a boolean if a field has been set.

### SetPositionNil

`func (o *RacksPutRequestRackValue) SetPositionNil(b bool)`

 SetPositionNil sets the value for Position to be an explicit nil

### UnsetPosition
`func (o *RacksPutRequestRackValue) UnsetPosition()`

UnsetPosition ensures that no value is present for Position, not even an explicit nil
### GetSu

`func (o *RacksPutRequestRackValue) GetSu() string`

GetSu returns the Su field if non-nil, zero value otherwise.

### GetSuOk

`func (o *RacksPutRequestRackValue) GetSuOk() (*string, bool)`

GetSuOk returns a tuple with the Su field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSu

`func (o *RacksPutRequestRackValue) SetSu(v string)`

SetSu sets Su field to given value.

### HasSu

`func (o *RacksPutRequestRackValue) HasSu() bool`

HasSu returns a boolean if a field has been set.

### GetSuRefType

`func (o *RacksPutRequestRackValue) GetSuRefType() string`

GetSuRefType returns the SuRefType field if non-nil, zero value otherwise.

### GetSuRefTypeOk

`func (o *RacksPutRequestRackValue) GetSuRefTypeOk() (*string, bool)`

GetSuRefTypeOk returns a tuple with the SuRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSuRefType

`func (o *RacksPutRequestRackValue) SetSuRefType(v string)`

SetSuRefType sets SuRefType field to given value.

### HasSuRefType

`func (o *RacksPutRequestRackValue) HasSuRefType() bool`

HasSuRefType returns a boolean if a field has been set.

### GetObjectProperties

`func (o *RacksPutRequestRackValue) GetObjectProperties() AclsPutRequestIpFilterValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *RacksPutRequestRackValue) GetObjectPropertiesOk() (*AclsPutRequestIpFilterValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *RacksPutRequestRackValue) SetObjectProperties(v AclsPutRequestIpFilterValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *RacksPutRequestRackValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


