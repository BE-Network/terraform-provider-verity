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
	var diagnostics diag.Diagnostics

	if time.Since(m.lastOperationTime) < BatchCollectionWindow {
		remaining := BatchCollectionWindow - time.Since(m.lastOperationTime)
		tflog.Debug(ctx, fmt.Sprintf("Waiting %v to collect more operations before executing", remaining))
		time.Sleep(remaining)
	}

	var opsDiags diag.Diagnostics
	var operationsPerformed bool

	switch m.mode {
	case "datacenter":
		tflog.Debug(ctx, "Executing pending operations in 'datacenter' mode")
		opsDiags, operationsPerformed = m.ExecuteDatacenterOperations(ctx)
	case "campus":
		tflog.Debug(ctx, "Executing pending operations in 'campus' mode")
		opsDiags, operationsPerformed = m.ExecuteCampusOperations(ctx)
	default:
		tflog.Warn(ctx, fmt.Sprintf("Unknown mode '%s', defaulting to 'datacenter' mode", m.mode))
		opsDiags, operationsPerformed = m.ExecuteDatacenterOperations(ctx)
	}

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
			m.clearCacheFunc(ctx, m.contextProvider(), "packet_queues")
			m.clearCacheFunc(ctx, m.contextProvider(), "eth_port_profiles")
			m.clearCacheFunc(ctx, m.contextProvider(), "eth_port_settings")
			m.clearCacheFunc(ctx, m.contextProvider(), "lags")
			m.clearCacheFunc(ctx, m.contextProvider(), "sflow_collectors")
			m.clearCacheFunc(ctx, m.contextProvider(), "diagnostics_profiles")
			m.clearCacheFunc(ctx, m.contextProvider(), "diagnostics_port_profiles")
			m.clearCacheFunc(ctx, m.contextProvider(), "bundles")
			m.clearCacheFunc(ctx, m.contextProvider(), "acls_ipv4")
			m.clearCacheFunc(ctx, m.contextProvider(), "acls_ipv6")
			m.clearCacheFunc(ctx, m.contextProvider(), "packet_brokers")
			m.clearCacheFunc(ctx, m.contextProvider(), "badges")
			m.clearCacheFunc(ctx, m.contextProvider(), "switchpoints")
			m.clearCacheFunc(ctx, m.contextProvider(), "device_controllers")
			m.clearCacheFunc(ctx, m.contextProvider(), "authenticated_eth_ports")
			m.clearCacheFunc(ctx, m.contextProvider(), "device_voice_settings")
			m.clearCacheFunc(ctx, m.contextProvider(), "service_port_profiles")
			m.clearCacheFunc(ctx, m.contextProvider(), "voice_port_profiles")
			m.clearCacheFunc(ctx, m.contextProvider(), "as_path_access_lists")
			m.clearCacheFunc(ctx, m.contextProvider(), "community_lists")
			m.clearCacheFunc(ctx, m.contextProvider(), "device_settings")
			m.clearCacheFunc(ctx, m.contextProvider(), "extended_community_lists")
			m.clearCacheFunc(ctx, m.contextProvider(), "ipv4_lists")
			m.clearCacheFunc(ctx, m.contextProvider(), "ipv4_prefix_lists")
			m.clearCacheFunc(ctx, m.contextProvider(), "ipv6_lists")
			m.clearCacheFunc(ctx, m.contextProvider(), "ipv6_prefix_lists")
			m.clearCacheFunc(ctx, m.contextProvider(), "route_map_clauses")
			m.clearCacheFunc(ctx, m.contextProvider(), "route_maps")
			m.clearCacheFunc(ctx, m.contextProvider(), "sfp_breakouts")
			m.clearCacheFunc(ctx, m.contextProvider(), "sites")
			m.clearCacheFunc(ctx, m.contextProvider(), "pods")
			m.clearCacheFunc(ctx, m.contextProvider(), "port_acls")
			m.clearCacheFunc(ctx, m.contextProvider(), "pb_routing")
			m.clearCacheFunc(ctx, m.contextProvider(), "pb_routing_acl")
			m.clearCacheFunc(ctx, m.contextProvider(), "spine_planes")
			m.clearCacheFunc(ctx, m.contextProvider(), "grouping_rules")
			m.clearCacheFunc(ctx, m.contextProvider(), "threshold_groups")
			m.clearCacheFunc(ctx, m.contextProvider(), "thresholds")
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

	// PUT operations - DC Order (note: sfp_breakout and site are skipped - they only support GET and PATCH)
	// 2. ipv6_prefix_list
	if !execute("PUT", m.getOperationCount("ipv6_prefix_list", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "ipv6_prefix_list", "PUT") }, "IPv6 Prefix List") {
		return diagnostics, operationsPerformed
	}
	// 3. community_list
	if !execute("PUT", m.getOperationCount("community_list", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "community_list", "PUT") }, "Community List") {
		return diagnostics, operationsPerformed
	}
	// 4. ipv4_prefix_list
	if !execute("PUT", m.getOperationCount("ipv4_prefix_list", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "ipv4_prefix_list", "PUT") }, "IPv4 Prefix List") {
		return diagnostics, operationsPerformed
	}
	// 5. extended_community_list
	if !execute("PUT", m.getOperationCount("extended_community_list", "PUT"), func(ctx context.Context) diag.Diagnostics {
		return m.ExecuteBulk(ctx, "extended_community_list", "PUT")
	}, "Extended Community List") {
		return diagnostics, operationsPerformed
	}
	// 6. as_path_access_list
	if !execute("PUT", m.getOperationCount("as_path_access_list", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "as_path_access_list", "PUT") }, "AS Path Access List") {
		return diagnostics, operationsPerformed
	}

	// BEFORE executing route_map_clause PUT, check for circular reference scenario
	circularPutInfo := m.detectCircularReferenceScenario(ctx)

	if circularPutInfo.NeedsFix {
		tflog.Info(ctx, "Applying circular reference fix for route_map_clause and tenant", map[string]interface{}{
			"affected_clauses":   circularPutInfo.ClauseNames,
			"referenced_tenants": circularPutInfo.TenantNames,
		})

		// Temporarily replace affected route_map_clause PUT data with versions having empty match_vrf
		m.mutex.Lock()
		routeMapClauseOps := m.resources["route_map_clause"]
		for name, data := range circularPutInfo.AffectedClauses {
			modifiedData := m.createRouteMapClauseWithEmptyMatchVrf(data)
			routeMapClauseOps.Put[name] = modifiedData
			tflog.Debug(ctx, fmt.Sprintf("Temporarily setting match_vrf to empty for route_map_clause: %s", name))
		}
		m.mutex.Unlock()
	}

	// 7. route_map_clause (PUT with empty match_vrf if circular ref detected)
	if !execute("PUT", m.getOperationCount("route_map_clause", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "route_map_clause", "PUT") }, "Route Map Clause") {
		return diagnostics, operationsPerformed
	}
	// 8-9. acl (both ipv6 and ipv4)
	if !execute("PUT", m.getOperationCount("acl", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "acl", "PUT") }, "ACL") {
		return diagnostics, operationsPerformed
	}
	// 10. route_map
	if !execute("PUT", m.getOperationCount("route_map", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "route_map", "PUT") }, "Route Map") {
		return diagnostics, operationsPerformed
	}
	// 11. pb_routing_acl
	if !execute("PUT", m.getOperationCount("pb_routing_acl", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "pb_routing_acl", "PUT") }, "PB Routing ACL") {
		return diagnostics, operationsPerformed
	}
	// 12. tenant
	if !execute("PUT", m.getOperationCount("tenant", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "tenant", "PUT") }, "Tenant") {
		return diagnostics, operationsPerformed
	}

	// If circular reference fix was applied, now PATCH route_map_clause with match_vrf
	if circularPutInfo.NeedsFix && len(circularPutInfo.AffectedClauses) > 0 {
		tflog.Info(ctx, "Applying PATCH to restore match_vrf fields in route_map_clause")

		m.mutex.Lock()
		// Add PATCH operations for the affected route_map_clauses
		routeMapClauseOps := m.resources["route_map_clause"]
		if routeMapClauseOps.Patch == nil {
			routeMapClauseOps.Patch = make(map[string]interface{})
		}

		affectedNames := make([]string, 0, len(circularPutInfo.AffectedClauses))
		for name, originalData := range circularPutInfo.AffectedClauses {
			// Extract the original match_vrf value to restore
			matchVrfValue := m.getMatchVrfValue(originalData)

			// Create PATCH data to restore the match_vrf
			patchData := m.createMatchVrfPatchData(originalData, matchVrfValue)
			if patchData != nil {
				routeMapClauseOps.Patch[name] = patchData
				affectedNames = append(affectedNames, name)
				tflog.Debug(ctx, fmt.Sprintf("Prepared PATCH for route_map_clause: %s", name))
			}

			// Restore original PUT data for future reference
			routeMapClauseOps.Put[name] = originalData
		}
		m.mutex.Unlock()

		tflog.Info(ctx, "Executing PATCH to restore match_vrf", map[string]interface{}{
			"affected_resources": affectedNames,
		})
		if !execute("PATCH", m.getOperationCount("route_map_clause", "PATCH"), func(ctx context.Context) diag.Diagnostics {
			return m.ExecuteBulk(ctx, "route_map_clause", "PATCH")
		}, "Route Map Clause (match_vrf restore)") {
			return diagnostics, operationsPerformed
		}

		// Clean up the PATCH operations
		m.mutex.Lock()
		for name := range circularPutInfo.AffectedClauses {
			delete(routeMapClauseOps.Patch, name)
		}
		m.mutex.Unlock()

		tflog.Info(ctx, "Successfully restored match_vrf fields via PATCH")
	}

	// 13. pb_routing
	if !execute("PUT", m.getOperationCount("pb_routing", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "pb_routing", "PUT") }, "PB Routing") {
		return diagnostics, operationsPerformed
	}
	// 14. ipv4_list
	if !execute("PUT", m.getOperationCount("ipv4_list", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "ipv4_list", "PUT") }, "IPv4 List") {
		return diagnostics, operationsPerformed
	}
	// 15. ipv6_list
	if !execute("PUT", m.getOperationCount("ipv6_list", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "ipv6_list", "PUT") }, "IPv6 List") {
		return diagnostics, operationsPerformed
	}
	// 16. service
	if !execute("PUT", m.getOperationCount("service", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "service", "PUT") }, "Service") {
		return diagnostics, operationsPerformed
	}
	// 17. port_acl
	if !execute("PUT", m.getOperationCount("port_acl", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "port_acl", "PUT") }, "Port ACL") {
		return diagnostics, operationsPerformed
	}
	// 18. packet_broker
	if !execute("PUT", m.getOperationCount("packet_broker", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "packet_broker", "PUT") }, "Packet Broker") {
		return diagnostics, operationsPerformed
	}
	// 19. eth_port_profile
	if !execute("PUT", m.getOperationCount("eth_port_profile", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "eth_port_profile", "PUT") }, "Eth Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 20. packet_queue
	if !execute("PUT", m.getOperationCount("packet_queue", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "packet_queue", "PUT") }, "Packet Queue") {
		return diagnostics, operationsPerformed
	}
	// 21. sflow_collector
	if !execute("PUT", m.getOperationCount("sflow_collector", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "sflow_collector", "PUT") }, "SFlow Collector") {
		return diagnostics, operationsPerformed
	}
	// 22. gateway
	if !execute("PUT", m.getOperationCount("gateway", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "gateway", "PUT") }, "Gateway") {
		return diagnostics, operationsPerformed
	}
	// 23. lag
	if !execute("PUT", m.getOperationCount("lag", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "lag", "PUT") }, "LAG") {
		return diagnostics, operationsPerformed
	}
	// 24. eth_port_settings
	if !execute("PUT", m.getOperationCount("eth_port_settings", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "eth_port_settings", "PUT") }, "Eth Port Settings") {
		return diagnostics, operationsPerformed
	}
	// 25. diagnostics_profile
	if !execute("PUT", m.getOperationCount("diagnostics_profile", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "diagnostics_profile", "PUT") }, "Diagnostics Profile") {
		return diagnostics, operationsPerformed
	}
	// 26. gateway_profile
	if !execute("PUT", m.getOperationCount("gateway_profile", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "gateway_profile", "PUT") }, "Gateway Profile") {
		return diagnostics, operationsPerformed
	}
	// 27. diagnostics_port_profile
	if !execute("PUT", m.getOperationCount("diagnostics_port_profile", "PUT"), func(ctx context.Context) diag.Diagnostics {
		return m.ExecuteBulk(ctx, "diagnostics_port_profile", "PUT")
	}, "Diagnostics Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 28. bundle
	if !execute("PUT", m.getOperationCount("bundle", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "bundle", "PUT") }, "Bundle") {
		return diagnostics, operationsPerformed
	}
	// 29. pod
	if !execute("PUT", m.getOperationCount("pod", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "pod", "PUT") }, "Pod") {
		return diagnostics, operationsPerformed
	}
	// 30. badge
	if !execute("PUT", m.getOperationCount("badge", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "badge", "PUT") }, "Badge") {
		return diagnostics, operationsPerformed
	}
	// 31. spine_plane
	if !execute("PUT", m.getOperationCount("spine_plane", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "spine_plane", "PUT") }, "Spine Plane") {
		return diagnostics, operationsPerformed
	}
	// 32. switchpoint
	if !execute("PUT", m.getOperationCount("switchpoint", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "switchpoint", "PUT") }, "Switchpoint") {
		return diagnostics, operationsPerformed
	}
	// 33. device_settings
	if !execute("PUT", m.getOperationCount("device_settings", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "device_settings", "PUT") }, "Device Settings") {
		return diagnostics, operationsPerformed
	}
	// 34. threshold
	if !execute("PUT", m.getOperationCount("threshold", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "threshold", "PUT") }, "Threshold") {
		return diagnostics, operationsPerformed
	}
	// 35. grouping_rule
	if !execute("PUT", m.getOperationCount("grouping_rule", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "grouping_rule", "PUT") }, "Grouping Rule") {
		return diagnostics, operationsPerformed
	}
	// 36. threshold_group
	if !execute("PUT", m.getOperationCount("threshold_group", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "threshold_group", "PUT") }, "Threshold Group") {
		return diagnostics, operationsPerformed
	}
	// 38. device_controller
	if !execute("PUT", m.getOperationCount("device_controller", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "device_controller", "PUT") }, "Device Controller") {
		return diagnostics, operationsPerformed
	}

	// PATCH operations - DC Order
	// 1. sfp_breakout
	if !execute("PATCH", m.getOperationCount("sfp_breakout", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "sfp_breakout", "PATCH") }, "SFP Breakout") {
		return diagnostics, operationsPerformed
	}
	// 2. ipv6_prefix_list
	if !execute("PATCH", m.getOperationCount("ipv6_prefix_list", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "ipv6_prefix_list", "PATCH") }, "IPv6 Prefix List") {
		return diagnostics, operationsPerformed
	}
	// 3. community_list
	if !execute("PATCH", m.getOperationCount("community_list", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "community_list", "PATCH") }, "Community List") {
		return diagnostics, operationsPerformed
	}
	// 4. ipv4_prefix_list
	if !execute("PATCH", m.getOperationCount("ipv4_prefix_list", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "ipv4_prefix_list", "PATCH") }, "IPv4 Prefix List") {
		return diagnostics, operationsPerformed
	}
	// 5. extended_community_list
	if !execute("PATCH", m.getOperationCount("extended_community_list", "PATCH"), func(ctx context.Context) diag.Diagnostics {
		return m.ExecuteBulk(ctx, "extended_community_list", "PATCH")
	}, "Extended Community List") {
		return diagnostics, operationsPerformed
	}
	// 6. as_path_access_list
	if !execute("PATCH", m.getOperationCount("as_path_access_list", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "as_path_access_list", "PATCH") }, "AS Path Access List") {
		return diagnostics, operationsPerformed
	}
	// 7. route_map_clause
	if !execute("PATCH", m.getOperationCount("route_map_clause", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "route_map_clause", "PATCH") }, "Route Map Clause") {
		return diagnostics, operationsPerformed
	}
	// 8-9. acl (both ipv6 and ipv4)
	if !execute("PATCH", m.getOperationCount("acl", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "acl", "PATCH") }, "ACL") {
		return diagnostics, operationsPerformed
	}
	// 10. route_map
	if !execute("PATCH", m.getOperationCount("route_map", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "route_map", "PATCH") }, "Route Map") {
		return diagnostics, operationsPerformed
	}
	// 11. pb_routing_acl
	if !execute("PATCH", m.getOperationCount("pb_routing_acl", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "pb_routing_acl", "PATCH") }, "PB Routing ACL") {
		return diagnostics, operationsPerformed
	}
	// 12. tenant
	if !execute("PATCH", m.getOperationCount("tenant", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "tenant", "PATCH") }, "Tenant") {
		return diagnostics, operationsPerformed
	}
	// 13. pb_routing
	if !execute("PATCH", m.getOperationCount("pb_routing", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "pb_routing", "PATCH") }, "PB Routing") {
		return diagnostics, operationsPerformed
	}
	// 14. ipv4_list
	if !execute("PATCH", m.getOperationCount("ipv4_list", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "ipv4_list", "PATCH") }, "IPv4 List") {
		return diagnostics, operationsPerformed
	}
	// 15. ipv6_list
	if !execute("PATCH", m.getOperationCount("ipv6_list", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "ipv6_list", "PATCH") }, "IPv6 List") {
		return diagnostics, operationsPerformed
	}
	// 16. service
	if !execute("PATCH", m.getOperationCount("service", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "service", "PATCH") }, "Service") {
		return diagnostics, operationsPerformed
	}
	// 17. port_acl
	if !execute("PATCH", m.getOperationCount("port_acl", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "port_acl", "PATCH") }, "Port ACL") {
		return diagnostics, operationsPerformed
	}
	// 18. packet_broker
	if !execute("PATCH", m.getOperationCount("packet_broker", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "packet_broker", "PATCH") }, "Packet Broker") {
		return diagnostics, operationsPerformed
	}
	// 19. eth_port_profile
	if !execute("PATCH", m.getOperationCount("eth_port_profile", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "eth_port_profile", "PATCH") }, "Eth Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 20. packet_queue
	if !execute("PATCH", m.getOperationCount("packet_queue", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "packet_queue", "PATCH") }, "Packet Queue") {
		return diagnostics, operationsPerformed
	}
	// 21. sflow_collector
	if !execute("PATCH", m.getOperationCount("sflow_collector", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "sflow_collector", "PATCH") }, "SFlow Collector") {
		return diagnostics, operationsPerformed
	}
	// 22. gateway
	if !execute("PATCH", m.getOperationCount("gateway", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "gateway", "PATCH") }, "Gateway") {
		return diagnostics, operationsPerformed
	}
	// 23. lag
	if !execute("PATCH", m.getOperationCount("lag", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "lag", "PATCH") }, "LAG") {
		return diagnostics, operationsPerformed
	}
	// 24. eth_port_settings
	if !execute("PATCH", m.getOperationCount("eth_port_settings", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "eth_port_settings", "PATCH") }, "Eth Port Settings") {
		return diagnostics, operationsPerformed
	}
	// 25. diagnostics_profile
	if !execute("PATCH", m.getOperationCount("diagnostics_profile", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "diagnostics_profile", "PATCH") }, "Diagnostics Profile") {
		return diagnostics, operationsPerformed
	}
	// 26. gateway_profile
	if !execute("PATCH", m.getOperationCount("gateway_profile", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "gateway_profile", "PATCH") }, "Gateway Profile") {
		return diagnostics, operationsPerformed
	}
	// 27. diagnostics_port_profile
	if !execute("PATCH", m.getOperationCount("diagnostics_port_profile", "PATCH"), func(ctx context.Context) diag.Diagnostics {
		return m.ExecuteBulk(ctx, "diagnostics_port_profile", "PATCH")
	}, "Diagnostics Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 28. bundle
	if !execute("PATCH", m.getOperationCount("bundle", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "bundle", "PATCH") }, "Bundle") {
		return diagnostics, operationsPerformed
	}
	// 29. pod
	if !execute("PATCH", m.getOperationCount("pod", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "pod", "PATCH") }, "Pod") {
		return diagnostics, operationsPerformed
	}
	// 30. badge
	if !execute("PATCH", m.getOperationCount("badge", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "badge", "PATCH") }, "Badge") {
		return diagnostics, operationsPerformed
	}
	// 31. spine_plane
	if !execute("PATCH", m.getOperationCount("spine_plane", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "spine_plane", "PATCH") }, "Spine Plane") {
		return diagnostics, operationsPerformed
	}
	// 32. switchpoint
	if !execute("PATCH", m.getOperationCount("switchpoint", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "switchpoint", "PATCH") }, "Switchpoint") {
		return diagnostics, operationsPerformed
	}
	// 33. device_settings
	if !execute("PATCH", m.getOperationCount("device_settings", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "device_settings", "PATCH") }, "Device Settings") {
		return diagnostics, operationsPerformed
	}
	// 34. threshold
	if !execute("PATCH", m.getOperationCount("threshold", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "threshold", "PATCH") }, "Threshold") {
		return diagnostics, operationsPerformed
	}
	// 35. grouping_rule
	if !execute("PATCH", m.getOperationCount("grouping_rule", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "grouping_rule", "PATCH") }, "Grouping Rule") {
		return diagnostics, operationsPerformed
	}
	// 36. threshold_group
	if !execute("PATCH", m.getOperationCount("threshold_group", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "threshold_group", "PATCH") }, "Threshold Group") {
		return diagnostics, operationsPerformed
	}
	// 37. site
	if !execute("PATCH", m.getOperationCount("site", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "site", "PATCH") }, "Site") {
		return diagnostics, operationsPerformed
	}
	// 38. device_controller
	if !execute("PATCH", m.getOperationCount("device_controller", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "device_controller", "PATCH") }, "Device Controller") {
		return diagnostics, operationsPerformed
	}

	// DELETE operations - Reverse DC Order (note: sfp_breakout and site are skipped - they only support GET and PATCH)

	// CIRCULAR REFERENCE FIX FOR DELETE OPERATIONS
	// Before starting DELETE operations, check if we need to handle circular references
	// between route_map_clause (match_vrf) and tenant (export/import_route_map)
	circularInfo := m.detectCompleteCircularDeleteSet(ctx)

	if circularInfo.IsCompleteSet {
		tflog.Info(ctx, "Complete circular reference set detected in DELETE operations - applying fix")
		tflog.Info(ctx, "Clearing match_vrf references before deletion", map[string]interface{}{
			"affected_clauses": circularInfo.ClauseNames,
			"affected_tenants": circularInfo.TenantNames,
		})

		m.mutex.Lock()
		// Add PATCH operations to clear match_vrf from route_map_clauses being deleted
		routeMapClauseOps := m.resources["route_map_clause"]
		if routeMapClauseOps.Patch == nil {
			routeMapClauseOps.Patch = make(map[string]interface{})
		}

		// For each clause in the circular set, create PATCH to clear match_vrf
		for _, clauseName := range circularInfo.ClauseNames {
			// Get the clause data from circularInfo.AffectedClauses (fetched from API)
			var originalData interface{}
			if clauseData, exists := circularInfo.AffectedClauses[clauseName]; exists {
				originalData = clauseData
			} else if putData, exists := routeMapClauseOps.Put[clauseName]; exists {
				originalData = putData
			}

			// Create PATCH data with empty match_vrf
			patchData := m.createMatchVrfPatchData(originalData, "")
			routeMapClauseOps.Patch[clauseName] = patchData
			tflog.Debug(ctx, fmt.Sprintf("Clearing match_vrf for route_map_clause: %s before deletion", clauseName))
		}
		m.mutex.Unlock()

		tflog.Info(ctx, "Executing PATCH to clear match_vrf references")
		if !execute("PATCH", m.getOperationCount("route_map_clause", "PATCH"),
			func(ctx context.Context) diag.Diagnostics {
				return m.ExecuteBulk(ctx, "route_map_clause", "PATCH")
			}, "Route Map Clause (clear match_vrf before deletion)") {
			return diagnostics, operationsPerformed
		}

		// Clean up PATCH operations
		m.mutex.Lock()
		for _, clauseName := range circularInfo.ClauseNames {
			delete(routeMapClauseOps.Patch, clauseName)
		}
		m.mutex.Unlock()

		tflog.Info(ctx, "Successfully cleared match_vrf references, proceeding with deletions")
	}

	// 38. device_controller
	if !execute("DELETE", m.getOperationCount("device_controller", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "device_controller", "DELETE") }, "Device Controller") {
		return diagnostics, operationsPerformed
	}
	// 36. threshold_group
	if !execute("DELETE", m.getOperationCount("threshold_group", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "threshold_group", "DELETE") }, "Threshold Group") {
		return diagnostics, operationsPerformed
	}
	// 35. grouping_rule
	if !execute("DELETE", m.getOperationCount("grouping_rule", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "grouping_rule", "DELETE") }, "Grouping Rule") {
		return diagnostics, operationsPerformed
	}
	// 34. threshold
	if !execute("DELETE", m.getOperationCount("threshold", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "threshold", "DELETE") }, "Threshold") {
		return diagnostics, operationsPerformed
	}
	// 33. device_settings
	if !execute("DELETE", m.getOperationCount("device_settings", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "device_settings", "DELETE") }, "Device Settings") {
		return diagnostics, operationsPerformed
	}
	// 32. switchpoint
	if !execute("DELETE", m.getOperationCount("switchpoint", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "switchpoint", "DELETE") }, "Switchpoint") {
		return diagnostics, operationsPerformed
	}
	// 31. spine_plane
	if !execute("DELETE", m.getOperationCount("spine_plane", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "spine_plane", "DELETE") }, "Spine Plane") {
		return diagnostics, operationsPerformed
	}
	// 30. badge
	if !execute("DELETE", m.getOperationCount("badge", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "badge", "DELETE") }, "Badge") {
		return diagnostics, operationsPerformed
	}
	// 29. pod
	if !execute("DELETE", m.getOperationCount("pod", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "pod", "DELETE") }, "Pod") {
		return diagnostics, operationsPerformed
	}
	// 28. bundle
	if !execute("DELETE", m.getOperationCount("bundle", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "bundle", "DELETE") }, "Bundle") {
		return diagnostics, operationsPerformed
	}
	// 27. diagnostics_port_profile
	if !execute("DELETE", m.getOperationCount("diagnostics_port_profile", "DELETE"), func(ctx context.Context) diag.Diagnostics {
		return m.ExecuteBulk(ctx, "diagnostics_port_profile", "DELETE")
	}, "Diagnostics Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 26. gateway_profile
	if !execute("DELETE", m.getOperationCount("gateway_profile", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "gateway_profile", "DELETE") }, "Gateway Profile") {
		return diagnostics, operationsPerformed
	}
	// 25. diagnostics_profile
	if !execute("DELETE", m.getOperationCount("diagnostics_profile", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "diagnostics_profile", "DELETE") }, "Diagnostics Profile") {
		return diagnostics, operationsPerformed
	}
	// 24. eth_port_settings
	if !execute("DELETE", m.getOperationCount("eth_port_settings", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "eth_port_settings", "DELETE") }, "Eth Port Settings") {
		return diagnostics, operationsPerformed
	}
	// 23. lag
	if !execute("DELETE", m.getOperationCount("lag", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "lag", "DELETE") }, "LAG") {
		return diagnostics, operationsPerformed
	}
	// 22. gateway
	if !execute("DELETE", m.getOperationCount("gateway", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "gateway", "DELETE") }, "Gateway") {
		return diagnostics, operationsPerformed
	}
	// 21. sflow_collector
	if !execute("DELETE", m.getOperationCount("sflow_collector", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "sflow_collector", "DELETE") }, "SFlow Collector") {
		return diagnostics, operationsPerformed
	}
	// 20. packet_queue
	if !execute("DELETE", m.getOperationCount("packet_queue", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "packet_queue", "DELETE") }, "Packet Queue") {
		return diagnostics, operationsPerformed
	}
	// 19. eth_port_profile
	if !execute("DELETE", m.getOperationCount("eth_port_profile", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "eth_port_profile", "DELETE") }, "Eth Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 18. packet_broker
	if !execute("DELETE", m.getOperationCount("packet_broker", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "packet_broker", "DELETE") }, "Packet Broker") {
		return diagnostics, operationsPerformed
	}
	// 17. port_acl
	if !execute("DELETE", m.getOperationCount("port_acl", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "port_acl", "DELETE") }, "Port ACL") {
		return diagnostics, operationsPerformed
	}
	// 16. service
	if !execute("DELETE", m.getOperationCount("service", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "service", "DELETE") }, "Service") {
		return diagnostics, operationsPerformed
	}
	// 15. ipv6_list
	if !execute("DELETE", m.getOperationCount("ipv6_list", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "ipv6_list", "DELETE") }, "IPv6 List Filter") {
		return diagnostics, operationsPerformed
	}
	// 14. ipv4_list
	if !execute("DELETE", m.getOperationCount("ipv4_list", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "ipv4_list", "DELETE") }, "IPv4 List Filter") {
		return diagnostics, operationsPerformed
	}
	// 13. pb_routing
	if !execute("DELETE", m.getOperationCount("pb_routing", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "pb_routing", "DELETE") }, "PB Routing") {
		return diagnostics, operationsPerformed
	}
	// 12. tenant
	if !execute("DELETE", m.getOperationCount("tenant", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "tenant", "DELETE") }, "Tenant") {
		return diagnostics, operationsPerformed
	}
	// 11. pb_routing_acl
	if !execute("DELETE", m.getOperationCount("pb_routing_acl", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "pb_routing_acl", "DELETE") }, "PB Routing ACL") {
		return diagnostics, operationsPerformed
	}
	// 10. route_map
	if !execute("DELETE", m.getOperationCount("route_map", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "route_map", "DELETE") }, "Route Map") {
		return diagnostics, operationsPerformed
	}
	// 8-9. acl (both ipv4 and ipv6)
	if !execute("DELETE", m.getOperationCount("acl", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "acl", "DELETE") }, "ACL") {
		return diagnostics, operationsPerformed
	}
	// 7. route_map_clause
	if !execute("DELETE", m.getOperationCount("route_map_clause", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "route_map_clause", "DELETE") }, "Route Map Clause") {
		return diagnostics, operationsPerformed
	}
	// 6. as_path_access_list
	if !execute("DELETE", m.getOperationCount("as_path_access_list", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "as_path_access_list", "DELETE") }, "AS Path Access List") {
		return diagnostics, operationsPerformed
	}
	// 5. extended_community_list
	if !execute("DELETE", m.getOperationCount("extended_community_list", "DELETE"), func(ctx context.Context) diag.Diagnostics {
		return m.ExecuteBulk(ctx, "extended_community_list", "DELETE")
	}, "Extended Community List") {
		return diagnostics, operationsPerformed
	}
	// 4. ipv4_prefix_list
	if !execute("DELETE", m.getOperationCount("ipv4_prefix_list", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "ipv4_prefix_list", "DELETE") }, "IPv4 Prefix List") {
		return diagnostics, operationsPerformed
	}
	// 3. community_list
	if !execute("DELETE", m.getOperationCount("community_list", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "community_list", "DELETE") }, "Community List") {
		return diagnostics, operationsPerformed
	}
	// 2. ipv6_prefix_list
	if !execute("DELETE", m.getOperationCount("ipv6_prefix_list", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "ipv6_prefix_list", "DELETE") }, "IPv6 Prefix List") {
		return diagnostics, operationsPerformed
	}

	return diagnostics, operationsPerformed
}

