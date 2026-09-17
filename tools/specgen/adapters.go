package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"os"
	"reflect"
	"sort"
	"strings"

	"terraform-provider-verity/internal/genericresource"
	"terraform-provider-verity/internal/spec"
)

type adapterOptions struct {
	Registry   string
	OpenAPIDir string
	Output     string
	Check      bool
}

// goStruct is one generated OpenAPI struct as the generator sees it: its fields
// indexed by the JSON name the API uses, which is the name the registry also
// speaks. Matching on the tag rather than on a transformed identifier is what
// keeps the generator independent of the SDK generator's naming rules.
type goStruct struct {
	Name   string
	Fields map[string]goField
}

type goField struct {
	GoName string
	GoType string
}

// discoverStructs reads the generated SDK package and indexes every struct in
// it. This runs at generation time, so a type or tag that moves is a build
// failure in regenerated code rather than a runtime surprise — which is the
// objection the plan raises against reflecting over these structs live.
func discoverStructs(dir string) (map[string]goStruct, error) {
	fileSet := token.NewFileSet()
	packages, err := parser.ParseDir(fileSet, dir, nil, 0)
	if err != nil {
		return nil, fmt.Errorf("parse %s: %w", dir, err)
	}
	structs := make(map[string]goStruct)
	for _, pkg := range packages {
		for _, file := range pkg.Files {
			for _, decl := range file.Decls {
				genDecl, ok := decl.(*ast.GenDecl)
				if !ok || genDecl.Tok != token.TYPE {
					continue
				}
				for _, s := range genDecl.Specs {
					typeSpec, ok := s.(*ast.TypeSpec)
					if !ok {
						continue
					}
					structType, ok := typeSpec.Type.(*ast.StructType)
					if !ok {
						continue
					}
					structs[typeSpec.Name.Name] = goStruct{
						Name:   typeSpec.Name.Name,
						Fields: structFields(structType),
					}
				}
			}
		}
	}
	if len(structs) == 0 {
		return nil, fmt.Errorf("no structs found in %s", dir)
	}
	return structs, nil
}

func structFields(structType *ast.StructType) map[string]goField {
	fields := make(map[string]goField)
	for _, field := range structType.Fields.List {
		if field.Tag == nil || len(field.Names) != 1 {
			continue
		}
		tag := reflect.StructTag(strings.Trim(field.Tag.Value, "`")).Get("json")
		name := strings.Split(tag, ",")[0]
		if name == "" || name == "-" {
			continue
		}
		fields[name] = goField{GoName: field.Names[0].Name, GoType: typeString(field.Type)}
	}
	return fields
}

func typeString(expr ast.Expr) string {
	var buf bytes.Buffer
	if err := format.Node(&buf, token.NewFileSet(), expr); err != nil {
		return ""
	}
	return buf.String()
}

// resourceAdapter is one resource the generator can emit an adapter for.
type resourceAdapter struct {
	TerraformType string
	GoTypeName    string
	Fields        []adapterField
}

type adapterField struct {
	APIName string
	GoName  string
	Setter  string
	// Nested is set for a singleton object or an indexed list. Its members are
	// converted by a generated function named by Setter, into the struct the SDK
	// declares for the object or for each list entry.
	Nested *nestedAdapter
}

type nestedAdapter struct {
	GoTypeName string
	Fields     []adapterField
	// List is true for an indexed collection, whose Go type is a slice of the
	// entry struct rather than a pointer to it.
	List bool
}

// putRequestTypeName derives the request type from the endpoint the way the SDK
// generator names it. The derivation is checked rather than trusted: the
// resulting struct's wrapper tag must equal the registry's request wrapper key,
// and a resource whose does not is reported and skipped.
func putRequestTypeName(endpointPath string) string {
	trimmed := strings.TrimPrefix(endpointPath, "/")
	if trimmed == "" {
		return ""
	}
	return strings.ToUpper(trimmed[:1]) + trimmed[1:] + "PutRequest"
}

// setterFor maps a generated field's Go type onto the transport helper that
// fills it. An unrecognised type is refused: guessing at a wrapper would
// produce an adapter that compiles and sends the wrong shape.
func setterFor(goType string) (string, bool) {
	switch goType {
	case "*string":
		return "wireStringPtr", true
	case "*bool":
		return "wireBoolPtr", true
	case "*int32":
		return "wireInt32Ptr", true
	case "*int64":
		return "wireInt64Ptr", true
	case "*float32":
		return "wireFloat32Ptr", true
	case "*float64":
		return "wireFloat64Ptr", true
	case "NullableInt32":
		return "wireNullableInt32", true
	case "NullableInt64":
		return "wireNullableInt64", true
	case "NullableFloat32":
		return "wireNullableFloat32", true
	case "NullableFloat64":
		return "wireNullableFloat64", true
	case "NullableString":
		return "wireNullableString", true
	default:
		return "", false
	}
}

