package bulkops_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"terraform-provider-verity/internal/bulkops"
	"terraform-provider-verity/openapi"
)

type requestRecord struct {
	Method string
	Path   string
}

// resourceAPIPath maps resource types to their API paths.
var resourceAPIPath = map[string]string{
	"ipv6_prefix_list":         "/ipv6prefixlists",
	"community_list":           "/communitylists",
	"ipv4_prefix_list":         "/ipv4prefixlists",
	"extended_community_list":  "/extendedcommunitylists",
	"as_path_access_list":      "/aspathaccesslists",
	"route_map_clause":         "/routemapclauses",
	"acl":                      "/acls",
	"route_map":                "/routemaps",
	"pb_routing_acl":           "/policybasedroutingacl",
	"tenant":                   "/tenants",
	"pb_routing":               "/policybasedrouting",
	"ipv4_list":                "/ipv4lists",
	"ipv6_list":                "/ipv6lists",
	"service":                  "/services",
	"port_acl":                 "/portacls",
	"packet_broker":            "/packetbroker",
	"eth_port_profile":         "/ethportprofiles",
	"packet_queue":             "/packetqueues",
	"sflow_collector":          "/sflowcollectors",
	"gateway":                  "/gateways",
	"lag":                      "/lags",
	"eth_port_settings":        "/ethportsettings",
	"diagnostics_profile":      "/diagnosticsprofiles",
	"gateway_profile":          "/gatewayprofiles",
	"device_settings":          "/devicesettings",
	"diagnostics_port_profile": "/diagnosticsportprofiles",
	"bundle":                   "/bundles",
	"pod":                      "/pods",
	"badge":                    "/badges",
	"spine_plane":              "/spineplanes",
	"switchpoint":              "/switchpoints",
	"threshold":                "/thresholds",
	"grouping_rule":            "/groupingrules",
	"threshold_group":          "/thresholdgroups",
	"device_controller":        "/devicecontrollers",
	"sfp_breakout":             "/sfpbreakouts",
	"site":                     "/sites",
	"service_port_profile":     "/serviceportprofiles",
	"device_voice_settings":    "/devicevoicesettings",
	"authenticated_eth_port":   "/authenticatedethports",
	"voice_port_profile":       "/voiceportprofiles",
}

var dcPutOrder = []string{
	"ipv6_prefix_list",
	"community_list",
	"ipv4_prefix_list",
	"extended_community_list",
	"as_path_access_list",
	"route_map_clause",
	"acl",
	"route_map",
	"pb_routing_acl",
	"tenant",
	"pb_routing",
	"ipv4_list",
	"ipv6_list",
	"service",
	"port_acl",
	"packet_broker",
	"eth_port_profile",
	"packet_queue",
	"sflow_collector",
	"gateway",
	"lag",
	"eth_port_settings",
	"diagnostics_profile",
	"gateway_profile",
	"device_settings",
	"diagnostics_port_profile",
	"bundle",
	"pod",
	"badge",
	"spine_plane",
	"switchpoint",
	"threshold",
	"grouping_rule",
	"threshold_group",
	"device_controller",
}

var campusPutOrder = []string{
	"ipv4_list",
	"ipv6_list",
	"acl",
	"pb_routing_acl",
	"pb_routing",
	"port_acl",
	"service",
	"eth_port_profile",
	"sflow_collector",
	"packet_queue",
	"service_port_profile",
	"diagnostics_port_profile",
	"device_voice_settings",
	"authenticated_eth_port",
	"diagnostics_profile",
	"eth_port_settings",
	"voice_port_profile",
	"device_settings",
	"lag",
	"bundle",
	"badge",
	"switchpoint",
	"threshold",
	"grouping_rule",
	"threshold_group",
	"device_controller",
}

var dcDeleteOrder = func() []string {
	r := make([]string, len(dcPutOrder))
	for i, v := range dcPutOrder {
		r[len(dcPutOrder)-1-i] = v
	}
	return r
}()

var campusDeleteOrder = func() []string {
	r := make([]string, len(campusPutOrder))
	for i, v := range campusPutOrder {
		r[len(campusPutOrder)-1-i] = v
	}
	return r
}()

// dcPatchOrder is dcPutOrder with "sfp_breakout" prepended and "site" inserted
// before "device_controller" (both are PATCH-only resources).
var dcPatchOrder = func() []string {
	result := make([]string, 0, len(dcPutOrder)+2)
	result = append(result, "sfp_breakout")
	for _, rt := range dcPutOrder {
		if rt == "device_controller" {
			result = append(result, "site")
		}
		result = append(result, rt)
	}
	return result
}()

