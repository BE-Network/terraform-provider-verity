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

// checks if the ConfigPutRequestEndpointEndpointName type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ConfigPutRequestEndpointEndpointName{}

// ConfigPutRequestEndpointEndpointName struct for ConfigPutRequestEndpointEndpointName
type ConfigPutRequestEndpointEndpointName struct {
	// Object Name. Must be unique.
	Name *string `json:"name,omitempty"`
	// Device Serial Number
	DeviceSerialNumber *string `json:"device_serial_number,omitempty"`
	// Connected Bundle
	ConnectedBundle *string `json:"connected_bundle,omitempty"`
	// Object type for connected_bundle field
	ConnectedBundleRefType *string `json:"connected_bundle_ref_type_,omitempty"`
	// When Read Only Mode is checked, vNetC will perform all functions except writing database updates to the target hardware
	ReadOnlyMode *bool `json:"read_only_mode,omitempty"`
	// Permission lock
	Locked *bool `json:"locked,omitempty"`
	// Disabled Ports It's a comma separated list of ports to disable.
	DisabledPorts *string `json:"disabled_ports,omitempty"`
	// For Switch Endpoints. Denotes a Switch that is Fabric rather than an Edge Device
	IsFabric *bool `json:"is_fabric,omitempty"`
	// For Switch Endpoints. Denotes a Switch is managed out of band via the management port
	OutOfBandManagement *bool `json:"out_of_band_management,omitempty"`
	Badges []ConfigPutRequestSwitchpointSwitchpointNameBadgesInner `json:"badges,omitempty"`
	Children []ConfigPutRequestEndpointEndpointNameChildrenInner `json:"children,omitempty"`
	TrafficMirrors []ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner `json:"traffic_mirrors,omitempty"`
	Eths []ConfigPutRequestEndpointEndpointNameEthsInner `json:"eths,omitempty"`
	ObjectProperties *ConfigPutRequestEndpointEndpointNameObjectProperties `json:"object_properties,omitempty"`
}

// NewConfigPutRequestEndpointEndpointName instantiates a new ConfigPutRequestEndpointEndpointName object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConfigPutRequestEndpointEndpointName() *ConfigPutRequestEndpointEndpointName {
	this := ConfigPutRequestEndpointEndpointName{}
	var name string = ""
	this.Name = &name
	var deviceSerialNumber string = ""
	this.DeviceSerialNumber = &deviceSerialNumber
	var connectedBundle string = ""
	this.ConnectedBundle = &connectedBundle
	var readOnlyMode bool = false
	this.ReadOnlyMode = &readOnlyMode
	var locked bool = false
	this.Locked = &locked
	var disabledPorts string = ""
	this.DisabledPorts = &disabledPorts
	var isFabric bool = false
	this.IsFabric = &isFabric
	var outOfBandManagement bool = false
	this.OutOfBandManagement = &outOfBandManagement
	return &this
}

// NewConfigPutRequestEndpointEndpointNameWithDefaults instantiates a new ConfigPutRequestEndpointEndpointName object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConfigPutRequestEndpointEndpointNameWithDefaults() *ConfigPutRequestEndpointEndpointName {
	this := ConfigPutRequestEndpointEndpointName{}
	var name string = ""
	this.Name = &name
	var deviceSerialNumber string = ""
	this.DeviceSerialNumber = &deviceSerialNumber
	var connectedBundle string = ""
	this.ConnectedBundle = &connectedBundle
	var readOnlyMode bool = false
	this.ReadOnlyMode = &readOnlyMode
	var locked bool = false
	this.Locked = &locked
	var disabledPorts string = ""
	this.DisabledPorts = &disabledPorts
	var isFabric bool = false
	this.IsFabric = &isFabric
	var outOfBandManagement bool = false
	this.OutOfBandManagement = &outOfBandManagement
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *ConfigPutRequestEndpointEndpointName) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEndpointEndpointName) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *ConfigPutRequestEndpointEndpointName) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *ConfigPutRequestEndpointEndpointName) SetName(v string) {
	o.Name = &v
}

// GetDeviceSerialNumber returns the DeviceSerialNumber field value if set, zero value otherwise.
func (o *ConfigPutRequestEndpointEndpointName) GetDeviceSerialNumber() string {
	if o == nil || IsNil(o.DeviceSerialNumber) {
		var ret string
		return ret
	}
	return *o.DeviceSerialNumber
}

