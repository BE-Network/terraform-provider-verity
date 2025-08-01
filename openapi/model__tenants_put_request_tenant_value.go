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

// checks if the TenantsPutRequestTenantValue type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TenantsPutRequestTenantValue{}

// TenantsPutRequestTenantValue struct for TenantsPutRequestTenantValue
type TenantsPutRequestTenantValue struct {
	// Object Name. Must be unique.
	Name *string `json:"name,omitempty"`
	// Enable object. It's highly recommended to set this value to true so that validation on the object will be ran.
	Enable *bool `json:"enable,omitempty"`
	// VNI value used to transport traffic between services of a Tenant 
	Layer3Vni NullableInt32 `json:"layer_3_vni,omitempty"`
	// Whether or not the value in layer_3_vni field has been automatically assigned or not. Set to false and change layer_3_vni value to edit.
	Layer3VniAutoAssigned *bool `json:"layer_3_vni_auto_assigned_,omitempty"`
	// VLAN value used to transport traffic between services of a Tenant 
	Layer3Vlan NullableInt32 `json:"layer_3_vlan,omitempty"`
	// Whether or not the value in layer_3_vlan field has been automatically assigned or not. Set to false and change layer_3_vlan value to edit.
	Layer3VlanAutoAssigned *bool `json:"layer_3_vlan_auto_assigned_,omitempty"`
	// Range of IPv4 addresses (represented in IPv4 subnet format) used to configure the source IP of each DHCP Relay on each switch that this Tenant is provisioned on.
	DhcpRelaySourceIpv4sSubnet *string `json:"dhcp_relay_source_ipv4s_subnet,omitempty"`
	// Range of IPv6 addresses (represented in IPv6 subnet format) used to configure the source IP of each DHCP Relay on each switch that this Tenant is provisioned on.
	DhcpRelaySourceIpv6sSubnet *string `json:"dhcp_relay_source_ipv6s_subnet,omitempty"`
	// Route Distinguishers are used to maintain uniqueness among identical routes from different routers.  If set, then routes from this Tenant will be identified with this Route Distinguisher (BGP Community).  It should be two numbers separated by a colon.
	RouteDistinguisher *string `json:"route_distinguisher,omitempty"`
	// A route-target (BGP Community) to attach while importing routes into the current tenant. It should be a comma-separated list of BGP Communities: each Community being two numbers separated by a colon.
	RouteTargetImport *string `json:"route_target_import,omitempty"`
	// A route-target (BGP Community) to attach while exporting routes from the current tenant. It should be a comma-separated list of BGP Communities: each Community being two numbers separated by a colon.
	RouteTargetExport *string `json:"route_target_export,omitempty"`
	// A route-map applied to routes imported into the current tenant from other tenants with the purpose of filtering or modifying the routes
	ImportRouteMap *string `json:"import_route_map,omitempty"`
	// Object type for import_route_map field
	ImportRouteMapRefType *string `json:"import_route_map_ref_type_,omitempty"`
	// A route-map applied to routes exported into the current tenant from other tenants with the purpose of filtering or modifying the routes
	ExportRouteMap *string `json:"export_route_map,omitempty"`
	// Object type for export_route_map field
	ExportRouteMapRefType *string `json:"export_route_map_ref_type_,omitempty"`
	// Virtual Routing and Forwarding instance name associated to tenants 
	VrfName *string `json:"vrf_name,omitempty"`
	// Whether or not the value in vrf_name field has been automatically assigned or not. Set to false and change vrf_name value to edit.
	VrfNameAutoAssigned *bool `json:"vrf_name_auto_assigned_,omitempty"`
	RouteTenants []TenantsPutRequestTenantValueRouteTenantsInner `json:"route_tenants,omitempty"`
	ObjectProperties *DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties `json:"object_properties,omitempty"`
	// Enables a leaf switch to originate IPv4 default type-5 EVPN routes across the switching fabric.
	DefaultOriginate *bool `json:"default_originate,omitempty"`
}

