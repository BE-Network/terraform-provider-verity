# ServicesPutRequestServiceValue

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
**AnycastIpv4Mask** | Pointer to **string** | Comma separated list of Static anycast gateway addresses(IPv4) for service  | [optional] [default to ""]
**AnycastIpv6Mask** | Pointer to **string** | Comma separated list of Static anycast gateway addresses(IPv6) for service  | [optional] [default to ""]
**DhcpServerIpv4** | Pointer to **string** | IPv4 address(s) of the DHCP server for service.  May have up to four separated by commas. | [optional] [default to ""]
**DhcpServerIpv6** | Pointer to **string** | IPv6 address(s) of the DHCP server for service.  May have up to four separated by commas. | [optional] [default to ""]
**Mtu** | Pointer to **NullableInt32** | MTU (Maximum Transmission Unit) The size used by a switch to determine when large packets must be broken up into smaller packets for delivery. If mismatched within a single vlan network, can cause dropped packets. | [optional] [default to 1500]
**ObjectProperties** | Pointer to [**ServicesPutRequestServiceValueObjectProperties**](ServicesPutRequestServiceValueObjectProperties.md) |  | [optional] 
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
**PolicyBasedRouting** | Pointer to **string** | Policy Based Routing | [optional] [default to ""]
**PolicyBasedRoutingRefType** | Pointer to **string** | Object type for policy_based_routing field | [optional] 

## Methods

### NewServicesPutRequestServiceValue

`func NewServicesPutRequestServiceValue() *ServicesPutRequestServiceValue`

NewServicesPutRequestServiceValue instantiates a new ServicesPutRequestServiceValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewServicesPutRequestServiceValueWithDefaults

`func NewServicesPutRequestServiceValueWithDefaults() *ServicesPutRequestServiceValue`

NewServicesPutRequestServiceValueWithDefaults instantiates a new ServicesPutRequestServiceValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ServicesPutRequestServiceValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ServicesPutRequestServiceValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ServicesPutRequestServiceValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ServicesPutRequestServiceValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ServicesPutRequestServiceValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ServicesPutRequestServiceValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ServicesPutRequestServiceValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ServicesPutRequestServiceValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetVlan

`func (o *ServicesPutRequestServiceValue) GetVlan() int32`

GetVlan returns the Vlan field if non-nil, zero value otherwise.

### GetVlanOk

`func (o *ServicesPutRequestServiceValue) GetVlanOk() (*int32, bool)`

GetVlanOk returns a tuple with the Vlan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVlan

`func (o *ServicesPutRequestServiceValue) SetVlan(v int32)`

SetVlan sets Vlan field to given value.

### HasVlan

`func (o *ServicesPutRequestServiceValue) HasVlan() bool`

HasVlan returns a boolean if a field has been set.

### SetVlanNil

`func (o *ServicesPutRequestServiceValue) SetVlanNil(b bool)`

 SetVlanNil sets the value for Vlan to be an explicit nil

### UnsetVlan
`func (o *ServicesPutRequestServiceValue) UnsetVlan()`

UnsetVlan ensures that no value is present for Vlan, not even an explicit nil
### GetVni

`func (o *ServicesPutRequestServiceValue) GetVni() int32`

GetVni returns the Vni field if non-nil, zero value otherwise.

### GetVniOk

`func (o *ServicesPutRequestServiceValue) GetVniOk() (*int32, bool)`

GetVniOk returns a tuple with the Vni field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVni

`func (o *ServicesPutRequestServiceValue) SetVni(v int32)`

SetVni sets Vni field to given value.

### HasVni

`func (o *ServicesPutRequestServiceValue) HasVni() bool`

HasVni returns a boolean if a field has been set.

### SetVniNil

`func (o *ServicesPutRequestServiceValue) SetVniNil(b bool)`

 SetVniNil sets the value for Vni to be an explicit nil

### UnsetVni
`func (o *ServicesPutRequestServiceValue) UnsetVni()`

UnsetVni ensures that no value is present for Vni, not even an explicit nil
### GetVniAutoAssigned

`func (o *ServicesPutRequestServiceValue) GetVniAutoAssigned() bool`

