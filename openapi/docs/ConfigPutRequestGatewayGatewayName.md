# ConfigPutRequestGatewayGatewayName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**Tenant** | Pointer to **string** | Tenant | [optional] [default to ""]
**TenantRefType** | Pointer to **string** | Object type for tenant field | [optional] 
**NeighborIpAddress** | Pointer to **string** | IP address of remote BGP peer | [optional] [default to ""]
**NeighborAsNumber** | Pointer to **NullableInt32** | Autonomous System Number of remote BGP peer  | [optional] 
**FabricInterconnect** | Pointer to **bool** |  | [optional] [default to false]
**KeepaliveTimer** | Pointer to **int32** | Interval in seconds between Keepalive messages sent to remote BGP peer | [optional] 
**HoldTimer** | Pointer to **int32** | Time, in seconds,  used to determine failure of session Keepalive messages received from remote BGP peer  | [optional] 
**ConnectTimer** | Pointer to **int32** | Time in seconds between sucessive attempts to Establish BGP session | [optional] 
**AdvertisementInterval** | Pointer to **int32** | The minimum time in seconds between sending route updates to BGP neighbor  | [optional] 
**EbgpMultihop** | Pointer to **int32** | Allows external BGP neighbors to establish peering session multiple network hops away.  | [optional] 
**EgressVlan** | Pointer to **int32** | VLAN used to carry BGP TCP session | [optional] 
**SourceIpAddress** | Pointer to **string** | Source IP address used to override the default source address calculation for BGP TCP session | [optional] [default to ""]
**AnycastIpMask** | Pointer to **string** | The Anycast Address will be used to enable an IP routing redundancy mechanism designed to allow for transparent failover across a leaf pair at the first-hop IP router. | [optional] [default to ""]
**Md5Password** | Pointer to **string** | MD5 password | [optional] [default to ""]
**ImportRouteMap** | Pointer to **string** | A route-map applied to routes imported into the current tenant from the targeted BGP router with the purpose of filtering or modifying the routes | [optional] [default to ""]
**ImportRouteMapRefType** | Pointer to **string** | Object type for import_route_map field | [optional] 
**ExportRouteMap** | Pointer to **string** | A route-map applied to routes exported into the current tenant from the targeted BGP router with the purpose of filtering or modifying the routes | [optional] [default to ""]
**ExportRouteMapRefType** | Pointer to **string** | Object type for export_route_map field | [optional] 
**GatewayMode** | Pointer to **string** | Gateway Mode. Can be BGP, Static, or Default | [optional] [default to "Static BGP"]
**LocalAsNumber** | Pointer to **int32** | Local AS Number | [optional] 
**LocalAsNoPrepend** | Pointer to **bool** | Do not prepend the local-as number to the AS-PATH for routes advertised through this BGP gateway. The Local AS Number must be set for this to be able to be set. | [optional] [default to false]
**ReplaceAs** | Pointer to **bool** | Prepend only Local AS in updates to EBGP peers. | [optional] [default to false]
**MaxLocalAsOccurrences** | Pointer to **NullableInt32** | Allow routes with the local AS number in the AS-path, specifying the maximum occurrences permitted before declaring a routing loop. Leave blank or &#39;0&#39; to disable. | [optional] 
**DynamicBgpSubnet** | Pointer to **string** | Dynamic BGP Subnet | [optional] [default to ""]
**DynamicBgpLimits** | Pointer to **NullableInt32** | Dynamic BGP Limits | [optional] 
**HelperHopIpAddress** | Pointer to **string** | Helper Hop IP Address | [optional] [default to ""]
**EnableBfd** | Pointer to **bool** | Enable BFD(Bi-Directional Forwarding) | [optional] [default to false]
**BfdReceiveInterval** | Pointer to **NullableInt32** | Configure the minimum interval during which the system can receive BFD control packets | [optional] 
**BfdTransmissionInterval** | Pointer to **NullableInt32** | Configure the minimum transmission interval during which the system can send BFD control packets | [optional] 
**BfdDetectMultiplier** | Pointer to **NullableInt32** | Configure the detection multiplier to determine packet loss | [optional] 
**NextHopSelf** | Pointer to **bool** | Optional attribute that disables the normal BGP calculation of next-hops for advertised routes and instead sets the next-hops for advertised routes to the IP address of the switch itself. | [optional] [default to false]
**StaticRoutes** | Pointer to [**[]ConfigPutRequestGatewayGatewayNameStaticRoutesInner**](ConfigPutRequestGatewayGatewayNameStaticRoutesInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties**](ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties.md) |  | [optional] 
**Md5PasswordEncrypted** | Pointer to **string** | MD5 password | [optional] [default to ""]

