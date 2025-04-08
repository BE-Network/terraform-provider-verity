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

// checks if the BundlesPatchRequestEndpointBundleValueEthPortPathsInner type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &BundlesPatchRequestEndpointBundleValueEthPortPathsInner{}

// BundlesPatchRequestEndpointBundleValueEthPortPathsInner struct for BundlesPatchRequestEndpointBundleValueEthPortPathsInner
type BundlesPatchRequestEndpointBundleValueEthPortPathsInner struct {
	// Eth Port Profile Or LAG for Eth Port
	EthPortNumEthPortProfile *string `json:"eth_port_num_eth_port_profile,omitempty"`
	// Object type for eth_port_num_eth_port_profile field
	EthPortNumEthPortProfileRefType *string `json:"eth_port_num_eth_port_profile_ref_type_,omitempty"`
	// Choose an Eth Port Settings
	EthPortNumEthPortSettings *string `json:"eth_port_num_eth_port_settings,omitempty"`
	// Object type for eth_port_num_eth_port_settings field
	EthPortNumEthPortSettingsRefType *string `json:"eth_port_num_eth_port_settings_ref_type_,omitempty"`
	// Gateway Profile or LAG for Eth Port
	EthPortNumGatewayProfile *string `json:"eth_port_num_gateway_profile,omitempty"`
	// Object type for eth_port_num_gateway_profile field
	EthPortNumGatewayProfileRefType *string `json:"eth_port_num_gateway_profile_ref_type_,omitempty"`
	// The index identifying the object. Zero if you want to add an object to the list.
	Index *int32 `json:"index,omitempty"`
	// The name identifying the port. Used for reference only, it won't actually change the port name.
	PortName *string `json:"port_name,omitempty"`
}

// NewBundlesPatchRequestEndpointBundleValueEthPortPathsInner instantiates a new BundlesPatchRequestEndpointBundleValueEthPortPathsInner object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBundlesPatchRequestEndpointBundleValueEthPortPathsInner() *BundlesPatchRequestEndpointBundleValueEthPortPathsInner {
	this := BundlesPatchRequestEndpointBundleValueEthPortPathsInner{}
	var ethPortNumEthPortProfile string = ""
	this.EthPortNumEthPortProfile = &ethPortNumEthPortProfile
	var ethPortNumEthPortSettings string = ""
	this.EthPortNumEthPortSettings = &ethPortNumEthPortSettings
	var ethPortNumGatewayProfile string = ""
	this.EthPortNumGatewayProfile = &ethPortNumGatewayProfile
	return &this
}

// NewBundlesPatchRequestEndpointBundleValueEthPortPathsInnerWithDefaults instantiates a new BundlesPatchRequestEndpointBundleValueEthPortPathsInner object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBundlesPatchRequestEndpointBundleValueEthPortPathsInnerWithDefaults() *BundlesPatchRequestEndpointBundleValueEthPortPathsInner {
	this := BundlesPatchRequestEndpointBundleValueEthPortPathsInner{}
	var ethPortNumEthPortProfile string = ""
	this.EthPortNumEthPortProfile = &ethPortNumEthPortProfile
	var ethPortNumEthPortSettings string = ""
	this.EthPortNumEthPortSettings = &ethPortNumEthPortSettings
	var ethPortNumGatewayProfile string = ""
	this.EthPortNumGatewayProfile = &ethPortNumGatewayProfile
	return &this
}

// GetEthPortNumEthPortProfile returns the EthPortNumEthPortProfile field value if set, zero value otherwise.
func (o *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) GetEthPortNumEthPortProfile() string {
	if o == nil || IsNil(o.EthPortNumEthPortProfile) {
		var ret string
		return ret
	}
	return *o.EthPortNumEthPortProfile
}

// GetEthPortNumEthPortProfileOk returns a tuple with the EthPortNumEthPortProfile field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) GetEthPortNumEthPortProfileOk() (*string, bool) {
	if o == nil || IsNil(o.EthPortNumEthPortProfile) {
		return nil, false
	}
	return o.EthPortNumEthPortProfile, true
}

// HasEthPortNumEthPortProfile returns a boolean if a field has been set.
func (o *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) HasEthPortNumEthPortProfile() bool {
	if o != nil && !IsNil(o.EthPortNumEthPortProfile) {
		return true
	}

	return false
}

// SetEthPortNumEthPortProfile gets a reference to the given string and assigns it to the EthPortNumEthPortProfile field.
func (o *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) SetEthPortNumEthPortProfile(v string) {
	o.EthPortNumEthPortProfile = &v
}

// GetEthPortNumEthPortProfileRefType returns the EthPortNumEthPortProfileRefType field value if set, zero value otherwise.
func (o *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) GetEthPortNumEthPortProfileRefType() string {
	if o == nil || IsNil(o.EthPortNumEthPortProfileRefType) {
		var ret string
		return ret
	}
	return *o.EthPortNumEthPortProfileRefType
}

