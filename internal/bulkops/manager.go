package bulkops

import (
	"context"
	"fmt"
	"sync"
	"terraform-provider-verity/openapi"
	"time"
)

// ================================================================================================
// MANAGER - MAIN STRUCT
// ================================================================================================

// Manager coordinates bulk operations across multiple resource types.
// It provides efficient batching, automatic retry, dependency ordering, and response caching.
//
// Key Features:
//   - Generic operations for multiple resource types
//   - Automatic batching with configurable timing windows
//   - Retry logic with exponential backoff for resilient operations
//   - Response caching for auto-generated fields
//   - Thread-safe operation queuing and execution
//   - Dependency-aware operation ordering (datacenter vs campus modes)
//   - Reflection-based generic data access patterns

type Manager struct {
	client            *openapi.APIClient
	contextProvider   ContextProviderFunc
	clearCacheFunc    ClearCacheFunc
	mode              string
	mutex             sync.Mutex
	lastOperationTime time.Time
	batchStartTime    time.Time
	resources         map[string]*ResourceOperations

	// resourceHeaderParams tracks header parameters for operations that need them
	// Key format: "resourceType:compositeKey" -> map of header params
	// Example: "acl:my_filter_ip_version4" -> {"ip_version": "4"}
	resourceHeaderParams map[string]map[string]string

	// resourceOriginalNames tracks the original resource names for composite keys
	// Key format: "resourceType:compositeKey" -> original resource name
	// Example: "acl:my_filter_ip_version4" -> "my_filter"
	resourceOriginalNames map[string]string

	// For tracking operations
	pendingOperations     map[string]*Operation
	operationResults      map[string]bool // true = success, false = failure
	operationErrors       map[string]error
	operationWaitChannels map[string]chan struct{}
	operationMutex        sync.Mutex
	closedChannels        map[string]bool
}

// ================================================================================================
// CONSTRUCTOR AND INITIALIZATION
// ================================================================================================

// initializeResourceOperations creates ResourceOperations for all registered resource types.
func initializeResourceOperations() map[string]*ResourceOperations {
	resources := make(map[string]*ResourceOperations)

	// Initialize ResourceOperations for each resource type in the registry
	for resourceType := range resourceRegistry {
		resources[resourceType] = NewResourceOperations()
	}

	return resources
}

// NewManager creates a new bulk operation manager.
func NewManager(client *openapi.APIClient, contextProvider ContextProviderFunc, clearCacheFunc ClearCacheFunc, mode string) *Manager {
	return &Manager{
		client:                client,
		contextProvider:       contextProvider,
		clearCacheFunc:        clearCacheFunc,
		mode:                  mode,
		lastOperationTime:     time.Now(),
		resources:             initializeResourceOperations(),
		resourceHeaderParams:  make(map[string]map[string]string),
		resourceOriginalNames: make(map[string]string),
		pendingOperations:     make(map[string]*Operation),
		operationResults:      make(map[string]bool),
		operationErrors:       make(map[string]error),
		operationWaitChannels: make(map[string]chan struct{}),
		closedChannels:        make(map[string]bool),
	}
}

// GetManager creates a bulk operation manager with a context provider wrapper.
func GetManager(client *openapi.APIClient, clearCacheFunc ClearCacheFunc, providerContext interface{}, mode string) *Manager {
	contextProvider := func() interface{} {
		return providerContext
	}

	return NewManager(client, contextProvider, clearCacheFunc, mode)
}

// ================================================================================================
// PUBLIC ACCESSOR METHODS
// ================================================================================================

// GetResourceResponse retrieves cached response data for a resource.
func (m *Manager) GetResourceResponse(resourceType, resourceName string) (map[string]interface{}, bool) {
	// Handle special case for ACL versioning - map to base "acl" type
	if resourceType == "acl_v4" || resourceType == "acl_v6" {
		resourceType = "acl"
	}

	// Get the resource operations from the unified map
	res, exists := m.resources[resourceType]
	if !exists {
		return nil, false
	}

	res.ResponsesMutex.RLock()
	defer res.ResponsesMutex.RUnlock()
	response, exists := res.Responses[resourceName]
	return response, exists
}

// HasPendingOrRecentOperations checks if a resource type has pending or recent operations.
func (m *Manager) HasPendingOrRecentOperations(resourceType string) bool {
	return m.hasPendingOrRecentOperations(resourceType)
}

// FailAllPendingOperations marks all pending operations as failed.
func (m *Manager) FailAllPendingOperations(ctx context.Context, err error) {
	m.operationMutex.Lock()
	var idsToClose []string
	failCount := 0

	for opID, op := range m.pendingOperations {
		if op.Status == OperationPending {
			updatedOp := op
			updatedOp.Status = OperationFailed
			updatedOp.Error = fmt.Errorf("operation aborted due to previous failure: %v", err)
			m.pendingOperations[opID] = updatedOp
			m.operationErrors[opID] = updatedOp.Error
			m.operationResults[opID] = false
			idsToClose = append(idsToClose, opID)
			failCount++
		}
	}

	m.operationMutex.Unlock()

	// Close channels outside the lock to avoid deadlock
	for _, opID := range idsToClose {
		m.safeCloseChannel(opID, false)
	}

	if failCount > 0 {
		fmt.Printf("Failed %d pending operations due to error: %v\n", failCount, err)
	}
}

// safeCloseChannel safely closes an operation's wait channel.
// The lockAlreadyHeld parameter indicates if the caller already holds the operationMutex.
func (m *Manager) safeCloseChannel(opID string, lockAlreadyHeld bool) {
	if !lockAlreadyHeld {
		m.operationMutex.Lock()
		defer m.operationMutex.Unlock()
	}

	if waitCh, ok := m.operationWaitChannels[opID]; ok {
		if _, closed := m.closedChannels[opID]; !closed {
			close(waitCh)
			m.closedChannels[opID] = true
		}
	}
}
