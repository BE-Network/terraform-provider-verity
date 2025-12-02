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
		case "bundle":
			// Bundles only support PATCH operations
			patchRequest := openapi.NewBundlesPatchRequest()
			bundleMap := make(map[string]openapi.BundlesPatchRequestEndpointBundleValue)
			for name, props := range filteredData {
				bundleMap[name] = props.(openapi.BundlesPatchRequestEndpointBundleValue)
			}
			patchRequest.SetEndpointBundle(bundleMap)
			return patchRequest
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
