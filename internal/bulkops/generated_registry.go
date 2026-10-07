package bulkops

import (
	"context"
	"net/http"
	"reflect"
	"terraform-provider-verity/openapi"
)

var resourceRegistry = map[string]ResourceConfig{
	"acl": {ResourceType: "acl", HeaderSplitKey: "ip_version",
		ResponseCollectionKey: "ipv4_filter",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.AclsPutRequest, values map[string]openapi.AclsPutRequestIpFilterValue) {
				request.IpFilter = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.AclsPutRequest, values map[string]openapi.AclsPutRequestIpFilterValue) {
				request.IpFilter = &values
			})
		},
		PutRequestType:   reflect.TypeOf(openapi.AclsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.AclsPutRequest{}),
		APIClientGetter:  func(c *openapi.APIClient) ResourceAPIClient { return &GenericAPIClient{client: c, resourceType: "acl"} },
		HeaderPutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}, headers map[string]string) (*http.Response, error) {
			request := c.ACLsAPI.AclsPut(ctx).AclsPutRequest(*req.(*openapi.AclsPutRequest))
			if value, ok := headers["ip_version"]; ok {
				request = request.IpVersion(value)
			}
			return request.Execute()
		},
		HeaderPatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}, headers map[string]string) (*http.Response, error) {
			request := c.ACLsAPI.AclsPatch(ctx).AclsPutRequest(*req.(*openapi.AclsPutRequest))
			if value, ok := headers["ip_version"]; ok {
				request = request.IpVersion(value)
			}
			return request.Execute()
		},
		HeaderDeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string, headers map[string]string) (*http.Response, error) {
			request := c.ACLsAPI.AclsDelete(ctx).IpFilterName(names)
			if value, ok := headers["ip_version"]; ok {
				request = request.IpVersion(value)
			}
			return request.Execute()
		},
		HeaderGetFunc: func(c *openapi.APIClient, ctx context.Context, headers map[string]string) (*http.Response, error) {
			request := c.ACLsAPI.AclsGet(ctx)
			if value, ok := headers["ip_version"]; ok {
				request = request.IpVersion(value)
			}
			return request.Execute()
		},
		HeaderResponseExtractor: func(raw map[string]interface{}, headers map[string]string) (map[string]interface{}, error) {
			key := "ipv4_filter"
			switch headers["ip_version"] {
			case "4":
				key = "ipv4_filter"
			case "6":
				key = "ipv6_filter"
			}
			if data, ok := raw[key].(map[string]interface{}); ok {
				return data, nil
			}
			return make(map[string]interface{}), nil
		},
	},
	"as_path_access_list": {ResourceType: "as_path_access_list", HeaderSplitKey: "",
		ResponseCollectionKey: "as_path_access_list",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.AspathaccesslistsPutRequest, values map[string]openapi.AspathaccesslistsPutRequestAsPathAccessListValue) {
				request.AsPathAccessList = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.AspathaccesslistsPutRequest, values map[string]openapi.AspathaccesslistsPutRequestAsPathAccessListValue) {
				request.AsPathAccessList = &values
			})
		},
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
	"authenticated_eth_port": {ResourceType: "authenticated_eth_port", HeaderSplitKey: "",
		ResponseCollectionKey: "authenticated_eth_port",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.AuthenticatedethportsPutRequest, values map[string]openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValue) {
				request.AuthenticatedEthPort = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.AuthenticatedethportsPutRequest, values map[string]openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValue) {
				request.AuthenticatedEthPort = &values
			})
		},
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
	"badge": {ResourceType: "badge", HeaderSplitKey: "",
		ResponseCollectionKey: "badge",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.BadgesPutRequest, values map[string]openapi.BadgesPutRequestBadgeValue) {
				request.Badge = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.BadgesPutRequest, values map[string]openapi.BadgesPutRequestBadgeValue) {
				request.Badge = &values
			})
		},
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
	"bundle": {ResourceType: "bundle", HeaderSplitKey: "",
		ResponseCollectionKey: "endpoint_bundle",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.BundlesPutRequest, values map[string]openapi.BundlesPutRequestEndpointBundleValue) {
				request.EndpointBundle = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.BundlesPutRequest, values map[string]openapi.BundlesPutRequestEndpointBundleValue) {
				request.EndpointBundle = &values
			})
		},
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
	"community_list": {ResourceType: "community_list", HeaderSplitKey: "",
		ResponseCollectionKey: "community_list",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.CommunitylistsPutRequest, values map[string]openapi.CommunitylistsPutRequestCommunityListValue) {
				request.CommunityList = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.CommunitylistsPutRequest, values map[string]openapi.CommunitylistsPutRequestCommunityListValue) {
				request.CommunityList = &values
			})
		},
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
	"device_aaa_profile": {ResourceType: "device_aaa_profile", HeaderSplitKey: "",
		ResponseCollectionKey: "device_aaa_profile",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.DeviceaaaprofilesPutRequest, values map[string]openapi.DeviceaaaprofilesPutRequestDeviceAaaProfileValue) {
				request.DeviceAaaProfile = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.DeviceaaaprofilesPutRequest, values map[string]openapi.DeviceaaaprofilesPutRequestDeviceAaaProfileValue) {
				request.DeviceAaaProfile = &values
			})
		},
		PutRequestType:   reflect.TypeOf(openapi.DeviceaaaprofilesPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.DeviceaaaprofilesPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "device_aaa_profile"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.DeviceAAAProfilesAPI.DeviceaaaprofilesPut(ctx).DeviceaaaprofilesPutRequest(*req.(*openapi.DeviceaaaprofilesPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.DeviceAAAProfilesAPI.DeviceaaaprofilesPatch(ctx).DeviceaaaprofilesPutRequest(*req.(*openapi.DeviceaaaprofilesPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.DeviceAAAProfilesAPI.DeviceaaaprofilesDelete(ctx).DeviceAaaProfileName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.DeviceAAAProfilesAPI.DeviceaaaprofilesGet(ctx).Execute()
		},
	},
	"device_settings": {ResourceType: "device_settings", HeaderSplitKey: "",
		ResponseCollectionKey: "eth_device_profiles",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.DevicesettingsPutRequest, values map[string]openapi.DevicesettingsPutRequestEthDeviceProfilesValue) {
				request.EthDeviceProfiles = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.DevicesettingsPutRequest, values map[string]openapi.DevicesettingsPutRequestEthDeviceProfilesValue) {
				request.EthDeviceProfiles = &values
			})
		},
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
	"device_voice_settings": {ResourceType: "device_voice_settings", HeaderSplitKey: "",
		ResponseCollectionKey: "device_voice_settings",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.DevicevoicesettingsPutRequest, values map[string]openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValue) {
				request.DeviceVoiceSettings = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.DevicevoicesettingsPutRequest, values map[string]openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValue) {
				request.DeviceVoiceSettings = &values
			})
		},
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
	"diagnostics_port_profile": {ResourceType: "diagnostics_port_profile", HeaderSplitKey: "",
		ResponseCollectionKey: "diagnostics_port_profile",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.DiagnosticsportprofilesPutRequest, values map[string]openapi.DiagnosticsportprofilesPutRequestDiagnosticsPortProfileValue) {
				request.DiagnosticsPortProfile = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.DiagnosticsportprofilesPutRequest, values map[string]openapi.DiagnosticsportprofilesPutRequestDiagnosticsPortProfileValue) {
				request.DiagnosticsPortProfile = &values
			})
		},
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
	"diagnostics_profile": {ResourceType: "diagnostics_profile", HeaderSplitKey: "",
		ResponseCollectionKey: "diagnostics_profile",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.DiagnosticsprofilesPutRequest, values map[string]openapi.DiagnosticsprofilesPutRequestDiagnosticsProfileValue) {
				request.DiagnosticsProfile = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.DiagnosticsprofilesPutRequest, values map[string]openapi.DiagnosticsprofilesPutRequestDiagnosticsProfileValue) {
				request.DiagnosticsProfile = &values
			})
		},
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
	"eth_port_profile": {ResourceType: "eth_port_profile", HeaderSplitKey: "",
		ResponseCollectionKey: "eth_port_profile_",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.EthportprofilesPutRequest, values map[string]openapi.EthportprofilesPutRequestEthPortProfileValue) {
				request.EthPortProfile = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.EthportprofilesPutRequest, values map[string]openapi.EthportprofilesPutRequestEthPortProfileValue) {
				request.EthPortProfile = &values
			})
		},
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
	"eth_port_settings": {ResourceType: "eth_port_settings", HeaderSplitKey: "",
		ResponseCollectionKey: "eth_port_settings",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.EthportsettingsPutRequest, values map[string]openapi.EthportsettingsPutRequestEthPortSettingsValue) {
				request.EthPortSettings = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.EthportsettingsPutRequest, values map[string]openapi.EthportsettingsPutRequestEthPortSettingsValue) {
				request.EthPortSettings = &values
			})
		},
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
	"extended_community_list": {ResourceType: "extended_community_list", HeaderSplitKey: "",
		ResponseCollectionKey: "extended_community_list",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.ExtendedcommunitylistsPutRequest, values map[string]openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValue) {
				request.ExtendedCommunityList = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.ExtendedcommunitylistsPutRequest, values map[string]openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValue) {
				request.ExtendedCommunityList = &values
			})
		},
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
	"fabric": {ResourceType: "fabric", HeaderSplitKey: "",
		ResponseCollectionKey: "fabric",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.FabricsPutRequest, values map[string]openapi.FabricsPutRequestFabricValue) {
				request.Fabric = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.FabricsPutRequest, values map[string]openapi.FabricsPutRequestFabricValue) {
				request.Fabric = &values
			})
		},
		PutRequestType:   reflect.TypeOf(openapi.FabricsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.FabricsPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "fabric"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.FabricsAPI.FabricsPut(ctx).FabricsPutRequest(*req.(*openapi.FabricsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.FabricsAPI.FabricsPatch(ctx).FabricsPutRequest(*req.(*openapi.FabricsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.FabricsAPI.FabricsDelete(ctx).FabricName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.FabricsAPI.FabricsGet(ctx).Execute()
		},
	},
	"gateway": {ResourceType: "gateway", HeaderSplitKey: "",
		ResponseCollectionKey: "gateway",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.GatewaysPutRequest, values map[string]openapi.GatewaysPutRequestGatewayValue) {
				request.Gateway = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.GatewaysPutRequest, values map[string]openapi.GatewaysPutRequestGatewayValue) {
				request.Gateway = &values
			})
		},
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
	"gateway_profile": {ResourceType: "gateway_profile", HeaderSplitKey: "",
		ResponseCollectionKey: "gateway_profile",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.GatewayprofilesPutRequest, values map[string]openapi.GatewayprofilesPutRequestGatewayProfileValue) {
				request.GatewayProfile = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.GatewayprofilesPutRequest, values map[string]openapi.GatewayprofilesPutRequestGatewayProfileValue) {
				request.GatewayProfile = &values
			})
		},
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
	"grouping_rule": {ResourceType: "grouping_rule", HeaderSplitKey: "",
		ResponseCollectionKey: "grouping_rules",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.GroupingrulesPutRequest, values map[string]openapi.GroupingrulesPutRequestGroupingRulesValue) {
				request.GroupingRules = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.GroupingrulesPutRequest, values map[string]openapi.GroupingrulesPutRequestGroupingRulesValue) {
				request.GroupingRules = &values
			})
		},
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
	"ipv4_list": {ResourceType: "ipv4_list", HeaderSplitKey: "",
		ResponseCollectionKey: "ipv4_list_filter",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.Ipv4listsPutRequest, values map[string]openapi.Ipv4listsPutRequestIpv4ListFilterValue) {
				request.Ipv4ListFilter = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.Ipv4listsPutRequest, values map[string]openapi.Ipv4listsPutRequestIpv4ListFilterValue) {
				request.Ipv4ListFilter = &values
			})
		},
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
	"ipv4_prefix_list": {ResourceType: "ipv4_prefix_list", HeaderSplitKey: "",
		ResponseCollectionKey: "ipv4_prefix_list",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.Ipv4prefixlistsPutRequest, values map[string]openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValue) {
				request.Ipv4PrefixList = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.Ipv4prefixlistsPutRequest, values map[string]openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValue) {
				request.Ipv4PrefixList = &values
			})
		},
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
	"ipv6_list": {ResourceType: "ipv6_list", HeaderSplitKey: "",
		ResponseCollectionKey: "ipv6_list_filter",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.Ipv6listsPutRequest, values map[string]openapi.Ipv6listsPutRequestIpv6ListFilterValue) {
				request.Ipv6ListFilter = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.Ipv6listsPutRequest, values map[string]openapi.Ipv6listsPutRequestIpv6ListFilterValue) {
				request.Ipv6ListFilter = &values
			})
		},
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
	"ipv6_prefix_list": {ResourceType: "ipv6_prefix_list", HeaderSplitKey: "",
		ResponseCollectionKey: "ipv6_prefix_list",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.Ipv6prefixlistsPutRequest, values map[string]openapi.Ipv6prefixlistsPutRequestIpv6PrefixListValue) {
				request.Ipv6PrefixList = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.Ipv6prefixlistsPutRequest, values map[string]openapi.Ipv6prefixlistsPutRequestIpv6PrefixListValue) {
				request.Ipv6PrefixList = &values
			})
		},
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
	"lag": {ResourceType: "lag", HeaderSplitKey: "",
		ResponseCollectionKey: "lag",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.LagsPutRequest, values map[string]openapi.LagsPutRequestLagValue) {
				request.Lag = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.LagsPutRequest, values map[string]openapi.LagsPutRequestLagValue) {
				request.Lag = &values
			})
		},
		PutRequestType:   reflect.TypeOf(openapi.LagsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.LagsPutRequest{}),
		APIClientGetter:  func(c *openapi.APIClient) ResourceAPIClient { return &GenericAPIClient{client: c, resourceType: "lag"} },
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
	"ldap_profile": {ResourceType: "ldap_profile", HeaderSplitKey: "",
		ResponseCollectionKey: "ldap_profile",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.LdapprofilesPutRequest, values map[string]openapi.LdapprofilesPutRequestLdapProfileValue) {
				request.LdapProfile = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.LdapprofilesPutRequest, values map[string]openapi.LdapprofilesPutRequestLdapProfileValue) {
				request.LdapProfile = &values
			})
		},
		PutRequestType:   reflect.TypeOf(openapi.LdapprofilesPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.LdapprofilesPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "ldap_profile"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.LDAPProfilesAPI.LdapprofilesPut(ctx).LdapprofilesPutRequest(*req.(*openapi.LdapprofilesPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.LDAPProfilesAPI.LdapprofilesPatch(ctx).LdapprofilesPutRequest(*req.(*openapi.LdapprofilesPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.LDAPProfilesAPI.LdapprofilesDelete(ctx).LdapProfileName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.LDAPProfilesAPI.LdapprofilesGet(ctx).Execute()
		},
	},
	"mac_filter": {ResourceType: "mac_filter", HeaderSplitKey: "",
		ResponseCollectionKey: "mac_filter",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.MacfiltersPutRequest, values map[string]openapi.MacfiltersPutRequestMacFilterValue) {
				request.MacFilter = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.MacfiltersPutRequest, values map[string]openapi.MacfiltersPutRequestMacFilterValue) {
				request.MacFilter = &values
			})
		},
		PutRequestType:   reflect.TypeOf(openapi.MacfiltersPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.MacfiltersPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "mac_filter"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.MACFiltersAPI.MacfiltersPut(ctx).MacfiltersPutRequest(*req.(*openapi.MacfiltersPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.MACFiltersAPI.MacfiltersPatch(ctx).MacfiltersPutRequest(*req.(*openapi.MacfiltersPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.MACFiltersAPI.MacfiltersDelete(ctx).MacFilterName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.MACFiltersAPI.MacfiltersGet(ctx).Execute()
		},
	},
	"packet_broker": {ResourceType: "packet_broker", HeaderSplitKey: "",
		ResponseCollectionKey: "pb_egress_profile",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.PacketbrokerPutRequest, values map[string]openapi.PacketbrokerPutRequestPbEgressProfileValue) {
				request.PbEgressProfile = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.PacketbrokerPutRequest, values map[string]openapi.PacketbrokerPutRequestPbEgressProfileValue) {
				request.PbEgressProfile = &values
			})
		},
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
	"packet_queue": {ResourceType: "packet_queue", HeaderSplitKey: "",
		ResponseCollectionKey: "packet_queue",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.PacketqueuesPutRequest, values map[string]openapi.PacketqueuesPutRequestPacketQueueValue) {
				request.PacketQueue = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.PacketqueuesPutRequest, values map[string]openapi.PacketqueuesPutRequestPacketQueueValue) {
				request.PacketQueue = &values
			})
		},
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
	"pair": {ResourceType: "pair", HeaderSplitKey: "",
		ResponseCollectionKey: "switch_pair",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.PairsPutRequest, values map[string]openapi.PairsPutRequestSwitchPairValue) {
				request.SwitchPair = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.PairsPutRequest, values map[string]openapi.PairsPutRequestSwitchPairValue) {
				request.SwitchPair = &values
			})
		},
		PutRequestType:   reflect.TypeOf(openapi.PairsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.PairsPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "pair"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.SwitchPairsAPI.PairsPut(ctx).PairsPutRequest(*req.(*openapi.PairsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.SwitchPairsAPI.PairsPatch(ctx).PairsPutRequest(*req.(*openapi.PairsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.SwitchPairsAPI.PairsDelete(ctx).SwitchPairName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.SwitchPairsAPI.PairsGet(ctx).Execute()
		},
	},
	"pb_routing": {ResourceType: "pb_routing", HeaderSplitKey: "",
		ResponseCollectionKey: "pb_routing",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.PolicybasedroutingPutRequest, values map[string]openapi.PolicybasedroutingPutRequestPbRoutingValue) {
				request.PbRouting = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.PolicybasedroutingPutRequest, values map[string]openapi.PolicybasedroutingPutRequestPbRoutingValue) {
				request.PbRouting = &values
			})
		},
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
	"pb_routing_acl": {ResourceType: "pb_routing_acl", HeaderSplitKey: "",
		ResponseCollectionKey: "pb_routing_acl",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.PolicybasedroutingaclPutRequest, values map[string]openapi.PolicybasedroutingaclPutRequestPbRoutingAclValue) {
				request.PbRoutingAcl = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.PolicybasedroutingaclPutRequest, values map[string]openapi.PolicybasedroutingaclPutRequestPbRoutingAclValue) {
				request.PbRoutingAcl = &values
			})
		},
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
	"plane": {ResourceType: "plane", HeaderSplitKey: "",
		ResponseCollectionKey: "plane",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.PlanesPutRequest, values map[string]openapi.PlanesPutRequestPlaneValue) {
				request.Plane = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.PlanesPutRequest, values map[string]openapi.PlanesPutRequestPlaneValue) {
				request.Plane = &values
			})
		},
		PutRequestType:   reflect.TypeOf(openapi.PlanesPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.PlanesPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "plane"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.PlanesAPI.PlanesPut(ctx).PlanesPutRequest(*req.(*openapi.PlanesPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.PlanesAPI.PlanesPatch(ctx).PlanesPutRequest(*req.(*openapi.PlanesPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.PlanesAPI.PlanesDelete(ctx).PlaneName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.PlanesAPI.PlanesGet(ctx).Execute()
		},
	},
	"pod": {ResourceType: "pod", HeaderSplitKey: "",
		ResponseCollectionKey: "pod",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.PodsPutRequest, values map[string]openapi.PodsPutRequestPodValue) {
				request.Pod = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.PodsPutRequest, values map[string]openapi.PodsPutRequestPodValue) {
				request.Pod = &values
			})
		},
		PutRequestType:   reflect.TypeOf(openapi.PodsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.PodsPutRequest{}),
		APIClientGetter:  func(c *openapi.APIClient) ResourceAPIClient { return &GenericAPIClient{client: c, resourceType: "pod"} },
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
	"port_acl": {ResourceType: "port_acl", HeaderSplitKey: "",
		ResponseCollectionKey: "port_acl",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.PortaclsPutRequest, values map[string]openapi.PortaclsPutRequestPortAclValue) {
				request.PortAcl = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.PortaclsPutRequest, values map[string]openapi.PortaclsPutRequestPortAclValue) {
				request.PortAcl = &values
			})
		},
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
	"rack": {ResourceType: "rack", HeaderSplitKey: "",
		ResponseCollectionKey: "rack",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.RacksPutRequest, values map[string]openapi.RacksPutRequestRackValue) {
				request.Rack = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.RacksPutRequest, values map[string]openapi.RacksPutRequestRackValue) {
				request.Rack = &values
			})
		},
		PutRequestType:   reflect.TypeOf(openapi.RacksPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.RacksPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "rack"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.RacksAPI.RacksPut(ctx).RacksPutRequest(*req.(*openapi.RacksPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.RacksAPI.RacksPatch(ctx).RacksPutRequest(*req.(*openapi.RacksPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.RacksAPI.RacksDelete(ctx).RackName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.RacksAPI.RacksGet(ctx).Execute()
		},
	},
	"route_map": {ResourceType: "route_map", HeaderSplitKey: "",
		ResponseCollectionKey: "route_map",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.RoutemapsPutRequest, values map[string]openapi.RoutemapsPutRequestRouteMapValue) {
				request.RouteMap = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.RoutemapsPutRequest, values map[string]openapi.RoutemapsPutRequestRouteMapValue) {
				request.RouteMap = &values
			})
		},
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
	"route_map_clause": {ResourceType: "route_map_clause", HeaderSplitKey: "",
		ResponseCollectionKey: "route_map_clause",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.RoutemapclausesPutRequest, values map[string]openapi.RoutemapclausesPutRequestRouteMapClauseValue) {
				request.RouteMapClause = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.RoutemapclausesPutRequest, values map[string]openapi.RoutemapclausesPutRequestRouteMapClauseValue) {
				request.RouteMapClause = &values
			})
		},
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
	"service": {ResourceType: "service", HeaderSplitKey: "",
		ResponseCollectionKey: "service",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.ServicesPutRequest, values map[string]openapi.ServicesPutRequestServiceValue) {
				request.Service = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.ServicesPutRequest, values map[string]openapi.ServicesPutRequestServiceValue) {
				request.Service = &values
			})
		},
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
	"service_port_profile": {ResourceType: "service_port_profile", HeaderSplitKey: "",
		ResponseCollectionKey: "service_port_profile",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.ServiceportprofilesPutRequest, values map[string]openapi.ServiceportprofilesPutRequestServicePortProfileValue) {
				request.ServicePortProfile = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.ServiceportprofilesPutRequest, values map[string]openapi.ServiceportprofilesPutRequestServicePortProfileValue) {
				request.ServicePortProfile = &values
			})
		},
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
	"sflow_collector": {ResourceType: "sflow_collector", HeaderSplitKey: "",
		ResponseCollectionKey: "sflow_collector",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.SflowcollectorsPutRequest, values map[string]openapi.SflowcollectorsPutRequestSflowCollectorValue) {
				request.SflowCollector = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.SflowcollectorsPutRequest, values map[string]openapi.SflowcollectorsPutRequestSflowCollectorValue) {
				request.SflowCollector = &values
			})
		},
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
	"sfp_breakout": {ResourceType: "sfp_breakout", HeaderSplitKey: "",
		ResponseCollectionKey: "sfp_breakouts",
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.SfpbreakoutsPatchRequest, values map[string]openapi.SfpbreakoutsPatchRequestSfpBreakoutsValue) {
				request.SfpBreakouts = &values
			})
		},
		PutRequestType:   reflect.TypeOf(openapi.SfpbreakoutsPatchRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.SfpbreakoutsPatchRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "sfp_breakout"}
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.SFPBreakoutsAPI.SfpbreakoutsPatch(ctx).SfpbreakoutsPatchRequest(*req.(*openapi.SfpbreakoutsPatchRequest)).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.SFPBreakoutsAPI.SfpbreakoutsGet(ctx).Execute()
		},
	},
	"spine_plane": {ResourceType: "spine_plane", HeaderSplitKey: "",
		ResponseCollectionKey: "spine_plane",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.SpineplanesPutRequest, values map[string]openapi.SpineplanesPutRequestSpinePlaneValue) {
				request.SpinePlane = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.SpineplanesPutRequest, values map[string]openapi.SpineplanesPutRequestSpinePlaneValue) {
				request.SpinePlane = &values
			})
		},
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
	"ssp_group": {ResourceType: "ssp_group", HeaderSplitKey: "",
		ResponseCollectionKey: "superspine_group",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.SspgroupsPutRequest, values map[string]openapi.SspgroupsPutRequestSuperspineGroupValue) {
				request.SuperspineGroup = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.SspgroupsPutRequest, values map[string]openapi.SspgroupsPutRequestSuperspineGroupValue) {
				request.SuperspineGroup = &values
			})
		},
		PutRequestType:   reflect.TypeOf(openapi.SspgroupsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.SspgroupsPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "ssp_group"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.SuperSpineGroupsAPI.SspgroupsPut(ctx).SspgroupsPutRequest(*req.(*openapi.SspgroupsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.SuperSpineGroupsAPI.SspgroupsPatch(ctx).SspgroupsPutRequest(*req.(*openapi.SspgroupsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.SuperSpineGroupsAPI.SspgroupsDelete(ctx).SuperspineGroupName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.SuperSpineGroupsAPI.SspgroupsGet(ctx).Execute()
		},
	},
	"su": {ResourceType: "su", HeaderSplitKey: "",
		ResponseCollectionKey: "su",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.SusPutRequest, values map[string]openapi.SusPutRequestSuValue) {
				request.Su = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.SusPutRequest, values map[string]openapi.SusPutRequestSuValue) {
				request.Su = &values
			})
		},
		PutRequestType:   reflect.TypeOf(openapi.SusPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.SusPutRequest{}),
		APIClientGetter:  func(c *openapi.APIClient) ResourceAPIClient { return &GenericAPIClient{client: c, resourceType: "su"} },
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.SUsAPI.SusPut(ctx).SusPutRequest(*req.(*openapi.SusPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.SUsAPI.SusPatch(ctx).SusPutRequest(*req.(*openapi.SusPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.SUsAPI.SusDelete(ctx).SuName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.SUsAPI.SusGet(ctx).Execute()
		},
	},
	"switchpoint": {ResourceType: "switchpoint", HeaderSplitKey: "",
		ResponseCollectionKey: "switchpoint",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.SwitchpointsPutRequest, values map[string]openapi.SwitchpointsPutRequestSwitchpointValue) {
				request.Switchpoint = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.SwitchpointsPutRequest, values map[string]openapi.SwitchpointsPutRequestSwitchpointValue) {
				request.Switchpoint = &values
			})
		},
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
	"tacacs_profile": {ResourceType: "tacacs_profile", HeaderSplitKey: "",
		ResponseCollectionKey: "tacacs_profile",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.TacacsprofilesPutRequest, values map[string]openapi.TacacsprofilesPutRequestTacacsProfileValue) {
				request.TacacsProfile = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.TacacsprofilesPutRequest, values map[string]openapi.TacacsprofilesPutRequestTacacsProfileValue) {
				request.TacacsProfile = &values
			})
		},
		PutRequestType:   reflect.TypeOf(openapi.TacacsprofilesPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.TacacsprofilesPutRequest{}),
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "tacacs_profile"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.TACACSProfilesAPI.TacacsprofilesPut(ctx).TacacsprofilesPutRequest(*req.(*openapi.TacacsprofilesPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.TACACSProfilesAPI.TacacsprofilesPatch(ctx).TacacsprofilesPutRequest(*req.(*openapi.TacacsprofilesPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.TACACSProfilesAPI.TacacsprofilesDelete(ctx).TacacsProfileName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.TACACSProfilesAPI.TacacsprofilesGet(ctx).Execute()
		},
	},
	"tenant": {ResourceType: "tenant", HeaderSplitKey: "",
		ResponseCollectionKey: "tenant",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.TenantsPutRequest, values map[string]openapi.TenantsPutRequestTenantValue) {
				request.Tenant = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.TenantsPutRequest, values map[string]openapi.TenantsPutRequestTenantValue) {
				request.Tenant = &values
			})
		},
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
	"threshold": {ResourceType: "threshold", HeaderSplitKey: "",
		ResponseCollectionKey: "threshold",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.ThresholdsPutRequest, values map[string]openapi.ThresholdsPutRequestThresholdValue) {
				request.Threshold = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.ThresholdsPutRequest, values map[string]openapi.ThresholdsPutRequestThresholdValue) {
				request.Threshold = &values
			})
		},
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
	"threshold_group": {ResourceType: "threshold_group", HeaderSplitKey: "",
		ResponseCollectionKey: "threshold_group",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.ThresholdgroupsPutRequest, values map[string]openapi.ThresholdgroupsPutRequestThresholdGroupValue) {
				request.ThresholdGroup = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.ThresholdgroupsPutRequest, values map[string]openapi.ThresholdgroupsPutRequestThresholdGroupValue) {
				request.ThresholdGroup = &values
			})
		},
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
	"voice_port_profile": {ResourceType: "voice_port_profile", HeaderSplitKey: "",
		ResponseCollectionKey: "voice_port_profiles",
		PreparePut: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.VoiceportprofilesPutRequest, values map[string]openapi.VoiceportprofilesPutRequestVoicePortProfilesValue) {
				request.VoicePortProfiles = &values
			})
		},
		PreparePatch: func(data map[string]interface{}) (interface{}, error) {
			return prepareTypedBulkRequest(data, func(request *openapi.VoiceportprofilesPutRequest, values map[string]openapi.VoiceportprofilesPutRequestVoicePortProfilesValue) {
				request.VoicePortProfiles = &values
			})
		},
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
}

