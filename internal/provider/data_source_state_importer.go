package provider

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"terraform-provider-verity/internal/importer"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"

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
	imp := importer.NewImporter(client, d.client.mode).WithSupportedFields(importerSupportedFields(ctx))
	result, err := imp.ImportAll(absPath)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Resources",
			fmt.Sprintf("Error importing resources: %v", err),
		)
		return
	}
	warning := unsupportedFieldsWarning(imp.UnsupportedFields())
	if warning != "" {
		tflog.Warn(ctx, warning)
		resp.Diagnostics.AddWarning("Some arguments are not supported by this provider version", warning)
	}
	if err := writeUnsupportedArgumentsFile(absPath, warning); err != nil {
		resp.Diagnostics.AddWarning("Error Writing "+unsupportedArgumentsFile, err.Error())
	}

	data.ImportedFiles = []types.String{}
	sort.Strings(result.Files)
	for _, filePath := range result.Files {
		data.ImportedFiles = append(data.ImportedFiles, types.StringValue(filePath))
	}

	importBlocksFile, err := createImportBlocks(absPath, result.Resources)
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

func createImportBlocks(dirPath string, resources []importer.ImportedResource) (string, error) {
	var output strings.Builder
	output.WriteString("# Import blocks for Verity resources\n\n")
	previousType := ""
	for _, resource := range resources {
		if resource.TerraformType != previousType {
			fmt.Fprintf(&output, "# %s imports\n", resource.TerraformType)
			previousType = resource.TerraformType
		}
		id := hclwrite.TokensForValue(cty.StringVal(resource.ID)).Bytes()
		fmt.Fprintf(&output, "import {\n  to = %s.%s\n  id = %s\n}\n\n",
			resource.TerraformType, resource.TerraformName, id)
	}

	outputFile := filepath.Join(dirPath, "import_blocks.tf")
	if err := os.WriteFile(outputFile, []byte(output.String()), 0644); err != nil {
		return "", fmt.Errorf("error writing import blocks: %w", err)
	}
	return outputFile, nil
}
