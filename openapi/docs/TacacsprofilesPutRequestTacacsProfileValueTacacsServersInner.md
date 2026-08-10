# TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Enabled** | Pointer to **bool** | Enable TACACS+ server | [optional] [default to false]
**Server** | Pointer to **string** | IPv4, IPv6, or DNS name for TACACS+ server | [optional] [default to ""]
**AuthType** | Pointer to **string** | TACACS+ authentication type | [optional] [default to "pap"]
**Port** | Pointer to **string** | TACACS+ server port | [optional] [default to ""]
**Timeout** | Pointer to **NullableInt64** | TACACS+ server timeout in seconds | [optional] 
**Secret** | Pointer to **string** | TACACS+ shared secret | [optional] [default to ""]
**EncSecret** | Pointer to **string** | TACACS+ shared secret (encrypted) | [optional] [default to ""]
**Index** | Pointer to **int64** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 

## Methods

### NewTacacsprofilesPutRequestTacacsProfileValueTacacsServersInner

`func NewTacacsprofilesPutRequestTacacsProfileValueTacacsServersInner() *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner`

NewTacacsprofilesPutRequestTacacsProfileValueTacacsServersInner instantiates a new TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTacacsprofilesPutRequestTacacsProfileValueTacacsServersInnerWithDefaults

`func NewTacacsprofilesPutRequestTacacsProfileValueTacacsServersInnerWithDefaults() *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner`

NewTacacsprofilesPutRequestTacacsProfileValueTacacsServersInnerWithDefaults instantiates a new TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEnabled

`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) GetEnabled() bool`

GetEnabled returns the Enabled field if non-nil, zero value otherwise.

### GetEnabledOk

`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) GetEnabledOk() (*bool, bool)`

GetEnabledOk returns a tuple with the Enabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnabled

`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) SetEnabled(v bool)`

SetEnabled sets Enabled field to given value.

### HasEnabled

`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) HasEnabled() bool`

HasEnabled returns a boolean if a field has been set.

### GetServer

`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) GetServer() string`

GetServer returns the Server field if non-nil, zero value otherwise.

### GetServerOk

`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) GetServerOk() (*string, bool)`

GetServerOk returns a tuple with the Server field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServer

`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) SetServer(v string)`

SetServer sets Server field to given value.

### HasServer

`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) HasServer() bool`

HasServer returns a boolean if a field has been set.

### GetAuthType

`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) GetAuthType() string`

GetAuthType returns the AuthType field if non-nil, zero value otherwise.

### GetAuthTypeOk

`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) GetAuthTypeOk() (*string, bool)`

GetAuthTypeOk returns a tuple with the AuthType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthType

`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) SetAuthType(v string)`

SetAuthType sets AuthType field to given value.

### HasAuthType

`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) HasAuthType() bool`

HasAuthType returns a boolean if a field has been set.

### GetPort

`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) GetPort() string`

GetPort returns the Port field if non-nil, zero value otherwise.

### GetPortOk

`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) GetPortOk() (*string, bool)`

GetPortOk returns a tuple with the Port field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPort

`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) SetPort(v string)`

SetPort sets Port field to given value.

### HasPort

`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) HasPort() bool`

HasPort returns a boolean if a field has been set.

### GetTimeout

`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) GetTimeout() int64`

GetTimeout returns the Timeout field if non-nil, zero value otherwise.

### GetTimeoutOk

`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) GetTimeoutOk() (*int64, bool)`

GetTimeoutOk returns a tuple with the Timeout field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTimeout

`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) SetTimeout(v int64)`

SetTimeout sets Timeout field to given value.

### HasTimeout

`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) HasTimeout() bool`

HasTimeout returns a boolean if a field has been set.

### SetTimeoutNil

`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) SetTimeoutNil(b bool)`

 SetTimeoutNil sets the value for Timeout to be an explicit nil

### UnsetTimeout
`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) UnsetTimeout()`

UnsetTimeout ensures that no value is present for Timeout, not even an explicit nil
### GetSecret

`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) GetSecret() string`

GetSecret returns the Secret field if non-nil, zero value otherwise.

### GetSecretOk

`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) GetSecretOk() (*string, bool)`

GetSecretOk returns a tuple with the Secret field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSecret

`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) SetSecret(v string)`

SetSecret sets Secret field to given value.

### HasSecret

`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) HasSecret() bool`

HasSecret returns a boolean if a field has been set.

### GetEncSecret

`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) GetEncSecret() string`

GetEncSecret returns the EncSecret field if non-nil, zero value otherwise.

### GetEncSecretOk

`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) GetEncSecretOk() (*string, bool)`

GetEncSecretOk returns a tuple with the EncSecret field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEncSecret

`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) SetEncSecret(v string)`

SetEncSecret sets EncSecret field to given value.

### HasEncSecret

`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) HasEncSecret() bool`

HasEncSecret returns a boolean if a field has been set.

### GetIndex

`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) GetIndex() int64`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) GetIndexOk() (*int64, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) SetIndex(v int64)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *TacacsprofilesPutRequestTacacsProfileValueTacacsServersInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


