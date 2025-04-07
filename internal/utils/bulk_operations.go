package utils

import (
	"context"
	"fmt"
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

		// Initialize with no recent operations
		recentGatewayOps:         false,
		recentLagOps:             false,
		recentServiceOps:         false,
		recentTenantOps:          false,
		recentGatewayProfileOps:  false,
		recentEthPortProfileOps:  false,
		recentEthPortSettingsOps: false,
		recentBundleOps:          false,
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
	defer b.operationMutex.Unlock()

	failCount := 0
	for opID, op := range b.pendingOperations {
		if op.Status == OperationPending {
			op.Status = OperationFailed
			op.Error = fmt.Errorf("Operation aborted due to previous failure: %v", err)
			b.operationErrors[opID] = op.Error
			b.operationResults[opID] = false

			if waitCh, ok := b.operationWaitChannels[opID]; ok {
				close(waitCh)
			}
			failCount++
		}
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

func (b *BulkOperationManager) HasPendingOrRecentGatewayOperations() bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Check if any gateway operations are pending
	hasPending := len(b.gatewayPut) > 0 || len(b.gatewayPatch) > 0 || len(b.gatewayDelete) > 0

	// Check if we've recently had operations (within the last 5 seconds)
	hasRecent := b.recentGatewayOps && time.Since(b.recentGatewayOpTime) < 5*time.Second

	return hasPending || hasRecent
}

func (b *BulkOperationManager) HasPendingOrRecentLagOperations() bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Check if any LAG operations are pending
	hasPending := len(b.lagPut) > 0 || len(b.lagPatch) > 0 || len(b.lagDelete) > 0

	// Check if we've recently had operations (within the last 5 seconds)
	hasRecent := b.recentLagOps && time.Since(b.recentLagOpTime) < 5*time.Second

	return hasPending || hasRecent
}

func (b *BulkOperationManager) HasPendingOrRecentServiceOperations() bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Check if any Service operations are pending
	hasPending := len(b.servicePut) > 0 || len(b.servicePatch) > 0 || len(b.serviceDelete) > 0

	// Check if we've recently had operations (within the last 5 seconds)
	hasRecent := b.recentServiceOps && time.Since(b.recentServiceOpTime) < 5*time.Second

	return hasPending || hasRecent
}

func (b *BulkOperationManager) HasPendingOrRecentTenantOperations() bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Check if any Tenant operations are pending
	hasPending := len(b.tenantPut) > 0 || len(b.tenantPatch) > 0 || len(b.tenantDelete) > 0

	// Check if we've recently had operations (within the last 5 seconds)
	hasRecent := b.recentTenantOps && time.Since(b.recentTenantOpTime) < 5*time.Second

	return hasPending || hasRecent
}

func (b *BulkOperationManager) HasPendingOrRecentGatewayProfileOperations() bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Check if any Gateway Profile operations are pending
	hasPending := len(b.gatewayProfilePut) > 0 || len(b.gatewayProfilePatch) > 0 || len(b.gatewayProfileDelete) > 0

	// Check if we've recently had operations (within the last 5 seconds)
	hasRecent := b.recentGatewayProfileOps && time.Since(b.recentGatewayProfileOpTime) < 5*time.Second

	return hasPending || hasRecent
}

func (b *BulkOperationManager) HasPendingOrRecentEthPortProfileOperations() bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Check if any Eth Port Profile operations are pending
	hasPending := len(b.ethPortProfilePut) > 0 || len(b.ethPortProfilePatch) > 0 || len(b.ethPortProfileDelete) > 0

	// Check if we've recently had operations (within the last 5 seconds)
	hasRecent := b.recentEthPortProfileOps && time.Since(b.recentEthPortProfileOpTime) < 5*time.Second

	return hasPending || hasRecent
}

func (b *BulkOperationManager) HasPendingOrRecentEthPortSettingsOperations() bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Check if any Eth Port Profile operations are pending
	hasPending := len(b.ethPortSettingsPut) > 0 || len(b.ethPortSettingsPatch) > 0 || len(b.ethPortSettingsDelete) > 0

	// Check if we've recently had operations (within the last 5 seconds)
	hasRecent := b.recentEthPortSettingsOps && time.Since(b.recentEthPortSettingsOpTime) < 5*time.Second

	return hasPending || hasRecent
}

func (b *BulkOperationManager) HasPendingOrRecentBundleOperations() bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Check if any bundle operations are pending
	hasPending := len(b.bundlePatch) > 0

	// Check if we've recently had operations (within the last 5 seconds)
	hasRecent := b.recentBundleOps && time.Since(b.recentBundleOpTime) < 5*time.Second

	return hasPending || hasRecent
}

func (b *BulkOperationManager) AddBundlePatch(ctx context.Context, bundleName string, props openapi.BundlesPatchRequestEndpointBundleValue) string {
	b.mutex.Lock()
	if b.bundlePatch == nil {
		b.bundlePatch = make(map[string]openapi.BundlesPatchRequestEndpointBundleValue)
	}
	b.bundlePatch[bundleName] = props
	b.mutex.Unlock()

	operationID := generateOperationID("bundle", bundleName, "PATCH")

	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	b.pendingOperations[operationID] = &Operation{
		ResourceType:  "bundle",
		ResourceName:  bundleName,
		OperationType: "PATCH",
		Status:        OperationPending,
	}

	b.operationWaitChannels[operationID] = make(chan struct{})

	now := time.Now()
	b.lastOperationTime = now
	if b.batchStartTime.IsZero() {
		b.batchStartTime = now
	}

	tflog.Debug(ctx, "Added Bundle to PATCH batch", map[string]interface{}{
		"bundle_name":  bundleName,
		"batch_size":   len(b.bundlePatch),
		"operation_id": operationID,
	})

	return operationID
}

func (b *BulkOperationManager) AddGatewayPut(ctx context.Context, gatewayName string, props openapi.ConfigPutRequestGatewayGatewayName) string {
	b.mutex.Lock()
	b.gatewayPut[gatewayName] = props
	b.mutex.Unlock()

	operationID := generateOperationID("gateway", gatewayName, "PUT")

	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	b.pendingOperations[operationID] = &Operation{
		ResourceType:  "gateway",
		ResourceName:  gatewayName,
		OperationType: "PUT",
		Status:        OperationPending,
	}

	b.operationWaitChannels[operationID] = make(chan struct{})

	now := time.Now()
	b.lastOperationTime = now
	if b.batchStartTime.IsZero() {
		b.batchStartTime = now
	}

	return operationID
}

func (b *BulkOperationManager) AddGatewayPatch(ctx context.Context, gatewayName string, props openapi.ConfigPutRequestGatewayGatewayName) string {
	b.mutex.Lock()
	b.gatewayPatch[gatewayName] = props
	b.mutex.Unlock()

	operationID := generateOperationID("gateway", gatewayName, "PATCH")

	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	b.pendingOperations[operationID] = &Operation{
		ResourceType:  "gateway",
		ResourceName:  gatewayName,
		OperationType: "PATCH",
		Status:        OperationPending,
	}

	b.operationWaitChannels[operationID] = make(chan struct{})

	now := time.Now()
	b.lastOperationTime = now
	if b.batchStartTime.IsZero() {
		b.batchStartTime = now
	}

	tflog.Debug(ctx, "Added Gateway to PATCH batch", map[string]interface{}{
		"gateway_name": gatewayName,
		"batch_size":   len(b.gatewayPatch),
		"operation_id": operationID,
	})

	return operationID
}

func (b *BulkOperationManager) AddGatewayDelete(ctx context.Context, gatewayName string) string {
	b.mutex.Lock()
	b.gatewayDelete = append(b.gatewayDelete, gatewayName)
	b.mutex.Unlock()

	operationID := generateOperationID("gateway", gatewayName, "DELETE")

	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	b.pendingOperations[operationID] = &Operation{
		ResourceType:  "gateway",
		ResourceName:  gatewayName,
		OperationType: "DELETE",
		Status:        OperationPending,
	}

	b.operationWaitChannels[operationID] = make(chan struct{})

	now := time.Now()
	b.lastOperationTime = now
	if b.batchStartTime.IsZero() {
		b.batchStartTime = now
	}

	tflog.Debug(ctx, "Added Gateway to DELETE batch", map[string]interface{}{
		"gateway_name": gatewayName,
		"batch_size":   len(b.gatewayDelete),
		"operation_id": operationID,
	})

	return operationID
}

func (b *BulkOperationManager) AddLagPut(ctx context.Context, lagName string, props openapi.ConfigPutRequestLagLagName) string {
	b.mutex.Lock()
	b.lagPut[lagName] = props
	b.mutex.Unlock()

	operationID := generateOperationID("lag", lagName, "PUT")

	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	b.pendingOperations[operationID] = &Operation{
		ResourceType:  "lag",
		ResourceName:  lagName,
		OperationType: "PUT",
		Status:        OperationPending,
	}

	b.operationWaitChannels[operationID] = make(chan struct{})

	now := time.Now()
	b.lastOperationTime = now
	if b.batchStartTime.IsZero() {
		b.batchStartTime = now
	}

	tflog.Debug(ctx, "Added LAG to PUT batch", map[string]interface{}{
		"lag_name":     lagName,
		"batch_size":   len(b.lagPut),
		"operation_id": operationID,
	})

	return operationID
}

func (b *BulkOperationManager) AddLagPatch(ctx context.Context, lagName string, props openapi.ConfigPutRequestLagLagName) string {
	b.mutex.Lock()
	b.lagPatch[lagName] = props
	b.mutex.Unlock()

	operationID := generateOperationID("lag", lagName, "PATCH")

	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	b.pendingOperations[operationID] = &Operation{
		ResourceType:  "lag",
		ResourceName:  lagName,
		OperationType: "PATCH",
		Status:        OperationPending,
	}

	b.operationWaitChannels[operationID] = make(chan struct{})

	now := time.Now()
	b.lastOperationTime = now
	if b.batchStartTime.IsZero() {
		b.batchStartTime = now
	}

	tflog.Debug(ctx, "Added LAG to PATCH batch", map[string]interface{}{
		"lag_name":     lagName,
		"batch_size":   len(b.lagPatch),
		"operation_id": operationID,
	})

	return operationID
}

func (b *BulkOperationManager) AddLagDelete(ctx context.Context, lagName string) string {
	b.mutex.Lock()
	b.lagDelete = append(b.lagDelete, lagName)
	b.mutex.Unlock()

	operationID := generateOperationID("lag", lagName, "DELETE")

	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	b.pendingOperations[operationID] = &Operation{
		ResourceType:  "lag",
		ResourceName:  lagName,
		OperationType: "DELETE",
		Status:        OperationPending,
	}

	b.operationWaitChannels[operationID] = make(chan struct{})

	now := time.Now()
	b.lastOperationTime = now
	if b.batchStartTime.IsZero() {
		b.batchStartTime = now
	}

	tflog.Debug(ctx, "Added LAG to DELETE batch", map[string]interface{}{
		"lag_name":     lagName,
		"batch_size":   len(b.lagDelete),
		"operation_id": operationID,
	})

	return operationID
}

func (b *BulkOperationManager) AddTenantPut(ctx context.Context, tenantName string, props openapi.ConfigPutRequestTenantTenantName) string {
	b.mutex.Lock()
	b.tenantPut[tenantName] = props
	b.mutex.Unlock()

	operationID := generateOperationID("tenant", tenantName, "PUT")

	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	b.pendingOperations[operationID] = &Operation{
		ResourceType:  "tenant",
		ResourceName:  tenantName,
		OperationType: "PUT",
		Status:        OperationPending,
	}

	b.operationWaitChannels[operationID] = make(chan struct{})

	now := time.Now()
	b.lastOperationTime = now
	if b.batchStartTime.IsZero() {
		b.batchStartTime = now
	}

	tflog.Debug(ctx, "Added Tenant to PUT batch", map[string]interface{}{
		"tenant_name":  tenantName,
		"batch_size":   len(b.tenantPut),
		"operation_id": operationID,
	})

	return operationID
}

