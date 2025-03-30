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

// checks if the ConfigPutRequestEndpointViewEndpointViewName type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ConfigPutRequestEndpointViewEndpointViewName{}

// ConfigPutRequestEndpointViewEndpointViewName struct for ConfigPutRequestEndpointViewEndpointViewName
type ConfigPutRequestEndpointViewEndpointViewName struct {
	// Object Name. Must be unique.
	Name *string `json:"name,omitempty"`
	// Enable object.
	Enable *bool `json:"enable,omitempty"`
	// Always Endpoints
	Type *string `json:"type,omitempty"`
	// One of Remote|Endpoint|Preprovisioned
	Location *string `json:"location,omitempty"`
	// Show on the summary view
	OnSummary *bool `json:"on_summary,omitempty"`
	// Show Upgrader Pie Chart on Summary
	UpgraderSummary *bool `json:"upgrader_summary,omitempty"`
	// Show Installation Pie Chart on Summary
	InstallationSummary *bool `json:"installation_summary,omitempty"`
	// Show Comm Pie Chart on Summary
	CommSummary *bool `json:"comm_summary,omitempty"`
	// Show Provisioning Pie Chart on Summary
	ProvisioningSummary *bool `json:"provisioning_summary,omitempty"`
	OrRules []ConfigPutRequestEndpointViewEndpointViewNameOrRulesInner `json:"orRules,omitempty"`
	ObjectProperties map[string]interface{} `json:"object_properties,omitempty"`
}

// NewConfigPutRequestEndpointViewEndpointViewName instantiates a new ConfigPutRequestEndpointViewEndpointViewName object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConfigPutRequestEndpointViewEndpointViewName() *ConfigPutRequestEndpointViewEndpointViewName {
	this := ConfigPutRequestEndpointViewEndpointViewName{}
	var name string = ""
	this.Name = &name
	var enable bool = true
	this.Enable = &enable
	var type_ string = "Endpoints"
	this.Type = &type_
	var location string = ""
	this.Location = &location
	var onSummary bool = true
	this.OnSummary = &onSummary
	var upgraderSummary bool = true
	this.UpgraderSummary = &upgraderSummary
	var installationSummary bool = true
	this.InstallationSummary = &installationSummary
	var commSummary bool = true
	this.CommSummary = &commSummary
	var provisioningSummary bool = true
	this.ProvisioningSummary = &provisioningSummary
	return &this
}

// NewConfigPutRequestEndpointViewEndpointViewNameWithDefaults instantiates a new ConfigPutRequestEndpointViewEndpointViewName object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConfigPutRequestEndpointViewEndpointViewNameWithDefaults() *ConfigPutRequestEndpointViewEndpointViewName {
	this := ConfigPutRequestEndpointViewEndpointViewName{}
	var name string = ""
	this.Name = &name
	var enable bool = true
	this.Enable = &enable
	var type_ string = "Endpoints"
	this.Type = &type_
	var location string = ""
	this.Location = &location
	var onSummary bool = true
	this.OnSummary = &onSummary
	var upgraderSummary bool = true
	this.UpgraderSummary = &upgraderSummary
	var installationSummary bool = true
	this.InstallationSummary = &installationSummary
	var commSummary bool = true
	this.CommSummary = &commSummary
	var provisioningSummary bool = true
	this.ProvisioningSummary = &provisioningSummary
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *ConfigPutRequestEndpointViewEndpointViewName) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEndpointViewEndpointViewName) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *ConfigPutRequestEndpointViewEndpointViewName) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *ConfigPutRequestEndpointViewEndpointViewName) SetName(v string) {
	o.Name = &v
}

// GetEnable returns the Enable field value if set, zero value otherwise.
func (o *ConfigPutRequestEndpointViewEndpointViewName) GetEnable() bool {
	if o == nil || IsNil(o.Enable) {
		var ret bool
		return ret
	}
	return *o.Enable
}

// GetEnableOk returns a tuple with the Enable field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEndpointViewEndpointViewName) GetEnableOk() (*bool, bool) {
	if o == nil || IsNil(o.Enable) {
		return nil, false
	}
	return o.Enable, true
}

