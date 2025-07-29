package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"sync"
	"terraform-provider-verity/openapi"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ================================================================================================
// TYPE DEFINITIONS AND CONSTANTS
// ================================================================================================

// ContextProviderFunc provides access to provider context data.
type ContextProviderFunc func() interface{}

// ClearCacheFunc clears cached data when operations complete.
type ClearCacheFunc func(ctx context.Context, provider interface{}, cacheKey string)

// OperationStatus represents the current state of a bulk operation.
type OperationStatus int

// Operation represents a single pending bulk operation.
type Operation struct {
	ResourceType  string          // The type of resource (e.g., "gateway", "tenant")
	ResourceName  string          // The name/identifier of the specific resource
	OperationType string          // The operation type ("PUT", "PATCH", "DELETE")
	Status        OperationStatus // Current status of the operation
	Error         error           // Error if operation failed
}

// ResourceExistenceCheck provides configuration for checking if resources already exist.
type ResourceExistenceCheck struct {
	FetchResources func(ctx context.Context) (map[string]interface{}, error) // Function to fetch existing resources
	ResourceType   string                                                    // Type of resource to check
	OperationType  string                                                    // Operation being performed
}

// BulkOperationConfig encapsulates all configuration needed to execute a bulk operation.
type BulkOperationConfig struct {
	ResourceType      string                                                                                                                                 // Type of resource being operated on
	OperationType     string                                                                                                                                 // Type of operation (PUT/PATCH/DELETE)
	ExtractOperations func() (map[string]interface{}, []string)                                                                                              // Extracts pending operations
	CheckPreExistence func(ctx context.Context, resourceNames []string, originalOperations map[string]interface{}) ([]string, map[string]interface{}, error) // Filters out existing resources
	PrepareRequest    func(filteredData map[string]interface{}) interface{}                                                                                  // Prepares API request
	ExecuteRequest    func(ctx context.Context, request interface{}) (*http.Response, error)                                                                 // Executes API request
	ProcessResponse   func(ctx context.Context, resp *http.Response) error                                                                                   // Processes API response
	UpdateRecentOps   func()                                                                                                                                 // Updates recent operation tracking
}

// ResourceConfig defines configuration for a specific resource type.
type ResourceConfig struct {
	ResourceType     string                                     // String identifier for the resource
	PutRequestType   reflect.Type                               // Type for PUT requests
	PatchRequestType reflect.Type                               // Type for PATCH requests
	HasAutoGen       bool                                       // Whether resource has auto-generated fields
	APIClientGetter  func(*openapi.APIClient) ResourceAPIClient // Function to get API client
}

// GenericAPIClient implements ResourceAPIClient for all resource types using reflection.
type GenericAPIClient struct {
	client       *openapi.APIClient // The underlying OpenAPI client
	resourceType string             // The resource type this client handles
}

// ResourceAPIClient provides a unified interface for all resource API operations.
type ResourceAPIClient interface {
	Put(ctx context.Context, request interface{}) (*http.Response, error)   // Execute PUT operation
	Patch(ctx context.Context, request interface{}) (*http.Response, error) // Execute PATCH operation
	Delete(ctx context.Context, names []string) (*http.Response, error)     // Execute DELETE operation
	Get(ctx context.Context) (*http.Response, error)                        // Execute GET operation
}

// Operation status constants.
const (
	OperationPending OperationStatus = iota
	OperationSucceeded
	OperationFailed
)

// Configuration constants for bulk operation timing and limits.
const (
	MaxBatchSize          = 1000                    // Maximum number of resources per batch
	DefaultBatchDelay     = 2 * time.Second         // Default delay between operations
	BatchCollectionWindow = 2000 * time.Millisecond // Time to collect operations into batches
	MaxBatchDelay         = 5 * time.Second         // Maximum time to wait before forcing execution
	OperationTimeout      = 300 * time.Second       // Timeout for individual API operations
)

// ================================================================================================
// BULK OPERATION MANAGER - MAIN STRUCT
// ================================================================================================

// BulkOperationManager coordinates bulk operations across multiple resource types.
// It provides efficient batching, automatic retry, dependency ordering, and response caching.
//
// Key Features:
//   - Generic operations for multiple resource types
//   - Automatic batching with configurable timing windows
//   - Retry logic with exponential backoff for resilient operations
//   - Response caching for auto-generated fields (tenant, service, switchpoint)
//   - Thread-safe operation queuing and execution
//   - Dependency-aware operation ordering (datacenter vs campus modes)
//   - Reflection-based generic data access patterns
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

	// Store API responses=
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

// resourceRegistry is a mapping of resource types to their configuration details.
// It provides a centralized registry for all resources that can be managed by the Verity provider.
//
// Each entry contains:
//   - ResourceType: String identifier for the resource type
//   - PutRequestType: The reflect.Type for PUT API requests for this resource
//   - PatchRequestType: The reflect.Type for PATCH API requests for this resource
//   - HasAutoGen: Boolean indicating whether the resource supports automatic generation
//   - APIClientGetter: Function that returns a ResourceAPIClient for the resource type
var resourceRegistry = map[string]ResourceConfig{
	"gateway": {
		ResourceType:     "gateway",
		PutRequestType:   reflect.TypeOf(openapi.ConfigPutRequestGatewayGatewayName{}),
		PatchRequestType: reflect.TypeOf(openapi.ConfigPutRequestGatewayGatewayName{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "gateway"}
		},
	},
	"lag": {
		ResourceType:     "lag",
		PutRequestType:   reflect.TypeOf(openapi.ConfigPutRequestLagLagName{}),
		PatchRequestType: reflect.TypeOf(openapi.ConfigPutRequestLagLagName{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "lag"}
		},
	},
	"tenant": {
		ResourceType:     "tenant",
		PutRequestType:   reflect.TypeOf(openapi.ConfigPutRequestTenantTenantName{}),
		PatchRequestType: reflect.TypeOf(openapi.ConfigPutRequestTenantTenantName{}),
		HasAutoGen:       true,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "tenant"}
		},
	},
	"service": {
		ResourceType:     "service",
		PutRequestType:   reflect.TypeOf(openapi.ConfigPutRequestServiceServiceName{}),
		PatchRequestType: reflect.TypeOf(openapi.ConfigPutRequestServiceServiceName{}),
		HasAutoGen:       true,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "service"}
		},
	},
	"gateway_profile": {
		ResourceType:     "gateway_profile",
		PutRequestType:   reflect.TypeOf(openapi.GatewayprofilesPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.GatewayprofilesPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "gateway_profile"}
		},
	},
	"packet_queue": {
		ResourceType:     "packet_queue",
		PutRequestType:   reflect.TypeOf(openapi.PacketqueuesPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.PacketqueuesPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "packet_queue"}
		},
	},
	"eth_port_profile": {
		ResourceType:     "eth_port_profile",
		PutRequestType:   reflect.TypeOf(openapi.EthportprofilesPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.EthportprofilesPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "eth_port_profile"}
		},
	},
	"eth_port_settings": {
		ResourceType:     "eth_port_settings",
		PutRequestType:   reflect.TypeOf(openapi.EthportsettingsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.EthportsettingsPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "eth_port_settings"}
		},
	},
	"bundle": {
		ResourceType:     "bundle",
		PutRequestType:   reflect.TypeOf(openapi.BundlesPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.BundlesPatchRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "bundle"}
		},
	},
	"acl": {
		ResourceType:     "acl",
		PutRequestType:   reflect.TypeOf(openapi.AclsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.AclsPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "acl"}
		},
	},
	"packet_broker": {
		ResourceType:     "packet_broker",
		PutRequestType:   reflect.TypeOf(openapi.PacketbrokerPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.PacketbrokerPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "packet_broker"}
		},
	},
	"badge": {
		ResourceType:     "badge",
		PutRequestType:   reflect.TypeOf(openapi.BadgesPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.BadgesPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "badge"}
		},
	},
	"switchpoint": {
		ResourceType:     "switchpoint",
		PutRequestType:   reflect.TypeOf(openapi.SwitchpointsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.SwitchpointsPutRequest{}),
		HasAutoGen:       true,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "switchpoint"}
		},
	},
	"device_controller": {
		ResourceType:     "device_controller",
		PutRequestType:   reflect.TypeOf(openapi.DevicecontrollersPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.DevicecontrollersPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "device_controller"}
		},
	},
	"authenticated_eth_port": {
		ResourceType:     "authenticated_eth_port",
		PutRequestType:   reflect.TypeOf(openapi.AuthenticatedethportsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.AuthenticatedethportsPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "authenticated_eth_port"}
		},
	},
	"device_voice_settings": {
		ResourceType:     "device_voice_settings",
		PutRequestType:   reflect.TypeOf(openapi.DevicevoicesettingsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.DevicevoicesettingsPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "device_voice_settings"}
		},
	},
	"voice_port_profile": {
		ResourceType:     "voice_port_profile",
		PutRequestType:   reflect.TypeOf(openapi.VoiceportprofilesPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.VoiceportprofilesPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "voice_port_profile"}
		},
	},
	"service_port_profile": {
		ResourceType:     "service_port_profile",
		PutRequestType:   reflect.TypeOf(openapi.ServiceportprofilesPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.ServiceportprofilesPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "service_port_profile"}
		},
	},
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
	if !execute("PUT", len(b.tenantPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "tenant", "PUT") }, "Tenant") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.gatewayPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "gateway", "PUT") }, "Gateway") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.gatewayProfilePut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "gateway_profile", "PUT") }, "Gateway Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.servicePut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "service", "PUT") }, "Service") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.packetQueuePut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "packet_queue", "PUT") }, "Packet Queue") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.ethPortProfilePut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "eth_port_profile", "PUT") }, "Eth Port Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.ethPortSettingsPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "eth_port_settings", "PUT") }, "Eth Port Settings") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.lagPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "lag", "PUT") }, "LAG") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.bundlePut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "bundle", "PUT") }, "Bundle") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.aclPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "acl", "PUT") }, "ACL") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.packetBrokerPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "packet_broker", "PUT") }, "Packet Broker") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.badgePut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "badge", "PUT") }, "Badge") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.switchpointPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "switchpoint", "PUT") }, "Switchpoint") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.deviceControllerPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "device_controller", "PUT") }, "Device Controller") {
		return diagnostics, operationsPerformed
	}

	// PATCH operations - DC Order
	tflog.Debug(ctx, "Checking PATCH operation counts", map[string]interface{}{
		"tenant_patch_count":            len(b.tenantPatch),
		"gateway_patch_count":           len(b.gatewayPatch),
		"gateway_profile_patch_count":   len(b.gatewayProfilePatch),
		"service_patch_count":           len(b.servicePatch),
		"packet_queue_patch_count":      len(b.packetQueuePatch),
		"eth_port_profile_patch_count":  len(b.ethPortProfilePatch),
		"eth_port_settings_patch_count": len(b.ethPortSettingsPatch),
		"lag_patch_count":               len(b.lagPatch),
		"bundle_patch_count":            len(b.bundlePatch),
		"acl_patch_count":               len(b.aclPatch),
		"packet_broker_patch_count":     len(b.packetBrokerPatch),
		"badge_patch_count":             len(b.badgePatch),
		"device_controller_patch_count": len(b.deviceControllerPatch),
	})

	if !execute("PATCH", len(b.tenantPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "tenant", "PATCH") }, "Tenant") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.gatewayPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "gateway", "PATCH") }, "Gateway") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.gatewayProfilePatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "gateway_profile", "PATCH") }, "Gateway Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.servicePatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "service", "PATCH") }, "Service") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.packetQueuePatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "packet_queue", "PATCH") }, "Packet Queue") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.ethPortProfilePatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "eth_port_profile", "PATCH") }, "Eth Port Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.ethPortSettingsPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "eth_port_settings", "PATCH") }, "Eth Port Settings") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.lagPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "lag", "PATCH") }, "LAG") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.bundlePatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "bundle", "PATCH") }, "Bundle") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.aclPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "acl", "PATCH") }, "ACL") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.packetBrokerPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "packet_broker", "PATCH") }, "Packet Broker") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.badgePatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "badge", "PATCH") }, "Badge") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.switchpointPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "switchpoint", "PATCH") }, "Switchpoint") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.deviceControllerPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "device_controller", "PATCH") }, "Device Controller") {
		return diagnostics, operationsPerformed
	}

	// DELETE operations - Reverse DC Order
	if !execute("DELETE", len(b.deviceControllerDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "device_controller", "DELETE") }, "Device Controller") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.switchpointDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "switchpoint", "DELETE") }, "Switchpoint") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.badgeDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "badge", "DELETE") }, "Badge") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.packetBrokerDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "packet_broker", "DELETE") }, "Packet Broker") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.aclDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "acl", "DELETE") }, "ACL") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.bundleDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "bundle", "DELETE") }, "Bundle") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.lagDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "lag", "DELETE") }, "LAG") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.ethPortSettingsDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "eth_port_settings", "DELETE") }, "Eth Port Settings") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.ethPortProfileDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "eth_port_profile", "DELETE") }, "Eth Port Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.packetQueueDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "packet_queue", "DELETE") }, "Packet Queue") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.serviceDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "service", "DELETE") }, "Service") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.gatewayProfileDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "gateway_profile", "DELETE") }, "Gateway Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.gatewayDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "gateway", "DELETE") }, "Gateway") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.tenantDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "tenant", "DELETE") }, "Tenant") {
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
	if !execute("PUT", len(b.servicePut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "service", "PUT") }, "Service") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.ethPortProfilePut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "eth_port_profile", "PUT") }, "Eth Port Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.authenticatedEthPortPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "authenticated_eth_port", "PUT") }, "Authenticated Eth-Port") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.deviceVoiceSettingsPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "device_voice_settings", "PUT") }, "Device Voice Settings") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.packetQueuePut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "packet_queue", "PUT") }, "Packet Queue") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.servicePortProfilePut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "service_port_profile", "PUT") }, "Service Port Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.voicePortProfilePut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "voice_port_profile", "PUT") }, "Voice-Port Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.ethPortSettingsPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "eth_port_settings", "PUT") }, "Eth Port Settings") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.lagPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "lag", "PUT") }, "LAG") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.bundlePut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "bundle", "PUT") }, "Bundle") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.badgePut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "badge", "PUT") }, "Badge") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.switchpointPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "switchpoint", "PUT") }, "Switchpoint") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.deviceControllerPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "device_controller", "PUT") }, "Device Controller") {
		return diagnostics, operationsPerformed
	}

	// PATCH operations - Campus Order
	if !execute("PATCH", len(b.servicePatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "service", "PATCH") }, "Service") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.ethPortProfilePatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "eth_port_profile", "PATCH") }, "Eth Port Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.authenticatedEthPortPatch), func(ctx context.Context) diag.Diagnostics {
		return b.ExecuteBulk(ctx, "authenticated_eth_port", "PATCH")
	}, "Authenticated Eth-Port") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.deviceVoiceSettingsPatch), func(ctx context.Context) diag.Diagnostics {
		return b.ExecuteBulk(ctx, "device_voice_settings", "PATCH")
	}, "Device Voice Settings") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.packetQueuePatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "packet_queue", "PATCH") }, "Packet Queue") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.servicePortProfilePatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "service_port_profile", "PATCH") }, "Service Port Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.voicePortProfilePatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "voice_port_profile", "PATCH") }, "Voice-Port Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.ethPortSettingsPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "eth_port_settings", "PATCH") }, "Eth Port Settings") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.lagPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "lag", "PATCH") }, "LAG") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.bundlePatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "bundle", "PATCH") }, "Bundle") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.badgePatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "badge", "PATCH") }, "Badge") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.switchpointPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "switchpoint", "PATCH") }, "Switchpoint") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.deviceControllerPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "device_controller", "PATCH") }, "Device Controller") {
		return diagnostics, operationsPerformed
	}

	// DELETE operations - Reverse Campus Order
	if !execute("DELETE", len(b.deviceControllerDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "device_controller", "DELETE") }, "Device Controller") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.switchpointDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "switchpoint", "DELETE") }, "Switchpoint") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.badgeDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "badge", "DELETE") }, "Badge") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.bundleDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "bundle", "DELETE") }, "Bundle") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.lagDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "lag", "DELETE") }, "LAG") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.ethPortSettingsDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "eth_port_settings", "DELETE") }, "Eth Port Settings") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.voicePortProfileDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "voice_port_profile", "DELETE") }, "Voice-Port Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.servicePortProfileDelete), func(ctx context.Context) diag.Diagnostics {
		return b.ExecuteBulk(ctx, "service_port_profile", "DELETE")
	}, "Service Port Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.packetQueueDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "packet_queue", "DELETE") }, "Packet Queue") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.deviceVoiceSettingsDelete), func(ctx context.Context) diag.Diagnostics {
		return b.ExecuteBulk(ctx, "device_voice_settings", "DELETE")
	}, "Device Voice Settings") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.authenticatedEthPortDelete), func(ctx context.Context) diag.Diagnostics {
		return b.ExecuteBulk(ctx, "authenticated_eth_port", "DELETE")
	}, "Authenticated Eth-Port") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.ethPortProfileDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "eth_port_profile", "DELETE") }, "Eth Port Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.serviceDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "service", "DELETE") }, "Service") {
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

