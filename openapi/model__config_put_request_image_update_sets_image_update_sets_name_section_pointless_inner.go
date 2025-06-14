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

// checks if the ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner{}

// ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner struct for ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner
type ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner struct {
	// The name of the Endpoint Set
	EndpointSetNumName *string `json:"endpoint_set_num_name,omitempty"`
	// The target SW version for member devices of the Endpoint Set
	EndpointSetForEndpointlessTargetUpgradeVersion *string `json:"endpoint_set_for_endpointless_target_upgrade_version,omitempty"`
	// Unique Identifier - not editable
	EndpointSetForEndpointlessUniqueIdentifier *string `json:"endpoint_set_for_endpointless_unique_identifier,omitempty"`
	// Include on the Summary
	EndpointSetNumOnSummary *bool `json:"endpoint_set_num_on_summary,omitempty"`
	// The time to update to the target SW version
	EndpointSetForEndpointlessTargetUpgradeVersionTime *string `json:"endpoint_set_for_endpointless_target_upgrade_version_time,omitempty"`
}

// NewConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner instantiates a new ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner() *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner {
	this := ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner{}
	var endpointSetNumName string = "Unassigned Devices"
	this.EndpointSetNumName = &endpointSetNumName
	var endpointSetForEndpointlessTargetUpgradeVersion string = "unmanaged"
	this.EndpointSetForEndpointlessTargetUpgradeVersion = &endpointSetForEndpointlessTargetUpgradeVersion
	var endpointSetForEndpointlessUniqueIdentifier string = "pointless"
	this.EndpointSetForEndpointlessUniqueIdentifier = &endpointSetForEndpointlessUniqueIdentifier
	var endpointSetNumOnSummary bool = true
	this.EndpointSetNumOnSummary = &endpointSetNumOnSummary
	var endpointSetForEndpointlessTargetUpgradeVersionTime string = ""
	this.EndpointSetForEndpointlessTargetUpgradeVersionTime = &endpointSetForEndpointlessTargetUpgradeVersionTime
	return &this
}

// NewConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInnerWithDefaults instantiates a new ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInnerWithDefaults() *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner {
	this := ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner{}
	var endpointSetNumName string = "Unassigned Devices"
	this.EndpointSetNumName = &endpointSetNumName
	var endpointSetForEndpointlessTargetUpgradeVersion string = "unmanaged"
	this.EndpointSetForEndpointlessTargetUpgradeVersion = &endpointSetForEndpointlessTargetUpgradeVersion
	var endpointSetForEndpointlessUniqueIdentifier string = "pointless"
	this.EndpointSetForEndpointlessUniqueIdentifier = &endpointSetForEndpointlessUniqueIdentifier
	var endpointSetNumOnSummary bool = true
	this.EndpointSetNumOnSummary = &endpointSetNumOnSummary
	var endpointSetForEndpointlessTargetUpgradeVersionTime string = ""
	this.EndpointSetForEndpointlessTargetUpgradeVersionTime = &endpointSetForEndpointlessTargetUpgradeVersionTime
	return &this
}

// GetEndpointSetNumName returns the EndpointSetNumName field value if set, zero value otherwise.
func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner) GetEndpointSetNumName() string {
	if o == nil || IsNil(o.EndpointSetNumName) {
		var ret string
		return ret
	}
	return *o.EndpointSetNumName
}

// GetEndpointSetNumNameOk returns a tuple with the EndpointSetNumName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner) GetEndpointSetNumNameOk() (*string, bool) {
	if o == nil || IsNil(o.EndpointSetNumName) {
		return nil, false
	}
	return o.EndpointSetNumName, true
}

// HasEndpointSetNumName returns a boolean if a field has been set.
func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner) HasEndpointSetNumName() bool {
	if o != nil && !IsNil(o.EndpointSetNumName) {
		return true
	}

	return false
}

