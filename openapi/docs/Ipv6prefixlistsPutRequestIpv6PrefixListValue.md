# Ipv6prefixlistsPutRequestIpv6PrefixListValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**Lists** | Pointer to [**[]Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner**](Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**AclsPutRequestIpFilterValueObjectProperties**](AclsPutRequestIpFilterValueObjectProperties.md) |  | [optional] 

## Methods

### NewIpv6prefixlistsPutRequestIpv6PrefixListValue

`func NewIpv6prefixlistsPutRequestIpv6PrefixListValue() *Ipv6prefixlistsPutRequestIpv6PrefixListValue`

NewIpv6prefixlistsPutRequestIpv6PrefixListValue instantiates a new Ipv6prefixlistsPutRequestIpv6PrefixListValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewIpv6prefixlistsPutRequestIpv6PrefixListValueWithDefaults

`func NewIpv6prefixlistsPutRequestIpv6PrefixListValueWithDefaults() *Ipv6prefixlistsPutRequestIpv6PrefixListValue`

NewIpv6prefixlistsPutRequestIpv6PrefixListValueWithDefaults instantiates a new Ipv6prefixlistsPutRequestIpv6PrefixListValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetLists

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValue) GetLists() []Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner`

GetLists returns the Lists field if non-nil, zero value otherwise.

### GetListsOk

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValue) GetListsOk() (*[]Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner, bool)`

GetListsOk returns a tuple with the Lists field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLists

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValue) SetLists(v []Ipv6prefixlistsPutRequestIpv6PrefixListValueListsInner)`

SetLists sets Lists field to given value.

### HasLists

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValue) HasLists() bool`

HasLists returns a boolean if a field has been set.

### GetObjectProperties

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValue) GetObjectProperties() AclsPutRequestIpFilterValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValue) GetObjectPropertiesOk() (*AclsPutRequestIpFilterValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValue) SetObjectProperties(v AclsPutRequestIpFilterValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *Ipv6prefixlistsPutRequestIpv6PrefixListValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


