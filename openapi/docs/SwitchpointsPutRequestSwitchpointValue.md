# SwitchpointsPutRequestSwitchpointValue

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Name** | Pointer to **string** | Template Name. Must be unique within type. | [optional] [default to ""]
**Enable** | Pointer to **bool** | Enable object. It&#39;s highly recommended to set this value to true so that validation on the object will be ran. | [optional] [default to true]
**Tenant** | Pointer to **string** | The Tenant of this Device | [optional] [default to ""]
**TenantRefType** | Pointer to **string** | Object type for tenant field | [optional] 
**DeviceSerialNumber** | Pointer to **string** | Device Serial Number | [optional] [default to ""]
**ConnectedBundle** | Pointer to **string** | Connected Bundle | [optional] [default to ""]
**ConnectedBundleRefType** | Pointer to **string** | Object type for connected_bundle field | [optional] 
**IsTopOfIsland** | Pointer to **bool** | Mark this Switchpoint as Top of Island | [optional] [default to false]
**ReadOnlyMode** | Pointer to **bool** | When Read Only Mode is checked, vNetC will perform all functions except writing database updates to the target hardware | [optional] [default to false]
**Locked** | Pointer to **bool** | Permission lock | [optional] [default to false]
**ExpectedSite** | Pointer to **string** | Expected Fabric | [optional] [default to ""]
**ExpectedSiteRefType** | Pointer to **string** | Object type for expected_site field | [optional] 
**OutOfBandManagement** | Pointer to **bool** | For Switch Endpoints. Denotes a Switch is managed out of band via the management port | [optional] [default to false]
**Type** | Pointer to **string** | Type of Switchpoint | [optional] [default to "leaf"]
**SpinePlane** | Pointer to **string** | Spine Plane - subgrouping of super spine and spine | [optional] [default to ""]
**SpinePlaneRefType** | Pointer to **string** | Object type for spine_plane field | [optional] 
**Pod** | Pointer to **string** | Pod - subgrouping of spine and leaf switches | [optional] [default to ""]
**PodRefType** | Pointer to **string** | Object type for pod field | [optional] 
**Su** | Pointer to **string** | SU | [optional] [default to ""]
**SuRefType** | Pointer to **string** | Object type for su field | [optional] 
**SspGroup** | Pointer to **string** | SuperSpine Group - grouping of superspines in 3-tier config | [optional] [default to ""]
**SspGroupRefType** | Pointer to **string** | Object type for ssp_group field | [optional] 
**Rack** | Pointer to **string** | Physical Rack location of the Switch  | [optional] [default to ""]
**Position** | Pointer to **NullableFloat32** | Position of the Switch | [optional] 
**RailGroup** | Pointer to **NullableFloat32** | Rail Group the Switch is part of | [optional] 
**SwitchRouterIdIpMask** | Pointer to **string** | Switch BGP Router Identifier | [optional] [default to "(auto)"]
**SwitchRouterIdIpMaskAutoAssigned** | Pointer to **bool** | Whether or not the value in switch_router_id_ip_mask field has been automatically assigned or not. Set to false and change switch_router_id_ip_mask value to edit. | [optional] 
**SwitchVtepIdIpMask** | Pointer to **string** | Switch VETP Identifier | [optional] [default to "(auto)"]
**SwitchVtepIdIpMaskAutoAssigned** | Pointer to **bool** | Whether or not the value in switch_vtep_id_ip_mask field has been automatically assigned or not. Set to false and change switch_vtep_id_ip_mask value to edit. | [optional] 
**BgpAsNumber** | Pointer to **NullableInt32** | BGP Autonomous System Number for the Fabric Underlay  | [optional] 
**BgpAsNumberAutoAssigned** | Pointer to **bool** | Whether or not the value in bgp_as_number field has been automatically assigned or not. Set to false and change bgp_as_number value to edit. | [optional] 
**Badges** | Pointer to [**[]SwitchpointsPutRequestSwitchpointValueBadgesInner**](SwitchpointsPutRequestSwitchpointValueBadgesInner.md) |  | [optional] 
**Children** | Pointer to [**[]SwitchpointsPutRequestSwitchpointValueChildrenInner**](SwitchpointsPutRequestSwitchpointValueChildrenInner.md) |  | [optional] 
**TrafficMirrors** | Pointer to [**[]SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner**](SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner.md) |  | [optional] 
**Eths** | Pointer to [**[]SwitchpointsPutRequestSwitchpointValueEthsInner**](SwitchpointsPutRequestSwitchpointValueEthsInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**SwitchpointsPutRequestSwitchpointValueObjectProperties**](SwitchpointsPutRequestSwitchpointValueObjectProperties.md) |  | [optional] 
**IsFabric** | Pointer to **bool** | For Switch Endpoints. Denotes a Switch that is Fabric rather than an Edge Device | [optional] [default to false]
**Pots** | Pointer to [**[]SwitchpointsPutRequestSwitchpointValuePotsInner**](SwitchpointsPutRequestSwitchpointValuePotsInner.md) |  | [optional] 

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

