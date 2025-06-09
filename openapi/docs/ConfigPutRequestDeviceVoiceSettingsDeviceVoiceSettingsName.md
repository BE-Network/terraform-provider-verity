# ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**DtmfMethod** | Pointer to **string** | Specifies how DTMF signals are carried | [optional] [default to "Inband"]
**Region** | Pointer to **string** | Region | [optional] [default to "US"]
**Protocol** | Pointer to **string** | Voice Protocol: MGCP or SIP | [optional] [default to "SIP"]
**ProxyServer** | Pointer to **string** | IP address or URI of the SIP proxy server for SIP signalling messages | [optional] [default to ""]
**ProxyServerPort** | Pointer to **NullableInt32** | Proxy Server Port | [optional] [default to 0]
**ProxyServerSecondary** | Pointer to **string** | IP address or URI of the secondary SIP proxy server for SIP signalling messages | [optional] [default to ""]
**ProxyServerSecondaryPort** | Pointer to **NullableInt32** | Secondary Proxy Server Port | [optional] [default to 0]
**RegistrarServer** | Pointer to **string** | Name or IP address or resolved name of the registrar server for SIP signalling messages. Examples: 10.10.10.10 and proxy.voip.net | [optional] [default to ""]
**RegistrarServerPort** | Pointer to **NullableInt32** | Registrar Server Port | [optional] [default to 0]
**RegistrarServerSecondary** | Pointer to **string** | Name or IP address or resolved name of the secondary registrar server for SIP signalling messages. Examples: 10.10.10.10 and proxy.voip.net | [optional] [default to ""]
**RegistrarServerSecondaryPort** | Pointer to **NullableInt32** | Secondary Registrar Server Port | [optional] [default to 0]
**UserAgentDomain** | Pointer to **string** | User Agent Domain | [optional] [default to ""]
**UserAgentTransport** | Pointer to **string** | User Agent Transport | [optional] [default to "UDP"]
**UserAgentPort** | Pointer to **NullableInt32** | User Agent Port | [optional] [default to 0]
**OutboundProxy** | Pointer to **string** | IP address or URI of the outbound proxy server for SIP signalling messages. An outbound SIP proxy may or may not be required within a given network | [optional] [default to ""]
**OutboundProxyPort** | Pointer to **NullableInt32** | Outbound Proxy Port | [optional] [default to 0]
**OutboundProxySecondary** | Pointer to **string** | IP address or URI of the secondary outbound proxy server for SIP signalling messages. An outbound SIP proxy may or may not be required within a given network | [optional] [default to ""]
**OutboundProxySecondaryPort** | Pointer to **NullableInt32** | Secondary Outbound Proxy Port | [optional] [default to 0]
**RegistrationPeriod** | Pointer to **NullableInt32** | Specifies the time in seconds to start the re-registration process. The default value is 3240 seconds | [optional] [default to 3240]
**RegisterExpires** | Pointer to **NullableInt32** | SIP registration expiration time in seconds. If value is 0, the SIP agent does not add an expiration time to the registration requests and does not perform re-registration. The default value is 3600 seconds | [optional] [default to 3600]
**VoicemailServer** | Pointer to **string** | Name or IP address or resolved name of the external voicemail server if not provided by SIP server for MWI control. Examples: 10.10.10.10 and proxy.voip.net | [optional] [default to ""]
**VoicemailServerPort** | Pointer to **NullableInt32** | Voicemail Server Port | [optional] [default to 0]
**VoicemailServerExpires** | Pointer to **NullableInt32** | Voicemail server expiration time in seconds. If value is 0, the Register Expires time is used instead. The default value is 3600 seconds | [optional] [default to 3600]
**SipDscpMark** | Pointer to **NullableInt32** | Sip Differentiated Services Code point (DSCP) | [optional] [default to 0]
**CallAgent1** | Pointer to **string** | Call Agent 1 | [optional] [default to ""]
**CallAgentPort1** | Pointer to **NullableInt32** | Call Agent Port 1 | [optional] [default to 0]
**CallAgent2** | Pointer to **string** | Call Agent 2 | [optional] [default to ""]
**CallAgentPort2** | Pointer to **NullableInt32** | Call Agent Port 2 | [optional] [default to 0]
**Domain** | Pointer to **string** | Domain | [optional] [default to ""]
**MgcpDscpMark** | Pointer to **NullableInt32** | MGCP Differentiated Services Code point (DSCP) | [optional] [default to 0]
**TerminationBase** | Pointer to **string** | Base string for the MGCP physical termination id(s) | [optional] [default to "aaln/"]
**LocalPortMin** | Pointer to **NullableInt32** | Defines the base RTP port that should be used for voice traffic | [optional] [default to 30000]
**LocalPortMax** | Pointer to **NullableInt32** | Defines the highest RTP port used for voice traffic, must be greater than local Local Port Min | [optional] [default to 30200]
**EventPayloadType** | Pointer to **NullableInt32** | Telephone Event Payload Type | [optional] [default to 101]
**CasEvents** | Pointer to **int32** | Enables or disables handling of CAS via RTP CAS events. Valid values are 0 &#x3D; off and 1 &#x3D; on | [optional] [default to 0]
**DscpMark** | Pointer to **NullableInt32** | Differentiated Services Code Point (DSCP) to be used for outgoing RTP packets | [optional] [default to 0]
**Rtcp** | Pointer to **bool** | RTCP Enable | [optional] [default to true]
**FaxT38** | Pointer to **bool** | Fax T.38 Enable | [optional] [default to false]
**BitRate** | Pointer to **string** | T.38 Bit Rate in bps. Most available fax machines support up to 14,400bps | [optional] [default to "14400"]
**CancelCallWaiting** | Pointer to **string** | Cancel Call waiting | [optional] [default to "*70"]
**CallHold** | Pointer to **string** | Call hold | [optional] [default to "*9"]
**CidsActivate** | Pointer to **string** | Caller ID Delivery Blocking (single call)  Activate | [optional] [default to "*67"]
**CidsDeactivate** | Pointer to **string** | Caller ID Delivery Blocking (single call) Deactivate | [optional] [default to "*82"]
**DoNotDisturbActivate** | Pointer to **string** | Do not Disturb Activate | [optional] [default to "*78"]
**DoNotDisturbDeactivate** | Pointer to **string** | Do not Disturb Deactivate | [optional] [default to "*79"]
**DoNotDisturbPinChange** | Pointer to **string** | Do not Disturb PIN Change | [optional] [default to "*10"]
**EmergencyServiceNumber** | Pointer to **string** | Emergency Service Number | [optional] [default to "911"]
**AnonCidBlockActivate** | Pointer to **string** | Anonymoes Caller ID Block Activate | [optional] [default to "*77"]
**AnonCidBlockDeactivate** | Pointer to **string** | Anonymous Caller ID Block Deactivate | [optional] [default to "*87"]
**CallForwardUnconditionalActivate** | Pointer to **string** | Call Forward Unconditional Activate | [optional] [default to "*72"]
**CallForwardUnconditionalDeactivate** | Pointer to **string** | Call Forward Unconditional Deactivate | [optional] [default to "*73"]
**CallForwardOnBusyActivate** | Pointer to **string** | Call Forward On Busy Activate | [optional] [default to "*90"]
**CallForwardOnBusyDeactivate** | Pointer to **string** | Call Forward On Busy Deactivate | [optional] [default to "*91"]
**CallForwardOnNoAnswerActivate** | Pointer to **string** | Call Forward On No Answer Activate | [optional] [default to "*92"]
**CallForwardOnNoAnswerDeactivate** | Pointer to **string** | Call Forward On No Answer Deactivate | [optional] [default to "*93"]
**Intercom1** | Pointer to **string** | Intercom 1 | [optional] [default to "*53"]
**Intercom2** | Pointer to **string** | Intercom 2 | [optional] [default to "*54"]
**Intercom3** | Pointer to **string** | Intercom 3 | [optional] [default to "*55"]
**Codecs** | Pointer to [**[]ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInner**](ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**ConfigPutRequestPacketQueuePacketQueueNameObjectProperties**](ConfigPutRequestPacketQueuePacketQueueNameObjectProperties.md) |  | [optional] 

