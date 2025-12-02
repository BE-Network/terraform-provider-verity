package bulkops

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"terraform-provider-verity/openapi"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (m *Manager) storeOperation(resourceType, resourceName, operationType string, props interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	res, exists := m.resources[resourceType]
	if !exists {
		return
	}

	switch operationType {
	case "PUT":
		if res.Put == nil {
			res.Put = make(map[string]interface{})
		}
		res.Put[resourceName] = props
	case "PATCH":
		if res.Patch == nil {
			res.Patch = make(map[string]interface{})
		}
		res.Patch[resourceName] = props
	case "DELETE":
		res.Delete = append(res.Delete, resourceName)
	}

	res.RecentOps = true
	res.RecentOpTime = time.Now()
}

func (m *Manager) getBatchSize(resourceType, operationType string) int {
	data := m.GetResourceOperationData(resourceType)
	if data == nil {
		return 0
	}

	switch operationType {
	case "PUT":
		if v := reflect.ValueOf(data.PutOperations); v.IsValid() && !v.IsNil() {
			return v.Len()
		}
	case "PATCH":
		if v := reflect.ValueOf(data.PatchOperations); v.IsValid() && !v.IsNil() {
			return v.Len()
		}
	case "DELETE":
		if data.DeleteOperations != nil {
			return len(*data.DeleteOperations)
		}
		return 0
	}
	return 0
}

// Dynamic function factories that create the appropriate behavior for each resource
func (m *Manager) createExtractor(resourceType, operationType string) func() (map[string]interface{}, []string) {
	return func() (map[string]interface{}, []string) {
		m.mutex.Lock()
		defer m.mutex.Unlock()

		res, exists := m.resources[resourceType]
		if !exists {
			return make(map[string]interface{}), []string{}
		}

		switch operationType {
		case "PUT":
			if res.Put == nil {
				return make(map[string]interface{}), []string{}
			}
			originalOperations := make(map[string]interface{})
			names := make([]string, 0, len(res.Put))
			for k, v := range res.Put {
				originalOperations[k] = v
				names = append(names, k)
			}
			// Clear the unified structure
			res.Put = make(map[string]interface{})
			return originalOperations, names

		case "PATCH":
			if res.Patch == nil {
				return make(map[string]interface{}), []string{}
			}
			originalOperations := make(map[string]interface{})
			names := make([]string, 0, len(res.Patch))
			for k, v := range res.Patch {
				originalOperations[k] = v
				names = append(names, k)
			}
			// Clear the unified structure
			res.Patch = make(map[string]interface{})
			return originalOperations, names

		case "DELETE":
			if len(res.Delete) == 0 {
				return make(map[string]interface{}), []string{}
			}
			names := make([]string, len(res.Delete))
			copy(names, res.Delete)
			result := make(map[string]interface{})
			for _, name := range names {
				result[name] = true
			}
			// Clear the unified structure
			res.Delete = res.Delete[:0]
			return result, names
		}

		return make(map[string]interface{}), []string{}
	}
}

