package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

const extractionFormatVersion = 1

type extractOptions struct {
	InputDir string
	Output   string
	Check    bool
}

type coverageManifest struct {
	FormatVersion int                `json:"format_version"`
	APIVersion    string             `json:"api_version"`
	Resources     []coverageResource `json:"resources"`
}

type coverageResource struct {
	Path                  string                   `json:"path"`
	Modes                 []string                 `json:"modes"`
	Operations            []string                 `json:"operations"`
	RequestWrapperKey     string                   `json:"request_wrapper_key,omitempty"`
	ResponseCollectionKey string                   `json:"response_collection_key,omitempty"`
	DeleteParameters      []coverageQueryParameter `json:"delete_parameters,omitempty"`
	Fields                []coverageField          `json:"fields"`
	ReviewRequired        []string                 `json:"review_required"`
}

type coverageQueryParameter struct {
	Name     string   `json:"name"`
	Required bool     `json:"required"`
	Kind     string   `json:"kind"`
	ItemKind string   `json:"item_kind,omitempty"`
	Modes    []string `json:"modes"`
}

type coverageField struct {
	APIName     string          `json:"api_name"`
	Kind        string          `json:"kind"`
	ItemKind    string          `json:"item_kind,omitempty"`
	Description string          `json:"description,omitempty"`
	Nullable    bool            `json:"nullable"`
	Modes       []string        `json:"modes"`
	Review      []string        `json:"review_required"`
	Fields      []coverageField `json:"fields,omitempty"`
}

func extract(opts extractOptions) error {
	if opts.InputDir == "" || opts.Output == "" {
		return errors.New("--input-dir and --output are required")
	}
	if err := verify(opts.InputDir); err != nil {
		return fmt.Errorf("verify canonical inputs: %w", err)
	}
	report, err := extractCoverage(opts.InputDir)
	if err != nil {
		return err
	}
	encoded, err := marshalCanonical(report)
	if err != nil {
		return fmt.Errorf("encode coverage manifest: %w", err)
	}
	if opts.Check {
		existing, err := os.ReadFile(opts.Output)
		if err != nil {
			return fmt.Errorf("read generated manifest for check: %w", err)
		}
		if !bytes.Equal(existing, encoded) {
			return fmt.Errorf("generated manifest differs: rerun specgen extract --input-dir %s --output %s", opts.InputDir, opts.Output)
		}
		return nil
	}
	if err := writeFile(opts.Output, encoded); err != nil {
		return fmt.Errorf("write coverage manifest: %w", err)
	}
	return nil
}

func extractCoverage(inputDir string) (coverageManifest, error) {
	rawManifest, err := os.ReadFile(filepath.Join(inputDir, "manifest.json"))
	if err != nil {
		return coverageManifest{}, fmt.Errorf("read input manifest: %w", err)
	}
	var inputs manifest
	if err := json.Unmarshal(rawManifest, &inputs); err != nil {
		return coverageManifest{}, fmt.Errorf("decode input manifest: %w", err)
	}
	resources := map[string]*coverageResource{}
	for _, source := range inputs.Sources {
		document, err := readJSONDocument(filepath.Join(inputDir, source.File))
		if err != nil {
			return coverageManifest{}, fmt.Errorf("read %s OpenAPI: %w", source.Mode, err)
		}
		paths, ok := object(document["paths"])
		if !ok {
			return coverageManifest{}, fmt.Errorf("%s OpenAPI has no paths object", source.Mode)
		}
		for path, rawPathItem := range paths {
			pathItem, ok := object(rawPathItem)
			if !ok {
				continue
			}
			operations := resourceOperations(pathItem)
			if len(operations) == 0 {
				continue
			}
			resource := resources[path]
			if resource == nil {
				resource = &coverageResource{Path: path}
				resources[path] = resource
			}
			resource.Modes = appendUnique(resource.Modes, source.Mode)
			for _, operation := range operations {
				resource.Operations = appendUnique(resource.Operations, operation)
				if operation == "put" || operation == "patch" {
					mergeRequestShape(resource, source.Mode, pathItem[operation])
				}
				if operation == "get" {
					mergeResponseShape(resource, source.Mode, pathItem[operation])
				}
				if operation == "delete" {
					mergeDeleteParameters(resource, source.Mode, pathItem[operation])
				}
			}
		}
	}

	result := coverageManifest{FormatVersion: extractionFormatVersion, APIVersion: inputs.APIVersion, Resources: make([]coverageResource, 0, len(resources))}
	for _, resource := range resources {
		sort.Strings(resource.Modes)
		sort.Strings(resource.Operations)
		sortCoverageFields(resource.Fields)
		sort.Slice(resource.DeleteParameters, func(i, j int) bool { return resource.DeleteParameters[i].Name < resource.DeleteParameters[j].Name })
		for index := range resource.DeleteParameters {
			sort.Strings(resource.DeleteParameters[index].Modes)
		}
		resource.ReviewRequired = appendUnique(resource.ReviewRequired, "Terraform type, identity, aliases, and lifecycle policies", "cache key requires reviewed override")
		sort.Strings(resource.ReviewRequired)
		result.Resources = append(result.Resources, *resource)
	}
	sort.Slice(result.Resources, func(i, j int) bool { return result.Resources[i].Path < result.Resources[j].Path })
	return result, nil
}

