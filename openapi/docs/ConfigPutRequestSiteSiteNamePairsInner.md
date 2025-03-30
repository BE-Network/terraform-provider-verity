# ConfigPutRequestSiteSiteNamePairsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Switchpoint1** | Pointer to **string** | Switchpoint | [optional] [default to ""]
**Switchpoint1RefType** | Pointer to **string** | Object type for switchpoint_1 field | [optional] 
**Switchpoint2** | Pointer to **string** | Switchpoint | [optional] [default to ""]
**Switchpoint2RefType** | Pointer to **string** | Object type for switchpoint_2 field | [optional] 
**LagGroup** | Pointer to **string** | LAG Group | [optional] [default to ""]
**LagGroupRefType** | Pointer to **string** | Object type for lag_group field | [optional] 
**IsWhiteboxPair** | Pointer to **bool** | LAG Pair | [optional] [default to false]
**Index** | Pointer to **int32** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 

## Methods

### NewConfigPutRequestSiteSiteNamePairsInner

`func NewConfigPutRequestSiteSiteNamePairsInner() *ConfigPutRequestSiteSiteNamePairsInner`

NewConfigPutRequestSiteSiteNamePairsInner instantiates a new ConfigPutRequestSiteSiteNamePairsInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestSiteSiteNamePairsInnerWithDefaults

`func NewConfigPutRequestSiteSiteNamePairsInnerWithDefaults() *ConfigPutRequestSiteSiteNamePairsInner`

