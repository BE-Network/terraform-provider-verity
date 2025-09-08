# DiagnosticsprofilesPutRequestDiagnosticsProfileValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**EnableSflow** | Pointer to **bool** | Enable sFlow for this Diagnostics Profile  | [optional] [default to false]
**FlowCollector** | Pointer to **string** | Flow Collector for this Diagnostics Profile  | [optional] [default to ""]
**FlowCollectorRefType** | Pointer to **string** | Object type for flow_collector field | [optional] 
**PollInterval** | Pointer to **int32** | The sampling rate for sFlow polling (seconds) | [optional] [default to 20]
**VrfType** | Pointer to **string** | Management or Underlay | [optional] [default to "management"]

## Methods

### NewDiagnosticsprofilesPutRequestDiagnosticsProfileValue

`func NewDiagnosticsprofilesPutRequestDiagnosticsProfileValue() *DiagnosticsprofilesPutRequestDiagnosticsProfileValue`

NewDiagnosticsprofilesPutRequestDiagnosticsProfileValue instantiates a new DiagnosticsprofilesPutRequestDiagnosticsProfileValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDiagnosticsprofilesPutRequestDiagnosticsProfileValueWithDefaults

`func NewDiagnosticsprofilesPutRequestDiagnosticsProfileValueWithDefaults() *DiagnosticsprofilesPutRequestDiagnosticsProfileValue`

NewDiagnosticsprofilesPutRequestDiagnosticsProfileValueWithDefaults instantiates a new DiagnosticsprofilesPutRequestDiagnosticsProfileValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetEnableSflow

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) GetEnableSflow() bool`

GetEnableSflow returns the EnableSflow field if non-nil, zero value otherwise.

### GetEnableSflowOk

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) GetEnableSflowOk() (*bool, bool)`

GetEnableSflowOk returns a tuple with the EnableSflow field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnableSflow

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) SetEnableSflow(v bool)`

SetEnableSflow sets EnableSflow field to given value.

### HasEnableSflow

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) HasEnableSflow() bool`

HasEnableSflow returns a boolean if a field has been set.

### GetFlowCollector

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) GetFlowCollector() string`

GetFlowCollector returns the FlowCollector field if non-nil, zero value otherwise.

### GetFlowCollectorOk

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) GetFlowCollectorOk() (*string, bool)`

GetFlowCollectorOk returns a tuple with the FlowCollector field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFlowCollector

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) SetFlowCollector(v string)`

SetFlowCollector sets FlowCollector field to given value.

### HasFlowCollector

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) HasFlowCollector() bool`

HasFlowCollector returns a boolean if a field has been set.

### GetFlowCollectorRefType

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) GetFlowCollectorRefType() string`

GetFlowCollectorRefType returns the FlowCollectorRefType field if non-nil, zero value otherwise.

### GetFlowCollectorRefTypeOk

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) GetFlowCollectorRefTypeOk() (*string, bool)`

GetFlowCollectorRefTypeOk returns a tuple with the FlowCollectorRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFlowCollectorRefType

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) SetFlowCollectorRefType(v string)`

SetFlowCollectorRefType sets FlowCollectorRefType field to given value.

### HasFlowCollectorRefType

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) HasFlowCollectorRefType() bool`

HasFlowCollectorRefType returns a boolean if a field has been set.

### GetPollInterval

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) GetPollInterval() int32`

GetPollInterval returns the PollInterval field if non-nil, zero value otherwise.

### GetPollIntervalOk

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) GetPollIntervalOk() (*int32, bool)`

GetPollIntervalOk returns a tuple with the PollInterval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPollInterval

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) SetPollInterval(v int32)`

SetPollInterval sets PollInterval field to given value.

### HasPollInterval

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) HasPollInterval() bool`

HasPollInterval returns a boolean if a field has been set.

### GetVrfType

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) GetVrfType() string`

GetVrfType returns the VrfType field if non-nil, zero value otherwise.

### GetVrfTypeOk

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) GetVrfTypeOk() (*string, bool)`

GetVrfTypeOk returns a tuple with the VrfType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVrfType

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) SetVrfType(v string)`

SetVrfType sets VrfType field to given value.

### HasVrfType

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) HasVrfType() bool`

HasVrfType returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


