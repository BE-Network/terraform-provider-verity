# BundlesPatchRequestEndpointBundleValueUserServicesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**RowAppEnable** | Pointer to **bool** | Enable of this User application | [optional] [default to false]
**RowAppConnectedService** | Pointer to **string** | Service connected to this User application | [optional] [default to ""]
**RowAppConnectedServiceRefType** | Pointer to **string** | Object type for row_app_connected_service field | [optional] 
**RowAppCliCommands** | Pointer to **string** | CLI Commands of this User application | [optional] [default to ""]
**RowIpMask** | Pointer to **string** | IP/Mask | [optional] [default to ""]
**Index** | Pointer to **int32** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 

## Methods

### NewBundlesPatchRequestEndpointBundleValueUserServicesInner

`func NewBundlesPatchRequestEndpointBundleValueUserServicesInner() *BundlesPatchRequestEndpointBundleValueUserServicesInner`

NewBundlesPatchRequestEndpointBundleValueUserServicesInner instantiates a new BundlesPatchRequestEndpointBundleValueUserServicesInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBundlesPatchRequestEndpointBundleValueUserServicesInnerWithDefaults

`func NewBundlesPatchRequestEndpointBundleValueUserServicesInnerWithDefaults() *BundlesPatchRequestEndpointBundleValueUserServicesInner`

NewBundlesPatchRequestEndpointBundleValueUserServicesInnerWithDefaults instantiates a new BundlesPatchRequestEndpointBundleValueUserServicesInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRowAppEnable

`func (o *BundlesPatchRequestEndpointBundleValueUserServicesInner) GetRowAppEnable() bool`

GetRowAppEnable returns the RowAppEnable field if non-nil, zero value otherwise.

### GetRowAppEnableOk

`func (o *BundlesPatchRequestEndpointBundleValueUserServicesInner) GetRowAppEnableOk() (*bool, bool)`

GetRowAppEnableOk returns a tuple with the RowAppEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowAppEnable

`func (o *BundlesPatchRequestEndpointBundleValueUserServicesInner) SetRowAppEnable(v bool)`

SetRowAppEnable sets RowAppEnable field to given value.

### HasRowAppEnable

`func (o *BundlesPatchRequestEndpointBundleValueUserServicesInner) HasRowAppEnable() bool`

HasRowAppEnable returns a boolean if a field has been set.

### GetRowAppConnectedService

`func (o *BundlesPatchRequestEndpointBundleValueUserServicesInner) GetRowAppConnectedService() string`

GetRowAppConnectedService returns the RowAppConnectedService field if non-nil, zero value otherwise.

### GetRowAppConnectedServiceOk

`func (o *BundlesPatchRequestEndpointBundleValueUserServicesInner) GetRowAppConnectedServiceOk() (*string, bool)`

GetRowAppConnectedServiceOk returns a tuple with the RowAppConnectedService field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowAppConnectedService

`func (o *BundlesPatchRequestEndpointBundleValueUserServicesInner) SetRowAppConnectedService(v string)`

SetRowAppConnectedService sets RowAppConnectedService field to given value.

### HasRowAppConnectedService

`func (o *BundlesPatchRequestEndpointBundleValueUserServicesInner) HasRowAppConnectedService() bool`

HasRowAppConnectedService returns a boolean if a field has been set.

### GetRowAppConnectedServiceRefType

`func (o *BundlesPatchRequestEndpointBundleValueUserServicesInner) GetRowAppConnectedServiceRefType() string`

GetRowAppConnectedServiceRefType returns the RowAppConnectedServiceRefType field if non-nil, zero value otherwise.

### GetRowAppConnectedServiceRefTypeOk

`func (o *BundlesPatchRequestEndpointBundleValueUserServicesInner) GetRowAppConnectedServiceRefTypeOk() (*string, bool)`

GetRowAppConnectedServiceRefTypeOk returns a tuple with the RowAppConnectedServiceRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowAppConnectedServiceRefType

`func (o *BundlesPatchRequestEndpointBundleValueUserServicesInner) SetRowAppConnectedServiceRefType(v string)`

SetRowAppConnectedServiceRefType sets RowAppConnectedServiceRefType field to given value.

### HasRowAppConnectedServiceRefType

`func (o *BundlesPatchRequestEndpointBundleValueUserServicesInner) HasRowAppConnectedServiceRefType() bool`

HasRowAppConnectedServiceRefType returns a boolean if a field has been set.

### GetRowAppCliCommands

`func (o *BundlesPatchRequestEndpointBundleValueUserServicesInner) GetRowAppCliCommands() string`

GetRowAppCliCommands returns the RowAppCliCommands field if non-nil, zero value otherwise.

### GetRowAppCliCommandsOk

`func (o *BundlesPatchRequestEndpointBundleValueUserServicesInner) GetRowAppCliCommandsOk() (*string, bool)`

GetRowAppCliCommandsOk returns a tuple with the RowAppCliCommands field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowAppCliCommands

`func (o *BundlesPatchRequestEndpointBundleValueUserServicesInner) SetRowAppCliCommands(v string)`

SetRowAppCliCommands sets RowAppCliCommands field to given value.

### HasRowAppCliCommands

`func (o *BundlesPatchRequestEndpointBundleValueUserServicesInner) HasRowAppCliCommands() bool`

HasRowAppCliCommands returns a boolean if a field has been set.

### GetRowIpMask

`func (o *BundlesPatchRequestEndpointBundleValueUserServicesInner) GetRowIpMask() string`

GetRowIpMask returns the RowIpMask field if non-nil, zero value otherwise.

### GetRowIpMaskOk

`func (o *BundlesPatchRequestEndpointBundleValueUserServicesInner) GetRowIpMaskOk() (*string, bool)`

GetRowIpMaskOk returns a tuple with the RowIpMask field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowIpMask

`func (o *BundlesPatchRequestEndpointBundleValueUserServicesInner) SetRowIpMask(v string)`

SetRowIpMask sets RowIpMask field to given value.

### HasRowIpMask

`func (o *BundlesPatchRequestEndpointBundleValueUserServicesInner) HasRowIpMask() bool`

HasRowIpMask returns a boolean if a field has been set.

### GetIndex

`func (o *BundlesPatchRequestEndpointBundleValueUserServicesInner) GetIndex() int32`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *BundlesPatchRequestEndpointBundleValueUserServicesInner) GetIndexOk() (*int32, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *BundlesPatchRequestEndpointBundleValueUserServicesInner) SetIndex(v int32)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *BundlesPatchRequestEndpointBundleValueUserServicesInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


