# ConfigPutRequestSwitchpointSwitchpointName

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
**OutOfBandManagement** | Pointer to **bool** | For Switch Endpoints. Denotes a Switch is managed out of band via the management port | [optional] [default to false]
**Type** | Pointer to **string** | Type of Switchpoint | [optional] [default to "leaf"]
**SuperPod** | Pointer to **string** | Super Pod  subgrouping of super spines and pods | [optional] [default to ""]
**Pod** | Pointer to **string** | Pod  subgrouping of spine and leaf switches  | [optional] [default to ""]
**Rack** | Pointer to **string** | Physical Rack location of the Switch  | [optional] [default to ""]
**SwitchRouterIdIpMask** | Pointer to **string** | Switch BGP Router Identifier | [optional] [default to "(auto)"]
**SwitchRouterIdIpMaskAutoAssigned** | Pointer to **bool** | Whether or not the value in switch_router_id_ip_mask field has been automatically assigned or not. Set to false and change switch_router_id_ip_mask value to edit. | [optional] 
**SwitchVtepIdIpMask** | Pointer to **string** | Switch VETP Identifier | [optional] [default to "(auto)"]
**SwitchVtepIdIpMaskAutoAssigned** | Pointer to **bool** | Whether or not the value in switch_vtep_id_ip_mask field has been automatically assigned or not. Set to false and change switch_vtep_id_ip_mask value to edit. | [optional] 
**BgpAsNumber** | Pointer to **int32** | BGP Autonomous System Number for the site underlay  | [optional] 
**BgpAsNumberAutoAssigned** | Pointer to **bool** | Whether or not the value in bgp_as_number field has been automatically assigned or not. Set to false and change bgp_as_number value to edit. | [optional] 
**Badges** | Pointer to [**[]ConfigPutRequestSwitchpointSwitchpointNameBadgesInner**](ConfigPutRequestSwitchpointSwitchpointNameBadgesInner.md) |  | [optional] 
**Children** | Pointer to [**[]ConfigPutRequestSwitchpointSwitchpointNameChildrenInner**](ConfigPutRequestSwitchpointSwitchpointNameChildrenInner.md) |  | [optional] 
**TrafficMirrors** | Pointer to [**[]ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner**](ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner.md) |  | [optional] 
**Eths** | Pointer to [**[]ConfigPutRequestSwitchpointSwitchpointNameEthsInner**](ConfigPutRequestSwitchpointSwitchpointNameEthsInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**ConfigPutRequestSwitchpointSwitchpointNameObjectProperties**](ConfigPutRequestSwitchpointSwitchpointNameObjectProperties.md) |  | [optional] 

## Methods

### NewConfigPutRequestSwitchpointSwitchpointName

`func NewConfigPutRequestSwitchpointSwitchpointName() *ConfigPutRequestSwitchpointSwitchpointName`

NewConfigPutRequestSwitchpointSwitchpointName instantiates a new ConfigPutRequestSwitchpointSwitchpointName object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestSwitchpointSwitchpointNameWithDefaults

`func NewConfigPutRequestSwitchpointSwitchpointNameWithDefaults() *ConfigPutRequestSwitchpointSwitchpointName`

NewConfigPutRequestSwitchpointSwitchpointNameWithDefaults instantiates a new ConfigPutRequestSwitchpointSwitchpointName object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *ConfigPutRequestSwitchpointSwitchpointName) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *ConfigPutRequestSwitchpointSwitchpointName) HasName() bool`

HasName returns a boolean if a field has been set.

### GetDeviceSerialNumber

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetDeviceSerialNumber() string`

GetDeviceSerialNumber returns the DeviceSerialNumber field if non-nil, zero value otherwise.

