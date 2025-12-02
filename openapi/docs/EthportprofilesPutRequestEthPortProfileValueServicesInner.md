# EthportprofilesPutRequestEthPortProfileValueServicesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**RowNumEnable** | Pointer to **bool** | Enable row | [optional] [default to false]
**RowNumService** | Pointer to **string** | Choose a Service to connect | [optional] [default to ""]
**RowNumServiceRefType** | Pointer to **string** | Object type for row_num_service field | [optional] 
**RowNumExternalVlan** | Pointer to **NullableInt32** | Choose an external vlan A value of 0 will make the VLAN untagged, while in case null is provided, the VLAN will be the one associated with the service. | [optional] 
**Index** | Pointer to **int32** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 

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


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


