package utils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/zclconf/go-cty/cty"
)

var workDirRegistry sync.Map

func RegisterWorkDir(uri, dir string) {
	workDirRegistry.Store(uri, dir)
}

func UnregisterWorkDir(uri string) {
	workDirRegistry.Delete(uri)
}

func GetWorkDirForProvider(uri string) string {
	if dir, ok := workDirRegistry.Load(uri); ok {
		return dir.(string)
	}
	if dir := os.Getenv("VERITY_TF_WORKDIR"); dir != "" {
		return dir
	}
	wd, err := os.Getwd()
	if err != nil {
		return ""
	}
	return wd
}

type ConfiguredAttributes struct {
	Attributes map[string]bool

	BlockAttributes map[string]bool

	Blocks map[string]bool

	IndexedBlockAttributes map[string]map[int64]map[string]bool
}

func (c *ConfiguredAttributes) IsConfigured(attrName string) bool {
	if c == nil || c.Attributes == nil {
		return false
	}
	return c.Attributes[attrName]
}

func (c *ConfiguredAttributes) IsBlockConfigured(blockName string) bool {
	if c == nil || c.Blocks == nil {
		return false
	}
	return c.Blocks[blockName]
}

func (c *ConfiguredAttributes) IsBlockAttributeConfigured(path string) bool {
	if c == nil || c.BlockAttributes == nil {
		return false
	}
	return c.BlockAttributes[path]
}

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

func ParseResourceConfiguredAttributes(ctx context.Context, workDir string, resourceType string, resourceName string) *ConfiguredAttributes {
	result := &ConfiguredAttributes{
		Attributes:             make(map[string]bool),
		BlockAttributes:        make(map[string]bool),
		Blocks:                 make(map[string]bool),
		IndexedBlockAttributes: make(map[string]map[int64]map[string]bool),
	}

	index := getWorkDirIndex(ctx, workDir)
	if index == nil {
		return result
	}

	sanitizedName := SanitizeResourceName(resourceName)

	for _, file := range index.files {
		for _, block := range file.byType[resourceType] {
			if block.nameAttrOK {
				if block.nameAttr != resourceName {
					continue
				}
			} else if block.label != resourceName && block.label != sanitizedName {
				continue
			}

			mergeConfiguredAttributes(result, block.attrs)
			logMatchedResource(ctx, resourceType, resourceName, sanitizedName, block)
			break
		}
	}

	return result
}

func mergeConfiguredAttributes(dst, src *ConfiguredAttributes) {
	for attr := range src.Attributes {
		dst.Attributes[attr] = true
	}
	for attr := range src.BlockAttributes {
		dst.BlockAttributes[attr] = true
	}
	for block := range src.Blocks {
		dst.Blocks[block] = true
	}
	for blockType, indexMap := range src.IndexedBlockAttributes {
		if _, exists := dst.IndexedBlockAttributes[blockType]; !exists {
			dst.IndexedBlockAttributes[blockType] = make(map[int64]map[string]bool)
		}
		for idx, attrMap := range indexMap {
			if _, exists := dst.IndexedBlockAttributes[blockType][idx]; !exists {
				dst.IndexedBlockAttributes[blockType][idx] = make(map[string]bool)
			}
			for attrName := range attrMap {
				dst.IndexedBlockAttributes[blockType][idx][attrName] = true
			}
		}
	}
}

func logMatchedResource(ctx context.Context, resourceType, resourceName, sanitizedName string, block *parsedResourceBlock) {
	indexedAttrsStr := ""
	for blockType, indexMap := range block.attrs.IndexedBlockAttributes {
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
		"matched_label":       block.label,
		"sanitized_name":      sanitizedName,
		"attributes":          mapKeysToString(block.attrs.Attributes),
		"blocks":              mapKeysToString(block.attrs.Blocks),
		"indexed_block_attrs": indexedAttrsStr,
	})
}

func mapKeysToString(m map[string]bool) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return strings.Join(keys, ", ")
}

type parsedResourceBlock struct {
	label      string
	nameAttr   string
	nameAttrOK bool
	attrs      *ConfiguredAttributes
}

type parsedFile struct {
	byType map[string][]*parsedResourceBlock
}

type workDirIndex struct {
	fingerprint string
	files       []*parsedFile
}

var (
	configIndexMutex sync.Mutex
	configIndexCache = make(map[string]*workDirIndex)
)

func ClearConfigIndexCache() {
	configIndexMutex.Lock()
	defer configIndexMutex.Unlock()
	configIndexCache = make(map[string]*workDirIndex)
}

func InvalidateConfigIndex(workDir string) {
	configIndexMutex.Lock()
	defer configIndexMutex.Unlock()
	delete(configIndexCache, workDir)
}

func fingerprintTfFiles(tfFiles []string) string {
	var b strings.Builder
	for _, f := range tfFiles {
		info, err := os.Stat(f)
		if err != nil {
			fmt.Fprintf(&b, "%s|missing\n", f)
			continue
		}
		fmt.Fprintf(&b, "%s|%d|%d\n", f, info.Size(), info.ModTime().UnixNano())
	}
	return b.String()
}

