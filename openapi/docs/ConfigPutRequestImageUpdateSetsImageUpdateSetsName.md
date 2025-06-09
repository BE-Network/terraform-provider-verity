# ConfigPutRequestImageUpdateSetsImageUpdateSetsName

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. | [optional] [default to true]
**UpgraderOnSummary** | Pointer to **bool** | Show Upgrader Pie Chart on Summary | [optional] [default to true]
**InstallationOnSummary** | Pointer to **bool** | Show Installation Pie Chart on Summary | [optional] [default to true]
**CommOnSummary** | Pointer to **bool** | Show Comm Pie Chart on Summary | [optional] [default to true]
**ProvisioningOnSummary** | Pointer to **bool** | Show Provisioning Pie Chart on Summary | [optional] [default to true]
**Type** | Pointer to **string** | Type of Image Update Sets | [optional] [default to "whitebox"]
**SectionPointless** | Pointer to [**[]ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner**](ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner.md) |  | [optional] 
**Section** | Pointer to [**[]ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner**](ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner.md) |  | [optional] 
**SectionElse** | Pointer to [**[]ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionElseInner**](ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionElseInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**ConfigPutRequestImageUpdateSetsImageUpdateSetsNameObjectProperties**](ConfigPutRequestImageUpdateSetsImageUpdateSetsNameObjectProperties.md) |  | [optional] 

## Methods

### NewConfigPutRequestImageUpdateSetsImageUpdateSetsName

`func NewConfigPutRequestImageUpdateSetsImageUpdateSetsName() *ConfigPutRequestImageUpdateSetsImageUpdateSetsName`

NewConfigPutRequestImageUpdateSetsImageUpdateSetsName instantiates a new ConfigPutRequestImageUpdateSetsImageUpdateSetsName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestImageUpdateSetsImageUpdateSetsNameWithDefaults

`func NewConfigPutRequestImageUpdateSetsImageUpdateSetsNameWithDefaults() *ConfigPutRequestImageUpdateSetsImageUpdateSetsName`

NewConfigPutRequestImageUpdateSetsImageUpdateSetsNameWithDefaults instantiates a new ConfigPutRequestImageUpdateSetsImageUpdateSetsName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetUpgraderOnSummary

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) GetUpgraderOnSummary() bool`

GetUpgraderOnSummary returns the UpgraderOnSummary field if non-nil, zero value otherwise.

### GetUpgraderOnSummaryOk

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) GetUpgraderOnSummaryOk() (*bool, bool)`

GetUpgraderOnSummaryOk returns a tuple with the UpgraderOnSummary field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpgraderOnSummary

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) SetUpgraderOnSummary(v bool)`

SetUpgraderOnSummary sets UpgraderOnSummary field to given value.

### HasUpgraderOnSummary

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) HasUpgraderOnSummary() bool`

HasUpgraderOnSummary returns a boolean if a field has been set.

### GetInstallationOnSummary

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) GetInstallationOnSummary() bool`

GetInstallationOnSummary returns the InstallationOnSummary field if non-nil, zero value otherwise.

### GetInstallationOnSummaryOk

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) GetInstallationOnSummaryOk() (*bool, bool)`

GetInstallationOnSummaryOk returns a tuple with the InstallationOnSummary field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstallationOnSummary

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) SetInstallationOnSummary(v bool)`

SetInstallationOnSummary sets InstallationOnSummary field to given value.

### HasInstallationOnSummary

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) HasInstallationOnSummary() bool`

HasInstallationOnSummary returns a boolean if a field has been set.

### GetCommOnSummary

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) GetCommOnSummary() bool`

GetCommOnSummary returns the CommOnSummary field if non-nil, zero value otherwise.

### GetCommOnSummaryOk

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) GetCommOnSummaryOk() (*bool, bool)`

GetCommOnSummaryOk returns a tuple with the CommOnSummary field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommOnSummary

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) SetCommOnSummary(v bool)`

SetCommOnSummary sets CommOnSummary field to given value.

### HasCommOnSummary

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) HasCommOnSummary() bool`

HasCommOnSummary returns a boolean if a field has been set.

### GetProvisioningOnSummary

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) GetProvisioningOnSummary() bool`

GetProvisioningOnSummary returns the ProvisioningOnSummary field if non-nil, zero value otherwise.

### GetProvisioningOnSummaryOk

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) GetProvisioningOnSummaryOk() (*bool, bool)`

GetProvisioningOnSummaryOk returns a tuple with the ProvisioningOnSummary field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProvisioningOnSummary

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) SetProvisioningOnSummary(v bool)`

SetProvisioningOnSummary sets ProvisioningOnSummary field to given value.

### HasProvisioningOnSummary

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) HasProvisioningOnSummary() bool`

HasProvisioningOnSummary returns a boolean if a field has been set.

### GetType

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) HasType() bool`

HasType returns a boolean if a field has been set.

### GetSectionPointless

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) GetSectionPointless() []ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner`

GetSectionPointless returns the SectionPointless field if non-nil, zero value otherwise.

### GetSectionPointlessOk

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) GetSectionPointlessOk() (*[]ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner, bool)`

GetSectionPointlessOk returns a tuple with the SectionPointless field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSectionPointless

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) SetSectionPointless(v []ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionPointlessInner)`

SetSectionPointless sets SectionPointless field to given value.

### HasSectionPointless

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) HasSectionPointless() bool`

HasSectionPointless returns a boolean if a field has been set.

### GetSection

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) GetSection() []ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner`

GetSection returns the Section field if non-nil, zero value otherwise.

### GetSectionOk

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) GetSectionOk() (*[]ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner, bool)`

GetSectionOk returns a tuple with the Section field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSection

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) SetSection(v []ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionInner)`

SetSection sets Section field to given value.

### HasSection

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) HasSection() bool`

HasSection returns a boolean if a field has been set.

### GetSectionElse

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) GetSectionElse() []ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionElseInner`

GetSectionElse returns the SectionElse field if non-nil, zero value otherwise.

### GetSectionElseOk

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) GetSectionElseOk() (*[]ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionElseInner, bool)`

GetSectionElseOk returns a tuple with the SectionElse field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSectionElse

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) SetSectionElse(v []ConfigPutRequestImageUpdateSetsImageUpdateSetsNameSectionElseInner)`

SetSectionElse sets SectionElse field to given value.

### HasSectionElse

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) HasSectionElse() bool`

HasSectionElse returns a boolean if a field has been set.

### GetObjectProperties

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) GetObjectProperties() ConfigPutRequestImageUpdateSetsImageUpdateSetsNameObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) GetObjectPropertiesOk() (*ConfigPutRequestImageUpdateSetsImageUpdateSetsNameObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) SetObjectProperties(v ConfigPutRequestImageUpdateSetsImageUpdateSetsNameObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ConfigPutRequestImageUpdateSetsImageUpdateSetsName) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


