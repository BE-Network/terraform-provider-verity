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

// checks if the ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner{}

// ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner struct for ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner
type ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner struct {
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

// NewImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner instantiates a new ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner() *ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner {
	this := ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner{}
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

// NewImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInnerWithDefaults instantiates a new ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInnerWithDefaults() *ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner {
	this := ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner{}
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
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner) GetEndpointSetNumName() string {
	if o == nil || IsNil(o.EndpointSetNumName) {
		var ret string
		return ret
	}
	return *o.EndpointSetNumName
}

// GetEndpointSetNumNameOk returns a tuple with the EndpointSetNumName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner) GetEndpointSetNumNameOk() (*string, bool) {
	if o == nil || IsNil(o.EndpointSetNumName) {
		return nil, false
	}
	return o.EndpointSetNumName, true
}

// HasEndpointSetNumName returns a boolean if a field has been set.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner) HasEndpointSetNumName() bool {
	if o != nil && !IsNil(o.EndpointSetNumName) {
		return true
	}

	return false
}

// SetEndpointSetNumName gets a reference to the given string and assigns it to the EndpointSetNumName field.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner) SetEndpointSetNumName(v string) {
	o.EndpointSetNumName = &v
}

// GetEndpointSetForEndpointlessTargetUpgradeVersion returns the EndpointSetForEndpointlessTargetUpgradeVersion field value if set, zero value otherwise.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner) GetEndpointSetForEndpointlessTargetUpgradeVersion() string {
	if o == nil || IsNil(o.EndpointSetForEndpointlessTargetUpgradeVersion) {
		var ret string
		return ret
	}
	return *o.EndpointSetForEndpointlessTargetUpgradeVersion
}

// GetEndpointSetForEndpointlessTargetUpgradeVersionOk returns a tuple with the EndpointSetForEndpointlessTargetUpgradeVersion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner) GetEndpointSetForEndpointlessTargetUpgradeVersionOk() (*string, bool) {
	if o == nil || IsNil(o.EndpointSetForEndpointlessTargetUpgradeVersion) {
		return nil, false
	}
	return o.EndpointSetForEndpointlessTargetUpgradeVersion, true
}

// HasEndpointSetForEndpointlessTargetUpgradeVersion returns a boolean if a field has been set.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner) HasEndpointSetForEndpointlessTargetUpgradeVersion() bool {
	if o != nil && !IsNil(o.EndpointSetForEndpointlessTargetUpgradeVersion) {
		return true
	}

	return false
}

// SetEndpointSetForEndpointlessTargetUpgradeVersion gets a reference to the given string and assigns it to the EndpointSetForEndpointlessTargetUpgradeVersion field.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner) SetEndpointSetForEndpointlessTargetUpgradeVersion(v string) {
	o.EndpointSetForEndpointlessTargetUpgradeVersion = &v
}

// GetEndpointSetForEndpointlessUniqueIdentifier returns the EndpointSetForEndpointlessUniqueIdentifier field value if set, zero value otherwise.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner) GetEndpointSetForEndpointlessUniqueIdentifier() string {
	if o == nil || IsNil(o.EndpointSetForEndpointlessUniqueIdentifier) {
		var ret string
		return ret
	}
	return *o.EndpointSetForEndpointlessUniqueIdentifier
}

// GetEndpointSetForEndpointlessUniqueIdentifierOk returns a tuple with the EndpointSetForEndpointlessUniqueIdentifier field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner) GetEndpointSetForEndpointlessUniqueIdentifierOk() (*string, bool) {
	if o == nil || IsNil(o.EndpointSetForEndpointlessUniqueIdentifier) {
		return nil, false
	}
	return o.EndpointSetForEndpointlessUniqueIdentifier, true
}

