package genericresource

import (
	"context"

	"terraform-provider-verity/internal/bulkops"
	"terraform-provider-verity/internal/spec"
	"terraform-provider-verity/internal/transport"
	"terraform-provider-verity/internal/utils"
)

type Runtime interface {
	Mode() string
	EnsureAuthenticated(ctx context.Context) error

	ClearCache(ctx context.Context, cacheKey string)
	BulkManager() *bulkops.Manager
	NotifyOperationAdded()

	FetchCollection(ctx context.Context, resource spec.ResourceSpec, resourceName string) (map[string]interface{}, error)

	ConfiguredAttributes(ctx context.Context, terraformType, resourceName string) *utils.ConfiguredAttributes
}

type TransportAdapter interface {
	ResourceValue(object transport.WireObject) (interface{}, error)
}
