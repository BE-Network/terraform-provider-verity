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

// checks if the ConfigPutRequestMacFilterMacFilterNameAclInner type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ConfigPutRequestMacFilterMacFilterNameAclInner{}

// ConfigPutRequestMacFilterMacFilterNameAclInner struct for ConfigPutRequestMacFilterMacFilterNameAclInner
type ConfigPutRequestMacFilterMacFilterNameAclInner struct {
	// MAC address descriptor including colons example 01:23:45:67:9a:ab. and * notation accepted example 12:*
	FilterNumMac *string `json:"filter_num_mac,omitempty"`
	// Hexidecimal mask including colons example ff:ff:fe:00:00:00. /n and * notation accepted example /16 or 12:*
	FilterNumMask *string `json:"filter_num_mask,omitempty"`
	// Enable of this MAC Filter 
	FilterNumEnable *bool `json:"filter_num_enable,omitempty"`
	// The index identifying the object. Zero if you want to add an object to the list.
	Index *int32 `json:"index,omitempty"`
}

// NewConfigPutRequestMacFilterMacFilterNameAclInner instantiates a new ConfigPutRequestMacFilterMacFilterNameAclInner object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConfigPutRequestMacFilterMacFilterNameAclInner() *ConfigPutRequestMacFilterMacFilterNameAclInner {
	this := ConfigPutRequestMacFilterMacFilterNameAclInner{}
	var filterNumMac string = ""
	this.FilterNumMac = &filterNumMac
	var filterNumMask string = ""
	this.FilterNumMask = &filterNumMask
	var filterNumEnable bool = false
	this.FilterNumEnable = &filterNumEnable
	return &this
}

// NewConfigPutRequestMacFilterMacFilterNameAclInnerWithDefaults instantiates a new ConfigPutRequestMacFilterMacFilterNameAclInner object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConfigPutRequestMacFilterMacFilterNameAclInnerWithDefaults() *ConfigPutRequestMacFilterMacFilterNameAclInner {
	this := ConfigPutRequestMacFilterMacFilterNameAclInner{}
	var filterNumMac string = ""
	this.FilterNumMac = &filterNumMac
	var filterNumMask string = ""
	this.FilterNumMask = &filterNumMask
	var filterNumEnable bool = false
	this.FilterNumEnable = &filterNumEnable
	return &this
}

// GetFilterNumMac returns the FilterNumMac field value if set, zero value otherwise.
func (o *ConfigPutRequestMacFilterMacFilterNameAclInner) GetFilterNumMac() string {
	if o == nil || IsNil(o.FilterNumMac) {
		var ret string
		return ret
	}
	return *o.FilterNumMac
}

// GetFilterNumMacOk returns a tuple with the FilterNumMac field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestMacFilterMacFilterNameAclInner) GetFilterNumMacOk() (*string, bool) {
	if o == nil || IsNil(o.FilterNumMac) {
		return nil, false
	}
	return o.FilterNumMac, true
}

// HasFilterNumMac returns a boolean if a field has been set.
func (o *ConfigPutRequestMacFilterMacFilterNameAclInner) HasFilterNumMac() bool {
	if o != nil && !IsNil(o.FilterNumMac) {
		return true
	}

	return false
}

// SetFilterNumMac gets a reference to the given string and assigns it to the FilterNumMac field.
func (o *ConfigPutRequestMacFilterMacFilterNameAclInner) SetFilterNumMac(v string) {
	o.FilterNumMac = &v
}

// GetFilterNumMask returns the FilterNumMask field value if set, zero value otherwise.
func (o *ConfigPutRequestMacFilterMacFilterNameAclInner) GetFilterNumMask() string {
	if o == nil || IsNil(o.FilterNumMask) {
		var ret string
		return ret
	}
	return *o.FilterNumMask
}

// GetFilterNumMaskOk returns a tuple with the FilterNumMask field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestMacFilterMacFilterNameAclInner) GetFilterNumMaskOk() (*string, bool) {
	if o == nil || IsNil(o.FilterNumMask) {
		return nil, false
	}
	return o.FilterNumMask, true
}

// HasFilterNumMask returns a boolean if a field has been set.
func (o *ConfigPutRequestMacFilterMacFilterNameAclInner) HasFilterNumMask() bool {
	if o != nil && !IsNil(o.FilterNumMask) {
		return true
	}

	return false
}

