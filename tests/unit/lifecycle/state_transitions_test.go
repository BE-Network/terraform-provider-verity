package lifecycle

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	"terraform-provider-verity/internal/registry"
	"terraform-provider-verity/internal/spec"
)

type stateTransitionPolicy struct {
	FormatVersion int                      `json:"format_version"`
	Releases      []stateTransitionRelease `json:"releases"`
}

type stateTransitionRelease struct {
	APIVersion     string                  `json:"api_version"`
	BaselineSHA256 string                  `json:"baseline_sha256"`
	Strategy       string                  `json:"strategy"`
	UpgradeGuide   string                  `json:"upgrade_guide"`
	Changes        []stateTransitionChange `json:"changes"`
}

type stateTransitionChange struct {
	TerraformType       string `json:"terraform_type"`
	ExpectedStateSHA256 string `json:"expected_state_sha256"`
	Reason              string `json:"reason"`
}

func loadStateTransitionPolicy(path, repoRoot string) (stateTransitionPolicy, error) {
	var policy stateTransitionPolicy
	raw, err := os.ReadFile(path)
	if err != nil {
		return policy, err
	}
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&policy); err != nil {
		return policy, err
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return policy, fmt.Errorf("state transition policy must contain exactly one JSON document")
	}
	for _, release := range policy.Releases {
		guide := release.UpgradeGuide
		if filepath.ToSlash(filepath.Clean(guide)) != guide || !strings.HasPrefix(guide, "docs/guides/") || !strings.HasSuffix(guide, ".md") {
			return policy, fmt.Errorf("%s: upgrade_guide must be a Markdown path under docs/guides", release.APIVersion)
		}
		contents, err := os.ReadFile(filepath.Join(repoRoot, filepath.FromSlash(guide)))
		if err != nil {
			return policy, fmt.Errorf("%s: read upgrade guide: %w", release.APIVersion, err)
		}
		if len(bytes.TrimSpace(contents)) == 0 {
			return policy, fmt.Errorf("%s: upgrade guide is empty", release.APIVersion)
		}
	}
	return policy, nil
}

func stateFingerprint(res resourceSchemaGolden) string {
	projection := stateRepresentation{TerraformType: res.TerraformType, SchemaVersion: res.SchemaVersion, Fields: stateFields(res.Attributes, res.Blocks, "")}
	raw, _ := json.Marshal(projection)
	return fmt.Sprintf("%x", sha256.Sum256(raw))
}

func validateStateTransitionPolicy(policy stateTransitionPolicy, baseline stateCompatibilityBaseline) error {
	if policy.FormatVersion != 1 {
		return fmt.Errorf("unsupported state transition policy format %d", policy.FormatVersion)
	}
	known := make(map[string]bool)
	for _, res := range baseline.Resources {
		known[res.TerraformType] = true
	}
	versions := make(map[string]bool)
	for _, release := range policy.Releases {
		if _, err := spec.ParseAPIVersion(release.APIVersion); err != nil {
			return fmt.Errorf("state transition API version: %w", err)
		}
		if versions[release.APIVersion] {
			return fmt.Errorf("duplicate state transition API version %s", release.APIVersion)
		}
		versions[release.APIVersion] = true
		if release.BaselineSHA256 != stateCompatibilitySHA256 {
			return fmt.Errorf("%s: baseline_sha256 does not match the pinned historical baseline", release.APIVersion)
		}
		if release.Strategy != "fresh_import" {
			return fmt.Errorf("%s: unsupported state transition strategy %q; only fresh_import is implemented", release.APIVersion, release.Strategy)
		}
		if strings.TrimSpace(release.UpgradeGuide) == "" || len(release.Changes) == 0 {
			return fmt.Errorf("%s: upgrade guide and explicit resource changes are required", release.APIVersion)
		}
		seen := make(map[string]bool)
		for _, change := range release.Changes {
			if !known[change.TerraformType] || seen[change.TerraformType] {
				return fmt.Errorf("%s: unknown or duplicate historical resource %q", release.APIVersion, change.TerraformType)
			}
			seen[change.TerraformType] = true
			if strings.TrimSpace(change.Reason) == "" {
				return fmt.Errorf("%s: %s needs a reason", release.APIVersion, change.TerraformType)
			}
			if change.ExpectedStateSHA256 != "removed" && (len(change.ExpectedStateSHA256) != 64 || strings.Trim(change.ExpectedStateSHA256, "0123456789abcdef") != "") {
				return fmt.Errorf("%s: %s needs an expected state SHA256 or removed", release.APIVersion, change.TerraformType)
			}
		}
	}
	return nil
}

