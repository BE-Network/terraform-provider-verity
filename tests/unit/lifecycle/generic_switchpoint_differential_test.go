package lifecycle

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

type switchpointPair struct {
	value, first, second string
}

var switchpointPairs = []switchpointPair{
	{"bgp_as_number", "65001", "65002"},
	{"switch_vtep_id_ip_mask", `"10.1.0.1/32"`, `"10.1.0.2/32"`},
	{"switch_router_id_ip_mask", `"10.2.0.1/32"`, `"10.2.0.2/32"`},
	{"controller_ip_and_mask", `"10.3.0.1/24"`, `"10.3.0.2/24"`},
	{"switch_ip_and_mask", `"10.4.0.1/24"`, `"10.4.0.2/24"`},
	{"ssh_key_or_password_encrypted", `"secret-one"`, `"secret-two"`},
	{"gateway", `"10.5.0.1"`, `"10.5.0.2"`},
	{"switch_gateway", `"10.6.0.1"`, `"10.6.0.2"`},
	{"lldp_search_string", `"uplink-one"`, `"uplink-two"`},
	{"username", `"operator-one"`, `"operator-two"`},
}

func (p switchpointPair) flag() string { return p.value + "_auto_assigned_" }

func switchpointBase(t *testing.T) string {
	t.Helper()
	entry := coverageEntry(t, "verity_switchpoint")
	rs := inspectSchema(entry.Factory)
	base := generateCoverageHCL(rs, entry.TerraformType, "diffsp", entry.Mode, entry.modeFieldsKey(), entry.Overrides)
	for _, pair := range switchpointPairs {
		line := regexp.MustCompile(`(?m)^  (` + pair.value + `|` + pair.flag() + `) = .*\n`)
		if len(line.FindAllString(base, -1)) != 2 {
			t.Fatalf("the harness configuration no longer writes %s and its flag the way this test expects", pair.value)
		}
		base = line.ReplaceAllString(base, "")
	}
	if !strings.Contains(base, "  name = \"diffsp\"\n") {
		t.Fatal("the harness configuration no longer writes the name the way this test expects")
	}
	return base
}

func withPairs(base string, lines func(switchpointPair) string) string {
	var body strings.Builder
	for _, pair := range switchpointPairs {
		body.WriteString(lines(pair))
	}
	return strings.Replace(base, "  name = \"diffsp\"\n", "  name = \"diffsp\"\n"+body.String(), 1)
}

func TestGenericMatchesLegacyOnSwitchpointAutoAssignment(t *testing.T) {
	base := switchpointBase(t)
	manual := func(value func(switchpointPair) string) string {
		return withPairs(base, func(p switchpointPair) string {
			return fmt.Sprintf("  %s = %s\n  %s = false\n", p.value, value(p), p.flag())
		})
	}
	first := func(p switchpointPair) string { return p.first }
	second := func(p switchpointPair) string { return p.second }
	assigned := withPairs(base, func(p switchpointPair) string { return fmt.Sprintf("  %s = true\n", p.flag()) })
	offWithoutValue := withPairs(base, func(p switchpointPair) string { return fmt.Sprintf("  %s = false\n", p.flag()) })
	valuesOnly := func(value func(switchpointPair) string) string {
		return withPairs(base, func(p switchpointPair) string { return fmt.Sprintf("  %s = %s\n", p.value, value(p)) })
	}

	scenarios := []struct {
		name           string
		create, update string
		outcome        lifecycleOutcome

		legacyZeroes map[string]interface{}
	}{

		{name: "values written without flags", create: valuesOnly(first), update: valuesOnly(second)},
		{name: "flags removed while values change", create: manual(first), update: valuesOnly(second)},

		{name: "values change with assignment off", create: manual(first), update: manual(second)},
		{name: "assignment is turned on", create: manual(first), update: assigned},

		{name: "assignment is turned off with new values", create: assigned, update: manual(second)},

		{name: "assignment is turned off without values", create: manual(first), update: offWithoutValue,
			outcome: lifecycleOutcome{intermediate: []string{assigned}, updatePlanChecks: unknownPairValues()},
			legacyZeroes: map[string]interface{}{
				"bgp_as_number": float64(0), "switch_vtep_id_ip_mask": "", "switch_router_id_ip_mask": "",
				"controller_ip_and_mask": "", "switch_ip_and_mask": "",
			}},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			legacy := captureLifecycle(t, "verity_switchpoint", false, scenario.create, scenario.update, scenario.outcome)
			generic := captureLifecycle(t, "verity_switchpoint", true, scenario.create, scenario.update, scenario.outcome)
			for _, operation := range []string{"PUT", "PATCH"} {
				pinned := map[string]bool{}
				for _, pair := range switchpointPairs {
					for _, key := range []string{pair.value, pair.flag()} {
						want, wantHeld := switchpointKey(legacy[operation], key)
						got, gotHeld := switchpointKey(generic[operation], key)
						if zero, known := scenario.legacyZeroes[key]; known && operation == "PATCH" {
							pinned[key] = true
							state := strings.Trim(pair.first, `"`)
							if !wantHeld || fmt.Sprint(want) != fmt.Sprint(zero) {
								t.Errorf("legacy PATCH %s = %s, want the zero value %#v; if the handwritten resend now skips an unknown plan, the two agree and this pin should go", key, describeKey(want, wantHeld), zero)
							}
							if !gotHeld || fmt.Sprint(got) != state {
								t.Errorf("generic PATCH %s = %s, want state's %s", key, describeKey(got, gotHeld), state)
							}
							continue
						}
						if wantHeld != gotHeld || fmt.Sprint(want) != fmt.Sprint(got) {
							t.Errorf("%s %s differs: legacy %s, generic %s", operation, key, describeKey(want, wantHeld), describeKey(got, gotHeld))
						}
					}
				}
				if want, got := canonical(t, withoutKeys(legacy[operation], pinned)), canonical(t, withoutKeys(generic[operation], pinned)); want != got {
					t.Errorf("%s differs between implementations\n  legacy:  %s\n  generic: %s", operation, want, got)
				}
			}
		})
	}
}

