# ConfigPutRequestSiteSiteName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to true]
**ServiceForSite** | Pointer to **string** | Service for Site | [optional] [default to "service|Management|"]
**ServiceForSiteRefType** | Pointer to **string** | Object type for service_for_site field | [optional] 
**SpanningTreeType** | Pointer to **string** | Sets the spanning tree type for all Ports in this Site with Spanning Tree enabled | [optional] [default to "pvst"]
**RegionName** | Pointer to **string** | Defines the logical boundary of the network. All switches in an MSTP region must have the same configured region name | [optional] [default to ""]
**Revision** | Pointer to **NullableInt32** | A logical number that signifies a revision for the MSTP configuration. All switches in an MSTP region must have the same revision number | [optional] 
**ForceSpanningTreeOnFabricPorts** | Pointer to **bool** | Enable spanning tree on all fabric connections.  This overrides the Eth Port Settings for Fabric ports | [optional] [default to false]
**ReadOnlyMode** | Pointer to **bool** | When Read Only Mode is checked, vNetC will perform all functions except writing database updates to the target hardware | [optional] [default to false]
**DscpToPBitMap** | Pointer to **string** | For any Service that is using DSCP to p-bit map packet prioritization. A string of length 64 with a 0-7 in each position | [optional] [default to "0000000011111111222222223333333344444444555555556666666677777777"]
**AnycastMacAddress** | Pointer to **string** | Site Level MAC Address for Anycast | [optional] [default to "(auto)"]
**AnycastMacAddressAutoAssigned** | Pointer to **bool** | Whether or not the value in anycast_mac_address field has been automatically assigned or not. Set to false and change anycast_mac_address value to edit. | [optional] 
**MacAddressAgingTime** | Pointer to **int32** | MAC Address Aging Time | [optional] 
**MlagDelayRestoreTimer** | Pointer to **int32** | MLAG Delay Restore Timer | [optional] 
**BgpKeepaliveTimer** | Pointer to **int32** | Spine BGP Keepalive Timer | [optional] 
**BgpHoldDownTimer** | Pointer to **int32** | Spine BGP Hold Down Timer | [optional] 
**SpineBgpAdvertisementInterval** | Pointer to **int32** | BGP Advertisement Interval for spines/superspines. Use \&quot;0\&quot; for immediate updates | [optional] 
**SpineBgpConnectTimer** | Pointer to **int32** | BGP Connect Timer | [optional] 
**LeafBgpKeepAliveTimer** | Pointer to **int32** | Leaf BGP Keep Alive Timer | [optional] 
**LeafBgpHoldDownTimer** | Pointer to **int32** | Leaf BGP Hold Down Timer | [optional] 
**LeafBgpAdvertisementInterval** | Pointer to **int32** | BGP Advertisement Interval for leafs. Use \&quot;0\&quot; for immediate updates | [optional] 
**LeafBgpConnectTimer** | Pointer to **int32** | BGP Connect Timer | [optional] 
**LinkStateTimeoutValue** | Pointer to **NullableInt32** | Link State Timeout Value | [optional] 
**EvpnMultihomingStartupDelay** | Pointer to **NullableInt32** | Startup Delay | [optional] 
**EvpnMacHoldtime** | Pointer to **NullableInt32** | MAC Holdtime | [optional] 
**EvpnNeighborHoldtime** | Pointer to **NullableInt32** | Neighbor Holdtime | [optional] 
**AggressiveReporting** | Pointer to **bool** | Fast Reporting of Switch Communications, Link Up/Down, and BGP Status | [optional] [default to true]
**CrcFailureThreshold** | Pointer to **NullableInt32** | Threshold in Errors per second that when met will disable the links as part of LAGs | [optional] 
**Islands** | Pointer to [**[]ConfigPutRequestSiteSiteNameIslandsInner**](ConfigPutRequestSiteSiteNameIslandsInner.md) |  | [optional] 
**Pairs** | Pointer to [**[]ConfigPutRequestSiteSiteNamePairsInner**](ConfigPutRequestSiteSiteNamePairsInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**ConfigPutRequestSiteSiteNameObjectProperties**](ConfigPutRequestSiteSiteNameObjectProperties.md) |  | [optional] 

## Methods

### NewConfigPutRequestSiteSiteName