// ResourceOperationData holds operation data for a resource type
type ResourceOperationData struct {
	PutOperations    interface{}
	PatchOperations  interface{}
	DeleteOperations *[]string
	RecentOps        *bool
	RecentOpTime     *time.Time
}

// GetResourceOperationData returns operation data for a resource type
func (b *BulkOperationManager) GetResourceOperationData(resourceType string) *ResourceOperationData {
	switch resourceType {
	case "gateway":
		return &ResourceOperationData{
			PutOperations:    b.gatewayPut,
			PatchOperations:  b.gatewayPatch,
			DeleteOperations: &b.gatewayDelete,
			RecentOps:        &b.recentGatewayOps,
			RecentOpTime:     &b.recentGatewayOpTime,
		}
	case "lag":
		return &ResourceOperationData{
			PutOperations:    b.lagPut,
			PatchOperations:  b.lagPatch,
			DeleteOperations: &b.lagDelete,
			RecentOps:        &b.recentLagOps,
			RecentOpTime:     &b.recentLagOpTime,
		}
	case "tenant":
		return &ResourceOperationData{
			PutOperations:    b.tenantPut,
			PatchOperations:  b.tenantPatch,
			DeleteOperations: &b.tenantDelete,
			RecentOps:        &b.recentTenantOps,
			RecentOpTime:     &b.recentTenantOpTime,
		}
	case "service":
		return &ResourceOperationData{
			PutOperations:    b.servicePut,
			PatchOperations:  b.servicePatch,
			DeleteOperations: &b.serviceDelete,
			RecentOps:        &b.recentServiceOps,
			RecentOpTime:     &b.recentServiceOpTime,
		}
	case "gateway_profile":
		return &ResourceOperationData{
			PutOperations:    b.gatewayProfilePut,
			PatchOperations:  b.gatewayProfilePatch,
			DeleteOperations: &b.gatewayProfileDelete,
			RecentOps:        &b.recentGatewayProfileOps,
			RecentOpTime:     &b.recentGatewayProfileOpTime,
		}
	case "eth_port_profile":
		return &ResourceOperationData{
			PutOperations:    b.ethPortProfilePut,
			PatchOperations:  b.ethPortProfilePatch,
			DeleteOperations: &b.ethPortProfileDelete,
			RecentOps:        &b.recentEthPortProfileOps,
			RecentOpTime:     &b.recentEthPortProfileOpTime,
		}
	case "eth_port_settings":
		return &ResourceOperationData{
			PutOperations:    b.ethPortSettingsPut,
			PatchOperations:  b.ethPortSettingsPatch,
			DeleteOperations: &b.ethPortSettingsDelete,
			RecentOps:        &b.recentEthPortSettingsOps,
			RecentOpTime:     &b.recentEthPortSettingsOpTime,
		}
	case "bundle":
		return &ResourceOperationData{
			PutOperations:    b.bundlePut,
			PatchOperations:  b.bundlePatch,
			DeleteOperations: &b.bundleDelete,
			RecentOps:        &b.recentBundleOps,
			RecentOpTime:     &b.recentBundleOpTime,
		}
	case "acl":
		return &ResourceOperationData{
			PutOperations:    b.aclPut,
			PatchOperations:  b.aclPatch,
			DeleteOperations: &b.aclDelete,
			RecentOps:        &b.recentAclOps,
			RecentOpTime:     &b.recentAclOpTime,
		}
	case "authenticated_eth_port":
		return &ResourceOperationData{
			PutOperations:    b.authenticatedEthPortPut,
			PatchOperations:  b.authenticatedEthPortPatch,
			DeleteOperations: &b.authenticatedEthPortDelete,
			RecentOps:        &b.recentAuthenticatedEthPortOps,
			RecentOpTime:     &b.recentAuthenticatedEthPortOpTime,
		}
	case "badge":
		return &ResourceOperationData{
			PutOperations:    b.badgePut,
			PatchOperations:  b.badgePatch,
			DeleteOperations: &b.badgeDelete,
			RecentOps:        &b.recentBadgeOps,
			RecentOpTime:     &b.recentBadgeOpTime,
		}
	case "device_voice_settings":
		return &ResourceOperationData{
			PutOperations:    b.deviceVoiceSettingsPut,
			PatchOperations:  b.deviceVoiceSettingsPatch,
			DeleteOperations: &b.deviceVoiceSettingsDelete,
			RecentOps:        &b.recentDeviceVoiceSettingsOps,
			RecentOpTime:     &b.recentDeviceVoiceSettingsOpTime,
		}
	case "packet_broker":
		return &ResourceOperationData{
			PutOperations:    b.packetBrokerPut,
			PatchOperations:  b.packetBrokerPatch,
			DeleteOperations: &b.packetBrokerDelete,
			RecentOps:        &b.recentPacketBrokerOps,
			RecentOpTime:     &b.recentPacketBrokerOpTime,
		}
	case "packet_queue":
		return &ResourceOperationData{
			PutOperations:    b.packetQueuePut,
			PatchOperations:  b.packetQueuePatch,
			DeleteOperations: &b.packetQueueDelete,
			RecentOps:        &b.recentPacketQueueOps,
			RecentOpTime:     &b.recentPacketQueueOpTime,
		}
	case "service_port_profile":
		return &ResourceOperationData{
			PutOperations:    b.servicePortProfilePut,
			PatchOperations:  b.servicePortProfilePatch,
			DeleteOperations: &b.servicePortProfileDelete,
			RecentOps:        &b.recentServicePortProfileOps,
			RecentOpTime:     &b.recentServicePortProfileOpTime,
		}
	case "switchpoint":
		return &ResourceOperationData{
			PutOperations:    b.switchpointPut,
			PatchOperations:  b.switchpointPatch,
			DeleteOperations: &b.switchpointDelete,
			RecentOps:        &b.recentSwitchpointOps,
			RecentOpTime:     &b.recentSwitchpointOpTime,
		}
	case "voice_port_profile":
		return &ResourceOperationData{
			PutOperations:    b.voicePortProfilePut,
			PatchOperations:  b.voicePortProfilePatch,
			DeleteOperations: &b.voicePortProfileDelete,
			RecentOps:        &b.recentVoicePortProfileOps,
			RecentOpTime:     &b.recentVoicePortProfileOpTime,
		}
	case "device_controller":
		return &ResourceOperationData{
			PutOperations:    b.deviceControllerPut,
			PatchOperations:  b.deviceControllerPatch,
			DeleteOperations: &b.deviceControllerDelete,
			RecentOps:        &b.recentDeviceControllerOps,
			RecentOpTime:     &b.recentDeviceControllerOpTime,
		}
	}
	return nil
}

