package spec

import (
	"strings"
	"testing"
)

func TestResourceSpecValidateAcceptsCompleteIPv4List(t *testing.T) {
	if err := validIPv4ListSpec().Validate(); err != nil {
		t.Fatalf("valid spec rejected: %v", err)
	}
}

func TestResourceSpecValidateRejectsIncompleteAndContradictoryFields(t *testing.T) {
	tests := []struct {
		name string
		edit func(*ResourceSpec)
		want string
	}{
		{"missing policy", func(spec *ResourceSpec) { spec.Fields[0].UnknownPlan = "" }, "all lifecycle policies"},
		{"invalid ownership", func(spec *ResourceSpec) { spec.Fields[0].StateOwnership = StateServer }, "requires state ownership"},
		{"field range escapes resource", func(spec *ResourceSpec) { spec.Fields[0].Versions.MinInclusive = APIVersion{Major: 6, Minor: 5} }, "contained by its parent"},
		{"invalid deterministic integer", func(spec *ResourceSpec) {
			spec.Fields[1].Default = &LiteralSpec{Kind: LiteralInt64, Value: "01"}
			spec.Fields[1].Kind = FieldKindInt64
		}, "canonical base-10"},
		{"missing nested strategy", func(spec *ResourceSpec) { spec.Fields = append(spec.Fields, nestedField(nil)) }, "explicit collection policy"},
		{"missing nested identity field", func(spec *ResourceSpec) {
			spec.Fields = append(spec.Fields, nestedField(&CollectionSpec{Strategy: CollectionReplace, Ordering: CollectionOrdered, IdentityField: "missing"}))
		}, "does not name a nested field"},
		// An unmanaged field carries no Terraform behavior, but it still records
		// API shape and applicability, so those must be validated like any other.
		{"unmanaged field without a kind", func(spec *ResourceSpec) {
			spec.Fields = append(spec.Fields, unmanagedField(func(f *FieldSpec) { f.Kind = "" }))
		}, "field kind is required"},
		{"unmanaged field with an empty version range", func(spec *ResourceSpec) {
			spec.Fields = append(spec.Fields, unmanagedField(func(f *FieldSpec) { f.Versions = VersionRange{} }))
		}, "versions"},
		{"unmanaged field outside the resource version range", func(spec *ResourceSpec) {
			spec.Fields = append(spec.Fields, unmanagedField(func(f *FieldSpec) {
				f.Versions = VersionRange{MinInclusive: APIVersion{Major: 6, Minor: 5}, MaxExclusive: APIVersion{Major: 6, Minor: 7}}
			}))
		}, "contained by its parent"},
		{"unmanaged field outside the resource modes", func(spec *ResourceSpec) {
			spec.Fields = append(spec.Fields, unmanagedField(func(f *FieldSpec) { f.Modes = []Mode{ModeCampus} }))
		}, "subset of parent modes"},
		{"unmanaged field declaring Terraform behavior", func(spec *ResourceSpec) {
			spec.Fields = append(spec.Fields, unmanagedField(func(f *FieldSpec) { f.Access = AccessOptional }))
		}, "cannot declare access or lifecycle policies"},
		{"missing reference target", func(spec *ResourceSpec) {
			spec.Fields[1].Reference = &ReferenceSpec{TypeField: "missing", AllowedTypes: []string{"tenant"}}
		}, "does not exist"},
		{"field outside resource modes", func(spec *ResourceSpec) { spec.Fields[1].Modes = []Mode{ModeCampus} }, "subset of parent modes"},
		{"nested field outside parent modes", func(spec *ResourceSpec) {
			spec.Fields = append(spec.Fields, nestedField(&CollectionSpec{Strategy: CollectionSingleton, Ordering: CollectionOrdered}))
			spec.Fields[3].Fields[0].Modes = []Mode{ModeCampus}
		}, "subset of parent modes"},
		{"reference companion unavailable in field mode", func(spec *ResourceSpec) {
			spec.Modes = []Mode{ModeDatacenter, ModeCampus}
			spec.Fields[1].Reference = &ReferenceSpec{TypeField: "name", AllowedTypes: []string{"tenant"}}
			spec.Fields[1].Modes = []Mode{ModeDatacenter, ModeCampus}
		}, "unavailable in one or more field modes"},
		{"auto-assignment companion unavailable in field mode", func(spec *ResourceSpec) {
			spec.Modes = []Mode{ModeDatacenter, ModeCampus}
			spec.Fields[2].AutoAssignment = &AutoAssignmentSpec{FlagField: "enable"}
			spec.Fields[2].Modes = []Mode{ModeDatacenter, ModeCampus}
		}, "unavailable in one or more field modes"},
		{"reference companion unavailable in API version", func(spec *ResourceSpec) {
			spec.Versions.MaxExclusive = APIVersion{Major: 6, Minor: 8}
			spec.Fields[1].Reference = &ReferenceSpec{TypeField: "name", AllowedTypes: []string{"tenant"}}
			spec.Fields[1].Versions.MaxExclusive = APIVersion{Major: 6, Minor: 7}
			spec.Fields[0].Versions = VersionRange{MinInclusive: APIVersion{Major: 6, Minor: 7}, MaxExclusive: APIVersion{Major: 6, Minor: 8}}
		}, "unavailable in one or more field API versions"},
		{"auto-assignment companion unavailable in API version", func(spec *ResourceSpec) {
			spec.Versions.MaxExclusive = APIVersion{Major: 6, Minor: 8}
			spec.Fields[2].AutoAssignment = &AutoAssignmentSpec{FlagField: "enable"}
			spec.Fields[2].Versions.MaxExclusive = APIVersion{Major: 6, Minor: 7}
			spec.Fields[1].Versions = VersionRange{MinInclusive: APIVersion{Major: 6, Minor: 7}, MaxExclusive: APIVersion{Major: 6, Minor: 8}}
		}, "unavailable in one or more field API versions"},
		{"identity path missing", func(spec *ResourceSpec) { spec.IdentityPath = "missing" }, "does not name a top-level field"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			spec := validIPv4ListSpec()
			test.edit(&spec)
			err := spec.Validate()
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("Validate() error = %v, want substring %q", err, test.want)
			}
		})
	}
}

