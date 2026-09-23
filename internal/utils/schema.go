package utils

import (
	"strings"

	"terraform-provider-verity/internal/spec"
)

func FieldAppliesToMode(endpointKey, fieldPath, mode string) bool {
	resource, found := registryResourceByEndpoint(endpointKey)
	if !found {
		return true
	}
	field, found := registryField(resource.Fields, strings.Split(fieldPath, "."))
	if !found {
		return true
	}
	modes := field.Modes
	if len(modes) == 0 {
		modes = resource.Modes
	}
	if len(modes) == 0 {
		return true
	}
	return specModeSet(modes)[mode]
}

func registryField(fields []spec.FieldSpec, path []string) (spec.FieldSpec, bool) {
	if len(path) == 0 {
		return spec.FieldSpec{}, false
	}
	for _, field := range fields {
		if field.TerraformName != path[0] {
			continue
		}
		if len(path) == 1 {
			return field, true
		}
		return registryField(field.Fields, path[1:])
	}
	return spec.FieldSpec{}, false
}
