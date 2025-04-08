# ConfigPutRequestAsPathAccessListAsPathAccessListName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**PermitDeny** | Pointer to **string** | Action upon match of Community Strings. | [optional] [default to "permit"]
**Lists** | Pointer to [**[]ConfigPutRequestAsPathAccessListAsPathAccessListNameListsInner**](ConfigPutRequestAsPathAccessListAsPathAccessListNameListsInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**ConfigPutRequestIpv4PrefixListIpv4PrefixListNameObjectProperties**](ConfigPutRequestIpv4PrefixListIpv4PrefixListNameObjectProperties.md) |  | [optional] 

## Methods

### NewConfigPutRequestAsPathAccessListAsPathAccessListName

`func NewConfigPutRequestAsPathAccessListAsPathAccessListName() *ConfigPutRequestAsPathAccessListAsPathAccessListName`

NewConfigPutRequestAsPathAccessListAsPathAccessListName instantiates a new ConfigPutRequestAsPathAccessListAsPathAccessListName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestAsPathAccessListAsPathAccessListNameWithDefaults

`func NewConfigPutRequestAsPathAccessListAsPathAccessListNameWithDefaults() *ConfigPutRequestAsPathAccessListAsPathAccessListName`

NewConfigPutRequestAsPathAccessListAsPathAccessListNameWithDefaults instantiates a new ConfigPutRequestAsPathAccessListAsPathAccessListName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetPermitDeny

`func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) GetPermitDeny() string`

GetPermitDeny returns the PermitDeny field if non-nil, zero value otherwise.

### GetPermitDenyOk

`func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) GetPermitDenyOk() (*string, bool)`

GetPermitDenyOk returns a tuple with the PermitDeny field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPermitDeny

`func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) SetPermitDeny(v string)`

SetPermitDeny sets PermitDeny field to given value.

### HasPermitDeny

`func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) HasPermitDeny() bool`

HasPermitDeny returns a boolean if a field has been set.

### GetLists

`func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) GetLists() []ConfigPutRequestAsPathAccessListAsPathAccessListNameListsInner`

GetLists returns the Lists field if non-nil, zero value otherwise.

### GetListsOk

`func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) GetListsOk() (*[]ConfigPutRequestAsPathAccessListAsPathAccessListNameListsInner, bool)`

GetListsOk returns a tuple with the Lists field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLists

`func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) SetLists(v []ConfigPutRequestAsPathAccessListAsPathAccessListNameListsInner)`

SetLists sets Lists field to given value.

### HasLists

`func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) HasLists() bool`

HasLists returns a boolean if a field has been set.

### GetObjectProperties

`func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) GetObjectProperties() ConfigPutRequestIpv4PrefixListIpv4PrefixListNameObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) GetObjectPropertiesOk() (*ConfigPutRequestIpv4PrefixListIpv4PrefixListNameObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) SetObjectProperties(v ConfigPutRequestIpv4PrefixListIpv4PrefixListNameObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


