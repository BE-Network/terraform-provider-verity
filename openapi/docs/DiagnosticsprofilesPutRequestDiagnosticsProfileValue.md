# DiagnosticsprofilesPutRequestDiagnosticsProfileValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Template Name. Must be unique within type. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**EnableSflow** | Pointer to **bool** | Enable sFlow for this Diagnostics Profile  | [optional] [default to false]
**UseInternalCollector** | Pointer to **bool** | Use the internal Collector as the flow collector | [optional] [default to false]
**FlowCollector** | Pointer to **string** | Flow Collector for this Diagnostics Profile  | [optional] [default to ""]
**FlowCollectorRefType** | Pointer to **string** | Object type for flow_collector field | [optional] 
**PollInterval** | Pointer to **NullableInt64** | The sampling rate for sFlow polling (seconds) | [optional] [default to 20]
**VrfType** | Pointer to **string** | Management or Underlay | [optional] [default to "management"]
**MonitoringAcl** | Pointer to **string** | Monitoring ACL whose Service VLANs are mirrored to the ERSPAN destination | [optional] [default to ""]
**MonitoringAclRefType** | Pointer to **string** | Object type for monitoring_acl field | [optional] 
**ErspanDestinationIp** | Pointer to **string** | IPv4 address of the remote ERSPAN collector | [optional] [default to ""]
**ErspanDscp** | Pointer to **NullableInt64** | DSCP value for ERSPAN packets (0-63) | [optional] 
**ErspanTtl** | Pointer to **NullableInt64** | Time-to-live value for ERSPAN packets (0-255) | [optional] 
**ErspanGreType** | Pointer to **string** | GRE protocol type as a 0x-prefixed hexadecimal value | [optional] [default to ""]
**ErspanQueue** | Pointer to **NullableInt64** | Output queue for ERSPAN packets (0-63) | [optional] 

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

### GetUseInternalCollector

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) GetUseInternalCollector() bool`

GetUseInternalCollector returns the UseInternalCollector field if non-nil, zero value otherwise.

### GetUseInternalCollectorOk

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) GetUseInternalCollectorOk() (*bool, bool)`

GetUseInternalCollectorOk returns a tuple with the UseInternalCollector field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUseInternalCollector

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) SetUseInternalCollector(v bool)`

SetUseInternalCollector sets UseInternalCollector field to given value.

### HasUseInternalCollector

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) HasUseInternalCollector() bool`

HasUseInternalCollector returns a boolean if a field has been set.

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

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) GetPollInterval() int64`

GetPollInterval returns the PollInterval field if non-nil, zero value otherwise.

### GetPollIntervalOk

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) GetPollIntervalOk() (*int64, bool)`

GetPollIntervalOk returns a tuple with the PollInterval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPollInterval

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) SetPollInterval(v int64)`

SetPollInterval sets PollInterval field to given value.

### HasPollInterval

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) HasPollInterval() bool`

HasPollInterval returns a boolean if a field has been set.

### SetPollIntervalNil

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) SetPollIntervalNil(b bool)`

 SetPollIntervalNil sets the value for PollInterval to be an explicit nil

### UnsetPollInterval
`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) UnsetPollInterval()`

UnsetPollInterval ensures that no value is present for PollInterval, not even an explicit nil
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

### GetMonitoringAcl

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) GetMonitoringAcl() string`

GetMonitoringAcl returns the MonitoringAcl field if non-nil, zero value otherwise.

### GetMonitoringAclOk

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) GetMonitoringAclOk() (*string, bool)`

GetMonitoringAclOk returns a tuple with the MonitoringAcl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMonitoringAcl

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) SetMonitoringAcl(v string)`

SetMonitoringAcl sets MonitoringAcl field to given value.

### HasMonitoringAcl

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) HasMonitoringAcl() bool`

HasMonitoringAcl returns a boolean if a field has been set.

### GetMonitoringAclRefType

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) GetMonitoringAclRefType() string`

GetMonitoringAclRefType returns the MonitoringAclRefType field if non-nil, zero value otherwise.

### GetMonitoringAclRefTypeOk

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) GetMonitoringAclRefTypeOk() (*string, bool)`

GetMonitoringAclRefTypeOk returns a tuple with the MonitoringAclRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMonitoringAclRefType

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) SetMonitoringAclRefType(v string)`

