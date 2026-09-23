package importer

import "sort"

type SchemaFields struct {
	Attributes map[string]bool
	Blocks     map[string]*SchemaFields
}

func (i *Importer) WithSupportedFields(fields map[string]*SchemaFields) *Importer {
	i.supported = fields
	return i
}

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

func (i *Importer) PruneUnsupported(resourceType string, objects map[string]map[string]interface{}) {
	fields, known := i.supported[resourceType]
	if !known {
		return
	}
	config, err := i.resourceConfig(resourceType)
	if err != nil {
		return
	}
	for _, object := range objects {
		i.pruneObject(resourceType, fields, object, "", config.SkipTopLevelKeys, config.FieldMappings)
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
