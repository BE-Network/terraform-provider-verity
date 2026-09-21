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
