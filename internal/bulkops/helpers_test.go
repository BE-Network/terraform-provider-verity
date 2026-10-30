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

func TestMissingSubmittedResponseKeysRequiresCollectionKey(t *testing.T) {
	flag := true
	operations := map[string]interface{}{"terraform-name": readinessOperation{Flag: &flag}}
	response := map[string]interface{}{"api-key": map[string]interface{}{"name": "terraform-name", "flag": true}}
	if missing := missingSubmittedResponseKeys(operations, response); len(missing["terraform-name"]) != 1 {
		t.Fatalf("an inner name must not satisfy a missing collection key: %#v", missing)
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

func TestVerificationWithoutInnerName(t *testing.T) {
	for _, tc := range []struct {
		name     string
		response map[string]interface{}
		missing  bool
	}{
		{name: "key satisfies identity", response: map[string]interface{}{"one": map[string]interface{}{"flag": true}}},
		{name: "missing object", response: map[string]interface{}{}, missing: true},
		{name: "other fields still required", response: map[string]interface{}{"one": map[string]interface{}{}}, missing: true},
		{name: "body names do not override keys", response: map[string]interface{}{"one": map[string]interface{}{"flag": true}, "other": map[string]interface{}{"name": "one"}}},
		{name: "inner name is not a collection key", response: map[string]interface{}{"other": map[string]interface{}{"name": "one", "flag": true}}, missing: true},
		{name: "invalid object", response: map[string]interface{}{"one": 42}, missing: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			missing := missingSubmittedResponseKeys(map[string]interface{}{"one": map[string]interface{}{"name": "one", "flag": true}}, tc.response)
			if (len(missing) != 0) != tc.missing {
				t.Fatalf("missing=%v", missing)
			}
		})
	}
}

func TestCacheUsesOnlyCollectionKeys(t *testing.T) {
	res := NewResourceOperations()
	collection := map[string]interface{}{
		"one":   map[string]interface{}{"flag": true},
		"other": map[string]interface{}{"name": "one", "flag": false},
		"third": map[string]interface{}{"name": "not-a-key"},
	}
	(&Manager{}).cacheResourceResponse(context.Background(), ResourceConfig{}, res, nil, collection)
	if len(res.Responses) != len(collection) || res.Responses["one"]["flag"] != true || res.Responses["other"]["flag"] != false {
		t.Fatalf("cache must use exact collection keys: %v", res.Responses)
	}
	if _, found := res.Responses["not-a-key"]; found {
		t.Fatal("inner name became a cache alias")
	}
}

func TestVerificationDoesNotRetryForMissingInnerName(t *testing.T) {
	oldRetries, oldBackoff := PostOperationVerificationRetries, PostOperationVerificationBackoff
	PostOperationVerificationRetries, PostOperationVerificationBackoff = 2, 0
	defer func() { PostOperationVerificationRetries, PostOperationVerificationBackoff = oldRetries, oldBackoff }()
	calls := 0
	config := ResourceConfig{ResourceType: "gateway", GetFunc: func(_ *openapi.APIClient, _ context.Context) (*http.Response, error) {
		calls++
		return &http.Response{Body: io.NopCloser(strings.NewReader(`{"gateway":{"one":{"flag":true}}}`))}, nil
	}}
	res := NewResourceOperations()
	err := (&Manager{}).fetchVerifyAndCacheResourceResponse(context.Background(), context.Background(), config, res, nil, "PUT", map[string]interface{}{"one": map[string]interface{}{"name": "one", "flag": true}})
	if err != nil || calls != 1 || res.Responses["one"]["flag"] != true {
		t.Fatalf("calls=%d err=%v responses=%v", calls, err, res.Responses)
	}
}
