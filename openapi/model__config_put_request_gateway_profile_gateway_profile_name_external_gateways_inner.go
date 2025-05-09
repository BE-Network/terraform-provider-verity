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

// checks if the ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner{}

// ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner struct for ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner
type ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner struct {
	// Enable row
	Enable *bool `json:"enable,omitempty"`
	// BGP Gateway referenced for port profile
	Gateway *string `json:"gateway,omitempty"`
	// Object type for gateway field
	GatewayRefType *string `json:"gateway_ref_type_,omitempty"`
	// Source address on the port if untagged or on the VLAN if tagged used for the outgoing BGP session 
	SourceIpMask *string `json:"source_ip_mask,omitempty"`
	// Setting for paired switches only. Flag indicating that this gateway is a peer gateway. For each gateway profile referencing a BGP session on a member of a leaf pair, the peer should have a gateway profile entry indicating the IP address for the peers gateway.
	PeerGw *bool `json:"peer_gw,omitempty"`
	// The index identifying the object. Zero if you want to add an object to the list.
	Index *int32 `json:"index,omitempty"`
}

// NewConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner instantiates a new ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner() *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner {
	this := ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner{}
	var enable bool = false
	this.Enable = &enable
	var gateway string = ""
	this.Gateway = &gateway
	var sourceIpMask string = ""
	this.SourceIpMask = &sourceIpMask
	var peerGw bool = false
	this.PeerGw = &peerGw
	return &this
}

// NewConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInnerWithDefaults instantiates a new ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInnerWithDefaults() *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner {
	this := ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner{}
	var enable bool = false
	this.Enable = &enable
	var gateway string = ""
	this.Gateway = &gateway
	var sourceIpMask string = ""
	this.SourceIpMask = &sourceIpMask
	var peerGw bool = false
	this.PeerGw = &peerGw
	return &this
}

// GetEnable returns the Enable field value if set, zero value otherwise.
func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) GetEnable() bool {
	if o == nil || IsNil(o.Enable) {
		var ret bool
		return ret
	}
	return *o.Enable
}

// GetEnableOk returns a tuple with the Enable field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) GetEnableOk() (*bool, bool) {
	if o == nil || IsNil(o.Enable) {
		return nil, false
	}
	return o.Enable, true
}

// HasEnable returns a boolean if a field has been set.
func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) HasEnable() bool {
	if o != nil && !IsNil(o.Enable) {
		return true
	}

	return false
}

// SetEnable gets a reference to the given bool and assigns it to the Enable field.
func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) SetEnable(v bool) {
	o.Enable = &v
}

// GetGateway returns the Gateway field value if set, zero value otherwise.
func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) GetGateway() string {
	if o == nil || IsNil(o.Gateway) {
		var ret string
		return ret
	}
	return *o.Gateway
}

// GetGatewayOk returns a tuple with the Gateway field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) GetGatewayOk() (*string, bool) {
	if o == nil || IsNil(o.Gateway) {
		return nil, false
	}
	return o.Gateway, true
}

// HasGateway returns a boolean if a field has been set.
func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) HasGateway() bool {
	if o != nil && !IsNil(o.Gateway) {
		return true
	}

	return false
}

// SetGateway gets a reference to the given string and assigns it to the Gateway field.
func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) SetGateway(v string) {
	o.Gateway = &v
}

// GetGatewayRefType returns the GatewayRefType field value if set, zero value otherwise.
func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) GetGatewayRefType() string {
	if o == nil || IsNil(o.GatewayRefType) {
		var ret string
		return ret
	}
	return *o.GatewayRefType
}

// GetGatewayRefTypeOk returns a tuple with the GatewayRefType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) GetGatewayRefTypeOk() (*string, bool) {
	if o == nil || IsNil(o.GatewayRefType) {
		return nil, false
	}
	return o.GatewayRefType, true
}

