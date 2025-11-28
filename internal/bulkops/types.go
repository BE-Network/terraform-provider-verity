package bulkops

import (
	"context"
	"net/http"
	"reflect"
	"sync"
	"terraform-provider-verity/openapi"
	"time"
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

	// HeaderSplitKey specifies which header param to use for splitting operations into separate batches
	// Example: "ip_version" for ACLs (splits into IPv4/IPv6 batches)
	HeaderSplitKey string

	// HeaderAwareFuncs provide header-aware API call functions that accept header params
	HeaderPutFunc    func(*openapi.APIClient, context.Context, interface{}, map[string]string) (*http.Response, error)
	HeaderPatchFunc  func(*openapi.APIClient, context.Context, interface{}, map[string]string) (*http.Response, error)
	HeaderDeleteFunc func(*openapi.APIClient, context.Context, []string, map[string]string) (*http.Response, error)
	HeaderGetFunc    func(*openapi.APIClient, context.Context, map[string]string) (*http.Response, error)

	// HeaderResponseExtractor extracts resource data from GET response based on header values
	// Used when response format differs based on headers (e.g., ACL returns ipv4_filter or ipv6_filter)
	HeaderResponseExtractor func(rawResponse map[string]interface{}, headers map[string]string) (map[string]interface{}, error)
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

// ResourceOperationData holds operation data for a resource type
type ResourceOperationData struct {
	PutOperations    interface{}
	PatchOperations  interface{}
	DeleteOperations *[]string
	RecentOps        *bool
	RecentOpTime     *time.Time
}
