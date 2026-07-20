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
**Plane** | Pointer to **string** | Plane | [optional] [default to ""]
**PlaneRefType** | Pointer to **string** | Object type for plane field | [optional] 
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
**BbSwitch** | Pointer to **bool** | Expose fields for Device Management | [optional] [default to false]
**PasswordEncrypted** | Pointer to **string** | Password | [optional] [default to ""]
**EnablePasswordEncrypted** | Pointer to **string** | Enable Password - to enable privileged CLI operations | [optional] [default to ""]
**SshKeyOrPasswordEncrypted** | Pointer to **string** | SSH Key or Password | [optional] [default to ""]
**PassphraseEncrypted** | Pointer to **string** | Passphrase | [optional] [default to ""]
**PrivatePasswordEncrypted** | Pointer to **string** | Password | [optional] [default to ""]
**IpSource** | Pointer to **string** | IP Source | [optional] [default to "dhcp"]
**ControllerIpAndMask** | Pointer to **string** | Controller IP and Mask | [optional] [default to ""]
**ControllerIpAndMaskAutoAssigned** | Pointer to **bool** | Whether or not the value in controller_ip_and_mask field has been automatically assigned or not. Set to false and change controller_ip_and_mask value to edit. | [optional] 
**Gateway** | Pointer to **string** | Gateway | [optional] [default to ""]
**SwitchIpAndMask** | Pointer to **string** | Switch IP and Mask | [optional] [default to ""]
**SwitchIpAndMaskAutoAssigned** | Pointer to **bool** | Whether or not the value in switch_ip_and_mask field has been automatically assigned or not. Set to false and change switch_ip_and_mask value to edit. | [optional] 
**SwitchGateway** | Pointer to **string** | Gateway of Managed Device | [optional] [default to ""]
**CommType** | Pointer to **string** | Comm Type | [optional] [default to "snmpv2"]
**SnmpCommunityString** | Pointer to **string** | Comm Credentials | [optional] [default to ""]
**UplinkPort** | Pointer to **string** | Uplink Port of Managed Device | [optional] [default to ""]
**LldpSearchString** | Pointer to **string** | Optional unless Located By is \&quot;LLDP\&quot; or Device managed as \&quot;Active SFP\&quot;. Must be either the chassis-id or the hostname of the LLDP from the managed device. Used to detect connections between managed devices. If blank, the chassis-id detected by the Device Controller via SNMP/CLI is used | [optional] [default to ""]
**ZtpIdentification** | Pointer to **string** | Service Tag or Serial Number to identify device for Zero Touch Provisioning | [optional] [default to ""]
**LocatedBy** | Pointer to **string** | Controls how the system locates this Device within its LAN | [optional] [default to "LLDP"]
**PowerState** | Pointer to **string** | Power state of Switch Controller | [optional] [default to "on"]
**CommunicationMode** | Pointer to **string** | Select the network operating system (NOS) type for this endpoint. | [optional] [default to "generic_snmp"]
**CliAccessMode** | Pointer to **string** | CLI Access Mode | [optional] [default to "SSH"]
**Username** | Pointer to **string** | Username | [optional] [default to ""]
**Password** | Pointer to **string** | Password | [optional] [default to ""]
**EnablePassword** | Pointer to **string** | Enable Password - to enable privileged CLI operations | [optional] [default to ""]
**SshKeyOrPassword** | Pointer to **string** | SSH Key or Password | [optional] [default to ""]
**ManagedOnNativeVlan** | Pointer to **bool** | Managed on native VLAN | [optional] [default to false]
**Sdlc** | Pointer to **string** | SDLC that Device Controller belongs to | [optional] [default to ""]
**SecurityType** | Pointer to **string** | Security level | [optional] [default to "noAuthNoPriv"]
**Snmpv3Username** | Pointer to **string** | Username | [optional] [default to ""]
**AuthenticationProtocol** | Pointer to **string** | Protocol | [optional] [default to "MD5"]
**Passphrase** | Pointer to **string** | Passphrase | [optional] [default to ""]
**PrivateProtocol** | Pointer to **string** | Protocol | [optional] [default to "DES"]
**PrivatePassword** | Pointer to **string** | Password | [optional] [default to ""]
**Badges** | Pointer to [**[]SwitchpointsPutRequestSwitchpointValueBadgesInner**](SwitchpointsPutRequestSwitchpointValueBadgesInner.md) |  | [optional] 
**Children** | Pointer to [**[]SwitchpointsPutRequestSwitchpointValueChildrenInner**](SwitchpointsPutRequestSwitchpointValueChildrenInner.md) |  | [optional] 
**TrafficMirrors** | Pointer to [**[]SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner**](SwitchpointsPutRequestSwitchpointValueTrafficMirrorsInner.md) |  | [optional] 
**Eths** | Pointer to [**[]SwitchpointsPutRequestSwitchpointValueEthsInner**](SwitchpointsPutRequestSwitchpointValueEthsInner.md) |  | [optional] 
**ObjectProperties** | Pointer to [**SwitchpointsPutRequestSwitchpointValueObjectProperties**](SwitchpointsPutRequestSwitchpointValueObjectProperties.md) |  | [optional] 
**IsFabric** | Pointer to **bool** | For Switch Endpoints. Denotes a Switch that is Fabric rather than an Edge Device | [optional] [default to false]
**DeviceManagedAs** | Pointer to **string** | Device managed as | [optional] [default to "switch"]
**Switch** | Pointer to **string** | Switchpoint locating the Switch to be controlled | [optional] [default to ""]
**SwitchRefType** | Pointer to **string** | Object type for switch field | [optional] 
**ConnectionService** | Pointer to **string** | Connect a Service | [optional] [default to ""]
**ConnectionServiceRefType** | Pointer to **string** | Object type for connection_service field | [optional] 
**Port** | Pointer to **string** | Port locating the Switch to be controlled | [optional] [default to ""]
**UsesTaggedPackets** | Pointer to **bool** | Indicates if the direct interface expects tagged or untagged packets | [optional] [default to true]
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

