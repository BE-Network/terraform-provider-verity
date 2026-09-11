package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const normalizationToolVersion = "1"

type sourceEntry struct {
	Mode             string `json:"mode"`
	File             string `json:"file"`
	SHA256           string `json:"sha256"`
	SourceSHA256     string `json:"source_sha256"`
	SourceExportDate string `json:"source_export_date"`
	Provenance       string `json:"provenance"`
}

type manifest struct {
	FormatVersion     int           `json:"format_version"`
	APIVersion        string        `json:"api_version"`
	NormalizationTool string        `json:"normalization_tool"`
	NormalizationVer  string        `json:"normalization_tool_version"`
	VolatileMetadata  []string      `json:"volatile_metadata_allowlist"`
	Sources           []sourceEntry `json:"sources"`
}

type normalizeOptions struct {
	Version          string
	Datacenter       string
	Campus           string
	OutputDir        string
	SourceExportDate string
	Provenance       string
}

func main() {
	if len(os.Args) < 2 {
		fail("usage: specgen <normalize|verify|extract|registry|metadata> [flags]")
	}

	switch os.Args[1] {
	case "normalize":
		fs := flag.NewFlagSet("normalize", flag.ExitOnError)
		opts := normalizeOptions{}
		fs.StringVar(&opts.Version, "version", "", "API version, for example 6.6")
		fs.StringVar(&opts.Datacenter, "datacenter", "", "datacenter OpenAPI JSON export")
		fs.StringVar(&opts.Campus, "campus", "", "campus OpenAPI JSON export")
		fs.StringVar(&opts.OutputDir, "output-dir", "", "canonical input directory")
		fs.StringVar(&opts.SourceExportDate, "source-export-date", "", "source export date (YYYY-MM-DD)")
		fs.StringVar(&opts.Provenance, "provenance", "", "source provenance note")
		_ = fs.Parse(os.Args[2:])
		if err := normalize(opts); err != nil {
			fail(err.Error())
		}
	case "verify":
		fs := flag.NewFlagSet("verify", flag.ExitOnError)
		inputDir := fs.String("input-dir", "", "canonical input directory")
		_ = fs.Parse(os.Args[2:])
		if err := verify(*inputDir); err != nil {
			fail(err.Error())
		}
	case "extract":
		fs := flag.NewFlagSet("extract", flag.ExitOnError)
		opts := extractOptions{}
		fs.StringVar(&opts.InputDir, "input-dir", "", "canonical input directory")
		fs.StringVar(&opts.Output, "output", "", "coverage manifest output path")
		fs.BoolVar(&opts.Check, "check", false, "fail if output differs from deterministic extraction")
		_ = fs.Parse(os.Args[2:])
		if err := extract(opts); err != nil {
			fail(err.Error())
		}
	case "registry":
		fs := flag.NewFlagSet("registry", flag.ExitOnError)
		opts := registryOptions{}
		fs.StringVar(&opts.InputDir, "input-dir", "", "canonical input directory")
		fs.StringVar(&opts.Overrides, "overrides", "", "reviewed override YAML file")
		fs.StringVar(&opts.Output, "output", "", "generated registry output path")
		fs.BoolVar(&opts.Check, "check", false, "fail if output differs from deterministic generation")
		_ = fs.Parse(os.Args[2:])
		if err := generateRegistry(opts); err != nil {
			fail(err.Error())
		}
	case "metadata":
		fs := flag.NewFlagSet("metadata", flag.ExitOnError)
		opts := metadataOptions{}
		fs.StringVar(&opts.Registry, "registry", "", "generated registry input path")
		fs.StringVar(&opts.Output, "output", "", "generated Go metadata output path")
		fs.StringVar(&opts.BulkOutput, "bulk-output", "", "generated bulk metadata output path")
		fs.BoolVar(&opts.Check, "check", false, "fail if output differs from deterministic generation")
		_ = fs.Parse(os.Args[2:])
		if err := generateModeMetadata(opts); err != nil {
			fail(err.Error())
		}
	default:
		fail(fmt.Sprintf("unknown command %q; expected normalize, verify, extract, registry, or metadata", os.Args[1]))
	}
}

func fail(message string) {
	fmt.Fprintln(os.Stderr, "specgen:", message)
	os.Exit(1)
}

func normalize(opts normalizeOptions) error {
	if opts.Version == "" || opts.Datacenter == "" || opts.Campus == "" || opts.OutputDir == "" || opts.SourceExportDate == "" || opts.Provenance == "" {
		return errors.New("--version, --datacenter, --campus, --output-dir, --source-export-date, and --provenance are required")
	}
	if _, err := time.Parse("2006-01-02", opts.SourceExportDate); err != nil {
		return fmt.Errorf("invalid --source-export-date: %w", err)
	}
	if err := os.MkdirAll(opts.OutputDir, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}

	sources := []struct {
		mode string
		path string
	}{
		{mode: "campus", path: opts.Campus},
		{mode: "datacenter", path: opts.Datacenter},
	}
	entries := make([]sourceEntry, 0, len(sources))
	for _, source := range sources {
		raw, err := os.ReadFile(source.path)
		if err != nil {
			return fmt.Errorf("read %s source: %w", source.mode, err)
		}
		canonical, err := canonicalJSON(raw)
		if err != nil {
			return fmt.Errorf("normalize %s source: %w", source.mode, err)
		}

		fileName := source.mode + ".json"
		if err := writeFile(filepath.Join(opts.OutputDir, fileName), canonical); err != nil {
			return err
		}
		entries = append(entries, sourceEntry{
			Mode:             source.mode,
			File:             fileName,
			SHA256:           sha256Hex(canonical),
			SourceSHA256:     sha256Hex(raw),
			SourceExportDate: opts.SourceExportDate,
			Provenance:       opts.Provenance,
		})
	}

	encoded, err := marshalCanonical(manifest{
		FormatVersion:     1,
		APIVersion:        opts.Version,
		NormalizationTool: "tools/specgen",
		NormalizationVer:  normalizationToolVersion,
		VolatileMetadata:  []string{},
		Sources:           entries,
	})
	if err != nil {
		return fmt.Errorf("encode manifest: %w", err)
	}
	if err := writeFile(filepath.Join(opts.OutputDir, "manifest.json"), encoded); err != nil {
		return err
	}
	return verify(opts.OutputDir)
}