func stateDifferences(prior stateRepresentation, current resourceSchemaGolden, exists bool) []string {
	if !exists {
		return []string{"resource removed"}
	}
	var differences []string
	if current.SchemaVersion != prior.SchemaVersion {
		differences = append(differences, fmt.Sprintf("schema version changed from %d to %d", prior.SchemaVersion, current.SchemaVersion))
	}
	fields := stateFields(current.Attributes, current.Blocks, "")
	paths := make([]string, 0, len(prior.Fields))
	for path := range prior.Fields {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	for _, path := range paths {
		if got, exists := fields[path]; !exists {
			differences = append(differences, path+" removed")
		} else if got != prior.Fields[path] {
			differences = append(differences, fmt.Sprintf("%s changed from %s to %s", path, prior.Fields[path], got))
		}
	}
	return differences
}

func checkStateCompatibility(baseline stateCompatibilityBaseline, current schemaGolden, apiVersion string, policy stateTransitionPolicy) error {
	if err := validateStateTransitionPolicy(policy, baseline); err != nil {
		return err
	}
	byName := make(map[string]resourceSchemaGolden)
	for _, res := range current.Resources {
		byName[res.TerraformType] = res
	}
	approved := make(map[string]stateTransitionChange)
	for _, release := range policy.Releases {
		if release.APIVersion != apiVersion {
			continue
		}
		for _, change := range release.Changes {
			res, exists := byName[change.TerraformType]
			actual := "removed"
			if exists {
				actual = stateFingerprint(res)
			}
			if actual != change.ExpectedStateSHA256 {
				return fmt.Errorf("%s: declared fresh-import target is stale: expected %s, got %s; review the new changes", change.TerraformType, change.ExpectedStateSHA256, actual)
			}
			approved[change.TerraformType] = change
		}
	}
	var problems []string
	for _, prior := range baseline.Resources {
		res, exists := byName[prior.TerraformType]
		differences := stateDifferences(prior, res, exists)
		if len(differences) == 0 {
			continue
		}
		if _, ok := approved[prior.TerraformType]; ok {
			continue
		}
		target := "removed"
		if exists {
			target = stateFingerprint(res)
		}
		problems = append(problems, fmt.Sprintf("%s: %s (expected_state_sha256: %s)", prior.TerraformType, strings.Join(differences, "; "), target))
	}
	if len(problems) != 0 {
		return fmt.Errorf("undeclared state change for API %s:\n%s\nReview and declare a fresh_import transition in specs/state_transitions.json with an upgrade guide; keep the historical baseline. Snapshot updates alone do not approve breaking changes", apiVersion, strings.Join(problems, "\n"))
	}
	return nil
}

func TestStateCompatibilityReport(t *testing.T) {
	baseline, err := loadStateCompatibilityBaseline()
	if err != nil {
		t.Fatal(err)
	}
	version, err := registry.Version()
	if err != nil {
		t.Fatal(err)
	}
	current := buildSchemaGolden(t)
	byName := make(map[string]resourceSchemaGolden)
	for _, res := range current.Resources {
		byName[res.TerraformType] = res
	}
	proposal := stateTransitionRelease{APIVersion: version.String(), BaselineSHA256: stateCompatibilitySHA256, Strategy: "fresh_import", Changes: []stateTransitionChange{}}
	for _, prior := range baseline.Resources {
		res, exists := byName[prior.TerraformType]
		if len(stateDifferences(prior, res, exists)) == 0 {
			continue
		}
		target := "removed"
		if exists {
			target = stateFingerprint(res)
		}
		proposal.Changes = append(proposal.Changes, stateTransitionChange{TerraformType: prior.TerraformType, ExpectedStateSHA256: target})
	}
	raw, err := json.MarshalIndent(proposal, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("Review this proposal and fill in upgrade_guide and every reason; this report grants no exceptions:\n%s", raw)
}