// NewTenantsPutRequestTenantValue instantiates a new TenantsPutRequestTenantValue object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTenantsPutRequestTenantValue() *TenantsPutRequestTenantValue {
	this := TenantsPutRequestTenantValue{}
	var name string = ""
	this.Name = &name
	var enable bool = true
	this.Enable = &enable
	var dhcpRelaySourceIpv4sSubnet string = ""
	this.DhcpRelaySourceIpv4sSubnet = &dhcpRelaySourceIpv4sSubnet
	var dhcpRelaySourceIpv6sSubnet string = ""
	this.DhcpRelaySourceIpv6sSubnet = &dhcpRelaySourceIpv6sSubnet
	var routeDistinguisher string = ""
	this.RouteDistinguisher = &routeDistinguisher
	var routeTargetImport string = ""
	this.RouteTargetImport = &routeTargetImport
	var routeTargetExport string = ""
	this.RouteTargetExport = &routeTargetExport
	var importRouteMap string = ""
	this.ImportRouteMap = &importRouteMap
	var exportRouteMap string = ""
	this.ExportRouteMap = &exportRouteMap
	var vrfName string = "(auto)"
	this.VrfName = &vrfName
	var defaultOriginate bool = false
	this.DefaultOriginate = &defaultOriginate
	return &this
}

// NewTenantsPutRequestTenantValueWithDefaults instantiates a new TenantsPutRequestTenantValue object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTenantsPutRequestTenantValueWithDefaults() *TenantsPutRequestTenantValue {
	this := TenantsPutRequestTenantValue{}
	var name string = ""
	this.Name = &name
	var enable bool = true
	this.Enable = &enable
	var dhcpRelaySourceIpv4sSubnet string = ""
	this.DhcpRelaySourceIpv4sSubnet = &dhcpRelaySourceIpv4sSubnet
	var dhcpRelaySourceIpv6sSubnet string = ""
	this.DhcpRelaySourceIpv6sSubnet = &dhcpRelaySourceIpv6sSubnet
	var routeDistinguisher string = ""
	this.RouteDistinguisher = &routeDistinguisher
	var routeTargetImport string = ""
	this.RouteTargetImport = &routeTargetImport
	var routeTargetExport string = ""
	this.RouteTargetExport = &routeTargetExport
	var importRouteMap string = ""
	this.ImportRouteMap = &importRouteMap
	var exportRouteMap string = ""
	this.ExportRouteMap = &exportRouteMap
	var vrfName string = "(auto)"
	this.VrfName = &vrfName
	var defaultOriginate bool = false
	this.DefaultOriginate = &defaultOriginate
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *TenantsPutRequestTenantValue) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TenantsPutRequestTenantValue) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *TenantsPutRequestTenantValue) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *TenantsPutRequestTenantValue) SetName(v string) {
	o.Name = &v
}

// GetEnable returns the Enable field value if set, zero value otherwise.
func (o *TenantsPutRequestTenantValue) GetEnable() bool {
	if o == nil || IsNil(o.Enable) {
		var ret bool
		return ret
	}
	return *o.Enable
}

// GetEnableOk returns a tuple with the Enable field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TenantsPutRequestTenantValue) GetEnableOk() (*bool, bool) {
	if o == nil || IsNil(o.Enable) {
		return nil, false
	}
	return o.Enable, true
}

// HasEnable returns a boolean if a field has been set.
func (o *TenantsPutRequestTenantValue) HasEnable() bool {
	if o != nil && !IsNil(o.Enable) {
		return true
	}

	return false
}

// SetEnable gets a reference to the given bool and assigns it to the Enable field.
func (o *TenantsPutRequestTenantValue) SetEnable(v bool) {
	o.Enable = &v
}

// GetLayer3Vni returns the Layer3Vni field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *TenantsPutRequestTenantValue) GetLayer3Vni() int32 {
	if o == nil || IsNil(o.Layer3Vni.Get()) {
		var ret int32
		return ret
	}
	return *o.Layer3Vni.Get()
}

// GetLayer3VniOk returns a tuple with the Layer3Vni field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *TenantsPutRequestTenantValue) GetLayer3VniOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.Layer3Vni.Get(), o.Layer3Vni.IsSet()
}

// HasLayer3Vni returns a boolean if a field has been set.
func (o *TenantsPutRequestTenantValue) HasLayer3Vni() bool {
	if o != nil && o.Layer3Vni.IsSet() {
		return true
	}

	return false
}