### GetPlane

`func (o *SwitchpointsPutRequestSwitchpointValue) GetPlane() string`

GetPlane returns the Plane field if non-nil, zero value otherwise.

### GetPlaneOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetPlaneOk() (*string, bool)`

GetPlaneOk returns a tuple with the Plane field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPlane

`func (o *SwitchpointsPutRequestSwitchpointValue) SetPlane(v string)`

SetPlane sets Plane field to given value.

### HasPlane

`func (o *SwitchpointsPutRequestSwitchpointValue) HasPlane() bool`

HasPlane returns a boolean if a field has been set.

### GetPlaneRefType

`func (o *SwitchpointsPutRequestSwitchpointValue) GetPlaneRefType() string`

GetPlaneRefType returns the PlaneRefType field if non-nil, zero value otherwise.

### GetPlaneRefTypeOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetPlaneRefTypeOk() (*string, bool)`

GetPlaneRefTypeOk returns a tuple with the PlaneRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPlaneRefType

`func (o *SwitchpointsPutRequestSwitchpointValue) SetPlaneRefType(v string)`

SetPlaneRefType sets PlaneRefType field to given value.

### HasPlaneRefType

`func (o *SwitchpointsPutRequestSwitchpointValue) HasPlaneRefType() bool`

HasPlaneRefType returns a boolean if a field has been set.

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

### GetBbSwitch

`func (o *SwitchpointsPutRequestSwitchpointValue) GetBbSwitch() bool`

GetBbSwitch returns the BbSwitch field if non-nil, zero value otherwise.

### GetBbSwitchOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetBbSwitchOk() (*bool, bool)`

GetBbSwitchOk returns a tuple with the BbSwitch field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBbSwitch

`func (o *SwitchpointsPutRequestSwitchpointValue) SetBbSwitch(v bool)`

SetBbSwitch sets BbSwitch field to given value.

### HasBbSwitch

`func (o *SwitchpointsPutRequestSwitchpointValue) HasBbSwitch() bool`

HasBbSwitch returns a boolean if a field has been set.

### GetPasswordEncrypted

`func (o *SwitchpointsPutRequestSwitchpointValue) GetPasswordEncrypted() string`

GetPasswordEncrypted returns the PasswordEncrypted field if non-nil, zero value otherwise.

### GetPasswordEncryptedOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetPasswordEncryptedOk() (*string, bool)`

GetPasswordEncryptedOk returns a tuple with the PasswordEncrypted field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPasswordEncrypted

`func (o *SwitchpointsPutRequestSwitchpointValue) SetPasswordEncrypted(v string)`

SetPasswordEncrypted sets PasswordEncrypted field to given value.

### HasPasswordEncrypted

`func (o *SwitchpointsPutRequestSwitchpointValue) HasPasswordEncrypted() bool`

HasPasswordEncrypted returns a boolean if a field has been set.

### GetEnablePasswordEncrypted

`func (o *SwitchpointsPutRequestSwitchpointValue) GetEnablePasswordEncrypted() string`

GetEnablePasswordEncrypted returns the EnablePasswordEncrypted field if non-nil, zero value otherwise.

### GetEnablePasswordEncryptedOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetEnablePasswordEncryptedOk() (*string, bool)`

GetEnablePasswordEncryptedOk returns a tuple with the EnablePasswordEncrypted field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnablePasswordEncrypted

`func (o *SwitchpointsPutRequestSwitchpointValue) SetEnablePasswordEncrypted(v string)`

SetEnablePasswordEncrypted sets EnablePasswordEncrypted field to given value.

### HasEnablePasswordEncrypted

`func (o *SwitchpointsPutRequestSwitchpointValue) HasEnablePasswordEncrypted() bool`

HasEnablePasswordEncrypted returns a boolean if a field has been set.

### GetSshKeyOrPasswordEncrypted

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSshKeyOrPasswordEncrypted() string`

GetSshKeyOrPasswordEncrypted returns the SshKeyOrPasswordEncrypted field if non-nil, zero value otherwise.

### GetSshKeyOrPasswordEncryptedOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSshKeyOrPasswordEncryptedOk() (*string, bool)`

GetSshKeyOrPasswordEncryptedOk returns a tuple with the SshKeyOrPasswordEncrypted field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSshKeyOrPasswordEncrypted

`func (o *SwitchpointsPutRequestSwitchpointValue) SetSshKeyOrPasswordEncrypted(v string)`

SetSshKeyOrPasswordEncrypted sets SshKeyOrPasswordEncrypted field to given value.

### HasSshKeyOrPasswordEncrypted

`func (o *SwitchpointsPutRequestSwitchpointValue) HasSshKeyOrPasswordEncrypted() bool`

HasSshKeyOrPasswordEncrypted returns a boolean if a field has been set.

### GetPassphraseEncrypted

`func (o *SwitchpointsPutRequestSwitchpointValue) GetPassphraseEncrypted() string`

GetPassphraseEncrypted returns the PassphraseEncrypted field if non-nil, zero value otherwise.

### GetPassphraseEncryptedOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetPassphraseEncryptedOk() (*string, bool)`

GetPassphraseEncryptedOk returns a tuple with the PassphraseEncrypted field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassphraseEncrypted

`func (o *SwitchpointsPutRequestSwitchpointValue) SetPassphraseEncrypted(v string)`

SetPassphraseEncrypted sets PassphraseEncrypted field to given value.

### HasPassphraseEncrypted

`func (o *SwitchpointsPutRequestSwitchpointValue) HasPassphraseEncrypted() bool`

HasPassphraseEncrypted returns a boolean if a field has been set.

### GetPrivatePasswordEncrypted

`func (o *SwitchpointsPutRequestSwitchpointValue) GetPrivatePasswordEncrypted() string`

GetPrivatePasswordEncrypted returns the PrivatePasswordEncrypted field if non-nil, zero value otherwise.

### GetPrivatePasswordEncryptedOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetPrivatePasswordEncryptedOk() (*string, bool)`

GetPrivatePasswordEncryptedOk returns a tuple with the PrivatePasswordEncrypted field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrivatePasswordEncrypted

`func (o *SwitchpointsPutRequestSwitchpointValue) SetPrivatePasswordEncrypted(v string)`

SetPrivatePasswordEncrypted sets PrivatePasswordEncrypted field to given value.

### HasPrivatePasswordEncrypted

`func (o *SwitchpointsPutRequestSwitchpointValue) HasPrivatePasswordEncrypted() bool`

HasPrivatePasswordEncrypted returns a boolean if a field has been set.

### GetIpSource

`func (o *SwitchpointsPutRequestSwitchpointValue) GetIpSource() string`

GetIpSource returns the IpSource field if non-nil, zero value otherwise.

### GetIpSourceOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetIpSourceOk() (*string, bool)`

GetIpSourceOk returns a tuple with the IpSource field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpSource

`func (o *SwitchpointsPutRequestSwitchpointValue) SetIpSource(v string)`

SetIpSource sets IpSource field to given value.

### HasIpSource

`func (o *SwitchpointsPutRequestSwitchpointValue) HasIpSource() bool`

HasIpSource returns a boolean if a field has been set.

### GetControllerIpAndMask

`func (o *SwitchpointsPutRequestSwitchpointValue) GetControllerIpAndMask() string`

GetControllerIpAndMask returns the ControllerIpAndMask field if non-nil, zero value otherwise.

### GetControllerIpAndMaskOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetControllerIpAndMaskOk() (*string, bool)`

GetControllerIpAndMaskOk returns a tuple with the ControllerIpAndMask field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetControllerIpAndMask

`func (o *SwitchpointsPutRequestSwitchpointValue) SetControllerIpAndMask(v string)`

SetControllerIpAndMask sets ControllerIpAndMask field to given value.

### HasControllerIpAndMask

`func (o *SwitchpointsPutRequestSwitchpointValue) HasControllerIpAndMask() bool`

HasControllerIpAndMask returns a boolean if a field has been set.

### GetControllerIpAndMaskAutoAssigned

