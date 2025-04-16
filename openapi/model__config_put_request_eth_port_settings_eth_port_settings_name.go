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

// checks if the ConfigPutRequestEthPortSettingsEthPortSettingsName type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ConfigPutRequestEthPortSettingsEthPortSettingsName{}

// ConfigPutRequestEthPortSettingsEthPortSettingsName struct for ConfigPutRequestEthPortSettingsEthPortSettingsName
type ConfigPutRequestEthPortSettingsEthPortSettingsName struct {
	// Object Name. Must be unique.
	Name *string `json:"name,omitempty"`
	// Enable object. It's highly recommended to set this value to true so that validation on the object will be ran.
	Enable *bool `json:"enable,omitempty"`
	// Indicates if port speed and duplex mode should be auto negotiated
	AutoNegotiation *bool `json:"auto_negotiation,omitempty"`
	// Maximum Bit Rate allowed
	MaxBitRate *string `json:"max_bit_rate,omitempty"`
	// Duplex Mode
	DuplexMode *string `json:"duplex_mode,omitempty"`
	// Enable Spanning Tree on the port.  Note: the Spanning Tree Type (VLAN, Port, MST) is controlled in the Site Settings
	StpEnable *bool `json:"stp_enable,omitempty"`
	// Enable Immediate Transition to Forwarding
	FastLearningMode *bool `json:"fast_learning_mode,omitempty"`
	// Block port on BPDU Receive
	BpduGuard *bool `json:"bpdu_guard,omitempty"`
	// Drop all Rx and Tx BPDUs
	BpduFilter *bool `json:"bpdu_filter,omitempty"`
	// Enable Cisco Guard Loop
	GuardLoop *bool `json:"guard_loop,omitempty"`
	// PoE Enable
	PoeEnable *bool `json:"poe_enable,omitempty"`
	// Priority given when assigning power in a limited power situation
	Priority *string `json:"priority,omitempty"`
	// Power the PoE system will attempt to allocate on this port
	AllocatedPower *string `json:"allocated_power,omitempty"`
	// Enable Traffic Storm Protection which prevents excessive broadcast/multicast/unknown-unicast traffic from overwhelming the Switch CPU
	BspEnable *bool `json:"bsp_enable,omitempty"`
	// Broadcast
	Broadcast *bool `json:"broadcast,omitempty"`
	// Multicast
	Multicast *bool `json:"multicast,omitempty"`
	// Max Percentage of the ports bandwidth allowed for broadcast/multicast/unknown-unicast traffic before invoking the protective action
	MaxAllowedValue *int32 `json:"max_allowed_value,omitempty"`
	// Max Percentage of the ports bandwidth allowed for broadcast/multicast/unknown-unicast traffic before invoking the protective action <br>                                                 <div class=\"tab\">                                                     %: Percentage.<br>                                                                                                                                                     kbps: kilobits per second <br>                                                     mbps: megabits per second <br>                                                     gbps: gigabits per second <br>                                                     pps: packet per second <br>                                                     kpps: kilopacket per second <br>                                                 </div>                                                 
	MaxAllowedUnit *string `json:"max_allowed_unit,omitempty"`
	// Action taken if broadcast/multicast/unknown-unicast traffic excedes the Max. One of: <br>                                                 <div class=\"tab\">                                                     Protect: Broadcast/Multicast packets beyond the percent rate are silently dropped. QOS drop counters should indicate the drops.<br><br>                                                     Restrict: Broadcast/Multicast packets beyond the percent rate are dropped. QOS drop counters should indicate the drops.                                                      Alarm is raised . Alarm automatically clears when rate is below configured threshold. <br><br>                                                     Shutdown: Alarm is raised and port is taken out of service. User must administratively Disable and Enable the port to restore service. <br>                                                 </div>                                             
	Action *string `json:"action,omitempty"`
	// FEC is Forward Error Correction which is error correction on the fiber link.                                                 <div class=\"tab\">                                                     Any: Allows switch Negotiation between FC and RS <br>                                                     None: Disables FEC on an interface.<br>                                                     FC: Enables FEC on supported interfaces. FC stands for fire code.<br>                                                     RS: Enables FEC on supported interfaces. RS stands for Reed-Solomon code. <br>                                                     None: VnetC doesn't alter the Switch Value.<br>                                                 </div>                                             
	Fec *string `json:"fec,omitempty"`
	// Ports with this setting will be disabled when link state tracking takes effect
	SingleLink *bool `json:"single_link,omitempty"`
	ObjectProperties *ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties `json:"object_properties,omitempty"`
}

