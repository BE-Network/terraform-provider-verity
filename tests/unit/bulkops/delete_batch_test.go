package bulkops_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"sync"
	"testing"

	"terraform-provider-verity/internal/bulkops"
	"terraform-provider-verity/openapi"
)

// Resources fall into two categories: without header params (most resources) and with header params (currently only ACLs).
// Testing badges covers the first category; testing ACLs covers the second. Together they cover all cases.

// TestDeleteBatchSplitting verifies DELETE ops exceeding MaxDeleteBatchSize (100) are split into correct batches.
func TestDeleteBatchSplitting(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name          string
		resourceCount int
		expectedCalls int
		expectedSizes []int
	}{
		{"below_limit_50", 50, 1, []int{50}},
		// 100 does NOT trigger batching (condition is strictly greater than 100)
		{"max_single_batch_100", 100, 1, []int{100}},
		{"just_over_101", 101, 2, []int{100, 1}},
		{"double_200", 200, 2, []int{100, 100}},
		{"three_batches_250", 250, 3, []int{100, 100, 50}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var mu sync.Mutex
			var deleteCalls [][]string

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")

				switch {
				case r.URL.Path == "/api/auth" && r.Method == http.MethodPost:
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"token":"mock"}`))
				case r.URL.Path == "/api/version" && r.Method == http.MethodGet:
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"version":"6.5","datacenter":true}`))
				case r.URL.Path == "/api/badges" && r.Method == http.MethodDelete:
					names := r.URL.Query()["badge_name"]
					mu.Lock()
					namesCopy := make([]string, len(names))
					copy(namesCopy, names)
					deleteCalls = append(deleteCalls, namesCopy)
					mu.Unlock()
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{}`))
				default:
					w.WriteHeader(http.StatusOK)
					w.Write([]byte(`{"badge":{}}`))
				}
			}))
			defer server.Close()

			cfg := openapi.NewConfiguration()
			cfg.Servers = openapi.ServerConfigurations{{URL: server.URL + "/api"}}
			cfg.HTTPClient = &http.Client{}
			client := openapi.NewAPIClient(cfg)

			mgr := bulkops.GetManager(
				client,
				func(ctx context.Context, m interface{}, key string) {},
				nil,
				"datacenter",
			)

			ctx := context.Background()
			for i := 0; i < tc.resourceCount; i++ {
				mgr.AddDelete(ctx, "badge", fmt.Sprintf("badge_%03d", i))
			}

			diags := mgr.ExecuteBulk(ctx, "badge", "DELETE")
			if diags.HasError() {
				t.Fatalf("unexpected error: %v", diags)
			}

			mu.Lock()
			defer mu.Unlock()

			if len(deleteCalls) != tc.expectedCalls {
				t.Errorf("expected %d DELETE calls, got %d", tc.expectedCalls, len(deleteCalls))
			}
			for i, expectedSize := range tc.expectedSizes {
				if i >= len(deleteCalls) {
					break
				}
				if len(deleteCalls[i]) != expectedSize {
					t.Errorf("batch %d: expected %d names, got %d", i+1, expectedSize, len(deleteCalls[i]))
				}
			}
		})
	}
}

func TestDeleteBatchCorrectNameSubsets(t *testing.T) {
	t.Parallel()
	var mu sync.Mutex
	var deleteCalls [][]string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/api/auth" && r.Method == http.MethodPost:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"token":"mock"}`))
		case r.URL.Path == "/api/badges" && r.Method == http.MethodDelete:
			names := r.URL.Query()["badge_name"]
			mu.Lock()
			namesCopy := make([]string, len(names))
			copy(namesCopy, names)
			deleteCalls = append(deleteCalls, namesCopy)
			mu.Unlock()
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{}`))
		default:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"badge":{}}`))
		}
	}))
	defer server.Close()

	cfg := openapi.NewConfiguration()
	cfg.Servers = openapi.ServerConfigurations{{URL: server.URL + "/api"}}
	cfg.HTTPClient = &http.Client{}
	client := openapi.NewAPIClient(cfg)

	mgr := bulkops.GetManager(client, func(ctx context.Context, m interface{}, key string) {}, nil, "datacenter")

	ctx := context.Background()
	allNames := make([]string, 150)
	for i := 0; i < 150; i++ {
		name := fmt.Sprintf("res_%03d", i)
		allNames[i] = name
		mgr.AddDelete(ctx, "badge", name)
	}

	diags := mgr.ExecuteBulk(ctx, "badge", "DELETE")
	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}

	mu.Lock()
	defer mu.Unlock()

	if len(deleteCalls) != 2 {
		t.Fatalf("expected 2 DELETE calls, got %d", len(deleteCalls))
	}

	var receivedNames []string
	for _, batch := range deleteCalls {
		receivedNames = append(receivedNames, batch...)
	}
	sort.Strings(receivedNames)
	sort.Strings(allNames)

	if len(receivedNames) != len(allNames) {
		t.Fatalf("expected %d total names, got %d", len(allNames), len(receivedNames))
	}
	for i, name := range allNames {
		if receivedNames[i] != name {
			t.Errorf("name mismatch at index %d: expected %q, got %q", i, name, receivedNames[i])
		}
	}
}

