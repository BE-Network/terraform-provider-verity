package bulkops

import (
	"context"
	"fmt"
	"sync"
	"terraform-provider-verity/openapi"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
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

	flushed     map[string]bool
	lateFlushes map[string]int

	pendingOperations     map[string]*Operation
	operationResults      map[string]bool
	operationErrors       map[string]error
	operationWaitChannels map[string]chan struct{}
	operationMutex        sync.Mutex
	closedChannels        map[string]bool

	heldPuts          map[string]matchVrfValue
	matchVrfHoldArmed bool
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
		flushed:               make(map[string]bool),
		lateFlushes:           make(map[string]int),
		pendingOperations:     make(map[string]*Operation),
		operationResults:      make(map[string]bool),
		operationErrors:       make(map[string]error),
		operationWaitChannels: make(map[string]chan struct{}),
		closedChannels:        make(map[string]bool),
		heldPuts:              make(map[string]matchVrfValue),
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
	if cancelStartedHook != nil {
		cancelStartedHook()
	}
	m.mutex.Lock()
	m.operationMutex.Lock()
	var idsToClose []string
	failCount := 0

	for opID, op := range m.pendingOperations {
		if op.Status == OperationPending {
			updatedOp := op
			updatedOp.Status = OperationFailed
			updatedOp.Error = fmt.Errorf("not sent because %v", err)
			m.pendingOperations[opID] = updatedOp
			m.operationErrors[opID] = updatedOp.Error
			m.operationResults[opID] = false
			m.removeQueuedPayloadLocked(op)
			idsToClose = append(idsToClose, opID)
			failCount++
		}
	}

	m.operationMutex.Unlock()
	m.mutex.Unlock()

	for _, opID := range idsToClose {
		m.safeCloseChannel(opID, false)
	}

	if failCount > 0 {
		tflog.Warn(ctx, fmt.Sprintf("[BULK-OPS] Cancelled %d queued operation(s): %v", failCount, err))
	}
}

var cancelStartedHook func()

func (m *Manager) removeQueuedPayloadLocked(op *Operation) {
	key := op.StoreKey
	if key == "" {
		key = op.ResourceName
	}
	if res, exists := m.resources[op.ResourceType]; exists {
		switch op.OperationType {
		case "PUT":
			delete(res.Put, key)
		case "PATCH":
			delete(res.Patch, key)
		case "DELETE":
			kept := res.Delete[:0]
			for _, name := range res.Delete {
				if name != key {
					kept = append(kept, name)
				}
			}
			res.Delete = kept
		}
	}
	paramKey := fmt.Sprintf("%s:%s", op.ResourceType, key)
	delete(m.resourceHeaderParams, paramKey)
	delete(m.resourceOriginalNames, paramKey)
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

func (m *Manager) noteFlush(ctx context.Context, resourceType, operationType string, count int) {
	if count == 0 {
		return
	}
	key := resourceType + " " + operationType
	m.mutex.Lock()
	late := m.flushed[key]
	m.flushed[key] = true
	if late {
		m.lateFlushes[key]++
	}
	m.mutex.Unlock()
	if late {
		tflog.Warn(ctx, fmt.Sprintf("[BULK-OPS] %d %s %s operation(s) arrived after the %s %s request was already sent, so they are sent in a separate request", count, resourceType, operationType, resourceType, operationType))
	}
}

func (m *Manager) LateFlushes() map[string]int {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	result := make(map[string]int, len(m.lateFlushes))
	for key, count := range m.lateFlushes {
		result[key] = count
	}
	return result
}