`func NewConfigPutRequestSiteSiteName() *ConfigPutRequestSiteSiteName`

NewConfigPutRequestSiteSiteName instantiates a new ConfigPutRequestSiteSiteName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestSiteSiteNameWithDefaults

`func NewConfigPutRequestSiteSiteNameWithDefaults() *ConfigPutRequestSiteSiteName`

NewConfigPutRequestSiteSiteNameWithDefaults instantiates a new ConfigPutRequestSiteSiteName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestSiteSiteName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestSiteSiteName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestSiteSiteName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestSiteSiteName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ConfigPutRequestSiteSiteName) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestSiteSiteName) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestSiteSiteName) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestSiteSiteName) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetServiceForSite

`func (o *ConfigPutRequestSiteSiteName) GetServiceForSite() string`

GetServiceForSite returns the ServiceForSite field if non-nil, zero value otherwise.

### GetServiceForSiteOk

`func (o *ConfigPutRequestSiteSiteName) GetServiceForSiteOk() (*string, bool)`

GetServiceForSiteOk returns a tuple with the ServiceForSite field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServiceForSite

`func (o *ConfigPutRequestSiteSiteName) SetServiceForSite(v string)`

SetServiceForSite sets ServiceForSite field to given value.

### HasServiceForSite

`func (o *ConfigPutRequestSiteSiteName) HasServiceForSite() bool`

HasServiceForSite returns a boolean if a field has been set.

### GetServiceForSiteRefType

`func (o *ConfigPutRequestSiteSiteName) GetServiceForSiteRefType() string`

GetServiceForSiteRefType returns the ServiceForSiteRefType field if non-nil, zero value otherwise.

### GetServiceForSiteRefTypeOk

`func (o *ConfigPutRequestSiteSiteName) GetServiceForSiteRefTypeOk() (*string, bool)`

GetServiceForSiteRefTypeOk returns a tuple with the ServiceForSiteRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServiceForSiteRefType

`func (o *ConfigPutRequestSiteSiteName) SetServiceForSiteRefType(v string)`

SetServiceForSiteRefType sets ServiceForSiteRefType field to given value.

### HasServiceForSiteRefType

`func (o *ConfigPutRequestSiteSiteName) HasServiceForSiteRefType() bool`

HasServiceForSiteRefType returns a boolean if a field has been set.

### GetSpanningTreeType

`func (o *ConfigPutRequestSiteSiteName) GetSpanningTreeType() string`

GetSpanningTreeType returns the SpanningTreeType field if non-nil, zero value otherwise.

### GetSpanningTreeTypeOk

`func (o *ConfigPutRequestSiteSiteName) GetSpanningTreeTypeOk() (*string, bool)`

GetSpanningTreeTypeOk returns a tuple with the SpanningTreeType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpanningTreeType

`func (o *ConfigPutRequestSiteSiteName) SetSpanningTreeType(v string)`

SetSpanningTreeType sets SpanningTreeType field to given value.

### HasSpanningTreeType

`func (o *ConfigPutRequestSiteSiteName) HasSpanningTreeType() bool`

HasSpanningTreeType returns a boolean if a field has been set.

### GetRegionName

`func (o *ConfigPutRequestSiteSiteName) GetRegionName() string`

GetRegionName returns the RegionName field if non-nil, zero value otherwise.

### GetRegionNameOk

`func (o *ConfigPutRequestSiteSiteName) GetRegionNameOk() (*string, bool)`

GetRegionNameOk returns a tuple with the RegionName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegionName

`func (o *ConfigPutRequestSiteSiteName) SetRegionName(v string)`

SetRegionName sets RegionName field to given value.

### HasRegionName

`func (o *ConfigPutRequestSiteSiteName) HasRegionName() bool`

HasRegionName returns a boolean if a field has been set.

### GetRevision

`func (o *ConfigPutRequestSiteSiteName) GetRevision() int32`

GetRevision returns the Revision field if non-nil, zero value otherwise.

### GetRevisionOk

`func (o *ConfigPutRequestSiteSiteName) GetRevisionOk() (*int32, bool)`

GetRevisionOk returns a tuple with the Revision field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRevision

`func (o *ConfigPutRequestSiteSiteName) SetRevision(v int32)`

SetRevision sets Revision field to given value.

