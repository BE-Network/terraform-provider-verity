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

// checks if the ConfigPutRequestSwitchpointSwitchpointNameObjectProperties type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ConfigPutRequestSwitchpointSwitchpointNameObjectProperties{}

// ConfigPutRequestSwitchpointSwitchpointNameObjectProperties struct for ConfigPutRequestSwitchpointSwitchpointNameObjectProperties
type ConfigPutRequestSwitchpointSwitchpointNameObjectProperties struct {
	// Notes writen by User about the site
	UserNotes *string `json:"user_notes,omitempty"`
	// Expected Parent Endpoint
	ExpectedParentEndpoint *string `json:"expected_parent_endpoint,omitempty"`
	// Object type for expected_parent_endpoint field
	ExpectedParentEndpointRefType *string `json:"expected_parent_endpoint_ref_type_,omitempty"`
	// Number of Multipoints
	NumberOfMultipoints NullableInt32 `json:"number_of_multipoints,omitempty"`
	// For Switch Endpoints. Denotes switch aggregated with all of its sub switches
	Aggregate *bool `json:"aggregate,omitempty"`
	// For Switch Endpoints. Denotes the Host Switch
	IsHost *bool `json:"is_host,omitempty"`
	Eths *ConfigPutRequestSwitchpointSwitchpointNameObjectPropertiesEths `json:"eths,omitempty"`
}

// NewConfigPutRequestSwitchpointSwitchpointNameObjectProperties instantiates a new ConfigPutRequestSwitchpointSwitchpointNameObjectProperties object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConfigPutRequestSwitchpointSwitchpointNameObjectProperties() *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties {
	this := ConfigPutRequestSwitchpointSwitchpointNameObjectProperties{}
	var userNotes string = ""
	this.UserNotes = &userNotes
	var expectedParentEndpoint string = ""
	this.ExpectedParentEndpoint = &expectedParentEndpoint
	var aggregate bool = false
	this.Aggregate = &aggregate
	var isHost bool = false
	this.IsHost = &isHost
	return &this
}

// NewConfigPutRequestSwitchpointSwitchpointNameObjectPropertiesWithDefaults instantiates a new ConfigPutRequestSwitchpointSwitchpointNameObjectProperties object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConfigPutRequestSwitchpointSwitchpointNameObjectPropertiesWithDefaults() *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties {
	this := ConfigPutRequestSwitchpointSwitchpointNameObjectProperties{}
	var userNotes string = ""
	this.UserNotes = &userNotes
	var expectedParentEndpoint string = ""
	this.ExpectedParentEndpoint = &expectedParentEndpoint
	var aggregate bool = false
	this.Aggregate = &aggregate
	var isHost bool = false
	this.IsHost = &isHost
	return &this
}

// GetUserNotes returns the UserNotes field value if set, zero value otherwise.
func (o *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) GetUserNotes() string {
	if o == nil || IsNil(o.UserNotes) {
		var ret string
		return ret
	}
	return *o.UserNotes
}

// GetUserNotesOk returns a tuple with the UserNotes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) GetUserNotesOk() (*string, bool) {
	if o == nil || IsNil(o.UserNotes) {
		return nil, false
	}
	return o.UserNotes, true
}

// HasUserNotes returns a boolean if a field has been set.
func (o *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) HasUserNotes() bool {
	if o != nil && !IsNil(o.UserNotes) {
		return true
	}

	return false
}

// SetUserNotes gets a reference to the given string and assigns it to the UserNotes field.
func (o *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) SetUserNotes(v string) {
	o.UserNotes = &v
}

// GetExpectedParentEndpoint returns the ExpectedParentEndpoint field value if set, zero value otherwise.
func (o *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) GetExpectedParentEndpoint() string {
	if o == nil || IsNil(o.ExpectedParentEndpoint) {
		var ret string
		return ret
	}
	return *o.ExpectedParentEndpoint
}

