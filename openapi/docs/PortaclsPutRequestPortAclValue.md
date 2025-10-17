# PortaclsPutRequestPortAclValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**Ipv4Permit** | Pointer to [**[]PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner**](PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner.md) |  | [optional] 
**Ipv4Deny** | Pointer to [**[]PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner**](PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner.md) |  | [optional] 
**Ipv6Permit** | Pointer to [**[]PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner**](PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner.md) |  | [optional] 
**Ipv6Deny** | Pointer to [**[]PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner**](PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner.md) |  | [optional] 

## Methods

### NewPortaclsPutRequestPortAclValue

`func NewPortaclsPutRequestPortAclValue() *PortaclsPutRequestPortAclValue`

NewPortaclsPutRequestPortAclValue instantiates a new PortaclsPutRequestPortAclValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPortaclsPutRequestPortAclValueWithDefaults

`func NewPortaclsPutRequestPortAclValueWithDefaults() *PortaclsPutRequestPortAclValue`

NewPortaclsPutRequestPortAclValueWithDefaults instantiates a new PortaclsPutRequestPortAclValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *PortaclsPutRequestPortAclValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *PortaclsPutRequestPortAclValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *PortaclsPutRequestPortAclValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *PortaclsPutRequestPortAclValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *PortaclsPutRequestPortAclValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *PortaclsPutRequestPortAclValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *PortaclsPutRequestPortAclValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *PortaclsPutRequestPortAclValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetIpv4Permit

`func (o *PortaclsPutRequestPortAclValue) GetIpv4Permit() []PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner`

GetIpv4Permit returns the Ipv4Permit field if non-nil, zero value otherwise.

### GetIpv4PermitOk

`func (o *PortaclsPutRequestPortAclValue) GetIpv4PermitOk() (*[]PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner, bool)`

GetIpv4PermitOk returns a tuple with the Ipv4Permit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpv4Permit

`func (o *PortaclsPutRequestPortAclValue) SetIpv4Permit(v []PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner)`

SetIpv4Permit sets Ipv4Permit field to given value.

### HasIpv4Permit

`func (o *PortaclsPutRequestPortAclValue) HasIpv4Permit() bool`

HasIpv4Permit returns a boolean if a field has been set.

### GetIpv4Deny

`func (o *PortaclsPutRequestPortAclValue) GetIpv4Deny() []PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner`

GetIpv4Deny returns the Ipv4Deny field if non-nil, zero value otherwise.

### GetIpv4DenyOk

`func (o *PortaclsPutRequestPortAclValue) GetIpv4DenyOk() (*[]PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner, bool)`

GetIpv4DenyOk returns a tuple with the Ipv4Deny field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpv4Deny

`func (o *PortaclsPutRequestPortAclValue) SetIpv4Deny(v []PolicybasedroutingaclPutRequestPbRoutingAclValueIpv4PermitInner)`

SetIpv4Deny sets Ipv4Deny field to given value.

### HasIpv4Deny

`func (o *PortaclsPutRequestPortAclValue) HasIpv4Deny() bool`

HasIpv4Deny returns a boolean if a field has been set.

### GetIpv6Permit

`func (o *PortaclsPutRequestPortAclValue) GetIpv6Permit() []PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner`

GetIpv6Permit returns the Ipv6Permit field if non-nil, zero value otherwise.

### GetIpv6PermitOk

`func (o *PortaclsPutRequestPortAclValue) GetIpv6PermitOk() (*[]PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner, bool)`

GetIpv6PermitOk returns a tuple with the Ipv6Permit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpv6Permit

`func (o *PortaclsPutRequestPortAclValue) SetIpv6Permit(v []PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner)`

SetIpv6Permit sets Ipv6Permit field to given value.

### HasIpv6Permit

`func (o *PortaclsPutRequestPortAclValue) HasIpv6Permit() bool`

HasIpv6Permit returns a boolean if a field has been set.

### GetIpv6Deny

`func (o *PortaclsPutRequestPortAclValue) GetIpv6Deny() []PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner`

GetIpv6Deny returns the Ipv6Deny field if non-nil, zero value otherwise.

### GetIpv6DenyOk

`func (o *PortaclsPutRequestPortAclValue) GetIpv6DenyOk() (*[]PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner, bool)`

GetIpv6DenyOk returns a tuple with the Ipv6Deny field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpv6Deny

`func (o *PortaclsPutRequestPortAclValue) SetIpv6Deny(v []PolicybasedroutingaclPutRequestPbRoutingAclValueIpv6PermitInner)`

SetIpv6Deny sets Ipv6Deny field to given value.

### HasIpv6Deny

`func (o *PortaclsPutRequestPortAclValue) HasIpv6Deny() bool`

HasIpv6Deny returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


