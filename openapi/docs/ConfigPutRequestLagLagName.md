# ConfigPutRequestLagLagName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to false]
**IsPeerLink** | Pointer to **bool** | Indicates this LAG is used for peer-to-peer Peer-LAG/IDS link | [optional] [default to false]
**Color** | Pointer to **string** | Choose the color to display the connectors on the network view | [optional] [default to "anakiwa"]
**Lacp** | Pointer to **bool** | LACP | [optional] [default to true]
**EthPortProfile** | Pointer to **string** | Choose an Eth Port Profile | [optional] [default to ""]
**EthPortProfileRefType** | Pointer to **string** | Object type for eth_port_profile field | [optional] 
**PeerLinkVlan** | Pointer to **NullableInt32** | For peer-peer LAGs. The VLAN used for control | [optional] 
**Fallback** | Pointer to **bool** | Allows an active member interface to establish a connection with a peer interface before the port channel receives the LACP protocol negotiation from the peer. | [optional] [default to false]
**FastRate** | Pointer to **bool** | Send LACP packets every second (if disabled, packets are sent every 30 seconds) | [optional] [default to false]
**ObjectProperties** | Pointer to **map[string]interface{}** |  | [optional] 

## Methods

### NewConfigPutRequestLagLagName

`func NewConfigPutRequestLagLagName() *ConfigPutRequestLagLagName`

NewConfigPutRequestLagLagName instantiates a new ConfigPutRequestLagLagName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestLagLagNameWithDefaults

`func NewConfigPutRequestLagLagNameWithDefaults() *ConfigPutRequestLagLagName`

NewConfigPutRequestLagLagNameWithDefaults instantiates a new ConfigPutRequestLagLagName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestLagLagName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestLagLagName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestLagLagName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestLagLagName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ConfigPutRequestLagLagName) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestLagLagName) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestLagLagName) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestLagLagName) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetIsPeerLink

`func (o *ConfigPutRequestLagLagName) GetIsPeerLink() bool`

GetIsPeerLink returns the IsPeerLink field if non-nil, zero value otherwise.

### GetIsPeerLinkOk

`func (o *ConfigPutRequestLagLagName) GetIsPeerLinkOk() (*bool, bool)`

GetIsPeerLinkOk returns a tuple with the IsPeerLink field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsPeerLink

`func (o *ConfigPutRequestLagLagName) SetIsPeerLink(v bool)`

SetIsPeerLink sets IsPeerLink field to given value.

### HasIsPeerLink

`func (o *ConfigPutRequestLagLagName) HasIsPeerLink() bool`

HasIsPeerLink returns a boolean if a field has been set.

### GetColor

`func (o *ConfigPutRequestLagLagName) GetColor() string`

GetColor returns the Color field if non-nil, zero value otherwise.

### GetColorOk

`func (o *ConfigPutRequestLagLagName) GetColorOk() (*string, bool)`

GetColorOk returns a tuple with the Color field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetColor

`func (o *ConfigPutRequestLagLagName) SetColor(v string)`

SetColor sets Color field to given value.

### HasColor

`func (o *ConfigPutRequestLagLagName) HasColor() bool`

HasColor returns a boolean if a field has been set.

### GetLacp

`func (o *ConfigPutRequestLagLagName) GetLacp() bool`

GetLacp returns the Lacp field if non-nil, zero value otherwise.

### GetLacpOk

`func (o *ConfigPutRequestLagLagName) GetLacpOk() (*bool, bool)`

GetLacpOk returns a tuple with the Lacp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLacp

`func (o *ConfigPutRequestLagLagName) SetLacp(v bool)`

SetLacp sets Lacp field to given value.

### HasLacp

`func (o *ConfigPutRequestLagLagName) HasLacp() bool`

HasLacp returns a boolean if a field has been set.

### GetEthPortProfile

`func (o *ConfigPutRequestLagLagName) GetEthPortProfile() string`

GetEthPortProfile returns the EthPortProfile field if non-nil, zero value otherwise.

### GetEthPortProfileOk

`func (o *ConfigPutRequestLagLagName) GetEthPortProfileOk() (*string, bool)`

