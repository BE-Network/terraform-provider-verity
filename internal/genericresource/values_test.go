package genericresource

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-verity/internal/spec"
	"terraform-provider-verity/internal/transport"
)

// scalarFields mirrors the reviewed spec for verity_ipv4_list: a Required
// identity that refuses null and unknown, and two optional-computed fields that
// omit a null on create and clear by wire type on update.
func scalarFields() []spec.FieldSpec {
	return []spec.FieldSpec{
		{
			TerraformName: "name", APIName: "name", Kind: spec.FieldKindString,
			Access: spec.AccessRequired, Replace: true, Modes: []spec.Mode{spec.ModeDatacenter},
			CreateNull: spec.CreateNullReject, UpdateClear: spec.UpdateClearReject,
			UnknownPlan: spec.UnknownPlanReject, ResponseAbsence: spec.ResponseAbsenceTerraformNull,
		},
		{
			TerraformName: "enable", APIName: "enable", Kind: spec.FieldKindBool,
			Access: spec.AccessOptionalComputed, Modes: []spec.Mode{spec.ModeDatacenter},
			CreateNull: spec.CreateNullOmit, UpdateClear: spec.UpdateClearFalse,
			UnknownPlan: spec.UnknownPlanOmitAndRead, ResponseAbsence: spec.ResponseAbsenceTerraformNull,
		},
		{
			TerraformName: "ipv4_list", APIName: "ipv4_list", Kind: spec.FieldKindString,
			Access: spec.AccessOptionalComputed, Modes: []spec.Mode{spec.ModeDatacenter},
			CreateNull: spec.CreateNullOmit, UpdateClear: spec.UpdateClearEmptyString,
			UnknownPlan: spec.UnknownPlanOmitAndRead, ResponseAbsence: spec.ResponseAbsenceTerraformNull,
		},
	}
}

func encode(t *testing.T, object transport.WireObject) string {
	t.Helper()
	encoded, err := json.Marshal(object)
	if err != nil {
		t.Fatal(err)
	}
	return string(encoded)
}

// A create sends what the configuration says and leaves out what it does not.
// The omitted field is Computed, so the read that follows supplies it; sending a
// zero value instead would store a value nobody asked for.
func TestBuildCreateOmitsNullAndUnknown(t *testing.T) {
	t.Parallel()

	object, err := buildCreate(scalarFields(), map[string]attr.Value{
		"name":      types.StringValue("list-a"),
		"enable":    types.BoolNull(),
		"ipv4_list": types.StringUnknown(),
	}, nullableSource{})
	if err != nil {
		t.Fatalf("buildCreate: %v", err)
	}
	if got, want := encode(t, object), `{"name":"list-a"}`; got != want {
		t.Fatalf("create sent %s, want %s: a null omits by create_null and an unknown by unknown_plan", got, want)
	}
}

func TestBuildCreateSendsKnownValues(t *testing.T) {
	t.Parallel()

	object, err := buildCreate(scalarFields(), map[string]attr.Value{
		"name":      types.StringValue("list-a"),
		"enable":    types.BoolValue(true),
		"ipv4_list": types.StringValue("10.0.0.1"),
	}, nullableSource{})
	if err != nil {
		t.Fatalf("buildCreate: %v", err)
	}
	if got, want := encode(t, object), `{"enable":true,"ipv4_list":"10.0.0.1","name":"list-a"}`; got != want {
		t.Fatalf("create sent %s, want %s", got, want)
	}
}

// The identity refuses both null and unknown rather than sending something the
// server would store under the wrong name.
func TestBuildCreateRejectsUnknownIdentity(t *testing.T) {
	t.Parallel()

	if _, err := buildCreate(scalarFields(), map[string]attr.Value{
		"name": types.StringUnknown(),
	}, nullableSource{}); err == nil {
		t.Fatal("an unknown identity was accepted; unknown_plan: reject must refuse it")
	}
	if _, err := buildCreate(scalarFields(), map[string]attr.Value{
		"name": types.StringNull(),
	}, nullableSource{}); err == nil {
		t.Fatal("a null identity was accepted; create_null: reject must refuse it")
	}
}

