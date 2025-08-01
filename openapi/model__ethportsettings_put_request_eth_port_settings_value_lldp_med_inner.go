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

// checks if the EthportsettingsPutRequestEthPortSettingsValueLldpMedInner type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &EthportsettingsPutRequestEthPortSettingsValueLldpMedInner{}

// EthportsettingsPutRequestEthPortSettingsValueLldpMedInner struct for EthportsettingsPutRequestEthPortSettingsValueLldpMedInner
type EthportsettingsPutRequestEthPortSettingsValueLldpMedInner struct {
	// Per LLDP Med row enable
	LldpMedRowNumEnable *bool `json:"lldp_med_row_num_enable,omitempty"`
	// Advertised application
	LldpMedRowNumAdvertisedApplicatio *string `json:"lldp_med_row_num_advertised_applicatio,omitempty"`
	// LLDP DSCP Mark
	LldpMedRowNumDscpMark *int32 `json:"lldp_med_row_num_dscp_mark,omitempty"`
	// LLDP Priority
	LldpMedRowNumPriority *int32 `json:"lldp_med_row_num_priority,omitempty"`
	// LLDP Service
	LldpMedRowNumService *string `json:"lldp_med_row_num_service,omitempty"`
	// Object type for lldp_med_row_num_service field
	LldpMedRowNumServiceRefType *string `json:"lldp_med_row_num_service_ref_type_,omitempty"`
	// The index identifying the object. Zero if you want to add an object to the list.
	Index *int32 `json:"index,omitempty"`
}

// NewEthportsettingsPutRequestEthPortSettingsValueLldpMedInner instantiates a new EthportsettingsPutRequestEthPortSettingsValueLldpMedInner object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewEthportsettingsPutRequestEthPortSettingsValueLldpMedInner() *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner {
	this := EthportsettingsPutRequestEthPortSettingsValueLldpMedInner{}
	var lldpMedRowNumEnable bool = false
	this.LldpMedRowNumEnable = &lldpMedRowNumEnable
	var lldpMedRowNumAdvertisedApplicatio string = ""
	this.LldpMedRowNumAdvertisedApplicatio = &lldpMedRowNumAdvertisedApplicatio
	var lldpMedRowNumDscpMark int32 = 0
	this.LldpMedRowNumDscpMark = &lldpMedRowNumDscpMark
	var lldpMedRowNumPriority int32 = 0
	this.LldpMedRowNumPriority = &lldpMedRowNumPriority
	var lldpMedRowNumService string = ""
	this.LldpMedRowNumService = &lldpMedRowNumService
	return &this
}

// NewEthportsettingsPutRequestEthPortSettingsValueLldpMedInnerWithDefaults instantiates a new EthportsettingsPutRequestEthPortSettingsValueLldpMedInner object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewEthportsettingsPutRequestEthPortSettingsValueLldpMedInnerWithDefaults() *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner {
	this := EthportsettingsPutRequestEthPortSettingsValueLldpMedInner{}
	var lldpMedRowNumEnable bool = false
	this.LldpMedRowNumEnable = &lldpMedRowNumEnable
	var lldpMedRowNumAdvertisedApplicatio string = ""
	this.LldpMedRowNumAdvertisedApplicatio = &lldpMedRowNumAdvertisedApplicatio
	var lldpMedRowNumDscpMark int32 = 0
	this.LldpMedRowNumDscpMark = &lldpMedRowNumDscpMark
	var lldpMedRowNumPriority int32 = 0
	this.LldpMedRowNumPriority = &lldpMedRowNumPriority
	var lldpMedRowNumService string = ""
	this.LldpMedRowNumService = &lldpMedRowNumService
	return &this
}

// GetLldpMedRowNumEnable returns the LldpMedRowNumEnable field value if set, zero value otherwise.
func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) GetLldpMedRowNumEnable() bool {
	if o == nil || IsNil(o.LldpMedRowNumEnable) {
		var ret bool
		return ret
	}
	return *o.LldpMedRowNumEnable
}

// GetLldpMedRowNumEnableOk returns a tuple with the LldpMedRowNumEnable field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) GetLldpMedRowNumEnableOk() (*bool, bool) {
	if o == nil || IsNil(o.LldpMedRowNumEnable) {
		return nil, false
	}
	return o.LldpMedRowNumEnable, true
}

// HasLldpMedRowNumEnable returns a boolean if a field has been set.
func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) HasLldpMedRowNumEnable() bool {
	if o != nil && !IsNil(o.LldpMedRowNumEnable) {
		return true
	}

	return false
}

