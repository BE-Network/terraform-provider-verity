package mock

import "testing"

// mergeIndexedArray models how the server places entries in an indexed
// collection, and index zero is the part worth testing directly.
//
// The API documents zero as the request to append: "The index identifying the
// object. Zero if you want to add an object to the list." Several new entries
// can therefore arrive in one request all carrying zero, and it is the server
// that tells them apart by giving each the next free index.
//
// This merge used to key every incoming entry by the index it carried, so two
// arriving at zero landed on the same map key and the second overwrote the
// first. That silently modelled a server which cannot accept more than one new
// entry per request, and it is why the index-zero case in
// tests/unit/lifecycle/index_zero_test.go could not be seen at all.
//
// That lifecycle test observes a Terraform error, and the same error would
// appear if the entries were collapsed again, so it cannot pin this behavior on
// its own. These assertions do.
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

	// The entry that was already there keeps its index, and the two new ones take
	// the next free ones after it.
	for _, want := range []float64{1, 2, 3} {
		if !indexes[want] {
			t.Errorf("no entry at index %v; got %v", want, indexes)
		}
	}

	// The assigned indexes must attach to the right entries, not merely exist:
	// collapsing the two and then padding the result would satisfy the counts
	// above but lose one of the bodies.
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

// An entry that names its index is placed at that index, and one that names an
// index already present updates it rather than being added again. Without this,
// assigning indexes to zeros could be widened into assigning them to everything.
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
