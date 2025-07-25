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
	EthPortPaths []BundlesPutRequestEndpointBundleValueEthPortPathsInner `json:"eth_port_paths,omitempty"`
	UserServices []BundlesPutRequestEndpointBundleValueUserServicesInner `json:"user_services,omitempty"`
	ObjectProperties *BundlesPutRequestEndpointBundleValueObjectProperties `json:"object_properties,omitempty"`
	RgServices []BundlesPutRequestEndpointBundleValueRgServicesInner `json:"rg_services,omitempty"`
	// Enable object. It's highly recommended to set this value to true so that validation on the object will be ran.
	Enable *bool `json:"enable,omitempty"`
	// Voice Protocol: MGCP or SIP
	Protocol *string `json:"protocol,omitempty"`
	// Device Voice Settings for device
	DeviceVoiceSettings *string `json:"device_voice_settings,omitempty"`
	// Object type for device_voice_settings field
	DeviceVoiceSettingsRefType *string `json:"device_voice_settings_ref_type_,omitempty"`
	VoicePortProfilePaths []BundlesPutRequestEndpointBundleValueVoicePortProfilePathsInner `json:"voice_port_profile_paths,omitempty"`
}

// NewBundlesPatchRequestEndpointBundleValue instantiates a new BundlesPatchRequestEndpointBundleValue object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBundlesPatchRequestEndpointBundleValue() *BundlesPatchRequestEndpointBundleValue {
	this := BundlesPatchRequestEndpointBundleValue{}
	var name string = ""
	this.Name = &name
	var deviceSettings string = "eth_device_profile|(Device Settings)|"
	this.DeviceSettings = &deviceSettings
	var cliCommands string = ""
	this.CliCommands = &cliCommands
	var enable bool = false
	this.Enable = &enable
	var protocol string = "SIP"
	this.Protocol = &protocol
	var deviceVoiceSettings string = "voice_device_profile|(SIP Voice Device)|"
	this.DeviceVoiceSettings = &deviceVoiceSettings
	return &this
}

// NewBundlesPatchRequestEndpointBundleValueWithDefaults instantiates a new BundlesPatchRequestEndpointBundleValue object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBundlesPatchRequestEndpointBundleValueWithDefaults() *BundlesPatchRequestEndpointBundleValue {
	this := BundlesPatchRequestEndpointBundleValue{}
	var name string = ""
	this.Name = &name
	var deviceSettings string = "eth_device_profile|(Device Settings)|"
	this.DeviceSettings = &deviceSettings
	var cliCommands string = ""
	this.CliCommands = &cliCommands
	var enable bool = false
	this.Enable = &enable
	var protocol string = "SIP"
	this.Protocol = &protocol
	var deviceVoiceSettings string = "voice_device_profile|(SIP Voice Device)|"
	this.DeviceVoiceSettings = &deviceVoiceSettings
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
func (o *BundlesPatchRequestEndpointBundleValue) GetEthPortPaths() []BundlesPutRequestEndpointBundleValueEthPortPathsInner {
	if o == nil || IsNil(o.EthPortPaths) {
		var ret []BundlesPutRequestEndpointBundleValueEthPortPathsInner
		return ret
	}
	return o.EthPortPaths
}

// GetEthPortPathsOk returns a tuple with the EthPortPaths field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BundlesPatchRequestEndpointBundleValue) GetEthPortPathsOk() ([]BundlesPutRequestEndpointBundleValueEthPortPathsInner, bool) {
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

// SetEthPortPaths gets a reference to the given []BundlesPutRequestEndpointBundleValueEthPortPathsInner and assigns it to the EthPortPaths field.
func (o *BundlesPatchRequestEndpointBundleValue) SetEthPortPaths(v []BundlesPutRequestEndpointBundleValueEthPortPathsInner) {
	o.EthPortPaths = v
}

// GetUserServices returns the UserServices field value if set, zero value otherwise.
func (o *BundlesPatchRequestEndpointBundleValue) GetUserServices() []BundlesPutRequestEndpointBundleValueUserServicesInner {
	if o == nil || IsNil(o.UserServices) {
		var ret []BundlesPutRequestEndpointBundleValueUserServicesInner
		return ret
	}
	return o.UserServices
}

// GetUserServicesOk returns a tuple with the UserServices field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BundlesPatchRequestEndpointBundleValue) GetUserServicesOk() ([]BundlesPutRequestEndpointBundleValueUserServicesInner, bool) {
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

// SetUserServices gets a reference to the given []BundlesPutRequestEndpointBundleValueUserServicesInner and assigns it to the UserServices field.
func (o *BundlesPatchRequestEndpointBundleValue) SetUserServices(v []BundlesPutRequestEndpointBundleValueUserServicesInner) {
	o.UserServices = v
}

