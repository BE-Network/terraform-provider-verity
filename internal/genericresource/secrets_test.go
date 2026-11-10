package genericresource

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"terraform-provider-verity/internal/spec"
	"terraform-provider-verity/internal/transport"
)

const placeholder = "this field has been set"

func secretTestFields() []spec.FieldSpec {
	both := []spec.Mode{spec.ModeDatacenter, spec.ModeCampus}
	optional := func(name string, kind spec.FieldKind, secret spec.SecretKind) spec.FieldSpec {
		clear := spec.UpdateClearEmptyString
		if kind == spec.FieldKindBool {
			clear = spec.UpdateClearFalse
		}
		return spec.FieldSpec{
			TerraformName: name, APIName: name, Kind: kind, Access: spec.AccessOptionalComputed, Modes: both,
			CreateNull: spec.CreateNullOmit, UpdateClear: clear, UnknownPlan: spec.UnknownPlanOmitAndRead,
			ResponseAbsence: spec.ResponseAbsenceTerraformNull, StateOwnership: spec.StateConfigurationOrServer,
			Sensitive: secret != "", Secret: secret,
		}
	}
	encrypted := optional("ssh_key_or_password_encrypted", spec.FieldKindString, spec.SecretMasked)
	encrypted.AutoAssignment = &spec.AutoAssignmentSpec{FlagField: "ssh_key_or_password_encrypted_auto_assigned_"}
	return []spec.FieldSpec{
		{
			TerraformName: "name", APIName: "name", Kind: spec.FieldKindString, Access: spec.AccessRequired,
			Replace: true, Modes: both, CreateNull: spec.CreateNullReject, UpdateClear: spec.UpdateClearReject,
			UnknownPlan: spec.UnknownPlanReject, ResponseAbsence: spec.ResponseAbsenceTerraformNull,
			StateOwnership: spec.StateConfiguration,
		},
		optional("enable", spec.FieldKindBool, ""),
		optional("ssh_key_or_password", spec.FieldKindString, spec.SecretWriteOnly),
		encrypted,
		optional("ssh_key_or_password_encrypted_auto_assigned_", spec.FieldKindBool, ""),
	}
}

func secretValues(overrides map[string]attr.Value) map[string]attr.Value {
	values := map[string]attr.Value{
		"name":                          types.StringValue("leaf1"),
		"enable":                        types.BoolValue(true),
		"ssh_key_or_password":           types.StringNull(),
		"ssh_key_or_password_encrypted": types.StringNull(),
		"ssh_key_or_password_encrypted_auto_assigned_": types.BoolValue(false),
	}
	for key, value := range overrides {
		values[key] = value
	}
	return values
}