func (b *BulkOperationManager) hasPendingOrRecentOperations(
	resourceType string,
) bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	data := b.GetResourceOperationData(resourceType)
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

	hasPending = hasPending || len(*data.DeleteOperations) > 0

	// Check if we've recently had operations (within the last 5 seconds)
	hasRecent := *data.RecentOps && time.Since(*data.RecentOpTime) < 5*time.Second

	return hasPending || hasRecent
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

	retryConfig := DefaultRetryConfig()
	var opErr error
	var apiResp *http.Response

	for retry := 0; retry < retryConfig.MaxRetries; retry++ {
		if retry > 0 {
			delayTime := CalculateBackoff(retry, retryConfig)
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

		if !IsRetriableError(opErr) {
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
		if processErr := config.ProcessResponse(ctx, apiResp); processErr != nil {
			tflog.Warn(ctx, fmt.Sprintf("Post-processing failed for bulk %s %s operation: %v",
				config.ResourceType, config.OperationType, processErr))
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

// GenericAPIClient implements ResourceAPIClient for all resource types
func (g *GenericAPIClient) Put(ctx context.Context, request interface{}) (*http.Response, error) {
	switch g.resourceType {
	case "gateway":
		req := g.client.GatewaysAPI.GatewaysPut(ctx).GatewaysPutRequest(*request.(*openapi.GatewaysPutRequest))
		return req.Execute()
	case "lag":
		req := g.client.LAGsAPI.LagsPut(ctx).LagsPutRequest(*request.(*openapi.LagsPutRequest))
		return req.Execute()
	case "tenant":
		req := g.client.TenantsAPI.TenantsPut(ctx).TenantsPutRequest(*request.(*openapi.TenantsPutRequest))
		return req.Execute()
	case "service":
		req := g.client.ServicesAPI.ServicesPut(ctx).ServicesPutRequest(*request.(*openapi.ServicesPutRequest))
		return req.Execute()
	case "gateway_profile":
		req := g.client.GatewayProfilesAPI.GatewayprofilesPut(ctx).GatewayprofilesPutRequest(*request.(*openapi.GatewayprofilesPutRequest))
		return req.Execute()
	case "packet_queue":
		req := g.client.PacketQueuesAPI.PacketqueuesPut(ctx).PacketqueuesPutRequest(*request.(*openapi.PacketqueuesPutRequest))
		return req.Execute()
	case "eth_port_profile":
		req := g.client.EthPortProfilesAPI.EthportprofilesPut(ctx).EthportprofilesPutRequest(*request.(*openapi.EthportprofilesPutRequest))
		return req.Execute()
	case "eth_port_settings":
		req := g.client.EthPortSettingsAPI.EthportsettingsPut(ctx).EthportsettingsPutRequest(*request.(*openapi.EthportsettingsPutRequest))
		return req.Execute()
	case "bundle":
		req := g.client.BundlesAPI.BundlesPut(ctx).BundlesPutRequest(*request.(*openapi.BundlesPutRequest))
		return req.Execute()
	case "acl":
		req := g.client.ACLsAPI.AclsPut(ctx).AclsPutRequest(*request.(*openapi.AclsPutRequest))
		return req.Execute()
	case "packet_broker":
		req := g.client.PacketBrokerAPI.PacketbrokerPut(ctx).PacketbrokerPutRequest(*request.(*openapi.PacketbrokerPutRequest))
		return req.Execute()
	case "badge":
		req := g.client.BadgesAPI.BadgesPut(ctx).BadgesPutRequest(*request.(*openapi.BadgesPutRequest))
		return req.Execute()
	case "switchpoint":
		req := g.client.SwitchpointsAPI.SwitchpointsPut(ctx).SwitchpointsPutRequest(*request.(*openapi.SwitchpointsPutRequest))
		return req.Execute()
	case "device_controller":
		req := g.client.DeviceControllersAPI.DevicecontrollersPut(ctx).DevicecontrollersPutRequest(*request.(*openapi.DevicecontrollersPutRequest))
		return req.Execute()
	case "authenticated_eth_port":
		req := g.client.AuthenticatedEthPortsAPI.AuthenticatedethportsPut(ctx).AuthenticatedethportsPutRequest(*request.(*openapi.AuthenticatedethportsPutRequest))
		return req.Execute()
	case "device_voice_settings":
		req := g.client.DeviceVoiceSettingsAPI.DevicevoicesettingsPut(ctx).DevicevoicesettingsPutRequest(*request.(*openapi.DevicevoicesettingsPutRequest))
		return req.Execute()
	case "voice_port_profile":
		req := g.client.VoicePortProfilesAPI.VoiceportprofilesPut(ctx).VoiceportprofilesPutRequest(*request.(*openapi.VoiceportprofilesPutRequest))
		return req.Execute()
	case "service_port_profile":
		req := g.client.ServicePortProfilesAPI.ServiceportprofilesPut(ctx).ServiceportprofilesPutRequest(*request.(*openapi.ServiceportprofilesPutRequest))
		return req.Execute()
	default:
		return nil, fmt.Errorf("unknown resource type: %s", g.resourceType)
	}
}

func (g *GenericAPIClient) Patch(ctx context.Context, request interface{}) (*http.Response, error) {
	switch g.resourceType {
	case "gateway":
		req := g.client.GatewaysAPI.GatewaysPatch(ctx).GatewaysPutRequest(*request.(*openapi.GatewaysPutRequest))
		return req.Execute()
	case "lag":
		req := g.client.LAGsAPI.LagsPatch(ctx).LagsPutRequest(*request.(*openapi.LagsPutRequest))
		return req.Execute()
	case "tenant":
		req := g.client.TenantsAPI.TenantsPatch(ctx).TenantsPutRequest(*request.(*openapi.TenantsPutRequest))
		return req.Execute()
	case "service":
		req := g.client.ServicesAPI.ServicesPatch(ctx).ServicesPutRequest(*request.(*openapi.ServicesPutRequest))
		return req.Execute()
	case "gateway_profile":
		req := g.client.GatewayProfilesAPI.GatewayprofilesPatch(ctx).GatewayprofilesPutRequest(*request.(*openapi.GatewayprofilesPutRequest))
		return req.Execute()
	case "packet_queue":
		req := g.client.PacketQueuesAPI.PacketqueuesPatch(ctx).PacketqueuesPutRequest(*request.(*openapi.PacketqueuesPutRequest))
		return req.Execute()
	case "eth_port_profile":
		req := g.client.EthPortProfilesAPI.EthportprofilesPatch(ctx).EthportprofilesPutRequest(*request.(*openapi.EthportprofilesPutRequest))
		return req.Execute()
	case "eth_port_settings":
		req := g.client.EthPortSettingsAPI.EthportsettingsPatch(ctx).EthportsettingsPutRequest(*request.(*openapi.EthportsettingsPutRequest))
		return req.Execute()
	case "bundle":
		req := g.client.BundlesAPI.BundlesPatch(ctx).BundlesPatchRequest(*request.(*openapi.BundlesPatchRequest))
		return req.Execute()
	case "acl":
		req := g.client.ACLsAPI.AclsPatch(ctx).AclsPutRequest(*request.(*openapi.AclsPutRequest))
		return req.Execute()
	case "packet_broker":
		req := g.client.PacketBrokerAPI.PacketbrokerPatch(ctx).PacketbrokerPutRequest(*request.(*openapi.PacketbrokerPutRequest))
		return req.Execute()
	case "badge":
		req := g.client.BadgesAPI.BadgesPatch(ctx).BadgesPutRequest(*request.(*openapi.BadgesPutRequest))
		return req.Execute()
	case "switchpoint":
		req := g.client.SwitchpointsAPI.SwitchpointsPatch(ctx).SwitchpointsPutRequest(*request.(*openapi.SwitchpointsPutRequest))
		return req.Execute()
	case "device_controller":
		req := g.client.DeviceControllersAPI.DevicecontrollersPatch(ctx).DevicecontrollersPutRequest(*request.(*openapi.DevicecontrollersPutRequest))
		return req.Execute()
	case "authenticated_eth_port":
		req := g.client.AuthenticatedEthPortsAPI.AuthenticatedethportsPatch(ctx).AuthenticatedethportsPutRequest(*request.(*openapi.AuthenticatedethportsPutRequest))
		return req.Execute()
	case "device_voice_settings":
		req := g.client.DeviceVoiceSettingsAPI.DevicevoicesettingsPatch(ctx).DevicevoicesettingsPutRequest(*request.(*openapi.DevicevoicesettingsPutRequest))
		return req.Execute()
	case "voice_port_profile":
		req := g.client.VoicePortProfilesAPI.VoiceportprofilesPatch(ctx).VoiceportprofilesPutRequest(*request.(*openapi.VoiceportprofilesPutRequest))
		return req.Execute()
	case "service_port_profile":
		req := g.client.ServicePortProfilesAPI.ServiceportprofilesPatch(ctx).ServiceportprofilesPutRequest(*request.(*openapi.ServiceportprofilesPutRequest))
		return req.Execute()
	default:
		return nil, fmt.Errorf("unknown resource type: %s", g.resourceType)
	}
}

func (g *GenericAPIClient) Delete(ctx context.Context, names []string) (*http.Response, error) {
	switch g.resourceType {
	case "gateway":
		req := g.client.GatewaysAPI.GatewaysDelete(ctx).GatewayName(names)
		return req.Execute()
	case "lag":
		req := g.client.LAGsAPI.LagsDelete(ctx).LagName(names)
		return req.Execute()
	case "tenant":
		req := g.client.TenantsAPI.TenantsDelete(ctx).TenantName(names)
		return req.Execute()
	case "service":
		req := g.client.ServicesAPI.ServicesDelete(ctx).ServiceName(names)
		return req.Execute()
	case "gateway_profile":
		req := g.client.GatewayProfilesAPI.GatewayprofilesDelete(ctx).ProfileName(names)
		return req.Execute()
	case "packet_queue":
		req := g.client.PacketQueuesAPI.PacketqueuesDelete(ctx).PacketQueueName(names)
		return req.Execute()
	case "eth_port_profile":
		req := g.client.EthPortProfilesAPI.EthportprofilesDelete(ctx).ProfileName(names)
		return req.Execute()
	case "eth_port_settings":
		req := g.client.EthPortSettingsAPI.EthportsettingsDelete(ctx).PortName(names)
		return req.Execute()
	case "bundle":
		req := g.client.BundlesAPI.BundlesDelete(ctx).BundleName(names)
		return req.Execute()
	case "acl":
		req := g.client.ACLsAPI.AclsDelete(ctx).IpFilterName(names)
		return req.Execute()
	case "packet_broker":
		req := g.client.PacketBrokerAPI.PacketbrokerDelete(ctx).PbEgressProfileName(names)
		return req.Execute()
	case "badge":
		req := g.client.BadgesAPI.BadgesDelete(ctx).BadgeName(names)
		return req.Execute()
	case "switchpoint":
		req := g.client.SwitchpointsAPI.SwitchpointsDelete(ctx).SwitchpointName(names)
		return req.Execute()
	case "device_controller":
		req := g.client.DeviceControllersAPI.DevicecontrollersDelete(ctx).DeviceControllerName(names)
		return req.Execute()
	case "authenticated_eth_port":
		req := g.client.AuthenticatedEthPortsAPI.AuthenticatedethportsDelete(ctx).AuthenticatedEthPortName(names)
		return req.Execute()
	case "device_voice_settings":
		req := g.client.DeviceVoiceSettingsAPI.DevicevoicesettingsDelete(ctx).DeviceVoiceSettingsName(names)
		return req.Execute()
	case "voice_port_profile":
		req := g.client.VoicePortProfilesAPI.VoiceportprofilesDelete(ctx).VoicePortProfileName(names)
		return req.Execute()
	case "service_port_profile":
		req := g.client.ServicePortProfilesAPI.ServiceportprofilesDelete(ctx).ServicePortProfileName(names)
		return req.Execute()
	default:
		return nil, fmt.Errorf("unknown resource type: %s", g.resourceType)
	}
}

func (g *GenericAPIClient) Get(ctx context.Context) (*http.Response, error) {
	switch g.resourceType {
	case "gateway":
		req := g.client.GatewaysAPI.GatewaysGet(ctx)
		return req.Execute()
	case "lag":
		req := g.client.LAGsAPI.LagsGet(ctx)
		return req.Execute()
	case "tenant":
		req := g.client.TenantsAPI.TenantsGet(ctx)
		return req.Execute()
	case "service":
		req := g.client.ServicesAPI.ServicesGet(ctx)
		return req.Execute()
	case "gateway_profile":
		req := g.client.GatewayProfilesAPI.GatewayprofilesGet(ctx)
		return req.Execute()
	case "packet_queue":
		req := g.client.PacketQueuesAPI.PacketqueuesGet(ctx)
		return req.Execute()
	case "eth_port_profile":
		req := g.client.EthPortProfilesAPI.EthportprofilesGet(ctx)
		return req.Execute()
	case "eth_port_settings":
		req := g.client.EthPortSettingsAPI.EthportsettingsGet(ctx)
		return req.Execute()
	case "bundle":
		req := g.client.BundlesAPI.BundlesGet(ctx)
		return req.Execute()
	case "acl":
		req := g.client.ACLsAPI.AclsGet(ctx)
		return req.Execute()
	case "packet_broker":
		req := g.client.PacketBrokerAPI.PacketbrokerGet(ctx)
		return req.Execute()
	case "badge":
		req := g.client.BadgesAPI.BadgesGet(ctx)
		return req.Execute()
	case "switchpoint":
		req := g.client.SwitchpointsAPI.SwitchpointsGet(ctx)
		return req.Execute()
	case "device_controller":
		req := g.client.DeviceControllersAPI.DevicecontrollersGet(ctx)
		return req.Execute()
	case "authenticated_eth_port":
		req := g.client.AuthenticatedEthPortsAPI.AuthenticatedethportsGet(ctx)
		return req.Execute()
	case "device_voice_settings":
		req := g.client.DeviceVoiceSettingsAPI.DevicevoicesettingsGet(ctx)
		return req.Execute()
	case "voice_port_profile":
		req := g.client.VoicePortProfilesAPI.VoiceportprofilesGet(ctx)
		return req.Execute()
	case "service_port_profile":
		req := g.client.ServicePortProfilesAPI.ServiceportprofilesGet(ctx)
		return req.Execute()
	default:
		return nil, fmt.Errorf("unknown resource type: %s", g.resourceType)
	}
}

// ================================================================================================
// GENERIC OPERATIONS
// ================================================================================================
func (b *BulkOperationManager) AddPut(ctx context.Context, resourceType, resourceName string, props interface{}) string {
	return b.addGenericOperation(ctx, resourceType, resourceName, "PUT", props, "")
}

func (b *BulkOperationManager) AddPatch(ctx context.Context, resourceType, resourceName string, props interface{}) string {
	return b.addGenericOperation(ctx, resourceType, resourceName, "PATCH", props, "")
}

func (b *BulkOperationManager) AddDelete(ctx context.Context, resourceType, resourceName string) string {
	return b.addGenericOperation(ctx, resourceType, resourceName, "DELETE", nil, "")
}

// Special ACL methods that handle IP version tracking
func (b *BulkOperationManager) AddAclPut(ctx context.Context, aclName string, props openapi.ConfigPutRequestIpv4FilterIpv4FilterName, ipVersion string) string {
	return b.addGenericOperation(ctx, "acl", aclName, "PUT", props, ipVersion)
}

func (b *BulkOperationManager) AddAclPatch(ctx context.Context, aclName string, props openapi.ConfigPutRequestIpv4FilterIpv4FilterName, ipVersion string) string {
	return b.addGenericOperation(ctx, "acl", aclName, "PATCH", props, ipVersion)
}

func (b *BulkOperationManager) AddAclDelete(ctx context.Context, aclName string, ipVersion string) string {
	return b.addGenericOperation(ctx, "acl", aclName, "DELETE", nil, ipVersion)
}

// Internal method that handles both generic operations and ACL-specific IP version tracking
func (b *BulkOperationManager) addGenericOperation(ctx context.Context, resourceType, resourceName, operationType string, props interface{}, ipVersion string) string {
	// For ACL operations, create a composite key that includes the IP version to prevent overwrites
	storeKey := resourceName
	if resourceType == "acl" && ipVersion != "" {
		storeKey = fmt.Sprintf("%s_v%s", resourceName, ipVersion)
	}

	storeFunc := func() {
		b.storeOperation(resourceType, storeKey, operationType, props)
		if resourceType == "acl" && ipVersion != "" {
			if b.aclIpVersion == nil {
				b.aclIpVersion = make(map[string]string)
			}
			b.aclIpVersion[storeKey] = ipVersion
		}
	}

	logDetails := map[string]interface{}{
		fmt.Sprintf("%s_name", resourceType): resourceName,
		"batch_size":                         b.getBatchSize(resourceType, operationType) + 1,
	}

	if resourceType == "acl" && ipVersion != "" {
		logDetails["ip_version"] = ipVersion
	}

	return b.addOperation(ctx, resourceType, resourceName, operationType, storeFunc, logDetails)
}

func (b *BulkOperationManager) ExecuteBulk(ctx context.Context, resourceType, operationType string) diag.Diagnostics {
	// Special handling for ACLs which need IP version separation
	if resourceType == "acl" {
		return b.executeBulkAcl(ctx, operationType)
	}

	config, exists := resourceRegistry[resourceType]
	if !exists {
		var diags diag.Diagnostics
		diags.AddError("Unknown resource type", fmt.Sprintf("Resource type %s is not registered", resourceType))
		return diags
	}

	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:      resourceType,
		OperationType:     operationType,
		ExtractOperations: b.createExtractor(resourceType, operationType),
		CheckPreExistence: b.createPreExistenceChecker(config, operationType),
		PrepareRequest:    b.createRequestPreparer(config, operationType),
		ExecuteRequest:    b.createRequestExecutor(config, operationType),
		ProcessResponse:   b.createResponseProcessor(config, operationType),
		UpdateRecentOps:   b.createRecentOpsUpdater(resourceType),
	})
}

// Special handling for ACL operations that need to be executed separately
func (b *BulkOperationManager) executeBulkAcl(ctx context.Context, operationType string) diag.Diagnostics {
	var diagnostics diag.Diagnostics
	b.mutex.Lock()
	var originalOperations map[string]openapi.ConfigPutRequestIpv4FilterIpv4FilterName
	var originalIpVersions map[string]string

	switch operationType {
	case "PUT":
		originalOperations = make(map[string]openapi.ConfigPutRequestIpv4FilterIpv4FilterName)
		for k, v := range b.aclPut {
			originalOperations[k] = v
		}
		b.aclPut = make(map[string]openapi.ConfigPutRequestIpv4FilterIpv4FilterName)
	case "PATCH":
		originalOperations = make(map[string]openapi.ConfigPutRequestIpv4FilterIpv4FilterName)
		for k, v := range b.aclPatch {
			originalOperations[k] = v
		}
		b.aclPatch = make(map[string]openapi.ConfigPutRequestIpv4FilterIpv4FilterName)
	case "DELETE":
		originalOperations = make(map[string]openapi.ConfigPutRequestIpv4FilterIpv4FilterName)
		for _, name := range b.aclDelete {
			originalOperations[name] = openapi.ConfigPutRequestIpv4FilterIpv4FilterName{}
		}
		b.aclDelete = b.aclDelete[:0]
	}

	originalIpVersions = make(map[string]string)
	for k, v := range b.aclIpVersion {
		if _, exists := originalOperations[k]; exists {
			originalIpVersions[k] = v
		}
	}

	for k := range originalOperations {
		delete(b.aclIpVersion, k)
	}
	b.mutex.Unlock()

	if len(originalOperations) == 0 {
		return diagnostics
	}

	ipv4Data := make(map[string]openapi.ConfigPutRequestIpv4FilterIpv4FilterName)
	ipv6Data := make(map[string]openapi.ConfigPutRequestIpv4FilterIpv4FilterName)

	// Process operations and extract original resource names from composite keys
	for compositeKey, props := range originalOperations {
		ipVersion := originalIpVersions[compositeKey]

		// Extract original resource name by removing the _v4 or _v6 suffix
		var originalName string
		if ipVersion == "4" && strings.HasSuffix(compositeKey, "_v4") {
			originalName = strings.TrimSuffix(compositeKey, "_v4")
		} else if ipVersion == "6" && strings.HasSuffix(compositeKey, "_v6") {
			originalName = strings.TrimSuffix(compositeKey, "_v6")
		} else {
			// Fallback for operations without composite keys
			originalName = compositeKey
		}

		if ipVersion == "6" {
			ipv6Data[originalName] = props
		} else {
			ipv4Data[originalName] = props
		}
	}

	// Process IPv4 ACLs
	if len(ipv4Data) > 0 {
		ipv4Diagnostics := b.executeAclForIpVersion(ctx, ipv4Data, operationType, "4")
		diagnostics = append(diagnostics, ipv4Diagnostics...)
	}

	// Process IPv6 ACLs
	if len(ipv6Data) > 0 {
		ipv6Diagnostics := b.executeAclForIpVersion(ctx, ipv6Data, operationType, "6")
		diagnostics = append(diagnostics, ipv6Diagnostics...)
	}

	// Update recent operations
	b.recentAclOps = true
	b.recentAclOpTime = time.Now()

	return diagnostics
}

func (b *BulkOperationManager) executeAclForIpVersion(ctx context.Context, aclData map[string]openapi.ConfigPutRequestIpv4FilterIpv4FilterName, operationType, ipVersion string) diag.Diagnostics {
	return b.executeBulkOperation(ctx, BulkOperationConfig{
		ResourceType:  fmt.Sprintf("acl_v%s", ipVersion),
		OperationType: operationType,

		ExtractOperations: func() (map[string]interface{}, []string) {
			result := make(map[string]interface{})
			names := make([]string, 0, len(aclData))
			for k, v := range aclData {
				result[k] = v
				names = append(names, k)
			}
			return result, names
		},

		CheckPreExistence: nil,

		PrepareRequest: func(filteredData map[string]interface{}) interface{} {
			if operationType == "DELETE" {
				names := make([]string, 0, len(filteredData))
				for name := range filteredData {
					names = append(names, name)
				}
				return names
			} else {
				request := openapi.NewAclsPutRequest()
				aclMap := make(map[string]openapi.ConfigPutRequestIpv4FilterIpv4FilterName)
				for name, props := range filteredData {
					aclMap[name] = props.(openapi.ConfigPutRequestIpv4FilterIpv4FilterName)
				}
				request.SetIpFilter(aclMap)
				return request
			}
		},

		ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
			switch operationType {
			case "PUT":
				req := b.client.ACLsAPI.AclsPut(ctx).IpVersion(ipVersion).AclsPutRequest(*request.(*openapi.AclsPutRequest))
				return req.Execute()
			case "PATCH":
				req := b.client.ACLsAPI.AclsPatch(ctx).IpVersion(ipVersion).AclsPutRequest(*request.(*openapi.AclsPutRequest))
				return req.Execute()
			case "DELETE":
				req := b.client.ACLsAPI.AclsDelete(ctx).IpFilterName(request.([]string)).IpVersion(ipVersion)
				return req.Execute()
			default:
				return nil, fmt.Errorf("unsupported ACL operation type: %s", operationType)
			}
		},

		ProcessResponse: nil,

		UpdateRecentOps: func() {
			// Already handled in parent function
		},
	})
}

func (b *BulkOperationManager) HasPendingOrRecentOperations(resourceType string) bool {
	return b.hasPendingOrRecentOperations(resourceType)
}

// ================================================================================================
// HELPER INFRASTRUCTURE - Generic data structure access and function factories
// ================================================================================================

// Generic data structure access methods

func (b *BulkOperationManager) storeOperation(resourceType, resourceName, operationType string, props interface{}) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	switch operationType {
	case "PUT", "PATCH":
		b.storeInTypedMap(resourceType, resourceName, operationType, props)
	case "DELETE":
		deleteSlice := b.getDeleteSlice(resourceType)
		if deleteSlice != nil {
			*deleteSlice = append(*deleteSlice, resourceName)
		}
	}
}