// SetLayer3Vni gets a reference to the given NullableInt32 and assigns it to the Layer3Vni field.
func (o *TenantsPutRequestTenantValue) SetLayer3Vni(v int32) {
	o.Layer3Vni.Set(&v)
}
// SetLayer3VniNil sets the value for Layer3Vni to be an explicit nil
func (o *TenantsPutRequestTenantValue) SetLayer3VniNil() {
	o.Layer3Vni.Set(nil)
}

// UnsetLayer3Vni ensures that no value is present for Layer3Vni, not even an explicit nil
func (o *TenantsPutRequestTenantValue) UnsetLayer3Vni() {
	o.Layer3Vni.Unset()
}

// GetLayer3VniAutoAssigned returns the Layer3VniAutoAssigned field value if set, zero value otherwise.
func (o *TenantsPutRequestTenantValue) GetLayer3VniAutoAssigned() bool {
	if o == nil || IsNil(o.Layer3VniAutoAssigned) {
		var ret bool
		return ret
	}
	return *o.Layer3VniAutoAssigned
}

// GetLayer3VniAutoAssignedOk returns a tuple with the Layer3VniAutoAssigned field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TenantsPutRequestTenantValue) GetLayer3VniAutoAssignedOk() (*bool, bool) {
	if o == nil || IsNil(o.Layer3VniAutoAssigned) {
		return nil, false
	}
	return o.Layer3VniAutoAssigned, true
}

// HasLayer3VniAutoAssigned returns a boolean if a field has been set.
func (o *TenantsPutRequestTenantValue) HasLayer3VniAutoAssigned() bool {
	if o != nil && !IsNil(o.Layer3VniAutoAssigned) {
		return true
	}

	return false
}

// SetLayer3VniAutoAssigned gets a reference to the given bool and assigns it to the Layer3VniAutoAssigned field.
func (o *TenantsPutRequestTenantValue) SetLayer3VniAutoAssigned(v bool) {
	o.Layer3VniAutoAssigned = &v
}

// GetLayer3Vlan returns the Layer3Vlan field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *TenantsPutRequestTenantValue) GetLayer3Vlan() int32 {
	if o == nil || IsNil(o.Layer3Vlan.Get()) {
		var ret int32
		return ret
	}
	return *o.Layer3Vlan.Get()
}

// GetLayer3VlanOk returns a tuple with the Layer3Vlan field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *TenantsPutRequestTenantValue) GetLayer3VlanOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.Layer3Vlan.Get(), o.Layer3Vlan.IsSet()
}

// HasLayer3Vlan returns a boolean if a field has been set.
func (o *TenantsPutRequestTenantValue) HasLayer3Vlan() bool {
	if o != nil && o.Layer3Vlan.IsSet() {
		return true
	}

	return false
}

// SetLayer3Vlan gets a reference to the given NullableInt32 and assigns it to the Layer3Vlan field.
func (o *TenantsPutRequestTenantValue) SetLayer3Vlan(v int32) {
	o.Layer3Vlan.Set(&v)
}
// SetLayer3VlanNil sets the value for Layer3Vlan to be an explicit nil
func (o *TenantsPutRequestTenantValue) SetLayer3VlanNil() {
	o.Layer3Vlan.Set(nil)
}

// UnsetLayer3Vlan ensures that no value is present for Layer3Vlan, not even an explicit nil
func (o *TenantsPutRequestTenantValue) UnsetLayer3Vlan() {
	o.Layer3Vlan.Unset()
}

// GetLayer3VlanAutoAssigned returns the Layer3VlanAutoAssigned field value if set, zero value otherwise.
func (o *TenantsPutRequestTenantValue) GetLayer3VlanAutoAssigned() bool {
	if o == nil || IsNil(o.Layer3VlanAutoAssigned) {
		var ret bool
		return ret
	}
	return *o.Layer3VlanAutoAssigned
}

// GetLayer3VlanAutoAssignedOk returns a tuple with the Layer3VlanAutoAssigned field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TenantsPutRequestTenantValue) GetLayer3VlanAutoAssignedOk() (*bool, bool) {
	if o == nil || IsNil(o.Layer3VlanAutoAssigned) {
		return nil, false
	}
	return o.Layer3VlanAutoAssigned, true
}

