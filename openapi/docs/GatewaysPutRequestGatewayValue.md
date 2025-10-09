# GatewaysPutRequestGatewayValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. It&#39;s highly recommended to set this value to true so that validation on the object will be ran. | [optional] [default to false]
**Tenant** | Pointer to **string** | Tenant | [optional] [default to ""]
**TenantRefType** | Pointer to **string** | Object type for tenant field | [optional] 
**NeighborIpAddress** | Pointer to **string** | IP address of remote BGP peer | [optional] [default to ""]
**NeighborAsNumber** | Pointer to **NullableInt32** | Autonomous System Number of remote BGP peer  | [optional] 
**FabricInterconnect** | Pointer to **bool** |  | [optional] [default to false]
**KeepaliveTimer** | Pointer to **int32** | Interval in seconds between Keepalive messages sent to remote BGP peer | [optional] [default to 60]
**HoldTimer** | Pointer to **int32** | Time, in seconds,  used to determine failure of session Keepalive messages received from remote BGP peer  | [optional] [default to 180]
**ConnectTimer** | Pointer to **int32** | Time in seconds between sucessive attempts to Establish BGP session | [optional] [default to 120]
**AdvertisementInterval** | Pointer to **int32** | The minimum time in seconds between sending route updates to BGP neighbor  | [optional] [default to 30]
**EbgpMultihop** | Pointer to **int32** | Allows external BGP neighbors to establish peering session multiple network hops away.  | [optional] [default to 255]
**EgressVlan** | Pointer to **NullableInt32** | VLAN used to carry BGP TCP session | [optional] 
**SourceIpAddress** | Pointer to **string** | Source IP address used to override the default source address calculation for BGP TCP session | [optional] [default to ""]
**AnycastIpMask** | Pointer to **string** | The Anycast Address will be used to enable an IP routing redundancy mechanism designed to allow for transparent failover across a leaf pair at the first-hop IP router. | [optional] [default to ""]
**Md5Password** | Pointer to **string** | MD5 password | [optional] [default to ""]
**ImportRouteMap** | Pointer to **string** | A route-map applied to routes imported into the current tenant from the targeted BGP router with the purpose of filtering or modifying the routes | [optional] [default to ""]
**ImportRouteMapRefType** | Pointer to **string** | Object type for import_route_map field | [optional] 
**ExportRouteMap** | Pointer to **string** | A route-map applied to routes exported into the current tenant from the targeted BGP router with the purpose of filtering or modifying the routes | [optional] [default to ""]
**ExportRouteMapRefType** | Pointer to **string** | Object type for export_route_map field | [optional] 
**GatewayMode** | Pointer to **string** | Gateway Mode. Can be BGP, Static, or Default | [optional] [default to "Static BGP"]
**LocalAsNumber** | Pointer to **NullableInt32** | Local AS Number | [optional] 
**LocalAsNoPrepend** | Pointer to **bool** | Do not prepend the local-as number to the AS-PATH for routes advertised through this BGP gateway. The Local AS Number must be set for this to be able to be set. | [optional] [default to false]
**ReplaceAs** | Pointer to **bool** | Prepend only Local AS in updates to EBGP peers. | [optional] [default to false]
**MaxLocalAsOccurrences** | Pointer to **NullableInt32** | Allow routes with the local AS number in the AS-path, specifying the maximum occurrences permitted before declaring a routing loop. Leave blank or &#39;0&#39; to disable. | [optional] [default to 0]
**DynamicBgpSubnet** | Pointer to **string** | Dynamic BGP Subnet | [optional] [default to ""]
**DynamicBgpLimits** | Pointer to **NullableInt32** | Dynamic BGP Limits | [optional] [default to 0]
**HelperHopIpAddress** | Pointer to **string** | Helper Hop IP Address | [optional] [default to ""]
**EnableBfd** | Pointer to **bool** | Enable BFD(Bi-Directional Forwarding) | [optional] [default to false]
**BfdReceiveInterval** | Pointer to **NullableInt32** | Configure the minimum interval during which the system can receive BFD control packets | [optional] [default to 300]
**BfdTransmissionInterval** | Pointer to **NullableInt32** | Configure the minimum transmission interval during which the system can send BFD control packets | [optional] [default to 300]
**BfdDetectMultiplier** | Pointer to **NullableInt32** | Configure the detection multiplier to determine packet loss | [optional] [default to 3]
**NextHopSelf** | Pointer to **bool** | Optional attribute that disables the normal BGP calculation of next-hops for advertised routes and instead sets the next-hops for advertised routes to the IP address of the switch itself. | [optional] [default to false]
**StaticRoutes** | Pointer to [**[]GatewaysPutRequestGatewayValueStaticRoutesInner**](GatewaysPutRequestGatewayValueStaticRoutesInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties**](DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties.md) |  | [optional] 
**DefaultOriginate** | Pointer to **bool** | Instructs BGP to generate and send a default route 0.0.0.0/0 to the specified neighbor. | [optional] [default to false]
**BfdMultihop** | Pointer to **bool** | Enable BFD Multi-Hop for Neighbor. This is used to detect failures in the forwarding path between the BGP peers. | [optional] [default to false]

