package genericresource

import (
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-verity/internal/spec"
)

func listFieldWithPortName() spec.FieldSpec {
	field := listField()
	portName := notesMember()
	portName.TerraformName, portName.APIName = "port_name", "port_name"
	portName.Access = spec.AccessComputed
	portName.StateOwnership = spec.StateServer
	portName.CreateNull = spec.CreateNullOmit
	portName.UpdateClear = spec.UpdateClearOmitUnmanaged
	field.Fields = append(field.Fields, portName)
	return field
}

func entryWithPortName(index int64, notes string, port attr.Value) map[string]attr.Value {
	values := entry(types.Int64Value(index), notes)
	values["port_name"] = port
	return values
}

func portList(t *testing.T, entries ...map[string]attr.Value) attr.Value {
	t.Helper()
	value, err := listValue(listFieldWithPortName(), entries)
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func TestServerOwnedListMemberIsNeverSentOnCreate(t *testing.T) {
	t.Parallel()
	fields := []spec.FieldSpec{nullableFields()[0], listFieldWithPortName()}
	plan := map[string]attr.Value{
		"name":    types.StringValue("p"),
		"entries": portList(t, entryWithPortName(1, "a", types.StringValue("Eth1/1")), entryWithPortName(2, "b", types.StringUnknown())),
	}
	object, err := buildCreate(fields, plan, nullableSource{})
	if err != nil {
		t.Fatal(err)
	}
	if got := encode(t, object); strings.Contains(got, "port_name") {
		t.Fatalf("create sent a read-only member: %s", got)
	}
}

func TestServerOwnedListMemberIsNeverSentOnUpdate(t *testing.T) {
	t.Parallel()
	fields := []spec.FieldSpec{nullableFields()[0], listFieldWithPortName()}
	state := portList(t, entryWithPortName(1, "a", types.StringValue("")))
	for name, plan := range map[string]attr.Value{
		"server now reports a port name":          portList(t, entryWithPortName(1, "a", types.StringValue("Eth1/1"))),
		"port name unknown, other member changed": portList(t, entryWithPortName(1, "changed", types.StringUnknown())),
	} {
		var diagnostics diag.Diagnostics
		object, _, err := buildUpdate(fields,
			map[string]attr.Value{"name": types.StringValue("p"), "entries": plan},
			map[string]attr.Value{"name": types.StringValue("p"), "entries": state},
			nullableSource{}, &diagnostics)
		if err != nil || diagnostics.HasError() {
			t.Fatalf("%s: %v %v", name, err, diagnostics)
		}
		if got := encode(t, object); strings.Contains(got, "port_name") {
			t.Errorf("%s: update sent a read-only member: %s", name, got)
		}
	}
}

func TestServerOwnedTopLevelFieldIsNeverSent(t *testing.T) {
	t.Parallel()
	fields := scalarFields()
	fields[2].Access, fields[2].StateOwnership, fields[2].UpdateClear = spec.AccessComputed, spec.StateServer, spec.UpdateClearOmitUnmanaged
	plan := map[string]attr.Value{"name": types.StringValue("x"), "enable": types.BoolValue(true), "ipv4_list": types.StringValue("server")}
	object, err := buildCreate(fields, plan, nullableSource{})
	if err != nil {
		t.Fatal(err)
	}
	if got := encode(t, object); strings.Contains(got, "ipv4_list") {
		t.Fatalf("create sent a read-only field: %s", got)
	}
	var diagnostics diag.Diagnostics
	state := map[string]attr.Value{"name": types.StringValue("x"), "enable": types.BoolValue(false), "ipv4_list": types.StringValue("old")}
	object, _, err = buildUpdate(fields, plan, state, nullableSource{}, &diagnostics)
	if err != nil || diagnostics.HasError() {
		t.Fatal(err, diagnostics)
	}
	if got := encode(t, object); strings.Contains(got, "ipv4_list") {
		t.Fatalf("update sent a read-only field: %s", got)
	}
}
