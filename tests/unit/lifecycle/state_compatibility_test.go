package lifecycle

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"terraform-provider-verity/internal/registry"
)

const stateCompatibilityPath = "../testdata/state_compatibility_v6_6.json"
const stateCompatibilitySHA256 = "017c8873b354249f4898671ae3cc50cd79ce4f9099832a831af9d97dde69266e"

type stateCompatibilityBaseline struct {
	FormatVersion        int                   `json:"format_version"`
	SourceCommit         string                `json:"source_commit"`
	SourceSnapshotSHA256 string                `json:"source_snapshot_sha256"`
	Resources            []stateRepresentation `json:"resources"`
}

type stateRepresentation struct {
	TerraformType string            `json:"terraform_type"`
	SchemaVersion int64             `json:"schema_version"`
	Fields        map[string]string `json:"fields"`
}

func stateFields(attributes []attributeGolden, blocks []blockSchemaGolden, prefix string) map[string]string {
	fields := make(map[string]string)
	for _, attribute := range attributes {
		fieldType := attribute.Type
		if fieldType == "int64" {
			fieldType = "number"
		}
		if attribute.ElementType != "" {
			fieldType += ":" + attribute.ElementType
		}
		fields[prefix+attribute.Name] = fieldType
	}
	for _, block := range blocks {
		path := prefix + block.Name
		fields[path] = "block:" + block.Nesting
		for child, fieldType := range stateFields(block.Attributes, block.Blocks, path+".") {
			fields[child] = fieldType
		}
	}
	return fields
}

func loadStateCompatibilityBaseline() (stateCompatibilityBaseline, error) {
	var baseline stateCompatibilityBaseline
	raw, err := os.ReadFile(stateCompatibilityPath)
	if err != nil {
		return baseline, err
	}
	if got := fmt.Sprintf("%x", sha256.Sum256(raw)); got != stateCompatibilitySHA256 {
		return baseline, fmt.Errorf("historical state baseline changed: SHA256 %s; preserve the released baseline instead of regenerating it", got)
	}
	if err := json.Unmarshal(raw, &baseline); err != nil {
		return baseline, err
	}
	if baseline.FormatVersion != 1 || len(baseline.Resources) == 0 || baseline.SourceCommit == "" || baseline.SourceSnapshotSHA256 == "" {
		return baseline, fmt.Errorf("invalid historical state baseline")
	}
	return baseline, nil
}

func requireStateCompatibility(t *testing.T, current schemaGolden) {
	t.Helper()
	baseline, err := loadStateCompatibilityBaseline()
	if err != nil {
		t.Fatalf("load state compatibility baseline: %v", err)
	}
	policy, err := loadStateTransitionPolicy("../../../specs/state_transitions.json", "../../..")
	if err != nil {
		t.Fatalf("load state transition policy: %v", err)
	}
	version, err := registry.Version()
	if err != nil {
		t.Fatalf("load selected API version: %v", err)
	}
	if err := checkStateCompatibility(baseline, current, version.String(), policy); err != nil {
		t.Fatal(err)
	}
}

func TestStateCompatibility(t *testing.T) {
	requireStateCompatibility(t, buildSchemaGolden(t))
}

func TestStateCompatibilityGuard(t *testing.T) {
	for _, test := range []struct {
		name   string
		change func(*resourceSchemaGolden)
		want   string
	}{
		{name: "unchanged"},
		{name: "metadata", change: func(r *resourceSchemaGolden) { r.Attributes[0].Sensitive = true; r.Attributes[0].Computed = true }},
		{name: "additive", change: func(r *resourceSchemaGolden) {
			r.Attributes = append(r.Attributes, attributeGolden{Name: "enabled", Type: "bool"})
		}},
		{name: "type", change: func(r *resourceSchemaGolden) { r.Attributes[0].Type = "number" }, want: "undeclared state change"},
		{name: "nested additive", change: func(r *resourceSchemaGolden) {
			r.Blocks[0].Attributes = append(r.Blocks[0].Attributes, attributeGolden{Name: "enabled", Type: "bool"})
		}},
		{name: "removed", change: func(r *resourceSchemaGolden) { r.Attributes = nil }, want: "name removed"},
		{name: "renamed", change: func(r *resourceSchemaGolden) { r.Attributes[0].Name = "label" }, want: "name removed"},
		{name: "nested type", change: func(r *resourceSchemaGolden) { r.Blocks[0].Attributes[0].Type = "bool" }, want: "entries.value changed"},
		{name: "nesting", change: func(r *resourceSchemaGolden) { r.Blocks[0].Nesting = "single" }, want: "entries changed"},
		{name: "collection element", change: func(r *resourceSchemaGolden) { r.Attributes[1].ElementType = "tftypes.Number" }, want: "labels changed"},
		{name: "version only", change: func(r *resourceSchemaGolden) { r.SchemaVersion = 1 }, want: "undeclared state change"},
	} {
		t.Run(test.name, func(t *testing.T) {
			res := resourceSchemaGolden{TerraformType: "verity_example", Attributes: []attributeGolden{{Name: "name", Type: "string"}, {Name: "labels", Type: "list", ElementType: "tftypes.String"}}, Blocks: []blockSchemaGolden{{Name: "entries", Nesting: "list", Attributes: []attributeGolden{{Name: "value", Type: "string"}}}}}
			baseline := stateCompatibilityBaseline{Resources: []stateRepresentation{{TerraformType: res.TerraformType, Fields: stateFields(res.Attributes, res.Blocks, "")}}}
			if test.change != nil {
				test.change(&res)
			}
			err := checkStateCompatibility(baseline, schemaGolden{Resources: []resourceSchemaGolden{res}}, "6.6", stateTransitionPolicy{FormatVersion: 1})
			if test.want == "" {
				if err != nil {
					t.Fatal(err)
				}
			} else if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("error = %v, want %q", err, test.want)
			}
		})
	}
}

func TestStateCompatibilityResourceAndVersionChanges(t *testing.T) {
	baseline := stateCompatibilityBaseline{Resources: []stateRepresentation{{TerraformType: "verity_existing", SchemaVersion: 1}}}
	for _, test := range []struct {
		name    string
		current []resourceSchemaGolden
		want    string
	}{
		{name: "removed resource", want: "resource removed"},
		{name: "version decreased", current: []resourceSchemaGolden{{TerraformType: "verity_existing", SchemaVersion: 0}}, want: "schema version changed"},
		{name: "added resource", current: []resourceSchemaGolden{{TerraformType: "verity_existing", SchemaVersion: 1}, {TerraformType: "verity_new"}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			err := checkStateCompatibility(baseline, schemaGolden{Resources: test.current}, "6.6", stateTransitionPolicy{FormatVersion: 1})
			if test.want == "" {
				if err != nil {
					t.Fatal(err)
				}
			} else if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("error = %v, want %q", err, test.want)
			}
		})
	}
}

func TestStateCompatibilityNumericRepresentation(t *testing.T) {
	baseline := stateCompatibilityBaseline{Resources: []stateRepresentation{{TerraformType: "verity_example", Fields: map[string]string{"index": "number"}}}}
	for _, attributeType := range []string{"int64", "number"} {
		current := schemaGolden{Resources: []resourceSchemaGolden{{TerraformType: "verity_example", Attributes: []attributeGolden{{Name: "index", Type: attributeType}}}}}
		if err := checkStateCompatibility(baseline, current, "6.6", stateTransitionPolicy{FormatVersion: 1}); err != nil {
			t.Fatalf("%s uses the Terraform number representation: %v", attributeType, err)
		}
	}
}