var finalCacheRefreshKeys = []string{
	"acls_ipv4",
	"acls_ipv6",
	"as_path_access_lists",
	"authenticated_eth_ports",
	"badges",
	"bundles",
	"community_lists",
	"device_aaa_profiles",
	"device_settings",
	"device_voice_settings",
	"diagnostics_port_profiles",
	"diagnostics_profiles",
	"eth_port_profiles",
	"eth_port_settings",
	"extended_community_lists",
	"fabrics",
	"gateway_profiles",
	"gateways",
	"grouping_rules",
	"ipv4_lists",
	"ipv4_prefix_lists",
	"ipv6_lists",
	"ipv6_prefix_lists",
	"lags",
	"ldap_profiles",
	"mac_filters",
	"packet_brokers",
	"packet_queues",
	"pairs",
	"pb_routing",
	"pb_routing_acl",
	"planes",
	"pods",
	"port_acls",
	"racks",
	"route_map_clauses",
	"route_maps",
	"service_port_profiles",
	"services",
	"sflow_collectors",
	"sfp_breakouts",
	"spine_planes",
	"ssp_groups",
	"sus",
	"switchpoints",
	"tacacs_profiles",
	"tenants",
	"threshold_groups",
	"thresholds",
	"voice_port_profiles",
}

var datacenterPutOrder = []string{
	"community_list",
	"as_path_access_list",
	"ipv6_prefix_list",
	"ipv4_prefix_list",
	"extended_community_list",
	"acl",
	"route_map_clause",
	"pb_routing_acl",
	"route_map",
	"pb_routing",
	"tenant",
	"service",
	"fabric",
	"tacacs_profile",
	"ldap_profile",
	"port_acl",
	"ipv6_list",
	"ipv4_list",
	"pod",
	"packet_queue",
	"device_aaa_profile",
	"eth_port_profile",
	"packet_broker",
	"sflow_collector",
	"gateway",
	"su",
	"diagnostics_port_profile",
	"device_settings",
	"lag",
	"diagnostics_profile",
	"gateway_profile",
	"eth_port_settings",
	"badge",
	"plane",
	"spine_plane",
	"rack",
	"bundle",
	"ssp_group",
	"grouping_rule",
	"switchpoint",
	"threshold",
	"threshold_group",
	"pair",
}

