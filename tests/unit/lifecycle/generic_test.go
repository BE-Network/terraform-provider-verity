package lifecycle

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"

	"terraform-provider-verity/internal/utils"
	"terraform-provider-verity/tests/unit/mock"
)

// deleteParamOverrides maps WrapperKey to the DELETE query parameter name
// for resources where it doesn't follow the convention of wrapperKey+"_name".
var deleteParamOverrides = map[string]string{
	"endpoint_bundle":     "bundle_name",
	"eth_port_profile_":   "profile_name",
	"eth_port_settings":   "port_name",
	"gateway_profile":     "profile_name",
	"voice_port_profiles": "voice_port_profile_name",
}

func detectRefPairs(rs resourceSchemaInfo) [][2]string {
	refTypeNames := make(map[string]bool)
	for _, fi := range rs.Attributes {
		if strings.HasSuffix(fi.Name, "_ref_type_") {
			refTypeNames[fi.Name] = true
		}
	}
	var pairs [][2]string
	for refType := range refTypeNames {
		base := strings.TrimSuffix(refType, "_ref_type_")
		for _, fi := range rs.Attributes {
			if fi.Name == base {
				pairs = append(pairs, [2]string{base, refType})
				break
			}
		}
	}
	return pairs
}

func detectNullableFields(rs resourceSchemaInfo) []fieldInfo {
	autoAssignedValues := make(map[string]bool)
	for _, fi := range rs.Attributes {
		if strings.HasSuffix(fi.Name, "_auto_assigned_") {
			autoAssignedValues[strings.TrimSuffix(fi.Name, "_auto_assigned_")] = true
		}
	}
	var result []fieldInfo
	for _, fi := range rs.Attributes {
		if fi.Type != "int64" && fi.Type != "number" {
			continue
		}
		if fi.Name == "name" || fi.Name == "index" {
			continue
		}
		if autoAssignedValues[fi.Name] {
			continue
		}
		result = append(result, fi)
	}
	return result
}

func detectAutoAssignedPairs(rs resourceSchemaInfo) map[string]string {
	result := make(map[string]string)
	for _, fi := range rs.Attributes {
		if fi.Type == "bool" && strings.HasSuffix(fi.Name, "_auto_assigned_") {
			valueName := strings.TrimSuffix(fi.Name, "_auto_assigned_")
			for _, vf := range rs.Attributes {
				if vf.Name == valueName {
					result[fi.Name] = valueName
					break
				}
			}
		}
	}
	return result
}

func generateNameOnlyHCL(tfType, resourceName string) string {
	return fmt.Sprintf("resource %q %q {\n  name = %q\n}\n", tfType, "test", resourceName)
}

func generateHCLWithExcludes(rs resourceSchemaInfo, tfType, resourceName, mode, modeFieldsKey string,
	overrides map[string]string, excludes map[string]bool) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("resource %q %q {\n", tfType, "test"))

	for _, fi := range rs.Attributes {
		if excludes != nil && excludes[fi.Name] {
			continue
		}
		if !utils.FieldAppliesToMode(modeFieldsKey, fi.Name, mode) {
			continue
		}
		val := defaultHCLValue(fi)
		if fi.Name == "name" {
			val = fmt.Sprintf("%q", resourceName)
		}
		if override, ok := overrides[fi.Name]; ok {
			val = override
		}
		b.WriteString(fmt.Sprintf("  %s = %s\n", fi.Name, val))
	}

	for _, block := range rs.Blocks {
		if excludes != nil && excludes[block.Name] {
			continue
		}
		if !utils.FieldAppliesToMode(modeFieldsKey, block.Name, mode) {
			continue
		}
		if len(block.Fields) == 0 {
			continue
		}
		b.WriteString(fmt.Sprintf("\n  %s {\n", block.Name))
		for _, fi := range block.Fields {
			nestedKey := block.Name + "." + fi.Name
			if excludes != nil && excludes[nestedKey] {
				continue
			}
			if !utils.FieldAppliesToMode(modeFieldsKey, nestedKey, mode) {
				continue
			}
			val := defaultHCLValue(fi)
			if override, ok := overrides[nestedKey]; ok {
				val = override
			}
			b.WriteString(fmt.Sprintf("    %s = %s\n", fi.Name, val))
		}
		b.WriteString("  }\n")
	}

	b.WriteString("}\n")
	return b.String()
}

func mergeOverrides(base, extra map[string]string) map[string]string {
	result := make(map[string]string)
	for k, v := range base {
		result[k] = v
	}
	for k, v := range extra {
		result[k] = v
	}
	return result
}

func generateHCLWithBlockCount(rs resourceSchemaInfo, tfType, resourceName, mode, modeFieldsKey string,
	overrides map[string]string, indexedBlockCount int) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("resource %q %q {\n", tfType, "test"))

	for _, fi := range rs.Attributes {
		if !utils.FieldAppliesToMode(modeFieldsKey, fi.Name, mode) {
			continue
		}
		val := defaultHCLValue(fi)
		if fi.Name == "name" {
			val = fmt.Sprintf("%q", resourceName)
		}
		if override, ok := overrides[fi.Name]; ok {
			val = override
		}
		b.WriteString(fmt.Sprintf("  %s = %s\n", fi.Name, val))
	}

	for _, block := range rs.Blocks {
		if !utils.FieldAppliesToMode(modeFieldsKey, block.Name, mode) {
			continue
		}
		if len(block.Fields) == 0 {
			continue
		}

		hasIndex := false
		for _, fi := range block.Fields {
			if fi.Name == "index" {
				hasIndex = true
				break
			}
		}

		count := 1
		if hasIndex && indexedBlockCount > 1 {
			count = indexedBlockCount
		}

		for i := 1; i <= count; i++ {
			b.WriteString(fmt.Sprintf("\n  %s {\n", block.Name))
			for _, fi := range block.Fields {
				nestedKey := block.Name + "." + fi.Name
				if !utils.FieldAppliesToMode(modeFieldsKey, nestedKey, mode) {
					continue
				}
				val := defaultHCLValue(fi)
				if fi.Name == "index" {
					val = fmt.Sprintf("%d", i)
				}
				if override, ok := overrides[nestedKey]; ok {
					val = override
				}
				b.WriteString(fmt.Sprintf("    %s = %s\n", fi.Name, val))
			}
			b.WriteString("  }\n")
		}
	}

	b.WriteString("}\n")
	return b.String()
}

