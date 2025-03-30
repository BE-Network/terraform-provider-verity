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

// checks if the ConfigPutRequestThresholdsThresholdsNameThresholdInner type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ConfigPutRequestThresholdsThresholdsNameThresholdInner{}

// ConfigPutRequestThresholdsThresholdsNameThresholdInner struct for ConfigPutRequestThresholdsThresholdsNameThresholdInner
type ConfigPutRequestThresholdsThresholdsNameThresholdInner struct {
	// Enable/Disable the threshold
	ThresholdNumEnabled *bool `json:"threshold_num_enabled,omitempty"`
	// Stat Type
	ThresholdNumStatType *string `json:"threshold_num_stat_type,omitempty"`
	// Method
	ThresholdNumMethod *string `json:"threshold_num_method,omitempty"`
	// Value
	ThresholdNumValue *string `json:"threshold_num_value,omitempty"`
	// Action
	ThresholdNumAction *string `json:"threshold_num_action,omitempty"`
	// Alarm clear method
	ThresholdNumAlarmClearMethod *string `json:"threshold_num_alarm_clear_method,omitempty"`
	// Alarm clear value
	ThresholdNumAlarmClearValue *string `json:"threshold_num_alarm_clear_value,omitempty"`
	// The index identifying the object. Zero if you want to add an object to the list.
	Index *int32 `json:"index,omitempty"`
}

// NewConfigPutRequestThresholdsThresholdsNameThresholdInner instantiates a new ConfigPutRequestThresholdsThresholdsNameThresholdInner object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConfigPutRequestThresholdsThresholdsNameThresholdInner() *ConfigPutRequestThresholdsThresholdsNameThresholdInner {
	this := ConfigPutRequestThresholdsThresholdsNameThresholdInner{}
	var thresholdNumEnabled bool = false
	this.ThresholdNumEnabled = &thresholdNumEnabled
	var thresholdNumStatType string = ""
	this.ThresholdNumStatType = &thresholdNumStatType
	var thresholdNumMethod string = ""
	this.ThresholdNumMethod = &thresholdNumMethod
	var thresholdNumValue string = ""
	this.ThresholdNumValue = &thresholdNumValue
	var thresholdNumAction string = ""
	this.ThresholdNumAction = &thresholdNumAction
	var thresholdNumAlarmClearMethod string = ""
	this.ThresholdNumAlarmClearMethod = &thresholdNumAlarmClearMethod
	var thresholdNumAlarmClearValue string = ""
	this.ThresholdNumAlarmClearValue = &thresholdNumAlarmClearValue
	return &this
}

// NewConfigPutRequestThresholdsThresholdsNameThresholdInnerWithDefaults instantiates a new ConfigPutRequestThresholdsThresholdsNameThresholdInner object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConfigPutRequestThresholdsThresholdsNameThresholdInnerWithDefaults() *ConfigPutRequestThresholdsThresholdsNameThresholdInner {
	this := ConfigPutRequestThresholdsThresholdsNameThresholdInner{}
	var thresholdNumEnabled bool = false
	this.ThresholdNumEnabled = &thresholdNumEnabled
	var thresholdNumStatType string = ""
	this.ThresholdNumStatType = &thresholdNumStatType
	var thresholdNumMethod string = ""
	this.ThresholdNumMethod = &thresholdNumMethod
	var thresholdNumValue string = ""
	this.ThresholdNumValue = &thresholdNumValue
	var thresholdNumAction string = ""
	this.ThresholdNumAction = &thresholdNumAction
	var thresholdNumAlarmClearMethod string = ""
	this.ThresholdNumAlarmClearMethod = &thresholdNumAlarmClearMethod
	var thresholdNumAlarmClearValue string = ""
	this.ThresholdNumAlarmClearValue = &thresholdNumAlarmClearValue
	return &this
}

// GetThresholdNumEnabled returns the ThresholdNumEnabled field value if set, zero value otherwise.
func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) GetThresholdNumEnabled() bool {
	if o == nil || IsNil(o.ThresholdNumEnabled) {
		var ret bool
		return ret
	}
	return *o.ThresholdNumEnabled
}

// GetThresholdNumEnabledOk returns a tuple with the ThresholdNumEnabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) GetThresholdNumEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.ThresholdNumEnabled) {
		return nil, false
	}
	return o.ThresholdNumEnabled, true
}

// HasThresholdNumEnabled returns a boolean if a field has been set.
func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) HasThresholdNumEnabled() bool {
	if o != nil && !IsNil(o.ThresholdNumEnabled) {
		return true
	}

	return false
}

// SetThresholdNumEnabled gets a reference to the given bool and assigns it to the ThresholdNumEnabled field.
func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) SetThresholdNumEnabled(v bool) {
	o.ThresholdNumEnabled = &v
}