func (m *Manager) createPreExistenceChecker(config ResourceConfig, operationType string) func(context.Context, []string, map[string]interface{}) ([]string, map[string]interface{}, error) {
	if operationType != "PUT" {
		return nil // Only PUT operations need pre-existence checking
	}

	return func(ctx context.Context, resourceNames []string, originalOperations map[string]interface{}) ([]string, map[string]interface{}, error) {
		checker := ResourceExistenceCheck{
			ResourceType:  config.ResourceType,
			OperationType: "PUT",
			FetchResources: func(ctx context.Context) (map[string]interface{}, error) {
				apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
				defer cancel()

				switch config.ResourceType {
				case "gateway":
					resp, err := m.client.GatewaysAPI.GatewaysGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						Gateway map[string]interface{} `json:"gateway"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.Gateway, nil

				case "lag":
					resp, err := m.client.LAGsAPI.LagsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						Lag map[string]interface{} `json:"lag"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.Lag, nil

				case "tenant":
					resp, err := m.client.TenantsAPI.TenantsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						Tenant map[string]interface{} `json:"tenant"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.Tenant, nil

				case "service":
					resp, err := m.client.ServicesAPI.ServicesGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						Service map[string]interface{} `json:"service"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.Service, nil

				case "gateway_profile":
					resp, err := m.client.GatewayProfilesAPI.GatewayprofilesGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						GatewayProfile map[string]interface{} `json:"gateway_profile"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.GatewayProfile, nil

				case "eth_port_profile":
					resp, err := m.client.EthPortProfilesAPI.EthportprofilesGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						EthPortProfile map[string]interface{} `json:"eth_port_profile_"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.EthPortProfile, nil

				case "eth_port_settings":
					resp, err := m.client.EthPortSettingsAPI.EthportsettingsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						EthPortSettings map[string]interface{} `json:"eth_port_settings"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.EthPortSettings, nil

				case "device_settings":
					resp, err := m.client.DeviceSettingsAPI.DevicesettingsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						EthDeviceProfiles map[string]interface{} `json:"eth_device_profiles"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.EthDeviceProfiles, nil

				case "bundle":
					resp, err := m.client.BundlesAPI.BundlesGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						Bundle map[string]interface{} `json:"endpoint_bundle"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.Bundle, nil

				case "acl":
					// ACLs require IP version checking - need to check both IPv4 and IPv6

					// Try IPv4
					resp4, err4 := m.client.ACLsAPI.AclsGet(apiCtx).IpVersion("4").Execute()
					existingResources := make(map[string]interface{})

					if err4 == nil {
						defer resp4.Body.Close()
						var result4 struct {
							IpFilter map[string]interface{} `json:"ip_filter"`
						}
						if err := json.NewDecoder(resp4.Body).Decode(&result4); err == nil {
							for name, props := range result4.IpFilter {
								existingResources[name] = props
							}
						}
					}

					// Try IPv6
					resp6, err6 := m.client.ACLsAPI.AclsGet(apiCtx).IpVersion("6").Execute()
					if err6 == nil {
						defer resp6.Body.Close()
						var result6 struct {
							IpFilter map[string]interface{} `json:"ip_filter"`
						}
						if err := json.NewDecoder(resp6.Body).Decode(&result6); err == nil {
							for name, props := range result6.IpFilter {
								existingResources[name] = props
							}
						}
					}

					if err4 != nil && err6 != nil {
						return nil, err4
					}

					return existingResources, nil

				case "ipv4_list":
					resp, err := m.client.IPv4ListFiltersAPI.Ipv4listsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						Ipv4ListFilter map[string]interface{} `json:"ipv4_list_filter"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.Ipv4ListFilter, nil

				case "ipv4_prefix_list":
					resp, err := m.client.IPv4PrefixListsAPI.Ipv4prefixlistsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						Ipv4PrefixList map[string]interface{} `json:"ipv4_prefix_list"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.Ipv4PrefixList, nil

				case "ipv6_list":
					resp, err := m.client.IPv6ListFiltersAPI.Ipv6listsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						Ipv6ListFilter map[string]interface{} `json:"ipv6_list_filter"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.Ipv6ListFilter, nil

				case "ipv6_prefix_list":
					resp, err := m.client.IPv6PrefixListsAPI.Ipv6prefixlistsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						Ipv6PrefixList map[string]interface{} `json:"ipv6_prefix_list"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.Ipv6PrefixList, nil

				case "authenticated_eth_port":
					resp, err := m.client.AuthenticatedEthPortsAPI.AuthenticatedethportsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						AuthenticatedEthPort map[string]interface{} `json:"authenticated_eth_port"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.AuthenticatedEthPort, nil

				case "badge":
					resp, err := m.client.BadgesAPI.BadgesGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						Badge map[string]interface{} `json:"badge"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.Badge, nil

				case "device_voice_settings":
					resp, err := m.client.DeviceVoiceSettingsAPI.DevicevoicesettingsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						DeviceVoiceSettings map[string]interface{} `json:"device_voice_settings"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.DeviceVoiceSettings, nil

				case "as_path_access_list":
					resp, err := m.client.ASPathAccessListsAPI.AspathaccesslistsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						AsPathAccessList map[string]interface{} `json:"as_path_access_list"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.AsPathAccessList, nil

				case "community_list":
					resp, err := m.client.CommunityListsAPI.CommunitylistsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						CommunityList map[string]interface{} `json:"community_list"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.CommunityList, nil

				case "extended_community_list":
					resp, err := m.client.ExtendedCommunityListsAPI.ExtendedcommunitylistsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						ExtendedCommunityList map[string]interface{} `json:"extended_community_list"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.ExtendedCommunityList, nil

				case "route_map_clause":
					resp, err := m.client.RouteMapClausesAPI.RoutemapclausesGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						RouteMapClause map[string]interface{} `json:"route_map_clause"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.RouteMapClause, nil

				case "route_map":
					resp, err := m.client.RouteMapsAPI.RoutemapsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						RouteMap map[string]interface{} `json:"route_map"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.RouteMap, nil

				case "packet_broker":
					resp, err := m.client.PacketBrokerAPI.PacketbrokerGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						PacketBroker map[string]interface{} `json:"pb_egress_profile"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.PacketBroker, nil

				case "packet_queue":
					resp, err := m.client.PacketQueuesAPI.PacketqueuesGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						PacketQueue map[string]interface{} `json:"packet_queue"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.PacketQueue, nil

				case "service_port_profile":
					resp, err := m.client.ServicePortProfilesAPI.ServiceportprofilesGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						ServicePortProfile map[string]interface{} `json:"service_port_profile"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.ServicePortProfile, nil

				case "switchpoint":
					resp, err := m.client.SwitchpointsAPI.SwitchpointsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						Switchpoint map[string]interface{} `json:"switchpoint"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.Switchpoint, nil

				case "voice_port_profile":
					resp, err := m.client.VoicePortProfilesAPI.VoiceportprofilesGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						VoicePortProfile map[string]interface{} `json:"voice_port_profiles"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.VoicePortProfile, nil

				case "device_controller":
					resp, err := m.client.DeviceControllersAPI.DevicecontrollersGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						DeviceController map[string]interface{} `json:"device_controller"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.DeviceController, nil

				case "pod":
					resp, err := m.client.PodsAPI.PodsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						Pod map[string]interface{} `json:"pod"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.Pod, nil

				case "port_acl":
					resp, err := m.client.PortACLsAPI.PortaclsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						PortAcl map[string]interface{} `json:"port_acl"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.PortAcl, nil

				case "sflow_collector":
					resp, err := m.client.SFlowCollectorsAPI.SflowcollectorsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						SflowCollector map[string]interface{} `json:"sflow_collector"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.SflowCollector, nil

				case "diagnostics_profile":
					resp, err := m.client.DiagnosticsProfilesAPI.DiagnosticsprofilesGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						DiagnosticsProfile map[string]interface{} `json:"diagnostics_profile"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.DiagnosticsProfile, nil

				case "diagnostics_port_profile":
					resp, err := m.client.DiagnosticsPortProfilesAPI.DiagnosticsportprofilesGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						DiagnosticsPortProfile map[string]interface{} `json:"diagnostics_port_profile"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.DiagnosticsPortProfile, nil

				case "pb_routing":
					resp, err := m.client.PBRoutingAPI.PolicybasedroutingGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						PbRouting map[string]interface{} `json:"pb_routing"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.PbRouting, nil

				case "pb_routing_acl":
					resp, err := m.client.PBRoutingACLAPI.PolicybasedroutingaclGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						PbRoutingAcl map[string]interface{} `json:"pb_routing_acl"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.PbRoutingAcl, nil

				case "spine_plane":
					resp, err := m.client.SpinePlanesAPI.SpineplanesGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						SpinePlane map[string]interface{} `json:"spine_plane"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.SpinePlane, nil

				case "grouping_rule":
					resp, err := m.client.GroupingRulesAPI.GroupingrulesGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						GroupingRule map[string]interface{} `json:"grouping_rule"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.GroupingRule, nil

				case "threshold_group":
					resp, err := m.client.ThresholdGroupsAPI.ThresholdgroupsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						ThresholdGroup map[string]interface{} `json:"threshold_group"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.ThresholdGroup, nil

				case "threshold":
					resp, err := m.client.ThresholdsAPI.ThresholdsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						Threshold map[string]interface{} `json:"threshold"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.Threshold, nil

				default:
					// For unknown resource types, assume no existing resources to avoid errors
					return make(map[string]interface{}), nil
				}
			},
		}

		filteredNames, err := m.FilterPreExistingResources(ctx, resourceNames, checker)
		if err != nil {
			return resourceNames, nil, err
		}

		filteredOperations := make(map[string]interface{})
		for _, name := range filteredNames {
			if val, ok := originalOperations[name]; ok {
				filteredOperations[name] = val
			}
		}

		return filteredNames, filteredOperations, nil
	}
}

