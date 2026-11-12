package bulkops

import (
	"context"
	"fmt"
	"net/http"
	"sort"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-verity/openapi"
)

func (m *Manager) clearMatchVrfBeforeDeletes(ctx context.Context) {
	m.mutex.Lock()
	tenants := map[string]bool{}
	var clauses []string
	if res, exists := m.resources["tenant"]; exists {
		for _, name := range res.Delete {
			tenants[name] = true
		}
	}
	if res, exists := m.resources["route_map_clause"]; exists {
		clauses = append(clauses, res.Delete...)
	}
	m.mutex.Unlock()
	if len(tenants) == 0 || len(clauses) == 0 {
		return
	}

	config, exists := resourceRegistry["route_map_clause"]
	if !exists || config.GetFunc == nil || config.PreparePatch == nil || config.PatchFunc == nil {
		return
	}
	fetchCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
	defer cancel()
	current, err := m.fetchResourceResponse(fetchCtx, ctx, config, nil)
	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf("[BULK-OPS] Could not read route map clauses before deleting tenants: %v", err))
		return
	}

	selected := map[string]interface{}{}
	var names []string
	for _, name := range clauses {
		object, _ := current[name].(map[string]interface{})
		tenant, _ := object["match_vrf"].(string)
		refType, _ := object["match_vrf_ref_type_"].(string)
		if tenant == "" || !tenants[tenant] || (refType != "" && refType != "tenant") {
			continue
		}
		value := openapi.RoutemapclausesPutRequestRouteMapClauseValue{}
		value.SetMatchVrf("")
		value.SetMatchVrfRefType("")
		selected[name] = value
		names = append(names, name)
	}
	if len(selected) == 0 {
		return
	}
	sort.Strings(names)

	request, err := config.PreparePatch(selected)
	if err == nil {
		_, err = m.sendWithRetries(ctx, BulkOperationConfig{
			ResourceType:  "route_map_clause",
			OperationType: "PATCH",
			ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
				return config.PatchFunc(m.client, ctx, request)
			},
		}, request)
	}
	if err != nil {
		tflog.Warn(ctx, fmt.Sprintf("[BULK-OPS] Could not clear match_vrf on route map clauses %v before deleting their tenants: %v", names, err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("[BULK-OPS] Cleared match_vrf on route map clauses %v before deleting them together with their tenants", names))
}