// GetThresholdNumStatType returns the ThresholdNumStatType field value if set, zero value otherwise.
func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) GetThresholdNumStatType() string {
	if o == nil || IsNil(o.ThresholdNumStatType) {
		var ret string
		return ret
	}
	return *o.ThresholdNumStatType
}

// GetThresholdNumStatTypeOk returns a tuple with the ThresholdNumStatType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) GetThresholdNumStatTypeOk() (*string, bool) {
	if o == nil || IsNil(o.ThresholdNumStatType) {
		return nil, false
	}
	return o.ThresholdNumStatType, true
}

// HasThresholdNumStatType returns a boolean if a field has been set.
func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) HasThresholdNumStatType() bool {
	if o != nil && !IsNil(o.ThresholdNumStatType) {
		return true
	}

	return false
}

// SetThresholdNumStatType gets a reference to the given string and assigns it to the ThresholdNumStatType field.
func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) SetThresholdNumStatType(v string) {
	o.ThresholdNumStatType = &v
}

// GetThresholdNumMethod returns the ThresholdNumMethod field value if set, zero value otherwise.
func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) GetThresholdNumMethod() string {
	if o == nil || IsNil(o.ThresholdNumMethod) {
		var ret string
		return ret
	}
	return *o.ThresholdNumMethod
}

// GetThresholdNumMethodOk returns a tuple with the ThresholdNumMethod field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) GetThresholdNumMethodOk() (*string, bool) {
	if o == nil || IsNil(o.ThresholdNumMethod) {
		return nil, false
	}
	return o.ThresholdNumMethod, true
}

// HasThresholdNumMethod returns a boolean if a field has been set.
func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) HasThresholdNumMethod() bool {
	if o != nil && !IsNil(o.ThresholdNumMethod) {
		return true
	}

	return false
}

// SetThresholdNumMethod gets a reference to the given string and assigns it to the ThresholdNumMethod field.
func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) SetThresholdNumMethod(v string) {
	o.ThresholdNumMethod = &v
}

// GetThresholdNumValue returns the ThresholdNumValue field value if set, zero value otherwise.
func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) GetThresholdNumValue() string {
	if o == nil || IsNil(o.ThresholdNumValue) {
		var ret string
		return ret
	}
	return *o.ThresholdNumValue
}

// GetThresholdNumValueOk returns a tuple with the ThresholdNumValue field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) GetThresholdNumValueOk() (*string, bool) {
	if o == nil || IsNil(o.ThresholdNumValue) {
		return nil, false
	}
	return o.ThresholdNumValue, true
}

// HasThresholdNumValue returns a boolean if a field has been set.
func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) HasThresholdNumValue() bool {
	if o != nil && !IsNil(o.ThresholdNumValue) {
		return true
	}

	return false
}

// SetThresholdNumValue gets a reference to the given string and assigns it to the ThresholdNumValue field.
func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) SetThresholdNumValue(v string) {
	o.ThresholdNumValue = &v
}

// GetThresholdNumAction returns the ThresholdNumAction field value if set, zero value otherwise.
func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) GetThresholdNumAction() string {
	if o == nil || IsNil(o.ThresholdNumAction) {
		var ret string
		return ret
	}
	return *o.ThresholdNumAction
}

// GetThresholdNumActionOk returns a tuple with the ThresholdNumAction field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) GetThresholdNumActionOk() (*string, bool) {
	if o == nil || IsNil(o.ThresholdNumAction) {
		return nil, false
	}
	return o.ThresholdNumAction, true
}

// HasThresholdNumAction returns a boolean if a field has been set.
func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) HasThresholdNumAction() bool {
	if o != nil && !IsNil(o.ThresholdNumAction) {
		return true
	}

	return false
}

// SetThresholdNumAction gets a reference to the given string and assigns it to the ThresholdNumAction field.
func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) SetThresholdNumAction(v string) {
	o.ThresholdNumAction = &v
}

// GetThresholdNumAlarmClearMethod returns the ThresholdNumAlarmClearMethod field value if set, zero value otherwise.
func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) GetThresholdNumAlarmClearMethod() string {
	if o == nil || IsNil(o.ThresholdNumAlarmClearMethod) {
		var ret string
		return ret
	}
	return *o.ThresholdNumAlarmClearMethod
}

// GetThresholdNumAlarmClearMethodOk returns a tuple with the ThresholdNumAlarmClearMethod field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) GetThresholdNumAlarmClearMethodOk() (*string, bool) {
	if o == nil || IsNil(o.ThresholdNumAlarmClearMethod) {
		return nil, false
	}
	return o.ThresholdNumAlarmClearMethod, true
}

// HasThresholdNumAlarmClearMethod returns a boolean if a field has been set.
func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) HasThresholdNumAlarmClearMethod() bool {
	if o != nil && !IsNil(o.ThresholdNumAlarmClearMethod) {
		return true
	}

	return false
}

