# RoutemapsPutRequestRouteMapValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**RouteMapClauses** | Pointer to [**[]RoutemapsPutRequestRouteMapValueRouteMapClausesInner**](RoutemapsPutRequestRouteMapValueRouteMapClausesInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**AclsPutRequestIpFilterValueObjectProperties**](AclsPutRequestIpFilterValueObjectProperties.md) |  | [optional] 

## Methods

### NewRoutemapsPutRequestRouteMapValue

`func NewRoutemapsPutRequestRouteMapValue() *RoutemapsPutRequestRouteMapValue`

NewRoutemapsPutRequestRouteMapValue instantiates a new RoutemapsPutRequestRouteMapValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRoutemapsPutRequestRouteMapValueWithDefaults

`func NewRoutemapsPutRequestRouteMapValueWithDefaults() *RoutemapsPutRequestRouteMapValue`

NewRoutemapsPutRequestRouteMapValueWithDefaults instantiates a new RoutemapsPutRequestRouteMapValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *RoutemapsPutRequestRouteMapValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *RoutemapsPutRequestRouteMapValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *RoutemapsPutRequestRouteMapValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *RoutemapsPutRequestRouteMapValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *RoutemapsPutRequestRouteMapValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *RoutemapsPutRequestRouteMapValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *RoutemapsPutRequestRouteMapValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *RoutemapsPutRequestRouteMapValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetRouteMapClauses

`func (o *RoutemapsPutRequestRouteMapValue) GetRouteMapClauses() []RoutemapsPutRequestRouteMapValueRouteMapClausesInner`

GetRouteMapClauses returns the RouteMapClauses field if non-nil, zero value otherwise.

### GetRouteMapClausesOk

`func (o *RoutemapsPutRequestRouteMapValue) GetRouteMapClausesOk() (*[]RoutemapsPutRequestRouteMapValueRouteMapClausesInner, bool)`

GetRouteMapClausesOk returns a tuple with the RouteMapClauses field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRouteMapClauses

`func (o *RoutemapsPutRequestRouteMapValue) SetRouteMapClauses(v []RoutemapsPutRequestRouteMapValueRouteMapClausesInner)`

SetRouteMapClauses sets RouteMapClauses field to given value.

### HasRouteMapClauses

`func (o *RoutemapsPutRequestRouteMapValue) HasRouteMapClauses() bool`

HasRouteMapClauses returns a boolean if a field has been set.

### GetObjectProperties

`func (o *RoutemapsPutRequestRouteMapValue) GetObjectProperties() AclsPutRequestIpFilterValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *RoutemapsPutRequestRouteMapValue) GetObjectPropertiesOk() (*AclsPutRequestIpFilterValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *RoutemapsPutRequestRouteMapValue) SetObjectProperties(v AclsPutRequestIpFilterValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *RoutemapsPutRequestRouteMapValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


