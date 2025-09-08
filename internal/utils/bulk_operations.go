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
	gatewayPut    map[string]openapi.GatewaysPutRequestGatewayValue
	gatewayPatch  map[string]openapi.GatewaysPutRequestGatewayValue
	gatewayDelete []string

	// LAG operations
	lagPut    map[string]openapi.LagsPutRequestLagValue
	lagPatch  map[string]openapi.LagsPutRequestLagValue
	lagDelete []string

	// Tenant operations
	tenantPut    map[string]openapi.TenantsPutRequestTenantValue
	tenantPatch  map[string]openapi.TenantsPutRequestTenantValue
	tenantDelete []string

	// Service operations
	servicePut    map[string]openapi.ServicesPutRequestServiceValue
	servicePatch  map[string]openapi.ServicesPutRequestServiceValue
	serviceDelete []string

	// Gateway Profile operations
	gatewayProfilePut    map[string]openapi.GatewayprofilesPutRequestGatewayProfileValue
	gatewayProfilePatch  map[string]openapi.GatewayprofilesPutRequestGatewayProfileValue
	gatewayProfileDelete []string

	// EthPortProfile operations
	ethPortProfilePut    map[string]openapi.EthportprofilesPutRequestEthPortProfileValue
	ethPortProfilePatch  map[string]openapi.EthportprofilesPutRequestEthPortProfileValue
	ethPortProfileDelete []string

	// EthPortSettings operations
	ethPortSettingsPut    map[string]openapi.EthportsettingsPutRequestEthPortSettingsValue
	ethPortSettingsPatch  map[string]openapi.EthportsettingsPutRequestEthPortSettingsValue
	ethPortSettingsDelete []string

	// Bundles operations
	bundlePut    map[string]openapi.BundlesPutRequestEndpointBundleValue
	bundlePatch  map[string]openapi.BundlesPutRequestEndpointBundleValue
	bundleDelete []string

	// ACL operations (unified for IPv4 and IPv6)
	aclPut       map[string]openapi.AclsPutRequestIpFilterValue
	aclPatch     map[string]openapi.AclsPutRequestIpFilterValue
	aclDelete    []string
	aclIpVersion map[string]string // Track which IP version each ACL operation uses

	// Authenticated Eth-Port operations
	authenticatedEthPortPut    map[string]openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValue
	authenticatedEthPortPatch  map[string]openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValue
	authenticatedEthPortDelete []string

	// Badge operations
	badgePut    map[string]openapi.BadgesPutRequestBadgeValue
	badgePatch  map[string]openapi.BadgesPutRequestBadgeValue
	badgeDelete []string

	// Device Port Settings operations
	deviceVoiceSettingsPut    map[string]openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValue
	deviceVoiceSettingsPatch  map[string]openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValue
	deviceVoiceSettingsDelete []string

	// Packet Broker operations
	packetBrokerPut    map[string]openapi.PacketbrokerPutRequestPbEgressProfileValue
	packetBrokerPatch  map[string]openapi.PacketbrokerPutRequestPbEgressProfileValue
	packetBrokerDelete []string

	// Packet Queues operations
	packetQueuePut    map[string]openapi.PacketqueuesPutRequestPacketQueueValue
	packetQueuePatch  map[string]openapi.PacketqueuesPutRequestPacketQueueValue
	packetQueueDelete []string

	// ServicePort Profile operations
	servicePortProfilePut    map[string]openapi.ServiceportprofilesPutRequestServicePortProfileValue
	servicePortProfilePatch  map[string]openapi.ServiceportprofilesPutRequestServicePortProfileValue
	servicePortProfileDelete []string

	// Switchpoint operations
	switchpointPut    map[string]openapi.SwitchpointsPutRequestSwitchpointValue
	switchpointPatch  map[string]openapi.SwitchpointsPutRequestSwitchpointValue
	switchpointDelete []string

	// Voice Port Profile operations
	voicePortProfilePut    map[string]openapi.VoiceportprofilesPutRequestVoicePortProfilesValue
	voicePortProfilePatch  map[string]openapi.VoiceportprofilesPutRequestVoicePortProfilesValue
	voicePortProfileDelete []string

	// Device Controller operations
	deviceControllerPut    map[string]openapi.DevicecontrollersPutRequestDeviceControllerValue
	deviceControllerPatch  map[string]openapi.DevicecontrollersPutRequestDeviceControllerValue
	deviceControllerDelete []string

	// AS Path Access List operations
	asPathAccessListPut    map[string]openapi.AspathaccesslistsPutRequestAsPathAccessListValue
	asPathAccessListPatch  map[string]openapi.AspathaccesslistsPutRequestAsPathAccessListValue
	asPathAccessListDelete []string

	// Community List operations
	communityListPut    map[string]openapi.CommunitylistsPutRequestCommunityListValue
	communityListPatch  map[string]openapi.CommunitylistsPutRequestCommunityListValue
	communityListDelete []string

	// Device Settings operations
	deviceSettingsPut    map[string]openapi.DevicesettingsPutRequestEthDeviceProfilesValue
	deviceSettingsPatch  map[string]openapi.DevicesettingsPutRequestEthDeviceProfilesValue
	deviceSettingsDelete []string

	// Extended Community List operations
	extendedCommunityListPut    map[string]openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValue
	extendedCommunityListPatch  map[string]openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValue
	extendedCommunityListDelete []string

	// IPv4 List operations
	ipv4ListPut    map[string]openapi.Ipv4listsPutRequestIpv4ListFilterValue
	ipv4ListPatch  map[string]openapi.Ipv4listsPutRequestIpv4ListFilterValue
	ipv4ListDelete []string

	// IPv4 Prefix List operations
	ipv4PrefixListPut    map[string]openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValue
	ipv4PrefixListPatch  map[string]openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValue
	ipv4PrefixListDelete []string

	// IPv6 List operations
	ipv6ListPut    map[string]openapi.Ipv6listsPutRequestIpv6ListFilterValue
	ipv6ListPatch  map[string]openapi.Ipv6listsPutRequestIpv6ListFilterValue
	ipv6ListDelete []string

	// IPv6 Prefix List operations
	ipv6PrefixListPut    map[string]openapi.Ipv6prefixlistsPutRequestIpv6PrefixListValue
	ipv6PrefixListPatch  map[string]openapi.Ipv6prefixlistsPutRequestIpv6PrefixListValue
	ipv6PrefixListDelete []string

	// Route Map Clause operations
	routeMapClausePut    map[string]openapi.RoutemapclausesPutRequestRouteMapClauseValue
	routeMapClausePatch  map[string]openapi.RoutemapclausesPutRequestRouteMapClauseValue
	routeMapClauseDelete []string

	// Route Map operations
	routeMapPut    map[string]openapi.RoutemapsPutRequestRouteMapValue
	routeMapPatch  map[string]openapi.RoutemapsPutRequestRouteMapValue
	routeMapDelete []string

	// SFP Breakout operations (only PATCH available)
	sfpBreakoutPatch map[string]openapi.SfpbreakoutsPatchRequestSfpBreakoutsValue

	// Site operations (only PATCH available)
	sitePatch map[string]openapi.SitesPatchRequestSiteValue

	// Pod operations
	podPut    map[string]openapi.PodsPutRequestPodValue
	podPatch  map[string]openapi.PodsPutRequestPodValue
	podDelete []string

	// Port ACL operations
	portAclPut    map[string]openapi.PortaclsPutRequestPortAclValue
	portAclPatch  map[string]openapi.PortaclsPutRequestPortAclValue
	portAclDelete []string

	// SFlow Collector operations
	sflowCollectorPut    map[string]openapi.SflowcollectorsPutRequestSflowCollectorValue
	sflowCollectorPatch  map[string]openapi.SflowcollectorsPutRequestSflowCollectorValue
	sflowCollectorDelete []string

	// Diagnostics Profile operations
	diagnosticsProfilePut    map[string]openapi.DiagnosticsprofilesPutRequestDiagnosticsProfileValue
	diagnosticsProfilePatch  map[string]openapi.DiagnosticsprofilesPutRequestDiagnosticsProfileValue
	diagnosticsProfileDelete []string

	// Diagnostics Port Profile operations
	diagnosticsPortProfilePut    map[string]openapi.DiagnosticsportprofilesPutRequestDiagnosticsPortProfileValue
	diagnosticsPortProfilePatch  map[string]openapi.DiagnosticsportprofilesPutRequestDiagnosticsPortProfileValue
	diagnosticsPortProfileDelete []string

	// Track recent operations to avoid race conditions
	recentGatewayOps                   bool
	recentGatewayOpTime                time.Time
	recentLagOps                       bool
	recentLagOpTime                    time.Time
	recentServiceOps                   bool
	recentServiceOpTime                time.Time
	recentTenantOps                    bool
	recentTenantOpTime                 time.Time
	recentGatewayProfileOps            bool
	recentGatewayProfileOpTime         time.Time
	recentEthPortProfileOps            bool
	recentEthPortProfileOpTime         time.Time
	recentEthPortSettingsOps           bool
	recentEthPortSettingsOpTime        time.Time
	recentBundleOps                    bool
	recentBundleOpTime                 time.Time
	recentAclOps                       bool
	recentAclOpTime                    time.Time
	recentAuthenticatedEthPortOps      bool
	recentAuthenticatedEthPortOpTime   time.Time
	recentBadgeOps                     bool
	recentBadgeOpTime                  time.Time
	recentVoicePortProfileOps          bool
	recentVoicePortProfileOpTime       time.Time
	recentSwitchpointOps               bool
	recentSwitchpointOpTime            time.Time
	recentServicePortProfileOps        bool
	recentServicePortProfileOpTime     time.Time
	recentPacketBrokerOps              bool
	recentPacketBrokerOpTime           time.Time
	recentPacketQueueOps               bool
	recentPacketQueueOpTime            time.Time
	recentDeviceVoiceSettingsOps       bool
	recentDeviceVoiceSettingsOpTime    time.Time
	recentDeviceControllerOps          bool
	recentDeviceControllerOpTime       time.Time
	recentAsPathAccessListOps          bool
	recentAsPathAccessListOpTime       time.Time
	recentCommunityListOps             bool
	recentCommunityListOpTime          time.Time
	recentDeviceSettingsOps            bool
	recentDeviceSettingsOpTime         time.Time
	recentExtendedCommunityListOps     bool
	recentExtendedCommunityListOpTime  time.Time
	recentIpv4ListOps                  bool
	recentIpv4ListOpTime               time.Time
	recentIpv4PrefixListOps            bool
	recentIpv4PrefixListOpTime         time.Time
	recentIpv6ListOps                  bool
	recentIpv6ListOpTime               time.Time
	recentIpv6PrefixListOps            bool
	recentIpv6PrefixListOpTime         time.Time
	recentRouteMapClauseOps            bool
	recentRouteMapClauseOpTime         time.Time
	recentRouteMapOps                  bool
	recentRouteMapOpTime               time.Time
	recentSfpBreakoutOps               bool
	recentSfpBreakoutOpTime            time.Time
	recentSiteOps                      bool
	recentSiteOpTime                   time.Time
	recentPodOps                       bool
	recentPodOpTime                    time.Time
	recentPortAclOps                   bool
	recentPortAclOpTime                time.Time
	recentSflowCollectorOps            bool
	recentSflowCollectorOpTime         time.Time
	recentDiagnosticsProfileOps        bool
	recentDiagnosticsProfileOpTime     time.Time
	recentDiagnosticsPortProfileOps    bool
	recentDiagnosticsPortProfileOpTime time.Time

	// For tracking operations
	pendingOperations     map[string]*Operation
	operationResults      map[string]bool // true = success, false = failure
	operationErrors       map[string]error
	operationWaitChannels map[string]chan struct{}
	operationMutex        sync.Mutex
	closedChannels        map[string]bool

	// Store API responses=
	gatewayResponses                     map[string]map[string]interface{}
	gatewayResponsesMutex                sync.RWMutex
	lagResponses                         map[string]map[string]interface{}
	lagResponsesMutex                    sync.RWMutex
	serviceResponses                     map[string]map[string]interface{}
	serviceResponsesMutex                sync.RWMutex
	tenantResponses                      map[string]map[string]interface{}
	tenantResponsesMutex                 sync.RWMutex
	gatewayProfileResponses              map[string]map[string]interface{}
	gatewayProfileResponsesMutex         sync.RWMutex
	ethPortProfileResponses              map[string]map[string]interface{}
	ethPortProfileResponsesMutex         sync.RWMutex
	ethPortSettingsResponses             map[string]map[string]interface{}
	ethPortSettingsResponsesMutex        sync.RWMutex
	bundleResponses                      map[string]map[string]interface{}
	bundleResponsesMutex                 sync.RWMutex
	aclResponses                         map[string]map[string]interface{}
	aclResponsesMutex                    sync.RWMutex
	authenticatedEthPortResponses        map[string]map[string]interface{}
	authenticatedEthPortResponsesMutex   sync.RWMutex
	badgeResponses                       map[string]map[string]interface{}
	badgeResponsesMutex                  sync.RWMutex
	voicePortProfileResponses            map[string]map[string]interface{}
	voicePortProfileResponsesMutex       sync.RWMutex
	switchpointResponses                 map[string]map[string]interface{}
	switchpointResponsesMutex            sync.RWMutex
	servicePortProfileResponses          map[string]map[string]interface{}
	servicePortProfileResponsesMutex     sync.RWMutex
	packetBrokerResponses                map[string]map[string]interface{}
	packetBrokerResponsesMutex           sync.RWMutex
	packetQueueResponses                 map[string]map[string]interface{}
	packetQueueResponsesMutex            sync.RWMutex
	deviceVoiceSettingsResponses         map[string]map[string]interface{}
	deviceVoiceSettingsResponsesMutex    sync.RWMutex
	deviceControllerResponses            map[string]map[string]interface{}
	deviceControllerResponsesMutex       sync.RWMutex
	asPathAccessListResponses            map[string]map[string]interface{}
	asPathAccessListResponsesMutex       sync.RWMutex
	communityListResponses               map[string]map[string]interface{}
	communityListResponsesMutex          sync.RWMutex
	deviceSettingsResponses              map[string]map[string]interface{}
	deviceSettingsResponsesMutex         sync.RWMutex
	extendedCommunityListResponses       map[string]map[string]interface{}
	extendedCommunityListResponsesMutex  sync.RWMutex
	ipv4ListResponses                    map[string]map[string]interface{}
	ipv4ListResponsesMutex               sync.RWMutex
	ipv4PrefixListResponses              map[string]map[string]interface{}
	ipv4PrefixListResponsesMutex         sync.RWMutex
	ipv6ListResponses                    map[string]map[string]interface{}
	ipv6ListResponsesMutex               sync.RWMutex
	ipv6PrefixListResponses              map[string]map[string]interface{}
	ipv6PrefixListResponsesMutex         sync.RWMutex
	routeMapClauseResponses              map[string]map[string]interface{}
	routeMapClauseResponsesMutex         sync.RWMutex
	routeMapResponses                    map[string]map[string]interface{}
	routeMapResponsesMutex               sync.RWMutex
	sfpBreakoutResponses                 map[string]map[string]interface{}
	sfpBreakoutResponsesMutex            sync.RWMutex
	siteResponses                        map[string]map[string]interface{}
	siteResponsesMutex                   sync.RWMutex
	podResponses                         map[string]map[string]interface{}
	podResponsesMutex                    sync.RWMutex
	portAclResponses                     map[string]map[string]interface{}
	portAclResponsesMutex                sync.RWMutex
	sflowCollectorResponses              map[string]map[string]interface{}
	sflowCollectorResponsesMutex         sync.RWMutex
	diagnosticsProfileResponses          map[string]map[string]interface{}
	diagnosticsProfileResponsesMutex     sync.RWMutex
	diagnosticsPortProfileResponses      map[string]map[string]interface{}
	diagnosticsPortProfileResponsesMutex sync.RWMutex
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
		PutRequestType:   reflect.TypeOf(openapi.GatewaysPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.GatewaysPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "gateway"}
		},
	},
	"lag": {
		ResourceType:     "lag",
		PutRequestType:   reflect.TypeOf(openapi.LagsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.LagsPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "lag"}
		},
	},
	"tenant": {
		ResourceType:     "tenant",
		PutRequestType:   reflect.TypeOf(openapi.TenantsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.TenantsPutRequest{}),
		HasAutoGen:       true,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "tenant"}
		},
	},
	"service": {
		ResourceType:     "service",
		PutRequestType:   reflect.TypeOf(openapi.ServicesPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.ServicesPutRequest{}),
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
		PatchRequestType: reflect.TypeOf(openapi.BundlesPutRequest{}),
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
	"as_path_access_list": {
		ResourceType:     "as_path_access_list",
		PutRequestType:   reflect.TypeOf(openapi.AspathaccesslistsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.AspathaccesslistsPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "as_path_access_list"}
		},
	},
	"community_list": {
		ResourceType:     "community_list",
		PutRequestType:   reflect.TypeOf(openapi.CommunitylistsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.CommunitylistsPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "community_list"}
		},
	},
	"device_settings": {
		ResourceType:     "device_settings",
		PutRequestType:   reflect.TypeOf(openapi.DevicesettingsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.DevicesettingsPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "device_settings"}
		},
	},
	"extended_community_list": {
		ResourceType:     "extended_community_list",
		PutRequestType:   reflect.TypeOf(openapi.ExtendedcommunitylistsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.ExtendedcommunitylistsPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "extended_community_list"}
		},
	},
	"ipv4_list": {
		ResourceType:     "ipv4_list",
		PutRequestType:   reflect.TypeOf(openapi.Ipv4listsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.Ipv4listsPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "ipv4_list"}
		},
	},
	"ipv4_prefix_list": {
		ResourceType:     "ipv4_prefix_list",
		PutRequestType:   reflect.TypeOf(openapi.Ipv4prefixlistsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.Ipv4prefixlistsPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "ipv4_prefix_list"}
		},
	},
	"ipv6_list": {
		ResourceType:     "ipv6_list",
		PutRequestType:   reflect.TypeOf(openapi.Ipv6listsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.Ipv6listsPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "ipv6_list"}
		},
	},
	"ipv6_prefix_list": {
		ResourceType:     "ipv6_prefix_list",
		PutRequestType:   reflect.TypeOf(openapi.Ipv6prefixlistsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.Ipv6prefixlistsPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "ipv6_prefix_list"}
		},
	},
	"route_map_clause": {
		ResourceType:     "route_map_clause",
		PutRequestType:   reflect.TypeOf(openapi.RoutemapclausesPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.RoutemapclausesPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "route_map_clause"}
		},
	},
	"route_map": {
		ResourceType:     "route_map",
		PutRequestType:   reflect.TypeOf(openapi.RoutemapsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.RoutemapsPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "route_map"}
		},
	},
	"sfp_breakout": {
		ResourceType:     "sfp_breakout",
		PutRequestType:   reflect.TypeOf(openapi.SfpbreakoutsPatchRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.SfpbreakoutsPatchRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "sfp_breakout"}
		},
	},
	"site": {
		ResourceType:     "site",
		PutRequestType:   reflect.TypeOf(openapi.SitesPatchRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.SitesPatchRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "site"}
		},
	},
	"pod": {
		ResourceType:     "pod",
		PutRequestType:   reflect.TypeOf(openapi.PodsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.PodsPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "pod"}
		},
	},
	"port_acl": {
		ResourceType:     "port_acl",
		PutRequestType:   reflect.TypeOf(openapi.PortaclsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.PortaclsPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "port_acl"}
		},
	},
	"sflow_collector": {
		ResourceType:     "sflow_collector",
		PutRequestType:   reflect.TypeOf(openapi.SflowcollectorsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.SflowcollectorsPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "sflow_collector"}
		},
	},
	"diagnostics_profile": {
		ResourceType:     "diagnostics_profile",
		PutRequestType:   reflect.TypeOf(openapi.DiagnosticsprofilesPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.DiagnosticsprofilesPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "diagnostics_profile"}
		},
	},
	"diagnostics_port_profile": {
		ResourceType:     "diagnostics_port_profile",
		PutRequestType:   reflect.TypeOf(openapi.DiagnosticsportprofilesPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.DiagnosticsportprofilesPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "diagnostics_port_profile"}
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
		client:                       client,
		contextProvider:              contextProvider,
		clearCacheFunc:               clearCacheFunc,
		mode:                         mode,
		lastOperationTime:            time.Now(),
		gatewayPut:                   make(map[string]openapi.GatewaysPutRequestGatewayValue),
		gatewayPatch:                 make(map[string]openapi.GatewaysPutRequestGatewayValue),
		gatewayDelete:                make([]string, 0),
		lagPut:                       make(map[string]openapi.LagsPutRequestLagValue),
		lagPatch:                     make(map[string]openapi.LagsPutRequestLagValue),
		lagDelete:                    make([]string, 0),
		tenantPut:                    make(map[string]openapi.TenantsPutRequestTenantValue),
		tenantPatch:                  make(map[string]openapi.TenantsPutRequestTenantValue),
		tenantDelete:                 make([]string, 0),
		servicePut:                   make(map[string]openapi.ServicesPutRequestServiceValue),
		servicePatch:                 make(map[string]openapi.ServicesPutRequestServiceValue),
		serviceDelete:                make([]string, 0),
		gatewayProfilePut:            make(map[string]openapi.GatewayprofilesPutRequestGatewayProfileValue),
		gatewayProfilePatch:          make(map[string]openapi.GatewayprofilesPutRequestGatewayProfileValue),
		gatewayProfileDelete:         make([]string, 0),
		ethPortProfilePut:            make(map[string]openapi.EthportprofilesPutRequestEthPortProfileValue),
		ethPortProfilePatch:          make(map[string]openapi.EthportprofilesPutRequestEthPortProfileValue),
		ethPortProfileDelete:         make([]string, 0),
		ethPortSettingsPut:           make(map[string]openapi.EthportsettingsPutRequestEthPortSettingsValue),
		ethPortSettingsPatch:         make(map[string]openapi.EthportsettingsPutRequestEthPortSettingsValue),
		ethPortSettingsDelete:        make([]string, 0),
		bundlePut:                    make(map[string]openapi.BundlesPutRequestEndpointBundleValue),
		bundlePatch:                  make(map[string]openapi.BundlesPutRequestEndpointBundleValue),
		bundleDelete:                 make([]string, 0),
		aclPut:                       make(map[string]openapi.AclsPutRequestIpFilterValue),
		aclPatch:                     make(map[string]openapi.AclsPutRequestIpFilterValue),
		aclDelete:                    make([]string, 0),
		aclIpVersion:                 make(map[string]string),
		authenticatedEthPortPut:      make(map[string]openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValue),
		authenticatedEthPortPatch:    make(map[string]openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValue),
		authenticatedEthPortDelete:   make([]string, 0),
		badgePut:                     make(map[string]openapi.BadgesPutRequestBadgeValue),
		badgePatch:                   make(map[string]openapi.BadgesPutRequestBadgeValue),
		badgeDelete:                  make([]string, 0),
		voicePortProfilePut:          make(map[string]openapi.VoiceportprofilesPutRequestVoicePortProfilesValue),
		voicePortProfilePatch:        make(map[string]openapi.VoiceportprofilesPutRequestVoicePortProfilesValue),
		voicePortProfileDelete:       make([]string, 0),
		switchpointPut:               make(map[string]openapi.SwitchpointsPutRequestSwitchpointValue),
		switchpointPatch:             make(map[string]openapi.SwitchpointsPutRequestSwitchpointValue),
		switchpointDelete:            make([]string, 0),
		servicePortProfilePut:        make(map[string]openapi.ServiceportprofilesPutRequestServicePortProfileValue),
		servicePortProfilePatch:      make(map[string]openapi.ServiceportprofilesPutRequestServicePortProfileValue),
		servicePortProfileDelete:     make([]string, 0),
		packetBrokerPut:              make(map[string]openapi.PacketbrokerPutRequestPbEgressProfileValue),
		packetBrokerPatch:            make(map[string]openapi.PacketbrokerPutRequestPbEgressProfileValue),
		packetBrokerDelete:           make([]string, 0),
		packetQueuePut:               make(map[string]openapi.PacketqueuesPutRequestPacketQueueValue),
		packetQueuePatch:             make(map[string]openapi.PacketqueuesPutRequestPacketQueueValue),
		packetQueueDelete:            make([]string, 0),
		deviceVoiceSettingsPut:       make(map[string]openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValue),
		deviceVoiceSettingsPatch:     make(map[string]openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValue),
		deviceVoiceSettingsDelete:    make([]string, 0),
		deviceControllerPut:          make(map[string]openapi.DevicecontrollersPutRequestDeviceControllerValue),
		deviceControllerPatch:        make(map[string]openapi.DevicecontrollersPutRequestDeviceControllerValue),
		deviceControllerDelete:       make([]string, 0),
		asPathAccessListPut:          make(map[string]openapi.AspathaccesslistsPutRequestAsPathAccessListValue),
		asPathAccessListPatch:        make(map[string]openapi.AspathaccesslistsPutRequestAsPathAccessListValue),
		asPathAccessListDelete:       make([]string, 0),
		communityListPut:             make(map[string]openapi.CommunitylistsPutRequestCommunityListValue),
		communityListPatch:           make(map[string]openapi.CommunitylistsPutRequestCommunityListValue),
		communityListDelete:          make([]string, 0),
		deviceSettingsPut:            make(map[string]openapi.DevicesettingsPutRequestEthDeviceProfilesValue),
		deviceSettingsPatch:          make(map[string]openapi.DevicesettingsPutRequestEthDeviceProfilesValue),
		deviceSettingsDelete:         make([]string, 0),
		extendedCommunityListPut:     make(map[string]openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValue),
		extendedCommunityListPatch:   make(map[string]openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValue),
		extendedCommunityListDelete:  make([]string, 0),
		ipv4ListPut:                  make(map[string]openapi.Ipv4listsPutRequestIpv4ListFilterValue),
		ipv4ListPatch:                make(map[string]openapi.Ipv4listsPutRequestIpv4ListFilterValue),
		ipv4ListDelete:               make([]string, 0),
		ipv4PrefixListPut:            make(map[string]openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValue),
		ipv4PrefixListPatch:          make(map[string]openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValue),
		ipv4PrefixListDelete:         make([]string, 0),
		ipv6ListPut:                  make(map[string]openapi.Ipv6listsPutRequestIpv6ListFilterValue),
		ipv6ListPatch:                make(map[string]openapi.Ipv6listsPutRequestIpv6ListFilterValue),
		ipv6ListDelete:               make([]string, 0),
		ipv6PrefixListPut:            make(map[string]openapi.Ipv6prefixlistsPutRequestIpv6PrefixListValue),
		ipv6PrefixListPatch:          make(map[string]openapi.Ipv6prefixlistsPutRequestIpv6PrefixListValue),
		ipv6PrefixListDelete:         make([]string, 0),
		routeMapClausePut:            make(map[string]openapi.RoutemapclausesPutRequestRouteMapClauseValue),
		routeMapClausePatch:          make(map[string]openapi.RoutemapclausesPutRequestRouteMapClauseValue),
		routeMapClauseDelete:         make([]string, 0),
		routeMapPut:                  make(map[string]openapi.RoutemapsPutRequestRouteMapValue),
		routeMapPatch:                make(map[string]openapi.RoutemapsPutRequestRouteMapValue),
		routeMapDelete:               make([]string, 0),
		sfpBreakoutPatch:             make(map[string]openapi.SfpbreakoutsPatchRequestSfpBreakoutsValue),
		sitePatch:                    make(map[string]openapi.SitesPatchRequestSiteValue),
		podPut:                       make(map[string]openapi.PodsPutRequestPodValue),
		podPatch:                     make(map[string]openapi.PodsPutRequestPodValue),
		podDelete:                    make([]string, 0),
		portAclPut:                   make(map[string]openapi.PortaclsPutRequestPortAclValue),
		portAclPatch:                 make(map[string]openapi.PortaclsPutRequestPortAclValue),
		portAclDelete:                make([]string, 0),
		sflowCollectorPut:            make(map[string]openapi.SflowcollectorsPutRequestSflowCollectorValue),
		sflowCollectorPatch:          make(map[string]openapi.SflowcollectorsPutRequestSflowCollectorValue),
		sflowCollectorDelete:         make([]string, 0),
		diagnosticsProfilePut:        make(map[string]openapi.DiagnosticsprofilesPutRequestDiagnosticsProfileValue),
		diagnosticsProfilePatch:      make(map[string]openapi.DiagnosticsprofilesPutRequestDiagnosticsProfileValue),
		diagnosticsProfileDelete:     make([]string, 0),
		diagnosticsPortProfilePut:    make(map[string]openapi.DiagnosticsportprofilesPutRequestDiagnosticsPortProfileValue),
		diagnosticsPortProfilePatch:  make(map[string]openapi.DiagnosticsportprofilesPutRequestDiagnosticsPortProfileValue),
		diagnosticsPortProfileDelete: make([]string, 0),

		pendingOperations:     make(map[string]*Operation),
		operationResults:      make(map[string]bool),
		operationErrors:       make(map[string]error),
		operationWaitChannels: make(map[string]chan struct{}),
		closedChannels:        make(map[string]bool),

		// Initialize with no recent operations
		recentGatewayOps:                false,
		recentLagOps:                    false,
		recentServiceOps:                false,
		recentTenantOps:                 false,
		recentGatewayProfileOps:         false,
		recentEthPortProfileOps:         false,
		recentEthPortSettingsOps:        false,
		recentBundleOps:                 false,
		recentAclOps:                    false,
		recentAuthenticatedEthPortOps:   false,
		recentBadgeOps:                  false,
		recentVoicePortProfileOps:       false,
		recentSwitchpointOps:            false,
		recentServicePortProfileOps:     false,
		recentPacketBrokerOps:           false,
		recentPacketQueueOps:            false,
		recentDeviceVoiceSettingsOps:    false,
		recentDeviceControllerOps:       false,
		recentAsPathAccessListOps:       false,
		recentCommunityListOps:          false,
		recentDeviceSettingsOps:         false,
		recentExtendedCommunityListOps:  false,
		recentIpv4ListOps:               false,
		recentIpv4PrefixListOps:         false,
		recentIpv6ListOps:               false,
		recentIpv6PrefixListOps:         false,
		recentRouteMapClauseOps:         false,
		recentRouteMapOps:               false,
		recentSfpBreakoutOps:            false,
		recentSiteOps:                   false,
		recentPodOps:                    false,
		recentPortAclOps:                false,
		recentSflowCollectorOps:         false,
		recentDiagnosticsProfileOps:     false,
		recentDiagnosticsPortProfileOps: false,

		// Initialize response caches
		gatewayResponses:                     make(map[string]map[string]interface{}),
		gatewayResponsesMutex:                sync.RWMutex{},
		lagResponses:                         make(map[string]map[string]interface{}),
		lagResponsesMutex:                    sync.RWMutex{},
		serviceResponses:                     make(map[string]map[string]interface{}),
		serviceResponsesMutex:                sync.RWMutex{},
		tenantResponses:                      make(map[string]map[string]interface{}),
		tenantResponsesMutex:                 sync.RWMutex{},
		gatewayProfileResponses:              make(map[string]map[string]interface{}),
		gatewayProfileResponsesMutex:         sync.RWMutex{},
		ethPortProfileResponses:              make(map[string]map[string]interface{}),
		ethPortProfileResponsesMutex:         sync.RWMutex{},
		ethPortSettingsResponses:             make(map[string]map[string]interface{}),
		ethPortSettingsResponsesMutex:        sync.RWMutex{},
		bundleResponses:                      make(map[string]map[string]interface{}),
		bundleResponsesMutex:                 sync.RWMutex{},
		aclResponses:                         make(map[string]map[string]interface{}),
		aclResponsesMutex:                    sync.RWMutex{},
		authenticatedEthPortResponses:        make(map[string]map[string]interface{}),
		authenticatedEthPortResponsesMutex:   sync.RWMutex{},
		badgeResponses:                       make(map[string]map[string]interface{}),
		badgeResponsesMutex:                  sync.RWMutex{},
		voicePortProfileResponses:            make(map[string]map[string]interface{}),
		voicePortProfileResponsesMutex:       sync.RWMutex{},
		switchpointResponses:                 make(map[string]map[string]interface{}),
		switchpointResponsesMutex:            sync.RWMutex{},
		servicePortProfileResponses:          make(map[string]map[string]interface{}),
		servicePortProfileResponsesMutex:     sync.RWMutex{},
		packetBrokerResponses:                make(map[string]map[string]interface{}),
		packetBrokerResponsesMutex:           sync.RWMutex{},
		packetQueueResponses:                 make(map[string]map[string]interface{}),
		packetQueueResponsesMutex:            sync.RWMutex{},
		deviceVoiceSettingsResponses:         make(map[string]map[string]interface{}),
		deviceVoiceSettingsResponsesMutex:    sync.RWMutex{},
		deviceControllerResponses:            make(map[string]map[string]interface{}),
		deviceControllerResponsesMutex:       sync.RWMutex{},
		asPathAccessListResponses:            make(map[string]map[string]interface{}),
		asPathAccessListResponsesMutex:       sync.RWMutex{},
		communityListResponses:               make(map[string]map[string]interface{}),
		communityListResponsesMutex:          sync.RWMutex{},
		deviceSettingsResponses:              make(map[string]map[string]interface{}),
		deviceSettingsResponsesMutex:         sync.RWMutex{},
		extendedCommunityListResponses:       make(map[string]map[string]interface{}),
		extendedCommunityListResponsesMutex:  sync.RWMutex{},
		ipv4ListResponses:                    make(map[string]map[string]interface{}),
		ipv4ListResponsesMutex:               sync.RWMutex{},
		ipv4PrefixListResponses:              make(map[string]map[string]interface{}),
		ipv4PrefixListResponsesMutex:         sync.RWMutex{},
		ipv6ListResponses:                    make(map[string]map[string]interface{}),
		ipv6ListResponsesMutex:               sync.RWMutex{},
		ipv6PrefixListResponses:              make(map[string]map[string]interface{}),
		ipv6PrefixListResponsesMutex:         sync.RWMutex{},
		routeMapClauseResponses:              make(map[string]map[string]interface{}),
		routeMapClauseResponsesMutex:         sync.RWMutex{},
		routeMapResponses:                    make(map[string]map[string]interface{}),
		routeMapResponsesMutex:               sync.RWMutex{},
		sfpBreakoutResponses:                 make(map[string]map[string]interface{}),
		sfpBreakoutResponsesMutex:            sync.RWMutex{},
		siteResponses:                        make(map[string]map[string]interface{}),
		siteResponsesMutex:                   sync.RWMutex{},
		podResponses:                         make(map[string]map[string]interface{}),
		podResponsesMutex:                    sync.RWMutex{},
		portAclResponses:                     make(map[string]map[string]interface{}),
		portAclResponsesMutex:                sync.RWMutex{},
		sflowCollectorResponses:              make(map[string]map[string]interface{}),
		sflowCollectorResponsesMutex:         sync.RWMutex{},
		diagnosticsProfileResponses:          make(map[string]map[string]interface{}),
		diagnosticsProfileResponsesMutex:     sync.RWMutex{},
		diagnosticsPortProfileResponses:      make(map[string]map[string]interface{}),
		diagnosticsPortProfileResponsesMutex: sync.RWMutex{},
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

	case "as_path_access_list":
		b.asPathAccessListResponsesMutex.RLock()
		defer b.asPathAccessListResponsesMutex.RUnlock()
		response, exists := b.asPathAccessListResponses[resourceName]
		return response, exists

	case "community_list":
		b.communityListResponsesMutex.RLock()
		defer b.communityListResponsesMutex.RUnlock()
		response, exists := b.communityListResponses[resourceName]
		return response, exists

	case "device_settings":
		b.deviceSettingsResponsesMutex.RLock()
		defer b.deviceSettingsResponsesMutex.RUnlock()
		response, exists := b.deviceSettingsResponses[resourceName]
		return response, exists

	case "extended_community_list":
		b.extendedCommunityListResponsesMutex.RLock()
		defer b.extendedCommunityListResponsesMutex.RUnlock()
		response, exists := b.extendedCommunityListResponses[resourceName]
		return response, exists

	case "ipv4_list":
		b.ipv4ListResponsesMutex.RLock()
		defer b.ipv4ListResponsesMutex.RUnlock()
		response, exists := b.ipv4ListResponses[resourceName]
		return response, exists

	case "ipv4_prefix_list":
		b.ipv4PrefixListResponsesMutex.RLock()
		defer b.ipv4PrefixListResponsesMutex.RUnlock()
		response, exists := b.ipv4PrefixListResponses[resourceName]
		return response, exists

	case "ipv6_list":
		b.ipv6ListResponsesMutex.RLock()
		defer b.ipv6ListResponsesMutex.RUnlock()
		response, exists := b.ipv6ListResponses[resourceName]
		return response, exists

	case "ipv6_prefix_list":
		b.ipv6PrefixListResponsesMutex.RLock()
		defer b.ipv6PrefixListResponsesMutex.RUnlock()
		response, exists := b.ipv6PrefixListResponses[resourceName]
		return response, exists

	case "route_map_clause":
		b.routeMapClauseResponsesMutex.RLock()
		defer b.routeMapClauseResponsesMutex.RUnlock()
		response, exists := b.routeMapClauseResponses[resourceName]
		return response, exists

	case "route_map":
		b.routeMapResponsesMutex.RLock()
		defer b.routeMapResponsesMutex.RUnlock()
		response, exists := b.routeMapResponses[resourceName]
		return response, exists

	case "sfp_breakout":
		b.sfpBreakoutResponsesMutex.RLock()
		defer b.sfpBreakoutResponsesMutex.RUnlock()
		response, exists := b.sfpBreakoutResponses[resourceName]
		return response, exists

	case "site":
		b.siteResponsesMutex.RLock()
		defer b.siteResponsesMutex.RUnlock()
		response, exists := b.siteResponses[resourceName]
		return response, exists

	case "pod":
		b.podResponsesMutex.RLock()
		defer b.podResponsesMutex.RUnlock()
		response, exists := b.podResponses[resourceName]
		return response, exists

	case "port_acl":
		b.portAclResponsesMutex.RLock()
		defer b.portAclResponsesMutex.RUnlock()
		response, exists := b.portAclResponses[resourceName]
		return response, exists

	case "sflow_collector":
		b.sflowCollectorResponsesMutex.RLock()
		defer b.sflowCollectorResponsesMutex.RUnlock()
		response, exists := b.sflowCollectorResponses[resourceName]
		return response, exists

	case "diagnostics_profile":
		b.diagnosticsProfileResponsesMutex.RLock()
		defer b.diagnosticsProfileResponsesMutex.RUnlock()
		response, exists := b.diagnosticsProfileResponses[resourceName]
		return response, exists

	case "diagnostics_port_profile":
		b.diagnosticsPortProfileResponsesMutex.RLock()
		defer b.diagnosticsPortProfileResponsesMutex.RUnlock()
		response, exists := b.diagnosticsPortProfileResponses[resourceName]
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
			b.clearCacheFunc(ctx, b.contextProvider(), "sflowcollectors")
			b.clearCacheFunc(ctx, b.contextProvider(), "diagnosticsprofiles")
			b.clearCacheFunc(ctx, b.contextProvider(), "diagnosticsportprofiles")
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
			b.clearCacheFunc(ctx, b.contextProvider(), "aspathaccesslists")
			b.clearCacheFunc(ctx, b.contextProvider(), "communitylists")
			b.clearCacheFunc(ctx, b.contextProvider(), "devicesettings")
			b.clearCacheFunc(ctx, b.contextProvider(), "extendedcommunitylists")
			b.clearCacheFunc(ctx, b.contextProvider(), "ipv4lists")
			b.clearCacheFunc(ctx, b.contextProvider(), "ipv4prefixlists")
			b.clearCacheFunc(ctx, b.contextProvider(), "ipv6lists")
			b.clearCacheFunc(ctx, b.contextProvider(), "ipv6prefixlists")
			b.clearCacheFunc(ctx, b.contextProvider(), "routemapclauses")
			b.clearCacheFunc(ctx, b.contextProvider(), "routemaps")
			b.clearCacheFunc(ctx, b.contextProvider(), "sfpbreakouts")
			b.clearCacheFunc(ctx, b.contextProvider(), "sites")
			b.clearCacheFunc(ctx, b.contextProvider(), "pods")
			b.clearCacheFunc(ctx, b.contextProvider(), "portacls")
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
	if !execute("PUT", len(b.deviceSettingsPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "device_settings", "PUT") }, "Device Settings") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.lagPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "lag", "PUT") }, "LAG") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.sflowCollectorPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "sflow_collector", "PUT") }, "SFlow Collector") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.diagnosticsProfilePut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "diagnostics_profile", "PUT") }, "Diagnostics Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.diagnosticsPortProfilePut), func(ctx context.Context) diag.Diagnostics {
		return b.ExecuteBulk(ctx, "diagnostics_port_profile", "PUT")
	}, "Diagnostics Port Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.bundlePut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "bundle", "PUT") }, "Bundle") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.aclPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "acl", "PUT") }, "ACL") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.ipv4PrefixListPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv4_prefix_list", "PUT") }, "IPv4 Prefix List") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.ipv6PrefixListPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv6_prefix_list", "PUT") }, "IPv6 Prefix List") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.ipv4ListPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv4_list", "PUT") }, "IPv4 List") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.ipv6ListPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv6_list", "PUT") }, "IPv6 List") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.packetBrokerPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "packet_broker", "PUT") }, "Packet Broker") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.portAclPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "port_acl", "PUT") }, "Port ACL") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.badgePut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "badge", "PUT") }, "Badge") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.podPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "pod", "PUT") }, "Pod") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.switchpointPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "switchpoint", "PUT") }, "Switchpoint") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.deviceControllerPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "device_controller", "PUT") }, "Device Controller") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.asPathAccessListPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "as_path_access_list", "PUT") }, "AS Path Access List") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.communityListPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "community_list", "PUT") }, "Community List") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.extendedCommunityListPut), func(ctx context.Context) diag.Diagnostics {
		return b.ExecuteBulk(ctx, "extended_community_list", "PUT")
	}, "Extended Community List") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.routeMapClausePut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "route_map_clause", "PUT") }, "Route Map Clause") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.routeMapPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "route_map", "PUT") }, "Route Map") {
		return diagnostics, operationsPerformed
	}

	// PATCH operations - DC Order
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
	if !execute("PATCH", len(b.deviceSettingsPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "device_settings", "PATCH") }, "Device Settings") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.lagPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "lag", "PATCH") }, "LAG") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.sflowCollectorPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "sflow_collector", "PATCH") }, "SFlow Collector") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.diagnosticsProfilePatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "diagnostics_profile", "PATCH") }, "Diagnostics Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.diagnosticsPortProfilePatch), func(ctx context.Context) diag.Diagnostics {
		return b.ExecuteBulk(ctx, "diagnostics_port_profile", "PATCH")
	}, "Diagnostics Port Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.bundlePatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "bundle", "PATCH") }, "Bundle") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.aclPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "acl", "PATCH") }, "ACL") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.ipv4PrefixListPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv4_prefix_list", "PATCH") }, "IPv4 Prefix List") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.ipv6PrefixListPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv6_prefix_list", "PATCH") }, "IPv6 Prefix List") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.ipv4ListPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv4_list", "PATCH") }, "IPv4 List") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.ipv6ListPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv6_list", "PATCH") }, "IPv6 List") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.packetBrokerPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "packet_broker", "PATCH") }, "Packet Broker") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.portAclPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "port_acl", "PATCH") }, "Port ACL") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.badgePatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "badge", "PATCH") }, "Badge") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.podPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "pod", "PATCH") }, "Pod") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.switchpointPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "switchpoint", "PATCH") }, "Switchpoint") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.deviceControllerPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "device_controller", "PATCH") }, "Device Controller") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.asPathAccessListPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "as_path_access_list", "PATCH") }, "AS Path Access List") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.communityListPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "community_list", "PATCH") }, "Community List") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.extendedCommunityListPatch), func(ctx context.Context) diag.Diagnostics {
		return b.ExecuteBulk(ctx, "extended_community_list", "PATCH")
	}, "Extended Community List") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.routeMapClausePatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "route_map_clause", "PATCH") }, "Route Map Clause") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.routeMapPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "route_map", "PATCH") }, "Route Map") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.sfpBreakoutPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "sfp_breakout", "PATCH") }, "SFP Breakout") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.sitePatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "site", "PATCH") }, "Site") {
		return diagnostics, operationsPerformed
	}

	// DELETE operations - Reverse DC Order
	if !execute("DELETE", len(b.routeMapDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "route_map", "DELETE") }, "Route Map") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.routeMapClauseDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "route_map_clause", "DELETE") }, "Route Map Clause") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.extendedCommunityListDelete), func(ctx context.Context) diag.Diagnostics {
		return b.ExecuteBulk(ctx, "extended_community_list", "DELETE")
	}, "Extended Community List") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.communityListDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "community_list", "DELETE") }, "Community List") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.asPathAccessListDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "as_path_access_list", "DELETE") }, "AS Path Access List") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.deviceControllerDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "device_controller", "DELETE") }, "Device Controller") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.switchpointDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "switchpoint", "DELETE") }, "Switchpoint") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.podDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "pod", "DELETE") }, "Pod") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.badgeDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "badge", "DELETE") }, "Badge") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.portAclDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "port_acl", "DELETE") }, "Port ACL") {
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
	if !execute("DELETE", len(b.diagnosticsPortProfileDelete), func(ctx context.Context) diag.Diagnostics {
		return b.ExecuteBulk(ctx, "diagnostics_port_profile", "DELETE")
	}, "Diagnostics Port Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.diagnosticsProfileDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "diagnostics_profile", "DELETE") }, "Diagnostics Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.sflowCollectorDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "sflow_collector", "DELETE") }, "SFlow Collector") {
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
	if !execute("DELETE", len(b.ipv6PrefixListDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv6_prefix_list", "DELETE") }, "IPv6 Prefix List") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.ipv6ListDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv6_list_filter", "DELETE") }, "IPv6 List Filter") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.ipv4PrefixListDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv4_prefix_list", "DELETE") }, "IPv4 Prefix List") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.ipv4ListDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv4_list_filter", "DELETE") }, "IPv4 List Filter") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.deviceSettingsDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "device_settings", "DELETE") }, "Device Settings") {
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
	if !execute("PUT", len(b.deviceSettingsPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "device_settings", "PUT") }, "Device Settings") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.lagPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "lag", "PUT") }, "LAG") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.sflowCollectorPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "sflow_collector", "PUT") }, "SFlow Collector") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.diagnosticsProfilePut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "diagnostics_profile", "PUT") }, "Diagnostics Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.diagnosticsPortProfilePut), func(ctx context.Context) diag.Diagnostics {
		return b.ExecuteBulk(ctx, "diagnostics_port_profile", "PUT")
	}, "Diagnostics Port Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.bundlePut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "bundle", "PUT") }, "Bundle") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.aclPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "acl", "PUT") }, "ACL") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.ipv4ListPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv4_list", "PUT") }, "IPv4 List") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.ipv6ListPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv6_list", "PUT") }, "IPv6 List") {
		return diagnostics, operationsPerformed
	}
	if !execute("PUT", len(b.portAclPut), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "port_acl", "PUT") }, "Port ACL") {
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
	if !execute("PATCH", len(b.deviceSettingsPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "device_settings", "PATCH") }, "Device Settings") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.lagPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "lag", "PATCH") }, "LAG") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.sflowCollectorPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "sflow_collector", "PATCH") }, "SFlow Collector") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.diagnosticsProfilePatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "diagnostics_profile", "PATCH") }, "Diagnostics Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.diagnosticsPortProfilePatch), func(ctx context.Context) diag.Diagnostics {
		return b.ExecuteBulk(ctx, "diagnostics_port_profile", "PATCH")
	}, "Diagnostics Port Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.bundlePatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "bundle", "PATCH") }, "Bundle") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.aclPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "acl", "PATCH") }, "ACL") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.ipv4ListPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv4_list", "PATCH") }, "IPv4 List") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.ipv6ListPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv6_list", "PATCH") }, "IPv6 List") {
		return diagnostics, operationsPerformed
	}
	if !execute("PATCH", len(b.portAclPatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "port_acl", "PATCH") }, "Port ACL") {
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
	if !execute("PATCH", len(b.sitePatch), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "site", "PATCH") }, "Site") {
		return diagnostics, operationsPerformed
	}

	// DELETE operations - Reverse Campus Order
	if !execute("DELETE", len(b.deviceControllerDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "device_controller", "DELETE") }, "Device Controller") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.switchpointDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "switchpoint", "DELETE") }, "Switchpoint") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.podDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "pod", "DELETE") }, "Pod") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.portAclDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "port_acl", "DELETE") }, "Port ACL") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.badgeDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "badge", "DELETE") }, "Badge") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.bundleDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "bundle", "DELETE") }, "Bundle") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.diagnosticsPortProfileDelete), func(ctx context.Context) diag.Diagnostics {
		return b.ExecuteBulk(ctx, "diagnostics_port_profile", "DELETE")
	}, "Diagnostics Port Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.diagnosticsProfileDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "diagnostics_profile", "DELETE") }, "Diagnostics Profile") {
		return diagnostics, operationsPerformed
	}
	if !execute("DELETE", len(b.sflowCollectorDelete), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "sflow_collector", "DELETE") }, "SFlow Collector") {
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
		len(b.deviceSettingsPut) == 0 && len(b.deviceSettingsPatch) == 0 && len(b.deviceSettingsDelete) == 0 &&
		len(b.bundlePut) == 0 && len(b.bundlePatch) == 0 && len(b.bundleDelete) == 0 && len(b.authenticatedEthPortPut) == 0 && len(b.authenticatedEthPortPatch) == 0 &&
		len(b.authenticatedEthPortDelete) == 0 && len(b.aclPut) == 0 && len(b.aclPatch) == 0 && len(b.aclDelete) == 0 &&
		len(b.ipv4ListPut) == 0 && len(b.ipv4ListPatch) == 0 && len(b.ipv4ListDelete) == 0 &&
		len(b.ipv4PrefixListPut) == 0 && len(b.ipv4PrefixListPatch) == 0 && len(b.ipv4PrefixListDelete) == 0 &&
		len(b.ipv6ListPut) == 0 && len(b.ipv6ListPatch) == 0 && len(b.ipv6ListDelete) == 0 &&
		len(b.ipv6PrefixListPut) == 0 && len(b.ipv6PrefixListPatch) == 0 && len(b.ipv6PrefixListDelete) == 0 &&
		len(b.badgePut) == 0 && len(b.badgePatch) == 0 && len(b.badgeDelete) == 0 &&
		len(b.voicePortProfilePut) == 0 && len(b.voicePortProfilePatch) == 0 && len(b.voicePortProfileDelete) == 0 &&
		len(b.switchpointPut) == 0 && len(b.switchpointPatch) == 0 && len(b.switchpointDelete) == 0 &&
		len(b.servicePortProfilePut) == 0 && len(b.servicePortProfilePatch) == 0 && len(b.servicePortProfileDelete) == 0 &&
		len(b.packetBrokerPut) == 0 && len(b.packetBrokerPatch) == 0 && len(b.packetBrokerDelete) == 0 &&
		len(b.packetQueuePut) == 0 && len(b.packetQueuePatch) == 0 && len(b.packetQueueDelete) == 0 &&
		len(b.deviceVoiceSettingsPut) == 0 && len(b.deviceVoiceSettingsPatch) == 0 && len(b.deviceVoiceSettingsDelete) == 0 &&
		len(b.asPathAccessListPut) == 0 && len(b.asPathAccessListPatch) == 0 && len(b.asPathAccessListDelete) == 0 &&
		len(b.communityListPut) == 0 && len(b.communityListPatch) == 0 && len(b.communityListDelete) == 0 &&
		len(b.extendedCommunityListPut) == 0 && len(b.extendedCommunityListPatch) == 0 && len(b.extendedCommunityListDelete) == 0 &&
		len(b.routeMapClausePut) == 0 && len(b.routeMapClausePatch) == 0 && len(b.routeMapClauseDelete) == 0 &&
		len(b.routeMapPut) == 0 && len(b.routeMapPatch) == 0 && len(b.routeMapDelete) == 0 &&
		len(b.sfpBreakoutPatch) == 0 &&
		len(b.sitePatch) == 0 &&
		len(b.podPut) == 0 && len(b.podPatch) == 0 && len(b.podDelete) == 0 &&
		len(b.portAclPut) == 0 && len(b.portAclPatch) == 0 && len(b.portAclDelete) == 0 &&
		len(b.sflowCollectorPut) == 0 && len(b.sflowCollectorPatch) == 0 && len(b.sflowCollectorDelete) == 0 &&
		len(b.diagnosticsProfilePut) == 0 && len(b.diagnosticsProfilePatch) == 0 && len(b.diagnosticsProfileDelete) == 0 &&
		len(b.diagnosticsPortProfilePut) == 0 && len(b.diagnosticsPortProfilePatch) == 0 && len(b.diagnosticsPortProfileDelete) == 0 &&
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

	deviceSettingsPutCount := len(b.deviceSettingsPut)
	deviceSettingsPatchCount := len(b.deviceSettingsPatch)
	deviceSettingsDeleteCount := len(b.deviceSettingsDelete)

	sflowCollectorPutCount := len(b.sflowCollectorPut)
	sflowCollectorPatchCount := len(b.sflowCollectorPatch)
	sflowCollectorDeleteCount := len(b.sflowCollectorDelete)

	diagnosticsProfilePutCount := len(b.diagnosticsProfilePut)
	diagnosticsProfilePatchCount := len(b.diagnosticsProfilePatch)
	diagnosticsProfileDeleteCount := len(b.diagnosticsProfileDelete)

	diagnosticsPortProfilePutCount := len(b.diagnosticsPortProfilePut)
	diagnosticsPortProfilePatchCount := len(b.diagnosticsPortProfilePatch)
	diagnosticsPortProfileDeleteCount := len(b.diagnosticsPortProfileDelete)

	bundlePutCount := len(b.bundlePut)
	bundlePatchCount := len(b.bundlePatch)
	bundleDeleteCount := len(b.bundleDelete)

	aclPutCount := len(b.aclPut)
	aclPatchCount := len(b.aclPatch)
	aclDeleteCount := len(b.aclDelete)

	ipv4ListPutCount := len(b.ipv4ListPut)
	ipv4ListPatchCount := len(b.ipv4ListPatch)
	ipv4ListDeleteCount := len(b.ipv4ListDelete)

	ipv4PrefixListPutCount := len(b.ipv4PrefixListPut)
	ipv4PrefixListPatchCount := len(b.ipv4PrefixListPatch)
	ipv4PrefixListDeleteCount := len(b.ipv4PrefixListDelete)

	ipv6ListPutCount := len(b.ipv6ListPut)
	ipv6ListPatchCount := len(b.ipv6ListPatch)
	ipv6ListDeleteCount := len(b.ipv6ListDelete)

	ipv6PrefixListPutCount := len(b.ipv6PrefixListPut)
	ipv6PrefixListPatchCount := len(b.ipv6PrefixListPatch)
	ipv6PrefixListDeleteCount := len(b.ipv6PrefixListDelete)

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

	asPathAccessListPutCount := len(b.asPathAccessListPut)
	asPathAccessListPatchCount := len(b.asPathAccessListPatch)
	asPathAccessListDeleteCount := len(b.asPathAccessListDelete)

	communityListPutCount := len(b.communityListPut)
	communityListPatchCount := len(b.communityListPatch)
	communityListDeleteCount := len(b.communityListDelete)

	extendedCommunityListPutCount := len(b.extendedCommunityListPut)
	extendedCommunityListPatchCount := len(b.extendedCommunityListPatch)
	extendedCommunityListDeleteCount := len(b.extendedCommunityListDelete)

	routeMapClausePutCount := len(b.routeMapClausePut)
	routeMapClausePatchCount := len(b.routeMapClausePatch)
	routeMapClauseDeleteCount := len(b.routeMapClauseDelete)

	routeMapPutCount := len(b.routeMapPut)
	routeMapPatchCount := len(b.routeMapPatch)
	routeMapDeleteCount := len(b.routeMapDelete)

	sfpBreakoutPatchCount := len(b.sfpBreakoutPatch)

	sitePatchCount := len(b.sitePatch)

	podPutCount := len(b.podPut)
	podPatchCount := len(b.podPatch)
	podDeleteCount := len(b.podDelete)

	portAclPutCount := len(b.portAclPut)
	portAclPatchCount := len(b.portAclPatch)
	portAclDeleteCount := len(b.portAclDelete)

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
		portAclPutCount + portAclPatchCount + portAclDeleteCount +
		authenticatedEthPortPutCount + authenticatedEthPortPatchCount + authenticatedEthPortDeleteCount +
		deviceControllerPutCount + deviceControllerPatchCount + deviceControllerDeleteCount

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
			"total_count":                           totalCount,
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
	case "device_settings":
		return &ResourceOperationData{
			PutOperations:    b.deviceSettingsPut,
			PatchOperations:  b.deviceSettingsPatch,
			DeleteOperations: &b.deviceSettingsDelete,
			RecentOps:        &b.recentDeviceSettingsOps,
			RecentOpTime:     &b.recentDeviceSettingsOpTime,
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
	case "ipv4_list_filter":
		return &ResourceOperationData{
			PutOperations:    b.ipv4ListPut,
			PatchOperations:  b.ipv4ListPatch,
			DeleteOperations: &b.ipv4ListDelete,
			RecentOps:        &b.recentIpv4ListOps,
			RecentOpTime:     &b.recentIpv4ListOpTime,
		}
	case "ipv4_prefix_list":
		return &ResourceOperationData{
			PutOperations:    b.ipv4PrefixListPut,
			PatchOperations:  b.ipv4PrefixListPatch,
			DeleteOperations: &b.ipv4PrefixListDelete,
			RecentOps:        &b.recentIpv4PrefixListOps,
			RecentOpTime:     &b.recentIpv4PrefixListOpTime,
		}
	case "ipv6_list_filter":
		return &ResourceOperationData{
			PutOperations:    b.ipv6ListPut,
			PatchOperations:  b.ipv6ListPatch,
			DeleteOperations: &b.ipv6ListDelete,
			RecentOps:        &b.recentIpv6ListOps,
			RecentOpTime:     &b.recentIpv6ListOpTime,
		}
	case "ipv6_prefix_list":
		return &ResourceOperationData{
			PutOperations:    b.ipv6PrefixListPut,
			PatchOperations:  b.ipv6PrefixListPatch,
			DeleteOperations: &b.ipv6PrefixListDelete,
			RecentOps:        &b.recentIpv6PrefixListOps,
			RecentOpTime:     &b.recentIpv6PrefixListOpTime,
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
	case "as_path_access_list":
		return &ResourceOperationData{
			PutOperations:    b.asPathAccessListPut,
			PatchOperations:  b.asPathAccessListPatch,
			DeleteOperations: &b.asPathAccessListDelete,
			RecentOps:        &b.recentAsPathAccessListOps,
			RecentOpTime:     &b.recentAsPathAccessListOpTime,
		}
	case "community_list":
		return &ResourceOperationData{
			PutOperations:    b.communityListPut,
			PatchOperations:  b.communityListPatch,
			DeleteOperations: &b.communityListDelete,
			RecentOps:        &b.recentCommunityListOps,
			RecentOpTime:     &b.recentCommunityListOpTime,
		}
	case "extended_community_list":
		return &ResourceOperationData{
			PutOperations:    b.extendedCommunityListPut,
			PatchOperations:  b.extendedCommunityListPatch,
			DeleteOperations: &b.extendedCommunityListDelete,
			RecentOps:        &b.recentExtendedCommunityListOps,
			RecentOpTime:     &b.recentExtendedCommunityListOpTime,
		}
	case "route_map_clause":
		return &ResourceOperationData{
			PutOperations:    b.routeMapClausePut,
			PatchOperations:  b.routeMapClausePatch,
			DeleteOperations: &b.routeMapClauseDelete,
			RecentOps:        &b.recentRouteMapClauseOps,
			RecentOpTime:     &b.recentRouteMapClauseOpTime,
		}
	case "route_map":
		return &ResourceOperationData{
			PutOperations:    b.routeMapPut,
			PatchOperations:  b.routeMapPatch,
			DeleteOperations: &b.routeMapDelete,
			RecentOps:        &b.recentRouteMapOps,
			RecentOpTime:     &b.recentRouteMapOpTime,
		}
	case "sfp_breakout":
		return &ResourceOperationData{
			PutOperations:    nil, // SFP Breakouts only support PATCH
			PatchOperations:  b.sfpBreakoutPatch,
			DeleteOperations: nil, // No DELETE operations for SFP Breakouts
			RecentOps:        &b.recentSfpBreakoutOps,
			RecentOpTime:     &b.recentSfpBreakoutOpTime,
		}
	case "site":
		return &ResourceOperationData{
			PutOperations:    nil, // Sites only support PATCH
			PatchOperations:  b.sitePatch,
			DeleteOperations: nil, // No DELETE operations for Sites
			RecentOps:        &b.recentSiteOps,
			RecentOpTime:     &b.recentSiteOpTime,
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
	case "pod":
		return &ResourceOperationData{
			PutOperations:    b.podPut,
			PatchOperations:  b.podPatch,
			DeleteOperations: &b.podDelete,
			RecentOps:        &b.recentPodOps,
			RecentOpTime:     &b.recentPodOpTime,
		}
	case "port_acl":
		return &ResourceOperationData{
			PutOperations:    b.portAclPut,
			PatchOperations:  b.portAclPatch,
			DeleteOperations: &b.portAclDelete,
			RecentOps:        &b.recentPortAclOps,
			RecentOpTime:     &b.recentPortAclOpTime,
		}
	case "sflow_collector":
		return &ResourceOperationData{
			PutOperations:    b.sflowCollectorPut,
			PatchOperations:  b.sflowCollectorPatch,
			DeleteOperations: &b.sflowCollectorDelete,
			RecentOps:        &b.recentSflowCollectorOps,
			RecentOpTime:     &b.recentSflowCollectorOpTime,
		}
	case "diagnostics_profile":
		return &ResourceOperationData{
			PutOperations:    b.diagnosticsProfilePut,
			PatchOperations:  b.diagnosticsProfilePatch,
			DeleteOperations: &b.diagnosticsProfileDelete,
			RecentOps:        &b.recentDiagnosticsProfileOps,
			RecentOpTime:     &b.recentDiagnosticsProfileOpTime,
		}
	case "diagnostics_port_profile":
		return &ResourceOperationData{
			PutOperations:    b.diagnosticsPortProfilePut,
			PatchOperations:  b.diagnosticsPortProfilePatch,
			DeleteOperations: &b.diagnosticsPortProfileDelete,
			RecentOps:        &b.recentDiagnosticsPortProfileOps,
			RecentOpTime:     &b.recentDiagnosticsPortProfileOpTime,
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

	hasPending = hasPending || (data.DeleteOperations != nil && len(*data.DeleteOperations) > 0)

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
	case "device_settings":
		req := g.client.DeviceSettingsAPI.DevicesettingsPut(ctx).DevicesettingsPutRequest(*request.(*openapi.DevicesettingsPutRequest))
		return req.Execute()
	case "ipv4_list_filter":
		req := g.client.IPv4ListFiltersAPI.Ipv4listsPut(ctx).Ipv4listsPutRequest(*request.(*openapi.Ipv4listsPutRequest))
		return req.Execute()
	case "ipv4_prefix_list":
		req := g.client.IPv4PrefixListsAPI.Ipv4prefixlistsPut(ctx).Ipv4prefixlistsPutRequest(*request.(*openapi.Ipv4prefixlistsPutRequest))
		return req.Execute()
	case "ipv6_list_filter":
		req := g.client.IPv6ListFiltersAPI.Ipv6listsPut(ctx).Ipv6listsPutRequest(*request.(*openapi.Ipv6listsPutRequest))
		return req.Execute()
	case "ipv6_prefix_list":
		req := g.client.IPv6PrefixListsAPI.Ipv6prefixlistsPut(ctx).Ipv6prefixlistsPutRequest(*request.(*openapi.Ipv6prefixlistsPutRequest))
		return req.Execute()
	case "as_path_access_list":
		req := g.client.ASPathAccessListsAPI.AspathaccesslistsPut(ctx).AspathaccesslistsPutRequest(*request.(*openapi.AspathaccesslistsPutRequest))
		return req.Execute()
	case "community_list":
		req := g.client.CommunityListsAPI.CommunitylistsPut(ctx).CommunitylistsPutRequest(*request.(*openapi.CommunitylistsPutRequest))
		return req.Execute()
	case "extended_community_list":
		req := g.client.ExtendedCommunityListsAPI.ExtendedcommunitylistsPut(ctx).ExtendedcommunitylistsPutRequest(*request.(*openapi.ExtendedcommunitylistsPutRequest))
		return req.Execute()
	case "route_map_clause":
		req := g.client.RouteMapClausesAPI.RoutemapclausesPut(ctx).RoutemapclausesPutRequest(*request.(*openapi.RoutemapclausesPutRequest))
		return req.Execute()
	case "route_map":
		req := g.client.RouteMapsAPI.RoutemapsPut(ctx).RoutemapsPutRequest(*request.(*openapi.RoutemapsPutRequest))
		return req.Execute()
	case "pod":
		req := g.client.PodsAPI.PodsPut(ctx).PodsPutRequest(*request.(*openapi.PodsPutRequest))
		return req.Execute()
	case "port_acl":
		req := g.client.PortACLsAPI.PortaclsPut(ctx).PortaclsPutRequest(*request.(*openapi.PortaclsPutRequest))
		return req.Execute()
	case "sflow_collector":
		req := g.client.SFlowCollectorsAPI.SflowcollectorsPut(ctx).SflowcollectorsPutRequest(*request.(*openapi.SflowcollectorsPutRequest))
		return req.Execute()
	case "diagnostics_profile":
		req := g.client.DiagnosticsProfilesAPI.DiagnosticsprofilesPut(ctx).DiagnosticsprofilesPutRequest(*request.(*openapi.DiagnosticsprofilesPutRequest))
		return req.Execute()
	case "diagnostics_port_profile":
		req := g.client.DiagnosticsPortProfilesAPI.DiagnosticsportprofilesPut(ctx).DiagnosticsportprofilesPutRequest(*request.(*openapi.DiagnosticsportprofilesPutRequest))
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
		req := g.client.BundlesAPI.BundlesPatch(ctx).BundlesPutRequest(*request.(*openapi.BundlesPutRequest))
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
	case "device_settings":
		req := g.client.DeviceSettingsAPI.DevicesettingsPatch(ctx).DevicesettingsPutRequest(*request.(*openapi.DevicesettingsPutRequest))
		return req.Execute()
	case "ipv4_list_filter":
		req := g.client.IPv4ListFiltersAPI.Ipv4listsPatch(ctx).Ipv4listsPutRequest(*request.(*openapi.Ipv4listsPutRequest))
		return req.Execute()
	case "ipv4_prefix_list":
		req := g.client.IPv4PrefixListsAPI.Ipv4prefixlistsPatch(ctx).Ipv4prefixlistsPutRequest(*request.(*openapi.Ipv4prefixlistsPutRequest))
		return req.Execute()
	case "ipv6_list_filter":
		req := g.client.IPv6ListFiltersAPI.Ipv6listsPatch(ctx).Ipv6listsPutRequest(*request.(*openapi.Ipv6listsPutRequest))
		return req.Execute()
	case "ipv6_prefix_list":
		req := g.client.IPv6PrefixListsAPI.Ipv6prefixlistsPatch(ctx).Ipv6prefixlistsPutRequest(*request.(*openapi.Ipv6prefixlistsPutRequest))
		return req.Execute()
	case "as_path_access_list":
		req := g.client.ASPathAccessListsAPI.AspathaccesslistsPatch(ctx).AspathaccesslistsPutRequest(*request.(*openapi.AspathaccesslistsPutRequest))
		return req.Execute()
	case "community_list":
		req := g.client.CommunityListsAPI.CommunitylistsPatch(ctx).CommunitylistsPutRequest(*request.(*openapi.CommunitylistsPutRequest))
		return req.Execute()
	case "extended_community_list":
		req := g.client.ExtendedCommunityListsAPI.ExtendedcommunitylistsPatch(ctx).ExtendedcommunitylistsPutRequest(*request.(*openapi.ExtendedcommunitylistsPutRequest))
		return req.Execute()
	case "route_map_clause":
		req := g.client.RouteMapClausesAPI.RoutemapclausesPatch(ctx).RoutemapclausesPutRequest(*request.(*openapi.RoutemapclausesPutRequest))
		return req.Execute()
	case "route_map":
		req := g.client.RouteMapsAPI.RoutemapsPatch(ctx).RoutemapsPutRequest(*request.(*openapi.RoutemapsPutRequest))
		return req.Execute()
	case "sfp_breakout":
		req := g.client.SFPBreakoutsAPI.SfpbreakoutsPatch(ctx).SfpbreakoutsPatchRequest(*request.(*openapi.SfpbreakoutsPatchRequest))
		return req.Execute()
	case "site":
		req := g.client.SitesAPI.SitesPatch(ctx).SitesPatchRequest(*request.(*openapi.SitesPatchRequest))
		return req.Execute()
	case "pod":
		req := g.client.PodsAPI.PodsPatch(ctx).PodsPutRequest(*request.(*openapi.PodsPutRequest))
		return req.Execute()
	case "port_acl":
		req := g.client.PortACLsAPI.PortaclsPatch(ctx).PortaclsPutRequest(*request.(*openapi.PortaclsPutRequest))
		return req.Execute()
	case "sflow_collector":
		req := g.client.SFlowCollectorsAPI.SflowcollectorsPatch(ctx).SflowcollectorsPutRequest(*request.(*openapi.SflowcollectorsPutRequest))
		return req.Execute()
	case "diagnostics_profile":
		req := g.client.DiagnosticsProfilesAPI.DiagnosticsprofilesPatch(ctx).DiagnosticsprofilesPutRequest(*request.(*openapi.DiagnosticsprofilesPutRequest))
		return req.Execute()
	case "diagnostics_port_profile":
		req := g.client.DiagnosticsPortProfilesAPI.DiagnosticsportprofilesPatch(ctx).DiagnosticsportprofilesPutRequest(*request.(*openapi.DiagnosticsportprofilesPutRequest))
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
	case "device_settings":
		req := g.client.DeviceSettingsAPI.DevicesettingsDelete(ctx).EthDeviceProfilesName(names)
		return req.Execute()
	case "ipv4_list_filter":
		req := g.client.IPv4ListFiltersAPI.Ipv4listsDelete(ctx).Ipv4ListFilterName(names)
		return req.Execute()
	case "ipv4_prefix_list":
		req := g.client.IPv4PrefixListsAPI.Ipv4prefixlistsDelete(ctx).Ipv4PrefixListName(names)
		return req.Execute()
	case "ipv6_list_filter":
		req := g.client.IPv6ListFiltersAPI.Ipv6listsDelete(ctx).Ipv6ListFilterName(names)
		return req.Execute()
	case "ipv6_prefix_list":
		req := g.client.IPv6PrefixListsAPI.Ipv6prefixlistsDelete(ctx).Ipv6PrefixListName(names)
		return req.Execute()
	case "as_path_access_list":
		req := g.client.ASPathAccessListsAPI.AspathaccesslistsDelete(ctx).AsPathAccessListName(names)
		return req.Execute()
	case "community_list":
		req := g.client.CommunityListsAPI.CommunitylistsDelete(ctx).CommunityListName(names)
		return req.Execute()
	case "extended_community_list":
		req := g.client.ExtendedCommunityListsAPI.ExtendedcommunitylistsDelete(ctx).ExtendedCommunityListName(names)
		return req.Execute()
	case "route_map_clause":
		req := g.client.RouteMapClausesAPI.RoutemapclausesDelete(ctx).RouteMapClauseName(names)
		return req.Execute()
	case "route_map":
		req := g.client.RouteMapsAPI.RoutemapsDelete(ctx).RouteMapName(names)
		return req.Execute()
	case "pod":
		req := g.client.PodsAPI.PodsDelete(ctx).PodName(names)
		return req.Execute()
	case "port_acl":
		req := g.client.PortACLsAPI.PortaclsDelete(ctx).PortAclName(names)
		return req.Execute()
	case "sflow_collector":
		req := g.client.SFlowCollectorsAPI.SflowcollectorsDelete(ctx).SflowCollectorName(names)
		return req.Execute()
	case "diagnostics_profile":
		req := g.client.DiagnosticsProfilesAPI.DiagnosticsprofilesDelete(ctx).DiagnosticsProfileName(names)
		return req.Execute()
	case "diagnostics_port_profile":
		req := g.client.DiagnosticsPortProfilesAPI.DiagnosticsportprofilesDelete(ctx).DiagnosticsPortProfileName(names)
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
	case "device_settings":
		req := g.client.DeviceSettingsAPI.DevicesettingsGet(ctx)
		return req.Execute()
	case "ipv4_list_filter":
		req := g.client.IPv4ListFiltersAPI.Ipv4listsGet(ctx)
		return req.Execute()
	case "ipv4_prefix_list":
		req := g.client.IPv4PrefixListsAPI.Ipv4prefixlistsGet(ctx)
		return req.Execute()
	case "ipv6_list_filter":
		req := g.client.IPv6ListFiltersAPI.Ipv6listsGet(ctx)
		return req.Execute()
	case "ipv6_prefix_list":
		req := g.client.IPv6PrefixListsAPI.Ipv6prefixlistsGet(ctx)
		return req.Execute()
	case "as_path_access_list":
		req := g.client.ASPathAccessListsAPI.AspathaccesslistsGet(ctx)
		return req.Execute()
	case "community_list":
		req := g.client.CommunityListsAPI.CommunitylistsGet(ctx)
		return req.Execute()
	case "extended_community_list":
		req := g.client.ExtendedCommunityListsAPI.ExtendedcommunitylistsGet(ctx)
		return req.Execute()
	case "route_map_clause":
		req := g.client.RouteMapClausesAPI.RoutemapclausesGet(ctx)
		return req.Execute()
	case "route_map":
		req := g.client.RouteMapsAPI.RoutemapsGet(ctx)
		return req.Execute()
	case "sfp_breakout":
		req := g.client.SFPBreakoutsAPI.SfpbreakoutsGet(ctx)
		return req.Execute()
	case "site":
		req := g.client.SitesAPI.SitesGet(ctx)
		return req.Execute()
	case "pod":
		req := g.client.PodsAPI.PodsGet(ctx)
		return req.Execute()
	case "port_acl":
		req := g.client.PortACLsAPI.PortaclsGet(ctx)
		return req.Execute()
	case "sflow_collector":
		req := g.client.SFlowCollectorsAPI.SflowcollectorsGet(ctx)
		return req.Execute()
	case "diagnostics_profile":
		req := g.client.DiagnosticsProfilesAPI.DiagnosticsprofilesGet(ctx)
		return req.Execute()
	case "diagnostics_port_profile":
		req := g.client.DiagnosticsPortProfilesAPI.DiagnosticsportprofilesGet(ctx)
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
func (b *BulkOperationManager) AddAclPut(ctx context.Context, aclName string, props openapi.AclsPutRequestIpFilterValue, ipVersion string) string {
	return b.addGenericOperation(ctx, "acl", aclName, "PUT", props, ipVersion)
}

func (b *BulkOperationManager) AddAclPatch(ctx context.Context, aclName string, props openapi.AclsPutRequestIpFilterValue, ipVersion string) string {
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
	var originalOperations map[string]openapi.AclsPutRequestIpFilterValue
	var originalIpVersions map[string]string

	switch operationType {
	case "PUT":
		originalOperations = make(map[string]openapi.AclsPutRequestIpFilterValue)
		for k, v := range b.aclPut {
			originalOperations[k] = v
		}
		b.aclPut = make(map[string]openapi.AclsPutRequestIpFilterValue)
	case "PATCH":
		originalOperations = make(map[string]openapi.AclsPutRequestIpFilterValue)
		for k, v := range b.aclPatch {
			originalOperations[k] = v
		}
		b.aclPatch = make(map[string]openapi.AclsPutRequestIpFilterValue)
	case "DELETE":
		originalOperations = make(map[string]openapi.AclsPutRequestIpFilterValue)
		for _, name := range b.aclDelete {
			originalOperations[name] = openapi.AclsPutRequestIpFilterValue{}
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

	ipv4Data := make(map[string]openapi.AclsPutRequestIpFilterValue)
	ipv6Data := make(map[string]openapi.AclsPutRequestIpFilterValue)

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

func (b *BulkOperationManager) executeAclForIpVersion(ctx context.Context, aclData map[string]openapi.AclsPutRequestIpFilterValue, operationType, ipVersion string) diag.Diagnostics {
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
				aclMap := make(map[string]openapi.AclsPutRequestIpFilterValue)
				for name, props := range filteredData {
					aclMap[name] = props.(openapi.AclsPutRequestIpFilterValue)
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
				b.gatewayPut = make(map[string]openapi.GatewaysPutRequestGatewayValue)
			}
			b.gatewayPut[resourceName] = props.(openapi.GatewaysPutRequestGatewayValue)
		} else {
			if b.gatewayPatch == nil {
				b.gatewayPatch = make(map[string]openapi.GatewaysPutRequestGatewayValue)
			}
			b.gatewayPatch[resourceName] = props.(openapi.GatewaysPutRequestGatewayValue)
		}
	case "lag":
		if operationType == "PUT" {
			if b.lagPut == nil {
				b.lagPut = make(map[string]openapi.LagsPutRequestLagValue)
			}
			b.lagPut[resourceName] = props.(openapi.LagsPutRequestLagValue)
		} else {
			if b.lagPatch == nil {
				b.lagPatch = make(map[string]openapi.LagsPutRequestLagValue)
			}
			b.lagPatch[resourceName] = props.(openapi.LagsPutRequestLagValue)
		}
	case "tenant":
		if operationType == "PUT" {
			if b.tenantPut == nil {
				b.tenantPut = make(map[string]openapi.TenantsPutRequestTenantValue)
			}
			b.tenantPut[resourceName] = props.(openapi.TenantsPutRequestTenantValue)
		} else {
			if b.tenantPatch == nil {
				b.tenantPatch = make(map[string]openapi.TenantsPutRequestTenantValue)
			}
			b.tenantPatch[resourceName] = props.(openapi.TenantsPutRequestTenantValue)
		}
	case "service":
		if operationType == "PUT" {
			if b.servicePut == nil {
				b.servicePut = make(map[string]openapi.ServicesPutRequestServiceValue)
			}
			b.servicePut[resourceName] = props.(openapi.ServicesPutRequestServiceValue)
		} else {
			if b.servicePatch == nil {
				b.servicePatch = make(map[string]openapi.ServicesPutRequestServiceValue)
			}
			b.servicePatch[resourceName] = props.(openapi.ServicesPutRequestServiceValue)
		}
	case "gateway_profile":
		if operationType == "PUT" {
			if b.gatewayProfilePut == nil {
				b.gatewayProfilePut = make(map[string]openapi.GatewayprofilesPutRequestGatewayProfileValue)
			}
			b.gatewayProfilePut[resourceName] = props.(openapi.GatewayprofilesPutRequestGatewayProfileValue)
		} else {
			if b.gatewayProfilePatch == nil {
				b.gatewayProfilePatch = make(map[string]openapi.GatewayprofilesPutRequestGatewayProfileValue)
			}
			b.gatewayProfilePatch[resourceName] = props.(openapi.GatewayprofilesPutRequestGatewayProfileValue)
		}
	case "eth_port_profile":
		if operationType == "PUT" {
			if b.ethPortProfilePut == nil {
				b.ethPortProfilePut = make(map[string]openapi.EthportprofilesPutRequestEthPortProfileValue)
			}
			b.ethPortProfilePut[resourceName] = props.(openapi.EthportprofilesPutRequestEthPortProfileValue)
		} else {
			if b.ethPortProfilePatch == nil {
				b.ethPortProfilePatch = make(map[string]openapi.EthportprofilesPutRequestEthPortProfileValue)
			}
			b.ethPortProfilePatch[resourceName] = props.(openapi.EthportprofilesPutRequestEthPortProfileValue)
		}
	case "eth_port_settings":
		if operationType == "PUT" {
			if b.ethPortSettingsPut == nil {
				b.ethPortSettingsPut = make(map[string]openapi.EthportsettingsPutRequestEthPortSettingsValue)
			}
			b.ethPortSettingsPut[resourceName] = props.(openapi.EthportsettingsPutRequestEthPortSettingsValue)
		} else {
			if b.ethPortSettingsPatch == nil {
				b.ethPortSettingsPatch = make(map[string]openapi.EthportsettingsPutRequestEthPortSettingsValue)
			}
			b.ethPortSettingsPatch[resourceName] = props.(openapi.EthportsettingsPutRequestEthPortSettingsValue)
		}
	case "device_settings":
		if operationType == "PUT" {
			if b.deviceSettingsPut == nil {
				b.deviceSettingsPut = make(map[string]openapi.DevicesettingsPutRequestEthDeviceProfilesValue)
			}
			b.deviceSettingsPut[resourceName] = props.(openapi.DevicesettingsPutRequestEthDeviceProfilesValue)
		} else {
			if b.deviceSettingsPatch == nil {
				b.deviceSettingsPatch = make(map[string]openapi.DevicesettingsPutRequestEthDeviceProfilesValue)
			}
			b.deviceSettingsPatch[resourceName] = props.(openapi.DevicesettingsPutRequestEthDeviceProfilesValue)
		}
	case "bundle":
		if operationType == "PUT" {
			if b.bundlePut == nil {
				b.bundlePut = make(map[string]openapi.BundlesPutRequestEndpointBundleValue)
			}
			b.bundlePut[resourceName] = props.(openapi.BundlesPutRequestEndpointBundleValue)
		} else {
			if b.bundlePatch == nil {
				b.bundlePatch = make(map[string]openapi.BundlesPutRequestEndpointBundleValue)
			}
			b.bundlePatch[resourceName] = props.(openapi.BundlesPutRequestEndpointBundleValue)
		}
	case "acl":
		if operationType == "PUT" {
			if b.aclPut == nil {
				b.aclPut = make(map[string]openapi.AclsPutRequestIpFilterValue)
			}
			b.aclPut[resourceName] = props.(openapi.AclsPutRequestIpFilterValue)
		} else {
			if b.aclPatch == nil {
				b.aclPatch = make(map[string]openapi.AclsPutRequestIpFilterValue)
			}
			b.aclPatch[resourceName] = props.(openapi.AclsPutRequestIpFilterValue)
		}
	case "ipv4_list_filter":
		if operationType == "PUT" {
			if b.ipv4ListPut == nil {
				b.ipv4ListPut = make(map[string]openapi.Ipv4listsPutRequestIpv4ListFilterValue)
			}
			b.ipv4ListPut[resourceName] = props.(openapi.Ipv4listsPutRequestIpv4ListFilterValue)
		} else {
			if b.ipv4ListPatch == nil {
				b.ipv4ListPatch = make(map[string]openapi.Ipv4listsPutRequestIpv4ListFilterValue)
			}
			b.ipv4ListPatch[resourceName] = props.(openapi.Ipv4listsPutRequestIpv4ListFilterValue)
		}
	case "ipv4_prefix_list":
		if operationType == "PUT" {
			if b.ipv4PrefixListPut == nil {
				b.ipv4PrefixListPut = make(map[string]openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValue)
			}
			b.ipv4PrefixListPut[resourceName] = props.(openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValue)
		} else {
			if b.ipv4PrefixListPatch == nil {
				b.ipv4PrefixListPatch = make(map[string]openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValue)
			}
			b.ipv4PrefixListPatch[resourceName] = props.(openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValue)
		}
	case "ipv6_list_filter":
		if operationType == "PUT" {
			if b.ipv6ListPut == nil {
				b.ipv6ListPut = make(map[string]openapi.Ipv6listsPutRequestIpv6ListFilterValue)
			}
			b.ipv6ListPut[resourceName] = props.(openapi.Ipv6listsPutRequestIpv6ListFilterValue)
		} else {
			if b.ipv6ListPatch == nil {
				b.ipv6ListPatch = make(map[string]openapi.Ipv6listsPutRequestIpv6ListFilterValue)
			}
			b.ipv6ListPatch[resourceName] = props.(openapi.Ipv6listsPutRequestIpv6ListFilterValue)
		}
	case "ipv6_prefix_list":
		if operationType == "PUT" {
			if b.ipv6PrefixListPut == nil {
				b.ipv6PrefixListPut = make(map[string]openapi.Ipv6prefixlistsPutRequestIpv6PrefixListValue)
			}
			b.ipv6PrefixListPut[resourceName] = props.(openapi.Ipv6prefixlistsPutRequestIpv6PrefixListValue)
		} else {
			if b.ipv6PrefixListPatch == nil {
				b.ipv6PrefixListPatch = make(map[string]openapi.Ipv6prefixlistsPutRequestIpv6PrefixListValue)
			}
			b.ipv6PrefixListPatch[resourceName] = props.(openapi.Ipv6prefixlistsPutRequestIpv6PrefixListValue)
		}
	case "authenticated_eth_port":
		if operationType == "PUT" {
			if b.authenticatedEthPortPut == nil {
				b.authenticatedEthPortPut = make(map[string]openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValue)
			}
			b.authenticatedEthPortPut[resourceName] = props.(openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValue)
		} else {
			if b.authenticatedEthPortPatch == nil {
				b.authenticatedEthPortPatch = make(map[string]openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValue)
			}
			b.authenticatedEthPortPatch[resourceName] = props.(openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValue)
		}
	case "badge":
		if operationType == "PUT" {
			if b.badgePut == nil {
				b.badgePut = make(map[string]openapi.BadgesPutRequestBadgeValue)
			}
			b.badgePut[resourceName] = props.(openapi.BadgesPutRequestBadgeValue)
		} else {
			if b.badgePatch == nil {
				b.badgePatch = make(map[string]openapi.BadgesPutRequestBadgeValue)
			}
			b.badgePatch[resourceName] = props.(openapi.BadgesPutRequestBadgeValue)
		}
	case "device_voice_settings":
		if operationType == "PUT" {
			if b.deviceVoiceSettingsPut == nil {
				b.deviceVoiceSettingsPut = make(map[string]openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValue)
			}
			b.deviceVoiceSettingsPut[resourceName] = props.(openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValue)
		} else {
			if b.deviceVoiceSettingsPatch == nil {
				b.deviceVoiceSettingsPatch = make(map[string]openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValue)
			}
			b.deviceVoiceSettingsPatch[resourceName] = props.(openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValue)
		}
	case "as_path_access_list":
		if operationType == "PUT" {
			if b.asPathAccessListPut == nil {
				b.asPathAccessListPut = make(map[string]openapi.AspathaccesslistsPutRequestAsPathAccessListValue)
			}
			b.asPathAccessListPut[resourceName] = props.(openapi.AspathaccesslistsPutRequestAsPathAccessListValue)
		} else {
			if b.asPathAccessListPatch == nil {
				b.asPathAccessListPatch = make(map[string]openapi.AspathaccesslistsPutRequestAsPathAccessListValue)
			}
			b.asPathAccessListPatch[resourceName] = props.(openapi.AspathaccesslistsPutRequestAsPathAccessListValue)
		}
	case "community_list":
		if operationType == "PUT" {
			if b.communityListPut == nil {
				b.communityListPut = make(map[string]openapi.CommunitylistsPutRequestCommunityListValue)
			}
			b.communityListPut[resourceName] = props.(openapi.CommunitylistsPutRequestCommunityListValue)
		} else {
			if b.communityListPatch == nil {
				b.communityListPatch = make(map[string]openapi.CommunitylistsPutRequestCommunityListValue)
			}
			b.communityListPatch[resourceName] = props.(openapi.CommunitylistsPutRequestCommunityListValue)
		}
	case "extended_community_list":
		if operationType == "PUT" {
			if b.extendedCommunityListPut == nil {
				b.extendedCommunityListPut = make(map[string]openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValue)
			}
			b.extendedCommunityListPut[resourceName] = props.(openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValue)
		} else {
			if b.extendedCommunityListPatch == nil {
				b.extendedCommunityListPatch = make(map[string]openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValue)
			}
			b.extendedCommunityListPatch[resourceName] = props.(openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValue)
		}
	case "route_map_clause":
		if operationType == "PUT" {
			if b.routeMapClausePut == nil {
				b.routeMapClausePut = make(map[string]openapi.RoutemapclausesPutRequestRouteMapClauseValue)
			}
			b.routeMapClausePut[resourceName] = props.(openapi.RoutemapclausesPutRequestRouteMapClauseValue)
		} else {
			if b.routeMapClausePatch == nil {
				b.routeMapClausePatch = make(map[string]openapi.RoutemapclausesPutRequestRouteMapClauseValue)
			}
			b.routeMapClausePatch[resourceName] = props.(openapi.RoutemapclausesPutRequestRouteMapClauseValue)
		}
	case "route_map":
		if operationType == "PUT" {
			if b.routeMapPut == nil {
				b.routeMapPut = make(map[string]openapi.RoutemapsPutRequestRouteMapValue)
			}
			b.routeMapPut[resourceName] = props.(openapi.RoutemapsPutRequestRouteMapValue)
		} else {
			if b.routeMapPatch == nil {
				b.routeMapPatch = make(map[string]openapi.RoutemapsPutRequestRouteMapValue)
			}
			b.routeMapPatch[resourceName] = props.(openapi.RoutemapsPutRequestRouteMapValue)
		}
	case "sfp_breakout":
		// SFP Breakouts only support PATCH operations
		if operationType == "PATCH" {
			if b.sfpBreakoutPatch == nil {
				b.sfpBreakoutPatch = make(map[string]openapi.SfpbreakoutsPatchRequestSfpBreakoutsValue)
			}
			b.sfpBreakoutPatch[resourceName] = props.(openapi.SfpbreakoutsPatchRequestSfpBreakoutsValue)
		}
	case "site":
		// Sites only support PATCH operations
		if operationType == "PATCH" {
			if b.sitePatch == nil {
				b.sitePatch = make(map[string]openapi.SitesPatchRequestSiteValue)
			}
			b.sitePatch[resourceName] = props.(openapi.SitesPatchRequestSiteValue)
		}
	case "packet_broker":
		if operationType == "PUT" {
			if b.packetBrokerPut == nil {
				b.packetBrokerPut = make(map[string]openapi.PacketbrokerPutRequestPbEgressProfileValue)
			}
			b.packetBrokerPut[resourceName] = props.(openapi.PacketbrokerPutRequestPbEgressProfileValue)
		} else {
			if b.packetBrokerPatch == nil {
				b.packetBrokerPatch = make(map[string]openapi.PacketbrokerPutRequestPbEgressProfileValue)
			}
			b.packetBrokerPatch[resourceName] = props.(openapi.PacketbrokerPutRequestPbEgressProfileValue)
		}
	case "packet_queue":
		if operationType == "PUT" {
			if b.packetQueuePut == nil {
				b.packetQueuePut = make(map[string]openapi.PacketqueuesPutRequestPacketQueueValue)
			}
			b.packetQueuePut[resourceName] = props.(openapi.PacketqueuesPutRequestPacketQueueValue)
		} else {
			if b.packetQueuePatch == nil {
				b.packetQueuePatch = make(map[string]openapi.PacketqueuesPutRequestPacketQueueValue)
			}
			b.packetQueuePatch[resourceName] = props.(openapi.PacketqueuesPutRequestPacketQueueValue)
		}
	case "service_port_profile":
		if operationType == "PUT" {
			if b.servicePortProfilePut == nil {
				b.servicePortProfilePut = make(map[string]openapi.ServiceportprofilesPutRequestServicePortProfileValue)
			}
			b.servicePortProfilePut[resourceName] = props.(openapi.ServiceportprofilesPutRequestServicePortProfileValue)
		} else {
			if b.servicePortProfilePatch == nil {
				b.servicePortProfilePatch = make(map[string]openapi.ServiceportprofilesPutRequestServicePortProfileValue)
			}
			b.servicePortProfilePatch[resourceName] = props.(openapi.ServiceportprofilesPutRequestServicePortProfileValue)
		}
	case "switchpoint":
		if operationType == "PUT" {
			if b.switchpointPut == nil {
				b.switchpointPut = make(map[string]openapi.SwitchpointsPutRequestSwitchpointValue)
			}
			b.switchpointPut[resourceName] = props.(openapi.SwitchpointsPutRequestSwitchpointValue)
		} else {
			if b.switchpointPatch == nil {
				b.switchpointPatch = make(map[string]openapi.SwitchpointsPutRequestSwitchpointValue)
			}
			b.switchpointPatch[resourceName] = props.(openapi.SwitchpointsPutRequestSwitchpointValue)
		}
	case "device_controller":
		if operationType == "PUT" {
			if b.deviceControllerPut == nil {
				b.deviceControllerPut = make(map[string]openapi.DevicecontrollersPutRequestDeviceControllerValue)
			}
			b.deviceControllerPut[resourceName] = props.(openapi.DevicecontrollersPutRequestDeviceControllerValue)
		} else {
			if b.deviceControllerPatch == nil {
				b.deviceControllerPatch = make(map[string]openapi.DevicecontrollersPutRequestDeviceControllerValue)
			}
			b.deviceControllerPatch[resourceName] = props.(openapi.DevicecontrollersPutRequestDeviceControllerValue)
		}
	case "voice_port_profile":
		if operationType == "PUT" {
			if b.voicePortProfilePut == nil {
				b.voicePortProfilePut = make(map[string]openapi.VoiceportprofilesPutRequestVoicePortProfilesValue)
			}
			b.voicePortProfilePut[resourceName] = props.(openapi.VoiceportprofilesPutRequestVoicePortProfilesValue)
		} else {
			if b.voicePortProfilePatch == nil {
				b.voicePortProfilePatch = make(map[string]openapi.VoiceportprofilesPutRequestVoicePortProfilesValue)
			}
			b.voicePortProfilePatch[resourceName] = props.(openapi.VoiceportprofilesPutRequestVoicePortProfilesValue)
		}
	case "pod":
		if operationType == "PUT" {
			if b.podPut == nil {
				b.podPut = make(map[string]openapi.PodsPutRequestPodValue)
			}
			b.podPut[resourceName] = props.(openapi.PodsPutRequestPodValue)
		} else {
			if b.podPatch == nil {
				b.podPatch = make(map[string]openapi.PodsPutRequestPodValue)
			}
			b.podPatch[resourceName] = props.(openapi.PodsPutRequestPodValue)
		}
	case "port_acl":
		if operationType == "PUT" {
			if b.portAclPut == nil {
				b.portAclPut = make(map[string]openapi.PortaclsPutRequestPortAclValue)
			}
			b.portAclPut[resourceName] = props.(openapi.PortaclsPutRequestPortAclValue)
		} else {
			if b.portAclPatch == nil {
				b.portAclPatch = make(map[string]openapi.PortaclsPutRequestPortAclValue)
			}
			b.portAclPatch[resourceName] = props.(openapi.PortaclsPutRequestPortAclValue)
		}
	case "sflow_collector":
		if operationType == "PUT" {
			if b.sflowCollectorPut == nil {
				b.sflowCollectorPut = make(map[string]openapi.SflowcollectorsPutRequestSflowCollectorValue)
			}
			b.sflowCollectorPut[resourceName] = props.(openapi.SflowcollectorsPutRequestSflowCollectorValue)
		} else {
			if b.sflowCollectorPatch == nil {
				b.sflowCollectorPatch = make(map[string]openapi.SflowcollectorsPutRequestSflowCollectorValue)
			}
			b.sflowCollectorPatch[resourceName] = props.(openapi.SflowcollectorsPutRequestSflowCollectorValue)
		}
	case "diagnostics_profile":
		if operationType == "PUT" {
			if b.diagnosticsProfilePut == nil {
				b.diagnosticsProfilePut = make(map[string]openapi.DiagnosticsprofilesPutRequestDiagnosticsProfileValue)
			}
			b.diagnosticsProfilePut[resourceName] = props.(openapi.DiagnosticsprofilesPutRequestDiagnosticsProfileValue)
		} else {
			if b.diagnosticsProfilePatch == nil {
				b.diagnosticsProfilePatch = make(map[string]openapi.DiagnosticsprofilesPutRequestDiagnosticsProfileValue)
			}
			b.diagnosticsProfilePatch[resourceName] = props.(openapi.DiagnosticsprofilesPutRequestDiagnosticsProfileValue)
		}
	case "diagnostics_port_profile":
		if operationType == "PUT" {
			if b.diagnosticsPortProfilePut == nil {
				b.diagnosticsPortProfilePut = make(map[string]openapi.DiagnosticsportprofilesPutRequestDiagnosticsPortProfileValue)
			}
			b.diagnosticsPortProfilePut[resourceName] = props.(openapi.DiagnosticsportprofilesPutRequestDiagnosticsPortProfileValue)
		} else {
			if b.diagnosticsPortProfilePatch == nil {
				b.diagnosticsPortProfilePatch = make(map[string]openapi.DiagnosticsportprofilesPutRequestDiagnosticsPortProfileValue)
			}
			b.diagnosticsPortProfilePatch[resourceName] = props.(openapi.DiagnosticsportprofilesPutRequestDiagnosticsPortProfileValue)
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
	case "device_settings":
		return &b.deviceSettingsDelete
	case "bundle":
		return &b.bundleDelete
	case "acl":
		return &b.aclDelete
	case "ipv4_list_filter":
		return &b.ipv4ListDelete
	case "ipv4_prefix_list":
		return &b.ipv4PrefixListDelete
	case "ipv6_list_filter":
		return &b.ipv6ListDelete
	case "ipv6_prefix_list":
		return &b.ipv6PrefixListDelete
	case "authenticated_eth_port":
		return &b.authenticatedEthPortDelete
	case "badge":
		return &b.badgeDelete
	case "device_voice_settings":
		return &b.deviceVoiceSettingsDelete
	case "as_path_access_list":
		return &b.asPathAccessListDelete
	case "community_list":
		return &b.communityListDelete
	case "extended_community_list":
		return &b.extendedCommunityListDelete
	case "route_map_clause":
		return &b.routeMapClauseDelete
	case "route_map":
		return &b.routeMapDelete
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
	case "pod":
		return &b.podDelete
	case "port_acl":
		return &b.portAclDelete
	case "sflow_collector":
		return &b.sflowCollectorDelete
	case "diagnostics_profile":
		return &b.diagnosticsProfileDelete
	case "diagnostics_port_profile":
		return &b.diagnosticsPortProfileDelete
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
		if data.DeleteOperations != nil {
			return len(*data.DeleteOperations)
		}
		return 0
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
	case "device_settings":
		if operationType == "PUT" {
			result := make(map[string]interface{})
			for k, v := range b.deviceSettingsPut {
				result[k] = v
			}
			return result
		} else {
			result := make(map[string]interface{})
			for k, v := range b.deviceSettingsPatch {
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
	case "ipv4_list_filter":
		if operationType == "PUT" {
			result := make(map[string]interface{})
			for k, v := range b.ipv4ListPut {
				result[k] = v
			}
			return result
		} else {
			result := make(map[string]interface{})
			for k, v := range b.ipv4ListPatch {
				result[k] = v
			}
			return result
		}
	case "ipv4_prefix_list":
		if operationType == "PUT" {
			result := make(map[string]interface{})
			for k, v := range b.ipv4PrefixListPut {
				result[k] = v
			}
			return result
		} else {
			result := make(map[string]interface{})
			for k, v := range b.ipv4PrefixListPatch {
				result[k] = v
			}
			return result
		}
	case "ipv6_list_filter":
		if operationType == "PUT" {
			result := make(map[string]interface{})
			for k, v := range b.ipv6ListPut {
				result[k] = v
			}
			return result
		} else {
			result := make(map[string]interface{})
			for k, v := range b.ipv6ListPatch {
				result[k] = v
			}
			return result
		}
	case "ipv6_prefix_list":
		if operationType == "PUT" {
			result := make(map[string]interface{})
			for k, v := range b.ipv6PrefixListPut {
				result[k] = v
			}
			return result
		} else {
			result := make(map[string]interface{})
			for k, v := range b.ipv6PrefixListPatch {
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
	case "as_path_access_list":
		if operationType == "PUT" {
			result := make(map[string]interface{})
			for k, v := range b.asPathAccessListPut {
				result[k] = v
			}
			return result
		} else {
			result := make(map[string]interface{})
			for k, v := range b.asPathAccessListPatch {
				result[k] = v
			}
			return result
		}
	case "community_list":
		if operationType == "PUT" {
			result := make(map[string]interface{})
			for k, v := range b.communityListPut {
				result[k] = v
			}
			return result
		} else {
			result := make(map[string]interface{})
			for k, v := range b.communityListPatch {
				result[k] = v
			}
			return result
		}
	case "extended_community_list":
		if operationType == "PUT" {
			result := make(map[string]interface{})
			for k, v := range b.extendedCommunityListPut {
				result[k] = v
			}
			return result
		} else {
			result := make(map[string]interface{})
			for k, v := range b.extendedCommunityListPatch {
				result[k] = v
			}
			return result
		}
	case "route_map_clause":
		if operationType == "PUT" {
			result := make(map[string]interface{})
			for k, v := range b.routeMapClausePut {
				result[k] = v
			}
			return result
		} else {
			result := make(map[string]interface{})
			for k, v := range b.routeMapClausePatch {
				result[k] = v
			}
			return result
		}
	case "route_map":
		if operationType == "PUT" {
			result := make(map[string]interface{})
			for k, v := range b.routeMapPut {
				result[k] = v
			}
			return result
		} else {
			result := make(map[string]interface{})
			for k, v := range b.routeMapPatch {
				result[k] = v
			}
			return result
		}
	case "sfp_breakout":
		// SFP Breakouts only support PATCH operations
		if operationType == "PATCH" {
			result := make(map[string]interface{})
			for k, v := range b.sfpBreakoutPatch {
				result[k] = v
			}
			return result
		}
		return make(map[string]interface{})
	case "site":
		// Sites only support PATCH operations
		if operationType == "PATCH" {
			result := make(map[string]interface{})
			for k, v := range b.sitePatch {
				result[k] = v
			}
			return result
		}
		return make(map[string]interface{})
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
	case "pod":
		if operationType == "PUT" {
			result := make(map[string]interface{})
			for k, v := range b.podPut {
				result[k] = v
			}
			return result
		} else {
			result := make(map[string]interface{})
			for k, v := range b.podPatch {
				result[k] = v
			}
			return result
		}
	case "port_acl":
		if operationType == "PUT" {
			result := make(map[string]interface{})
			for k, v := range b.portAclPut {
				result[k] = v
			}
			return result
		} else {
			result := make(map[string]interface{})
			for k, v := range b.portAclPatch {
				result[k] = v
			}
			return result
		}
	case "sflow_collector":
		if operationType == "PUT" {
			result := make(map[string]interface{})
			for k, v := range b.sflowCollectorPut {
				result[k] = v
			}
			return result
		} else {
			result := make(map[string]interface{})
			for k, v := range b.sflowCollectorPatch {
				result[k] = v
			}
			return result
		}
	case "diagnostics_profile":
		if operationType == "PUT" {
			result := make(map[string]interface{})
			for k, v := range b.diagnosticsProfilePut {
				result[k] = v
			}
			return result
		} else {
			result := make(map[string]interface{})
			for k, v := range b.diagnosticsProfilePatch {
				result[k] = v
			}
			return result
		}
	case "diagnostics_port_profile":
		if operationType == "PUT" {
			result := make(map[string]interface{})
			for k, v := range b.diagnosticsPortProfilePut {
				result[k] = v
			}
			return result
		} else {
			result := make(map[string]interface{})
			for k, v := range b.diagnosticsPortProfilePatch {
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
			b.gatewayPut = make(map[string]openapi.GatewaysPutRequestGatewayValue)
		} else {
			b.gatewayPatch = make(map[string]openapi.GatewaysPutRequestGatewayValue)
		}
	case "lag":
		if operationType == "PUT" {
			b.lagPut = make(map[string]openapi.LagsPutRequestLagValue)
		} else {
			b.lagPatch = make(map[string]openapi.LagsPutRequestLagValue)
		}
	case "tenant":
		if operationType == "PUT" {
			b.tenantPut = make(map[string]openapi.TenantsPutRequestTenantValue)
		} else {
			b.tenantPatch = make(map[string]openapi.TenantsPutRequestTenantValue)
		}
	case "service":
		if operationType == "PUT" {
			b.servicePut = make(map[string]openapi.ServicesPutRequestServiceValue)
		} else {
			b.servicePatch = make(map[string]openapi.ServicesPutRequestServiceValue)
		}
	case "gateway_profile":
		if operationType == "PUT" {
			b.gatewayProfilePut = make(map[string]openapi.GatewayprofilesPutRequestGatewayProfileValue)
		} else {
			b.gatewayProfilePatch = make(map[string]openapi.GatewayprofilesPutRequestGatewayProfileValue)
		}
	case "eth_port_profile":
		if operationType == "PUT" {
			b.ethPortProfilePut = make(map[string]openapi.EthportprofilesPutRequestEthPortProfileValue)
		} else {
			b.ethPortProfilePatch = make(map[string]openapi.EthportprofilesPutRequestEthPortProfileValue)
		}
	case "eth_port_settings":
		if operationType == "PUT" {
			b.ethPortSettingsPut = make(map[string]openapi.EthportsettingsPutRequestEthPortSettingsValue)
		} else {
			b.ethPortSettingsPatch = make(map[string]openapi.EthportsettingsPutRequestEthPortSettingsValue)
		}
	case "device_settings":
		if operationType == "PUT" {
			b.deviceSettingsPut = make(map[string]openapi.DevicesettingsPutRequestEthDeviceProfilesValue)
		} else {
			b.deviceSettingsPatch = make(map[string]openapi.DevicesettingsPutRequestEthDeviceProfilesValue)
		}
	case "bundle":
		if operationType == "PUT" {
			b.bundlePut = make(map[string]openapi.BundlesPutRequestEndpointBundleValue)
		} else {
			b.bundlePatch = make(map[string]openapi.BundlesPutRequestEndpointBundleValue)
		}
	case "acl":
		if operationType == "PUT" {
			b.aclPut = make(map[string]openapi.AclsPutRequestIpFilterValue)
		} else {
			b.aclPatch = make(map[string]openapi.AclsPutRequestIpFilterValue)
		}
	case "ipv4_list_filter":
		if operationType == "PUT" {
			b.ipv4ListPut = make(map[string]openapi.Ipv4listsPutRequestIpv4ListFilterValue)
		} else {
			b.ipv4ListPatch = make(map[string]openapi.Ipv4listsPutRequestIpv4ListFilterValue)
		}
	case "ipv4_prefix_list":
		if operationType == "PUT" {
			b.ipv4PrefixListPut = make(map[string]openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValue)
		} else {
			b.ipv4PrefixListPatch = make(map[string]openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValue)
		}
	case "ipv6_list_filter":
		if operationType == "PUT" {
			b.ipv6ListPut = make(map[string]openapi.Ipv6listsPutRequestIpv6ListFilterValue)
		} else {
			b.ipv6ListPatch = make(map[string]openapi.Ipv6listsPutRequestIpv6ListFilterValue)
		}
	case "ipv6_prefix_list":
		if operationType == "PUT" {
			b.ipv6PrefixListPut = make(map[string]openapi.Ipv6prefixlistsPutRequestIpv6PrefixListValue)
		} else {
			b.ipv6PrefixListPatch = make(map[string]openapi.Ipv6prefixlistsPutRequestIpv6PrefixListValue)
		}
	case "authenticated_eth_port":
		if operationType == "PUT" {
			b.authenticatedEthPortPut = make(map[string]openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValue)
		} else {
			b.authenticatedEthPortPatch = make(map[string]openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValue)
		}
	case "badge":
		if operationType == "PUT" {
			b.badgePut = make(map[string]openapi.BadgesPutRequestBadgeValue)
		} else {
			b.badgePatch = make(map[string]openapi.BadgesPutRequestBadgeValue)
		}
	case "device_voice_settings":
		if operationType == "PUT" {
			b.deviceVoiceSettingsPut = make(map[string]openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValue)
		} else {
			b.deviceVoiceSettingsPatch = make(map[string]openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValue)
		}
	case "as_path_access_list":
		if operationType == "PUT" {
			b.asPathAccessListPut = make(map[string]openapi.AspathaccesslistsPutRequestAsPathAccessListValue)
		} else {
			b.asPathAccessListPatch = make(map[string]openapi.AspathaccesslistsPutRequestAsPathAccessListValue)
		}
	case "community_list":
		if operationType == "PUT" {
			b.communityListPut = make(map[string]openapi.CommunitylistsPutRequestCommunityListValue)
		} else {
			b.communityListPatch = make(map[string]openapi.CommunitylistsPutRequestCommunityListValue)
		}
	case "extended_community_list":
		if operationType == "PUT" {
			b.extendedCommunityListPut = make(map[string]openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValue)
		} else {
			b.extendedCommunityListPatch = make(map[string]openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValue)
		}
	case "route_map_clause":
		if operationType == "PUT" {
			b.routeMapClausePut = make(map[string]openapi.RoutemapclausesPutRequestRouteMapClauseValue)
		} else {
			b.routeMapClausePatch = make(map[string]openapi.RoutemapclausesPutRequestRouteMapClauseValue)
		}
	case "route_map":
		if operationType == "PUT" {
			b.routeMapPut = make(map[string]openapi.RoutemapsPutRequestRouteMapValue)
		} else {
			b.routeMapPatch = make(map[string]openapi.RoutemapsPutRequestRouteMapValue)
		}
	case "sfp_breakout":
		// SFP Breakouts only support PATCH operations
		if operationType == "PATCH" {
			b.sfpBreakoutPatch = make(map[string]openapi.SfpbreakoutsPatchRequestSfpBreakoutsValue)
		}
	case "site":
		// Sites only support PATCH operations
		if operationType == "PATCH" {
			b.sitePatch = make(map[string]openapi.SitesPatchRequestSiteValue)
		}
	case "packet_broker":
		if operationType == "PUT" {
			b.packetBrokerPut = make(map[string]openapi.PacketbrokerPutRequestPbEgressProfileValue)
		} else {
			b.packetBrokerPatch = make(map[string]openapi.PacketbrokerPutRequestPbEgressProfileValue)
		}
	case "packet_queue":
		if operationType == "PUT" {
			b.packetQueuePut = make(map[string]openapi.PacketqueuesPutRequestPacketQueueValue)
		} else {
			b.packetQueuePatch = make(map[string]openapi.PacketqueuesPutRequestPacketQueueValue)
		}
	case "service_port_profile":
		if operationType == "PUT" {
			b.servicePortProfilePut = make(map[string]openapi.ServiceportprofilesPutRequestServicePortProfileValue)
		} else {
			b.servicePortProfilePatch = make(map[string]openapi.ServiceportprofilesPutRequestServicePortProfileValue)
		}
	case "switchpoint":
		if operationType == "PUT" {
			b.switchpointPut = make(map[string]openapi.SwitchpointsPutRequestSwitchpointValue)
		} else {
			b.switchpointPatch = make(map[string]openapi.SwitchpointsPutRequestSwitchpointValue)
		}
	case "voice_port_profile":
		if operationType == "PUT" {
			b.voicePortProfilePut = make(map[string]openapi.VoiceportprofilesPutRequestVoicePortProfilesValue)
		} else {
			b.voicePortProfilePatch = make(map[string]openapi.VoiceportprofilesPutRequestVoicePortProfilesValue)
		}
	case "device_controller":
		if operationType == "PUT" {
			b.deviceControllerPut = make(map[string]openapi.DevicecontrollersPutRequestDeviceControllerValue)
		} else {
			b.deviceControllerPatch = make(map[string]openapi.DevicecontrollersPutRequestDeviceControllerValue)
		}
	case "pod":
		if operationType == "PUT" {
			b.podPut = make(map[string]openapi.PodsPutRequestPodValue)
		} else {
			b.podPatch = make(map[string]openapi.PodsPutRequestPodValue)
		}
	case "port_acl":
		if operationType == "PUT" {
			b.portAclPut = make(map[string]openapi.PortaclsPutRequestPortAclValue)
		} else {
			b.portAclPatch = make(map[string]openapi.PortaclsPutRequestPortAclValue)
		}
	case "sflow_collector":
		if operationType == "PUT" {
			b.sflowCollectorPut = make(map[string]openapi.SflowcollectorsPutRequestSflowCollectorValue)
		} else {
			b.sflowCollectorPatch = make(map[string]openapi.SflowcollectorsPutRequestSflowCollectorValue)
		}
	case "diagnostics_profile":
		if operationType == "PUT" {
			b.diagnosticsProfilePut = make(map[string]openapi.DiagnosticsprofilesPutRequestDiagnosticsProfileValue)
		} else {
			b.diagnosticsProfilePatch = make(map[string]openapi.DiagnosticsprofilesPutRequestDiagnosticsProfileValue)
		}
	case "diagnostics_port_profile":
		if operationType == "PUT" {
			b.diagnosticsPortProfilePut = make(map[string]openapi.DiagnosticsportprofilesPutRequestDiagnosticsPortProfileValue)
		} else {
			b.diagnosticsPortProfilePatch = make(map[string]openapi.DiagnosticsportprofilesPutRequestDiagnosticsPortProfileValue)
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

				case "device_settings":
					resp, err := b.client.DeviceSettingsAPI.DevicesettingsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						EthDeviceProfiles map[string]interface{} `json:"eth_device_profiles"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.EthDeviceProfiles, nil

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

				case "ipv4_list_filter":
					resp, err := b.client.IPv4ListFiltersAPI.Ipv4listsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						Ipv4ListFilter map[string]interface{} `json:"ipv4_list_filter"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.Ipv4ListFilter, nil

				case "ipv4_prefix_list":
					resp, err := b.client.IPv4PrefixListsAPI.Ipv4prefixlistsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						Ipv4PrefixList map[string]interface{} `json:"ipv4_prefix_list"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.Ipv4PrefixList, nil

				case "ipv6_list_filter":
					resp, err := b.client.IPv6ListFiltersAPI.Ipv6listsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						Ipv6ListFilter map[string]interface{} `json:"ipv6_list_filter"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.Ipv6ListFilter, nil

				case "ipv6_prefix_list":
					resp, err := b.client.IPv6PrefixListsAPI.Ipv6prefixlistsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						Ipv6PrefixList map[string]interface{} `json:"ipv6_prefix_list"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.Ipv6PrefixList, nil

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

				case "as_path_access_list":
					resp, err := b.client.ASPathAccessListsAPI.AspathaccesslistsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						AsPathAccessList map[string]interface{} `json:"as_path_access_list"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.AsPathAccessList, nil

				case "community_list":
					resp, err := b.client.CommunityListsAPI.CommunitylistsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						CommunityList map[string]interface{} `json:"community_list"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.CommunityList, nil

				case "extended_community_list":
					resp, err := b.client.ExtendedCommunityListsAPI.ExtendedcommunitylistsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						ExtendedCommunityList map[string]interface{} `json:"extended_community_list"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.ExtendedCommunityList, nil

				case "route_map_clause":
					resp, err := b.client.RouteMapClausesAPI.RoutemapclausesGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						RouteMapClause map[string]interface{} `json:"route_map_clause"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.RouteMapClause, nil

				case "route_map":
					resp, err := b.client.RouteMapsAPI.RoutemapsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						RouteMap map[string]interface{} `json:"route_map"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.RouteMap, nil

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

				case "pod":
					resp, err := b.client.PodsAPI.PodsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						Pod map[string]interface{} `json:"pod"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.Pod, nil

				case "port_acl":
					resp, err := b.client.PortACLsAPI.PortaclsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						PortAcl map[string]interface{} `json:"port_acl"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.PortAcl, nil

				case "sflow_collector":
					resp, err := b.client.SFlowCollectorsAPI.SflowcollectorsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						SflowCollector map[string]interface{} `json:"sflow_collector"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.SflowCollector, nil

				case "diagnostics_profile":
					resp, err := b.client.DiagnosticsProfilesAPI.DiagnosticsprofilesGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						DiagnosticsProfile map[string]interface{} `json:"diagnostics_profile"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.DiagnosticsProfile, nil

				case "diagnostics_port_profile":
					resp, err := b.client.DiagnosticsPortProfilesAPI.DiagnosticsportprofilesGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						DiagnosticsPortProfile map[string]interface{} `json:"diagnostics_port_profile"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.DiagnosticsPortProfile, nil

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
			gatewayMap := make(map[string]openapi.GatewaysPutRequestGatewayValue)
			for name, props := range filteredData {
				gatewayMap[name] = props.(openapi.GatewaysPutRequestGatewayValue)
			}
			putRequest.SetGateway(gatewayMap)
			return putRequest
		case "lag":
			putRequest := openapi.NewLagsPutRequest()
			lagMap := make(map[string]openapi.LagsPutRequestLagValue)
			for name, props := range filteredData {
				lagMap[name] = props.(openapi.LagsPutRequestLagValue)
			}
			putRequest.SetLag(lagMap)
			return putRequest
		case "tenant":
			putRequest := openapi.NewTenantsPutRequest()
			tenantMap := make(map[string]openapi.TenantsPutRequestTenantValue)
			for name, props := range filteredData {
				tenantMap[name] = props.(openapi.TenantsPutRequestTenantValue)
			}
			putRequest.SetTenant(tenantMap)
			return putRequest
		case "service":
			putRequest := openapi.NewServicesPutRequest()
			serviceMap := make(map[string]openapi.ServicesPutRequestServiceValue)
			for name, props := range filteredData {
				serviceMap[name] = props.(openapi.ServicesPutRequestServiceValue)
			}
			putRequest.SetService(serviceMap)
			return putRequest
		case "gateway_profile":
			putRequest := openapi.NewGatewayprofilesPutRequest()
			profileMap := make(map[string]openapi.GatewayprofilesPutRequestGatewayProfileValue)
			for name, props := range filteredData {
				profileMap[name] = props.(openapi.GatewayprofilesPutRequestGatewayProfileValue)
			}
			putRequest.SetGatewayProfile(profileMap)
			return putRequest
		case "eth_port_profile":
			putRequest := openapi.NewEthportprofilesPutRequest()
			profileMap := make(map[string]openapi.EthportprofilesPutRequestEthPortProfileValue)
			for name, props := range filteredData {
				profileMap[name] = props.(openapi.EthportprofilesPutRequestEthPortProfileValue)
			}
			putRequest.SetEthPortProfile(profileMap)
			return putRequest
		case "eth_port_settings":
			putRequest := openapi.NewEthportsettingsPutRequest()
			settingsMap := make(map[string]openapi.EthportsettingsPutRequestEthPortSettingsValue)
			for name, props := range filteredData {
				settingsMap[name] = props.(openapi.EthportsettingsPutRequestEthPortSettingsValue)
			}
			putRequest.SetEthPortSettings(settingsMap)
			return putRequest
		case "device_settings":
			putRequest := openapi.NewDevicesettingsPutRequest()
			settingsMap := make(map[string]openapi.DevicesettingsPutRequestEthDeviceProfilesValue)
			for name, props := range filteredData {
				settingsMap[name] = props.(openapi.DevicesettingsPutRequestEthDeviceProfilesValue)
			}
			putRequest.SetEthDeviceProfiles(settingsMap)
			return putRequest
		case "bundle":
			putRequest := openapi.NewBundlesPutRequest()
			bundleMap := make(map[string]openapi.BundlesPutRequestEndpointBundleValue)
			for name, props := range filteredData {
				bundleMap[name] = props.(openapi.BundlesPutRequestEndpointBundleValue)
			}
			putRequest.SetEndpointBundle(bundleMap)
			return putRequest
		case "acl":
			putRequest := openapi.NewAclsPutRequest()
			aclMap := make(map[string]openapi.AclsPutRequestIpFilterValue)
			for name, props := range filteredData {
				aclMap[name] = props.(openapi.AclsPutRequestIpFilterValue)
			}
			putRequest.SetIpFilter(aclMap)
			return putRequest
		case "ipv4_list_filter":
			putRequest := openapi.NewIpv4listsPutRequest()
			ipv4Map := make(map[string]openapi.Ipv4listsPutRequestIpv4ListFilterValue)
			for name, props := range filteredData {
				ipv4Map[name] = props.(openapi.Ipv4listsPutRequestIpv4ListFilterValue)
			}
			putRequest.SetIpv4ListFilter(ipv4Map)
			return putRequest
		case "ipv4_prefix_list":
			putRequest := openapi.NewIpv4prefixlistsPutRequest()
			ipv4PrefixMap := make(map[string]openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValue)
			for name, props := range filteredData {
				ipv4PrefixMap[name] = props.(openapi.Ipv4prefixlistsPutRequestIpv4PrefixListValue)
			}
			putRequest.SetIpv4PrefixList(ipv4PrefixMap)
			return putRequest
		case "ipv6_list_filter":
			putRequest := openapi.NewIpv6listsPutRequest()
			ipv6Map := make(map[string]openapi.Ipv6listsPutRequestIpv6ListFilterValue)
			for name, props := range filteredData {
				ipv6Map[name] = props.(openapi.Ipv6listsPutRequestIpv6ListFilterValue)
			}
			putRequest.SetIpv6ListFilter(ipv6Map)
			return putRequest
		case "ipv6_prefix_list":
			putRequest := openapi.NewIpv6prefixlistsPutRequest()
			ipv6PrefixMap := make(map[string]openapi.Ipv6prefixlistsPutRequestIpv6PrefixListValue)
			for name, props := range filteredData {
				ipv6PrefixMap[name] = props.(openapi.Ipv6prefixlistsPutRequestIpv6PrefixListValue)
			}
			putRequest.SetIpv6PrefixList(ipv6PrefixMap)
			return putRequest
		case "authenticated_eth_port":
			putRequest := openapi.NewAuthenticatedethportsPutRequest()
			portMap := make(map[string]openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValue)
			for name, props := range filteredData {
				portMap[name] = props.(openapi.AuthenticatedethportsPutRequestAuthenticatedEthPortValue)
			}
			putRequest.SetAuthenticatedEthPort(portMap)
			return putRequest
		case "badge":
			putRequest := openapi.NewBadgesPutRequest()
			badgeMap := make(map[string]openapi.BadgesPutRequestBadgeValue)
			for name, props := range filteredData {
				badgeMap[name] = props.(openapi.BadgesPutRequestBadgeValue)
			}
			putRequest.SetBadge(badgeMap)
			return putRequest
		case "device_voice_settings":
			putRequest := openapi.NewDevicevoicesettingsPutRequest()
			settingsMap := make(map[string]openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValue)
			for name, props := range filteredData {
				settingsMap[name] = props.(openapi.DevicevoicesettingsPutRequestDeviceVoiceSettingsValue)
			}
			putRequest.SetDeviceVoiceSettings(settingsMap)
			return putRequest
		case "as_path_access_list":
			putRequest := openapi.NewAspathaccesslistsPutRequest()
			asPathMap := make(map[string]openapi.AspathaccesslistsPutRequestAsPathAccessListValue)
			for name, props := range filteredData {
				asPathMap[name] = props.(openapi.AspathaccesslistsPutRequestAsPathAccessListValue)
			}
			putRequest.SetAsPathAccessList(asPathMap)
			return putRequest
		case "community_list":
			putRequest := openapi.NewCommunitylistsPutRequest()
			communityMap := make(map[string]openapi.CommunitylistsPutRequestCommunityListValue)
			for name, props := range filteredData {
				communityMap[name] = props.(openapi.CommunitylistsPutRequestCommunityListValue)
			}
			putRequest.SetCommunityList(communityMap)
			return putRequest
		case "extended_community_list":
			putRequest := openapi.NewExtendedcommunitylistsPutRequest()
			extCommunityMap := make(map[string]openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValue)
			for name, props := range filteredData {
				extCommunityMap[name] = props.(openapi.ExtendedcommunitylistsPutRequestExtendedCommunityListValue)
			}
			putRequest.SetExtendedCommunityList(extCommunityMap)
			return putRequest
		case "route_map_clause":
			putRequest := openapi.NewRoutemapclausesPutRequest()
			routeMapClauseMap := make(map[string]openapi.RoutemapclausesPutRequestRouteMapClauseValue)
			for name, props := range filteredData {
				routeMapClauseMap[name] = props.(openapi.RoutemapclausesPutRequestRouteMapClauseValue)
			}
			putRequest.SetRouteMapClause(routeMapClauseMap)
			return putRequest
		case "route_map":
			putRequest := openapi.NewRoutemapsPutRequest()
			routeMapMap := make(map[string]openapi.RoutemapsPutRequestRouteMapValue)
			for name, props := range filteredData {
				routeMapMap[name] = props.(openapi.RoutemapsPutRequestRouteMapValue)
			}
			putRequest.SetRouteMap(routeMapMap)
			return putRequest
		case "sfp_breakout":
			// SFP Breakouts only support PATCH operations
			patchRequest := openapi.NewSfpbreakoutsPatchRequest()
			sfpMap := make(map[string]openapi.SfpbreakoutsPatchRequestSfpBreakoutsValue)
			for name, props := range filteredData {
				sfpMap[name] = props.(openapi.SfpbreakoutsPatchRequestSfpBreakoutsValue)
			}
			patchRequest.SetSfpBreakouts(sfpMap)
			return patchRequest
		case "site":
			// Sites only support PATCH operations
			patchRequest := openapi.NewSitesPatchRequest()
			siteMap := make(map[string]openapi.SitesPatchRequestSiteValue)
			for name, props := range filteredData {
				siteMap[name] = props.(openapi.SitesPatchRequestSiteValue)
			}
			patchRequest.SetSite(siteMap)
			return patchRequest
		case "packet_broker":
			putRequest := openapi.NewPacketbrokerPutRequest()
			brokerMap := make(map[string]openapi.PacketbrokerPutRequestPbEgressProfileValue)
			for name, props := range filteredData {
				brokerMap[name] = props.(openapi.PacketbrokerPutRequestPbEgressProfileValue)
			}
			putRequest.SetPbEgressProfile(brokerMap)
			return putRequest
		case "packet_queue":
			putRequest := openapi.NewPacketqueuesPutRequest()
			queueMap := make(map[string]openapi.PacketqueuesPutRequestPacketQueueValue)
			for name, props := range filteredData {
				queueMap[name] = props.(openapi.PacketqueuesPutRequestPacketQueueValue)
			}
			putRequest.SetPacketQueue(queueMap)
			return putRequest
		case "service_port_profile":
			putRequest := openapi.NewServiceportprofilesPutRequest()
			profileMap := make(map[string]openapi.ServiceportprofilesPutRequestServicePortProfileValue)
			for name, props := range filteredData {
				profileMap[name] = props.(openapi.ServiceportprofilesPutRequestServicePortProfileValue)
			}
			putRequest.SetServicePortProfile(profileMap)
			return putRequest
		case "switchpoint":
			putRequest := openapi.NewSwitchpointsPutRequest()
			switchpointMap := make(map[string]openapi.SwitchpointsPutRequestSwitchpointValue)
			for name, props := range filteredData {
				switchpointMap[name] = props.(openapi.SwitchpointsPutRequestSwitchpointValue)
			}
			putRequest.SetSwitchpoint(switchpointMap)
			return putRequest
		case "voice_port_profile":
			putRequest := openapi.NewVoiceportprofilesPutRequest()
			profileMap := make(map[string]openapi.VoiceportprofilesPutRequestVoicePortProfilesValue)
			for name, props := range filteredData {
				profileMap[name] = props.(openapi.VoiceportprofilesPutRequestVoicePortProfilesValue)
			}
			putRequest.SetVoicePortProfiles(profileMap)
			return putRequest
		case "device_controller":
			putRequest := openapi.NewDevicecontrollersPutRequest()
			deviceMap := make(map[string]openapi.DevicecontrollersPutRequestDeviceControllerValue)
			for name, props := range filteredData {
				deviceMap[name] = props.(openapi.DevicecontrollersPutRequestDeviceControllerValue)
			}
			putRequest.SetDeviceController(deviceMap)
			return putRequest
		case "pod":
			putRequest := openapi.NewPodsPutRequest()
			podMap := make(map[string]openapi.PodsPutRequestPodValue)
			for name, props := range filteredData {
				podMap[name] = props.(openapi.PodsPutRequestPodValue)
			}
			putRequest.SetPod(podMap)
			return putRequest
		case "port_acl":
			putRequest := openapi.NewPortaclsPutRequest()
			portAclMap := make(map[string]openapi.PortaclsPutRequestPortAclValue)
			for name, props := range filteredData {
				portAclMap[name] = props.(openapi.PortaclsPutRequestPortAclValue)
			}
			putRequest.SetPortAcl(portAclMap)
			return putRequest
		case "sflow_collector":
			putRequest := openapi.NewSflowcollectorsPutRequest()
			sflowMap := make(map[string]openapi.SflowcollectorsPutRequestSflowCollectorValue)
			for name, props := range filteredData {
				sflowMap[name] = props.(openapi.SflowcollectorsPutRequestSflowCollectorValue)
			}
			putRequest.SetSflowCollector(sflowMap)
			return putRequest
		case "diagnostics_profile":
			putRequest := openapi.NewDiagnosticsprofilesPutRequest()
			diagnosticsMap := make(map[string]openapi.DiagnosticsprofilesPutRequestDiagnosticsProfileValue)
			for name, props := range filteredData {
				diagnosticsMap[name] = props.(openapi.DiagnosticsprofilesPutRequestDiagnosticsProfileValue)
			}
			putRequest.SetDiagnosticsProfile(diagnosticsMap)
			return putRequest
		case "diagnostics_port_profile":
			putRequest := openapi.NewDiagnosticsportprofilesPutRequest()
			diagnosticsPortMap := make(map[string]openapi.DiagnosticsportprofilesPutRequestDiagnosticsPortProfileValue)
			for name, props := range filteredData {
				diagnosticsPortMap[name] = props.(openapi.DiagnosticsportprofilesPutRequestDiagnosticsPortProfileValue)
			}
			putRequest.SetDiagnosticsPortProfile(diagnosticsPortMap)
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

		case "site":
			sitesReq := b.client.SitesAPI.SitesGet(fetchCtx)
			sitesResp, fetchErr := sitesReq.Execute()

			if fetchErr != nil {
				tflog.Error(ctx, "Failed to fetch sites after PATCH for auto-generated fields", map[string]interface{}{
					"error": fetchErr.Error(),
				})
				return fetchErr
			}

			defer sitesResp.Body.Close()

			var sitesData struct {
				Site map[string]map[string]interface{} `json:"site"`
			}

			if respErr := json.NewDecoder(sitesResp.Body).Decode(&sitesData); respErr != nil {
				tflog.Error(ctx, "Failed to decode sites response for auto-generated fields", map[string]interface{}{
					"error": respErr.Error(),
				})
				return respErr
			}

			b.siteResponsesMutex.Lock()
			for siteName, siteData := range sitesData.Site {
				b.siteResponses[siteName] = siteData

				if name, ok := siteData["name"].(string); ok && name != siteName {
					b.siteResponses[name] = siteData
				}
			}
			b.siteResponsesMutex.Unlock()

			tflog.Debug(ctx, "Successfully stored site data for auto-generated fields", map[string]interface{}{
				"site_count": len(sitesData.Site),
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
		case "device_settings":
			b.recentDeviceSettingsOps = true
			b.recentDeviceSettingsOpTime = now
		case "bundle":
			b.recentBundleOps = true
			b.recentBundleOpTime = now
		case "acl":
			b.recentAclOps = true
			b.recentAclOpTime = now
		case "ipv4_list_filter":
			b.recentIpv4ListOps = true
			b.recentIpv4ListOpTime = now
		case "ipv4_prefix_list":
			b.recentIpv4PrefixListOps = true
			b.recentIpv4PrefixListOpTime = now
		case "ipv6_list_filter":
			b.recentIpv6ListOps = true
			b.recentIpv6ListOpTime = now
		case "ipv6_prefix_list":
			b.recentIpv6PrefixListOps = true
			b.recentIpv6PrefixListOpTime = now
		case "authenticated_eth_port":
			b.recentAuthenticatedEthPortOps = true
			b.recentAuthenticatedEthPortOpTime = now
		case "badge":
			b.recentBadgeOps = true
			b.recentBadgeOpTime = now
		case "device_voice_settings":
			b.recentDeviceVoiceSettingsOps = true
			b.recentDeviceVoiceSettingsOpTime = now
		case "as_path_access_list":
			b.recentAsPathAccessListOps = true
			b.recentAsPathAccessListOpTime = now
		case "community_list":
			b.recentCommunityListOps = true
			b.recentCommunityListOpTime = now
		case "extended_community_list":
			b.recentExtendedCommunityListOps = true
			b.recentExtendedCommunityListOpTime = now
		case "route_map_clause":
			b.recentRouteMapClauseOps = true
			b.recentRouteMapClauseOpTime = now
		case "route_map":
			b.recentRouteMapOps = true
			b.recentRouteMapOpTime = now
		case "sfp_breakout":
			b.recentSfpBreakoutOps = true
			b.recentSfpBreakoutOpTime = now
		case "site":
			b.recentSiteOps = true
			b.recentSiteOpTime = now
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
		case "pod":
			b.recentPodOps = true
			b.recentPodOpTime = now
		case "port_acl":
			b.recentPortAclOps = true
			b.recentPortAclOpTime = now
		case "sflow_collector":
			b.recentSflowCollectorOps = true
			b.recentSflowCollectorOpTime = now
		case "diagnostics_profile":
			b.recentDiagnosticsProfileOps = true
			b.recentDiagnosticsProfileOpTime = now
		case "diagnostics_port_profile":
			b.recentDiagnosticsPortProfileOps = true
			b.recentDiagnosticsPortProfileOpTime = now
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