NewConfigPutRequestSiteSiteNamePairsInnerWithDefaults instantiates a new ConfigPutRequestSiteSiteNamePairsInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestSiteSiteNamePairsInner) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestSiteSiteNamePairsInner) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestSiteSiteNamePairsInner) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestSiteSiteNamePairsInner) HasName() bool`

HasName returns a boolean if a field has been set.

### GetSwitchpoint1

`func (o *ConfigPutRequestSiteSiteNamePairsInner) GetSwitchpoint1() string`

GetSwitchpoint1 returns the Switchpoint1 field if non-nil, zero value otherwise.

### GetSwitchpoint1Ok

`func (o *ConfigPutRequestSiteSiteNamePairsInner) GetSwitchpoint1Ok() (*string, bool)`

GetSwitchpoint1Ok returns a tuple with the Switchpoint1 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchpoint1

`func (o *ConfigPutRequestSiteSiteNamePairsInner) SetSwitchpoint1(v string)`

SetSwitchpoint1 sets Switchpoint1 field to given value.

### HasSwitchpoint1

`func (o *ConfigPutRequestSiteSiteNamePairsInner) HasSwitchpoint1() bool`

HasSwitchpoint1 returns a boolean if a field has been set.

### GetSwitchpoint1RefType

`func (o *ConfigPutRequestSiteSiteNamePairsInner) GetSwitchpoint1RefType() string`

GetSwitchpoint1RefType returns the Switchpoint1RefType field if non-nil, zero value otherwise.

### GetSwitchpoint1RefTypeOk

`func (o *ConfigPutRequestSiteSiteNamePairsInner) GetSwitchpoint1RefTypeOk() (*string, bool)`

GetSwitchpoint1RefTypeOk returns a tuple with the Switchpoint1RefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchpoint1RefType

`func (o *ConfigPutRequestSiteSiteNamePairsInner) SetSwitchpoint1RefType(v string)`

SetSwitchpoint1RefType sets Switchpoint1RefType field to given value.

### HasSwitchpoint1RefType

`func (o *ConfigPutRequestSiteSiteNamePairsInner) HasSwitchpoint1RefType() bool`

HasSwitchpoint1RefType returns a boolean if a field has been set.

### GetSwitchpoint2

`func (o *ConfigPutRequestSiteSiteNamePairsInner) GetSwitchpoint2() string`

GetSwitchpoint2 returns the Switchpoint2 field if non-nil, zero value otherwise.

### GetSwitchpoint2Ok

`func (o *ConfigPutRequestSiteSiteNamePairsInner) GetSwitchpoint2Ok() (*string, bool)`

GetSwitchpoint2Ok returns a tuple with the Switchpoint2 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchpoint2

`func (o *ConfigPutRequestSiteSiteNamePairsInner) SetSwitchpoint2(v string)`

SetSwitchpoint2 sets Switchpoint2 field to given value.

### HasSwitchpoint2

`func (o *ConfigPutRequestSiteSiteNamePairsInner) HasSwitchpoint2() bool`

HasSwitchpoint2 returns a boolean if a field has been set.

### GetSwitchpoint2RefType

`func (o *ConfigPutRequestSiteSiteNamePairsInner) GetSwitchpoint2RefType() string`

GetSwitchpoint2RefType returns the Switchpoint2RefType field if non-nil, zero value otherwise.

### GetSwitchpoint2RefTypeOk

`func (o *ConfigPutRequestSiteSiteNamePairsInner) GetSwitchpoint2RefTypeOk() (*string, bool)`

GetSwitchpoint2RefTypeOk returns a tuple with the Switchpoint2RefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchpoint2RefType

`func (o *ConfigPutRequestSiteSiteNamePairsInner) SetSwitchpoint2RefType(v string)`

SetSwitchpoint2RefType sets Switchpoint2RefType field to given value.

### HasSwitchpoint2RefType

`func (o *ConfigPutRequestSiteSiteNamePairsInner) HasSwitchpoint2RefType() bool`

HasSwitchpoint2RefType returns a boolean if a field has been set.

### GetLagGroup

`func (o *ConfigPutRequestSiteSiteNamePairsInner) GetLagGroup() string`

GetLagGroup returns the LagGroup field if non-nil, zero value otherwise.

### GetLagGroupOk

`func (o *ConfigPutRequestSiteSiteNamePairsInner) GetLagGroupOk() (*string, bool)`

GetLagGroupOk returns a tuple with the LagGroup field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLagGroup

`func (o *ConfigPutRequestSiteSiteNamePairsInner) SetLagGroup(v string)`

SetLagGroup sets LagGroup field to given value.

### HasLagGroup

`func (o *ConfigPutRequestSiteSiteNamePairsInner) HasLagGroup() bool`

HasLagGroup returns a boolean if a field has been set.

### GetLagGroupRefType

`func (o *ConfigPutRequestSiteSiteNamePairsInner) GetLagGroupRefType() string`

GetLagGroupRefType returns the LagGroupRefType field if non-nil, zero value otherwise.

### GetLagGroupRefTypeOk

`func (o *ConfigPutRequestSiteSiteNamePairsInner) GetLagGroupRefTypeOk() (*string, bool)`

GetLagGroupRefTypeOk returns a tuple with the LagGroupRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLagGroupRefType

`func (o *ConfigPutRequestSiteSiteNamePairsInner) SetLagGroupRefType(v string)`

SetLagGroupRefType sets LagGroupRefType field to given value.

### HasLagGroupRefType

`func (o *ConfigPutRequestSiteSiteNamePairsInner) HasLagGroupRefType() bool`

HasLagGroupRefType returns a boolean if a field has been set.

### GetIsWhiteboxPair

`func (o *ConfigPutRequestSiteSiteNamePairsInner) GetIsWhiteboxPair() bool`

GetIsWhiteboxPair returns the IsWhiteboxPair field if non-nil, zero value otherwise.

### GetIsWhiteboxPairOk

`func (o *ConfigPutRequestSiteSiteNamePairsInner) GetIsWhiteboxPairOk() (*bool, bool)`

GetIsWhiteboxPairOk returns a tuple with the IsWhiteboxPair field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsWhiteboxPair

`func (o *ConfigPutRequestSiteSiteNamePairsInner) SetIsWhiteboxPair(v bool)`

SetIsWhiteboxPair sets IsWhiteboxPair field to given value.

### HasIsWhiteboxPair

`func (o *ConfigPutRequestSiteSiteNamePairsInner) HasIsWhiteboxPair() bool`

HasIsWhiteboxPair returns a boolean if a field has been set.

### GetIndex

`func (o *ConfigPutRequestSiteSiteNamePairsInner) GetIndex() int32`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *ConfigPutRequestSiteSiteNamePairsInner) GetIndexOk() (*int32, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *ConfigPutRequestSiteSiteNamePairsInner) SetIndex(v int32)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *ConfigPutRequestSiteSiteNamePairsInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