### GetEnable

`func (o *SwitchpointsPutRequestSwitchpointValue) GetEnable() bool`

GetEnable returns the Enable field if non-nil, zero value otherwise.

### GetEnableOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetEnableOk() (*bool, bool)`

GetEnableOk returns a tuple with the Enable field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnable

`func (o *SwitchpointsPutRequestSwitchpointValue) SetEnable(v bool)`

SetEnable sets Enable field to given value.

### HasEnable

`func (o *SwitchpointsPutRequestSwitchpointValue) HasEnable() bool`

HasEnable returns a boolean if a field has been set.

### GetTenant

`func (o *SwitchpointsPutRequestSwitchpointValue) GetTenant() string`

GetTenant returns the Tenant field if non-nil, zero value otherwise.

### GetTenantOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetTenantOk() (*string, bool)`

GetTenantOk returns a tuple with the Tenant field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenant

`func (o *SwitchpointsPutRequestSwitchpointValue) SetTenant(v string)`

SetTenant sets Tenant field to given value.

### HasTenant

`func (o *SwitchpointsPutRequestSwitchpointValue) HasTenant() bool`

HasTenant returns a boolean if a field has been set.

### GetTenantRefType

`func (o *SwitchpointsPutRequestSwitchpointValue) GetTenantRefType() string`

GetTenantRefType returns the TenantRefType field if non-nil, zero value otherwise.

### GetTenantRefTypeOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetTenantRefTypeOk() (*string, bool)`

GetTenantRefTypeOk returns a tuple with the TenantRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenantRefType

`func (o *SwitchpointsPutRequestSwitchpointValue) SetTenantRefType(v string)`

SetTenantRefType sets TenantRefType field to given value.

### HasTenantRefType

`func (o *SwitchpointsPutRequestSwitchpointValue) HasTenantRefType() bool`

HasTenantRefType returns a boolean if a field has been set.

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

### GetIsTopOfIsland

`func (o *SwitchpointsPutRequestSwitchpointValue) GetIsTopOfIsland() bool`

GetIsTopOfIsland returns the IsTopOfIsland field if non-nil, zero value otherwise.

### GetIsTopOfIslandOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetIsTopOfIslandOk() (*bool, bool)`

GetIsTopOfIslandOk returns a tuple with the IsTopOfIsland field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsTopOfIsland

`func (o *SwitchpointsPutRequestSwitchpointValue) SetIsTopOfIsland(v bool)`

SetIsTopOfIsland sets IsTopOfIsland field to given value.

### HasIsTopOfIsland

`func (o *SwitchpointsPutRequestSwitchpointValue) HasIsTopOfIsland() bool`

HasIsTopOfIsland returns a boolean if a field has been set.

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

### GetExpectedSite

`func (o *SwitchpointsPutRequestSwitchpointValue) GetExpectedSite() string`

GetExpectedSite returns the ExpectedSite field if non-nil, zero value otherwise.

### GetExpectedSiteOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetExpectedSiteOk() (*string, bool)`

GetExpectedSiteOk returns a tuple with the ExpectedSite field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpectedSite

`func (o *SwitchpointsPutRequestSwitchpointValue) SetExpectedSite(v string)`

SetExpectedSite sets ExpectedSite field to given value.

### HasExpectedSite

`func (o *SwitchpointsPutRequestSwitchpointValue) HasExpectedSite() bool`

HasExpectedSite returns a boolean if a field has been set.

### GetExpectedSiteRefType

`func (o *SwitchpointsPutRequestSwitchpointValue) GetExpectedSiteRefType() string`

GetExpectedSiteRefType returns the ExpectedSiteRefType field if non-nil, zero value otherwise.

### GetExpectedSiteRefTypeOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetExpectedSiteRefTypeOk() (*string, bool)`

