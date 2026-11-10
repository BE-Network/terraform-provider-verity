package lifecycle

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func legacyReference(t *testing.T) map[string]map[string]interface{} {
	t.Helper()
	path := legacyReferencePath(t)
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("no recorded legacy reference for %s: %v", t.Name(), err)
	}
	var captured map[string]map[string]interface{}
	if err := json.Unmarshal(raw, &captured); err != nil {
		t.Fatalf("%s: %v", path, err)
	}
	for _, body := range captured {
		for wrapper, collection := range body {
			objects, _ := collection.(map[string]interface{})
			for _, object := range objects {
				fields, ok := object.(map[string]interface{})
				if !ok {
					continue
				}
				delete(fields, "name")
				for _, removed := range removedRequestFields[wrapper] {
					container, member, nested := strings.Cut(removed, ".")
					if !nested {
						delete(fields, removed)
						continue
					}
					switch value := fields[container].(type) {
					case map[string]interface{}:
						delete(value, member)
					case []interface{}:
						for _, entry := range value {
							if entryFields, ok := entry.(map[string]interface{}); ok {
								delete(entryFields, member)
							}
						}
					}
				}
			}
		}
	}
	return captured
}

var removedRequestFields = map[string][]string{
	"switchpoint":      {"tenant", "tenant_ref_type_", "eths.port_name"},
	"fabric":           {"enable"},
	"route_map_clause": {"object_properties.match_fields_shown"},
}

func legacyReferencePath(t *testing.T) string {
	name := strings.NewReplacer("/", "__", " ", "_", ":", "_", "'", "", "\"", "").Replace(t.Name())
	return filepath.Join("testdata", "legacy_reference", name+"-1.json")
}

func inspectLegacySchema(t *testing.T, entry ResourceCoverageEntry) resourceSchemaInfo {
	t.Helper()
	baseline, err := loadStateCompatibilityBaseline()
	if err != nil {
		t.Fatal(err)
	}
	var paths map[string]string
	for _, resource := range baseline.Resources {
		if resource.TerraformType == entry.TerraformType {
			paths = resource.Fields
			break
		}
	}
	if paths == nil {
		t.Fatalf("no historical schema for %s", entry.TerraformType)
	}
	filterFields := func(fields []fieldInfo, prefix string) []fieldInfo {
		var historical []fieldInfo
		for _, field := range fields {
			if _, existed := paths[prefix+field.Name]; existed {
				historical = append(historical, field)
			}
		}
		return historical
	}
	var filterBlocks func([]blockInfo, string) []blockInfo
	filterBlocks = func(blocks []blockInfo, prefix string) []blockInfo {
		var historical []blockInfo
		for _, block := range blocks {
			path := prefix + block.Name
			if _, existed := paths[path]; !existed {
				continue
			}
			block.Fields = filterFields(block.Fields, path+".")
			block.Blocks = filterBlocks(block.Blocks, path+".")
			historical = append(historical, block)
		}
		return historical
	}
	schema := inspectSchema(entry.Factory)
	schema.Attributes = filterFields(schema.Attributes, "")
	schema.Blocks = filterBlocks(schema.Blocks, "")
	return schema
}
