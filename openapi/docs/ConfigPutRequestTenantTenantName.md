# ConfigPutRequestTenantTenantName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to true]
**Layer3Vni** | Pointer to **int32** | VNI value used to transport traffic between services of a Tenant  | [optional] 
**Layer3VniAutoAssigned** | Pointer to **bool** | Whether or not the value in layer_3_vni field has been automatically assigned or not. Set to false and change layer_3_vni value to edit. | [optional] 
**Layer3Vlan** | Pointer to **int32** | VLAN value used to transport traffic between services of a Tenant  | [optional] 
**Layer3VlanAutoAssigned** | Pointer to **bool** | Whether or not the value in layer_3_vlan field has been automatically assigned or not. Set to false and change layer_3_vlan value to edit. | [optional] 
**DhcpRelaySourceIpsSubnet** | Pointer to **string** | Range of IP addresses (represented in IP subnet format) used to configure the source IP of each DHCP Relay on each switch that this Tenant is provisioned on. | [optional] [default to ""]
**RouteDistinguisher** | Pointer to **string** | Route Distinguishers are used to maintain uniqueness among identical routes from different routers.  If set, then routes from this Tenant will be identified with this Route Distinguisher (BGP Community).  It should be two numbers separated by a colon. | [optional] [default to ""]
**RouteTargetImport** | Pointer to **string** | A route-target (BGP Community) to attach while importing routes into the current tenant. It should be a comma-separated list of BGP Communities: each Community being two numbers separated by a colon. | [optional] [default to ""]
**RouteTargetExport** | Pointer to **string** | A route-target (BGP Community) to attach while exporting routes from the current tenant. It should be a comma-separated list of BGP Communities: each Community being two numbers separated by a colon. | [optional] [default to ""]
**ImportRouteMap** | Pointer to **string** | A route-map applied to routes imported into the current tenant from other tenants with the purpose of filtering or modifying the routes | [optional] [default to ""]
**ImportRouteMapRefType** | Pointer to **string** | Object type for import_route_map field | [optional] 
**ExportRouteMap** | Pointer to **string** | A route-map applied to routes exported into the current tenant from other tenants with the purpose of filtering or modifying the routes | [optional] [default to ""]
**ExportRouteMapRefType** | Pointer to **string** | Object type for export_route_map field | [optional] 
**VrfName** | Pointer to **string** | Virtual Routing and Forwarding instance name associated to tenants  | [optional] [default to "(auto)"]
**VrfNameAutoAssigned** | Pointer to **bool** | Whether or not the value in vrf_name field has been automatically assigned or not. Set to false and change vrf_name value to edit. | [optional] 
**RouteTenants** | Pointer to [**[]ConfigPutRequestTenantTenantNameRouteTenantsInner**](ConfigPutRequestTenantTenantNameRouteTenantsInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties**](ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties.md) |  | [optional] 

## Methods

### NewConfigPutRequestTenantTenantName

`func NewConfigPutRequestTenantTenantName() *ConfigPutRequestTenantTenantName`

NewConfigPutRequestTenantTenantName instantiates a new ConfigPutRequestTenantTenantName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestTenantTenantNameWithDefaults

`func NewConfigPutRequestTenantTenantNameWithDefaults() *ConfigPutRequestTenantTenantName`

NewConfigPutRequestTenantTenantNameWithDefaults instantiates a new ConfigPutRequestTenantTenantName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestTenantTenantName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestTenantTenantName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestTenantTenantName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestTenantTenantName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ConfigPutRequestTenantTenantName) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestTenantTenantName) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestTenantTenantName) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestTenantTenantName) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetLayer3Vni

`func (o *ConfigPutRequestTenantTenantName) GetLayer3Vni() int32`

GetLayer3Vni returns the Layer3Vni field if non-nil, zero value otherwise.

### GetLayer3VniOk

`func (o *ConfigPutRequestTenantTenantName) GetLayer3VniOk() (*int32, bool)`

GetLayer3VniOk returns a tuple with the Layer3Vni field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLayer3Vni

`func (o *ConfigPutRequestTenantTenantName) SetLayer3Vni(v int32)`

SetLayer3Vni sets Layer3Vni field to given value.

### HasLayer3Vni

