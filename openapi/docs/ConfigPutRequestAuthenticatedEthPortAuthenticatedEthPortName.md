# ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**ConnectionMode** | Pointer to **string** | Choose connection mode for Authenticated Eth-Port&lt;br&gt;&lt;b&gt;Port Mode&lt;/b&gt;  Standard mode. The last authenticated clients VLAN access is applied.&lt;br&gt;&lt;b&gt;Single Client Mode&lt;/b&gt;  MAC filtered client. Only the authenticated clients traffic can pass. No traffic from a second client may pass. Only when the first client deauthenticates can a new authentication take place.&lt;br&gt;&lt;b&gt;Multiple Client Mode&lt;/b&gt;  MAC filtered clients. Only authenticated client traffic can pass. Multiple clients can authenticate and gain access to individual service offerings. MAC-based authentication is not supported. | [optional] [default to "PortMode"]
**ReauthorizationPeriodSec** | Pointer to **int32** | Amount of time in seconds before 802.1X requires reauthorization of an active session. \&quot;0\&quot; disables reauthorization (not recommended) | [optional] [default to 3600]
**AllowMacBasedAuthentication** | Pointer to **bool** | Enables 802.1x to capture the connected MAC address and send it tothe Radius Server instead of requesting credentials.  Useful for printers and similar devices | [optional] [default to false]
**MacAuthenticationHoldoffSec** | Pointer to **int32** | Amount of time in seconds 802.1X authentication is allowed to run before MAC-based authentication has begun | [optional] [default to 60]
**TrustedPort** | Pointer to **bool** | Trusted Ports do not participate in IP Source Guard, Dynamic ARP Inspection, nor DHCP Snooping, meaning all packets are forwarded without any checks. | [optional] [default to false]
**EthPorts** | Pointer to [**[]ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortNameEthPortsInner**](ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortNameEthPortsInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortNameObjectProperties**](ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortNameObjectProperties.md) |  | [optional] 

## Methods

### NewConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName

`func NewConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName() *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName`

NewConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName instantiates a new ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortNameWithDefaults

`func NewConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortNameWithDefaults() *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName`

NewConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortNameWithDefaults instantiates a new ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetConnectionMode

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) GetConnectionMode() string`

GetConnectionMode returns the ConnectionMode field if non-nil, zero value otherwise.

### GetConnectionModeOk

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) GetConnectionModeOk() (*string, bool)`

GetConnectionModeOk returns a tuple with the ConnectionMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnectionMode

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) SetConnectionMode(v string)`

SetConnectionMode sets ConnectionMode field to given value.

### HasConnectionMode

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) HasConnectionMode() bool`

HasConnectionMode returns a boolean if a field has been set.

### GetReauthorizationPeriodSec

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) GetReauthorizationPeriodSec() int32`

GetReauthorizationPeriodSec returns the ReauthorizationPeriodSec field if non-nil, zero value otherwise.

### GetReauthorizationPeriodSecOk

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) GetReauthorizationPeriodSecOk() (*int32, bool)`

GetReauthorizationPeriodSecOk returns a tuple with the ReauthorizationPeriodSec field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReauthorizationPeriodSec

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) SetReauthorizationPeriodSec(v int32)`

SetReauthorizationPeriodSec sets ReauthorizationPeriodSec field to given value.

### HasReauthorizationPeriodSec

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) HasReauthorizationPeriodSec() bool`

HasReauthorizationPeriodSec returns a boolean if a field has been set.

### GetAllowMacBasedAuthentication

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) GetAllowMacBasedAuthentication() bool`

GetAllowMacBasedAuthentication returns the AllowMacBasedAuthentication field if non-nil, zero value otherwise.

### GetAllowMacBasedAuthenticationOk

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) GetAllowMacBasedAuthenticationOk() (*bool, bool)`

GetAllowMacBasedAuthenticationOk returns a tuple with the AllowMacBasedAuthentication field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllowMacBasedAuthentication

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) SetAllowMacBasedAuthentication(v bool)`

SetAllowMacBasedAuthentication sets AllowMacBasedAuthentication field to given value.

### HasAllowMacBasedAuthentication

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) HasAllowMacBasedAuthentication() bool`

HasAllowMacBasedAuthentication returns a boolean if a field has been set.

### GetMacAuthenticationHoldoffSec

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) GetMacAuthenticationHoldoffSec() int32`

GetMacAuthenticationHoldoffSec returns the MacAuthenticationHoldoffSec field if non-nil, zero value otherwise.

### GetMacAuthenticationHoldoffSecOk

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) GetMacAuthenticationHoldoffSecOk() (*int32, bool)`

GetMacAuthenticationHoldoffSecOk returns a tuple with the MacAuthenticationHoldoffSec field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMacAuthenticationHoldoffSec

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) SetMacAuthenticationHoldoffSec(v int32)`

SetMacAuthenticationHoldoffSec sets MacAuthenticationHoldoffSec field to given value.

### HasMacAuthenticationHoldoffSec

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) HasMacAuthenticationHoldoffSec() bool`

HasMacAuthenticationHoldoffSec returns a boolean if a field has been set.

### GetTrustedPort

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) GetTrustedPort() bool`

GetTrustedPort returns the TrustedPort field if non-nil, zero value otherwise.

### GetTrustedPortOk

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) GetTrustedPortOk() (*bool, bool)`

GetTrustedPortOk returns a tuple with the TrustedPort field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTrustedPort

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) SetTrustedPort(v bool)`

SetTrustedPort sets TrustedPort field to given value.

### HasTrustedPort

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) HasTrustedPort() bool`

HasTrustedPort returns a boolean if a field has been set.

### GetEthPorts

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) GetEthPorts() []ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortNameEthPortsInner`

GetEthPorts returns the EthPorts field if non-nil, zero value otherwise.

### GetEthPortsOk

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) GetEthPortsOk() (*[]ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortNameEthPortsInner, bool)`

GetEthPortsOk returns a tuple with the EthPorts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEthPorts

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) SetEthPorts(v []ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortNameEthPortsInner)`

SetEthPorts sets EthPorts field to given value.

### HasEthPorts

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) HasEthPorts() bool`

HasEthPorts returns a boolean if a field has been set.

### GetObjectProperties

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) GetObjectProperties() ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortNameObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) GetObjectPropertiesOk() (*ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortNameObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) SetObjectProperties(v ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortNameObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


