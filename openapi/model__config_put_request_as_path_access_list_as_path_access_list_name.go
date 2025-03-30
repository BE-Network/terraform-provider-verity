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

// checks if the ConfigPutRequestAsPathAccessListAsPathAccessListName type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ConfigPutRequestAsPathAccessListAsPathAccessListName{}

// ConfigPutRequestAsPathAccessListAsPathAccessListName struct for ConfigPutRequestAsPathAccessListAsPathAccessListName
type ConfigPutRequestAsPathAccessListAsPathAccessListName struct {
	// Object Name. Must be unique.
	Name *string `json:"name,omitempty"`
	// Enable object.
	Enable *bool `json:"enable,omitempty"`
	// Action upon match of Community Strings.
	PermitDeny *string `json:"permit_deny,omitempty"`
	Lists []ConfigPutRequestAsPathAccessListAsPathAccessListNameListsInner `json:"lists,omitempty"`
	ObjectProperties *ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties `json:"object_properties,omitempty"`
}

// NewConfigPutRequestAsPathAccessListAsPathAccessListName instantiates a new ConfigPutRequestAsPathAccessListAsPathAccessListName object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConfigPutRequestAsPathAccessListAsPathAccessListName() *ConfigPutRequestAsPathAccessListAsPathAccessListName {
	this := ConfigPutRequestAsPathAccessListAsPathAccessListName{}
	var name string = ""
	this.Name = &name
	var enable bool = false
	this.Enable = &enable
	var permitDeny string = "permit"
	this.PermitDeny = &permitDeny
	return &this
}

// NewConfigPutRequestAsPathAccessListAsPathAccessListNameWithDefaults instantiates a new ConfigPutRequestAsPathAccessListAsPathAccessListName object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConfigPutRequestAsPathAccessListAsPathAccessListNameWithDefaults() *ConfigPutRequestAsPathAccessListAsPathAccessListName {
	this := ConfigPutRequestAsPathAccessListAsPathAccessListName{}
	var name string = ""
	this.Name = &name
	var enable bool = false
	this.Enable = &enable
	var permitDeny string = "permit"
	this.PermitDeny = &permitDeny
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) SetName(v string) {
	o.Name = &v
}

// GetEnable returns the Enable field value if set, zero value otherwise.
func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) GetEnable() bool {
	if o == nil || IsNil(o.Enable) {
		var ret bool
		return ret
	}
	return *o.Enable
}

// GetEnableOk returns a tuple with the Enable field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) GetEnableOk() (*bool, bool) {
	if o == nil || IsNil(o.Enable) {
		return nil, false
	}
	return o.Enable, true
}

// HasEnable returns a boolean if a field has been set.
func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) HasEnable() bool {
	if o != nil && !IsNil(o.Enable) {
		return true
	}

	return false
}

// SetEnable gets a reference to the given bool and assigns it to the Enable field.
func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) SetEnable(v bool) {
	o.Enable = &v
}

// GetPermitDeny returns the PermitDeny field value if set, zero value otherwise.
func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) GetPermitDeny() string {
	if o == nil || IsNil(o.PermitDeny) {
		var ret string
		return ret
	}
	return *o.PermitDeny
}

// GetPermitDenyOk returns a tuple with the PermitDeny field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) GetPermitDenyOk() (*string, bool) {
	if o == nil || IsNil(o.PermitDeny) {
		return nil, false
	}
	return o.PermitDeny, true
}

// HasPermitDeny returns a boolean if a field has been set.
func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) HasPermitDeny() bool {
	if o != nil && !IsNil(o.PermitDeny) {
		return true
	}

	return false
}

// SetPermitDeny gets a reference to the given string and assigns it to the PermitDeny field.
func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) SetPermitDeny(v string) {
	o.PermitDeny = &v
}