// storeInTypedMap stores operations directly in the appropriate typed map
func (b *BulkOperationManager) storeInTypedMap(resourceType, resourceName, operationType string, props interface{}) {
	switch resourceType {
	case "gateway":
		if operationType == "PUT" {
			if b.gatewayPut == nil {
				b.gatewayPut = make(map[string]openapi.ConfigPutRequestGatewayGatewayName)
			}
			b.gatewayPut[resourceName] = props.(openapi.ConfigPutRequestGatewayGatewayName)
		} else {
			if b.gatewayPatch == nil {
				b.gatewayPatch = make(map[string]openapi.ConfigPutRequestGatewayGatewayName)
			}
			b.gatewayPatch[resourceName] = props.(openapi.ConfigPutRequestGatewayGatewayName)
		}
	case "lag":
		if operationType == "PUT" {
			if b.lagPut == nil {
				b.lagPut = make(map[string]openapi.ConfigPutRequestLagLagName)
			}
			b.lagPut[resourceName] = props.(openapi.ConfigPutRequestLagLagName)
		} else {
			if b.lagPatch == nil {
				b.lagPatch = make(map[string]openapi.ConfigPutRequestLagLagName)
			}
			b.lagPatch[resourceName] = props.(openapi.ConfigPutRequestLagLagName)
		}
	case "tenant":
		if operationType == "PUT" {
			if b.tenantPut == nil {
				b.tenantPut = make(map[string]openapi.ConfigPutRequestTenantTenantName)
			}
			b.tenantPut[resourceName] = props.(openapi.ConfigPutRequestTenantTenantName)
		} else {
			if b.tenantPatch == nil {
				b.tenantPatch = make(map[string]openapi.ConfigPutRequestTenantTenantName)
			}
			b.tenantPatch[resourceName] = props.(openapi.ConfigPutRequestTenantTenantName)
		}
	case "service":
		if operationType == "PUT" {
			if b.servicePut == nil {
				b.servicePut = make(map[string]openapi.ConfigPutRequestServiceServiceName)
			}
			b.servicePut[resourceName] = props.(openapi.ConfigPutRequestServiceServiceName)
		} else {
			if b.servicePatch == nil {
				b.servicePatch = make(map[string]openapi.ConfigPutRequestServiceServiceName)
			}
			b.servicePatch[resourceName] = props.(openapi.ConfigPutRequestServiceServiceName)
		}
	case "gateway_profile":
		if operationType == "PUT" {
			if b.gatewayProfilePut == nil {
				b.gatewayProfilePut = make(map[string]openapi.ConfigPutRequestGatewayProfileGatewayProfileName)
			}
			b.gatewayProfilePut[resourceName] = props.(openapi.ConfigPutRequestGatewayProfileGatewayProfileName)
		} else {
			if b.gatewayProfilePatch == nil {
				b.gatewayProfilePatch = make(map[string]openapi.ConfigPutRequestGatewayProfileGatewayProfileName)
			}
			b.gatewayProfilePatch[resourceName] = props.(openapi.ConfigPutRequestGatewayProfileGatewayProfileName)
		}
	case "eth_port_profile":
		if operationType == "PUT" {
			if b.ethPortProfilePut == nil {
				b.ethPortProfilePut = make(map[string]openapi.ConfigPutRequestEthPortProfileEthPortProfileName)
			}
			b.ethPortProfilePut[resourceName] = props.(openapi.ConfigPutRequestEthPortProfileEthPortProfileName)
		} else {
			if b.ethPortProfilePatch == nil {
				b.ethPortProfilePatch = make(map[string]openapi.ConfigPutRequestEthPortProfileEthPortProfileName)
			}
			b.ethPortProfilePatch[resourceName] = props.(openapi.ConfigPutRequestEthPortProfileEthPortProfileName)
		}
	case "eth_port_settings":
		if operationType == "PUT" {
			if b.ethPortSettingsPut == nil {
				b.ethPortSettingsPut = make(map[string]openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName)
			}
			b.ethPortSettingsPut[resourceName] = props.(openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName)
		} else {
			if b.ethPortSettingsPatch == nil {
				b.ethPortSettingsPatch = make(map[string]openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName)
			}
			b.ethPortSettingsPatch[resourceName] = props.(openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName)
		}
	case "bundle":
		if operationType == "PUT" {
			if b.bundlePut == nil {
				b.bundlePut = make(map[string]openapi.BundlesPutRequestEndpointBundleValue)
			}
			b.bundlePut[resourceName] = props.(openapi.BundlesPutRequestEndpointBundleValue)
		} else {
			if b.bundlePatch == nil {
				b.bundlePatch = make(map[string]openapi.BundlesPatchRequestEndpointBundleValue)
			}
			b.bundlePatch[resourceName] = props.(openapi.BundlesPatchRequestEndpointBundleValue)
		}
	case "acl":
		if operationType == "PUT" {
			if b.aclPut == nil {
				b.aclPut = make(map[string]openapi.ConfigPutRequestIpv4FilterIpv4FilterName)
			}
			b.aclPut[resourceName] = props.(openapi.ConfigPutRequestIpv4FilterIpv4FilterName)
		} else {
			if b.aclPatch == nil {
				b.aclPatch = make(map[string]openapi.ConfigPutRequestIpv4FilterIpv4FilterName)
			}
			b.aclPatch[resourceName] = props.(openapi.ConfigPutRequestIpv4FilterIpv4FilterName)
		}
	case "authenticated_eth_port":
		if operationType == "PUT" {
			if b.authenticatedEthPortPut == nil {
				b.authenticatedEthPortPut = make(map[string]openapi.ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName)
			}
			b.authenticatedEthPortPut[resourceName] = props.(openapi.ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName)
		} else {
			if b.authenticatedEthPortPatch == nil {
				b.authenticatedEthPortPatch = make(map[string]openapi.ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName)
			}
			b.authenticatedEthPortPatch[resourceName] = props.(openapi.ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName)
		}
	case "badge":
		if operationType == "PUT" {
			if b.badgePut == nil {
				b.badgePut = make(map[string]openapi.ConfigPutRequestBadgeBadgeName)
			}
			b.badgePut[resourceName] = props.(openapi.ConfigPutRequestBadgeBadgeName)
		} else {
			if b.badgePatch == nil {
				b.badgePatch = make(map[string]openapi.ConfigPutRequestBadgeBadgeName)
			}
			b.badgePatch[resourceName] = props.(openapi.ConfigPutRequestBadgeBadgeName)
		}
	case "device_voice_settings":
		if operationType == "PUT" {
			if b.deviceVoiceSettingsPut == nil {
				b.deviceVoiceSettingsPut = make(map[string]openapi.ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName)
			}
			b.deviceVoiceSettingsPut[resourceName] = props.(openapi.ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName)
		} else {
			if b.deviceVoiceSettingsPatch == nil {
				b.deviceVoiceSettingsPatch = make(map[string]openapi.ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName)
			}
			b.deviceVoiceSettingsPatch[resourceName] = props.(openapi.ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName)
		}
	case "packet_broker":
		if operationType == "PUT" {
			if b.packetBrokerPut == nil {
				b.packetBrokerPut = make(map[string]openapi.ConfigPutRequestPbEgressProfilePbEgressProfileName)
			}
			b.packetBrokerPut[resourceName] = props.(openapi.ConfigPutRequestPbEgressProfilePbEgressProfileName)
		} else {
			if b.packetBrokerPatch == nil {
				b.packetBrokerPatch = make(map[string]openapi.ConfigPutRequestPbEgressProfilePbEgressProfileName)
			}
			b.packetBrokerPatch[resourceName] = props.(openapi.ConfigPutRequestPbEgressProfilePbEgressProfileName)
		}
	case "packet_queue":
		if operationType == "PUT" {
			if b.packetQueuePut == nil {
				b.packetQueuePut = make(map[string]openapi.ConfigPutRequestPacketQueuePacketQueueName)
			}
			b.packetQueuePut[resourceName] = props.(openapi.ConfigPutRequestPacketQueuePacketQueueName)
		} else {
			if b.packetQueuePatch == nil {
				b.packetQueuePatch = make(map[string]openapi.ConfigPutRequestPacketQueuePacketQueueName)
			}
			b.packetQueuePatch[resourceName] = props.(openapi.ConfigPutRequestPacketQueuePacketQueueName)
		}
	case "service_port_profile":
		if operationType == "PUT" {
			if b.servicePortProfilePut == nil {
				b.servicePortProfilePut = make(map[string]openapi.ConfigPutRequestServicePortProfileServicePortProfileName)
			}
			b.servicePortProfilePut[resourceName] = props.(openapi.ConfigPutRequestServicePortProfileServicePortProfileName)
		} else {
			if b.servicePortProfilePatch == nil {
				b.servicePortProfilePatch = make(map[string]openapi.ConfigPutRequestServicePortProfileServicePortProfileName)
			}
			b.servicePortProfilePatch[resourceName] = props.(openapi.ConfigPutRequestServicePortProfileServicePortProfileName)
		}
	case "switchpoint":
		if operationType == "PUT" {
			if b.switchpointPut == nil {
				b.switchpointPut = make(map[string]openapi.ConfigPutRequestSwitchpointSwitchpointName)
			}
			b.switchpointPut[resourceName] = props.(openapi.ConfigPutRequestSwitchpointSwitchpointName)
		} else {
			if b.switchpointPatch == nil {
				b.switchpointPatch = make(map[string]openapi.ConfigPutRequestSwitchpointSwitchpointName)
			}
			b.switchpointPatch[resourceName] = props.(openapi.ConfigPutRequestSwitchpointSwitchpointName)
		}
	case "device_controller":
		if operationType == "PUT" {
			if b.deviceControllerPut == nil {
				b.deviceControllerPut = make(map[string]openapi.ConfigPutRequestDeviceControllerDeviceControllerName)
			}
			b.deviceControllerPut[resourceName] = props.(openapi.ConfigPutRequestDeviceControllerDeviceControllerName)
		} else {
			if b.deviceControllerPatch == nil {
				b.deviceControllerPatch = make(map[string]openapi.ConfigPutRequestDeviceControllerDeviceControllerName)
			}
			b.deviceControllerPatch[resourceName] = props.(openapi.ConfigPutRequestDeviceControllerDeviceControllerName)
		}
	case "voice_port_profile":
		if operationType == "PUT" {
			if b.voicePortProfilePut == nil {
				b.voicePortProfilePut = make(map[string]openapi.ConfigPutRequestVoicePortProfilesVoicePortProfilesName)
			}
			b.voicePortProfilePut[resourceName] = props.(openapi.ConfigPutRequestVoicePortProfilesVoicePortProfilesName)
		} else {
			if b.voicePortProfilePatch == nil {
				b.voicePortProfilePatch = make(map[string]openapi.ConfigPutRequestVoicePortProfilesVoicePortProfilesName)
			}
			b.voicePortProfilePatch[resourceName] = props.(openapi.ConfigPutRequestVoicePortProfilesVoicePortProfilesName)
		}
	}
}

