package provider

import (
	"context"
	"fmt"

	"terraform-provider-verity/internal/bulkops"
	"terraform-provider-verity/internal/genericresource"
	"terraform-provider-verity/internal/spec"
	"terraform-provider-verity/internal/transport"
	"terraform-provider-verity/internal/utils"
)

type genericRuntime struct {
	provCtx *providerContext
}

func bindGenericRuntime(providerData interface{}) (genericresource.Runtime, error) {
	provCtx, ok := providerData.(*providerContext)
	if !ok {
		return nil, fmt.Errorf("expected *providerContext, got: %T", providerData)
	}
	return genericRuntime{provCtx: provCtx}, nil
}

func (g genericRuntime) Mode() string { return g.provCtx.mode }

func (g genericRuntime) EnsureAuthenticated(ctx context.Context) error {
	return ensureAuthenticated(ctx, g.provCtx)
}

func (g genericRuntime) ClearCache(ctx context.Context, cacheKey string) {
	clearCache(ctx, g.provCtx, cacheKey)
}

func (g genericRuntime) BulkManager() *bulkops.Manager { return g.provCtx.bulkOpsMgr }

func (g genericRuntime) NotifyOperationAdded() { g.provCtx.NotifyOperationAdded() }

func (g genericRuntime) ConfiguredAttributes(ctx context.Context, terraformType, resourceName string) *utils.ConfiguredAttributes {
	return utils.ParseResourceConfiguredAttributes(ctx, g.provCtx.workDir, terraformType, resourceName)
}

func (g genericRuntime) FetchCollection(ctx context.Context, resourceSpec spec.ResourceSpec, resourceName string) (map[string]interface{}, error) {
	return utils.FetchResourceWithRetry(ctx, g.provCtx, resourceSpec.API.CacheKey, resourceName,
		func() (map[string]interface{}, error) { return g.get(ctx, resourceSpec) },
		getCachedResponse,
	)
}

func (g genericRuntime) get(ctx context.Context, resourceSpec spec.ResourceSpec) (map[string]interface{}, error) {
	return transport.FetchCollection(ctx, g.provCtx.client, resourceSpec.TerraformType,
		resourceSpec.API.EndpointPath, resourceSpec.API.FixedHeaders, resourceSpec.API.ResponseCollectionKey)
}
