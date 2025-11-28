package bulkops

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

func (m *Manager) AddPut(ctx context.Context, resourceType, resourceName string, props interface{}, headerParams ...map[string]string) string {
	var params map[string]string
	if len(headerParams) > 0 {
		params = headerParams[0]
	}
	return m.addGenericOperation(ctx, resourceType, resourceName, "PUT", props, params)
}

func (m *Manager) AddPatch(ctx context.Context, resourceType, resourceName string, props interface{}, headerParams ...map[string]string) string {
	var params map[string]string
	if len(headerParams) > 0 {
		params = headerParams[0]
	}
	return m.addGenericOperation(ctx, resourceType, resourceName, "PATCH", props, params)
}

func (m *Manager) AddDelete(ctx context.Context, resourceType, resourceName string, headerParams ...map[string]string) string {
	var params map[string]string
	if len(headerParams) > 0 {
		params = headerParams[0]
	}
	return m.addGenericOperation(ctx, resourceType, resourceName, "DELETE", nil, params)
}

func (m *Manager) addGenericOperation(ctx context.Context, resourceType, resourceName, operationType string, props interface{}, headerParams map[string]string) string {
	// For resources with HeaderSplitKey, create a composite key to prevent overwrites
	storeKey := resourceName
	if config, exists := resourceRegistry[resourceType]; exists && config.HeaderSplitKey != "" && headerParams != nil {
		if headerValue, exists := headerParams[config.HeaderSplitKey]; exists && headerValue != "" {
			storeKey = fmt.Sprintf("%s_%s%s", resourceName, config.HeaderSplitKey, headerValue)
		}
	}

	storeFunc := func() {
		m.storeOperation(resourceType, storeKey, operationType, props)
		// Store header params and original name for resources that need them
		if len(headerParams) > 0 {
			if m.resourceHeaderParams == nil {
				m.resourceHeaderParams = make(map[string]map[string]string)
			}
			if m.resourceOriginalNames == nil {
				m.resourceOriginalNames = make(map[string]string)
			}
			paramKey := fmt.Sprintf("%s:%s", resourceType, storeKey)
			m.resourceHeaderParams[paramKey] = headerParams
			// Store original name if composite key was created
			if storeKey != resourceName {
				m.resourceOriginalNames[paramKey] = resourceName
			}
		}
	}

	logDetails := map[string]interface{}{
		fmt.Sprintf("%s_name", resourceType): resourceName,
		"batch_size":                         m.getBatchSize(resourceType, operationType) + 1,
	}

	for key, value := range headerParams {
		logDetails[key] = value
	}

	return m.addOperation(ctx, resourceType, resourceName, operationType, storeFunc, logDetails)
}

func (m *Manager) ExecuteBulk(ctx context.Context, resourceType, operationType string) diag.Diagnostics {
	config, exists := resourceRegistry[resourceType]
	if !exists {
		var diags diag.Diagnostics
		diags.AddError("Unknown resource type", fmt.Sprintf("Resource type %s is not registered", resourceType))
		return diags
	}

	// Check if this resource type needs header-based operation splitting
	if config.HeaderSplitKey != "" {
		return m.executeBulkWithHeaderSplit(ctx, resourceType, operationType, config)
	}

	return m.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:      resourceType,
		OperationType:     operationType,
		ExtractOperations: m.createExtractor(resourceType, operationType),
		CheckPreExistence: m.createPreExistenceChecker(config, operationType),
		PrepareRequest:    m.createRequestPreparer(config, operationType),
		ExecuteRequest:    m.createRequestExecutor(config, operationType),
		ProcessResponse:   m.createResponseProcessor(config, operationType),
		UpdateRecentOps:   m.createRecentOpsUpdater(resourceType),
	})
}

