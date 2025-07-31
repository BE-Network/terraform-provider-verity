# SwitchpointsPutRequestEndpointValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Object Name. Must be unique. | [optional] [default to ""]
**DeviceSerialNumber** | Pointer to **string** | Device Serial Number | [optional] [default to ""]
**ConnectedBundle** | Pointer to **string** | Connected Bundle | [optional] [default to ""]
**ConnectedBundleRefType** | Pointer to **string** | Object type for connected_bundle field | [optional] 
**ReadOnlyMode** | Pointer to **bool** | When Read Only Mode is checked, vNetC will perform all functions except writing database updates to the target hardware | [optional] [default to false]
**Locked** | Pointer to **bool** | Permission lock | [optional] [default to false]
**DisabledPorts** | Pointer to **string** | Disabled Ports It&#39;s a comma separated list of ports to disable. | [optional] [default to ""]
**IsFabric** | Pointer to **bool** | For Switch Endpoints. Denotes a Switch that is Fabric rather than an Edge Device | [optional] [default to false]
**OutOfBandManagement** | Pointer to **bool** | For Switch Endpoints. Denotes a Switch is managed out of band via the management port | [optional] [default to false]
**Badges** | Pointer to [**[]SwitchpointsPutRequestSwitchpointValueBadgesInner**](SwitchpointsPutRequestSwitchpointValueBadgesInner.md) |  | [optional] 
**Children** | Pointer to [**[]SwitchpointsPutRequestEndpointValueChildrenInner**](SwitchpointsPutRequestEndpointValueChildrenInner.md) |  | [optional] 
**TrafficMirrors** | Pointer to [**[]SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner**](SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner.md) |  | [optional] 
**Eths** | Pointer to [**[]SwitchpointsPutRequestSwitchpointValueEthsInner**](SwitchpointsPutRequestSwitchpointValueEthsInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**SwitchpointsPutRequestEndpointValueObjectProperties**](SwitchpointsPutRequestEndpointValueObjectProperties.md) |  | [optional] 

## Methods

### NewSwitchpointsPutRequestEndpointValue

`func NewSwitchpointsPutRequestEndpointValue() *SwitchpointsPutRequestEndpointValue`

NewSwitchpointsPutRequestEndpointValue instantiates a new SwitchpointsPutRequestEndpointValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSwitchpointsPutRequestEndpointValueWithDefaults

`func NewSwitchpointsPutRequestEndpointValueWithDefaults() *SwitchpointsPutRequestEndpointValue`

NewSwitchpointsPutRequestEndpointValueWithDefaults instantiates a new SwitchpointsPutRequestEndpointValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *SwitchpointsPutRequestEndpointValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *SwitchpointsPutRequestEndpointValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *SwitchpointsPutRequestEndpointValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *SwitchpointsPutRequestEndpointValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetDeviceSerialNumber

`func (o *SwitchpointsPutRequestEndpointValue) GetDeviceSerialNumber() string`

GetDeviceSerialNumber returns the DeviceSerialNumber field if non-nil, zero value otherwise.

### GetDeviceSerialNumberOk

`func (o *SwitchpointsPutRequestEndpointValue) GetDeviceSerialNumberOk() (*string, bool)`

GetDeviceSerialNumberOk returns a tuple with the DeviceSerialNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceSerialNumber

`func (o *SwitchpointsPutRequestEndpointValue) SetDeviceSerialNumber(v string)`

SetDeviceSerialNumber sets DeviceSerialNumber field to given value.

### HasDeviceSerialNumber

`func (o *SwitchpointsPutRequestEndpointValue) HasDeviceSerialNumber() bool`

HasDeviceSerialNumber returns a boolean if a field has been set.

### GetConnectedBundle

`func (o *SwitchpointsPutRequestEndpointValue) GetConnectedBundle() string`

GetConnectedBundle returns the ConnectedBundle field if non-nil, zero value otherwise.

### GetConnectedBundleOk

`func (o *SwitchpointsPutRequestEndpointValue) GetConnectedBundleOk() (*string, bool)`

GetConnectedBundleOk returns a tuple with the ConnectedBundle field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnectedBundle

`func (o *SwitchpointsPutRequestEndpointValue) SetConnectedBundle(v string)`

SetConnectedBundle sets ConnectedBundle field to given value.

### HasConnectedBundle