// HasLayer3VlanAutoAssigned returns a boolean if a field has been set.
func (o *TenantsPutRequestTenantValue) HasLayer3VlanAutoAssigned() bool {
	if o != nil && !IsNil(o.Layer3VlanAutoAssigned) {
		return true
	}

	return false
}

// SetLayer3VlanAutoAssigned gets a reference to the given bool and assigns it to the Layer3VlanAutoAssigned field.
func (o *TenantsPutRequestTenantValue) SetLayer3VlanAutoAssigned(v bool) {
	o.Layer3VlanAutoAssigned = &v
}

// GetDhcpRelaySourceIpv4sSubnet returns the DhcpRelaySourceIpv4sSubnet field value if set, zero value otherwise.
func (o *TenantsPutRequestTenantValue) GetDhcpRelaySourceIpv4sSubnet() string {
	if o == nil || IsNil(o.DhcpRelaySourceIpv4sSubnet) {
		var ret string
		return ret
	}
	return *o.DhcpRelaySourceIpv4sSubnet
}

// GetDhcpRelaySourceIpv4sSubnetOk returns a tuple with the DhcpRelaySourceIpv4sSubnet field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TenantsPutRequestTenantValue) GetDhcpRelaySourceIpv4sSubnetOk() (*string, bool) {
	if o == nil || IsNil(o.DhcpRelaySourceIpv4sSubnet) {
		return nil, false
	}
	return o.DhcpRelaySourceIpv4sSubnet, true
}

// HasDhcpRelaySourceIpv4sSubnet returns a boolean if a field has been set.
func (o *TenantsPutRequestTenantValue) HasDhcpRelaySourceIpv4sSubnet() bool {
	if o != nil && !IsNil(o.DhcpRelaySourceIpv4sSubnet) {
		return true
	}

	return false
}

// SetDhcpRelaySourceIpv4sSubnet gets a reference to the given string and assigns it to the DhcpRelaySourceIpv4sSubnet field.
func (o *TenantsPutRequestTenantValue) SetDhcpRelaySourceIpv4sSubnet(v string) {
	o.DhcpRelaySourceIpv4sSubnet = &v
}

// GetDhcpRelaySourceIpv6sSubnet returns the DhcpRelaySourceIpv6sSubnet field value if set, zero value otherwise.
func (o *TenantsPutRequestTenantValue) GetDhcpRelaySourceIpv6sSubnet() string {
	if o == nil || IsNil(o.DhcpRelaySourceIpv6sSubnet) {
		var ret string
		return ret
	}
	return *o.DhcpRelaySourceIpv6sSubnet
}

// GetDhcpRelaySourceIpv6sSubnetOk returns a tuple with the DhcpRelaySourceIpv6sSubnet field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TenantsPutRequestTenantValue) GetDhcpRelaySourceIpv6sSubnetOk() (*string, bool) {
	if o == nil || IsNil(o.DhcpRelaySourceIpv6sSubnet) {
		return nil, false
	}
	return o.DhcpRelaySourceIpv6sSubnet, true
}

// HasDhcpRelaySourceIpv6sSubnet returns a boolean if a field has been set.
func (o *TenantsPutRequestTenantValue) HasDhcpRelaySourceIpv6sSubnet() bool {
	if o != nil && !IsNil(o.DhcpRelaySourceIpv6sSubnet) {
		return true
	}

	return false
}

// SetDhcpRelaySourceIpv6sSubnet gets a reference to the given string and assigns it to the DhcpRelaySourceIpv6sSubnet field.
func (o *TenantsPutRequestTenantValue) SetDhcpRelaySourceIpv6sSubnet(v string) {
	o.DhcpRelaySourceIpv6sSubnet = &v
}

// GetRouteDistinguisher returns the RouteDistinguisher field value if set, zero value otherwise.
func (o *TenantsPutRequestTenantValue) GetRouteDistinguisher() string {
	if o == nil || IsNil(o.RouteDistinguisher) {
		var ret string
		return ret
	}
	return *o.RouteDistinguisher
}