// planAdapters decides which resources can have an adapter generated and why the
// rest cannot. A skip is reported rather than silently dropped, so the set the
// engine can serve is visible instead of inferred.
func planAdapters(registry spec.Registry, structs map[string]goStruct) (adapters []resourceAdapter, skipped []string) {
	for _, resource := range registry {
		adapter, reason := planOne(resource, structs)
		if reason != "" {
			skipped = append(skipped, fmt.Sprintf("%s: %s", resource.TerraformType, reason))
			continue
		}
		adapters = append(adapters, adapter)
	}
	sort.Slice(adapters, func(i, j int) bool { return adapters[i].TerraformType < adapters[j].TerraformType })
	sort.Strings(skipped)
	return adapters, skipped
}

func planOne(resource spec.ResourceSpec, structs map[string]goStruct) (resourceAdapter, string) {
	// The engine's own support rule decides servability, so the adapters and the
	// engine cannot disagree about which resources can be migrated.
	if err := genericresource.Supported(resource); err != nil {
		return resourceAdapter{}, err.Error()
	}
	requestName := putRequestTypeName(resource.API.EndpointPath)
	request, found := structs[requestName]
	if !found {
		return resourceAdapter{}, fmt.Sprintf("no %s in the generated SDK", requestName)
	}
	wrapper, found := request.Fields[resource.API.RequestWrapperKey]
	if !found {
		return resourceAdapter{}, fmt.Sprintf("%s carries no %q field", requestName, resource.API.RequestWrapperKey)
	}
	valueTypeName := strings.TrimPrefix(wrapper.GoType, "*map[string]")
	if valueTypeName == wrapper.GoType {
		return resourceAdapter{}, fmt.Sprintf("%s.%s is %s, not a name-keyed map", requestName, wrapper.GoName, wrapper.GoType)
	}
	value, found := structs[valueTypeName]
	if !found {
		return resourceAdapter{}, fmt.Sprintf("no %s in the generated SDK", valueTypeName)
	}

	fields, reason := planFields(resource.TerraformType, resource.Fields, value, structs)
	if reason != "" {
		return resourceAdapter{}, reason
	}
	return resourceAdapter{TerraformType: resource.TerraformType, GoTypeName: valueTypeName, Fields: fields}, ""
}

// planFields maps a level of spec fields onto the struct the SDK declares for
// it, recursing into a singleton's own struct. The engine's support rule has
// already decided what can appear here, so an unexpected shape is reported as a
// mismatch between the SDK and the registry rather than accommodated.
func planFields(terraformType string, fields []spec.FieldSpec, value goStruct, structs map[string]goStruct) ([]adapterField, string) {
	planned := make([]adapterField, 0, len(fields))
	for _, field := range fields {
		if field.Unmanaged {
			continue
		}
		target, found := value.Fields[field.APIName]
		if !found {
			return nil, fmt.Sprintf("%s carries no %q field", value.Name, field.APIName)
		}
		if field.Kind == spec.FieldKindObject || field.Kind == spec.FieldKindList {
			prefix, shape := "*", "a pointer to"
			if field.Kind == spec.FieldKindList {
				prefix, shape = "[]", "a slice of"
			}
			nestedName := strings.TrimPrefix(target.GoType, prefix)
			nested, found := structs[nestedName]
			if nestedName == target.GoType || !found {
				return nil, fmt.Sprintf("%s.%s is %s, not %s a generated struct", value.Name, target.GoName, target.GoType, shape)
			}
			members, reason := planFields(terraformType, field.Fields, nested, structs)
			if reason != "" {
				return nil, reason
			}
			planned = append(planned, adapterField{
				APIName: field.APIName,
				GoName:  target.GoName,
				Setter:  strings.TrimSuffix(adapterTypeName(terraformType), "Adapter") + target.GoName + "Value",
				Nested:  &nestedAdapter{GoTypeName: nestedName, Fields: members, List: field.Kind == spec.FieldKindList},
			})
			continue
		}
		setter, supported := setterFor(target.GoType)
		if !supported {
			return nil, fmt.Sprintf("%s.%s has unsupported type %s", value.Name, target.GoName, target.GoType)
		}
		planned = append(planned, adapterField{APIName: field.APIName, GoName: target.GoName, Setter: setter})
	}
	sort.Slice(planned, func(i, j int) bool { return planned[i].APIName < planned[j].APIName })
	return planned, ""
}

