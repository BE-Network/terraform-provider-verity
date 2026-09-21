package provider

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"

	"terraform-provider-verity/internal/importer"
)

// importerSupportedFields returns the schema of every resource the provider
// registers, in the shape the importer checks generated configuration against.
// It reads the schemas the provider actually serves, so the check always agrees
// with what Terraform will accept.
func importerSupportedFields(ctx context.Context) map[string]*importer.SchemaFields {
	fields := make(map[string]*importer.SchemaFields)
	for _, factory := range getAllResources() {
		r := factory()
		var metadata resource.MetadataResponse
		r.Metadata(ctx, resource.MetadataRequest{ProviderTypeName: "verity"}, &metadata)
		var schemaResp resource.SchemaResponse
		r.Schema(ctx, resource.SchemaRequest{}, &schemaResp)
		if schemaResp.Diagnostics.HasError() {
			continue
		}
		fields[metadata.TypeName] = schemaFields(schemaResp.Schema.Attributes, schemaResp.Schema.Blocks)
	}
	return fields
}

func schemaFields(attributes map[string]schema.Attribute, blocks map[string]schema.Block) *importer.SchemaFields {
	fields := &importer.SchemaFields{
		Attributes: make(map[string]bool, len(attributes)),
		Blocks:     make(map[string]*importer.SchemaFields, len(blocks)),
	}
	for name := range attributes {
		fields.Attributes[name] = true
	}
	for name, block := range blocks {
		switch b := block.(type) {
		case schema.ListNestedBlock:
			fields.Blocks[name] = schemaFields(b.NestedObject.Attributes, b.NestedObject.Blocks)
		case schema.SetNestedBlock:
			fields.Blocks[name] = schemaFields(b.NestedObject.Attributes, b.NestedObject.Blocks)
		case schema.SingleNestedBlock:
			fields.Blocks[name] = schemaFields(b.Attributes, b.Blocks)
		default:
			// A block kind the provider does not use: accept it whole rather
			// than strip arguments that may be valid.
			fields.Attributes[name] = true
		}
	}
	return fields
}

// unsupportedFieldsWarning describes what the importer left out, or returns ""
// when it left out nothing.
func unsupportedFieldsWarning(unsupported map[string][]string) string {
	if len(unsupported) == 0 {
		return ""
	}
	types := make([]string, 0, len(unsupported))
	for resourceType := range unsupported {
		types = append(types, resourceType)
	}
	sort.Strings(types)

	var detail strings.Builder
	detail.WriteString("The Verity API returned arguments that this provider version does not support. ")
	detail.WriteString("They were left out of the generated configuration, so Terraform will not manage them:\n")
	for _, resourceType := range types {
		fmt.Fprintf(&detail, "\n  %s: %s", resourceType, strings.Join(unsupported[resourceType], ", "))
	}
	detail.WriteString("\n\nPlease check for a newer provider version that supports them.")
	return detail.String()
}

// unsupportedArgumentsFile holds the warning beside the generated configuration.
// Terraform prints the warning when the importer runs, but the import scripts go
// on to a second apply whose output buries it; they print this file at the end
// instead.
const unsupportedArgumentsFile = "unsupported_arguments.txt"

// writeUnsupportedArgumentsFile writes the warning to the output directory, or,
// when nothing was left out, removes the file an earlier import left, so it
// never describes a different system or provider version.
func writeUnsupportedArgumentsFile(dir, warning string) error {
	path := filepath.Join(dir, unsupportedArgumentsFile)
	if warning == "" {
		if err := os.Remove(path); err != nil && !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		return nil
	}
	return os.WriteFile(path, []byte(warning+"\n"), 0644)
}
