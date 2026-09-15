// Package registry holds the reviewed resource registry the provider binary
// carries.
//
// The plan puts the validated registry in the binary so the generic engine reads
// one embedded description rather than loading a file at runtime. go:embed
// cannot reach outside its own package directory, so specgen writes the same
// bytes here and to specs/, and CI drift-checks both.
package registry

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"sync"

	"terraform-provider-verity/internal/spec"
)

//go:embed registry.json
var embedded []byte

type artifact struct {
	APIVersion string        `json:"api_version"`
	Resources  spec.Registry `json:"resources"`
}

var (
	once      sync.Once
	loaded    spec.Registry
	loadError error
)

// Load returns the embedded registry, validated. It is parsed once: the bytes
// never change for the life of the process, and every resource construction
// would otherwise re-parse them.
func Load() (spec.Registry, error) {
	once.Do(func() {
		var decoded artifact
		if err := json.Unmarshal(embedded, &decoded); err != nil {
			loadError = fmt.Errorf("decode embedded registry: %w", err)
			return
		}
		if err := decoded.Resources.Validate(); err != nil {
			loadError = fmt.Errorf("validate embedded registry: %w", err)
			return
		}
		loaded = decoded.Resources
	})
	return loaded, loadError
}

// Lookup returns one resource's reviewed spec by Terraform type.
func Lookup(terraformType string) (spec.ResourceSpec, error) {
	all, err := Load()
	if err != nil {
		return spec.ResourceSpec{}, err
	}
	for _, candidate := range all {
		if candidate.TerraformType == terraformType {
			return candidate, nil
		}
	}
	return spec.ResourceSpec{}, fmt.Errorf("no registry entry for %s", terraformType)
}