func (m *Manager) createRequestPreparer(config ResourceConfig, operationType string) func(map[string]interface{}) interface{} {
	return func(filteredData map[string]interface{}) interface{} {
		// Handle DELETE operations
		if operationType == "DELETE" {
			names := make([]string, 0, len(filteredData))
			for name := range filteredData {
				names = append(names, name)
			}
			return names
		}

		// Handle PUT and PATCH operations
		switch config.ResourceType {
		case "gateway":
			putRequest := openapi.NewGatewaysPutRequest()
			gatewayMap := make(map[string]openapi.GatewaysPutRequestGatewayValue)
			for name, props := range filteredData {
				gatewayMap[name] = props.(openapi.GatewaysPutRequestGatewayValue)
			}
			putRequest.SetGateway(gatewayMap)
			return putRequest
		case "lag":
			putRequest := openapi.NewLagsPutRequest()
			lagMap := make(map[string]openapi.LagsPutRequestLagValue)
			for name, props := range filteredData {
				lagMap[name] = props.(openapi.LagsPutRequestLagValue)
			}
			putRequest.SetLag(lagMap)
			return putRequest
		case "tenant":
			putRequest := openapi.NewTenantsPutRequest()
			tenantMap := make(map[string]openapi.TenantsPutRequestTenantValue)
			for name, props := range filteredData {
				tenantMap[name] = props.(openapi.TenantsPutRequestTenantValue)
			}
			putRequest.SetTenant(tenantMap)
			return putRequest
		case "service":
			putRequest := openapi.NewServicesPutRequest()
			serviceMap := make(map[string]openapi.ServicesPutRequestServiceValue)
			for name, props := range filteredData {
				serviceMap[name] = props.(openapi.ServicesPutRequestServiceValue)
			}
			putRequest.SetService(serviceMap)
			return putRequest
		case "gateway_profile":
			putRequest := openapi.NewGatewayprofilesPutRequest()
			profileMap := make(map[string]openapi.GatewayprofilesPutRequestGatewayProfileValue)
			for name, props := range filteredData {
				profileMap[name] = props.(openapi.GatewayprofilesPutRequestGatewayProfileValue)
			}
			putRequest.SetGatewayProfile(profileMap)
			return putRequest
		case "eth_port_profile":
			putRequest := openapi.NewEthportprofilesPutRequest()
			profileMap := make(map[string]openapi.EthportprofilesPutRequestEthPortProfileValue)
			for name, props := range filteredData {
				profileMap[name] = props.(openapi.EthportprofilesPutRequestEthPortProfileValue)
			}
			putRequest.SetEthPortProfile(profileMap)
			return putRequest
		case "eth_port_settings":
			putRequest := openapi.NewEthportsettingsPutRequest()
			settingsMap := make(map[string]openapi.EthportsettingsPutRequestEthPortSettingsValue)
			for name, props := range filteredData {
				settingsMap[name] = props.(openapi.EthportsettingsPutRequestEthPortSettingsValue)
			}
			putRequest.SetEthPortSettings(settingsMap)
			return putRequest
		case "device_settings":
			putRequest := openapi.NewDevicesettingsPutRequest()
			settingsMap := make(map[string]openapi.DevicesettingsPutRequestEthDeviceProfilesValue)
			for name, props := range filteredData {
				settingsMap[name] = props.(openapi.DevicesettingsPutRequestEthDeviceProfilesValue)
			}
			putRequest.SetEthDeviceProfiles(settingsMap)
			return putRequest
		case "bundle":
			putRequest := openapi.NewBundlesPutRequest()
			bundleMap := make(map[string]openapi.BundlesPutRequestEndpointBundleValue)
			for name, props := range filteredData {
				bundleMap[name] = props.(openapi.BundlesPutRequestEndpointBundleValue)
			}
			putRequest.SetEndpointBundle(bundleMap)
			return putRequest
		case "acl":
			putRequest := openapi.NewAclsPutRequest()
			aclMap := make(map[string]openapi.AclsPutRequestIpFilterValue)
			for name, props := range filteredData {
				aclMap[name] = props.(openapi.AclsPutRequestIpFilterValue)
			}
			putRequest.SetIpFilter(aclMap)
			return putRequest
		case "ipv4_list":
			putRequest := openapi.NewIpv4listsPutRequest()
			ipv4Map := make(map[string]openapi.Ipv4listsPutRequestIpv4ListFilterValue)
			for name, props := range filteredData {
				ipv4Map[name] = props.(openapi.Ipv4listsPutRequestIpv4ListFilterValue)
			}
			putRequest.SetIpv4ListFilter(ipv4Map)
			return putRequest
		case "ipv4_prefix_list":
			putRequest := openapi.NewIpv4prefixlistsPutRequest()
			ipv4PrefixMap := make(map[string]openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValue)
			for name, props := range filteredData {
				ipv4PrefixMap[name] = props.(openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValue)
			}
			putRequest.SetIpv4PrefixList(ipv4PrefixMap)
			return putRequest
		case "ipv6_list":
			putRequest := openapi.NewIpv6listsPutRequest()
			ipv6Map := make(map[string]openapi.Ipv6listsPutRequestIpv6ListFilterValue)
			for name, props := range filteredData {
				ipv6Map[name] = props.(openapi.Ipv6listsPutRequestIpv6ListFilterValue)
			}
			putRequest.SetIpv6ListFilter(ipv6Map)
			return putRequest
		case "ipv6_prefix_list":
			putRequest := openapi.NewIpv6prefixlistsPutRequest()
			ipv6PrefixMap := make(map[string]openapi.Ipv6prefixlistsPutRequestIpv6PrefixListValue)
			for name, props := range filteredData {
				ipv6PrefixMap[name] = props.(openapi.Ipv6prefixlistsPutRequestIpv6PrefixListValue)
			}
			putRequest.SetIpv6PrefixList(ipv6PrefixMap)
			return putRequest
		case "authenticated_eth_port":
			putRequest := openapi.NewAuthenticatedethportsPutRequest()
			portMap := make(map[string]openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValue)
			for name, props := range filteredData {
				portMap[name] = props.(openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValue)
			}
			putRequest.SetAuthenticatedEthPort(portMap)
			return putRequest
		case "badge":
			putRequest := openapi.NewBadgesPutRequest()
			badgeMap := make(map[string]openapi.BadgesPutRequestBadgeValue)
			for name, props := range filteredData {
				badgeMap[name] = props.(openapi.BadgesPutRequestBadgeValue)
			}
			putRequest.SetBadge(badgeMap)
			return putRequest
		case "device_voice_settings":
			putRequest := openapi.NewDevicevoicesettingsPutRequest()
			settingsMap := make(map[string]openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValue)
			for name, props := range filteredData {
				settingsMap[name] = props.(openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValue)
			}
			putRequest.SetDeviceVoiceSettings(settingsMap)
			return putRequest
		case "as_path_access_list":
			putRequest := openapi.NewAspathaccesslistsPutRequest()
			asPathMap := make(map[string]openapi.AspathaccesslistsPutRequestAsPathAccessListValue)
			for name, props := range filteredData {
				asPathMap[name] = props.(openapi.AspathaccesslistsPutRequestAsPathAccessListValue)
			}
			putRequest.SetAsPathAccessList(asPathMap)
			return putRequest
		case "community_list":
			putRequest := openapi.NewCommunitylistsPutRequest()
			communityMap := make(map[string]openapi.CommunitylistsPutRequestCommunityListValue)
			for name, props := range filteredData {
				communityMap[name] = props.(openapi.CommunitylistsPutRequestCommunityListValue)
			}
			putRequest.SetCommunityList(communityMap)
			return putRequest
		case "extended_community_list":
			putRequest := openapi.NewExtendedcommunitylistsPutRequest()
			extCommunityMap := make(map[string]openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValue)
			for name, props := range filteredData {
				extCommunityMap[name] = props.(openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValue)
			}
			putRequest.SetExtendedCommunityList(extCommunityMap)
			return putRequest
		case "route_map_clause":
			putRequest := openapi.NewRoutemapclausesPutRequest()
			routeMapClauseMap := make(map[string]openapi.RoutemapclausesPutRequestRouteMapClauseValue)
			for name, props := range filteredData {
				routeMapClauseMap[name] = props.(openapi.RoutemapclausesPutRequestRouteMapClauseValue)
			}
			putRequest.SetRouteMapClause(routeMapClauseMap)
			return putRequest
		case "route_map":
			putRequest := openapi.NewRoutemapsPutRequest()
			routeMapMap := make(map[string]openapi.RoutemapsPutRequestRouteMapValue)
			for name, props := range filteredData {
				routeMapMap[name] = props.(openapi.RoutemapsPutRequestRouteMapValue)
			}
			putRequest.SetRouteMap(routeMapMap)
			return putRequest
		case "sfp_breakout":
			// SFP Breakouts only support PATCH operations
			patchRequest := openapi.NewSfpbreakoutsPatchRequest()
			sfpMap := make(map[string]openapi.SfpbreakoutsPatchRequestSfpBreakoutsValue)
			for name, props := range filteredData {
				sfpMap[name] = props.(openapi.SfpbreakoutsPatchRequestSfpBreakoutsValue)
			}
			patchRequest.SetSfpBreakouts(sfpMap)
			return patchRequest
		case "site":
			// Sites only support PATCH operations
			patchRequest := openapi.NewSitesPatchRequest()
			siteMap := make(map[string]openapi.SitesPatchRequestSiteValue)
			for name, props := range filteredData {
				siteMap[name] = props.(openapi.SitesPatchRequestSiteValue)
			}
			patchRequest.SetSite(siteMap)
			return patchRequest
		case "packet_broker":
			putRequest := openapi.NewPacketbrokerPutRequest()
			brokerMap := make(map[string]openapi.PacketbrokerPutRequestPortAclValue)
			for name, props := range filteredData {
				brokerMap[name] = props.(openapi.PacketbrokerPutRequestPortAclValue)
			}
			putRequest.SetPortAcl(brokerMap)
			return putRequest
		case "packet_queue":
			putRequest := openapi.NewPacketqueuesPutRequest()
			queueMap := make(map[string]openapi.PacketqueuesPutRequestPacketQueueValue)
			for name, props := range filteredData {
				queueMap[name] = props.(openapi.PacketqueuesPutRequestPacketQueueValue)
			}
			putRequest.SetPacketQueue(queueMap)
			return putRequest
		case "service_port_profile":
			putRequest := openapi.NewServiceportprofilesPutRequest()
			profileMap := make(map[string]openapi.ServiceportprofilesPutRequestServicePortProfileValue)
			for name, props := range filteredData {
				profileMap[name] = props.(openapi.ServiceportprofilesPutRequestServicePortProfileValue)
			}
			putRequest.SetServicePortProfile(profileMap)
			return putRequest
		case "switchpoint":
			putRequest := openapi.NewSwitchpointsPutRequest()
			switchpointMap := make(map[string]openapi.SwitchpointsPutRequestSwitchpointValue)
			for name, props := range filteredData {
				switchpointMap[name] = props.(openapi.SwitchpointsPutRequestSwitchpointValue)
			}
			putRequest.SetSwitchpoint(switchpointMap)
			return putRequest
		case "voice_port_profile":
			putRequest := openapi.NewVoiceportprofilesPutRequest()
			profileMap := make(map[string]openapi.VoiceportprofilesPutRequestVoicePortProfilesValue)
			for name, props := range filteredData {
				profileMap[name] = props.(openapi.VoiceportprofilesPutRequestVoicePortProfilesValue)
			}
			putRequest.SetVoicePortProfiles(profileMap)
			return putRequest
		case "device_controller":
			putRequest := openapi.NewDevicecontrollersPutRequest()
			deviceMap := make(map[string]openapi.DevicecontrollersPutRequestDeviceControllerValue)
			for name, props := range filteredData {
				deviceMap[name] = props.(openapi.DevicecontrollersPutRequestDeviceControllerValue)
			}
			putRequest.SetDeviceController(deviceMap)
			return putRequest
		case "pod":
			putRequest := openapi.NewPodsPutRequest()
			podMap := make(map[string]openapi.PodsPutRequestPodValue)
			for name, props := range filteredData {
				podMap[name] = props.(openapi.PodsPutRequestPodValue)
			}
			putRequest.SetPod(podMap)
			return putRequest
		case "port_acl":
			putRequest := openapi.NewPortaclsPutRequest()
			portAclMap := make(map[string]openapi.PortaclsPutRequestPortAclValue)
			for name, props := range filteredData {
				portAclMap[name] = props.(openapi.PortaclsPutRequestPortAclValue)
			}
			putRequest.SetPortAcl(portAclMap)
			return putRequest
		case "sflow_collector":
			putRequest := openapi.NewSflowcollectorsPutRequest()
			sflowMap := make(map[string]openapi.SflowcollectorsPutRequestSflowCollectorValue)
			for name, props := range filteredData {
				sflowMap[name] = props.(openapi.SflowcollectorsPutRequestSflowCollectorValue)
			}
			putRequest.SetSflowCollector(sflowMap)
			return putRequest
		case "diagnostics_profile":
			putRequest := openapi.NewDiagnosticsprofilesPutRequest()
			diagnosticsMap := make(map[string]openapi.DiagnosticsprofilesPutRequestDiagnosticsProfileValue)
			for name, props := range filteredData {
				diagnosticsMap[name] = props.(openapi.DiagnosticsprofilesPutRequestDiagnosticsProfileValue)
			}
			putRequest.SetDiagnosticsProfile(diagnosticsMap)
			return putRequest
		case "diagnostics_port_profile":
			putRequest := openapi.NewDiagnosticsportprofilesPutRequest()
			diagnosticsPortMap := make(map[string]openapi.DiagnosticsportprofilesPutRequestDiagnosticsPortProfileValue)
			for name, props := range filteredData {
				diagnosticsPortMap[name] = props.(openapi.DiagnosticsportprofilesPutRequestDiagnosticsPortProfileValue)
			}
			putRequest.SetDiagnosticsPortProfile(diagnosticsPortMap)
			return putRequest
		case "pb_routing":
			putRequest := openapi.NewPolicybasedroutingPutRequest()
			pbRoutingMap := make(map[string]openapi.PolicybasedroutingPutRequestPbRoutingValue)
			for name, props := range filteredData {
				pbRoutingMap[name] = props.(openapi.PolicybasedroutingPutRequestPbRoutingValue)
			}
			putRequest.SetPbRouting(pbRoutingMap)
			return putRequest
		case "pb_routing_acl":
			putRequest := openapi.NewPolicybasedroutingaclPutRequest()
			pbRoutingAclMap := make(map[string]openapi.PolicybasedroutingaclPutRequestPbRoutingAclValue)
			for name, props := range filteredData {
				pbRoutingAclMap[name] = props.(openapi.PolicybasedroutingaclPutRequestPbRoutingAclValue)
			}
			putRequest.SetPbRoutingAcl(pbRoutingAclMap)
			return putRequest
		case "spine_plane":
			putRequest := openapi.NewSpineplanesPutRequest()
			spinePlaneMap := make(map[string]openapi.SpineplanesPutRequestSpinePlaneValue)
			for name, props := range filteredData {
				spinePlaneMap[name] = props.(openapi.SpineplanesPutRequestSpinePlaneValue)
			}
			putRequest.SetSpinePlane(spinePlaneMap)
			return putRequest
		case "grouping_rule":
			putRequest := openapi.NewGroupingrulesPutRequest()
			groupingRuleMap := make(map[string]openapi.GroupingrulesPutRequestGroupingRulesValue)
			for name, props := range filteredData {
				groupingRuleMap[name] = props.(openapi.GroupingrulesPutRequestGroupingRulesValue)
			}
			putRequest.SetGroupingRules(groupingRuleMap)
			return putRequest
		case "threshold_group":
			putRequest := openapi.NewThresholdgroupsPutRequest()
			thresholdGroupMap := make(map[string]openapi.ThresholdgroupsPutRequestThresholdGroupValue)
			for name, props := range filteredData {
				thresholdGroupMap[name] = props.(openapi.ThresholdgroupsPutRequestThresholdGroupValue)
			}
			putRequest.SetThresholdGroup(thresholdGroupMap)
			return putRequest
		case "threshold":
			putRequest := openapi.NewThresholdsPutRequest()
			thresholdMap := make(map[string]openapi.ThresholdsPutRequestThresholdValue)
			for name, props := range filteredData {
				thresholdMap[name] = props.(openapi.ThresholdsPutRequestThresholdValue)
			}
			putRequest.SetThreshold(thresholdMap)
			return putRequest
		}
		return nil
	}
}

