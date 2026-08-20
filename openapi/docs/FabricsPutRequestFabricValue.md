# FabricsPutRequestFabricValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Template Name. Must be unique within type. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to true]
**SwitchUsername** | Pointer to **string** | Default username for managed switches in this Fabric | [optional] [default to ""]
**SwitchPassword** | Pointer to **string** | Default password for managed switches in this Fabric | [optional] [default to ""]
**SwitchPasswordEncrypted** | Pointer to **string** | Default password for managed switches in this Fabric | [optional] [default to ""]
**HgxUsername** | Pointer to **string** | Default username for HGX devices in this Fabric | [optional] [default to ""]
**HgxPassword** | Pointer to **string** | Default password for HGX devices in this Fabric | [optional] [default to ""]
**HgxPasswordEncrypted** | Pointer to **string** | Default password for HGX devices in this Fabric | [optional] [default to ""]
**SwitchGateway** | Pointer to **string** | Default switch management gateway IP for devices in this Fabric | [optional] [default to ""]
**ControllerGateway** | Pointer to **string** | Default Device Management VM gateway IP for devices in this Fabric | [optional] [default to ""]
**HgxGateway** | Pointer to **string** | Default HGX management gateway IP for devices in this Fabric | [optional] [default to ""]
**PlaneCount** | Pointer to **string** | Number of planes in this Fabric | [optional] [default to "1"]
**SuSize** | Pointer to **string** | Number of HGXs per SU | [optional] [default to "32"]
**SuSupport** | Pointer to **bool** | Support grouping leaf switches in SUs | [optional] [default to false]
**GpuArchitecture** | Pointer to **string** | GPU Architecture used within this Fabric | [optional] [default to "hgx"]
**ServerManagement** | Pointer to **bool** | Support managing servers | [optional] [default to true]
**AllowAllUnderlayConnections** | Pointer to **bool** | Allows underlay connections between PODs | [optional] [default to false]
**FabricType** | Pointer to **string** | Type of Fabric | [optional] [default to "enterprise"]
**DuplicateAddressDetectionMaxNumberOfMoves** | Pointer to **NullableInt64** | Controls duplicate MAC address detection (DAD) Max Number of Moves for EVPN (Ethernet VPN) within the BGP address-family. Number of moves (2 to 1000; default 5 if left blank) | [optional] [default to 5]
**DuplicateAddressDetectionTime** | Pointer to **NullableInt64** | Controls duplicate MAC address detection (DAD) time for EVPN (Ethernet VPN) within the BGP address-family. Time in seconds (2 to 1800; default 180 if left blank) | [optional] [default to 180]
**PortAdminPollingInterval** | Pointer to **NullableInt64** | Polling interval values in seconds, set if aggressive reporting is not enabled | [optional] [default to 0]
**PortStatusPollingInterval** | Pointer to **NullableInt64** | Polling interval values in seconds, set if aggressive reporting is not enabled | [optional] [default to 0]
**ServiceForFabric** | Pointer to **string** | Service for Fabric | [optional] [default to "(predefined):Management"]
**ServiceForFabricRefType** | Pointer to **string** | Object type for service_for_fabric field | [optional] 
**SpanningTreeType** | Pointer to **string** | Sets the spanning tree type for all Ports in this Fabric with Spanning Tree enabled | [optional] [default to "pvst"]
**RegionName** | Pointer to **string** | Defines the logical boundary of the network. All switches in an MSTP region must have the same configured region name | [optional] [default to ""]
**Revision** | Pointer to **NullableInt64** | A logical number that signifies a revision for the MSTP configuration. All switches in an MSTP region must have the same revision number | [optional] [default to 0]
**ForceSpanningTreeOnFabricPorts** | Pointer to **bool** | Enable spanning tree on all fabric connections.  This overrides the Eth Port Settings for Fabric ports | [optional] [default to false]
**ReadOnlyMode** | Pointer to **bool** | When Read Only Mode is checked, vNetC will perform all functions except writing database updates to the target hardware | [optional] [default to false]
**DomainForFabric** | Pointer to **string** | Fabric Collection for Fabric | [optional] [default to ""]
**DomainForFabricRefType** | Pointer to **string** | Object type for domain_for_fabric field | [optional] 
**EnableDscp** | Pointer to **bool** | Enable DSCP to p-bit/TC configuration. When enabled, DSCP to p-bit/TC mappings are applied. | [optional] [default to true]
**DscpToPBitMap** | Pointer to **string** | For any Service that is using DSCP to p-bit map packet prioritization. A string of length 64 with a 0-7 in each position | [optional] [default to "0000000011111111222222223333333344444444555555556666666677777777"]
**AnycastMacAddress** | Pointer to **string** | Fabric Level MAC Address for Anycast | [optional] [default to "(auto)"]
**AnycastMacAddressAutoAssigned** | Pointer to **bool** | Whether or not the value in anycast_mac_address field has been automatically assigned or not. Set to false and change anycast_mac_address value to edit. | [optional] 
**MacAddressAgingTime** | Pointer to **NullableInt64** | MAC Address Aging Time (between 1-100000) | [optional] [default to 600]
**MlagDelayRestoreTimer** | Pointer to **NullableInt64** | MLAG Delay Restore Timer | [optional] [default to 300]
**BgpKeepaliveTimer** | Pointer to **NullableInt64** | Spine BGP Keepalive Timer | [optional] [default to 60]
**BgpHoldDownTimer** | Pointer to **NullableInt64** | Spine BGP Hold Down Timer | [optional] [default to 180]
**SpineBgpAdvertisementInterval** | Pointer to **NullableInt64** | BGP Advertisement Interval for spines/superspines. Use \&quot;0\&quot; for immediate updates | [optional] [default to 1]
**SpineBgpConnectTimer** | Pointer to **NullableInt64** | BGP Connect Timer | [optional] [default to 120]
**SpineAsNumber** | Pointer to **NullableInt64** | BGP AS number applied uniformly to all spine endpoints in this CLOS fabric on save. Leave blank to manage spine AS numbers individually. | [optional] 
**LeafBgpKeepAliveTimer** | Pointer to **NullableInt64** | Leaf BGP Keep Alive Timer | [optional] [default to 60]
**LeafBgpHoldDownTimer** | Pointer to **NullableInt64** | Leaf BGP Hold Down Timer | [optional] [default to 180]
**LeafBgpAdvertisementInterval** | Pointer to **NullableInt64** | BGP Advertisement Interval for leafs. Use \&quot;0\&quot; for immediate updates | [optional] [default to 1]
**LeafBgpConnectTimer** | Pointer to **NullableInt64** | BGP Connect Timer | [optional] [default to 120]
**LinkStateTimeoutValue** | Pointer to **NullableInt64** | Link State Timeout Value | [optional] [default to 60]
**EvpnMultihomingStartupDelay** | Pointer to **NullableInt64** | Startup Delay | [optional] [default to 300]
**EvpnMacHoldtime** | Pointer to **NullableInt64** | MAC Holdtime | [optional] [default to 1080]
**AggressiveReporting** | Pointer to **bool** | Fast Reporting of Switch Communications, Link Up/Down, and BGP Status | [optional] [default to true]
**SwitchIpBase** | Pointer to **string** | Base IPv4 address for switch IPs in this Fabric | [optional] [default to ""]
**ControllerIpBase** | Pointer to **string** | Base IPv4 address for the Device Management VM IPs in this Fabric | [optional] [default to ""]
**MultiTenant** | Pointer to **bool** | Allow multiple tenants to HGX endpoints on this fabric. | [optional] [default to true]
**BaseBgpAsNumber** | Pointer to **string** | Base BGP Autonomous System Number used for switches in the fabric  | [optional] [default to "61000"]
**RouterIdBasePrefix** | Pointer to **string** | Router ID starting IP address  | [optional] [default to "172.16.0.0"]
**VtepIdBasePrefix** | Pointer to **string** | Vtep ID starting IP address  | [optional] [default to "172.16.10.0"]
**PairedIpSubnet** | Pointer to **string** | IP address range reserved for communication between paired switches  | [optional] [default to "192.168.254.0/24"]
**MaxSwitches** | Pointer to **string** | Max number Switches to support in this site  | [optional] [default to "2000"]
**PauseValidationAlarms** | Pointer to **bool** | Validation still runs, but validation alarms are not raised for this Fabric while enabled. | [optional] [default to false]
**StartingOctet** | Pointer to **NullableInt64** | Starting Octet for HGX Port IPs | [optional] 
**MaxSus** | Pointer to **NullableInt64** | Maximum number of SUs allowed per POD | [optional] 
**MaxPods** | Pointer to **NullableInt64** | Maximum number of PODs allowed in the Fabric | [optional] 
**ObjectProperties** | Pointer to [**FabricsPutRequestFabricValueObjectProperties**](FabricsPutRequestFabricValueObjectProperties.md) |  | [optional] 
**SetLeafRouterIdOnBgp** | Pointer to **bool** | Enabling this will use the endpoint loopback0 address as the router ID for all BGP sessions on leaf switches | [optional] [default to false]
**RouteAggregation** | Pointer to **string** | Route Aggregation configuration for this fabric | [optional] [default to ""]
**RouteAggregators** | Pointer to [**[]FabricsPutRequestFabricValueRouteAggregatorsInner**](FabricsPutRequestFabricValueRouteAggregatorsInner.md) |  | [optional] 
**IpSourceGuard** | Pointer to **bool** | On untrusted ports, only allow known traffic from known IP addresses. IP addresses are discovered via DHCP snooping or with static IP settings | [optional] [default to false]
**EnableDhcpSnooping** | Pointer to **bool** | Enables the switches to monitor DHCP traffic and collect assigned IP addresses which are then placed in the DHCP assigned IPs report. | [optional] [default to false]

