package transport

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

	"terraform-provider-verity/openapi"
)

func FetchCollection(ctx context.Context, client *openapi.APIClient, label, endpointPath string, fixedHeaders map[string]string, collectionKey string) (map[string]interface{}, error) {
	config := client.GetConfig()
	base, err := config.ServerURLWithContext(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("resolve server URL: %w", err)
	}
	endpoint := strings.TrimSuffix(base, "/") + endpointPath

	if len(fixedHeaders) != 0 {
		query := url.Values{}
		for name, value := range fixedHeaders {
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

	tflog.Debug(ctx, fmt.Sprintf("Generic API request for %s: %s %s", label, request.Method, request.URL.EscapedPath()))
	response, err := config.HTTPClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %w", label, err)
	}
	defer response.Body.Close()
	if config.Debug {
		dump, err := httputil.DumpResponse(response, true)
		if err != nil {
			return nil, fmt.Errorf("dump %s response: %w", label, err)
		}
		log.Printf("\n%s\n", string(dump))
	}
	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return nil, fmt.Errorf("error reading %s: unexpected status %s", label, response.Status)
	}

	var decoded map[string]interface{}
	if err := json.NewDecoder(response.Body).Decode(&decoded); err != nil {
		return nil, fmt.Errorf("failed to decode %s response: %w", label, err)
	}
	collection, ok := decoded[collectionKey].(map[string]interface{})
	if !ok {
		if _, present := decoded[collectionKey]; !present {
			return map[string]interface{}{}, nil
		}
		return nil, fmt.Errorf("%s response key %q is not an object", label, collectionKey)
	}
	return collection, nil
}
