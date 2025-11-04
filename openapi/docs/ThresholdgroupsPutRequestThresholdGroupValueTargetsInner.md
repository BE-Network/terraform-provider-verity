# ThresholdgroupsPutRequestThresholdGroupValueTargetsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Enable** | Pointer to **bool** | Enable | [optional] [default to false]
**Type** | Pointer to **string** | Specific element or Grouping Rules to apply thresholds to | [optional] [default to "grouping_rules"]
**GroupingRules** | Pointer to **string** | Elements to apply thresholds to | [optional] [default to ""]
**GroupingRulesRefType** | Pointer to **string** | Object type for grouping_rules field | [optional] 
**Switchpoint** | Pointer to **string** | Switchpoint to apply thresholds to | [optional] [default to ""]
**SwitchpointRefType** | Pointer to **string** | Object type for switchpoint field | [optional] 
**Port** | Pointer to **string** | Port to apply thresholds to | [optional] [default to ""]
**Index** | Pointer to **int32** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 

## Methods

### NewThresholdgroupsPutRequestThresholdGroupValueTargetsInner

`func NewThresholdgroupsPutRequestThresholdGroupValueTargetsInner() *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner`

NewThresholdgroupsPutRequestThresholdGroupValueTargetsInner instantiates a new ThresholdgroupsPutRequestThresholdGroupValueTargetsInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewThresholdgroupsPutRequestThresholdGroupValueTargetsInnerWithDefaults

`func NewThresholdgroupsPutRequestThresholdGroupValueTargetsInnerWithDefaults() *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner`

NewThresholdgroupsPutRequestThresholdGroupValueTargetsInnerWithDefaults instantiates a new ThresholdgroupsPutRequestThresholdGroupValueTargetsInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEnable

`func (o *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetType

`func (o *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner) HasType() bool`

HasType returns a boolean if a field has been set.

### GetGroupingRules

`func (o *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner) GetGroupingRules() string`

GetGroupingRules returns the GroupingRules field if non-nil, zero value otherwise.

### GetGroupingRulesOk

`func (o *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner) GetGroupingRulesOk() (*string, bool)`

GetGroupingRulesOk returns a tuple with the GroupingRules field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGroupingRules

`func (o *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner) SetGroupingRules(v string)`

SetGroupingRules sets GroupingRules field to given value.

### HasGroupingRules

`func (o *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner) HasGroupingRules() bool`

HasGroupingRules returns a boolean if a field has been set.

### GetGroupingRulesRefType

`func (o *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner) GetGroupingRulesRefType() string`

GetGroupingRulesRefType returns the GroupingRulesRefType field if non-nil, zero value otherwise.

### GetGroupingRulesRefTypeOk

`func (o *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner) GetGroupingRulesRefTypeOk() (*string, bool)`

GetGroupingRulesRefTypeOk returns a tuple with the GroupingRulesRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGroupingRulesRefType

`func (o *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner) SetGroupingRulesRefType(v string)`

SetGroupingRulesRefType sets GroupingRulesRefType field to given value.

### HasGroupingRulesRefType

`func (o *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner) HasGroupingRulesRefType() bool`

HasGroupingRulesRefType returns a boolean if a field has been set.

### GetSwitchpoint

`func (o *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner) GetSwitchpoint() string`

GetSwitchpoint returns the Switchpoint field if non-nil, zero value otherwise.

### GetSwitchpointOk

`func (o *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner) GetSwitchpointOk() (*string, bool)`

GetSwitchpointOk returns a tuple with the Switchpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchpoint

`func (o *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner) SetSwitchpoint(v string)`

SetSwitchpoint sets Switchpoint field to given value.

### HasSwitchpoint

`func (o *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner) HasSwitchpoint() bool`

HasSwitchpoint returns a boolean if a field has been set.

### GetSwitchpointRefType

`func (o *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner) GetSwitchpointRefType() string`

GetSwitchpointRefType returns the SwitchpointRefType field if non-nil, zero value otherwise.

### GetSwitchpointRefTypeOk

`func (o *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner) GetSwitchpointRefTypeOk() (*string, bool)`

GetSwitchpointRefTypeOk returns a tuple with the SwitchpointRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchpointRefType

`func (o *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner) SetSwitchpointRefType(v string)`

SetSwitchpointRefType sets SwitchpointRefType field to given value.

### HasSwitchpointRefType

`func (o *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner) HasSwitchpointRefType() bool`

HasSwitchpointRefType returns a boolean if a field has been set.

### GetPort

`func (o *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner) GetPort() string`

GetPort returns the Port field if non-nil, zero value otherwise.

### GetPortOk

`func (o *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner) GetPortOk() (*string, bool)`

GetPortOk returns a tuple with the Port field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPort

`func (o *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner) SetPort(v string)`

SetPort sets Port field to given value.

### HasPort

`func (o *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner) HasPort() bool`

HasPort returns a boolean if a field has been set.

### GetIndex

`func (o *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner) GetIndex() int32`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner) GetIndexOk() (*int32, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner) SetIndex(v int32)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *ThresholdgroupsPutRequestThresholdGroupValueTargetsInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