var datacenterPatchOrder = []string{
	"sfp_breakout",
	"community_list",
	"as_path_access_list",
	"ipv6_prefix_list",
	"ipv4_prefix_list",
	"extended_community_list",
	"acl",
	"route_map_clause",
	"pb_routing_acl",
	"route_map",
	"pb_routing",
	"tenant",
	"service",
	"fabric",
	"tacacs_profile",
	"ldap_profile",
	"port_acl",
	"ipv6_list",
	"ipv4_list",
	"pod",
	"packet_queue",
	"device_aaa_profile",
	"eth_port_profile",
	"packet_broker",
	"sflow_collector",
	"gateway",
	"su",
	"diagnostics_port_profile",
	"device_settings",
	"lag",
	"diagnostics_profile",
	"gateway_profile",
	"eth_port_settings",
	"badge",
	"plane",
	"spine_plane",
	"rack",
	"bundle",
	"ssp_group",
	"grouping_rule",
	"switchpoint",
	"threshold",
	"threshold_group",
	"pair",
}

var datacenterDeleteOrder = []string{
	"pair",
	"threshold_group",
	"threshold",
	"switchpoint",
	"grouping_rule",
	"ssp_group",
	"bundle",
	"rack",
	"spine_plane",
	"plane",
	"badge",
	"eth_port_settings",
	"gateway_profile",
	"diagnostics_profile",
	"lag",
	"device_settings",
	"diagnostics_port_profile",
	"su",
	"gateway",
	"sflow_collector",
	"packet_broker",
	"eth_port_profile",
	"device_aaa_profile",
	"packet_queue",
	"pod",
	"ipv4_list",
	"ipv6_list",
	"port_acl",
	"ldap_profile",
	"tacacs_profile",
	"fabric",
	"service",
	"tenant",
	"pb_routing",
	"route_map",
	"pb_routing_acl",
	"route_map_clause",
	"acl",
	"extended_community_list",
	"ipv4_prefix_list",
	"ipv6_prefix_list",
	"as_path_access_list",
	"community_list",
}