## Methods

### NewConfigPutRequestGatewayGatewayName

`func NewConfigPutRequestGatewayGatewayName() *ConfigPutRequestGatewayGatewayName`

NewConfigPutRequestGatewayGatewayName instantiates a new ConfigPutRequestGatewayGatewayName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestGatewayGatewayNameWithDefaults

`func NewConfigPutRequestGatewayGatewayNameWithDefaults() *ConfigPutRequestGatewayGatewayName`

NewConfigPutRequestGatewayGatewayNameWithDefaults instantiates a new ConfigPutRequestGatewayGatewayName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestGatewayGatewayName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestGatewayGatewayName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestGatewayGatewayName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestGatewayGatewayName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ConfigPutRequestGatewayGatewayName) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestGatewayGatewayName) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestGatewayGatewayName) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestGatewayGatewayName) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetTenant

`func (o *ConfigPutRequestGatewayGatewayName) GetTenant() string`

GetTenant returns the Tenant field if non-nil, zero value otherwise.

### GetTenantOk

`func (o *ConfigPutRequestGatewayGatewayName) GetTenantOk() (*string, bool)`

GetTenantOk returns a tuple with the Tenant field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenant

`func (o *ConfigPutRequestGatewayGatewayName) SetTenant(v string)`

SetTenant sets Tenant field to given value.

### HasTenant

`func (o *ConfigPutRequestGatewayGatewayName) HasTenant() bool`

HasTenant returns a boolean if a field has been set.

### GetTenantRefType

`func (o *ConfigPutRequestGatewayGatewayName) GetTenantRefType() string`

GetTenantRefType returns the TenantRefType field if non-nil, zero value otherwise.

### GetTenantRefTypeOk

`func (o *ConfigPutRequestGatewayGatewayName) GetTenantRefTypeOk() (*string, bool)`

GetTenantRefTypeOk returns a tuple with the TenantRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantRefType

`func (o *ConfigPutRequestGatewayGatewayName) SetTenantRefType(v string)`

SetTenantRefType sets TenantRefType field to given value.

### HasTenantRefType

`func (o *ConfigPutRequestGatewayGatewayName) HasTenantRefType() bool`

HasTenantRefType returns a boolean if a field has been set.

### GetNeighborIpAddress

`func (o *ConfigPutRequestGatewayGatewayName) GetNeighborIpAddress() string`

GetNeighborIpAddress returns the NeighborIpAddress field if non-nil, zero value otherwise.

### GetNeighborIpAddressOk

`func (o *ConfigPutRequestGatewayGatewayName) GetNeighborIpAddressOk() (*string, bool)`

GetNeighborIpAddressOk returns a tuple with the NeighborIpAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNeighborIpAddress

`func (o *ConfigPutRequestGatewayGatewayName) SetNeighborIpAddress(v string)`

SetNeighborIpAddress sets NeighborIpAddress field to given value.

### HasNeighborIpAddress

`func (o *ConfigPutRequestGatewayGatewayName) HasNeighborIpAddress() bool`

HasNeighborIpAddress returns a boolean if a field has been set.

### GetNeighborAsNumber

`func (o *ConfigPutRequestGatewayGatewayName) GetNeighborAsNumber() int32`

GetNeighborAsNumber returns the NeighborAsNumber field if non-nil, zero value otherwise.

### GetNeighborAsNumberOk

`func (o *ConfigPutRequestGatewayGatewayName) GetNeighborAsNumberOk() (*int32, bool)`

GetNeighborAsNumberOk returns a tuple with the NeighborAsNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNeighborAsNumber

`func (o *ConfigPutRequestGatewayGatewayName) SetNeighborAsNumber(v int32)`

SetNeighborAsNumber sets NeighborAsNumber field to given value.

### HasNeighborAsNumber

`func (o *ConfigPutRequestGatewayGatewayName) HasNeighborAsNumber() bool`

HasNeighborAsNumber returns a boolean if a field has been set.

### SetNeighborAsNumberNil

`func (o *ConfigPutRequestGatewayGatewayName) SetNeighborAsNumberNil(b bool)`

 SetNeighborAsNumberNil sets the value for NeighborAsNumber to be an explicit nil

### UnsetNeighborAsNumber
`func (o *ConfigPutRequestGatewayGatewayName) UnsetNeighborAsNumber()`

UnsetNeighborAsNumber ensures that no value is present for NeighborAsNumber, not even an explicit nil
### GetFabricInterconnect

`func (o *ConfigPutRequestGatewayGatewayName) GetFabricInterconnect() bool`