// GetExpectedParentEndpointOk returns a tuple with the ExpectedParentEndpoint field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) GetExpectedParentEndpointOk() (*string, bool) {
	if o == nil || IsNil(o.ExpectedParentEndpoint) {
		return nil, false
	}
	return o.ExpectedParentEndpoint, true
}

// HasExpectedParentEndpoint returns a boolean if a field has been set.
func (o *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) HasExpectedParentEndpoint() bool {
	if o != nil && !IsNil(o.ExpectedParentEndpoint) {
		return true
	}

	return false
}

// SetExpectedParentEndpoint gets a reference to the given string and assigns it to the ExpectedParentEndpoint field.
func (o *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) SetExpectedParentEndpoint(v string) {
	o.ExpectedParentEndpoint = &v
}

// GetExpectedParentEndpointRefType returns the ExpectedParentEndpointRefType field value if set, zero value otherwise.
func (o *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) GetExpectedParentEndpointRefType() string {
	if o == nil || IsNil(o.ExpectedParentEndpointRefType) {
		var ret string
		return ret
	}
	return *o.ExpectedParentEndpointRefType
}

// GetExpectedParentEndpointRefTypeOk returns a tuple with the ExpectedParentEndpointRefType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) GetExpectedParentEndpointRefTypeOk() (*string, bool) {
	if o == nil || IsNil(o.ExpectedParentEndpointRefType) {
		return nil, false
	}
	return o.ExpectedParentEndpointRefType, true
}

// HasExpectedParentEndpointRefType returns a boolean if a field has been set.
func (o *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) HasExpectedParentEndpointRefType() bool {
	if o != nil && !IsNil(o.ExpectedParentEndpointRefType) {
		return true
	}

	return false
}

// SetExpectedParentEndpointRefType gets a reference to the given string and assigns it to the ExpectedParentEndpointRefType field.
func (o *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) SetExpectedParentEndpointRefType(v string) {
	o.ExpectedParentEndpointRefType = &v
}

// GetNumberOfMultipoints returns the NumberOfMultipoints field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) GetNumberOfMultipoints() int32 {
	if o == nil || IsNil(o.NumberOfMultipoints.Get()) {
		var ret int32
		return ret
	}
	return *o.NumberOfMultipoints.Get()
}

// GetNumberOfMultipointsOk returns a tuple with the NumberOfMultipoints field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) GetNumberOfMultipointsOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.NumberOfMultipoints.Get(), o.NumberOfMultipoints.IsSet()
}

// HasNumberOfMultipoints returns a boolean if a field has been set.
func (o *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) HasNumberOfMultipoints() bool {
	if o != nil && o.NumberOfMultipoints.IsSet() {
		return true
	}

	return false
}

// SetNumberOfMultipoints gets a reference to the given NullableInt32 and assigns it to the NumberOfMultipoints field.
func (o *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) SetNumberOfMultipoints(v int32) {
	o.NumberOfMultipoints.Set(&v)
}
// SetNumberOfMultipointsNil sets the value for NumberOfMultipoints to be an explicit nil
func (o *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) SetNumberOfMultipointsNil() {
	o.NumberOfMultipoints.Set(nil)
}

// UnsetNumberOfMultipoints ensures that no value is present for NumberOfMultipoints, not even an explicit nil
func (o *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) UnsetNumberOfMultipoints() {
	o.NumberOfMultipoints.Unset()
}

// GetAggregate returns the Aggregate field value if set, zero value otherwise.
func (o *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) GetAggregate() bool {
	if o == nil || IsNil(o.Aggregate) {
		var ret bool
		return ret
	}
	return *o.Aggregate
}

// GetAggregateOk returns a tuple with the Aggregate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) GetAggregateOk() (*bool, bool) {
	if o == nil || IsNil(o.Aggregate) {
		return nil, false
	}
	return o.Aggregate, true
}

// HasAggregate returns a boolean if a field has been set.
func (o *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) HasAggregate() bool {
	if o != nil && !IsNil(o.Aggregate) {
		return true
	}

	return false
}

// SetAggregate gets a reference to the given bool and assigns it to the Aggregate field.
func (o *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) SetAggregate(v bool) {
	o.Aggregate = &v
}

