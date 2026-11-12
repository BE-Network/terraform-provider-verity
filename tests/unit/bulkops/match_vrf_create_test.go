package bulkops_test

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"terraform-provider-verity/internal/bulkops"
	"terraform-provider-verity/openapi"
)

type clauseRef struct {
	MatchVrf string `json:"match_vrf"`
	RefType  string `json:"match_vrf_ref_type_"`
}

type gate struct {
	entered chan struct{}
	release chan struct{}
	once    sync.Once
}

func newGate() *gate {
	return &gate{entered: make(chan struct{}), release: make(chan struct{})}
}

type recordedRequest struct {
	key  string
	body string
}

type createAPI struct {
	mu           sync.Mutex
	clauses      map[string]clauseRef
	tenants      map[string]bool
	fail         map[string]string
	gates        map[string]*gate
	staleGets    int
	failGets     bool
	patched      bool
	requests     []recordedRequest
	onServicePut func()
}

const afterPatchGet = "GET /routemapclauses after PATCH"

func newCreateAPI(t *testing.T, tenants ...string) (*createAPI, *httptest.Server) {
	t.Helper()
	api := &createAPI{clauses: map[string]clauseRef{}, tenants: map[string]bool{}, fail: map[string]string{}, gates: map[string]*gate{}}
	for _, tenant := range tenants {
		api.tenants[tenant] = true
	}
	server := httptest.NewServer(http.HandlerFunc(api.handle))
	t.Cleanup(func() {
		api.mu.Lock()
		for _, g := range api.gates {
			g.once.Do(func() { close(g.release) })
		}
		api.mu.Unlock()
		server.Close()
	})
	return api, server
}

func (a *createAPI) block(key string) *gate {
	a.mu.Lock()
	defer a.mu.Unlock()
	g := newGate()
	a.gates[key] = g
	return g
}

func (g *gate) open() {
	g.once.Do(func() { close(g.release) })
}

func (g *gate) waitEntered(t *testing.T) {
	t.Helper()
	select {
	case <-g.entered:
	case <-time.After(5 * time.Second):
		t.Fatal("the blocked request was never sent")
	}
}

func (a *createAPI) handle(w http.ResponseWriter, r *http.Request) {
	raw, _ := io.ReadAll(r.Body)
	key := r.Method + " " + r.URL.Path

	a.mu.Lock()
	gateKey := key
	if key == "GET /routemapclauses" && a.patched {
		gateKey = afterPatchGet
	}
	g := a.gates[gateKey]
	if g != nil {
		delete(a.gates, gateKey)
	}
	onServicePut := a.onServicePut
	a.mu.Unlock()
	if g != nil {
		close(g.entered)
		<-g.release
	}
	if key == "PUT /services" && onServicePut != nil {
		onServicePut()
	}

	a.mu.Lock()
	defer a.mu.Unlock()
	a.requests = append(a.requests, recordedRequest{key: key, body: string(raw)})
	w.Header().Set("Content-Type", "application/json")
	if reason, failing := a.fail[key]; failing {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"status": "failed", "message": %q}`, reason)
		return
	}

	switch key {
	case "GET /routemapclauses":
		if a.patched && a.failGets {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"message": "read unavailable"}`))
			return
		}
		objects := map[string]clauseRef{}
		for name, ref := range a.clauses {
			if a.patched && a.staleGets > 0 {
				ref = clauseRef{}
			}
			objects[name] = ref
		}
		if a.patched && a.staleGets > 0 {
			a.staleGets--
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"route_map_clause": objects})
		return
	case "GET /tenants":
		objects := map[string]map[string]string{}
		for name := range a.tenants {
			objects[name] = map[string]string{}
		}
		json.NewEncoder(w).Encode(map[string]interface{}{"tenant": objects})
		return
	case "PUT /routemapclauses":
		var body map[string]map[string]clauseRef
		json.Unmarshal(raw, &body)
		for name, ref := range body["route_map_clause"] {
			if ref.MatchVrf != "" && !a.tenants[ref.MatchVrf] {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, `{"status": "failed", "message": %q}`, fmt.Sprintf("Cannot resolve object reference: tenant %q is missing from the database or import payload, but is referenced by route_map_clause %q", ref.MatchVrf, name))
				return
			}
		}
		for name, ref := range body["route_map_clause"] {
			a.clauses[name] = ref
		}
	case "PATCH /routemapclauses":
		var body map[string]map[string]clauseRef
		json.Unmarshal(raw, &body)
		for name, ref := range body["route_map_clause"] {
			a.clauses[name] = ref
		}
		a.patched = true
	case "PUT /tenants":
		var body map[string]map[string]json.RawMessage
		json.Unmarshal(raw, &body)
		for name := range body["tenant"] {
			a.tenants[name] = true
		}
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}

