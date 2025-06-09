# ConfigPutRequest

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DeviceController** | Pointer to [**ConfigPutRequestDeviceController**](ConfigPutRequestDeviceController.md) |  | [optional] 
**SfpBreakouts** | Pointer to [**ConfigPutRequestSfpBreakouts**](ConfigPutRequestSfpBreakouts.md) |  | [optional] 
**Site** | Pointer to [**ConfigPutRequestSite**](ConfigPutRequestSite.md) |  | [optional] 
**Ipv4PrefixList** | Pointer to [**ConfigPutRequestIpv4PrefixList**](ConfigPutRequestIpv4PrefixList.md) |  | [optional] 
**EthDeviceProfiles** | Pointer to [**ConfigPutRequestEthDeviceProfiles**](ConfigPutRequestEthDeviceProfiles.md) |  | [optional] 
**Ipv4Filter** | Pointer to [**ConfigPutRequestIpv4Filter**](ConfigPutRequestIpv4Filter.md) |  | [optional] 
**GatewayProfile** | Pointer to [**ConfigPutRequestGatewayProfile**](ConfigPutRequestGatewayProfile.md) |  | [optional] 
**EndpointView** | Pointer to [**ConfigPutRequestEndpointView**](ConfigPutRequestEndpointView.md) |  | [optional] 
**Gateway** | Pointer to [**ConfigPutRequestGateway**](ConfigPutRequestGateway.md) |  | [optional] 
**Ipv6Filter** | Pointer to [**ConfigPutRequestIpv6Filter**](ConfigPutRequestIpv6Filter.md) |  | [optional] 
**RouteMapClause** | Pointer to [**ConfigPutRequestRouteMapClause**](ConfigPutRequestRouteMapClause.md) |  | [optional] 
**Service** | Pointer to [**ConfigPutRequestService**](ConfigPutRequestService.md) |  | [optional] 
**RouteMap** | Pointer to [**ConfigPutRequestRouteMap**](ConfigPutRequestRouteMap.md) |  | [optional] 
**BizdConfig** | Pointer to [**ConfigPutRequestBizdConfig**](ConfigPutRequestBizdConfig.md) |  | [optional] 
**PbEgressProfile** | Pointer to [**ConfigPutRequestPbEgressProfile**](ConfigPutRequestPbEgressProfile.md) |  | [optional] 
**Badge** | Pointer to [**ConfigPutRequestBadge**](ConfigPutRequestBadge.md) |  | [optional] 
**Lag** | Pointer to [**ConfigPutRequestLag**](ConfigPutRequestLag.md) |  | [optional] 
**StaticIp** | Pointer to [**ConfigPutRequestStaticIp**](ConfigPutRequestStaticIp.md) |  | [optional] 
**CommunityList** | Pointer to [**ConfigPutRequestCommunityList**](ConfigPutRequestCommunityList.md) |  | [optional] 
**Tenant** | Pointer to [**ConfigPutRequestTenant**](ConfigPutRequestTenant.md) |  | [optional] 
**Ipv6PrefixList** | Pointer to [**ConfigPutRequestIpv6PrefixList**](ConfigPutRequestIpv6PrefixList.md) |  | [optional] 
**Ipv4ListFilter** | Pointer to [**ConfigPutRequestIpv4ListFilter**](ConfigPutRequestIpv4ListFilter.md) |  | [optional] 
**AsPathAccessList** | Pointer to [**ConfigPutRequestAsPathAccessList**](ConfigPutRequestAsPathAccessList.md) |  | [optional] 
**Ipv6ListFilter** | Pointer to [**ConfigPutRequestIpv6ListFilter**](ConfigPutRequestIpv6ListFilter.md) |  | [optional] 
**EthPortProfile** | Pointer to [**ConfigPutRequestEthPortProfile**](ConfigPutRequestEthPortProfile.md) |  | [optional] 
**EthPortSettings** | Pointer to [**ConfigPutRequestEthPortSettings**](ConfigPutRequestEthPortSettings.md) |  | [optional] 
**EndpointBundle** | Pointer to [**ConfigPutRequestEndpointBundle**](ConfigPutRequestEndpointBundle.md) |  | [optional] 
**StaticConnections** | Pointer to [**ConfigPutRequestStaticConnections**](ConfigPutRequestStaticConnections.md) |  | [optional] 
**Switchpoint** | Pointer to [**ConfigPutRequestSwitchpoint**](ConfigPutRequestSwitchpoint.md) |  | [optional] 
**ExtendedCommunityList** | Pointer to [**ConfigPutRequestExtendedCommunityList**](ConfigPutRequestExtendedCommunityList.md) |  | [optional] 
**ImageUpdateSets** | Pointer to [**ConfigPutRequestImageUpdateSets**](ConfigPutRequestImageUpdateSets.md) |  | [optional] 
**IpFilter** | Pointer to [**ConfigPutRequestIpFilter**](ConfigPutRequestIpFilter.md) |  | [optional] 
**FeatureFlag** | Pointer to [**ConfigPutRequestFeatureFlag**](ConfigPutRequestFeatureFlag.md) |  | [optional] 
**PacketQueue** | Pointer to [**ConfigPutRequestPacketQueue**](ConfigPutRequestPacketQueue.md) |  | [optional] 
**ServicePortProfile** | Pointer to [**ConfigPutRequestServicePortProfile**](ConfigPutRequestServicePortProfile.md) |  | [optional] 
**DeviceVoiceSettings** | Pointer to [**ConfigPutRequestDeviceVoiceSettings**](ConfigPutRequestDeviceVoiceSettings.md) |  | [optional] 
**AuthenticatedEthPort** | Pointer to [**ConfigPutRequestAuthenticatedEthPort**](ConfigPutRequestAuthenticatedEthPort.md) |  | [optional] 
**VoicePortProfiles** | Pointer to [**ConfigPutRequestVoicePortProfiles**](ConfigPutRequestVoicePortProfiles.md) |  | [optional] 
**Endpoint** | Pointer to [**ConfigPutRequestEndpoint**](ConfigPutRequestEndpoint.md) |  | [optional] 
**MacFilter** | Pointer to [**ConfigPutRequestMacFilter**](ConfigPutRequestMacFilter.md) |  | [optional] 

