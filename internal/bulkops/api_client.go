package bulkops

import (
	"context"
	"fmt"
	"net/http"
)

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