func (a *createAPI) log() []recordedRequest {
	a.mu.Lock()
	defer a.mu.Unlock()
	return append([]recordedRequest(nil), a.requests...)
}

func (a *createAPI) keys() []string {
	var keys []string
	for _, request := range a.log() {
		keys = append(keys, request.key)
	}
	return keys
}

func (a *createAPI) bodyOf(t *testing.T, key string) string {
	t.Helper()
	for _, request := range a.log() {
		if request.key == key {
			return canonicalJSON(t, request.body)
		}
	}
	t.Fatalf("no %s request in %v", key, a.keys())
	return ""
}

func clauseValue(tenant string) openapi.RoutemapclausesPutRequestRouteMapClauseValue {
	value := openapi.RoutemapclausesPutRequestRouteMapClauseValue{}
	value.SetMatchVrf(tenant)
	value.SetMatchVrfRefType("tenant")
	return value
}

func newDatacenterManager(server *httptest.Server) *bulkops.Manager {
	return bulkops.GetManager(newTestClient(server.URL), nopClearCache, nil, "datacenter")
}

func cachedMatchVrf(mgr *bulkops.Manager, name string) (string, bool) {
	data, exists := mgr.GetResourceResponse("route_map_clause", name)
	if !exists {
		return "", false
	}
	value, _ := data["match_vrf"].(string)
	return value, true
}

func requireOrder(t *testing.T, keys []string, ordered ...string) {
	t.Helper()
	last := -1
	for _, key := range ordered {
		at := indexOf(keys, key)
		if at < 0 || at < last {
			t.Fatalf("requests = %v, want %v in this order", keys, ordered)
		}
		last = at
	}
}

func TestClauseCreatedWithItsTenantGetsMatchVrfAfterTheTenant(t *testing.T) {
	t.Parallel()
	api, server := newCreateAPI(t)
	mgr := newDatacenterManager(server)
	ctx := context.Background()
	clause := mgr.AddPut(ctx, "route_map_clause", "vrf_clause", clauseValue("tenant_2"))
	tenant := mgr.AddPut(ctx, "tenant", "tenant_2", zeroPutValue("tenant"))

	if diags, _ := mgr.ExecuteDatacenterOperations(ctx); diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	for _, id := range []string{clause, tenant} {
		if err := settled(t, mgr, id); err != nil {
			t.Fatal(err)
		}
	}
	keys := api.keys()
	requireOrder(t, keys, "PUT /routemapclauses", "PUT /tenants", "PATCH /routemapclauses")
	if got := api.bodyOf(t, "PUT /routemapclauses"); !strings.Contains(got, `"match_vrf":""`) || !strings.Contains(got, `"match_vrf_ref_type_":""`) {
		t.Fatalf("clause PUT = %s, want match_vrf and its ref type sent empty", got)
	}
	if got, want := api.bodyOf(t, "PATCH /routemapclauses"), `{"route_map_clause":{"vrf_clause":{"match_vrf":"tenant_2","match_vrf_ref_type_":"tenant"}}}`; got != want {
		t.Fatalf("PATCH = %s, want %s", got, want)
	}
	if value, _ := cachedMatchVrf(mgr, "vrf_clause"); value != "tenant_2" {
		t.Fatalf("cached match_vrf = %q, want the restored value", value)
	}
}

func TestMatchVrfIsRestoredBeforeLaterPutTypes(t *testing.T) {
	t.Parallel()
	api, server := newCreateAPI(t)
	mgr := newDatacenterManager(server)
	ctx := context.Background()
	clause := mgr.AddPut(ctx, "route_map_clause", "vrf_clause", clauseValue("tenant_2"))
	mgr.AddPut(ctx, "tenant", "tenant_2", zeroPutValue("tenant"))
	service := mgr.AddPut(ctx, "service", "service_1", zeroPutValue("service"))
	settledBeforeService := make(chan bool, 1)
	api.mu.Lock()
	api.onServicePut = func() {
		waitCtx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
		defer cancel()
		settledBeforeService <- mgr.WaitForOperation(waitCtx, clause, time.Second) == nil
	}
	api.mu.Unlock()

	if diags, _ := mgr.ExecuteDatacenterOperations(ctx); diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	if err := settled(t, mgr, service); err != nil {
		t.Fatal(err)
	}
	requireOrder(t, api.keys(), "PUT /tenants", "PATCH /routemapclauses", "PUT /services")
	if !<-settledBeforeService {
		t.Fatal("the clause create must be settled before later PUT types are sent")
	}
}