// executeBulkWithHeaderSplit handles operations for resources that need header-based batch splitting
// Resources with HeaderSplitKey will have their operations grouped by that header parameter value
// Example: ACLs split by "ip_version" into separate IPv4 and IPv6 batches
func (m *Manager) executeBulkWithHeaderSplit(ctx context.Context, resourceType, operationType string, config ResourceConfig) diag.Diagnostics {
	var diagnostics diag.Diagnostics
	m.mutex.Lock()

	res, exists := m.resources[resourceType]
	if !exists {
		m.mutex.Unlock()
		return diagnostics
	}

	// Extract operations based on type
	var originalOperations map[string]interface{}
	switch operationType {
	case "PUT":
		if res.Put == nil {
			m.mutex.Unlock()
			return diagnostics
		}
		originalOperations = make(map[string]interface{})
		for k, v := range res.Put {
			originalOperations[k] = v
		}
		res.Put = make(map[string]interface{})
	case "PATCH":
		if res.Patch == nil {
			m.mutex.Unlock()
			return diagnostics
		}
		originalOperations = make(map[string]interface{})
		for k, v := range res.Patch {
			originalOperations[k] = v
		}
		res.Patch = make(map[string]interface{})
	case "DELETE":
		if len(res.Delete) == 0 {
			m.mutex.Unlock()
			return diagnostics
		}
		originalOperations = make(map[string]interface{})
		for _, name := range res.Delete {
			// Use empty struct as placeholder for DELETE operations
			originalOperations[name] = struct{}{}
		}
		res.Delete = res.Delete[:0]
	}

	// Extract header values and original names for all operations
	headerValues := make(map[string]string)
	originalNames := make(map[string]string)
	for k := range originalOperations {
		paramKey := fmt.Sprintf("%s:%s", resourceType, k)
		if params, exists := m.resourceHeaderParams[paramKey]; exists {
			if headerValue, ok := params[config.HeaderSplitKey]; ok {
				headerValues[k] = headerValue
			}
		}
		// Get original name if it was stored
		if origName, exists := m.resourceOriginalNames[paramKey]; exists {
			originalNames[k] = origName
		} else {
			// Composite key not used, k is the original name
			originalNames[k] = k
		}
		// Clean up header params and original names
		delete(m.resourceHeaderParams, paramKey)
		delete(m.resourceOriginalNames, paramKey)
	}
	m.mutex.Unlock()

	if len(originalOperations) == 0 {
		return diagnostics
	}

	// Group operations by header value using original resource names
	groupedOps := make(map[string]map[string]interface{})
	for compositeKey, props := range originalOperations {
		headerValue := headerValues[compositeKey]
		originalName := originalNames[compositeKey]

		if groupedOps[headerValue] == nil {
			groupedOps[headerValue] = make(map[string]interface{})
		}
		groupedOps[headerValue][originalName] = props
	}

	// Execute operations for each header value group
	for headerValue, ops := range groupedOps {
		headers := map[string]string{config.HeaderSplitKey: headerValue}
		groupDiags := m.executeOperationsWithHeaders(ctx, resourceType, operationType, ops, headers, config)
		diagnostics = append(diagnostics, groupDiags...)
	}

	// Update recent operations
	res.RecentOps = true
	res.RecentOpTime = time.Now()

	return diagnostics
}

// executeOperationsWithHeaders executes a batch of operations with specific header parameters
func (m *Manager) executeOperationsWithHeaders(ctx context.Context, resourceType, operationType string, operations map[string]interface{}, headers map[string]string, config ResourceConfig) diag.Diagnostics {
	return m.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  resourceType,
		OperationType: operationType,

		ExtractOperations: func() (map[string]interface{}, []string) {
			names := make([]string, 0, len(operations))
			for k := range operations {
				names = append(names, k)
			}
			return operations, names
		},

		CheckPreExistence: func(ctx context.Context, resourceNames []string, originalOperations map[string]interface{}) ([]string, map[string]interface{}, error) {
			if operationType != "PUT" || config.HeaderGetFunc == nil {
				return nil, nil, nil
			}

			checker := ResourceExistenceCheck{
				ResourceType:  resourceType,
				OperationType: "PUT",
				FetchResources: func(ctx context.Context) (map[string]interface{}, error) {
					apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
					defer cancel()

					resp, err := config.HeaderGetFunc(m.client, apiCtx, headers)
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result map[string]interface{}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}

					// Use custom response extractor if configured
					if config.HeaderResponseExtractor != nil {
						return config.HeaderResponseExtractor(result, headers)
					}

					return result, nil
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
		},

		PrepareRequest: func(filteredData map[string]interface{}) interface{} {
			if operationType == "DELETE" {
				names := make([]string, 0, len(filteredData))
				for name := range filteredData {
					names = append(names, name)
				}
				return names
			}
			return m.createRequestPreparer(config, operationType)(filteredData)
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			switch operationType {
			case "PUT":
				if config.HeaderPutFunc != nil {
					return config.HeaderPutFunc(m.client, ctx, request, headers)
				}
			case "PATCH":
				if config.HeaderPatchFunc != nil {
					return config.HeaderPatchFunc(m.client, ctx, request, headers)
				}
			case "DELETE":
				if config.HeaderDeleteFunc != nil {
					return config.HeaderDeleteFunc(m.client, ctx, request.([]string), headers)
				}
			}
			return nil, fmt.Errorf("unsupported operation type for header-aware resource: %s", operationType)
		},

		ProcessResponse: nil,

		UpdateRecentOps: func() {
			// Already handled in parent function
		},
	})
}
