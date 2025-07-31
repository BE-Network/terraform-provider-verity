# SitesPatchRequestSiteValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to true]
**ServiceForSite** | Pointer to **string** | Service for Site | [optional] [default to "service|Management|"]
**ServiceForSiteRefType** | Pointer to **string** | Object type for service_for_site field | [optional] 
**SpanningTreeType** | Pointer to **string** | Sets the spanning tree type for all Ports in this Site with Spanning Tree enabled | [optional] [default to "pvst"]
**RegionName** | Pointer to **string** | Defines the logical boundary of the network. All switches in an MSTP region must have the same configured region name | [optional] [default to ""]
**Revision** | Pointer to **NullableInt32** | A logical number that signifies a revision for the MSTP configuration. All switches in an MSTP region must have the same revision number | [optional] [default to 0]
**ForceSpanningTreeOnFabricPorts** | Pointer to **bool** | Enable spanning tree on all fabric connections.  This overrides the Eth Port Settings for Fabric ports | [optional] [default to false]
**ReadOnlyMode** | Pointer to **bool** | When Read Only Mode is checked, vNetC will perform all functions except writing database updates to the target hardware | [optional] [default to false]
**DscpToPBitMap** | Pointer to **string** | For any Service that is using DSCP to p-bit map packet prioritization. A string of length 64 with a 0-7 in each position | [optional] [default to "0000000011111111222222223333333344444444555555556666666677777777"]
**AnycastMacAddress** | Pointer to **string** | Site Level MAC Address for Anycast | [optional] [default to "(auto)"]
**AnycastMacAddressAutoAssigned** | Pointer to **bool** | Whether or not the value in anycast_mac_address field has been automatically assigned or not. Set to false and change anycast_mac_address value to edit. | [optional] 
**MacAddressAgingTime** | Pointer to **int32** | MAC Address Aging Time (between 1-100000) | [optional] [default to 600]
**MlagDelayRestoreTimer** | Pointer to **int32** | MLAG Delay Restore Timer | [optional] [default to 300]
**BgpKeepaliveTimer** | Pointer to **int32** | Spine BGP Keepalive Timer | [optional] [default to 60]
**BgpHoldDownTimer** | Pointer to **int32** | Spine BGP Hold Down Timer | [optional] [default to 180]
**SpineBgpAdvertisementInterval** | Pointer to **int32** | BGP Advertisement Interval for spines/superspines. Use \&quot;0\&quot; for immediate updates | [optional] [default to 1]
**SpineBgpConnectTimer** | Pointer to **int32** | BGP Connect Timer | [optional] [default to 120]
**LeafBgpKeepAliveTimer** | Pointer to **int32** | Leaf BGP Keep Alive Timer | [optional] [default to 60]
**LeafBgpHoldDownTimer** | Pointer to **int32** | Leaf BGP Hold Down Timer | [optional] [default to 180]
**LeafBgpAdvertisementInterval** | Pointer to **int32** | BGP Advertisement Interval for leafs. Use \&quot;0\&quot; for immediate updates | [optional] [default to 1]
**LeafBgpConnectTimer** | Pointer to **int32** | BGP Connect Timer | [optional] [default to 120]
**LinkStateTimeoutValue** | Pointer to **NullableInt32** | Link State Timeout Value | [optional] [default to 60]
**EvpnMultihomingStartupDelay** | Pointer to **NullableInt32** | Startup Delay | [optional] [default to 300]
**EvpnMacHoldtime** | Pointer to **NullableInt32** | MAC Holdtime | [optional] [default to 1080]
**AggressiveReporting** | Pointer to **bool** | Fast Reporting of Switch Communications, Link Up/Down, and BGP Status | [optional] [default to true]
**CrcFailureThreshold** | Pointer to **NullableInt32** | Threshold in Errors per second that when met will disable the links as part of LAGs | [optional] [default to 5]
**Islands** | Pointer to [**[]SitesPatchRequestSiteValueIslandsInner**](SitesPatchRequestSiteValueIslandsInner.md) |  | [optional] 
**Pairs** | Pointer to [**[]SitesPatchRequestSiteValuePairsInner**](SitesPatchRequestSiteValuePairsInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**SitesPatchRequestSiteValueObjectProperties**](SitesPatchRequestSiteValueObjectProperties.md) |  | [optional] 
**EnableDhcpSnooping** | Pointer to **bool** | Enables the switches to monitor DHCP traffic and collect assigned IP addresses which are then placed in the DHCP assigned IPs report. | [optional] [default to false]