func (b *BulkOperationManager) AddTenantPatch(ctx context.Context, tenantName string, props openapi.ConfigPutRequestTenantTenantName) string {
	b.mutex.Lock()
	b.tenantPatch[tenantName] = props
	b.mutex.Unlock()

	operationID := generateOperationID("tenant", tenantName, "PATCH")

	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	b.pendingOperations[operationID] = &Operation{
		ResourceType:  "tenant",
		ResourceName:  tenantName,
		OperationType: "PATCH",
		Status:        OperationPending,
	}

	b.operationWaitChannels[operationID] = make(chan struct{})

	now := time.Now()
	b.lastOperationTime = now
	if b.batchStartTime.IsZero() {
		b.batchStartTime = now
	}

	tflog.Debug(ctx, "Added Tenant to PATCH batch", map[string]interface{}{
		"tenant_name":  tenantName,
		"batch_size":   len(b.tenantPatch),
		"operation_id": operationID,
	})

	return operationID
}

func (b *BulkOperationManager) AddTenantDelete(ctx context.Context, tenantName string) string {
	b.mutex.Lock()
	b.tenantDelete = append(b.tenantDelete, tenantName)
	b.mutex.Unlock()

	operationID := generateOperationID("tenant", tenantName, "DELETE")

	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	b.pendingOperations[operationID] = &Operation{
		ResourceType:  "tenant",
		ResourceName:  tenantName,
		OperationType: "DELETE",
		Status:        OperationPending,
	}

	b.operationWaitChannels[operationID] = make(chan struct{})

	now := time.Now()
	b.lastOperationTime = now
	if b.batchStartTime.IsZero() {
		b.batchStartTime = now
	}

	tflog.Debug(ctx, "Added Tenant to DELETE batch", map[string]interface{}{
		"tenant_name":  tenantName,
		"batch_size":   len(b.tenantDelete),
		"operation_id": operationID,
	})

	return operationID
}

func (b *BulkOperationManager) AddServicePut(ctx context.Context, serviceName string, props openapi.ConfigPutRequestServiceServiceName) string {
	b.mutex.Lock()
	b.servicePut[serviceName] = props
	b.mutex.Unlock()

	operationID := generateOperationID("service", serviceName, "PUT")

	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	b.pendingOperations[operationID] = &Operation{
		ResourceType:  "service",
		ResourceName:  serviceName,
		OperationType: "PUT",
		Status:        OperationPending,
	}

	b.operationWaitChannels[operationID] = make(chan struct{})

	now := time.Now()
	b.lastOperationTime = now
	if b.batchStartTime.IsZero() {
		b.batchStartTime = now
	}

	tflog.Debug(ctx, "Added Service to PUT batch", map[string]interface{}{
		"service_name": serviceName,
		"batch_size":   len(b.servicePut),
		"operation_id": operationID,
	})

	return operationID
}

func (b *BulkOperationManager) AddServicePatch(ctx context.Context, serviceName string, props openapi.ConfigPutRequestServiceServiceName) string {
	b.mutex.Lock()
	b.servicePatch[serviceName] = props
	b.mutex.Unlock()

	operationID := generateOperationID("service", serviceName, "PATCH")

	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	b.pendingOperations[operationID] = &Operation{
		ResourceType:  "service",
		ResourceName:  serviceName,
		OperationType: "PATCH",
		Status:        OperationPending,
	}

	b.operationWaitChannels[operationID] = make(chan struct{})

	now := time.Now()
	b.lastOperationTime = now
	if b.batchStartTime.IsZero() {
		b.batchStartTime = now
	}

	tflog.Debug(ctx, "Added Service to PATCH batch", map[string]interface{}{
		"service_name": serviceName,
		"batch_size":   len(b.servicePatch),
		"operation_id": operationID,
	})

	return operationID
}

func (b *BulkOperationManager) AddServiceDelete(ctx context.Context, serviceName string) string {
	b.mutex.Lock()
	b.serviceDelete = append(b.serviceDelete, serviceName)
	b.mutex.Unlock()

	operationID := generateOperationID("service", serviceName, "DELETE")

	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	b.pendingOperations[operationID] = &Operation{
		ResourceType:  "service",
		ResourceName:  serviceName,
		OperationType: "DELETE",
		Status:        OperationPending,
	}

	b.operationWaitChannels[operationID] = make(chan struct{})

	now := time.Now()
	b.lastOperationTime = now
	if b.batchStartTime.IsZero() {
		b.batchStartTime = now
	}

	tflog.Debug(ctx, "Added Service to DELETE batch", map[string]interface{}{
		"service_name": serviceName,
		"batch_size":   len(b.serviceDelete),
		"operation_id": operationID,
	})

	return operationID
}

func (b *BulkOperationManager) AddGatewayProfilePut(ctx context.Context, profileName string, props openapi.ConfigPutRequestGatewayProfileGatewayProfileName) string {
	b.mutex.Lock()
	b.gatewayProfilePut[profileName] = props
	b.mutex.Unlock()

	operationID := generateOperationID("gateway_profile", profileName, "PUT")

	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	b.pendingOperations[operationID] = &Operation{
		ResourceType:  "gateway_profile",
		ResourceName:  profileName,
		OperationType: "PUT",
		Status:        OperationPending,
	}

	b.operationWaitChannels[operationID] = make(chan struct{})

	now := time.Now()
	b.lastOperationTime = now
	if b.batchStartTime.IsZero() {
		b.batchStartTime = now
	}

	tflog.Debug(ctx, "Added Gateway Profile to PUT batch", map[string]interface{}{
		"profile_name": profileName,
		"batch_size":   len(b.gatewayProfilePut),
		"operation_id": operationID,
	})

	return operationID
}

func (b *BulkOperationManager) AddGatewayProfilePatch(ctx context.Context, profileName string, props openapi.ConfigPutRequestGatewayProfileGatewayProfileName) string {
	b.mutex.Lock()
	b.gatewayProfilePatch[profileName] = props
	b.mutex.Unlock()

	operationID := generateOperationID("gateway_profile", profileName, "PATCH")

	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	b.pendingOperations[operationID] = &Operation{
		ResourceType:  "gateway_profile",
		ResourceName:  profileName,
		OperationType: "PATCH",
		Status:        OperationPending,
	}

	b.operationWaitChannels[operationID] = make(chan struct{})

	now := time.Now()
	b.lastOperationTime = now
	if b.batchStartTime.IsZero() {
		b.batchStartTime = now
	}

	tflog.Debug(ctx, "Added Gateway Profile to PATCH batch", map[string]interface{}{
		"profile_name": profileName,
		"batch_size":   len(b.gatewayProfilePatch),
		"operation_id": operationID,
	})

	return operationID
}

func (b *BulkOperationManager) AddGatewayProfileDelete(ctx context.Context, profileName string) string {
	b.mutex.Lock()
	b.gatewayProfileDelete = append(b.gatewayProfileDelete, profileName)
	b.mutex.Unlock()

	operationID := generateOperationID("gateway_profile", profileName, "DELETE")

	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	b.pendingOperations[operationID] = &Operation{
		ResourceType:  "gateway_profile",
		ResourceName:  profileName,
		OperationType: "DELETE",
		Status:        OperationPending,
	}

	b.operationWaitChannels[operationID] = make(chan struct{})

	now := time.Now()
	b.lastOperationTime = now
	if b.batchStartTime.IsZero() {
		b.batchStartTime = now
	}

	tflog.Debug(ctx, "Added Gateway Profile to DELETE batch", map[string]interface{}{
		"profile_name": profileName,
		"batch_size":   len(b.gatewayProfileDelete),
		"operation_id": operationID,
	})

	return operationID
}

func (b *BulkOperationManager) AddEthPortProfilePut(ctx context.Context, ethPortProfileName string, props openapi.ConfigPutRequestEthPortProfileEthPortProfileName) string {
	b.mutex.Lock()
	b.ethPortProfilePut[ethPortProfileName] = props
	b.mutex.Unlock()

	operationID := generateOperationID("eth_port_profile", ethPortProfileName, "PUT")

	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	b.pendingOperations[operationID] = &Operation{
		ResourceType:  "eth_port_profile",
		ResourceName:  ethPortProfileName,
		OperationType: "PUT",
		Status:        OperationPending,
	}

	b.operationWaitChannels[operationID] = make(chan struct{})

	now := time.Now()
	b.lastOperationTime = now
	if b.batchStartTime.IsZero() {
		b.batchStartTime = now
	}

	tflog.Debug(ctx, "Added EthPortProfile to PUT batch", map[string]interface{}{
		"eth_port_profile_name": ethPortProfileName,
		"batch_size":            len(b.ethPortProfilePut),
		"operation_id":          operationID,
	})

	return operationID
}

func (b *BulkOperationManager) AddEthPortProfilePatch(ctx context.Context, ethPortProfileName string, props openapi.ConfigPutRequestEthPortProfileEthPortProfileName) string {
	b.mutex.Lock()
	b.ethPortProfilePatch[ethPortProfileName] = props
	b.mutex.Unlock()

	operationID := generateOperationID("eth_port_profile", ethPortProfileName, "PATCH")

	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	b.pendingOperations[operationID] = &Operation{
		ResourceType:  "eth_port_profile",
		ResourceName:  ethPortProfileName,
		OperationType: "PATCH",
		Status:        OperationPending,
	}

	b.operationWaitChannels[operationID] = make(chan struct{})

	now := time.Now()
	b.lastOperationTime = now
	if b.batchStartTime.IsZero() {
		b.batchStartTime = now
	}

	tflog.Debug(ctx, "Added EthPortProfile to PATCH batch", map[string]interface{}{
		"eth_port_profile_name": ethPortProfileName,
		"batch_size":            len(b.ethPortProfilePatch),
		"operation_id":          operationID,
	})

	return operationID
}

func (b *BulkOperationManager) AddEthPortProfileDelete(ctx context.Context, ethPortProfileName string) string {
	b.mutex.Lock()
	b.ethPortProfileDelete = append(b.ethPortProfileDelete, ethPortProfileName)
	b.mutex.Unlock()

	operationID := generateOperationID("eth_port_profile", ethPortProfileName, "DELETE")

	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	b.pendingOperations[operationID] = &Operation{
		ResourceType:  "eth_port_profile",
		ResourceName:  ethPortProfileName,
		OperationType: "DELETE",
		Status:        OperationPending,
	}

	b.operationWaitChannels[operationID] = make(chan struct{})

	now := time.Now()
	b.lastOperationTime = now
	if b.batchStartTime.IsZero() {
		b.batchStartTime = now
	}

	tflog.Debug(ctx, "Added EthPortProfile to DELETE batch", map[string]interface{}{
		"eth_port_profile_name": ethPortProfileName,
		"batch_size":            len(b.ethPortProfileDelete),
		"operation_id":          operationID,
	})

	return operationID
}

func (b *BulkOperationManager) AddEthPortSettingsPut(ctx context.Context, ethPortSettingsName string, props openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName) string {
	b.mutex.Lock()
	b.ethPortSettingsPut[ethPortSettingsName] = props
	b.mutex.Unlock()

	operationID := generateOperationID("eth_port_settings", ethPortSettingsName, "PUT")

	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	b.pendingOperations[operationID] = &Operation{
		ResourceType:  "eth_port_settings",
		ResourceName:  ethPortSettingsName,
		OperationType: "PUT",
		Status:        OperationPending,
	}

	b.operationWaitChannels[operationID] = make(chan struct{})

	now := time.Now()
	b.lastOperationTime = now
	if b.batchStartTime.IsZero() {
		b.batchStartTime = now
	}

	tflog.Debug(ctx, "Added EthPortSettings to PUT batch", map[string]interface{}{
		"eth_port_settings_name": ethPortSettingsName,
		"batch_size":             len(b.ethPortSettingsPut),
		"operation_id":           operationID,
	})

	return operationID
}