// SetLldpMedRowNumEnable gets a reference to the given bool and assigns it to the LldpMedRowNumEnable field.
func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) SetLldpMedRowNumEnable(v bool) {
	o.LldpMedRowNumEnable = &v
}

// GetLldpMedRowNumAdvertisedApplicatio returns the LldpMedRowNumAdvertisedApplicatio field value if set, zero value otherwise.
func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) GetLldpMedRowNumAdvertisedApplicatio() string {
	if o == nil || IsNil(o.LldpMedRowNumAdvertisedApplicatio) {
		var ret string
		return ret
	}
	return *o.LldpMedRowNumAdvertisedApplicatio
}

// GetLldpMedRowNumAdvertisedApplicatioOk returns a tuple with the LldpMedRowNumAdvertisedApplicatio field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) GetLldpMedRowNumAdvertisedApplicatioOk() (*string, bool) {
	if o == nil || IsNil(o.LldpMedRowNumAdvertisedApplicatio) {
		return nil, false
	}
	return o.LldpMedRowNumAdvertisedApplicatio, true
}

// HasLldpMedRowNumAdvertisedApplicatio returns a boolean if a field has been set.
func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) HasLldpMedRowNumAdvertisedApplicatio() bool {
	if o != nil && !IsNil(o.LldpMedRowNumAdvertisedApplicatio) {
		return true
	}

	return false
}

// SetLldpMedRowNumAdvertisedApplicatio gets a reference to the given string and assigns it to the LldpMedRowNumAdvertisedApplicatio field.
func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) SetLldpMedRowNumAdvertisedApplicatio(v string) {
	o.LldpMedRowNumAdvertisedApplicatio = &v
}

// GetLldpMedRowNumDscpMark returns the LldpMedRowNumDscpMark field value if set, zero value otherwise.
func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) GetLldpMedRowNumDscpMark() int32 {
	if o == nil || IsNil(o.LldpMedRowNumDscpMark) {
		var ret int32
		return ret
	}
	return *o.LldpMedRowNumDscpMark
}

// GetLldpMedRowNumDscpMarkOk returns a tuple with the LldpMedRowNumDscpMark field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) GetLldpMedRowNumDscpMarkOk() (*int32, bool) {
	if o == nil || IsNil(o.LldpMedRowNumDscpMark) {
		return nil, false
	}
	return o.LldpMedRowNumDscpMark, true
}

// HasLldpMedRowNumDscpMark returns a boolean if a field has been set.
func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) HasLldpMedRowNumDscpMark() bool {
	if o != nil && !IsNil(o.LldpMedRowNumDscpMark) {
		return true
	}

	return false
}

// SetLldpMedRowNumDscpMark gets a reference to the given int32 and assigns it to the LldpMedRowNumDscpMark field.
func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) SetLldpMedRowNumDscpMark(v int32) {
	o.LldpMedRowNumDscpMark = &v
}

// GetLldpMedRowNumPriority returns the LldpMedRowNumPriority field value if set, zero value otherwise.
func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) GetLldpMedRowNumPriority() int32 {
	if o == nil || IsNil(o.LldpMedRowNumPriority) {
		var ret int32
		return ret
	}
	return *o.LldpMedRowNumPriority
}

// GetLldpMedRowNumPriorityOk returns a tuple with the LldpMedRowNumPriority field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) GetLldpMedRowNumPriorityOk() (*int32, bool) {
	if o == nil || IsNil(o.LldpMedRowNumPriority) {
		return nil, false
	}
	return o.LldpMedRowNumPriority, true
}

// HasLldpMedRowNumPriority returns a boolean if a field has been set.
func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) HasLldpMedRowNumPriority() bool {
	if o != nil && !IsNil(o.LldpMedRowNumPriority) {
		return true
	}

	return false
}

// SetLldpMedRowNumPriority gets a reference to the given int32 and assigns it to the LldpMedRowNumPriority field.
func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) SetLldpMedRowNumPriority(v int32) {
	o.LldpMedRowNumPriority = &v
}

// GetLldpMedRowNumService returns the LldpMedRowNumService field value if set, zero value otherwise.
func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) GetLldpMedRowNumService() string {
	if o == nil || IsNil(o.LldpMedRowNumService) {
		var ret string
		return ret
	}
	return *o.LldpMedRowNumService
}

// GetLldpMedRowNumServiceOk returns a tuple with the LldpMedRowNumService field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) GetLldpMedRowNumServiceOk() (*string, bool) {
	if o == nil || IsNil(o.LldpMedRowNumService) {
		return nil, false
	}
	return o.LldpMedRowNumService, true
}

