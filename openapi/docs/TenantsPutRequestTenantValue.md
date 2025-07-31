# TenantsPutRequestTenantValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. It&#39;s highly recommended to set this value to true so that validation on the object will be ran. | [optional] [default to true]
**Layer3Vni** | Pointer to **NullableInt32** | VNI value used to transport traffic between services of a Tenant  | [optional] 
**Layer3VniAutoAssigned** | Pointer to **bool** | Whether or not the value in layer_3_vni field has been automatically assigned or not. Set to false and change layer_3_vni value to edit. | [optional] 
**Layer3Vlan** | Pointer to **NullableInt32** | VLAN value used to transport traffic between services of a Tenant  | [optional] 
**Layer3VlanAutoAssigned** | Pointer to **bool** | Whether or not the value in layer_3_vlan field has been automatically assigned or not. Set to false and change layer_3_vlan value to edit. | [optional] 
**DhcpRelaySourceIpv4sSubnet** | Pointer to **string** | Range of IPv4 addresses (represented in IPv4 subnet format) used to configure the source IP of each DHCP Relay on each switch that this Tenant is provisioned on. | [optional] [default to ""]
**DhcpRelaySourceIpv6sSubnet** | Pointer to **string** | Range of IPv6 addresses (represented in IPv6 subnet format) used to configure the source IP of each DHCP Relay on each switch that this Tenant is provisioned on. | [optional] [default to ""]
**RouteDistinguisher** | Pointer to **string** | Route Distinguishers are used to maintain uniqueness among identical routes from different routers.  If set, then routes from this Tenant will be identified with this Route Distinguisher (BGP Community).  It should be two numbers separated by a colon. | [optional] [default to ""]
**RouteTargetImport** | Pointer to **string** | A route-target (BGP Community) to attach while importing routes into the current tenant. It should be a comma-separated list of BGP Communities: each Community being two numbers separated by a colon. | [optional] [default to ""]
**RouteTargetExport** | Pointer to **string** | A route-target (BGP Community) to attach while exporting routes from the current tenant. It should be a comma-separated list of BGP Communities: each Community being two numbers separated by a colon. | [optional] [default to ""]
**ImportRouteMap** | Pointer to **string** | A route-map applied to routes imported into the current tenant from other tenants with the purpose of filtering or modifying the routes | [optional] [default to ""]
**ImportRouteMapRefType** | Pointer to **string** | Object type for import_route_map field | [optional] 
**ExportRouteMap** | Pointer to **string** | A route-map applied to routes exported into the current tenant from other tenants with the purpose of filtering or modifying the routes | [optional] [default to ""]
**ExportRouteMapRefType** | Pointer to **string** | Object type for export_route_map field | [optional] 
**VrfName** | Pointer to **string** | Virtual Routing and Forwarding instance name associated to tenants  | [optional] [default to "(auto)"]
**VrfNameAutoAssigned** | Pointer to **bool** | Whether or not the value in vrf_name field has been automatically assigned or not. Set to false and change vrf_name value to edit. | [optional] 
**RouteTenants** | Pointer to [**[]TenantsPutRequestTenantValueRouteTenantsInner**](TenantsPutRequestTenantValueRouteTenantsInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties**](DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties.md) |  | [optional] 
**DefaultOriginate** | Pointer to **bool** | Enables a leaf switch to originate IPv4 default type-5 EVPN routes across the switching fabric. | [optional] [default to false]

## Methods

### NewTenantsPutRequestTenantValue

`func NewTenantsPutRequestTenantValue() *TenantsPutRequestTenantValue`

NewTenantsPutRequestTenantValue instantiates a new TenantsPutRequestTenantValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTenantsPutRequestTenantValueWithDefaults

`func NewTenantsPutRequestTenantValueWithDefaults() *TenantsPutRequestTenantValue`