func (b *BulkOperationManager) AddEthPortSettingsPatch(ctx context.Context, ethPortSettingsName string, props openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName) string {
	b.mutex.Lock()
	b.ethPortSettingsPatch[ethPortSettingsName] = props
	b.mutex.Unlock()

	operationID := generateOperationID("eth_port_settings", ethPortSettingsName, "PATCH")

	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	b.pendingOperations[operationID] = &Operation{
		ResourceType:  "eth_port_settings",
		ResourceName:  ethPortSettingsName,
		OperationType: "PATCH",
		Status:        OperationPending,
	}

	b.operationWaitChannels[operationID] = make(chan struct{})

	now := time.Now()
	b.lastOperationTime = now
	if b.batchStartTime.IsZero() {
		b.batchStartTime = now
	}

	tflog.Debug(ctx, "Added EthPortSettings to PATCH batch", map[string]interface{}{
		"eth_port_settings_name": ethPortSettingsName,
		"batch_size":             len(b.ethPortSettingsPatch),
		"operation_id":           operationID,
	})

	return operationID
}

func (b *BulkOperationManager) AddEthPortSettingsDelete(ctx context.Context, ethPortSettingsName string) string {
	b.mutex.Lock()
	b.ethPortSettingsDelete = append(b.ethPortSettingsDelete, ethPortSettingsName)
	b.mutex.Unlock()

	operationID := generateOperationID("eth_port_settings", ethPortSettingsName, "DELETE")

	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	b.pendingOperations[operationID] = &Operation{
		ResourceType:  "eth_port_settings",
		ResourceName:  ethPortSettingsName,
		OperationType: "DELETE",
		Status:        OperationPending,
	}

	b.operationWaitChannels[operationID] = make(chan struct{})

	now := time.Now()
	b.lastOperationTime = now
	if b.batchStartTime.IsZero() {
		b.batchStartTime = now
	}

	tflog.Debug(ctx, "Added EthPortSettings to DELETE batch", map[string]interface{}{
		"eth_port_settings_name": ethPortSettingsName,
		"batch_size":             len(b.ethPortSettingsDelete),
		"operation_id":           operationID,
	})

	return operationID
}

func (b *BulkOperationManager) ExecuteBulkGatewayPut(ctx context.Context) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	b.mutex.Lock()

	gatewayPut := make(map[string]openapi.ConfigPutRequestGatewayGatewayName)
	for k, v := range b.gatewayPut {
		gatewayPut[k] = v
	}

	b.gatewayPut = make(map[string]openapi.ConfigPutRequestGatewayGatewayName)

	b.mutex.Unlock()

	if len(gatewayPut) == 0 {
		return diagnostics
	}

	gatewayNames := make([]string, 0, len(gatewayPut))
	for name := range gatewayPut {
		gatewayNames = append(gatewayNames, name)
	}

	tflog.Debug(ctx, "Executing bulk gateway PUT operation", map[string]interface{}{
		"gateway_count": len(gatewayPut),
		"gateway_names": gatewayNames,
	})

	putRequest := openapi.NewGatewaysPutRequest()
	gatewayMap := make(map[string]openapi.ConfigPutRequestGatewayGatewayName)

	for name, props := range gatewayPut {
		gatewayMap[name] = props
	}
	putRequest.SetGateway(gatewayMap)
	retryConfig := DefaultRetryConfig()
	var err error

	for retry := 0; retry < retryConfig.MaxRetries; retry++ {
		// Create a separate context for the API call
		apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)

		req := b.client.GatewaysAPI.GatewaysPut(apiCtx).GatewaysPutRequest(*putRequest)
		_, err := req.Execute()

		// Release the API call context
		cancel()

		if err == nil {
			tflog.Debug(ctx, "Bulk gateway PUT operation succeeded", map[string]interface{}{
				"attempt": retry + 1,
			})
			break
		}

		if IsRetriableError(err) {
			delayTime := CalculateBackoff(retry, retryConfig)
			tflog.Debug(ctx, "Bulk gateway PUT operation failed with retriable error, retrying", map[string]interface{}{
				"attempt":     retry + 1,
				"error":       err.Error(),
				"delay_ms":    delayTime.Milliseconds(),
				"max_retries": retryConfig.MaxRetries,
			})

			time.Sleep(delayTime)
			continue
		}

		tflog.Error(ctx, "Bulk gateway PUT operation failed with non-retriable error", map[string]interface{}{
			"error": err.Error(),
		})
		break
	}

	// Update operation statuses based on the result
	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	for opID, op := range b.pendingOperations {
		// Only process gateway PUT operations
		if op.ResourceType == "gateway" && op.OperationType == "PUT" {
			// Check if this operation's gateway name is in our batch
			if _, exists := gatewayPut[op.ResourceName]; exists {
				if err == nil {
					// Mark operation as successful
					op.Status = OperationSucceeded
					b.operationResults[opID] = true

					// Signal waiting resources
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				} else {
					// Mark operation as failed
					op.Status = OperationFailed
					op.Error = err
					b.operationErrors[opID] = err
					b.operationResults[opID] = false

					// Signal waiting resources with the error
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				}
			}
		}
	}

	if err != nil {
		diagnostics.AddError(
			"Failed to execute bulk gateway PUT operation",
			fmt.Sprintf("Error: %s", err),
		)
		return diagnostics
	}

	b.recentGatewayOps = true
	b.recentGatewayOpTime = time.Now()
	return diagnostics
}

func (b *BulkOperationManager) ExecuteBulkGatewayPatch(ctx context.Context) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	b.mutex.Lock()

	gatewayPatch := make(map[string]openapi.ConfigPutRequestGatewayGatewayName)
	for k, v := range b.gatewayPatch {
		gatewayPatch[k] = v
	}

	b.gatewayPatch = make(map[string]openapi.ConfigPutRequestGatewayGatewayName)

	b.mutex.Unlock()

	if len(gatewayPatch) == 0 {
		return diagnostics
	}

	gatewayNames := make([]string, 0, len(gatewayPatch))
	for name := range gatewayPatch {
		gatewayNames = append(gatewayNames, name)
	}

	tflog.Debug(ctx, "Executing bulk gateway PATCH operation", map[string]interface{}{
		"gateway_count": len(gatewayPatch),
		"gateway_names": gatewayNames,
	})

	patchRequest := openapi.NewGatewaysPutRequest()
	gatewayMap := make(map[string]openapi.ConfigPutRequestGatewayGatewayName)

	for name, props := range gatewayPatch {
		gatewayMap[name] = props
	}
	patchRequest.SetGateway(gatewayMap)
	retryConfig := DefaultRetryConfig()
	var err error

	for retry := 0; retry < retryConfig.MaxRetries; retry++ {
		if retry > 0 {
			delay := CalculateBackoff(retry-1, retryConfig)

			tflog.Debug(ctx, "Retrying bulk Gateway PATCH operation after delay", map[string]interface{}{
				"retry": retry,
				"delay": delay,
			})

			time.Sleep(delay)
		}

		// Create a separate context for the API call
		apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)

		req := b.client.GatewaysAPI.GatewaysPatch(apiCtx).GatewaysPutRequest(*patchRequest)
		_, err = req.Execute()

		// Release the API call context
		cancel()

		if err == nil {
			tflog.Debug(ctx, "Bulk Gateway PATCH operation successful", map[string]interface{}{
				"count": len(gatewayPatch),
			})
			break
		}

		if !IsRetriableError(err) {
			tflog.Error(ctx, "Bulk Gateway PATCH operation failed with non-retriable error", map[string]interface{}{
				"error": err.Error(),
			})
			break
		}
	}

	// Update operation statuses based on the result
	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	for opID, op := range b.pendingOperations {
		// Only process gateway PATCH operations
		if op.ResourceType == "gateway" && op.OperationType == "PATCH" {
			// Check if this operation's gateway name is in our batch
			if _, exists := gatewayPatch[op.ResourceName]; exists {
				if err == nil {
					// Mark operation as successful
					op.Status = OperationSucceeded
					b.operationResults[opID] = true

					// Signal waiting resources
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				} else {
					// Mark operation as failed
					op.Status = OperationFailed
					op.Error = err
					b.operationErrors[opID] = err
					b.operationResults[opID] = false

					// Signal waiting resources with the error
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				}
			}
		}
	}

	if err != nil {
		diagnostics.AddError(
			"Failed to execute bulk Gateway PATCH operation",
			fmt.Sprintf("Error: %s", err),
		)
		return diagnostics
	}

	b.recentGatewayOps = true
	b.recentGatewayOpTime = time.Now()
	return diagnostics
}

func (b *BulkOperationManager) ExecuteBulkLagPut(ctx context.Context) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	b.mutex.Lock()

	lagPut := make(map[string]openapi.ConfigPutRequestLagLagName)
	for k, v := range b.lagPut {
		lagPut[k] = v
	}

	b.lagPut = make(map[string]openapi.ConfigPutRequestLagLagName)

	b.mutex.Unlock()

	if len(lagPut) == 0 {
		return diagnostics
	}

	lagNames := make([]string, 0, len(lagPut))
	for name := range lagPut {
		lagNames = append(lagNames, name)
	}

	tflog.Debug(ctx, "Executing bulk LAG PUT operation", map[string]interface{}{
		"lag_count": len(lagPut),
		"lag_names": lagNames,
	})

	putRequest := openapi.NewLagsPutRequest()
	lagMap := make(map[string]openapi.ConfigPutRequestLagLagName)

	for name, props := range lagPut {
		lagMap[name] = props
	}
	putRequest.SetLag(lagMap)
	retryConfig := DefaultRetryConfig()
	var err error

	for retry := 0; retry < retryConfig.MaxRetries; retry++ {
		// Create a separate context for the API call
		apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)

		req := b.client.LAGsAPI.LagsPut(apiCtx).LagsPutRequest(*putRequest)
		_, err = req.Execute()

		// Release the API call context
		cancel()

		if err == nil {
			tflog.Debug(ctx, "Bulk LAG PUT operation succeeded", map[string]interface{}{
				"attempt": retry + 1,
			})
			break
		}

		if IsRetriableError(err) {
			delayTime := CalculateBackoff(retry, retryConfig)
			tflog.Debug(ctx, "Bulk LAG PUT operation failed with retriable error, retrying", map[string]interface{}{
				"attempt":     retry + 1,
				"error":       err.Error(),
				"delay_ms":    delayTime.Milliseconds(),
				"max_retries": retryConfig.MaxRetries,
			})

			time.Sleep(delayTime)
			continue
		}

		tflog.Error(ctx, "Bulk LAG PUT operation failed with non-retriable error", map[string]interface{}{
			"error": err.Error(),
		})
		break
	}

	// Update operation statuses based on the result
	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	for opID, op := range b.pendingOperations {
		// Only process LAG PUT operations
		if op.ResourceType == "lag" && op.OperationType == "PUT" {
			// Check if this operation's lag name is in our batch
			if _, exists := lagPut[op.ResourceName]; exists {
				if err == nil {
					// Mark operation as successful
					op.Status = OperationSucceeded
					b.operationResults[opID] = true

					// Signal waiting resources
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				} else {
					// Mark operation as failed
					op.Status = OperationFailed
					op.Error = err
					b.operationErrors[opID] = err
					b.operationResults[opID] = false

					// Signal waiting resources with the error
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				}
			}
		}
	}

	if err != nil {
		diagnostics.AddError(
			"Failed to execute bulk LAG PUT operation",
			fmt.Sprintf("Error: %s", err),
		)
		return diagnostics
	}

	b.recentLagOps = true
	b.recentLagOpTime = time.Now()
	return diagnostics
}

