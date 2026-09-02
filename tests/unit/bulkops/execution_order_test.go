package bulkops_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"terraform-provider-verity/internal/bulkops"
	"terraform-provider-verity/openapi"
)

type requestRecord struct {
	Method    string
	Path      string
	IPVersion string
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
	"tacacs_profile":           "/tacacsprofiles",
	"ldap_profile":             "/ldapprofiles",
	"eth_port_profile":         "/ethportprofiles",
	"packet_queue":             "/packetqueues",
	"sflow_collector":          "/sflowcollectors",
	"gateway":                  "/gateways",
	"device_aaa_profile":       "/deviceaaaprofiles",
	"lag":                      "/lags",
	"eth_port_settings":        "/ethportsettings",
	"diagnostics_profile":      "/diagnosticsprofiles",
	"gateway_profile":          "/gatewayprofiles",
	"device_settings":          "/devicesettings",
	"diagnostics_port_profile": "/diagnosticsportprofiles",
	"bundle":                   "/bundles",
	"pod":                      "/pods",
	"badge":                    "/badges",
	"su":                       "/sus",
	"ssp_group":                "/sspgroups",
	"pair":                     "/pairs",
	"mac_filter":               "/macfilters",
	"spine_plane":              "/spineplanes",
	"switchpoint":              "/switchpoints",
	"threshold":                "/thresholds",
	"grouping_rule":            "/groupingrules",
	"threshold_group":          "/thresholdgroups",
	"sfp_breakout":             "/sfpbreakouts",
	"fabric":                   "/fabrics",
	"plane":                    "/planes",
	"rack":                     "/racks",
	"service_port_profile":     "/serviceportprofiles",
	"device_voice_settings":    "/devicevoicesettings",
	"authenticated_eth_port":   "/authenticatedethports",
	"voice_port_profile":       "/voiceportprofiles",
}

var dcPutOrder = []string{
	"community_list",
	"as_path_access_list",
	"ipv6_prefix_list",
	"ipv4_prefix_list",
	"extended_community_list",
	"acl",
	"route_map_clause",
	"pb_routing_acl",
	"route_map",
	"pb_routing",
	"tenant",
	"service",
	"fabric",
	"tacacs_profile",
	"ldap_profile",
	"port_acl",
	"ipv6_list",
	"ipv4_list",
	"pod",
	"packet_queue",
	"device_aaa_profile",
	"eth_port_profile",
	"packet_broker",
	"sflow_collector",
	"gateway",
	"su",
	"diagnostics_port_profile",
	"device_settings",
	"lag",
	"diagnostics_profile",
	"gateway_profile",
	"eth_port_settings",
	"badge",
	"plane",
	"spine_plane",
	"rack",
	"bundle",
	"ssp_group",
	"grouping_rule",
	"switchpoint",
	"threshold",
	"threshold_group",
	"pair",
}