// NewConfigPutRequestEthPortSettingsEthPortSettingsName instantiates a new ConfigPutRequestEthPortSettingsEthPortSettingsName object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewConfigPutRequestEthPortSettingsEthPortSettingsName() *ConfigPutRequestEthPortSettingsEthPortSettingsName {
	this := ConfigPutRequestEthPortSettingsEthPortSettingsName{}
	var name string = ""
	this.Name = &name
	var enable bool = false
	this.Enable = &enable
	var autoNegotiation bool = true
	this.AutoNegotiation = &autoNegotiation
	var maxBitRate string = "-1"
	this.MaxBitRate = &maxBitRate
	var duplexMode string = "Auto"
	this.DuplexMode = &duplexMode
	var stpEnable bool = false
	this.StpEnable = &stpEnable
	var fastLearningMode bool = true
	this.FastLearningMode = &fastLearningMode
	var bpduGuard bool = false
	this.BpduGuard = &bpduGuard
	var bpduFilter bool = false
	this.BpduFilter = &bpduFilter
	var guardLoop bool = false
	this.GuardLoop = &guardLoop
	var poeEnable bool = false
	this.PoeEnable = &poeEnable
	var priority string = "High"
	this.Priority = &priority
	var allocatedPower string = "0.0"
	this.AllocatedPower = &allocatedPower
	var bspEnable bool = false
	this.BspEnable = &bspEnable
	var broadcast bool = true
	this.Broadcast = &broadcast
	var multicast bool = true
	this.Multicast = &multicast
	var maxAllowedValue int32 = 1000
	this.MaxAllowedValue = &maxAllowedValue
	var maxAllowedUnit string = "pps"
	this.MaxAllowedUnit = &maxAllowedUnit
	var action string = "Protect"
	this.Action = &action
	var fec string = "unaltered"
	this.Fec = &fec
	var singleLink bool = false
	this.SingleLink = &singleLink
	return &this
}

// NewConfigPutRequestEthPortSettingsEthPortSettingsNameWithDefaults instantiates a new ConfigPutRequestEthPortSettingsEthPortSettingsName object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewConfigPutRequestEthPortSettingsEthPortSettingsNameWithDefaults() *ConfigPutRequestEthPortSettingsEthPortSettingsName {
	this := ConfigPutRequestEthPortSettingsEthPortSettingsName{}
	var name string = ""
	this.Name = &name
	var enable bool = false
	this.Enable = &enable
	var autoNegotiation bool = true
	this.AutoNegotiation = &autoNegotiation
	var maxBitRate string = "-1"
	this.MaxBitRate = &maxBitRate
	var duplexMode string = "Auto"
	this.DuplexMode = &duplexMode
	var stpEnable bool = false
	this.StpEnable = &stpEnable
	var fastLearningMode bool = true
	this.FastLearningMode = &fastLearningMode
	var bpduGuard bool = false
	this.BpduGuard = &bpduGuard
	var bpduFilter bool = false
	this.BpduFilter = &bpduFilter
	var guardLoop bool = false
	this.GuardLoop = &guardLoop
	var poeEnable bool = false
	this.PoeEnable = &poeEnable
	var priority string = "High"
	this.Priority = &priority
	var allocatedPower string = "0.0"
	this.AllocatedPower = &allocatedPower
	var bspEnable bool = false
	this.BspEnable = &bspEnable
	var broadcast bool = true
	this.Broadcast = &broadcast
	var multicast bool = true
	this.Multicast = &multicast
	var maxAllowedValue int32 = 1000
	this.MaxAllowedValue = &maxAllowedValue
	var maxAllowedUnit string = "pps"
	this.MaxAllowedUnit = &maxAllowedUnit
	var action string = "Protect"
	this.Action = &action
	var fec string = "unaltered"
	this.Fec = &fec
	var singleLink bool = false
	this.SingleLink = &singleLink
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetName(v string) {
	o.Name = &v
}

