# ConfigPutRequestServiceServiceName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. It&#39;s highly recommended to set this value to true so that validation on the object will be ran. | [optional] [default to false]
**Vlan** | Pointer to **NullableInt32** | A Value between 1 and 4096 | [optional] 
**Vni** | Pointer to **NullableInt32** | Indication of the outgoing VLAN layer 2 service | [optional] 
**VniAutoAssigned** | Pointer to **bool** | Whether or not the value in vni field has been automatically assigned or not. Set to false and change vni value to edit. | [optional] 
**Tenant** | Pointer to **string** | Tenant | [optional] [default to ""]
**TenantRefType** | Pointer to **string** | Object type for tenant field | [optional] 
**AnycastIpMask** | Pointer to **string** | Static anycast gateway address for service  | [optional] [default to ""]
**DhcpServerIp** | Pointer to **string** | IP address(s) of the DHCP server for service.  May have up to four separated by commas. | [optional] [default to ""]
**Mtu** | Pointer to **NullableInt32** | MTU (Maximum Transmission Unit) The size used by a switch to determine when large packets must be broken up into smaller packets for delivery. If mismatched within a single vlan network, can cause dropped packets. | [optional] [default to 1500]
**ObjectProperties** | Pointer to [**ConfigPutRequestServiceServiceNameObjectProperties**](ConfigPutRequestServiceServiceNameObjectProperties.md) |  | [optional] 
**AnycastIpv4Mask** | Pointer to **string** | Static anycast gateway addresses(IPv4) for service  | [optional] [default to ""]
**AnycastIpv6Mask** | Pointer to **string** | Static anycast gateway addresses(IPv6) for service  | [optional] [default to ""]
**MaxUpstreamRateMbps** | Pointer to **int32** | Bandwidth allocated per port in the upstream direction. (Max 10000 Mbps) | [optional] 
**MaxDownstreamRateMbps** | Pointer to **int32** | Bandwidth allocated per port in the downstream direction. (Max 10000 Mbps) | [optional] 
**PacketPriority** | Pointer to **string** | Priority untagged packets will be tagged with on ingress to the network. If the network is flooded packets of lower priority will be dropped | [optional] [default to "0"]
**MulticastManagementMode** | Pointer to **string** | Determines how undefined handle multicast packet for Service&lt;ul&gt;&lt;li&gt;* \&quot;Multicast Flooding (Normal)\&quot; Multicast packets are broadcast&lt;/li&gt;&lt;li&gt;* \&quot;Multicast Flooding (AVB/PTP/Cobranet)\&quot; Multicast packets are broadcast with special treatment for critical latency packets such as used by AVB, PTP, and Cobranet&lt;/li&gt;&lt;li&gt;* \&quot;IPTV Filtering (IGMP Snooping)\&quot; Multicast packets are propagated via IGMP Snooping&lt;/li&gt;&lt;li&gt;* \&quot;IPTV Filtering (IGMP Report/Leave Flooding)\&quot; Multicast packets are propagated via IGMP Snooping. except that IGMP Report/Leave packets are broadcast&lt;/li&gt;&lt;/ul&gt; | [optional] [default to "flooding"]
**TaggedPackets** | Pointer to **bool** | Overrides priority bits on incoming tagged packets. Always done for untagged packets | [optional] [default to false]
**Tls** | Pointer to **bool** | Is a Transparent LAN Service? | [optional] [default to false]
**AllowLocalSwitching** | Pointer to **bool** | Allow Edge Devices to communicate with each other. Disabling this forces upstream traffic to the router | [optional] [default to true]
**ActAsMulticastQuerier** | Pointer to **bool** | Multicast managment through IGMP requires a multicast querier. Check this box if SD LAN should provide a multicast querier | [optional] [default to false]
**BlockUnknownUnicastFlood** | Pointer to **bool** | Block unknown unicast traffic flooding and only permits egress traffic with MAC addresses that are known to exit on the port | [optional] [default to false]
**BlockDownstreamDhcpServer** | Pointer to **bool** | Block inbound packets sent by Downstream DHCP servers | [optional] [default to true]
**IsManagementService** | Pointer to **bool** | Denotes a Management Service | [optional] [default to false]
**UseDscpToPBitMappingForL3PacketsIfAvailable** | Pointer to **bool** | use DSCP to p-bit Mapping for L3 packets if available | [optional] [default to false]
**AllowFastLeave** | Pointer to **bool** | The Fast Leave feature causes the switch to immediately remove a port from the forwarding list for a IGMP multicast group when the port receives a leave message. Not recommended unless there is only a single receiver present on every point in the VLAN | [optional] [default to false]
**MstInstance** | Pointer to **int32** | MST Instance ID (0-4094) | [optional] [default to 0]

