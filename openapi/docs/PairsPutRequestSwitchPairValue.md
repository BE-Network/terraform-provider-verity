# PairsPutRequestSwitchPairValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Template Name. Must be unique within type. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to true]
**Switchpoint1** | Pointer to **string** | Switchpoint | [optional] [default to ""]
**Switchpoint1RefType** | Pointer to **string** | Object type for switchpoint_1 field | [optional] 
**Switchpoint2** | Pointer to **string** | Switchpoint | [optional] [default to ""]
**Switchpoint2RefType** | Pointer to **string** | Object type for switchpoint_2 field | [optional] 
**Lag** | Pointer to **string** | LAG | [optional] [default to ""]
**LagRefType** | Pointer to **string** | Object type for lag field | [optional] 
**IsWhiteboxPair** | Pointer to **bool** | Is Whitebox Pair | [optional] [default to false]

## Methods

### NewPairsPutRequestSwitchPairValue

`func NewPairsPutRequestSwitchPairValue() *PairsPutRequestSwitchPairValue`

NewPairsPutRequestSwitchPairValue instantiates a new PairsPutRequestSwitchPairValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPairsPutRequestSwitchPairValueWithDefaults

`func NewPairsPutRequestSwitchPairValueWithDefaults() *PairsPutRequestSwitchPairValue`

NewPairsPutRequestSwitchPairValueWithDefaults instantiates a new PairsPutRequestSwitchPairValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *PairsPutRequestSwitchPairValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *PairsPutRequestSwitchPairValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *PairsPutRequestSwitchPairValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *PairsPutRequestSwitchPairValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *PairsPutRequestSwitchPairValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *PairsPutRequestSwitchPairValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *PairsPutRequestSwitchPairValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *PairsPutRequestSwitchPairValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetSwitchpoint1

`func (o *PairsPutRequestSwitchPairValue) GetSwitchpoint1() string`

GetSwitchpoint1 returns the Switchpoint1 field if non-nil, zero value otherwise.

### GetSwitchpoint1Ok

`func (o *PairsPutRequestSwitchPairValue) GetSwitchpoint1Ok() (*string, bool)`

GetSwitchpoint1Ok returns a tuple with the Switchpoint1 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchpoint1

`func (o *PairsPutRequestSwitchPairValue) SetSwitchpoint1(v string)`

SetSwitchpoint1 sets Switchpoint1 field to given value.

### HasSwitchpoint1

`func (o *PairsPutRequestSwitchPairValue) HasSwitchpoint1() bool`

HasSwitchpoint1 returns a boolean if a field has been set.

### GetSwitchpoint1RefType

`func (o *PairsPutRequestSwitchPairValue) GetSwitchpoint1RefType() string`

GetSwitchpoint1RefType returns the Switchpoint1RefType field if non-nil, zero value otherwise.

### GetSwitchpoint1RefTypeOk

`func (o *PairsPutRequestSwitchPairValue) GetSwitchpoint1RefTypeOk() (*string, bool)`

GetSwitchpoint1RefTypeOk returns a tuple with the Switchpoint1RefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchpoint1RefType

`func (o *PairsPutRequestSwitchPairValue) SetSwitchpoint1RefType(v string)`

SetSwitchpoint1RefType sets Switchpoint1RefType field to given value.

### HasSwitchpoint1RefType

`func (o *PairsPutRequestSwitchPairValue) HasSwitchpoint1RefType() bool`

HasSwitchpoint1RefType returns a boolean if a field has been set.

### GetSwitchpoint2

`func (o *PairsPutRequestSwitchPairValue) GetSwitchpoint2() string`

GetSwitchpoint2 returns the Switchpoint2 field if non-nil, zero value otherwise.

### GetSwitchpoint2Ok

`func (o *PairsPutRequestSwitchPairValue) GetSwitchpoint2Ok() (*string, bool)`

GetSwitchpoint2Ok returns a tuple with the Switchpoint2 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchpoint2

`func (o *PairsPutRequestSwitchPairValue) SetSwitchpoint2(v string)`

SetSwitchpoint2 sets Switchpoint2 field to given value.

### HasSwitchpoint2

`func (o *PairsPutRequestSwitchPairValue) HasSwitchpoint2() bool`

HasSwitchpoint2 returns a boolean if a field has been set.

### GetSwitchpoint2RefType

`func (o *PairsPutRequestSwitchPairValue) GetSwitchpoint2RefType() string`

GetSwitchpoint2RefType returns the Switchpoint2RefType field if non-nil, zero value otherwise.

### GetSwitchpoint2RefTypeOk

`func (o *PairsPutRequestSwitchPairValue) GetSwitchpoint2RefTypeOk() (*string, bool)`

GetSwitchpoint2RefTypeOk returns a tuple with the Switchpoint2RefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchpoint2RefType

`func (o *PairsPutRequestSwitchPairValue) SetSwitchpoint2RefType(v string)`

SetSwitchpoint2RefType sets Switchpoint2RefType field to given value.

### HasSwitchpoint2RefType

`func (o *PairsPutRequestSwitchPairValue) HasSwitchpoint2RefType() bool`

HasSwitchpoint2RefType returns a boolean if a field has been set.

### GetLag

`func (o *PairsPutRequestSwitchPairValue) GetLag() string`

GetLag returns the Lag field if non-nil, zero value otherwise.

### GetLagOk

`func (o *PairsPutRequestSwitchPairValue) GetLagOk() (*string, bool)`

GetLagOk returns a tuple with the Lag field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLag

`func (o *PairsPutRequestSwitchPairValue) SetLag(v string)`

SetLag sets Lag field to given value.

### HasLag

`func (o *PairsPutRequestSwitchPairValue) HasLag() bool`

HasLag returns a boolean if a field has been set.

### GetLagRefType

`func (o *PairsPutRequestSwitchPairValue) GetLagRefType() string`

GetLagRefType returns the LagRefType field if non-nil, zero value otherwise.

### GetLagRefTypeOk

`func (o *PairsPutRequestSwitchPairValue) GetLagRefTypeOk() (*string, bool)`

GetLagRefTypeOk returns a tuple with the LagRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLagRefType

`func (o *PairsPutRequestSwitchPairValue) SetLagRefType(v string)`

SetLagRefType sets LagRefType field to given value.

### HasLagRefType

`func (o *PairsPutRequestSwitchPairValue) HasLagRefType() bool`

HasLagRefType returns a boolean if a field has been set.

### GetIsWhiteboxPair

`func (o *PairsPutRequestSwitchPairValue) GetIsWhiteboxPair() bool`

GetIsWhiteboxPair returns the IsWhiteboxPair field if non-nil, zero value otherwise.

### GetIsWhiteboxPairOk

`func (o *PairsPutRequestSwitchPairValue) GetIsWhiteboxPairOk() (*bool, bool)`

GetIsWhiteboxPairOk returns a tuple with the IsWhiteboxPair field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsWhiteboxPair

`func (o *PairsPutRequestSwitchPairValue) SetIsWhiteboxPair(v bool)`

SetIsWhiteboxPair sets IsWhiteboxPair field to given value.

### HasIsWhiteboxPair

`func (o *PairsPutRequestSwitchPairValue) HasIsWhiteboxPair() bool`

HasIsWhiteboxPair returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