SetMonitoringAclRefType sets MonitoringAclRefType field to given value.

### HasMonitoringAclRefType

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) HasMonitoringAclRefType() bool`

HasMonitoringAclRefType returns a boolean if a field has been set.

### GetErspanDestinationIp

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) GetErspanDestinationIp() string`

GetErspanDestinationIp returns the ErspanDestinationIp field if non-nil, zero value otherwise.

### GetErspanDestinationIpOk

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) GetErspanDestinationIpOk() (*string, bool)`

GetErspanDestinationIpOk returns a tuple with the ErspanDestinationIp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErspanDestinationIp

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) SetErspanDestinationIp(v string)`

SetErspanDestinationIp sets ErspanDestinationIp field to given value.

### HasErspanDestinationIp

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) HasErspanDestinationIp() bool`

HasErspanDestinationIp returns a boolean if a field has been set.

### GetErspanDscp

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) GetErspanDscp() int64`

GetErspanDscp returns the ErspanDscp field if non-nil, zero value otherwise.

### GetErspanDscpOk

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) GetErspanDscpOk() (*int64, bool)`

GetErspanDscpOk returns a tuple with the ErspanDscp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErspanDscp

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) SetErspanDscp(v int64)`

SetErspanDscp sets ErspanDscp field to given value.

### HasErspanDscp

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) HasErspanDscp() bool`

HasErspanDscp returns a boolean if a field has been set.

### SetErspanDscpNil

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) SetErspanDscpNil(b bool)`

 SetErspanDscpNil sets the value for ErspanDscp to be an explicit nil

### UnsetErspanDscp
`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) UnsetErspanDscp()`

UnsetErspanDscp ensures that no value is present for ErspanDscp, not even an explicit nil
### GetErspanTtl

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) GetErspanTtl() int64`

GetErspanTtl returns the ErspanTtl field if non-nil, zero value otherwise.

### GetErspanTtlOk

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) GetErspanTtlOk() (*int64, bool)`

GetErspanTtlOk returns a tuple with the ErspanTtl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErspanTtl

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) SetErspanTtl(v int64)`

SetErspanTtl sets ErspanTtl field to given value.

### HasErspanTtl

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) HasErspanTtl() bool`

HasErspanTtl returns a boolean if a field has been set.

### SetErspanTtlNil

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) SetErspanTtlNil(b bool)`

 SetErspanTtlNil sets the value for ErspanTtl to be an explicit nil

### UnsetErspanTtl
`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) UnsetErspanTtl()`

UnsetErspanTtl ensures that no value is present for ErspanTtl, not even an explicit nil
### GetErspanGreType

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) GetErspanGreType() string`

GetErspanGreType returns the ErspanGreType field if non-nil, zero value otherwise.

### GetErspanGreTypeOk

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) GetErspanGreTypeOk() (*string, bool)`

GetErspanGreTypeOk returns a tuple with the ErspanGreType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErspanGreType

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) SetErspanGreType(v string)`

SetErspanGreType sets ErspanGreType field to given value.

### HasErspanGreType

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) HasErspanGreType() bool`

HasErspanGreType returns a boolean if a field has been set.

### GetErspanQueue

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) GetErspanQueue() int64`

GetErspanQueue returns the ErspanQueue field if non-nil, zero value otherwise.

### GetErspanQueueOk

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) GetErspanQueueOk() (*int64, bool)`

GetErspanQueueOk returns a tuple with the ErspanQueue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetErspanQueue

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) SetErspanQueue(v int64)`

SetErspanQueue sets ErspanQueue field to given value.

### HasErspanQueue

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) HasErspanQueue() bool`

HasErspanQueue returns a boolean if a field has been set.

### SetErspanQueueNil

`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) SetErspanQueueNil(b bool)`

 SetErspanQueueNil sets the value for ErspanQueue to be an explicit nil

### UnsetErspanQueue
`func (o *DiagnosticsprofilesPutRequestDiagnosticsProfileValue) UnsetErspanQueue()`

UnsetErspanQueue ensures that no value is present for ErspanQueue, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


