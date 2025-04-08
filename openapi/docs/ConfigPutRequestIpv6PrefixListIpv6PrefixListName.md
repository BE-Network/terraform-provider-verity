# ConfigPutRequestIpv6PrefixListIpv6PrefixListName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**Lists** | Pointer to [**[]ConfigPutRequestIpv6PrefixListIpv6PrefixListNameListsInner**](ConfigPutRequestIpv6PrefixListIpv6PrefixListNameListsInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**ConfigPutRequestIpv4PrefixListIpv4PrefixListNameObjectProperties**](ConfigPutRequestIpv4PrefixListIpv4PrefixListNameObjectProperties.md) |  | [optional] 

## Methods

### NewConfigPutRequestIpv6PrefixListIpv6PrefixListName

`func NewConfigPutRequestIpv6PrefixListIpv6PrefixListName() *ConfigPutRequestIpv6PrefixListIpv6PrefixListName`

NewConfigPutRequestIpv6PrefixListIpv6PrefixListName instantiates a new ConfigPutRequestIpv6PrefixListIpv6PrefixListName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestIpv6PrefixListIpv6PrefixListNameWithDefaults

`func NewConfigPutRequestIpv6PrefixListIpv6PrefixListNameWithDefaults() *ConfigPutRequestIpv6PrefixListIpv6PrefixListName`

NewConfigPutRequestIpv6PrefixListIpv6PrefixListNameWithDefaults instantiates a new ConfigPutRequestIpv6PrefixListIpv6PrefixListName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestIpv6PrefixListIpv6PrefixListName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestIpv6PrefixListIpv6PrefixListName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestIpv6PrefixListIpv6PrefixListName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestIpv6PrefixListIpv6PrefixListName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ConfigPutRequestIpv6PrefixListIpv6PrefixListName) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestIpv6PrefixListIpv6PrefixListName) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestIpv6PrefixListIpv6PrefixListName) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestIpv6PrefixListIpv6PrefixListName) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetLists

`func (o *ConfigPutRequestIpv6PrefixListIpv6PrefixListName) GetLists() []ConfigPutRequestIpv6PrefixListIpv6PrefixListNameListsInner`

GetLists returns the Lists field if non-nil, zero value otherwise.

### GetListsOk

`func (o *ConfigPutRequestIpv6PrefixListIpv6PrefixListName) GetListsOk() (*[]ConfigPutRequestIpv6PrefixListIpv6PrefixListNameListsInner, bool)`

GetListsOk returns a tuple with the Lists field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLists

`func (o *ConfigPutRequestIpv6PrefixListIpv6PrefixListName) SetLists(v []ConfigPutRequestIpv6PrefixListIpv6PrefixListNameListsInner)`

SetLists sets Lists field to given value.

### HasLists

`func (o *ConfigPutRequestIpv6PrefixListIpv6PrefixListName) HasLists() bool`

HasLists returns a boolean if a field has been set.

### GetObjectProperties

`func (o *ConfigPutRequestIpv6PrefixListIpv6PrefixListName) GetObjectProperties() ConfigPutRequestIpv4PrefixListIpv4PrefixListNameObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ConfigPutRequestIpv6PrefixListIpv6PrefixListName) GetObjectPropertiesOk() (*ConfigPutRequestIpv4PrefixListIpv4PrefixListNameObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ConfigPutRequestIpv6PrefixListIpv6PrefixListName) SetObjectProperties(v ConfigPutRequestIpv4PrefixListIpv4PrefixListNameObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ConfigPutRequestIpv6PrefixListIpv6PrefixListName) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