## Methods

### NewConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName

`func NewConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName() *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName`

NewConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName instantiates a new ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameWithDefaults

`func NewConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameWithDefaults() *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName`

NewConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameWithDefaults instantiates a new ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetDtmfMethod

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetDtmfMethod() string`

GetDtmfMethod returns the DtmfMethod field if non-nil, zero value otherwise.

### GetDtmfMethodOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetDtmfMethodOk() (*string, bool)`

GetDtmfMethodOk returns a tuple with the DtmfMethod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDtmfMethod

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetDtmfMethod(v string)`

SetDtmfMethod sets DtmfMethod field to given value.

### HasDtmfMethod

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasDtmfMethod() bool`

HasDtmfMethod returns a boolean if a field has been set.

### GetRegion

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetRegion() string`

GetRegion returns the Region field if non-nil, zero value otherwise.

### GetRegionOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetRegionOk() (*string, bool)`

GetRegionOk returns a tuple with the Region field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegion

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetRegion(v string)`

SetRegion sets Region field to given value.

### HasRegion

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasRegion() bool`

HasRegion returns a boolean if a field has been set.

### GetProtocol

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetProtocol() string`

GetProtocol returns the Protocol field if non-nil, zero value otherwise.

### GetProtocolOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetProtocolOk() (*string, bool)`

GetProtocolOk returns a tuple with the Protocol field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProtocol

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetProtocol(v string)`

SetProtocol sets Protocol field to given value.

### HasProtocol

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasProtocol() bool`

HasProtocol returns a boolean if a field has been set.

### GetProxyServer

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetProxyServer() string`

GetProxyServer returns the ProxyServer field if non-nil, zero value otherwise.

### GetProxyServerOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetProxyServerOk() (*string, bool)`

GetProxyServerOk returns a tuple with the ProxyServer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProxyServer

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetProxyServer(v string)`

SetProxyServer sets ProxyServer field to given value.

### HasProxyServer

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasProxyServer() bool`

HasProxyServer returns a boolean if a field has been set.

### GetProxyServerPort

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetProxyServerPort() int32`

GetProxyServerPort returns the ProxyServerPort field if non-nil, zero value otherwise.

### GetProxyServerPortOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetProxyServerPortOk() (*int32, bool)`

GetProxyServerPortOk returns a tuple with the ProxyServerPort field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProxyServerPort

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetProxyServerPort(v int32)`

SetProxyServerPort sets ProxyServerPort field to given value.

### HasProxyServerPort

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasProxyServerPort() bool`

HasProxyServerPort returns a boolean if a field has been set.

### SetProxyServerPortNil

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetProxyServerPortNil(b bool)`

 SetProxyServerPortNil sets the value for ProxyServerPort to be an explicit nil

### UnsetProxyServerPort
`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) UnsetProxyServerPort()`

UnsetProxyServerPort ensures that no value is present for ProxyServerPort, not even an explicit nil
### GetProxyServerSecondary

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetProxyServerSecondary() string`

GetProxyServerSecondary returns the ProxyServerSecondary field if non-nil, zero value otherwise.

### GetProxyServerSecondaryOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetProxyServerSecondaryOk() (*string, bool)`

GetProxyServerSecondaryOk returns a tuple with the ProxyServerSecondary field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProxyServerSecondary

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetProxyServerSecondary(v string)`

SetProxyServerSecondary sets ProxyServerSecondary field to given value.

### HasProxyServerSecondary

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasProxyServerSecondary() bool`

HasProxyServerSecondary returns a boolean if a field has been set.

### GetProxyServerSecondaryPort

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetProxyServerSecondaryPort() int32`

GetProxyServerSecondaryPort returns the ProxyServerSecondaryPort field if non-nil, zero value otherwise.

### GetProxyServerSecondaryPortOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetProxyServerSecondaryPortOk() (*int32, bool)`

GetProxyServerSecondaryPortOk returns a tuple with the ProxyServerSecondaryPort field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProxyServerSecondaryPort

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetProxyServerSecondaryPort(v int32)`

SetProxyServerSecondaryPort sets ProxyServerSecondaryPort field to given value.

### HasProxyServerSecondaryPort

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasProxyServerSecondaryPort() bool`

HasProxyServerSecondaryPort returns a boolean if a field has been set.