## Methods

### NewGatewaysPutRequestGatewayValue

`func NewGatewaysPutRequestGatewayValue() *GatewaysPutRequestGatewayValue`

NewGatewaysPutRequestGatewayValue instantiates a new GatewaysPutRequestGatewayValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGatewaysPutRequestGatewayValueWithDefaults

`func NewGatewaysPutRequestGatewayValueWithDefaults() *GatewaysPutRequestGatewayValue`

NewGatewaysPutRequestGatewayValueWithDefaults instantiates a new GatewaysPutRequestGatewayValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *GatewaysPutRequestGatewayValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *GatewaysPutRequestGatewayValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *GatewaysPutRequestGatewayValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *GatewaysPutRequestGatewayValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *GatewaysPutRequestGatewayValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *GatewaysPutRequestGatewayValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *GatewaysPutRequestGatewayValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *GatewaysPutRequestGatewayValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetTenant

`func (o *GatewaysPutRequestGatewayValue) GetTenant() string`

GetTenant returns the Tenant field if non-nil, zero value otherwise.

### GetTenantOk

`func (o *GatewaysPutRequestGatewayValue) GetTenantOk() (*string, bool)`

GetTenantOk returns a tuple with the Tenant field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenant

`func (o *GatewaysPutRequestGatewayValue) SetTenant(v string)`

SetTenant sets Tenant field to given value.

### HasTenant

`func (o *GatewaysPutRequestGatewayValue) HasTenant() bool`

HasTenant returns a boolean if a field has been set.

### GetTenantRefType

`func (o *GatewaysPutRequestGatewayValue) GetTenantRefType() string`

GetTenantRefType returns the TenantRefType field if non-nil, zero value otherwise.

### GetTenantRefTypeOk

`func (o *GatewaysPutRequestGatewayValue) GetTenantRefTypeOk() (*string, bool)`

GetTenantRefTypeOk returns a tuple with the TenantRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantRefType

`func (o *GatewaysPutRequestGatewayValue) SetTenantRefType(v string)`

SetTenantRefType sets TenantRefType field to given value.

### HasTenantRefType

`func (o *GatewaysPutRequestGatewayValue) HasTenantRefType() bool`

HasTenantRefType returns a boolean if a field has been set.

### GetNeighborIpAddress

`func (o *GatewaysPutRequestGatewayValue) GetNeighborIpAddress() string`

GetNeighborIpAddress returns the NeighborIpAddress field if non-nil, zero value otherwise.

### GetNeighborIpAddressOk

`func (o *GatewaysPutRequestGatewayValue) GetNeighborIpAddressOk() (*string, bool)`

GetNeighborIpAddressOk returns a tuple with the NeighborIpAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNeighborIpAddress

`func (o *GatewaysPutRequestGatewayValue) SetNeighborIpAddress(v string)`

SetNeighborIpAddress sets NeighborIpAddress field to given value.

### HasNeighborIpAddress

`func (o *GatewaysPutRequestGatewayValue) HasNeighborIpAddress() bool`

HasNeighborIpAddress returns a boolean if a field has been set.

### GetNeighborAsNumber

`func (o *GatewaysPutRequestGatewayValue) GetNeighborAsNumber() int32`

GetNeighborAsNumber returns the NeighborAsNumber field if non-nil, zero value otherwise.

### GetNeighborAsNumberOk

`func (o *GatewaysPutRequestGatewayValue) GetNeighborAsNumberOk() (*int32, bool)`

GetNeighborAsNumberOk returns a tuple with the NeighborAsNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNeighborAsNumber

`func (o *GatewaysPutRequestGatewayValue) SetNeighborAsNumber(v int32)`

SetNeighborAsNumber sets NeighborAsNumber field to given value.

### HasNeighborAsNumber

`func (o *GatewaysPutRequestGatewayValue) HasNeighborAsNumber() bool`

HasNeighborAsNumber returns a boolean if a field has been set.

### SetNeighborAsNumberNil

`func (o *GatewaysPutRequestGatewayValue) SetNeighborAsNumberNil(b bool)`

 SetNeighborAsNumberNil sets the value for NeighborAsNumber to be an explicit nil

### UnsetNeighborAsNumber
`func (o *GatewaysPutRequestGatewayValue) UnsetNeighborAsNumber()`

