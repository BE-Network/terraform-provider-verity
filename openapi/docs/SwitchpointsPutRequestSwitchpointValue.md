# SwitchpointsPutRequestSwitchpointValue

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
**Badges** | Pointer to [**[]SwitchpointsPutRequestSwitchpointValueBadgesInner**](SwitchpointsPutRequestSwitchpointValueBadgesInner.md) |  | [optional] 
**Children** | Pointer to [**[]SwitchpointsPutRequestSwitchpointValueChildrenInner**](SwitchpointsPutRequestSwitchpointValueChildrenInner.md) |  | [optional] 
**TrafficMirrors** | Pointer to [**[]SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner**](SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner.md) |  | [optional] 
**Eths** | Pointer to [**[]SwitchpointsPutRequestSwitchpointValueEthsInner**](SwitchpointsPutRequestSwitchpointValueEthsInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**SwitchpointsPutRequestSwitchpointValueObjectProperties**](SwitchpointsPutRequestSwitchpointValueObjectProperties.md) |  | [optional] 

## Methods

### NewSwitchpointsPutRequestSwitchpointValue

`func NewSwitchpointsPutRequestSwitchpointValue() *SwitchpointsPutRequestSwitchpointValue`

NewSwitchpointsPutRequestSwitchpointValue instantiates a new SwitchpointsPutRequestSwitchpointValue object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewSwitchpointsPutRequestSwitchpointValueWithDefaults

`func NewSwitchpointsPutRequestSwitchpointValueWithDefaults() *SwitchpointsPutRequestSwitchpointValue`

NewSwitchpointsPutRequestSwitchpointValueWithDefaults instantiates a new SwitchpointsPutRequestSwitchpointValue object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetName

`func (o *SwitchpointsPutRequestSwitchpointValue) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *SwitchpointsPutRequestSwitchpointValue) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *SwitchpointsPutRequestSwitchpointValue) HasName() bool`

HasName returns a boolean if a field has been set.

### GetDeviceSerialNumber

`func (o *SwitchpointsPutRequestSwitchpointValue) GetDeviceSerialNumber() string`

GetDeviceSerialNumber returns the DeviceSerialNumber field if non-nil, zero value otherwise.

### GetDeviceSerialNumberOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetDeviceSerialNumberOk() (*string, bool)`

GetDeviceSerialNumberOk returns a tuple with the DeviceSerialNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceSerialNumber

`func (o *SwitchpointsPutRequestSwitchpointValue) SetDeviceSerialNumber(v string)`

SetDeviceSerialNumber sets DeviceSerialNumber field to given value.

### HasDeviceSerialNumber

`func (o *SwitchpointsPutRequestSwitchpointValue) HasDeviceSerialNumber() bool`

HasDeviceSerialNumber returns a boolean if a field has been set.

### GetConnectedBundle

`func (o *SwitchpointsPutRequestSwitchpointValue) GetConnectedBundle() string`

GetConnectedBundle returns the ConnectedBundle field if non-nil, zero value otherwise.

### GetConnectedBundleOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetConnectedBundleOk() (*string, bool)`

GetConnectedBundleOk returns a tuple with the ConnectedBundle field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnectedBundle

`func (o *SwitchpointsPutRequestSwitchpointValue) SetConnectedBundle(v string)`

SetConnectedBundle sets ConnectedBundle field to given value.

### HasConnectedBundle

`func (o *SwitchpointsPutRequestSwitchpointValue) HasConnectedBundle() bool`

HasConnectedBundle returns a boolean if a field has been set.

### GetConnectedBundleRefType

`func (o *SwitchpointsPutRequestSwitchpointValue) GetConnectedBundleRefType() string`

GetConnectedBundleRefType returns the ConnectedBundleRefType field if non-nil, zero value otherwise.

### GetConnectedBundleRefTypeOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetConnectedBundleRefTypeOk() (*string, bool)`