## Methods

### NewFabricsPutRequestFabricValue

`func NewFabricsPutRequestFabricValue() *FabricsPutRequestFabricValue`

NewFabricsPutRequestFabricValue instantiates a new FabricsPutRequestFabricValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewFabricsPutRequestFabricValueWithDefaults

`func NewFabricsPutRequestFabricValueWithDefaults() *FabricsPutRequestFabricValue`

NewFabricsPutRequestFabricValueWithDefaults instantiates a new FabricsPutRequestFabricValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *FabricsPutRequestFabricValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *FabricsPutRequestFabricValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *FabricsPutRequestFabricValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *FabricsPutRequestFabricValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *FabricsPutRequestFabricValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *FabricsPutRequestFabricValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *FabricsPutRequestFabricValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *FabricsPutRequestFabricValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetSwitchUsername

`func (o *FabricsPutRequestFabricValue) GetSwitchUsername() string`

GetSwitchUsername returns the SwitchUsername field if non-nil, zero value otherwise.

### GetSwitchUsernameOk

`func (o *FabricsPutRequestFabricValue) GetSwitchUsernameOk() (*string, bool)`

GetSwitchUsernameOk returns a tuple with the SwitchUsername field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchUsername

`func (o *FabricsPutRequestFabricValue) SetSwitchUsername(v string)`

SetSwitchUsername sets SwitchUsername field to given value.

### HasSwitchUsername

`func (o *FabricsPutRequestFabricValue) HasSwitchUsername() bool`

HasSwitchUsername returns a boolean if a field has been set.

### GetSwitchPassword

`func (o *FabricsPutRequestFabricValue) GetSwitchPassword() string`

GetSwitchPassword returns the SwitchPassword field if non-nil, zero value otherwise.

### GetSwitchPasswordOk

`func (o *FabricsPutRequestFabricValue) GetSwitchPasswordOk() (*string, bool)`

GetSwitchPasswordOk returns a tuple with the SwitchPassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchPassword

`func (o *FabricsPutRequestFabricValue) SetSwitchPassword(v string)`

SetSwitchPassword sets SwitchPassword field to given value.

### HasSwitchPassword

`func (o *FabricsPutRequestFabricValue) HasSwitchPassword() bool`

HasSwitchPassword returns a boolean if a field has been set.

### GetSwitchPasswordEncrypted

`func (o *FabricsPutRequestFabricValue) GetSwitchPasswordEncrypted() string`

GetSwitchPasswordEncrypted returns the SwitchPasswordEncrypted field if non-nil, zero value otherwise.

### GetSwitchPasswordEncryptedOk

`func (o *FabricsPutRequestFabricValue) GetSwitchPasswordEncryptedOk() (*string, bool)`

GetSwitchPasswordEncryptedOk returns a tuple with the SwitchPasswordEncrypted field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchPasswordEncrypted

`func (o *FabricsPutRequestFabricValue) SetSwitchPasswordEncrypted(v string)`

SetSwitchPasswordEncrypted sets SwitchPasswordEncrypted field to given value.

### HasSwitchPasswordEncrypted

`func (o *FabricsPutRequestFabricValue) HasSwitchPasswordEncrypted() bool`

HasSwitchPasswordEncrypted returns a boolean if a field has been set.

### GetHgxUsername

`func (o *FabricsPutRequestFabricValue) GetHgxUsername() string`

GetHgxUsername returns the HgxUsername field if non-nil, zero value otherwise.

### GetHgxUsernameOk

`func (o *FabricsPutRequestFabricValue) GetHgxUsernameOk() (*string, bool)`

GetHgxUsernameOk returns a tuple with the HgxUsername field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHgxUsername

`func (o *FabricsPutRequestFabricValue) SetHgxUsername(v string)`

SetHgxUsername sets HgxUsername field to given value.

### HasHgxUsername

`func (o *FabricsPutRequestFabricValue) HasHgxUsername() bool`

HasHgxUsername returns a boolean if a field has been set.

### GetHgxPassword

`func (o *FabricsPutRequestFabricValue) GetHgxPassword() string`

GetHgxPassword returns the HgxPassword field if non-nil, zero value otherwise.

### GetHgxPasswordOk

`func (o *FabricsPutRequestFabricValue) GetHgxPasswordOk() (*string, bool)`

GetHgxPasswordOk returns a tuple with the HgxPassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHgxPassword

`func (o *FabricsPutRequestFabricValue) SetHgxPassword(v string)`

SetHgxPassword sets HgxPassword field to given value.

### HasHgxPassword

`func (o *FabricsPutRequestFabricValue) HasHgxPassword() bool`

HasHgxPassword returns a boolean if a field has been set.

### GetHgxPasswordEncrypted

`func (o *FabricsPutRequestFabricValue) GetHgxPasswordEncrypted() string`

GetHgxPasswordEncrypted returns the HgxPasswordEncrypted field if non-nil, zero value otherwise.

### GetHgxPasswordEncryptedOk

`func (o *FabricsPutRequestFabricValue) GetHgxPasswordEncryptedOk() (*string, bool)`

GetHgxPasswordEncryptedOk returns a tuple with the HgxPasswordEncrypted field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHgxPasswordEncrypted

`func (o *FabricsPutRequestFabricValue) SetHgxPasswordEncrypted(v string)`

SetHgxPasswordEncrypted sets HgxPasswordEncrypted field to given value.

### HasHgxPasswordEncrypted

`func (o *FabricsPutRequestFabricValue) HasHgxPasswordEncrypted() bool`

HasHgxPasswordEncrypted returns a boolean if a field has been set.

### GetSwitchGateway

`func (o *FabricsPutRequestFabricValue) GetSwitchGateway() string`

GetSwitchGateway returns the SwitchGateway field if non-nil, zero value otherwise.

### GetSwitchGatewayOk

`func (o *FabricsPutRequestFabricValue) GetSwitchGatewayOk() (*string, bool)`

GetSwitchGatewayOk returns a tuple with the SwitchGateway field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchGateway

`func (o *FabricsPutRequestFabricValue) SetSwitchGateway(v string)`

SetSwitchGateway sets SwitchGateway field to given value.

### HasSwitchGateway

`func (o *FabricsPutRequestFabricValue) HasSwitchGateway() bool`

HasSwitchGateway returns a boolean if a field has been set.

### GetControllerGateway

`func (o *FabricsPutRequestFabricValue) GetControllerGateway() string`

GetControllerGateway returns the ControllerGateway field if non-nil, zero value otherwise.

### GetControllerGatewayOk

`func (o *FabricsPutRequestFabricValue) GetControllerGatewayOk() (*string, bool)`

GetControllerGatewayOk returns a tuple with the ControllerGateway field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetControllerGateway

`func (o *FabricsPutRequestFabricValue) SetControllerGateway(v string)`

SetControllerGateway sets ControllerGateway field to given value.

### HasControllerGateway

`func (o *FabricsPutRequestFabricValue) HasControllerGateway() bool`

HasControllerGateway returns a boolean if a field has been set.

### GetHgxGateway

`func (o *FabricsPutRequestFabricValue) GetHgxGateway() string`

GetHgxGateway returns the HgxGateway field if non-nil, zero value otherwise.

### GetHgxGatewayOk

`func (o *FabricsPutRequestFabricValue) GetHgxGatewayOk() (*string, bool)`

GetHgxGatewayOk returns a tuple with the HgxGateway field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHgxGateway

`func (o *FabricsPutRequestFabricValue) SetHgxGateway(v string)`

SetHgxGateway sets HgxGateway field to given value.

### HasHgxGateway

`func (o *FabricsPutRequestFabricValue) HasHgxGateway() bool`

HasHgxGateway returns a boolean if a field has been set.

### GetPlaneCount

`func (o *FabricsPutRequestFabricValue) GetPlaneCount() string`

GetPlaneCount returns the PlaneCount field if non-nil, zero value otherwise.

### GetPlaneCountOk

`func (o *FabricsPutRequestFabricValue) GetPlaneCountOk() (*string, bool)`

GetPlaneCountOk returns a tuple with the PlaneCount field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPlaneCount

`func (o *FabricsPutRequestFabricValue) SetPlaneCount(v string)`

SetPlaneCount sets PlaneCount field to given value.

### HasPlaneCount

`func (o *FabricsPutRequestFabricValue) HasPlaneCount() bool`

HasPlaneCount returns a boolean if a field has been set.

### GetSuSize

`func (o *FabricsPutRequestFabricValue) GetSuSize() string`

GetSuSize returns the SuSize field if non-nil, zero value otherwise.

### GetSuSizeOk

`func (o *FabricsPutRequestFabricValue) GetSuSizeOk() (*string, bool)`

GetSuSizeOk returns a tuple with the SuSize field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSuSize

`func (o *FabricsPutRequestFabricValue) SetSuSize(v string)`

SetSuSize sets SuSize field to given value.

### HasSuSize

`func (o *FabricsPutRequestFabricValue) HasSuSize() bool`

HasSuSize returns a boolean if a field has been set.

### GetSuSupport

`func (o *FabricsPutRequestFabricValue) GetSuSupport() bool`

GetSuSupport returns the SuSupport field if non-nil, zero value otherwise.

### GetSuSupportOk

`func (o *FabricsPutRequestFabricValue) GetSuSupportOk() (*bool, bool)`