### HasRevision

`func (o *ConfigPutRequestSiteSiteName) HasRevision() bool`

HasRevision returns a boolean if a field has been set.

### SetRevisionNil

`func (o *ConfigPutRequestSiteSiteName) SetRevisionNil(b bool)`

 SetRevisionNil sets the value for Revision to be an explicit nil

### UnsetRevision
`func (o *ConfigPutRequestSiteSiteName) UnsetRevision()`

UnsetRevision ensures that no value is present for Revision, not even an explicit nil
### GetForceSpanningTreeOnFabricPorts

`func (o *ConfigPutRequestSiteSiteName) GetForceSpanningTreeOnFabricPorts() bool`

GetForceSpanningTreeOnFabricPorts returns the ForceSpanningTreeOnFabricPorts field if non-nil, zero value otherwise.

### GetForceSpanningTreeOnFabricPortsOk

`func (o *ConfigPutRequestSiteSiteName) GetForceSpanningTreeOnFabricPortsOk() (*bool, bool)`

GetForceSpanningTreeOnFabricPortsOk returns a tuple with the ForceSpanningTreeOnFabricPorts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetForceSpanningTreeOnFabricPorts

`func (o *ConfigPutRequestSiteSiteName) SetForceSpanningTreeOnFabricPorts(v bool)`

SetForceSpanningTreeOnFabricPorts sets ForceSpanningTreeOnFabricPorts field to given value.

### HasForceSpanningTreeOnFabricPorts

`func (o *ConfigPutRequestSiteSiteName) HasForceSpanningTreeOnFabricPorts() bool`

HasForceSpanningTreeOnFabricPorts returns a boolean if a field has been set.

### GetReadOnlyMode

`func (o *ConfigPutRequestSiteSiteName) GetReadOnlyMode() bool`

GetReadOnlyMode returns the ReadOnlyMode field if non-nil, zero value otherwise.

### GetReadOnlyModeOk

`func (o *ConfigPutRequestSiteSiteName) GetReadOnlyModeOk() (*bool, bool)`

GetReadOnlyModeOk returns a tuple with the ReadOnlyMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReadOnlyMode

`func (o *ConfigPutRequestSiteSiteName) SetReadOnlyMode(v bool)`

SetReadOnlyMode sets ReadOnlyMode field to given value.

### HasReadOnlyMode

`func (o *ConfigPutRequestSiteSiteName) HasReadOnlyMode() bool`

HasReadOnlyMode returns a boolean if a field has been set.

### GetDscpToPBitMap

`func (o *ConfigPutRequestSiteSiteName) GetDscpToPBitMap() string`

GetDscpToPBitMap returns the DscpToPBitMap field if non-nil, zero value otherwise.

### GetDscpToPBitMapOk

`func (o *ConfigPutRequestSiteSiteName) GetDscpToPBitMapOk() (*string, bool)`

GetDscpToPBitMapOk returns a tuple with the DscpToPBitMap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDscpToPBitMap

`func (o *ConfigPutRequestSiteSiteName) SetDscpToPBitMap(v string)`

SetDscpToPBitMap sets DscpToPBitMap field to given value.

### HasDscpToPBitMap

`func (o *ConfigPutRequestSiteSiteName) HasDscpToPBitMap() bool`

HasDscpToPBitMap returns a boolean if a field has been set.

### GetAnycastMacAddress

`func (o *ConfigPutRequestSiteSiteName) GetAnycastMacAddress() string`

GetAnycastMacAddress returns the AnycastMacAddress field if non-nil, zero value otherwise.

### GetAnycastMacAddressOk

`func (o *ConfigPutRequestSiteSiteName) GetAnycastMacAddressOk() (*string, bool)`

GetAnycastMacAddressOk returns a tuple with the AnycastMacAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnycastMacAddress

`func (o *ConfigPutRequestSiteSiteName) SetAnycastMacAddress(v string)`

SetAnycastMacAddress sets AnycastMacAddress field to given value.

### HasAnycastMacAddress

`func (o *ConfigPutRequestSiteSiteName) HasAnycastMacAddress() bool`

HasAnycastMacAddress returns a boolean if a field has been set.

### GetAnycastMacAddressAutoAssigned

`func (o *ConfigPutRequestSiteSiteName) GetAnycastMacAddressAutoAssigned() bool`

