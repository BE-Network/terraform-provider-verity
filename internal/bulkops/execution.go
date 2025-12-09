package bulkops

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"terraform-provider-verity/internal/utils"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ================================================================================================
// OPERATION EXECUTION
// ================================================================================================

// GetResourceOperationData returns operation data for a resource type
func (m *Manager) GetResourceOperationData(resourceType string) *ResourceOperationData {
	// Handle special cases and aliases
	switch resourceType {
	case "acl_v4", "acl_v6":
		resourceType = "acl"
	}

	// Get the resource operations from the unified map
	res, exists := m.resources[resourceType]
	if !exists {
		return nil
	}

	// Return a ResourceOperationData that points to the unified structure fields
	return &ResourceOperationData{
		PutOperations:    res.Put,
		PatchOperations:  res.Patch,
		DeleteOperations: &res.Delete,
		RecentOps:        &res.RecentOps,
		RecentOpTime:     &res.RecentOpTime,
	}

}

func (m *Manager) hasPendingOrRecentOperations(
	resourceType string,
) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	data := m.GetResourceOperationData(resourceType)
	if data == nil {
		return false
	}

	// Check if any operations are pending
	var hasPending bool
	if putMap, ok := data.PutOperations.(map[string]interface{}); ok {
		hasPending = hasPending || len(putMap) > 0
	} else {
		if v := reflect.ValueOf(data.PutOperations); v.IsValid() && !v.IsNil() {
			hasPending = hasPending || v.Len() > 0
		}
	}

	if patchMap, ok := data.PatchOperations.(map[string]interface{}); ok {
		hasPending = hasPending || len(patchMap) > 0
	} else {
		if v := reflect.ValueOf(data.PatchOperations); v.IsValid() && !v.IsNil() {
			hasPending = hasPending || v.Len() > 0
		}
	}

	hasPending = hasPending || (data.DeleteOperations != nil && len(*data.DeleteOperations) > 0)

	// Check if we've recently had operations (within the last 5 seconds)
	hasRecent := *data.RecentOps && time.Since(*data.RecentOpTime) < 5*time.Second

	return hasPending || hasRecent
}

func (m *Manager) addOperation(
	ctx context.Context,
	resourceType string,
	resourceName string,
	operationType string,
	storeFunc func(),
	logDetails map[string]interface{},
) string {
	storeFunc()

	operationID := generateOperationID(resourceType, resourceName, operationType)
	m.operationMutex.Lock()
	defer m.operationMutex.Unlock()

	m.pendingOperations[operationID] = &Operation{
		ResourceType:  resourceType,
		ResourceName:  resourceName,
		OperationType: operationType,
		Status:        OperationPending,
	}

	m.operationWaitChannels[operationID] = make(chan struct{})

	now := time.Now()
	m.lastOperationTime = now
	if m.batchStartTime.IsZero() {
		m.batchStartTime = now
	}

	if logDetails != nil {
		logDetails["operation_id"] = operationID
		tflog.Debug(ctx, fmt.Sprintf("Added %s to %s batch", resourceType, operationType), logDetails)
	}

	return operationID
}

