# LagsPutRequestLagValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. It&#39;s highly recommended to set this value to true so that validation on the object will be ran. | [optional] [default to false]
**IsPeerLink** | Pointer to **bool** | Indicates this LAG is used for peer-to-peer Peer-LAG/IDS link | [optional] [default to false]
**Color** | Pointer to **string** | Choose the color to display the connectors on the network view | [optional] [default to "anakiwa"]
**Lacp** | Pointer to **bool** | LACP | [optional] [default to true]
**EthPortProfile** | Pointer to **string** | Choose an Eth Port Profile | [optional] [default to ""]
**EthPortProfileRefType** | Pointer to **string** | Object type for eth_port_profile field | [optional] 
**PeerLinkVlan** | Pointer to **NullableInt32** | For peer-peer LAGs. The VLAN used for control | [optional] 
**Fallback** | Pointer to **bool** | Allows an active member interface to establish a connection with a peer interface before the port channel receives the LACP protocol negotiation from the peer. | [optional] [default to false]
**FastRate** | Pointer to **bool** | Send LACP packets every second (if disabled, packets are sent every 30 seconds) | [optional] [default to false]
**ObjectProperties** | Pointer to **map[string]interface{}** |  | [optional] 
**Uplink** | Pointer to **bool** | Indicates this LAG is designated as an uplink in the case of a spineless pod. Link State Tracking will be applied to BGP Egress VLANs/Interfaces and the MCLAG Peer Link VLAN | [optional] [default to false]

## Methods

### NewLagsPutRequestLagValue

`func NewLagsPutRequestLagValue() *LagsPutRequestLagValue`

NewLagsPutRequestLagValue instantiates a new LagsPutRequestLagValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewLagsPutRequestLagValueWithDefaults

`func NewLagsPutRequestLagValueWithDefaults() *LagsPutRequestLagValue`

NewLagsPutRequestLagValueWithDefaults instantiates a new LagsPutRequestLagValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *LagsPutRequestLagValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *LagsPutRequestLagValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *LagsPutRequestLagValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *LagsPutRequestLagValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *LagsPutRequestLagValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *LagsPutRequestLagValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *LagsPutRequestLagValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *LagsPutRequestLagValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetIsPeerLink

`func (o *LagsPutRequestLagValue) GetIsPeerLink() bool`

GetIsPeerLink returns the IsPeerLink field if non-nil, zero value otherwise.

### GetIsPeerLinkOk

`func (o *LagsPutRequestLagValue) GetIsPeerLinkOk() (*bool, bool)`

GetIsPeerLinkOk returns a tuple with the IsPeerLink field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsPeerLink

`func (o *LagsPutRequestLagValue) SetIsPeerLink(v bool)`

SetIsPeerLink sets IsPeerLink field to given value.

### HasIsPeerLink

`func (o *LagsPutRequestLagValue) HasIsPeerLink() bool`

HasIsPeerLink returns a boolean if a field has been set.

### GetColor

`func (o *LagsPutRequestLagValue) GetColor() string`

GetColor returns the Color field if non-nil, zero value otherwise.

### GetColorOk

`func (o *LagsPutRequestLagValue) GetColorOk() (*string, bool)`

GetColorOk returns a tuple with the Color field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetColor

`func (o *LagsPutRequestLagValue) SetColor(v string)`

SetColor sets Color field to given value.

### HasColor

`func (o *LagsPutRequestLagValue) HasColor() bool`

HasColor returns a boolean if a field has been set.

### GetLacp

`func (o *LagsPutRequestLagValue) GetLacp() bool`

GetLacp returns the Lacp field if non-nil, zero value otherwise.

### GetLacpOk

`func (o *LagsPutRequestLagValue) GetLacpOk() (*bool, bool)`

GetLacpOk returns a tuple with the Lacp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLacp

`func (o *LagsPutRequestLagValue) SetLacp(v bool)`

SetLacp sets Lacp field to given value.

### HasLacp

`func (o *LagsPutRequestLagValue) HasLacp() bool`

HasLacp returns a boolean if a field has been set.

### GetEthPortProfile

`func (o *LagsPutRequestLagValue) GetEthPortProfile() string`

GetEthPortProfile returns the EthPortProfile field if non-nil, zero value otherwise.

### GetEthPortProfileOk

`func (o *LagsPutRequestLagValue) GetEthPortProfileOk() (*string, bool)`

GetEthPortProfileOk returns a tuple with the EthPortProfile field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEthPortProfile

