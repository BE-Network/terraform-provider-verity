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

type orderedExecutionState struct {
	manager             *Manager
	ctx                 context.Context
	diagnostics         *diag.Diagnostics
	operationsPerformed *bool
}

var finalCacheRefreshKeys = []string{
	"tenants",
	"gateways",
	"gateway_profiles",
	"device_aaa_profiles",
	"ldap_profiles",
	"services",
	"packet_queues",
	"tacacs_profiles",
	"eth_port_profiles",
	"eth_port_settings",
	"lags",
	"sflow_collectors",
	"diagnostics_profiles",
	"diagnostics_port_profiles",
	"bundles",
	"acls_ipv4",
	"acls_ipv6",
	"packet_brokers",
	"badges",
	"switchpoints",
	"authenticated_eth_ports",
	"device_voice_settings",
	"service_port_profiles",
	"voice_port_profiles",
	"as_path_access_lists",
	"community_lists",
	"mac_filters",
	"device_settings",
	"extended_community_lists",
	"ipv4_lists",
	"ipv4_prefix_lists",
	"ipv6_lists",
	"ipv6_prefix_lists",
	"route_map_clauses",
	"route_maps",
	"sfp_breakouts",
	"fabrics",
	"planes",
	"racks",
	"pairs",
	"pods",
	"port_acls",
	"pb_routing",
	"pb_routing_acl",
	"ssp_groups",
	"spine_planes",
	"sus",
	"grouping_rules",
	"threshold_groups",
	"thresholds",
}

var datacenterPutOrder = []string{
	"extended_community_list",
	"ipv6_prefix_list",
	"community_list",
	"as_path_access_list",
	"ipv4_prefix_list",
	"route_map_clause",
	"acl",
	"route_map",
	"pb_routing_acl",
	"tenant",
	"pb_routing",
	"service",
	"ldap_profile",
	"tacacs_profile",
	"ipv6_list",
	"ipv4_list",
	"port_acl",
	"fabric",
	"packet_queue",
	"device_aaa_profile",
	"sflow_collector",
	"packet_broker",
	"eth_port_profile",
	"gateway",
	"pod",
	"device_settings",
	"diagnostics_port_profile",
	"diagnostics_profile",
	"eth_port_settings",
	"lag",
	"gateway_profile",
	"su",
	"plane",
	"bundle",
	"ssp_group",
	"badge",
	"rack",
	"spine_plane",
	"switchpoint",
	"threshold",
	"grouping_rule",
	"threshold_group",
	"pair",
}

var datacenterPatchOrder = joinOperationOrders([]string{"sfp_breakout"}, datacenterPutOrder)
var datacenterDeleteOrder = reverseOperationOrder(datacenterPutOrder)

var campusPutOrder = []string{
	"ipv4_list",
	"ipv6_list",
	"acl",
	"service",
	"port_acl",
	"mac_filter",
	"ldap_profile",
	"tacacs_profile",
	"service_port_profile",
	"fabric",
	"eth_port_profile",
	"device_aaa_profile",
	"packet_queue",
	"sflow_collector",
	"diagnostics_port_profile",
	"lag",
	"voice_port_profile",
	"device_settings",
	"eth_port_settings",
	"authenticated_eth_port",
	"device_voice_settings",
	"diagnostics_profile",
	"bundle",
	"badge",
	"switchpoint",
	"grouping_rule",
	"threshold",
	"threshold_group",
	"pair",
}

var campusPatchOrder = joinOperationOrders([]string{"sfp_breakout"}, campusPutOrder)
var campusDeleteOrder = reverseOperationOrder(campusPutOrder)

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

	// Mark all operations as executing - this sets the ExecutionStartTime so that
	// WaitForOperation can track timeout from when the API call actually starts
	m.markOperationsAsExecuting(config.ResourceType, config.OperationType, filteredResourceNames)

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

	// Mark all operations as executing - this sets the ExecutionStartTime so that
	// WaitForOperation can track timeout from when the API call actually starts
	m.markOperationsAsExecuting(config.ResourceType, config.OperationType, resourceNames)

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

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-waitCh:
			// Operation completed
			m.operationMutex.Lock()
			defer m.operationMutex.Unlock()

			if err, hasError := m.operationErrors[operationID]; hasError {
				return err
			}
			return nil

		case <-ticker.C:
			// Check if operation has started executing and if timeout has elapsed
			m.operationMutex.Lock()
			op, opExists := m.pendingOperations[operationID]
			if opExists && op != nil && !op.ExecutionStartTime.IsZero() {
				// Operation has started executing - check if timeout elapsed from start time
				elapsed := time.Since(op.ExecutionStartTime)
				if elapsed >= timeout {
					m.operationMutex.Unlock()
					return fmt.Errorf("timeout waiting for operation %s (elapsed %v since execution started)", operationID, elapsed)
				}
			}
			m.operationMutex.Unlock()

		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

