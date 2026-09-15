package lifecycle

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"terraform-provider-verity/tests/unit/mock"
)

// TestPacketBrokerIPv6PermitEnableReachesPatch covers a reconciliation bug.
//
// Packet Broker has four structurally identical filter collections that share one
// model type. Three of them compared and sent the entry's enable flag on update;
// ipv6_permit did not, leaving the comment in place with the call removed. The
// effect was that Terraform planned a change to the flag, the PATCH omitted it,
// and the next read restored the old value, so the plan never converged.
//
// The other three collections define the intended contract, so this asserts the
// fixed behavior directly rather than through the field-coverage harness: an
// existing entry changing only enable must carry it in the PATCH body.
func TestPacketBrokerIPv6PermitEnableReachesPatch(t *testing.T) {
	t.Parallel()

	var entry ResourceCoverageEntry
	for _, candidate := range allResourceTests {
		if candidate.TerraformType == "verity_packet_broker" {
			entry = candidate
			break
		}
	}
	if entry.TerraformType == "" {
		t.Fatal("verity_packet_broker has no coverage entry")
	}

	ms := mock.NewMockServer(entry.Mode)
	defer ms.Close()
	ms.SetTestLogger(t)
	if err := ms.LoadResponsesFromDir(mock.ResponsesDir(entry.Mode)); err != nil {
		t.Fatalf("failed to load responses: %v", err)
	}

	schema := inspectSchema(entry.Factory)
	modeKey := entry.modeFieldsKey()
	name := "pb_ipv6_enable"

	createConfig := mock.ProviderConfig(ms.URL(), entry.Mode) +
		generateCoverageHCL(schema, entry.TerraformType, name, entry.Mode, modeKey,
			mergeOverrides(entry.Overrides, map[string]string{"ipv6_permit.enable": "true"}))
	updateConfig := mock.ProviderConfig(ms.URL(), entry.Mode) +
		generateCoverageHCL(schema, entry.TerraformType, name, entry.Mode, modeKey,
			mergeOverrides(entry.Overrides, map[string]string{"ipv6_permit.enable": "false"}))

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), createConfig) },
				Config:    createConfig,
			},
			{
				PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), updateConfig) },
				Config:    updateConfig,
				Check: func(*terraform.State) error {
					patches := ms.GetRequestsByMethodAndPath("PATCH", entry.APIPath)
					if len(patches) == 0 {
						return fmt.Errorf("no PATCH request captured for %s", entry.APIPath)
					}
					body := patches[len(patches)-1].Body
					wrapper, ok := body[entry.WrapperKey].(map[string]interface{})
					if !ok {
						return fmt.Errorf("wrapper %q absent from PATCH body", entry.WrapperKey)
					}
					object, ok := wrapper[name].(map[string]interface{})
					if !ok {
						return fmt.Errorf("resource %q absent from PATCH body", name)
					}
					entries, ok := object["ipv6_permit"].([]interface{})
					if !ok || len(entries) == 0 {
						return fmt.Errorf("ipv6_permit absent from PATCH body; the enable change was dropped")
					}
					first, ok := entries[0].(map[string]interface{})
					if !ok {
						return fmt.Errorf("ipv6_permit[0] is not an object")
					}
					value, present := first["enable"]
					if !present {
						return fmt.Errorf("ipv6_permit[0] carries no enable field; Terraform planned the change but the request omitted it")
					}
					if value != false {
						return fmt.Errorf("ipv6_permit[0].enable = %v, want false", value)
					}
					return nil
				},
			},
		},
	})
}
