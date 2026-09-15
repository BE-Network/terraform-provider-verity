package lifecycle

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"terraform-provider-verity/tests/unit/mock"
)

// Adding an indexed child without naming its index does not work today.
//
// The migration plan asks Phase 0 to characterize "multiple new index-zero
// items", and nothing else exercised it: every other nested-block test writes
// indexes explicitly as 1..N and adds at most one entry per step, which never
// produces two entries sharing an index.
//
// The API documents the contract on the attribute itself — "The index
// identifying the object. Zero if you want to add an object to the list." — and
// `index` is Optional and Computed in all 54 collections, so both spellings of
// that intent are things a configuration can express. Neither survives an apply:
//
//   - Written as `index = 0`, the value is known, so Terraform holds the provider
//     to it. The request carries zero, the server assigns a real index, and the
//     apply fails with "inconsistent result after apply".
//   - Left out, the value plans as unknown. SetInt64Fields skips an unknown, so
//     the entry reaches the wire with no index member at all rather than the zero
//     the API asks for, and nothing correlates it back.
//
// So an indexed child can only be added by writing an index the server has not
// assigned yet, which means guessing the server's numbering.
//
// These tests pin that as it stands rather than assert the behavior anyone would
// want. They are characterization, in the Phase 0 sense: the generic collection
// engine in Phase 4 is where this gets fixed, and when it does, these tests fail
// and are rewritten deliberately instead of the change passing unnoticed.
const (
	indexZeroType         = "verity_packet_broker"
	indexZeroAPIPath      = "/api/packetbroker"
	indexZeroWrapperKey   = "pb_egress_profile"
	indexZeroResourceName = "idx0_pb"
	indexZeroBlock        = "ipv4_permit"
)

func indexZeroBaseline(url string) string {
	return mock.ProviderConfig(url, "datacenter") + fmt.Sprintf(`
resource %q "test" {
  name   = %q
  enable = true

  ipv4_permit {
    index  = 1
    enable = true
  }
}
`, indexZeroType, indexZeroResourceName)
}

// indexZeroAdded appends two entries to the baseline, both asking for an index.
// newEntry supplies the body of each, so the two spellings differ only in whether
// the index is written.
func indexZeroAdded(url, newEntry string) string {
	return mock.ProviderConfig(url, "datacenter") + fmt.Sprintf(`
resource %q "test" {
  name   = %q
  enable = true

  ipv4_permit {
    index  = 1
    enable = true
  }

  ipv4_permit {
%s
  }

  ipv4_permit {
%s
  }
}
`, indexZeroType, indexZeroResourceName, newEntry, newEntry)
}

func newIndexZeroServer(t *testing.T) *mock.MockServer {
	t.Helper()
	ms := mock.NewMockServer("datacenter")
	t.Cleanup(ms.Close)
	ms.SetTestLogger(t)
	if err := ms.LoadResponsesFromDir(mock.ResponsesDir("datacenter")); err != nil {
		t.Fatalf("failed to load responses: %v", err)
	}
	return ms
}

// TestIndexZeroAddRejectedAsInconsistentResult records the documented spelling:
// `index = 0` is a known value, so Terraform requires the applied state to match
// it, and the server's assigned index never will.
func TestIndexZeroAddRejectedAsInconsistentResult(t *testing.T) {
	t.Parallel()

	ms := newIndexZeroServer(t)
	baseline := indexZeroBaseline(ms.URL())
	added := indexZeroAdded(ms.URL(), "    index  = 0\n    enable = true")

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), baseline) },
				Config:    baseline,
			},
			{
				PreConfig: func() { ms.Reset(); mock.WriteTFConfig(t, ms.URL(), added) },
				Config:    added,
				// Both new entries do reach the request, each carrying index zero;
				// what fails is reconciling the indexes the server assigns back into
				// a configuration that pinned them to zero.
				ExpectError: regexp.MustCompile(`(?s)inconsistent result after apply.*index`),
			},
		},
	})
}

// TestIndexZeroAddOmitsTheIndexEntirely records the other spelling. Leaving the
// attribute out plans an unknown, and the guarded setter drops an unknown, so the
// entry is sent with no index member rather than the zero that would ask the
// server for one.
func TestIndexZeroAddOmitsTheIndexEntirely(t *testing.T) {
	t.Parallel()

	ms := newIndexZeroServer(t)
	baseline := indexZeroBaseline(ms.URL())
	added := indexZeroAdded(ms.URL(), "    enable = true")

	var patched bool
	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), baseline) },
				Config:    baseline,
			},
			{
				PreConfig: func() { ms.Reset(); mock.WriteTFConfig(t, ms.URL(), added) },
				Config:    added,
				Check: func(*terraform.State) error {
					patched = true
					return nil
				},
				// The entries carry no index, so the server cannot place them and the
				// read returns the one entry it already had.
				ExpectError: regexp.MustCompile(`(?s)inconsistent result after apply.*block count changed from 3 to 1`),
			},
		},
	})
	if patched {
		t.Fatal("the apply was expected to fail before its state check ran")
	}

	// The request itself is the evidence: two new entries were sent, and neither
	// named an index.
	patches := ms.GetRequestsByMethodAndPath("PATCH", indexZeroAPIPath)
	if len(patches) == 0 {
		t.Fatalf("no PATCH captured for %s", indexZeroAPIPath)
	}
	object, err := extractResourceFromBody(patches[len(patches)-1].Body, indexZeroWrapperKey, indexZeroResourceName)
	if err != nil {
		t.Fatal(err)
	}
	entries, ok := object[indexZeroBlock].([]interface{})
	if !ok {
		t.Fatalf("%s absent from the PATCH body", indexZeroBlock)
	}
	withoutIndex := 0
	for _, raw := range entries {
		item, ok := raw.(map[string]interface{})
		if !ok {
			t.Fatalf("%s entry is not an object", indexZeroBlock)
		}
		if _, named := item["index"]; !named {
			withoutIndex++
		}
	}
	if withoutIndex != 2 {
		t.Fatalf("%d of %d %s entries were sent without an index, want 2; if this is now 0 the provider "+
			"has started sending the zero the API asks for, and this test should assert that instead",
			withoutIndex, len(entries), indexZeroBlock)
	}
}
