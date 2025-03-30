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

// checks if the ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner{}

// ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner struct for ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner
type ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner struct {
	// 1st Switchpoint for a Static Connection
	Endpoint1ForAStaticConnection *string `json:"endpoint1_for_a_static_connection,omitempty"`
	// Object type for endpoint1_for_a_static_connection field
	Endpoint1ForAStaticConnectionRefType *string `json:"endpoint1_for_a_static_connection_ref_type_,omitempty"`
	// 1st Port Number for a Static Connection
	Port1ForAStaticConnection *string `json:"port1_for_a_static_connection,omitempty"`
	// Object type for port1_for_a_static_connection field
	Port1ForAStaticConnectionRefType *string `json:"port1_for_a_static_connection_ref_type_,omitempty"`
	// 2nd Switchpoint for a Static Connection
	Endpoint2ForAStaticConnection *string `json:"endpoint2_for_a_static_connection,omitempty"`
	// Object type for endpoint2_for_a_static_connection field
	Endpoint2ForAStaticConnectionRefType *string `json:"endpoint2_for_a_static_connection_ref_type_,omitempty"`
	// 2nd Port Number for a Static Connection
	Port2ForAStaticConnection *string `json:"port2_for_a_static_connection,omitempty"`
	// Object type for port2_for_a_static_connection field
	Port2ForAStaticConnectionRefType *string `json:"port2_for_a_static_connection_ref_type_,omitempty"`
	// The index identifying the object. Zero if you want to add an object to the list.
	Index *int32 `json:"index,omitempty"`
}

// NewConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner instantiates a new ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner() *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner {
	this := ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner{}
	var endpoint1ForAStaticConnection string = ""
	this.Endpoint1ForAStaticConnection = &endpoint1ForAStaticConnection
	var port1ForAStaticConnection string = ""
	this.Port1ForAStaticConnection = &port1ForAStaticConnection
	var endpoint2ForAStaticConnection string = ""
	this.Endpoint2ForAStaticConnection = &endpoint2ForAStaticConnection
	var port2ForAStaticConnection string = ""
	this.Port2ForAStaticConnection = &port2ForAStaticConnection
	return &this
}

// NewConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInnerWithDefaults instantiates a new ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInnerWithDefaults() *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner {
	this := ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner{}
	var endpoint1ForAStaticConnection string = ""
	this.Endpoint1ForAStaticConnection = &endpoint1ForAStaticConnection
	var port1ForAStaticConnection string = ""
	this.Port1ForAStaticConnection = &port1ForAStaticConnection
	var endpoint2ForAStaticConnection string = ""
	this.Endpoint2ForAStaticConnection = &endpoint2ForAStaticConnection
	var port2ForAStaticConnection string = ""
	this.Port2ForAStaticConnection = &port2ForAStaticConnection
	return &this
}

// GetEndpoint1ForAStaticConnection returns the Endpoint1ForAStaticConnection field value if set, zero value otherwise.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) GetEndpoint1ForAStaticConnection() string {
	if o == nil || IsNil(o.Endpoint1ForAStaticConnection) {
		var ret string
		return ret
	}
	return *o.Endpoint1ForAStaticConnection
}

// GetEndpoint1ForAStaticConnectionOk returns a tuple with the Endpoint1ForAStaticConnection field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) GetEndpoint1ForAStaticConnectionOk() (*string, bool) {
	if o == nil || IsNil(o.Endpoint1ForAStaticConnection) {
		return nil, false
	}
	return o.Endpoint1ForAStaticConnection, true
}

// HasEndpoint1ForAStaticConnection returns a boolean if a field has been set.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) HasEndpoint1ForAStaticConnection() bool {
	if o != nil && !IsNil(o.Endpoint1ForAStaticConnection) {
		return true
	}

	return false
}

