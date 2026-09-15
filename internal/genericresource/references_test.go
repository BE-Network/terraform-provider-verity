package genericresource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-verity/internal/spec"
)

// pairFields models one reference: a value naming an object and a companion
// naming its type, as verity_pair carries three of.
func pairFields(allowed ...string) []spec.FieldSpec {
	if len(allowed) == 0 {
		allowed = []string{"lag"}
	}
	return []spec.FieldSpec{
		{
			TerraformName: "name", APIName: "name", Kind: spec.FieldKindString,
			Access: spec.AccessRequired, Modes: []spec.Mode{spec.ModeDatacenter},
			CreateNull: spec.CreateNullReject, UpdateClear: spec.UpdateClearReject,
			UnknownPlan: spec.UnknownPlanReject, ResponseAbsence: spec.ResponseAbsenceTerraformNull,
		},
		{
			TerraformName: "lag", APIName: "lag", Kind: spec.FieldKindString,
			Access: spec.AccessOptionalComputed, Modes: []spec.Mode{spec.ModeDatacenter},
			CreateNull: spec.CreateNullOmit, UpdateClear: spec.UpdateClearEmptyString,
			UnknownPlan: spec.UnknownPlanOmitAndRead, ResponseAbsence: spec.ResponseAbsenceTerraformNull,
			Reference: &spec.ReferenceSpec{TypeField: "lag_ref_type_", AllowedTypes: allowed},
		},
		{
			TerraformName: "lag_ref_type_", APIName: "lag_ref_type_", Kind: spec.FieldKindString,
			Access: spec.AccessOptionalComputed, Modes: []spec.Mode{spec.ModeDatacenter},
			CreateNull: spec.CreateNullOmit, UpdateClear: spec.UpdateClearEmptyString,
			UnknownPlan: spec.UnknownPlanOmitAndRead, ResponseAbsence: spec.ResponseAbsenceTerraformNull,
		},
	}
}

func pairState(base, refType string) map[string]attr.Value {
	return map[string]attr.Value{
		"name":          types.StringValue("p"),
		"lag":           types.StringValue(base),
		"lag_ref_type_": types.StringValue(refType),
	}
}

func updatePair(t *testing.T, fields []spec.FieldSpec, plan, state map[string]attr.Value) (string, bool, diag.Diagnostics) {
	t.Helper()
	var diagnostics diag.Diagnostics
	object, changed, err := buildUpdate(fields, plan, state, nullableSource{}, &diagnostics)
	if err != nil {
		t.Fatalf("buildUpdate: %v", err)
	}
	if diagnostics.HasError() {
		return "", changed, diagnostics
	}
	return encode(t, object), changed, diagnostics
}

// An untouched pair is not a change, so an update that only alters something
// else does not restate it.
func TestReferencePairUnchangedSendsNothing(t *testing.T) {
	t.Parallel()

	state := pairState("lag-a", "lag")
	body, changed, _ := updatePair(t, pairFields(), state, state)
	if changed {
		t.Fatalf("an untouched pair reported a change and would send %s", body)
	}
}

// With one permitted type the companion is implied by the value, so changing
// the value alone does not restate it. This is the handwritten behavior.
func TestReferencePairValueOnlyChangeSendsValueOnly(t *testing.T) {
	t.Parallel()

	state := pairState("lag-a", "lag")
	plan := pairState("lag-b", "lag")

	body, changed, _ := updatePair(t, pairFields(), plan, state)
	if !changed {
		t.Fatal("changing the referenced object reported no change")
	}
	if want := `{"lag":"lag-b"}`; body != want {
		t.Fatalf("update sent %s, want %s: with one permitted type the companion is implied", body, want)
	}
}

