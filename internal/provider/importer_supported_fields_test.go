package provider

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"terraform-provider-verity/internal/importer"
	"terraform-provider-verity/internal/registry"
)

// The importer leaves out what the provider's schema does not have. That must
// never include an argument the provider supports, or an import that works
// today would silently lose it. The golden request bodies are generated from
// configurations that write every argument of every resource, so checking them
// against the schemas the provider serves, with either engine, must leave out
// nothing.
//
// The mock API responses under tests/unit/testdata/responses are not used: some
// carry keys the 6.6 API does not return, such as a top-level enable on an SFP
// breakout, which the importer is right to leave out.
func TestImporterKeepsEverySupportedArgument(t *testing.T) {
	resources, err := registry.Load()
	if err != nil {
		t.Fatal(err)
	}
	typesByWrapper := map[string][]string{}
	for _, resource := range resources {
		// Requests and responses may wrap objects under different keys, as the
		// ACLs do; the fixtures hold both.
		keys := map[string]bool{resource.API.RequestWrapperKey: true, resource.API.ResponseCollectionKey: true}
		for key := range keys {
			if key != "" {
				typesByWrapper[key] = append(typesByWrapper[key], resource.TerraformType)
			}
		}
	}

	var bodies []string
	for _, pattern := range []string{
		"../../tests/unit/lifecycle/testdata/golden/*/put.json",
		"../../tests/unit/lifecycle/testdata/golden/*/patch.json",
	} {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			t.Fatal(err)
		}
		bodies = append(bodies, matches...)
	}
	if len(bodies) < 50 {
		t.Fatalf("found only %d bodies to check; the test data moved", len(bodies))
	}

	for _, selection := range []string{"", "all"} {
		t.Run("generic="+selection, func(t *testing.T) {
			t.Setenv(GenericResourcesEnvVar, selection)
			imp := importer.NewImporter(nil, "datacenter").WithSupportedFields(importerSupportedFields(context.Background()))
			checked := 0
			for _, path := range bodies {
				raw, err := os.ReadFile(path)
				if err != nil {
					t.Fatal(err)
				}
				var body map[string]map[string]map[string]interface{}
				if err := json.Unmarshal(raw, &body); err != nil {
					t.Fatalf("%s: %v", path, err)
				}
				for wrapper, objects := range body {
					resourceTypes := typesByWrapper[wrapper]
					if len(resourceTypes) == 0 {
						t.Fatalf("%s: no registry resource uses wrapper key %q", path, wrapper)
					}
					for _, resourceType := range resourceTypes {
						imp.PruneUnsupported(resourceType, objects)
						checked++
					}
				}
			}
			if checked == 0 {
				t.Fatal("nothing was checked")
			}
			for resourceType, paths := range imp.UnsupportedFields() {
				t.Errorf("%s: supported arguments were left out: %s", resourceType, strings.Join(paths, ", "))
			}
		})
	}
}

// What the user saw from a newer Verity system: arguments the provider does not
// support, at the top level and inside a list entry, are left out and named in
// one warning.
func TestImporterLeavesOutArgumentsFromANewerAPI(t *testing.T) {
	imp := importer.NewImporter(nil, "datacenter").WithSupportedFields(importerSupportedFields(context.Background()))
	tenants := map[string]map[string]interface{}{
		"t1": {"name": "t1", "enable": true, "maximum_ebgp_paths": nil, "maximum_ebgp_paths_mode": "automated"},
	}
	switchpoints := map[string]map[string]interface{}{
		"Leaf_01": {"name": "Leaf_01", "traffic_mirrors": []interface{}{map[string]interface{}{
			"index":                                       float64(1),
			"traffic_mirror_num_enable":                   true,
			"traffic_mirror_num_monitoring_acl":           "montest",
			"traffic_mirror_num_monitoring_acl_ref_type_": "monitoring_acl",
		}}},
	}
	imp.PruneUnsupported("verity_tenant", tenants)
	imp.PruneUnsupported("verity_switchpoint", switchpoints)

	if _, kept := tenants["t1"]["maximum_ebgp_paths"]; kept {
		t.Error("maximum_ebgp_paths was written")
	}
	if _, kept := tenants["t1"]["enable"]; !kept {
		t.Error("enable was left out")
	}
	entry := switchpoints["Leaf_01"]["traffic_mirrors"].([]interface{})[0].(map[string]interface{})
	if _, kept := entry["traffic_mirror_num_monitoring_acl"]; kept {
		t.Error("traffic_mirror_num_monitoring_acl was written")
	}
	if _, kept := entry["traffic_mirror_num_enable"]; !kept {
		t.Error("traffic_mirror_num_enable was left out")
	}

	want := "The Verity API returned arguments that this provider version does not support. " +
		"They were left out of the generated configuration, so Terraform will not manage them:\n" +
		"\n  verity_switchpoint: traffic_mirrors.traffic_mirror_num_monitoring_acl, traffic_mirrors.traffic_mirror_num_monitoring_acl_ref_type_" +
		"\n  verity_tenant: maximum_ebgp_paths, maximum_ebgp_paths_mode" +
		"\n\nPlease check for a newer provider version that supports them."
	if got := unsupportedFieldsWarning(imp.UnsupportedFields()); got != want {
		t.Errorf("warning =\n%s\nwant\n%s", got, want)
	}
	if got := unsupportedFieldsWarning(nil); got != "" {
		t.Errorf("warning with nothing left out = %q, want none", got)
	}
}

// The file the import scripts print holds the warning after an import that
// left arguments out, and is removed by one that left out nothing, so it never
// describes an earlier system.
func TestUnsupportedArgumentsFileFollowsTheLastImport(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, unsupportedArgumentsFile)

	if err := writeUnsupportedArgumentsFile(dir, "left out"); err != nil {
		t.Fatal(err)
	}
	if written, err := os.ReadFile(path); err != nil || string(written) != "left out\n" {
		t.Fatalf("file = %q, %v; want the warning", written, err)
	}
	if err := writeUnsupportedArgumentsFile(dir, ""); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("the file from the earlier import is still there: %v", err)
	}
	if err := writeUnsupportedArgumentsFile(dir, ""); err != nil {
		t.Fatalf("removing a file that is not there failed: %v", err)
	}
}