// campusPatchOrder is campusPutOrder with "site" inserted before "device_controller".
var campusPatchOrder = func() []string {
	result := make([]string, 0, len(campusPutOrder)+1)
	for _, rt := range campusPutOrder {
		if rt == "device_controller" {
			result = append(result, "site")
		}
		result = append(result, rt)
	}
	return result
}()

func orderTrackingServer(t *testing.T) (*httptest.Server, *[]requestRecord, *sync.Mutex) {
	t.Helper()
	var mu sync.Mutex
	var records []requestRecord

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		records = append(records, requestRecord{Method: r.Method, Path: r.URL.Path})
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)
	return server, &records, &mu
}

func newTestClient(serverURL string) *openapi.APIClient {
	cfg := openapi.NewConfiguration()
	cfg.Servers = openapi.ServerConfigurations{{URL: serverURL}}
	cfg.HTTPClient = &http.Client{}
	return openapi.NewAPIClient(cfg)
}

func nopClearCache(_ context.Context, _ interface{}, _ string) {}

func zeroPutValue(resourceType string) interface{} {
	switch resourceType {
	case "gateway":
		return *openapi.NewGatewaysPutRequestGatewayValue()
	case "lag":
		return *openapi.NewLagsPutRequestLagValue()
	case "tenant":
		return *openapi.NewTenantsPutRequestTenantValue()
	case "service":
		return *openapi.NewServicesPutRequestServiceValue()
	case "gateway_profile":
		return *openapi.NewGatewayprofilesPutRequestGatewayProfileValue()
	case "eth_port_profile":
		return *openapi.NewEthportprofilesPutRequestEthPortProfileValue()
	case "eth_port_settings":
		return *openapi.NewEthportsettingsPutRequestEthPortSettingsValue()
	case "device_settings":
		return *openapi.NewDevicesettingsPutRequestEthDeviceProfilesValue()
	case "bundle":
		return *openapi.NewBundlesPutRequestEndpointBundleValue()
	case "acl":
		return *openapi.NewAclsPutRequestIpFilterValue()
	case "ipv4_list":
		return *openapi.NewIpv4listsPutRequestIpv4ListFilterValue()
	case "ipv4_prefix_list":
		return *openapi.NewIpv4prefixlistsPutRequestIpv4PrefixListValue()
	case "ipv6_list":
		return *openapi.NewIpv6listsPutRequestIpv6ListFilterValue()
	case "ipv6_prefix_list":
		return *openapi.NewIpv6prefixlistsPutRequestIpv6PrefixListValue()
	case "authenticated_eth_port":
		return *openapi.NewAuthenticatedethportsPutRequestAuthenticatedEthPortValue()
	case "badge":
		return *openapi.NewBadgesPutRequestBadgeValue()
	case "device_voice_settings":
		return *openapi.NewDevicevoicesettingsPutRequestDeviceVoiceSettingsValue()
	case "as_path_access_list":
		return *openapi.NewAspathaccesslistsPutRequestAsPathAccessListValue()
	case "community_list":
		return *openapi.NewCommunitylistsPutRequestCommunityListValue()
	case "extended_community_list":
		return *openapi.NewExtendedcommunitylistsPutRequestExtendedCommunityListValue()
	case "route_map_clause":
		return *openapi.NewRoutemapclausesPutRequestRouteMapClauseValue()
	case "route_map":
		return *openapi.NewRoutemapsPutRequestRouteMapValue()
	case "packet_broker":
		return *openapi.NewPacketbrokerPutRequestPbEgressProfileValue()
	case "packet_queue":
		return *openapi.NewPacketqueuesPutRequestPacketQueueValue()
	case "service_port_profile":
		return *openapi.NewServiceportprofilesPutRequestServicePortProfileValue()
	case "switchpoint":
		return *openapi.NewSwitchpointsPutRequestSwitchpointValue()
	case "voice_port_profile":
		return *openapi.NewVoiceportprofilesPutRequestVoicePortProfilesValue()
	case "device_controller":
		return *openapi.NewDevicecontrollersPutRequestDeviceControllerValue()
	case "pod":
		return *openapi.NewPodsPutRequestPodValue()
	case "port_acl":
		return *openapi.NewPortaclsPutRequestPortAclValue()
	case "sflow_collector":
		return *openapi.NewSflowcollectorsPutRequestSflowCollectorValue()
	case "diagnostics_profile":
		return *openapi.NewDiagnosticsprofilesPutRequestDiagnosticsProfileValue()
	case "diagnostics_port_profile":
		return *openapi.NewDiagnosticsportprofilesPutRequestDiagnosticsPortProfileValue()
	case "pb_routing":
		return *openapi.NewPolicybasedroutingPutRequestPbRoutingValue()
	case "pb_routing_acl":
		return *openapi.NewPolicybasedroutingaclPutRequestPbRoutingAclValue()
	case "spine_plane":
		return *openapi.NewSpineplanesPutRequestSpinePlaneValue()
	case "grouping_rule":
		return *openapi.NewGroupingrulesPutRequestGroupingRulesValue()
	case "threshold_group":
		return *openapi.NewThresholdgroupsPutRequestThresholdGroupValue()
	case "threshold":
		return *openapi.NewThresholdsPutRequestThresholdValue()
	default:
		panic(fmt.Sprintf("unknown resource type in zeroPutValue: %s", resourceType))
	}
}

