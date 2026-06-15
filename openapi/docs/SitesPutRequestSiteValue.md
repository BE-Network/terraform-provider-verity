# SitesPutRequestSiteValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Template Name. Must be unique within type. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to true]
**SuSupport** | Pointer to **bool** | Support grouping leaf switches in SUs | [optional] [default to false]
**AllowAllUnderlayConnections** | Pointer to **bool** | Allows underlay connections between PODs | [optional] [default to false]
**SiteType** | Pointer to **string** | Type of Fabric | [optional] [default to "enterprise"]
**DuplicateAddressDetectionMaxNumberOfMoves** | Pointer to **NullableInt32** | Controls duplicate MAC address detection (DAD) Max Number of Moves for EVPN (Ethernet VPN) within the BGP address-family. Number of moves (2 to 1000; default 5 if left blank) | [optional] [default to 5]
**DuplicateAddressDetectionTime** | Pointer to **NullableInt32** | Controls duplicate MAC address detection (DAD) time for EVPN (Ethernet VPN) within the BGP address-family. Time in seconds (2 to 1800; default 180 if left blank) | [optional] [default to 180]
**PortAdminPollingInterval** | Pointer to **NullableInt32** | polling interval values in seconds, set if aggressive reporting is not enabled | [optional] [default to 0]
**PortStatusPollingInterval** | Pointer to **NullableInt32** | polling interval values in seconds, set if aggressive reporting is not enabled | [optional] [default to 0]
**ServiceForSite** | Pointer to **string** | Service for Fabric | [optional] [default to "(predefined):Management"]
**ServiceForSiteRefType** | Pointer to **string** | Object type for service_for_site field | [optional] 
**SpanningTreeType** | Pointer to **string** | Sets the spanning tree type for all Ports in this Fabric with Spanning Tree enabled | [optional] [default to "pvst"]
**RegionName** | Pointer to **string** | Defines the logical boundary of the network. All switches in an MSTP region must have the same configured region name | [optional] [default to ""]
**Revision** | Pointer to **NullableInt32** | A logical number that signifies a revision for the MSTP configuration. All switches in an MSTP region must have the same revision number | [optional] [default to 0]
**ForceSpanningTreeOnFabricPorts** | Pointer to **bool** | Enable spanning tree on all fabric connections.  This overrides the Eth Port Settings for Fabric ports | [optional] [default to false]
**ReadOnlyMode** | Pointer to **bool** | When Read Only Mode is checked, vNetC will perform all functions except writing database updates to the target hardware | [optional] [default to false]
**DomainForSite** | Pointer to **string** | Fabric Collection for Fabric | [optional] [default to ""]
**DomainForSiteRefType** | Pointer to **string** | Object type for domain_for_site field | [optional] 
**EnableDscp** | Pointer to **bool** | Enable DSCP to p-bit/TC configuration. When enabled, DSCP to p-bit/TC mappings are applied. | [optional] [default to true]
**DscpToPBitMap** | Pointer to **string** | For any Service that is using DSCP to p-bit map packet prioritization. A string of length 64 with a 0-7 in each position | [optional] [default to "0000000011111111222222223333333344444444555555556666666677777777"]
**AnycastMacAddress** | Pointer to **string** | Fabric Level MAC Address for Anycast | [optional] [default to "(auto)"]
**AnycastMacAddressAutoAssigned** | Pointer to **bool** | Whether or not the value in anycast_mac_address field has been automatically assigned or not. Set to false and change anycast_mac_address value to edit. | [optional] 
**MacAddressAgingTime** | Pointer to **NullableInt32** | MAC Address Aging Time (between 1-100000) | [optional] [default to 600]
**MlagDelayRestoreTimer** | Pointer to **NullableInt32** | MLAG Delay Restore Timer | [optional] [default to 300]
**BgpKeepaliveTimer** | Pointer to **NullableInt32** | Spine BGP Keepalive Timer | [optional] [default to 60]
**BgpHoldDownTimer** | Pointer to **NullableInt32** | Spine BGP Hold Down Timer | [optional] [default to 180]
**SpineBgpAdvertisementInterval** | Pointer to **NullableInt32** | BGP Advertisement Interval for spines/superspines. Use \&quot;0\&quot; for immediate updates | [optional] [default to 1]
**SpineBgpConnectTimer** | Pointer to **NullableInt32** | BGP Connect Timer | [optional] [default to 120]
**SpineAsNumber** | Pointer to **NullableInt32** | BGP AS number applied uniformly to all spine endpoints in this CLOS fabric on save. Leave blank to manage spine AS numbers individually. | [optional] 
**LeafBgpKeepAliveTimer** | Pointer to **NullableInt32** | Leaf BGP Keep Alive Timer | [optional] [default to 60]
**LeafBgpHoldDownTimer** | Pointer to **NullableInt32** | Leaf BGP Hold Down Timer | [optional] [default to 180]
**LeafBgpAdvertisementInterval** | Pointer to **NullableInt32** | BGP Advertisement Interval for leafs. Use \&quot;0\&quot; for immediate updates | [optional] [default to 1]
**LeafBgpConnectTimer** | Pointer to **NullableInt32** | BGP Connect Timer | [optional] [default to 120]
**LinkStateTimeoutValue** | Pointer to **NullableInt32** | Link State Timeout Value | [optional] [default to 60]
**EvpnMultihomingStartupDelay** | Pointer to **NullableInt32** | Startup Delay | [optional] [default to 300]
**EvpnMacHoldtime** | Pointer to **NullableInt32** | MAC Holdtime | [optional] [default to 1080]
**AggressiveReporting** | Pointer to **bool** | Fast Reporting of Switch Communications, Link Up/Down, and BGP Status | [optional] [default to true]
**SwitchIpBase** | Pointer to **string** | Base IPv4 address for switch IPs in this Fabric | [optional] [default to ""]
**ControllerIpBase** | Pointer to **string** | Base IPv4 address for controller IPs in this Fabric | [optional] [default to ""]
**MultiTenant** | Pointer to **bool** | Allow multiple tenants to HGX endpoints on this fabric. | [optional] [default to true]
**BaseBgpAsNumber** | Pointer to **string** | Base BGP Autonomous System Number used for switches in the fabric  | [optional] [default to "61000"]
**RouterIdBasePrefix** | Pointer to **string** | Router ID starting IP address  | [optional] [default to "172.16.0.0"]
**VtepIdBasePrefix** | Pointer to **string** | Vtep ID starting IP address  | [optional] [default to "172.16.10.0"]
**PairedIpSubnet** | Pointer to **string** | IP address range reserved for communication between paired switches  | [optional] [default to "192.168.254.0/24"]
**MaxSwitches** | Pointer to **string** | Max number Switches to support in this site  | [optional] [default to "2000"]
**PauseValidationAlarms** | Pointer to **bool** | Validation still runs, but validation alarms are not raised for this Fabric while enabled. | [optional] [default to false]
**StartingOctet** | Pointer to **NullableInt32** | Starting Octet for HGX Port IPs | [optional] 
**MaxSus** | Pointer to **NullableInt32** | Maximum number of SUs allowed per POD | [optional] 
**MaxPods** | Pointer to **NullableInt32** | Maximum number of PODs allowed in the Fabric | [optional] 
**ObjectProperties** | Pointer to [**SitesPutRequestSiteValueObjectProperties**](SitesPutRequestSiteValueObjectProperties.md) |  | [optional] 
**SwitchUsername** | Pointer to **string** | Default username for managed switches in this Fabric | [optional] [default to ""]
**SwitchPassword** | Pointer to **string** | Default password for managed switches in this Fabric | [optional] [default to ""]
**SwitchPasswordEncrypted** | Pointer to **string** | Default password for managed switches in this Fabric | [optional] [default to ""]
**HgxUsername** | Pointer to **string** | Default username for HGX devices in this Fabric | [optional] [default to ""]
**HgxPassword** | Pointer to **string** | Default password for HGX devices in this Fabric | [optional] [default to ""]
**HgxPasswordEncrypted** | Pointer to **string** | Default password for HGX devices in this Fabric | [optional] [default to ""]
**SwitchGateway** | Pointer to **string** | Default switch management gateway IP for devices in this Fabric | [optional] [default to ""]
**ControllerGateway** | Pointer to **string** | Default controller gateway IP for devices in this Fabric | [optional] [default to ""]
**HgxGateway** | Pointer to **string** | Default HGX management gateway IP for devices in this Fabric | [optional] [default to ""]
**IpSourceGuard** | Pointer to **bool** | On untrusted ports, only allow known traffic from known IP addresses. IP addresses are discovered via DHCP snooping or with static IP settings | [optional] [default to false]
**EnableDhcpSnooping** | Pointer to **bool** | Enables the switches to monitor DHCP traffic and collect assigned IP addresses which are then placed in the DHCP assigned IPs report. | [optional] [default to false]

