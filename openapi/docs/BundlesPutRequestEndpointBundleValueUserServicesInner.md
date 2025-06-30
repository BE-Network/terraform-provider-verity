# BundlesPutRequestEndpointBundleValueUserServicesInner

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

### NewBundlesPutRequestEndpointBundleValueUserServicesInner

`func NewBundlesPutRequestEndpointBundleValueUserServicesInner() *BundlesPutRequestEndpointBundleValueUserServicesInner`

NewBundlesPutRequestEndpointBundleValueUserServicesInner instantiates a new BundlesPutRequestEndpointBundleValueUserServicesInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBundlesPutRequestEndpointBundleValueUserServicesInnerWithDefaults

`func NewBundlesPutRequestEndpointBundleValueUserServicesInnerWithDefaults() *BundlesPutRequestEndpointBundleValueUserServicesInner`

NewBundlesPutRequestEndpointBundleValueUserServicesInnerWithDefaults instantiates a new BundlesPutRequestEndpointBundleValueUserServicesInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRowAppEnable

`func (o *BundlesPutRequestEndpointBundleValueUserServicesInner) GetRowAppEnable() bool`

GetRowAppEnable returns the RowAppEnable field if non-nil, zero value otherwise.

### GetRowAppEnableOk

`func (o *BundlesPutRequestEndpointBundleValueUserServicesInner) GetRowAppEnableOk() (*bool, bool)`

GetRowAppEnableOk returns a tuple with the RowAppEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowAppEnable

`func (o *BundlesPutRequestEndpointBundleValueUserServicesInner) SetRowAppEnable(v bool)`

SetRowAppEnable sets RowAppEnable field to given value.

### HasRowAppEnable

`func (o *BundlesPutRequestEndpointBundleValueUserServicesInner) HasRowAppEnable() bool`

HasRowAppEnable returns a boolean if a field has been set.

### GetRowAppConnectedService

`func (o *BundlesPutRequestEndpointBundleValueUserServicesInner) GetRowAppConnectedService() string`

GetRowAppConnectedService returns the RowAppConnectedService field if non-nil, zero value otherwise.

### GetRowAppConnectedServiceOk

`func (o *BundlesPutRequestEndpointBundleValueUserServicesInner) GetRowAppConnectedServiceOk() (*string, bool)`

GetRowAppConnectedServiceOk returns a tuple with the RowAppConnectedService field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowAppConnectedService

`func (o *BundlesPutRequestEndpointBundleValueUserServicesInner) SetRowAppConnectedService(v string)`

SetRowAppConnectedService sets RowAppConnectedService field to given value.

### HasRowAppConnectedService

`func (o *BundlesPutRequestEndpointBundleValueUserServicesInner) HasRowAppConnectedService() bool`

HasRowAppConnectedService returns a boolean if a field has been set.

### GetRowAppConnectedServiceRefType

`func (o *BundlesPutRequestEndpointBundleValueUserServicesInner) GetRowAppConnectedServiceRefType() string`

GetRowAppConnectedServiceRefType returns the RowAppConnectedServiceRefType field if non-nil, zero value otherwise.

### GetRowAppConnectedServiceRefTypeOk

`func (o *BundlesPutRequestEndpointBundleValueUserServicesInner) GetRowAppConnectedServiceRefTypeOk() (*string, bool)`

GetRowAppConnectedServiceRefTypeOk returns a tuple with the RowAppConnectedServiceRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowAppConnectedServiceRefType

`func (o *BundlesPutRequestEndpointBundleValueUserServicesInner) SetRowAppConnectedServiceRefType(v string)`

SetRowAppConnectedServiceRefType sets RowAppConnectedServiceRefType field to given value.

### HasRowAppConnectedServiceRefType

`func (o *BundlesPutRequestEndpointBundleValueUserServicesInner) HasRowAppConnectedServiceRefType() bool`

HasRowAppConnectedServiceRefType returns a boolean if a field has been set.

### GetRowAppCliCommands

`func (o *BundlesPutRequestEndpointBundleValueUserServicesInner) GetRowAppCliCommands() string`

GetRowAppCliCommands returns the RowAppCliCommands field if non-nil, zero value otherwise.

### GetRowAppCliCommandsOk

`func (o *BundlesPutRequestEndpointBundleValueUserServicesInner) GetRowAppCliCommandsOk() (*string, bool)`

GetRowAppCliCommandsOk returns a tuple with the RowAppCliCommands field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowAppCliCommands

`func (o *BundlesPutRequestEndpointBundleValueUserServicesInner) SetRowAppCliCommands(v string)`

SetRowAppCliCommands sets RowAppCliCommands field to given value.

### HasRowAppCliCommands

`func (o *BundlesPutRequestEndpointBundleValueUserServicesInner) HasRowAppCliCommands() bool`

HasRowAppCliCommands returns a boolean if a field has been set.

### GetRowIpMask

`func (o *BundlesPutRequestEndpointBundleValueUserServicesInner) GetRowIpMask() string`

GetRowIpMask returns the RowIpMask field if non-nil, zero value otherwise.

### GetRowIpMaskOk

`func (o *BundlesPutRequestEndpointBundleValueUserServicesInner) GetRowIpMaskOk() (*string, bool)`

GetRowIpMaskOk returns a tuple with the RowIpMask field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowIpMask

`func (o *BundlesPutRequestEndpointBundleValueUserServicesInner) SetRowIpMask(v string)`

SetRowIpMask sets RowIpMask field to given value.

### HasRowIpMask

`func (o *BundlesPutRequestEndpointBundleValueUserServicesInner) HasRowIpMask() bool`

HasRowIpMask returns a boolean if a field has been set.

### GetIndex

`func (o *BundlesPutRequestEndpointBundleValueUserServicesInner) GetIndex() int32`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *BundlesPutRequestEndpointBundleValueUserServicesInner) GetIndexOk() (*int32, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *BundlesPutRequestEndpointBundleValueUserServicesInner) SetIndex(v int32)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *BundlesPutRequestEndpointBundleValueUserServicesInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


