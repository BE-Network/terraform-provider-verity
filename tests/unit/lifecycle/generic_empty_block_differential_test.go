package lifecycle

import (
	"regexp"
	"strings"
	"testing"
)

// An object the API declares with no properties is exposed by
// verity_device_settings as an empty block, and only its presence can change.
// The golden fixtures never write it — the harness skips blocks with no members —
// so these cases compare both implementations when it is written at create,
// added later, and removed.
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
		// The handwritten resource counts removing the block as a change but has
		// nothing to send, so the server keeps the object and the read restores
		// the block. The engine reproduces that failure for parity.
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