// SetEndpoint1ForAStaticConnection gets a reference to the given string and assigns it to the Endpoint1ForAStaticConnection field.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) SetEndpoint1ForAStaticConnection(v string) {
	o.Endpoint1ForAStaticConnection = &v
}

// GetEndpoint1ForAStaticConnectionRefType returns the Endpoint1ForAStaticConnectionRefType field value if set, zero value otherwise.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) GetEndpoint1ForAStaticConnectionRefType() string {
	if o == nil || IsNil(o.Endpoint1ForAStaticConnectionRefType) {
		var ret string
		return ret
	}
	return *o.Endpoint1ForAStaticConnectionRefType
}

// GetEndpoint1ForAStaticConnectionRefTypeOk returns a tuple with the Endpoint1ForAStaticConnectionRefType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) GetEndpoint1ForAStaticConnectionRefTypeOk() (*string, bool) {
	if o == nil || IsNil(o.Endpoint1ForAStaticConnectionRefType) {
		return nil, false
	}
	return o.Endpoint1ForAStaticConnectionRefType, true
}

// HasEndpoint1ForAStaticConnectionRefType returns a boolean if a field has been set.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) HasEndpoint1ForAStaticConnectionRefType() bool {
	if o != nil && !IsNil(o.Endpoint1ForAStaticConnectionRefType) {
		return true
	}

	return false
}

// SetEndpoint1ForAStaticConnectionRefType gets a reference to the given string and assigns it to the Endpoint1ForAStaticConnectionRefType field.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) SetEndpoint1ForAStaticConnectionRefType(v string) {
	o.Endpoint1ForAStaticConnectionRefType = &v
}

// GetPort1ForAStaticConnection returns the Port1ForAStaticConnection field value if set, zero value otherwise.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) GetPort1ForAStaticConnection() string {
	if o == nil || IsNil(o.Port1ForAStaticConnection) {
		var ret string
		return ret
	}
	return *o.Port1ForAStaticConnection
}

// GetPort1ForAStaticConnectionOk returns a tuple with the Port1ForAStaticConnection field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) GetPort1ForAStaticConnectionOk() (*string, bool) {
	if o == nil || IsNil(o.Port1ForAStaticConnection) {
		return nil, false
	}
	return o.Port1ForAStaticConnection, true
}

// HasPort1ForAStaticConnection returns a boolean if a field has been set.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) HasPort1ForAStaticConnection() bool {
	if o != nil && !IsNil(o.Port1ForAStaticConnection) {
		return true
	}

	return false
}

// SetPort1ForAStaticConnection gets a reference to the given string and assigns it to the Port1ForAStaticConnection field.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) SetPort1ForAStaticConnection(v string) {
	o.Port1ForAStaticConnection = &v
}

// GetPort1ForAStaticConnectionRefType returns the Port1ForAStaticConnectionRefType field value if set, zero value otherwise.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) GetPort1ForAStaticConnectionRefType() string {
	if o == nil || IsNil(o.Port1ForAStaticConnectionRefType) {
		var ret string
		return ret
	}
	return *o.Port1ForAStaticConnectionRefType
}

// GetPort1ForAStaticConnectionRefTypeOk returns a tuple with the Port1ForAStaticConnectionRefType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) GetPort1ForAStaticConnectionRefTypeOk() (*string, bool) {
	if o == nil || IsNil(o.Port1ForAStaticConnectionRefType) {
		return nil, false
	}
	return o.Port1ForAStaticConnectionRefType, true
}

// HasPort1ForAStaticConnectionRefType returns a boolean if a field has been set.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) HasPort1ForAStaticConnectionRefType() bool {
	if o != nil && !IsNil(o.Port1ForAStaticConnectionRefType) {
		return true
	}

	return false
}

// SetPort1ForAStaticConnectionRefType gets a reference to the given string and assigns it to the Port1ForAStaticConnectionRefType field.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) SetPort1ForAStaticConnectionRefType(v string) {
	o.Port1ForAStaticConnectionRefType = &v
}

