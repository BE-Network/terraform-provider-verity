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
	mode              string
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
	bundlePut    map[string]openapi.BundlesPutRequestEndpointBundleValue
	bundlePatch  map[string]openapi.BundlesPatchRequestEndpointBundleValue
	bundleDelete []string

	// ACL operations (unified for IPv4 and IPv6)
	aclPut       map[string]openapi.ConfigPutRequestIpv4FilterIpv4FilterName
	aclPatch     map[string]openapi.ConfigPutRequestIpv4FilterIpv4FilterName
	aclDelete    []string
	aclIpVersion map[string]string // Track which IP version each ACL operation uses

	// Authenticated Eth-Port operations
	authenticatedEthPortPut    map[string]openapi.ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName
	authenticatedEthPortPatch  map[string]openapi.ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName
	authenticatedEthPortDelete []string

	// Badge operations
	badgePut    map[string]openapi.ConfigPutRequestBadgeBadgeName
	badgePatch  map[string]openapi.ConfigPutRequestBadgeBadgeName
	badgeDelete []string

	// Device Port Settings operations
	deviceVoiceSettingsPut    map[string]openapi.ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName
	deviceVoiceSettingsPatch  map[string]openapi.ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName
	deviceVoiceSettingsDelete []string

	// Packet Broker operations
	packetBrokerPut    map[string]openapi.ConfigPutRequestPbEgressProfilePbEgressProfileName
	packetBrokerPatch  map[string]openapi.ConfigPutRequestPbEgressProfilePbEgressProfileName
	packetBrokerDelete []string

	// Packet Queues operations
	packetQueuePut    map[string]openapi.ConfigPutRequestPacketQueuePacketQueueName
	packetQueuePatch  map[string]openapi.ConfigPutRequestPacketQueuePacketQueueName
	packetQueueDelete []string

	// ServicePort Profile operations
	servicePortProfilePut    map[string]openapi.ConfigPutRequestServicePortProfileServicePortProfileName
	servicePortProfilePatch  map[string]openapi.ConfigPutRequestServicePortProfileServicePortProfileName
	servicePortProfileDelete []string

	// Switchpoint operations
	switchpointPut    map[string]openapi.ConfigPutRequestSwitchpointSwitchpointName
	switchpointPatch  map[string]openapi.ConfigPutRequestSwitchpointSwitchpointName
	switchpointDelete []string

	// Voice Port Profile operations
	voicePortProfilePut    map[string]openapi.ConfigPutRequestVoicePortProfilesVoicePortProfilesName
	voicePortProfilePatch  map[string]openapi.ConfigPutRequestVoicePortProfilesVoicePortProfilesName
	voicePortProfileDelete []string

	// Device Controller operations
	deviceControllerPut    map[string]openapi.ConfigPutRequestDeviceControllerDeviceControllerName
	deviceControllerPatch  map[string]openapi.ConfigPutRequestDeviceControllerDeviceControllerName
	deviceControllerDelete []string

	// Track recent operations to avoid race conditions
	recentGatewayOps                 bool
	recentGatewayOpTime              time.Time
	recentLagOps                     bool
	recentLagOpTime                  time.Time
	recentServiceOps                 bool
	recentServiceOpTime              time.Time
	recentTenantOps                  bool
	recentTenantOpTime               time.Time
	recentGatewayProfileOps          bool
	recentGatewayProfileOpTime       time.Time
	recentEthPortProfileOps          bool
	recentEthPortProfileOpTime       time.Time
	recentEthPortSettingsOps         bool
	recentEthPortSettingsOpTime      time.Time
	recentBundleOps                  bool
	recentBundleOpTime               time.Time
	recentAclOps                     bool
	recentAclOpTime                  time.Time
	recentAuthenticatedEthPortOps    bool
	recentAuthenticatedEthPortOpTime time.Time
	recentBadgeOps                   bool
	recentBadgeOpTime                time.Time
	recentVoicePortProfileOps        bool
	recentVoicePortProfileOpTime     time.Time
	recentSwitchpointOps             bool
	recentSwitchpointOpTime          time.Time
	recentServicePortProfileOps      bool
	recentServicePortProfileOpTime   time.Time
	recentPacketBrokerOps            bool
	recentPacketBrokerOpTime         time.Time
	recentPacketQueueOps             bool
	recentPacketQueueOpTime          time.Time
	recentDeviceVoiceSettingsOps     bool
	recentDeviceVoiceSettingsOpTime  time.Time
	recentDeviceControllerOps        bool
	recentDeviceControllerOpTime     time.Time

	// For tracking operations
	pendingOperations     map[string]*Operation
	operationResults      map[string]bool // true = success, false = failure
	operationErrors       map[string]error
	operationWaitChannels map[string]chan struct{}
	operationMutex        sync.Mutex
	closedChannels        map[string]bool

	// Store API responses
	gatewayResponses                   map[string]map[string]interface{}
	gatewayResponsesMutex              sync.RWMutex
	lagResponses                       map[string]map[string]interface{}
	lagResponsesMutex                  sync.RWMutex
	serviceResponses                   map[string]map[string]interface{}
	serviceResponsesMutex              sync.RWMutex
	tenantResponses                    map[string]map[string]interface{}
	tenantResponsesMutex               sync.RWMutex
	gatewayProfileResponses            map[string]map[string]interface{}
	gatewayProfileResponsesMutex       sync.RWMutex
	ethPortProfileResponses            map[string]map[string]interface{}
	ethPortProfileResponsesMutex       sync.RWMutex
	ethPortSettingsResponses           map[string]map[string]interface{}
	ethPortSettingsResponsesMutex      sync.RWMutex
	bundleResponses                    map[string]map[string]interface{}
	bundleResponsesMutex               sync.RWMutex
	aclResponses                       map[string]map[string]interface{}
	aclResponsesMutex                  sync.RWMutex
	authenticatedEthPortResponses      map[string]map[string]interface{}
	authenticatedEthPortResponsesMutex sync.RWMutex
	badgeResponses                     map[string]map[string]interface{}
	badgeResponsesMutex                sync.RWMutex
	voicePortProfileResponses          map[string]map[string]interface{}
	voicePortProfileResponsesMutex     sync.RWMutex
	switchpointResponses               map[string]map[string]interface{}
	switchpointResponsesMutex          sync.RWMutex
	servicePortProfileResponses        map[string]map[string]interface{}
	servicePortProfileResponsesMutex   sync.RWMutex
	packetBrokerResponses              map[string]map[string]interface{}
	packetBrokerResponsesMutex         sync.RWMutex
	packetQueueResponses               map[string]map[string]interface{}
	packetQueueResponsesMutex          sync.RWMutex
	deviceVoiceSettingsResponses       map[string]map[string]interface{}
	deviceVoiceSettingsResponsesMutex  sync.RWMutex
	deviceControllerResponses          map[string]map[string]interface{}
	deviceControllerResponsesMutex     sync.RWMutex
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

func NewBulkOperationManager(client *openapi.APIClient, contextProvider ContextProviderFunc, clearCacheFunc ClearCacheFunc, mode string) *BulkOperationManager {
	return &BulkOperationManager{
		client:                     client,
		contextProvider:            contextProvider,
		clearCacheFunc:             clearCacheFunc,
		mode:                       mode,
		lastOperationTime:          time.Now(),
		gatewayPut:                 make(map[string]openapi.ConfigPutRequestGatewayGatewayName),
		gatewayPatch:               make(map[string]openapi.ConfigPutRequestGatewayGatewayName),
		gatewayDelete:              make([]string, 0),
		lagPut:                     make(map[string]openapi.ConfigPutRequestLagLagName),
		lagPatch:                   make(map[string]openapi.ConfigPutRequestLagLagName),
		lagDelete:                  make([]string, 0),
		tenantPut:                  make(map[string]openapi.ConfigPutRequestTenantTenantName),
		tenantPatch:                make(map[string]openapi.ConfigPutRequestTenantTenantName),
		tenantDelete:               make([]string, 0),
		servicePut:                 make(map[string]openapi.ConfigPutRequestServiceServiceName),
		servicePatch:               make(map[string]openapi.ConfigPutRequestServiceServiceName),
		serviceDelete:              make([]string, 0),
		gatewayProfilePut:          make(map[string]openapi.ConfigPutRequestGatewayProfileGatewayProfileName),
		gatewayProfilePatch:        make(map[string]openapi.ConfigPutRequestGatewayProfileGatewayProfileName),
		gatewayProfileDelete:       make([]string, 0),
		ethPortProfilePut:          make(map[string]openapi.ConfigPutRequestEthPortProfileEthPortProfileName),
		ethPortProfilePatch:        make(map[string]openapi.ConfigPutRequestEthPortProfileEthPortProfileName),
		ethPortProfileDelete:       make([]string, 0),
		ethPortSettingsPut:         make(map[string]openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName),
		ethPortSettingsPatch:       make(map[string]openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName),
		ethPortSettingsDelete:      make([]string, 0),
		bundlePut:                  make(map[string]openapi.BundlesPutRequestEndpointBundleValue),
		bundlePatch:                make(map[string]openapi.BundlesPatchRequestEndpointBundleValue),
		bundleDelete:               make([]string, 0),
		aclPut:                     make(map[string]openapi.ConfigPutRequestIpv4FilterIpv4FilterName),
		aclPatch:                   make(map[string]openapi.ConfigPutRequestIpv4FilterIpv4FilterName),
		aclDelete:                  make([]string, 0),
		aclIpVersion:               make(map[string]string),
		authenticatedEthPortPut:    make(map[string]openapi.ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName),
		authenticatedEthPortPatch:  make(map[string]openapi.ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName),
		authenticatedEthPortDelete: make([]string, 0),
		badgePut:                   make(map[string]openapi.ConfigPutRequestBadgeBadgeName),
		badgePatch:                 make(map[string]openapi.ConfigPutRequestBadgeBadgeName),
		badgeDelete:                make([]string, 0),
		voicePortProfilePut:        make(map[string]openapi.ConfigPutRequestVoicePortProfilesVoicePortProfilesName),
		voicePortProfilePatch:      make(map[string]openapi.ConfigPutRequestVoicePortProfilesVoicePortProfilesName),
		voicePortProfileDelete:     make([]string, 0),
		switchpointPut:             make(map[string]openapi.ConfigPutRequestSwitchpointSwitchpointName),
		switchpointPatch:           make(map[string]openapi.ConfigPutRequestSwitchpointSwitchpointName),
		switchpointDelete:          make([]string, 0),
		servicePortProfilePut:      make(map[string]openapi.ConfigPutRequestServicePortProfileServicePortProfileName),
		servicePortProfilePatch:    make(map[string]openapi.ConfigPutRequestServicePortProfileServicePortProfileName),
		servicePortProfileDelete:   make([]string, 0),
		packetBrokerPut:            make(map[string]openapi.ConfigPutRequestPbEgressProfilePbEgressProfileName),
		packetBrokerPatch:          make(map[string]openapi.ConfigPutRequestPbEgressProfilePbEgressProfileName),
		packetBrokerDelete:         make([]string, 0),
		packetQueuePut:             make(map[string]openapi.ConfigPutRequestPacketQueuePacketQueueName),
		packetQueuePatch:           make(map[string]openapi.ConfigPutRequestPacketQueuePacketQueueName),
		packetQueueDelete:          make([]string, 0),
		deviceVoiceSettingsPut:     make(map[string]openapi.ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName),
		deviceVoiceSettingsPatch:   make(map[string]openapi.ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName),
		deviceVoiceSettingsDelete:  make([]string, 0),
		deviceControllerPut:        make(map[string]openapi.ConfigPutRequestDeviceControllerDeviceControllerName),
		deviceControllerPatch:      make(map[string]openapi.ConfigPutRequestDeviceControllerDeviceControllerName),
		deviceControllerDelete:     make([]string, 0),
		pendingOperations:          make(map[string]*Operation),
		operationResults:           make(map[string]bool),
		operationErrors:            make(map[string]error),
		operationWaitChannels:      make(map[string]chan struct{}),
		closedChannels:             make(map[string]bool),

		// Initialize with no recent operations
		recentGatewayOps:              false,
		recentLagOps:                  false,
		recentServiceOps:              false,
		recentTenantOps:               false,
		recentGatewayProfileOps:       false,
		recentEthPortProfileOps:       false,
		recentEthPortSettingsOps:      false,
		recentBundleOps:               false,
		recentAclOps:                  false,
		recentAuthenticatedEthPortOps: false,
		recentBadgeOps:                false,
		recentVoicePortProfileOps:     false,
		recentSwitchpointOps:          false,
		recentServicePortProfileOps:   false,
		recentPacketBrokerOps:         false,
		recentPacketQueueOps:          false,
		recentDeviceVoiceSettingsOps:  false,
		recentDeviceControllerOps:     false,

		// Initialize response caches
		gatewayResponses:                   make(map[string]map[string]interface{}),
		gatewayResponsesMutex:              sync.RWMutex{},
		lagResponses:                       make(map[string]map[string]interface{}),
		lagResponsesMutex:                  sync.RWMutex{},
		serviceResponses:                   make(map[string]map[string]interface{}),
		serviceResponsesMutex:              sync.RWMutex{},
		tenantResponses:                    make(map[string]map[string]interface{}),
		tenantResponsesMutex:               sync.RWMutex{},
		gatewayProfileResponses:            make(map[string]map[string]interface{}),
		gatewayProfileResponsesMutex:       sync.RWMutex{},
		ethPortProfileResponses:            make(map[string]map[string]interface{}),
		ethPortProfileResponsesMutex:       sync.RWMutex{},
		ethPortSettingsResponses:           make(map[string]map[string]interface{}),
		ethPortSettingsResponsesMutex:      sync.RWMutex{},
		bundleResponses:                    make(map[string]map[string]interface{}),
		bundleResponsesMutex:               sync.RWMutex{},
		aclResponses:                       make(map[string]map[string]interface{}),
		aclResponsesMutex:                  sync.RWMutex{},
		authenticatedEthPortResponses:      make(map[string]map[string]interface{}),
		authenticatedEthPortResponsesMutex: sync.RWMutex{},
		badgeResponses:                     make(map[string]map[string]interface{}),
		badgeResponsesMutex:                sync.RWMutex{},
		voicePortProfileResponses:          make(map[string]map[string]interface{}),
		voicePortProfileResponsesMutex:     sync.RWMutex{},
		switchpointResponses:               make(map[string]map[string]interface{}),
		switchpointResponsesMutex:          sync.RWMutex{},
		servicePortProfileResponses:        make(map[string]map[string]interface{}),
		servicePortProfileResponsesMutex:   sync.RWMutex{},
		packetBrokerResponses:              make(map[string]map[string]interface{}),
		packetBrokerResponsesMutex:         sync.RWMutex{},
		packetQueueResponses:               make(map[string]map[string]interface{}),
		packetQueueResponsesMutex:          sync.RWMutex{},
		deviceVoiceSettingsResponses:       make(map[string]map[string]interface{}),
		deviceVoiceSettingsResponsesMutex:  sync.RWMutex{},
		deviceControllerResponses:          make(map[string]map[string]interface{}),
		deviceControllerResponsesMutex:     sync.RWMutex{},
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

	case "acl_v4", "acl_v6":
		b.aclResponsesMutex.RLock()
		defer b.aclResponsesMutex.RUnlock()
		response, exists := b.aclResponses[resourceName]
		return response, exists

	case "authenticated_eth_port":
		b.authenticatedEthPortResponsesMutex.RLock()
		defer b.authenticatedEthPortResponsesMutex.RUnlock()
		response, exists := b.authenticatedEthPortResponses[resourceName]
		return response, exists

	case "badge":
		b.badgeResponsesMutex.RLock()
		defer b.badgeResponsesMutex.RUnlock()
		response, exists := b.badgeResponses[resourceName]
		return response, exists

	case "switchpoint":
		b.switchpointResponsesMutex.RLock()
		defer b.switchpointResponsesMutex.RUnlock()
		response, exists := b.switchpointResponses[resourceName]
		return response, exists

	case "service_port_profile":
		b.servicePortProfileResponsesMutex.RLock()
		defer b.servicePortProfileResponsesMutex.RUnlock()
		response, exists := b.servicePortProfileResponses[resourceName]
		return response, exists

	case "packet_broker":
		b.packetBrokerResponsesMutex.RLock()
		defer b.packetBrokerResponsesMutex.RUnlock()
		response, exists := b.packetBrokerResponses[resourceName]
		return response, exists

	case "packet_queue":
		b.packetQueueResponsesMutex.RLock()
		defer b.packetQueueResponsesMutex.RUnlock()
		response, exists := b.packetQueueResponses[resourceName]
		return response, exists

	case "device_voice_settings":
		b.deviceVoiceSettingsResponsesMutex.RLock()
		defer b.deviceVoiceSettingsResponsesMutex.RUnlock()
		response, exists := b.deviceVoiceSettingsResponses[resourceName]
		return response, exists

	case "voice_port_profile":
		b.voicePortProfileResponsesMutex.RLock()
		defer b.voicePortProfileResponsesMutex.RUnlock()
		response, exists := b.voicePortProfileResponses[resourceName]
		return response, exists

	case "device_controller":
		b.deviceControllerResponsesMutex.RLock()
		defer b.deviceControllerResponsesMutex.RUnlock()
		response, exists := b.deviceControllerResponses[resourceName]
		return response, exists

	default:
		return nil, false
	}
}

func GetBulkOperationManager(client *openapi.APIClient, clearCacheFunc ClearCacheFunc, providerContext interface{}, mode string) *BulkOperationManager {
	contextProvider := func() interface{} {
		return providerContext
	}

	return NewBulkOperationManager(client, contextProvider, clearCacheFunc, mode)
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

	var opsDiags diag.Diagnostics
	var operationsPerformed bool

	switch b.mode {
	case "datacenter":
		tflog.Debug(ctx, "Executing pending operations in 'datacenter' mode")
		opsDiags, operationsPerformed = b.ExecuteDatacenterOperations(ctx)
	case "campus":
		tflog.Debug(ctx, "Executing pending operations in 'campus' mode")
		opsDiags, operationsPerformed = b.ExecuteCampusOperations(ctx)
	default:
		tflog.Warn(ctx, fmt.Sprintf("Unknown mode '%s', defaulting to 'datacenter' mode", b.mode))
		opsDiags, operationsPerformed = b.ExecuteDatacenterOperations(ctx)
	}

	diagnostics.Append(opsDiags...)

	if operationsPerformed {
		waitDuration := 800 * time.Millisecond
		tflog.Debug(ctx, fmt.Sprintf("Waiting %v for all operations to propagate before final cache refresh", waitDuration))
		time.Sleep(waitDuration)

		tflog.Debug(ctx, "Final cache clear after all operations to ensure verification with fresh data")
		if b.clearCacheFunc != nil && b.contextProvider != nil {
			b.clearCacheFunc(ctx, b.contextProvider(), "tenants")
			b.clearCacheFunc(ctx, b.contextProvider(), "gateways")
			b.clearCacheFunc(ctx, b.contextProvider(), "gatewayprofiles")
			b.clearCacheFunc(ctx, b.contextProvider(), "services")
			b.clearCacheFunc(ctx, b.contextProvider(), "packetqueues")
			b.clearCacheFunc(ctx, b.contextProvider(), "ethportprofiles")
			b.clearCacheFunc(ctx, b.contextProvider(), "ethportsettings")
			b.clearCacheFunc(ctx, b.contextProvider(), "lags")
			b.clearCacheFunc(ctx, b.contextProvider(), "bundles")
			b.clearCacheFunc(ctx, b.contextProvider(), "acls_ipv4")
			b.clearCacheFunc(ctx, b.contextProvider(), "acls_ipv6")
			b.clearCacheFunc(ctx, b.contextProvider(), "packetbroker")
			b.clearCacheFunc(ctx, b.contextProvider(), "badges")
			b.clearCacheFunc(ctx, b.contextProvider(), "switchpoints")
			b.clearCacheFunc(ctx, b.contextProvider(), "devicecontrollers")
			b.clearCacheFunc(ctx, b.contextProvider(), "authenticatedethports")
			b.clearCacheFunc(ctx, b.contextProvider(), "devicevoicesettings")
			b.clearCacheFunc(ctx, b.contextProvider(), "serviceportprofiles")
			b.clearCacheFunc(ctx, b.contextProvider(), "voiceportprofiles")
		}
	}

	return diagnostics
}

