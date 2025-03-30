# ConfigPutRequestAsPathAccessListAsPathAccessListNameListsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Enable** | Pointer to **bool** | Enable this AS Path Access List | [optional] [default to false]
**RegularExpression** | Pointer to **string** | Regular Expression to match BGP Community Strings | [optional] [default to ""]
**Index** | Pointer to **int32** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 

## Methods

### NewConfigPutRequestAsPathAccessListAsPathAccessListNameListsInner

`func NewConfigPutRequestAsPathAccessListAsPathAccessListNameListsInner() *ConfigPutRequestAsPathAccessListAsPathAccessListNameListsInner`

NewConfigPutRequestAsPathAccessListAsPathAccessListNameListsInner instantiates a new ConfigPutRequestAsPathAccessListAsPathAccessListNameListsInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestAsPathAccessListAsPathAccessListNameListsInnerWithDefaults

`func NewConfigPutRequestAsPathAccessListAsPathAccessListNameListsInnerWithDefaults() *ConfigPutRequestAsPathAccessListAsPathAccessListNameListsInner`

NewConfigPutRequestAsPathAccessListAsPathAccessListNameListsInnerWithDefaults instantiates a new ConfigPutRequestAsPathAccessListAsPathAccessListNameListsInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEnable

`func (o *ConfigPutRequestAsPathAccessListAsPathAccessListNameListsInner) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestAsPathAccessListAsPathAccessListNameListsInner) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestAsPathAccessListAsPathAccessListNameListsInner) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestAsPathAccessListAsPathAccessListNameListsInner) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetRegularExpression

`func (o *ConfigPutRequestAsPathAccessListAsPathAccessListNameListsInner) GetRegularExpression() string`

GetRegularExpression returns the RegularExpression field if non-nil, zero value otherwise.

### GetRegularExpressionOk

`func (o *ConfigPutRequestAsPathAccessListAsPathAccessListNameListsInner) GetRegularExpressionOk() (*string, bool)`

GetRegularExpressionOk returns a tuple with the RegularExpression field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegularExpression

`func (o *ConfigPutRequestAsPathAccessListAsPathAccessListNameListsInner) SetRegularExpression(v string)`

SetRegularExpression sets RegularExpression field to given value.

### HasRegularExpression

`func (o *ConfigPutRequestAsPathAccessListAsPathAccessListNameListsInner) HasRegularExpression() bool`

HasRegularExpression returns a boolean if a field has been set.

### GetIndex

`func (o *ConfigPutRequestAsPathAccessListAsPathAccessListNameListsInner) GetIndex() int32`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *ConfigPutRequestAsPathAccessListAsPathAccessListNameListsInner) GetIndexOk() (*int32, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *ConfigPutRequestAsPathAccessListAsPathAccessListNameListsInner) SetIndex(v int32)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *ConfigPutRequestAsPathAccessListAsPathAccessListNameListsInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


