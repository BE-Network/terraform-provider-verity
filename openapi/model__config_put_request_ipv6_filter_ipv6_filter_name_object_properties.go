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

// checks if the ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties{}

// ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties struct for ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties
type ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties struct {
	// User Notes.
	Notes *string `json:"notes,omitempty"`
}

// NewConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties instantiates a new ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties() *ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties {
	this := ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties{}
	var notes string = ""
	this.Notes = &notes
	return &this
}

// NewConfigPutRequestIpv6FilterIpv6FilterNameObjectPropertiesWithDefaults instantiates a new ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConfigPutRequestIpv6FilterIpv6FilterNameObjectPropertiesWithDefaults() *ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties {
	this := ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties{}
	var notes string = ""
	this.Notes = &notes
	return &this
}

// GetNotes returns the Notes field value if set, zero value otherwise.
func (o *ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties) GetNotes() string {
	if o == nil || IsNil(o.Notes) {
		var ret string
		return ret
	}
	return *o.Notes
}

// GetNotesOk returns a tuple with the Notes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties) GetNotesOk() (*string, bool) {
	if o == nil || IsNil(o.Notes) {
		return nil, false
	}
	return o.Notes, true
}

// HasNotes returns a boolean if a field has been set.
func (o *ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties) HasNotes() bool {
	if o != nil && !IsNil(o.Notes) {
		return true
	}

	return false
}

// SetNotes gets a reference to the given string and assigns it to the Notes field.
func (o *ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties) SetNotes(v string) {
	o.Notes = &v
}

func (o ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Notes) {
		toSerialize["notes"] = o.Notes
	}
	return toSerialize, nil
}

type NullableConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties struct {
	value *ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties
	isSet bool
}

func (v NullableConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties) Get() *ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties {
	return v.value
}

func (v *NullableConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties) Set(val *ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties) {
	v.value = val
	v.isSet = true
}

func (v NullableConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties) IsSet() bool {
	return v.isSet
}

func (v *NullableConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties(val *ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties) *NullableConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties {
	return &NullableConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties{value: val, isSet: true}
}

func (v NullableConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