// GetRouteDistinguisherOk returns a tuple with the RouteDistinguisher field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TenantsPutRequestTenantValue) GetRouteDistinguisherOk() (*string, bool) {
	if o == nil || IsNil(o.RouteDistinguisher) {
		return nil, false
	}
	return o.RouteDistinguisher, true
}

// HasRouteDistinguisher returns a boolean if a field has been set.
func (o *TenantsPutRequestTenantValue) HasRouteDistinguisher() bool {
	if o != nil && !IsNil(o.RouteDistinguisher) {
		return true
	}

	return false
}

// SetRouteDistinguisher gets a reference to the given string and assigns it to the RouteDistinguisher field.
func (o *TenantsPutRequestTenantValue) SetRouteDistinguisher(v string) {
	o.RouteDistinguisher = &v
}

// GetRouteTargetImport returns the RouteTargetImport field value if set, zero value otherwise.
func (o *TenantsPutRequestTenantValue) GetRouteTargetImport() string {
	if o == nil || IsNil(o.RouteTargetImport) {
		var ret string
		return ret
	}
	return *o.RouteTargetImport
}

// GetRouteTargetImportOk returns a tuple with the RouteTargetImport field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TenantsPutRequestTenantValue) GetRouteTargetImportOk() (*string, bool) {
	if o == nil || IsNil(o.RouteTargetImport) {
		return nil, false
	}
	return o.RouteTargetImport, true
}

// HasRouteTargetImport returns a boolean if a field has been set.
func (o *TenantsPutRequestTenantValue) HasRouteTargetImport() bool {
	if o != nil && !IsNil(o.RouteTargetImport) {
		return true
	}

	return false
}

// SetRouteTargetImport gets a reference to the given string and assigns it to the RouteTargetImport field.
func (o *TenantsPutRequestTenantValue) SetRouteTargetImport(v string) {
	o.RouteTargetImport = &v
}

// GetRouteTargetExport returns the RouteTargetExport field value if set, zero value otherwise.
func (o *TenantsPutRequestTenantValue) GetRouteTargetExport() string {
	if o == nil || IsNil(o.RouteTargetExport) {
		var ret string
		return ret
	}
	return *o.RouteTargetExport
}

// GetRouteTargetExportOk returns a tuple with the RouteTargetExport field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TenantsPutRequestTenantValue) GetRouteTargetExportOk() (*string, bool) {
	if o == nil || IsNil(o.RouteTargetExport) {
		return nil, false
	}
	return o.RouteTargetExport, true
}

// HasRouteTargetExport returns a boolean if a field has been set.
func (o *TenantsPutRequestTenantValue) HasRouteTargetExport() bool {
	if o != nil && !IsNil(o.RouteTargetExport) {
		return true
	}

	return false
}

// SetRouteTargetExport gets a reference to the given string and assigns it to the RouteTargetExport field.
func (o *TenantsPutRequestTenantValue) SetRouteTargetExport(v string) {
	o.RouteTargetExport = &v
}

// GetImportRouteMap returns the ImportRouteMap field value if set, zero value otherwise.
func (o *TenantsPutRequestTenantValue) GetImportRouteMap() string {
	if o == nil || IsNil(o.ImportRouteMap) {
		var ret string
		return ret
	}
	return *o.ImportRouteMap
}

// GetImportRouteMapOk returns a tuple with the ImportRouteMap field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TenantsPutRequestTenantValue) GetImportRouteMapOk() (*string, bool) {
	if o == nil || IsNil(o.ImportRouteMap) {
		return nil, false
	}
	return o.ImportRouteMap, true
}

// HasImportRouteMap returns a boolean if a field has been set.
func (o *TenantsPutRequestTenantValue) HasImportRouteMap() bool {
	if o != nil && !IsNil(o.ImportRouteMap) {
		return true
	}

	return false
}

// SetImportRouteMap gets a reference to the given string and assigns it to the ImportRouteMap field.
func (o *TenantsPutRequestTenantValue) SetImportRouteMap(v string) {
	o.ImportRouteMap = &v
}

// GetImportRouteMapRefType returns the ImportRouteMapRefType field value if set, zero value otherwise.
func (o *TenantsPutRequestTenantValue) GetImportRouteMapRefType() string {
	if o == nil || IsNil(o.ImportRouteMapRefType) {
		var ret string
		return ret
	}
	return *o.ImportRouteMapRefType
}