`func (o *ConfigPutRequestTenantTenantName) HasLayer3Vni() bool`

HasLayer3Vni returns a boolean if a field has been set.

### GetLayer3VniAutoAssigned

`func (o *ConfigPutRequestTenantTenantName) GetLayer3VniAutoAssigned() bool`

GetLayer3VniAutoAssigned returns the Layer3VniAutoAssigned field if non-nil, zero value otherwise.

### GetLayer3VniAutoAssignedOk

`func (o *ConfigPutRequestTenantTenantName) GetLayer3VniAutoAssignedOk() (*bool, bool)`

GetLayer3VniAutoAssignedOk returns a tuple with the Layer3VniAutoAssigned field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLayer3VniAutoAssigned

`func (o *ConfigPutRequestTenantTenantName) SetLayer3VniAutoAssigned(v bool)`

SetLayer3VniAutoAssigned sets Layer3VniAutoAssigned field to given value.

### HasLayer3VniAutoAssigned

`func (o *ConfigPutRequestTenantTenantName) HasLayer3VniAutoAssigned() bool`

HasLayer3VniAutoAssigned returns a boolean if a field has been set.

### GetLayer3Vlan

`func (o *ConfigPutRequestTenantTenantName) GetLayer3Vlan() int32`

GetLayer3Vlan returns the Layer3Vlan field if non-nil, zero value otherwise.

### GetLayer3VlanOk

`func (o *ConfigPutRequestTenantTenantName) GetLayer3VlanOk() (*int32, bool)`

GetLayer3VlanOk returns a tuple with the Layer3Vlan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLayer3Vlan

`func (o *ConfigPutRequestTenantTenantName) SetLayer3Vlan(v int32)`

SetLayer3Vlan sets Layer3Vlan field to given value.

### HasLayer3Vlan

`func (o *ConfigPutRequestTenantTenantName) HasLayer3Vlan() bool`

HasLayer3Vlan returns a boolean if a field has been set.

### GetLayer3VlanAutoAssigned

`func (o *ConfigPutRequestTenantTenantName) GetLayer3VlanAutoAssigned() bool`

GetLayer3VlanAutoAssigned returns the Layer3VlanAutoAssigned field if non-nil, zero value otherwise.

### GetLayer3VlanAutoAssignedOk

`func (o *ConfigPutRequestTenantTenantName) GetLayer3VlanAutoAssignedOk() (*bool, bool)`

GetLayer3VlanAutoAssignedOk returns a tuple with the Layer3VlanAutoAssigned field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLayer3VlanAutoAssigned

`func (o *ConfigPutRequestTenantTenantName) SetLayer3VlanAutoAssigned(v bool)`

SetLayer3VlanAutoAssigned sets Layer3VlanAutoAssigned field to given value.

### HasLayer3VlanAutoAssigned

`func (o *ConfigPutRequestTenantTenantName) HasLayer3VlanAutoAssigned() bool`

HasLayer3VlanAutoAssigned returns a boolean if a field has been set.

### GetDhcpRelaySourceIpsSubnet

`func (o *ConfigPutRequestTenantTenantName) GetDhcpRelaySourceIpsSubnet() string`

GetDhcpRelaySourceIpsSubnet returns the DhcpRelaySourceIpsSubnet field if non-nil, zero value otherwise.

### GetDhcpRelaySourceIpsSubnetOk

`func (o *ConfigPutRequestTenantTenantName) GetDhcpRelaySourceIpsSubnetOk() (*string, bool)`

GetDhcpRelaySourceIpsSubnetOk returns a tuple with the DhcpRelaySourceIpsSubnet field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDhcpRelaySourceIpsSubnet

`func (o *ConfigPutRequestTenantTenantName) SetDhcpRelaySourceIpsSubnet(v string)`

SetDhcpRelaySourceIpsSubnet sets DhcpRelaySourceIpsSubnet field to given value.

### HasDhcpRelaySourceIpsSubnet

`func (o *ConfigPutRequestTenantTenantName) HasDhcpRelaySourceIpsSubnet() bool`

HasDhcpRelaySourceIpsSubnet returns a boolean if a field has been set.

### GetRouteDistinguisher

`func (o *ConfigPutRequestTenantTenantName) GetRouteDistinguisher() string`

