package bulkops

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"terraform-provider-verity/openapi"
)

type readinessOperation struct {
	Flag     *bool `json:"flag,omitempty"`
	Nullable *bool `json:"nullable,omitempty"`
}

func TestMissingSubmittedResponseKeys(t *testing.T) {
	flag := false
	operations := map[string]interface{}{"one": readinessOperation{Flag: &flag}}
	missing := missingSubmittedResponseKeys(operations, map[string]interface{}{"one": map[string]interface{}{}})
	if got := missing["one"]; len(got) != 1 || got[0] != "flag" {
		t.Fatalf("missing keys = %#v", missing)
	}
}

func TestMissingSubmittedResponseKeysTreatsNullAsPresent(t *testing.T) {
	flag := true
	operations := map[string]interface{}{"one": readinessOperation{Flag: &flag}}
	response := map[string]interface{}{"one": map[string]interface{}{"flag": nil}}
	if missing := missingSubmittedResponseKeys(operations, response); len(missing) != 0 {
		t.Fatalf("explicit null must be present, got %#v", missing)
	}
}

func TestMissingSubmittedResponseKeysIgnoresRequestNull(t *testing.T) {
	operations := map[string]interface{}{"one": readinessOperation{}}
	if missing := missingSubmittedResponseKeys(operations, map[string]interface{}{"one": map[string]interface{}{}}); len(missing) != 0 {
		t.Fatalf("request null must not be expected, got %#v", missing)
	}
}

func TestMissingSubmittedResponseKeysUsesResponseNameAlias(t *testing.T) {
	flag := true
	operations := map[string]interface{}{"terraform-name": readinessOperation{Flag: &flag}}
	response := map[string]interface{}{"api-key": map[string]interface{}{"name": "terraform-name", "flag": true}}
	if missing := missingSubmittedResponseKeys(operations, response); len(missing) != 0 {
		t.Fatalf("name alias must satisfy readiness: %#v", missing)
	}
}

func TestFetchVerifyAndCacheRetriesUntilSubmittedKeyAppears(t *testing.T) {
	oldRetries, oldBackoff := PostOperationVerificationRetries, PostOperationVerificationBackoff
	PostOperationVerificationRetries, PostOperationVerificationBackoff = 2, 0
	t.Cleanup(func() { PostOperationVerificationRetries, PostOperationVerificationBackoff = oldRetries, oldBackoff })

	calls := 0
	config := ResourceConfig{ResourceType: "switchpoint", GetFunc: func(_ *openapi.APIClient, _ context.Context) (*http.Response, error) {
		calls++
		body := `{"switchpoint":{"one":{"name":"one"}}}`
		if calls == 2 {
			body = `{"switchpoint":{"one":{"name":"one","flag":true}}}`
		}
		return &http.Response{Body: io.NopCloser(strings.NewReader(body))}, nil
	}}
	flag := true
	res := NewResourceOperations()
	err := (&Manager{}).fetchVerifyAndCacheResourceResponse(context.Background(), context.Background(), config, res, nil, "PUT", map[string]interface{}{"one": readinessOperation{Flag: &flag}})
	if err != nil {
		t.Fatal(err)
	}
	if calls != 2 {
		t.Fatalf("GET calls = %d, want 2", calls)
	}
	if got := res.Responses["one"]["flag"]; got != true {
		t.Fatalf("cached flag = %#v", got)
	}
}

func TestFetchVerifyAndCacheExhaustsRetriesAndCachesFinalResponse(t *testing.T) {
	oldRetries, oldBackoff := PostOperationVerificationRetries, PostOperationVerificationBackoff
	PostOperationVerificationRetries, PostOperationVerificationBackoff = 2, 0
	t.Cleanup(func() { PostOperationVerificationRetries, PostOperationVerificationBackoff = oldRetries, oldBackoff })

	calls := 0
	config := ResourceConfig{ResourceType: "switchpoint", GetFunc: func(_ *openapi.APIClient, _ context.Context) (*http.Response, error) {
		calls++
		return &http.Response{Body: io.NopCloser(strings.NewReader(`{"switchpoint":{"one":{"name":"one"}}}`))}, nil
	}}
	flag := true
	res := NewResourceOperations()
	err := (&Manager{}).fetchVerifyAndCacheResourceResponse(context.Background(), context.Background(), config, res, nil, "PATCH", map[string]interface{}{"one": readinessOperation{Flag: &flag}})
	if err != nil {
		t.Fatal(err)
	}
	if calls != 3 {
		t.Fatalf("GET calls = %d, want 3", calls)
	}
	if _, exists := res.Responses["one"]["flag"]; exists {
		t.Fatal("incomplete final response was not cached")
	}
}

func TestFetchVerifyAndCacheDeleteDoesNotRetry(t *testing.T) {
	oldRetries, oldBackoff := PostOperationVerificationRetries, PostOperationVerificationBackoff
	PostOperationVerificationRetries, PostOperationVerificationBackoff = 2, 0
	t.Cleanup(func() { PostOperationVerificationRetries, PostOperationVerificationBackoff = oldRetries, oldBackoff })

	calls := 0
	config := ResourceConfig{ResourceType: "switchpoint", GetFunc: func(_ *openapi.APIClient, _ context.Context) (*http.Response, error) {
		calls++
		return &http.Response{Body: io.NopCloser(strings.NewReader(`{"switchpoint":{}}`))}, nil
	}}
	if err := (&Manager{}).fetchVerifyAndCacheResourceResponse(context.Background(), context.Background(), config, NewResourceOperations(), nil, "DELETE", map[string]interface{}{"one": readinessOperation{}}); err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("DELETE GET calls = %d, want 1", calls)
	}
}