// GetImportRouteMapRefTypeOk returns a tuple with the ImportRouteMapRefType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TenantsPutRequestTenantValue) GetImportRouteMapRefTypeOk() (*string, bool) {
	if o == nil || IsNil(o.ImportRouteMapRefType) {
		return nil, false
	}
	return o.ImportRouteMapRefType, true
}

// HasImportRouteMapRefType returns a boolean if a field has been set.
func (o *TenantsPutRequestTenantValue) HasImportRouteMapRefType() bool {
	if o != nil && !IsNil(o.ImportRouteMapRefType) {
		return true
	}

	return false
}

// SetImportRouteMapRefType gets a reference to the given string and assigns it to the ImportRouteMapRefType field.
func (o *TenantsPutRequestTenantValue) SetImportRouteMapRefType(v string) {
	o.ImportRouteMapRefType = &v
}

// GetExportRouteMap returns the ExportRouteMap field value if set, zero value otherwise.
func (o *TenantsPutRequestTenantValue) GetExportRouteMap() string {
	if o == nil || IsNil(o.ExportRouteMap) {
		var ret string
		return ret
	}
	return *o.ExportRouteMap
}

// GetExportRouteMapOk returns a tuple with the ExportRouteMap field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TenantsPutRequestTenantValue) GetExportRouteMapOk() (*string, bool) {
	if o == nil || IsNil(o.ExportRouteMap) {
		return nil, false
	}
	return o.ExportRouteMap, true
}

// HasExportRouteMap returns a boolean if a field has been set.
func (o *TenantsPutRequestTenantValue) HasExportRouteMap() bool {
	if o != nil && !IsNil(o.ExportRouteMap) {
		return true
	}

	return false
}

// SetExportRouteMap gets a reference to the given string and assigns it to the ExportRouteMap field.
func (o *TenantsPutRequestTenantValue) SetExportRouteMap(v string) {
	o.ExportRouteMap = &v
}

// GetExportRouteMapRefType returns the ExportRouteMapRefType field value if set, zero value otherwise.
func (o *TenantsPutRequestTenantValue) GetExportRouteMapRefType() string {
	if o == nil || IsNil(o.ExportRouteMapRefType) {
		var ret string
		return ret
	}
	return *o.ExportRouteMapRefType
}

// GetExportRouteMapRefTypeOk returns a tuple with the ExportRouteMapRefType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TenantsPutRequestTenantValue) GetExportRouteMapRefTypeOk() (*string, bool) {
	if o == nil || IsNil(o.ExportRouteMapRefType) {
		return nil, false
	}
	return o.ExportRouteMapRefType, true
}

// HasExportRouteMapRefType returns a boolean if a field has been set.
func (o *TenantsPutRequestTenantValue) HasExportRouteMapRefType() bool {
	if o != nil && !IsNil(o.ExportRouteMapRefType) {
		return true
	}

	return false
}

// SetExportRouteMapRefType gets a reference to the given string and assigns it to the ExportRouteMapRefType field.
func (o *TenantsPutRequestTenantValue) SetExportRouteMapRefType(v string) {
	o.ExportRouteMapRefType = &v
}

// GetVrfName returns the VrfName field value if set, zero value otherwise.
func (o *TenantsPutRequestTenantValue) GetVrfName() string {
	if o == nil || IsNil(o.VrfName) {
		var ret string
		return ret
	}
	return *o.VrfName
}

// GetVrfNameOk returns a tuple with the VrfName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TenantsPutRequestTenantValue) GetVrfNameOk() (*string, bool) {
	if o == nil || IsNil(o.VrfName) {
		return nil, false
	}
	return o.VrfName, true
}

// HasVrfName returns a boolean if a field has been set.
func (o *TenantsPutRequestTenantValue) HasVrfName() bool {
	if o != nil && !IsNil(o.VrfName) {
		return true
	}

	return false
}

// SetVrfName gets a reference to the given string and assigns it to the VrfName field.
func (o *TenantsPutRequestTenantValue) SetVrfName(v string) {
	o.VrfName = &v
}

