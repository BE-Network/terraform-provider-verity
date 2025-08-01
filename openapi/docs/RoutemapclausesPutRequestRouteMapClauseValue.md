# RoutemapclausesPutRequestRouteMapClauseValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable flag of this provisioning object | [optional] [default to false]
**PermitDeny** | Pointer to **string** | Action upon match of Community Strings. | [optional] [default to "permit"]
**MatchAsPathAccessList** | Pointer to **string** | Match AS Path Access List | [optional] [default to ""]
**MatchAsPathAccessListRefType** | Pointer to **string** | Object type for match_as_path_access_list field | [optional] 
**MatchCommunityList** | Pointer to **string** | Match Community List | [optional] [default to ""]
**MatchCommunityListRefType** | Pointer to **string** | Object type for match_community_list field | [optional] 
**MatchExtendedCommunityList** | Pointer to **string** | Match Extended Community List | [optional] [default to ""]
**MatchExtendedCommunityListRefType** | Pointer to **string** | Object type for match_extended_community_list field | [optional] 
**MatchInterfaceNumber** | Pointer to **NullableInt32** | Match Interface Number | [optional] 
**MatchInterfaceVlan** | Pointer to **NullableInt32** | Match Interface VLAN | [optional] 
**MatchIpv4AddressIpPrefixList** | Pointer to **string** | Match IPv4 Address IPv4 Prefix List | [optional] [default to ""]
**MatchIpv4AddressIpPrefixListRefType** | Pointer to **string** | Object type for match_ipv4_address_ip_prefix_list field | [optional] 
**MatchIpv4NextHopIpPrefixList** | Pointer to **string** | Match IPv4 Next Hop IPv4 Prefix List | [optional] [default to ""]
**MatchIpv4NextHopIpPrefixListRefType** | Pointer to **string** | Object type for match_ipv4_next_hop_ip_prefix_list field | [optional] 
**MatchLocalPreference** | Pointer to **NullableInt32** | Match BGP Local Preference value on the route  | [optional] 
**MatchMetric** | Pointer to **NullableInt32** | Match Metric of the IP route entry  | [optional] 
**MatchOrigin** | Pointer to **string** | Match routes based on the value of the BGP Origin attribute  | [optional] [default to ""]
**MatchPeerIpAddress** | Pointer to **string** | Match BGP Peer IP Address the route was learned from  | [optional] [default to ""]
**MatchPeerInterface** | Pointer to **NullableInt32** | Match BGP Peer port the route was learned from  | [optional] 
**MatchPeerVlan** | Pointer to **NullableInt32** | Match BGP Peer VLAN over which the route was learned  | [optional] 
**MatchSourceProtocol** | Pointer to **string** | Match Routing  Protocol the route originated from  | [optional] [default to ""]
**MatchVrf** | Pointer to **string** | Match VRF the route is associated with  | [optional] [default to ""]
**MatchVrfRefType** | Pointer to **string** | Object type for match_vrf field | [optional] 
**MatchTag** | Pointer to **NullableInt32** | Match routes that have this value for a Tag attribute | [optional] 
**MatchEvpnRouteTypeDefault** | Pointer to **bool** | Match based on the type of EVPN Route Type being Default\&quot; | [optional] 
**MatchEvpnRouteType** | Pointer to **string** | Match based on the indicated EVPN Route Type | [optional] [default to ""]
**MatchVni** | Pointer to **NullableInt32** | Match based on the VNI value  | [optional] 
**ObjectProperties** | Pointer to [**RoutemapclausesPutRequestRouteMapClauseValueObjectProperties**](RoutemapclausesPutRequestRouteMapClauseValueObjectProperties.md) |  | [optional] 
**MatchIpv6AddressIpv6PrefixList** | Pointer to **string** | Match IPv4 Address IPv6 Prefix List | [optional] [default to ""]
**MatchIpv6AddressIpv6PrefixListRefType** | Pointer to **string** | Object type for match_ipv6_address_ipv6_prefix_list field | [optional] 
**MatchIpv6NextHopIpv6PrefixList** | Pointer to **string** | Match IPv6 Next Hop IPv6 Prefix List | [optional] [default to ""]
**MatchIpv6NextHopIpv6PrefixListRefType** | Pointer to **string** | Object type for match_ipv6_next_hop_ipv6_prefix_list field | [optional] 

## Methods

### NewRoutemapclausesPutRequestRouteMapClauseValue

