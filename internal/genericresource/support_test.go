package genericresource

import (
	"strings"
	"testing"

	"terraform-provider-verity/internal/spec"
)

func singletonSpec(members ...spec.FieldSpec) spec.ResourceSpec {
	return spec.ResourceSpec{
		TerraformType: "verity_synthetic",
		Fields: []spec.FieldSpec{
			nullableFields()[0],
			{
				TerraformName: "object_properties", APIName: "object_properties", Kind: spec.FieldKindObject,
				Access: spec.AccessOptional, Modes: []spec.Mode{spec.ModeDatacenter},
				Collection: &spec.CollectionSpec{Strategy: spec.CollectionSingleton},
				CreateNull: spec.CreateNullOmit, UpdateClear: spec.UpdateClearOmit,
				UnknownPlan: spec.UnknownPlanReject, ResponseAbsence: spec.ResponseAbsenceTerraformNull,
				Fields: members,
			},
		},
	}
}

func notesMember() spec.FieldSpec {
	return spec.FieldSpec{
		TerraformName: "notes", APIName: "notes", Kind: spec.FieldKindString,
		Access: spec.AccessOptionalComputed, Modes: []spec.Mode{spec.ModeDatacenter},
		CreateNull: spec.CreateNullOmit, UpdateClear: spec.UpdateClearOmit,
		UnknownPlan: spec.UnknownPlanOmitAndRead, ResponseAbsence: spec.ResponseAbsenceTerraformNull,
	}
}

// Each refusal is something the engine would otherwise get silently wrong, so
// each is asserted with its reason rather than only as an error.
func TestSupportedRefusesWhatTheEngineWouldGetWrong(t *testing.T) {
	t.Parallel()

	if err := Supported(singletonSpec(notesMember())); err != nil {
		t.Fatalf("a singleton with a scalar member was refused: %v", err)
	}

	// A fixed parameter selecting between resources on one endpoint is served: the
	// write path passes it to the bulk manager and the read sends it as a query
	// parameter.
	withParameter := singletonSpec(notesMember())
	withParameter.API.FixedHeaders = map[string]string{"ip_version": "4"}
	if err := Supported(withParameter); err != nil {
		t.Fatalf("a resource selected by a fixed parameter was refused: %v", err)
	}

	indexMember := notesMember()
	indexMember.TerraformName, indexMember.APIName, indexMember.Kind = "index", "index", spec.FieldKindInt64
	indexedList := func(members ...spec.FieldSpec) spec.FieldSpec {
		return spec.FieldSpec{
			TerraformName: "entries", APIName: "entries", Kind: spec.FieldKindList, ElementKind: spec.FieldKindObject,
			Access: spec.AccessOptional, Modes: []spec.Mode{spec.ModeDatacenter},
			Collection: &spec.CollectionSpec{Strategy: spec.CollectionIndexedPatch, IdentityField: "index"},
			Fields:     members,
		}
	}
	withList := func(list spec.FieldSpec) spec.ResourceSpec {
		resource := singletonSpec(notesMember())
		resource.Fields = append(resource.Fields, list)
		return resource
	}
	if err := Supported(withList(indexedList(indexMember, notesMember()))); err != nil {
		t.Fatalf("an indexed_patch list of flat members was refused: %v", err)
	}

	otherStrategy := indexedList(indexMember, notesMember())
	otherStrategy.Collection = &spec.CollectionSpec{Strategy: spec.CollectionReplace, IdentityField: "index"}

	stringIndex := notesMember()
	stringIndex.TerraformName, stringIndex.APIName = "index", "index"

	nullableEntryMember := notesMember()
	nullableEntryMember.Kind, nullableEntryMember.Nullable = spec.FieldKindInt64, true

	listInBlock := singletonSpec(indexedList(indexMember, notesMember()))

	nestedAuto := notesMember()
	nestedAuto.AutoAssignment = &spec.AutoAssignmentSpec{FlagField: "notes_auto_assigned_"}

	nullableMember := notesMember()
	nullableMember.Kind, nullableMember.Nullable = spec.FieldKindInt64, true

	nestedObject := notesMember()
	nestedObject.Kind = spec.FieldKindObject

	notSingleton := singletonSpec(notesMember())
	notSingleton.Fields[1].Collection = &spec.CollectionSpec{Strategy: spec.CollectionIndexedPatch}

	cases := []struct {
		name     string
		resource spec.ResourceSpec
		reason   string
	}{
		{"list with a strategy no resource uses", withList(otherStrategy), "without the indexed_patch strategy"},
		{"list identified by a non-integer member", withList(indexedList(stringIndex, notesMember())), "not an int64 member"},
		{"nullable member of a list entry", withList(indexedList(indexMember, nullableEntryMember)), "nullable member of a block"},
		{"list inside a block", listInBlock, "list inside a block"},
		{"auto-assignment inside an object", singletonSpec(nestedAuto), "auto-assignment pair inside an object"},
		{"nullable member of an object", singletonSpec(nullableMember), "nullable member of a block"},
		{"object inside an object", singletonSpec(nestedObject), "object inside an object"},
		{"object without the singleton strategy", notSingleton, "without the singleton strategy"},
		{"object with no members", singletonSpec(), "no members"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := Supported(tc.resource)
			if err == nil {
				t.Fatal("accepted")
			}
			if !strings.Contains(err.Error(), tc.reason) {
				t.Fatalf("refused for %q, want a reason mentioning %q", err, tc.reason)
			}
		})
	}
}