GetRouteDistinguisher returns the RouteDistinguisher field if non-nil, zero value otherwise.

### GetRouteDistinguisherOk

`func (o *ConfigPutRequestTenantTenantName) GetRouteDistinguisherOk() (*string, bool)`

GetRouteDistinguisherOk returns a tuple with the RouteDistinguisher field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRouteDistinguisher

`func (o *ConfigPutRequestTenantTenantName) SetRouteDistinguisher(v string)`

SetRouteDistinguisher sets RouteDistinguisher field to given value.

### HasRouteDistinguisher

`func (o *ConfigPutRequestTenantTenantName) HasRouteDistinguisher() bool`

HasRouteDistinguisher returns a boolean if a field has been set.

### GetRouteTargetImport

`func (o *ConfigPutRequestTenantTenantName) GetRouteTargetImport() string`

GetRouteTargetImport returns the RouteTargetImport field if non-nil, zero value otherwise.

### GetRouteTargetImportOk

`func (o *ConfigPutRequestTenantTenantName) GetRouteTargetImportOk() (*string, bool)`

GetRouteTargetImportOk returns a tuple with the RouteTargetImport field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRouteTargetImport

`func (o *ConfigPutRequestTenantTenantName) SetRouteTargetImport(v string)`

SetRouteTargetImport sets RouteTargetImport field to given value.

### HasRouteTargetImport

`func (o *ConfigPutRequestTenantTenantName) HasRouteTargetImport() bool`

HasRouteTargetImport returns a boolean if a field has been set.

### GetRouteTargetExport

`func (o *ConfigPutRequestTenantTenantName) GetRouteTargetExport() string`

GetRouteTargetExport returns the RouteTargetExport field if non-nil, zero value otherwise.

### GetRouteTargetExportOk

`func (o *ConfigPutRequestTenantTenantName) GetRouteTargetExportOk() (*string, bool)`

GetRouteTargetExportOk returns a tuple with the RouteTargetExport field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRouteTargetExport

`func (o *ConfigPutRequestTenantTenantName) SetRouteTargetExport(v string)`

SetRouteTargetExport sets RouteTargetExport field to given value.

### HasRouteTargetExport

`func (o *ConfigPutRequestTenantTenantName) HasRouteTargetExport() bool`

HasRouteTargetExport returns a boolean if a field has been set.

### GetImportRouteMap

`func (o *ConfigPutRequestTenantTenantName) GetImportRouteMap() string`

GetImportRouteMap returns the ImportRouteMap field if non-nil, zero value otherwise.

### GetImportRouteMapOk

`func (o *ConfigPutRequestTenantTenantName) GetImportRouteMapOk() (*string, bool)`

GetImportRouteMapOk returns a tuple with the ImportRouteMap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImportRouteMap

`func (o *ConfigPutRequestTenantTenantName) SetImportRouteMap(v string)`

SetImportRouteMap sets ImportRouteMap field to given value.

### HasImportRouteMap

`func (o *ConfigPutRequestTenantTenantName) HasImportRouteMap() bool`

HasImportRouteMap returns a boolean if a field has been set.

### GetImportRouteMapRefType

`func (o *ConfigPutRequestTenantTenantName) GetImportRouteMapRefType() string`

GetImportRouteMapRefType returns the ImportRouteMapRefType field if non-nil, zero value otherwise.

### GetImportRouteMapRefTypeOk

`func (o *ConfigPutRequestTenantTenantName) GetImportRouteMapRefTypeOk() (*string, bool)`

GetImportRouteMapRefTypeOk returns a tuple with the ImportRouteMapRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImportRouteMapRefType

`func (o *ConfigPutRequestTenantTenantName) SetImportRouteMapRefType(v string)`

SetImportRouteMapRefType sets ImportRouteMapRefType field to given value.

### HasImportRouteMapRefType

`func (o *ConfigPutRequestTenantTenantName) HasImportRouteMapRefType() bool`

HasImportRouteMapRefType returns a boolean if a field has been set.

### GetExportRouteMap

`func (o *ConfigPutRequestTenantTenantName) GetExportRouteMap() string`

GetExportRouteMap returns the ExportRouteMap field if non-nil, zero value otherwise.

### GetExportRouteMapOk

