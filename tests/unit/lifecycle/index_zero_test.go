package lifecycle

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"terraform-provider-verity/tests/unit/mock"
)

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

				ExpectError: regexp.MustCompile(`(?s)inconsistent result after apply.*index`),
			},
		},
	})
}

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

				ExpectError: regexp.MustCompile(`(?s)inconsistent result after apply.*block count changed from 3 to 1`),
			},
		},
	})
	if patched {
		t.Fatal("the apply was expected to fail before its state check ran")
	}

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