`func (o *SwitchpointsPutRequestEndpointValue) HasConnectedBundle() bool`

HasConnectedBundle returns a boolean if a field has been set.

### GetConnectedBundleRefType

`func (o *SwitchpointsPutRequestEndpointValue) GetConnectedBundleRefType() string`

GetConnectedBundleRefType returns the ConnectedBundleRefType field if non-nil, zero value otherwise.

### GetConnectedBundleRefTypeOk

`func (o *SwitchpointsPutRequestEndpointValue) GetConnectedBundleRefTypeOk() (*string, bool)`

GetConnectedBundleRefTypeOk returns a tuple with the ConnectedBundleRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnectedBundleRefType

`func (o *SwitchpointsPutRequestEndpointValue) SetConnectedBundleRefType(v string)`

SetConnectedBundleRefType sets ConnectedBundleRefType field to given value.

### HasConnectedBundleRefType

`func (o *SwitchpointsPutRequestEndpointValue) HasConnectedBundleRefType() bool`

HasConnectedBundleRefType returns a boolean if a field has been set.

### GetReadOnlyMode

`func (o *SwitchpointsPutRequestEndpointValue) GetReadOnlyMode() bool`

GetReadOnlyMode returns the ReadOnlyMode field if non-nil, zero value otherwise.

### GetReadOnlyModeOk

`func (o *SwitchpointsPutRequestEndpointValue) GetReadOnlyModeOk() (*bool, bool)`

GetReadOnlyModeOk returns a tuple with the ReadOnlyMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReadOnlyMode

`func (o *SwitchpointsPutRequestEndpointValue) SetReadOnlyMode(v bool)`

SetReadOnlyMode sets ReadOnlyMode field to given value.

### HasReadOnlyMode

`func (o *SwitchpointsPutRequestEndpointValue) HasReadOnlyMode() bool`

HasReadOnlyMode returns a boolean if a field has been set.

### GetLocked

`func (o *SwitchpointsPutRequestEndpointValue) GetLocked() bool`

GetLocked returns the Locked field if non-nil, zero value otherwise.

### GetLockedOk

`func (o *SwitchpointsPutRequestEndpointValue) GetLockedOk() (*bool, bool)`

GetLockedOk returns a tuple with the Locked field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocked

`func (o *SwitchpointsPutRequestEndpointValue) SetLocked(v bool)`

SetLocked sets Locked field to given value.

### HasLocked

`func (o *SwitchpointsPutRequestEndpointValue) HasLocked() bool`

HasLocked returns a boolean if a field has been set.

### GetDisabledPorts

`func (o *SwitchpointsPutRequestEndpointValue) GetDisabledPorts() string`

GetDisabledPorts returns the DisabledPorts field if non-nil, zero value otherwise.

### GetDisabledPortsOk

`func (o *SwitchpointsPutRequestEndpointValue) GetDisabledPortsOk() (*string, bool)`

GetDisabledPortsOk returns a tuple with the DisabledPorts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDisabledPorts

`func (o *SwitchpointsPutRequestEndpointValue) SetDisabledPorts(v string)`

SetDisabledPorts sets DisabledPorts field to given value.

### HasDisabledPorts

`func (o *SwitchpointsPutRequestEndpointValue) HasDisabledPorts() bool`

HasDisabledPorts returns a boolean if a field has been set.

### GetIsFabric

`func (o *SwitchpointsPutRequestEndpointValue) GetIsFabric() bool`

GetIsFabric returns the IsFabric field if non-nil, zero value otherwise.

### GetIsFabricOk

`func (o *SwitchpointsPutRequestEndpointValue) GetIsFabricOk() (*bool, bool)`

GetIsFabricOk returns a tuple with the IsFabric field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsFabric

`func (o *SwitchpointsPutRequestEndpointValue) SetIsFabric(v bool)`

SetIsFabric sets IsFabric field to given value.

### HasIsFabric

`func (o *SwitchpointsPutRequestEndpointValue) HasIsFabric() bool`

HasIsFabric returns a boolean if a field has been set.

### GetOutOfBandManagement

`func (o *SwitchpointsPutRequestEndpointValue) GetOutOfBandManagement() bool`

GetOutOfBandManagement returns the OutOfBandManagement field if non-nil, zero value otherwise.

### GetOutOfBandManagementOk