`func NewRoutemapclausesPutRequestRouteMapClauseValue() *RoutemapclausesPutRequestRouteMapClauseValue`

NewRoutemapclausesPutRequestRouteMapClauseValue instantiates a new RoutemapclausesPutRequestRouteMapClauseValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRoutemapclausesPutRequestRouteMapClauseValueWithDefaults

`func NewRoutemapclausesPutRequestRouteMapClauseValueWithDefaults() *RoutemapclausesPutRequestRouteMapClauseValue`

NewRoutemapclausesPutRequestRouteMapClauseValueWithDefaults instantiates a new RoutemapclausesPutRequestRouteMapClauseValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetPermitDeny

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetPermitDeny() string`

GetPermitDeny returns the PermitDeny field if non-nil, zero value otherwise.

### GetPermitDenyOk

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetPermitDenyOk() (*string, bool)`

GetPermitDenyOk returns a tuple with the PermitDeny field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPermitDeny

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetPermitDeny(v string)`

SetPermitDeny sets PermitDeny field to given value.

### HasPermitDeny

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) HasPermitDeny() bool`

HasPermitDeny returns a boolean if a field has been set.

### GetMatchAsPathAccessList

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchAsPathAccessList() string`

GetMatchAsPathAccessList returns the MatchAsPathAccessList field if non-nil, zero value otherwise.

### GetMatchAsPathAccessListOk

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchAsPathAccessListOk() (*string, bool)`

GetMatchAsPathAccessListOk returns a tuple with the MatchAsPathAccessList field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchAsPathAccessList

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchAsPathAccessList(v string)`

SetMatchAsPathAccessList sets MatchAsPathAccessList field to given value.

### HasMatchAsPathAccessList

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) HasMatchAsPathAccessList() bool`

HasMatchAsPathAccessList returns a boolean if a field has been set.

### GetMatchAsPathAccessListRefType

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchAsPathAccessListRefType() string`

GetMatchAsPathAccessListRefType returns the MatchAsPathAccessListRefType field if non-nil, zero value otherwise.

### GetMatchAsPathAccessListRefTypeOk

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchAsPathAccessListRefTypeOk() (*string, bool)`

GetMatchAsPathAccessListRefTypeOk returns a tuple with the MatchAsPathAccessListRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchAsPathAccessListRefType

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchAsPathAccessListRefType(v string)`

SetMatchAsPathAccessListRefType sets MatchAsPathAccessListRefType field to given value.

### HasMatchAsPathAccessListRefType

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) HasMatchAsPathAccessListRefType() bool`

HasMatchAsPathAccessListRefType returns a boolean if a field has been set.

### GetMatchCommunityList

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchCommunityList() string`

GetMatchCommunityList returns the MatchCommunityList field if non-nil, zero value otherwise.

### GetMatchCommunityListOk

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchCommunityListOk() (*string, bool)`

GetMatchCommunityListOk returns a tuple with the MatchCommunityList field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchCommunityList

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchCommunityList(v string)`

SetMatchCommunityList sets MatchCommunityList field to given value.

### HasMatchCommunityList

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) HasMatchCommunityList() bool`

HasMatchCommunityList returns a boolean if a field has been set.

### GetMatchCommunityListRefType

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchCommunityListRefType() string`

GetMatchCommunityListRefType returns the MatchCommunityListRefType field if non-nil, zero value otherwise.

### GetMatchCommunityListRefTypeOk

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchCommunityListRefTypeOk() (*string, bool)`

GetMatchCommunityListRefTypeOk returns a tuple with the MatchCommunityListRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchCommunityListRefType

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchCommunityListRefType(v string)`

SetMatchCommunityListRefType sets MatchCommunityListRefType field to given value.

### HasMatchCommunityListRefType

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) HasMatchCommunityListRefType() bool`

HasMatchCommunityListRefType returns a boolean if a field has been set.

### GetMatchExtendedCommunityList

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchExtendedCommunityList() string`

GetMatchExtendedCommunityList returns the MatchExtendedCommunityList field if non-nil, zero value otherwise.

### GetMatchExtendedCommunityListOk

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchExtendedCommunityListOk() (*string, bool)`

GetMatchExtendedCommunityListOk returns a tuple with the MatchExtendedCommunityList field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchExtendedCommunityList

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchExtendedCommunityList(v string)`

SetMatchExtendedCommunityList sets MatchExtendedCommunityList field to given value.

### HasMatchExtendedCommunityList

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) HasMatchExtendedCommunityList() bool`

