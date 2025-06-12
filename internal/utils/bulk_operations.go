package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"terraform-provider-verity/openapi"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type ContextProviderFunc func() interface{}

type ClearCacheFunc func(ctx context.Context, provider interface{}, cacheKey string)

type OperationStatus int

type Operation struct {
	ResourceType  string
	ResourceName  string
	OperationType string
	Status        OperationStatus
	Error         error
}

type ResourceExistenceCheck struct {
	FetchResources func(ctx context.Context) (map[string]interface{}, error)
	ResourceType   string
	OperationType  string
}

type BulkOperationConfig struct {
	ResourceType      string
	OperationType     string
	ExtractOperations func() (map[string]interface{}, []string)
	CheckPreExistence func(ctx context.Context, resourceNames []string) ([]string, map[string]interface{}, error)
	PrepareRequest    func(filteredData map[string]interface{}) interface{}
	ExecuteRequest    func(ctx context.Context, request interface{}) (*http.Response, error)
	ProcessResponse   func(ctx context.Context, resp *http.Response) error
	UpdateRecentOps   func()
}

const (
	OperationPending OperationStatus = iota
	OperationSucceeded
	OperationFailed
)

const (
	MaxBatchSize          = 1000
	DefaultBatchDelay     = 2 * time.Second
	BatchCollectionWindow = 2000 * time.Millisecond
	MaxBatchDelay         = 5 * time.Second
	OperationTimeout      = 300 * time.Second
)

type BulkOperationManager struct {
	client            *openapi.APIClient
	contextProvider   ContextProviderFunc
	clearCacheFunc    ClearCacheFunc
	mutex             sync.Mutex
	lastOperationTime time.Time
	batchStartTime    time.Time

	// Gateway operations
	gatewayPut    map[string]openapi.ConfigPutRequestGatewayGatewayName
	gatewayPatch  map[string]openapi.ConfigPutRequestGatewayGatewayName
	gatewayDelete []string

	// LAG operations
	lagPut    map[string]openapi.ConfigPutRequestLagLagName
	lagPatch  map[string]openapi.ConfigPutRequestLagLagName
	lagDelete []string

	// Tenant operations
	tenantPut    map[string]openapi.ConfigPutRequestTenantTenantName
	tenantPatch  map[string]openapi.ConfigPutRequestTenantTenantName
	tenantDelete []string

	// Service operations
	servicePut    map[string]openapi.ConfigPutRequestServiceServiceName
	servicePatch  map[string]openapi.ConfigPutRequestServiceServiceName
	serviceDelete []string

	// Gateway Profile operations
	gatewayProfilePut    map[string]openapi.ConfigPutRequestGatewayProfileGatewayProfileName
	gatewayProfilePatch  map[string]openapi.ConfigPutRequestGatewayProfileGatewayProfileName
	gatewayProfileDelete []string

	// EthPortProfile operations
	ethPortProfilePut    map[string]openapi.ConfigPutRequestEthPortProfileEthPortProfileName
	ethPortProfilePatch  map[string]openapi.ConfigPutRequestEthPortProfileEthPortProfileName
	ethPortProfileDelete []string

	// EthPortSettings operations
	ethPortSettingsPut    map[string]openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName
	ethPortSettingsPatch  map[string]openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName
	ethPortSettingsDelete []string

	// Bundles operations
	bundlePatch map[string]openapi.BundlesPatchRequestEndpointBundleValue

	// Track recent operations to avoid race conditions
	recentGatewayOps            bool
	recentGatewayOpTime         time.Time
	recentLagOps                bool
	recentLagOpTime             time.Time
	recentServiceOps            bool
	recentServiceOpTime         time.Time
	recentTenantOps             bool
	recentTenantOpTime          time.Time
	recentGatewayProfileOps     bool
	recentGatewayProfileOpTime  time.Time
	recentEthPortProfileOps     bool
	recentEthPortProfileOpTime  time.Time
	recentEthPortSettingsOps    bool
	recentEthPortSettingsOpTime time.Time
	recentBundleOps             bool
	recentBundleOpTime          time.Time

	// For tracking operations
	pendingOperations     map[string]*Operation
	operationResults      map[string]bool // true = success, false = failure
	operationErrors       map[string]error
	operationWaitChannels map[string]chan struct{}
	operationMutex        sync.Mutex
	closedChannels        map[string]bool

	// Store API responses
	gatewayResponses              map[string]map[string]interface{}
	gatewayResponsesMutex         sync.RWMutex
	lagResponses                  map[string]map[string]interface{}
	lagResponsesMutex             sync.RWMutex
	serviceResponses              map[string]map[string]interface{}
	serviceResponsesMutex         sync.RWMutex
	tenantResponses               map[string]map[string]interface{}
	tenantResponsesMutex          sync.RWMutex
	gatewayProfileResponses       map[string]map[string]interface{}
	gatewayProfileResponsesMutex  sync.RWMutex
	ethPortProfileResponses       map[string]map[string]interface{}
	ethPortProfileResponsesMutex  sync.RWMutex
	ethPortSettingsResponses      map[string]map[string]interface{}
	ethPortSettingsResponsesMutex sync.RWMutex
	bundleResponses               map[string]map[string]interface{}
	bundleResponsesMutex          sync.RWMutex
}

func (b *BulkOperationManager) FilterPreExistingResources(
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
		b.operationMutex.Lock()
		defer b.operationMutex.Unlock()

		for opID, op := range b.pendingOperations {
			if op.ResourceType == checker.ResourceType &&
				op.OperationType == checker.OperationType &&
				alreadyExistingResources[op.ResourceName] {
				// Mark operation as successful
				updatedOp := *op
				updatedOp.Status = OperationSucceeded
				b.pendingOperations[opID] = &updatedOp
				b.operationResults[opID] = true

				b.safeCloseChannel(opID, true)

				tflog.Debug(ctx, fmt.Sprintf("Marked operation %s as successful since resource already exists", opID))
			}
		}
	}

	return notExistingResources, nil
}

func (b *BulkOperationManager) safeCloseChannel(opID string, lockAlreadyHeld ...bool) {
	// Only lock if the caller doesn't already hold the lock
	alreadyLocked := len(lockAlreadyHeld) > 0 && lockAlreadyHeld[0]
	if !alreadyLocked {
		b.operationMutex.Lock()
		defer b.operationMutex.Unlock()
	}

	if waitCh, ok := b.operationWaitChannels[opID]; ok {
		if _, closed := b.closedChannels[opID]; !closed {
			close(waitCh)
			b.closedChannels[opID] = true
		}
	}
}

func generateOperationID(resourceType, resourceName, operationType string) string {
	return fmt.Sprintf("%s-%s-%s-%s", resourceType, resourceName, operationType, uuid.New().String())
}

func (b *BulkOperationManager) WaitForOperation(ctx context.Context, operationID string, timeout time.Duration) error {
	b.operationMutex.Lock()
	waitCh, exists := b.operationWaitChannels[operationID]
	if !exists {
		b.operationMutex.Unlock()
		return fmt.Errorf("operation %s not found", operationID)
	}

	if closed, ok := b.closedChannels[operationID]; ok && closed {
		var err error
		if errorVal, hasError := b.operationErrors[operationID]; hasError {
			err = errorVal
		}
		b.operationMutex.Unlock()
		return err
	}
	b.operationMutex.Unlock()

	select {
	case <-waitCh:
		// Operation completed
		b.operationMutex.Lock()
		defer b.operationMutex.Unlock()

		if err, hasError := b.operationErrors[operationID]; hasError {
			return err
		}
		return nil

	case <-time.After(timeout):
		return fmt.Errorf("timeout waiting for operation %s", operationID)

	case <-ctx.Done():
		return ctx.Err()
	}
}

func NewBulkOperationManager(client *openapi.APIClient, contextProvider ContextProviderFunc, clearCacheFunc ClearCacheFunc) *BulkOperationManager {
	return &BulkOperationManager{
		client:                client,
		contextProvider:       contextProvider,
		clearCacheFunc:        clearCacheFunc,
		lastOperationTime:     time.Now(),
		gatewayPut:            make(map[string]openapi.ConfigPutRequestGatewayGatewayName),
		gatewayPatch:          make(map[string]openapi.ConfigPutRequestGatewayGatewayName),
		gatewayDelete:         make([]string, 0),
		lagPut:                make(map[string]openapi.ConfigPutRequestLagLagName),
		lagPatch:              make(map[string]openapi.ConfigPutRequestLagLagName),
		lagDelete:             make([]string, 0),
		tenantPut:             make(map[string]openapi.ConfigPutRequestTenantTenantName),
		tenantPatch:           make(map[string]openapi.ConfigPutRequestTenantTenantName),
		tenantDelete:          make([]string, 0),
		servicePut:            make(map[string]openapi.ConfigPutRequestServiceServiceName),
		servicePatch:          make(map[string]openapi.ConfigPutRequestServiceServiceName),
		serviceDelete:         make([]string, 0),
		gatewayProfilePut:     make(map[string]openapi.ConfigPutRequestGatewayProfileGatewayProfileName),
		gatewayProfilePatch:   make(map[string]openapi.ConfigPutRequestGatewayProfileGatewayProfileName),
		gatewayProfileDelete:  make([]string, 0),
		ethPortProfilePut:     make(map[string]openapi.ConfigPutRequestEthPortProfileEthPortProfileName),
		ethPortProfilePatch:   make(map[string]openapi.ConfigPutRequestEthPortProfileEthPortProfileName),
		ethPortProfileDelete:  make([]string, 0),
		ethPortSettingsPut:    make(map[string]openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName),
		ethPortSettingsPatch:  make(map[string]openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName),
		ethPortSettingsDelete: make([]string, 0),
		bundlePatch:           make(map[string]openapi.BundlesPatchRequestEndpointBundleValue),
		pendingOperations:     make(map[string]*Operation),
		operationResults:      make(map[string]bool),
		operationErrors:       make(map[string]error),
		operationWaitChannels: make(map[string]chan struct{}),
		closedChannels:        make(map[string]bool),

		// Initialize with no recent operations
		recentGatewayOps:         false,
		recentLagOps:             false,
		recentServiceOps:         false,
		recentTenantOps:          false,
		recentGatewayProfileOps:  false,
		recentEthPortProfileOps:  false,
		recentEthPortSettingsOps: false,
		recentBundleOps:          false,

		// Initialize response caches
		gatewayResponses:              make(map[string]map[string]interface{}),
		gatewayResponsesMutex:         sync.RWMutex{},
		lagResponses:                  make(map[string]map[string]interface{}),
		lagResponsesMutex:             sync.RWMutex{},
		serviceResponses:              make(map[string]map[string]interface{}),
		serviceResponsesMutex:         sync.RWMutex{},
		tenantResponses:               make(map[string]map[string]interface{}),
		tenantResponsesMutex:          sync.RWMutex{},
		gatewayProfileResponses:       make(map[string]map[string]interface{}),
		gatewayProfileResponsesMutex:  sync.RWMutex{},
		ethPortProfileResponses:       make(map[string]map[string]interface{}),
		ethPortProfileResponsesMutex:  sync.RWMutex{},
		ethPortSettingsResponses:      make(map[string]map[string]interface{}),
		ethPortSettingsResponsesMutex: sync.RWMutex{},
		bundleResponses:               make(map[string]map[string]interface{}),
		bundleResponsesMutex:          sync.RWMutex{},
	}
}

