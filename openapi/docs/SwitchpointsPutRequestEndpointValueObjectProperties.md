# SwitchpointsPutRequestEndpointValueObjectProperties

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**UserNotes** | Pointer to **string** | Notes writen by User about the site | [optional] [default to ""]
**ExpectedParentEndpoint** | Pointer to **string** | Expected Parent Endpoint | [optional] [default to ""]
**ExpectedParentEndpointRefType** | Pointer to **string** | Object type for expected_parent_endpoint field | [optional] 
**NumberOfMultipoints** | Pointer to **NullableInt32** | Number of Multipoints | [optional] [default to 0]
**DrawAsEdgeDevice** | Pointer to **bool** | Turn on to display the switch as an edge device instead of as a switch | [optional] [default to false]
**Aggregate** | Pointer to **bool** | For Switch Endpoints. Denotes switch aggregated with all of its sub switches | [optional] [default to false]
**IsHost** | Pointer to **bool** | For Switch Endpoints. Denotes the Host Switch | [optional] [default to false]
**Eths** | Pointer to [**[]SwitchpointsPutRequestSwitchpointValueObjectPropertiesEthsInner**](SwitchpointsPutRequestSwitchpointValueObjectPropertiesEthsInner.md) |  | [optional] 

## Methods

### NewSwitchpointsPutRequestEndpointValueObjectProperties

`func NewSwitchpointsPutRequestEndpointValueObjectProperties() *SwitchpointsPutRequestEndpointValueObjectProperties`

NewSwitchpointsPutRequestEndpointValueObjectProperties instantiates a new SwitchpointsPutRequestEndpointValueObjectProperties object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSwitchpointsPutRequestEndpointValueObjectPropertiesWithDefaults

`func NewSwitchpointsPutRequestEndpointValueObjectPropertiesWithDefaults() *SwitchpointsPutRequestEndpointValueObjectProperties`

NewSwitchpointsPutRequestEndpointValueObjectPropertiesWithDefaults instantiates a new SwitchpointsPutRequestEndpointValueObjectProperties object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUserNotes

`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) GetUserNotes() string`

GetUserNotes returns the UserNotes field if non-nil, zero value otherwise.

### GetUserNotesOk

`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) GetUserNotesOk() (*string, bool)`

GetUserNotesOk returns a tuple with the UserNotes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserNotes

`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) SetUserNotes(v string)`

SetUserNotes sets UserNotes field to given value.

### HasUserNotes

`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) HasUserNotes() bool`

HasUserNotes returns a boolean if a field has been set.

### GetExpectedParentEndpoint

`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) GetExpectedParentEndpoint() string`

GetExpectedParentEndpoint returns the ExpectedParentEndpoint field if non-nil, zero value otherwise.

### GetExpectedParentEndpointOk

`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) GetExpectedParentEndpointOk() (*string, bool)`

GetExpectedParentEndpointOk returns a tuple with the ExpectedParentEndpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpectedParentEndpoint

`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) SetExpectedParentEndpoint(v string)`

SetExpectedParentEndpoint sets ExpectedParentEndpoint field to given value.

### HasExpectedParentEndpoint

`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) HasExpectedParentEndpoint() bool`

HasExpectedParentEndpoint returns a boolean if a field has been set.

### GetExpectedParentEndpointRefType

`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) GetExpectedParentEndpointRefType() string`

GetExpectedParentEndpointRefType returns the ExpectedParentEndpointRefType field if non-nil, zero value otherwise.

### GetExpectedParentEndpointRefTypeOk

`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) GetExpectedParentEndpointRefTypeOk() (*string, bool)`

GetExpectedParentEndpointRefTypeOk returns a tuple with the ExpectedParentEndpointRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpectedParentEndpointRefType

`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) SetExpectedParentEndpointRefType(v string)`

SetExpectedParentEndpointRefType sets ExpectedParentEndpointRefType field to given value.

### HasExpectedParentEndpointRefType

`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) HasExpectedParentEndpointRefType() bool`

HasExpectedParentEndpointRefType returns a boolean if a field has been set.

### GetNumberOfMultipoints

`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) GetNumberOfMultipoints() int32`

GetNumberOfMultipoints returns the NumberOfMultipoints field if non-nil, zero value otherwise.

### GetNumberOfMultipointsOk

`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) GetNumberOfMultipointsOk() (*int32, bool)`

GetNumberOfMultipointsOk returns a tuple with the NumberOfMultipoints field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNumberOfMultipoints

`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) SetNumberOfMultipoints(v int32)`

SetNumberOfMultipoints sets NumberOfMultipoints field to given value.

### HasNumberOfMultipoints

`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) HasNumberOfMultipoints() bool`

HasNumberOfMultipoints returns a boolean if a field has been set.

### SetNumberOfMultipointsNil

`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) SetNumberOfMultipointsNil(b bool)`

 SetNumberOfMultipointsNil sets the value for NumberOfMultipoints to be an explicit nil

### UnsetNumberOfMultipoints
`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) UnsetNumberOfMultipoints()`

UnsetNumberOfMultipoints ensures that no value is present for NumberOfMultipoints, not even an explicit nil
### GetDrawAsEdgeDevice

`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) GetDrawAsEdgeDevice() bool`

GetDrawAsEdgeDevice returns the DrawAsEdgeDevice field if non-nil, zero value otherwise.

### GetDrawAsEdgeDeviceOk

`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) GetDrawAsEdgeDeviceOk() (*bool, bool)`

GetDrawAsEdgeDeviceOk returns a tuple with the DrawAsEdgeDevice field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDrawAsEdgeDevice

`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) SetDrawAsEdgeDevice(v bool)`

SetDrawAsEdgeDevice sets DrawAsEdgeDevice field to given value.

### HasDrawAsEdgeDevice

`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) HasDrawAsEdgeDevice() bool`

HasDrawAsEdgeDevice returns a boolean if a field has been set.

### GetAggregate

`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) GetAggregate() bool`

GetAggregate returns the Aggregate field if non-nil, zero value otherwise.

### GetAggregateOk

`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) GetAggregateOk() (*bool, bool)`

GetAggregateOk returns a tuple with the Aggregate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAggregate

`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) SetAggregate(v bool)`

SetAggregate sets Aggregate field to given value.

### HasAggregate

`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) HasAggregate() bool`

HasAggregate returns a boolean if a field has been set.

### GetIsHost

`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) GetIsHost() bool`

GetIsHost returns the IsHost field if non-nil, zero value otherwise.

### GetIsHostOk

`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) GetIsHostOk() (*bool, bool)`

GetIsHostOk returns a tuple with the IsHost field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsHost

`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) SetIsHost(v bool)`

SetIsHost sets IsHost field to given value.

### HasIsHost

`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) HasIsHost() bool`

HasIsHost returns a boolean if a field has been set.

### GetEths

`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) GetEths() []SwitchpointsPutRequestSwitchpointValueObjectPropertiesEthsInner`

GetEths returns the Eths field if non-nil, zero value otherwise.

### GetEthsOk

`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) GetEthsOk() (*[]SwitchpointsPutRequestSwitchpointValueObjectPropertiesEthsInner, bool)`

GetEthsOk returns a tuple with the Eths field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEths

`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) SetEths(v []SwitchpointsPutRequestSwitchpointValueObjectPropertiesEthsInner)`

SetEths sets Eths field to given value.

### HasEths

`func (o *SwitchpointsPutRequestEndpointValueObjectProperties) HasEths() bool`

HasEths returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