func zeroPatchValue(resourceType string) interface{} {
	switch resourceType {
	case "sfp_breakout":
		return *openapi.NewSfpbreakoutsPatchRequestSfpBreakoutsValue()
	case "site":
		return *openapi.NewSitesPatchRequestSiteValue()
	default:
		return zeroPutValue(resourceType)
	}
}

func filterRecords(records []requestRecord, method string) []string {
	var paths []string
	for _, r := range records {
		if r.Method == method {
			paths = append(paths, r.Path)
		}
	}
	return paths
}

func assertOrderedSubset(t *testing.T, label string, actual, expected []string) {
	t.Helper()
	idxMap := make(map[string]int)
	for i, p := range actual {
		if _, exists := idxMap[p]; !exists {
			idxMap[p] = i
		}
	}
	for i, path := range expected {
		if _, exists := idxMap[path]; !exists {
			t.Errorf("%s: expected path %q not found in actual requests", label, path)
			return
		}
		if i > 0 {
			prevIdx := idxMap[expected[i-1]]
			curIdx := idxMap[path]
			if curIdx <= prevIdx {
				t.Errorf("%s: expected %q (idx %d) after %q (idx %d), but order is wrong.\n  actual order: %v",
					label, path, curIdx, expected[i-1], prevIdx, actual)
				return
			}
		}
	}
}

func snapshotRecords(mu *sync.Mutex, records *[]requestRecord) []requestRecord {
	mu.Lock()
	defer mu.Unlock()
	out := make([]requestRecord, len(*records))
	copy(out, *records)
	return out
}

func toPaths(resources []string) []string {
	var paths []string
	for _, rt := range resources {
		if path, ok := resourceAPIPath[rt]; ok {
			paths = append(paths, path)
		}
	}
	return paths
}

func assertPhaseOrder(t *testing.T, records []requestRecord) {
	t.Helper()
	lastPutIdx := -1
	firstPatchIdx := -1
	lastPatchIdx := -1
	firstDeleteIdx := -1

	for i, r := range records {
		switch r.Method {
		case http.MethodPut:
			lastPutIdx = i
		case http.MethodPatch:
			if firstPatchIdx == -1 {
				firstPatchIdx = i
			}
			lastPatchIdx = i
		case http.MethodDelete:
			if firstDeleteIdx == -1 {
				firstDeleteIdx = i
			}
		}
	}

	if lastPutIdx == -1 || firstPatchIdx == -1 || firstDeleteIdx == -1 {
		t.Fatalf("expected PUT, PATCH, and DELETE requests; got PUTs=%v PATCHes=%v DELETEs=%v",
			lastPutIdx != -1, firstPatchIdx != -1, firstDeleteIdx != -1)
	}
	if lastPutIdx >= firstPatchIdx {
		t.Errorf("last PUT (idx %d) should come before first PATCH (idx %d)", lastPutIdx, firstPatchIdx)
	}
	if lastPatchIdx >= firstDeleteIdx {
		t.Errorf("last PATCH (idx %d) should come before first DELETE (idx %d)", lastPatchIdx, firstDeleteIdx)
	}
}