### GetDeviceSerialNumberOk

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetDeviceSerialNumberOk() (*string, bool)`

GetDeviceSerialNumberOk returns a tuple with the DeviceSerialNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceSerialNumber

`func (o *ConfigPutRequestSwitchpointSwitchpointName) SetDeviceSerialNumber(v string)`

SetDeviceSerialNumber sets DeviceSerialNumber field to given value.

### HasDeviceSerialNumber

`func (o *ConfigPutRequestSwitchpointSwitchpointName) HasDeviceSerialNumber() bool`

HasDeviceSerialNumber returns a boolean if a field has been set.

### GetConnectedBundle

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetConnectedBundle() string`

GetConnectedBundle returns the ConnectedBundle field if non-nil, zero value otherwise.

### GetConnectedBundleOk

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetConnectedBundleOk() (*string, bool)`

GetConnectedBundleOk returns a tuple with the ConnectedBundle field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnectedBundle

`func (o *ConfigPutRequestSwitchpointSwitchpointName) SetConnectedBundle(v string)`

SetConnectedBundle sets ConnectedBundle field to given value.

### HasConnectedBundle

`func (o *ConfigPutRequestSwitchpointSwitchpointName) HasConnectedBundle() bool`

HasConnectedBundle returns a boolean if a field has been set.

### GetConnectedBundleRefType

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetConnectedBundleRefType() string`

GetConnectedBundleRefType returns the ConnectedBundleRefType field if non-nil, zero value otherwise.

### GetConnectedBundleRefTypeOk

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetConnectedBundleRefTypeOk() (*string, bool)`

GetConnectedBundleRefTypeOk returns a tuple with the ConnectedBundleRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnectedBundleRefType

`func (o *ConfigPutRequestSwitchpointSwitchpointName) SetConnectedBundleRefType(v string)`

SetConnectedBundleRefType sets ConnectedBundleRefType field to given value.

### HasConnectedBundleRefType

`func (o *ConfigPutRequestSwitchpointSwitchpointName) HasConnectedBundleRefType() bool`

HasConnectedBundleRefType returns a boolean if a field has been set.

### GetReadOnlyMode

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetReadOnlyMode() bool`

GetReadOnlyMode returns the ReadOnlyMode field if non-nil, zero value otherwise.

### GetReadOnlyModeOk

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetReadOnlyModeOk() (*bool, bool)`

GetReadOnlyModeOk returns a tuple with the ReadOnlyMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReadOnlyMode

`func (o *ConfigPutRequestSwitchpointSwitchpointName) SetReadOnlyMode(v bool)`

SetReadOnlyMode sets ReadOnlyMode field to given value.

### HasReadOnlyMode

`func (o *ConfigPutRequestSwitchpointSwitchpointName) HasReadOnlyMode() bool`

HasReadOnlyMode returns a boolean if a field has been set.

### GetLocked

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetLocked() bool`

GetLocked returns the Locked field if non-nil, zero value otherwise.

### GetLockedOk

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetLockedOk() (*bool, bool)`

GetLockedOk returns a tuple with the Locked field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocked

`func (o *ConfigPutRequestSwitchpointSwitchpointName) SetLocked(v bool)`

SetLocked sets Locked field to given value.

### HasLocked

`func (o *ConfigPutRequestSwitchpointSwitchpointName) HasLocked() bool`

HasLocked returns a boolean if a field has been set.

### GetDisabledPorts

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetDisabledPorts() string`

GetDisabledPorts returns the DisabledPorts field if non-nil, zero value otherwise.

### GetDisabledPortsOk

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetDisabledPortsOk() (*string, bool)`

GetDisabledPortsOk returns a tuple with the DisabledPorts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDisabledPorts

`func (o *ConfigPutRequestSwitchpointSwitchpointName) SetDisabledPorts(v string)`

SetDisabledPorts sets DisabledPorts field to given value.

### HasDisabledPorts

`func (o *ConfigPutRequestSwitchpointSwitchpointName) HasDisabledPorts() bool`

HasDisabledPorts returns a boolean if a field has been set.

### GetOutOfBandManagement

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetOutOfBandManagement() bool`

GetOutOfBandManagement returns the OutOfBandManagement field if non-nil, zero value otherwise.