// HasLldpMedRowNumService returns a boolean if a field has been set.
func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) HasLldpMedRowNumService() bool {
	if o != nil && !IsNil(o.LldpMedRowNumService) {
		return true
	}

	return false
}

// SetLldpMedRowNumService gets a reference to the given string and assigns it to the LldpMedRowNumService field.
func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) SetLldpMedRowNumService(v string) {
	o.LldpMedRowNumService = &v
}

// GetLldpMedRowNumServiceRefType returns the LldpMedRowNumServiceRefType field value if set, zero value otherwise.
func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) GetLldpMedRowNumServiceRefType() string {
	if o == nil || IsNil(o.LldpMedRowNumServiceRefType) {
		var ret string
		return ret
	}
	return *o.LldpMedRowNumServiceRefType
}

// GetLldpMedRowNumServiceRefTypeOk returns a tuple with the LldpMedRowNumServiceRefType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) GetLldpMedRowNumServiceRefTypeOk() (*string, bool) {
	if o == nil || IsNil(o.LldpMedRowNumServiceRefType) {
		return nil, false
	}
	return o.LldpMedRowNumServiceRefType, true
}

// HasLldpMedRowNumServiceRefType returns a boolean if a field has been set.
func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) HasLldpMedRowNumServiceRefType() bool {
	if o != nil && !IsNil(o.LldpMedRowNumServiceRefType) {
		return true
	}

	return false
}

// SetLldpMedRowNumServiceRefType gets a reference to the given string and assigns it to the LldpMedRowNumServiceRefType field.
func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) SetLldpMedRowNumServiceRefType(v string) {
	o.LldpMedRowNumServiceRefType = &v
}

// GetIndex returns the Index field value if set, zero value otherwise.
func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) GetIndex() int32 {
	if o == nil || IsNil(o.Index) {
		var ret int32
		return ret
	}
	return *o.Index
}

// GetIndexOk returns a tuple with the Index field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) GetIndexOk() (*int32, bool) {
	if o == nil || IsNil(o.Index) {
		return nil, false
	}
	return o.Index, true
}

// HasIndex returns a boolean if a field has been set.
func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) HasIndex() bool {
	if o != nil && !IsNil(o.Index) {
		return true
	}

	return false
}

// SetIndex gets a reference to the given int32 and assigns it to the Index field.
func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) SetIndex(v int32) {
	o.Index = &v
}

func (o EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.LldpMedRowNumEnable) {
		toSerialize["lldp_med_row_num_enable"] = o.LldpMedRowNumEnable
	}
	if !IsNil(o.LldpMedRowNumAdvertisedApplicatio) {
		toSerialize["lldp_med_row_num_advertised_applicatio"] = o.LldpMedRowNumAdvertisedApplicatio
	}
	if !IsNil(o.LldpMedRowNumDscpMark) {
		toSerialize["lldp_med_row_num_dscp_mark"] = o.LldpMedRowNumDscpMark
	}
	if !IsNil(o.LldpMedRowNumPriority) {
		toSerialize["lldp_med_row_num_priority"] = o.LldpMedRowNumPriority
	}
	if !IsNil(o.LldpMedRowNumService) {
		toSerialize["lldp_med_row_num_service"] = o.LldpMedRowNumService
	}
	if !IsNil(o.LldpMedRowNumServiceRefType) {
		toSerialize["lldp_med_row_num_service_ref_type_"] = o.LldpMedRowNumServiceRefType
	}
	if !IsNil(o.Index) {
		toSerialize["index"] = o.Index
	}
	return toSerialize, nil
}

type NullableEthportsettingsPutRequestEthPortSettingsValueLldpMedInner struct {
	value *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner
	isSet bool
}

func (v NullableEthportsettingsPutRequestEthPortSettingsValueLldpMedInner) Get() *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner {
	return v.value
}

func (v *NullableEthportsettingsPutRequestEthPortSettingsValueLldpMedInner) Set(val *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) {
	v.value = val
	v.isSet = true
}

func (v NullableEthportsettingsPutRequestEthPortSettingsValueLldpMedInner) IsSet() bool {
	return v.isSet
}

func (v *NullableEthportsettingsPutRequestEthPortSettingsValueLldpMedInner) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableEthportsettingsPutRequestEthPortSettingsValueLldpMedInner(val *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) *NullableEthportsettingsPutRequestEthPortSettingsValueLldpMedInner {
	return &NullableEthportsettingsPutRequestEthPortSettingsValueLldpMedInner{value: val, isSet: true}
}

func (v NullableEthportsettingsPutRequestEthPortSettingsValueLldpMedInner) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableEthportsettingsPutRequestEthPortSettingsValueLldpMedInner) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


