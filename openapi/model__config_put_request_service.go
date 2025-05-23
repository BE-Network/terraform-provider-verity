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

// checks if the ConfigPutRequestService type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ConfigPutRequestService{}

// ConfigPutRequestService struct for ConfigPutRequestService
type ConfigPutRequestService struct {
	ServiceName *ConfigPutRequestServiceServiceName `json:"service_name,omitempty"`
}

// NewConfigPutRequestService instantiates a new ConfigPutRequestService object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConfigPutRequestService() *ConfigPutRequestService {
	this := ConfigPutRequestService{}
	return &this
}

// NewConfigPutRequestServiceWithDefaults instantiates a new ConfigPutRequestService object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConfigPutRequestServiceWithDefaults() *ConfigPutRequestService {
	this := ConfigPutRequestService{}
	return &this
}

// GetServiceName returns the ServiceName field value if set, zero value otherwise.
func (o *ConfigPutRequestService) GetServiceName() ConfigPutRequestServiceServiceName {
	if o == nil || IsNil(o.ServiceName) {
		var ret ConfigPutRequestServiceServiceName
		return ret
	}
	return *o.ServiceName
}

// GetServiceNameOk returns a tuple with the ServiceName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestService) GetServiceNameOk() (*ConfigPutRequestServiceServiceName, bool) {
	if o == nil || IsNil(o.ServiceName) {
		return nil, false
	}
	return o.ServiceName, true
}

// HasServiceName returns a boolean if a field has been set.
func (o *ConfigPutRequestService) HasServiceName() bool {
	if o != nil && !IsNil(o.ServiceName) {
		return true
	}

	return false
}

// SetServiceName gets a reference to the given ConfigPutRequestServiceServiceName and assigns it to the ServiceName field.
func (o *ConfigPutRequestService) SetServiceName(v ConfigPutRequestServiceServiceName) {
	o.ServiceName = &v
}

func (o ConfigPutRequestService) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ConfigPutRequestService) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.ServiceName) {
		toSerialize["service_name"] = o.ServiceName
	}
	return toSerialize, nil
}

type NullableConfigPutRequestService struct {
	value *ConfigPutRequestService
	isSet bool
}

func (v NullableConfigPutRequestService) Get() *ConfigPutRequestService {
	return v.value
}

func (v *NullableConfigPutRequestService) Set(val *ConfigPutRequestService) {
	v.value = val
	v.isSet = true
}

func (v NullableConfigPutRequestService) IsSet() bool {
	return v.isSet
}

func (v *NullableConfigPutRequestService) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConfigPutRequestService(val *ConfigPutRequestService) *NullableConfigPutRequestService {
	return &NullableConfigPutRequestService{value: val, isSet: true}
}

func (v NullableConfigPutRequestService) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConfigPutRequestService) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


