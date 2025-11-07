# SwitchpointsPutRequestSwitchpointValueEthsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Breakout** | Pointer to **string** | Breakout Port Override. Available options determined by Switch capability, Installed SFP and the capacity of the pipeline. | [optional] [default to ""]
**Index** | Pointer to **int32** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 
**EthNumIcon** | Pointer to **string** | Icon of this Eth Port | [optional] [default to "empty"]
**EthNumLabel** | Pointer to **string** | Label of this Eth Port | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable port.  | [optional] [default to true]

## Methods

### NewSwitchpointsPutRequestSwitchpointValueEthsInner

`func NewSwitchpointsPutRequestSwitchpointValueEthsInner() *SwitchpointsPutRequestSwitchpointValueEthsInner`

NewSwitchpointsPutRequestSwitchpointValueEthsInner instantiates a new SwitchpointsPutRequestSwitchpointValueEthsInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSwitchpointsPutRequestSwitchpointValueEthsInnerWithDefaults

`func NewSwitchpointsPutRequestSwitchpointValueEthsInnerWithDefaults() *SwitchpointsPutRequestSwitchpointValueEthsInner`

NewSwitchpointsPutRequestSwitchpointValueEthsInnerWithDefaults instantiates a new SwitchpointsPutRequestSwitchpointValueEthsInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBreakout

`func (o *SwitchpointsPutRequestSwitchpointValueEthsInner) GetBreakout() string`

GetBreakout returns the Breakout field if non-nil, zero value otherwise.

### GetBreakoutOk

`func (o *SwitchpointsPutRequestSwitchpointValueEthsInner) GetBreakoutOk() (*string, bool)`

GetBreakoutOk returns a tuple with the Breakout field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBreakout

`func (o *SwitchpointsPutRequestSwitchpointValueEthsInner) SetBreakout(v string)`

SetBreakout sets Breakout field to given value.

### HasBreakout

`func (o *SwitchpointsPutRequestSwitchpointValueEthsInner) HasBreakout() bool`

HasBreakout returns a boolean if a field has been set.

### GetIndex

`func (o *SwitchpointsPutRequestSwitchpointValueEthsInner) GetIndex() int32`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *SwitchpointsPutRequestSwitchpointValueEthsInner) GetIndexOk() (*int32, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *SwitchpointsPutRequestSwitchpointValueEthsInner) SetIndex(v int32)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *SwitchpointsPutRequestSwitchpointValueEthsInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.

### GetEthNumIcon

`func (o *SwitchpointsPutRequestSwitchpointValueEthsInner) GetEthNumIcon() string`

GetEthNumIcon returns the EthNumIcon field if non-nil, zero value otherwise.

### GetEthNumIconOk

`func (o *SwitchpointsPutRequestSwitchpointValueEthsInner) GetEthNumIconOk() (*string, bool)`

GetEthNumIconOk returns a tuple with the EthNumIcon field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEthNumIcon

`func (o *SwitchpointsPutRequestSwitchpointValueEthsInner) SetEthNumIcon(v string)`

SetEthNumIcon sets EthNumIcon field to given value.

### HasEthNumIcon

`func (o *SwitchpointsPutRequestSwitchpointValueEthsInner) HasEthNumIcon() bool`

HasEthNumIcon returns a boolean if a field has been set.

### GetEthNumLabel

`func (o *SwitchpointsPutRequestSwitchpointValueEthsInner) GetEthNumLabel() string`

GetEthNumLabel returns the EthNumLabel field if non-nil, zero value otherwise.

### GetEthNumLabelOk

`func (o *SwitchpointsPutRequestSwitchpointValueEthsInner) GetEthNumLabelOk() (*string, bool)`

GetEthNumLabelOk returns a tuple with the EthNumLabel field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEthNumLabel

`func (o *SwitchpointsPutRequestSwitchpointValueEthsInner) SetEthNumLabel(v string)`

SetEthNumLabel sets EthNumLabel field to given value.

### HasEthNumLabel

`func (o *SwitchpointsPutRequestSwitchpointValueEthsInner) HasEthNumLabel() bool`

HasEthNumLabel returns a boolean if a field has been set.

### GetEnable

`func (o *SwitchpointsPutRequestSwitchpointValueEthsInner) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *SwitchpointsPutRequestSwitchpointValueEthsInner) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *SwitchpointsPutRequestSwitchpointValueEthsInner) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *SwitchpointsPutRequestSwitchpointValueEthsInner) HasEnable() bool`

HasEnable returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