## Methods

### NewConfigPutRequest

`func NewConfigPutRequest() *ConfigPutRequest`

NewConfigPutRequest instantiates a new ConfigPutRequest object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewConfigPutRequestWithDefaults

`func NewConfigPutRequestWithDefaults() *ConfigPutRequest`

NewConfigPutRequestWithDefaults instantiates a new ConfigPutRequest object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDeviceController

`func (o *ConfigPutRequest) GetDeviceController() ConfigPutRequestDeviceController`

GetDeviceController returns the DeviceController field if non-nil, zero value otherwise.

### GetDeviceControllerOk

`func (o *ConfigPutRequest) GetDeviceControllerOk() (*ConfigPutRequestDeviceController, bool)`

GetDeviceControllerOk returns a tuple with the DeviceController field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceController

`func (o *ConfigPutRequest) SetDeviceController(v ConfigPutRequestDeviceController)`

SetDeviceController sets DeviceController field to given value.

### HasDeviceController

`func (o *ConfigPutRequest) HasDeviceController() bool`

HasDeviceController returns a boolean if a field has been set.

### GetSfpBreakouts

`func (o *ConfigPutRequest) GetSfpBreakouts() ConfigPutRequestSfpBreakouts`

GetSfpBreakouts returns the SfpBreakouts field if non-nil, zero value otherwise.

### GetSfpBreakoutsOk

`func (o *ConfigPutRequest) GetSfpBreakoutsOk() (*ConfigPutRequestSfpBreakouts, bool)`

GetSfpBreakoutsOk returns a tuple with the SfpBreakouts field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSfpBreakouts

`func (o *ConfigPutRequest) SetSfpBreakouts(v ConfigPutRequestSfpBreakouts)`

SetSfpBreakouts sets SfpBreakouts field to given value.

### HasSfpBreakouts

`func (o *ConfigPutRequest) HasSfpBreakouts() bool`

HasSfpBreakouts returns a boolean if a field has been set.

### GetSite

`func (o *ConfigPutRequest) GetSite() ConfigPutRequestSite`

GetSite returns the Site field if non-nil, zero value otherwise.

### GetSiteOk

`func (o *ConfigPutRequest) GetSiteOk() (*ConfigPutRequestSite, bool)`

GetSiteOk returns a tuple with the Site field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSite

`func (o *ConfigPutRequest) SetSite(v ConfigPutRequestSite)`

SetSite sets Site field to given value.

### HasSite

`func (o *ConfigPutRequest) HasSite() bool`

HasSite returns a boolean if a field has been set.

### GetIpv4PrefixList

`func (o *ConfigPutRequest) GetIpv4PrefixList() ConfigPutRequestIpv4PrefixList`

GetIpv4PrefixList returns the Ipv4PrefixList field if non-nil, zero value otherwise.

### GetIpv4PrefixListOk

`func (o *ConfigPutRequest) GetIpv4PrefixListOk() (*ConfigPutRequestIpv4PrefixList, bool)`

GetIpv4PrefixListOk returns a tuple with the Ipv4PrefixList field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpv4PrefixList

`func (o *ConfigPutRequest) SetIpv4PrefixList(v ConfigPutRequestIpv4PrefixList)`

SetIpv4PrefixList sets Ipv4PrefixList field to given value.

