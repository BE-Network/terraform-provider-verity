# ServicesPutRequestServiceValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. It&#39;s highly recommended to set this value to true so that validation on the object will be ran. | [optional] [default to false]
**Vlan** | Pointer to **NullableInt32** | A Value between 1 and 4096 | [optional] 
**Vni** | Pointer to **NullableInt32** | Indication of the outgoing VLAN layer 2 service | [optional] 
**VniAutoAssigned** | Pointer to **bool** | Whether or not the value in vni field has been automatically assigned or not. Set to false and change vni value to edit. | [optional] 
**Tenant** | Pointer to **string** | Tenant | [optional] [default to ""]
**TenantRefType** | Pointer to **string** | Object type for tenant field | [optional] 
**AnycastIpMask** | Pointer to **string** | Static anycast gateway address for service  | [optional] [default to ""]
**DhcpServerIp** | Pointer to **string** | IP address(s) of the DHCP server for service.  May have up to four separated by commas. | [optional] [default to ""]
**Mtu** | Pointer to **NullableInt32** | MTU (Maximum Transmission Unit) The size used by a switch to determine when large packets must be broken up into smaller packets for delivery. If mismatched within a single vlan network, can cause dropped packets. | [optional] [default to 1500]
**ObjectProperties** | Pointer to [**EthportsettingsPutRequestEthPortSettingsValueObjectProperties**](EthportsettingsPutRequestEthPortSettingsValueObjectProperties.md) |  | [optional] 

## Methods

### NewServicesPutRequestServiceValue

`func NewServicesPutRequestServiceValue() *ServicesPutRequestServiceValue`

NewServicesPutRequestServiceValue instantiates a new ServicesPutRequestServiceValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewServicesPutRequestServiceValueWithDefaults

`func NewServicesPutRequestServiceValueWithDefaults() *ServicesPutRequestServiceValue`

NewServicesPutRequestServiceValueWithDefaults instantiates a new ServicesPutRequestServiceValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ServicesPutRequestServiceValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ServicesPutRequestServiceValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ServicesPutRequestServiceValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ServicesPutRequestServiceValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ServicesPutRequestServiceValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ServicesPutRequestServiceValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ServicesPutRequestServiceValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ServicesPutRequestServiceValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetVlan

`func (o *ServicesPutRequestServiceValue) GetVlan() int32`

GetVlan returns the Vlan field if non-nil, zero value otherwise.

### GetVlanOk

`func (o *ServicesPutRequestServiceValue) GetVlanOk() (*int32, bool)`

GetVlanOk returns a tuple with the Vlan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVlan

`func (o *ServicesPutRequestServiceValue) SetVlan(v int32)`

SetVlan sets Vlan field to given value.

### HasVlan

`func (o *ServicesPutRequestServiceValue) HasVlan() bool`

HasVlan returns a boolean if a field has been set.

### SetVlanNil

`func (o *ServicesPutRequestServiceValue) SetVlanNil(b bool)`

 SetVlanNil sets the value for Vlan to be an explicit nil

### UnsetVlan
`func (o *ServicesPutRequestServiceValue) UnsetVlan()`

UnsetVlan ensures that no value is present for Vlan, not even an explicit nil
### GetVni

`func (o *ServicesPutRequestServiceValue) GetVni() int32`

GetVni returns the Vni field if non-nil, zero value otherwise.

### GetVniOk

`func (o *ServicesPutRequestServiceValue) GetVniOk() (*int32, bool)`

GetVniOk returns a tuple with the Vni field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVni

`func (o *ServicesPutRequestServiceValue) SetVni(v int32)`

SetVni sets Vni field to given value.

### HasVni

`func (o *ServicesPutRequestServiceValue) HasVni() bool`

HasVni returns a boolean if a field has been set.

### SetVniNil

`func (o *ServicesPutRequestServiceValue) SetVniNil(b bool)`

 SetVniNil sets the value for Vni to be an explicit nil

### UnsetVni
`func (o *ServicesPutRequestServiceValue) UnsetVni()`

UnsetVni ensures that no value is present for Vni, not even an explicit nil
### GetVniAutoAssigned

`func (o *ServicesPutRequestServiceValue) GetVniAutoAssigned() bool`

GetVniAutoAssigned returns the VniAutoAssigned field if non-nil, zero value otherwise.

### GetVniAutoAssignedOk

`func (o *ServicesPutRequestServiceValue) GetVniAutoAssignedOk() (*bool, bool)`

GetVniAutoAssignedOk returns a tuple with the VniAutoAssigned field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVniAutoAssigned

