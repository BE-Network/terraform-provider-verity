package bulkops

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"sync"
	"terraform-provider-verity/openapi"
	"time"
)

type ContextProviderFunc func() interface{}

type ClearCacheFunc func(ctx context.Context, provider interface{}, cacheKey string)

type OperationStatus int

type Operation struct {
	ResourceType       string
	ResourceName       string
	OperationType      string
	Status             OperationStatus
	Error              error
	ExecutionStartTime time.Time
}

type ResourceExistenceCheck struct {
	FetchResources func(ctx context.Context) (map[string]interface{}, error)
	ResourceType   string
	OperationType  string
}

type BulkOperationConfig struct {
	ResourceType            string
	OperationType           string
	ExtractOperations       func() (map[string]interface{}, []string)
	CheckPreExistence       func(ctx context.Context, resourceNames []string, originalOperations map[string]interface{}) ([]string, map[string]interface{}, error)
	PrepareRequest          func(filteredData map[string]interface{}) interface{}
	PrepareRequestWithError func(filteredData map[string]interface{}) (interface{}, error)
	ExecuteRequest          func(ctx context.Context, request interface{}) (*http.Response, error)
	ProcessResponse         func(ctx context.Context, resp *http.Response, operations map[string]interface{}) error
	UpdateRecentOps         func()
}

func (c BulkOperationConfig) prepareRequest(filteredData map[string]interface{}) (interface{}, error) {
	if c.PrepareRequestWithError != nil {
		return c.PrepareRequestWithError(filteredData)
	}
	if c.PrepareRequest == nil {
		return nil, fmt.Errorf("no request preparer configured for %s %s", c.ResourceType, c.OperationType)
	}
	return c.PrepareRequest(filteredData), nil
}

type ResourceOperations struct {
	Put            map[string]interface{}
	Patch          map[string]interface{}
	Delete         []string
	RecentOps      bool
	RecentOpTime   time.Time
	Responses      map[string]map[string]interface{}
	ResponsesMutex sync.RWMutex
}

func NewResourceOperations() *ResourceOperations {
	return &ResourceOperations{
		Put:       make(map[string]interface{}),
		Patch:     make(map[string]interface{}),
		Delete:    make([]string, 0),
		RecentOps: false,
		Responses: make(map[string]map[string]interface{}),
	}
}

type ResourceConfig struct {
	ResourceType     string
	PutRequestType   reflect.Type
	PatchRequestType reflect.Type
	APIClientGetter  func(*openapi.APIClient) ResourceAPIClient
	PutFunc          func(*openapi.APIClient, context.Context, interface{}) (*http.Response, error)
	PatchFunc        func(*openapi.APIClient, context.Context, interface{}) (*http.Response, error)
	DeleteFunc       func(*openapi.APIClient, context.Context, []string) (*http.Response, error)
	GetFunc          func(*openapi.APIClient, context.Context) (*http.Response, error)

	HeaderSplitKey string

	HeaderPutFunc    func(*openapi.APIClient, context.Context, interface{}, map[string]string) (*http.Response, error)
	HeaderPatchFunc  func(*openapi.APIClient, context.Context, interface{}, map[string]string) (*http.Response, error)
	HeaderDeleteFunc func(*openapi.APIClient, context.Context, []string, map[string]string) (*http.Response, error)
	HeaderGetFunc    func(*openapi.APIClient, context.Context, map[string]string) (*http.Response, error)

	HeaderResponseExtractor func(rawResponse map[string]interface{}, headers map[string]string) (map[string]interface{}, error)
}

type GenericAPIClient struct {
	client       *openapi.APIClient
	resourceType string
}

type ResourceAPIClient interface {
	Put(ctx context.Context, request interface{}) (*http.Response, error)
	Patch(ctx context.Context, request interface{}) (*http.Response, error)
	Delete(ctx context.Context, names []string) (*http.Response, error)
	Get(ctx context.Context) (*http.Response, error)
}

const (
	OperationPending OperationStatus = iota
	OperationExecuting
	OperationSucceeded
	OperationFailed
)

const (
	MaxBatchSize       = 1000
	MaxDeleteBatchSize = 100
	OperationTimeout   = 300 * time.Second
)

var (
	DefaultBatchDelay                = parseDuration("VERITY_DEFAULT_BATCH_DELAY", 2*time.Second)
	BatchCollectionWindow            = parseDuration("VERITY_BATCH_COLLECTION_WINDOW", 2000*time.Millisecond)
	MaxBatchDelay                    = parseDuration("VERITY_MAX_BATCH_DELAY", 5*time.Second)
	ResponseProcessorDelay           = parseDuration("VERITY_RESPONSE_PROCESSOR_DELAY", 5*time.Second)
	PostOperationVerificationRetries = parseInt("VERITY_POST_OPERATION_VERIFICATION_RETRIES", 2)
	PostOperationVerificationBackoff = parseDuration("VERITY_POST_OPERATION_VERIFICATION_BACKOFF", 2500*time.Millisecond)
	DebounceDelay                    = parseDuration("VERITY_DEBOUNCE_DELAY", 15*time.Second)
)

func parseDuration(envVar string, defaultVal time.Duration) time.Duration {
	if v := os.Getenv(envVar); v != "" {
		if d, err := time.ParseDuration(v); err == nil {
			return d
		}
	}
	return defaultVal
}

func parseInt(envVar string, defaultVal int) int {
	if v := os.Getenv(envVar); v != "" {
		if parsed, err := strconv.Atoi(v); err == nil && parsed >= 0 {
			return parsed
		}
	}
	return defaultVal
}

type ResourceOperationData struct {
	PutOperations    interface{}
	PatchOperations  interface{}
	DeleteOperations *[]string
	RecentOps        *bool
	RecentOpTime     *time.Time
}