GetFabricInterconnect returns the FabricInterconnect field if non-nil, zero value otherwise.

### GetFabricInterconnectOk

`func (o *ConfigPutRequestGatewayGatewayName) GetFabricInterconnectOk() (*bool, bool)`

GetFabricInterconnectOk returns a tuple with the FabricInterconnect field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFabricInterconnect

`func (o *ConfigPutRequestGatewayGatewayName) SetFabricInterconnect(v bool)`

SetFabricInterconnect sets FabricInterconnect field to given value.

### HasFabricInterconnect

`func (o *ConfigPutRequestGatewayGatewayName) HasFabricInterconnect() bool`

HasFabricInterconnect returns a boolean if a field has been set.

### GetKeepaliveTimer

`func (o *ConfigPutRequestGatewayGatewayName) GetKeepaliveTimer() int32`

GetKeepaliveTimer returns the KeepaliveTimer field if non-nil, zero value otherwise.

### GetKeepaliveTimerOk

`func (o *ConfigPutRequestGatewayGatewayName) GetKeepaliveTimerOk() (*int32, bool)`

GetKeepaliveTimerOk returns a tuple with the KeepaliveTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeepaliveTimer

`func (o *ConfigPutRequestGatewayGatewayName) SetKeepaliveTimer(v int32)`

SetKeepaliveTimer sets KeepaliveTimer field to given value.

### HasKeepaliveTimer

`func (o *ConfigPutRequestGatewayGatewayName) HasKeepaliveTimer() bool`

HasKeepaliveTimer returns a boolean if a field has been set.

### GetHoldTimer

`func (o *ConfigPutRequestGatewayGatewayName) GetHoldTimer() int32`

GetHoldTimer returns the HoldTimer field if non-nil, zero value otherwise.

### GetHoldTimerOk

`func (o *ConfigPutRequestGatewayGatewayName) GetHoldTimerOk() (*int32, bool)`

GetHoldTimerOk returns a tuple with the HoldTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHoldTimer

`func (o *ConfigPutRequestGatewayGatewayName) SetHoldTimer(v int32)`

SetHoldTimer sets HoldTimer field to given value.

### HasHoldTimer

`func (o *ConfigPutRequestGatewayGatewayName) HasHoldTimer() bool`

HasHoldTimer returns a boolean if a field has been set.

### GetConnectTimer

`func (o *ConfigPutRequestGatewayGatewayName) GetConnectTimer() int32`

GetConnectTimer returns the ConnectTimer field if non-nil, zero value otherwise.

### GetConnectTimerOk

`func (o *ConfigPutRequestGatewayGatewayName) GetConnectTimerOk() (*int32, bool)`

GetConnectTimerOk returns a tuple with the ConnectTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnectTimer

`func (o *ConfigPutRequestGatewayGatewayName) SetConnectTimer(v int32)`

SetConnectTimer sets ConnectTimer field to given value.

### HasConnectTimer

`func (o *ConfigPutRequestGatewayGatewayName) HasConnectTimer() bool`

HasConnectTimer returns a boolean if a field has been set.

### GetAdvertisementInterval

`func (o *ConfigPutRequestGatewayGatewayName) GetAdvertisementInterval() int32`

GetAdvertisementInterval returns the AdvertisementInterval field if non-nil, zero value otherwise.

### GetAdvertisementIntervalOk

`func (o *ConfigPutRequestGatewayGatewayName) GetAdvertisementIntervalOk() (*int32, bool)`

GetAdvertisementIntervalOk returns a tuple with the AdvertisementInterval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAdvertisementInterval

`func (o *ConfigPutRequestGatewayGatewayName) SetAdvertisementInterval(v int32)`

SetAdvertisementInterval sets AdvertisementInterval field to given value.

### HasAdvertisementInterval

`func (o *ConfigPutRequestGatewayGatewayName) HasAdvertisementInterval() bool`

HasAdvertisementInterval returns a boolean if a field has been set.

### GetEbgpMultihop

`func (o *ConfigPutRequestGatewayGatewayName) GetEbgpMultihop() int32`

GetEbgpMultihop returns the EbgpMultihop field if non-nil, zero value otherwise.

### GetEbgpMultihopOk

`func (o *ConfigPutRequestGatewayGatewayName) GetEbgpMultihopOk() (*int32, bool)`

GetEbgpMultihopOk returns a tuple with the EbgpMultihop field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEbgpMultihop

`func (o *ConfigPutRequestGatewayGatewayName) SetEbgpMultihop(v int32)`

SetEbgpMultihop sets EbgpMultihop field to given value.

### HasEbgpMultihop

