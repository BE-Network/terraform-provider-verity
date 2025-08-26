# EthportprofilesPutRequestEthPortProfileValueServicesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**RowNumEnable** | Pointer to **bool** | Enable row | [optional] [default to false]
**RowNumService** | Pointer to **string** | Choose a Service to connect | [optional] [default to ""]
**RowNumServiceRefType** | Pointer to **string** | Object type for row_num_service field | [optional] 
**RowNumExternalVlan** | Pointer to **NullableInt32** | Choose an external vlan A value of 0 will make the VLAN untagged, while in case null is provided, the VLAN will be the one associated with the service. | [optional] 
**RowNumIngressAcl** | Pointer to **string** | Choose an ingress access control list | [optional] [default to ""]
**RowNumIngressAclRefType** | Pointer to **string** | Object type for row_num_ingress_acl field | [optional] 
**RowNumEgressAcl** | Pointer to **string** | Choose an egress access control list | [optional] [default to ""]
**RowNumEgressAclRefType** | Pointer to **string** | Object type for row_num_egress_acl field | [optional] 
**Index** | Pointer to **int32** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 
**RowNumMacFilter** | Pointer to **string** | Choose an access control list | [optional] [default to ""]
**RowNumMacFilterRefType** | Pointer to **string** | Object type for row_num_mac_filter field | [optional] 
**RowNumLanIptv** | Pointer to **string** | Denotes a LAN or IPTV service | [optional] [default to ""]

## Methods

### NewEthportprofilesPutRequestEthPortProfileValueServicesInner

`func NewEthportprofilesPutRequestEthPortProfileValueServicesInner() *EthportprofilesPutRequestEthPortProfileValueServicesInner`

NewEthportprofilesPutRequestEthPortProfileValueServicesInner instantiates a new EthportprofilesPutRequestEthPortProfileValueServicesInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewEthportprofilesPutRequestEthPortProfileValueServicesInnerWithDefaults

`func NewEthportprofilesPutRequestEthPortProfileValueServicesInnerWithDefaults() *EthportprofilesPutRequestEthPortProfileValueServicesInner`

NewEthportprofilesPutRequestEthPortProfileValueServicesInnerWithDefaults instantiates a new EthportprofilesPutRequestEthPortProfileValueServicesInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRowNumEnable

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) GetRowNumEnable() bool`

GetRowNumEnable returns the RowNumEnable field if non-nil, zero value otherwise.

### GetRowNumEnableOk

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) GetRowNumEnableOk() (*bool, bool)`

GetRowNumEnableOk returns a tuple with the RowNumEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowNumEnable

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) SetRowNumEnable(v bool)`

SetRowNumEnable sets RowNumEnable field to given value.

### HasRowNumEnable

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) HasRowNumEnable() bool`

HasRowNumEnable returns a boolean if a field has been set.

### GetRowNumService

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) GetRowNumService() string`

GetRowNumService returns the RowNumService field if non-nil, zero value otherwise.

### GetRowNumServiceOk

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) GetRowNumServiceOk() (*string, bool)`

GetRowNumServiceOk returns a tuple with the RowNumService field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowNumService

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) SetRowNumService(v string)`

SetRowNumService sets RowNumService field to given value.

### HasRowNumService

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) HasRowNumService() bool`

HasRowNumService returns a boolean if a field has been set.

### GetRowNumServiceRefType

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) GetRowNumServiceRefType() string`

GetRowNumServiceRefType returns the RowNumServiceRefType field if non-nil, zero value otherwise.

### GetRowNumServiceRefTypeOk

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) GetRowNumServiceRefTypeOk() (*string, bool)`

GetRowNumServiceRefTypeOk returns a tuple with the RowNumServiceRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowNumServiceRefType

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) SetRowNumServiceRefType(v string)`

SetRowNumServiceRefType sets RowNumServiceRefType field to given value.

### HasRowNumServiceRefType

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) HasRowNumServiceRefType() bool`

HasRowNumServiceRefType returns a boolean if a field has been set.

### GetRowNumExternalVlan

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) GetRowNumExternalVlan() int32`

GetRowNumExternalVlan returns the RowNumExternalVlan field if non-nil, zero value otherwise.

### GetRowNumExternalVlanOk

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) GetRowNumExternalVlanOk() (*int32, bool)`

GetRowNumExternalVlanOk returns a tuple with the RowNumExternalVlan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowNumExternalVlan

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) SetRowNumExternalVlan(v int32)`

SetRowNumExternalVlan sets RowNumExternalVlan field to given value.

### HasRowNumExternalVlan

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) HasRowNumExternalVlan() bool`

HasRowNumExternalVlan returns a boolean if a field has been set.

### SetRowNumExternalVlanNil

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) SetRowNumExternalVlanNil(b bool)`

 SetRowNumExternalVlanNil sets the value for RowNumExternalVlan to be an explicit nil

### UnsetRowNumExternalVlan
`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) UnsetRowNumExternalVlan()`

UnsetRowNumExternalVlan ensures that no value is present for RowNumExternalVlan, not even an explicit nil
### GetRowNumIngressAcl

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) GetRowNumIngressAcl() string`

