package main

import (
	"bytes"
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"text/template"

	"terraform-provider-verity/internal/spec"
)

//go:embed templates/*.tmpl
var docTemplates embed.FS

type docsOptions struct {
	Registry  string
	OutputDir string
	Keep      string
	Check     bool
}

type docArgument struct {
	Name        string
	Type        string
	Description string
	Notes       []string
	Children    []docArgument
}

type docReference struct {
	Field        string
	TypeField    string
	AllowedTypes []string
	Path         string
}

type docAutoAssignment struct {
	Field          string
	Flag           string
	RecomputedWhen []string
	Path           string
}

type docExample struct {
	Mode string
	HCL  string
}

type docResource struct {
	TerraformType string
	Description   string
	Modes         []string
	UpdateOnly    bool
	Identity      string
	Examples      []docExample
	Required      []docArgument
	Optional      []docArgument
	ReadOnly      []docArgument
	References    []docReference
	AutoAssigned  []docAutoAssignment
}

func generateDocs(opts docsOptions) error {
	if opts.Registry == "" || opts.OutputDir == "" {
		return fmt.Errorf("--registry and --output-dir are required")
	}
	registry, err := readRegistry(opts.Registry)
	if err != nil {
		return err
	}
	templates, err := template.New("docs").Funcs(template.FuncMap{
		"join":  strings.Join,
		"code":  codeList,
		"title": modeTitle,
		"dict":  templateDict,
	}).ParseFS(docTemplates, "templates/*.tmpl")
	if err != nil {
		return fmt.Errorf("parse doc templates: %w", err)
	}

	outputs := map[string][]byte{}
	for _, resource := range registry {
		var buf bytes.Buffer
		if err := templates.ExecuteTemplate(&buf, "resource.md.tmpl", buildDocResource(resource)); err != nil {
			return fmt.Errorf("render %s: %w", resource.TerraformType, err)
		}
		outputs[filepath.Join(opts.OutputDir, resource.TerraformType+".md")] = tidyMarkdown(buf.Bytes())
	}
	keep := map[string]bool{}
	for _, name := range strings.Split(opts.Keep, ",") {
		if name = strings.TrimSpace(name); name != "" {
			keep[name] = true
		}
	}
	stale, err := staleDocs(opts.OutputDir, outputs, keep)
	if err != nil {
		return err
	}

	paths := make([]string, 0, len(outputs))
	for path := range outputs {
		paths = append(paths, path)
	}
	sort.Strings(paths)

	if opts.Check {
		var problems []string
		for _, path := range paths {
			existing, err := os.ReadFile(path)
			if err != nil || !bytes.Equal(existing, outputs[path]) {
				problems = append(problems, path+" differs")
			}
		}
		for _, path := range stale {
			problems = append(problems, path+" is not generated from the registry")
		}
		if len(problems) > 0 {
			return fmt.Errorf("%s: rerun specgen docs --registry %s", strings.Join(problems, "; "), opts.Registry)
		}
		return nil
	}
	for _, path := range stale {
		if err := os.Remove(path); err != nil {
			return fmt.Errorf("remove stale %s: %w", path, err)
		}
	}
	for _, path := range paths {
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return fmt.Errorf("create %s: %w", filepath.Dir(path), err)
		}
		if err := writeFile(path, outputs[path]); err != nil {
			return err
		}
	}
	return nil
}

func staleDocs(dir string, outputs map[string][]byte, keep map[string]bool) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("read %s: %w", dir, err)
	}
	var stale []string
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") || keep[entry.Name()] {
			continue
		}
		path := filepath.Join(dir, entry.Name())
		if _, generated := outputs[path]; !generated {
			stale = append(stale, path)
		}
	}
	sort.Strings(stale)
	return stale, nil
}

func buildDocResource(resource spec.ResourceSpec) docResource {
	doc := docResource{
		TerraformType: resource.TerraformType,
		Description:   docText(resource.Description),
		UpdateOnly:    !resource.Operations.Create,
		Identity:      resource.IdentityPath,
	}
	for _, mode := range resource.Modes {
		doc.Modes = append(doc.Modes, string(mode))
	}
	for _, field := range resource.Fields {
		if field.Unmanaged {
			continue
		}
		argument := docArgumentFor(field, resource.Modes)
		switch {
		case field.Access == spec.AccessRequired:
			doc.Required = append(doc.Required, argument)
		case field.Access == spec.AccessComputed:
			doc.ReadOnly = append(doc.ReadOnly, argument)
		default:
			doc.Optional = append(doc.Optional, argument)
		}
	}
	collectDocPairs(resource.Fields, "", &doc)
	for _, mode := range resource.Modes {
		doc.Examples = append(doc.Examples, docExample{Mode: modeTitle(string(mode)), HCL: renderExample(resource, mode)})
	}
	if len(doc.Examples) > 0 {
		same := true
		for _, example := range doc.Examples[1:] {
			same = same && example.HCL == doc.Examples[0].HCL
		}
		if same {
			doc.Examples = []docExample{{HCL: doc.Examples[0].HCL}}
		}
	}
	return doc
}

