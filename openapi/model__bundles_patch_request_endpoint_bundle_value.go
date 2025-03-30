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

// checks if the BundlesPatchRequestEndpointBundleValue type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &BundlesPatchRequestEndpointBundleValue{}

// BundlesPatchRequestEndpointBundleValue struct for BundlesPatchRequestEndpointBundleValue
type BundlesPatchRequestEndpointBundleValue struct {
	// Object Name. Must be unique.
	Name *string `json:"name,omitempty"`
	// Device Settings for device
	DeviceSettings *string `json:"device_settings,omitempty"`
	// Object type for device_settings field
	DeviceSettingsRefType *string `json:"device_settings_ref_type_,omitempty"`
	// CLI Commands
	CliCommands *string `json:"cli_commands,omitempty"`
	EthPortPaths []BundlesPatchRequestEndpointBundleValueEthPortPathsInner `json:"eth_port_paths,omitempty"`
	UserServices []BundlesPatchRequestEndpointBundleValueUserServicesInner `json:"user_services,omitempty"`
	ObjectProperties *BundlesPatchRequestEndpointBundleValueObjectProperties `json:"object_properties,omitempty"`
}

// NewBundlesPatchRequestEndpointBundleValue instantiates a new BundlesPatchRequestEndpointBundleValue object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBundlesPatchRequestEndpointBundleValue() *BundlesPatchRequestEndpointBundleValue {
	this := BundlesPatchRequestEndpointBundleValue{}
	var name string = ""
	this.Name = &name
	var deviceSettings string = "eth_device_profile|(Default)|"
	this.DeviceSettings = &deviceSettings
	var cliCommands string = ""
	this.CliCommands = &cliCommands
	return &this
}

// NewBundlesPatchRequestEndpointBundleValueWithDefaults instantiates a new BundlesPatchRequestEndpointBundleValue object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBundlesPatchRequestEndpointBundleValueWithDefaults() *BundlesPatchRequestEndpointBundleValue {
	this := BundlesPatchRequestEndpointBundleValue{}
	var name string = ""
	this.Name = &name
	var deviceSettings string = "eth_device_profile|(Default)|"
	this.DeviceSettings = &deviceSettings
	var cliCommands string = ""
	this.CliCommands = &cliCommands
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *BundlesPatchRequestEndpointBundleValue) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BundlesPatchRequestEndpointBundleValue) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *BundlesPatchRequestEndpointBundleValue) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *BundlesPatchRequestEndpointBundleValue) SetName(v string) {
	o.Name = &v
}

// GetDeviceSettings returns the DeviceSettings field value if set, zero value otherwise.
func (o *BundlesPatchRequestEndpointBundleValue) GetDeviceSettings() string {
	if o == nil || IsNil(o.DeviceSettings) {
		var ret string
		return ret
	}
	return *o.DeviceSettings
}

// GetDeviceSettingsOk returns a tuple with the DeviceSettings field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BundlesPatchRequestEndpointBundleValue) GetDeviceSettingsOk() (*string, bool) {
	if o == nil || IsNil(o.DeviceSettings) {
		return nil, false
	}
	return o.DeviceSettings, true
}

// HasDeviceSettings returns a boolean if a field has been set.
func (o *BundlesPatchRequestEndpointBundleValue) HasDeviceSettings() bool {
	if o != nil && !IsNil(o.DeviceSettings) {
		return true
	}

	return false
}

// SetDeviceSettings gets a reference to the given string and assigns it to the DeviceSettings field.
func (o *BundlesPatchRequestEndpointBundleValue) SetDeviceSettings(v string) {
	o.DeviceSettings = &v
}

// GetDeviceSettingsRefType returns the DeviceSettingsRefType field value if set, zero value otherwise.
func (o *BundlesPatchRequestEndpointBundleValue) GetDeviceSettingsRefType() string {
	if o == nil || IsNil(o.DeviceSettingsRefType) {
		var ret string
		return ret
	}
	return *o.DeviceSettingsRefType
}