// updateOperationStatuses updates the status of pending operations based on the bulk operation result
func (m *Manager) updateOperationStatuses(ctx context.Context, resourceType, operationType string, resourceNames []string, opErr error) {
	var idsToClose []string

	m.operationMutex.Lock()

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
			// For PUT, we need to check that operation is pending or executing
			if operationType == "PUT" && op.Status != OperationPending && op.Status != OperationExecuting {
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
				idsToClose = append(idsToClose, opID)
			}
		}
	}

	m.operationMutex.Unlock()

	for _, opID := range idsToClose {
		m.safeCloseChannel(opID, false)
	}
}

// markOperationsAsExecuting sets the ExecutionStartTime for all operations in the batch.
// This allows WaitForOperation to track timeout from when the API call actually begins,
// not from when the operation was queued.
func (m *Manager) markOperationsAsExecuting(resourceType, operationType string, resourceNames []string) {
	m.operationMutex.Lock()
	defer m.operationMutex.Unlock()

	startTime := time.Now()
	resourceMap := make(map[string]bool)
	for _, name := range resourceNames {
		resourceMap[name] = true
	}

	for opID, op := range m.pendingOperations {
		matchesResourceType := false
		if resourceType == "acl_v4" || resourceType == "acl_v6" {
			matchesResourceType = op.ResourceType == "acl"
		} else {
			matchesResourceType = op.ResourceType == resourceType
		}

		if matchesResourceType && op.OperationType == operationType && op.Status == OperationPending {
			if resourceMap[op.ResourceName] {
				op.Status = OperationExecuting
				op.ExecutionStartTime = startTime
				m.pendingOperations[opID] = op
			}
		}
	}
}

// WaitAndFlushAllOperations is called from stage resources acting as barriers between
// resource type groups. Waits for sibling resources to queue ops, flushes them, and
// returns only after consecutive quiet periods confirm no more ops will arrive.
func (m *Manager) WaitAndFlushAllOperations(ctx context.Context) diag.Diagnostics {
	var allDiags diag.Diagnostics

	// Phase 1: Wait for operations to arrive from concurrent sibling resources.
	initialWaitDeadline := MaxBatchDelay
	checkInterval := 500 * time.Millisecond
	elapsed := time.Duration(0)

	tflog.Info(ctx, "[BULK-OPS] Stage barrier activated — waiting for operations from concurrent resources")

	for elapsed < initialWaitDeadline {
		if m.hasPendingOperations() {
			tflog.Debug(ctx, fmt.Sprintf("Stage barrier: operations detected after %v", elapsed))
			break
		}
		time.Sleep(checkInterval)
		elapsed += checkInterval
	}

	if !m.hasPendingOperations() {
		tflog.Info(ctx, "[BULK-OPS] Stage barrier: no operations arrived during initial wait — proceeding to next stage")
		return allDiags
	}

	// Phase 2: Repeatedly flush until quiet (requires consecutive quiet periods).
	const maxFlushCycles = 200
	const requiredQuietPeriods = 3
	consecutiveQuiet := 0

	for cycle := 0; cycle < maxFlushCycles; cycle++ {
		tflog.Debug(ctx, fmt.Sprintf("Stage barrier: flush cycle %d (quiet streak: %d/%d)", cycle, consecutiveQuiet, requiredQuietPeriods))

		diags := m.ExecuteAllPendingOperations(ctx)
		allDiags.Append(diags...)
		if diags.HasError() {
			return allDiags
		}

		// Wait to let concurrent resources queue more operations
		time.Sleep(BatchCollectionWindow)

		if m.hasPendingOperations() {
			consecutiveQuiet = 0
			continue
		}

		consecutiveQuiet++
		if consecutiveQuiet >= requiredQuietPeriods {
			tflog.Info(ctx, fmt.Sprintf("[BULK-OPS] Stage barrier: all operations flushed (%d quiet periods confirmed) — proceeding to next stage", requiredQuietPeriods))
			break
		}
	}

	return allDiags
}