func docArgumentFor(field spec.FieldSpec, resourceModes []spec.Mode) docArgument {
	argument := docArgument{
		Name:        field.TerraformName,
		Type:        docType(field),
		Description: docText(field.Description),
	}
	if field.Collection != nil && field.Collection.Strategy == spec.CollectionSingleton {
		argument.Notes = append(argument.Notes, "At most one block.")
	}
	if field.Collection != nil && field.Collection.IdentityField != "" {
		argument.Notes = append(argument.Notes, fmt.Sprintf("Entries are matched by `%s`.", field.Collection.IdentityField))
	}
	if only := modeOnly(field.Modes, resourceModes); only != "" {
		argument.Notes = append(argument.Notes, only)
	}
	if field.Replace {
		argument.Notes = append(argument.Notes, "Changing it replaces the resource.")
	}
	if field.Sensitive {
		argument.Notes = append(argument.Notes, "Sensitive.")
	}
	if field.Nullable {
		argument.Notes = append(argument.Notes, "Set it to `null` to clear it.")
	}
	if field.AutoAssignment != nil {
		argument.Notes = append(argument.Notes, fmt.Sprintf("Assigned by the server while `%s` is `true`; it cannot be set then.", field.AutoAssignment.FlagField))
	}
	if field.Reference != nil {
		argument.Notes = append(argument.Notes, fmt.Sprintf("Set together with `%s`.", field.Reference.TypeField))
	}
	for _, member := range field.Fields {
		if member.Unmanaged {
			continue
		}
		argument.Children = append(argument.Children, docArgumentFor(member, field.Modes))
	}
	return argument
}

func collectDocPairs(fields []spec.FieldSpec, prefix string, doc *docResource) {
	for _, field := range fields {
		if field.Unmanaged {
			continue
		}
		if field.Reference != nil {
			doc.References = append(doc.References, docReference{
				Field:        field.TerraformName,
				TypeField:    field.Reference.TypeField,
				AllowedTypes: field.Reference.AllowedTypes,
				Path:         strings.TrimSuffix(prefix, "."),
			})
		}
		if field.AutoAssignment != nil {
			doc.AutoAssigned = append(doc.AutoAssigned, docAutoAssignment{
				Field:          field.TerraformName,
				Flag:           field.AutoAssignment.FlagField,
				RecomputedWhen: field.AutoAssignment.RecomputedWhen,
				Path:           strings.TrimSuffix(prefix, "."),
			})
		}
		if len(field.Fields) > 0 {
			collectDocPairs(field.Fields, prefix+field.TerraformName+".", doc)
		}
	}
}

func docType(field spec.FieldSpec) string {
	switch field.Kind {
	case spec.FieldKindString:
		return "String"
	case spec.FieldKindBool:
		return "Boolean"
	case spec.FieldKindInt64:
		return "Integer"
	case spec.FieldKindNumber:
		return "Number"
	case spec.FieldKindObject:
		return "Block"
	case spec.FieldKindList:
		if field.ElementKind == spec.FieldKindObject || len(field.Fields) > 0 {
			return "Block List"
		}
		return "List of " + docType(spec.FieldSpec{Kind: field.ElementKind})
	default:
		return string(field.Kind)
	}
}

func modeOnly(fieldModes, resourceModes []spec.Mode) string {
	if len(fieldModes) == 0 || len(fieldModes) >= len(resourceModes) {
		return ""
	}
	names := make([]string, 0, len(fieldModes))
	for _, mode := range fieldModes {
		names = append(names, modeTitle(string(mode)))
	}
	return strings.Join(names, " and ") + " mode only."
}

func modeTitle(mode string) string {
	if mode == "" {
		return mode
	}
	return strings.ToUpper(mode[:1]) + mode[1:]
}

func codeList(values []string) string {
	quoted := make([]string, 0, len(values))
	for _, value := range values {
		quoted = append(quoted, "`"+value+"`")
	}
	return strings.Join(quoted, ", ")
}