func TestPreserveSecretsRules(t *testing.T) {
	t.Parallel()
	fields := secretTestFields()
	for _, tc := range []struct {
		name   string
		data   map[string]interface{}
		source map[string]attr.Value
		want   map[string]interface{}
	}{
		{
			name:   "write-only returned empty keeps the known value",
			data:   map[string]interface{}{"ssh_key_or_password": ""},
			source: map[string]attr.Value{"ssh_key_or_password": types.StringValue("test-password-1")},
			want:   map[string]interface{}{"ssh_key_or_password": "test-password-1"},
		},
		{
			name:   "write-only missing from the response keeps the known value",
			data:   map[string]interface{}{},
			source: map[string]attr.Value{"ssh_key_or_password": types.StringValue("test-password-1")},
			want:   map[string]interface{}{"ssh_key_or_password": "test-password-1"},
		},
		{
			name:   "write-only with nothing known keeps the server value",
			data:   map[string]interface{}{"ssh_key_or_password": ""},
			source: map[string]attr.Value{"ssh_key_or_password": types.StringNull()},
			want:   map[string]interface{}{"ssh_key_or_password": ""},
		},
		{
			name:   "write-only with an unknown plan keeps the server value",
			data:   map[string]interface{}{"ssh_key_or_password": ""},
			source: map[string]attr.Value{"ssh_key_or_password": types.StringUnknown()},
			want:   map[string]interface{}{"ssh_key_or_password": ""},
		},
		{
			name:   "masked placeholder keeps the known value",
			data:   map[string]interface{}{"ssh_key_or_password_encrypted": placeholder},
			source: map[string]attr.Value{"ssh_key_or_password_encrypted": types.StringValue("$enc$abc")},
			want:   map[string]interface{}{"ssh_key_or_password_encrypted": "$enc$abc"},
		},
		{
			name:   "masked placeholder matches regardless of case",
			data:   map[string]interface{}{"ssh_key_or_password_encrypted": "This Field Has Been Set"},
			source: map[string]attr.Value{"ssh_key_or_password_encrypted": types.StringValue("$enc$abc")},
			want:   map[string]interface{}{"ssh_key_or_password_encrypted": "$enc$abc"},
		},
		{
			name:   "masked returned empty while a value is planned is not hidden",
			data:   map[string]interface{}{"ssh_key_or_password_encrypted": ""},
			source: map[string]attr.Value{"ssh_key_or_password_encrypted": types.StringValue("$enc$abc")},
			want:   map[string]interface{}{"ssh_key_or_password_encrypted": ""},
		},
		{
			name:   "masked with nothing known keeps the placeholder",
			data:   map[string]interface{}{"ssh_key_or_password_encrypted": placeholder},
			source: map[string]attr.Value{"ssh_key_or_password_encrypted": types.StringNull()},
			want:   map[string]interface{}{"ssh_key_or_password_encrypted": placeholder},
		},
		{
			name:   "ordinary fields are untouched",
			data:   map[string]interface{}{"enable": false},
			source: map[string]attr.Value{"enable": types.BoolValue(true)},
			want:   map[string]interface{}{"enable": false},
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			original := map[string]interface{}{}
			for key, value := range tc.data {
				original[key] = value
			}
			got := preserveSecrets(fields, tc.data, tc.source)
			for key, want := range tc.want {
				if got[key] != want {
					t.Errorf("%s = %#v, want %#v", key, got[key], want)
				}
			}
			for key, value := range original {
				if tc.data[key] != value {
					t.Errorf("input data was modified: %s", key)
				}
			}
		})
	}
}

func secretResource(t *testing.T, runtime *readRuntime) *Resource {
	t.Helper()
	r := &Resource{
		spec: spec.ResourceSpec{
			TerraformType: "verity_switchpoint", IdentityPath: "name", Fields: secretTestFields(),
			Operations: spec.OperationSpec{Create: true, Read: true, Update: true, Delete: true},
			API:        spec.APIResourceSpec{BulkKey: "switchpoint"},
		},
		runtime: runtime,
	}
	compiled, err := CompileSchema(r.spec)
	if err != nil {
		t.Fatal(err)
	}
	r.schema = compiled
	return r
}

func stateValue(t *testing.T, r *Resource, state tfsdk.State, name string) attr.Value {
	t.Helper()
	values, diags := readScalars(context.Background(), state, r.spec.Fields)
	if diags.HasError() {
		t.Fatal(diags)
	}
	return values[name]
}

func serverObject() map[string]interface{} {
	return map[string]interface{}{
		"enable":                        true,
		"ssh_key_or_password":           "",
		"ssh_key_or_password_encrypted": placeholder,
		"ssh_key_or_password_encrypted_auto_assigned_": false,
	}
}

func TestSecretsSurviveRefresh(t *testing.T) {
	t.Parallel()
	runtime := &readRuntime{mode: "datacenter", collection: map[string]interface{}{"leaf1": serverObject()}}
	r := secretResource(t, runtime)
	state := tfsdk.State{Schema: r.schema}
	if diags := r.setState(context.Background(), &state, secretValues(map[string]attr.Value{
		"ssh_key_or_password":           types.StringValue("test-password-1"),
		"ssh_key_or_password_encrypted": types.StringValue("$enc$abc"),
	})); diags.HasError() {
		t.Fatal(diags)
	}
	resp := resource.ReadResponse{State: state}
	r.Read(context.Background(), resource.ReadRequest{State: state}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatal(resp.Diagnostics)
	}
	if got := stateValue(t, r, resp.State, "ssh_key_or_password"); !got.Equal(types.StringValue("test-password-1")) {
		t.Errorf("plaintext password after refresh = %v", got)
	}
	if got := stateValue(t, r, resp.State, "ssh_key_or_password_encrypted"); !got.Equal(types.StringValue("$enc$abc")) {
		t.Errorf("encrypted password after refresh = %v", got)
	}
}

