# ConfigPutRequestEndpointEndpointNameObjectProperties

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
**Eths** | Pointer to [**ConfigPutRequestSwitchpointSwitchpointNameObjectPropertiesEths**](ConfigPutRequestSwitchpointSwitchpointNameObjectPropertiesEths.md) |  | [optional] 

## Methods

### NewConfigPutRequestEndpointEndpointNameObjectProperties

`func NewConfigPutRequestEndpointEndpointNameObjectProperties() *ConfigPutRequestEndpointEndpointNameObjectProperties`

NewConfigPutRequestEndpointEndpointNameObjectProperties instantiates a new ConfigPutRequestEndpointEndpointNameObjectProperties object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestEndpointEndpointNameObjectPropertiesWithDefaults

`func NewConfigPutRequestEndpointEndpointNameObjectPropertiesWithDefaults() *ConfigPutRequestEndpointEndpointNameObjectProperties`

NewConfigPutRequestEndpointEndpointNameObjectPropertiesWithDefaults instantiates a new ConfigPutRequestEndpointEndpointNameObjectProperties object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUserNotes

`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) GetUserNotes() string`

GetUserNotes returns the UserNotes field if non-nil, zero value otherwise.

### GetUserNotesOk

`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) GetUserNotesOk() (*string, bool)`

GetUserNotesOk returns a tuple with the UserNotes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserNotes

`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) SetUserNotes(v string)`

SetUserNotes sets UserNotes field to given value.

### HasUserNotes

`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) HasUserNotes() bool`

HasUserNotes returns a boolean if a field has been set.

### GetExpectedParentEndpoint

`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) GetExpectedParentEndpoint() string`

GetExpectedParentEndpoint returns the ExpectedParentEndpoint field if non-nil, zero value otherwise.

### GetExpectedParentEndpointOk

`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) GetExpectedParentEndpointOk() (*string, bool)`

GetExpectedParentEndpointOk returns a tuple with the ExpectedParentEndpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpectedParentEndpoint

`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) SetExpectedParentEndpoint(v string)`

SetExpectedParentEndpoint sets ExpectedParentEndpoint field to given value.

### HasExpectedParentEndpoint

`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) HasExpectedParentEndpoint() bool`

HasExpectedParentEndpoint returns a boolean if a field has been set.

### GetExpectedParentEndpointRefType

`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) GetExpectedParentEndpointRefType() string`

GetExpectedParentEndpointRefType returns the ExpectedParentEndpointRefType field if non-nil, zero value otherwise.

### GetExpectedParentEndpointRefTypeOk

`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) GetExpectedParentEndpointRefTypeOk() (*string, bool)`

GetExpectedParentEndpointRefTypeOk returns a tuple with the ExpectedParentEndpointRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpectedParentEndpointRefType

`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) SetExpectedParentEndpointRefType(v string)`

SetExpectedParentEndpointRefType sets ExpectedParentEndpointRefType field to given value.

### HasExpectedParentEndpointRefType

`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) HasExpectedParentEndpointRefType() bool`

HasExpectedParentEndpointRefType returns a boolean if a field has been set.

### GetNumberOfMultipoints

`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) GetNumberOfMultipoints() int32`

GetNumberOfMultipoints returns the NumberOfMultipoints field if non-nil, zero value otherwise.

### GetNumberOfMultipointsOk

`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) GetNumberOfMultipointsOk() (*int32, bool)`

GetNumberOfMultipointsOk returns a tuple with the NumberOfMultipoints field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNumberOfMultipoints

`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) SetNumberOfMultipoints(v int32)`

SetNumberOfMultipoints sets NumberOfMultipoints field to given value.

### HasNumberOfMultipoints

`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) HasNumberOfMultipoints() bool`

HasNumberOfMultipoints returns a boolean if a field has been set.

### SetNumberOfMultipointsNil

`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) SetNumberOfMultipointsNil(b bool)`

 SetNumberOfMultipointsNil sets the value for NumberOfMultipoints to be an explicit nil

### UnsetNumberOfMultipoints
`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) UnsetNumberOfMultipoints()`

UnsetNumberOfMultipoints ensures that no value is present for NumberOfMultipoints, not even an explicit nil
### GetDrawAsEdgeDevice

`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) GetDrawAsEdgeDevice() bool`

GetDrawAsEdgeDevice returns the DrawAsEdgeDevice field if non-nil, zero value otherwise.

### GetDrawAsEdgeDeviceOk

`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) GetDrawAsEdgeDeviceOk() (*bool, bool)`

GetDrawAsEdgeDeviceOk returns a tuple with the DrawAsEdgeDevice field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDrawAsEdgeDevice

`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) SetDrawAsEdgeDevice(v bool)`

SetDrawAsEdgeDevice sets DrawAsEdgeDevice field to given value.

### HasDrawAsEdgeDevice

`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) HasDrawAsEdgeDevice() bool`

HasDrawAsEdgeDevice returns a boolean if a field has been set.

### GetAggregate

`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) GetAggregate() bool`

GetAggregate returns the Aggregate field if non-nil, zero value otherwise.

### GetAggregateOk

`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) GetAggregateOk() (*bool, bool)`

GetAggregateOk returns a tuple with the Aggregate field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAggregate

`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) SetAggregate(v bool)`

SetAggregate sets Aggregate field to given value.

### HasAggregate

`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) HasAggregate() bool`

HasAggregate returns a boolean if a field has been set.

### GetIsHost

`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) GetIsHost() bool`

GetIsHost returns the IsHost field if non-nil, zero value otherwise.

### GetIsHostOk

`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) GetIsHostOk() (*bool, bool)`

GetIsHostOk returns a tuple with the IsHost field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsHost

`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) SetIsHost(v bool)`

SetIsHost sets IsHost field to given value.

### HasIsHost

`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) HasIsHost() bool`

HasIsHost returns a boolean if a field has been set.

### GetEths

`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) GetEths() ConfigPutRequestSwitchpointSwitchpointNameObjectPropertiesEths`

GetEths returns the Eths field if non-nil, zero value otherwise.

### GetEthsOk

`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) GetEthsOk() (*ConfigPutRequestSwitchpointSwitchpointNameObjectPropertiesEths, bool)`

GetEthsOk returns a tuple with the Eths field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEths

`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) SetEths(v ConfigPutRequestSwitchpointSwitchpointNameObjectPropertiesEths)`

SetEths sets Eths field to given value.

### HasEths

`func (o *ConfigPutRequestEndpointEndpointNameObjectProperties) HasEths() bool`

HasEths returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


