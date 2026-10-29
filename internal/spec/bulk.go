package spec

import (
	"fmt"
	"sort"
)

func (r Registry) ValidateBulkOrders() error {
	for _, resource := range r {
		for mode := range resource.BulkOrder {
			if !resourceSupportsMode(resource, mode) {
				return fmt.Errorf("%s declares bulk order for unsupported mode %s", resource.TerraformType, mode)
			}
		}
		for _, mode := range resource.Modes {
			order, exists := resource.BulkOrder[mode]
			if !exists {
				return fmt.Errorf("%s has no %s bulk order", resource.TerraformType, mode)
			}
			for _, operation := range []struct {
				name      string
				rank      int
				supported bool
			}{
				{"PUT", order.Put, resource.Operations.Create},
				{"PATCH", order.Patch, resource.Operations.Update},
				{"DELETE", order.Delete, resource.Operations.Delete},
			} {
				if operation.rank < 0 || (operation.rank > 0) != operation.supported {
					return fmt.Errorf("%s %s %s bulk order %d disagrees with operation support %t", resource.TerraformType, mode, operation.name, operation.rank, operation.supported)
				}
			}
		}
	}
	for _, mode := range []Mode{ModeCampus, ModeDatacenter} {
		for _, operation := range []string{"PUT", "PATCH", "DELETE"} {
			byKey := make(map[string]int)
			byRank := make(map[int]string)
			for _, resource := range r {
				rank := resource.BulkOrder[mode].Rank(operation)
				if rank == 0 {
					continue
				}
				key := resource.API.BulkKey
				if previous, exists := byKey[key]; exists && previous != rank {
					return fmt.Errorf("%s %s bulk variants for %s disagree on order: %d and %d", mode, operation, key, previous, rank)
				}
				if previous, exists := byRank[rank]; exists && previous != key {
					return fmt.Errorf("%s %s bulk order %d is shared by %s and %s", mode, operation, rank, previous, key)
				}
				byKey[key], byRank[rank] = rank, key
			}
		}
	}
	for _, mode := range []Mode{ModeCampus, ModeDatacenter} {
		put, deleted := r.BulkOperationOrder(mode, "PUT"), r.BulkOperationOrder(mode, "DELETE")
		if len(put) != len(deleted) {
			return fmt.Errorf("%s DELETE order must reverse PUT order", mode)
		}
		for index, key := range put {
			if deleted[len(deleted)-1-index] != key {
				return fmt.Errorf("%s DELETE order must reverse PUT order", mode)
			}
		}
	}
	return nil
}

func (o BulkOrderSpec) Rank(operation string) int {
	switch operation {
	case "PUT":
		return o.Put
	case "PATCH":
		return o.Patch
	case "DELETE":
		return o.Delete
	default:
		return 0
	}
}

func (r Registry) BulkOperationOrder(mode Mode, operation string) []string {
	ranks := make(map[string]int)
	for _, resource := range r {
		if rank := resource.BulkOrder[mode].Rank(operation); rank > 0 {
			ranks[resource.API.BulkKey] = rank
		}
	}
	keys := make([]string, 0, len(ranks))
	for key := range ranks {
		keys = append(keys, key)
	}
	sort.Slice(keys, func(i, j int) bool {
		if ranks[keys[i]] == ranks[keys[j]] {
			return keys[i] < keys[j]
		}
		return ranks[keys[i]] < ranks[keys[j]]
	})
	return keys
}
