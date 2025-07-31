# CommunitylistsPutRequestCommunityListValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**PermitDeny** | Pointer to **string** | Action upon match of Community Strings. | [optional] [default to "permit"]
**AnyAll** | Pointer to **string** | BGP does not advertise any or all routes that do not match the Community String | [optional] [default to "any"]
**StandardExpanded** | Pointer to **string** | Used Community String or Expanded Expression | [optional] [default to "standard"]
**Lists** | Pointer to [**[]CommunitylistsPutRequestCommunityListValueListsInner**](CommunitylistsPutRequestCommunityListValueListsInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**AclsPutRequestIpFilterValueObjectProperties**](AclsPutRequestIpFilterValueObjectProperties.md) |  | [optional] 

## Methods

### NewCommunitylistsPutRequestCommunityListValue

`func NewCommunitylistsPutRequestCommunityListValue() *CommunitylistsPutRequestCommunityListValue`

NewCommunitylistsPutRequestCommunityListValue instantiates a new CommunitylistsPutRequestCommunityListValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewCommunitylistsPutRequestCommunityListValueWithDefaults

`func NewCommunitylistsPutRequestCommunityListValueWithDefaults() *CommunitylistsPutRequestCommunityListValue`

NewCommunitylistsPutRequestCommunityListValueWithDefaults instantiates a new CommunitylistsPutRequestCommunityListValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *CommunitylistsPutRequestCommunityListValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *CommunitylistsPutRequestCommunityListValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *CommunitylistsPutRequestCommunityListValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *CommunitylistsPutRequestCommunityListValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *CommunitylistsPutRequestCommunityListValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *CommunitylistsPutRequestCommunityListValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *CommunitylistsPutRequestCommunityListValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *CommunitylistsPutRequestCommunityListValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetPermitDeny

`func (o *CommunitylistsPutRequestCommunityListValue) GetPermitDeny() string`

GetPermitDeny returns the PermitDeny field if non-nil, zero value otherwise.

### GetPermitDenyOk

`func (o *CommunitylistsPutRequestCommunityListValue) GetPermitDenyOk() (*string, bool)`

GetPermitDenyOk returns a tuple with the PermitDeny field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPermitDeny

`func (o *CommunitylistsPutRequestCommunityListValue) SetPermitDeny(v string)`

SetPermitDeny sets PermitDeny field to given value.

### HasPermitDeny

`func (o *CommunitylistsPutRequestCommunityListValue) HasPermitDeny() bool`

HasPermitDeny returns a boolean if a field has been set.

### GetAnyAll

`func (o *CommunitylistsPutRequestCommunityListValue) GetAnyAll() string`

GetAnyAll returns the AnyAll field if non-nil, zero value otherwise.

### GetAnyAllOk

`func (o *CommunitylistsPutRequestCommunityListValue) GetAnyAllOk() (*string, bool)`

GetAnyAllOk returns a tuple with the AnyAll field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnyAll

`func (o *CommunitylistsPutRequestCommunityListValue) SetAnyAll(v string)`

SetAnyAll sets AnyAll field to given value.

### HasAnyAll

`func (o *CommunitylistsPutRequestCommunityListValue) HasAnyAll() bool`

HasAnyAll returns a boolean if a field has been set.

### GetStandardExpanded

`func (o *CommunitylistsPutRequestCommunityListValue) GetStandardExpanded() string`

GetStandardExpanded returns the StandardExpanded field if non-nil, zero value otherwise.

### GetStandardExpandedOk

`func (o *CommunitylistsPutRequestCommunityListValue) GetStandardExpandedOk() (*string, bool)`

GetStandardExpandedOk returns a tuple with the StandardExpanded field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStandardExpanded

`func (o *CommunitylistsPutRequestCommunityListValue) SetStandardExpanded(v string)`

SetStandardExpanded sets StandardExpanded field to given value.

### HasStandardExpanded

`func (o *CommunitylistsPutRequestCommunityListValue) HasStandardExpanded() bool`

HasStandardExpanded returns a boolean if a field has been set.

### GetLists

`func (o *CommunitylistsPutRequestCommunityListValue) GetLists() []CommunitylistsPutRequestCommunityListValueListsInner`

GetLists returns the Lists field if non-nil, zero value otherwise.

### GetListsOk

`func (o *CommunitylistsPutRequestCommunityListValue) GetListsOk() (*[]CommunitylistsPutRequestCommunityListValueListsInner, bool)`

GetListsOk returns a tuple with the Lists field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLists

`func (o *CommunitylistsPutRequestCommunityListValue) SetLists(v []CommunitylistsPutRequestCommunityListValueListsInner)`

SetLists sets Lists field to given value.

### HasLists

`func (o *CommunitylistsPutRequestCommunityListValue) HasLists() bool`

HasLists returns a boolean if a field has been set.

### GetObjectProperties

`func (o *CommunitylistsPutRequestCommunityListValue) GetObjectProperties() AclsPutRequestIpFilterValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *CommunitylistsPutRequestCommunityListValue) GetObjectPropertiesOk() (*AclsPutRequestIpFilterValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *CommunitylistsPutRequestCommunityListValue) SetObjectProperties(v AclsPutRequestIpFilterValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *CommunitylistsPutRequestCommunityListValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