// GetEthPortNumEthPortProfileRefTypeOk returns a tuple with the EthPortNumEthPortProfileRefType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) GetEthPortNumEthPortProfileRefTypeOk() (*string, bool) {
	if o == nil || IsNil(o.EthPortNumEthPortProfileRefType) {
		return nil, false
	}
	return o.EthPortNumEthPortProfileRefType, true
}

// HasEthPortNumEthPortProfileRefType returns a boolean if a field has been set.
func (o *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) HasEthPortNumEthPortProfileRefType() bool {
	if o != nil && !IsNil(o.EthPortNumEthPortProfileRefType) {
		return true
	}

	return false
}

// SetEthPortNumEthPortProfileRefType gets a reference to the given string and assigns it to the EthPortNumEthPortProfileRefType field.
func (o *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) SetEthPortNumEthPortProfileRefType(v string) {
	o.EthPortNumEthPortProfileRefType = &v
}

// GetEthPortNumEthPortSettings returns the EthPortNumEthPortSettings field value if set, zero value otherwise.
func (o *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) GetEthPortNumEthPortSettings() string {
	if o == nil || IsNil(o.EthPortNumEthPortSettings) {
		var ret string
		return ret
	}
	return *o.EthPortNumEthPortSettings
}

// GetEthPortNumEthPortSettingsOk returns a tuple with the EthPortNumEthPortSettings field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) GetEthPortNumEthPortSettingsOk() (*string, bool) {
	if o == nil || IsNil(o.EthPortNumEthPortSettings) {
		return nil, false
	}
	return o.EthPortNumEthPortSettings, true
}

// HasEthPortNumEthPortSettings returns a boolean if a field has been set.
func (o *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) HasEthPortNumEthPortSettings() bool {
	if o != nil && !IsNil(o.EthPortNumEthPortSettings) {
		return true
	}

	return false
}

// SetEthPortNumEthPortSettings gets a reference to the given string and assigns it to the EthPortNumEthPortSettings field.
func (o *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) SetEthPortNumEthPortSettings(v string) {
	o.EthPortNumEthPortSettings = &v
}

// GetEthPortNumEthPortSettingsRefType returns the EthPortNumEthPortSettingsRefType field value if set, zero value otherwise.
func (o *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) GetEthPortNumEthPortSettingsRefType() string {
	if o == nil || IsNil(o.EthPortNumEthPortSettingsRefType) {
		var ret string
		return ret
	}
	return *o.EthPortNumEthPortSettingsRefType
}

// GetEthPortNumEthPortSettingsRefTypeOk returns a tuple with the EthPortNumEthPortSettingsRefType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) GetEthPortNumEthPortSettingsRefTypeOk() (*string, bool) {
	if o == nil || IsNil(o.EthPortNumEthPortSettingsRefType) {
		return nil, false
	}
	return o.EthPortNumEthPortSettingsRefType, true
}

// HasEthPortNumEthPortSettingsRefType returns a boolean if a field has been set.
func (o *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) HasEthPortNumEthPortSettingsRefType() bool {
	if o != nil && !IsNil(o.EthPortNumEthPortSettingsRefType) {
		return true
	}

	return false
}

// SetEthPortNumEthPortSettingsRefType gets a reference to the given string and assigns it to the EthPortNumEthPortSettingsRefType field.
func (o *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) SetEthPortNumEthPortSettingsRefType(v string) {
	o.EthPortNumEthPortSettingsRefType = &v
}

// GetEthPortNumGatewayProfile returns the EthPortNumGatewayProfile field value if set, zero value otherwise.
func (o *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) GetEthPortNumGatewayProfile() string {
	if o == nil || IsNil(o.EthPortNumGatewayProfile) {
		var ret string
		return ret
	}
	return *o.EthPortNumGatewayProfile
}

// GetEthPortNumGatewayProfileOk returns a tuple with the EthPortNumGatewayProfile field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) GetEthPortNumGatewayProfileOk() (*string, bool) {
	if o == nil || IsNil(o.EthPortNumGatewayProfile) {
		return nil, false
	}
	return o.EthPortNumGatewayProfile, true
}

// HasEthPortNumGatewayProfile returns a boolean if a field has been set.
func (o *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) HasEthPortNumGatewayProfile() bool {
	if o != nil && !IsNil(o.EthPortNumGatewayProfile) {
		return true
	}

	return false
}

// SetEthPortNumGatewayProfile gets a reference to the given string and assigns it to the EthPortNumGatewayProfile field.
func (o *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) SetEthPortNumGatewayProfile(v string) {
	o.EthPortNumGatewayProfile = &v
}

// GetEthPortNumGatewayProfileRefType returns the EthPortNumGatewayProfileRefType field value if set, zero value otherwise.
func (o *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) GetEthPortNumGatewayProfileRefType() string {
	if o == nil || IsNil(o.EthPortNumGatewayProfileRefType) {
		var ret string
		return ret
	}
	return *o.EthPortNumGatewayProfileRefType
}

