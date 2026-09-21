package lifecycle

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

// verity_fabric's object_properties.system_graphs is the provider's only list
// nested inside a singleton, and the plan asks for dedicated tests of that path
// before the resource migrates. Its handwritten update reconciles the nested list
// whenever either side holds the block — unlike a singleton's scalar members — so
// adding or removing the block creates or deletes its entries. These cases run
// both implementations over each transition and compare what they send.
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
		// Removing the block deletes its graphs, but the object itself stays on
		// the server, so the read restores the block and Terraform rejects the
		// apply. The handwritten resource fails this way, as it does for an empty
		// block, and the engine reproduces it.
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

// fabricWithBlock replaces the harness's object_properties block with one
// holding the given system_graphs indexes, or removes it.
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