func (b *BulkOperationManager) ExecuteDatacenterOperations(ctx context.Context) (diag.Diagnostics, bool) {
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
				b.FailAllPendingOperations(ctx, err)
				return false
			}
			operationsPerformed = true
		}
		return true
	}

	// PUT operations - DC Order
	if !execute("PUT", len(b.tenantPut), b.ExecuteBulkTenantPut, "Tenant") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.gatewayPut), b.ExecuteBulkGatewayPut, "Gateway") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.gatewayProfilePut), b.ExecuteBulkGatewayProfilePut, "Gateway Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.servicePut), b.ExecuteBulkServicePut, "Service") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.packetQueuePut), b.ExecuteBulkPacketQueuePut, "Packet Queue") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.ethPortProfilePut), b.ExecuteBulkEthPortProfilePut, "Eth Port Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.ethPortSettingsPut), b.ExecuteBulkEthPortSettingsPut, "Eth Port Settings") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.lagPut), b.ExecuteBulkLagPut, "LAG") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.bundlePut), b.ExecuteBulkBundlePut, "Bundle") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.aclPut), b.ExecuteBulkAclPut, "ACL") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.packetBrokerPut), b.ExecuteBulkPacketBrokerPut, "Packet Broker") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.badgePut), b.ExecuteBulkBadgePut, "Badge") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.switchpointPut), b.ExecuteBulkSwitchpointPut, "Switchpoint") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.deviceControllerPut), b.ExecuteBulkDeviceControllerPut, "Device Controller") {
		return diagnostics, operationsPerformed
	}

	// PATCH operations - DC Order
	if !execute("PATCH", len(b.tenantPatch), b.ExecuteBulkTenantPatch, "Tenant") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.gatewayPatch), b.ExecuteBulkGatewayPatch, "Gateway") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.gatewayProfilePatch), b.ExecuteBulkGatewayProfilePatch, "Gateway Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.servicePatch), b.ExecuteBulkServicePatch, "Service") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.packetQueuePatch), b.ExecuteBulkPacketQueuePatch, "Packet Queue") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.ethPortProfilePatch), b.ExecuteBulkEthPortProfilePatch, "Eth Port Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.ethPortSettingsPatch), b.ExecuteBulkEthPortSettingsPatch, "Eth Port Settings") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.lagPatch), b.ExecuteBulkLagPatch, "LAG") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.bundlePatch), b.ExecuteBulkBundlePatch, "Bundle") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.aclPatch), b.ExecuteBulkAclPatch, "ACL") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.packetBrokerPatch), b.ExecuteBulkPacketBrokerPatch, "Packet Broker") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.badgePatch), b.ExecuteBulkBadgePatch, "Badge") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.switchpointPatch), b.ExecuteBulkSwitchpointPatch, "Switchpoint") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.deviceControllerPatch), b.ExecuteBulkDeviceControllerPatch, "Device Controller") {
		return diagnostics, operationsPerformed
	}

	// DELETE operations - Reverse DC Order
	if !execute("DELETE", len(b.deviceControllerDelete), b.ExecuteBulkDeviceControllerDelete, "Device Controller") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.switchpointDelete), b.ExecuteBulkSwitchpointDelete, "Switchpoint") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.badgeDelete), b.ExecuteBulkBadgeDelete, "Badge") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.packetBrokerDelete), b.ExecuteBulkPacketBrokerDelete, "Packet Broker") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.aclDelete), b.ExecuteBulkAclDelete, "ACL") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.bundleDelete), b.ExecuteBulkBundleDelete, "Bundle") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.lagDelete), b.ExecuteBulkLagDelete, "LAG") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.ethPortSettingsDelete), b.ExecuteBulkEthPortSettingsDelete, "Eth Port Settings") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.ethPortProfileDelete), b.ExecuteBulkEthPortProfileDelete, "Eth Port Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.packetQueueDelete), b.ExecuteBulkPacketQueueDelete, "Packet Queue") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.serviceDelete), b.ExecuteBulkServiceDelete, "Service") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.gatewayProfileDelete), b.ExecuteBulkGatewayProfileDelete, "Gateway Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.gatewayDelete), b.ExecuteBulkGatewayDelete, "Gateway") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.tenantDelete), b.ExecuteBulkTenantDelete, "Tenant") {
		return diagnostics, operationsPerformed
	}

	return diagnostics, operationsPerformed
}

func (b *BulkOperationManager) ExecuteCampusOperations(ctx context.Context) (diag.Diagnostics, bool) {
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
				b.FailAllPendingOperations(ctx, err)
				return false
			}
			operationsPerformed = true
		}
		return true
	}

	// PUT operations - Campus Order
	if !execute("PUT", len(b.servicePut), b.ExecuteBulkServicePut, "Service") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.ethPortProfilePut), b.ExecuteBulkEthPortProfilePut, "Eth Port Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.authenticatedEthPortPut), b.ExecuteBulkAuthenticatedEthPortPut, "Authenticated Eth-Port") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.deviceVoiceSettingsPut), b.ExecuteBulkDeviceVoiceSettingsPut, "Device Voice Settings") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.packetQueuePut), b.ExecuteBulkPacketQueuePut, "Packet Queue") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.servicePortProfilePut), b.ExecuteBulkServicePortProfilePut, "Service Port Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.voicePortProfilePut), b.ExecuteBulkVoicePortProfilePut, "Voice-Port Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.ethPortSettingsPut), b.ExecuteBulkEthPortSettingsPut, "Eth Port Settings") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.lagPut), b.ExecuteBulkLagPut, "LAG") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.bundlePut), b.ExecuteBulkBundlePut, "Bundle") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.badgePut), b.ExecuteBulkBadgePut, "Badge") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.switchpointPut), b.ExecuteBulkSwitchpointPut, "Switchpoint") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.deviceControllerPut), b.ExecuteBulkDeviceControllerPut, "Device Controller") {
		return diagnostics, operationsPerformed
	}

	// PATCH operations - Campus Order
	if !execute("PATCH", len(b.servicePatch), b.ExecuteBulkServicePatch, "Service") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.ethPortProfilePatch), b.ExecuteBulkEthPortProfilePatch, "Eth Port Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.authenticatedEthPortPatch), b.ExecuteBulkAuthenticatedEthPortPatch, "Authenticated Eth-Port") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.deviceVoiceSettingsPatch), b.ExecuteBulkDeviceVoiceSettingsPatch, "Device Voice Settings") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.packetQueuePatch), b.ExecuteBulkPacketQueuePatch, "Packet Queue") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.servicePortProfilePatch), b.ExecuteBulkServicePortProfilePatch, "Service Port Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.voicePortProfilePatch), b.ExecuteBulkVoicePortProfilePatch, "Voice-Port Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.ethPortSettingsPatch), b.ExecuteBulkEthPortSettingsPatch, "Eth Port Settings") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.lagPatch), b.ExecuteBulkLagPatch, "LAG") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.bundlePatch), b.ExecuteBulkBundlePatch, "Bundle") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.badgePatch), b.ExecuteBulkBadgePatch, "Badge") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.switchpointPatch), b.ExecuteBulkSwitchpointPatch, "Switchpoint") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.deviceControllerPatch), b.ExecuteBulkDeviceControllerPatch, "Device Controller") {
		return diagnostics, operationsPerformed
	}

	// DELETE operations - Reverse Campus Order
	if !execute("DELETE", len(b.deviceControllerDelete), b.ExecuteBulkDeviceControllerDelete, "Device Controller") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.switchpointDelete), b.ExecuteBulkSwitchpointDelete, "Switchpoint") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.badgeDelete), b.ExecuteBulkBadgeDelete, "Badge") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.bundleDelete), b.ExecuteBulkBundleDelete, "Bundle") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.lagDelete), b.ExecuteBulkLagDelete, "LAG") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.ethPortSettingsDelete), b.ExecuteBulkEthPortSettingsDelete, "Eth Port Settings") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.voicePortProfileDelete), b.ExecuteBulkVoicePortProfileDelete, "Voice-Port Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.servicePortProfileDelete), b.ExecuteBulkServicePortProfileDelete, "Service Port Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.packetQueueDelete), b.ExecuteBulkPacketQueueDelete, "Packet Queue") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.deviceVoiceSettingsDelete), b.ExecuteBulkDeviceVoiceSettingsDelete, "Device Voice Settings") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.authenticatedEthPortDelete), b.ExecuteBulkAuthenticatedEthPortDelete, "Authenticated Eth-Port") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.ethPortProfileDelete), b.ExecuteBulkEthPortProfileDelete, "Eth Port Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.serviceDelete), b.ExecuteBulkServiceDelete, "Service") {
		return diagnostics, operationsPerformed
	}

	return diagnostics, operationsPerformed
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
		len(b.bundlePut) == 0 && len(b.bundlePatch) == 0 && len(b.bundleDelete) == 0 && len(b.authenticatedEthPortPut) == 0 && len(b.authenticatedEthPortPatch) == 0 &&
		len(b.authenticatedEthPortDelete) == 0 && len(b.aclPut) == 0 && len(b.aclPatch) == 0 && len(b.aclDelete) == 0 &&
		len(b.badgePut) == 0 && len(b.badgePatch) == 0 && len(b.badgeDelete) == 0 &&
		len(b.voicePortProfilePut) == 0 && len(b.voicePortProfilePatch) == 0 && len(b.voicePortProfileDelete) == 0 &&
		len(b.switchpointPut) == 0 && len(b.switchpointPatch) == 0 && len(b.switchpointDelete) == 0 &&
		len(b.servicePortProfilePut) == 0 && len(b.servicePortProfilePatch) == 0 && len(b.servicePortProfileDelete) == 0 &&
		len(b.packetBrokerPut) == 0 && len(b.packetBrokerPatch) == 0 && len(b.packetBrokerDelete) == 0 &&
		len(b.packetQueuePut) == 0 && len(b.packetQueuePatch) == 0 && len(b.packetQueueDelete) == 0 &&
		len(b.deviceVoiceSettingsPut) == 0 && len(b.deviceVoiceSettingsPatch) == 0 && len(b.deviceVoiceSettingsDelete) == 0 &&
		len(b.deviceControllerPut) == 0 && len(b.deviceControllerPatch) == 0 && len(b.deviceControllerDelete) == 0 {
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

	bundlePutCount := len(b.bundlePut)
	bundlePatchCount := len(b.bundlePatch)
	bundleDeleteCount := len(b.bundleDelete)

	aclPutCount := len(b.aclPut)
	aclPatchCount := len(b.aclPatch)
	aclDeleteCount := len(b.aclDelete)

	authenticatedEthPortPutCount := len(b.authenticatedEthPortPut)
	authenticatedEthPortPatchCount := len(b.authenticatedEthPortPatch)
	authenticatedEthPortDeleteCount := len(b.authenticatedEthPortDelete)

	badgePutCount := len(b.badgePut)
	badgePatchCount := len(b.badgePatch)
	badgeDeleteCount := len(b.badgeDelete)

	voicePortProfilePutCount := len(b.voicePortProfilePut)
	voicePortProfilePatchCount := len(b.voicePortProfilePatch)
	voicePortProfileDeleteCount := len(b.voicePortProfileDelete)

	switchpointPutCount := len(b.switchpointPut)
	switchpointPatchCount := len(b.switchpointPatch)
	switchpointDeleteCount := len(b.switchpointDelete)

	servicePortProfilePutCount := len(b.servicePortProfilePut)
	servicePortProfilePatchCount := len(b.servicePortProfilePatch)
	servicePortProfileDeleteCount := len(b.servicePortProfileDelete)

	packetBrokerPutCount := len(b.packetBrokerPut)
	packetBrokerPatchCount := len(b.packetBrokerPatch)
	packetBrokerDeleteCount := len(b.packetBrokerDelete)

	packetQueuePutCount := len(b.packetQueuePut)
	packetQueuePatchCount := len(b.packetQueuePatch)
	packetQueueDeleteCount := len(b.packetQueueDelete)

	deviceVoiceSettingsPutCount := len(b.deviceVoiceSettingsPut)
	deviceVoiceSettingsPatchCount := len(b.deviceVoiceSettingsPatch)
	deviceVoiceSettingsDeleteCount := len(b.deviceVoiceSettingsDelete)

	deviceControllerPutCount := len(b.deviceControllerPut)
	deviceControllerPatchCount := len(b.deviceControllerPatch)
	deviceControllerDeleteCount := len(b.deviceControllerDelete)

	b.mutex.Unlock()

	totalCount := gatewayPutCount + gatewayPatchCount + gatewayDeleteCount +
		lagPutCount + lagPatchCount + lagDeleteCount +
		tenantPutCount + tenantPatchCount + tenantDeleteCount +
		servicePutCount + servicePatchCount + serviceDeleteCount +
		gatewayProfilePutCount + gatewayProfilePatchCount + gatewayProfileDeleteCount +
		ethPortProfilePutCount + ethPortProfilePatchCount + ethPortProfileDeleteCount +
		ethPortSettingsPutCount + ethPortSettingsPatchCount + ethPortSettingsDeleteCount +
		bundlePutCount + bundlePatchCount + bundleDeleteCount + aclPutCount + aclPatchCount + aclDeleteCount +
		badgePutCount + badgePatchCount + badgeDeleteCount +
		voicePortProfilePutCount + voicePortProfilePatchCount + voicePortProfileDeleteCount +
		switchpointPutCount + switchpointPatchCount + switchpointDeleteCount +
		servicePortProfilePutCount + servicePortProfilePatchCount + servicePortProfileDeleteCount +
		packetBrokerPutCount + packetBrokerPatchCount + packetBrokerDeleteCount +
		packetQueuePutCount + packetQueuePatchCount + packetQueueDeleteCount +
		deviceVoiceSettingsPutCount + deviceVoiceSettingsPatchCount + deviceVoiceSettingsDeleteCount +
		authenticatedEthPortPutCount + authenticatedEthPortPatchCount + authenticatedEthPortDeleteCount +
		deviceControllerPutCount + deviceControllerPatchCount + deviceControllerDeleteCount

	if totalCount > 0 {
		tflog.Debug(ctx, "Multiple operations detected, executing in sequence", map[string]interface{}{
			"gateway_put_count":                   gatewayPutCount,
			"gateway_patch_count":                 gatewayPatchCount,
			"gateway_delete_count":                gatewayDeleteCount,
			"lag_put_count":                       lagPutCount,
			"lag_patch_count":                     lagPatchCount,
			"lag_delete_count":                    lagDeleteCount,
			"tenant_put_count":                    tenantPutCount,
			"tenant_patch_count":                  tenantPatchCount,
			"tenant_delete_count":                 tenantDeleteCount,
			"service_put_count":                   servicePutCount,
			"service_patch_count":                 servicePatchCount,
			"service_delete_count":                serviceDeleteCount,
			"gateway_profile_put_count":           gatewayProfilePutCount,
			"gateway_profile_patch_count":         gatewayProfilePatchCount,
			"gateway_profile_delete_count":        gatewayProfileDeleteCount,
			"eth_port_profile_put_count":          ethPortProfilePutCount,
			"eth_port_profile_patch_count":        ethPortProfilePatchCount,
			"eth_port_profile_delete_count":       ethPortProfileDeleteCount,
			"eth_port_settings_put_count":         ethPortSettingsPutCount,
			"eth_port_settings_patch_count":       ethPortSettingsPatchCount,
			"eth_port_settings_delete_count":      ethPortSettingsDeleteCount,
			"bundle_put_count":                    bundlePutCount,
			"bundle_patch_count":                  bundlePatchCount,
			"bundle_delete_count":                 bundleDeleteCount,
			"acl_put_count":                       aclPutCount,
			"acl_patch_count":                     aclPatchCount,
			"acl_delete_count":                    aclDeleteCount,
			"badge_put_count":                     badgePutCount,
			"badge_patch_count":                   badgePatchCount,
			"badge_delete_count":                  badgeDeleteCount,
			"voice_port_profile_put_count":        voicePortProfilePutCount,
			"voice_port_profile_patch_count":      voicePortProfilePatchCount,
			"voice_port_profile_delete_count":     voicePortProfileDeleteCount,
			"switchpoint_put_count":               switchpointPutCount,
			"switchpoint_patch_count":             switchpointPatchCount,
			"switchpoint_delete_count":            switchpointDeleteCount,
			"service_port_profile_put_count":      servicePortProfilePutCount,
			"service_port_profile_patch_count":    servicePortProfilePatchCount,
			"service_port_profile_delete_count":   servicePortProfileDeleteCount,
			"packet_broker_put_count":             packetBrokerPutCount,
			"packet_broker_patch_count":           packetBrokerPatchCount,
			"packet_broker_delete_count":          packetBrokerDeleteCount,
			"packet_queue_put_count":              packetQueuePutCount,
			"packet_queue_patch_count":            packetQueuePatchCount,
			"packet_queue_delete_count":           packetQueueDeleteCount,
			"device_voice_settings_put_count":     deviceVoiceSettingsPutCount,
			"device_voice_settings_patch_count":   deviceVoiceSettingsPatchCount,
			"device_voice_settings_delete_count":  deviceVoiceSettingsDeleteCount,
			"authenticated_eth_port_put_count":    authenticatedEthPortPutCount,
			"authenticated_eth_port_patch_count":  authenticatedEthPortPatchCount,
			"authenticated_eth_port_delete_count": authenticatedEthPortDeleteCount,
			"device_controller_put_count":         deviceControllerPutCount,
			"device_controller_patch_count":       deviceControllerPatchCount,
			"device_controller_delete_count":      deviceControllerDeleteCount,
			"total_count":                         totalCount,
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
		putLen = len(b.bundlePut)
		patchLen = len(b.bundlePatch)
		deleteLen = len(b.bundleDelete)
		recentOps = b.recentBundleOps
		recentOpTime = b.recentBundleOpTime
	case "acl":
		putLen = len(b.aclPut)
		patchLen = len(b.aclPatch)
		deleteLen = len(b.aclDelete)
		recentOps = b.recentAclOps
		recentOpTime = b.recentAclOpTime
	case "authenticated_eth_port":
		putLen = len(b.authenticatedEthPortPut)
		patchLen = len(b.authenticatedEthPortPatch)
		deleteLen = len(b.authenticatedEthPortDelete)
		recentOps = b.recentAuthenticatedEthPortOps
		recentOpTime = b.recentAuthenticatedEthPortOpTime
	case "badge":
		putLen = len(b.badgePut)
		patchLen = len(b.badgePatch)
		deleteLen = len(b.badgeDelete)
		recentOps = b.recentBadgeOps
		recentOpTime = b.recentBadgeOpTime
	case "voice_port_profile":
		putLen = len(b.voicePortProfilePut)
		patchLen = len(b.voicePortProfilePatch)
		deleteLen = len(b.voicePortProfileDelete)
		recentOps = b.recentVoicePortProfileOps
		recentOpTime = b.recentVoicePortProfileOpTime
	case "switchpoint":
		putLen = len(b.switchpointPut)
		patchLen = len(b.switchpointPatch)
		deleteLen = len(b.switchpointDelete)
		recentOps = b.recentSwitchpointOps
		recentOpTime = b.recentSwitchpointOpTime
	case "service_port_profile":
		putLen = len(b.servicePortProfilePut)
		patchLen = len(b.servicePortProfilePatch)
		deleteLen = len(b.servicePortProfileDelete)
		recentOps = b.recentServicePortProfileOps
		recentOpTime = b.recentServicePortProfileOpTime
	case "packet_broker":
		putLen = len(b.packetBrokerPut)
		patchLen = len(b.packetBrokerPatch)
		deleteLen = len(b.packetBrokerDelete)
		recentOps = b.recentPacketBrokerOps
		recentOpTime = b.recentPacketBrokerOpTime
	case "packet_queue":
		putLen = len(b.packetQueuePut)
		patchLen = len(b.packetQueuePatch)
		deleteLen = len(b.packetQueueDelete)
		recentOps = b.recentPacketQueueOps
		recentOpTime = b.recentPacketQueueOpTime
	case "device_voice_settings":
		putLen = len(b.deviceVoiceSettingsPut)
		patchLen = len(b.deviceVoiceSettingsPatch)
		deleteLen = len(b.deviceVoiceSettingsDelete)
		recentOps = b.recentDeviceVoiceSettingsOps
		recentOpTime = b.recentDeviceVoiceSettingsOpTime
	case "device_controller":
		putLen = len(b.deviceControllerPut)
		patchLen = len(b.deviceControllerPatch)
		deleteLen = len(b.deviceControllerDelete)
		recentOps = b.recentDeviceControllerOps
		recentOpTime = b.recentDeviceControllerOpTime
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

func (b *BulkOperationManager) HasPendingOrRecentAclOperations() bool {
	return b.hasPendingOrRecentOperations("acl")
}

func (b *BulkOperationManager) HasPendingOrRecentAuthenticatedEthPortOperations() bool {
	return b.hasPendingOrRecentOperations("authenticated_eth_port")
}

func (b *BulkOperationManager) HasPendingOrRecentBadgeOperations() bool {
	return b.hasPendingOrRecentOperations("badge")
}

func (b *BulkOperationManager) HasPendingOrRecentVoicePortProfileOperations() bool {
	return b.hasPendingOrRecentOperations("voice_port_profile")
}

func (b *BulkOperationManager) HasPendingOrRecentSwitchpointOperations() bool {
	return b.hasPendingOrRecentOperations("switchpoint")
}

func (b *BulkOperationManager) HasPendingOrRecentServicePortProfileOperations() bool {
	return b.hasPendingOrRecentOperations("service_port_profile")
}

func (b *BulkOperationManager) HasPendingOrRecentPacketBrokerOperations() bool {
	return b.hasPendingOrRecentOperations("packet_broker")
}

func (b *BulkOperationManager) HasPendingOrRecentPacketQueueOperations() bool {
	return b.hasPendingOrRecentOperations("packet_queue")
}

func (b *BulkOperationManager) HasPendingOrRecentDeviceVoiceSettingsOperations() bool {
	return b.hasPendingOrRecentOperations("device_voice_settings")
}

func (b *BulkOperationManager) HasPendingOrRecentDeviceControllerOperations() bool {
	return b.hasPendingOrRecentOperations("device_controller")
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

func (b *BulkOperationManager) AddBundlePut(ctx context.Context, bundleName string, props openapi.BundlesPutRequestEndpointBundleValue) string {
	return b.addOperation(
		ctx,
		"bundle",
		bundleName,
		"PUT",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.bundlePut[bundleName] = props
		},
		map[string]interface{}{
			"bundle_name": bundleName,
			"batch_size":  len(b.bundlePut) + 1,
		},
	)
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

func (b *BulkOperationManager) AddBundleDelete(ctx context.Context, bundleName string) string {
	return b.addOperation(
		ctx,
		"bundle",
		bundleName,
		"DELETE",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.bundleDelete = append(b.bundleDelete, bundleName)
		},
		map[string]interface{}{
			"bundle_name": bundleName,
			"batch_size":  len(b.bundleDelete) + 1,
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

func (b *BulkOperationManager) AddAclPut(ctx context.Context, aclName string, props openapi.ConfigPutRequestIpv4FilterIpv4FilterName, ipVersion string) string {
	resourceType := "acl_v4"
	if ipVersion == "6" {
		resourceType = "acl_v6"
	}

	return b.addOperation(
		ctx,
		resourceType,
		aclName,
		"PUT",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.aclPut[aclName] = props
			b.aclIpVersion[aclName] = ipVersion
		},
		map[string]interface{}{
			"acl_name":   aclName,
			"ip_version": ipVersion,
			"batch_size": len(b.aclPut) + 1,
		},
	)
}

func (b *BulkOperationManager) AddAclPatch(ctx context.Context, aclName string, props openapi.ConfigPutRequestIpv4FilterIpv4FilterName, ipVersion string) string {
	resourceType := "acl_v4"
	if ipVersion == "6" {
		resourceType = "acl_v6"
	}

	return b.addOperation(
		ctx,
		resourceType,
		aclName,
		"PATCH",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.aclPatch[aclName] = props
			b.aclIpVersion[aclName] = ipVersion
		},
		map[string]interface{}{
			"acl_name":   aclName,
			"ip_version": ipVersion,
			"batch_size": len(b.aclPatch) + 1,
		},
	)
}

func (b *BulkOperationManager) AddAclDelete(ctx context.Context, aclName string, ipVersion string) string {
	resourceType := "acl_v4"
	if ipVersion == "6" {
		resourceType = "acl_v6"
	}

	return b.addOperation(
		ctx,
		resourceType,
		aclName,
		"DELETE",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.aclDelete = append(b.aclDelete, aclName)
			b.aclIpVersion[aclName] = ipVersion
		},
		map[string]interface{}{
			"acl_name":   aclName,
			"ip_version": ipVersion,
			"batch_size": len(b.aclDelete) + 1,
		},
	)
}

func (b *BulkOperationManager) AddAuthenticatedEthPortPut(ctx context.Context, portName string, props openapi.ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) string {
	return b.addOperation(
		ctx,
		"authenticated_eth_port",
		portName,
		"PUT",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.authenticatedEthPortPut[portName] = props
		},
		map[string]interface{}{
			"port_name":  portName,
			"batch_size": len(b.authenticatedEthPortPut) + 1,
		},
	)
}

func (b *BulkOperationManager) AddAuthenticatedEthPortPatch(ctx context.Context, portName string, props openapi.ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName) string {
	return b.addOperation(
		ctx,
		"authenticated_eth_port",
		portName,
		"PATCH",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.authenticatedEthPortPatch[portName] = props
		},
		map[string]interface{}{
			"port_name":  portName,
			"batch_size": len(b.authenticatedEthPortPatch) + 1,
		},
	)
}

func (b *BulkOperationManager) AddAuthenticatedEthPortDelete(ctx context.Context, portName string) string {
	return b.addOperation(
		ctx,
		"authenticated_eth_port",
		portName,
		"DELETE",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.authenticatedEthPortDelete = append(b.authenticatedEthPortDelete, portName)
		},
		map[string]interface{}{
			"port_name":  portName,
			"batch_size": len(b.authenticatedEthPortDelete) + 1,
		},
	)
}

func (b *BulkOperationManager) AddBadgePut(ctx context.Context, badgeName string, props openapi.ConfigPutRequestBadgeBadgeName) string {
	return b.addOperation(
		ctx,
		"badge",
		badgeName,
		"PUT",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.badgePut[badgeName] = props
		},
		map[string]interface{}{
			"badge_name": badgeName,
			"batch_size": len(b.badgePut) + 1,
		},
	)
}