GetRowNumIngressAcl returns the RowNumIngressAcl field if non-nil, zero value otherwise.

### GetRowNumIngressAclOk

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) GetRowNumIngressAclOk() (*string, bool)`

GetRowNumIngressAclOk returns a tuple with the RowNumIngressAcl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowNumIngressAcl

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) SetRowNumIngressAcl(v string)`

SetRowNumIngressAcl sets RowNumIngressAcl field to given value.

### HasRowNumIngressAcl

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) HasRowNumIngressAcl() bool`

HasRowNumIngressAcl returns a boolean if a field has been set.

### GetRowNumIngressAclRefType

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) GetRowNumIngressAclRefType() string`

GetRowNumIngressAclRefType returns the RowNumIngressAclRefType field if non-nil, zero value otherwise.

### GetRowNumIngressAclRefTypeOk

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) GetRowNumIngressAclRefTypeOk() (*string, bool)`

GetRowNumIngressAclRefTypeOk returns a tuple with the RowNumIngressAclRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowNumIngressAclRefType

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) SetRowNumIngressAclRefType(v string)`

SetRowNumIngressAclRefType sets RowNumIngressAclRefType field to given value.

### HasRowNumIngressAclRefType

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) HasRowNumIngressAclRefType() bool`

HasRowNumIngressAclRefType returns a boolean if a field has been set.

### GetRowNumEgressAcl

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) GetRowNumEgressAcl() string`

GetRowNumEgressAcl returns the RowNumEgressAcl field if non-nil, zero value otherwise.

### GetRowNumEgressAclOk

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) GetRowNumEgressAclOk() (*string, bool)`

GetRowNumEgressAclOk returns a tuple with the RowNumEgressAcl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowNumEgressAcl

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) SetRowNumEgressAcl(v string)`

SetRowNumEgressAcl sets RowNumEgressAcl field to given value.

### HasRowNumEgressAcl

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) HasRowNumEgressAcl() bool`

HasRowNumEgressAcl returns a boolean if a field has been set.

### GetRowNumEgressAclRefType

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) GetRowNumEgressAclRefType() string`

GetRowNumEgressAclRefType returns the RowNumEgressAclRefType field if non-nil, zero value otherwise.

### GetRowNumEgressAclRefTypeOk

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) GetRowNumEgressAclRefTypeOk() (*string, bool)`

GetRowNumEgressAclRefTypeOk returns a tuple with the RowNumEgressAclRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowNumEgressAclRefType

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) SetRowNumEgressAclRefType(v string)`

SetRowNumEgressAclRefType sets RowNumEgressAclRefType field to given value.

### HasRowNumEgressAclRefType

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) HasRowNumEgressAclRefType() bool`

HasRowNumEgressAclRefType returns a boolean if a field has been set.

### GetIndex

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) GetIndex() int32`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) GetIndexOk() (*int32, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) SetIndex(v int32)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.

### GetRowNumMacFilter

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) GetRowNumMacFilter() string`

GetRowNumMacFilter returns the RowNumMacFilter field if non-nil, zero value otherwise.

### GetRowNumMacFilterOk

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) GetRowNumMacFilterOk() (*string, bool)`

GetRowNumMacFilterOk returns a tuple with the RowNumMacFilter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowNumMacFilter

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) SetRowNumMacFilter(v string)`

SetRowNumMacFilter sets RowNumMacFilter field to given value.

### HasRowNumMacFilter

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) HasRowNumMacFilter() bool`

HasRowNumMacFilter returns a boolean if a field has been set.

### GetRowNumMacFilterRefType

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) GetRowNumMacFilterRefType() string`

GetRowNumMacFilterRefType returns the RowNumMacFilterRefType field if non-nil, zero value otherwise.

### GetRowNumMacFilterRefTypeOk

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) GetRowNumMacFilterRefTypeOk() (*string, bool)`

GetRowNumMacFilterRefTypeOk returns a tuple with the RowNumMacFilterRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowNumMacFilterRefType

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) SetRowNumMacFilterRefType(v string)`

SetRowNumMacFilterRefType sets RowNumMacFilterRefType field to given value.

### HasRowNumMacFilterRefType

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) HasRowNumMacFilterRefType() bool`

HasRowNumMacFilterRefType returns a boolean if a field has been set.

### GetRowNumLanIptv

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) GetRowNumLanIptv() string`

GetRowNumLanIptv returns the RowNumLanIptv field if non-nil, zero value otherwise.

### GetRowNumLanIptvOk

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) GetRowNumLanIptvOk() (*string, bool)`

GetRowNumLanIptvOk returns a tuple with the RowNumLanIptv field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowNumLanIptv

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) SetRowNumLanIptv(v string)`

SetRowNumLanIptv sets RowNumLanIptv field to given value.

### HasRowNumLanIptv

`func (o *EthportprofilesPutRequestEthPortProfileValueServicesInner) HasRowNumLanIptv() bool`

HasRowNumLanIptv returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