func (m *Manager) createRequestExecutor(config ResourceConfig, operationType string) func(context.Context, interface{}) (*http.Response, error) {
	return func(ctx context.Context, request interface{}) (*http.Response, error) {
		apiClient := config.APIClientGetter(m.client)

		switch operationType {
		case "PUT":
			return apiClient.Put(ctx, request)
		case "PATCH":
			return apiClient.Patch(ctx, request)
		case "DELETE":
			return apiClient.Delete(ctx, request.([]string))
		}
		return nil, fmt.Errorf("unknown operation type: %s", operationType)
	}
}

func (m *Manager) createResponseProcessor(config ResourceConfig, operationType string) func(context.Context, *http.Response) error {
	if !config.HasAutoGen {
		return nil // No post-processing needed for resources without auto-generated fields
	}

	return func(ctx context.Context, resp *http.Response) error {
		delayTime := 2 * time.Second
		tflog.Debug(ctx, fmt.Sprintf("Waiting %v for auto-generated values to be assigned before fetching %s", delayTime, config.ResourceType))
		time.Sleep(delayTime)

		fetchCtx, fetchCancel := context.WithTimeout(context.Background(), OperationTimeout)
		defer fetchCancel()

		tflog.Debug(ctx, fmt.Sprintf("Fetching %s after successful PUT operation to retrieve auto-generated values", config.ResourceType))

		res, exists := m.resources[config.ResourceType]
		if !exists {
			return fmt.Errorf("resource type %s not found in unified structure", config.ResourceType)
		}

		switch config.ResourceType {
		case "tenant":
			tenantsReq := m.client.TenantsAPI.TenantsGet(fetchCtx)
			tenantsResp, fetchErr := tenantsReq.Execute()

			if fetchErr != nil {
				tflog.Error(ctx, "Failed to fetch tenants after PUT for auto-generated fields", map[string]interface{}{
					"error": fetchErr.Error(),
				})
				return fetchErr
			}

			defer tenantsResp.Body.Close()

			var tenantsData struct {
				Tenant map[string]map[string]interface{} `json:"tenant"`
			}

			if respErr := json.NewDecoder(tenantsResp.Body).Decode(&tenantsData); respErr != nil {
				tflog.Error(ctx, "Failed to decode tenants response for auto-generated fields", map[string]interface{}{
					"error": respErr.Error(),
				})
				return respErr
			}

			res.ResponsesMutex.Lock()
			for tenantName, tenantData := range tenantsData.Tenant {
				res.Responses[tenantName] = tenantData
				if name, ok := tenantData["name"].(string); ok && name != tenantName {
					res.Responses[name] = tenantData
				}
			}
			res.ResponsesMutex.Unlock()

			tflog.Debug(ctx, "Successfully stored tenant data for auto-generated fields", map[string]interface{}{
				"tenant_count": len(tenantsData.Tenant),
			})

		case "service":
			servicesReq := m.client.ServicesAPI.ServicesGet(fetchCtx)
			servicesResp, fetchErr := servicesReq.Execute()

			if fetchErr != nil {
				tflog.Error(ctx, "Failed to fetch services after PUT for auto-generated fields", map[string]interface{}{
					"error": fetchErr.Error(),
				})
				return fetchErr
			}

			defer servicesResp.Body.Close()

			var servicesData struct {
				Service map[string]map[string]interface{} `json:"service"`
			}

			if respErr := json.NewDecoder(servicesResp.Body).Decode(&servicesData); respErr != nil {
				tflog.Error(ctx, "Failed to decode services response for auto-generated fields", map[string]interface{}{
					"error": respErr.Error(),
				})
				return respErr
			}

			res.ResponsesMutex.Lock()
			for serviceName, serviceData := range servicesData.Service {
				res.Responses[serviceName] = serviceData
				if name, ok := serviceData["name"].(string); ok && name != serviceName {
					res.Responses[name] = serviceData
				}
			}
			res.ResponsesMutex.Unlock()

			tflog.Debug(ctx, "Successfully stored service data for auto-generated fields", map[string]interface{}{
				"service_count": len(servicesData.Service),
			})

		case "switchpoint":
			switchpointsReq := m.client.SwitchpointsAPI.SwitchpointsGet(fetchCtx)
			switchpointsResp, fetchErr := switchpointsReq.Execute()

			if fetchErr != nil {
				tflog.Error(ctx, "Failed to fetch switchpoints after PUT for auto-generated fields", map[string]interface{}{
					"error": fetchErr.Error(),
				})
				return fetchErr
			}

			defer switchpointsResp.Body.Close()

			var switchpointsData struct {
				Switchpoint map[string]map[string]interface{} `json:"switchpoint"`
			}

			if respErr := json.NewDecoder(switchpointsResp.Body).Decode(&switchpointsData); respErr != nil {
				tflog.Error(ctx, "Failed to decode switchpoints response for auto-generated fields", map[string]interface{}{
					"error": respErr.Error(),
				})
				return respErr
			}

			res.ResponsesMutex.Lock()
			for switchpointName, switchpointData := range switchpointsData.Switchpoint {
				res.Responses[switchpointName] = switchpointData
				if name, ok := switchpointData["name"].(string); ok && name != switchpointName {
					res.Responses[name] = switchpointData
				}
			}
			res.ResponsesMutex.Unlock()

			tflog.Debug(ctx, "Successfully stored switchpoint data for auto-generated fields", map[string]interface{}{
				"switchpoint_count": len(switchpointsData.Switchpoint),
			})

		case "site":
			sitesReq := m.client.SitesAPI.SitesGet(fetchCtx)
			sitesResp, fetchErr := sitesReq.Execute()

			if fetchErr != nil {
				tflog.Error(ctx, "Failed to fetch sites after PATCH for auto-generated fields", map[string]interface{}{
					"error": fetchErr.Error(),
				})
				return fetchErr
			}

			defer sitesResp.Body.Close()

			var sitesData struct {
				Site map[string]map[string]interface{} `json:"site"`
			}

			if respErr := json.NewDecoder(sitesResp.Body).Decode(&sitesData); respErr != nil {
				tflog.Error(ctx, "Failed to decode sites response for auto-generated fields", map[string]interface{}{
					"error": respErr.Error(),
				})
				return respErr
			}

			res.ResponsesMutex.Lock()
			for siteName, siteData := range sitesData.Site {
				res.Responses[siteName] = siteData
				if name, ok := siteData["name"].(string); ok && name != siteName {
					res.Responses[name] = siteData
				}
			}
			res.ResponsesMutex.Unlock()

			tflog.Debug(ctx, "Successfully stored site data for auto-generated fields", map[string]interface{}{
				"site_count": len(sitesData.Site),
			})

		default:
			tflog.Warn(ctx, fmt.Sprintf("Unknown resource type with auto-generated fields: %s", config.ResourceType))
		}

		return nil
	}
}