func (m *Manager) executeBulkOperation(ctx context.Context, config BulkOperationConfig) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	operations, resourceNames := config.ExtractOperations()

	if len(operations) == 0 {
		return diagnostics
	}

	// For DELETE operations with many resources, batch them to avoid URL length limits
	// DELETE operations use query parameters which can exceed server URL limits (~8KB for Apache)
	if config.OperationType == "DELETE" && len(resourceNames) > MaxDeleteBatchSize {
		return m.executeBatchedDeleteOperation(ctx, config, operations, resourceNames)
	}

	// For PUT operations, filter out resources that already exist
	var filteredOperations map[string]interface{}
	var filteredResourceNames []string

	if config.OperationType == "PUT" && config.CheckPreExistence != nil {
		var err error
		filteredResourceNames, filteredOperations, err = config.CheckPreExistence(ctx, resourceNames, operations)
		if err != nil {
			tflog.Warn(ctx, fmt.Sprintf("Error checking for existing %s: %v - proceeding with all resources",
				config.ResourceType, err))
			filteredResourceNames = resourceNames
			filteredOperations = operations
		}

		if len(filteredOperations) == 0 && len(filteredResourceNames) > 0 {
			filteredOperations = make(map[string]interface{})
			for _, name := range filteredResourceNames {
				if val, ok := operations[name]; ok {
					filteredOperations[name] = val
				}
			}
		}

		if len(filteredResourceNames) == 0 {
			tflog.Info(ctx, fmt.Sprintf("All %s already exist, skipping bulk %s operation",
				config.ResourceType, config.OperationType))
			config.UpdateRecentOps()
			return diagnostics
		}
	} else {
		filteredOperations = operations
		filteredResourceNames = resourceNames
	}

	tflog.Debug(ctx, fmt.Sprintf("Executing bulk %s %s operation", config.ResourceType, config.OperationType),
		map[string]interface{}{
			fmt.Sprintf("%s_count", config.ResourceType): len(filteredOperations),
			fmt.Sprintf("%s_names", config.ResourceType): filteredResourceNames,
		})

	request := config.PrepareRequest(filteredOperations)

	retryConfig := utils.DefaultRetryConfig()
	var opErr error
	var apiResp *http.Response

	for retry := 0; retry < retryConfig.MaxRetries; retry++ {
		if retry > 0 {
			delayTime := utils.CalculateBackoff(retry, retryConfig)
			tflog.Debug(ctx, fmt.Sprintf("Retrying bulk %s %s operation after %v",
				config.ResourceType, config.OperationType, delayTime))
			time.Sleep(delayTime)
		}

		apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
		apiResp, opErr = config.ExecuteRequest(apiCtx, request)
		cancel()

		if opErr == nil {
			break
		}

		if !utils.IsRetriableError(opErr) {
			break
		}

		delayTime := utils.CalculateBackoff(retry, retryConfig)
		tflog.Debug(ctx, fmt.Sprintf("Bulk %s %s operation failed with retriable error, retrying",
			config.ResourceType, config.OperationType),
			map[string]interface{}{
				"attempt":     retry + 1,
				"error":       opErr.Error(),
				"delay_ms":    delayTime.Milliseconds(),
				"max_retries": retryConfig.MaxRetries,
			})
	}

	if opErr == nil && apiResp != nil && config.ProcessResponse != nil {
		if processErr := config.ProcessResponse(ctx, apiResp); processErr != nil {
			tflog.Warn(ctx, fmt.Sprintf("Post-processing failed for bulk %s %s operation: %v",
				config.ResourceType, config.OperationType, processErr))
		}
	}

	m.updateOperationStatuses(ctx, config.ResourceType, config.OperationType, filteredResourceNames, opErr)

	if opErr != nil {
		diagnostics.AddError(
			fmt.Sprintf("Failed to execute bulk %s %s operation", config.ResourceType, config.OperationType),
			fmt.Sprintf("Error: %s", opErr),
		)
		return diagnostics
	}

	config.UpdateRecentOps()
	return diagnostics
}

