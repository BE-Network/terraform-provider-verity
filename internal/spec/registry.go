package spec

import (
	"fmt"
	"sort"
	"strings"
)

type Registry []ResourceSpec

func (r Registry) Validate() error {
	if len(r) == 0 {
		return fmt.Errorf("resource registry is empty")
	}
	terraformTypes := make(map[string]bool, len(r))

	bulkRoutes := make(map[string]string, len(r))
	for _, resource := range r {
		if err := resource.Validate(); err != nil {
			return err
		}
		if terraformTypes[resource.TerraformType] {
			return fmt.Errorf("duplicate Terraform type %q", resource.TerraformType)
		}
		terraformTypes[resource.TerraformType] = true
		route := resource.API.BulkKey + "\x00" + fixedHeaderSignature(resource.API.FixedHeaders)
		if owner, exists := bulkRoutes[route]; exists {
			if len(resource.API.FixedHeaders) == 0 {
				return fmt.Errorf("duplicate bulk key %q shared by %s and %s without a fixed header to discriminate them", resource.API.BulkKey, owner, resource.TerraformType)
			}
			return fmt.Errorf("duplicate bulk key %q and fixed headers shared by %s and %s", resource.API.BulkKey, owner, resource.TerraformType)
		}
		bulkRoutes[route] = resource.TerraformType
	}
	if err := r.validateImportStages(); err != nil {
		return err
	}
	return nil
}

func (r Registry) validateImportStages() error {
	for _, resource := range r {
		for mode := range resource.ImportStages {
			if mode != ModeCampus && mode != ModeDatacenter {
				return fmt.Errorf("%s declares an import stage for unknown mode %q", resource.TerraformType, mode)
			}
			if !resourceSupportsMode(resource, mode) {
				return fmt.Errorf("%s declares an import stage for unsupported mode %s", resource.TerraformType, mode)
			}
		}
	}
	for _, mode := range []Mode{ModeCampus, ModeDatacenter} {
		names := make(map[string]string)
		orders := make(map[int]string)
		expected := 0
		for _, resource := range r {
			if !resourceSupportsMode(resource, mode) {
				continue
			}
			expected++
			stage, declared := resource.ImportStages[mode]
			if !declared {
				return fmt.Errorf("%s supports %s but declares no import stage for it", resource.TerraformType, mode)
			}
			if stage.Name == "" {
				return fmt.Errorf("%s declares an unnamed %s import stage", resource.TerraformType, mode)
			}
			if owner, taken := names[stage.Name]; taken {
				return fmt.Errorf("%s import stage %q is declared by both %s and %s", mode, stage.Name, owner, resource.TerraformType)
			}
			names[stage.Name] = resource.TerraformType
			if owner, taken := orders[stage.Order]; taken {
				return fmt.Errorf("%s import stage order %d is declared by both %s and %s", mode, stage.Order, owner, resource.TerraformType)
			}
			orders[stage.Order] = resource.TerraformType
		}
		for order := 1; order <= expected; order++ {
			if _, declared := orders[order]; !declared {
				return fmt.Errorf("%s import stage order %d is missing; the %d resources in that mode must form one sequence", mode, order, expected)
			}
		}
	}
	return nil
}

func fixedHeaderSignature(headers map[string]string) string {
	if len(headers) == 0 {
		return ""
	}
	pairs := make([]string, 0, len(headers))
	for name, value := range headers {
		pairs = append(pairs, name+"="+value)
	}
	sort.Strings(pairs)
	return strings.Join(pairs, ",")
}

func (r Registry) ForVersion(version APIVersion) Registry {
	selected := make(Registry, 0, len(r))
	for _, resource := range r {
		if resource.Versions.Contains(version) {
			selected = append(selected, resource)
		}
	}
	return selected
}

func resourceSupportsMode(resource ResourceSpec, mode Mode) bool {
	for _, candidate := range resource.Modes {
		if candidate == mode {
			return true
		}
	}
	return false
}