`func (o *ConfigPutRequestGatewayGatewayName) HasEbgpMultihop() bool`

HasEbgpMultihop returns a boolean if a field has been set.

### GetEgressVlan

`func (o *ConfigPutRequestGatewayGatewayName) GetEgressVlan() int32`

GetEgressVlan returns the EgressVlan field if non-nil, zero value otherwise.

### GetEgressVlanOk

`func (o *ConfigPutRequestGatewayGatewayName) GetEgressVlanOk() (*int32, bool)`

GetEgressVlanOk returns a tuple with the EgressVlan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEgressVlan

`func (o *ConfigPutRequestGatewayGatewayName) SetEgressVlan(v int32)`

SetEgressVlan sets EgressVlan field to given value.

### HasEgressVlan

`func (o *ConfigPutRequestGatewayGatewayName) HasEgressVlan() bool`

HasEgressVlan returns a boolean if a field has been set.

### GetSourceIpAddress

`func (o *ConfigPutRequestGatewayGatewayName) GetSourceIpAddress() string`

GetSourceIpAddress returns the SourceIpAddress field if non-nil, zero value otherwise.

### GetSourceIpAddressOk

`func (o *ConfigPutRequestGatewayGatewayName) GetSourceIpAddressOk() (*string, bool)`

GetSourceIpAddressOk returns a tuple with the SourceIpAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSourceIpAddress

`func (o *ConfigPutRequestGatewayGatewayName) SetSourceIpAddress(v string)`

SetSourceIpAddress sets SourceIpAddress field to given value.

### HasSourceIpAddress

`func (o *ConfigPutRequestGatewayGatewayName) HasSourceIpAddress() bool`

HasSourceIpAddress returns a boolean if a field has been set.

### GetAnycastIpMask

`func (o *ConfigPutRequestGatewayGatewayName) GetAnycastIpMask() string`

GetAnycastIpMask returns the AnycastIpMask field if non-nil, zero value otherwise.

### GetAnycastIpMaskOk

`func (o *ConfigPutRequestGatewayGatewayName) GetAnycastIpMaskOk() (*string, bool)`

GetAnycastIpMaskOk returns a tuple with the AnycastIpMask field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnycastIpMask

`func (o *ConfigPutRequestGatewayGatewayName) SetAnycastIpMask(v string)`

SetAnycastIpMask sets AnycastIpMask field to given value.

### HasAnycastIpMask

`func (o *ConfigPutRequestGatewayGatewayName) HasAnycastIpMask() bool`

HasAnycastIpMask returns a boolean if a field has been set.

### GetMd5Password

`func (o *ConfigPutRequestGatewayGatewayName) GetMd5Password() string`

GetMd5Password returns the Md5Password field if non-nil, zero value otherwise.

### GetMd5PasswordOk

`func (o *ConfigPutRequestGatewayGatewayName) GetMd5PasswordOk() (*string, bool)`

GetMd5PasswordOk returns a tuple with the Md5Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMd5Password

`func (o *ConfigPutRequestGatewayGatewayName) SetMd5Password(v string)`

SetMd5Password sets Md5Password field to given value.

### HasMd5Password

`func (o *ConfigPutRequestGatewayGatewayName) HasMd5Password() bool`

HasMd5Password returns a boolean if a field has been set.

### GetImportRouteMap

`func (o *ConfigPutRequestGatewayGatewayName) GetImportRouteMap() string`

GetImportRouteMap returns the ImportRouteMap field if non-nil, zero value otherwise.

### GetImportRouteMapOk

`func (o *ConfigPutRequestGatewayGatewayName) GetImportRouteMapOk() (*string, bool)`

GetImportRouteMapOk returns a tuple with the ImportRouteMap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImportRouteMap

`func (o *ConfigPutRequestGatewayGatewayName) SetImportRouteMap(v string)`

SetImportRouteMap sets ImportRouteMap field to given value.

### HasImportRouteMap

`func (o *ConfigPutRequestGatewayGatewayName) HasImportRouteMap() bool`

HasImportRouteMap returns a boolean if a field has been set.

### GetImportRouteMapRefType

`func (o *ConfigPutRequestGatewayGatewayName) GetImportRouteMapRefType() string`

GetImportRouteMapRefType returns the ImportRouteMapRefType field if non-nil, zero value otherwise.

### GetImportRouteMapRefTypeOk

`func (o *ConfigPutRequestGatewayGatewayName) GetImportRouteMapRefTypeOk() (*string, bool)`

GetImportRouteMapRefTypeOk returns a tuple with the ImportRouteMapRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImportRouteMapRefType

`func (o *ConfigPutRequestGatewayGatewayName) SetImportRouteMapRefType(v string)`

SetImportRouteMapRefType sets ImportRouteMapRefType field to given value.

### HasImportRouteMapRefType

`func (o *ConfigPutRequestGatewayGatewayName) HasImportRouteMapRefType() bool`

HasImportRouteMapRefType returns a boolean if a field has been set.

### GetExportRouteMap

`func (o *ConfigPutRequestGatewayGatewayName) GetExportRouteMap() string`

GetExportRouteMap returns the ExportRouteMap field if non-nil, zero value otherwise.

### GetExportRouteMapOk

`func (o *ConfigPutRequestGatewayGatewayName) GetExportRouteMapOk() (*string, bool)`

GetExportRouteMapOk returns a tuple with the ExportRouteMap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExportRouteMap

`func (o *ConfigPutRequestGatewayGatewayName) SetExportRouteMap(v string)`

SetExportRouteMap sets ExportRouteMap field to given value.

### HasExportRouteMap

`func (o *ConfigPutRequestGatewayGatewayName) HasExportRouteMap() bool`

HasExportRouteMap returns a boolean if a field has been set.

### GetExportRouteMapRefType

`func (o *ConfigPutRequestGatewayGatewayName) GetExportRouteMapRefType() string`

GetExportRouteMapRefType returns the ExportRouteMapRefType field if non-nil, zero value otherwise.

### GetExportRouteMapRefTypeOk

`func (o *ConfigPutRequestGatewayGatewayName) GetExportRouteMapRefTypeOk() (*string, bool)`

GetExportRouteMapRefTypeOk returns a tuple with the ExportRouteMapRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExportRouteMapRefType

`func (o *ConfigPutRequestGatewayGatewayName) SetExportRouteMapRefType(v string)`

SetExportRouteMapRefType sets ExportRouteMapRefType field to given value.

### HasExportRouteMapRefType

`func (o *ConfigPutRequestGatewayGatewayName) HasExportRouteMapRefType() bool`

HasExportRouteMapRefType returns a boolean if a field has been set.

### GetGatewayMode

`func (o *ConfigPutRequestGatewayGatewayName) GetGatewayMode() string`

GetGatewayMode returns the GatewayMode field if non-nil, zero value otherwise.

### GetGatewayModeOk

`func (o *ConfigPutRequestGatewayGatewayName) GetGatewayModeOk() (*string, bool)`

GetGatewayModeOk returns a tuple with the GatewayMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGatewayMode

`func (o *ConfigPutRequestGatewayGatewayName) SetGatewayMode(v string)`

SetGatewayMode sets GatewayMode field to given value.

### HasGatewayMode

`func (o *ConfigPutRequestGatewayGatewayName) HasGatewayMode() bool`

HasGatewayMode returns a boolean if a field has been set.

### GetLocalAsNumber

`func (o *ConfigPutRequestGatewayGatewayName) GetLocalAsNumber() int32`

GetLocalAsNumber returns the LocalAsNumber field if non-nil, zero value otherwise.

### GetLocalAsNumberOk

`func (o *ConfigPutRequestGatewayGatewayName) GetLocalAsNumberOk() (*int32, bool)`

GetLocalAsNumberOk returns a tuple with the LocalAsNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocalAsNumber

`func (o *ConfigPutRequestGatewayGatewayName) SetLocalAsNumber(v int32)`

SetLocalAsNumber sets LocalAsNumber field to given value.

### HasLocalAsNumber

`func (o *ConfigPutRequestGatewayGatewayName) HasLocalAsNumber() bool`

HasLocalAsNumber returns a boolean if a field has been set.

### GetLocalAsNoPrepend

`func (o *ConfigPutRequestGatewayGatewayName) GetLocalAsNoPrepend() bool`

GetLocalAsNoPrepend returns the LocalAsNoPrepend field if non-nil, zero value otherwise.

### GetLocalAsNoPrependOk

`func (o *ConfigPutRequestGatewayGatewayName) GetLocalAsNoPrependOk() (*bool, bool)`

GetLocalAsNoPrependOk returns a tuple with the LocalAsNoPrepend field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocalAsNoPrepend

`func (o *ConfigPutRequestGatewayGatewayName) SetLocalAsNoPrepend(v bool)`

SetLocalAsNoPrepend sets LocalAsNoPrepend field to given value.

### HasLocalAsNoPrepend

`func (o *ConfigPutRequestGatewayGatewayName) HasLocalAsNoPrepend() bool`

HasLocalAsNoPrepend returns a boolean if a field has been set.

### GetReplaceAs