// executeBatchedDeleteOperation handles DELETE operations that exceed MaxDeleteBatchSize
// by splitting them into smaller batches to avoid URL length limits.
func (m *Manager) executeBatchedDeleteOperation(ctx context.Context, config BulkOperationConfig, operations map[string]interface{}, resourceNames []string) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	totalResources := len(resourceNames)
	batchCount := (totalResources + MaxDeleteBatchSize - 1) / MaxDeleteBatchSize

	tflog.Info(ctx, fmt.Sprintf("Splitting bulk %s DELETE into %d batches of max %d resources each (total: %d)",
		config.ResourceType, batchCount, MaxDeleteBatchSize, totalResources))

	// Process each batch
	for batchNum := 0; batchNum < batchCount; batchNum++ {
		start := batchNum * MaxDeleteBatchSize
		end := start + MaxDeleteBatchSize
		if end > totalResources {
			end = totalResources
		}

		batchNames := resourceNames[start:end]
		batchOperations := make(map[string]interface{})
		for _, name := range batchNames {
			if val, exists := operations[name]; exists {
				batchOperations[name] = val
			}
		}

		tflog.Debug(ctx, fmt.Sprintf("Executing DELETE batch %d/%d for %s", batchNum+1, batchCount, config.ResourceType),
			map[string]interface{}{
				"batch_size":     len(batchNames),
				"resource_names": batchNames,
			})

		// Create a batch-specific config that returns only this batch's operations
		batchConfig := BulkOperationConfig{
			ResourceType:  config.ResourceType,
			OperationType: config.OperationType,
			ExtractOperations: func() (map[string]interface{}, []string) {
				return batchOperations, batchNames
			},
			CheckPreExistence: config.CheckPreExistence,
			PrepareRequest:    config.PrepareRequest,
			ExecuteRequest:    config.ExecuteRequest,
			ProcessResponse:   config.ProcessResponse,
			UpdateRecentOps:   func() {}, // Don't update until all batches complete
		}

		// Execute this batch using the standard execution path (won't recurse since batch size <= MaxDeleteBatchSize)
		batchDiags := m.executeSingleDeleteBatch(ctx, batchConfig, batchOperations, batchNames)
		diagnostics.Append(batchDiags...)

		if batchDiags.HasError() {
			tflog.Error(ctx, fmt.Sprintf("DELETE batch %d/%d failed for %s, stopping further batches",
				batchNum+1, batchCount, config.ResourceType))
			return diagnostics
		}

		// Small delay between batches to avoid overwhelming the server
		if batchNum < batchCount-1 {
			time.Sleep(100 * time.Millisecond)
		}
	}

	// Update recent ops after all batches complete successfully
	config.UpdateRecentOps()

	tflog.Info(ctx, fmt.Sprintf("Successfully completed all %d DELETE batches for %s", batchCount, config.ResourceType))
	return diagnostics
}

// executeSingleDeleteBatch executes a single batch of DELETE operations
func (m *Manager) executeSingleDeleteBatch(ctx context.Context, config BulkOperationConfig, operations map[string]interface{}, resourceNames []string) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	tflog.Debug(ctx, fmt.Sprintf("Executing bulk %s %s operation", config.ResourceType, config.OperationType),
		map[string]interface{}{
			fmt.Sprintf("%s_count", config.ResourceType): len(operations),
			fmt.Sprintf("%s_names", config.ResourceType): resourceNames,
		})

	request := config.PrepareRequest(operations)

	retryConfig := utils.DefaultRetryConfig()
	var opErr error
	var apiResp *http.Response

	for retry := 0; retry < retryConfig.MaxRetries; retry++ {
		if retry > 0 {
			delayTime := utils.CalculateBackoff(retry, retryConfig)
			tflog.Debug(ctx, fmt.Sprintf("Retrying bulk %s %s operation after %v",
				config.ResourceType, config.OperationType, delayTime))
			time.Sleep(delayTime)
		}

		apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
		apiResp, opErr = config.ExecuteRequest(apiCtx, request)
		cancel()

		if opErr == nil {
			break
		}

		if !utils.IsRetriableError(opErr) {
			break
		}

		delayTime := utils.CalculateBackoff(retry, retryConfig)
		tflog.Debug(ctx, fmt.Sprintf("Bulk %s %s operation failed with retriable error, retrying",
			config.ResourceType, config.OperationType),
			map[string]interface{}{
				"attempt":     retry + 1,
				"error":       opErr.Error(),
				"delay_ms":    delayTime.Milliseconds(),
				"max_retries": retryConfig.MaxRetries,
			})
	}

	if opErr == nil && apiResp != nil && config.ProcessResponse != nil {
		if processErr := config.ProcessResponse(ctx, apiResp); processErr != nil {
			tflog.Warn(ctx, fmt.Sprintf("Post-processing failed for bulk %s %s operation: %v",
				config.ResourceType, config.OperationType, processErr))
		}
	}

	m.updateOperationStatuses(ctx, config.ResourceType, config.OperationType, resourceNames, opErr)

	if opErr != nil {
		diagnostics.AddError(
			fmt.Sprintf("Failed to execute bulk %s %s operation", config.ResourceType, config.OperationType),
			fmt.Sprintf("Error: %s", opErr),
		)
	}

	return diagnostics
}