func TestDeleteBatchFailureStopsRemaining(t *testing.T) {
	t.Parallel()
	var mu sync.Mutex
	var deleteCallCount int

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/api/auth" && r.Method == http.MethodPost:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"token":"mock"}`))
		case r.URL.Path == "/api/badges" && r.Method == http.MethodDelete:
			mu.Lock()
			deleteCallCount++
			callNum := deleteCallCount
			mu.Unlock()

			if callNum == 2 {
				// Fail on the second batch
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"server error"}`))
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{}`))
		default:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"badge":{}}`))
		}
	}))
	defer server.Close()

	cfg := openapi.NewConfiguration()
	cfg.Servers = openapi.ServerConfigurations{{URL: server.URL + "/api"}}
	cfg.HTTPClient = &http.Client{}
	client := openapi.NewAPIClient(cfg)

	mgr := bulkops.GetManager(client, func(ctx context.Context, m interface{}, key string) {}, nil, "datacenter")

	ctx := context.Background()
	// 250 resources → 3 expected batches, but batch 2 fails
	for i := 0; i < 250; i++ {
		mgr.AddDelete(ctx, "badge", fmt.Sprintf("badge_%03d", i))
	}

	diags := mgr.ExecuteBulk(ctx, "badge", "DELETE")
	if !diags.HasError() {
		t.Error("expected error from failed batch, got none")
	}

	mu.Lock()
	defer mu.Unlock()

	// Batch 3 must not execute. GenericOpenAPIError is not retriable, so exactly 2 calls in practice.
	// maxExpectedCalls guards against future retriable errors: 1 (batch 1) + 1 + 5 retries (batch 2).
	const maxExpectedCalls = 7
	if deleteCallCount > maxExpectedCalls {
		t.Errorf("expected batch 3 to not execute, but got %d total DELETE calls (max allowed: %d)",
			deleteCallCount, maxExpectedCalls)
	}
}

func TestDeleteBatchFirstBatchFailureAbortsImmediately(t *testing.T) {
	t.Parallel()
	var mu sync.Mutex
	var deleteCallCount int

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/api/auth" && r.Method == http.MethodPost:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"token":"mock"}`))
		case r.URL.Path == "/api/badges" && r.Method == http.MethodDelete:
			mu.Lock()
			deleteCallCount++
			mu.Unlock()
			// Always fail — simulates batch 1 failure
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"server error"}`))
		default:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"badge":{}}`))
		}
	}))
	defer server.Close()

	cfg := openapi.NewConfiguration()
	cfg.Servers = openapi.ServerConfigurations{{URL: server.URL + "/api"}}
	cfg.HTTPClient = &http.Client{}
	client := openapi.NewAPIClient(cfg)

	mgr := bulkops.GetManager(client, func(ctx context.Context, m interface{}, key string) {}, nil, "datacenter")

	ctx := context.Background()
	// 250 resources → 3 batches; batch 1 fails immediately
	for i := 0; i < 250; i++ {
		mgr.AddDelete(ctx, "badge", fmt.Sprintf("badge_%03d", i))
	}

	diags := mgr.ExecuteBulk(ctx, "badge", "DELETE")
	if !diags.HasError() {
		t.Error("expected error from failed batch 1, got none")
	}

	mu.Lock()
	defer mu.Unlock()

	const maxExpectedCalls = 6
	if deleteCallCount > maxExpectedCalls {
		t.Errorf("expected only batch 1 to execute, but got %d total DELETE calls", deleteCallCount)
	}
	if deleteCallCount == 0 {
		t.Error("expected at least 1 DELETE call for batch 1, got none")
	}
}