// HasEnable returns a boolean if a field has been set.
func (o *ConfigPutRequestEndpointViewEndpointViewName) HasEnable() bool {
	if o != nil && !IsNil(o.Enable) {
		return true
	}

	return false
}

// SetEnable gets a reference to the given bool and assigns it to the Enable field.
func (o *ConfigPutRequestEndpointViewEndpointViewName) SetEnable(v bool) {
	o.Enable = &v
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *ConfigPutRequestEndpointViewEndpointViewName) GetType() string {
	if o == nil || IsNil(o.Type) {
		var ret string
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEndpointViewEndpointViewName) GetTypeOk() (*string, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *ConfigPutRequestEndpointViewEndpointViewName) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given string and assigns it to the Type field.
func (o *ConfigPutRequestEndpointViewEndpointViewName) SetType(v string) {
	o.Type = &v
}

// GetLocation returns the Location field value if set, zero value otherwise.
func (o *ConfigPutRequestEndpointViewEndpointViewName) GetLocation() string {
	if o == nil || IsNil(o.Location) {
		var ret string
		return ret
	}
	return *o.Location
}

// GetLocationOk returns a tuple with the Location field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEndpointViewEndpointViewName) GetLocationOk() (*string, bool) {
	if o == nil || IsNil(o.Location) {
		return nil, false
	}
	return o.Location, true
}

// HasLocation returns a boolean if a field has been set.
func (o *ConfigPutRequestEndpointViewEndpointViewName) HasLocation() bool {
	if o != nil && !IsNil(o.Location) {
		return true
	}

	return false
}

// SetLocation gets a reference to the given string and assigns it to the Location field.
func (o *ConfigPutRequestEndpointViewEndpointViewName) SetLocation(v string) {
	o.Location = &v
}

// GetOnSummary returns the OnSummary field value if set, zero value otherwise.
func (o *ConfigPutRequestEndpointViewEndpointViewName) GetOnSummary() bool {
	if o == nil || IsNil(o.OnSummary) {
		var ret bool
		return ret
	}
	return *o.OnSummary
}

// GetOnSummaryOk returns a tuple with the OnSummary field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEndpointViewEndpointViewName) GetOnSummaryOk() (*bool, bool) {
	if o == nil || IsNil(o.OnSummary) {
		return nil, false
	}
	return o.OnSummary, true
}

// HasOnSummary returns a boolean if a field has been set.
func (o *ConfigPutRequestEndpointViewEndpointViewName) HasOnSummary() bool {
	if o != nil && !IsNil(o.OnSummary) {
		return true
	}

	return false
}

// SetOnSummary gets a reference to the given bool and assigns it to the OnSummary field.
func (o *ConfigPutRequestEndpointViewEndpointViewName) SetOnSummary(v bool) {
	o.OnSummary = &v
}

// GetUpgraderSummary returns the UpgraderSummary field value if set, zero value otherwise.
func (o *ConfigPutRequestEndpointViewEndpointViewName) GetUpgraderSummary() bool {
	if o == nil || IsNil(o.UpgraderSummary) {
		var ret bool
		return ret
	}
	return *o.UpgraderSummary
}

// GetUpgraderSummaryOk returns a tuple with the UpgraderSummary field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEndpointViewEndpointViewName) GetUpgraderSummaryOk() (*bool, bool) {
	if o == nil || IsNil(o.UpgraderSummary) {
		return nil, false
	}
	return o.UpgraderSummary, true
}

// HasUpgraderSummary returns a boolean if a field has been set.
func (o *ConfigPutRequestEndpointViewEndpointViewName) HasUpgraderSummary() bool {
	if o != nil && !IsNil(o.UpgraderSummary) {
		return true
	}

	return false
}

// SetUpgraderSummary gets a reference to the given bool and assigns it to the UpgraderSummary field.
func (o *ConfigPutRequestEndpointViewEndpointViewName) SetUpgraderSummary(v bool) {
	o.UpgraderSummary = &v
}