// GetIsHost returns the IsHost field value if set, zero value otherwise.
func (o *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) GetIsHost() bool {
	if o == nil || IsNil(o.IsHost) {
		var ret bool
		return ret
	}
	return *o.IsHost
}

// GetIsHostOk returns a tuple with the IsHost field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) GetIsHostOk() (*bool, bool) {
	if o == nil || IsNil(o.IsHost) {
		return nil, false
	}
	return o.IsHost, true
}

// HasIsHost returns a boolean if a field has been set.
func (o *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) HasIsHost() bool {
	if o != nil && !IsNil(o.IsHost) {
		return true
	}

	return false
}

// SetIsHost gets a reference to the given bool and assigns it to the IsHost field.
func (o *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) SetIsHost(v bool) {
	o.IsHost = &v
}

// GetEths returns the Eths field value if set, zero value otherwise.
func (o *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) GetEths() ConfigPutRequestSwitchpointSwitchpointNameObjectPropertiesEths {
	if o == nil || IsNil(o.Eths) {
		var ret ConfigPutRequestSwitchpointSwitchpointNameObjectPropertiesEths
		return ret
	}
	return *o.Eths
}

// GetEthsOk returns a tuple with the Eths field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) GetEthsOk() (*ConfigPutRequestSwitchpointSwitchpointNameObjectPropertiesEths, bool) {
	if o == nil || IsNil(o.Eths) {
		return nil, false
	}
	return o.Eths, true
}

// HasEths returns a boolean if a field has been set.
func (o *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) HasEths() bool {
	if o != nil && !IsNil(o.Eths) {
		return true
	}

	return false
}

// SetEths gets a reference to the given ConfigPutRequestSwitchpointSwitchpointNameObjectPropertiesEths and assigns it to the Eths field.
func (o *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) SetEths(v ConfigPutRequestSwitchpointSwitchpointNameObjectPropertiesEths) {
	o.Eths = &v
}

func (o ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.UserNotes) {
		toSerialize["user_notes"] = o.UserNotes
	}
	if !IsNil(o.ExpectedParentEndpoint) {
		toSerialize["expected_parent_endpoint"] = o.ExpectedParentEndpoint
	}
	if !IsNil(o.ExpectedParentEndpointRefType) {
		toSerialize["expected_parent_endpoint_ref_type_"] = o.ExpectedParentEndpointRefType
	}
	if o.NumberOfMultipoints.IsSet() {
		toSerialize["number_of_multipoints"] = o.NumberOfMultipoints.Get()
	}
	if !IsNil(o.Aggregate) {
		toSerialize["aggregate"] = o.Aggregate
	}
	if !IsNil(o.IsHost) {
		toSerialize["is_host"] = o.IsHost
	}
	if !IsNil(o.Eths) {
		toSerialize["eths"] = o.Eths
	}
	return toSerialize, nil
}

type NullableConfigPutRequestSwitchpointSwitchpointNameObjectProperties struct {
	value *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties
	isSet bool
}

func (v NullableConfigPutRequestSwitchpointSwitchpointNameObjectProperties) Get() *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties {
	return v.value
}

func (v *NullableConfigPutRequestSwitchpointSwitchpointNameObjectProperties) Set(val *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) {
	v.value = val
	v.isSet = true
}

func (v NullableConfigPutRequestSwitchpointSwitchpointNameObjectProperties) IsSet() bool {
	return v.isSet
}

func (v *NullableConfigPutRequestSwitchpointSwitchpointNameObjectProperties) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConfigPutRequestSwitchpointSwitchpointNameObjectProperties(val *ConfigPutRequestSwitchpointSwitchpointNameObjectProperties) *NullableConfigPutRequestSwitchpointSwitchpointNameObjectProperties {
	return &NullableConfigPutRequestSwitchpointSwitchpointNameObjectProperties{value: val, isSet: true}
}

func (v NullableConfigPutRequestSwitchpointSwitchpointNameObjectProperties) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConfigPutRequestSwitchpointSwitchpointNameObjectProperties) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