// GetEnable returns the Enable field value if set, zero value otherwise.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetEnable() bool {
	if o == nil || IsNil(o.Enable) {
		var ret bool
		return ret
	}
	return *o.Enable
}

// GetEnableOk returns a tuple with the Enable field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetEnableOk() (*bool, bool) {
	if o == nil || IsNil(o.Enable) {
		return nil, false
	}
	return o.Enable, true
}

// HasEnable returns a boolean if a field has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasEnable() bool {
	if o != nil && !IsNil(o.Enable) {
		return true
	}

	return false
}

// SetEnable gets a reference to the given bool and assigns it to the Enable field.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetEnable(v bool) {
	o.Enable = &v
}

// GetAutoNegotiation returns the AutoNegotiation field value if set, zero value otherwise.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetAutoNegotiation() bool {
	if o == nil || IsNil(o.AutoNegotiation) {
		var ret bool
		return ret
	}
	return *o.AutoNegotiation
}

// GetAutoNegotiationOk returns a tuple with the AutoNegotiation field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetAutoNegotiationOk() (*bool, bool) {
	if o == nil || IsNil(o.AutoNegotiation) {
		return nil, false
	}
	return o.AutoNegotiation, true
}

// HasAutoNegotiation returns a boolean if a field has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasAutoNegotiation() bool {
	if o != nil && !IsNil(o.AutoNegotiation) {
		return true
	}

	return false
}

// SetAutoNegotiation gets a reference to the given bool and assigns it to the AutoNegotiation field.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetAutoNegotiation(v bool) {
	o.AutoNegotiation = &v
}

// GetMaxBitRate returns the MaxBitRate field value if set, zero value otherwise.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetMaxBitRate() string {
	if o == nil || IsNil(o.MaxBitRate) {
		var ret string
		return ret
	}
	return *o.MaxBitRate
}

// GetMaxBitRateOk returns a tuple with the MaxBitRate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetMaxBitRateOk() (*string, bool) {
	if o == nil || IsNil(o.MaxBitRate) {
		return nil, false
	}
	return o.MaxBitRate, true
}

// HasMaxBitRate returns a boolean if a field has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasMaxBitRate() bool {
	if o != nil && !IsNil(o.MaxBitRate) {
		return true
	}

	return false
}

// SetMaxBitRate gets a reference to the given string and assigns it to the MaxBitRate field.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetMaxBitRate(v string) {
	o.MaxBitRate = &v
}

// GetDuplexMode returns the DuplexMode field value if set, zero value otherwise.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetDuplexMode() string {
	if o == nil || IsNil(o.DuplexMode) {
		var ret string
		return ret
	}
	return *o.DuplexMode
}

// GetDuplexModeOk returns a tuple with the DuplexMode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetDuplexModeOk() (*string, bool) {
	if o == nil || IsNil(o.DuplexMode) {
		return nil, false
	}
	return o.DuplexMode, true
}

// HasDuplexMode returns a boolean if a field has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasDuplexMode() bool {
	if o != nil && !IsNil(o.DuplexMode) {
		return true
	}

	return false
}

// SetDuplexMode gets a reference to the given string and assigns it to the DuplexMode field.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetDuplexMode(v string) {
	o.DuplexMode = &v
}

// GetStpEnable returns the StpEnable field value if set, zero value otherwise.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetStpEnable() bool {
	if o == nil || IsNil(o.StpEnable) {
		var ret bool
		return ret
	}
	return *o.StpEnable
}

// GetStpEnableOk returns a tuple with the StpEnable field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetStpEnableOk() (*bool, bool) {
	if o == nil || IsNil(o.StpEnable) {
		return nil, false
	}
	return o.StpEnable, true
}

// HasStpEnable returns a boolean if a field has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasStpEnable() bool {
	if o != nil && !IsNil(o.StpEnable) {
		return true
	}

	return false
}

// SetStpEnable gets a reference to the given bool and assigns it to the StpEnable field.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetStpEnable(v bool) {
	o.StpEnable = &v
}

