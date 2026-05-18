# MacfiltersPutRequestMacFilterValueFiltersInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**FilterNumMac** | Pointer to **string** | MAC address descriptor including colons example 01:23:45:67:9a:ab. and * notation accepted example 12:* | [optional] [default to ""]
**FilterNumMask** | Pointer to **string** | Hexidecimal mask including colons example ff:ff:fe:00:00:00. /n and * notation accepted example /16 or 12:* | [optional] [default to ""]
**FilterNumEnable** | Pointer to **bool** | Enable of this MAC Filter  | [optional] [default to false]
**Index** | Pointer to **int32** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 

## Methods

### NewMacfiltersPutRequestMacFilterValueFiltersInner

`func NewMacfiltersPutRequestMacFilterValueFiltersInner() *MacfiltersPutRequestMacFilterValueFiltersInner`

NewMacfiltersPutRequestMacFilterValueFiltersInner instantiates a new MacfiltersPutRequestMacFilterValueFiltersInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMacfiltersPutRequestMacFilterValueFiltersInnerWithDefaults

`func NewMacfiltersPutRequestMacFilterValueFiltersInnerWithDefaults() *MacfiltersPutRequestMacFilterValueFiltersInner`

NewMacfiltersPutRequestMacFilterValueFiltersInnerWithDefaults instantiates a new MacfiltersPutRequestMacFilterValueFiltersInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetFilterNumMac

`func (o *MacfiltersPutRequestMacFilterValueFiltersInner) GetFilterNumMac() string`

GetFilterNumMac returns the FilterNumMac field if non-nil, zero value otherwise.

### GetFilterNumMacOk

`func (o *MacfiltersPutRequestMacFilterValueFiltersInner) GetFilterNumMacOk() (*string, bool)`

GetFilterNumMacOk returns a tuple with the FilterNumMac field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFilterNumMac

`func (o *MacfiltersPutRequestMacFilterValueFiltersInner) SetFilterNumMac(v string)`

SetFilterNumMac sets FilterNumMac field to given value.

### HasFilterNumMac

`func (o *MacfiltersPutRequestMacFilterValueFiltersInner) HasFilterNumMac() bool`

HasFilterNumMac returns a boolean if a field has been set.

### GetFilterNumMask

`func (o *MacfiltersPutRequestMacFilterValueFiltersInner) GetFilterNumMask() string`

GetFilterNumMask returns the FilterNumMask field if non-nil, zero value otherwise.

### GetFilterNumMaskOk

`func (o *MacfiltersPutRequestMacFilterValueFiltersInner) GetFilterNumMaskOk() (*string, bool)`

GetFilterNumMaskOk returns a tuple with the FilterNumMask field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFilterNumMask

`func (o *MacfiltersPutRequestMacFilterValueFiltersInner) SetFilterNumMask(v string)`

SetFilterNumMask sets FilterNumMask field to given value.

### HasFilterNumMask

`func (o *MacfiltersPutRequestMacFilterValueFiltersInner) HasFilterNumMask() bool`

HasFilterNumMask returns a boolean if a field has been set.

### GetFilterNumEnable

`func (o *MacfiltersPutRequestMacFilterValueFiltersInner) GetFilterNumEnable() bool`

GetFilterNumEnable returns the FilterNumEnable field if non-nil, zero value otherwise.

### GetFilterNumEnableOk

`func (o *MacfiltersPutRequestMacFilterValueFiltersInner) GetFilterNumEnableOk() (*bool, bool)`

GetFilterNumEnableOk returns a tuple with the FilterNumEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFilterNumEnable

`func (o *MacfiltersPutRequestMacFilterValueFiltersInner) SetFilterNumEnable(v bool)`

SetFilterNumEnable sets FilterNumEnable field to given value.

### HasFilterNumEnable

`func (o *MacfiltersPutRequestMacFilterValueFiltersInner) HasFilterNumEnable() bool`

HasFilterNumEnable returns a boolean if a field has been set.

### GetIndex

`func (o *MacfiltersPutRequestMacFilterValueFiltersInner) GetIndex() int32`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *MacfiltersPutRequestMacFilterValueFiltersInner) GetIndexOk() (*int32, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *MacfiltersPutRequestMacFilterValueFiltersInner) SetIndex(v int32)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *MacfiltersPutRequestMacFilterValueFiltersInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