// GetDeviceSettingsRefTypeOk returns a tuple with the DeviceSettingsRefType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BundlesPatchRequestEndpointBundleValue) GetDeviceSettingsRefTypeOk() (*string, bool) {
	if o == nil || IsNil(o.DeviceSettingsRefType) {
		return nil, false
	}
	return o.DeviceSettingsRefType, true
}

// HasDeviceSettingsRefType returns a boolean if a field has been set.
func (o *BundlesPatchRequestEndpointBundleValue) HasDeviceSettingsRefType() bool {
	if o != nil && !IsNil(o.DeviceSettingsRefType) {
		return true
	}

	return false
}

// SetDeviceSettingsRefType gets a reference to the given string and assigns it to the DeviceSettingsRefType field.
func (o *BundlesPatchRequestEndpointBundleValue) SetDeviceSettingsRefType(v string) {
	o.DeviceSettingsRefType = &v
}

// GetCliCommands returns the CliCommands field value if set, zero value otherwise.
func (o *BundlesPatchRequestEndpointBundleValue) GetCliCommands() string {
	if o == nil || IsNil(o.CliCommands) {
		var ret string
		return ret
	}
	return *o.CliCommands
}

// GetCliCommandsOk returns a tuple with the CliCommands field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BundlesPatchRequestEndpointBundleValue) GetCliCommandsOk() (*string, bool) {
	if o == nil || IsNil(o.CliCommands) {
		return nil, false
	}
	return o.CliCommands, true
}

// HasCliCommands returns a boolean if a field has been set.
func (o *BundlesPatchRequestEndpointBundleValue) HasCliCommands() bool {
	if o != nil && !IsNil(o.CliCommands) {
		return true
	}

	return false
}

// SetCliCommands gets a reference to the given string and assigns it to the CliCommands field.
func (o *BundlesPatchRequestEndpointBundleValue) SetCliCommands(v string) {
	o.CliCommands = &v
}

// GetEthPortPaths returns the EthPortPaths field value if set, zero value otherwise.
func (o *BundlesPatchRequestEndpointBundleValue) GetEthPortPaths() []BundlesPatchRequestEndpointBundleValueEthPortPathsInner {
	if o == nil || IsNil(o.EthPortPaths) {
		var ret []BundlesPatchRequestEndpointBundleValueEthPortPathsInner
		return ret
	}
	return o.EthPortPaths
}

// GetEthPortPathsOk returns a tuple with the EthPortPaths field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BundlesPatchRequestEndpointBundleValue) GetEthPortPathsOk() ([]BundlesPatchRequestEndpointBundleValueEthPortPathsInner, bool) {
	if o == nil || IsNil(o.EthPortPaths) {
		return nil, false
	}
	return o.EthPortPaths, true
}

// HasEthPortPaths returns a boolean if a field has been set.
func (o *BundlesPatchRequestEndpointBundleValue) HasEthPortPaths() bool {
	if o != nil && !IsNil(o.EthPortPaths) {
		return true
	}

	return false
}

// SetEthPortPaths gets a reference to the given []BundlesPatchRequestEndpointBundleValueEthPortPathsInner and assigns it to the EthPortPaths field.
func (o *BundlesPatchRequestEndpointBundleValue) SetEthPortPaths(v []BundlesPatchRequestEndpointBundleValueEthPortPathsInner) {
	o.EthPortPaths = v
}

// GetUserServices returns the UserServices field value if set, zero value otherwise.
func (o *BundlesPatchRequestEndpointBundleValue) GetUserServices() []BundlesPatchRequestEndpointBundleValueUserServicesInner {
	if o == nil || IsNil(o.UserServices) {
		var ret []BundlesPatchRequestEndpointBundleValueUserServicesInner
		return ret
	}
	return o.UserServices
}

// GetUserServicesOk returns a tuple with the UserServices field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BundlesPatchRequestEndpointBundleValue) GetUserServicesOk() ([]BundlesPatchRequestEndpointBundleValueUserServicesInner, bool) {
	if o == nil || IsNil(o.UserServices) {
		return nil, false
	}
	return o.UserServices, true
}