## Methods

### NewSitesPatchRequestSiteValue

`func NewSitesPatchRequestSiteValue() *SitesPatchRequestSiteValue`

NewSitesPatchRequestSiteValue instantiates a new SitesPatchRequestSiteValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSitesPatchRequestSiteValueWithDefaults

`func NewSitesPatchRequestSiteValueWithDefaults() *SitesPatchRequestSiteValue`

NewSitesPatchRequestSiteValueWithDefaults instantiates a new SitesPatchRequestSiteValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *SitesPatchRequestSiteValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *SitesPatchRequestSiteValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *SitesPatchRequestSiteValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *SitesPatchRequestSiteValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *SitesPatchRequestSiteValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *SitesPatchRequestSiteValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *SitesPatchRequestSiteValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *SitesPatchRequestSiteValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetServiceForSite

`func (o *SitesPatchRequestSiteValue) GetServiceForSite() string`

GetServiceForSite returns the ServiceForSite field if non-nil, zero value otherwise.

### GetServiceForSiteOk

`func (o *SitesPatchRequestSiteValue) GetServiceForSiteOk() (*string, bool)`

GetServiceForSiteOk returns a tuple with the ServiceForSite field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServiceForSite

`func (o *SitesPatchRequestSiteValue) SetServiceForSite(v string)`

SetServiceForSite sets ServiceForSite field to given value.

### HasServiceForSite

`func (o *SitesPatchRequestSiteValue) HasServiceForSite() bool`

HasServiceForSite returns a boolean if a field has been set.

### GetServiceForSiteRefType

`func (o *SitesPatchRequestSiteValue) GetServiceForSiteRefType() string`

GetServiceForSiteRefType returns the ServiceForSiteRefType field if non-nil, zero value otherwise.

### GetServiceForSiteRefTypeOk

`func (o *SitesPatchRequestSiteValue) GetServiceForSiteRefTypeOk() (*string, bool)`

GetServiceForSiteRefTypeOk returns a tuple with the ServiceForSiteRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServiceForSiteRefType

`func (o *SitesPatchRequestSiteValue) SetServiceForSiteRefType(v string)`

SetServiceForSiteRefType sets ServiceForSiteRefType field to given value.

### HasServiceForSiteRefType

`func (o *SitesPatchRequestSiteValue) HasServiceForSiteRefType() bool`

HasServiceForSiteRefType returns a boolean if a field has been set.

### GetSpanningTreeType

`func (o *SitesPatchRequestSiteValue) GetSpanningTreeType() string`

GetSpanningTreeType returns the SpanningTreeType field if non-nil, zero value otherwise.

### GetSpanningTreeTypeOk

`func (o *SitesPatchRequestSiteValue) GetSpanningTreeTypeOk() (*string, bool)`

GetSpanningTreeTypeOk returns a tuple with the SpanningTreeType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpanningTreeType

`func (o *SitesPatchRequestSiteValue) SetSpanningTreeType(v string)`

SetSpanningTreeType sets SpanningTreeType field to given value.

### HasSpanningTreeType

`func (o *SitesPatchRequestSiteValue) HasSpanningTreeType() bool`

HasSpanningTreeType returns a boolean if a field has been set.

### GetRegionName

`func (o *SitesPatchRequestSiteValue) GetRegionName() string`

GetRegionName returns the RegionName field if non-nil, zero value otherwise.

### GetRegionNameOk

`func (o *SitesPatchRequestSiteValue) GetRegionNameOk() (*string, bool)`

GetRegionNameOk returns a tuple with the RegionName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegionName

`func (o *SitesPatchRequestSiteValue) SetRegionName(v string)`

SetRegionName sets RegionName field to given value.

### HasRegionName

`func (o *SitesPatchRequestSiteValue) HasRegionName() bool`