### GetOutOfBandManagementOk

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetOutOfBandManagementOk() (*bool, bool)`

GetOutOfBandManagementOk returns a tuple with the OutOfBandManagement field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutOfBandManagement

`func (o *ConfigPutRequestSwitchpointSwitchpointName) SetOutOfBandManagement(v bool)`

SetOutOfBandManagement sets OutOfBandManagement field to given value.

### HasOutOfBandManagement

`func (o *ConfigPutRequestSwitchpointSwitchpointName) HasOutOfBandManagement() bool`

HasOutOfBandManagement returns a boolean if a field has been set.

### GetType

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *ConfigPutRequestSwitchpointSwitchpointName) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *ConfigPutRequestSwitchpointSwitchpointName) HasType() bool`

HasType returns a boolean if a field has been set.

### GetSuperPod

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetSuperPod() string`

GetSuperPod returns the SuperPod field if non-nil, zero value otherwise.

### GetSuperPodOk

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetSuperPodOk() (*string, bool)`

GetSuperPodOk returns a tuple with the SuperPod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSuperPod

`func (o *ConfigPutRequestSwitchpointSwitchpointName) SetSuperPod(v string)`

SetSuperPod sets SuperPod field to given value.

### HasSuperPod

`func (o *ConfigPutRequestSwitchpointSwitchpointName) HasSuperPod() bool`

HasSuperPod returns a boolean if a field has been set.

### GetPod

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetPod() string`

GetPod returns the Pod field if non-nil, zero value otherwise.

### GetPodOk

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetPodOk() (*string, bool)`

GetPodOk returns a tuple with the Pod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPod

`func (o *ConfigPutRequestSwitchpointSwitchpointName) SetPod(v string)`

SetPod sets Pod field to given value.

### HasPod

`func (o *ConfigPutRequestSwitchpointSwitchpointName) HasPod() bool`

HasPod returns a boolean if a field has been set.

### GetRack

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetRack() string`

GetRack returns the Rack field if non-nil, zero value otherwise.

### GetRackOk

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetRackOk() (*string, bool)`

GetRackOk returns a tuple with the Rack field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRack

`func (o *ConfigPutRequestSwitchpointSwitchpointName) SetRack(v string)`

SetRack sets Rack field to given value.

### HasRack

`func (o *ConfigPutRequestSwitchpointSwitchpointName) HasRack() bool`

HasRack returns a boolean if a field has been set.

### GetSwitchRouterIdIpMask

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetSwitchRouterIdIpMask() string`

GetSwitchRouterIdIpMask returns the SwitchRouterIdIpMask field if non-nil, zero value otherwise.

### GetSwitchRouterIdIpMaskOk

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetSwitchRouterIdIpMaskOk() (*string, bool)`

GetSwitchRouterIdIpMaskOk returns a tuple with the SwitchRouterIdIpMask field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchRouterIdIpMask

`func (o *ConfigPutRequestSwitchpointSwitchpointName) SetSwitchRouterIdIpMask(v string)`

SetSwitchRouterIdIpMask sets SwitchRouterIdIpMask field to given value.

### HasSwitchRouterIdIpMask

`func (o *ConfigPutRequestSwitchpointSwitchpointName) HasSwitchRouterIdIpMask() bool`

HasSwitchRouterIdIpMask returns a boolean if a field has been set.

### GetSwitchRouterIdIpMaskAutoAssigned

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetSwitchRouterIdIpMaskAutoAssigned() bool`

GetSwitchRouterIdIpMaskAutoAssigned returns the SwitchRouterIdIpMaskAutoAssigned field if non-nil, zero value otherwise.

### GetSwitchRouterIdIpMaskAutoAssignedOk

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetSwitchRouterIdIpMaskAutoAssignedOk() (*bool, bool)`

GetSwitchRouterIdIpMaskAutoAssignedOk returns a tuple with the SwitchRouterIdIpMaskAutoAssigned field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchRouterIdIpMaskAutoAssigned

