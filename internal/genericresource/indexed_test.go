package genericresource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-verity/internal/spec"
)

func listField() spec.FieldSpec {
	index := notesMember()
	index.TerraformName, index.APIName, index.Kind = "index", "index", spec.FieldKindInt64
	index.UpdateClear = spec.UpdateClearZero
	pair := pairFields()
	filter, filterType := pair[1], pair[2]
	return spec.FieldSpec{
		TerraformName: "entries", APIName: "entries", Kind: spec.FieldKindList, ElementKind: spec.FieldKindObject,
		Access: spec.AccessOptional, Modes: []spec.Mode{spec.ModeDatacenter},
		Collection: &spec.CollectionSpec{Strategy: spec.CollectionIndexedPatch, IdentityField: "index"},
		CreateNull: spec.CreateNullOmit, UpdateClear: spec.UpdateClearOmit,
		UnknownPlan: spec.UnknownPlanReject, ResponseAbsence: spec.ResponseAbsenceTerraformNull,
		Fields: []spec.FieldSpec{index, notesMember(), filter, filterType},
	}
}

func entry(index attr.Value, notes string) map[string]attr.Value {
	return map[string]attr.Value{
		"index": index, "notes": types.StringValue(notes),
		"lag": types.StringValue(""), "lag_ref_type_": types.StringValue(""),
	}
}

func list(t *testing.T, entries ...map[string]attr.Value) attr.Value {
	t.Helper()
	field := listField()
	if len(entries) == 0 {
		return types.ListValueMust(elementType(field), nil)
	}
	value, err := listValue(field, entries)
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func listUpdate(t *testing.T, plan, state attr.Value) (string, bool) {
	t.Helper()
	fields := []spec.FieldSpec{nullableFields()[0], listField()}
	var diagnostics diag.Diagnostics
	object, changed, err := buildUpdate(fields,
		map[string]attr.Value{"name": types.StringValue("p"), "entries": plan},
		map[string]attr.Value{"name": types.StringValue("p"), "entries": state},
		nullableSource{}, &diagnostics)
	if err != nil || diagnostics.HasError() {
		t.Fatalf("buildUpdate: %v %v", err, diagnostics)
	}
	return encode(t, object), changed
}

func TestCreateListSendsEveryEntry(t *testing.T) {
	t.Parallel()

	fields := []spec.FieldSpec{nullableFields()[0], listField()}
	plan := map[string]attr.Value{
		"name":    types.StringValue("p"),
		"entries": list(t, entry(types.Int64Value(1), "a"), entry(types.Int64Unknown(), "b")),
	}
	object, err := buildCreate(fields, plan, nullableSource{})
	if err != nil {
		t.Fatalf("buildCreate: %v", err)
	}

	want := `{"entries":[{"index":1,"lag":"","lag_ref_type_":"","notes":"a"},{"lag":"","lag_ref_type_":"","notes":"b"}],"name":"p"}`
	if got := encode(t, object); got != want {
		t.Fatalf("create sent %s, want %s", got, want)
	}

	plan["entries"] = list(t)
	object, err = buildCreate(fields, plan, nullableSource{})
	if err != nil {
		t.Fatalf("buildCreate: %v", err)
	}
	if got := encode(t, object); got != `{"name":"p"}` {
		t.Fatalf("an empty list sent %s, want it omitted", got)
	}
}

func TestUpdateListReconcilesByIndex(t *testing.T) {
	t.Parallel()

	state := list(t, entry(types.Int64Value(1), "a"), entry(types.Int64Value(2), "b"), entry(types.Int64Value(3), "c"))

	cases := []struct {
		name        string
		plan        attr.Value
		want        string
		wantChanged bool
	}{
		{"a changed entry sends its index and the changed member",
			list(t, entry(types.Int64Value(1), "A"), entry(types.Int64Value(2), "b"), entry(types.Int64Value(3), "c")),
			`{"entries":[{"index":1,"notes":"A"}]}`, true},
		{"a new index is sent whole",
			list(t, entry(types.Int64Value(1), "a"), entry(types.Int64Value(2), "b"), entry(types.Int64Value(3), "c"), entry(types.Int64Value(9), "z")),
			`{"entries":[{"index":9,"lag":"","lag_ref_type_":"","notes":"z"}]}`, true},
		{"removed entries are sent as their index, in ascending order",
			list(t, entry(types.Int64Value(2), "b")),
			`{"entries":[{"index":1},{"index":3}]}`, true},
		{"entries matched by index ignore their order",
			list(t, entry(types.Int64Value(3), "c"), entry(types.Int64Value(1), "a"), entry(types.Int64Value(2), "b")),
			`{}`, false},
		{"an unknown index reads as zero, a create with no index",
			list(t, entry(types.Int64Value(1), "a"), entry(types.Int64Value(2), "b"), entry(types.Int64Value(3), "c"), entry(types.Int64Unknown(), "n")),
			`{"entries":[{"lag":"","lag_ref_type_":"","notes":"n"}]}`, true},
		{"a null index is ignored",
			list(t, entry(types.Int64Value(1), "a"), entry(types.Int64Value(2), "b"), entry(types.Int64Value(3), "c"), entry(types.Int64Null(), "n")),
			`{}`, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, changed := listUpdate(t, tc.plan, state)
			if changed != tc.wantChanged {
				t.Fatalf("changed = %v, want %v", changed, tc.wantChanged)
			}
			if got != tc.want {
				t.Fatalf("update sent %s, want %s", got, tc.want)
			}
		})
	}
}

