package lifecycle

import (
	"strings"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"terraform-provider-verity/tests/unit/mock"
)

func TestDestroyingPredefinedObjectRemovesItFromStateOnly(t *testing.T) {
	entry := coverageEntry(t, "verity_packet_queue")
	assertServedGenerically(t, entry.TerraformType)

	ms := mock.NewMockServer(entry.Mode)
	defer ms.Close()
	ms.SetTestLogger(t)
	if err := ms.LoadResponsesFromDir(mock.ResponsesDir(entry.Mode)); err != nil {
		t.Fatal(err)
	}
	ms.SetPredefinedObjects(entry.APIPath, "packet_queue_name", "packet_queue", "(Packet Queue)")

	rs := inspectLegacySchema(t, entry)
	object := func(label, name string) string {
		hcl := generateCoverageHCL(t, rs, entry.TerraformType, name, entry.Mode, entry.modeFieldsKey(), entry.Overrides)
		return strings.Replace(hcl, `resource "verity_packet_queue" "test"`, `resource "verity_packet_queue" "`+label+`"`, 1)
	}
	provider := mock.ProviderConfig(ms.URL(), entry.Mode)
	both := provider + object("predefined", "(Packet Queue)") + object("ordinary", "Queue A")

	fwresource.UnitTest(t, fwresource.TestCase{
		ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
		Steps: []fwresource.TestStep{
			{
				PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), both) },
				Config:    both,
			},
			{
				PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), provider) },
				Config:    provider,
			},
		},
	})

	deletes := ms.GetRequestsByMethodAndPath("DELETE", entry.APIPath)
	if len(deletes) != 2 {
		t.Fatalf("got %d DELETE requests, want the rejected one and a retry", len(deletes))
	}
	first, retry := deletes[0].QueryParams["packet_queue_name"], deletes[1].QueryParams["packet_queue_name"]
	if len(first) != 2 {
		t.Fatalf("first DELETE names = %v, want both objects", first)
	}
	if len(retry) != 1 || retry[0] != "Queue A" {
		t.Fatalf("retry DELETE names = %v, want only the ordinary object", retry)
	}
}