func (b *BulkOperationManager) ExecuteBulkLagPatch(ctx context.Context) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	b.mutex.Lock()

	lagPatch := make(map[string]openapi.ConfigPutRequestLagLagName)
	for k, v := range b.lagPatch {
		lagPatch[k] = v
	}

	b.lagPatch = make(map[string]openapi.ConfigPutRequestLagLagName)

	b.mutex.Unlock()

	if len(lagPatch) == 0 {
		return diagnostics
	}

	patchRequest := openapi.NewLagsPutRequest()
	lagMap := make(map[string]openapi.ConfigPutRequestLagLagName)

	for name, props := range lagPatch {
		lagMap[name] = props
	}
	patchRequest.SetLag(lagMap)
	retryConfig := DefaultRetryConfig()
	var err error

	for retry := 0; retry < retryConfig.MaxRetries; retry++ {
		if retry > 0 {
			delay := CalculateBackoff(retry-1, retryConfig)

			tflog.Debug(ctx, "Retrying bulk LAG PATCH operation after delay", map[string]interface{}{
				"retry": retry,
				"delay": delay,
			})

			time.Sleep(delay)
		}

		// Create a separate context for the API call
		apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)

		req := b.client.LAGsAPI.LagsPatch(apiCtx).LagsPutRequest(*patchRequest)
		_, err = req.Execute()

		// Release the API call context
		cancel()

		if err == nil {
			tflog.Debug(ctx, "Bulk LAG PATCH operation successful", map[string]interface{}{
				"count": len(lagPatch),
			})
			break
		}

		if !IsRetriableError(err) {
			tflog.Error(ctx, "Bulk LAG PATCH operation failed with non-retriable error", map[string]interface{}{
				"error": err.Error(),
			})
			break
		}
	}

	// Update operation statuses based on the result
	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	for opID, op := range b.pendingOperations {
		// Only process LAG PATCH operations
		if op.ResourceType == "lag" && op.OperationType == "PATCH" {
			// Check if this operation's lag name is in our batch
			if _, exists := lagPatch[op.ResourceName]; exists {
				if err == nil {
					// Mark operation as successful
					op.Status = OperationSucceeded
					b.operationResults[opID] = true

					// Signal waiting resources
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				} else {
					// Mark operation as failed
					op.Status = OperationFailed
					op.Error = err
					b.operationErrors[opID] = err
					b.operationResults[opID] = false

					// Signal waiting resources with the error
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				}
			}
		}
	}

	if err != nil {
		diagnostics.AddError(
			"Failed to execute bulk LAG PATCH operation",
			fmt.Sprintf("Error: %s", err),
		)
		return diagnostics
	}

	b.recentLagOps = true
	b.recentLagOpTime = time.Now()
	return diagnostics
}

func (b *BulkOperationManager) ExecuteBulkLagDelete(ctx context.Context) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	b.mutex.Lock()

	lagNames := make([]string, len(b.lagDelete))
	copy(lagNames, b.lagDelete)

	lagDeleteMap := make(map[string]bool)
	for _, name := range lagNames {
		lagDeleteMap[name] = true
	}

	b.lagDelete = make([]string, 0)

	b.mutex.Unlock()

	if len(lagNames) == 0 {
		return diagnostics
	}

	retryConfig := DefaultRetryConfig()
	var err error

	for retry := 0; retry < retryConfig.MaxRetries; retry++ {
		if retry > 0 {
			delay := CalculateBackoff(retry-1, retryConfig)

			tflog.Debug(ctx, "Retrying bulk LAG DELETE operation after delay", map[string]interface{}{
				"retry": retry,
				"delay": delay,
			})

			time.Sleep(delay)
		}

		// Create a separate context for the API call
		apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)

		req := b.client.LAGsAPI.LagsDelete(apiCtx).LagName(lagNames)
		_, err = req.Execute()

		// Release the API call context
		cancel()

		if err == nil {
			tflog.Debug(ctx, "Bulk LAG DELETE operation successful", map[string]interface{}{
				"count": len(lagNames),
			})
			break
		}

		if !IsRetriableError(err) {
			tflog.Error(ctx, "Bulk LAG DELETE operation failed with non-retriable error", map[string]interface{}{
				"error": err.Error(),
			})
			break
		}
	}

	// Update operation statuses based on the result
	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	for opID, op := range b.pendingOperations {
		// Only process LAG DELETE operations
		if op.ResourceType == "lag" && op.OperationType == "DELETE" {
			// Check if this operation's lag name is in our batch
			if _, exists := lagDeleteMap[op.ResourceName]; exists {
				if err == nil {
					// Mark operation as successful
					op.Status = OperationSucceeded
					b.operationResults[opID] = true

					// Signal waiting resources
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				} else {
					// Mark operation as failed
					op.Status = OperationFailed
					op.Error = err
					b.operationErrors[opID] = err
					b.operationResults[opID] = false

					// Signal waiting resources with the error
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				}
			}
		}
	}

	if err != nil {
		diagnostics.AddError(
			"Failed to execute bulk LAG DELETE operation",
			fmt.Sprintf("Error: %s", err),
		)
		return diagnostics
	}

	b.recentLagOps = true
	b.recentLagOpTime = time.Now()
	return diagnostics
}

func (b *BulkOperationManager) ExecuteBulkTenantPut(ctx context.Context) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	b.mutex.Lock()

	tenantPut := make(map[string]openapi.ConfigPutRequestTenantTenantName)
	for k, v := range b.tenantPut {
		tenantPut[k] = v
	}

	b.tenantPut = make(map[string]openapi.ConfigPutRequestTenantTenantName)

	b.mutex.Unlock()

	if len(tenantPut) == 0 {
		return diagnostics
	}

	tenantNames := make([]string, 0, len(tenantPut))
	for name := range tenantPut {
		tenantNames = append(tenantNames, name)
	}

	tflog.Debug(ctx, "Executing bulk Tenant PUT operation", map[string]interface{}{
		"tenant_count": len(tenantPut),
		"tenant_names": tenantNames,
	})

	putRequest := openapi.NewTenantsPutRequest()
	tenantMap := make(map[string]openapi.ConfigPutRequestTenantTenantName)

	for name, props := range tenantPut {
		tenantMap[name] = props
	}
	putRequest.SetTenant(tenantMap)
	retryConfig := DefaultRetryConfig()
	var err error

	for retry := 0; retry < retryConfig.MaxRetries; retry++ {
		// Create a separate context for the API call
		apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)

		req := b.client.TenantsAPI.TenantsPut(apiCtx).TenantsPutRequest(*putRequest)
		_, err = req.Execute()

		// Release the API call context
		cancel()

		if err == nil {
			tflog.Debug(ctx, "Bulk Tenant PUT operation succeeded", map[string]interface{}{
				"attempt": retry + 1,
			})
			break
		}

		if IsRetriableError(err) {
			delayTime := CalculateBackoff(retry, retryConfig)
			tflog.Debug(ctx, "Bulk Tenant PUT operation failed with retriable error, retrying", map[string]interface{}{
				"attempt":     retry + 1,
				"error":       err.Error(),
				"delay_ms":    delayTime.Milliseconds(),
				"max_retries": retryConfig.MaxRetries,
			})

			time.Sleep(delayTime)
			continue
		}

		tflog.Error(ctx, "Bulk Tenant PUT operation failed with non-retriable error", map[string]interface{}{
			"error": err.Error(),
		})
		break
	}

	// Update operation statuses based on the result
	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	for opID, op := range b.pendingOperations {
		// Only process Tenant PUT operations
		if op.ResourceType == "tenant" && op.OperationType == "PUT" {
			// Check if this operation's tenant name is in our batch
			if _, exists := tenantPut[op.ResourceName]; exists {
				if err == nil {
					// Mark operation as successful
					op.Status = OperationSucceeded
					b.operationResults[opID] = true

					// Signal waiting resources
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				} else {
					// Mark operation as failed
					op.Status = OperationFailed
					op.Error = err
					b.operationErrors[opID] = err
					b.operationResults[opID] = false

					// Signal waiting resources with the error
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				}
			}
		}
	}

	if err != nil {
		diagnostics.AddError(
			"Failed to execute bulk Tenant PUT operation",
			fmt.Sprintf("Error: %s", err),
		)
		return diagnostics
	}

	b.recentTenantOps = true
	b.recentTenantOpTime = time.Now()
	return diagnostics
}

func (b *BulkOperationManager) ExecuteBulkTenantPatch(ctx context.Context) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	b.mutex.Lock()

	tenantPatch := make(map[string]openapi.ConfigPutRequestTenantTenantName)
	for k, v := range b.tenantPatch {
		tenantPatch[k] = v
	}

	b.tenantPatch = make(map[string]openapi.ConfigPutRequestTenantTenantName)

	b.mutex.Unlock()

	if len(tenantPatch) == 0 {
		return diagnostics
	}

	tenantNames := make([]string, 0, len(tenantPatch))
	for name := range tenantPatch {
		tenantNames = append(tenantNames, name)
	}

	tflog.Debug(ctx, "Executing bulk Tenant PATCH operation", map[string]interface{}{
		"tenant_count": len(tenantPatch),
		"tenant_names": tenantNames,
	})

	patchRequest := openapi.NewTenantsPutRequest()
	tenantMap := make(map[string]openapi.ConfigPutRequestTenantTenantName)

	for name, props := range tenantPatch {
		tenantMap[name] = props
	}
	patchRequest.SetTenant(tenantMap)
	retryConfig := DefaultRetryConfig()
	var err error

	for retry := 0; retry < retryConfig.MaxRetries; retry++ {
		if retry > 0 {
			delay := CalculateBackoff(retry-1, retryConfig)

			tflog.Debug(ctx, "Retrying bulk Tenant PATCH operation after delay", map[string]interface{}{
				"retry": retry,
				"delay": delay,
			})

			time.Sleep(delay)
		}

		// Create a separate context for the API call
		apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)

		req := b.client.TenantsAPI.TenantsPatch(apiCtx).TenantsPutRequest(*patchRequest)
		_, err = req.Execute()

		// Release the API call context
		cancel()

		if err == nil {
			tflog.Debug(ctx, "Bulk Tenant PATCH operation successful", map[string]interface{}{
				"count": len(tenantPatch),
			})
			break
		}

		if !IsRetriableError(err) {
			tflog.Error(ctx, "Bulk Tenant PATCH operation failed with non-retriable error", map[string]interface{}{
				"error": err.Error(),
			})
			break
		}
	}

	// Update operation statuses based on the result
	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	for opID, op := range b.pendingOperations {
		// Only process Tenant PATCH operations
		if op.ResourceType == "tenant" && op.OperationType == "PATCH" {
			// Check if this operation's tenant name is in our batch
			if _, exists := tenantPatch[op.ResourceName]; exists {
				if err == nil {
					// Mark operation as successful
					op.Status = OperationSucceeded
					b.operationResults[opID] = true

					// Signal waiting resources
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				} else {
					// Mark operation as failed
					op.Status = OperationFailed
					op.Error = err
					b.operationErrors[opID] = err
					b.operationResults[opID] = false

					// Signal waiting resources with the error
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				}
			}
		}
	}

	if err != nil {
		diagnostics.AddError(
			"Failed to execute bulk Tenant PATCH operation",
			fmt.Sprintf("Error: %s", err),
		)
		return diagnostics
	}

	b.recentTenantOps = true
	b.recentTenantOpTime = time.Now()
	return diagnostics
}