UnsetNeighborAsNumber ensures that no value is present for NeighborAsNumber, not even an explicit nil
### GetFabricInterconnect

`func (o *GatewaysPutRequestGatewayValue) GetFabricInterconnect() bool`

GetFabricInterconnect returns the FabricInterconnect field if non-nil, zero value otherwise.

### GetFabricInterconnectOk

`func (o *GatewaysPutRequestGatewayValue) GetFabricInterconnectOk() (*bool, bool)`

GetFabricInterconnectOk returns a tuple with the FabricInterconnect field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFabricInterconnect

`func (o *GatewaysPutRequestGatewayValue) SetFabricInterconnect(v bool)`

SetFabricInterconnect sets FabricInterconnect field to given value.

### HasFabricInterconnect

`func (o *GatewaysPutRequestGatewayValue) HasFabricInterconnect() bool`

HasFabricInterconnect returns a boolean if a field has been set.

### GetKeepaliveTimer

`func (o *GatewaysPutRequestGatewayValue) GetKeepaliveTimer() int32`

GetKeepaliveTimer returns the KeepaliveTimer field if non-nil, zero value otherwise.

### GetKeepaliveTimerOk

`func (o *GatewaysPutRequestGatewayValue) GetKeepaliveTimerOk() (*int32, bool)`

GetKeepaliveTimerOk returns a tuple with the KeepaliveTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeepaliveTimer

`func (o *GatewaysPutRequestGatewayValue) SetKeepaliveTimer(v int32)`

SetKeepaliveTimer sets KeepaliveTimer field to given value.

### HasKeepaliveTimer

`func (o *GatewaysPutRequestGatewayValue) HasKeepaliveTimer() bool`

HasKeepaliveTimer returns a boolean if a field has been set.

### GetHoldTimer

`func (o *GatewaysPutRequestGatewayValue) GetHoldTimer() int32`

GetHoldTimer returns the HoldTimer field if non-nil, zero value otherwise.

### GetHoldTimerOk

`func (o *GatewaysPutRequestGatewayValue) GetHoldTimerOk() (*int32, bool)`

GetHoldTimerOk returns a tuple with the HoldTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHoldTimer

`func (o *GatewaysPutRequestGatewayValue) SetHoldTimer(v int32)`

SetHoldTimer sets HoldTimer field to given value.

### HasHoldTimer

`func (o *GatewaysPutRequestGatewayValue) HasHoldTimer() bool`

HasHoldTimer returns a boolean if a field has been set.

### GetConnectTimer

`func (o *GatewaysPutRequestGatewayValue) GetConnectTimer() int32`

GetConnectTimer returns the ConnectTimer field if non-nil, zero value otherwise.

### GetConnectTimerOk

`func (o *GatewaysPutRequestGatewayValue) GetConnectTimerOk() (*int32, bool)`

GetConnectTimerOk returns a tuple with the ConnectTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnectTimer

`func (o *GatewaysPutRequestGatewayValue) SetConnectTimer(v int32)`

SetConnectTimer sets ConnectTimer field to given value.

### HasConnectTimer

`func (o *GatewaysPutRequestGatewayValue) HasConnectTimer() bool`

HasConnectTimer returns a boolean if a field has been set.

### GetAdvertisementInterval

`func (o *GatewaysPutRequestGatewayValue) GetAdvertisementInterval() int32`

GetAdvertisementInterval returns the AdvertisementInterval field if non-nil, zero value otherwise.

### GetAdvertisementIntervalOk

`func (o *GatewaysPutRequestGatewayValue) GetAdvertisementIntervalOk() (*int32, bool)`

GetAdvertisementIntervalOk returns a tuple with the AdvertisementInterval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAdvertisementInterval

`func (o *GatewaysPutRequestGatewayValue) SetAdvertisementInterval(v int32)`

SetAdvertisementInterval sets AdvertisementInterval field to given value.

### HasAdvertisementInterval

`func (o *GatewaysPutRequestGatewayValue) HasAdvertisementInterval() bool`

HasAdvertisementInterval returns a boolean if a field has been set.

### GetEbgpMultihop

`func (o *GatewaysPutRequestGatewayValue) GetEbgpMultihop() int32`

GetEbgpMultihop returns the EbgpMultihop field if non-nil, zero value otherwise.

### GetEbgpMultihopOk

`func (o *GatewaysPutRequestGatewayValue) GetEbgpMultihopOk() (*int32, bool)`

GetEbgpMultihopOk returns a tuple with the EbgpMultihop field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEbgpMultihop

`func (o *GatewaysPutRequestGatewayValue) SetEbgpMultihop(v int32)`

SetEbgpMultihop sets EbgpMultihop field to given value.

### HasEbgpMultihop

`func (o *GatewaysPutRequestGatewayValue) HasEbgpMultihop() bool`