// SetEndpointSetNumName gets a reference to the given string and assigns it to the EndpointSetNumName field.
func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner) SetEndpointSetNumName(v string) {
	o.EndpointSetNumName = &v
}

// GetEndpointSetForEndpointlessTargetUpgradeVersion returns the EndpointSetForEndpointlessTargetUpgradeVersion field value if set, zero value otherwise.
func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner) GetEndpointSetForEndpointlessTargetUpgradeVersion() string {
	if o == nil || IsNil(o.EndpointSetForEndpointlessTargetUpgradeVersion) {
		var ret string
		return ret
	}
	return *o.EndpointSetForEndpointlessTargetUpgradeVersion
}

// GetEndpointSetForEndpointlessTargetUpgradeVersionOk returns a tuple with the EndpointSetForEndpointlessTargetUpgradeVersion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner) GetEndpointSetForEndpointlessTargetUpgradeVersionOk() (*string, bool) {
	if o == nil || IsNil(o.EndpointSetForEndpointlessTargetUpgradeVersion) {
		return nil, false
	}
	return o.EndpointSetForEndpointlessTargetUpgradeVersion, true
}

// HasEndpointSetForEndpointlessTargetUpgradeVersion returns a boolean if a field has been set.
func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner) HasEndpointSetForEndpointlessTargetUpgradeVersion() bool {
	if o != nil && !IsNil(o.EndpointSetForEndpointlessTargetUpgradeVersion) {
		return true
	}

	return false
}

// SetEndpointSetForEndpointlessTargetUpgradeVersion gets a reference to the given string and assigns it to the EndpointSetForEndpointlessTargetUpgradeVersion field.
func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner) SetEndpointSetForEndpointlessTargetUpgradeVersion(v string) {
	o.EndpointSetForEndpointlessTargetUpgradeVersion = &v
}

// GetEndpointSetForEndpointlessUniqueIdentifier returns the EndpointSetForEndpointlessUniqueIdentifier field value if set, zero value otherwise.
func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner) GetEndpointSetForEndpointlessUniqueIdentifier() string {
	if o == nil || IsNil(o.EndpointSetForEndpointlessUniqueIdentifier) {
		var ret string
		return ret
	}
	return *o.EndpointSetForEndpointlessUniqueIdentifier
}

// GetEndpointSetForEndpointlessUniqueIdentifierOk returns a tuple with the EndpointSetForEndpointlessUniqueIdentifier field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner) GetEndpointSetForEndpointlessUniqueIdentifierOk() (*string, bool) {
	if o == nil || IsNil(o.EndpointSetForEndpointlessUniqueIdentifier) {
		return nil, false
	}
	return o.EndpointSetForEndpointlessUniqueIdentifier, true
}

// HasEndpointSetForEndpointlessUniqueIdentifier returns a boolean if a field has been set.
func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner) HasEndpointSetForEndpointlessUniqueIdentifier() bool {
	if o != nil && !IsNil(o.EndpointSetForEndpointlessUniqueIdentifier) {
		return true
	}

	return false
}

// SetEndpointSetForEndpointlessUniqueIdentifier gets a reference to the given string and assigns it to the EndpointSetForEndpointlessUniqueIdentifier field.
func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner) SetEndpointSetForEndpointlessUniqueIdentifier(v string) {
	o.EndpointSetForEndpointlessUniqueIdentifier = &v
}

// GetEndpointSetNumOnSummary returns the EndpointSetNumOnSummary field value if set, zero value otherwise.
func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner) GetEndpointSetNumOnSummary() bool {
	if o == nil || IsNil(o.EndpointSetNumOnSummary) {
		var ret bool
		return ret
	}
	return *o.EndpointSetNumOnSummary
}

// GetEndpointSetNumOnSummaryOk returns a tuple with the EndpointSetNumOnSummary field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner) GetEndpointSetNumOnSummaryOk() (*bool, bool) {
	if o == nil || IsNil(o.EndpointSetNumOnSummary) {
		return nil, false
	}
	return o.EndpointSetNumOnSummary, true
}

