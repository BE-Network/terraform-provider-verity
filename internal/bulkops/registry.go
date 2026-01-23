package bulkops

import (
	"context"
	"net/http"
	"reflect"
	"terraform-provider-verity/openapi"
)

// ================================================================================================
// RESOURCE REGISTRY
// ================================================================================================

// resourceRegistry is a mapping of resource types to their configuration details.
// It provides a centralized registry for all resources that can be managed by the Verity provider.
//
// Each entry contains:
//   - ResourceType: String identifier for the resource type
//   - PutRequestType: The reflect.Type for PUT API requests for this resource
//   - PatchRequestType: The reflect.Type for PATCH API requests for this resource
//   - APIClientGetter: Function that returns a ResourceAPIClient for the resource type
var resourceRegistry = map[string]ResourceConfig{
	"gateway": {
		ResourceType:     "gateway",
		PutRequestType:   reflect.TypeOf(openapi.GatewaysPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.GatewaysPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "gateway"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.GatewaysAPI.GatewaysPut(ctx).GatewaysPutRequest(*req.(*openapi.GatewaysPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.GatewaysAPI.GatewaysPatch(ctx).GatewaysPutRequest(*req.(*openapi.GatewaysPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.GatewaysAPI.GatewaysDelete(ctx).GatewayName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.GatewaysAPI.GatewaysGet(ctx).Execute()
		},
	},
	"lag": {
		ResourceType:     "lag",
		PutRequestType:   reflect.TypeOf(openapi.LagsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.LagsPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "lag"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.LAGsAPI.LagsPut(ctx).LagsPutRequest(*req.(*openapi.LagsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.LAGsAPI.LagsPatch(ctx).LagsPutRequest(*req.(*openapi.LagsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.LAGsAPI.LagsDelete(ctx).LagName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.LAGsAPI.LagsGet(ctx).Execute()
		},
	},
	"tenant": {
		ResourceType:     "tenant",
		PutRequestType:   reflect.TypeOf(openapi.TenantsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.TenantsPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "tenant"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.TenantsAPI.TenantsPut(ctx).TenantsPutRequest(*req.(*openapi.TenantsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.TenantsAPI.TenantsPatch(ctx).TenantsPutRequest(*req.(*openapi.TenantsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.TenantsAPI.TenantsDelete(ctx).TenantName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.TenantsAPI.TenantsGet(ctx).Execute()
		},
	},
	"service": {
		ResourceType:     "service",
		PutRequestType:   reflect.TypeOf(openapi.ServicesPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.ServicesPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "service"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.ServicesAPI.ServicesPut(ctx).ServicesPutRequest(*req.(*openapi.ServicesPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.ServicesAPI.ServicesPatch(ctx).ServicesPutRequest(*req.(*openapi.ServicesPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.ServicesAPI.ServicesDelete(ctx).ServiceName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.ServicesAPI.ServicesGet(ctx).Execute()
		},
	},
	"gateway_profile": {
		ResourceType:     "gateway_profile",
		PutRequestType:   reflect.TypeOf(openapi.GatewayprofilesPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.GatewayprofilesPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "gateway_profile"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.GatewayProfilesAPI.GatewayprofilesPut(ctx).GatewayprofilesPutRequest(*req.(*openapi.GatewayprofilesPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.GatewayProfilesAPI.GatewayprofilesPatch(ctx).GatewayprofilesPutRequest(*req.(*openapi.GatewayprofilesPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.GatewayProfilesAPI.GatewayprofilesDelete(ctx).ProfileName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.GatewayProfilesAPI.GatewayprofilesGet(ctx).Execute()
		},
	},
	"grouping_rule": {
		ResourceType:     "grouping_rule",
		PutRequestType:   reflect.TypeOf(openapi.GroupingrulesPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.GroupingrulesPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "grouping_rule"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.GroupingRulesAPI.GroupingrulesPut(ctx).GroupingrulesPutRequest(*req.(*openapi.GroupingrulesPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.GroupingRulesAPI.GroupingrulesPatch(ctx).GroupingrulesPutRequest(*req.(*openapi.GroupingrulesPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.GroupingRulesAPI.GroupingrulesDelete(ctx).GroupingRulesName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.GroupingRulesAPI.GroupingrulesGet(ctx).Execute()
		},
	},
	"threshold_group": {
		ResourceType:     "threshold_group",
		PutRequestType:   reflect.TypeOf(openapi.ThresholdgroupsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.ThresholdgroupsPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "threshold_group"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.ThresholdGroupsAPI.ThresholdgroupsPut(ctx).ThresholdgroupsPutRequest(*req.(*openapi.ThresholdgroupsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.ThresholdGroupsAPI.ThresholdgroupsPatch(ctx).ThresholdgroupsPutRequest(*req.(*openapi.ThresholdgroupsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.ThresholdGroupsAPI.ThresholdgroupsDelete(ctx).ThresholdGroupName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.ThresholdGroupsAPI.ThresholdgroupsGet(ctx).Execute()
		},
	},
	"threshold": {
		ResourceType:     "threshold",
		PutRequestType:   reflect.TypeOf(openapi.ThresholdsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.ThresholdsPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "threshold"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.ThresholdsAPI.ThresholdsPut(ctx).ThresholdsPutRequest(*req.(*openapi.ThresholdsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.ThresholdsAPI.ThresholdsPatch(ctx).ThresholdsPutRequest(*req.(*openapi.ThresholdsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.ThresholdsAPI.ThresholdsDelete(ctx).ThresholdName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.ThresholdsAPI.ThresholdsGet(ctx).Execute()
		},
	},
	"packet_queue": {
		ResourceType:     "packet_queue",
		PutRequestType:   reflect.TypeOf(openapi.PacketqueuesPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.PacketqueuesPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "packet_queue"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.PacketQueuesAPI.PacketqueuesPut(ctx).PacketqueuesPutRequest(*req.(*openapi.PacketqueuesPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.PacketQueuesAPI.PacketqueuesPatch(ctx).PacketqueuesPutRequest(*req.(*openapi.PacketqueuesPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.PacketQueuesAPI.PacketqueuesDelete(ctx).PacketQueueName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.PacketQueuesAPI.PacketqueuesGet(ctx).Execute()
		},
	},
	"eth_port_profile": {
		ResourceType:     "eth_port_profile",
		PutRequestType:   reflect.TypeOf(openapi.EthportprofilesPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.EthportprofilesPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "eth_port_profile"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.EthPortProfilesAPI.EthportprofilesPut(ctx).EthportprofilesPutRequest(*req.(*openapi.EthportprofilesPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.EthPortProfilesAPI.EthportprofilesPatch(ctx).EthportprofilesPutRequest(*req.(*openapi.EthportprofilesPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.EthPortProfilesAPI.EthportprofilesDelete(ctx).ProfileName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.EthPortProfilesAPI.EthportprofilesGet(ctx).Execute()
		},
	},
	"eth_port_settings": {
		ResourceType:     "eth_port_settings",
		PutRequestType:   reflect.TypeOf(openapi.EthportsettingsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.EthportsettingsPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "eth_port_settings"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.EthPortSettingsAPI.EthportsettingsPut(ctx).EthportsettingsPutRequest(*req.(*openapi.EthportsettingsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.EthPortSettingsAPI.EthportsettingsPatch(ctx).EthportsettingsPutRequest(*req.(*openapi.EthportsettingsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.EthPortSettingsAPI.EthportsettingsDelete(ctx).PortName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.EthPortSettingsAPI.EthportsettingsGet(ctx).Execute()
		},
	},
	"bundle": {
		ResourceType:     "bundle",
		PutRequestType:   reflect.TypeOf(openapi.BundlesPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.BundlesPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "bundle"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.BundlesAPI.BundlesPut(ctx).BundlesPutRequest(*req.(*openapi.BundlesPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.BundlesAPI.BundlesPatch(ctx).BundlesPutRequest(*req.(*openapi.BundlesPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.BundlesAPI.BundlesDelete(ctx).BundleName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.BundlesAPI.BundlesGet(ctx).Execute()
		},
	},
	"acl": {
		ResourceType:     "acl",
		PutRequestType:   reflect.TypeOf(openapi.AclsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.AclsPutRequest{}),
		HeaderSplitKey:   "ip_version",
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "acl"}
		},
		HeaderPutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}, headers map[string]string) (*http.Response, error) {
			request := c.ACLsAPI.AclsPut(ctx).AclsPutRequest(*req.(*openapi.AclsPutRequest))
			if ipVersion, ok := headers["ip_version"]; ok {
				request = request.IpVersion(ipVersion)
			}
			return request.Execute()
		},
		HeaderPatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}, headers map[string]string) (*http.Response, error) {
			request := c.ACLsAPI.AclsPatch(ctx).AclsPutRequest(*req.(*openapi.AclsPutRequest))
			if ipVersion, ok := headers["ip_version"]; ok {
				request = request.IpVersion(ipVersion)
			}
			return request.Execute()
		},
		HeaderDeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string, headers map[string]string) (*http.Response, error) {
			request := c.ACLsAPI.AclsDelete(ctx).IpFilterName(names)
			if ipVersion, ok := headers["ip_version"]; ok {
				request = request.IpVersion(ipVersion)
			}
			return request.Execute()
		},
		HeaderGetFunc: func(c *openapi.APIClient, ctx context.Context, headers map[string]string) (*http.Response, error) {
			request := c.ACLsAPI.AclsGet(ctx)
			if ipVersion, ok := headers["ip_version"]; ok {
				request = request.IpVersion(ipVersion)
			}
			return request.Execute()
		},
		HeaderResponseExtractor: func(rawResponse map[string]interface{}, headers map[string]string) (map[string]interface{}, error) {
			// ACL response has different field names based on IP version
			var filterKey string
			if headers["ip_version"] == "6" {
				filterKey = "ipv6_filter"
			} else {
				filterKey = "ipv4_filter"
			}
			if ipFilter, ok := rawResponse[filterKey].(map[string]interface{}); ok {
				return ipFilter, nil
			}
			return make(map[string]interface{}), nil
		},
	},
	"packet_broker": {
		ResourceType:     "packet_broker",
		PutRequestType:   reflect.TypeOf(openapi.PacketbrokerPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.PacketbrokerPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "packet_broker"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.PacketBrokerAPI.PacketbrokerPut(ctx).PacketbrokerPutRequest(*req.(*openapi.PacketbrokerPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.PacketBrokerAPI.PacketbrokerPatch(ctx).PacketbrokerPutRequest(*req.(*openapi.PacketbrokerPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.PacketBrokerAPI.PacketbrokerDelete(ctx).PbEgressProfileName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.PacketBrokerAPI.PacketbrokerGet(ctx).Execute()
		},
	},
	"badge": {
		ResourceType:     "badge",
		PutRequestType:   reflect.TypeOf(openapi.BadgesPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.BadgesPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "badge"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.BadgesAPI.BadgesPut(ctx).BadgesPutRequest(*req.(*openapi.BadgesPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.BadgesAPI.BadgesPatch(ctx).BadgesPutRequest(*req.(*openapi.BadgesPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.BadgesAPI.BadgesDelete(ctx).BadgeName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.BadgesAPI.BadgesGet(ctx).Execute()
		},
	},
	"switchpoint": {
		ResourceType:     "switchpoint",
		PutRequestType:   reflect.TypeOf(openapi.SwitchpointsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.SwitchpointsPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "switchpoint"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.SwitchpointsAPI.SwitchpointsPut(ctx).SwitchpointsPutRequest(*req.(*openapi.SwitchpointsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.SwitchpointsAPI.SwitchpointsPatch(ctx).SwitchpointsPutRequest(*req.(*openapi.SwitchpointsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.SwitchpointsAPI.SwitchpointsDelete(ctx).SwitchpointName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.SwitchpointsAPI.SwitchpointsGet(ctx).Execute()
		},
	},
	"device_controller": {
		ResourceType:     "device_controller",
		PutRequestType:   reflect.TypeOf(openapi.DevicecontrollersPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.DevicecontrollersPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "device_controller"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.DeviceControllersAPI.DevicecontrollersPut(ctx).DevicecontrollersPutRequest(*req.(*openapi.DevicecontrollersPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.DeviceControllersAPI.DevicecontrollersPatch(ctx).DevicecontrollersPutRequest(*req.(*openapi.DevicecontrollersPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.DeviceControllersAPI.DevicecontrollersDelete(ctx).DeviceControllerName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.DeviceControllersAPI.DevicecontrollersGet(ctx).Execute()
		},
	},
	"authenticated_eth_port": {
		ResourceType:     "authenticated_eth_port",
		PutRequestType:   reflect.TypeOf(openapi.AuthenticatedethportsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.AuthenticatedethportsPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "authenticated_eth_port"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.AuthenticatedEthPortsAPI.AuthenticatedethportsPut(ctx).AuthenticatedethportsPutRequest(*req.(*openapi.AuthenticatedethportsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.AuthenticatedEthPortsAPI.AuthenticatedethportsPatch(ctx).AuthenticatedethportsPutRequest(*req.(*openapi.AuthenticatedethportsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.AuthenticatedEthPortsAPI.AuthenticatedethportsDelete(ctx).AuthenticatedEthPortName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.AuthenticatedEthPortsAPI.AuthenticatedethportsGet(ctx).Execute()
		},
	},
	"device_voice_settings": {
		ResourceType:     "device_voice_settings",
		PutRequestType:   reflect.TypeOf(openapi.DevicevoicesettingsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.DevicevoicesettingsPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "device_voice_settings"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.DeviceVoiceSettingsAPI.DevicevoicesettingsPut(ctx).DevicevoicesettingsPutRequest(*req.(*openapi.DevicevoicesettingsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.DeviceVoiceSettingsAPI.DevicevoicesettingsPatch(ctx).DevicevoicesettingsPutRequest(*req.(*openapi.DevicevoicesettingsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.DeviceVoiceSettingsAPI.DevicevoicesettingsDelete(ctx).DeviceVoiceSettingsName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.DeviceVoiceSettingsAPI.DevicevoicesettingsGet(ctx).Execute()
		},
	},
	"voice_port_profile": {
		ResourceType:     "voice_port_profile",
		PutRequestType:   reflect.TypeOf(openapi.VoiceportprofilesPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.VoiceportprofilesPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "voice_port_profile"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.VoicePortProfilesAPI.VoiceportprofilesPut(ctx).VoiceportprofilesPutRequest(*req.(*openapi.VoiceportprofilesPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.VoicePortProfilesAPI.VoiceportprofilesPatch(ctx).VoiceportprofilesPutRequest(*req.(*openapi.VoiceportprofilesPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.VoicePortProfilesAPI.VoiceportprofilesDelete(ctx).VoicePortProfileName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.VoicePortProfilesAPI.VoiceportprofilesGet(ctx).Execute()
		},
	},
	"service_port_profile": {
		ResourceType:     "service_port_profile",
		PutRequestType:   reflect.TypeOf(openapi.ServiceportprofilesPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.ServiceportprofilesPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "service_port_profile"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.ServicePortProfilesAPI.ServiceportprofilesPut(ctx).ServiceportprofilesPutRequest(*req.(*openapi.ServiceportprofilesPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.ServicePortProfilesAPI.ServiceportprofilesPatch(ctx).ServiceportprofilesPutRequest(*req.(*openapi.ServiceportprofilesPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.ServicePortProfilesAPI.ServiceportprofilesDelete(ctx).ServicePortProfileName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.ServicePortProfilesAPI.ServiceportprofilesGet(ctx).Execute()
		},
	},
	"as_path_access_list": {
		ResourceType:     "as_path_access_list",
		PutRequestType:   reflect.TypeOf(openapi.AspathaccesslistsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.AspathaccesslistsPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "as_path_access_list"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.ASPathAccessListsAPI.AspathaccesslistsPut(ctx).AspathaccesslistsPutRequest(*req.(*openapi.AspathaccesslistsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.ASPathAccessListsAPI.AspathaccesslistsPatch(ctx).AspathaccesslistsPutRequest(*req.(*openapi.AspathaccesslistsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.ASPathAccessListsAPI.AspathaccesslistsDelete(ctx).AsPathAccessListName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.ASPathAccessListsAPI.AspathaccesslistsGet(ctx).Execute()
		},
	},
	"community_list": {
		ResourceType:     "community_list",
		PutRequestType:   reflect.TypeOf(openapi.CommunitylistsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.CommunitylistsPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "community_list"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.CommunityListsAPI.CommunitylistsPut(ctx).CommunitylistsPutRequest(*req.(*openapi.CommunitylistsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.CommunityListsAPI.CommunitylistsPatch(ctx).CommunitylistsPutRequest(*req.(*openapi.CommunitylistsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.CommunityListsAPI.CommunitylistsDelete(ctx).CommunityListName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.CommunityListsAPI.CommunitylistsGet(ctx).Execute()
		},
	},
	"device_settings": {
		ResourceType:     "device_settings",
		PutRequestType:   reflect.TypeOf(openapi.DevicesettingsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.DevicesettingsPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "device_settings"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.DeviceSettingsAPI.DevicesettingsPut(ctx).DevicesettingsPutRequest(*req.(*openapi.DevicesettingsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.DeviceSettingsAPI.DevicesettingsPatch(ctx).DevicesettingsPutRequest(*req.(*openapi.DevicesettingsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.DeviceSettingsAPI.DevicesettingsDelete(ctx).EthDeviceProfilesName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.DeviceSettingsAPI.DevicesettingsGet(ctx).Execute()
		},
	},
	"extended_community_list": {
		ResourceType:     "extended_community_list",
		PutRequestType:   reflect.TypeOf(openapi.ExtendedcommunitylistsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.ExtendedcommunitylistsPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "extended_community_list"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.ExtendedCommunityListsAPI.ExtendedcommunitylistsPut(ctx).ExtendedcommunitylistsPutRequest(*req.(*openapi.ExtendedcommunitylistsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.ExtendedCommunityListsAPI.ExtendedcommunitylistsPatch(ctx).ExtendedcommunitylistsPutRequest(*req.(*openapi.ExtendedcommunitylistsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.ExtendedCommunityListsAPI.ExtendedcommunitylistsDelete(ctx).ExtendedCommunityListName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.ExtendedCommunityListsAPI.ExtendedcommunitylistsGet(ctx).Execute()
		},
	},
	"ipv4_list": {
		ResourceType:     "ipv4_list",
		PutRequestType:   reflect.TypeOf(openapi.Ipv4listsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.Ipv4listsPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "ipv4_list"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.IPv4ListFiltersAPI.Ipv4listsPut(ctx).Ipv4listsPutRequest(*req.(*openapi.Ipv4listsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.IPv4ListFiltersAPI.Ipv4listsPatch(ctx).Ipv4listsPutRequest(*req.(*openapi.Ipv4listsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.IPv4ListFiltersAPI.Ipv4listsDelete(ctx).Ipv4ListFilterName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.IPv4ListFiltersAPI.Ipv4listsGet(ctx).Execute()
		},
	},
	"ipv4_prefix_list": {
		ResourceType:     "ipv4_prefix_list",
		PutRequestType:   reflect.TypeOf(openapi.Ipv4prefixlistsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.Ipv4prefixlistsPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "ipv4_prefix_list"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.IPv4PrefixListsAPI.Ipv4prefixlistsPut(ctx).Ipv4prefixlistsPutRequest(*req.(*openapi.Ipv4prefixlistsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.IPv4PrefixListsAPI.Ipv4prefixlistsPatch(ctx).Ipv4prefixlistsPutRequest(*req.(*openapi.Ipv4prefixlistsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.IPv4PrefixListsAPI.Ipv4prefixlistsDelete(ctx).Ipv4PrefixListName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.IPv4PrefixListsAPI.Ipv4prefixlistsGet(ctx).Execute()
		},
	},
	"ipv6_list": {
		ResourceType:     "ipv6_list",
		PutRequestType:   reflect.TypeOf(openapi.Ipv6listsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.Ipv6listsPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "ipv6_list"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.IPv6ListFiltersAPI.Ipv6listsPut(ctx).Ipv6listsPutRequest(*req.(*openapi.Ipv6listsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.IPv6ListFiltersAPI.Ipv6listsPatch(ctx).Ipv6listsPutRequest(*req.(*openapi.Ipv6listsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.IPv6ListFiltersAPI.Ipv6listsDelete(ctx).Ipv6ListFilterName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.IPv6ListFiltersAPI.Ipv6listsGet(ctx).Execute()
		},
	},
	"ipv6_prefix_list": {
		ResourceType:     "ipv6_prefix_list",
		PutRequestType:   reflect.TypeOf(openapi.Ipv6prefixlistsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.Ipv6prefixlistsPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "ipv6_prefix_list"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.IPv6PrefixListsAPI.Ipv6prefixlistsPut(ctx).Ipv6prefixlistsPutRequest(*req.(*openapi.Ipv6prefixlistsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.IPv6PrefixListsAPI.Ipv6prefixlistsPatch(ctx).Ipv6prefixlistsPutRequest(*req.(*openapi.Ipv6prefixlistsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.IPv6PrefixListsAPI.Ipv6prefixlistsDelete(ctx).Ipv6PrefixListName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.IPv6PrefixListsAPI.Ipv6prefixlistsGet(ctx).Execute()
		},
	},
	"route_map_clause": {
		ResourceType:     "route_map_clause",
		PutRequestType:   reflect.TypeOf(openapi.RoutemapclausesPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.RoutemapclausesPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "route_map_clause"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.RouteMapClausesAPI.RoutemapclausesPut(ctx).RoutemapclausesPutRequest(*req.(*openapi.RoutemapclausesPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.RouteMapClausesAPI.RoutemapclausesPatch(ctx).RoutemapclausesPutRequest(*req.(*openapi.RoutemapclausesPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.RouteMapClausesAPI.RoutemapclausesDelete(ctx).RouteMapClauseName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.RouteMapClausesAPI.RoutemapclausesGet(ctx).Execute()
		},
	},
	"route_map": {
		ResourceType:     "route_map",
		PutRequestType:   reflect.TypeOf(openapi.RoutemapsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.RoutemapsPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "route_map"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.RouteMapsAPI.RoutemapsPut(ctx).RoutemapsPutRequest(*req.(*openapi.RoutemapsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.RouteMapsAPI.RoutemapsPatch(ctx).RoutemapsPutRequest(*req.(*openapi.RoutemapsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.RouteMapsAPI.RoutemapsDelete(ctx).RouteMapName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.RouteMapsAPI.RoutemapsGet(ctx).Execute()
		},
	},
	"sfp_breakout": {
		ResourceType:     "sfp_breakout",
		PutRequestType:   reflect.TypeOf(openapi.SfpbreakoutsPatchRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.SfpbreakoutsPatchRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "sfp_breakout"}
		},
		PutFunc: nil, // SFP Breakouts only support PATCH
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.SFPBreakoutsAPI.SfpbreakoutsPatch(ctx).SfpbreakoutsPatchRequest(*req.(*openapi.SfpbreakoutsPatchRequest)).Execute()
		},
		DeleteFunc: nil, // No DELETE operation
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.SFPBreakoutsAPI.SfpbreakoutsGet(ctx).Execute()
		},
	},
	"site": {
		ResourceType:     "site",
		PutRequestType:   reflect.TypeOf(openapi.SitesPatchRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.SitesPatchRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "site"}
		},
		PutFunc: nil, // Sites only support PATCH
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.SitesAPI.SitesPatch(ctx).SitesPatchRequest(*req.(*openapi.SitesPatchRequest)).Execute()
		},
		DeleteFunc: nil, // No DELETE operation
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.SitesAPI.SitesGet(ctx).Execute()
		},
	},
	"pod": {
		ResourceType:     "pod",
		PutRequestType:   reflect.TypeOf(openapi.PodsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.PodsPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "pod"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.PodsAPI.PodsPut(ctx).PodsPutRequest(*req.(*openapi.PodsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.PodsAPI.PodsPatch(ctx).PodsPutRequest(*req.(*openapi.PodsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.PodsAPI.PodsDelete(ctx).PodName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.PodsAPI.PodsGet(ctx).Execute()
		},
	},
	"port_acl": {
		ResourceType:     "port_acl",
		PutRequestType:   reflect.TypeOf(openapi.PortaclsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.PortaclsPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "port_acl"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.PortACLsAPI.PortaclsPut(ctx).PortaclsPutRequest(*req.(*openapi.PortaclsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.PortACLsAPI.PortaclsPatch(ctx).PortaclsPutRequest(*req.(*openapi.PortaclsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.PortACLsAPI.PortaclsDelete(ctx).PortAclName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.PortACLsAPI.PortaclsGet(ctx).Execute()
		},
	},
	"sflow_collector": {
		ResourceType:     "sflow_collector",
		PutRequestType:   reflect.TypeOf(openapi.SflowcollectorsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.SflowcollectorsPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "sflow_collector"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.SFlowCollectorsAPI.SflowcollectorsPut(ctx).SflowcollectorsPutRequest(*req.(*openapi.SflowcollectorsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.SFlowCollectorsAPI.SflowcollectorsPatch(ctx).SflowcollectorsPutRequest(*req.(*openapi.SflowcollectorsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.SFlowCollectorsAPI.SflowcollectorsDelete(ctx).SflowCollectorName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.SFlowCollectorsAPI.SflowcollectorsGet(ctx).Execute()
		},
	},
	"diagnostics_profile": {
		ResourceType:     "diagnostics_profile",
		PutRequestType:   reflect.TypeOf(openapi.DiagnosticsprofilesPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.DiagnosticsprofilesPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "diagnostics_profile"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.DiagnosticsProfilesAPI.DiagnosticsprofilesPut(ctx).DiagnosticsprofilesPutRequest(*req.(*openapi.DiagnosticsprofilesPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.DiagnosticsProfilesAPI.DiagnosticsprofilesPatch(ctx).DiagnosticsprofilesPutRequest(*req.(*openapi.DiagnosticsprofilesPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.DiagnosticsProfilesAPI.DiagnosticsprofilesDelete(ctx).DiagnosticsProfileName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.DiagnosticsProfilesAPI.DiagnosticsprofilesGet(ctx).Execute()
		},
	},
	"diagnostics_port_profile": {
		ResourceType:     "diagnostics_port_profile",
		PutRequestType:   reflect.TypeOf(openapi.DiagnosticsportprofilesPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.DiagnosticsportprofilesPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "diagnostics_port_profile"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.DiagnosticsPortProfilesAPI.DiagnosticsportprofilesPut(ctx).DiagnosticsportprofilesPutRequest(*req.(*openapi.DiagnosticsportprofilesPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.DiagnosticsPortProfilesAPI.DiagnosticsportprofilesPatch(ctx).DiagnosticsportprofilesPutRequest(*req.(*openapi.DiagnosticsportprofilesPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.DiagnosticsPortProfilesAPI.DiagnosticsportprofilesDelete(ctx).DiagnosticsPortProfileName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.DiagnosticsPortProfilesAPI.DiagnosticsportprofilesGet(ctx).Execute()
		},
	},
	"pb_routing": {
		ResourceType:     "pb_routing",
		PutRequestType:   reflect.TypeOf(openapi.PolicybasedroutingPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.PolicybasedroutingPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "pb_routing"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.PBRoutingAPI.PolicybasedroutingPut(ctx).PolicybasedroutingPutRequest(*req.(*openapi.PolicybasedroutingPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.PBRoutingAPI.PolicybasedroutingPatch(ctx).PolicybasedroutingPutRequest(*req.(*openapi.PolicybasedroutingPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.PBRoutingAPI.PolicybasedroutingDelete(ctx).PbRoutingName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.PBRoutingAPI.PolicybasedroutingGet(ctx).Execute()
		},
	},
	"pb_routing_acl": {
		ResourceType:     "pb_routing_acl",
		PutRequestType:   reflect.TypeOf(openapi.PolicybasedroutingaclPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.PolicybasedroutingaclPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "pb_routing_acl"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.PBRoutingACLAPI.PolicybasedroutingaclPut(ctx).PolicybasedroutingaclPutRequest(*req.(*openapi.PolicybasedroutingaclPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.PBRoutingACLAPI.PolicybasedroutingaclPatch(ctx).PolicybasedroutingaclPutRequest(*req.(*openapi.PolicybasedroutingaclPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.PBRoutingACLAPI.PolicybasedroutingaclDelete(ctx).PbRoutingAclName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.PBRoutingACLAPI.PolicybasedroutingaclGet(ctx).Execute()
		},
	},
	"spine_plane": {
		ResourceType:     "spine_plane",
		PutRequestType:   reflect.TypeOf(openapi.SpineplanesPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.SpineplanesPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "spine_plane"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.SpinePlanesAPI.SpineplanesPut(ctx).SpineplanesPutRequest(*req.(*openapi.SpineplanesPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.SpinePlanesAPI.SpineplanesPatch(ctx).SpineplanesPutRequest(*req.(*openapi.SpineplanesPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.SpinePlanesAPI.SpineplanesDelete(ctx).SpinePlaneName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.SpinePlanesAPI.SpineplanesGet(ctx).Execute()
		},
	},
}