func readRegistry(path string) (spec.Registry, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read registry: %w", err)
	}
	var artifact struct {
		Resources spec.Registry `json:"resources"`
	}
	if err := json.Unmarshal(raw, &artifact); err != nil {
		return nil, fmt.Errorf("decode registry: %w", err)
	}
	if err := artifact.Resources.Validate(); err != nil {
		return nil, fmt.Errorf("validate registry: %w", err)
	}
	return artifact.Resources, nil
}

// generateAdapters emits one transport adapter per resource the scalar engine
// can serve.
//
// The plan puts a generated adapter at the boundary between the codec and the
// bulk manager: the codec owns field names and values, the adapter owns the
// generated SDK types. Writing them by hand would reinstate the per-resource
// duplication the registry exists to remove, and reflecting over the SDK structs
// at runtime would move the mistakes to runtime. Generating from the structs'
// own JSON tags does neither.
func generateAdapters(opts adapterOptions) error {
	if opts.Registry == "" || opts.OpenAPIDir == "" || opts.Output == "" {
		return fmt.Errorf("--registry, --openapi-dir and --output are required")
	}
	registry, err := readRegistry(opts.Registry)
	if err != nil {
		return err
	}
	structs, err := discoverStructs(opts.OpenAPIDir)
	if err != nil {
		return err
	}
	adapters, skipped := planAdapters(registry, structs)
	if len(adapters) == 0 {
		return fmt.Errorf("no resource could be given an adapter; skipped:\n  %s", strings.Join(skipped, "\n  "))
	}

	var buf bytes.Buffer
	buf.WriteString("// Code generated by tools/specgen adapters. DO NOT EDIT.\n")
	buf.WriteString("//\n")
	buf.WriteString("// Source: specs/generated_registry.json and the generated SDK's own JSON tags.\n")
	buf.WriteString("//\n")
	buf.WriteString("// One adapter per resource the scalar engine can serve. Each converts the\n")
	buf.WriteString("// codec's canonical object into the typed value the bulk manager asserts, which\n")
	buf.WriteString("// is the boundary the plan puts between the two.\n")
	buf.WriteString("//\n")
	buf.WriteString("// Resources without an adapter, and why:\n")
	for _, reason := range skipped {
		fmt.Fprintf(&buf, "//   %s\n", reason)
	}
	buf.WriteString("\npackage transport\n\n")
	buf.WriteString("import (\n\t\"fmt\"\n\n\t\"terraform-provider-verity/openapi\"\n)\n\n")

	buf.WriteString("// GeneratedAdapters holds every generated adapter by Terraform type.\n")
	buf.WriteString("var GeneratedAdapters = map[string]ResourceValueAdapter{\n")
	for _, adapter := range adapters {
		fmt.Fprintf(&buf, "\t%q: %s{},\n", adapter.TerraformType, adapterTypeName(adapter.TerraformType))
	}
	buf.WriteString("}\n")

	for _, adapter := range adapters {
		typeName := adapterTypeName(adapter.TerraformType)
		fmt.Fprintf(&buf, "\ntype %s struct{}\n\n", typeName)
		fmt.Fprintf(&buf, "func (%s) ResourceValue(object WireObject) (interface{}, error) {\n", typeName)
		fmt.Fprintf(&buf, "\tvar value openapi.%s\n", adapter.GoTypeName)
		buf.WriteString("\tfor name, wire := range object {\n\t\tswitch name {\n")
		for _, field := range adapter.Fields {
			fmt.Fprintf(&buf, "\t\tcase %q:\n", field.APIName)
			fmt.Fprintf(&buf, "\t\t\tif err := %s(wire, &value.%s); err != nil {\n", field.Setter, field.GoName)
			fmt.Fprintf(&buf, "\t\t\t\treturn nil, fmt.Errorf(\"%%s: %%w\", name, err)\n\t\t\t}\n")
		}
		buf.WriteString("\t\tdefault:\n")
		buf.WriteString("\t\t\t// The codec only emits fields the spec declares, so an unknown one\n")
		buf.WriteString("\t\t\t// means the registry and this adapter were generated apart.\n")
		fmt.Fprintf(&buf, "\t\t\treturn nil, fmt.Errorf(\"%s has no field %%q\", name)\n", adapter.GoTypeName)
		buf.WriteString("\t\t}\n\t}\n\treturn value, nil\n}\n")
		for _, field := range adapter.Fields {
			if field.Nested != nil {
				writeNestedAdapter(&buf, field)
			}
		}
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("format generated adapters: %w", err)
	}
	if opts.Check {
		existing, err := os.ReadFile(opts.Output)
		if err != nil {
			return fmt.Errorf("read %s for check: %w", opts.Output, err)
		}
		if !bytes.Equal(existing, formatted) {
			return fmt.Errorf("%s differs: rerun specgen adapters", opts.Output)
		}
		return nil
	}
	return writeFile(opts.Output, formatted)
}