func (b *BulkOperationManager) ExecuteBulkTenantDelete(ctx context.Context) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	b.mutex.Lock()

	tenantNames := make([]string, len(b.tenantDelete))
	copy(tenantNames, b.tenantDelete)

	tenantDeleteMap := make(map[string]bool)
	for _, name := range tenantNames {
		tenantDeleteMap[name] = true
	}

	b.tenantDelete = make([]string, 0)

	b.mutex.Unlock()

	if len(tenantNames) == 0 {
		return diagnostics
	}

	retryConfig := DefaultRetryConfig()
	var err error

	for retry := 0; retry < retryConfig.MaxRetries; retry++ {
		if retry > 0 {
			delay := CalculateBackoff(retry-1, retryConfig)

			tflog.Debug(ctx, "Retrying bulk Tenant DELETE operation after delay", map[string]interface{}{
				"retry":       retry,
				"delay":       delay,
				"max_retries": retryConfig.MaxRetries,
			})

			time.Sleep(delay)
		}

		tflog.Debug(ctx, "Executing bulk Tenant DELETE operation", map[string]interface{}{
			"tenant_count": len(tenantNames),
			"tenant_names": tenantNames,
		})

		// Create a separate context for the API call
		apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)

		req := b.client.TenantsAPI.TenantsDelete(apiCtx).TenantName(tenantNames)
		_, err = req.Execute()

		// Release the API call context
		cancel()

		if err == nil {
			tflog.Debug(ctx, "Bulk Tenant DELETE operation successful", map[string]interface{}{
				"count": len(tenantNames),
			})
			break
		}

		if !IsRetriableError(err) {
			tflog.Error(ctx, "Bulk Tenant DELETE operation failed with non-retriable error", map[string]interface{}{
				"error": err.Error(),
			})
			break
		}
	}

	// Update operation statuses based on the result
	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	for opID, op := range b.pendingOperations {
		// Only process Tenant DELETE operations
		if op.ResourceType == "tenant" && op.OperationType == "DELETE" {
			// Check if this operation's tenant name is in our batch
			if _, exists := tenantDeleteMap[op.ResourceName]; exists {
				if err == nil {
					// Mark operation as successful
					op.Status = OperationSucceeded
					b.operationResults[opID] = true

					// Signal waiting resources
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				} else {
					// Mark operation as failed
					op.Status = OperationFailed
					op.Error = err
					b.operationErrors[opID] = err
					b.operationResults[opID] = false

					// Signal waiting resources with the error
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				}
			}
		}
	}

	if err != nil {
		diagnostics.AddError(
			"Failed to execute bulk Tenant DELETE operation",
			fmt.Sprintf("Error: %s", err),
		)
		return diagnostics
	}

	b.recentTenantOps = true
	b.recentTenantOpTime = time.Now()
	return diagnostics
}

func (b *BulkOperationManager) ExecuteBulkGatewayDelete(ctx context.Context) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	b.mutex.Lock()

	gatewayNames := make([]string, len(b.gatewayDelete))
	copy(gatewayNames, b.gatewayDelete)

	gatewayDeleteMap := make(map[string]bool)
	for _, name := range gatewayNames {
		gatewayDeleteMap[name] = true
	}

	b.gatewayDelete = make([]string, 0)

	b.mutex.Unlock()

	if len(gatewayNames) == 0 {
		return diagnostics
	}

	retryConfig := DefaultRetryConfig()
	var err error

	for retry := 0; retry < retryConfig.MaxRetries; retry++ {
		if retry > 0 {
			delay := CalculateBackoff(retry-1, retryConfig)

			tflog.Debug(ctx, "Retrying bulk Gateway DELETE operation after delay", map[string]interface{}{
				"retry":       retry,
				"delay":       delay,
				"max_retries": retryConfig.MaxRetries,
			})

			time.Sleep(delay)
		}

		tflog.Debug(ctx, "Executing bulk Gateway DELETE operation", map[string]interface{}{
			"gateway_count": len(gatewayNames),
			"gateway_names": gatewayNames,
		})

		// Create a separate context for the API call
		apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)

		req := b.client.GatewaysAPI.GatewaysDelete(apiCtx).GatewayName(gatewayNames)
		_, err = req.Execute()

		// Release the API call context
		cancel()

		if err == nil {
			tflog.Debug(ctx, "Bulk Gateway DELETE operation successful", map[string]interface{}{
				"count": len(gatewayNames),
			})
			break
		}

		if !IsRetriableError(err) {
			tflog.Error(ctx, "Bulk Gateway DELETE operation failed with non-retriable error", map[string]interface{}{
				"error": err.Error(),
			})
			break
		}
	}

	// Update operation statuses based on the result
	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	for opID, op := range b.pendingOperations {
		// Only process gateway DELETE operations
		if op.ResourceType == "gateway" && op.OperationType == "DELETE" {
			// Check if this operation's gateway name is in our batch
			if _, exists := gatewayDeleteMap[op.ResourceName]; exists {
				if err == nil {
					// Mark operation as successful
					op.Status = OperationSucceeded
					b.operationResults[opID] = true

					// Signal waiting resources
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				} else {
					// Mark operation as failed
					op.Status = OperationFailed
					op.Error = err
					b.operationErrors[opID] = err
					b.operationResults[opID] = false

					// Signal waiting resources with the error
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				}
			}
		}
	}

	if err != nil {
		diagnostics.AddError(
			"Failed to execute bulk Gateway DELETE operation",
			fmt.Sprintf("Error: %s", err),
		)
		return diagnostics
	}

	b.recentGatewayOps = true
	b.recentGatewayOpTime = time.Now()
	return diagnostics
}

func (b *BulkOperationManager) ExecuteBulkServicePut(ctx context.Context) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	b.mutex.Lock()

	servicePut := make(map[string]openapi.ConfigPutRequestServiceServiceName)
	for k, v := range b.servicePut {
		servicePut[k] = v
	}

	b.servicePut = make(map[string]openapi.ConfigPutRequestServiceServiceName)

	b.mutex.Unlock()

	if len(servicePut) == 0 {
		return diagnostics
	}

	serviceNames := make([]string, 0, len(servicePut))
	for name := range servicePut {
		serviceNames = append(serviceNames, name)
	}

	tflog.Debug(ctx, "Executing bulk Service PUT operation", map[string]interface{}{
		"service_count": len(servicePut),
		"service_names": serviceNames,
	})

	putRequest := openapi.NewServicesPutRequest()
	serviceMap := make(map[string]openapi.ConfigPutRequestServiceServiceName)

	for name, props := range servicePut {
		serviceMap[name] = props
	}
	putRequest.SetService(serviceMap)
	retryConfig := DefaultRetryConfig()
	var err error

	for retry := 0; retry < retryConfig.MaxRetries; retry++ {
		// Create a separate context for the API call
		apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)

		req := b.client.ServicesAPI.ServicesPut(apiCtx).ServicesPutRequest(*putRequest)
		_, err = req.Execute()

		// Release the API call context
		cancel()

		if err == nil {
			tflog.Debug(ctx, "Bulk Service PUT operation succeeded", map[string]interface{}{
				"attempt": retry + 1,
			})
			break
		}

		if IsRetriableError(err) {
			delayTime := CalculateBackoff(retry, retryConfig)
			tflog.Debug(ctx, "Bulk Service PUT operation failed with retriable error, retrying", map[string]interface{}{
				"attempt":     retry + 1,
				"error":       err.Error(),
				"delay_ms":    delayTime.Milliseconds(),
				"max_retries": retryConfig.MaxRetries,
			})

			time.Sleep(delayTime)
			continue
		}

		tflog.Error(ctx, "Bulk Service PUT operation failed with non-retriable error", map[string]interface{}{
			"error": err.Error(),
		})
		break
	}

	// Update operation statuses based on the result
	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	for opID, op := range b.pendingOperations {
		// Only process Service PUT operations
		if op.ResourceType == "service" && op.OperationType == "PUT" {
			// Check if this operation's service name is in our batch
			if _, exists := servicePut[op.ResourceName]; exists {
				if err == nil {
					// Mark operation as successful
					op.Status = OperationSucceeded
					b.operationResults[opID] = true

					// Signal waiting resources
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				} else {
					// Mark operation as failed
					op.Status = OperationFailed
					op.Error = err
					b.operationErrors[opID] = err
					b.operationResults[opID] = false

					// Signal waiting resources with the error
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				}
			}
		}
	}

	if err != nil {
		diagnostics.AddError(
			"Failed to execute bulk Service PUT operation",
			fmt.Sprintf("Error: %s", err),
		)
		return diagnostics
	}

	b.recentServiceOps = true
	b.recentServiceOpTime = time.Now()
	return diagnostics
}

func (b *BulkOperationManager) ExecuteBulkServicePatch(ctx context.Context) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	b.mutex.Lock()

	servicePatch := make(map[string]openapi.ConfigPutRequestServiceServiceName)
	for k, v := range b.servicePatch {
		servicePatch[k] = v
	}

	b.servicePatch = make(map[string]openapi.ConfigPutRequestServiceServiceName)

	b.mutex.Unlock()

	if len(servicePatch) == 0 {
		return diagnostics
	}

	serviceNames := make([]string, 0, len(servicePatch))
	for name := range servicePatch {
		serviceNames = append(serviceNames, name)
	}

	tflog.Debug(ctx, "Executing bulk Service PATCH operation", map[string]interface{}{
		"service_count": len(servicePatch),
		"service_names": serviceNames,
	})

	patchRequest := openapi.NewServicesPutRequest()
	serviceMap := make(map[string]openapi.ConfigPutRequestServiceServiceName)

	for name, props := range servicePatch {
		serviceMap[name] = props
	}
	patchRequest.SetService(serviceMap)
	retryConfig := DefaultRetryConfig()
	var err error

	for retry := 0; retry < retryConfig.MaxRetries; retry++ {
		if retry > 0 {
			delay := CalculateBackoff(retry-1, retryConfig)

			tflog.Debug(ctx, "Retrying bulk Service PATCH operation after delay", map[string]interface{}{
				"retry": retry,
				"delay": delay,
			})

			time.Sleep(delay)
		}

		// Create a separate context for the API call
		apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)

		req := b.client.ServicesAPI.ServicesPatch(apiCtx).ServicesPutRequest(*patchRequest)
		_, err = req.Execute()

		// Release the API call context
		cancel()

		if err == nil {
			tflog.Debug(ctx, "Bulk Service PATCH operation successful", map[string]interface{}{
				"count": len(servicePatch),
			})
			break
		}

		if !IsRetriableError(err) {
			tflog.Error(ctx, "Bulk Service PATCH operation failed with non-retriable error", map[string]interface{}{
				"error": err.Error(),
			})
			break
		}
	}

	// Update operation statuses based on the result
	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	for opID, op := range b.pendingOperations {
		// Only process Service PATCH operations
		if op.ResourceType == "service" && op.OperationType == "PATCH" {
			// Check if this operation's service name is in our batch
			if _, exists := servicePatch[op.ResourceName]; exists {
				if err == nil {
					// Mark operation as successful
					op.Status = OperationSucceeded
					b.operationResults[opID] = true

					// Signal waiting resources
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				} else {
					// Mark operation as failed
					op.Status = OperationFailed
					op.Error = err
					b.operationErrors[opID] = err
					b.operationResults[opID] = false

					// Signal waiting resources with the error
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				}
			}
		}
	}

	if err != nil {
		diagnostics.AddError(
			"Failed to execute bulk Service PATCH operation",
			fmt.Sprintf("Error: %s", err),
		)
		return diagnostics
	}

	b.recentServiceOps = true
	b.recentServiceOpTime = time.Now()
	return diagnostics
}