HasEbgpMultihop returns a boolean if a field has been set.

### GetEgressVlan

`func (o *GatewaysPutRequestGatewayValue) GetEgressVlan() int32`

GetEgressVlan returns the EgressVlan field if non-nil, zero value otherwise.

### GetEgressVlanOk

`func (o *GatewaysPutRequestGatewayValue) GetEgressVlanOk() (*int32, bool)`

GetEgressVlanOk returns a tuple with the EgressVlan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEgressVlan

`func (o *GatewaysPutRequestGatewayValue) SetEgressVlan(v int32)`

SetEgressVlan sets EgressVlan field to given value.

### HasEgressVlan

`func (o *GatewaysPutRequestGatewayValue) HasEgressVlan() bool`

HasEgressVlan returns a boolean if a field has been set.

### SetEgressVlanNil

`func (o *GatewaysPutRequestGatewayValue) SetEgressVlanNil(b bool)`

 SetEgressVlanNil sets the value for EgressVlan to be an explicit nil

### UnsetEgressVlan
`func (o *GatewaysPutRequestGatewayValue) UnsetEgressVlan()`

UnsetEgressVlan ensures that no value is present for EgressVlan, not even an explicit nil
### GetSourceIpAddress

`func (o *GatewaysPutRequestGatewayValue) GetSourceIpAddress() string`

GetSourceIpAddress returns the SourceIpAddress field if non-nil, zero value otherwise.

### GetSourceIpAddressOk

`func (o *GatewaysPutRequestGatewayValue) GetSourceIpAddressOk() (*string, bool)`

GetSourceIpAddressOk returns a tuple with the SourceIpAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSourceIpAddress

`func (o *GatewaysPutRequestGatewayValue) SetSourceIpAddress(v string)`

SetSourceIpAddress sets SourceIpAddress field to given value.

### HasSourceIpAddress

`func (o *GatewaysPutRequestGatewayValue) HasSourceIpAddress() bool`

HasSourceIpAddress returns a boolean if a field has been set.

### GetAnycastIpMask

`func (o *GatewaysPutRequestGatewayValue) GetAnycastIpMask() string`

GetAnycastIpMask returns the AnycastIpMask field if non-nil, zero value otherwise.

### GetAnycastIpMaskOk

`func (o *GatewaysPutRequestGatewayValue) GetAnycastIpMaskOk() (*string, bool)`

GetAnycastIpMaskOk returns a tuple with the AnycastIpMask field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnycastIpMask

`func (o *GatewaysPutRequestGatewayValue) SetAnycastIpMask(v string)`

SetAnycastIpMask sets AnycastIpMask field to given value.

### HasAnycastIpMask

`func (o *GatewaysPutRequestGatewayValue) HasAnycastIpMask() bool`

HasAnycastIpMask returns a boolean if a field has been set.

### GetMd5Password

`func (o *GatewaysPutRequestGatewayValue) GetMd5Password() string`

GetMd5Password returns the Md5Password field if non-nil, zero value otherwise.

### GetMd5PasswordOk

`func (o *GatewaysPutRequestGatewayValue) GetMd5PasswordOk() (*string, bool)`

GetMd5PasswordOk returns a tuple with the Md5Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMd5Password

`func (o *GatewaysPutRequestGatewayValue) SetMd5Password(v string)`

SetMd5Password sets Md5Password field to given value.

### HasMd5Password

`func (o *GatewaysPutRequestGatewayValue) HasMd5Password() bool`

HasMd5Password returns a boolean if a field has been set.

### GetImportRouteMap

`func (o *GatewaysPutRequestGatewayValue) GetImportRouteMap() string`

GetImportRouteMap returns the ImportRouteMap field if non-nil, zero value otherwise.

### GetImportRouteMapOk

`func (o *GatewaysPutRequestGatewayValue) GetImportRouteMapOk() (*string, bool)`

GetImportRouteMapOk returns a tuple with the ImportRouteMap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImportRouteMap

`func (o *GatewaysPutRequestGatewayValue) SetImportRouteMap(v string)`

SetImportRouteMap sets ImportRouteMap field to given value.

### HasImportRouteMap

`func (o *GatewaysPutRequestGatewayValue) HasImportRouteMap() bool`

HasImportRouteMap returns a boolean if a field has been set.

### GetImportRouteMapRefType

`func (o *GatewaysPutRequestGatewayValue) GetImportRouteMapRefType() string`

GetImportRouteMapRefType returns the ImportRouteMapRefType field if non-nil, zero value otherwise.

### GetImportRouteMapRefTypeOk

`func (o *GatewaysPutRequestGatewayValue) GetImportRouteMapRefTypeOk() (*string, bool)`

GetImportRouteMapRefTypeOk returns a tuple with the ImportRouteMapRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImportRouteMapRefType

