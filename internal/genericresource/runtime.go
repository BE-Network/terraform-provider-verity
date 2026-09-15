package genericresource

import (
	"context"

	"terraform-provider-verity/internal/bulkops"
	"terraform-provider-verity/internal/spec"
	"terraform-provider-verity/internal/transport"
)

// Runtime is what the generic engine needs from the provider: the ambient facts
// of a configured session and the two shared services every resource uses.
//
// It is an interface rather than the provider's own context type so the engine
// does not depend on the package it is registered from, and so a test can drive
// a lifecycle without standing up a provider.
type Runtime interface {
	// Mode is the operation mode the session is configured for. It decides which
	// fields apply, so it is read on every decode.
	Mode() string
	EnsureAuthenticated(ctx context.Context) error
	// ClearCache drops the cached collection for a key after a write, so the next
	// read sees the change rather than the response that preceded it.
	ClearCache(ctx context.Context, cacheKey string)
	BulkManager() *bulkops.Manager
	NotifyOperationAdded()
	// FetchCollection returns the objects at a resource's endpoint, keyed by their
	// API name, going through the provider's cache and retry policy.
	FetchCollection(ctx context.Context, resource spec.ResourceSpec, resourceName string) (map[string]interface{}, error)
}

// TransportAdapter converts the engine's canonical object into the typed value
// the bulk manager asserts for this resource.
//
// The bulk manager type-asserts the value it is given, so the canonical object
// cannot reach it directly. The plan's design keeps that boundary and generates
// an adapter per resource to cross it; this is the interface those adapters
// implement, with the IPv4 List one written by hand as the pilot.
type TransportAdapter interface {
	ResourceValue(object transport.WireObject) (interface{}, error)
}
