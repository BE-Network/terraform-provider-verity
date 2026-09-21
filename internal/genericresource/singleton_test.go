package genericresource

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-verity/internal/spec"
	"terraform-provider-verity/internal/transport"
)

// singletonFields is a resource with an identity and one singleton block holding
// a plain member and a reference pair, the shape verity_lag carries.
func singletonFields() []spec.FieldSpec {
	pair := pairFields()
	lag, lagType := pair[1], pair[2]
	lag.UpdateClear, lagType.UpdateClear = spec.UpdateClearEmptyString, spec.UpdateClearEmptyString
	return singletonSpec(notesMember(), lag, lagType).Fields
}

func block(t *testing.T, members map[string]attr.Value) attr.Value {
	t.Helper()
	field := singletonFields()[1]
	value, err := singletonValue(field, members)
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func blockMembers(notes, lag, lagType attr.Value) map[string]attr.Value {
	return map[string]attr.Value{"notes": notes, "lag": lag, "lag_ref_type_": lagType}
}

func TestCreateSingleton(t *testing.T) {
	t.Parallel()

	fields := singletonFields()
	cases := []struct {
		name  string
		block attr.Value
		want  string
	}{
		{"written block sends its members", block(t, blockMembers(types.StringValue("n"), types.StringValue("l"), types.StringValue("lag"))),
			`{"name":"p","object_properties":{"lag":"l","lag_ref_type_":"lag","notes":"n"}}`},
		{"unknown members are left to the read", block(t, blockMembers(types.StringUnknown(), types.StringUnknown(), types.StringUnknown())),
			`{"name":"p","object_properties":{}}`},
		{"absent block is not sent", types.ListValueMust(elementType(fields[1]), nil),
			`{"name":"p"}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			plan := map[string]attr.Value{"name": types.StringValue("p"), "object_properties": tc.block}
			object, err := buildCreate(fields, plan, nullableSource{})
			if err != nil {
				t.Fatalf("buildCreate: %v", err)
			}
			if got := encode(t, object); got != tc.want {
				t.Fatalf("create sent %s, want %s", got, tc.want)
			}
		})
	}
}

func TestUpdateSingleton(t *testing.T) {
	t.Parallel()

	fields := singletonFields()
	state := map[string]attr.Value{
		"name":              types.StringValue("p"),
		"object_properties": block(t, blockMembers(types.StringValue("n"), types.StringValue("l"), types.StringValue("lag"))),
	}
	absent := types.ListValueMust(elementType(fields[1]), nil)

	cases := []struct {
		name        string
		planBlock   attr.Value
		stateBlock  attr.Value
		want        string
		wantChanged bool
	}{
		{"only the changed member travels", block(t, blockMembers(types.StringValue("m"), types.StringValue("l"), types.StringValue("lag"))), nil,
			`{"object_properties":{"notes":"m"}}`, true},
		{"a member cleared by omission still sends the object", block(t, blockMembers(types.StringNull(), types.StringValue("l"), types.StringValue("lag"))), nil,
			`{"object_properties":{}}`, true},
		{"a reference pair inside the block clears to empty strings", block(t, blockMembers(types.StringValue("n"), types.StringNull(), types.StringNull())), nil,
			`{"object_properties":{"lag":"","lag_ref_type_":""}}`, true},
		{"removing the block sends nothing", absent, nil, `{}`, false},
		{"adding the block to a state without one sends nothing", state["object_properties"], absent, `{}`, false},
		{"an unchanged block sends nothing", state["object_properties"], nil, `{}`, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			previous := state
			if tc.stateBlock != nil {
				previous = map[string]attr.Value{"name": types.StringValue("p"), "object_properties": tc.stateBlock}
			}
			plan := map[string]attr.Value{"name": types.StringValue("p"), "object_properties": tc.planBlock}
			var diagnostics diag.Diagnostics
			object, changed, err := buildUpdate(fields, plan, previous, nullableSource{}, &diagnostics)
			if err != nil || diagnostics.HasError() {
				t.Fatalf("buildUpdate: %v %v", err, diagnostics)
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

// A pair inside the block refuses the same combinations as at the top level, in
// the handwritten helper's words.
func TestUpdateSingletonValidatesNestedReferencePairs(t *testing.T) {
	t.Parallel()

	fields := singletonFields()
	state := map[string]attr.Value{
		"name":              types.StringValue("p"),
		"object_properties": block(t, blockMembers(types.StringValue("n"), types.StringValue(""), types.StringNull())),
	}
	plan := map[string]attr.Value{
		"name":              types.StringValue("p"),
		"object_properties": block(t, blockMembers(types.StringValue("n"), types.StringValue("lag-a"), types.StringNull())),
	}
	var diagnostics diag.Diagnostics
	if _, _, err := buildUpdate(fields, plan, state, nullableSource{}, &diagnostics); err != nil {
		t.Fatalf("buildUpdate: %v", err)
	}
	if !diagnostics.HasError() || diagnostics.Errors()[0].Summary() != "Missing reference type" {
		t.Fatalf("diagnostics = %v, want the handwritten helper's Missing reference type", diagnostics)
	}
}

func TestSingletonFromAPI(t *testing.T) {
	t.Parallel()

	fields := singletonFields()
	values, err := stateFromAPI(fields, map[string]interface{}{
		"name":              "p",
		"object_properties": map[string]interface{}{"notes": "n"},
	}, "datacenter", nil)
	if err != nil {
		t.Fatalf("stateFromAPI: %v", err)
	}
	want := block(t, blockMembers(types.StringValue("n"), types.StringNull(), types.StringNull()))
	if got := values["object_properties"]; !got.Equal(want) {
		t.Fatalf("object_properties = %v, want %v: a response object is a one-entry block, absent members null", got, want)
	}

	missing, err := stateFromAPI(fields, map[string]interface{}{"name": "p"}, "datacenter", nil)
	if err != nil {
		t.Fatalf("stateFromAPI: %v", err)
	}
	if got := missing["object_properties"]; !got.IsNull() {
		t.Fatalf("object_properties = %v with no response object, want a null block", got)
	}

	campus, err := stateFromAPI(fields, map[string]interface{}{
		"name":              "p",
		"object_properties": map[string]interface{}{"notes": "n"},
	}, "campus", nil)
	if err != nil {
		t.Fatalf("stateFromAPI: %v", err)
	}
	if got := campus["object_properties"]; !got.IsNull() {
		t.Fatalf("object_properties = %v in a mode it does not apply to, want a null block", got)
	}
}

// Writing a plan to state settles unknown members inside the block, not only at
// the top level; an unknown cannot be stored.
func TestNullifyUnknownSettlesSingletonMembers(t *testing.T) {
	t.Parallel()

	fields := singletonFields()
	settled := nullifyUnknown(fields, map[string]attr.Value{
		"name":              types.StringValue("p"),
		"object_properties": block(t, blockMembers(types.StringUnknown(), types.StringValue("l"), types.StringValue("lag"))),
	})
	want := block(t, blockMembers(types.StringNull(), types.StringValue("l"), types.StringValue("lag")))
	if got := settled["object_properties"]; !got.Equal(want) {
		t.Fatalf("object_properties = %v, want %v", got, want)
	}
}

// An object the API declares with no properties can only be present or absent.
// Adding it sends an empty object; removing it is a change with nothing to send,
// as the handwritten resources treat it.
func TestUpdateEmptySingletonFollowsPresence(t *testing.T) {
	t.Parallel()

	field := singletonSpec().Fields[1]
	present, err := singletonValue(field, map[string]attr.Value{})
	if err != nil {
		t.Fatal(err)
	}
	absent := types.ListValueMust(elementType(field), nil)

	cases := []struct {
		name              string
		plan, state       attr.Value
		wantSend, wantChg bool
	}{
		{"added", present, absent, true, true},
		{"removed", absent, present, false, true},
		{"present in both", present, present, false, false},
		{"absent in both", absent, absent, false, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			wire, send, changed, err := updateEmptySingleton(field, tc.plan, tc.state)
			if err != nil {
				t.Fatalf("updateEmptySingleton: %v", err)
			}
			if send != tc.wantSend || changed != tc.wantChg {
				t.Fatalf("send=%v changed=%v, want send=%v changed=%v", send, changed, tc.wantSend, tc.wantChg)
			}
			if send {
				if got := encode(t, transport.WireObject{"object_properties": wire}); got != `{"object_properties":{}}` {
					t.Fatalf("sent %s, want an empty object", got)
				}
			}
		})
	}
}

// graphBlock is verity_fabric's object_properties: a singleton holding only an
// indexed list whose entries carry nothing but their index.
func graphBlock() spec.FieldSpec {
	index := notesMember()
	index.TerraformName, index.APIName, index.Kind = "index", "index", spec.FieldKindInt64
	graphs := spec.FieldSpec{
		TerraformName: "system_graphs", APIName: "system_graphs", Kind: spec.FieldKindList, ElementKind: spec.FieldKindObject,
		Access: spec.AccessOptional, Modes: []spec.Mode{spec.ModeDatacenter},
		Collection: &spec.CollectionSpec{Strategy: spec.CollectionIndexedPatch, IdentityField: "index"},
		CreateNull: spec.CreateNullOmit, UpdateClear: spec.UpdateClearOmit,
		UnknownPlan: spec.UnknownPlanReject, ResponseAbsence: spec.ResponseAbsenceTerraformNull,
		Fields: []spec.FieldSpec{index},
	}
	return singletonSpec(graphs).Fields[1]
}

func graphs(t *testing.T, indexes ...int64) attr.Value {
	t.Helper()
	field := graphBlock()
	list := field.Fields[0]
	entries := make([]map[string]attr.Value, 0, len(indexes))
	for _, index := range indexes {
		entries = append(entries, map[string]attr.Value{"index": types.Int64Value(index)})
	}
	value := types.ListValueMust(elementType(list), nil)
	if len(entries) != 0 {
		var err error
		if value, err = listValueTyped(list, entries); err != nil {
			t.Fatal(err)
		}
	}
	block, err := singletonValue(field, map[string]attr.Value{"system_graphs": value})
	if err != nil {
		t.Fatal(err)
	}
	return block
}

func listValueTyped(field spec.FieldSpec, entries []map[string]attr.Value) (types.List, error) {
	value, err := listValue(field, entries)
	if err != nil {
		return types.List{}, err
	}
	return value.(types.List), nil
}

// A written block with no entries still sends its list, empty, as the
// handwritten fabric create does.
func TestCreateSingletonSendsAnEmptyNestedList(t *testing.T) {
	t.Parallel()

	wire, send, err := createSingleton(graphBlock(), graphs(t))
	if err != nil || !send {
		t.Fatalf("createSingleton: send=%v err=%v", send, err)
	}
	if got := encode(t, transport.WireObject{"object_properties": wire}); got != `{"object_properties":{"system_graphs":[]}}` {
		t.Fatalf("create sent %s, want an empty system_graphs list", got)
	}
}

// A nested list is reconciled whenever either side holds the block: adding the
// block creates its entries, and removing it deletes them.
func TestUpdateSingletonReconcilesANestedListAcrossPresence(t *testing.T) {
	t.Parallel()

	field := graphBlock()
	absent := types.ListValueMust(elementType(field), nil)
	cases := []struct {
		name        string
		plan, state attr.Value
		want        string
	}{
		{"added with an entry", graphs(t, 1), absent, `{"object_properties":{"system_graphs":[{"index":1}]}}`},
		{"removed with an entry", absent, graphs(t, 1), `{"object_properties":{"system_graphs":[{"index":1}]}}`},
		{"an entry added to a present block", graphs(t, 1, 2), graphs(t, 1), `{"object_properties":{"system_graphs":[{"index":2}]}}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var diagnostics diag.Diagnostics
			wire, changed, err := updateSingleton(field, tc.plan, tc.state, &diagnostics)
			if err != nil || diagnostics.HasError() || !changed {
				t.Fatalf("updateSingleton: changed=%v err=%v %v", changed, err, diagnostics)
			}
			if got := encode(t, transport.WireObject{"object_properties": wire}); got != tc.want {
				t.Fatalf("update sent %s, want %s", got, tc.want)
			}
		})
	}
}