func (b *BulkOperationManager) ExecuteBulkServiceDelete(ctx context.Context) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	b.mutex.Lock()

	serviceNames := make([]string, len(b.serviceDelete))
	copy(serviceNames, b.serviceDelete)

	serviceDeleteMap := make(map[string]bool)
	for _, name := range serviceNames {
		serviceDeleteMap[name] = true
	}

	b.serviceDelete = make([]string, 0)

	b.mutex.Unlock()

	if len(serviceNames) == 0 {
		return diagnostics
	}

	retryConfig := DefaultRetryConfig()
	var err error

	for retry := 0; retry < retryConfig.MaxRetries; retry++ {
		if retry > 0 {
			delay := CalculateBackoff(retry-1, retryConfig)

			tflog.Debug(ctx, "Retrying bulk Service DELETE operation after delay", map[string]interface{}{
				"retry":       retry,
				"delay":       delay,
				"max_retries": retryConfig.MaxRetries,
			})

			time.Sleep(delay)
		}

		tflog.Debug(ctx, "Executing bulk Service DELETE operation", map[string]interface{}{
			"service_count": len(serviceNames),
			"service_names": serviceNames,
		})

		// Create a separate context for the API call
		apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)

		req := b.client.ServicesAPI.ServicesDelete(apiCtx).ServiceName(serviceNames)
		_, err = req.Execute()

		// Release the API call context
		cancel()

		if err == nil {
			tflog.Debug(ctx, "Bulk Service DELETE operation successful", map[string]interface{}{
				"count": len(serviceNames),
			})
			break
		}

		if !IsRetriableError(err) {
			tflog.Error(ctx, "Bulk Service DELETE operation failed with non-retriable error", map[string]interface{}{
				"error": err.Error(),
			})
			break
		}
	}

	// Update operation statuses based on the result
	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	for opID, op := range b.pendingOperations {
		// Only process Service DELETE operations
		if op.ResourceType == "service" && op.OperationType == "DELETE" {
			// Check if this operation's service name is in our batch
			if _, exists := serviceDeleteMap[op.ResourceName]; exists {
				if err == nil {
					// Mark operation as successful
					op.Status = OperationSucceeded
					b.operationResults[opID] = true

					// Signal waiting resources
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				} else {
					// Mark operation as failed
					op.Status = OperationFailed
					op.Error = err
					b.operationErrors[opID] = err
					b.operationResults[opID] = false

					// Signal waiting resources with the error
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				}
			}
		}
	}

	if err != nil {
		diagnostics.AddError(
			"Failed to execute bulk Service DELETE operation",
			fmt.Sprintf("Error: %s", err),
		)
		return diagnostics
	}

	b.recentServiceOps = true
	b.recentServiceOpTime = time.Now()
	return diagnostics
}

func (b *BulkOperationManager) ExecuteBulkGatewayProfilePut(ctx context.Context) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	b.mutex.Lock()

	gatewayProfilePut := make(map[string]openapi.ConfigPutRequestGatewayProfileGatewayProfileName)
	for k, v := range b.gatewayProfilePut {
		gatewayProfilePut[k] = v
	}

	b.gatewayProfilePut = make(map[string]openapi.ConfigPutRequestGatewayProfileGatewayProfileName)

	b.mutex.Unlock()

	if len(gatewayProfilePut) == 0 {
		return diagnostics
	}

	gatewayProfileNames := make([]string, 0, len(gatewayProfilePut))
	for name := range gatewayProfilePut {
		gatewayProfileNames = append(gatewayProfileNames, name)
	}

	tflog.Debug(ctx, "Executing bulk Gateway Profile PUT operation", map[string]interface{}{
		"gateway_profile_count": len(gatewayProfilePut),
		"gateway_profile_names": gatewayProfileNames,
	})

	putRequest := openapi.NewGatewayprofilesPutRequest()
	gatewayProfileMap := make(map[string]openapi.ConfigPutRequestGatewayProfileGatewayProfileName)

	for name, props := range gatewayProfilePut {
		gatewayProfileMap[name] = props
	}
	putRequest.SetGatewayProfile(gatewayProfileMap)
	retryConfig := DefaultRetryConfig()
	var err error

	for retry := 0; retry < retryConfig.MaxRetries; retry++ {
		// Create a separate context for the API call
		apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)

		req := b.client.GatewayProfilesAPI.GatewayprofilesPut(apiCtx).GatewayprofilesPutRequest(*putRequest)
		_, err = req.Execute()

		// Release the API call context
		cancel()

		if err == nil {
			tflog.Debug(ctx, "Bulk Gateway Profile PUT operation succeeded", map[string]interface{}{
				"attempt": retry + 1,
			})
			break
		}

		if IsRetriableError(err) {
			delayTime := CalculateBackoff(retry, retryConfig)
			tflog.Debug(ctx, "Bulk Gateway Profile PUT operation failed with retriable error, retrying", map[string]interface{}{
				"attempt":     retry + 1,
				"error":       err.Error(),
				"delay_ms":    delayTime.Milliseconds(),
				"max_retries": retryConfig.MaxRetries,
			})

			time.Sleep(delayTime)
			continue
		}

		tflog.Error(ctx, "Bulk Gateway Profile PUT operation failed with non-retriable error", map[string]interface{}{
			"error": err.Error(),
		})
		break
	}

	// Update operation statuses based on the result
	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	for opID, op := range b.pendingOperations {
		// Only process Gateway Profile PUT operations
		if op.ResourceType == "gateway_profile" && op.OperationType == "PUT" {
			// Check if this operation's gateway profile name is in our batch
			if _, exists := gatewayProfilePut[op.ResourceName]; exists {
				if err == nil {
					// Mark operation as successful
					op.Status = OperationSucceeded
					b.operationResults[opID] = true

					// Signal waiting resources
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				} else {
					// Mark operation as failed
					op.Status = OperationFailed
					op.Error = err
					b.operationErrors[opID] = err
					b.operationResults[opID] = false

					// Signal waiting resources with the error
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				}
			}
		}
	}

	if err != nil {
		diagnostics.AddError(
			"Failed to execute bulk Gateway Profile PUT operation",
			fmt.Sprintf("Error: %s", err),
		)
		return diagnostics
	}

	b.recentGatewayProfileOps = true
	b.recentGatewayProfileOpTime = time.Now()
	return diagnostics
}

func (b *BulkOperationManager) ExecuteBulkGatewayProfilePatch(ctx context.Context) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	b.mutex.Lock()

	gatewayProfilePatch := make(map[string]openapi.ConfigPutRequestGatewayProfileGatewayProfileName)
	for k, v := range b.gatewayProfilePatch {
		gatewayProfilePatch[k] = v
	}

	b.gatewayProfilePatch = make(map[string]openapi.ConfigPutRequestGatewayProfileGatewayProfileName)

	b.mutex.Unlock()

	if len(gatewayProfilePatch) == 0 {
		return diagnostics
	}

	gatewayProfileNames := make([]string, 0, len(gatewayProfilePatch))
	for name := range gatewayProfilePatch {
		gatewayProfileNames = append(gatewayProfileNames, name)
	}

	tflog.Debug(ctx, "Executing bulk Gateway Profile PATCH operation", map[string]interface{}{
		"gateway_profile_count": len(gatewayProfilePatch),
		"gateway_profile_names": gatewayProfileNames,
	})

	patchRequest := openapi.NewGatewayprofilesPutRequest()
	gatewayProfileMap := make(map[string]openapi.ConfigPutRequestGatewayProfileGatewayProfileName)

	for name, props := range gatewayProfilePatch {
		gatewayProfileMap[name] = props
	}
	patchRequest.SetGatewayProfile(gatewayProfileMap)
	retryConfig := DefaultRetryConfig()
	var err error

	for retry := 0; retry < retryConfig.MaxRetries; retry++ {
		if retry > 0 {
			delay := CalculateBackoff(retry-1, retryConfig)

			tflog.Debug(ctx, "Retrying bulk Gateway Profile PATCH operation after delay", map[string]interface{}{
				"retry": retry,
				"delay": delay,
			})

			time.Sleep(delay)
		}

		// Create a separate context for the API call
		apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)

		req := b.client.GatewayProfilesAPI.GatewayprofilesPatch(apiCtx).GatewayprofilesPutRequest(*patchRequest)
		_, err = req.Execute()

		// Release the API call context
		cancel()

		if err == nil {
			tflog.Debug(ctx, "Bulk Gateway Profile PATCH operation successful", map[string]interface{}{
				"count": len(gatewayProfilePatch),
			})
			break
		}

		if !IsRetriableError(err) {
			tflog.Error(ctx, "Bulk Gateway Profile PATCH operation failed with non-retriable error", map[string]interface{}{
				"error": err.Error(),
			})
			break
		}
	}

	// Update operation statuses based on the result
	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	for opID, op := range b.pendingOperations {
		// Only process Gateway Profile PATCH operations
		if op.ResourceType == "gateway_profile" && op.OperationType == "PATCH" {
			// Check if this operation's gateway profile name is in our batch
			if _, exists := gatewayProfilePatch[op.ResourceName]; exists {
				if err == nil {
					// Mark operation as successful
					op.Status = OperationSucceeded
					b.operationResults[opID] = true

					// Signal waiting resources
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				} else {
					// Mark operation as failed
					op.Status = OperationFailed
					op.Error = err
					b.operationErrors[opID] = err
					b.operationResults[opID] = false

					// Signal waiting resources with the error
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				}
			}
		}
	}

	if err != nil {
		diagnostics.AddError(
			"Failed to execute bulk Gateway Profile PATCH operation",
			fmt.Sprintf("Error: %s", err),
		)
		return diagnostics
	}

	b.recentGatewayProfileOps = true
	b.recentGatewayProfileOpTime = time.Now()
	return diagnostics
}

func (b *BulkOperationManager) ExecuteBulkGatewayProfileDelete(ctx context.Context) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	b.mutex.Lock()

	gatewayProfileNames := make([]string, len(b.gatewayProfileDelete))
	copy(gatewayProfileNames, b.gatewayProfileDelete)

	gatewayProfileDeleteMap := make(map[string]bool)
	for _, name := range gatewayProfileNames {
		gatewayProfileDeleteMap[name] = true
	}

	// Clear the pending operations list early to avoid duplicates
	b.gatewayProfileDelete = make([]string, 0)

	b.mutex.Unlock()

	if len(gatewayProfileNames) == 0 {
		return diagnostics
	}

	retryConfig := DefaultRetryConfig()
	var err error

	for retry := 0; retry < retryConfig.MaxRetries; retry++ {
		if retry > 0 {
			delay := CalculateBackoff(retry-1, retryConfig)

			tflog.Debug(ctx, "Retrying bulk Gateway Profile DELETE operation after delay", map[string]interface{}{
				"retry":       retry,
				"delay":       delay,
				"max_retries": retryConfig.MaxRetries,
			})

			time.Sleep(delay)
		}

		tflog.Debug(ctx, "Executing bulk Gateway Profile DELETE operation", map[string]interface{}{
			"gateway_profile_count": len(gatewayProfileNames),
			"gateway_profile_names": gatewayProfileNames,
		})

		// Create a separate context for the API call
		apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)

		req := b.client.GatewayProfilesAPI.GatewayprofilesDelete(apiCtx).ProfileName(gatewayProfileNames)
		_, err = req.Execute()

		// Release the API call context
		cancel()

		if err == nil {
			tflog.Debug(ctx, "Bulk Gateway Profile DELETE operation successful", map[string]interface{}{
				"count": len(gatewayProfileNames),
			})
			break
		}

		if !IsRetriableError(err) {
			tflog.Error(ctx, "Bulk Gateway Profile DELETE operation failed with non-retriable error", map[string]interface{}{
				"error": err.Error(),
			})
			break
		}
	}

	// Update operation statuses based on the result
	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	for opID, op := range b.pendingOperations {
		// Only process Gateway Profile DELETE operations
		if op.ResourceType == "gateway_profile" && op.OperationType == "DELETE" {
			// Check if this operation's gateway profile name is in our batch
			if _, exists := gatewayProfileDeleteMap[op.ResourceName]; exists {
				if err == nil {
					// Mark operation as successful
					op.Status = OperationSucceeded
					b.operationResults[opID] = true

					// Signal waiting resources
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				} else {
					// Mark operation as failed
					op.Status = OperationFailed
					op.Error = err
					b.operationErrors[opID] = err
					b.operationResults[opID] = false

					// Signal waiting resources with the error
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				}
			}
		}
	}

	if err != nil {
		diagnostics.AddError(
			"Failed to execute bulk Gateway Profile DELETE operation",
			fmt.Sprintf("Error: %s", err),
		)
		return diagnostics
	}

	b.recentGatewayProfileOps = true
	b.recentGatewayProfileOpTime = time.Now()
	return diagnostics
}