var campusPutOrder = []string{
	"acl",
	"mac_filter",
	"service",
	"port_acl",
	"tacacs_profile",
	"ldap_profile",
	"sflow_collector",
	"eth_port_profile",
	"packet_queue",
	"device_aaa_profile",
	"fabric",
	"service_port_profile",
	"diagnostics_profile",
	"authenticated_eth_port",
	"device_settings",
	"voice_port_profile",
	"lag",
	"device_voice_settings",
	"eth_port_settings",
	"diagnostics_port_profile",
	"bundle",
	"badge",
	"grouping_rule",
	"switchpoint",
	"threshold",
	"threshold_group",
	"pair",
}

var campusPatchOrder = []string{
	"sfp_breakout",
	"acl",
	"mac_filter",
	"service",
	"port_acl",
	"tacacs_profile",
	"ldap_profile",
	"sflow_collector",
	"eth_port_profile",
	"packet_queue",
	"device_aaa_profile",
	"fabric",
	"service_port_profile",
	"diagnostics_profile",
	"authenticated_eth_port",
	"device_settings",
	"voice_port_profile",
	"lag",
	"device_voice_settings",
	"eth_port_settings",
	"diagnostics_port_profile",
	"bundle",
	"badge",
	"grouping_rule",
	"switchpoint",
	"threshold",
	"threshold_group",
	"pair",
}

var campusDeleteOrder = []string{
	"pair",
	"threshold_group",
	"threshold",
	"switchpoint",
	"grouping_rule",
	"badge",
	"bundle",
	"diagnostics_port_profile",
	"eth_port_settings",
	"device_voice_settings",
	"lag",
	"voice_port_profile",
	"device_settings",
	"authenticated_eth_port",
	"diagnostics_profile",
	"service_port_profile",
	"fabric",
	"device_aaa_profile",
	"packet_queue",
	"eth_port_profile",
	"sflow_collector",
	"ldap_profile",
	"tacacs_profile",
	"port_acl",
	"service",
	"mac_filter",
	"acl",
}