func TestSecretsSurviveFallbackRead(t *testing.T) {
	t.Parallel()
	runtime := &readRuntime{mode: "datacenter", collection: map[string]interface{}{"leaf1": serverObject()}}
	r := secretResource(t, runtime)
	plan := secretValues(map[string]attr.Value{"ssh_key_or_password": types.StringValue("test-password-1")})

	minimal := map[string]attr.Value{}
	for _, field := range r.spec.Fields {
		minimal[field.TerraformName] = nullFor(field)
	}
	minimal["name"] = types.StringValue("leaf1")
	state := tfsdk.State{Schema: r.schema}
	if diags := r.setState(context.Background(), &state, minimal); diags.HasError() {
		t.Fatal(diags)
	}
	resp := resource.ReadResponse{State: state}
	r.Read(r.withFallback(context.Background(), plan), resource.ReadRequest{State: state}, &resp)
	if resp.Diagnostics.HasError() {
		t.Fatal(resp.Diagnostics)
	}
	if got := stateValue(t, r, resp.State, "ssh_key_or_password"); !got.Equal(types.StringValue("test-password-1")) {
		t.Errorf("plaintext password after fallback read = %v, want the planned value", got)
	}
	if got := stateValue(t, r, resp.State, "ssh_key_or_password_encrypted"); !got.Equal(types.StringValue(placeholder)) {
		t.Errorf("encrypted password with nothing planned = %v, want the server placeholder", got)
	}
}

func TestSecretsSurviveCachedWriteResponse(t *testing.T) {
	t.Parallel()
	r := secretResource(t, &readRuntime{mode: "datacenter"})
	plan := secretValues(map[string]attr.Value{
		"ssh_key_or_password":           types.StringValue("test-password-1"),
		"ssh_key_or_password_encrypted": types.StringValue("$enc$abc"),
	})
	state := tfsdk.State{Schema: r.schema}
	var diags diag.Diagnostics
	r.applyResponse(context.Background(), r.mergePlanScalars(serverObject(), plan), plan, &state, &diags)
	if diags.HasError() {
		t.Fatal(diags)
	}
	for name, want := range map[string]string{"ssh_key_or_password": "test-password-1", "ssh_key_or_password_encrypted": "$enc$abc"} {
		if got := stateValue(t, r, state, name); !got.Equal(types.StringValue(want)) {
			t.Errorf("%s = %v, want %q", name, got, want)
		}
	}
}

func TestGuardOutboundSecrets(t *testing.T) {
	t.Parallel()
	fields := secretTestFields()

	object := transport.WireObject{"ssh_key_or_password_encrypted": transport.String(placeholder), "enable": transport.Bool(true)}
	var createDiags diag.Diagnostics
	guardOutboundSecrets(fields, object, true, "verity_switchpoint leaf1", &createDiags)
	if !createDiags.HasError() {
		t.Error("create with the placeholder must be an error")
	}
	if _, present := object["ssh_key_or_password_encrypted"]; present {
		t.Error("placeholder still in the create request")
	}

	object = transport.WireObject{"ssh_key_or_password_encrypted": transport.String(placeholder), "enable": transport.Bool(true)}
	var updateDiags diag.Diagnostics
	guardOutboundSecrets(fields, object, false, "verity_switchpoint leaf1", &updateDiags)
	if updateDiags.HasError() || updateDiags.WarningsCount() != 1 {
		t.Errorf("update with the placeholder: want one warning, got %v", updateDiags)
	}
	if _, present := object["ssh_key_or_password_encrypted"]; present {
		t.Error("placeholder still in the update request")
	}
	if _, present := object["enable"]; !present {
		t.Error("the rest of the update must still be sent")
	}

	object = transport.WireObject{"ssh_key_or_password_encrypted": transport.String("$enc$abc")}
	var realDiags diag.Diagnostics
	guardOutboundSecrets(fields, object, false, "verity_switchpoint leaf1", &realDiags)
	if realDiags.HasError() || realDiags.WarningsCount() != 0 || object["ssh_key_or_password_encrypted"].StringValue() != "$enc$abc" {
		t.Error("a real encrypted password must be sent unchanged")
	}
}