// GetFastLearningMode returns the FastLearningMode field value if set, zero value otherwise.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetFastLearningMode() bool {
	if o == nil || IsNil(o.FastLearningMode) {
		var ret bool
		return ret
	}
	return *o.FastLearningMode
}

// GetFastLearningModeOk returns a tuple with the FastLearningMode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetFastLearningModeOk() (*bool, bool) {
	if o == nil || IsNil(o.FastLearningMode) {
		return nil, false
	}
	return o.FastLearningMode, true
}

// HasFastLearningMode returns a boolean if a field has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasFastLearningMode() bool {
	if o != nil && !IsNil(o.FastLearningMode) {
		return true
	}

	return false
}

// SetFastLearningMode gets a reference to the given bool and assigns it to the FastLearningMode field.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetFastLearningMode(v bool) {
	o.FastLearningMode = &v
}

// GetBpduGuard returns the BpduGuard field value if set, zero value otherwise.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetBpduGuard() bool {
	if o == nil || IsNil(o.BpduGuard) {
		var ret bool
		return ret
	}
	return *o.BpduGuard
}

// GetBpduGuardOk returns a tuple with the BpduGuard field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetBpduGuardOk() (*bool, bool) {
	if o == nil || IsNil(o.BpduGuard) {
		return nil, false
	}
	return o.BpduGuard, true
}

// HasBpduGuard returns a boolean if a field has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasBpduGuard() bool {
	if o != nil && !IsNil(o.BpduGuard) {
		return true
	}

	return false
}

// SetBpduGuard gets a reference to the given bool and assigns it to the BpduGuard field.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetBpduGuard(v bool) {
	o.BpduGuard = &v
}

// GetBpduFilter returns the BpduFilter field value if set, zero value otherwise.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetBpduFilter() bool {
	if o == nil || IsNil(o.BpduFilter) {
		var ret bool
		return ret
	}
	return *o.BpduFilter
}

// GetBpduFilterOk returns a tuple with the BpduFilter field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetBpduFilterOk() (*bool, bool) {
	if o == nil || IsNil(o.BpduFilter) {
		return nil, false
	}
	return o.BpduFilter, true
}

// HasBpduFilter returns a boolean if a field has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasBpduFilter() bool {
	if o != nil && !IsNil(o.BpduFilter) {
		return true
	}

	return false
}

// SetBpduFilter gets a reference to the given bool and assigns it to the BpduFilter field.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetBpduFilter(v bool) {
	o.BpduFilter = &v
}

// GetGuardLoop returns the GuardLoop field value if set, zero value otherwise.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetGuardLoop() bool {
	if o == nil || IsNil(o.GuardLoop) {
		var ret bool
		return ret
	}
	return *o.GuardLoop
}

// GetGuardLoopOk returns a tuple with the GuardLoop field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetGuardLoopOk() (*bool, bool) {
	if o == nil || IsNil(o.GuardLoop) {
		return nil, false
	}
	return o.GuardLoop, true
}

// HasGuardLoop returns a boolean if a field has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasGuardLoop() bool {
	if o != nil && !IsNil(o.GuardLoop) {
		return true
	}

	return false
}

// SetGuardLoop gets a reference to the given bool and assigns it to the GuardLoop field.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetGuardLoop(v bool) {
	o.GuardLoop = &v
}

// GetPoeEnable returns the PoeEnable field value if set, zero value otherwise.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetPoeEnable() bool {
	if o == nil || IsNil(o.PoeEnable) {
		var ret bool
		return ret
	}
	return *o.PoeEnable
}

// GetPoeEnableOk returns a tuple with the PoeEnable field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetPoeEnableOk() (*bool, bool) {
	if o == nil || IsNil(o.PoeEnable) {
		return nil, false
	}
	return o.PoeEnable, true
}

// HasPoeEnable returns a boolean if a field has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasPoeEnable() bool {
	if o != nil && !IsNil(o.PoeEnable) {
		return true
	}

	return false
}

// SetPoeEnable gets a reference to the given bool and assigns it to the PoeEnable field.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetPoeEnable(v bool) {
	o.PoeEnable = &v
}