func addPutsForResources(ctx context.Context, mgr *bulkops.Manager, resources []string) {
	for _, rt := range resources {
		if rt == "acl" {
			mgr.AddPut(ctx, rt, fmt.Sprintf("test_%s_v4", rt), zeroPutValue(rt), map[string]string{"ip_version": "4"})
		} else {
			mgr.AddPut(ctx, rt, fmt.Sprintf("test_%s", rt), zeroPutValue(rt))
		}
	}
}

func addPatchesForResources(ctx context.Context, mgr *bulkops.Manager, resources []string) {
	for _, rt := range resources {
		if rt == "acl" {
			mgr.AddPatch(ctx, rt, fmt.Sprintf("test_%s_v4", rt), zeroPatchValue(rt), map[string]string{"ip_version": "4"})
		} else {
			mgr.AddPatch(ctx, rt, fmt.Sprintf("test_%s", rt), zeroPatchValue(rt))
		}
	}
}

func addDeletesForResources(ctx context.Context, mgr *bulkops.Manager, resources []string) {
	for _, rt := range resources {
		if rt == "acl" {
			mgr.AddDelete(ctx, rt, fmt.Sprintf("test_%s", rt), map[string]string{"ip_version": "4"})
		} else {
			mgr.AddDelete(ctx, rt, fmt.Sprintf("test_%s", rt))
		}
	}
}

func failingOrderTrackingServer(t *testing.T, shouldFail func(*http.Request) bool) (*httptest.Server, *[]requestRecord, *sync.Mutex) {
	t.Helper()
	var mu sync.Mutex
	var records []requestRecord

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		records = append(records, requestRecord{Method: r.Method, Path: r.URL.Path})
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		if shouldFail(r) {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(`{"error":"simulated failure"}`))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	t.Cleanup(server.Close)
	return server, &records, &mu
}

func TestDatacenterPutOrder(t *testing.T) {
	t.Parallel()
	server, records, mu := orderTrackingServer(t)
	client := newTestClient(server.URL)
	mgr := bulkops.GetManager(client, nopClearCache, nil, "datacenter")

	ctx := context.Background()

	addPutsForResources(ctx, mgr, dcPutOrder)

	diags, _ := mgr.ExecuteDatacenterOperations(ctx)
	if diags.HasError() {
		t.Fatalf("ExecuteDatacenterOperations returned errors: %v", diags)
	}

	putPaths := filterRecords(snapshotRecords(mu, records), http.MethodPut)
	assertOrderedSubset(t, "DC PUT order", putPaths, toPaths(dcPutOrder))
}

func TestDatacenterDeleteOrder(t *testing.T) {
	t.Parallel()
	server, records, mu := orderTrackingServer(t)
	client := newTestClient(server.URL)
	mgr := bulkops.GetManager(client, nopClearCache, nil, "datacenter")

	ctx := context.Background()

	addDeletesForResources(ctx, mgr, dcDeleteOrder)

	diags, _ := mgr.ExecuteDatacenterOperations(ctx)
	if diags.HasError() {
		t.Fatalf("ExecuteDatacenterOperations returned errors: %v", diags)
	}

	deletePaths := filterRecords(snapshotRecords(mu, records), http.MethodDelete)
	assertOrderedSubset(t, "DC DELETE order", deletePaths, toPaths(dcDeleteOrder))
}

func TestCampusPutOrder(t *testing.T) {
	t.Parallel()
	server, records, mu := orderTrackingServer(t)
	client := newTestClient(server.URL)
	mgr := bulkops.GetManager(client, nopClearCache, nil, "campus")

	ctx := context.Background()

	addPutsForResources(ctx, mgr, campusPutOrder)

	diags, _ := mgr.ExecuteCampusOperations(ctx)
	if diags.HasError() {
		t.Fatalf("ExecuteCampusOperations returned errors: %v", diags)
	}

	putPaths := filterRecords(snapshotRecords(mu, records), http.MethodPut)
	assertOrderedSubset(t, "campus PUT order", putPaths, toPaths(campusPutOrder))
}

func TestCampusDeleteOrder(t *testing.T) {
	t.Parallel()
	server, records, mu := orderTrackingServer(t)
	client := newTestClient(server.URL)
	mgr := bulkops.GetManager(client, nopClearCache, nil, "campus")

	ctx := context.Background()

	addDeletesForResources(ctx, mgr, campusDeleteOrder)

	diags, _ := mgr.ExecuteCampusOperations(ctx)
	if diags.HasError() {
		t.Fatalf("ExecuteCampusOperations returned errors: %v", diags)
	}

	deletePaths := filterRecords(snapshotRecords(mu, records), http.MethodDelete)
	assertOrderedSubset(t, "campus DELETE order", deletePaths, toPaths(campusDeleteOrder))
}

