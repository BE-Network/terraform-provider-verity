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

// checks if the SwitchpointsPutRequestSwitchpointValueBadgesInner type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &SwitchpointsPutRequestSwitchpointValueBadgesInner{}

// SwitchpointsPutRequestSwitchpointValueBadgesInner struct for SwitchpointsPutRequestSwitchpointValueBadgesInner
type SwitchpointsPutRequestSwitchpointValueBadgesInner struct {
	// Enable of this POTS port
	Badge *string `json:"badge,omitempty"`
	// Object type for badge field
	BadgeRefType *string `json:"badge_ref_type_,omitempty"`
	// The index identifying the object. Zero if you want to add an object to the list.
	Index *int32 `json:"index,omitempty"`
}

// NewSwitchpointsPutRequestSwitchpointValueBadgesInner instantiates a new SwitchpointsPutRequestSwitchpointValueBadgesInner object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSwitchpointsPutRequestSwitchpointValueBadgesInner() *SwitchpointsPutRequestSwitchpointValueBadgesInner {
	this := SwitchpointsPutRequestSwitchpointValueBadgesInner{}
	var badge string = ""
	this.Badge = &badge
	return &this
}

// NewSwitchpointsPutRequestSwitchpointValueBadgesInnerWithDefaults instantiates a new SwitchpointsPutRequestSwitchpointValueBadgesInner object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSwitchpointsPutRequestSwitchpointValueBadgesInnerWithDefaults() *SwitchpointsPutRequestSwitchpointValueBadgesInner {
	this := SwitchpointsPutRequestSwitchpointValueBadgesInner{}
	var badge string = ""
	this.Badge = &badge
	return &this
}

// GetBadge returns the Badge field value if set, zero value otherwise.
func (o *SwitchpointsPutRequestSwitchpointValueBadgesInner) GetBadge() string {
	if o == nil || IsNil(o.Badge) {
		var ret string
		return ret
	}
	return *o.Badge
}

// GetBadgeOk returns a tuple with the Badge field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SwitchpointsPutRequestSwitchpointValueBadgesInner) GetBadgeOk() (*string, bool) {
	if o == nil || IsNil(o.Badge) {
		return nil, false
	}
	return o.Badge, true
}

// HasBadge returns a boolean if a field has been set.
func (o *SwitchpointsPutRequestSwitchpointValueBadgesInner) HasBadge() bool {
	if o != nil && !IsNil(o.Badge) {
		return true
	}

	return false
}

// SetBadge gets a reference to the given string and assigns it to the Badge field.
func (o *SwitchpointsPutRequestSwitchpointValueBadgesInner) SetBadge(v string) {
	o.Badge = &v
}

// GetBadgeRefType returns the BadgeRefType field value if set, zero value otherwise.
func (o *SwitchpointsPutRequestSwitchpointValueBadgesInner) GetBadgeRefType() string {
	if o == nil || IsNil(o.BadgeRefType) {
		var ret string
		return ret
	}
	return *o.BadgeRefType
}

// GetBadgeRefTypeOk returns a tuple with the BadgeRefType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SwitchpointsPutRequestSwitchpointValueBadgesInner) GetBadgeRefTypeOk() (*string, bool) {
	if o == nil || IsNil(o.BadgeRefType) {
		return nil, false
	}
	return o.BadgeRefType, true
}

// HasBadgeRefType returns a boolean if a field has been set.
func (o *SwitchpointsPutRequestSwitchpointValueBadgesInner) HasBadgeRefType() bool {
	if o != nil && !IsNil(o.BadgeRefType) {
		return true
	}

	return false
}

// SetBadgeRefType gets a reference to the given string and assigns it to the BadgeRefType field.
func (o *SwitchpointsPutRequestSwitchpointValueBadgesInner) SetBadgeRefType(v string) {
	o.BadgeRefType = &v
}

// GetIndex returns the Index field value if set, zero value otherwise.
func (o *SwitchpointsPutRequestSwitchpointValueBadgesInner) GetIndex() int32 {
	if o == nil || IsNil(o.Index) {
		var ret int32
		return ret
	}
	return *o.Index
}

// GetIndexOk returns a tuple with the Index field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *SwitchpointsPutRequestSwitchpointValueBadgesInner) GetIndexOk() (*int32, bool) {
	if o == nil || IsNil(o.Index) {
		return nil, false
	}
	return o.Index, true
}

// HasIndex returns a boolean if a field has been set.
func (o *SwitchpointsPutRequestSwitchpointValueBadgesInner) HasIndex() bool {
	if o != nil && !IsNil(o.Index) {
		return true
	}

	return false
}

// SetIndex gets a reference to the given int32 and assigns it to the Index field.
func (o *SwitchpointsPutRequestSwitchpointValueBadgesInner) SetIndex(v int32) {
	o.Index = &v
}

func (o SwitchpointsPutRequestSwitchpointValueBadgesInner) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o SwitchpointsPutRequestSwitchpointValueBadgesInner) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Badge) {
		toSerialize["badge"] = o.Badge
	}
	if !IsNil(o.BadgeRefType) {
		toSerialize["badge_ref_type_"] = o.BadgeRefType
	}
	if !IsNil(o.Index) {
		toSerialize["index"] = o.Index
	}
	return toSerialize, nil
}

type NullableSwitchpointsPutRequestSwitchpointValueBadgesInner struct {
	value *SwitchpointsPutRequestSwitchpointValueBadgesInner
	isSet bool
}

func (v NullableSwitchpointsPutRequestSwitchpointValueBadgesInner) Get() *SwitchpointsPutRequestSwitchpointValueBadgesInner {
	return v.value
}

func (v *NullableSwitchpointsPutRequestSwitchpointValueBadgesInner) Set(val *SwitchpointsPutRequestSwitchpointValueBadgesInner) {
	v.value = val
	v.isSet = true
}

func (v NullableSwitchpointsPutRequestSwitchpointValueBadgesInner) IsSet() bool {
	return v.isSet
}

func (v *NullableSwitchpointsPutRequestSwitchpointValueBadgesInner) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSwitchpointsPutRequestSwitchpointValueBadgesInner(val *SwitchpointsPutRequestSwitchpointValueBadgesInner) *NullableSwitchpointsPutRequestSwitchpointValueBadgesInner {
	return &NullableSwitchpointsPutRequestSwitchpointValueBadgesInner{value: val, isSet: true}
}

func (v NullableSwitchpointsPutRequestSwitchpointValueBadgesInner) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSwitchpointsPutRequestSwitchpointValueBadgesInner) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