func (b *BulkOperationManager) getDeleteSlice(resourceType string) *[]string {
	switch resourceType {
	case "gateway":
		return &b.gatewayDelete
	case "lag":
		return &b.lagDelete
	case "tenant":
		return &b.tenantDelete
	case "service":
		return &b.serviceDelete
	case "gateway_profile":
		return &b.gatewayProfileDelete
	case "eth_port_profile":
		return &b.ethPortProfileDelete
	case "eth_port_settings":
		return &b.ethPortSettingsDelete
	case "bundle":
		return &b.bundleDelete
	case "acl":
		return &b.aclDelete
	case "authenticated_eth_port":
		return &b.authenticatedEthPortDelete
	case "badge":
		return &b.badgeDelete
	case "device_voice_settings":
		return &b.deviceVoiceSettingsDelete
	case "packet_broker":
		return &b.packetBrokerDelete
	case "packet_queue":
		return &b.packetQueueDelete
	case "service_port_profile":
		return &b.servicePortProfileDelete
	case "switchpoint":
		return &b.switchpointDelete
	case "voice_port_profile":
		return &b.voicePortProfileDelete
	case "device_controller":
		return &b.deviceControllerDelete
	}
	return nil
}

func (b *BulkOperationManager) getBatchSize(resourceType, operationType string) int {
	data := b.GetResourceOperationData(resourceType)
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
		return len(*data.DeleteOperations)
	}
	return 0
}