GetAnycastMacAddressAutoAssigned returns the AnycastMacAddressAutoAssigned field if non-nil, zero value otherwise.

### GetAnycastMacAddressAutoAssignedOk

`func (o *ConfigPutRequestSiteSiteName) GetAnycastMacAddressAutoAssignedOk() (*bool, bool)`

GetAnycastMacAddressAutoAssignedOk returns a tuple with the AnycastMacAddressAutoAssigned field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnycastMacAddressAutoAssigned

`func (o *ConfigPutRequestSiteSiteName) SetAnycastMacAddressAutoAssigned(v bool)`

SetAnycastMacAddressAutoAssigned sets AnycastMacAddressAutoAssigned field to given value.

### HasAnycastMacAddressAutoAssigned

`func (o *ConfigPutRequestSiteSiteName) HasAnycastMacAddressAutoAssigned() bool`

HasAnycastMacAddressAutoAssigned returns a boolean if a field has been set.

### GetMacAddressAgingTime

`func (o *ConfigPutRequestSiteSiteName) GetMacAddressAgingTime() int32`

GetMacAddressAgingTime returns the MacAddressAgingTime field if non-nil, zero value otherwise.

### GetMacAddressAgingTimeOk

`func (o *ConfigPutRequestSiteSiteName) GetMacAddressAgingTimeOk() (*int32, bool)`

GetMacAddressAgingTimeOk returns a tuple with the MacAddressAgingTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMacAddressAgingTime

`func (o *ConfigPutRequestSiteSiteName) SetMacAddressAgingTime(v int32)`

SetMacAddressAgingTime sets MacAddressAgingTime field to given value.

### HasMacAddressAgingTime

`func (o *ConfigPutRequestSiteSiteName) HasMacAddressAgingTime() bool`

HasMacAddressAgingTime returns a boolean if a field has been set.

### GetMlagDelayRestoreTimer

`func (o *ConfigPutRequestSiteSiteName) GetMlagDelayRestoreTimer() int32`

GetMlagDelayRestoreTimer returns the MlagDelayRestoreTimer field if non-nil, zero value otherwise.

### GetMlagDelayRestoreTimerOk

`func (o *ConfigPutRequestSiteSiteName) GetMlagDelayRestoreTimerOk() (*int32, bool)`

GetMlagDelayRestoreTimerOk returns a tuple with the MlagDelayRestoreTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMlagDelayRestoreTimer

`func (o *ConfigPutRequestSiteSiteName) SetMlagDelayRestoreTimer(v int32)`

SetMlagDelayRestoreTimer sets MlagDelayRestoreTimer field to given value.

### HasMlagDelayRestoreTimer

`func (o *ConfigPutRequestSiteSiteName) HasMlagDelayRestoreTimer() bool`

HasMlagDelayRestoreTimer returns a boolean if a field has been set.

### GetBgpKeepaliveTimer

`func (o *ConfigPutRequestSiteSiteName) GetBgpKeepaliveTimer() int32`

GetBgpKeepaliveTimer returns the BgpKeepaliveTimer field if non-nil, zero value otherwise.

### GetBgpKeepaliveTimerOk

`func (o *ConfigPutRequestSiteSiteName) GetBgpKeepaliveTimerOk() (*int32, bool)`

GetBgpKeepaliveTimerOk returns a tuple with the BgpKeepaliveTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBgpKeepaliveTimer

`func (o *ConfigPutRequestSiteSiteName) SetBgpKeepaliveTimer(v int32)`

SetBgpKeepaliveTimer sets BgpKeepaliveTimer field to given value.

### HasBgpKeepaliveTimer

`func (o *ConfigPutRequestSiteSiteName) HasBgpKeepaliveTimer() bool`

HasBgpKeepaliveTimer returns a boolean if a field has been set.

### GetBgpHoldDownTimer

`func (o *ConfigPutRequestSiteSiteName) GetBgpHoldDownTimer() int32`

GetBgpHoldDownTimer returns the BgpHoldDownTimer field if non-nil, zero value otherwise.

### GetBgpHoldDownTimerOk

`func (o *ConfigPutRequestSiteSiteName) GetBgpHoldDownTimerOk() (*int32, bool)`

