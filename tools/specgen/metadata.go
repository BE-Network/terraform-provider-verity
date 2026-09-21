package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go/format"
	"os"
	"sort"
	"strings"

	"terraform-provider-verity/internal/spec"
)

type metadataOptions struct {
	Registry   string
	Output     string
	BulkOutput string
	KeysOutput string
	Check      bool
}

func generateModeMetadata(opts metadataOptions) error {
	if opts.Registry == "" || opts.Output == "" {
		return fmt.Errorf("--registry and --output are required")
	}
	raw, err := os.ReadFile(opts.Registry)
	if err != nil {
		return fmt.Errorf("read registry: %w", err)
	}
	var artifact struct {
		Resources spec.Registry `json:"resources"`
	}
	if err := json.Unmarshal(raw, &artifact); err != nil {
		return fmt.Errorf("decode registry: %w", err)
	}
	if err := artifact.Resources.Validate(); err != nil {
		return fmt.Errorf("validate registry: %w", err)
	}

	type resourceEntry struct {
		TerraformType string
		Mode          string
		EndpointKey   string
		Fields        [][2]string
	}
	entries := make([]resourceEntry, 0, len(artifact.Resources))
	for _, resource := range artifact.Resources {
		fields := make([][2]string, 0)
		collectModeFields(resource.Fields, "", resource.Modes, &fields)
		sort.Slice(fields, func(i, j int) bool { return fields[i][0] < fields[j][0] })
		entries = append(entries, resourceEntry{
			TerraformType: resource.TerraformType,
			Mode:          modeConstant(resource.Modes),
			EndpointKey:   strings.TrimPrefix(resource.API.EndpointPath, "/"),
			Fields:        fields,
		})
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].TerraformType < entries[j].TerraformType })

	var buf bytes.Buffer
	buf.WriteString("package utils\n\n")

	buf.WriteString("var generatedResourceCompatibility = map[string]ResourceMode{\n")
	for _, entry := range entries {
		fmt.Fprintf(&buf, "\t%q: %s,\n", entry.TerraformType, entry.Mode)
	}
	buf.WriteString("}\n\n")

	buf.WriteString("var generatedModeFields = map[string]map[string]FieldMode{\n")
	byEndpoint := make(map[string][][2]string, len(entries))
	endpoints := make([]string, 0, len(entries))
	for _, entry := range entries {
		if _, seen := byEndpoint[entry.EndpointKey]; !seen {
			endpoints = append(endpoints, entry.EndpointKey)
		}
		byEndpoint[entry.EndpointKey] = mergeFieldModes(byEndpoint[entry.EndpointKey], entry.Fields)
	}
	sort.Strings(endpoints)
	for _, endpoint := range endpoints {
		fmt.Fprintf(&buf, "\t%q: {\n", endpoint)
		for _, field := range byEndpoint[endpoint] {
			fmt.Fprintf(&buf, "\t\t%q: %s,\n", field[0], field[1])
		}
		buf.WriteString("\t},\n")
	}
	buf.WriteString("}\n")

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return fmt.Errorf("format generated metadata: %w", err)
	}
	if err := emitOrCheck(opts, opts.Output, formatted); err != nil {
		return err
	}
	if opts.BulkOutput == "" {
		return nil
	}
	bulk, err := renderBulkMetadata(artifact.Resources)
	if err != nil {
		return err
	}
	if err := emitOrCheck(opts, opts.BulkOutput, bulk); err != nil {
		return err
	}
	if opts.KeysOutput == "" {
		return nil
	}
	keys, err := renderResourceKeys(artifact.Resources)
	if err != nil {
		return err
	}
	return emitOrCheck(opts, opts.KeysOutput, keys)
}

func renderResourceKeys(resources spec.Registry) ([]byte, error) {
	sorted := make(spec.Registry, len(resources))
	copy(sorted, resources)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].TerraformType < sorted[j].TerraformType })

	var buf bytes.Buffer
	buf.WriteString("package provider\n\n")
	buf.WriteString("type generatedResourceKey struct {\n\tEndpoint              string\n\tCacheKey              string\n\tResponseCollectionKey string\n}\n\n")
	buf.WriteString("var generatedResourceKeys = map[string]generatedResourceKey{\n")
	for _, resource := range sorted {
		fmt.Fprintf(&buf, "\t%q: {Endpoint: %q, CacheKey: %q, ResponseCollectionKey: %q},\n",
			resource.TerraformType,
			strings.TrimPrefix(resource.API.EndpointPath, "/"),
			resource.API.CacheKey,
			resource.API.ResponseCollectionKey)
	}
	buf.WriteString("}\n\n")

	buf.WriteString("var generatedResourceOrder = []string{\n")
	for _, resource := range sorted {
		fmt.Fprintf(&buf, "\t%q,\n", resource.TerraformType)
	}
	buf.WriteString("}\n")
	return format.Source(buf.Bytes())
}

