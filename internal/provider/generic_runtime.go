package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-verity/internal/bulkops"
	"terraform-provider-verity/internal/genericresource"
	"terraform-provider-verity/internal/spec"
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
	config := g.provCtx.client.GetConfig()
	base, err := config.ServerURLWithContext(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("resolve server URL: %w", err)
	}
	endpoint := strings.TrimSuffix(base, "/") + resourceSpec.API.EndpointPath

	if len(resourceSpec.API.FixedHeaders) != 0 {
		query := url.Values{}
		for name, value := range resourceSpec.API.FixedHeaders {
			query.Set(name, value)
		}
		endpoint += "?" + query.Encode()
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("build request for %s: %w", endpoint, err)
	}
	request.Header.Set("Accept", "application/json")
	for name, value := range config.DefaultHeader {
		request.Header.Set(name, value)
	}

	tflog.Debug(ctx, fmt.Sprintf("Generic API request for %s: %s %s",
		resourceSpec.TerraformType, request.Method, request.URL.EscapedPath()))
	response, err := config.HTTPClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %w", resourceSpec.TerraformType, err)
	}
	defer response.Body.Close()
	if config.Debug {

		dump, err := httputil.DumpResponse(response, true)
		if err != nil {
			return nil, fmt.Errorf("dump %s response: %w", resourceSpec.TerraformType, err)
		}
		log.Printf("\n%s\n", string(dump))
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("error reading %s: unexpected status %s", resourceSpec.TerraformType, response.Status)
	}

	var decoded map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return nil, fmt.Errorf("failed to decode %s response: %w", resourceSpec.TerraformType, err)
	}
	collection, ok := decoded[resourceSpec.API.ResponseCollectionKey].(map[string]interface{})
	if !ok {

		if _, present := decoded[resourceSpec.API.ResponseCollectionKey]; !present {
			return map[string]interface{}{}, nil
		}
		return nil, fmt.Errorf("%s response key %q is not an object",
			resourceSpec.TerraformType, resourceSpec.API.ResponseCollectionKey)
	}
	return collection, nil
}