GetBgpHoldDownTimerOk returns a tuple with the BgpHoldDownTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBgpHoldDownTimer

`func (o *ConfigPutRequestSiteSiteName) SetBgpHoldDownTimer(v int32)`

SetBgpHoldDownTimer sets BgpHoldDownTimer field to given value.

### HasBgpHoldDownTimer

`func (o *ConfigPutRequestSiteSiteName) HasBgpHoldDownTimer() bool`

HasBgpHoldDownTimer returns a boolean if a field has been set.

### GetSpineBgpAdvertisementInterval

`func (o *ConfigPutRequestSiteSiteName) GetSpineBgpAdvertisementInterval() int32`

GetSpineBgpAdvertisementInterval returns the SpineBgpAdvertisementInterval field if non-nil, zero value otherwise.

### GetSpineBgpAdvertisementIntervalOk

`func (o *ConfigPutRequestSiteSiteName) GetSpineBgpAdvertisementIntervalOk() (*int32, bool)`

GetSpineBgpAdvertisementIntervalOk returns a tuple with the SpineBgpAdvertisementInterval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpineBgpAdvertisementInterval

`func (o *ConfigPutRequestSiteSiteName) SetSpineBgpAdvertisementInterval(v int32)`

SetSpineBgpAdvertisementInterval sets SpineBgpAdvertisementInterval field to given value.

### HasSpineBgpAdvertisementInterval

`func (o *ConfigPutRequestSiteSiteName) HasSpineBgpAdvertisementInterval() bool`

HasSpineBgpAdvertisementInterval returns a boolean if a field has been set.

### GetSpineBgpConnectTimer

`func (o *ConfigPutRequestSiteSiteName) GetSpineBgpConnectTimer() int32`

GetSpineBgpConnectTimer returns the SpineBgpConnectTimer field if non-nil, zero value otherwise.

### GetSpineBgpConnectTimerOk

`func (o *ConfigPutRequestSiteSiteName) GetSpineBgpConnectTimerOk() (*int32, bool)`

GetSpineBgpConnectTimerOk returns a tuple with the SpineBgpConnectTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpineBgpConnectTimer

`func (o *ConfigPutRequestSiteSiteName) SetSpineBgpConnectTimer(v int32)`

SetSpineBgpConnectTimer sets SpineBgpConnectTimer field to given value.

### HasSpineBgpConnectTimer

`func (o *ConfigPutRequestSiteSiteName) HasSpineBgpConnectTimer() bool`

HasSpineBgpConnectTimer returns a boolean if a field has been set.

### GetLeafBgpKeepAliveTimer

`func (o *ConfigPutRequestSiteSiteName) GetLeafBgpKeepAliveTimer() int32`

GetLeafBgpKeepAliveTimer returns the LeafBgpKeepAliveTimer field if non-nil, zero value otherwise.

### GetLeafBgpKeepAliveTimerOk

`func (o *ConfigPutRequestSiteSiteName) GetLeafBgpKeepAliveTimerOk() (*int32, bool)`

GetLeafBgpKeepAliveTimerOk returns a tuple with the LeafBgpKeepAliveTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLeafBgpKeepAliveTimer

`func (o *ConfigPutRequestSiteSiteName) SetLeafBgpKeepAliveTimer(v int32)`

SetLeafBgpKeepAliveTimer sets LeafBgpKeepAliveTimer field to given value.

### HasLeafBgpKeepAliveTimer

`func (o *ConfigPutRequestSiteSiteName) HasLeafBgpKeepAliveTimer() bool`

HasLeafBgpKeepAliveTimer returns a boolean if a field has been set.

### GetLeafBgpHoldDownTimer

`func (o *ConfigPutRequestSiteSiteName) GetLeafBgpHoldDownTimer() int32`

GetLeafBgpHoldDownTimer returns the LeafBgpHoldDownTimer field if non-nil, zero value otherwise.

### GetLeafBgpHoldDownTimerOk

`func (o *ConfigPutRequestSiteSiteName) GetLeafBgpHoldDownTimerOk() (*int32, bool)`

GetLeafBgpHoldDownTimerOk returns a tuple with the LeafBgpHoldDownTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLeafBgpHoldDownTimer

`func (o *ConfigPutRequestSiteSiteName) SetLeafBgpHoldDownTimer(v int32)`