// SetThresholdNumAlarmClearMethod gets a reference to the given string and assigns it to the ThresholdNumAlarmClearMethod field.
func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) SetThresholdNumAlarmClearMethod(v string) {
	o.ThresholdNumAlarmClearMethod = &v
}

// GetThresholdNumAlarmClearValue returns the ThresholdNumAlarmClearValue field value if set, zero value otherwise.
func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) GetThresholdNumAlarmClearValue() string {
	if o == nil || IsNil(o.ThresholdNumAlarmClearValue) {
		var ret string
		return ret
	}
	return *o.ThresholdNumAlarmClearValue
}

// GetThresholdNumAlarmClearValueOk returns a tuple with the ThresholdNumAlarmClearValue field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) GetThresholdNumAlarmClearValueOk() (*string, bool) {
	if o == nil || IsNil(o.ThresholdNumAlarmClearValue) {
		return nil, false
	}
	return o.ThresholdNumAlarmClearValue, true
}

// HasThresholdNumAlarmClearValue returns a boolean if a field has been set.
func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) HasThresholdNumAlarmClearValue() bool {
	if o != nil && !IsNil(o.ThresholdNumAlarmClearValue) {
		return true
	}

	return false
}

// SetThresholdNumAlarmClearValue gets a reference to the given string and assigns it to the ThresholdNumAlarmClearValue field.
func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) SetThresholdNumAlarmClearValue(v string) {
	o.ThresholdNumAlarmClearValue = &v
}

// GetIndex returns the Index field value if set, zero value otherwise.
func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) GetIndex() int32 {
	if o == nil || IsNil(o.Index) {
		var ret int32
		return ret
	}
	return *o.Index
}

// GetIndexOk returns a tuple with the Index field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) GetIndexOk() (*int32, bool) {
	if o == nil || IsNil(o.Index) {
		return nil, false
	}
	return o.Index, true
}

// HasIndex returns a boolean if a field has been set.
func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) HasIndex() bool {
	if o != nil && !IsNil(o.Index) {
		return true
	}

	return false
}

// SetIndex gets a reference to the given int32 and assigns it to the Index field.
func (o *ConfigPutRequestThresholdsThresholdsNameThresholdInner) SetIndex(v int32) {
	o.Index = &v
}

func (o ConfigPutRequestThresholdsThresholdsNameThresholdInner) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ConfigPutRequestThresholdsThresholdsNameThresholdInner) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.ThresholdNumEnabled) {
		toSerialize["threshold_num_enabled"] = o.ThresholdNumEnabled
	}
	if !IsNil(o.ThresholdNumStatType) {
		toSerialize["threshold_num_stat_type"] = o.ThresholdNumStatType
	}
	if !IsNil(o.ThresholdNumMethod) {
		toSerialize["threshold_num_method"] = o.ThresholdNumMethod
	}
	if !IsNil(o.ThresholdNumValue) {
		toSerialize["threshold_num_value"] = o.ThresholdNumValue
	}
	if !IsNil(o.ThresholdNumAction) {
		toSerialize["threshold_num_action"] = o.ThresholdNumAction
	}
	if !IsNil(o.ThresholdNumAlarmClearMethod) {
		toSerialize["threshold_num_alarm_clear_method"] = o.ThresholdNumAlarmClearMethod
	}
	if !IsNil(o.ThresholdNumAlarmClearValue) {
		toSerialize["threshold_num_alarm_clear_value"] = o.ThresholdNumAlarmClearValue
	}
	if !IsNil(o.Index) {
		toSerialize["index"] = o.Index
	}
	return toSerialize, nil
}

type NullableConfigPutRequestThresholdsThresholdsNameThresholdInner struct {
	value *ConfigPutRequestThresholdsThresholdsNameThresholdInner
	isSet bool
}

func (v NullableConfigPutRequestThresholdsThresholdsNameThresholdInner) Get() *ConfigPutRequestThresholdsThresholdsNameThresholdInner {
	return v.value
}

func (v *NullableConfigPutRequestThresholdsThresholdsNameThresholdInner) Set(val *ConfigPutRequestThresholdsThresholdsNameThresholdInner) {
	v.value = val
	v.isSet = true
}

func (v NullableConfigPutRequestThresholdsThresholdsNameThresholdInner) IsSet() bool {
	return v.isSet
}

func (v *NullableConfigPutRequestThresholdsThresholdsNameThresholdInner) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConfigPutRequestThresholdsThresholdsNameThresholdInner(val *ConfigPutRequestThresholdsThresholdsNameThresholdInner) *NullableConfigPutRequestThresholdsThresholdsNameThresholdInner {
	return &NullableConfigPutRequestThresholdsThresholdsNameThresholdInner{value: val, isSet: true}
}

func (v NullableConfigPutRequestThresholdsThresholdsNameThresholdInner) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConfigPutRequestThresholdsThresholdsNameThresholdInner) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


