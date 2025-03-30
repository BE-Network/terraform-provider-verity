# ConfigPutRequestExtendedCommunityListExtendedCommunityListNameListsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Enable** | Pointer to **bool** | Enable of this Extended Community List | [optional] [default to false]
**Mode** | Pointer to **string** | Mode | [optional] [default to "route"]
**RouteTargetExpandedExpression** | Pointer to **string** | Match against a BGP extended community of type Route Target | [optional] [default to ""]
**Index** | Pointer to **int32** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 

## Methods

### NewConfigPutRequestExtendedCommunityListExtendedCommunityListNameListsInner

`func NewConfigPutRequestExtendedCommunityListExtendedCommunityListNameListsInner() *ConfigPutRequestExtendedCommunityListExtendedCommunityListNameListsInner`

NewConfigPutRequestExtendedCommunityListExtendedCommunityListNameListsInner instantiates a new ConfigPutRequestExtendedCommunityListExtendedCommunityListNameListsInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestExtendedCommunityListExtendedCommunityListNameListsInnerWithDefaults

`func NewConfigPutRequestExtendedCommunityListExtendedCommunityListNameListsInnerWithDefaults() *ConfigPutRequestExtendedCommunityListExtendedCommunityListNameListsInner`

NewConfigPutRequestExtendedCommunityListExtendedCommunityListNameListsInnerWithDefaults instantiates a new ConfigPutRequestExtendedCommunityListExtendedCommunityListNameListsInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEnable

`func (o *ConfigPutRequestExtendedCommunityListExtendedCommunityListNameListsInner) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestExtendedCommunityListExtendedCommunityListNameListsInner) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestExtendedCommunityListExtendedCommunityListNameListsInner) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestExtendedCommunityListExtendedCommunityListNameListsInner) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetMode

`func (o *ConfigPutRequestExtendedCommunityListExtendedCommunityListNameListsInner) GetMode() string`

GetMode returns the Mode field if non-nil, zero value otherwise.

### GetModeOk

`func (o *ConfigPutRequestExtendedCommunityListExtendedCommunityListNameListsInner) GetModeOk() (*string, bool)`

GetModeOk returns a tuple with the Mode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMode

`func (o *ConfigPutRequestExtendedCommunityListExtendedCommunityListNameListsInner) SetMode(v string)`

SetMode sets Mode field to given value.

### HasMode

`func (o *ConfigPutRequestExtendedCommunityListExtendedCommunityListNameListsInner) HasMode() bool`

HasMode returns a boolean if a field has been set.

### GetRouteTargetExpandedExpression

`func (o *ConfigPutRequestExtendedCommunityListExtendedCommunityListNameListsInner) GetRouteTargetExpandedExpression() string`

GetRouteTargetExpandedExpression returns the RouteTargetExpandedExpression field if non-nil, zero value otherwise.

### GetRouteTargetExpandedExpressionOk

`func (o *ConfigPutRequestExtendedCommunityListExtendedCommunityListNameListsInner) GetRouteTargetExpandedExpressionOk() (*string, bool)`

GetRouteTargetExpandedExpressionOk returns a tuple with the RouteTargetExpandedExpression field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRouteTargetExpandedExpression

`func (o *ConfigPutRequestExtendedCommunityListExtendedCommunityListNameListsInner) SetRouteTargetExpandedExpression(v string)`

SetRouteTargetExpandedExpression sets RouteTargetExpandedExpression field to given value.

### HasRouteTargetExpandedExpression

`func (o *ConfigPutRequestExtendedCommunityListExtendedCommunityListNameListsInner) HasRouteTargetExpandedExpression() bool`

HasRouteTargetExpandedExpression returns a boolean if a field has been set.

### GetIndex

`func (o *ConfigPutRequestExtendedCommunityListExtendedCommunityListNameListsInner) GetIndex() int32`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *ConfigPutRequestExtendedCommunityListExtendedCommunityListNameListsInner) GetIndexOk() (*int32, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *ConfigPutRequestExtendedCommunityListExtendedCommunityListNameListsInner) SetIndex(v int32)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *ConfigPutRequestExtendedCommunityListExtendedCommunityListNameListsInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


