package bulkops_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"terraform-provider-verity/internal/bulkops"
)

type deleteCall struct {
	path  string
	names []string
}

func deleteRejectingServer(t *testing.T, param string, reply func(path string, names []string) string) (*httptest.Server, func() []deleteCall) {
	t.Helper()
	var mu sync.Mutex
	var calls []deleteCall
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{}`))
			return
		}
		names := append([]string(nil), r.URL.Query()[param]...)
		sort.Strings(names)
		mu.Lock()
		calls = append(calls, deleteCall{path: r.URL.Path, names: names})
		mu.Unlock()
		if message := reply(r.URL.Path, names); message != "" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"status": "failed", "message": %q}`, message)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)
	return server, func() []deleteCall {
		mu.Lock()
		defer mu.Unlock()
		return append([]deleteCall(nil), calls...)
	}
}

func contains(names []string, name string) bool {
	for _, candidate := range names {
		if candidate == name {
			return true
		}
	}
	return false
}

func waitFor(t *testing.T, mgr *bulkops.Manager, id string) error {
	t.Helper()
	return mgr.WaitForOperation(context.Background(), id, time.Second)
}

func TestPartialDeletePredefinedObjectIsSettledWithWarning(t *testing.T) {
	t.Parallel()
	server, calls := deleteRejectingServer(t, "packet_queue_name", func(_ string, names []string) string {
		if contains(names, "(Packet Queue)") {
			return "Object(s) could not be deleted as they are predefined objects:</br>packet_queue|(Packet Queue)|"
		}
		return ""
	})
	mgr := bulkops.GetManager(newTestClient(server.URL), nopClearCache, nil, "datacenter")
	ctx := context.Background()
	predefined := mgr.AddDelete(ctx, "packet_queue", "(Packet Queue)")
	ordinary := mgr.AddDelete(ctx, "packet_queue", "Queue A")

	if diags, _ := mgr.ExecuteDatacenterOperations(ctx); diags.HasError() {
		t.Fatalf("a predefined object alone must not fail the delete: %v", diags)
	}
	if err := waitFor(t, mgr, predefined); err != nil {
		t.Fatalf("predefined object: %v", err)
	}
	if warning := mgr.OperationWarning(predefined); !strings.Contains(warning, "predefined") || !strings.Contains(warning, "Terraform state only") {
		t.Fatalf("predefined object warning = %q", warning)
	}
	if err := waitFor(t, mgr, ordinary); err != nil {
		t.Fatalf("ordinary object: %v", err)
	}
	if mgr.OperationWarning(ordinary) != "" {
		t.Error("the ordinary object must not carry a warning")
	}
	got := calls()
	if len(got) != 2 || !contains(got[1].names, "Queue A") || contains(got[1].names, "(Packet Queue)") {
		t.Fatalf("calls = %+v, want the original request and one retry without the predefined object", got)
	}
}

