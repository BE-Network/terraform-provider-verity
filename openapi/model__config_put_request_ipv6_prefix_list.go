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

// checks if the ConfigPutRequestIpv6PrefixList type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ConfigPutRequestIpv6PrefixList{}

// ConfigPutRequestIpv6PrefixList struct for ConfigPutRequestIpv6PrefixList
type ConfigPutRequestIpv6PrefixList struct {
	Ipv6PrefixListName *ConfigPutRequestIpv6PrefixListIpv6PrefixListName `json:"ipv6_prefix_list_name,omitempty"`
}

// NewConfigPutRequestIpv6PrefixList instantiates a new ConfigPutRequestIpv6PrefixList object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConfigPutRequestIpv6PrefixList() *ConfigPutRequestIpv6PrefixList {
	this := ConfigPutRequestIpv6PrefixList{}
	return &this
}

// NewConfigPutRequestIpv6PrefixListWithDefaults instantiates a new ConfigPutRequestIpv6PrefixList object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConfigPutRequestIpv6PrefixListWithDefaults() *ConfigPutRequestIpv6PrefixList {
	this := ConfigPutRequestIpv6PrefixList{}
	return &this
}

// GetIpv6PrefixListName returns the Ipv6PrefixListName field value if set, zero value otherwise.
func (o *ConfigPutRequestIpv6PrefixList) GetIpv6PrefixListName() ConfigPutRequestIpv6PrefixListIpv6PrefixListName {
	if o == nil || IsNil(o.Ipv6PrefixListName) {
		var ret ConfigPutRequestIpv6PrefixListIpv6PrefixListName
		return ret
	}
	return *o.Ipv6PrefixListName
}

// GetIpv6PrefixListNameOk returns a tuple with the Ipv6PrefixListName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestIpv6PrefixList) GetIpv6PrefixListNameOk() (*ConfigPutRequestIpv6PrefixListIpv6PrefixListName, bool) {
	if o == nil || IsNil(o.Ipv6PrefixListName) {
		return nil, false
	}
	return o.Ipv6PrefixListName, true
}

// HasIpv6PrefixListName returns a boolean if a field has been set.
func (o *ConfigPutRequestIpv6PrefixList) HasIpv6PrefixListName() bool {
	if o != nil && !IsNil(o.Ipv6PrefixListName) {
		return true
	}

	return false
}

// SetIpv6PrefixListName gets a reference to the given ConfigPutRequestIpv6PrefixListIpv6PrefixListName and assigns it to the Ipv6PrefixListName field.
func (o *ConfigPutRequestIpv6PrefixList) SetIpv6PrefixListName(v ConfigPutRequestIpv6PrefixListIpv6PrefixListName) {
	o.Ipv6PrefixListName = &v
}

func (o ConfigPutRequestIpv6PrefixList) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ConfigPutRequestIpv6PrefixList) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Ipv6PrefixListName) {
		toSerialize["ipv6_prefix_list_name"] = o.Ipv6PrefixListName
	}
	return toSerialize, nil
}

type NullableConfigPutRequestIpv6PrefixList struct {
	value *ConfigPutRequestIpv6PrefixList
	isSet bool
}

func (v NullableConfigPutRequestIpv6PrefixList) Get() *ConfigPutRequestIpv6PrefixList {
	return v.value
}

func (v *NullableConfigPutRequestIpv6PrefixList) Set(val *ConfigPutRequestIpv6PrefixList) {
	v.value = val
	v.isSet = true
}

func (v NullableConfigPutRequestIpv6PrefixList) IsSet() bool {
	return v.isSet
}

func (v *NullableConfigPutRequestIpv6PrefixList) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConfigPutRequestIpv6PrefixList(val *ConfigPutRequestIpv6PrefixList) *NullableConfigPutRequestIpv6PrefixList {
	return &NullableConfigPutRequestIpv6PrefixList{value: val, isSet: true}
}

func (v NullableConfigPutRequestIpv6PrefixList) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConfigPutRequestIpv6PrefixList) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