// writeNestedAdapter emits the conversion for one singleton object: the codec's
// canonical object into the struct the SDK declares for it, set through a
// pointer because the SDK omits an absent object rather than sending it empty.
func writeNestedAdapter(buf *bytes.Buffer, field adapterField) {
	nested := field.Nested
	if nested.List {
		writeListAdapter(buf, field)
		return
	}
	fmt.Fprintf(buf, "\nfunc %s(wire WireValue, target **openapi.%s) error {\n", field.Setter, nested.GoTypeName)
	buf.WriteString("\tmembers, err := wireObject(wire)\n\tif err != nil {\n\t\treturn err\n\t}\n")
	fmt.Fprintf(buf, "\tvar value openapi.%s\n", nested.GoTypeName)
	buf.WriteString("\tfor name, member := range members {\n\t\tswitch name {\n")
	for _, member := range nested.Fields {
		fmt.Fprintf(buf, "\t\tcase %q:\n", member.APIName)
		fmt.Fprintf(buf, "\t\t\tif err := %s(member, &value.%s); err != nil {\n", member.Setter, member.GoName)
		buf.WriteString("\t\t\t\treturn fmt.Errorf(\"%s: %w\", name, err)\n\t\t\t}\n")
	}
	buf.WriteString("\t\tdefault:\n")
	fmt.Fprintf(buf, "\t\t\treturn fmt.Errorf(\"%s has no field %%q\", name)\n", nested.GoTypeName)
	buf.WriteString("\t\t}\n\t}\n\t*target = &value\n\treturn nil\n}\n")
}

// writeListAdapter emits the conversion for one indexed collection: each entry of
// the codec's canonical list into the struct the SDK declares for an entry, in
// the order the codec produced them.
func writeListAdapter(buf *bytes.Buffer, field adapterField) {
	nested := field.Nested
	fmt.Fprintf(buf, "\nfunc %s(wire WireValue, target *[]openapi.%s) error {\n", field.Setter, nested.GoTypeName)
	buf.WriteString("\tentries, err := wireList(wire)\n\tif err != nil {\n\t\treturn err\n\t}\n")
	fmt.Fprintf(buf, "\tvalues := make([]openapi.%s, 0, len(entries))\n", nested.GoTypeName)
	buf.WriteString("\tfor position, entry := range entries {\n")
	buf.WriteString("\t\tmembers, err := wireObject(entry)\n\t\tif err != nil {\n\t\t\treturn fmt.Errorf(\"[%d]: %w\", position, err)\n\t\t}\n")
	fmt.Fprintf(buf, "\t\tvar value openapi.%s\n", nested.GoTypeName)
	buf.WriteString("\t\tfor name, member := range members {\n\t\t\tswitch name {\n")
	for _, member := range nested.Fields {
		fmt.Fprintf(buf, "\t\t\tcase %q:\n", member.APIName)
		fmt.Fprintf(buf, "\t\t\t\tif err := %s(member, &value.%s); err != nil {\n", member.Setter, member.GoName)
		buf.WriteString("\t\t\t\t\treturn fmt.Errorf(\"[%d].%s: %w\", position, name, err)\n\t\t\t\t}\n")
	}
	buf.WriteString("\t\t\tdefault:\n")
	fmt.Fprintf(buf, "\t\t\t\treturn fmt.Errorf(\"%s has no field %%q\", name)\n", nested.GoTypeName)
	buf.WriteString("\t\t\t}\n\t\t}\n\t\tvalues = append(values, value)\n\t}\n\t*target = values\n\treturn nil\n}\n")
}

// adapterTypeName turns verity_ipv4_list into ipv4ListAdapter.
func adapterTypeName(terraformType string) string {
	parts := strings.Split(strings.TrimPrefix(terraformType, "verity_"), "_")
	name := parts[0]
	for _, part := range parts[1:] {
		if part == "" {
			continue
		}
		name += strings.ToUpper(part[:1]) + part[1:]
	}
	return name + "Adapter"
}