func collectExcludedModeFields(rs resourceSchemaInfo, modeFieldsKey, mode string) []string {
	var excluded []string
	for _, fi := range rs.Attributes {
		if !utils.FieldAppliesToMode(modeFieldsKey, fi.Name, mode) {
			excluded = append(excluded, fi.Name)
		}
	}
	for _, block := range rs.Blocks {
		if !utils.FieldAppliesToMode(modeFieldsKey, block.Name, mode) {
			excluded = append(excluded, block.Name)
			continue
		}
		for _, fi := range block.Fields {
			nestedKey := block.Name + "." + fi.Name
			if !utils.FieldAppliesToMode(modeFieldsKey, nestedKey, mode) {
				excluded = append(excluded, nestedKey)
			}
		}
	}
	return excluded
}

func deleteParamName(wrapperKey string) string {
	if override, ok := deleteParamOverrides[wrapperKey]; ok {
		return override
	}
	return wrapperKey + "_name"
}

func TestGeneric_PutOnlyName(t *testing.T) {
	t.Parallel()
	for _, tc := range allResourceTests {
		t.Run(tc.TerraformType, func(t *testing.T) {
			t.Parallel()
			if tc.SkipCreate {
				t.Skipf("%s is update-only, skipping", tc.TerraformType)
			}

			ms := mock.NewMockServer(tc.Mode)
			defer ms.Close()
			if err := ms.LoadResponsesFromDir(mock.ResponsesDir(tc.Mode)); err != nil {
				t.Fatalf("failed to load responses: %v", err)
			}

			rs := inspectSchema(tc.Factory)
			autoAssigned := detectAutoAssignedPairs(rs)

			resourceName := "nom_" + tc.ResourceName
			hcl := generateNameOnlyHCL(tc.TerraformType, resourceName)
			config := mock.ProviderConfig(ms.URL(), tc.Mode) + hcl

			allowedOptional := make(map[string]bool)
			for flag, value := range autoAssigned {
				allowedOptional[flag] = true
				allowedOptional[value] = true
			}

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
				Steps: []resource.TestStep{
					{
						Config: config,
						Check: func(s *terraform.State) error {
							puts := ms.GetRequestsByMethodAndPath("PUT", tc.APIPath)
							if len(puts) == 0 {
								return fmt.Errorf("no PUT request captured for %s", tc.APIPath)
							}
							body := puts[len(puts)-1].Body
							path := tc.WrapperKey + "." + resourceName
							mock.AssertFieldEquals(t, body, path+".name", resourceName)
							// Custom check: "name" required, auto-assigned fields allowed but optional, anything else is unexpected
							val, found := body[tc.WrapperKey]
							if !found {
								return fmt.Errorf("wrapper %q not found", tc.WrapperKey)
							}
							wrapper := val.(map[string]interface{})
							obj := wrapper[resourceName].(map[string]interface{})
							for key := range obj {
								if key == "name" || allowedOptional[key] {
									continue
								}
								t.Errorf("unexpected field %q in PUT for name-only create", key)
							}
							return nil
						},
					},
				},
			})
		})
	}
}

func TestGeneric_PatchEnableField(t *testing.T) {
	t.Parallel()
	for _, tc := range allResourceTests {
		t.Run(tc.TerraformType, func(t *testing.T) {
			t.Parallel()
			if tc.SkipCreate {
				t.Skipf("%s is update-only, skipping", tc.TerraformType)
			}

			rs := inspectSchema(tc.Factory)

			ms := mock.NewMockServer(tc.Mode)
			defer ms.Close()
			if err := ms.LoadResponsesFromDir(mock.ResponsesDir(tc.Mode)); err != nil {
				t.Fatalf("failed to load responses: %v", err)
			}

			resourceName := "pat_" + tc.ResourceName
			modeKey := tc.modeFieldsKey()

			createOverrides := mergeOverrides(tc.Overrides, map[string]string{"enable": "true"})
			updateOverrides := mergeOverrides(tc.Overrides, map[string]string{"enable": "false"})

			createHCL := generateHCLWithExcludes(rs, tc.TerraformType, resourceName, tc.Mode, modeKey, createOverrides, nil)
			updateHCL := generateHCLWithExcludes(rs, tc.TerraformType, resourceName, tc.Mode, modeKey, updateOverrides, nil)

			createConfig := mock.ProviderConfig(ms.URL(), tc.Mode) + createHCL
			updateConfig := mock.ProviderConfig(ms.URL(), tc.Mode) + updateHCL

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
				Steps: []resource.TestStep{
					{
						PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), createConfig) },
						Config:    createConfig,
						Check: func(s *terraform.State) error {
							ms.Reset()
							return nil
						},
					},
					{
						PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), updateConfig) },
						Config:    updateConfig,
						Check: func(s *terraform.State) error {
							patches := ms.GetRequestsByMethodAndPath("PATCH", tc.APIPath)
							if len(patches) == 0 {
								return fmt.Errorf("no PATCH request captured for %s", tc.APIPath)
							}
							body := patches[len(patches)-1].Body
							path := tc.WrapperKey + "." + resourceName
							mock.AssertFieldEquals(t, body, path+".enable", false)
							mock.AssertOnlyFields(t, body, path, []string{"enable"})
							return nil
						},
					},
				},
			})
		})
	}
}

func TestGeneric_RequiredQueryParams(t *testing.T) {
	t.Parallel()
	for _, tc := range allResourceTests {
		if len(tc.RequiredQueryParams) == 0 {
			continue
		}
		t.Run(tc.TerraformType, func(t *testing.T) {
			t.Parallel()

			ms := mock.NewMockServer(tc.Mode)
			defer ms.Close()
			if err := ms.LoadResponsesFromDir(mock.ResponsesDir(tc.Mode)); err != nil {
				t.Fatalf("failed to load responses: %v", err)
			}

			rs := inspectSchema(tc.Factory)
			modeKey := tc.modeFieldsKey()
			resourceName := "qp_" + tc.ResourceName
			basePath := tc.WrapperKey + "." + resourceName

			createOverrides := mergeOverrides(tc.Overrides, map[string]string{"enable": "true"})
			updateOverrides := mergeOverrides(tc.Overrides, map[string]string{"enable": "false"})

			createHCL := generateHCLWithExcludes(rs, tc.TerraformType, resourceName, tc.Mode, modeKey, createOverrides, nil)
			updateHCL := generateHCLWithExcludes(rs, tc.TerraformType, resourceName, tc.Mode, modeKey, updateOverrides, nil)

			createConfig := mock.ProviderConfig(ms.URL(), tc.Mode) + createHCL
			updateConfig := mock.ProviderConfig(ms.URL(), tc.Mode) + updateHCL

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
				Steps: []resource.TestStep{
					{
						PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), createConfig) },
						Config:    createConfig,
						Check: func(s *terraform.State) error {
							ms.Reset()
							return nil
						},
					},
					{
						PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), updateConfig) },
						Config:    updateConfig,
						Check: func(s *terraform.State) error {
							patches := ms.GetRequestsByMethodAndPath("PATCH", tc.APIPath)
							if len(patches) == 0 {
								return fmt.Errorf("no PATCH request captured for %s", tc.APIPath)
							}
							patch := patches[len(patches)-1]

							mock.AssertFieldEquals(t, patch.Body, basePath+".enable", false)
							mock.AssertOnlyFields(t, patch.Body, basePath, []string{"enable"})

							for param, expected := range tc.RequiredQueryParams {
								values := patch.QueryParams[param]
								if len(values) == 0 || values[0] != expected {
									t.Errorf("PATCH query param %q = %v, want [%s]", param, values, expected)
								}
							}
							return nil
						},
					},
				},
			})
		})
	}
}