## Methods

### NewSitesPutRequestSiteValue

`func NewSitesPutRequestSiteValue() *SitesPutRequestSiteValue`

NewSitesPutRequestSiteValue instantiates a new SitesPutRequestSiteValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSitesPutRequestSiteValueWithDefaults

`func NewSitesPutRequestSiteValueWithDefaults() *SitesPutRequestSiteValue`

NewSitesPutRequestSiteValueWithDefaults instantiates a new SitesPutRequestSiteValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *SitesPutRequestSiteValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *SitesPutRequestSiteValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *SitesPutRequestSiteValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *SitesPutRequestSiteValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *SitesPutRequestSiteValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *SitesPutRequestSiteValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *SitesPutRequestSiteValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *SitesPutRequestSiteValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetSuSupport

`func (o *SitesPutRequestSiteValue) GetSuSupport() bool`

GetSuSupport returns the SuSupport field if non-nil, zero value otherwise.

### GetSuSupportOk

`func (o *SitesPutRequestSiteValue) GetSuSupportOk() (*bool, bool)`

GetSuSupportOk returns a tuple with the SuSupport field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSuSupport

`func (o *SitesPutRequestSiteValue) SetSuSupport(v bool)`

SetSuSupport sets SuSupport field to given value.

### HasSuSupport

`func (o *SitesPutRequestSiteValue) HasSuSupport() bool`

HasSuSupport returns a boolean if a field has been set.

### GetAllowAllUnderlayConnections

`func (o *SitesPutRequestSiteValue) GetAllowAllUnderlayConnections() bool`

GetAllowAllUnderlayConnections returns the AllowAllUnderlayConnections field if non-nil, zero value otherwise.

### GetAllowAllUnderlayConnectionsOk

`func (o *SitesPutRequestSiteValue) GetAllowAllUnderlayConnectionsOk() (*bool, bool)`

GetAllowAllUnderlayConnectionsOk returns a tuple with the AllowAllUnderlayConnections field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllowAllUnderlayConnections

`func (o *SitesPutRequestSiteValue) SetAllowAllUnderlayConnections(v bool)`

SetAllowAllUnderlayConnections sets AllowAllUnderlayConnections field to given value.

### HasAllowAllUnderlayConnections

`func (o *SitesPutRequestSiteValue) HasAllowAllUnderlayConnections() bool`

HasAllowAllUnderlayConnections returns a boolean if a field has been set.

### GetSiteType

`func (o *SitesPutRequestSiteValue) GetSiteType() string`

GetSiteType returns the SiteType field if non-nil, zero value otherwise.

### GetSiteTypeOk

`func (o *SitesPutRequestSiteValue) GetSiteTypeOk() (*string, bool)`

GetSiteTypeOk returns a tuple with the SiteType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSiteType

`func (o *SitesPutRequestSiteValue) SetSiteType(v string)`

SetSiteType sets SiteType field to given value.

### HasSiteType

`func (o *SitesPutRequestSiteValue) HasSiteType() bool`

HasSiteType returns a boolean if a field has been set.

### GetDuplicateAddressDetectionMaxNumberOfMoves

`func (o *SitesPutRequestSiteValue) GetDuplicateAddressDetectionMaxNumberOfMoves() int32`

GetDuplicateAddressDetectionMaxNumberOfMoves returns the DuplicateAddressDetectionMaxNumberOfMoves field if non-nil, zero value otherwise.

### GetDuplicateAddressDetectionMaxNumberOfMovesOk

`func (o *SitesPutRequestSiteValue) GetDuplicateAddressDetectionMaxNumberOfMovesOk() (*int32, bool)`

GetDuplicateAddressDetectionMaxNumberOfMovesOk returns a tuple with the DuplicateAddressDetectionMaxNumberOfMoves field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDuplicateAddressDetectionMaxNumberOfMoves

`func (o *SitesPutRequestSiteValue) SetDuplicateAddressDetectionMaxNumberOfMoves(v int32)`

SetDuplicateAddressDetectionMaxNumberOfMoves sets DuplicateAddressDetectionMaxNumberOfMoves field to given value.

### HasDuplicateAddressDetectionMaxNumberOfMoves

`func (o *SitesPutRequestSiteValue) HasDuplicateAddressDetectionMaxNumberOfMoves() bool`

HasDuplicateAddressDetectionMaxNumberOfMoves returns a boolean if a field has been set.

### SetDuplicateAddressDetectionMaxNumberOfMovesNil

`func (o *SitesPutRequestSiteValue) SetDuplicateAddressDetectionMaxNumberOfMovesNil(b bool)`

 SetDuplicateAddressDetectionMaxNumberOfMovesNil sets the value for DuplicateAddressDetectionMaxNumberOfMoves to be an explicit nil

### UnsetDuplicateAddressDetectionMaxNumberOfMoves
`func (o *SitesPutRequestSiteValue) UnsetDuplicateAddressDetectionMaxNumberOfMoves()`

UnsetDuplicateAddressDetectionMaxNumberOfMoves ensures that no value is present for DuplicateAddressDetectionMaxNumberOfMoves, not even an explicit nil
### GetDuplicateAddressDetectionTime

`func (o *SitesPutRequestSiteValue) GetDuplicateAddressDetectionTime() int32`

GetDuplicateAddressDetectionTime returns the DuplicateAddressDetectionTime field if non-nil, zero value otherwise.

### GetDuplicateAddressDetectionTimeOk

`func (o *SitesPutRequestSiteValue) GetDuplicateAddressDetectionTimeOk() (*int32, bool)`

GetDuplicateAddressDetectionTimeOk returns a tuple with the DuplicateAddressDetectionTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDuplicateAddressDetectionTime

`func (o *SitesPutRequestSiteValue) SetDuplicateAddressDetectionTime(v int32)`

SetDuplicateAddressDetectionTime sets DuplicateAddressDetectionTime field to given value.

### HasDuplicateAddressDetectionTime

`func (o *SitesPutRequestSiteValue) HasDuplicateAddressDetectionTime() bool`

HasDuplicateAddressDetectionTime returns a boolean if a field has been set.

### SetDuplicateAddressDetectionTimeNil

`func (o *SitesPutRequestSiteValue) SetDuplicateAddressDetectionTimeNil(b bool)`

 SetDuplicateAddressDetectionTimeNil sets the value for DuplicateAddressDetectionTime to be an explicit nil

### UnsetDuplicateAddressDetectionTime
`func (o *SitesPutRequestSiteValue) UnsetDuplicateAddressDetectionTime()`

UnsetDuplicateAddressDetectionTime ensures that no value is present for DuplicateAddressDetectionTime, not even an explicit nil
### GetPortAdminPollingInterval

`func (o *SitesPutRequestSiteValue) GetPortAdminPollingInterval() int32`

GetPortAdminPollingInterval returns the PortAdminPollingInterval field if non-nil, zero value otherwise.

### GetPortAdminPollingIntervalOk

`func (o *SitesPutRequestSiteValue) GetPortAdminPollingIntervalOk() (*int32, bool)`

GetPortAdminPollingIntervalOk returns a tuple with the PortAdminPollingInterval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPortAdminPollingInterval

`func (o *SitesPutRequestSiteValue) SetPortAdminPollingInterval(v int32)`

SetPortAdminPollingInterval sets PortAdminPollingInterval field to given value.

### HasPortAdminPollingInterval

`func (o *SitesPutRequestSiteValue) HasPortAdminPollingInterval() bool`

HasPortAdminPollingInterval returns a boolean if a field has been set.

### SetPortAdminPollingIntervalNil

`func (o *SitesPutRequestSiteValue) SetPortAdminPollingIntervalNil(b bool)`

 SetPortAdminPollingIntervalNil sets the value for PortAdminPollingInterval to be an explicit nil

### UnsetPortAdminPollingInterval
`func (o *SitesPutRequestSiteValue) UnsetPortAdminPollingInterval()`

UnsetPortAdminPollingInterval ensures that no value is present for PortAdminPollingInterval, not even an explicit nil
### GetPortStatusPollingInterval

`func (o *SitesPutRequestSiteValue) GetPortStatusPollingInterval() int32`

GetPortStatusPollingInterval returns the PortStatusPollingInterval field if non-nil, zero value otherwise.

### GetPortStatusPollingIntervalOk

`func (o *SitesPutRequestSiteValue) GetPortStatusPollingIntervalOk() (*int32, bool)`