func (m *Manager) refreshAllResourceCaches(ctx context.Context) {
	if m.clearCacheFunc == nil || m.contextProvider == nil {
		return
	}
	for _, cacheKey := range finalCacheRefreshKeys {
		m.clearCacheFunc(ctx, m.contextProvider(), cacheKey)
	}
}

func (m *Manager) ExecuteAllPendingOperations(ctx context.Context) diag.Diagnostics {
	// Ensure only one execution runs at a time - prevents race conditions when multiple
	// timer callbacks fire due to low parallelism causing resource waves
	m.executionMutex.Lock()
	defer m.executionMutex.Unlock()

	var diagnostics diag.Diagnostics
	anyOperationsPerformed := false

	// After executing all types in order, check if more operations arrived
	// during execution (due to low parallelism releasing Terraform goroutines in waves).
	const maxPasses = 100
	for pass := 0; pass < maxPasses; pass++ {
		tflog.Info(ctx, fmt.Sprintf("[BULK-OPS] ExecuteAllPendingOperations pass %d — checking for queued operations", pass))

		if time.Since(m.lastOperationTime) < BatchCollectionWindow {
			remaining := BatchCollectionWindow - time.Since(m.lastOperationTime)
			tflog.Debug(ctx, fmt.Sprintf("Pass %d: Waiting %v to collect more operations before executing", pass, remaining))
			time.Sleep(remaining)
		}

		var opsDiags diag.Diagnostics
		var operationsPerformed bool

		switch m.mode {
		case "datacenter":
			tflog.Debug(ctx, fmt.Sprintf("Pass %d: Executing pending operations in 'datacenter' mode", pass))
			opsDiags, operationsPerformed = m.ExecuteDatacenterOperations(ctx)
		case "campus":
			tflog.Debug(ctx, fmt.Sprintf("Pass %d: Executing pending operations in 'campus' mode", pass))
			opsDiags, operationsPerformed = m.ExecuteCampusOperations(ctx)
		default:
			tflog.Warn(ctx, fmt.Sprintf("Unknown mode '%s', defaulting to 'datacenter' mode", m.mode))
			opsDiags, operationsPerformed = m.ExecuteDatacenterOperations(ctx)
		}

		diagnostics.Append(opsDiags...)

		if operationsPerformed {
			anyOperationsPerformed = true
		}

		if opsDiags.HasError() {
			tflog.Debug(ctx, fmt.Sprintf("Pass %d: Stopping due to errors", pass))
			break
		}

		if !operationsPerformed {
			tflog.Debug(ctx, fmt.Sprintf("Pass %d: No operations performed, done", pass))
			break
		}

		// Wait briefly to let in-flight Terraform goroutines submit new operations
		time.Sleep(BatchCollectionWindow)

		// Check if more operations arrived during execution
		if !m.hasPendingOperations() {
			tflog.Debug(ctx, fmt.Sprintf("Pass %d: No more pending operations, done", pass))
			break
		}

		tflog.Debug(ctx, fmt.Sprintf("Pass %d: More operations arrived during execution, starting next pass", pass))
	}

	if anyOperationsPerformed {
		waitDuration := 800 * time.Millisecond
		tflog.Debug(ctx, fmt.Sprintf("Waiting %v for all operations to propagate before final cache refresh", waitDuration))
		time.Sleep(waitDuration)

		tflog.Debug(ctx, "Final cache clear after all operations to ensure verification with fresh data")
		m.refreshAllResourceCaches(ctx)
	}

	return diagnostics
}

func (s orderedExecutionState) execute(opType string, count int, execFunc func(context.Context) diag.Diagnostics, resourceName string) bool {
	if count == 0 {
		return true
	}

	tflog.Info(s.ctx, fmt.Sprintf("[BULK-OPS] >>> Proceeding with %s %s for %d resource(s) — sending API request...", resourceName, opType, count))
	diags := execFunc(s.ctx)
	s.diagnostics.Append(diags...)
	if diags.HasError() {
		tflog.Error(s.ctx, fmt.Sprintf("[BULK-OPS] <<< FAILED %s %s — aborting remaining operations", resourceName, opType))
		s.manager.FailAllPendingOperations(s.ctx, fmt.Errorf("bulk %s %s operation failed", resourceName, opType))
		return false
	}

	tflog.Info(s.ctx, fmt.Sprintf("[BULK-OPS] <<< Completed %s %s for %d resource(s) — success, moving to next type", resourceName, opType, count))
	*s.operationsPerformed = true
	return true
}