// GetEndpoint2ForAStaticConnection returns the Endpoint2ForAStaticConnection field value if set, zero value otherwise.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) GetEndpoint2ForAStaticConnection() string {
	if o == nil || IsNil(o.Endpoint2ForAStaticConnection) {
		var ret string
		return ret
	}
	return *o.Endpoint2ForAStaticConnection
}

// GetEndpoint2ForAStaticConnectionOk returns a tuple with the Endpoint2ForAStaticConnection field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) GetEndpoint2ForAStaticConnectionOk() (*string, bool) {
	if o == nil || IsNil(o.Endpoint2ForAStaticConnection) {
		return nil, false
	}
	return o.Endpoint2ForAStaticConnection, true
}

// HasEndpoint2ForAStaticConnection returns a boolean if a field has been set.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) HasEndpoint2ForAStaticConnection() bool {
	if o != nil && !IsNil(o.Endpoint2ForAStaticConnection) {
		return true
	}

	return false
}

// SetEndpoint2ForAStaticConnection gets a reference to the given string and assigns it to the Endpoint2ForAStaticConnection field.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) SetEndpoint2ForAStaticConnection(v string) {
	o.Endpoint2ForAStaticConnection = &v
}

// GetEndpoint2ForAStaticConnectionRefType returns the Endpoint2ForAStaticConnectionRefType field value if set, zero value otherwise.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) GetEndpoint2ForAStaticConnectionRefType() string {
	if o == nil || IsNil(o.Endpoint2ForAStaticConnectionRefType) {
		var ret string
		return ret
	}
	return *o.Endpoint2ForAStaticConnectionRefType
}

// GetEndpoint2ForAStaticConnectionRefTypeOk returns a tuple with the Endpoint2ForAStaticConnectionRefType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) GetEndpoint2ForAStaticConnectionRefTypeOk() (*string, bool) {
	if o == nil || IsNil(o.Endpoint2ForAStaticConnectionRefType) {
		return nil, false
	}
	return o.Endpoint2ForAStaticConnectionRefType, true
}

// HasEndpoint2ForAStaticConnectionRefType returns a boolean if a field has been set.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) HasEndpoint2ForAStaticConnectionRefType() bool {
	if o != nil && !IsNil(o.Endpoint2ForAStaticConnectionRefType) {
		return true
	}

	return false
}

// SetEndpoint2ForAStaticConnectionRefType gets a reference to the given string and assigns it to the Endpoint2ForAStaticConnectionRefType field.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) SetEndpoint2ForAStaticConnectionRefType(v string) {
	o.Endpoint2ForAStaticConnectionRefType = &v
}

// GetPort2ForAStaticConnection returns the Port2ForAStaticConnection field value if set, zero value otherwise.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) GetPort2ForAStaticConnection() string {
	if o == nil || IsNil(o.Port2ForAStaticConnection) {
		var ret string
		return ret
	}
	return *o.Port2ForAStaticConnection
}

// GetPort2ForAStaticConnectionOk returns a tuple with the Port2ForAStaticConnection field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) GetPort2ForAStaticConnectionOk() (*string, bool) {
	if o == nil || IsNil(o.Port2ForAStaticConnection) {
		return nil, false
	}
	return o.Port2ForAStaticConnection, true
}

// HasPort2ForAStaticConnection returns a boolean if a field has been set.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) HasPort2ForAStaticConnection() bool {
	if o != nil && !IsNil(o.Port2ForAStaticConnection) {
		return true
	}

	return false
}

// SetPort2ForAStaticConnection gets a reference to the given string and assigns it to the Port2ForAStaticConnection field.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) SetPort2ForAStaticConnection(v string) {
	o.Port2ForAStaticConnection = &v
}

// GetPort2ForAStaticConnectionRefType returns the Port2ForAStaticConnectionRefType field value if set, zero value otherwise.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) GetPort2ForAStaticConnectionRefType() string {
	if o == nil || IsNil(o.Port2ForAStaticConnectionRefType) {
		var ret string
		return ret
	}
	return *o.Port2ForAStaticConnectionRefType
}