### SetProxyServerSecondaryPortNil

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetProxyServerSecondaryPortNil(b bool)`

 SetProxyServerSecondaryPortNil sets the value for ProxyServerSecondaryPort to be an explicit nil

### UnsetProxyServerSecondaryPort
`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) UnsetProxyServerSecondaryPort()`

UnsetProxyServerSecondaryPort ensures that no value is present for ProxyServerSecondaryPort, not even an explicit nil
### GetRegistrarServer

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetRegistrarServer() string`

GetRegistrarServer returns the RegistrarServer field if non-nil, zero value otherwise.

### GetRegistrarServerOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetRegistrarServerOk() (*string, bool)`

GetRegistrarServerOk returns a tuple with the RegistrarServer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegistrarServer

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetRegistrarServer(v string)`

SetRegistrarServer sets RegistrarServer field to given value.

### HasRegistrarServer

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasRegistrarServer() bool`

HasRegistrarServer returns a boolean if a field has been set.

### GetRegistrarServerPort

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetRegistrarServerPort() int32`

GetRegistrarServerPort returns the RegistrarServerPort field if non-nil, zero value otherwise.

### GetRegistrarServerPortOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetRegistrarServerPortOk() (*int32, bool)`

GetRegistrarServerPortOk returns a tuple with the RegistrarServerPort field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegistrarServerPort

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetRegistrarServerPort(v int32)`

SetRegistrarServerPort sets RegistrarServerPort field to given value.

### HasRegistrarServerPort

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasRegistrarServerPort() bool`

HasRegistrarServerPort returns a boolean if a field has been set.

### SetRegistrarServerPortNil

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetRegistrarServerPortNil(b bool)`

 SetRegistrarServerPortNil sets the value for RegistrarServerPort to be an explicit nil

### UnsetRegistrarServerPort
`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) UnsetRegistrarServerPort()`

UnsetRegistrarServerPort ensures that no value is present for RegistrarServerPort, not even an explicit nil
### GetRegistrarServerSecondary

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetRegistrarServerSecondary() string`

GetRegistrarServerSecondary returns the RegistrarServerSecondary field if non-nil, zero value otherwise.

### GetRegistrarServerSecondaryOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetRegistrarServerSecondaryOk() (*string, bool)`

GetRegistrarServerSecondaryOk returns a tuple with the RegistrarServerSecondary field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegistrarServerSecondary

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetRegistrarServerSecondary(v string)`

SetRegistrarServerSecondary sets RegistrarServerSecondary field to given value.

### HasRegistrarServerSecondary

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasRegistrarServerSecondary() bool`

HasRegistrarServerSecondary returns a boolean if a field has been set.

### GetRegistrarServerSecondaryPort

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetRegistrarServerSecondaryPort() int32`

GetRegistrarServerSecondaryPort returns the RegistrarServerSecondaryPort field if non-nil, zero value otherwise.

### GetRegistrarServerSecondaryPortOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetRegistrarServerSecondaryPortOk() (*int32, bool)`

GetRegistrarServerSecondaryPortOk returns a tuple with the RegistrarServerSecondaryPort field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegistrarServerSecondaryPort

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetRegistrarServerSecondaryPort(v int32)`

SetRegistrarServerSecondaryPort sets RegistrarServerSecondaryPort field to given value.

### HasRegistrarServerSecondaryPort

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasRegistrarServerSecondaryPort() bool`

HasRegistrarServerSecondaryPort returns a boolean if a field has been set.

### SetRegistrarServerSecondaryPortNil

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetRegistrarServerSecondaryPortNil(b bool)`

 SetRegistrarServerSecondaryPortNil sets the value for RegistrarServerSecondaryPort to be an explicit nil

### UnsetRegistrarServerSecondaryPort
`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) UnsetRegistrarServerSecondaryPort()`

UnsetRegistrarServerSecondaryPort ensures that no value is present for RegistrarServerSecondaryPort, not even an explicit nil
### GetUserAgentDomain

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetUserAgentDomain() string`

GetUserAgentDomain returns the UserAgentDomain field if non-nil, zero value otherwise.

### GetUserAgentDomainOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetUserAgentDomainOk() (*string, bool)`

GetUserAgentDomainOk returns a tuple with the UserAgentDomain field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserAgentDomain

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetUserAgentDomain(v string)`

SetUserAgentDomain sets UserAgentDomain field to given value.

### HasUserAgentDomain

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasUserAgentDomain() bool`

HasUserAgentDomain returns a boolean if a field has been set.

### GetUserAgentTransport

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetUserAgentTransport() string`

GetUserAgentTransport returns the UserAgentTransport field if non-nil, zero value otherwise.

### GetUserAgentTransportOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetUserAgentTransportOk() (*string, bool)`

GetUserAgentTransportOk returns a tuple with the UserAgentTransport field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserAgentTransport

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetUserAgentTransport(v string)`

SetUserAgentTransport sets UserAgentTransport field to given value.

### HasUserAgentTransport

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasUserAgentTransport() bool`

HasUserAgentTransport returns a boolean if a field has been set.

### GetUserAgentPort

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetUserAgentPort() int32`

GetUserAgentPort returns the UserAgentPort field if non-nil, zero value otherwise.

### GetUserAgentPortOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetUserAgentPortOk() (*int32, bool)`

GetUserAgentPortOk returns a tuple with the UserAgentPort field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserAgentPort

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetUserAgentPort(v int32)`

SetUserAgentPort sets UserAgentPort field to given value.

### HasUserAgentPort

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasUserAgentPort() bool`

HasUserAgentPort returns a boolean if a field has been set.

### SetUserAgentPortNil

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetUserAgentPortNil(b bool)`

 SetUserAgentPortNil sets the value for UserAgentPort to be an explicit nil

### UnsetUserAgentPort
`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) UnsetUserAgentPort()`

UnsetUserAgentPort ensures that no value is present for UserAgentPort, not even an explicit nil
### GetOutboundProxy

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetOutboundProxy() string`

GetOutboundProxy returns the OutboundProxy field if non-nil, zero value otherwise.

### GetOutboundProxyOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetOutboundProxyOk() (*string, bool)`

GetOutboundProxyOk returns a tuple with the OutboundProxy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutboundProxy

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetOutboundProxy(v string)`

SetOutboundProxy sets OutboundProxy field to given value.

### HasOutboundProxy

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasOutboundProxy() bool`

HasOutboundProxy returns a boolean if a field has been set.

