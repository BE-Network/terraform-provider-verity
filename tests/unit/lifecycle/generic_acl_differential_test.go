package lifecycle

import (
	"fmt"
	"testing"
)

// The ACLs are the one pair of resources that share an endpoint. Both write to
// /acls, and ip_version in the query string decides which objects a request
// touches. The golden fixtures pin each version's create in isolation, so they
// cannot show that an update or a delete still carries the right version, nor
// that the nullable port fields behave as they do in the handwritten resource.
//
// Each case runs both implementations and compares the PUT and PATCH bodies and
// the query of every write, including the DELETE the harness issues when it
// tears the case down.
func TestGenericMatchesLegacyOnACLUpdates(t *testing.T) {
	acl := func(version, body string) string {
		return fmt.Sprintf(`resource "verity_acl_v%s" "test" {
  name = "diffacl"
  bidirectional = true
  destination_ip = "10.0.0.2"
  destination_port_operator = "equal"
  enable = true
  protocol = "tcp"
  source_ip = "10.0.0.1"
  source_port_operator = "equal"
  object_properties {
    notes = "acl"
  }
%s}
`, version, body)
	}

	for _, version := range []string{"4", "6"} {
		terraformType := "verity_acl_v" + version
		scenarios := []struct {
			name           string
			create, update string
		}{
			{
				name:   "nullable ports change",
				create: acl(version, "  source_port_1 = 80\n  source_port_2 = 81\n  destination_port_1 = 443\n  destination_port_2 = 444\n"),
				update: acl(version, "  source_port_1 = 8080\n  source_port_2 = 81\n  destination_port_1 = 443\n  destination_port_2 = 444\n"),
			},
			{
				name:   "nullable port is written as null",
				create: acl(version, "  source_port_1 = 80\n  source_port_2 = 81\n  destination_port_1 = 443\n  destination_port_2 = 444\n"),
				update: acl(version, "  source_port_1 = null\n  source_port_2 = 81\n  destination_port_1 = 443\n  destination_port_2 = 444\n"),
			},
			{
				name:   "nullable port is removed from configuration",
				create: acl(version, "  source_port_1 = 80\n  source_port_2 = 81\n  destination_port_1 = 443\n  destination_port_2 = 444\n"),
				update: acl(version, "  source_port_2 = 81\n  destination_port_1 = 443\n  destination_port_2 = 444\n"),
			},
		}
		for _, scenario := range scenarios {
			t.Run(terraformType+"/"+scenario.name, func(t *testing.T) {
				legacy := captureLifecycle(t, terraformType, false, scenario.create, scenario.update)
				generic := captureLifecycle(t, terraformType, true, scenario.create, scenario.update)

				for _, key := range []string{"PUT", "PUT query", "PATCH", "PATCH query", "DELETE query"} {
					want, got := canonical(t, legacy[key]), canonical(t, generic[key])
					if want != got {
						t.Errorf("%s differs between implementations\n  legacy:  %s\n  generic: %s", key, want, got)
					}
				}
				if got := canonical(t, generic["PUT query"]); got != fmt.Sprintf(`{"ip_version":["%s"]}`, version) {
					t.Errorf("generic PUT query = %s, want ip_version %s", got, version)
				}
			})
		}
	}
}