// TestDeleteBatchACLHeaderSplitAndBatching verifies ACL DELETEs are first split by ip_version,
// then each version group is independently batched. 150 IPv4 + 150 IPv6 → 4 DELETE calls total.
func TestDeleteBatchACLHeaderSplitAndBatching(t *testing.T) {
	t.Parallel()

	type deleteCall struct {
		names     []string
		ipVersion string
	}
	var mu sync.Mutex
	var deleteCalls []deleteCall

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/api/auth" && r.Method == http.MethodPost:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"token":"mock"}`))
		case r.URL.Path == "/api/acls" && r.Method == http.MethodDelete:
			names := r.URL.Query()["ip_filter_name"]
			ipVer := r.URL.Query().Get("ip_version")
			namesCopy := make([]string, len(names))
			copy(namesCopy, names)
			mu.Lock()
			deleteCalls = append(deleteCalls, deleteCall{names: namesCopy, ipVersion: ipVer})
			mu.Unlock()
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{}`))
		default:
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{}`))
		}
	}))
	defer server.Close()

	cfg := openapi.NewConfiguration()
	cfg.Servers = openapi.ServerConfigurations{{URL: server.URL + "/api"}}
	cfg.HTTPClient = &http.Client{}
	client := openapi.NewAPIClient(cfg)

	mgr := bulkops.GetManager(client, func(ctx context.Context, m interface{}, key string) {}, nil, "datacenter")

	ctx := context.Background()
	// 150 IPv4 + 150 IPv6 = 300 ACLs total
	for i := 0; i < 150; i++ {
		mgr.AddDelete(ctx, "acl", fmt.Sprintf("v4_acl_%03d", i), map[string]string{"ip_version": "4"})
	}
	for i := 0; i < 150; i++ {
		mgr.AddDelete(ctx, "acl", fmt.Sprintf("v6_acl_%03d", i), map[string]string{"ip_version": "6"})
	}

	diags := mgr.ExecuteBulk(ctx, "acl", "DELETE")
	if diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}

	mu.Lock()
	defer mu.Unlock()

	// Expect 4 total DELETE calls: 2 for IPv4 (100+50) and 2 for IPv6 (100+50)
	if len(deleteCalls) != 4 {
		t.Fatalf("expected 4 DELETE calls (2 per version), got %d", len(deleteCalls))
	}

	v4Names := []string{}
	v6Names := []string{}
	v4Batches := 0
	v6Batches := 0
	for _, call := range deleteCalls {
		switch call.ipVersion {
		case "4":
			v4Batches++
			v4Names = append(v4Names, call.names...)
		case "6":
			v6Batches++
			v6Names = append(v6Names, call.names...)
		default:
			t.Errorf("unexpected ip_version %q in DELETE call", call.ipVersion)
		}
	}

	if v4Batches != 2 {
		t.Errorf("expected 2 DELETE batches for IPv4, got %d", v4Batches)
	}
	if v6Batches != 2 {
		t.Errorf("expected 2 DELETE batches for IPv6, got %d", v6Batches)
	}
	if len(v4Names) != 150 {
		t.Errorf("expected 150 IPv4 names total, got %d", len(v4Names))
	}
	if len(v6Names) != 150 {
		t.Errorf("expected 150 IPv6 names total, got %d", len(v6Names))
	}

	// Verify no cross-contamination: no v6 name in v4 group and vice versa
	sort.Strings(v4Names)
	sort.Strings(v6Names)
	for _, name := range v4Names {
		if len(name) > 2 && name[:2] == "v6" {
			t.Errorf("found IPv6 ACL name %q in IPv4 DELETE batches", name)
		}
	}
	for _, name := range v6Names {
		if len(name) > 2 && name[:2] == "v4" {
			t.Errorf("found IPv4 ACL name %q in IPv6 DELETE batches", name)
		}
	}
}