HasRegionName returns a boolean if a field has been set.

### GetRevision

`func (o *SitesPatchRequestSiteValue) GetRevision() int32`

GetRevision returns the Revision field if non-nil, zero value otherwise.

### GetRevisionOk

`func (o *SitesPatchRequestSiteValue) GetRevisionOk() (*int32, bool)`

GetRevisionOk returns a tuple with the Revision field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRevision

`func (o *SitesPatchRequestSiteValue) SetRevision(v int32)`

SetRevision sets Revision field to given value.

### HasRevision

`func (o *SitesPatchRequestSiteValue) HasRevision() bool`

HasRevision returns a boolean if a field has been set.

### SetRevisionNil

`func (o *SitesPatchRequestSiteValue) SetRevisionNil(b bool)`

 SetRevisionNil sets the value for Revision to be an explicit nil

### UnsetRevision
`func (o *SitesPatchRequestSiteValue) UnsetRevision()`

UnsetRevision ensures that no value is present for Revision, not even an explicit nil
### GetForceSpanningTreeOnFabricPorts

`func (o *SitesPatchRequestSiteValue) GetForceSpanningTreeOnFabricPorts() bool`

GetForceSpanningTreeOnFabricPorts returns the ForceSpanningTreeOnFabricPorts field if non-nil, zero value otherwise.

### GetForceSpanningTreeOnFabricPortsOk

`func (o *SitesPatchRequestSiteValue) GetForceSpanningTreeOnFabricPortsOk() (*bool, bool)`

GetForceSpanningTreeOnFabricPortsOk returns a tuple with the ForceSpanningTreeOnFabricPorts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetForceSpanningTreeOnFabricPorts

`func (o *SitesPatchRequestSiteValue) SetForceSpanningTreeOnFabricPorts(v bool)`

SetForceSpanningTreeOnFabricPorts sets ForceSpanningTreeOnFabricPorts field to given value.

### HasForceSpanningTreeOnFabricPorts

`func (o *SitesPatchRequestSiteValue) HasForceSpanningTreeOnFabricPorts() bool`

HasForceSpanningTreeOnFabricPorts returns a boolean if a field has been set.

### GetReadOnlyMode

`func (o *SitesPatchRequestSiteValue) GetReadOnlyMode() bool`

GetReadOnlyMode returns the ReadOnlyMode field if non-nil, zero value otherwise.

### GetReadOnlyModeOk

`func (o *SitesPatchRequestSiteValue) GetReadOnlyModeOk() (*bool, bool)`

GetReadOnlyModeOk returns a tuple with the ReadOnlyMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReadOnlyMode

`func (o *SitesPatchRequestSiteValue) SetReadOnlyMode(v bool)`

SetReadOnlyMode sets ReadOnlyMode field to given value.

### HasReadOnlyMode

`func (o *SitesPatchRequestSiteValue) HasReadOnlyMode() bool`

HasReadOnlyMode returns a boolean if a field has been set.

### GetDscpToPBitMap

`func (o *SitesPatchRequestSiteValue) GetDscpToPBitMap() string`

GetDscpToPBitMap returns the DscpToPBitMap field if non-nil, zero value otherwise.

### GetDscpToPBitMapOk

`func (o *SitesPatchRequestSiteValue) GetDscpToPBitMapOk() (*string, bool)`

GetDscpToPBitMapOk returns a tuple with the DscpToPBitMap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDscpToPBitMap

`func (o *SitesPatchRequestSiteValue) SetDscpToPBitMap(v string)`

SetDscpToPBitMap sets DscpToPBitMap field to given value.

### HasDscpToPBitMap

`func (o *SitesPatchRequestSiteValue) HasDscpToPBitMap() bool`

HasDscpToPBitMap returns a boolean if a field has been set.

### GetAnycastMacAddress

`func (o *SitesPatchRequestSiteValue) GetAnycastMacAddress() string`

GetAnycastMacAddress returns the AnycastMacAddress field if non-nil, zero value otherwise.

### GetAnycastMacAddressOk

`func (o *SitesPatchRequestSiteValue) GetAnycastMacAddressOk() (*string, bool)`