func TestGeneric_RefFieldPairsInPut(t *testing.T) {
	t.Parallel()
	for _, tc := range allResourceTests {
		t.Run(tc.TerraformType, func(t *testing.T) {
			t.Parallel()
			if tc.SkipCreate {
				t.Skipf("%s is update-only, skipping", tc.TerraformType)
			}

			rs := inspectSchema(tc.Factory)
			modeKey := tc.modeFieldsKey()
			refPairs := detectRefPairs(rs)

			var applicablePairs [][2]string
			for _, pair := range refPairs {
				if utils.FieldAppliesToMode(modeKey, pair[0], tc.Mode) {
					applicablePairs = append(applicablePairs, pair)
				}
			}
			if len(applicablePairs) == 0 {
				t.Skipf("%s has no ref_type_ fields, skipping", tc.TerraformType)
			}

			ms := mock.NewMockServer(tc.Mode)
			defer ms.Close()
			if err := ms.LoadResponsesFromDir(mock.ResponsesDir(tc.Mode)); err != nil {
				t.Fatalf("failed to load responses: %v", err)
			}

			resourceName := "ref_" + tc.ResourceName
			overrides := mergeOverrides(tc.Overrides, nil)
			for _, pair := range applicablePairs {
				overrides[pair[0]] = `"TestRefValue"`
				overrides[pair[1]] = `"test_ref_type"`
			}

			hcl := generateHCLWithExcludes(rs, tc.TerraformType, resourceName, tc.Mode, modeKey, overrides, nil)
			config := mock.ProviderConfig(ms.URL(), tc.Mode) + hcl

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
				Steps: []resource.TestStep{
					{
						PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), config) },
						Config:    config,
						Check: func(s *terraform.State) error {
							puts := ms.GetRequestsByMethodAndPath("PUT", tc.APIPath)
							if len(puts) == 0 {
								return fmt.Errorf("no PUT request captured for %s", tc.APIPath)
							}
							body := puts[len(puts)-1].Body
							basePath := tc.WrapperKey + "." + resourceName

							for _, pair := range applicablePairs {
								mock.AssertFieldEquals(t, body, basePath+"."+pair[0], "TestRefValue")
								mock.AssertFieldEquals(t, body, basePath+"."+pair[1], "test_ref_type")
							}
							return nil
						},
					},
				},
			})
		})
	}
}

func TestGeneric_NullableFieldTransitions(t *testing.T) {
	t.Parallel()
	for _, tc := range allResourceTests {
		t.Run(tc.TerraformType, func(t *testing.T) {
			t.Parallel()
			if tc.SkipCreate {
				t.Skipf("%s is update-only, skipping", tc.TerraformType)
			}

			rs := inspectSchema(tc.Factory)
			modeKey := tc.modeFieldsKey()

			nullableFields := detectNullableFields(rs)
			var applicable []fieldInfo
			for _, fi := range nullableFields {
				if utils.FieldAppliesToMode(modeKey, fi.Name, tc.Mode) {
					applicable = append(applicable, fi)
				}
			}
			if len(applicable) == 0 {
				t.Skipf("%s has no nullable fields, skipping", tc.TerraformType)
			}

			target := applicable[0]
			ms := mock.NewMockServer(tc.Mode)
			defer ms.Close()
			if err := ms.LoadResponsesFromDir(mock.ResponsesDir(tc.Mode)); err != nil {
				t.Fatalf("failed to load responses: %v", err)
			}

			resourceName := "nul_" + tc.ResourceName
			basePath := tc.WrapperKey + "." + resourceName

			initVal := "42"
			if target.Type == "number" {
				initVal = "1.5"
			}

			createOverrides := mergeOverrides(tc.Overrides, map[string]string{target.Name: initVal})
			nullOverrides := mergeOverrides(tc.Overrides, map[string]string{target.Name: "null"})
			zeroOverrides := mergeOverrides(tc.Overrides, map[string]string{target.Name: "0"})

			createHCL := generateHCLWithExcludes(rs, tc.TerraformType, resourceName, tc.Mode, modeKey, createOverrides, nil)
			nullHCL := generateHCLWithExcludes(rs, tc.TerraformType, resourceName, tc.Mode, modeKey, nullOverrides, nil)
			zeroHCL := generateHCLWithExcludes(rs, tc.TerraformType, resourceName, tc.Mode, modeKey, zeroOverrides, nil)

			createConfig := mock.ProviderConfig(ms.URL(), tc.Mode) + createHCL
			nullConfig := mock.ProviderConfig(ms.URL(), tc.Mode) + nullHCL
			zeroConfig := mock.ProviderConfig(ms.URL(), tc.Mode) + zeroHCL

			t.Logf("Testing nullable transitions for field %q (%s)", target.Name, target.Type)

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
				Steps: []resource.TestStep{
					// Step 1: Create with initial value
					{
						PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), createConfig) },
						Config:    createConfig,
						Check: func(s *terraform.State) error {
							ms.Reset()
							return nil
						},
					},
					// Step 2: Change to null
					{
						PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), nullConfig) },
						Config:    nullConfig,
						Check: func(s *terraform.State) error {
							patches := ms.GetRequestsByMethodAndPath("PATCH", tc.APIPath)
							if len(patches) == 0 {
								return fmt.Errorf("no PATCH for %s after %s→null", tc.APIPath, target.Name)
							}
							body := patches[len(patches)-1].Body
							mock.AssertFieldNull(t, body, basePath+"."+target.Name)
							ms.Reset()
							return nil
						},
					},
					// Step 3: Change from null to zero
					{
						PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), zeroConfig) },
						Config:    zeroConfig,
						Check: func(s *terraform.State) error {
							patches := ms.GetRequestsByMethodAndPath("PATCH", tc.APIPath)
							if len(patches) == 0 {
								return fmt.Errorf("no PATCH for %s after null→0", tc.APIPath)
							}
							body := patches[len(patches)-1].Body
							mock.AssertFieldEquals(t, body, basePath+"."+target.Name, 0)
							return nil
						},
					},
				},
			})
		})
	}
}