// GetObjectProperties returns the ObjectProperties field value if set, zero value otherwise.
func (o *BundlesPatchRequestEndpointBundleValue) GetObjectProperties() BundlesPutRequestEndpointBundleValueObjectProperties {
	if o == nil || IsNil(o.ObjectProperties) {
		var ret BundlesPutRequestEndpointBundleValueObjectProperties
		return ret
	}
	return *o.ObjectProperties
}

// GetObjectPropertiesOk returns a tuple with the ObjectProperties field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BundlesPatchRequestEndpointBundleValue) GetObjectPropertiesOk() (*BundlesPutRequestEndpointBundleValueObjectProperties, bool) {
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

// SetObjectProperties gets a reference to the given BundlesPutRequestEndpointBundleValueObjectProperties and assigns it to the ObjectProperties field.
func (o *BundlesPatchRequestEndpointBundleValue) SetObjectProperties(v BundlesPutRequestEndpointBundleValueObjectProperties) {
	o.ObjectProperties = &v
}

// GetRgServices returns the RgServices field value if set, zero value otherwise.
func (o *BundlesPatchRequestEndpointBundleValue) GetRgServices() []BundlesPutRequestEndpointBundleValueRgServicesInner {
	if o == nil || IsNil(o.RgServices) {
		var ret []BundlesPutRequestEndpointBundleValueRgServicesInner
		return ret
	}
	return o.RgServices
}

// GetRgServicesOk returns a tuple with the RgServices field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BundlesPatchRequestEndpointBundleValue) GetRgServicesOk() ([]BundlesPutRequestEndpointBundleValueRgServicesInner, bool) {
	if o == nil || IsNil(o.RgServices) {
		return nil, false
	}
	return o.RgServices, true
}

// HasRgServices returns a boolean if a field has been set.
func (o *BundlesPatchRequestEndpointBundleValue) HasRgServices() bool {
	if o != nil && !IsNil(o.RgServices) {
		return true
	}

	return false
}

// SetRgServices gets a reference to the given []BundlesPutRequestEndpointBundleValueRgServicesInner and assigns it to the RgServices field.
func (o *BundlesPatchRequestEndpointBundleValue) SetRgServices(v []BundlesPutRequestEndpointBundleValueRgServicesInner) {
	o.RgServices = v
}

// GetEnable returns the Enable field value if set, zero value otherwise.
func (o *BundlesPatchRequestEndpointBundleValue) GetEnable() bool {
	if o == nil || IsNil(o.Enable) {
		var ret bool
		return ret
	}
	return *o.Enable
}

// GetEnableOk returns a tuple with the Enable field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BundlesPatchRequestEndpointBundleValue) GetEnableOk() (*bool, bool) {
	if o == nil || IsNil(o.Enable) {
		return nil, false
	}
	return o.Enable, true
}

// HasEnable returns a boolean if a field has been set.
func (o *BundlesPatchRequestEndpointBundleValue) HasEnable() bool {
	if o != nil && !IsNil(o.Enable) {
		return true
	}

	return false
}

// SetEnable gets a reference to the given bool and assigns it to the Enable field.
func (o *BundlesPatchRequestEndpointBundleValue) SetEnable(v bool) {
	o.Enable = &v
}

// GetProtocol returns the Protocol field value if set, zero value otherwise.
func (o *BundlesPatchRequestEndpointBundleValue) GetProtocol() string {
	if o == nil || IsNil(o.Protocol) {
		var ret string
		return ret
	}
	return *o.Protocol
}

// GetProtocolOk returns a tuple with the Protocol field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BundlesPatchRequestEndpointBundleValue) GetProtocolOk() (*string, bool) {
	if o == nil || IsNil(o.Protocol) {
		return nil, false
	}
	return o.Protocol, true
}

// HasProtocol returns a boolean if a field has been set.
func (o *BundlesPatchRequestEndpointBundleValue) HasProtocol() bool {
	if o != nil && !IsNil(o.Protocol) {
		return true
	}

	return false
}

// SetProtocol gets a reference to the given string and assigns it to the Protocol field.
func (o *BundlesPatchRequestEndpointBundleValue) SetProtocol(v string) {
	o.Protocol = &v
}

// GetDeviceVoiceSettings returns the DeviceVoiceSettings field value if set, zero value otherwise.
func (o *BundlesPatchRequestEndpointBundleValue) GetDeviceVoiceSettings() string {
	if o == nil || IsNil(o.DeviceVoiceSettings) {
		var ret string
		return ret
	}
	return *o.DeviceVoiceSettings
}

