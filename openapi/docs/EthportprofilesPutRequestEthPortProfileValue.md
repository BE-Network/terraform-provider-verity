# EthportprofilesPutRequestEthPortProfileValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. It&#39;s highly recommended to set this value to true so that validation on the object will be ran. | [optional] [default to false]
**TenantSliceManaged** | Pointer to **bool** | Profiles that Tenant Slice creates and manages | [optional] [default to false]
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

### GetTenantSliceManaged

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetTenantSliceManaged() bool`

GetTenantSliceManaged returns the TenantSliceManaged field if non-nil, zero value otherwise.

### GetTenantSliceManagedOk

`func (o *EthportprofilesPutRequestEthPortProfileValue) GetTenantSliceManagedOk() (*bool, bool)`

GetTenantSliceManagedOk returns a tuple with the TenantSliceManaged field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantSliceManaged

`func (o *EthportprofilesPutRequestEthPortProfileValue) SetTenantSliceManaged(v bool)`

SetTenantSliceManaged sets TenantSliceManaged field to given value.

### HasTenantSliceManaged

`func (o *EthportprofilesPutRequestEthPortProfileValue) HasTenantSliceManaged() bool`

HasTenantSliceManaged returns a boolean if a field has been set.

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