`func (o *LagsPutRequestLagValue) SetEthPortProfile(v string)`

SetEthPortProfile sets EthPortProfile field to given value.

### HasEthPortProfile

`func (o *LagsPutRequestLagValue) HasEthPortProfile() bool`

HasEthPortProfile returns a boolean if a field has been set.

### GetEthPortProfileRefType

`func (o *LagsPutRequestLagValue) GetEthPortProfileRefType() string`

GetEthPortProfileRefType returns the EthPortProfileRefType field if non-nil, zero value otherwise.

### GetEthPortProfileRefTypeOk

`func (o *LagsPutRequestLagValue) GetEthPortProfileRefTypeOk() (*string, bool)`

GetEthPortProfileRefTypeOk returns a tuple with the EthPortProfileRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEthPortProfileRefType

`func (o *LagsPutRequestLagValue) SetEthPortProfileRefType(v string)`

SetEthPortProfileRefType sets EthPortProfileRefType field to given value.

### HasEthPortProfileRefType

`func (o *LagsPutRequestLagValue) HasEthPortProfileRefType() bool`

HasEthPortProfileRefType returns a boolean if a field has been set.

### GetPeerLinkVlan

`func (o *LagsPutRequestLagValue) GetPeerLinkVlan() int32`

GetPeerLinkVlan returns the PeerLinkVlan field if non-nil, zero value otherwise.

### GetPeerLinkVlanOk

`func (o *LagsPutRequestLagValue) GetPeerLinkVlanOk() (*int32, bool)`

GetPeerLinkVlanOk returns a tuple with the PeerLinkVlan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPeerLinkVlan

`func (o *LagsPutRequestLagValue) SetPeerLinkVlan(v int32)`

SetPeerLinkVlan sets PeerLinkVlan field to given value.

### HasPeerLinkVlan

`func (o *LagsPutRequestLagValue) HasPeerLinkVlan() bool`

HasPeerLinkVlan returns a boolean if a field has been set.

### SetPeerLinkVlanNil

`func (o *LagsPutRequestLagValue) SetPeerLinkVlanNil(b bool)`

 SetPeerLinkVlanNil sets the value for PeerLinkVlan to be an explicit nil

### UnsetPeerLinkVlan
`func (o *LagsPutRequestLagValue) UnsetPeerLinkVlan()`

UnsetPeerLinkVlan ensures that no value is present for PeerLinkVlan, not even an explicit nil
### GetFallback

`func (o *LagsPutRequestLagValue) GetFallback() bool`

GetFallback returns the Fallback field if non-nil, zero value otherwise.

### GetFallbackOk

`func (o *LagsPutRequestLagValue) GetFallbackOk() (*bool, bool)`

GetFallbackOk returns a tuple with the Fallback field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFallback

`func (o *LagsPutRequestLagValue) SetFallback(v bool)`

SetFallback sets Fallback field to given value.

### HasFallback

`func (o *LagsPutRequestLagValue) HasFallback() bool`

HasFallback returns a boolean if a field has been set.

### GetFastRate

`func (o *LagsPutRequestLagValue) GetFastRate() bool`

GetFastRate returns the FastRate field if non-nil, zero value otherwise.

### GetFastRateOk

`func (o *LagsPutRequestLagValue) GetFastRateOk() (*bool, bool)`

GetFastRateOk returns a tuple with the FastRate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFastRate

`func (o *LagsPutRequestLagValue) SetFastRate(v bool)`

SetFastRate sets FastRate field to given value.

### HasFastRate

`func (o *LagsPutRequestLagValue) HasFastRate() bool`

HasFastRate returns a boolean if a field has been set.

### GetObjectProperties

`func (o *LagsPutRequestLagValue) GetObjectProperties() map[string]interface{}`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *LagsPutRequestLagValue) GetObjectPropertiesOk() (*map[string]interface{}, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *LagsPutRequestLagValue) SetObjectProperties(v map[string]interface{})`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *LagsPutRequestLagValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.

### GetUplink

`func (o *LagsPutRequestLagValue) GetUplink() bool`

GetUplink returns the Uplink field if non-nil, zero value otherwise.

### GetUplinkOk

`func (o *LagsPutRequestLagValue) GetUplinkOk() (*bool, bool)`

GetUplinkOk returns a tuple with the Uplink field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUplink

`func (o *LagsPutRequestLagValue) SetUplink(v bool)`

SetUplink sets Uplink field to given value.

### HasUplink

`func (o *LagsPutRequestLagValue) HasUplink() bool`

HasUplink returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