GetPortStatusPollingIntervalOk returns a tuple with the PortStatusPollingInterval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPortStatusPollingInterval

`func (o *SitesPutRequestSiteValue) SetPortStatusPollingInterval(v int32)`

SetPortStatusPollingInterval sets PortStatusPollingInterval field to given value.

### HasPortStatusPollingInterval

`func (o *SitesPutRequestSiteValue) HasPortStatusPollingInterval() bool`

HasPortStatusPollingInterval returns a boolean if a field has been set.

### SetPortStatusPollingIntervalNil

`func (o *SitesPutRequestSiteValue) SetPortStatusPollingIntervalNil(b bool)`

 SetPortStatusPollingIntervalNil sets the value for PortStatusPollingInterval to be an explicit nil

### UnsetPortStatusPollingInterval
`func (o *SitesPutRequestSiteValue) UnsetPortStatusPollingInterval()`

UnsetPortStatusPollingInterval ensures that no value is present for PortStatusPollingInterval, not even an explicit nil
### GetServiceForSite

`func (o *SitesPutRequestSiteValue) GetServiceForSite() string`

GetServiceForSite returns the ServiceForSite field if non-nil, zero value otherwise.

### GetServiceForSiteOk

`func (o *SitesPutRequestSiteValue) GetServiceForSiteOk() (*string, bool)`

GetServiceForSiteOk returns a tuple with the ServiceForSite field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServiceForSite

`func (o *SitesPutRequestSiteValue) SetServiceForSite(v string)`

SetServiceForSite sets ServiceForSite field to given value.

### HasServiceForSite

`func (o *SitesPutRequestSiteValue) HasServiceForSite() bool`

HasServiceForSite returns a boolean if a field has been set.

### GetServiceForSiteRefType

`func (o *SitesPutRequestSiteValue) GetServiceForSiteRefType() string`

GetServiceForSiteRefType returns the ServiceForSiteRefType field if non-nil, zero value otherwise.

### GetServiceForSiteRefTypeOk

`func (o *SitesPutRequestSiteValue) GetServiceForSiteRefTypeOk() (*string, bool)`

GetServiceForSiteRefTypeOk returns a tuple with the ServiceForSiteRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServiceForSiteRefType

`func (o *SitesPutRequestSiteValue) SetServiceForSiteRefType(v string)`

SetServiceForSiteRefType sets ServiceForSiteRefType field to given value.

### HasServiceForSiteRefType

`func (o *SitesPutRequestSiteValue) HasServiceForSiteRefType() bool`

HasServiceForSiteRefType returns a boolean if a field has been set.

### GetSpanningTreeType

`func (o *SitesPutRequestSiteValue) GetSpanningTreeType() string`

GetSpanningTreeType returns the SpanningTreeType field if non-nil, zero value otherwise.

### GetSpanningTreeTypeOk

`func (o *SitesPutRequestSiteValue) GetSpanningTreeTypeOk() (*string, bool)`

GetSpanningTreeTypeOk returns a tuple with the SpanningTreeType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpanningTreeType

`func (o *SitesPutRequestSiteValue) SetSpanningTreeType(v string)`

SetSpanningTreeType sets SpanningTreeType field to given value.

### HasSpanningTreeType

`func (o *SitesPutRequestSiteValue) HasSpanningTreeType() bool`

HasSpanningTreeType returns a boolean if a field has been set.

### GetRegionName

`func (o *SitesPutRequestSiteValue) GetRegionName() string`

GetRegionName returns the RegionName field if non-nil, zero value otherwise.

### GetRegionNameOk

`func (o *SitesPutRequestSiteValue) GetRegionNameOk() (*string, bool)`

GetRegionNameOk returns a tuple with the RegionName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegionName

`func (o *SitesPutRequestSiteValue) SetRegionName(v string)`

SetRegionName sets RegionName field to given value.

### HasRegionName

`func (o *SitesPutRequestSiteValue) HasRegionName() bool`

HasRegionName returns a boolean if a field has been set.

### GetRevision

`func (o *SitesPutRequestSiteValue) GetRevision() int32`

GetRevision returns the Revision field if non-nil, zero value otherwise.

### GetRevisionOk

`func (o *SitesPutRequestSiteValue) GetRevisionOk() (*int32, bool)`

GetRevisionOk returns a tuple with the Revision field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRevision

`func (o *SitesPutRequestSiteValue) SetRevision(v int32)`

SetRevision sets Revision field to given value.

### HasRevision

`func (o *SitesPutRequestSiteValue) HasRevision() bool`

HasRevision returns a boolean if a field has been set.

### SetRevisionNil

`func (o *SitesPutRequestSiteValue) SetRevisionNil(b bool)`

 SetRevisionNil sets the value for Revision to be an explicit nil

### UnsetRevision
`func (o *SitesPutRequestSiteValue) UnsetRevision()`

UnsetRevision ensures that no value is present for Revision, not even an explicit nil
### GetForceSpanningTreeOnFabricPorts

`func (o *SitesPutRequestSiteValue) GetForceSpanningTreeOnFabricPorts() bool`

GetForceSpanningTreeOnFabricPorts returns the ForceSpanningTreeOnFabricPorts field if non-nil, zero value otherwise.

### GetForceSpanningTreeOnFabricPortsOk

`func (o *SitesPutRequestSiteValue) GetForceSpanningTreeOnFabricPortsOk() (*bool, bool)`

GetForceSpanningTreeOnFabricPortsOk returns a tuple with the ForceSpanningTreeOnFabricPorts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetForceSpanningTreeOnFabricPorts

`func (o *SitesPutRequestSiteValue) SetForceSpanningTreeOnFabricPorts(v bool)`

SetForceSpanningTreeOnFabricPorts sets ForceSpanningTreeOnFabricPorts field to given value.

### HasForceSpanningTreeOnFabricPorts

`func (o *SitesPutRequestSiteValue) HasForceSpanningTreeOnFabricPorts() bool`

HasForceSpanningTreeOnFabricPorts returns a boolean if a field has been set.

### GetReadOnlyMode

`func (o *SitesPutRequestSiteValue) GetReadOnlyMode() bool`

GetReadOnlyMode returns the ReadOnlyMode field if non-nil, zero value otherwise.

### GetReadOnlyModeOk

`func (o *SitesPutRequestSiteValue) GetReadOnlyModeOk() (*bool, bool)`

GetReadOnlyModeOk returns a tuple with the ReadOnlyMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReadOnlyMode

`func (o *SitesPutRequestSiteValue) SetReadOnlyMode(v bool)`

SetReadOnlyMode sets ReadOnlyMode field to given value.

### HasReadOnlyMode

`func (o *SitesPutRequestSiteValue) HasReadOnlyMode() bool`

HasReadOnlyMode returns a boolean if a field has been set.

### GetDomainForSite

`func (o *SitesPutRequestSiteValue) GetDomainForSite() string`

GetDomainForSite returns the DomainForSite field if non-nil, zero value otherwise.

### GetDomainForSiteOk

`func (o *SitesPutRequestSiteValue) GetDomainForSiteOk() (*string, bool)`

GetDomainForSiteOk returns a tuple with the DomainForSite field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDomainForSite

`func (o *SitesPutRequestSiteValue) SetDomainForSite(v string)`

SetDomainForSite sets DomainForSite field to given value.

### HasDomainForSite

`func (o *SitesPutRequestSiteValue) HasDomainForSite() bool`

HasDomainForSite returns a boolean if a field has been set.

### GetDomainForSiteRefType

`func (o *SitesPutRequestSiteValue) GetDomainForSiteRefType() string`

GetDomainForSiteRefType returns the DomainForSiteRefType field if non-nil, zero value otherwise.

### GetDomainForSiteRefTypeOk

`func (o *SitesPutRequestSiteValue) GetDomainForSiteRefTypeOk() (*string, bool)`

GetDomainForSiteRefTypeOk returns a tuple with the DomainForSiteRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDomainForSiteRefType

`func (o *SitesPutRequestSiteValue) SetDomainForSiteRefType(v string)`

SetDomainForSiteRefType sets DomainForSiteRefType field to given value.

### HasDomainForSiteRefType

`func (o *SitesPutRequestSiteValue) HasDomainForSiteRefType() bool`

HasDomainForSiteRefType returns a boolean if a field has been set.

### GetEnableDscp

`func (o *SitesPutRequestSiteValue) GetEnableDscp() bool`

GetEnableDscp returns the EnableDscp field if non-nil, zero value otherwise.

### GetEnableDscpOk

`func (o *SitesPutRequestSiteValue) GetEnableDscpOk() (*bool, bool)`

GetEnableDscpOk returns a tuple with the EnableDscp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnableDscp

`func (o *SitesPutRequestSiteValue) SetEnableDscp(v bool)`

SetEnableDscp sets EnableDscp field to given value.

### HasEnableDscp

`func (o *SitesPutRequestSiteValue) HasEnableDscp() bool`

HasEnableDscp returns a boolean if a field has been set.

### GetDscpToPBitMap

`func (o *SitesPutRequestSiteValue) GetDscpToPBitMap() string`

GetDscpToPBitMap returns the DscpToPBitMap field if non-nil, zero value otherwise.

### GetDscpToPBitMapOk

`func (o *SitesPutRequestSiteValue) GetDscpToPBitMapOk() (*string, bool)`

GetDscpToPBitMapOk returns a tuple with the DscpToPBitMap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDscpToPBitMap

`func (o *SitesPutRequestSiteValue) SetDscpToPBitMap(v string)`

SetDscpToPBitMap sets DscpToPBitMap field to given value.

### HasDscpToPBitMap

`func (o *SitesPutRequestSiteValue) HasDscpToPBitMap() bool`

HasDscpToPBitMap returns a boolean if a field has been set.

### GetAnycastMacAddress

`func (o *SitesPutRequestSiteValue) GetAnycastMacAddress() string`

GetAnycastMacAddress returns the AnycastMacAddress field if non-nil, zero value otherwise.

### GetAnycastMacAddressOk

`func (o *SitesPutRequestSiteValue) GetAnycastMacAddressOk() (*string, bool)`

GetAnycastMacAddressOk returns a tuple with the AnycastMacAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnycastMacAddress

`func (o *SitesPutRequestSiteValue) SetAnycastMacAddress(v string)`

SetAnycastMacAddress sets AnycastMacAddress field to given value.

### HasAnycastMacAddress

`func (o *SitesPutRequestSiteValue) HasAnycastMacAddress() bool`

HasAnycastMacAddress returns a boolean if a field has been set.

### GetAnycastMacAddressAutoAssigned

`func (o *SitesPutRequestSiteValue) GetAnycastMacAddressAutoAssigned() bool`

GetAnycastMacAddressAutoAssigned returns the AnycastMacAddressAutoAssigned field if non-nil, zero value otherwise.

### GetAnycastMacAddressAutoAssignedOk

`func (o *SitesPutRequestSiteValue) GetAnycastMacAddressAutoAssignedOk() (*bool, bool)`

GetAnycastMacAddressAutoAssignedOk returns a tuple with the AnycastMacAddressAutoAssigned field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnycastMacAddressAutoAssigned

`func (o *SitesPutRequestSiteValue) SetAnycastMacAddressAutoAssigned(v bool)`

SetAnycastMacAddressAutoAssigned sets AnycastMacAddressAutoAssigned field to given value.

### HasAnycastMacAddressAutoAssigned

`func (o *SitesPutRequestSiteValue) HasAnycastMacAddressAutoAssigned() bool`

HasAnycastMacAddressAutoAssigned returns a boolean if a field has been set.

### GetMacAddressAgingTime

`func (o *SitesPutRequestSiteValue) GetMacAddressAgingTime() int32`

GetMacAddressAgingTime returns the MacAddressAgingTime field if non-nil, zero value otherwise.

### GetMacAddressAgingTimeOk

`func (o *SitesPutRequestSiteValue) GetMacAddressAgingTimeOk() (*int32, bool)`

GetMacAddressAgingTimeOk returns a tuple with the MacAddressAgingTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMacAddressAgingTime

`func (o *SitesPutRequestSiteValue) SetMacAddressAgingTime(v int32)`

SetMacAddressAgingTime sets MacAddressAgingTime field to given value.

### HasMacAddressAgingTime

`func (o *SitesPutRequestSiteValue) HasMacAddressAgingTime() bool`

HasMacAddressAgingTime returns a boolean if a field has been set.

### SetMacAddressAgingTimeNil

`func (o *SitesPutRequestSiteValue) SetMacAddressAgingTimeNil(b bool)`

 SetMacAddressAgingTimeNil sets the value for MacAddressAgingTime to be an explicit nil

### UnsetMacAddressAgingTime
`func (o *SitesPutRequestSiteValue) UnsetMacAddressAgingTime()`

UnsetMacAddressAgingTime ensures that no value is present for MacAddressAgingTime, not even an explicit nil
### GetMlagDelayRestoreTimer

`func (o *SitesPutRequestSiteValue) GetMlagDelayRestoreTimer() int32`

GetMlagDelayRestoreTimer returns the MlagDelayRestoreTimer field if non-nil, zero value otherwise.

### GetMlagDelayRestoreTimerOk

`func (o *SitesPutRequestSiteValue) GetMlagDelayRestoreTimerOk() (*int32, bool)`

GetMlagDelayRestoreTimerOk returns a tuple with the MlagDelayRestoreTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMlagDelayRestoreTimer

`func (o *SitesPutRequestSiteValue) SetMlagDelayRestoreTimer(v int32)`

SetMlagDelayRestoreTimer sets MlagDelayRestoreTimer field to given value.

### HasMlagDelayRestoreTimer

`func (o *SitesPutRequestSiteValue) HasMlagDelayRestoreTimer() bool`

HasMlagDelayRestoreTimer returns a boolean if a field has been set.

### SetMlagDelayRestoreTimerNil

`func (o *SitesPutRequestSiteValue) SetMlagDelayRestoreTimerNil(b bool)`

 SetMlagDelayRestoreTimerNil sets the value for MlagDelayRestoreTimer to be an explicit nil

### UnsetMlagDelayRestoreTimer
`func (o *SitesPutRequestSiteValue) UnsetMlagDelayRestoreTimer()`

UnsetMlagDelayRestoreTimer ensures that no value is present for MlagDelayRestoreTimer, not even an explicit nil
### GetBgpKeepaliveTimer

`func (o *SitesPutRequestSiteValue) GetBgpKeepaliveTimer() int32`

GetBgpKeepaliveTimer returns the BgpKeepaliveTimer field if non-nil, zero value otherwise.

### GetBgpKeepaliveTimerOk

`func (o *SitesPutRequestSiteValue) GetBgpKeepaliveTimerOk() (*int32, bool)`

GetBgpKeepaliveTimerOk returns a tuple with the BgpKeepaliveTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBgpKeepaliveTimer

`func (o *SitesPutRequestSiteValue) SetBgpKeepaliveTimer(v int32)`

SetBgpKeepaliveTimer sets BgpKeepaliveTimer field to given value.

### HasBgpKeepaliveTimer

`func (o *SitesPutRequestSiteValue) HasBgpKeepaliveTimer() bool`

HasBgpKeepaliveTimer returns a boolean if a field has been set.

### SetBgpKeepaliveTimerNil

`func (o *SitesPutRequestSiteValue) SetBgpKeepaliveTimerNil(b bool)`

 SetBgpKeepaliveTimerNil sets the value for BgpKeepaliveTimer to be an explicit nil

### UnsetBgpKeepaliveTimer
`func (o *SitesPutRequestSiteValue) UnsetBgpKeepaliveTimer()`

UnsetBgpKeepaliveTimer ensures that no value is present for BgpKeepaliveTimer, not even an explicit nil
### GetBgpHoldDownTimer

`func (o *SitesPutRequestSiteValue) GetBgpHoldDownTimer() int32`

GetBgpHoldDownTimer returns the BgpHoldDownTimer field if non-nil, zero value otherwise.

### GetBgpHoldDownTimerOk

`func (o *SitesPutRequestSiteValue) GetBgpHoldDownTimerOk() (*int32, bool)`

GetBgpHoldDownTimerOk returns a tuple with the BgpHoldDownTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBgpHoldDownTimer

`func (o *SitesPutRequestSiteValue) SetBgpHoldDownTimer(v int32)`

SetBgpHoldDownTimer sets BgpHoldDownTimer field to given value.

### HasBgpHoldDownTimer

`func (o *SitesPutRequestSiteValue) HasBgpHoldDownTimer() bool`

HasBgpHoldDownTimer returns a boolean if a field has been set.

### SetBgpHoldDownTimerNil

`func (o *SitesPutRequestSiteValue) SetBgpHoldDownTimerNil(b bool)`

 SetBgpHoldDownTimerNil sets the value for BgpHoldDownTimer to be an explicit nil

### UnsetBgpHoldDownTimer
`func (o *SitesPutRequestSiteValue) UnsetBgpHoldDownTimer()`

UnsetBgpHoldDownTimer ensures that no value is present for BgpHoldDownTimer, not even an explicit nil
### GetSpineBgpAdvertisementInterval

`func (o *SitesPutRequestSiteValue) GetSpineBgpAdvertisementInterval() int32`

