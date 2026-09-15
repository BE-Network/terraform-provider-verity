package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"terraform-provider-verity/internal/bulkops"
	"terraform-provider-verity/internal/genericresource"
	"terraform-provider-verity/internal/spec"
	"terraform-provider-verity/internal/utils"
)

// genericRuntime adapts the provider's session to what the generic engine needs.
//
// The engine takes an interface so it does not depend on this package; this is
// the only place the two meet, and it is deliberately thin — every method here
// forwards to the same helper the handwritten resources call, so a migrated
// resource shares their caching, retry, and authentication behavior rather than
// getting a parallel implementation of it.
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

// FetchCollection reads a resource's endpoint and returns the objects under its
// response collection key, through the same cache and retry the handwritten
// resources use.
//
// The request is issued directly rather than through a generated SDK method. The
// typed methods are one per endpoint, so reaching them generically would need a
// table mapping every resource to its own function — the duplication the
// registry exists to remove. A GET has no request body to type, so the endpoint
// path from the spec is enough to build it, and the response is decoded as data
// because that is what the engine decodes against the spec anyway.
func (g genericRuntime) FetchCollection(ctx context.Context, resourceSpec spec.ResourceSpec, resourceName string) (map[string]interface{}, error) {
	return utils.FetchResourceWithRetry(ctx, g.provCtx, resourceSpec.API.CacheKey, resourceName,
		func() (map[string]interface{}, error) { return g.get(ctx, resourceSpec) },
		getCachedResponse,
	)
}

func (g genericRuntime) get(ctx context.Context, resourceSpec spec.ResourceSpec) (map[string]interface{}, error) {
	config := g.provCtx.client.GetConfig()
	base, err := config.ServerURLWithContext(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("resolve server URL: %w", err)
	}
	endpoint := strings.TrimSuffix(base, "/") + resourceSpec.API.EndpointPath

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("build request for %s: %w", endpoint, err)
	}
	request.Header.Set("Accept", "application/json")
	for name, value := range config.DefaultHeader {
		request.Header.Set(name, value)
	}
	// One endpoint can back several resources, discriminated by a fixed header;
	// ACLs are the case, separated by ip_version.
	for name, value := range resourceSpec.API.FixedHeaders {
		request.Header.Set(name, value)
	}

	response, err := config.HTTPClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %w", resourceSpec.TerraformType, err)
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("error reading %s: unexpected status %s", resourceSpec.TerraformType, response.Status)
	}

	var decoded map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return nil, fmt.Errorf("failed to decode %s response: %w", resourceSpec.TerraformType, err)
	}
	collection, ok := decoded[resourceSpec.API.ResponseCollectionKey].(map[string]interface{})
	if !ok {
		// An endpoint that returns no objects at all is a valid empty collection,
		// not a malformed response.
		if _, present := decoded[resourceSpec.API.ResponseCollectionKey]; !present {
			return map[string]interface{}{}, nil
		}
		return nil, fmt.Errorf("%s response key %q is not an object",
			resourceSpec.TerraformType, resourceSpec.API.ResponseCollectionKey)
	}
	return collection, nil
}