// GetLists returns the Lists field value if set, zero value otherwise.
func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) GetLists() []ConfigPutRequestAsPathAccessListAsPathAccessListNameListsInner {
	if o == nil || IsNil(o.Lists) {
		var ret []ConfigPutRequestAsPathAccessListAsPathAccessListNameListsInner
		return ret
	}
	return o.Lists
}

// GetListsOk returns a tuple with the Lists field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) GetListsOk() ([]ConfigPutRequestAsPathAccessListAsPathAccessListNameListsInner, bool) {
	if o == nil || IsNil(o.Lists) {
		return nil, false
	}
	return o.Lists, true
}

// HasLists returns a boolean if a field has been set.
func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) HasLists() bool {
	if o != nil && !IsNil(o.Lists) {
		return true
	}

	return false
}

// SetLists gets a reference to the given []ConfigPutRequestAsPathAccessListAsPathAccessListNameListsInner and assigns it to the Lists field.
func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) SetLists(v []ConfigPutRequestAsPathAccessListAsPathAccessListNameListsInner) {
	o.Lists = v
}

// GetObjectProperties returns the ObjectProperties field value if set, zero value otherwise.
func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) GetObjectProperties() ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties {
	if o == nil || IsNil(o.ObjectProperties) {
		var ret ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties
		return ret
	}
	return *o.ObjectProperties
}

// GetObjectPropertiesOk returns a tuple with the ObjectProperties field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) GetObjectPropertiesOk() (*ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties, bool) {
	if o == nil || IsNil(o.ObjectProperties) {
		return nil, false
	}
	return o.ObjectProperties, true
}

// HasObjectProperties returns a boolean if a field has been set.
func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) HasObjectProperties() bool {
	if o != nil && !IsNil(o.ObjectProperties) {
		return true
	}

	return false
}

// SetObjectProperties gets a reference to the given ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties and assigns it to the ObjectProperties field.
func (o *ConfigPutRequestAsPathAccessListAsPathAccessListName) SetObjectProperties(v ConfigPutRequestIpv6FilterIpv6FilterNameObjectProperties) {
	o.ObjectProperties = &v
}

func (o ConfigPutRequestAsPathAccessListAsPathAccessListName) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ConfigPutRequestAsPathAccessListAsPathAccessListName) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.Enable) {
		toSerialize["enable"] = o.Enable
	}
	if !IsNil(o.PermitDeny) {
		toSerialize["permit_deny"] = o.PermitDeny
	}
	if !IsNil(o.Lists) {
		toSerialize["lists"] = o.Lists
	}
	if !IsNil(o.ObjectProperties) {
		toSerialize["object_properties"] = o.ObjectProperties
	}
	return toSerialize, nil
}

type NullableConfigPutRequestAsPathAccessListAsPathAccessListName struct {
	value *ConfigPutRequestAsPathAccessListAsPathAccessListName
	isSet bool
}

func (v NullableConfigPutRequestAsPathAccessListAsPathAccessListName) Get() *ConfigPutRequestAsPathAccessListAsPathAccessListName {
	return v.value
}

func (v *NullableConfigPutRequestAsPathAccessListAsPathAccessListName) Set(val *ConfigPutRequestAsPathAccessListAsPathAccessListName) {
	v.value = val
	v.isSet = true
}

func (v NullableConfigPutRequestAsPathAccessListAsPathAccessListName) IsSet() bool {
	return v.isSet
}

func (v *NullableConfigPutRequestAsPathAccessListAsPathAccessListName) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConfigPutRequestAsPathAccessListAsPathAccessListName(val *ConfigPutRequestAsPathAccessListAsPathAccessListName) *NullableConfigPutRequestAsPathAccessListAsPathAccessListName {
	return &NullableConfigPutRequestAsPathAccessListAsPathAccessListName{value: val, isSet: true}
}

func (v NullableConfigPutRequestAsPathAccessListAsPathAccessListName) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConfigPutRequestAsPathAccessListAsPathAccessListName) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


