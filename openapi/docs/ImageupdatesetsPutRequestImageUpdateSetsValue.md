# ImageupdatesetsPutRequestImageUpdateSetsValue

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
**SectionPointless** | Pointer to [**[]ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner**](ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner.md) |  | [optional] 
**Section** | Pointer to [**[]ImageupdatesetsPutRequestImageUpdateSetsValueSectionInner**](ImageupdatesetsPutRequestImageUpdateSetsValueSectionInner.md) |  | [optional] 
**SectionElse** | Pointer to [**[]ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner**](ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**ImageupdatesetsPutRequestImageUpdateSetsValueObjectProperties**](ImageupdatesetsPutRequestImageUpdateSetsValueObjectProperties.md) |  | [optional] 

## Methods

### NewImageupdatesetsPutRequestImageUpdateSetsValue

`func NewImageupdatesetsPutRequestImageUpdateSetsValue() *ImageupdatesetsPutRequestImageUpdateSetsValue`

NewImageupdatesetsPutRequestImageUpdateSetsValue instantiates a new ImageupdatesetsPutRequestImageUpdateSetsValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewImageupdatesetsPutRequestImageUpdateSetsValueWithDefaults

`func NewImageupdatesetsPutRequestImageUpdateSetsValueWithDefaults() *ImageupdatesetsPutRequestImageUpdateSetsValue`

NewImageupdatesetsPutRequestImageUpdateSetsValueWithDefaults instantiates a new ImageupdatesetsPutRequestImageUpdateSetsValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetUpgraderOnSummary

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) GetUpgraderOnSummary() bool`

GetUpgraderOnSummary returns the UpgraderOnSummary field if non-nil, zero value otherwise.

### GetUpgraderOnSummaryOk

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) GetUpgraderOnSummaryOk() (*bool, bool)`

GetUpgraderOnSummaryOk returns a tuple with the UpgraderOnSummary field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpgraderOnSummary

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) SetUpgraderOnSummary(v bool)`

SetUpgraderOnSummary sets UpgraderOnSummary field to given value.

### HasUpgraderOnSummary

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) HasUpgraderOnSummary() bool`

HasUpgraderOnSummary returns a boolean if a field has been set.

### GetInstallationOnSummary

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) GetInstallationOnSummary() bool`

GetInstallationOnSummary returns the InstallationOnSummary field if non-nil, zero value otherwise.

### GetInstallationOnSummaryOk

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) GetInstallationOnSummaryOk() (*bool, bool)`

GetInstallationOnSummaryOk returns a tuple with the InstallationOnSummary field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstallationOnSummary

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) SetInstallationOnSummary(v bool)`

SetInstallationOnSummary sets InstallationOnSummary field to given value.

### HasInstallationOnSummary

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) HasInstallationOnSummary() bool`

HasInstallationOnSummary returns a boolean if a field has been set.

### GetCommOnSummary

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) GetCommOnSummary() bool`

GetCommOnSummary returns the CommOnSummary field if non-nil, zero value otherwise.

### GetCommOnSummaryOk

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) GetCommOnSummaryOk() (*bool, bool)`

GetCommOnSummaryOk returns a tuple with the CommOnSummary field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommOnSummary

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) SetCommOnSummary(v bool)`

SetCommOnSummary sets CommOnSummary field to given value.

### HasCommOnSummary

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) HasCommOnSummary() bool`

HasCommOnSummary returns a boolean if a field has been set.

### GetProvisioningOnSummary

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) GetProvisioningOnSummary() bool`

GetProvisioningOnSummary returns the ProvisioningOnSummary field if non-nil, zero value otherwise.

### GetProvisioningOnSummaryOk

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) GetProvisioningOnSummaryOk() (*bool, bool)`

GetProvisioningOnSummaryOk returns a tuple with the ProvisioningOnSummary field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProvisioningOnSummary

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) SetProvisioningOnSummary(v bool)`

SetProvisioningOnSummary sets ProvisioningOnSummary field to given value.

### HasProvisioningOnSummary

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) HasProvisioningOnSummary() bool`

HasProvisioningOnSummary returns a boolean if a field has been set.

### GetType

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) HasType() bool`

HasType returns a boolean if a field has been set.

### GetSectionPointless

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) GetSectionPointless() []ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner`

GetSectionPointless returns the SectionPointless field if non-nil, zero value otherwise.

### GetSectionPointlessOk

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) GetSectionPointlessOk() (*[]ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner, bool)`

GetSectionPointlessOk returns a tuple with the SectionPointless field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSectionPointless

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) SetSectionPointless(v []ImageupdatesetsPutRequestImageUpdateSetsValueSectionPointlessInner)`

SetSectionPointless sets SectionPointless field to given value.

### HasSectionPointless

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) HasSectionPointless() bool`

HasSectionPointless returns a boolean if a field has been set.

### GetSection

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) GetSection() []ImageupdatesetsPutRequestImageUpdateSetsValueSectionInner`

GetSection returns the Section field if non-nil, zero value otherwise.

### GetSectionOk

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) GetSectionOk() (*[]ImageupdatesetsPutRequestImageUpdateSetsValueSectionInner, bool)`

GetSectionOk returns a tuple with the Section field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSection

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) SetSection(v []ImageupdatesetsPutRequestImageUpdateSetsValueSectionInner)`

SetSection sets Section field to given value.

### HasSection

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) HasSection() bool`

HasSection returns a boolean if a field has been set.

### GetSectionElse

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) GetSectionElse() []ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner`

GetSectionElse returns the SectionElse field if non-nil, zero value otherwise.

### GetSectionElseOk

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) GetSectionElseOk() (*[]ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner, bool)`

GetSectionElseOk returns a tuple with the SectionElse field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSectionElse

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) SetSectionElse(v []ImageupdatesetsPutRequestImageUpdateSetsValueSectionElseInner)`

SetSectionElse sets SectionElse field to given value.

### HasSectionElse

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) HasSectionElse() bool`

HasSectionElse returns a boolean if a field has been set.

### GetObjectProperties

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) GetObjectProperties() ImageupdatesetsPutRequestImageUpdateSetsValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) GetObjectPropertiesOk() (*ImageupdatesetsPutRequestImageUpdateSetsValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) SetObjectProperties(v ImageupdatesetsPutRequestImageUpdateSetsValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ImageupdatesetsPutRequestImageUpdateSetsValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