`func (o *SwitchpointsPutRequestEndpointValue) GetOutOfBandManagementOk() (*bool, bool)`

GetOutOfBandManagementOk returns a tuple with the OutOfBandManagement field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutOfBandManagement

`func (o *SwitchpointsPutRequestEndpointValue) SetOutOfBandManagement(v bool)`

SetOutOfBandManagement sets OutOfBandManagement field to given value.

### HasOutOfBandManagement

`func (o *SwitchpointsPutRequestEndpointValue) HasOutOfBandManagement() bool`

HasOutOfBandManagement returns a boolean if a field has been set.

### GetBadges

`func (o *SwitchpointsPutRequestEndpointValue) GetBadges() []SwitchpointsPutRequestSwitchpointValueBadgesInner`

GetBadges returns the Badges field if non-nil, zero value otherwise.

### GetBadgesOk

`func (o *SwitchpointsPutRequestEndpointValue) GetBadgesOk() (*[]SwitchpointsPutRequestSwitchpointValueBadgesInner, bool)`

GetBadgesOk returns a tuple with the Badges field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBadges

`func (o *SwitchpointsPutRequestEndpointValue) SetBadges(v []SwitchpointsPutRequestSwitchpointValueBadgesInner)`

SetBadges sets Badges field to given value.

### HasBadges

`func (o *SwitchpointsPutRequestEndpointValue) HasBadges() bool`

HasBadges returns a boolean if a field has been set.

### GetChildren

`func (o *SwitchpointsPutRequestEndpointValue) GetChildren() []SwitchpointsPutRequestEndpointValueChildrenInner`

GetChildren returns the Children field if non-nil, zero value otherwise.

### GetChildrenOk

`func (o *SwitchpointsPutRequestEndpointValue) GetChildrenOk() (*[]SwitchpointsPutRequestEndpointValueChildrenInner, bool)`

GetChildrenOk returns a tuple with the Children field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChildren

`func (o *SwitchpointsPutRequestEndpointValue) SetChildren(v []SwitchpointsPutRequestEndpointValueChildrenInner)`

SetChildren sets Children field to given value.

### HasChildren

`func (o *SwitchpointsPutRequestEndpointValue) HasChildren() bool`

HasChildren returns a boolean if a field has been set.

### GetTrafficMirrors

`func (o *SwitchpointsPutRequestEndpointValue) GetTrafficMirrors() []SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner`

GetTrafficMirrors returns the TrafficMirrors field if non-nil, zero value otherwise.

### GetTrafficMirrorsOk

`func (o *SwitchpointsPutRequestEndpointValue) GetTrafficMirrorsOk() (*[]SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner, bool)`

GetTrafficMirrorsOk returns a tuple with the TrafficMirrors field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTrafficMirrors

`func (o *SwitchpointsPutRequestEndpointValue) SetTrafficMirrors(v []SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner)`

SetTrafficMirrors sets TrafficMirrors field to given value.

### HasTrafficMirrors

`func (o *SwitchpointsPutRequestEndpointValue) HasTrafficMirrors() bool`

HasTrafficMirrors returns a boolean if a field has been set.

### GetEths

`func (o *SwitchpointsPutRequestEndpointValue) GetEths() []SwitchpointsPutRequestSwitchpointValueEthsInner`

GetEths returns the Eths field if non-nil, zero value otherwise.

### GetEthsOk

`func (o *SwitchpointsPutRequestEndpointValue) GetEthsOk() (*[]SwitchpointsPutRequestSwitchpointValueEthsInner, bool)`

GetEthsOk returns a tuple with the Eths field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEths

`func (o *SwitchpointsPutRequestEndpointValue) SetEths(v []SwitchpointsPutRequestSwitchpointValueEthsInner)`

SetEths sets Eths field to given value.

### HasEths

`func (o *SwitchpointsPutRequestEndpointValue) HasEths() bool`

HasEths returns a boolean if a field has been set.

### GetObjectProperties

`func (o *SwitchpointsPutRequestEndpointValue) GetObjectProperties() SwitchpointsPutRequestEndpointValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *SwitchpointsPutRequestEndpointValue) GetObjectPropertiesOk() (*SwitchpointsPutRequestEndpointValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *SwitchpointsPutRequestEndpointValue) SetObjectProperties(v SwitchpointsPutRequestEndpointValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *SwitchpointsPutRequestEndpointValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


