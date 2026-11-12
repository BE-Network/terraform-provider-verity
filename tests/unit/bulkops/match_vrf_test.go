package bulkops_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"terraform-provider-verity/internal/bulkops"
)

type referenceAPI struct {
	mu        sync.Mutex
	clauses   map[string]string
	tenants   map[string]bool
	failGet   bool
	requests  []string
	patchBody string
}

func newReferenceAPI(t *testing.T, clauses map[string]string, tenants ...string) (*referenceAPI, *httptest.Server) {
	t.Helper()
	api := &referenceAPI{clauses: clauses, tenants: map[string]bool{}}
	for _, tenant := range tenants {
		api.tenants[tenant] = true
	}
	server := httptest.NewServer(http.HandlerFunc(api.handle))
	t.Cleanup(server.Close)
	return api, server
}

func (a *referenceAPI) handle(w http.ResponseWriter, r *http.Request) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.requests = append(a.requests, r.Method+" "+r.URL.Path)
	w.Header().Set("Content-Type", "application/json")
	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/routemapclauses":
		if a.failGet {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"unavailable"}`))
			return
		}
		objects := map[string]map[string]string{}
		for name, tenant := range a.clauses {
			refType := ""
			if tenant != "" {
				refType = "tenant"
			}
			objects[name] = map[string]string{"match_vrf": tenant, "match_vrf_ref_type_": refType}
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"route_map_clause": objects})
		return
	case r.Method == http.MethodPatch && r.URL.Path == "/routemapclauses":
		raw, _ := io.ReadAll(r.Body)
		a.patchBody = string(raw)
		var body map[string]map[string]map[string]string
		json.Unmarshal(raw, &body)
		for name, fields := range body["route_map_clause"] {
			if value, present := fields["match_vrf"]; present {
				a.clauses[name] = value
			}
		}
	case r.Method == http.MethodDelete && r.URL.Path == "/tenants":
		var blocked []string
		for _, tenant := range r.URL.Query()["tenant_name"] {
			for clause, referenced := range a.clauses {
				if referenced == tenant {
					blocked = append(blocked, fmt.Sprintf("tenant|%s| could not be deleted as it was referenced by route_map_clause|%s|", tenant, clause))
				}
			}
		}
		if len(blocked) > 0 {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, `{"status": "failed", "message": %q}`, "Object(s) could not be deleted as they are referenced by other objects:</br>"+strings.Join(blocked, "</br>"))
			return
		}
		for _, tenant := range r.URL.Query()["tenant_name"] {
			delete(a.tenants, tenant)
		}
	case r.Method == http.MethodDelete && r.URL.Path == "/routemapclauses":
		for _, clause := range r.URL.Query()["route_map_clause_name"] {
			delete(a.clauses, clause)
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}

func (a *referenceAPI) snapshot() ([]string, string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	return append([]string(nil), a.requests...), a.patchBody
}

func indexOf(requests []string, request string) int {
	for i, candidate := range requests {
		if candidate == request {
			return i
		}
	}
	return -1
}

func TestMatchVrfIsClearedBeforeDeletingClauseAndTenantTogether(t *testing.T) {
	t.Parallel()
	api, server := newReferenceAPI(t, map[string]string{"vrf_clause": "tenant_2", "other_clause": ""}, "tenant_2")
	mgr := bulkops.GetManager(newTestClient(server.URL), nopClearCache, nil, "datacenter")
	ctx := context.Background()
	tenant := mgr.AddDelete(ctx, "tenant", "tenant_2")
	vrfClause := mgr.AddDelete(ctx, "route_map_clause", "vrf_clause")
	otherClause := mgr.AddDelete(ctx, "route_map_clause", "other_clause")

	if diags, _ := mgr.ExecuteDatacenterOperations(ctx); diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	for _, id := range []string{tenant, vrfClause, otherClause} {
		if err := settled(t, mgr, id); err != nil {
			t.Fatal(err)
		}
	}
	requests, patch := api.snapshot()
	patchAt, tenantAt := indexOf(requests, "PATCH /routemapclauses"), indexOf(requests, "DELETE /tenants")
	if patchAt < 0 || tenantAt < 0 || patchAt > tenantAt {
		t.Fatalf("requests = %v, want the clause PATCH before the tenant DELETE", requests)
	}
	if got, want := canonicalJSON(t, patch), `{"route_map_clause":{"vrf_clause":{"match_vrf":"","match_vrf_ref_type_":""}}}`; got != want {
		t.Fatalf("PATCH body = %s, want %s", got, want)
	}
}

func TestMatchVrfOfAKeptClauseIsNotChanged(t *testing.T) {
	t.Parallel()
	api, server := newReferenceAPI(t, map[string]string{"kept_clause": "tenant_2", "other_clause": ""}, "tenant_2")
	mgr := bulkops.GetManager(newTestClient(server.URL), nopClearCache, nil, "datacenter")
	ctx := context.Background()
	tenant := mgr.AddDelete(ctx, "tenant", "tenant_2")
	mgr.AddDelete(ctx, "route_map_clause", "other_clause")

	if diags, _ := mgr.ExecuteDatacenterOperations(ctx); !diags.HasError() {
		t.Fatal("the tenant delete must fail while a kept clause references it")
	}
	if err := settled(t, mgr, tenant); err == nil || !strings.Contains(err.Error(), "kept_clause") {
		t.Fatalf("tenant error = %v, want the API reason naming kept_clause", err)
	}
	if requests, _ := api.snapshot(); indexOf(requests, "PATCH /routemapclauses") >= 0 {
		t.Fatalf("requests = %v: a clause that stays must never be changed", requests)
	}
}

func TestMatchVrfIsNotReadWithoutTenantDeletes(t *testing.T) {
	t.Parallel()
	api, server := newReferenceAPI(t, map[string]string{"vrf_clause": "tenant_2"}, "tenant_2")
	mgr := bulkops.GetManager(newTestClient(server.URL), nopClearCache, nil, "datacenter")
	ctx := context.Background()
	mgr.AddDelete(ctx, "route_map_clause", "vrf_clause")

	if diags, _ := mgr.ExecuteDatacenterOperations(ctx); diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	requests, _ := api.snapshot()
	getAt, deleteAt := indexOf(requests, "GET /routemapclauses"), indexOf(requests, "DELETE /routemapclauses")
	if (getAt >= 0 && getAt < deleteAt) || indexOf(requests, "PATCH /routemapclauses") >= 0 {
		t.Fatalf("requests = %v, want no clause read before the DELETE and no PATCH without tenant deletes", requests)
	}
}

func TestMatchVrfReadFailureKeepsTodaysBehaviour(t *testing.T) {
	t.Parallel()
	api, server := newReferenceAPI(t, map[string]string{"vrf_clause": "tenant_2"}, "tenant_2")
	api.failGet = true
	mgr := bulkops.GetManager(newTestClient(server.URL), nopClearCache, nil, "datacenter")
	ctx := context.Background()
	tenant := mgr.AddDelete(ctx, "tenant", "tenant_2")
	mgr.AddDelete(ctx, "route_map_clause", "vrf_clause")

	mgr.ExecuteDatacenterOperations(ctx)
	if err := settled(t, mgr, tenant); err == nil || !strings.Contains(err.Error(), "vrf_clause") {
		t.Fatalf("tenant error = %v, want the API reason", err)
	}
	if requests, _ := api.snapshot(); indexOf(requests, "PATCH /routemapclauses") >= 0 {
		t.Fatalf("requests = %v, want no PATCH when the clauses can't be read", requests)
	}
}

func canonicalJSON(t *testing.T, raw string) string {
	t.Helper()
	var value interface{}
	if err := json.Unmarshal([]byte(raw), &value); err != nil {
		t.Fatalf("%s: %v", raw, err)
	}
	out, _ := json.Marshal(value)
	return string(out)
}
