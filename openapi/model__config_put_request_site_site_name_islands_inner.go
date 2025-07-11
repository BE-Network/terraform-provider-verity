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

// checks if the ConfigPutRequestSiteSiteNameIslandsInner type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ConfigPutRequestSiteSiteNameIslandsInner{}

// ConfigPutRequestSiteSiteNameIslandsInner struct for ConfigPutRequestSiteSiteNameIslandsInner
type ConfigPutRequestSiteSiteNameIslandsInner struct {
	// TOI Endpoint
	ToiSwitchpoint *string `json:"toi_switchpoint,omitempty"`
	// Object type for toi_switchpoint field
	ToiSwitchpointRefType *string `json:"toi_switchpoint_ref_type_,omitempty"`
	// The index identifying the object. Zero if you want to add an object to the list.
	Index *int32 `json:"index,omitempty"`
}

// NewConfigPutRequestSiteSiteNameIslandsInner instantiates a new ConfigPutRequestSiteSiteNameIslandsInner object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConfigPutRequestSiteSiteNameIslandsInner() *ConfigPutRequestSiteSiteNameIslandsInner {
	this := ConfigPutRequestSiteSiteNameIslandsInner{}
	var toiSwitchpoint string = ""
	this.ToiSwitchpoint = &toiSwitchpoint
	return &this
}

// NewConfigPutRequestSiteSiteNameIslandsInnerWithDefaults instantiates a new ConfigPutRequestSiteSiteNameIslandsInner object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConfigPutRequestSiteSiteNameIslandsInnerWithDefaults() *ConfigPutRequestSiteSiteNameIslandsInner {
	this := ConfigPutRequestSiteSiteNameIslandsInner{}
	var toiSwitchpoint string = ""
	this.ToiSwitchpoint = &toiSwitchpoint
	return &this
}

// GetToiSwitchpoint returns the ToiSwitchpoint field value if set, zero value otherwise.
func (o *ConfigPutRequestSiteSiteNameIslandsInner) GetToiSwitchpoint() string {
	if o == nil || IsNil(o.ToiSwitchpoint) {
		var ret string
		return ret
	}
	return *o.ToiSwitchpoint
}

// GetToiSwitchpointOk returns a tuple with the ToiSwitchpoint field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestSiteSiteNameIslandsInner) GetToiSwitchpointOk() (*string, bool) {
	if o == nil || IsNil(o.ToiSwitchpoint) {
		return nil, false
	}
	return o.ToiSwitchpoint, true
}

// HasToiSwitchpoint returns a boolean if a field has been set.
func (o *ConfigPutRequestSiteSiteNameIslandsInner) HasToiSwitchpoint() bool {
	if o != nil && !IsNil(o.ToiSwitchpoint) {
		return true
	}

	return false
}

// SetToiSwitchpoint gets a reference to the given string and assigns it to the ToiSwitchpoint field.
func (o *ConfigPutRequestSiteSiteNameIslandsInner) SetToiSwitchpoint(v string) {
	o.ToiSwitchpoint = &v
}

// GetToiSwitchpointRefType returns the ToiSwitchpointRefType field value if set, zero value otherwise.
func (o *ConfigPutRequestSiteSiteNameIslandsInner) GetToiSwitchpointRefType() string {
	if o == nil || IsNil(o.ToiSwitchpointRefType) {
		var ret string
		return ret
	}
	return *o.ToiSwitchpointRefType
}

// GetToiSwitchpointRefTypeOk returns a tuple with the ToiSwitchpointRefType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestSiteSiteNameIslandsInner) GetToiSwitchpointRefTypeOk() (*string, bool) {
	if o == nil || IsNil(o.ToiSwitchpointRefType) {
		return nil, false
	}
	return o.ToiSwitchpointRefType, true
}

// HasToiSwitchpointRefType returns a boolean if a field has been set.
func (o *ConfigPutRequestSiteSiteNameIslandsInner) HasToiSwitchpointRefType() bool {
	if o != nil && !IsNil(o.ToiSwitchpointRefType) {
		return true
	}

	return false
}

// SetToiSwitchpointRefType gets a reference to the given string and assigns it to the ToiSwitchpointRefType field.
func (o *ConfigPutRequestSiteSiteNameIslandsInner) SetToiSwitchpointRefType(v string) {
	o.ToiSwitchpointRefType = &v
}

// GetIndex returns the Index field value if set, zero value otherwise.
func (o *ConfigPutRequestSiteSiteNameIslandsInner) GetIndex() int32 {
	if o == nil || IsNil(o.Index) {
		var ret int32
		return ret
	}
	return *o.Index
}

// GetIndexOk returns a tuple with the Index field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestSiteSiteNameIslandsInner) GetIndexOk() (*int32, bool) {
	if o == nil || IsNil(o.Index) {
		return nil, false
	}
	return o.Index, true
}

// HasIndex returns a boolean if a field has been set.
func (o *ConfigPutRequestSiteSiteNameIslandsInner) HasIndex() bool {
	if o != nil && !IsNil(o.Index) {
		return true
	}

	return false
}

// SetIndex gets a reference to the given int32 and assigns it to the Index field.
func (o *ConfigPutRequestSiteSiteNameIslandsInner) SetIndex(v int32) {
	o.Index = &v
}

func (o ConfigPutRequestSiteSiteNameIslandsInner) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ConfigPutRequestSiteSiteNameIslandsInner) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.ToiSwitchpoint) {
		toSerialize["toi_switchpoint"] = o.ToiSwitchpoint
	}
	if !IsNil(o.ToiSwitchpointRefType) {
		toSerialize["toi_switchpoint_ref_type_"] = o.ToiSwitchpointRefType
	}
	if !IsNil(o.Index) {
		toSerialize["index"] = o.Index
	}
	return toSerialize, nil
}

type NullableConfigPutRequestSiteSiteNameIslandsInner struct {
	value *ConfigPutRequestSiteSiteNameIslandsInner
	isSet bool
}

func (v NullableConfigPutRequestSiteSiteNameIslandsInner) Get() *ConfigPutRequestSiteSiteNameIslandsInner {
	return v.value
}

func (v *NullableConfigPutRequestSiteSiteNameIslandsInner) Set(val *ConfigPutRequestSiteSiteNameIslandsInner) {
	v.value = val
	v.isSet = true
}

func (v NullableConfigPutRequestSiteSiteNameIslandsInner) IsSet() bool {
	return v.isSet
}

func (v *NullableConfigPutRequestSiteSiteNameIslandsInner) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConfigPutRequestSiteSiteNameIslandsInner(val *ConfigPutRequestSiteSiteNameIslandsInner) *NullableConfigPutRequestSiteSiteNameIslandsInner {
	return &NullableConfigPutRequestSiteSiteNameIslandsInner{value: val, isSet: true}
}

func (v NullableConfigPutRequestSiteSiteNameIslandsInner) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConfigPutRequestSiteSiteNameIslandsInner) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