GetConnectedBundleRefTypeOk returns a tuple with the ConnectedBundleRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnectedBundleRefType

`func (o *SwitchpointsPutRequestSwitchpointValue) SetConnectedBundleRefType(v string)`

SetConnectedBundleRefType sets ConnectedBundleRefType field to given value.

### HasConnectedBundleRefType

`func (o *SwitchpointsPutRequestSwitchpointValue) HasConnectedBundleRefType() bool`

HasConnectedBundleRefType returns a boolean if a field has been set.

### GetReadOnlyMode

`func (o *SwitchpointsPutRequestSwitchpointValue) GetReadOnlyMode() bool`

GetReadOnlyMode returns the ReadOnlyMode field if non-nil, zero value otherwise.

### GetReadOnlyModeOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetReadOnlyModeOk() (*bool, bool)`

GetReadOnlyModeOk returns a tuple with the ReadOnlyMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetReadOnlyMode

`func (o *SwitchpointsPutRequestSwitchpointValue) SetReadOnlyMode(v bool)`

SetReadOnlyMode sets ReadOnlyMode field to given value.

### HasReadOnlyMode

`func (o *SwitchpointsPutRequestSwitchpointValue) HasReadOnlyMode() bool`

HasReadOnlyMode returns a boolean if a field has been set.

### GetLocked

`func (o *SwitchpointsPutRequestSwitchpointValue) GetLocked() bool`

GetLocked returns the Locked field if non-nil, zero value otherwise.

### GetLockedOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetLockedOk() (*bool, bool)`

GetLockedOk returns a tuple with the Locked field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocked

`func (o *SwitchpointsPutRequestSwitchpointValue) SetLocked(v bool)`

SetLocked sets Locked field to given value.

### HasLocked

`func (o *SwitchpointsPutRequestSwitchpointValue) HasLocked() bool`

HasLocked returns a boolean if a field has been set.

### GetDisabledPorts

`func (o *SwitchpointsPutRequestSwitchpointValue) GetDisabledPorts() string`

GetDisabledPorts returns the DisabledPorts field if non-nil, zero value otherwise.

### GetDisabledPortsOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetDisabledPortsOk() (*string, bool)`

GetDisabledPortsOk returns a tuple with the DisabledPorts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDisabledPorts

`func (o *SwitchpointsPutRequestSwitchpointValue) SetDisabledPorts(v string)`

SetDisabledPorts sets DisabledPorts field to given value.

### HasDisabledPorts

`func (o *SwitchpointsPutRequestSwitchpointValue) HasDisabledPorts() bool`

HasDisabledPorts returns a boolean if a field has been set.

### GetOutOfBandManagement

`func (o *SwitchpointsPutRequestSwitchpointValue) GetOutOfBandManagement() bool`

GetOutOfBandManagement returns the OutOfBandManagement field if non-nil, zero value otherwise.

### GetOutOfBandManagementOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetOutOfBandManagementOk() (*bool, bool)`

GetOutOfBandManagementOk returns a tuple with the OutOfBandManagement field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOutOfBandManagement

`func (o *SwitchpointsPutRequestSwitchpointValue) SetOutOfBandManagement(v bool)`

SetOutOfBandManagement sets OutOfBandManagement field to given value.

### HasOutOfBandManagement

`func (o *SwitchpointsPutRequestSwitchpointValue) HasOutOfBandManagement() bool`

HasOutOfBandManagement returns a boolean if a field has been set.

### GetType

`func (o *SwitchpointsPutRequestSwitchpointValue) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *SwitchpointsPutRequestSwitchpointValue) SetType(v string)`

SetType sets Type field to given value.

### HasType

`func (o *SwitchpointsPutRequestSwitchpointValue) HasType() bool`

HasType returns a boolean if a field has been set.