func (b *BulkOperationManager) GetResourceResponse(resourceType, resourceName string) (map[string]interface{}, bool) {
	switch resourceType {
	case "gateway":
		b.gatewayResponsesMutex.RLock()
		defer b.gatewayResponsesMutex.RUnlock()
		response, exists := b.gatewayResponses[resourceName]
		return response, exists

	case "lag":
		b.lagResponsesMutex.RLock()
		defer b.lagResponsesMutex.RUnlock()
		response, exists := b.lagResponses[resourceName]
		return response, exists

	case "service":
		b.serviceResponsesMutex.RLock()
		defer b.serviceResponsesMutex.RUnlock()
		response, exists := b.serviceResponses[resourceName]
		return response, exists

	case "tenant":
		b.tenantResponsesMutex.RLock()
		defer b.tenantResponsesMutex.RUnlock()
		response, exists := b.tenantResponses[resourceName]
		return response, exists

	case "gateway_profile":
		b.gatewayProfileResponsesMutex.RLock()
		defer b.gatewayProfileResponsesMutex.RUnlock()
		response, exists := b.gatewayProfileResponses[resourceName]
		return response, exists

	case "eth_port_profile":
		b.ethPortProfileResponsesMutex.RLock()
		defer b.ethPortProfileResponsesMutex.RUnlock()
		response, exists := b.ethPortProfileResponses[resourceName]
		return response, exists

	case "eth_port_settings":
		b.ethPortSettingsResponsesMutex.RLock()
		defer b.ethPortSettingsResponsesMutex.RUnlock()
		response, exists := b.ethPortSettingsResponses[resourceName]
		return response, exists

	case "bundle":
		b.bundleResponsesMutex.RLock()
		defer b.bundleResponsesMutex.RUnlock()
		response, exists := b.bundleResponses[resourceName]
		return response, exists

	default:
		return nil, false
	}
}

func GetBulkOperationManager(client *openapi.APIClient, clearCacheFunc ClearCacheFunc, providerContext interface{}) *BulkOperationManager {
	contextProvider := func() interface{} {
		return providerContext
	}

	return NewBulkOperationManager(client, contextProvider, clearCacheFunc)
}

func (b *BulkOperationManager) FailAllPendingOperations(ctx context.Context, err error) {
	b.operationMutex.Lock()
	var idsToClose []string
	failCount := 0

	for opID, op := range b.pendingOperations {
		if op.Status == OperationPending {
			updatedOp := op
			updatedOp.Status = OperationFailed
			updatedOp.Error = fmt.Errorf("Operation aborted due to previous failure: %v", err)
			b.pendingOperations[opID] = updatedOp
			b.operationErrors[opID] = updatedOp.Error
			b.operationResults[opID] = false
			idsToClose = append(idsToClose, opID)
			failCount++
		}
	}
	b.operationMutex.Unlock()

	for _, opID := range idsToClose {
		b.safeCloseChannel(opID)
	}

	if failCount > 0 {
		tflog.Error(ctx, fmt.Sprintf("Failed %d pending operations due to a previous operation failure", failCount), map[string]interface{}{
			"error": err.Error(),
		})
	}
}

func (b *BulkOperationManager) ExecuteAllPendingOperations(ctx context.Context) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	if time.Since(b.lastOperationTime) < BatchCollectionWindow {
		remaining := BatchCollectionWindow - time.Since(b.lastOperationTime)
		tflog.Debug(ctx, fmt.Sprintf("Waiting %v to collect more operations before executing", remaining))
		time.Sleep(remaining)
	}

	// Get counts of each operation type
	tenantPutCount := len(b.tenantPut)
	tenantPatchCount := len(b.tenantPatch)
	tenantDeleteCount := len(b.tenantDelete)

	gatewayPutCount := len(b.gatewayPut)
	gatewayPatchCount := len(b.gatewayPatch)
	gatewayDeleteCount := len(b.gatewayDelete)

	gatewayProfilePutCount := len(b.gatewayProfilePut)
	gatewayProfilePatchCount := len(b.gatewayProfilePatch)
	gatewayProfileDeleteCount := len(b.gatewayProfileDelete)

	servicePutCount := len(b.servicePut)
	servicePatchCount := len(b.servicePatch)
	serviceDeleteCount := len(b.serviceDelete)

	ethPortProfilePutCount := len(b.ethPortProfilePut)
	ethPortProfilePatchCount := len(b.ethPortProfilePatch)
	ethPortProfileDeleteCount := len(b.ethPortProfileDelete)

	ethPortSettingsPutCount := len(b.ethPortSettingsPut)
	ethPortSettingsPatchCount := len(b.ethPortSettingsPatch)
	ethPortSettingsDeleteCount := len(b.ethPortSettingsDelete)

	lagPutCount := len(b.lagPut)
	lagPatchCount := len(b.lagPatch)
	lagDeleteCount := len(b.lagDelete)

	bundlePatchCount := len(b.bundlePatch)

	operationsPerformed := false

	// Step 1: Execute all PUT operations in the desired order
	// Order:
	// 1. Tenants
	// 2. Gateways
	// 3. Gateway Profiles
	// 4. Services
	// 5. Eth Port Profiles
	// 6. Eth Port Settings
	// 7. Lags

	// First execute all PUT operations
	if tenantPutCount > 0 {
		tflog.Debug(ctx, "Executing Tenant PUT operations", map[string]interface{}{
			"operation_count": tenantPutCount,
		})
		putDiags := b.ExecuteBulkTenantPut(ctx)
		diagnostics.Append(putDiags...)

		if putDiags.HasError() {
			err := fmt.Errorf("bulk tenant PUT operation failed")
			b.FailAllPendingOperations(ctx, err)
			return diagnostics
		}
		operationsPerformed = true
	}

	if gatewayPutCount > 0 {
		tflog.Debug(ctx, "Executing Gateway PUT operations", map[string]interface{}{
			"operation_count": gatewayPutCount,
		})
		putDiags := b.ExecuteBulkGatewayPut(ctx)
		diagnostics.Append(putDiags...)

		if putDiags.HasError() {
			err := fmt.Errorf("bulk gateway PUT operation failed")
			b.FailAllPendingOperations(ctx, err)
			return diagnostics
		}
		operationsPerformed = true
	}

	if gatewayProfilePutCount > 0 {
		tflog.Debug(ctx, "Executing Gateway Profile PUT operations", map[string]interface{}{
			"operation_count": gatewayProfilePutCount,
		})
		putDiags := b.ExecuteBulkGatewayProfilePut(ctx)
		diagnostics.Append(putDiags...)

		if putDiags.HasError() {
			err := fmt.Errorf("bulk gateway profile PUT operation failed")
			b.FailAllPendingOperations(ctx, err)
			return diagnostics
		}
		operationsPerformed = true
	}

	if servicePutCount > 0 {
		tflog.Debug(ctx, "Executing Service PUT operations", map[string]interface{}{
			"operation_count": servicePutCount,
		})
		putDiags := b.ExecuteBulkServicePut(ctx)
		diagnostics.Append(putDiags...)

		if putDiags.HasError() {
			err := fmt.Errorf("bulk service PUT operation failed")
			b.FailAllPendingOperations(ctx, err)
			return diagnostics
		}
		operationsPerformed = true
	}

	if ethPortProfilePutCount > 0 {
		tflog.Debug(ctx, "Executing Eth Port Profile PUT operations", map[string]interface{}{
			"operation_count": ethPortProfilePutCount,
		})
		putDiags := b.ExecuteBulkEthPortProfilePut(ctx)
		diagnostics.Append(putDiags...)

		if putDiags.HasError() {
			err := fmt.Errorf("bulk eth port profile PUT operation failed")
			b.FailAllPendingOperations(ctx, err)
			return diagnostics
		}
		operationsPerformed = true
	}

	if ethPortSettingsPutCount > 0 {
		tflog.Debug(ctx, "Executing Eth Port Settings PUT operations", map[string]interface{}{
			"operation_count": ethPortSettingsPutCount,
		})
		putDiags := b.ExecuteBulkEthPortSettingsPut(ctx)
		diagnostics.Append(putDiags...)

		if putDiags.HasError() {
			err := fmt.Errorf("bulk eth port settings PUT operation failed")
			b.FailAllPendingOperations(ctx, err)
			return diagnostics
		}
		operationsPerformed = true
	}

	if lagPutCount > 0 {
		tflog.Debug(ctx, "Executing LAG PUT operations", map[string]interface{}{
			"operation_count": lagPutCount,
		})
		putDiags := b.ExecuteBulkLagPut(ctx)
		diagnostics.Append(putDiags...)

		if putDiags.HasError() {
			err := fmt.Errorf("bulk LAG PUT operation failed")
			b.FailAllPendingOperations(ctx, err)
			return diagnostics
		}
		operationsPerformed = true
	}

	// Step 2: Execute all PATCH operations in the desired order
	// Order:
	// 1. Tenants
	// 2. Gateways
	// 3. Gateway Profiles
	// 4. Services
	// 5. Eth Port Profiles
	// 6. Eth Port Settings
	// 7. Lags
	// 8. Bundles

	if tenantPatchCount > 0 {
		tflog.Debug(ctx, "Executing Tenant PATCH operations", map[string]interface{}{
			"operation_count": tenantPatchCount,
		})
		patchDiags := b.ExecuteBulkTenantPatch(ctx)
		diagnostics.Append(patchDiags...)

		if patchDiags.HasError() {
			err := fmt.Errorf("bulk tenant PATCH operation failed")
			b.FailAllPendingOperations(ctx, err)
			return diagnostics
		}
		operationsPerformed = true
	}

	if gatewayPatchCount > 0 {
		tflog.Debug(ctx, "Executing Gateway PATCH operations", map[string]interface{}{
			"operation_count": gatewayPatchCount,
		})
		patchDiags := b.ExecuteBulkGatewayPatch(ctx)
		diagnostics.Append(patchDiags...)

		if patchDiags.HasError() {
			err := fmt.Errorf("bulk gateway PATCH operation failed")
			b.FailAllPendingOperations(ctx, err)
			return diagnostics
		}
		operationsPerformed = true
	}

	if gatewayProfilePatchCount > 0 {
		tflog.Debug(ctx, "Executing Gateway Profile PATCH operations", map[string]interface{}{
			"operation_count": gatewayProfilePatchCount,
		})
		patchDiags := b.ExecuteBulkGatewayProfilePatch(ctx)
		diagnostics.Append(patchDiags...)

		if patchDiags.HasError() {
			err := fmt.Errorf("bulk gateway profile PATCH operation failed")
			b.FailAllPendingOperations(ctx, err)
			return diagnostics
		}
		operationsPerformed = true
	}

	if servicePatchCount > 0 {
		tflog.Debug(ctx, "Executing Service PATCH operations", map[string]interface{}{
			"operation_count": servicePatchCount,
		})
		patchDiags := b.ExecuteBulkServicePatch(ctx)
		diagnostics.Append(patchDiags...)

		if patchDiags.HasError() {
			err := fmt.Errorf("bulk service PATCH operation failed")
			b.FailAllPendingOperations(ctx, err)
			return diagnostics
		}
		operationsPerformed = true
	}

	if ethPortProfilePatchCount > 0 {
		tflog.Debug(ctx, "Executing Eth Port Profile PATCH operations", map[string]interface{}{
			"operation_count": ethPortProfilePatchCount,
		})
		patchDiags := b.ExecuteBulkEthPortProfilePatch(ctx)
		diagnostics.Append(patchDiags...)

		if patchDiags.HasError() {
			err := fmt.Errorf("bulk eth port profile PATCH operation failed")
			b.FailAllPendingOperations(ctx, err)
			return diagnostics
		}
		operationsPerformed = true
	}

	if ethPortSettingsPatchCount > 0 {
		tflog.Debug(ctx, "Executing Eth Port Settings PATCH operations", map[string]interface{}{
			"operation_count": ethPortSettingsPatchCount,
		})
		patchDiags := b.ExecuteBulkEthPortSettingsPatch(ctx)
		diagnostics.Append(patchDiags...)

		if patchDiags.HasError() {
			err := fmt.Errorf("bulk eth port settings PATCH operation failed")
			b.FailAllPendingOperations(ctx, err)
			return diagnostics
		}
		operationsPerformed = true
	}

	if lagPatchCount > 0 {
		tflog.Debug(ctx, "Executing LAG PATCH operations", map[string]interface{}{
			"operation_count": lagPatchCount,
		})
		patchDiags := b.ExecuteBulkLagPatch(ctx)
		diagnostics.Append(patchDiags...)

		if patchDiags.HasError() {
			err := fmt.Errorf("bulk lag PATCH operation failed")
			b.FailAllPendingOperations(ctx, err)
			return diagnostics
		}
		operationsPerformed = true
	}

	if bundlePatchCount > 0 {
		tflog.Debug(ctx, "Executing Bundle PATCH operations", map[string]interface{}{
			"operation_count": bundlePatchCount,
		})
		patchDiags := b.ExecuteBulkBundlePatch(ctx)
		diagnostics.Append(patchDiags...)

		if patchDiags.HasError() {
			err := fmt.Errorf("bulk bundle PATCH operation failed")
			b.FailAllPendingOperations(ctx, err)
			return diagnostics
		}
		operationsPerformed = true
	}

	// Step 3: Execute all DELETE operations in desired order.
	// DELETE order:
	// 1. Bundles
	// 2. Lags
	// 3. Eth Port Settings
	// 4. Eth Port Profiles
	// 5. Services
	// 6. Gateway Profiles
	// 7. Gateways
	// 8. Tenants

	// Skipping Bundles DELETE operations as it's not supported by API

	if lagDeleteCount > 0 {
		tflog.Debug(ctx, "Executing LAG DELETE operations", map[string]interface{}{
			"operation_count": lagDeleteCount,
		})
		deleteDiags := b.ExecuteBulkLagDelete(ctx)
		diagnostics.Append(deleteDiags...)

		if deleteDiags.HasError() {
			err := fmt.Errorf("bulk lag DELETE operation failed")
			b.FailAllPendingOperations(ctx, err)
			return diagnostics
		}
		operationsPerformed = true
	}

	if ethPortSettingsDeleteCount > 0 {
		tflog.Debug(ctx, "Executing Eth Port Settings DELETE operations", map[string]interface{}{
			"operation_count": ethPortSettingsDeleteCount,
		})
		deleteDiags := b.ExecuteBulkEthPortSettingsDelete(ctx)
		diagnostics.Append(deleteDiags...)

		if deleteDiags.HasError() {
			err := fmt.Errorf("bulk eth port settings DELETE operation failed")
			b.FailAllPendingOperations(ctx, err)
			return diagnostics
		}
		operationsPerformed = true
	}

	if ethPortProfileDeleteCount > 0 {
		tflog.Debug(ctx, "Executing Eth Port Profile DELETE operations", map[string]interface{}{
			"operation_count": ethPortProfileDeleteCount,
		})
		deleteDiags := b.ExecuteBulkEthPortProfileDelete(ctx)
		diagnostics.Append(deleteDiags...)

		if deleteDiags.HasError() {
			err := fmt.Errorf("bulk eth port profile DELETE operation failed")
			b.FailAllPendingOperations(ctx, err)
			return diagnostics
		}
		operationsPerformed = true
	}

	if serviceDeleteCount > 0 {
		tflog.Debug(ctx, "Executing Service DELETE operations", map[string]interface{}{
			"operation_count": serviceDeleteCount,
		})
		deleteDiags := b.ExecuteBulkServiceDelete(ctx)
		diagnostics.Append(deleteDiags...)

		if deleteDiags.HasError() {
			err := fmt.Errorf("bulk service DELETE operation failed")
			b.FailAllPendingOperations(ctx, err)
			return diagnostics
		}
		operationsPerformed = true
	}

	if gatewayProfileDeleteCount > 0 {
		tflog.Debug(ctx, "Executing Gateway Profile DELETE operations", map[string]interface{}{
			"operation_count": gatewayProfileDeleteCount,
		})
		deleteDiags := b.ExecuteBulkGatewayProfileDelete(ctx)
		diagnostics.Append(deleteDiags...)

		if deleteDiags.HasError() {
			err := fmt.Errorf("bulk gateway profile DELETE operation failed")
			b.FailAllPendingOperations(ctx, err)
			return diagnostics
		}
		operationsPerformed = true
	}

	if gatewayDeleteCount > 0 {
		tflog.Debug(ctx, "Executing Gateway DELETE operations", map[string]interface{}{
			"operation_count": gatewayDeleteCount,
		})
		deleteDiags := b.ExecuteBulkGatewayDelete(ctx)
		diagnostics.Append(deleteDiags...)

		if deleteDiags.HasError() {
			err := fmt.Errorf("bulk gateway DELETE operation failed")
			b.FailAllPendingOperations(ctx, err)
			return diagnostics
		}
		operationsPerformed = true
	}

	if tenantDeleteCount > 0 {
		tflog.Debug(ctx, "Executing Tenant DELETE operations", map[string]interface{}{
			"operation_count": tenantDeleteCount,
		})
		deleteDiags := b.ExecuteBulkTenantDelete(ctx)
		diagnostics.Append(deleteDiags...)

		if deleteDiags.HasError() {
			err := fmt.Errorf("bulk tenant DELETE operation failed")
			b.FailAllPendingOperations(ctx, err)
			return diagnostics
		}
		operationsPerformed = true
	}

	if operationsPerformed {
		waitDuration := 800 * time.Millisecond
		tflog.Debug(ctx, fmt.Sprintf("Waiting %v for all operations to propagate before final cache refresh", waitDuration))
		time.Sleep(waitDuration)

		tflog.Debug(ctx, "Final cache clear after all operations to ensure verification with fresh data")
		if b.clearCacheFunc != nil && b.contextProvider != nil {
			b.clearCacheFunc(ctx, b.contextProvider(), "gateways")
			b.clearCacheFunc(ctx, b.contextProvider(), "lags")
			b.clearCacheFunc(ctx, b.contextProvider(), "tenants")
			b.clearCacheFunc(ctx, b.contextProvider(), "services")
			b.clearCacheFunc(ctx, b.contextProvider(), "gatewayprofiles")
			b.clearCacheFunc(ctx, b.contextProvider(), "ethportprofiles")
			b.clearCacheFunc(ctx, b.contextProvider(), "ethportsettings")
			b.clearCacheFunc(ctx, b.contextProvider(), "bundles")
		}
	}

	return diagnostics
}