### HasIpv4PrefixList

`func (o *ConfigPutRequest) HasIpv4PrefixList() bool`

HasIpv4PrefixList returns a boolean if a field has been set.

### GetEthDeviceProfiles

`func (o *ConfigPutRequest) GetEthDeviceProfiles() ConfigPutRequestEthDeviceProfiles`

GetEthDeviceProfiles returns the EthDeviceProfiles field if non-nil, zero value otherwise.

### GetEthDeviceProfilesOk

`func (o *ConfigPutRequest) GetEthDeviceProfilesOk() (*ConfigPutRequestEthDeviceProfiles, bool)`

GetEthDeviceProfilesOk returns a tuple with the EthDeviceProfiles field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEthDeviceProfiles

`func (o *ConfigPutRequest) SetEthDeviceProfiles(v ConfigPutRequestEthDeviceProfiles)`

SetEthDeviceProfiles sets EthDeviceProfiles field to given value.

### HasEthDeviceProfiles

`func (o *ConfigPutRequest) HasEthDeviceProfiles() bool`

HasEthDeviceProfiles returns a boolean if a field has been set.

### GetIpv4Filter

`func (o *ConfigPutRequest) GetIpv4Filter() ConfigPutRequestIpv4Filter`

GetIpv4Filter returns the Ipv4Filter field if non-nil, zero value otherwise.

### GetIpv4FilterOk

`func (o *ConfigPutRequest) GetIpv4FilterOk() (*ConfigPutRequestIpv4Filter, bool)`

GetIpv4FilterOk returns a tuple with the Ipv4Filter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpv4Filter

`func (o *ConfigPutRequest) SetIpv4Filter(v ConfigPutRequestIpv4Filter)`

SetIpv4Filter sets Ipv4Filter field to given value.

### HasIpv4Filter

`func (o *ConfigPutRequest) HasIpv4Filter() bool`

HasIpv4Filter returns a boolean if a field has been set.

### GetGatewayProfile

`func (o *ConfigPutRequest) GetGatewayProfile() ConfigPutRequestGatewayProfile`

GetGatewayProfile returns the GatewayProfile field if non-nil, zero value otherwise.

### GetGatewayProfileOk

`func (o *ConfigPutRequest) GetGatewayProfileOk() (*ConfigPutRequestGatewayProfile, bool)`

GetGatewayProfileOk returns a tuple with the GatewayProfile field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGatewayProfile

`func (o *ConfigPutRequest) SetGatewayProfile(v ConfigPutRequestGatewayProfile)`

SetGatewayProfile sets GatewayProfile field to given value.

### HasGatewayProfile

`func (o *ConfigPutRequest) HasGatewayProfile() bool`

HasGatewayProfile returns a boolean if a field has been set.

### GetEndpointView

`func (o *ConfigPutRequest) GetEndpointView() ConfigPutRequestEndpointView`

GetEndpointView returns the EndpointView field if non-nil, zero value otherwise.

### GetEndpointViewOk

`func (o *ConfigPutRequest) GetEndpointViewOk() (*ConfigPutRequestEndpointView, bool)`

GetEndpointViewOk returns a tuple with the EndpointView field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpointView

`func (o *ConfigPutRequest) SetEndpointView(v ConfigPutRequestEndpointView)`

SetEndpointView sets EndpointView field to given value.

### HasEndpointView

`func (o *ConfigPutRequest) HasEndpointView() bool`

HasEndpointView returns a boolean if a field has been set.

### GetGateway

`func (o *ConfigPutRequest) GetGateway() ConfigPutRequestGateway`

GetGateway returns the Gateway field if non-nil, zero value otherwise.

### GetGatewayOk

`func (o *ConfigPutRequest) GetGatewayOk() (*ConfigPutRequestGateway, bool)`

GetGatewayOk returns a tuple with the Gateway field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGateway

`func (o *ConfigPutRequest) SetGateway(v ConfigPutRequestGateway)`

SetGateway sets Gateway field to given value.

### HasGateway

`func (o *ConfigPutRequest) HasGateway() bool`

HasGateway returns a boolean if a field has been set.

### GetIpv6Filter

`func (o *ConfigPutRequest) GetIpv6Filter() ConfigPutRequestIpv6Filter`

GetIpv6Filter returns the Ipv6Filter field if non-nil, zero value otherwise.

### GetIpv6FilterOk

`func (o *ConfigPutRequest) GetIpv6FilterOk() (*ConfigPutRequestIpv6Filter, bool)`

GetIpv6FilterOk returns a tuple with the Ipv6Filter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpv6Filter

`func (o *ConfigPutRequest) SetIpv6Filter(v ConfigPutRequestIpv6Filter)`

SetIpv6Filter sets Ipv6Filter field to given value.