GetExpectedSiteRefTypeOk returns a tuple with the ExpectedSiteRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpectedSiteRefType

`func (o *SwitchpointsPutRequestSwitchpointValue) SetExpectedSiteRefType(v string)`

SetExpectedSiteRefType sets ExpectedSiteRefType field to given value.

### HasExpectedSiteRefType

`func (o *SwitchpointsPutRequestSwitchpointValue) HasExpectedSiteRefType() bool`

HasExpectedSiteRefType returns a boolean if a field has been set.

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

### GetSpinePlane

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSpinePlane() string`

GetSpinePlane returns the SpinePlane field if non-nil, zero value otherwise.

### GetSpinePlaneOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSpinePlaneOk() (*string, bool)`

GetSpinePlaneOk returns a tuple with the SpinePlane field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpinePlane

`func (o *SwitchpointsPutRequestSwitchpointValue) SetSpinePlane(v string)`

SetSpinePlane sets SpinePlane field to given value.

### HasSpinePlane

`func (o *SwitchpointsPutRequestSwitchpointValue) HasSpinePlane() bool`

HasSpinePlane returns a boolean if a field has been set.

### GetSpinePlaneRefType

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSpinePlaneRefType() string`

GetSpinePlaneRefType returns the SpinePlaneRefType field if non-nil, zero value otherwise.

### GetSpinePlaneRefTypeOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSpinePlaneRefTypeOk() (*string, bool)`

GetSpinePlaneRefTypeOk returns a tuple with the SpinePlaneRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSpinePlaneRefType

`func (o *SwitchpointsPutRequestSwitchpointValue) SetSpinePlaneRefType(v string)`

SetSpinePlaneRefType sets SpinePlaneRefType field to given value.

### HasSpinePlaneRefType

`func (o *SwitchpointsPutRequestSwitchpointValue) HasSpinePlaneRefType() bool`

HasSpinePlaneRefType returns a boolean if a field has been set.

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

### GetPodRefType

`func (o *SwitchpointsPutRequestSwitchpointValue) GetPodRefType() string`

GetPodRefType returns the PodRefType field if non-nil, zero value otherwise.

### GetPodRefTypeOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetPodRefTypeOk() (*string, bool)`

GetPodRefTypeOk returns a tuple with the PodRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPodRefType

`func (o *SwitchpointsPutRequestSwitchpointValue) SetPodRefType(v string)`

SetPodRefType sets PodRefType field to given value.

### HasPodRefType

`func (o *SwitchpointsPutRequestSwitchpointValue) HasPodRefType() bool`

HasPodRefType returns a boolean if a field has been set.

### GetSu

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSu() string`

GetSu returns the Su field if non-nil, zero value otherwise.

### GetSuOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSuOk() (*string, bool)`

GetSuOk returns a tuple with the Su field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSu

`func (o *SwitchpointsPutRequestSwitchpointValue) SetSu(v string)`

SetSu sets Su field to given value.

### HasSu

`func (o *SwitchpointsPutRequestSwitchpointValue) HasSu() bool`

HasSu returns a boolean if a field has been set.

### GetSuRefType

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSuRefType() string`

GetSuRefType returns the SuRefType field if non-nil, zero value otherwise.

### GetSuRefTypeOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSuRefTypeOk() (*string, bool)`

GetSuRefTypeOk returns a tuple with the SuRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSuRefType

`func (o *SwitchpointsPutRequestSwitchpointValue) SetSuRefType(v string)`

SetSuRefType sets SuRefType field to given value.

### HasSuRefType

`func (o *SwitchpointsPutRequestSwitchpointValue) HasSuRefType() bool`

HasSuRefType returns a boolean if a field has been set.

### GetSspGroup

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSspGroup() string`

GetSspGroup returns the SspGroup field if non-nil, zero value otherwise.

### GetSspGroupOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSspGroupOk() (*string, bool)`

GetSspGroupOk returns a tuple with the SspGroup field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSspGroup

`func (o *SwitchpointsPutRequestSwitchpointValue) SetSspGroup(v string)`

SetSspGroup sets SspGroup field to given value.

### HasSspGroup

`func (o *SwitchpointsPutRequestSwitchpointValue) HasSspGroup() bool`

HasSspGroup returns a boolean if a field has been set.

### GetSspGroupRefType

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSspGroupRefType() string`

GetSspGroupRefType returns the SspGroupRefType field if non-nil, zero value otherwise.

### GetSspGroupRefTypeOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSspGroupRefTypeOk() (*string, bool)`