func (m *Manager) createRecentOpsUpdater(resourceType string) func() {
	return func() {
		now := time.Now()

		if res, exists := m.resources[resourceType]; exists {
			res.RecentOps = true
			res.RecentOpTime = now
		}
	}
}

// FilterPreExistingResources filters out resources that already exist in the system.
func (m *Manager) FilterPreExistingResources(
	ctx context.Context,
	resourceNames []string,
	checker ResourceExistenceCheck,
) ([]string, error) {
	existingResources, err := checker.FetchResources(ctx)
	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf("Failed to fetch existing %s for pre-flight check: %v",
			checker.ResourceType, err))
		return resourceNames, nil
	}

	var notExistingResources []string
	alreadyExistingResources := make(map[string]bool)

	for _, name := range resourceNames {
		if _, exists := existingResources[name]; exists {
			// Resource already exists
			alreadyExistingResources[name] = true
			tflog.Info(ctx, fmt.Sprintf("Skipping creation of %s '%s' as it already exists",
				checker.ResourceType, name))
		} else {
			// Resource doesn't exist - add to filtered list
			notExistingResources = append(notExistingResources, name)
		}
	}

	// Update operation tracking for already existing resources
	if len(alreadyExistingResources) > 0 {
		m.operationMutex.Lock()
		defer m.operationMutex.Unlock()

		for opID, op := range m.pendingOperations {
			if op.ResourceType == checker.ResourceType &&
				op.OperationType == checker.OperationType &&
				alreadyExistingResources[op.ResourceName] {
				// Mark operation as successful
				updatedOp := *op
				updatedOp.Status = OperationSucceeded
				m.pendingOperations[opID] = &updatedOp
				m.operationResults[opID] = true

				m.safeCloseChannel(opID, true)

				tflog.Debug(ctx, fmt.Sprintf("Marked operation %s as successful since resource already exists", opID))
			}
		}
	}

	return notExistingResources, nil
}

