package mock

import (
	"encoding/json"
	"fmt"
	"testing"
)

func AssertFieldEquals(t *testing.T, body map[string]interface{}, jsonPath string, expected interface{}) {
	t.Helper()
	val, found := navigatePath(body, jsonPath)
	if !found {
		t.Errorf("field %q not found in body: %s", jsonPath, jsonMarshal(body))
		return
	}

	switch e := expected.(type) {
	case int:
		expected = float64(e)
	case int64:
		expected = float64(e)
	case int32:
		expected = float64(e)
	}

	if fmt.Sprintf("%v", val) != fmt.Sprintf("%v", expected) {
		t.Errorf("field %q = %v (%T), want %v (%T)", jsonPath, val, val, expected, expected)
	}
}

func AssertFieldAbsent(t *testing.T, body map[string]interface{}, jsonPath string) {
	t.Helper()
	_, found := navigatePath(body, jsonPath)
	if found {
		t.Errorf("field %q should be absent but was found in body: %s", jsonPath, jsonMarshal(body))
	}
}

func AssertFieldNull(t *testing.T, body map[string]interface{}, jsonPath string) {
	t.Helper()
	val, found := navigatePathRaw(body, jsonPath)
	if !found {
		t.Errorf("field %q not found (expected null), body: %s", jsonPath, jsonMarshal(body))
		return
	}
	if val != nil {
		t.Errorf("field %q = %v, want null", jsonPath, val)
	}
}

func AssertOnlyFields(t *testing.T, body map[string]interface{}, objectPath string, allowedFields []string) {
	t.Helper()
	val, found := navigatePath(body, objectPath)
	if !found {
		t.Errorf("object %q not found in body", objectPath)
		return
	}

	obj, ok := val.(map[string]interface{})
	if !ok {
		t.Errorf("object %q is not a map, got %T", objectPath, val)
		return
	}

	allowed := make(map[string]bool)
	for _, f := range allowedFields {
		allowed[f] = true
	}

	for key := range obj {
		if !allowed[key] {
			t.Errorf("unexpected field %q in %q (allowed: %v)", key, objectPath, allowedFields)
		}
	}

	for _, f := range allowedFields {
		if _, exists := obj[f]; !exists {
			t.Errorf("expected field %q missing from %q", f, objectPath)
		}
	}
}

func AssertDeleteQueryParams(t *testing.T, req CapturedRequest, paramName string, expectedValues []string) {
	t.Helper()
	values, ok := req.QueryParams[paramName]
	if !ok {
		t.Errorf("query param %q not found in DELETE request to %s", paramName, req.Path)
		return
	}

	if len(values) != len(expectedValues) {
		t.Errorf("query param %q has %d values, want %d: got %v, want %v",
			paramName, len(values), len(expectedValues), values, expectedValues)
		return
	}

	valSet := make(map[string]bool)
	for _, v := range values {
		valSet[v] = true
	}
	for _, ev := range expectedValues {
		if !valSet[ev] {
			t.Errorf("query param %q missing value %q, got %v", paramName, ev, values)
		}
	}
}

func AssertRequestCount(t *testing.T, requests []CapturedRequest, method string, expectedCount int) {
	t.Helper()
	count := 0
	for _, r := range requests {
		if r.Method == method {
			count++
		}
	}
	if count != expectedCount {
		t.Errorf("expected %d %s requests, got %d", expectedCount, method, count)
	}
}

func AssertNoRequests(t *testing.T, requests []CapturedRequest, method string) {
	t.Helper()
	AssertRequestCount(t, requests, method, 0)
}

func navigatePath(obj map[string]interface{}, path string) (interface{}, bool) {
	keys := splitPath(path)
	var current interface{} = obj

	for _, key := range keys {
		m, ok := current.(map[string]interface{})
		if !ok {
			return nil, false
		}
		current, ok = m[key]
		if !ok {
			return nil, false
		}
	}
	return current, true
}

func navigatePathRaw(obj map[string]interface{}, path string) (interface{}, bool) {
	keys := splitPath(path)
	var current interface{} = obj

	for i, key := range keys {
		m, ok := current.(map[string]interface{})
		if !ok {
			return nil, false
		}
		if i == len(keys)-1 {
			val, exists := m[key]
			return val, exists
		}
		current, ok = m[key]
		if !ok {
			return nil, false
		}
	}
	return current, true
}

func splitPath(path string) []string {
	var parts []string
	current := ""
	for _, ch := range path {
		if ch == '.' {
			if current != "" {
				parts = append(parts, current)
			}
			current = ""
		} else {
			current += string(ch)
		}
	}
	if current != "" {
		parts = append(parts, current)
	}
	return parts
}

func jsonMarshal(v interface{}) string {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(data)
}
