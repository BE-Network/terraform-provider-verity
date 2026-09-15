package genericresource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-verity/internal/spec"
	"terraform-provider-verity/internal/utils"
)

// modeFields is a synthetic resource whose endpoint exists in no generated
// table. That is the point: the engine must answer the mode question from the
// spec it holds, and a lookup keyed by endpoint name would find nothing here.
func modeFields() []spec.FieldSpec {
	return []spec.FieldSpec{
		{
			TerraformName: "name", APIName: "name", Kind: spec.FieldKindString,
			Access: spec.AccessRequired, Modes: []spec.Mode{spec.ModeDatacenter, spec.ModeCampus},
			CreateNull: spec.CreateNullReject, UpdateClear: spec.UpdateClearReject,
			UnknownPlan: spec.UnknownPlanReject, ResponseAbsence: spec.ResponseAbsenceTerraformNull,
		},
		{
			TerraformName: "campus_only", APIName: "campus_only", Kind: spec.FieldKindString,
			Access: spec.AccessOptionalComputed, Modes: []spec.Mode{spec.ModeCampus},
			CreateNull: spec.CreateNullOmit, UpdateClear: spec.UpdateClearEmptyString,
			UnknownPlan: spec.UnknownPlanOmitAndRead, ResponseAbsence: spec.ResponseAbsenceTerraformNull,
		},
		{
			TerraformName: "both_modes", APIName: "both_modes", Kind: spec.FieldKindString,
			Access: spec.AccessOptionalComputed, Modes: []spec.Mode{spec.ModeDatacenter, spec.ModeCampus},
			CreateNull: spec.CreateNullOmit, UpdateClear: spec.UpdateClearEmptyString,
			UnknownPlan: spec.UnknownPlanOmitAndRead, ResponseAbsence: spec.ResponseAbsenceTerraformNull,
		},
	}
}

// A field outside the running mode reads as null even when the response carries
// a value for it.
//
// This is what keeps a campus-only attribute from showing up in a datacenter
// plan. The engine decides it from FieldSpec.Modes; deciding it from the
// generated lookup table instead would answer "applies" for every field here,
// because that table fails open on an endpoint it does not know and this
// resource is in no table at all.
func TestStateFromAPIHonorsFieldModes(t *testing.T) {
	t.Parallel()

	response := map[string]interface{}{
		"name":        "synthetic",
		"campus_only": "served anyway",
		"both_modes":  "kept",
	}

	values, err := stateFromAPI(modeFields(), response, "datacenter", nil)
	if err != nil {
		t.Fatalf("stateFromAPI: %v", err)
	}
	if got := values["campus_only"]; !got.Equal(types.StringNull()) {
		t.Errorf("campus_only = %v in datacenter mode, want null: the response value must not reach state", got)
	}
	if got := values["both_modes"]; !got.Equal(types.StringValue("kept")) {
		t.Errorf("both_modes = %v, want \"kept\": a field in both modes is read normally", got)
	}

	campus, err := stateFromAPI(modeFields(), response, "campus", nil)
	if err != nil {
		t.Fatalf("stateFromAPI: %v", err)
	}
	if got := campus["campus_only"]; !got.Equal(types.StringValue("served anyway")) {
		t.Errorf("campus_only = %v in campus mode, want the response value", got)
	}
}

// The same decision drives ModifyPlan, which nullifies out-of-mode fields so a
// plan does not show "known after apply" for something the API will never
// return. Asserting the predicate directly covers both callers.
func TestAppliesToModeReadsTheSpec(t *testing.T) {
	t.Parallel()

	fields := modeFields()
	campusOnly, bothModes := fields[1], fields[2]

	if appliesToMode(campusOnly, "datacenter") {
		t.Error("a campus-only field was reported as applying to datacenter")
	}
	if !appliesToMode(campusOnly, "campus") {
		t.Error("a campus-only field was reported as not applying to campus")
	}
	if !appliesToMode(bothModes, "datacenter") || !appliesToMode(bothModes, "campus") {
		t.Error("a field declared for both modes was reported as not applying to one")
	}

	// An unrecognised mode matches nothing, rather than everything. The lookup
	// table this replaced returned true for an unknown mode, so a typo there
	// silently exposed every field.
	if appliesToMode(bothModes, "") || appliesToMode(bothModes, "datacentre") {
		t.Error("an unrecognised mode was treated as matching")
	}
}

// The distinction this turns on: the table the engine used to consult answers
// "applies" for a resource it has never heard of. Recording that here is what
// makes the choice above legible as a fix rather than a preference — if
// FieldAppliesToMode ever fails closed instead, this test says so and the
// comment above it can be revisited.
func TestLegacyModeTableFailsOpenForUnknownResources(t *testing.T) {
	t.Parallel()

	if !utils.FieldAppliesToMode("no_such_endpoint", "campus_only", "datacenter") {
		t.Skip("FieldAppliesToMode no longer fails open; the engine's reason for not using it has changed")
	}
}