func generateOperationID(resourceType, resourceName, operationType string) string {
	return fmt.Sprintf("%s-%s-%s-%s", resourceType, resourceName, operationType, uuid.New().String())
}

func (m *Manager) WaitForOperation(ctx context.Context, operationID string, timeout time.Duration) error {
	m.operationMutex.Lock()
	waitCh, exists := m.operationWaitChannels[operationID]
	if !exists {
		m.operationMutex.Unlock()
		return fmt.Errorf("operation %s not found", operationID)
	}

	if closed, ok := m.closedChannels[operationID]; ok && closed {
		var err error
		if errorVal, hasError := m.operationErrors[operationID]; hasError {
			err = errorVal
		}
		m.operationMutex.Unlock()
		return err
	}
	m.operationMutex.Unlock()

	select {
	case <-waitCh:
		// Operation completed
		m.operationMutex.Lock()
		defer m.operationMutex.Unlock()

		if err, hasError := m.operationErrors[operationID]; hasError {
			return err
		}
		return nil

	case <-time.After(timeout):
		return fmt.Errorf("timeout waiting for operation %s", operationID)

	case <-ctx.Done():
		return ctx.Err()
	}
}

// updateOperationStatuses updates the status of pending operations based on the bulk operation result
func (m *Manager) updateOperationStatuses(ctx context.Context, resourceType, operationType string, resourceNames []string, opErr error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	resourceMap := make(map[string]bool)
	for _, name := range resourceNames {
		resourceMap[name] = true
	}

	for opID, op := range m.pendingOperations {
		matchesResourceType := false
		if resourceType == "acl_v4" || resourceType == "acl_v6" {
			// For ACL operations, match if the stored type is "acl"
			matchesResourceType = op.ResourceType == "acl"
		} else {
			matchesResourceType = op.ResourceType == resourceType
		}

		if matchesResourceType && op.OperationType == operationType {
			// For PUT, we need to check pending status
			if operationType == "PUT" && op.Status != OperationPending {
				continue
			}

			// Check if this operation's resource name is in our filtered batch
			if resourceMap[op.ResourceName] {
				updatedOp := op
				if opErr == nil {
					// Mark operation as successful
					updatedOp.Status = OperationSucceeded
					m.pendingOperations[opID] = updatedOp
					m.operationResults[opID] = true
				} else {
					// Mark operation as failed
					updatedOp.Status = OperationFailed
					updatedOp.Error = opErr
					m.pendingOperations[opID] = updatedOp
					m.operationErrors[opID] = opErr
					m.operationResults[opID] = false
				}
				m.safeCloseChannel(opID, true)
			}
		}
	}
}

// ================================================================================================
// BULK OPERATION EXECUTION ORCHESTRATION
// ================================================================================================

func (m *Manager) ExecuteAllPendingOperations(ctx context.Context) diag.Diagnostics {
	// Ensure only one execution runs at a time - prevents race conditions when multiple
	// timer callbacks fire due to low parallelism causing resource waves
	m.executionMutex.Lock()
	defer m.executionMutex.Unlock()

	var diagnostics diag.Diagnostics

	if time.Since(m.lastOperationTime) < BatchCollectionWindow {
		remaining := BatchCollectionWindow - time.Since(m.lastOperationTime)
		tflog.Debug(ctx, fmt.Sprintf("Waiting %v to collect more operations before executing", remaining))
		time.Sleep(remaining)
	}

	var opsDiags diag.Diagnostics
	var operationsPerformed bool

	tflog.Debug(ctx, "Executing pending operations in 'datacenter' mode")
	opsDiags, operationsPerformed = m.ExecuteDatacenterOperations(ctx)

	diagnostics.Append(opsDiags...)

	if operationsPerformed {
		waitDuration := 800 * time.Millisecond
		tflog.Debug(ctx, fmt.Sprintf("Waiting %v for all operations to propagate before final cache refresh", waitDuration))
		time.Sleep(waitDuration)

		tflog.Debug(ctx, "Final cache clear after all operations to ensure verification with fresh data")
		if m.clearCacheFunc != nil && m.contextProvider != nil {
			m.clearCacheFunc(ctx, m.contextProvider(), "tenants")
			m.clearCacheFunc(ctx, m.contextProvider(), "gateways")
			m.clearCacheFunc(ctx, m.contextProvider(), "gateway_profiles")
			m.clearCacheFunc(ctx, m.contextProvider(), "services")
			m.clearCacheFunc(ctx, m.contextProvider(), "eth_port_profiles")
			m.clearCacheFunc(ctx, m.contextProvider(), "eth_port_settings")
			m.clearCacheFunc(ctx, m.contextProvider(), "lags")
			m.clearCacheFunc(ctx, m.contextProvider(), "bundles")
		}
	}

	return diagnostics
}