GetVniAutoAssigned returns the VniAutoAssigned field if non-nil, zero value otherwise.

### GetVniAutoAssignedOk

`func (o *ServicesPutRequestServiceValue) GetVniAutoAssignedOk() (*bool, bool)`

GetVniAutoAssignedOk returns a tuple with the VniAutoAssigned field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVniAutoAssigned

`func (o *ServicesPutRequestServiceValue) SetVniAutoAssigned(v bool)`

SetVniAutoAssigned sets VniAutoAssigned field to given value.

### HasVniAutoAssigned

`func (o *ServicesPutRequestServiceValue) HasVniAutoAssigned() bool`

HasVniAutoAssigned returns a boolean if a field has been set.

### GetTenant

`func (o *ServicesPutRequestServiceValue) GetTenant() string`

GetTenant returns the Tenant field if non-nil, zero value otherwise.

### GetTenantOk

`func (o *ServicesPutRequestServiceValue) GetTenantOk() (*string, bool)`

GetTenantOk returns a tuple with the Tenant field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenant

`func (o *ServicesPutRequestServiceValue) SetTenant(v string)`

SetTenant sets Tenant field to given value.

### HasTenant

`func (o *ServicesPutRequestServiceValue) HasTenant() bool`

HasTenant returns a boolean if a field has been set.

### GetTenantRefType

`func (o *ServicesPutRequestServiceValue) GetTenantRefType() string`

GetTenantRefType returns the TenantRefType field if non-nil, zero value otherwise.

### GetTenantRefTypeOk

`func (o *ServicesPutRequestServiceValue) GetTenantRefTypeOk() (*string, bool)`

GetTenantRefTypeOk returns a tuple with the TenantRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantRefType

`func (o *ServicesPutRequestServiceValue) SetTenantRefType(v string)`

SetTenantRefType sets TenantRefType field to given value.

### HasTenantRefType

`func (o *ServicesPutRequestServiceValue) HasTenantRefType() bool`

HasTenantRefType returns a boolean if a field has been set.

### GetAnycastIpv4Mask

`func (o *ServicesPutRequestServiceValue) GetAnycastIpv4Mask() string`

GetAnycastIpv4Mask returns the AnycastIpv4Mask field if non-nil, zero value otherwise.

### GetAnycastIpv4MaskOk

`func (o *ServicesPutRequestServiceValue) GetAnycastIpv4MaskOk() (*string, bool)`

GetAnycastIpv4MaskOk returns a tuple with the AnycastIpv4Mask field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnycastIpv4Mask

`func (o *ServicesPutRequestServiceValue) SetAnycastIpv4Mask(v string)`

SetAnycastIpv4Mask sets AnycastIpv4Mask field to given value.

### HasAnycastIpv4Mask

`func (o *ServicesPutRequestServiceValue) HasAnycastIpv4Mask() bool`

HasAnycastIpv4Mask returns a boolean if a field has been set.

### GetAnycastIpv6Mask

`func (o *ServicesPutRequestServiceValue) GetAnycastIpv6Mask() string`

GetAnycastIpv6Mask returns the AnycastIpv6Mask field if non-nil, zero value otherwise.

### GetAnycastIpv6MaskOk

`func (o *ServicesPutRequestServiceValue) GetAnycastIpv6MaskOk() (*string, bool)`

GetAnycastIpv6MaskOk returns a tuple with the AnycastIpv6Mask field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAnycastIpv6Mask

`func (o *ServicesPutRequestServiceValue) SetAnycastIpv6Mask(v string)`

SetAnycastIpv6Mask sets AnycastIpv6Mask field to given value.

### HasAnycastIpv6Mask

`func (o *ServicesPutRequestServiceValue) HasAnycastIpv6Mask() bool`

HasAnycastIpv6Mask returns a boolean if a field has been set.

### GetDhcpServerIpv4

`func (o *ServicesPutRequestServiceValue) GetDhcpServerIpv4() string`

GetDhcpServerIpv4 returns the DhcpServerIpv4 field if non-nil, zero value otherwise.

### GetDhcpServerIpv4Ok

`func (o *ServicesPutRequestServiceValue) GetDhcpServerIpv4Ok() (*string, bool)`

