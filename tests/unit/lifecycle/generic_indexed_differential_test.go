package lifecycle

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"testing"
)

func TestGenericMatchesLegacyOnIndexedCollections(t *testing.T) {
	macFilter := func(topType string, entries ...string) string {
		return fmt.Sprintf(`resource "verity_mac_filter" "test" {
  name = "diffmac"
  enable = true
  type = %q
%s}
`, topType, strings.Join(entries, ""))
	}
	filter := func(index, mac string, enable bool) string {
		indexLine := ""
		if index != "" {
			indexLine = "    index = " + index + "\n"
		}
		return fmt.Sprintf("  filters {\n%s    filter_num_mac = %q\n    filter_num_mask = \"ff:ff:ff:00:00:00\"\n    filter_num_enable = %v\n  }\n", indexLine, mac, enable)
	}

	packetBroker := func(ipv4Permit ...string) string {
		fixed := ""
		for _, block := range []string{"ipv4_deny", "ipv6_deny", "ipv6_permit"} {
			fixed += fmt.Sprintf("  %s {\n    index = 1\n    enable = true\n    filter = \"\"\n    filter_ref_type_ = \"\"\n  }\n", block)
		}
		return fmt.Sprintf("resource \"verity_packet_broker\" \"test\" {\n  name = \"diffpb\"\n  enable = true\n%s%s}\n", fixed, strings.Join(ipv4Permit, ""))
	}
	permit := func(index int, filterName, filterType string) string {
		return fmt.Sprintf("  ipv4_permit {\n    index = %d\n    enable = true\n    filter = %q\n    filter_ref_type_ = %q\n  }\n", index, filterName, filterType)
	}

	scenarios := []struct {
		name           string
		terraformType  string
		create, update string
		outcome        lifecycleOutcome

		unordered bool
	}{
		{
			name: "an entry member changes", terraformType: "verity_mac_filter",
			create: macFilter("", filter("1", "00:11:22:33:44:55", true)),
			update: macFilter("", filter("1", "66:77:88:99:aa:bb", true)),
		},
		{
			name: "entry members are cleared", terraformType: "verity_mac_filter",
			create: macFilter("", filter("1", "00:11:22:33:44:55", true)),
			update: macFilter("", filter("1", "", false)),
		},
		{
			name: "an entry is added with its index", terraformType: "verity_mac_filter",
			create: macFilter("", filter("1", "00:11:22:33:44:55", true)),
			update: macFilter("", filter("1", "00:11:22:33:44:55", true), filter("2", "66:77:88:99:aa:bb", false)),
		},
		{
			name: "an entry is removed", terraformType: "verity_mac_filter",
			create: macFilter("", filter("1", "00:11:22:33:44:55", true), filter("2", "66:77:88:99:aa:bb", false)),
			update: macFilter("", filter("1", "00:11:22:33:44:55", true)),
		},
		{
			name: "several entries are removed", terraformType: "verity_mac_filter",
			create:    macFilter("", filter("1", "00:11:22:33:44:55", true), filter("2", "66:77:88:99:aa:bb", false), filter("3", "aa:bb:cc:dd:ee:ff", true)),
			update:    macFilter("", filter("1", "00:11:22:33:44:55", true)),
			unordered: true,
		},
		{
			name: "every entry is removed", terraformType: "verity_mac_filter",
			create: macFilter("", filter("1", "00:11:22:33:44:55", true)),
			update: macFilter(""),
		},
		{
			name: "the list is untouched while another field changes", terraformType: "verity_mac_filter",
			create: macFilter("", filter("1", "00:11:22:33:44:55", true)),
			update: macFilter("white", filter("1", "00:11:22:33:44:55", true)),
		},
		{

			name: "entries are reordered", terraformType: "verity_mac_filter",
			create:  macFilter("", filter("1", "00:11:22:33:44:55", true), filter("2", "66:77:88:99:aa:bb", false)),
			update:  macFilter("", filter("2", "66:77:88:99:aa:bb", false), filter("1", "00:11:22:33:44:55", true)),
			outcome: lifecycleOutcome{driftAfterApply: true},
		},
		{

			name: "an entry is added without an index", terraformType: "verity_mac_filter",
			create:  macFilter("", filter("1", "00:11:22:33:44:55", true)),
			update:  macFilter("", filter("1", "00:11:22:33:44:55", true), filter("", "66:77:88:99:aa:bb", false)),
			outcome: lifecycleOutcome{applyError: regexp.MustCompile(`(?s)inconsistent result after apply`)},
		},
		{
			name: "a reference pair inside an entry changes its value", terraformType: "verity_packet_broker",
			create: packetBroker(permit(1, "filter-a", "ipv4_filter")),
			update: packetBroker(permit(1, "filter-b", "ipv4_filter")),
		},
		{
			name: "a reference pair inside an entry is cleared", terraformType: "verity_packet_broker",
			create: packetBroker(permit(1, "filter-a", "ipv4_filter")),
			update: packetBroker(permit(1, "", "")),
		},
		{
			name: "an entry is added to one of several lists", terraformType: "verity_packet_broker",
			create: packetBroker(permit(1, "filter-a", "ipv4_filter")),
			update: packetBroker(permit(1, "filter-a", "ipv4_filter"), permit(2, "filter-b", "ipv4_filter")),
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			legacy := captureLifecycle(t, scenario.terraformType, false, scenario.create, scenario.update, scenario.outcome)
			generic := captureLifecycle(t, scenario.terraformType, true, scenario.create, scenario.update, scenario.outcome)

			render := canonical
			if scenario.unordered {
				render = canonicalUnordered
			}
			for _, operation := range []string{"PUT", "PATCH"} {
				want, got := render(t, legacy[operation]), render(t, generic[operation])
				if want != got {
					t.Errorf("%s differs between implementations\n  legacy:  %s\n  generic: %s", operation, want, got)
				}
			}
		})
	}
}

func canonicalUnordered(t *testing.T, body map[string]interface{}) string {
	t.Helper()
	if body == nil {
		return "<no request>"
	}
	var sortArrays func(value interface{}) interface{}
	sortArrays = func(value interface{}) interface{} {
		switch typed := value.(type) {
		case map[string]interface{}:
			for key, member := range typed {
				typed[key] = sortArrays(member)
			}
			return typed
		case []interface{}:
			encoded := make([]string, len(typed))
			for i, element := range typed {
				bytes, _ := json.Marshal(sortArrays(element))
				encoded[i] = string(bytes)
			}
			sort.Strings(encoded)
			sorted := make([]interface{}, len(encoded))
			for i, element := range encoded {
				var decoded interface{}
				_ = json.Unmarshal([]byte(element), &decoded)
				sorted[i] = decoded
			}
			return sorted
		default:
			return value
		}
	}
	copied := map[string]interface{}{}
	raw, _ := json.Marshal(body)
	_ = json.Unmarshal(raw, &copied)
	encoded, err := json.Marshal(sortArrays(copied))
	if err != nil {
		t.Fatal(err)
	}
	return string(encoded)
}
