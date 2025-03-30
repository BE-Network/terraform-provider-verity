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

// checks if the ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner{}

// ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner struct for ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner
type ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner struct {
	// Enable
	Enable *bool `json:"enable,omitempty"`
	// Route Map Clause is a collection match and set rules
	RouteMapClause *string `json:"route_map_clause,omitempty"`
	// Object type for route_map_clause field
	RouteMapClauseRefType *string `json:"route_map_clause_ref_type_,omitempty"`
	// The index identifying the object. Zero if you want to add an object to the list.
	Index *int32 `json:"index,omitempty"`
}

// NewConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner instantiates a new ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner() *ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner {
	this := ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner{}
	var enable bool = false
	this.Enable = &enable
	var routeMapClause string = ""
	this.RouteMapClause = &routeMapClause
	return &this
}

// NewConfigPutRequestRouteMapRouteMapNameRouteMapClausesInnerWithDefaults instantiates a new ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConfigPutRequestRouteMapRouteMapNameRouteMapClausesInnerWithDefaults() *ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner {
	this := ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner{}
	var enable bool = false
	this.Enable = &enable
	var routeMapClause string = ""
	this.RouteMapClause = &routeMapClause
	return &this
}

// GetEnable returns the Enable field value if set, zero value otherwise.
func (o *ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner) GetEnable() bool {
	if o == nil || IsNil(o.Enable) {
		var ret bool
		return ret
	}
	return *o.Enable
}

// GetEnableOk returns a tuple with the Enable field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner) GetEnableOk() (*bool, bool) {
	if o == nil || IsNil(o.Enable) {
		return nil, false
	}
	return o.Enable, true
}

// HasEnable returns a boolean if a field has been set.
func (o *ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner) HasEnable() bool {
	if o != nil && !IsNil(o.Enable) {
		return true
	}

	return false
}

// SetEnable gets a reference to the given bool and assigns it to the Enable field.
func (o *ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner) SetEnable(v bool) {
	o.Enable = &v
}

// GetRouteMapClause returns the RouteMapClause field value if set, zero value otherwise.
func (o *ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner) GetRouteMapClause() string {
	if o == nil || IsNil(o.RouteMapClause) {
		var ret string
		return ret
	}
	return *o.RouteMapClause
}

// GetRouteMapClauseOk returns a tuple with the RouteMapClause field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner) GetRouteMapClauseOk() (*string, bool) {
	if o == nil || IsNil(o.RouteMapClause) {
		return nil, false
	}
	return o.RouteMapClause, true
}

// HasRouteMapClause returns a boolean if a field has been set.
func (o *ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner) HasRouteMapClause() bool {
	if o != nil && !IsNil(o.RouteMapClause) {
		return true
	}

	return false
}

// SetRouteMapClause gets a reference to the given string and assigns it to the RouteMapClause field.
func (o *ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner) SetRouteMapClause(v string) {
	o.RouteMapClause = &v
}

// GetRouteMapClauseRefType returns the RouteMapClauseRefType field value if set, zero value otherwise.
func (o *ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner) GetRouteMapClauseRefType() string {
	if o == nil || IsNil(o.RouteMapClauseRefType) {
		var ret string
		return ret
	}
	return *o.RouteMapClauseRefType
}

// GetRouteMapClauseRefTypeOk returns a tuple with the RouteMapClauseRefType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner) GetRouteMapClauseRefTypeOk() (*string, bool) {
	if o == nil || IsNil(o.RouteMapClauseRefType) {
		return nil, false
	}
	return o.RouteMapClauseRefType, true
}

// HasRouteMapClauseRefType returns a boolean if a field has been set.
func (o *ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner) HasRouteMapClauseRefType() bool {
	if o != nil && !IsNil(o.RouteMapClauseRefType) {
		return true
	}

	return false
}

// SetRouteMapClauseRefType gets a reference to the given string and assigns it to the RouteMapClauseRefType field.
func (o *ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner) SetRouteMapClauseRefType(v string) {
	o.RouteMapClauseRefType = &v
}

// GetIndex returns the Index field value if set, zero value otherwise.
func (o *ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner) GetIndex() int32 {
	if o == nil || IsNil(o.Index) {
		var ret int32
		return ret
	}
	return *o.Index
}

// GetIndexOk returns a tuple with the Index field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner) GetIndexOk() (*int32, bool) {
	if o == nil || IsNil(o.Index) {
		return nil, false
	}
	return o.Index, true
}

// HasIndex returns a boolean if a field has been set.
func (o *ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner) HasIndex() bool {
	if o != nil && !IsNil(o.Index) {
		return true
	}

	return false
}

// SetIndex gets a reference to the given int32 and assigns it to the Index field.
func (o *ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner) SetIndex(v int32) {
	o.Index = &v
}

func (o ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Enable) {
		toSerialize["enable"] = o.Enable
	}
	if !IsNil(o.RouteMapClause) {
		toSerialize["route_map_clause"] = o.RouteMapClause
	}
	if !IsNil(o.RouteMapClauseRefType) {
		toSerialize["route_map_clause_ref_type_"] = o.RouteMapClauseRefType
	}
	if !IsNil(o.Index) {
		toSerialize["index"] = o.Index
	}
	return toSerialize, nil
}

type NullableConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner struct {
	value *ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner
	isSet bool
}

func (v NullableConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner) Get() *ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner {
	return v.value
}

func (v *NullableConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner) Set(val *ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner) {
	v.value = val
	v.isSet = true
}

func (v NullableConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner) IsSet() bool {
	return v.isSet
}

func (v *NullableConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner(val *ConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner) *NullableConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner {
	return &NullableConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner{value: val, isSet: true}
}

func (v NullableConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConfigPutRequestRouteMapRouteMapNameRouteMapClausesInner) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