func TestDatacenterPatchOrder(t *testing.T) {
	t.Parallel()
	server, records, mu := orderTrackingServer(t)
	client := newTestClient(server.URL)
	mgr := bulkops.GetManager(client, nopClearCache, nil, "datacenter")

	ctx := context.Background()

	addPatchesForResources(ctx, mgr, dcPatchOrder)

	diags, _ := mgr.ExecuteDatacenterOperations(ctx)
	if diags.HasError() {
		t.Fatalf("ExecuteDatacenterOperations returned errors: %v", diags)
	}

	patchPaths := filterRecords(snapshotRecords(mu, records), http.MethodPatch)
	assertOrderedSubset(t, "DC PATCH order", patchPaths, toPaths(dcPatchOrder))
}

func TestErrorAbortsRemainingOperations(t *testing.T) {
	t.Parallel()
	server, records, mu := failingOrderTrackingServer(t, func(r *http.Request) bool {
		return r.Method == http.MethodPut && r.URL.Path == "/tenants"
	})

	client := newTestClient(server.URL)
	mgr := bulkops.GetManager(client, nopClearCache, nil, "datacenter")

	ctx := context.Background()

	// Add PUTs for resources before and after tenant in the DC order
	mgr.AddPut(ctx, "ipv6_prefix_list", "test_prefix", zeroPutValue("ipv6_prefix_list"))
	mgr.AddPut(ctx, "route_map_clause", "test_clause", zeroPutValue("route_map_clause"))
	mgr.AddPut(ctx, "tenant", "test_tenant", zeroPutValue("tenant"))
	// These come AFTER tenant in the DC order — should NOT execute
	mgr.AddPut(ctx, "pb_routing", "test_pbr", zeroPutValue("pb_routing"))
	mgr.AddPut(ctx, "service", "test_service", zeroPutValue("service"))
	mgr.AddPut(ctx, "gateway", "test_gw", zeroPutValue("gateway"))

	diags, _ := mgr.ExecuteDatacenterOperations(ctx)
	if !diags.HasError() {
		t.Fatal("expected errors from failed tenant PUT, got none")
	}

	putPaths := filterRecords(snapshotRecords(mu, records), http.MethodPut)

	for _, path := range putPaths {
		if path == "/policybasedrouting" || path == "/services" || path == "/gateways" {
			t.Errorf("resource at path %q should NOT have been called after tenant failure", path)
		}
	}

	foundPrefix := false
	foundClause := false
	foundTenant := false
	for _, path := range putPaths {
		switch path {
		case "/ipv6prefixlists":
			foundPrefix = true
		case "/routemapclauses":
			foundClause = true
		case "/tenants":
			foundTenant = true
		}
	}
	if !foundPrefix {
		t.Error("ipv6_prefix_list PUT should have been called before tenant failure")
	}
	if !foundClause {
		t.Error("route_map_clause PUT should have been called before tenant failure")
	}
	if !foundTenant {
		t.Error("tenant PUT should have been attempted")
	}
}