`func (o *SwitchpointsPutRequestSwitchpointValue) GetControllerIpAndMaskAutoAssigned() bool`

GetControllerIpAndMaskAutoAssigned returns the ControllerIpAndMaskAutoAssigned field if non-nil, zero value otherwise.

### GetControllerIpAndMaskAutoAssignedOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetControllerIpAndMaskAutoAssignedOk() (*bool, bool)`

GetControllerIpAndMaskAutoAssignedOk returns a tuple with the ControllerIpAndMaskAutoAssigned field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetControllerIpAndMaskAutoAssigned

`func (o *SwitchpointsPutRequestSwitchpointValue) SetControllerIpAndMaskAutoAssigned(v bool)`

SetControllerIpAndMaskAutoAssigned sets ControllerIpAndMaskAutoAssigned field to given value.

### HasControllerIpAndMaskAutoAssigned

`func (o *SwitchpointsPutRequestSwitchpointValue) HasControllerIpAndMaskAutoAssigned() bool`

HasControllerIpAndMaskAutoAssigned returns a boolean if a field has been set.

### GetGateway

`func (o *SwitchpointsPutRequestSwitchpointValue) GetGateway() string`

GetGateway returns the Gateway field if non-nil, zero value otherwise.

### GetGatewayOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetGatewayOk() (*string, bool)`

GetGatewayOk returns a tuple with the Gateway field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGateway

`func (o *SwitchpointsPutRequestSwitchpointValue) SetGateway(v string)`

SetGateway sets Gateway field to given value.

### HasGateway

`func (o *SwitchpointsPutRequestSwitchpointValue) HasGateway() bool`

HasGateway returns a boolean if a field has been set.

### GetSwitchIpAndMask

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSwitchIpAndMask() string`

GetSwitchIpAndMask returns the SwitchIpAndMask field if non-nil, zero value otherwise.

### GetSwitchIpAndMaskOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSwitchIpAndMaskOk() (*string, bool)`

GetSwitchIpAndMaskOk returns a tuple with the SwitchIpAndMask field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchIpAndMask

`func (o *SwitchpointsPutRequestSwitchpointValue) SetSwitchIpAndMask(v string)`

SetSwitchIpAndMask sets SwitchIpAndMask field to given value.

### HasSwitchIpAndMask

`func (o *SwitchpointsPutRequestSwitchpointValue) HasSwitchIpAndMask() bool`

HasSwitchIpAndMask returns a boolean if a field has been set.

### GetSwitchIpAndMaskAutoAssigned

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSwitchIpAndMaskAutoAssigned() bool`

GetSwitchIpAndMaskAutoAssigned returns the SwitchIpAndMaskAutoAssigned field if non-nil, zero value otherwise.

### GetSwitchIpAndMaskAutoAssignedOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSwitchIpAndMaskAutoAssignedOk() (*bool, bool)`

GetSwitchIpAndMaskAutoAssignedOk returns a tuple with the SwitchIpAndMaskAutoAssigned field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchIpAndMaskAutoAssigned

`func (o *SwitchpointsPutRequestSwitchpointValue) SetSwitchIpAndMaskAutoAssigned(v bool)`

SetSwitchIpAndMaskAutoAssigned sets SwitchIpAndMaskAutoAssigned field to given value.

### HasSwitchIpAndMaskAutoAssigned

`func (o *SwitchpointsPutRequestSwitchpointValue) HasSwitchIpAndMaskAutoAssigned() bool`

HasSwitchIpAndMaskAutoAssigned returns a boolean if a field has been set.

### GetSwitchGateway

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSwitchGateway() string`

GetSwitchGateway returns the SwitchGateway field if non-nil, zero value otherwise.

### GetSwitchGatewayOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSwitchGatewayOk() (*string, bool)`

GetSwitchGatewayOk returns a tuple with the SwitchGateway field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchGateway

`func (o *SwitchpointsPutRequestSwitchpointValue) SetSwitchGateway(v string)`

SetSwitchGateway sets SwitchGateway field to given value.

### HasSwitchGateway

`func (o *SwitchpointsPutRequestSwitchpointValue) HasSwitchGateway() bool`

HasSwitchGateway returns a boolean if a field has been set.

### GetCommType

`func (o *SwitchpointsPutRequestSwitchpointValue) GetCommType() string`

GetCommType returns the CommType field if non-nil, zero value otherwise.

### GetCommTypeOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetCommTypeOk() (*string, bool)`

GetCommTypeOk returns a tuple with the CommType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommType

`func (o *SwitchpointsPutRequestSwitchpointValue) SetCommType(v string)`

SetCommType sets CommType field to given value.

### HasCommType

`func (o *SwitchpointsPutRequestSwitchpointValue) HasCommType() bool`

