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

// checks if the ConfigPutRequestIpPrefixListIpPrefixListName type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ConfigPutRequestIpPrefixListIpPrefixListName{}

// ConfigPutRequestIpPrefixListIpPrefixListName struct for ConfigPutRequestIpPrefixListIpPrefixListName
type ConfigPutRequestIpPrefixListIpPrefixListName struct {
	// Object Name. Must be unique.
	Name *string `json:"name,omitempty"`
	// Enable object.
	Enable *bool `json:"enable,omitempty"`
	Lists []ConfigPutRequestIpPrefixListIpPrefixListNameListsInner `json:"lists,omitempty"`
	ObjectProperties *ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties `json:"object_properties,omitempty"`
}

// NewConfigPutRequestIpPrefixListIpPrefixListName instantiates a new ConfigPutRequestIpPrefixListIpPrefixListName object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConfigPutRequestIpPrefixListIpPrefixListName() *ConfigPutRequestIpPrefixListIpPrefixListName {
	this := ConfigPutRequestIpPrefixListIpPrefixListName{}
	var name string = ""
	this.Name = &name
	var enable bool = false
	this.Enable = &enable
	return &this
}

// NewConfigPutRequestIpPrefixListIpPrefixListNameWithDefaults instantiates a new ConfigPutRequestIpPrefixListIpPrefixListName object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConfigPutRequestIpPrefixListIpPrefixListNameWithDefaults() *ConfigPutRequestIpPrefixListIpPrefixListName {
	this := ConfigPutRequestIpPrefixListIpPrefixListName{}
	var name string = ""
	this.Name = &name
	var enable bool = false
	this.Enable = &enable
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *ConfigPutRequestIpPrefixListIpPrefixListName) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestIpPrefixListIpPrefixListName) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *ConfigPutRequestIpPrefixListIpPrefixListName) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *ConfigPutRequestIpPrefixListIpPrefixListName) SetName(v string) {
	o.Name = &v
}

// GetEnable returns the Enable field value if set, zero value otherwise.
func (o *ConfigPutRequestIpPrefixListIpPrefixListName) GetEnable() bool {
	if o == nil || IsNil(o.Enable) {
		var ret bool
		return ret
	}
	return *o.Enable
}

// GetEnableOk returns a tuple with the Enable field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestIpPrefixListIpPrefixListName) GetEnableOk() (*bool, bool) {
	if o == nil || IsNil(o.Enable) {
		return nil, false
	}
	return o.Enable, true
}

// HasEnable returns a boolean if a field has been set.
func (o *ConfigPutRequestIpPrefixListIpPrefixListName) HasEnable() bool {
	if o != nil && !IsNil(o.Enable) {
		return true
	}

	return false
}

// SetEnable gets a reference to the given bool and assigns it to the Enable field.
func (o *ConfigPutRequestIpPrefixListIpPrefixListName) SetEnable(v bool) {
	o.Enable = &v
}

// GetLists returns the Lists field value if set, zero value otherwise.
func (o *ConfigPutRequestIpPrefixListIpPrefixListName) GetLists() []ConfigPutRequestIpPrefixListIpPrefixListNameListsInner {
	if o == nil || IsNil(o.Lists) {
		var ret []ConfigPutRequestIpPrefixListIpPrefixListNameListsInner
		return ret
	}
	return o.Lists
}

// GetListsOk returns a tuple with the Lists field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestIpPrefixListIpPrefixListName) GetListsOk() ([]ConfigPutRequestIpPrefixListIpPrefixListNameListsInner, bool) {
	if o == nil || IsNil(o.Lists) {
		return nil, false
	}
	return o.Lists, true
}

// HasLists returns a boolean if a field has been set.
func (o *ConfigPutRequestIpPrefixListIpPrefixListName) HasLists() bool {
	if o != nil && !IsNil(o.Lists) {
		return true
	}

	return false
}

// SetLists gets a reference to the given []ConfigPutRequestIpPrefixListIpPrefixListNameListsInner and assigns it to the Lists field.
func (o *ConfigPutRequestIpPrefixListIpPrefixListName) SetLists(v []ConfigPutRequestIpPrefixListIpPrefixListNameListsInner) {
	o.Lists = v
}

// GetObjectProperties returns the ObjectProperties field value if set, zero value otherwise.
func (o *ConfigPutRequestIpPrefixListIpPrefixListName) GetObjectProperties() ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties {
	if o == nil || IsNil(o.ObjectProperties) {
		var ret ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties
		return ret
	}
	return *o.ObjectProperties
}

// GetObjectPropertiesOk returns a tuple with the ObjectProperties field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestIpPrefixListIpPrefixListName) GetObjectPropertiesOk() (*ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties, bool) {
	if o == nil || IsNil(o.ObjectProperties) {
		return nil, false
	}
	return o.ObjectProperties, true
}

// HasObjectProperties returns a boolean if a field has been set.
func (o *ConfigPutRequestIpPrefixListIpPrefixListName) HasObjectProperties() bool {
	if o != nil && !IsNil(o.ObjectProperties) {
		return true
	}

	return false
}

// SetObjectProperties gets a reference to the given ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties and assigns it to the ObjectProperties field.
func (o *ConfigPutRequestIpPrefixListIpPrefixListName) SetObjectProperties(v ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties) {
	o.ObjectProperties = &v
}

func (o ConfigPutRequestIpPrefixListIpPrefixListName) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ConfigPutRequestIpPrefixListIpPrefixListName) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.Enable) {
		toSerialize["enable"] = o.Enable
	}
	if !IsNil(o.Lists) {
		toSerialize["lists"] = o.Lists
	}
	if !IsNil(o.ObjectProperties) {
		toSerialize["object_properties"] = o.ObjectProperties
	}
	return toSerialize, nil
}

type NullableConfigPutRequestIpPrefixListIpPrefixListName struct {
	value *ConfigPutRequestIpPrefixListIpPrefixListName
	isSet bool
}

func (v NullableConfigPutRequestIpPrefixListIpPrefixListName) Get() *ConfigPutRequestIpPrefixListIpPrefixListName {
	return v.value
}

func (v *NullableConfigPutRequestIpPrefixListIpPrefixListName) Set(val *ConfigPutRequestIpPrefixListIpPrefixListName) {
	v.value = val
	v.isSet = true
}

func (v NullableConfigPutRequestIpPrefixListIpPrefixListName) IsSet() bool {
	return v.isSet
}

func (v *NullableConfigPutRequestIpPrefixListIpPrefixListName) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConfigPutRequestIpPrefixListIpPrefixListName(val *ConfigPutRequestIpPrefixListIpPrefixListName) *NullableConfigPutRequestIpPrefixListIpPrefixListName {
	return &NullableConfigPutRequestIpPrefixListIpPrefixListName{value: val, isSet: true}
}

func (v NullableConfigPutRequestIpPrefixListIpPrefixListName) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConfigPutRequestIpPrefixListIpPrefixListName) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