`func (o *ServicesPutRequestServiceValue) SetVniAutoAssigned(v bool)`

SetVniAutoAssigned sets VniAutoAssigned field to given value.

### HasVniAutoAssigned

`func (o *ServicesPutRequestServiceValue) HasVniAutoAssigned() bool`

HasVniAutoAssigned returns a boolean if a field has been set.

### GetTenant

`func (o *ServicesPutRequestServiceValue) GetTenant() string`

GetTenant returns the Tenant field if non-nil, zero value otherwise.

### GetTenantOk

`func (o *ServicesPutRequestServiceValue) GetTenantOk() (*string, bool)`

GetTenantOk returns a tuple with the Tenant field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenant

`func (o *ServicesPutRequestServiceValue) SetTenant(v string)`

SetTenant sets Tenant field to given value.

### HasTenant

`func (o *ServicesPutRequestServiceValue) HasTenant() bool`

HasTenant returns a boolean if a field has been set.

### GetTenantRefType

`func (o *ServicesPutRequestServiceValue) GetTenantRefType() string`

GetTenantRefType returns the TenantRefType field if non-nil, zero value otherwise.

### GetTenantRefTypeOk

`func (o *ServicesPutRequestServiceValue) GetTenantRefTypeOk() (*string, bool)`

GetTenantRefTypeOk returns a tuple with the TenantRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantRefType

`func (o *ServicesPutRequestServiceValue) SetTenantRefType(v string)`

SetTenantRefType sets TenantRefType field to given value.

### HasTenantRefType

`func (o *ServicesPutRequestServiceValue) HasTenantRefType() bool`

HasTenantRefType returns a boolean if a field has been set.

### GetAnycastIpMask

`func (o *ServicesPutRequestServiceValue) GetAnycastIpMask() string`

GetAnycastIpMask returns the AnycastIpMask field if non-nil, zero value otherwise.

### GetAnycastIpMaskOk

`func (o *ServicesPutRequestServiceValue) GetAnycastIpMaskOk() (*string, bool)`

GetAnycastIpMaskOk returns a tuple with the AnycastIpMask field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnycastIpMask

`func (o *ServicesPutRequestServiceValue) SetAnycastIpMask(v string)`

SetAnycastIpMask sets AnycastIpMask field to given value.

### HasAnycastIpMask

`func (o *ServicesPutRequestServiceValue) HasAnycastIpMask() bool`

HasAnycastIpMask returns a boolean if a field has been set.

### GetDhcpServerIp

`func (o *ServicesPutRequestServiceValue) GetDhcpServerIp() string`

GetDhcpServerIp returns the DhcpServerIp field if non-nil, zero value otherwise.

### GetDhcpServerIpOk

`func (o *ServicesPutRequestServiceValue) GetDhcpServerIpOk() (*string, bool)`

GetDhcpServerIpOk returns a tuple with the DhcpServerIp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDhcpServerIp

`func (o *ServicesPutRequestServiceValue) SetDhcpServerIp(v string)`

SetDhcpServerIp sets DhcpServerIp field to given value.

### HasDhcpServerIp

`func (o *ServicesPutRequestServiceValue) HasDhcpServerIp() bool`

HasDhcpServerIp returns a boolean if a field has been set.

### GetMtu

`func (o *ServicesPutRequestServiceValue) GetMtu() int32`

GetMtu returns the Mtu field if non-nil, zero value otherwise.

### GetMtuOk

`func (o *ServicesPutRequestServiceValue) GetMtuOk() (*int32, bool)`

GetMtuOk returns a tuple with the Mtu field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMtu

`func (o *ServicesPutRequestServiceValue) SetMtu(v int32)`

SetMtu sets Mtu field to given value.

### HasMtu

`func (o *ServicesPutRequestServiceValue) HasMtu() bool`

HasMtu returns a boolean if a field has been set.

### SetMtuNil

`func (o *ServicesPutRequestServiceValue) SetMtuNil(b bool)`

 SetMtuNil sets the value for Mtu to be an explicit nil

### UnsetMtu
`func (o *ServicesPutRequestServiceValue) UnsetMtu()`

UnsetMtu ensures that no value is present for Mtu, not even an explicit nil
### GetObjectProperties

`func (o *ServicesPutRequestServiceValue) GetObjectProperties() EthportsettingsPutRequestEthPortSettingsValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ServicesPutRequestServiceValue) GetObjectPropertiesOk() (*EthportsettingsPutRequestEthPortSettingsValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ServicesPutRequestServiceValue) SetObjectProperties(v EthportsettingsPutRequestEthPortSettingsValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ServicesPutRequestServiceValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


