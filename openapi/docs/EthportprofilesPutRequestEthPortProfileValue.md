# EthportprofilesPutRequestEthPortProfileValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. It&#39;s highly recommended to set this value to true so that validation on the object will be ran. | [optional] [default to false]
**IngressAcl** | Pointer to **string** | Choose an ingress access control list | [optional] [default to ""]
**IngressAclRefType** | Pointer to **string** | Object type for ingress_acl field | [optional] 
**EgressAcl** | Pointer to **string** | Choose an egress access control list | [optional] [default to ""]
**EgressAclRefType** | Pointer to **string** | Object type for egress_acl field | [optional] 
**Services** | Pointer to [**[]EthportprofilesPutRequestEthPortProfileValueServicesInner**](EthportprofilesPutRequestEthPortProfileValueServicesInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**EthportprofilesPutRequestEthPortProfileValueObjectProperties**](EthportprofilesPutRequestEthPortProfileValueObjectProperties.md) |  | [optional] 
**Tls** | Pointer to **bool** | Transparent LAN Service Trunk | [optional] [default to false]
**TlsService** | Pointer to **string** | Choose a Service supporting Transparent LAN Service | [optional] [default to ""]
**TlsServiceRefType** | Pointer to **string** | Object type for tls_service field | [optional] 
**TrustedPort** | Pointer to **bool** | Trusted Ports do not participate in IP Source Guard, Dynamic ARP Inspection, nor DHCP Snooping, meaning all packets are forwarded without any checks. | [optional] [default to false]

## Methods

### NewEthportprofilesPutRequestEthPortProfileValue

`func NewEthportprofilesPutRequestEthPortProfileValue() *EthportprofilesPutRequestEthPortProfileValue`

NewEthportprofilesPutRequestEthPortProfileValue instantiates a new EthportprofilesPutRequestEthPortProfileValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewEthportprofilesPutRequestEthPortProfileValueWithDefaults

`func NewEthportprofilesPutRequestEthPortProfileValueWithDefaults() *EthportprofilesPutRequestEthPortProfileValue`

NewEthportprofilesPutRequestEthPortProfileValueWithDefaults instantiates a new EthportprofilesPutRequestEthPortProfileValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *EthportprofilesPutRequestEthPortProfileValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *EthportprofilesPutRequestEthPortProfileValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *EthportprofilesPutRequestEthPortProfileValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *EthportprofilesPutRequestEthPortProfileValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetIngressAcl

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetIngressAcl() string`

GetIngressAcl returns the IngressAcl field if non-nil, zero value otherwise.

### GetIngressAclOk

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetIngressAclOk() (*string, bool)`

GetIngressAclOk returns a tuple with the IngressAcl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIngressAcl

`func (o *EthportprofilesPutRequestEthPortProfileValue) SetIngressAcl(v string)`

SetIngressAcl sets IngressAcl field to given value.

### HasIngressAcl

`func (o *EthportprofilesPutRequestEthPortProfileValue) HasIngressAcl() bool`

HasIngressAcl returns a boolean if a field has been set.

### GetIngressAclRefType

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetIngressAclRefType() string`

GetIngressAclRefType returns the IngressAclRefType field if non-nil, zero value otherwise.

### GetIngressAclRefTypeOk

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetIngressAclRefTypeOk() (*string, bool)`

GetIngressAclRefTypeOk returns a tuple with the IngressAclRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIngressAclRefType

`func (o *EthportprofilesPutRequestEthPortProfileValue) SetIngressAclRefType(v string)`

SetIngressAclRefType sets IngressAclRefType field to given value.

### HasIngressAclRefType

`func (o *EthportprofilesPutRequestEthPortProfileValue) HasIngressAclRefType() bool`

HasIngressAclRefType returns a boolean if a field has been set.

### GetEgressAcl

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetEgressAcl() string`

GetEgressAcl returns the EgressAcl field if non-nil, zero value otherwise.

### GetEgressAclOk

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetEgressAclOk() (*string, bool)`

GetEgressAclOk returns a tuple with the EgressAcl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEgressAcl

`func (o *EthportprofilesPutRequestEthPortProfileValue) SetEgressAcl(v string)`

SetEgressAcl sets EgressAcl field to given value.

### HasEgressAcl

`func (o *EthportprofilesPutRequestEthPortProfileValue) HasEgressAcl() bool`

