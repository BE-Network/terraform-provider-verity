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
