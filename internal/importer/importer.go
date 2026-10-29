package importer

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"terraform-provider-verity/internal/registry"
	"terraform-provider-verity/internal/spec"
	"terraform-provider-verity/internal/transport"
	"terraform-provider-verity/internal/utils"
	"terraform-provider-verity/openapi"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type Importer struct {
	client *openapi.APIClient
	ctx    context.Context
	Mode   string

	supported   map[string]*SchemaFields
	unsupported map[string]map[string]bool
}

type ResourceConfig struct {
	ResourceType     string
	StageName        string
	Fields           []spec.FieldSpec
	SkipTopLevelKeys map[string]bool
}

type ImportedResource struct {
	TerraformType string
	TerraformName string
	ID            string
}

type ImportResult struct {
	Files     []string
	Resources []ImportedResource
}

var rootIndexSkipped = map[string]bool{
	"verity_gateway_profile":  true,
	"verity_eth_port_profile": true,
	"verity_bundle":           true,
}

func (i *Importer) resourceConfig(terraformType string) (ResourceConfig, error) {
	resource, err := registry.Lookup(terraformType)
	if err != nil {
		return ResourceConfig{}, err
	}
	stageName, err := stageNameFor(i.Mode, terraformType)
	if err != nil {
		return ResourceConfig{}, err
	}
	config := ResourceConfig{
		ResourceType:     strings.TrimPrefix(terraformType, "verity_"),
		StageName:        stageName,
		Fields:           resource.Fields,
		SkipTopLevelKeys: map[string]bool{"name": true},
	}
	if rootIndexSkipped[terraformType] {
		config.SkipTopLevelKeys["index"] = true
	}
	return config, nil
}

func stageNameFor(mode, terraformType string) (string, error) {
	stages, err := stageOrder(mode)
	if err != nil {
		return "", err
	}
	for _, stage := range stages {
		if stage.ResourceType == terraformType {
			return stage.StageName, nil
		}
	}
	return "", fmt.Errorf("%s has no import stage for mode %q", terraformType, mode)
}

var nameSplitRE = regexp.MustCompile(`(\d+|\D+)`)

func getNaturalSortParts(s string) []interface{} {
	matches := nameSplitRE.FindAllString(s, -1)
	parts := make([]interface{}, len(matches))
	for i, match := range matches {
		if num, err := strconv.Atoi(match); err == nil {
			parts[i] = num
		} else {
			parts[i] = match
		}
	}
	return parts
}

func NewImporter(client *openapi.APIClient, mode string) *Importer {
	return &Importer{
		client: client,
		ctx:    context.Background(),
		Mode:   mode,
	}
}

