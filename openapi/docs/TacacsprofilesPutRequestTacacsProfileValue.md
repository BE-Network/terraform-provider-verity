# TacacsprofilesPutRequestTacacsProfileValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Template Name. Must be unique within type. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**TacacsServers** | Pointer to [**[]TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner**](TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner.md) |  | [optional] 

## Methods

### NewTacacsprofilesPutRequestTacacsProfileValue

`func NewTacacsprofilesPutRequestTacacsProfileValue() *TacacsprofilesPutRequestTacacsProfileValue`

NewTacacsprofilesPutRequestTacacsProfileValue instantiates a new TacacsprofilesPutRequestTacacsProfileValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTacacsprofilesPutRequestTacacsProfileValueWithDefaults

`func NewTacacsprofilesPutRequestTacacsProfileValueWithDefaults() *TacacsprofilesPutRequestTacacsProfileValue`

NewTacacsprofilesPutRequestTacacsProfileValueWithDefaults instantiates a new TacacsprofilesPutRequestTacacsProfileValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *TacacsprofilesPutRequestTacacsProfileValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *TacacsprofilesPutRequestTacacsProfileValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *TacacsprofilesPutRequestTacacsProfileValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *TacacsprofilesPutRequestTacacsProfileValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *TacacsprofilesPutRequestTacacsProfileValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *TacacsprofilesPutRequestTacacsProfileValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *TacacsprofilesPutRequestTacacsProfileValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *TacacsprofilesPutRequestTacacsProfileValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetTacacsServers

`func (o *TacacsprofilesPutRequestTacacsProfileValue) GetTacacsServers() []TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner`

GetTacacsServers returns the TacacsServers field if non-nil, zero value otherwise.

### GetTacacsServersOk

`func (o *TacacsprofilesPutRequestTacacsProfileValue) GetTacacsServersOk() (*[]TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner, bool)`

GetTacacsServersOk returns a tuple with the TacacsServers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTacacsServers

`func (o *TacacsprofilesPutRequestTacacsProfileValue) SetTacacsServers(v []TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner)`

SetTacacsServers sets TacacsServers field to given value.

### HasTacacsServers

`func (o *TacacsprofilesPutRequestTacacsProfileValue) HasTacacsServers() bool`

HasTacacsServers returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