### GetSuperPod

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSuperPod() string`

GetSuperPod returns the SuperPod field if non-nil, zero value otherwise.

### GetSuperPodOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSuperPodOk() (*string, bool)`

GetSuperPodOk returns a tuple with the SuperPod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSuperPod

`func (o *SwitchpointsPutRequestSwitchpointValue) SetSuperPod(v string)`

SetSuperPod sets SuperPod field to given value.

### HasSuperPod

`func (o *SwitchpointsPutRequestSwitchpointValue) HasSuperPod() bool`

HasSuperPod returns a boolean if a field has been set.

### GetPod

`func (o *SwitchpointsPutRequestSwitchpointValue) GetPod() string`

GetPod returns the Pod field if non-nil, zero value otherwise.

### GetPodOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetPodOk() (*string, bool)`

GetPodOk returns a tuple with the Pod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPod

`func (o *SwitchpointsPutRequestSwitchpointValue) SetPod(v string)`

SetPod sets Pod field to given value.

### HasPod

`func (o *SwitchpointsPutRequestSwitchpointValue) HasPod() bool`

HasPod returns a boolean if a field has been set.

### GetRack

`func (o *SwitchpointsPutRequestSwitchpointValue) GetRack() string`

GetRack returns the Rack field if non-nil, zero value otherwise.

### GetRackOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetRackOk() (*string, bool)`

GetRackOk returns a tuple with the Rack field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRack

`func (o *SwitchpointsPutRequestSwitchpointValue) SetRack(v string)`

SetRack sets Rack field to given value.

### HasRack

`func (o *SwitchpointsPutRequestSwitchpointValue) HasRack() bool`

HasRack returns a boolean if a field has been set.

### GetSwitchRouterIdIpMask

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSwitchRouterIdIpMask() string`

GetSwitchRouterIdIpMask returns the SwitchRouterIdIpMask field if non-nil, zero value otherwise.

### GetSwitchRouterIdIpMaskOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSwitchRouterIdIpMaskOk() (*string, bool)`

GetSwitchRouterIdIpMaskOk returns a tuple with the SwitchRouterIdIpMask field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchRouterIdIpMask

`func (o *SwitchpointsPutRequestSwitchpointValue) SetSwitchRouterIdIpMask(v string)`

SetSwitchRouterIdIpMask sets SwitchRouterIdIpMask field to given value.

### HasSwitchRouterIdIpMask

`func (o *SwitchpointsPutRequestSwitchpointValue) HasSwitchRouterIdIpMask() bool`

HasSwitchRouterIdIpMask returns a boolean if a field has been set.

### GetSwitchRouterIdIpMaskAutoAssigned

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSwitchRouterIdIpMaskAutoAssigned() bool`

GetSwitchRouterIdIpMaskAutoAssigned returns the SwitchRouterIdIpMaskAutoAssigned field if non-nil, zero value otherwise.

### GetSwitchRouterIdIpMaskAutoAssignedOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSwitchRouterIdIpMaskAutoAssignedOk() (*bool, bool)`

GetSwitchRouterIdIpMaskAutoAssignedOk returns a tuple with the SwitchRouterIdIpMaskAutoAssigned field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchRouterIdIpMaskAutoAssigned

`func (o *SwitchpointsPutRequestSwitchpointValue) SetSwitchRouterIdIpMaskAutoAssigned(v bool)`

SetSwitchRouterIdIpMaskAutoAssigned sets SwitchRouterIdIpMaskAutoAssigned field to given value.

### HasSwitchRouterIdIpMaskAutoAssigned

`func (o *SwitchpointsPutRequestSwitchpointValue) HasSwitchRouterIdIpMaskAutoAssigned() bool`

HasSwitchRouterIdIpMaskAutoAssigned returns a boolean if a field has been set.

### GetSwitchVtepIdIpMask

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSwitchVtepIdIpMask() string`

GetSwitchVtepIdIpMask returns the SwitchVtepIdIpMask field if non-nil, zero value otherwise.

### GetSwitchVtepIdIpMaskOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSwitchVtepIdIpMaskOk() (*string, bool)`

