package importer

import "sort"

// SchemaFields is the part of a resource schema the importer checks generated
// configuration against: the argument names at one level, and the blocks nested
// in it.
//
// The importer writes whatever the API returns. A Verity system newer than the
// API version this provider supports returns arguments the provider's schema
// does not have, and a single one of them makes Terraform reject the whole
// generated configuration. Checking against the schema keeps the configuration
// importable; the arguments left out are reported so the user knows a newer
// provider version would manage them.
type SchemaFields struct {
	Attributes map[string]bool
	Blocks     map[string]*SchemaFields
}

// WithSupportedFields makes the importer leave out any argument the given
// schemas, keyed by Terraform resource type, do not have. Without it every
// argument the API returns is written, as before.
func (i *Importer) WithSupportedFields(fields map[string]*SchemaFields) *Importer {
	i.supported = fields
	return i
}

// UnsupportedFields returns, per Terraform resource type, the argument paths the
// importer left out because the schema does not have them. A nested argument is
// written as "block.argument".
func (i *Importer) UnsupportedFields() map[string][]string {
	result := make(map[string][]string, len(i.unsupported))
	for resourceType, paths := range i.unsupported {
		sorted := make([]string, 0, len(paths))
		for path := range paths {
			sorted = append(sorted, path)
		}
		sort.Strings(sorted)
		result[resourceType] = sorted
	}
	return result
}

// PruneUnsupported removes from every object of a Terraform resource type, as
// the API returns them keyed by name, the arguments the resource's schema does
// not have, and records what it removed. ImportAll calls it before writing each
// file.
func (i *Importer) PruneUnsupported(resourceType string, objects map[string]map[string]interface{}) {
	fields, known := i.supported[resourceType]
	if !known {
		return
	}
	config := resourceConfigs[terraformTypeToResourceKey[resourceType]]
	// These keys are never written as arguments, or are written by the header,
	// so the schema has no say over them.
	skip := map[string]bool{"name": true}
	for _, key := range config.AdditionalTopLevelSkipKeys {
		skip[key] = true
	}
	for _, object := range objects {
		i.pruneObject(resourceType, fields, object, "", skip, config.FieldMappings)
	}
}

func (i *Importer) pruneObject(resourceType string, fields *SchemaFields, object map[string]interface{}, prefix string, skip map[string]bool, mappings map[string]string) {
	for key, value := range object {
		if skip[key] {
			continue
		}
		name := key
		if mapped, found := mappings[key]; found {
			name = mapped
		}
		if fields.Attributes[name] {
			continue
		}
		if nested := fields.Blocks[name]; nested != nil {
			switch entries := value.(type) {
			case map[string]interface{}:
				i.pruneObject(resourceType, nested, entries, prefix+name+".", nil, nil)
			case []interface{}:
				for _, entry := range entries {
					if entryObject, isObject := entry.(map[string]interface{}); isObject {
						i.pruneObject(resourceType, nested, entryObject, prefix+name+".", nil, nil)
					}
				}
			}
			continue
		}
		delete(object, key)
		if i.unsupported == nil {
			i.unsupported = map[string]map[string]bool{}
		}
		if i.unsupported[resourceType] == nil {
			i.unsupported[resourceType] = map[string]bool{}
		}
		i.unsupported[resourceType][prefix+name] = true
	}
}