// GetInstallationSummary returns the InstallationSummary field value if set, zero value otherwise.
func (o *ConfigPutRequestEndpointViewEndpointViewName) GetInstallationSummary() bool {
	if o == nil || IsNil(o.InstallationSummary) {
		var ret bool
		return ret
	}
	return *o.InstallationSummary
}

// GetInstallationSummaryOk returns a tuple with the InstallationSummary field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEndpointViewEndpointViewName) GetInstallationSummaryOk() (*bool, bool) {
	if o == nil || IsNil(o.InstallationSummary) {
		return nil, false
	}
	return o.InstallationSummary, true
}

// HasInstallationSummary returns a boolean if a field has been set.
func (o *ConfigPutRequestEndpointViewEndpointViewName) HasInstallationSummary() bool {
	if o != nil && !IsNil(o.InstallationSummary) {
		return true
	}

	return false
}

// SetInstallationSummary gets a reference to the given bool and assigns it to the InstallationSummary field.
func (o *ConfigPutRequestEndpointViewEndpointViewName) SetInstallationSummary(v bool) {
	o.InstallationSummary = &v
}

// GetCommSummary returns the CommSummary field value if set, zero value otherwise.
func (o *ConfigPutRequestEndpointViewEndpointViewName) GetCommSummary() bool {
	if o == nil || IsNil(o.CommSummary) {
		var ret bool
		return ret
	}
	return *o.CommSummary
}

// GetCommSummaryOk returns a tuple with the CommSummary field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEndpointViewEndpointViewName) GetCommSummaryOk() (*bool, bool) {
	if o == nil || IsNil(o.CommSummary) {
		return nil, false
	}
	return o.CommSummary, true
}

// HasCommSummary returns a boolean if a field has been set.
func (o *ConfigPutRequestEndpointViewEndpointViewName) HasCommSummary() bool {
	if o != nil && !IsNil(o.CommSummary) {
		return true
	}

	return false
}

// SetCommSummary gets a reference to the given bool and assigns it to the CommSummary field.
func (o *ConfigPutRequestEndpointViewEndpointViewName) SetCommSummary(v bool) {
	o.CommSummary = &v
}

// GetProvisioningSummary returns the ProvisioningSummary field value if set, zero value otherwise.
func (o *ConfigPutRequestEndpointViewEndpointViewName) GetProvisioningSummary() bool {
	if o == nil || IsNil(o.ProvisioningSummary) {
		var ret bool
		return ret
	}
	return *o.ProvisioningSummary
}

// GetProvisioningSummaryOk returns a tuple with the ProvisioningSummary field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEndpointViewEndpointViewName) GetProvisioningSummaryOk() (*bool, bool) {
	if o == nil || IsNil(o.ProvisioningSummary) {
		return nil, false
	}
	return o.ProvisioningSummary, true
}

// HasProvisioningSummary returns a boolean if a field has been set.
func (o *ConfigPutRequestEndpointViewEndpointViewName) HasProvisioningSummary() bool {
	if o != nil && !IsNil(o.ProvisioningSummary) {
		return true
	}

	return false
}

// SetProvisioningSummary gets a reference to the given bool and assigns it to the ProvisioningSummary field.
func (o *ConfigPutRequestEndpointViewEndpointViewName) SetProvisioningSummary(v bool) {
	o.ProvisioningSummary = &v
}

// GetOrRules returns the OrRules field value if set, zero value otherwise.
func (o *ConfigPutRequestEndpointViewEndpointViewName) GetOrRules() []ConfigPutRequestEndpointViewEndpointViewNameOrRulesInner {
	if o == nil || IsNil(o.OrRules) {
		var ret []ConfigPutRequestEndpointViewEndpointViewNameOrRulesInner
		return ret
	}
	return o.OrRules
}

// GetOrRulesOk returns a tuple with the OrRules field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEndpointViewEndpointViewName) GetOrRulesOk() ([]ConfigPutRequestEndpointViewEndpointViewNameOrRulesInner, bool) {
	if o == nil || IsNil(o.OrRules) {
		return nil, false
	}
	return o.OrRules, true
}