### HasIpv6Filter

`func (o *ConfigPutRequest) HasIpv6Filter() bool`

HasIpv6Filter returns a boolean if a field has been set.

### GetRouteMapClause

`func (o *ConfigPutRequest) GetRouteMapClause() ConfigPutRequestRouteMapClause`

GetRouteMapClause returns the RouteMapClause field if non-nil, zero value otherwise.

### GetRouteMapClauseOk

`func (o *ConfigPutRequest) GetRouteMapClauseOk() (*ConfigPutRequestRouteMapClause, bool)`

GetRouteMapClauseOk returns a tuple with the RouteMapClause field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRouteMapClause

`func (o *ConfigPutRequest) SetRouteMapClause(v ConfigPutRequestRouteMapClause)`

SetRouteMapClause sets RouteMapClause field to given value.

### HasRouteMapClause

`func (o *ConfigPutRequest) HasRouteMapClause() bool`

HasRouteMapClause returns a boolean if a field has been set.

### GetService

`func (o *ConfigPutRequest) GetService() ConfigPutRequestService`

GetService returns the Service field if non-nil, zero value otherwise.

### GetServiceOk

`func (o *ConfigPutRequest) GetServiceOk() (*ConfigPutRequestService, bool)`

GetServiceOk returns a tuple with the Service field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetService

`func (o *ConfigPutRequest) SetService(v ConfigPutRequestService)`

SetService sets Service field to given value.

### HasService

`func (o *ConfigPutRequest) HasService() bool`

HasService returns a boolean if a field has been set.

### GetRouteMap

`func (o *ConfigPutRequest) GetRouteMap() ConfigPutRequestRouteMap`

GetRouteMap returns the RouteMap field if non-nil, zero value otherwise.

### GetRouteMapOk

`func (o *ConfigPutRequest) GetRouteMapOk() (*ConfigPutRequestRouteMap, bool)`

GetRouteMapOk returns a tuple with the RouteMap field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRouteMap

`func (o *ConfigPutRequest) SetRouteMap(v ConfigPutRequestRouteMap)`

SetRouteMap sets RouteMap field to given value.

### HasRouteMap

`func (o *ConfigPutRequest) HasRouteMap() bool`

HasRouteMap returns a boolean if a field has been set.

### GetBizdConfig

`func (o *ConfigPutRequest) GetBizdConfig() ConfigPutRequestBizdConfig`

GetBizdConfig returns the BizdConfig field if non-nil, zero value otherwise.

### GetBizdConfigOk

`func (o *ConfigPutRequest) GetBizdConfigOk() (*ConfigPutRequestBizdConfig, bool)`

GetBizdConfigOk returns a tuple with the BizdConfig field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBizdConfig

`func (o *ConfigPutRequest) SetBizdConfig(v ConfigPutRequestBizdConfig)`

SetBizdConfig sets BizdConfig field to given value.

### HasBizdConfig

`func (o *ConfigPutRequest) HasBizdConfig() bool`

HasBizdConfig returns a boolean if a field has been set.

### GetPbEgressProfile

`func (o *ConfigPutRequest) GetPbEgressProfile() ConfigPutRequestPbEgressProfile`

GetPbEgressProfile returns the PbEgressProfile field if non-nil, zero value otherwise.

### GetPbEgressProfileOk

`func (o *ConfigPutRequest) GetPbEgressProfileOk() (*ConfigPutRequestPbEgressProfile, bool)`

GetPbEgressProfileOk returns a tuple with the PbEgressProfile field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPbEgressProfile

`func (o *ConfigPutRequest) SetPbEgressProfile(v ConfigPutRequestPbEgressProfile)`

SetPbEgressProfile sets PbEgressProfile field to given value.

### HasPbEgressProfile

`func (o *ConfigPutRequest) HasPbEgressProfile() bool`

HasPbEgressProfile returns a boolean if a field has been set.

### GetBadge

`func (o *ConfigPutRequest) GetBadge() ConfigPutRequestBadge`

GetBadge returns the Badge field if non-nil, zero value otherwise.

### GetBadgeOk

`func (o *ConfigPutRequest) GetBadgeOk() (*ConfigPutRequestBadge, bool)`

GetBadgeOk returns a tuple with the Badge field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBadge

`func (o *ConfigPutRequest) SetBadge(v ConfigPutRequestBadge)`

SetBadge sets Badge field to given value.

### HasBadge

`func (o *ConfigPutRequest) HasBadge() bool`

HasBadge returns a boolean if a field has been set.

### GetLag

`func (o *ConfigPutRequest) GetLag() ConfigPutRequestLag`

GetLag returns the Lag field if non-nil, zero value otherwise.

### GetLagOk

`func (o *ConfigPutRequest) GetLagOk() (*ConfigPutRequestLag, bool)`

