package transport

import (
	"fmt"

	"terraform-provider-verity/openapi"
)

type IPv4ListAdapter struct{}

func (IPv4ListAdapter) BuildPut(resources map[string]WireObject) (*openapi.Ipv4listsPutRequest, error) {
	values, err := ipv4ListValues(resources, true)
	if err != nil {
		return nil, err
	}
	request := openapi.NewIpv4listsPutRequest()
	request.SetIpv4ListFilter(values)
	return request, nil
}

func (IPv4ListAdapter) BuildPatch(resources map[string]WireObject) (*openapi.Ipv4listsPutRequest, error) {
	values, err := ipv4ListValues(resources, false)
	if err != nil {
		return nil, err
	}
	request := openapi.NewIpv4listsPutRequest()
	request.SetIpv4ListFilter(values)
	return request, nil
}

func ipv4ListValues(resources map[string]WireObject, requireName bool) (map[string]openapi.Ipv4listsPutRequestIpv4ListFilterValue, error) {
	values := make(map[string]openapi.Ipv4listsPutRequestIpv4ListFilterValue, len(resources))
	for resourceName, object := range resources {
		value := openapi.Ipv4listsPutRequestIpv4ListFilterValue{}
		name, err := wireString(object, "name", requireName)
		if err != nil {
			return nil, fmt.Errorf("ipv4 list %q: %w", resourceName, err)
		}
		value.Name = name
		enable, err := wireBool(object, "enable")
		if err != nil {
			return nil, fmt.Errorf("ipv4 list %q: %w", resourceName, err)
		}
		value.Enable = enable
		ipv4List, err := wireString(object, "ipv4_list", false)
		if err != nil {
			return nil, fmt.Errorf("ipv4 list %q: %w", resourceName, err)
		}
		value.Ipv4List = ipv4List
		values[resourceName] = value
	}
	return values, nil
}

func wireString(object WireObject, field string, required bool) (*string, error) {
	value, exists := object[field]
	if !exists {
		if required {
			return nil, fmt.Errorf("missing required field %q", field)
		}
		return nil, nil
	}
	if value.Kind() == ValueKindNull {
		return nil, fmt.Errorf("field %q is explicit null, but the generated IPv4 List SDK field cannot represent null", field)
	}
	if value.Kind() != ValueKindString {
		return nil, fmt.Errorf("field %q must be string, got %s", field, value.Kind())
	}
	result := value.StringValue()
	return &result, nil
}

func wireBool(object WireObject, field string) (*bool, error) {
	value, exists := object[field]
	if !exists {
		return nil, nil
	}
	if value.Kind() == ValueKindNull {
		return nil, fmt.Errorf("field %q is explicit null, but the generated IPv4 List SDK field cannot represent null", field)
	}
	if value.Kind() != ValueKindBool {
		return nil, fmt.Errorf("field %q must be bool, got %s", field, value.Kind())
	}
	result := value.BoolValue()
	return &result, nil
}
