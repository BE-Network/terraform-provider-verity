# EthportsettingsPutRequestEthPortSettingsValueLldpMedInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**LldpMedRowNumEnable** | Pointer to **bool** | Per LLDP Med row enable | [optional] [default to false]
**LldpMedRowNumAdvertisedApplicatio** | Pointer to **string** | Advertised application | [optional] [default to ""]
**LldpMedRowNumDscpMark** | Pointer to **NullableInt32** | LLDP DSCP Mark | [optional] [default to 0]
**LldpMedRowNumPriority** | Pointer to **NullableInt32** | LLDP Priority | [optional] [default to 0]
**LldpMedRowNumService** | Pointer to **string** | LLDP Service | [optional] [default to ""]
**LldpMedRowNumServiceRefType** | Pointer to **string** | Object type for lldp_med_row_num_service field | [optional] 
**Index** | Pointer to **int32** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 

## Methods

### NewEthportsettingsPutRequestEthPortSettingsValueLldpMedInner

`func NewEthportsettingsPutRequestEthPortSettingsValueLldpMedInner() *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner`

NewEthportsettingsPutRequestEthPortSettingsValueLldpMedInner instantiates a new EthportsettingsPutRequestEthPortSettingsValueLldpMedInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewEthportsettingsPutRequestEthPortSettingsValueLldpMedInnerWithDefaults

`func NewEthportsettingsPutRequestEthPortSettingsValueLldpMedInnerWithDefaults() *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner`

NewEthportsettingsPutRequestEthPortSettingsValueLldpMedInnerWithDefaults instantiates a new EthportsettingsPutRequestEthPortSettingsValueLldpMedInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLldpMedRowNumEnable

`func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) GetLldpMedRowNumEnable() bool`

GetLldpMedRowNumEnable returns the LldpMedRowNumEnable field if non-nil, zero value otherwise.

### GetLldpMedRowNumEnableOk

`func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) GetLldpMedRowNumEnableOk() (*bool, bool)`

GetLldpMedRowNumEnableOk returns a tuple with the LldpMedRowNumEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLldpMedRowNumEnable

`func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) SetLldpMedRowNumEnable(v bool)`

SetLldpMedRowNumEnable sets LldpMedRowNumEnable field to given value.

### HasLldpMedRowNumEnable

`func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) HasLldpMedRowNumEnable() bool`

HasLldpMedRowNumEnable returns a boolean if a field has been set.

### GetLldpMedRowNumAdvertisedApplicatio

`func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) GetLldpMedRowNumAdvertisedApplicatio() string`

GetLldpMedRowNumAdvertisedApplicatio returns the LldpMedRowNumAdvertisedApplicatio field if non-nil, zero value otherwise.

### GetLldpMedRowNumAdvertisedApplicatioOk

`func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) GetLldpMedRowNumAdvertisedApplicatioOk() (*string, bool)`

GetLldpMedRowNumAdvertisedApplicatioOk returns a tuple with the LldpMedRowNumAdvertisedApplicatio field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLldpMedRowNumAdvertisedApplicatio

`func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) SetLldpMedRowNumAdvertisedApplicatio(v string)`

SetLldpMedRowNumAdvertisedApplicatio sets LldpMedRowNumAdvertisedApplicatio field to given value.

### HasLldpMedRowNumAdvertisedApplicatio

`func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) HasLldpMedRowNumAdvertisedApplicatio() bool`

HasLldpMedRowNumAdvertisedApplicatio returns a boolean if a field has been set.

### GetLldpMedRowNumDscpMark

`func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) GetLldpMedRowNumDscpMark() int32`

GetLldpMedRowNumDscpMark returns the LldpMedRowNumDscpMark field if non-nil, zero value otherwise.

### GetLldpMedRowNumDscpMarkOk

`func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) GetLldpMedRowNumDscpMarkOk() (*int32, bool)`

GetLldpMedRowNumDscpMarkOk returns a tuple with the LldpMedRowNumDscpMark field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLldpMedRowNumDscpMark

`func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) SetLldpMedRowNumDscpMark(v int32)`

SetLldpMedRowNumDscpMark sets LldpMedRowNumDscpMark field to given value.

### HasLldpMedRowNumDscpMark

`func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) HasLldpMedRowNumDscpMark() bool`

HasLldpMedRowNumDscpMark returns a boolean if a field has been set.

### SetLldpMedRowNumDscpMarkNil

`func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) SetLldpMedRowNumDscpMarkNil(b bool)`

 SetLldpMedRowNumDscpMarkNil sets the value for LldpMedRowNumDscpMark to be an explicit nil

### UnsetLldpMedRowNumDscpMark
`func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) UnsetLldpMedRowNumDscpMark()`

UnsetLldpMedRowNumDscpMark ensures that no value is present for LldpMedRowNumDscpMark, not even an explicit nil
### GetLldpMedRowNumPriority

`func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) GetLldpMedRowNumPriority() int32`

GetLldpMedRowNumPriority returns the LldpMedRowNumPriority field if non-nil, zero value otherwise.

### GetLldpMedRowNumPriorityOk

`func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) GetLldpMedRowNumPriorityOk() (*int32, bool)`

GetLldpMedRowNumPriorityOk returns a tuple with the LldpMedRowNumPriority field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLldpMedRowNumPriority

`func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) SetLldpMedRowNumPriority(v int32)`

SetLldpMedRowNumPriority sets LldpMedRowNumPriority field to given value.

### HasLldpMedRowNumPriority

`func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) HasLldpMedRowNumPriority() bool`

HasLldpMedRowNumPriority returns a boolean if a field has been set.

### SetLldpMedRowNumPriorityNil

`func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) SetLldpMedRowNumPriorityNil(b bool)`

 SetLldpMedRowNumPriorityNil sets the value for LldpMedRowNumPriority to be an explicit nil

### UnsetLldpMedRowNumPriority
`func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) UnsetLldpMedRowNumPriority()`

UnsetLldpMedRowNumPriority ensures that no value is present for LldpMedRowNumPriority, not even an explicit nil
### GetLldpMedRowNumService

`func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) GetLldpMedRowNumService() string`

GetLldpMedRowNumService returns the LldpMedRowNumService field if non-nil, zero value otherwise.

### GetLldpMedRowNumServiceOk

`func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) GetLldpMedRowNumServiceOk() (*string, bool)`

GetLldpMedRowNumServiceOk returns a tuple with the LldpMedRowNumService field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLldpMedRowNumService

`func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) SetLldpMedRowNumService(v string)`

SetLldpMedRowNumService sets LldpMedRowNumService field to given value.

### HasLldpMedRowNumService

`func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) HasLldpMedRowNumService() bool`

HasLldpMedRowNumService returns a boolean if a field has been set.

### GetLldpMedRowNumServiceRefType

`func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) GetLldpMedRowNumServiceRefType() string`

GetLldpMedRowNumServiceRefType returns the LldpMedRowNumServiceRefType field if non-nil, zero value otherwise.

### GetLldpMedRowNumServiceRefTypeOk

`func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) GetLldpMedRowNumServiceRefTypeOk() (*string, bool)`

GetLldpMedRowNumServiceRefTypeOk returns a tuple with the LldpMedRowNumServiceRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLldpMedRowNumServiceRefType

`func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) SetLldpMedRowNumServiceRefType(v string)`

SetLldpMedRowNumServiceRefType sets LldpMedRowNumServiceRefType field to given value.

### HasLldpMedRowNumServiceRefType

`func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) HasLldpMedRowNumServiceRefType() bool`

HasLldpMedRowNumServiceRefType returns a boolean if a field has been set.

### GetIndex

`func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) GetIndex() int32`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) GetIndexOk() (*int32, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) SetIndex(v int32)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *EthportsettingsPutRequestEthPortSettingsValueLldpMedInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


