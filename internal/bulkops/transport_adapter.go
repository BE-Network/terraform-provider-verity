package bulkops

import (
	"fmt"

	"terraform-provider-verity/internal/transport"
)

func (m *Manager) createRequestPreparerWithError(config ResourceConfig, operationType string) func(map[string]interface{}) (interface{}, error) {
	return func(filteredData map[string]interface{}) (interface{}, error) {
		if operationType == "DELETE" {
			names := make([]string, 0, len(filteredData))
			for name := range filteredData {
				names = append(names, name)
			}
			return names, nil
		}
		preparer := config.PreparePut
		if operationType == "PATCH" {
			preparer = config.PreparePatch
		}
		if (operationType != "PUT" && operationType != "PATCH") || preparer == nil {
			return nil, fmt.Errorf("%s has no %s request preparer", config.ResourceType, operationType)
		}
		if config.ResourceType != "ipv4_list" {
			return preparer(filteredData)
		}
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
			return preparer(filteredData)
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

func prepareTypedBulkRequest[Request, Value any](data map[string]interface{}, set func(*Request, map[string]Value)) (interface{}, error) {
	values := make(map[string]Value, len(data))
	for name, raw := range data {
		value, ok := raw.(Value)
		if !ok {
			return nil, fmt.Errorf("resource %q has value type %T, expected %T", name, raw, *new(Value))
		}
		values[name] = value
	}
	request := new(Request)
	set(request, values)
	return request, nil
}