func (m *Manager) ExecuteDatacenterOperations(ctx context.Context) (diag.Diagnostics, bool) {
	var diagnostics diag.Diagnostics
	operationsPerformed := false

	execute := func(opType string, count int, execFunc func(context.Context) diag.Diagnostics, resourceName string) bool {
		if count > 0 {
			tflog.Debug(ctx, fmt.Sprintf("Executing %s %s operations", resourceName, opType), map[string]interface{}{
				"operation_count": count,
			})
			diags := execFunc(ctx)
			diagnostics.Append(diags...)
			if diags.HasError() {
				err := fmt.Errorf("bulk %s %s operation failed", resourceName, opType)
				m.FailAllPendingOperations(ctx, err)
				return false
			}
			operationsPerformed = true
		}
		return true
	}

	// PUT operations (Bundles do not support PUT)
	// 1. Tenants
	if !execute("PUT", m.getOperationCount("tenant", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "tenant", "PUT") }, "Tenant") {
		return diagnostics, operationsPerformed
	}
	// 2. Gateways
	if !execute("PUT", m.getOperationCount("gateway", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "gateway", "PUT") }, "Gateway") {
		return diagnostics, operationsPerformed
	}
	// 3. Gateway Profiles
	if !execute("PUT", m.getOperationCount("gateway_profile", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "gateway_profile", "PUT") }, "Gateway Profile") {
		return diagnostics, operationsPerformed
	}
	// 4. Services
	if !execute("PUT", m.getOperationCount("service", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "service", "PUT") }, "Service") {
		return diagnostics, operationsPerformed
	}
	// 5. Eth Port Profiles
	if !execute("PUT", m.getOperationCount("eth_port_profile", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "eth_port_profile", "PUT") }, "Eth Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 6. Eth Port Settings
	if !execute("PUT", m.getOperationCount("eth_port_settings", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "eth_port_settings", "PUT") }, "Eth Port Settings") {
		return diagnostics, operationsPerformed
	}
	// 7. Lags
	if !execute("PUT", m.getOperationCount("lag", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "lag", "PUT") }, "LAG") {
		return diagnostics, operationsPerformed
	}

	// PATCH operations
	// 1. Tenants
	if !execute("PATCH", m.getOperationCount("tenant", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "tenant", "PATCH") }, "Tenant") {
		return diagnostics, operationsPerformed
	}
	// 2. Gateways
	if !execute("PATCH", m.getOperationCount("gateway", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "gateway", "PATCH") }, "Gateway") {
		return diagnostics, operationsPerformed
	}
	// 3. Gateway Profiles
	if !execute("PATCH", m.getOperationCount("gateway_profile", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "gateway_profile", "PATCH") }, "Gateway Profile") {
		return diagnostics, operationsPerformed
	}
	// 4. Services
	if !execute("PATCH", m.getOperationCount("service", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "service", "PATCH") }, "Service") {
		return diagnostics, operationsPerformed
	}
	// 5. Eth Port Profiles
	if !execute("PATCH", m.getOperationCount("eth_port_profile", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "eth_port_profile", "PATCH") }, "Eth Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 6. Eth Port Settings
	if !execute("PATCH", m.getOperationCount("eth_port_settings", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "eth_port_settings", "PATCH") }, "Eth Port Settings") {
		return diagnostics, operationsPerformed
	}
	// 7. Lags
	if !execute("PATCH", m.getOperationCount("lag", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "lag", "PATCH") }, "LAG") {
		return diagnostics, operationsPerformed
	}
	// 8. Bundles
	if !execute("PATCH", m.getOperationCount("bundle", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "bundle", "PATCH") }, "Bundle") {
		return diagnostics, operationsPerformed
	}

	// DELETE operations - Reverse order (Bundles do not support DELETE)
	// 7. Lags
	if !execute("DELETE", m.getOperationCount("lag", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "lag", "DELETE") }, "LAG") {
		return diagnostics, operationsPerformed
	}
	// 6. Eth Port Settings
	if !execute("DELETE", m.getOperationCount("eth_port_settings", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "eth_port_settings", "DELETE") }, "Eth Port Settings") {
		return diagnostics, operationsPerformed
	}
	// 5. Eth Port Profiles
	if !execute("DELETE", m.getOperationCount("eth_port_profile", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "eth_port_profile", "DELETE") }, "Eth Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 4. Services
	if !execute("DELETE", m.getOperationCount("service", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "service", "DELETE") }, "Service") {
		return diagnostics, operationsPerformed
	}
	// 3. Gateway Profiles
	if !execute("DELETE", m.getOperationCount("gateway_profile", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "gateway_profile", "DELETE") }, "Gateway Profile") {
		return diagnostics, operationsPerformed
	}
	// 2. Gateways
	if !execute("DELETE", m.getOperationCount("gateway", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "gateway", "DELETE") }, "Gateway") {
		return diagnostics, operationsPerformed
	}
	// 1. Tenants
	if !execute("DELETE", m.getOperationCount("tenant", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "tenant", "DELETE") }, "Tenant") {
		return diagnostics, operationsPerformed
	}

	return diagnostics, operationsPerformed
}