func (b *BulkOperationManager) ExecuteBulkEthPortProfilePut(ctx context.Context) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	b.mutex.Lock()

	ethPortProfilePut := make(map[string]openapi.ConfigPutRequestEthPortProfileEthPortProfileName)
	for k, v := range b.ethPortProfilePut {
		ethPortProfilePut[k] = v
	}

	b.ethPortProfilePut = make(map[string]openapi.ConfigPutRequestEthPortProfileEthPortProfileName)

	b.mutex.Unlock()

	if len(ethPortProfilePut) == 0 {
		return diagnostics
	}

	ethPortProfileNames := make([]string, 0, len(ethPortProfilePut))
	for name := range ethPortProfilePut {
		ethPortProfileNames = append(ethPortProfileNames, name)
	}

	tflog.Debug(ctx, "Executing bulk EthPort Profile PUT operation", map[string]interface{}{
		"eth_port_profile_count": len(ethPortProfilePut),
		"eth_port_profile_names": ethPortProfileNames,
	})

	putRequest := openapi.NewEthportprofilesPutRequest()
	ethPortProfileMap := make(map[string]openapi.ConfigPutRequestEthPortProfileEthPortProfileName)

	for name, props := range ethPortProfilePut {
		ethPortProfileMap[name] = props
	}
	putRequest.SetEthPortProfile(ethPortProfileMap)
	retryConfig := DefaultRetryConfig()
	var err error

	for retry := 0; retry < retryConfig.MaxRetries; retry++ {
		// Create a separate context for the API call
		apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)

		req := b.client.EthPortProfilesAPI.EthportprofilesPut(apiCtx).EthportprofilesPutRequest(*putRequest)
		_, err = req.Execute()

		// Release the API call context
		cancel()

		if err == nil {
			tflog.Debug(ctx, "Bulk EthPortProfile PUT operation succeeded", map[string]interface{}{
				"attempt": retry + 1,
			})
			break
		}

		if IsRetriableError(err) {
			delayTime := CalculateBackoff(retry, retryConfig)
			tflog.Debug(ctx, "Bulk Eth Port Profile PUT operation failed with retriable error, retrying", map[string]interface{}{
				"attempt":     retry + 1,
				"error":       err.Error(),
				"delay_ms":    delayTime.Milliseconds(),
				"max_retries": retryConfig.MaxRetries,
			})

			time.Sleep(delayTime)
			continue
		}

		tflog.Error(ctx, "Bulk Eth Port Profile PUT operation failed with non-retriable error", map[string]interface{}{
			"error": err.Error(),
		})
		break
	}

	// Update operation statuses based on the result
	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	for opID, op := range b.pendingOperations {
		// Only process EthPortProfile PUT operations
		if op.ResourceType == "eth_port_profile" && op.OperationType == "PUT" {
			// Check if this operation's profile name is in our batch
			if _, exists := ethPortProfilePut[op.ResourceName]; exists {
				if err == nil {
					// Mark operation as successful
					op.Status = OperationSucceeded
					b.operationResults[opID] = true

					// Signal waiting resources
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				} else {
					// Mark operation as failed
					op.Status = OperationFailed
					op.Error = err
					b.operationErrors[opID] = err
					b.operationResults[opID] = false

					// Signal waiting resources with the error
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				}
			}
		}
	}

	if err != nil {
		diagnostics.AddError(
			"Failed to execute bulk Eth Port Profile PUT operation",
			fmt.Sprintf("Error: %s", err),
		)
		return diagnostics
	}

	b.recentEthPortProfileOps = true
	b.recentEthPortProfileOpTime = time.Now()
	return diagnostics
}

func (b *BulkOperationManager) ExecuteBulkEthPortSettingsPut(ctx context.Context) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	b.mutex.Lock()

	ethPortSettingsPut := make(map[string]openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName)
	for k, v := range b.ethPortSettingsPut {
		ethPortSettingsPut[k] = v
	}

	b.ethPortSettingsPut = make(map[string]openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName)

	b.mutex.Unlock()

	if len(ethPortSettingsPut) == 0 {
		return diagnostics
	}

	ethPortSettingsNames := make([]string, 0, len(ethPortSettingsPut))
	for name := range ethPortSettingsPut {
		ethPortSettingsNames = append(ethPortSettingsNames, name)
	}

	tflog.Debug(ctx, "Executing bulk EthPort Settings PUT operation", map[string]interface{}{
		"eth_port_settings_count": len(ethPortSettingsPut),
		"eth_port_settings_names": ethPortSettingsNames,
	})

	putRequest := openapi.NewEthportsettingsPutRequest()
	ethPortSettingsMap := make(map[string]openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName)

	for name, props := range ethPortSettingsPut {
		ethPortSettingsMap[name] = props
	}
	putRequest.SetEthPortSettings(ethPortSettingsMap)
	retryConfig := DefaultRetryConfig()
	var err error

	for retry := 0; retry < retryConfig.MaxRetries; retry++ {
		// Create a separate context for the API call
		apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)

		req := b.client.EthPortSettingsAPI.EthportsettingsPut(apiCtx).EthportsettingsPutRequest(*putRequest)
		_, err = req.Execute()

		// Release the API call context
		cancel()

		if err == nil {
			tflog.Debug(ctx, "Bulk EthPort Settings PUT operation succeeded", map[string]interface{}{
				"attempt": retry + 1,
			})
			break
		}

		if IsRetriableError(err) {
			delayTime := CalculateBackoff(retry, retryConfig)
			tflog.Debug(ctx, "Bulk EthPort Settings PUT operation failed with retriable error, retrying", map[string]interface{}{
				"attempt":     retry + 1,
				"error":       err.Error(),
				"delay_ms":    delayTime.Milliseconds(),
				"max_retries": retryConfig.MaxRetries,
			})

			time.Sleep(delayTime)
			continue
		}

		tflog.Error(ctx, "Bulk EthPort Settings PUT operation failed with non-retriable error", map[string]interface{}{
			"error": err.Error(),
		})
		break
	}

	// Update operation statuses based on the result
	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	for opID, op := range b.pendingOperations {
		// Only process EthPortSettings PUT operations
		if op.ResourceType == "eth_port_settings" && op.OperationType == "PUT" {
			// Check if this operation's settings name is in our batch
			if _, exists := ethPortSettingsPut[op.ResourceName]; exists {
				if err == nil {
					// Mark operation as successful
					op.Status = OperationSucceeded
					b.operationResults[opID] = true

					// Signal waiting resources
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				} else {
					// Mark operation as failed
					op.Status = OperationFailed
					op.Error = err
					b.operationErrors[opID] = err
					b.operationResults[opID] = false

					// Signal waiting resources with the error
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				}
			}
		}
	}

	if err != nil {
		diagnostics.AddError(
			"Failed to execute bulk EthPort Settings PUT operation",
			fmt.Sprintf("Error: %s", err),
		)
		return diagnostics
	}

	b.recentEthPortSettingsOps = true
	b.recentEthPortSettingsOpTime = time.Now()
	return diagnostics
}

func (b *BulkOperationManager) ExecuteBulkEthPortProfilePatch(ctx context.Context) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	b.mutex.Lock()

	ethPortProfilePatch := make(map[string]openapi.ConfigPutRequestEthPortProfileEthPortProfileName)
	for k, v := range b.ethPortProfilePatch {
		ethPortProfilePatch[k] = v
	}

	b.ethPortProfilePatch = make(map[string]openapi.ConfigPutRequestEthPortProfileEthPortProfileName)

	b.mutex.Unlock()

	if len(ethPortProfilePatch) == 0 {
		return diagnostics
	}

	ethPortProfileNames := make([]string, 0, len(ethPortProfilePatch))
	for name := range ethPortProfilePatch {
		ethPortProfileNames = append(ethPortProfileNames, name)
	}

	tflog.Debug(ctx, "Executing bulk EthPort Profile PATCH operation", map[string]interface{}{
		"eth_port_profile_count": len(ethPortProfilePatch),
		"eth_port_profile_names": ethPortProfileNames,
	})

	patchRequest := openapi.NewEthportprofilesPutRequest()
	ethPortProfileMap := make(map[string]openapi.ConfigPutRequestEthPortProfileEthPortProfileName)

	for name, props := range ethPortProfilePatch {
		ethPortProfileMap[name] = props
	}
	patchRequest.SetEthPortProfile(ethPortProfileMap)
	retryConfig := DefaultRetryConfig()
	var err error

	for retry := 0; retry < retryConfig.MaxRetries; retry++ {
		if retry > 0 {
			delay := CalculateBackoff(retry-1, retryConfig)

			tflog.Debug(ctx, "Retrying bulk Eth Port Profile PATCH operation after delay", map[string]interface{}{
				"retry": retry,
				"delay": delay,
			})

			time.Sleep(delay)
		}

		// Create a separate context for the API call
		apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)

		req := b.client.EthPortProfilesAPI.EthportprofilesPatch(apiCtx).EthportprofilesPutRequest(*patchRequest)
		_, err = req.Execute()

		// Release the API call context
		cancel()

		if err == nil {
			tflog.Debug(ctx, "Bulk Eth Port Profile PATCH operation successful", map[string]interface{}{
				"count": len(ethPortProfilePatch),
			})
			break
		}

		if !IsRetriableError(err) {
			tflog.Error(ctx, "Bulk Eth Port Profile PATCH operation failed with non-retriable error", map[string]interface{}{
				"error": err.Error(),
			})
			break
		}
	}

	// Update operation statuses based on the result
	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	for opID, op := range b.pendingOperations {
		// Only process EthPortProfile PATCH operations
		if op.ResourceType == "eth_port_profile" && op.OperationType == "PATCH" {
			// Check if this operation's profile name is in our batch
			if _, exists := ethPortProfilePatch[op.ResourceName]; exists {
				if err == nil {
					// Mark operation as successful
					op.Status = OperationSucceeded
					b.operationResults[opID] = true

					// Signal waiting resources
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				} else {
					// Mark operation as failed
					op.Status = OperationFailed
					op.Error = err
					b.operationErrors[opID] = err
					b.operationResults[opID] = false

					// Signal waiting resources with the error
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				}
			}
		}
	}

	if err != nil {
		diagnostics.AddError(
			"Failed to execute bulk Eth Port Profile PATCH operation",
			fmt.Sprintf("Error: %s", err),
		)
		return diagnostics
	}

	b.recentEthPortProfileOps = true
	b.recentEthPortProfileOpTime = time.Now()
	return diagnostics
}

