package utils

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// NestedBlockUpdateHandler defines the interface for handling nested block updates
// This provides the generic pattern for add/update/delete operations on indexed nested blocks
type NestedBlockUpdateHandler[TModel, TApi any] interface {
	// GetIndex extracts the index from a model, returns -1 if no valid index
	GetIndex(TModel) int64

	// CreateNew creates a new API object from a model (for adding new items)
	CreateNew(TModel, *diag.Diagnostics) (TApi, error)

	// UpdateExisting creates an API object with only changed fields (for updating existing items)
	UpdateExisting(plan, state TModel, diags *diag.Diagnostics) (TApi, bool, error)

	// DeleteByIndex creates an API object for deletion (usually just contains the index)
	DeleteByIndex(int64) TApi
}

// ProcessNestedBlockUpdates implements the generic nested block update pattern:
// 1. Build index map of current state
// 2. For each plan item: check if new/existing and handle accordingly
// 3. For each state item not in plan: mark for deletion
// 4. Return changes and whether any changes occurred
func ProcessNestedBlockUpdates[TModel, TApi any](
	planItems, stateItems []TModel,
	handler NestedBlockUpdateHandler[TModel, TApi],
	diagnostics *diag.Diagnostics,
) ([]TApi, bool) {
	// Step 1: Build index map of current state
	stateByIndex := make(map[int64]TModel)
	for _, item := range stateItems {
		if index := handler.GetIndex(item); index >= 0 {
			stateByIndex[index] = item
		}
	}

	var changedItems []TApi
	hasChanges := false

	// Step 2: Process plan items (new or updated)
	for _, planItem := range planItems {
		index := handler.GetIndex(planItem)
		if index < 0 {
			continue // Skip items without valid index
		}

		stateItem, existsInState := stateByIndex[index]

		if !existsInState {
			// New item - include all fields
			apiItem, err := handler.CreateNew(planItem, diagnostics)
			if err != nil {
				return nil, false
			}
			changedItems = append(changedItems, apiItem)
			hasChanges = true
		} else {
			// Existing item - check for changes
			apiItem, itemChanged, err := handler.UpdateExisting(planItem, stateItem, diagnostics)
			if err != nil {
				return nil, false
			}
			if itemChanged {
				changedItems = append(changedItems, apiItem)
				hasChanges = true
			}
		}
	}

	// Step 3: Check for deleted items (in state but not in plan)
	planIndexes := make(map[int64]bool)
	for _, planItem := range planItems {
		if index := handler.GetIndex(planItem); index >= 0 {
			planIndexes[index] = true
		}
	}

	for index := range stateByIndex {
		if !planIndexes[index] {
			// Item was deleted
			deletedItem := handler.DeleteByIndex(index)
			changedItems = append(changedItems, deletedItem)
			hasChanges = true
		}
	}

	return changedItems, hasChanges
}