var campusPutOrder = []string{
	"acl",
	"mac_filter",
	"service",
	"port_acl",
	"tacacs_profile",
	"ldap_profile",
	"sflow_collector",
	"eth_port_profile",
	"packet_queue",
	"device_aaa_profile",
	"fabric",
	"service_port_profile",
	"diagnostics_profile",
	"authenticated_eth_port",
	"device_settings",
	"voice_port_profile",
	"lag",
	"device_voice_settings",
	"eth_port_settings",
	"diagnostics_port_profile",
	"bundle",
	"badge",
	"grouping_rule",
	"switchpoint",
	"threshold",
	"threshold_group",
	"pair",
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

// dcPatchOrder is dcPutOrder with "sfp_breakout" prepended.
var dcPatchOrder = func() []string {
	result := make([]string, 0, len(dcPutOrder)+1)
	result = append(result, "sfp_breakout")
	result = append(result, dcPutOrder...)
	return result
}()

var campusPatchOrder = append([]string{"sfp_breakout"}, campusPutOrder...)

func orderTrackingServer(t *testing.T) (*httptest.Server, *[]requestRecord, *sync.Mutex) {
	t.Helper()
	var mu sync.Mutex
	var records []requestRecord

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		records = append(records, requestRecord{
			Method:    r.Method,
			Path:      r.URL.Path,
			IPVersion: r.Header.Get("ip_version"),
		})
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
	case "device_aaa_profile":
		return *openapi.NewDeviceaaaprofilesPutRequestDeviceAaaProfileValue()
	case "ldap_profile":
		return *openapi.NewLdapprofilesPutRequestLdapProfileValue()
	case "tacacs_profile":
		return *openapi.NewTacacsprofilesPutRequestTacacsProfileValue()
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
	case "mac_filter":
		return *openapi.NewMacfiltersPutRequestMacFilterValue()
	case "pair":
		return *openapi.NewPairsPutRequestSwitchPairValue()
	case "packet_queue":
		return *openapi.NewPacketqueuesPutRequestPacketQueueValue()
	case "service_port_profile":
		return *openapi.NewServiceportprofilesPutRequestServicePortProfileValue()
	case "switchpoint":
		return *openapi.NewSwitchpointsPutRequestSwitchpointValue()
	case "fabric":
		return *openapi.NewFabricsPutRequestFabricValue()
	case "plane":
		return *openapi.NewPlanesPutRequestPlaneValue()
	case "rack":
		return *openapi.NewRacksPutRequestRackValue()
	case "voice_port_profile":
		return *openapi.NewVoiceportprofilesPutRequestVoicePortProfilesValue()
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
	case "ssp_group":
		return *openapi.NewSspgroupsPutRequestSuperspineGroupValue()
	case "su":
		return *openapi.NewSusPutRequestSuValue()
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

func assertPathsAbsent(t *testing.T, label string, actual, forbidden []string) {
	t.Helper()
	for _, path := range forbidden {
		for _, actualPath := range actual {
			if actualPath == path {
				t.Errorf("%s: unexpected request to %q", label, path)
				break
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

func TestACLsPutIPv6BeforeIPv4InEveryMode(t *testing.T) {
	for _, mode := range []string{"campus", "datacenter"} {
		t.Run(mode, func(t *testing.T) {
			server, records, mu := orderTrackingServer(t)
			mgr := bulkops.GetManager(newTestClient(server.URL), nopClearCache, nil, mode)
			ctx := context.Background()

			mgr.AddPut(ctx, "acl", "ipv4", zeroPutValue("acl"), map[string]string{"ip_version": "4"})
			mgr.AddPut(ctx, "acl", "ipv6", zeroPutValue("acl"), map[string]string{"ip_version": "6"})

			var hasErrors bool
			if mode == "campus" {
				d, _ := mgr.ExecuteCampusOperations(ctx)
				hasErrors = d.HasError()
			} else {
				d, _ := mgr.ExecuteDatacenterOperations(ctx)
				hasErrors = d.HasError()
			}
			if hasErrors {
				t.Fatalf("ACL PUT execution returned errors")
			}

			var versions []string
			for _, record := range snapshotRecords(mu, records) {
				if record.Method == http.MethodPut && record.Path == "/acls" {
					versions = append(versions, record.IPVersion)
				}
			}
			if fmt.Sprint(versions) != "[6 4]" {
				t.Fatalf("expected IPv6 ACL PUT before IPv4, got %v", versions)
			}
		})
	}
}

func TestDatacenterFabricPrecedesGateway(t *testing.T) {
	t.Parallel()
	server, records, mu := orderTrackingServer(t)
	client := newTestClient(server.URL)
	mgr := bulkops.GetManager(client, nopClearCache, nil, "datacenter")

	ctx := context.Background()
	mgr.AddPut(ctx, "gateway", "test_gateway", zeroPutValue("gateway"))
	mgr.AddPut(ctx, "fabric", "test_fabric", zeroPutValue("fabric"))

	diags, _ := mgr.ExecuteDatacenterOperations(ctx)
	if diags.HasError() {
		t.Fatalf("ExecuteDatacenterOperations returned errors: %v", diags)
	}

	putPaths := filterRecords(snapshotRecords(mu, records), http.MethodPut)
	assertOrderedSubset(t, "DC Fabric/Gateway dependency", putPaths, []string{"/fabrics", "/gateways"})
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
	mgr.AddPut(ctx, "ipv4_list", "test_ipv4_list", zeroPutValue("ipv4_list"))
	mgr.AddPut(ctx, "ipv6_list", "test_ipv6_list", zeroPutValue("ipv6_list"))

	diags, _ := mgr.ExecuteCampusOperations(ctx)
	if diags.HasError() {
		t.Fatalf("ExecuteCampusOperations returned errors: %v", diags)
	}

	putPaths := filterRecords(snapshotRecords(mu, records), http.MethodPut)
	assertOrderedSubset(t, "campus PUT order", putPaths, toPaths(campusPutOrder))
	assertPathsAbsent(t, "campus PUT order", putPaths, []string{"/ipv4lists", "/ipv6lists"})
}

func TestCampusDeleteOrder(t *testing.T) {
	t.Parallel()
	server, records, mu := orderTrackingServer(t)
	client := newTestClient(server.URL)
	mgr := bulkops.GetManager(client, nopClearCache, nil, "campus")

	ctx := context.Background()

	addDeletesForResources(ctx, mgr, campusDeleteOrder)
	mgr.AddDelete(ctx, "ipv4_list", "test_ipv4_list")
	mgr.AddDelete(ctx, "ipv6_list", "test_ipv6_list")

	diags, _ := mgr.ExecuteCampusOperations(ctx)
	if diags.HasError() {
		t.Fatalf("ExecuteCampusOperations returned errors: %v", diags)
	}

	deletePaths := filterRecords(snapshotRecords(mu, records), http.MethodDelete)
	assertOrderedSubset(t, "campus DELETE order", deletePaths, toPaths(campusDeleteOrder))
	assertPathsAbsent(t, "campus DELETE order", deletePaths, []string{"/ipv4lists", "/ipv6lists"})
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

func TestACLHeaderSplitOrder(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		mode      string
		operation string
		want      []string
	}{
		{"datacenter put", "datacenter", "PUT", []string{"4", "6"}},
		{"campus put", "campus", "PUT", []string{"6", "4"}},
		{"datacenter delete", "datacenter", "DELETE", []string{"6", "4"}},
		{"campus delete", "campus", "DELETE", []string{"4", "6"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var mu sync.Mutex
			var aclHeaders []string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/acls" && r.Method == tt.operation {
					mu.Lock()
					aclHeaders = append(aclHeaders, r.URL.Query().Get("ip_version"))
					mu.Unlock()
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{}`))
			}))
			t.Cleanup(server.Close)

			client := newTestClient(server.URL)
			mgr := bulkops.GetManager(client, nopClearCache, nil, tt.mode)
			ctx := context.Background()
			for _, version := range []string{"4", "6"} {
				if tt.operation == "PUT" {
					mgr.AddPut(ctx, "acl", "filter_v"+version, zeroPutValue("acl"), map[string]string{"ip_version": version})
				} else {
					mgr.AddDelete(ctx, "acl", "filter_v"+version, map[string]string{"ip_version": version})
				}
			}

			var diags interface{ HasError() bool }
			if tt.mode == "datacenter" {
				d, _ := mgr.ExecuteDatacenterOperations(ctx)
				diags = d
			} else {
				d, _ := mgr.ExecuteCampusOperations(ctx)
				diags = d
			}
			if diags.HasError() {
				t.Fatalf("execution returned errors")
			}

			mu.Lock()
			got := append([]string(nil), aclHeaders...)
			mu.Unlock()
			if len(got) != len(tt.want) {
				t.Fatalf("ACL %s request count = %d, want %d (%v)", tt.operation, len(got), len(tt.want), got)
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Fatalf("ACL %s header order = %v, want %v", tt.operation, got, tt.want)
				}
			}
		})
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
	mgr.AddPatch(ctx, "ipv4_list", "test_ipv4_list", zeroPatchValue("ipv4_list"))
	mgr.AddPatch(ctx, "ipv6_list", "test_ipv6_list", zeroPatchValue("ipv6_list"))

	diags, _ := mgr.ExecuteCampusOperations(ctx)
	if diags.HasError() {
		t.Fatalf("ExecuteCampusOperations returned errors: %v", diags)
	}

	patchPaths := filterRecords(snapshotRecords(mu, records), http.MethodPatch)
	assertOrderedSubset(t, "campus PATCH order", patchPaths, toPaths(campusPatchOrder))
	assertPathsAbsent(t, "campus PATCH order", patchPaths, []string{"/ipv4lists", "/ipv6lists"})
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
		return r.Method == http.MethodDelete && r.URL.Path == "/thresholdgroups"
	})

	client := newTestClient(server.URL)
	mgr := bulkops.GetManager(client, nopClearCache, nil, "datacenter")

	ctx := context.Background()

	// threshold_group is first in reverse DC order — its failure should abort the rest
	mgr.AddDelete(ctx, "threshold_group", "test_tg")
	// These come AFTER threshold_group in reverse DC order — should NOT execute
	mgr.AddDelete(ctx, "grouping_rule", "test_gr")
	mgr.AddDelete(ctx, "badge", "test_badge")

	diags, _ := mgr.ExecuteDatacenterOperations(ctx)
	if !diags.HasError() {
		t.Fatal("expected errors from failed threshold_group DELETE, got none")
	}

	deletePaths := filterRecords(snapshotRecords(mu, records), http.MethodDelete)

	for _, path := range deletePaths {
		if path == "/groupingrules" || path == "/badges" {
			t.Errorf("path %q should NOT have been deleted after threshold_group DELETE failure", path)
		}
	}

	foundTG := false
	for _, path := range deletePaths {
		if path == "/thresholdgroups" {
			foundTG = true
		}
	}
	if !foundTG {
		t.Error("threshold_group DELETE should have been attempted")
	}
}
