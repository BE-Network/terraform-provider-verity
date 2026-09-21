package lifecycle

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestGenericMatchesLegacyOnAutoAssignment(t *testing.T) {
	service := func(vlan int, pair string) string {
		return fmt.Sprintf(`resource "verity_service" "test" {
  name = "diffsvc"
  anycast_ipv4_mask = ""
  anycast_ipv6_mask = ""
  dhcp_server_ipv4 = ""
  dhcp_server_ipv6 = ""
  enable = true
  ip_attach_host_advertise = 42
  mtu = 1500
  policy_based_routing = ""
  policy_based_routing_ref_type_ = ""
  tenant = ""
  tenant_ref_type_ = ""
  vlan = %d
%s}
`, vlan, pair)
	}

	scenarios := []struct {
		name           string
		create, update string
		outcome        lifecycleOutcome
	}{
		{
			name:   "assignment on at create",
			create: service(10, "  vni_auto_assigned_ = true\n"),
			update: service(10, "  vni_auto_assigned_ = true\n"),
		},
		{
			name:   "manual vni changes",
			create: service(10, "  vni = 100\n  vni_auto_assigned_ = false\n"),
			update: service(10, "  vni = 200\n  vni_auto_assigned_ = false\n"),
		},
		{
			name:   "assignment is turned on",
			create: service(10, "  vni = 100\n  vni_auto_assigned_ = false\n"),
			update: service(10, "  vni_auto_assigned_ = true\n"),
			outcome: lifecycleOutcome{updatePlanChecks: []plancheck.PlanCheck{
				plancheck.ExpectUnknownValue("verity_service.test", tfjsonpath.New("vni")),
			}},
		},
		{
			name:   "assignment is turned off with a new vni",
			create: service(10, "  vni_auto_assigned_ = true\n"),
			update: service(10, "  vni = 300\n  vni_auto_assigned_ = false\n"),
		},
		{
			name:   "assignment is turned off without a vni",
			create: service(10, "  vni = 100\n  vni_auto_assigned_ = false\n"),
			update: service(10, "  vni_auto_assigned_ = false\n"),
		},
		{
			name:   "vlan changes with assignment on",
			create: service(10, "  vni_auto_assigned_ = true\n"),
			update: service(20, "  vni_auto_assigned_ = true\n"),
			outcome: lifecycleOutcome{updatePlanChecks: []plancheck.PlanCheck{
				plancheck.ExpectUnknownValue("verity_service.test", tfjsonpath.New("vni")),
			}},
		},
		{

			name:   "vlan changes with assignment on and a vni in state",
			create: service(10, "  vni = 100\n  vni_auto_assigned_ = false\n"),
			update: service(20, "  vni_auto_assigned_ = true\n"),
			outcome: lifecycleOutcome{
				intermediate: []string{service(10, "  vni_auto_assigned_ = true\n")},
				updatePlanChecks: []plancheck.PlanCheck{
					plancheck.ExpectUnknownValue("verity_service.test", tfjsonpath.New("vni")),
				},
			},
		},
		{

			name:   "vlan changes with assignment off and vni not written",
			create: service(10, "  vni = 100\n  vni_auto_assigned_ = false\n"),
			update: service(20, "  vni_auto_assigned_ = false\n"),
			outcome: lifecycleOutcome{updatePlanChecks: []plancheck.PlanCheck{
				plancheck.ExpectUnknownValue("verity_service.test", tfjsonpath.New("vni")),
			}},
		},
		{
			name:   "vlan changes with assignment off and vni written",
			create: service(10, "  vni = 100\n  vni_auto_assigned_ = false\n"),
			update: service(20, "  vni = 100\n  vni_auto_assigned_ = false\n"),
			outcome: lifecycleOutcome{updatePlanChecks: []plancheck.PlanCheck{
				plancheck.ExpectKnownValue("verity_service.test", tfjsonpath.New("vni"), knownvalue.Int64Exact(100)),
			}},
		},
		{
			name:    "vni written while assignment is on is refused",
			create:  service(10, "  vni = 100\n  vni_auto_assigned_ = false\n"),
			update:  service(10, "  vni = 100\n  vni_auto_assigned_ = true\n"),
			outcome: lifecycleOutcome{applyError: regexp.MustCompile(`(?s)'vni' field cannot be specified in the configuration when\s+'vni_auto_assigned_' is set to true`)},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			legacy := captureLifecycle(t, "verity_service", false, scenario.create, scenario.update, scenario.outcome)
			generic := captureLifecycle(t, "verity_service", true, scenario.create, scenario.update, scenario.outcome)

			for _, operation := range []string{"PUT", "PATCH"} {
				want, got := canonical(t, legacy[operation]), canonical(t, generic[operation])
				if want != got {
					t.Errorf("%s differs between implementations\n  legacy:  %s\n  generic: %s", operation, want, got)
				}
			}
		})
	}
}