// HasOrRules returns a boolean if a field has been set.
func (o *ConfigPutRequestEndpointViewEndpointViewName) HasOrRules() bool {
	if o != nil && !IsNil(o.OrRules) {
		return true
	}

	return false
}

// SetOrRules gets a reference to the given []ConfigPutRequestEndpointViewEndpointViewNameOrRulesInner and assigns it to the OrRules field.
func (o *ConfigPutRequestEndpointViewEndpointViewName) SetOrRules(v []ConfigPutRequestEndpointViewEndpointViewNameOrRulesInner) {
	o.OrRules = v
}

// GetObjectProperties returns the ObjectProperties field value if set, zero value otherwise.
func (o *ConfigPutRequestEndpointViewEndpointViewName) GetObjectProperties() map[string]interface{} {
	if o == nil || IsNil(o.ObjectProperties) {
		var ret map[string]interface{}
		return ret
	}
	return o.ObjectProperties
}

// GetObjectPropertiesOk returns a tuple with the ObjectProperties field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEndpointViewEndpointViewName) GetObjectPropertiesOk() (map[string]interface{}, bool) {
	if o == nil || IsNil(o.ObjectProperties) {
		return map[string]interface{}{}, false
	}
	return o.ObjectProperties, true
}

// HasObjectProperties returns a boolean if a field has been set.
func (o *ConfigPutRequestEndpointViewEndpointViewName) HasObjectProperties() bool {
	if o != nil && !IsNil(o.ObjectProperties) {
		return true
	}

	return false
}

// SetObjectProperties gets a reference to the given map[string]interface{} and assigns it to the ObjectProperties field.
func (o *ConfigPutRequestEndpointViewEndpointViewName) SetObjectProperties(v map[string]interface{}) {
	o.ObjectProperties = v
}

func (o ConfigPutRequestEndpointViewEndpointViewName) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ConfigPutRequestEndpointViewEndpointViewName) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.Enable) {
		toSerialize["enable"] = o.Enable
	}
	if !IsNil(o.Type) {
		toSerialize["type"] = o.Type
	}
	if !IsNil(o.Location) {
		toSerialize["location"] = o.Location
	}
	if !IsNil(o.OnSummary) {
		toSerialize["on_summary"] = o.OnSummary
	}
	if !IsNil(o.UpgraderSummary) {
		toSerialize["upgrader_summary"] = o.UpgraderSummary
	}
	if !IsNil(o.InstallationSummary) {
		toSerialize["installation_summary"] = o.InstallationSummary
	}
	if !IsNil(o.CommSummary) {
		toSerialize["comm_summary"] = o.CommSummary
	}
	if !IsNil(o.ProvisioningSummary) {
		toSerialize["provisioning_summary"] = o.ProvisioningSummary
	}
	if !IsNil(o.OrRules) {
		toSerialize["orRules"] = o.OrRules
	}
	if !IsNil(o.ObjectProperties) {
		toSerialize["object_properties"] = o.ObjectProperties
	}
	return toSerialize, nil
}

type NullableConfigPutRequestEndpointViewEndpointViewName struct {
	value *ConfigPutRequestEndpointViewEndpointViewName
	isSet bool
}

func (v NullableConfigPutRequestEndpointViewEndpointViewName) Get() *ConfigPutRequestEndpointViewEndpointViewName {
	return v.value
}

func (v *NullableConfigPutRequestEndpointViewEndpointViewName) Set(val *ConfigPutRequestEndpointViewEndpointViewName) {
	v.value = val
	v.isSet = true
}

func (v NullableConfigPutRequestEndpointViewEndpointViewName) IsSet() bool {
	return v.isSet
}

func (v *NullableConfigPutRequestEndpointViewEndpointViewName) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConfigPutRequestEndpointViewEndpointViewName(val *ConfigPutRequestEndpointViewEndpointViewName) *NullableConfigPutRequestEndpointViewEndpointViewName {
	return &NullableConfigPutRequestEndpointViewEndpointViewName{value: val, isSet: true}
}

func (v NullableConfigPutRequestEndpointViewEndpointViewName) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConfigPutRequestEndpointViewEndpointViewName) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