func emitOrCheck(opts metadataOptions, path string, content []byte) error {
	if opts.Check {
		existing, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s for check: %w", path, err)
		}
		if !bytes.Equal(existing, content) {
			return fmt.Errorf("%s differs: rerun specgen metadata --registry %s", path, opts.Registry)
		}
		return nil
	}
	return writeFile(path, content)
}

func renderBulkMetadata(resources spec.Registry) ([]byte, error) {
	type bulkEntry struct {
		HeaderSplitKey string
		Endpoint       string
	}
	entries := make(map[string]bulkEntry, len(resources))
	keys := make([]string, 0, len(resources))
	for _, resource := range resources {

		if len(resource.API.FixedHeaders) > 1 {
			names := make([]string, 0, len(resource.API.FixedHeaders))
			for name := range resource.API.FixedHeaders {
				names = append(names, name)
			}
			sort.Strings(names)
			return nil, fmt.Errorf("%s declares %d fixed headers (%s); the bulk transport supports one split key",
				resource.TerraformType, len(names), strings.Join(names, ", "))
		}
		split := ""
		for name := range resource.API.FixedHeaders {
			split = name
		}
		existing, seen := entries[resource.API.BulkKey]
		if !seen {
			keys = append(keys, resource.API.BulkKey)
		} else if existing.HeaderSplitKey != split {
			return nil, fmt.Errorf("bulk key %q has conflicting header split keys %q and %q", resource.API.BulkKey, existing.HeaderSplitKey, split)
		}
		entries[resource.API.BulkKey] = bulkEntry{HeaderSplitKey: split, Endpoint: resource.API.EndpointPath}
	}
	sort.Strings(keys)

	var buf bytes.Buffer
	buf.WriteString("package bulkops\n\n")
	buf.WriteString("type generatedBulkResource struct {\n\tEndpoint       string\n\tHeaderSplitKey string\n}\n\n")
	buf.WriteString("var generatedBulkMetadata = map[string]generatedBulkResource{\n")
	for _, key := range keys {
		fmt.Fprintf(&buf, "\t%q: {Endpoint: %q, HeaderSplitKey: %q},\n", key, entries[key].Endpoint, entries[key].HeaderSplitKey)
	}
	buf.WriteString("}\n")
	return format.Source(buf.Bytes())
}

func collectModeFields(fields []spec.FieldSpec, prefix string, parentModes []spec.Mode, out *[][2]string) {
	for _, field := range fields {
		path := prefix + field.APIName
		*out = append(*out, [2]string{path, fieldModeConstant(field.Modes)})
		if len(field.Fields) != 0 {
			collectModeFields(field.Fields, path+".", field.Modes, out)
		}
	}
}

func mergeFieldModes(existing [][2]string, incoming [][2]string) [][2]string {
	if existing == nil {
		return incoming
	}
	index := make(map[string]string, len(existing))
	for _, field := range existing {
		index[field[0]] = field[1]
	}
	for _, field := range incoming {
		if current, seen := index[field[0]]; seen && current != field[1] {
			index[field[0]] = "FieldModeBoth"
			continue
		}
		index[field[0]] = field[1]
	}
	merged := make([][2]string, 0, len(index))
	for path, mode := range index {
		merged = append(merged, [2]string{path, mode})
	}
	sort.Slice(merged, func(i, j int) bool { return merged[i][0] < merged[j][0] })
	return merged
}

func modeConstant(modes []spec.Mode) string {
	datacenter, campus := hasMode(modes, spec.ModeDatacenter), hasMode(modes, spec.ModeCampus)
	switch {
	case datacenter && campus:
		return "ResourceModeBoth"
	case campus:
		return "ResourceModeCampus"
	default:
		return "ResourceModeDatacenter"
	}
}

func fieldModeConstant(modes []spec.Mode) string {
	datacenter, campus := hasMode(modes, spec.ModeDatacenter), hasMode(modes, spec.ModeCampus)
	switch {
	case datacenter && campus:
		return "FieldModeBoth"
	case campus:
		return "FieldModeCampus"
	default:
		return "FieldModeDatacenter"
	}
}

func hasMode(modes []spec.Mode, want spec.Mode) bool {
	for _, mode := range modes {
		if mode == want {
			return true
		}
	}
	return false
}