// Dynamic function factories that create the appropriate behavior for each resource
func (b *BulkOperationManager) createExtractor(resourceType, operationType string) func() (map[string]interface{}, []string) {
	var originalOperations map[string]interface{}

	return func() (map[string]interface{}, []string) {
		b.mutex.Lock()
		defer b.mutex.Unlock()

		switch operationType {
		case "PUT", "PATCH":
			operationMap := b.getOriginalOperationMap(resourceType, operationType)
			originalOperations = make(map[string]interface{})
			names := make([]string, 0, len(operationMap))

			for k, v := range operationMap {
				originalOperations[k] = v
				names = append(names, k)
			}

			b.clearOperationMap(resourceType, operationType)
			return originalOperations, names

		case "DELETE":
			deleteSlice := b.getDeleteSlice(resourceType)
			if deleteSlice == nil {
				return make(map[string]interface{}), []string{}
			}

			names := make([]string, len(*deleteSlice))
			copy(names, *deleteSlice)

			result := make(map[string]interface{})
			for _, name := range names {
				result[name] = true
			}

			*deleteSlice = (*deleteSlice)[:0]
			return result, names
		}

		return make(map[string]interface{}), []string{}
	}
}

func (b *BulkOperationManager) getOriginalOperationMap(resourceType, operationType string) map[string]interface{} {
	switch resourceType {
	case "gateway":
		if operationType == "PUT" {
			result := make(map[string]interface{})
			for k, v := range b.gatewayPut {
				result[k] = v
			}
			return result
		} else {
			result := make(map[string]interface{})
			for k, v := range b.gatewayPatch {
				result[k] = v
			}
			return result
		}
	case "lag":
		if operationType == "PUT" {
			result := make(map[string]interface{})
			for k, v := range b.lagPut {
				result[k] = v
			}
			return result
		} else {
			result := make(map[string]interface{})
			for k, v := range b.lagPatch {
				result[k] = v
			}
			return result
		}
	case "tenant":
		if operationType == "PUT" {
			result := make(map[string]interface{})
			for k, v := range b.tenantPut {
				result[k] = v
			}
			return result
		} else {
			result := make(map[string]interface{})
			for k, v := range b.tenantPatch {
				result[k] = v
			}
			return result
		}
	case "service":
		if operationType == "PUT" {
			result := make(map[string]interface{})
			for k, v := range b.servicePut {
				result[k] = v
			}
			return result
		} else {
			result := make(map[string]interface{})
			for k, v := range b.servicePatch {
				result[k] = v
			}
			return result
		}
	case "gateway_profile":
		if operationType == "PUT" {
			result := make(map[string]interface{})
			for k, v := range b.gatewayProfilePut {
				result[k] = v
			}
			return result
		} else {
			result := make(map[string]interface{})
			for k, v := range b.gatewayProfilePatch {
				result[k] = v
			}
			return result
		}
	case "eth_port_profile":
		if operationType == "PUT" {
			result := make(map[string]interface{})
			for k, v := range b.ethPortProfilePut {
				result[k] = v
			}
			return result
		} else {
			result := make(map[string]interface{})
			for k, v := range b.ethPortProfilePatch {
				result[k] = v
			}
			return result
		}
	case "eth_port_settings":
		if operationType == "PUT" {
			result := make(map[string]interface{})
			for k, v := range b.ethPortSettingsPut {
				result[k] = v
			}
			return result
		} else {
			result := make(map[string]interface{})
			for k, v := range b.ethPortSettingsPatch {
				result[k] = v
			}
			return result
		}
	case "bundle":
		if operationType == "PUT" {
			result := make(map[string]interface{})
			for k, v := range b.bundlePut {
				result[k] = v
			}
			return result
		} else {
			result := make(map[string]interface{})
			for k, v := range b.bundlePatch {
				result[k] = v
			}
			return result
		}
	case "acl":
		if operationType == "PUT" {
			result := make(map[string]interface{})
			for k, v := range b.aclPut {
				result[k] = v
			}
			return result
		} else {
			result := make(map[string]interface{})
			for k, v := range b.aclPatch {
				result[k] = v
			}
			return result
		}
	case "authenticated_eth_port":
		if operationType == "PUT" {
			result := make(map[string]interface{})
			for k, v := range b.authenticatedEthPortPut {
				result[k] = v
			}
			return result
		} else {
			result := make(map[string]interface{})
			for k, v := range b.authenticatedEthPortPatch {
				result[k] = v
			}
			return result
		}
	case "badge":
		if operationType == "PUT" {
			result := make(map[string]interface{})
			for k, v := range b.badgePut {
				result[k] = v
			}
			return result
		} else {
			result := make(map[string]interface{})
			for k, v := range b.badgePatch {
				result[k] = v
			}
			return result
		}
	case "device_voice_settings":
		if operationType == "PUT" {
			result := make(map[string]interface{})
			for k, v := range b.deviceVoiceSettingsPut {
				result[k] = v
			}
			return result
		} else {
			result := make(map[string]interface{})
			for k, v := range b.deviceVoiceSettingsPatch {
				result[k] = v
			}
			return result
		}
	case "packet_broker":
		if operationType == "PUT" {
			result := make(map[string]interface{})
			for k, v := range b.packetBrokerPut {
				result[k] = v
			}
			return result
		} else {
			result := make(map[string]interface{})
			for k, v := range b.packetBrokerPatch {
				result[k] = v
			}
			return result
		}
	case "packet_queue":
		if operationType == "PUT" {
			result := make(map[string]interface{})
			for k, v := range b.packetQueuePut {
				result[k] = v
			}
			return result
		} else {
			result := make(map[string]interface{})
			for k, v := range b.packetQueuePatch {
				result[k] = v
			}
			return result
		}
	case "service_port_profile":
		if operationType == "PUT" {
			result := make(map[string]interface{})
			for k, v := range b.servicePortProfilePut {
				result[k] = v
			}
			return result
		} else {
			result := make(map[string]interface{})
			for k, v := range b.servicePortProfilePatch {
				result[k] = v
			}
			return result
		}
	case "switchpoint":
		if operationType == "PUT" {
			result := make(map[string]interface{})
			for k, v := range b.switchpointPut {
				result[k] = v
			}
			return result
		} else {
			result := make(map[string]interface{})
			for k, v := range b.switchpointPatch {
				result[k] = v
			}
			return result
		}
	case "voice_port_profile":
		if operationType == "PUT" {
			result := make(map[string]interface{})
			for k, v := range b.voicePortProfilePut {
				result[k] = v
			}
			return result
		} else {
			result := make(map[string]interface{})
			for k, v := range b.voicePortProfilePatch {
				result[k] = v
			}
			return result
		}
	case "device_controller":
		if operationType == "PUT" {
			result := make(map[string]interface{})
			for k, v := range b.deviceControllerPut {
				result[k] = v
			}
			return result
		} else {
			result := make(map[string]interface{})
			for k, v := range b.deviceControllerPatch {
				result[k] = v
			}
			return result
		}
	}
	return make(map[string]interface{})
}

