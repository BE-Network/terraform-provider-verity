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

// checks if the PacketqueuesPutRequestPacketQueueValueObjectProperties type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PacketqueuesPutRequestPacketQueueValueObjectProperties{}

// PacketqueuesPutRequestPacketQueueValueObjectProperties struct for PacketqueuesPutRequestPacketQueueValueObjectProperties
type PacketqueuesPutRequestPacketQueueValueObjectProperties struct {
	// Default object.
	Isdefault *bool `json:"isdefault,omitempty"`
	// Group
	Group *string `json:"group,omitempty"`
}

// NewPacketqueuesPutRequestPacketQueueValueObjectProperties instantiates a new PacketqueuesPutRequestPacketQueueValueObjectProperties object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPacketqueuesPutRequestPacketQueueValueObjectProperties() *PacketqueuesPutRequestPacketQueueValueObjectProperties {
	this := PacketqueuesPutRequestPacketQueueValueObjectProperties{}
	var isdefault bool = false
	this.Isdefault = &isdefault
	var group string = ""
	this.Group = &group
	return &this
}

// NewPacketqueuesPutRequestPacketQueueValueObjectPropertiesWithDefaults instantiates a new PacketqueuesPutRequestPacketQueueValueObjectProperties object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPacketqueuesPutRequestPacketQueueValueObjectPropertiesWithDefaults() *PacketqueuesPutRequestPacketQueueValueObjectProperties {
	this := PacketqueuesPutRequestPacketQueueValueObjectProperties{}
	var isdefault bool = false
	this.Isdefault = &isdefault
	var group string = ""
	this.Group = &group
	return &this
}

// GetIsdefault returns the Isdefault field value if set, zero value otherwise.
func (o *PacketqueuesPutRequestPacketQueueValueObjectProperties) GetIsdefault() bool {
	if o == nil || IsNil(o.Isdefault) {
		var ret bool
		return ret
	}
	return *o.Isdefault
}

// GetIsdefaultOk returns a tuple with the Isdefault field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PacketqueuesPutRequestPacketQueueValueObjectProperties) GetIsdefaultOk() (*bool, bool) {
	if o == nil || IsNil(o.Isdefault) {
		return nil, false
	}
	return o.Isdefault, true
}

// HasIsdefault returns a boolean if a field has been set.
func (o *PacketqueuesPutRequestPacketQueueValueObjectProperties) HasIsdefault() bool {
	if o != nil && !IsNil(o.Isdefault) {
		return true
	}

	return false
}

// SetIsdefault gets a reference to the given bool and assigns it to the Isdefault field.
func (o *PacketqueuesPutRequestPacketQueueValueObjectProperties) SetIsdefault(v bool) {
	o.Isdefault = &v
}

// GetGroup returns the Group field value if set, zero value otherwise.
func (o *PacketqueuesPutRequestPacketQueueValueObjectProperties) GetGroup() string {
	if o == nil || IsNil(o.Group) {
		var ret string
		return ret
	}
	return *o.Group
}

// GetGroupOk returns a tuple with the Group field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PacketqueuesPutRequestPacketQueueValueObjectProperties) GetGroupOk() (*string, bool) {
	if o == nil || IsNil(o.Group) {
		return nil, false
	}
	return o.Group, true
}

// HasGroup returns a boolean if a field has been set.
func (o *PacketqueuesPutRequestPacketQueueValueObjectProperties) HasGroup() bool {
	if o != nil && !IsNil(o.Group) {
		return true
	}

	return false
}

// SetGroup gets a reference to the given string and assigns it to the Group field.
func (o *PacketqueuesPutRequestPacketQueueValueObjectProperties) SetGroup(v string) {
	o.Group = &v
}

func (o PacketqueuesPutRequestPacketQueueValueObjectProperties) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PacketqueuesPutRequestPacketQueueValueObjectProperties) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Isdefault) {
		toSerialize["isdefault"] = o.Isdefault
	}
	if !IsNil(o.Group) {
		toSerialize["group"] = o.Group
	}
	return toSerialize, nil
}

type NullablePacketqueuesPutRequestPacketQueueValueObjectProperties struct {
	value *PacketqueuesPutRequestPacketQueueValueObjectProperties
	isSet bool
}

func (v NullablePacketqueuesPutRequestPacketQueueValueObjectProperties) Get() *PacketqueuesPutRequestPacketQueueValueObjectProperties {
	return v.value
}

func (v *NullablePacketqueuesPutRequestPacketQueueValueObjectProperties) Set(val *PacketqueuesPutRequestPacketQueueValueObjectProperties) {
	v.value = val
	v.isSet = true
}

func (v NullablePacketqueuesPutRequestPacketQueueValueObjectProperties) IsSet() bool {
	return v.isSet
}

func (v *NullablePacketqueuesPutRequestPacketQueueValueObjectProperties) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePacketqueuesPutRequestPacketQueueValueObjectProperties(val *PacketqueuesPutRequestPacketQueueValueObjectProperties) *NullablePacketqueuesPutRequestPacketQueueValueObjectProperties {
	return &NullablePacketqueuesPutRequestPacketQueueValueObjectProperties{value: val, isSet: true}
}

func (v NullablePacketqueuesPutRequestPacketQueueValueObjectProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePacketqueuesPutRequestPacketQueueValueObjectProperties) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


