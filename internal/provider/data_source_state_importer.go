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
	skipFiles := map[string]bool{
		"provider.tf":      true,
		"versions.tf":      true,
		"variables.tf":     true,
		"terraform.tfvars": true,
		"import_blocks.tf": true,
	}

	outputFile := filepath.Join(dirPath, "import_blocks.tf")

	file, err := os.Create(outputFile)
	if err != nil {
		return "", fmt.Errorf("error creating output file: %w", err)
	}
	defer file.Close()

	if _, err := file.WriteString("# Import blocks for Verity resources\n\n"); err != nil {
		return "", fmt.Errorf("error writing to output file: %w", err)
	}

	resourceOrder := []string{
		// Datacenter mode resources
		"verity_tenant",
		"verity_service",
		"verity_eth_port_settings",
		"verity_eth_port_profile",
		"verity_gateway_profile",
		"verity_gateway",
		"verity_lag",
		"verity_bundle",
		"verity_acl",
		"verity_packet_broker",
		// Campus mode resources
		"verity_authenticated_eth_port",
		"verity_device_voice_settings",
		"verity_packet_queue",
		"verity_service_port_profile",
		"verity_voice_port_profile",
		// Both modes
		"verity_badge",
		"verity_switchpoint",
		"verity_device_controller",
	}

	importBlocks := make(map[string]string)
	for _, res := range resourceOrder {
		importBlocks[res] = ""
	}

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

		if skipFiles[entry.Name()] {
			tflog.Info(ctx, "Skipping file", map[string]any{"file": entry.Name()})
			continue
		}

		tflog.Info(ctx, "Processing file", map[string]any{"file": entry.Name()})
		filePath := filepath.Join(dirPath, entry.Name())

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
			hclName := resourceMatches[2]

			nameValue := hclName
			nameMatches := nameRegex.FindStringSubmatch(block)
			if len(nameMatches) > 1 {
				nameValue = nameMatches[1]
			}

			importBlock := fmt.Sprintf("import {\n  to = %s.%s\n  id = \"%s\"\n}\n\n",
				resourceType, hclName, nameValue)

			if _, exists := importBlocks[resourceType]; exists {
				importBlocks[resourceType] += importBlock
			}
		}
	}

	for _, resourceType := range resourceOrder {
		blocks := importBlocks[resourceType]
		if blocks != "" {
			if _, err := file.WriteString(fmt.Sprintf("# %s imports\n%s", resourceType, blocks)); err != nil {
				return "", fmt.Errorf("error writing to output file: %w", err)
			}
			delete(importBlocks, resourceType)
		}
	}

	for resourceType, blocks := range importBlocks {
		if blocks != "" {
			if _, err := file.WriteString(fmt.Sprintf("# %s imports\n%s", resourceType, blocks)); err != nil {
				return "", fmt.Errorf("error writing to output file: %w", err)
			}
		}
	}

	return outputFile, nil
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