// GetPriority returns the Priority field value if set, zero value otherwise.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetPriority() string {
	if o == nil || IsNil(o.Priority) {
		var ret string
		return ret
	}
	return *o.Priority
}

// GetPriorityOk returns a tuple with the Priority field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetPriorityOk() (*string, bool) {
	if o == nil || IsNil(o.Priority) {
		return nil, false
	}
	return o.Priority, true
}

// HasPriority returns a boolean if a field has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasPriority() bool {
	if o != nil && !IsNil(o.Priority) {
		return true
	}

	return false
}

// SetPriority gets a reference to the given string and assigns it to the Priority field.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetPriority(v string) {
	o.Priority = &v
}

// GetAllocatedPower returns the AllocatedPower field value if set, zero value otherwise.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetAllocatedPower() string {
	if o == nil || IsNil(o.AllocatedPower) {
		var ret string
		return ret
	}
	return *o.AllocatedPower
}

// GetAllocatedPowerOk returns a tuple with the AllocatedPower field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetAllocatedPowerOk() (*string, bool) {
	if o == nil || IsNil(o.AllocatedPower) {
		return nil, false
	}
	return o.AllocatedPower, true
}

// HasAllocatedPower returns a boolean if a field has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasAllocatedPower() bool {
	if o != nil && !IsNil(o.AllocatedPower) {
		return true
	}

	return false
}

// SetAllocatedPower gets a reference to the given string and assigns it to the AllocatedPower field.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetAllocatedPower(v string) {
	o.AllocatedPower = &v
}

// GetBspEnable returns the BspEnable field value if set, zero value otherwise.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetBspEnable() bool {
	if o == nil || IsNil(o.BspEnable) {
		var ret bool
		return ret
	}
	return *o.BspEnable
}

// GetBspEnableOk returns a tuple with the BspEnable field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetBspEnableOk() (*bool, bool) {
	if o == nil || IsNil(o.BspEnable) {
		return nil, false
	}
	return o.BspEnable, true
}

// HasBspEnable returns a boolean if a field has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasBspEnable() bool {
	if o != nil && !IsNil(o.BspEnable) {
		return true
	}

	return false
}

// SetBspEnable gets a reference to the given bool and assigns it to the BspEnable field.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetBspEnable(v bool) {
	o.BspEnable = &v
}

// GetBroadcast returns the Broadcast field value if set, zero value otherwise.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetBroadcast() bool {
	if o == nil || IsNil(o.Broadcast) {
		var ret bool
		return ret
	}
	return *o.Broadcast
}

// GetBroadcastOk returns a tuple with the Broadcast field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetBroadcastOk() (*bool, bool) {
	if o == nil || IsNil(o.Broadcast) {
		return nil, false
	}
	return o.Broadcast, true
}

// HasBroadcast returns a boolean if a field has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasBroadcast() bool {
	if o != nil && !IsNil(o.Broadcast) {
		return true
	}

	return false
}

// SetBroadcast gets a reference to the given bool and assigns it to the Broadcast field.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetBroadcast(v bool) {
	o.Broadcast = &v
}

// GetMulticast returns the Multicast field value if set, zero value otherwise.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetMulticast() bool {
	if o == nil || IsNil(o.Multicast) {
		var ret bool
		return ret
	}
	return *o.Multicast
}

// GetMulticastOk returns a tuple with the Multicast field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetMulticastOk() (*bool, bool) {
	if o == nil || IsNil(o.Multicast) {
		return nil, false
	}
	return o.Multicast, true
}

// HasMulticast returns a boolean if a field has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasMulticast() bool {
	if o != nil && !IsNil(o.Multicast) {
		return true
	}

	return false
}

// SetMulticast gets a reference to the given bool and assigns it to the Multicast field.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetMulticast(v bool) {
	o.Multicast = &v
}

// GetMaxAllowedValue returns the MaxAllowedValue field value if set, zero value otherwise.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetMaxAllowedValue() int32 {
	if o == nil || IsNil(o.MaxAllowedValue) {
		var ret int32
		return ret
	}
	return *o.MaxAllowedValue
}

// GetMaxAllowedValueOk returns a tuple with the MaxAllowedValue field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetMaxAllowedValueOk() (*int32, bool) {
	if o == nil || IsNil(o.MaxAllowedValue) {
		return nil, false
	}
	return o.MaxAllowedValue, true
}

