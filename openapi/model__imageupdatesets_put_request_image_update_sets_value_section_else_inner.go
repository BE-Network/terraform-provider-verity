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

// checks if the ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner{}

// ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner struct for ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner
type ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner struct {
	// The name of the Endpoint Set
	EndpointSetNumName *string `json:"endpoint_set_num_name,omitempty"`
	// The target SW version for member devices of the Endpoint Set
	EndpointSetForAllOthersTargetUpgradeVersion *string `json:"endpoint_set_for_all_others_target_upgrade_version,omitempty"`
	// Unique Identifier - not editable
	EndpointSetForAllOthersUniqueIdentifier *string `json:"endpoint_set_for_all_others_unique_identifier,omitempty"`
	// Include on the Summary
	EndpointSetNumOnSummary *bool `json:"endpoint_set_num_on_summary,omitempty"`
	// The time to update to the target SW version
	EndpointSetForAllOthersTargetUpgradeVersionTime *string `json:"endpoint_set_for_all_others_target_upgrade_version_time,omitempty"`
}

// NewImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner instantiates a new ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner() *ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner {
	this := ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner{}
	var endpointSetNumName string = "All Others"
	this.EndpointSetNumName = &endpointSetNumName
	var endpointSetForAllOthersTargetUpgradeVersion string = "unmanaged"
	this.EndpointSetForAllOthersTargetUpgradeVersion = &endpointSetForAllOthersTargetUpgradeVersion
	var endpointSetForAllOthersUniqueIdentifier string = "else"
	this.EndpointSetForAllOthersUniqueIdentifier = &endpointSetForAllOthersUniqueIdentifier
	var endpointSetNumOnSummary bool = true
	this.EndpointSetNumOnSummary = &endpointSetNumOnSummary
	var endpointSetForAllOthersTargetUpgradeVersionTime string = ""
	this.EndpointSetForAllOthersTargetUpgradeVersionTime = &endpointSetForAllOthersTargetUpgradeVersionTime
	return &this
}

// NewImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInnerWithDefaults instantiates a new ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInnerWithDefaults() *ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner {
	this := ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner{}
	var endpointSetNumName string = "All Others"
	this.EndpointSetNumName = &endpointSetNumName
	var endpointSetForAllOthersTargetUpgradeVersion string = "unmanaged"
	this.EndpointSetForAllOthersTargetUpgradeVersion = &endpointSetForAllOthersTargetUpgradeVersion
	var endpointSetForAllOthersUniqueIdentifier string = "else"
	this.EndpointSetForAllOthersUniqueIdentifier = &endpointSetForAllOthersUniqueIdentifier
	var endpointSetNumOnSummary bool = true
	this.EndpointSetNumOnSummary = &endpointSetNumOnSummary
	var endpointSetForAllOthersTargetUpgradeVersionTime string = ""
	this.EndpointSetForAllOthersTargetUpgradeVersionTime = &endpointSetForAllOthersTargetUpgradeVersionTime
	return &this
}

// GetEndpointSetNumName returns the EndpointSetNumName field value if set, zero value otherwise.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner) GetEndpointSetNumName() string {
	if o == nil || IsNil(o.EndpointSetNumName) {
		var ret string
		return ret
	}
	return *o.EndpointSetNumName
}

// GetEndpointSetNumNameOk returns a tuple with the EndpointSetNumName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner) GetEndpointSetNumNameOk() (*string, bool) {
	if o == nil || IsNil(o.EndpointSetNumName) {
		return nil, false
	}
	return o.EndpointSetNumName, true
}

// HasEndpointSetNumName returns a boolean if a field has been set.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner) HasEndpointSetNumName() bool {
	if o != nil && !IsNil(o.EndpointSetNumName) {
		return true
	}

	return false
}

// SetEndpointSetNumName gets a reference to the given string and assigns it to the EndpointSetNumName field.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner) SetEndpointSetNumName(v string) {
	o.EndpointSetNumName = &v
}

// GetEndpointSetForAllOthersTargetUpgradeVersion returns the EndpointSetForAllOthersTargetUpgradeVersion field value if set, zero value otherwise.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner) GetEndpointSetForAllOthersTargetUpgradeVersion() string {
	if o == nil || IsNil(o.EndpointSetForAllOthersTargetUpgradeVersion) {
		var ret string
		return ret
	}
	return *o.EndpointSetForAllOthersTargetUpgradeVersion
}

// GetEndpointSetForAllOthersTargetUpgradeVersionOk returns a tuple with the EndpointSetForAllOthersTargetUpgradeVersion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner) GetEndpointSetForAllOthersTargetUpgradeVersionOk() (*string, bool) {
	if o == nil || IsNil(o.EndpointSetForAllOthersTargetUpgradeVersion) {
		return nil, false
	}
	return o.EndpointSetForAllOthersTargetUpgradeVersion, true
}

// HasEndpointSetForAllOthersTargetUpgradeVersion returns a boolean if a field has been set.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner) HasEndpointSetForAllOthersTargetUpgradeVersion() bool {
	if o != nil && !IsNil(o.EndpointSetForAllOthersTargetUpgradeVersion) {
		return true
	}

	return false
}

// SetEndpointSetForAllOthersTargetUpgradeVersion gets a reference to the given string and assigns it to the EndpointSetForAllOthersTargetUpgradeVersion field.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner) SetEndpointSetForAllOthersTargetUpgradeVersion(v string) {
	o.EndpointSetForAllOthersTargetUpgradeVersion = &v
}