HasMatchExtendedCommunityList returns a boolean if a field has been set.

### GetMatchExtendedCommunityListRefType

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchExtendedCommunityListRefType() string`

GetMatchExtendedCommunityListRefType returns the MatchExtendedCommunityListRefType field if non-nil, zero value otherwise.

### GetMatchExtendedCommunityListRefTypeOk

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchExtendedCommunityListRefTypeOk() (*string, bool)`

GetMatchExtendedCommunityListRefTypeOk returns a tuple with the MatchExtendedCommunityListRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchExtendedCommunityListRefType

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchExtendedCommunityListRefType(v string)`

SetMatchExtendedCommunityListRefType sets MatchExtendedCommunityListRefType field to given value.

### HasMatchExtendedCommunityListRefType

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) HasMatchExtendedCommunityListRefType() bool`

HasMatchExtendedCommunityListRefType returns a boolean if a field has been set.

### GetMatchInterfaceNumber

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchInterfaceNumber() int32`

GetMatchInterfaceNumber returns the MatchInterfaceNumber field if non-nil, zero value otherwise.

### GetMatchInterfaceNumberOk

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchInterfaceNumberOk() (*int32, bool)`

GetMatchInterfaceNumberOk returns a tuple with the MatchInterfaceNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchInterfaceNumber

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchInterfaceNumber(v int32)`

SetMatchInterfaceNumber sets MatchInterfaceNumber field to given value.

### HasMatchInterfaceNumber

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) HasMatchInterfaceNumber() bool`

HasMatchInterfaceNumber returns a boolean if a field has been set.

### SetMatchInterfaceNumberNil

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchInterfaceNumberNil(b bool)`

 SetMatchInterfaceNumberNil sets the value for MatchInterfaceNumber to be an explicit nil

### UnsetMatchInterfaceNumber
`func (o *RoutemapclausesPutRequestRouteMapClauseValue) UnsetMatchInterfaceNumber()`

UnsetMatchInterfaceNumber ensures that no value is present for MatchInterfaceNumber, not even an explicit nil
### GetMatchInterfaceVlan

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchInterfaceVlan() int32`

GetMatchInterfaceVlan returns the MatchInterfaceVlan field if non-nil, zero value otherwise.

### GetMatchInterfaceVlanOk

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchInterfaceVlanOk() (*int32, bool)`

GetMatchInterfaceVlanOk returns a tuple with the MatchInterfaceVlan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchInterfaceVlan

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchInterfaceVlan(v int32)`

SetMatchInterfaceVlan sets MatchInterfaceVlan field to given value.

### HasMatchInterfaceVlan

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) HasMatchInterfaceVlan() bool`

HasMatchInterfaceVlan returns a boolean if a field has been set.

### SetMatchInterfaceVlanNil

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchInterfaceVlanNil(b bool)`

 SetMatchInterfaceVlanNil sets the value for MatchInterfaceVlan to be an explicit nil

### UnsetMatchInterfaceVlan
`func (o *RoutemapclausesPutRequestRouteMapClauseValue) UnsetMatchInterfaceVlan()`

UnsetMatchInterfaceVlan ensures that no value is present for MatchInterfaceVlan, not even an explicit nil
### GetMatchIpv4AddressIpPrefixList

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchIpv4AddressIpPrefixList() string`

GetMatchIpv4AddressIpPrefixList returns the MatchIpv4AddressIpPrefixList field if non-nil, zero value otherwise.

### GetMatchIpv4AddressIpPrefixListOk

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchIpv4AddressIpPrefixListOk() (*string, bool)`

GetMatchIpv4AddressIpPrefixListOk returns a tuple with the MatchIpv4AddressIpPrefixList field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchIpv4AddressIpPrefixList

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchIpv4AddressIpPrefixList(v string)`

SetMatchIpv4AddressIpPrefixList sets MatchIpv4AddressIpPrefixList field to given value.

### HasMatchIpv4AddressIpPrefixList

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) HasMatchIpv4AddressIpPrefixList() bool`

HasMatchIpv4AddressIpPrefixList returns a boolean if a field has been set.

