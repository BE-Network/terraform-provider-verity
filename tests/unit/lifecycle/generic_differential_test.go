package lifecycle

import (
	"encoding/json"
	"fmt"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"terraform-provider-verity/internal/provider"
	"terraform-provider-verity/tests/unit/mock"
)

// The golden fixtures pin a create and an enable flip, which leaves the update
// paths that matter most to Phase 3 uncovered: a reference pair changing, and a
// nullable numeric changing or being cleared. No fixture records those, so there
// is nothing to compare a migrated resource against.
//
// This runs both implementations over the same configuration against the same
// mock responses and compares what each sends, which is what the plan asks for
// and what a fixture cannot give when no fixture exists. The legacy run is the
// definition of correct; the generic run has to match it.
//
// It also answers a question reasoning alone does not settle. The handwritten
// resources read nullable values from req.Config together with a flag parsed out
// of the .tf file, while the engine reads req.Plan. Whether those agree is a
// fact about Terraform's behavior, and this measures it rather than assuming it.
//
// Every optional field is written in both configurations on purpose. An
// attribute left out plans as unknown, and the two implementations genuinely
// disagree about an unknown on update: the legacy compare helpers serialize one
// as its zero value, and the engine omits it. That difference is deliberate and
// is covered on its own by TestLegacyUpdateSerializesUnknownsAndTheEngineOmits,
// so leaving fields out here would drown the semantics under test in it.
func TestGenericMatchesLegacyOnPairAndNullableUpdates(t *testing.T) {
	const terraformType = "verity_diagnostics_profile"

	base := func(body string) string {
		return fmt.Sprintf("resource %q \"test\" {\n  name = \"diffp\"\n%s}\n", terraformType, body)
	}

	// Written out in full so no attribute plans as unknown; see the note above.
	const fixed = `  enable = true
  enable_sflow = true
  erspan_destination_ip = "10.0.0.9"
  erspan_gre_type = "0x88be"
  erspan_dscp = 10
  erspan_queue = 2
  erspan_ttl = 64
  monitoring_acl = "acl-a"
  monitoring_acl_ref_type_ = "monitoring_acl"
  use_internal_collector = false
  vrf_type = "default"
`

	scenarios := []struct {
		name           string
		create, update string
	}{
		{
			name: "reference value changes and the type does not",
			create: fixed + `  poll_interval = 30
  flow_collector = "collector-a"
  flow_collector_ref_type_ = "sflow_collector"
`,
			update: fixed + `  poll_interval = 30
  flow_collector = "collector-b"
  flow_collector_ref_type_ = "sflow_collector"
`,
		},
		{
			name: "reference pair is cleared",
			create: fixed + `  poll_interval = 30
  flow_collector = "collector-a"
  flow_collector_ref_type_ = "sflow_collector"
`,
			update: fixed + `  poll_interval = 30
  flow_collector = ""
  flow_collector_ref_type_ = ""
`,
		},
		{
			name: "nullable numeric changes",
			create: fixed + `  poll_interval = 30
  flow_collector = "collector-a"
  flow_collector_ref_type_ = "sflow_collector"
`,
			update: fixed + `  poll_interval = 60
  flow_collector = "collector-a"
  flow_collector_ref_type_ = "sflow_collector"
`,
		},
		{
			name: "nullable numeric is removed from configuration",
			create: fixed + `  poll_interval = 30
  flow_collector = "collector-a"
  flow_collector_ref_type_ = "sflow_collector"
`,
			// Removing the attribute is not the same as clearing it: neither
			// implementation sends anything, because an unwritten nullable is left
			// to the server rather than reset.
			update: fixed + `  flow_collector = "collector-a"
  flow_collector_ref_type_ = "sflow_collector"
`,
		},
		{
			name: "nullable numeric is written as null",
			create: fixed + `  poll_interval = 30
  flow_collector = "collector-a"
  flow_collector_ref_type_ = "sflow_collector"
`,
			update: fixed + `  poll_interval = null
  flow_collector = "collector-a"
  flow_collector_ref_type_ = "sflow_collector"
`,
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			legacy := captureLifecycle(t, terraformType, false, base(scenario.create), base(scenario.update))
			generic := captureLifecycle(t, terraformType, true, base(scenario.create), base(scenario.update))

			for _, operation := range []string{"PUT", "PATCH"} {
				want, got := canonical(t, legacy[operation]), canonical(t, generic[operation])
				if want != got {
					t.Errorf("%s differs between implementations\n  legacy:  %s\n  generic: %s", operation, want, got)
				}
			}
		})
	}
}

// captureLifecycle applies a create then an update and returns the last body of
// each method. The switch is set per run, so the two runs differ in nothing but
// which implementation serves the resource.
func captureLifecycle(t *testing.T, terraformType string, generic bool, createConfig, updateConfig string) map[string]map[string]interface{} {
	t.Helper()

	selection := ""
	if generic {
		selection = terraformType
	}
	t.Setenv(provider.GenericResourcesEnvVar, selection)
	if generic {
		assertServedGenerically(t, terraformType)
	}

	entry := coverageEntry(t, terraformType)
	ms := mock.NewMockServer(entry.Mode)
	defer ms.Close()
	ms.SetTestLogger(t)
	if err := ms.LoadResponsesFromDir(mock.ResponsesDir(entry.Mode)); err != nil {
		t.Fatalf("failed to load responses: %v", err)
	}

	provider := mock.ProviderConfig(ms.URL(), entry.Mode)
	captured := map[string]map[string]interface{}{}

	fwresource.UnitTest(t, fwresource.TestCase{
		ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
		Steps: []fwresource.TestStep{
			{
				PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), provider+createConfig) },
				Config:    provider + createConfig,
			},
			{
				PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), provider+updateConfig) },
				Config:    provider + updateConfig,
				Check: func(*terraform.State) error {
					for _, method := range []string{"PUT", "PATCH"} {
						requests := ms.GetRequestsByMethodAndPath(method, entry.APIPath)
						if len(requests) > 0 {
							captured[method] = requests[len(requests)-1].Body
						}
					}
					return nil
				},
			},
		},
	})
	return captured
}

// canonical renders a body so a difference in key order is not read as a
// difference in what was sent.
func canonical(t *testing.T, body map[string]interface{}) string {
	t.Helper()
	if body == nil {
		return "<no request>"
	}
	encoded, err := json.Marshal(body)
	if err != nil {
		t.Fatal(err)
	}
	return string(encoded)
}