HasCommType returns a boolean if a field has been set.

### GetSnmpCommunityString

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSnmpCommunityString() string`

GetSnmpCommunityString returns the SnmpCommunityString field if non-nil, zero value otherwise.

### GetSnmpCommunityStringOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSnmpCommunityStringOk() (*string, bool)`

GetSnmpCommunityStringOk returns a tuple with the SnmpCommunityString field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSnmpCommunityString

`func (o *SwitchpointsPutRequestSwitchpointValue) SetSnmpCommunityString(v string)`

SetSnmpCommunityString sets SnmpCommunityString field to given value.

### HasSnmpCommunityString

`func (o *SwitchpointsPutRequestSwitchpointValue) HasSnmpCommunityString() bool`

HasSnmpCommunityString returns a boolean if a field has been set.

### GetUplinkPort

`func (o *SwitchpointsPutRequestSwitchpointValue) GetUplinkPort() string`

GetUplinkPort returns the UplinkPort field if non-nil, zero value otherwise.

### GetUplinkPortOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetUplinkPortOk() (*string, bool)`

GetUplinkPortOk returns a tuple with the UplinkPort field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUplinkPort

`func (o *SwitchpointsPutRequestSwitchpointValue) SetUplinkPort(v string)`

SetUplinkPort sets UplinkPort field to given value.

### HasUplinkPort

`func (o *SwitchpointsPutRequestSwitchpointValue) HasUplinkPort() bool`

HasUplinkPort returns a boolean if a field has been set.

### GetLldpSearchString

`func (o *SwitchpointsPutRequestSwitchpointValue) GetLldpSearchString() string`

GetLldpSearchString returns the LldpSearchString field if non-nil, zero value otherwise.

### GetLldpSearchStringOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetLldpSearchStringOk() (*string, bool)`

GetLldpSearchStringOk returns a tuple with the LldpSearchString field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLldpSearchString

`func (o *SwitchpointsPutRequestSwitchpointValue) SetLldpSearchString(v string)`

SetLldpSearchString sets LldpSearchString field to given value.

### HasLldpSearchString

`func (o *SwitchpointsPutRequestSwitchpointValue) HasLldpSearchString() bool`

HasLldpSearchString returns a boolean if a field has been set.

### GetZtpIdentification

`func (o *SwitchpointsPutRequestSwitchpointValue) GetZtpIdentification() string`

GetZtpIdentification returns the ZtpIdentification field if non-nil, zero value otherwise.

### GetZtpIdentificationOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetZtpIdentificationOk() (*string, bool)`

GetZtpIdentificationOk returns a tuple with the ZtpIdentification field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetZtpIdentification

`func (o *SwitchpointsPutRequestSwitchpointValue) SetZtpIdentification(v string)`

SetZtpIdentification sets ZtpIdentification field to given value.

### HasZtpIdentification

`func (o *SwitchpointsPutRequestSwitchpointValue) HasZtpIdentification() bool`

HasZtpIdentification returns a boolean if a field has been set.

### GetLocatedBy

`func (o *SwitchpointsPutRequestSwitchpointValue) GetLocatedBy() string`

GetLocatedBy returns the LocatedBy field if non-nil, zero value otherwise.

### GetLocatedByOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetLocatedByOk() (*string, bool)`

GetLocatedByOk returns a tuple with the LocatedBy field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLocatedBy

`func (o *SwitchpointsPutRequestSwitchpointValue) SetLocatedBy(v string)`

SetLocatedBy sets LocatedBy field to given value.

### HasLocatedBy

`func (o *SwitchpointsPutRequestSwitchpointValue) HasLocatedBy() bool`

HasLocatedBy returns a boolean if a field has been set.

### GetPowerState

`func (o *SwitchpointsPutRequestSwitchpointValue) GetPowerState() string`

GetPowerState returns the PowerState field if non-nil, zero value otherwise.

### GetPowerStateOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetPowerStateOk() (*string, bool)`

GetPowerStateOk returns a tuple with the PowerState field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPowerState

`func (o *SwitchpointsPutRequestSwitchpointValue) SetPowerState(v string)`

SetPowerState sets PowerState field to given value.

### HasPowerState

`func (o *SwitchpointsPutRequestSwitchpointValue) HasPowerState() bool`

HasPowerState returns a boolean if a field has been set.

### GetCommunicationMode

`func (o *SwitchpointsPutRequestSwitchpointValue) GetCommunicationMode() string`

GetCommunicationMode returns the CommunicationMode field if non-nil, zero value otherwise.

### GetCommunicationModeOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetCommunicationModeOk() (*string, bool)`

GetCommunicationModeOk returns a tuple with the CommunicationMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommunicationMode

