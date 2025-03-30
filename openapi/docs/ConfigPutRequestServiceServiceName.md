# ConfigPutRequestServiceServiceName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**Vlan** | Pointer to **int32** | A Value between 1 and 4096 | [optional] 
**Vni** | Pointer to **int32** | Indication of the outgoing VLAN layer 2 service | [optional] 
**VniAutoAssigned** | Pointer to **bool** | Whether or not the value in vni field has been automatically assigned or not. Set to false and change vni value to edit. | [optional] 
**Tenant** | Pointer to **string** | Tenant | [optional] [default to ""]
**TenantRefType** | Pointer to **string** | Object type for tenant field | [optional] 
**AnycastIpMask** | Pointer to **string** | Static anycast gateway address for service  | [optional] [default to ""]
**DhcpServerIp** | Pointer to **string** | IP address(s) of the DHCP server for service.  May have up to four separated by commas. | [optional] [default to ""]
**Mtu** | Pointer to **NullableInt32** | MTU (Maximum Transmission Unit) The size used by a switch to determine when large packets must be broken up into smaller packets for delivery. If mismatched within a single vlan network, can cause dropped packets. | [optional] 
**ObjectProperties** | Pointer to [**ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties**](ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties.md) |  | [optional] 

## Methods

### NewConfigPutRequestServiceServiceName

`func NewConfigPutRequestServiceServiceName() *ConfigPutRequestServiceServiceName`

NewConfigPutRequestServiceServiceName instantiates a new ConfigPutRequestServiceServiceName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestServiceServiceNameWithDefaults

`func NewConfigPutRequestServiceServiceNameWithDefaults() *ConfigPutRequestServiceServiceName`

NewConfigPutRequestServiceServiceNameWithDefaults instantiates a new ConfigPutRequestServiceServiceName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestServiceServiceName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestServiceServiceName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestServiceServiceName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestServiceServiceName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ConfigPutRequestServiceServiceName) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestServiceServiceName) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestServiceServiceName) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestServiceServiceName) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetVlan

`func (o *ConfigPutRequestServiceServiceName) GetVlan() int32`

GetVlan returns the Vlan field if non-nil, zero value otherwise.

### GetVlanOk

`func (o *ConfigPutRequestServiceServiceName) GetVlanOk() (*int32, bool)`

GetVlanOk returns a tuple with the Vlan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVlan

`func (o *ConfigPutRequestServiceServiceName) SetVlan(v int32)`

SetVlan sets Vlan field to given value.

### HasVlan

`func (o *ConfigPutRequestServiceServiceName) HasVlan() bool`

HasVlan returns a boolean if a field has been set.

### GetVni

`func (o *ConfigPutRequestServiceServiceName) GetVni() int32`

GetVni returns the Vni field if non-nil, zero value otherwise.

### GetVniOk

`func (o *ConfigPutRequestServiceServiceName) GetVniOk() (*int32, bool)`

GetVniOk returns a tuple with the Vni field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVni

`func (o *ConfigPutRequestServiceServiceName) SetVni(v int32)`

SetVni sets Vni field to given value.

### HasVni

`func (o *ConfigPutRequestServiceServiceName) HasVni() bool`

HasVni returns a boolean if a field has been set.

### GetVniAutoAssigned

`func (o *ConfigPutRequestServiceServiceName) GetVniAutoAssigned() bool`

GetVniAutoAssigned returns the VniAutoAssigned field if non-nil, zero value otherwise.

### GetVniAutoAssignedOk

`func (o *ConfigPutRequestServiceServiceName) GetVniAutoAssignedOk() (*bool, bool)`

GetVniAutoAssignedOk returns a tuple with the VniAutoAssigned field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVniAutoAssigned

`func (o *ConfigPutRequestServiceServiceName) SetVniAutoAssigned(v bool)`

SetVniAutoAssigned sets VniAutoAssigned field to given value.

### HasVniAutoAssigned

`func (o *ConfigPutRequestServiceServiceName) HasVniAutoAssigned() bool`

HasVniAutoAssigned returns a boolean if a field has been set.

### GetTenant

`func (o *ConfigPutRequestServiceServiceName) GetTenant() string`