func (b *BulkOperationManager) clearOperationMap(resourceType, operationType string) {
	switch resourceType {
	case "gateway":
		if operationType == "PUT" {
			b.gatewayPut = make(map[string]openapi.ConfigPutRequestGatewayGatewayName)
		} else {
			b.gatewayPatch = make(map[string]openapi.ConfigPutRequestGatewayGatewayName)
		}
	case "lag":
		if operationType == "PUT" {
			b.lagPut = make(map[string]openapi.ConfigPutRequestLagLagName)
		} else {
			b.lagPatch = make(map[string]openapi.ConfigPutRequestLagLagName)
		}
	case "tenant":
		if operationType == "PUT" {
			b.tenantPut = make(map[string]openapi.ConfigPutRequestTenantTenantName)
		} else {
			b.tenantPatch = make(map[string]openapi.ConfigPutRequestTenantTenantName)
		}
	case "service":
		if operationType == "PUT" {
			b.servicePut = make(map[string]openapi.ConfigPutRequestServiceServiceName)
		} else {
			b.servicePatch = make(map[string]openapi.ConfigPutRequestServiceServiceName)
		}
	case "gateway_profile":
		if operationType == "PUT" {
			b.gatewayProfilePut = make(map[string]openapi.ConfigPutRequestGatewayProfileGatewayProfileName)
		} else {
			b.gatewayProfilePatch = make(map[string]openapi.ConfigPutRequestGatewayProfileGatewayProfileName)
		}
	case "eth_port_profile":
		if operationType == "PUT" {
			b.ethPortProfilePut = make(map[string]openapi.ConfigPutRequestEthPortProfileEthPortProfileName)
		} else {
			b.ethPortProfilePatch = make(map[string]openapi.ConfigPutRequestEthPortProfileEthPortProfileName)
		}
	case "eth_port_settings":
		if operationType == "PUT" {
			b.ethPortSettingsPut = make(map[string]openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName)
		} else {
			b.ethPortSettingsPatch = make(map[string]openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName)
		}
	case "bundle":
		if operationType == "PUT" {
			b.bundlePut = make(map[string]openapi.BundlesPutRequestEndpointBundleValue)
		} else {
			b.bundlePatch = make(map[string]openapi.BundlesPatchRequestEndpointBundleValue)
		}
	case "acl":
		if operationType == "PUT" {
			b.aclPut = make(map[string]openapi.ConfigPutRequestIpv4FilterIpv4FilterName)
		} else {
			b.aclPatch = make(map[string]openapi.ConfigPutRequestIpv4FilterIpv4FilterName)
		}
	case "authenticated_eth_port":
		if operationType == "PUT" {
			b.authenticatedEthPortPut = make(map[string]openapi.ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName)
		} else {
			b.authenticatedEthPortPatch = make(map[string]openapi.ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName)
		}
	case "badge":
		if operationType == "PUT" {
			b.badgePut = make(map[string]openapi.ConfigPutRequestBadgeBadgeName)
		} else {
			b.badgePatch = make(map[string]openapi.ConfigPutRequestBadgeBadgeName)
		}
	case "device_voice_settings":
		if operationType == "PUT" {
			b.deviceVoiceSettingsPut = make(map[string]openapi.ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName)
		} else {
			b.deviceVoiceSettingsPatch = make(map[string]openapi.ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName)
		}
	case "packet_broker":
		if operationType == "PUT" {
			b.packetBrokerPut = make(map[string]openapi.ConfigPutRequestPbEgressProfilePbEgressProfileName)
		} else {
			b.packetBrokerPatch = make(map[string]openapi.ConfigPutRequestPbEgressProfilePbEgressProfileName)
		}
	case "packet_queue":
		if operationType == "PUT" {
			b.packetQueuePut = make(map[string]openapi.ConfigPutRequestPacketQueuePacketQueueName)
		} else {
			b.packetQueuePatch = make(map[string]openapi.ConfigPutRequestPacketQueuePacketQueueName)
		}
	case "service_port_profile":
		if operationType == "PUT" {
			b.servicePortProfilePut = make(map[string]openapi.ConfigPutRequestServicePortProfileServicePortProfileName)
		} else {
			b.servicePortProfilePatch = make(map[string]openapi.ConfigPutRequestServicePortProfileServicePortProfileName)
		}
	case "switchpoint":
		if operationType == "PUT" {
			b.switchpointPut = make(map[string]openapi.ConfigPutRequestSwitchpointSwitchpointName)
		} else {
			b.switchpointPatch = make(map[string]openapi.ConfigPutRequestSwitchpointSwitchpointName)
		}
	case "voice_port_profile":
		if operationType == "PUT" {
			b.voicePortProfilePut = make(map[string]openapi.ConfigPutRequestVoicePortProfilesVoicePortProfilesName)
		} else {
			b.voicePortProfilePatch = make(map[string]openapi.ConfigPutRequestVoicePortProfilesVoicePortProfilesName)
		}
	case "device_controller":
		if operationType == "PUT" {
			b.deviceControllerPut = make(map[string]openapi.ConfigPutRequestDeviceControllerDeviceControllerName)
		} else {
			b.deviceControllerPatch = make(map[string]openapi.ConfigPutRequestDeviceControllerDeviceControllerName)
		}
	}
}

