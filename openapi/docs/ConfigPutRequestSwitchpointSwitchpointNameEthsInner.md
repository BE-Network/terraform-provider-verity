# ConfigPutRequestSwitchpointSwitchpointNameEthsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Breakout** | Pointer to **string** | Breakout Port Override. Available options determined by Switch capability, Installed SFP and the capacity of the pipeline. | [optional] [default to ""]
**Index** | Pointer to **int32** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 

## Methods

### NewConfigPutRequestSwitchpointSwitchpointNameEthsInner

`func NewConfigPutRequestSwitchpointSwitchpointNameEthsInner() *ConfigPutRequestSwitchpointSwitchpointNameEthsInner`

NewConfigPutRequestSwitchpointSwitchpointNameEthsInner instantiates a new ConfigPutRequestSwitchpointSwitchpointNameEthsInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestSwitchpointSwitchpointNameEthsInnerWithDefaults

`func NewConfigPutRequestSwitchpointSwitchpointNameEthsInnerWithDefaults() *ConfigPutRequestSwitchpointSwitchpointNameEthsInner`

NewConfigPutRequestSwitchpointSwitchpointNameEthsInnerWithDefaults instantiates a new ConfigPutRequestSwitchpointSwitchpointNameEthsInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBreakout

`func (o *ConfigPutRequestSwitchpointSwitchpointNameEthsInner) GetBreakout() string`

GetBreakout returns the Breakout field if non-nil, zero value otherwise.

### GetBreakoutOk

`func (o *ConfigPutRequestSwitchpointSwitchpointNameEthsInner) GetBreakoutOk() (*string, bool)`

GetBreakoutOk returns a tuple with the Breakout field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBreakout

`func (o *ConfigPutRequestSwitchpointSwitchpointNameEthsInner) SetBreakout(v string)`

SetBreakout sets Breakout field to given value.

### HasBreakout

`func (o *ConfigPutRequestSwitchpointSwitchpointNameEthsInner) HasBreakout() bool`

HasBreakout returns a boolean if a field has been set.

### GetIndex

`func (o *ConfigPutRequestSwitchpointSwitchpointNameEthsInner) GetIndex() int32`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *ConfigPutRequestSwitchpointSwitchpointNameEthsInner) GetIndexOk() (*int32, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *ConfigPutRequestSwitchpointSwitchpointNameEthsInner) SetIndex(v int32)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *ConfigPutRequestSwitchpointSwitchpointNameEthsInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