GetLagOk returns a tuple with the Lag field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLag

`func (o *ConfigPutRequest) SetLag(v ConfigPutRequestLag)`

SetLag sets Lag field to given value.

### HasLag

`func (o *ConfigPutRequest) HasLag() bool`

HasLag returns a boolean if a field has been set.

### GetStaticIp

`func (o *ConfigPutRequest) GetStaticIp() ConfigPutRequestStaticIp`

GetStaticIp returns the StaticIp field if non-nil, zero value otherwise.

### GetStaticIpOk

`func (o *ConfigPutRequest) GetStaticIpOk() (*ConfigPutRequestStaticIp, bool)`

GetStaticIpOk returns a tuple with the StaticIp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStaticIp

`func (o *ConfigPutRequest) SetStaticIp(v ConfigPutRequestStaticIp)`

SetStaticIp sets StaticIp field to given value.

### HasStaticIp

`func (o *ConfigPutRequest) HasStaticIp() bool`

HasStaticIp returns a boolean if a field has been set.

### GetCommunityList

`func (o *ConfigPutRequest) GetCommunityList() ConfigPutRequestCommunityList`

GetCommunityList returns the CommunityList field if non-nil, zero value otherwise.

### GetCommunityListOk

`func (o *ConfigPutRequest) GetCommunityListOk() (*ConfigPutRequestCommunityList, bool)`

GetCommunityListOk returns a tuple with the CommunityList field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCommunityList

`func (o *ConfigPutRequest) SetCommunityList(v ConfigPutRequestCommunityList)`

SetCommunityList sets CommunityList field to given value.

### HasCommunityList

`func (o *ConfigPutRequest) HasCommunityList() bool`

HasCommunityList returns a boolean if a field has been set.

### GetTenant

`func (o *ConfigPutRequest) GetTenant() ConfigPutRequestTenant`

GetTenant returns the Tenant field if non-nil, zero value otherwise.

### GetTenantOk

`func (o *ConfigPutRequest) GetTenantOk() (*ConfigPutRequestTenant, bool)`

GetTenantOk returns a tuple with the Tenant field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTenant

`func (o *ConfigPutRequest) SetTenant(v ConfigPutRequestTenant)`

SetTenant sets Tenant field to given value.

### HasTenant

`func (o *ConfigPutRequest) HasTenant() bool`

HasTenant returns a boolean if a field has been set.

### GetIpv6PrefixList

`func (o *ConfigPutRequest) GetIpv6PrefixList() ConfigPutRequestIpv6PrefixList`

GetIpv6PrefixList returns the Ipv6PrefixList field if non-nil, zero value otherwise.

### GetIpv6PrefixListOk

`func (o *ConfigPutRequest) GetIpv6PrefixListOk() (*ConfigPutRequestIpv6PrefixList, bool)`

GetIpv6PrefixListOk returns a tuple with the Ipv6PrefixList field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpv6PrefixList

`func (o *ConfigPutRequest) SetIpv6PrefixList(v ConfigPutRequestIpv6PrefixList)`

SetIpv6PrefixList sets Ipv6PrefixList field to given value.

### HasIpv6PrefixList

`func (o *ConfigPutRequest) HasIpv6PrefixList() bool`

HasIpv6PrefixList returns a boolean if a field has been set.

### GetIpv4ListFilter

`func (o *ConfigPutRequest) GetIpv4ListFilter() ConfigPutRequestIpv4ListFilter`

GetIpv4ListFilter returns the Ipv4ListFilter field if non-nil, zero value otherwise.

### GetIpv4ListFilterOk

`func (o *ConfigPutRequest) GetIpv4ListFilterOk() (*ConfigPutRequestIpv4ListFilter, bool)`

GetIpv4ListFilterOk returns a tuple with the Ipv4ListFilter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpv4ListFilter

`func (o *ConfigPutRequest) SetIpv4ListFilter(v ConfigPutRequestIpv4ListFilter)`

SetIpv4ListFilter sets Ipv4ListFilter field to given value.

### HasIpv4ListFilter

`func (o *ConfigPutRequest) HasIpv4ListFilter() bool`

HasIpv4ListFilter returns a boolean if a field has been set.

### GetAsPathAccessList

`func (o *ConfigPutRequest) GetAsPathAccessList() ConfigPutRequestAsPathAccessList`

GetAsPathAccessList returns the AsPathAccessList field if non-nil, zero value otherwise.

### GetAsPathAccessListOk

`func (o *ConfigPutRequest) GetAsPathAccessListOk() (*ConfigPutRequestAsPathAccessList, bool)`

GetAsPathAccessListOk returns a tuple with the AsPathAccessList field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAsPathAccessList

`func (o *ConfigPutRequest) SetAsPathAccessList(v ConfigPutRequestAsPathAccessList)`