`func (o *SwitchpointsPutRequestSwitchpointValue) SetCommunicationMode(v string)`

SetCommunicationMode sets CommunicationMode field to given value.

### HasCommunicationMode

`func (o *SwitchpointsPutRequestSwitchpointValue) HasCommunicationMode() bool`

HasCommunicationMode returns a boolean if a field has been set.

### GetCliAccessMode

`func (o *SwitchpointsPutRequestSwitchpointValue) GetCliAccessMode() string`

GetCliAccessMode returns the CliAccessMode field if non-nil, zero value otherwise.

### GetCliAccessModeOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetCliAccessModeOk() (*string, bool)`

GetCliAccessModeOk returns a tuple with the CliAccessMode field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCliAccessMode

`func (o *SwitchpointsPutRequestSwitchpointValue) SetCliAccessMode(v string)`

SetCliAccessMode sets CliAccessMode field to given value.

### HasCliAccessMode

`func (o *SwitchpointsPutRequestSwitchpointValue) HasCliAccessMode() bool`

HasCliAccessMode returns a boolean if a field has been set.

### GetUsername

`func (o *SwitchpointsPutRequestSwitchpointValue) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsername

`func (o *SwitchpointsPutRequestSwitchpointValue) SetUsername(v string)`

SetUsername sets Username field to given value.

### HasUsername

`func (o *SwitchpointsPutRequestSwitchpointValue) HasUsername() bool`

HasUsername returns a boolean if a field has been set.

### GetPassword

`func (o *SwitchpointsPutRequestSwitchpointValue) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *SwitchpointsPutRequestSwitchpointValue) SetPassword(v string)`

SetPassword sets Password field to given value.

### HasPassword

`func (o *SwitchpointsPutRequestSwitchpointValue) HasPassword() bool`

HasPassword returns a boolean if a field has been set.

### GetEnablePassword

`func (o *SwitchpointsPutRequestSwitchpointValue) GetEnablePassword() string`

GetEnablePassword returns the EnablePassword field if non-nil, zero value otherwise.

### GetEnablePasswordOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetEnablePasswordOk() (*string, bool)`

GetEnablePasswordOk returns a tuple with the EnablePassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnablePassword

`func (o *SwitchpointsPutRequestSwitchpointValue) SetEnablePassword(v string)`

SetEnablePassword sets EnablePassword field to given value.

### HasEnablePassword

`func (o *SwitchpointsPutRequestSwitchpointValue) HasEnablePassword() bool`

HasEnablePassword returns a boolean if a field has been set.

### GetSshKeyOrPassword

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSshKeyOrPassword() string`

GetSshKeyOrPassword returns the SshKeyOrPassword field if non-nil, zero value otherwise.

### GetSshKeyOrPasswordOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSshKeyOrPasswordOk() (*string, bool)`

GetSshKeyOrPasswordOk returns a tuple with the SshKeyOrPassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSshKeyOrPassword

`func (o *SwitchpointsPutRequestSwitchpointValue) SetSshKeyOrPassword(v string)`

SetSshKeyOrPassword sets SshKeyOrPassword field to given value.

### HasSshKeyOrPassword

`func (o *SwitchpointsPutRequestSwitchpointValue) HasSshKeyOrPassword() bool`

HasSshKeyOrPassword returns a boolean if a field has been set.

### GetManagedOnNativeVlan

`func (o *SwitchpointsPutRequestSwitchpointValue) GetManagedOnNativeVlan() bool`

GetManagedOnNativeVlan returns the ManagedOnNativeVlan field if non-nil, zero value otherwise.

### GetManagedOnNativeVlanOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetManagedOnNativeVlanOk() (*bool, bool)`

GetManagedOnNativeVlanOk returns a tuple with the ManagedOnNativeVlan field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetManagedOnNativeVlan

`func (o *SwitchpointsPutRequestSwitchpointValue) SetManagedOnNativeVlan(v bool)`

SetManagedOnNativeVlan sets ManagedOnNativeVlan field to given value.

### HasManagedOnNativeVlan

`func (o *SwitchpointsPutRequestSwitchpointValue) HasManagedOnNativeVlan() bool`

HasManagedOnNativeVlan returns a boolean if a field has been set.

### GetSdlc

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSdlc() string`

GetSdlc returns the Sdlc field if non-nil, zero value otherwise.

### GetSdlcOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSdlcOk() (*string, bool)`

GetSdlcOk returns a tuple with the Sdlc field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSdlc

`func (o *SwitchpointsPutRequestSwitchpointValue) SetSdlc(v string)`

SetSdlc sets Sdlc field to given value.

### HasSdlc

`func (o *SwitchpointsPutRequestSwitchpointValue) HasSdlc() bool`

