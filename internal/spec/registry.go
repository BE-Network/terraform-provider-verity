package spec

import "fmt"

type Registry []ResourceSpec

func (r Registry) Validate() error {
	if len(r) == 0 {
		return fmt.Errorf("resource registry is empty")
	}
	terraformTypes := make(map[string]bool, len(r))
	bulkKeys := make(map[string]bool, len(r))
	for _, resource := range r {
		if err := resource.Validate(); err != nil {
			return err
		}
		if terraformTypes[resource.TerraformType] {
			return fmt.Errorf("duplicate Terraform type %q", resource.TerraformType)
		}
		terraformTypes[resource.TerraformType] = true
		if bulkKeys[resource.API.BulkKey] {
			return fmt.Errorf("duplicate bulk key %q", resource.API.BulkKey)
		}
		bulkKeys[resource.API.BulkKey] = true
	}
	return nil
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
