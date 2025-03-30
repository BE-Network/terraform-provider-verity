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

// checks if the ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner{}

// ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner struct for ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner
type ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner struct {
	// Enable row
	RowNumEnable *bool `json:"row_num_enable,omitempty"`
	// Choose a Service to connect
	RowNumService *string `json:"row_num_service,omitempty"`
	// Object type for row_num_service field
	RowNumServiceRefType *string `json:"row_num_service_ref_type_,omitempty"`
	// Choose an external vlan A value of 0 will make the VLAN untagged, while in case null is provided, the VLAN will be the one associated with the service.
	RowNumExternalVlan NullableInt32 `json:"row_num_external_vlan,omitempty"`
	// The index identifying the object. Zero if you want to add an object to the list.
	Index *int32 `json:"index,omitempty"`
}

// NewConfigPutRequestEthPortProfileEthPortProfileNameServicesInner instantiates a new ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConfigPutRequestEthPortProfileEthPortProfileNameServicesInner() *ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner {
	this := ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner{}
	var rowNumEnable bool = false
	this.RowNumEnable = &rowNumEnable
	var rowNumService string = ""
	this.RowNumService = &rowNumService
	return &this
}

// NewConfigPutRequestEthPortProfileEthPortProfileNameServicesInnerWithDefaults instantiates a new ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConfigPutRequestEthPortProfileEthPortProfileNameServicesInnerWithDefaults() *ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner {
	this := ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner{}
	var rowNumEnable bool = false
	this.RowNumEnable = &rowNumEnable
	var rowNumService string = ""
	this.RowNumService = &rowNumService
	return &this
}

// GetRowNumEnable returns the RowNumEnable field value if set, zero value otherwise.
func (o *ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner) GetRowNumEnable() bool {
	if o == nil || IsNil(o.RowNumEnable) {
		var ret bool
		return ret
	}
	return *o.RowNumEnable
}

// GetRowNumEnableOk returns a tuple with the RowNumEnable field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner) GetRowNumEnableOk() (*bool, bool) {
	if o == nil || IsNil(o.RowNumEnable) {
		return nil, false
	}
	return o.RowNumEnable, true
}

// HasRowNumEnable returns a boolean if a field has been set.
func (o *ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner) HasRowNumEnable() bool {
	if o != nil && !IsNil(o.RowNumEnable) {
		return true
	}

	return false
}

// SetRowNumEnable gets a reference to the given bool and assigns it to the RowNumEnable field.
func (o *ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner) SetRowNumEnable(v bool) {
	o.RowNumEnable = &v
}

// GetRowNumService returns the RowNumService field value if set, zero value otherwise.
func (o *ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner) GetRowNumService() string {
	if o == nil || IsNil(o.RowNumService) {
		var ret string
		return ret
	}
	return *o.RowNumService
}

// GetRowNumServiceOk returns a tuple with the RowNumService field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner) GetRowNumServiceOk() (*string, bool) {
	if o == nil || IsNil(o.RowNumService) {
		return nil, false
	}
	return o.RowNumService, true
}

// HasRowNumService returns a boolean if a field has been set.
func (o *ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner) HasRowNumService() bool {
	if o != nil && !IsNil(o.RowNumService) {
		return true
	}

	return false
}

// SetRowNumService gets a reference to the given string and assigns it to the RowNumService field.
func (o *ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner) SetRowNumService(v string) {
	o.RowNumService = &v
}

// GetRowNumServiceRefType returns the RowNumServiceRefType field value if set, zero value otherwise.
func (o *ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner) GetRowNumServiceRefType() string {
	if o == nil || IsNil(o.RowNumServiceRefType) {
		var ret string
		return ret
	}
	return *o.RowNumServiceRefType
}

// GetRowNumServiceRefTypeOk returns a tuple with the RowNumServiceRefType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner) GetRowNumServiceRefTypeOk() (*string, bool) {
	if o == nil || IsNil(o.RowNumServiceRefType) {
		return nil, false
	}
	return o.RowNumServiceRefType, true
}

// HasRowNumServiceRefType returns a boolean if a field has been set.
func (o *ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner) HasRowNumServiceRefType() bool {
	if o != nil && !IsNil(o.RowNumServiceRefType) {
		return true
	}

	return false
}