func (b *BulkOperationManager) ExecuteBulkEthPortSettingsPatch(ctx context.Context) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	b.mutex.Lock()

	ethPortSettingsPatch := make(map[string]openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName)
	for k, v := range b.ethPortSettingsPatch {
		ethPortSettingsPatch[k] = v
	}

	b.ethPortSettingsPatch = make(map[string]openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName)

	b.mutex.Unlock()

	if len(ethPortSettingsPatch) == 0 {
		return diagnostics
	}

	ethPortSettingsNames := make([]string, 0, len(ethPortSettingsPatch))
	for name := range ethPortSettingsPatch {
		ethPortSettingsNames = append(ethPortSettingsNames, name)
	}

	tflog.Debug(ctx, "Executing bulk EthPort Settings PATCH operation", map[string]interface{}{
		"eth_port_settings_count": len(ethPortSettingsPatch),
		"eth_port_settings_names": ethPortSettingsNames,
	})

	patchRequest := openapi.NewEthportsettingsPutRequest()
	ethPortSettingsMap := make(map[string]openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName)

	for name, props := range ethPortSettingsPatch {
		ethPortSettingsMap[name] = props
	}
	patchRequest.SetEthPortSettings(ethPortSettingsMap)
	retryConfig := DefaultRetryConfig()
	var err error

	for retry := 0; retry < retryConfig.MaxRetries; retry++ {
		if retry > 0 {
			delay := CalculateBackoff(retry-1, retryConfig)

			tflog.Debug(ctx, "Retrying bulk EthPort Settings PATCH operation after delay", map[string]interface{}{
				"retry": retry,
				"delay": delay,
			})

			time.Sleep(delay)
		}

		// Create a separate context for the API call
		apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)

		req := b.client.EthPortSettingsAPI.EthportsettingsPatch(apiCtx).EthportsettingsPutRequest(*patchRequest)
		_, err = req.Execute()

		// Release the API call context
		cancel()

		if err == nil {
			tflog.Debug(ctx, "Bulk EthPort Settings PATCH operation successful", map[string]interface{}{
				"count": len(ethPortSettingsPatch),
			})
			break
		}

		if !IsRetriableError(err) {
			tflog.Error(ctx, "Bulk EthPort Settings PATCH operation failed with non-retriable error", map[string]interface{}{
				"error": err.Error(),
			})
			break
		}
	}

	// Update operation statuses based on the result
	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	for opID, op := range b.pendingOperations {
		// Only process EthPortSettings PATCH operations
		if op.ResourceType == "eth_port_settings" && op.OperationType == "PATCH" {
			// Check if this operation's settings name is in our batch
			if _, exists := ethPortSettingsPatch[op.ResourceName]; exists {
				if err == nil {
					// Mark operation as successful
					op.Status = OperationSucceeded
					b.operationResults[opID] = true

					// Signal waiting resources
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				} else {
					// Mark operation as failed
					op.Status = OperationFailed
					op.Error = err
					b.operationErrors[opID] = err
					b.operationResults[opID] = false

					// Signal waiting resources with the error
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				}
			}
		}
	}

	if err != nil {
		diagnostics.AddError(
			"Failed to execute bulk EthPort Settings PATCH operation",
			fmt.Sprintf("Error: %s", err),
		)
		return diagnostics
	}

	b.recentEthPortSettingsOps = true
	b.recentEthPortSettingsOpTime = time.Now()
	return diagnostics
}

func (b *BulkOperationManager) ExecuteBulkEthPortProfileDelete(ctx context.Context) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	b.mutex.Lock()

	ethPortProfileNames := make([]string, len(b.ethPortProfileDelete))
	copy(ethPortProfileNames, b.ethPortProfileDelete)

	ethPortProfileDeleteMap := make(map[string]bool)
	for _, name := range ethPortProfileNames {
		ethPortProfileDeleteMap[name] = true
	}

	b.ethPortProfileDelete = make([]string, 0)

	b.mutex.Unlock()

	if len(ethPortProfileNames) == 0 {
		return diagnostics
	}

	retryConfig := DefaultRetryConfig()
	var err error

	for retry := 0; retry < retryConfig.MaxRetries; retry++ {
		if retry > 0 {
			delay := CalculateBackoff(retry-1, retryConfig)

			tflog.Debug(ctx, "Retrying bulk Eth Port Profile DELETE operation after delay", map[string]interface{}{
				"retry":       retry,
				"delay":       delay,
				"max_retries": retryConfig.MaxRetries,
			})

			time.Sleep(delay)
		}

		tflog.Debug(ctx, "Executing bulk Eth Port Profile DELETE operation", map[string]interface{}{
			"eth_port_profile_count": len(ethPortProfileNames),
			"eth_port_profile_names": ethPortProfileNames,
		})

		// Create a separate context for the API call
		apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)

		req := b.client.EthPortProfilesAPI.EthportprofilesDelete(apiCtx).ProfileName(ethPortProfileNames)
		_, err = req.Execute()

		// Release the API call context
		cancel()

		if err == nil {
			tflog.Debug(ctx, "Bulk Eth Port Profile DELETE operation successful", map[string]interface{}{
				"count": len(ethPortProfileNames),
			})
			break
		}

		if !IsRetriableError(err) {
			tflog.Error(ctx, "Bulk Eth Port Profile DELETE operation failed with non-retriable error", map[string]interface{}{
				"error": err.Error(),
			})
			break
		}
	}

	// Update operation statuses based on the result
	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	for opID, op := range b.pendingOperations {
		// Only process EthPortProfile DELETE operations
		if op.ResourceType == "eth_port_profile" && op.OperationType == "DELETE" {
			// Check if this operation's profile name is in our batch
			if _, exists := ethPortProfileDeleteMap[op.ResourceName]; exists {
				if err == nil {
					// Mark operation as successful
					op.Status = OperationSucceeded
					b.operationResults[opID] = true

					// Signal waiting resources
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				} else {
					// Mark operation as failed
					op.Status = OperationFailed
					op.Error = err
					b.operationErrors[opID] = err
					b.operationResults[opID] = false

					// Signal waiting resources with the error
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				}
			}
		}
	}

	if err != nil {
		diagnostics.AddError(
			"Failed to execute bulk Eth Port Profile DELETE operation",
			fmt.Sprintf("Error: %s", err),
		)
		return diagnostics
	}

	b.recentEthPortProfileOps = true
	b.recentEthPortProfileOpTime = time.Now()
	return diagnostics
}

func (b *BulkOperationManager) ExecuteBulkEthPortSettingsDelete(ctx context.Context) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	b.mutex.Lock()

	ethPortSettingsNames := make([]string, len(b.ethPortSettingsDelete))
	copy(ethPortSettingsNames, b.ethPortSettingsDelete)

	ethPortSettingsDeleteMap := make(map[string]bool)
	for _, name := range ethPortSettingsNames {
		ethPortSettingsDeleteMap[name] = true
	}

	b.ethPortSettingsDelete = make([]string, 0)

	b.mutex.Unlock()

	if len(ethPortSettingsNames) == 0 {
		return diagnostics
	}

	retryConfig := DefaultRetryConfig()
	var err error

	for retry := 0; retry < retryConfig.MaxRetries; retry++ {
		if retry > 0 {
			delay := CalculateBackoff(retry-1, retryConfig)

			tflog.Debug(ctx, "Retrying bulk EthPort Settings DELETE operation after delay", map[string]interface{}{
				"retry":       retry,
				"delay":       delay,
				"max_retries": retryConfig.MaxRetries,
			})

			time.Sleep(delay)
		}

		tflog.Debug(ctx, "Executing bulk EthPort Settings DELETE operation", map[string]interface{}{
			"eth_port_settings_count": len(ethPortSettingsNames),
			"eth_port_settings_names": ethPortSettingsNames,
		})

		// Create a separate context for the API call
		apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)

		req := b.client.EthPortSettingsAPI.EthportsettingsDelete(apiCtx).PortName(ethPortSettingsNames)
		_, err = req.Execute()

		// Release the API call context
		cancel()

		if err == nil {
			tflog.Debug(ctx, "Bulk EthPort Settings DELETE operation successful", map[string]interface{}{
				"count": len(ethPortSettingsNames),
			})
			break
		}

		if !IsRetriableError(err) {
			tflog.Error(ctx, "Bulk EthPort Settings DELETE operation failed with non-retriable error", map[string]interface{}{
				"error": err.Error(),
			})
			break
		}
	}

	// Update operation statuses based on the result
	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	for opID, op := range b.pendingOperations {
		// Only process EthPortSettings DELETE operations
		if op.ResourceType == "eth_port_settings" && op.OperationType == "DELETE" {
			// Check if this operation's settings name is in our batch
			if _, exists := ethPortSettingsDeleteMap[op.ResourceName]; exists {
				if err == nil {
					// Mark operation as successful
					op.Status = OperationSucceeded
					b.operationResults[opID] = true

					// Signal waiting resources
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				} else {
					// Mark operation as failed
					op.Status = OperationFailed
					op.Error = err
					b.operationErrors[opID] = err
					b.operationResults[opID] = false

					// Signal waiting resources with the error
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				}
			}
		}
	}

	if err != nil {
		diagnostics.AddError(
			"Failed to execute bulk EthPort Settings DELETE operation",
			fmt.Sprintf("Error: %s", err),
		)
		return diagnostics
	}

	b.recentEthPortSettingsOps = true
	b.recentEthPortSettingsOpTime = time.Now()
	return diagnostics
}

func (b *BulkOperationManager) ExecuteBulkBundlePatch(ctx context.Context) diag.Diagnostics {
	var diagnostics diag.Diagnostics

	b.mutex.Lock()

	bundlePatch := make(map[string]openapi.BundlesPatchRequestEndpointBundleValue)
	for k, v := range b.bundlePatch {
		bundlePatch[k] = v
	}

	b.bundlePatch = make(map[string]openapi.BundlesPatchRequestEndpointBundleValue)

	b.mutex.Unlock()

	if len(bundlePatch) == 0 {
		return diagnostics
	}

	bundleNames := make([]string, 0, len(bundlePatch))
	for name := range bundlePatch {
		bundleNames = append(bundleNames, name)
	}

	tflog.Debug(ctx, "Executing bulk Bundle PATCH operation", map[string]interface{}{
		"bundle_count": len(bundlePatch),
		"bundle_names": bundleNames,
	})

	patchRequest := openapi.NewBundlesPatchRequest()
	patchRequest.SetEndpointBundle(bundlePatch)

	retryConfig := DefaultRetryConfig()
	var err error

	for retry := 0; retry < retryConfig.MaxRetries; retry++ {
		if retry > 0 {
			delay := CalculateBackoff(retry-1, retryConfig)

			tflog.Debug(ctx, "Retrying bulk Bundle PATCH operation after delay", map[string]interface{}{
				"retry": retry,
				"delay": delay,
			})

			time.Sleep(delay)
		}

		// Create a separate context for the API call
		apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)

		req := b.client.BundlesAPI.BundlesPatch(apiCtx).BundlesPatchRequest(*patchRequest)
		_, err = req.Execute()

		// Release the API call context
		cancel()

		if err == nil {
			tflog.Debug(ctx, "Bulk Bundle PATCH operation successful", map[string]interface{}{
				"count": len(bundlePatch),
			})
			break
		}

		if !IsRetriableError(err) {
			tflog.Error(ctx, "Bulk Bundle PATCH operation failed with non-retriable error", map[string]interface{}{
				"error": err.Error(),
			})
			break
		}
	}

	// Update operation statuses based on the result
	b.operationMutex.Lock()
	defer b.operationMutex.Unlock()

	for opID, op := range b.pendingOperations {
		// Only process Bundle PATCH operations
		if op.ResourceType == "bundle" && op.OperationType == "PATCH" {
			// Check if this operation's bundle name is in our batch
			if _, exists := bundlePatch[op.ResourceName]; exists {
				if err == nil {
					// Mark operation as successful
					op.Status = OperationSucceeded
					b.operationResults[opID] = true

					// Signal waiting resources
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				} else {
					// Mark operation as failed
					op.Status = OperationFailed
					op.Error = err
					b.operationErrors[opID] = err
					b.operationResults[opID] = false

					// Signal waiting resources with the error
					if waitCh, ok := b.operationWaitChannels[opID]; ok {
						close(waitCh)
					}
				}
			}
		}
	}

	if err != nil {
		diagnostics.AddError(
			"Failed to execute bulk Bundle PATCH operation",
			fmt.Sprintf("Error: %s", err),
		)
		return diagnostics
	}

	b.recentBundleOps = true
	b.recentBundleOpTime = time.Now()
	return diagnostics
}
