package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func selectedVersion(path string) (string, error) {
	overrides, err := readOverrides(path)
	if err != nil {
		return "", err
	}
	if overrides.FormatVersion != registryFormatVersion {
		return "", fmt.Errorf("unsupported overrides format_version %d", overrides.FormatVersion)
	}
	version, err := parseAPIVersion(overrides.APIVersion)
	if err != nil {
		return "", err
	}
	return version.String(), nil
}

func verifyVersion(inputDir, expected string) error {
	if expected != "" {
		if _, err := parseAPIVersion(expected); err != nil {
			return err
		}
	}
	if err := verify(inputDir); err != nil {
		return err
	}
	if expected == "" {
		return nil
	}
	data, err := os.ReadFile(filepath.Join(inputDir, "manifest.json"))
	if err != nil {
		return err
	}
	var inputs manifest
	if err := json.Unmarshal(data, &inputs); err != nil {
		return err
	}
	if inputs.APIVersion != expected {
		return fmt.Errorf("input API version %q does not match selected version %q", inputs.APIVersion, expected)
	}
	return nil
}