`func (o *ConfigPutRequestGatewayGatewayName) GetReplaceAs() bool`

GetReplaceAs returns the ReplaceAs field if non-nil, zero value otherwise.

### GetReplaceAsOk

`func (o *ConfigPutRequestGatewayGatewayName) GetReplaceAsOk() (*bool, bool)`

GetReplaceAsOk returns a tuple with the ReplaceAs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReplaceAs

`func (o *ConfigPutRequestGatewayGatewayName) SetReplaceAs(v bool)`

SetReplaceAs sets ReplaceAs field to given value.

### HasReplaceAs

`func (o *ConfigPutRequestGatewayGatewayName) HasReplaceAs() bool`

HasReplaceAs returns a boolean if a field has been set.

### GetMaxLocalAsOccurrences

`func (o *ConfigPutRequestGatewayGatewayName) GetMaxLocalAsOccurrences() int32`

GetMaxLocalAsOccurrences returns the MaxLocalAsOccurrences field if non-nil, zero value otherwise.

### GetMaxLocalAsOccurrencesOk

`func (o *ConfigPutRequestGatewayGatewayName) GetMaxLocalAsOccurrencesOk() (*int32, bool)`

GetMaxLocalAsOccurrencesOk returns a tuple with the MaxLocalAsOccurrences field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxLocalAsOccurrences

`func (o *ConfigPutRequestGatewayGatewayName) SetMaxLocalAsOccurrences(v int32)`

SetMaxLocalAsOccurrences sets MaxLocalAsOccurrences field to given value.

### HasMaxLocalAsOccurrences

`func (o *ConfigPutRequestGatewayGatewayName) HasMaxLocalAsOccurrences() bool`

HasMaxLocalAsOccurrences returns a boolean if a field has been set.

### SetMaxLocalAsOccurrencesNil

`func (o *ConfigPutRequestGatewayGatewayName) SetMaxLocalAsOccurrencesNil(b bool)`

 SetMaxLocalAsOccurrencesNil sets the value for MaxLocalAsOccurrences to be an explicit nil

### UnsetMaxLocalAsOccurrences
`func (o *ConfigPutRequestGatewayGatewayName) UnsetMaxLocalAsOccurrences()`

UnsetMaxLocalAsOccurrences ensures that no value is present for MaxLocalAsOccurrences, not even an explicit nil
### GetDynamicBgpSubnet

`func (o *ConfigPutRequestGatewayGatewayName) GetDynamicBgpSubnet() string`

GetDynamicBgpSubnet returns the DynamicBgpSubnet field if non-nil, zero value otherwise.

### GetDynamicBgpSubnetOk

`func (o *ConfigPutRequestGatewayGatewayName) GetDynamicBgpSubnetOk() (*string, bool)`

GetDynamicBgpSubnetOk returns a tuple with the DynamicBgpSubnet field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDynamicBgpSubnet

`func (o *ConfigPutRequestGatewayGatewayName) SetDynamicBgpSubnet(v string)`

SetDynamicBgpSubnet sets DynamicBgpSubnet field to given value.

### HasDynamicBgpSubnet

`func (o *ConfigPutRequestGatewayGatewayName) HasDynamicBgpSubnet() bool`

HasDynamicBgpSubnet returns a boolean if a field has been set.

### GetDynamicBgpLimits

`func (o *ConfigPutRequestGatewayGatewayName) GetDynamicBgpLimits() int32`

GetDynamicBgpLimits returns the DynamicBgpLimits field if non-nil, zero value otherwise.

### GetDynamicBgpLimitsOk

`func (o *ConfigPutRequestGatewayGatewayName) GetDynamicBgpLimitsOk() (*int32, bool)`

GetDynamicBgpLimitsOk returns a tuple with the DynamicBgpLimits field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDynamicBgpLimits

`func (o *ConfigPutRequestGatewayGatewayName) SetDynamicBgpLimits(v int32)`

SetDynamicBgpLimits sets DynamicBgpLimits field to given value.

### HasDynamicBgpLimits

`func (o *ConfigPutRequestGatewayGatewayName) HasDynamicBgpLimits() bool`

HasDynamicBgpLimits returns a boolean if a field has been set.

### SetDynamicBgpLimitsNil

`func (o *ConfigPutRequestGatewayGatewayName) SetDynamicBgpLimitsNil(b bool)`

 SetDynamicBgpLimitsNil sets the value for DynamicBgpLimits to be an explicit nil

### UnsetDynamicBgpLimits
`func (o *ConfigPutRequestGatewayGatewayName) UnsetDynamicBgpLimits()`

