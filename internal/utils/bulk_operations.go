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

// ResourceOperations holds all operation data for a single resource type.
type ResourceOperations struct {
	Put            map[string]interface{}            // Pending PUT operations
	Patch          map[string]interface{}            // Pending PATCH operations
	Delete         []string                          // Pending DELETE operations
	RecentOps      bool                              // Whether recent operations occurred
	RecentOpTime   time.Time                         // Time of most recent operation
	Responses      map[string]map[string]interface{} // Cached API responses
	ResponsesMutex sync.RWMutex                      // Mutex for response cache
}

// NewResourceOperations creates a new ResourceOperations instance with initialized maps.
func NewResourceOperations() *ResourceOperations {
	return &ResourceOperations{
		Put:       make(map[string]interface{}),
		Patch:     make(map[string]interface{}),
		Delete:    make([]string, 0),
		RecentOps: false,
		Responses: make(map[string]map[string]interface{}),
	}
}

// ResourceConfig defines configuration for a specific resource type.
type ResourceConfig struct {
	ResourceType     string                                                                         // String identifier for the resource
	PutRequestType   reflect.Type                                                                   // Type for PUT requests
	PatchRequestType reflect.Type                                                                   // Type for PATCH requests
	HasAutoGen       bool                                                                           // Whether resource has auto-generated fields
	APIClientGetter  func(*openapi.APIClient) ResourceAPIClient                                     // Function to get API client
	PutFunc          func(*openapi.APIClient, context.Context, interface{}) (*http.Response, error) // Direct PUT API call
	PatchFunc        func(*openapi.APIClient, context.Context, interface{}) (*http.Response, error) // Direct PATCH API call
	DeleteFunc       func(*openapi.APIClient, context.Context, []string) (*http.Response, error)    // Direct DELETE API call
	GetFunc          func(*openapi.APIClient, context.Context) (*http.Response, error)              // Direct GET API call
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
	resources         map[string]*ResourceOperations

	aclIpVersion map[string]string // Track which IP version each ACL operation uses

	// For tracking operations
	pendingOperations     map[string]*Operation
	operationResults      map[string]bool // true = success, false = failure
	operationErrors       map[string]error
	operationWaitChannels map[string]chan struct{}
	operationMutex        sync.Mutex
	closedChannels        map[string]bool
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
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.GatewaysAPI.GatewaysPut(ctx).GatewaysPutRequest(*req.(*openapi.GatewaysPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.GatewaysAPI.GatewaysPatch(ctx).GatewaysPutRequest(*req.(*openapi.GatewaysPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.GatewaysAPI.GatewaysDelete(ctx).GatewayName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.GatewaysAPI.GatewaysGet(ctx).Execute()
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
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.LAGsAPI.LagsPut(ctx).LagsPutRequest(*req.(*openapi.LagsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.LAGsAPI.LagsPatch(ctx).LagsPutRequest(*req.(*openapi.LagsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.LAGsAPI.LagsDelete(ctx).LagName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.LAGsAPI.LagsGet(ctx).Execute()
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
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.TenantsAPI.TenantsPut(ctx).TenantsPutRequest(*req.(*openapi.TenantsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.TenantsAPI.TenantsPatch(ctx).TenantsPutRequest(*req.(*openapi.TenantsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.TenantsAPI.TenantsDelete(ctx).TenantName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.TenantsAPI.TenantsGet(ctx).Execute()
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
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.ServicesAPI.ServicesPut(ctx).ServicesPutRequest(*req.(*openapi.ServicesPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.ServicesAPI.ServicesPatch(ctx).ServicesPutRequest(*req.(*openapi.ServicesPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.ServicesAPI.ServicesDelete(ctx).ServiceName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.ServicesAPI.ServicesGet(ctx).Execute()
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
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.GatewayProfilesAPI.GatewayprofilesPut(ctx).GatewayprofilesPutRequest(*req.(*openapi.GatewayprofilesPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.GatewayProfilesAPI.GatewayprofilesPatch(ctx).GatewayprofilesPutRequest(*req.(*openapi.GatewayprofilesPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.GatewayProfilesAPI.GatewayprofilesDelete(ctx).ProfileName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.GatewayProfilesAPI.GatewayprofilesGet(ctx).Execute()
		},
	},
	"grouping_rule": {
		ResourceType:     "grouping_rule",
		PutRequestType:   reflect.TypeOf(openapi.GroupingrulesPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.GroupingrulesPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "grouping_rule"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.GroupingRulesAPI.GroupingrulesPut(ctx).GroupingrulesPutRequest(*req.(*openapi.GroupingrulesPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.GroupingRulesAPI.GroupingrulesPatch(ctx).GroupingrulesPutRequest(*req.(*openapi.GroupingrulesPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.GroupingRulesAPI.GroupingrulesDelete(ctx).GroupingRulesName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.GroupingRulesAPI.GroupingrulesGet(ctx).Execute()
		},
	},
	"threshold_group": {
		ResourceType:     "threshold_group",
		PutRequestType:   reflect.TypeOf(openapi.ThresholdgroupsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.ThresholdgroupsPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "threshold_group"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.ThresholdGroupsAPI.ThresholdgroupsPut(ctx).ThresholdgroupsPutRequest(*req.(*openapi.ThresholdgroupsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.ThresholdGroupsAPI.ThresholdgroupsPatch(ctx).ThresholdgroupsPutRequest(*req.(*openapi.ThresholdgroupsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.ThresholdGroupsAPI.ThresholdgroupsDelete(ctx).ThresholdGroupName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.ThresholdGroupsAPI.ThresholdgroupsGet(ctx).Execute()
		},
	},
	"threshold": {
		ResourceType:     "threshold",
		PutRequestType:   reflect.TypeOf(openapi.ThresholdsPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.ThresholdsPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "threshold"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.ThresholdsAPI.ThresholdsPut(ctx).ThresholdsPutRequest(*req.(*openapi.ThresholdsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.ThresholdsAPI.ThresholdsPatch(ctx).ThresholdsPutRequest(*req.(*openapi.ThresholdsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.ThresholdsAPI.ThresholdsDelete(ctx).ThresholdName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.ThresholdsAPI.ThresholdsGet(ctx).Execute()
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
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.PacketQueuesAPI.PacketqueuesPut(ctx).PacketqueuesPutRequest(*req.(*openapi.PacketqueuesPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.PacketQueuesAPI.PacketqueuesPatch(ctx).PacketqueuesPutRequest(*req.(*openapi.PacketqueuesPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.PacketQueuesAPI.PacketqueuesDelete(ctx).PacketQueueName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.PacketQueuesAPI.PacketqueuesGet(ctx).Execute()
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
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.EthPortProfilesAPI.EthportprofilesPut(ctx).EthportprofilesPutRequest(*req.(*openapi.EthportprofilesPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.EthPortProfilesAPI.EthportprofilesPatch(ctx).EthportprofilesPutRequest(*req.(*openapi.EthportprofilesPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.EthPortProfilesAPI.EthportprofilesDelete(ctx).ProfileName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.EthPortProfilesAPI.EthportprofilesGet(ctx).Execute()
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
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.EthPortSettingsAPI.EthportsettingsPut(ctx).EthportsettingsPutRequest(*req.(*openapi.EthportsettingsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.EthPortSettingsAPI.EthportsettingsPatch(ctx).EthportsettingsPutRequest(*req.(*openapi.EthportsettingsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.EthPortSettingsAPI.EthportsettingsDelete(ctx).PortName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.EthPortSettingsAPI.EthportsettingsGet(ctx).Execute()
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
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.BundlesAPI.BundlesPut(ctx).BundlesPutRequest(*req.(*openapi.BundlesPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.BundlesAPI.BundlesPatch(ctx).BundlesPutRequest(*req.(*openapi.BundlesPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.BundlesAPI.BundlesDelete(ctx).BundleName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.BundlesAPI.BundlesGet(ctx).Execute()
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
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.ACLsAPI.AclsPut(ctx).AclsPutRequest(*req.(*openapi.AclsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.ACLsAPI.AclsPatch(ctx).AclsPutRequest(*req.(*openapi.AclsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.ACLsAPI.AclsDelete(ctx).IpFilterName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.ACLsAPI.AclsGet(ctx).Execute()
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
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.PacketBrokerAPI.PacketbrokerPut(ctx).PacketbrokerPutRequest(*req.(*openapi.PacketbrokerPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.PacketBrokerAPI.PacketbrokerPatch(ctx).PacketbrokerPutRequest(*req.(*openapi.PacketbrokerPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.PacketBrokerAPI.PacketbrokerDelete(ctx).PbEgressProfileName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.PacketBrokerAPI.PacketbrokerGet(ctx).Execute()
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
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.BadgesAPI.BadgesPut(ctx).BadgesPutRequest(*req.(*openapi.BadgesPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.BadgesAPI.BadgesPatch(ctx).BadgesPutRequest(*req.(*openapi.BadgesPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.BadgesAPI.BadgesDelete(ctx).BadgeName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.BadgesAPI.BadgesGet(ctx).Execute()
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
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.SwitchpointsAPI.SwitchpointsPut(ctx).SwitchpointsPutRequest(*req.(*openapi.SwitchpointsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.SwitchpointsAPI.SwitchpointsPatch(ctx).SwitchpointsPutRequest(*req.(*openapi.SwitchpointsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.SwitchpointsAPI.SwitchpointsDelete(ctx).SwitchpointName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.SwitchpointsAPI.SwitchpointsGet(ctx).Execute()
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
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.DeviceControllersAPI.DevicecontrollersPut(ctx).DevicecontrollersPutRequest(*req.(*openapi.DevicecontrollersPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.DeviceControllersAPI.DevicecontrollersPatch(ctx).DevicecontrollersPutRequest(*req.(*openapi.DevicecontrollersPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.DeviceControllersAPI.DevicecontrollersDelete(ctx).DeviceControllerName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.DeviceControllersAPI.DevicecontrollersGet(ctx).Execute()
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
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.AuthenticatedEthPortsAPI.AuthenticatedethportsPut(ctx).AuthenticatedethportsPutRequest(*req.(*openapi.AuthenticatedethportsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.AuthenticatedEthPortsAPI.AuthenticatedethportsPatch(ctx).AuthenticatedethportsPutRequest(*req.(*openapi.AuthenticatedethportsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.AuthenticatedEthPortsAPI.AuthenticatedethportsDelete(ctx).AuthenticatedEthPortName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.AuthenticatedEthPortsAPI.AuthenticatedethportsGet(ctx).Execute()
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
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.DeviceVoiceSettingsAPI.DevicevoicesettingsPut(ctx).DevicevoicesettingsPutRequest(*req.(*openapi.DevicevoicesettingsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.DeviceVoiceSettingsAPI.DevicevoicesettingsPatch(ctx).DevicevoicesettingsPutRequest(*req.(*openapi.DevicevoicesettingsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.DeviceVoiceSettingsAPI.DevicevoicesettingsDelete(ctx).DeviceVoiceSettingsName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.DeviceVoiceSettingsAPI.DevicevoicesettingsGet(ctx).Execute()
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
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.VoicePortProfilesAPI.VoiceportprofilesPut(ctx).VoiceportprofilesPutRequest(*req.(*openapi.VoiceportprofilesPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.VoicePortProfilesAPI.VoiceportprofilesPatch(ctx).VoiceportprofilesPutRequest(*req.(*openapi.VoiceportprofilesPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.VoicePortProfilesAPI.VoiceportprofilesDelete(ctx).VoicePortProfileName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.VoicePortProfilesAPI.VoiceportprofilesGet(ctx).Execute()
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
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.ServicePortProfilesAPI.ServiceportprofilesPut(ctx).ServiceportprofilesPutRequest(*req.(*openapi.ServiceportprofilesPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.ServicePortProfilesAPI.ServiceportprofilesPatch(ctx).ServiceportprofilesPutRequest(*req.(*openapi.ServiceportprofilesPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.ServicePortProfilesAPI.ServiceportprofilesDelete(ctx).ServicePortProfileName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.ServicePortProfilesAPI.ServiceportprofilesGet(ctx).Execute()
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
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.ASPathAccessListsAPI.AspathaccesslistsPut(ctx).AspathaccesslistsPutRequest(*req.(*openapi.AspathaccesslistsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.ASPathAccessListsAPI.AspathaccesslistsPatch(ctx).AspathaccesslistsPutRequest(*req.(*openapi.AspathaccesslistsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.ASPathAccessListsAPI.AspathaccesslistsDelete(ctx).AsPathAccessListName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.ASPathAccessListsAPI.AspathaccesslistsGet(ctx).Execute()
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
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.CommunityListsAPI.CommunitylistsPut(ctx).CommunitylistsPutRequest(*req.(*openapi.CommunitylistsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.CommunityListsAPI.CommunitylistsPatch(ctx).CommunitylistsPutRequest(*req.(*openapi.CommunitylistsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.CommunityListsAPI.CommunitylistsDelete(ctx).CommunityListName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.CommunityListsAPI.CommunitylistsGet(ctx).Execute()
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
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.DeviceSettingsAPI.DevicesettingsPut(ctx).DevicesettingsPutRequest(*req.(*openapi.DevicesettingsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.DeviceSettingsAPI.DevicesettingsPatch(ctx).DevicesettingsPutRequest(*req.(*openapi.DevicesettingsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.DeviceSettingsAPI.DevicesettingsDelete(ctx).EthDeviceProfilesName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.DeviceSettingsAPI.DevicesettingsGet(ctx).Execute()
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
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.ExtendedCommunityListsAPI.ExtendedcommunitylistsPut(ctx).ExtendedcommunitylistsPutRequest(*req.(*openapi.ExtendedcommunitylistsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.ExtendedCommunityListsAPI.ExtendedcommunitylistsPatch(ctx).ExtendedcommunitylistsPutRequest(*req.(*openapi.ExtendedcommunitylistsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.ExtendedCommunityListsAPI.ExtendedcommunitylistsDelete(ctx).ExtendedCommunityListName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.ExtendedCommunityListsAPI.ExtendedcommunitylistsGet(ctx).Execute()
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
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.IPv4ListFiltersAPI.Ipv4listsPut(ctx).Ipv4listsPutRequest(*req.(*openapi.Ipv4listsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.IPv4ListFiltersAPI.Ipv4listsPatch(ctx).Ipv4listsPutRequest(*req.(*openapi.Ipv4listsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.IPv4ListFiltersAPI.Ipv4listsDelete(ctx).Ipv4ListFilterName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.IPv4ListFiltersAPI.Ipv4listsGet(ctx).Execute()
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
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.IPv4PrefixListsAPI.Ipv4prefixlistsPut(ctx).Ipv4prefixlistsPutRequest(*req.(*openapi.Ipv4prefixlistsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.IPv4PrefixListsAPI.Ipv4prefixlistsPatch(ctx).Ipv4prefixlistsPutRequest(*req.(*openapi.Ipv4prefixlistsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.IPv4PrefixListsAPI.Ipv4prefixlistsDelete(ctx).Ipv4PrefixListName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.IPv4PrefixListsAPI.Ipv4prefixlistsGet(ctx).Execute()
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
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.IPv6ListFiltersAPI.Ipv6listsPut(ctx).Ipv6listsPutRequest(*req.(*openapi.Ipv6listsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.IPv6ListFiltersAPI.Ipv6listsPatch(ctx).Ipv6listsPutRequest(*req.(*openapi.Ipv6listsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.IPv6ListFiltersAPI.Ipv6listsDelete(ctx).Ipv6ListFilterName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.IPv6ListFiltersAPI.Ipv6listsGet(ctx).Execute()
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
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.IPv6PrefixListsAPI.Ipv6prefixlistsPut(ctx).Ipv6prefixlistsPutRequest(*req.(*openapi.Ipv6prefixlistsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.IPv6PrefixListsAPI.Ipv6prefixlistsPatch(ctx).Ipv6prefixlistsPutRequest(*req.(*openapi.Ipv6prefixlistsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.IPv6PrefixListsAPI.Ipv6prefixlistsDelete(ctx).Ipv6PrefixListName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.IPv6PrefixListsAPI.Ipv6prefixlistsGet(ctx).Execute()
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
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.RouteMapClausesAPI.RoutemapclausesPut(ctx).RoutemapclausesPutRequest(*req.(*openapi.RoutemapclausesPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.RouteMapClausesAPI.RoutemapclausesPatch(ctx).RoutemapclausesPutRequest(*req.(*openapi.RoutemapclausesPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.RouteMapClausesAPI.RoutemapclausesDelete(ctx).RouteMapClauseName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.RouteMapClausesAPI.RoutemapclausesGet(ctx).Execute()
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
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.RouteMapsAPI.RoutemapsPut(ctx).RoutemapsPutRequest(*req.(*openapi.RoutemapsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.RouteMapsAPI.RoutemapsPatch(ctx).RoutemapsPutRequest(*req.(*openapi.RoutemapsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.RouteMapsAPI.RoutemapsDelete(ctx).RouteMapName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.RouteMapsAPI.RoutemapsGet(ctx).Execute()
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
		PutFunc: nil, // SFP Breakouts only support PATCH
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.SFPBreakoutsAPI.SfpbreakoutsPatch(ctx).SfpbreakoutsPatchRequest(*req.(*openapi.SfpbreakoutsPatchRequest)).Execute()
		},
		DeleteFunc: nil, // No DELETE operation
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.SFPBreakoutsAPI.SfpbreakoutsGet(ctx).Execute()
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
		PutFunc: nil, // Sites only support PATCH
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.SitesAPI.SitesPatch(ctx).SitesPatchRequest(*req.(*openapi.SitesPatchRequest)).Execute()
		},
		DeleteFunc: nil, // No DELETE operation
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.SitesAPI.SitesGet(ctx).Execute()
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
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.PodsAPI.PodsPut(ctx).PodsPutRequest(*req.(*openapi.PodsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.PodsAPI.PodsPatch(ctx).PodsPutRequest(*req.(*openapi.PodsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.PodsAPI.PodsDelete(ctx).PodName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.PodsAPI.PodsGet(ctx).Execute()
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
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.PortACLsAPI.PortaclsPut(ctx).PortaclsPutRequest(*req.(*openapi.PortaclsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.PortACLsAPI.PortaclsPatch(ctx).PortaclsPutRequest(*req.(*openapi.PortaclsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.PortACLsAPI.PortaclsDelete(ctx).PortAclName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.PortACLsAPI.PortaclsGet(ctx).Execute()
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
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.SFlowCollectorsAPI.SflowcollectorsPut(ctx).SflowcollectorsPutRequest(*req.(*openapi.SflowcollectorsPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.SFlowCollectorsAPI.SflowcollectorsPatch(ctx).SflowcollectorsPutRequest(*req.(*openapi.SflowcollectorsPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.SFlowCollectorsAPI.SflowcollectorsDelete(ctx).SflowCollectorName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.SFlowCollectorsAPI.SflowcollectorsGet(ctx).Execute()
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
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.DiagnosticsProfilesAPI.DiagnosticsprofilesPut(ctx).DiagnosticsprofilesPutRequest(*req.(*openapi.DiagnosticsprofilesPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.DiagnosticsProfilesAPI.DiagnosticsprofilesPatch(ctx).DiagnosticsprofilesPutRequest(*req.(*openapi.DiagnosticsprofilesPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.DiagnosticsProfilesAPI.DiagnosticsprofilesDelete(ctx).DiagnosticsProfileName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.DiagnosticsProfilesAPI.DiagnosticsprofilesGet(ctx).Execute()
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
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.DiagnosticsPortProfilesAPI.DiagnosticsportprofilesPut(ctx).DiagnosticsportprofilesPutRequest(*req.(*openapi.DiagnosticsportprofilesPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.DiagnosticsPortProfilesAPI.DiagnosticsportprofilesPatch(ctx).DiagnosticsportprofilesPutRequest(*req.(*openapi.DiagnosticsportprofilesPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.DiagnosticsPortProfilesAPI.DiagnosticsportprofilesDelete(ctx).DiagnosticsPortProfileName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.DiagnosticsPortProfilesAPI.DiagnosticsportprofilesGet(ctx).Execute()
		},
	},
	"pb_routing": {
		ResourceType:     "pb_routing",
		PutRequestType:   reflect.TypeOf(openapi.PolicybasedroutingPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.PolicybasedroutingPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "pb_routing"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.PBRoutingAPI.PolicybasedroutingPut(ctx).PolicybasedroutingPutRequest(*req.(*openapi.PolicybasedroutingPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.PBRoutingAPI.PolicybasedroutingPatch(ctx).PolicybasedroutingPutRequest(*req.(*openapi.PolicybasedroutingPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.PBRoutingAPI.PolicybasedroutingDelete(ctx).PbRoutingName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.PBRoutingAPI.PolicybasedroutingGet(ctx).Execute()
		},
	},
	"pb_routing_acl": {
		ResourceType:     "pb_routing_acl",
		PutRequestType:   reflect.TypeOf(openapi.PolicybasedroutingaclPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.PolicybasedroutingaclPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "pb_routing_acl"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.PBRoutingACLAPI.PolicybasedroutingaclPut(ctx).PolicybasedroutingaclPutRequest(*req.(*openapi.PolicybasedroutingaclPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.PBRoutingACLAPI.PolicybasedroutingaclPatch(ctx).PolicybasedroutingaclPutRequest(*req.(*openapi.PolicybasedroutingaclPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.PBRoutingACLAPI.PolicybasedroutingaclDelete(ctx).PbRoutingAclName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.PBRoutingACLAPI.PolicybasedroutingaclGet(ctx).Execute()
		},
	},
	"spine_plane": {
		ResourceType:     "spine_plane",
		PutRequestType:   reflect.TypeOf(openapi.SpineplanesPutRequest{}),
		PatchRequestType: reflect.TypeOf(openapi.SpineplanesPutRequest{}),
		HasAutoGen:       false,
		APIClientGetter: func(c *openapi.APIClient) ResourceAPIClient {
			return &GenericAPIClient{client: c, resourceType: "spine_plane"}
		},
		PutFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.SpinePlanesAPI.SpineplanesPut(ctx).SpineplanesPutRequest(*req.(*openapi.SpineplanesPutRequest)).Execute()
		},
		PatchFunc: func(c *openapi.APIClient, ctx context.Context, req interface{}) (*http.Response, error) {
			return c.SpinePlanesAPI.SpineplanesPatch(ctx).SpineplanesPutRequest(*req.(*openapi.SpineplanesPutRequest)).Execute()
		},
		DeleteFunc: func(c *openapi.APIClient, ctx context.Context, names []string) (*http.Response, error) {
			return c.SpinePlanesAPI.SpineplanesDelete(ctx).SpinePlaneName(names).Execute()
		},
		GetFunc: func(c *openapi.APIClient, ctx context.Context) (*http.Response, error) {
			return c.SpinePlanesAPI.SpineplanesGet(ctx).Execute()
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

func initializeResourceOperations() map[string]*ResourceOperations {
	resources := make(map[string]*ResourceOperations)

	// Initialize ResourceOperations for each resource type in the registry
	for resourceType := range resourceRegistry {
		resources[resourceType] = NewResourceOperations()
	}

	return resources
}

func NewBulkOperationManager(client *openapi.APIClient, contextProvider ContextProviderFunc, clearCacheFunc ClearCacheFunc, mode string) *BulkOperationManager {
	return &BulkOperationManager{
		client:                client,
		contextProvider:       contextProvider,
		clearCacheFunc:        clearCacheFunc,
		mode:                  mode,
		lastOperationTime:     time.Now(),
		resources:             initializeResourceOperations(),
		aclIpVersion:          make(map[string]string),
		pendingOperations:     make(map[string]*Operation),
		operationResults:      make(map[string]bool),
		operationErrors:       make(map[string]error),
		operationWaitChannels: make(map[string]chan struct{}),
		closedChannels:        make(map[string]bool),
	}
}

func (b *BulkOperationManager) GetResourceResponse(resourceType, resourceName string) (map[string]interface{}, bool) {
	// Handle special case for ACL versioning - map to base "acl" type
	if resourceType == "acl_v4" || resourceType == "acl_v6" {
		resourceType = "acl"
	}

	// Get the resource operations from the unified map
	res, exists := b.resources[resourceType]
	if !exists {
		return nil, false
	}

	res.ResponsesMutex.RLock()
	defer res.ResponsesMutex.RUnlock()
	response, exists := res.Responses[resourceName]
	return response, exists
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
			b.clearCacheFunc(ctx, b.contextProvider(), "gateway_profiles")
			b.clearCacheFunc(ctx, b.contextProvider(), "services")
			b.clearCacheFunc(ctx, b.contextProvider(), "packet_queues")
			b.clearCacheFunc(ctx, b.contextProvider(), "eth_port_profiles")
			b.clearCacheFunc(ctx, b.contextProvider(), "eth_port_settings")
			b.clearCacheFunc(ctx, b.contextProvider(), "lags")
			b.clearCacheFunc(ctx, b.contextProvider(), "sflow_collectors")
			b.clearCacheFunc(ctx, b.contextProvider(), "diagnostics_profiles")
			b.clearCacheFunc(ctx, b.contextProvider(), "diagnostics_port_profiles")
			b.clearCacheFunc(ctx, b.contextProvider(), "bundles")
			b.clearCacheFunc(ctx, b.contextProvider(), "acls_ipv4")
			b.clearCacheFunc(ctx, b.contextProvider(), "acls_ipv6")
			b.clearCacheFunc(ctx, b.contextProvider(), "packet_brokers")
			b.clearCacheFunc(ctx, b.contextProvider(), "badges")
			b.clearCacheFunc(ctx, b.contextProvider(), "switchpoints")
			b.clearCacheFunc(ctx, b.contextProvider(), "device_controllers")
			b.clearCacheFunc(ctx, b.contextProvider(), "authenticated_eth_ports")
			b.clearCacheFunc(ctx, b.contextProvider(), "device_voice_settings")
			b.clearCacheFunc(ctx, b.contextProvider(), "service_port_profiles")
			b.clearCacheFunc(ctx, b.contextProvider(), "voice_port_profiles")
			b.clearCacheFunc(ctx, b.contextProvider(), "as_path_access_lists")
			b.clearCacheFunc(ctx, b.contextProvider(), "community_lists")
			b.clearCacheFunc(ctx, b.contextProvider(), "device_settings")
			b.clearCacheFunc(ctx, b.contextProvider(), "extended_community_lists")
			b.clearCacheFunc(ctx, b.contextProvider(), "ipv4_lists")
			b.clearCacheFunc(ctx, b.contextProvider(), "ipv4_prefix_lists")
			b.clearCacheFunc(ctx, b.contextProvider(), "ipv6_lists")
			b.clearCacheFunc(ctx, b.contextProvider(), "ipv6_prefix_lists")
			b.clearCacheFunc(ctx, b.contextProvider(), "route_map_clauses")
			b.clearCacheFunc(ctx, b.contextProvider(), "route_maps")
			b.clearCacheFunc(ctx, b.contextProvider(), "sfp_breakouts")
			b.clearCacheFunc(ctx, b.contextProvider(), "sites")
			b.clearCacheFunc(ctx, b.contextProvider(), "pods")
			b.clearCacheFunc(ctx, b.contextProvider(), "port_acls")
			b.clearCacheFunc(ctx, b.contextProvider(), "pb_routing")
			b.clearCacheFunc(ctx, b.contextProvider(), "pb_routing_acl")
			b.clearCacheFunc(ctx, b.contextProvider(), "spine_planes")
			b.clearCacheFunc(ctx, b.contextProvider(), "grouping_rules")
			b.clearCacheFunc(ctx, b.contextProvider(), "threshold_groups")
			b.clearCacheFunc(ctx, b.contextProvider(), "thresholds")
		}
	}

	return diagnostics
}

func (b *BulkOperationManager) getOperationCount(resourceType, operationType string) int {
	res, exists := b.resources[resourceType]
	if !exists {
		return 0
	}

	switch operationType {
	case "PUT":
		if res.Put == nil {
			return 0
		}
		return len(res.Put)
	case "PATCH":
		if res.Patch == nil {
			return 0
		}
		return len(res.Patch)
	case "DELETE":
		return len(res.Delete)
	}
	return 0
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

	// PUT operations - DC Order (note: sfp_breakout and site are skipped - they only support GET and PATCH)
	// 2. ipv6_prefix_list
	if !execute("PUT", b.getOperationCount("ipv6_prefix_list", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv6_prefix_list", "PUT") }, "IPv6 Prefix List") {
		return diagnostics, operationsPerformed
	}
	// 3. community_list
	if !execute("PUT", b.getOperationCount("community_list", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "community_list", "PUT") }, "Community List") {
		return diagnostics, operationsPerformed
	}
	// 4. ipv4_prefix_list
	if !execute("PUT", b.getOperationCount("ipv4_prefix_list", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv4_prefix_list", "PUT") }, "IPv4 Prefix List") {
		return diagnostics, operationsPerformed
	}
	// 5. extended_community_list
	if !execute("PUT", b.getOperationCount("extended_community_list", "PUT"), func(ctx context.Context) diag.Diagnostics {
		return b.ExecuteBulk(ctx, "extended_community_list", "PUT")
	}, "Extended Community List") {
		return diagnostics, operationsPerformed
	}
	// 6. as_path_access_list
	if !execute("PUT", b.getOperationCount("as_path_access_list", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "as_path_access_list", "PUT") }, "AS Path Access List") {
		return diagnostics, operationsPerformed
	}
	// 7. route_map_clause
	if !execute("PUT", b.getOperationCount("route_map_clause", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "route_map_clause", "PUT") }, "Route Map Clause") {
		return diagnostics, operationsPerformed
	}
	// 8-9. acl (both ipv6 and ipv4)
	if !execute("PUT", b.getOperationCount("acl", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "acl", "PUT") }, "ACL") {
		return diagnostics, operationsPerformed
	}
	// 10. route_map
	if !execute("PUT", b.getOperationCount("route_map", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "route_map", "PUT") }, "Route Map") {
		return diagnostics, operationsPerformed
	}
	// 11. pb_routing_acl
	if !execute("PUT", b.getOperationCount("pb_routing_acl", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "pb_routing_acl", "PUT") }, "PB Routing ACL") {
		return diagnostics, operationsPerformed
	}
	// 12. tenant
	if !execute("PUT", b.getOperationCount("tenant", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "tenant", "PUT") }, "Tenant") {
		return diagnostics, operationsPerformed
	}
	// 13. pb_routing
	if !execute("PUT", b.getOperationCount("pb_routing", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "pb_routing", "PUT") }, "PB Routing") {
		return diagnostics, operationsPerformed
	}
	// 14. ipv4_list
	if !execute("PUT", b.getOperationCount("ipv4_list", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv4_list", "PUT") }, "IPv4 List") {
		return diagnostics, operationsPerformed
	}
	// 15. ipv6_list
	if !execute("PUT", b.getOperationCount("ipv6_list", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv6_list", "PUT") }, "IPv6 List") {
		return diagnostics, operationsPerformed
	}
	// 16. service
	if !execute("PUT", b.getOperationCount("service", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "service", "PUT") }, "Service") {
		return diagnostics, operationsPerformed
	}
	// 17. port_acl
	if !execute("PUT", b.getOperationCount("port_acl", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "port_acl", "PUT") }, "Port ACL") {
		return diagnostics, operationsPerformed
	}
	// 18. packet_broker
	if !execute("PUT", b.getOperationCount("packet_broker", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "packet_broker", "PUT") }, "Packet Broker") {
		return diagnostics, operationsPerformed
	}
	// 19. eth_port_profile
	if !execute("PUT", b.getOperationCount("eth_port_profile", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "eth_port_profile", "PUT") }, "Eth Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 20. packet_queue
	if !execute("PUT", b.getOperationCount("packet_queue", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "packet_queue", "PUT") }, "Packet Queue") {
		return diagnostics, operationsPerformed
	}
	// 21. sflow_collector
	if !execute("PUT", b.getOperationCount("sflow_collector", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "sflow_collector", "PUT") }, "SFlow Collector") {
		return diagnostics, operationsPerformed
	}
	// 22. gateway
	if !execute("PUT", b.getOperationCount("gateway", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "gateway", "PUT") }, "Gateway") {
		return diagnostics, operationsPerformed
	}
	// 23. lag
	if !execute("PUT", b.getOperationCount("lag", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "lag", "PUT") }, "LAG") {
		return diagnostics, operationsPerformed
	}
	// 24. eth_port_settings
	if !execute("PUT", b.getOperationCount("eth_port_settings", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "eth_port_settings", "PUT") }, "Eth Port Settings") {
		return diagnostics, operationsPerformed
	}
	// 25. diagnostics_profile
	if !execute("PUT", b.getOperationCount("diagnostics_profile", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "diagnostics_profile", "PUT") }, "Diagnostics Profile") {
		return diagnostics, operationsPerformed
	}
	// 26. gateway_profile
	if !execute("PUT", b.getOperationCount("gateway_profile", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "gateway_profile", "PUT") }, "Gateway Profile") {
		return diagnostics, operationsPerformed
	}
	// 27. diagnostics_port_profile
	if !execute("PUT", b.getOperationCount("diagnostics_port_profile", "PUT"), func(ctx context.Context) diag.Diagnostics {
		return b.ExecuteBulk(ctx, "diagnostics_port_profile", "PUT")
	}, "Diagnostics Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 28. bundle
	if !execute("PUT", b.getOperationCount("bundle", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "bundle", "PUT") }, "Bundle") {
		return diagnostics, operationsPerformed
	}
	// 29. pod
	if !execute("PUT", b.getOperationCount("pod", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "pod", "PUT") }, "Pod") {
		return diagnostics, operationsPerformed
	}
	// 30. badge
	if !execute("PUT", b.getOperationCount("badge", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "badge", "PUT") }, "Badge") {
		return diagnostics, operationsPerformed
	}
	// 31. spine_plane
	if !execute("PUT", b.getOperationCount("spine_plane", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "spine_plane", "PUT") }, "Spine Plane") {
		return diagnostics, operationsPerformed
	}
	// 32. switchpoint
	if !execute("PUT", b.getOperationCount("switchpoint", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "switchpoint", "PUT") }, "Switchpoint") {
		return diagnostics, operationsPerformed
	}
	// 33. device_settings
	if !execute("PUT", b.getOperationCount("device_settings", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "device_settings", "PUT") }, "Device Settings") {
		return diagnostics, operationsPerformed
	}
	// 34. threshold
	if !execute("PUT", b.getOperationCount("threshold", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "threshold", "PUT") }, "Threshold") {
		return diagnostics, operationsPerformed
	}
	// 35. grouping_rule
	if !execute("PUT", b.getOperationCount("grouping_rule", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "grouping_rule", "PUT") }, "Grouping Rule") {
		return diagnostics, operationsPerformed
	}
	// 36. threshold_group
	if !execute("PUT", b.getOperationCount("threshold_group", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "threshold_group", "PUT") }, "Threshold Group") {
		return diagnostics, operationsPerformed
	}
	// 38. device_controller
	if !execute("PUT", b.getOperationCount("device_controller", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "device_controller", "PUT") }, "Device Controller") {
		return diagnostics, operationsPerformed
	}

	// PATCH operations - DC Order
	// 1. sfp_breakout
	if !execute("PATCH", b.getOperationCount("sfp_breakout", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "sfp_breakout", "PATCH") }, "SFP Breakout") {
		return diagnostics, operationsPerformed
	}
	// 2. ipv6_prefix_list
	if !execute("PATCH", b.getOperationCount("ipv6_prefix_list", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv6_prefix_list", "PATCH") }, "IPv6 Prefix List") {
		return diagnostics, operationsPerformed
	}
	// 3. community_list
	if !execute("PATCH", b.getOperationCount("community_list", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "community_list", "PATCH") }, "Community List") {
		return diagnostics, operationsPerformed
	}
	// 4. ipv4_prefix_list
	if !execute("PATCH", b.getOperationCount("ipv4_prefix_list", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv4_prefix_list", "PATCH") }, "IPv4 Prefix List") {
		return diagnostics, operationsPerformed
	}
	// 5. extended_community_list
	if !execute("PATCH", b.getOperationCount("extended_community_list", "PATCH"), func(ctx context.Context) diag.Diagnostics {
		return b.ExecuteBulk(ctx, "extended_community_list", "PATCH")
	}, "Extended Community List") {
		return diagnostics, operationsPerformed
	}
	// 6. as_path_access_list
	if !execute("PATCH", b.getOperationCount("as_path_access_list", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "as_path_access_list", "PATCH") }, "AS Path Access List") {
		return diagnostics, operationsPerformed
	}
	// 7. route_map_clause
	if !execute("PATCH", b.getOperationCount("route_map_clause", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "route_map_clause", "PATCH") }, "Route Map Clause") {
		return diagnostics, operationsPerformed
	}
	// 8-9. acl (both ipv6 and ipv4)
	if !execute("PATCH", b.getOperationCount("acl", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "acl", "PATCH") }, "ACL") {
		return diagnostics, operationsPerformed
	}
	// 10. route_map
	if !execute("PATCH", b.getOperationCount("route_map", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "route_map", "PATCH") }, "Route Map") {
		return diagnostics, operationsPerformed
	}
	// 11. pb_routing_acl
	if !execute("PATCH", b.getOperationCount("pb_routing_acl", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "pb_routing_acl", "PATCH") }, "PB Routing ACL") {
		return diagnostics, operationsPerformed
	}
	// 12. tenant
	if !execute("PATCH", b.getOperationCount("tenant", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "tenant", "PATCH") }, "Tenant") {
		return diagnostics, operationsPerformed
	}
	// 13. pb_routing
	if !execute("PATCH", b.getOperationCount("pb_routing", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "pb_routing", "PATCH") }, "PB Routing") {
		return diagnostics, operationsPerformed
	}
	// 14. ipv4_list
	if !execute("PATCH", b.getOperationCount("ipv4_list", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv4_list", "PATCH") }, "IPv4 List") {
		return diagnostics, operationsPerformed
	}
	// 15. ipv6_list
	if !execute("PATCH", b.getOperationCount("ipv6_list", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv6_list", "PATCH") }, "IPv6 List") {
		return diagnostics, operationsPerformed
	}
	// 16. service
	if !execute("PATCH", b.getOperationCount("service", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "service", "PATCH") }, "Service") {
		return diagnostics, operationsPerformed
	}
	// 17. port_acl
	if !execute("PATCH", b.getOperationCount("port_acl", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "port_acl", "PATCH") }, "Port ACL") {
		return diagnostics, operationsPerformed
	}
	// 18. packet_broker
	if !execute("PATCH", b.getOperationCount("packet_broker", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "packet_broker", "PATCH") }, "Packet Broker") {
		return diagnostics, operationsPerformed
	}
	// 19. eth_port_profile
	if !execute("PATCH", b.getOperationCount("eth_port_profile", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "eth_port_profile", "PATCH") }, "Eth Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 20. packet_queue
	if !execute("PATCH", b.getOperationCount("packet_queue", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "packet_queue", "PATCH") }, "Packet Queue") {
		return diagnostics, operationsPerformed
	}
	// 21. sflow_collector
	if !execute("PATCH", b.getOperationCount("sflow_collector", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "sflow_collector", "PATCH") }, "SFlow Collector") {
		return diagnostics, operationsPerformed
	}
	// 22. gateway
	if !execute("PATCH", b.getOperationCount("gateway", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "gateway", "PATCH") }, "Gateway") {
		return diagnostics, operationsPerformed
	}
	// 23. lag
	if !execute("PATCH", b.getOperationCount("lag", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "lag", "PATCH") }, "LAG") {
		return diagnostics, operationsPerformed
	}
	// 24. eth_port_settings
	if !execute("PATCH", b.getOperationCount("eth_port_settings", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "eth_port_settings", "PATCH") }, "Eth Port Settings") {
		return diagnostics, operationsPerformed
	}
	// 25. diagnostics_profile
	if !execute("PATCH", b.getOperationCount("diagnostics_profile", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "diagnostics_profile", "PATCH") }, "Diagnostics Profile") {
		return diagnostics, operationsPerformed
	}
	// 26. gateway_profile
	if !execute("PATCH", b.getOperationCount("gateway_profile", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "gateway_profile", "PATCH") }, "Gateway Profile") {
		return diagnostics, operationsPerformed
	}
	// 27. diagnostics_port_profile
	if !execute("PATCH", b.getOperationCount("diagnostics_port_profile", "PATCH"), func(ctx context.Context) diag.Diagnostics {
		return b.ExecuteBulk(ctx, "diagnostics_port_profile", "PATCH")
	}, "Diagnostics Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 28. bundle
	if !execute("PATCH", b.getOperationCount("bundle", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "bundle", "PATCH") }, "Bundle") {
		return diagnostics, operationsPerformed
	}
	// 29. pod
	if !execute("PATCH", b.getOperationCount("pod", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "pod", "PATCH") }, "Pod") {
		return diagnostics, operationsPerformed
	}
	// 30. badge
	if !execute("PATCH", b.getOperationCount("badge", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "badge", "PATCH") }, "Badge") {
		return diagnostics, operationsPerformed
	}
	// 31. spine_plane
	if !execute("PATCH", b.getOperationCount("spine_plane", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "spine_plane", "PATCH") }, "Spine Plane") {
		return diagnostics, operationsPerformed
	}
	// 32. switchpoint
	if !execute("PATCH", b.getOperationCount("switchpoint", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "switchpoint", "PATCH") }, "Switchpoint") {
		return diagnostics, operationsPerformed
	}
	// 33. device_settings
	if !execute("PATCH", b.getOperationCount("device_settings", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "device_settings", "PATCH") }, "Device Settings") {
		return diagnostics, operationsPerformed
	}
	// 34. threshold
	if !execute("PATCH", b.getOperationCount("threshold", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "threshold", "PATCH") }, "Threshold") {
		return diagnostics, operationsPerformed
	}
	// 35. grouping_rule
	if !execute("PATCH", b.getOperationCount("grouping_rule", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "grouping_rule", "PATCH") }, "Grouping Rule") {
		return diagnostics, operationsPerformed
	}
	// 36. threshold_group
	if !execute("PATCH", b.getOperationCount("threshold_group", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "threshold_group", "PATCH") }, "Threshold Group") {
		return diagnostics, operationsPerformed
	}
	// 37. site
	if !execute("PATCH", b.getOperationCount("site", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "site", "PATCH") }, "Site") {
		return diagnostics, operationsPerformed
	}
	// 38. device_controller
	if !execute("PATCH", b.getOperationCount("device_controller", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "device_controller", "PATCH") }, "Device Controller") {
		return diagnostics, operationsPerformed
	}

	// DELETE operations - Reverse DC Order (note: sfp_breakout and site are skipped - they only support GET and PATCH)
	// 38. device_controller
	if !execute("DELETE", b.getOperationCount("device_controller", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "device_controller", "DELETE") }, "Device Controller") {
		return diagnostics, operationsPerformed
	}
	// 36. threshold_group
	if !execute("DELETE", b.getOperationCount("threshold_group", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "threshold_group", "DELETE") }, "Threshold Group") {
		return diagnostics, operationsPerformed
	}
	// 35. grouping_rule
	if !execute("DELETE", b.getOperationCount("grouping_rule", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "grouping_rule", "DELETE") }, "Grouping Rule") {
		return diagnostics, operationsPerformed
	}
	// 34. threshold
	if !execute("DELETE", b.getOperationCount("threshold", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "threshold", "DELETE") }, "Threshold") {
		return diagnostics, operationsPerformed
	}
	// 33. device_settings
	if !execute("DELETE", b.getOperationCount("device_settings", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "device_settings", "DELETE") }, "Device Settings") {
		return diagnostics, operationsPerformed
	}
	// 32. switchpoint
	if !execute("DELETE", b.getOperationCount("switchpoint", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "switchpoint", "DELETE") }, "Switchpoint") {
		return diagnostics, operationsPerformed
	}
	// 31. spine_plane
	if !execute("DELETE", b.getOperationCount("spine_plane", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "spine_plane", "DELETE") }, "Spine Plane") {
		return diagnostics, operationsPerformed
	}
	// 30. badge
	if !execute("DELETE", b.getOperationCount("badge", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "badge", "DELETE") }, "Badge") {
		return diagnostics, operationsPerformed
	}
	// 29. pod
	if !execute("DELETE", b.getOperationCount("pod", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "pod", "DELETE") }, "Pod") {
		return diagnostics, operationsPerformed
	}
	// 28. bundle
	if !execute("DELETE", b.getOperationCount("bundle", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "bundle", "DELETE") }, "Bundle") {
		return diagnostics, operationsPerformed
	}
	// 27. diagnostics_port_profile
	if !execute("DELETE", b.getOperationCount("diagnostics_port_profile", "DELETE"), func(ctx context.Context) diag.Diagnostics {
		return b.ExecuteBulk(ctx, "diagnostics_port_profile", "DELETE")
	}, "Diagnostics Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 26. gateway_profile
	if !execute("DELETE", b.getOperationCount("gateway_profile", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "gateway_profile", "DELETE") }, "Gateway Profile") {
		return diagnostics, operationsPerformed
	}
	// 25. diagnostics_profile
	if !execute("DELETE", b.getOperationCount("diagnostics_profile", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "diagnostics_profile", "DELETE") }, "Diagnostics Profile") {
		return diagnostics, operationsPerformed
	}
	// 24. eth_port_settings
	if !execute("DELETE", b.getOperationCount("eth_port_settings", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "eth_port_settings", "DELETE") }, "Eth Port Settings") {
		return diagnostics, operationsPerformed
	}
	// 23. lag
	if !execute("DELETE", b.getOperationCount("lag", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "lag", "DELETE") }, "LAG") {
		return diagnostics, operationsPerformed
	}
	// 22. gateway
	if !execute("DELETE", b.getOperationCount("gateway", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "gateway", "DELETE") }, "Gateway") {
		return diagnostics, operationsPerformed
	}
	// 21. sflow_collector
	if !execute("DELETE", b.getOperationCount("sflow_collector", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "sflow_collector", "DELETE") }, "SFlow Collector") {
		return diagnostics, operationsPerformed
	}
	// 20. packet_queue
	if !execute("DELETE", b.getOperationCount("packet_queue", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "packet_queue", "DELETE") }, "Packet Queue") {
		return diagnostics, operationsPerformed
	}
	// 19. eth_port_profile
	if !execute("DELETE", b.getOperationCount("eth_port_profile", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "eth_port_profile", "DELETE") }, "Eth Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 18. packet_broker
	if !execute("DELETE", b.getOperationCount("packet_broker", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "packet_broker", "DELETE") }, "Packet Broker") {
		return diagnostics, operationsPerformed
	}
	// 17. port_acl
	if !execute("DELETE", b.getOperationCount("port_acl", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "port_acl", "DELETE") }, "Port ACL") {
		return diagnostics, operationsPerformed
	}
	// 16. service
	if !execute("DELETE", b.getOperationCount("service", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "service", "DELETE") }, "Service") {
		return diagnostics, operationsPerformed
	}
	// 15. ipv6_list
	if !execute("DELETE", b.getOperationCount("ipv6_list", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv6_list", "DELETE") }, "IPv6 List Filter") {
		return diagnostics, operationsPerformed
	}
	// 14. ipv4_list
	if !execute("DELETE", b.getOperationCount("ipv4_list", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv4_list", "DELETE") }, "IPv4 List Filter") {
		return diagnostics, operationsPerformed
	}
	// 13. pb_routing
	if !execute("DELETE", b.getOperationCount("pb_routing", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "pb_routing", "DELETE") }, "PB Routing") {
		return diagnostics, operationsPerformed
	}
	// 12. tenant
	if !execute("DELETE", b.getOperationCount("tenant", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "tenant", "DELETE") }, "Tenant") {
		return diagnostics, operationsPerformed
	}
	// 11. pb_routing_acl
	if !execute("DELETE", b.getOperationCount("pb_routing_acl", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "pb_routing_acl", "DELETE") }, "PB Routing ACL") {
		return diagnostics, operationsPerformed
	}
	// 10. route_map
	if !execute("DELETE", b.getOperationCount("route_map", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "route_map", "DELETE") }, "Route Map") {
		return diagnostics, operationsPerformed
	}
	// 8-9. acl (both ipv4 and ipv6)
	if !execute("DELETE", b.getOperationCount("acl", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "acl", "DELETE") }, "ACL") {
		return diagnostics, operationsPerformed
	}
	// 7. route_map_clause
	if !execute("DELETE", b.getOperationCount("route_map_clause", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "route_map_clause", "DELETE") }, "Route Map Clause") {
		return diagnostics, operationsPerformed
	}
	// 6. as_path_access_list
	if !execute("DELETE", b.getOperationCount("as_path_access_list", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "as_path_access_list", "DELETE") }, "AS Path Access List") {
		return diagnostics, operationsPerformed
	}
	// 5. extended_community_list
	if !execute("DELETE", b.getOperationCount("extended_community_list", "DELETE"), func(ctx context.Context) diag.Diagnostics {
		return b.ExecuteBulk(ctx, "extended_community_list", "DELETE")
	}, "Extended Community List") {
		return diagnostics, operationsPerformed
	}
	// 4. ipv4_prefix_list
	if !execute("DELETE", b.getOperationCount("ipv4_prefix_list", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv4_prefix_list", "DELETE") }, "IPv4 Prefix List") {
		return diagnostics, operationsPerformed
	}
	// 3. community_list
	if !execute("DELETE", b.getOperationCount("community_list", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "community_list", "DELETE") }, "Community List") {
		return diagnostics, operationsPerformed
	}
	// 2. ipv6_prefix_list
	if !execute("DELETE", b.getOperationCount("ipv6_prefix_list", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv6_prefix_list", "DELETE") }, "IPv6 Prefix List") {
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

	// PUT operations - Campus Order (note: site is skipped - it only supports GET and PATCH)
	// 1. ipv4_list
	if !execute("PUT", b.getOperationCount("ipv4_list", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv4_list", "PUT") }, "IPv4 List") {
		return diagnostics, operationsPerformed
	}
	// 2. ipv6_list
	if !execute("PUT", b.getOperationCount("ipv6_list", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv6_list", "PUT") }, "IPv6 List") {
		return diagnostics, operationsPerformed
	}
	// 3-4. acl (both ipv4 and ipv6)
	if !execute("PUT", b.getOperationCount("acl", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "acl", "PUT") }, "ACL") {
		return diagnostics, operationsPerformed
	}
	// 5. pb_routing_acl
	if !execute("PUT", b.getOperationCount("pb_routing_acl", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "pb_routing_acl", "PUT") }, "PB Routing ACL") {
		return diagnostics, operationsPerformed
	}
	// 6. pb_routing
	if !execute("PUT", b.getOperationCount("pb_routing", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "pb_routing", "PUT") }, "PB Routing") {
		return diagnostics, operationsPerformed
	}
	// 7. port_acl
	if !execute("PUT", b.getOperationCount("port_acl", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "port_acl", "PUT") }, "Port ACL") {
		return diagnostics, operationsPerformed
	}
	// 8. service
	if !execute("PUT", b.getOperationCount("service", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "service", "PUT") }, "Service") {
		return diagnostics, operationsPerformed
	}
	// 9. eth_port_profile
	if !execute("PUT", b.getOperationCount("eth_port_profile", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "eth_port_profile", "PUT") }, "Eth Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 10. sflow_collector
	if !execute("PUT", b.getOperationCount("sflow_collector", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "sflow_collector", "PUT") }, "SFlow Collector") {
		return diagnostics, operationsPerformed
	}
	// 11. packet_queue
	if !execute("PUT", b.getOperationCount("packet_queue", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "packet_queue", "PUT") }, "Packet Queue") {
		return diagnostics, operationsPerformed
	}
	// 12. service_port_profile
	if !execute("PUT", b.getOperationCount("service_port_profile", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "service_port_profile", "PUT") }, "Service Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 13. diagnostics_port_profile
	if !execute("PUT", b.getOperationCount("diagnostics_port_profile", "PUT"), func(ctx context.Context) diag.Diagnostics {
		return b.ExecuteBulk(ctx, "diagnostics_port_profile", "PUT")
	}, "Diagnostics Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 14. device_voice_settings
	if !execute("PUT", b.getOperationCount("device_voice_settings", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "device_voice_settings", "PUT") }, "Device Voice Settings") {
		return diagnostics, operationsPerformed
	}
	// 15. authenticated_eth_port
	if !execute("PUT", b.getOperationCount("authenticated_eth_port", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "authenticated_eth_port", "PUT") }, "Authenticated Eth-Port") {
		return diagnostics, operationsPerformed
	}
	// 16. diagnostics_profile
	if !execute("PUT", b.getOperationCount("diagnostics_profile", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "diagnostics_profile", "PUT") }, "Diagnostics Profile") {
		return diagnostics, operationsPerformed
	}
	// 17. eth_port_settings
	if !execute("PUT", b.getOperationCount("eth_port_settings", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "eth_port_settings", "PUT") }, "Eth Port Settings") {
		return diagnostics, operationsPerformed
	}
	// 18. voice_port_profile
	if !execute("PUT", b.getOperationCount("voice_port_profile", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "voice_port_profile", "PUT") }, "Voice-Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 19. device_settings
	if !execute("PUT", b.getOperationCount("device_settings", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "device_settings", "PUT") }, "Device Settings") {
		return diagnostics, operationsPerformed
	}
	// 20. lag
	if !execute("PUT", b.getOperationCount("lag", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "lag", "PUT") }, "LAG") {
		return diagnostics, operationsPerformed
	}
	// 21. bundle
	if !execute("PUT", b.getOperationCount("bundle", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "bundle", "PUT") }, "Bundle") {
		return diagnostics, operationsPerformed
	}
	// 22. badge
	if !execute("PUT", b.getOperationCount("badge", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "badge", "PUT") }, "Badge") {
		return diagnostics, operationsPerformed
	}
	// 23. switchpoint
	if !execute("PUT", b.getOperationCount("switchpoint", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "switchpoint", "PUT") }, "Switchpoint") {
		return diagnostics, operationsPerformed
	}
	// 24. threshold
	if !execute("PUT", b.getOperationCount("threshold", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "threshold", "PUT") }, "Threshold") {
		return diagnostics, operationsPerformed
	}
	// 25. grouping_rule
	if !execute("PUT", b.getOperationCount("grouping_rule", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "grouping_rule", "PUT") }, "Grouping Rule") {
		return diagnostics, operationsPerformed
	}
	// 26. threshold_group
	if !execute("PUT", b.getOperationCount("threshold_group", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "threshold_group", "PUT") }, "Threshold Group") {
		return diagnostics, operationsPerformed
	}
	// 28. device_controller
	if !execute("PUT", b.getOperationCount("device_controller", "PUT"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "device_controller", "PUT") }, "Device Controller") {
		return diagnostics, operationsPerformed
	}

	// PATCH operations - Campus Order
	// 1. ipv4_list
	if !execute("PATCH", b.getOperationCount("ipv4_list", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv4_list", "PATCH") }, "IPv4 List") {
		return diagnostics, operationsPerformed
	}
	// 2. ipv6_list
	if !execute("PATCH", b.getOperationCount("ipv6_list", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv6_list", "PATCH") }, "IPv6 List") {
		return diagnostics, operationsPerformed
	}
	// 3-4. acl (both ipv4 and ipv6)
	if !execute("PATCH", b.getOperationCount("acl", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "acl", "PATCH") }, "ACL") {
		return diagnostics, operationsPerformed
	}
	// 5. pb_routing_acl
	if !execute("PATCH", b.getOperationCount("pb_routing_acl", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "pb_routing_acl", "PATCH") }, "PB Routing ACL") {
		return diagnostics, operationsPerformed
	}
	// 6. pb_routing
	if !execute("PATCH", b.getOperationCount("pb_routing", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "pb_routing", "PATCH") }, "PB Routing") {
		return diagnostics, operationsPerformed
	}
	// 7. port_acl
	if !execute("PATCH", b.getOperationCount("port_acl", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "port_acl", "PATCH") }, "Port ACL") {
		return diagnostics, operationsPerformed
	}
	// 8. service
	if !execute("PATCH", b.getOperationCount("service", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "service", "PATCH") }, "Service") {
		return diagnostics, operationsPerformed
	}
	// 9. eth_port_profile
	if !execute("PATCH", b.getOperationCount("eth_port_profile", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "eth_port_profile", "PATCH") }, "Eth Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 10. sflow_collector
	if !execute("PATCH", b.getOperationCount("sflow_collector", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "sflow_collector", "PATCH") }, "SFlow Collector") {
		return diagnostics, operationsPerformed
	}
	// 11. packet_queue
	if !execute("PATCH", b.getOperationCount("packet_queue", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "packet_queue", "PATCH") }, "Packet Queue") {
		return diagnostics, operationsPerformed
	}
	// 12. service_port_profile
	if !execute("PATCH", b.getOperationCount("service_port_profile", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "service_port_profile", "PATCH") }, "Service Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 13. diagnostics_port_profile
	if !execute("PATCH", b.getOperationCount("diagnostics_port_profile", "PATCH"), func(ctx context.Context) diag.Diagnostics {
		return b.ExecuteBulk(ctx, "diagnostics_port_profile", "PATCH")
	}, "Diagnostics Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 14. device_voice_settings
	if !execute("PATCH", b.getOperationCount("device_voice_settings", "PATCH"), func(ctx context.Context) diag.Diagnostics {
		return b.ExecuteBulk(ctx, "device_voice_settings", "PATCH")
	}, "Device Voice Settings") {
		return diagnostics, operationsPerformed
	}
	// 15. authenticated_eth_port
	if !execute("PATCH", b.getOperationCount("authenticated_eth_port", "PATCH"), func(ctx context.Context) diag.Diagnostics {
		return b.ExecuteBulk(ctx, "authenticated_eth_port", "PATCH")
	}, "Authenticated Eth-Port") {
		return diagnostics, operationsPerformed
	}
	// 16. diagnostics_profile
	if !execute("PATCH", b.getOperationCount("diagnostics_profile", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "diagnostics_profile", "PATCH") }, "Diagnostics Profile") {
		return diagnostics, operationsPerformed
	}
	// 17. eth_port_settings
	if !execute("PATCH", b.getOperationCount("eth_port_settings", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "eth_port_settings", "PATCH") }, "Eth Port Settings") {
		return diagnostics, operationsPerformed
	}
	// 18. voice_port_profile
	if !execute("PATCH", b.getOperationCount("voice_port_profile", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "voice_port_profile", "PATCH") }, "Voice-Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 19. device_settings
	if !execute("PATCH", b.getOperationCount("device_settings", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "device_settings", "PATCH") }, "Device Settings") {
		return diagnostics, operationsPerformed
	}
	// 20. lag
	if !execute("PATCH", b.getOperationCount("lag", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "lag", "PATCH") }, "LAG") {
		return diagnostics, operationsPerformed
	}
	// 21. bundle
	if !execute("PATCH", b.getOperationCount("bundle", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "bundle", "PATCH") }, "Bundle") {
		return diagnostics, operationsPerformed
	}
	// 22. badge
	if !execute("PATCH", b.getOperationCount("badge", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "badge", "PATCH") }, "Badge") {
		return diagnostics, operationsPerformed
	}
	// 23. switchpoint
	if !execute("PATCH", b.getOperationCount("switchpoint", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "switchpoint", "PATCH") }, "Switchpoint") {
		return diagnostics, operationsPerformed
	}
	// 24. threshold
	if !execute("PATCH", b.getOperationCount("threshold", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "threshold", "PATCH") }, "Threshold") {
		return diagnostics, operationsPerformed
	}
	// 25. grouping_rule
	if !execute("PATCH", b.getOperationCount("grouping_rule", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "grouping_rule", "PATCH") }, "Grouping Rule") {
		return diagnostics, operationsPerformed
	}
	// 26. threshold_group
	if !execute("PATCH", b.getOperationCount("threshold_group", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "threshold_group", "PATCH") }, "Threshold Group") {
		return diagnostics, operationsPerformed
	}
	// 27. site
	if !execute("PATCH", b.getOperationCount("site", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "site", "PATCH") }, "Site") {
		return diagnostics, operationsPerformed
	}
	// 28. device_controller
	if !execute("PATCH", b.getOperationCount("device_controller", "PATCH"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "device_controller", "PATCH") }, "Device Controller") {
		return diagnostics, operationsPerformed
	}

	// DELETE operations - Reverse Campus Order (note: site is skipped - it only supports GET and PATCH)
	// 28. device_controller
	if !execute("DELETE", b.getOperationCount("device_controller", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "device_controller", "DELETE") }, "Device Controller") {
		return diagnostics, operationsPerformed
	}
	// 26. threshold_group
	if !execute("DELETE", b.getOperationCount("threshold_group", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "threshold_group", "DELETE") }, "Threshold Group") {
		return diagnostics, operationsPerformed
	}
	// 25. grouping_rule
	if !execute("DELETE", b.getOperationCount("grouping_rule", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "grouping_rule", "DELETE") }, "Grouping Rule") {
		return diagnostics, operationsPerformed
	}
	// 24. threshold
	if !execute("DELETE", b.getOperationCount("threshold", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "threshold", "DELETE") }, "Threshold") {
		return diagnostics, operationsPerformed
	}
	// 23. switchpoint
	if !execute("DELETE", b.getOperationCount("switchpoint", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "switchpoint", "DELETE") }, "Switchpoint") {
		return diagnostics, operationsPerformed
	}
	// 22. badge
	if !execute("DELETE", b.getOperationCount("badge", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "badge", "DELETE") }, "Badge") {
		return diagnostics, operationsPerformed
	}
	// 21. bundle
	if !execute("DELETE", b.getOperationCount("bundle", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "bundle", "DELETE") }, "Bundle") {
		return diagnostics, operationsPerformed
	}
	// 20. lag
	if !execute("DELETE", b.getOperationCount("lag", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "lag", "DELETE") }, "LAG") {
		return diagnostics, operationsPerformed
	}
	// 19. device_settings
	if !execute("DELETE", b.getOperationCount("device_settings", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "device_settings", "DELETE") }, "Device Settings") {
		return diagnostics, operationsPerformed
	}
	// 18. voice_port_profile
	if !execute("DELETE", b.getOperationCount("voice_port_profile", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "voice_port_profile", "DELETE") }, "Voice-Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 17. eth_port_settings
	if !execute("DELETE", b.getOperationCount("eth_port_settings", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "eth_port_settings", "DELETE") }, "Eth Port Settings") {
		return diagnostics, operationsPerformed
	}
	// 16. diagnostics_profile
	if !execute("DELETE", b.getOperationCount("diagnostics_profile", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "diagnostics_profile", "DELETE") }, "Diagnostics Profile") {
		return diagnostics, operationsPerformed
	}
	// 15. authenticated_eth_port
	if !execute("DELETE", b.getOperationCount("authenticated_eth_port", "DELETE"), func(ctx context.Context) diag.Diagnostics {
		return b.ExecuteBulk(ctx, "authenticated_eth_port", "DELETE")
	}, "Authenticated Eth-Port") {
		return diagnostics, operationsPerformed
	}
	// 14. device_voice_settings
	if !execute("DELETE", b.getOperationCount("device_voice_settings", "DELETE"), func(ctx context.Context) diag.Diagnostics {
		return b.ExecuteBulk(ctx, "device_voice_settings", "DELETE")
	}, "Device Voice Settings") {
		return diagnostics, operationsPerformed
	}
	// 13. diagnostics_port_profile
	if !execute("DELETE", b.getOperationCount("diagnostics_port_profile", "DELETE"), func(ctx context.Context) diag.Diagnostics {
		return b.ExecuteBulk(ctx, "diagnostics_port_profile", "DELETE")
	}, "Diagnostics Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 12. service_port_profile
	if !execute("DELETE", b.getOperationCount("service_port_profile", "DELETE"), func(ctx context.Context) diag.Diagnostics {
		return b.ExecuteBulk(ctx, "service_port_profile", "DELETE")
	}, "Service Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 11. packet_queue
	if !execute("DELETE", b.getOperationCount("packet_queue", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "packet_queue", "DELETE") }, "Packet Queue") {
		return diagnostics, operationsPerformed
	}
	// 10. sflow_collector
	if !execute("DELETE", b.getOperationCount("sflow_collector", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "sflow_collector", "DELETE") }, "SFlow Collector") {
		return diagnostics, operationsPerformed
	}
	// 9. eth_port_profile
	if !execute("DELETE", b.getOperationCount("eth_port_profile", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "eth_port_profile", "DELETE") }, "Eth Port Profile") {
		return diagnostics, operationsPerformed
	}
	// 8. service
	if !execute("DELETE", b.getOperationCount("service", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "service", "DELETE") }, "Service") {
		return diagnostics, operationsPerformed
	}
	// 7. port_acl
	if !execute("DELETE", b.getOperationCount("port_acl", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "port_acl", "DELETE") }, "Port ACL") {
		return diagnostics, operationsPerformed
	}
	// 6. pb_routing
	if !execute("DELETE", b.getOperationCount("pb_routing", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "pb_routing", "DELETE") }, "PB Routing") {
		return diagnostics, operationsPerformed
	}
	// 5. pb_routing_acl
	if !execute("DELETE", b.getOperationCount("pb_routing_acl", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "pb_routing_acl", "DELETE") }, "PB Routing ACL") {
		return diagnostics, operationsPerformed
	}
	// 3-4. acl (both ipv4 and ipv6)
	if !execute("DELETE", b.getOperationCount("acl", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "acl", "DELETE") }, "ACL") {
		return diagnostics, operationsPerformed
	}
	// 2. ipv6_list
	if !execute("DELETE", b.getOperationCount("ipv6_list", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv6_list", "DELETE") }, "IPv6 List Filter") {
		return diagnostics, operationsPerformed
	}
	// 1. ipv4_list
	if !execute("DELETE", b.getOperationCount("ipv4_list", "DELETE"), func(ctx context.Context) diag.Diagnostics { return b.ExecuteBulk(ctx, "ipv4_list", "DELETE") }, "IPv4 List Filter") {
		return diagnostics, operationsPerformed
	}

	return diagnostics, operationsPerformed
}

func (b *BulkOperationManager) ShouldExecuteOperations(ctx context.Context) bool {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// If there are no pending operations, no need to execute
	if b.getOperationCount("gateway", "PUT") == 0 && b.getOperationCount("gateway", "PATCH") == 0 && b.getOperationCount("gateway", "DELETE") == 0 &&
		b.getOperationCount("lag", "PUT") == 0 && b.getOperationCount("lag", "PATCH") == 0 && b.getOperationCount("lag", "DELETE") == 0 &&
		b.getOperationCount("tenant", "PUT") == 0 && b.getOperationCount("tenant", "PATCH") == 0 && b.getOperationCount("tenant", "DELETE") == 0 &&
		b.getOperationCount("service", "PUT") == 0 && b.getOperationCount("service", "PATCH") == 0 && b.getOperationCount("service", "DELETE") == 0 &&
		b.getOperationCount("gateway_profile", "PUT") == 0 && b.getOperationCount("gateway_profile", "PATCH") == 0 && b.getOperationCount("gateway_profile", "DELETE") == 0 &&
		b.getOperationCount("eth_port_profile", "PUT") == 0 && b.getOperationCount("eth_port_profile", "PATCH") == 0 && b.getOperationCount("eth_port_profile", "DELETE") == 0 &&
		b.getOperationCount("eth_port_settings", "PUT") == 0 && b.getOperationCount("eth_port_settings", "PATCH") == 0 && b.getOperationCount("eth_port_settings", "DELETE") == 0 &&
		b.getOperationCount("device_settings", "PUT") == 0 && b.getOperationCount("device_settings", "PATCH") == 0 && b.getOperationCount("device_settings", "DELETE") == 0 &&
		b.getOperationCount("bundle", "PUT") == 0 && b.getOperationCount("bundle", "PATCH") == 0 && b.getOperationCount("bundle", "DELETE") == 0 && b.getOperationCount("authenticated_eth_port", "PUT") == 0 && b.getOperationCount("authenticated_eth_port", "PATCH") == 0 &&
		b.getOperationCount("authenticated_eth_port", "DELETE") == 0 && b.getOperationCount("acl", "PUT") == 0 && b.getOperationCount("acl", "PATCH") == 0 && b.getOperationCount("acl", "DELETE") == 0 &&
		b.getOperationCount("ipv4_list", "PUT") == 0 && b.getOperationCount("ipv4_list", "PATCH") == 0 && b.getOperationCount("ipv4_list", "DELETE") == 0 &&
		b.getOperationCount("ipv4_prefix_list", "PUT") == 0 && b.getOperationCount("ipv4_prefix_list", "PATCH") == 0 && b.getOperationCount("ipv4_prefix_list", "DELETE") == 0 &&
		b.getOperationCount("ipv6_list", "PUT") == 0 && b.getOperationCount("ipv6_list", "PATCH") == 0 && b.getOperationCount("ipv6_list", "DELETE") == 0 &&
		b.getOperationCount("ipv6_prefix_list", "PUT") == 0 && b.getOperationCount("ipv6_prefix_list", "PATCH") == 0 && b.getOperationCount("ipv6_prefix_list", "DELETE") == 0 &&
		b.getOperationCount("badge", "PUT") == 0 && b.getOperationCount("badge", "PATCH") == 0 && b.getOperationCount("badge", "DELETE") == 0 &&
		b.getOperationCount("voice_port_profile", "PUT") == 0 && b.getOperationCount("voice_port_profile", "PATCH") == 0 && b.getOperationCount("voice_port_profile", "DELETE") == 0 &&
		b.getOperationCount("switchpoint", "PUT") == 0 && b.getOperationCount("switchpoint", "PATCH") == 0 && b.getOperationCount("switchpoint", "DELETE") == 0 &&
		b.getOperationCount("service_port_profile", "PUT") == 0 && b.getOperationCount("service_port_profile", "PATCH") == 0 && b.getOperationCount("service_port_profile", "DELETE") == 0 &&
		b.getOperationCount("packet_broker", "PUT") == 0 && b.getOperationCount("packet_broker", "PATCH") == 0 && b.getOperationCount("packet_broker", "DELETE") == 0 &&
		b.getOperationCount("packet_queue", "PUT") == 0 && b.getOperationCount("packet_queue", "PATCH") == 0 && b.getOperationCount("packet_queue", "DELETE") == 0 &&
		b.getOperationCount("device_voice_settings", "PUT") == 0 && b.getOperationCount("device_voice_settings", "PATCH") == 0 && b.getOperationCount("device_voice_settings", "DELETE") == 0 &&
		b.getOperationCount("as_path_access_list", "PUT") == 0 && b.getOperationCount("as_path_access_list", "PATCH") == 0 && b.getOperationCount("as_path_access_list", "DELETE") == 0 &&
		b.getOperationCount("community_list", "PUT") == 0 && b.getOperationCount("community_list", "PATCH") == 0 && b.getOperationCount("community_list", "DELETE") == 0 &&
		b.getOperationCount("extended_community_list", "PUT") == 0 && b.getOperationCount("extended_community_list", "PATCH") == 0 && b.getOperationCount("extended_community_list", "DELETE") == 0 &&
		b.getOperationCount("route_map_clause", "PUT") == 0 && b.getOperationCount("route_map_clause", "PATCH") == 0 && b.getOperationCount("route_map_clause", "DELETE") == 0 &&
		b.getOperationCount("route_map", "PUT") == 0 && b.getOperationCount("route_map", "PATCH") == 0 && b.getOperationCount("route_map", "DELETE") == 0 &&
		b.getOperationCount("sfp_breakout", "PATCH") == 0 &&
		b.getOperationCount("site", "PATCH") == 0 &&
		b.getOperationCount("pod", "PUT") == 0 && b.getOperationCount("pod", "PATCH") == 0 && b.getOperationCount("pod", "DELETE") == 0 &&
		b.getOperationCount("pb_routing", "PUT") == 0 && b.getOperationCount("pb_routing", "PATCH") == 0 && b.getOperationCount("pb_routing", "DELETE") == 0 &&
		b.getOperationCount("pb_routing_acl", "PUT") == 0 && b.getOperationCount("pb_routing_acl", "PATCH") == 0 && b.getOperationCount("pb_routing_acl", "DELETE") == 0 &&
		b.getOperationCount("spine_plane", "PUT") == 0 && b.getOperationCount("spine_plane", "PATCH") == 0 && b.getOperationCount("spine_plane", "DELETE") == 0 &&
		b.getOperationCount("port_acl", "PUT") == 0 && b.getOperationCount("port_acl", "PATCH") == 0 && b.getOperationCount("port_acl", "DELETE") == 0 &&
		b.getOperationCount("sflow_collector", "PUT") == 0 && b.getOperationCount("sflow_collector", "PATCH") == 0 && b.getOperationCount("sflow_collector", "DELETE") == 0 &&
		b.getOperationCount("diagnostics_profile", "PUT") == 0 && b.getOperationCount("diagnostics_profile", "PATCH") == 0 && b.getOperationCount("diagnostics_profile", "DELETE") == 0 &&
		b.getOperationCount("diagnostics_port_profile", "PUT") == 0 && b.getOperationCount("diagnostics_port_profile", "PATCH") == 0 && b.getOperationCount("diagnostics_port_profile", "DELETE") == 0 &&
		b.getOperationCount("device_controller", "PUT") == 0 && b.getOperationCount("device_controller", "PATCH") == 0 && b.getOperationCount("device_controller", "DELETE") == 0 &&
		b.getOperationCount("grouping_rule", "PUT") == 0 && b.getOperationCount("grouping_rule", "PATCH") == 0 && b.getOperationCount("grouping_rule", "DELETE") == 0 &&
		b.getOperationCount("threshold_group", "PUT") == 0 && b.getOperationCount("threshold_group", "PATCH") == 0 && b.getOperationCount("threshold_group", "DELETE") == 0 &&
		b.getOperationCount("threshold", "PUT") == 0 && b.getOperationCount("threshold", "PATCH") == 0 && b.getOperationCount("threshold", "DELETE") == 0 {
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
	gatewayPutCount := b.getOperationCount("gateway", "PUT")
	gatewayPatchCount := b.getOperationCount("gateway", "PATCH")
	gatewayDeleteCount := b.getOperationCount("gateway", "DELETE")

	lagPutCount := b.getOperationCount("lag", "PUT")
	lagPatchCount := b.getOperationCount("lag", "PATCH")
	lagDeleteCount := b.getOperationCount("lag", "DELETE")

	tenantPutCount := b.getOperationCount("tenant", "PUT")
	tenantPatchCount := b.getOperationCount("tenant", "PATCH")
	tenantDeleteCount := b.getOperationCount("tenant", "DELETE")

	servicePutCount := b.getOperationCount("service", "PUT")
	servicePatchCount := b.getOperationCount("service", "PATCH")
	serviceDeleteCount := b.getOperationCount("service", "DELETE")

	gatewayProfilePutCount := b.getOperationCount("gateway_profile", "PUT")
	gatewayProfilePatchCount := b.getOperationCount("gateway_profile", "PATCH")
	gatewayProfileDeleteCount := b.getOperationCount("gateway_profile", "DELETE")

	ethPortProfilePutCount := b.getOperationCount("eth_port_profile", "PUT")
	ethPortProfilePatchCount := b.getOperationCount("eth_port_profile", "PATCH")
	ethPortProfileDeleteCount := b.getOperationCount("eth_port_profile", "DELETE")

	ethPortSettingsPutCount := b.getOperationCount("eth_port_settings", "PUT")
	ethPortSettingsPatchCount := b.getOperationCount("eth_port_settings", "PATCH")
	ethPortSettingsDeleteCount := b.getOperationCount("eth_port_settings", "DELETE")

	deviceSettingsPutCount := b.getOperationCount("device_settings", "PUT")
	deviceSettingsPatchCount := b.getOperationCount("device_settings", "PATCH")
	deviceSettingsDeleteCount := b.getOperationCount("device_settings", "DELETE")

	sflowCollectorPutCount := b.getOperationCount("sflow_collector", "PUT")
	sflowCollectorPatchCount := b.getOperationCount("sflow_collector", "PATCH")
	sflowCollectorDeleteCount := b.getOperationCount("sflow_collector", "DELETE")

	diagnosticsProfilePutCount := b.getOperationCount("diagnostics_profile", "PUT")
	diagnosticsProfilePatchCount := b.getOperationCount("diagnostics_profile", "PATCH")
	diagnosticsProfileDeleteCount := b.getOperationCount("diagnostics_profile", "DELETE")

	diagnosticsPortProfilePutCount := b.getOperationCount("diagnostics_port_profile", "PUT")
	diagnosticsPortProfilePatchCount := b.getOperationCount("diagnostics_port_profile", "PATCH")
	diagnosticsPortProfileDeleteCount := b.getOperationCount("diagnostics_port_profile", "DELETE")

	bundlePutCount := b.getOperationCount("bundle", "PUT")
	bundlePatchCount := b.getOperationCount("bundle", "PATCH")
	bundleDeleteCount := b.getOperationCount("bundle", "DELETE")

	aclPutCount := b.getOperationCount("acl", "PUT")
	aclPatchCount := b.getOperationCount("acl", "PATCH")
	aclDeleteCount := b.getOperationCount("acl", "DELETE")

	ipv4ListPutCount := b.getOperationCount("ipv4_list", "PUT")
	ipv4ListPatchCount := b.getOperationCount("ipv4_list", "PATCH")
	ipv4ListDeleteCount := b.getOperationCount("ipv4_list", "DELETE")

	ipv4PrefixListPutCount := b.getOperationCount("ipv4_prefix_list", "PUT")
	ipv4PrefixListPatchCount := b.getOperationCount("ipv4_prefix_list", "PATCH")
	ipv4PrefixListDeleteCount := b.getOperationCount("ipv4_prefix_list", "DELETE")

	ipv6ListPutCount := b.getOperationCount("ipv6_list", "PUT")
	ipv6ListPatchCount := b.getOperationCount("ipv6_list", "PATCH")
	ipv6ListDeleteCount := b.getOperationCount("ipv6_list", "DELETE")

	ipv6PrefixListPutCount := b.getOperationCount("ipv6_prefix_list", "PUT")
	ipv6PrefixListPatchCount := b.getOperationCount("ipv6_prefix_list", "PATCH")
	ipv6PrefixListDeleteCount := b.getOperationCount("ipv6_prefix_list", "DELETE")

	authenticatedEthPortPutCount := b.getOperationCount("authenticated_eth_port", "PUT")
	authenticatedEthPortPatchCount := b.getOperationCount("authenticated_eth_port", "PATCH")
	authenticatedEthPortDeleteCount := b.getOperationCount("authenticated_eth_port", "DELETE")

	badgePutCount := b.getOperationCount("badge", "PUT")
	badgePatchCount := b.getOperationCount("badge", "PATCH")
	badgeDeleteCount := b.getOperationCount("badge", "DELETE")

	voicePortProfilePutCount := b.getOperationCount("voice_port_profile", "PUT")
	voicePortProfilePatchCount := b.getOperationCount("voice_port_profile", "PATCH")
	voicePortProfileDeleteCount := b.getOperationCount("voice_port_profile", "DELETE")

	switchpointPutCount := b.getOperationCount("switchpoint", "PUT")
	switchpointPatchCount := b.getOperationCount("switchpoint", "PATCH")
	switchpointDeleteCount := b.getOperationCount("switchpoint", "DELETE")

	servicePortProfilePutCount := b.getOperationCount("service_port_profile", "PUT")
	servicePortProfilePatchCount := b.getOperationCount("service_port_profile", "PATCH")
	servicePortProfileDeleteCount := b.getOperationCount("service_port_profile", "DELETE")

	packetBrokerPutCount := b.getOperationCount("packet_broker", "PUT")
	packetBrokerPatchCount := b.getOperationCount("packet_broker", "PATCH")
	packetBrokerDeleteCount := b.getOperationCount("packet_broker", "DELETE")

	packetQueuePutCount := b.getOperationCount("packet_queue", "PUT")
	packetQueuePatchCount := b.getOperationCount("packet_queue", "PATCH")
	packetQueueDeleteCount := b.getOperationCount("packet_queue", "DELETE")

	deviceVoiceSettingsPutCount := b.getOperationCount("device_voice_settings", "PUT")
	deviceVoiceSettingsPatchCount := b.getOperationCount("device_voice_settings", "PATCH")
	deviceVoiceSettingsDeleteCount := b.getOperationCount("device_voice_settings", "DELETE")

	asPathAccessListPutCount := b.getOperationCount("as_path_access_list", "PUT")
	asPathAccessListPatchCount := b.getOperationCount("as_path_access_list", "PATCH")
	asPathAccessListDeleteCount := b.getOperationCount("as_path_access_list", "DELETE")

	communityListPutCount := b.getOperationCount("community_list", "PUT")
	communityListPatchCount := b.getOperationCount("community_list", "PATCH")
	communityListDeleteCount := b.getOperationCount("community_list", "DELETE")

	extendedCommunityListPutCount := b.getOperationCount("extended_community_list", "PUT")
	extendedCommunityListPatchCount := b.getOperationCount("extended_community_list", "PATCH")
	extendedCommunityListDeleteCount := b.getOperationCount("extended_community_list", "DELETE")

	routeMapClausePutCount := b.getOperationCount("route_map_clause", "PUT")
	routeMapClausePatchCount := b.getOperationCount("route_map_clause", "PATCH")
	routeMapClauseDeleteCount := b.getOperationCount("route_map_clause", "DELETE")

	routeMapPutCount := b.getOperationCount("route_map", "PUT")
	routeMapPatchCount := b.getOperationCount("route_map", "PATCH")
	routeMapDeleteCount := b.getOperationCount("route_map", "DELETE")

	sfpBreakoutPatchCount := b.getOperationCount("sfp_breakout", "PATCH")

	sitePatchCount := b.getOperationCount("site", "PATCH")

	podPutCount := b.getOperationCount("pod", "PUT")
	podPatchCount := b.getOperationCount("pod", "PATCH")
	podDeleteCount := b.getOperationCount("pod", "DELETE")

	pbRoutingPutCount := b.getOperationCount("pb_routing", "PUT")
	pbRoutingPatchCount := b.getOperationCount("pb_routing", "PATCH")
	pbRoutingDeleteCount := b.getOperationCount("pb_routing", "DELETE")

	pbRoutingAclPutCount := b.getOperationCount("pb_routing_acl", "PUT")
	pbRoutingAclPatchCount := b.getOperationCount("pb_routing_acl", "PATCH")
	pbRoutingAclDeleteCount := b.getOperationCount("pb_routing_acl", "DELETE")

	spinePlanePutCount := b.getOperationCount("spine_plane", "PUT")
	spinePlanePatchCount := b.getOperationCount("spine_plane", "PATCH")
	spinePlaneDeleteCount := b.getOperationCount("spine_plane", "DELETE")

	portAclPutCount := b.getOperationCount("port_acl", "PUT")
	portAclPatchCount := b.getOperationCount("port_acl", "PATCH")
	portAclDeleteCount := b.getOperationCount("port_acl", "DELETE")

	deviceControllerPutCount := b.getOperationCount("device_controller", "PUT")
	deviceControllerPatchCount := b.getOperationCount("device_controller", "PATCH")
	deviceControllerDeleteCount := b.getOperationCount("device_controller", "DELETE")

	groupingRulePutCount := b.getOperationCount("grouping_rule", "PUT")
	groupingRulePatchCount := b.getOperationCount("grouping_rule", "PATCH")
	groupingRuleDeleteCount := b.getOperationCount("grouping_rule", "DELETE")

	thresholdGroupPutCount := b.getOperationCount("threshold_group", "PUT")
	thresholdGroupPatchCount := b.getOperationCount("threshold_group", "PATCH")
	thresholdGroupDeleteCount := b.getOperationCount("threshold_group", "DELETE")

	thresholdPutCount := b.getOperationCount("threshold", "PUT")
	thresholdPatchCount := b.getOperationCount("threshold", "PATCH")
	thresholdDeleteCount := b.getOperationCount("threshold", "DELETE")

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
	// Handle special cases and aliases
	switch resourceType {
	case "acl_v4", "acl_v6":
		resourceType = "acl"
	}

	// Get the resource operations from the unified map
	res, exists := b.resources[resourceType]
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
	config, exists := resourceRegistry[g.resourceType]
	if !exists {
		return nil, fmt.Errorf("unknown resource type: %s", g.resourceType)
	}
	if config.PutFunc == nil {
		return nil, fmt.Errorf("PUT operation not supported for resource type: %s", g.resourceType)
	}
	return config.PutFunc(g.client, ctx, request)
}

func (g *GenericAPIClient) Patch(ctx context.Context, request interface{}) (*http.Response, error) {
	config, exists := resourceRegistry[g.resourceType]
	if !exists {
		return nil, fmt.Errorf("unknown resource type: %s", g.resourceType)
	}
	if config.PatchFunc == nil {
		return nil, fmt.Errorf("PATCH operation not supported for resource type: %s", g.resourceType)
	}
	return config.PatchFunc(g.client, ctx, request)
}

func (g *GenericAPIClient) Delete(ctx context.Context, names []string) (*http.Response, error) {
	config, exists := resourceRegistry[g.resourceType]
	if !exists {
		return nil, fmt.Errorf("unknown resource type: %s", g.resourceType)
	}
	if config.DeleteFunc == nil {
		return nil, fmt.Errorf("DELETE operation not supported for resource type: %s", g.resourceType)
	}
	return config.DeleteFunc(g.client, ctx, names)
}

func (g *GenericAPIClient) Get(ctx context.Context) (*http.Response, error) {
	config, exists := resourceRegistry[g.resourceType]
	if !exists {
		return nil, fmt.Errorf("unknown resource type: %s", g.resourceType)
	}
	if config.GetFunc == nil {
		return nil, fmt.Errorf("GET operation not supported for resource type: %s", g.resourceType)
	}
	return config.GetFunc(g.client, ctx)
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

	res, exists := b.resources["acl"]
	if !exists {
		b.mutex.Unlock()
		return diagnostics
	}

	var originalOperations map[string]interface{}
	var originalIpVersions map[string]string

	switch operationType {
	case "PUT":
		if res.Put == nil {
			b.mutex.Unlock()
			return diagnostics
		}
		originalOperations = make(map[string]interface{})
		for k, v := range res.Put {
			originalOperations[k] = v
		}
		res.Put = make(map[string]interface{})
	case "PATCH":
		if res.Patch == nil {
			b.mutex.Unlock()
			return diagnostics
		}
		originalOperations = make(map[string]interface{})
		for k, v := range res.Patch {
			originalOperations[k] = v
		}
		res.Patch = make(map[string]interface{})
	case "DELETE":
		if len(res.Delete) == 0 {
			b.mutex.Unlock()
			return diagnostics
		}
		originalOperations = make(map[string]interface{})
		for _, name := range res.Delete {
			originalOperations[name] = openapi.AclsPutRequestIpFilterValue{}
		}
		res.Delete = res.Delete[:0]
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
			ipv6Data[originalName] = props.(openapi.AclsPutRequestIpFilterValue)
		} else {
			ipv4Data[originalName] = props.(openapi.AclsPutRequestIpFilterValue)
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
	res.RecentOps = true
	res.RecentOpTime = time.Now()

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

		CheckPreExistence: func(ctx context.Context, resourceNames []string, originalOperations map[string]interface{}) ([]string, map[string]interface{}, error) {
			if operationType != "PUT" {
				return nil, nil, nil
			}

			checker := ResourceExistenceCheck{
				ResourceType:  "acl",
				OperationType: "PUT",
				FetchResources: func(ctx context.Context) (map[string]interface{}, error) {
					apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
					defer cancel()

					resp, err := b.client.ACLsAPI.AclsGet(apiCtx).IpVersion(ipVersion).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result map[string]interface{}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}

					// Extract the correct field based on IP version
					var filterKey string
					if ipVersion == "6" {
						filterKey = "ipv6_filter"
					} else {
						filterKey = "ipv4_filter"
					}

					if ipFilter, ok := result[filterKey].(map[string]interface{}); ok {
						return ipFilter, nil
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
				if val, ok := originalOperations[name]; ok {
					filteredOperations[name] = val
				}
			}

			return filteredNames, filteredOperations, nil
		},

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

	res, exists := b.resources[resourceType]
	if !exists {
		return
	}

	switch operationType {
	case "PUT":
		if res.Put == nil {
			res.Put = make(map[string]interface{})
		}
		res.Put[resourceName] = props
	case "PATCH":
		if res.Patch == nil {
			res.Patch = make(map[string]interface{})
		}
		res.Patch[resourceName] = props
	case "DELETE":
		res.Delete = append(res.Delete, resourceName)
	}

	res.RecentOps = true
	res.RecentOpTime = time.Now()
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
	return func() (map[string]interface{}, []string) {
		b.mutex.Lock()
		defer b.mutex.Unlock()

		res, exists := b.resources[resourceType]
		if !exists {
			return make(map[string]interface{}), []string{}
		}

		switch operationType {
		case "PUT":
			if res.Put == nil {
				return make(map[string]interface{}), []string{}
			}
			originalOperations := make(map[string]interface{})
			names := make([]string, 0, len(res.Put))
			for k, v := range res.Put {
				originalOperations[k] = v
				names = append(names, k)
			}
			// Clear the unified structure
			res.Put = make(map[string]interface{})
			return originalOperations, names

		case "PATCH":
			if res.Patch == nil {
				return make(map[string]interface{}), []string{}
			}
			originalOperations := make(map[string]interface{})
			names := make([]string, 0, len(res.Patch))
			for k, v := range res.Patch {
				originalOperations[k] = v
				names = append(names, k)
			}
			// Clear the unified structure
			res.Patch = make(map[string]interface{})
			return originalOperations, names

		case "DELETE":
			if len(res.Delete) == 0 {
				return make(map[string]interface{}), []string{}
			}
			names := make([]string, len(res.Delete))
			copy(names, res.Delete)
			result := make(map[string]interface{})
			for _, name := range names {
				result[name] = true
			}
			// Clear the unified structure
			res.Delete = res.Delete[:0]
			return result, names
		}

		return make(map[string]interface{}), []string{}
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
				apiCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
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

				case "ipv4_list":
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

				case "ipv6_list":
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

				case "pb_routing":
					resp, err := b.client.PBRoutingAPI.PolicybasedroutingGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						PbRouting map[string]interface{} `json:"pb_routing"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.PbRouting, nil

				case "pb_routing_acl":
					resp, err := b.client.PBRoutingACLAPI.PolicybasedroutingaclGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						PbRoutingAcl map[string]interface{} `json:"pb_routing_acl"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.PbRoutingAcl, nil

				case "spine_plane":
					resp, err := b.client.SpinePlanesAPI.SpineplanesGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						SpinePlane map[string]interface{} `json:"spine_plane"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.SpinePlane, nil

				case "grouping_rule":
					resp, err := b.client.GroupingRulesAPI.GroupingrulesGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						GroupingRule map[string]interface{} `json:"grouping_rule"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.GroupingRule, nil

				case "threshold_group":
					resp, err := b.client.ThresholdGroupsAPI.ThresholdgroupsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						ThresholdGroup map[string]interface{} `json:"threshold_group"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.ThresholdGroup, nil

				case "threshold":
					resp, err := b.client.ThresholdsAPI.ThresholdsGet(apiCtx).Execute()
					if err != nil {
						return nil, err
					}
					defer resp.Body.Close()

					var result struct {
						Threshold map[string]interface{} `json:"threshold"`
					}
					if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
						return nil, err
					}
					return result.Threshold, nil

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
		case "ipv4_list":
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
		case "ipv6_list":
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
			brokerMap := make(map[string]openapi.PacketbrokerPutRequestPortAclValue)
			for name, props := range filteredData {
				brokerMap[name] = props.(openapi.PacketbrokerPutRequestPortAclValue)
			}
			putRequest.SetPortAcl(brokerMap)
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
		case "pb_routing":
			putRequest := openapi.NewPolicybasedroutingPutRequest()
			pbRoutingMap := make(map[string]openapi.PolicybasedroutingPutRequestPbRoutingValue)
			for name, props := range filteredData {
				pbRoutingMap[name] = props.(openapi.PolicybasedroutingPutRequestPbRoutingValue)
			}
			putRequest.SetPbRouting(pbRoutingMap)
			return putRequest
		case "pb_routing_acl":
			putRequest := openapi.NewPolicybasedroutingaclPutRequest()
			pbRoutingAclMap := make(map[string]openapi.PolicybasedroutingaclPutRequestPbRoutingAclValue)
			for name, props := range filteredData {
				pbRoutingAclMap[name] = props.(openapi.PolicybasedroutingaclPutRequestPbRoutingAclValue)
			}
			putRequest.SetPbRoutingAcl(pbRoutingAclMap)
			return putRequest
		case "spine_plane":
			putRequest := openapi.NewSpineplanesPutRequest()
			spinePlaneMap := make(map[string]openapi.SpineplanesPutRequestSpinePlaneValue)
			for name, props := range filteredData {
				spinePlaneMap[name] = props.(openapi.SpineplanesPutRequestSpinePlaneValue)
			}
			putRequest.SetSpinePlane(spinePlaneMap)
			return putRequest
		case "grouping_rule":
			putRequest := openapi.NewGroupingrulesPutRequest()
			groupingRuleMap := make(map[string]openapi.GroupingrulesPutRequestGroupingRulesValue)
			for name, props := range filteredData {
				groupingRuleMap[name] = props.(openapi.GroupingrulesPutRequestGroupingRulesValue)
			}
			putRequest.SetGroupingRules(groupingRuleMap)
			return putRequest
		case "threshold_group":
			putRequest := openapi.NewThresholdgroupsPutRequest()
			thresholdGroupMap := make(map[string]openapi.ThresholdgroupsPutRequestThresholdGroupValue)
			for name, props := range filteredData {
				thresholdGroupMap[name] = props.(openapi.ThresholdgroupsPutRequestThresholdGroupValue)
			}
			putRequest.SetThresholdGroup(thresholdGroupMap)
			return putRequest
		case "threshold":
			putRequest := openapi.NewThresholdsPutRequest()
			thresholdMap := make(map[string]openapi.ThresholdsPutRequestThresholdValue)
			for name, props := range filteredData {
				thresholdMap[name] = props.(openapi.ThresholdsPutRequestThresholdValue)
			}
			putRequest.SetThreshold(thresholdMap)
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

		res, exists := b.resources[config.ResourceType]
		if !exists {
			return fmt.Errorf("resource type %s not found in unified structure", config.ResourceType)
		}

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

			res.ResponsesMutex.Lock()
			for tenantName, tenantData := range tenantsData.Tenant {
				res.Responses[tenantName] = tenantData
				if name, ok := tenantData["name"].(string); ok && name != tenantName {
					res.Responses[name] = tenantData
				}
			}
			res.ResponsesMutex.Unlock()

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

			res.ResponsesMutex.Lock()
			for serviceName, serviceData := range servicesData.Service {
				res.Responses[serviceName] = serviceData
				if name, ok := serviceData["name"].(string); ok && name != serviceName {
					res.Responses[name] = serviceData
				}
			}
			res.ResponsesMutex.Unlock()

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

			res.ResponsesMutex.Lock()
			for switchpointName, switchpointData := range switchpointsData.Switchpoint {
				res.Responses[switchpointName] = switchpointData
				if name, ok := switchpointData["name"].(string); ok && name != switchpointName {
					res.Responses[name] = switchpointData
				}
			}
			res.ResponsesMutex.Unlock()

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

			res.ResponsesMutex.Lock()
			for siteName, siteData := range sitesData.Site {
				res.Responses[siteName] = siteData
				if name, ok := siteData["name"].(string); ok && name != siteName {
					res.Responses[name] = siteData
				}
			}
			res.ResponsesMutex.Unlock()

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

		if res, exists := b.resources[resourceType]; exists {
			res.RecentOps = true
			res.RecentOpTime = now
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