func (b *BulkOperationManager) createPreExistenceChecker(config ResourceConfig, operationType string) func(context.Context, []string, map[string]interface{}) ([]string, map[string]interface{}, error) {
	if operationType != "PUT" {
		return nil // Only PUT operations need pre-existence checking
	}

	return func(ctx context.Context, resourceNames []string, originalOperations map[string]interface{}) ([]string, map[string]interface{}, error) {
		checker := ResourceExistenceCheck{
			ResourceType:  config.ResourceType,
			OperationType: "PUT",
			FetchResources: func(ctx context.Context) (map[string]interface{}, error) {
				apiCtx, cancel := context.WithTimeout(ctx, OperationTimeout)
				defer cancel()

				switch config.ResourceType {
				case "gateway":
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

				case "lag":
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
					return result.Lag, nil

				case "tenant":
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
					return result.Tenant, nil

				case "service":
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
					return result.Service, nil

				case "gateway_profile":
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
					return result.GatewayProfile, nil

				case "eth_port_profile":
					resp, err := b.client.EthPortProfilesAPI.EthportprofilesGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						EthPortProfile map[string]interface{} `json:"eth_port_profile_"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.EthPortProfile, nil

				case "eth_port_settings":
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
					return result.EthPortSettings, nil

				case "bundle":
					resp, err := b.client.BundlesAPI.BundlesGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						Bundle map[string]interface{} `json:"endpoint_bundle"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.Bundle, nil

				case "acl":
					// ACLs require IP version checking - need to check both IPv4 and IPv6

					// Try IPv4
					resp4, err4 := b.client.ACLsAPI.AclsGet(apiCtx).IpVersion("4").Execute()
					existingResources := make(map[string]interface{})

					if err4 == nil {
						defer resp4.Body.Close()
						var result4 struct {
							IpFilter map[string]interface{} `json:"ip_filter"`
						}
						if err := json.NewDecoder(resp4.Body).Decode(&result4); err == nil {
							for name, props := range result4.IpFilter {
								existingResources[name] = props
							}
						}
					}

					// Try IPv6
					resp6, err6 := b.client.ACLsAPI.AclsGet(apiCtx).IpVersion("6").Execute()
					if err6 == nil {
						defer resp6.Body.Close()
						var result6 struct {
							IpFilter map[string]interface{} `json:"ip_filter"`
						}
						if err := json.NewDecoder(resp6.Body).Decode(&result6); err == nil {
							for name, props := range result6.IpFilter {
								existingResources[name] = props
							}
						}
					}

					if err4 != nil && err6 != nil {
						return nil, err4
					}

					return existingResources, nil

				case "authenticated_eth_port":
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
					return result.AuthenticatedEthPort, nil

				case "badge":
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
					return result.Badge, nil

				case "device_voice_settings":
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
					return result.DeviceVoiceSettings, nil

				case "packet_broker":
					resp, err := b.client.PacketBrokerAPI.PacketbrokerGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						PacketBroker map[string]interface{} `json:"pb_egress_profile"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.PacketBroker, nil

				case "packet_queue":
					resp, err := b.client.PacketQueuesAPI.PacketqueuesGet(apiCtx).Execute()
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
					return result.PacketQueue, nil

				case "service_port_profile":
					resp, err := b.client.ServicePortProfilesAPI.ServiceportprofilesGet(apiCtx).Execute()
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

				case "switchpoint":
					resp, err := b.client.SwitchpointsAPI.SwitchpointsGet(apiCtx).Execute()
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
					return result.Switchpoint, nil

				case "voice_port_profile":
					resp, err := b.client.VoicePortProfilesAPI.VoiceportprofilesGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						VoicePortProfile map[string]interface{} `json:"voice_port_profiles"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.VoicePortProfile, nil

				case "device_controller":
					resp, err := b.client.DeviceControllersAPI.DevicecontrollersGet(apiCtx).Execute()
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
					return result.DeviceController, nil

				default:
					// For unknown resource types, assume no existing resources to avoid errors
					return make(map[string]interface{}), nil
				}
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
	}
}

func (b *BulkOperationManager) createRequestPreparer(config ResourceConfig, operationType string) func(map[string]interface{}) interface{} {
	return func(filteredData map[string]interface{}) interface{} {
		// Handle DELETE operations
		if operationType == "DELETE" {
			names := make([]string, 0, len(filteredData))
			for name := range filteredData {
				names = append(names, name)
			}
			return names
		}

		// Handle PUT and PATCH operations
		switch config.ResourceType {
		case "gateway":
			putRequest := openapi.NewGatewaysPutRequest()
			gatewayMap := make(map[string]openapi.ConfigPutRequestGatewayGatewayName)
			for name, props := range filteredData {
				gatewayMap[name] = props.(openapi.ConfigPutRequestGatewayGatewayName)
			}
			putRequest.SetGateway(gatewayMap)
			return putRequest
		case "lag":
			putRequest := openapi.NewLagsPutRequest()
			lagMap := make(map[string]openapi.ConfigPutRequestLagLagName)
			for name, props := range filteredData {
				lagMap[name] = props.(openapi.ConfigPutRequestLagLagName)
			}
			putRequest.SetLag(lagMap)
			return putRequest
		case "tenant":
			putRequest := openapi.NewTenantsPutRequest()
			tenantMap := make(map[string]openapi.ConfigPutRequestTenantTenantName)
			for name, props := range filteredData {
				tenantMap[name] = props.(openapi.ConfigPutRequestTenantTenantName)
			}
			putRequest.SetTenant(tenantMap)
			return putRequest
		case "service":
			putRequest := openapi.NewServicesPutRequest()
			serviceMap := make(map[string]openapi.ConfigPutRequestServiceServiceName)
			for name, props := range filteredData {
				serviceMap[name] = props.(openapi.ConfigPutRequestServiceServiceName)
			}
			putRequest.SetService(serviceMap)
			return putRequest
		case "gateway_profile":
			putRequest := openapi.NewGatewayprofilesPutRequest()
			profileMap := make(map[string]openapi.ConfigPutRequestGatewayProfileGatewayProfileName)
			for name, props := range filteredData {
				profileMap[name] = props.(openapi.ConfigPutRequestGatewayProfileGatewayProfileName)
			}
			putRequest.SetGatewayProfile(profileMap)
			return putRequest
		case "eth_port_profile":
			putRequest := openapi.NewEthportprofilesPutRequest()
			profileMap := make(map[string]openapi.ConfigPutRequestEthPortProfileEthPortProfileName)
			for name, props := range filteredData {
				profileMap[name] = props.(openapi.ConfigPutRequestEthPortProfileEthPortProfileName)
			}
			putRequest.SetEthPortProfile(profileMap)
			return putRequest
		case "eth_port_settings":
			putRequest := openapi.NewEthportsettingsPutRequest()
			settingsMap := make(map[string]openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName)
			for name, props := range filteredData {
				settingsMap[name] = props.(openapi.ConfigPutRequestEthPortSettingsEthPortSettingsName)
			}
			putRequest.SetEthPortSettings(settingsMap)
			return putRequest
		case "bundle":
			if operationType == "PUT" {
				putRequest := openapi.NewBundlesPutRequest()
				bundleMap := make(map[string]openapi.BundlesPutRequestEndpointBundleValue)
				for name, props := range filteredData {
					bundleMap[name] = props.(openapi.BundlesPutRequestEndpointBundleValue)
				}
				putRequest.SetEndpointBundle(bundleMap)
				return putRequest
			} else {
				patchRequest := openapi.NewBundlesPatchRequest()
				bundleMap := make(map[string]openapi.BundlesPatchRequestEndpointBundleValue)
				for name, props := range filteredData {
					bundleMap[name] = props.(openapi.BundlesPatchRequestEndpointBundleValue)
				}
				patchRequest.SetEndpointBundle(bundleMap)
				return patchRequest
			}
		case "acl":
			putRequest := openapi.NewAclsPutRequest()
			aclMap := make(map[string]openapi.ConfigPutRequestIpv4FilterIpv4FilterName)
			for name, props := range filteredData {
				aclMap[name] = props.(openapi.ConfigPutRequestIpv4FilterIpv4FilterName)
			}
			putRequest.SetIpFilter(aclMap)
			return putRequest
		case "authenticated_eth_port":
			putRequest := openapi.NewAuthenticatedethportsPutRequest()
			portMap := make(map[string]openapi.ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName)
			for name, props := range filteredData {
				portMap[name] = props.(openapi.ConfigPutRequestAuthenticatedEthPortAuthenticatedEthPortName)
			}
			putRequest.SetAuthenticatedEthPort(portMap)
			return putRequest
		case "badge":
			putRequest := openapi.NewBadgesPutRequest()
			badgeMap := make(map[string]openapi.ConfigPutRequestBadgeBadgeName)
			for name, props := range filteredData {
				badgeMap[name] = props.(openapi.ConfigPutRequestBadgeBadgeName)
			}
			putRequest.SetBadge(badgeMap)
			return putRequest
		case "device_voice_settings":
			putRequest := openapi.NewDevicevoicesettingsPutRequest()
			settingsMap := make(map[string]openapi.ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName)
			for name, props := range filteredData {
				settingsMap[name] = props.(openapi.ConfigPutRequestDeviceVoiceSettingsDeviceVoiceSettingsName)
			}
			putRequest.SetDeviceVoiceSettings(settingsMap)
			return putRequest
		case "packet_broker":
			putRequest := openapi.NewPacketbrokerPutRequest()
			brokerMap := make(map[string]openapi.ConfigPutRequestPbEgressProfilePbEgressProfileName)
			for name, props := range filteredData {
				brokerMap[name] = props.(openapi.ConfigPutRequestPbEgressProfilePbEgressProfileName)
			}
			putRequest.SetPbEgressProfile(brokerMap)
			return putRequest
		case "packet_queue":
			putRequest := openapi.NewPacketqueuesPutRequest()
			queueMap := make(map[string]openapi.ConfigPutRequestPacketQueuePacketQueueName)
			for name, props := range filteredData {
				queueMap[name] = props.(openapi.ConfigPutRequestPacketQueuePacketQueueName)
			}
			putRequest.SetPacketQueue(queueMap)
			return putRequest
		case "service_port_profile":
			putRequest := openapi.NewServiceportprofilesPutRequest()
			profileMap := make(map[string]openapi.ConfigPutRequestServicePortProfileServicePortProfileName)
			for name, props := range filteredData {
				profileMap[name] = props.(openapi.ConfigPutRequestServicePortProfileServicePortProfileName)
			}
			putRequest.SetServicePortProfile(profileMap)
			return putRequest
		case "switchpoint":
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

				switchpointMap[name] = switchpointValue
			}
			putRequest.SetSwitchpoint(switchpointMap)
			return putRequest
		case "voice_port_profile":
			putRequest := openapi.NewVoiceportprofilesPutRequest()
			profileMap := make(map[string]openapi.ConfigPutRequestVoicePortProfilesVoicePortProfilesName)
			for name, props := range filteredData {
				profileMap[name] = props.(openapi.ConfigPutRequestVoicePortProfilesVoicePortProfilesName)
			}
			putRequest.SetVoicePortProfiles(profileMap)
			return putRequest
		case "device_controller":
			putRequest := openapi.NewDevicecontrollersPutRequest()
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
			putRequest.SetDeviceController(deviceMap)
			return putRequest
		}
		return nil
	}
}

func (b *BulkOperationManager) createRequestExecutor(config ResourceConfig, operationType string) func(context.Context, interface{}) (*http.Response, error) {
	return func(ctx context.Context, request interface{}) (*http.Response, error) {
		apiClient := config.APIClientGetter(b.client)

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

func (b *BulkOperationManager) createResponseProcessor(config ResourceConfig, operationType string) func(context.Context, *http.Response) error {
	if !config.HasAutoGen {
		return nil // No post-processing needed for resources without auto-generated fields
	}

	return func(ctx context.Context, resp *http.Response) error {
		delayTime := 2 * time.Second
		tflog.Debug(ctx, fmt.Sprintf("Waiting %v for auto-generated values to be assigned before fetching %s", delayTime, config.ResourceType))
		time.Sleep(delayTime)

		fetchCtx, fetchCancel := context.WithTimeout(context.Background(), OperationTimeout)
		defer fetchCancel()

		tflog.Debug(ctx, fmt.Sprintf("Fetching %s after successful PUT operation to retrieve auto-generated values", config.ResourceType))

		switch config.ResourceType {
		case "tenant":
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

		case "service":
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

		case "switchpoint":
			switchpointsReq := b.client.SwitchpointsAPI.SwitchpointsGet(fetchCtx)
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
				tflog.Error(ctx, "Failed to decode switchpoints response for auto-generated fields", map[string]interface{}{
					"error": respErr.Error(),
				})
				return respErr
			}

			b.switchpointResponsesMutex.Lock()
			for switchpointName, switchpointData := range switchpointsData.Switchpoint {
				b.switchpointResponses[switchpointName] = switchpointData

				if name, ok := switchpointData["name"].(string); ok && name != switchpointName {
					b.switchpointResponses[name] = switchpointData
				}
			}
			b.switchpointResponsesMutex.Unlock()

			tflog.Debug(ctx, "Successfully stored switchpoint data for auto-generated fields", map[string]interface{}{
				"switchpoint_count": len(switchpointsData.Switchpoint),
			})

		default:
			tflog.Warn(ctx, fmt.Sprintf("Unknown resource type with auto-generated fields: %s", config.ResourceType))
		}

		return nil
	}
}

func (b *BulkOperationManager) createRecentOpsUpdater(resourceType string) func() {
	return func() {
		now := time.Now()
		switch resourceType {
		case "gateway":
			b.recentGatewayOps = true
			b.recentGatewayOpTime = now
		case "lag":
			b.recentLagOps = true
			b.recentLagOpTime = now
		case "tenant":
			b.recentTenantOps = true
			b.recentTenantOpTime = now
		case "service":
			b.recentServiceOps = true
			b.recentServiceOpTime = now
		case "gateway_profile":
			b.recentGatewayProfileOps = true
			b.recentGatewayProfileOpTime = now
		case "eth_port_profile":
			b.recentEthPortProfileOps = true
			b.recentEthPortProfileOpTime = now
		case "eth_port_settings":
			b.recentEthPortSettingsOps = true
			b.recentEthPortSettingsOpTime = now
		case "bundle":
			b.recentBundleOps = true
			b.recentBundleOpTime = now
		case "acl":
			b.recentAclOps = true
			b.recentAclOpTime = now
		case "authenticated_eth_port":
			b.recentAuthenticatedEthPortOps = true
			b.recentAuthenticatedEthPortOpTime = now
		case "badge":
			b.recentBadgeOps = true
			b.recentBadgeOpTime = now
		case "device_voice_settings":
			b.recentDeviceVoiceSettingsOps = true
			b.recentDeviceVoiceSettingsOpTime = now
		case "packet_broker":
			b.recentPacketBrokerOps = true
			b.recentPacketBrokerOpTime = now
		case "packet_queue":
			b.recentPacketQueueOps = true
			b.recentPacketQueueOpTime = now
		case "service_port_profile":
			b.recentServicePortProfileOps = true
			b.recentServicePortProfileOpTime = now
		case "switchpoint":
			b.recentSwitchpointOps = true
			b.recentSwitchpointOpTime = now
		case "voice_port_profile":
			b.recentVoicePortProfileOps = true
			b.recentVoicePortProfileOpTime = now
		case "device_controller":
			b.recentDeviceControllerOps = true
			b.recentDeviceControllerOpTime = now
		}
	}
}

// updateOperationStatuses updates the status of pending operations based on the bulk operation result
func (b *BulkOperationManager) updateOperationStatuses(ctx context.Context, resourceType, operationType string, resourceNames []string, opErr error) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	resourceMap := make(map[string]bool)
	for _, name := range resourceNames {
		resourceMap[name] = true
	}

	for opID, op := range b.pendingOperations {
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