func TestPartialDeleteReferencedObjectsFailAndTheRestIsDeleted(t *testing.T) {
	t.Parallel()
	server, calls := deleteRejectingServer(t, "threshold_name", func(path string, names []string) string {
		if path == "/thresholds" && contains(names, "Threshold A") {
			return "Object(s) could not be deleted as they are referenced by other objects:</br>" +
				"threshold|Threshold A| could not be deleted as it was referenced by threshold_group|Group A|</br>" +
				"threshold|Threshold B| could not be deleted as it was referenced by threshold_group|Group B|</br>"
		}
		return ""
	})
	mgr := bulkops.GetManager(newTestClient(server.URL), nopClearCache, nil, "datacenter")
	ctx := context.Background()
	blockedA := mgr.AddDelete(ctx, "threshold", "Threshold A")
	blockedB := mgr.AddDelete(ctx, "threshold", "Threshold B")
	free := mgr.AddDelete(ctx, "threshold", "Threshold C")
	laterType := mgr.AddDelete(ctx, "tenant", "Tenant A")

	diags, _ := mgr.ExecuteDatacenterOperations(ctx)
	if !diags.HasError() {
		t.Fatal("blocked objects must be reported as an error")
	}
	for id, group := range map[string]string{blockedA: "Group A", blockedB: "Group B"} {
		err := waitFor(t, mgr, id)
		if err == nil || !strings.Contains(err.Error(), "referenced by threshold_group|"+group+"|") {
			t.Errorf("blocked object error = %v, want the API reason naming %s", err, group)
		}
		if strings.Contains(fmt.Sprint(err), "aborted") {
			t.Errorf("a settled object was overwritten by the cancellation: %v", err)
		}
	}
	if err := waitFor(t, mgr, free); err != nil {
		t.Errorf("the unblocked object must be deleted by the retry: %v", err)
	}
	if err := waitFor(t, mgr, laterType); err == nil {
		t.Error("the later DELETE type should have been cancelled after the failure")
	}
	got := calls()
	if len(got) != 2 {
		t.Fatalf("calls = %+v, want the original request and one retry", got)
	}
	if len(got[1].names) != 1 || got[1].names[0] != "Threshold C" {
		t.Fatalf("retry = %+v, want only Threshold C", got[1])
	}
}

func TestPartialDeleteUnparseableReplyFailsTheWholeRequest(t *testing.T) {
	t.Parallel()
	server, calls := deleteRejectingServer(t, "threshold_name", func(string, []string) string {
		return "Object(s) could not be deleted as they are locked by an upgrade:</br>threshold|Threshold A|"
	})
	mgr := bulkops.GetManager(newTestClient(server.URL), nopClearCache, nil, "datacenter")
	ctx := context.Background()
	a := mgr.AddDelete(ctx, "threshold", "Threshold A")
	b := mgr.AddDelete(ctx, "threshold", "Threshold B")

	if diags, _ := mgr.ExecuteDatacenterOperations(ctx); !diags.HasError() {
		t.Fatal("expected an error")
	}
	for _, id := range []string{a, b} {
		if err := waitFor(t, mgr, id); err == nil {
			t.Errorf("operation %s: expected the whole request to fail", id)
		}
		if mgr.OperationWarning(id) != "" {
			t.Errorf("operation %s: unexpected warning", id)
		}
	}
	if got := calls(); len(got) != 1 {
		t.Fatalf("calls = %+v, want no retry for a reply that wasn't understood", got)
	}
}

func TestPartialDeleteInBatchedDeleteKeepsLaterBatchesRunning(t *testing.T) {
	t.Parallel()
	server, calls := deleteRejectingServer(t, "threshold_name", func(_ string, names []string) string {
		if contains(names, "T007") {
			return "Object(s) could not be deleted as they are referenced by other objects:</br>" +
				"threshold|T007| could not be deleted as it was referenced by threshold_group|Group A|"
		}
		return ""
	})
	mgr := bulkops.GetManager(newTestClient(server.URL), nopClearCache, nil, "datacenter")
	ctx := context.Background()
	ids := map[string]string{}
	for i := 0; i < 150; i++ {
		name := fmt.Sprintf("T%03d", i)
		ids[name] = mgr.AddDelete(ctx, "threshold", name)
	}

	if diags, _ := mgr.ExecuteDatacenterOperations(ctx); !diags.HasError() {
		t.Fatal("the blocked object must be reported as an error")
	}
	for name, id := range ids {
		err := waitFor(t, mgr, id)
		if name == "T007" {
			if err == nil || !strings.Contains(err.Error(), "Group A") {
				t.Errorf("T007 error = %v", err)
			}
			continue
		}
		if err != nil {
			t.Errorf("%s: %v", name, err)
		}
	}
	total := 0
	for _, call := range calls() {
		if !contains(call.names, "T007") {
			total += len(call.names)
		}
	}
	if total != 149 {
		t.Fatalf("successfully sent %d names, want all 149 unblocked ones across sub-batches and the retry", total)
	}
}
