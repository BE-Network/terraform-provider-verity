# ConfigPutRequestEndpointEndpointName

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
**Badges** | Pointer to [**[]ConfigPutRequestSwitchpointSwitchpointNameBadgesInner**](ConfigPutRequestSwitchpointSwitchpointNameBadgesInner.md) |  | [optional] 
**Children** | Pointer to [**[]ConfigPutRequestEndpointEndpointNameChildrenInner**](ConfigPutRequestEndpointEndpointNameChildrenInner.md) |  | [optional] 
**TrafficMirrors** | Pointer to [**[]ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner**](ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner.md) |  | [optional] 
**Eths** | Pointer to [**[]ConfigPutRequestEndpointEndpointNameEthsInner**](ConfigPutRequestEndpointEndpointNameEthsInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**ConfigPutRequestEndpointEndpointNameObjectProperties**](ConfigPutRequestEndpointEndpointNameObjectProperties.md) |  | [optional] 

## Methods

### NewConfigPutRequestEndpointEndpointName

`func NewConfigPutRequestEndpointEndpointName() *ConfigPutRequestEndpointEndpointName`

NewConfigPutRequestEndpointEndpointName instantiates a new ConfigPutRequestEndpointEndpointName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestEndpointEndpointNameWithDefaults

`func NewConfigPutRequestEndpointEndpointNameWithDefaults() *ConfigPutRequestEndpointEndpointName`

NewConfigPutRequestEndpointEndpointNameWithDefaults instantiates a new ConfigPutRequestEndpointEndpointName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestEndpointEndpointName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestEndpointEndpointName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestEndpointEndpointName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestEndpointEndpointName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetDeviceSerialNumber

`func (o *ConfigPutRequestEndpointEndpointName) GetDeviceSerialNumber() string`

GetDeviceSerialNumber returns the DeviceSerialNumber field if non-nil, zero value otherwise.

### GetDeviceSerialNumberOk

`func (o *ConfigPutRequestEndpointEndpointName) GetDeviceSerialNumberOk() (*string, bool)`

GetDeviceSerialNumberOk returns a tuple with the DeviceSerialNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceSerialNumber

`func (o *ConfigPutRequestEndpointEndpointName) SetDeviceSerialNumber(v string)`

SetDeviceSerialNumber sets DeviceSerialNumber field to given value.

### HasDeviceSerialNumber

`func (o *ConfigPutRequestEndpointEndpointName) HasDeviceSerialNumber() bool`

HasDeviceSerialNumber returns a boolean if a field has been set.

### GetConnectedBundle

`func (o *ConfigPutRequestEndpointEndpointName) GetConnectedBundle() string`

GetConnectedBundle returns the ConnectedBundle field if non-nil, zero value otherwise.

### GetConnectedBundleOk

`func (o *ConfigPutRequestEndpointEndpointName) GetConnectedBundleOk() (*string, bool)`

GetConnectedBundleOk returns a tuple with the ConnectedBundle field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnectedBundle

`func (o *ConfigPutRequestEndpointEndpointName) SetConnectedBundle(v string)`

SetConnectedBundle sets ConnectedBundle field to given value.

### HasConnectedBundle

`func (o *ConfigPutRequestEndpointEndpointName) HasConnectedBundle() bool`

HasConnectedBundle returns a boolean if a field has been set.

### GetConnectedBundleRefType

`func (o *ConfigPutRequestEndpointEndpointName) GetConnectedBundleRefType() string`

GetConnectedBundleRefType returns the ConnectedBundleRefType field if non-nil, zero value otherwise.

### GetConnectedBundleRefTypeOk

`func (o *ConfigPutRequestEndpointEndpointName) GetConnectedBundleRefTypeOk() (*string, bool)`

GetConnectedBundleRefTypeOk returns a tuple with the ConnectedBundleRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnectedBundleRefType

`func (o *ConfigPutRequestEndpointEndpointName) SetConnectedBundleRefType(v string)`