## Methods

### NewConfigPutRequestServiceServiceName

`func NewConfigPutRequestServiceServiceName() *ConfigPutRequestServiceServiceName`

NewConfigPutRequestServiceServiceName instantiates a new ConfigPutRequestServiceServiceName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestServiceServiceNameWithDefaults

`func NewConfigPutRequestServiceServiceNameWithDefaults() *ConfigPutRequestServiceServiceName`

NewConfigPutRequestServiceServiceNameWithDefaults instantiates a new ConfigPutRequestServiceServiceName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestServiceServiceName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestServiceServiceName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestServiceServiceName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestServiceServiceName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ConfigPutRequestServiceServiceName) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestServiceServiceName) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestServiceServiceName) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestServiceServiceName) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetVlan

`func (o *ConfigPutRequestServiceServiceName) GetVlan() int32`

GetVlan returns the Vlan field if non-nil, zero value otherwise.

### GetVlanOk

`func (o *ConfigPutRequestServiceServiceName) GetVlanOk() (*int32, bool)`

GetVlanOk returns a tuple with the Vlan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVlan

`func (o *ConfigPutRequestServiceServiceName) SetVlan(v int32)`

SetVlan sets Vlan field to given value.

### HasVlan

`func (o *ConfigPutRequestServiceServiceName) HasVlan() bool`

HasVlan returns a boolean if a field has been set.

### SetVlanNil

`func (o *ConfigPutRequestServiceServiceName) SetVlanNil(b bool)`

 SetVlanNil sets the value for Vlan to be an explicit nil

### UnsetVlan
`func (o *ConfigPutRequestServiceServiceName) UnsetVlan()`

UnsetVlan ensures that no value is present for Vlan, not even an explicit nil
### GetVni

`func (o *ConfigPutRequestServiceServiceName) GetVni() int32`

GetVni returns the Vni field if non-nil, zero value otherwise.

### GetVniOk

`func (o *ConfigPutRequestServiceServiceName) GetVniOk() (*int32, bool)`

GetVniOk returns a tuple with the Vni field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVni

`func (o *ConfigPutRequestServiceServiceName) SetVni(v int32)`

SetVni sets Vni field to given value.

### HasVni

`func (o *ConfigPutRequestServiceServiceName) HasVni() bool`

HasVni returns a boolean if a field has been set.

### SetVniNil

`func (o *ConfigPutRequestServiceServiceName) SetVniNil(b bool)`

 SetVniNil sets the value for Vni to be an explicit nil

### UnsetVni
`func (o *ConfigPutRequestServiceServiceName) UnsetVni()`

UnsetVni ensures that no value is present for Vni, not even an explicit nil
### GetVniAutoAssigned

`func (o *ConfigPutRequestServiceServiceName) GetVniAutoAssigned() bool`

GetVniAutoAssigned returns the VniAutoAssigned field if non-nil, zero value otherwise.

### GetVniAutoAssignedOk

`func (o *ConfigPutRequestServiceServiceName) GetVniAutoAssignedOk() (*bool, bool)`

GetVniAutoAssignedOk returns a tuple with the VniAutoAssigned field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVniAutoAssigned

`func (o *ConfigPutRequestServiceServiceName) SetVniAutoAssigned(v bool)`

SetVniAutoAssigned sets VniAutoAssigned field to given value.

### HasVniAutoAssigned

`func (o *ConfigPutRequestServiceServiceName) HasVniAutoAssigned() bool`

HasVniAutoAssigned returns a boolean if a field has been set.

### GetTenant

`func (o *ConfigPutRequestServiceServiceName) GetTenant() string`

GetTenant returns the Tenant field if non-nil, zero value otherwise.

### GetTenantOk

`func (o *ConfigPutRequestServiceServiceName) GetTenantOk() (*string, bool)`

GetTenantOk returns a tuple with the Tenant field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenant

`func (o *ConfigPutRequestServiceServiceName) SetTenant(v string)`

SetTenant sets Tenant field to given value.

### HasTenant

`func (o *ConfigPutRequestServiceServiceName) HasTenant() bool`

HasTenant returns a boolean if a field has been set.

### GetTenantRefType

`func (o *ConfigPutRequestServiceServiceName) GetTenantRefType() string`

