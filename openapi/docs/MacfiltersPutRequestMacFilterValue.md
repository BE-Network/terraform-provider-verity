# MacfiltersPutRequestMacFilterValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Template Name. Must be unique within type. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**Type** | Pointer to **string** | Black vs White MAC Filter | [optional] [default to "White"]
**Filters** | Pointer to [**[]MacfiltersPutRequestMacFilterValueFiltersInner**](MacfiltersPutRequestMacFilterValueFiltersInner.md) |  | [optional] 

## Methods

### NewMacfiltersPutRequestMacFilterValue

`func NewMacfiltersPutRequestMacFilterValue() *MacfiltersPutRequestMacFilterValue`

NewMacfiltersPutRequestMacFilterValue instantiates a new MacfiltersPutRequestMacFilterValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewMacfiltersPutRequestMacFilterValueWithDefaults

`func NewMacfiltersPutRequestMacFilterValueWithDefaults() *MacfiltersPutRequestMacFilterValue`

NewMacfiltersPutRequestMacFilterValueWithDefaults instantiates a new MacfiltersPutRequestMacFilterValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *MacfiltersPutRequestMacFilterValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *MacfiltersPutRequestMacFilterValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *MacfiltersPutRequestMacFilterValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *MacfiltersPutRequestMacFilterValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *MacfiltersPutRequestMacFilterValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *MacfiltersPutRequestMacFilterValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *MacfiltersPutRequestMacFilterValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *MacfiltersPutRequestMacFilterValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetType

`func (o *MacfiltersPutRequestMacFilterValue) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *MacfiltersPutRequestMacFilterValue) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *MacfiltersPutRequestMacFilterValue) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *MacfiltersPutRequestMacFilterValue) HasType() bool`

HasType returns a boolean if a field has been set.

### GetFilters

`func (o *MacfiltersPutRequestMacFilterValue) GetFilters() []MacfiltersPutRequestMacFilterValueFiltersInner`

GetFilters returns the Filters field if non-nil, zero value otherwise.

### GetFiltersOk

`func (o *MacfiltersPutRequestMacFilterValue) GetFiltersOk() (*[]MacfiltersPutRequestMacFilterValueFiltersInner, bool)`

GetFiltersOk returns a tuple with the Filters field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFilters

`func (o *MacfiltersPutRequestMacFilterValue) SetFilters(v []MacfiltersPutRequestMacFilterValueFiltersInner)`

SetFilters sets Filters field to given value.

### HasFilters

`func (o *MacfiltersPutRequestMacFilterValue) HasFilters() bool`

HasFilters returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