`func (o *GatewaysPutRequestGatewayValue) SetImportRouteMapRefType(v string)`

SetImportRouteMapRefType sets ImportRouteMapRefType field to given value.

### HasImportRouteMapRefType

`func (o *GatewaysPutRequestGatewayValue) HasImportRouteMapRefType() bool`

HasImportRouteMapRefType returns a boolean if a field has been set.

### GetExportRouteMap

`func (o *GatewaysPutRequestGatewayValue) GetExportRouteMap() string`

GetExportRouteMap returns the ExportRouteMap field if non-nil, zero value otherwise.

### GetExportRouteMapOk

`func (o *GatewaysPutRequestGatewayValue) GetExportRouteMapOk() (*string, bool)`

GetExportRouteMapOk returns a tuple with the ExportRouteMap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExportRouteMap

`func (o *GatewaysPutRequestGatewayValue) SetExportRouteMap(v string)`

SetExportRouteMap sets ExportRouteMap field to given value.

### HasExportRouteMap

`func (o *GatewaysPutRequestGatewayValue) HasExportRouteMap() bool`

HasExportRouteMap returns a boolean if a field has been set.

### GetExportRouteMapRefType

`func (o *GatewaysPutRequestGatewayValue) GetExportRouteMapRefType() string`

GetExportRouteMapRefType returns the ExportRouteMapRefType field if non-nil, zero value otherwise.

### GetExportRouteMapRefTypeOk

`func (o *GatewaysPutRequestGatewayValue) GetExportRouteMapRefTypeOk() (*string, bool)`

GetExportRouteMapRefTypeOk returns a tuple with the ExportRouteMapRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExportRouteMapRefType

`func (o *GatewaysPutRequestGatewayValue) SetExportRouteMapRefType(v string)`

SetExportRouteMapRefType sets ExportRouteMapRefType field to given value.

### HasExportRouteMapRefType

`func (o *GatewaysPutRequestGatewayValue) HasExportRouteMapRefType() bool`

HasExportRouteMapRefType returns a boolean if a field has been set.

### GetGatewayMode

`func (o *GatewaysPutRequestGatewayValue) GetGatewayMode() string`

GetGatewayMode returns the GatewayMode field if non-nil, zero value otherwise.

### GetGatewayModeOk

`func (o *GatewaysPutRequestGatewayValue) GetGatewayModeOk() (*string, bool)`

GetGatewayModeOk returns a tuple with the GatewayMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGatewayMode

`func (o *GatewaysPutRequestGatewayValue) SetGatewayMode(v string)`

SetGatewayMode sets GatewayMode field to given value.

### HasGatewayMode

`func (o *GatewaysPutRequestGatewayValue) HasGatewayMode() bool`

HasGatewayMode returns a boolean if a field has been set.

### GetLocalAsNumber

`func (o *GatewaysPutRequestGatewayValue) GetLocalAsNumber() int32`

GetLocalAsNumber returns the LocalAsNumber field if non-nil, zero value otherwise.

### GetLocalAsNumberOk

`func (o *GatewaysPutRequestGatewayValue) GetLocalAsNumberOk() (*int32, bool)`

GetLocalAsNumberOk returns a tuple with the LocalAsNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocalAsNumber

`func (o *GatewaysPutRequestGatewayValue) SetLocalAsNumber(v int32)`

SetLocalAsNumber sets LocalAsNumber field to given value.

### HasLocalAsNumber

`func (o *GatewaysPutRequestGatewayValue) HasLocalAsNumber() bool`

HasLocalAsNumber returns a boolean if a field has been set.

### SetLocalAsNumberNil

`func (o *GatewaysPutRequestGatewayValue) SetLocalAsNumberNil(b bool)`

 SetLocalAsNumberNil sets the value for LocalAsNumber to be an explicit nil

### UnsetLocalAsNumber
`func (o *GatewaysPutRequestGatewayValue) UnsetLocalAsNumber()`

UnsetLocalAsNumber ensures that no value is present for LocalAsNumber, not even an explicit nil
### GetLocalAsNoPrepend

`func (o *GatewaysPutRequestGatewayValue) GetLocalAsNoPrepend() bool`

GetLocalAsNoPrepend returns the LocalAsNoPrepend field if non-nil, zero value otherwise.

### GetLocalAsNoPrependOk

`func (o *GatewaysPutRequestGatewayValue) GetLocalAsNoPrependOk() (*bool, bool)`

GetLocalAsNoPrependOk returns a tuple with the LocalAsNoPrepend field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocalAsNoPrepend

`func (o *GatewaysPutRequestGatewayValue) SetLocalAsNoPrepend(v bool)`

SetLocalAsNoPrepend sets LocalAsNoPrepend field to given value.

### HasLocalAsNoPrepend

