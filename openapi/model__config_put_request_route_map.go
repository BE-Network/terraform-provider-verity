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

// checks if the ConfigPutRequestRouteMap type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ConfigPutRequestRouteMap{}

// ConfigPutRequestRouteMap struct for ConfigPutRequestRouteMap
type ConfigPutRequestRouteMap struct {
	RouteMapName *ConfigPutRequestRouteMapRouteMapName `json:"route_map_name,omitempty"`
}

// NewConfigPutRequestRouteMap instantiates a new ConfigPutRequestRouteMap object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConfigPutRequestRouteMap() *ConfigPutRequestRouteMap {
	this := ConfigPutRequestRouteMap{}
	return &this
}

// NewConfigPutRequestRouteMapWithDefaults instantiates a new ConfigPutRequestRouteMap object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConfigPutRequestRouteMapWithDefaults() *ConfigPutRequestRouteMap {
	this := ConfigPutRequestRouteMap{}
	return &this
}

// GetRouteMapName returns the RouteMapName field value if set, zero value otherwise.
func (o *ConfigPutRequestRouteMap) GetRouteMapName() ConfigPutRequestRouteMapRouteMapName {
	if o == nil || IsNil(o.RouteMapName) {
		var ret ConfigPutRequestRouteMapRouteMapName
		return ret
	}
	return *o.RouteMapName
}

// GetRouteMapNameOk returns a tuple with the RouteMapName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestRouteMap) GetRouteMapNameOk() (*ConfigPutRequestRouteMapRouteMapName, bool) {
	if o == nil || IsNil(o.RouteMapName) {
		return nil, false
	}
	return o.RouteMapName, true
}

// HasRouteMapName returns a boolean if a field has been set.
func (o *ConfigPutRequestRouteMap) HasRouteMapName() bool {
	if o != nil && !IsNil(o.RouteMapName) {
		return true
	}

	return false
}

// SetRouteMapName gets a reference to the given ConfigPutRequestRouteMapRouteMapName and assigns it to the RouteMapName field.
func (o *ConfigPutRequestRouteMap) SetRouteMapName(v ConfigPutRequestRouteMapRouteMapName) {
	o.RouteMapName = &v
}

func (o ConfigPutRequestRouteMap) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ConfigPutRequestRouteMap) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.RouteMapName) {
		toSerialize["route_map_name"] = o.RouteMapName
	}
	return toSerialize, nil
}

type NullableConfigPutRequestRouteMap struct {
	value *ConfigPutRequestRouteMap
	isSet bool
}

func (v NullableConfigPutRequestRouteMap) Get() *ConfigPutRequestRouteMap {
	return v.value
}

func (v *NullableConfigPutRequestRouteMap) Set(val *ConfigPutRequestRouteMap) {
	v.value = val
	v.isSet = true
}

func (v NullableConfigPutRequestRouteMap) IsSet() bool {
	return v.isSet
}

func (v *NullableConfigPutRequestRouteMap) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConfigPutRequestRouteMap(val *ConfigPutRequestRouteMap) *NullableConfigPutRequestRouteMap {
	return &NullableConfigPutRequestRouteMap{value: val, isSet: true}
}

func (v NullableConfigPutRequestRouteMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConfigPutRequestRouteMap) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