GetSwitchVtepIdIpMaskOk returns a tuple with the SwitchVtepIdIpMask field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchVtepIdIpMask

`func (o *SwitchpointsPutRequestSwitchpointValue) SetSwitchVtepIdIpMask(v string)`

SetSwitchVtepIdIpMask sets SwitchVtepIdIpMask field to given value.

### HasSwitchVtepIdIpMask

`func (o *SwitchpointsPutRequestSwitchpointValue) HasSwitchVtepIdIpMask() bool`

HasSwitchVtepIdIpMask returns a boolean if a field has been set.

### GetSwitchVtepIdIpMaskAutoAssigned

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSwitchVtepIdIpMaskAutoAssigned() bool`

GetSwitchVtepIdIpMaskAutoAssigned returns the SwitchVtepIdIpMaskAutoAssigned field if non-nil, zero value otherwise.

### GetSwitchVtepIdIpMaskAutoAssignedOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSwitchVtepIdIpMaskAutoAssignedOk() (*bool, bool)`

GetSwitchVtepIdIpMaskAutoAssignedOk returns a tuple with the SwitchVtepIdIpMaskAutoAssigned field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchVtepIdIpMaskAutoAssigned

`func (o *SwitchpointsPutRequestSwitchpointValue) SetSwitchVtepIdIpMaskAutoAssigned(v bool)`

SetSwitchVtepIdIpMaskAutoAssigned sets SwitchVtepIdIpMaskAutoAssigned field to given value.

### HasSwitchVtepIdIpMaskAutoAssigned

`func (o *SwitchpointsPutRequestSwitchpointValue) HasSwitchVtepIdIpMaskAutoAssigned() bool`

HasSwitchVtepIdIpMaskAutoAssigned returns a boolean if a field has been set.

### GetBgpAsNumber

`func (o *SwitchpointsPutRequestSwitchpointValue) GetBgpAsNumber() int32`

GetBgpAsNumber returns the BgpAsNumber field if non-nil, zero value otherwise.

### GetBgpAsNumberOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetBgpAsNumberOk() (*int32, bool)`

GetBgpAsNumberOk returns a tuple with the BgpAsNumber field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBgpAsNumber

`func (o *SwitchpointsPutRequestSwitchpointValue) SetBgpAsNumber(v int32)`

SetBgpAsNumber sets BgpAsNumber field to given value.

### HasBgpAsNumber

`func (o *SwitchpointsPutRequestSwitchpointValue) HasBgpAsNumber() bool`

HasBgpAsNumber returns a boolean if a field has been set.

### GetBgpAsNumberAutoAssigned

`func (o *SwitchpointsPutRequestSwitchpointValue) GetBgpAsNumberAutoAssigned() bool`

GetBgpAsNumberAutoAssigned returns the BgpAsNumberAutoAssigned field if non-nil, zero value otherwise.

### GetBgpAsNumberAutoAssignedOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetBgpAsNumberAutoAssignedOk() (*bool, bool)`

GetBgpAsNumberAutoAssignedOk returns a tuple with the BgpAsNumberAutoAssigned field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBgpAsNumberAutoAssigned

`func (o *SwitchpointsPutRequestSwitchpointValue) SetBgpAsNumberAutoAssigned(v bool)`

SetBgpAsNumberAutoAssigned sets BgpAsNumberAutoAssigned field to given value.

### HasBgpAsNumberAutoAssigned

`func (o *SwitchpointsPutRequestSwitchpointValue) HasBgpAsNumberAutoAssigned() bool`

HasBgpAsNumberAutoAssigned returns a boolean if a field has been set.

### GetBadges

`func (o *SwitchpointsPutRequestSwitchpointValue) GetBadges() []SwitchpointsPutRequestSwitchpointValueBadgesInner`

