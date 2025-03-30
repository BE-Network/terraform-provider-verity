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

// checks if the BundlesPatchRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &BundlesPatchRequest{}

// BundlesPatchRequest struct for BundlesPatchRequest
type BundlesPatchRequest struct {
	EndpointBundle *map[string]BundlesPatchRequestEndpointBundleValue `json:"endpoint_bundle,omitempty"`
}

// NewBundlesPatchRequest instantiates a new BundlesPatchRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBundlesPatchRequest() *BundlesPatchRequest {
	this := BundlesPatchRequest{}
	return &this
}

// NewBundlesPatchRequestWithDefaults instantiates a new BundlesPatchRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBundlesPatchRequestWithDefaults() *BundlesPatchRequest {
	this := BundlesPatchRequest{}
	return &this
}

// GetEndpointBundle returns the EndpointBundle field value if set, zero value otherwise.
func (o *BundlesPatchRequest) GetEndpointBundle() map[string]BundlesPatchRequestEndpointBundleValue {
	if o == nil || IsNil(o.EndpointBundle) {
		var ret map[string]BundlesPatchRequestEndpointBundleValue
		return ret
	}
	return *o.EndpointBundle
}

// GetEndpointBundleOk returns a tuple with the EndpointBundle field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BundlesPatchRequest) GetEndpointBundleOk() (*map[string]BundlesPatchRequestEndpointBundleValue, bool) {
	if o == nil || IsNil(o.EndpointBundle) {
		return nil, false
	}
	return o.EndpointBundle, true
}

// HasEndpointBundle returns a boolean if a field has been set.
func (o *BundlesPatchRequest) HasEndpointBundle() bool {
	if o != nil && !IsNil(o.EndpointBundle) {
		return true
	}

	return false
}

// SetEndpointBundle gets a reference to the given map[string]BundlesPatchRequestEndpointBundleValue and assigns it to the EndpointBundle field.
func (o *BundlesPatchRequest) SetEndpointBundle(v map[string]BundlesPatchRequestEndpointBundleValue) {
	o.EndpointBundle = &v
}

func (o BundlesPatchRequest) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o BundlesPatchRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.EndpointBundle) {
		toSerialize["endpoint_bundle"] = o.EndpointBundle
	}
	return toSerialize, nil
}

type NullableBundlesPatchRequest struct {
	value *BundlesPatchRequest
	isSet bool
}

func (v NullableBundlesPatchRequest) Get() *BundlesPatchRequest {
	return v.value
}

func (v *NullableBundlesPatchRequest) Set(val *BundlesPatchRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableBundlesPatchRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableBundlesPatchRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBundlesPatchRequest(val *BundlesPatchRequest) *NullableBundlesPatchRequest {
	return &NullableBundlesPatchRequest{value: val, isSet: true}
}

func (v NullableBundlesPatchRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBundlesPatchRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


