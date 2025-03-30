package provider

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"terraform-provider-verity/internal/importer"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
				Description: "Directory where the TF files will be saved",
				Required:    true,
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
	imp := importer.NewImporter(client)
	err = imp.ImportAll(absPath)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Importing Resources",
			fmt.Sprintf("Error importing resources: %v", err),
		)
		return
	}

	data.ID = types.StringValue(time.Now().UTC().String())
	data.ImportedFiles = []types.String{}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
