# ConfigPutRequestCommunityListCommunityListName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**PermitDeny** | Pointer to **string** | Action upon match of Community Strings. | [optional] [default to "permit"]
**AnyAll** | Pointer to **string** | BGP does not advertise any or all routes that do not match the Community String | [optional] [default to "any"]
**StandardExpanded** | Pointer to **string** | Used Community String or Expanded Expression | [optional] [default to "standard"]
**Lists** | Pointer to [**[]ConfigPutRequestCommunityListCommunityListNameListsInner**](ConfigPutRequestCommunityListCommunityListNameListsInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**ConfigPutRequestIpv4PrefixListIpv4PrefixListNameObjectProperties**](ConfigPutRequestIpv4PrefixListIpv4PrefixListNameObjectProperties.md) |  | [optional] 

## Methods

### NewConfigPutRequestCommunityListCommunityListName

`func NewConfigPutRequestCommunityListCommunityListName() *ConfigPutRequestCommunityListCommunityListName`

NewConfigPutRequestCommunityListCommunityListName instantiates a new ConfigPutRequestCommunityListCommunityListName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestCommunityListCommunityListNameWithDefaults

`func NewConfigPutRequestCommunityListCommunityListNameWithDefaults() *ConfigPutRequestCommunityListCommunityListName`

NewConfigPutRequestCommunityListCommunityListNameWithDefaults instantiates a new ConfigPutRequestCommunityListCommunityListName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestCommunityListCommunityListName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestCommunityListCommunityListName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestCommunityListCommunityListName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestCommunityListCommunityListName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ConfigPutRequestCommunityListCommunityListName) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestCommunityListCommunityListName) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestCommunityListCommunityListName) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestCommunityListCommunityListName) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetPermitDeny

`func (o *ConfigPutRequestCommunityListCommunityListName) GetPermitDeny() string`

GetPermitDeny returns the PermitDeny field if non-nil, zero value otherwise.

### GetPermitDenyOk

`func (o *ConfigPutRequestCommunityListCommunityListName) GetPermitDenyOk() (*string, bool)`

GetPermitDenyOk returns a tuple with the PermitDeny field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPermitDeny

`func (o *ConfigPutRequestCommunityListCommunityListName) SetPermitDeny(v string)`

SetPermitDeny sets PermitDeny field to given value.

### HasPermitDeny

`func (o *ConfigPutRequestCommunityListCommunityListName) HasPermitDeny() bool`

HasPermitDeny returns a boolean if a field has been set.

### GetAnyAll

`func (o *ConfigPutRequestCommunityListCommunityListName) GetAnyAll() string`

GetAnyAll returns the AnyAll field if non-nil, zero value otherwise.

### GetAnyAllOk

`func (o *ConfigPutRequestCommunityListCommunityListName) GetAnyAllOk() (*string, bool)`

GetAnyAllOk returns a tuple with the AnyAll field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnyAll

`func (o *ConfigPutRequestCommunityListCommunityListName) SetAnyAll(v string)`

SetAnyAll sets AnyAll field to given value.

### HasAnyAll

`func (o *ConfigPutRequestCommunityListCommunityListName) HasAnyAll() bool`

HasAnyAll returns a boolean if a field has been set.

### GetStandardExpanded

`func (o *ConfigPutRequestCommunityListCommunityListName) GetStandardExpanded() string`

GetStandardExpanded returns the StandardExpanded field if non-nil, zero value otherwise.

### GetStandardExpandedOk

`func (o *ConfigPutRequestCommunityListCommunityListName) GetStandardExpandedOk() (*string, bool)`

GetStandardExpandedOk returns a tuple with the StandardExpanded field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStandardExpanded

`func (o *ConfigPutRequestCommunityListCommunityListName) SetStandardExpanded(v string)`

SetStandardExpanded sets StandardExpanded field to given value.

### HasStandardExpanded

`func (o *ConfigPutRequestCommunityListCommunityListName) HasStandardExpanded() bool`

HasStandardExpanded returns a boolean if a field has been set.

### GetLists

`func (o *ConfigPutRequestCommunityListCommunityListName) GetLists() []ConfigPutRequestCommunityListCommunityListNameListsInner`

GetLists returns the Lists field if non-nil, zero value otherwise.

### GetListsOk

`func (o *ConfigPutRequestCommunityListCommunityListName) GetListsOk() (*[]ConfigPutRequestCommunityListCommunityListNameListsInner, bool)`

GetListsOk returns a tuple with the Lists field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLists

`func (o *ConfigPutRequestCommunityListCommunityListName) SetLists(v []ConfigPutRequestCommunityListCommunityListNameListsInner)`

SetLists sets Lists field to given value.

### HasLists

`func (o *ConfigPutRequestCommunityListCommunityListName) HasLists() bool`

HasLists returns a boolean if a field has been set.

### GetObjectProperties

`func (o *ConfigPutRequestCommunityListCommunityListName) GetObjectProperties() ConfigPutRequestIpv4PrefixListIpv4PrefixListNameObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ConfigPutRequestCommunityListCommunityListName) GetObjectPropertiesOk() (*ConfigPutRequestIpv4PrefixListIpv4PrefixListNameObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ConfigPutRequestCommunityListCommunityListName) SetObjectProperties(v ConfigPutRequestIpv4PrefixListIpv4PrefixListNameObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ConfigPutRequestCommunityListCommunityListName) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


