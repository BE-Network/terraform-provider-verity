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

// checks if the ConfigPutRequestExtendedCommunityList type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ConfigPutRequestExtendedCommunityList{}

// ConfigPutRequestExtendedCommunityList struct for ConfigPutRequestExtendedCommunityList
type ConfigPutRequestExtendedCommunityList struct {
	ExtendedCommunityListName *ConfigPutRequestExtendedCommunityListExtendedCommunityListName `json:"extended_community_list_name,omitempty"`
}

// NewConfigPutRequestExtendedCommunityList instantiates a new ConfigPutRequestExtendedCommunityList object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConfigPutRequestExtendedCommunityList() *ConfigPutRequestExtendedCommunityList {
	this := ConfigPutRequestExtendedCommunityList{}
	return &this
}

// NewConfigPutRequestExtendedCommunityListWithDefaults instantiates a new ConfigPutRequestExtendedCommunityList object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConfigPutRequestExtendedCommunityListWithDefaults() *ConfigPutRequestExtendedCommunityList {
	this := ConfigPutRequestExtendedCommunityList{}
	return &this
}

// GetExtendedCommunityListName returns the ExtendedCommunityListName field value if set, zero value otherwise.
func (o *ConfigPutRequestExtendedCommunityList) GetExtendedCommunityListName() ConfigPutRequestExtendedCommunityListExtendedCommunityListName {
	if o == nil || IsNil(o.ExtendedCommunityListName) {
		var ret ConfigPutRequestExtendedCommunityListExtendedCommunityListName
		return ret
	}
	return *o.ExtendedCommunityListName
}

// GetExtendedCommunityListNameOk returns a tuple with the ExtendedCommunityListName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestExtendedCommunityList) GetExtendedCommunityListNameOk() (*ConfigPutRequestExtendedCommunityListExtendedCommunityListName, bool) {
	if o == nil || IsNil(o.ExtendedCommunityListName) {
		return nil, false
	}
	return o.ExtendedCommunityListName, true
}

// HasExtendedCommunityListName returns a boolean if a field has been set.
func (o *ConfigPutRequestExtendedCommunityList) HasExtendedCommunityListName() bool {
	if o != nil && !IsNil(o.ExtendedCommunityListName) {
		return true
	}

	return false
}

// SetExtendedCommunityListName gets a reference to the given ConfigPutRequestExtendedCommunityListExtendedCommunityListName and assigns it to the ExtendedCommunityListName field.
func (o *ConfigPutRequestExtendedCommunityList) SetExtendedCommunityListName(v ConfigPutRequestExtendedCommunityListExtendedCommunityListName) {
	o.ExtendedCommunityListName = &v
}

func (o ConfigPutRequestExtendedCommunityList) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ConfigPutRequestExtendedCommunityList) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.ExtendedCommunityListName) {
		toSerialize["extended_community_list_name"] = o.ExtendedCommunityListName
	}
	return toSerialize, nil
}

type NullableConfigPutRequestExtendedCommunityList struct {
	value *ConfigPutRequestExtendedCommunityList
	isSet bool
}

func (v NullableConfigPutRequestExtendedCommunityList) Get() *ConfigPutRequestExtendedCommunityList {
	return v.value
}

func (v *NullableConfigPutRequestExtendedCommunityList) Set(val *ConfigPutRequestExtendedCommunityList) {
	v.value = val
	v.isSet = true
}

func (v NullableConfigPutRequestExtendedCommunityList) IsSet() bool {
	return v.isSet
}

func (v *NullableConfigPutRequestExtendedCommunityList) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConfigPutRequestExtendedCommunityList(val *ConfigPutRequestExtendedCommunityList) *NullableConfigPutRequestExtendedCommunityList {
	return &NullableConfigPutRequestExtendedCommunityList{value: val, isSet: true}
}

func (v NullableConfigPutRequestExtendedCommunityList) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConfigPutRequestExtendedCommunityList) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