GetSuSupportOk returns a tuple with the SuSupport field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSuSupport

`func (o *FabricsPutRequestFabricValue) SetSuSupport(v bool)`

SetSuSupport sets SuSupport field to given value.

### HasSuSupport

`func (o *FabricsPutRequestFabricValue) HasSuSupport() bool`

HasSuSupport returns a boolean if a field has been set.

### GetGpuArchitecture

`func (o *FabricsPutRequestFabricValue) GetGpuArchitecture() string`

GetGpuArchitecture returns the GpuArchitecture field if non-nil, zero value otherwise.

### GetGpuArchitectureOk

`func (o *FabricsPutRequestFabricValue) GetGpuArchitectureOk() (*string, bool)`

GetGpuArchitectureOk returns a tuple with the GpuArchitecture field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGpuArchitecture

`func (o *FabricsPutRequestFabricValue) SetGpuArchitecture(v string)`

SetGpuArchitecture sets GpuArchitecture field to given value.

### HasGpuArchitecture

`func (o *FabricsPutRequestFabricValue) HasGpuArchitecture() bool`

HasGpuArchitecture returns a boolean if a field has been set.

### GetServerManagement

`func (o *FabricsPutRequestFabricValue) GetServerManagement() bool`

GetServerManagement returns the ServerManagement field if non-nil, zero value otherwise.

### GetServerManagementOk

`func (o *FabricsPutRequestFabricValue) GetServerManagementOk() (*bool, bool)`

GetServerManagementOk returns a tuple with the ServerManagement field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServerManagement

`func (o *FabricsPutRequestFabricValue) SetServerManagement(v bool)`

SetServerManagement sets ServerManagement field to given value.

### HasServerManagement

`func (o *FabricsPutRequestFabricValue) HasServerManagement() bool`

HasServerManagement returns a boolean if a field has been set.

### GetAllowAllUnderlayConnections

`func (o *FabricsPutRequestFabricValue) GetAllowAllUnderlayConnections() bool`

GetAllowAllUnderlayConnections returns the AllowAllUnderlayConnections field if non-nil, zero value otherwise.

### GetAllowAllUnderlayConnectionsOk

`func (o *FabricsPutRequestFabricValue) GetAllowAllUnderlayConnectionsOk() (*bool, bool)`

GetAllowAllUnderlayConnectionsOk returns a tuple with the AllowAllUnderlayConnections field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllowAllUnderlayConnections

`func (o *FabricsPutRequestFabricValue) SetAllowAllUnderlayConnections(v bool)`

SetAllowAllUnderlayConnections sets AllowAllUnderlayConnections field to given value.

### HasAllowAllUnderlayConnections

`func (o *FabricsPutRequestFabricValue) HasAllowAllUnderlayConnections() bool`

HasAllowAllUnderlayConnections returns a boolean if a field has been set.

### GetFabricType

`func (o *FabricsPutRequestFabricValue) GetFabricType() string`

GetFabricType returns the FabricType field if non-nil, zero value otherwise.

### GetFabricTypeOk

`func (o *FabricsPutRequestFabricValue) GetFabricTypeOk() (*string, bool)`

GetFabricTypeOk returns a tuple with the FabricType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFabricType

`func (o *FabricsPutRequestFabricValue) SetFabricType(v string)`

SetFabricType sets FabricType field to given value.

### HasFabricType

`func (o *FabricsPutRequestFabricValue) HasFabricType() bool`

HasFabricType returns a boolean if a field has been set.

### GetDuplicateAddressDetectionMaxNumberOfMoves

`func (o *FabricsPutRequestFabricValue) GetDuplicateAddressDetectionMaxNumberOfMoves() int64`

GetDuplicateAddressDetectionMaxNumberOfMoves returns the DuplicateAddressDetectionMaxNumberOfMoves field if non-nil, zero value otherwise.

### GetDuplicateAddressDetectionMaxNumberOfMovesOk

`func (o *FabricsPutRequestFabricValue) GetDuplicateAddressDetectionMaxNumberOfMovesOk() (*int64, bool)`

GetDuplicateAddressDetectionMaxNumberOfMovesOk returns a tuple with the DuplicateAddressDetectionMaxNumberOfMoves field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDuplicateAddressDetectionMaxNumberOfMoves

`func (o *FabricsPutRequestFabricValue) SetDuplicateAddressDetectionMaxNumberOfMoves(v int64)`

SetDuplicateAddressDetectionMaxNumberOfMoves sets DuplicateAddressDetectionMaxNumberOfMoves field to given value.

### HasDuplicateAddressDetectionMaxNumberOfMoves

`func (o *FabricsPutRequestFabricValue) HasDuplicateAddressDetectionMaxNumberOfMoves() bool`

HasDuplicateAddressDetectionMaxNumberOfMoves returns a boolean if a field has been set.

### SetDuplicateAddressDetectionMaxNumberOfMovesNil

`func (o *FabricsPutRequestFabricValue) SetDuplicateAddressDetectionMaxNumberOfMovesNil(b bool)`

 SetDuplicateAddressDetectionMaxNumberOfMovesNil sets the value for DuplicateAddressDetectionMaxNumberOfMoves to be an explicit nil

### UnsetDuplicateAddressDetectionMaxNumberOfMoves
`func (o *FabricsPutRequestFabricValue) UnsetDuplicateAddressDetectionMaxNumberOfMoves()`

UnsetDuplicateAddressDetectionMaxNumberOfMoves ensures that no value is present for DuplicateAddressDetectionMaxNumberOfMoves, not even an explicit nil
### GetDuplicateAddressDetectionTime

`func (o *FabricsPutRequestFabricValue) GetDuplicateAddressDetectionTime() int64`

GetDuplicateAddressDetectionTime returns the DuplicateAddressDetectionTime field if non-nil, zero value otherwise.

### GetDuplicateAddressDetectionTimeOk

`func (o *FabricsPutRequestFabricValue) GetDuplicateAddressDetectionTimeOk() (*int64, bool)`

GetDuplicateAddressDetectionTimeOk returns a tuple with the DuplicateAddressDetectionTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDuplicateAddressDetectionTime

`func (o *FabricsPutRequestFabricValue) SetDuplicateAddressDetectionTime(v int64)`

SetDuplicateAddressDetectionTime sets DuplicateAddressDetectionTime field to given value.

### HasDuplicateAddressDetectionTime

`func (o *FabricsPutRequestFabricValue) HasDuplicateAddressDetectionTime() bool`

HasDuplicateAddressDetectionTime returns a boolean if a field has been set.

### SetDuplicateAddressDetectionTimeNil

`func (o *FabricsPutRequestFabricValue) SetDuplicateAddressDetectionTimeNil(b bool)`

 SetDuplicateAddressDetectionTimeNil sets the value for DuplicateAddressDetectionTime to be an explicit nil

### UnsetDuplicateAddressDetectionTime
`func (o *FabricsPutRequestFabricValue) UnsetDuplicateAddressDetectionTime()`

UnsetDuplicateAddressDetectionTime ensures that no value is present for DuplicateAddressDetectionTime, not even an explicit nil
### GetPortAdminPollingInterval

`func (o *FabricsPutRequestFabricValue) GetPortAdminPollingInterval() int64`

GetPortAdminPollingInterval returns the PortAdminPollingInterval field if non-nil, zero value otherwise.

### GetPortAdminPollingIntervalOk

`func (o *FabricsPutRequestFabricValue) GetPortAdminPollingIntervalOk() (*int64, bool)`

GetPortAdminPollingIntervalOk returns a tuple with the PortAdminPollingInterval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPortAdminPollingInterval

`func (o *FabricsPutRequestFabricValue) SetPortAdminPollingInterval(v int64)`

SetPortAdminPollingInterval sets PortAdminPollingInterval field to given value.

### HasPortAdminPollingInterval

`func (o *FabricsPutRequestFabricValue) HasPortAdminPollingInterval() bool`

HasPortAdminPollingInterval returns a boolean if a field has been set.

### SetPortAdminPollingIntervalNil

`func (o *FabricsPutRequestFabricValue) SetPortAdminPollingIntervalNil(b bool)`

 SetPortAdminPollingIntervalNil sets the value for PortAdminPollingInterval to be an explicit nil

### UnsetPortAdminPollingInterval
`func (o *FabricsPutRequestFabricValue) UnsetPortAdminPollingInterval()`

UnsetPortAdminPollingInterval ensures that no value is present for PortAdminPollingInterval, not even an explicit nil
### GetPortStatusPollingInterval

`func (o *FabricsPutRequestFabricValue) GetPortStatusPollingInterval() int64`

GetPortStatusPollingInterval returns the PortStatusPollingInterval field if non-nil, zero value otherwise.

### GetPortStatusPollingIntervalOk

`func (o *FabricsPutRequestFabricValue) GetPortStatusPollingIntervalOk() (*int64, bool)`

GetPortStatusPollingIntervalOk returns a tuple with the PortStatusPollingInterval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPortStatusPollingInterval

`func (o *FabricsPutRequestFabricValue) SetPortStatusPollingInterval(v int64)`

SetPortStatusPollingInterval sets PortStatusPollingInterval field to given value.

### HasPortStatusPollingInterval

`func (o *FabricsPutRequestFabricValue) HasPortStatusPollingInterval() bool`

HasPortStatusPollingInterval returns a boolean if a field has been set.

### SetPortStatusPollingIntervalNil

`func (o *FabricsPutRequestFabricValue) SetPortStatusPollingIntervalNil(b bool)`

 SetPortStatusPollingIntervalNil sets the value for PortStatusPollingInterval to be an explicit nil