GetAnycastMacAddressOk returns a tuple with the AnycastMacAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnycastMacAddress

`func (o *SitesPatchRequestSiteValue) SetAnycastMacAddress(v string)`

SetAnycastMacAddress sets AnycastMacAddress field to given value.

### HasAnycastMacAddress

`func (o *SitesPatchRequestSiteValue) HasAnycastMacAddress() bool`

HasAnycastMacAddress returns a boolean if a field has been set.

### GetAnycastMacAddressAutoAssigned

`func (o *SitesPatchRequestSiteValue) GetAnycastMacAddressAutoAssigned() bool`

GetAnycastMacAddressAutoAssigned returns the AnycastMacAddressAutoAssigned field if non-nil, zero value otherwise.

### GetAnycastMacAddressAutoAssignedOk

`func (o *SitesPatchRequestSiteValue) GetAnycastMacAddressAutoAssignedOk() (*bool, bool)`

GetAnycastMacAddressAutoAssignedOk returns a tuple with the AnycastMacAddressAutoAssigned field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnycastMacAddressAutoAssigned

`func (o *SitesPatchRequestSiteValue) SetAnycastMacAddressAutoAssigned(v bool)`

SetAnycastMacAddressAutoAssigned sets AnycastMacAddressAutoAssigned field to given value.

### HasAnycastMacAddressAutoAssigned

`func (o *SitesPatchRequestSiteValue) HasAnycastMacAddressAutoAssigned() bool`

HasAnycastMacAddressAutoAssigned returns a boolean if a field has been set.

### GetMacAddressAgingTime

`func (o *SitesPatchRequestSiteValue) GetMacAddressAgingTime() int32`

GetMacAddressAgingTime returns the MacAddressAgingTime field if non-nil, zero value otherwise.

### GetMacAddressAgingTimeOk

`func (o *SitesPatchRequestSiteValue) GetMacAddressAgingTimeOk() (*int32, bool)`

GetMacAddressAgingTimeOk returns a tuple with the MacAddressAgingTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMacAddressAgingTime

`func (o *SitesPatchRequestSiteValue) SetMacAddressAgingTime(v int32)`

SetMacAddressAgingTime sets MacAddressAgingTime field to given value.

### HasMacAddressAgingTime

`func (o *SitesPatchRequestSiteValue) HasMacAddressAgingTime() bool`

HasMacAddressAgingTime returns a boolean if a field has been set.

### GetMlagDelayRestoreTimer

`func (o *SitesPatchRequestSiteValue) GetMlagDelayRestoreTimer() int32`

GetMlagDelayRestoreTimer returns the MlagDelayRestoreTimer field if non-nil, zero value otherwise.

### GetMlagDelayRestoreTimerOk

`func (o *SitesPatchRequestSiteValue) GetMlagDelayRestoreTimerOk() (*int32, bool)`

GetMlagDelayRestoreTimerOk returns a tuple with the MlagDelayRestoreTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMlagDelayRestoreTimer

`func (o *SitesPatchRequestSiteValue) SetMlagDelayRestoreTimer(v int32)`

SetMlagDelayRestoreTimer sets MlagDelayRestoreTimer field to given value.

### HasMlagDelayRestoreTimer

`func (o *SitesPatchRequestSiteValue) HasMlagDelayRestoreTimer() bool`

HasMlagDelayRestoreTimer returns a boolean if a field has been set.

### GetBgpKeepaliveTimer

`func (o *SitesPatchRequestSiteValue) GetBgpKeepaliveTimer() int32`

GetBgpKeepaliveTimer returns the BgpKeepaliveTimer field if non-nil, zero value otherwise.

### GetBgpKeepaliveTimerOk

`func (o *SitesPatchRequestSiteValue) GetBgpKeepaliveTimerOk() (*int32, bool)`

GetBgpKeepaliveTimerOk returns a tuple with the BgpKeepaliveTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBgpKeepaliveTimer

`func (o *SitesPatchRequestSiteValue) SetBgpKeepaliveTimer(v int32)`

SetBgpKeepaliveTimer sets BgpKeepaliveTimer field to given value.

### HasBgpKeepaliveTimer