func renderExample(resource spec.ResourceSpec, mode spec.Mode) string {
	var buf strings.Builder
	fmt.Fprintf(&buf, "resource %q \"example\" {\n", resource.TerraformType)
	writeExampleFields(&buf, resource.Fields, mode, "  ", resource.IdentityPath)
	buf.WriteString("}\n")
	return buf.String()
}

func appliesTo(field spec.FieldSpec, mode spec.Mode) bool {
	if len(field.Modes) == 0 {
		return true
	}
	for _, candidate := range field.Modes {
		if candidate == mode {
			return true
		}
	}
	return false
}

func writeExampleFields(buf *strings.Builder, fields []spec.FieldSpec, mode spec.Mode, indent, identity string) {
	assigned := map[string]bool{}
	for _, field := range fields {
		if field.AutoAssignment != nil {
			assigned[field.TerraformName] = true
		}
	}
	var blocks []spec.FieldSpec
	ordered := make([]spec.FieldSpec, 0, len(fields))
	for _, field := range fields {
		if field.TerraformName == identity {
			ordered = append([]spec.FieldSpec{field}, ordered...)
			continue
		}
		ordered = append(ordered, field)
	}
	for _, field := range ordered {
		if field.Unmanaged || field.Access == spec.AccessComputed || assigned[field.TerraformName] {
			continue
		}
		if !appliesTo(field, mode) {
			continue
		}
		if field.Kind == spec.FieldKindObject || (field.Kind == spec.FieldKindList && len(field.Fields) > 0) {
			blocks = append(blocks, field)
			continue
		}
		fmt.Fprintf(buf, "%s%s = %s\n", indent, field.TerraformName, exampleValue(field, fields, identity))
	}
	for _, block := range blocks {
		var body strings.Builder
		blockIdentity := ""
		if block.Collection != nil {
			blockIdentity = block.Collection.IdentityField
		}
		writeExampleFields(&body, block.Fields, mode, indent+"  ", blockIdentity)
		if body.Len() == 0 && len(block.Fields) > 0 {
			continue
		}
		fmt.Fprintf(buf, "\n%s%s {\n%s%s}\n", indent, block.TerraformName, body.String(), indent)
	}
}

func exampleValue(field spec.FieldSpec, siblings []spec.FieldSpec, identity string) string {
	if field.TerraformName == identity && field.Kind == spec.FieldKindString {
		return `"example"`
	}
	for _, sibling := range siblings {
		if sibling.AutoAssignment != nil && sibling.AutoAssignment.FlagField == field.TerraformName {
			return "true"
		}
		if sibling.Reference != nil && sibling.Reference.TypeField == field.TerraformName && len(sibling.Reference.AllowedTypes) == 1 {
			return fmt.Sprintf("%q", sibling.Reference.AllowedTypes[0])
		}
	}
	if field.Collection != nil || field.TerraformName == "index" {
		return "1"
	}
	switch field.Kind {
	case spec.FieldKindBool:
		return "false"
	case spec.FieldKindInt64, spec.FieldKindNumber:
		if field.Nullable {
			return "null"
		}
		return "0"
	case spec.FieldKindList:
		return "[]"
	default:
		return `""`
	}
}

func templateDict(pairs ...interface{}) (map[string]interface{}, error) {
	if len(pairs)%2 != 0 {
		return nil, fmt.Errorf("dict needs key and value pairs")
	}
	values := make(map[string]interface{}, len(pairs)/2)
	for i := 0; i < len(pairs); i += 2 {
		key, ok := pairs[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict key %v is not a string", pairs[i])
		}
		values[key] = pairs[i+1]
	}
	return values, nil
}

var docMarkup = strings.NewReplacer(
	"<br>", " ", "<br/>", " ", "<br />", " ",
	"</li><li>", "; ", "<ul>", " ", "</ul>", " ", "<li>", "", "</li>", "",
	"<b>", "", "</b>", "", "<i>", "", "</i>", "",
)

func docText(text string) string {
	text = strings.Join(strings.Fields(docMarkup.Replace(text)), " ")
	text = strings.ReplaceAll(text, "; * ", "; ")
	text = strings.ReplaceAll(text, " * \"", " \"")
	if text == "" {
		return ""
	}
	if last := text[len(text)-1]; last != '.' && last != '?' && last != '!' && last != ':' {
		text += "."
	}
	return text
}

var extraBlankLines = regexp.MustCompile(`\n{3,}`)

func tidyMarkdown(content []byte) []byte {
	return append(bytes.TrimRight(extraBlankLines.ReplaceAll(content, []byte("\n\n")), "\n"), '\n')
}
