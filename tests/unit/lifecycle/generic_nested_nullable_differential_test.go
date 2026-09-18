package lifecycle

import (
	"fmt"
	"strings"
	"testing"
)

// A nullable member of a list entry has the same problem as a top-level one —
// Terraform delivers `x = null` and an absent x identically — with an extra
// dependency: the configuration scan records nested attributes under the literal
// index each entry is written with. These cases run both implementations over the
// same configuration for verity_tacacs_profile.tacacs_servers.timeout, on create,
// on update, and for an entry added to an existing list.
func TestGenericMatchesLegacyOnNullableEntryMembers(t *testing.T) {
	profile := func(servers ...string) string {
		return fmt.Sprintf("resource \"verity_tacacs_profile\" \"test\" {\n  name = \"difftacacs\"\n  enable = true\n%s}\n", strings.Join(servers, ""))
	}
	server := func(index int, timeout string) string {
		timeoutLine := ""
		if timeout != "" {
			timeoutLine = "    timeout = " + timeout + "\n"
		}
		return fmt.Sprintf("  tacacs_servers {\n    index = %d\n    auth_type = \"pap\"\n    enabled = true\n    enc_secret = \"\"\n    port = \"49\"\n    secret = \"\"\n    server = \"10.0.0.%d\"\n%s  }\n", index, index, timeoutLine)
	}

	scenarios := []struct {
		name           string
		create, update string
	}{
		{"written on create", profile(server(1, "5")), profile(server(1, "5"))},
		{"written as null on create", profile(server(1, "null")), profile(server(1, "null"))},
		{"not written on create", profile(server(1, "")), profile(server(1, ""))},
		{"changed", profile(server(1, "5")), profile(server(1, "10"))},
		{"written as null on update", profile(server(1, "5")), profile(server(1, "null"))},
		{"removed from configuration", profile(server(1, "5")), profile(server(1, ""))},
		{"written as null in an added entry", profile(server(1, "5")), profile(server(1, "5"), server(2, "null"))},
		{"written as null in one of two entries", profile(server(1, "5"), server(2, "7")), profile(server(1, "5"), server(2, "null"))},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			legacy := captureLifecycle(t, "verity_tacacs_profile", false, scenario.create, scenario.update)
			generic := captureLifecycle(t, "verity_tacacs_profile", true, scenario.create, scenario.update)
			for _, operation := range []string{"PUT", "PATCH"} {
				want, got := canonical(t, legacy[operation]), canonical(t, generic[operation])
				if want != got {
					t.Errorf("%s differs between implementations\n  legacy:  %s\n  generic: %s", operation, want, got)
				}
			}
		})
	}
}