// GetDeviceSerialNumberOk returns a tuple with the DeviceSerialNumber field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEndpointEndpointName) GetDeviceSerialNumberOk() (*string, bool) {
	if o == nil || IsNil(o.DeviceSerialNumber) {
		return nil, false
	}
	return o.DeviceSerialNumber, true
}

// HasDeviceSerialNumber returns a boolean if a field has been set.
func (o *ConfigPutRequestEndpointEndpointName) HasDeviceSerialNumber() bool {
	if o != nil && !IsNil(o.DeviceSerialNumber) {
		return true
	}

	return false
}

// SetDeviceSerialNumber gets a reference to the given string and assigns it to the DeviceSerialNumber field.
func (o *ConfigPutRequestEndpointEndpointName) SetDeviceSerialNumber(v string) {
	o.DeviceSerialNumber = &v
}

// GetConnectedBundle returns the ConnectedBundle field value if set, zero value otherwise.
func (o *ConfigPutRequestEndpointEndpointName) GetConnectedBundle() string {
	if o == nil || IsNil(o.ConnectedBundle) {
		var ret string
		return ret
	}
	return *o.ConnectedBundle
}

// GetConnectedBundleOk returns a tuple with the ConnectedBundle field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEndpointEndpointName) GetConnectedBundleOk() (*string, bool) {
	if o == nil || IsNil(o.ConnectedBundle) {
		return nil, false
	}
	return o.ConnectedBundle, true
}

// HasConnectedBundle returns a boolean if a field has been set.
func (o *ConfigPutRequestEndpointEndpointName) HasConnectedBundle() bool {
	if o != nil && !IsNil(o.ConnectedBundle) {
		return true
	}

	return false
}

// SetConnectedBundle gets a reference to the given string and assigns it to the ConnectedBundle field.
func (o *ConfigPutRequestEndpointEndpointName) SetConnectedBundle(v string) {
	o.ConnectedBundle = &v
}

// GetConnectedBundleRefType returns the ConnectedBundleRefType field value if set, zero value otherwise.
func (o *ConfigPutRequestEndpointEndpointName) GetConnectedBundleRefType() string {
	if o == nil || IsNil(o.ConnectedBundleRefType) {
		var ret string
		return ret
	}
	return *o.ConnectedBundleRefType
}

// GetConnectedBundleRefTypeOk returns a tuple with the ConnectedBundleRefType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEndpointEndpointName) GetConnectedBundleRefTypeOk() (*string, bool) {
	if o == nil || IsNil(o.ConnectedBundleRefType) {
		return nil, false
	}
	return o.ConnectedBundleRefType, true
}

// HasConnectedBundleRefType returns a boolean if a field has been set.
func (o *ConfigPutRequestEndpointEndpointName) HasConnectedBundleRefType() bool {
	if o != nil && !IsNil(o.ConnectedBundleRefType) {
		return true
	}

	return false
}

// SetConnectedBundleRefType gets a reference to the given string and assigns it to the ConnectedBundleRefType field.
func (o *ConfigPutRequestEndpointEndpointName) SetConnectedBundleRefType(v string) {
	o.ConnectedBundleRefType = &v
}

// GetReadOnlyMode returns the ReadOnlyMode field value if set, zero value otherwise.
func (o *ConfigPutRequestEndpointEndpointName) GetReadOnlyMode() bool {
	if o == nil || IsNil(o.ReadOnlyMode) {
		var ret bool
		return ret
	}
	return *o.ReadOnlyMode
}

// GetReadOnlyModeOk returns a tuple with the ReadOnlyMode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEndpointEndpointName) GetReadOnlyModeOk() (*bool, bool) {
	if o == nil || IsNil(o.ReadOnlyMode) {
		return nil, false
	}
	return o.ReadOnlyMode, true
}

// HasReadOnlyMode returns a boolean if a field has been set.
func (o *ConfigPutRequestEndpointEndpointName) HasReadOnlyMode() bool {
	if o != nil && !IsNil(o.ReadOnlyMode) {
		return true
	}

	return false
}

// SetReadOnlyMode gets a reference to the given bool and assigns it to the ReadOnlyMode field.
func (o *ConfigPutRequestEndpointEndpointName) SetReadOnlyMode(v bool) {
	o.ReadOnlyMode = &v
}

