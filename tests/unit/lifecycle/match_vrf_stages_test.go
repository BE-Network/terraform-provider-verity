package lifecycle

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"terraform-provider-verity/tests/unit/mock"
)

const routeMapTenantChainWithStages = `
resource "verity_operation_stage" "acl_v4_stage" {
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "route_map_tenant_stage" {
  depends_on = [verity_operation_stage.acl_v4_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_operation_stage" "service_stage" {
  depends_on = [verity_operation_stage.route_map_tenant_stage]
  lifecycle {
    create_before_destroy = true
  }
}

resource "verity_route_map_clause" "chain_filter_clause" {
  name       = "chain_filter_clause"
  depends_on = [verity_operation_stage.route_map_tenant_stage]
}

resource "verity_route_map_clause" "chain_vrf_clause" {
  name                = "chain_vrf_clause"
  match_vrf           = "chain_tenant"
  match_vrf_ref_type_ = "tenant"
  depends_on          = [verity_operation_stage.route_map_tenant_stage]
}

resource "verity_route_map" "chain_filter" {
  name = "chain_filter"
  route_map_clauses {
    index                      = 1
    route_map_clause           = "chain_filter_clause"
    route_map_clause_ref_type_ = "route_map_clause"
  }
  depends_on = [verity_operation_stage.route_map_tenant_stage]
}

resource "verity_route_map" "chain_extra" {
  name = "chain_extra"
  route_map_clauses {
    index                      = 1
    route_map_clause           = "chain_vrf_clause"
    route_map_clause_ref_type_ = "route_map_clause"
  }
  depends_on = [verity_operation_stage.route_map_tenant_stage]
}

resource "verity_tenant" "chain_tenant" {
  name                       = "chain_tenant"
  import_route_map           = "chain_filter"
  import_route_map_ref_type_ = "route_map"
  depends_on                 = [verity_operation_stage.route_map_tenant_stage]
}
`

func newChainServer(t *testing.T) (*mock.MockServer, string) {
	t.Helper()
	ms := mock.NewMockServer("datacenter")
	t.Cleanup(ms.Close)
	ms.SetTestLogger(t)
	if err := ms.LoadResponsesFromDir(mock.ResponsesDir("datacenter")); err != nil {
		t.Fatal(err)
	}
	ms.SetReferenceCheck("/api/routemapclauses", "route_map_clause", "match_vrf", "/api/tenants", "tenant")
	ms.SetReferenceCheck("/api/tenants", "tenant", "import_route_map", "/api/routemaps", "route_map")
	return ms, mock.ProviderConfig(ms.URL(), "datacenter") + routeMapTenantChainWithStages
}

func requestPosition(ms *mock.MockServer, method, path string) int {
	for position, request := range ms.GetRequests() {
		if request.Method == method && request.Path == path {
			return position
		}
	}
	return -1
}

func patchSettingMatchVrf(ms *mock.MockServer, clause, value string) int {
	for position, request := range ms.GetRequests() {
		if request.Method != "PATCH" || request.Path != "/api/routemapclauses" {
			continue
		}
		objects, _ := request.Body["route_map_clause"].(map[string]interface{})
		fields, _ := objects[clause].(map[string]interface{})
		if fields["match_vrf"] == value {
			return position
		}
	}
	return -1
}

func TestRouteMapTenantChainMatchesTheImporterStages(t *testing.T) {
	generated, err := os.ReadFile(filepath.Join("testdata", "importer_golden", "datacenter", "stages.tf"))
	if err != nil {
		t.Fatal(err)
	}
	for _, block := range []string{
		`resource "verity_operation_stage" "route_map_tenant_stage" {
  depends_on = [verity_operation_stage.acl_v4_stage]`,
		`resource "verity_operation_stage" "service_stage" {
  depends_on = [verity_operation_stage.route_map_tenant_stage]`,
	} {
		if !strings.Contains(string(generated), block) || !strings.Contains(routeMapTenantChainWithStages, block) {
			t.Fatalf("the importer's stages.tf and this test's configuration must both contain:\n%s", block)
		}
	}
	for _, file := range []string{"routemapclauses.tf", "policybasedroutingacl.tf", "routemaps.tf", "policybasedrouting.tf", "tenants.tf"} {
		content, err := os.ReadFile(filepath.Join("testdata", "importer_golden", "datacenter", file))
		if err != nil {
			t.Fatal(err)
		}
		if strings.Count(string(content), "depends_on = [verity_operation_stage.") != strings.Count(string(content), "depends_on = [verity_operation_stage.route_map_tenant_stage]") {
			t.Fatalf("%s: every resource must depend on route_map_tenant_stage", file)
		}
	}
}