`func (o *SitesPatchRequestSiteValue) HasBgpKeepaliveTimer() bool`

HasBgpKeepaliveTimer returns a boolean if a field has been set.

### GetBgpHoldDownTimer

`func (o *SitesPatchRequestSiteValue) GetBgpHoldDownTimer() int32`

GetBgpHoldDownTimer returns the BgpHoldDownTimer field if non-nil, zero value otherwise.

### GetBgpHoldDownTimerOk

`func (o *SitesPatchRequestSiteValue) GetBgpHoldDownTimerOk() (*int32, bool)`

GetBgpHoldDownTimerOk returns a tuple with the BgpHoldDownTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBgpHoldDownTimer

`func (o *SitesPatchRequestSiteValue) SetBgpHoldDownTimer(v int32)`

SetBgpHoldDownTimer sets BgpHoldDownTimer field to given value.

### HasBgpHoldDownTimer

`func (o *SitesPatchRequestSiteValue) HasBgpHoldDownTimer() bool`

HasBgpHoldDownTimer returns a boolean if a field has been set.

### GetSpineBgpAdvertisementInterval

`func (o *SitesPatchRequestSiteValue) GetSpineBgpAdvertisementInterval() int32`

GetSpineBgpAdvertisementInterval returns the SpineBgpAdvertisementInterval field if non-nil, zero value otherwise.

### GetSpineBgpAdvertisementIntervalOk

`func (o *SitesPatchRequestSiteValue) GetSpineBgpAdvertisementIntervalOk() (*int32, bool)`

GetSpineBgpAdvertisementIntervalOk returns a tuple with the SpineBgpAdvertisementInterval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpineBgpAdvertisementInterval

`func (o *SitesPatchRequestSiteValue) SetSpineBgpAdvertisementInterval(v int32)`

SetSpineBgpAdvertisementInterval sets SpineBgpAdvertisementInterval field to given value.

### HasSpineBgpAdvertisementInterval

`func (o *SitesPatchRequestSiteValue) HasSpineBgpAdvertisementInterval() bool`

HasSpineBgpAdvertisementInterval returns a boolean if a field has been set.

### GetSpineBgpConnectTimer

`func (o *SitesPatchRequestSiteValue) GetSpineBgpConnectTimer() int32`

GetSpineBgpConnectTimer returns the SpineBgpConnectTimer field if non-nil, zero value otherwise.

### GetSpineBgpConnectTimerOk

`func (o *SitesPatchRequestSiteValue) GetSpineBgpConnectTimerOk() (*int32, bool)`

GetSpineBgpConnectTimerOk returns a tuple with the SpineBgpConnectTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpineBgpConnectTimer

`func (o *SitesPatchRequestSiteValue) SetSpineBgpConnectTimer(v int32)`

SetSpineBgpConnectTimer sets SpineBgpConnectTimer field to given value.

### HasSpineBgpConnectTimer

`func (o *SitesPatchRequestSiteValue) HasSpineBgpConnectTimer() bool`

HasSpineBgpConnectTimer returns a boolean if a field has been set.

### GetLeafBgpKeepAliveTimer

`func (o *SitesPatchRequestSiteValue) GetLeafBgpKeepAliveTimer() int32`

GetLeafBgpKeepAliveTimer returns the LeafBgpKeepAliveTimer field if non-nil, zero value otherwise.

### GetLeafBgpKeepAliveTimerOk

`func (o *SitesPatchRequestSiteValue) GetLeafBgpKeepAliveTimerOk() (*int32, bool)`

GetLeafBgpKeepAliveTimerOk returns a tuple with the LeafBgpKeepAliveTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLeafBgpKeepAliveTimer

`func (o *SitesPatchRequestSiteValue) SetLeafBgpKeepAliveTimer(v int32)`

SetLeafBgpKeepAliveTimer sets LeafBgpKeepAliveTimer field to given value.

### HasLeafBgpKeepAliveTimer

`func (o *SitesPatchRequestSiteValue) HasLeafBgpKeepAliveTimer() bool`

HasLeafBgpKeepAliveTimer returns a boolean if a field has been set.

### GetLeafBgpHoldDownTimer