func readJSONDocument(path string) (map[string]any, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.UseNumber()
	var document map[string]any
	if err := decoder.Decode(&document); err != nil {
		return nil, err
	}
	return document, nil
}

func resourceOperations(pathItem map[string]any) []string {
	mutable := false
	for _, operation := range []string{"put", "patch", "delete"} {
		if _, ok := object(pathItem[operation]); ok {
			mutable = true
			break
		}
	}
	if !mutable {
		return nil
	}
	operations := make([]string, 0, 4)
	for _, operation := range []string{"get", "put", "patch", "delete"} {
		if _, ok := object(pathItem[operation]); ok {
			operations = append(operations, operation)
		}
	}
	return operations
}

func mergeResponseShape(resource *coverageResource, mode string, rawOperation any) {
	operation, ok := object(rawOperation)
	if !ok {
		return
	}
	responses, ok := object(operation["responses"])
	if !ok {
		resource.ReviewRequired = appendUnique(resource.ReviewRequired, "GET response shape is missing")
		return
	}
	for _, status := range []string{"200", "201", "default"} {
		response, exists := object(responses[status])
		if !exists {
			continue
		}
		content, exists := object(response["content"])
		if !exists {
			resource.ReviewRequired = appendUnique(resource.ReviewRequired, "GET response schema is unavailable")
			return
		}
		jsonContent, exists := object(content["application/json"])
		if !exists {
			resource.ReviewRequired = appendUnique(resource.ReviewRequired, "GET application/json response schema is unavailable")
			return
		}
		schema, exists := object(jsonContent["schema"])
		if !exists {
			resource.ReviewRequired = appendUnique(resource.ReviewRequired, "GET response schema is unavailable")
			return
		}
		properties, exists := object(schema["properties"])
		if !exists || len(properties) != 1 {
			resource.ReviewRequired = appendUnique(resource.ReviewRequired, "GET response collection key is ambiguous")
			return
		}
		for key := range properties {
			resource.ResponseCollectionKey = mergeString(resource.ResponseCollectionKey, key, &resource.ReviewRequired, "response collection key differs between modes")
		}
		return
	}
	resource.ReviewRequired = appendUnique(resource.ReviewRequired, "GET success response is undocumented")
}