### GetMatchIpv4AddressIpPrefixListRefType

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchIpv4AddressIpPrefixListRefType() string`

GetMatchIpv4AddressIpPrefixListRefType returns the MatchIpv4AddressIpPrefixListRefType field if non-nil, zero value otherwise.

### GetMatchIpv4AddressIpPrefixListRefTypeOk

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchIpv4AddressIpPrefixListRefTypeOk() (*string, bool)`

GetMatchIpv4AddressIpPrefixListRefTypeOk returns a tuple with the MatchIpv4AddressIpPrefixListRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchIpv4AddressIpPrefixListRefType

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchIpv4AddressIpPrefixListRefType(v string)`

SetMatchIpv4AddressIpPrefixListRefType sets MatchIpv4AddressIpPrefixListRefType field to given value.

### HasMatchIpv4AddressIpPrefixListRefType

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) HasMatchIpv4AddressIpPrefixListRefType() bool`

HasMatchIpv4AddressIpPrefixListRefType returns a boolean if a field has been set.

### GetMatchIpv4NextHopIpPrefixList

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchIpv4NextHopIpPrefixList() string`

GetMatchIpv4NextHopIpPrefixList returns the MatchIpv4NextHopIpPrefixList field if non-nil, zero value otherwise.

### GetMatchIpv4NextHopIpPrefixListOk

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchIpv4NextHopIpPrefixListOk() (*string, bool)`

GetMatchIpv4NextHopIpPrefixListOk returns a tuple with the MatchIpv4NextHopIpPrefixList field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchIpv4NextHopIpPrefixList

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchIpv4NextHopIpPrefixList(v string)`

SetMatchIpv4NextHopIpPrefixList sets MatchIpv4NextHopIpPrefixList field to given value.

### HasMatchIpv4NextHopIpPrefixList

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) HasMatchIpv4NextHopIpPrefixList() bool`

HasMatchIpv4NextHopIpPrefixList returns a boolean if a field has been set.

### GetMatchIpv4NextHopIpPrefixListRefType

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchIpv4NextHopIpPrefixListRefType() string`

GetMatchIpv4NextHopIpPrefixListRefType returns the MatchIpv4NextHopIpPrefixListRefType field if non-nil, zero value otherwise.

### GetMatchIpv4NextHopIpPrefixListRefTypeOk

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchIpv4NextHopIpPrefixListRefTypeOk() (*string, bool)`

GetMatchIpv4NextHopIpPrefixListRefTypeOk returns a tuple with the MatchIpv4NextHopIpPrefixListRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchIpv4NextHopIpPrefixListRefType

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchIpv4NextHopIpPrefixListRefType(v string)`

SetMatchIpv4NextHopIpPrefixListRefType sets MatchIpv4NextHopIpPrefixListRefType field to given value.

### HasMatchIpv4NextHopIpPrefixListRefType

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) HasMatchIpv4NextHopIpPrefixListRefType() bool`

HasMatchIpv4NextHopIpPrefixListRefType returns a boolean if a field has been set.

### GetMatchLocalPreference

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchLocalPreference() int32`

GetMatchLocalPreference returns the MatchLocalPreference field if non-nil, zero value otherwise.

### GetMatchLocalPreferenceOk

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchLocalPreferenceOk() (*int32, bool)`

GetMatchLocalPreferenceOk returns a tuple with the MatchLocalPreference field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchLocalPreference

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchLocalPreference(v int32)`

SetMatchLocalPreference sets MatchLocalPreference field to given value.

### HasMatchLocalPreference

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) HasMatchLocalPreference() bool`

HasMatchLocalPreference returns a boolean if a field has been set.

### SetMatchLocalPreferenceNil

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchLocalPreferenceNil(b bool)`

 SetMatchLocalPreferenceNil sets the value for MatchLocalPreference to be an explicit nil

### UnsetMatchLocalPreference
`func (o *RoutemapclausesPutRequestRouteMapClauseValue) UnsetMatchLocalPreference()`

UnsetMatchLocalPreference ensures that no value is present for MatchLocalPreference, not even an explicit nil
### GetMatchMetric

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchMetric() int32`

GetMatchMetric returns the MatchMetric field if non-nil, zero value otherwise.

### GetMatchMetricOk

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchMetricOk() (*int32, bool)`

GetMatchMetricOk returns a tuple with the MatchMetric field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchMetric

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchMetric(v int32)`

SetMatchMetric sets MatchMetric field to given value.

