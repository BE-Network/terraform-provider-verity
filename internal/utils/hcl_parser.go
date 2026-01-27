package utils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zclconf/go-cty/cty"
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
	// IndexedBlockAttributes maps: blockType -> indexValue -> attrName -> true
	// Used to track which attributes are configured for each indexed block instance
	IndexedBlockAttributes map[string]map[int64]map[string]bool
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

// IsIndexedBlockAttributeConfigured returns true if an attribute within a specific
// indexed block instance was configured. Used for blocks with an "index" field.
func (c *ConfiguredAttributes) IsIndexedBlockAttributeConfigured(blockType string, indexValue int64, attrName string) bool {
	if c == nil || c.IndexedBlockAttributes == nil {
		return false
	}
	if indexMap, ok := c.IndexedBlockAttributes[blockType]; ok {
		if attrMap, ok := indexMap[indexValue]; ok {
			return attrMap[attrName]
		}
	}
	return false
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
		Attributes:             make(map[string]bool),
		BlockAttributes:        make(map[string]bool),
		Blocks:                 make(map[string]bool),
		IndexedBlockAttributes: make(map[string]map[int64]map[string]bool),
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
			// Merge IndexedBlockAttributes
			for blockType, indexMap := range found.IndexedBlockAttributes {
				if _, exists := result.IndexedBlockAttributes[blockType]; !exists {
					result.IndexedBlockAttributes[blockType] = make(map[int64]map[string]bool)
				}
				for idx, attrMap := range indexMap {
					if _, exists := result.IndexedBlockAttributes[blockType][idx]; !exists {
						result.IndexedBlockAttributes[blockType][idx] = make(map[string]bool)
					}
					for attrName := range attrMap {
						result.IndexedBlockAttributes[blockType][idx][attrName] = true
					}
				}
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
					Attributes:             make(map[string]bool),
					BlockAttributes:        make(map[string]bool),
					Blocks:                 make(map[string]bool),
					IndexedBlockAttributes: make(map[string]map[int64]map[string]bool),
				}

				// Extract top-level attributes
				for name := range block.Body.Attributes {
					result.Attributes[name] = true
				}

				// Extract nested blocks and their attributes
				extractNestedBlocks(block.Body.Blocks, result, "")

				// Build a string representation of IndexedBlockAttributes for debugging
				indexedAttrsStr := ""
				for blockType, indexMap := range result.IndexedBlockAttributes {
					for idx, attrMap := range indexMap {
						attrs := make([]string, 0, len(attrMap))
						for attr := range attrMap {
							attrs = append(attrs, attr)
						}
						indexedAttrsStr += fmt.Sprintf("%s[%d]:{%s} ", blockType, idx, strings.Join(attrs, ","))
					}
				}

				tflog.Debug(ctx, "HCL parser found resource", map[string]interface{}{
					"resource":            resourceType + "." + resourceName,
					"matched_label":       blockName,
					"sanitized_name":      sanitizedName,
					"attributes":          mapKeysToString(result.Attributes),
					"blocks":              mapKeysToString(result.Blocks),
					"indexed_block_attrs": indexedAttrsStr,
				})
				return result
			}
		}
	}

	return nil
}

// extractNestedBlocks recursively extracts block names and their attributes.
// For blocks with an "index" attribute, it also populates IndexedBlockAttributes
// to track which attributes are configured per-block-index.
func extractNestedBlocks(blocks []*hclsyntax.Block, result *ConfiguredAttributes, prefix string) {
	for _, block := range blocks {
		fullBlockPath := block.Type
		if prefix != "" {
			fullBlockPath = prefix + "." + block.Type
		}

		// Record block names
		result.Blocks[block.Type] = true
		result.Blocks[fullBlockPath] = true

		// Try to extract the "index" attribute value for indexed blocks
		var blockIndex int64 = -1
		if indexAttr, hasIndex := block.Body.Attributes["index"]; hasIndex {
			val, diags := indexAttr.Expr.Value(nil)
			if diags.HasErrors() {
				// Log the error but continue
				for _, d := range diags {
					_ = d // can't log without context
				}
			} else if val.Type() == cty.Number && !val.IsNull() && val.IsKnown() {
				bf := val.AsBigFloat()
				if bf.IsInt() {
					blockIndex, _ = bf.Int64()
				}
			}
		}

		// If we have a valid index, populate IndexedBlockAttributes
		if blockIndex >= 0 {
			blockType := block.Type
			if prefix != "" {
				blockType = fullBlockPath
			}
			if _, exists := result.IndexedBlockAttributes[blockType]; !exists {
				result.IndexedBlockAttributes[blockType] = make(map[int64]map[string]bool)
			}
			if _, exists := result.IndexedBlockAttributes[blockType][blockIndex]; !exists {
				result.IndexedBlockAttributes[blockType][blockIndex] = make(map[string]bool)
			}
			for attrName := range block.Body.Attributes {
				result.IndexedBlockAttributes[blockType][blockIndex][attrName] = true
			}
		}

		// Extract attributes from this block (existing behavior)
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