GetSspGroupRefTypeOk returns a tuple with the SspGroupRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSspGroupRefType

`func (o *SwitchpointsPutRequestSwitchpointValue) SetSspGroupRefType(v string)`

SetSspGroupRefType sets SspGroupRefType field to given value.

### HasSspGroupRefType

`func (o *SwitchpointsPutRequestSwitchpointValue) HasSspGroupRefType() bool`

HasSspGroupRefType returns a boolean if a field has been set.

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

### GetPosition

`func (o *SwitchpointsPutRequestSwitchpointValue) GetPosition() float32`

GetPosition returns the Position field if non-nil, zero value otherwise.

### GetPositionOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetPositionOk() (*float32, bool)`

GetPositionOk returns a tuple with the Position field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPosition

`func (o *SwitchpointsPutRequestSwitchpointValue) SetPosition(v float32)`

SetPosition sets Position field to given value.

### HasPosition

`func (o *SwitchpointsPutRequestSwitchpointValue) HasPosition() bool`

HasPosition returns a boolean if a field has been set.

### SetPositionNil

`func (o *SwitchpointsPutRequestSwitchpointValue) SetPositionNil(b bool)`

 SetPositionNil sets the value for Position to be an explicit nil

### UnsetPosition
`func (o *SwitchpointsPutRequestSwitchpointValue) UnsetPosition()`

UnsetPosition ensures that no value is present for Position, not even an explicit nil
### GetRailGroup

`func (o *SwitchpointsPutRequestSwitchpointValue) GetRailGroup() float32`

GetRailGroup returns the RailGroup field if non-nil, zero value otherwise.

### GetRailGroupOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetRailGroupOk() (*float32, bool)`

GetRailGroupOk returns a tuple with the RailGroup field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRailGroup

`func (o *SwitchpointsPutRequestSwitchpointValue) SetRailGroup(v float32)`

SetRailGroup sets RailGroup field to given value.

### HasRailGroup

`func (o *SwitchpointsPutRequestSwitchpointValue) HasRailGroup() bool`

HasRailGroup returns a boolean if a field has been set.

### SetRailGroupNil

`func (o *SwitchpointsPutRequestSwitchpointValue) SetRailGroupNil(b bool)`

 SetRailGroupNil sets the value for RailGroup to be an explicit nil

### UnsetRailGroup
`func (o *SwitchpointsPutRequestSwitchpointValue) UnsetRailGroup()`

UnsetRailGroup ensures that no value is present for RailGroup, not even an explicit nil
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

### SetBgpAsNumberNil

`func (o *SwitchpointsPutRequestSwitchpointValue) SetBgpAsNumberNil(b bool)`

 SetBgpAsNumberNil sets the value for BgpAsNumber to be an explicit nil

### UnsetBgpAsNumber
`func (o *SwitchpointsPutRequestSwitchpointValue) UnsetBgpAsNumber()`

UnsetBgpAsNumber ensures that no value is present for BgpAsNumber, not even an explicit nil
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

### GetIsFabric

`func (o *SwitchpointsPutRequestSwitchpointValue) GetIsFabric() bool`

GetIsFabric returns the IsFabric field if non-nil, zero value otherwise.

### GetIsFabricOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetIsFabricOk() (*bool, bool)`

GetIsFabricOk returns a tuple with the IsFabric field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsFabric

`func (o *SwitchpointsPutRequestSwitchpointValue) SetIsFabric(v bool)`

SetIsFabric sets IsFabric field to given value.

### HasIsFabric

`func (o *SwitchpointsPutRequestSwitchpointValue) HasIsFabric() bool`

HasIsFabric returns a boolean if a field has been set.

### GetPots

`func (o *SwitchpointsPutRequestSwitchpointValue) GetPots() []SwitchpointsPutRequestSwitchpointValuePotsInner`

GetPots returns the Pots field if non-nil, zero value otherwise.

### GetPotsOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetPotsOk() (*[]SwitchpointsPutRequestSwitchpointValuePotsInner, bool)`

GetPotsOk returns a tuple with the Pots field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPots

`func (o *SwitchpointsPutRequestSwitchpointValue) SetPots(v []SwitchpointsPutRequestSwitchpointValuePotsInner)`

SetPots sets Pots field to given value.

### HasPots

`func (o *SwitchpointsPutRequestSwitchpointValue) HasPots() bool`

HasPots returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