// GetLocked returns the Locked field value if set, zero value otherwise.
func (o *ConfigPutRequestEndpointEndpointName) GetLocked() bool {
	if o == nil || IsNil(o.Locked) {
		var ret bool
		return ret
	}
	return *o.Locked
}

// GetLockedOk returns a tuple with the Locked field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEndpointEndpointName) GetLockedOk() (*bool, bool) {
	if o == nil || IsNil(o.Locked) {
		return nil, false
	}
	return o.Locked, true
}

// HasLocked returns a boolean if a field has been set.
func (o *ConfigPutRequestEndpointEndpointName) HasLocked() bool {
	if o != nil && !IsNil(o.Locked) {
		return true
	}

	return false
}

// SetLocked gets a reference to the given bool and assigns it to the Locked field.
func (o *ConfigPutRequestEndpointEndpointName) SetLocked(v bool) {
	o.Locked = &v
}

// GetDisabledPorts returns the DisabledPorts field value if set, zero value otherwise.
func (o *ConfigPutRequestEndpointEndpointName) GetDisabledPorts() string {
	if o == nil || IsNil(o.DisabledPorts) {
		var ret string
		return ret
	}
	return *o.DisabledPorts
}

// GetDisabledPortsOk returns a tuple with the DisabledPorts field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEndpointEndpointName) GetDisabledPortsOk() (*string, bool) {
	if o == nil || IsNil(o.DisabledPorts) {
		return nil, false
	}
	return o.DisabledPorts, true
}

// HasDisabledPorts returns a boolean if a field has been set.
func (o *ConfigPutRequestEndpointEndpointName) HasDisabledPorts() bool {
	if o != nil && !IsNil(o.DisabledPorts) {
		return true
	}

	return false
}

// SetDisabledPorts gets a reference to the given string and assigns it to the DisabledPorts field.
func (o *ConfigPutRequestEndpointEndpointName) SetDisabledPorts(v string) {
	o.DisabledPorts = &v
}

// GetIsFabric returns the IsFabric field value if set, zero value otherwise.
func (o *ConfigPutRequestEndpointEndpointName) GetIsFabric() bool {
	if o == nil || IsNil(o.IsFabric) {
		var ret bool
		return ret
	}
	return *o.IsFabric
}

// GetIsFabricOk returns a tuple with the IsFabric field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEndpointEndpointName) GetIsFabricOk() (*bool, bool) {
	if o == nil || IsNil(o.IsFabric) {
		return nil, false
	}
	return o.IsFabric, true
}

// HasIsFabric returns a boolean if a field has been set.
func (o *ConfigPutRequestEndpointEndpointName) HasIsFabric() bool {
	if o != nil && !IsNil(o.IsFabric) {
		return true
	}

	return false
}

// SetIsFabric gets a reference to the given bool and assigns it to the IsFabric field.
func (o *ConfigPutRequestEndpointEndpointName) SetIsFabric(v bool) {
	o.IsFabric = &v
}

// GetOutOfBandManagement returns the OutOfBandManagement field value if set, zero value otherwise.
func (o *ConfigPutRequestEndpointEndpointName) GetOutOfBandManagement() bool {
	if o == nil || IsNil(o.OutOfBandManagement) {
		var ret bool
		return ret
	}
	return *o.OutOfBandManagement
}

// GetOutOfBandManagementOk returns a tuple with the OutOfBandManagement field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEndpointEndpointName) GetOutOfBandManagementOk() (*bool, bool) {
	if o == nil || IsNil(o.OutOfBandManagement) {
		return nil, false
	}
	return o.OutOfBandManagement, true
}

// HasOutOfBandManagement returns a boolean if a field has been set.
func (o *ConfigPutRequestEndpointEndpointName) HasOutOfBandManagement() bool {
	if o != nil && !IsNil(o.OutOfBandManagement) {
		return true
	}

	return false
}

// SetOutOfBandManagement gets a reference to the given bool and assigns it to the OutOfBandManagement field.
func (o *ConfigPutRequestEndpointEndpointName) SetOutOfBandManagement(v bool) {
	o.OutOfBandManagement = &v
}

