/*
Verity API

This application demonstrates the usage of Verity API. 

API version: 2.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
	"bytes"
	"fmt"
)

// checks if the SwitchpointsUpgradePatchRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &SwitchpointsUpgradePatchRequest{}

// SwitchpointsUpgradePatchRequest struct for SwitchpointsUpgradePatchRequest
type SwitchpointsUpgradePatchRequest struct {
	// Version to upgrade to
	PackageVersion string `json:"package_version"`
	DeviceNames []string `json:"device_names"`
}

type _SwitchpointsUpgradePatchRequest SwitchpointsUpgradePatchRequest

// NewSwitchpointsUpgradePatchRequest instantiates a new SwitchpointsUpgradePatchRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewSwitchpointsUpgradePatchRequest(packageVersion string, deviceNames []string) *SwitchpointsUpgradePatchRequest {
	this := SwitchpointsUpgradePatchRequest{}
	this.PackageVersion = packageVersion
	this.DeviceNames = deviceNames
	return &this
}

// NewSwitchpointsUpgradePatchRequestWithDefaults instantiates a new SwitchpointsUpgradePatchRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewSwitchpointsUpgradePatchRequestWithDefaults() *SwitchpointsUpgradePatchRequest {
	this := SwitchpointsUpgradePatchRequest{}
	var packageVersion string = ""
	this.PackageVersion = packageVersion
	return &this
}

// GetPackageVersion returns the PackageVersion field value
func (o *SwitchpointsUpgradePatchRequest) GetPackageVersion() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.PackageVersion
}

// GetPackageVersionOk returns a tuple with the PackageVersion field value
// and a boolean to check if the value has been set.
func (o *SwitchpointsUpgradePatchRequest) GetPackageVersionOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.PackageVersion, true
}

// SetPackageVersion sets field value
func (o *SwitchpointsUpgradePatchRequest) SetPackageVersion(v string) {
	o.PackageVersion = v
}

// GetDeviceNames returns the DeviceNames field value
func (o *SwitchpointsUpgradePatchRequest) GetDeviceNames() []string {
	if o == nil {
		var ret []string
		return ret
	}

	return o.DeviceNames
}

// GetDeviceNamesOk returns a tuple with the DeviceNames field value
// and a boolean to check if the value has been set.
func (o *SwitchpointsUpgradePatchRequest) GetDeviceNamesOk() ([]string, bool) {
	if o == nil {
		return nil, false
	}
	return o.DeviceNames, true
}

// SetDeviceNames sets field value
func (o *SwitchpointsUpgradePatchRequest) SetDeviceNames(v []string) {
	o.DeviceNames = v
}

func (o SwitchpointsUpgradePatchRequest) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o SwitchpointsUpgradePatchRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["package_version"] = o.PackageVersion
	toSerialize["device_names"] = o.DeviceNames
	return toSerialize, nil
}

func (o *SwitchpointsUpgradePatchRequest) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"package_version",
		"device_names",
	}

	allProperties := make(map[string]interface{})

	err = json.Unmarshal(data, &allProperties)

	if err != nil {
		return err;
	}

	for _, requiredProperty := range(requiredProperties) {
		if _, exists := allProperties[requiredProperty]; !exists {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	varSwitchpointsUpgradePatchRequest := _SwitchpointsUpgradePatchRequest{}

	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	err = decoder.Decode(&varSwitchpointsUpgradePatchRequest)

	if err != nil {
		return err
	}

	*o = SwitchpointsUpgradePatchRequest(varSwitchpointsUpgradePatchRequest)

	return err
}

type NullableSwitchpointsUpgradePatchRequest struct {
	value *SwitchpointsUpgradePatchRequest
	isSet bool
}

func (v NullableSwitchpointsUpgradePatchRequest) Get() *SwitchpointsUpgradePatchRequest {
	return v.value
}

func (v *NullableSwitchpointsUpgradePatchRequest) Set(val *SwitchpointsUpgradePatchRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableSwitchpointsUpgradePatchRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableSwitchpointsUpgradePatchRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableSwitchpointsUpgradePatchRequest(val *SwitchpointsUpgradePatchRequest) *NullableSwitchpointsUpgradePatchRequest {
	return &NullableSwitchpointsUpgradePatchRequest{value: val, isSet: true}
}

func (v NullableSwitchpointsUpgradePatchRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableSwitchpointsUpgradePatchRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