// HasMaxAllowedValue returns a boolean if a field has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasMaxAllowedValue() bool {
	if o != nil && !IsNil(o.MaxAllowedValue) {
		return true
	}

	return false
}

// SetMaxAllowedValue gets a reference to the given int32 and assigns it to the MaxAllowedValue field.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetMaxAllowedValue(v int32) {
	o.MaxAllowedValue = &v
}

// GetMaxAllowedUnit returns the MaxAllowedUnit field value if set, zero value otherwise.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetMaxAllowedUnit() string {
	if o == nil || IsNil(o.MaxAllowedUnit) {
		var ret string
		return ret
	}
	return *o.MaxAllowedUnit
}

// GetMaxAllowedUnitOk returns a tuple with the MaxAllowedUnit field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetMaxAllowedUnitOk() (*string, bool) {
	if o == nil || IsNil(o.MaxAllowedUnit) {
		return nil, false
	}
	return o.MaxAllowedUnit, true
}

// HasMaxAllowedUnit returns a boolean if a field has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasMaxAllowedUnit() bool {
	if o != nil && !IsNil(o.MaxAllowedUnit) {
		return true
	}

	return false
}

// SetMaxAllowedUnit gets a reference to the given string and assigns it to the MaxAllowedUnit field.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetMaxAllowedUnit(v string) {
	o.MaxAllowedUnit = &v
}

// GetAction returns the Action field value if set, zero value otherwise.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetAction() string {
	if o == nil || IsNil(o.Action) {
		var ret string
		return ret
	}
	return *o.Action
}

// GetActionOk returns a tuple with the Action field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetActionOk() (*string, bool) {
	if o == nil || IsNil(o.Action) {
		return nil, false
	}
	return o.Action, true
}

// HasAction returns a boolean if a field has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasAction() bool {
	if o != nil && !IsNil(o.Action) {
		return true
	}

	return false
}

// SetAction gets a reference to the given string and assigns it to the Action field.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetAction(v string) {
	o.Action = &v
}

// GetFec returns the Fec field value if set, zero value otherwise.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetFec() string {
	if o == nil || IsNil(o.Fec) {
		var ret string
		return ret
	}
	return *o.Fec
}

// GetFecOk returns a tuple with the Fec field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetFecOk() (*string, bool) {
	if o == nil || IsNil(o.Fec) {
		return nil, false
	}
	return o.Fec, true
}

// HasFec returns a boolean if a field has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasFec() bool {
	if o != nil && !IsNil(o.Fec) {
		return true
	}

	return false
}

// SetFec gets a reference to the given string and assigns it to the Fec field.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetFec(v string) {
	o.Fec = &v
}

// GetSingleLink returns the SingleLink field value if set, zero value otherwise.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetSingleLink() bool {
	if o == nil || IsNil(o.SingleLink) {
		var ret bool
		return ret
	}
	return *o.SingleLink
}

// GetSingleLinkOk returns a tuple with the SingleLink field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetSingleLinkOk() (*bool, bool) {
	if o == nil || IsNil(o.SingleLink) {
		return nil, false
	}
	return o.SingleLink, true
}

// HasSingleLink returns a boolean if a field has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasSingleLink() bool {
	if o != nil && !IsNil(o.SingleLink) {
		return true
	}

	return false
}

// SetSingleLink gets a reference to the given bool and assigns it to the SingleLink field.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetSingleLink(v bool) {
	o.SingleLink = &v
}

// GetObjectProperties returns the ObjectProperties field value if set, zero value otherwise.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetObjectProperties() ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties {
	if o == nil || IsNil(o.ObjectProperties) {
		var ret ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties
		return ret
	}
	return *o.ObjectProperties
}

// GetObjectPropertiesOk returns a tuple with the ObjectProperties field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) GetObjectPropertiesOk() (*ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties, bool) {
	if o == nil || IsNil(o.ObjectProperties) {
		return nil, false
	}
	return o.ObjectProperties, true
}

// HasObjectProperties returns a boolean if a field has been set.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) HasObjectProperties() bool {
	if o != nil && !IsNil(o.ObjectProperties) {
		return true
	}

	return false
}

