# SwitchpointsPutRequestSwitchpointValueObjectProperties

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**UserNotes** | Pointer to **string** | Notes writen by User about the site | [optional] [default to ""]
**ExpectedParentEndpoint** | Pointer to **string** | Expected Parent Endpoint | [optional] [default to ""]
**ExpectedParentEndpointRefType** | Pointer to **string** | Object type for expected_parent_endpoint field | [optional] 
**NumberOfMultipoints** | Pointer to **NullableInt32** | Number of Multipoints | [optional] [default to 0]
**Aggregate** | Pointer to **bool** | For Switch Endpoints. Denotes switch aggregated with all of its sub switches | [optional] [default to false]
**IsHost** | Pointer to **bool** | For Switch Endpoints. Denotes the Host Switch | [optional] [default to false]
**Eths** | Pointer to [**[]ConfigPutRequestEndpointEndpointNameObjectPropertiesEthsInner**](ConfigPutRequestEndpointEndpointNameObjectPropertiesEthsInner.md) |  | [optional] 

## Methods

### NewSwitchpointsPutRequestSwitchpointValueObjectProperties

`func NewSwitchpointsPutRequestSwitchpointValueObjectProperties() *SwitchpointsPutRequestSwitchpointValueObjectProperties`

NewSwitchpointsPutRequestSwitchpointValueObjectProperties instantiates a new SwitchpointsPutRequestSwitchpointValueObjectProperties object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSwitchpointsPutRequestSwitchpointValueObjectPropertiesWithDefaults

`func NewSwitchpointsPutRequestSwitchpointValueObjectPropertiesWithDefaults() *SwitchpointsPutRequestSwitchpointValueObjectProperties`

NewSwitchpointsPutRequestSwitchpointValueObjectPropertiesWithDefaults instantiates a new SwitchpointsPutRequestSwitchpointValueObjectProperties object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUserNotes

`func (o *SwitchpointsPutRequestSwitchpointValueObjectProperties) GetUserNotes() string`

GetUserNotes returns the UserNotes field if non-nil, zero value otherwise.

### GetUserNotesOk

`func (o *SwitchpointsPutRequestSwitchpointValueObjectProperties) GetUserNotesOk() (*string, bool)`

GetUserNotesOk returns a tuple with the UserNotes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserNotes

`func (o *SwitchpointsPutRequestSwitchpointValueObjectProperties) SetUserNotes(v string)`

SetUserNotes sets UserNotes field to given value.

### HasUserNotes

`func (o *SwitchpointsPutRequestSwitchpointValueObjectProperties) HasUserNotes() bool`

HasUserNotes returns a boolean if a field has been set.

### GetExpectedParentEndpoint

`func (o *SwitchpointsPutRequestSwitchpointValueObjectProperties) GetExpectedParentEndpoint() string`

GetExpectedParentEndpoint returns the ExpectedParentEndpoint field if non-nil, zero value otherwise.

### GetExpectedParentEndpointOk

`func (o *SwitchpointsPutRequestSwitchpointValueObjectProperties) GetExpectedParentEndpointOk() (*string, bool)`

GetExpectedParentEndpointOk returns a tuple with the ExpectedParentEndpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpectedParentEndpoint

`func (o *SwitchpointsPutRequestSwitchpointValueObjectProperties) SetExpectedParentEndpoint(v string)`

SetExpectedParentEndpoint sets ExpectedParentEndpoint field to given value.

### HasExpectedParentEndpoint

`func (o *SwitchpointsPutRequestSwitchpointValueObjectProperties) HasExpectedParentEndpoint() bool`

HasExpectedParentEndpoint returns a boolean if a field has been set.

### GetExpectedParentEndpointRefType

`func (o *SwitchpointsPutRequestSwitchpointValueObjectProperties) GetExpectedParentEndpointRefType() string`

GetExpectedParentEndpointRefType returns the ExpectedParentEndpointRefType field if non-nil, zero value otherwise.

### GetExpectedParentEndpointRefTypeOk

`func (o *SwitchpointsPutRequestSwitchpointValueObjectProperties) GetExpectedParentEndpointRefTypeOk() (*string, bool)`

GetExpectedParentEndpointRefTypeOk returns a tuple with the ExpectedParentEndpointRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpectedParentEndpointRefType

`func (o *SwitchpointsPutRequestSwitchpointValueObjectProperties) SetExpectedParentEndpointRefType(v string)`