`func (o *GatewaysPutRequestGatewayValue) HasLocalAsNoPrepend() bool`

HasLocalAsNoPrepend returns a boolean if a field has been set.

### GetReplaceAs

`func (o *GatewaysPutRequestGatewayValue) GetReplaceAs() bool`

GetReplaceAs returns the ReplaceAs field if non-nil, zero value otherwise.

### GetReplaceAsOk

`func (o *GatewaysPutRequestGatewayValue) GetReplaceAsOk() (*bool, bool)`

GetReplaceAsOk returns a tuple with the ReplaceAs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReplaceAs

`func (o *GatewaysPutRequestGatewayValue) SetReplaceAs(v bool)`

SetReplaceAs sets ReplaceAs field to given value.

### HasReplaceAs

`func (o *GatewaysPutRequestGatewayValue) HasReplaceAs() bool`

HasReplaceAs returns a boolean if a field has been set.

### GetMaxLocalAsOccurrences

`func (o *GatewaysPutRequestGatewayValue) GetMaxLocalAsOccurrences() int32`

GetMaxLocalAsOccurrences returns the MaxLocalAsOccurrences field if non-nil, zero value otherwise.

### GetMaxLocalAsOccurrencesOk

`func (o *GatewaysPutRequestGatewayValue) GetMaxLocalAsOccurrencesOk() (*int32, bool)`

GetMaxLocalAsOccurrencesOk returns a tuple with the MaxLocalAsOccurrences field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxLocalAsOccurrences

`func (o *GatewaysPutRequestGatewayValue) SetMaxLocalAsOccurrences(v int32)`

SetMaxLocalAsOccurrences sets MaxLocalAsOccurrences field to given value.

### HasMaxLocalAsOccurrences

`func (o *GatewaysPutRequestGatewayValue) HasMaxLocalAsOccurrences() bool`

HasMaxLocalAsOccurrences returns a boolean if a field has been set.

### SetMaxLocalAsOccurrencesNil

`func (o *GatewaysPutRequestGatewayValue) SetMaxLocalAsOccurrencesNil(b bool)`

 SetMaxLocalAsOccurrencesNil sets the value for MaxLocalAsOccurrences to be an explicit nil

### UnsetMaxLocalAsOccurrences
`func (o *GatewaysPutRequestGatewayValue) UnsetMaxLocalAsOccurrences()`

UnsetMaxLocalAsOccurrences ensures that no value is present for MaxLocalAsOccurrences, not even an explicit nil
### GetDynamicBgpSubnet

`func (o *GatewaysPutRequestGatewayValue) GetDynamicBgpSubnet() string`

GetDynamicBgpSubnet returns the DynamicBgpSubnet field if non-nil, zero value otherwise.

### GetDynamicBgpSubnetOk

`func (o *GatewaysPutRequestGatewayValue) GetDynamicBgpSubnetOk() (*string, bool)`

GetDynamicBgpSubnetOk returns a tuple with the DynamicBgpSubnet field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDynamicBgpSubnet

`func (o *GatewaysPutRequestGatewayValue) SetDynamicBgpSubnet(v string)`

SetDynamicBgpSubnet sets DynamicBgpSubnet field to given value.

### HasDynamicBgpSubnet

`func (o *GatewaysPutRequestGatewayValue) HasDynamicBgpSubnet() bool`

HasDynamicBgpSubnet returns a boolean if a field has been set.

### GetDynamicBgpLimits

`func (o *GatewaysPutRequestGatewayValue) GetDynamicBgpLimits() int32`

GetDynamicBgpLimits returns the DynamicBgpLimits field if non-nil, zero value otherwise.

### GetDynamicBgpLimitsOk

`func (o *GatewaysPutRequestGatewayValue) GetDynamicBgpLimitsOk() (*int32, bool)`

GetDynamicBgpLimitsOk returns a tuple with the DynamicBgpLimits field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDynamicBgpLimits

`func (o *GatewaysPutRequestGatewayValue) SetDynamicBgpLimits(v int32)`

SetDynamicBgpLimits sets DynamicBgpLimits field to given value.

### HasDynamicBgpLimits

`func (o *GatewaysPutRequestGatewayValue) HasDynamicBgpLimits() bool`

HasDynamicBgpLimits returns a boolean if a field has been set.

### SetDynamicBgpLimitsNil

`func (o *GatewaysPutRequestGatewayValue) SetDynamicBgpLimitsNil(b bool)`

 SetDynamicBgpLimitsNil sets the value for DynamicBgpLimits to be an explicit nil

### UnsetDynamicBgpLimits
`func (o *GatewaysPutRequestGatewayValue) UnsetDynamicBgpLimits()`

UnsetDynamicBgpLimits ensures that no value is present for DynamicBgpLimits, not even an explicit nil
### GetHelperHopIpAddress