func TestRouteMapTenantChainIsCreatedAndDestroyedWithStages(t *testing.T) {
	ms, config := newChainServer(t)

	fwresource.UnitTest(t, fwresource.TestCase{
		ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
		Steps: []fwresource.TestStep{
			{
				PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), config) },
				Config:    config,
				Check: fwresource.ComposeTestCheckFunc(
					fwresource.TestCheckResourceAttr("verity_route_map_clause.chain_vrf_clause", "match_vrf", "chain_tenant"),
					fwresource.TestCheckResourceAttr("verity_tenant.chain_tenant", "import_route_map", "chain_filter"),
				),
			},
			{Config: config, PlanOnly: true},
		},
	})

	clausePut := requestPosition(ms, "PUT", "/api/routemapclauses")
	routeMapPut := requestPosition(ms, "PUT", "/api/routemaps")
	tenantPut := requestPosition(ms, "PUT", "/api/tenants")
	restore := patchSettingMatchVrf(ms, "chain_vrf_clause", "chain_tenant")
	if clausePut < 0 || routeMapPut < clausePut || tenantPut < routeMapPut || restore < tenantPut {
		t.Fatalf("positions: clause PUT %d, route map PUT %d, tenant PUT %d, restore PATCH %d; want them in this order", clausePut, routeMapPut, tenantPut, restore)
	}
	if puts := ms.GetRequestsByMethodAndPath("PUT", "/api/routemapclauses"); len(puts) != 1 {
		t.Fatalf("got %d clause PUTs, want one batch for both clauses", len(puts))
	}
	cleared := patchSettingMatchVrf(ms, "chain_vrf_clause", "")
	tenantDelete := requestPosition(ms, "DELETE", "/api/tenants")
	if cleared < restore || tenantDelete < cleared {
		t.Fatalf("positions: clearing PATCH %d, tenant DELETE %d; want match_vrf cleared before the tenant is deleted", cleared, tenantDelete)
	}
}

func TestRouteMapTenantChainRecoversFromAFailedRestoreWithStages(t *testing.T) {
	ms, config := newChainServer(t)

	fwresource.UnitTest(t, fwresource.TestCase{
		ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
		Steps: []fwresource.TestStep{
			{
				PreConfig: func() {
					mock.WriteTFConfig(t, ms.URL(), config)
					ms.FailRequests("PATCH", "/api/routemapclauses", 1, "patch rejected by the test")
				},
				Config:      config,
				ExpectError: regexp.MustCompile(`patch rejected by the test`),
			},
			{
				Config: config,
				Check:  fwresource.TestCheckResourceAttr("verity_route_map_clause.chain_vrf_clause", "match_vrf", "chain_tenant"),
			},
			{Config: config, PlanOnly: true},
		},
	})

	if puts := ms.GetRequestsByMethodAndPath("PUT", "/api/tenants"); len(puts) != 1 {
		t.Fatalf("got %d tenant PUTs, want only the first apply's", len(puts))
	}
}

const routeMapTenantChainWithPerTypeStages = `
resource "verity_operation_stage" "acl_v4_stage" {}

resource "verity_operation_stage" "route_map_clause_stage" {
  depends_on = [verity_operation_stage.acl_v4_stage]
}

resource "verity_operation_stage" "route_map_stage" {
  depends_on = [verity_operation_stage.route_map_clause_stage]
}

resource "verity_operation_stage" "tenant_stage" {
  depends_on = [verity_operation_stage.route_map_stage]
}

resource "verity_route_map_clause" "chain_filter_clause" {
  name       = "chain_filter_clause"
  depends_on = [verity_operation_stage.route_map_clause_stage]
}

resource "verity_route_map_clause" "chain_vrf_clause" {
  name                = "chain_vrf_clause"
  match_vrf           = "chain_tenant"
  match_vrf_ref_type_ = "tenant"
  depends_on          = [verity_operation_stage.route_map_clause_stage]
}

resource "verity_route_map" "chain_filter" {
  name = "chain_filter"
  route_map_clauses {
    index                      = 1
    route_map_clause           = "chain_filter_clause"
    route_map_clause_ref_type_ = "route_map_clause"
  }
  depends_on = [verity_operation_stage.route_map_stage]
}

resource "verity_tenant" "chain_tenant" {
  name                       = "chain_tenant"
  import_route_map           = "chain_filter"
  import_route_map_ref_type_ = "route_map"
  depends_on                 = [verity_operation_stage.tenant_stage]
}
`

func TestRouteMapTenantChainWithTheOldPerTypeStagesFailsClearly(t *testing.T) {
	ms, _ := newChainServer(t)
	config := mock.ProviderConfig(ms.URL(), "datacenter") + routeMapTenantChainWithPerTypeStages

	fwresource.UnitTest(t, fwresource.TestCase{
		ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
		Steps: []fwresource.TestStep{{
			PreConfig:   func() { mock.WriteTFConfig(t, ms.URL(), config) },
			Config:      config,
			ExpectError: regexp.MustCompile(`Cannot resolve object reference`),
		}},
	})

	for _, request := range ms.GetRequestsByMethodAndPath("PUT", "/api/routemapclauses") {
		objects, _ := request.Body["route_map_clause"].(map[string]interface{})
		if fields, sent := objects["chain_vrf_clause"].(map[string]interface{}); sent && fields["match_vrf"] == "" {
			t.Fatal("with the tenant in a later stage, the clause must not be created with a blank reference")
		}
	}
}