// GetEndpointSetForAllOthersUniqueIdentifier returns the EndpointSetForAllOthersUniqueIdentifier field value if set, zero value otherwise.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner) GetEndpointSetForAllOthersUniqueIdentifier() string {
	if o == nil || IsNil(o.EndpointSetForAllOthersUniqueIdentifier) {
		var ret string
		return ret
	}
	return *o.EndpointSetForAllOthersUniqueIdentifier
}

// GetEndpointSetForAllOthersUniqueIdentifierOk returns a tuple with the EndpointSetForAllOthersUniqueIdentifier field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner) GetEndpointSetForAllOthersUniqueIdentifierOk() (*string, bool) {
	if o == nil || IsNil(o.EndpointSetForAllOthersUniqueIdentifier) {
		return nil, false
	}
	return o.EndpointSetForAllOthersUniqueIdentifier, true
}

// HasEndpointSetForAllOthersUniqueIdentifier returns a boolean if a field has been set.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner) HasEndpointSetForAllOthersUniqueIdentifier() bool {
	if o != nil && !IsNil(o.EndpointSetForAllOthersUniqueIdentifier) {
		return true
	}

	return false
}

// SetEndpointSetForAllOthersUniqueIdentifier gets a reference to the given string and assigns it to the EndpointSetForAllOthersUniqueIdentifier field.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner) SetEndpointSetForAllOthersUniqueIdentifier(v string) {
	o.EndpointSetForAllOthersUniqueIdentifier = &v
}

// GetEndpointSetNumOnSummary returns the EndpointSetNumOnSummary field value if set, zero value otherwise.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner) GetEndpointSetNumOnSummary() bool {
	if o == nil || IsNil(o.EndpointSetNumOnSummary) {
		var ret bool
		return ret
	}
	return *o.EndpointSetNumOnSummary
}

// GetEndpointSetNumOnSummaryOk returns a tuple with the EndpointSetNumOnSummary field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner) GetEndpointSetNumOnSummaryOk() (*bool, bool) {
	if o == nil || IsNil(o.EndpointSetNumOnSummary) {
		return nil, false
	}
	return o.EndpointSetNumOnSummary, true
}

// HasEndpointSetNumOnSummary returns a boolean if a field has been set.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner) HasEndpointSetNumOnSummary() bool {
	if o != nil && !IsNil(o.EndpointSetNumOnSummary) {
		return true
	}

	return false
}

// SetEndpointSetNumOnSummary gets a reference to the given bool and assigns it to the EndpointSetNumOnSummary field.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner) SetEndpointSetNumOnSummary(v bool) {
	o.EndpointSetNumOnSummary = &v
}

// GetEndpointSetForAllOthersTargetUpgradeVersionTime returns the EndpointSetForAllOthersTargetUpgradeVersionTime field value if set, zero value otherwise.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner) GetEndpointSetForAllOthersTargetUpgradeVersionTime() string {
	if o == nil || IsNil(o.EndpointSetForAllOthersTargetUpgradeVersionTime) {
		var ret string
		return ret
	}
	return *o.EndpointSetForAllOthersTargetUpgradeVersionTime
}

// GetEndpointSetForAllOthersTargetUpgradeVersionTimeOk returns a tuple with the EndpointSetForAllOthersTargetUpgradeVersionTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner) GetEndpointSetForAllOthersTargetUpgradeVersionTimeOk() (*string, bool) {
	if o == nil || IsNil(o.EndpointSetForAllOthersTargetUpgradeVersionTime) {
		return nil, false
	}
	return o.EndpointSetForAllOthersTargetUpgradeVersionTime, true
}

// HasEndpointSetForAllOthersTargetUpgradeVersionTime returns a boolean if a field has been set.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner) HasEndpointSetForAllOthersTargetUpgradeVersionTime() bool {
	if o != nil && !IsNil(o.EndpointSetForAllOthersTargetUpgradeVersionTime) {
		return true
	}

	return false
}

// SetEndpointSetForAllOthersTargetUpgradeVersionTime gets a reference to the given string and assigns it to the EndpointSetForAllOthersTargetUpgradeVersionTime field.
func (o *ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner) SetEndpointSetForAllOthersTargetUpgradeVersionTime(v string) {
	o.EndpointSetForAllOthersTargetUpgradeVersionTime = &v
}

func (o ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.EndpointSetNumName) {
		toSerialize["endpoint_set_num_name"] = o.EndpointSetNumName
	}
	if !IsNil(o.EndpointSetForAllOthersTargetUpgradeVersion) {
		toSerialize["endpoint_set_for_all_others_target_upgrade_version"] = o.EndpointSetForAllOthersTargetUpgradeVersion
	}
	if !IsNil(o.EndpointSetForAllOthersUniqueIdentifier) {
		toSerialize["endpoint_set_for_all_others_unique_identifier"] = o.EndpointSetForAllOthersUniqueIdentifier
	}
	if !IsNil(o.EndpointSetNumOnSummary) {
		toSerialize["endpoint_set_num_on_summary"] = o.EndpointSetNumOnSummary
	}
	if !IsNil(o.EndpointSetForAllOthersTargetUpgradeVersionTime) {
		toSerialize["endpoint_set_for_all_others_target_upgrade_version_time"] = o.EndpointSetForAllOthersTargetUpgradeVersionTime
	}
	return toSerialize, nil
}

type NullableImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner struct {
	value *ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner
	isSet bool
}

func (v NullableImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner) Get() *ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner {
	return v.value
}

func (v *NullableImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner) Set(val *ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner) {
	v.value = val
	v.isSet = true
}

func (v NullableImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner) IsSet() bool {
	return v.isSet
}

func (v *NullableImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner(val *ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner) *NullableImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner {
	return &NullableImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner{value: val, isSet: true}
}

func (v NullableImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