`func (o *GatewaysPutRequestGatewayValue) GetHelperHopIpAddress() string`

GetHelperHopIpAddress returns the HelperHopIpAddress field if non-nil, zero value otherwise.

### GetHelperHopIpAddressOk

`func (o *GatewaysPutRequestGatewayValue) GetHelperHopIpAddressOk() (*string, bool)`

GetHelperHopIpAddressOk returns a tuple with the HelperHopIpAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHelperHopIpAddress

`func (o *GatewaysPutRequestGatewayValue) SetHelperHopIpAddress(v string)`

SetHelperHopIpAddress sets HelperHopIpAddress field to given value.

### HasHelperHopIpAddress

`func (o *GatewaysPutRequestGatewayValue) HasHelperHopIpAddress() bool`

HasHelperHopIpAddress returns a boolean if a field has been set.

### GetEnableBfd

`func (o *GatewaysPutRequestGatewayValue) GetEnableBfd() bool`

GetEnableBfd returns the EnableBfd field if non-nil, zero value otherwise.

### GetEnableBfdOk

`func (o *GatewaysPutRequestGatewayValue) GetEnableBfdOk() (*bool, bool)`

GetEnableBfdOk returns a tuple with the EnableBfd field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnableBfd

`func (o *GatewaysPutRequestGatewayValue) SetEnableBfd(v bool)`

SetEnableBfd sets EnableBfd field to given value.

### HasEnableBfd

`func (o *GatewaysPutRequestGatewayValue) HasEnableBfd() bool`

HasEnableBfd returns a boolean if a field has been set.

### GetBfdReceiveInterval

`func (o *GatewaysPutRequestGatewayValue) GetBfdReceiveInterval() int32`

GetBfdReceiveInterval returns the BfdReceiveInterval field if non-nil, zero value otherwise.

### GetBfdReceiveIntervalOk

`func (o *GatewaysPutRequestGatewayValue) GetBfdReceiveIntervalOk() (*int32, bool)`

GetBfdReceiveIntervalOk returns a tuple with the BfdReceiveInterval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBfdReceiveInterval

`func (o *GatewaysPutRequestGatewayValue) SetBfdReceiveInterval(v int32)`

SetBfdReceiveInterval sets BfdReceiveInterval field to given value.

### HasBfdReceiveInterval

`func (o *GatewaysPutRequestGatewayValue) HasBfdReceiveInterval() bool`

HasBfdReceiveInterval returns a boolean if a field has been set.

### SetBfdReceiveIntervalNil

`func (o *GatewaysPutRequestGatewayValue) SetBfdReceiveIntervalNil(b bool)`

 SetBfdReceiveIntervalNil sets the value for BfdReceiveInterval to be an explicit nil

### UnsetBfdReceiveInterval
`func (o *GatewaysPutRequestGatewayValue) UnsetBfdReceiveInterval()`

UnsetBfdReceiveInterval ensures that no value is present for BfdReceiveInterval, not even an explicit nil
### GetBfdTransmissionInterval

`func (o *GatewaysPutRequestGatewayValue) GetBfdTransmissionInterval() int32`

GetBfdTransmissionInterval returns the BfdTransmissionInterval field if non-nil, zero value otherwise.

### GetBfdTransmissionIntervalOk

`func (o *GatewaysPutRequestGatewayValue) GetBfdTransmissionIntervalOk() (*int32, bool)`

GetBfdTransmissionIntervalOk returns a tuple with the BfdTransmissionInterval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBfdTransmissionInterval

`func (o *GatewaysPutRequestGatewayValue) SetBfdTransmissionInterval(v int32)`

SetBfdTransmissionInterval sets BfdTransmissionInterval field to given value.

### HasBfdTransmissionInterval

`func (o *GatewaysPutRequestGatewayValue) HasBfdTransmissionInterval() bool`

HasBfdTransmissionInterval returns a boolean if a field has been set.

### SetBfdTransmissionIntervalNil

`func (o *GatewaysPutRequestGatewayValue) SetBfdTransmissionIntervalNil(b bool)`

 SetBfdTransmissionIntervalNil sets the value for BfdTransmissionInterval to be an explicit nil

### UnsetBfdTransmissionInterval
`func (o *GatewaysPutRequestGatewayValue) UnsetBfdTransmissionInterval()`

UnsetBfdTransmissionInterval ensures that no value is present for BfdTransmissionInterval, not even an explicit nil
### GetBfdDetectMultiplier

`func (o *GatewaysPutRequestGatewayValue) GetBfdDetectMultiplier() int32`

GetBfdDetectMultiplier returns the BfdDetectMultiplier field if non-nil, zero value otherwise.

### GetBfdDetectMultiplierOk