func (m *Manager) ShouldExecuteOperations(ctx context.Context) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// If there are no pending operations, no need to execute
	if m.getOperationCount("gateway", "PUT") == 0 && m.getOperationCount("gateway", "PATCH") == 0 && m.getOperationCount("gateway", "DELETE") == 0 &&
		m.getOperationCount("lag", "PUT") == 0 && m.getOperationCount("lag", "PATCH") == 0 && m.getOperationCount("lag", "DELETE") == 0 &&
		m.getOperationCount("tenant", "PUT") == 0 && m.getOperationCount("tenant", "PATCH") == 0 && m.getOperationCount("tenant", "DELETE") == 0 &&
		m.getOperationCount("service", "PUT") == 0 && m.getOperationCount("service", "PATCH") == 0 && m.getOperationCount("service", "DELETE") == 0 &&
		m.getOperationCount("gateway_profile", "PUT") == 0 && m.getOperationCount("gateway_profile", "PATCH") == 0 && m.getOperationCount("gateway_profile", "DELETE") == 0 &&
		m.getOperationCount("eth_port_profile", "PUT") == 0 && m.getOperationCount("eth_port_profile", "PATCH") == 0 && m.getOperationCount("eth_port_profile", "DELETE") == 0 &&
		m.getOperationCount("eth_port_settings", "PUT") == 0 && m.getOperationCount("eth_port_settings", "PATCH") == 0 && m.getOperationCount("eth_port_settings", "DELETE") == 0 &&
		m.getOperationCount("bundle", "PATCH") == 0 {
		return false
	}

	elapsedSinceLast := time.Since(m.lastOperationTime)
	elapsedSinceBatchStart := time.Since(m.batchStartTime)

	// Only flush if either sufficient time has passed since the last operation
	// OR the batch has been open for too long
	if elapsedSinceLast < BatchCollectionWindow && elapsedSinceBatchStart < MaxBatchDelay {
		return false
	}

	return true
}

