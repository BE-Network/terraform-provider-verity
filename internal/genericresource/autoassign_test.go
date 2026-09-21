package genericresource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-verity/internal/spec"
	"terraform-provider-verity/internal/transport"
)

func autoFields() []spec.FieldSpec {
	base := nullableFields()
	vlan := base[1]
	vlan.TerraformName, vlan.APIName = "vlan", "vlan"
	vni := base[1]
	vni.TerraformName, vni.APIName = "vni", "vni"
	vni.AutoAssignment = &spec.AutoAssignmentSpec{FlagField: "vni_auto_assigned_", RecomputedWhen: []string{"vlan"}}
	flag := spec.FieldSpec{
		TerraformName: "vni_auto_assigned_", APIName: "vni_auto_assigned_", Kind: spec.FieldKindBool,
		Access: spec.AccessOptionalComputed, Modes: []spec.Mode{spec.ModeDatacenter},
		CreateNull: spec.CreateNullOmit, UpdateClear: spec.UpdateClearFalse,
		UnknownPlan: spec.UnknownPlanOmitAndRead, ResponseAbsence: spec.ResponseAbsenceTerraformNull,
	}
	return []spec.FieldSpec{base[0], vlan, vni, flag}
}

func autoPair(t *testing.T) autoAssignedPair {
	t.Helper()
	pairs, _, err := autoAssignmentPairs(autoFields())
	if err != nil || len(pairs) != 1 {
		t.Fatalf("autoAssignmentPairs = %v, %v", pairs, err)
	}
	return pairs[0]
}

func writtenVNI(value attr.Value) nullableSource {
	return nullableSource{
		config:     map[string]attr.Value{"vni": value},
		configured: func(name string) bool { return name == "vni" },
	}
}

func TestCreateAutoAssigned(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name      string
		vni, flag attr.Value
		source    nullableSource
		want      string
	}{
		{"flag on sends the flag and not the value", types.Int64Unknown(), types.BoolValue(true), nullableSource{}, `{"vni_auto_assigned_":true}`},
		{"a written value travels with the flag", types.Int64Value(100), types.BoolValue(false), writtenVNI(types.Int64Value(100)), `{"vni":100,"vni_auto_assigned_":false}`},
		{"an unwritten value is left to the server", types.Int64Unknown(), types.BoolValue(false), nullableSource{configured: func(string) bool { return false }}, `{"vni_auto_assigned_":false}`},
		{"a written null is sent", types.Int64Unknown(), types.BoolValue(false), writtenVNI(types.Int64Null()), `{"vni":null,"vni_auto_assigned_":false}`},
		{"an unknown flag is not sent", types.Int64Unknown(), types.BoolUnknown(), nullableSource{}, `{}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			object := transport.WireObject{}
			plan := map[string]attr.Value{"vni": tc.vni, "vni_auto_assigned_": tc.flag}
			if err := createAutoAssigned(autoPair(t), plan, tc.source, object); err != nil {
				t.Fatalf("createAutoAssigned: %v", err)
			}
			if got := encode(t, object); got != tc.want {
				t.Fatalf("create sent %s, want %s", got, tc.want)
			}
		})
	}
}

func TestUpdateAutoAssigned(t *testing.T) {
	t.Parallel()

	flagWritten := map[string]attr.Value{"vni_auto_assigned_": types.BoolValue(false)}
	flagNotWritten := map[string]attr.Value{"vni_auto_assigned_": types.BoolNull()}

	cases := []struct {
		name                string
		planVNI, planFlag   attr.Value
		stateVNI, stateFlag attr.Value
		config              map[string]attr.Value
		want                string
		wantChanged         bool
	}{
		{"a changed value carries the planned flag",
			types.Int64Value(200), types.BoolValue(false), types.Int64Value(100), types.BoolValue(false), flagWritten,
			`{"vni":200,"vni_auto_assigned_":false}`, true},
		{"turning assignment off resends the value from state",
			types.Int64Unknown(), types.BoolValue(false), types.Int64Value(4242), types.BoolValue(true), flagWritten,
			`{"vni":4242,"vni_auto_assigned_":false}`, true},
		{"a flag change the configuration does not state is not sent",
			types.Int64Value(100), types.BoolValue(true), types.Int64Value(100), types.BoolValue(false), flagNotWritten,
			`{}`, true},
		{"an unknown value alone is no change",
			types.Int64Unknown(), types.BoolValue(true), types.Int64Value(100), types.BoolValue(true), flagWritten,
			`{}`, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			object := transport.WireObject{}
			plan := map[string]attr.Value{"vni": tc.planVNI, "vni_auto_assigned_": tc.planFlag}
			state := map[string]attr.Value{"vni": tc.stateVNI, "vni_auto_assigned_": tc.stateFlag}
			changed, err := updateAutoAssigned(autoPair(t), plan, state, tc.config, object)
			if err != nil {
				t.Fatalf("updateAutoAssigned: %v", err)
			}
			if changed != tc.wantChanged {
				t.Fatalf("changed = %v, want %v", changed, tc.wantChanged)
			}
			if got := encode(t, object); got != tc.want {
				t.Fatalf("update sent %s, want %s", got, tc.want)
			}
		})
	}
}

func TestChangedTriggerAndIsSet(t *testing.T) {
	t.Parallel()

	pair := autoPair(t)
	same := map[string]attr.Value{"vlan": types.Int64Value(10)}
	if got := changedTrigger(pair, same, same); got != "" {
		t.Errorf("changedTrigger with an unchanged vlan = %q, want none", got)
	}
	if got := changedTrigger(pair, map[string]attr.Value{"vlan": types.Int64Value(20)}, same); got != "vlan" {
		t.Errorf("changedTrigger with a changed vlan = %q, want vlan", got)
	}

	if isSet(types.StringValue("")) || isSet(types.Int64Unknown()) || isSet(types.Int64Null()) {
		t.Error("an empty, unknown, or null value was treated as set")
	}
	if !isSet(types.StringValue("x")) || !isSet(types.Int64Value(0)) {
		t.Error("a written value was treated as not set")
	}
}

func TestAutoAssignmentPairsRefusesAMalformedFlag(t *testing.T) {
	t.Parallel()

	fields := autoFields()
	fields = fields[:3]
	if _, _, err := autoAssignmentPairs(fields); err == nil {
		t.Fatal("an auto-assigned value with no flag field was accepted")
	}
}