// SetRowNumServiceRefType gets a reference to the given string and assigns it to the RowNumServiceRefType field.
func (o *ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner) SetRowNumServiceRefType(v string) {
	o.RowNumServiceRefType = &v
}

// GetRowNumExternalVlan returns the RowNumExternalVlan field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner) GetRowNumExternalVlan() int32 {
	if o == nil || IsNil(o.RowNumExternalVlan.Get()) {
		var ret int32
		return ret
	}
	return *o.RowNumExternalVlan.Get()
}

// GetRowNumExternalVlanOk returns a tuple with the RowNumExternalVlan field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner) GetRowNumExternalVlanOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.RowNumExternalVlan.Get(), o.RowNumExternalVlan.IsSet()
}

// HasRowNumExternalVlan returns a boolean if a field has been set.
func (o *ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner) HasRowNumExternalVlan() bool {
	if o != nil && o.RowNumExternalVlan.IsSet() {
		return true
	}

	return false
}

// SetRowNumExternalVlan gets a reference to the given NullableInt32 and assigns it to the RowNumExternalVlan field.
func (o *ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner) SetRowNumExternalVlan(v int32) {
	o.RowNumExternalVlan.Set(&v)
}
// SetRowNumExternalVlanNil sets the value for RowNumExternalVlan to be an explicit nil
func (o *ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner) SetRowNumExternalVlanNil() {
	o.RowNumExternalVlan.Set(nil)
}

// UnsetRowNumExternalVlan ensures that no value is present for RowNumExternalVlan, not even an explicit nil
func (o *ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner) UnsetRowNumExternalVlan() {
	o.RowNumExternalVlan.Unset()
}

// GetIndex returns the Index field value if set, zero value otherwise.
func (o *ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner) GetIndex() int32 {
	if o == nil || IsNil(o.Index) {
		var ret int32
		return ret
	}
	return *o.Index
}

// GetIndexOk returns a tuple with the Index field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner) GetIndexOk() (*int32, bool) {
	if o == nil || IsNil(o.Index) {
		return nil, false
	}
	return o.Index, true
}

// HasIndex returns a boolean if a field has been set.
func (o *ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner) HasIndex() bool {
	if o != nil && !IsNil(o.Index) {
		return true
	}

	return false
}

// SetIndex gets a reference to the given int32 and assigns it to the Index field.
func (o *ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner) SetIndex(v int32) {
	o.Index = &v
}

func (o ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.RowNumEnable) {
		toSerialize["row_num_enable"] = o.RowNumEnable
	}
	if !IsNil(o.RowNumService) {
		toSerialize["row_num_service"] = o.RowNumService
	}
	if !IsNil(o.RowNumServiceRefType) {
		toSerialize["row_num_service_ref_type_"] = o.RowNumServiceRefType
	}
	if o.RowNumExternalVlan.IsSet() {
		toSerialize["row_num_external_vlan"] = o.RowNumExternalVlan.Get()
	}
	if !IsNil(o.Index) {
		toSerialize["index"] = o.Index
	}
	return toSerialize, nil
}

type NullableConfigPutRequestEthPortProfileEthPortProfileNameServicesInner struct {
	value *ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner
	isSet bool
}

func (v NullableConfigPutRequestEthPortProfileEthPortProfileNameServicesInner) Get() *ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner {
	return v.value
}

func (v *NullableConfigPutRequestEthPortProfileEthPortProfileNameServicesInner) Set(val *ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner) {
	v.value = val
	v.isSet = true
}

func (v NullableConfigPutRequestEthPortProfileEthPortProfileNameServicesInner) IsSet() bool {
	return v.isSet
}

func (v *NullableConfigPutRequestEthPortProfileEthPortProfileNameServicesInner) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConfigPutRequestEthPortProfileEthPortProfileNameServicesInner(val *ConfigPutRequestEthPortProfileEthPortProfileNameServicesInner) *NullableConfigPutRequestEthPortProfileEthPortProfileNameServicesInner {
	return &NullableConfigPutRequestEthPortProfileEthPortProfileNameServicesInner{value: val, isSet: true}
}

func (v NullableConfigPutRequestEthPortProfileEthPortProfileNameServicesInner) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConfigPutRequestEthPortProfileEthPortProfileNameServicesInner) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


