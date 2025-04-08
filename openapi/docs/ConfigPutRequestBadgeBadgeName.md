# ConfigPutRequestBadgeBadgeName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
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