GetSpineBgpAdvertisementInterval returns the SpineBgpAdvertisementInterval field if non-nil, zero value otherwise.

### GetSpineBgpAdvertisementIntervalOk

`func (o *SitesPutRequestSiteValue) GetSpineBgpAdvertisementIntervalOk() (*int32, bool)`

GetSpineBgpAdvertisementIntervalOk returns a tuple with the SpineBgpAdvertisementInterval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpineBgpAdvertisementInterval

`func (o *SitesPutRequestSiteValue) SetSpineBgpAdvertisementInterval(v int32)`

SetSpineBgpAdvertisementInterval sets SpineBgpAdvertisementInterval field to given value.

### HasSpineBgpAdvertisementInterval

`func (o *SitesPutRequestSiteValue) HasSpineBgpAdvertisementInterval() bool`

HasSpineBgpAdvertisementInterval returns a boolean if a field has been set.

### SetSpineBgpAdvertisementIntervalNil

`func (o *SitesPutRequestSiteValue) SetSpineBgpAdvertisementIntervalNil(b bool)`

 SetSpineBgpAdvertisementIntervalNil sets the value for SpineBgpAdvertisementInterval to be an explicit nil

### UnsetSpineBgpAdvertisementInterval
`func (o *SitesPutRequestSiteValue) UnsetSpineBgpAdvertisementInterval()`

UnsetSpineBgpAdvertisementInterval ensures that no value is present for SpineBgpAdvertisementInterval, not even an explicit nil
### GetSpineBgpConnectTimer

`func (o *SitesPutRequestSiteValue) GetSpineBgpConnectTimer() int32`

GetSpineBgpConnectTimer returns the SpineBgpConnectTimer field if non-nil, zero value otherwise.

### GetSpineBgpConnectTimerOk

`func (o *SitesPutRequestSiteValue) GetSpineBgpConnectTimerOk() (*int32, bool)`

GetSpineBgpConnectTimerOk returns a tuple with the SpineBgpConnectTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpineBgpConnectTimer

`func (o *SitesPutRequestSiteValue) SetSpineBgpConnectTimer(v int32)`

SetSpineBgpConnectTimer sets SpineBgpConnectTimer field to given value.

### HasSpineBgpConnectTimer

`func (o *SitesPutRequestSiteValue) HasSpineBgpConnectTimer() bool`

HasSpineBgpConnectTimer returns a boolean if a field has been set.

### SetSpineBgpConnectTimerNil

`func (o *SitesPutRequestSiteValue) SetSpineBgpConnectTimerNil(b bool)`

 SetSpineBgpConnectTimerNil sets the value for SpineBgpConnectTimer to be an explicit nil

### UnsetSpineBgpConnectTimer
`func (o *SitesPutRequestSiteValue) UnsetSpineBgpConnectTimer()`

UnsetSpineBgpConnectTimer ensures that no value is present for SpineBgpConnectTimer, not even an explicit nil
### GetSpineAsNumber

`func (o *SitesPutRequestSiteValue) GetSpineAsNumber() int32`

GetSpineAsNumber returns the SpineAsNumber field if non-nil, zero value otherwise.

### GetSpineAsNumberOk

`func (o *SitesPutRequestSiteValue) GetSpineAsNumberOk() (*int32, bool)`

GetSpineAsNumberOk returns a tuple with the SpineAsNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpineAsNumber

`func (o *SitesPutRequestSiteValue) SetSpineAsNumber(v int32)`

SetSpineAsNumber sets SpineAsNumber field to given value.

### HasSpineAsNumber

`func (o *SitesPutRequestSiteValue) HasSpineAsNumber() bool`

HasSpineAsNumber returns a boolean if a field has been set.

### SetSpineAsNumberNil

`func (o *SitesPutRequestSiteValue) SetSpineAsNumberNil(b bool)`

 SetSpineAsNumberNil sets the value for SpineAsNumber to be an explicit nil

### UnsetSpineAsNumber
`func (o *SitesPutRequestSiteValue) UnsetSpineAsNumber()`

UnsetSpineAsNumber ensures that no value is present for SpineAsNumber, not even an explicit nil
### GetLeafBgpKeepAliveTimer

`func (o *SitesPutRequestSiteValue) GetLeafBgpKeepAliveTimer() int32`

GetLeafBgpKeepAliveTimer returns the LeafBgpKeepAliveTimer field if non-nil, zero value otherwise.

### GetLeafBgpKeepAliveTimerOk

`func (o *SitesPutRequestSiteValue) GetLeafBgpKeepAliveTimerOk() (*int32, bool)`

GetLeafBgpKeepAliveTimerOk returns a tuple with the LeafBgpKeepAliveTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLeafBgpKeepAliveTimer

`func (o *SitesPutRequestSiteValue) SetLeafBgpKeepAliveTimer(v int32)`

SetLeafBgpKeepAliveTimer sets LeafBgpKeepAliveTimer field to given value.

### HasLeafBgpKeepAliveTimer

`func (o *SitesPutRequestSiteValue) HasLeafBgpKeepAliveTimer() bool`

HasLeafBgpKeepAliveTimer returns a boolean if a field has been set.

### SetLeafBgpKeepAliveTimerNil

`func (o *SitesPutRequestSiteValue) SetLeafBgpKeepAliveTimerNil(b bool)`

 SetLeafBgpKeepAliveTimerNil sets the value for LeafBgpKeepAliveTimer to be an explicit nil

### UnsetLeafBgpKeepAliveTimer
`func (o *SitesPutRequestSiteValue) UnsetLeafBgpKeepAliveTimer()`

UnsetLeafBgpKeepAliveTimer ensures that no value is present for LeafBgpKeepAliveTimer, not even an explicit nil
### GetLeafBgpHoldDownTimer

`func (o *SitesPutRequestSiteValue) GetLeafBgpHoldDownTimer() int32`

GetLeafBgpHoldDownTimer returns the LeafBgpHoldDownTimer field if non-nil, zero value otherwise.

### GetLeafBgpHoldDownTimerOk

`func (o *SitesPutRequestSiteValue) GetLeafBgpHoldDownTimerOk() (*int32, bool)`

GetLeafBgpHoldDownTimerOk returns a tuple with the LeafBgpHoldDownTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLeafBgpHoldDownTimer

`func (o *SitesPutRequestSiteValue) SetLeafBgpHoldDownTimer(v int32)`

SetLeafBgpHoldDownTimer sets LeafBgpHoldDownTimer field to given value.

### HasLeafBgpHoldDownTimer

`func (o *SitesPutRequestSiteValue) HasLeafBgpHoldDownTimer() bool`

HasLeafBgpHoldDownTimer returns a boolean if a field has been set.

### SetLeafBgpHoldDownTimerNil

`func (o *SitesPutRequestSiteValue) SetLeafBgpHoldDownTimerNil(b bool)`

 SetLeafBgpHoldDownTimerNil sets the value for LeafBgpHoldDownTimer to be an explicit nil

### UnsetLeafBgpHoldDownTimer
`func (o *SitesPutRequestSiteValue) UnsetLeafBgpHoldDownTimer()`

UnsetLeafBgpHoldDownTimer ensures that no value is present for LeafBgpHoldDownTimer, not even an explicit nil
### GetLeafBgpAdvertisementInterval

`func (o *SitesPutRequestSiteValue) GetLeafBgpAdvertisementInterval() int32`

GetLeafBgpAdvertisementInterval returns the LeafBgpAdvertisementInterval field if non-nil, zero value otherwise.

### GetLeafBgpAdvertisementIntervalOk

`func (o *SitesPutRequestSiteValue) GetLeafBgpAdvertisementIntervalOk() (*int32, bool)`

GetLeafBgpAdvertisementIntervalOk returns a tuple with the LeafBgpAdvertisementInterval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLeafBgpAdvertisementInterval

`func (o *SitesPutRequestSiteValue) SetLeafBgpAdvertisementInterval(v int32)`

SetLeafBgpAdvertisementInterval sets LeafBgpAdvertisementInterval field to given value.

### HasLeafBgpAdvertisementInterval

`func (o *SitesPutRequestSiteValue) HasLeafBgpAdvertisementInterval() bool`

HasLeafBgpAdvertisementInterval returns a boolean if a field has been set.

### SetLeafBgpAdvertisementIntervalNil

`func (o *SitesPutRequestSiteValue) SetLeafBgpAdvertisementIntervalNil(b bool)`

 SetLeafBgpAdvertisementIntervalNil sets the value for LeafBgpAdvertisementInterval to be an explicit nil