### GetOutboundProxyPort

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetOutboundProxyPort() int32`

GetOutboundProxyPort returns the OutboundProxyPort field if non-nil, zero value otherwise.

### GetOutboundProxyPortOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetOutboundProxyPortOk() (*int32, bool)`

GetOutboundProxyPortOk returns a tuple with the OutboundProxyPort field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutboundProxyPort

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetOutboundProxyPort(v int32)`

SetOutboundProxyPort sets OutboundProxyPort field to given value.

### HasOutboundProxyPort

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasOutboundProxyPort() bool`

HasOutboundProxyPort returns a boolean if a field has been set.

### SetOutboundProxyPortNil

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetOutboundProxyPortNil(b bool)`

 SetOutboundProxyPortNil sets the value for OutboundProxyPort to be an explicit nil

### UnsetOutboundProxyPort
`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) UnsetOutboundProxyPort()`

UnsetOutboundProxyPort ensures that no value is present for OutboundProxyPort, not even an explicit nil
### GetOutboundProxySecondary

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetOutboundProxySecondary() string`

GetOutboundProxySecondary returns the OutboundProxySecondary field if non-nil, zero value otherwise.

### GetOutboundProxySecondaryOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetOutboundProxySecondaryOk() (*string, bool)`

GetOutboundProxySecondaryOk returns a tuple with the OutboundProxySecondary field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutboundProxySecondary

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetOutboundProxySecondary(v string)`

SetOutboundProxySecondary sets OutboundProxySecondary field to given value.

### HasOutboundProxySecondary

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasOutboundProxySecondary() bool`

HasOutboundProxySecondary returns a boolean if a field has been set.

### GetOutboundProxySecondaryPort

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetOutboundProxySecondaryPort() int32`

GetOutboundProxySecondaryPort returns the OutboundProxySecondaryPort field if non-nil, zero value otherwise.

### GetOutboundProxySecondaryPortOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetOutboundProxySecondaryPortOk() (*int32, bool)`

GetOutboundProxySecondaryPortOk returns a tuple with the OutboundProxySecondaryPort field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutboundProxySecondaryPort

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetOutboundProxySecondaryPort(v int32)`

SetOutboundProxySecondaryPort sets OutboundProxySecondaryPort field to given value.

### HasOutboundProxySecondaryPort

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasOutboundProxySecondaryPort() bool`

HasOutboundProxySecondaryPort returns a boolean if a field has been set.

### SetOutboundProxySecondaryPortNil

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetOutboundProxySecondaryPortNil(b bool)`

 SetOutboundProxySecondaryPortNil sets the value for OutboundProxySecondaryPort to be an explicit nil

### UnsetOutboundProxySecondaryPort
`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) UnsetOutboundProxySecondaryPort()`

UnsetOutboundProxySecondaryPort ensures that no value is present for OutboundProxySecondaryPort, not even an explicit nil
### GetRegistrationPeriod

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetRegistrationPeriod() int32`

GetRegistrationPeriod returns the RegistrationPeriod field if non-nil, zero value otherwise.

### GetRegistrationPeriodOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetRegistrationPeriodOk() (*int32, bool)`

GetRegistrationPeriodOk returns a tuple with the RegistrationPeriod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegistrationPeriod

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetRegistrationPeriod(v int32)`

SetRegistrationPeriod sets RegistrationPeriod field to given value.

### HasRegistrationPeriod

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasRegistrationPeriod() bool`

HasRegistrationPeriod returns a boolean if a field has been set.

### SetRegistrationPeriodNil

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetRegistrationPeriodNil(b bool)`

 SetRegistrationPeriodNil sets the value for RegistrationPeriod to be an explicit nil

### UnsetRegistrationPeriod
`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) UnsetRegistrationPeriod()`

UnsetRegistrationPeriod ensures that no value is present for RegistrationPeriod, not even an explicit nil
### GetRegisterExpires

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetRegisterExpires() int32`

GetRegisterExpires returns the RegisterExpires field if non-nil, zero value otherwise.

### GetRegisterExpiresOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetRegisterExpiresOk() (*int32, bool)`

GetRegisterExpiresOk returns a tuple with the RegisterExpires field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegisterExpires

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetRegisterExpires(v int32)`

SetRegisterExpires sets RegisterExpires field to given value.

### HasRegisterExpires

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasRegisterExpires() bool`

HasRegisterExpires returns a boolean if a field has been set.

### SetRegisterExpiresNil

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetRegisterExpiresNil(b bool)`

 SetRegisterExpiresNil sets the value for RegisterExpires to be an explicit nil

### UnsetRegisterExpires
`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) UnsetRegisterExpires()`

UnsetRegisterExpires ensures that no value is present for RegisterExpires, not even an explicit nil
### GetVoicemailServer

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetVoicemailServer() string`

GetVoicemailServer returns the VoicemailServer field if non-nil, zero value otherwise.

### GetVoicemailServerOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetVoicemailServerOk() (*string, bool)`

GetVoicemailServerOk returns a tuple with the VoicemailServer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVoicemailServer

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetVoicemailServer(v string)`

SetVoicemailServer sets VoicemailServer field to given value.

### HasVoicemailServer

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasVoicemailServer() bool`

HasVoicemailServer returns a boolean if a field has been set.

### GetVoicemailServerPort

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetVoicemailServerPort() int32`

GetVoicemailServerPort returns the VoicemailServerPort field if non-nil, zero value otherwise.

### GetVoicemailServerPortOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetVoicemailServerPortOk() (*int32, bool)`

GetVoicemailServerPortOk returns a tuple with the VoicemailServerPort field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVoicemailServerPort

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetVoicemailServerPort(v int32)`

SetVoicemailServerPort sets VoicemailServerPort field to given value.

### HasVoicemailServerPort

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasVoicemailServerPort() bool`

HasVoicemailServerPort returns a boolean if a field has been set.

### SetVoicemailServerPortNil

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetVoicemailServerPortNil(b bool)`

 SetVoicemailServerPortNil sets the value for VoicemailServerPort to be an explicit nil

### UnsetVoicemailServerPort
`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) UnsetVoicemailServerPort()`

UnsetVoicemailServerPort ensures that no value is present for VoicemailServerPort, not even an explicit nil
### GetVoicemailServerExpires

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetVoicemailServerExpires() int32`

GetVoicemailServerExpires returns the VoicemailServerExpires field if non-nil, zero value otherwise.

### GetVoicemailServerExpiresOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetVoicemailServerExpiresOk() (*int32, bool)`

GetVoicemailServerExpiresOk returns a tuple with the VoicemailServerExpires field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVoicemailServerExpires

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetVoicemailServerExpires(v int32)`

SetVoicemailServerExpires sets VoicemailServerExpires field to given value.

### HasVoicemailServerExpires

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasVoicemailServerExpires() bool`

HasVoicemailServerExpires returns a boolean if a field has been set.

### SetVoicemailServerExpiresNil

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetVoicemailServerExpiresNil(b bool)`

 SetVoicemailServerExpiresNil sets the value for VoicemailServerExpires to be an explicit nil

### UnsetVoicemailServerExpires
`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) UnsetVoicemailServerExpires()`

UnsetVoicemailServerExpires ensures that no value is present for VoicemailServerExpires, not even an explicit nil
### GetSipDscpMark

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetSipDscpMark() int32`

GetSipDscpMark returns the SipDscpMark field if non-nil, zero value otherwise.

### GetSipDscpMarkOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetSipDscpMarkOk() (*int32, bool)`

GetSipDscpMarkOk returns a tuple with the SipDscpMark field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSipDscpMark

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetSipDscpMark(v int32)`

SetSipDscpMark sets SipDscpMark field to given value.

### HasSipDscpMark

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasSipDscpMark() bool`

HasSipDscpMark returns a boolean if a field has been set.

### SetSipDscpMarkNil

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetSipDscpMarkNil(b bool)`

 SetSipDscpMarkNil sets the value for SipDscpMark to be an explicit nil

### UnsetSipDscpMark
`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) UnsetSipDscpMark()`

UnsetSipDscpMark ensures that no value is present for SipDscpMark, not even an explicit nil
### GetCallAgent1

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetCallAgent1() string`

GetCallAgent1 returns the CallAgent1 field if non-nil, zero value otherwise.

### GetCallAgent1Ok

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetCallAgent1Ok() (*string, bool)`

GetCallAgent1Ok returns a tuple with the CallAgent1 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallAgent1

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetCallAgent1(v string)`

SetCallAgent1 sets CallAgent1 field to given value.

### HasCallAgent1

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasCallAgent1() bool`

HasCallAgent1 returns a boolean if a field has been set.

### GetCallAgentPort1

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetCallAgentPort1() int32`

GetCallAgentPort1 returns the CallAgentPort1 field if non-nil, zero value otherwise.

### GetCallAgentPort1Ok

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetCallAgentPort1Ok() (*int32, bool)`

GetCallAgentPort1Ok returns a tuple with the CallAgentPort1 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallAgentPort1

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetCallAgentPort1(v int32)`

SetCallAgentPort1 sets CallAgentPort1 field to given value.

### HasCallAgentPort1

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasCallAgentPort1() bool`

HasCallAgentPort1 returns a boolean if a field has been set.

### SetCallAgentPort1Nil

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetCallAgentPort1Nil(b bool)`

 SetCallAgentPort1Nil sets the value for CallAgentPort1 to be an explicit nil

### UnsetCallAgentPort1
`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) UnsetCallAgentPort1()`

UnsetCallAgentPort1 ensures that no value is present for CallAgentPort1, not even an explicit nil
### GetCallAgent2

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetCallAgent2() string`

GetCallAgent2 returns the CallAgent2 field if non-nil, zero value otherwise.

### GetCallAgent2Ok

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetCallAgent2Ok() (*string, bool)`

GetCallAgent2Ok returns a tuple with the CallAgent2 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallAgent2

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetCallAgent2(v string)`

SetCallAgent2 sets CallAgent2 field to given value.

### HasCallAgent2

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasCallAgent2() bool`

HasCallAgent2 returns a boolean if a field has been set.

### GetCallAgentPort2

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetCallAgentPort2() int32`

GetCallAgentPort2 returns the CallAgentPort2 field if non-nil, zero value otherwise.

### GetCallAgentPort2Ok

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetCallAgentPort2Ok() (*int32, bool)`

GetCallAgentPort2Ok returns a tuple with the CallAgentPort2 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallAgentPort2

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetCallAgentPort2(v int32)`

SetCallAgentPort2 sets CallAgentPort2 field to given value.

### HasCallAgentPort2

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasCallAgentPort2() bool`

HasCallAgentPort2 returns a boolean if a field has been set.

### SetCallAgentPort2Nil

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetCallAgentPort2Nil(b bool)`

 SetCallAgentPort2Nil sets the value for CallAgentPort2 to be an explicit nil

### UnsetCallAgentPort2
`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) UnsetCallAgentPort2()`

UnsetCallAgentPort2 ensures that no value is present for CallAgentPort2, not even an explicit nil
### GetDomain

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetDomain() string`

GetDomain returns the Domain field if non-nil, zero value otherwise.

### GetDomainOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetDomainOk() (*string, bool)`

GetDomainOk returns a tuple with the Domain field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDomain

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetDomain(v string)`

SetDomain sets Domain field to given value.

### HasDomain

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasDomain() bool`

HasDomain returns a boolean if a field has been set.

### GetMgcpDscpMark

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetMgcpDscpMark() int32`

GetMgcpDscpMark returns the MgcpDscpMark field if non-nil, zero value otherwise.

### GetMgcpDscpMarkOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetMgcpDscpMarkOk() (*int32, bool)`

GetMgcpDscpMarkOk returns a tuple with the MgcpDscpMark field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMgcpDscpMark

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetMgcpDscpMark(v int32)`

SetMgcpDscpMark sets MgcpDscpMark field to given value.

### HasMgcpDscpMark

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasMgcpDscpMark() bool`

HasMgcpDscpMark returns a boolean if a field has been set.

### SetMgcpDscpMarkNil

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetMgcpDscpMarkNil(b bool)`

 SetMgcpDscpMarkNil sets the value for MgcpDscpMark to be an explicit nil

### UnsetMgcpDscpMark
`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) UnsetMgcpDscpMark()`

UnsetMgcpDscpMark ensures that no value is present for MgcpDscpMark, not even an explicit nil
### GetTerminationBase

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetTerminationBase() string`

GetTerminationBase returns the TerminationBase field if non-nil, zero value otherwise.

### GetTerminationBaseOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetTerminationBaseOk() (*string, bool)`