UnsetDynamicBgpLimits ensures that no value is present for DynamicBgpLimits, not even an explicit nil
### GetHelperHopIpAddress

`func (o *ConfigPutRequestGatewayGatewayName) GetHelperHopIpAddress() string`

GetHelperHopIpAddress returns the HelperHopIpAddress field if non-nil, zero value otherwise.

### GetHelperHopIpAddressOk

`func (o *ConfigPutRequestGatewayGatewayName) GetHelperHopIpAddressOk() (*string, bool)`

GetHelperHopIpAddressOk returns a tuple with the HelperHopIpAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHelperHopIpAddress

`func (o *ConfigPutRequestGatewayGatewayName) SetHelperHopIpAddress(v string)`

SetHelperHopIpAddress sets HelperHopIpAddress field to given value.

### HasHelperHopIpAddress

`func (o *ConfigPutRequestGatewayGatewayName) HasHelperHopIpAddress() bool`

HasHelperHopIpAddress returns a boolean if a field has been set.

### GetEnableBfd

`func (o *ConfigPutRequestGatewayGatewayName) GetEnableBfd() bool`

GetEnableBfd returns the EnableBfd field if non-nil, zero value otherwise.

### GetEnableBfdOk

`func (o *ConfigPutRequestGatewayGatewayName) GetEnableBfdOk() (*bool, bool)`

GetEnableBfdOk returns a tuple with the EnableBfd field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnableBfd

`func (o *ConfigPutRequestGatewayGatewayName) SetEnableBfd(v bool)`

SetEnableBfd sets EnableBfd field to given value.

### HasEnableBfd

`func (o *ConfigPutRequestGatewayGatewayName) HasEnableBfd() bool`

HasEnableBfd returns a boolean if a field has been set.

### GetBfdReceiveInterval

`func (o *ConfigPutRequestGatewayGatewayName) GetBfdReceiveInterval() int32`

GetBfdReceiveInterval returns the BfdReceiveInterval field if non-nil, zero value otherwise.

### GetBfdReceiveIntervalOk

`func (o *ConfigPutRequestGatewayGatewayName) GetBfdReceiveIntervalOk() (*int32, bool)`

GetBfdReceiveIntervalOk returns a tuple with the BfdReceiveInterval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBfdReceiveInterval

`func (o *ConfigPutRequestGatewayGatewayName) SetBfdReceiveInterval(v int32)`

SetBfdReceiveInterval sets BfdReceiveInterval field to given value.

### HasBfdReceiveInterval

`func (o *ConfigPutRequestGatewayGatewayName) HasBfdReceiveInterval() bool`

HasBfdReceiveInterval returns a boolean if a field has been set.

### SetBfdReceiveIntervalNil

`func (o *ConfigPutRequestGatewayGatewayName) SetBfdReceiveIntervalNil(b bool)`

 SetBfdReceiveIntervalNil sets the value for BfdReceiveInterval to be an explicit nil

### UnsetBfdReceiveInterval
`func (o *ConfigPutRequestGatewayGatewayName) UnsetBfdReceiveInterval()`

UnsetBfdReceiveInterval ensures that no value is present for BfdReceiveInterval, not even an explicit nil
### GetBfdTransmissionInterval

`func (o *ConfigPutRequestGatewayGatewayName) GetBfdTransmissionInterval() int32`

GetBfdTransmissionInterval returns the BfdTransmissionInterval field if non-nil, zero value otherwise.

### GetBfdTransmissionIntervalOk

`func (o *ConfigPutRequestGatewayGatewayName) GetBfdTransmissionIntervalOk() (*int32, bool)`

GetBfdTransmissionIntervalOk returns a tuple with the BfdTransmissionInterval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBfdTransmissionInterval

`func (o *ConfigPutRequestGatewayGatewayName) SetBfdTransmissionInterval(v int32)`

SetBfdTransmissionInterval sets BfdTransmissionInterval field to given value.

### HasBfdTransmissionInterval

`func (o *ConfigPutRequestGatewayGatewayName) HasBfdTransmissionInterval() bool`

HasBfdTransmissionInterval returns a boolean if a field has been set.

### SetBfdTransmissionIntervalNil

`func (o *ConfigPutRequestGatewayGatewayName) SetBfdTransmissionIntervalNil(b bool)`

 SetBfdTransmissionIntervalNil sets the value for BfdTransmissionInterval to be an explicit nil

### UnsetBfdTransmissionInterval
`func (o *ConfigPutRequestGatewayGatewayName) UnsetBfdTransmissionInterval()`

UnsetBfdTransmissionInterval ensures that no value is present for BfdTransmissionInterval, not even an explicit nil
### GetBfdDetectMultiplier