### UnsetLeafBgpAdvertisementInterval
`func (o *SitesPutRequestSiteValue) UnsetLeafBgpAdvertisementInterval()`

UnsetLeafBgpAdvertisementInterval ensures that no value is present for LeafBgpAdvertisementInterval, not even an explicit nil
### GetLeafBgpConnectTimer

`func (o *SitesPutRequestSiteValue) GetLeafBgpConnectTimer() int32`

GetLeafBgpConnectTimer returns the LeafBgpConnectTimer field if non-nil, zero value otherwise.

### GetLeafBgpConnectTimerOk

`func (o *SitesPutRequestSiteValue) GetLeafBgpConnectTimerOk() (*int32, bool)`

GetLeafBgpConnectTimerOk returns a tuple with the LeafBgpConnectTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLeafBgpConnectTimer

`func (o *SitesPutRequestSiteValue) SetLeafBgpConnectTimer(v int32)`

SetLeafBgpConnectTimer sets LeafBgpConnectTimer field to given value.

### HasLeafBgpConnectTimer

`func (o *SitesPutRequestSiteValue) HasLeafBgpConnectTimer() bool`

HasLeafBgpConnectTimer returns a boolean if a field has been set.

### SetLeafBgpConnectTimerNil

`func (o *SitesPutRequestSiteValue) SetLeafBgpConnectTimerNil(b bool)`

 SetLeafBgpConnectTimerNil sets the value for LeafBgpConnectTimer to be an explicit nil

### UnsetLeafBgpConnectTimer
`func (o *SitesPutRequestSiteValue) UnsetLeafBgpConnectTimer()`

UnsetLeafBgpConnectTimer ensures that no value is present for LeafBgpConnectTimer, not even an explicit nil
### GetLinkStateTimeoutValue

`func (o *SitesPutRequestSiteValue) GetLinkStateTimeoutValue() int32`

GetLinkStateTimeoutValue returns the LinkStateTimeoutValue field if non-nil, zero value otherwise.

### GetLinkStateTimeoutValueOk

`func (o *SitesPutRequestSiteValue) GetLinkStateTimeoutValueOk() (*int32, bool)`

GetLinkStateTimeoutValueOk returns a tuple with the LinkStateTimeoutValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLinkStateTimeoutValue

`func (o *SitesPutRequestSiteValue) SetLinkStateTimeoutValue(v int32)`

SetLinkStateTimeoutValue sets LinkStateTimeoutValue field to given value.

### HasLinkStateTimeoutValue

`func (o *SitesPutRequestSiteValue) HasLinkStateTimeoutValue() bool`

HasLinkStateTimeoutValue returns a boolean if a field has been set.

### SetLinkStateTimeoutValueNil

`func (o *SitesPutRequestSiteValue) SetLinkStateTimeoutValueNil(b bool)`

 SetLinkStateTimeoutValueNil sets the value for LinkStateTimeoutValue to be an explicit nil

### UnsetLinkStateTimeoutValue
`func (o *SitesPutRequestSiteValue) UnsetLinkStateTimeoutValue()`

UnsetLinkStateTimeoutValue ensures that no value is present for LinkStateTimeoutValue, not even an explicit nil
### GetEvpnMultihomingStartupDelay

`func (o *SitesPutRequestSiteValue) GetEvpnMultihomingStartupDelay() int32`

GetEvpnMultihomingStartupDelay returns the EvpnMultihomingStartupDelay field if non-nil, zero value otherwise.

### GetEvpnMultihomingStartupDelayOk

`func (o *SitesPutRequestSiteValue) GetEvpnMultihomingStartupDelayOk() (*int32, bool)`

GetEvpnMultihomingStartupDelayOk returns a tuple with the EvpnMultihomingStartupDelay field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEvpnMultihomingStartupDelay

`func (o *SitesPutRequestSiteValue) SetEvpnMultihomingStartupDelay(v int32)`

SetEvpnMultihomingStartupDelay sets EvpnMultihomingStartupDelay field to given value.

### HasEvpnMultihomingStartupDelay

`func (o *SitesPutRequestSiteValue) HasEvpnMultihomingStartupDelay() bool`

HasEvpnMultihomingStartupDelay returns a boolean if a field has been set.

### SetEvpnMultihomingStartupDelayNil

`func (o *SitesPutRequestSiteValue) SetEvpnMultihomingStartupDelayNil(b bool)`

 SetEvpnMultihomingStartupDelayNil sets the value for EvpnMultihomingStartupDelay to be an explicit nil

### UnsetEvpnMultihomingStartupDelay
`func (o *SitesPutRequestSiteValue) UnsetEvpnMultihomingStartupDelay()`

UnsetEvpnMultihomingStartupDelay ensures that no value is present for EvpnMultihomingStartupDelay, not even an explicit nil
### GetEvpnMacHoldtime

`func (o *SitesPutRequestSiteValue) GetEvpnMacHoldtime() int32`

GetEvpnMacHoldtime returns the EvpnMacHoldtime field if non-nil, zero value otherwise.

### GetEvpnMacHoldtimeOk

`func (o *SitesPutRequestSiteValue) GetEvpnMacHoldtimeOk() (*int32, bool)`

GetEvpnMacHoldtimeOk returns a tuple with the EvpnMacHoldtime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEvpnMacHoldtime

`func (o *SitesPutRequestSiteValue) SetEvpnMacHoldtime(v int32)`

SetEvpnMacHoldtime sets EvpnMacHoldtime field to given value.

### HasEvpnMacHoldtime

`func (o *SitesPutRequestSiteValue) HasEvpnMacHoldtime() bool`

HasEvpnMacHoldtime returns a boolean if a field has been set.

### SetEvpnMacHoldtimeNil

`func (o *SitesPutRequestSiteValue) SetEvpnMacHoldtimeNil(b bool)`

 SetEvpnMacHoldtimeNil sets the value for EvpnMacHoldtime to be an explicit nil

### UnsetEvpnMacHoldtime
`func (o *SitesPutRequestSiteValue) UnsetEvpnMacHoldtime()`

UnsetEvpnMacHoldtime ensures that no value is present for EvpnMacHoldtime, not even an explicit nil
### GetAggressiveReporting

`func (o *SitesPutRequestSiteValue) GetAggressiveReporting() bool`

GetAggressiveReporting returns the AggressiveReporting field if non-nil, zero value otherwise.

### GetAggressiveReportingOk

`func (o *SitesPutRequestSiteValue) GetAggressiveReportingOk() (*bool, bool)`

GetAggressiveReportingOk returns a tuple with the AggressiveReporting field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAggressiveReporting

`func (o *SitesPutRequestSiteValue) SetAggressiveReporting(v bool)`

SetAggressiveReporting sets AggressiveReporting field to given value.

### HasAggressiveReporting

`func (o *SitesPutRequestSiteValue) HasAggressiveReporting() bool`

HasAggressiveReporting returns a boolean if a field has been set.

### GetSwitchIpBase

`func (o *SitesPutRequestSiteValue) GetSwitchIpBase() string`

GetSwitchIpBase returns the SwitchIpBase field if non-nil, zero value otherwise.

### GetSwitchIpBaseOk

`func (o *SitesPutRequestSiteValue) GetSwitchIpBaseOk() (*string, bool)`

GetSwitchIpBaseOk returns a tuple with the SwitchIpBase field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchIpBase

`func (o *SitesPutRequestSiteValue) SetSwitchIpBase(v string)`

SetSwitchIpBase sets SwitchIpBase field to given value.

### HasSwitchIpBase

`func (o *SitesPutRequestSiteValue) HasSwitchIpBase() bool`

HasSwitchIpBase returns a boolean if a field has been set.

### GetControllerIpBase

`func (o *SitesPutRequestSiteValue) GetControllerIpBase() string`

GetControllerIpBase returns the ControllerIpBase field if non-nil, zero value otherwise.

### GetControllerIpBaseOk

`func (o *SitesPutRequestSiteValue) GetControllerIpBaseOk() (*string, bool)`

GetControllerIpBaseOk returns a tuple with the ControllerIpBase field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetControllerIpBase

`func (o *SitesPutRequestSiteValue) SetControllerIpBase(v string)`

SetControllerIpBase sets ControllerIpBase field to given value.

### HasControllerIpBase

`func (o *SitesPutRequestSiteValue) HasControllerIpBase() bool`

HasControllerIpBase returns a boolean if a field has been set.

### GetMultiTenant

`func (o *SitesPutRequestSiteValue) GetMultiTenant() bool`

GetMultiTenant returns the MultiTenant field if non-nil, zero value otherwise.

### GetMultiTenantOk

`func (o *SitesPutRequestSiteValue) GetMultiTenantOk() (*bool, bool)`