// GetBadges returns the Badges field value if set, zero value otherwise.
func (o *ConfigPutRequestEndpointEndpointName) GetBadges() []ConfigPutRequestSwitchpointSwitchpointNameBadgesInner {
	if o == nil || IsNil(o.Badges) {
		var ret []ConfigPutRequestSwitchpointSwitchpointNameBadgesInner
		return ret
	}
	return o.Badges
}

// GetBadgesOk returns a tuple with the Badges field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEndpointEndpointName) GetBadgesOk() ([]ConfigPutRequestSwitchpointSwitchpointNameBadgesInner, bool) {
	if o == nil || IsNil(o.Badges) {
		return nil, false
	}
	return o.Badges, true
}

// HasBadges returns a boolean if a field has been set.
func (o *ConfigPutRequestEndpointEndpointName) HasBadges() bool {
	if o != nil && !IsNil(o.Badges) {
		return true
	}

	return false
}

// SetBadges gets a reference to the given []ConfigPutRequestSwitchpointSwitchpointNameBadgesInner and assigns it to the Badges field.
func (o *ConfigPutRequestEndpointEndpointName) SetBadges(v []ConfigPutRequestSwitchpointSwitchpointNameBadgesInner) {
	o.Badges = v
}

// GetChildren returns the Children field value if set, zero value otherwise.
func (o *ConfigPutRequestEndpointEndpointName) GetChildren() []ConfigPutRequestEndpointEndpointNameChildrenInner {
	if o == nil || IsNil(o.Children) {
		var ret []ConfigPutRequestEndpointEndpointNameChildrenInner
		return ret
	}
	return o.Children
}

// GetChildrenOk returns a tuple with the Children field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEndpointEndpointName) GetChildrenOk() ([]ConfigPutRequestEndpointEndpointNameChildrenInner, bool) {
	if o == nil || IsNil(o.Children) {
		return nil, false
	}
	return o.Children, true
}

// HasChildren returns a boolean if a field has been set.
func (o *ConfigPutRequestEndpointEndpointName) HasChildren() bool {
	if o != nil && !IsNil(o.Children) {
		return true
	}

	return false
}

// SetChildren gets a reference to the given []ConfigPutRequestEndpointEndpointNameChildrenInner and assigns it to the Children field.
func (o *ConfigPutRequestEndpointEndpointName) SetChildren(v []ConfigPutRequestEndpointEndpointNameChildrenInner) {
	o.Children = v
}

// GetTrafficMirrors returns the TrafficMirrors field value if set, zero value otherwise.
func (o *ConfigPutRequestEndpointEndpointName) GetTrafficMirrors() []ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner {
	if o == nil || IsNil(o.TrafficMirrors) {
		var ret []ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner
		return ret
	}
	return o.TrafficMirrors
}

// GetTrafficMirrorsOk returns a tuple with the TrafficMirrors field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEndpointEndpointName) GetTrafficMirrorsOk() ([]ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner, bool) {
	if o == nil || IsNil(o.TrafficMirrors) {
		return nil, false
	}
	return o.TrafficMirrors, true
}

// HasTrafficMirrors returns a boolean if a field has been set.
func (o *ConfigPutRequestEndpointEndpointName) HasTrafficMirrors() bool {
	if o != nil && !IsNil(o.TrafficMirrors) {
		return true
	}

	return false
}

// SetTrafficMirrors gets a reference to the given []ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner and assigns it to the TrafficMirrors field.
func (o *ConfigPutRequestEndpointEndpointName) SetTrafficMirrors(v []ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner) {
	o.TrafficMirrors = v
}

// GetEths returns the Eths field value if set, zero value otherwise.
func (o *ConfigPutRequestEndpointEndpointName) GetEths() []ConfigPutRequestEndpointEndpointNameEthsInner {
	if o == nil || IsNil(o.Eths) {
		var ret []ConfigPutRequestEndpointEndpointNameEthsInner
		return ret
	}
	return o.Eths
}

// GetEthsOk returns a tuple with the Eths field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEndpointEndpointName) GetEthsOk() ([]ConfigPutRequestEndpointEndpointNameEthsInner, bool) {
	if o == nil || IsNil(o.Eths) {
		return nil, false
	}
	return o.Eths, true
}

// HasEths returns a boolean if a field has been set.
func (o *ConfigPutRequestEndpointEndpointName) HasEths() bool {
	if o != nil && !IsNil(o.Eths) {
		return true
	}

	return false
}

