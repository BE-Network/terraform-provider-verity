package bulkops

import (
	"context"
	"encoding/json"
	"fmt"
	"terraform-provider-verity/openapi"

	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// ================================================================================================
// CIRCULAR REFERENCE FIX
// ================================================================================================
//
// Problem:
// - route_map_clause.match_vrf → references tenant.name
// - tenant.export_route_map/import_route_map → references route_map.name
// - route_map.route_map_clauses → references route_map_clause.name
// This creates: route_map_clause → tenant → route_map → route_map_clause
//
// Solution:
// - PUT (Create): Create route_map_clause with empty match_vrf, then PATCH to restore after tenant exists
// - DELETE: PATCH to clear match_vrf before deleting the complete circular set
//
// Key Design Principles:
// 1. Selective Application: Only resources with actual circular dependencies are affected
// 2. No .tf File Modification: All fixes happen via in-memory API manipulation or backend PATCH
// 3. No Auto-Fix for Invalid Configs: If user references non-existent tenant, API error surfaces
// 4. Proper Errors: Invalid configurations fail with clear API errors
// ================================================================================================

// detectCircularReferenceScenario checks if we need to apply the circular reference fix
// for route_map_clause and tenant resources during PUT operations.
func (m *Manager) detectCircularReferenceScenario(ctx context.Context) CircularPutInfo {
	result := CircularPutInfo{
		NeedsFix:        false,
		AffectedClauses: make(map[string]interface{}),
		ClauseNames:     []string{},
		TenantNames:     []string{},
	}

	m.mutex.Lock()
	routeMapClauseOps, hasClauseOps := m.resources["route_map_clause"]
	tenantOps, hasTenantOps := m.resources["tenant"]
	m.mutex.Unlock()

	// Check if both route_map_clause and tenant have PUT operations
	if !hasClauseOps || !hasTenantOps {
		return result
	}
	if len(routeMapClauseOps.Put) == 0 || len(tenantOps.Put) == 0 {
		return result
	}

	tenantsBeingCreated := make(map[string]bool)
	for tenantName := range tenantOps.Put {
		tenantsBeingCreated[tenantName] = true
	}

	// Only apply the fix to route_map_clauses that have match_vrf referencing
	// tenants that are also being created in this same batch.
	//
	// If a clause references an existing tenant, we do not apply the fix:
	// - If the tenant exists, the API should accept the request as-is
	// - If it doesn't exist, the user's .tf configuration is wrong and should fail with API error
	// We never auto-fix improperly structured requests - only handle the circular dependency case.
	circularRefClauses := make([]string, 0)
	referencedTenants := make(map[string]bool)

	// Check PUT operations - only include clauses referencing tenants being created
	for name, clauseData := range routeMapClauseOps.Put {
		if matchVrfValue := m.getMatchVrfValue(clauseData); matchVrfValue != "" {
			// Only apply fix if the referenced tenant is also being created (circular dependency)
			if tenantsBeingCreated[matchVrfValue] {
				result.AffectedClauses[name] = clauseData
				referencedTenants[matchVrfValue] = true
				circularRefClauses = append(circularRefClauses, name)
				tflog.Debug(ctx, fmt.Sprintf("Circular ref: route_map_clause '%s' references tenant '%s' being created", name, matchVrfValue))
			} else {
				// Clause references an existing tenant - let it go through as-is
				// If the tenant doesn't exist, API will return an error
				tflog.Debug(ctx, fmt.Sprintf("Skipping fix: route_map_clause '%s' references tenant '%s' (not being created)", name, matchVrfValue))
			}
		}
	}

	// Also check PATCH operations for match_vrf changes
	if len(routeMapClauseOps.Patch) > 0 {
		for name, patchData := range routeMapClauseOps.Patch {
			if matchVrfValue := m.getMatchVrfValue(patchData); matchVrfValue != "" {
				// Only apply fix if the referenced tenant is also being created
				if tenantsBeingCreated[matchVrfValue] {
					result.AffectedClauses[name] = patchData
					referencedTenants[matchVrfValue] = true
					circularRefClauses = append(circularRefClauses, name)
					tflog.Debug(ctx, fmt.Sprintf("Circular ref in PATCH: route_map_clause '%s' references tenant '%s' being created", name, matchVrfValue))
				} else {
					tflog.Debug(ctx, fmt.Sprintf("Skipping fix in PATCH: route_map_clause '%s' references tenant '%s' (not being created)", name, matchVrfValue))
				}
			}
		}
	}

	result.NeedsFix = len(result.AffectedClauses) > 0

	if result.NeedsFix {
		// Populate ClauseNames and TenantNames
		for clauseName := range result.AffectedClauses {
			result.ClauseNames = append(result.ClauseNames, clauseName)
		}
		for tenantName := range referencedTenants {
			result.TenantNames = append(result.TenantNames, tenantName)
		}

		tflog.Info(ctx, fmt.Sprintf("Applying circular reference fix for %d route_map_clause(s)", len(circularRefClauses)))
		tflog.Debug(ctx, "Circular reference fix details", map[string]interface{}{
			"route_map_clause_put_count":   len(routeMapClauseOps.Put),
			"route_map_clause_patch_count": len(routeMapClauseOps.Patch),
			"tenant_put_count":             len(tenantOps.Put),
			"tenants_being_created":        getMapKeys(tenantsBeingCreated),
			"circular_ref_clauses":         circularRefClauses,
			"total_affected_clauses":       len(result.AffectedClauses),
			"affected_clause_names":        result.ClauseNames,
			"referenced_tenant_names":      result.TenantNames,
		})
	}

	return result
}

// detectCompleteCircularDeleteSet checks if a complete circular reference set is being deleted.
// Returns info about the circular set only if all resources (clause, route_map, tenant) are being deleted together.
func (m *Manager) detectCompleteCircularDeleteSet(ctx context.Context) CircularDeleteInfo {
	result := CircularDeleteInfo{
		IsCompleteSet:   false,
		ClauseNames:     []string{},
		TenantNames:     []string{},
		RouteMapNames:   []string{},
		AffectedClauses: make(map[string]interface{}),
	}

	m.mutex.Lock()
	routeMapClauseOps, hasClauseOps := m.resources["route_map_clause"]
	tenantOps, hasTenantOps := m.resources["tenant"]
	routeMapOps, hasRouteMapOps := m.resources["route_map"]
	m.mutex.Unlock()

	// If any of the three resource types don't have DELETE operations, no circular set
	if !hasClauseOps || !hasTenantOps || !hasRouteMapOps {
		return result
	}
	if len(routeMapClauseOps.Delete) == 0 || len(tenantOps.Delete) == 0 || len(routeMapOps.Delete) == 0 {
		return result
	}

	tflog.Debug(ctx, "Checking for complete circular reference set in DELETE operations")

	// Fetch current state of route_map_clauses to check match_vrf values
	// (DELETE operations only store names, not full data)
	clausesWithMatchVrf := make(map[string]string) // clauseName -> tenantName

	tflog.Debug(ctx, "Fetching current state of route_map_clauses")

	clausesData, err := m.fetchCurrentRouteMapClauses(ctx)
	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf("Failed to fetch route_map_clauses: %v", err))
		return result
	}

	if len(clausesData) == 0 {
		tflog.Debug(ctx, "No route_map_clause data in response")
		return result
	}

	// Check each clause being deleted for match_vrf values
	for _, clauseName := range routeMapClauseOps.Delete {
		if clauseData, exists := clausesData[clauseName]; exists {
			matchVrfValue := m.getMatchVrfValue(clauseData)
			if matchVrfValue != "" {
				clausesWithMatchVrf[clauseName] = matchVrfValue
				// Store the full clause data for later PATCH generation
				result.AffectedClauses[clauseName] = clauseData
				tflog.Debug(ctx, fmt.Sprintf("Found route_map_clause being deleted with match_vrf: %s -> %s", clauseName, matchVrfValue))
			}
		}
	}

	if len(clausesWithMatchVrf) == 0 {
		tflog.Debug(ctx, "No route_map_clauses with match_vrf being deleted")
		return result
	}

	// Check if the tenants referenced by these clauses are also being deleted
	tenantsBeingDeleted := make(map[string]bool)
	for _, tenantName := range tenantOps.Delete {
		tenantsBeingDeleted[tenantName] = true
	}

	for clauseName, tenantName := range clausesWithMatchVrf {
		if tenantsBeingDeleted[tenantName] {
			// This clause+tenant pair forms part of a circular reference being deleted
			result.ClauseNames = append(result.ClauseNames, clauseName)
			if !contains(result.TenantNames, tenantName) {
				result.TenantNames = append(result.TenantNames, tenantName)
			}
			tflog.Debug(ctx, fmt.Sprintf("Circular reference detected: clause '%s' and tenant '%s' both being deleted", clauseName, tenantName))
		}
	}

	if len(result.ClauseNames) == 0 {
		tflog.Debug(ctx, "No complete circular references found (clauses reference tenants not being deleted)")
		return result
	}

	// Check if route_maps are also being deleted
	// This completes the circular chain: clause -> tenant -> route_map -> clause
	for _, routeMapName := range routeMapOps.Delete {
		if !contains(result.RouteMapNames, routeMapName) {
			result.RouteMapNames = append(result.RouteMapNames, routeMapName)
		}
	}

	if len(result.RouteMapNames) > 0 {
		result.IsCompleteSet = true
		tflog.Info(ctx, "Complete circular reference set detected in DELETE operations", map[string]interface{}{
			"clauses":    result.ClauseNames,
			"tenants":    result.TenantNames,
			"route_maps": result.RouteMapNames,
		})
	}

	return result
}

