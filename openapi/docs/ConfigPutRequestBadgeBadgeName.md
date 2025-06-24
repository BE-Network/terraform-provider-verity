# ConfigPutRequestBadgeBadgeName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Color** | Pointer to **string** | Badge color | [optional] 
**Number** | Pointer to **int32** | Badge number | [optional] 
**ObjectProperties** | Pointer to [**ConfigPutRequestIpv4PrefixListIpv4PrefixListNameObjectProperties**](ConfigPutRequestIpv4PrefixListIpv4PrefixListNameObjectProperties.md) |  | [optional] 

## Methods

### NewConfigPutRequestBadgeBadgeName

`func NewConfigPutRequestBadgeBadgeName() *ConfigPutRequestBadgeBadgeName`

NewConfigPutRequestBadgeBadgeName instantiates a new ConfigPutRequestBadgeBadgeName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestBadgeBadgeNameWithDefaults

`func NewConfigPutRequestBadgeBadgeNameWithDefaults() *ConfigPutRequestBadgeBadgeName`

NewConfigPutRequestBadgeBadgeNameWithDefaults instantiates a new ConfigPutRequestBadgeBadgeName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestBadgeBadgeName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestBadgeBadgeName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestBadgeBadgeName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestBadgeBadgeName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetColor

`func (o *ConfigPutRequestBadgeBadgeName) GetColor() string`

GetColor returns the Color field if non-nil, zero value otherwise.

### GetColorOk

`func (o *ConfigPutRequestBadgeBadgeName) GetColorOk() (*string, bool)`

GetColorOk returns a tuple with the Color field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetColor

`func (o *ConfigPutRequestBadgeBadgeName) SetColor(v string)`

SetColor sets Color field to given value.

### HasColor

`func (o *ConfigPutRequestBadgeBadgeName) HasColor() bool`

HasColor returns a boolean if a field has been set.

### GetNumber

`func (o *ConfigPutRequestBadgeBadgeName) GetNumber() int32`

GetNumber returns the Number field if non-nil, zero value otherwise.

### GetNumberOk

`func (o *ConfigPutRequestBadgeBadgeName) GetNumberOk() (*int32, bool)`

GetNumberOk returns a tuple with the Number field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNumber

`func (o *ConfigPutRequestBadgeBadgeName) SetNumber(v int32)`

SetNumber sets Number field to given value.

### HasNumber

`func (o *ConfigPutRequestBadgeBadgeName) HasNumber() bool`

HasNumber returns a boolean if a field has been set.

### GetObjectProperties

`func (o *ConfigPutRequestBadgeBadgeName) GetObjectProperties() ConfigPutRequestIpv4PrefixListIpv4PrefixListNameObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ConfigPutRequestBadgeBadgeName) GetObjectPropertiesOk() (*ConfigPutRequestIpv4PrefixListIpv4PrefixListNameObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ConfigPutRequestBadgeBadgeName) SetObjectProperties(v ConfigPutRequestIpv4PrefixListIpv4PrefixListNameObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ConfigPutRequestBadgeBadgeName) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


