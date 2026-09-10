package bulkops

import (
	"fmt"

	"terraform-provider-verity/internal/transport"
)

func (m *Manager) createRequestPreparerWithError(config ResourceConfig, operationType string) func(map[string]interface{}) (interface{}, error) {
	if config.ResourceType != "ipv4_list" || (operationType != "PUT" && operationType != "PATCH") {
		return nil
	}

	legacyPreparer := m.createRequestPreparer(config, operationType)
	return func(filteredData map[string]interface{}) (interface{}, error) {
		wireResources := make(map[string]transport.WireObject, len(filteredData))
		for name, value := range filteredData {
			wireObject, ok := value.(transport.WireObject)
			if !ok {
				continue
			}
			wireResources[name] = wireObject
		}

		switch len(wireResources) {
		case 0:
			return legacyPreparer(filteredData), nil
		case len(filteredData):
			adapter := transport.IPv4ListAdapter{}
			if operationType == "PUT" {
				return adapter.BuildPut(wireResources)
			}
			return adapter.BuildPatch(wireResources)
		default:
			return nil, fmt.Errorf("ipv4_list %s batch mixes legacy typed values and transport WireObject values", operationType)
		}
	}
}