`func (o *GatewaysPutRequestGatewayValue) GetBfdDetectMultiplierOk() (*int32, bool)`

GetBfdDetectMultiplierOk returns a tuple with the BfdDetectMultiplier field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBfdDetectMultiplier

`func (o *GatewaysPutRequestGatewayValue) SetBfdDetectMultiplier(v int32)`

SetBfdDetectMultiplier sets BfdDetectMultiplier field to given value.

### HasBfdDetectMultiplier

`func (o *GatewaysPutRequestGatewayValue) HasBfdDetectMultiplier() bool`

HasBfdDetectMultiplier returns a boolean if a field has been set.

### SetBfdDetectMultiplierNil

`func (o *GatewaysPutRequestGatewayValue) SetBfdDetectMultiplierNil(b bool)`

 SetBfdDetectMultiplierNil sets the value for BfdDetectMultiplier to be an explicit nil

### UnsetBfdDetectMultiplier
`func (o *GatewaysPutRequestGatewayValue) UnsetBfdDetectMultiplier()`

UnsetBfdDetectMultiplier ensures that no value is present for BfdDetectMultiplier, not even an explicit nil
### GetNextHopSelf

`func (o *GatewaysPutRequestGatewayValue) GetNextHopSelf() bool`

GetNextHopSelf returns the NextHopSelf field if non-nil, zero value otherwise.

### GetNextHopSelfOk

`func (o *GatewaysPutRequestGatewayValue) GetNextHopSelfOk() (*bool, bool)`

GetNextHopSelfOk returns a tuple with the NextHopSelf field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNextHopSelf

`func (o *GatewaysPutRequestGatewayValue) SetNextHopSelf(v bool)`

SetNextHopSelf sets NextHopSelf field to given value.

### HasNextHopSelf

`func (o *GatewaysPutRequestGatewayValue) HasNextHopSelf() bool`

HasNextHopSelf returns a boolean if a field has been set.

### GetStaticRoutes

`func (o *GatewaysPutRequestGatewayValue) GetStaticRoutes() []GatewaysPutRequestGatewayValueStaticRoutesInner`

GetStaticRoutes returns the StaticRoutes field if non-nil, zero value otherwise.

### GetStaticRoutesOk

`func (o *GatewaysPutRequestGatewayValue) GetStaticRoutesOk() (*[]GatewaysPutRequestGatewayValueStaticRoutesInner, bool)`

GetStaticRoutesOk returns a tuple with the StaticRoutes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStaticRoutes

`func (o *GatewaysPutRequestGatewayValue) SetStaticRoutes(v []GatewaysPutRequestGatewayValueStaticRoutesInner)`

SetStaticRoutes sets StaticRoutes field to given value.

### HasStaticRoutes

`func (o *GatewaysPutRequestGatewayValue) HasStaticRoutes() bool`

HasStaticRoutes returns a boolean if a field has been set.

### GetObjectProperties

`func (o *GatewaysPutRequestGatewayValue) GetObjectProperties() DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *GatewaysPutRequestGatewayValue) GetObjectPropertiesOk() (*DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *GatewaysPutRequestGatewayValue) SetObjectProperties(v DevicesettingsPutRequestEthDeviceProfilesValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *GatewaysPutRequestGatewayValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.

### GetDefaultOriginate

`func (o *GatewaysPutRequestGatewayValue) GetDefaultOriginate() bool`

GetDefaultOriginate returns the DefaultOriginate field if non-nil, zero value otherwise.

### GetDefaultOriginateOk

`func (o *GatewaysPutRequestGatewayValue) GetDefaultOriginateOk() (*bool, bool)`

GetDefaultOriginateOk returns a tuple with the DefaultOriginate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefaultOriginate

`func (o *GatewaysPutRequestGatewayValue) SetDefaultOriginate(v bool)`

SetDefaultOriginate sets DefaultOriginate field to given value.

### HasDefaultOriginate

`func (o *GatewaysPutRequestGatewayValue) HasDefaultOriginate() bool`

HasDefaultOriginate returns a boolean if a field has been set.

### GetBfdMultihop

`func (o *GatewaysPutRequestGatewayValue) GetBfdMultihop() bool`

GetBfdMultihop returns the BfdMultihop field if non-nil, zero value otherwise.

### GetBfdMultihopOk

`func (o *GatewaysPutRequestGatewayValue) GetBfdMultihopOk() (*bool, bool)`

GetBfdMultihopOk returns a tuple with the BfdMultihop field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBfdMultihop

`func (o *GatewaysPutRequestGatewayValue) SetBfdMultihop(v bool)`

SetBfdMultihop sets BfdMultihop field to given value.

### HasBfdMultihop

`func (o *GatewaysPutRequestGatewayValue) HasBfdMultihop() bool`

HasBfdMultihop returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