SetLeafBgpHoldDownTimer sets LeafBgpHoldDownTimer field to given value.

### HasLeafBgpHoldDownTimer

`func (o *ConfigPutRequestSiteSiteName) HasLeafBgpHoldDownTimer() bool`

HasLeafBgpHoldDownTimer returns a boolean if a field has been set.

### GetLeafBgpAdvertisementInterval

`func (o *ConfigPutRequestSiteSiteName) GetLeafBgpAdvertisementInterval() int32`

GetLeafBgpAdvertisementInterval returns the LeafBgpAdvertisementInterval field if non-nil, zero value otherwise.

### GetLeafBgpAdvertisementIntervalOk

`func (o *ConfigPutRequestSiteSiteName) GetLeafBgpAdvertisementIntervalOk() (*int32, bool)`

GetLeafBgpAdvertisementIntervalOk returns a tuple with the LeafBgpAdvertisementInterval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLeafBgpAdvertisementInterval

`func (o *ConfigPutRequestSiteSiteName) SetLeafBgpAdvertisementInterval(v int32)`

SetLeafBgpAdvertisementInterval sets LeafBgpAdvertisementInterval field to given value.

### HasLeafBgpAdvertisementInterval

`func (o *ConfigPutRequestSiteSiteName) HasLeafBgpAdvertisementInterval() bool`

HasLeafBgpAdvertisementInterval returns a boolean if a field has been set.

### GetLeafBgpConnectTimer

`func (o *ConfigPutRequestSiteSiteName) GetLeafBgpConnectTimer() int32`

GetLeafBgpConnectTimer returns the LeafBgpConnectTimer field if non-nil, zero value otherwise.

### GetLeafBgpConnectTimerOk

`func (o *ConfigPutRequestSiteSiteName) GetLeafBgpConnectTimerOk() (*int32, bool)`

GetLeafBgpConnectTimerOk returns a tuple with the LeafBgpConnectTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLeafBgpConnectTimer

`func (o *ConfigPutRequestSiteSiteName) SetLeafBgpConnectTimer(v int32)`

SetLeafBgpConnectTimer sets LeafBgpConnectTimer field to given value.

### HasLeafBgpConnectTimer

`func (o *ConfigPutRequestSiteSiteName) HasLeafBgpConnectTimer() bool`

HasLeafBgpConnectTimer returns a boolean if a field has been set.

### GetLinkStateTimeoutValue

`func (o *ConfigPutRequestSiteSiteName) GetLinkStateTimeoutValue() int32`

GetLinkStateTimeoutValue returns the LinkStateTimeoutValue field if non-nil, zero value otherwise.

### GetLinkStateTimeoutValueOk

`func (o *ConfigPutRequestSiteSiteName) GetLinkStateTimeoutValueOk() (*int32, bool)`

GetLinkStateTimeoutValueOk returns a tuple with the LinkStateTimeoutValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLinkStateTimeoutValue

`func (o *ConfigPutRequestSiteSiteName) SetLinkStateTimeoutValue(v int32)`

SetLinkStateTimeoutValue sets LinkStateTimeoutValue field to given value.

### HasLinkStateTimeoutValue

`func (o *ConfigPutRequestSiteSiteName) HasLinkStateTimeoutValue() bool`

HasLinkStateTimeoutValue returns a boolean if a field has been set.

### SetLinkStateTimeoutValueNil

`func (o *ConfigPutRequestSiteSiteName) SetLinkStateTimeoutValueNil(b bool)`

 SetLinkStateTimeoutValueNil sets the value for LinkStateTimeoutValue to be an explicit nil

### UnsetLinkStateTimeoutValue
`func (o *ConfigPutRequestSiteSiteName) UnsetLinkStateTimeoutValue()`

UnsetLinkStateTimeoutValue ensures that no value is present for LinkStateTimeoutValue, not even an explicit nil
### GetEvpnMultihomingStartupDelay

`func (o *ConfigPutRequestSiteSiteName) GetEvpnMultihomingStartupDelay() int32`

GetEvpnMultihomingStartupDelay returns the EvpnMultihomingStartupDelay field if non-nil, zero value otherwise.

### GetEvpnMultihomingStartupDelayOk

