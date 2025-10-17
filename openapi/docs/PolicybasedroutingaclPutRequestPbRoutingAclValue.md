# PolicybasedroutingaclPutRequestPbRoutingAclValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**IpvProtocol** | Pointer to **string** | IPv4 or IPv6 | [optional] [default to "ipv4"]
**NextHopIps** | Pointer to **string** | Next hop IP addresses | [optional] [default to ""]
**Ipv4Permit** | Pointer to [**[]PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner**](PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner.md) |  | [optional] 
**Ipv4Deny** | Pointer to [**[]PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner**](PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner.md) |  | [optional] 
**Ipv6Permit** | Pointer to [**[]PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner**](PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner.md) |  | [optional] 
**Ipv6Deny** | Pointer to [**[]PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner**](PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner.md) |  | [optional] 

## Methods

### NewPolicybasedroutingaclPutRequestPbRoutingAclValue

`func NewPolicybasedroutingaclPutRequestPbRoutingAclValue() *PolicybasedroutingaclPutRequestPbRoutingAclValue`

NewPolicybasedroutingaclPutRequestPbRoutingAclValue instantiates a new PolicybasedroutingaclPutRequestPbRoutingAclValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPolicybasedroutingaclPutRequestPbRoutingAclValueWithDefaults

`func NewPolicybasedroutingaclPutRequestPbRoutingAclValueWithDefaults() *PolicybasedroutingaclPutRequestPbRoutingAclValue`

NewPolicybasedroutingaclPutRequestPbRoutingAclValueWithDefaults instantiates a new PolicybasedroutingaclPutRequestPbRoutingAclValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *PolicybasedroutingaclPutRequestPbRoutingAclValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *PolicybasedroutingaclPutRequestPbRoutingAclValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *PolicybasedroutingaclPutRequestPbRoutingAclValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *PolicybasedroutingaclPutRequestPbRoutingAclValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *PolicybasedroutingaclPutRequestPbRoutingAclValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *PolicybasedroutingaclPutRequestPbRoutingAclValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *PolicybasedroutingaclPutRequestPbRoutingAclValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *PolicybasedroutingaclPutRequestPbRoutingAclValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetIpvProtocol

`func (o *PolicybasedroutingaclPutRequestPbRoutingAclValue) GetIpvProtocol() string`

GetIpvProtocol returns the IpvProtocol field if non-nil, zero value otherwise.

### GetIpvProtocolOk

`func (o *PolicybasedroutingaclPutRequestPbRoutingAclValue) GetIpvProtocolOk() (*string, bool)`

GetIpvProtocolOk returns a tuple with the IpvProtocol field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpvProtocol

`func (o *PolicybasedroutingaclPutRequestPbRoutingAclValue) SetIpvProtocol(v string)`

SetIpvProtocol sets IpvProtocol field to given value.

### HasIpvProtocol

`func (o *PolicybasedroutingaclPutRequestPbRoutingAclValue) HasIpvProtocol() bool`

HasIpvProtocol returns a boolean if a field has been set.

### GetNextHopIps

`func (o *PolicybasedroutingaclPutRequestPbRoutingAclValue) GetNextHopIps() string`

GetNextHopIps returns the NextHopIps field if non-nil, zero value otherwise.

### GetNextHopIpsOk

`func (o *PolicybasedroutingaclPutRequestPbRoutingAclValue) GetNextHopIpsOk() (*string, bool)`

GetNextHopIpsOk returns a tuple with the NextHopIps field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNextHopIps

`func (o *PolicybasedroutingaclPutRequestPbRoutingAclValue) SetNextHopIps(v string)`

SetNextHopIps sets NextHopIps field to given value.

### HasNextHopIps

`func (o *PolicybasedroutingaclPutRequestPbRoutingAclValue) HasNextHopIps() bool`

HasNextHopIps returns a boolean if a field has been set.

### GetIpv4Permit

`func (o *PolicybasedroutingaclPutRequestPbRoutingAclValue) GetIpv4Permit() []PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner`

GetIpv4Permit returns the Ipv4Permit field if non-nil, zero value otherwise.

### GetIpv4PermitOk

`func (o *PolicybasedroutingaclPutRequestPbRoutingAclValue) GetIpv4PermitOk() (*[]PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner, bool)`

GetIpv4PermitOk returns a tuple with the Ipv4Permit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpv4Permit

`func (o *PolicybasedroutingaclPutRequestPbRoutingAclValue) SetIpv4Permit(v []PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner)`

SetIpv4Permit sets Ipv4Permit field to given value.

### HasIpv4Permit

`func (o *PolicybasedroutingaclPutRequestPbRoutingAclValue) HasIpv4Permit() bool`

HasIpv4Permit returns a boolean if a field has been set.

### GetIpv4Deny

`func (o *PolicybasedroutingaclPutRequestPbRoutingAclValue) GetIpv4Deny() []PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner`

GetIpv4Deny returns the Ipv4Deny field if non-nil, zero value otherwise.

### GetIpv4DenyOk

`func (o *PolicybasedroutingaclPutRequestPbRoutingAclValue) GetIpv4DenyOk() (*[]PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner, bool)`

GetIpv4DenyOk returns a tuple with the Ipv4Deny field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpv4Deny

`func (o *PolicybasedroutingaclPutRequestPbRoutingAclValue) SetIpv4Deny(v []PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner)`

SetIpv4Deny sets Ipv4Deny field to given value.

### HasIpv4Deny

`func (o *PolicybasedroutingaclPutRequestPbRoutingAclValue) HasIpv4Deny() bool`

HasIpv4Deny returns a boolean if a field has been set.

### GetIpv6Permit

`func (o *PolicybasedroutingaclPutRequestPbRoutingAclValue) GetIpv6Permit() []PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner`

GetIpv6Permit returns the Ipv6Permit field if non-nil, zero value otherwise.

### GetIpv6PermitOk

`func (o *PolicybasedroutingaclPutRequestPbRoutingAclValue) GetIpv6PermitOk() (*[]PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner, bool)`

GetIpv6PermitOk returns a tuple with the Ipv6Permit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpv6Permit

`func (o *PolicybasedroutingaclPutRequestPbRoutingAclValue) SetIpv6Permit(v []PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner)`

SetIpv6Permit sets Ipv6Permit field to given value.

### HasIpv6Permit

`func (o *PolicybasedroutingaclPutRequestPbRoutingAclValue) HasIpv6Permit() bool`

HasIpv6Permit returns a boolean if a field has been set.

### GetIpv6Deny

`func (o *PolicybasedroutingaclPutRequestPbRoutingAclValue) GetIpv6Deny() []PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner`

GetIpv6Deny returns the Ipv6Deny field if non-nil, zero value otherwise.

### GetIpv6DenyOk

`func (o *PolicybasedroutingaclPutRequestPbRoutingAclValue) GetIpv6DenyOk() (*[]PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner, bool)`

GetIpv6DenyOk returns a tuple with the Ipv6Deny field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpv6Deny

`func (o *PolicybasedroutingaclPutRequestPbRoutingAclValue) SetIpv6Deny(v []PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner)`

SetIpv6Deny sets Ipv6Deny field to given value.

### HasIpv6Deny

`func (o *PolicybasedroutingaclPutRequestPbRoutingAclValue) HasIpv6Deny() bool`

HasIpv6Deny returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