// HasEndpointSetForEndpointlessUniqueIdentifier returns a boolean if a field has been set.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner) HasEndpointSetForEndpointlessUniqueIdentifier() bool {
	if o != nil && !IsNil(o.EndpointSetForEndpointlessUniqueIdentifier) {
		return true
	}

	return false
}

// SetEndpointSetForEndpointlessUniqueIdentifier gets a reference to the given string and assigns it to the EndpointSetForEndpointlessUniqueIdentifier field.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner) SetEndpointSetForEndpointlessUniqueIdentifier(v string) {
	o.EndpointSetForEndpointlessUniqueIdentifier = &v
}

// GetEndpointSetNumOnSummary returns the EndpointSetNumOnSummary field value if set, zero value otherwise.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner) GetEndpointSetNumOnSummary() bool {
	if o == nil || IsNil(o.EndpointSetNumOnSummary) {
		var ret bool
		return ret
	}
	return *o.EndpointSetNumOnSummary
}

// GetEndpointSetNumOnSummaryOk returns a tuple with the EndpointSetNumOnSummary field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner) GetEndpointSetNumOnSummaryOk() (*bool, bool) {
	if o == nil || IsNil(o.EndpointSetNumOnSummary) {
		return nil, false
	}
	return o.EndpointSetNumOnSummary, true
}

// HasEndpointSetNumOnSummary returns a boolean if a field has been set.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner) HasEndpointSetNumOnSummary() bool {
	if o != nil && !IsNil(o.EndpointSetNumOnSummary) {
		return true
	}

	return false
}

// SetEndpointSetNumOnSummary gets a reference to the given bool and assigns it to the EndpointSetNumOnSummary field.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner) SetEndpointSetNumOnSummary(v bool) {
	o.EndpointSetNumOnSummary = &v
}

// GetEndpointSetForEndpointlessTargetUpgradeVersionTime returns the EndpointSetForEndpointlessTargetUpgradeVersionTime field value if set, zero value otherwise.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner) GetEndpointSetForEndpointlessTargetUpgradeVersionTime() string {
	if o == nil || IsNil(o.EndpointSetForEndpointlessTargetUpgradeVersionTime) {
		var ret string
		return ret
	}
	return *o.EndpointSetForEndpointlessTargetUpgradeVersionTime
}

// GetEndpointSetForEndpointlessTargetUpgradeVersionTimeOk returns a tuple with the EndpointSetForEndpointlessTargetUpgradeVersionTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner) GetEndpointSetForEndpointlessTargetUpgradeVersionTimeOk() (*string, bool) {
	if o == nil || IsNil(o.EndpointSetForEndpointlessTargetUpgradeVersionTime) {
		return nil, false
	}
	return o.EndpointSetForEndpointlessTargetUpgradeVersionTime, true
}

// HasEndpointSetForEndpointlessTargetUpgradeVersionTime returns a boolean if a field has been set.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner) HasEndpointSetForEndpointlessTargetUpgradeVersionTime() bool {
	if o != nil && !IsNil(o.EndpointSetForEndpointlessTargetUpgradeVersionTime) {
		return true
	}

	return false
}

// SetEndpointSetForEndpointlessTargetUpgradeVersionTime gets a reference to the given string and assigns it to the EndpointSetForEndpointlessTargetUpgradeVersionTime field.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner) SetEndpointSetForEndpointlessTargetUpgradeVersionTime(v string) {
	o.EndpointSetForEndpointlessTargetUpgradeVersionTime = &v
}

func (o ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner) ToMap() (map[string]interface{}, error) {
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

type NullableImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner struct {
	value *ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner
	isSet bool
}

func (v NullableImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner) Get() *ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner {
	return v.value
}

func (v *NullableImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner) Set(val *ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner) {
	v.value = val
	v.isSet = true
}

func (v NullableImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner) IsSet() bool {
	return v.isSet
}

func (v *NullableImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner(val *ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner) *NullableImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner {
	return &NullableImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner{value: val, isSet: true}
}

func (v NullableImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