`func (o *ConfigPutRequestTenantTenantName) GetExportRouteMapOk() (*string, bool)`

GetExportRouteMapOk returns a tuple with the ExportRouteMap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExportRouteMap

`func (o *ConfigPutRequestTenantTenantName) SetExportRouteMap(v string)`

SetExportRouteMap sets ExportRouteMap field to given value.

### HasExportRouteMap

`func (o *ConfigPutRequestTenantTenantName) HasExportRouteMap() bool`

HasExportRouteMap returns a boolean if a field has been set.

### GetExportRouteMapRefType

`func (o *ConfigPutRequestTenantTenantName) GetExportRouteMapRefType() string`

GetExportRouteMapRefType returns the ExportRouteMapRefType field if non-nil, zero value otherwise.

### GetExportRouteMapRefTypeOk

`func (o *ConfigPutRequestTenantTenantName) GetExportRouteMapRefTypeOk() (*string, bool)`

GetExportRouteMapRefTypeOk returns a tuple with the ExportRouteMapRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExportRouteMapRefType

`func (o *ConfigPutRequestTenantTenantName) SetExportRouteMapRefType(v string)`

SetExportRouteMapRefType sets ExportRouteMapRefType field to given value.

### HasExportRouteMapRefType

`func (o *ConfigPutRequestTenantTenantName) HasExportRouteMapRefType() bool`

HasExportRouteMapRefType returns a boolean if a field has been set.

### GetVrfName

`func (o *ConfigPutRequestTenantTenantName) GetVrfName() string`

GetVrfName returns the VrfName field if non-nil, zero value otherwise.

### GetVrfNameOk

`func (o *ConfigPutRequestTenantTenantName) GetVrfNameOk() (*string, bool)`

GetVrfNameOk returns a tuple with the VrfName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVrfName

`func (o *ConfigPutRequestTenantTenantName) SetVrfName(v string)`

SetVrfName sets VrfName field to given value.

### HasVrfName

`func (o *ConfigPutRequestTenantTenantName) HasVrfName() bool`

HasVrfName returns a boolean if a field has been set.

### GetVrfNameAutoAssigned

`func (o *ConfigPutRequestTenantTenantName) GetVrfNameAutoAssigned() bool`

GetVrfNameAutoAssigned returns the VrfNameAutoAssigned field if non-nil, zero value otherwise.

### GetVrfNameAutoAssignedOk

`func (o *ConfigPutRequestTenantTenantName) GetVrfNameAutoAssignedOk() (*bool, bool)`

GetVrfNameAutoAssignedOk returns a tuple with the VrfNameAutoAssigned field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVrfNameAutoAssigned

`func (o *ConfigPutRequestTenantTenantName) SetVrfNameAutoAssigned(v bool)`

SetVrfNameAutoAssigned sets VrfNameAutoAssigned field to given value.

### HasVrfNameAutoAssigned

`func (o *ConfigPutRequestTenantTenantName) HasVrfNameAutoAssigned() bool`

HasVrfNameAutoAssigned returns a boolean if a field has been set.

### GetRouteTenants

`func (o *ConfigPutRequestTenantTenantName) GetRouteTenants() []ConfigPutRequestTenantTenantNameRouteTenantsInner`

GetRouteTenants returns the RouteTenants field if non-nil, zero value otherwise.

### GetRouteTenantsOk

`func (o *ConfigPutRequestTenantTenantName) GetRouteTenantsOk() (*[]ConfigPutRequestTenantTenantNameRouteTenantsInner, bool)`

GetRouteTenantsOk returns a tuple with the RouteTenants field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRouteTenants

`func (o *ConfigPutRequestTenantTenantName) SetRouteTenants(v []ConfigPutRequestTenantTenantNameRouteTenantsInner)`

SetRouteTenants sets RouteTenants field to given value.

### HasRouteTenants

`func (o *ConfigPutRequestTenantTenantName) HasRouteTenants() bool`

HasRouteTenants returns a boolean if a field has been set.

### GetObjectProperties

`func (o *ConfigPutRequestTenantTenantName) GetObjectProperties() ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ConfigPutRequestTenantTenantName) GetObjectPropertiesOk() (*ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ConfigPutRequestTenantTenantName) SetObjectProperties(v ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ConfigPutRequestTenantTenantName) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


