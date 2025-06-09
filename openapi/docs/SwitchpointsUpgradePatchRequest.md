# SwitchpointsUpgradePatchRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**PackageVersion** | **string** | Version to upgrade to | [default to ""]
**DeviceNames** | **[]string** |  | 

## Methods

### NewSwitchpointsUpgradePatchRequest

`func NewSwitchpointsUpgradePatchRequest(packageVersion string, deviceNames []string, ) *SwitchpointsUpgradePatchRequest`

NewSwitchpointsUpgradePatchRequest instantiates a new SwitchpointsUpgradePatchRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSwitchpointsUpgradePatchRequestWithDefaults

`func NewSwitchpointsUpgradePatchRequestWithDefaults() *SwitchpointsUpgradePatchRequest`

NewSwitchpointsUpgradePatchRequestWithDefaults instantiates a new SwitchpointsUpgradePatchRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPackageVersion

`func (o *SwitchpointsUpgradePatchRequest) GetPackageVersion() string`

GetPackageVersion returns the PackageVersion field if non-nil, zero value otherwise.

### GetPackageVersionOk

`func (o *SwitchpointsUpgradePatchRequest) GetPackageVersionOk() (*string, bool)`

GetPackageVersionOk returns a tuple with the PackageVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPackageVersion

`func (o *SwitchpointsUpgradePatchRequest) SetPackageVersion(v string)`

SetPackageVersion sets PackageVersion field to given value.


### GetDeviceNames

`func (o *SwitchpointsUpgradePatchRequest) GetDeviceNames() []string`

GetDeviceNames returns the DeviceNames field if non-nil, zero value otherwise.

### GetDeviceNamesOk

`func (o *SwitchpointsUpgradePatchRequest) GetDeviceNamesOk() (*[]string, bool)`

GetDeviceNamesOk returns a tuple with the DeviceNames field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceNames

`func (o *SwitchpointsUpgradePatchRequest) SetDeviceNames(v []string)`

SetDeviceNames sets DeviceNames field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