// SetFilterNumMask gets a reference to the given string and assigns it to the FilterNumMask field.
func (o *ConfigPutRequestMacFilterMacFilterNameAclInner) SetFilterNumMask(v string) {
	o.FilterNumMask = &v
}

// GetFilterNumEnable returns the FilterNumEnable field value if set, zero value otherwise.
func (o *ConfigPutRequestMacFilterMacFilterNameAclInner) GetFilterNumEnable() bool {
	if o == nil || IsNil(o.FilterNumEnable) {
		var ret bool
		return ret
	}
	return *o.FilterNumEnable
}

// GetFilterNumEnableOk returns a tuple with the FilterNumEnable field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestMacFilterMacFilterNameAclInner) GetFilterNumEnableOk() (*bool, bool) {
	if o == nil || IsNil(o.FilterNumEnable) {
		return nil, false
	}
	return o.FilterNumEnable, true
}

// HasFilterNumEnable returns a boolean if a field has been set.
func (o *ConfigPutRequestMacFilterMacFilterNameAclInner) HasFilterNumEnable() bool {
	if o != nil && !IsNil(o.FilterNumEnable) {
		return true
	}

	return false
}

// SetFilterNumEnable gets a reference to the given bool and assigns it to the FilterNumEnable field.
func (o *ConfigPutRequestMacFilterMacFilterNameAclInner) SetFilterNumEnable(v bool) {
	o.FilterNumEnable = &v
}

// GetIndex returns the Index field value if set, zero value otherwise.
func (o *ConfigPutRequestMacFilterMacFilterNameAclInner) GetIndex() int32 {
	if o == nil || IsNil(o.Index) {
		var ret int32
		return ret
	}
	return *o.Index
}

// GetIndexOk returns a tuple with the Index field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestMacFilterMacFilterNameAclInner) GetIndexOk() (*int32, bool) {
	if o == nil || IsNil(o.Index) {
		return nil, false
	}
	return o.Index, true
}

// HasIndex returns a boolean if a field has been set.
func (o *ConfigPutRequestMacFilterMacFilterNameAclInner) HasIndex() bool {
	if o != nil && !IsNil(o.Index) {
		return true
	}

	return false
}

// SetIndex gets a reference to the given int32 and assigns it to the Index field.
func (o *ConfigPutRequestMacFilterMacFilterNameAclInner) SetIndex(v int32) {
	o.Index = &v
}

func (o ConfigPutRequestMacFilterMacFilterNameAclInner) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ConfigPutRequestMacFilterMacFilterNameAclInner) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.FilterNumMac) {
		toSerialize["filter_num_mac"] = o.FilterNumMac
	}
	if !IsNil(o.FilterNumMask) {
		toSerialize["filter_num_mask"] = o.FilterNumMask
	}
	if !IsNil(o.FilterNumEnable) {
		toSerialize["filter_num_enable"] = o.FilterNumEnable
	}
	if !IsNil(o.Index) {
		toSerialize["index"] = o.Index
	}
	return toSerialize, nil
}

type NullableConfigPutRequestMacFilterMacFilterNameAclInner struct {
	value *ConfigPutRequestMacFilterMacFilterNameAclInner
	isSet bool
}

func (v NullableConfigPutRequestMacFilterMacFilterNameAclInner) Get() *ConfigPutRequestMacFilterMacFilterNameAclInner {
	return v.value
}

func (v *NullableConfigPutRequestMacFilterMacFilterNameAclInner) Set(val *ConfigPutRequestMacFilterMacFilterNameAclInner) {
	v.value = val
	v.isSet = true
}

func (v NullableConfigPutRequestMacFilterMacFilterNameAclInner) IsSet() bool {
	return v.isSet
}

func (v *NullableConfigPutRequestMacFilterMacFilterNameAclInner) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConfigPutRequestMacFilterMacFilterNameAclInner(val *ConfigPutRequestMacFilterMacFilterNameAclInner) *NullableConfigPutRequestMacFilterMacFilterNameAclInner {
	return &NullableConfigPutRequestMacFilterMacFilterNameAclInner{value: val, isSet: true}
}

func (v NullableConfigPutRequestMacFilterMacFilterNameAclInner) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConfigPutRequestMacFilterMacFilterNameAclInner) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


