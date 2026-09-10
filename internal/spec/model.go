package spec

import "fmt"

type Mode string

const (
	ModeDatacenter Mode = "datacenter"
	ModeCampus     Mode = "campus"
)

type APIVersion struct {
	Major int `json:"major"`
	Minor int `json:"minor"`
}

func (v APIVersion) String() string { return fmt.Sprintf("%d.%d", v.Major, v.Minor) }

func (v APIVersion) Compare(other APIVersion) int {
	if v.Major != other.Major {
		if v.Major < other.Major {
			return -1
		}
		return 1
	}
	if v.Minor < other.Minor {
		return -1
	}
	if v.Minor > other.Minor {
		return 1
	}
	return 0
}

type VersionRange struct {
	MinInclusive APIVersion `json:"min_inclusive"`
	MaxExclusive APIVersion `json:"max_exclusive"`
}

func (r VersionRange) Contains(version APIVersion) bool {
	return r.MinInclusive.Compare(version) <= 0 && version.Compare(r.MaxExclusive) < 0
}

type FieldKind string

const (
	FieldKindString FieldKind = "string"
	FieldKindBool   FieldKind = "bool"
	FieldKindInt64  FieldKind = "int64"
	FieldKindNumber FieldKind = "number"
	FieldKindObject FieldKind = "object"
	FieldKindList   FieldKind = "list"
)

type Access string

const (
	AccessRequired         Access = "required"
	AccessOptional         Access = "optional"
	AccessComputed         Access = "computed"
	AccessOptionalComputed Access = "optional_computed"
)

type LiteralKind string

const (
	LiteralNull    LiteralKind = "null"
	LiteralString  LiteralKind = "string"
	LiteralBool    LiteralKind = "bool"
	LiteralInt64   LiteralKind = "int64"
	LiteralDecimal LiteralKind = "decimal"
)

type LiteralSpec struct {
	Kind  LiteralKind `json:"kind"`
	Value string      `json:"value"`
}

type ResponseAbsencePolicy string

const (
	ResponseAbsenceTerraformNull ResponseAbsencePolicy = "terraform_null"
	ResponseAbsenceDefault       ResponseAbsencePolicy = "default"
	ResponseAbsencePreserve      ResponseAbsencePolicy = "preserve"
	ResponseAbsenceError         ResponseAbsencePolicy = "error"
)

type CreateNullPolicy string

const (
	CreateNullOmit    CreateNullPolicy = "omit"
	CreateNullAPINull CreateNullPolicy = "api_null"
	CreateNullDefault CreateNullPolicy = "default"
	CreateNullReject  CreateNullPolicy = "reject"
)

type UpdateClearPolicy string

const (
	UpdateClearAPINull       UpdateClearPolicy = "api_null"
	UpdateClearEmptyString   UpdateClearPolicy = "empty_string"
	UpdateClearZero          UpdateClearPolicy = "zero"
	UpdateClearFalse         UpdateClearPolicy = "false"
	UpdateClearDefault       UpdateClearPolicy = "default"
	UpdateClearOmitUnmanaged UpdateClearPolicy = "omit_unmanaged"
	UpdateClearReject        UpdateClearPolicy = "reject"
)

type UnknownPlanPolicy string

const (
	UnknownPlanOmitAndRead UnknownPlanPolicy = "omit_and_read"
	UnknownPlanPreserve    UnknownPlanPolicy = "preserve_state"
	UnknownPlanReject      UnknownPlanPolicy = "reject"
)

type StateOwnershipPolicy string

const (
	StateConfiguration         StateOwnershipPolicy = "configuration"
	StateServer                StateOwnershipPolicy = "server"
	StateConfigurationOrServer StateOwnershipPolicy = "configuration_or_server"
)

type CollectionStrategy string

const (
	CollectionSingleton             CollectionStrategy = "singleton"
	CollectionIndexedPatch          CollectionStrategy = "indexed_patch"
	CollectionIndexedServerAssigned CollectionStrategy = "indexed_server_assigned"
	CollectionReplace               CollectionStrategy = "replace"
	CollectionComputedSubset        CollectionStrategy = "computed_subset"
)

type CollectionOrdering string

const (
	CollectionOrdered   CollectionOrdering = "ordered"
	CollectionUnordered CollectionOrdering = "unordered"
)

type CollectionSpec struct {
	Strategy      CollectionStrategy `json:"strategy"`
	Ordering      CollectionOrdering `json:"ordering"`
	IdentityField string             `json:"identity_field,omitempty"`
}

type ReferenceSpec struct {
	TypeField    string   `json:"type_field"`
	AllowedTypes []string `json:"allowed_types"`
}

type AutoAssignmentSpec struct {
	FlagField string `json:"flag_field"`
}

type ValidatorKind string

const (
	ValidatorOneOf        ValidatorKind = "one_of"
	ValidatorStringLength ValidatorKind = "string_length"
	ValidatorInt64Range   ValidatorKind = "int64_range"
	ValidatorNumberRange  ValidatorKind = "number_range"
	ValidatorStringRegex  ValidatorKind = "string_regex"
)

type ValidatorSpec struct {
	Kind     ValidatorKind `json:"kind"`
	Values   []LiteralSpec `json:"values,omitempty"`    // one_of
	Min      *int64        `json:"min,omitempty"`       // string_length, int64_range
	Max      *int64        `json:"max,omitempty"`       // string_length, int64_range
	MinValue *LiteralSpec  `json:"min_value,omitempty"` // number_range
	MaxValue *LiteralSpec  `json:"max_value,omitempty"` // number_range
	Pattern  string        `json:"pattern,omitempty"`   // string_regex
}

type FieldSpec struct {
	TerraformName string `json:"terraform_name"`
	APIName       string `json:"api_name"`
	Kind          FieldKind
	Access        Access
	Description   string
	Nullable      bool
	Sensitive     bool
	Replace       bool
	Modes         []Mode
	Versions      VersionRange
	Default       *LiteralSpec
	Validators    []ValidatorSpec

	ResponseAbsence ResponseAbsencePolicy
	CreateNull      CreateNullPolicy
	UpdateClear     UpdateClearPolicy
	UnknownPlan     UnknownPlanPolicy
	StateOwnership  StateOwnershipPolicy

	Fields         []FieldSpec
	Collection     *CollectionSpec
	Reference      *ReferenceSpec
	AutoAssignment *AutoAssignmentSpec
}

type APIResourceSpec struct {
	EndpointPath          string            `json:"endpoint_path"`
	BulkKey               string            `json:"bulk_key"`
	RequestWrapperKey     string            `json:"request_wrapper_key"`
	ResponseCollectionKey string            `json:"response_collection_key"`
	DeleteParameter       string            `json:"delete_parameter"`
	CacheKey              string            `json:"cache_key"`
	FixedHeaders          map[string]string `json:"fixed_headers,omitempty"`
}

type OperationSpec struct {
	Create bool `json:"create"`
	Read   bool `json:"read"`
	Update bool `json:"update"`
	Delete bool `json:"delete"`
}

type DependencySpec struct {
	Before []string `json:"before,omitempty"`
	After  []string `json:"after,omitempty"`
}

type ResourceSpec struct {
	TerraformType string
	Description   string
	Modes         []Mode
	Versions      VersionRange
	IdentityPath  string
	SchemaVersion int64
	API           APIResourceSpec
	Operations    OperationSpec
	Fields        []FieldSpec
	Dependencies  DependencySpec
	Hooks         []string
}
