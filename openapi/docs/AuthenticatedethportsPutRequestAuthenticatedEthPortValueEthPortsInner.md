# AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**EthPortProfileNumEnable** | Pointer to **bool** | Enable row | [optional] [default to false]
**EthPortProfileNumEthPort** | Pointer to **string** | Choose an Eth Port Profile | [optional] [default to ""]
**EthPortProfileNumEthPortRefType** | Pointer to **string** | Object type for eth_port_profile_num_eth_port field | [optional] 
**EthPortProfileNumWalledGardenSet** | Pointer to **bool** | Flag indicating this Eth Port Profile is the Walled Garden | [optional] [default to false]
**EthPortProfileNumRadiusFilterId** | Pointer to **string** | The value of filter-id in the RADIUS response which will evoke this Eth Port Profile | [optional] [default to ""]
**Index** | Pointer to **int32** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 

## Methods

### NewAuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner

`func NewAuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner() *AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner`

NewAuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner instantiates a new AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInnerWithDefaults

`func NewAuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInnerWithDefaults() *AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner`

NewAuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInnerWithDefaults instantiates a new AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEthPortProfileNumEnable

`func (o *AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner) GetEthPortProfileNumEnable() bool`

GetEthPortProfileNumEnable returns the EthPortProfileNumEnable field if non-nil, zero value otherwise.

### GetEthPortProfileNumEnableOk

`func (o *AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner) GetEthPortProfileNumEnableOk() (*bool, bool)`

GetEthPortProfileNumEnableOk returns a tuple with the EthPortProfileNumEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEthPortProfileNumEnable

`func (o *AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner) SetEthPortProfileNumEnable(v bool)`

SetEthPortProfileNumEnable sets EthPortProfileNumEnable field to given value.

### HasEthPortProfileNumEnable

`func (o *AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner) HasEthPortProfileNumEnable() bool`

HasEthPortProfileNumEnable returns a boolean if a field has been set.

### GetEthPortProfileNumEthPort

`func (o *AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner) GetEthPortProfileNumEthPort() string`

GetEthPortProfileNumEthPort returns the EthPortProfileNumEthPort field if non-nil, zero value otherwise.

### GetEthPortProfileNumEthPortOk

`func (o *AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner) GetEthPortProfileNumEthPortOk() (*string, bool)`

GetEthPortProfileNumEthPortOk returns a tuple with the EthPortProfileNumEthPort field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEthPortProfileNumEthPort

`func (o *AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner) SetEthPortProfileNumEthPort(v string)`

SetEthPortProfileNumEthPort sets EthPortProfileNumEthPort field to given value.

### HasEthPortProfileNumEthPort

`func (o *AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner) HasEthPortProfileNumEthPort() bool`

HasEthPortProfileNumEthPort returns a boolean if a field has been set.

### GetEthPortProfileNumEthPortRefType

`func (o *AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner) GetEthPortProfileNumEthPortRefType() string`

GetEthPortProfileNumEthPortRefType returns the EthPortProfileNumEthPortRefType field if non-nil, zero value otherwise.

### GetEthPortProfileNumEthPortRefTypeOk

`func (o *AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner) GetEthPortProfileNumEthPortRefTypeOk() (*string, bool)`

GetEthPortProfileNumEthPortRefTypeOk returns a tuple with the EthPortProfileNumEthPortRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEthPortProfileNumEthPortRefType

`func (o *AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner) SetEthPortProfileNumEthPortRefType(v string)`

SetEthPortProfileNumEthPortRefType sets EthPortProfileNumEthPortRefType field to given value.

### HasEthPortProfileNumEthPortRefType

`func (o *AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner) HasEthPortProfileNumEthPortRefType() bool`

HasEthPortProfileNumEthPortRefType returns a boolean if a field has been set.

### GetEthPortProfileNumWalledGardenSet

`func (o *AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner) GetEthPortProfileNumWalledGardenSet() bool`

GetEthPortProfileNumWalledGardenSet returns the EthPortProfileNumWalledGardenSet field if non-nil, zero value otherwise.

### GetEthPortProfileNumWalledGardenSetOk

`func (o *AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner) GetEthPortProfileNumWalledGardenSetOk() (*bool, bool)`

GetEthPortProfileNumWalledGardenSetOk returns a tuple with the EthPortProfileNumWalledGardenSet field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEthPortProfileNumWalledGardenSet

`func (o *AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner) SetEthPortProfileNumWalledGardenSet(v bool)`

SetEthPortProfileNumWalledGardenSet sets EthPortProfileNumWalledGardenSet field to given value.

### HasEthPortProfileNumWalledGardenSet

`func (o *AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner) HasEthPortProfileNumWalledGardenSet() bool`

HasEthPortProfileNumWalledGardenSet returns a boolean if a field has been set.

### GetEthPortProfileNumRadiusFilterId

`func (o *AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner) GetEthPortProfileNumRadiusFilterId() string`

GetEthPortProfileNumRadiusFilterId returns the EthPortProfileNumRadiusFilterId field if non-nil, zero value otherwise.

### GetEthPortProfileNumRadiusFilterIdOk

`func (o *AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner) GetEthPortProfileNumRadiusFilterIdOk() (*string, bool)`

GetEthPortProfileNumRadiusFilterIdOk returns a tuple with the EthPortProfileNumRadiusFilterId field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEthPortProfileNumRadiusFilterId

`func (o *AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner) SetEthPortProfileNumRadiusFilterId(v string)`

SetEthPortProfileNumRadiusFilterId sets EthPortProfileNumRadiusFilterId field to given value.

### HasEthPortProfileNumRadiusFilterId

`func (o *AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner) HasEthPortProfileNumRadiusFilterId() bool`

HasEthPortProfileNumRadiusFilterId returns a boolean if a field has been set.

### GetIndex

`func (o *AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner) GetIndex() int32`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner) GetIndexOk() (*int32, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner) SetIndex(v int32)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *AuthenticatedethportsPutRequestAuthenticatedEthPortValueEthPortsInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