func TestAutoAssignedValueNotWrittenOnCreateIsOmitted(t *testing.T) {
	create := `resource "verity_service" "test" {
  name = "diffsvc"
  anycast_ipv4_mask = ""
  anycast_ipv6_mask = ""
  dhcp_server_ipv4 = ""
  dhcp_server_ipv6 = ""
  enable = true
  ip_attach_host_advertise = 42
  mtu = 1500
  policy_based_routing = ""
  policy_based_routing_ref_type_ = ""
  tenant = ""
  tenant_ref_type_ = ""
  vlan = 10
  vni_auto_assigned_ = false
}
`
	legacy := captureLifecycle(t, "verity_service", false, create, create)
	generic := captureLifecycle(t, "verity_service", true, create, create)

	const common = `"anycast_ipv4_mask":"","anycast_ipv6_mask":"","dhcp_server_ipv4":"","dhcp_server_ipv6":"","enable":true,"ip_attach_host_advertise":42,"mtu":1500,"name":"diffsvc","policy_based_routing":"","policy_based_routing_ref_type_":"","tenant":"","tenant_ref_type_":"","vlan":10`
	wantLegacy := `{"service":{"diffsvc":{` + common + `,"vni":0,"vni_auto_assigned_":false}}}`
	wantGeneric := `{"service":{"diffsvc":{` + common + `,"vni_auto_assigned_":false}}}`

	if got := canonical(t, legacy["PUT"]); got != wantLegacy {
		t.Errorf("legacy PUT = %s\nwant %s\nif the handwritten create now skips an unknown vni, the two agree and this test should assert equality", got, wantLegacy)
	}
	if got := canonical(t, generic["PUT"]); got != wantGeneric {
		t.Errorf("generic PUT = %s\nwant %s\nan unwritten vni must follow unknown_plan", got, wantGeneric)
	}
}

func TestAutoAssignedValueWrittenAsNullIsSent(t *testing.T) {
	config := func(vni string) string {
		return `resource "verity_service" "test" {
  name = "diffsvc"
  anycast_ipv4_mask = ""
  anycast_ipv6_mask = ""
  dhcp_server_ipv4 = ""
  dhcp_server_ipv6 = ""
  enable = true
  ip_attach_host_advertise = 42
  mtu = 1500
  policy_based_routing = ""
  policy_based_routing_ref_type_ = ""
  tenant = ""
  tenant_ref_type_ = ""
  vlan = 10
  vni = ` + vni + `
  vni_auto_assigned_ = false
}
`
	}

	legacy := captureLifecycle(t, "verity_service", false, config("100"), config("null"))
	generic := captureLifecycle(t, "verity_service", true, config("100"), config("null"))

	if got := canonical(t, legacy["PATCH"]); got != "<no request>" {
		t.Errorf("legacy PATCH = %s, want no request: if the handwritten service now sends the null, the two agree and this test should assert equality", got)
	}
	if got, want := canonical(t, generic["PATCH"]), `{"service":{"diffsvc":{"vni":null,"vni_auto_assigned_":false}}}`; got != want {
		t.Errorf("generic PATCH = %s, want %s: a written null must be sent", got, want)
	}
}