### HasMatchMetric

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) HasMatchMetric() bool`

HasMatchMetric returns a boolean if a field has been set.

### SetMatchMetricNil

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchMetricNil(b bool)`

 SetMatchMetricNil sets the value for MatchMetric to be an explicit nil

### UnsetMatchMetric
`func (o *RoutemapclausesPutRequestRouteMapClauseValue) UnsetMatchMetric()`

UnsetMatchMetric ensures that no value is present for MatchMetric, not even an explicit nil
### GetMatchOrigin

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchOrigin() string`

GetMatchOrigin returns the MatchOrigin field if non-nil, zero value otherwise.

### GetMatchOriginOk

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchOriginOk() (*string, bool)`

GetMatchOriginOk returns a tuple with the MatchOrigin field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchOrigin

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchOrigin(v string)`

SetMatchOrigin sets MatchOrigin field to given value.

### HasMatchOrigin

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) HasMatchOrigin() bool`

HasMatchOrigin returns a boolean if a field has been set.

### GetMatchPeerIpAddress

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchPeerIpAddress() string`

GetMatchPeerIpAddress returns the MatchPeerIpAddress field if non-nil, zero value otherwise.

### GetMatchPeerIpAddressOk

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchPeerIpAddressOk() (*string, bool)`

GetMatchPeerIpAddressOk returns a tuple with the MatchPeerIpAddress field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchPeerIpAddress

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchPeerIpAddress(v string)`

SetMatchPeerIpAddress sets MatchPeerIpAddress field to given value.

### HasMatchPeerIpAddress

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) HasMatchPeerIpAddress() bool`

HasMatchPeerIpAddress returns a boolean if a field has been set.

### GetMatchPeerInterface

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchPeerInterface() int32`

GetMatchPeerInterface returns the MatchPeerInterface field if non-nil, zero value otherwise.

### GetMatchPeerInterfaceOk

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchPeerInterfaceOk() (*int32, bool)`

GetMatchPeerInterfaceOk returns a tuple with the MatchPeerInterface field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchPeerInterface

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchPeerInterface(v int32)`

SetMatchPeerInterface sets MatchPeerInterface field to given value.

### HasMatchPeerInterface

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) HasMatchPeerInterface() bool`

HasMatchPeerInterface returns a boolean if a field has been set.

### SetMatchPeerInterfaceNil

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchPeerInterfaceNil(b bool)`

 SetMatchPeerInterfaceNil sets the value for MatchPeerInterface to be an explicit nil

### UnsetMatchPeerInterface
`func (o *RoutemapclausesPutRequestRouteMapClauseValue) UnsetMatchPeerInterface()`

UnsetMatchPeerInterface ensures that no value is present for MatchPeerInterface, not even an explicit nil
### GetMatchPeerVlan

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchPeerVlan() int32`

GetMatchPeerVlan returns the MatchPeerVlan field if non-nil, zero value otherwise.

### GetMatchPeerVlanOk

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchPeerVlanOk() (*int32, bool)`

GetMatchPeerVlanOk returns a tuple with the MatchPeerVlan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchPeerVlan

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchPeerVlan(v int32)`

SetMatchPeerVlan sets MatchPeerVlan field to given value.

### HasMatchPeerVlan

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) HasMatchPeerVlan() bool`

HasMatchPeerVlan returns a boolean if a field has been set.

### SetMatchPeerVlanNil

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchPeerVlanNil(b bool)`

 SetMatchPeerVlanNil sets the value for MatchPeerVlan to be an explicit nil

### UnsetMatchPeerVlan
`func (o *RoutemapclausesPutRequestRouteMapClauseValue) UnsetMatchPeerVlan()`

UnsetMatchPeerVlan ensures that no value is present for MatchPeerVlan, not even an explicit nil
### GetMatchSourceProtocol

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchSourceProtocol() string`

GetMatchSourceProtocol returns the MatchSourceProtocol field if non-nil, zero value otherwise.

### GetMatchSourceProtocolOk

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchSourceProtocolOk() (*string, bool)`

GetMatchSourceProtocolOk returns a tuple with the MatchSourceProtocol field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchSourceProtocol

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchSourceProtocol(v string)`

SetMatchSourceProtocol sets MatchSourceProtocol field to given value.

### HasMatchSourceProtocol

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) HasMatchSourceProtocol() bool`

HasMatchSourceProtocol returns a boolean if a field has been set.