func (b *BulkOperationManager) ShouldExecuteOperations(ctx context.Context) bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// If there are no pending operations, no need to execute
	if len(b.gatewayPut) == 0 && len(b.gatewayPatch) == 0 && len(b.gatewayDelete) == 0 &&
		len(b.lagPut) == 0 && len(b.lagPatch) == 0 && len(b.lagDelete) == 0 &&
		len(b.tenantPut) == 0 && len(b.tenantPatch) == 0 && len(b.tenantDelete) == 0 &&
		len(b.servicePut) == 0 && len(b.servicePatch) == 0 && len(b.serviceDelete) == 0 &&
		len(b.gatewayProfilePut) == 0 && len(b.gatewayProfilePatch) == 0 && len(b.gatewayProfileDelete) == 0 &&
		len(b.ethPortProfilePut) == 0 && len(b.ethPortProfilePatch) == 0 && len(b.ethPortProfileDelete) == 0 &&
		len(b.ethPortSettingsPut) == 0 && len(b.ethPortSettingsPatch) == 0 && len(b.ethPortSettingsDelete) == 0 &&
		len(b.bundlePatch) == 0 {
		return false
	}

	elapsedSinceLast := time.Since(b.lastOperationTime)
	elapsedSinceBatchStart := time.Since(b.batchStartTime)

	// Only flush if either sufficient time has passed since the last operation
	// OR the batch has been open for too long
	if elapsedSinceLast < BatchCollectionWindow && elapsedSinceBatchStart < MaxBatchDelay {
		return false
	}

	return true
}

func (b *BulkOperationManager) ExecuteIfMultipleOperations(ctx context.Context) diag.Diagnostics {
	b.mutex.Lock()
	gatewayPutCount := len(b.gatewayPut)
	gatewayPatchCount := len(b.gatewayPatch)
	gatewayDeleteCount := len(b.gatewayDelete)

	lagPutCount := len(b.lagPut)
	lagPatchCount := len(b.lagPatch)
	lagDeleteCount := len(b.lagDelete)

	tenantPutCount := len(b.tenantPut)
	tenantPatchCount := len(b.tenantPatch)
	tenantDeleteCount := len(b.tenantDelete)

	servicePutCount := len(b.servicePut)
	servicePatchCount := len(b.servicePatch)
	serviceDeleteCount := len(b.serviceDelete)

	gatewayProfilePutCount := len(b.gatewayProfilePut)
	gatewayProfilePatchCount := len(b.gatewayProfilePatch)
	gatewayProfileDeleteCount := len(b.gatewayProfileDelete)

	ethPortProfilePutCount := len(b.ethPortProfilePut)
	ethPortProfilePatchCount := len(b.ethPortProfilePatch)
	ethPortProfileDeleteCount := len(b.ethPortProfileDelete)

	ethPortSettingsPutCount := len(b.ethPortSettingsPut)
	ethPortSettingsPatchCount := len(b.ethPortSettingsPatch)
	ethPortSettingsDeleteCount := len(b.ethPortSettingsDelete)

	bundlePatchCount := len(b.bundlePatch)

	b.mutex.Unlock()

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

		return b.ExecuteAllPendingOperations(ctx)
	}

	return nil
}

func (b *BulkOperationManager) hasPendingOrRecentOperations(
	resourceType string,
) bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	var putLen, patchLen, deleteLen int
	var recentOps bool
	var recentOpTime time.Time

	switch resourceType {
	case "gateway":
		putLen = len(b.gatewayPut)
		patchLen = len(b.gatewayPatch)
		deleteLen = len(b.gatewayDelete)
		recentOps = b.recentGatewayOps
		recentOpTime = b.recentGatewayOpTime
	case "lag":
		putLen = len(b.lagPut)
		patchLen = len(b.lagPatch)
		deleteLen = len(b.lagDelete)
		recentOps = b.recentLagOps
		recentOpTime = b.recentLagOpTime
	case "service":
		putLen = len(b.servicePut)
		patchLen = len(b.servicePatch)
		deleteLen = len(b.serviceDelete)
		recentOps = b.recentServiceOps
		recentOpTime = b.recentServiceOpTime
	case "tenant":
		putLen = len(b.tenantPut)
		patchLen = len(b.tenantPatch)
		deleteLen = len(b.tenantDelete)
		recentOps = b.recentTenantOps
		recentOpTime = b.recentTenantOpTime
	case "gateway_profile":
		putLen = len(b.gatewayProfilePut)
		patchLen = len(b.gatewayProfilePatch)
		deleteLen = len(b.gatewayProfileDelete)
		recentOps = b.recentGatewayProfileOps
		recentOpTime = b.recentGatewayProfileOpTime
	case "eth_port_profile":
		putLen = len(b.ethPortProfilePut)
		patchLen = len(b.ethPortProfilePatch)
		deleteLen = len(b.ethPortProfileDelete)
		recentOps = b.recentEthPortProfileOps
		recentOpTime = b.recentEthPortProfileOpTime
	case "eth_port_settings":
		putLen = len(b.ethPortSettingsPut)
		patchLen = len(b.ethPortSettingsPatch)
		deleteLen = len(b.ethPortSettingsDelete)
		recentOps = b.recentEthPortSettingsOps
		recentOpTime = b.recentEthPortSettingsOpTime
	case "bundle":
		putLen = 0
		patchLen = len(b.bundlePatch)
		deleteLen = 0
		recentOps = b.recentBundleOps
		recentOpTime = b.recentBundleOpTime
	}

	// Check if any operations are pending
	hasPending := putLen > 0 || patchLen > 0 || deleteLen > 0

	// Check if we've recently had operations (within the last 5 seconds)
	hasRecent := recentOps && time.Since(recentOpTime) < 5*time.Second

	return hasPending || hasRecent
}

func (b *BulkOperationManager) HasPendingOrRecentGatewayOperations() bool {
	return b.hasPendingOrRecentOperations("gateway")
}

func (b *BulkOperationManager) HasPendingOrRecentLagOperations() bool {
	return b.hasPendingOrRecentOperations("lag")
}

func (b *BulkOperationManager) HasPendingOrRecentServiceOperations() bool {
	return b.hasPendingOrRecentOperations("service")
}