GetTerminationBaseOk returns a tuple with the TerminationBase field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTerminationBase

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetTerminationBase(v string)`

SetTerminationBase sets TerminationBase field to given value.

### HasTerminationBase

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasTerminationBase() bool`

HasTerminationBase returns a boolean if a field has been set.

### GetLocalPortMin

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetLocalPortMin() int32`

GetLocalPortMin returns the LocalPortMin field if non-nil, zero value otherwise.

### GetLocalPortMinOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetLocalPortMinOk() (*int32, bool)`

GetLocalPortMinOk returns a tuple with the LocalPortMin field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocalPortMin

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetLocalPortMin(v int32)`

SetLocalPortMin sets LocalPortMin field to given value.

### HasLocalPortMin

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasLocalPortMin() bool`

HasLocalPortMin returns a boolean if a field has been set.

### SetLocalPortMinNil

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetLocalPortMinNil(b bool)`

 SetLocalPortMinNil sets the value for LocalPortMin to be an explicit nil

### UnsetLocalPortMin
`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) UnsetLocalPortMin()`

UnsetLocalPortMin ensures that no value is present for LocalPortMin, not even an explicit nil
### GetLocalPortMax

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetLocalPortMax() int32`

GetLocalPortMax returns the LocalPortMax field if non-nil, zero value otherwise.

### GetLocalPortMaxOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetLocalPortMaxOk() (*int32, bool)`

GetLocalPortMaxOk returns a tuple with the LocalPortMax field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocalPortMax

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetLocalPortMax(v int32)`

SetLocalPortMax sets LocalPortMax field to given value.

### HasLocalPortMax

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasLocalPortMax() bool`

HasLocalPortMax returns a boolean if a field has been set.

### SetLocalPortMaxNil

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetLocalPortMaxNil(b bool)`

 SetLocalPortMaxNil sets the value for LocalPortMax to be an explicit nil

### UnsetLocalPortMax
`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) UnsetLocalPortMax()`

UnsetLocalPortMax ensures that no value is present for LocalPortMax, not even an explicit nil
### GetEventPayloadType

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetEventPayloadType() int32`

GetEventPayloadType returns the EventPayloadType field if non-nil, zero value otherwise.

### GetEventPayloadTypeOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetEventPayloadTypeOk() (*int32, bool)`

GetEventPayloadTypeOk returns a tuple with the EventPayloadType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEventPayloadType

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetEventPayloadType(v int32)`

SetEventPayloadType sets EventPayloadType field to given value.

### HasEventPayloadType

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasEventPayloadType() bool`

HasEventPayloadType returns a boolean if a field has been set.

### SetEventPayloadTypeNil

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetEventPayloadTypeNil(b bool)`

 SetEventPayloadTypeNil sets the value for EventPayloadType to be an explicit nil

### UnsetEventPayloadType
`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) UnsetEventPayloadType()`

UnsetEventPayloadType ensures that no value is present for EventPayloadType, not even an explicit nil
### GetCasEvents

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetCasEvents() int32`

GetCasEvents returns the CasEvents field if non-nil, zero value otherwise.

### GetCasEventsOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetCasEventsOk() (*int32, bool)`

GetCasEventsOk returns a tuple with the CasEvents field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCasEvents

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetCasEvents(v int32)`

SetCasEvents sets CasEvents field to given value.

### HasCasEvents

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasCasEvents() bool`

HasCasEvents returns a boolean if a field has been set.

### GetDscpMark

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetDscpMark() int32`

GetDscpMark returns the DscpMark field if non-nil, zero value otherwise.

### GetDscpMarkOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetDscpMarkOk() (*int32, bool)`

GetDscpMarkOk returns a tuple with the DscpMark field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDscpMark

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetDscpMark(v int32)`

SetDscpMark sets DscpMark field to given value.

### HasDscpMark

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasDscpMark() bool`

HasDscpMark returns a boolean if a field has been set.

### SetDscpMarkNil

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetDscpMarkNil(b bool)`

 SetDscpMarkNil sets the value for DscpMark to be an explicit nil

### UnsetDscpMark
`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) UnsetDscpMark()`

UnsetDscpMark ensures that no value is present for DscpMark, not even an explicit nil
### GetRtcp

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetRtcp() bool`

GetRtcp returns the Rtcp field if non-nil, zero value otherwise.

### GetRtcpOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetRtcpOk() (*bool, bool)`

GetRtcpOk returns a tuple with the Rtcp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRtcp

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetRtcp(v bool)`

SetRtcp sets Rtcp field to given value.

### HasRtcp

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasRtcp() bool`

HasRtcp returns a boolean if a field has been set.

### GetFaxT38

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetFaxT38() bool`

GetFaxT38 returns the FaxT38 field if non-nil, zero value otherwise.

### GetFaxT38Ok

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetFaxT38Ok() (*bool, bool)`

GetFaxT38Ok returns a tuple with the FaxT38 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFaxT38

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetFaxT38(v bool)`

SetFaxT38 sets FaxT38 field to given value.

### HasFaxT38

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasFaxT38() bool`

HasFaxT38 returns a boolean if a field has been set.

### GetBitRate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetBitRate() string`

GetBitRate returns the BitRate field if non-nil, zero value otherwise.

### GetBitRateOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetBitRateOk() (*string, bool)`

GetBitRateOk returns a tuple with the BitRate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBitRate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetBitRate(v string)`

SetBitRate sets BitRate field to given value.

### HasBitRate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasBitRate() bool`

HasBitRate returns a boolean if a field has been set.

### GetCancelCallWaiting

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetCancelCallWaiting() string`

GetCancelCallWaiting returns the CancelCallWaiting field if non-nil, zero value otherwise.

### GetCancelCallWaitingOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetCancelCallWaitingOk() (*string, bool)`

GetCancelCallWaitingOk returns a tuple with the CancelCallWaiting field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCancelCallWaiting

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetCancelCallWaiting(v string)`

SetCancelCallWaiting sets CancelCallWaiting field to given value.

### HasCancelCallWaiting

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasCancelCallWaiting() bool`

HasCancelCallWaiting returns a boolean if a field has been set.

### GetCallHold

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetCallHold() string`

GetCallHold returns the CallHold field if non-nil, zero value otherwise.

### GetCallHoldOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetCallHoldOk() (*string, bool)`