### GetMatchVrf

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchVrf() string`

GetMatchVrf returns the MatchVrf field if non-nil, zero value otherwise.

### GetMatchVrfOk

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchVrfOk() (*string, bool)`

GetMatchVrfOk returns a tuple with the MatchVrf field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchVrf

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchVrf(v string)`

SetMatchVrf sets MatchVrf field to given value.

### HasMatchVrf

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) HasMatchVrf() bool`

HasMatchVrf returns a boolean if a field has been set.

### GetMatchVrfRefType

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchVrfRefType() string`

GetMatchVrfRefType returns the MatchVrfRefType field if non-nil, zero value otherwise.

### GetMatchVrfRefTypeOk

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchVrfRefTypeOk() (*string, bool)`

GetMatchVrfRefTypeOk returns a tuple with the MatchVrfRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchVrfRefType

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchVrfRefType(v string)`

SetMatchVrfRefType sets MatchVrfRefType field to given value.

### HasMatchVrfRefType

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) HasMatchVrfRefType() bool`

HasMatchVrfRefType returns a boolean if a field has been set.

### GetMatchTag

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchTag() int32`

GetMatchTag returns the MatchTag field if non-nil, zero value otherwise.

### GetMatchTagOk

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchTagOk() (*int32, bool)`

GetMatchTagOk returns a tuple with the MatchTag field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchTag

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchTag(v int32)`

SetMatchTag sets MatchTag field to given value.

### HasMatchTag

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) HasMatchTag() bool`

HasMatchTag returns a boolean if a field has been set.

### SetMatchTagNil

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchTagNil(b bool)`

 SetMatchTagNil sets the value for MatchTag to be an explicit nil

### UnsetMatchTag
`func (o *RoutemapclausesPutRequestRouteMapClauseValue) UnsetMatchTag()`

UnsetMatchTag ensures that no value is present for MatchTag, not even an explicit nil
### GetMatchEvpnRouteTypeDefault

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchEvpnRouteTypeDefault() bool`

GetMatchEvpnRouteTypeDefault returns the MatchEvpnRouteTypeDefault field if non-nil, zero value otherwise.

### GetMatchEvpnRouteTypeDefaultOk

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchEvpnRouteTypeDefaultOk() (*bool, bool)`

GetMatchEvpnRouteTypeDefaultOk returns a tuple with the MatchEvpnRouteTypeDefault field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchEvpnRouteTypeDefault

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchEvpnRouteTypeDefault(v bool)`

SetMatchEvpnRouteTypeDefault sets MatchEvpnRouteTypeDefault field to given value.

### HasMatchEvpnRouteTypeDefault

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) HasMatchEvpnRouteTypeDefault() bool`

HasMatchEvpnRouteTypeDefault returns a boolean if a field has been set.

### GetMatchEvpnRouteType

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchEvpnRouteType() string`

GetMatchEvpnRouteType returns the MatchEvpnRouteType field if non-nil, zero value otherwise.

### GetMatchEvpnRouteTypeOk

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchEvpnRouteTypeOk() (*string, bool)`

GetMatchEvpnRouteTypeOk returns a tuple with the MatchEvpnRouteType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchEvpnRouteType

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchEvpnRouteType(v string)`

SetMatchEvpnRouteType sets MatchEvpnRouteType field to given value.

### HasMatchEvpnRouteType

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) HasMatchEvpnRouteType() bool`

HasMatchEvpnRouteType returns a boolean if a field has been set.

### GetMatchVni

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchVni() int32`

GetMatchVni returns the MatchVni field if non-nil, zero value otherwise.

### GetMatchVniOk

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchVniOk() (*int32, bool)`

GetMatchVniOk returns a tuple with the MatchVni field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchVni

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchVni(v int32)`

SetMatchVni sets MatchVni field to given value.

### HasMatchVni

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) HasMatchVni() bool`

HasMatchVni returns a boolean if a field has been set.

### SetMatchVniNil

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchVniNil(b bool)`

 SetMatchVniNil sets the value for MatchVni to be an explicit nil

### UnsetMatchVni
`func (o *RoutemapclausesPutRequestRouteMapClauseValue) UnsetMatchVni()`