SetAsPathAccessList sets AsPathAccessList field to given value.

### HasAsPathAccessList

`func (o *ConfigPutRequest) HasAsPathAccessList() bool`

HasAsPathAccessList returns a boolean if a field has been set.

### GetIpv6ListFilter

`func (o *ConfigPutRequest) GetIpv6ListFilter() ConfigPutRequestIpv6ListFilter`

GetIpv6ListFilter returns the Ipv6ListFilter field if non-nil, zero value otherwise.

### GetIpv6ListFilterOk

`func (o *ConfigPutRequest) GetIpv6ListFilterOk() (*ConfigPutRequestIpv6ListFilter, bool)`

GetIpv6ListFilterOk returns a tuple with the Ipv6ListFilter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpv6ListFilter

`func (o *ConfigPutRequest) SetIpv6ListFilter(v ConfigPutRequestIpv6ListFilter)`

SetIpv6ListFilter sets Ipv6ListFilter field to given value.

### HasIpv6ListFilter

`func (o *ConfigPutRequest) HasIpv6ListFilter() bool`

HasIpv6ListFilter returns a boolean if a field has been set.

### GetEthPortProfile

`func (o *ConfigPutRequest) GetEthPortProfile() ConfigPutRequestEthPortProfile`

GetEthPortProfile returns the EthPortProfile field if non-nil, zero value otherwise.

### GetEthPortProfileOk

`func (o *ConfigPutRequest) GetEthPortProfileOk() (*ConfigPutRequestEthPortProfile, bool)`

GetEthPortProfileOk returns a tuple with the EthPortProfile field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEthPortProfile

`func (o *ConfigPutRequest) SetEthPortProfile(v ConfigPutRequestEthPortProfile)`

SetEthPortProfile sets EthPortProfile field to given value.

### HasEthPortProfile

`func (o *ConfigPutRequest) HasEthPortProfile() bool`

HasEthPortProfile returns a boolean if a field has been set.

### GetEthPortSettings

`func (o *ConfigPutRequest) GetEthPortSettings() ConfigPutRequestEthPortSettings`

GetEthPortSettings returns the EthPortSettings field if non-nil, zero value otherwise.

### GetEthPortSettingsOk

`func (o *ConfigPutRequest) GetEthPortSettingsOk() (*ConfigPutRequestEthPortSettings, bool)`

GetEthPortSettingsOk returns a tuple with the EthPortSettings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEthPortSettings

`func (o *ConfigPutRequest) SetEthPortSettings(v ConfigPutRequestEthPortSettings)`

SetEthPortSettings sets EthPortSettings field to given value.

### HasEthPortSettings

`func (o *ConfigPutRequest) HasEthPortSettings() bool`

HasEthPortSettings returns a boolean if a field has been set.

### GetEndpointBundle

`func (o *ConfigPutRequest) GetEndpointBundle() ConfigPutRequestEndpointBundle`

GetEndpointBundle returns the EndpointBundle field if non-nil, zero value otherwise.

### GetEndpointBundleOk

`func (o *ConfigPutRequest) GetEndpointBundleOk() (*ConfigPutRequestEndpointBundle, bool)`

GetEndpointBundleOk returns a tuple with the EndpointBundle field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpointBundle

`func (o *ConfigPutRequest) SetEndpointBundle(v ConfigPutRequestEndpointBundle)`

SetEndpointBundle sets EndpointBundle field to given value.

### HasEndpointBundle

`func (o *ConfigPutRequest) HasEndpointBundle() bool`

HasEndpointBundle returns a boolean if a field has been set.

### GetStaticConnections

`func (o *ConfigPutRequest) GetStaticConnections() ConfigPutRequestStaticConnections`

GetStaticConnections returns the StaticConnections field if non-nil, zero value otherwise.

### GetStaticConnectionsOk

`func (o *ConfigPutRequest) GetStaticConnectionsOk() (*ConfigPutRequestStaticConnections, bool)`

GetStaticConnectionsOk returns a tuple with the StaticConnections field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStaticConnections

`func (o *ConfigPutRequest) SetStaticConnections(v ConfigPutRequestStaticConnections)`

SetStaticConnections sets StaticConnections field to given value.

### HasStaticConnections

`func (o *ConfigPutRequest) HasStaticConnections() bool`

HasStaticConnections returns a boolean if a field has been set.

### GetSwitchpoint

`func (o *ConfigPutRequest) GetSwitchpoint() ConfigPutRequestSwitchpoint`

GetSwitchpoint returns the Switchpoint field if non-nil, zero value otherwise.

### GetSwitchpointOk

`func (o *ConfigPutRequest) GetSwitchpointOk() (*ConfigPutRequestSwitchpoint, bool)`

