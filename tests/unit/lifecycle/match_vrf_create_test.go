package lifecycle

import (
	"regexp"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"terraform-provider-verity/tests/unit/mock"
)

const matchVrfConfig = `
resource "verity_tenant" "tenant" {
  name = "fix5_tenant"
}

resource "verity_route_map_clause" "clause" {
  name                = "fix5_clause"
  match_vrf           = "fix5_tenant"
  match_vrf_ref_type_ = "tenant"
}
`

func newMatchVrfServer(t *testing.T) (*mock.MockServer, string) {
	t.Helper()
	ms := mock.NewMockServer("datacenter")
	t.Cleanup(ms.Close)
	ms.SetTestLogger(t)
	if err := ms.LoadResponsesFromDir(mock.ResponsesDir("datacenter")); err != nil {
		t.Fatal(err)
	}
	ms.SetReferenceCheck("/api/routemapclauses", "route_map_clause", "match_vrf", "/api/tenants", "tenant")
	return ms, mock.ProviderConfig(ms.URL(), "datacenter") + matchVrfConfig
}

func matchVrfIn(request mock.CapturedRequest) interface{} {
	objects, _ := request.Body["route_map_clause"].(map[string]interface{})
	clause, _ := objects["fix5_clause"].(map[string]interface{})
	return clause["match_vrf"]
}

func patchesSetting(ms *mock.MockServer, value string) int {
	count := 0
	for _, request := range ms.GetRequestsByMethodAndPath("PATCH", "/api/routemapclauses") {
		if matchVrfIn(request) == value {
			count++
		}
	}
	return count
}

func TestClauseAndTenantCreatedTogether(t *testing.T) {
	ms, config := newMatchVrfServer(t)

	fwresource.UnitTest(t, fwresource.TestCase{
		ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
		Steps: []fwresource.TestStep{
			{
				PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), config) },
				Config:    config,
				Check:     fwresource.TestCheckResourceAttr("verity_route_map_clause.clause", "match_vrf", "fix5_tenant"),
			},
			{Config: config, PlanOnly: true},
		},
	})

	puts := ms.GetRequestsByMethodAndPath("PUT", "/api/routemapclauses")
	if len(puts) != 1 || matchVrfIn(puts[0]) != "" {
		t.Fatalf("clause PUTs = %d, want one without match_vrf", len(puts))
	}
	if got := patchesSetting(ms, "fix5_tenant"); got != 1 {
		t.Fatalf("clause PATCHes setting match_vrf = %d, want one", got)
	}
}

func TestFailedMatchVrfRestoreIsRepairedByTheNextApply(t *testing.T) {
	ms, config := newMatchVrfServer(t)

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
				Check: fwresource.ComposeTestCheckFunc(
					fwresource.TestCheckResourceAttr("verity_tenant.tenant", "name", "fix5_tenant"),
					fwresource.TestCheckResourceAttr("verity_route_map_clause.clause", "match_vrf", "fix5_tenant"),
				),
			},
			{Config: config, PlanOnly: true},
		},
	})

	if puts := ms.GetRequestsByMethodAndPath("PUT", "/api/tenants"); len(puts) != 1 {
		t.Fatalf("tenant PUTs = %d, want only the first apply's", len(puts))
	}
	if puts := ms.GetRequestsByMethodAndPath("PUT", "/api/routemapclauses"); len(puts) != 1 {
		t.Fatalf("clause PUTs = %d, want only the first apply's", len(puts))
	}
	if got := patchesSetting(ms, "fix5_tenant"); got != 2 {
		t.Fatalf("clause PATCHes setting match_vrf = %d, want the failed one and the repair", got)
	}
}