GetMultiTenantOk returns a tuple with the MultiTenant field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMultiTenant

`func (o *SitesPutRequestSiteValue) SetMultiTenant(v bool)`

SetMultiTenant sets MultiTenant field to given value.

### HasMultiTenant

`func (o *SitesPutRequestSiteValue) HasMultiTenant() bool`

HasMultiTenant returns a boolean if a field has been set.

### GetBaseBgpAsNumber

`func (o *SitesPutRequestSiteValue) GetBaseBgpAsNumber() string`

GetBaseBgpAsNumber returns the BaseBgpAsNumber field if non-nil, zero value otherwise.

### GetBaseBgpAsNumberOk

`func (o *SitesPutRequestSiteValue) GetBaseBgpAsNumberOk() (*string, bool)`

GetBaseBgpAsNumberOk returns a tuple with the BaseBgpAsNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBaseBgpAsNumber

`func (o *SitesPutRequestSiteValue) SetBaseBgpAsNumber(v string)`

SetBaseBgpAsNumber sets BaseBgpAsNumber field to given value.

### HasBaseBgpAsNumber

`func (o *SitesPutRequestSiteValue) HasBaseBgpAsNumber() bool`

HasBaseBgpAsNumber returns a boolean if a field has been set.

### GetRouterIdBasePrefix

`func (o *SitesPutRequestSiteValue) GetRouterIdBasePrefix() string`

GetRouterIdBasePrefix returns the RouterIdBasePrefix field if non-nil, zero value otherwise.

### GetRouterIdBasePrefixOk

`func (o *SitesPutRequestSiteValue) GetRouterIdBasePrefixOk() (*string, bool)`

GetRouterIdBasePrefixOk returns a tuple with the RouterIdBasePrefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRouterIdBasePrefix

`func (o *SitesPutRequestSiteValue) SetRouterIdBasePrefix(v string)`

SetRouterIdBasePrefix sets RouterIdBasePrefix field to given value.

### HasRouterIdBasePrefix

`func (o *SitesPutRequestSiteValue) HasRouterIdBasePrefix() bool`

HasRouterIdBasePrefix returns a boolean if a field has been set.

### GetVtepIdBasePrefix

`func (o *SitesPutRequestSiteValue) GetVtepIdBasePrefix() string`

GetVtepIdBasePrefix returns the VtepIdBasePrefix field if non-nil, zero value otherwise.

### GetVtepIdBasePrefixOk

`func (o *SitesPutRequestSiteValue) GetVtepIdBasePrefixOk() (*string, bool)`

GetVtepIdBasePrefixOk returns a tuple with the VtepIdBasePrefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVtepIdBasePrefix

`func (o *SitesPutRequestSiteValue) SetVtepIdBasePrefix(v string)`

SetVtepIdBasePrefix sets VtepIdBasePrefix field to given value.

### HasVtepIdBasePrefix

`func (o *SitesPutRequestSiteValue) HasVtepIdBasePrefix() bool`

HasVtepIdBasePrefix returns a boolean if a field has been set.

### GetPairedIpSubnet

`func (o *SitesPutRequestSiteValue) GetPairedIpSubnet() string`

GetPairedIpSubnet returns the PairedIpSubnet field if non-nil, zero value otherwise.

### GetPairedIpSubnetOk

`func (o *SitesPutRequestSiteValue) GetPairedIpSubnetOk() (*string, bool)`

GetPairedIpSubnetOk returns a tuple with the PairedIpSubnet field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPairedIpSubnet

`func (o *SitesPutRequestSiteValue) SetPairedIpSubnet(v string)`

SetPairedIpSubnet sets PairedIpSubnet field to given value.

### HasPairedIpSubnet

`func (o *SitesPutRequestSiteValue) HasPairedIpSubnet() bool`

HasPairedIpSubnet returns a boolean if a field has been set.

### GetMaxSwitches

`func (o *SitesPutRequestSiteValue) GetMaxSwitches() string`

GetMaxSwitches returns the MaxSwitches field if non-nil, zero value otherwise.

### GetMaxSwitchesOk

`func (o *SitesPutRequestSiteValue) GetMaxSwitchesOk() (*string, bool)`

GetMaxSwitchesOk returns a tuple with the MaxSwitches field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxSwitches

`func (o *SitesPutRequestSiteValue) SetMaxSwitches(v string)`

SetMaxSwitches sets MaxSwitches field to given value.

### HasMaxSwitches

`func (o *SitesPutRequestSiteValue) HasMaxSwitches() bool`

HasMaxSwitches returns a boolean if a field has been set.

### GetPauseValidationAlarms

`func (o *SitesPutRequestSiteValue) GetPauseValidationAlarms() bool`

GetPauseValidationAlarms returns the PauseValidationAlarms field if non-nil, zero value otherwise.

### GetPauseValidationAlarmsOk

`func (o *SitesPutRequestSiteValue) GetPauseValidationAlarmsOk() (*bool, bool)`

GetPauseValidationAlarmsOk returns a tuple with the PauseValidationAlarms field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPauseValidationAlarms

`func (o *SitesPutRequestSiteValue) SetPauseValidationAlarms(v bool)`

SetPauseValidationAlarms sets PauseValidationAlarms field to given value.

### HasPauseValidationAlarms

`func (o *SitesPutRequestSiteValue) HasPauseValidationAlarms() bool`

HasPauseValidationAlarms returns a boolean if a field has been set.

### GetStartingOctet

`func (o *SitesPutRequestSiteValue) GetStartingOctet() int32`

GetStartingOctet returns the StartingOctet field if non-nil, zero value otherwise.

### GetStartingOctetOk

`func (o *SitesPutRequestSiteValue) GetStartingOctetOk() (*int32, bool)`

GetStartingOctetOk returns a tuple with the StartingOctet field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStartingOctet

`func (o *SitesPutRequestSiteValue) SetStartingOctet(v int32)`

SetStartingOctet sets StartingOctet field to given value.

### HasStartingOctet

`func (o *SitesPutRequestSiteValue) HasStartingOctet() bool`

HasStartingOctet returns a boolean if a field has been set.

### SetStartingOctetNil

`func (o *SitesPutRequestSiteValue) SetStartingOctetNil(b bool)`

 SetStartingOctetNil sets the value for StartingOctet to be an explicit nil

### UnsetStartingOctet
`func (o *SitesPutRequestSiteValue) UnsetStartingOctet()`

UnsetStartingOctet ensures that no value is present for StartingOctet, not even an explicit nil
### GetMaxSus

`func (o *SitesPutRequestSiteValue) GetMaxSus() int32`

GetMaxSus returns the MaxSus field if non-nil, zero value otherwise.

### GetMaxSusOk

`func (o *SitesPutRequestSiteValue) GetMaxSusOk() (*int32, bool)`

GetMaxSusOk returns a tuple with the MaxSus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxSus

`func (o *SitesPutRequestSiteValue) SetMaxSus(v int32)`

SetMaxSus sets MaxSus field to given value.

### HasMaxSus

`func (o *SitesPutRequestSiteValue) HasMaxSus() bool`

HasMaxSus returns a boolean if a field has been set.

### SetMaxSusNil

`func (o *SitesPutRequestSiteValue) SetMaxSusNil(b bool)`

 SetMaxSusNil sets the value for MaxSus to be an explicit nil

### UnsetMaxSus
`func (o *SitesPutRequestSiteValue) UnsetMaxSus()`

UnsetMaxSus ensures that no value is present for MaxSus, not even an explicit nil
### GetMaxPods

`func (o *SitesPutRequestSiteValue) GetMaxPods() int32`

GetMaxPods returns the MaxPods field if non-nil, zero value otherwise.

### GetMaxPodsOk

`func (o *SitesPutRequestSiteValue) GetMaxPodsOk() (*int32, bool)`

GetMaxPodsOk returns a tuple with the MaxPods field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxPods

`func (o *SitesPutRequestSiteValue) SetMaxPods(v int32)`

SetMaxPods sets MaxPods field to given value.

### HasMaxPods

`func (o *SitesPutRequestSiteValue) HasMaxPods() bool`

HasMaxPods returns a boolean if a field has been set.

### SetMaxPodsNil

`func (o *SitesPutRequestSiteValue) SetMaxPodsNil(b bool)`

 SetMaxPodsNil sets the value for MaxPods to be an explicit nil

### UnsetMaxPods
`func (o *SitesPutRequestSiteValue) UnsetMaxPods()`

UnsetMaxPods ensures that no value is present for MaxPods, not even an explicit nil
### GetObjectProperties

