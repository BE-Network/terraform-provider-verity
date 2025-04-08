# ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**TrafficMirrorNumEnable** | Pointer to **bool** | Enable Traffic Mirror | [optional] [default to false]
**TrafficMirrorNumSourcePort** | Pointer to **string** | Source Port for Traffic Mirror | [optional] [default to ""]
**TrafficMirrorNumSourceLagIndicator** | Pointer to **bool** | Source LAG Indicator for Traffic Mirror | [optional] [default to false]
**TrafficMirrorNumDestinationPort** | Pointer to **string** | Destination Port for Traffic Mirror | [optional] [default to ""]
**TrafficMirrorNumInboundTraffic** | Pointer to **bool** | Boolean value indicating if the mirror is for inbound traffic | [optional] [default to false]
**TrafficMirrorNumOutboundTraffic** | Pointer to **bool** | Boolean value indicating if the mirror is for outbound traffic | [optional] [default to false]

## Methods

### NewConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner

`func NewConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner() *ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner`

NewConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner instantiates a new ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInnerWithDefaults

`func NewConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInnerWithDefaults() *ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner`

NewConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInnerWithDefaults instantiates a new ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetTrafficMirrorNumEnable

`func (o *ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner) GetTrafficMirrorNumEnable() bool`

GetTrafficMirrorNumEnable returns the TrafficMirrorNumEnable field if non-nil, zero value otherwise.

### GetTrafficMirrorNumEnableOk

`func (o *ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner) GetTrafficMirrorNumEnableOk() (*bool, bool)`

GetTrafficMirrorNumEnableOk returns a tuple with the TrafficMirrorNumEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTrafficMirrorNumEnable

`func (o *ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner) SetTrafficMirrorNumEnable(v bool)`

SetTrafficMirrorNumEnable sets TrafficMirrorNumEnable field to given value.

### HasTrafficMirrorNumEnable

`func (o *ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner) HasTrafficMirrorNumEnable() bool`

HasTrafficMirrorNumEnable returns a boolean if a field has been set.

### GetTrafficMirrorNumSourcePort

`func (o *ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner) GetTrafficMirrorNumSourcePort() string`

GetTrafficMirrorNumSourcePort returns the TrafficMirrorNumSourcePort field if non-nil, zero value otherwise.

### GetTrafficMirrorNumSourcePortOk

`func (o *ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner) GetTrafficMirrorNumSourcePortOk() (*string, bool)`

GetTrafficMirrorNumSourcePortOk returns a tuple with the TrafficMirrorNumSourcePort field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTrafficMirrorNumSourcePort

`func (o *ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner) SetTrafficMirrorNumSourcePort(v string)`

SetTrafficMirrorNumSourcePort sets TrafficMirrorNumSourcePort field to given value.

### HasTrafficMirrorNumSourcePort

`func (o *ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner) HasTrafficMirrorNumSourcePort() bool`

HasTrafficMirrorNumSourcePort returns a boolean if a field has been set.

### GetTrafficMirrorNumSourceLagIndicator

`func (o *ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner) GetTrafficMirrorNumSourceLagIndicator() bool`

GetTrafficMirrorNumSourceLagIndicator returns the TrafficMirrorNumSourceLagIndicator field if non-nil, zero value otherwise.

### GetTrafficMirrorNumSourceLagIndicatorOk

`func (o *ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner) GetTrafficMirrorNumSourceLagIndicatorOk() (*bool, bool)`

GetTrafficMirrorNumSourceLagIndicatorOk returns a tuple with the TrafficMirrorNumSourceLagIndicator field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTrafficMirrorNumSourceLagIndicator

`func (o *ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner) SetTrafficMirrorNumSourceLagIndicator(v bool)`

SetTrafficMirrorNumSourceLagIndicator sets TrafficMirrorNumSourceLagIndicator field to given value.

### HasTrafficMirrorNumSourceLagIndicator

`func (o *ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner) HasTrafficMirrorNumSourceLagIndicator() bool`

HasTrafficMirrorNumSourceLagIndicator returns a boolean if a field has been set.

### GetTrafficMirrorNumDestinationPort

`func (o *ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner) GetTrafficMirrorNumDestinationPort() string`

GetTrafficMirrorNumDestinationPort returns the TrafficMirrorNumDestinationPort field if non-nil, zero value otherwise.

### GetTrafficMirrorNumDestinationPortOk

`func (o *ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner) GetTrafficMirrorNumDestinationPortOk() (*string, bool)`

GetTrafficMirrorNumDestinationPortOk returns a tuple with the TrafficMirrorNumDestinationPort field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTrafficMirrorNumDestinationPort

`func (o *ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner) SetTrafficMirrorNumDestinationPort(v string)`

SetTrafficMirrorNumDestinationPort sets TrafficMirrorNumDestinationPort field to given value.

### HasTrafficMirrorNumDestinationPort

`func (o *ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner) HasTrafficMirrorNumDestinationPort() bool`

HasTrafficMirrorNumDestinationPort returns a boolean if a field has been set.

### GetTrafficMirrorNumInboundTraffic

`func (o *ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner) GetTrafficMirrorNumInboundTraffic() bool`

GetTrafficMirrorNumInboundTraffic returns the TrafficMirrorNumInboundTraffic field if non-nil, zero value otherwise.

### GetTrafficMirrorNumInboundTrafficOk

`func (o *ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner) GetTrafficMirrorNumInboundTrafficOk() (*bool, bool)`

GetTrafficMirrorNumInboundTrafficOk returns a tuple with the TrafficMirrorNumInboundTraffic field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTrafficMirrorNumInboundTraffic

`func (o *ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner) SetTrafficMirrorNumInboundTraffic(v bool)`

SetTrafficMirrorNumInboundTraffic sets TrafficMirrorNumInboundTraffic field to given value.

### HasTrafficMirrorNumInboundTraffic

`func (o *ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner) HasTrafficMirrorNumInboundTraffic() bool`

HasTrafficMirrorNumInboundTraffic returns a boolean if a field has been set.

### GetTrafficMirrorNumOutboundTraffic

`func (o *ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner) GetTrafficMirrorNumOutboundTraffic() bool`

GetTrafficMirrorNumOutboundTraffic returns the TrafficMirrorNumOutboundTraffic field if non-nil, zero value otherwise.

### GetTrafficMirrorNumOutboundTrafficOk

`func (o *ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner) GetTrafficMirrorNumOutboundTrafficOk() (*bool, bool)`

GetTrafficMirrorNumOutboundTrafficOk returns a tuple with the TrafficMirrorNumOutboundTraffic field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTrafficMirrorNumOutboundTraffic

`func (o *ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner) SetTrafficMirrorNumOutboundTraffic(v bool)`

SetTrafficMirrorNumOutboundTraffic sets TrafficMirrorNumOutboundTraffic field to given value.

### HasTrafficMirrorNumOutboundTraffic

`func (o *ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner) HasTrafficMirrorNumOutboundTraffic() bool`

HasTrafficMirrorNumOutboundTraffic returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