SetConnectedBundleRefType sets ConnectedBundleRefType field to given value.

### HasConnectedBundleRefType

`func (o *ConfigPutRequestEndpointEndpointName) HasConnectedBundleRefType() bool`

HasConnectedBundleRefType returns a boolean if a field has been set.

### GetReadOnlyMode

`func (o *ConfigPutRequestEndpointEndpointName) GetReadOnlyMode() bool`

GetReadOnlyMode returns the ReadOnlyMode field if non-nil, zero value otherwise.

### GetReadOnlyModeOk

`func (o *ConfigPutRequestEndpointEndpointName) GetReadOnlyModeOk() (*bool, bool)`

GetReadOnlyModeOk returns a tuple with the ReadOnlyMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReadOnlyMode

`func (o *ConfigPutRequestEndpointEndpointName) SetReadOnlyMode(v bool)`

SetReadOnlyMode sets ReadOnlyMode field to given value.

### HasReadOnlyMode

`func (o *ConfigPutRequestEndpointEndpointName) HasReadOnlyMode() bool`

HasReadOnlyMode returns a boolean if a field has been set.

### GetLocked

`func (o *ConfigPutRequestEndpointEndpointName) GetLocked() bool`

GetLocked returns the Locked field if non-nil, zero value otherwise.

### GetLockedOk

`func (o *ConfigPutRequestEndpointEndpointName) GetLockedOk() (*bool, bool)`

GetLockedOk returns a tuple with the Locked field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocked

`func (o *ConfigPutRequestEndpointEndpointName) SetLocked(v bool)`

SetLocked sets Locked field to given value.

### HasLocked

`func (o *ConfigPutRequestEndpointEndpointName) HasLocked() bool`

HasLocked returns a boolean if a field has been set.

### GetDisabledPorts

`func (o *ConfigPutRequestEndpointEndpointName) GetDisabledPorts() string`

GetDisabledPorts returns the DisabledPorts field if non-nil, zero value otherwise.

### GetDisabledPortsOk

`func (o *ConfigPutRequestEndpointEndpointName) GetDisabledPortsOk() (*string, bool)`

GetDisabledPortsOk returns a tuple with the DisabledPorts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDisabledPorts

`func (o *ConfigPutRequestEndpointEndpointName) SetDisabledPorts(v string)`

SetDisabledPorts sets DisabledPorts field to given value.

### HasDisabledPorts

`func (o *ConfigPutRequestEndpointEndpointName) HasDisabledPorts() bool`

HasDisabledPorts returns a boolean if a field has been set.

### GetIsFabric

`func (o *ConfigPutRequestEndpointEndpointName) GetIsFabric() bool`

GetIsFabric returns the IsFabric field if non-nil, zero value otherwise.

### GetIsFabricOk

`func (o *ConfigPutRequestEndpointEndpointName) GetIsFabricOk() (*bool, bool)`

GetIsFabricOk returns a tuple with the IsFabric field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsFabric

`func (o *ConfigPutRequestEndpointEndpointName) SetIsFabric(v bool)`

SetIsFabric sets IsFabric field to given value.

### HasIsFabric

`func (o *ConfigPutRequestEndpointEndpointName) HasIsFabric() bool`

HasIsFabric returns a boolean if a field has been set.

### GetOutOfBandManagement

`func (o *ConfigPutRequestEndpointEndpointName) GetOutOfBandManagement() bool`

GetOutOfBandManagement returns the OutOfBandManagement field if non-nil, zero value otherwise.

### GetOutOfBandManagementOk

`func (o *ConfigPutRequestEndpointEndpointName) GetOutOfBandManagementOk() (*bool, bool)`

GetOutOfBandManagementOk returns a tuple with the OutOfBandManagement field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutOfBandManagement

`func (o *ConfigPutRequestEndpointEndpointName) SetOutOfBandManagement(v bool)`

SetOutOfBandManagement sets OutOfBandManagement field to given value.

### HasOutOfBandManagement