func TestGenericMatchesLegacyOnSwitchpointNullableMember(t *testing.T) {
	base := switchpointBase(t)
	member := regexp.MustCompile(`(?m)^    number_of_multipoints = 42\n`)
	if !member.MatchString(base) {
		t.Fatal("the harness configuration no longer writes number_of_multipoints the way this test expects")
	}
	multipoints := func(value string) string {
		if value == "" {
			return member.ReplaceAllString(base, "")
		}
		return member.ReplaceAllLiteralString(base, "    number_of_multipoints = "+value+"\n")
	}
	enableOff := func(config string) string {
		return strings.Replace(config, "  enable = true\n", "  enable = false\n", 1)
	}

	scenarios := []struct {
		name           string
		create, update string
	}{
		{name: "written as null on create, then given a value", create: multipoints("null"), update: multipoints("7")},
		{name: "a value is changed", create: multipoints("5"), update: multipoints("7")},
		{name: "a value is cleared with a written null", create: multipoints("5"), update: multipoints("null")},
		{name: "a value is no longer written", create: multipoints("5"), update: enableOff(multipoints(""))},
		{name: "never written", create: multipoints(""), update: enableOff(multipoints(""))},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			legacy := captureLifecycle(t, "verity_switchpoint", false, scenario.create, scenario.update)
			generic := captureLifecycle(t, "verity_switchpoint", true, scenario.create, scenario.update)
			for _, operation := range []string{"PUT", "PATCH"} {
				if want, got := canonical(t, legacy[operation]), canonical(t, generic[operation]); want != got {
					t.Errorf("%s differs between implementations\n  legacy:  %s\n  generic: %s", operation, want, got)
				}
			}
		})
	}
}

func switchpointKey(body map[string]interface{}, key string) (interface{}, bool) {
	wrapper, _ := body["switchpoint"].(map[string]interface{})
	object, _ := wrapper["diffsp"].(map[string]interface{})
	value, held := object[key]
	return value, held
}

func describeKey(value interface{}, held bool) string {
	if !held {
		return "absent"
	}
	return fmt.Sprintf("%#v", value)
}

func unknownPairValues() []plancheck.PlanCheck {
	checks := make([]plancheck.PlanCheck, 0, len(switchpointPairs))
	for _, pair := range switchpointPairs {
		checks = append(checks, plancheck.ExpectUnknownValue("verity_switchpoint.test", tfjsonpath.New(pair.value)))
	}
	return checks
}

func withoutKeys(body map[string]interface{}, keys map[string]bool) map[string]interface{} {
	if len(keys) == 0 || body == nil {
		return body
	}
	wrapper, _ := body["switchpoint"].(map[string]interface{})
	object, _ := wrapper["diffsp"].(map[string]interface{})
	trimmed := make(map[string]interface{}, len(object))
	for key, value := range object {
		if !keys[key] {
			trimmed[key] = value
		}
	}
	return map[string]interface{}{"switchpoint": map[string]interface{}{"diffsp": trimmed}}
}