HasSdlc returns a boolean if a field has been set.

### GetSecurityType

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSecurityType() string`

GetSecurityType returns the SecurityType field if non-nil, zero value otherwise.

### GetSecurityTypeOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSecurityTypeOk() (*string, bool)`

GetSecurityTypeOk returns a tuple with the SecurityType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSecurityType

`func (o *SwitchpointsPutRequestSwitchpointValue) SetSecurityType(v string)`

SetSecurityType sets SecurityType field to given value.

### HasSecurityType

`func (o *SwitchpointsPutRequestSwitchpointValue) HasSecurityType() bool`

HasSecurityType returns a boolean if a field has been set.

### GetSnmpv3Username

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSnmpv3Username() string`

GetSnmpv3Username returns the Snmpv3Username field if non-nil, zero value otherwise.

### GetSnmpv3UsernameOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSnmpv3UsernameOk() (*string, bool)`

GetSnmpv3UsernameOk returns a tuple with the Snmpv3Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSnmpv3Username

`func (o *SwitchpointsPutRequestSwitchpointValue) SetSnmpv3Username(v string)`

SetSnmpv3Username sets Snmpv3Username field to given value.

### HasSnmpv3Username

`func (o *SwitchpointsPutRequestSwitchpointValue) HasSnmpv3Username() bool`

HasSnmpv3Username returns a boolean if a field has been set.

### GetAuthenticationProtocol

`func (o *SwitchpointsPutRequestSwitchpointValue) GetAuthenticationProtocol() string`

GetAuthenticationProtocol returns the AuthenticationProtocol field if non-nil, zero value otherwise.

### GetAuthenticationProtocolOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetAuthenticationProtocolOk() (*string, bool)`

GetAuthenticationProtocolOk returns a tuple with the AuthenticationProtocol field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthenticationProtocol

`func (o *SwitchpointsPutRequestSwitchpointValue) SetAuthenticationProtocol(v string)`

SetAuthenticationProtocol sets AuthenticationProtocol field to given value.

### HasAuthenticationProtocol

`func (o *SwitchpointsPutRequestSwitchpointValue) HasAuthenticationProtocol() bool`

HasAuthenticationProtocol returns a boolean if a field has been set.

### GetPassphrase

`func (o *SwitchpointsPutRequestSwitchpointValue) GetPassphrase() string`

GetPassphrase returns the Passphrase field if non-nil, zero value otherwise.

### GetPassphraseOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetPassphraseOk() (*string, bool)`

GetPassphraseOk returns a tuple with the Passphrase field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassphrase

`func (o *SwitchpointsPutRequestSwitchpointValue) SetPassphrase(v string)`

SetPassphrase sets Passphrase field to given value.

### HasPassphrase

`func (o *SwitchpointsPutRequestSwitchpointValue) HasPassphrase() bool`

HasPassphrase returns a boolean if a field has been set.

### GetPrivateProtocol

`func (o *SwitchpointsPutRequestSwitchpointValue) GetPrivateProtocol() string`

GetPrivateProtocol returns the PrivateProtocol field if non-nil, zero value otherwise.

### GetPrivateProtocolOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetPrivateProtocolOk() (*string, bool)`

GetPrivateProtocolOk returns a tuple with the PrivateProtocol field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrivateProtocol

`func (o *SwitchpointsPutRequestSwitchpointValue) SetPrivateProtocol(v string)`

SetPrivateProtocol sets PrivateProtocol field to given value.

### HasPrivateProtocol

`func (o *SwitchpointsPutRequestSwitchpointValue) HasPrivateProtocol() bool`

HasPrivateProtocol returns a boolean if a field has been set.

### GetPrivatePassword

`func (o *SwitchpointsPutRequestSwitchpointValue) GetPrivatePassword() string`

GetPrivatePassword returns the PrivatePassword field if non-nil, zero value otherwise.

### GetPrivatePasswordOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetPrivatePasswordOk() (*string, bool)`

GetPrivatePasswordOk returns a tuple with the PrivatePassword field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrivatePassword

`func (o *SwitchpointsPutRequestSwitchpointValue) SetPrivatePassword(v string)`

SetPrivatePassword sets PrivatePassword field to given value.

### HasPrivatePassword

`func (o *SwitchpointsPutRequestSwitchpointValue) HasPrivatePassword() bool`

HasPrivatePassword returns a boolean if a field has been set.

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

### GetDeviceManagedAs

`func (o *SwitchpointsPutRequestSwitchpointValue) GetDeviceManagedAs() string`

GetDeviceManagedAs returns the DeviceManagedAs field if non-nil, zero value otherwise.

### GetDeviceManagedAsOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetDeviceManagedAsOk() (*string, bool)`