// SetObjectProperties gets a reference to the given ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties and assigns it to the ObjectProperties field.
func (o *ConfigPutRequestEthPortSettingsEthPortSettingsName) SetObjectProperties(v ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties) {
	o.ObjectProperties = &v
}

func (o ConfigPutRequestEthPortSettingsEthPortSettingsName) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ConfigPutRequestEthPortSettingsEthPortSettingsName) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.Enable) {
		toSerialize["enable"] = o.Enable
	}
	if !IsNil(o.AutoNegotiation) {
		toSerialize["auto_negotiation"] = o.AutoNegotiation
	}
	if !IsNil(o.MaxBitRate) {
		toSerialize["max_bit_rate"] = o.MaxBitRate
	}
	if !IsNil(o.DuplexMode) {
		toSerialize["duplex_mode"] = o.DuplexMode
	}
	if !IsNil(o.StpEnable) {
		toSerialize["stp_enable"] = o.StpEnable
	}
	if !IsNil(o.FastLearningMode) {
		toSerialize["fast_learning_mode"] = o.FastLearningMode
	}
	if !IsNil(o.BpduGuard) {
		toSerialize["bpdu_guard"] = o.BpduGuard
	}
	if !IsNil(o.BpduFilter) {
		toSerialize["bpdu_filter"] = o.BpduFilter
	}
	if !IsNil(o.GuardLoop) {
		toSerialize["guard_loop"] = o.GuardLoop
	}
	if !IsNil(o.PoeEnable) {
		toSerialize["poe_enable"] = o.PoeEnable
	}
	if !IsNil(o.Priority) {
		toSerialize["priority"] = o.Priority
	}
	if !IsNil(o.AllocatedPower) {
		toSerialize["allocated_power"] = o.AllocatedPower
	}
	if !IsNil(o.BspEnable) {
		toSerialize["bsp_enable"] = o.BspEnable
	}
	if !IsNil(o.Broadcast) {
		toSerialize["broadcast"] = o.Broadcast
	}
	if !IsNil(o.Multicast) {
		toSerialize["multicast"] = o.Multicast
	}
	if !IsNil(o.MaxAllowedValue) {
		toSerialize["max_allowed_value"] = o.MaxAllowedValue
	}
	if !IsNil(o.MaxAllowedUnit) {
		toSerialize["max_allowed_unit"] = o.MaxAllowedUnit
	}
	if !IsNil(o.Action) {
		toSerialize["action"] = o.Action
	}
	if !IsNil(o.Fec) {
		toSerialize["fec"] = o.Fec
	}
	if !IsNil(o.SingleLink) {
		toSerialize["single_link"] = o.SingleLink
	}
	if !IsNil(o.ObjectProperties) {
		toSerialize["object_properties"] = o.ObjectProperties
	}
	return toSerialize, nil
}

type NullableConfigPutRequestEthPortSettingsEthPortSettingsName struct {
	value *ConfigPutRequestEthPortSettingsEthPortSettingsName
	isSet bool
}

func (v NullableConfigPutRequestEthPortSettingsEthPortSettingsName) Get() *ConfigPutRequestEthPortSettingsEthPortSettingsName {
	return v.value
}

func (v *NullableConfigPutRequestEthPortSettingsEthPortSettingsName) Set(val *ConfigPutRequestEthPortSettingsEthPortSettingsName) {
	v.value = val
	v.isSet = true
}

func (v NullableConfigPutRequestEthPortSettingsEthPortSettingsName) IsSet() bool {
	return v.isSet
}

func (v *NullableConfigPutRequestEthPortSettingsEthPortSettingsName) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableConfigPutRequestEthPortSettingsEthPortSettingsName(val *ConfigPutRequestEthPortSettingsEthPortSettingsName) *NullableConfigPutRequestEthPortSettingsEthPortSettingsName {
	return &NullableConfigPutRequestEthPortSettingsEthPortSettingsName{value: val, isSet: true}
}

func (v NullableConfigPutRequestEthPortSettingsEthPortSettingsName) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableConfigPutRequestEthPortSettingsEthPortSettingsName) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