// getOperationCount returns the number of pending operations for a resource type
func (m *Manager) getOperationCount(resourceType, operationType string) int {
	res, exists := m.resources[resourceType]
	if !exists {
		return 0
	}

	switch operationType {
	case "PUT":
		if res.Put == nil {
			return 0
		}
		return len(res.Put)
	case "PATCH":
		if res.Patch == nil {
			return 0
		}
		return len(res.Patch)
	case "DELETE":
		return len(res.Delete)
	}
	return 0
}

// checks if we need to apply the circular reference fix for route_map_clause and tenant resources
func (m *Manager) detectCircularReferenceScenario(ctx context.Context) (bool, map[string]interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Check if both route_map_clause and tenant have PUT operations
	routeMapClauseOps := m.resources["route_map_clause"]
	tenantOps := m.resources["tenant"]

	if routeMapClauseOps == nil || tenantOps == nil {
		return false, nil
	}

	if len(routeMapClauseOps.Put) == 0 || len(tenantOps.Put) == 0 {
		return false, nil
	}

	tenantsBeingCreated := make(map[string]bool)
	for tenantName := range tenantOps.Put {
		tenantsBeingCreated[tenantName] = true
	}

	// When tenants are being created, we need to apply the fix to all route_map_clauses
	// with match_vrf (not just those referencing tenants being created) because:
	// 1. Clauses referencing tenants being created need the fix (circular dependency)
	// 2. Clauses referencing other tenants need the fix too (execution order - they come before tenant)
	// This ensures all match_vrf references are set after all tenants exist
	affectedClauses := make(map[string]interface{})
	circularRefClauses := make([]string, 0)
	otherRefClauses := make([]string, 0)

	// Check PUT operations
	for name, clauseData := range routeMapClauseOps.Put {
		if matchVrfValue := m.getMatchVrfValue(clauseData); matchVrfValue != "" {
			affectedClauses[name] = clauseData
			if tenantsBeingCreated[matchVrfValue] {
				circularRefClauses = append(circularRefClauses, name)
				tflog.Debug(ctx, fmt.Sprintf("Circular ref: route_map_clause '%s' references tenant '%s' being created", name, matchVrfValue))
			} else {
				otherRefClauses = append(otherRefClauses, name)
				tflog.Debug(ctx, fmt.Sprintf("Other ref: route_map_clause '%s' references tenant '%s' (execution order fix)", name, matchVrfValue))
			}
		}
	}

	// Also check PATCH operations for match_vrf changes
	if len(routeMapClauseOps.Patch) > 0 {
		for name, patchData := range routeMapClauseOps.Patch {
			if matchVrfValue := m.getMatchVrfValue(patchData); matchVrfValue != "" {
				affectedClauses[name] = patchData
				if tenantsBeingCreated[matchVrfValue] {
					circularRefClauses = append(circularRefClauses, name)
					tflog.Debug(ctx, fmt.Sprintf("Circular ref in PATCH: route_map_clause '%s' references tenant '%s' being created", name, matchVrfValue))
				} else {
					otherRefClauses = append(otherRefClauses, name)
					tflog.Debug(ctx, fmt.Sprintf("Other ref in PATCH: route_map_clause '%s' references tenant '%s' (execution order fix)", name, matchVrfValue))
				}
			}
		}
	}

	needsFix := len(affectedClauses) > 0

	if needsFix {
		tflog.Info(ctx, fmt.Sprintf("Applying match_vrf fix: %d circular refs, %d execution order refs", len(circularRefClauses), len(otherRefClauses)))
		tflog.Debug(ctx, "Match VRF fix details", map[string]interface{}{
			"route_map_clause_put_count":   len(routeMapClauseOps.Put),
			"route_map_clause_patch_count": len(routeMapClauseOps.Patch),
			"tenant_put_count":             len(tenantOps.Put),
			"tenants_being_created":        getMapKeys(tenantsBeingCreated),
			"circular_ref_clauses":         circularRefClauses,
			"other_ref_clauses":            otherRefClauses,
			"total_affected_clauses":       len(affectedClauses),
		})
	}

	return needsFix, affectedClauses
}

