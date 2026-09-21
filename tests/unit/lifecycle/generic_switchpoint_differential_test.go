package lifecycle

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

// switchpointPair is one of verity_switchpoint's ten auto-assignment pairs, with
// two distinct values for it.
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

// switchpointBase is the harness configuration with every pair attribute taken
// out, so each case writes exactly the pair lines it means to.
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

// withPairs writes the lines every pair gets into the base configuration.
func withPairs(base string, lines func(switchpointPair) string) string {
	var body strings.Builder
	for _, pair := range switchpointPairs {
		body.WriteString(lines(pair))
	}
	return strings.Replace(base, "  name = \"diffsp\"\n", "  name = \"diffsp\"\n"+body.String(), 1)
}

// The handwritten verity_switchpoint writes its ten auto-assignment pairs in two
// code shapes — three read the configuration locally, seven reuse an earlier
// read — which is structure, not policy. An audit read all ten as the one
// effective contract the engine implements for every pair:
//
//   - the flag is sent only when the configuration states it;
//   - turning assignment off resends the value, from the plan when it is known
//     and from state otherwise;
//   - a value change with the flag unchanged includes the flag.
//
// These cases write every pair the same way in one configuration, run both
// implementations, and compare each pair's keys as well as the whole body, so a
// pair that departs from the contract is named rather than lost in a large diff.
//
// One transition departs, and only for five pairs: turning assignment off
// without writing the value. The value is unwritten, Optional and Computed, so
// it plans as unknown. For bgp_as_number, switch_vtep_id_ip_mask,
// switch_router_id_ip_mask, controller_ip_and_mask, and switch_ip_and_mask the
// handwritten resend tests only !plan.X.IsNull(), which an unknown passes, and
// ValueInt64 or ValueString reports the zero value for it, so the PATCH clears
// the value it means to keep: "bgp_as_number": 0, "controller_ip_and_mask": "".
// The other five also test !plan.X.IsUnknown() and fall back to state, as the
// engine does for all ten. It is the unknown-as-zero defect recorded for the
// service's vni. The engine's behavior is accepted as a bug fix, and no registry
// exception models the defect. The case pins both exact bodies, pair by pair.
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
		// legacyZeroes names the pairs whose PATCH the handwritten resource sends
		// with the value's zero where the engine resends state's value.
		legacyZeroes map[string]interface{}
	}{
		// The flag is sent only when configured.
		{name: "values written without flags", create: valuesOnly(first), update: valuesOnly(second)},
		{name: "flags removed while values change", create: manual(first), update: valuesOnly(second)},
		// A value change with the flag unchanged includes the flag.
		{name: "values change with assignment off", create: manual(first), update: manual(second)},
		{name: "assignment is turned on", create: manual(first), update: assigned},
		// Turning assignment off resends the value: from the plan when written…
		{name: "assignment is turned off with new values", create: assigned, update: manual(second)},
		// …and otherwise from what state holds. Assignment is turned on in between
		// so the mock keeps the values the create stored.
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

// object_properties.number_of_multipoints is the provider's one nullable member
// of a singleton block. The handwritten resource sends it only when the .tf
// files write it, a written null included, updates it only when both the plan
// and state hold the block, and plans a written null over a state value as null
// so the update sends it.
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

// switchpointKey returns one top-level key of the diffsp object in a captured
// body, and whether the body holds it.
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

// unknownPairValues requires every pair's value to plan as unknown: unwritten,
// Optional and Computed, on a resource that changes.
func unknownPairValues() []plancheck.PlanCheck {
	checks := make([]plancheck.PlanCheck, 0, len(switchpointPairs))
	for _, pair := range switchpointPairs {
		checks = append(checks, plancheck.ExpectUnknownValue("verity_switchpoint.test", tfjsonpath.New(pair.value)))
	}
	return checks
}

// withoutKeys returns a copy of a captured body with the given keys of the
// diffsp object removed, so the rest of the body is still compared exactly.
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
