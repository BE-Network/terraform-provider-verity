package utils

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"terraform-provider-verity/internal/registry"
	"terraform-provider-verity/internal/spec"
)

var (
	registryOnce       sync.Once
	registryByType     map[string]spec.ResourceSpec
	registryByEndpoint map[string]spec.ResourceSpec
)

func loadRegistryIndexes() {
	registryOnce.Do(func() {
		registryByType = map[string]spec.ResourceSpec{}
		registryByEndpoint = map[string]spec.ResourceSpec{}
		resources, err := registry.Load()
		if err != nil {
			return
		}
		for _, resource := range resources {
			registryByType[resource.TerraformType] = resource
			endpoint := strings.Trim(resource.API.EndpointPath, "/")
			if _, shared := registryByEndpoint[endpoint]; !shared {
				registryByEndpoint[endpoint] = resource
			}
		}
	})
}

func registryResource(terraformType string) (spec.ResourceSpec, bool) {
	loadRegistryIndexes()
	resource, found := registryByType[terraformType]
	return resource, found
}

func registryResourceByEndpoint(endpoint string) (spec.ResourceSpec, bool) {
	loadRegistryIndexes()
	resource, found := registryByEndpoint[strings.Trim(endpoint, "/")]
	return resource, found
}

func specModeSet(modes []spec.Mode) map[string]bool {
	set := make(map[string]bool, len(modes))
	for _, mode := range modes {
		set[string(mode)] = true
	}
	return set
}

func ResponseCollectionKeyForBulkKey(bulkKey string) string {
	loadRegistryIndexes()
	key := ""
	for _, resource := range registryByType {
		if resource.API.BulkKey != bulkKey {
			continue
		}
		if key != "" && key != resource.API.ResponseCollectionKey {
			return ""
		}
		key = resource.API.ResponseCollectionKey
	}
	return key
}

func ResponseCollectionKeyForEndpoint(endpoint string) string {
	resource, found := registryResourceByEndpoint(endpoint)
	if !found {
		return ""
	}
	return resource.API.ResponseCollectionKey
}

func ResponseCollectionKeyForType(terraformType string) string {
	resource, found := registryResource(terraformType)
	if !found {
		return ""
	}
	return resource.API.ResponseCollectionKey
}

func HeaderSplitKeyForBulkKey(bulkKey string) (string, error) {
	loadRegistryIndexes()
	resources := make([]spec.ResourceSpec, 0, len(registryByType))
	for _, resource := range registryByType {
		resources = append(resources, resource)
	}
	return HeaderSplitKey(resources, bulkKey)
}

func HeaderSplitKey(resources []spec.ResourceSpec, bulkKey string) (string, error) {
	found := false
	splitKey := ""
	for _, resource := range resources {
		if resource.API.BulkKey != bulkKey {
			continue
		}
		found = true
		if len(resource.API.FixedHeaders) > 1 {
			headers := make([]string, 0, len(resource.API.FixedHeaders))
			for header := range resource.API.FixedHeaders {
				headers = append(headers, header)
			}
			sort.Strings(headers)
			return "", fmt.Errorf("bulk key %q supports one split key, %s has %s", bulkKey, resource.TerraformType, strings.Join(headers, ", "))
		}
		for header := range resource.API.FixedHeaders {
			if splitKey != "" && splitKey != header {
				return "", fmt.Errorf("bulk key %q supports one split key, resources disagree: %s and %s", bulkKey, splitKey, header)
			}
			splitKey = header
		}
	}
	if !found {
		return "", fmt.Errorf("the resource registry describes no bulk key %q", bulkKey)
	}
	return splitKey, nil
}