// GetVrfNameAutoAssigned returns the VrfNameAutoAssigned field value if set, zero value otherwise.
func (o *TenantsPutRequestTenantValue) GetVrfNameAutoAssigned() bool {
	if o == nil || IsNil(o.VrfNameAutoAssigned) {
		var ret bool
		return ret
	}
	return *o.VrfNameAutoAssigned
}

// GetVrfNameAutoAssignedOk returns a tuple with the VrfNameAutoAssigned field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TenantsPutRequestTenantValue) GetVrfNameAutoAssignedOk() (*bool, bool) {
	if o == nil || IsNil(o.VrfNameAutoAssigned) {
		return nil, false
	}
	return o.VrfNameAutoAssigned, true
}

// HasVrfNameAutoAssigned returns a boolean if a field has been set.
func (o *TenantsPutRequestTenantValue) HasVrfNameAutoAssigned() bool {
	if o != nil && !IsNil(o.VrfNameAutoAssigned) {
		return true
	}

	return false
}

// SetVrfNameAutoAssigned gets a reference to the given bool and assigns it to the VrfNameAutoAssigned field.
func (o *TenantsPutRequestTenantValue) SetVrfNameAutoAssigned(v bool) {
	o.VrfNameAutoAssigned = &v
}

// GetRouteTenants returns the RouteTenants field value if set, zero value otherwise.
func (o *TenantsPutRequestTenantValue) GetRouteTenants() []TenantsPutRequestTenantValueRouteTenantsInner {
	if o == nil || IsNil(o.RouteTenants) {
		var ret []TenantsPutRequestTenantValueRouteTenantsInner
		return ret
	}
	return o.RouteTenants
}

// GetRouteTenantsOk returns a tuple with the RouteTenants field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TenantsPutRequestTenantValue) GetRouteTenantsOk() ([]TenantsPutRequestTenantValueRouteTenantsInner, bool) {
	if o == nil || IsNil(o.RouteTenants) {
		return nil, false
	}
	return o.RouteTenants, true
}

// HasRouteTenants returns a boolean if a field has been set.
func (o *TenantsPutRequestTenantValue) HasRouteTenants() bool {
	if o != nil && !IsNil(o.RouteTenants) {
		return true
	}

	return false
}

// SetRouteTenants gets a reference to the given []TenantsPutRequestTenantValueRouteTenantsInner and assigns it to the RouteTenants field.
func (o *TenantsPutRequestTenantValue) SetRouteTenants(v []TenantsPutRequestTenantValueRouteTenantsInner) {
	o.RouteTenants = v
}

// GetObjectProperties returns the ObjectProperties field value if set, zero value otherwise.
func (o *TenantsPutRequestTenantValue) GetObjectProperties() DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties {
	if o == nil || IsNil(o.ObjectProperties) {
		var ret DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties
		return ret
	}
	return *o.ObjectProperties
}

// GetObjectPropertiesOk returns a tuple with the ObjectProperties field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TenantsPutRequestTenantValue) GetObjectPropertiesOk() (*DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties, bool) {
	if o == nil || IsNil(o.ObjectProperties) {
		return nil, false
	}
	return o.ObjectProperties, true
}

// HasObjectProperties returns a boolean if a field has been set.
func (o *TenantsPutRequestTenantValue) HasObjectProperties() bool {
	if o != nil && !IsNil(o.ObjectProperties) {
		return true
	}

	return false
}

// SetObjectProperties gets a reference to the given DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties and assigns it to the ObjectProperties field.
func (o *TenantsPutRequestTenantValue) SetObjectProperties(v DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties) {
	o.ObjectProperties = &v
}

// GetDefaultOriginate returns the DefaultOriginate field value if set, zero value otherwise.
func (o *TenantsPutRequestTenantValue) GetDefaultOriginate() bool {
	if o == nil || IsNil(o.DefaultOriginate) {
		var ret bool
		return ret
	}
	return *o.DefaultOriginate
}

// GetDefaultOriginateOk returns a tuple with the DefaultOriginate field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TenantsPutRequestTenantValue) GetDefaultOriginateOk() (*bool, bool) {
	if o == nil || IsNil(o.DefaultOriginate) {
		return nil, false
	}
	return o.DefaultOriginate, true
}