func TestGeneric_AutoAssignedFieldsExcluded(t *testing.T) {
	t.Parallel()
	for _, tc := range allResourceTests {
		t.Run(tc.TerraformType, func(t *testing.T) {
			t.Parallel()
			if tc.SkipCreate {
				t.Skipf("%s is update-only, skipping", tc.TerraformType)
			}

			rs := inspectSchema(tc.Factory)
			autoAssigned := detectAutoAssignedPairs(rs)
			if len(autoAssigned) == 0 {
				t.Skipf("%s has no auto-assigned fields, skipping", tc.TerraformType)
			}

			ms := mock.NewMockServer(tc.Mode)
			defer ms.Close()
			if err := ms.LoadResponsesFromDir(mock.ResponsesDir(tc.Mode)); err != nil {
				t.Fatalf("failed to load responses: %v", err)
			}

			resourceName := "aut_" + tc.ResourceName
			modeKey := tc.modeFieldsKey()
			overrides := make(map[string]string)
			excludes := make(map[string]bool)
			for k, v := range tc.Overrides {
				overrides[k] = v
			}
			for flag, valueField := range autoAssigned {
				overrides[flag] = "true"
				excludes[valueField] = true
			}

			enrichment := make(map[string]interface{})
			for _, valueField := range autoAssigned {
				for _, fi := range rs.Attributes {
					if fi.Name == valueField {
						switch fi.Type {
						case "int64":
							enrichment[valueField] = float64(99999)
						case "number":
							enrichment[valueField] = float64(99.99)
						case "string":
							enrichment[valueField] = "auto_generated"
						}
						break
					}
				}
			}
			ms.SetPostPutEnrichment(tc.APIPath, tc.WrapperKey, resourceName, enrichment)

			hcl := generateHCLWithExcludes(rs, tc.TerraformType, resourceName, tc.Mode, modeKey, overrides, excludes)
			config := mock.ProviderConfig(ms.URL(), tc.Mode) + hcl

			t.Logf("Auto-assigned pairs: %v", autoAssigned)

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
				Steps: []resource.TestStep{
					{
						PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), config) },
						Config:    config,
						Check: func(s *terraform.State) error {
							puts := ms.GetRequestsByMethodAndPath("PUT", tc.APIPath)
							if len(puts) == 0 {
								return fmt.Errorf("no PUT request captured for %s", tc.APIPath)
							}
							body := puts[len(puts)-1].Body
							basePath := tc.WrapperKey + "." + resourceName

							for flag := range autoAssigned {
								mock.AssertFieldEquals(t, body, basePath+"."+flag, true)
							}
							for _, valueField := range autoAssigned {
								mock.AssertFieldAbsent(t, body, basePath+"."+valueField)
							}
							return nil
						},
					},
				},
			})
		})
	}
}

func TestGeneric_Delete(t *testing.T) {
	t.Parallel()
	for _, tc := range allResourceTests {
		t.Run(tc.TerraformType, func(t *testing.T) {
			t.Parallel()
			if tc.SkipCreate {
				t.Skipf("%s is update-only (no delete), skipping", tc.TerraformType)
			}

			ms := mock.NewMockServer(tc.Mode)
			defer ms.Close()
			if err := ms.LoadResponsesFromDir(mock.ResponsesDir(tc.Mode)); err != nil {
				t.Fatalf("failed to load responses: %v", err)
			}

			resourceName := "del_" + tc.ResourceName
			hcl := generateNameOnlyHCL(tc.TerraformType, resourceName)
			createConfig := mock.ProviderConfig(ms.URL(), tc.Mode) + hcl
			emptyConfig := mock.ProviderConfig(ms.URL(), tc.Mode)

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
				Steps: []resource.TestStep{
					{
						Config: createConfig,
						Check: func(s *terraform.State) error {
							ms.Reset()
							return nil
						},
					},
					{
						Config: emptyConfig,
						Check: func(s *terraform.State) error {
							deletes := ms.GetRequestsByMethodAndPath("DELETE", tc.APIPath)
							if len(deletes) == 0 {
								return fmt.Errorf("no DELETE requests captured for %s", tc.APIPath)
							}
							del := deletes[len(deletes)-1]
							paramName := deleteParamName(tc.WrapperKey)
							mock.AssertDeleteQueryParams(t, del, paramName, []string{resourceName})
							return nil
						},
					},
				},
			})
		})
	}
}

