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

// checks if the ConfigPutRequestEthPortProfile type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ConfigPutRequestEthPortProfile{}

// ConfigPutRequestEthPortProfile struct for ConfigPutRequestEthPortProfile
type ConfigPutRequestEthPortProfile struct {
	EthPortProfileName *ConfigPutRequestEthPortProfileEthPortProfileName `json:"eth_port_profile__name,omitempty"`
}

// NewConfigPutRequestEthPortProfile instantiates a new ConfigPutRequestEthPortProfile object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConfigPutRequestEthPortProfile() *ConfigPutRequestEthPortProfile {
	this := ConfigPutRequestEthPortProfile{}
	return &this
}

// NewConfigPutRequestEthPortProfileWithDefaults instantiates a new ConfigPutRequestEthPortProfile object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConfigPutRequestEthPortProfileWithDefaults() *ConfigPutRequestEthPortProfile {
	this := ConfigPutRequestEthPortProfile{}
	return &this
}

// GetEthPortProfileName returns the EthPortProfileName field value if set, zero value otherwise.
func (o *ConfigPutRequestEthPortProfile) GetEthPortProfileName() ConfigPutRequestEthPortProfileEthPortProfileName {
	if o == nil || IsNil(o.EthPortProfileName) {
		var ret ConfigPutRequestEthPortProfileEthPortProfileName
		return ret
	}
	return *o.EthPortProfileName
}

// GetEthPortProfileNameOk returns a tuple with the EthPortProfileName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEthPortProfile) GetEthPortProfileNameOk() (*ConfigPutRequestEthPortProfileEthPortProfileName, bool) {
	if o == nil || IsNil(o.EthPortProfileName) {
		return nil, false
	}
	return o.EthPortProfileName, true
}

// HasEthPortProfileName returns a boolean if a field has been set.
func (o *ConfigPutRequestEthPortProfile) HasEthPortProfileName() bool {
	if o != nil && !IsNil(o.EthPortProfileName) {
		return true
	}

	return false
}

// SetEthPortProfileName gets a reference to the given ConfigPutRequestEthPortProfileEthPortProfileName and assigns it to the EthPortProfileName field.
func (o *ConfigPutRequestEthPortProfile) SetEthPortProfileName(v ConfigPutRequestEthPortProfileEthPortProfileName) {
	o.EthPortProfileName = &v
}

func (o ConfigPutRequestEthPortProfile) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ConfigPutRequestEthPortProfile) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.EthPortProfileName) {
		toSerialize["eth_port_profile__name"] = o.EthPortProfileName
	}
	return toSerialize, nil
}

type NullableConfigPutRequestEthPortProfile struct {
	value *ConfigPutRequestEthPortProfile
	isSet bool
}

func (v NullableConfigPutRequestEthPortProfile) Get() *ConfigPutRequestEthPortProfile {
	return v.value
}

func (v *NullableConfigPutRequestEthPortProfile) Set(val *ConfigPutRequestEthPortProfile) {
	v.value = val
	v.isSet = true
}

func (v NullableConfigPutRequestEthPortProfile) IsSet() bool {
	return v.isSet
}

func (v *NullableConfigPutRequestEthPortProfile) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConfigPutRequestEthPortProfile(val *ConfigPutRequestEthPortProfile) *NullableConfigPutRequestEthPortProfile {
	return &NullableConfigPutRequestEthPortProfile{value: val, isSet: true}
}

func (v NullableConfigPutRequestEthPortProfile) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConfigPutRequestEthPortProfile) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