func TestUpdateListAppliesReferencePairsInsideEntries(t *testing.T) {
	t.Parallel()

	before := entry(types.Int64Value(1), "a")
	after := entry(types.Int64Value(1), "a")
	after["lag"], after["lag_ref_type_"] = types.StringValue("lag-a"), types.StringValue("lag")

	got, changed := listUpdate(t, list(t, after), list(t, before))
	if !changed || got != `{"entries":[{"index":1,"lag":"lag-a","lag_ref_type_":"lag"}]}` {
		t.Fatalf("update sent %s (changed %v), want both halves of the pair with the index", got, changed)
	}
}

func TestListFromAPI(t *testing.T) {
	t.Parallel()

	field := listField()
	decoded, err := listFromAPI(field, []interface{}{
		map[string]interface{}{"index": float64(2), "notes": "second"},
		"not an object",
		map[string]interface{}{"index": float64(1), "notes": "first"},
	}, "datacenter")
	if err != nil {
		t.Fatalf("listFromAPI: %v", err)
	}
	entries, present, err := listEntries(field, decoded)
	if err != nil || !present || len(entries) != 2 {
		t.Fatalf("decoded %v entries (present %v, err %v), want 2", len(entries), present, err)
	}

	if !entries[0]["index"].Equal(types.Int64Value(2)) || !entries[1]["index"].Equal(types.Int64Value(1)) {
		t.Fatalf("entries decoded out of the response's order: %v", entries)
	}

	for name, raw := range map[string]interface{}{"empty": []interface{}{}, "missing": nil, "not a list": "x"} {
		value, err := listFromAPI(field, raw, "datacenter")
		if err != nil || !value.IsNull() {
			t.Errorf("%s response decoded to %v (err %v), want an absent block", name, value, err)
		}
	}
}

func TestSettleListNullsUnknownMembers(t *testing.T) {
	t.Parallel()

	field := listField()
	settled := settleList(field, list(t, entry(types.Int64Unknown(), "a")))
	entries, _, err := listEntries(field, settled)
	if err != nil || len(entries) != 1 || !entries[0]["index"].IsNull() {
		t.Fatalf("settled list = %v (err %v), want the unknown index as null", entries, err)
	}
}

func TestEntryMemberReadsTheConfigurationEntryByIndex(t *testing.T) {
	t.Parallel()

	field := listField()
	field.Fields[1].Kind, field.Fields[1].Nullable = spec.FieldKindInt64, true
	member := field.Fields[1]

	configEntry := func(index int64, notes attr.Value) map[string]attr.Value {
		return map[string]attr.Value{"index": types.Int64Value(index), "notes": notes, "lag": types.StringValue(""), "lag_ref_type_": types.StringValue("")}
	}
	config, err := listValue(field, []map[string]attr.Value{configEntry(2, types.Int64Null()), configEntry(1, types.Int64Value(7))})
	if err != nil {
		t.Fatal(err)
	}
	written := map[int64]bool{1: true, 2: true}
	source := nullableSource{
		config: map[string]attr.Value{"entries": config},
		indexed: func(block string, index int64, name string) bool {
			return block == "entries" && name == "notes" && written[index]
		},
	}

	planEntry := map[string]attr.Value{"index": types.Int64Value(2), "notes": types.Int64Unknown()}
	value, held := source.entryMember(field, member, planEntry)
	if !held || !value.IsNull() {
		t.Fatalf("index 2 read %v (held %v), want the configuration entry's written null", value, held)
	}

	planEntry["index"] = types.Int64Value(1)
	value, held = source.entryMember(field, member, planEntry)
	if !held || !value.Equal(types.Int64Value(7)) {
		t.Fatalf("index 1 read %v (held %v), want 7 from the matching entry, not the first one", value, held)
	}

	written[3] = true
	planEntry = map[string]attr.Value{"index": types.Int64Value(3), "notes": types.Int64Value(9)}
	value, held = source.entryMember(field, member, planEntry)
	if !held || !value.Equal(types.Int64Value(9)) {
		t.Fatalf("index 3 read %v (held %v), want the plan entry when no configuration entry carries the index", value, held)
	}

	planEntry["index"] = types.Int64Value(4)
	if _, held := source.entryMember(field, member, planEntry); held {
		t.Fatal("a member the scan did not record under the entry's index was treated as written")
	}
}
