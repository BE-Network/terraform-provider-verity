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

// checks if the ConfigPutRequestIpv6ListFilterIpv6ListFilterName type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ConfigPutRequestIpv6ListFilterIpv6ListFilterName{}

// ConfigPutRequestIpv6ListFilterIpv6ListFilterName struct for ConfigPutRequestIpv6ListFilterIpv6ListFilterName
type ConfigPutRequestIpv6ListFilterIpv6ListFilterName struct {
	// Object Name. Must be unique.
	Name *string `json:"name,omitempty"`
	// Enable object.
	Enable *bool `json:"enable,omitempty"`
	// Comma separated list of IPv6 addresses
	Ipv6List *string `json:"ipv6_list,omitempty" validate:"regexp=^.*$"`
}

// NewConfigPutRequestIpv6ListFilterIpv6ListFilterName instantiates a new ConfigPutRequestIpv6ListFilterIpv6ListFilterName object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConfigPutRequestIpv6ListFilterIpv6ListFilterName() *ConfigPutRequestIpv6ListFilterIpv6ListFilterName {
	this := ConfigPutRequestIpv6ListFilterIpv6ListFilterName{}
	var name string = ""
	this.Name = &name
	var enable bool = false
	this.Enable = &enable
	var ipv6List string = ""
	this.Ipv6List = &ipv6List
	return &this
}

// NewConfigPutRequestIpv6ListFilterIpv6ListFilterNameWithDefaults instantiates a new ConfigPutRequestIpv6ListFilterIpv6ListFilterName object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConfigPutRequestIpv6ListFilterIpv6ListFilterNameWithDefaults() *ConfigPutRequestIpv6ListFilterIpv6ListFilterName {
	this := ConfigPutRequestIpv6ListFilterIpv6ListFilterName{}
	var name string = ""
	this.Name = &name
	var enable bool = false
	this.Enable = &enable
	var ipv6List string = ""
	this.Ipv6List = &ipv6List
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *ConfigPutRequestIpv6ListFilterIpv6ListFilterName) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestIpv6ListFilterIpv6ListFilterName) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *ConfigPutRequestIpv6ListFilterIpv6ListFilterName) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *ConfigPutRequestIpv6ListFilterIpv6ListFilterName) SetName(v string) {
	o.Name = &v
}

// GetEnable returns the Enable field value if set, zero value otherwise.
func (o *ConfigPutRequestIpv6ListFilterIpv6ListFilterName) GetEnable() bool {
	if o == nil || IsNil(o.Enable) {
		var ret bool
		return ret
	}
	return *o.Enable
}

// GetEnableOk returns a tuple with the Enable field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestIpv6ListFilterIpv6ListFilterName) GetEnableOk() (*bool, bool) {
	if o == nil || IsNil(o.Enable) {
		return nil, false
	}
	return o.Enable, true
}

// HasEnable returns a boolean if a field has been set.
func (o *ConfigPutRequestIpv6ListFilterIpv6ListFilterName) HasEnable() bool {
	if o != nil && !IsNil(o.Enable) {
		return true
	}

	return false
}

// SetEnable gets a reference to the given bool and assigns it to the Enable field.
func (o *ConfigPutRequestIpv6ListFilterIpv6ListFilterName) SetEnable(v bool) {
	o.Enable = &v
}

// GetIpv6List returns the Ipv6List field value if set, zero value otherwise.
func (o *ConfigPutRequestIpv6ListFilterIpv6ListFilterName) GetIpv6List() string {
	if o == nil || IsNil(o.Ipv6List) {
		var ret string
		return ret
	}
	return *o.Ipv6List
}

// GetIpv6ListOk returns a tuple with the Ipv6List field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestIpv6ListFilterIpv6ListFilterName) GetIpv6ListOk() (*string, bool) {
	if o == nil || IsNil(o.Ipv6List) {
		return nil, false
	}
	return o.Ipv6List, true
}

// HasIpv6List returns a boolean if a field has been set.
func (o *ConfigPutRequestIpv6ListFilterIpv6ListFilterName) HasIpv6List() bool {
	if o != nil && !IsNil(o.Ipv6List) {
		return true
	}

	return false
}

// SetIpv6List gets a reference to the given string and assigns it to the Ipv6List field.
func (o *ConfigPutRequestIpv6ListFilterIpv6ListFilterName) SetIpv6List(v string) {
	o.Ipv6List = &v
}

func (o ConfigPutRequestIpv6ListFilterIpv6ListFilterName) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ConfigPutRequestIpv6ListFilterIpv6ListFilterName) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.Enable) {
		toSerialize["enable"] = o.Enable
	}
	if !IsNil(o.Ipv6List) {
		toSerialize["ipv6_list"] = o.Ipv6List
	}
	return toSerialize, nil
}

type NullableConfigPutRequestIpv6ListFilterIpv6ListFilterName struct {
	value *ConfigPutRequestIpv6ListFilterIpv6ListFilterName
	isSet bool
}

func (v NullableConfigPutRequestIpv6ListFilterIpv6ListFilterName) Get() *ConfigPutRequestIpv6ListFilterIpv6ListFilterName {
	return v.value
}

func (v *NullableConfigPutRequestIpv6ListFilterIpv6ListFilterName) Set(val *ConfigPutRequestIpv6ListFilterIpv6ListFilterName) {
	v.value = val
	v.isSet = true
}

func (v NullableConfigPutRequestIpv6ListFilterIpv6ListFilterName) IsSet() bool {
	return v.isSet
}

func (v *NullableConfigPutRequestIpv6ListFilterIpv6ListFilterName) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConfigPutRequestIpv6ListFilterIpv6ListFilterName(val *ConfigPutRequestIpv6ListFilterIpv6ListFilterName) *NullableConfigPutRequestIpv6ListFilterIpv6ListFilterName {
	return &NullableConfigPutRequestIpv6ListFilterIpv6ListFilterName{value: val, isSet: true}
}

func (v NullableConfigPutRequestIpv6ListFilterIpv6ListFilterName) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConfigPutRequestIpv6ListFilterIpv6ListFilterName) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