SetExpectedParentEndpointRefType sets ExpectedParentEndpointRefType field to given value.

### HasExpectedParentEndpointRefType

`func (o *SwitchpointsPutRequestSwitchpointValueObjectProperties) HasExpectedParentEndpointRefType() bool`

HasExpectedParentEndpointRefType returns a boolean if a field has been set.

### GetNumberOfMultipoints

`func (o *SwitchpointsPutRequestSwitchpointValueObjectProperties) GetNumberOfMultipoints() int32`

GetNumberOfMultipoints returns the NumberOfMultipoints field if non-nil, zero value otherwise.

### GetNumberOfMultipointsOk

`func (o *SwitchpointsPutRequestSwitchpointValueObjectProperties) GetNumberOfMultipointsOk() (*int32, bool)`

GetNumberOfMultipointsOk returns a tuple with the NumberOfMultipoints field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNumberOfMultipoints

`func (o *SwitchpointsPutRequestSwitchpointValueObjectProperties) SetNumberOfMultipoints(v int32)`

SetNumberOfMultipoints sets NumberOfMultipoints field to given value.

### HasNumberOfMultipoints

`func (o *SwitchpointsPutRequestSwitchpointValueObjectProperties) HasNumberOfMultipoints() bool`

HasNumberOfMultipoints returns a boolean if a field has been set.

### SetNumberOfMultipointsNil

`func (o *SwitchpointsPutRequestSwitchpointValueObjectProperties) SetNumberOfMultipointsNil(b bool)`

 SetNumberOfMultipointsNil sets the value for NumberOfMultipoints to be an explicit nil

### UnsetNumberOfMultipoints
`func (o *SwitchpointsPutRequestSwitchpointValueObjectProperties) UnsetNumberOfMultipoints()`

UnsetNumberOfMultipoints ensures that no value is present for NumberOfMultipoints, not even an explicit nil
### GetAggregate

`func (o *SwitchpointsPutRequestSwitchpointValueObjectProperties) GetAggregate() bool`

GetAggregate returns the Aggregate field if non-nil, zero value otherwise.

### GetAggregateOk

`func (o *SwitchpointsPutRequestSwitchpointValueObjectProperties) GetAggregateOk() (*bool, bool)`

GetAggregateOk returns a tuple with the Aggregate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAggregate

`func (o *SwitchpointsPutRequestSwitchpointValueObjectProperties) SetAggregate(v bool)`

SetAggregate sets Aggregate field to given value.

### HasAggregate

`func (o *SwitchpointsPutRequestSwitchpointValueObjectProperties) HasAggregate() bool`

HasAggregate returns a boolean if a field has been set.

### GetIsHost

`func (o *SwitchpointsPutRequestSwitchpointValueObjectProperties) GetIsHost() bool`

GetIsHost returns the IsHost field if non-nil, zero value otherwise.

### GetIsHostOk

`func (o *SwitchpointsPutRequestSwitchpointValueObjectProperties) GetIsHostOk() (*bool, bool)`

GetIsHostOk returns a tuple with the IsHost field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsHost

`func (o *SwitchpointsPutRequestSwitchpointValueObjectProperties) SetIsHost(v bool)`

SetIsHost sets IsHost field to given value.

### HasIsHost

`func (o *SwitchpointsPutRequestSwitchpointValueObjectProperties) HasIsHost() bool`

HasIsHost returns a boolean if a field has been set.

### GetEths

`func (o *SwitchpointsPutRequestSwitchpointValueObjectProperties) GetEths() []ConfigPutRequestEndpointEndpointNameObjectPropertiesEthsInner`

GetEths returns the Eths field if non-nil, zero value otherwise.

### GetEthsOk

`func (o *SwitchpointsPutRequestSwitchpointValueObjectProperties) GetEthsOk() (*[]ConfigPutRequestEndpointEndpointNameObjectPropertiesEthsInner, bool)`

GetEthsOk returns a tuple with the Eths field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEths

`func (o *SwitchpointsPutRequestSwitchpointValueObjectProperties) SetEths(v []ConfigPutRequestEndpointEndpointNameObjectPropertiesEthsInner)`

SetEths sets Eths field to given value.

### HasEths

`func (o *SwitchpointsPutRequestSwitchpointValueObjectProperties) HasEths() bool`

HasEths returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