`func (o *ConfigPutRequestGatewayGatewayName) GetBfdDetectMultiplier() int32`

GetBfdDetectMultiplier returns the BfdDetectMultiplier field if non-nil, zero value otherwise.

### GetBfdDetectMultiplierOk

`func (o *ConfigPutRequestGatewayGatewayName) GetBfdDetectMultiplierOk() (*int32, bool)`

GetBfdDetectMultiplierOk returns a tuple with the BfdDetectMultiplier field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBfdDetectMultiplier

`func (o *ConfigPutRequestGatewayGatewayName) SetBfdDetectMultiplier(v int32)`

SetBfdDetectMultiplier sets BfdDetectMultiplier field to given value.

### HasBfdDetectMultiplier

`func (o *ConfigPutRequestGatewayGatewayName) HasBfdDetectMultiplier() bool`

HasBfdDetectMultiplier returns a boolean if a field has been set.

### SetBfdDetectMultiplierNil

`func (o *ConfigPutRequestGatewayGatewayName) SetBfdDetectMultiplierNil(b bool)`

 SetBfdDetectMultiplierNil sets the value for BfdDetectMultiplier to be an explicit nil

### UnsetBfdDetectMultiplier
`func (o *ConfigPutRequestGatewayGatewayName) UnsetBfdDetectMultiplier()`

UnsetBfdDetectMultiplier ensures that no value is present for BfdDetectMultiplier, not even an explicit nil
### GetNextHopSelf

`func (o *ConfigPutRequestGatewayGatewayName) GetNextHopSelf() bool`

GetNextHopSelf returns the NextHopSelf field if non-nil, zero value otherwise.

### GetNextHopSelfOk

`func (o *ConfigPutRequestGatewayGatewayName) GetNextHopSelfOk() (*bool, bool)`

GetNextHopSelfOk returns a tuple with the NextHopSelf field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNextHopSelf

`func (o *ConfigPutRequestGatewayGatewayName) SetNextHopSelf(v bool)`

SetNextHopSelf sets NextHopSelf field to given value.

### HasNextHopSelf

`func (o *ConfigPutRequestGatewayGatewayName) HasNextHopSelf() bool`

HasNextHopSelf returns a boolean if a field has been set.

### GetStaticRoutes

`func (o *ConfigPutRequestGatewayGatewayName) GetStaticRoutes() []ConfigPutRequestGatewayGatewayNameStaticRoutesInner`

GetStaticRoutes returns the StaticRoutes field if non-nil, zero value otherwise.

### GetStaticRoutesOk

`func (o *ConfigPutRequestGatewayGatewayName) GetStaticRoutesOk() (*[]ConfigPutRequestGatewayGatewayNameStaticRoutesInner, bool)`

GetStaticRoutesOk returns a tuple with the StaticRoutes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStaticRoutes

`func (o *ConfigPutRequestGatewayGatewayName) SetStaticRoutes(v []ConfigPutRequestGatewayGatewayNameStaticRoutesInner)`

SetStaticRoutes sets StaticRoutes field to given value.

### HasStaticRoutes

`func (o *ConfigPutRequestGatewayGatewayName) HasStaticRoutes() bool`

HasStaticRoutes returns a boolean if a field has been set.

### GetObjectProperties

`func (o *ConfigPutRequestGatewayGatewayName) GetObjectProperties() ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ConfigPutRequestGatewayGatewayName) GetObjectPropertiesOk() (*ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ConfigPutRequestGatewayGatewayName) SetObjectProperties(v ConfigPutRequestEthDeviceProfilesEthDeviceProfilesNameObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ConfigPutRequestGatewayGatewayName) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.

### GetMd5PasswordEncrypted

`func (o *ConfigPutRequestGatewayGatewayName) GetMd5PasswordEncrypted() string`

GetMd5PasswordEncrypted returns the Md5PasswordEncrypted field if non-nil, zero value otherwise.

### GetMd5PasswordEncryptedOk

`func (o *ConfigPutRequestGatewayGatewayName) GetMd5PasswordEncryptedOk() (*string, bool)`

GetMd5PasswordEncryptedOk returns a tuple with the Md5PasswordEncrypted field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMd5PasswordEncrypted

`func (o *ConfigPutRequestGatewayGatewayName) SetMd5PasswordEncrypted(v string)`

SetMd5PasswordEncrypted sets Md5PasswordEncrypted field to given value.

### HasMd5PasswordEncrypted

`func (o *ConfigPutRequestGatewayGatewayName) HasMd5PasswordEncrypted() bool`

HasMd5PasswordEncrypted returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