// TestCircularReferencePutFix verifies that when route_map_clauses reference tenants
// being created in the same batch, the circular reference fix is applied:
// 1. route_map_clause PUT with empty match_vrf
// 2. tenant PUT
// 3. route_map_clause PATCH to restore match_vrf
func TestCircularReferencePutFix(t *testing.T) {
	t.Parallel()
	var mu sync.Mutex
	var records []requestRecord
	var capturedBodies []struct {
		Method string
		Path   string
		Body   map[string]interface{}
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		records = append(records, requestRecord{Method: r.Method, Path: r.URL.Path})

		if r.URL.Path == "/routemapclauses" && (r.Method == http.MethodPut || r.Method == http.MethodPatch) {
			var body map[string]interface{}
			if err := json.NewDecoder(r.Body).Decode(&body); err == nil {
				capturedBodies = append(capturedBodies, struct {
					Method string
					Path   string
					Body   map[string]interface{}
				}{Method: r.Method, Path: r.URL.Path, Body: body})
			}
		}
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	mgr := bulkops.GetManager(client, nopClearCache, nil, "datacenter")

	ctx := context.Background()

	// Add a route_map_clause with match_vrf referencing a tenant being created
	clauseVal := openapi.RoutemapclausesPutRequestRouteMapClauseValue{
		Name:     openapi.PtrString("clause_1"),
		MatchVrf: openapi.PtrString("tenant_1"),
	}
	mgr.AddPut(ctx, "route_map_clause", "clause_1", clauseVal)

	// Add the tenant that the clause references
	tenantVal := openapi.TenantsPutRequestTenantValue{
		Name: openapi.PtrString("tenant_1"),
	}
	mgr.AddPut(ctx, "tenant", "tenant_1", tenantVal)

	diags, _ := mgr.ExecuteDatacenterOperations(ctx)
	if diags.HasError() {
		t.Fatalf("ExecuteDatacenterOperations returned errors: %v", diags)
	}

	mu.Lock()
	allRecords := make([]requestRecord, len(records))
	copy(allRecords, records)
	bodies := make([]struct {
		Method string
		Path   string
		Body   map[string]interface{}
	}, len(capturedBodies))
	copy(bodies, capturedBodies)
	mu.Unlock()

	var clausePutIdx, tenantPutIdx, clausePatchIdx int
	clausePutIdx, tenantPutIdx, clausePatchIdx = -1, -1, -1

	for i, r := range allRecords {
		if r.Method == http.MethodPut && r.Path == "/routemapclauses" {
			clausePutIdx = i
		}
		if r.Method == http.MethodPut && r.Path == "/tenants" {
			tenantPutIdx = i
		}
		if r.Method == http.MethodPatch && r.Path == "/routemapclauses" {
			clausePatchIdx = i
		}
	}

	if clausePutIdx == -1 {
		t.Fatal("route_map_clause PUT not found")
	}
	if tenantPutIdx == -1 {
		t.Fatal("tenant PUT not found")
	}
	if clausePatchIdx == -1 {
		t.Fatal("route_map_clause PATCH not found (circular ref fix should restore match_vrf)")
	}

	if clausePutIdx >= tenantPutIdx {
		t.Errorf("route_map_clause PUT (idx %d) should come before tenant PUT (idx %d)", clausePutIdx, tenantPutIdx)
	}

	if tenantPutIdx >= clausePatchIdx {
		t.Errorf("tenant PUT (idx %d) should come before route_map_clause PATCH (idx %d)", tenantPutIdx, clausePatchIdx)
	}

	if len(bodies) > 0 {
		for _, b := range bodies {
			if b.Method == http.MethodPut {
				if clause, ok := b.Body["route_map_clause"]; ok {
					if clauseMap, ok := clause.(map[string]interface{}); ok {
						for _, v := range clauseMap {
							if vMap, ok := v.(map[string]interface{}); ok {
								if vrf, exists := vMap["match_vrf"]; exists {
									if vrfStr, ok := vrf.(string); ok && vrfStr != "" {
										t.Errorf("route_map_clause PUT should have empty match_vrf, got %q", vrfStr)
									}
								}
							}
						}
					}
				}
			}
		}
	}
}

func TestMixedOperationsOrder(t *testing.T) {
	t.Parallel()
	server, records, mu := orderTrackingServer(t)
	client := newTestClient(server.URL)
	mgr := bulkops.GetManager(client, nopClearCache, nil, "datacenter")

	ctx := context.Background()

	// Add a mix of PUT, PATCH, and DELETE for different resources
	mgr.AddPut(ctx, "badge", "put_badge", zeroPutValue("badge"))
	mgr.AddPut(ctx, "gateway", "put_gw", zeroPutValue("gateway"))
	mgr.AddPatch(ctx, "badge", "patch_badge", zeroPatchValue("badge"))
	mgr.AddPatch(ctx, "gateway", "patch_gw", zeroPatchValue("gateway"))
	mgr.AddDelete(ctx, "badge", "del_badge")
	mgr.AddDelete(ctx, "gateway", "del_gw")

	diags, _ := mgr.ExecuteDatacenterOperations(ctx)
	if diags.HasError() {
		t.Fatalf("ExecuteDatacenterOperations returned errors: %v", diags)
	}

	assertPhaseOrder(t, snapshotRecords(mu, records))
}

func TestACLHeaderSplitExecution(t *testing.T) {
	t.Parallel()
	var mu sync.Mutex
	var records []requestRecord
	var aclHeaders []string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		records = append(records, requestRecord{Method: r.Method, Path: r.URL.Path})
		if r.URL.Path == "/acls" && r.Method == http.MethodPut {
			aclHeaders = append(aclHeaders, r.URL.Query().Get("ip_version"))
		}
		mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}))
	defer server.Close()

	client := newTestClient(server.URL)
	mgr := bulkops.GetManager(client, nopClearCache, nil, "datacenter")

	ctx := context.Background()

	mgr.AddPut(ctx, "acl", "v4_filter", zeroPutValue("acl"), map[string]string{"ip_version": "4"})
	mgr.AddPut(ctx, "acl", "v6_filter", zeroPutValue("acl"), map[string]string{"ip_version": "6"})

	diags, _ := mgr.ExecuteDatacenterOperations(ctx)
	if diags.HasError() {
		t.Fatalf("ExecuteDatacenterOperations returned errors: %v", diags)
	}

	mu.Lock()
	headersCopy := make([]string, len(aclHeaders))
	copy(headersCopy, aclHeaders)
	allRecords := make([]requestRecord, len(records))
	copy(allRecords, records)
	mu.Unlock()

	aclPutCount := 0
	for _, r := range allRecords {
		if r.Path == "/acls" && r.Method == http.MethodPut {
			aclPutCount++
		}
	}

	if aclPutCount != 2 {
		t.Errorf("expected 2 ACL PUT calls (one per ip_version), got %d", aclPutCount)
	}

	has4, has6 := false, false
	for _, h := range headersCopy {
		if h == "4" {
			has4 = true
		}
		if h == "6" {
			has6 = true
		}
	}
	if !has4 {
		t.Error("expected ACL PUT with ip_version=4")
	}
	if !has6 {
		t.Error("expected ACL PUT with ip_version=6")
	}
}

