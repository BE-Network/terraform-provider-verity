# ConfigPutRequestIpPrefixListIpPrefixListName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**Lists** | Pointer to [**[]ConfigPutRequestIpPrefixListIpPrefixListNameListsInner**](ConfigPutRequestIpPrefixListIpPrefixListNameListsInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties**](ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties.md) |  | [optional] 

## Methods

### NewConfigPutRequestIpPrefixListIpPrefixListName

`func NewConfigPutRequestIpPrefixListIpPrefixListName() *ConfigPutRequestIpPrefixListIpPrefixListName`

NewConfigPutRequestIpPrefixListIpPrefixListName instantiates a new ConfigPutRequestIpPrefixListIpPrefixListName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestIpPrefixListIpPrefixListNameWithDefaults

`func NewConfigPutRequestIpPrefixListIpPrefixListNameWithDefaults() *ConfigPutRequestIpPrefixListIpPrefixListName`

NewConfigPutRequestIpPrefixListIpPrefixListNameWithDefaults instantiates a new ConfigPutRequestIpPrefixListIpPrefixListName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestIpPrefixListIpPrefixListName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestIpPrefixListIpPrefixListName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestIpPrefixListIpPrefixListName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestIpPrefixListIpPrefixListName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ConfigPutRequestIpPrefixListIpPrefixListName) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestIpPrefixListIpPrefixListName) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestIpPrefixListIpPrefixListName) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestIpPrefixListIpPrefixListName) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetLists

`func (o *ConfigPutRequestIpPrefixListIpPrefixListName) GetLists() []ConfigPutRequestIpPrefixListIpPrefixListNameListsInner`

GetLists returns the Lists field if non-nil, zero value otherwise.

### GetListsOk

`func (o *ConfigPutRequestIpPrefixListIpPrefixListName) GetListsOk() (*[]ConfigPutRequestIpPrefixListIpPrefixListNameListsInner, bool)`

GetListsOk returns a tuple with the Lists field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLists

`func (o *ConfigPutRequestIpPrefixListIpPrefixListName) SetLists(v []ConfigPutRequestIpPrefixListIpPrefixListNameListsInner)`

SetLists sets Lists field to given value.

### HasLists

`func (o *ConfigPutRequestIpPrefixListIpPrefixListName) HasLists() bool`

HasLists returns a boolean if a field has been set.

### GetObjectProperties

`func (o *ConfigPutRequestIpPrefixListIpPrefixListName) GetObjectProperties() ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ConfigPutRequestIpPrefixListIpPrefixListName) GetObjectPropertiesOk() (*ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ConfigPutRequestIpPrefixListIpPrefixListName) SetObjectProperties(v ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ConfigPutRequestIpPrefixListIpPrefixListName) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


