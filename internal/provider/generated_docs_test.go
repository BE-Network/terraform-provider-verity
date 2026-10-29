package provider

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"

	"terraform-provider-verity/internal/importer"
)

var docExample = regexp.MustCompile("(?s)```hcl\n(.*?)```")

func TestGeneratedDocExamplesMatchTheSchema(t *testing.T) {
	schemas := importerSupportedFields(context.Background())
	pages, err := filepath.Glob(filepath.Join("..", "..", "docs", "resources", "verity_*.md"))
	if err != nil {
		t.Fatal(err)
	}
	checked := 0
	for _, page := range pages {
		raw, err := os.ReadFile(page)
		if err != nil {
			t.Fatal(err)
		}
		for _, match := range docExample.FindAllSubmatch(raw, -1) {
			file, diags := hclsyntax.ParseConfig(match[1], page, hcl.InitialPos)
			if diags.HasErrors() {
				t.Errorf("%s: example does not parse: %s", page, diags.Error())
				continue
			}
			for _, block := range file.Body.(*hclsyntax.Body).Blocks {
				if block.Type != "resource" || len(block.Labels) == 0 {
					continue
				}
				fields, known := schemas[block.Labels[0]]
				if !known {
					if block.Labels[0] != "verity_operation_stage" {
						t.Errorf("%s: example uses unknown resource type %s", page, block.Labels[0])
					}
					continue
				}
				checkExampleBody(t, page, block.Labels[0], block.Body, fields)
				checked++
			}
		}
	}
	if want := len(registryResourceOrder()); checked < want {
		t.Fatalf("checked %d examples, want at least one for each of the %d registry resources", checked, want)
	}
}

func checkExampleBody(t *testing.T, page, path string, body *hclsyntax.Body, fields *importer.SchemaFields) {
	t.Helper()
	for name := range body.Attributes {
		if name == "depends_on" || name == "lifecycle" {
			continue
		}
		if !fields.Attributes[name] {
			t.Errorf("%s: %s.%s is not an argument of the schema", page, path, name)
		}
	}
	for _, block := range body.Blocks {
		nested, known := fields.Blocks[block.Type]
		if !known {
			if !strings.HasPrefix(block.Type, "lifecycle") {
				t.Errorf("%s: %s.%s is not a block of the schema", page, path, block.Type)
			}
			continue
		}
		checkExampleBody(t, page, path+"."+block.Type, block.Body, nested)
	}
}
