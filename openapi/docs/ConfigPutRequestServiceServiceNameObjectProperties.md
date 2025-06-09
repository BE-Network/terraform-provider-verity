# ConfigPutRequestServiceServiceNameObjectProperties

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Group** | Pointer to **string** | Group | [optional] [default to ""]
**OnSummary** | Pointer to **bool** | Show on the summary view | [optional] [default to true]
**WarnOnNoExternalSource** | Pointer to **bool** | Warn if there is not outbound path for service in SD-Router or a Service Port Profile | [optional] [default to true]

## Methods

### NewConfigPutRequestServiceServiceNameObjectProperties

`func NewConfigPutRequestServiceServiceNameObjectProperties() *ConfigPutRequestServiceServiceNameObjectProperties`

NewConfigPutRequestServiceServiceNameObjectProperties instantiates a new ConfigPutRequestServiceServiceNameObjectProperties object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestServiceServiceNameObjectPropertiesWithDefaults

`func NewConfigPutRequestServiceServiceNameObjectPropertiesWithDefaults() *ConfigPutRequestServiceServiceNameObjectProperties`

NewConfigPutRequestServiceServiceNameObjectPropertiesWithDefaults instantiates a new ConfigPutRequestServiceServiceNameObjectProperties object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetGroup

`func (o *ConfigPutRequestServiceServiceNameObjectProperties) GetGroup() string`

GetGroup returns the Group field if non-nil, zero value otherwise.

### GetGroupOk

`func (o *ConfigPutRequestServiceServiceNameObjectProperties) GetGroupOk() (*string, bool)`

GetGroupOk returns a tuple with the Group field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGroup

`func (o *ConfigPutRequestServiceServiceNameObjectProperties) SetGroup(v string)`

SetGroup sets Group field to given value.

### HasGroup

`func (o *ConfigPutRequestServiceServiceNameObjectProperties) HasGroup() bool`

HasGroup returns a boolean if a field has been set.

### GetOnSummary

`func (o *ConfigPutRequestServiceServiceNameObjectProperties) GetOnSummary() bool`

GetOnSummary returns the OnSummary field if non-nil, zero value otherwise.

### GetOnSummaryOk

`func (o *ConfigPutRequestServiceServiceNameObjectProperties) GetOnSummaryOk() (*bool, bool)`

GetOnSummaryOk returns a tuple with the OnSummary field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOnSummary

`func (o *ConfigPutRequestServiceServiceNameObjectProperties) SetOnSummary(v bool)`

SetOnSummary sets OnSummary field to given value.

### HasOnSummary

`func (o *ConfigPutRequestServiceServiceNameObjectProperties) HasOnSummary() bool`

HasOnSummary returns a boolean if a field has been set.

### GetWarnOnNoExternalSource

`func (o *ConfigPutRequestServiceServiceNameObjectProperties) GetWarnOnNoExternalSource() bool`

GetWarnOnNoExternalSource returns the WarnOnNoExternalSource field if non-nil, zero value otherwise.

### GetWarnOnNoExternalSourceOk

`func (o *ConfigPutRequestServiceServiceNameObjectProperties) GetWarnOnNoExternalSourceOk() (*bool, bool)`

GetWarnOnNoExternalSourceOk returns a tuple with the WarnOnNoExternalSource field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWarnOnNoExternalSource

`func (o *ConfigPutRequestServiceServiceNameObjectProperties) SetWarnOnNoExternalSource(v bool)`

SetWarnOnNoExternalSource sets WarnOnNoExternalSource field to given value.

### HasWarnOnNoExternalSource

`func (o *ConfigPutRequestServiceServiceNameObjectProperties) HasWarnOnNoExternalSource() bool`

HasWarnOnNoExternalSource returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