func (m *Manager) ExecuteCampusOperations(ctx context.Context) (diag.Diagnostics, bool) {
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

	// PUT operations - Campus Order (note: site is skipped - it only supports GET and PATCH)
	// 1. ipv4_list
	if !execute("PUT", m.getOperationCount("ipv4_list", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "ipv4_list", "PUT") }, "IPv4 List") {
		return diagnostics, operationsPerformed
	}
	// 2. ipv6_list
	if !execute("PUT", m.getOperationCount("ipv6_list", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "ipv6_list", "PUT") }, "IPv6 List") {
		return diagnostics, operationsPerformed
	}
	// 3-4. acl (both ipv4 and ipv6)
	if !execute("PUT", m.getOperationCount("acl", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "acl", "PUT") }, "ACL") {
		return diagnostics, operationsPerformed
	}
	// 5. pb_routing_acl
	if !execute("PUT", m.getOperationCount("pb_routing_acl", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "pb_routing_acl", "PUT") }, "PB Routing ACL") {
		return diagnostics, operationsPerformed
	}
	// 6. pb_routing
	if !execute("PUT", m.getOperationCount("pb_routing", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "pb_routing", "PUT") }, "PB Routing") {
		return diagnostics, operationsPerformed
	}
	// 7. port_acl
	if !execute("PUT", m.getOperationCount("port_acl", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "port_acl", "PUT") }, "Port ACL") {
		return diagnostics, operationsPerformed
	}
	// 8. service
	if !execute("PUT", m.getOperationCount("service", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "service", "PUT") }, "Service") {
		return diagnostics, operationsPerformed
	}
	// 9. eth_port_profile
	if !execute("PUT", m.getOperationCount("eth_port_profile", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "eth_port_profile", "PUT") }, "Eth Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 10. sflow_collector
	if !execute("PUT", m.getOperationCount("sflow_collector", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "sflow_collector", "PUT") }, "SFlow Collector") {
		return diagnostics, operationsPerformed
	}
	// 11. packet_queue
	if !execute("PUT", m.getOperationCount("packet_queue", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "packet_queue", "PUT") }, "Packet Queue") {
		return diagnostics, operationsPerformed
	}
	// 12. service_port_profile
	if !execute("PUT", m.getOperationCount("service_port_profile", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "service_port_profile", "PUT") }, "Service Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 13. diagnostics_port_profile
	if !execute("PUT", m.getOperationCount("diagnostics_port_profile", "PUT"), func(ctx context.Context) diag.Diagnostics {
		return m.ExecuteBulk(ctx, "diagnostics_port_profile", "PUT")
	}, "Diagnostics Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 14. device_voice_settings
	if !execute("PUT", m.getOperationCount("device_voice_settings", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "device_voice_settings", "PUT") }, "Device Voice Settings") {
		return diagnostics, operationsPerformed
	}
	// 15. authenticated_eth_port
	if !execute("PUT", m.getOperationCount("authenticated_eth_port", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "authenticated_eth_port", "PUT") }, "Authenticated Eth-Port") {
		return diagnostics, operationsPerformed
	}
	// 16. diagnostics_profile
	if !execute("PUT", m.getOperationCount("diagnostics_profile", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "diagnostics_profile", "PUT") }, "Diagnostics Profile") {
		return diagnostics, operationsPerformed
	}
	// 17. eth_port_settings
	if !execute("PUT", m.getOperationCount("eth_port_settings", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "eth_port_settings", "PUT") }, "Eth Port Settings") {
		return diagnostics, operationsPerformed
	}
	// 18. voice_port_profile
	if !execute("PUT", m.getOperationCount("voice_port_profile", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "voice_port_profile", "PUT") }, "Voice-Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 19. device_settings
	if !execute("PUT", m.getOperationCount("device_settings", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "device_settings", "PUT") }, "Device Settings") {
		return diagnostics, operationsPerformed
	}
	// 20. lag
	if !execute("PUT", m.getOperationCount("lag", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "lag", "PUT") }, "LAG") {
		return diagnostics, operationsPerformed
	}
	// 21. bundle
	if !execute("PUT", m.getOperationCount("bundle", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "bundle", "PUT") }, "Bundle") {
		return diagnostics, operationsPerformed
	}
	// 22. badge
	if !execute("PUT", m.getOperationCount("badge", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "badge", "PUT") }, "Badge") {
		return diagnostics, operationsPerformed
	}
	// 23. switchpoint
	if !execute("PUT", m.getOperationCount("switchpoint", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "switchpoint", "PUT") }, "Switchpoint") {
		return diagnostics, operationsPerformed
	}
	// 24. threshold
	if !execute("PUT", m.getOperationCount("threshold", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "threshold", "PUT") }, "Threshold") {
		return diagnostics, operationsPerformed
	}
	// 25. grouping_rule
	if !execute("PUT", m.getOperationCount("grouping_rule", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "grouping_rule", "PUT") }, "Grouping Rule") {
		return diagnostics, operationsPerformed
	}
	// 26. threshold_group
	if !execute("PUT", m.getOperationCount("threshold_group", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "threshold_group", "PUT") }, "Threshold Group") {
		return diagnostics, operationsPerformed
	}
	// 28. device_controller
	if !execute("PUT", m.getOperationCount("device_controller", "PUT"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "device_controller", "PUT") }, "Device Controller") {
		return diagnostics, operationsPerformed
	}

	// PATCH operations - Campus Order
	// 1. ipv4_list
	if !execute("PATCH", m.getOperationCount("ipv4_list", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "ipv4_list", "PATCH") }, "IPv4 List") {
		return diagnostics, operationsPerformed
	}
	// 2. ipv6_list
	if !execute("PATCH", m.getOperationCount("ipv6_list", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "ipv6_list", "PATCH") }, "IPv6 List") {
		return diagnostics, operationsPerformed
	}
	// 3-4. acl (both ipv4 and ipv6)
	if !execute("PATCH", m.getOperationCount("acl", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "acl", "PATCH") }, "ACL") {
		return diagnostics, operationsPerformed
	}
	// 5. pb_routing_acl
	if !execute("PATCH", m.getOperationCount("pb_routing_acl", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "pb_routing_acl", "PATCH") }, "PB Routing ACL") {
		return diagnostics, operationsPerformed
	}
	// 6. pb_routing
	if !execute("PATCH", m.getOperationCount("pb_routing", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "pb_routing", "PATCH") }, "PB Routing") {
		return diagnostics, operationsPerformed
	}
	// 7. port_acl
	if !execute("PATCH", m.getOperationCount("port_acl", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "port_acl", "PATCH") }, "Port ACL") {
		return diagnostics, operationsPerformed
	}
	// 8. service
	if !execute("PATCH", m.getOperationCount("service", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "service", "PATCH") }, "Service") {
		return diagnostics, operationsPerformed
	}
	// 9. eth_port_profile
	if !execute("PATCH", m.getOperationCount("eth_port_profile", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "eth_port_profile", "PATCH") }, "Eth Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 10. sflow_collector
	if !execute("PATCH", m.getOperationCount("sflow_collector", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "sflow_collector", "PATCH") }, "SFlow Collector") {
		return diagnostics, operationsPerformed
	}
	// 11. packet_queue
	if !execute("PATCH", m.getOperationCount("packet_queue", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "packet_queue", "PATCH") }, "Packet Queue") {
		return diagnostics, operationsPerformed
	}
	// 12. service_port_profile
	if !execute("PATCH", m.getOperationCount("service_port_profile", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "service_port_profile", "PATCH") }, "Service Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 13. diagnostics_port_profile
	if !execute("PATCH", m.getOperationCount("diagnostics_port_profile", "PATCH"), func(ctx context.Context) diag.Diagnostics {
		return m.ExecuteBulk(ctx, "diagnostics_port_profile", "PATCH")
	}, "Diagnostics Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 14. device_voice_settings
	if !execute("PATCH", m.getOperationCount("device_voice_settings", "PATCH"), func(ctx context.Context) diag.Diagnostics {
		return m.ExecuteBulk(ctx, "device_voice_settings", "PATCH")
	}, "Device Voice Settings") {
		return diagnostics, operationsPerformed
	}
	// 15. authenticated_eth_port
	if !execute("PATCH", m.getOperationCount("authenticated_eth_port", "PATCH"), func(ctx context.Context) diag.Diagnostics {
		return m.ExecuteBulk(ctx, "authenticated_eth_port", "PATCH")
	}, "Authenticated Eth-Port") {
		return diagnostics, operationsPerformed
	}
	// 16. diagnostics_profile
	if !execute("PATCH", m.getOperationCount("diagnostics_profile", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "diagnostics_profile", "PATCH") }, "Diagnostics Profile") {
		return diagnostics, operationsPerformed
	}
	// 17. eth_port_settings
	if !execute("PATCH", m.getOperationCount("eth_port_settings", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "eth_port_settings", "PATCH") }, "Eth Port Settings") {
		return diagnostics, operationsPerformed
	}
	// 18. voice_port_profile
	if !execute("PATCH", m.getOperationCount("voice_port_profile", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "voice_port_profile", "PATCH") }, "Voice-Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 19. device_settings
	if !execute("PATCH", m.getOperationCount("device_settings", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "device_settings", "PATCH") }, "Device Settings") {
		return diagnostics, operationsPerformed
	}
	// 20. lag
	if !execute("PATCH", m.getOperationCount("lag", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "lag", "PATCH") }, "LAG") {
		return diagnostics, operationsPerformed
	}
	// 21. bundle
	if !execute("PATCH", m.getOperationCount("bundle", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "bundle", "PATCH") }, "Bundle") {
		return diagnostics, operationsPerformed
	}
	// 22. badge
	if !execute("PATCH", m.getOperationCount("badge", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "badge", "PATCH") }, "Badge") {
		return diagnostics, operationsPerformed
	}
	// 23. switchpoint
	if !execute("PATCH", m.getOperationCount("switchpoint", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "switchpoint", "PATCH") }, "Switchpoint") {
		return diagnostics, operationsPerformed
	}
	// 24. threshold
	if !execute("PATCH", m.getOperationCount("threshold", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "threshold", "PATCH") }, "Threshold") {
		return diagnostics, operationsPerformed
	}
	// 25. grouping_rule
	if !execute("PATCH", m.getOperationCount("grouping_rule", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "grouping_rule", "PATCH") }, "Grouping Rule") {
		return diagnostics, operationsPerformed
	}
	// 26. threshold_group
	if !execute("PATCH", m.getOperationCount("threshold_group", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "threshold_group", "PATCH") }, "Threshold Group") {
		return diagnostics, operationsPerformed
	}
	// 27. site
	if !execute("PATCH", m.getOperationCount("site", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "site", "PATCH") }, "Site") {
		return diagnostics, operationsPerformed
	}
	// 28. device_controller
	if !execute("PATCH", m.getOperationCount("device_controller", "PATCH"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "device_controller", "PATCH") }, "Device Controller") {
		return diagnostics, operationsPerformed
	}

	// DELETE operations - Reverse Campus Order (note: site is skipped - it only supports GET and PATCH)
	// 28. device_controller
	if !execute("DELETE", m.getOperationCount("device_controller", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "device_controller", "DELETE") }, "Device Controller") {
		return diagnostics, operationsPerformed
	}
	// 26. threshold_group
	if !execute("DELETE", m.getOperationCount("threshold_group", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "threshold_group", "DELETE") }, "Threshold Group") {
		return diagnostics, operationsPerformed
	}
	// 25. grouping_rule
	if !execute("DELETE", m.getOperationCount("grouping_rule", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "grouping_rule", "DELETE") }, "Grouping Rule") {
		return diagnostics, operationsPerformed
	}
	// 24. threshold
	if !execute("DELETE", m.getOperationCount("threshold", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "threshold", "DELETE") }, "Threshold") {
		return diagnostics, operationsPerformed
	}
	// 23. switchpoint
	if !execute("DELETE", m.getOperationCount("switchpoint", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "switchpoint", "DELETE") }, "Switchpoint") {
		return diagnostics, operationsPerformed
	}
	// 22. badge
	if !execute("DELETE", m.getOperationCount("badge", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "badge", "DELETE") }, "Badge") {
		return diagnostics, operationsPerformed
	}
	// 21. bundle
	if !execute("DELETE", m.getOperationCount("bundle", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "bundle", "DELETE") }, "Bundle") {
		return diagnostics, operationsPerformed
	}
	// 20. lag
	if !execute("DELETE", m.getOperationCount("lag", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "lag", "DELETE") }, "LAG") {
		return diagnostics, operationsPerformed
	}
	// 19. device_settings
	if !execute("DELETE", m.getOperationCount("device_settings", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "device_settings", "DELETE") }, "Device Settings") {
		return diagnostics, operationsPerformed
	}
	// 18. voice_port_profile
	if !execute("DELETE", m.getOperationCount("voice_port_profile", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "voice_port_profile", "DELETE") }, "Voice-Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 17. eth_port_settings
	if !execute("DELETE", m.getOperationCount("eth_port_settings", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "eth_port_settings", "DELETE") }, "Eth Port Settings") {
		return diagnostics, operationsPerformed
	}
	// 16. diagnostics_profile
	if !execute("DELETE", m.getOperationCount("diagnostics_profile", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "diagnostics_profile", "DELETE") }, "Diagnostics Profile") {
		return diagnostics, operationsPerformed
	}
	// 15. authenticated_eth_port
	if !execute("DELETE", m.getOperationCount("authenticated_eth_port", "DELETE"), func(ctx context.Context) diag.Diagnostics {
		return m.ExecuteBulk(ctx, "authenticated_eth_port", "DELETE")
	}, "Authenticated Eth-Port") {
		return diagnostics, operationsPerformed
	}
	// 14. device_voice_settings
	if !execute("DELETE", m.getOperationCount("device_voice_settings", "DELETE"), func(ctx context.Context) diag.Diagnostics {
		return m.ExecuteBulk(ctx, "device_voice_settings", "DELETE")
	}, "Device Voice Settings") {
		return diagnostics, operationsPerformed
	}
	// 13. diagnostics_port_profile
	if !execute("DELETE", m.getOperationCount("diagnostics_port_profile", "DELETE"), func(ctx context.Context) diag.Diagnostics {
		return m.ExecuteBulk(ctx, "diagnostics_port_profile", "DELETE")
	}, "Diagnostics Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 12. service_port_profile
	if !execute("DELETE", m.getOperationCount("service_port_profile", "DELETE"), func(ctx context.Context) diag.Diagnostics {
		return m.ExecuteBulk(ctx, "service_port_profile", "DELETE")
	}, "Service Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 11. packet_queue
	if !execute("DELETE", m.getOperationCount("packet_queue", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "packet_queue", "DELETE") }, "Packet Queue") {
		return diagnostics, operationsPerformed
	}
	// 10. sflow_collector
	if !execute("DELETE", m.getOperationCount("sflow_collector", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "sflow_collector", "DELETE") }, "SFlow Collector") {
		return diagnostics, operationsPerformed
	}
	// 9. eth_port_profile
	if !execute("DELETE", m.getOperationCount("eth_port_profile", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "eth_port_profile", "DELETE") }, "Eth Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 8. service
	if !execute("DELETE", m.getOperationCount("service", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "service", "DELETE") }, "Service") {
		return diagnostics, operationsPerformed
	}
	// 7. port_acl
	if !execute("DELETE", m.getOperationCount("port_acl", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "port_acl", "DELETE") }, "Port ACL") {
		return diagnostics, operationsPerformed
	}
	// 6. pb_routing
	if !execute("DELETE", m.getOperationCount("pb_routing", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "pb_routing", "DELETE") }, "PB Routing") {
		return diagnostics, operationsPerformed
	}
	// 5. pb_routing_acl
	if !execute("DELETE", m.getOperationCount("pb_routing_acl", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "pb_routing_acl", "DELETE") }, "PB Routing ACL") {
		return diagnostics, operationsPerformed
	}
	// 3-4. acl (both ipv4 and ipv6)
	if !execute("DELETE", m.getOperationCount("acl", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "acl", "DELETE") }, "ACL") {
		return diagnostics, operationsPerformed
	}
	// 2. ipv6_list
	if !execute("DELETE", m.getOperationCount("ipv6_list", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "ipv6_list", "DELETE") }, "IPv6 List Filter") {
		return diagnostics, operationsPerformed
	}
	// 1. ipv4_list
	if !execute("DELETE", m.getOperationCount("ipv4_list", "DELETE"), func(ctx context.Context) diag.Diagnostics { return m.ExecuteBulk(ctx, "ipv4_list", "DELETE") }, "IPv4 List Filter") {
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
		m.getOperationCount("device_settings", "PUT") == 0 && m.getOperationCount("device_settings", "PATCH") == 0 && m.getOperationCount("device_settings", "DELETE") == 0 &&
		m.getOperationCount("bundle", "PUT") == 0 && m.getOperationCount("bundle", "PATCH") == 0 && m.getOperationCount("bundle", "DELETE") == 0 && m.getOperationCount("authenticated_eth_port", "PUT") == 0 && m.getOperationCount("authenticated_eth_port", "PATCH") == 0 &&
		m.getOperationCount("authenticated_eth_port", "DELETE") == 0 && m.getOperationCount("acl", "PUT") == 0 && m.getOperationCount("acl", "PATCH") == 0 && m.getOperationCount("acl", "DELETE") == 0 &&
		m.getOperationCount("ipv4_list", "PUT") == 0 && m.getOperationCount("ipv4_list", "PATCH") == 0 && m.getOperationCount("ipv4_list", "DELETE") == 0 &&
		m.getOperationCount("ipv4_prefix_list", "PUT") == 0 && m.getOperationCount("ipv4_prefix_list", "PATCH") == 0 && m.getOperationCount("ipv4_prefix_list", "DELETE") == 0 &&
		m.getOperationCount("ipv6_list", "PUT") == 0 && m.getOperationCount("ipv6_list", "PATCH") == 0 && m.getOperationCount("ipv6_list", "DELETE") == 0 &&
		m.getOperationCount("ipv6_prefix_list", "PUT") == 0 && m.getOperationCount("ipv6_prefix_list", "PATCH") == 0 && m.getOperationCount("ipv6_prefix_list", "DELETE") == 0 &&
		m.getOperationCount("badge", "PUT") == 0 && m.getOperationCount("badge", "PATCH") == 0 && m.getOperationCount("badge", "DELETE") == 0 &&
		m.getOperationCount("voice_port_profile", "PUT") == 0 && m.getOperationCount("voice_port_profile", "PATCH") == 0 && m.getOperationCount("voice_port_profile", "DELETE") == 0 &&
		m.getOperationCount("switchpoint", "PUT") == 0 && m.getOperationCount("switchpoint", "PATCH") == 0 && m.getOperationCount("switchpoint", "DELETE") == 0 &&
		m.getOperationCount("service_port_profile", "PUT") == 0 && m.getOperationCount("service_port_profile", "PATCH") == 0 && m.getOperationCount("service_port_profile", "DELETE") == 0 &&
		m.getOperationCount("packet_broker", "PUT") == 0 && m.getOperationCount("packet_broker", "PATCH") == 0 && m.getOperationCount("packet_broker", "DELETE") == 0 &&
		m.getOperationCount("packet_queue", "PUT") == 0 && m.getOperationCount("packet_queue", "PATCH") == 0 && m.getOperationCount("packet_queue", "DELETE") == 0 &&
		m.getOperationCount("device_voice_settings", "PUT") == 0 && m.getOperationCount("device_voice_settings", "PATCH") == 0 && m.getOperationCount("device_voice_settings", "DELETE") == 0 &&
		m.getOperationCount("as_path_access_list", "PUT") == 0 && m.getOperationCount("as_path_access_list", "PATCH") == 0 && m.getOperationCount("as_path_access_list", "DELETE") == 0 &&
		m.getOperationCount("community_list", "PUT") == 0 && m.getOperationCount("community_list", "PATCH") == 0 && m.getOperationCount("community_list", "DELETE") == 0 &&
		m.getOperationCount("extended_community_list", "PUT") == 0 && m.getOperationCount("extended_community_list", "PATCH") == 0 && m.getOperationCount("extended_community_list", "DELETE") == 0 &&
		m.getOperationCount("route_map_clause", "PUT") == 0 && m.getOperationCount("route_map_clause", "PATCH") == 0 && m.getOperationCount("route_map_clause", "DELETE") == 0 &&
		m.getOperationCount("route_map", "PUT") == 0 && m.getOperationCount("route_map", "PATCH") == 0 && m.getOperationCount("route_map", "DELETE") == 0 &&
		m.getOperationCount("sfp_breakout", "PATCH") == 0 &&
		m.getOperationCount("site", "PATCH") == 0 &&
		m.getOperationCount("pod", "PUT") == 0 && m.getOperationCount("pod", "PATCH") == 0 && m.getOperationCount("pod", "DELETE") == 0 &&
		m.getOperationCount("pb_routing", "PUT") == 0 && m.getOperationCount("pb_routing", "PATCH") == 0 && m.getOperationCount("pb_routing", "DELETE") == 0 &&
		m.getOperationCount("pb_routing_acl", "PUT") == 0 && m.getOperationCount("pb_routing_acl", "PATCH") == 0 && m.getOperationCount("pb_routing_acl", "DELETE") == 0 &&
		m.getOperationCount("spine_plane", "PUT") == 0 && m.getOperationCount("spine_plane", "PATCH") == 0 && m.getOperationCount("spine_plane", "DELETE") == 0 &&
		m.getOperationCount("port_acl", "PUT") == 0 && m.getOperationCount("port_acl", "PATCH") == 0 && m.getOperationCount("port_acl", "DELETE") == 0 &&
		m.getOperationCount("sflow_collector", "PUT") == 0 && m.getOperationCount("sflow_collector", "PATCH") == 0 && m.getOperationCount("sflow_collector", "DELETE") == 0 &&
		m.getOperationCount("diagnostics_profile", "PUT") == 0 && m.getOperationCount("diagnostics_profile", "PATCH") == 0 && m.getOperationCount("diagnostics_profile", "DELETE") == 0 &&
		m.getOperationCount("diagnostics_port_profile", "PUT") == 0 && m.getOperationCount("diagnostics_port_profile", "PATCH") == 0 && m.getOperationCount("diagnostics_port_profile", "DELETE") == 0 &&
		m.getOperationCount("device_controller", "PUT") == 0 && m.getOperationCount("device_controller", "PATCH") == 0 && m.getOperationCount("device_controller", "DELETE") == 0 &&
		m.getOperationCount("grouping_rule", "PUT") == 0 && m.getOperationCount("grouping_rule", "PATCH") == 0 && m.getOperationCount("grouping_rule", "DELETE") == 0 &&
		m.getOperationCount("threshold_group", "PUT") == 0 && m.getOperationCount("threshold_group", "PATCH") == 0 && m.getOperationCount("threshold_group", "DELETE") == 0 &&
		m.getOperationCount("threshold", "PUT") == 0 && m.getOperationCount("threshold", "PATCH") == 0 && m.getOperationCount("threshold", "DELETE") == 0 {
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

	deviceSettingsPutCount := m.getOperationCount("device_settings", "PUT")
	deviceSettingsPatchCount := m.getOperationCount("device_settings", "PATCH")
	deviceSettingsDeleteCount := m.getOperationCount("device_settings", "DELETE")

	sflowCollectorPutCount := m.getOperationCount("sflow_collector", "PUT")
	sflowCollectorPatchCount := m.getOperationCount("sflow_collector", "PATCH")
	sflowCollectorDeleteCount := m.getOperationCount("sflow_collector", "DELETE")

	diagnosticsProfilePutCount := m.getOperationCount("diagnostics_profile", "PUT")
	diagnosticsProfilePatchCount := m.getOperationCount("diagnostics_profile", "PATCH")
	diagnosticsProfileDeleteCount := m.getOperationCount("diagnostics_profile", "DELETE")

	diagnosticsPortProfilePutCount := m.getOperationCount("diagnostics_port_profile", "PUT")
	diagnosticsPortProfilePatchCount := m.getOperationCount("diagnostics_port_profile", "PATCH")
	diagnosticsPortProfileDeleteCount := m.getOperationCount("diagnostics_port_profile", "DELETE")

	bundlePutCount := m.getOperationCount("bundle", "PUT")
	bundlePatchCount := m.getOperationCount("bundle", "PATCH")
	bundleDeleteCount := m.getOperationCount("bundle", "DELETE")

	aclPutCount := m.getOperationCount("acl", "PUT")
	aclPatchCount := m.getOperationCount("acl", "PATCH")
	aclDeleteCount := m.getOperationCount("acl", "DELETE")

	ipv4ListPutCount := m.getOperationCount("ipv4_list", "PUT")
	ipv4ListPatchCount := m.getOperationCount("ipv4_list", "PATCH")
	ipv4ListDeleteCount := m.getOperationCount("ipv4_list", "DELETE")

	ipv4PrefixListPutCount := m.getOperationCount("ipv4_prefix_list", "PUT")
	ipv4PrefixListPatchCount := m.getOperationCount("ipv4_prefix_list", "PATCH")
	ipv4PrefixListDeleteCount := m.getOperationCount("ipv4_prefix_list", "DELETE")

	ipv6ListPutCount := m.getOperationCount("ipv6_list", "PUT")
	ipv6ListPatchCount := m.getOperationCount("ipv6_list", "PATCH")
	ipv6ListDeleteCount := m.getOperationCount("ipv6_list", "DELETE")

	ipv6PrefixListPutCount := m.getOperationCount("ipv6_prefix_list", "PUT")
	ipv6PrefixListPatchCount := m.getOperationCount("ipv6_prefix_list", "PATCH")
	ipv6PrefixListDeleteCount := m.getOperationCount("ipv6_prefix_list", "DELETE")

	authenticatedEthPortPutCount := m.getOperationCount("authenticated_eth_port", "PUT")
	authenticatedEthPortPatchCount := m.getOperationCount("authenticated_eth_port", "PATCH")
	authenticatedEthPortDeleteCount := m.getOperationCount("authenticated_eth_port", "DELETE")

	badgePutCount := m.getOperationCount("badge", "PUT")
	badgePatchCount := m.getOperationCount("badge", "PATCH")
	badgeDeleteCount := m.getOperationCount("badge", "DELETE")

	voicePortProfilePutCount := m.getOperationCount("voice_port_profile", "PUT")
	voicePortProfilePatchCount := m.getOperationCount("voice_port_profile", "PATCH")
	voicePortProfileDeleteCount := m.getOperationCount("voice_port_profile", "DELETE")

	switchpointPutCount := m.getOperationCount("switchpoint", "PUT")
	switchpointPatchCount := m.getOperationCount("switchpoint", "PATCH")
	switchpointDeleteCount := m.getOperationCount("switchpoint", "DELETE")

	servicePortProfilePutCount := m.getOperationCount("service_port_profile", "PUT")
	servicePortProfilePatchCount := m.getOperationCount("service_port_profile", "PATCH")
	servicePortProfileDeleteCount := m.getOperationCount("service_port_profile", "DELETE")

	packetBrokerPutCount := m.getOperationCount("packet_broker", "PUT")
	packetBrokerPatchCount := m.getOperationCount("packet_broker", "PATCH")
	packetBrokerDeleteCount := m.getOperationCount("packet_broker", "DELETE")

	packetQueuePutCount := m.getOperationCount("packet_queue", "PUT")
	packetQueuePatchCount := m.getOperationCount("packet_queue", "PATCH")
	packetQueueDeleteCount := m.getOperationCount("packet_queue", "DELETE")

	deviceVoiceSettingsPutCount := m.getOperationCount("device_voice_settings", "PUT")
	deviceVoiceSettingsPatchCount := m.getOperationCount("device_voice_settings", "PATCH")
	deviceVoiceSettingsDeleteCount := m.getOperationCount("device_voice_settings", "DELETE")

	asPathAccessListPutCount := m.getOperationCount("as_path_access_list", "PUT")
	asPathAccessListPatchCount := m.getOperationCount("as_path_access_list", "PATCH")
	asPathAccessListDeleteCount := m.getOperationCount("as_path_access_list", "DELETE")

	communityListPutCount := m.getOperationCount("community_list", "PUT")
	communityListPatchCount := m.getOperationCount("community_list", "PATCH")
	communityListDeleteCount := m.getOperationCount("community_list", "DELETE")

	extendedCommunityListPutCount := m.getOperationCount("extended_community_list", "PUT")
	extendedCommunityListPatchCount := m.getOperationCount("extended_community_list", "PATCH")
	extendedCommunityListDeleteCount := m.getOperationCount("extended_community_list", "DELETE")

	routeMapClausePutCount := m.getOperationCount("route_map_clause", "PUT")
	routeMapClausePatchCount := m.getOperationCount("route_map_clause", "PATCH")
	routeMapClauseDeleteCount := m.getOperationCount("route_map_clause", "DELETE")

	routeMapPutCount := m.getOperationCount("route_map", "PUT")
	routeMapPatchCount := m.getOperationCount("route_map", "PATCH")
	routeMapDeleteCount := m.getOperationCount("route_map", "DELETE")

	sfpBreakoutPatchCount := m.getOperationCount("sfp_breakout", "PATCH")

	sitePatchCount := m.getOperationCount("site", "PATCH")

	podPutCount := m.getOperationCount("pod", "PUT")
	podPatchCount := m.getOperationCount("pod", "PATCH")
	podDeleteCount := m.getOperationCount("pod", "DELETE")

	pbRoutingPutCount := m.getOperationCount("pb_routing", "PUT")
	pbRoutingPatchCount := m.getOperationCount("pb_routing", "PATCH")
	pbRoutingDeleteCount := m.getOperationCount("pb_routing", "DELETE")

	pbRoutingAclPutCount := m.getOperationCount("pb_routing_acl", "PUT")
	pbRoutingAclPatchCount := m.getOperationCount("pb_routing_acl", "PATCH")
	pbRoutingAclDeleteCount := m.getOperationCount("pb_routing_acl", "DELETE")

	spinePlanePutCount := m.getOperationCount("spine_plane", "PUT")
	spinePlanePatchCount := m.getOperationCount("spine_plane", "PATCH")
	spinePlaneDeleteCount := m.getOperationCount("spine_plane", "DELETE")

	portAclPutCount := m.getOperationCount("port_acl", "PUT")
	portAclPatchCount := m.getOperationCount("port_acl", "PATCH")
	portAclDeleteCount := m.getOperationCount("port_acl", "DELETE")

	deviceControllerPutCount := m.getOperationCount("device_controller", "PUT")
	deviceControllerPatchCount := m.getOperationCount("device_controller", "PATCH")
	deviceControllerDeleteCount := m.getOperationCount("device_controller", "DELETE")

	groupingRulePutCount := m.getOperationCount("grouping_rule", "PUT")
	groupingRulePatchCount := m.getOperationCount("grouping_rule", "PATCH")
	groupingRuleDeleteCount := m.getOperationCount("grouping_rule", "DELETE")

	thresholdGroupPutCount := m.getOperationCount("threshold_group", "PUT")
	thresholdGroupPatchCount := m.getOperationCount("threshold_group", "PATCH")
	thresholdGroupDeleteCount := m.getOperationCount("threshold_group", "DELETE")

	thresholdPutCount := m.getOperationCount("threshold", "PUT")
	thresholdPatchCount := m.getOperationCount("threshold", "PATCH")
	thresholdDeleteCount := m.getOperationCount("threshold", "DELETE")

	m.mutex.Unlock()

	totalCount := gatewayPutCount + gatewayPatchCount + gatewayDeleteCount +
		lagPutCount + lagPatchCount + lagDeleteCount +
		tenantPutCount + tenantPatchCount + tenantDeleteCount +
		servicePutCount + servicePatchCount + serviceDeleteCount +
		gatewayProfilePutCount + gatewayProfilePatchCount + gatewayProfileDeleteCount +
		ethPortProfilePutCount + ethPortProfilePatchCount + ethPortProfileDeleteCount +
		ethPortSettingsPutCount + ethPortSettingsPatchCount + ethPortSettingsDeleteCount +
		deviceSettingsPutCount + deviceSettingsPatchCount + deviceSettingsDeleteCount +
		sflowCollectorPutCount + sflowCollectorPatchCount + sflowCollectorDeleteCount +
		diagnosticsProfilePutCount + diagnosticsProfilePatchCount + diagnosticsProfileDeleteCount +
		diagnosticsPortProfilePutCount + diagnosticsPortProfilePatchCount + diagnosticsPortProfileDeleteCount +
		bundlePutCount + bundlePatchCount + bundleDeleteCount + aclPutCount + aclPatchCount + aclDeleteCount +
		ipv4ListPutCount + ipv4ListPatchCount + ipv4ListDeleteCount +
		ipv4PrefixListPutCount + ipv4PrefixListPatchCount + ipv4PrefixListDeleteCount +
		ipv6ListPutCount + ipv6ListPatchCount + ipv6ListDeleteCount +
		ipv6PrefixListPutCount + ipv6PrefixListPatchCount + ipv6PrefixListDeleteCount +
		badgePutCount + badgePatchCount + badgeDeleteCount +
		voicePortProfilePutCount + voicePortProfilePatchCount + voicePortProfileDeleteCount +
		switchpointPutCount + switchpointPatchCount + switchpointDeleteCount +
		servicePortProfilePutCount + servicePortProfilePatchCount + servicePortProfileDeleteCount +
		packetBrokerPutCount + packetBrokerPatchCount + packetBrokerDeleteCount +
		packetQueuePutCount + packetQueuePatchCount + packetQueueDeleteCount +
		deviceVoiceSettingsPutCount + deviceVoiceSettingsPatchCount + deviceVoiceSettingsDeleteCount +
		asPathAccessListPutCount + asPathAccessListPatchCount + asPathAccessListDeleteCount +
		communityListPutCount + communityListPatchCount + communityListDeleteCount +
		extendedCommunityListPutCount + extendedCommunityListPatchCount + extendedCommunityListDeleteCount +
		routeMapClausePutCount + routeMapClausePatchCount + routeMapClauseDeleteCount +
		routeMapPutCount + routeMapPatchCount + routeMapDeleteCount +
		sfpBreakoutPatchCount + sitePatchCount +
		podPutCount + podPatchCount + podDeleteCount +
		pbRoutingPutCount + pbRoutingPatchCount + pbRoutingDeleteCount +
		pbRoutingAclPutCount + pbRoutingAclPatchCount + pbRoutingAclDeleteCount +
		spinePlanePutCount + spinePlanePatchCount + spinePlaneDeleteCount +
		portAclPutCount + portAclPatchCount + portAclDeleteCount +
		authenticatedEthPortPutCount + authenticatedEthPortPatchCount + authenticatedEthPortDeleteCount +
		deviceControllerPutCount + deviceControllerPatchCount + deviceControllerDeleteCount +
		groupingRulePutCount + groupingRulePatchCount + groupingRuleDeleteCount +
		thresholdGroupPutCount + thresholdGroupPatchCount + thresholdGroupDeleteCount +
		thresholdPutCount + thresholdPatchCount + thresholdDeleteCount

	if totalCount > 0 {
		tflog.Debug(ctx, "Multiple operations detected, executing in sequence", map[string]interface{}{
			"gateway_put_count":                     gatewayPutCount,
			"gateway_patch_count":                   gatewayPatchCount,
			"gateway_delete_count":                  gatewayDeleteCount,
			"lag_put_count":                         lagPutCount,
			"lag_patch_count":                       lagPatchCount,
			"lag_delete_count":                      lagDeleteCount,
			"tenant_put_count":                      tenantPutCount,
			"tenant_patch_count":                    tenantPatchCount,
			"tenant_delete_count":                   tenantDeleteCount,
			"service_put_count":                     servicePutCount,
			"service_patch_count":                   servicePatchCount,
			"service_delete_count":                  serviceDeleteCount,
			"gateway_profile_put_count":             gatewayProfilePutCount,
			"gateway_profile_patch_count":           gatewayProfilePatchCount,
			"gateway_profile_delete_count":          gatewayProfileDeleteCount,
			"eth_port_profile_put_count":            ethPortProfilePutCount,
			"eth_port_profile_patch_count":          ethPortProfilePatchCount,
			"eth_port_profile_delete_count":         ethPortProfileDeleteCount,
			"eth_port_settings_put_count":           ethPortSettingsPutCount,
			"eth_port_settings_patch_count":         ethPortSettingsPatchCount,
			"eth_port_settings_delete_count":        ethPortSettingsDeleteCount,
			"device_settings_put_count":             deviceSettingsPutCount,
			"device_settings_patch_count":           deviceSettingsPatchCount,
			"device_settings_delete_count":          deviceSettingsDeleteCount,
			"bundle_put_count":                      bundlePutCount,
			"bundle_patch_count":                    bundlePatchCount,
			"bundle_delete_count":                   bundleDeleteCount,
			"acl_put_count":                         aclPutCount,
			"acl_patch_count":                       aclPatchCount,
			"acl_delete_count":                      aclDeleteCount,
			"ipv4_list_put_count":                   ipv4ListPutCount,
			"ipv4_list_patch_count":                 ipv4ListPatchCount,
			"ipv4_list_delete_count":                ipv4ListDeleteCount,
			"ipv4_prefix_list_put_count":            ipv4PrefixListPutCount,
			"ipv4_prefix_list_patch_count":          ipv4PrefixListPatchCount,
			"ipv4_prefix_list_delete_count":         ipv4PrefixListDeleteCount,
			"ipv6_list_put_count":                   ipv6ListPutCount,
			"ipv6_list_patch_count":                 ipv6ListPatchCount,
			"ipv6_list_delete_count":                ipv6ListDeleteCount,
			"ipv6_prefix_list_put_count":            ipv6PrefixListPutCount,
			"ipv6_prefix_list_patch_count":          ipv6PrefixListPatchCount,
			"ipv6_prefix_list_delete_count":         ipv6PrefixListDeleteCount,
			"badge_put_count":                       badgePutCount,
			"badge_patch_count":                     badgePatchCount,
			"badge_delete_count":                    badgeDeleteCount,
			"voice_port_profile_put_count":          voicePortProfilePutCount,
			"voice_port_profile_patch_count":        voicePortProfilePatchCount,
			"voice_port_profile_delete_count":       voicePortProfileDeleteCount,
			"switchpoint_put_count":                 switchpointPutCount,
			"switchpoint_patch_count":               switchpointPatchCount,
			"switchpoint_delete_count":              switchpointDeleteCount,
			"service_port_profile_put_count":        servicePortProfilePutCount,
			"service_port_profile_patch_count":      servicePortProfilePatchCount,
			"service_port_profile_delete_count":     servicePortProfileDeleteCount,
			"packet_broker_put_count":               packetBrokerPutCount,
			"packet_broker_patch_count":             packetBrokerPatchCount,
			"packet_broker_delete_count":            packetBrokerDeleteCount,
			"packet_queue_put_count":                packetQueuePutCount,
			"packet_queue_patch_count":              packetQueuePatchCount,
			"packet_queue_delete_count":             packetQueueDeleteCount,
			"device_voice_settings_put_count":       deviceVoiceSettingsPutCount,
			"device_voice_settings_patch_count":     deviceVoiceSettingsPatchCount,
			"device_voice_settings_delete_count":    deviceVoiceSettingsDeleteCount,
			"as_path_access_list_put_count":         asPathAccessListPutCount,
			"as_path_access_list_patch_count":       asPathAccessListPatchCount,
			"as_path_access_list_delete_count":      asPathAccessListDeleteCount,
			"community_list_put_count":              communityListPutCount,
			"community_list_patch_count":            communityListPatchCount,
			"community_list_delete_count":           communityListDeleteCount,
			"extended_community_list_put_count":     extendedCommunityListPutCount,
			"extended_community_list_patch_count":   extendedCommunityListPatchCount,
			"extended_community_list_delete_count":  extendedCommunityListDeleteCount,
			"route_map_clause_put_count":            routeMapClausePutCount,
			"route_map_clause_patch_count":          routeMapClausePatchCount,
			"route_map_clause_delete_count":         routeMapClauseDeleteCount,
			"route_map_put_count":                   routeMapPutCount,
			"route_map_patch_count":                 routeMapPatchCount,
			"route_map_delete_count":                routeMapDeleteCount,
			"sfp_breakout_patch_count":              sfpBreakoutPatchCount,
			"site_patch_count":                      sitePatchCount,
			"pod_put_count":                         podPutCount,
			"pod_patch_count":                       podPatchCount,
			"pod_delete_count":                      podDeleteCount,
			"pb_routing_put_count":                  pbRoutingPutCount,
			"pb_routing_patch_count":                pbRoutingPatchCount,
			"pb_routing_delete_count":               pbRoutingDeleteCount,
			"pb_routing_acl_put_count":              pbRoutingAclPutCount,
			"pb_routing_acl_patch_count":            pbRoutingAclPatchCount,
			"pb_routing_acl_delete_count":           pbRoutingAclDeleteCount,
			"spine_plane_put_count":                 spinePlanePutCount,
			"spine_plane_patch_count":               spinePlanePatchCount,
			"spine_plane_delete_count":              spinePlaneDeleteCount,
			"port_acl_put_count":                    portAclPutCount,
			"port_acl_patch_count":                  portAclPatchCount,
			"port_acl_delete_count":                 portAclDeleteCount,
			"authenticated_eth_port_put_count":      authenticatedEthPortPutCount,
			"authenticated_eth_port_patch_count":    authenticatedEthPortPatchCount,
			"authenticated_eth_port_delete_count":   authenticatedEthPortDeleteCount,
			"device_controller_put_count":           deviceControllerPutCount,
			"device_controller_patch_count":         deviceControllerPatchCount,
			"device_controller_delete_count":        deviceControllerDeleteCount,
			"sflow_collector_put_count":             sflowCollectorPutCount,
			"sflow_collector_patch_count":           sflowCollectorPatchCount,
			"sflow_collector_delete_count":          sflowCollectorDeleteCount,
			"diagnostics_profile_put_count":         diagnosticsProfilePutCount,
			"diagnostics_profile_patch_count":       diagnosticsProfilePatchCount,
			"diagnostics_profile_delete_count":      diagnosticsProfileDeleteCount,
			"diagnostics_port_profile_put_count":    diagnosticsPortProfilePutCount,
			"diagnostics_port_profile_patch_count":  diagnosticsPortProfilePatchCount,
			"diagnostics_port_profile_delete_count": diagnosticsPortProfileDeleteCount,
			"grouping_rule_put_count":               groupingRulePutCount,
			"grouping_rule_patch_count":             groupingRulePatchCount,
			"grouping_rule_delete_count":            groupingRuleDeleteCount,
			"threshold_group_put_count":             thresholdGroupPutCount,
			"threshold_group_patch_count":           thresholdGroupPatchCount,
			"threshold_group_delete_count":          thresholdGroupDeleteCount,
			"threshold_put_count":                   thresholdPutCount,
			"threshold_patch_count":                 thresholdPatchCount,
			"threshold_delete_count":                thresholdDeleteCount,
			"total_count":                           totalCount,
		})

		return m.ExecuteAllPendingOperations(ctx)
	}

	return nil
}