func TestGeneric_PatchRefFieldPairs(t *testing.T) {
	t.Parallel()
	for _, tc := range allResourceTests {
		t.Run(tc.TerraformType, func(t *testing.T) {
			t.Parallel()
			if tc.SkipCreate {
				t.Skipf("%s is update-only, skipping", tc.TerraformType)
			}

			rs := inspectSchema(tc.Factory)
			modeKey := tc.modeFieldsKey()

			refPairs := detectRefPairs(rs)
			var applicablePairs [][2]string
			for _, pair := range refPairs {
				if utils.FieldAppliesToMode(modeKey, pair[0], tc.Mode) {
					applicablePairs = append(applicablePairs, pair)
				}
			}
			if len(applicablePairs) == 0 {
				t.Skipf("%s has no ref_type_ fields, skipping", tc.TerraformType)
			}

			ms := mock.NewMockServer(tc.Mode)
			defer ms.Close()
			if err := ms.LoadResponsesFromDir(mock.ResponsesDir(tc.Mode)); err != nil {
				t.Fatalf("failed to load responses: %v", err)
			}

			resourceName := "prf_" + tc.ResourceName
			basePath := tc.WrapperKey + "." + resourceName

			// Create with initial ref values
			createOverrides := mergeOverrides(tc.Overrides, nil)
			for _, pair := range applicablePairs {
				createOverrides[pair[0]] = `"OldRefValue"`
				createOverrides[pair[1]] = `"old_ref_type"`
			}

			// Update: change only the first ref pair
			updateOverrides := mergeOverrides(tc.Overrides, nil)
			for i, pair := range applicablePairs {
				if i == 0 {
					updateOverrides[pair[0]] = `"NewRefValue"`
					updateOverrides[pair[1]] = `"new_ref_type"`
				} else {
					updateOverrides[pair[0]] = `"OldRefValue"`
					updateOverrides[pair[1]] = `"old_ref_type"`
				}
			}

			createHCL := generateHCLWithExcludes(rs, tc.TerraformType, resourceName, tc.Mode, modeKey, createOverrides, nil)
			updateHCL := generateHCLWithExcludes(rs, tc.TerraformType, resourceName, tc.Mode, modeKey, updateOverrides, nil)

			createConfig := mock.ProviderConfig(ms.URL(), tc.Mode) + createHCL
			updateConfig := mock.ProviderConfig(ms.URL(), tc.Mode) + updateHCL

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
				Steps: []resource.TestStep{
					{
						PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), createConfig) },
						Config:    createConfig,
						Check: func(s *terraform.State) error {
							ms.Reset()
							return nil
						},
					},
					{
						PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), updateConfig) },
						Config:    updateConfig,
						Check: func(s *terraform.State) error {
							patches := ms.GetRequestsByMethodAndPath("PATCH", tc.APIPath)
							if len(patches) == 0 {
								return fmt.Errorf("no PATCH request for %s after ref field change", tc.APIPath)
							}
							body := patches[len(patches)-1].Body

							// Changed ref pair should be in PATCH
							mock.AssertFieldEquals(t, body, basePath+"."+applicablePairs[0][0], "NewRefValue")
							mock.AssertFieldEquals(t, body, basePath+"."+applicablePairs[0][1], "new_ref_type")

							// PATCH must contain ONLY the changed ref pair
							mock.AssertOnlyFields(t, body, basePath, []string{applicablePairs[0][0], applicablePairs[0][1]})
							return nil
						},
					},
				},
			})
		})
	}
}

func TestGeneric_PatchSingleStringField(t *testing.T) {
	t.Parallel()
	for _, tc := range allResourceTests {
		t.Run(tc.TerraformType, func(t *testing.T) {
			t.Parallel()
			if tc.SkipCreate {
				t.Skipf("%s is update-only, skipping", tc.TerraformType)
			}

			rs := inspectSchema(tc.Factory)
			modeKey := tc.modeFieldsKey()

			refBases := make(map[string]bool)
			for _, pair := range detectRefPairs(rs) {
				refBases[pair[0]] = true
			}

			// Find a non-name, non-ref-type, non-ref-base string field
			var target string
			for _, fi := range rs.Attributes {
				if fi.Type != "string" || fi.Name == "name" {
					continue
				}
				if strings.HasSuffix(fi.Name, "_ref_type_") {
					continue
				}
				if refBases[fi.Name] {
					continue
				}
				if !utils.FieldAppliesToMode(modeKey, fi.Name, tc.Mode) {
					continue
				}
				target = fi.Name
				break
			}
			if target == "" {
				t.Skipf("%s has no patchable string field, skipping", tc.TerraformType)
			}

			ms := mock.NewMockServer(tc.Mode)
			defer ms.Close()
			if err := ms.LoadResponsesFromDir(mock.ResponsesDir(tc.Mode)); err != nil {
				t.Fatalf("failed to load responses: %v", err)
			}

			resourceName := "pss_" + tc.ResourceName
			basePath := tc.WrapperKey + "." + resourceName

			createOverrides := mergeOverrides(tc.Overrides, map[string]string{target: `"initial_value"`})
			updateOverrides := mergeOverrides(tc.Overrides, map[string]string{target: `"updated_value"`})

			createHCL := generateHCLWithExcludes(rs, tc.TerraformType, resourceName, tc.Mode, modeKey, createOverrides, nil)
			updateHCL := generateHCLWithExcludes(rs, tc.TerraformType, resourceName, tc.Mode, modeKey, updateOverrides, nil)

			createConfig := mock.ProviderConfig(ms.URL(), tc.Mode) + createHCL
			updateConfig := mock.ProviderConfig(ms.URL(), tc.Mode) + updateHCL

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
				Steps: []resource.TestStep{
					{
						PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), createConfig) },
						Config:    createConfig,
						Check: func(s *terraform.State) error {
							ms.Reset()
							return nil
						},
					},
					{
						PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), updateConfig) },
						Config:    updateConfig,
						Check: func(s *terraform.State) error {
							patches := ms.GetRequestsByMethodAndPath("PATCH", tc.APIPath)
							if len(patches) == 0 {
								return fmt.Errorf("no PATCH for %s after string field %s change", tc.APIPath, target)
							}
							body := patches[len(patches)-1].Body
							mock.AssertFieldEquals(t, body, basePath+"."+target, "updated_value")
							mock.AssertOnlyFields(t, body, basePath, []string{target})
							return nil
						},
					},
				},
			})
		})
	}
}

func TestGeneric_EdgeCaseZeroValues(t *testing.T) {
	t.Parallel()
	for _, tc := range allResourceTests {
		t.Run(tc.TerraformType, func(t *testing.T) {
			t.Parallel()
			if tc.SkipCreate {
				t.Skipf("%s is update-only, skipping", tc.TerraformType)
			}

			rs := inspectSchema(tc.Factory)
			modeKey := tc.modeFieldsKey()

			var intField, boolField string
			autoAssigned := detectAutoAssignedPairs(rs)
			autoValues := make(map[string]bool)
			for _, v := range autoAssigned {
				autoValues[v] = true
			}

			for _, fi := range rs.Attributes {
				if fi.Name == "name" || fi.Name == "index" {
					continue
				}
				if !utils.FieldAppliesToMode(modeKey, fi.Name, tc.Mode) {
					continue
				}
				if fi.Type == "int64" && intField == "" && !autoValues[fi.Name] {
					intField = fi.Name
				}
				if fi.Type == "bool" && boolField == "" && !strings.HasSuffix(fi.Name, "_auto_assigned_") && fi.Name != "enable" {
					boolField = fi.Name
				}
			}

			if intField == "" && boolField == "" {
				t.Skipf("%s has no testable zero-value fields, skipping", tc.TerraformType)
			}

			ms := mock.NewMockServer(tc.Mode)
			defer ms.Close()
			if err := ms.LoadResponsesFromDir(mock.ResponsesDir(tc.Mode)); err != nil {
				t.Fatalf("failed to load responses: %v", err)
			}

			resourceName := "zer_" + tc.ResourceName

			overrides := mergeOverrides(tc.Overrides, nil)
			if intField != "" {
				overrides[intField] = "0"
			}
			if boolField != "" {
				overrides[boolField] = "false"
			}

			hcl := generateHCLWithExcludes(rs, tc.TerraformType, resourceName, tc.Mode, modeKey, overrides, nil)
			config := mock.ProviderConfig(ms.URL(), tc.Mode) + hcl

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
				Steps: []resource.TestStep{
					{
						PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), config) },
						Config:    config,
						Check: func(s *terraform.State) error {
							puts := ms.GetRequestsByMethodAndPath("PUT", tc.APIPath)
							if len(puts) == 0 {
								return fmt.Errorf("no PUT request captured for %s", tc.APIPath)
							}
							body := puts[len(puts)-1].Body
							basePath := tc.WrapperKey + "." + resourceName

							// Zero values must be PRESENT in PUT, not omitted
							if intField != "" {
								mock.AssertFieldEquals(t, body, basePath+"."+intField, 0)
							}
							if boolField != "" {
								mock.AssertFieldEquals(t, body, basePath+"."+boolField, false)
							}
							return nil
						},
					},
				},
			})
		})
	}
}

