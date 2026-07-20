# LdapprofilesPutRequestLdapProfileValueAttributeMapsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Enabled** | Pointer to **bool** | Enable this mapping entry | [optional] [default to false]
**MapName** | Pointer to **string** | Category of mapping override | [optional] [default to "attribute"]
**From** | Pointer to **string** | Original RFC2307 attribute or class name to map from | [optional] [default to ""]
**To** | Pointer to **string** | Replacement attribute/class name or value to map to | [optional] [default to ""]
**Index** | Pointer to **int32** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 

## Methods

### NewLdapprofilesPutRequestLdapProfileValueAttributeMapsInner

`func NewLdapprofilesPutRequestLdapProfileValueAttributeMapsInner() *LdapprofilesPutRequestLdapProfileValueAttributeMapsInner`

NewLdapprofilesPutRequestLdapProfileValueAttributeMapsInner instantiates a new LdapprofilesPutRequestLdapProfileValueAttributeMapsInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewLdapprofilesPutRequestLdapProfileValueAttributeMapsInnerWithDefaults

`func NewLdapprofilesPutRequestLdapProfileValueAttributeMapsInnerWithDefaults() *LdapprofilesPutRequestLdapProfileValueAttributeMapsInner`

NewLdapprofilesPutRequestLdapProfileValueAttributeMapsInnerWithDefaults instantiates a new LdapprofilesPutRequestLdapProfileValueAttributeMapsInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEnabled

`func (o *LdapprofilesPutRequestLdapProfileValueAttributeMapsInner) GetEnabled() bool`

GetEnabled returns the Enabled field if non-nil, zero value otherwise.

### GetEnabledOk

`func (o *LdapprofilesPutRequestLdapProfileValueAttributeMapsInner) GetEnabledOk() (*bool, bool)`

GetEnabledOk returns a tuple with the Enabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnabled

`func (o *LdapprofilesPutRequestLdapProfileValueAttributeMapsInner) SetEnabled(v bool)`

SetEnabled sets Enabled field to given value.

### HasEnabled

`func (o *LdapprofilesPutRequestLdapProfileValueAttributeMapsInner) HasEnabled() bool`

HasEnabled returns a boolean if a field has been set.

### GetMapName

`func (o *LdapprofilesPutRequestLdapProfileValueAttributeMapsInner) GetMapName() string`

GetMapName returns the MapName field if non-nil, zero value otherwise.

### GetMapNameOk

`func (o *LdapprofilesPutRequestLdapProfileValueAttributeMapsInner) GetMapNameOk() (*string, bool)`

GetMapNameOk returns a tuple with the MapName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMapName

`func (o *LdapprofilesPutRequestLdapProfileValueAttributeMapsInner) SetMapName(v string)`

SetMapName sets MapName field to given value.

### HasMapName

`func (o *LdapprofilesPutRequestLdapProfileValueAttributeMapsInner) HasMapName() bool`

HasMapName returns a boolean if a field has been set.

### GetFrom

`func (o *LdapprofilesPutRequestLdapProfileValueAttributeMapsInner) GetFrom() string`

GetFrom returns the From field if non-nil, zero value otherwise.

### GetFromOk

`func (o *LdapprofilesPutRequestLdapProfileValueAttributeMapsInner) GetFromOk() (*string, bool)`

GetFromOk returns a tuple with the From field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFrom

`func (o *LdapprofilesPutRequestLdapProfileValueAttributeMapsInner) SetFrom(v string)`

SetFrom sets From field to given value.

### HasFrom

`func (o *LdapprofilesPutRequestLdapProfileValueAttributeMapsInner) HasFrom() bool`

HasFrom returns a boolean if a field has been set.

### GetTo

`func (o *LdapprofilesPutRequestLdapProfileValueAttributeMapsInner) GetTo() string`

GetTo returns the To field if non-nil, zero value otherwise.

### GetToOk

`func (o *LdapprofilesPutRequestLdapProfileValueAttributeMapsInner) GetToOk() (*string, bool)`

GetToOk returns a tuple with the To field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTo

`func (o *LdapprofilesPutRequestLdapProfileValueAttributeMapsInner) SetTo(v string)`

SetTo sets To field to given value.

### HasTo

`func (o *LdapprofilesPutRequestLdapProfileValueAttributeMapsInner) HasTo() bool`

HasTo returns a boolean if a field has been set.

### GetIndex

`func (o *LdapprofilesPutRequestLdapProfileValueAttributeMapsInner) GetIndex() int32`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *LdapprofilesPutRequestLdapProfileValueAttributeMapsInner) GetIndexOk() (*int32, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *LdapprofilesPutRequestLdapProfileValueAttributeMapsInner) SetIndex(v int32)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *LdapprofilesPutRequestLdapProfileValueAttributeMapsInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


