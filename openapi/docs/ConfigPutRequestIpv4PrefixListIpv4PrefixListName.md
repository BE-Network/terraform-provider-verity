# ConfigPutRequestIpv4PrefixListIpv4PrefixListName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**Lists** | Pointer to [**[]ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner**](ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**ConfigPutRequestIpv4PrefixListIpv4PrefixListNameObjectProperties**](ConfigPutRequestIpv4PrefixListIpv4PrefixListNameObjectProperties.md) |  | [optional] 

## Methods

### NewConfigPutRequestIpv4PrefixListIpv4PrefixListName

`func NewConfigPutRequestIpv4PrefixListIpv4PrefixListName() *ConfigPutRequestIpv4PrefixListIpv4PrefixListName`

NewConfigPutRequestIpv4PrefixListIpv4PrefixListName instantiates a new ConfigPutRequestIpv4PrefixListIpv4PrefixListName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestIpv4PrefixListIpv4PrefixListNameWithDefaults

`func NewConfigPutRequestIpv4PrefixListIpv4PrefixListNameWithDefaults() *ConfigPutRequestIpv4PrefixListIpv4PrefixListName`

NewConfigPutRequestIpv4PrefixListIpv4PrefixListNameWithDefaults instantiates a new ConfigPutRequestIpv4PrefixListIpv4PrefixListName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListName) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListName) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListName) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListName) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetLists

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListName) GetLists() []ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner`

GetLists returns the Lists field if non-nil, zero value otherwise.

### GetListsOk

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListName) GetListsOk() (*[]ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner, bool)`

GetListsOk returns a tuple with the Lists field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLists

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListName) SetLists(v []ConfigPutRequestIpv4PrefixListIpv4PrefixListNameListsInner)`

SetLists sets Lists field to given value.

### HasLists

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListName) HasLists() bool`

HasLists returns a boolean if a field has been set.

### GetObjectProperties

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListName) GetObjectProperties() ConfigPutRequestIpv4PrefixListIpv4PrefixListNameObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListName) GetObjectPropertiesOk() (*ConfigPutRequestIpv4PrefixListIpv4PrefixListNameObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListName) SetObjectProperties(v ConfigPutRequestIpv4PrefixListIpv4PrefixListNameObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ConfigPutRequestIpv4PrefixListIpv4PrefixListName) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