// ================================================================================================
// CIRCULAR REFERENCE - HELPER FUNCTIONS
// ================================================================================================

// getMatchVrfValue extracts the match_vrf value from route_map_clause data.
func (m *Manager) getMatchVrfValue(data interface{}) string {
	dataMap := dataToMap(data)
	if dataMap == nil {
		return ""
	}
	return extractStringField(dataMap, "match_vrf")
}

// createRouteMapClauseWithEmptyMatchVrf creates a copy of route_map_clause data with match_vrf set to empty.
func (m *Manager) createRouteMapClauseWithEmptyMatchVrf(originalData interface{}) interface{} {
	jsonData, err := json.Marshal(originalData)
	if err != nil {
		return originalData
	}

	var result openapi.RoutemapclausesPutRequestRouteMapClauseValue
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return originalData
	}

	result.MatchVrf = openapi.PtrString("")
	return result
}

// fetchCurrentRouteMapClauses fetches the current state of route_map_clauses from the API.
// Returns a map of clause name to clause data.
func (m *Manager) fetchCurrentRouteMapClauses(ctx context.Context) (map[string]interface{}, error) {
	bgCtx := context.Background()

	resp, err := m.client.RouteMapClausesAPI.RoutemapclausesGet(bgCtx).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch route_map_clauses: %w", err)
	}
	defer resp.Body.Close()

	var clausesResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&clausesResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	clausesData, ok := clausesResponse["route_map_clause"].(map[string]interface{})
	if !ok {
		return make(map[string]interface{}), nil
	}

	return clausesData, nil
}

