# ServiceportprofilesPutRequestServicePortProfileValue

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
**Services** | Pointer to [**[]ServiceportprofilesPutRequestServicePortProfileValueServicesInner**](ServiceportprofilesPutRequestServicePortProfileValueServicesInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**ServiceportprofilesPutRequestServicePortProfileValueObjectProperties**](ServiceportprofilesPutRequestServicePortProfileValueObjectProperties.md) |  | [optional] 

## Methods

### NewServiceportprofilesPutRequestServicePortProfileValue

`func NewServiceportprofilesPutRequestServicePortProfileValue() *ServiceportprofilesPutRequestServicePortProfileValue`

NewServiceportprofilesPutRequestServicePortProfileValue instantiates a new ServiceportprofilesPutRequestServicePortProfileValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewServiceportprofilesPutRequestServicePortProfileValueWithDefaults

`func NewServiceportprofilesPutRequestServicePortProfileValueWithDefaults() *ServiceportprofilesPutRequestServicePortProfileValue`

NewServiceportprofilesPutRequestServicePortProfileValueWithDefaults instantiates a new ServiceportprofilesPutRequestServicePortProfileValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetPortType

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) GetPortType() string`

GetPortType returns the PortType field if non-nil, zero value otherwise.

### GetPortTypeOk

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) GetPortTypeOk() (*string, bool)`

GetPortTypeOk returns a tuple with the PortType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPortType

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) SetPortType(v string)`

SetPortType sets PortType field to given value.

### HasPortType

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) HasPortType() bool`

HasPortType returns a boolean if a field has been set.

### GetTlsLimitIn

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) GetTlsLimitIn() int32`

GetTlsLimitIn returns the TlsLimitIn field if non-nil, zero value otherwise.

### GetTlsLimitInOk

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) GetTlsLimitInOk() (*int32, bool)`

GetTlsLimitInOk returns a tuple with the TlsLimitIn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTlsLimitIn

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) SetTlsLimitIn(v int32)`

SetTlsLimitIn sets TlsLimitIn field to given value.

### HasTlsLimitIn

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) HasTlsLimitIn() bool`

HasTlsLimitIn returns a boolean if a field has been set.

### SetTlsLimitInNil

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) SetTlsLimitInNil(b bool)`

 SetTlsLimitInNil sets the value for TlsLimitIn to be an explicit nil

### UnsetTlsLimitIn
`func (o *ServiceportprofilesPutRequestServicePortProfileValue) UnsetTlsLimitIn()`

UnsetTlsLimitIn ensures that no value is present for TlsLimitIn, not even an explicit nil
### GetTlsService

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) GetTlsService() string`

GetTlsService returns the TlsService field if non-nil, zero value otherwise.

### GetTlsServiceOk

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) GetTlsServiceOk() (*string, bool)`

GetTlsServiceOk returns a tuple with the TlsService field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTlsService

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) SetTlsService(v string)`

SetTlsService sets TlsService field to given value.

### HasTlsService

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) HasTlsService() bool`

HasTlsService returns a boolean if a field has been set.

### GetTlsServiceRefType

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) GetTlsServiceRefType() string`

GetTlsServiceRefType returns the TlsServiceRefType field if non-nil, zero value otherwise.

### GetTlsServiceRefTypeOk

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) GetTlsServiceRefTypeOk() (*string, bool)`

GetTlsServiceRefTypeOk returns a tuple with the TlsServiceRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTlsServiceRefType

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) SetTlsServiceRefType(v string)`

SetTlsServiceRefType sets TlsServiceRefType field to given value.

### HasTlsServiceRefType

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) HasTlsServiceRefType() bool`

HasTlsServiceRefType returns a boolean if a field has been set.

### GetTrustedPort

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) GetTrustedPort() bool`

GetTrustedPort returns the TrustedPort field if non-nil, zero value otherwise.

### GetTrustedPortOk

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) GetTrustedPortOk() (*bool, bool)`

GetTrustedPortOk returns a tuple with the TrustedPort field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTrustedPort

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) SetTrustedPort(v bool)`

SetTrustedPort sets TrustedPort field to given value.

### HasTrustedPort

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) HasTrustedPort() bool`

HasTrustedPort returns a boolean if a field has been set.

### GetIpMask

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) GetIpMask() string`

GetIpMask returns the IpMask field if non-nil, zero value otherwise.

### GetIpMaskOk

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) GetIpMaskOk() (*string, bool)`

GetIpMaskOk returns a tuple with the IpMask field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpMask

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) SetIpMask(v string)`

SetIpMask sets IpMask field to given value.

### HasIpMask

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) HasIpMask() bool`

HasIpMask returns a boolean if a field has been set.

### GetServices

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) GetServices() []ServiceportprofilesPutRequestServicePortProfileValueServicesInner`

GetServices returns the Services field if non-nil, zero value otherwise.

### GetServicesOk

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) GetServicesOk() (*[]ServiceportprofilesPutRequestServicePortProfileValueServicesInner, bool)`

GetServicesOk returns a tuple with the Services field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServices

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) SetServices(v []ServiceportprofilesPutRequestServicePortProfileValueServicesInner)`

SetServices sets Services field to given value.

### HasServices

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) HasServices() bool`

HasServices returns a boolean if a field has been set.

### GetObjectProperties

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) GetObjectProperties() ServiceportprofilesPutRequestServicePortProfileValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) GetObjectPropertiesOk() (*ServiceportprofilesPutRequestServicePortProfileValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) SetObjectProperties(v ServiceportprofilesPutRequestServicePortProfileValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ServiceportprofilesPutRequestServicePortProfileValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


