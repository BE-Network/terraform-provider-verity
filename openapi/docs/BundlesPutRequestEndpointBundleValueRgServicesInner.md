# BundlesPutRequestEndpointBundleValueRgServicesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**RowIpMask** | Pointer to **string** | IP/Mask | [optional] [default to ""]
**Index** | Pointer to **int32** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 
**RowAppEnable** | Pointer to **bool** | Enable of this ONT application | [optional] [default to false]
**RowAppConnectedService** | Pointer to **string** | Service connected to this ONT application | [optional] [default to ""]
**RowAppConnectedServiceRefType** | Pointer to **string** | Object type for row_app_connected_service field | [optional] 
**RowAppType** | Pointer to **string** | Type of ONT Application | [optional] [default to ""]

## Methods

### NewBundlesPutRequestEndpointBundleValueRgServicesInner

`func NewBundlesPutRequestEndpointBundleValueRgServicesInner() *BundlesPutRequestEndpointBundleValueRgServicesInner`

NewBundlesPutRequestEndpointBundleValueRgServicesInner instantiates a new BundlesPutRequestEndpointBundleValueRgServicesInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBundlesPutRequestEndpointBundleValueRgServicesInnerWithDefaults

`func NewBundlesPutRequestEndpointBundleValueRgServicesInnerWithDefaults() *BundlesPutRequestEndpointBundleValueRgServicesInner`

NewBundlesPutRequestEndpointBundleValueRgServicesInnerWithDefaults instantiates a new BundlesPutRequestEndpointBundleValueRgServicesInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRowIpMask

`func (o *BundlesPutRequestEndpointBundleValueRgServicesInner) GetRowIpMask() string`

GetRowIpMask returns the RowIpMask field if non-nil, zero value otherwise.

### GetRowIpMaskOk

`func (o *BundlesPutRequestEndpointBundleValueRgServicesInner) GetRowIpMaskOk() (*string, bool)`

GetRowIpMaskOk returns a tuple with the RowIpMask field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowIpMask

`func (o *BundlesPutRequestEndpointBundleValueRgServicesInner) SetRowIpMask(v string)`

SetRowIpMask sets RowIpMask field to given value.

### HasRowIpMask

`func (o *BundlesPutRequestEndpointBundleValueRgServicesInner) HasRowIpMask() bool`

HasRowIpMask returns a boolean if a field has been set.

### GetIndex

`func (o *BundlesPutRequestEndpointBundleValueRgServicesInner) GetIndex() int32`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *BundlesPutRequestEndpointBundleValueRgServicesInner) GetIndexOk() (*int32, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *BundlesPutRequestEndpointBundleValueRgServicesInner) SetIndex(v int32)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *BundlesPutRequestEndpointBundleValueRgServicesInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.

### GetRowAppEnable

`func (o *BundlesPutRequestEndpointBundleValueRgServicesInner) GetRowAppEnable() bool`

GetRowAppEnable returns the RowAppEnable field if non-nil, zero value otherwise.

### GetRowAppEnableOk

`func (o *BundlesPutRequestEndpointBundleValueRgServicesInner) GetRowAppEnableOk() (*bool, bool)`

GetRowAppEnableOk returns a tuple with the RowAppEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowAppEnable

`func (o *BundlesPutRequestEndpointBundleValueRgServicesInner) SetRowAppEnable(v bool)`

SetRowAppEnable sets RowAppEnable field to given value.

### HasRowAppEnable

`func (o *BundlesPutRequestEndpointBundleValueRgServicesInner) HasRowAppEnable() bool`

HasRowAppEnable returns a boolean if a field has been set.

### GetRowAppConnectedService

`func (o *BundlesPutRequestEndpointBundleValueRgServicesInner) GetRowAppConnectedService() string`

GetRowAppConnectedService returns the RowAppConnectedService field if non-nil, zero value otherwise.

### GetRowAppConnectedServiceOk

`func (o *BundlesPutRequestEndpointBundleValueRgServicesInner) GetRowAppConnectedServiceOk() (*string, bool)`

GetRowAppConnectedServiceOk returns a tuple with the RowAppConnectedService field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowAppConnectedService

`func (o *BundlesPutRequestEndpointBundleValueRgServicesInner) SetRowAppConnectedService(v string)`

SetRowAppConnectedService sets RowAppConnectedService field to given value.

### HasRowAppConnectedService

`func (o *BundlesPutRequestEndpointBundleValueRgServicesInner) HasRowAppConnectedService() bool`

HasRowAppConnectedService returns a boolean if a field has been set.

### GetRowAppConnectedServiceRefType

`func (o *BundlesPutRequestEndpointBundleValueRgServicesInner) GetRowAppConnectedServiceRefType() string`

GetRowAppConnectedServiceRefType returns the RowAppConnectedServiceRefType field if non-nil, zero value otherwise.

### GetRowAppConnectedServiceRefTypeOk

`func (o *BundlesPutRequestEndpointBundleValueRgServicesInner) GetRowAppConnectedServiceRefTypeOk() (*string, bool)`

GetRowAppConnectedServiceRefTypeOk returns a tuple with the RowAppConnectedServiceRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowAppConnectedServiceRefType

`func (o *BundlesPutRequestEndpointBundleValueRgServicesInner) SetRowAppConnectedServiceRefType(v string)`

SetRowAppConnectedServiceRefType sets RowAppConnectedServiceRefType field to given value.

### HasRowAppConnectedServiceRefType

`func (o *BundlesPutRequestEndpointBundleValueRgServicesInner) HasRowAppConnectedServiceRefType() bool`

HasRowAppConnectedServiceRefType returns a boolean if a field has been set.

### GetRowAppType

`func (o *BundlesPutRequestEndpointBundleValueRgServicesInner) GetRowAppType() string`

GetRowAppType returns the RowAppType field if non-nil, zero value otherwise.

### GetRowAppTypeOk

`func (o *BundlesPutRequestEndpointBundleValueRgServicesInner) GetRowAppTypeOk() (*string, bool)`

GetRowAppTypeOk returns a tuple with the RowAppType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowAppType

`func (o *BundlesPutRequestEndpointBundleValueRgServicesInner) SetRowAppType(v string)`

SetRowAppType sets RowAppType field to given value.

### HasRowAppType

`func (o *BundlesPutRequestEndpointBundleValueRgServicesInner) HasRowAppType() bool`

HasRowAppType returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