func TestResourceSpecValidateValidators(t *testing.T) {
	min, max := int64(1), int64(10)
	tests := []struct {
		name      string
		kind      FieldKind
		validator ValidatorSpec
		wantError bool
	}{
		{"one of strings", FieldKindString, ValidatorSpec{Kind: ValidatorOneOf, Values: []LiteralSpec{{Kind: LiteralString, Value: "a"}}}, false},
		{"string length", FieldKindString, ValidatorSpec{Kind: ValidatorStringLength, Min: &min, Max: &max}, false},
		{"int64 range", FieldKindInt64, ValidatorSpec{Kind: ValidatorInt64Range, Min: &min, Max: &max}, false},
		{"number range", FieldKindNumber, ValidatorSpec{Kind: ValidatorNumberRange, MinValue: &LiteralSpec{Kind: LiteralDecimal, Value: "1.5"}}, false},
		{"number range null bound", FieldKindNumber, ValidatorSpec{Kind: ValidatorNumberRange, MinValue: &LiteralSpec{Kind: LiteralNull}, MaxValue: &LiteralSpec{Kind: LiteralDecimal, Value: "2"}}, true},
		{"string regex", FieldKindString, ValidatorSpec{Kind: ValidatorStringRegex, Pattern: "^[a-z]+$"}, false},
		{"unknown kind", FieldKindString, ValidatorSpec{Kind: "anything"}, true},
		{"wrong field type", FieldKindBool, ValidatorSpec{Kind: ValidatorStringLength, Min: &min}, true},
		{"bad regex", FieldKindString, ValidatorSpec{Kind: ValidatorStringRegex, Pattern: "["}, true},
		{"one of untyped value", FieldKindString, ValidatorSpec{Kind: ValidatorOneOf, Values: []LiteralSpec{{Kind: LiteralInt64, Value: "1"}}}, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			spec := validIPv4ListSpec()
			spec.Fields[1].Kind = test.kind
			spec.Fields[1].Validators = []ValidatorSpec{test.validator}
			err := spec.Validate()
			if (err != nil) != test.wantError {
				t.Fatalf("Validate() error = %v, want error=%v", err, test.wantError)
			}
		})
	}
}

func TestVersionRangeUsesNumericComparison(t *testing.T) {
	rangeSpec := VersionRange{MinInclusive: APIVersion{Major: 6, Minor: 6}, MaxExclusive: APIVersion{Major: 6, Minor: 11}}
	if !rangeSpec.Contains(APIVersion{Major: 6, Minor: 10}) {
		t.Fatal("6.10 should be inside 6.6-6.11")
	}
	if rangeSpec.Contains(APIVersion{Major: 6, Minor: 11}) {
		t.Fatal("exclusive maximum must not be included")
	}
}