func TestGeneric_EmptyStringInPut(t *testing.T) {
	t.Parallel()
	for _, tc := range allResourceTests {
		t.Run(tc.TerraformType, func(t *testing.T) {
			t.Parallel()
			if tc.SkipCreate {
				t.Skipf("%s is update-only, skipping", tc.TerraformType)
			}

			rs := inspectSchema(tc.Factory)
			modeKey := tc.modeFieldsKey()

			refBases := make(map[string]bool)
			for _, pair := range detectRefPairs(rs) {
				refBases[pair[0]] = true
			}

			var target string
			for _, fi := range rs.Attributes {
				if fi.Type != "string" || fi.Name == "name" {
					continue
				}
				if strings.HasSuffix(fi.Name, "_ref_type_") {
					continue
				}
				if refBases[fi.Name] {
					continue
				}
				if !utils.FieldAppliesToMode(modeKey, fi.Name, tc.Mode) {
					continue
				}
				target = fi.Name
				break
			}
			if target == "" {
				t.Skipf("%s has no non-name string field, skipping", tc.TerraformType)
			}

			ms := mock.NewMockServer(tc.Mode)
			defer ms.Close()
			if err := ms.LoadResponsesFromDir(mock.ResponsesDir(tc.Mode)); err != nil {
				t.Fatalf("failed to load responses: %v", err)
			}

			resourceName := "emp_" + tc.ResourceName
			overrides := mergeOverrides(tc.Overrides, map[string]string{target: `""`})
			hcl := generateHCLWithExcludes(rs, tc.TerraformType, resourceName, tc.Mode, modeKey, overrides, nil)
			config := mock.ProviderConfig(ms.URL(), tc.Mode) + hcl

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
				Steps: []resource.TestStep{
					{
						PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), config) },
						Config:    config,
						Check: func(s *terraform.State) error {
							puts := ms.GetRequestsByMethodAndPath("PUT", tc.APIPath)
							if len(puts) == 0 {
								return fmt.Errorf("no PUT request captured for %s", tc.APIPath)
							}
							body := puts[len(puts)-1].Body
							basePath := tc.WrapperKey + "." + resourceName

							// Empty string must be PRESENT in PUT, not omitted
							mock.AssertFieldEquals(t, body, basePath+"."+target, "")
							return nil
						},
					},
				},
			})
		})
	}
}

func TestGeneric_Import(t *testing.T) {
	t.Parallel()
	for _, tc := range allResourceTests {
		t.Run(tc.TerraformType, func(t *testing.T) {
			t.Parallel()
			if tc.SkipCreate {
				t.Skipf("%s is update-only, skipping import", tc.TerraformType)
			}

			ms := mock.NewMockServer(tc.Mode)
			defer ms.Close()
			ms.SetTestLogger(t)
			if err := ms.LoadResponsesFromDir(mock.ResponsesDir(tc.Mode)); err != nil {
				t.Fatalf("failed to load responses: %v", err)
			}

			rs := inspectSchema(tc.Factory)
			modeKey := tc.modeFieldsKey()
			resourceName := "imp_" + tc.ResourceName

			hcl := generateCoverageHCL(rs, tc.TerraformType, resourceName, tc.Mode, modeKey, tc.Overrides)
			config := mock.ProviderConfig(ms.URL(), tc.Mode) + hcl

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
				Steps: []resource.TestStep{
					{
						PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), config) },
						Config:    config,
					},
					{
						PreConfig:                            func() { mock.WriteTFConfig(t, ms.URL(), config) },
						ResourceName:                         tc.TerraformType + ".test",
						ImportState:                          true,
						ImportStateId:                        resourceName,
						ImportStateVerify:                    true,
						ImportStateVerifyIdentifierAttribute: "name",
					},
				},
			})
		})
	}
}