`func (o *ConfigPutRequestSiteSiteName) GetEvpnMultihomingStartupDelayOk() (*int32, bool)`

GetEvpnMultihomingStartupDelayOk returns a tuple with the EvpnMultihomingStartupDelay field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEvpnMultihomingStartupDelay

`func (o *ConfigPutRequestSiteSiteName) SetEvpnMultihomingStartupDelay(v int32)`

SetEvpnMultihomingStartupDelay sets EvpnMultihomingStartupDelay field to given value.

### HasEvpnMultihomingStartupDelay

`func (o *ConfigPutRequestSiteSiteName) HasEvpnMultihomingStartupDelay() bool`

HasEvpnMultihomingStartupDelay returns a boolean if a field has been set.

### SetEvpnMultihomingStartupDelayNil

`func (o *ConfigPutRequestSiteSiteName) SetEvpnMultihomingStartupDelayNil(b bool)`

 SetEvpnMultihomingStartupDelayNil sets the value for EvpnMultihomingStartupDelay to be an explicit nil

### UnsetEvpnMultihomingStartupDelay
`func (o *ConfigPutRequestSiteSiteName) UnsetEvpnMultihomingStartupDelay()`

UnsetEvpnMultihomingStartupDelay ensures that no value is present for EvpnMultihomingStartupDelay, not even an explicit nil
### GetEvpnMacHoldtime

`func (o *ConfigPutRequestSiteSiteName) GetEvpnMacHoldtime() int32`

GetEvpnMacHoldtime returns the EvpnMacHoldtime field if non-nil, zero value otherwise.

### GetEvpnMacHoldtimeOk

`func (o *ConfigPutRequestSiteSiteName) GetEvpnMacHoldtimeOk() (*int32, bool)`

GetEvpnMacHoldtimeOk returns a tuple with the EvpnMacHoldtime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEvpnMacHoldtime

`func (o *ConfigPutRequestSiteSiteName) SetEvpnMacHoldtime(v int32)`

SetEvpnMacHoldtime sets EvpnMacHoldtime field to given value.

### HasEvpnMacHoldtime

`func (o *ConfigPutRequestSiteSiteName) HasEvpnMacHoldtime() bool`

HasEvpnMacHoldtime returns a boolean if a field has been set.

### SetEvpnMacHoldtimeNil

`func (o *ConfigPutRequestSiteSiteName) SetEvpnMacHoldtimeNil(b bool)`

 SetEvpnMacHoldtimeNil sets the value for EvpnMacHoldtime to be an explicit nil

### UnsetEvpnMacHoldtime
`func (o *ConfigPutRequestSiteSiteName) UnsetEvpnMacHoldtime()`

UnsetEvpnMacHoldtime ensures that no value is present for EvpnMacHoldtime, not even an explicit nil
### GetEvpnNeighborHoldtime

`func (o *ConfigPutRequestSiteSiteName) GetEvpnNeighborHoldtime() int32`

GetEvpnNeighborHoldtime returns the EvpnNeighborHoldtime field if non-nil, zero value otherwise.

### GetEvpnNeighborHoldtimeOk

`func (o *ConfigPutRequestSiteSiteName) GetEvpnNeighborHoldtimeOk() (*int32, bool)`

GetEvpnNeighborHoldtimeOk returns a tuple with the EvpnNeighborHoldtime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEvpnNeighborHoldtime

`func (o *ConfigPutRequestSiteSiteName) SetEvpnNeighborHoldtime(v int32)`

SetEvpnNeighborHoldtime sets EvpnNeighborHoldtime field to given value.

### HasEvpnNeighborHoldtime

`func (o *ConfigPutRequestSiteSiteName) HasEvpnNeighborHoldtime() bool`

HasEvpnNeighborHoldtime returns a boolean if a field has been set.

### SetEvpnNeighborHoldtimeNil

`func (o *ConfigPutRequestSiteSiteName) SetEvpnNeighborHoldtimeNil(b bool)`

 SetEvpnNeighborHoldtimeNil sets the value for EvpnNeighborHoldtime to be an explicit nil

### UnsetEvpnNeighborHoldtime
`func (o *ConfigPutRequestSiteSiteName) UnsetEvpnNeighborHoldtime()`

UnsetEvpnNeighborHoldtime ensures that no value is present for EvpnNeighborHoldtime, not even an explicit nil
### GetAggressiveReporting

