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

// checks if the ConfigPutRequestPacketQueuePacketQueueNameQueueInner type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ConfigPutRequestPacketQueuePacketQueueNameQueueInner{}

// ConfigPutRequestPacketQueuePacketQueueNameQueueInner struct for ConfigPutRequestPacketQueuePacketQueueNameQueueInner
type ConfigPutRequestPacketQueuePacketQueueNameQueueInner struct {
	// Percentage bandwidth allocated to Queue. 0 is no limit
	BandwidthForQueue NullableInt32 `json:"bandwidth_for_queue,omitempty"`
	// Scheduler Type for Queue
	SchedulerType *string `json:"scheduler_type,omitempty"`
	// Weight associated with WRR or DWRR scheduler
	SchedulerWeight NullableInt32 `json:"scheduler_weight,omitempty"`
	// The index identifying the object. Zero if you want to add an object to the list.
	Index *int32 `json:"index,omitempty"`
}

// NewConfigPutRequestPacketQueuePacketQueueNameQueueInner instantiates a new ConfigPutRequestPacketQueuePacketQueueNameQueueInner object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConfigPutRequestPacketQueuePacketQueueNameQueueInner() *ConfigPutRequestPacketQueuePacketQueueNameQueueInner {
	this := ConfigPutRequestPacketQueuePacketQueueNameQueueInner{}
	var bandwidthForQueue int32 = 0
	this.BandwidthForQueue = *NewNullableInt32(&bandwidthForQueue)
	var schedulerType string = "SP"
	this.SchedulerType = &schedulerType
	var schedulerWeight int32 = 0
	this.SchedulerWeight = *NewNullableInt32(&schedulerWeight)
	return &this
}

// NewConfigPutRequestPacketQueuePacketQueueNameQueueInnerWithDefaults instantiates a new ConfigPutRequestPacketQueuePacketQueueNameQueueInner object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConfigPutRequestPacketQueuePacketQueueNameQueueInnerWithDefaults() *ConfigPutRequestPacketQueuePacketQueueNameQueueInner {
	this := ConfigPutRequestPacketQueuePacketQueueNameQueueInner{}
	var bandwidthForQueue int32 = 0
	this.BandwidthForQueue = *NewNullableInt32(&bandwidthForQueue)
	var schedulerType string = "SP"
	this.SchedulerType = &schedulerType
	var schedulerWeight int32 = 0
	this.SchedulerWeight = *NewNullableInt32(&schedulerWeight)
	return &this
}

// GetBandwidthForQueue returns the BandwidthForQueue field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) GetBandwidthForQueue() int32 {
	if o == nil || IsNil(o.BandwidthForQueue.Get()) {
		var ret int32
		return ret
	}
	return *o.BandwidthForQueue.Get()
}

// GetBandwidthForQueueOk returns a tuple with the BandwidthForQueue field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) GetBandwidthForQueueOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.BandwidthForQueue.Get(), o.BandwidthForQueue.IsSet()
}

// HasBandwidthForQueue returns a boolean if a field has been set.
func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) HasBandwidthForQueue() bool {
	if o != nil && o.BandwidthForQueue.IsSet() {
		return true
	}

	return false
}

// SetBandwidthForQueue gets a reference to the given NullableInt32 and assigns it to the BandwidthForQueue field.
func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) SetBandwidthForQueue(v int32) {
	o.BandwidthForQueue.Set(&v)
}
// SetBandwidthForQueueNil sets the value for BandwidthForQueue to be an explicit nil
func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) SetBandwidthForQueueNil() {
	o.BandwidthForQueue.Set(nil)
}

// UnsetBandwidthForQueue ensures that no value is present for BandwidthForQueue, not even an explicit nil
func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) UnsetBandwidthForQueue() {
	o.BandwidthForQueue.Unset()
}

// GetSchedulerType returns the SchedulerType field value if set, zero value otherwise.
func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) GetSchedulerType() string {
	if o == nil || IsNil(o.SchedulerType) {
		var ret string
		return ret
	}
	return *o.SchedulerType
}

// GetSchedulerTypeOk returns a tuple with the SchedulerType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) GetSchedulerTypeOk() (*string, bool) {
	if o == nil || IsNil(o.SchedulerType) {
		return nil, false
	}
	return o.SchedulerType, true
}

