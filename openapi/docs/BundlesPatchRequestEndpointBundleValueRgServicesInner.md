# BundlesPatchRequestEndpointBundleValueRgServicesInner

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**IpMask** | Pointer to **string** | IP/Mask | [optional] [default to ""]
**Index** | Pointer to **int32** | The index identifying the object. Zero if you want to add an object to the list. | [optional] 
**RowAppEnable** | Pointer to **bool** | Enable of this ONT application | [optional] [default to false]
**RowAppConnectedService** | Pointer to **string** | Service connected to this ONT application | [optional] [default to ""]
**RowAppConnectedServiceRefType** | Pointer to **string** | Object type for row_app_connected_service field | [optional] 
**RowAppType** | Pointer to **string** | Type of ONT Application | [optional] [default to ""]

## Methods

### NewBundlesPatchRequestEndpointBundleValueRgServicesInner

`func NewBundlesPatchRequestEndpointBundleValueRgServicesInner() *BundlesPatchRequestEndpointBundleValueRgServicesInner`

NewBundlesPatchRequestEndpointBundleValueRgServicesInner instantiates a new BundlesPatchRequestEndpointBundleValueRgServicesInner object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBundlesPatchRequestEndpointBundleValueRgServicesInnerWithDefaults

`func NewBundlesPatchRequestEndpointBundleValueRgServicesInnerWithDefaults() *BundlesPatchRequestEndpointBundleValueRgServicesInner`

NewBundlesPatchRequestEndpointBundleValueRgServicesInnerWithDefaults instantiates a new BundlesPatchRequestEndpointBundleValueRgServicesInner object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetIpMask

`func (o *BundlesPatchRequestEndpointBundleValueRgServicesInner) GetIpMask() string`

GetIpMask returns the IpMask field if non-nil, zero value otherwise.

### GetIpMaskOk

`func (o *BundlesPatchRequestEndpointBundleValueRgServicesInner) GetIpMaskOk() (*string, bool)`

GetIpMaskOk returns a tuple with the IpMask field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpMask

`func (o *BundlesPatchRequestEndpointBundleValueRgServicesInner) SetIpMask(v string)`

SetIpMask sets IpMask field to given value.

### HasIpMask

`func (o *BundlesPatchRequestEndpointBundleValueRgServicesInner) HasIpMask() bool`

HasIpMask returns a boolean if a field has been set.

### GetIndex

`func (o *BundlesPatchRequestEndpointBundleValueRgServicesInner) GetIndex() int32`

GetIndex returns the Index field if non-nil, zero value otherwise.

### GetIndexOk

`func (o *BundlesPatchRequestEndpointBundleValueRgServicesInner) GetIndexOk() (*int32, bool)`

GetIndexOk returns a tuple with the Index field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIndex

`func (o *BundlesPatchRequestEndpointBundleValueRgServicesInner) SetIndex(v int32)`

SetIndex sets Index field to given value.

### HasIndex

`func (o *BundlesPatchRequestEndpointBundleValueRgServicesInner) HasIndex() bool`

HasIndex returns a boolean if a field has been set.

### GetRowAppEnable

`func (o *BundlesPatchRequestEndpointBundleValueRgServicesInner) GetRowAppEnable() bool`

GetRowAppEnable returns the RowAppEnable field if non-nil, zero value otherwise.

### GetRowAppEnableOk

`func (o *BundlesPatchRequestEndpointBundleValueRgServicesInner) GetRowAppEnableOk() (*bool, bool)`

GetRowAppEnableOk returns a tuple with the RowAppEnable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowAppEnable

`func (o *BundlesPatchRequestEndpointBundleValueRgServicesInner) SetRowAppEnable(v bool)`

SetRowAppEnable sets RowAppEnable field to given value.

### HasRowAppEnable

`func (o *BundlesPatchRequestEndpointBundleValueRgServicesInner) HasRowAppEnable() bool`

HasRowAppEnable returns a boolean if a field has been set.

### GetRowAppConnectedService

`func (o *BundlesPatchRequestEndpointBundleValueRgServicesInner) GetRowAppConnectedService() string`

GetRowAppConnectedService returns the RowAppConnectedService field if non-nil, zero value otherwise.

### GetRowAppConnectedServiceOk

`func (o *BundlesPatchRequestEndpointBundleValueRgServicesInner) GetRowAppConnectedServiceOk() (*string, bool)`

GetRowAppConnectedServiceOk returns a tuple with the RowAppConnectedService field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowAppConnectedService

`func (o *BundlesPatchRequestEndpointBundleValueRgServicesInner) SetRowAppConnectedService(v string)`

SetRowAppConnectedService sets RowAppConnectedService field to given value.

### HasRowAppConnectedService

`func (o *BundlesPatchRequestEndpointBundleValueRgServicesInner) HasRowAppConnectedService() bool`

HasRowAppConnectedService returns a boolean if a field has been set.

### GetRowAppConnectedServiceRefType

`func (o *BundlesPatchRequestEndpointBundleValueRgServicesInner) GetRowAppConnectedServiceRefType() string`

GetRowAppConnectedServiceRefType returns the RowAppConnectedServiceRefType field if non-nil, zero value otherwise.

### GetRowAppConnectedServiceRefTypeOk

`func (o *BundlesPatchRequestEndpointBundleValueRgServicesInner) GetRowAppConnectedServiceRefTypeOk() (*string, bool)`

GetRowAppConnectedServiceRefTypeOk returns a tuple with the RowAppConnectedServiceRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowAppConnectedServiceRefType

`func (o *BundlesPatchRequestEndpointBundleValueRgServicesInner) SetRowAppConnectedServiceRefType(v string)`

SetRowAppConnectedServiceRefType sets RowAppConnectedServiceRefType field to given value.

### HasRowAppConnectedServiceRefType

`func (o *BundlesPatchRequestEndpointBundleValueRgServicesInner) HasRowAppConnectedServiceRefType() bool`

HasRowAppConnectedServiceRefType returns a boolean if a field has been set.

### GetRowAppType

`func (o *BundlesPatchRequestEndpointBundleValueRgServicesInner) GetRowAppType() string`

GetRowAppType returns the RowAppType field if non-nil, zero value otherwise.

### GetRowAppTypeOk

`func (o *BundlesPatchRequestEndpointBundleValueRgServicesInner) GetRowAppTypeOk() (*string, bool)`

GetRowAppTypeOk returns a tuple with the RowAppType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowAppType

`func (o *BundlesPatchRequestEndpointBundleValueRgServicesInner) SetRowAppType(v string)`

SetRowAppType sets RowAppType field to given value.

### HasRowAppType

`func (o *BundlesPatchRequestEndpointBundleValueRgServicesInner) HasRowAppType() bool`

HasRowAppType returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