GetTenant returns the Tenant field if non-nil, zero value otherwise.

### GetTenantOk

`func (o *ConfigPutRequestServiceServiceName) GetTenantOk() (*string, bool)`

GetTenantOk returns a tuple with the Tenant field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenant

`func (o *ConfigPutRequestServiceServiceName) SetTenant(v string)`

SetTenant sets Tenant field to given value.

### HasTenant

`func (o *ConfigPutRequestServiceServiceName) HasTenant() bool`

HasTenant returns a boolean if a field has been set.

### GetTenantRefType

`func (o *ConfigPutRequestServiceServiceName) GetTenantRefType() string`

GetTenantRefType returns the TenantRefType field if non-nil, zero value otherwise.

### GetTenantRefTypeOk

`func (o *ConfigPutRequestServiceServiceName) GetTenantRefTypeOk() (*string, bool)`

GetTenantRefTypeOk returns a tuple with the TenantRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantRefType

`func (o *ConfigPutRequestServiceServiceName) SetTenantRefType(v string)`

SetTenantRefType sets TenantRefType field to given value.

### HasTenantRefType

`func (o *ConfigPutRequestServiceServiceName) HasTenantRefType() bool`

HasTenantRefType returns a boolean if a field has been set.

### GetAnycastIpMask

`func (o *ConfigPutRequestServiceServiceName) GetAnycastIpMask() string`

GetAnycastIpMask returns the AnycastIpMask field if non-nil, zero value otherwise.

### GetAnycastIpMaskOk

`func (o *ConfigPutRequestServiceServiceName) GetAnycastIpMaskOk() (*string, bool)`

GetAnycastIpMaskOk returns a tuple with the AnycastIpMask field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnycastIpMask

`func (o *ConfigPutRequestServiceServiceName) SetAnycastIpMask(v string)`

SetAnycastIpMask sets AnycastIpMask field to given value.

### HasAnycastIpMask

`func (o *ConfigPutRequestServiceServiceName) HasAnycastIpMask() bool`

HasAnycastIpMask returns a boolean if a field has been set.

### GetDhcpServerIp

`func (o *ConfigPutRequestServiceServiceName) GetDhcpServerIp() string`

GetDhcpServerIp returns the DhcpServerIp field if non-nil, zero value otherwise.

### GetDhcpServerIpOk

`func (o *ConfigPutRequestServiceServiceName) GetDhcpServerIpOk() (*string, bool)`

GetDhcpServerIpOk returns a tuple with the DhcpServerIp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDhcpServerIp

`func (o *ConfigPutRequestServiceServiceName) SetDhcpServerIp(v string)`

SetDhcpServerIp sets DhcpServerIp field to given value.

### HasDhcpServerIp

`func (o *ConfigPutRequestServiceServiceName) HasDhcpServerIp() bool`

HasDhcpServerIp returns a boolean if a field has been set.

### GetMtu

`func (o *ConfigPutRequestServiceServiceName) GetMtu() int32`

GetMtu returns the Mtu field if non-nil, zero value otherwise.

### GetMtuOk

`func (o *ConfigPutRequestServiceServiceName) GetMtuOk() (*int32, bool)`

GetMtuOk returns a tuple with the Mtu field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMtu

`func (o *ConfigPutRequestServiceServiceName) SetMtu(v int32)`

SetMtu sets Mtu field to given value.

### HasMtu

`func (o *ConfigPutRequestServiceServiceName) HasMtu() bool`

HasMtu returns a boolean if a field has been set.

### SetMtuNil

`func (o *ConfigPutRequestServiceServiceName) SetMtuNil(b bool)`

 SetMtuNil sets the value for Mtu to be an explicit nil

### UnsetMtu
`func (o *ConfigPutRequestServiceServiceName) UnsetMtu()`

UnsetMtu ensures that no value is present for Mtu, not even an explicit nil
### GetObjectProperties

`func (o *ConfigPutRequestServiceServiceName) GetObjectProperties() ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ConfigPutRequestServiceServiceName) GetObjectPropertiesOk() (*ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ConfigPutRequestServiceServiceName) SetObjectProperties(v ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ConfigPutRequestServiceServiceName) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


