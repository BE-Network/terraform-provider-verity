package bulkops

import (
	"context"
	"fmt"
	"sync"
	"terraform-provider-verity/openapi"
	"time"
)

type Manager struct {
	client            *openapi.APIClient
	contextProvider   ContextProviderFunc
	clearCacheFunc    ClearCacheFunc
	mode              string
	mutex             sync.Mutex
	executionMutex    sync.Mutex
	lastOperationTime time.Time
	batchStartTime    time.Time
	resources         map[string]*ResourceOperations

	resourceHeaderParams map[string]map[string]string

	resourceOriginalNames map[string]string

	pendingOperations     map[string]*Operation
	operationResults      map[string]bool
	operationErrors       map[string]error
	operationWaitChannels map[string]chan struct{}
	operationMutex        sync.Mutex
	closedChannels        map[string]bool
}

func initializeResourceOperations() map[string]*ResourceOperations {
	resources := make(map[string]*ResourceOperations)

	for resourceType := range resourceRegistry {
		resources[resourceType] = NewResourceOperations()
	}

	return resources
}

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

func GetManager(client *openapi.APIClient, clearCacheFunc ClearCacheFunc, providerContext interface{}, mode string) *Manager {
	contextProvider := func() interface{} {
		return providerContext
	}

	return NewManager(client, contextProvider, clearCacheFunc, mode)
}

func (m *Manager) GetResourceResponse(resourceType, resourceName string) (map[string]interface{}, bool) {

	if resourceType == "acl_v4" || resourceType == "acl_v6" {
		resourceType = "acl"
	}

	res, exists := m.resources[resourceType]
	if !exists {
		return nil, false
	}

	res.ResponsesMutex.RLock()
	defer res.ResponsesMutex.RUnlock()
	response, exists := res.Responses[resourceName]
	return response, exists
}

func (m *Manager) HasPendingOrRecentOperations(resourceType string) bool {
	return m.hasPendingOrRecentOperations(resourceType)
}

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

	for _, opID := range idsToClose {
		m.safeCloseChannel(opID, false)
	}

	if failCount > 0 {
		fmt.Printf("Failed %d pending operations due to error: %v\n", failCount, err)
	}
}

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
