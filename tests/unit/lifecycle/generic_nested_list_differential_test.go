package lifecycle

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func TestGenericMatchesLegacyOnListInsideSingleton(t *testing.T) {
	entry := coverageEntry(t, "verity_fabric")
	rs := inspectSchema(entry.Factory)
	base := generateCoverageHCL(rs, entry.TerraformType, "difffabric", entry.Mode, entry.modeFieldsKey(), entry.Overrides)
	block := regexp.MustCompile(`(?s)\n  object_properties \{\n.*?\n  \}\n`)
	if !block.MatchString(base) {
		t.Fatal("the harness configuration no longer writes object_properties the way this test expects")
	}
	fabric := func(graphs ...int) string {
		return fabricWithBlock(t, base, block, true, graphs...)
	}
	noBlock := fabricWithBlock(t, base, block, false)

	scenarios := []struct {
		name           string
		create, update string
		outcome        lifecycleOutcome
	}{
		{name: "a graph is added", create: fabric(1), update: fabric(1, 2)},
		{name: "a graph is removed", create: fabric(1, 2), update: fabric(1)},
		{name: "the block is written with no graphs", create: fabric(), update: fabric()},
		{name: "the block is added with a graph", create: noBlock, update: fabric(1)},

		{name: "the block is removed with its graph", create: fabric(1), update: noBlock,
			outcome: lifecycleOutcome{applyError: regexp.MustCompile(`(?s)inconsistent result after apply.*object_properties: block count changed from 0 to 1`)}},
		{name: "graphs are untouched while another field changes", create: fabric(1),
			update: strings.Replace(fabric(1), "enable = true", "enable = false", 1)},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			legacy := captureLifecycle(t, entry.TerraformType, false, scenario.create, scenario.update, scenario.outcome)
			generic := captureLifecycle(t, entry.TerraformType, true, scenario.create, scenario.update, scenario.outcome)
			for _, operation := range []string{"PUT", "PATCH"} {
				want, got := canonical(t, legacy[operation]), canonical(t, generic[operation])
				if want != got {
					t.Errorf("%s differs between implementations\n  legacy:  %s\n  generic: %s", operation, want, got)
				}
			}
		})
	}
}

func fabricWithBlock(t *testing.T, base string, block *regexp.Regexp, present bool, graphs ...int) string {
	t.Helper()
	if !present {
		return block.ReplaceAllString(base, "\n")
	}
	var body strings.Builder
	body.WriteString("\n  object_properties {\n")
	for _, index := range graphs {
		fmt.Fprintf(&body, "    system_graphs {\n      index = %d\n    }\n", index)
	}
	body.WriteString("  }\n")
	return block.ReplaceAllLiteralString(base, body.String())
}