func TestHeldClauseWaiterDoesNotTimeOut(t *testing.T) {
	for _, blocked := range []string{"PUT /tenants", "PATCH /routemapclauses", afterPatchGet} {
		t.Run(blocked, func(t *testing.T) {
			t.Parallel()
			api, server := newCreateAPI(t)
			mgr := newDatacenterManager(server)
			ctx := context.Background()
			clause := mgr.AddPut(ctx, "route_map_clause", "vrf_clause", clauseValue("tenant_2"))
			mgr.AddPut(ctx, "tenant", "tenant_2", zeroPutValue("tenant"))
			g := api.block(blocked)

			done := make(chan struct{})
			go func() {
				defer close(done)
				mgr.ExecuteDatacenterOperations(ctx)
			}()
			g.waitEntered(t)

			early, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
			if err := mgr.WaitForOperation(early, clause, time.Second); !errors.Is(err, context.DeadlineExceeded) {
				t.Fatalf("held clause settled early: %v", err)
			}
			cancel()

			waited := make(chan error, 1)
			go func() { waited <- mgr.WaitForOperation(context.Background(), clause, 50*time.Millisecond) }()
			time.Sleep(300 * time.Millisecond)
			g.open()
			if err := <-waited; err != nil {
				t.Fatalf("waiter returned %v, want success after the blocked %s", err, blocked)
			}
			<-done
		})
	}
}

func TestInFlightClausePutStillTimesOut(t *testing.T) {
	t.Parallel()
	api, server := newCreateAPI(t)
	mgr := newDatacenterManager(server)
	ctx := context.Background()
	clause := mgr.AddPut(ctx, "route_map_clause", "vrf_clause", clauseValue("tenant_2"))
	mgr.AddPut(ctx, "tenant", "tenant_2", zeroPutValue("tenant"))
	g := api.block("PUT /routemapclauses")

	done := make(chan struct{})
	go func() {
		defer close(done)
		mgr.ExecuteDatacenterOperations(ctx)
	}()
	g.waitEntered(t)
	err := mgr.WaitForOperation(context.Background(), clause, 50*time.Millisecond)
	g.open()
	<-done
	if err == nil || !strings.Contains(err.Error(), "timeout") {
		t.Fatalf("waiter returned %v, want the normal timeout while the clause PUT is in flight", err)
	}
}

func TestClauseQueuedDuringAnEarlierPutIsStillHeld(t *testing.T) {
	t.Parallel()
	api, server := newCreateAPI(t)
	mgr := newDatacenterManager(server)
	ctx := context.Background()
	mgr.AddPut(ctx, "community_list", "list_1", zeroPutValue("community_list"))
	g := api.block("PUT /communitylists")

	done := make(chan struct{})
	go func() {
		defer close(done)
		if diags, _ := mgr.ExecuteDatacenterOperations(ctx); diags.HasError() {
			t.Errorf("unexpected error: %v", diags)
		}
	}()
	g.waitEntered(t)
	clause := mgr.AddPut(ctx, "route_map_clause", "vrf_clause", clauseValue("tenant_2"))
	mgr.AddPut(ctx, "tenant", "tenant_2", zeroPutValue("tenant"))
	g.open()
	<-done

	if err := settled(t, mgr, clause); err != nil {
		t.Fatal(err)
	}
	if got := api.bodyOf(t, "PUT /routemapclauses"); !strings.Contains(got, `"match_vrf":""`) {
		t.Fatalf("clause PUT = %s, want match_vrf held back", got)
	}
	requireOrder(t, api.keys(), "PUT /routemapclauses", "PUT /tenants", "PATCH /routemapclauses")
}