`func (o *ConfigPutRequestSwitchpointSwitchpointName) SetSwitchRouterIdIpMaskAutoAssigned(v bool)`

SetSwitchRouterIdIpMaskAutoAssigned sets SwitchRouterIdIpMaskAutoAssigned field to given value.

### HasSwitchRouterIdIpMaskAutoAssigned

`func (o *ConfigPutRequestSwitchpointSwitchpointName) HasSwitchRouterIdIpMaskAutoAssigned() bool`

HasSwitchRouterIdIpMaskAutoAssigned returns a boolean if a field has been set.

### GetSwitchVtepIdIpMask

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetSwitchVtepIdIpMask() string`

GetSwitchVtepIdIpMask returns the SwitchVtepIdIpMask field if non-nil, zero value otherwise.

### GetSwitchVtepIdIpMaskOk

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetSwitchVtepIdIpMaskOk() (*string, bool)`

GetSwitchVtepIdIpMaskOk returns a tuple with the SwitchVtepIdIpMask field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchVtepIdIpMask

`func (o *ConfigPutRequestSwitchpointSwitchpointName) SetSwitchVtepIdIpMask(v string)`

SetSwitchVtepIdIpMask sets SwitchVtepIdIpMask field to given value.

### HasSwitchVtepIdIpMask

`func (o *ConfigPutRequestSwitchpointSwitchpointName) HasSwitchVtepIdIpMask() bool`

HasSwitchVtepIdIpMask returns a boolean if a field has been set.

### GetSwitchVtepIdIpMaskAutoAssigned

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetSwitchVtepIdIpMaskAutoAssigned() bool`

GetSwitchVtepIdIpMaskAutoAssigned returns the SwitchVtepIdIpMaskAutoAssigned field if non-nil, zero value otherwise.

### GetSwitchVtepIdIpMaskAutoAssignedOk

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetSwitchVtepIdIpMaskAutoAssignedOk() (*bool, bool)`

GetSwitchVtepIdIpMaskAutoAssignedOk returns a tuple with the SwitchVtepIdIpMaskAutoAssigned field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchVtepIdIpMaskAutoAssigned

`func (o *ConfigPutRequestSwitchpointSwitchpointName) SetSwitchVtepIdIpMaskAutoAssigned(v bool)`

SetSwitchVtepIdIpMaskAutoAssigned sets SwitchVtepIdIpMaskAutoAssigned field to given value.

### HasSwitchVtepIdIpMaskAutoAssigned

`func (o *ConfigPutRequestSwitchpointSwitchpointName) HasSwitchVtepIdIpMaskAutoAssigned() bool`

HasSwitchVtepIdIpMaskAutoAssigned returns a boolean if a field has been set.

### GetBgpAsNumber

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetBgpAsNumber() int32`

GetBgpAsNumber returns the BgpAsNumber field if non-nil, zero value otherwise.

### GetBgpAsNumberOk

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetBgpAsNumberOk() (*int32, bool)`

GetBgpAsNumberOk returns a tuple with the BgpAsNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBgpAsNumber

`func (o *ConfigPutRequestSwitchpointSwitchpointName) SetBgpAsNumber(v int32)`

SetBgpAsNumber sets BgpAsNumber field to given value.

### HasBgpAsNumber

`func (o *ConfigPutRequestSwitchpointSwitchpointName) HasBgpAsNumber() bool`

HasBgpAsNumber returns a boolean if a field has been set.

### GetBgpAsNumberAutoAssigned

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetBgpAsNumberAutoAssigned() bool`

GetBgpAsNumberAutoAssigned returns the BgpAsNumberAutoAssigned field if non-nil, zero value otherwise.

### GetBgpAsNumberAutoAssignedOk

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetBgpAsNumberAutoAssignedOk() (*bool, bool)`

GetBgpAsNumberAutoAssignedOk returns a tuple with the BgpAsNumberAutoAssigned field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBgpAsNumberAutoAssigned

`func (o *ConfigPutRequestSwitchpointSwitchpointName) SetBgpAsNumberAutoAssigned(v bool)`

