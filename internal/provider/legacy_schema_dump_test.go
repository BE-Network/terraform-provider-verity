package provider

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type dumpAttr struct {
	Kind        string `json:"kind"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Optional    bool   `json:"optional"`
	Computed    bool   `json:"computed"`
	Sensitive   bool   `json:"sensitive"`
	Modifiers   int    `json:"modifiers"`
}

type dumpBlock struct {
	Description string               `json:"description"`
	Attributes  map[string]dumpAttr  `json:"attributes"`
	Blocks      map[string]dumpBlock `json:"blocks,omitempty"`
}

// blocksOf recurses: a nested block may itself contain blocks, as Fabric does
// with object_properties.system_graphs.
func blocksOf(blocks map[string]schema.Block) map[string]dumpBlock {
	out := map[string]dumpBlock{}
	for name, block := range blocks {
		list, ok := block.(schema.ListNestedBlock)
		if !ok {
			out[name] = dumpBlock{Description: "UNSUPPORTED_BLOCK_TYPE"}
			continue
		}
		nested := dumpBlock{Description: list.Description, Attributes: map[string]dumpAttr{}}
		for attrName, attr := range list.NestedObject.Attributes {
			nested.Attributes[attrName] = attrOf(attr)
		}
		if len(list.NestedObject.Blocks) != 0 {
			nested.Blocks = blocksOf(list.NestedObject.Blocks)
		}
		out[name] = nested
	}
	return out
}

type dumpResource struct {
	Description string               `json:"description"`
	Attributes  map[string]dumpAttr  `json:"attributes"`
	Blocks      map[string]dumpBlock `json:"blocks"`
}

func attrOf(a schema.Attribute) dumpAttr {
	d := dumpAttr{Description: a.GetDescription(), Required: a.IsRequired(), Optional: a.IsOptional(), Computed: a.IsComputed(), Sensitive: a.IsSensitive()}
	switch v := a.(type) {
	case schema.StringAttribute:
		d.Kind, d.Modifiers = "string", len(v.PlanModifiers)
	case schema.BoolAttribute:
		d.Kind, d.Modifiers = "bool", len(v.PlanModifiers)
	case schema.Int64Attribute:
		d.Kind, d.Modifiers = "int64", len(v.PlanModifiers)
	case schema.NumberAttribute:
		d.Kind, d.Modifiers = "number", len(v.PlanModifiers)
	case schema.ListAttribute:
		d.Kind, d.Modifiers = "list", len(v.PlanModifiers)
	default:
		d.Kind = "UNSUPPORTED"
	}
	return d
}

// TestDumpLegacySchemas is an authoring aid rather than an assertion. Reviewed
// overrides must carry each field's legacy description, access, and replacement
// behavior exactly, and transcribing those by hand is error prone, so this dumps
// the shipped schemas for comparison. It is inert unless DUMP_OUT is set.
func TestDumpLegacySchemas(t *testing.T) {
	if os.Getenv("DUMP_OUT") == "" {
		t.Skip("set DUMP_OUT to dump legacy schemas for override authoring")
	}
	out := map[string]dumpResource{}
	ctx := context.Background()
	for _, ctor := range getAllResources() {
		r := ctor()
		var md resource.MetadataResponse
		r.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "verity"}, &md)
		var resp resource.SchemaResponse
		r.Schema(ctx, resource.SchemaRequest{}, &resp)
		if resp.Diagnostics.HasError() {
			t.Fatalf("%s: %s", md.TypeName, resp.Diagnostics)
		}
		d := dumpResource{Description: resp.Schema.Description, Attributes: map[string]dumpAttr{}, Blocks: map[string]dumpBlock{}}
		for n, a := range resp.Schema.Attributes {
			d.Attributes[n] = attrOf(a)
		}
		d.Blocks = blocksOf(resp.Schema.Blocks)
		out[md.TypeName] = d
	}
	raw, _ := json.MarshalIndent(out, "", " ")
	if err := os.WriteFile(os.Getenv("DUMP_OUT"), raw, 0o644); err != nil {
		t.Fatal(err)
	}
	t.Logf("dumped %d resources", len(out))
}