// An update carries only what differs. A field equal to its state is not a
// change, and when nothing differs there is no request to make at all.
func TestBuildUpdateSendsOnlyChangedFields(t *testing.T) {
	t.Parallel()

	state := map[string]attr.Value{
		"name":      types.StringValue("list-a"),
		"enable":    types.BoolValue(true),
		"ipv4_list": types.StringValue("10.0.0.1"),
	}

	unchanged, changed, err := buildUpdate(scalarFields(), state, state, nullableSource{}, &diag.Diagnostics{})
	if err != nil {
		t.Fatalf("buildUpdate: %v", err)
	}
	if changed {
		t.Fatalf("an identical plan reported a change and would send %s", encode(t, unchanged))
	}

	plan := map[string]attr.Value{
		"name":      types.StringValue("list-a"),
		"enable":    types.BoolValue(false),
		"ipv4_list": types.StringValue("10.0.0.1"),
	}
	object, changed, err := buildUpdate(scalarFields(), plan, state, nullableSource{}, &diag.Diagnostics{})
	if err != nil {
		t.Fatalf("buildUpdate: %v", err)
	}
	if !changed {
		t.Fatal("a changed plan reported no change")
	}
	if got, want := encode(t, object), `{"enable":false}`; got != want {
		t.Fatalf("update sent %s, want %s: only the field that differs", got, want)
	}
}

// Clearing is per field, and the wire type is what decides it: a string clears
// to empty and a bool to false. Neither is a null on this API.
func TestBuildUpdateClearsByDeclaredPolicy(t *testing.T) {
	t.Parallel()

	state := map[string]attr.Value{
		"name":      types.StringValue("list-a"),
		"enable":    types.BoolValue(true),
		"ipv4_list": types.StringValue("10.0.0.1"),
	}
	plan := map[string]attr.Value{
		"name":      types.StringValue("list-a"),
		"enable":    types.BoolNull(),
		"ipv4_list": types.StringNull(),
	}

	object, changed, err := buildUpdate(scalarFields(), plan, state, nullableSource{}, &diag.Diagnostics{})
	if err != nil {
		t.Fatalf("buildUpdate: %v", err)
	}
	if !changed {
		t.Fatal("clearing two fields reported no change")
	}
	if got, want := encode(t, object), `{"enable":false,"ipv4_list":""}`; got != want {
		t.Fatalf("update sent %s, want %s", got, want)
	}
}

// An unknown on update is left out, the same as on create, rather than
// serialized as the zero its Go type would produce. This is the policy the
// registry records for these fields, and the legacy compare helpers do not
// implement it; see the Phase 2 notes in status.md.
func TestBuildUpdateOmitsUnknown(t *testing.T) {
	t.Parallel()

	state := map[string]attr.Value{
		"name":      types.StringValue("list-a"),
		"enable":    types.BoolValue(true),
		"ipv4_list": types.StringValue("10.0.0.1"),
	}
	plan := map[string]attr.Value{
		"name":      types.StringValue("list-a"),
		"enable":    types.BoolValue(true),
		"ipv4_list": types.StringUnknown(),
	}

	object, changed, err := buildUpdate(scalarFields(), plan, state, nullableSource{}, &diag.Diagnostics{})
	if err != nil {
		t.Fatalf("buildUpdate: %v", err)
	}
	if changed {
		t.Fatalf("an unknown alone reported a change and would send %s", encode(t, object))
	}
	if _, sent := object["ipv4_list"]; sent {
		t.Fatalf("an unknown reached the request as %s", encode(t, object))
	}
}

// A clear that must not happen is refused rather than sent as something else.
func TestBuildUpdateRefusesToClearTheIdentity(t *testing.T) {
	t.Parallel()

	state := map[string]attr.Value{"name": types.StringValue("list-a")}
	plan := map[string]attr.Value{"name": types.StringNull()}

	if _, _, err := buildUpdate(scalarFields(), plan, state, nullableSource{}, &diag.Diagnostics{}); err == nil {
		t.Fatal("clearing the identity was accepted; update_clear: reject must refuse it")
	}
}

// Decoding a response is the other direction, and an absent member is the case
// worth pinning: it reads as a Terraform null rather than a zero value.
func TestStateFromAPIDecodesPresentAndAbsent(t *testing.T) {
	t.Parallel()

	values, err := stateFromAPI(scalarFields(), map[string]interface{}{
		"name":   "list-a",
		"enable": true,
	}, "datacenter", nil)
	if err != nil {
		t.Fatalf("stateFromAPI: %v", err)
	}

	if got := values["name"]; !got.Equal(types.StringValue("list-a")) {
		t.Errorf("name = %v, want \"list-a\"", got)
	}
	if got := values["enable"]; !got.Equal(types.BoolValue(true)) {
		t.Errorf("enable = %v, want true", got)
	}
	if got := values["ipv4_list"]; !got.Equal(types.StringNull()) {
		t.Errorf("ipv4_list = %v, want null: the response carried no such member", got)
	}
}
