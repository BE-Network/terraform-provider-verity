package utils

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ConfiguredAttributes holds the set of attributes that were explicitly written
// in the .tf file for a specific resource.
type ConfiguredAttributes struct {
	// Attributes contains top-level attribute names
	Attributes map[string]bool
	// BlockAttributes contains attributes within nested blocks, keyed by "blockName.attrName"
	// For repeated blocks, uses "blockName.attrName" without index
	BlockAttributes map[string]bool
	// Blocks contains the names of blocks that are present
	Blocks map[string]bool
}

// IsConfigured returns true if the attribute was explicitly written in the .tf file.
func (c *ConfiguredAttributes) IsConfigured(attrName string) bool {
	if c == nil || c.Attributes == nil {
		return false
	}
	return c.Attributes[attrName]
}

// IsBlockConfigured returns true if the block was explicitly written in the .tf file.
func (c *ConfiguredAttributes) IsBlockConfigured(blockName string) bool {
	if c == nil || c.Blocks == nil {
		return false
	}
	return c.Blocks[blockName]
}

// IsBlockAttributeConfigured returns true if an attribute within a nested block was configured.
func (c *ConfiguredAttributes) IsBlockAttributeConfigured(path string) bool {
	if c == nil || c.BlockAttributes == nil {
		return false
	}
	return c.BlockAttributes[path]
}

// ParseResourceConfiguredAttributes parses all .tf files in the given directory
// and returns the set of attributes that were explicitly configured for the
// specified resource type and name.
//
// This is the ONLY reliable way to distinguish between:
//   - Field not in .tf file → omit from API request (use server default)
//   - `field = null` in .tf file → send null to API (explicit null)
//   - `field = value` in .tf file → send value to API
//
// The terraform plugin framework does not provide this distinction - it fills
// all schema attributes with null for unspecified fields.
func ParseResourceConfiguredAttributes(ctx context.Context, workDir string, resourceType string, resourceName string) *ConfiguredAttributes {
	result := &ConfiguredAttributes{
		Attributes:      make(map[string]bool),
		BlockAttributes: make(map[string]bool),
		Blocks:          make(map[string]bool),
	}

	// Find all .tf files in the working directory
	tfFiles, err := filepath.Glob(filepath.Join(workDir, "*.tf"))
	if err != nil {
		tflog.Warn(ctx, "Failed to find .tf files", map[string]interface{}{
			"error":   err.Error(),
			"workDir": workDir,
		})
		return result
	}

	parser := hclparse.NewParser()

	for _, tfFile := range tfFiles {
		src, err := os.ReadFile(tfFile)
		if err != nil {
			tflog.Debug(ctx, "Failed to read .tf file", map[string]interface{}{
				"file":  tfFile,
				"error": err.Error(),
			})
			continue
		}

		file, diags := parser.ParseHCL(src, tfFile)
		if diags.HasErrors() {
			tflog.Debug(ctx, "Failed to parse .tf file", map[string]interface{}{
				"file":  tfFile,
				"error": diags.Error(),
			})
			continue
		}

		// Find the resource block and extract all configured attributes
		found := findResourceAttributes(ctx, file.Body, resourceType, resourceName)
		if found != nil {
			// Merge found attributes
			for attr := range found.Attributes {
				result.Attributes[attr] = true
			}
			for attr := range found.BlockAttributes {
				result.BlockAttributes[attr] = true
			}
			for block := range found.Blocks {
				result.Blocks[block] = true
			}
		}
	}

	return result
}

// findResourceAttributes looks for a resource block with the given type and name
// and returns all attributes and blocks that are explicitly set in it.
// Returns nil if the resource is not found or cannot be fully parsed.
func findResourceAttributes(ctx context.Context, body hcl.Body, resourceType string, resourceName string) *ConfiguredAttributes {
	syntaxBody, ok := body.(*hclsyntax.Body)
	if !ok {
		tflog.Warn(ctx, "HCL parser: body is not hclsyntax.Body, cannot parse")
		return nil
	}

	// Also try the sanitized version of the name, since HCL resource labels
	// use sanitized names while the API uses the original names.
	// Example: API name "(Device Settings)" -> HCL label "_Device_Settings_"
	sanitizedName := SanitizeResourceName(resourceName)

	for _, block := range syntaxBody.Blocks {
		if block.Type == "resource" && len(block.Labels) >= 2 {
			blockType := block.Labels[0]
			blockName := block.Labels[1]

			// Match either the exact name or the sanitized version
			if blockType == resourceType && (blockName == resourceName || blockName == sanitizedName) {
				result := &ConfiguredAttributes{
					Attributes:      make(map[string]bool),
					BlockAttributes: make(map[string]bool),
					Blocks:          make(map[string]bool),
				}

				// Extract top-level attributes
				for name := range block.Body.Attributes {
					result.Attributes[name] = true
				}

				// Extract nested blocks and their attributes
				extractNestedBlocks(block.Body.Blocks, result, "")

				tflog.Debug(ctx, "HCL parser found resource", map[string]interface{}{
					"resource":       resourceType + "." + resourceName,
					"matched_label":  blockName,
					"sanitized_name": sanitizedName,
					"attributes":     mapKeysToString(result.Attributes),
					"blocks":         mapKeysToString(result.Blocks),
				})
				return result
			}
		}
	}

	return nil
}

// extractNestedBlocks recursively extracts block names and their attributes.
func extractNestedBlocks(blocks []*hclsyntax.Block, result *ConfiguredAttributes, prefix string) {
	for _, block := range blocks {
		fullBlockPath := block.Type
		if prefix != "" {
			fullBlockPath = prefix + "." + block.Type
		}

		// Record block names
		result.Blocks[block.Type] = true
		result.Blocks[fullBlockPath] = true

		// Extract attributes from this block
		for attrName := range block.Body.Attributes {
			simplePath := block.Type + "." + attrName
			result.BlockAttributes[simplePath] = true

			if prefix != "" {
				fullPath := fullBlockPath + "." + attrName
				result.BlockAttributes[fullPath] = true
			}
		}

		// Recursively process nested blocks
		if len(block.Body.Blocks) > 0 {
			extractNestedBlocks(block.Body.Blocks, result, fullBlockPath)
		}
	}
}

func mapKeysToString(m map[string]bool) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}

func GetWorkingDirectory() string {
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return wd
}
