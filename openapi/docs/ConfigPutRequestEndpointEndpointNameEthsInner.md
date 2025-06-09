# ConfigPutRequestEndpointEndpointNameEthsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Breakout** | Pointer to **string** | Breakout Port Override. Available options determined by Switch capability, Installed SFP and the capacity of the pipeline. | [optional] [default to ""]
**Index** | Pointer to **int32** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 

## Methods

### NewConfigPutRequestEndpointEndpointNameEthsInner

`func NewConfigPutRequestEndpointEndpointNameEthsInner() *ConfigPutRequestEndpointEndpointNameEthsInner`

NewConfigPutRequestEndpointEndpointNameEthsInner instantiates a new ConfigPutRequestEndpointEndpointNameEthsInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestEndpointEndpointNameEthsInnerWithDefaults

`func NewConfigPutRequestEndpointEndpointNameEthsInnerWithDefaults() *ConfigPutRequestEndpointEndpointNameEthsInner`

NewConfigPutRequestEndpointEndpointNameEthsInnerWithDefaults instantiates a new ConfigPutRequestEndpointEndpointNameEthsInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBreakout

`func (o *ConfigPutRequestEndpointEndpointNameEthsInner) GetBreakout() string`

GetBreakout returns the Breakout field if non-nil, zero value otherwise.

### GetBreakoutOk

`func (o *ConfigPutRequestEndpointEndpointNameEthsInner) GetBreakoutOk() (*string, bool)`

GetBreakoutOk returns a tuple with the Breakout field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBreakout

`func (o *ConfigPutRequestEndpointEndpointNameEthsInner) SetBreakout(v string)`

SetBreakout sets Breakout field to given value.

### HasBreakout

`func (o *ConfigPutRequestEndpointEndpointNameEthsInner) HasBreakout() bool`

HasBreakout returns a boolean if a field has been set.

### GetIndex

`func (o *ConfigPutRequestEndpointEndpointNameEthsInner) GetIndex() int32`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *ConfigPutRequestEndpointEndpointNameEthsInner) GetIndexOk() (*int32, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *ConfigPutRequestEndpointEndpointNameEthsInner) SetIndex(v int32)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *ConfigPutRequestEndpointEndpointNameEthsInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


