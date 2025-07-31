# ServicesPutRequestServiceValueObjectProperties

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Group** | Pointer to **string** | Group | [optional] [default to ""]
**OnSummary** | Pointer to **bool** | Show on the summary view | [optional] [default to true]
**WarnOnNoExternalSource** | Pointer to **bool** | Warn if there is not outbound path for service in SD-Router or a Service Port Profile | [optional] [default to true]

## Methods

### NewServicesPutRequestServiceValueObjectProperties

`func NewServicesPutRequestServiceValueObjectProperties() *ServicesPutRequestServiceValueObjectProperties`

NewServicesPutRequestServiceValueObjectProperties instantiates a new ServicesPutRequestServiceValueObjectProperties object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewServicesPutRequestServiceValueObjectPropertiesWithDefaults

`func NewServicesPutRequestServiceValueObjectPropertiesWithDefaults() *ServicesPutRequestServiceValueObjectProperties`

NewServicesPutRequestServiceValueObjectPropertiesWithDefaults instantiates a new ServicesPutRequestServiceValueObjectProperties object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetGroup

`func (o *ServicesPutRequestServiceValueObjectProperties) GetGroup() string`

GetGroup returns the Group field if non-nil, zero value otherwise.

### GetGroupOk

`func (o *ServicesPutRequestServiceValueObjectProperties) GetGroupOk() (*string, bool)`

GetGroupOk returns a tuple with the Group field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGroup

`func (o *ServicesPutRequestServiceValueObjectProperties) SetGroup(v string)`

SetGroup sets Group field to given value.

### HasGroup

`func (o *ServicesPutRequestServiceValueObjectProperties) HasGroup() bool`

HasGroup returns a boolean if a field has been set.

### GetOnSummary

`func (o *ServicesPutRequestServiceValueObjectProperties) GetOnSummary() bool`

GetOnSummary returns the OnSummary field if non-nil, zero value otherwise.

### GetOnSummaryOk

`func (o *ServicesPutRequestServiceValueObjectProperties) GetOnSummaryOk() (*bool, bool)`

GetOnSummaryOk returns a tuple with the OnSummary field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOnSummary

`func (o *ServicesPutRequestServiceValueObjectProperties) SetOnSummary(v bool)`

SetOnSummary sets OnSummary field to given value.

### HasOnSummary

`func (o *ServicesPutRequestServiceValueObjectProperties) HasOnSummary() bool`

HasOnSummary returns a boolean if a field has been set.

### GetWarnOnNoExternalSource

`func (o *ServicesPutRequestServiceValueObjectProperties) GetWarnOnNoExternalSource() bool`

GetWarnOnNoExternalSource returns the WarnOnNoExternalSource field if non-nil, zero value otherwise.

### GetWarnOnNoExternalSourceOk

`func (o *ServicesPutRequestServiceValueObjectProperties) GetWarnOnNoExternalSourceOk() (*bool, bool)`

GetWarnOnNoExternalSourceOk returns a tuple with the WarnOnNoExternalSource field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWarnOnNoExternalSource

`func (o *ServicesPutRequestServiceValueObjectProperties) SetWarnOnNoExternalSource(v bool)`

SetWarnOnNoExternalSource sets WarnOnNoExternalSource field to given value.

### HasWarnOnNoExternalSource

`func (o *ServicesPutRequestServiceValueObjectProperties) HasWarnOnNoExternalSource() bool`

HasWarnOnNoExternalSource returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