func TestGeneric_NestedBlockOperations(t *testing.T) {
	t.Parallel()
	for _, tc := range allResourceTests {
		t.Run(tc.TerraformType, func(t *testing.T) {
			t.Parallel()
			if tc.SkipCreate {
				t.Skipf("%s is update-only, skipping", tc.TerraformType)
			}

			rs := inspectSchema(tc.Factory)
			modeKey := tc.modeFieldsKey()

			type modifiableBlock struct {
				block        blockInfo
				modField     string
				modFieldType string
			}
			var blocks []modifiableBlock
			for _, block := range rs.Blocks {
				if block.Name == "object_properties" {
					continue
				}
				if !utils.FieldAppliesToMode(modeKey, block.Name, tc.Mode) {
					continue
				}
				hasIndex := false
				for _, fi := range block.Fields {
					if fi.Name == "index" {
						hasIndex = true
						break
					}
				}
				if !hasIndex {
					continue
				}

				refPairedFields := make(map[string]bool)
				for _, fi := range block.Fields {
					if strings.HasSuffix(fi.Name, "_ref_type_") {
						refPairedFields[strings.TrimSuffix(fi.Name, "_ref_type_")] = true
					}
				}

				for _, fi := range block.Fields {
					if fi.Name == "index" {
						continue
					}
					if strings.HasSuffix(fi.Name, "_ref_type_") {
						continue
					}
					if refPairedFields[fi.Name] {
						continue
					}
					nestedKey := block.Name + "." + fi.Name
					if !utils.FieldAppliesToMode(modeKey, nestedKey, tc.Mode) {
						continue
					}
					if _, ok := tc.Overrides[nestedKey]; ok {
						continue
					}
					blocks = append(blocks, modifiableBlock{
						block:        block,
						modField:     fi.Name,
						modFieldType: fi.Type,
					})
					break
				}
			}
			if len(blocks) == 0 {
				t.Skipf("%s has no indexed nested blocks with modifiable fields, skipping", tc.TerraformType)
			}

			modifiedOverrides := make(map[string]string)
			for k, v := range tc.Overrides {
				modifiedOverrides[k] = v
			}
			mb := blocks[0]
			{
				nestedKey := mb.block.Name + "." + mb.modField
				switch mb.modFieldType {
				case "string":
					modifiedOverrides[nestedKey] = `"modified"`
				case "bool":
					modifiedOverrides[nestedKey] = "false"
				case "int64":
					modifiedOverrides[nestedKey] = "43"
				case "number":
					modifiedOverrides[nestedKey] = "2.5"
				default:
					modifiedOverrides[nestedKey] = `"modified"`
				}
			}

			ms := mock.NewMockServer(tc.Mode)
			defer ms.Close()
			ms.SetTestLogger(t)
			if err := ms.LoadResponsesFromDir(mock.ResponsesDir(tc.Mode)); err != nil {
				t.Fatalf("failed to load responses: %v", err)
			}

			resourceName := "blk_" + tc.ResourceName
			providerCfg := mock.ProviderConfig(ms.URL(), tc.Mode)

			// Config 1: 1 item per block (baseline create)
			config1 := providerCfg + generateHCLWithBlockCount(rs, tc.TerraformType, resourceName, tc.Mode, modeKey, tc.Overrides, 1)
			// Config 2: 2 items per block (add item 2)
			config2 := providerCfg + generateHCLWithBlockCount(rs, tc.TerraformType, resourceName, tc.Mode, modeKey, tc.Overrides, 2)
			// Config 3: 2 items per block with modified field (modify both items)
			config3 := providerCfg + generateHCLWithBlockCount(rs, tc.TerraformType, resourceName, tc.Mode, modeKey, modifiedOverrides, 2)
			// Config 4: 1 item per block with modified field (delete item 2)
			config4 := providerCfg + generateHCLWithBlockCount(rs, tc.TerraformType, resourceName, tc.Mode, modeKey, modifiedOverrides, 1)

			blockNames := make([]string, len(blocks))
			for i, mb := range blocks {
				blockNames[i] = fmt.Sprintf("%s (modify: %s)", mb.block.Name, mb.modField)
			}
			t.Logf("Testing %d indexed blocks: %v", len(blocks), blockNames)

			resource.UnitTest(t, resource.TestCase{
				ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
				Steps: []resource.TestStep{
					// Step 1: Create with 1 item (PUT)
					{
						PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), config1) },
						Config:    config1,
						Check: func(s *terraform.State) error {
							puts := ms.GetRequestsByMethodAndPath("PUT", tc.APIPath)
							if len(puts) == 0 {
								return fmt.Errorf("step 1: no PUT for %s", tc.APIPath)
							}
							return nil
						},
					},
					// Step 2: Add second item (PATCH — CreateNew pattern)
					{
						PreConfig: func() { ms.Reset(); mock.WriteTFConfig(t, ms.URL(), config2) },
						Config:    config2,
						Check: func(s *terraform.State) error {
							patches := ms.GetRequestsByMethodAndPath("PATCH", tc.APIPath)
							if len(patches) == 0 {
								return fmt.Errorf("step 2: no PATCH for %s", tc.APIPath)
							}
							body := patches[len(patches)-1].Body
							res, err := extractResourceFromBody(body, tc.WrapperKey, resourceName)
							if err != nil {
								return fmt.Errorf("step 2: %w", err)
							}
							for _, mb := range blocks {
								arr, ok := res[mb.block.Name].([]interface{})
								if !ok {
									return fmt.Errorf("step 2: block %q not found or not array in PATCH", mb.block.Name)
								}
								// Should have exactly 1 new item (item 2, CreateNew)
								if len(arr) != 1 {
									return fmt.Errorf("step 2: block %q expected 1 new item in PATCH, got %d", mb.block.Name, len(arr))
								}
								item, ok := arr[0].(map[string]interface{})
								if !ok {
									return fmt.Errorf("step 2: block %q[0] not a map", mb.block.Name)
								}
								if item["index"] != float64(2) {
									return fmt.Errorf("step 2: block %q new item index = %v, want 2", mb.block.Name, item["index"])
								}
								// CreateNew: should have ALL fields (more than just index)
								if len(item) < 2 {
									return fmt.Errorf("step 2: block %q new item should have all fields (CreateNew), got only %d", mb.block.Name, len(item))
								}
							}
							return nil
						},
					},
					// Step 3: Modify field in first block's items (PATCH — UpdateExisting pattern)
					{
						PreConfig: func() { ms.Reset(); mock.WriteTFConfig(t, ms.URL(), config3) },
						Config:    config3,
						Check: func(s *terraform.State) error {
							patches := ms.GetRequestsByMethodAndPath("PATCH", tc.APIPath)
							if len(patches) == 0 {
								return fmt.Errorf("step 3: no PATCH for %s", tc.APIPath)
							}
							body := patches[len(patches)-1].Body
							res, err := extractResourceFromBody(body, tc.WrapperKey, resourceName)
							if err != nil {
								return fmt.Errorf("step 3: %w", err)
							}
							// Only the first block was modified
							mb := blocks[0]
							arr, ok := res[mb.block.Name].([]interface{})
							if !ok {
								return fmt.Errorf("step 3: block %q not found or not array in PATCH", mb.block.Name)
							}
							for i, rawItem := range arr {
								item, ok := rawItem.(map[string]interface{})
								if !ok {
									return fmt.Errorf("step 3: block %q[%d] not a map", mb.block.Name, i)
								}
								if _, hasIndex := item["index"]; !hasIndex {
									return fmt.Errorf("step 3: block %q[%d] missing index", mb.block.Name, i)
								}
								if _, hasModField := item[mb.modField]; !hasModField {
									return fmt.Errorf("step 3: block %q[%d] missing modified field %q", mb.block.Name, i, mb.modField)
								}
								// UpdateExisting: should have index + changed field(s) only
								if len(item) > 3 {
									t.Logf("step 3 note: block %q[%d] has %d fields (expected ~2: index + %s), fields: %v",
										mb.block.Name, i, len(item), mb.modField, fieldKeys(item))
								}
							}
							return nil
						},
					},
					// Step 4: Remove item 2 (PATCH — CreateDeleted pattern)
					{
						PreConfig: func() { ms.Reset(); mock.WriteTFConfig(t, ms.URL(), config4) },
						Config:    config4,
						Check: func(s *terraform.State) error {
							patches := ms.GetRequestsByMethodAndPath("PATCH", tc.APIPath)
							if len(patches) == 0 {
								return fmt.Errorf("step 4: no PATCH for %s", tc.APIPath)
							}
							body := patches[len(patches)-1].Body
							res, err := extractResourceFromBody(body, tc.WrapperKey, resourceName)
							if err != nil {
								return fmt.Errorf("step 4: %w", err)
							}
							for _, mb := range blocks {
								arr, ok := res[mb.block.Name].([]interface{})
								if !ok {
									return fmt.Errorf("step 4: block %q not found or not array in PATCH", mb.block.Name)
								}
								// Should have exactly 1 item (the deleted item 2)
								if len(arr) != 1 {
									return fmt.Errorf("step 4: block %q expected 1 deleted item in PATCH, got %d", mb.block.Name, len(arr))
								}
								item, ok := arr[0].(map[string]interface{})
								if !ok {
									return fmt.Errorf("step 4: block %q[0] not a map", mb.block.Name)
								}
								if item["index"] != float64(2) {
									return fmt.Errorf("step 4: block %q deleted item index = %v, want 2", mb.block.Name, item["index"])
								}
								// CreateDeleted: should have ONLY index
								if len(item) != 1 {
									return fmt.Errorf("step 4: block %q deleted item should have only index (CreateDeleted), got %d fields: %v",
										mb.block.Name, len(item), fieldKeys(item))
								}
							}
							return nil
						},
					},
				},
			})
		})
	}
}