func TestFlagToggleDoesNotResendPlaceholder(t *testing.T) {
	t.Parallel()
	fields := secretTestFields()
	state := secretValues(map[string]attr.Value{
		"ssh_key_or_password_encrypted":                types.StringValue(placeholder),
		"ssh_key_or_password_encrypted_auto_assigned_": types.BoolValue(true),
	})
	plan := secretValues(map[string]attr.Value{
		"ssh_key_or_password_encrypted":                types.StringValue(placeholder),
		"ssh_key_or_password_encrypted_auto_assigned_": types.BoolValue(false),
	})
	var diags diag.Diagnostics
	object, changed, err := buildUpdate(fields, plan, state, nullableSource{config: plan}, &diags)
	if err != nil || diags.HasError() {
		t.Fatal(err, diags)
	}
	if !changed {
		t.Fatal("expected the flag change to produce an update")
	}
	if value, present := object["ssh_key_or_password_encrypted"]; !present || value.StringValue() != placeholder {
		t.Fatalf("precondition: updateAutoAssigned should resend the unchanged value, got %v", object)
	}
	guardOutboundSecrets(fields, object, false, "verity_switchpoint leaf1", &diags)
	if _, present := object["ssh_key_or_password_encrypted"]; present {
		t.Error("placeholder was resent with the flag change")
	}
	if flag, present := object["ssh_key_or_password_encrypted_auto_assigned_"]; !present || flag.BoolValue() {
		t.Error("the flag change itself must still be sent")
	}
}

func runModifyPlan(t *testing.T, r *Resource, prior map[string]attr.Value, planned map[string]attr.Value) diag.Diagnostics {
	t.Helper()
	ctx := context.Background()
	object := func(values map[string]attr.Value) types.Object {
		value, diags := types.ObjectValue(attributeTypes(r.spec.Fields), values)
		if diags.HasError() {
			t.Fatal(diags)
		}
		return value
	}
	plan := tfsdk.Plan{Schema: r.schema}
	if diags := plan.Set(ctx, object(planned)); diags.HasError() {
		t.Fatal(diags)
	}
	config := tfsdk.Config{Schema: r.schema, Raw: plan.Raw}
	state := tfsdk.State{Schema: r.schema, Raw: tftypes.NewValue(r.schema.Type().TerraformType(ctx), nil)}
	if prior != nil {
		if diags := state.Set(ctx, object(prior)); diags.HasError() {
			t.Fatal(diags)
		}
	}
	resp := resource.ModifyPlanResponse{Plan: plan}
	r.ModifyPlan(ctx, resource.ModifyPlanRequest{Plan: plan, State: state, Config: config}, &resp)
	return resp.Diagnostics
}

func diagnosticText(diags diag.Diagnostics) string {
	var parts []string
	for _, d := range diags {
		parts = append(parts, d.Summary()+": "+d.Detail())
	}
	return strings.Join(parts, "\n")
}

func TestPlanSecretsPlaceholderOnCreateIsAnError(t *testing.T) {
	t.Parallel()
	r := secretResource(t, &readRuntime{mode: "datacenter"})
	diags := runModifyPlan(t, r, nil, secretValues(map[string]attr.Value{
		"ssh_key_or_password_encrypted": types.StringValue(placeholder),
	}))
	if !diags.HasError() {
		t.Fatalf("expected an error, got: %s", diagnosticText(diags))
	}
}

func TestPlanSecretsChangeToPlaceholderIsAnError(t *testing.T) {
	t.Parallel()
	r := secretResource(t, &readRuntime{mode: "datacenter"})
	diags := runModifyPlan(t, r,
		secretValues(map[string]attr.Value{"ssh_key_or_password_encrypted": types.StringValue("$enc$abc")}),
		secretValues(map[string]attr.Value{"ssh_key_or_password_encrypted": types.StringValue(placeholder)}))
	if !diags.HasError() {
		t.Fatalf("expected an error, got: %s", diagnosticText(diags))
	}
}