func (i *Importer) ImportAll(outputDir string) (ImportResult, error) {
	var result ImportResult
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return ImportResult{}, fmt.Errorf("failed to create output directory: %w", err)
	}

	tflog.Info(i.ctx, "Starting importer with mode", map[string]interface{}{
		"mode":        i.Mode,
		"api_version": utils.GetSupportedAPIVersionString(),
	})

	tasks, err := i.importTasks()
	if err != nil {
		return ImportResult{}, err
	}

	for _, task := range tasks {
		tflog.Info(i.ctx, "Importing resource", map[string]interface{}{
			"resource_name":           task.name,
			"terraform_resource_type": task.terraformResourceType,
		})

		data, err := i.fetchResource(task.resource)
		if err != nil {
			tflog.Error(i.ctx, "Failed to import resource", map[string]interface{}{"resource_name": task.name, "error": err})
			return ImportResult{}, fmt.Errorf("failed to import %s: %w", task.name, err)
		}
		if len(data) == 0 {
			tflog.Info(i.ctx, "No data found for resource, skipping TF generation", map[string]interface{}{"resource_name": task.name})
			continue
		}

		config, err := i.resourceConfig(task.terraformResourceType)
		if err != nil {
			tflog.Error(i.ctx, "No registry entry for terraform type", map[string]interface{}{
				"resource_name":  task.name,
				"terraform_type": task.terraformResourceType,
				"error":          err,
			})
			return ImportResult{}, fmt.Errorf("no registry entry for %s: %w", task.terraformResourceType, err)
		}

		i.PruneUnsupported(task.terraformResourceType, data)

		tfConfig, resources, err := i.generateResourceTF(data, config)
		if err != nil {
			tflog.Error(i.ctx, "Failed to generate Terraform config", map[string]interface{}{"resource_name": task.name, "error": err})
			return ImportResult{}, fmt.Errorf("failed to generate terraform config for %s: %w", task.name, err)
		}

		if strings.TrimSpace(tfConfig) == "" {
			tflog.Info(i.ctx, "Generated TF config is empty, skipping file write", map[string]interface{}{"resource_name": task.name})
			continue
		}

		outputFile := filepath.Join(outputDir, fmt.Sprintf("%s.tf", task.name))
		if err := os.WriteFile(outputFile, []byte(tfConfig), 0644); err != nil {
			tflog.Error(i.ctx, "Failed to write TF config to file", map[string]interface{}{"resource_name": task.name, "file": outputFile, "error": err})
			return ImportResult{}, fmt.Errorf("failed to write %s terraform config: %w", task.name, err)
		}
		result.Files = append(result.Files, outputFile)
		result.Resources = append(result.Resources, resources...)
		tflog.Info(i.ctx, "Successfully wrote TF config for resource", map[string]interface{}{"resource_name": task.name, "file": outputFile})
	}
	if err := i.checkStaleResourceFiles(outputDir, result.Files); err != nil {
		return ImportResult{}, err
	}

	stagesTF, stages, err := i.generateStagesTF()
	if err != nil {
		tflog.Error(i.ctx, "Failed to generate stages TF", map[string]interface{}{"error": err})
		return ImportResult{}, fmt.Errorf("failed to generate stages: %w", err)
	}

	stagesFile := filepath.Join(outputDir, "stages.tf")
	if err := os.WriteFile(stagesFile, []byte(stagesTF), 0644); err != nil {
		tflog.Error(i.ctx, "Failed to write stages TF config", map[string]interface{}{"error": err, "file": stagesFile})
		return ImportResult{}, fmt.Errorf("failed to write stages terraform config: %w", err)
	}
	result.Files = append(result.Files, stagesFile)
	result.Resources = append(result.Resources, stages...)

	return result, nil
}

type importTask struct {
	name                  string
	terraformResourceType string
	resource              spec.ResourceSpec
}

var importFileNames = map[string]string{
	"verity_acl_v4": "acls_ipv4",
	"verity_acl_v6": "acls_ipv6",
}

func importFileName(resource spec.ResourceSpec) string {
	if name, named := importFileNames[resource.TerraformType]; named {
		return name
	}
	return strings.Trim(resource.API.EndpointPath, "/")
}

func (i *Importer) checkStaleResourceFiles(outputDir string, writtenFiles []string) error {
	written := make(map[string]bool, len(writtenFiles))
	for _, file := range writtenFiles {
		written[file] = true
	}
	resources, err := registry.Load()
	if err != nil {
		return fmt.Errorf("check stale resource files: %w", err)
	}
	for _, resource := range resources {
		file := filepath.Join(outputDir, importFileName(resource)+".tf")
		if written[file] {
			continue
		}
		if _, err := os.Lstat(file); err == nil {
			return fmt.Errorf("resource file %q was not regenerated: no %s objects were fetched in mode %q; Terraform could recreate deleted objects from this file. Review and move or remove it before rerunning the importer", file, resource.TerraformType, i.Mode)
		} else if !os.IsNotExist(err) {
			return fmt.Errorf("check resource file %q: %w", file, err)
		}
	}
	return nil
}

func (i *Importer) importTasks() ([]importTask, error) {
	tasks := make([]importTask, 0, 50)
	resourceTypes, err := ResourceTypeOrder(i.Mode)
	if err != nil {
		return nil, err
	}
	for _, terraformType := range resourceTypes {
		resource, err := registry.Lookup(terraformType)
		if err != nil {
			return nil, fmt.Errorf("no registry entry for %s: %w", terraformType, err)
		}
		tasks = append(tasks, importTask{name: importFileName(resource), terraformResourceType: terraformType, resource: resource})
	}
	return tasks, nil
}

