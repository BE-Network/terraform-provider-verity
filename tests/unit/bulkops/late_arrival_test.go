package bulkops_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"terraform-provider-verity/internal/bulkops"
)

func okServer(t *testing.T) *httptest.Server {
	t.Helper()
	server, _, _ := orderTrackingServer(t)
	return server
}

func TestLateArrivalIsCounted(t *testing.T) {
	t.Parallel()
	mgr := bulkops.GetManager(newTestClient(okServer(t).URL), nopClearCache, nil, "datacenter")
	ctx := context.Background()

	mgr.AddPut(ctx, "bundle", "b1", zeroPutValue("bundle"))
	mgr.ExecuteDatacenterOperations(ctx)
	mgr.AddPut(ctx, "bundle", "b2", zeroPutValue("bundle"))
	mgr.ExecuteDatacenterOperations(ctx)

	if got := mgr.LateFlushes(); got["bundle PUT"] != 1 || len(got) != 1 {
		t.Fatalf("late flushes = %v, want one late bundle PUT", got)
	}
}

func TestDifferentMethodsAreNotLateArrivals(t *testing.T) {
	t.Parallel()
	mgr := bulkops.GetManager(newTestClient(okServer(t).URL), nopClearCache, nil, "datacenter")
	ctx := context.Background()

	mgr.AddDelete(ctx, "bundle", "b1")
	mgr.ExecuteDatacenterOperations(ctx)
	mgr.AddPut(ctx, "bundle", "b1", zeroPutValue("bundle"))
	mgr.ExecuteDatacenterOperations(ctx)

	if got := mgr.LateFlushes(); len(got) != 0 {
		t.Fatalf("late flushes = %v, want none for a replacement (DELETE then PUT)", got)
	}
}

func TestIntentionalSplitsAreNotLateArrivals(t *testing.T) {
	t.Parallel()
	t.Run("DELETE over 100 names", func(t *testing.T) {
		t.Parallel()
		mgr := bulkops.GetManager(newTestClient(okServer(t).URL), nopClearCache, nil, "datacenter")
		ctx := context.Background()
		for i := 0; i < 250; i++ {
			mgr.AddDelete(ctx, "badge", fmt.Sprintf("badge_%03d", i))
		}
		mgr.ExecuteDatacenterOperations(ctx)
		if got := mgr.LateFlushes(); len(got) != 0 {
			t.Fatalf("late flushes = %v, want none for DELETE sub-batches", got)
		}
	})
	t.Run("header split", func(t *testing.T) {
		t.Parallel()
		mgr := bulkops.GetManager(newTestClient(okServer(t).URL), nopClearCache, nil, "datacenter")
		ctx := context.Background()
		mgr.AddPut(ctx, "acl", "v4", zeroPutValue("acl"), map[string]string{"ip_version": "4"})
		mgr.AddPut(ctx, "acl", "v6", zeroPutValue("acl"), map[string]string{"ip_version": "6"})
		mgr.ExecuteDatacenterOperations(ctx)
		if got := mgr.LateFlushes(); len(got) != 0 {
			t.Fatalf("late flushes = %v, want none for the ACL v4/v6 split", got)
		}
	})
	t.Run("retry", func(t *testing.T) {
		t.Parallel()
		var attempts atomic.Int32
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			if r.Method == http.MethodPut && attempts.Add(1) == 1 {
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(`{"payload":"The system is currently being modified"}`))
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{}`))
		}))
		t.Cleanup(server.Close)
		mgr := bulkops.GetManager(newTestClient(server.URL), nopClearCache, nil, "datacenter")
		ctx := context.Background()
		id := mgr.AddPut(ctx, "bundle", "b1", zeroPutValue("bundle"))
		mgr.ExecuteDatacenterOperations(ctx)
		if err := settled(t, mgr, id); err != nil {
			t.Fatalf("the retried PUT failed: %v", err)
		}
		if attempts.Load() < 2 {
			t.Fatalf("the PUT was not retried (attempts=%d)", attempts.Load())
		}
		if got := mgr.LateFlushes(); len(got) != 0 {
			t.Fatalf("late flushes = %v, want none for a retry", got)
		}
	})
}