// HasDefaultOriginate returns a boolean if a field has been set.
func (o *TenantsPutRequestTenantValue) HasDefaultOriginate() bool {
	if o != nil && !IsNil(o.DefaultOriginate) {
		return true
	}

	return false
}

// SetDefaultOriginate gets a reference to the given bool and assigns it to the DefaultOriginate field.
func (o *TenantsPutRequestTenantValue) SetDefaultOriginate(v bool) {
	o.DefaultOriginate = &v
}

func (o TenantsPutRequestTenantValue) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TenantsPutRequestTenantValue) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.Enable) {
		toSerialize["enable"] = o.Enable
	}
	if o.Layer3Vni.IsSet() {
		toSerialize["layer_3_vni"] = o.Layer3Vni.Get()
	}
	if !IsNil(o.Layer3VniAutoAssigned) {
		toSerialize["layer_3_vni_auto_assigned_"] = o.Layer3VniAutoAssigned
	}
	if o.Layer3Vlan.IsSet() {
		toSerialize["layer_3_vlan"] = o.Layer3Vlan.Get()
	}
	if !IsNil(o.Layer3VlanAutoAssigned) {
		toSerialize["layer_3_vlan_auto_assigned_"] = o.Layer3VlanAutoAssigned
	}
	if !IsNil(o.DhcpRelaySourceIpv4sSubnet) {
		toSerialize["dhcp_relay_source_ipv4s_subnet"] = o.DhcpRelaySourceIpv4sSubnet
	}
	if !IsNil(o.DhcpRelaySourceIpv6sSubnet) {
		toSerialize["dhcp_relay_source_ipv6s_subnet"] = o.DhcpRelaySourceIpv6sSubnet
	}
	if !IsNil(o.RouteDistinguisher) {
		toSerialize["route_distinguisher"] = o.RouteDistinguisher
	}
	if !IsNil(o.RouteTargetImport) {
		toSerialize["route_target_import"] = o.RouteTargetImport
	}
	if !IsNil(o.RouteTargetExport) {
		toSerialize["route_target_export"] = o.RouteTargetExport
	}
	if !IsNil(o.ImportRouteMap) {
		toSerialize["import_route_map"] = o.ImportRouteMap
	}
	if !IsNil(o.ImportRouteMapRefType) {
		toSerialize["import_route_map_ref_type_"] = o.ImportRouteMapRefType
	}
	if !IsNil(o.ExportRouteMap) {
		toSerialize["export_route_map"] = o.ExportRouteMap
	}
	if !IsNil(o.ExportRouteMapRefType) {
		toSerialize["export_route_map_ref_type_"] = o.ExportRouteMapRefType
	}
	if !IsNil(o.VrfName) {
		toSerialize["vrf_name"] = o.VrfName
	}
	if !IsNil(o.VrfNameAutoAssigned) {
		toSerialize["vrf_name_auto_assigned_"] = o.VrfNameAutoAssigned
	}
	if !IsNil(o.RouteTenants) {
		toSerialize["route_tenants"] = o.RouteTenants
	}
	if !IsNil(o.ObjectProperties) {
		toSerialize["object_properties"] = o.ObjectProperties
	}
	if !IsNil(o.DefaultOriginate) {
		toSerialize["default_originate"] = o.DefaultOriginate
	}
	return toSerialize, nil
}

type NullableTenantsPutRequestTenantValue struct {
	value *TenantsPutRequestTenantValue
	isSet bool
}

func (v NullableTenantsPutRequestTenantValue) Get() *TenantsPutRequestTenantValue {
	return v.value
}

func (v *NullableTenantsPutRequestTenantValue) Set(val *TenantsPutRequestTenantValue) {
	v.value = val
	v.isSet = true
}

func (v NullableTenantsPutRequestTenantValue) IsSet() bool {
	return v.isSet
}

func (v *NullableTenantsPutRequestTenantValue) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTenantsPutRequestTenantValue(val *TenantsPutRequestTenantValue) *NullableTenantsPutRequestTenantValue {
	return &NullableTenantsPutRequestTenantValue{value: val, isSet: true}
}

func (v NullableTenantsPutRequestTenantValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTenantsPutRequestTenantValue) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