func mergeDeleteParameters(resource *coverageResource, mode string, rawOperation any) {
	operation, ok := object(rawOperation)
	if !ok {
		return
	}
	parameters, ok := operation["parameters"].([]any)
	if !ok {
		resource.ReviewRequired = appendUnique(resource.ReviewRequired, "delete parameter semantics are missing")
		return
	}
	for _, rawParameter := range parameters {
		parameter, ok := object(rawParameter)
		if !ok {
			resource.ReviewRequired = appendUnique(resource.ReviewRequired, "delete parameter reference is unresolved")
			continue
		}
		if stringValue(parameter["in"]) != "query" {
			continue
		}
		schema, ok := object(parameter["schema"])
		if !ok {
			resource.ReviewRequired = appendUnique(resource.ReviewRequired, "delete parameter schema is missing")
			continue
		}
		incoming := coverageQueryParameter{Name: stringValue(parameter["name"]), Required: boolValue(parameter["required"]), Kind: fieldKind(schema), Modes: []string{mode}}
		if items, ok := object(schema["items"]); ok {
			incoming.ItemKind = fieldKind(items)
		}
		if incoming.Name == "" {
			resource.ReviewRequired = appendUnique(resource.ReviewRequired, "delete parameter name is missing")
			continue
		}
		mergeDeleteParameter(resource, incoming)
	}
}

func mergeDeleteParameter(resource *coverageResource, incoming coverageQueryParameter) {
	for index := range resource.DeleteParameters {
		parameter := &resource.DeleteParameters[index]
		if parameter.Name != incoming.Name {
			continue
		}
		parameter.Modes = appendUnique(parameter.Modes, incoming.Modes...)
		if parameter.Kind != incoming.Kind || parameter.ItemKind != incoming.ItemKind || parameter.Required != incoming.Required {
			resource.ReviewRequired = appendUnique(resource.ReviewRequired, "delete parameter differs between modes")
		}
		return
	}
	resource.DeleteParameters = append(resource.DeleteParameters, incoming)
}

func mergeRequestShape(resource *coverageResource, mode string, rawOperation any) {
	operation, ok := object(rawOperation)
	if !ok {
		return
	}
	requestBody, ok := object(operation["requestBody"])
	if !ok {
		resource.ReviewRequired = appendUnique(resource.ReviewRequired, "request body shape is missing")
		return
	}
	content, ok := object(requestBody["content"])
	if !ok {
		resource.ReviewRequired = appendUnique(resource.ReviewRequired, "request body content is missing")
		return
	}
	jsonContent, ok := object(content["application/json"])
	if !ok {
		resource.ReviewRequired = appendUnique(resource.ReviewRequired, "application/json request body is missing")
		return
	}
	schema, ok := object(jsonContent["schema"])
	if !ok {
		resource.ReviewRequired = appendUnique(resource.ReviewRequired, "request schema is missing")
		return
	}
	properties, ok := object(schema["properties"])
	if !ok || len(properties) == 0 {
		resource.ReviewRequired = appendUnique(resource.ReviewRequired, "request schema has no inline properties")
		return
	}

	fieldProperties := properties
	if len(properties) == 1 {
		for wrapper, rawWrapper := range properties {
			wrapperSchema, wrapperOK := object(rawWrapper)
			inner, innerOK := object(wrapperSchema["properties"])
			if wrapperOK && innerOK && len(inner) == 1 {
				resource.RequestWrapperKey = mergeString(resource.RequestWrapperKey, wrapper, &resource.ReviewRequired, "request wrapper differs between modes or operations")
				for _, rawValue := range inner {
					if valueSchema, ok := object(rawValue); ok {
						if valueProperties, ok := object(valueSchema["properties"]); ok {
							fieldProperties = valueProperties
						} else {
							resource.ReviewRequired = appendUnique(resource.ReviewRequired, "request value schema has no inline properties")
						}
					}
				}
			}
		}
	} else {
		resource.ReviewRequired = appendUnique(resource.ReviewRequired, "request wrapper is ambiguous")
	}
	for name, rawField := range fieldProperties {
		fieldSchema, ok := object(rawField)
		if !ok {
			continue
		}
		mergeField(resource, coverageFieldFromSchema(name, fieldSchema, mode))
	}
}

