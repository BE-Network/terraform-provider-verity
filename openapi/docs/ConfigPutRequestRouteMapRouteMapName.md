# ConfigPutRequestRouteMapRouteMapName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**RouteMapClauses** | Pointer to [**[]ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner**](ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties**](ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties.md) |  | [optional] 

## Methods

### NewConfigPutRequestRouteMapRouteMapName

`func NewConfigPutRequestRouteMapRouteMapName() *ConfigPutRequestRouteMapRouteMapName`

NewConfigPutRequestRouteMapRouteMapName instantiates a new ConfigPutRequestRouteMapRouteMapName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestRouteMapRouteMapNameWithDefaults

`func NewConfigPutRequestRouteMapRouteMapNameWithDefaults() *ConfigPutRequestRouteMapRouteMapName`

NewConfigPutRequestRouteMapRouteMapNameWithDefaults instantiates a new ConfigPutRequestRouteMapRouteMapName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestRouteMapRouteMapName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestRouteMapRouteMapName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestRouteMapRouteMapName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestRouteMapRouteMapName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ConfigPutRequestRouteMapRouteMapName) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestRouteMapRouteMapName) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestRouteMapRouteMapName) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestRouteMapRouteMapName) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetRouteMapClauses

`func (o *ConfigPutRequestRouteMapRouteMapName) GetRouteMapClauses() []ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner`

GetRouteMapClauses returns the RouteMapClauses field if non-nil, zero value otherwise.

### GetRouteMapClausesOk

`func (o *ConfigPutRequestRouteMapRouteMapName) GetRouteMapClausesOk() (*[]ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner, bool)`

GetRouteMapClausesOk returns a tuple with the RouteMapClauses field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRouteMapClauses

`func (o *ConfigPutRequestRouteMapRouteMapName) SetRouteMapClauses(v []ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner)`

SetRouteMapClauses sets RouteMapClauses field to given value.

### HasRouteMapClauses

`func (o *ConfigPutRequestRouteMapRouteMapName) HasRouteMapClauses() bool`

HasRouteMapClauses returns a boolean if a field has been set.

### GetObjectProperties

`func (o *ConfigPutRequestRouteMapRouteMapName) GetObjectProperties() ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ConfigPutRequestRouteMapRouteMapName) GetObjectPropertiesOk() (*ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ConfigPutRequestRouteMapRouteMapName) SetObjectProperties(v ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ConfigPutRequestRouteMapRouteMapName) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