func TestClauseReferencingAnExistingTenantIsSentUnchanged(t *testing.T) {
	t.Parallel()
	api, server := newCreateAPI(t, "tenant_1")
	mgr := newDatacenterManager(server)
	ctx := context.Background()
	kept := mgr.AddPut(ctx, "route_map_clause", "kept_clause", clauseValue("tenant_1"))
	held := mgr.AddPut(ctx, "route_map_clause", "held_clause", clauseValue("tenant_2"))
	mgr.AddPut(ctx, "tenant", "tenant_2", zeroPutValue("tenant"))

	if diags, _ := mgr.ExecuteDatacenterOperations(ctx); diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	for _, id := range []string{kept, held} {
		if err := settled(t, mgr, id); err != nil {
			t.Fatal(err)
		}
	}
	put := api.bodyOf(t, "PUT /routemapclauses")
	if !strings.Contains(put, `"kept_clause":{"match_vrf":"tenant_1","match_vrf_ref_type_":"tenant"}`) || !strings.Contains(put, `"held_clause":{"match_vrf":"","match_vrf_ref_type_":""}`) {
		t.Fatalf("clause PUT = %s, want only held_clause stripped", put)
	}
	if got, want := api.bodyOf(t, "PATCH /routemapclauses"), `{"route_map_clause":{"held_clause":{"match_vrf":"tenant_2","match_vrf_ref_type_":"tenant"}}}`; got != want {
		t.Fatalf("PATCH = %s, want %s", got, want)
	}
}

func TestClauseWithoutTenantCreatesIsNotPatched(t *testing.T) {
	t.Parallel()
	api, server := newCreateAPI(t, "tenant_1")
	mgr := newDatacenterManager(server)
	ctx := context.Background()
	clause := mgr.AddPut(ctx, "route_map_clause", "vrf_clause", clauseValue("tenant_1"))

	if diags, _ := mgr.ExecuteDatacenterOperations(ctx); diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	if err := settled(t, mgr, clause); err != nil {
		t.Fatal(err)
	}
	if got := api.bodyOf(t, "PUT /routemapclauses"); !strings.Contains(got, `"match_vrf":"tenant_1"`) {
		t.Fatalf("clause PUT = %s, want it unchanged", got)
	}
	if indexOf(api.keys(), "PATCH /routemapclauses") >= 0 {
		t.Fatalf("requests = %v, want no PATCH", api.keys())
	}
}

func TestExistingClauseWithoutMatchVrfIsPatchedWithOnlyTheClauseQueued(t *testing.T) {
	t.Parallel()
	api, server := newCreateAPI(t, "tenant_2")
	api.clauses["vrf_clause"] = clauseRef{}
	mgr := newDatacenterManager(server)
	ctx := context.Background()
	clause := mgr.AddPut(ctx, "route_map_clause", "vrf_clause", clauseValue("tenant_2"))

	if diags, _ := mgr.ExecuteDatacenterOperations(ctx); diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	if err := settled(t, mgr, clause); err != nil {
		t.Fatal(err)
	}
	keys := api.keys()
	if indexOf(keys, "PUT /routemapclauses") >= 0 || indexOf(keys, "PUT /tenants") >= 0 {
		t.Fatalf("requests = %v, want no clause or tenant PUT", keys)
	}
	if got, want := api.bodyOf(t, "PATCH /routemapclauses"), `{"route_map_clause":{"vrf_clause":{"match_vrf":"tenant_2","match_vrf_ref_type_":"tenant"}}}`; got != want {
		t.Fatalf("PATCH = %s, want %s", got, want)
	}
	if value, _ := cachedMatchVrf(mgr, "vrf_clause"); value != "tenant_2" {
		t.Fatalf("cached match_vrf = %q, want the restored value", value)
	}
}

func TestExistingClauseWithTheSameMatchVrfIsAdoptedAsToday(t *testing.T) {
	t.Parallel()
	api, server := newCreateAPI(t, "tenant_2")
	api.clauses["vrf_clause"] = clauseRef{MatchVrf: "tenant_2", RefType: "tenant"}
	mgr := newDatacenterManager(server)
	ctx := context.Background()
	clause := mgr.AddPut(ctx, "route_map_clause", "vrf_clause", clauseValue("tenant_2"))

	if diags, _ := mgr.ExecuteDatacenterOperations(ctx); diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	if err := settled(t, mgr, clause); err != nil {
		t.Fatal(err)
	}
	if keys := api.keys(); indexOf(keys, "PATCH /routemapclauses") >= 0 || indexOf(keys, "PUT /routemapclauses") >= 0 {
		t.Fatalf("requests = %v, want neither a PUT nor a PATCH", keys)
	}
}