func (b *BulkOperationManager) HasPendingOrRecentTenantOperations() bool {
	return b.hasPendingOrRecentOperations("tenant")
}

func (b *BulkOperationManager) HasPendingOrRecentGatewayProfileOperations() bool {
	return b.hasPendingOrRecentOperations("gateway_profile")
}

func (b *BulkOperationManager) HasPendingOrRecentEthPortProfileOperations() bool {
	return b.hasPendingOrRecentOperations("eth_port_profile")
}

func (b *BulkOperationManager) HasPendingOrRecentEthPortSettingsOperations() bool {
	return b.hasPendingOrRecentOperations("eth_port_settings")
}

func (b *BulkOperationManager) HasPendingOrRecentBundleOperations() bool {
	return b.hasPendingOrRecentOperations("bundle")
}

func (b *BulkOperationManager) addOperation(
	ctx context.Context,
	resourceType string,
	resourceName string,
	operationType string,
	storeFunc func(),
	logDetails map[string]interface{},
) string {
	storeFunc()

	operationID := generateOperationID(resourceType, resourceName, operationType)
	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	b.pendingOperations[operationID] = &Operation{
		ResourceType:  resourceType,
		ResourceName:  resourceName,
		OperationType: operationType,
		Status:        OperationPending,
	}

	b.operationWaitChannels[operationID] = make(chan struct{})

	now := time.Now()
	b.lastOperationTime = now
	if b.batchStartTime.IsZero() {
		b.batchStartTime = now
	}

	if logDetails != nil {
		logDetails["operation_id"] = operationID
		tflog.Debug(ctx, fmt.Sprintf("Added %s to %s batch", resourceType, operationType), logDetails)
	}

	return operationID
}

func (b *BulkOperationManager) AddBundlePatch(ctx context.Context, bundleName string, props openapi.BundlesPatchRequestEndpointBundleValue) string {
	return b.addOperation(
		ctx,
		"bundle",
		bundleName,
		"PATCH",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.bundlePatch[bundleName] = props
		},
		map[string]interface{}{
			"bundle_name": bundleName,
			"batch_size":  len(b.bundlePatch) + 1,
		},
	)
}

func (b *BulkOperationManager) AddGatewayPut(ctx context.Context, gatewayName string, props openapi.ConfigPutRequestGatewayGatewayName) string {
	return b.addOperation(
		ctx,
		"gateway",
		gatewayName,
		"PUT",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.gatewayPut[gatewayName] = props
		},
		map[string]interface{}{
			"gateway_name": gatewayName,
			"batch_size":   len(b.gatewayPatch) + 1,
		},
	)
}

func (b *BulkOperationManager) AddGatewayPatch(ctx context.Context, gatewayName string, props openapi.ConfigPutRequestGatewayGatewayName) string {
	return b.addOperation(
		ctx,
		"gateway",
		gatewayName,
		"PATCH",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.gatewayPatch[gatewayName] = props
		},
		map[string]interface{}{
			"gateway_name": gatewayName,
			"batch_size":   len(b.gatewayPatch) + 1,
		},
	)
}

func (b *BulkOperationManager) AddGatewayDelete(ctx context.Context, gatewayName string) string {
	return b.addOperation(
		ctx,
		"gateway",
		gatewayName,
		"DELETE",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.gatewayDelete = append(b.gatewayDelete, gatewayName)
		},
		map[string]interface{}{
			"gateway_name": gatewayName,
			"batch_size":   len(b.gatewayDelete) + 1,
		},
	)
}

func (b *BulkOperationManager) AddLagPut(ctx context.Context, lagName string, props openapi.ConfigPutRequestLagLagName) string {
	return b.addOperation(
		ctx,
		"lag",
		lagName,
		"PUT",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.lagPut[lagName] = props
		},
		map[string]interface{}{
			"lag_name":   lagName,
			"batch_size": len(b.lagPut) + 1,
		},
	)
}

func (b *BulkOperationManager) AddLagPatch(ctx context.Context, lagName string, props openapi.ConfigPutRequestLagLagName) string {
	return b.addOperation(
		ctx,
		"lag",
		lagName,
		"PATCH",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.lagPatch[lagName] = props
		},
		map[string]interface{}{
			"lag_name":   lagName,
			"batch_size": len(b.lagPatch) + 1,
		},
	)
}

func (b *BulkOperationManager) AddLagDelete(ctx context.Context, lagName string) string {
	return b.addOperation(
		ctx,
		"lag",
		lagName,
		"DELETE",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.lagDelete = append(b.lagDelete, lagName)
		},
		map[string]interface{}{
			"lag_name":   lagName,
			"batch_size": len(b.lagDelete) + 1,
		},
	)
}

func (b *BulkOperationManager) AddTenantPut(ctx context.Context, tenantName string, props openapi.ConfigPutRequestTenantTenantName) string {
	return b.addOperation(
		ctx,
		"tenant",
		tenantName,
		"PUT",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.tenantPut[tenantName] = props
		},
		map[string]interface{}{
			"tenant_name": tenantName,
			"batch_size":  len(b.tenantPut) + 1,
		},
	)
}

func (b *BulkOperationManager) AddTenantPatch(ctx context.Context, tenantName string, props openapi.ConfigPutRequestTenantTenantName) string {
	return b.addOperation(
		ctx,
		"tenant",
		tenantName,
		"PATCH",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.tenantPatch[tenantName] = props
		},
		map[string]interface{}{
			"tenant_name": tenantName,
			"batch_size":  len(b.tenantPatch) + 1,
		},
	)
}

func (b *BulkOperationManager) AddTenantDelete(ctx context.Context, tenantName string) string {
	return b.addOperation(
		ctx,
		"tenant",
		tenantName,
		"DELETE",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.tenantDelete = append(b.tenantDelete, tenantName)
		},
		map[string]interface{}{
			"tenant_name": tenantName,
			"batch_size":  len(b.tenantDelete) + 1,
		},
	)
}

func (b *BulkOperationManager) AddServicePut(ctx context.Context, serviceName string, props openapi.ConfigPutRequestServiceServiceName) string {
	return b.addOperation(
		ctx,
		"service",
		serviceName,
		"PUT",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.servicePut[serviceName] = props
		},
		map[string]interface{}{
			"service_name": serviceName,
			"batch_size":   len(b.servicePut) + 1,
		},
	)
}

func (b *BulkOperationManager) AddServicePatch(ctx context.Context, serviceName string, props openapi.ConfigPutRequestServiceServiceName) string {
	return b.addOperation(
		ctx,
		"service",
		serviceName,
		"PATCH",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.servicePatch[serviceName] = props
		},
		map[string]interface{}{
			"service_name": serviceName,
			"batch_size":   len(b.servicePatch) + 1,
		},
	)
}

func (b *BulkOperationManager) AddServiceDelete(ctx context.Context, serviceName string) string {
	return b.addOperation(
		ctx,
		"service",
		serviceName,
		"DELETE",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.serviceDelete = append(b.serviceDelete, serviceName)
		},
		map[string]interface{}{
			"service_name": serviceName,
			"batch_size":   len(b.serviceDelete) + 1,
		},
	)
}

func (b *BulkOperationManager) AddGatewayProfilePut(ctx context.Context, profileName string, props openapi.ConfigPutRequestGatewayProfileGatewayProfileName) string {
	return b.addOperation(
		ctx,
		"gateway_profile",
		profileName,
		"PUT",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.gatewayProfilePut[profileName] = props
		},
		map[string]interface{}{
			"profile_name": profileName,
			"batch_size":   len(b.gatewayProfilePut) + 1,
		},
	)
}

func (b *BulkOperationManager) AddGatewayProfilePatch(ctx context.Context, profileName string, props openapi.ConfigPutRequestGatewayProfileGatewayProfileName) string {
	return b.addOperation(
		ctx,
		"gateway_profile",
		profileName,
		"PATCH",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.gatewayProfilePatch[profileName] = props
		},
		map[string]interface{}{
			"profile_name": profileName,
			"batch_size":   len(b.gatewayProfilePatch) + 1,
		},
	)
}

func (b *BulkOperationManager) AddGatewayProfileDelete(ctx context.Context, profileName string) string {
	return b.addOperation(
		ctx,
		"gateway_profile",
		profileName,
		"DELETE",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.gatewayProfileDelete = append(b.gatewayProfileDelete, profileName)
		},
		map[string]interface{}{
			"profile_name": profileName,
			"batch_size":   len(b.gatewayProfileDelete) + 1,
		},
	)
}

func (b *BulkOperationManager) AddEthPortProfilePut(ctx context.Context, ethPortProfileName string, props openapi.ConfigPutRequestEthPortProfileEthPortProfileName) string {
	return b.addOperation(
		ctx,
		"eth_port_profile",
		ethPortProfileName,
		"PUT",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.ethPortProfilePut[ethPortProfileName] = props
		},
		map[string]interface{}{
			"eth_port_profile_name": ethPortProfileName,
			"batch_size":            len(b.ethPortProfilePut) + 1,
		},
	)
}

func (b *BulkOperationManager) AddEthPortProfilePatch(ctx context.Context, ethPortProfileName string, props openapi.ConfigPutRequestEthPortProfileEthPortProfileName) string {
	return b.addOperation(
		ctx,
		"eth_port_profile",
		ethPortProfileName,
		"PATCH",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.ethPortProfilePatch[ethPortProfileName] = props
		},
		map[string]interface{}{
			"eth_port_profile_name": ethPortProfileName,
			"batch_size":            len(b.ethPortProfilePatch) + 1,
		},
	)
}

func (b *BulkOperationManager) AddEthPortProfileDelete(ctx context.Context, ethPortProfileName string) string {
	return b.addOperation(
		ctx,
		"eth_port_profile",
		ethPortProfileName,
		"DELETE",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.ethPortProfileDelete = append(b.ethPortProfileDelete, ethPortProfileName)
		},
		map[string]interface{}{
			"eth_port_profile_name": ethPortProfileName,
			"batch_size":            len(b.ethPortProfileDelete) + 1,
		},
	)
}

func (b *BulkOperationManager) AddEthPortSettingsPut(ctx context.Context, ethPortSettingsName string, props openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName) string {
	return b.addOperation(
		ctx,
		"eth_port_settings",
		ethPortSettingsName,
		"PUT",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.ethPortSettingsPut[ethPortSettingsName] = props
		},
		map[string]interface{}{
			"eth_port_settings_name": ethPortSettingsName,
			"batch_size":             len(b.ethPortSettingsPut) + 1,
		},
	)
}

func (b *BulkOperationManager) AddEthPortSettingsPatch(ctx context.Context, ethPortSettingsName string, props openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName) string {
	return b.addOperation(
		ctx,
		"eth_port_settings",
		ethPortSettingsName,
		"PATCH",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.ethPortSettingsPatch[ethPortSettingsName] = props
		},
		map[string]interface{}{
			"eth_port_settings_name": ethPortSettingsName,
			"batch_size":             len(b.ethPortSettingsPatch) + 1,
		},
	)
}

func (b *BulkOperationManager) AddEthPortSettingsDelete(ctx context.Context, ethPortSettingsName string) string {
	return b.addOperation(
		ctx,
		"eth_port_settings",
		ethPortSettingsName,
		"DELETE",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.ethPortSettingsDelete = append(b.ethPortSettingsDelete, ethPortSettingsName)
		},
		map[string]interface{}{
			"eth_port_settings_name": ethPortSettingsName,
			"batch_size":             len(b.ethPortSettingsDelete) + 1,
		},
	)
}