func (m *Manager) getMatchVrfValue(data interface{}) string {
	clauseMap, ok := data.(map[string]interface{})
	if !ok {
		jsonData, err := json.Marshal(data)
		if err != nil {
			return ""
		}
		var tempMap map[string]interface{}
		if err := json.Unmarshal(jsonData, &tempMap); err != nil {
			return ""
		}
		clauseMap = tempMap
	}

	if matchVrf, exists := clauseMap["match_vrf"]; exists {
		if matchVrfStr, ok := matchVrf.(string); ok {
			return matchVrfStr
		} else if matchVrfPtr, ok := matchVrf.(*string); ok && matchVrfPtr != nil {
			return *matchVrfPtr
		}
	}

	return ""
}

func getMapKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// creates a copy of route_map_clause data with match_vrf set to empty
func (m *Manager) createRouteMapClauseWithEmptyMatchVrf(originalData interface{}) interface{} {
	jsonData, err := json.Marshal(originalData)
	if err != nil {
		return originalData
	}

	var result openapi.RoutemapclausesPutRequestRouteMapClauseValue
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return originalData
	}

	// Set match_vrf to empty string
	result.MatchVrf = openapi.PtrString("")

	return result
}

// creates PATCH data with only match_vrf and match_vrf_ref_type_ fields
func (m *Manager) createMatchVrfPatchData(originalData interface{}) interface{} {
	jsonData, err := json.Marshal(originalData)
	if err != nil {
		return nil
	}

	var dataMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &dataMap); err != nil {
		return nil
	}

	patchData := openapi.RoutemapclausesPutRequestRouteMapClauseValue{}

	if name, exists := dataMap["name"]; exists {
		if nameStr, ok := name.(string); ok {
			patchData.Name = openapi.PtrString(nameStr)
		}
	}

	// Include match_vrf if it exists
	if matchVrf, exists := dataMap["match_vrf"]; exists {
		if matchVrfStr, ok := matchVrf.(string); ok {
			patchData.MatchVrf = openapi.PtrString(matchVrfStr)
		}
	}

	// Include match_vrf_ref_type_ if it exists
	if matchVrfRefType, exists := dataMap["match_vrf_ref_type_"]; exists {
		if matchVrfRefTypeStr, ok := matchVrfRefType.(string); ok {
			patchData.MatchVrfRefType = openapi.PtrString(matchVrfRefTypeStr)
		}
	}

	return patchData
}