// GetDeviceVoiceSettingsOk returns a tuple with the DeviceVoiceSettings field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BundlesPatchRequestEndpointBundleValue) GetDeviceVoiceSettingsOk() (*string, bool) {
	if o == nil || IsNil(o.DeviceVoiceSettings) {
		return nil, false
	}
	return o.DeviceVoiceSettings, true
}

// HasDeviceVoiceSettings returns a boolean if a field has been set.
func (o *BundlesPatchRequestEndpointBundleValue) HasDeviceVoiceSettings() bool {
	if o != nil && !IsNil(o.DeviceVoiceSettings) {
		return true
	}

	return false
}

// SetDeviceVoiceSettings gets a reference to the given string and assigns it to the DeviceVoiceSettings field.
func (o *BundlesPatchRequestEndpointBundleValue) SetDeviceVoiceSettings(v string) {
	o.DeviceVoiceSettings = &v
}

// GetDeviceVoiceSettingsRefType returns the DeviceVoiceSettingsRefType field value if set, zero value otherwise.
func (o *BundlesPatchRequestEndpointBundleValue) GetDeviceVoiceSettingsRefType() string {
	if o == nil || IsNil(o.DeviceVoiceSettingsRefType) {
		var ret string
		return ret
	}
	return *o.DeviceVoiceSettingsRefType
}

// GetDeviceVoiceSettingsRefTypeOk returns a tuple with the DeviceVoiceSettingsRefType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BundlesPatchRequestEndpointBundleValue) GetDeviceVoiceSettingsRefTypeOk() (*string, bool) {
	if o == nil || IsNil(o.DeviceVoiceSettingsRefType) {
		return nil, false
	}
	return o.DeviceVoiceSettingsRefType, true
}

// HasDeviceVoiceSettingsRefType returns a boolean if a field has been set.
func (o *BundlesPatchRequestEndpointBundleValue) HasDeviceVoiceSettingsRefType() bool {
	if o != nil && !IsNil(o.DeviceVoiceSettingsRefType) {
		return true
	}

	return false
}

// SetDeviceVoiceSettingsRefType gets a reference to the given string and assigns it to the DeviceVoiceSettingsRefType field.
func (o *BundlesPatchRequestEndpointBundleValue) SetDeviceVoiceSettingsRefType(v string) {
	o.DeviceVoiceSettingsRefType = &v
}

// GetVoicePortProfilePaths returns the VoicePortProfilePaths field value if set, zero value otherwise.
func (o *BundlesPatchRequestEndpointBundleValue) GetVoicePortProfilePaths() []BundlesPutRequestEndpointBundleValueVoicePortProfilePathsInner {
	if o == nil || IsNil(o.VoicePortProfilePaths) {
		var ret []BundlesPutRequestEndpointBundleValueVoicePortProfilePathsInner
		return ret
	}
	return o.VoicePortProfilePaths
}

// GetVoicePortProfilePathsOk returns a tuple with the VoicePortProfilePaths field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BundlesPatchRequestEndpointBundleValue) GetVoicePortProfilePathsOk() ([]BundlesPutRequestEndpointBundleValueVoicePortProfilePathsInner, bool) {
	if o == nil || IsNil(o.VoicePortProfilePaths) {
		return nil, false
	}
	return o.VoicePortProfilePaths, true
}

// HasVoicePortProfilePaths returns a boolean if a field has been set.
func (o *BundlesPatchRequestEndpointBundleValue) HasVoicePortProfilePaths() bool {
	if o != nil && !IsNil(o.VoicePortProfilePaths) {
		return true
	}

	return false
}

// SetVoicePortProfilePaths gets a reference to the given []BundlesPutRequestEndpointBundleValueVoicePortProfilePathsInner and assigns it to the VoicePortProfilePaths field.
func (o *BundlesPatchRequestEndpointBundleValue) SetVoicePortProfilePaths(v []BundlesPutRequestEndpointBundleValueVoicePortProfilePathsInner) {
	o.VoicePortProfilePaths = v
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
	if !IsNil(o.RgServices) {
		toSerialize["rg_services"] = o.RgServices
	}
	if !IsNil(o.Enable) {
		toSerialize["enable"] = o.Enable
	}
	if !IsNil(o.Protocol) {
		toSerialize["protocol"] = o.Protocol
	}
	if !IsNil(o.DeviceVoiceSettings) {
		toSerialize["device_voice_settings"] = o.DeviceVoiceSettings
	}
	if !IsNil(o.DeviceVoiceSettingsRefType) {
		toSerialize["device_voice_settings_ref_type_"] = o.DeviceVoiceSettingsRefType
	}
	if !IsNil(o.VoicePortProfilePaths) {
		toSerialize["voice_port_profile_paths"] = o.VoicePortProfilePaths
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