func (b *BulkOperationManager) executeBulkOperation(ctx context.Context, config BulkOperationConfig) diag.Diagnostics {
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
		filteredResourceNames, filteredOperations, err = config.CheckPreExistence(ctx, resourceNames)
		if err != nil {
			tflog.Warn(ctx, fmt.Sprintf("Error checking for existing %s: %v - proceeding with all resources",
				config.ResourceType, err))
			filteredResourceNames = resourceNames
			filteredOperations = operations
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

	retryConfig := DefaultRetryConfig()
	var opErr error
	var apiResp *http.Response

	for retry := 0; retry < retryConfig.MaxRetries; retry++ {
		if retry > 0 {
			delay := CalculateBackoff(retry-1, retryConfig)
			tflog.Debug(ctx, fmt.Sprintf("Retrying bulk %s %s operation after delay",
				config.ResourceType, config.OperationType),
				map[string]interface{}{
					"retry": retry,
					"delay": delay,
				})
			time.Sleep(delay)
		}

		apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
		apiResp, opErr = config.ExecuteRequest(apiCtx, request)
		cancel()

		if opErr == nil {
			tflog.Debug(ctx, fmt.Sprintf("Bulk %s %s operation succeeded",
				config.ResourceType, config.OperationType),
				map[string]interface{}{
					"attempt": retry + 1,
				})
			break
		}

		if !IsRetriableError(opErr) {
			tflog.Error(ctx, fmt.Sprintf("Bulk %s %s operation failed with non-retriable error",
				config.ResourceType, config.OperationType),
				map[string]interface{}{
					"error": opErr.Error(),
				})
			break
		}

		delayTime := CalculateBackoff(retry, retryConfig)
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
		if err := config.ProcessResponse(ctx, apiResp); err != nil {
			tflog.Error(ctx, fmt.Sprintf("Error processing %s response: %v", config.ResourceType, err))
		}
	}

	b.updateOperationStatuses(ctx, config.ResourceType, config.OperationType, filteredResourceNames, opErr)

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

func (b *BulkOperationManager) updateOperationStatuses(ctx context.Context, resourceType, operationType string,
	resourceNames []string, opErr error) {
	resourceMap := make(map[string]bool)

	for _, name := range resourceNames {
		resourceMap[name] = true
	}

	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	for opID, op := range b.pendingOperations {
		if op.ResourceType == resourceType && op.OperationType == operationType {
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
					b.pendingOperations[opID] = updatedOp
					b.operationResults[opID] = true
				} else {
					// Mark operation as failed
					updatedOp.Status = OperationFailed
					updatedOp.Error = opErr
					b.pendingOperations[opID] = updatedOp
					b.operationErrors[opID] = opErr
					b.operationResults[opID] = false
				}
				b.safeCloseChannel(opID, true)
			}
		}
	}
}