// SetEths gets a reference to the given []ConfigPutRequestEndpointEndpointNameEthsInner and assigns it to the Eths field.
func (o *ConfigPutRequestEndpointEndpointName) SetEths(v []ConfigPutRequestEndpointEndpointNameEthsInner) {
	o.Eths = v
}

// GetObjectProperties returns the ObjectProperties field value if set, zero value otherwise.
func (o *ConfigPutRequestEndpointEndpointName) GetObjectProperties() ConfigPutRequestEndpointEndpointNameObjectProperties {
	if o == nil || IsNil(o.ObjectProperties) {
		var ret ConfigPutRequestEndpointEndpointNameObjectProperties
		return ret
	}
	return *o.ObjectProperties
}

// GetObjectPropertiesOk returns a tuple with the ObjectProperties field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEndpointEndpointName) GetObjectPropertiesOk() (*ConfigPutRequestEndpointEndpointNameObjectProperties, bool) {
	if o == nil || IsNil(o.ObjectProperties) {
		return nil, false
	}
	return o.ObjectProperties, true
}

// HasObjectProperties returns a boolean if a field has been set.
func (o *ConfigPutRequestEndpointEndpointName) HasObjectProperties() bool {
	if o != nil && !IsNil(o.ObjectProperties) {
		return true
	}

	return false
}

// SetObjectProperties gets a reference to the given ConfigPutRequestEndpointEndpointNameObjectProperties and assigns it to the ObjectProperties field.
func (o *ConfigPutRequestEndpointEndpointName) SetObjectProperties(v ConfigPutRequestEndpointEndpointNameObjectProperties) {
	o.ObjectProperties = &v
}

func (o ConfigPutRequestEndpointEndpointName) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ConfigPutRequestEndpointEndpointName) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.DeviceSerialNumber) {
		toSerialize["device_serial_number"] = o.DeviceSerialNumber
	}
	if !IsNil(o.ConnectedBundle) {
		toSerialize["connected_bundle"] = o.ConnectedBundle
	}
	if !IsNil(o.ConnectedBundleRefType) {
		toSerialize["connected_bundle_ref_type_"] = o.ConnectedBundleRefType
	}
	if !IsNil(o.ReadOnlyMode) {
		toSerialize["read_only_mode"] = o.ReadOnlyMode
	}
	if !IsNil(o.Locked) {
		toSerialize["locked"] = o.Locked
	}
	if !IsNil(o.DisabledPorts) {
		toSerialize["disabled_ports"] = o.DisabledPorts
	}
	if !IsNil(o.IsFabric) {
		toSerialize["is_fabric"] = o.IsFabric
	}
	if !IsNil(o.OutOfBandManagement) {
		toSerialize["out_of_band_management"] = o.OutOfBandManagement
	}
	if !IsNil(o.Badges) {
		toSerialize["badges"] = o.Badges
	}
	if !IsNil(o.Children) {
		toSerialize["children"] = o.Children
	}
	if !IsNil(o.TrafficMirrors) {
		toSerialize["traffic_mirrors"] = o.TrafficMirrors
	}
	if !IsNil(o.Eths) {
		toSerialize["eths"] = o.Eths
	}
	if !IsNil(o.ObjectProperties) {
		toSerialize["object_properties"] = o.ObjectProperties
	}
	return toSerialize, nil
}

type NullableConfigPutRequestEndpointEndpointName struct {
	value *ConfigPutRequestEndpointEndpointName
	isSet bool
}

func (v NullableConfigPutRequestEndpointEndpointName) Get() *ConfigPutRequestEndpointEndpointName {
	return v.value
}

func (v *NullableConfigPutRequestEndpointEndpointName) Set(val *ConfigPutRequestEndpointEndpointName) {
	v.value = val
	v.isSet = true
}

func (v NullableConfigPutRequestEndpointEndpointName) IsSet() bool {
	return v.isSet
}

func (v *NullableConfigPutRequestEndpointEndpointName) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConfigPutRequestEndpointEndpointName(val *ConfigPutRequestEndpointEndpointName) *NullableConfigPutRequestEndpointEndpointName {
	return &NullableConfigPutRequestEndpointEndpointName{value: val, isSet: true}
}

func (v NullableConfigPutRequestEndpointEndpointName) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConfigPutRequestEndpointEndpointName) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


