# ImageupdatesetsPatchRequestImageUpdateSetsValue

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
**SectionPointless** | Pointer to [**[]ImageupdatesetsPatchRequestImageUpdateSetsValueSectionPointlessInner**](ImageupdatesetsPatchRequestImageUpdateSetsValueSectionPointlessInner.md) |  | [optional] 
**Section** | Pointer to [**[]ImageupdatesetsPatchRequestImageUpdateSetsValueSectionInner**](ImageupdatesetsPatchRequestImageUpdateSetsValueSectionInner.md) |  | [optional] 
**SectionElse** | Pointer to [**[]ImageupdatesetsPatchRequestImageUpdateSetsValueSectionElseInner**](ImageupdatesetsPatchRequestImageUpdateSetsValueSectionElseInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**ImageupdatesetsPatchRequestImageUpdateSetsValueObjectProperties**](ImageupdatesetsPatchRequestImageUpdateSetsValueObjectProperties.md) |  | [optional] 

## Methods

### NewImageupdatesetsPatchRequestImageUpdateSetsValue

`func NewImageupdatesetsPatchRequestImageUpdateSetsValue() *ImageupdatesetsPatchRequestImageUpdateSetsValue`

NewImageupdatesetsPatchRequestImageUpdateSetsValue instantiates a new ImageupdatesetsPatchRequestImageUpdateSetsValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewImageupdatesetsPatchRequestImageUpdateSetsValueWithDefaults

`func NewImageupdatesetsPatchRequestImageUpdateSetsValueWithDefaults() *ImageupdatesetsPatchRequestImageUpdateSetsValue`

NewImageupdatesetsPatchRequestImageUpdateSetsValueWithDefaults instantiates a new ImageupdatesetsPatchRequestImageUpdateSetsValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetEnable

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetUpgraderOnSummary

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) GetUpgraderOnSummary() bool`

GetUpgraderOnSummary returns the UpgraderOnSummary field if non-nil, zero value otherwise.

### GetUpgraderOnSummaryOk

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) GetUpgraderOnSummaryOk() (*bool, bool)`

GetUpgraderOnSummaryOk returns a tuple with the UpgraderOnSummary field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUpgraderOnSummary

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) SetUpgraderOnSummary(v bool)`

SetUpgraderOnSummary sets UpgraderOnSummary field to given value.

### HasUpgraderOnSummary

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) HasUpgraderOnSummary() bool`

HasUpgraderOnSummary returns a boolean if a field has been set.

### GetInstallationOnSummary

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) GetInstallationOnSummary() bool`

GetInstallationOnSummary returns the InstallationOnSummary field if non-nil, zero value otherwise.

### GetInstallationOnSummaryOk

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) GetInstallationOnSummaryOk() (*bool, bool)`

GetInstallationOnSummaryOk returns a tuple with the InstallationOnSummary field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstallationOnSummary

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) SetInstallationOnSummary(v bool)`

SetInstallationOnSummary sets InstallationOnSummary field to given value.

### HasInstallationOnSummary

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) HasInstallationOnSummary() bool`

HasInstallationOnSummary returns a boolean if a field has been set.

### GetCommOnSummary

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) GetCommOnSummary() bool`

GetCommOnSummary returns the CommOnSummary field if non-nil, zero value otherwise.

### GetCommOnSummaryOk

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) GetCommOnSummaryOk() (*bool, bool)`

GetCommOnSummaryOk returns a tuple with the CommOnSummary field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommOnSummary

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) SetCommOnSummary(v bool)`

SetCommOnSummary sets CommOnSummary field to given value.

### HasCommOnSummary

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) HasCommOnSummary() bool`

HasCommOnSummary returns a boolean if a field has been set.

### GetProvisioningOnSummary

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) GetProvisioningOnSummary() bool`

GetProvisioningOnSummary returns the ProvisioningOnSummary field if non-nil, zero value otherwise.

### GetProvisioningOnSummaryOk

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) GetProvisioningOnSummaryOk() (*bool, bool)`

GetProvisioningOnSummaryOk returns a tuple with the ProvisioningOnSummary field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetProvisioningOnSummary

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) SetProvisioningOnSummary(v bool)`

SetProvisioningOnSummary sets ProvisioningOnSummary field to given value.

### HasProvisioningOnSummary

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) HasProvisioningOnSummary() bool`

HasProvisioningOnSummary returns a boolean if a field has been set.

### GetType

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) HasType() bool`

HasType returns a boolean if a field has been set.

### GetSectionPointless

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) GetSectionPointless() []ImageupdatesetsPatchRequestImageUpdateSetsValueSectionPointlessInner`

GetSectionPointless returns the SectionPointless field if non-nil, zero value otherwise.

### GetSectionPointlessOk

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) GetSectionPointlessOk() (*[]ImageupdatesetsPatchRequestImageUpdateSetsValueSectionPointlessInner, bool)`

GetSectionPointlessOk returns a tuple with the SectionPointless field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSectionPointless

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) SetSectionPointless(v []ImageupdatesetsPatchRequestImageUpdateSetsValueSectionPointlessInner)`

SetSectionPointless sets SectionPointless field to given value.

### HasSectionPointless

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) HasSectionPointless() bool`

HasSectionPointless returns a boolean if a field has been set.

### GetSection

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) GetSection() []ImageupdatesetsPatchRequestImageUpdateSetsValueSectionInner`

GetSection returns the Section field if non-nil, zero value otherwise.

### GetSectionOk

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) GetSectionOk() (*[]ImageupdatesetsPatchRequestImageUpdateSetsValueSectionInner, bool)`

GetSectionOk returns a tuple with the Section field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSection

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) SetSection(v []ImageupdatesetsPatchRequestImageUpdateSetsValueSectionInner)`

SetSection sets Section field to given value.

### HasSection

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) HasSection() bool`

HasSection returns a boolean if a field has been set.

### GetSectionElse

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) GetSectionElse() []ImageupdatesetsPatchRequestImageUpdateSetsValueSectionElseInner`

GetSectionElse returns the SectionElse field if non-nil, zero value otherwise.

### GetSectionElseOk

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) GetSectionElseOk() (*[]ImageupdatesetsPatchRequestImageUpdateSetsValueSectionElseInner, bool)`

GetSectionElseOk returns a tuple with the SectionElse field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSectionElse

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) SetSectionElse(v []ImageupdatesetsPatchRequestImageUpdateSetsValueSectionElseInner)`

SetSectionElse sets SectionElse field to given value.

### HasSectionElse

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) HasSectionElse() bool`

HasSectionElse returns a boolean if a field has been set.

### GetObjectProperties

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) GetObjectProperties() ImageupdatesetsPatchRequestImageUpdateSetsValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) GetObjectPropertiesOk() (*ImageupdatesetsPatchRequestImageUpdateSetsValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) SetObjectProperties(v ImageupdatesetsPatchRequestImageUpdateSetsValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ImageupdatesetsPatchRequestImageUpdateSetsValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


