# ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Enable** | Pointer to **bool** | Enable of this Extended Community List | [optional] [default to false]
**Mode** | Pointer to **string** | Mode | [optional] [default to "route"]
**RouteTargetExpandedExpression** | Pointer to **string** | Match against a BGP extended community of type Route Target | [optional] [default to ""]
**Index** | Pointer to **int32** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 

## Methods

### NewExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner

`func NewExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner() *ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner`

NewExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner instantiates a new ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInnerWithDefaults

`func NewExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInnerWithDefaults() *ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner`

NewExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInnerWithDefaults instantiates a new ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEnable

`func (o *ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetMode

`func (o *ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner) GetMode() string`

GetMode returns the Mode field if non-nil, zero value otherwise.

### GetModeOk

`func (o *ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner) GetModeOk() (*string, bool)`

GetModeOk returns a tuple with the Mode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMode

`func (o *ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner) SetMode(v string)`

SetMode sets Mode field to given value.

### HasMode

`func (o *ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner) HasMode() bool`

HasMode returns a boolean if a field has been set.

### GetRouteTargetExpandedExpression

`func (o *ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner) GetRouteTargetExpandedExpression() string`

GetRouteTargetExpandedExpression returns the RouteTargetExpandedExpression field if non-nil, zero value otherwise.

### GetRouteTargetExpandedExpressionOk

`func (o *ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner) GetRouteTargetExpandedExpressionOk() (*string, bool)`

GetRouteTargetExpandedExpressionOk returns a tuple with the RouteTargetExpandedExpression field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRouteTargetExpandedExpression

`func (o *ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner) SetRouteTargetExpandedExpression(v string)`

SetRouteTargetExpandedExpression sets RouteTargetExpandedExpression field to given value.

### HasRouteTargetExpandedExpression

`func (o *ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner) HasRouteTargetExpandedExpression() bool`

HasRouteTargetExpandedExpression returns a boolean if a field has been set.

### GetIndex

`func (o *ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner) GetIndex() int32`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner) GetIndexOk() (*int32, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner) SetIndex(v int32)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *ExtendedcommunitylistsPutRequestExtendedCommunityListValueListsInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