### UnsetPortStatusPollingInterval
`func (o *FabricsPutRequestFabricValue) UnsetPortStatusPollingInterval()`

UnsetPortStatusPollingInterval ensures that no value is present for PortStatusPollingInterval, not even an explicit nil
### GetServiceForFabric

`func (o *FabricsPutRequestFabricValue) GetServiceForFabric() string`

GetServiceForFabric returns the ServiceForFabric field if non-nil, zero value otherwise.

### GetServiceForFabricOk

`func (o *FabricsPutRequestFabricValue) GetServiceForFabricOk() (*string, bool)`

GetServiceForFabricOk returns a tuple with the ServiceForFabric field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServiceForFabric

`func (o *FabricsPutRequestFabricValue) SetServiceForFabric(v string)`

SetServiceForFabric sets ServiceForFabric field to given value.

### HasServiceForFabric

`func (o *FabricsPutRequestFabricValue) HasServiceForFabric() bool`

HasServiceForFabric returns a boolean if a field has been set.

### GetServiceForFabricRefType

`func (o *FabricsPutRequestFabricValue) GetServiceForFabricRefType() string`

GetServiceForFabricRefType returns the ServiceForFabricRefType field if non-nil, zero value otherwise.

### GetServiceForFabricRefTypeOk

`func (o *FabricsPutRequestFabricValue) GetServiceForFabricRefTypeOk() (*string, bool)`

GetServiceForFabricRefTypeOk returns a tuple with the ServiceForFabricRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServiceForFabricRefType

`func (o *FabricsPutRequestFabricValue) SetServiceForFabricRefType(v string)`

SetServiceForFabricRefType sets ServiceForFabricRefType field to given value.

### HasServiceForFabricRefType

`func (o *FabricsPutRequestFabricValue) HasServiceForFabricRefType() bool`

HasServiceForFabricRefType returns a boolean if a field has been set.

### GetSpanningTreeType

`func (o *FabricsPutRequestFabricValue) GetSpanningTreeType() string`

GetSpanningTreeType returns the SpanningTreeType field if non-nil, zero value otherwise.

### GetSpanningTreeTypeOk

`func (o *FabricsPutRequestFabricValue) GetSpanningTreeTypeOk() (*string, bool)`

GetSpanningTreeTypeOk returns a tuple with the SpanningTreeType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpanningTreeType

`func (o *FabricsPutRequestFabricValue) SetSpanningTreeType(v string)`

SetSpanningTreeType sets SpanningTreeType field to given value.

### HasSpanningTreeType

`func (o *FabricsPutRequestFabricValue) HasSpanningTreeType() bool`

HasSpanningTreeType returns a boolean if a field has been set.

### GetRegionName

`func (o *FabricsPutRequestFabricValue) GetRegionName() string`

GetRegionName returns the RegionName field if non-nil, zero value otherwise.

### GetRegionNameOk

`func (o *FabricsPutRequestFabricValue) GetRegionNameOk() (*string, bool)`

GetRegionNameOk returns a tuple with the RegionName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRegionName

`func (o *FabricsPutRequestFabricValue) SetRegionName(v string)`

SetRegionName sets RegionName field to given value.

### HasRegionName

`func (o *FabricsPutRequestFabricValue) HasRegionName() bool`

HasRegionName returns a boolean if a field has been set.

### GetRevision

`func (o *FabricsPutRequestFabricValue) GetRevision() int64`

GetRevision returns the Revision field if non-nil, zero value otherwise.

### GetRevisionOk

`func (o *FabricsPutRequestFabricValue) GetRevisionOk() (*int64, bool)`

GetRevisionOk returns a tuple with the Revision field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRevision

`func (o *FabricsPutRequestFabricValue) SetRevision(v int64)`

SetRevision sets Revision field to given value.

### HasRevision

`func (o *FabricsPutRequestFabricValue) HasRevision() bool`

HasRevision returns a boolean if a field has been set.

### SetRevisionNil

`func (o *FabricsPutRequestFabricValue) SetRevisionNil(b bool)`

 SetRevisionNil sets the value for Revision to be an explicit nil

### UnsetRevision
`func (o *FabricsPutRequestFabricValue) UnsetRevision()`

UnsetRevision ensures that no value is present for Revision, not even an explicit nil
### GetForceSpanningTreeOnFabricPorts

`func (o *FabricsPutRequestFabricValue) GetForceSpanningTreeOnFabricPorts() bool`

GetForceSpanningTreeOnFabricPorts returns the ForceSpanningTreeOnFabricPorts field if non-nil, zero value otherwise.

### GetForceSpanningTreeOnFabricPortsOk

`func (o *FabricsPutRequestFabricValue) GetForceSpanningTreeOnFabricPortsOk() (*bool, bool)`

GetForceSpanningTreeOnFabricPortsOk returns a tuple with the ForceSpanningTreeOnFabricPorts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetForceSpanningTreeOnFabricPorts

`func (o *FabricsPutRequestFabricValue) SetForceSpanningTreeOnFabricPorts(v bool)`

SetForceSpanningTreeOnFabricPorts sets ForceSpanningTreeOnFabricPorts field to given value.

### HasForceSpanningTreeOnFabricPorts

`func (o *FabricsPutRequestFabricValue) HasForceSpanningTreeOnFabricPorts() bool`

HasForceSpanningTreeOnFabricPorts returns a boolean if a field has been set.

### GetReadOnlyMode

`func (o *FabricsPutRequestFabricValue) GetReadOnlyMode() bool`

GetReadOnlyMode returns the ReadOnlyMode field if non-nil, zero value otherwise.

### GetReadOnlyModeOk

`func (o *FabricsPutRequestFabricValue) GetReadOnlyModeOk() (*bool, bool)`

GetReadOnlyModeOk returns a tuple with the ReadOnlyMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReadOnlyMode

`func (o *FabricsPutRequestFabricValue) SetReadOnlyMode(v bool)`

SetReadOnlyMode sets ReadOnlyMode field to given value.

### HasReadOnlyMode

`func (o *FabricsPutRequestFabricValue) HasReadOnlyMode() bool`

HasReadOnlyMode returns a boolean if a field has been set.

### GetDomainForFabric

`func (o *FabricsPutRequestFabricValue) GetDomainForFabric() string`

GetDomainForFabric returns the DomainForFabric field if non-nil, zero value otherwise.

### GetDomainForFabricOk

`func (o *FabricsPutRequestFabricValue) GetDomainForFabricOk() (*string, bool)`

GetDomainForFabricOk returns a tuple with the DomainForFabric field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDomainForFabric

`func (o *FabricsPutRequestFabricValue) SetDomainForFabric(v string)`

SetDomainForFabric sets DomainForFabric field to given value.

### HasDomainForFabric

`func (o *FabricsPutRequestFabricValue) HasDomainForFabric() bool`

HasDomainForFabric returns a boolean if a field has been set.

### GetDomainForFabricRefType

`func (o *FabricsPutRequestFabricValue) GetDomainForFabricRefType() string`

GetDomainForFabricRefType returns the DomainForFabricRefType field if non-nil, zero value otherwise.

### GetDomainForFabricRefTypeOk

`func (o *FabricsPutRequestFabricValue) GetDomainForFabricRefTypeOk() (*string, bool)`

GetDomainForFabricRefTypeOk returns a tuple with the DomainForFabricRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDomainForFabricRefType

`func (o *FabricsPutRequestFabricValue) SetDomainForFabricRefType(v string)`

SetDomainForFabricRefType sets DomainForFabricRefType field to given value.

### HasDomainForFabricRefType

`func (o *FabricsPutRequestFabricValue) HasDomainForFabricRefType() bool`

HasDomainForFabricRefType returns a boolean if a field has been set.

### GetEnableDscp

`func (o *FabricsPutRequestFabricValue) GetEnableDscp() bool`

GetEnableDscp returns the EnableDscp field if non-nil, zero value otherwise.

### GetEnableDscpOk

`func (o *FabricsPutRequestFabricValue) GetEnableDscpOk() (*bool, bool)`

GetEnableDscpOk returns a tuple with the EnableDscp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnableDscp

`func (o *FabricsPutRequestFabricValue) SetEnableDscp(v bool)`

SetEnableDscp sets EnableDscp field to given value.

### HasEnableDscp

`func (o *FabricsPutRequestFabricValue) HasEnableDscp() bool`

HasEnableDscp returns a boolean if a field has been set.

### GetDscpToPBitMap

`func (o *FabricsPutRequestFabricValue) GetDscpToPBitMap() string`

GetDscpToPBitMap returns the DscpToPBitMap field if non-nil, zero value otherwise.

### GetDscpToPBitMapOk

`func (o *FabricsPutRequestFabricValue) GetDscpToPBitMapOk() (*string, bool)`

GetDscpToPBitMapOk returns a tuple with the DscpToPBitMap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDscpToPBitMap

`func (o *FabricsPutRequestFabricValue) SetDscpToPBitMap(v string)`

SetDscpToPBitMap sets DscpToPBitMap field to given value.

### HasDscpToPBitMap

`func (o *FabricsPutRequestFabricValue) HasDscpToPBitMap() bool`

HasDscpToPBitMap returns a boolean if a field has been set.

### GetAnycastMacAddress

`func (o *FabricsPutRequestFabricValue) GetAnycastMacAddress() string`

GetAnycastMacAddress returns the AnycastMacAddress field if non-nil, zero value otherwise.

### GetAnycastMacAddressOk