// createMatchVrfPatchData creates PATCH data with only match_vrf and match_vrf_ref_type_ fields.
// Can be used for both restoring match_vrf (PUT fix) and clearing it (DELETE fix).
func (m *Manager) createMatchVrfPatchData(originalData interface{}, matchVrfValue string) interface{} {
	dataMap := dataToMap(originalData)
	if dataMap == nil {
		return nil
	}

	patchData := openapi.RoutemapclausesPutRequestRouteMapClauseValue{}

	name := extractStringField(dataMap, "name")
	if name == "" {
		return nil // Cannot create PATCH data without a name
	}
	patchData.Name = openapi.PtrString(name)

	// Set match_vrf to provided value (can be empty string for clearing)
	patchData.MatchVrf = openapi.PtrString(matchVrfValue)

	// Include match_vrf_ref_type_ if it exists
	if refType := extractStringField(dataMap, "match_vrf_ref_type_"); refType != "" {
		patchData.MatchVrfRefType = openapi.PtrString(refType)
	}

	return patchData
}

// ================================================================================================
// CIRCULAR REFERENCE - UTILITY FUNCTIONS
// ================================================================================================

// dataToMap converts any data structure to map[string]interface{} for inspection.
// Returns nil if conversion fails.
func dataToMap(data interface{}) map[string]interface{} {
	if dataMap, ok := data.(map[string]interface{}); ok {
		return dataMap
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil
	}

	var result map[string]interface{}
	if err := json.Unmarshal(jsonData, &result); err != nil {
		return nil
	}
	return result
}

// extractStringField safely extracts a string field from a map or pointer.
func extractStringField(dataMap map[string]interface{}, fieldName string) string {
	if value, exists := dataMap[fieldName]; exists {
		if strValue, ok := value.(string); ok {
			return strValue
		}
		if ptrValue, ok := value.(*string); ok && ptrValue != nil {
			return *ptrValue
		}
	}
	return ""
}

// getMapKeys returns the keys of a map[string]bool as a slice.
func getMapKeys(m map[string]bool) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// contains checks if a string slice contains a specific string.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
