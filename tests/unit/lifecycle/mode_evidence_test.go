package lifecycle

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"
)

type modeEvidenceField struct {
	Name   string              `json:"api_name"`
	Modes  []string            `json:"modes"`
	Fields []modeEvidenceField `json:"fields"`
}

type modeEvidenceResource struct {
	Path   string              `json:"path"`
	Modes  []string            `json:"modes"`
	Fields []modeEvidenceField `json:"fields"`
}

type modeEvidence map[string]map[string]map[string]bool

var modeEvidenceOnce sync.Once
var extractedModeEvidence modeEvidence
var extractedModeError error

func parseModeEvidence(raw []byte) (modeEvidence, error) {
	var manifest struct {
		FormatVersion int                    `json:"format_version"`
		Resources     []modeEvidenceResource `json:"resources"`
	}
	if err := json.Unmarshal(raw, &manifest); err != nil {
		return nil, err
	}
	if manifest.FormatVersion != 1 || len(manifest.Resources) == 0 {
		return nil, fmt.Errorf("invalid OpenAPI extraction manifest")
	}
	index := modeEvidence{}
	for _, resource := range manifest.Resources {
		endpoint := strings.Trim(resource.Path, "/")
		if endpoint == "" || index[endpoint] != nil {
			return nil, fmt.Errorf("missing or duplicate manifest endpoint %q", endpoint)
		}
		fields := map[string]map[string]bool{}
		var walk func([]modeEvidenceField, string, []string) error
		walk = func(children []modeEvidenceField, prefix string, parentModes []string) error {
			for _, field := range children {
				if field.Name == "" {
					continue
				}
				path := prefix + field.Name
				if fields[path] != nil {
					return fmt.Errorf("duplicate manifest field %s.%s", endpoint, path)
				}
				modes := map[string]bool{}
				for _, mode := range field.Modes {
					if mode != "campus" && mode != "datacenter" {
						return fmt.Errorf("unknown manifest mode %q", mode)
					}
					for _, parent := range parentModes {
						if mode == parent {
							modes[mode] = true
						}
					}
				}
				if len(modes) == 0 {
					return fmt.Errorf("manifest field %s.%s has no applicable mode", endpoint, path)
				}
				fields[path] = modes
				applicable := make([]string, 0, len(modes))
				for mode := range modes {
					applicable = append(applicable, mode)
				}
				if err := walk(field.Fields, path+".", applicable); err != nil {
					return err
				}
			}
			return nil
		}
		if err := walk(resource.Fields, "", resource.Modes); err != nil {
			return nil, err
		}
		index[endpoint] = fields
	}
	return index, nil
}

func (evidence modeEvidence) applies(endpoint, path, mode string) (bool, error) {
	if mode != "campus" && mode != "datacenter" {
		return false, fmt.Errorf("unknown mode %q", mode)
	}
	fields, found := evidence[strings.Trim(endpoint, "/")]
	if !found {
		return false, fmt.Errorf("unknown manifest endpoint %q", endpoint)
	}
	modes, found := fields[path]
	if !found {
		return false, fmt.Errorf("unknown manifest field %s.%s", endpoint, path)
	}
	return modes[mode], nil
}

func fieldAppliesToMode(t *testing.T, endpoint, path, mode string) bool {
	t.Helper()
	modeEvidenceOnce.Do(func() {
		raw, err := os.ReadFile("../../../specs/generated_manifest.json")
		if err != nil {
			extractedModeError = err
			return
		}
		extractedModeEvidence, extractedModeError = parseModeEvidence(raw)
	})
	if extractedModeError != nil {
		t.Fatalf("load extracted mode evidence: %v", extractedModeError)
	}
	applies, err := extractedModeEvidence.applies(endpoint, path, mode)
	if err != nil {
		t.Fatal(err)
	}
	return applies
}

func TestModeEvidenceRefusesUnknownLookups(t *testing.T) {
	evidence, err := parseModeEvidence([]byte(`{"format_version":1,"resources":[{"path":"/widgets","modes":["campus","datacenter"],"fields":[{"api_name":"","modes":["campus"]},{"api_name":"settings","modes":["campus"],"fields":[{"api_name":"enable","modes":["campus","datacenter"]}]}]}]}`))
	if err != nil {
		t.Fatal(err)
	}
	for _, lookup := range []struct {
		endpoint, path, mode string
		want                 bool
		invalid              bool
	}{
		{"widgets", "settings.enable", "campus", true, false},
		{"widgets", "settings.enable", "datacenter", false, false},
		{"widget", "settings.enable", "campus", false, true},
		{"widgets", "settings.enabel", "campus", false, true},
		{"widgets", "settings.enable", "datacentre", false, true},
	} {
		got, err := evidence.applies(lookup.endpoint, lookup.path, lookup.mode)
		if lookup.invalid {
			if err == nil {
				t.Errorf("invalid lookup %+v accepted", lookup)
			}
			continue
		}
		if err != nil || got != lookup.want {
			t.Errorf("lookup %+v = %t, %v", lookup, got, err)
		}
	}
}

func TestModeEvidenceLoadsExtractionManifest(t *testing.T) {
	for _, test := range []struct {
		path, mode string
		want       bool
	}{
		{"anycast_ipv4_mask", "datacenter", true},
		{"anycast_ipv4_mask", "campus", false},
		{"packet_priority", "campus", true},
		{"packet_priority", "datacenter", false},
		{"object_properties.warn_on_no_external_source", "campus", true},
		{"object_properties.warn_on_no_external_source", "datacenter", false},
	} {
		if got := fieldAppliesToMode(t, "services", test.path, test.mode); got != test.want {
			t.Errorf("services.%s in %s = %t, want %t", test.path, test.mode, got, test.want)
		}
	}
}