GetDhcpServerIpv4Ok returns a tuple with the DhcpServerIpv4 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDhcpServerIpv4

`func (o *ServicesPutRequestServiceValue) SetDhcpServerIpv4(v string)`

SetDhcpServerIpv4 sets DhcpServerIpv4 field to given value.

### HasDhcpServerIpv4

`func (o *ServicesPutRequestServiceValue) HasDhcpServerIpv4() bool`

HasDhcpServerIpv4 returns a boolean if a field has been set.

### GetDhcpServerIpv6

`func (o *ServicesPutRequestServiceValue) GetDhcpServerIpv6() string`

GetDhcpServerIpv6 returns the DhcpServerIpv6 field if non-nil, zero value otherwise.

### GetDhcpServerIpv6Ok

`func (o *ServicesPutRequestServiceValue) GetDhcpServerIpv6Ok() (*string, bool)`

GetDhcpServerIpv6Ok returns a tuple with the DhcpServerIpv6 field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDhcpServerIpv6

`func (o *ServicesPutRequestServiceValue) SetDhcpServerIpv6(v string)`

SetDhcpServerIpv6 sets DhcpServerIpv6 field to given value.

### HasDhcpServerIpv6

`func (o *ServicesPutRequestServiceValue) HasDhcpServerIpv6() bool`

HasDhcpServerIpv6 returns a boolean if a field has been set.

### GetMtu

`func (o *ServicesPutRequestServiceValue) GetMtu() int32`

GetMtu returns the Mtu field if non-nil, zero value otherwise.

### GetMtuOk

`func (o *ServicesPutRequestServiceValue) GetMtuOk() (*int32, bool)`

GetMtuOk returns a tuple with the Mtu field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMtu

`func (o *ServicesPutRequestServiceValue) SetMtu(v int32)`

SetMtu sets Mtu field to given value.

### HasMtu

`func (o *ServicesPutRequestServiceValue) HasMtu() bool`

HasMtu returns a boolean if a field has been set.

### SetMtuNil

`func (o *ServicesPutRequestServiceValue) SetMtuNil(b bool)`

 SetMtuNil sets the value for Mtu to be an explicit nil

### UnsetMtu
`func (o *ServicesPutRequestServiceValue) UnsetMtu()`

UnsetMtu ensures that no value is present for Mtu, not even an explicit nil
### GetObjectProperties

