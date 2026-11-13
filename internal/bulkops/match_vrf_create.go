package bulkops

import (
	"context"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-verity/openapi"
)

const routeMapClauseType = "route_map_clause"

type matchVrfValue struct {
	matchVrf string
	refType  string
}

func splitAfter(order []string, last string) ([]string, []string) {
	for i, resourceType := range order {
		if resourceType == last {
			return order[:i+1], order[i+1:]
		}
	}
	return order, nil
}

func (m *Manager) armMatchVrfHold() {
	m.operationMutex.Lock()
	m.matchVrfHoldArmed = true
	m.operationMutex.Unlock()
}

func (m *Manager) disarmMatchVrfHold(err error) {
	m.failHeld(err)
	m.operationMutex.Lock()
	m.matchVrfHoldArmed = false
	m.operationMutex.Unlock()
}

func (m *Manager) holdMatchVrfAtExtraction(ctx context.Context, operations map[string]interface{}) {
	m.mutex.Lock()
	tenants := map[string]bool{}
	if res, exists := m.resources["tenant"]; exists {
		for name := range res.Put {
			tenants[name] = true
		}
	}
	m.mutex.Unlock()
	if len(tenants) == 0 {
		return
	}

	m.operationMutex.Lock()
	defer m.operationMutex.Unlock()
	if !m.matchVrfHoldArmed {
		return
	}
	var held []string
	for name, raw := range operations {
		value, ok := raw.(openapi.RoutemapclausesPutRequestRouteMapClauseValue)
		if !ok {
			continue
		}
		tenant, refType := value.GetMatchVrf(), value.GetMatchVrfRefType()
		if tenant == "" || !tenants[tenant] || (refType != "" && refType != "tenant") {
			continue
		}
		m.heldPuts[name] = matchVrfValue{matchVrf: tenant, refType: refType}
		value.SetMatchVrf("")
		value.SetMatchVrfRefType("")
		operations[name] = value
		held = append(held, name)
	}
	if len(held) > 0 {
		sort.Strings(held)
		tflog.Info(ctx, fmt.Sprintf("[BULK-OPS] Creating route map clauses %v without match_vrf; it is set once their tenants exist", held))
	}
}

func (m *Manager) keepClauseOpen(name string, existing interface{}, queued interface{}) bool {
	m.operationMutex.Lock()
	defer m.operationMutex.Unlock()
	if !m.matchVrfHoldArmed {
		return false
	}
	if _, held := m.heldPuts[name]; held {
		return true
	}
	value, ok := queued.(openapi.RoutemapclausesPutRequestRouteMapClauseValue)
	if !ok || value.GetMatchVrf() == "" {
		return false
	}
	want, wantRef := value.GetMatchVrf(), value.GetMatchVrfRefType()
	if matchVrfShows(existing, want, wantRef) {
		return false
	}
	m.heldPuts[name] = matchVrfValue{matchVrf: want, refType: wantRef}
	return true
}

func matchVrfShows(object interface{}, want, wantRef string) bool {
	fields, _ := object.(map[string]interface{})
	got, _ := fields["match_vrf"].(string)
	gotRef, _ := fields["match_vrf_ref_type_"].(string)
	return got == want && (wantRef == "" || gotRef == wantRef)
}

func (m *Manager) isHeldLocked(resourceType, operationType, name string) bool {
	if resourceType != routeMapClauseType || operationType != "PUT" {
		return false
	}
	_, held := m.heldPuts[name]
	return held
}

func (m *Manager) heldSnapshot() (map[string]matchVrfValue, []string) {
	m.operationMutex.Lock()
	defer m.operationMutex.Unlock()
	held := make(map[string]matchVrfValue, len(m.heldPuts))
	names := make([]string, 0, len(m.heldPuts))
	for name, value := range m.heldPuts {
		held[name] = value
		names = append(names, name)
	}
	sort.Strings(names)
	return held, names
}

func (m *Manager) failHeld(err error) {
	m.operationMutex.Lock()
	names := make([]string, 0, len(m.heldPuts))
	for name := range m.heldPuts {
		names = append(names, name)
	}
	m.heldPuts = make(map[string]matchVrfValue)
	m.operationMutex.Unlock()
	if len(names) == 0 {
		return
	}
	m.forgetResponses(routeMapClauseType, names)
	m.settleOperations(routeMapClauseType, "PUT", names, err, "")
}

func (m *Manager) forgetResponses(resourceType string, names []string) {
	res, exists := m.resources[resourceType]
	if !exists {
		return
	}
	res.ResponsesMutex.Lock()
	for _, name := range names {
		delete(res.Responses, name)
	}
	res.ResponsesMutex.Unlock()
}

func (m *Manager) restoreMatchVrf(ctx context.Context) error {
	held, names := m.heldSnapshot()
	if len(held) == 0 {
		return nil
	}
	config, exists := resourceRegistry[routeMapClauseType]
	if !exists || config.PreparePatch == nil || config.PatchFunc == nil || config.GetFunc == nil {
		return fmt.Errorf("route map clauses %v were created without match_vrf, and it can't be set", names)
	}

	patch := make(map[string]interface{}, len(held))
	for name, value := range held {
		object := openapi.RoutemapclausesPutRequestRouteMapClauseValue{}
		object.SetMatchVrf(value.matchVrf)
		object.SetMatchVrfRefType(value.refType)
		patch[name] = object
	}
	request, err := config.PreparePatch(patch)
	if err == nil {
		_, err = m.sendWithRetries(ctx, BulkOperationConfig{
			ResourceType:  routeMapClauseType,
			OperationType: "PATCH",
			ExecuteRequest: func(ctx context.Context, request interface{}) (*http.Response, error) {
				return config.PatchFunc(m.client, ctx, request)
			},
		}, request)
	}
	if err != nil {
		return fmt.Errorf("route map clauses %v were created, but match_vrf could not be set: %s", names, bulkErrorDetail(err))
	}

	confirmed, err := m.confirmMatchVrf(ctx, config, held)
	if err != nil {
		return fmt.Errorf("match_vrf was sent for route map clauses %v but could not be confirmed: %v", names, err)
	}
	if res, exists := m.resources[routeMapClauseType]; exists {
		m.cacheResourceResponse(ctx, config, res, nil, confirmed)
	}

	m.operationMutex.Lock()
	m.heldPuts = make(map[string]matchVrfValue)
	m.operationMutex.Unlock()
	m.settleOperations(routeMapClauseType, "PUT", names, nil, "")
	tflog.Info(ctx, fmt.Sprintf("[BULK-OPS] Set match_vrf on route map clauses %v after creating their tenants", names))
	return nil
}

func (m *Manager) confirmMatchVrf(ctx context.Context, config ResourceConfig, held map[string]matchVrfValue) (map[string]interface{}, error) {
	var lastErr error
	for attempt := 0; attempt <= PostOperationVerificationRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(PostOperationVerificationBackoff)
		}
		fetchCtx, cancel := context.WithTimeout(context.Background(), OperationTimeout)
		current, err := m.fetchResourceResponse(fetchCtx, ctx, config, nil)
		cancel()
		if err != nil {
			lastErr = err
			continue
		}
		var different []string
		for name, value := range held {
			if !matchVrfShows(current[name], value.matchVrf, value.refType) {
				different = append(different, name)
			}
		}
		if len(different) == 0 {
			return current, nil
		}
		sort.Strings(different)
		lastErr = fmt.Errorf("the API still shows a different match_vrf for %v", different)
	}
	return nil, lastErr
}
