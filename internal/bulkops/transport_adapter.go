package bulkops

import "fmt"

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
		return preparer(filteredData)
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