func TestNoOpsSkipped(t *testing.T) {
	t.Parallel()
	server, records, mu := orderTrackingServer(t)
	client := newTestClient(server.URL)
	mgr := bulkops.GetManager(client, nopClearCache, nil, "datacenter")

	ctx := context.Background()

	// Only add badge PUT — no other resources
	mgr.AddPut(ctx, "badge", "test_badge", zeroPutValue("badge"))

	diags, _ := mgr.ExecuteDatacenterOperations(ctx)
	if diags.HasError() {
		t.Fatalf("ExecuteDatacenterOperations returned errors: %v", diags)
	}

	allRecords := snapshotRecords(mu, records)

	putPaths := filterRecords(allRecords, http.MethodPut)
	if len(putPaths) != 1 {
		t.Errorf("expected exactly 1 PUT call, got %d: %v", len(putPaths), putPaths)
	}
	if len(putPaths) == 1 && putPaths[0] != "/badges" {
		t.Errorf("expected PUT to /badges, got %s", putPaths[0])
	}

	patchPaths := filterRecords(allRecords, http.MethodPatch)
	deletePaths := filterRecords(allRecords, http.MethodDelete)
	if len(patchPaths) > 0 {
		t.Errorf("expected no PATCH calls, got %d: %v", len(patchPaths), patchPaths)
	}
	if len(deletePaths) > 0 {
		t.Errorf("expected no DELETE calls, got %d: %v", len(deletePaths), deletePaths)
	}
}

func TestCampusPatchOrder(t *testing.T) {
	t.Parallel()
	server, records, mu := orderTrackingServer(t)
	client := newTestClient(server.URL)
	mgr := bulkops.GetManager(client, nopClearCache, nil, "campus")

	ctx := context.Background()

	addPatchesForResources(ctx, mgr, campusPatchOrder)

	diags, _ := mgr.ExecuteCampusOperations(ctx)
	if diags.HasError() {
		t.Fatalf("ExecuteCampusOperations returned errors: %v", diags)
	}

	patchPaths := filterRecords(snapshotRecords(mu, records), http.MethodPatch)
	assertOrderedSubset(t, "campus PATCH order", patchPaths, toPaths(campusPatchOrder))
}

func TestCampusMixedOperationsOrder(t *testing.T) {
	t.Parallel()
	server, records, mu := orderTrackingServer(t)
	client := newTestClient(server.URL)
	mgr := bulkops.GetManager(client, nopClearCache, nil, "campus")

	ctx := context.Background()

	mgr.AddPut(ctx, "badge", "put_badge", zeroPutValue("badge"))
	mgr.AddPut(ctx, "service", "put_service", zeroPutValue("service"))
	mgr.AddPatch(ctx, "badge", "patch_badge", zeroPatchValue("badge"))
	mgr.AddPatch(ctx, "service", "patch_service", zeroPatchValue("service"))
	mgr.AddDelete(ctx, "badge", "del_badge")
	mgr.AddDelete(ctx, "service", "del_service")

	diags, _ := mgr.ExecuteCampusOperations(ctx)
	if diags.HasError() {
		t.Fatalf("ExecuteCampusOperations returned errors: %v", diags)
	}

	assertPhaseOrder(t, snapshotRecords(mu, records))
}