`func (o *SitesPatchRequestSiteValue) GetLeafBgpHoldDownTimer() int32`

GetLeafBgpHoldDownTimer returns the LeafBgpHoldDownTimer field if non-nil, zero value otherwise.

### GetLeafBgpHoldDownTimerOk

`func (o *SitesPatchRequestSiteValue) GetLeafBgpHoldDownTimerOk() (*int32, bool)`

GetLeafBgpHoldDownTimerOk returns a tuple with the LeafBgpHoldDownTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLeafBgpHoldDownTimer

`func (o *SitesPatchRequestSiteValue) SetLeafBgpHoldDownTimer(v int32)`

SetLeafBgpHoldDownTimer sets LeafBgpHoldDownTimer field to given value.

### HasLeafBgpHoldDownTimer

`func (o *SitesPatchRequestSiteValue) HasLeafBgpHoldDownTimer() bool`

HasLeafBgpHoldDownTimer returns a boolean if a field has been set.

### GetLeafBgpAdvertisementInterval

`func (o *SitesPatchRequestSiteValue) GetLeafBgpAdvertisementInterval() int32`

GetLeafBgpAdvertisementInterval returns the LeafBgpAdvertisementInterval field if non-nil, zero value otherwise.

### GetLeafBgpAdvertisementIntervalOk

`func (o *SitesPatchRequestSiteValue) GetLeafBgpAdvertisementIntervalOk() (*int32, bool)`

GetLeafBgpAdvertisementIntervalOk returns a tuple with the LeafBgpAdvertisementInterval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLeafBgpAdvertisementInterval

`func (o *SitesPatchRequestSiteValue) SetLeafBgpAdvertisementInterval(v int32)`

SetLeafBgpAdvertisementInterval sets LeafBgpAdvertisementInterval field to given value.

### HasLeafBgpAdvertisementInterval

`func (o *SitesPatchRequestSiteValue) HasLeafBgpAdvertisementInterval() bool`

HasLeafBgpAdvertisementInterval returns a boolean if a field has been set.

### GetLeafBgpConnectTimer

`func (o *SitesPatchRequestSiteValue) GetLeafBgpConnectTimer() int32`

GetLeafBgpConnectTimer returns the LeafBgpConnectTimer field if non-nil, zero value otherwise.

### GetLeafBgpConnectTimerOk

`func (o *SitesPatchRequestSiteValue) GetLeafBgpConnectTimerOk() (*int32, bool)`

GetLeafBgpConnectTimerOk returns a tuple with the LeafBgpConnectTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLeafBgpConnectTimer

`func (o *SitesPatchRequestSiteValue) SetLeafBgpConnectTimer(v int32)`

SetLeafBgpConnectTimer sets LeafBgpConnectTimer field to given value.

### HasLeafBgpConnectTimer

`func (o *SitesPatchRequestSiteValue) HasLeafBgpConnectTimer() bool`

HasLeafBgpConnectTimer returns a boolean if a field has been set.

### GetLinkStateTimeoutValue

`func (o *SitesPatchRequestSiteValue) GetLinkStateTimeoutValue() int32`

GetLinkStateTimeoutValue returns the LinkStateTimeoutValue field if non-nil, zero value otherwise.

### GetLinkStateTimeoutValueOk

`func (o *SitesPatchRequestSiteValue) GetLinkStateTimeoutValueOk() (*int32, bool)`

GetLinkStateTimeoutValueOk returns a tuple with the LinkStateTimeoutValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLinkStateTimeoutValue

`func (o *SitesPatchRequestSiteValue) SetLinkStateTimeoutValue(v int32)`

SetLinkStateTimeoutValue sets LinkStateTimeoutValue field to given value.

### HasLinkStateTimeoutValue

`func (o *SitesPatchRequestSiteValue) HasLinkStateTimeoutValue() bool`

HasLinkStateTimeoutValue returns a boolean if a field has been set.

### SetLinkStateTimeoutValueNil

`func (o *SitesPatchRequestSiteValue) SetLinkStateTimeoutValueNil(b bool)`

 SetLinkStateTimeoutValueNil sets the value for LinkStateTimeoutValue to be an explicit nil