`func (o *FabricsPutRequestFabricValue) GetAnycastMacAddressOk() (*string, bool)`

GetAnycastMacAddressOk returns a tuple with the AnycastMacAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnycastMacAddress

`func (o *FabricsPutRequestFabricValue) SetAnycastMacAddress(v string)`

SetAnycastMacAddress sets AnycastMacAddress field to given value.

### HasAnycastMacAddress

`func (o *FabricsPutRequestFabricValue) HasAnycastMacAddress() bool`

HasAnycastMacAddress returns a boolean if a field has been set.

### GetAnycastMacAddressAutoAssigned

`func (o *FabricsPutRequestFabricValue) GetAnycastMacAddressAutoAssigned() bool`

GetAnycastMacAddressAutoAssigned returns the AnycastMacAddressAutoAssigned field if non-nil, zero value otherwise.

### GetAnycastMacAddressAutoAssignedOk

`func (o *FabricsPutRequestFabricValue) GetAnycastMacAddressAutoAssignedOk() (*bool, bool)`

GetAnycastMacAddressAutoAssignedOk returns a tuple with the AnycastMacAddressAutoAssigned field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnycastMacAddressAutoAssigned

`func (o *FabricsPutRequestFabricValue) SetAnycastMacAddressAutoAssigned(v bool)`

SetAnycastMacAddressAutoAssigned sets AnycastMacAddressAutoAssigned field to given value.

### HasAnycastMacAddressAutoAssigned

`func (o *FabricsPutRequestFabricValue) HasAnycastMacAddressAutoAssigned() bool`

HasAnycastMacAddressAutoAssigned returns a boolean if a field has been set.

### GetMacAddressAgingTime

`func (o *FabricsPutRequestFabricValue) GetMacAddressAgingTime() int64`

GetMacAddressAgingTime returns the MacAddressAgingTime field if non-nil, zero value otherwise.

### GetMacAddressAgingTimeOk

`func (o *FabricsPutRequestFabricValue) GetMacAddressAgingTimeOk() (*int64, bool)`

GetMacAddressAgingTimeOk returns a tuple with the MacAddressAgingTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMacAddressAgingTime

`func (o *FabricsPutRequestFabricValue) SetMacAddressAgingTime(v int64)`

SetMacAddressAgingTime sets MacAddressAgingTime field to given value.

### HasMacAddressAgingTime

`func (o *FabricsPutRequestFabricValue) HasMacAddressAgingTime() bool`

HasMacAddressAgingTime returns a boolean if a field has been set.

### SetMacAddressAgingTimeNil

`func (o *FabricsPutRequestFabricValue) SetMacAddressAgingTimeNil(b bool)`

 SetMacAddressAgingTimeNil sets the value for MacAddressAgingTime to be an explicit nil

### UnsetMacAddressAgingTime
`func (o *FabricsPutRequestFabricValue) UnsetMacAddressAgingTime()`

UnsetMacAddressAgingTime ensures that no value is present for MacAddressAgingTime, not even an explicit nil
### GetMlagDelayRestoreTimer

`func (o *FabricsPutRequestFabricValue) GetMlagDelayRestoreTimer() int64`

GetMlagDelayRestoreTimer returns the MlagDelayRestoreTimer field if non-nil, zero value otherwise.

### GetMlagDelayRestoreTimerOk

`func (o *FabricsPutRequestFabricValue) GetMlagDelayRestoreTimerOk() (*int64, bool)`

GetMlagDelayRestoreTimerOk returns a tuple with the MlagDelayRestoreTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMlagDelayRestoreTimer

`func (o *FabricsPutRequestFabricValue) SetMlagDelayRestoreTimer(v int64)`

SetMlagDelayRestoreTimer sets MlagDelayRestoreTimer field to given value.

### HasMlagDelayRestoreTimer

`func (o *FabricsPutRequestFabricValue) HasMlagDelayRestoreTimer() bool`

HasMlagDelayRestoreTimer returns a boolean if a field has been set.

### SetMlagDelayRestoreTimerNil

`func (o *FabricsPutRequestFabricValue) SetMlagDelayRestoreTimerNil(b bool)`

 SetMlagDelayRestoreTimerNil sets the value for MlagDelayRestoreTimer to be an explicit nil

### UnsetMlagDelayRestoreTimer
`func (o *FabricsPutRequestFabricValue) UnsetMlagDelayRestoreTimer()`

UnsetMlagDelayRestoreTimer ensures that no value is present for MlagDelayRestoreTimer, not even an explicit nil
### GetBgpKeepaliveTimer

`func (o *FabricsPutRequestFabricValue) GetBgpKeepaliveTimer() int64`

GetBgpKeepaliveTimer returns the BgpKeepaliveTimer field if non-nil, zero value otherwise.

### GetBgpKeepaliveTimerOk

`func (o *FabricsPutRequestFabricValue) GetBgpKeepaliveTimerOk() (*int64, bool)`

GetBgpKeepaliveTimerOk returns a tuple with the BgpKeepaliveTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBgpKeepaliveTimer

`func (o *FabricsPutRequestFabricValue) SetBgpKeepaliveTimer(v int64)`

SetBgpKeepaliveTimer sets BgpKeepaliveTimer field to given value.

### HasBgpKeepaliveTimer

`func (o *FabricsPutRequestFabricValue) HasBgpKeepaliveTimer() bool`

HasBgpKeepaliveTimer returns a boolean if a field has been set.

### SetBgpKeepaliveTimerNil

`func (o *FabricsPutRequestFabricValue) SetBgpKeepaliveTimerNil(b bool)`

 SetBgpKeepaliveTimerNil sets the value for BgpKeepaliveTimer to be an explicit nil

### UnsetBgpKeepaliveTimer
`func (o *FabricsPutRequestFabricValue) UnsetBgpKeepaliveTimer()`

UnsetBgpKeepaliveTimer ensures that no value is present for BgpKeepaliveTimer, not even an explicit nil
### GetBgpHoldDownTimer

`func (o *FabricsPutRequestFabricValue) GetBgpHoldDownTimer() int64`

GetBgpHoldDownTimer returns the BgpHoldDownTimer field if non-nil, zero value otherwise.

### GetBgpHoldDownTimerOk

`func (o *FabricsPutRequestFabricValue) GetBgpHoldDownTimerOk() (*int64, bool)`

GetBgpHoldDownTimerOk returns a tuple with the BgpHoldDownTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBgpHoldDownTimer

`func (o *FabricsPutRequestFabricValue) SetBgpHoldDownTimer(v int64)`

SetBgpHoldDownTimer sets BgpHoldDownTimer field to given value.

### HasBgpHoldDownTimer

`func (o *FabricsPutRequestFabricValue) HasBgpHoldDownTimer() bool`

HasBgpHoldDownTimer returns a boolean if a field has been set.

### SetBgpHoldDownTimerNil

`func (o *FabricsPutRequestFabricValue) SetBgpHoldDownTimerNil(b bool)`

 SetBgpHoldDownTimerNil sets the value for BgpHoldDownTimer to be an explicit nil

### UnsetBgpHoldDownTimer
`func (o *FabricsPutRequestFabricValue) UnsetBgpHoldDownTimer()`

UnsetBgpHoldDownTimer ensures that no value is present for BgpHoldDownTimer, not even an explicit nil
### GetSpineBgpAdvertisementInterval

`func (o *FabricsPutRequestFabricValue) GetSpineBgpAdvertisementInterval() int64`

GetSpineBgpAdvertisementInterval returns the SpineBgpAdvertisementInterval field if non-nil, zero value otherwise.

### GetSpineBgpAdvertisementIntervalOk

`func (o *FabricsPutRequestFabricValue) GetSpineBgpAdvertisementIntervalOk() (*int64, bool)`

GetSpineBgpAdvertisementIntervalOk returns a tuple with the SpineBgpAdvertisementInterval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpineBgpAdvertisementInterval

`func (o *FabricsPutRequestFabricValue) SetSpineBgpAdvertisementInterval(v int64)`

SetSpineBgpAdvertisementInterval sets SpineBgpAdvertisementInterval field to given value.

### HasSpineBgpAdvertisementInterval

`func (o *FabricsPutRequestFabricValue) HasSpineBgpAdvertisementInterval() bool`

HasSpineBgpAdvertisementInterval returns a boolean if a field has been set.

### SetSpineBgpAdvertisementIntervalNil

`func (o *FabricsPutRequestFabricValue) SetSpineBgpAdvertisementIntervalNil(b bool)`

 SetSpineBgpAdvertisementIntervalNil sets the value for SpineBgpAdvertisementInterval to be an explicit nil

### UnsetSpineBgpAdvertisementInterval
`func (o *FabricsPutRequestFabricValue) UnsetSpineBgpAdvertisementInterval()`

UnsetSpineBgpAdvertisementInterval ensures that no value is present for SpineBgpAdvertisementInterval, not even an explicit nil
### GetSpineBgpConnectTimer

`func (o *FabricsPutRequestFabricValue) GetSpineBgpConnectTimer() int64`

GetSpineBgpConnectTimer returns the SpineBgpConnectTimer field if non-nil, zero value otherwise.

### GetSpineBgpConnectTimerOk

`func (o *FabricsPutRequestFabricValue) GetSpineBgpConnectTimerOk() (*int64, bool)`

GetSpineBgpConnectTimerOk returns a tuple with the SpineBgpConnectTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpineBgpConnectTimer

`func (o *FabricsPutRequestFabricValue) SetSpineBgpConnectTimer(v int64)`

