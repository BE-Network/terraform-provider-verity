package spec

var SelectedAPIVersion = APIVersion{Major: 6, Minor: 6}

var SeedRegistry = Registry{ipv4ListSpec()}

func ipv4ListSpec() ResourceSpec {
	versions := VersionRange{MinInclusive: APIVersion{Major: 6, Minor: 6}, MaxExclusive: APIVersion{Major: 6, Minor: 7}}
	policies := func(access Access, ownership StateOwnershipPolicy) FieldSpec {
		return FieldSpec{
			Access: access, Modes: []Mode{ModeDatacenter}, Versions: versions,
			ResponseAbsence: ResponseAbsenceTerraformNull, CreateNull: CreateNullOmit, UpdateClear: UpdateClearAPINull,
			UnknownPlan: UnknownPlanPreserve, StateOwnership: ownership,
		}
	}
	name := policies(AccessRequired, StateConfiguration)
	name.TerraformName, name.APIName, name.Kind, name.Description = "name", "name", FieldKindString, "Object Name. Must be unique."
	name.Replace, name.ResponseAbsence, name.CreateNull, name.UpdateClear, name.UnknownPlan = true, ResponseAbsenceError, CreateNullReject, UpdateClearReject, UnknownPlanReject
	enable := policies(AccessOptionalComputed, StateConfigurationOrServer)
	enable.TerraformName, enable.APIName, enable.Kind, enable.Description = "enable", "enable", FieldKindBool, "Enable object."
	ipv4List := policies(AccessOptionalComputed, StateConfigurationOrServer)
	ipv4List.TerraformName, ipv4List.APIName, ipv4List.Kind, ipv4List.Description = "ipv4_list", "ipv4_list", FieldKindString, "Comma separated list of IPv4 addresses"

	return ResourceSpec{
		TerraformType: "verity_ipv4_list", Description: "Manages a Verity IPv4 List Filter",
		Modes: []Mode{ModeDatacenter}, Versions: versions, IdentityPath: "name", SchemaVersion: 0,
		API: APIResourceSpec{
			EndpointPath: "/ipv4lists", BulkKey: "ipv4_list", RequestWrapperKey: "ipv4_list_filter",
			ResponseCollectionKey: "ipv4_list_filter", DeleteParameter: "ipv4_list_filter_name", CacheKey: "ipv4_lists",
		},
		Operations: OperationSpec{Create: true, Read: true, Update: true, Delete: true},
		Fields:     []FieldSpec{name, enable, ipv4List},
	}
}