func TestHeldClauseFailsWhenAPutBeforeTheRestoreFails(t *testing.T) {
	for _, failing := range []string{"PUT /tenants", "PUT /routemaps"} {
		t.Run(failing, func(t *testing.T) {
			t.Parallel()
			api, server := newCreateAPI(t)
			api.fail[failing] = "rejected by the test"
			mgr := newDatacenterManager(server)
			ctx := context.Background()
			clause := mgr.AddPut(ctx, "route_map_clause", "vrf_clause", clauseValue("tenant_2"))
			mgr.AddPut(ctx, "route_map", "map_1", zeroPutValue("route_map"))
			mgr.AddPut(ctx, "tenant", "tenant_2", zeroPutValue("tenant"))

			start := time.Now()
			if diags, _ := mgr.ExecuteDatacenterOperations(ctx); !diags.HasError() {
				t.Fatal("the failed PUT must be reported")
			}
			err := settled(t, mgr, clause)
			if err == nil || !strings.Contains(err.Error(), "rejected by the test") {
				t.Fatalf("clause error = %v, want the failed PUT's reason", err)
			}
			if elapsed := time.Since(start); elapsed > 3*time.Second {
				t.Fatalf("took %v, want the held clause to settle right away", elapsed)
			}
			if indexOf(api.keys(), "PATCH /routemapclauses") >= 0 {
				t.Fatalf("requests = %v, want no PATCH", api.keys())
			}
			if _, cached := cachedMatchVrf(mgr, "vrf_clause"); cached {
				t.Fatal("the stale clause response must not stay cached")
			}
		})
	}
}

func TestFailedRestoreStopsThePass(t *testing.T) {
	for _, tc := range []struct {
		name    string
		setup   func(*createAPI)
		reason  string
		patches int
	}{
		{"PATCH fails", func(a *createAPI) { a.fail["PATCH /routemapclauses"] = "patch rejected" }, "patch rejected", 1},
		{"GET fails", func(a *createAPI) { a.failGets = true }, "could not be confirmed", 1},
		{"GET keeps the old value", func(a *createAPI) { a.staleGets = 1000 }, "could not be confirmed", 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			api, server := newCreateAPI(t)
			tc.setup(api)
			mgr := newDatacenterManager(server)
			ctx := context.Background()
			clause := mgr.AddPut(ctx, "route_map_clause", "vrf_clause", clauseValue("tenant_2"))
			tenant := mgr.AddPut(ctx, "tenant", "tenant_2", zeroPutValue("tenant"))
			service := mgr.AddPut(ctx, "service", "service_1", zeroPutValue("service"))

			diags, _ := mgr.ExecuteDatacenterOperations(ctx)
			if !diags.HasError() || !strings.Contains(fmt.Sprint(diags), tc.reason) {
				t.Fatalf("diagnostics = %v, want the restore failure", diags)
			}
			if err := settled(t, mgr, tenant); err != nil {
				t.Fatalf("tenant error = %v, want it created", err)
			}
			if err := settled(t, mgr, clause); err == nil || !strings.Contains(err.Error(), tc.reason) {
				t.Fatalf("clause error = %v, want %q", err, tc.reason)
			}
			if err := settled(t, mgr, service); err == nil || !strings.Contains(err.Error(), "not sent") {
				t.Fatalf("service error = %v, want it cancelled", err)
			}
			if indexOf(api.keys(), "PUT /services") >= 0 {
				t.Fatalf("requests = %v, want nothing sent after the failed restore", api.keys())
			}
			if _, cached := cachedMatchVrf(mgr, "vrf_clause"); cached {
				t.Fatal("the stale clause response must not stay cached")
			}
		})
	}
}

func TestRestoreRetriesAStaleRead(t *testing.T) {
	t.Parallel()
	api, server := newCreateAPI(t)
	api.staleGets = 1
	mgr := newDatacenterManager(server)
	ctx := context.Background()
	clause := mgr.AddPut(ctx, "route_map_clause", "vrf_clause", clauseValue("tenant_2"))
	mgr.AddPut(ctx, "tenant", "tenant_2", zeroPutValue("tenant"))

	if diags, _ := mgr.ExecuteDatacenterOperations(ctx); diags.HasError() {
		t.Fatalf("unexpected error: %v", diags)
	}
	if err := settled(t, mgr, clause); err != nil {
		t.Fatal(err)
	}
	if value, _ := cachedMatchVrf(mgr, "vrf_clause"); value != "tenant_2" {
		t.Fatalf("cached match_vrf = %q, want the confirmed value", value)
	}
}
