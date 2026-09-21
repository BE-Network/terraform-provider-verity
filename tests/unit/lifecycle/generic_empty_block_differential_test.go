package lifecycle

import (
	"regexp"
	"strings"
	"testing"
)

func TestGenericMatchesLegacyOnEmptyBlock(t *testing.T) {
	entry := coverageEntry(t, "verity_device_settings")
	rs := inspectSchema(entry.Factory)
	base := generateCoverageHCL(rs, entry.TerraformType, "diffds", entry.Mode, entry.modeFieldsKey(), entry.Overrides)
	withBlock := func(present bool) string {
		if !present {
			return base
		}
		closing := strings.LastIndex(base, "}")
		return base[:closing] + "  object_properties {\n  }\n" + base[closing:]
	}

	scenarios := []struct {
		name           string
		create, update bool
		outcome        lifecycleOutcome
	}{
		{name: "written at create", create: true, update: true},
		{name: "added on update", create: false, update: true},

		{name: "removed on update", create: true, update: false,
			outcome: lifecycleOutcome{applyError: regexp.MustCompile(`(?s)inconsistent result after apply.*object_properties: block count changed from 0 to 1`)}},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			legacy := captureLifecycle(t, entry.TerraformType, false, withBlock(scenario.create), withBlock(scenario.update), scenario.outcome)
			generic := captureLifecycle(t, entry.TerraformType, true, withBlock(scenario.create), withBlock(scenario.update), scenario.outcome)
			for _, operation := range []string{"PUT", "PATCH"} {
				want, got := canonical(t, legacy[operation]), canonical(t, generic[operation])
				if want != got {
					t.Errorf("%s differs between implementations\n  legacy:  %s\n  generic: %s", operation, want, got)
				}
			}
		})
	}
}
