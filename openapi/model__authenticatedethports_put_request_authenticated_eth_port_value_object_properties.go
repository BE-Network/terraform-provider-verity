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

// checks if the AuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &AuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties{}

// AuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties struct for AuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties
type AuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties struct {
	// Group
	Group *string `json:"group,omitempty"`
	// Defines importance of Link Down on this port
	PortMonitoring *string `json:"port_monitoring,omitempty"`
}

// NewAuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties instantiates a new AuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties() *AuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties {
	this := AuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties{}
	var group string = ""
	this.Group = &group
	var portMonitoring string = ""
	this.PortMonitoring = &portMonitoring
	return &this
}

// NewAuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectPropertiesWithDefaults instantiates a new AuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectPropertiesWithDefaults() *AuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties {
	this := AuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties{}
	var group string = ""
	this.Group = &group
	var portMonitoring string = ""
	this.PortMonitoring = &portMonitoring
	return &this
}

// GetGroup returns the Group field value if set, zero value otherwise.
func (o *AuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties) GetGroup() string {
	if o == nil || IsNil(o.Group) {
		var ret string
		return ret
	}
	return *o.Group
}

// GetGroupOk returns a tuple with the Group field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties) GetGroupOk() (*string, bool) {
	if o == nil || IsNil(o.Group) {
		return nil, false
	}
	return o.Group, true
}

// HasGroup returns a boolean if a field has been set.
func (o *AuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties) HasGroup() bool {
	if o != nil && !IsNil(o.Group) {
		return true
	}

	return false
}

// SetGroup gets a reference to the given string and assigns it to the Group field.
func (o *AuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties) SetGroup(v string) {
	o.Group = &v
}

// GetPortMonitoring returns the PortMonitoring field value if set, zero value otherwise.
func (o *AuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties) GetPortMonitoring() string {
	if o == nil || IsNil(o.PortMonitoring) {
		var ret string
		return ret
	}
	return *o.PortMonitoring
}

// GetPortMonitoringOk returns a tuple with the PortMonitoring field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties) GetPortMonitoringOk() (*string, bool) {
	if o == nil || IsNil(o.PortMonitoring) {
		return nil, false
	}
	return o.PortMonitoring, true
}

// HasPortMonitoring returns a boolean if a field has been set.
func (o *AuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties) HasPortMonitoring() bool {
	if o != nil && !IsNil(o.PortMonitoring) {
		return true
	}

	return false
}

// SetPortMonitoring gets a reference to the given string and assigns it to the PortMonitoring field.
func (o *AuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties) SetPortMonitoring(v string) {
	o.PortMonitoring = &v
}

func (o AuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o AuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Group) {
		toSerialize["group"] = o.Group
	}
	if !IsNil(o.PortMonitoring) {
		toSerialize["port_monitoring"] = o.PortMonitoring
	}
	return toSerialize, nil
}

type NullableAuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties struct {
	value *AuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties
	isSet bool
}

func (v NullableAuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties) Get() *AuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties {
	return v.value
}

func (v *NullableAuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties) Set(val *AuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties) {
	v.value = val
	v.isSet = true
}

func (v NullableAuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties) IsSet() bool {
	return v.isSet
}

func (v *NullableAuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties(val *AuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties) *NullableAuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties {
	return &NullableAuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties{value: val, isSet: true}
}

func (v NullableAuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAuthenticatedethportsPutRequestAuthenticatedEthPortValueObjectProperties) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