func (m *Manager) ExecuteIfMultipleOperations(ctx context.Context) diag.Diagnostics {
	m.mutex.Lock()
	gatewayPutCount := m.getOperationCount("gateway", "PUT")
	gatewayPatchCount := m.getOperationCount("gateway", "PATCH")
	gatewayDeleteCount := m.getOperationCount("gateway", "DELETE")

	lagPutCount := m.getOperationCount("lag", "PUT")
	lagPatchCount := m.getOperationCount("lag", "PATCH")
	lagDeleteCount := m.getOperationCount("lag", "DELETE")

	tenantPutCount := m.getOperationCount("tenant", "PUT")
	tenantPatchCount := m.getOperationCount("tenant", "PATCH")
	tenantDeleteCount := m.getOperationCount("tenant", "DELETE")

	servicePutCount := m.getOperationCount("service", "PUT")
	servicePatchCount := m.getOperationCount("service", "PATCH")
	serviceDeleteCount := m.getOperationCount("service", "DELETE")

	gatewayProfilePutCount := m.getOperationCount("gateway_profile", "PUT")
	gatewayProfilePatchCount := m.getOperationCount("gateway_profile", "PATCH")
	gatewayProfileDeleteCount := m.getOperationCount("gateway_profile", "DELETE")

	ethPortProfilePutCount := m.getOperationCount("eth_port_profile", "PUT")
	ethPortProfilePatchCount := m.getOperationCount("eth_port_profile", "PATCH")
	ethPortProfileDeleteCount := m.getOperationCount("eth_port_profile", "DELETE")

	ethPortSettingsPutCount := m.getOperationCount("eth_port_settings", "PUT")
	ethPortSettingsPatchCount := m.getOperationCount("eth_port_settings", "PATCH")
	ethPortSettingsDeleteCount := m.getOperationCount("eth_port_settings", "DELETE")

	bundlePatchCount := m.getOperationCount("bundle", "PATCH")

	m.mutex.Unlock()

	totalCount := gatewayPutCount + gatewayPatchCount + gatewayDeleteCount +
		lagPutCount + lagPatchCount + lagDeleteCount +
		tenantPutCount + tenantPatchCount + tenantDeleteCount +
		servicePutCount + servicePatchCount + serviceDeleteCount +
		gatewayProfilePutCount + gatewayProfilePatchCount + gatewayProfileDeleteCount +
		ethPortProfilePutCount + ethPortProfilePatchCount + ethPortProfileDeleteCount +
		ethPortSettingsPutCount + ethPortSettingsPatchCount + ethPortSettingsDeleteCount +
		bundlePatchCount

	if totalCount > 0 {
		tflog.Debug(ctx, "Multiple operations detected, executing in sequence", map[string]interface{}{
			"gateway_put_count":              gatewayPutCount,
			"gateway_patch_count":            gatewayPatchCount,
			"gateway_delete_count":           gatewayDeleteCount,
			"lag_put_count":                  lagPutCount,
			"lag_patch_count":                lagPatchCount,
			"lag_delete_count":               lagDeleteCount,
			"tenant_put_count":               tenantPutCount,
			"tenant_patch_count":             tenantPatchCount,
			"tenant_delete_count":            tenantDeleteCount,
			"service_put_count":              servicePutCount,
			"service_patch_count":            servicePatchCount,
			"service_delete_count":           serviceDeleteCount,
			"gateway_profile_put_count":      gatewayProfilePutCount,
			"gateway_profile_patch_count":    gatewayProfilePatchCount,
			"gateway_profile_delete_count":   gatewayProfileDeleteCount,
			"eth_port_profile_put_count":     ethPortProfilePutCount,
			"eth_port_profile_patch_count":   ethPortProfilePatchCount,
			"eth_port_profile_delete_count":  ethPortProfileDeleteCount,
			"eth_port_settings_put_count":    ethPortSettingsPutCount,
			"eth_port_settings_patch_count":  ethPortSettingsPatchCount,
			"eth_port_settings_delete_count": ethPortSettingsDeleteCount,
			"bundle_patch_count":             bundlePatchCount,
			"total_count":                    totalCount,
		})

		return m.ExecuteAllPendingOperations(ctx)
	}

	return nil
}