// HasSchedulerType returns a boolean if a field has been set.
func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) HasSchedulerType() bool {
	if o != nil && !IsNil(o.SchedulerType) {
		return true
	}

	return false
}

// SetSchedulerType gets a reference to the given string and assigns it to the SchedulerType field.
func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) SetSchedulerType(v string) {
	o.SchedulerType = &v
}

// GetSchedulerWeight returns the SchedulerWeight field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) GetSchedulerWeight() int32 {
	if o == nil || IsNil(o.SchedulerWeight.Get()) {
		var ret int32
		return ret
	}
	return *o.SchedulerWeight.Get()
}

// GetSchedulerWeightOk returns a tuple with the SchedulerWeight field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) GetSchedulerWeightOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.SchedulerWeight.Get(), o.SchedulerWeight.IsSet()
}

// HasSchedulerWeight returns a boolean if a field has been set.
func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) HasSchedulerWeight() bool {
	if o != nil && o.SchedulerWeight.IsSet() {
		return true
	}

	return false
}

// SetSchedulerWeight gets a reference to the given NullableInt32 and assigns it to the SchedulerWeight field.
func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) SetSchedulerWeight(v int32) {
	o.SchedulerWeight.Set(&v)
}
// SetSchedulerWeightNil sets the value for SchedulerWeight to be an explicit nil
func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) SetSchedulerWeightNil() {
	o.SchedulerWeight.Set(nil)
}

// UnsetSchedulerWeight ensures that no value is present for SchedulerWeight, not even an explicit nil
func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) UnsetSchedulerWeight() {
	o.SchedulerWeight.Unset()
}

// GetIndex returns the Index field value if set, zero value otherwise.
func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) GetIndex() int32 {
	if o == nil || IsNil(o.Index) {
		var ret int32
		return ret
	}
	return *o.Index
}

// GetIndexOk returns a tuple with the Index field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) GetIndexOk() (*int32, bool) {
	if o == nil || IsNil(o.Index) {
		return nil, false
	}
	return o.Index, true
}

// HasIndex returns a boolean if a field has been set.
func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) HasIndex() bool {
	if o != nil && !IsNil(o.Index) {
		return true
	}

	return false
}

// SetIndex gets a reference to the given int32 and assigns it to the Index field.
func (o *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) SetIndex(v int32) {
	o.Index = &v
}

func (o ConfigPutRequestPacketQueuePacketQueueNameQueueInner) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ConfigPutRequestPacketQueuePacketQueueNameQueueInner) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if o.BandwidthForQueue.IsSet() {
		toSerialize["bandwidth_for_queue"] = o.BandwidthForQueue.Get()
	}
	if !IsNil(o.SchedulerType) {
		toSerialize["scheduler_type"] = o.SchedulerType
	}
	if o.SchedulerWeight.IsSet() {
		toSerialize["scheduler_weight"] = o.SchedulerWeight.Get()
	}
	if !IsNil(o.Index) {
		toSerialize["index"] = o.Index
	}
	return toSerialize, nil
}

type NullableConfigPutRequestPacketQueuePacketQueueNameQueueInner struct {
	value *ConfigPutRequestPacketQueuePacketQueueNameQueueInner
	isSet bool
}

func (v NullableConfigPutRequestPacketQueuePacketQueueNameQueueInner) Get() *ConfigPutRequestPacketQueuePacketQueueNameQueueInner {
	return v.value
}

func (v *NullableConfigPutRequestPacketQueuePacketQueueNameQueueInner) Set(val *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) {
	v.value = val
	v.isSet = true
}

func (v NullableConfigPutRequestPacketQueuePacketQueueNameQueueInner) IsSet() bool {
	return v.isSet
}

func (v *NullableConfigPutRequestPacketQueuePacketQueueNameQueueInner) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConfigPutRequestPacketQueuePacketQueueNameQueueInner(val *ConfigPutRequestPacketQueuePacketQueueNameQueueInner) *NullableConfigPutRequestPacketQueuePacketQueueNameQueueInner {
	return &NullableConfigPutRequestPacketQueuePacketQueueNameQueueInner{value: val, isSet: true}
}

func (v NullableConfigPutRequestPacketQueuePacketQueueNameQueueInner) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConfigPutRequestPacketQueuePacketQueueNameQueueInner) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


