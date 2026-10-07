package lifecycle

import (
	"fmt"
	"reflect"
	"testing"

	fwresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"terraform-provider-verity/tests/unit/mock"
)

func TestMonitoringACLUpdatesServicesByIndex(t *testing.T) {
	for _, mode := range []string{"datacenter", "campus"} {
		t.Run(mode, func(t *testing.T) {
			ms := newImporterServer(t, mode)
			config := func(services string) string {
				return mock.ProviderConfig(ms.URL(), mode) + `
resource "verity_monitoring_acl" "test" {
  name = "monitoring-list-update"
  enable = true
` + services + "\n}\n"
			}
			service := func(index int, name string, enable bool) string {
				return fmt.Sprintf(`
  services {
    index = %d
    service = %q
    service_ref_type_ = "service"
    enable = %t
  }
`, index, name, enable)
			}
			check := func(want ...map[string]interface{}) fwresource.TestCheckFunc {
				return func(_ *terraform.State) error {
					patches := ms.GetRequestsByMethodAndPath("PATCH", "/api/monitoringacls")
					if len(patches) == 0 {
						return fmt.Errorf("no monitoring ACL PATCH recorded")
					}
					body := patches[len(patches)-1].Body
					objects := body["monitoring_acl"].(map[string]interface{})
					object := objects["monitoring-list-update"].(map[string]interface{})
					rows, ok := object["services"].([]interface{})
					if !ok || len(rows) != len(want) {
						return fmt.Errorf("services = %#v, want %v", object["services"], want)
					}
					for index, raw := range rows {
						if !reflect.DeepEqual(raw, want[index]) {
							return fmt.Errorf("service row = %#v, want %#v", raw, want[index])
						}
					}
					return nil
				}
			}
			fwresource.UnitTest(t, fwresource.TestCase{
				ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
				Steps: []fwresource.TestStep{
					{Config: config(service(1, "a", true) + service(2, "b", true))},
					{Config: config(service(1, "a", false) + service(2, "b", true)), Check: check(map[string]interface{}{"index": float64(1), "enable": false})},
					{Config: config(service(2, "b", true)), Check: check(map[string]interface{}{"index": float64(1)})},
					{Config: config(""), Check: check(map[string]interface{}{"index": float64(2)})},
				},
			})
		})
	}
}
