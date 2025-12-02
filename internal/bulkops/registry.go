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
//   - HasAutoGen: Boolean indicating whether the resource supports automatic generation
//   - APIClientGetter: Function that returns a ResourceAPIClient for the resource type
var resourceRegistry = map[string]ResourceConfig{
	"gateway": {
		ResourceType:     "gateway",
		PutRequestType:   reflect.TypeOf(openapi.GatewaysPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.GatewaysPutRequest{}),
		HasAutoGen:       false,
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
		HasAutoGen:       false,
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
		HasAutoGen:       true,
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
		HasAutoGen:       true,
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
		HasAutoGen:       false,
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
	"eth_port_profile": {
		ResourceType:     "eth_port_profile",
		PutRequestType:   reflect.TypeOf(openapi.EthportprofilesPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.EthportprofilesPutRequest{}),
		HasAutoGen:       false,
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
		HasAutoGen:       false,
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
		PutRequestType:   reflect.TypeOf(openapi.BundlesPatchRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.BundlesPatchRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "bundle"}
		},
		PutFunc: nil, // Bundles only support PATCH
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.BundlesAPI.BundlesPatch(ctx).BundlesPatchRequest(*req.(*openapi.BundlesPatchRequest)).Execute()
		},
		DeleteFunc: nil, // No DELETE operation
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.BundlesAPI.BundlesGet(ctx).Execute()
		},
	},
}
