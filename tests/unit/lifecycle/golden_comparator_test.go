package lifecycle

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// The golden comparison is what the Phase 2 parity claim rests on, so its
// ability to detect a difference is worth asserting rather than assuming.
//
// TestGenericIPv4ListMatchesLegacyGoldenFixtures passes when the generic engine
// reproduces the recorded bytes. That is only evidence if the comparison would
// have failed had it not. These tests exercise the comparator directly against
// the committed fixtures: identical content reports no difference, and each way
// the engine could plausibly diverge reports one.
//
// The two mutations below are the failure modes that matter. A dropped field is
// the request going out incomplete; a changed value is the response being
// decoded wrongly on the way back. Both were verified by hand against the live
// engine during Phase 2; these keep that verification repeatable.
func loadGolden(t *testing.T, terraformType, operation string) map[string]interface{} {
	t.Helper()
	path := filepath.Join("testdata", "golden", terraformType, operation+".json")
	raw, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read %s: %v", path, err)
	}
	var decoded map[string]interface{}
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("decode %s: %v", path, err)
	}
	return decoded
}

// guardAgainstRegeneration keeps these tests from rewriting the fixtures they
// are checking. compareGoldenValue writes instead of comparing when
// UPDATE_GOLDEN is set, which would turn a deliberate mutation into a committed
// one.
func guardAgainstRegeneration(t *testing.T) {
	t.Helper()
	if os.Getenv("UPDATE_GOLDEN") != "" {
		t.Skip("UPDATE_GOLDEN is set; the comparator writes rather than compares")
	}
}

func TestGoldenComparatorAcceptsIdenticalContent(t *testing.T) {
	guardAgainstRegeneration(t)

	for _, operation := range []string{"put", "patch", "state"} {
		fixture := loadGolden(t, "verity_ipv4_list", operation)
		if err := compareGoldenValue(t, "verity_ipv4_list", operation, fixture); err != nil {
			t.Errorf("%s: the fixture did not match itself: %v", operation, err)
		}
	}
}

// A field missing from the request must be reported. This is the shape of a
// codec that stopped sending something the configuration set.
func TestGoldenComparatorDetectsADroppedField(t *testing.T) {
	guardAgainstRegeneration(t)

	fixture := loadGolden(t, "verity_ipv4_list", "put")
	wrapper, ok := fixture["ipv4_list_filter"].(map[string]interface{})
	if !ok {
		t.Fatal("the put fixture has no ipv4_list_filter wrapper")
	}
	object, ok := wrapper["cov_ipv4l"].(map[string]interface{})
	if !ok {
		t.Fatal("the put fixture has no cov_ipv4l object")
	}
	if _, present := object["enable"]; !present {
		t.Fatal("the put fixture carries no enable field, so dropping it proves nothing")
	}
	delete(object, "enable")

	if err := compareGoldenValue(t, "verity_ipv4_list", "put", fixture); err == nil {
		t.Fatal("a request missing a field compared equal to its fixture")
	}
}

// A value decoded differently must be reported. This is the shape of a reader
// that mapped a response member wrongly, which no request comparison would see.
func TestGoldenComparatorDetectsAChangedStateValue(t *testing.T) {
	guardAgainstRegeneration(t)

	fixture := loadGolden(t, "verity_ipv4_list", "state")
	if _, present := fixture["name"]; !present {
		t.Fatal("the state fixture carries no name, so changing it proves nothing")
	}
	fixture["name"] = "perturbed"

	if err := compareGoldenValue(t, "verity_ipv4_list", "state", fixture); err == nil {
		t.Fatal("state with a changed value compared equal to its fixture")
	}
}
