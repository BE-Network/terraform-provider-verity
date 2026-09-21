package lifecycle

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"terraform-provider-verity/internal/provider"
	"terraform-provider-verity/tests/unit/mock"
)

func TestGenericMatchesLegacyOnPairAndNullableUpdates(t *testing.T) {
	const terraformType = "verity_diagnostics_profile"

	base := func(body string) string {
		return fmt.Sprintf("resource %q \"test\" {\n  name = \"diffp\"\n%s}\n", terraformType, body)
	}

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

			update: fixed + `  flow_collector = "collector-a"
  flow_collector_ref_type_ = "sflow_collector"
`,
		},
		{

			name: "nullable numeric is written as null on create",
			create: fixed + `  poll_interval = null
  flow_collector = "collector-a"
  flow_collector_ref_type_ = "sflow_collector"
`,
			update: fixed + `  poll_interval = null
  flow_collector = "collector-a"
  flow_collector_ref_type_ = "sflow_collector"
`,
		},
		{
			name: "nullable numeric is absent on create",
			create: fixed + `  flow_collector = "collector-a"
  flow_collector_ref_type_ = "sflow_collector"
`,
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

type lifecycleOutcome struct {
	driftAfterApply bool

	applyError *regexp.Regexp

	updatePlanChecks []plancheck.PlanCheck

	intermediate []string
}

func captureLifecycle(t *testing.T, terraformType string, generic bool, createConfig, updateConfig string, outcome ...lifecycleOutcome) map[string]map[string]interface{} {
	t.Helper()

	var expect lifecycleOutcome
	if len(outcome) > 0 {
		expect = outcome[0]
	}

	selection := ""
	if generic {
		selection = terraformType
	}
	t.Setenv(provider.GenericResourcesEnvVar, selection)
	if generic {
		assertServedGenerically(t, terraformType)
	}

	entry := coverageEntry(t, terraformType)
	t.Logf("=== lifecycle for %s, generic=%v", terraformType, generic)
	ms := mock.NewMockServer(entry.Mode)
	defer ms.Close()
	ms.SetTestLogger(t)
	if err := ms.LoadResponsesFromDir(mock.ResponsesDir(entry.Mode)); err != nil {
		t.Fatalf("failed to load responses: %v", err)
	}

	provider := mock.ProviderConfig(ms.URL(), entry.Mode)
	update := fwresource.TestStep{
		PreConfig:   func() { mock.WriteTFConfig(t, ms.URL(), provider+updateConfig) },
		Config:      provider + updateConfig,
		ExpectError: expect.applyError,
	}
	if len(expect.updatePlanChecks) != 0 {
		update.ConfigPlanChecks.PreApply = expect.updatePlanChecks
	}
	if expect.driftAfterApply {

		update.ExpectNonEmptyPlan = true
		update.ConfigPlanChecks.PostApplyPostRefresh = []plancheck.PlanCheck{plancheck.ExpectNonEmptyPlan()}
	}

	fwresource.UnitTest(t, fwresource.TestCase{
		ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
		Steps:                    lifecycleSteps(t, ms.URL(), provider, createConfig, expect.intermediate, update),
	})

	captured := map[string]map[string]interface{}{}
	for _, method := range []string{"PUT", "PATCH", "DELETE"} {
		requests := ms.GetRequestsByMethodAndPath(method, entry.APIPath)
		if len(requests) == 0 {
			continue
		}
		last := requests[len(requests)-1]
		if method != "DELETE" {
			captured[method] = last.Body
		}

		query := make(map[string]interface{}, len(last.QueryParams))
		for name, values := range last.QueryParams {
			query[name] = values
		}
		captured[method+" query"] = query
	}
	return captured
}

func lifecycleSteps(t *testing.T, url, provider, createConfig string, intermediate []string, update fwresource.TestStep) []fwresource.TestStep {
	t.Helper()
	steps := []fwresource.TestStep{{
		PreConfig: func() { mock.WriteTFConfig(t, url, provider+createConfig) },
		Config:    provider + createConfig,
	}}
	for _, config := range intermediate {
		config := config
		steps = append(steps, fwresource.TestStep{
			PreConfig: func() { mock.WriteTFConfig(t, url, provider+config) },
			Config:    provider + config,
		})
	}
	return append(steps, update)
}

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
