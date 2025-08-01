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

// checks if the AclsPutRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &AclsPutRequest{}

// AclsPutRequest struct for AclsPutRequest
type AclsPutRequest struct {
	IpFilter *map[string]AclsPutRequestIpFilterValue `json:"ip_filter,omitempty"`
}

// NewAclsPutRequest instantiates a new AclsPutRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAclsPutRequest() *AclsPutRequest {
	this := AclsPutRequest{}
	return &this
}

// NewAclsPutRequestWithDefaults instantiates a new AclsPutRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAclsPutRequestWithDefaults() *AclsPutRequest {
	this := AclsPutRequest{}
	return &this
}

// GetIpFilter returns the IpFilter field value if set, zero value otherwise.
func (o *AclsPutRequest) GetIpFilter() map[string]AclsPutRequestIpFilterValue {
	if o == nil || IsNil(o.IpFilter) {
		var ret map[string]AclsPutRequestIpFilterValue
		return ret
	}
	return *o.IpFilter
}

// GetIpFilterOk returns a tuple with the IpFilter field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AclsPutRequest) GetIpFilterOk() (*map[string]AclsPutRequestIpFilterValue, bool) {
	if o == nil || IsNil(o.IpFilter) {
		return nil, false
	}
	return o.IpFilter, true
}

// HasIpFilter returns a boolean if a field has been set.
func (o *AclsPutRequest) HasIpFilter() bool {
	if o != nil && !IsNil(o.IpFilter) {
		return true
	}

	return false
}

// SetIpFilter gets a reference to the given map[string]AclsPutRequestIpFilterValue and assigns it to the IpFilter field.
func (o *AclsPutRequest) SetIpFilter(v map[string]AclsPutRequestIpFilterValue) {
	o.IpFilter = &v
}

func (o AclsPutRequest) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o AclsPutRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.IpFilter) {
		toSerialize["ip_filter"] = o.IpFilter
	}
	return toSerialize, nil
}

type NullableAclsPutRequest struct {
	value *AclsPutRequest
	isSet bool
}

func (v NullableAclsPutRequest) Get() *AclsPutRequest {
	return v.value
}

func (v *NullableAclsPutRequest) Set(val *AclsPutRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableAclsPutRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableAclsPutRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAclsPutRequest(val *AclsPutRequest) *NullableAclsPutRequest {
	return &NullableAclsPutRequest{value: val, isSet: true}
}

func (v NullableAclsPutRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAclsPutRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