func TestPlanSecretsUnchangedLegacyPlaceholderIsAWarning(t *testing.T) {
	t.Parallel()
	r := secretResource(t, &readRuntime{mode: "datacenter"})
	values := secretValues(map[string]attr.Value{"ssh_key_or_password_encrypted": types.StringValue(placeholder)})
	diags := runModifyPlan(t, r, values, values)
	if diags.HasError() || diags.WarningsCount() != 1 {
		t.Fatalf("expected exactly one warning, got: %s", diagnosticText(diags))
	}
}

func TestPlanSecretsFlagOffWithoutPasswordIsAnError(t *testing.T) {
	t.Parallel()
	r := secretResource(t, &readRuntime{mode: "datacenter"})
	prior := secretValues(map[string]attr.Value{
		"ssh_key_or_password_encrypted":                types.StringValue(placeholder),
		"ssh_key_or_password_encrypted_auto_assigned_": types.BoolValue(true),
	})
	planned := secretValues(map[string]attr.Value{
		"ssh_key_or_password_encrypted":                types.StringUnknown(),
		"ssh_key_or_password_encrypted_auto_assigned_": types.BoolValue(false),
	})
	if diags := runModifyPlan(t, r, prior, planned); !diags.HasError() || !strings.Contains(diagnosticText(diags), "requires setting") {
		t.Fatalf("expected a 'password required' error, got: %s", diagnosticText(diags))
	}
}

func TestPlanSecretsFlagOffResendingKnownRealPasswordIsAllowed(t *testing.T) {
	t.Parallel()
	r := secretResource(t, &readRuntime{mode: "datacenter"})
	prior := secretValues(map[string]attr.Value{
		"ssh_key_or_password_encrypted":                types.StringValue("$enc$known"),
		"ssh_key_or_password_encrypted_auto_assigned_": types.BoolValue(true),
	})
	planned := secretValues(map[string]attr.Value{
		"ssh_key_or_password_encrypted":                types.StringUnknown(),
		"ssh_key_or_password_encrypted_auto_assigned_": types.BoolValue(false),
	})
	if diags := runModifyPlan(t, r, prior, planned); diags.HasError() {
		t.Fatalf("resending a known real password must be allowed, got: %s", diagnosticText(diags))
	}
}

func TestPlanSecretsFlagOffWithRealPasswordIsAllowed(t *testing.T) {
	t.Parallel()
	r := secretResource(t, &readRuntime{mode: "datacenter"})
	prior := secretValues(map[string]attr.Value{
		"ssh_key_or_password_encrypted":                types.StringValue(placeholder),
		"ssh_key_or_password_encrypted_auto_assigned_": types.BoolValue(true),
	})
	planned := secretValues(map[string]attr.Value{
		"ssh_key_or_password_encrypted":                types.StringValue("$enc$new"),
		"ssh_key_or_password_encrypted_auto_assigned_": types.BoolValue(false),
	})
	if diags := runModifyPlan(t, r, prior, planned); diags.HasError() {
		t.Fatalf("expected no error, got: %s", diagnosticText(diags))
	}
}

func TestPlanSecretsFlagOffWithRealPlaintextPasswordIsAllowed(t *testing.T) {
	t.Parallel()
	r := secretResource(t, &readRuntime{mode: "datacenter"})
	prior := secretValues(map[string]attr.Value{
		"ssh_key_or_password_encrypted":                types.StringValue(placeholder),
		"ssh_key_or_password_encrypted_auto_assigned_": types.BoolValue(true),
	})
	planned := secretValues(map[string]attr.Value{
		"ssh_key_or_password":                          types.StringValue("test-password-1"),
		"ssh_key_or_password_encrypted":                types.StringUnknown(),
		"ssh_key_or_password_encrypted_auto_assigned_": types.BoolValue(false),
	})
	if diags := runModifyPlan(t, r, prior, planned); diags.HasError() {
		t.Fatalf("a real plaintext password must be accepted, got: %s", diagnosticText(diags))
	}
}