// HasGatewayRefType returns a boolean if a field has been set.
func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) HasGatewayRefType() bool {
	if o != nil && !IsNil(o.GatewayRefType) {
		return true
	}

	return false
}

// SetGatewayRefType gets a reference to the given string and assigns it to the GatewayRefType field.
func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) SetGatewayRefType(v string) {
	o.GatewayRefType = &v
}

// GetSourceIpMask returns the SourceIpMask field value if set, zero value otherwise.
func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) GetSourceIpMask() string {
	if o == nil || IsNil(o.SourceIpMask) {
		var ret string
		return ret
	}
	return *o.SourceIpMask
}

// GetSourceIpMaskOk returns a tuple with the SourceIpMask field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) GetSourceIpMaskOk() (*string, bool) {
	if o == nil || IsNil(o.SourceIpMask) {
		return nil, false
	}
	return o.SourceIpMask, true
}

// HasSourceIpMask returns a boolean if a field has been set.
func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) HasSourceIpMask() bool {
	if o != nil && !IsNil(o.SourceIpMask) {
		return true
	}

	return false
}

// SetSourceIpMask gets a reference to the given string and assigns it to the SourceIpMask field.
func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) SetSourceIpMask(v string) {
	o.SourceIpMask = &v
}

// GetPeerGw returns the PeerGw field value if set, zero value otherwise.
func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) GetPeerGw() bool {
	if o == nil || IsNil(o.PeerGw) {
		var ret bool
		return ret
	}
	return *o.PeerGw
}

// GetPeerGwOk returns a tuple with the PeerGw field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) GetPeerGwOk() (*bool, bool) {
	if o == nil || IsNil(o.PeerGw) {
		return nil, false
	}
	return o.PeerGw, true
}

// HasPeerGw returns a boolean if a field has been set.
func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) HasPeerGw() bool {
	if o != nil && !IsNil(o.PeerGw) {
		return true
	}

	return false
}

// SetPeerGw gets a reference to the given bool and assigns it to the PeerGw field.
func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) SetPeerGw(v bool) {
	o.PeerGw = &v
}

// GetIndex returns the Index field value if set, zero value otherwise.
func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) GetIndex() int32 {
	if o == nil || IsNil(o.Index) {
		var ret int32
		return ret
	}
	return *o.Index
}

// GetIndexOk returns a tuple with the Index field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) GetIndexOk() (*int32, bool) {
	if o == nil || IsNil(o.Index) {
		return nil, false
	}
	return o.Index, true
}

// HasIndex returns a boolean if a field has been set.
func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) HasIndex() bool {
	if o != nil && !IsNil(o.Index) {
		return true
	}

	return false
}

// SetIndex gets a reference to the given int32 and assigns it to the Index field.
func (o *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) SetIndex(v int32) {
	o.Index = &v
}

func (o ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Enable) {
		toSerialize["enable"] = o.Enable
	}
	if !IsNil(o.Gateway) {
		toSerialize["gateway"] = o.Gateway
	}
	if !IsNil(o.GatewayRefType) {
		toSerialize["gateway_ref_type_"] = o.GatewayRefType
	}
	if !IsNil(o.SourceIpMask) {
		toSerialize["source_ip_mask"] = o.SourceIpMask
	}
	if !IsNil(o.PeerGw) {
		toSerialize["peer_gw"] = o.PeerGw
	}
	if !IsNil(o.Index) {
		toSerialize["index"] = o.Index
	}
	return toSerialize, nil
}

type NullableConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner struct {
	value *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner
	isSet bool
}

func (v NullableConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) Get() *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner {
	return v.value
}

func (v *NullableConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) Set(val *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) {
	v.value = val
	v.isSet = true
}

func (v NullableConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) IsSet() bool {
	return v.isSet
}

func (v *NullableConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner(val *ConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) *NullableConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner {
	return &NullableConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner{value: val, isSet: true}
}

func (v NullableConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConfigPutRequestGatewayProfileGatewayProfileNameExternalGatewaysInner) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