func coverageFieldFromSchema(name string, value map[string]any, mode string) coverageField {
	field := coverageField{APIName: name, Kind: fieldKind(value), Description: stringValue(value["description"]), Nullable: boolValue(value["nullable"]), Modes: []string{mode}, Review: fieldReview(value)}
	nestedSchema := value
	if field.Kind == "array" {
		if items, ok := object(value["items"]); ok {
			field.ItemKind = fieldKind(items)
			nestedSchema = items
		}
	}
	if properties, ok := object(nestedSchema["properties"]); ok {
		for nestedName, rawNested := range properties {
			if nested, ok := object(rawNested); ok {
				field.Fields = append(field.Fields, coverageFieldFromSchema(nestedName, nested, mode))
			}
		}
		sortCoverageFields(field.Fields)
	}
	return field
}

func mergeField(resource *coverageResource, incoming coverageField) {
	for index := range resource.Fields {
		field := &resource.Fields[index]
		if field.APIName != incoming.APIName {
			continue
		}
		field.Modes = appendUnique(field.Modes, incoming.Modes...)
		field.Review = appendUnique(field.Review, incoming.Review...)
		if field.Kind != incoming.Kind {
			field.Review = appendUnique(field.Review, "type differs between modes or operations")
		}
		if field.ItemKind != incoming.ItemKind {
			field.Review = appendUnique(field.Review, "item type differs between modes or operations")
		}
		if field.Description == "" {
			field.Description = incoming.Description
		}
		field.Nullable = field.Nullable || incoming.Nullable
		mergeCoverageFields(&field.Fields, incoming.Fields)
		return
	}
	resource.Fields = append(resource.Fields, incoming)
}

func mergeCoverageFields(current *[]coverageField, incoming []coverageField) {
	for _, field := range incoming {
		found := false
		for index := range *current {
			if (*current)[index].APIName != field.APIName {
				continue
			}
			mergeCoverageField(&(*current)[index], field)
			found = true
			break
		}
		if !found {
			*current = append(*current, field)
		}
	}
}

func mergeCoverageField(current *coverageField, incoming coverageField) {
	current.Modes = appendUnique(current.Modes, incoming.Modes...)
	current.Review = appendUnique(current.Review, incoming.Review...)
	if current.Kind != incoming.Kind {
		current.Review = appendUnique(current.Review, "type differs between modes or operations")
	}
	if current.ItemKind != incoming.ItemKind {
		current.Review = appendUnique(current.Review, "item type differs between modes or operations")
	}
	if current.Description == "" {
		current.Description = incoming.Description
	}
	current.Nullable = current.Nullable || incoming.Nullable
	mergeCoverageFields(&current.Fields, incoming.Fields)
}

func sortCoverageFields(fields []coverageField) {
	sort.Slice(fields, func(i, j int) bool { return fields[i].APIName < fields[j].APIName })
	for index := range fields {
		sort.Strings(fields[index].Modes)
		sort.Strings(fields[index].Review)
		sortCoverageFields(fields[index].Fields)
	}
}

func fieldReview(schema map[string]any) []string {
	review := []string{"Terraform access and lifecycle policies"}
	if _, nested := object(schema["properties"]); nested {
		review = append(review, "nested collection strategy")
	}
	if stringValue(schema["type"]) == "array" {
		review = append(review, "collection identity and update strategy")
	}
	return review
}

func fieldKind(schema map[string]any) string {
	if kind := stringValue(schema["type"]); kind != "" {
		return kind
	}
	return "unknown"
}
func object(value any) (map[string]any, bool) {
	object, ok := value.(map[string]any)
	return object, ok
}
func stringValue(value any) string { result, _ := value.(string); return result }
func boolValue(value any) bool     { result, _ := value.(bool); return result }
func appendUnique(values []string, incoming ...string) []string {
	for _, value := range incoming {
		found := false
		for _, current := range values {
			if current == value {
				found = true
				break
			}
		}
		if !found {
			values = append(values, value)
		}
	}
	return values
}
func mergeString(current, incoming string, review *[]string, conflict string) string {
	if current == "" || current == incoming {
		return incoming
	}
	*review = appendUnique(*review, conflict)
	return current
}
