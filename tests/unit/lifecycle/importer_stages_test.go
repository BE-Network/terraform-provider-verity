package lifecycle

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
	fwresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"

	"terraform-provider-verity/internal/bulkops"
	"terraform-provider-verity/tests/unit/mock"
)

func TestImporterImportsStagesWithoutCreatingThem(t *testing.T) {
	for _, mode := range []string{"datacenter", "campus"} {
		t.Run(mode, func(t *testing.T) {
			ms := newImporterServer(t, mode)
			outputDir := t.TempDir()
			applyImporter(t, ms, mode, outputDir)
			stages, err := os.ReadFile(filepath.Join(outputDir, "stages.tf"))
			if err != nil {
				t.Fatal(err)
			}
			content, err := os.ReadFile(filepath.Join(outputDir, "import_blocks.tf"))
			if err != nil {
				t.Fatal(err)
			}
			file, diags := hclwrite.ParseConfig(content, "import_blocks.tf", hcl.InitialPos)
			if diags.HasErrors() {
				t.Fatal(diags)
			}
			imports := hclwrite.NewEmptyFile()
			var checks []fwresource.TestCheckFunc
			var planChecks []plancheck.PlanCheck
			for _, block := range file.Body().Blocks() {
				expr, diags := hclsyntax.ParseExpression(block.Body().GetAttribute("to").Expr().BuildTokens(nil).Bytes(), "import_blocks.tf", hcl.InitialPos)
				if diags.HasErrors() {
					t.Fatal(diags)
				}
				to, diags := hcl.AbsTraversalForExpr(expr)
				if diags.HasErrors() || len(to) != 2 {
					t.Fatalf("invalid import target: %v, %v", to, diags)
				}
				if to.RootName() != "verity_operation_stage" {
					continue
				}
				address := to.RootName() + "." + to[1].(hcl.TraverseAttr).Name
				imports.Body().AppendBlock(block)
				checks = append(checks, fwresource.TestCheckResourceAttr(address, "id", "stage"))
				planChecks = append(planChecks, plancheck.ExpectResourceAction(address, plancheck.ResourceActionNoop))
			}
			if len(checks) == 0 {
				t.Fatal("no stage import blocks generated")
			}
			config := mock.ProviderConfig(ms.URL(), mode) + string(stages) + string(imports.Bytes())
			oldDelay := bulkops.MaxBatchDelay
			t.Cleanup(func() { bulkops.MaxBatchDelay = oldDelay })
			fwresource.UnitTest(t, fwresource.TestCase{
				ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
				Steps: []fwresource.TestStep{
					{
						Config:           config,
						ConfigPlanChecks: fwresource.ConfigPlanChecks{PreApply: planChecks},
						Check:            fwresource.ComposeAggregateTestCheckFunc(checks...),
					},
					{
						Config:           config,
						ConfigPlanChecks: fwresource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
						Check:            fwresource.ComposeAggregateTestCheckFunc(checks...),
					},
					{
						Config:    config,
						Destroy:   true,
						PreConfig: func() { bulkops.MaxBatchDelay = 0 },
					},
					{
						Config: mock.ProviderConfig(ms.URL(), mode) + string(stages),
						Check:  fwresource.ComposeAggregateTestCheckFunc(checks...),
					},
					{
						Config:           config,
						PreConfig:        func() { bulkops.MaxBatchDelay = oldDelay },
						ConfigPlanChecks: fwresource.ConfigPlanChecks{PreApply: []plancheck.PlanCheck{plancheck.ExpectEmptyPlan()}},
						Check:            fwresource.ComposeAggregateTestCheckFunc(checks...),
					},
					{
						Config:    config,
						Destroy:   true,
						PreConfig: func() { bulkops.MaxBatchDelay = 0 },
					},
				},
			})
		})
	}
}