`func (o *ConfigPutRequestEndpointEndpointName) HasOutOfBandManagement() bool`

HasOutOfBandManagement returns a boolean if a field has been set.

### GetBadges

`func (o *ConfigPutRequestEndpointEndpointName) GetBadges() []ConfigPutRequestSwitchpointSwitchpointNameBadgesInner`

GetBadges returns the Badges field if non-nil, zero value otherwise.

### GetBadgesOk

`func (o *ConfigPutRequestEndpointEndpointName) GetBadgesOk() (*[]ConfigPutRequestSwitchpointSwitchpointNameBadgesInner, bool)`

GetBadgesOk returns a tuple with the Badges field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBadges

`func (o *ConfigPutRequestEndpointEndpointName) SetBadges(v []ConfigPutRequestSwitchpointSwitchpointNameBadgesInner)`

SetBadges sets Badges field to given value.

### HasBadges

`func (o *ConfigPutRequestEndpointEndpointName) HasBadges() bool`

HasBadges returns a boolean if a field has been set.

### GetChildren

`func (o *ConfigPutRequestEndpointEndpointName) GetChildren() []ConfigPutRequestEndpointEndpointNameChildrenInner`

GetChildren returns the Children field if non-nil, zero value otherwise.

### GetChildrenOk

`func (o *ConfigPutRequestEndpointEndpointName) GetChildrenOk() (*[]ConfigPutRequestEndpointEndpointNameChildrenInner, bool)`

GetChildrenOk returns a tuple with the Children field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChildren

`func (o *ConfigPutRequestEndpointEndpointName) SetChildren(v []ConfigPutRequestEndpointEndpointNameChildrenInner)`

SetChildren sets Children field to given value.

### HasChildren

`func (o *ConfigPutRequestEndpointEndpointName) HasChildren() bool`

HasChildren returns a boolean if a field has been set.

### GetTrafficMirrors

`func (o *ConfigPutRequestEndpointEndpointName) GetTrafficMirrors() []ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner`

GetTrafficMirrors returns the TrafficMirrors field if non-nil, zero value otherwise.

### GetTrafficMirrorsOk

`func (o *ConfigPutRequestEndpointEndpointName) GetTrafficMirrorsOk() (*[]ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner, bool)`

GetTrafficMirrorsOk returns a tuple with the TrafficMirrors field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTrafficMirrors

`func (o *ConfigPutRequestEndpointEndpointName) SetTrafficMirrors(v []ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner)`

SetTrafficMirrors sets TrafficMirrors field to given value.

### HasTrafficMirrors

`func (o *ConfigPutRequestEndpointEndpointName) HasTrafficMirrors() bool`

HasTrafficMirrors returns a boolean if a field has been set.

### GetEths

`func (o *ConfigPutRequestEndpointEndpointName) GetEths() []ConfigPutRequestEndpointEndpointNameEthsInner`

GetEths returns the Eths field if non-nil, zero value otherwise.

### GetEthsOk

`func (o *ConfigPutRequestEndpointEndpointName) GetEthsOk() (*[]ConfigPutRequestEndpointEndpointNameEthsInner, bool)`

GetEthsOk returns a tuple with the Eths field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEths

`func (o *ConfigPutRequestEndpointEndpointName) SetEths(v []ConfigPutRequestEndpointEndpointNameEthsInner)`

SetEths sets Eths field to given value.

### HasEths

`func (o *ConfigPutRequestEndpointEndpointName) HasEths() bool`

HasEths returns a boolean if a field has been set.

### GetObjectProperties

`func (o *ConfigPutRequestEndpointEndpointName) GetObjectProperties() ConfigPutRequestEndpointEndpointNameObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ConfigPutRequestEndpointEndpointName) GetObjectPropertiesOk() (*ConfigPutRequestEndpointEndpointNameObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ConfigPutRequestEndpointEndpointName) SetObjectProperties(v ConfigPutRequestEndpointEndpointNameObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ConfigPutRequestEndpointEndpointName) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