GetEthPortProfileOk returns a tuple with the EthPortProfile field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEthPortProfile

`func (o *ConfigPutRequestLagLagName) SetEthPortProfile(v string)`

SetEthPortProfile sets EthPortProfile field to given value.

### HasEthPortProfile

`func (o *ConfigPutRequestLagLagName) HasEthPortProfile() bool`

HasEthPortProfile returns a boolean if a field has been set.

### GetEthPortProfileRefType

`func (o *ConfigPutRequestLagLagName) GetEthPortProfileRefType() string`

GetEthPortProfileRefType returns the EthPortProfileRefType field if non-nil, zero value otherwise.

### GetEthPortProfileRefTypeOk

`func (o *ConfigPutRequestLagLagName) GetEthPortProfileRefTypeOk() (*string, bool)`

GetEthPortProfileRefTypeOk returns a tuple with the EthPortProfileRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEthPortProfileRefType

`func (o *ConfigPutRequestLagLagName) SetEthPortProfileRefType(v string)`

SetEthPortProfileRefType sets EthPortProfileRefType field to given value.

### HasEthPortProfileRefType

`func (o *ConfigPutRequestLagLagName) HasEthPortProfileRefType() bool`

HasEthPortProfileRefType returns a boolean if a field has been set.

### GetPeerLinkVlan

`func (o *ConfigPutRequestLagLagName) GetPeerLinkVlan() int32`

GetPeerLinkVlan returns the PeerLinkVlan field if non-nil, zero value otherwise.

### GetPeerLinkVlanOk

`func (o *ConfigPutRequestLagLagName) GetPeerLinkVlanOk() (*int32, bool)`

GetPeerLinkVlanOk returns a tuple with the PeerLinkVlan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPeerLinkVlan

`func (o *ConfigPutRequestLagLagName) SetPeerLinkVlan(v int32)`

SetPeerLinkVlan sets PeerLinkVlan field to given value.

### HasPeerLinkVlan

`func (o *ConfigPutRequestLagLagName) HasPeerLinkVlan() bool`

HasPeerLinkVlan returns a boolean if a field has been set.

### SetPeerLinkVlanNil

`func (o *ConfigPutRequestLagLagName) SetPeerLinkVlanNil(b bool)`

 SetPeerLinkVlanNil sets the value for PeerLinkVlan to be an explicit nil

### UnsetPeerLinkVlan
`func (o *ConfigPutRequestLagLagName) UnsetPeerLinkVlan()`

UnsetPeerLinkVlan ensures that no value is present for PeerLinkVlan, not even an explicit nil
### GetFallback

`func (o *ConfigPutRequestLagLagName) GetFallback() bool`

GetFallback returns the Fallback field if non-nil, zero value otherwise.

### GetFallbackOk

`func (o *ConfigPutRequestLagLagName) GetFallbackOk() (*bool, bool)`

GetFallbackOk returns a tuple with the Fallback field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFallback

`func (o *ConfigPutRequestLagLagName) SetFallback(v bool)`

SetFallback sets Fallback field to given value.

### HasFallback

`func (o *ConfigPutRequestLagLagName) HasFallback() bool`

HasFallback returns a boolean if a field has been set.

### GetFastRate

`func (o *ConfigPutRequestLagLagName) GetFastRate() bool`

GetFastRate returns the FastRate field if non-nil, zero value otherwise.

### GetFastRateOk

`func (o *ConfigPutRequestLagLagName) GetFastRateOk() (*bool, bool)`

GetFastRateOk returns a tuple with the FastRate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFastRate

`func (o *ConfigPutRequestLagLagName) SetFastRate(v bool)`

SetFastRate sets FastRate field to given value.

### HasFastRate

`func (o *ConfigPutRequestLagLagName) HasFastRate() bool`

HasFastRate returns a boolean if a field has been set.

### GetObjectProperties

`func (o *ConfigPutRequestLagLagName) GetObjectProperties() map[string]interface{}`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ConfigPutRequestLagLagName) GetObjectPropertiesOk() (*map[string]interface{}, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ConfigPutRequestLagLagName) SetObjectProperties(v map[string]interface{})`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ConfigPutRequestLagLagName) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