func extractResourceFromBody(body map[string]interface{}, wrapperKey, resourceName string) (map[string]interface{}, error) {
	wrapper, ok := body[wrapperKey].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("wrapper %q not found", wrapperKey)
	}
	res, ok := wrapper[resourceName].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("resource %q not found", resourceName)
	}
	return res, nil
}

func fieldKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

func TestGeneric_ModeFieldExclusion(t *testing.T) {
	t.Parallel()
	for _, tc := range allResourceTests {
		t.Run(tc.TerraformType, func(t *testing.T) {
			t.Parallel()
			if tc.SkipCreate {
				t.Skipf("%s is update-only, skipping", tc.TerraformType)
			}

			modeKey := tc.modeFieldsKey()
			rs := inspectSchema(tc.Factory)

			excluded := collectExcludedModeFields(rs, modeKey, tc.Mode)
			if len(excluded) > 0 {
				t.Run(tc.Mode, func(t *testing.T) {
					verifyModeExclusion(t, tc, rs, tc.Mode, "mfe_"+tc.ResourceName, excluded)
				})
			}

			altMode := "campus"
			if tc.Mode == "campus" {
				altMode = "datacenter"
			}

			altExcluded := collectExcludedModeFields(rs, modeKey, altMode)
			if len(altExcluded) == 0 {
				if len(excluded) == 0 {
					t.Skipf("%s has no mode-specific fields, skipping", tc.TerraformType)
				}
				return
			}

			hasApplicableFields := false
			for _, fi := range rs.Attributes {
				if utils.FieldAppliesToMode(modeKey, fi.Name, altMode) {
					hasApplicableFields = true
					break
				}
			}
			if !hasApplicableFields {
				return
			}

			pathSuffix := strings.TrimPrefix(tc.APIPath, "/api/")
			filename := pathSuffix + ".json"
			if tc.APIPath == "/api/acls" {
				if tc.TerraformType == "verity_acl_v4" {
					filename = "acls_ipv4.json"
				} else {
					filename = "acls_ipv6.json"
				}
			}
			if _, err := mock.LoadResponse(altMode, filename); err != nil {
				t.Logf("No %s response for %s, skipping alternate mode test", altMode, tc.TerraformType)
				return
			}

			t.Run(altMode, func(t *testing.T) {
				altTC := tc
				altTC.Mode = altMode
				verifyModeExclusion(t, altTC, rs, altMode, "alt_"+tc.ResourceName, altExcluded)
			})
		})
	}
}

func verifyModeExclusion(t *testing.T, tc ResourceCoverageEntry, rs resourceSchemaInfo, mode, resourceName string, excluded []string) {
	t.Helper()

	ms := mock.NewMockServer(mode)
	defer ms.Close()
	if err := ms.LoadResponsesFromDir(mock.ResponsesDir(mode)); err != nil {
		t.Fatalf("failed to load responses: %v", err)
	}

	modeKey := tc.modeFieldsKey()
	hcl := generateCoverageHCL(rs, tc.TerraformType, resourceName, mode, modeKey, tc.Overrides)
	config := mock.ProviderConfig(ms.URL(), mode) + hcl

	t.Logf("Mode: %s, excluded: %v", mode, excluded)

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: mock.ProtoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				PreConfig: func() { mock.WriteTFConfig(t, ms.URL(), config) },
				Config:    config,
				Check: func(s *terraform.State) error {
					puts := ms.GetRequestsByMethodAndPath("PUT", tc.APIPath)
					if len(puts) == 0 {
						return fmt.Errorf("no PUT for %s", tc.APIPath)
					}
					body := puts[len(puts)-1].Body

					wrapper, ok := body[tc.WrapperKey].(map[string]interface{})
					if !ok {
						return fmt.Errorf("wrapper %q not found", tc.WrapperKey)
					}
					res, ok := wrapper[resourceName].(map[string]interface{})
					if !ok {
						return fmt.Errorf("resource %q not found", resourceName)
					}

					for _, field := range excluded {
						if strings.Contains(field, ".") {
							parts := strings.SplitN(field, ".", 2)
							blockName, fieldName := parts[0], parts[1]
							if arr, ok := res[blockName].([]interface{}); ok {
								for i, item := range arr {
									if m, ok := item.(map[string]interface{}); ok {
										if _, exists := m[fieldName]; exists {
											t.Errorf("excluded field %s present in %s[%d]", field, blockName, i)
										}
									}
								}
							} else if obj, ok := res[blockName].(map[string]interface{}); ok {
								if _, exists := obj[fieldName]; exists {
									t.Errorf("excluded field %s present in %s object", field, blockName)
								}
							}
						} else {
							if _, exists := res[field]; exists {
								t.Errorf("excluded field %q present in PUT body (%s mode)", field, mode)
							}
						}
					}
					return nil
				},
			},
		},
	})
}