UnsetMatchVni ensures that no value is present for MatchVni, not even an explicit nil
### GetObjectProperties

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetObjectProperties() RoutemapclausesPutRequestRouteMapClauseValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetObjectPropertiesOk() (*RoutemapclausesPutRequestRouteMapClauseValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetObjectProperties(v RoutemapclausesPutRequestRouteMapClauseValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.

### GetMatchIpv6AddressIpv6PrefixList

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchIpv6AddressIpv6PrefixList() string`

GetMatchIpv6AddressIpv6PrefixList returns the MatchIpv6AddressIpv6PrefixList field if non-nil, zero value otherwise.

### GetMatchIpv6AddressIpv6PrefixListOk

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchIpv6AddressIpv6PrefixListOk() (*string, bool)`

GetMatchIpv6AddressIpv6PrefixListOk returns a tuple with the MatchIpv6AddressIpv6PrefixList field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchIpv6AddressIpv6PrefixList

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchIpv6AddressIpv6PrefixList(v string)`

SetMatchIpv6AddressIpv6PrefixList sets MatchIpv6AddressIpv6PrefixList field to given value.

### HasMatchIpv6AddressIpv6PrefixList

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) HasMatchIpv6AddressIpv6PrefixList() bool`

HasMatchIpv6AddressIpv6PrefixList returns a boolean if a field has been set.

### GetMatchIpv6AddressIpv6PrefixListRefType

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchIpv6AddressIpv6PrefixListRefType() string`

GetMatchIpv6AddressIpv6PrefixListRefType returns the MatchIpv6AddressIpv6PrefixListRefType field if non-nil, zero value otherwise.

### GetMatchIpv6AddressIpv6PrefixListRefTypeOk

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchIpv6AddressIpv6PrefixListRefTypeOk() (*string, bool)`

GetMatchIpv6AddressIpv6PrefixListRefTypeOk returns a tuple with the MatchIpv6AddressIpv6PrefixListRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchIpv6AddressIpv6PrefixListRefType

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchIpv6AddressIpv6PrefixListRefType(v string)`

SetMatchIpv6AddressIpv6PrefixListRefType sets MatchIpv6AddressIpv6PrefixListRefType field to given value.

### HasMatchIpv6AddressIpv6PrefixListRefType

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) HasMatchIpv6AddressIpv6PrefixListRefType() bool`

HasMatchIpv6AddressIpv6PrefixListRefType returns a boolean if a field has been set.

### GetMatchIpv6NextHopIpv6PrefixList

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchIpv6NextHopIpv6PrefixList() string`

GetMatchIpv6NextHopIpv6PrefixList returns the MatchIpv6NextHopIpv6PrefixList field if non-nil, zero value otherwise.

### GetMatchIpv6NextHopIpv6PrefixListOk

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchIpv6NextHopIpv6PrefixListOk() (*string, bool)`

GetMatchIpv6NextHopIpv6PrefixListOk returns a tuple with the MatchIpv6NextHopIpv6PrefixList field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchIpv6NextHopIpv6PrefixList

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchIpv6NextHopIpv6PrefixList(v string)`

SetMatchIpv6NextHopIpv6PrefixList sets MatchIpv6NextHopIpv6PrefixList field to given value.

### HasMatchIpv6NextHopIpv6PrefixList

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) HasMatchIpv6NextHopIpv6PrefixList() bool`

HasMatchIpv6NextHopIpv6PrefixList returns a boolean if a field has been set.

### GetMatchIpv6NextHopIpv6PrefixListRefType

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchIpv6NextHopIpv6PrefixListRefType() string`

GetMatchIpv6NextHopIpv6PrefixListRefType returns the MatchIpv6NextHopIpv6PrefixListRefType field if non-nil, zero value otherwise.

### GetMatchIpv6NextHopIpv6PrefixListRefTypeOk

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) GetMatchIpv6NextHopIpv6PrefixListRefTypeOk() (*string, bool)`

GetMatchIpv6NextHopIpv6PrefixListRefTypeOk returns a tuple with the MatchIpv6NextHopIpv6PrefixListRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMatchIpv6NextHopIpv6PrefixListRefType

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) SetMatchIpv6NextHopIpv6PrefixListRefType(v string)`

SetMatchIpv6NextHopIpv6PrefixListRefType sets MatchIpv6NextHopIpv6PrefixListRefType field to given value.

### HasMatchIpv6NextHopIpv6PrefixListRefType

`func (o *RoutemapclausesPutRequestRouteMapClauseValue) HasMatchIpv6NextHopIpv6PrefixListRefType() bool`

HasMatchIpv6NextHopIpv6PrefixListRefType returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