func TestFlagToggleWithPlaintextPasswordDropsPlaceholderSilently(t *testing.T) {
	t.Parallel()
	fields := secretTestFields()
	state := secretValues(map[string]attr.Value{
		"ssh_key_or_password_encrypted":                types.StringValue(placeholder),
		"ssh_key_or_password_encrypted_auto_assigned_": types.BoolValue(true),
	})
	plan := secretValues(map[string]attr.Value{
		"ssh_key_or_password":                          types.StringValue("test-password-1"),
		"ssh_key_or_password_encrypted":                types.StringUnknown(),
		"ssh_key_or_password_encrypted_auto_assigned_": types.BoolValue(false),
	})
	var diags diag.Diagnostics
	object, changed, err := buildUpdate(fields, plan, state, nullableSource{config: plan}, &diags)
	if err != nil || diags.HasError() || !changed {
		t.Fatal(err, diags, changed)
	}
	guardOutboundSecrets(fields, object, false, "verity_switchpoint leaf1", &diags)
	if _, present := object["ssh_key_or_password_encrypted"]; present {
		t.Errorf("the placeholder was resent: %v", object)
	}
	if value, present := object["ssh_key_or_password"]; !present || value.StringValue() != "test-password-1" {
		t.Errorf("the plaintext password was not sent: %v", object)
	}
	if flag, present := object["ssh_key_or_password_encrypted_auto_assigned_"]; !present || flag.BoolValue() {
		t.Errorf("the flag change was not sent: %v", object)
	}
	if diags.WarningsCount() != 0 || diags.HasError() {
		t.Errorf("expected no diagnostics, got: %s", diagnosticText(diags))
	}
}

func TestFlagToggleWithUnchangedPlaintextPasswordResendsIt(t *testing.T) {
	t.Parallel()
	fields := secretTestFields()
	state := secretValues(map[string]attr.Value{
		"ssh_key_or_password":                          types.StringValue("test-password-1"),
		"ssh_key_or_password_encrypted":                types.StringValue(placeholder),
		"ssh_key_or_password_encrypted_auto_assigned_": types.BoolValue(true),
	})
	plan := secretValues(map[string]attr.Value{
		"ssh_key_or_password":                          types.StringValue("test-password-1"),
		"ssh_key_or_password_encrypted":                types.StringUnknown(),
		"ssh_key_or_password_encrypted_auto_assigned_": types.BoolValue(false),
	})
	r := secretResource(t, &readRuntime{mode: "datacenter"})
	if diags := runModifyPlan(t, r, state, plan); diags.HasError() {
		t.Fatalf("plan: %s", diagnosticText(diags))
	}
	var diags diag.Diagnostics
	object, changed, err := buildUpdate(fields, plan, state, nullableSource{config: plan}, &diags)
	if err != nil || diags.HasError() || !changed {
		t.Fatal(err, diags, changed)
	}
	if _, present := object["ssh_key_or_password"]; present {
		t.Fatalf("precondition: an unchanged password is not part of the diff, got %v", object)
	}
	resendPlaintextOnFlagOff(fields, plan, state, object)
	guardOutboundSecrets(fields, object, false, "verity_switchpoint leaf1", &diags)
	if value, present := object["ssh_key_or_password"]; !present || value.StringValue() != "test-password-1" {
		t.Errorf("the known password must be sent with the flag change: %v", object)
	}
	if _, present := object["ssh_key_or_password_encrypted"]; present {
		t.Errorf("the placeholder was resent: %v", object)
	}
	if flag, present := object["ssh_key_or_password_encrypted_auto_assigned_"]; !present || flag.BoolValue() {
		t.Errorf("the flag change was not sent: %v", object)
	}
	if diags.WarningsCount() != 0 || diags.HasError() {
		t.Errorf("expected no diagnostics, got: %s", diagnosticText(diags))
	}
}

func TestResendPlaintextOnlyWhenTurningTheFlagOff(t *testing.T) {
	t.Parallel()
	fields := secretTestFields()
	values := secretValues(map[string]attr.Value{
		"ssh_key_or_password":                          types.StringValue("test-password-1"),
		"ssh_key_or_password_encrypted_auto_assigned_": types.BoolValue(false),
	})
	object := transport.WireObject{"enable": transport.Bool(false)}
	resendPlaintextOnFlagOff(fields, values, values, object)
	if _, present := object["ssh_key_or_password"]; present {
		t.Fatalf("the password must not be resent with an unrelated update: %v", object)
	}
}
