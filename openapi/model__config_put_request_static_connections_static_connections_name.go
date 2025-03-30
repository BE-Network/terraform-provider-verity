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

// checks if the ConfigPutRequestStaticConnectionsStaticConnectionsName type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ConfigPutRequestStaticConnectionsStaticConnectionsName{}

// ConfigPutRequestStaticConnectionsStaticConnectionsName struct for ConfigPutRequestStaticConnectionsStaticConnectionsName
type ConfigPutRequestStaticConnectionsStaticConnectionsName struct {
	// Enable object.
	Enable *bool `json:"enable,omitempty"`
	Connections []ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner `json:"connections,omitempty"`
	ObjectProperties map[string]interface{} `json:"object_properties,omitempty"`
}

// NewConfigPutRequestStaticConnectionsStaticConnectionsName instantiates a new ConfigPutRequestStaticConnectionsStaticConnectionsName object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConfigPutRequestStaticConnectionsStaticConnectionsName() *ConfigPutRequestStaticConnectionsStaticConnectionsName {
	this := ConfigPutRequestStaticConnectionsStaticConnectionsName{}
	var enable bool = true
	this.Enable = &enable
	return &this
}

// NewConfigPutRequestStaticConnectionsStaticConnectionsNameWithDefaults instantiates a new ConfigPutRequestStaticConnectionsStaticConnectionsName object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConfigPutRequestStaticConnectionsStaticConnectionsNameWithDefaults() *ConfigPutRequestStaticConnectionsStaticConnectionsName {
	this := ConfigPutRequestStaticConnectionsStaticConnectionsName{}
	var enable bool = true
	this.Enable = &enable
	return &this
}

// GetEnable returns the Enable field value if set, zero value otherwise.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsName) GetEnable() bool {
	if o == nil || IsNil(o.Enable) {
		var ret bool
		return ret
	}
	return *o.Enable
}

// GetEnableOk returns a tuple with the Enable field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsName) GetEnableOk() (*bool, bool) {
	if o == nil || IsNil(o.Enable) {
		return nil, false
	}
	return o.Enable, true
}

// HasEnable returns a boolean if a field has been set.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsName) HasEnable() bool {
	if o != nil && !IsNil(o.Enable) {
		return true
	}

	return false
}

// SetEnable gets a reference to the given bool and assigns it to the Enable field.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsName) SetEnable(v bool) {
	o.Enable = &v
}

// GetConnections returns the Connections field value if set, zero value otherwise.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsName) GetConnections() []ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner {
	if o == nil || IsNil(o.Connections) {
		var ret []ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner
		return ret
	}
	return o.Connections
}

// GetConnectionsOk returns a tuple with the Connections field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsName) GetConnectionsOk() ([]ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner, bool) {
	if o == nil || IsNil(o.Connections) {
		return nil, false
	}
	return o.Connections, true
}

// HasConnections returns a boolean if a field has been set.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsName) HasConnections() bool {
	if o != nil && !IsNil(o.Connections) {
		return true
	}

	return false
}

// SetConnections gets a reference to the given []ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner and assigns it to the Connections field.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsName) SetConnections(v []ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) {
	o.Connections = v
}

// GetObjectProperties returns the ObjectProperties field value if set, zero value otherwise.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsName) GetObjectProperties() map[string]interface{} {
	if o == nil || IsNil(o.ObjectProperties) {
		var ret map[string]interface{}
		return ret
	}
	return o.ObjectProperties
}

// GetObjectPropertiesOk returns a tuple with the ObjectProperties field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsName) GetObjectPropertiesOk() (map[string]interface{}, bool) {
	if o == nil || IsNil(o.ObjectProperties) {
		return map[string]interface{}{}, false
	}
	return o.ObjectProperties, true
}

// HasObjectProperties returns a boolean if a field has been set.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsName) HasObjectProperties() bool {
	if o != nil && !IsNil(o.ObjectProperties) {
		return true
	}

	return false
}

// SetObjectProperties gets a reference to the given map[string]interface{} and assigns it to the ObjectProperties field.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsName) SetObjectProperties(v map[string]interface{}) {
	o.ObjectProperties = v
}

func (o ConfigPutRequestStaticConnectionsStaticConnectionsName) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ConfigPutRequestStaticConnectionsStaticConnectionsName) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Enable) {
		toSerialize["enable"] = o.Enable
	}
	if !IsNil(o.Connections) {
		toSerialize["connections"] = o.Connections
	}
	if !IsNil(o.ObjectProperties) {
		toSerialize["object_properties"] = o.ObjectProperties
	}
	return toSerialize, nil
}

type NullableConfigPutRequestStaticConnectionsStaticConnectionsName struct {
	value *ConfigPutRequestStaticConnectionsStaticConnectionsName
	isSet bool
}

func (v NullableConfigPutRequestStaticConnectionsStaticConnectionsName) Get() *ConfigPutRequestStaticConnectionsStaticConnectionsName {
	return v.value
}

func (v *NullableConfigPutRequestStaticConnectionsStaticConnectionsName) Set(val *ConfigPutRequestStaticConnectionsStaticConnectionsName) {
	v.value = val
	v.isSet = true
}

func (v NullableConfigPutRequestStaticConnectionsStaticConnectionsName) IsSet() bool {
	return v.isSet
}

func (v *NullableConfigPutRequestStaticConnectionsStaticConnectionsName) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConfigPutRequestStaticConnectionsStaticConnectionsName(val *ConfigPutRequestStaticConnectionsStaticConnectionsName) *NullableConfigPutRequestStaticConnectionsStaticConnectionsName {
	return &NullableConfigPutRequestStaticConnectionsStaticConnectionsName{value: val, isSet: true}
}

func (v NullableConfigPutRequestStaticConnectionsStaticConnectionsName) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConfigPutRequestStaticConnectionsStaticConnectionsName) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