func (s orderedExecutionState) executeSequence(opType string, resourceTypes []string) bool {
	for _, resourceType := range resourceTypes {
		if !s.execute(opType, s.manager.getOperationCount(resourceType, opType), func(ctx context.Context) diag.Diagnostics {
			return s.manager.ExecuteBulk(ctx, resourceType, opType)
		}, resourceType) {
			return false
		}
	}
	return true
}

func joinOperationOrders(orders ...[]string) []string {
	length := 0
	for _, order := range orders {
		length += len(order)
	}

	joined := make([]string, 0, length)
	for _, order := range orders {
		joined = append(joined, order...)
	}
	return joined
}

func reverseOperationOrder(order []string) []string {
	reversed := make([]string, len(order))
	for i, resourceType := range order {
		reversed[len(order)-1-i] = resourceType
	}
	return reversed
}

func (m *Manager) ExecuteDatacenterOperations(ctx context.Context) (diag.Diagnostics, bool) {
	var diagnostics diag.Diagnostics
	operationsPerformed := false
	execution := orderedExecutionState{m, ctx, &diagnostics, &operationsPerformed}
	// PUT operations - DC Order
	// sfp_breakout is intentionally omitted because it only supports PATCH.
	if !execution.executeSequence("PUT", datacenterPutOrder) {
		return diagnostics, operationsPerformed
	}

	// PATCH operations - DC Order
	if !execution.executeSequence("PATCH", datacenterPatchOrder) {
		return diagnostics, operationsPerformed
	}

	// DELETE operations - reverse normal Data Center staging order.
	// sfp_breakout is intentionally omitted because it only supports PATCH.
	if !execution.executeSequence("DELETE", datacenterDeleteOrder) {
		return diagnostics, operationsPerformed
	}

	return diagnostics, operationsPerformed
}

func (m *Manager) ExecuteCampusOperations(ctx context.Context) (diag.Diagnostics, bool) {
	var diagnostics diag.Diagnostics
	operationsPerformed := false
	execution := orderedExecutionState{m, ctx, &diagnostics, &operationsPerformed}

	// PUT operations - Campus Order
	if !execution.executeSequence("PUT", campusPutOrder) {
		return diagnostics, operationsPerformed
	}

	// PATCH operations - Campus Order
	if !execution.executeSequence("PATCH", campusPatchOrder) {
		return diagnostics, operationsPerformed
	}

	// DELETE operations - Reverse Campus Order
	if !execution.executeSequence("DELETE", campusDeleteOrder) {
		return diagnostics, operationsPerformed
	}

	return diagnostics, operationsPerformed
}

func (m *Manager) ShouldExecuteOperations(ctx context.Context) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// If there are no pending operations, no need to execute
	if !m.hasPendingOperationsLocked() {
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

func (m *Manager) pendingOperationSummaryLocked() (int, map[string]interface{}) {
	totalCount := 0
	summary := make(map[string]interface{})

	for resourceType, operations := range m.resources {
		putCount := len(operations.Put)
		patchCount := len(operations.Patch)
		deleteCount := len(operations.Delete)
		totalCount += putCount + patchCount + deleteCount

		if putCount > 0 {
			summary[resourceType+"_put_count"] = putCount
		}
		if patchCount > 0 {
			summary[resourceType+"_patch_count"] = patchCount
		}
		if deleteCount > 0 {
			summary[resourceType+"_delete_count"] = deleteCount
		}
	}

	summary["total_count"] = totalCount
	return totalCount, summary
}

func (m *Manager) ExecuteIfMultipleOperations(ctx context.Context) diag.Diagnostics {
	m.mutex.Lock()
	totalCount, summary := m.pendingOperationSummaryLocked()
	m.mutex.Unlock()

	if totalCount == 0 {
		return nil
	}

	tflog.Debug(ctx, "Multiple operations detected, executing in sequence", summary)
	return m.ExecuteAllPendingOperations(ctx)
}
