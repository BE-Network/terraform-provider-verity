package mock

import "testing"

func TestMergeIndexedArrayAssignsAnIndexToEachZero(t *testing.T) {
	t.Parallel()

	existing := []interface{}{
		map[string]interface{}{"index": float64(1), "enable": true},
	}
	patch := []interface{}{
		map[string]interface{}{"index": float64(0), "enable": true},
		map[string]interface{}{"index": float64(0), "enable": false},
	}

	merged := mergeIndexedArray(existing, patch)

	if len(merged) != 3 {
		t.Fatalf("merged to %d entries, want 3: the entry that was already there plus both new ones", len(merged))
	}

	indexes := make(map[float64]bool, len(merged))
	for position, raw := range merged {
		item, ok := raw.(map[string]interface{})
		if !ok {
			t.Fatalf("entry %d is not an object", position)
		}
		index, ok := item["index"].(float64)
		if !ok {
			t.Fatalf("entry %d carries no index", position)
		}
		if index == 0 {
			t.Fatalf("entry %d was stored still carrying index zero; the server assigns a real index", position)
		}
		if indexes[index] {
			t.Fatalf("index %v was assigned twice; each new entry gets its own", index)
		}
		indexes[index] = true
	}

	for _, want := range []float64{1, 2, 3} {
		if !indexes[want] {
			t.Errorf("no entry at index %v; got %v", want, indexes)
		}
	}

	byIndex := make(map[float64]map[string]interface{}, len(merged))
	for _, raw := range merged {
		item := raw.(map[string]interface{})
		byIndex[item["index"].(float64)] = item
	}
	for index, wantEnable := range map[float64]bool{1: true, 2: true, 3: false} {
		if got := byIndex[index]["enable"]; got != wantEnable {
			t.Errorf("index %v enable = %v, want %v", index, got, wantEnable)
		}
	}
}

func TestMergeIndexedArrayKeepsNamedIndexes(t *testing.T) {
	t.Parallel()

	existing := []interface{}{
		map[string]interface{}{"index": float64(1), "enable": true, "filter": "keep"},
	}
	patch := []interface{}{
		map[string]interface{}{"index": float64(1), "enable": false},
		map[string]interface{}{"index": float64(7), "enable": true},
	}

	merged := mergeIndexedArray(existing, patch)

	if len(merged) != 2 {
		t.Fatalf("merged to %d entries, want 2: one updated in place and one added at 7", len(merged))
	}

	got := make(map[float64]map[string]interface{}, len(merged))
	for _, raw := range merged {
		item := raw.(map[string]interface{})
		got[item["index"].(float64)] = item
	}

	updated, present := got[1]
	if !present {
		t.Fatalf("no entry at index 1; got %v", keysOf(got))
	}
	if updated["enable"] != false {
		t.Errorf("index 1 enable = %v, want false: the patch updates it in place", updated["enable"])
	}
	if updated["filter"] != "keep" {
		t.Errorf("index 1 filter = %v, want %q: fields the patch omits are left alone", updated["filter"], "keep")
	}
	if _, present := got[7]; !present {
		t.Errorf("no entry at index 7; a named index is used as given, not reassigned: got %v", keysOf(got))
	}
}

func keysOf(m map[float64]map[string]interface{}) []float64 {
	out := make([]float64, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}

func TestPatchMergesObjectMembers(t *testing.T) {
	t.Parallel()

	ms := NewMockServer("datacenter")
	defer ms.Close()

	path := "/api/lags"
	ms.resourceState = map[string]map[string]map[string]interface{}{
		path: {"lag": {"l": map[string]interface{}{
			"name":              "l",
			"object_properties": map[string]interface{}{"fabric": "fabric-a", "fabric_ref_type_": "fabric"},
		}}},
	}
	ms.applyPatchState(path, map[string]interface{}{
		"lag": map[string]interface{}{"l": map[string]interface{}{
			"object_properties": map[string]interface{}{"fabric": "fabric-b"},
		}},
	}, nil)

	got := ms.resourceState[path]["lag"]["l"].(map[string]interface{})["object_properties"].(map[string]interface{})
	if got["fabric"] != "fabric-b" {
		t.Errorf("fabric = %v, want the patched value", got["fabric"])
	}
	if got["fabric_ref_type_"] != "fabric" {
		t.Errorf("fabric_ref_type_ = %v, want it kept: a member the PATCH did not carry is unchanged", got["fabric_ref_type_"])
	}
}

func TestPatchMergesIndexedArraysInsideObjects(t *testing.T) {
	t.Parallel()

	ms := NewMockServer("datacenter")
	defer ms.Close()

	path := "/api/fabrics"
	ms.resourceState = map[string]map[string]map[string]interface{}{
		path: {"fabric": {"f": map[string]interface{}{
			"name": "f",
			"object_properties": map[string]interface{}{
				"system_graphs": []interface{}{map[string]interface{}{"index": float64(1)}},
			},
		}}},
	}
	ms.applyPatchState(path, map[string]interface{}{
		"fabric": map[string]interface{}{"f": map[string]interface{}{
			"object_properties": map[string]interface{}{
				"system_graphs": []interface{}{map[string]interface{}{"index": float64(2)}},
			},
		}},
	}, nil)

	graphs := ms.resourceState[path]["fabric"]["f"].(map[string]interface{})["object_properties"].(map[string]interface{})["system_graphs"].([]interface{})
	if len(graphs) != 2 {
		t.Fatalf("system_graphs = %v, want entries 1 and 2: adding an entry must not drop the others", graphs)
	}
}

func TestMergeIndexedArrayIndexOnlyEntries(t *testing.T) {
	t.Parallel()

	existing := []interface{}{map[string]interface{}{"index": float64(1)}}
	merged := mergeIndexedArray(existing, []interface{}{
		map[string]interface{}{"index": float64(1)},
		map[string]interface{}{"index": float64(2)},
	})
	if len(merged) != 1 || merged[0].(map[string]interface{})["index"] != float64(2) {
		t.Fatalf("merged = %v, want index 1 deleted and index 2 added", merged)
	}
}