func TestGenericMatchesLegacyOnTenantAutoAssignment(t *testing.T) {
	tenant := func(pairs string) string {
		return `resource "verity_tenant" "test" {
  name = "difftenant"
  default_originate = true
  dhcp_relay_source_ipv4s_subnet = ""
  dhcp_relay_source_ipv6s_subnet = ""
  enable = true
  export_route_map = ""
  export_route_map_ref_type_ = ""
  import_route_map = ""
  import_route_map_ref_type_ = ""
  route_aggregation = ""
  route_distinguisher = ""
  route_target_export = ""
  route_target_import = ""
  tenant_type = ""
` + pairs + `}
`
	}
	const manual = "  layer_3_vlan = 42\n  layer_3_vlan_auto_assigned_ = false\n"

	scenarios := []struct {
		name           string
		create, update string
		outcome        lifecycleOutcome
	}{
		{
			name:   "numeric assignment on at create",
			create: tenant(manual + "  layer_3_vni_auto_assigned_ = true\n  vrf_name = \"TestVrf\"\n  vrf_name_auto_assigned_ = false\n"),
			update: tenant(manual + "  layer_3_vni_auto_assigned_ = true\n  vrf_name = \"TestVrf\"\n  vrf_name_auto_assigned_ = false\n"),
		},
		{
			name:   "numeric assignment turned off with a new value",
			create: tenant(manual + "  layer_3_vni_auto_assigned_ = true\n  vrf_name = \"TestVrf\"\n  vrf_name_auto_assigned_ = false\n"),
			update: tenant(manual + "  layer_3_vni = 5000\n  layer_3_vni_auto_assigned_ = false\n  vrf_name = \"TestVrf\"\n  vrf_name_auto_assigned_ = false\n"),
		},
		{
			name:   "string assignment turned on",
			create: tenant(manual + "  layer_3_vni = 5000\n  layer_3_vni_auto_assigned_ = false\n  vrf_name = \"TestVrf\"\n  vrf_name_auto_assigned_ = false\n"),
			update: tenant(manual + "  layer_3_vni = 5000\n  layer_3_vni_auto_assigned_ = false\n  vrf_name_auto_assigned_ = true\n"),
			outcome: lifecycleOutcome{updatePlanChecks: []plancheck.PlanCheck{
				plancheck.ExpectUnknownValue("verity_tenant.test", tfjsonpath.New("vrf_name")),
			}},
		},
		{

			name:    "an empty string written while string assignment is turned on",
			create:  tenant(manual + "  layer_3_vni = 5000\n  layer_3_vni_auto_assigned_ = false\n  vrf_name = \"TestVrf\"\n  vrf_name_auto_assigned_ = false\n"),
			update:  tenant(manual + "  layer_3_vni = 5000\n  layer_3_vni_auto_assigned_ = false\n  vrf_name = \"\"\n  vrf_name_auto_assigned_ = true\n"),
			outcome: lifecycleOutcome{applyError: regexp.MustCompile(`(?s)planned an invalid value\s+for verity_tenant.test.vrf_name`)},
		},
		{
			name:    "a value written while numeric assignment is on is refused",
			create:  tenant(manual + "  layer_3_vni = 5000\n  layer_3_vni_auto_assigned_ = false\n  vrf_name = \"TestVrf\"\n  vrf_name_auto_assigned_ = false\n"),
			update:  tenant(manual + "  layer_3_vni = 5000\n  layer_3_vni_auto_assigned_ = true\n  vrf_name = \"TestVrf\"\n  vrf_name_auto_assigned_ = false\n"),
			outcome: lifecycleOutcome{applyError: regexp.MustCompile(`(?s)'layer_3_vni' field cannot be specified in the configuration when\s+'layer_3_vni_auto_assigned_' is set to true`)},
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			legacy := captureLifecycle(t, "verity_tenant", false, scenario.create, scenario.update, scenario.outcome)
			generic := captureLifecycle(t, "verity_tenant", true, scenario.create, scenario.update, scenario.outcome)
			for _, operation := range []string{"PUT", "PATCH"} {
				want, got := canonical(t, legacy[operation]), canonical(t, generic[operation])
				if want != got {
					t.Errorf("%s differs between implementations\n  legacy:  %s\n  generic: %s", operation, want, got)
				}
			}
		})
	}
}

func TestGenericMatchesLegacyOnFabricAutoAssignment(t *testing.T) {
	entry := coverageEntry(t, "verity_fabric")
	rs := inspectSchema(entry.Factory)
	base := generateCoverageHCL(rs, entry.TerraformType, "difffabric", entry.Mode, entry.modeFieldsKey(), entry.Overrides)
	const value, flag = `  anycast_mac_address = ""` + "\n", `  anycast_mac_address_auto_assigned_ = false` + "\n"
	if !strings.Contains(base, value) || !strings.Contains(base, flag) {
		t.Fatal("the harness configuration no longer writes the anycast pair the way this test expects")
	}
	assignmentOn := strings.Replace(strings.Replace(base, value, "", 1), flag, "  anycast_mac_address_auto_assigned_ = true\n", 1)
	writtenWhileOn := strings.Replace(strings.Replace(base, value, "  anycast_mac_address = \"aa:bb:cc:dd:ee:ff\"\n", 1), flag, "  anycast_mac_address_auto_assigned_ = true\n", 1)

	scenarios := []struct {
		name    string
		update  string
		outcome lifecycleOutcome
	}{
		{"assignment turned on", assignmentOn, lifecycleOutcome{updatePlanChecks: []plancheck.PlanCheck{
			plancheck.ExpectUnknownValue("verity_fabric.test", tfjsonpath.New("anycast_mac_address")),
		}}},
		{"a value written while assignment is on is refused", writtenWhileOn, lifecycleOutcome{
			applyError: regexp.MustCompile(`(?s)'anycast_mac_address'\s+field\s+cannot\s+be\s+specified\s+in\s+the\s+configuration\s+when\s+'anycast_mac_address_auto_assigned_'\s+is\s+set\s+to\s+true`),
		}},
	}
	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			legacy := captureLifecycle(t, entry.TerraformType, false, base, scenario.update, scenario.outcome)
			generic := captureLifecycle(t, entry.TerraformType, true, base, scenario.update, scenario.outcome)
			for _, operation := range []string{"PUT", "PATCH"} {
				want, got := canonical(t, legacy[operation]), canonical(t, generic[operation])
				if want != got {
					t.Errorf("%s differs between implementations\n  legacy:  %s\n  generic: %s", operation, want, got)
				}
			}
		})
	}
}