// GetPort2ForAStaticConnectionRefTypeOk returns a tuple with the Port2ForAStaticConnectionRefType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) GetPort2ForAStaticConnectionRefTypeOk() (*string, bool) {
	if o == nil || IsNil(o.Port2ForAStaticConnectionRefType) {
		return nil, false
	}
	return o.Port2ForAStaticConnectionRefType, true
}

// HasPort2ForAStaticConnectionRefType returns a boolean if a field has been set.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) HasPort2ForAStaticConnectionRefType() bool {
	if o != nil && !IsNil(o.Port2ForAStaticConnectionRefType) {
		return true
	}

	return false
}

// SetPort2ForAStaticConnectionRefType gets a reference to the given string and assigns it to the Port2ForAStaticConnectionRefType field.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) SetPort2ForAStaticConnectionRefType(v string) {
	o.Port2ForAStaticConnectionRefType = &v
}

// GetIndex returns the Index field value if set, zero value otherwise.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) GetIndex() int32 {
	if o == nil || IsNil(o.Index) {
		var ret int32
		return ret
	}
	return *o.Index
}

// GetIndexOk returns a tuple with the Index field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) GetIndexOk() (*int32, bool) {
	if o == nil || IsNil(o.Index) {
		return nil, false
	}
	return o.Index, true
}

// HasIndex returns a boolean if a field has been set.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) HasIndex() bool {
	if o != nil && !IsNil(o.Index) {
		return true
	}

	return false
}

// SetIndex gets a reference to the given int32 and assigns it to the Index field.
func (o *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) SetIndex(v int32) {
	o.Index = &v
}

func (o ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Endpoint1ForAStaticConnection) {
		toSerialize["endpoint1_for_a_static_connection"] = o.Endpoint1ForAStaticConnection
	}
	if !IsNil(o.Endpoint1ForAStaticConnectionRefType) {
		toSerialize["endpoint1_for_a_static_connection_ref_type_"] = o.Endpoint1ForAStaticConnectionRefType
	}
	if !IsNil(o.Port1ForAStaticConnection) {
		toSerialize["port1_for_a_static_connection"] = o.Port1ForAStaticConnection
	}
	if !IsNil(o.Port1ForAStaticConnectionRefType) {
		toSerialize["port1_for_a_static_connection_ref_type_"] = o.Port1ForAStaticConnectionRefType
	}
	if !IsNil(o.Endpoint2ForAStaticConnection) {
		toSerialize["endpoint2_for_a_static_connection"] = o.Endpoint2ForAStaticConnection
	}
	if !IsNil(o.Endpoint2ForAStaticConnectionRefType) {
		toSerialize["endpoint2_for_a_static_connection_ref_type_"] = o.Endpoint2ForAStaticConnectionRefType
	}
	if !IsNil(o.Port2ForAStaticConnection) {
		toSerialize["port2_for_a_static_connection"] = o.Port2ForAStaticConnection
	}
	if !IsNil(o.Port2ForAStaticConnectionRefType) {
		toSerialize["port2_for_a_static_connection_ref_type_"] = o.Port2ForAStaticConnectionRefType
	}
	if !IsNil(o.Index) {
		toSerialize["index"] = o.Index
	}
	return toSerialize, nil
}

type NullableConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner struct {
	value *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner
	isSet bool
}

func (v NullableConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) Get() *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner {
	return v.value
}

func (v *NullableConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) Set(val *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) {
	v.value = val
	v.isSet = true
}

func (v NullableConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) IsSet() bool {
	return v.isSet
}

func (v *NullableConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner(val *ConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) *NullableConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner {
	return &NullableConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner{value: val, isSet: true}
}

func (v NullableConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConfigPutRequestStaticConnectionsStaticConnectionsNameConnectionsInner) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


