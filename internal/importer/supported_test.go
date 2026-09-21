package importer

import (
	"reflect"
	"testing"
)

// A newer Verity system returns arguments the provider's schema does not have.
// The importer leaves each one out, at the top level, in object_properties, in
// list entries, and in a list inside object_properties, and records its path.
// Everything the schema has is kept, including keys the importer never writes
// as arguments and keys it renames.
func TestPruneUnsupportedLeavesOutUnknownArguments(t *testing.T) {
	leaf := func(names ...string) *SchemaFields {
		fields := &SchemaFields{Attributes: map[string]bool{}, Blocks: map[string]*SchemaFields{}}
		for _, name := range names {
			fields.Attributes[name] = true
		}
		return fields
	}
	objectProperties := leaf("user_notes")
	objectProperties.Blocks["system_graphs"] = leaf("index", "graph_type")
	tenant := leaf("enable", "vrf_name")
	tenant.Blocks["object_properties"] = objectProperties
	tenant.Blocks["route_tenants"] = leaf("index", "tenant")

	imp := (&Importer{}).WithSupportedFields(map[string]*SchemaFields{"verity_tenant": tenant})
	objects := map[string]map[string]interface{}{
		"t1": {
			"name":                    "t1",
			"enable":                  true,
			"vrf_name":                "",
			"maximum_ebgp_paths":      nil,
			"maximum_ebgp_paths_mode": "automated",
			"object_properties": map[string]interface{}{
				"user_notes": "",
				"new_note":   "",
				"system_graphs": []interface{}{
					map[string]interface{}{"index": float64(1), "graph_type": "", "graph_color": ""},
				},
			},
			"route_tenants": []interface{}{
				map[string]interface{}{"index": float64(1), "tenant": "", "tenant_ref_type_": "tenant"},
			},
		},
		"t2": {"name": "t2", "enable": false, "maximum_ebgp_paths": float64(4)},
	}
	imp.PruneUnsupported("verity_tenant", objects)

	want := map[string]map[string]interface{}{
		"t1": {
			"name":     "t1",
			"enable":   true,
			"vrf_name": "",
			"object_properties": map[string]interface{}{
				"user_notes": "",
				"system_graphs": []interface{}{
					map[string]interface{}{"index": float64(1), "graph_type": ""},
				},
			},
			"route_tenants": []interface{}{
				map[string]interface{}{"index": float64(1), "tenant": ""},
			},
		},
		"t2": {"name": "t2", "enable": false},
	}
	if !reflect.DeepEqual(objects, want) {
		t.Fatalf("pruned objects = %#v\nwant %#v", objects, want)
	}

	wantPaths := map[string][]string{"verity_tenant": {
		"maximum_ebgp_paths",
		"maximum_ebgp_paths_mode",
		"object_properties.new_note",
		"object_properties.system_graphs.graph_color",
		"route_tenants.tenant_ref_type_",
	}}
	if got := imp.UnsupportedFields(); !reflect.DeepEqual(got, wantPaths) {
		t.Fatalf("UnsupportedFields() = %v, want %v", got, wantPaths)
	}
}

// Keys the importer skips or renames are judged by what it writes, not by the
// API's key.
func TestPruneUnsupportedHonorsSkipKeysAndFieldMappings(t *testing.T) {
	// device_voice_settings writes the API's "Codecs" as "codecs".
	voice := &SchemaFields{Attributes: map[string]bool{"codecs": true}, Blocks: map[string]*SchemaFields{}}
	imp := (&Importer{}).WithSupportedFields(map[string]*SchemaFields{"verity_device_voice_settings": voice})
	objects := map[string]map[string]interface{}{"v": {"name": "v", "Codecs": []interface{}{}}}
	imp.PruneUnsupported("verity_device_voice_settings", objects)
	if _, kept := objects["v"]["Codecs"]; !kept {
		t.Fatal("a renamed argument the schema has was left out")
	}
	if got := imp.UnsupportedFields(); len(got) != 0 {
		t.Fatalf("UnsupportedFields() = %v, want nothing", got)
	}
}

// Without schemas the importer writes everything the API returns, as before.
func TestPruneUnsupportedWithoutSchemasKeepsEverything(t *testing.T) {
	imp := &Importer{}
	objects := map[string]map[string]interface{}{"t": {"name": "t", "anything": true}}
	imp.PruneUnsupported("verity_tenant", objects)
	if _, kept := objects["t"]["anything"]; !kept {
		t.Fatal("an argument was left out with no schema to check it against")
	}
}