NewTenantsPutRequestTenantValueWithDefaults instantiates a new TenantsPutRequestTenantValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *TenantsPutRequestTenantValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *TenantsPutRequestTenantValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *TenantsPutRequestTenantValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *TenantsPutRequestTenantValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *TenantsPutRequestTenantValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *TenantsPutRequestTenantValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *TenantsPutRequestTenantValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *TenantsPutRequestTenantValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetLayer3Vni

`func (o *TenantsPutRequestTenantValue) GetLayer3Vni() int32`

GetLayer3Vni returns the Layer3Vni field if non-nil, zero value otherwise.

### GetLayer3VniOk

`func (o *TenantsPutRequestTenantValue) GetLayer3VniOk() (*int32, bool)`

GetLayer3VniOk returns a tuple with the Layer3Vni field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLayer3Vni

`func (o *TenantsPutRequestTenantValue) SetLayer3Vni(v int32)`

SetLayer3Vni sets Layer3Vni field to given value.

### HasLayer3Vni

`func (o *TenantsPutRequestTenantValue) HasLayer3Vni() bool`

HasLayer3Vni returns a boolean if a field has been set.

### SetLayer3VniNil

`func (o *TenantsPutRequestTenantValue) SetLayer3VniNil(b bool)`

 SetLayer3VniNil sets the value for Layer3Vni to be an explicit nil

### UnsetLayer3Vni
`func (o *TenantsPutRequestTenantValue) UnsetLayer3Vni()`

UnsetLayer3Vni ensures that no value is present for Layer3Vni, not even an explicit nil
### GetLayer3VniAutoAssigned

`func (o *TenantsPutRequestTenantValue) GetLayer3VniAutoAssigned() bool`

GetLayer3VniAutoAssigned returns the Layer3VniAutoAssigned field if non-nil, zero value otherwise.

### GetLayer3VniAutoAssignedOk

`func (o *TenantsPutRequestTenantValue) GetLayer3VniAutoAssignedOk() (*bool, bool)`

GetLayer3VniAutoAssignedOk returns a tuple with the Layer3VniAutoAssigned field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLayer3VniAutoAssigned

`func (o *TenantsPutRequestTenantValue) SetLayer3VniAutoAssigned(v bool)`

SetLayer3VniAutoAssigned sets Layer3VniAutoAssigned field to given value.

### HasLayer3VniAutoAssigned

`func (o *TenantsPutRequestTenantValue) HasLayer3VniAutoAssigned() bool`

HasLayer3VniAutoAssigned returns a boolean if a field has been set.

### GetLayer3Vlan

`func (o *TenantsPutRequestTenantValue) GetLayer3Vlan() int32`

GetLayer3Vlan returns the Layer3Vlan field if non-nil, zero value otherwise.

### GetLayer3VlanOk

`func (o *TenantsPutRequestTenantValue) GetLayer3VlanOk() (*int32, bool)`

GetLayer3VlanOk returns a tuple with the Layer3Vlan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLayer3Vlan

`func (o *TenantsPutRequestTenantValue) SetLayer3Vlan(v int32)`

SetLayer3Vlan sets Layer3Vlan field to given value.

### HasLayer3Vlan

`func (o *TenantsPutRequestTenantValue) HasLayer3Vlan() bool`

HasLayer3Vlan returns a boolean if a field has been set.

### SetLayer3VlanNil

`func (o *TenantsPutRequestTenantValue) SetLayer3VlanNil(b bool)`

 SetLayer3VlanNil sets the value for Layer3Vlan to be an explicit nil

### UnsetLayer3Vlan
`func (o *TenantsPutRequestTenantValue) UnsetLayer3Vlan()`

UnsetLayer3Vlan ensures that no value is present for Layer3Vlan, not even an explicit nil
### GetLayer3VlanAutoAssigned

`func (o *TenantsPutRequestTenantValue) GetLayer3VlanAutoAssigned() bool`

GetLayer3VlanAutoAssigned returns the Layer3VlanAutoAssigned field if non-nil, zero value otherwise.

### GetLayer3VlanAutoAssignedOk

`func (o *TenantsPutRequestTenantValue) GetLayer3VlanAutoAssignedOk() (*bool, bool)`