GetTenantRefType returns the TenantRefType field if non-nil, zero value otherwise.

### GetTenantRefTypeOk

`func (o *ConfigPutRequestServiceServiceName) GetTenantRefTypeOk() (*string, bool)`

GetTenantRefTypeOk returns a tuple with the TenantRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantRefType

`func (o *ConfigPutRequestServiceServiceName) SetTenantRefType(v string)`

SetTenantRefType sets TenantRefType field to given value.

### HasTenantRefType

`func (o *ConfigPutRequestServiceServiceName) HasTenantRefType() bool`

HasTenantRefType returns a boolean if a field has been set.

### GetAnycastIpMask

`func (o *ConfigPutRequestServiceServiceName) GetAnycastIpMask() string`

GetAnycastIpMask returns the AnycastIpMask field if non-nil, zero value otherwise.

### GetAnycastIpMaskOk

`func (o *ConfigPutRequestServiceServiceName) GetAnycastIpMaskOk() (*string, bool)`

GetAnycastIpMaskOk returns a tuple with the AnycastIpMask field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnycastIpMask

`func (o *ConfigPutRequestServiceServiceName) SetAnycastIpMask(v string)`

SetAnycastIpMask sets AnycastIpMask field to given value.

### HasAnycastIpMask

`func (o *ConfigPutRequestServiceServiceName) HasAnycastIpMask() bool`

HasAnycastIpMask returns a boolean if a field has been set.

### GetDhcpServerIp

`func (o *ConfigPutRequestServiceServiceName) GetDhcpServerIp() string`

GetDhcpServerIp returns the DhcpServerIp field if non-nil, zero value otherwise.

### GetDhcpServerIpOk

`func (o *ConfigPutRequestServiceServiceName) GetDhcpServerIpOk() (*string, bool)`

GetDhcpServerIpOk returns a tuple with the DhcpServerIp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDhcpServerIp

`func (o *ConfigPutRequestServiceServiceName) SetDhcpServerIp(v string)`

SetDhcpServerIp sets DhcpServerIp field to given value.

### HasDhcpServerIp

`func (o *ConfigPutRequestServiceServiceName) HasDhcpServerIp() bool`

HasDhcpServerIp returns a boolean if a field has been set.

### GetMtu

`func (o *ConfigPutRequestServiceServiceName) GetMtu() int32`

GetMtu returns the Mtu field if non-nil, zero value otherwise.

### GetMtuOk

`func (o *ConfigPutRequestServiceServiceName) GetMtuOk() (*int32, bool)`

GetMtuOk returns a tuple with the Mtu field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMtu

`func (o *ConfigPutRequestServiceServiceName) SetMtu(v int32)`

SetMtu sets Mtu field to given value.

### HasMtu

`func (o *ConfigPutRequestServiceServiceName) HasMtu() bool`

HasMtu returns a boolean if a field has been set.

### SetMtuNil

`func (o *ConfigPutRequestServiceServiceName) SetMtuNil(b bool)`

 SetMtuNil sets the value for Mtu to be an explicit nil

### UnsetMtu
`func (o *ConfigPutRequestServiceServiceName) UnsetMtu()`

UnsetMtu ensures that no value is present for Mtu, not even an explicit nil
### GetObjectProperties