func (i *Importer) fetchResource(resource spec.ResourceSpec) (map[string]map[string]interface{}, error) {
	return i.fetch(resource.TerraformType, resource.API.EndpointPath, resource.API.FixedHeaders, resource.API.ResponseCollectionKey)
}

func (i *Importer) fetch(label, endpointPath string, fixedHeaders map[string]string, collectionKey string) (map[string]map[string]interface{}, error) {
	collection, err := transport.FetchCollection(i.ctx, i.client, label, endpointPath, fixedHeaders, collectionKey)
	if err != nil {
		return nil, err
	}
	objects := make(map[string]map[string]interface{}, len(collection))
	for name, raw := range collection {
		object, ok := raw.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("%s object %q is not an object", label, name)
		}
		objects[name] = object
	}
	return objects, nil
}

type stageDefinition struct {
	StageName      string
	ResourceType   string
	DependsOnStage string
	order          int
}

func stageOrder(mode string) ([]stageDefinition, error) {
	resources, err := registry.Load()
	if err != nil {
		return nil, fmt.Errorf("load import stages: %w", err)
	}
	stages := make([]stageDefinition, 0, len(resources))
	for _, resource := range resources {
		stage, declared := resource.ImportStages[spec.Mode(mode)]
		if !declared {
			continue
		}
		stages = append(stages, stageDefinition{StageName: stage.Name, ResourceType: resource.TerraformType, order: stage.Order})
	}
	sort.Slice(stages, func(i, j int) bool { return stages[i].order < stages[j].order })
	for index := range stages {
		if index > 0 {
			stages[index].DependsOnStage = stages[index-1].StageName
		}
	}
	return stages, nil
}

func ResourceTypeOrder(mode string) ([]string, error) {
	stages, err := stageOrder(mode)
	if err != nil {
		return nil, err
	}
	order := make([]string, 0, len(stages))
	for _, stage := range stages {
		if utils.IsResourceCompatibleWithMode(stage.ResourceType, mode) {
			order = append(order, stage.ResourceType)
		}
	}
	return order, nil
}

func (i *Importer) generateStagesTF() (string, []ImportedResource, error) {
	var tfConfig strings.Builder

	tflog.Info(i.ctx, "Generating stages for mode", map[string]interface{}{
		"mode": i.Mode,
	})

	stages, err := stageOrder(i.Mode)
	if err != nil {
		return "", nil, err
	}

	var compatibleStages []stageDefinition
	var lastCompatibleStage string

	for _, stage := range stages {
		if utils.IsResourceCompatibleWithMode(stage.ResourceType, i.Mode) {
			if stage.DependsOnStage != "" && lastCompatibleStage != "" && stage.DependsOnStage != lastCompatibleStage {
				stage.DependsOnStage = lastCompatibleStage
			}
			compatibleStages = append(compatibleStages, stage)
			lastCompatibleStage = stage.StageName
		} else {
			tflog.Debug(i.ctx, "Excluding stage for incompatible resource", map[string]interface{}{
				"stage_name":    stage.StageName,
				"resource_type": stage.ResourceType,
				"mode":          i.Mode,
			})
		}
	}

	modeComment := strings.ToUpper(i.Mode)
	tfConfig.WriteString(fmt.Sprintf("\n# These resources establish ordering for bulk operations in %s mode\n", modeComment))

	var resources []ImportedResource
	for _, stage := range compatibleStages {
		resources = append(resources, ImportedResource{TerraformType: "verity_operation_stage", TerraformName: stage.StageName, ID: "stage"})
		tfConfig.WriteString(fmt.Sprintf("resource \"verity_operation_stage\" \"%s\" {\n", stage.StageName))

		if stage.DependsOnStage != "" {
			tfConfig.WriteString(fmt.Sprintf("  depends_on = [verity_operation_stage.%s]\n", stage.DependsOnStage))
		}

		tfConfig.WriteString("  lifecycle {\n")
		tfConfig.WriteString("    create_before_destroy = true\n")
		tfConfig.WriteString("  }\n")
		tfConfig.WriteString("}\n\n")
	}

	tflog.Info(i.ctx, "Generated stages", map[string]interface{}{
		"mode":              i.Mode,
		"total_stages":      len(stages),
		"compatible_stages": len(compatibleStages),
	})

	return tfConfig.String(), resources, nil
}