### UnsetLinkStateTimeoutValue
`func (o *SitesPatchRequestSiteValue) UnsetLinkStateTimeoutValue()`

UnsetLinkStateTimeoutValue ensures that no value is present for LinkStateTimeoutValue, not even an explicit nil
### GetEvpnMultihomingStartupDelay

`func (o *SitesPatchRequestSiteValue) GetEvpnMultihomingStartupDelay() int32`

GetEvpnMultihomingStartupDelay returns the EvpnMultihomingStartupDelay field if non-nil, zero value otherwise.

### GetEvpnMultihomingStartupDelayOk

`func (o *SitesPatchRequestSiteValue) GetEvpnMultihomingStartupDelayOk() (*int32, bool)`

GetEvpnMultihomingStartupDelayOk returns a tuple with the EvpnMultihomingStartupDelay field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEvpnMultihomingStartupDelay

`func (o *SitesPatchRequestSiteValue) SetEvpnMultihomingStartupDelay(v int32)`

SetEvpnMultihomingStartupDelay sets EvpnMultihomingStartupDelay field to given value.

### HasEvpnMultihomingStartupDelay

`func (o *SitesPatchRequestSiteValue) HasEvpnMultihomingStartupDelay() bool`

HasEvpnMultihomingStartupDelay returns a boolean if a field has been set.

### SetEvpnMultihomingStartupDelayNil

`func (o *SitesPatchRequestSiteValue) SetEvpnMultihomingStartupDelayNil(b bool)`

 SetEvpnMultihomingStartupDelayNil sets the value for EvpnMultihomingStartupDelay to be an explicit nil

### UnsetEvpnMultihomingStartupDelay
`func (o *SitesPatchRequestSiteValue) UnsetEvpnMultihomingStartupDelay()`

UnsetEvpnMultihomingStartupDelay ensures that no value is present for EvpnMultihomingStartupDelay, not even an explicit nil
### GetEvpnMacHoldtime

`func (o *SitesPatchRequestSiteValue) GetEvpnMacHoldtime() int32`

GetEvpnMacHoldtime returns the EvpnMacHoldtime field if non-nil, zero value otherwise.

### GetEvpnMacHoldtimeOk

`func (o *SitesPatchRequestSiteValue) GetEvpnMacHoldtimeOk() (*int32, bool)`

GetEvpnMacHoldtimeOk returns a tuple with the EvpnMacHoldtime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEvpnMacHoldtime

`func (o *SitesPatchRequestSiteValue) SetEvpnMacHoldtime(v int32)`

SetEvpnMacHoldtime sets EvpnMacHoldtime field to given value.

### HasEvpnMacHoldtime

`func (o *SitesPatchRequestSiteValue) HasEvpnMacHoldtime() bool`

HasEvpnMacHoldtime returns a boolean if a field has been set.

### SetEvpnMacHoldtimeNil

`func (o *SitesPatchRequestSiteValue) SetEvpnMacHoldtimeNil(b bool)`

 SetEvpnMacHoldtimeNil sets the value for EvpnMacHoldtime to be an explicit nil

### UnsetEvpnMacHoldtime
`func (o *SitesPatchRequestSiteValue) UnsetEvpnMacHoldtime()`

UnsetEvpnMacHoldtime ensures that no value is present for EvpnMacHoldtime, not even an explicit nil
### GetAggressiveReporting

`func (o *SitesPatchRequestSiteValue) GetAggressiveReporting() bool`

GetAggressiveReporting returns the AggressiveReporting field if non-nil, zero value otherwise.

### GetAggressiveReportingOk

`func (o *SitesPatchRequestSiteValue) GetAggressiveReportingOk() (*bool, bool)`

GetAggressiveReportingOk returns a tuple with the AggressiveReporting field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAggressiveReporting

`func (o *SitesPatchRequestSiteValue) SetAggressiveReporting(v bool)`

SetAggressiveReporting sets AggressiveReporting field to given value.

### HasAggressiveReporting

`func (o *SitesPatchRequestSiteValue) HasAggressiveReporting() bool`

HasAggressiveReporting returns a boolean if a field has been set.

