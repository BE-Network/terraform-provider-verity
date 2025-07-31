# Ipv4prefixlistsPutRequestIpv4PrefixListValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**Lists** | Pointer to [**[]Ipv4prefixlistsPutRequestIpv4PrefixListValueListsInner**](Ipv4prefixlistsPutRequestIpv4PrefixListValueListsInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**AclsPutRequestIpFilterValueObjectProperties**](AclsPutRequestIpFilterValueObjectProperties.md) |  | [optional] 

## Methods

### NewIpv4prefixlistsPutRequestIpv4PrefixListValue

`func NewIpv4prefixlistsPutRequestIpv4PrefixListValue() *Ipv4prefixlistsPutRequestIpv4PrefixListValue`

NewIpv4prefixlistsPutRequestIpv4PrefixListValue instantiates a new Ipv4prefixlistsPutRequestIpv4PrefixListValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewIpv4prefixlistsPutRequestIpv4PrefixListValueWithDefaults

`func NewIpv4prefixlistsPutRequestIpv4PrefixListValueWithDefaults() *Ipv4prefixlistsPutRequestIpv4PrefixListValue`

NewIpv4prefixlistsPutRequestIpv4PrefixListValueWithDefaults instantiates a new Ipv4prefixlistsPutRequestIpv4PrefixListValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *Ipv4prefixlistsPutRequestIpv4PrefixListValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *Ipv4prefixlistsPutRequestIpv4PrefixListValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *Ipv4prefixlistsPutRequestIpv4PrefixListValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *Ipv4prefixlistsPutRequestIpv4PrefixListValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *Ipv4prefixlistsPutRequestIpv4PrefixListValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *Ipv4prefixlistsPutRequestIpv4PrefixListValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *Ipv4prefixlistsPutRequestIpv4PrefixListValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *Ipv4prefixlistsPutRequestIpv4PrefixListValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetLists

`func (o *Ipv4prefixlistsPutRequestIpv4PrefixListValue) GetLists() []Ipv4prefixlistsPutRequestIpv4PrefixListValueListsInner`

GetLists returns the Lists field if non-nil, zero value otherwise.

### GetListsOk

`func (o *Ipv4prefixlistsPutRequestIpv4PrefixListValue) GetListsOk() (*[]Ipv4prefixlistsPutRequestIpv4PrefixListValueListsInner, bool)`

GetListsOk returns a tuple with the Lists field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLists

`func (o *Ipv4prefixlistsPutRequestIpv4PrefixListValue) SetLists(v []Ipv4prefixlistsPutRequestIpv4PrefixListValueListsInner)`

SetLists sets Lists field to given value.

### HasLists

`func (o *Ipv4prefixlistsPutRequestIpv4PrefixListValue) HasLists() bool`

HasLists returns a boolean if a field has been set.

### GetObjectProperties

`func (o *Ipv4prefixlistsPutRequestIpv4PrefixListValue) GetObjectProperties() AclsPutRequestIpFilterValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *Ipv4prefixlistsPutRequestIpv4PrefixListValue) GetObjectPropertiesOk() (*AclsPutRequestIpFilterValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *Ipv4prefixlistsPutRequestIpv4PrefixListValue) SetObjectProperties(v AclsPutRequestIpFilterValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *Ipv4prefixlistsPutRequestIpv4PrefixListValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


