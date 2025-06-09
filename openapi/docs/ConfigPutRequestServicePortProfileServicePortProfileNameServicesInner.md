# ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**RowNumEnable** | Pointer to **bool** | Enable row | [optional] [default to false]
**RowNumService** | Pointer to **string** | Connect a Service | [optional] [default to ""]
**RowNumServiceRefType** | Pointer to **string** | Object type for row_num_service field | [optional] 
**RowNumExternalVlan** | Pointer to **NullableInt32** | Choose an external vlan | [optional] 
**RowNumLimitIn** | Pointer to **NullableInt32** | Speed of ingress (Mbps) | [optional] 
**RowNumLimitOut** | Pointer to **NullableInt32** | Speed of egress (Mbps) | [optional] [default to 1000]
**Index** | Pointer to **int32** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 

## Methods

### NewConfigPutRequestServicePortProfileServicePortProfileNameServicesInner

`func NewConfigPutRequestServicePortProfileServicePortProfileNameServicesInner() *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner`

NewConfigPutRequestServicePortProfileServicePortProfileNameServicesInner instantiates a new ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestServicePortProfileServicePortProfileNameServicesInnerWithDefaults

`func NewConfigPutRequestServicePortProfileServicePortProfileNameServicesInnerWithDefaults() *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner`

NewConfigPutRequestServicePortProfileServicePortProfileNameServicesInnerWithDefaults instantiates a new ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRowNumEnable

`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) GetRowNumEnable() bool`

GetRowNumEnable returns the RowNumEnable field if non-nil, zero value otherwise.

### GetRowNumEnableOk

`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) GetRowNumEnableOk() (*bool, bool)`

GetRowNumEnableOk returns a tuple with the RowNumEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowNumEnable

`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) SetRowNumEnable(v bool)`

SetRowNumEnable sets RowNumEnable field to given value.

### HasRowNumEnable

`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) HasRowNumEnable() bool`

HasRowNumEnable returns a boolean if a field has been set.

### GetRowNumService

`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) GetRowNumService() string`

GetRowNumService returns the RowNumService field if non-nil, zero value otherwise.

### GetRowNumServiceOk

`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) GetRowNumServiceOk() (*string, bool)`

GetRowNumServiceOk returns a tuple with the RowNumService field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowNumService

`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) SetRowNumService(v string)`

SetRowNumService sets RowNumService field to given value.

### HasRowNumService

`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) HasRowNumService() bool`

HasRowNumService returns a boolean if a field has been set.

### GetRowNumServiceRefType

`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) GetRowNumServiceRefType() string`

GetRowNumServiceRefType returns the RowNumServiceRefType field if non-nil, zero value otherwise.

### GetRowNumServiceRefTypeOk

`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) GetRowNumServiceRefTypeOk() (*string, bool)`

GetRowNumServiceRefTypeOk returns a tuple with the RowNumServiceRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowNumServiceRefType

`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) SetRowNumServiceRefType(v string)`

SetRowNumServiceRefType sets RowNumServiceRefType field to given value.

### HasRowNumServiceRefType

`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) HasRowNumServiceRefType() bool`

HasRowNumServiceRefType returns a boolean if a field has been set.

### GetRowNumExternalVlan

`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) GetRowNumExternalVlan() int32`

GetRowNumExternalVlan returns the RowNumExternalVlan field if non-nil, zero value otherwise.

### GetRowNumExternalVlanOk

`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) GetRowNumExternalVlanOk() (*int32, bool)`

GetRowNumExternalVlanOk returns a tuple with the RowNumExternalVlan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowNumExternalVlan

`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) SetRowNumExternalVlan(v int32)`

SetRowNumExternalVlan sets RowNumExternalVlan field to given value.

### HasRowNumExternalVlan

`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) HasRowNumExternalVlan() bool`

HasRowNumExternalVlan returns a boolean if a field has been set.

### SetRowNumExternalVlanNil

`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) SetRowNumExternalVlanNil(b bool)`

 SetRowNumExternalVlanNil sets the value for RowNumExternalVlan to be an explicit nil

### UnsetRowNumExternalVlan
`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) UnsetRowNumExternalVlan()`

UnsetRowNumExternalVlan ensures that no value is present for RowNumExternalVlan, not even an explicit nil
### GetRowNumLimitIn

`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) GetRowNumLimitIn() int32`

GetRowNumLimitIn returns the RowNumLimitIn field if non-nil, zero value otherwise.

### GetRowNumLimitInOk

`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) GetRowNumLimitInOk() (*int32, bool)`

GetRowNumLimitInOk returns a tuple with the RowNumLimitIn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowNumLimitIn

`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) SetRowNumLimitIn(v int32)`

SetRowNumLimitIn sets RowNumLimitIn field to given value.

### HasRowNumLimitIn

`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) HasRowNumLimitIn() bool`

HasRowNumLimitIn returns a boolean if a field has been set.

### SetRowNumLimitInNil

`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) SetRowNumLimitInNil(b bool)`

 SetRowNumLimitInNil sets the value for RowNumLimitIn to be an explicit nil

### UnsetRowNumLimitIn
`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) UnsetRowNumLimitIn()`

UnsetRowNumLimitIn ensures that no value is present for RowNumLimitIn, not even an explicit nil
### GetRowNumLimitOut

`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) GetRowNumLimitOut() int32`

GetRowNumLimitOut returns the RowNumLimitOut field if non-nil, zero value otherwise.

### GetRowNumLimitOutOk

`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) GetRowNumLimitOutOk() (*int32, bool)`

GetRowNumLimitOutOk returns a tuple with the RowNumLimitOut field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowNumLimitOut

`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) SetRowNumLimitOut(v int32)`

SetRowNumLimitOut sets RowNumLimitOut field to given value.

### HasRowNumLimitOut

`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) HasRowNumLimitOut() bool`

HasRowNumLimitOut returns a boolean if a field has been set.

### SetRowNumLimitOutNil

`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) SetRowNumLimitOutNil(b bool)`

 SetRowNumLimitOutNil sets the value for RowNumLimitOut to be an explicit nil

### UnsetRowNumLimitOut
`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) UnsetRowNumLimitOut()`

UnsetRowNumLimitOut ensures that no value is present for RowNumLimitOut, not even an explicit nil
### GetIndex

`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) GetIndex() int32`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) GetIndexOk() (*int32, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) SetIndex(v int32)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *ConfigPutRequestServicePortProfileServicePortProfileNameServicesInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