func (b *BulkOperationManager) ExecuteBulkGatewayPut(ctx context.Context) diag.Diagnostics {
	var originalOperations map[string]openapi.ConfigPutRequestGatewayGatewayName

	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "gateway",
		OperationType: "PUT",

		ExtractOperations: func() (map[string]interface{}, []string) {
			b.mutex.Lock()
			originalOperations = make(map[string]openapi.ConfigPutRequestGatewayGatewayName)
			for k, v := range b.gatewayPut {
				originalOperations[k] = v
			}
			b.gatewayPut = make(map[string]openapi.ConfigPutRequestGatewayGatewayName)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			names := make([]string, 0, len(originalOperations))

			for k, v := range originalOperations {
				result[k] = v
				names = append(names, k)
			}

			return result, names
		},

		CheckPreExistence: func(ctx context.Context, resourceNames []string) ([]string, map[string]interface{}, error) {
			checker := ResourceExistenceCheck{
				ResourceType:  "gateway",
				OperationType: "PUT",
				FetchResources: func(ctx context.Context) (map[string]interface{}, error) {
					apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
					defer cancel()

					resp, err := b.client.GatewaysAPI.GatewaysGet(apiCtx).Execute()
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
				},
			}

			filteredNames, err := b.FilterPreExistingResources(ctx, resourceNames, checker)
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
			putRequest := openapi.NewGatewaysPutRequest()
			gatewayMap := make(map[string]openapi.ConfigPutRequestGatewayGatewayName)

			for name, props := range filteredData {
				gatewayMap[name] = props.(openapi.ConfigPutRequestGatewayGatewayName)
			}
			putRequest.SetGateway(gatewayMap)
			return putRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.GatewaysAPI.GatewaysPut(ctx).GatewaysPutRequest(
				*request.(*openapi.GatewaysPutRequest))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentGatewayOps = true
			b.recentGatewayOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkGatewayPatch(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "gateway",
		OperationType: "PATCH",

		ExtractOperations: func() (map[string]interface{}, []string) {

			b.mutex.Lock()
			gatewayPatch := make(map[string]openapi.ConfigPutRequestGatewayGatewayName)
			for k, v := range b.gatewayPatch {
				gatewayPatch[k] = v
			}
			b.gatewayPatch = make(map[string]openapi.ConfigPutRequestGatewayGatewayName)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			names := make([]string, 0, len(gatewayPatch))

			for k, v := range gatewayPatch {
				result[k] = v
				names = append(names, k)
			}

			return result, names
		},

		CheckPreExistence: nil,

		PrepareRequest: func(filteredData map[string]interface{}) interface{} {
			patchRequest := openapi.NewGatewaysPutRequest()
			gatewayMap := make(map[string]openapi.ConfigPutRequestGatewayGatewayName)

			for name, props := range filteredData {
				gatewayMap[name] = props.(openapi.ConfigPutRequestGatewayGatewayName)
			}
			patchRequest.SetGateway(gatewayMap)
			return patchRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.GatewaysAPI.GatewaysPatch(ctx).GatewaysPutRequest(
				*request.(*openapi.GatewaysPutRequest))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentGatewayOps = true
			b.recentGatewayOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkLagPut(ctx context.Context) diag.Diagnostics {
	var originalOperations map[string]openapi.ConfigPutRequestLagLagName

	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "lag",
		OperationType: "PUT",

		ExtractOperations: func() (map[string]interface{}, []string) {

			b.mutex.Lock()
			originalOperations = make(map[string]openapi.ConfigPutRequestLagLagName)
			for k, v := range b.lagPut {
				originalOperations[k] = v
			}
			b.lagPut = make(map[string]openapi.ConfigPutRequestLagLagName)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			names := make([]string, 0, len(originalOperations))

			for k, v := range originalOperations {
				result[k] = v
				names = append(names, k)
			}

			return result, names
		},

		CheckPreExistence: func(ctx context.Context, resourceNames []string) ([]string, map[string]interface{}, error) {
			checker := ResourceExistenceCheck{
				ResourceType:  "lag",
				OperationType: "PUT",
				FetchResources: func(ctx context.Context) (map[string]interface{}, error) {
					// First check if we have cached LAG data
					b.lagResponsesMutex.RLock()
					if len(b.lagResponses) > 0 {
						cachedData := make(map[string]interface{})
						for k, v := range b.lagResponses {
							cachedData[k] = v
						}
						b.lagResponsesMutex.RUnlock()

						tflog.Debug(ctx, "Using cached LAG data for pre-existence check", map[string]interface{}{
							"count": len(cachedData),
						})

						return cachedData, nil
					}
					b.lagResponsesMutex.RUnlock()

					// Fall back to API call if no cache
					apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
					defer cancel()

					resp, err := b.client.LAGsAPI.LagsGet(apiCtx).Execute()
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

					b.lagResponsesMutex.Lock()
					for k, v := range result.Lag {
						if vMap, ok := v.(map[string]interface{}); ok {
							b.lagResponses[k] = vMap

							if name, ok := vMap["name"].(string); ok && name != k {
								b.lagResponses[name] = vMap
							}
						}
					}
					b.lagResponsesMutex.Unlock()

					return result.Lag, nil
				},
			}

			filteredNames, err := b.FilterPreExistingResources(ctx, resourceNames, checker)
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
			putRequest := openapi.NewLagsPutRequest()
			lagMap := make(map[string]openapi.ConfigPutRequestLagLagName)

			for name, props := range filteredData {
				lagMap[name] = props.(openapi.ConfigPutRequestLagLagName)
			}
			putRequest.SetLag(lagMap)
			return putRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.LAGsAPI.LagsPut(ctx).LagsPutRequest(
				*request.(*openapi.LagsPutRequest))
			return req.Execute()
		},

		ProcessResponse: func(ctx context.Context, resp *http.Response) error {
			delayTime := 2 * time.Second
			tflog.Debug(ctx, fmt.Sprintf("Waiting %v for auto-generated values to be assigned before fetching LAGs", delayTime))
			time.Sleep(delayTime)

			fetchCtx, fetchCancel := context.WithTimeout(context.Background(), OperationTimeout)
			defer fetchCancel()

			tflog.Debug(ctx, "Fetching LAGs after successful PUT operation to retrieve auto-generated values")
			lagsReq := b.client.LAGsAPI.LagsGet(fetchCtx)
			lagsResp, fetchErr := lagsReq.Execute()

			if fetchErr != nil {
				tflog.Error(ctx, "Failed to fetch LAGs after PUT for auto-generated fields", map[string]interface{}{
					"error": fetchErr.Error(),
				})
				return fetchErr
			}

			defer lagsResp.Body.Close()

			var lagsData struct {
				Lag map[string]map[string]interface{} `json:"lag"`
			}

			if respErr := json.NewDecoder(lagsResp.Body).Decode(&lagsData); respErr != nil {
				tflog.Error(ctx, "Failed to decode LAGs response for auto-generated fields", map[string]interface{}{
					"error": respErr.Error(),
				})
				return respErr
			}

			b.lagResponsesMutex.Lock()
			for lagName, lagData := range lagsData.Lag {
				b.lagResponses[lagName] = lagData

				if name, ok := lagData["name"].(string); ok && name != lagName {
					b.lagResponses[name] = lagData
				}
			}
			b.lagResponsesMutex.Unlock()

			tflog.Debug(ctx, "Successfully stored LAG data for auto-generated fields", map[string]interface{}{
				"lag_count": len(lagsData.Lag),
			})

			return nil
		},

		UpdateRecentOps: func() {
			b.recentLagOps = true
			b.recentLagOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkLagPatch(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "lag",
		OperationType: "PATCH",

		ExtractOperations: func() (map[string]interface{}, []string) {

			b.mutex.Lock()
			lagPatch := make(map[string]openapi.ConfigPutRequestLagLagName)
			for k, v := range b.lagPatch {
				lagPatch[k] = v
			}
			b.lagPatch = make(map[string]openapi.ConfigPutRequestLagLagName)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			names := make([]string, 0, len(lagPatch))

			for k, v := range lagPatch {
				result[k] = v
				names = append(names, k)
			}

			return result, names
		},

		CheckPreExistence: nil,

		PrepareRequest: func(filteredData map[string]interface{}) interface{} {
			patchRequest := openapi.NewLagsPutRequest()
			lagMap := make(map[string]openapi.ConfigPutRequestLagLagName)

			for name, props := range filteredData {
				lagMap[name] = props.(openapi.ConfigPutRequestLagLagName)
			}
			patchRequest.SetLag(lagMap)
			return patchRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.LAGsAPI.LagsPatch(ctx).LagsPutRequest(
				*request.(*openapi.LagsPutRequest))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentLagOps = true
			b.recentLagOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkLagDelete(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "lag",
		OperationType: "DELETE",

		ExtractOperations: func() (map[string]interface{}, []string) {

			b.mutex.Lock()
			lagNames := make([]string, len(b.lagDelete))
			copy(lagNames, b.lagDelete)

			lagDeleteMap := make(map[string]bool)
			for _, name := range lagNames {
				lagDeleteMap[name] = true
			}

			b.lagDelete = make([]string, 0)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			for _, name := range lagNames {
				result[name] = true
			}

			return result, lagNames
		},

		CheckPreExistence: nil,

		PrepareRequest: func(filteredData map[string]interface{}) interface{} {
			names := make([]string, 0, len(filteredData))
			for name := range filteredData {
				names = append(names, name)
			}
			return names
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.LAGsAPI.LagsDelete(ctx).LagName(request.([]string))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentLagOps = true
			b.recentLagOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkTenantPut(ctx context.Context) diag.Diagnostics {
	var originalOperations map[string]openapi.ConfigPutRequestTenantTenantName

	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "tenant",
		OperationType: "PUT",

		ExtractOperations: func() (map[string]interface{}, []string) {

			b.mutex.Lock()
			originalOperations = make(map[string]openapi.ConfigPutRequestTenantTenantName)
			for k, v := range b.tenantPut {
				originalOperations[k] = v
			}
			b.tenantPut = make(map[string]openapi.ConfigPutRequestTenantTenantName)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			names := make([]string, 0, len(originalOperations))

			for k, v := range originalOperations {
				result[k] = v
				names = append(names, k)
			}

			return result, names
		},

		CheckPreExistence: func(ctx context.Context, resourceNames []string) ([]string, map[string]interface{}, error) {
			checker := ResourceExistenceCheck{
				ResourceType:  "tenant",
				OperationType: "PUT",
				FetchResources: func(ctx context.Context) (map[string]interface{}, error) {
					// First check if we have cached tenant data
					b.tenantResponsesMutex.RLock()
					if len(b.tenantResponses) > 0 {
						cachedData := make(map[string]interface{})
						for k, v := range b.tenantResponses {
							cachedData[k] = v
						}
						b.tenantResponsesMutex.RUnlock()

						tflog.Debug(ctx, "Using cached tenant data for pre-existence check", map[string]interface{}{
							"count": len(cachedData),
						})

						return cachedData, nil
					}
					b.tenantResponsesMutex.RUnlock()

					// Fall back to API call if no cache
					apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
					defer cancel()

					resp, err := b.client.TenantsAPI.TenantsGet(apiCtx).Execute()
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

					b.tenantResponsesMutex.Lock()
					for k, v := range result.Tenant {
						if vMap, ok := v.(map[string]interface{}); ok {
							b.tenantResponses[k] = vMap

							if name, ok := vMap["name"].(string); ok && name != k {
								b.tenantResponses[name] = vMap
							}
						}
					}
					b.tenantResponsesMutex.Unlock()

					return result.Tenant, nil
				},
			}

			filteredNames, err := b.FilterPreExistingResources(ctx, resourceNames, checker)
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
			putRequest := openapi.NewTenantsPutRequest()
			tenantMap := make(map[string]openapi.ConfigPutRequestTenantTenantName)

			for name, props := range filteredData {
				tenantMap[name] = props.(openapi.ConfigPutRequestTenantTenantName)
			}
			putRequest.SetTenant(tenantMap)
			return putRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.TenantsAPI.TenantsPut(ctx).TenantsPutRequest(
				*request.(*openapi.TenantsPutRequest))
			return req.Execute()
		},

		ProcessResponse: func(ctx context.Context, resp *http.Response) error {
			delayTime := 2 * time.Second
			tflog.Debug(ctx, fmt.Sprintf("Waiting %v for auto-generated values to be assigned before fetching tenants", delayTime))
			time.Sleep(delayTime)

			fetchCtx, fetchCancel := context.WithTimeout(context.Background(), OperationTimeout)
			defer fetchCancel()

			tflog.Debug(ctx, "Fetching tenants after successful PUT operation to retrieve auto-generated values")
			tenantsReq := b.client.TenantsAPI.TenantsGet(fetchCtx)
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

			b.tenantResponsesMutex.Lock()
			for tenantName, tenantData := range tenantsData.Tenant {
				b.tenantResponses[tenantName] = tenantData

				if name, ok := tenantData["name"].(string); ok && name != tenantName {
					b.tenantResponses[name] = tenantData
				}
			}
			b.tenantResponsesMutex.Unlock()

			tflog.Debug(ctx, "Successfully stored tenant data for auto-generated fields", map[string]interface{}{
				"tenant_count": len(tenantsData.Tenant),
			})

			return nil
		},

		UpdateRecentOps: func() {
			b.recentTenantOps = true
			b.recentTenantOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkTenantPatch(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "tenant",
		OperationType: "PATCH",

		ExtractOperations: func() (map[string]interface{}, []string) {

			b.mutex.Lock()
			tenantPatch := make(map[string]openapi.ConfigPutRequestTenantTenantName)
			for k, v := range b.tenantPatch {
				tenantPatch[k] = v
			}
			b.tenantPatch = make(map[string]openapi.ConfigPutRequestTenantTenantName)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			names := make([]string, 0, len(tenantPatch))

			for k, v := range tenantPatch {
				result[k] = v
				names = append(names, k)
			}

			return result, names
		},

		CheckPreExistence: nil,

		PrepareRequest: func(filteredData map[string]interface{}) interface{} {
			patchRequest := openapi.NewTenantsPutRequest()
			tenantMap := make(map[string]openapi.ConfigPutRequestTenantTenantName)

			for name, props := range filteredData {
				tenantMap[name] = props.(openapi.ConfigPutRequestTenantTenantName)
			}
			patchRequest.SetTenant(tenantMap)
			return patchRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.TenantsAPI.TenantsPatch(ctx).TenantsPutRequest(
				*request.(*openapi.TenantsPutRequest))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentTenantOps = true
			b.recentTenantOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkTenantDelete(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "tenant",
		OperationType: "DELETE",

		ExtractOperations: func() (map[string]interface{}, []string) {

			b.mutex.Lock()
			tenantNames := make([]string, len(b.tenantDelete))
			copy(tenantNames, b.tenantDelete)

			tenantDeleteMap := make(map[string]bool)
			for _, name := range tenantNames {
				tenantDeleteMap[name] = true
			}

			b.tenantDelete = make([]string, 0)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			for _, name := range tenantNames {
				result[name] = true
			}

			return result, tenantNames
		},

		CheckPreExistence: nil,

		PrepareRequest: func(filteredData map[string]interface{}) interface{} {
			names := make([]string, 0, len(filteredData))
			for name := range filteredData {
				names = append(names, name)
			}
			return names
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.TenantsAPI.TenantsDelete(ctx).TenantName(request.([]string))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentTenantOps = true
			b.recentTenantOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkGatewayDelete(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "gateway",
		OperationType: "DELETE",

		ExtractOperations: func() (map[string]interface{}, []string) {

			b.mutex.Lock()
			gatewayNames := make([]string, len(b.gatewayDelete))
			copy(gatewayNames, b.gatewayDelete)

			gatewayDeleteMap := make(map[string]bool)
			for _, name := range gatewayNames {
				gatewayDeleteMap[name] = true
			}

			b.gatewayDelete = make([]string, 0)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			for _, name := range gatewayNames {
				result[name] = true
			}

			return result, gatewayNames
		},

		CheckPreExistence: nil,

		PrepareRequest: func(filteredData map[string]interface{}) interface{} {
			names := make([]string, 0, len(filteredData))
			for name := range filteredData {
				names = append(names, name)
			}
			return names
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.GatewaysAPI.GatewaysDelete(ctx).GatewayName(request.([]string))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentGatewayOps = true
			b.recentGatewayOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkServicePut(ctx context.Context) diag.Diagnostics {
	var originalOperations map[string]openapi.ConfigPutRequestServiceServiceName

	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "service",
		OperationType: "PUT",

		ExtractOperations: func() (map[string]interface{}, []string) {

			b.mutex.Lock()
			originalOperations = make(map[string]openapi.ConfigPutRequestServiceServiceName)
			for k, v := range b.servicePut {
				originalOperations[k] = v
			}
			b.servicePut = make(map[string]openapi.ConfigPutRequestServiceServiceName)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			names := make([]string, 0, len(originalOperations))

			for k, v := range originalOperations {
				result[k] = v
				names = append(names, k)
			}

			return result, names
		},

		CheckPreExistence: func(ctx context.Context, resourceNames []string) ([]string, map[string]interface{}, error) {
			checker := ResourceExistenceCheck{
				ResourceType:  "service",
				OperationType: "PUT",
				FetchResources: func(ctx context.Context) (map[string]interface{}, error) {
					// First check if we have cached service data
					b.serviceResponsesMutex.RLock()
					if len(b.serviceResponses) > 0 {
						cachedData := make(map[string]interface{})
						for k, v := range b.serviceResponses {
							cachedData[k] = v
						}
						b.serviceResponsesMutex.RUnlock()

						tflog.Debug(ctx, "Using cached service data for pre-existence check", map[string]interface{}{
							"count": len(cachedData),
						})

						return cachedData, nil
					}
					b.serviceResponsesMutex.RUnlock()

					// Fall back to API call if no cache
					apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
					defer cancel()

					resp, err := b.client.ServicesAPI.ServicesGet(apiCtx).Execute()
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

					b.serviceResponsesMutex.Lock()
					for k, v := range result.Service {
						if vMap, ok := v.(map[string]interface{}); ok {
							b.serviceResponses[k] = vMap

							if name, ok := vMap["name"].(string); ok && name != k {
								b.serviceResponses[name] = vMap
							}
						}
					}
					b.serviceResponsesMutex.Unlock()

					return result.Service, nil
				},
			}

			filteredNames, err := b.FilterPreExistingResources(ctx, resourceNames, checker)
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
			putRequest := openapi.NewServicesPutRequest()
			serviceMap := make(map[string]openapi.ConfigPutRequestServiceServiceName)

			for name, props := range filteredData {
				serviceMap[name] = props.(openapi.ConfigPutRequestServiceServiceName)
			}
			putRequest.SetService(serviceMap)
			return putRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.ServicesAPI.ServicesPut(ctx).ServicesPutRequest(
				*request.(*openapi.ServicesPutRequest))
			return req.Execute()
		},

		ProcessResponse: func(ctx context.Context, resp *http.Response) error {
			delayTime := 2 * time.Second
			tflog.Debug(ctx, fmt.Sprintf("Waiting %v for auto-generated values to be assigned before fetching services", delayTime))
			time.Sleep(delayTime)

			fetchCtx, fetchCancel := context.WithTimeout(context.Background(), OperationTimeout)
			defer fetchCancel()

			tflog.Debug(ctx, "Fetching services after successful PUT operation to retrieve auto-generated values")
			servicesReq := b.client.ServicesAPI.ServicesGet(fetchCtx)
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

			b.serviceResponsesMutex.Lock()
			for serviceName, serviceData := range servicesData.Service {
				b.serviceResponses[serviceName] = serviceData

				if name, ok := serviceData["name"].(string); ok && name != serviceName {
					b.serviceResponses[name] = serviceData
				}
			}
			b.serviceResponsesMutex.Unlock()

			tflog.Debug(ctx, "Successfully stored service data for auto-generated fields", map[string]interface{}{
				"service_count": len(servicesData.Service),
			})

			return nil
		},

		UpdateRecentOps: func() {
			b.recentServiceOps = true
			b.recentServiceOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkServicePatch(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "service",
		OperationType: "PATCH",

		ExtractOperations: func() (map[string]interface{}, []string) {

			b.mutex.Lock()
			servicePatch := make(map[string]openapi.ConfigPutRequestServiceServiceName)
			for k, v := range b.servicePatch {
				servicePatch[k] = v
			}
			b.servicePatch = make(map[string]openapi.ConfigPutRequestServiceServiceName)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			names := make([]string, 0, len(servicePatch))

			for k, v := range servicePatch {
				result[k] = v
				names = append(names, k)
			}

			return result, names
		},

		CheckPreExistence: nil,

		PrepareRequest: func(filteredData map[string]interface{}) interface{} {
			patchRequest := openapi.NewServicesPutRequest()
			serviceMap := make(map[string]openapi.ConfigPutRequestServiceServiceName)

			for name, props := range filteredData {
				serviceMap[name] = props.(openapi.ConfigPutRequestServiceServiceName)
			}
			patchRequest.SetService(serviceMap)
			return patchRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.ServicesAPI.ServicesPatch(ctx).ServicesPutRequest(
				*request.(*openapi.ServicesPutRequest))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentServiceOps = true
			b.recentServiceOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkServiceDelete(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "service",
		OperationType: "DELETE",

		ExtractOperations: func() (map[string]interface{}, []string) {

			b.mutex.Lock()
			serviceNames := make([]string, len(b.serviceDelete))
			copy(serviceNames, b.serviceDelete)

			serviceDeleteMap := make(map[string]bool)
			for _, name := range serviceNames {
				serviceDeleteMap[name] = true
			}

			b.serviceDelete = make([]string, 0)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			for _, name := range serviceNames {
				result[name] = true
			}

			return result, serviceNames
		},

		CheckPreExistence: nil,

		PrepareRequest: func(filteredData map[string]interface{}) interface{} {
			names := make([]string, 0, len(filteredData))
			for name := range filteredData {
				names = append(names, name)
			}
			return names
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.ServicesAPI.ServicesDelete(ctx).ServiceName(request.([]string))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentServiceOps = true
			b.recentServiceOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkGatewayProfilePut(ctx context.Context) diag.Diagnostics {
	var originalOperations map[string]openapi.ConfigPutRequestGatewayProfileGatewayProfileName

	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "gateway_profile",
		OperationType: "PUT",

		ExtractOperations: func() (map[string]interface{}, []string) {

			b.mutex.Lock()
			originalOperations = make(map[string]openapi.ConfigPutRequestGatewayProfileGatewayProfileName)
			for k, v := range b.gatewayProfilePut {
				originalOperations[k] = v
			}
			b.gatewayProfilePut = make(map[string]openapi.ConfigPutRequestGatewayProfileGatewayProfileName)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			names := make([]string, 0, len(originalOperations))

			for k, v := range originalOperations {
				result[k] = v
				names = append(names, k)
			}

			return result, names
		},

		CheckPreExistence: func(ctx context.Context, resourceNames []string) ([]string, map[string]interface{}, error) {
			checker := ResourceExistenceCheck{
				ResourceType:  "gateway_profile",
				OperationType: "PUT",
				FetchResources: func(ctx context.Context) (map[string]interface{}, error) {
					// First check if we have cached Gateway Profile data
					b.gatewayProfileResponsesMutex.RLock()
					if len(b.gatewayProfileResponses) > 0 {
						cachedData := make(map[string]interface{})
						for k, v := range b.gatewayProfileResponses {
							cachedData[k] = v
						}
						b.gatewayProfileResponsesMutex.RUnlock()

						tflog.Debug(ctx, "Using cached Gateway Profile data for pre-existence check", map[string]interface{}{
							"count": len(cachedData),
						})

						return cachedData, nil
					}
					b.gatewayProfileResponsesMutex.RUnlock()

					// Fall back to API call if no cache
					apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
					defer cancel()

					resp, err := b.client.GatewayProfilesAPI.GatewayprofilesGet(apiCtx).Execute()
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

					b.gatewayProfileResponsesMutex.Lock()
					for k, v := range result.GatewayProfile {
						if vMap, ok := v.(map[string]interface{}); ok {
							b.gatewayProfileResponses[k] = vMap

							if name, ok := vMap["name"].(string); ok && name != k {
								b.gatewayProfileResponses[name] = vMap
							}
						}
					}
					b.gatewayProfileResponsesMutex.Unlock()

					return result.GatewayProfile, nil
				},
			}

			filteredNames, err := b.FilterPreExistingResources(ctx, resourceNames, checker)
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
			putRequest := openapi.NewGatewayprofilesPutRequest()
			gatewayProfileMap := make(map[string]openapi.ConfigPutRequestGatewayProfileGatewayProfileName)

			for name, props := range filteredData {
				gatewayProfileMap[name] = props.(openapi.ConfigPutRequestGatewayProfileGatewayProfileName)
			}
			putRequest.SetGatewayProfile(gatewayProfileMap)
			return putRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.GatewayProfilesAPI.GatewayprofilesPut(ctx).GatewayprofilesPutRequest(
				*request.(*openapi.GatewayprofilesPutRequest))
			return req.Execute()
		},

		ProcessResponse: func(ctx context.Context, resp *http.Response) error {
			delayTime := 2 * time.Second
			tflog.Debug(ctx, fmt.Sprintf("Waiting %v for auto-generated values to be assigned before fetching Gateways", delayTime))
			time.Sleep(delayTime)

			fetchCtx, fetchCancel := context.WithTimeout(context.Background(), OperationTimeout)
			defer fetchCancel()

			tflog.Debug(ctx, "Fetching Gateways after successful PUT operation to retrieve auto-generated values")
			gatewaysReq := b.client.GatewaysAPI.GatewaysGet(fetchCtx)
			gatewaysResp, fetchErr := gatewaysReq.Execute()

			if fetchErr != nil {
				tflog.Error(ctx, "Failed to fetch Gateways after PUT for auto-generated fields", map[string]interface{}{
					"error": fetchErr.Error(),
				})
				return fetchErr
			}

			defer gatewaysResp.Body.Close()

			var gatewaysData struct {
				Gateway map[string]map[string]interface{} `json:"gateway"`
			}

			if respErr := json.NewDecoder(gatewaysResp.Body).Decode(&gatewaysData); respErr != nil {
				tflog.Error(ctx, "Failed to decode Gateways response for auto-generated fields", map[string]interface{}{
					"error": respErr.Error(),
				})
				return respErr
			}

			b.gatewayResponsesMutex.Lock()
			for gatewayName, gatewayData := range gatewaysData.Gateway {
				b.gatewayResponses[gatewayName] = gatewayData

				if name, ok := gatewayData["name"].(string); ok && name != gatewayName {
					b.gatewayResponses[name] = gatewayData
				}
			}
			b.gatewayResponsesMutex.Unlock()

			tflog.Debug(ctx, "Successfully stored Gateway data for auto-generated fields", map[string]interface{}{
				"gateway_count": len(gatewaysData.Gateway),
			})

			return nil
		},

		UpdateRecentOps: func() {
			b.recentGatewayProfileOps = true
			b.recentGatewayProfileOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkGatewayProfilePatch(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "gateway_profile",
		OperationType: "PATCH",

		ExtractOperations: func() (map[string]interface{}, []string) {

			b.mutex.Lock()
			gatewayProfilePatch := make(map[string]openapi.ConfigPutRequestGatewayProfileGatewayProfileName)
			for k, v := range b.gatewayProfilePatch {
				gatewayProfilePatch[k] = v
			}
			b.gatewayProfilePatch = make(map[string]openapi.ConfigPutRequestGatewayProfileGatewayProfileName)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			names := make([]string, 0, len(gatewayProfilePatch))

			for k, v := range gatewayProfilePatch {
				result[k] = v
				names = append(names, k)
			}

			return result, names
		},

		CheckPreExistence: nil,

		PrepareRequest: func(filteredData map[string]interface{}) interface{} {
			patchRequest := openapi.NewGatewayprofilesPutRequest()
			gatewayProfileMap := make(map[string]openapi.ConfigPutRequestGatewayProfileGatewayProfileName)

			for name, props := range filteredData {
				gatewayProfileMap[name] = props.(openapi.ConfigPutRequestGatewayProfileGatewayProfileName)
			}
			patchRequest.SetGatewayProfile(gatewayProfileMap)
			return patchRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.GatewayProfilesAPI.GatewayprofilesPatch(ctx).GatewayprofilesPutRequest(
				*request.(*openapi.GatewayprofilesPutRequest))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentGatewayProfileOps = true
			b.recentGatewayProfileOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkGatewayProfileDelete(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "gateway_profile",
		OperationType: "DELETE",

		ExtractOperations: func() (map[string]interface{}, []string) {

			b.mutex.Lock()
			gatewayProfileNames := make([]string, len(b.gatewayProfileDelete))
			copy(gatewayProfileNames, b.gatewayProfileDelete)

			gatewayProfileDeleteMap := make(map[string]bool)
			for _, name := range gatewayProfileNames {
				gatewayProfileDeleteMap[name] = true
			}

			b.gatewayProfileDelete = make([]string, 0)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			for _, name := range gatewayProfileNames {
				result[name] = true
			}

			return result, gatewayProfileNames
		},

		CheckPreExistence: nil,

		PrepareRequest: func(filteredData map[string]interface{}) interface{} {
			names := make([]string, 0, len(filteredData))
			for name := range filteredData {
				names = append(names, name)
			}
			return names
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.GatewayProfilesAPI.GatewayprofilesDelete(ctx).ProfileName(request.([]string))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentGatewayProfileOps = true
			b.recentGatewayProfileOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkEthPortProfilePut(ctx context.Context) diag.Diagnostics {
	var originalOperations map[string]openapi.ConfigPutRequestEthPortProfileEthPortProfileName

	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "eth_port_profile",
		OperationType: "PUT",

		ExtractOperations: func() (map[string]interface{}, []string) {

			b.mutex.Lock()
			originalOperations = make(map[string]openapi.ConfigPutRequestEthPortProfileEthPortProfileName)
			for k, v := range b.ethPortProfilePut {
				originalOperations[k] = v
			}
			b.ethPortProfilePut = make(map[string]openapi.ConfigPutRequestEthPortProfileEthPortProfileName)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			names := make([]string, 0, len(originalOperations))

			for k, v := range originalOperations {
				result[k] = v
				names = append(names, k)
			}

			return result, names
		},

		CheckPreExistence: func(ctx context.Context, resourceNames []string) ([]string, map[string]interface{}, error) {
			checker := ResourceExistenceCheck{
				ResourceType:  "eth_port_profile",
				OperationType: "PUT",
				FetchResources: func(ctx context.Context) (map[string]interface{}, error) {
					// First check if we have cached EthPortProfile data
					b.ethPortProfileResponsesMutex.RLock()
					if len(b.ethPortProfileResponses) > 0 {
						cachedData := make(map[string]interface{})
						for k, v := range b.ethPortProfileResponses {
							cachedData[k] = v
						}
						b.ethPortProfileResponsesMutex.RUnlock()

						tflog.Debug(ctx, "Using cached EthPortProfile data for pre-existence check", map[string]interface{}{
							"count": len(cachedData),
						})

						return cachedData, nil
					}
					b.ethPortProfileResponsesMutex.RUnlock()

					// Fall back to API call if no cache
					apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
					defer cancel()

					resp, err := b.client.EthPortProfilesAPI.EthportprofilesGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						EthPortProfile map[string]interface{} `json:"eth_port_profile"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}

					b.ethPortProfileResponsesMutex.Lock()
					for k, v := range result.EthPortProfile {
						if vMap, ok := v.(map[string]interface{}); ok {
							b.ethPortProfileResponses[k] = vMap

							if name, ok := vMap["name"].(string); ok && name != k {
								b.ethPortProfileResponses[name] = vMap
							}
						}
					}
					b.ethPortProfileResponsesMutex.Unlock()

					return result.EthPortProfile, nil
				},
			}

			filteredNames, err := b.FilterPreExistingResources(ctx, resourceNames, checker)
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
			putRequest := openapi.NewEthportprofilesPutRequest()
			ethPortProfileMap := make(map[string]openapi.ConfigPutRequestEthPortProfileEthPortProfileName)

			for name, props := range filteredData {
				ethPortProfileMap[name] = props.(openapi.ConfigPutRequestEthPortProfileEthPortProfileName)
			}
			putRequest.SetEthPortProfile(ethPortProfileMap)
			return putRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.EthPortProfilesAPI.EthportprofilesPut(ctx).EthportprofilesPutRequest(
				*request.(*openapi.EthportprofilesPutRequest))
			return req.Execute()
		},

		ProcessResponse: func(ctx context.Context, resp *http.Response) error {
			delayTime := 2 * time.Second
			tflog.Debug(ctx, fmt.Sprintf("Waiting %v for auto-generated values to be assigned before fetching EthPortProfiles", delayTime))
			time.Sleep(delayTime)

			fetchCtx, fetchCancel := context.WithTimeout(context.Background(), OperationTimeout)
			defer fetchCancel()

			tflog.Debug(ctx, "Fetching EthPortProfiles after successful PUT operation to retrieve auto-generated values")
			profilesReq := b.client.EthPortProfilesAPI.EthportprofilesGet(fetchCtx)
			profilesResp, fetchErr := profilesReq.Execute()

			if fetchErr != nil {
				tflog.Error(ctx, "Failed to fetch EthPortProfiles after PUT for auto-generated fields", map[string]interface{}{
					"error": fetchErr.Error(),
				})
				return fetchErr
			}

			defer profilesResp.Body.Close()

			var profilesData struct {
				EthPortProfile map[string]map[string]interface{} `json:"eth_port_profile"`
			}

			if respErr := json.NewDecoder(profilesResp.Body).Decode(&profilesData); respErr != nil {
				tflog.Error(ctx, "Failed to decode EthPortProfiles response for auto-generated fields", map[string]interface{}{
					"error": respErr.Error(),
				})
				return respErr
			}

			b.ethPortProfileResponsesMutex.Lock()
			for profileName, profileData := range profilesData.EthPortProfile {
				b.ethPortProfileResponses[profileName] = profileData

				if name, ok := profileData["name"].(string); ok && name != profileName {
					b.ethPortProfileResponses[name] = profileData
				}
			}
			b.ethPortProfileResponsesMutex.Unlock()

			tflog.Debug(ctx, "Successfully stored EthPortProfile data for auto-generated fields", map[string]interface{}{
				"eth_port_profile_count": len(profilesData.EthPortProfile),
			})

			return nil
		},

		UpdateRecentOps: func() {
			b.recentEthPortProfileOps = true
			b.recentEthPortProfileOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkEthPortSettingsPut(ctx context.Context) diag.Diagnostics {
	var originalOperations map[string]openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName

	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "eth_port_settings",
		OperationType: "PUT",

		ExtractOperations: func() (map[string]interface{}, []string) {

			b.mutex.Lock()
			originalOperations = make(map[string]openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName)
			for k, v := range b.ethPortSettingsPut {
				originalOperations[k] = v
			}
			b.ethPortSettingsPut = make(map[string]openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			names := make([]string, 0, len(originalOperations))

			for k, v := range originalOperations {
				result[k] = v
				names = append(names, k)
			}

			return result, names
		},

		CheckPreExistence: func(ctx context.Context, resourceNames []string) ([]string, map[string]interface{}, error) {
			checker := ResourceExistenceCheck{
				ResourceType:  "eth_port_settings",
				OperationType: "PUT",
				FetchResources: func(ctx context.Context) (map[string]interface{}, error) {
					// First check if we have cached EthPortSettings data
					b.ethPortSettingsResponsesMutex.RLock()
					if len(b.ethPortSettingsResponses) > 0 {
						cachedData := make(map[string]interface{})
						for k, v := range b.ethPortSettingsResponses {
							cachedData[k] = v
						}
						b.ethPortSettingsResponsesMutex.RUnlock()

						tflog.Debug(ctx, "Using cached EthPortSettings data for pre-existence check", map[string]interface{}{
							"count": len(cachedData),
						})

						return cachedData, nil
					}
					b.ethPortSettingsResponsesMutex.RUnlock()

					// Fall back to API call if no cache
					apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
					defer cancel()

					resp, err := b.client.EthPortSettingsAPI.EthportsettingsGet(apiCtx).Execute()
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

					b.ethPortSettingsResponsesMutex.Lock()
					for k, v := range result.EthPortSettings {
						if vMap, ok := v.(map[string]interface{}); ok {
							b.ethPortSettingsResponses[k] = vMap

							if name, ok := vMap["name"].(string); ok && name != k {
								b.ethPortSettingsResponses[name] = vMap
							}
						}
					}
					b.ethPortSettingsResponsesMutex.Unlock()

					return result.EthPortSettings, nil
				},
			}

			filteredNames, err := b.FilterPreExistingResources(ctx, resourceNames, checker)
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
			putRequest := openapi.NewEthportsettingsPutRequest()
			ethPortSettingsMap := make(map[string]openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName)

			for name, props := range filteredData {
				ethPortSettingsMap[name] = props.(openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName)
			}
			putRequest.SetEthPortSettings(ethPortSettingsMap)
			return putRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.EthPortSettingsAPI.EthportsettingsPut(ctx).EthportsettingsPutRequest(
				*request.(*openapi.EthportsettingsPutRequest))
			return req.Execute()
		},

		ProcessResponse: func(ctx context.Context, resp *http.Response) error {
			delayTime := 2 * time.Second
			tflog.Debug(ctx, fmt.Sprintf("Waiting %v for auto-generated values to be assigned before fetching EthPortSettings", delayTime))
			time.Sleep(delayTime)

			fetchCtx, fetchCancel := context.WithTimeout(context.Background(), OperationTimeout)
			defer fetchCancel()

			tflog.Debug(ctx, "Fetching EthPortSettings after successful PUT operation to retrieve auto-generated values")
			settingsReq := b.client.EthPortSettingsAPI.EthportsettingsGet(fetchCtx)
			settingsResp, fetchErr := settingsReq.Execute()

			if fetchErr != nil {
				tflog.Error(ctx, "Failed to fetch EthPortSettings after PUT for auto-generated fields", map[string]interface{}{
					"error": fetchErr.Error(),
				})
				return fetchErr
			}

			defer settingsResp.Body.Close()

			var settingsData struct {
				EthPortSettings map[string]map[string]interface{} `json:"eth_port_settings"`
			}

			if respErr := json.NewDecoder(settingsResp.Body).Decode(&settingsData); respErr != nil {
				tflog.Error(ctx, "Failed to decode EthPortSettings response for auto-generated fields", map[string]interface{}{
					"error": respErr.Error(),
				})
				return respErr
			}

			b.ethPortSettingsResponsesMutex.Lock()
			for settingsName, settingsData := range settingsData.EthPortSettings {
				b.ethPortSettingsResponses[settingsName] = settingsData

				if name, ok := settingsData["name"].(string); ok && name != settingsName {
					b.ethPortSettingsResponses[name] = settingsData
				}
			}
			b.ethPortSettingsResponsesMutex.Unlock()

			tflog.Debug(ctx, "Successfully stored EthPortSettings data for auto-generated fields", map[string]interface{}{
				"eth_port_settings_count": len(settingsData.EthPortSettings),
			})

			return nil
		},

		UpdateRecentOps: func() {
			b.recentEthPortSettingsOps = true
			b.recentEthPortSettingsOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkEthPortProfilePatch(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "eth_port_profile",
		OperationType: "PATCH",

		ExtractOperations: func() (map[string]interface{}, []string) {

			b.mutex.Lock()
			ethPortProfilePatch := make(map[string]openapi.ConfigPutRequestEthPortProfileEthPortProfileName)
			for k, v := range b.ethPortProfilePatch {
				ethPortProfilePatch[k] = v
			}
			b.ethPortProfilePatch = make(map[string]openapi.ConfigPutRequestEthPortProfileEthPortProfileName)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			names := make([]string, 0, len(ethPortProfilePatch))

			for k, v := range ethPortProfilePatch {
				result[k] = v
				names = append(names, k)
			}

			return result, names
		},

		CheckPreExistence: nil,

		PrepareRequest: func(filteredData map[string]interface{}) interface{} {
			patchRequest := openapi.NewEthportprofilesPutRequest()
			ethPortProfileMap := make(map[string]openapi.ConfigPutRequestEthPortProfileEthPortProfileName)

			for name, props := range filteredData {
				ethPortProfileMap[name] = props.(openapi.ConfigPutRequestEthPortProfileEthPortProfileName)
			}
			patchRequest.SetEthPortProfile(ethPortProfileMap)
			return patchRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.EthPortProfilesAPI.EthportprofilesPatch(ctx).EthportprofilesPutRequest(
				*request.(*openapi.EthportprofilesPutRequest))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentEthPortProfileOps = true
			b.recentEthPortProfileOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkEthPortSettingsPatch(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "eth_port_settings",
		OperationType: "PATCH",

		ExtractOperations: func() (map[string]interface{}, []string) {

			b.mutex.Lock()
			ethPortSettingsPatch := make(map[string]openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName)
			for k, v := range b.ethPortSettingsPatch {
				ethPortSettingsPatch[k] = v
			}
			b.ethPortSettingsPatch = make(map[string]openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			names := make([]string, 0, len(ethPortSettingsPatch))

			for k, v := range ethPortSettingsPatch {
				result[k] = v
				names = append(names, k)
			}

			return result, names
		},

		CheckPreExistence: nil,

		PrepareRequest: func(filteredData map[string]interface{}) interface{} {
			patchRequest := openapi.NewEthportsettingsPutRequest()
			ethPortSettingsMap := make(map[string]openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName)

			for name, props := range filteredData {
				ethPortSettingsMap[name] = props.(openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName)
			}
			patchRequest.SetEthPortSettings(ethPortSettingsMap)
			return patchRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.EthPortSettingsAPI.EthportsettingsPatch(ctx).EthportsettingsPutRequest(
				*request.(*openapi.EthportsettingsPutRequest))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentEthPortSettingsOps = true
			b.recentEthPortSettingsOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkEthPortProfileDelete(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "eth_port_profile",
		OperationType: "DELETE",

		ExtractOperations: func() (map[string]interface{}, []string) {

			b.mutex.Lock()
			ethPortProfileNames := make([]string, len(b.ethPortProfileDelete))
			copy(ethPortProfileNames, b.ethPortProfileDelete)

			ethPortProfileDeleteMap := make(map[string]bool)
			for _, name := range ethPortProfileNames {
				ethPortProfileDeleteMap[name] = true
			}

			b.ethPortProfileDelete = make([]string, 0)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			for _, name := range ethPortProfileNames {
				result[name] = true
			}

			return result, ethPortProfileNames
		},

		CheckPreExistence: nil,

		PrepareRequest: func(filteredData map[string]interface{}) interface{} {
			names := make([]string, 0, len(filteredData))
			for name := range filteredData {
				names = append(names, name)
			}
			return names
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.EthPortProfilesAPI.EthportprofilesDelete(ctx).ProfileName(request.([]string))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentEthPortProfileOps = true
			b.recentEthPortProfileOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkEthPortSettingsDelete(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "eth_port_settings",
		OperationType: "DELETE",

		ExtractOperations: func() (map[string]interface{}, []string) {

			b.mutex.Lock()
			ethPortSettingsNames := make([]string, len(b.ethPortSettingsDelete))
			copy(ethPortSettingsNames, b.ethPortSettingsDelete)

			ethPortSettingsDeleteMap := make(map[string]bool)
			for _, name := range ethPortSettingsNames {
				ethPortSettingsDeleteMap[name] = true
			}

			b.ethPortSettingsDelete = make([]string, 0)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			for _, name := range ethPortSettingsNames {
				result[name] = true
			}

			return result, ethPortSettingsNames
		},

		CheckPreExistence: nil,

		PrepareRequest: func(filteredData map[string]interface{}) interface{} {
			names := make([]string, 0, len(filteredData))
			for name := range filteredData {
				names = append(names, name)
			}
			return names
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.EthPortSettingsAPI.EthportsettingsDelete(ctx).PortName(request.([]string))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentEthPortSettingsOps = true
			b.recentEthPortSettingsOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkBundlePatch(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "bundle",
		OperationType: "PATCH",

		ExtractOperations: func() (map[string]interface{}, []string) {

			b.mutex.Lock()
			bundlePatch := make(map[string]openapi.BundlesPatchRequestEndpointBundleValue)
			for k, v := range b.bundlePatch {
				bundlePatch[k] = v
			}
			b.bundlePatch = make(map[string]openapi.BundlesPatchRequestEndpointBundleValue)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			names := make([]string, 0, len(bundlePatch))

			for k, v := range bundlePatch {
				result[k] = v
				names = append(names, k)
			}

			return result, names
		},

		CheckPreExistence: nil,

		PrepareRequest: func(filteredData map[string]interface{}) interface{} {
			patchRequest := openapi.NewBundlesPatchRequest()
			bundleMap := make(map[string]openapi.BundlesPatchRequestEndpointBundleValue)

			for name, props := range filteredData {
				bundleMap[name] = props.(openapi.BundlesPatchRequestEndpointBundleValue)
			}
			patchRequest.SetEndpointBundle(bundleMap)
			return patchRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.BundlesAPI.BundlesPatch(ctx).BundlesPatchRequest(
				*request.(*openapi.BundlesPatchRequest))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentBundleOps = true
			b.recentBundleOpTime = time.Now()
		},
	})
}
