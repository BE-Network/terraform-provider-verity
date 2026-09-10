package genericresource

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

type PriorStateDescriptor struct {
	Version     int64
	PriorSchema schema.Schema
	Upgrade     resource.StateUpgrader
}

type StateUpgraderRegistry struct {
	upgraders map[int64]resource.StateUpgrader
}

func NewStateUpgraderRegistry(currentVersion int64, prior []PriorStateDescriptor) (StateUpgraderRegistry, error) {
	if currentVersion < 0 {
		return StateUpgraderRegistry{}, fmt.Errorf("current schema version cannot be negative")
	}
	upgraders := make(map[int64]resource.StateUpgrader, len(prior))
	for _, descriptor := range prior {
		if descriptor.Version < 0 || descriptor.Version >= currentVersion {
			return StateUpgraderRegistry{}, fmt.Errorf("prior schema version %d must be non-negative and less than current version %d", descriptor.Version, currentVersion)
		}
		if descriptor.Upgrade.StateUpgrader == nil {
			return StateUpgraderRegistry{}, fmt.Errorf("prior schema version %d has no upgrader", descriptor.Version)
		}
		if _, exists := upgraders[descriptor.Version]; exists {
			return StateUpgraderRegistry{}, fmt.Errorf("duplicate prior schema version %d", descriptor.Version)
		}
		descriptor.Upgrade.PriorSchema = &descriptor.PriorSchema
		upgraders[descriptor.Version] = descriptor.Upgrade
	}
	return StateUpgraderRegistry{upgraders: upgraders}, nil
}

func (r StateUpgraderRegistry) UpgradeState(context.Context) map[int64]resource.StateUpgrader {
	result := make(map[int64]resource.StateUpgrader, len(r.upgraders))
	for version, upgrader := range r.upgraders {
		result[version] = upgrader
	}
	return result
}
