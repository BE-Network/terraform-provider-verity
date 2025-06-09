# ConfigPutRequestServicePortProfileServicePortProfileName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**PortType** | Pointer to **string** | Determines what Service are provisioned on the port and if those Services are propagated upstream&lt;ul&gt;&lt;li&gt;* \&quot;Upstream Switchport\&quot; Services specified below.  Services are not propagated.&lt;/li&gt;&lt;li&gt;* \&quot;Downstream Switchport\&quot; Services specified below. Services are propagated.&lt;/li&gt;&lt;li&gt;* \&quot;Crosslink Switchport\&quot; Services is union of all Services on each switch.  Services are not propagated.&lt;/li&gt;&lt;li&gt;* \&quot;Upstream L3 (L2/L3 Switches Only\&quot; No Services.&lt;/li&gt;&lt;/ul&gt; | [optional] [default to "up"]
**TlsLimitIn** | Pointer to **NullableInt32** | Speed of ingress (Mbps) for TLS (Transparent LAN Service) | [optional] [default to 1000]
**TlsService** | Pointer to **string** | Service used for TLS (Transparent LAN Service) | [optional] [default to ""]
**TlsServiceRefType** | Pointer to **string** | Object type for tls_service field | [optional] 
**TrustedPort** | Pointer to **bool** | Trusted Ports do not participate in IP Source Guard, Dynamic ARP Inspection, nor DHCP Snooping, meaning all packets are forwarded without any checks. | [optional] [default to false]
**IpMask** | Pointer to **string** | IP/Mask | [optional] [default to ""]
**Services** | Pointer to [**[]ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner**](ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**ConfigPutRequestServicePortProfileServicePortProfileNameObjectProperties**](ConfigPutRequestServicePortProfileServicePortProfileNameObjectProperties.md) |  | [optional] 

## Methods

### NewConfigPutRequestServicePortProfileServicePortProfileName

`func NewConfigPutRequestServicePortProfileServicePortProfileName() *ConfigPutRequestServicePortProfileServicePortProfileName`

NewConfigPutRequestServicePortProfileServicePortProfileName instantiates a new ConfigPutRequestServicePortProfileServicePortProfileName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestServicePortProfileServicePortProfileNameWithDefaults

`func NewConfigPutRequestServicePortProfileServicePortProfileNameWithDefaults() *ConfigPutRequestServicePortProfileServicePortProfileName`

NewConfigPutRequestServicePortProfileServicePortProfileNameWithDefaults instantiates a new ConfigPutRequestServicePortProfileServicePortProfileName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetPortType

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) GetPortType() string`

GetPortType returns the PortType field if non-nil, zero value otherwise.

### GetPortTypeOk

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) GetPortTypeOk() (*string, bool)`

GetPortTypeOk returns a tuple with the PortType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPortType

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) SetPortType(v string)`

SetPortType sets PortType field to given value.

### HasPortType

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) HasPortType() bool`

HasPortType returns a boolean if a field has been set.

### GetTlsLimitIn

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) GetTlsLimitIn() int32`

GetTlsLimitIn returns the TlsLimitIn field if non-nil, zero value otherwise.

### GetTlsLimitInOk

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) GetTlsLimitInOk() (*int32, bool)`

GetTlsLimitInOk returns a tuple with the TlsLimitIn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTlsLimitIn

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) SetTlsLimitIn(v int32)`

SetTlsLimitIn sets TlsLimitIn field to given value.

### HasTlsLimitIn

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) HasTlsLimitIn() bool`

HasTlsLimitIn returns a boolean if a field has been set.

### SetTlsLimitInNil

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) SetTlsLimitInNil(b bool)`

 SetTlsLimitInNil sets the value for TlsLimitIn to be an explicit nil

### UnsetTlsLimitIn
`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) UnsetTlsLimitIn()`

UnsetTlsLimitIn ensures that no value is present for TlsLimitIn, not even an explicit nil
### GetTlsService

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) GetTlsService() string`

GetTlsService returns the TlsService field if non-nil, zero value otherwise.

### GetTlsServiceOk

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) GetTlsServiceOk() (*string, bool)`

GetTlsServiceOk returns a tuple with the TlsService field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTlsService

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) SetTlsService(v string)`

SetTlsService sets TlsService field to given value.

### HasTlsService

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) HasTlsService() bool`

HasTlsService returns a boolean if a field has been set.

### GetTlsServiceRefType

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) GetTlsServiceRefType() string`

GetTlsServiceRefType returns the TlsServiceRefType field if non-nil, zero value otherwise.

### GetTlsServiceRefTypeOk

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) GetTlsServiceRefTypeOk() (*string, bool)`

GetTlsServiceRefTypeOk returns a tuple with the TlsServiceRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTlsServiceRefType

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) SetTlsServiceRefType(v string)`

SetTlsServiceRefType sets TlsServiceRefType field to given value.

### HasTlsServiceRefType

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) HasTlsServiceRefType() bool`

HasTlsServiceRefType returns a boolean if a field has been set.

### GetTrustedPort

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) GetTrustedPort() bool`

GetTrustedPort returns the TrustedPort field if non-nil, zero value otherwise.

### GetTrustedPortOk

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) GetTrustedPortOk() (*bool, bool)`

GetTrustedPortOk returns a tuple with the TrustedPort field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTrustedPort

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) SetTrustedPort(v bool)`

SetTrustedPort sets TrustedPort field to given value.

### HasTrustedPort

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) HasTrustedPort() bool`

HasTrustedPort returns a boolean if a field has been set.

### GetIpMask

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) GetIpMask() string`

GetIpMask returns the IpMask field if non-nil, zero value otherwise.

### GetIpMaskOk

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) GetIpMaskOk() (*string, bool)`

GetIpMaskOk returns a tuple with the IpMask field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpMask

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) SetIpMask(v string)`

SetIpMask sets IpMask field to given value.

### HasIpMask

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) HasIpMask() bool`

HasIpMask returns a boolean if a field has been set.

### GetServices

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) GetServices() []ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner`

GetServices returns the Services field if non-nil, zero value otherwise.

### GetServicesOk

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) GetServicesOk() (*[]ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner, bool)`

GetServicesOk returns a tuple with the Services field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServices

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) SetServices(v []ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner)`

SetServices sets Services field to given value.

### HasServices

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) HasServices() bool`

HasServices returns a boolean if a field has been set.

### GetObjectProperties

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) GetObjectProperties() ConfigPutRequestServicePortProfileServicePortProfileNameObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) GetObjectPropertiesOk() (*ConfigPutRequestServicePortProfileServicePortProfileNameObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) SetObjectProperties(v ConfigPutRequestServicePortProfileServicePortProfileNameObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ConfigPutRequestServicePortProfileServicePortProfileName) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