`func (o *ConfigPutRequestSiteSiteName) GetAggressiveReporting() bool`

GetAggressiveReporting returns the AggressiveReporting field if non-nil, zero value otherwise.

### GetAggressiveReportingOk

`func (o *ConfigPutRequestSiteSiteName) GetAggressiveReportingOk() (*bool, bool)`

GetAggressiveReportingOk returns a tuple with the AggressiveReporting field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAggressiveReporting

`func (o *ConfigPutRequestSiteSiteName) SetAggressiveReporting(v bool)`

SetAggressiveReporting sets AggressiveReporting field to given value.

### HasAggressiveReporting

`func (o *ConfigPutRequestSiteSiteName) HasAggressiveReporting() bool`

HasAggressiveReporting returns a boolean if a field has been set.

### GetCrcFailureThreshold

`func (o *ConfigPutRequestSiteSiteName) GetCrcFailureThreshold() int32`

GetCrcFailureThreshold returns the CrcFailureThreshold field if non-nil, zero value otherwise.

### GetCrcFailureThresholdOk

`func (o *ConfigPutRequestSiteSiteName) GetCrcFailureThresholdOk() (*int32, bool)`

GetCrcFailureThresholdOk returns a tuple with the CrcFailureThreshold field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCrcFailureThreshold

`func (o *ConfigPutRequestSiteSiteName) SetCrcFailureThreshold(v int32)`

SetCrcFailureThreshold sets CrcFailureThreshold field to given value.

### HasCrcFailureThreshold

`func (o *ConfigPutRequestSiteSiteName) HasCrcFailureThreshold() bool`

HasCrcFailureThreshold returns a boolean if a field has been set.

### SetCrcFailureThresholdNil

`func (o *ConfigPutRequestSiteSiteName) SetCrcFailureThresholdNil(b bool)`

 SetCrcFailureThresholdNil sets the value for CrcFailureThreshold to be an explicit nil

### UnsetCrcFailureThreshold
`func (o *ConfigPutRequestSiteSiteName) UnsetCrcFailureThreshold()`

UnsetCrcFailureThreshold ensures that no value is present for CrcFailureThreshold, not even an explicit nil
### GetIslands

`func (o *ConfigPutRequestSiteSiteName) GetIslands() []ConfigPutRequestSiteSiteNameIslandsInner`

GetIslands returns the Islands field if non-nil, zero value otherwise.

### GetIslandsOk

`func (o *ConfigPutRequestSiteSiteName) GetIslandsOk() (*[]ConfigPutRequestSiteSiteNameIslandsInner, bool)`

GetIslandsOk returns a tuple with the Islands field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIslands

`func (o *ConfigPutRequestSiteSiteName) SetIslands(v []ConfigPutRequestSiteSiteNameIslandsInner)`

SetIslands sets Islands field to given value.

### HasIslands

`func (o *ConfigPutRequestSiteSiteName) HasIslands() bool`

HasIslands returns a boolean if a field has been set.

### GetPairs

`func (o *ConfigPutRequestSiteSiteName) GetPairs() []ConfigPutRequestSiteSiteNamePairsInner`

GetPairs returns the Pairs field if non-nil, zero value otherwise.

### GetPairsOk

`func (o *ConfigPutRequestSiteSiteName) GetPairsOk() (*[]ConfigPutRequestSiteSiteNamePairsInner, bool)`

GetPairsOk returns a tuple with the Pairs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPairs

`func (o *ConfigPutRequestSiteSiteName) SetPairs(v []ConfigPutRequestSiteSiteNamePairsInner)`

SetPairs sets Pairs field to given value.

### HasPairs

`func (o *ConfigPutRequestSiteSiteName) HasPairs() bool`

HasPairs returns a boolean if a field has been set.

### GetObjectProperties

`func (o *ConfigPutRequestSiteSiteName) GetObjectProperties() ConfigPutRequestSiteSiteNameObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ConfigPutRequestSiteSiteName) GetObjectPropertiesOk() (*ConfigPutRequestSiteSiteNameObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ConfigPutRequestSiteSiteName) SetObjectProperties(v ConfigPutRequestSiteSiteNameObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ConfigPutRequestSiteSiteName) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