SetSpineBgpConnectTimer sets SpineBgpConnectTimer field to given value.

### HasSpineBgpConnectTimer

`func (o *FabricsPutRequestFabricValue) HasSpineBgpConnectTimer() bool`

HasSpineBgpConnectTimer returns a boolean if a field has been set.

### SetSpineBgpConnectTimerNil

`func (o *FabricsPutRequestFabricValue) SetSpineBgpConnectTimerNil(b bool)`

 SetSpineBgpConnectTimerNil sets the value for SpineBgpConnectTimer to be an explicit nil

### UnsetSpineBgpConnectTimer
`func (o *FabricsPutRequestFabricValue) UnsetSpineBgpConnectTimer()`

UnsetSpineBgpConnectTimer ensures that no value is present for SpineBgpConnectTimer, not even an explicit nil
### GetSpineAsNumber

`func (o *FabricsPutRequestFabricValue) GetSpineAsNumber() int64`

GetSpineAsNumber returns the SpineAsNumber field if non-nil, zero value otherwise.

### GetSpineAsNumberOk

`func (o *FabricsPutRequestFabricValue) GetSpineAsNumberOk() (*int64, bool)`

GetSpineAsNumberOk returns a tuple with the SpineAsNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpineAsNumber

`func (o *FabricsPutRequestFabricValue) SetSpineAsNumber(v int64)`

SetSpineAsNumber sets SpineAsNumber field to given value.

### HasSpineAsNumber

`func (o *FabricsPutRequestFabricValue) HasSpineAsNumber() bool`

HasSpineAsNumber returns a boolean if a field has been set.

### SetSpineAsNumberNil

`func (o *FabricsPutRequestFabricValue) SetSpineAsNumberNil(b bool)`

 SetSpineAsNumberNil sets the value for SpineAsNumber to be an explicit nil

### UnsetSpineAsNumber
`func (o *FabricsPutRequestFabricValue) UnsetSpineAsNumber()`

UnsetSpineAsNumber ensures that no value is present for SpineAsNumber, not even an explicit nil
### GetLeafBgpKeepAliveTimer

`func (o *FabricsPutRequestFabricValue) GetLeafBgpKeepAliveTimer() int64`

GetLeafBgpKeepAliveTimer returns the LeafBgpKeepAliveTimer field if non-nil, zero value otherwise.

### GetLeafBgpKeepAliveTimerOk

`func (o *FabricsPutRequestFabricValue) GetLeafBgpKeepAliveTimerOk() (*int64, bool)`

GetLeafBgpKeepAliveTimerOk returns a tuple with the LeafBgpKeepAliveTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLeafBgpKeepAliveTimer

`func (o *FabricsPutRequestFabricValue) SetLeafBgpKeepAliveTimer(v int64)`

SetLeafBgpKeepAliveTimer sets LeafBgpKeepAliveTimer field to given value.

### HasLeafBgpKeepAliveTimer

`func (o *FabricsPutRequestFabricValue) HasLeafBgpKeepAliveTimer() bool`

HasLeafBgpKeepAliveTimer returns a boolean if a field has been set.

### SetLeafBgpKeepAliveTimerNil

`func (o *FabricsPutRequestFabricValue) SetLeafBgpKeepAliveTimerNil(b bool)`

 SetLeafBgpKeepAliveTimerNil sets the value for LeafBgpKeepAliveTimer to be an explicit nil

### UnsetLeafBgpKeepAliveTimer
`func (o *FabricsPutRequestFabricValue) UnsetLeafBgpKeepAliveTimer()`

UnsetLeafBgpKeepAliveTimer ensures that no value is present for LeafBgpKeepAliveTimer, not even an explicit nil
### GetLeafBgpHoldDownTimer

`func (o *FabricsPutRequestFabricValue) GetLeafBgpHoldDownTimer() int64`

GetLeafBgpHoldDownTimer returns the LeafBgpHoldDownTimer field if non-nil, zero value otherwise.

### GetLeafBgpHoldDownTimerOk

`func (o *FabricsPutRequestFabricValue) GetLeafBgpHoldDownTimerOk() (*int64, bool)`

GetLeafBgpHoldDownTimerOk returns a tuple with the LeafBgpHoldDownTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLeafBgpHoldDownTimer

`func (o *FabricsPutRequestFabricValue) SetLeafBgpHoldDownTimer(v int64)`

SetLeafBgpHoldDownTimer sets LeafBgpHoldDownTimer field to given value.

### HasLeafBgpHoldDownTimer

`func (o *FabricsPutRequestFabricValue) HasLeafBgpHoldDownTimer() bool`

HasLeafBgpHoldDownTimer returns a boolean if a field has been set.

### SetLeafBgpHoldDownTimerNil

`func (o *FabricsPutRequestFabricValue) SetLeafBgpHoldDownTimerNil(b bool)`

 SetLeafBgpHoldDownTimerNil sets the value for LeafBgpHoldDownTimer to be an explicit nil

### UnsetLeafBgpHoldDownTimer
`func (o *FabricsPutRequestFabricValue) UnsetLeafBgpHoldDownTimer()`

UnsetLeafBgpHoldDownTimer ensures that no value is present for LeafBgpHoldDownTimer, not even an explicit nil
### GetLeafBgpAdvertisementInterval

`func (o *FabricsPutRequestFabricValue) GetLeafBgpAdvertisementInterval() int64`

GetLeafBgpAdvertisementInterval returns the LeafBgpAdvertisementInterval field if non-nil, zero value otherwise.

### GetLeafBgpAdvertisementIntervalOk

`func (o *FabricsPutRequestFabricValue) GetLeafBgpAdvertisementIntervalOk() (*int64, bool)`

GetLeafBgpAdvertisementIntervalOk returns a tuple with the LeafBgpAdvertisementInterval field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLeafBgpAdvertisementInterval

`func (o *FabricsPutRequestFabricValue) SetLeafBgpAdvertisementInterval(v int64)`

SetLeafBgpAdvertisementInterval sets LeafBgpAdvertisementInterval field to given value.

### HasLeafBgpAdvertisementInterval

`func (o *FabricsPutRequestFabricValue) HasLeafBgpAdvertisementInterval() bool`

HasLeafBgpAdvertisementInterval returns a boolean if a field has been set.

### SetLeafBgpAdvertisementIntervalNil

`func (o *FabricsPutRequestFabricValue) SetLeafBgpAdvertisementIntervalNil(b bool)`

 SetLeafBgpAdvertisementIntervalNil sets the value for LeafBgpAdvertisementInterval to be an explicit nil

### UnsetLeafBgpAdvertisementInterval
`func (o *FabricsPutRequestFabricValue) UnsetLeafBgpAdvertisementInterval()`

UnsetLeafBgpAdvertisementInterval ensures that no value is present for LeafBgpAdvertisementInterval, not even an explicit nil
### GetLeafBgpConnectTimer

`func (o *FabricsPutRequestFabricValue) GetLeafBgpConnectTimer() int64`

GetLeafBgpConnectTimer returns the LeafBgpConnectTimer field if non-nil, zero value otherwise.

### GetLeafBgpConnectTimerOk

`func (o *FabricsPutRequestFabricValue) GetLeafBgpConnectTimerOk() (*int64, bool)`

GetLeafBgpConnectTimerOk returns a tuple with the LeafBgpConnectTimer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLeafBgpConnectTimer

`func (o *FabricsPutRequestFabricValue) SetLeafBgpConnectTimer(v int64)`

SetLeafBgpConnectTimer sets LeafBgpConnectTimer field to given value.

### HasLeafBgpConnectTimer

`func (o *FabricsPutRequestFabricValue) HasLeafBgpConnectTimer() bool`

HasLeafBgpConnectTimer returns a boolean if a field has been set.

### SetLeafBgpConnectTimerNil

`func (o *FabricsPutRequestFabricValue) SetLeafBgpConnectTimerNil(b bool)`

 SetLeafBgpConnectTimerNil sets the value for LeafBgpConnectTimer to be an explicit nil

### UnsetLeafBgpConnectTimer
`func (o *FabricsPutRequestFabricValue) UnsetLeafBgpConnectTimer()`

UnsetLeafBgpConnectTimer ensures that no value is present for LeafBgpConnectTimer, not even an explicit nil
### GetLinkStateTimeoutValue

`func (o *FabricsPutRequestFabricValue) GetLinkStateTimeoutValue() int64`

GetLinkStateTimeoutValue returns the LinkStateTimeoutValue field if non-nil, zero value otherwise.

### GetLinkStateTimeoutValueOk

`func (o *FabricsPutRequestFabricValue) GetLinkStateTimeoutValueOk() (*int64, bool)`

GetLinkStateTimeoutValueOk returns a tuple with the LinkStateTimeoutValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLinkStateTimeoutValue

`func (o *FabricsPutRequestFabricValue) SetLinkStateTimeoutValue(v int64)`

SetLinkStateTimeoutValue sets LinkStateTimeoutValue field to given value.

### HasLinkStateTimeoutValue

`func (o *FabricsPutRequestFabricValue) HasLinkStateTimeoutValue() bool`

HasLinkStateTimeoutValue returns a boolean if a field has been set.

### SetLinkStateTimeoutValueNil

`func (o *FabricsPutRequestFabricValue) SetLinkStateTimeoutValueNil(b bool)`

 SetLinkStateTimeoutValueNil sets the value for LinkStateTimeoutValue to be an explicit nil

### UnsetLinkStateTimeoutValue
`func (o *FabricsPutRequestFabricValue) UnsetLinkStateTimeoutValue()`