GetBadges returns the Badges field if non-nil, zero value otherwise.

### GetBadgesOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetBadgesOk() (*[]SwitchpointsPutRequestSwitchpointValueBadgesInner, bool)`

GetBadgesOk returns a tuple with the Badges field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBadges

`func (o *SwitchpointsPutRequestSwitchpointValue) SetBadges(v []SwitchpointsPutRequestSwitchpointValueBadgesInner)`

SetBadges sets Badges field to given value.

### HasBadges

`func (o *SwitchpointsPutRequestSwitchpointValue) HasBadges() bool`

HasBadges returns a boolean if a field has been set.

### GetChildren

`func (o *SwitchpointsPutRequestSwitchpointValue) GetChildren() []SwitchpointsPutRequestSwitchpointValueChildrenInner`

GetChildren returns the Children field if non-nil, zero value otherwise.

### GetChildrenOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetChildrenOk() (*[]SwitchpointsPutRequestSwitchpointValueChildrenInner, bool)`

GetChildrenOk returns a tuple with the Children field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetChildren

`func (o *SwitchpointsPutRequestSwitchpointValue) SetChildren(v []SwitchpointsPutRequestSwitchpointValueChildrenInner)`

SetChildren sets Children field to given value.

### HasChildren

`func (o *SwitchpointsPutRequestSwitchpointValue) HasChildren() bool`

HasChildren returns a boolean if a field has been set.

### GetTrafficMirrors

`func (o *SwitchpointsPutRequestSwitchpointValue) GetTrafficMirrors() []SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner`

GetTrafficMirrors returns the TrafficMirrors field if non-nil, zero value otherwise.

### GetTrafficMirrorsOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetTrafficMirrorsOk() (*[]SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner, bool)`

GetTrafficMirrorsOk returns a tuple with the TrafficMirrors field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTrafficMirrors

`func (o *SwitchpointsPutRequestSwitchpointValue) SetTrafficMirrors(v []SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner)`

SetTrafficMirrors sets TrafficMirrors field to given value.

### HasTrafficMirrors

`func (o *SwitchpointsPutRequestSwitchpointValue) HasTrafficMirrors() bool`

HasTrafficMirrors returns a boolean if a field has been set.

### GetEths

`func (o *SwitchpointsPutRequestSwitchpointValue) GetEths() []SwitchpointsPutRequestSwitchpointValueEthsInner`

GetEths returns the Eths field if non-nil, zero value otherwise.

### GetEthsOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetEthsOk() (*[]SwitchpointsPutRequestSwitchpointValueEthsInner, bool)`

GetEthsOk returns a tuple with the Eths field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEths

`func (o *SwitchpointsPutRequestSwitchpointValue) SetEths(v []SwitchpointsPutRequestSwitchpointValueEthsInner)`

SetEths sets Eths field to given value.

### HasEths

`func (o *SwitchpointsPutRequestSwitchpointValue) HasEths() bool`

HasEths returns a boolean if a field has been set.

### GetObjectProperties

`func (o *SwitchpointsPutRequestSwitchpointValue) GetObjectProperties() SwitchpointsPutRequestSwitchpointValueObjectProperties`

GetObjectProperties returns the ObjectProperties field if non-nil, zero value otherwise.

### GetObjectPropertiesOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetObjectPropertiesOk() (*SwitchpointsPutRequestSwitchpointValueObjectProperties, bool)`

GetObjectPropertiesOk returns a tuple with the ObjectProperties field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetObjectProperties

`func (o *SwitchpointsPutRequestSwitchpointValue) SetObjectProperties(v SwitchpointsPutRequestSwitchpointValueObjectProperties)`

SetObjectProperties sets ObjectProperties field to given value.

### HasObjectProperties

`func (o *SwitchpointsPutRequestSwitchpointValue) HasObjectProperties() bool`

HasObjectProperties returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