GetSwitchpointOk returns a tuple with the Switchpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSwitchpoint

`func (o *ConfigPutRequest) SetSwitchpoint(v ConfigPutRequestSwitchpoint)`

SetSwitchpoint sets Switchpoint field to given value.

### HasSwitchpoint

`func (o *ConfigPutRequest) HasSwitchpoint() bool`

HasSwitchpoint returns a boolean if a field has been set.

### GetExtendedCommunityList

`func (o *ConfigPutRequest) GetExtendedCommunityList() ConfigPutRequestExtendedCommunityList`

GetExtendedCommunityList returns the ExtendedCommunityList field if non-nil, zero value otherwise.

### GetExtendedCommunityListOk

`func (o *ConfigPutRequest) GetExtendedCommunityListOk() (*ConfigPutRequestExtendedCommunityList, bool)`

GetExtendedCommunityListOk returns a tuple with the ExtendedCommunityList field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExtendedCommunityList

`func (o *ConfigPutRequest) SetExtendedCommunityList(v ConfigPutRequestExtendedCommunityList)`

SetExtendedCommunityList sets ExtendedCommunityList field to given value.

### HasExtendedCommunityList

`func (o *ConfigPutRequest) HasExtendedCommunityList() bool`

HasExtendedCommunityList returns a boolean if a field has been set.

### GetImageUpdateSets

`func (o *ConfigPutRequest) GetImageUpdateSets() ConfigPutRequestImageUpdateSets`

GetImageUpdateSets returns the ImageUpdateSets field if non-nil, zero value otherwise.

### GetImageUpdateSetsOk

`func (o *ConfigPutRequest) GetImageUpdateSetsOk() (*ConfigPutRequestImageUpdateSets, bool)`

GetImageUpdateSetsOk returns a tuple with the ImageUpdateSets field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetImageUpdateSets

`func (o *ConfigPutRequest) SetImageUpdateSets(v ConfigPutRequestImageUpdateSets)`

SetImageUpdateSets sets ImageUpdateSets field to given value.

### HasImageUpdateSets

`func (o *ConfigPutRequest) HasImageUpdateSets() bool`

HasImageUpdateSets returns a boolean if a field has been set.

### GetIpFilter

`func (o *ConfigPutRequest) GetIpFilter() ConfigPutRequestIpFilter`

GetIpFilter returns the IpFilter field if non-nil, zero value otherwise.

### GetIpFilterOk

`func (o *ConfigPutRequest) GetIpFilterOk() (*ConfigPutRequestIpFilter, bool)`

GetIpFilterOk returns a tuple with the IpFilter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpFilter

`func (o *ConfigPutRequest) SetIpFilter(v ConfigPutRequestIpFilter)`

SetIpFilter sets IpFilter field to given value.

### HasIpFilter

`func (o *ConfigPutRequest) HasIpFilter() bool`

HasIpFilter returns a boolean if a field has been set.

### GetFeatureFlag

`func (o *ConfigPutRequest) GetFeatureFlag() ConfigPutRequestFeatureFlag`

GetFeatureFlag returns the FeatureFlag field if non-nil, zero value otherwise.

### GetFeatureFlagOk

`func (o *ConfigPutRequest) GetFeatureFlagOk() (*ConfigPutRequestFeatureFlag, bool)`

GetFeatureFlagOk returns a tuple with the FeatureFlag field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFeatureFlag

`func (o *ConfigPutRequest) SetFeatureFlag(v ConfigPutRequestFeatureFlag)`

SetFeatureFlag sets FeatureFlag field to given value.

### HasFeatureFlag

`func (o *ConfigPutRequest) HasFeatureFlag() bool`

HasFeatureFlag returns a boolean if a field has been set.

### GetPacketQueue

`func (o *ConfigPutRequest) GetPacketQueue() ConfigPutRequestPacketQueue`

GetPacketQueue returns the PacketQueue field if non-nil, zero value otherwise.

### GetPacketQueueOk

`func (o *ConfigPutRequest) GetPacketQueueOk() (*ConfigPutRequestPacketQueue, bool)`

GetPacketQueueOk returns a tuple with the PacketQueue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPacketQueue

`func (o *ConfigPutRequest) SetPacketQueue(v ConfigPutRequestPacketQueue)`

SetPacketQueue sets PacketQueue field to given value.

### HasPacketQueue

`func (o *ConfigPutRequest) HasPacketQueue() bool`

HasPacketQueue returns a boolean if a field has been set.

### GetServicePortProfile

`func (o *ConfigPutRequest) GetServicePortProfile() ConfigPutRequestServicePortProfile`

GetServicePortProfile returns the ServicePortProfile field if non-nil, zero value otherwise.

### GetServicePortProfileOk

`func (o *ConfigPutRequest) GetServicePortProfileOk() (*ConfigPutRequestServicePortProfile, bool)`

