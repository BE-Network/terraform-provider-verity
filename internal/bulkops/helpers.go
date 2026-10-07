package bulkops

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"terraform-provider-verity/internal/utils"
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

			res.Delete = res.Delete[:0]
			return result, names
		}

		return make(map[string]interface{}), []string{}
	}
}

func (m *Manager) createPreExistenceChecker(config ResourceConfig, operationType string) func(context.Context, []string, map[string]interface{}) ([]string, map[string]interface{}, error) {
	if operationType != "PUT" {
		return nil
	}

	return func(ctx context.Context, resourceNames []string, originalOperations map[string]interface{}) ([]string, map[string]interface{}, error) {
		if config.GetFunc == nil || config.ResponseCollectionKey == "" {
			return resourceNames, nil, fmt.Errorf("%s has no pre-existence binding", config.ResourceType)
		}
		checker := ResourceExistenceCheck{
			ResourceType:  config.ResourceType,
			OperationType: "PUT",
			FetchResources: func(ctx context.Context) (map[string]interface{}, error) {
				apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
				defer cancel()

				resp, err := config.GetFunc(m.client, apiCtx)
				if err != nil {
					return nil, err
				}
				defer resp.Body.Close()
				var result map[string]json.RawMessage
				if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
					return nil, err
				}
				var resources map[string]interface{}
				if data, exists := result[config.ResponseCollectionKey]; exists {
					if err := json.Unmarshal(data, &resources); err != nil {
						return nil, err
					}
				}
				return resources, nil
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
	return func(data map[string]interface{}) interface{} {
		request, _ := m.createRequestPreparerWithError(config, operationType)(data)
		return request
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

func (m *Manager) createResponseProcessor(config ResourceConfig, operationType string) func(context.Context, *http.Response, map[string]interface{}) error {

	return m.createResponseProcessorWithHeaders(config, operationType, nil)
}

func (m *Manager) createHeaderAwareResponseProcessor(config ResourceConfig, operationType string, headers map[string]string) func(context.Context, *http.Response, map[string]interface{}) error {
	return m.createResponseProcessorWithHeaders(config, operationType, headers)
}

func (m *Manager) createResponseProcessorWithHeaders(config ResourceConfig, operationType string, headers map[string]string) func(context.Context, *http.Response, map[string]interface{}) error {
	return func(ctx context.Context, resp *http.Response, operations map[string]interface{}) error {
		delayTime := ResponseProcessorDelay
		tflog.Debug(ctx, fmt.Sprintf("Waiting %v for server values to be assigned before fetching %s", delayTime, config.ResourceType))
		time.Sleep(delayTime)

		fetchCtx, fetchCancel := context.WithTimeout(context.Background(), OperationTimeout)
		defer fetchCancel()

		if headers != nil {
			tflog.Debug(ctx, fmt.Sprintf("Fetching %s after successful operation to retrieve server values (headers: %v)", config.ResourceType, headers))
		} else {
			tflog.Debug(ctx, fmt.Sprintf("Fetching %s after successful operation to retrieve server values", config.ResourceType))
		}

		res, exists := m.resources[config.ResourceType]
		if !exists {
			return fmt.Errorf("resource type %s not found in unified structure", config.ResourceType)
		}

		return m.fetchVerifyAndCacheResourceResponse(fetchCtx, ctx, config, res, headers, operationType, operations)
	}
}

func (m *Manager) fetchVerifyAndCacheResourceResponse(fetchCtx context.Context, logCtx context.Context, config ResourceConfig, res *ResourceOperations, headers map[string]string, operationType string, operations map[string]interface{}) error {
	resourceData, err := m.fetchResourceResponse(fetchCtx, logCtx, config, headers)
	if err != nil {
		return err
	}

	if (operationType == "PUT" || operationType == "PATCH") && resourceData != nil {
		for attempt := 1; attempt <= PostOperationVerificationRetries; attempt++ {
			missing := missingSubmittedResponseKeys(operations, resourceData)
			if len(missing) == 0 {
				break
			}
			tflog.Warn(logCtx, fmt.Sprintf("Post-operation %s verification response is incomplete; retrying", config.ResourceType), map[string]interface{}{
				"attempt":      attempt,
				"max_retries":  PostOperationVerificationRetries,
				"missing_keys": missing,
			})
			time.Sleep(PostOperationVerificationBackoff)
			resourceData, err = m.fetchResourceResponse(fetchCtx, logCtx, config, headers)
			if err != nil {
				return err
			}
		}
	}

	m.cacheResourceResponse(logCtx, config, res, headers, resourceData)
	return nil
}

func (m *Manager) fetchResourceResponse(fetchCtx context.Context, logCtx context.Context, config ResourceConfig, headers map[string]string) (map[string]interface{}, error) {
	var getResp *http.Response
	var fetchErr error

	if headers != nil {

		if config.HeaderGetFunc == nil {
			tflog.Debug(logCtx, fmt.Sprintf("No HeaderGetFunc defined for %s, skipping response caching", config.ResourceType))
			return nil, nil
		}
		getResp, fetchErr = config.HeaderGetFunc(m.client, fetchCtx, headers)
	} else {

		if config.GetFunc == nil {
			tflog.Debug(logCtx, fmt.Sprintf("No GetFunc defined for %s, skipping response caching", config.ResourceType))
			return nil, nil
		}
		getResp, fetchErr = config.GetFunc(m.client, fetchCtx)
	}

	if fetchErr != nil {
		logFields := map[string]interface{}{"error": fetchErr.Error()}
		if headers != nil {
			logFields["headers"] = headers
		}
		tflog.Error(logCtx, fmt.Sprintf("Failed to fetch %s after operation", config.ResourceType), logFields)
		return nil, fetchErr
	}
	defer getResp.Body.Close()

	var rawResponse map[string]interface{}
	if respErr := json.NewDecoder(getResp.Body).Decode(&rawResponse); respErr != nil {
		tflog.Error(logCtx, fmt.Sprintf("Failed to decode %s response", config.ResourceType), map[string]interface{}{
			"error": respErr.Error(),
		})
		return nil, respErr
	}

	var resourceData map[string]interface{}
	if headers != nil && config.HeaderResponseExtractor != nil {
		var err error
		resourceData, err = config.HeaderResponseExtractor(rawResponse, headers)
		if err != nil {
			tflog.Error(logCtx, fmt.Sprintf("Failed to extract %s response data", config.ResourceType), map[string]interface{}{
				"error": err.Error(),
			})
			return nil, err
		}
	} else {

		jsonKey := utils.ResponseCollectionKeyForBulkKey(config.ResourceType)
		if jsonKey == "" {
			tflog.Warn(logCtx, fmt.Sprintf("No JSON key mapping found for resource type: %s", config.ResourceType))
			return nil, nil
		}
		var ok bool
		resourceData, ok = rawResponse[jsonKey].(map[string]interface{})
		if !ok {
			tflog.Debug(logCtx, fmt.Sprintf("No %s data found in response or unexpected format", jsonKey))
			return nil, nil
		}
	}
	return resourceData, nil
}

func (m *Manager) cacheResourceResponse(logCtx context.Context, config ResourceConfig, res *ResourceOperations, headers map[string]string, resourceData map[string]interface{}) {

	res.ResponsesMutex.Lock()
	for resourceName, data := range resourceData {
		if resourceMap, ok := data.(map[string]interface{}); ok {
			res.Responses[resourceName] = resourceMap

			if name, nameOk := resourceMap["name"].(string); nameOk && name != resourceName {
				res.Responses[name] = resourceMap
			}
		}
	}
	res.ResponsesMutex.Unlock()

	if headers != nil {
		tflog.Debug(logCtx, fmt.Sprintf("Successfully cached %s data (headers: %v)", config.ResourceType, headers), map[string]interface{}{
			"count": len(resourceData),
		})
	} else {
		tflog.Debug(logCtx, fmt.Sprintf("Successfully cached %s data", config.ResourceType), map[string]interface{}{
			"count": len(resourceData),
		})
	}

}

func missingSubmittedResponseKeys(operations map[string]interface{}, resourceData map[string]interface{}) map[string][]string {
	missing := make(map[string][]string)

	responseByName := make(map[string]map[string]interface{}, len(resourceData))
	for responseKey, data := range resourceData {
		response, ok := data.(map[string]interface{})
		if !ok {
			continue
		}
		responseByName[responseKey] = response
		if name, ok := response["name"].(string); ok {
			responseByName[name] = response
		}
	}
	for name, operation := range operations {
		encoded, err := json.Marshal(operation)
		if err != nil {

			continue
		}

		var submitted map[string]interface{}
		if err := json.Unmarshal(encoded, &submitted); err != nil {
			continue
		}

		response, found := responseByName[name]
		for key, value := range submitted {
			if value == nil {
				continue
			}
			if !found {
				missing[name] = append(missing[name], key)
				continue
			}
			if _, present := response[key]; !present {
				missing[name] = append(missing[name], key)
			}
		}
	}
	return missing
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

			alreadyExistingResources[name] = true
			tflog.Info(ctx, fmt.Sprintf("Skipping creation of %s '%s' as it already exists",
				checker.ResourceType, name))
		} else {

			notExistingResources = append(notExistingResources, name)
		}
	}

	if len(alreadyExistingResources) > 0 {
		m.operationMutex.Lock()
		defer m.operationMutex.Unlock()

		for opID, op := range m.pendingOperations {
			if op.ResourceType == checker.ResourceType &&
				op.OperationType == checker.OperationType &&
				alreadyExistingResources[op.ResourceName] {

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

func (m *Manager) hasPendingOperationsLocked() bool {
	for _, res := range m.resources {
		if len(res.Put) > 0 || len(res.Patch) > 0 || len(res.Delete) > 0 {
			return true
		}
	}
	return false
}

func (m *Manager) hasPendingOperations() bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.hasPendingOperationsLocked()
}

func (m *Manager) getOperationCount(resourceType, operationType string) int {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	return m.getOperationCountLocked(resourceType, operationType)
}

func (m *Manager) getOperationCountLocked(resourceType, operationType string) int {
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
