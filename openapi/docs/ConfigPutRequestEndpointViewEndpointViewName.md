# ConfigPutRequestEndpointViewEndpointViewName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to true]
**Type** | Pointer to **string** | Always Endpoints | [optional] [default to "Endpoints"]
**Location** | Pointer to **string** | One of Remote|Endpoint|Preprovisioned | [optional] [default to ""]
**OnSummary** | Pointer to **bool** | Show on the summary view | [optional] [default to true]
**UpgraderSummary** | Pointer to **bool** | Show Upgrader Pie Chart on Summary | [optional] [default to true]
**InstallationSummary** | Pointer to **bool** | Show Installation Pie Chart on Summary | [optional] [default to true]
**CommSummary** | Pointer to **bool** | Show Comm Pie Chart on Summary | [optional] [default to true]
**ProvisioningSummary** | Pointer to **bool** | Show Provisioning Pie Chart on Summary | [optional] [default to true]
**OrRules** | Pointer to [**[]ConfigPutRequestEndpointViewEndpointViewNameOrRulesInner**](ConfigPutRequestEndpointViewEndpointViewNameOrRulesInner.md) |  | [optional] 
**ObjectProperties** | Pointer to **map[string]interface{}** |  | [optional] 

## Methods

### NewConfigPutRequestEndpointViewEndpointViewName

`func NewConfigPutRequestEndpointViewEndpointViewName() *ConfigPutRequestEndpointViewEndpointViewName`

NewConfigPutRequestEndpointViewEndpointViewName instantiates a new ConfigPutRequestEndpointViewEndpointViewName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestEndpointViewEndpointViewNameWithDefaults

`func NewConfigPutRequestEndpointViewEndpointViewNameWithDefaults() *ConfigPutRequestEndpointViewEndpointViewName`

NewConfigPutRequestEndpointViewEndpointViewNameWithDefaults instantiates a new ConfigPutRequestEndpointViewEndpointViewName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestEndpointViewEndpointViewName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestEndpointViewEndpointViewName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestEndpointViewEndpointViewName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestEndpointViewEndpointViewName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ConfigPutRequestEndpointViewEndpointViewName) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestEndpointViewEndpointViewName) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestEndpointViewEndpointViewName) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestEndpointViewEndpointViewName) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetType

`func (o *ConfigPutRequestEndpointViewEndpointViewName) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ConfigPutRequestEndpointViewEndpointViewName) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ConfigPutRequestEndpointViewEndpointViewName) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *ConfigPutRequestEndpointViewEndpointViewName) HasType() bool`

HasType returns a boolean if a field has been set.

### GetLocation

`func (o *ConfigPutRequestEndpointViewEndpointViewName) GetLocation() string`

GetLocation returns the Location field if non-nil, zero value otherwise.

### GetLocationOk

`func (o *ConfigPutRequestEndpointViewEndpointViewName) GetLocationOk() (*string, bool)`

GetLocationOk returns a tuple with the Location field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocation

`func (o *ConfigPutRequestEndpointViewEndpointViewName) SetLocation(v string)`

SetLocation sets Location field to given value.

### HasLocation

`func (o *ConfigPutRequestEndpointViewEndpointViewName) HasLocation() bool`

HasLocation returns a boolean if a field has been set.

### GetOnSummary

`func (o *ConfigPutRequestEndpointViewEndpointViewName) GetOnSummary() bool`

GetOnSummary returns the OnSummary field if non-nil, zero value otherwise.

### GetOnSummaryOk

`func (o *ConfigPutRequestEndpointViewEndpointViewName) GetOnSummaryOk() (*bool, bool)`

GetOnSummaryOk returns a tuple with the OnSummary field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOnSummary

`func (o *ConfigPutRequestEndpointViewEndpointViewName) SetOnSummary(v bool)`

SetOnSummary sets OnSummary field to given value.

### HasOnSummary

`func (o *ConfigPutRequestEndpointViewEndpointViewName) HasOnSummary() bool`

HasOnSummary returns a boolean if a field has been set.

### GetUpgraderSummary

`func (o *ConfigPutRequestEndpointViewEndpointViewName) GetUpgraderSummary() bool`

GetUpgraderSummary returns the UpgraderSummary field if non-nil, zero value otherwise.

### GetUpgraderSummaryOk

`func (o *ConfigPutRequestEndpointViewEndpointViewName) GetUpgraderSummaryOk() (*bool, bool)`