### GetCrcFailureThreshold

`func (o *SitesPatchRequestSiteValue) GetCrcFailureThreshold() int32`

GetCrcFailureThreshold returns the CrcFailureThreshold field if non-nil, zero value otherwise.

### GetCrcFailureThresholdOk

`func (o *SitesPatchRequestSiteValue) GetCrcFailureThresholdOk() (*int32, bool)`

GetCrcFailureThresholdOk returns a tuple with the CrcFailureThreshold field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCrcFailureThreshold

`func (o *SitesPatchRequestSiteValue) SetCrcFailureThreshold(v int32)`

SetCrcFailureThreshold sets CrcFailureThreshold field to given value.

### HasCrcFailureThreshold

`func (o *SitesPatchRequestSiteValue) HasCrcFailureThreshold() bool`

HasCrcFailureThreshold returns a boolean if a field has been set.

### SetCrcFailureThresholdNil

`func (o *SitesPatchRequestSiteValue) SetCrcFailureThresholdNil(b bool)`

 SetCrcFailureThresholdNil sets the value for CrcFailureThreshold to be an explicit nil

### UnsetCrcFailureThreshold
`func (o *SitesPatchRequestSiteValue) UnsetCrcFailureThreshold()`

UnsetCrcFailureThreshold ensures that no value is present for CrcFailureThreshold, not even an explicit nil
### GetIslands

`func (o *SitesPatchRequestSiteValue) GetIslands() []SitesPatchRequestSiteValueIslandsInner`

GetIslands returns the Islands field if non-nil, zero value otherwise.

### GetIslandsOk

`func (o *SitesPatchRequestSiteValue) GetIslandsOk() (*[]SitesPatchRequestSiteValueIslandsInner, bool)`

GetIslandsOk returns a tuple with the Islands field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIslands

`func (o *SitesPatchRequestSiteValue) SetIslands(v []SitesPatchRequestSiteValueIslandsInner)`

SetIslands sets Islands field to given value.

### HasIslands

`func (o *SitesPatchRequestSiteValue) HasIslands() bool`

HasIslands returns a boolean if a field has been set.

### GetPairs

`func (o *SitesPatchRequestSiteValue) GetPairs() []SitesPatchRequestSiteValuePairsInner`

GetPairs returns the Pairs field if non-nil, zero value otherwise.

### GetPairsOk

`func (o *SitesPatchRequestSiteValue) GetPairsOk() (*[]SitesPatchRequestSiteValuePairsInner, bool)`

GetPairsOk returns a tuple with the Pairs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPairs

`func (o *SitesPatchRequestSiteValue) SetPairs(v []SitesPatchRequestSiteValuePairsInner)`

SetPairs sets Pairs field to given value.

### HasPairs

`func (o *SitesPatchRequestSiteValue) HasPairs() bool`

HasPairs returns a boolean if a field has been set.

### GetObjectProperties

`func (o *SitesPatchRequestSiteValue) GetObjectProperties() SitesPatchRequestSiteValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *SitesPatchRequestSiteValue) GetObjectPropertiesOk() (*SitesPatchRequestSiteValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *SitesPatchRequestSiteValue) SetObjectProperties(v SitesPatchRequestSiteValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *SitesPatchRequestSiteValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.

### GetEnableDhcpSnooping

`func (o *SitesPatchRequestSiteValue) GetEnableDhcpSnooping() bool`

GetEnableDhcpSnooping returns the EnableDhcpSnooping field if non-nil, zero value otherwise.

### GetEnableDhcpSnoopingOk

`func (o *SitesPatchRequestSiteValue) GetEnableDhcpSnoopingOk() (*bool, bool)`

GetEnableDhcpSnoopingOk returns a tuple with the EnableDhcpSnooping field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnableDhcpSnooping

`func (o *SitesPatchRequestSiteValue) SetEnableDhcpSnooping(v bool)`

SetEnableDhcpSnooping sets EnableDhcpSnooping field to given value.

### HasEnableDhcpSnooping

`func (o *SitesPatchRequestSiteValue) HasEnableDhcpSnooping() bool`

HasEnableDhcpSnooping returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


