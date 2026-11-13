package bulkops_test

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"terraform-provider-verity/internal/bulkops"
)

func settled(t *testing.T, mgr *bulkops.Manager, id string) error {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := mgr.WaitForOperation(ctx, id, time.Second)
	if errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("operation %s never completed", id)
	}
	return err
}

func TestMixedPassFailureSettlesEveryWaiter(t *testing.T) {
	type outcome int
	const (
		succeeded outcome = iota
		failed
		cancelled
	)
	for _, tc := range []struct {
		name   string
		method string
		path   string
		want   map[string]outcome
	}{
		{
			name: "PUT failure stops the pass", method: http.MethodPut, path: "/tenants",
			want: map[string]outcome{"put tenant": failed, "put service": cancelled, "patch tenant": cancelled, "patch service": cancelled, "delete service": cancelled, "delete tenant": cancelled},
		},
		{
			name: "PATCH failure stops the pass", method: http.MethodPatch, path: "/tenants",
			want: map[string]outcome{"put tenant": succeeded, "put service": succeeded, "patch tenant": failed, "patch service": cancelled, "delete service": cancelled, "delete tenant": cancelled},
		},
		{
			name: "DELETE failure continues with later DELETE types", method: http.MethodDelete, path: "/services",
			want: map[string]outcome{"put tenant": succeeded, "put service": succeeded, "patch tenant": succeeded, "patch service": succeeded, "delete service": failed, "delete tenant": succeeded},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			server, _, _ := failingOrderTrackingServer(t, func(r *http.Request) bool {
				return r.Method == tc.method && r.URL.Path == tc.path
			})
			mgr := bulkops.GetManager(newTestClient(server.URL), nopClearCache, nil, "datacenter")
			ctx := context.Background()
			ids := map[string]string{
				"put tenant":     mgr.AddPut(ctx, "tenant", "put_tenant", zeroPutValue("tenant")),
				"put service":    mgr.AddPut(ctx, "service", "put_service", zeroPutValue("service")),
				"patch tenant":   mgr.AddPatch(ctx, "tenant", "patch_tenant", zeroPatchValue("tenant")),
				"patch service":  mgr.AddPatch(ctx, "service", "patch_service", zeroPatchValue("service")),
				"delete service": mgr.AddDelete(ctx, "service", "delete_service"),
				"delete tenant":  mgr.AddDelete(ctx, "tenant", "delete_tenant"),
			}

			if diags, _ := mgr.ExecuteDatacenterOperations(ctx); !diags.HasError() {
				t.Fatal("expected the failing request to be reported")
			}
			for label, id := range ids {
				err := settled(t, mgr, id)
				switch tc.want[label] {
				case succeeded:
					if err != nil {
						t.Errorf("%s: %v, want success", label, err)
					}
				case failed:
					if err == nil || strings.HasPrefix(err.Error(), "not sent") {
						t.Errorf("%s: %v, want the request's own failure", label, err)
					}
				case cancelled:
					if err == nil || !strings.HasPrefix(err.Error(), "not sent because the ") {
						t.Errorf("%s: %v, want a cancellation", label, err)
					}
				}
			}
		})
	}
}

func TestCancellationMessageNamesTheFailedRequestAndReason(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodPut && r.URL.Path == "/bundles" {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"Error":"Failed to decrypt some encrypted fields."}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)
	mgr := bulkops.GetManager(newTestClient(server.URL), nopClearCache, nil, "datacenter")
	ctx := context.Background()
	mgr.AddPut(ctx, "bundle", "b1", zeroPutValue("bundle"))
	switchpoint := mgr.AddPut(ctx, "switchpoint", "s1", zeroPutValue("switchpoint"))

	mgr.ExecuteDatacenterOperations(ctx)
	err := settled(t, mgr, switchpoint)
	if err == nil || !strings.Contains(err.Error(), "not sent because the bundle PUT request failed") ||
		!strings.Contains(err.Error(), "Failed to decrypt some encrypted fields.") {
		t.Fatalf("cancellation error = %v", err)
	}
}

func TestBatchedDeleteStopSettlesUnsentSubBatches(t *testing.T) {
	t.Parallel()
	var mu sync.Mutex
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodDelete {
			mu.Lock()
			calls++
			mu.Unlock()
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"simulated failure"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)
	mgr := bulkops.GetManager(newTestClient(server.URL), nopClearCache, nil, "datacenter")
	ctx := context.Background()
	var ids []string
	for i := 0; i < 250; i++ {
		ids = append(ids, mgr.AddDelete(ctx, "badge", fmt.Sprintf("badge_%03d", i)))
	}

	mgr.ExecuteDatacenterOperations(ctx)
	for _, id := range ids {
		if err := settled(t, mgr, id); err == nil {
			t.Fatalf("operation %s: expected a failure", id)
		}
	}
	mu.Lock()
	defer mu.Unlock()
	if calls == 0 || calls > 3 {
		t.Fatalf("got %d DELETE calls", calls)
	}
}

func TestArrivalDuringPassFollowsTypeTurns(t *testing.T) {
	t.Parallel()
	var mgr *bulkops.Manager
	var mu sync.Mutex
	var paths []string
	var lateService, lateList string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			mu.Lock()
			paths = append(paths, r.Method+" "+r.URL.Path)
			first := r.Method == http.MethodPut && r.URL.Path == "/tenants" && lateService == ""
			mu.Unlock()
			if first {
				lateService = mgr.AddPut(context.Background(), "service", "late_service", zeroPutValue("service"))
				lateList = mgr.AddPut(context.Background(), "community_list", "late_list", zeroPutValue("community_list"))
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)
	mgr = bulkops.GetManager(newTestClient(server.URL), nopClearCache, nil, "datacenter")
	ctx := context.Background()
	mgr.AddPut(ctx, "tenant", "t1", zeroPutValue("tenant"))

	if diags, _ := mgr.ExecuteDatacenterOperations(ctx); diags.HasError() {
		t.Fatal(diags)
	}
	if err := settled(t, mgr, lateService); err != nil {
		t.Fatalf("late service: %v", err)
	}
	mu.Lock()
	firstPass := append([]string(nil), paths...)
	mu.Unlock()
	if !contains(firstPass, "PUT /services") {
		t.Errorf("first pass %v: the service PUT arrived before its turn and should have run", firstPass)
	}
	if contains(firstPass, "PUT /communitylists") {
		t.Errorf("first pass %v: community_list had already had its turn and should wait", firstPass)
	}

	if diags, _ := mgr.ExecuteDatacenterOperations(ctx); diags.HasError() {
		t.Fatal(diags)
	}
	if err := settled(t, mgr, lateList); err != nil {
		t.Fatalf("late community list: %v", err)
	}
}
