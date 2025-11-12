package provider

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"terraform-provider-verity/internal/importer"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &stateImporterDataSource{}
	_ datasource.DataSourceWithConfigure = &stateImporterDataSource{}
)

func NewVerityStateImporterDataSource() datasource.DataSource {
	return &stateImporterDataSource{}
}

type stateImporterDataSource struct {
	client *providerContext
}

type stateImporterDataSourceModel struct {
	ID            types.String   `tfsdk:"id"`
	OutputDir     types.String   `tfsdk:"output_dir"`
	ImportedFiles []types.String `tfsdk:"imported_files"`
}

func (d *stateImporterDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_state_importer"
}

func (d *stateImporterDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Data source for importing existing resources into Terraform state",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "Identifier for this import operation",
				Computed:    true,
			},
			"output_dir": schema.StringAttribute{
				Description: "Directory where the TF files will be saved. Defaults to current directory.",
				Optional:    true,
			},
			"imported_files": schema.ListAttribute{
				Description: "List of files that were created during import",
				Computed:    true,
				ElementType: types.StringType,
			},
		},
	}
}

func (d *stateImporterDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*providerContext)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *providerContext, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *stateImporterDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data stateImporterDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	outputDir := data.OutputDir.ValueString()
	if outputDir == "" {
		currentDir, err := os.Getwd()
		if err != nil {
			resp.Diagnostics.AddError("Error getting current directory", err.Error())
			return
		}
		outputDir = currentDir
		data.OutputDir = types.StringValue(outputDir)
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Output Directory",
			fmt.Sprintf("Error creating output directory: %v", err),
		)
		return
	}

	absPath, err := filepath.Abs(outputDir)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Getting Absolute Path",
			fmt.Sprintf("Error getting absolute path: %v", err),
		)
		return
	}

	client := d.client.client
	imp := importer.NewImporter(client, d.client.mode)
	err = imp.ImportAll(absPath)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Resources",
			fmt.Sprintf("Error importing resources: %v", err),
		)
		return
	}

	data.ImportedFiles = []types.String{}

	entries, err := os.ReadDir(absPath)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Directory",
			fmt.Sprintf("Error reading directory: %v", err),
		)
		return
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".tf") {
			filePath := filepath.Join(absPath, entry.Name())
			data.ImportedFiles = append(data.ImportedFiles, types.StringValue(filePath))
		}
	}

	importBlocksFile, err := createImportBlocks(ctx, absPath)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Generating Import Blocks",
			fmt.Sprintf("Error generating import blocks: %v", err),
		)
		return
	}
	data.ImportedFiles = append(data.ImportedFiles, types.StringValue(importBlocksFile))
	tflog.Info(ctx, "Successfully generated import blocks", map[string]any{
		"file": importBlocksFile,
	})

	data.ID = types.StringValue(time.Now().UTC().String())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func createImportBlocks(ctx context.Context, dirPath string) (string, error) {
	outputFile := filepath.Join(dirPath, "import_blocks.tf")

	file, err := os.Create(outputFile)
	if err != nil {
		return "", fmt.Errorf("error creating output file: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString("# Import blocks for Verity resources\n\n"); err != nil {
		return "", fmt.Errorf("error writing to output file: %w", err)
	}

	supportedResources := map[string]struct{}{
		"verity_service":                  {},
		"verity_eth_port_profile":         {},
		"verity_authenticated_eth_port":   {},
		"verity_device_voice_settings":    {},
		"verity_packet_queue":             {},
		"verity_service_port_profile":     {},
		"verity_voice_port_profile":       {},
		"verity_eth_port_settings":        {},
		"verity_device_settings":          {},
		"verity_lag":                      {},
		"verity_sflow_collector":          {},
		"verity_diagnostics_profile":      {},
		"verity_diagnostics_port_profile": {},
		"verity_bundle":                   {},
		"verity_acl_v4":                   {},
		"verity_acl_v6":                   {},
		"verity_ipv4_list":                {},
		"verity_ipv6_list":                {},
		"verity_port_acl":                 {},
		"verity_badge":                    {},
		"verity_switchpoint":              {},
		"verity_device_controller":        {},
		"verity_site":                     {},
		"verity_tenant":                   {},
		"verity_gateway_profile":          {},
		"verity_gateway":                  {},
		"verity_ipv4_prefix_list":         {},
		"verity_ipv6_prefix_list":         {},
		"verity_packet_broker":            {},
		"verity_pod":                      {},
		"verity_as_path_access_list":      {},
		"verity_community_list":           {},
		"verity_extended_community_list":  {},
		"verity_route_map_clause":         {},
		"verity_route_map":                {},
		"verity_pb_routing":               {},
		"verity_pb_routing_acl":           {},
		"verity_spine_plane":              {},
		"verity_sfp_breakout":             {},
		"verity_grouping_rule":            {},
		"verity_threshold_group":          {},
		"verity_threshold":                {},
	}

	importBlocks := make(map[string][]string)

	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return "", fmt.Errorf("error reading directory: %w", err)
	}

	resourceRegex := regexp.MustCompile(`resource\s+"([^"]+)"\s+"([^"]+)"`)
	nameRegex := regexp.MustCompile(`name\s*=\s*"([^"]+)"`)

	tflog.Info(ctx, "Processing Terraform files for import blocks")

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".tf") {
			continue
		}

		if entry.Name() == "import_blocks.tf" {
			continue
		}

		filePath := filepath.Join(dirPath, entry.Name())

		hasVerityResources, err := containsVerityResources(filePath)
		if err != nil {
			return "", fmt.Errorf("error checking file %s: %w", entry.Name(), err)
		}

		if !hasVerityResources {
			tflog.Info(ctx, "Skipping file (no Verity resources)", map[string]any{"file": entry.Name()})
			continue
		}

		tflog.Info(ctx, "Processing file", map[string]any{"file": entry.Name()})

		resourceBlocks, err := findResourceBlocks(filePath)
		if err != nil {
			return "", fmt.Errorf("error parsing file %s: %w", entry.Name(), err)
		}

		for _, block := range resourceBlocks {
			resourceMatches := resourceRegex.FindStringSubmatch(block)
			if len(resourceMatches) < 3 {
				continue
			}

			resourceType := resourceMatches[1]

			// Only process supported Verity resources
			if _, isSupported := supportedResources[resourceType]; !isSupported {
				tflog.Debug(ctx, "Skipping unsupported resource type", map[string]any{
					"resource_type": resourceType,
					"file":          entry.Name(),
				})
				continue
			}

			hclName := resourceMatches[2]

			nameValue := hclName
			nameMatches := nameRegex.FindStringSubmatch(block)
			if len(nameMatches) > 1 {
				nameValue = nameMatches[1]
			}

			importBlock := fmt.Sprintf("import {\n  to = %s.%s\n  id = \"%s\"\n}\n",
				resourceType, hclName, nameValue)

			importBlocks[resourceType] = append(importBlocks[resourceType], importBlock)
		}
	}

	// Write all collected import blocks grouped by resource type
	for resourceType, blocks := range importBlocks {
		if len(blocks) > 0 {
			if _, err := file.WriteString(fmt.Sprintf("# %s imports\n", resourceType)); err != nil {
				return "", fmt.Errorf("error writing to output file: %w", err)
			}
			for _, block := range blocks {
				if _, err := file.WriteString(block + "\n"); err != nil {
					return "", fmt.Errorf("error writing to output file: %w", err)
				}
			}
		}
	}

	return outputFile, nil
}

func containsVerityResources(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	verityResourcePattern := regexp.MustCompile(`resource\s+"verity_`)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if verityResourcePattern.MatchString(line) {
			return true, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return false, fmt.Errorf("error reading file: %w", err)
	}

	return false, nil
}

func findResourceBlocks(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	var blocks []string
	var currentBlock strings.Builder
	inBlock := false
	braceCount := 0

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if !inBlock && strings.Contains(line, "resource ") {
			inBlock = true
			currentBlock.WriteString(line + "\n")
			braceCount += strings.Count(line, "{") - strings.Count(line, "}")
			continue
		}

		if inBlock {
			currentBlock.WriteString(line + "\n")
			braceCount += strings.Count(line, "{") - strings.Count(line, "}")

			if braceCount == 0 {
				blocks = append(blocks, currentBlock.String())
				currentBlock.Reset()
				inBlock = false
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file: %w", err)
	}

	return blocks, nil
}