`func (o *SitesPutRequestSiteValue) GetObjectProperties() SitesPutRequestSiteValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *SitesPutRequestSiteValue) GetObjectPropertiesOk() (*SitesPutRequestSiteValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *SitesPutRequestSiteValue) SetObjectProperties(v SitesPutRequestSiteValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *SitesPutRequestSiteValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.

### GetSwitchUsername

`func (o *SitesPutRequestSiteValue) GetSwitchUsername() string`

GetSwitchUsername returns the SwitchUsername field if non-nil, zero value otherwise.

### GetSwitchUsernameOk

`func (o *SitesPutRequestSiteValue) GetSwitchUsernameOk() (*string, bool)`

GetSwitchUsernameOk returns a tuple with the SwitchUsername field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchUsername

`func (o *SitesPutRequestSiteValue) SetSwitchUsername(v string)`

SetSwitchUsername sets SwitchUsername field to given value.

### HasSwitchUsername

`func (o *SitesPutRequestSiteValue) HasSwitchUsername() bool`

HasSwitchUsername returns a boolean if a field has been set.

### GetSwitchPassword

`func (o *SitesPutRequestSiteValue) GetSwitchPassword() string`

GetSwitchPassword returns the SwitchPassword field if non-nil, zero value otherwise.

### GetSwitchPasswordOk

`func (o *SitesPutRequestSiteValue) GetSwitchPasswordOk() (*string, bool)`

GetSwitchPasswordOk returns a tuple with the SwitchPassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchPassword

`func (o *SitesPutRequestSiteValue) SetSwitchPassword(v string)`

SetSwitchPassword sets SwitchPassword field to given value.

### HasSwitchPassword

`func (o *SitesPutRequestSiteValue) HasSwitchPassword() bool`

HasSwitchPassword returns a boolean if a field has been set.

### GetSwitchPasswordEncrypted

`func (o *SitesPutRequestSiteValue) GetSwitchPasswordEncrypted() string`

GetSwitchPasswordEncrypted returns the SwitchPasswordEncrypted field if non-nil, zero value otherwise.

### GetSwitchPasswordEncryptedOk

`func (o *SitesPutRequestSiteValue) GetSwitchPasswordEncryptedOk() (*string, bool)`

GetSwitchPasswordEncryptedOk returns a tuple with the SwitchPasswordEncrypted field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchPasswordEncrypted

`func (o *SitesPutRequestSiteValue) SetSwitchPasswordEncrypted(v string)`

SetSwitchPasswordEncrypted sets SwitchPasswordEncrypted field to given value.

### HasSwitchPasswordEncrypted

`func (o *SitesPutRequestSiteValue) HasSwitchPasswordEncrypted() bool`

HasSwitchPasswordEncrypted returns a boolean if a field has been set.

### GetHgxUsername

`func (o *SitesPutRequestSiteValue) GetHgxUsername() string`

GetHgxUsername returns the HgxUsername field if non-nil, zero value otherwise.

### GetHgxUsernameOk

`func (o *SitesPutRequestSiteValue) GetHgxUsernameOk() (*string, bool)`

GetHgxUsernameOk returns a tuple with the HgxUsername field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHgxUsername

`func (o *SitesPutRequestSiteValue) SetHgxUsername(v string)`

SetHgxUsername sets HgxUsername field to given value.

### HasHgxUsername

`func (o *SitesPutRequestSiteValue) HasHgxUsername() bool`

HasHgxUsername returns a boolean if a field has been set.

### GetHgxPassword

`func (o *SitesPutRequestSiteValue) GetHgxPassword() string`

GetHgxPassword returns the HgxPassword field if non-nil, zero value otherwise.

### GetHgxPasswordOk

`func (o *SitesPutRequestSiteValue) GetHgxPasswordOk() (*string, bool)`

GetHgxPasswordOk returns a tuple with the HgxPassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHgxPassword

`func (o *SitesPutRequestSiteValue) SetHgxPassword(v string)`

SetHgxPassword sets HgxPassword field to given value.

### HasHgxPassword

`func (o *SitesPutRequestSiteValue) HasHgxPassword() bool`

HasHgxPassword returns a boolean if a field has been set.

### GetHgxPasswordEncrypted

`func (o *SitesPutRequestSiteValue) GetHgxPasswordEncrypted() string`

GetHgxPasswordEncrypted returns the HgxPasswordEncrypted field if non-nil, zero value otherwise.

### GetHgxPasswordEncryptedOk

`func (o *SitesPutRequestSiteValue) GetHgxPasswordEncryptedOk() (*string, bool)`

GetHgxPasswordEncryptedOk returns a tuple with the HgxPasswordEncrypted field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHgxPasswordEncrypted

`func (o *SitesPutRequestSiteValue) SetHgxPasswordEncrypted(v string)`

SetHgxPasswordEncrypted sets HgxPasswordEncrypted field to given value.

### HasHgxPasswordEncrypted

`func (o *SitesPutRequestSiteValue) HasHgxPasswordEncrypted() bool`

HasHgxPasswordEncrypted returns a boolean if a field has been set.

### GetSwitchGateway

`func (o *SitesPutRequestSiteValue) GetSwitchGateway() string`

GetSwitchGateway returns the SwitchGateway field if non-nil, zero value otherwise.

### GetSwitchGatewayOk

`func (o *SitesPutRequestSiteValue) GetSwitchGatewayOk() (*string, bool)`

GetSwitchGatewayOk returns a tuple with the SwitchGateway field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchGateway

`func (o *SitesPutRequestSiteValue) SetSwitchGateway(v string)`

SetSwitchGateway sets SwitchGateway field to given value.

### HasSwitchGateway

`func (o *SitesPutRequestSiteValue) HasSwitchGateway() bool`

HasSwitchGateway returns a boolean if a field has been set.

### GetControllerGateway

`func (o *SitesPutRequestSiteValue) GetControllerGateway() string`

GetControllerGateway returns the ControllerGateway field if non-nil, zero value otherwise.

### GetControllerGatewayOk

`func (o *SitesPutRequestSiteValue) GetControllerGatewayOk() (*string, bool)`

GetControllerGatewayOk returns a tuple with the ControllerGateway field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetControllerGateway

`func (o *SitesPutRequestSiteValue) SetControllerGateway(v string)`

SetControllerGateway sets ControllerGateway field to given value.

### HasControllerGateway

`func (o *SitesPutRequestSiteValue) HasControllerGateway() bool`

HasControllerGateway returns a boolean if a field has been set.

### GetHgxGateway

`func (o *SitesPutRequestSiteValue) GetHgxGateway() string`

GetHgxGateway returns the HgxGateway field if non-nil, zero value otherwise.

### GetHgxGatewayOk

`func (o *SitesPutRequestSiteValue) GetHgxGatewayOk() (*string, bool)`

GetHgxGatewayOk returns a tuple with the HgxGateway field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHgxGateway

`func (o *SitesPutRequestSiteValue) SetHgxGateway(v string)`

SetHgxGateway sets HgxGateway field to given value.

### HasHgxGateway

`func (o *SitesPutRequestSiteValue) HasHgxGateway() bool`

HasHgxGateway returns a boolean if a field has been set.

### GetIpSourceGuard

`func (o *SitesPutRequestSiteValue) GetIpSourceGuard() bool`

GetIpSourceGuard returns the IpSourceGuard field if non-nil, zero value otherwise.

### GetIpSourceGuardOk

`func (o *SitesPutRequestSiteValue) GetIpSourceGuardOk() (*bool, bool)`

GetIpSourceGuardOk returns a tuple with the IpSourceGuard field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpSourceGuard

`func (o *SitesPutRequestSiteValue) SetIpSourceGuard(v bool)`

SetIpSourceGuard sets IpSourceGuard field to given value.

### HasIpSourceGuard

`func (o *SitesPutRequestSiteValue) HasIpSourceGuard() bool`

HasIpSourceGuard returns a boolean if a field has been set.

### GetEnableDhcpSnooping

`func (o *SitesPutRequestSiteValue) GetEnableDhcpSnooping() bool`

GetEnableDhcpSnooping returns the EnableDhcpSnooping field if non-nil, zero value otherwise.

### GetEnableDhcpSnoopingOk

`func (o *SitesPutRequestSiteValue) GetEnableDhcpSnoopingOk() (*bool, bool)`

GetEnableDhcpSnoopingOk returns a tuple with the EnableDhcpSnooping field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnableDhcpSnooping

`func (o *SitesPutRequestSiteValue) SetEnableDhcpSnooping(v bool)`

SetEnableDhcpSnooping sets EnableDhcpSnooping field to given value.

### HasEnableDhcpSnooping

`func (o *SitesPutRequestSiteValue) HasEnableDhcpSnooping() bool`

HasEnableDhcpSnooping returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