// HasUserServices returns a boolean if a field has been set.
func (o *BundlesPatchRequestEndpointBundleValue) HasUserServices() bool {
	if o != nil && !IsNil(o.UserServices) {
		return true
	}

	return false
}

// SetUserServices gets a reference to the given []BundlesPatchRequestEndpointBundleValueUserServicesInner and assigns it to the UserServices field.
func (o *BundlesPatchRequestEndpointBundleValue) SetUserServices(v []BundlesPatchRequestEndpointBundleValueUserServicesInner) {
	o.UserServices = v
}

// GetObjectProperties returns the ObjectProperties field value if set, zero value otherwise.
func (o *BundlesPatchRequestEndpointBundleValue) GetObjectProperties() BundlesPatchRequestEndpointBundleValueObjectProperties {
	if o == nil || IsNil(o.ObjectProperties) {
		var ret BundlesPatchRequestEndpointBundleValueObjectProperties
		return ret
	}
	return *o.ObjectProperties
}

// GetObjectPropertiesOk returns a tuple with the ObjectProperties field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BundlesPatchRequestEndpointBundleValue) GetObjectPropertiesOk() (*BundlesPatchRequestEndpointBundleValueObjectProperties, bool) {
	if o == nil || IsNil(o.ObjectProperties) {
		return nil, false
	}
	return o.ObjectProperties, true
}

// HasObjectProperties returns a boolean if a field has been set.
func (o *BundlesPatchRequestEndpointBundleValue) HasObjectProperties() bool {
	if o != nil && !IsNil(o.ObjectProperties) {
		return true
	}

	return false
}

// SetObjectProperties gets a reference to the given BundlesPatchRequestEndpointBundleValueObjectProperties and assigns it to the ObjectProperties field.
func (o *BundlesPatchRequestEndpointBundleValue) SetObjectProperties(v BundlesPatchRequestEndpointBundleValueObjectProperties) {
	o.ObjectProperties = &v
}

func (o BundlesPatchRequestEndpointBundleValue) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o BundlesPatchRequestEndpointBundleValue) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.DeviceSettings) {
		toSerialize["device_settings"] = o.DeviceSettings
	}
	if !IsNil(o.DeviceSettingsRefType) {
		toSerialize["device_settings_ref_type_"] = o.DeviceSettingsRefType
	}
	if !IsNil(o.CliCommands) {
		toSerialize["cli_commands"] = o.CliCommands
	}
	if !IsNil(o.EthPortPaths) {
		toSerialize["eth_port_paths"] = o.EthPortPaths
	}
	if !IsNil(o.UserServices) {
		toSerialize["user_services"] = o.UserServices
	}
	if !IsNil(o.ObjectProperties) {
		toSerialize["object_properties"] = o.ObjectProperties
	}
	return toSerialize, nil
}

type NullableBundlesPatchRequestEndpointBundleValue struct {
	value *BundlesPatchRequestEndpointBundleValue
	isSet bool
}

func (v NullableBundlesPatchRequestEndpointBundleValue) Get() *BundlesPatchRequestEndpointBundleValue {
	return v.value
}

func (v *NullableBundlesPatchRequestEndpointBundleValue) Set(val *BundlesPatchRequestEndpointBundleValue) {
	v.value = val
	v.isSet = true
}

func (v NullableBundlesPatchRequestEndpointBundleValue) IsSet() bool {
	return v.isSet
}

func (v *NullableBundlesPatchRequestEndpointBundleValue) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBundlesPatchRequestEndpointBundleValue(val *BundlesPatchRequestEndpointBundleValue) *NullableBundlesPatchRequestEndpointBundleValue {
	return &NullableBundlesPatchRequestEndpointBundleValue{value: val, isSet: true}
}

func (v NullableBundlesPatchRequestEndpointBundleValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBundlesPatchRequestEndpointBundleValue) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


