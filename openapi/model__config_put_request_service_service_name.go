/*
Verity API

This application demonstrates the usage of Verity API. 

API version: 2.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
)

// checks if the ConfigPutRequestServiceServiceName type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ConfigPutRequestServiceServiceName{}

// ConfigPutRequestServiceServiceName struct for ConfigPutRequestServiceServiceName
type ConfigPutRequestServiceServiceName struct {
	// Object Name. Must be unique.
	Name *string `json:"name,omitempty"`
	// Enable object.
	Enable *bool `json:"enable,omitempty"`
	// A Value between 1 and 4096
	Vlan NullableInt32 `json:"vlan,omitempty"`
	// Indication of the outgoing VLAN layer 2 service
	Vni NullableInt32 `json:"vni,omitempty"`
	// Whether or not the value in vni field has been automatically assigned or not. Set to false and change vni value to edit.
	VniAutoAssigned *bool `json:"vni_auto_assigned_,omitempty"`
	// Tenant
	Tenant *string `json:"tenant,omitempty"`
	// Object type for tenant field
	TenantRefType *string `json:"tenant_ref_type_,omitempty"`
	// Static anycast gateway address for service 
	AnycastIpMask *string `json:"anycast_ip_mask,omitempty"`
	// IP address(s) of the DHCP server for service.  May have up to four separated by commas.
	DhcpServerIp *string `json:"dhcp_server_ip,omitempty"`
	// MTU (Maximum Transmission Unit) The size used by a switch to determine when large packets must be broken up into smaller packets for delivery. If mismatched within a single vlan network, can cause dropped packets.
	Mtu NullableInt32 `json:"mtu,omitempty"`
	ObjectProperties *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties `json:"object_properties,omitempty"`
}

// NewConfigPutRequestServiceServiceName instantiates a new ConfigPutRequestServiceServiceName object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConfigPutRequestServiceServiceName() *ConfigPutRequestServiceServiceName {
	this := ConfigPutRequestServiceServiceName{}
	var name string = ""
	this.Name = &name
	var enable bool = false
	this.Enable = &enable
	var tenant string = ""
	this.Tenant = &tenant
	var anycastIpMask string = ""
	this.AnycastIpMask = &anycastIpMask
	var dhcpServerIp string = ""
	this.DhcpServerIp = &dhcpServerIp
	return &this
}

// NewConfigPutRequestServiceServiceNameWithDefaults instantiates a new ConfigPutRequestServiceServiceName object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConfigPutRequestServiceServiceNameWithDefaults() *ConfigPutRequestServiceServiceName {
	this := ConfigPutRequestServiceServiceName{}
	var name string = ""
	this.Name = &name
	var enable bool = false
	this.Enable = &enable
	var tenant string = ""
	this.Tenant = &tenant
	var anycastIpMask string = ""
	this.AnycastIpMask = &anycastIpMask
	var dhcpServerIp string = ""
	this.DhcpServerIp = &dhcpServerIp
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *ConfigPutRequestServiceServiceName) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestServiceServiceName) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *ConfigPutRequestServiceServiceName) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *ConfigPutRequestServiceServiceName) SetName(v string) {
	o.Name = &v
}

// GetEnable returns the Enable field value if set, zero value otherwise.
func (o *ConfigPutRequestServiceServiceName) GetEnable() bool {
	if o == nil || IsNil(o.Enable) {
		var ret bool
		return ret
	}
	return *o.Enable
}

// GetEnableOk returns a tuple with the Enable field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestServiceServiceName) GetEnableOk() (*bool, bool) {
	if o == nil || IsNil(o.Enable) {
		return nil, false
	}
	return o.Enable, true
}

// HasEnable returns a boolean if a field has been set.
func (o *ConfigPutRequestServiceServiceName) HasEnable() bool {
	if o != nil && !IsNil(o.Enable) {
		return true
	}

	return false
}

// SetEnable gets a reference to the given bool and assigns it to the Enable field.
func (o *ConfigPutRequestServiceServiceName) SetEnable(v bool) {
	o.Enable = &v
}

// GetVlan returns the Vlan field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *ConfigPutRequestServiceServiceName) GetVlan() int32 {
	if o == nil || IsNil(o.Vlan.Get()) {
		var ret int32
		return ret
	}
	return *o.Vlan.Get()
}

// GetVlanOk returns a tuple with the Vlan field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ConfigPutRequestServiceServiceName) GetVlanOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.Vlan.Get(), o.Vlan.IsSet()
}

// HasVlan returns a boolean if a field has been set.
func (o *ConfigPutRequestServiceServiceName) HasVlan() bool {
	if o != nil && o.Vlan.IsSet() {
		return true
	}

	return false
}

// SetVlan gets a reference to the given NullableInt32 and assigns it to the Vlan field.
func (o *ConfigPutRequestServiceServiceName) SetVlan(v int32) {
	o.Vlan.Set(&v)
}
// SetVlanNil sets the value for Vlan to be an explicit nil
func (o *ConfigPutRequestServiceServiceName) SetVlanNil() {
	o.Vlan.Set(nil)
}

// UnsetVlan ensures that no value is present for Vlan, not even an explicit nil
func (o *ConfigPutRequestServiceServiceName) UnsetVlan() {
	o.Vlan.Unset()
}

// GetVni returns the Vni field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *ConfigPutRequestServiceServiceName) GetVni() int32 {
	if o == nil || IsNil(o.Vni.Get()) {
		var ret int32
		return ret
	}
	return *o.Vni.Get()
}

// GetVniOk returns a tuple with the Vni field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ConfigPutRequestServiceServiceName) GetVniOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.Vni.Get(), o.Vni.IsSet()
}

// HasVni returns a boolean if a field has been set.
func (o *ConfigPutRequestServiceServiceName) HasVni() bool {
	if o != nil && o.Vni.IsSet() {
		return true
	}

	return false
}

// SetVni gets a reference to the given NullableInt32 and assigns it to the Vni field.
func (o *ConfigPutRequestServiceServiceName) SetVni(v int32) {
	o.Vni.Set(&v)
}
// SetVniNil sets the value for Vni to be an explicit nil
func (o *ConfigPutRequestServiceServiceName) SetVniNil() {
	o.Vni.Set(nil)
}

// UnsetVni ensures that no value is present for Vni, not even an explicit nil
func (o *ConfigPutRequestServiceServiceName) UnsetVni() {
	o.Vni.Unset()
}

// GetVniAutoAssigned returns the VniAutoAssigned field value if set, zero value otherwise.
func (o *ConfigPutRequestServiceServiceName) GetVniAutoAssigned() bool {
	if o == nil || IsNil(o.VniAutoAssigned) {
		var ret bool
		return ret
	}
	return *o.VniAutoAssigned
}

// GetVniAutoAssignedOk returns a tuple with the VniAutoAssigned field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestServiceServiceName) GetVniAutoAssignedOk() (*bool, bool) {
	if o == nil || IsNil(o.VniAutoAssigned) {
		return nil, false
	}
	return o.VniAutoAssigned, true
}

// HasVniAutoAssigned returns a boolean if a field has been set.
func (o *ConfigPutRequestServiceServiceName) HasVniAutoAssigned() bool {
	if o != nil && !IsNil(o.VniAutoAssigned) {
		return true
	}

	return false
}

// SetVniAutoAssigned gets a reference to the given bool and assigns it to the VniAutoAssigned field.
func (o *ConfigPutRequestServiceServiceName) SetVniAutoAssigned(v bool) {
	o.VniAutoAssigned = &v
}

// GetTenant returns the Tenant field value if set, zero value otherwise.
func (o *ConfigPutRequestServiceServiceName) GetTenant() string {
	if o == nil || IsNil(o.Tenant) {
		var ret string
		return ret
	}
	return *o.Tenant
}

// GetTenantOk returns a tuple with the Tenant field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestServiceServiceName) GetTenantOk() (*string, bool) {
	if o == nil || IsNil(o.Tenant) {
		return nil, false
	}
	return o.Tenant, true
}

// HasTenant returns a boolean if a field has been set.
func (o *ConfigPutRequestServiceServiceName) HasTenant() bool {
	if o != nil && !IsNil(o.Tenant) {
		return true
	}

	return false
}

// SetTenant gets a reference to the given string and assigns it to the Tenant field.
func (o *ConfigPutRequestServiceServiceName) SetTenant(v string) {
	o.Tenant = &v
}

// GetTenantRefType returns the TenantRefType field value if set, zero value otherwise.
func (o *ConfigPutRequestServiceServiceName) GetTenantRefType() string {
	if o == nil || IsNil(o.TenantRefType) {
		var ret string
		return ret
	}
	return *o.TenantRefType
}

// GetTenantRefTypeOk returns a tuple with the TenantRefType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestServiceServiceName) GetTenantRefTypeOk() (*string, bool) {
	if o == nil || IsNil(o.TenantRefType) {
		return nil, false
	}
	return o.TenantRefType, true
}

// HasTenantRefType returns a boolean if a field has been set.
func (o *ConfigPutRequestServiceServiceName) HasTenantRefType() bool {
	if o != nil && !IsNil(o.TenantRefType) {
		return true
	}

	return false
}

// SetTenantRefType gets a reference to the given string and assigns it to the TenantRefType field.
func (o *ConfigPutRequestServiceServiceName) SetTenantRefType(v string) {
	o.TenantRefType = &v
}

// GetAnycastIpMask returns the AnycastIpMask field value if set, zero value otherwise.
func (o *ConfigPutRequestServiceServiceName) GetAnycastIpMask() string {
	if o == nil || IsNil(o.AnycastIpMask) {
		var ret string
		return ret
	}
	return *o.AnycastIpMask
}

// GetAnycastIpMaskOk returns a tuple with the AnycastIpMask field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestServiceServiceName) GetAnycastIpMaskOk() (*string, bool) {
	if o == nil || IsNil(o.AnycastIpMask) {
		return nil, false
	}
	return o.AnycastIpMask, true
}

// HasAnycastIpMask returns a boolean if a field has been set.
func (o *ConfigPutRequestServiceServiceName) HasAnycastIpMask() bool {
	if o != nil && !IsNil(o.AnycastIpMask) {
		return true
	}

	return false
}

// SetAnycastIpMask gets a reference to the given string and assigns it to the AnycastIpMask field.
func (o *ConfigPutRequestServiceServiceName) SetAnycastIpMask(v string) {
	o.AnycastIpMask = &v
}

// GetDhcpServerIp returns the DhcpServerIp field value if set, zero value otherwise.
func (o *ConfigPutRequestServiceServiceName) GetDhcpServerIp() string {
	if o == nil || IsNil(o.DhcpServerIp) {
		var ret string
		return ret
	}
	return *o.DhcpServerIp
}

// GetDhcpServerIpOk returns a tuple with the DhcpServerIp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestServiceServiceName) GetDhcpServerIpOk() (*string, bool) {
	if o == nil || IsNil(o.DhcpServerIp) {
		return nil, false
	}
	return o.DhcpServerIp, true
}

// HasDhcpServerIp returns a boolean if a field has been set.
func (o *ConfigPutRequestServiceServiceName) HasDhcpServerIp() bool {
	if o != nil && !IsNil(o.DhcpServerIp) {
		return true
	}

	return false
}

// SetDhcpServerIp gets a reference to the given string and assigns it to the DhcpServerIp field.
func (o *ConfigPutRequestServiceServiceName) SetDhcpServerIp(v string) {
	o.DhcpServerIp = &v
}

// GetMtu returns the Mtu field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *ConfigPutRequestServiceServiceName) GetMtu() int32 {
	if o == nil || IsNil(o.Mtu.Get()) {
		var ret int32
		return ret
	}
	return *o.Mtu.Get()
}

// GetMtuOk returns a tuple with the Mtu field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ConfigPutRequestServiceServiceName) GetMtuOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.Mtu.Get(), o.Mtu.IsSet()
}

// HasMtu returns a boolean if a field has been set.
func (o *ConfigPutRequestServiceServiceName) HasMtu() bool {
	if o != nil && o.Mtu.IsSet() {
		return true
	}

	return false
}

// SetMtu gets a reference to the given NullableInt32 and assigns it to the Mtu field.
func (o *ConfigPutRequestServiceServiceName) SetMtu(v int32) {
	o.Mtu.Set(&v)
}
// SetMtuNil sets the value for Mtu to be an explicit nil
func (o *ConfigPutRequestServiceServiceName) SetMtuNil() {
	o.Mtu.Set(nil)
}

// UnsetMtu ensures that no value is present for Mtu, not even an explicit nil
func (o *ConfigPutRequestServiceServiceName) UnsetMtu() {
	o.Mtu.Unset()
}

// GetObjectProperties returns the ObjectProperties field value if set, zero value otherwise.
func (o *ConfigPutRequestServiceServiceName) GetObjectProperties() ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties {
	if o == nil || IsNil(o.ObjectProperties) {
		var ret ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties
		return ret
	}
	return *o.ObjectProperties
}

// GetObjectPropertiesOk returns a tuple with the ObjectProperties field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestServiceServiceName) GetObjectPropertiesOk() (*ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties, bool) {
	if o == nil || IsNil(o.ObjectProperties) {
		return nil, false
	}
	return o.ObjectProperties, true
}

// HasObjectProperties returns a boolean if a field has been set.
func (o *ConfigPutRequestServiceServiceName) HasObjectProperties() bool {
	if o != nil && !IsNil(o.ObjectProperties) {
		return true
	}

	return false
}

// SetObjectProperties gets a reference to the given ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties and assigns it to the ObjectProperties field.
func (o *ConfigPutRequestServiceServiceName) SetObjectProperties(v ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties) {
	o.ObjectProperties = &v
}

func (o ConfigPutRequestServiceServiceName) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ConfigPutRequestServiceServiceName) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.Enable) {
		toSerialize["enable"] = o.Enable
	}
	if o.Vlan.IsSet() {
		toSerialize["vlan"] = o.Vlan.Get()
	}
	if o.Vni.IsSet() {
		toSerialize["vni"] = o.Vni.Get()
	}
	if !IsNil(o.VniAutoAssigned) {
		toSerialize["vni_auto_assigned_"] = o.VniAutoAssigned
	}
	if !IsNil(o.Tenant) {
		toSerialize["tenant"] = o.Tenant
	}
	if !IsNil(o.TenantRefType) {
		toSerialize["tenant_ref_type_"] = o.TenantRefType
	}
	if !IsNil(o.AnycastIpMask) {
		toSerialize["anycast_ip_mask"] = o.AnycastIpMask
	}
	if !IsNil(o.DhcpServerIp) {
		toSerialize["dhcp_server_ip"] = o.DhcpServerIp
	}
	if o.Mtu.IsSet() {
		toSerialize["mtu"] = o.Mtu.Get()
	}
	if !IsNil(o.ObjectProperties) {
		toSerialize["object_properties"] = o.ObjectProperties
	}
	return toSerialize, nil
}

type NullableConfigPutRequestServiceServiceName struct {
	value *ConfigPutRequestServiceServiceName
	isSet bool
}

func (v NullableConfigPutRequestServiceServiceName) Get() *ConfigPutRequestServiceServiceName {
	return v.value
}

func (v *NullableConfigPutRequestServiceServiceName) Set(val *ConfigPutRequestServiceServiceName) {
	v.value = val
	v.isSet = true
}

func (v NullableConfigPutRequestServiceServiceName) IsSet() bool {
	return v.isSet
}

func (v *NullableConfigPutRequestServiceServiceName) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConfigPutRequestServiceServiceName(val *ConfigPutRequestServiceServiceName) *NullableConfigPutRequestServiceServiceName {
	return &NullableConfigPutRequestServiceServiceName{value: val, isSet: true}
}

func (v NullableConfigPutRequestServiceServiceName) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConfigPutRequestServiceServiceName) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


