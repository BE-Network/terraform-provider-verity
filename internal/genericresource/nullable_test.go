package genericresource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-verity/internal/spec"
)

func nullableFields() []spec.FieldSpec {
	return []spec.FieldSpec{
		{
			TerraformName: "name", APIName: "name", Kind: spec.FieldKindString,
			Access: spec.AccessRequired, Modes: []spec.Mode{spec.ModeDatacenter},
			CreateNull: spec.CreateNullReject, UpdateClear: spec.UpdateClearReject,
			UnknownPlan: spec.UnknownPlanReject, ResponseAbsence: spec.ResponseAbsenceTerraformNull,
		},
		{
			TerraformName: "poll_interval", APIName: "poll_interval", Kind: spec.FieldKindInt64,
			Access: spec.AccessOptionalComputed, Modes: []spec.Mode{spec.ModeDatacenter}, Nullable: true,
			CreateNull: spec.CreateNullAPINull, UpdateClear: spec.UpdateClearAPINull,
			UnknownPlan: spec.UnknownPlanOmitAndRead, ResponseAbsence: spec.ResponseAbsenceTerraformNull,
		},
	}
}

func written(value attr.Value) nullableSource {
	return nullableSource{
		config:     map[string]attr.Value{"poll_interval": value},
		configured: func(name string) bool { return name == "poll_interval" },
	}
}

func notWritten() nullableSource {
	return nullableSource{
		config:     map[string]attr.Value{"poll_interval": types.Int64Null()},
		configured: func(string) bool { return false },
	}
}

func TestBuildCreateNullableFollowsTheConfigurationSource(t *testing.T) {
	t.Parallel()

	plan := map[string]attr.Value{
		"name":          types.StringValue("p"),
		"poll_interval": types.Int64Unknown(),
	}

	cases := []struct {
		name   string
		source nullableSource
		want   string
	}{
		{"written as null is sent as an explicit null", written(types.Int64Null()), `{"name":"p","poll_interval":null}`},
		{"written as a value is sent", written(types.Int64Value(30)), `{"name":"p","poll_interval":30}`},
		{"not written is left to the server", notWritten(), `{"name":"p"}`},
		{"written but unknown is left to the read", written(types.Int64Unknown()), `{"name":"p"}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			object, err := buildCreate(nullableFields(), plan, tc.source)
			if err != nil {
				t.Fatalf("buildCreate: %v", err)
			}
			if got := encode(t, object); got != tc.want {
				t.Fatalf("create sent %s, want %s", got, tc.want)
			}
		})
	}
}

func TestBuildUpdateNullableFollowsTheConfigurationSource(t *testing.T) {
	t.Parallel()

	state := map[string]attr.Value{
		"name":          types.StringValue("p"),
		"poll_interval": types.Int64Value(30),
	}

	plan := state

	cases := []struct {
		name        string
		source      nullableSource
		want        string
		wantChanged bool
	}{
		{"written as null clears with an explicit null", written(types.Int64Null()), `{"poll_interval":null}`, true},
		{"written as a new value is sent", written(types.Int64Value(60)), `{"poll_interval":60}`, true},
		{"written as the same value is not a change", written(types.Int64Value(30)), `{}`, false},
		{"not written is left alone", notWritten(), `{}`, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			object, changed, err := buildUpdate(nullableFields(), plan, state, tc.source, &diag.Diagnostics{})
			if err != nil {
				t.Fatalf("buildUpdate: %v", err)
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

func nullableWithPolicies(createNull spec.CreateNullPolicy, updateClear spec.UpdateClearPolicy, unknown spec.UnknownPlanPolicy) []spec.FieldSpec {
	fields := nullableFields()
	fields[1].CreateNull = createNull
	fields[1].UpdateClear = updateClear
	fields[1].UnknownPlan = unknown
	return fields
}

func TestNullableCreateHonorsDeclaredPolicies(t *testing.T) {
	t.Parallel()

	plan := map[string]attr.Value{"name": types.StringValue("p"), "poll_interval": types.Int64Unknown()}

	omit := nullableWithPolicies(spec.CreateNullOmit, spec.UpdateClearAPINull, spec.UnknownPlanOmitAndRead)
	object, err := buildCreate(omit, plan, written(types.Int64Null()))
	if err != nil {
		t.Fatalf("buildCreate: %v", err)
	}
	if got, want := encode(t, object), `{"name":"p"}`; got != want {
		t.Fatalf("create_null: omit sent %s, want %s: a written null must follow the declared policy", got, want)
	}

	reject := nullableWithPolicies(spec.CreateNullReject, spec.UpdateClearAPINull, spec.UnknownPlanOmitAndRead)
	if _, err := buildCreate(reject, plan, written(types.Int64Null())); err == nil {
		t.Fatal("create_null: reject accepted a written null")
	}

	rejectUnknown := nullableWithPolicies(spec.CreateNullAPINull, spec.UpdateClearAPINull, spec.UnknownPlanReject)
	if _, err := buildCreate(rejectUnknown, plan, written(types.Int64Unknown())); err == nil {
		t.Fatal("unknown_plan: reject accepted a written unknown")
	}
}

func TestNullableUpdateHonorsDeclaredPolicies(t *testing.T) {
	t.Parallel()

	state := map[string]attr.Value{"name": types.StringValue("p"), "poll_interval": types.Int64Value(30)}

	cases := []struct {
		name        string
		clear       spec.UpdateClearPolicy
		want        string
		wantChanged bool
	}{
		{"update_clear: zero sends zero", spec.UpdateClearZero, `{"poll_interval":0}`, true},
		{"update_clear: omit is a change with nothing to send", spec.UpdateClearOmit, `{}`, true},
		{"update_clear: api_null sends null", spec.UpdateClearAPINull, `{"poll_interval":null}`, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			fields := nullableWithPolicies(spec.CreateNullAPINull, tc.clear, spec.UnknownPlanOmitAndRead)
			object, changed, err := buildUpdate(fields, state, state, written(types.Int64Null()), &diag.Diagnostics{})
			if err != nil {
				t.Fatalf("buildUpdate: %v", err)
			}
			if changed != tc.wantChanged {
				t.Fatalf("changed = %v, want %v", changed, tc.wantChanged)
			}
			if got := encode(t, object); got != tc.want {
				t.Fatalf("update sent %s, want %s", got, tc.want)
			}
		})
	}

	reject := nullableWithPolicies(spec.CreateNullAPINull, spec.UpdateClearReject, spec.UnknownPlanOmitAndRead)
	if _, _, err := buildUpdate(reject, state, state, written(types.Int64Null()), &diag.Diagnostics{}); err == nil {
		t.Fatal("update_clear: reject accepted a written null")
	}

	rejectUnknown := nullableWithPolicies(spec.CreateNullAPINull, spec.UpdateClearAPINull, spec.UnknownPlanReject)
	if _, _, err := buildUpdate(rejectUnknown, state, state, written(types.Int64Unknown()), &diag.Diagnostics{}); err == nil {
		t.Fatal("unknown_plan: reject accepted a written unknown on update")
	}
}