HasEgressAcl returns a boolean if a field has been set.

### GetEgressAclRefType

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetEgressAclRefType() string`

GetEgressAclRefType returns the EgressAclRefType field if non-nil, zero value otherwise.

### GetEgressAclRefTypeOk

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetEgressAclRefTypeOk() (*string, bool)`

GetEgressAclRefTypeOk returns a tuple with the EgressAclRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEgressAclRefType

`func (o *EthportprofilesPutRequestEthPortProfileValue) SetEgressAclRefType(v string)`

SetEgressAclRefType sets EgressAclRefType field to given value.

### HasEgressAclRefType

`func (o *EthportprofilesPutRequestEthPortProfileValue) HasEgressAclRefType() bool`

HasEgressAclRefType returns a boolean if a field has been set.

### GetServices

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetServices() []EthportprofilesPutRequestEthPortProfileValueServicesInner`

GetServices returns the Services field if non-nil, zero value otherwise.

### GetServicesOk

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetServicesOk() (*[]EthportprofilesPutRequestEthPortProfileValueServicesInner, bool)`

GetServicesOk returns a tuple with the Services field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServices

`func (o *EthportprofilesPutRequestEthPortProfileValue) SetServices(v []EthportprofilesPutRequestEthPortProfileValueServicesInner)`

SetServices sets Services field to given value.

### HasServices

`func (o *EthportprofilesPutRequestEthPortProfileValue) HasServices() bool`

HasServices returns a boolean if a field has been set.

### GetObjectProperties

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetObjectProperties() EthportprofilesPutRequestEthPortProfileValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetObjectPropertiesOk() (*EthportprofilesPutRequestEthPortProfileValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *EthportprofilesPutRequestEthPortProfileValue) SetObjectProperties(v EthportprofilesPutRequestEthPortProfileValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *EthportprofilesPutRequestEthPortProfileValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.

### GetTls

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetTls() bool`

GetTls returns the Tls field if non-nil, zero value otherwise.

### GetTlsOk

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetTlsOk() (*bool, bool)`

GetTlsOk returns a tuple with the Tls field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTls

`func (o *EthportprofilesPutRequestEthPortProfileValue) SetTls(v bool)`

SetTls sets Tls field to given value.

### HasTls

`func (o *EthportprofilesPutRequestEthPortProfileValue) HasTls() bool`

HasTls returns a boolean if a field has been set.

### GetTlsService

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetTlsService() string`

GetTlsService returns the TlsService field if non-nil, zero value otherwise.

### GetTlsServiceOk

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetTlsServiceOk() (*string, bool)`

GetTlsServiceOk returns a tuple with the TlsService field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTlsService

`func (o *EthportprofilesPutRequestEthPortProfileValue) SetTlsService(v string)`

SetTlsService sets TlsService field to given value.

### HasTlsService

`func (o *EthportprofilesPutRequestEthPortProfileValue) HasTlsService() bool`

HasTlsService returns a boolean if a field has been set.

### GetTlsServiceRefType

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetTlsServiceRefType() string`

GetTlsServiceRefType returns the TlsServiceRefType field if non-nil, zero value otherwise.

### GetTlsServiceRefTypeOk

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetTlsServiceRefTypeOk() (*string, bool)`

GetTlsServiceRefTypeOk returns a tuple with the TlsServiceRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTlsServiceRefType

`func (o *EthportprofilesPutRequestEthPortProfileValue) SetTlsServiceRefType(v string)`

SetTlsServiceRefType sets TlsServiceRefType field to given value.

### HasTlsServiceRefType

`func (o *EthportprofilesPutRequestEthPortProfileValue) HasTlsServiceRefType() bool`

HasTlsServiceRefType returns a boolean if a field has been set.

### GetTrustedPort

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetTrustedPort() bool`

GetTrustedPort returns the TrustedPort field if non-nil, zero value otherwise.

### GetTrustedPortOk

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetTrustedPortOk() (*bool, bool)`

GetTrustedPortOk returns a tuple with the TrustedPort field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTrustedPort

`func (o *EthportprofilesPutRequestEthPortProfileValue) SetTrustedPort(v bool)`

SetTrustedPort sets TrustedPort field to given value.

### HasTrustedPort

`func (o *EthportprofilesPutRequestEthPortProfileValue) HasTrustedPort() bool`

HasTrustedPort returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


