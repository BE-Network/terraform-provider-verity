package bulkops_test

import (
	"context"
	"net/http"
	"sync/atomic"
	"testing"
	"time"

	"terraform-provider-verity/internal/bulkops"
)

func TestCancelledOperationsAreNeverSent(t *testing.T) {
	t.Parallel()
	var failTenant atomic.Bool
	failTenant.Store(true)
	server, records, mu := failingOrderTrackingServer(t, func(r *http.Request) bool {
		return failTenant.Load() && r.Method == http.MethodPut && r.URL.Path == "/tenants"
	})
	mgr := bulkops.GetManager(newTestClient(server.URL), nopClearCache, nil, "datacenter")
	ctx := context.Background()

	mgr.AddPut(ctx, "tenant", "test_tenant", zeroPutValue("tenant"))
	cancelled := []string{
		mgr.AddPut(ctx, "service", "cancelled_service", zeroPutValue("service")),
		mgr.AddPatch(ctx, "gateway", "cancelled_gateway", zeroPatchValue("gateway")),
		mgr.AddDelete(ctx, "bundle", "cancelled_bundle"),
	}

	if diags, _ := mgr.ExecuteDatacenterOperations(ctx); !diags.HasError() {
		t.Fatal("expected the tenant PUT to fail")
	}
	for _, id := range cancelled {
		if err := mgr.WaitForOperation(ctx, id, time.Second); err == nil {
			t.Errorf("operation %s: expected a cancellation error", id)
		}
	}

	firstRun := len(snapshotRecords(mu, records))
	failTenant.Store(false)

	fresh := mgr.AddPut(ctx, "ipv4_list", "fresh_list", zeroPutValue("ipv4_list"))
	if diags, _ := mgr.ExecuteDatacenterOperations(ctx); diags.HasError() {
		t.Fatalf("second flush failed: %v", diags)
	}
	if err := mgr.WaitForOperation(ctx, fresh, time.Second); err != nil {
		t.Fatalf("new operation after the cancellation failed: %v", err)
	}

	second := snapshotRecords(mu, records)[firstRun:]
	sawFresh := false
	for _, record := range second {
		switch record.Path {
		case "/ipv4lists":
			sawFresh = true
		case "/services", "/gateways", "/bundles":
			t.Errorf("cancelled operation was sent in a later flush: %s %s", record.Method, record.Path)
		}
	}
	if !sawFresh {
		t.Error("the operation added after the cancellation was not sent")
	}
}