GetLayer3VlanAutoAssignedOk returns a tuple with the Layer3VlanAutoAssigned field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLayer3VlanAutoAssigned

`func (o *TenantsPutRequestTenantValue) SetLayer3VlanAutoAssigned(v bool)`

SetLayer3VlanAutoAssigned sets Layer3VlanAutoAssigned field to given value.

### HasLayer3VlanAutoAssigned

`func (o *TenantsPutRequestTenantValue) HasLayer3VlanAutoAssigned() bool`

HasLayer3VlanAutoAssigned returns a boolean if a field has been set.

### GetDhcpRelaySourceIpv4sSubnet

`func (o *TenantsPutRequestTenantValue) GetDhcpRelaySourceIpv4sSubnet() string`

GetDhcpRelaySourceIpv4sSubnet returns the DhcpRelaySourceIpv4sSubnet field if non-nil, zero value otherwise.

### GetDhcpRelaySourceIpv4sSubnetOk

`func (o *TenantsPutRequestTenantValue) GetDhcpRelaySourceIpv4sSubnetOk() (*string, bool)`

GetDhcpRelaySourceIpv4sSubnetOk returns a tuple with the DhcpRelaySourceIpv4sSubnet field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDhcpRelaySourceIpv4sSubnet

`func (o *TenantsPutRequestTenantValue) SetDhcpRelaySourceIpv4sSubnet(v string)`

SetDhcpRelaySourceIpv4sSubnet sets DhcpRelaySourceIpv4sSubnet field to given value.

### HasDhcpRelaySourceIpv4sSubnet

`func (o *TenantsPutRequestTenantValue) HasDhcpRelaySourceIpv4sSubnet() bool`

HasDhcpRelaySourceIpv4sSubnet returns a boolean if a field has been set.

### GetDhcpRelaySourceIpv6sSubnet

`func (o *TenantsPutRequestTenantValue) GetDhcpRelaySourceIpv6sSubnet() string`

GetDhcpRelaySourceIpv6sSubnet returns the DhcpRelaySourceIpv6sSubnet field if non-nil, zero value otherwise.

### GetDhcpRelaySourceIpv6sSubnetOk

`func (o *TenantsPutRequestTenantValue) GetDhcpRelaySourceIpv6sSubnetOk() (*string, bool)`

GetDhcpRelaySourceIpv6sSubnetOk returns a tuple with the DhcpRelaySourceIpv6sSubnet field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDhcpRelaySourceIpv6sSubnet

`func (o *TenantsPutRequestTenantValue) SetDhcpRelaySourceIpv6sSubnet(v string)`

SetDhcpRelaySourceIpv6sSubnet sets DhcpRelaySourceIpv6sSubnet field to given value.

### HasDhcpRelaySourceIpv6sSubnet

`func (o *TenantsPutRequestTenantValue) HasDhcpRelaySourceIpv6sSubnet() bool`

HasDhcpRelaySourceIpv6sSubnet returns a boolean if a field has been set.

### GetRouteDistinguisher

`func (o *TenantsPutRequestTenantValue) GetRouteDistinguisher() string`

GetRouteDistinguisher returns the RouteDistinguisher field if non-nil, zero value otherwise.

### GetRouteDistinguisherOk

`func (o *TenantsPutRequestTenantValue) GetRouteDistinguisherOk() (*string, bool)`

GetRouteDistinguisherOk returns a tuple with the RouteDistinguisher field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRouteDistinguisher

`func (o *TenantsPutRequestTenantValue) SetRouteDistinguisher(v string)`

SetRouteDistinguisher sets RouteDistinguisher field to given value.

### HasRouteDistinguisher

`func (o *TenantsPutRequestTenantValue) HasRouteDistinguisher() bool`

HasRouteDistinguisher returns a boolean if a field has been set.

### GetRouteTargetImport

`func (o *TenantsPutRequestTenantValue) GetRouteTargetImport() string`

GetRouteTargetImport returns the RouteTargetImport field if non-nil, zero value otherwise.

### GetRouteTargetImportOk

