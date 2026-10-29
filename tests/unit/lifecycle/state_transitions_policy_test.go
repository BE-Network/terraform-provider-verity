package lifecycle

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStateCompatibilityFreshImportTransitions(t *testing.T) {
	for _, test := range []struct {
		name   string
		mutate func(*schemaGolden, *stateTransitionPolicy)
		want   string
	}{
		{name: "removed endpoint and changed fields with new endpoint"},
		{name: "wrong API version", mutate: func(_ *schemaGolden, p *stateTransitionPolicy) { p.Releases[0].APIVersion = "6.8" }, want: "undeclared state change"},
		{name: "unlisted resource", mutate: func(_ *schemaGolden, p *stateTransitionPolicy) { p.Releases[0].Changes = p.Releases[0].Changes[:1] }, want: "verity_retired: resource removed"},
		{name: "unreviewed type change", mutate: func(c *schemaGolden, _ *stateTransitionPolicy) { c.Resources[0].Attributes[0].Type = "bool" }, want: "target is stale"},
		{name: "unreviewed added field", mutate: func(c *schemaGolden, _ *stateTransitionPolicy) {
			c.Resources[0].Attributes = append(c.Resources[0].Attributes, attributeGolden{Name: "extra", Type: "string"})
		}, want: "target is stale"},
		{name: "unreviewed version", mutate: func(c *schemaGolden, _ *stateTransitionPolicy) { c.Resources[0].SchemaVersion = 1 }, want: "target is stale"},
		{name: "resource restored", mutate: func(c *schemaGolden, _ *stateTransitionPolicy) {
			c.Resources = append(c.Resources, resourceSchemaGolden{TerraformType: "verity_retired"})
		}, want: "target is stale"},
		{name: "wrong baseline", mutate: func(_ *schemaGolden, p *stateTransitionPolicy) {
			p.Releases[0].BaselineSHA256 = strings.Repeat("0", 64)
		}, want: "pinned historical baseline"},
		{name: "missing reason", mutate: func(_ *schemaGolden, p *stateTransitionPolicy) { p.Releases[0].Changes[0].Reason = " " }, want: "needs a reason"},
		{name: "missing guide", mutate: func(_ *schemaGolden, p *stateTransitionPolicy) { p.Releases[0].UpgradeGuide = "" }, want: "upgrade guide"},
		{name: "unknown resource", mutate: func(_ *schemaGolden, p *stateTransitionPolicy) {
			p.Releases[0].Changes[0].TerraformType = "verity_typo"
		}, want: "unknown or duplicate"},
		{name: "duplicate resource", mutate: func(_ *schemaGolden, p *stateTransitionPolicy) {
			p.Releases[0].Changes = append(p.Releases[0].Changes, p.Releases[0].Changes[0])
		}, want: "unknown or duplicate"},
		{name: "duplicate release", mutate: func(_ *schemaGolden, p *stateTransitionPolicy) { p.Releases = append(p.Releases, p.Releases[0]) }, want: "duplicate state transition"},
		{name: "invalid fingerprint", mutate: func(_ *schemaGolden, p *stateTransitionPolicy) { p.Releases[0].Changes[0].ExpectedStateSHA256 = "*" }, want: "expected state SHA256"},
		{name: "unsupported strategy", mutate: func(_ *schemaGolden, p *stateTransitionPolicy) { p.Releases[0].Strategy = "ignore" }, want: "unsupported state transition strategy"},
		{name: "invalid version", mutate: func(_ *schemaGolden, p *stateTransitionPolicy) { p.Releases[0].APIVersion = "06.7" }, want: "API version"},
	} {
		t.Run(test.name, func(t *testing.T) {
			baseline := stateCompatibilityBaseline{Resources: []stateRepresentation{
				{TerraformType: "verity_changed", Fields: map[string]string{"name": "string", "old_field": "string"}},
				{TerraformType: "verity_retired", Fields: map[string]string{"name": "string"}},
			}}
			current := schemaGolden{Resources: []resourceSchemaGolden{
				{TerraformType: "verity_changed", Attributes: []attributeGolden{{Name: "name", Type: "string"}, {Name: "new_field", Type: "bool"}}},
				{TerraformType: "verity_new", Attributes: []attributeGolden{{Name: "name", Type: "string"}}},
			}}
			policy := stateTransitionPolicy{FormatVersion: 1, Releases: []stateTransitionRelease{{
				APIVersion: "6.7", BaselineSHA256: stateCompatibilitySHA256, Strategy: "fresh_import", UpgradeGuide: "docs/guides/upgrading-to-6-7.md",
				Changes: []stateTransitionChange{
					{TerraformType: "verity_changed", ExpectedStateSHA256: stateFingerprint(current.Resources[0]), Reason: "API removes old_field and adds new_field; reimport current objects"},
					{TerraformType: "verity_retired", ExpectedStateSHA256: "removed", Reason: "API endpoint retired; resource is no longer managed"},
				},
			}}}
			if test.mutate != nil {
				test.mutate(&current, &policy)
			}
			err := checkStateCompatibility(baseline, current, "6.7", policy)
			if test.want == "" {
				if err != nil {
					t.Fatal(err)
				}
			} else if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("error = %v, want %q", err, test.want)
			}
		})
	}
}

func TestStateCompatibilityPolicyLoading(t *testing.T) {
	for _, test := range []struct {
		name       string
		guide      string
		contents   string
		writeGuide bool
		extraJSON  string
		want       string
	}{
		{name: "documented release", guide: "docs/guides/upgrade.md", contents: "Fresh import instructions", writeGuide: true},
		{name: "missing guide", guide: "docs/guides/upgrade.md", want: "read upgrade guide"},
		{name: "empty guide", guide: "docs/guides/upgrade.md", contents: "\n", writeGuide: true, want: "guide is empty"},
		{name: "path traversal", guide: "docs/guides/../../../outside.md", want: "Markdown path"},
		{name: "absolute path", guide: "/tmp/upgrade.md", want: "Markdown path"},
		{name: "unknown key", extraJSON: `{"format_version":1,"releasez":[]}`, want: "unknown field"},
		{name: "trailing JSON", extraJSON: `{"format_version":1,"releases":[]} {}`, want: "exactly one JSON"},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			policy := stateTransitionPolicy{FormatVersion: 1, Releases: []stateTransitionRelease{{APIVersion: "6.7", UpgradeGuide: test.guide}}}
			raw, err := json.Marshal(policy)
			if err != nil {
				t.Fatal(err)
			}
			if test.extraJSON != "" {
				raw = []byte(test.extraJSON)
			}
			path := filepath.Join(root, "policy.json")
			if err := os.WriteFile(path, raw, 0o600); err != nil {
				t.Fatal(err)
			}
			if test.writeGuide {
				guide := filepath.Join(root, filepath.FromSlash(test.guide))
				if err := os.MkdirAll(filepath.Dir(guide), 0o700); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(guide, []byte(test.contents), 0o600); err != nil {
					t.Fatal(err)
				}
			}
			_, err = loadStateTransitionPolicy(path, root)
			if test.want == "" {
				if err != nil {
					t.Fatal(err)
				}
			} else if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("error = %v, want %q", err, test.want)
			}
		})
	}
}