`func (o *ServicesPutRequestServiceValue) GetObjectProperties() ServicesPutRequestServiceValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ServicesPutRequestServiceValue) GetObjectPropertiesOk() (*ServicesPutRequestServiceValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ServicesPutRequestServiceValue) SetObjectProperties(v ServicesPutRequestServiceValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ServicesPutRequestServiceValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.

### GetMaxUpstreamRateMbps

`func (o *ServicesPutRequestServiceValue) GetMaxUpstreamRateMbps() int32`

GetMaxUpstreamRateMbps returns the MaxUpstreamRateMbps field if non-nil, zero value otherwise.

### GetMaxUpstreamRateMbpsOk

`func (o *ServicesPutRequestServiceValue) GetMaxUpstreamRateMbpsOk() (*int32, bool)`

GetMaxUpstreamRateMbpsOk returns a tuple with the MaxUpstreamRateMbps field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxUpstreamRateMbps

`func (o *ServicesPutRequestServiceValue) SetMaxUpstreamRateMbps(v int32)`

SetMaxUpstreamRateMbps sets MaxUpstreamRateMbps field to given value.

### HasMaxUpstreamRateMbps

`func (o *ServicesPutRequestServiceValue) HasMaxUpstreamRateMbps() bool`

HasMaxUpstreamRateMbps returns a boolean if a field has been set.

### GetMaxDownstreamRateMbps

`func (o *ServicesPutRequestServiceValue) GetMaxDownstreamRateMbps() int32`

GetMaxDownstreamRateMbps returns the MaxDownstreamRateMbps field if non-nil, zero value otherwise.

### GetMaxDownstreamRateMbpsOk

`func (o *ServicesPutRequestServiceValue) GetMaxDownstreamRateMbpsOk() (*int32, bool)`

GetMaxDownstreamRateMbpsOk returns a tuple with the MaxDownstreamRateMbps field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxDownstreamRateMbps

`func (o *ServicesPutRequestServiceValue) SetMaxDownstreamRateMbps(v int32)`

SetMaxDownstreamRateMbps sets MaxDownstreamRateMbps field to given value.

### HasMaxDownstreamRateMbps

`func (o *ServicesPutRequestServiceValue) HasMaxDownstreamRateMbps() bool`

HasMaxDownstreamRateMbps returns a boolean if a field has been set.

### GetPacketPriority

`func (o *ServicesPutRequestServiceValue) GetPacketPriority() string`

GetPacketPriority returns the PacketPriority field if non-nil, zero value otherwise.

### GetPacketPriorityOk

`func (o *ServicesPutRequestServiceValue) GetPacketPriorityOk() (*string, bool)`

GetPacketPriorityOk returns a tuple with the PacketPriority field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPacketPriority

`func (o *ServicesPutRequestServiceValue) SetPacketPriority(v string)`

SetPacketPriority sets PacketPriority field to given value.

### HasPacketPriority

`func (o *ServicesPutRequestServiceValue) HasPacketPriority() bool`

HasPacketPriority returns a boolean if a field has been set.

### GetMulticastManagementMode

`func (o *ServicesPutRequestServiceValue) GetMulticastManagementMode() string`

GetMulticastManagementMode returns the MulticastManagementMode field if non-nil, zero value otherwise.

### GetMulticastManagementModeOk

`func (o *ServicesPutRequestServiceValue) GetMulticastManagementModeOk() (*string, bool)`

GetMulticastManagementModeOk returns a tuple with the MulticastManagementMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMulticastManagementMode

`func (o *ServicesPutRequestServiceValue) SetMulticastManagementMode(v string)`

SetMulticastManagementMode sets MulticastManagementMode field to given value.

### HasMulticastManagementMode

`func (o *ServicesPutRequestServiceValue) HasMulticastManagementMode() bool`

HasMulticastManagementMode returns a boolean if a field has been set.

### GetTaggedPackets

`func (o *ServicesPutRequestServiceValue) GetTaggedPackets() bool`

GetTaggedPackets returns the TaggedPackets field if non-nil, zero value otherwise.

### GetTaggedPacketsOk

`func (o *ServicesPutRequestServiceValue) GetTaggedPacketsOk() (*bool, bool)`

GetTaggedPacketsOk returns a tuple with the TaggedPackets field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTaggedPackets

`func (o *ServicesPutRequestServiceValue) SetTaggedPackets(v bool)`

SetTaggedPackets sets TaggedPackets field to given value.

### HasTaggedPackets

`func (o *ServicesPutRequestServiceValue) HasTaggedPackets() bool`

HasTaggedPackets returns a boolean if a field has been set.

### GetTls

`func (o *ServicesPutRequestServiceValue) GetTls() bool`

GetTls returns the Tls field if non-nil, zero value otherwise.

### GetTlsOk

`func (o *ServicesPutRequestServiceValue) GetTlsOk() (*bool, bool)`

GetTlsOk returns a tuple with the Tls field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTls

`func (o *ServicesPutRequestServiceValue) SetTls(v bool)`

SetTls sets Tls field to given value.

### HasTls

`func (o *ServicesPutRequestServiceValue) HasTls() bool`

HasTls returns a boolean if a field has been set.

### GetAllowLocalSwitching

`func (o *ServicesPutRequestServiceValue) GetAllowLocalSwitching() bool`

GetAllowLocalSwitching returns the AllowLocalSwitching field if non-nil, zero value otherwise.

### GetAllowLocalSwitchingOk

`func (o *ServicesPutRequestServiceValue) GetAllowLocalSwitchingOk() (*bool, bool)`

GetAllowLocalSwitchingOk returns a tuple with the AllowLocalSwitching field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllowLocalSwitching

`func (o *ServicesPutRequestServiceValue) SetAllowLocalSwitching(v bool)`

SetAllowLocalSwitching sets AllowLocalSwitching field to given value.

### HasAllowLocalSwitching

`func (o *ServicesPutRequestServiceValue) HasAllowLocalSwitching() bool`

HasAllowLocalSwitching returns a boolean if a field has been set.

### GetActAsMulticastQuerier

`func (o *ServicesPutRequestServiceValue) GetActAsMulticastQuerier() bool`

GetActAsMulticastQuerier returns the ActAsMulticastQuerier field if non-nil, zero value otherwise.

### GetActAsMulticastQuerierOk

`func (o *ServicesPutRequestServiceValue) GetActAsMulticastQuerierOk() (*bool, bool)`

GetActAsMulticastQuerierOk returns a tuple with the ActAsMulticastQuerier field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetActAsMulticastQuerier

`func (o *ServicesPutRequestServiceValue) SetActAsMulticastQuerier(v bool)`

SetActAsMulticastQuerier sets ActAsMulticastQuerier field to given value.

### HasActAsMulticastQuerier

`func (o *ServicesPutRequestServiceValue) HasActAsMulticastQuerier() bool`

HasActAsMulticastQuerier returns a boolean if a field has been set.

### GetBlockUnknownUnicastFlood

`func (o *ServicesPutRequestServiceValue) GetBlockUnknownUnicastFlood() bool`

GetBlockUnknownUnicastFlood returns the BlockUnknownUnicastFlood field if non-nil, zero value otherwise.

### GetBlockUnknownUnicastFloodOk

`func (o *ServicesPutRequestServiceValue) GetBlockUnknownUnicastFloodOk() (*bool, bool)`

GetBlockUnknownUnicastFloodOk returns a tuple with the BlockUnknownUnicastFlood field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBlockUnknownUnicastFlood

`func (o *ServicesPutRequestServiceValue) SetBlockUnknownUnicastFlood(v bool)`

SetBlockUnknownUnicastFlood sets BlockUnknownUnicastFlood field to given value.

### HasBlockUnknownUnicastFlood

`func (o *ServicesPutRequestServiceValue) HasBlockUnknownUnicastFlood() bool`

HasBlockUnknownUnicastFlood returns a boolean if a field has been set.

### GetBlockDownstreamDhcpServer

`func (o *ServicesPutRequestServiceValue) GetBlockDownstreamDhcpServer() bool`

GetBlockDownstreamDhcpServer returns the BlockDownstreamDhcpServer field if non-nil, zero value otherwise.

### GetBlockDownstreamDhcpServerOk

`func (o *ServicesPutRequestServiceValue) GetBlockDownstreamDhcpServerOk() (*bool, bool)`

GetBlockDownstreamDhcpServerOk returns a tuple with the BlockDownstreamDhcpServer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBlockDownstreamDhcpServer

`func (o *ServicesPutRequestServiceValue) SetBlockDownstreamDhcpServer(v bool)`

SetBlockDownstreamDhcpServer sets BlockDownstreamDhcpServer field to given value.

### HasBlockDownstreamDhcpServer

`func (o *ServicesPutRequestServiceValue) HasBlockDownstreamDhcpServer() bool`

HasBlockDownstreamDhcpServer returns a boolean if a field has been set.

### GetIsManagementService

`func (o *ServicesPutRequestServiceValue) GetIsManagementService() bool`

GetIsManagementService returns the IsManagementService field if non-nil, zero value otherwise.

### GetIsManagementServiceOk

`func (o *ServicesPutRequestServiceValue) GetIsManagementServiceOk() (*bool, bool)`

GetIsManagementServiceOk returns a tuple with the IsManagementService field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsManagementService

`func (o *ServicesPutRequestServiceValue) SetIsManagementService(v bool)`

SetIsManagementService sets IsManagementService field to given value.

### HasIsManagementService

`func (o *ServicesPutRequestServiceValue) HasIsManagementService() bool`

HasIsManagementService returns a boolean if a field has been set.

### GetUseDscpToPBitMappingForL3PacketsIfAvailable

`func (o *ServicesPutRequestServiceValue) GetUseDscpToPBitMappingForL3PacketsIfAvailable() bool`

GetUseDscpToPBitMappingForL3PacketsIfAvailable returns the UseDscpToPBitMappingForL3PacketsIfAvailable field if non-nil, zero value otherwise.

### GetUseDscpToPBitMappingForL3PacketsIfAvailableOk

`func (o *ServicesPutRequestServiceValue) GetUseDscpToPBitMappingForL3PacketsIfAvailableOk() (*bool, bool)`

GetUseDscpToPBitMappingForL3PacketsIfAvailableOk returns a tuple with the UseDscpToPBitMappingForL3PacketsIfAvailable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUseDscpToPBitMappingForL3PacketsIfAvailable

`func (o *ServicesPutRequestServiceValue) SetUseDscpToPBitMappingForL3PacketsIfAvailable(v bool)`

SetUseDscpToPBitMappingForL3PacketsIfAvailable sets UseDscpToPBitMappingForL3PacketsIfAvailable field to given value.

### HasUseDscpToPBitMappingForL3PacketsIfAvailable

`func (o *ServicesPutRequestServiceValue) HasUseDscpToPBitMappingForL3PacketsIfAvailable() bool`

HasUseDscpToPBitMappingForL3PacketsIfAvailable returns a boolean if a field has been set.

### GetAllowFastLeave

`func (o *ServicesPutRequestServiceValue) GetAllowFastLeave() bool`

GetAllowFastLeave returns the AllowFastLeave field if non-nil, zero value otherwise.

### GetAllowFastLeaveOk

`func (o *ServicesPutRequestServiceValue) GetAllowFastLeaveOk() (*bool, bool)`

GetAllowFastLeaveOk returns a tuple with the AllowFastLeave field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAllowFastLeave

`func (o *ServicesPutRequestServiceValue) SetAllowFastLeave(v bool)`

SetAllowFastLeave sets AllowFastLeave field to given value.

### HasAllowFastLeave

`func (o *ServicesPutRequestServiceValue) HasAllowFastLeave() bool`

HasAllowFastLeave returns a boolean if a field has been set.

### GetMstInstance

`func (o *ServicesPutRequestServiceValue) GetMstInstance() int32`

GetMstInstance returns the MstInstance field if non-nil, zero value otherwise.

### GetMstInstanceOk

`func (o *ServicesPutRequestServiceValue) GetMstInstanceOk() (*int32, bool)`

GetMstInstanceOk returns a tuple with the MstInstance field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMstInstance

`func (o *ServicesPutRequestServiceValue) SetMstInstance(v int32)`

SetMstInstance sets MstInstance field to given value.

### HasMstInstance

`func (o *ServicesPutRequestServiceValue) HasMstInstance() bool`

HasMstInstance returns a boolean if a field has been set.

### GetPolicyBasedRouting

`func (o *ServicesPutRequestServiceValue) GetPolicyBasedRouting() string`

GetPolicyBasedRouting returns the PolicyBasedRouting field if non-nil, zero value otherwise.

### GetPolicyBasedRoutingOk

`func (o *ServicesPutRequestServiceValue) GetPolicyBasedRoutingOk() (*string, bool)`

GetPolicyBasedRoutingOk returns a tuple with the PolicyBasedRouting field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPolicyBasedRouting

`func (o *ServicesPutRequestServiceValue) SetPolicyBasedRouting(v string)`

SetPolicyBasedRouting sets PolicyBasedRouting field to given value.

### HasPolicyBasedRouting

`func (o *ServicesPutRequestServiceValue) HasPolicyBasedRouting() bool`

HasPolicyBasedRouting returns a boolean if a field has been set.

### GetPolicyBasedRoutingRefType

`func (o *ServicesPutRequestServiceValue) GetPolicyBasedRoutingRefType() string`

GetPolicyBasedRoutingRefType returns the PolicyBasedRoutingRefType field if non-nil, zero value otherwise.

### GetPolicyBasedRoutingRefTypeOk

`func (o *ServicesPutRequestServiceValue) GetPolicyBasedRoutingRefTypeOk() (*string, bool)`

GetPolicyBasedRoutingRefTypeOk returns a tuple with the PolicyBasedRoutingRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPolicyBasedRoutingRefType

`func (o *ServicesPutRequestServiceValue) SetPolicyBasedRoutingRefType(v string)`

SetPolicyBasedRoutingRefType sets PolicyBasedRoutingRefType field to given value.

### HasPolicyBasedRoutingRefType

`func (o *ServicesPutRequestServiceValue) HasPolicyBasedRoutingRefType() bool`

HasPolicyBasedRoutingRefType returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


