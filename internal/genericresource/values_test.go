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

func TestBuildUpdateRefusesToClearTheIdentity(t *testing.T) {
	t.Parallel()

	state := map[string]attr.Value{"name": types.StringValue("list-a")}
	plan := map[string]attr.Value{"name": types.StringNull()}

	if _, _, err := buildUpdate(scalarFields(), plan, state, nullableSource{}, &diag.Diagnostics{}); err == nil {
		t.Fatal("clearing the identity was accepted; update_clear: reject must refuse it")
	}
}

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