`func (o *ConfigPutRequestServiceServiceName) GetObjectProperties() ConfigPutRequestServiceServiceNameObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ConfigPutRequestServiceServiceName) GetObjectPropertiesOk() (*ConfigPutRequestServiceServiceNameObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ConfigPutRequestServiceServiceName) SetObjectProperties(v ConfigPutRequestServiceServiceNameObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ConfigPutRequestServiceServiceName) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.

### GetAnycastIpv4Mask

`func (o *ConfigPutRequestServiceServiceName) GetAnycastIpv4Mask() string`

GetAnycastIpv4Mask returns the AnycastIpv4Mask field if non-nil, zero value otherwise.

### GetAnycastIpv4MaskOk

`func (o *ConfigPutRequestServiceServiceName) GetAnycastIpv4MaskOk() (*string, bool)`

GetAnycastIpv4MaskOk returns a tuple with the AnycastIpv4Mask field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnycastIpv4Mask

`func (o *ConfigPutRequestServiceServiceName) SetAnycastIpv4Mask(v string)`

SetAnycastIpv4Mask sets AnycastIpv4Mask field to given value.

### HasAnycastIpv4Mask

`func (o *ConfigPutRequestServiceServiceName) HasAnycastIpv4Mask() bool`

HasAnycastIpv4Mask returns a boolean if a field has been set.

### GetAnycastIpv6Mask

`func (o *ConfigPutRequestServiceServiceName) GetAnycastIpv6Mask() string`

GetAnycastIpv6Mask returns the AnycastIpv6Mask field if non-nil, zero value otherwise.

### GetAnycastIpv6MaskOk

`func (o *ConfigPutRequestServiceServiceName) GetAnycastIpv6MaskOk() (*string, bool)`

GetAnycastIpv6MaskOk returns a tuple with the AnycastIpv6Mask field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnycastIpv6Mask

`func (o *ConfigPutRequestServiceServiceName) SetAnycastIpv6Mask(v string)`

SetAnycastIpv6Mask sets AnycastIpv6Mask field to given value.

### HasAnycastIpv6Mask

`func (o *ConfigPutRequestServiceServiceName) HasAnycastIpv6Mask() bool`

HasAnycastIpv6Mask returns a boolean if a field has been set.

### GetMaxUpstreamRateMbps

`func (o *ConfigPutRequestServiceServiceName) GetMaxUpstreamRateMbps() int32`

GetMaxUpstreamRateMbps returns the MaxUpstreamRateMbps field if non-nil, zero value otherwise.

### GetMaxUpstreamRateMbpsOk

`func (o *ConfigPutRequestServiceServiceName) GetMaxUpstreamRateMbpsOk() (*int32, bool)`

GetMaxUpstreamRateMbpsOk returns a tuple with the MaxUpstreamRateMbps field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxUpstreamRateMbps

`func (o *ConfigPutRequestServiceServiceName) SetMaxUpstreamRateMbps(v int32)`

SetMaxUpstreamRateMbps sets MaxUpstreamRateMbps field to given value.

### HasMaxUpstreamRateMbps

`func (o *ConfigPutRequestServiceServiceName) HasMaxUpstreamRateMbps() bool`

HasMaxUpstreamRateMbps returns a boolean if a field has been set.

### GetMaxDownstreamRateMbps

`func (o *ConfigPutRequestServiceServiceName) GetMaxDownstreamRateMbps() int32`

GetMaxDownstreamRateMbps returns the MaxDownstreamRateMbps field if non-nil, zero value otherwise.

### GetMaxDownstreamRateMbpsOk

`func (o *ConfigPutRequestServiceServiceName) GetMaxDownstreamRateMbpsOk() (*int32, bool)`

GetMaxDownstreamRateMbpsOk returns a tuple with the MaxDownstreamRateMbps field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxDownstreamRateMbps

`func (o *ConfigPutRequestServiceServiceName) SetMaxDownstreamRateMbps(v int32)`

SetMaxDownstreamRateMbps sets MaxDownstreamRateMbps field to given value.

### HasMaxDownstreamRateMbps

`func (o *ConfigPutRequestServiceServiceName) HasMaxDownstreamRateMbps() bool`

HasMaxDownstreamRateMbps returns a boolean if a field has been set.

### GetPacketPriority

`func (o *ConfigPutRequestServiceServiceName) GetPacketPriority() string`

GetPacketPriority returns the PacketPriority field if non-nil, zero value otherwise.

### GetPacketPriorityOk

`func (o *ConfigPutRequestServiceServiceName) GetPacketPriorityOk() (*string, bool)`

GetPacketPriorityOk returns a tuple with the PacketPriority field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPacketPriority

`func (o *ConfigPutRequestServiceServiceName) SetPacketPriority(v string)`

SetPacketPriority sets PacketPriority field to given value.

### HasPacketPriority

`func (o *ConfigPutRequestServiceServiceName) HasPacketPriority() bool`

HasPacketPriority returns a boolean if a field has been set.

### GetMulticastManagementMode

`func (o *ConfigPutRequestServiceServiceName) GetMulticastManagementMode() string`

GetMulticastManagementMode returns the MulticastManagementMode field if non-nil, zero value otherwise.

### GetMulticastManagementModeOk

`func (o *ConfigPutRequestServiceServiceName) GetMulticastManagementModeOk() (*string, bool)`

GetMulticastManagementModeOk returns a tuple with the MulticastManagementMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMulticastManagementMode

`func (o *ConfigPutRequestServiceServiceName) SetMulticastManagementMode(v string)`

SetMulticastManagementMode sets MulticastManagementMode field to given value.

### HasMulticastManagementMode

`func (o *ConfigPutRequestServiceServiceName) HasMulticastManagementMode() bool`

HasMulticastManagementMode returns a boolean if a field has been set.

### GetTaggedPackets

`func (o *ConfigPutRequestServiceServiceName) GetTaggedPackets() bool`

GetTaggedPackets returns the TaggedPackets field if non-nil, zero value otherwise.

### GetTaggedPacketsOk

`func (o *ConfigPutRequestServiceServiceName) GetTaggedPacketsOk() (*bool, bool)`

GetTaggedPacketsOk returns a tuple with the TaggedPackets field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTaggedPackets

`func (o *ConfigPutRequestServiceServiceName) SetTaggedPackets(v bool)`

SetTaggedPackets sets TaggedPackets field to given value.

### HasTaggedPackets

`func (o *ConfigPutRequestServiceServiceName) HasTaggedPackets() bool`

HasTaggedPackets returns a boolean if a field has been set.

### GetTls

`func (o *ConfigPutRequestServiceServiceName) GetTls() bool`

GetTls returns the Tls field if non-nil, zero value otherwise.

### GetTlsOk

`func (o *ConfigPutRequestServiceServiceName) GetTlsOk() (*bool, bool)`

GetTlsOk returns a tuple with the Tls field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTls

`func (o *ConfigPutRequestServiceServiceName) SetTls(v bool)`

SetTls sets Tls field to given value.

### HasTls

`func (o *ConfigPutRequestServiceServiceName) HasTls() bool`

HasTls returns a boolean if a field has been set.

### GetAllowLocalSwitching

`func (o *ConfigPutRequestServiceServiceName) GetAllowLocalSwitching() bool`

GetAllowLocalSwitching returns the AllowLocalSwitching field if non-nil, zero value otherwise.

### GetAllowLocalSwitchingOk

`func (o *ConfigPutRequestServiceServiceName) GetAllowLocalSwitchingOk() (*bool, bool)`

GetAllowLocalSwitchingOk returns a tuple with the AllowLocalSwitching field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllowLocalSwitching

`func (o *ConfigPutRequestServiceServiceName) SetAllowLocalSwitching(v bool)`

SetAllowLocalSwitching sets AllowLocalSwitching field to given value.

### HasAllowLocalSwitching

`func (o *ConfigPutRequestServiceServiceName) HasAllowLocalSwitching() bool`

HasAllowLocalSwitching returns a boolean if a field has been set.

### GetActAsMulticastQuerier

`func (o *ConfigPutRequestServiceServiceName) GetActAsMulticastQuerier() bool`

GetActAsMulticastQuerier returns the ActAsMulticastQuerier field if non-nil, zero value otherwise.

### GetActAsMulticastQuerierOk

`func (o *ConfigPutRequestServiceServiceName) GetActAsMulticastQuerierOk() (*bool, bool)`

GetActAsMulticastQuerierOk returns a tuple with the ActAsMulticastQuerier field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActAsMulticastQuerier

`func (o *ConfigPutRequestServiceServiceName) SetActAsMulticastQuerier(v bool)`

SetActAsMulticastQuerier sets ActAsMulticastQuerier field to given value.

### HasActAsMulticastQuerier

`func (o *ConfigPutRequestServiceServiceName) HasActAsMulticastQuerier() bool`

HasActAsMulticastQuerier returns a boolean if a field has been set.

### GetBlockUnknownUnicastFlood

`func (o *ConfigPutRequestServiceServiceName) GetBlockUnknownUnicastFlood() bool`

GetBlockUnknownUnicastFlood returns the BlockUnknownUnicastFlood field if non-nil, zero value otherwise.

### GetBlockUnknownUnicastFloodOk

`func (o *ConfigPutRequestServiceServiceName) GetBlockUnknownUnicastFloodOk() (*bool, bool)`

GetBlockUnknownUnicastFloodOk returns a tuple with the BlockUnknownUnicastFlood field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBlockUnknownUnicastFlood

`func (o *ConfigPutRequestServiceServiceName) SetBlockUnknownUnicastFlood(v bool)`

SetBlockUnknownUnicastFlood sets BlockUnknownUnicastFlood field to given value.

### HasBlockUnknownUnicastFlood

`func (o *ConfigPutRequestServiceServiceName) HasBlockUnknownUnicastFlood() bool`

HasBlockUnknownUnicastFlood returns a boolean if a field has been set.

### GetBlockDownstreamDhcpServer

`func (o *ConfigPutRequestServiceServiceName) GetBlockDownstreamDhcpServer() bool`

GetBlockDownstreamDhcpServer returns the BlockDownstreamDhcpServer field if non-nil, zero value otherwise.

### GetBlockDownstreamDhcpServerOk

`func (o *ConfigPutRequestServiceServiceName) GetBlockDownstreamDhcpServerOk() (*bool, bool)`

GetBlockDownstreamDhcpServerOk returns a tuple with the BlockDownstreamDhcpServer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBlockDownstreamDhcpServer

`func (o *ConfigPutRequestServiceServiceName) SetBlockDownstreamDhcpServer(v bool)`

SetBlockDownstreamDhcpServer sets BlockDownstreamDhcpServer field to given value.

### HasBlockDownstreamDhcpServer

`func (o *ConfigPutRequestServiceServiceName) HasBlockDownstreamDhcpServer() bool`

HasBlockDownstreamDhcpServer returns a boolean if a field has been set.

### GetIsManagementService

`func (o *ConfigPutRequestServiceServiceName) GetIsManagementService() bool`

GetIsManagementService returns the IsManagementService field if non-nil, zero value otherwise.

### GetIsManagementServiceOk

`func (o *ConfigPutRequestServiceServiceName) GetIsManagementServiceOk() (*bool, bool)`

GetIsManagementServiceOk returns a tuple with the IsManagementService field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsManagementService

`func (o *ConfigPutRequestServiceServiceName) SetIsManagementService(v bool)`

SetIsManagementService sets IsManagementService field to given value.

### HasIsManagementService

`func (o *ConfigPutRequestServiceServiceName) HasIsManagementService() bool`

HasIsManagementService returns a boolean if a field has been set.

### GetUseDscpToPBitMappingForL3PacketsIfAvailable

`func (o *ConfigPutRequestServiceServiceName) GetUseDscpToPBitMappingForL3PacketsIfAvailable() bool`

GetUseDscpToPBitMappingForL3PacketsIfAvailable returns the UseDscpToPBitMappingForL3PacketsIfAvailable field if non-nil, zero value otherwise.

### GetUseDscpToPBitMappingForL3PacketsIfAvailableOk

`func (o *ConfigPutRequestServiceServiceName) GetUseDscpToPBitMappingForL3PacketsIfAvailableOk() (*bool, bool)`

GetUseDscpToPBitMappingForL3PacketsIfAvailableOk returns a tuple with the UseDscpToPBitMappingForL3PacketsIfAvailable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUseDscpToPBitMappingForL3PacketsIfAvailable

`func (o *ConfigPutRequestServiceServiceName) SetUseDscpToPBitMappingForL3PacketsIfAvailable(v bool)`

SetUseDscpToPBitMappingForL3PacketsIfAvailable sets UseDscpToPBitMappingForL3PacketsIfAvailable field to given value.

### HasUseDscpToPBitMappingForL3PacketsIfAvailable

`func (o *ConfigPutRequestServiceServiceName) HasUseDscpToPBitMappingForL3PacketsIfAvailable() bool`

HasUseDscpToPBitMappingForL3PacketsIfAvailable returns a boolean if a field has been set.

### GetAllowFastLeave

`func (o *ConfigPutRequestServiceServiceName) GetAllowFastLeave() bool`

GetAllowFastLeave returns the AllowFastLeave field if non-nil, zero value otherwise.

### GetAllowFastLeaveOk

`func (o *ConfigPutRequestServiceServiceName) GetAllowFastLeaveOk() (*bool, bool)`

GetAllowFastLeaveOk returns a tuple with the AllowFastLeave field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllowFastLeave

`func (o *ConfigPutRequestServiceServiceName) SetAllowFastLeave(v bool)`

SetAllowFastLeave sets AllowFastLeave field to given value.

### HasAllowFastLeave

`func (o *ConfigPutRequestServiceServiceName) HasAllowFastLeave() bool`

HasAllowFastLeave returns a boolean if a field has been set.

### GetMstInstance

`func (o *ConfigPutRequestServiceServiceName) GetMstInstance() int32`

GetMstInstance returns the MstInstance field if non-nil, zero value otherwise.

### GetMstInstanceOk

`func (o *ConfigPutRequestServiceServiceName) GetMstInstanceOk() (*int32, bool)`

GetMstInstanceOk returns a tuple with the MstInstance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMstInstance

`func (o *ConfigPutRequestServiceServiceName) SetMstInstance(v int32)`

SetMstInstance sets MstInstance field to given value.

### HasMstInstance

`func (o *ConfigPutRequestServiceServiceName) HasMstInstance() bool`

HasMstInstance returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