UnsetLinkStateTimeoutValue ensures that no value is present for LinkStateTimeoutValue, not even an explicit nil
### GetEvpnMultihomingStartupDelay

`func (o *FabricsPutRequestFabricValue) GetEvpnMultihomingStartupDelay() int64`

GetEvpnMultihomingStartupDelay returns the EvpnMultihomingStartupDelay field if non-nil, zero value otherwise.

### GetEvpnMultihomingStartupDelayOk

`func (o *FabricsPutRequestFabricValue) GetEvpnMultihomingStartupDelayOk() (*int64, bool)`

GetEvpnMultihomingStartupDelayOk returns a tuple with the EvpnMultihomingStartupDelay field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEvpnMultihomingStartupDelay

`func (o *FabricsPutRequestFabricValue) SetEvpnMultihomingStartupDelay(v int64)`

SetEvpnMultihomingStartupDelay sets EvpnMultihomingStartupDelay field to given value.

### HasEvpnMultihomingStartupDelay

`func (o *FabricsPutRequestFabricValue) HasEvpnMultihomingStartupDelay() bool`

HasEvpnMultihomingStartupDelay returns a boolean if a field has been set.

### SetEvpnMultihomingStartupDelayNil

`func (o *FabricsPutRequestFabricValue) SetEvpnMultihomingStartupDelayNil(b bool)`

 SetEvpnMultihomingStartupDelayNil sets the value for EvpnMultihomingStartupDelay to be an explicit nil

### UnsetEvpnMultihomingStartupDelay
`func (o *FabricsPutRequestFabricValue) UnsetEvpnMultihomingStartupDelay()`

UnsetEvpnMultihomingStartupDelay ensures that no value is present for EvpnMultihomingStartupDelay, not even an explicit nil
### GetEvpnMacHoldtime

`func (o *FabricsPutRequestFabricValue) GetEvpnMacHoldtime() int64`

GetEvpnMacHoldtime returns the EvpnMacHoldtime field if non-nil, zero value otherwise.

### GetEvpnMacHoldtimeOk

`func (o *FabricsPutRequestFabricValue) GetEvpnMacHoldtimeOk() (*int64, bool)`

GetEvpnMacHoldtimeOk returns a tuple with the EvpnMacHoldtime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEvpnMacHoldtime

`func (o *FabricsPutRequestFabricValue) SetEvpnMacHoldtime(v int64)`

SetEvpnMacHoldtime sets EvpnMacHoldtime field to given value.

### HasEvpnMacHoldtime

`func (o *FabricsPutRequestFabricValue) HasEvpnMacHoldtime() bool`

HasEvpnMacHoldtime returns a boolean if a field has been set.

### SetEvpnMacHoldtimeNil

`func (o *FabricsPutRequestFabricValue) SetEvpnMacHoldtimeNil(b bool)`

 SetEvpnMacHoldtimeNil sets the value for EvpnMacHoldtime to be an explicit nil

### UnsetEvpnMacHoldtime
`func (o *FabricsPutRequestFabricValue) UnsetEvpnMacHoldtime()`

UnsetEvpnMacHoldtime ensures that no value is present for EvpnMacHoldtime, not even an explicit nil
### GetAggressiveReporting

`func (o *FabricsPutRequestFabricValue) GetAggressiveReporting() bool`

GetAggressiveReporting returns the AggressiveReporting field if non-nil, zero value otherwise.

### GetAggressiveReportingOk

`func (o *FabricsPutRequestFabricValue) GetAggressiveReportingOk() (*bool, bool)`

GetAggressiveReportingOk returns a tuple with the AggressiveReporting field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAggressiveReporting

`func (o *FabricsPutRequestFabricValue) SetAggressiveReporting(v bool)`

SetAggressiveReporting sets AggressiveReporting field to given value.

### HasAggressiveReporting

`func (o *FabricsPutRequestFabricValue) HasAggressiveReporting() bool`

HasAggressiveReporting returns a boolean if a field has been set.

### GetSwitchIpBase

`func (o *FabricsPutRequestFabricValue) GetSwitchIpBase() string`

GetSwitchIpBase returns the SwitchIpBase field if non-nil, zero value otherwise.

### GetSwitchIpBaseOk

`func (o *FabricsPutRequestFabricValue) GetSwitchIpBaseOk() (*string, bool)`

GetSwitchIpBaseOk returns a tuple with the SwitchIpBase field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchIpBase

`func (o *FabricsPutRequestFabricValue) SetSwitchIpBase(v string)`

SetSwitchIpBase sets SwitchIpBase field to given value.

### HasSwitchIpBase

`func (o *FabricsPutRequestFabricValue) HasSwitchIpBase() bool`

HasSwitchIpBase returns a boolean if a field has been set.

### GetControllerIpBase

`func (o *FabricsPutRequestFabricValue) GetControllerIpBase() string`

GetControllerIpBase returns the ControllerIpBase field if non-nil, zero value otherwise.

### GetControllerIpBaseOk

`func (o *FabricsPutRequestFabricValue) GetControllerIpBaseOk() (*string, bool)`

GetControllerIpBaseOk returns a tuple with the ControllerIpBase field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetControllerIpBase

`func (o *FabricsPutRequestFabricValue) SetControllerIpBase(v string)`

SetControllerIpBase sets ControllerIpBase field to given value.

### HasControllerIpBase

`func (o *FabricsPutRequestFabricValue) HasControllerIpBase() bool`

HasControllerIpBase returns a boolean if a field has been set.

### GetMultiTenant

`func (o *FabricsPutRequestFabricValue) GetMultiTenant() bool`

GetMultiTenant returns the MultiTenant field if non-nil, zero value otherwise.

### GetMultiTenantOk

`func (o *FabricsPutRequestFabricValue) GetMultiTenantOk() (*bool, bool)`

GetMultiTenantOk returns a tuple with the MultiTenant field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMultiTenant

`func (o *FabricsPutRequestFabricValue) SetMultiTenant(v bool)`

SetMultiTenant sets MultiTenant field to given value.

### HasMultiTenant

`func (o *FabricsPutRequestFabricValue) HasMultiTenant() bool`

HasMultiTenant returns a boolean if a field has been set.

### GetBaseBgpAsNumber

`func (o *FabricsPutRequestFabricValue) GetBaseBgpAsNumber() string`

GetBaseBgpAsNumber returns the BaseBgpAsNumber field if non-nil, zero value otherwise.

### GetBaseBgpAsNumberOk

`func (o *FabricsPutRequestFabricValue) GetBaseBgpAsNumberOk() (*string, bool)`

GetBaseBgpAsNumberOk returns a tuple with the BaseBgpAsNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBaseBgpAsNumber

`func (o *FabricsPutRequestFabricValue) SetBaseBgpAsNumber(v string)`

SetBaseBgpAsNumber sets BaseBgpAsNumber field to given value.

### HasBaseBgpAsNumber

`func (o *FabricsPutRequestFabricValue) HasBaseBgpAsNumber() bool`

HasBaseBgpAsNumber returns a boolean if a field has been set.

### GetRouterIdBasePrefix

`func (o *FabricsPutRequestFabricValue) GetRouterIdBasePrefix() string`

GetRouterIdBasePrefix returns the RouterIdBasePrefix field if non-nil, zero value otherwise.

### GetRouterIdBasePrefixOk

`func (o *FabricsPutRequestFabricValue) GetRouterIdBasePrefixOk() (*string, bool)`

GetRouterIdBasePrefixOk returns a tuple with the RouterIdBasePrefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRouterIdBasePrefix

`func (o *FabricsPutRequestFabricValue) SetRouterIdBasePrefix(v string)`

SetRouterIdBasePrefix sets RouterIdBasePrefix field to given value.

### HasRouterIdBasePrefix

`func (o *FabricsPutRequestFabricValue) HasRouterIdBasePrefix() bool`

HasRouterIdBasePrefix returns a boolean if a field has been set.

### GetVtepIdBasePrefix

`func (o *FabricsPutRequestFabricValue) GetVtepIdBasePrefix() string`

GetVtepIdBasePrefix returns the VtepIdBasePrefix field if non-nil, zero value otherwise.

### GetVtepIdBasePrefixOk

`func (o *FabricsPutRequestFabricValue) GetVtepIdBasePrefixOk() (*string, bool)`

GetVtepIdBasePrefixOk returns a tuple with the VtepIdBasePrefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVtepIdBasePrefix

`func (o *FabricsPutRequestFabricValue) SetVtepIdBasePrefix(v string)`

SetVtepIdBasePrefix sets VtepIdBasePrefix field to given value.

### HasVtepIdBasePrefix

`func (o *FabricsPutRequestFabricValue) HasVtepIdBasePrefix() bool`

HasVtepIdBasePrefix returns a boolean if a field has been set.

### GetPairedIpSubnet

`func (o *FabricsPutRequestFabricValue) GetPairedIpSubnet() string`

GetPairedIpSubnet returns the PairedIpSubnet field if non-nil, zero value otherwise.

### GetPairedIpSubnetOk

`func (o *FabricsPutRequestFabricValue) GetPairedIpSubnetOk() (*string, bool)`

GetPairedIpSubnetOk returns a tuple with the PairedIpSubnet field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPairedIpSubnet

`func (o *FabricsPutRequestFabricValue) SetPairedIpSubnet(v string)`

SetPairedIpSubnet sets PairedIpSubnet field to given value.

### HasPairedIpSubnet

`func (o *FabricsPutRequestFabricValue) HasPairedIpSubnet() bool`

HasPairedIpSubnet returns a boolean if a field has been set.

