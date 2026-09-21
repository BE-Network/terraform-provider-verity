package genericresource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-verity/internal/spec"
	"terraform-provider-verity/internal/utils"
)

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

	if appliesToMode(bothModes, "") || appliesToMode(bothModes, "datacentre") {
		t.Error("an unrecognised mode was treated as matching")
	}
}

func TestLegacyModeTableFailsOpenForUnknownResources(t *testing.T) {
	t.Parallel()

	if !utils.FieldAppliesToMode("no_such_endpoint", "campus_only", "datacenter") {
		t.Skip("FieldAppliesToMode no longer fails open; the engine's reason for not using it has changed")
	}
}