// GetEthPortNumGatewayProfileRefTypeOk returns a tuple with the EthPortNumGatewayProfileRefType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) GetEthPortNumGatewayProfileRefTypeOk() (*string, bool) {
	if o == nil || IsNil(o.EthPortNumGatewayProfileRefType) {
		return nil, false
	}
	return o.EthPortNumGatewayProfileRefType, true
}

// HasEthPortNumGatewayProfileRefType returns a boolean if a field has been set.
func (o *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) HasEthPortNumGatewayProfileRefType() bool {
	if o != nil && !IsNil(o.EthPortNumGatewayProfileRefType) {
		return true
	}

	return false
}

// SetEthPortNumGatewayProfileRefType gets a reference to the given string and assigns it to the EthPortNumGatewayProfileRefType field.
func (o *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) SetEthPortNumGatewayProfileRefType(v string) {
	o.EthPortNumGatewayProfileRefType = &v
}

// GetIndex returns the Index field value if set, zero value otherwise.
func (o *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) GetIndex() int32 {
	if o == nil || IsNil(o.Index) {
		var ret int32
		return ret
	}
	return *o.Index
}

// GetIndexOk returns a tuple with the Index field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) GetIndexOk() (*int32, bool) {
	if o == nil || IsNil(o.Index) {
		return nil, false
	}
	return o.Index, true
}

// HasIndex returns a boolean if a field has been set.
func (o *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) HasIndex() bool {
	if o != nil && !IsNil(o.Index) {
		return true
	}

	return false
}

// SetIndex gets a reference to the given int32 and assigns it to the Index field.
func (o *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) SetIndex(v int32) {
	o.Index = &v
}

// GetPortName returns the PortName field value if set, zero value otherwise.
func (o *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) GetPortName() string {
	if o == nil || IsNil(o.PortName) {
		var ret string
		return ret
	}
	return *o.PortName
}

// GetPortNameOk returns a tuple with the PortName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) GetPortNameOk() (*string, bool) {
	if o == nil || IsNil(o.PortName) {
		return nil, false
	}
	return o.PortName, true
}

// HasPortName returns a boolean if a field has been set.
func (o *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) HasPortName() bool {
	if o != nil && !IsNil(o.PortName) {
		return true
	}

	return false
}

// SetPortName gets a reference to the given string and assigns it to the PortName field.
func (o *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) SetPortName(v string) {
	o.PortName = &v
}

func (o BundlesPatchRequestEndpointBundleValueEthPortPathsInner) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o BundlesPatchRequestEndpointBundleValueEthPortPathsInner) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.EthPortNumEthPortProfile) {
		toSerialize["eth_port_num_eth_port_profile"] = o.EthPortNumEthPortProfile
	}
	if !IsNil(o.EthPortNumEthPortProfileRefType) {
		toSerialize["eth_port_num_eth_port_profile_ref_type_"] = o.EthPortNumEthPortProfileRefType
	}
	if !IsNil(o.EthPortNumEthPortSettings) {
		toSerialize["eth_port_num_eth_port_settings"] = o.EthPortNumEthPortSettings
	}
	if !IsNil(o.EthPortNumEthPortSettingsRefType) {
		toSerialize["eth_port_num_eth_port_settings_ref_type_"] = o.EthPortNumEthPortSettingsRefType
	}
	if !IsNil(o.EthPortNumGatewayProfile) {
		toSerialize["eth_port_num_gateway_profile"] = o.EthPortNumGatewayProfile
	}
	if !IsNil(o.EthPortNumGatewayProfileRefType) {
		toSerialize["eth_port_num_gateway_profile_ref_type_"] = o.EthPortNumGatewayProfileRefType
	}
	if !IsNil(o.Index) {
		toSerialize["index"] = o.Index
	}
	if !IsNil(o.PortName) {
		toSerialize["port_name"] = o.PortName
	}
	return toSerialize, nil
}

type NullableBundlesPatchRequestEndpointBundleValueEthPortPathsInner struct {
	value *BundlesPatchRequestEndpointBundleValueEthPortPathsInner
	isSet bool
}

func (v NullableBundlesPatchRequestEndpointBundleValueEthPortPathsInner) Get() *BundlesPatchRequestEndpointBundleValueEthPortPathsInner {
	return v.value
}

func (v *NullableBundlesPatchRequestEndpointBundleValueEthPortPathsInner) Set(val *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) {
	v.value = val
	v.isSet = true
}

func (v NullableBundlesPatchRequestEndpointBundleValueEthPortPathsInner) IsSet() bool {
	return v.isSet
}

func (v *NullableBundlesPatchRequestEndpointBundleValueEthPortPathsInner) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBundlesPatchRequestEndpointBundleValueEthPortPathsInner(val *BundlesPatchRequestEndpointBundleValueEthPortPathsInner) *NullableBundlesPatchRequestEndpointBundleValueEthPortPathsInner {
	return &NullableBundlesPatchRequestEndpointBundleValueEthPortPathsInner{value: val, isSet: true}
}

func (v NullableBundlesPatchRequestEndpointBundleValueEthPortPathsInner) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBundlesPatchRequestEndpointBundleValueEthPortPathsInner) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