GetDeviceManagedAsOk returns a tuple with the DeviceManagedAs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceManagedAs

`func (o *SwitchpointsPutRequestSwitchpointValue) SetDeviceManagedAs(v string)`

SetDeviceManagedAs sets DeviceManagedAs field to given value.

### HasDeviceManagedAs

`func (o *SwitchpointsPutRequestSwitchpointValue) HasDeviceManagedAs() bool`

HasDeviceManagedAs returns a boolean if a field has been set.

### GetSwitch

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSwitch() string`

GetSwitch returns the Switch field if non-nil, zero value otherwise.

### GetSwitchOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSwitchOk() (*string, bool)`

GetSwitchOk returns a tuple with the Switch field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitch

`func (o *SwitchpointsPutRequestSwitchpointValue) SetSwitch(v string)`

SetSwitch sets Switch field to given value.

### HasSwitch

`func (o *SwitchpointsPutRequestSwitchpointValue) HasSwitch() bool`

HasSwitch returns a boolean if a field has been set.

### GetSwitchRefType

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSwitchRefType() string`

GetSwitchRefType returns the SwitchRefType field if non-nil, zero value otherwise.

### GetSwitchRefTypeOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetSwitchRefTypeOk() (*string, bool)`

GetSwitchRefTypeOk returns a tuple with the SwitchRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchRefType

`func (o *SwitchpointsPutRequestSwitchpointValue) SetSwitchRefType(v string)`

SetSwitchRefType sets SwitchRefType field to given value.

### HasSwitchRefType

`func (o *SwitchpointsPutRequestSwitchpointValue) HasSwitchRefType() bool`

HasSwitchRefType returns a boolean if a field has been set.

### GetConnectionService

`func (o *SwitchpointsPutRequestSwitchpointValue) GetConnectionService() string`

GetConnectionService returns the ConnectionService field if non-nil, zero value otherwise.

### GetConnectionServiceOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetConnectionServiceOk() (*string, bool)`

GetConnectionServiceOk returns a tuple with the ConnectionService field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnectionService

`func (o *SwitchpointsPutRequestSwitchpointValue) SetConnectionService(v string)`

SetConnectionService sets ConnectionService field to given value.

### HasConnectionService

`func (o *SwitchpointsPutRequestSwitchpointValue) HasConnectionService() bool`

HasConnectionService returns a boolean if a field has been set.

### GetConnectionServiceRefType

`func (o *SwitchpointsPutRequestSwitchpointValue) GetConnectionServiceRefType() string`

GetConnectionServiceRefType returns the ConnectionServiceRefType field if non-nil, zero value otherwise.

### GetConnectionServiceRefTypeOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetConnectionServiceRefTypeOk() (*string, bool)`

GetConnectionServiceRefTypeOk returns a tuple with the ConnectionServiceRefType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetConnectionServiceRefType

`func (o *SwitchpointsPutRequestSwitchpointValue) SetConnectionServiceRefType(v string)`

SetConnectionServiceRefType sets ConnectionServiceRefType field to given value.

### HasConnectionServiceRefType

`func (o *SwitchpointsPutRequestSwitchpointValue) HasConnectionServiceRefType() bool`

HasConnectionServiceRefType returns a boolean if a field has been set.

### GetPort

`func (o *SwitchpointsPutRequestSwitchpointValue) GetPort() string`

GetPort returns the Port field if non-nil, zero value otherwise.

### GetPortOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetPortOk() (*string, bool)`

GetPortOk returns a tuple with the Port field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPort

`func (o *SwitchpointsPutRequestSwitchpointValue) SetPort(v string)`

SetPort sets Port field to given value.

### HasPort

`func (o *SwitchpointsPutRequestSwitchpointValue) HasPort() bool`

HasPort returns a boolean if a field has been set.

### GetUsesTaggedPackets

`func (o *SwitchpointsPutRequestSwitchpointValue) GetUsesTaggedPackets() bool`

GetUsesTaggedPackets returns the UsesTaggedPackets field if non-nil, zero value otherwise.

### GetUsesTaggedPacketsOk

`func (o *SwitchpointsPutRequestSwitchpointValue) GetUsesTaggedPacketsOk() (*bool, bool)`

GetUsesTaggedPacketsOk returns a tuple with the UsesTaggedPackets field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsesTaggedPackets

`func (o *SwitchpointsPutRequestSwitchpointValue) SetUsesTaggedPackets(v bool)`

SetUsesTaggedPackets sets UsesTaggedPackets field to given value.

### HasUsesTaggedPackets

`func (o *SwitchpointsPutRequestSwitchpointValue) HasUsesTaggedPackets() bool`

HasUsesTaggedPackets returns a boolean if a field has been set.

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


