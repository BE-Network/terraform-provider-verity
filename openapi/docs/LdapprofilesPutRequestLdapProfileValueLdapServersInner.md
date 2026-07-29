# LdapprofilesPutRequestLdapProfileValueLdapServersInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Enabled** | Pointer to **bool** | Enable this LDAP server entry | [optional] [default to false]
**Server** | Pointer to **string** | IPv4, IPv6, or DNS hostname for LDAP server | [optional] [default to ""]
**Port** | Pointer to **NullableInt32** | Server port (overrides global default port) | [optional] 
**UseType** | Pointer to **string** | Which LDAP client(s) use this server | [optional] [default to "all"]
**Priority** | Pointer to **NullableInt32** | Server priority (1-99, lower &#x3D; higher priority) | [optional] 
**SslTlsMode** | Pointer to **string** | Per-server TLS mode (overrides global setting) | [optional] [default to "off"]
**RetransmitAttempts** | Pointer to **NullableInt32** | Per-server retransmit attempts (0-10) | [optional] 
**Index** | Pointer to **int32** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 

## Methods

### NewLdapprofilesPutRequestLdapProfileValueLdapServersInner

`func NewLdapprofilesPutRequestLdapProfileValueLdapServersInner() *LdapprofilesPutRequestLdapProfileValueLdapServersInner`

NewLdapprofilesPutRequestLdapProfileValueLdapServersInner instantiates a new LdapprofilesPutRequestLdapProfileValueLdapServersInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewLdapprofilesPutRequestLdapProfileValueLdapServersInnerWithDefaults

`func NewLdapprofilesPutRequestLdapProfileValueLdapServersInnerWithDefaults() *LdapprofilesPutRequestLdapProfileValueLdapServersInner`

NewLdapprofilesPutRequestLdapProfileValueLdapServersInnerWithDefaults instantiates a new LdapprofilesPutRequestLdapProfileValueLdapServersInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEnabled

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) GetEnabled() bool`

GetEnabled returns the Enabled field if non-nil, zero value otherwise.

### GetEnabledOk

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) GetEnabledOk() (*bool, bool)`

GetEnabledOk returns a tuple with the Enabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnabled

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) SetEnabled(v bool)`

SetEnabled sets Enabled field to given value.

### HasEnabled

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) HasEnabled() bool`

HasEnabled returns a boolean if a field has been set.

### GetServer

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) GetServer() string`

GetServer returns the Server field if non-nil, zero value otherwise.

### GetServerOk

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) GetServerOk() (*string, bool)`

GetServerOk returns a tuple with the Server field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServer

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) SetServer(v string)`

SetServer sets Server field to given value.

### HasServer

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) HasServer() bool`

HasServer returns a boolean if a field has been set.

### GetPort

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) GetPort() int32`

GetPort returns the Port field if non-nil, zero value otherwise.

### GetPortOk

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) GetPortOk() (*int32, bool)`

GetPortOk returns a tuple with the Port field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPort

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) SetPort(v int32)`

SetPort sets Port field to given value.

### HasPort

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) HasPort() bool`

HasPort returns a boolean if a field has been set.

### SetPortNil

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) SetPortNil(b bool)`

 SetPortNil sets the value for Port to be an explicit nil

### UnsetPort
`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) UnsetPort()`

UnsetPort ensures that no value is present for Port, not even an explicit nil
### GetUseType

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) GetUseType() string`

GetUseType returns the UseType field if non-nil, zero value otherwise.

### GetUseTypeOk

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) GetUseTypeOk() (*string, bool)`

GetUseTypeOk returns a tuple with the UseType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUseType

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) SetUseType(v string)`

SetUseType sets UseType field to given value.

### HasUseType

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) HasUseType() bool`

HasUseType returns a boolean if a field has been set.

### GetPriority

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) GetPriority() int32`

GetPriority returns the Priority field if non-nil, zero value otherwise.

### GetPriorityOk

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) GetPriorityOk() (*int32, bool)`

GetPriorityOk returns a tuple with the Priority field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPriority

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) SetPriority(v int32)`

SetPriority sets Priority field to given value.

### HasPriority

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) HasPriority() bool`

HasPriority returns a boolean if a field has been set.

### SetPriorityNil

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) SetPriorityNil(b bool)`

 SetPriorityNil sets the value for Priority to be an explicit nil

### UnsetPriority
`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) UnsetPriority()`

UnsetPriority ensures that no value is present for Priority, not even an explicit nil
### GetSslTlsMode

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) GetSslTlsMode() string`

GetSslTlsMode returns the SslTlsMode field if non-nil, zero value otherwise.

### GetSslTlsModeOk

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) GetSslTlsModeOk() (*string, bool)`

GetSslTlsModeOk returns a tuple with the SslTlsMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSslTlsMode

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) SetSslTlsMode(v string)`

SetSslTlsMode sets SslTlsMode field to given value.

### HasSslTlsMode

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) HasSslTlsMode() bool`

HasSslTlsMode returns a boolean if a field has been set.

### GetRetransmitAttempts

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) GetRetransmitAttempts() int32`

GetRetransmitAttempts returns the RetransmitAttempts field if non-nil, zero value otherwise.

### GetRetransmitAttemptsOk

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) GetRetransmitAttemptsOk() (*int32, bool)`

GetRetransmitAttemptsOk returns a tuple with the RetransmitAttempts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRetransmitAttempts

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) SetRetransmitAttempts(v int32)`

SetRetransmitAttempts sets RetransmitAttempts field to given value.

### HasRetransmitAttempts

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) HasRetransmitAttempts() bool`

HasRetransmitAttempts returns a boolean if a field has been set.

### SetRetransmitAttemptsNil

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) SetRetransmitAttemptsNil(b bool)`

 SetRetransmitAttemptsNil sets the value for RetransmitAttempts to be an explicit nil

### UnsetRetransmitAttempts
`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) UnsetRetransmitAttempts()`

UnsetRetransmitAttempts ensures that no value is present for RetransmitAttempts, not even an explicit nil
### GetIndex

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) GetIndex() int32`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) GetIndexOk() (*int32, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) SetIndex(v int32)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *LdapprofilesPutRequestLdapProfileValueLdapServersInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


