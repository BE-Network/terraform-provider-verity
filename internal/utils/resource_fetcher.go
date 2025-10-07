package utils

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// FetchResourceWithRetry implements the universal 3-retry pattern for fetching resources with caching
func FetchResourceWithRetry[T any](
	ctx context.Context,
	provCtx interface{},
	cacheKey, resourceName string,
	fetchFunc func() (T, error),
	getCachedResponseFunc func(context.Context, interface{}, string, func() (interface{}, error), ...bool) (interface{}, error),
) (T, error) {
	var result T
	var err error
	maxRetries := 3

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			sleepTime := time.Duration(100*(attempt+1)) * time.Millisecond
			tflog.Debug(ctx, fmt.Sprintf("Failed to fetch %s on attempt %d, retrying in %v", resourceName, attempt, sleepTime))
			time.Sleep(sleepTime)
		}

		data, fetchErr := getCachedResponseFunc(ctx, provCtx, cacheKey, func() (interface{}, error) {
			return fetchFunc()
		})

		if fetchErr == nil {
			if typedResult, ok := data.(T); ok {
				result = typedResult
				break
			} else {
				err = fmt.Errorf("failed to cast result to expected type")
			}
		}
		err = fetchErr
	}

	return result, err
}

// FindResourceByAPIName searches for a resource by its API name only (not by resource ID)
// This addresses the issue where resource IDs get sanitized for HCL compatibility
func FindResourceByAPIName[T any](
	resources map[string]T,
	targetName string,
	nameExtractor func(T) (string, bool),
) (resource T, actualAPIName string, exists bool) {
	for apiName, res := range resources {
		if extractedName, ok := nameExtractor(res); ok && extractedName == targetName {
			return res, apiName, true
		}
	}
	return resource, "", false
}