// When the type changes both halves go, because the server resolves them
// together and a type without its value names nothing.
func TestReferencePairTypeChangeSendsBothHalves(t *testing.T) {
	t.Parallel()

	fields := pairFields("lag", "bundle")
	state := pairState("thing-a", "lag")
	plan := pairState("thing-a", "bundle")

	body, changed, _ := updatePair(t, fields, plan, state)
	if !changed {
		t.Fatal("changing the reference type reported no change")
	}
	if want := `{"lag":"thing-a","lag_ref_type_":"bundle"}`; body != want {
		t.Fatalf("update sent %s, want %s: both halves travel together", body, want)
	}
}

// Both halves clear to an empty string. A reference is never cleared by
// omission and never by a null, whichever half it is.
func TestReferencePairClearsToEmptyStrings(t *testing.T) {
	t.Parallel()

	state := pairState("lag-a", "lag")
	plan := map[string]attr.Value{
		"name":          types.StringValue("p"),
		"lag":           types.StringNull(),
		"lag_ref_type_": types.StringNull(),
	}

	body, changed, _ := updatePair(t, pairFields(), plan, state)
	if !changed {
		t.Fatal("clearing a reference reported no change")
	}
	if want := `{"lag":"","lag_ref_type_":""}`; body != want {
		t.Fatalf("update sent %s, want %s", body, want)
	}
}

// A value with no type is refused rather than sent, because the server could not
// resolve what it names. Which refusal it is depends on whether the companion
// also changed, and both wordings are the handwritten helpers' own.
func TestReferencePairRefusesAValueWithNoType(t *testing.T) {
	t.Parallel()

	// The value alone changes and the companion stays absent: the pair is
	// incomplete, and the complaint names what is missing.
	valueOnlyState := map[string]attr.Value{
		"name":          types.StringValue("p"),
		"lag":           types.StringValue(""),
		"lag_ref_type_": types.StringNull(),
	}
	valueOnlyPlan := map[string]attr.Value{
		"name":          types.StringValue("p"),
		"lag":           types.StringValue("lag-a"),
		"lag_ref_type_": types.StringNull(),
	}
	_, _, diagnostics := updatePair(t, pairFields(), valueOnlyPlan, valueOnlyState)
	if !diagnostics.HasError() {
		t.Fatal("setting a reference value with no type was accepted")
	}
	if summary := diagnostics.Errors()[0].Summary(); summary != "Missing reference type" {
		t.Fatalf("diagnostic summary = %q, want %q", summary, "Missing reference type")
	}

	// Both halves change, the companion to nothing: the complaint is that the
	// two disagree, not that one is missing.
	_, _, both := updatePair(t, pairFields(), valueOnlyPlan, pairState("", ""))
	if !both.HasError() {
		t.Fatal("clearing the type while setting the value was accepted")
	}
	if summary := both.Errors()[0].Summary(); summary != "Inconsistent fields" {
		t.Fatalf("diagnostic summary = %q, want %q", summary, "Inconsistent fields")
	}
}

// With several permitted types the companion is never implied, so both halves
// are always sent and the unchanged half comes from state.
func TestReferencePairWithSeveralTypesAlwaysSendsBoth(t *testing.T) {
	t.Parallel()

	fields := pairFields("lag", "bundle", "switchpoint")
	state := pairState("thing-a", "lag")
	plan := pairState("thing-b", "lag")

	body, changed, _ := updatePair(t, fields, plan, state)
	if !changed {
		t.Fatal("changing the referenced object reported no change")
	}
	if want := `{"lag":"thing-b","lag_ref_type_":"lag"}`; body != want {
		t.Fatalf("update sent %s, want %s: the unchanged half is restated from state", body, want)
	}
}

// A companion the spec names but the resource does not carry would mean sending
// one half of a pair the API requires whole, so it is refused rather than
// silently treated as an ordinary string.
func TestReferencePairRefusesAMissingCompanion(t *testing.T) {
	t.Parallel()

	fields := pairFields()
	fields = fields[:2] // drop lag_ref_type_

	var diagnostics diag.Diagnostics
	if _, _, err := buildUpdate(fields, pairState("a", "lag"), pairState("b", "lag"), nullableSource{}, &diagnostics); err == nil {
		t.Fatal("a reference with no companion field was accepted")
	}
}