`func (o *TenantsPutRequestTenantValue) GetRouteTargetImportOk() (*string, bool)`

GetRouteTargetImportOk returns a tuple with the RouteTargetImport field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRouteTargetImport

`func (o *TenantsPutRequestTenantValue) SetRouteTargetImport(v string)`

SetRouteTargetImport sets RouteTargetImport field to given value.

### HasRouteTargetImport

`func (o *TenantsPutRequestTenantValue) HasRouteTargetImport() bool`

HasRouteTargetImport returns a boolean if a field has been set.

### GetRouteTargetExport

`func (o *TenantsPutRequestTenantValue) GetRouteTargetExport() string`

GetRouteTargetExport returns the RouteTargetExport field if non-nil, zero value otherwise.

### GetRouteTargetExportOk

`func (o *TenantsPutRequestTenantValue) GetRouteTargetExportOk() (*string, bool)`

GetRouteTargetExportOk returns a tuple with the RouteTargetExport field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRouteTargetExport

`func (o *TenantsPutRequestTenantValue) SetRouteTargetExport(v string)`

SetRouteTargetExport sets RouteTargetExport field to given value.

### HasRouteTargetExport

`func (o *TenantsPutRequestTenantValue) HasRouteTargetExport() bool`

HasRouteTargetExport returns a boolean if a field has been set.

### GetImportRouteMap

`func (o *TenantsPutRequestTenantValue) GetImportRouteMap() string`

GetImportRouteMap returns the ImportRouteMap field if non-nil, zero value otherwise.

### GetImportRouteMapOk

`func (o *TenantsPutRequestTenantValue) GetImportRouteMapOk() (*string, bool)`

GetImportRouteMapOk returns a tuple with the ImportRouteMap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImportRouteMap

`func (o *TenantsPutRequestTenantValue) SetImportRouteMap(v string)`

SetImportRouteMap sets ImportRouteMap field to given value.

### HasImportRouteMap

`func (o *TenantsPutRequestTenantValue) HasImportRouteMap() bool`

HasImportRouteMap returns a boolean if a field has been set.

### GetImportRouteMapRefType

`func (o *TenantsPutRequestTenantValue) GetImportRouteMapRefType() string`

GetImportRouteMapRefType returns the ImportRouteMapRefType field if non-nil, zero value otherwise.

### GetImportRouteMapRefTypeOk

`func (o *TenantsPutRequestTenantValue) GetImportRouteMapRefTypeOk() (*string, bool)`

GetImportRouteMapRefTypeOk returns a tuple with the ImportRouteMapRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImportRouteMapRefType

`func (o *TenantsPutRequestTenantValue) SetImportRouteMapRefType(v string)`

SetImportRouteMapRefType sets ImportRouteMapRefType field to given value.

### HasImportRouteMapRefType

`func (o *TenantsPutRequestTenantValue) HasImportRouteMapRefType() bool`

HasImportRouteMapRefType returns a boolean if a field has been set.

### GetExportRouteMap

`func (o *TenantsPutRequestTenantValue) GetExportRouteMap() string`

GetExportRouteMap returns the ExportRouteMap field if non-nil, zero value otherwise.

### GetExportRouteMapOk

`func (o *TenantsPutRequestTenantValue) GetExportRouteMapOk() (*string, bool)`

GetExportRouteMapOk returns a tuple with the ExportRouteMap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExportRouteMap

`func (o *TenantsPutRequestTenantValue) SetExportRouteMap(v string)`

SetExportRouteMap sets ExportRouteMap field to given value.

### HasExportRouteMap

`func (o *TenantsPutRequestTenantValue) HasExportRouteMap() bool`

HasExportRouteMap returns a boolean if a field has been set.

### GetExportRouteMapRefType

`func (o *TenantsPutRequestTenantValue) GetExportRouteMapRefType() string`

GetExportRouteMapRefType returns the ExportRouteMapRefType field if non-nil, zero value otherwise.

### GetExportRouteMapRefTypeOk

`func (o *TenantsPutRequestTenantValue) GetExportRouteMapRefTypeOk() (*string, bool)`