GetServicePortProfileOk returns a tuple with the ServicePortProfile field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServicePortProfile

`func (o *ConfigPutRequest) SetServicePortProfile(v ConfigPutRequestServicePortProfile)`

SetServicePortProfile sets ServicePortProfile field to given value.

### HasServicePortProfile

`func (o *ConfigPutRequest) HasServicePortProfile() bool`

HasServicePortProfile returns a boolean if a field has been set.

### GetDeviceVoiceSettings

`func (o *ConfigPutRequest) GetDeviceVoiceSettings() ConfigPutRequestDeviceVoiceSettings`

GetDeviceVoiceSettings returns the DeviceVoiceSettings field if non-nil, zero value otherwise.

### GetDeviceVoiceSettingsOk

`func (o *ConfigPutRequest) GetDeviceVoiceSettingsOk() (*ConfigPutRequestDeviceVoiceSettings, bool)`

GetDeviceVoiceSettingsOk returns a tuple with the DeviceVoiceSettings field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDeviceVoiceSettings

`func (o *ConfigPutRequest) SetDeviceVoiceSettings(v ConfigPutRequestDeviceVoiceSettings)`

SetDeviceVoiceSettings sets DeviceVoiceSettings field to given value.

### HasDeviceVoiceSettings

`func (o *ConfigPutRequest) HasDeviceVoiceSettings() bool`

HasDeviceVoiceSettings returns a boolean if a field has been set.

### GetAuthenticatedEthPort

`func (o *ConfigPutRequest) GetAuthenticatedEthPort() ConfigPutRequestAuthenticatedEthPort`

GetAuthenticatedEthPort returns the AuthenticatedEthPort field if non-nil, zero value otherwise.

### GetAuthenticatedEthPortOk

`func (o *ConfigPutRequest) GetAuthenticatedEthPortOk() (*ConfigPutRequestAuthenticatedEthPort, bool)`

GetAuthenticatedEthPortOk returns a tuple with the AuthenticatedEthPort field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthenticatedEthPort

`func (o *ConfigPutRequest) SetAuthenticatedEthPort(v ConfigPutRequestAuthenticatedEthPort)`

SetAuthenticatedEthPort sets AuthenticatedEthPort field to given value.

### HasAuthenticatedEthPort

`func (o *ConfigPutRequest) HasAuthenticatedEthPort() bool`

HasAuthenticatedEthPort returns a boolean if a field has been set.

### GetVoicePortProfiles

`func (o *ConfigPutRequest) GetVoicePortProfiles() ConfigPutRequestVoicePortProfiles`

GetVoicePortProfiles returns the VoicePortProfiles field if non-nil, zero value otherwise.

### GetVoicePortProfilesOk

`func (o *ConfigPutRequest) GetVoicePortProfilesOk() (*ConfigPutRequestVoicePortProfiles, bool)`

GetVoicePortProfilesOk returns a tuple with the VoicePortProfiles field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVoicePortProfiles

`func (o *ConfigPutRequest) SetVoicePortProfiles(v ConfigPutRequestVoicePortProfiles)`

SetVoicePortProfiles sets VoicePortProfiles field to given value.

### HasVoicePortProfiles

`func (o *ConfigPutRequest) HasVoicePortProfiles() bool`

HasVoicePortProfiles returns a boolean if a field has been set.

### GetEndpoint

`func (o *ConfigPutRequest) GetEndpoint() ConfigPutRequestEndpoint`

GetEndpoint returns the Endpoint field if non-nil, zero value otherwise.

### GetEndpointOk

`func (o *ConfigPutRequest) GetEndpointOk() (*ConfigPutRequestEndpoint, bool)`

GetEndpointOk returns a tuple with the Endpoint field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEndpoint

`func (o *ConfigPutRequest) SetEndpoint(v ConfigPutRequestEndpoint)`

SetEndpoint sets Endpoint field to given value.

### HasEndpoint

`func (o *ConfigPutRequest) HasEndpoint() bool`

HasEndpoint returns a boolean if a field has been set.

### GetMacFilter

`func (o *ConfigPutRequest) GetMacFilter() ConfigPutRequestMacFilter`

GetMacFilter returns the MacFilter field if non-nil, zero value otherwise.

### GetMacFilterOk

`func (o *ConfigPutRequest) GetMacFilterOk() (*ConfigPutRequestMacFilter, bool)`

GetMacFilterOk returns a tuple with the MacFilter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMacFilter

`func (o *ConfigPutRequest) SetMacFilter(v ConfigPutRequestMacFilter)`

SetMacFilter sets MacFilter field to given value.

### HasMacFilter

`func (o *ConfigPutRequest) HasMacFilter() bool`

HasMacFilter returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


