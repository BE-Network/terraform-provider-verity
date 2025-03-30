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

// checks if the ConfigPutRequestLag type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ConfigPutRequestLag{}

// ConfigPutRequestLag struct for ConfigPutRequestLag
type ConfigPutRequestLag struct {
	LagName *ConfigPutRequestLagLagName `json:"lag_name,omitempty"`
}

// NewConfigPutRequestLag instantiates a new ConfigPutRequestLag object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConfigPutRequestLag() *ConfigPutRequestLag {
	this := ConfigPutRequestLag{}
	return &this
}

// NewConfigPutRequestLagWithDefaults instantiates a new ConfigPutRequestLag object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConfigPutRequestLagWithDefaults() *ConfigPutRequestLag {
	this := ConfigPutRequestLag{}
	return &this
}

// GetLagName returns the LagName field value if set, zero value otherwise.
func (o *ConfigPutRequestLag) GetLagName() ConfigPutRequestLagLagName {
	if o == nil || IsNil(o.LagName) {
		var ret ConfigPutRequestLagLagName
		return ret
	}
	return *o.LagName
}

// GetLagNameOk returns a tuple with the LagName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestLag) GetLagNameOk() (*ConfigPutRequestLagLagName, bool) {
	if o == nil || IsNil(o.LagName) {
		return nil, false
	}
	return o.LagName, true
}

// HasLagName returns a boolean if a field has been set.
func (o *ConfigPutRequestLag) HasLagName() bool {
	if o != nil && !IsNil(o.LagName) {
		return true
	}

	return false
}

// SetLagName gets a reference to the given ConfigPutRequestLagLagName and assigns it to the LagName field.
func (o *ConfigPutRequestLag) SetLagName(v ConfigPutRequestLagLagName) {
	o.LagName = &v
}

func (o ConfigPutRequestLag) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ConfigPutRequestLag) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.LagName) {
		toSerialize["lag_name"] = o.LagName
	}
	return toSerialize, nil
}

type NullableConfigPutRequestLag struct {
	value *ConfigPutRequestLag
	isSet bool
}

func (v NullableConfigPutRequestLag) Get() *ConfigPutRequestLag {
	return v.value
}

func (v *NullableConfigPutRequestLag) Set(val *ConfigPutRequestLag) {
	v.value = val
	v.isSet = true
}

func (v NullableConfigPutRequestLag) IsSet() bool {
	return v.isSet
}

func (v *NullableConfigPutRequestLag) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConfigPutRequestLag(val *ConfigPutRequestLag) *NullableConfigPutRequestLag {
	return &NullableConfigPutRequestLag{value: val, isSet: true}
}

func (v NullableConfigPutRequestLag) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConfigPutRequestLag) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