GetUpgraderSummaryOk returns a tuple with the UpgraderSummary field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpgraderSummary

`func (o *ConfigPutRequestEndpointViewEndpointViewName) SetUpgraderSummary(v bool)`

SetUpgraderSummary sets UpgraderSummary field to given value.

### HasUpgraderSummary

`func (o *ConfigPutRequestEndpointViewEndpointViewName) HasUpgraderSummary() bool`

HasUpgraderSummary returns a boolean if a field has been set.

### GetInstallationSummary

`func (o *ConfigPutRequestEndpointViewEndpointViewName) GetInstallationSummary() bool`

GetInstallationSummary returns the InstallationSummary field if non-nil, zero value otherwise.

### GetInstallationSummaryOk

`func (o *ConfigPutRequestEndpointViewEndpointViewName) GetInstallationSummaryOk() (*bool, bool)`

GetInstallationSummaryOk returns a tuple with the InstallationSummary field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstallationSummary

`func (o *ConfigPutRequestEndpointViewEndpointViewName) SetInstallationSummary(v bool)`

SetInstallationSummary sets InstallationSummary field to given value.

### HasInstallationSummary

`func (o *ConfigPutRequestEndpointViewEndpointViewName) HasInstallationSummary() bool`

HasInstallationSummary returns a boolean if a field has been set.

### GetCommSummary

`func (o *ConfigPutRequestEndpointViewEndpointViewName) GetCommSummary() bool`

GetCommSummary returns the CommSummary field if non-nil, zero value otherwise.

### GetCommSummaryOk

`func (o *ConfigPutRequestEndpointViewEndpointViewName) GetCommSummaryOk() (*bool, bool)`

GetCommSummaryOk returns a tuple with the CommSummary field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommSummary

`func (o *ConfigPutRequestEndpointViewEndpointViewName) SetCommSummary(v bool)`

SetCommSummary sets CommSummary field to given value.

### HasCommSummary

`func (o *ConfigPutRequestEndpointViewEndpointViewName) HasCommSummary() bool`

HasCommSummary returns a boolean if a field has been set.

### GetProvisioningSummary

`func (o *ConfigPutRequestEndpointViewEndpointViewName) GetProvisioningSummary() bool`

GetProvisioningSummary returns the ProvisioningSummary field if non-nil, zero value otherwise.

### GetProvisioningSummaryOk

`func (o *ConfigPutRequestEndpointViewEndpointViewName) GetProvisioningSummaryOk() (*bool, bool)`

GetProvisioningSummaryOk returns a tuple with the ProvisioningSummary field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProvisioningSummary

`func (o *ConfigPutRequestEndpointViewEndpointViewName) SetProvisioningSummary(v bool)`

SetProvisioningSummary sets ProvisioningSummary field to given value.

### HasProvisioningSummary

`func (o *ConfigPutRequestEndpointViewEndpointViewName) HasProvisioningSummary() bool`

HasProvisioningSummary returns a boolean if a field has been set.

### GetOrRules

`func (o *ConfigPutRequestEndpointViewEndpointViewName) GetOrRules() []ConfigPutRequestEndpointViewEndpointViewNameOrRulesInner`

GetOrRules returns the OrRules field if non-nil, zero value otherwise.

### GetOrRulesOk

`func (o *ConfigPutRequestEndpointViewEndpointViewName) GetOrRulesOk() (*[]ConfigPutRequestEndpointViewEndpointViewNameOrRulesInner, bool)`

GetOrRulesOk returns a tuple with the OrRules field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOrRules

`func (o *ConfigPutRequestEndpointViewEndpointViewName) SetOrRules(v []ConfigPutRequestEndpointViewEndpointViewNameOrRulesInner)`

SetOrRules sets OrRules field to given value.

### HasOrRules

`func (o *ConfigPutRequestEndpointViewEndpointViewName) HasOrRules() bool`

HasOrRules returns a boolean if a field has been set.

### GetObjectProperties

`func (o *ConfigPutRequestEndpointViewEndpointViewName) GetObjectProperties() map[string]interface{}`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ConfigPutRequestEndpointViewEndpointViewName) GetObjectPropertiesOk() (*map[string]interface{}, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ConfigPutRequestEndpointViewEndpointViewName) SetObjectProperties(v map[string]interface{})`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ConfigPutRequestEndpointViewEndpointViewName) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


