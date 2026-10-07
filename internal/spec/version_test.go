package spec

import "testing"

func TestParseAPIVersion(t *testing.T) {
	for _, value := range []string{"0.0", "6.6", "6.10", "12.34"} {
		version, err := ParseAPIVersion(value)
		if err != nil || version.String() != value {
			t.Errorf("ParseAPIVersion(%q) = %v, %v", value, version, err)
		}
	}
	for _, value := range []string{"", "6", "6.6.1", "-1.6", "6.-1", "+6.6", "06.6", "6.06", "6.6/../../", "6.x"} {
		if _, err := ParseAPIVersion(value); err == nil {
			t.Errorf("ParseAPIVersion(%q) accepted noncanonical version", value)
		}
	}
}