func (b *BulkOperationManager) AddBadgePatch(ctx context.Context, badgeName string, props openapi.ConfigPutRequestBadgeBadgeName) string {
	return b.addOperation(
		ctx,
		"badge",
		badgeName,
		"PATCH",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.badgePatch[badgeName] = props
		},
		map[string]interface{}{
			"badge_name": badgeName,
			"batch_size": len(b.badgePatch) + 1,
		},
	)
}

func (b *BulkOperationManager) AddBadgeDelete(ctx context.Context, badgeName string) string {
	return b.addOperation(
		ctx,
		"badge",
		badgeName,
		"DELETE",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.badgeDelete = append(b.badgeDelete, badgeName)
		},
		map[string]interface{}{
			"badge_name": badgeName,
			"batch_size": len(b.badgeDelete) + 1,
		},
	)
}

func (b *BulkOperationManager) AddDeviceVoiceSettingsPut(ctx context.Context, deviceVoiceSettingsName string, props openapi.ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) string {
	return b.addOperation(
		ctx,
		"device_voice_settings",
		deviceVoiceSettingsName,
		"PUT",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.deviceVoiceSettingsPut[deviceVoiceSettingsName] = props
		},
		map[string]interface{}{
			"device_voice_settings_name": deviceVoiceSettingsName,
			"batch_size":                 len(b.deviceVoiceSettingsPut) + 1,
		},
	)
}

func (b *BulkOperationManager) AddDeviceVoiceSettingsPatch(ctx context.Context, deviceVoiceSettingsName string, props openapi.ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName) string {
	return b.addOperation(
		ctx,
		"device_voice_settings",
		deviceVoiceSettingsName,
		"PATCH",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.deviceVoiceSettingsPatch[deviceVoiceSettingsName] = props
		},
		map[string]interface{}{
			"device_voice_settings_name": deviceVoiceSettingsName,
			"batch_size":                 len(b.deviceVoiceSettingsPatch) + 1,
		},
	)
}

func (b *BulkOperationManager) AddDeviceVoiceSettingsDelete(ctx context.Context, deviceVoiceSettingsName string) string {
	return b.addOperation(
		ctx,
		"device_voice_settings",
		deviceVoiceSettingsName,
		"DELETE",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.deviceVoiceSettingsDelete = append(b.deviceVoiceSettingsDelete, deviceVoiceSettingsName)
		},
		map[string]interface{}{
			"device_voice_settings_name": deviceVoiceSettingsName,
			"batch_size":                 len(b.deviceVoiceSettingsDelete) + 1,
		},
	)
}

func (b *BulkOperationManager) AddPacketBrokerPut(ctx context.Context, packetBrokerName string, props openapi.ConfigPutRequestPbEgressProfilePbEgressProfileName) string {
	return b.addOperation(
		ctx,
		"packet_broker",
		packetBrokerName,
		"PUT",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.packetBrokerPut[packetBrokerName] = props
		},
		map[string]interface{}{
			"packet_broker_name": packetBrokerName,
			"batch_size":         len(b.packetBrokerPut) + 1,
		},
	)
}

func (b *BulkOperationManager) AddPacketBrokerPatch(ctx context.Context, packetBrokerName string, props openapi.ConfigPutRequestPbEgressProfilePbEgressProfileName) string {
	return b.addOperation(
		ctx,
		"packet_broker",
		packetBrokerName,
		"PATCH",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.packetBrokerPatch[packetBrokerName] = props
		},
		map[string]interface{}{
			"packet_broker_name": packetBrokerName,
			"batch_size":         len(b.packetBrokerPatch) + 1,
		},
	)
}

func (b *BulkOperationManager) AddPacketBrokerDelete(ctx context.Context, packetBrokerName string) string {
	return b.addOperation(
		ctx,
		"packet_broker",
		packetBrokerName,
		"DELETE",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.packetBrokerDelete = append(b.packetBrokerDelete, packetBrokerName)
		},
		map[string]interface{}{
			"packet_broker_name": packetBrokerName,
			"batch_size":         len(b.packetBrokerDelete) + 1,
		},
	)
}

func (b *BulkOperationManager) AddPacketQueuePut(ctx context.Context, packetQueueName string, props openapi.ConfigPutRequestPacketQueuePacketQueueName) string {
	return b.addOperation(
		ctx,
		"packet_queue",
		packetQueueName,
		"PUT",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.packetQueuePut[packetQueueName] = props
		},
		map[string]interface{}{
			"packet_queue_name": packetQueueName,
			"batch_size":        len(b.packetQueuePut) + 1,
		},
	)
}

func (b *BulkOperationManager) AddPacketQueuePatch(ctx context.Context, packetQueueName string, props openapi.ConfigPutRequestPacketQueuePacketQueueName) string {
	return b.addOperation(
		ctx,
		"packet_queue",
		packetQueueName,
		"PATCH",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.packetQueuePatch[packetQueueName] = props
		},
		map[string]interface{}{
			"packet_queue_name": packetQueueName,
			"batch_size":        len(b.packetQueuePatch) + 1,
		},
	)
}

func (b *BulkOperationManager) AddPacketQueueDelete(ctx context.Context, packetQueueName string) string {
	return b.addOperation(
		ctx,
		"packet_queue",
		packetQueueName,
		"DELETE",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.packetQueueDelete = append(b.packetQueueDelete, packetQueueName)
		},
		map[string]interface{}{
			"packet_queue_name": packetQueueName,
			"batch_size":        len(b.packetQueueDelete) + 1,
		},
	)
}

func (b *BulkOperationManager) AddServicePortProfilePut(ctx context.Context, servicePortProfileName string, props openapi.ConfigPutRequestServicePortProfileServicePortProfileName) string {
	return b.addOperation(
		ctx,
		"service_port_profile",
		servicePortProfileName,
		"PUT",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.servicePortProfilePut[servicePortProfileName] = props
		},
		map[string]interface{}{
			"service_port_profile_name": servicePortProfileName,
			"batch_size":                len(b.servicePortProfilePut) + 1,
		},
	)
}

func (b *BulkOperationManager) AddServicePortProfilePatch(ctx context.Context, servicePortProfileName string, props openapi.ConfigPutRequestServicePortProfileServicePortProfileName) string {
	return b.addOperation(
		ctx,
		"service_port_profile",
		servicePortProfileName,
		"PATCH",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.servicePortProfilePatch[servicePortProfileName] = props
		},
		map[string]interface{}{
			"service_port_profile_name": servicePortProfileName,
			"batch_size":                len(b.servicePortProfilePatch) + 1,
		},
	)
}

func (b *BulkOperationManager) AddServicePortProfileDelete(ctx context.Context, servicePortProfileName string) string {
	return b.addOperation(
		ctx,
		"service_port_profile",
		servicePortProfileName,
		"DELETE",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.servicePortProfileDelete = append(b.servicePortProfileDelete, servicePortProfileName)
		},
		map[string]interface{}{
			"service_port_profile_name": servicePortProfileName,
			"batch_size":                len(b.servicePortProfileDelete) + 1,
		},
	)
}

func (b *BulkOperationManager) AddSwitchpointPut(ctx context.Context, switchpointName string, props openapi.ConfigPutRequestSwitchpointSwitchpointName) string {
	return b.addOperation(
		ctx,
		"switchpoint",
		switchpointName,
		"PUT",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.switchpointPut[switchpointName] = props
		},
		map[string]interface{}{
			"switchpoint_name": switchpointName,
			"batch_size":       len(b.switchpointPut) + 1,
		},
	)
}

func (b *BulkOperationManager) AddSwitchpointPatch(ctx context.Context, switchpointName string, props openapi.ConfigPutRequestSwitchpointSwitchpointName) string {
	return b.addOperation(
		ctx,
		"switchpoint",
		switchpointName,
		"PATCH",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.switchpointPatch[switchpointName] = props
		},
		map[string]interface{}{
			"switchpoint_name": switchpointName,
			"batch_size":       len(b.switchpointPatch) + 1,
		},
	)
}

func (b *BulkOperationManager) AddSwitchpointDelete(ctx context.Context, switchpointName string) string {
	return b.addOperation(
		ctx,
		"switchpoint",
		switchpointName,
		"DELETE",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.switchpointDelete = append(b.switchpointDelete, switchpointName)
		},
		map[string]interface{}{
			"switchpoint_name": switchpointName,
			"batch_size":       len(b.switchpointDelete) + 1,
		},
	)
}

func (b *BulkOperationManager) AddVoicePortProfilePut(ctx context.Context, voicePortProfileName string, props openapi.ConfigPutRequestVoicePortProfilesVoicePortProfilesName) string {
	return b.addOperation(
		ctx,
		"voice_port_profile",
		voicePortProfileName,
		"PUT",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.voicePortProfilePut[voicePortProfileName] = props
		},
		map[string]interface{}{
			"voice_port_profile_name": voicePortProfileName,
			"batch_size":              len(b.voicePortProfilePut) + 1,
		},
	)
}

func (b *BulkOperationManager) AddVoicePortProfilePatch(ctx context.Context, voicePortProfileName string, props openapi.ConfigPutRequestVoicePortProfilesVoicePortProfilesName) string {
	return b.addOperation(
		ctx,
		"voice_port_profile",
		voicePortProfileName,
		"PATCH",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.voicePortProfilePatch[voicePortProfileName] = props
		},
		map[string]interface{}{
			"voice_port_profile_name": voicePortProfileName,
			"batch_size":              len(b.voicePortProfilePatch) + 1,
		},
	)
}

func (b *BulkOperationManager) AddVoicePortProfileDelete(ctx context.Context, voicePortProfileName string) string {
	return b.addOperation(
		ctx,
		"voice_port_profile",
		voicePortProfileName,
		"DELETE",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.voicePortProfileDelete = append(b.voicePortProfileDelete, voicePortProfileName)
		},
		map[string]interface{}{
			"voice_port_profile_name": voicePortProfileName,
			"batch_size":              len(b.voicePortProfileDelete) + 1,
		},
	)
}

func (b *BulkOperationManager) AddDeviceControllerPut(ctx context.Context, deviceControllerName string, props openapi.ConfigPutRequestDeviceControllerDeviceControllerName) string {
	return b.addOperation(
		ctx,
		"device_controller",
		deviceControllerName,
		"PUT",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.deviceControllerPut[deviceControllerName] = props
		},
		map[string]interface{}{
			"device_controller_name": deviceControllerName,
			"batch_size":             len(b.deviceControllerPut) + 1,
		},
	)
}