GetCallHoldOk returns a tuple with the CallHold field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallHold

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetCallHold(v string)`

SetCallHold sets CallHold field to given value.

### HasCallHold

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasCallHold() bool`

HasCallHold returns a boolean if a field has been set.

### GetCidsActivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetCidsActivate() string`

GetCidsActivate returns the CidsActivate field if non-nil, zero value otherwise.

### GetCidsActivateOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetCidsActivateOk() (*string, bool)`

GetCidsActivateOk returns a tuple with the CidsActivate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCidsActivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetCidsActivate(v string)`

SetCidsActivate sets CidsActivate field to given value.

### HasCidsActivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasCidsActivate() bool`

HasCidsActivate returns a boolean if a field has been set.

### GetCidsDeactivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetCidsDeactivate() string`

GetCidsDeactivate returns the CidsDeactivate field if non-nil, zero value otherwise.

### GetCidsDeactivateOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetCidsDeactivateOk() (*string, bool)`

GetCidsDeactivateOk returns a tuple with the CidsDeactivate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCidsDeactivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetCidsDeactivate(v string)`

SetCidsDeactivate sets CidsDeactivate field to given value.

### HasCidsDeactivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasCidsDeactivate() bool`

HasCidsDeactivate returns a boolean if a field has been set.

### GetDoNotDisturbActivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetDoNotDisturbActivate() string`

GetDoNotDisturbActivate returns the DoNotDisturbActivate field if non-nil, zero value otherwise.

### GetDoNotDisturbActivateOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetDoNotDisturbActivateOk() (*string, bool)`

GetDoNotDisturbActivateOk returns a tuple with the DoNotDisturbActivate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDoNotDisturbActivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetDoNotDisturbActivate(v string)`

SetDoNotDisturbActivate sets DoNotDisturbActivate field to given value.

### HasDoNotDisturbActivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasDoNotDisturbActivate() bool`

HasDoNotDisturbActivate returns a boolean if a field has been set.

### GetDoNotDisturbDeactivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetDoNotDisturbDeactivate() string`

GetDoNotDisturbDeactivate returns the DoNotDisturbDeactivate field if non-nil, zero value otherwise.

### GetDoNotDisturbDeactivateOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetDoNotDisturbDeactivateOk() (*string, bool)`

GetDoNotDisturbDeactivateOk returns a tuple with the DoNotDisturbDeactivate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDoNotDisturbDeactivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetDoNotDisturbDeactivate(v string)`

SetDoNotDisturbDeactivate sets DoNotDisturbDeactivate field to given value.

### HasDoNotDisturbDeactivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasDoNotDisturbDeactivate() bool`

HasDoNotDisturbDeactivate returns a boolean if a field has been set.

### GetDoNotDisturbPinChange

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetDoNotDisturbPinChange() string`

GetDoNotDisturbPinChange returns the DoNotDisturbPinChange field if non-nil, zero value otherwise.

### GetDoNotDisturbPinChangeOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetDoNotDisturbPinChangeOk() (*string, bool)`

GetDoNotDisturbPinChangeOk returns a tuple with the DoNotDisturbPinChange field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDoNotDisturbPinChange

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetDoNotDisturbPinChange(v string)`

SetDoNotDisturbPinChange sets DoNotDisturbPinChange field to given value.

### HasDoNotDisturbPinChange

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasDoNotDisturbPinChange() bool`

HasDoNotDisturbPinChange returns a boolean if a field has been set.

### GetEmergencyServiceNumber

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetEmergencyServiceNumber() string`

GetEmergencyServiceNumber returns the EmergencyServiceNumber field if non-nil, zero value otherwise.

### GetEmergencyServiceNumberOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetEmergencyServiceNumberOk() (*string, bool)`

GetEmergencyServiceNumberOk returns a tuple with the EmergencyServiceNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEmergencyServiceNumber

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetEmergencyServiceNumber(v string)`

SetEmergencyServiceNumber sets EmergencyServiceNumber field to given value.

### HasEmergencyServiceNumber

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasEmergencyServiceNumber() bool`

HasEmergencyServiceNumber returns a boolean if a field has been set.

### GetAnonCidBlockActivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetAnonCidBlockActivate() string`

GetAnonCidBlockActivate returns the AnonCidBlockActivate field if non-nil, zero value otherwise.

### GetAnonCidBlockActivateOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetAnonCidBlockActivateOk() (*string, bool)`

GetAnonCidBlockActivateOk returns a tuple with the AnonCidBlockActivate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnonCidBlockActivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetAnonCidBlockActivate(v string)`

SetAnonCidBlockActivate sets AnonCidBlockActivate field to given value.

### HasAnonCidBlockActivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasAnonCidBlockActivate() bool`

HasAnonCidBlockActivate returns a boolean if a field has been set.

### GetAnonCidBlockDeactivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetAnonCidBlockDeactivate() string`

GetAnonCidBlockDeactivate returns the AnonCidBlockDeactivate field if non-nil, zero value otherwise.

### GetAnonCidBlockDeactivateOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetAnonCidBlockDeactivateOk() (*string, bool)`

GetAnonCidBlockDeactivateOk returns a tuple with the AnonCidBlockDeactivate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnonCidBlockDeactivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetAnonCidBlockDeactivate(v string)`

SetAnonCidBlockDeactivate sets AnonCidBlockDeactivate field to given value.

### HasAnonCidBlockDeactivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasAnonCidBlockDeactivate() bool`

HasAnonCidBlockDeactivate returns a boolean if a field has been set.

### GetCallForwardUnconditionalActivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetCallForwardUnconditionalActivate() string`

GetCallForwardUnconditionalActivate returns the CallForwardUnconditionalActivate field if non-nil, zero value otherwise.

### GetCallForwardUnconditionalActivateOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetCallForwardUnconditionalActivateOk() (*string, bool)`

GetCallForwardUnconditionalActivateOk returns a tuple with the CallForwardUnconditionalActivate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallForwardUnconditionalActivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetCallForwardUnconditionalActivate(v string)`

SetCallForwardUnconditionalActivate sets CallForwardUnconditionalActivate field to given value.

### HasCallForwardUnconditionalActivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasCallForwardUnconditionalActivate() bool`

HasCallForwardUnconditionalActivate returns a boolean if a field has been set.

### GetCallForwardUnconditionalDeactivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetCallForwardUnconditionalDeactivate() string`

