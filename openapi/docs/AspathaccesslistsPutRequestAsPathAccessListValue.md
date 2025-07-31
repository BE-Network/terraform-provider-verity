# AspathaccesslistsPutRequestAsPathAccessListValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**PermitDeny** | Pointer to **string** | Action upon match of Community Strings. | [optional] [default to "permit"]
**Lists** | Pointer to [**[]AspathaccesslistsPutRequestAsPathAccessListValueListsInner**](AspathaccesslistsPutRequestAsPathAccessListValueListsInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**AclsPutRequestIpFilterValueObjectProperties**](AclsPutRequestIpFilterValueObjectProperties.md) |  | [optional] 

## Methods

### NewAspathaccesslistsPutRequestAsPathAccessListValue

`func NewAspathaccesslistsPutRequestAsPathAccessListValue() *AspathaccesslistsPutRequestAsPathAccessListValue`

NewAspathaccesslistsPutRequestAsPathAccessListValue instantiates a new AspathaccesslistsPutRequestAsPathAccessListValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAspathaccesslistsPutRequestAsPathAccessListValueWithDefaults

`func NewAspathaccesslistsPutRequestAsPathAccessListValueWithDefaults() *AspathaccesslistsPutRequestAsPathAccessListValue`

NewAspathaccesslistsPutRequestAsPathAccessListValueWithDefaults instantiates a new AspathaccesslistsPutRequestAsPathAccessListValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *AspathaccesslistsPutRequestAsPathAccessListValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *AspathaccesslistsPutRequestAsPathAccessListValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *AspathaccesslistsPutRequestAsPathAccessListValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *AspathaccesslistsPutRequestAsPathAccessListValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *AspathaccesslistsPutRequestAsPathAccessListValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *AspathaccesslistsPutRequestAsPathAccessListValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *AspathaccesslistsPutRequestAsPathAccessListValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *AspathaccesslistsPutRequestAsPathAccessListValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetPermitDeny

`func (o *AspathaccesslistsPutRequestAsPathAccessListValue) GetPermitDeny() string`

GetPermitDeny returns the PermitDeny field if non-nil, zero value otherwise.

### GetPermitDenyOk

`func (o *AspathaccesslistsPutRequestAsPathAccessListValue) GetPermitDenyOk() (*string, bool)`

GetPermitDenyOk returns a tuple with the PermitDeny field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPermitDeny

`func (o *AspathaccesslistsPutRequestAsPathAccessListValue) SetPermitDeny(v string)`

SetPermitDeny sets PermitDeny field to given value.

### HasPermitDeny

`func (o *AspathaccesslistsPutRequestAsPathAccessListValue) HasPermitDeny() bool`

HasPermitDeny returns a boolean if a field has been set.

### GetLists

`func (o *AspathaccesslistsPutRequestAsPathAccessListValue) GetLists() []AspathaccesslistsPutRequestAsPathAccessListValueListsInner`

GetLists returns the Lists field if non-nil, zero value otherwise.

### GetListsOk

`func (o *AspathaccesslistsPutRequestAsPathAccessListValue) GetListsOk() (*[]AspathaccesslistsPutRequestAsPathAccessListValueListsInner, bool)`

GetListsOk returns a tuple with the Lists field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLists

`func (o *AspathaccesslistsPutRequestAsPathAccessListValue) SetLists(v []AspathaccesslistsPutRequestAsPathAccessListValueListsInner)`

SetLists sets Lists field to given value.

### HasLists

`func (o *AspathaccesslistsPutRequestAsPathAccessListValue) HasLists() bool`

HasLists returns a boolean if a field has been set.

### GetObjectProperties

`func (o *AspathaccesslistsPutRequestAsPathAccessListValue) GetObjectProperties() AclsPutRequestIpFilterValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *AspathaccesslistsPutRequestAsPathAccessListValue) GetObjectPropertiesOk() (*AclsPutRequestIpFilterValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *AspathaccesslistsPutRequestAsPathAccessListValue) SetObjectProperties(v AclsPutRequestIpFilterValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *AspathaccesslistsPutRequestAsPathAccessListValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