func (b *BulkOperationManager) AddDeviceControllerPatch(ctx context.Context, deviceControllerName string, props openapi.ConfigPutRequestDeviceControllerDeviceControllerName) string {
	return b.addOperation(
		ctx,
		"device_controller",
		deviceControllerName,
		"PATCH",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.deviceControllerPatch[deviceControllerName] = props
		},
		map[string]interface{}{
			"device_controller_name": deviceControllerName,
			"batch_size":             len(b.deviceControllerPatch) + 1,
		},
	)
}

func (b *BulkOperationManager) AddDeviceControllerDelete(ctx context.Context, deviceControllerName string) string {
	return b.addOperation(
		ctx,
		"device_controller",
		deviceControllerName,
		"DELETE",
		func() {
			b.mutex.Lock()
			defer b.mutex.Unlock()
			b.deviceControllerDelete = append(b.deviceControllerDelete, deviceControllerName)
		},
		map[string]interface{}{
			"device_controller_name": deviceControllerName,
			"batch_size":             len(b.deviceControllerDelete) + 1,
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

func (b *BulkOperationManager) ExecuteBulkBundlePut(ctx context.Context) diag.Diagnostics {
	var originalOperations map[string]openapi.BundlesPutRequestEndpointBundleValue

	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "bundle",
		OperationType: "PUT",

		ExtractOperations: func() (map[string]interface{}, []string) {
			b.mutex.Lock()
			originalOperations = make(map[string]openapi.BundlesPutRequestEndpointBundleValue)
			for k, v := range b.bundlePut {
				originalOperations[k] = v
			}
			b.bundlePut = make(map[string]openapi.BundlesPutRequestEndpointBundleValue)
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
				ResourceType:  "bundle",
				OperationType: "PUT",
				FetchResources: func(ctx context.Context) (map[string]interface{}, error) {
					// First check if we have cached bundle data
					b.bundleResponsesMutex.RLock()
					if len(b.bundleResponses) > 0 {
						cachedData := make(map[string]interface{})
						for k, v := range b.bundleResponses {
							cachedData[k] = v
						}
						b.bundleResponsesMutex.RUnlock()

						tflog.Debug(ctx, "Using cached bundle data for pre-existence check", map[string]interface{}{
							"count": len(cachedData),
						})

						return cachedData, nil
					}
					b.bundleResponsesMutex.RUnlock()

					// Fall back to API call if no cache
					apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
					defer cancel()

					resp, err := b.client.BundlesAPI.BundlesGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						EndpointBundle map[string]interface{} `json:"endpoint_bundle"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}

					b.bundleResponsesMutex.Lock()
					for k, v := range result.EndpointBundle {
						if vMap, ok := v.(map[string]interface{}); ok {
							b.bundleResponses[k] = vMap

							if name, ok := vMap["name"].(string); ok && name != k {
								b.bundleResponses[name] = vMap
							}
						}
					}
					b.bundleResponsesMutex.Unlock()

					return result.EndpointBundle, nil
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
			putRequest := openapi.NewBundlesPutRequest()
			bundleMap := make(map[string]openapi.BundlesPutRequestEndpointBundleValue)

			for name, props := range filteredData {
				bundleMap[name] = props.(openapi.BundlesPutRequestEndpointBundleValue)
			}
			putRequest.SetEndpointBundle(bundleMap)
			return putRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.BundlesAPI.BundlesPut(ctx).BundlesPutRequest(
				*request.(*openapi.BundlesPutRequest))
			return req.Execute()
		},

		ProcessResponse: func(ctx context.Context, resp *http.Response) error {
			delayTime := 2 * time.Second
			tflog.Debug(ctx, fmt.Sprintf("Waiting %v for auto-generated values to be assigned before fetching bundles", delayTime))
			time.Sleep(delayTime)

			fetchCtx, fetchCancel := context.WithTimeout(context.Background(), OperationTimeout)
			defer fetchCancel()

			tflog.Debug(ctx, "Fetching bundles after successful PUT operation to retrieve auto-generated values")
			bundlesReq := b.client.BundlesAPI.BundlesGet(fetchCtx)
			bundlesResp, fetchErr := bundlesReq.Execute()

			if fetchErr != nil {
				tflog.Error(ctx, "Failed to fetch bundles after PUT for auto-generated fields", map[string]interface{}{
					"error": fetchErr.Error(),
				})
				return fetchErr
			}

			defer bundlesResp.Body.Close()

			var bundlesData struct {
				EndpointBundle map[string]map[string]interface{} `json:"endpoint_bundle"`
			}

			if respErr := json.NewDecoder(bundlesResp.Body).Decode(&bundlesData); respErr != nil {
				tflog.Error(ctx, "Failed to decode bundles response for auto-generated fields", map[string]interface{}{
					"error": respErr.Error(),
				})
				return respErr
			}

			b.bundleResponsesMutex.Lock()
			for bundleName, bundleData := range bundlesData.EndpointBundle {
				b.bundleResponses[bundleName] = bundleData

				if name, ok := bundleData["name"].(string); ok && name != bundleName {
					b.bundleResponses[name] = bundleData
				}
			}
			b.bundleResponsesMutex.Unlock()

			tflog.Debug(ctx, "Successfully stored bundle data for auto-generated fields", map[string]interface{}{
				"bundle_count": len(bundlesData.EndpointBundle),
			})

			return nil
		},

		UpdateRecentOps: func() {
			b.recentBundleOps = true
			b.recentBundleOpTime = time.Now()
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

func (b *BulkOperationManager) ExecuteBulkBundleDelete(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "bundle",
		OperationType: "DELETE",

		ExtractOperations: func() (map[string]interface{}, []string) {
			b.mutex.Lock()
			bundleNames := make([]string, len(b.bundleDelete))
			copy(bundleNames, b.bundleDelete)

			bundleDeleteMap := make(map[string]bool)
			for _, name := range bundleNames {
				bundleDeleteMap[name] = true
			}

			b.bundleDelete = make([]string, 0)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			for _, name := range bundleNames {
				result[name] = true
			}

			return result, bundleNames
		},

		CheckPreExistence: nil,

		PrepareRequest: func(filteredData map[string]interface{}) interface{} {
			bundleNames := make([]string, 0, len(filteredData))
			for name := range filteredData {
				bundleNames = append(bundleNames, name)
			}
			return bundleNames
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.BundlesAPI.BundlesDelete(ctx).BundleName(request.([]string))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentBundleOps = true
			b.recentBundleOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkAclPut(ctx context.Context) diag.Diagnostics {
	var diagnostics diag.Diagnostics
	b.mutex.Lock()
	originalOperations := make(map[string]openapi.ConfigPutRequestIpv4FilterIpv4FilterName)
	originalIpVersions := make(map[string]string)
	for k, v := range b.aclPut {
		originalOperations[k] = v
	}
	for k, v := range b.aclIpVersion {
		if _, exists := originalOperations[k]; exists {
			originalIpVersions[k] = v
		}
	}
	b.aclPut = make(map[string]openapi.ConfigPutRequestIpv4FilterIpv4FilterName)
	// Clean up IP version tracking for processed ACLs
	for k := range originalOperations {
		delete(b.aclIpVersion, k)
	}
	b.mutex.Unlock()

	if len(originalOperations) == 0 {
		return diagnostics
	}

	ipv4Data := make(map[string]openapi.ConfigPutRequestIpv4FilterIpv4FilterName)
	ipv6Data := make(map[string]openapi.ConfigPutRequestIpv4FilterIpv4FilterName)

	for name, props := range originalOperations {
		ipVersion := originalIpVersions[name]
		if ipVersion == "6" {
			ipv6Data[name] = props
		} else {
			ipv4Data[name] = props
		}
	}

	// Process IPv4 ACLs
	if len(ipv4Data) > 0 {
		ipv4Diagnostics := b.executeBulkOperation(ctx, BulkOperationConfig{
			ResourceType:  "acl_v4",
			OperationType: "PUT",

			ExtractOperations: func() (map[string]interface{}, []string) {
				result := make(map[string]interface{})
				names := make([]string, 0, len(ipv4Data))
				for k, v := range ipv4Data {
					result[k] = v
					names = append(names, k)
				}
				return result, names
			},

			CheckPreExistence: func(ctx context.Context, resourceNames []string) ([]string, map[string]interface{}, error) {
				checker := ResourceExistenceCheck{
					ResourceType:  "acl",
					OperationType: "PUT",
					FetchResources: func(ctx context.Context) (map[string]interface{}, error) {
						// Check cached ACL data first
						b.aclResponsesMutex.RLock()
						if len(b.aclResponses) > 0 {
							cachedData := make(map[string]interface{})
							for k, v := range b.aclResponses {
								cachedData[k] = v
							}
							b.aclResponsesMutex.RUnlock()
							return cachedData, nil
						}
						b.aclResponsesMutex.RUnlock()

						// Fetch IPv4 ACLs from API
						apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
						defer cancel()

						resp, err := b.client.ACLsAPI.AclsGet(apiCtx).IpVersion("4").Execute()
						if err != nil {
							return nil, err
						}
						defer resp.Body.Close()

						var result map[string]interface{}
						if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
							return nil, err
						}

						if ipv4Filter, ok := result["ipv4_filter"].(map[string]interface{}); ok {
							b.aclResponsesMutex.Lock()
							for k, v := range ipv4Filter {
								if vMap, ok := v.(map[string]interface{}); ok {
									b.aclResponses[k] = vMap
								}
							}
							b.aclResponsesMutex.Unlock()
							return ipv4Filter, nil
						}
						return make(map[string]interface{}), nil
					},
				}

				filteredNames, err := b.FilterPreExistingResources(ctx, resourceNames, checker)
				if err != nil {
					return resourceNames, nil, err
				}

				filteredOperations := make(map[string]interface{})
				for _, name := range filteredNames {
					if val, ok := ipv4Data[name]; ok {
						filteredOperations[name] = val
					}
				}
				return filteredNames, filteredOperations, nil
			},

			PrepareRequest: func(filteredData map[string]interface{}) interface{} {
				putRequest := openapi.NewAclsPutRequest()
				aclMap := make(map[string]openapi.ConfigPutRequestIpv4FilterIpv4FilterName)
				for name, props := range filteredData {
					aclMap[name] = props.(openapi.ConfigPutRequestIpv4FilterIpv4FilterName)
				}
				putRequest.SetIpFilter(aclMap)
				return putRequest
			},

			ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
				req := b.client.ACLsAPI.AclsPut(ctx).IpVersion("4").AclsPutRequest(
					*request.(*openapi.AclsPutRequest))
				return req.Execute()
			},

			ProcessResponse: func(ctx context.Context, resp *http.Response) error {
				// Fetch IPv4 ACLs after PUT to get auto-generated values
				delayTime := 2 * time.Second
				time.Sleep(delayTime)

				fetchCtx, fetchCancel := context.WithTimeout(context.Background(), OperationTimeout)
				defer fetchCancel()

				aclsReq := b.client.ACLsAPI.AclsGet(fetchCtx).IpVersion("4")
				aclsResp, fetchErr := aclsReq.Execute()
				if fetchErr != nil {
					return fetchErr
				}
				defer aclsResp.Body.Close()

				var aclsResponse map[string]interface{}
				if respErr := json.NewDecoder(aclsResp.Body).Decode(&aclsResponse); respErr != nil {
					return respErr
				}

				if ipv4Filter, ok := aclsResponse["ipv4_filter"].(map[string]interface{}); ok {
					b.aclResponsesMutex.Lock()
					for k, v := range ipv4Filter {
						if vMap, ok := v.(map[string]interface{}); ok {
							b.aclResponses[k] = vMap
						}
					}
					b.aclResponsesMutex.Unlock()
				}
				return nil
			},

			UpdateRecentOps: func() {
				b.recentAclOps = true
				b.recentAclOpTime = time.Now()
			},
		})
		diagnostics = append(diagnostics, ipv4Diagnostics...)
	}

	// Process IPv6 ACLs
	if len(ipv6Data) > 0 {
		ipv6Diagnostics := b.executeBulkOperation(ctx, BulkOperationConfig{
			ResourceType:  "acl_v6",
			OperationType: "PUT",

			ExtractOperations: func() (map[string]interface{}, []string) {
				result := make(map[string]interface{})
				names := make([]string, 0, len(ipv6Data))
				for k, v := range ipv6Data {
					result[k] = v
					names = append(names, k)
				}
				return result, names
			},

			CheckPreExistence: func(ctx context.Context, resourceNames []string) ([]string, map[string]interface{}, error) {
				checker := ResourceExistenceCheck{
					ResourceType:  "acl",
					OperationType: "PUT",
					FetchResources: func(ctx context.Context) (map[string]interface{}, error) {
						// Check cached ACL data first
						b.aclResponsesMutex.RLock()
						if len(b.aclResponses) > 0 {
							cachedData := make(map[string]interface{})
							for k, v := range b.aclResponses {
								cachedData[k] = v
							}
							b.aclResponsesMutex.RUnlock()
							return cachedData, nil
						}
						b.aclResponsesMutex.RUnlock()

						// Fetch IPv6 ACLs from API
						apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
						defer cancel()

						resp, err := b.client.ACLsAPI.AclsGet(apiCtx).IpVersion("6").Execute()
						if err != nil {
							return nil, err
						}
						defer resp.Body.Close()

						var result map[string]interface{}
						if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
							return nil, err
						}

						if ipv6Filter, ok := result["ipv6_filter"].(map[string]interface{}); ok {
							b.aclResponsesMutex.Lock()
							for k, v := range ipv6Filter {
								if vMap, ok := v.(map[string]interface{}); ok {
									b.aclResponses[k] = vMap
								}
							}
							b.aclResponsesMutex.Unlock()
							return ipv6Filter, nil
						}
						return make(map[string]interface{}), nil
					},
				}

				filteredNames, err := b.FilterPreExistingResources(ctx, resourceNames, checker)
				if err != nil {
					return resourceNames, nil, err
				}

				filteredOperations := make(map[string]interface{})
				for _, name := range filteredNames {
					if val, ok := ipv6Data[name]; ok {
						filteredOperations[name] = val
					}
				}
				return filteredNames, filteredOperations, nil
			},

			PrepareRequest: func(filteredData map[string]interface{}) interface{} {
				putRequest := openapi.NewAclsPutRequest()
				aclMap := make(map[string]openapi.ConfigPutRequestIpv4FilterIpv4FilterName)
				for name, props := range filteredData {
					aclMap[name] = props.(openapi.ConfigPutRequestIpv4FilterIpv4FilterName)
				}
				putRequest.SetIpFilter(aclMap)
				return putRequest
			},

			ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
				req := b.client.ACLsAPI.AclsPut(ctx).IpVersion("6").AclsPutRequest(
					*request.(*openapi.AclsPutRequest))
				return req.Execute()
			},

			ProcessResponse: func(ctx context.Context, resp *http.Response) error {
				// Fetch IPv6 ACLs after PUT to get auto-generated values
				delayTime := 2 * time.Second
				time.Sleep(delayTime)

				fetchCtx, fetchCancel := context.WithTimeout(context.Background(), OperationTimeout)
				defer fetchCancel()

				aclsReq := b.client.ACLsAPI.AclsGet(fetchCtx).IpVersion("6")
				aclsResp, fetchErr := aclsReq.Execute()
				if fetchErr != nil {
					return fetchErr
				}
				defer aclsResp.Body.Close()

				var aclsResponse map[string]interface{}
				if respErr := json.NewDecoder(aclsResp.Body).Decode(&aclsResponse); respErr != nil {
					return respErr
				}

				if ipv6Filter, ok := aclsResponse["ipv6_filter"].(map[string]interface{}); ok {
					b.aclResponsesMutex.Lock()
					for k, v := range ipv6Filter {
						if vMap, ok := v.(map[string]interface{}); ok {
							b.aclResponses[k] = vMap
						}
					}
					b.aclResponsesMutex.Unlock()
				}
				return nil
			},

			UpdateRecentOps: func() {
				b.recentAclOps = true
				b.recentAclOpTime = time.Now()
			},
		})
		diagnostics = append(diagnostics, ipv6Diagnostics...)
	}

	return diagnostics
}

func (b *BulkOperationManager) ExecuteBulkAclPatch(ctx context.Context) diag.Diagnostics {
	var diagnostics diag.Diagnostics
	b.mutex.Lock()
	originalOperations := make(map[string]openapi.ConfigPutRequestIpv4FilterIpv4FilterName)
	originalIpVersions := make(map[string]string)
	for k, v := range b.aclPatch {
		originalOperations[k] = v
	}
	for k, v := range b.aclIpVersion {
		if _, exists := originalOperations[k]; exists {
			originalIpVersions[k] = v
		}
	}
	b.aclPatch = make(map[string]openapi.ConfigPutRequestIpv4FilterIpv4FilterName)
	// Clean up IP version tracking for processed ACLs
	for k := range originalOperations {
		delete(b.aclIpVersion, k)
	}
	b.mutex.Unlock()

	if len(originalOperations) == 0 {
		return diagnostics
	}

	ipv4Data := make(map[string]openapi.ConfigPutRequestIpv4FilterIpv4FilterName)
	ipv6Data := make(map[string]openapi.ConfigPutRequestIpv4FilterIpv4FilterName)

	for name, props := range originalOperations {
		ipVersion := originalIpVersions[name]
		if ipVersion == "6" {
			ipv6Data[name] = props
		} else {
			ipv4Data[name] = props
		}
	}

	// Process IPv4 ACLs
	if len(ipv4Data) > 0 {
		ipv4Diagnostics := b.executeBulkOperation(ctx, BulkOperationConfig{
			ResourceType:  "acl_v4",
			OperationType: "PATCH",

			ExtractOperations: func() (map[string]interface{}, []string) {
				result := make(map[string]interface{})
				names := make([]string, 0, len(ipv4Data))
				for k, v := range ipv4Data {
					result[k] = v
					names = append(names, k)
				}
				return result, names
			},

			CheckPreExistence: nil,

			PrepareRequest: func(filteredData map[string]interface{}) interface{} {
				patchRequest := openapi.NewAclsPutRequest()
				aclMap := make(map[string]openapi.ConfigPutRequestIpv4FilterIpv4FilterName)
				for name, props := range filteredData {
					aclMap[name] = props.(openapi.ConfigPutRequestIpv4FilterIpv4FilterName)
				}
				patchRequest.SetIpFilter(aclMap)
				return patchRequest
			},

			ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
				req := b.client.ACLsAPI.AclsPatch(ctx).IpVersion("4").AclsPutRequest(
					*request.(*openapi.AclsPutRequest))
				return req.Execute()
			},

			UpdateRecentOps: func() {
				b.recentAclOps = true
				b.recentAclOpTime = time.Now()
			},
		})
		diagnostics = append(diagnostics, ipv4Diagnostics...)
	}

	// Process IPv6 ACLs
	if len(ipv6Data) > 0 {
		ipv6Diagnostics := b.executeBulkOperation(ctx, BulkOperationConfig{
			ResourceType:  "acl_v6",
			OperationType: "PATCH",

			ExtractOperations: func() (map[string]interface{}, []string) {
				result := make(map[string]interface{})
				names := make([]string, 0, len(ipv6Data))
				for k, v := range ipv6Data {
					result[k] = v
					names = append(names, k)
				}
				return result, names
			},

			CheckPreExistence: nil,

			PrepareRequest: func(filteredData map[string]interface{}) interface{} {
				patchRequest := openapi.NewAclsPutRequest()
				aclMap := make(map[string]openapi.ConfigPutRequestIpv4FilterIpv4FilterName)
				for name, props := range filteredData {
					aclMap[name] = props.(openapi.ConfigPutRequestIpv4FilterIpv4FilterName)
				}
				patchRequest.SetIpFilter(aclMap)
				return patchRequest
			},

			ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
				req := b.client.ACLsAPI.AclsPatch(ctx).IpVersion("6").AclsPutRequest(
					*request.(*openapi.AclsPutRequest))
				return req.Execute()
			},

			UpdateRecentOps: func() {
				b.recentAclOps = true
				b.recentAclOpTime = time.Now()
			},
		})
		diagnostics = append(diagnostics, ipv6Diagnostics...)
	}

	return diagnostics
}

func (b *BulkOperationManager) ExecuteBulkAclDelete(ctx context.Context) diag.Diagnostics {
	var diagnostics diag.Diagnostics
	b.mutex.Lock()
	aclNames := make([]string, len(b.aclDelete))
	copy(aclNames, b.aclDelete)

	originalIpVersions := make(map[string]string)
	for _, name := range aclNames {
		if ipVersion, exists := b.aclIpVersion[name]; exists {
			originalIpVersions[name] = ipVersion
		}
	}

	b.aclDelete = make([]string, 0)
	// Clean up IP version tracking for processed ACLs
	for _, name := range aclNames {
		delete(b.aclIpVersion, name)
	}
	b.mutex.Unlock()

	if len(aclNames) == 0 {
		return diagnostics
	}

	ipv4Names := make([]string, 0)
	ipv6Names := make([]string, 0)

	for _, name := range aclNames {
		ipVersion := originalIpVersions[name]
		if ipVersion == "6" {
			ipv6Names = append(ipv6Names, name)
		} else {
			ipv4Names = append(ipv4Names, name)
		}
	}

	// Process IPv4 ACLs
	if len(ipv4Names) > 0 {
		ipv4Diagnostics := b.executeBulkOperation(ctx, BulkOperationConfig{
			ResourceType:  "acl_v4",
			OperationType: "DELETE",

			ExtractOperations: func() (map[string]interface{}, []string) {
				result := make(map[string]interface{})
				for _, name := range ipv4Names {
					result[name] = true
				}
				return result, ipv4Names
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
				req := b.client.ACLsAPI.AclsDelete(ctx).
					IpFilterName(request.([]string)).
					IpVersion("4")
				return req.Execute()
			},

			UpdateRecentOps: func() {
				b.recentAclOps = true
				b.recentAclOpTime = time.Now()
			},
		})
		diagnostics = append(diagnostics, ipv4Diagnostics...)
	}

	// Process IPv6 ACLs
	if len(ipv6Names) > 0 {
		ipv6Diagnostics := b.executeBulkOperation(ctx, BulkOperationConfig{
			ResourceType:  "acl_v6",
			OperationType: "DELETE",

			ExtractOperations: func() (map[string]interface{}, []string) {
				result := make(map[string]interface{})
				for _, name := range ipv6Names {
					result[name] = true
				}
				return result, ipv6Names
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
				req := b.client.ACLsAPI.AclsDelete(ctx).
					IpFilterName(request.([]string)).
					IpVersion("6")
				return req.Execute()
			},

			UpdateRecentOps: func() {
				b.recentAclOps = true
				b.recentAclOpTime = time.Now()
			},
		})
		diagnostics = append(diagnostics, ipv6Diagnostics...)
	}

	return diagnostics
}

func (b *BulkOperationManager) ExecuteBulkAuthenticatedEthPortPut(ctx context.Context) diag.Diagnostics {
	var originalOperations map[string]openapi.ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName

	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "authenticated_eth_port",
		OperationType: "PUT",

		ExtractOperations: func() (map[string]interface{}, []string) {
			b.mutex.Lock()
			originalOperations = make(map[string]openapi.ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName)
			for k, v := range b.authenticatedEthPortPut {
				originalOperations[k] = v
			}
			b.authenticatedEthPortPut = make(map[string]openapi.ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName)
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
				ResourceType:  "authenticated_eth_port",
				OperationType: "PUT",
				FetchResources: func(ctx context.Context) (map[string]interface{}, error) {
					// First check if we have cached Authenticated Eth-Port data
					b.authenticatedEthPortResponsesMutex.RLock()
					if len(b.authenticatedEthPortResponses) > 0 {
						cachedData := make(map[string]interface{})
						for k, v := range b.authenticatedEthPortResponses {
							cachedData[k] = v
						}
						b.authenticatedEthPortResponsesMutex.RUnlock()

						tflog.Debug(ctx, "Using cached Authenticated Eth-Port data for pre-existence check", map[string]interface{}{
							"count": len(cachedData),
						})

						return cachedData, nil
					}
					b.authenticatedEthPortResponsesMutex.RUnlock()

					// Fall back to API call if no cache
					apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
					defer cancel()

					resp, err := b.client.AuthenticatedEthPortsAPI.AuthenticatedethportsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						AuthenticatedEthPort map[string]interface{} `json:"authenticated_eth_port"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}

					b.authenticatedEthPortResponsesMutex.Lock()
					for k, v := range result.AuthenticatedEthPort {
						if vMap, ok := v.(map[string]interface{}); ok {
							b.authenticatedEthPortResponses[k] = vMap

							if name, ok := vMap["name"].(string); ok && name != k {
								b.authenticatedEthPortResponses[name] = vMap
							}
						}
					}
					b.authenticatedEthPortResponsesMutex.Unlock()

					return result.AuthenticatedEthPort, nil
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
			putRequest := openapi.NewAuthenticatedethportsPutRequest()
			portMap := make(map[string]openapi.ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName)

			for name, props := range filteredData {
				portMap[name] = props.(openapi.ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName)
			}
			putRequest.SetAuthenticatedEthPort(portMap)
			return putRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.AuthenticatedEthPortsAPI.AuthenticatedethportsPut(ctx).AuthenticatedethportsPutRequest(
				*request.(*openapi.AuthenticatedethportsPutRequest))
			return req.Execute()
		},

		ProcessResponse: func(ctx context.Context, resp *http.Response) error {
			delayTime := 2 * time.Second
			tflog.Debug(ctx, fmt.Sprintf("Waiting %v for auto-generated values to be assigned before fetching Authenticated Eth-Ports", delayTime))
			time.Sleep(delayTime)

			fetchCtx, fetchCancel := context.WithTimeout(context.Background(), OperationTimeout)
			defer fetchCancel()

			tflog.Debug(ctx, "Fetching Authenticated Eth-Ports after successful PUT operation to retrieve auto-generated values")
			portsReq := b.client.AuthenticatedEthPortsAPI.AuthenticatedethportsGet(fetchCtx)
			portsResp, fetchErr := portsReq.Execute()

			if fetchErr != nil {
				tflog.Error(ctx, "Failed to fetch Authenticated Eth-Ports after PUT for auto-generated fields", map[string]interface{}{
					"error": fetchErr.Error(),
				})
				return fetchErr
			}

			defer portsResp.Body.Close()

			var portsData struct {
				AuthenticatedEthPort map[string]map[string]interface{} `json:"authenticated_eth_port"`
			}

			if respErr := json.NewDecoder(portsResp.Body).Decode(&portsData); respErr != nil {
				tflog.Error(ctx, "Failed to decode Authenticated Eth-Ports response for auto-generated fields", map[string]interface{}{
					"error": respErr.Error(),
				})
				return respErr
			}

			b.authenticatedEthPortResponsesMutex.Lock()
			for portName, portData := range portsData.AuthenticatedEthPort {
				b.authenticatedEthPortResponses[portName] = portData

				if name, ok := portData["name"].(string); ok && name != portName {
					b.authenticatedEthPortResponses[name] = portData
				}
			}
			b.authenticatedEthPortResponsesMutex.Unlock()

			tflog.Debug(ctx, "Successfully stored Authenticated Eth-Port data for auto-generated fields", map[string]interface{}{
				"authenticated_eth_port_count": len(portsData.AuthenticatedEthPort),
			})

			return nil
		},

		UpdateRecentOps: func() {
			b.recentAuthenticatedEthPortOps = true
			b.recentAuthenticatedEthPortOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkAuthenticatedEthPortPatch(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "authenticated_eth_port",
		OperationType: "PATCH",

		ExtractOperations: func() (map[string]interface{}, []string) {
			b.mutex.Lock()
			authenticatedEthPortPatch := make(map[string]openapi.ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName)
			for k, v := range b.authenticatedEthPortPatch {
				authenticatedEthPortPatch[k] = v
			}
			b.authenticatedEthPortPatch = make(map[string]openapi.ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			names := make([]string, 0, len(authenticatedEthPortPatch))

			for k, v := range authenticatedEthPortPatch {
				result[k] = v
				names = append(names, k)
			}

			return result, names
		},

		CheckPreExistence: nil,

		PrepareRequest: func(filteredData map[string]interface{}) interface{} {
			patchRequest := openapi.NewAuthenticatedethportsPutRequest()
			portMap := make(map[string]openapi.ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName)

			for name, props := range filteredData {
				portMap[name] = props.(openapi.ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName)
			}
			patchRequest.SetAuthenticatedEthPort(portMap)
			return patchRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.AuthenticatedEthPortsAPI.AuthenticatedethportsPatch(ctx).AuthenticatedethportsPutRequest(
				*request.(*openapi.AuthenticatedethportsPutRequest))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentAuthenticatedEthPortOps = true
			b.recentAuthenticatedEthPortOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkAuthenticatedEthPortDelete(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "authenticated_eth_port",
		OperationType: "DELETE",

		ExtractOperations: func() (map[string]interface{}, []string) {
			b.mutex.Lock()
			portNames := make([]string, len(b.authenticatedEthPortDelete))
			copy(portNames, b.authenticatedEthPortDelete)

			portDeleteMap := make(map[string]bool)
			for _, name := range portNames {
				portDeleteMap[name] = true
			}

			b.authenticatedEthPortDelete = make([]string, 0)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			for _, name := range portNames {
				result[name] = true
			}

			return result, portNames
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
			req := b.client.AuthenticatedEthPortsAPI.AuthenticatedethportsDelete(ctx).AuthenticatedEthPortName(request.([]string))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentAuthenticatedEthPortOps = true
			b.recentAuthenticatedEthPortOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkBadgePut(ctx context.Context) diag.Diagnostics {
	var originalOperations map[string]openapi.ConfigPutRequestBadgeBadgeName

	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "badge",
		OperationType: "PUT",

		ExtractOperations: func() (map[string]interface{}, []string) {
			b.mutex.Lock()
			originalOperations = make(map[string]openapi.ConfigPutRequestBadgeBadgeName)
			for k, v := range b.badgePut {
				originalOperations[k] = v
			}
			b.badgePut = make(map[string]openapi.ConfigPutRequestBadgeBadgeName)
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
				ResourceType:  "badge",
				OperationType: "PUT",
				FetchResources: func(ctx context.Context) (map[string]interface{}, error) {
					// First check if we have cached badge data
					b.badgeResponsesMutex.RLock()
					if len(b.badgeResponses) > 0 {
						cachedData := make(map[string]interface{})
						for k, v := range b.badgeResponses {
							cachedData[k] = v
						}
						b.badgeResponsesMutex.RUnlock()

						tflog.Debug(ctx, "Using cached badge data for pre-existence check", map[string]interface{}{
							"count": len(cachedData),
						})

						return cachedData, nil
					}
					b.badgeResponsesMutex.RUnlock()

					// Fall back to API call if no cache
					apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
					defer cancel()

					resp, err := b.client.BadgesAPI.BadgesGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						Badge map[string]interface{} `json:"badge"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}

					b.badgeResponsesMutex.Lock()
					for k, v := range result.Badge {
						if vMap, ok := v.(map[string]interface{}); ok {
							b.badgeResponses[k] = vMap

							if name, ok := vMap["name"].(string); ok && name != k {
								b.badgeResponses[name] = vMap
							}
						}
					}
					b.badgeResponsesMutex.Unlock()

					return result.Badge, nil
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
			putRequest := openapi.NewBadgesPutRequest()
			badgeMap := make(map[string]openapi.ConfigPutRequestBadgeBadgeName)

			for name, props := range filteredData {
				badgeMap[name] = props.(openapi.ConfigPutRequestBadgeBadgeName)
			}
			putRequest.SetBadge(badgeMap)
			return putRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.BadgesAPI.BadgesPut(ctx).BadgesPutRequest(
				*request.(*openapi.BadgesPutRequest))
			return req.Execute()
		},

		ProcessResponse: func(ctx context.Context, resp *http.Response) error {
			delayTime := 2 * time.Second
			tflog.Debug(ctx, fmt.Sprintf("Waiting %v for auto-generated values to be assigned before fetching badges", delayTime))
			time.Sleep(delayTime)

			fetchCtx, fetchCancel := context.WithTimeout(context.Background(), OperationTimeout)
			defer fetchCancel()

			tflog.Debug(ctx, "Fetching badges after successful PUT operation to retrieve auto-generated values")
			badgesReq := b.client.BadgesAPI.BadgesGet(fetchCtx)
			badgesResp, fetchErr := badgesReq.Execute()

			if fetchErr != nil {
				tflog.Error(ctx, "Failed to fetch badges after PUT for auto-generated fields", map[string]interface{}{
					"error": fetchErr.Error(),
				})
				return fetchErr
			}

			defer badgesResp.Body.Close()

			var badgesData struct {
				Badge map[string]map[string]interface{} `json:"badge"`
			}

			if respErr := json.NewDecoder(badgesResp.Body).Decode(&badgesData); respErr != nil {
				tflog.Error(ctx, "Failed to decode badges response for auto-generated fields", map[string]interface{}{
					"error": respErr.Error(),
				})
				return respErr
			}

			b.badgeResponsesMutex.Lock()
			for badgeName, badgeData := range badgesData.Badge {
				b.badgeResponses[badgeName] = badgeData

				if name, ok := badgeData["name"].(string); ok && name != badgeName {
					b.badgeResponses[name] = badgeData
				}
			}
			b.badgeResponsesMutex.Unlock()

			tflog.Debug(ctx, "Successfully stored badge data for auto-generated fields", map[string]interface{}{
				"badge_count": len(badgesData.Badge),
			})

			return nil
		},

		UpdateRecentOps: func() {
			b.recentBadgeOps = true
			b.recentBadgeOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkBadgePatch(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "badge",
		OperationType: "PATCH",

		ExtractOperations: func() (map[string]interface{}, []string) {
			b.mutex.Lock()
			badgePatch := make(map[string]openapi.ConfigPutRequestBadgeBadgeName)
			for k, v := range b.badgePatch {
				badgePatch[k] = v
			}
			b.badgePatch = make(map[string]openapi.ConfigPutRequestBadgeBadgeName)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			names := make([]string, 0, len(badgePatch))

			for k, v := range badgePatch {
				result[k] = v
				names = append(names, k)
			}

			return result, names
		},

		CheckPreExistence: nil,

		PrepareRequest: func(filteredData map[string]interface{}) interface{} {
			patchRequest := openapi.NewBadgesPutRequest()
			badgeMap := make(map[string]openapi.ConfigPutRequestBadgeBadgeName)

			for name, props := range filteredData {
				badgeMap[name] = props.(openapi.ConfigPutRequestBadgeBadgeName)
			}
			patchRequest.SetBadge(badgeMap)
			return patchRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.BadgesAPI.BadgesPatch(ctx).BadgesPutRequest(
				*request.(*openapi.BadgesPutRequest))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentBadgeOps = true
			b.recentBadgeOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkBadgeDelete(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "badge",
		OperationType: "DELETE",

		ExtractOperations: func() (map[string]interface{}, []string) {
			b.mutex.Lock()
			badgeNames := make([]string, len(b.badgeDelete))
			copy(badgeNames, b.badgeDelete)

			badgeDeleteMap := make(map[string]bool)
			for _, name := range badgeNames {
				badgeDeleteMap[name] = true
			}
			b.badgeDelete = make([]string, 0)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			for _, name := range badgeNames {
				result[name] = true
			}

			return result, badgeNames
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
			req := b.client.BadgesAPI.BadgesDelete(ctx).BadgeName(request.([]string))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentBadgeOps = true
			b.recentBadgeOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkDeviceVoiceSettingsPut(ctx context.Context) diag.Diagnostics {
	var originalOperations map[string]openapi.ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName

	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "device_voice_settings",
		OperationType: "PUT",

		ExtractOperations: func() (map[string]interface{}, []string) {
			b.mutex.Lock()
			originalOperations = make(map[string]openapi.ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName)
			for k, v := range b.deviceVoiceSettingsPut {
				originalOperations[k] = v
			}
			b.deviceVoiceSettingsPut = make(map[string]openapi.ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName)
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
				ResourceType:  "device_voice_settings",
				OperationType: "PUT",
				FetchResources: func(ctx context.Context) (map[string]interface{}, error) {
					// First check if we have cached device voice settings data
					b.deviceVoiceSettingsResponsesMutex.RLock()
					if len(b.deviceVoiceSettingsResponses) > 0 {
						cachedData := make(map[string]interface{})
						for k, v := range b.deviceVoiceSettingsResponses {
							cachedData[k] = v
						}
						b.deviceVoiceSettingsResponsesMutex.RUnlock()

						tflog.Debug(ctx, "Using cached device voice settings data for pre-existence check", map[string]interface{}{
							"count": len(cachedData),
						})

						return cachedData, nil
					}
					b.deviceVoiceSettingsResponsesMutex.RUnlock()

					// Fall back to API call if no cache
					apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
					defer cancel()

					resp, err := b.client.DeviceVoiceSettingsAPI.DevicevoicesettingsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						DeviceVoiceSettings map[string]interface{} `json:"device_voice_settings"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}

					b.deviceVoiceSettingsResponsesMutex.Lock()
					for k, v := range result.DeviceVoiceSettings {
						if vMap, ok := v.(map[string]interface{}); ok {
							b.deviceVoiceSettingsResponses[k] = vMap

							if name, ok := vMap["name"].(string); ok && name != k {
								b.deviceVoiceSettingsResponses[name] = vMap
							}
						}
					}
					b.deviceVoiceSettingsResponsesMutex.Unlock()

					return result.DeviceVoiceSettings, nil
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
			putRequest := openapi.NewDevicevoicesettingsPutRequest()
			deviceVoiceSettingsMap := make(map[string]openapi.ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName)

			for name, props := range filteredData {
				deviceVoiceSettingsMap[name] = props.(openapi.ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName)
			}
			putRequest.SetDeviceVoiceSettings(deviceVoiceSettingsMap)
			return putRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.DeviceVoiceSettingsAPI.DevicevoicesettingsPut(ctx).DevicevoicesettingsPutRequest(
				*request.(*openapi.DevicevoicesettingsPutRequest))
			return req.Execute()
		},

		ProcessResponse: func(ctx context.Context, resp *http.Response) error {
			delayTime := 2 * time.Second
			tflog.Debug(ctx, fmt.Sprintf("Waiting %v for auto-generated values to be assigned before fetching device voice settings", delayTime))
			time.Sleep(delayTime)

			fetchCtx, fetchCancel := context.WithTimeout(context.Background(), OperationTimeout)
			defer fetchCancel()

			tflog.Debug(ctx, "Fetching device voice settings after successful PUT operation to retrieve auto-generated values")
			settingsReq := b.client.DeviceVoiceSettingsAPI.DevicevoicesettingsGet(fetchCtx)
			settingsResp, fetchErr := settingsReq.Execute()

			if fetchErr != nil {
				tflog.Error(ctx, "Failed to fetch device voice settings after PUT for auto-generated fields", map[string]interface{}{
					"error": fetchErr.Error(),
				})
				return fetchErr
			}

			defer settingsResp.Body.Close()

			var settingsData struct {
				DeviceVoiceSettings map[string]map[string]interface{} `json:"device_voice_settings"`
			}

			if respErr := json.NewDecoder(settingsResp.Body).Decode(&settingsData); respErr != nil {
				tflog.Error(ctx, "Failed to decode device voice settings response for auto-generated fields", map[string]interface{}{
					"error": respErr.Error(),
				})
				return respErr
			}

			b.deviceVoiceSettingsResponsesMutex.Lock()
			for settingsName, settingsData := range settingsData.DeviceVoiceSettings {
				b.deviceVoiceSettingsResponses[settingsName] = settingsData

				if name, ok := settingsData["name"].(string); ok && name != settingsName {
					b.deviceVoiceSettingsResponses[name] = settingsData
				}
			}
			b.deviceVoiceSettingsResponsesMutex.Unlock()

			tflog.Debug(ctx, "Successfully stored device voice settings data for auto-generated fields", map[string]interface{}{
				"device_voice_settings_count": len(settingsData.DeviceVoiceSettings),
			})

			return nil
		},

		UpdateRecentOps: func() {
			b.recentDeviceVoiceSettingsOps = true
			b.recentDeviceVoiceSettingsOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkDeviceVoiceSettingsPatch(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "device_voice_settings",
		OperationType: "PATCH",

		ExtractOperations: func() (map[string]interface{}, []string) {
			b.mutex.Lock()
			deviceVoiceSettingsPatch := make(map[string]openapi.ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName)
			for k, v := range b.deviceVoiceSettingsPatch {
				deviceVoiceSettingsPatch[k] = v
			}
			b.deviceVoiceSettingsPatch = make(map[string]openapi.ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			names := make([]string, 0, len(deviceVoiceSettingsPatch))

			for k, v := range deviceVoiceSettingsPatch {
				result[k] = v
				names = append(names, k)
			}

			return result, names
		},

		CheckPreExistence: nil,

		PrepareRequest: func(filteredData map[string]interface{}) interface{} {
			patchRequest := openapi.NewDevicevoicesettingsPutRequest()
			deviceVoiceSettingsMap := make(map[string]openapi.ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName)

			for name, props := range filteredData {
				deviceVoiceSettingsMap[name] = props.(openapi.ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName)
			}
			patchRequest.SetDeviceVoiceSettings(deviceVoiceSettingsMap)
			return patchRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.DeviceVoiceSettingsAPI.DevicevoicesettingsPatch(ctx).DevicevoicesettingsPutRequest(
				*request.(*openapi.DevicevoicesettingsPutRequest))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentDeviceVoiceSettingsOps = true
			b.recentDeviceVoiceSettingsOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkDeviceVoiceSettingsDelete(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "device_voice_settings",
		OperationType: "DELETE",

		ExtractOperations: func() (map[string]interface{}, []string) {
			b.mutex.Lock()
			deviceVoiceSettingsNames := make([]string, len(b.deviceVoiceSettingsDelete))
			copy(deviceVoiceSettingsNames, b.deviceVoiceSettingsDelete)

			deviceVoiceSettingsDeleteMap := make(map[string]bool)
			for _, name := range deviceVoiceSettingsNames {
				deviceVoiceSettingsDeleteMap[name] = true
			}
			b.deviceVoiceSettingsDelete = make([]string, 0)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			for _, name := range deviceVoiceSettingsNames {
				result[name] = true
			}

			return result, deviceVoiceSettingsNames
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
			req := b.client.DeviceVoiceSettingsAPI.DevicevoicesettingsDelete(ctx).DeviceVoiceSettingsName(request.([]string))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentDeviceVoiceSettingsOps = true
			b.recentDeviceVoiceSettingsOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkPacketBrokerPut(ctx context.Context) diag.Diagnostics {
	var originalOperations map[string]openapi.ConfigPutRequestPbEgressProfilePbEgressProfileName

	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "packet_broker",
		OperationType: "PUT",

		ExtractOperations: func() (map[string]interface{}, []string) {
			b.mutex.Lock()
			originalOperations = make(map[string]openapi.ConfigPutRequestPbEgressProfilePbEgressProfileName)
			for k, v := range b.packetBrokerPut {
				originalOperations[k] = v
			}
			b.packetBrokerPut = make(map[string]openapi.ConfigPutRequestPbEgressProfilePbEgressProfileName)
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
				ResourceType:  "packet_broker",
				OperationType: "PUT",
				FetchResources: func(ctx context.Context) (map[string]interface{}, error) {
					apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
					defer cancel()

					resp, err := b.client.PacketBrokerAPI.PacketbrokerGet(apiCtx).IncludeData(true).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						PbEgressProfile map[string]interface{} `json:"pb_egress_profile"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.PbEgressProfile, nil
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
			putRequest := openapi.NewPacketbrokerPutRequest()
			packetBrokerMap := make(map[string]openapi.ConfigPutRequestPbEgressProfilePbEgressProfileName)

			for name, props := range filteredData {
				pbConfig := props.(openapi.ConfigPutRequestPbEgressProfilePbEgressProfileName)
				packetBrokerMap[name] = pbConfig
			}
			putRequest.SetPbEgressProfile(packetBrokerMap)
			return putRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.PacketBrokerAPI.PacketbrokerPut(ctx).PacketbrokerPutRequest(
				*request.(*openapi.PacketbrokerPutRequest))
			return req.Execute()
		},

		ProcessResponse: func(ctx context.Context, resp *http.Response) error {
			tflog.Debug(ctx, "Processing PacketBroker PUT response")

			// Fetch updated data for all resources
			getResp, err := b.client.PacketBrokerAPI.PacketbrokerGet(ctx).IncludeData(true).Execute()
			if err != nil {
				tflog.Error(ctx, "Failed to fetch PacketBroker data after PUT", map[string]interface{}{
					"error": err.Error(),
				})
				return fmt.Errorf("failed to fetch PacketBroker data after PUT: %v", err)
			}
			defer getResp.Body.Close()

			var result struct {
				PbEgressProfile map[string]map[string]interface{} `json:"pb_egress_profile"`
			}
			if err := json.NewDecoder(getResp.Body).Decode(&result); err != nil {
				tflog.Error(ctx, "Failed to decode PacketBroker response after PUT", map[string]interface{}{
					"error": err.Error(),
				})
				return fmt.Errorf("failed to decode PacketBroker response after PUT: %v", err)
			}

			// Update cache with fresh data
			b.packetBrokerResponsesMutex.Lock()
			for name, data := range result.PbEgressProfile {
				b.packetBrokerResponses[name] = data
			}
			b.packetBrokerResponsesMutex.Unlock()

			tflog.Debug(ctx, "Updated PacketBroker cache after PUT", map[string]interface{}{
				"cached_resources": len(result.PbEgressProfile),
			})

			return nil
		},

		UpdateRecentOps: func() {
			b.recentPacketBrokerOps = true
			b.recentPacketBrokerOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkPacketBrokerPatch(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "packet_broker",
		OperationType: "PATCH",

		ExtractOperations: func() (map[string]interface{}, []string) {
			b.mutex.Lock()
			packetBrokerPatch := make(map[string]openapi.ConfigPutRequestPbEgressProfilePbEgressProfileName)
			for k, v := range b.packetBrokerPatch {
				packetBrokerPatch[k] = v
			}
			b.packetBrokerPatch = make(map[string]openapi.ConfigPutRequestPbEgressProfilePbEgressProfileName)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			names := make([]string, 0, len(packetBrokerPatch))

			for k, v := range packetBrokerPatch {
				result[k] = v
				names = append(names, k)
			}

			return result, names
		},

		CheckPreExistence: nil,

		PrepareRequest: func(filteredData map[string]interface{}) interface{} {
			patchRequest := openapi.NewPacketbrokerPutRequest()
			packetBrokerMap := make(map[string]openapi.ConfigPutRequestPbEgressProfilePbEgressProfileName)

			for name, props := range filteredData {
				pbConfig := props.(openapi.ConfigPutRequestPbEgressProfilePbEgressProfileName)
				packetBrokerMap[name] = pbConfig
			}
			patchRequest.SetPbEgressProfile(packetBrokerMap)
			return patchRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.PacketBrokerAPI.PacketbrokerPatch(ctx).PacketbrokerPutRequest(
				*request.(*openapi.PacketbrokerPutRequest))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentPacketBrokerOps = true
			b.recentPacketBrokerOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkPacketBrokerDelete(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "packet_broker",
		OperationType: "DELETE",

		ExtractOperations: func() (map[string]interface{}, []string) {
			b.mutex.Lock()
			packetBrokerNames := make([]string, len(b.packetBrokerDelete))
			copy(packetBrokerNames, b.packetBrokerDelete)

			packetBrokerDeleteMap := make(map[string]bool)
			for _, name := range packetBrokerNames {
				packetBrokerDeleteMap[name] = true
			}
			b.packetBrokerDelete = make([]string, 0)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			for _, name := range packetBrokerNames {
				result[name] = true
			}

			return result, packetBrokerNames
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
			req := b.client.PacketBrokerAPI.PacketbrokerDelete(ctx).PbEgressProfileName(request.([]string))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentPacketBrokerOps = true
			b.recentPacketBrokerOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkPacketQueuePut(ctx context.Context) diag.Diagnostics {
	var originalOperations map[string]openapi.ConfigPutRequestPacketQueuePacketQueueName

	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "packet_queue",
		OperationType: "PUT",

		ExtractOperations: func() (map[string]interface{}, []string) {
			b.mutex.Lock()
			originalOperations = make(map[string]openapi.ConfigPutRequestPacketQueuePacketQueueName)
			for k, v := range b.packetQueuePut {
				originalOperations[k] = v
			}
			b.packetQueuePut = make(map[string]openapi.ConfigPutRequestPacketQueuePacketQueueName)
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
				ResourceType:  "packet_queue",
				OperationType: "PUT",
				FetchResources: func(ctx context.Context) (map[string]interface{}, error) {
					// First check if we have cached packet queue data
					b.packetQueueResponsesMutex.RLock()
					if len(b.packetQueueResponses) > 0 {
						cachedData := make(map[string]interface{})
						for k, v := range b.packetQueueResponses {
							cachedData[k] = v
						}
						b.packetQueueResponsesMutex.RUnlock()
						return cachedData, nil
					}
					b.packetQueueResponsesMutex.RUnlock()

					// Fall back to API call if no cache
					apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
					defer cancel()

					resp, err := b.client.PacketQueuesAPI.PacketqueuesGet(apiCtx).IncludeData(true).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						PacketQueue map[string]interface{} `json:"packet_queue"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}

					// Cache the results
					b.packetQueueResponsesMutex.Lock()
					for k, v := range result.PacketQueue {
						if vMap, ok := v.(map[string]interface{}); ok {
							b.packetQueueResponses[k] = vMap
						}
					}
					b.packetQueueResponsesMutex.Unlock()

					return result.PacketQueue, nil
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
			putRequest := openapi.NewPacketqueuesPutRequest()
			packetQueueMap := make(map[string]openapi.ConfigPutRequestPacketQueuePacketQueueName)

			for name, props := range filteredData {
				packetQueueMap[name] = props.(openapi.ConfigPutRequestPacketQueuePacketQueueName)
			}
			putRequest.SetPacketQueue(packetQueueMap)
			return putRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.PacketQueuesAPI.PacketqueuesPut(ctx).PacketqueuesPutRequest(
				*request.(*openapi.PacketqueuesPutRequest))
			return req.Execute()
		},

		ProcessResponse: func(ctx context.Context, resp *http.Response) error {
			delayTime := 2 * time.Second
			tflog.Debug(ctx, fmt.Sprintf("Waiting %v for auto-generated values to be assigned before fetching packet queues", delayTime))
			time.Sleep(delayTime)

			fetchCtx, fetchCancel := context.WithTimeout(context.Background(), OperationTimeout)
			defer fetchCancel()

			tflog.Debug(ctx, "Fetching packet queues after successful PUT operation to retrieve auto-generated values")
			queuesReq := b.client.PacketQueuesAPI.PacketqueuesGet(fetchCtx)
			queuesResp, fetchErr := queuesReq.Execute()

			if fetchErr != nil {
				tflog.Error(ctx, "Failed to fetch packet queues after PUT for auto-generated fields", map[string]interface{}{
					"error": fetchErr.Error(),
				})
				return fetchErr
			}

			defer queuesResp.Body.Close()

			var queuesData struct {
				PacketQueue map[string]map[string]interface{} `json:"packet_queue"`
			}

			if respErr := json.NewDecoder(queuesResp.Body).Decode(&queuesData); respErr != nil {
				tflog.Error(ctx, "Failed to decode packet queues response after PUT", map[string]interface{}{
					"error": respErr.Error(),
				})
				return respErr
			}

			b.packetQueueResponsesMutex.Lock()
			for queueName, queueData := range queuesData.PacketQueue {
				b.packetQueueResponses[queueName] = queueData
			}
			b.packetQueueResponsesMutex.Unlock()

			tflog.Debug(ctx, "Successfully stored packet queue data for auto-generated fields", map[string]interface{}{
				"packet_queue_count": len(queuesData.PacketQueue),
			})

			return nil
		},

		UpdateRecentOps: func() {
			b.recentPacketQueueOps = true
			b.recentPacketQueueOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkPacketQueuePatch(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "packet_queue",
		OperationType: "PATCH",

		ExtractOperations: func() (map[string]interface{}, []string) {
			b.mutex.Lock()
			packetQueuePatch := make(map[string]openapi.ConfigPutRequestPacketQueuePacketQueueName)
			for k, v := range b.packetQueuePatch {
				packetQueuePatch[k] = v
			}
			b.packetQueuePatch = make(map[string]openapi.ConfigPutRequestPacketQueuePacketQueueName)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			names := make([]string, 0, len(packetQueuePatch))

			for k, v := range packetQueuePatch {
				result[k] = v
				names = append(names, k)
			}

			return result, names
		},

		CheckPreExistence: nil,

		PrepareRequest: func(filteredData map[string]interface{}) interface{} {
			patchRequest := openapi.NewPacketqueuesPutRequest()
			packetQueueMap := make(map[string]openapi.ConfigPutRequestPacketQueuePacketQueueName)

			for name, props := range filteredData {
				packetQueueMap[name] = props.(openapi.ConfigPutRequestPacketQueuePacketQueueName)
			}
			patchRequest.SetPacketQueue(packetQueueMap)
			return patchRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.PacketQueuesAPI.PacketqueuesPatch(ctx).PacketqueuesPutRequest(
				*request.(*openapi.PacketqueuesPutRequest))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentPacketQueueOps = true
			b.recentPacketQueueOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkPacketQueueDelete(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "packet_queue",
		OperationType: "DELETE",

		ExtractOperations: func() (map[string]interface{}, []string) {
			b.mutex.Lock()
			packetQueueNames := make([]string, len(b.packetQueueDelete))
			copy(packetQueueNames, b.packetQueueDelete)

			packetQueueDeleteMap := make(map[string]bool)
			for _, name := range packetQueueNames {
				packetQueueDeleteMap[name] = true
			}
			b.packetQueueDelete = make([]string, 0)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			for _, name := range packetQueueNames {
				result[name] = true
			}

			return result, packetQueueNames
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
			req := b.client.PacketQueuesAPI.PacketqueuesDelete(ctx).PacketQueueName(request.([]string))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentPacketQueueOps = true
			b.recentPacketQueueOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkServicePortProfilePut(ctx context.Context) diag.Diagnostics {
	var originalOperations map[string]openapi.ConfigPutRequestServicePortProfileServicePortProfileName

	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "service_port_profile",
		OperationType: "PUT",

		ExtractOperations: func() (map[string]interface{}, []string) {
			b.mutex.Lock()
			originalOperations = make(map[string]openapi.ConfigPutRequestServicePortProfileServicePortProfileName)
			for k, v := range b.servicePortProfilePut {
				originalOperations[k] = v
			}
			b.servicePortProfilePut = make(map[string]openapi.ConfigPutRequestServicePortProfileServicePortProfileName)
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
				ResourceType:  "service_port_profile",
				OperationType: "PUT",
				FetchResources: func(ctx context.Context) (map[string]interface{}, error) {
					apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
					defer cancel()

					resp, err := b.client.ServicePortProfilesAPI.ServiceportprofilesGet(apiCtx).IncludeData(true).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						ServicePortProfile map[string]interface{} `json:"service_port_profile"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.ServicePortProfile, nil
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
			putRequest := openapi.NewServiceportprofilesPutRequest()
			servicePortProfileMap := make(map[string]openapi.ConfigPutRequestServicePortProfileServicePortProfileName)

			for name, props := range filteredData {
				servicePortProfileMap[name] = props.(openapi.ConfigPutRequestServicePortProfileServicePortProfileName)
			}
			putRequest.SetServicePortProfile(servicePortProfileMap)
			return putRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.ServicePortProfilesAPI.ServiceportprofilesPut(ctx).ServiceportprofilesPutRequest(
				*request.(*openapi.ServiceportprofilesPutRequest))
			return req.Execute()
		},

		ProcessResponse: func(ctx context.Context, resp *http.Response) error {
			delayTime := 2 * time.Second
			tflog.Debug(ctx, fmt.Sprintf("Waiting %v for auto-generated values to be assigned before fetching service port profiles", delayTime))
			time.Sleep(delayTime)

			fetchCtx, fetchCancel := context.WithTimeout(context.Background(), OperationTimeout)
			defer fetchCancel()

			tflog.Debug(ctx, "Fetching service port profiles after successful PUT operation to retrieve auto-generated values")
			profilesReq := b.client.ServicePortProfilesAPI.ServiceportprofilesGet(fetchCtx).IncludeData(true)
			profilesResp, fetchErr := profilesReq.Execute()

			if fetchErr != nil {
				tflog.Error(ctx, "Failed to fetch service port profiles after PUT for auto-generated fields", map[string]interface{}{
					"error": fetchErr.Error(),
				})
				return fetchErr
			}

			defer profilesResp.Body.Close()

			var profilesData struct {
				ServicePortProfile map[string]map[string]interface{} `json:"service_port_profile"`
			}

			if respErr := json.NewDecoder(profilesResp.Body).Decode(&profilesData); respErr != nil {
				tflog.Error(ctx, "Failed to decode service port profiles response after PUT", map[string]interface{}{
					"error": respErr.Error(),
				})
				return respErr
			}

			b.servicePortProfileResponsesMutex.Lock()
			for profileName, profileData := range profilesData.ServicePortProfile {
				b.servicePortProfileResponses[profileName] = profileData
			}
			b.servicePortProfileResponsesMutex.Unlock()

			tflog.Debug(ctx, "Successfully stored service port profile data for auto-generated fields", map[string]interface{}{
				"service_port_profile_count": len(profilesData.ServicePortProfile),
			})

			return nil
		},

		UpdateRecentOps: func() {
			b.recentServicePortProfileOps = true
			b.recentServicePortProfileOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkServicePortProfilePatch(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "service_port_profile",
		OperationType: "PATCH",

		ExtractOperations: func() (map[string]interface{}, []string) {
			b.mutex.Lock()
			servicePortProfilePatch := make(map[string]openapi.ConfigPutRequestServicePortProfileServicePortProfileName)
			for k, v := range b.servicePortProfilePatch {
				servicePortProfilePatch[k] = v
			}
			b.servicePortProfilePatch = make(map[string]openapi.ConfigPutRequestServicePortProfileServicePortProfileName)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			names := make([]string, 0, len(servicePortProfilePatch))

			for k, v := range servicePortProfilePatch {
				result[k] = v
				names = append(names, k)
			}

			return result, names
		},

		CheckPreExistence: nil,

		PrepareRequest: func(filteredData map[string]interface{}) interface{} {
			patchRequest := openapi.NewServiceportprofilesPutRequest()
			servicePortProfileMap := make(map[string]openapi.ConfigPutRequestServicePortProfileServicePortProfileName)

			for name, props := range filteredData {
				servicePortProfileMap[name] = props.(openapi.ConfigPutRequestServicePortProfileServicePortProfileName)
			}
			patchRequest.SetServicePortProfile(servicePortProfileMap)
			return patchRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.ServicePortProfilesAPI.ServiceportprofilesPatch(ctx).ServiceportprofilesPutRequest(
				*request.(*openapi.ServiceportprofilesPutRequest))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentServicePortProfileOps = true
			b.recentServicePortProfileOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkServicePortProfileDelete(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "service_port_profile",
		OperationType: "DELETE",

		ExtractOperations: func() (map[string]interface{}, []string) {
			b.mutex.Lock()
			servicePortProfileNames := make([]string, len(b.servicePortProfileDelete))
			copy(servicePortProfileNames, b.servicePortProfileDelete)

			servicePortProfileDeleteMap := make(map[string]bool)
			for _, name := range servicePortProfileNames {
				servicePortProfileDeleteMap[name] = true
			}
			b.servicePortProfileDelete = make([]string, 0)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			for _, name := range servicePortProfileNames {
				result[name] = true
			}

			return result, servicePortProfileNames
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
			req := b.client.ServicePortProfilesAPI.ServiceportprofilesDelete(ctx).ServicePortProfileName(request.([]string))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentServicePortProfileOps = true
			b.recentServicePortProfileOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkSwitchpointPut(ctx context.Context) diag.Diagnostics {
	var originalOperations map[string]openapi.ConfigPutRequestSwitchpointSwitchpointName

	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "switchpoint",
		OperationType: "PUT",

		ExtractOperations: func() (map[string]interface{}, []string) {
			b.mutex.Lock()
			originalOperations = make(map[string]openapi.ConfigPutRequestSwitchpointSwitchpointName)
			for k, v := range b.switchpointPut {
				originalOperations[k] = v
			}
			b.switchpointPut = make(map[string]openapi.ConfigPutRequestSwitchpointSwitchpointName)
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
				ResourceType:  "switchpoint",
				OperationType: "PUT",
				FetchResources: func(ctx context.Context) (map[string]interface{}, error) {
					// First check if we have cached Switchpoint data
					b.switchpointResponsesMutex.RLock()
					if len(b.switchpointResponses) > 0 {
						cachedData := make(map[string]interface{})
						for k, v := range b.switchpointResponses {
							cachedData[k] = v
						}
						b.switchpointResponsesMutex.RUnlock()

						tflog.Debug(ctx, "Using cached Switchpoint data for pre-existence check", map[string]interface{}{
							"count": len(cachedData),
						})

						return cachedData, nil
					}
					b.switchpointResponsesMutex.RUnlock()

					// Fall back to API call if no cache
					apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
					defer cancel()

					resp, err := b.client.SwitchpointsAPI.SwitchpointsGet(apiCtx).IncludeData(true).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						Switchpoint map[string]interface{} `json:"switchpoint"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}

					// Update cache with API data
					b.switchpointResponsesMutex.Lock()
					for k, v := range result.Switchpoint {
						if vMap, ok := v.(map[string]interface{}); ok {
							b.switchpointResponses[k] = vMap

							if name, ok := vMap["name"].(string); ok && name != k {
								b.switchpointResponses[name] = vMap
							}
						}
					}
					b.switchpointResponsesMutex.Unlock()

					return result.Switchpoint, nil
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
			putRequest := openapi.NewSwitchpointsPutRequest()
			switchpointMap := make(map[string]openapi.SwitchpointsPutRequestSwitchpointValue)

			for name, props := range filteredData {
				configProps := props.(openapi.ConfigPutRequestSwitchpointSwitchpointName)
				switchpointValue := openapi.SwitchpointsPutRequestSwitchpointValue{}

				if configProps.Name != nil {
					switchpointValue.SetName(*configProps.Name)
				}
				if configProps.DeviceSerialNumber != nil {
					switchpointValue.SetDeviceSerialNumber(*configProps.DeviceSerialNumber)
				}
				if configProps.ConnectedBundle != nil {
					switchpointValue.SetConnectedBundle(*configProps.ConnectedBundle)
				}
				if configProps.ReadOnlyMode != nil {
					switchpointValue.SetReadOnlyMode(*configProps.ReadOnlyMode)
				}
				if configProps.Locked != nil {
					switchpointValue.SetLocked(*configProps.Locked)
				}
				if configProps.DisabledPorts != nil {
					switchpointValue.SetDisabledPorts(*configProps.DisabledPorts)
				}
				if configProps.OutOfBandManagement != nil {
					switchpointValue.SetOutOfBandManagement(*configProps.OutOfBandManagement)
				}
				if configProps.Type != nil {
					switchpointValue.SetType(*configProps.Type)
				}
				if configProps.SuperPod != nil {
					switchpointValue.SetSuperPod(*configProps.SuperPod)
				}
				if configProps.Pod != nil {
					switchpointValue.SetPod(*configProps.Pod)
				}
				if configProps.Rack != nil {
					switchpointValue.SetRack(*configProps.Rack)
				}
				if configProps.SwitchRouterIdIpMask != nil {
					switchpointValue.SetSwitchRouterIdIpMask(*configProps.SwitchRouterIdIpMask)
				}

				switchpointMap[name] = switchpointValue
			}
			putRequest.SetSwitchpoint(switchpointMap)
			return putRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.SwitchpointsAPI.SwitchpointsPut(ctx).SwitchpointsPutRequest(
				*request.(*openapi.SwitchpointsPutRequest))
			return req.Execute()
		},

		ProcessResponse: func(ctx context.Context, resp *http.Response) error {
			delayTime := 2 * time.Second
			tflog.Debug(ctx, fmt.Sprintf("Waiting %v for auto-generated values to be assigned before fetching switchpoints", delayTime))
			time.Sleep(delayTime)

			fetchCtx, fetchCancel := context.WithTimeout(context.Background(), OperationTimeout)
			defer fetchCancel()

			tflog.Debug(ctx, "Fetching switchpoints after successful PUT operation to retrieve auto-generated values")
			switchpointsReq := b.client.SwitchpointsAPI.SwitchpointsGet(fetchCtx).IncludeData(true)
			switchpointsResp, fetchErr := switchpointsReq.Execute()

			if fetchErr != nil {
				tflog.Error(ctx, "Failed to fetch switchpoints after PUT for auto-generated fields", map[string]interface{}{
					"error": fetchErr.Error(),
				})
				return fetchErr
			}

			defer switchpointsResp.Body.Close()

			var switchpointsData struct {
				Switchpoint map[string]map[string]interface{} `json:"switchpoint"`
			}

			if respErr := json.NewDecoder(switchpointsResp.Body).Decode(&switchpointsData); respErr != nil {
				tflog.Error(ctx, "Failed to decode switchpoints response after PUT", map[string]interface{}{
					"error": respErr.Error(),
				})
				return respErr
			}

			b.switchpointResponsesMutex.Lock()
			for switchpointName, switchpointData := range switchpointsData.Switchpoint {
				b.switchpointResponses[switchpointName] = switchpointData
			}
			b.switchpointResponsesMutex.Unlock()

			tflog.Debug(ctx, "Successfully stored switchpoint data for auto-generated fields", map[string]interface{}{
				"switchpoint_count": len(switchpointsData.Switchpoint),
			})

			return nil
		},

		UpdateRecentOps: func() {
			b.recentSwitchpointOps = true
			b.recentSwitchpointOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkSwitchpointPatch(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "switchpoint",
		OperationType: "PATCH",

		ExtractOperations: func() (map[string]interface{}, []string) {
			b.mutex.Lock()
			switchpointPatch := make(map[string]openapi.ConfigPutRequestSwitchpointSwitchpointName)
			for k, v := range b.switchpointPatch {
				switchpointPatch[k] = v
			}
			b.switchpointPatch = make(map[string]openapi.ConfigPutRequestSwitchpointSwitchpointName)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			names := make([]string, 0, len(switchpointPatch))

			for k, v := range switchpointPatch {
				result[k] = v
				names = append(names, k)
			}

			return result, names
		},

		CheckPreExistence: nil,

		PrepareRequest: func(filteredData map[string]interface{}) interface{} {
			patchRequest := openapi.NewSwitchpointsPutRequest()
			switchpointMap := make(map[string]openapi.SwitchpointsPutRequestSwitchpointValue)

			for name, props := range filteredData {
				configProps := props.(openapi.ConfigPutRequestSwitchpointSwitchpointName)
				switchpointValue := openapi.SwitchpointsPutRequestSwitchpointValue{}

				if configProps.Name != nil {
					switchpointValue.SetName(*configProps.Name)
				}
				if configProps.DeviceSerialNumber != nil {
					switchpointValue.SetDeviceSerialNumber(*configProps.DeviceSerialNumber)
				}
				if configProps.ConnectedBundle != nil {
					switchpointValue.SetConnectedBundle(*configProps.ConnectedBundle)
				}
				if configProps.ReadOnlyMode != nil {
					switchpointValue.SetReadOnlyMode(*configProps.ReadOnlyMode)
				}
				if configProps.Locked != nil {
					switchpointValue.SetLocked(*configProps.Locked)
				}
				if configProps.DisabledPorts != nil {
					switchpointValue.SetDisabledPorts(*configProps.DisabledPorts)
				}
				if configProps.OutOfBandManagement != nil {
					switchpointValue.SetOutOfBandManagement(*configProps.OutOfBandManagement)
				}
				if configProps.Type != nil {
					switchpointValue.SetType(*configProps.Type)
				}
				if configProps.SuperPod != nil {
					switchpointValue.SetSuperPod(*configProps.SuperPod)
				}
				if configProps.Pod != nil {
					switchpointValue.SetPod(*configProps.Pod)
				}
				if configProps.Rack != nil {
					switchpointValue.SetRack(*configProps.Rack)
				}
				if configProps.SwitchRouterIdIpMask != nil {
					switchpointValue.SetSwitchRouterIdIpMask(*configProps.SwitchRouterIdIpMask)
				}

				switchpointMap[name] = switchpointValue
			}
			patchRequest.SetSwitchpoint(switchpointMap)
			return patchRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.SwitchpointsAPI.SwitchpointsPatch(ctx).SwitchpointsPutRequest(
				*request.(*openapi.SwitchpointsPutRequest))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentSwitchpointOps = true
			b.recentSwitchpointOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkSwitchpointDelete(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "switchpoint",
		OperationType: "DELETE",

		ExtractOperations: func() (map[string]interface{}, []string) {
			b.mutex.Lock()
			switchpointNames := make([]string, len(b.switchpointDelete))
			copy(switchpointNames, b.switchpointDelete)

			switchpointDeleteMap := make(map[string]bool)
			for _, name := range switchpointNames {
				switchpointDeleteMap[name] = true
			}
			b.switchpointDelete = make([]string, 0)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			for _, name := range switchpointNames {
				result[name] = true
			}

			return result, switchpointNames
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
			req := b.client.SwitchpointsAPI.SwitchpointsDelete(ctx).SwitchpointName(request.([]string))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentSwitchpointOps = true
			b.recentSwitchpointOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkVoicePortProfilePut(ctx context.Context) diag.Diagnostics {
	var originalOperations map[string]openapi.ConfigPutRequestVoicePortProfilesVoicePortProfilesName

	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "voice_port_profile",
		OperationType: "PUT",

		ExtractOperations: func() (map[string]interface{}, []string) {
			b.mutex.Lock()
			originalOperations = make(map[string]openapi.ConfigPutRequestVoicePortProfilesVoicePortProfilesName)
			for k, v := range b.voicePortProfilePut {
				originalOperations[k] = v
			}
			b.voicePortProfilePut = make(map[string]openapi.ConfigPutRequestVoicePortProfilesVoicePortProfilesName)
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
				ResourceType:  "voice_port_profile",
				OperationType: "PUT",
				FetchResources: func(ctx context.Context) (map[string]interface{}, error) {
					// First check if we have cached Voice Port Profile data
					b.voicePortProfileResponsesMutex.RLock()
					if len(b.voicePortProfileResponses) > 0 {
						cachedData := make(map[string]interface{})
						for k, v := range b.voicePortProfileResponses {
							cachedData[k] = v
						}
						b.voicePortProfileResponsesMutex.RUnlock()

						tflog.Debug(ctx, "Using cached Voice Port Profile data for pre-existence check", map[string]interface{}{
							"count": len(cachedData),
						})

						return cachedData, nil
					}
					b.voicePortProfileResponsesMutex.RUnlock()

					// Fall back to API call if no cache
					apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
					defer cancel()

					resp, err := b.client.VoicePortProfilesAPI.VoiceportprofilesGet(apiCtx).IncludeData(true).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						VoicePortProfiles map[string]interface{} `json:"voice_port_profiles"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}

					// Update cache with API data
					b.voicePortProfileResponsesMutex.Lock()
					for k, v := range result.VoicePortProfiles {
						if vMap, ok := v.(map[string]interface{}); ok {
							b.voicePortProfileResponses[k] = vMap

							if name, ok := vMap["name"].(string); ok && name != k {
								b.voicePortProfileResponses[name] = vMap
							}
						}
					}
					b.voicePortProfileResponsesMutex.Unlock()

					return result.VoicePortProfiles, nil
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
			putRequest := openapi.NewVoiceportprofilesPutRequest()
			profileMap := make(map[string]openapi.ConfigPutRequestVoicePortProfilesVoicePortProfilesName)

			for name, props := range filteredData {
				profileMap[name] = props.(openapi.ConfigPutRequestVoicePortProfilesVoicePortProfilesName)
			}
			putRequest.SetVoicePortProfiles(profileMap)
			return putRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.VoicePortProfilesAPI.VoiceportprofilesPut(ctx).VoiceportprofilesPutRequest(
				*request.(*openapi.VoiceportprofilesPutRequest))
			return req.Execute()
		},

		ProcessResponse: func(ctx context.Context, resp *http.Response) error {
			delayTime := 2 * time.Second
			tflog.Debug(ctx, fmt.Sprintf("Waiting %v for auto-generated values to be assigned before fetching voice port profiles", delayTime))
			time.Sleep(delayTime)

			fetchCtx, fetchCancel := context.WithTimeout(context.Background(), OperationTimeout)
			defer fetchCancel()

			tflog.Debug(ctx, "Fetching voice port profiles after successful PUT operation to retrieve auto-generated values")
			profilesReq := b.client.VoicePortProfilesAPI.VoiceportprofilesGet(fetchCtx).IncludeData(true)
			profilesResp, fetchErr := profilesReq.Execute()

			if fetchErr != nil {
				tflog.Error(ctx, "Failed to fetch voice port profiles after PUT for auto-generated fields", map[string]interface{}{
					"error": fetchErr.Error(),
				})
				return fetchErr
			}

			defer profilesResp.Body.Close()

			var profilesData struct {
				VoicePortProfiles map[string]map[string]interface{} `json:"voice_port_profiles"`
			}

			if respErr := json.NewDecoder(profilesResp.Body).Decode(&profilesData); respErr != nil {
				tflog.Error(ctx, "Failed to decode voice port profiles response after PUT", map[string]interface{}{
					"error": respErr.Error(),
				})
				return respErr
			}

			b.voicePortProfileResponsesMutex.Lock()
			for profileName, profileData := range profilesData.VoicePortProfiles {
				b.voicePortProfileResponses[profileName] = profileData
			}
			b.voicePortProfileResponsesMutex.Unlock()

			tflog.Debug(ctx, "Successfully stored voice port profile data for auto-generated fields", map[string]interface{}{
				"voice_port_profile_count": len(profilesData.VoicePortProfiles),
			})

			return nil
		},

		UpdateRecentOps: func() {
			b.recentVoicePortProfileOps = true
			b.recentVoicePortProfileOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkVoicePortProfilePatch(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "voice_port_profile",
		OperationType: "PATCH",

		ExtractOperations: func() (map[string]interface{}, []string) {
			b.mutex.Lock()
			voicePortProfilePatch := make(map[string]openapi.ConfigPutRequestVoicePortProfilesVoicePortProfilesName)
			for k, v := range b.voicePortProfilePatch {
				voicePortProfilePatch[k] = v
			}
			b.voicePortProfilePatch = make(map[string]openapi.ConfigPutRequestVoicePortProfilesVoicePortProfilesName)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			names := make([]string, 0, len(voicePortProfilePatch))

			for k, v := range voicePortProfilePatch {
				result[k] = v
				names = append(names, k)
			}

			return result, names
		},

		CheckPreExistence: nil,

		PrepareRequest: func(filteredData map[string]interface{}) interface{} {
			patchRequest := openapi.NewVoiceportprofilesPutRequest()
			profileMap := make(map[string]openapi.ConfigPutRequestVoicePortProfilesVoicePortProfilesName)

			for name, props := range filteredData {
				profileMap[name] = props.(openapi.ConfigPutRequestVoicePortProfilesVoicePortProfilesName)
			}
			patchRequest.SetVoicePortProfiles(profileMap)
			return patchRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.VoicePortProfilesAPI.VoiceportprofilesPatch(ctx).VoiceportprofilesPutRequest(
				*request.(*openapi.VoiceportprofilesPutRequest))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentVoicePortProfileOps = true
			b.recentVoicePortProfileOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkVoicePortProfileDelete(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "voice_port_profile",
		OperationType: "DELETE",

		ExtractOperations: func() (map[string]interface{}, []string) {
			b.mutex.Lock()
			voicePortProfileNames := make([]string, len(b.voicePortProfileDelete))
			copy(voicePortProfileNames, b.voicePortProfileDelete)

			voicePortProfileDeleteMap := make(map[string]bool)
			for _, name := range voicePortProfileNames {
				voicePortProfileDeleteMap[name] = true
			}
			b.voicePortProfileDelete = make([]string, 0)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			for _, name := range voicePortProfileNames {
				result[name] = true
			}

			return result, voicePortProfileNames
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
			req := b.client.VoicePortProfilesAPI.VoiceportprofilesDelete(ctx).VoicePortProfileName(request.([]string))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentVoicePortProfileOps = true
			b.recentVoicePortProfileOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkDeviceControllerPut(ctx context.Context) diag.Diagnostics {
	var originalOperations map[string]openapi.ConfigPutRequestDeviceControllerDeviceControllerName

	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "device_controller",
		OperationType: "PUT",

		ExtractOperations: func() (map[string]interface{}, []string) {
			b.mutex.Lock()
			originalOperations = make(map[string]openapi.ConfigPutRequestDeviceControllerDeviceControllerName)
			for k, v := range b.deviceControllerPut {
				originalOperations[k] = v
			}
			b.deviceControllerPut = make(map[string]openapi.ConfigPutRequestDeviceControllerDeviceControllerName)
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
				ResourceType:  "device_controller",
				OperationType: "PUT",
				FetchResources: func(ctx context.Context) (map[string]interface{}, error) {
					// First check if we have cached DeviceController data
					b.deviceControllerResponsesMutex.RLock()
					if len(b.deviceControllerResponses) > 0 {
						cachedData := make(map[string]interface{})
						for k, v := range b.deviceControllerResponses {
							cachedData[k] = v
						}
						b.deviceControllerResponsesMutex.RUnlock()

						tflog.Debug(ctx, "Using cached DeviceController data for pre-existence check", map[string]interface{}{
							"count": len(cachedData),
						})

						return cachedData, nil
					}
					b.deviceControllerResponsesMutex.RUnlock()

					// Fall back to API call if no cache
					apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
					defer cancel()

					resp, err := b.client.DeviceControllersAPI.DevicecontrollersGet(apiCtx).IncludeData(true).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						DeviceController map[string]interface{} `json:"device_controller"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}

					// Update cache with API data
					b.deviceControllerResponsesMutex.Lock()
					for k, v := range result.DeviceController {
						if vMap, ok := v.(map[string]interface{}); ok {
							b.deviceControllerResponses[k] = vMap

							if name, ok := vMap["name"].(string); ok && name != k {
								b.deviceControllerResponses[name] = vMap
							}
						}
					}
					b.deviceControllerResponsesMutex.Unlock()

					return result.DeviceController, nil
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
			patchRequest := openapi.NewDevicecontrollersPutRequest()
			deviceMap := make(map[string]openapi.DevicecontrollersPutRequestDeviceControllerValue)

			for name, props := range filteredData {
				sourceProps := props.(openapi.ConfigPutRequestDeviceControllerDeviceControllerName)
				targetProps := openapi.DevicecontrollersPutRequestDeviceControllerValue{}

				if sourceProps.Name != nil {
					targetProps.Name = sourceProps.Name
				}
				if sourceProps.Enable != nil {
					targetProps.Enable = sourceProps.Enable
				}
				if sourceProps.IpSource != nil {
					targetProps.IpSource = sourceProps.IpSource
				}
				if sourceProps.ControllerIpAndMask != nil {
					targetProps.ControllerIpAndMask = sourceProps.ControllerIpAndMask
				}
				if sourceProps.Gateway != nil {
					targetProps.Gateway = sourceProps.Gateway
				}
				if sourceProps.SwitchIpAndMask != nil {
					targetProps.SwitchIpAndMask = sourceProps.SwitchIpAndMask
				}
				if sourceProps.SwitchGateway != nil {
					targetProps.SwitchGateway = sourceProps.SwitchGateway
				}
				if sourceProps.CommType != nil {
					targetProps.CommType = sourceProps.CommType
				}
				if sourceProps.SnmpCommunityString != nil {
					targetProps.SnmpCommunityString = sourceProps.SnmpCommunityString
				}
				if sourceProps.UplinkPort != nil {
					targetProps.UplinkPort = sourceProps.UplinkPort
				}
				if sourceProps.LldpSearchString != nil {
					targetProps.LldpSearchString = sourceProps.LldpSearchString
				}
				if sourceProps.ZtpIdentification != nil {
					targetProps.ZtpIdentification = sourceProps.ZtpIdentification
				}
				if sourceProps.LocatedBy != nil {
					targetProps.LocatedBy = sourceProps.LocatedBy
				}
				if sourceProps.PowerState != nil {
					targetProps.PowerState = sourceProps.PowerState
				}
				if sourceProps.CommunicationMode != nil {
					targetProps.CommunicationMode = sourceProps.CommunicationMode
				}
				if sourceProps.CliAccessMode != nil {
					targetProps.CliAccessMode = sourceProps.CliAccessMode
				}
				if sourceProps.Username != nil {
					targetProps.Username = sourceProps.Username
				}
				if sourceProps.Password != nil {
					targetProps.Password = sourceProps.Password
				}
				if sourceProps.EnablePassword != nil {
					targetProps.EnablePassword = sourceProps.EnablePassword
				}
				if sourceProps.SshKeyOrPassword != nil {
					targetProps.SshKeyOrPassword = sourceProps.SshKeyOrPassword
				}
				if sourceProps.ManagedOnNativeVlan != nil {
					targetProps.ManagedOnNativeVlan = sourceProps.ManagedOnNativeVlan
				}
				if sourceProps.Sdlc != nil {
					targetProps.Sdlc = sourceProps.Sdlc
				}
				if sourceProps.Switchpoint != nil {
					targetProps.Switchpoint = sourceProps.Switchpoint
				}
				if sourceProps.SwitchpointRefType != nil {
					targetProps.SwitchpointRefType = sourceProps.SwitchpointRefType
				}
				if sourceProps.SecurityType != nil {
					targetProps.SecurityType = sourceProps.SecurityType
				}
				if sourceProps.Snmpv3Username != nil {
					targetProps.Snmpv3Username = sourceProps.Snmpv3Username
				}
				if sourceProps.AuthenticationProtocol != nil {
					targetProps.AuthenticationProtocol = sourceProps.AuthenticationProtocol
				}
				if sourceProps.Passphrase != nil {
					targetProps.Passphrase = sourceProps.Passphrase
				}
				if sourceProps.PrivateProtocol != nil {
					targetProps.PrivateProtocol = sourceProps.PrivateProtocol
				}
				if sourceProps.PrivatePassword != nil {
					targetProps.PrivatePassword = sourceProps.PrivatePassword
				}
				if sourceProps.PasswordEncrypted != nil {
					targetProps.PasswordEncrypted = sourceProps.PasswordEncrypted
				}
				if sourceProps.EnablePasswordEncrypted != nil {
					targetProps.EnablePasswordEncrypted = sourceProps.EnablePasswordEncrypted
				}
				if sourceProps.SshKeyOrPasswordEncrypted != nil {
					targetProps.SshKeyOrPasswordEncrypted = sourceProps.SshKeyOrPasswordEncrypted
				}
				if sourceProps.PassphraseEncrypted != nil {
					targetProps.PassphraseEncrypted = sourceProps.PassphraseEncrypted
				}
				if sourceProps.PrivatePasswordEncrypted != nil {
					targetProps.PrivatePasswordEncrypted = sourceProps.PrivatePasswordEncrypted
				}
				if sourceProps.DeviceManagedAs != nil {
					targetProps.DeviceManagedAs = sourceProps.DeviceManagedAs
				}
				if sourceProps.Switch != nil {
					targetProps.Switch = sourceProps.Switch
				}
				if sourceProps.SwitchRefType != nil {
					targetProps.SwitchRefType = sourceProps.SwitchRefType
				}
				if sourceProps.ConnectionService != nil {
					targetProps.ConnectionService = sourceProps.ConnectionService
				}
				if sourceProps.ConnectionServiceRefType != nil {
					targetProps.ConnectionServiceRefType = sourceProps.ConnectionServiceRefType
				}
				if sourceProps.Port != nil {
					targetProps.Port = sourceProps.Port
				}
				if sourceProps.SfpMacAddressOrSn != nil {
					targetProps.SfpMacAddressOrSn = sourceProps.SfpMacAddressOrSn
				}
				if sourceProps.UsesTaggedPackets != nil {
					targetProps.UsesTaggedPackets = sourceProps.UsesTaggedPackets
				}

				deviceMap[name] = targetProps
			}
			patchRequest.SetDeviceController(deviceMap)
			return patchRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.DeviceControllersAPI.DevicecontrollersPut(ctx).DevicecontrollersPutRequest(
				*request.(*openapi.DevicecontrollersPutRequest))
			return req.Execute()
		},

		ProcessResponse: func(ctx context.Context, resp *http.Response) error {
			delayTime := 2 * time.Second
			tflog.Debug(ctx, fmt.Sprintf("Waiting %v for auto-generated values to be assigned before fetching device controllers", delayTime))
			time.Sleep(delayTime)

			fetchCtx, fetchCancel := context.WithTimeout(context.Background(), OperationTimeout)
			defer fetchCancel()

			tflog.Debug(ctx, "Fetching device controllers after successful PUT operation to retrieve auto-generated values")
			deviceControllersReq := b.client.DeviceControllersAPI.DevicecontrollersGet(fetchCtx).IncludeData(true)
			deviceControllersResp, fetchErr := deviceControllersReq.Execute()

			if fetchErr != nil {
				tflog.Error(ctx, "Failed to fetch device controllers after PUT for auto-generated fields", map[string]interface{}{
					"error": fetchErr.Error(),
				})
				return fetchErr
			}

			defer deviceControllersResp.Body.Close()

			var deviceControllersData struct {
				DeviceController map[string]map[string]interface{} `json:"device_controller"`
			}

			if respErr := json.NewDecoder(deviceControllersResp.Body).Decode(&deviceControllersData); respErr != nil {
				tflog.Error(ctx, "Failed to decode device controllers response after PUT", map[string]interface{}{
					"error": respErr.Error(),
				})
				return respErr
			}

			b.deviceControllerResponsesMutex.Lock()
			for deviceControllerName, deviceControllerData := range deviceControllersData.DeviceController {
				b.deviceControllerResponses[deviceControllerName] = deviceControllerData
			}
			b.deviceControllerResponsesMutex.Unlock()

			tflog.Debug(ctx, "Successfully stored device controller data for auto-generated fields", map[string]interface{}{
				"device_controller_count": len(deviceControllersData.DeviceController),
			})

			return nil
		},

		UpdateRecentOps: func() {
			b.recentDeviceControllerOps = true
			b.recentDeviceControllerOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkDeviceControllerPatch(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "device_controller",
		OperationType: "PATCH",

		ExtractOperations: func() (map[string]interface{}, []string) {
			b.mutex.Lock()
			deviceControllerPatch := make(map[string]openapi.ConfigPutRequestDeviceControllerDeviceControllerName)
			for k, v := range b.deviceControllerPatch {
				deviceControllerPatch[k] = v
			}
			b.deviceControllerPatch = make(map[string]openapi.ConfigPutRequestDeviceControllerDeviceControllerName)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			names := make([]string, 0, len(deviceControllerPatch))

			for k, v := range deviceControllerPatch {
				result[k] = v
				names = append(names, k)
			}

			return result, names
		},

		CheckPreExistence: nil,

		PrepareRequest: func(filteredData map[string]interface{}) interface{} {
			patchRequest := openapi.NewDevicecontrollersPutRequest()
			deviceMap := make(map[string]openapi.DevicecontrollersPutRequestDeviceControllerValue)

			for name, props := range filteredData {
				sourceProps := props.(openapi.ConfigPutRequestDeviceControllerDeviceControllerName)
				targetProps := openapi.DevicecontrollersPutRequestDeviceControllerValue{}

				if sourceProps.Name != nil {
					targetProps.Name = sourceProps.Name
				}
				if sourceProps.Enable != nil {
					targetProps.Enable = sourceProps.Enable
				}
				if sourceProps.IpSource != nil {
					targetProps.IpSource = sourceProps.IpSource
				}
				if sourceProps.ControllerIpAndMask != nil {
					targetProps.ControllerIpAndMask = sourceProps.ControllerIpAndMask
				}
				if sourceProps.Gateway != nil {
					targetProps.Gateway = sourceProps.Gateway
				}
				if sourceProps.SwitchIpAndMask != nil {
					targetProps.SwitchIpAndMask = sourceProps.SwitchIpAndMask
				}
				if sourceProps.SwitchGateway != nil {
					targetProps.SwitchGateway = sourceProps.SwitchGateway
				}
				if sourceProps.CommType != nil {
					targetProps.CommType = sourceProps.CommType
				}
				if sourceProps.SnmpCommunityString != nil {
					targetProps.SnmpCommunityString = sourceProps.SnmpCommunityString
				}
				if sourceProps.UplinkPort != nil {
					targetProps.UplinkPort = sourceProps.UplinkPort
				}
				if sourceProps.LldpSearchString != nil {
					targetProps.LldpSearchString = sourceProps.LldpSearchString
				}
				if sourceProps.ZtpIdentification != nil {
					targetProps.ZtpIdentification = sourceProps.ZtpIdentification
				}
				if sourceProps.LocatedBy != nil {
					targetProps.LocatedBy = sourceProps.LocatedBy
				}
				if sourceProps.PowerState != nil {
					targetProps.PowerState = sourceProps.PowerState
				}
				if sourceProps.CommunicationMode != nil {
					targetProps.CommunicationMode = sourceProps.CommunicationMode
				}
				if sourceProps.CliAccessMode != nil {
					targetProps.CliAccessMode = sourceProps.CliAccessMode
				}
				if sourceProps.Username != nil {
					targetProps.Username = sourceProps.Username
				}
				if sourceProps.Password != nil {
					targetProps.Password = sourceProps.Password
				}
				if sourceProps.EnablePassword != nil {
					targetProps.EnablePassword = sourceProps.EnablePassword
				}
				if sourceProps.SshKeyOrPassword != nil {
					targetProps.SshKeyOrPassword = sourceProps.SshKeyOrPassword
				}
				if sourceProps.ManagedOnNativeVlan != nil {
					targetProps.ManagedOnNativeVlan = sourceProps.ManagedOnNativeVlan
				}
				if sourceProps.Sdlc != nil {
					targetProps.Sdlc = sourceProps.Sdlc
				}
				if sourceProps.Switchpoint != nil {
					targetProps.Switchpoint = sourceProps.Switchpoint
				}
				if sourceProps.SwitchpointRefType != nil {
					targetProps.SwitchpointRefType = sourceProps.SwitchpointRefType
				}
				if sourceProps.SecurityType != nil {
					targetProps.SecurityType = sourceProps.SecurityType
				}
				if sourceProps.Snmpv3Username != nil {
					targetProps.Snmpv3Username = sourceProps.Snmpv3Username
				}
				if sourceProps.AuthenticationProtocol != nil {
					targetProps.AuthenticationProtocol = sourceProps.AuthenticationProtocol
				}
				if sourceProps.Passphrase != nil {
					targetProps.Passphrase = sourceProps.Passphrase
				}
				if sourceProps.PrivateProtocol != nil {
					targetProps.PrivateProtocol = sourceProps.PrivateProtocol
				}
				if sourceProps.PrivatePassword != nil {
					targetProps.PrivatePassword = sourceProps.PrivatePassword
				}
				if sourceProps.PasswordEncrypted != nil {
					targetProps.PasswordEncrypted = sourceProps.PasswordEncrypted
				}
				if sourceProps.EnablePasswordEncrypted != nil {
					targetProps.EnablePasswordEncrypted = sourceProps.EnablePasswordEncrypted
				}
				if sourceProps.SshKeyOrPasswordEncrypted != nil {
					targetProps.SshKeyOrPasswordEncrypted = sourceProps.SshKeyOrPasswordEncrypted
				}
				if sourceProps.PassphraseEncrypted != nil {
					targetProps.PassphraseEncrypted = sourceProps.PassphraseEncrypted
				}
				if sourceProps.PrivatePasswordEncrypted != nil {
					targetProps.PrivatePasswordEncrypted = sourceProps.PrivatePasswordEncrypted
				}
				if sourceProps.DeviceManagedAs != nil {
					targetProps.DeviceManagedAs = sourceProps.DeviceManagedAs
				}
				if sourceProps.Switch != nil {
					targetProps.Switch = sourceProps.Switch
				}
				if sourceProps.SwitchRefType != nil {
					targetProps.SwitchRefType = sourceProps.SwitchRefType
				}
				if sourceProps.ConnectionService != nil {
					targetProps.ConnectionService = sourceProps.ConnectionService
				}
				if sourceProps.ConnectionServiceRefType != nil {
					targetProps.ConnectionServiceRefType = sourceProps.ConnectionServiceRefType
				}
				if sourceProps.Port != nil {
					targetProps.Port = sourceProps.Port
				}
				if sourceProps.SfpMacAddressOrSn != nil {
					targetProps.SfpMacAddressOrSn = sourceProps.SfpMacAddressOrSn
				}
				if sourceProps.UsesTaggedPackets != nil {
					targetProps.UsesTaggedPackets = sourceProps.UsesTaggedPackets
				}

				deviceMap[name] = targetProps
			}
			patchRequest.SetDeviceController(deviceMap)
			return patchRequest
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			req := b.client.DeviceControllersAPI.DevicecontrollersPatch(ctx).DevicecontrollersPutRequest(
				*request.(*openapi.DevicecontrollersPutRequest))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentDeviceControllerOps = true
			b.recentDeviceControllerOpTime = time.Now()
		},
	})
}

func (b *BulkOperationManager) ExecuteBulkDeviceControllerDelete(ctx context.Context) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  "device_controller",
		OperationType: "DELETE",

		ExtractOperations: func() (map[string]interface{}, []string) {
			b.mutex.Lock()
			deviceControllerNames := make([]string, len(b.deviceControllerDelete))
			copy(deviceControllerNames, b.deviceControllerDelete)

			deviceControllerDeleteMap := make(map[string]bool)
			for _, name := range deviceControllerNames {
				deviceControllerDeleteMap[name] = true
			}
			b.deviceControllerDelete = make([]string, 0)
			b.mutex.Unlock()

			result := make(map[string]interface{})
			for _, name := range deviceControllerNames {
				result[name] = true
			}

			return result, deviceControllerNames
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
			req := b.client.DeviceControllersAPI.DevicecontrollersDelete(ctx).DeviceControllerName(request.([]string))
			return req.Execute()
		},

		UpdateRecentOps: func() {
			b.recentDeviceControllerOps = true
			b.recentDeviceControllerOpTime = time.Now()
		},
	})
}