### GetMaxSwitches

`func (o *FabricsPutRequestFabricValue) GetMaxSwitches() string`

GetMaxSwitches returns the MaxSwitches field if non-nil, zero value otherwise.

### GetMaxSwitchesOk

`func (o *FabricsPutRequestFabricValue) GetMaxSwitchesOk() (*string, bool)`

GetMaxSwitchesOk returns a tuple with the MaxSwitches field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxSwitches

`func (o *FabricsPutRequestFabricValue) SetMaxSwitches(v string)`

SetMaxSwitches sets MaxSwitches field to given value.

### HasMaxSwitches

`func (o *FabricsPutRequestFabricValue) HasMaxSwitches() bool`

HasMaxSwitches returns a boolean if a field has been set.

### GetPauseValidationAlarms

`func (o *FabricsPutRequestFabricValue) GetPauseValidationAlarms() bool`

GetPauseValidationAlarms returns the PauseValidationAlarms field if non-nil, zero value otherwise.

### GetPauseValidationAlarmsOk

`func (o *FabricsPutRequestFabricValue) GetPauseValidationAlarmsOk() (*bool, bool)`

GetPauseValidationAlarmsOk returns a tuple with the PauseValidationAlarms field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPauseValidationAlarms

`func (o *FabricsPutRequestFabricValue) SetPauseValidationAlarms(v bool)`

SetPauseValidationAlarms sets PauseValidationAlarms field to given value.

### HasPauseValidationAlarms

`func (o *FabricsPutRequestFabricValue) HasPauseValidationAlarms() bool`

HasPauseValidationAlarms returns a boolean if a field has been set.

### GetStartingOctet

`func (o *FabricsPutRequestFabricValue) GetStartingOctet() int64`

GetStartingOctet returns the StartingOctet field if non-nil, zero value otherwise.

### GetStartingOctetOk

`func (o *FabricsPutRequestFabricValue) GetStartingOctetOk() (*int64, bool)`

GetStartingOctetOk returns a tuple with the StartingOctet field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStartingOctet

`func (o *FabricsPutRequestFabricValue) SetStartingOctet(v int64)`

SetStartingOctet sets StartingOctet field to given value.

### HasStartingOctet

`func (o *FabricsPutRequestFabricValue) HasStartingOctet() bool`

HasStartingOctet returns a boolean if a field has been set.

### SetStartingOctetNil

`func (o *FabricsPutRequestFabricValue) SetStartingOctetNil(b bool)`

 SetStartingOctetNil sets the value for StartingOctet to be an explicit nil

### UnsetStartingOctet
`func (o *FabricsPutRequestFabricValue) UnsetStartingOctet()`

UnsetStartingOctet ensures that no value is present for StartingOctet, not even an explicit nil
### GetMaxSus

`func (o *FabricsPutRequestFabricValue) GetMaxSus() int64`

GetMaxSus returns the MaxSus field if non-nil, zero value otherwise.

### GetMaxSusOk

`func (o *FabricsPutRequestFabricValue) GetMaxSusOk() (*int64, bool)`

GetMaxSusOk returns a tuple with the MaxSus field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxSus

`func (o *FabricsPutRequestFabricValue) SetMaxSus(v int64)`

SetMaxSus sets MaxSus field to given value.

### HasMaxSus

`func (o *FabricsPutRequestFabricValue) HasMaxSus() bool`

HasMaxSus returns a boolean if a field has been set.

### SetMaxSusNil

`func (o *FabricsPutRequestFabricValue) SetMaxSusNil(b bool)`

 SetMaxSusNil sets the value for MaxSus to be an explicit nil

### UnsetMaxSus
`func (o *FabricsPutRequestFabricValue) UnsetMaxSus()`

UnsetMaxSus ensures that no value is present for MaxSus, not even an explicit nil
### GetMaxPods

`func (o *FabricsPutRequestFabricValue) GetMaxPods() int64`

GetMaxPods returns the MaxPods field if non-nil, zero value otherwise.

### GetMaxPodsOk

`func (o *FabricsPutRequestFabricValue) GetMaxPodsOk() (*int64, bool)`

GetMaxPodsOk returns a tuple with the MaxPods field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxPods

`func (o *FabricsPutRequestFabricValue) SetMaxPods(v int64)`

SetMaxPods sets MaxPods field to given value.

### HasMaxPods

`func (o *FabricsPutRequestFabricValue) HasMaxPods() bool`

HasMaxPods returns a boolean if a field has been set.

### SetMaxPodsNil

`func (o *FabricsPutRequestFabricValue) SetMaxPodsNil(b bool)`

 SetMaxPodsNil sets the value for MaxPods to be an explicit nil

### UnsetMaxPods
`func (o *FabricsPutRequestFabricValue) UnsetMaxPods()`

UnsetMaxPods ensures that no value is present for MaxPods, not even an explicit nil
### GetObjectProperties

`func (o *FabricsPutRequestFabricValue) GetObjectProperties() FabricsPutRequestFabricValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *FabricsPutRequestFabricValue) GetObjectPropertiesOk() (*FabricsPutRequestFabricValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *FabricsPutRequestFabricValue) SetObjectProperties(v FabricsPutRequestFabricValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *FabricsPutRequestFabricValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.

### GetSetLeafRouterIdOnBgp

`func (o *FabricsPutRequestFabricValue) GetSetLeafRouterIdOnBgp() bool`

GetSetLeafRouterIdOnBgp returns the SetLeafRouterIdOnBgp field if non-nil, zero value otherwise.

### GetSetLeafRouterIdOnBgpOk

`func (o *FabricsPutRequestFabricValue) GetSetLeafRouterIdOnBgpOk() (*bool, bool)`

GetSetLeafRouterIdOnBgpOk returns a tuple with the SetLeafRouterIdOnBgp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSetLeafRouterIdOnBgp

`func (o *FabricsPutRequestFabricValue) SetSetLeafRouterIdOnBgp(v bool)`

SetSetLeafRouterIdOnBgp sets SetLeafRouterIdOnBgp field to given value.

### HasSetLeafRouterIdOnBgp

`func (o *FabricsPutRequestFabricValue) HasSetLeafRouterIdOnBgp() bool`

HasSetLeafRouterIdOnBgp returns a boolean if a field has been set.

### GetRouteAggregation

`func (o *FabricsPutRequestFabricValue) GetRouteAggregation() string`

GetRouteAggregation returns the RouteAggregation field if non-nil, zero value otherwise.

### GetRouteAggregationOk

`func (o *FabricsPutRequestFabricValue) GetRouteAggregationOk() (*string, bool)`

GetRouteAggregationOk returns a tuple with the RouteAggregation field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRouteAggregation

`func (o *FabricsPutRequestFabricValue) SetRouteAggregation(v string)`

SetRouteAggregation sets RouteAggregation field to given value.

### HasRouteAggregation

`func (o *FabricsPutRequestFabricValue) HasRouteAggregation() bool`

HasRouteAggregation returns a boolean if a field has been set.

### GetRouteAggregators

`func (o *FabricsPutRequestFabricValue) GetRouteAggregators() []FabricsPutRequestFabricValueRouteAggregatorsInner`

GetRouteAggregators returns the RouteAggregators field if non-nil, zero value otherwise.

### GetRouteAggregatorsOk

`func (o *FabricsPutRequestFabricValue) GetRouteAggregatorsOk() (*[]FabricsPutRequestFabricValueRouteAggregatorsInner, bool)`

GetRouteAggregatorsOk returns a tuple with the RouteAggregators field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRouteAggregators

`func (o *FabricsPutRequestFabricValue) SetRouteAggregators(v []FabricsPutRequestFabricValueRouteAggregatorsInner)`

SetRouteAggregators sets RouteAggregators field to given value.

### HasRouteAggregators

`func (o *FabricsPutRequestFabricValue) HasRouteAggregators() bool`

HasRouteAggregators returns a boolean if a field has been set.

### GetIpSourceGuard

`func (o *FabricsPutRequestFabricValue) GetIpSourceGuard() bool`

GetIpSourceGuard returns the IpSourceGuard field if non-nil, zero value otherwise.

### GetIpSourceGuardOk

`func (o *FabricsPutRequestFabricValue) GetIpSourceGuardOk() (*bool, bool)`

GetIpSourceGuardOk returns a tuple with the IpSourceGuard field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpSourceGuard

`func (o *FabricsPutRequestFabricValue) SetIpSourceGuard(v bool)`

SetIpSourceGuard sets IpSourceGuard field to given value.

### HasIpSourceGuard

`func (o *FabricsPutRequestFabricValue) HasIpSourceGuard() bool`

HasIpSourceGuard returns a boolean if a field has been set.

### GetEnableDhcpSnooping

`func (o *FabricsPutRequestFabricValue) GetEnableDhcpSnooping() bool`

GetEnableDhcpSnooping returns the EnableDhcpSnooping field if non-nil, zero value otherwise.

### GetEnableDhcpSnoopingOk

`func (o *FabricsPutRequestFabricValue) GetEnableDhcpSnoopingOk() (*bool, bool)`

GetEnableDhcpSnoopingOk returns a tuple with the EnableDhcpSnooping field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnableDhcpSnooping

`func (o *FabricsPutRequestFabricValue) SetEnableDhcpSnooping(v bool)`

SetEnableDhcpSnooping sets EnableDhcpSnooping field to given value.

### HasEnableDhcpSnooping

`func (o *FabricsPutRequestFabricValue) HasEnableDhcpSnooping() bool`

HasEnableDhcpSnooping returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