GetCallForwardUnconditionalDeactivate returns the CallForwardUnconditionalDeactivate field if non-nil, zero value otherwise.

### GetCallForwardUnconditionalDeactivateOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetCallForwardUnconditionalDeactivateOk() (*string, bool)`

GetCallForwardUnconditionalDeactivateOk returns a tuple with the CallForwardUnconditionalDeactivate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallForwardUnconditionalDeactivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetCallForwardUnconditionalDeactivate(v string)`

SetCallForwardUnconditionalDeactivate sets CallForwardUnconditionalDeactivate field to given value.

### HasCallForwardUnconditionalDeactivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasCallForwardUnconditionalDeactivate() bool`

HasCallForwardUnconditionalDeactivate returns a boolean if a field has been set.

### GetCallForwardOnBusyActivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetCallForwardOnBusyActivate() string`

GetCallForwardOnBusyActivate returns the CallForwardOnBusyActivate field if non-nil, zero value otherwise.

### GetCallForwardOnBusyActivateOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetCallForwardOnBusyActivateOk() (*string, bool)`

GetCallForwardOnBusyActivateOk returns a tuple with the CallForwardOnBusyActivate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallForwardOnBusyActivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetCallForwardOnBusyActivate(v string)`

SetCallForwardOnBusyActivate sets CallForwardOnBusyActivate field to given value.

### HasCallForwardOnBusyActivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasCallForwardOnBusyActivate() bool`

HasCallForwardOnBusyActivate returns a boolean if a field has been set.

### GetCallForwardOnBusyDeactivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetCallForwardOnBusyDeactivate() string`

GetCallForwardOnBusyDeactivate returns the CallForwardOnBusyDeactivate field if non-nil, zero value otherwise.

### GetCallForwardOnBusyDeactivateOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetCallForwardOnBusyDeactivateOk() (*string, bool)`

GetCallForwardOnBusyDeactivateOk returns a tuple with the CallForwardOnBusyDeactivate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallForwardOnBusyDeactivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetCallForwardOnBusyDeactivate(v string)`

SetCallForwardOnBusyDeactivate sets CallForwardOnBusyDeactivate field to given value.

### HasCallForwardOnBusyDeactivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasCallForwardOnBusyDeactivate() bool`

HasCallForwardOnBusyDeactivate returns a boolean if a field has been set.

### GetCallForwardOnNoAnswerActivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetCallForwardOnNoAnswerActivate() string`

GetCallForwardOnNoAnswerActivate returns the CallForwardOnNoAnswerActivate field if non-nil, zero value otherwise.

### GetCallForwardOnNoAnswerActivateOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetCallForwardOnNoAnswerActivateOk() (*string, bool)`

GetCallForwardOnNoAnswerActivateOk returns a tuple with the CallForwardOnNoAnswerActivate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallForwardOnNoAnswerActivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetCallForwardOnNoAnswerActivate(v string)`

SetCallForwardOnNoAnswerActivate sets CallForwardOnNoAnswerActivate field to given value.

### HasCallForwardOnNoAnswerActivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasCallForwardOnNoAnswerActivate() bool`

HasCallForwardOnNoAnswerActivate returns a boolean if a field has been set.

### GetCallForwardOnNoAnswerDeactivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetCallForwardOnNoAnswerDeactivate() string`

GetCallForwardOnNoAnswerDeactivate returns the CallForwardOnNoAnswerDeactivate field if non-nil, zero value otherwise.

### GetCallForwardOnNoAnswerDeactivateOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetCallForwardOnNoAnswerDeactivateOk() (*string, bool)`

GetCallForwardOnNoAnswerDeactivateOk returns a tuple with the CallForwardOnNoAnswerDeactivate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCallForwardOnNoAnswerDeactivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetCallForwardOnNoAnswerDeactivate(v string)`

SetCallForwardOnNoAnswerDeactivate sets CallForwardOnNoAnswerDeactivate field to given value.

### HasCallForwardOnNoAnswerDeactivate

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasCallForwardOnNoAnswerDeactivate() bool`

HasCallForwardOnNoAnswerDeactivate returns a boolean if a field has been set.

### GetIntercom1

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetIntercom1() string`

GetIntercom1 returns the Intercom1 field if non-nil, zero value otherwise.

### GetIntercom1Ok

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetIntercom1Ok() (*string, bool)`

GetIntercom1Ok returns a tuple with the Intercom1 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIntercom1

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetIntercom1(v string)`

SetIntercom1 sets Intercom1 field to given value.

### HasIntercom1

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasIntercom1() bool`

HasIntercom1 returns a boolean if a field has been set.

### GetIntercom2

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetIntercom2() string`

GetIntercom2 returns the Intercom2 field if non-nil, zero value otherwise.

### GetIntercom2Ok

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetIntercom2Ok() (*string, bool)`

GetIntercom2Ok returns a tuple with the Intercom2 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIntercom2

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetIntercom2(v string)`

SetIntercom2 sets Intercom2 field to given value.

### HasIntercom2

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasIntercom2() bool`

HasIntercom2 returns a boolean if a field has been set.

### GetIntercom3

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetIntercom3() string`

GetIntercom3 returns the Intercom3 field if non-nil, zero value otherwise.

### GetIntercom3Ok

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetIntercom3Ok() (*string, bool)`

GetIntercom3Ok returns a tuple with the Intercom3 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIntercom3

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetIntercom3(v string)`

SetIntercom3 sets Intercom3 field to given value.

### HasIntercom3

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasIntercom3() bool`

HasIntercom3 returns a boolean if a field has been set.

### GetCodecs

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetCodecs() []ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInner`

GetCodecs returns the Codecs field if non-nil, zero value otherwise.

### GetCodecsOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetCodecsOk() (*[]ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInner, bool)`

GetCodecsOk returns a tuple with the Codecs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCodecs

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetCodecs(v []ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsNameCodecsInner)`

SetCodecs sets Codecs field to given value.

### HasCodecs

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasCodecs() bool`

HasCodecs returns a boolean if a field has been set.

### GetObjectProperties

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetObjectProperties() ConfigPutRequestPacketQueuePacketQueueNameObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) GetObjectPropertiesOk() (*ConfigPutRequestPacketQueuePacketQueueNameObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) SetObjectProperties(v ConfigPutRequestPacketQueuePacketQueueNameObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