func TestPatchErrorAbortsRemainingOperations(t *testing.T) {
	t.Parallel()
	server, records, mu := failingOrderTrackingServer(t, func(r *http.Request) bool {
		return r.Method == http.MethodPatch && r.URL.Path == "/routemapclauses"
	})

	client := newTestClient(server.URL)
	mgr := bulkops.GetManager(client, nopClearCache, nil, "datacenter")

	ctx := context.Background()

	// PUT succeeds for a resource that is also PATCHed
	mgr.AddPut(ctx, "ipv6_prefix_list", "test_prefix", zeroPutValue("ipv6_prefix_list"))
	// PATCH: early one succeeds, route_map_clause fails, later ones should be skipped
	mgr.AddPatch(ctx, "ipv6_prefix_list", "test_prefix", zeroPatchValue("ipv6_prefix_list"))
	mgr.AddPatch(ctx, "route_map_clause", "test_clause", zeroPatchValue("route_map_clause"))
	// These come AFTER route_map_clause in DC PATCH order — should NOT execute
	mgr.AddPatch(ctx, "route_map", "test_rm", zeroPatchValue("route_map"))
	mgr.AddPatch(ctx, "tenant", "test_tenant", zeroPatchValue("tenant"))
	// DELETE should also NOT execute after PATCH failure
	mgr.AddDelete(ctx, "badge", "test_badge")

	diags, _ := mgr.ExecuteDatacenterOperations(ctx)
	if !diags.HasError() {
		t.Fatal("expected errors from failed route_map_clause PATCH, got none")
	}

	allRecords := snapshotRecords(mu, records)
	patchPaths := filterRecords(allRecords, http.MethodPatch)
	deletePaths := filterRecords(allRecords, http.MethodDelete)

	for _, path := range patchPaths {
		if path == "/routemaps" || path == "/tenants" {
			t.Errorf("path %q should NOT have been patched after route_map_clause PATCH failure", path)
		}
	}

	if len(deletePaths) > 0 {
		t.Errorf("expected no DELETE calls after PATCH failure, got %v", deletePaths)
	}

	foundPrefix := false
	foundClause := false
	for _, path := range patchPaths {
		switch path {
		case "/ipv6prefixlists":
			foundPrefix = true
		case "/routemapclauses":
			foundClause = true
		}
	}
	if !foundPrefix {
		t.Error("ipv6_prefix_list PATCH should have been called before route_map_clause failure")
	}
	if !foundClause {
		t.Error("route_map_clause PATCH should have been attempted")
	}
}

func TestDeleteErrorAbortsRemainingOperations(t *testing.T) {
	t.Parallel()
	server, records, mu := failingOrderTrackingServer(t, func(r *http.Request) bool {
		return r.Method == http.MethodDelete && r.URL.Path == "/devicecontrollers"
	})

	client := newTestClient(server.URL)
	mgr := bulkops.GetManager(client, nopClearCache, nil, "datacenter")

	ctx := context.Background()

	// device_controller is first in reverse DC order — its failure should abort the rest
	mgr.AddDelete(ctx, "device_controller", "test_dc")
	// These come AFTER device_controller in reverse DC order — should NOT execute
	mgr.AddDelete(ctx, "threshold_group", "test_tg")
	mgr.AddDelete(ctx, "badge", "test_badge")

	diags, _ := mgr.ExecuteDatacenterOperations(ctx)
	if !diags.HasError() {
		t.Fatal("expected errors from failed device_controller DELETE, got none")
	}

	deletePaths := filterRecords(snapshotRecords(mu, records), http.MethodDelete)

	for _, path := range deletePaths {
		if path == "/thresholdgroups" || path == "/badges" {
			t.Errorf("path %q should NOT have been deleted after device_controller DELETE failure", path)
		}
	}

	foundDC := false
	for _, path := range deletePaths {
		if path == "/devicecontrollers" {
			foundDC = true
		}
	}
	if !foundDC {
		t.Error("device_controller DELETE should have been attempted")
	}
}