GetExportRouteMapRefTypeOk returns a tuple with the ExportRouteMapRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExportRouteMapRefType

`func (o *TenantsPutRequestTenantValue) SetExportRouteMapRefType(v string)`

SetExportRouteMapRefType sets ExportRouteMapRefType field to given value.

### HasExportRouteMapRefType

`func (o *TenantsPutRequestTenantValue) HasExportRouteMapRefType() bool`

HasExportRouteMapRefType returns a boolean if a field has been set.

### GetVrfName

`func (o *TenantsPutRequestTenantValue) GetVrfName() string`

GetVrfName returns the VrfName field if non-nil, zero value otherwise.

### GetVrfNameOk

`func (o *TenantsPutRequestTenantValue) GetVrfNameOk() (*string, bool)`

GetVrfNameOk returns a tuple with the VrfName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVrfName

`func (o *TenantsPutRequestTenantValue) SetVrfName(v string)`

SetVrfName sets VrfName field to given value.

### HasVrfName

`func (o *TenantsPutRequestTenantValue) HasVrfName() bool`

HasVrfName returns a boolean if a field has been set.

### GetVrfNameAutoAssigned

`func (o *TenantsPutRequestTenantValue) GetVrfNameAutoAssigned() bool`

GetVrfNameAutoAssigned returns the VrfNameAutoAssigned field if non-nil, zero value otherwise.

### GetVrfNameAutoAssignedOk

`func (o *TenantsPutRequestTenantValue) GetVrfNameAutoAssignedOk() (*bool, bool)`

GetVrfNameAutoAssignedOk returns a tuple with the VrfNameAutoAssigned field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVrfNameAutoAssigned

`func (o *TenantsPutRequestTenantValue) SetVrfNameAutoAssigned(v bool)`

SetVrfNameAutoAssigned sets VrfNameAutoAssigned field to given value.

### HasVrfNameAutoAssigned

`func (o *TenantsPutRequestTenantValue) HasVrfNameAutoAssigned() bool`

HasVrfNameAutoAssigned returns a boolean if a field has been set.

### GetRouteTenants

`func (o *TenantsPutRequestTenantValue) GetRouteTenants() []TenantsPutRequestTenantValueRouteTenantsInner`

GetRouteTenants returns the RouteTenants field if non-nil, zero value otherwise.

### GetRouteTenantsOk

`func (o *TenantsPutRequestTenantValue) GetRouteTenantsOk() (*[]TenantsPutRequestTenantValueRouteTenantsInner, bool)`

GetRouteTenantsOk returns a tuple with the RouteTenants field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRouteTenants

`func (o *TenantsPutRequestTenantValue) SetRouteTenants(v []TenantsPutRequestTenantValueRouteTenantsInner)`

SetRouteTenants sets RouteTenants field to given value.

### HasRouteTenants

`func (o *TenantsPutRequestTenantValue) HasRouteTenants() bool`

HasRouteTenants returns a boolean if a field has been set.

### GetObjectProperties

`func (o *TenantsPutRequestTenantValue) GetObjectProperties() DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *TenantsPutRequestTenantValue) GetObjectPropertiesOk() (*DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *TenantsPutRequestTenantValue) SetObjectProperties(v DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *TenantsPutRequestTenantValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.

### GetDefaultOriginate

`func (o *TenantsPutRequestTenantValue) GetDefaultOriginate() bool`

GetDefaultOriginate returns the DefaultOriginate field if non-nil, zero value otherwise.

### GetDefaultOriginateOk

`func (o *TenantsPutRequestTenantValue) GetDefaultOriginateOk() (*bool, bool)`

GetDefaultOriginateOk returns a tuple with the DefaultOriginate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefaultOriginate

`func (o *TenantsPutRequestTenantValue) SetDefaultOriginate(v bool)`

SetDefaultOriginate sets DefaultOriginate field to given value.

### HasDefaultOriginate

`func (o *TenantsPutRequestTenantValue) HasDefaultOriginate() bool`

HasDefaultOriginate returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