func verify(inputDir string) error {
	if inputDir == "" {
		return errors.New("--input-dir is required")
	}
	rawManifest, err := os.ReadFile(filepath.Join(inputDir, "manifest.json"))
	if err != nil {
		return fmt.Errorf("read manifest: %w", err)
	}
	canonicalManifest, err := canonicalJSON(rawManifest)
	if err != nil {
		return fmt.Errorf("manifest is not valid JSON: %w", err)
	}
	if !bytes.Equal(rawManifest, canonicalManifest) {
		return errors.New("manifest.json is not canonically formatted; rerun specgen normalize")
	}

	var current manifest
	if err := json.Unmarshal(rawManifest, &current); err != nil {
		return fmt.Errorf("decode manifest: %w", err)
	}
	if current.FormatVersion != 1 || current.APIVersion == "" || current.NormalizationTool != "tools/specgen" || current.NormalizationVer != normalizationToolVersion {
		return errors.New("manifest has an unsupported format, missing API version, or unexpected normalization tool version")
	}
	if len(current.VolatileMetadata) != 0 {
		return errors.New("volatile metadata removal is not implemented; manifest allowlist must be empty")
	}
	if len(current.Sources) != 2 {
		return errors.New("manifest must contain exactly datacenter and campus sources")
	}

	expectedModes := []string{"campus", "datacenter"}
	for index, entry := range current.Sources {
		if entry.Mode != expectedModes[index] || entry.File != entry.Mode+".json" {
			return fmt.Errorf("manifest source %d must be %s/%s.json", index, expectedModes[index], expectedModes[index])
		}
		if entry.SourceSHA256 == "" || entry.SourceExportDate == "" || entry.Provenance == "" {
			return fmt.Errorf("manifest source %q is missing source provenance", entry.Mode)
		}
		if _, err := time.Parse("2006-01-02", entry.SourceExportDate); err != nil {
			return fmt.Errorf("manifest source %q has invalid source_export_date: %w", entry.Mode, err)
		}

		raw, err := os.ReadFile(filepath.Join(inputDir, entry.File))
		if err != nil {
			return fmt.Errorf("read %s: %w", entry.File, err)
		}
		canonical, err := canonicalJSON(raw)
		if err != nil {
			return fmt.Errorf("validate %s: %w", entry.File, err)
		}
		if !bytes.Equal(raw, canonical) {
			return fmt.Errorf("%s is not canonically formatted; rerun specgen normalize", entry.File)
		}
		if sha256Hex(raw) != entry.SHA256 {
			return fmt.Errorf("checksum mismatch for %s", entry.File)
		}
	}
	return nil
}

func canonicalJSON(input []byte) ([]byte, error) {
	decoder := json.NewDecoder(bytes.NewReader(input))
	decoder.UseNumber()
	var value any
	if err := decoder.Decode(&value); err != nil {
		return nil, err
	}
	var extra any
	err := decoder.Decode(&extra)
	if !errors.Is(err, io.EOF) {
		if err == nil {
			return nil, errors.New("multiple JSON values are not allowed")
		}
		return nil, err
	}
	return encodeIndented(value)
}

func marshalCanonical(value any) ([]byte, error) {
	raw, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return canonicalJSON(raw)
}

func encodeIndented(value any) ([]byte, error) {
	var output bytes.Buffer
	encoder := json.NewEncoder(&output)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(value); err != nil {
		return nil, err
	}
	return output.Bytes(), nil
}

func writeFile(path string, contents []byte) error {
	temporary, err := os.CreateTemp(filepath.Dir(path), ".specgen-*")
	if err != nil {
		return fmt.Errorf("create temporary file for %s: %w", path, err)
	}
	temporaryPath := temporary.Name()
	defer os.Remove(temporaryPath)
	if _, err := temporary.Write(contents); err != nil {
		temporary.Close()
		return fmt.Errorf("write %s: %w", path, err)
	}
	if err := temporary.Chmod(0o644); err != nil {
		temporary.Close()
		return fmt.Errorf("set permissions on %s: %w", path, err)
	}
	if err := temporary.Close(); err != nil {
		return fmt.Errorf("close %s: %w", path, err)
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		return fmt.Errorf("replace %s: %w", path, err)
	}
	return nil
}

func sha256Hex(value []byte) string {
	sum := sha256.Sum256(value)
	return hex.EncodeToString(sum[:])
}