SetBgpAsNumberAutoAssigned sets BgpAsNumberAutoAssigned field to given value.

### HasBgpAsNumberAutoAssigned

`func (o *ConfigPutRequestSwitchpointSwitchpointName) HasBgpAsNumberAutoAssigned() bool`

HasBgpAsNumberAutoAssigned returns a boolean if a field has been set.

### GetBadges

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetBadges() []ConfigPutRequestSwitchpointSwitchpointNameBadgesInner`

GetBadges returns the Badges field if non-nil, zero value otherwise.

### GetBadgesOk

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetBadgesOk() (*[]ConfigPutRequestSwitchpointSwitchpointNameBadgesInner, bool)`

GetBadgesOk returns a tuple with the Badges field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBadges

`func (o *ConfigPutRequestSwitchpointSwitchpointName) SetBadges(v []ConfigPutRequestSwitchpointSwitchpointNameBadgesInner)`

SetBadges sets Badges field to given value.

### HasBadges

`func (o *ConfigPutRequestSwitchpointSwitchpointName) HasBadges() bool`

HasBadges returns a boolean if a field has been set.

### GetChildren

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetChildren() []ConfigPutRequestSwitchpointSwitchpointNameChildrenInner`

GetChildren returns the Children field if non-nil, zero value otherwise.

### GetChildrenOk

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetChildrenOk() (*[]ConfigPutRequestSwitchpointSwitchpointNameChildrenInner, bool)`

GetChildrenOk returns a tuple with the Children field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChildren

`func (o *ConfigPutRequestSwitchpointSwitchpointName) SetChildren(v []ConfigPutRequestSwitchpointSwitchpointNameChildrenInner)`

SetChildren sets Children field to given value.

### HasChildren

`func (o *ConfigPutRequestSwitchpointSwitchpointName) HasChildren() bool`

HasChildren returns a boolean if a field has been set.

### GetTrafficMirrors

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetTrafficMirrors() []ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner`

GetTrafficMirrors returns the TrafficMirrors field if non-nil, zero value otherwise.

### GetTrafficMirrorsOk

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetTrafficMirrorsOk() (*[]ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner, bool)`

GetTrafficMirrorsOk returns a tuple with the TrafficMirrors field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTrafficMirrors

`func (o *ConfigPutRequestSwitchpointSwitchpointName) SetTrafficMirrors(v []ConfigPutRequestSwitchpointSwitchpointNameTrafficMirrorsInner)`

SetTrafficMirrors sets TrafficMirrors field to given value.

### HasTrafficMirrors

`func (o *ConfigPutRequestSwitchpointSwitchpointName) HasTrafficMirrors() bool`

HasTrafficMirrors returns a boolean if a field has been set.

### GetEths

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetEths() []ConfigPutRequestSwitchpointSwitchpointNameEthsInner`

GetEths returns the Eths field if non-nil, zero value otherwise.

### GetEthsOk

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetEthsOk() (*[]ConfigPutRequestSwitchpointSwitchpointNameEthsInner, bool)`

GetEthsOk returns a tuple with the Eths field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEths

`func (o *ConfigPutRequestSwitchpointSwitchpointName) SetEths(v []ConfigPutRequestSwitchpointSwitchpointNameEthsInner)`

SetEths sets Eths field to given value.

### HasEths

`func (o *ConfigPutRequestSwitchpointSwitchpointName) HasEths() bool`

HasEths returns a boolean if a field has been set.

### GetObjectProperties

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetObjectProperties() ConfigPutRequestSwitchpointSwitchpointNameObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *ConfigPutRequestSwitchpointSwitchpointName) GetObjectPropertiesOk() (*ConfigPutRequestSwitchpointSwitchpointNameObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *ConfigPutRequestSwitchpointSwitchpointName) SetObjectProperties(v ConfigPutRequestSwitchpointSwitchpointNameObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *ConfigPutRequestSwitchpointSwitchpointName) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