func validIPv4ListSpec() ResourceSpec {
	versionRange := VersionRange{MinInclusive: APIVersion{Major: 6, Minor: 6}, MaxExclusive: APIVersion{Major: 6, Minor: 7}}
	return ResourceSpec{
		TerraformType: "verity_ipv4_list",
		Description:   "Manages a Verity IPv4 List Filter",
		Modes:         []Mode{ModeDatacenter},
		Versions:      versionRange,
		IdentityPath:  "name",
		API: APIResourceSpec{
			EndpointPath: "/ipv4lists", BulkKey: "ipv4_list", RequestWrapperKey: "ipv4_list_filter",
			ResponseCollectionKey: "ipv4_list_filter", DeleteParameter: "ipv4_list_filter_name", CacheKey: "ipv4_lists",
		},
		Operations: OperationSpec{Create: true, Read: true, Update: true, Delete: true},
		Fields: []FieldSpec{
			{
				TerraformName: "name", APIName: "name", Kind: FieldKindString, Access: AccessRequired, Description: "Name",
				Replace: true, Modes: []Mode{ModeDatacenter}, Versions: versionRange,
				ResponseAbsence: ResponseAbsenceError, CreateNull: CreateNullReject, UpdateClear: UpdateClearReject, UnknownPlan: UnknownPlanReject, StateOwnership: StateConfiguration,
			},
			{
				TerraformName: "enable", APIName: "enable", Kind: FieldKindBool, Access: AccessOptionalComputed, Description: "Enable",
				Modes: []Mode{ModeDatacenter}, Versions: versionRange,
				ResponseAbsence: ResponseAbsenceTerraformNull, CreateNull: CreateNullOmit, UpdateClear: UpdateClearAPINull, UnknownPlan: UnknownPlanPreserve, StateOwnership: StateConfigurationOrServer,
			},
			{
				TerraformName: "ipv4_list", APIName: "ipv4_list", Kind: FieldKindString, Access: AccessOptionalComputed, Description: "IPv4 addresses",
				Modes: []Mode{ModeDatacenter}, Versions: versionRange,
				ResponseAbsence: ResponseAbsenceTerraformNull, CreateNull: CreateNullOmit, UpdateClear: UpdateClearAPINull, UnknownPlan: UnknownPlanPreserve, StateOwnership: StateConfigurationOrServer,
			},
		},
	}
}

func nestedField(collection *CollectionSpec) FieldSpec {
	versionRange := VersionRange{MinInclusive: APIVersion{Major: 6, Minor: 6}, MaxExclusive: APIVersion{Major: 6, Minor: 7}}
	return FieldSpec{
		TerraformName: "nested", APIName: "nested", Kind: FieldKindList, ElementKind: FieldKindObject, Access: AccessOptional, Description: "Nested", Modes: []Mode{ModeDatacenter}, Versions: versionRange,
		ResponseAbsence: ResponseAbsenceTerraformNull, CreateNull: CreateNullOmit, UpdateClear: UpdateClearAPINull, UnknownPlan: UnknownPlanReject, StateOwnership: StateConfiguration,
		Collection: collection,
		Fields: []FieldSpec{{
			TerraformName: "value", APIName: "value", Kind: FieldKindString, Access: AccessOptional, Description: "Value", Modes: []Mode{ModeDatacenter}, Versions: versionRange,
			ResponseAbsence: ResponseAbsenceTerraformNull, CreateNull: CreateNullOmit, UpdateClear: UpdateClearAPINull, UnknownPlan: UnknownPlanReject, StateOwnership: StateConfiguration,
		}},
	}
}

// unmanagedField builds a valid unmanaged field that the caller then breaks in
// exactly one way, so each table case isolates a single validation rule.
func unmanagedField(breakIt func(*FieldSpec)) FieldSpec {
	field := FieldSpec{
		APIName:   "object_properties",
		Kind:      FieldKindObject,
		Unmanaged: true,
		Modes:     []Mode{ModeDatacenter},
		Versions:  VersionRange{MinInclusive: APIVersion{Major: 6, Minor: 6}, MaxExclusive: APIVersion{Major: 6, Minor: 7}},
	}
	breakIt(&field)
	return field
}