func getWorkDirIndex(ctx context.Context, workDir string) *workDirIndex {
	tfFiles, err := filepath.Glob(filepath.Join(workDir, "*.tf"))
	if err != nil {
		tflog.Warn(ctx, "Failed to find .tf files", map[string]interface{}{
			"error":   err.Error(),
			"workDir": workDir,
		})
		return nil
	}

	fingerprint := fingerprintTfFiles(tfFiles)

	configIndexMutex.Lock()
	defer configIndexMutex.Unlock()

	if cached, ok := configIndexCache[workDir]; ok && cached.fingerprint == fingerprint {
		return cached
	}

	index := buildWorkDirIndex(ctx, tfFiles, fingerprint)
	configIndexCache[workDir] = index
	return index
}

func buildWorkDirIndex(ctx context.Context, tfFiles []string, fingerprint string) *workDirIndex {
	index := &workDirIndex{
		fingerprint: fingerprint,
		files:       make([]*parsedFile, 0, len(tfFiles)),
	}

	for _, tfFile := range tfFiles {
		src, err := os.ReadFile(tfFile)
		if err != nil {
			tflog.Debug(ctx, "Failed to read .tf file", map[string]interface{}{
				"file":  tfFile,
				"error": err.Error(),
			})
			continue
		}

		file, diags := hclparse.NewParser().ParseHCL(src, tfFile)
		if diags.HasErrors() {
			tflog.Debug(ctx, "Failed to parse .tf file", map[string]interface{}{
				"file":  tfFile,
				"error": diags.Error(),
			})
			continue
		}

		if parsed := indexResourceBlocks(ctx, file.Body); parsed != nil {
			index.files = append(index.files, parsed)
		}
	}

	tflog.Debug(ctx, "HCL parser built configuration index", map[string]interface{}{
		"files": len(index.files),
	})

	return index
}

func indexResourceBlocks(ctx context.Context, body hcl.Body) *parsedFile {
	syntaxBody, ok := body.(*hclsyntax.Body)
	if !ok {
		tflog.Warn(ctx, "HCL parser: body is not hclsyntax.Body, cannot parse")
		return nil
	}

	parsed := &parsedFile{byType: make(map[string][]*parsedResourceBlock)}

	for _, block := range syntaxBody.Blocks {
		if block.Type != "resource" || len(block.Labels) < 2 {
			continue
		}

		entry := &parsedResourceBlock{
			label: block.Labels[1],
			attrs: extractBlockAttributes(block.Body),
		}

		if nameAttr, hasName := block.Body.Attributes["name"]; hasName {
			val, diags := nameAttr.Expr.Value(nil)
			if !diags.HasErrors() && val.Type() == cty.String && !val.IsNull() && val.IsKnown() {
				entry.nameAttr = val.AsString()
				entry.nameAttrOK = true
			}
		}

		resourceType := block.Labels[0]
		parsed.byType[resourceType] = append(parsed.byType[resourceType], entry)
	}

	return parsed
}

func extractBlockAttributes(body *hclsyntax.Body) *ConfiguredAttributes {
	result := &ConfiguredAttributes{
		Attributes:             make(map[string]bool),
		BlockAttributes:        make(map[string]bool),
		Blocks:                 make(map[string]bool),
		IndexedBlockAttributes: make(map[string]map[int64]map[string]bool),
	}

	for name := range body.Attributes {
		result.Attributes[name] = true
	}

	extractNestedBlocks(body.Blocks, result, "")

	return result
}

func extractNestedBlocks(blocks []*hclsyntax.Block, result *ConfiguredAttributes, prefix string) {
	for _, block := range blocks {
		fullBlockPath := block.Type
		if prefix != "" {
			fullBlockPath = prefix + "." + block.Type
		}

		result.Blocks[block.Type] = true
		result.Blocks[fullBlockPath] = true

		var blockIndex int64 = -1
		if indexAttr, hasIndex := block.Body.Attributes["index"]; hasIndex {
			val, diags := indexAttr.Expr.Value(nil)
			if diags.HasErrors() {

				for _, d := range diags {
					_ = d
				}
			} else if val.Type() == cty.Number && !val.IsNull() && val.IsKnown() {
				bf := val.AsBigFloat()
				if bf.IsInt() {
					blockIndex, _ = bf.Int64()
				}
			}
		}

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

		for attrName := range block.Body.Attributes {
			simplePath := block.Type + "." + attrName
			result.BlockAttributes[simplePath] = true

			if prefix != "" {
				fullPath := fullBlockPath + "." + attrName
				result.BlockAttributes[fullPath] = true
			}
		}

		if len(block.Body.Blocks) > 0 {
			extractNestedBlocks(block.Body.Blocks, result, fullBlockPath)
		}
	}
}