// HasEndpointSetNumOnSummary returns a boolean if a field has been set.
func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner) HasEndpointSetNumOnSummary() bool {
	if o != nil && !IsNil(o.EndpointSetNumOnSummary) {
		return true
	}

	return false
}

// SetEndpointSetNumOnSummary gets a reference to the given bool and assigns it to the EndpointSetNumOnSummary field.
func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner) SetEndpointSetNumOnSummary(v bool) {
	o.EndpointSetNumOnSummary = &v
}

// GetEndpointSetForEndpointlessTargetUpgradeVersionTime returns the EndpointSetForEndpointlessTargetUpgradeVersionTime field value if set, zero value otherwise.
func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner) GetEndpointSetForEndpointlessTargetUpgradeVersionTime() string {
	if o == nil || IsNil(o.EndpointSetForEndpointlessTargetUpgradeVersionTime) {
		var ret string
		return ret
	}
	return *o.EndpointSetForEndpointlessTargetUpgradeVersionTime
}

// GetEndpointSetForEndpointlessTargetUpgradeVersionTimeOk returns a tuple with the EndpointSetForEndpointlessTargetUpgradeVersionTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner) GetEndpointSetForEndpointlessTargetUpgradeVersionTimeOk() (*string, bool) {
	if o == nil || IsNil(o.EndpointSetForEndpointlessTargetUpgradeVersionTime) {
		return nil, false
	}
	return o.EndpointSetForEndpointlessTargetUpgradeVersionTime, true
}

// HasEndpointSetForEndpointlessTargetUpgradeVersionTime returns a boolean if a field has been set.
func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner) HasEndpointSetForEndpointlessTargetUpgradeVersionTime() bool {
	if o != nil && !IsNil(o.EndpointSetForEndpointlessTargetUpgradeVersionTime) {
		return true
	}

	return false
}

// SetEndpointSetForEndpointlessTargetUpgradeVersionTime gets a reference to the given string and assigns it to the EndpointSetForEndpointlessTargetUpgradeVersionTime field.
func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner) SetEndpointSetForEndpointlessTargetUpgradeVersionTime(v string) {
	o.EndpointSetForEndpointlessTargetUpgradeVersionTime = &v
}

func (o ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.EndpointSetNumName) {
		toSerialize["endpoint_set_num_name"] = o.EndpointSetNumName
	}
	if !IsNil(o.EndpointSetForEndpointlessTargetUpgradeVersion) {
		toSerialize["endpoint_set_for_endpointless_target_upgrade_version"] = o.EndpointSetForEndpointlessTargetUpgradeVersion
	}
	if !IsNil(o.EndpointSetForEndpointlessUniqueIdentifier) {
		toSerialize["endpoint_set_for_endpointless_unique_identifier"] = o.EndpointSetForEndpointlessUniqueIdentifier
	}
	if !IsNil(o.EndpointSetNumOnSummary) {
		toSerialize["endpoint_set_num_on_summary"] = o.EndpointSetNumOnSummary
	}
	if !IsNil(o.EndpointSetForEndpointlessTargetUpgradeVersionTime) {
		toSerialize["endpoint_set_for_endpointless_target_upgrade_version_time"] = o.EndpointSetForEndpointlessTargetUpgradeVersionTime
	}
	return toSerialize, nil
}

type NullableConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner struct {
	value *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner
	isSet bool
}

func (v NullableConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner) Get() *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner {
	return v.value
}

func (v *NullableConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner) Set(val *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner) {
	v.value = val
	v.isSet = true
}

func (v NullableConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner) IsSet() bool {
	return v.isSet
}

func (v *NullableConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner(val *ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner) *NullableConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner {
	return &NullableConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner{value: val, isSet: true}
}

func (v NullableConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


