package utils

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// IndexedItem represents an item that has an Index field for identification
type IndexedItem interface {
	GetIndex() types.Int64
}

// IndexedItemHandler provides functions for creating and updating indexed items
type IndexedItemHandler[PlanType IndexedItem, APIType any] struct {
	// CreateNew creates a new API item from a plan item (for new items)
	CreateNew func(planItem PlanType) APIType

	// UpdateExisting creates an API item from plan and state items (for changed items)
	// Returns the updated item and whether any fields actually changed
	UpdateExisting func(planItem PlanType, stateItem PlanType) (APIType, bool)

	// CreateDeleted creates an API item for deletion (typically just with index)
	CreateDeleted func(index int64) APIType
}

// ProcessIndexedArrayUpdates processes updates for indexed array fields
// This is the common pattern used across static routes, pbits, queues, filters, etc.
func ProcessIndexedArrayUpdates[PlanType IndexedItem, APIType any](
	planItems []PlanType,
	stateItems []PlanType,
	handler IndexedItemHandler[PlanType, APIType],
) ([]APIType, bool) {
	// Build index map of existing items
	stateItemsByIndex := make(map[int64]PlanType)
	for _, item := range stateItems {
		if !item.GetIndex().IsNull() {
			stateItemsByIndex[item.GetIndex().ValueInt64()] = item
		}
	}

	var changedItems []APIType
	hasChanges := false

	// Process plan items (CREATE and UPDATE)
	for _, planItem := range planItems {
		if planItem.GetIndex().IsNull() {
			continue // Skip items without identifier
		}

		index := planItem.GetIndex().ValueInt64()
		stateItem, exists := stateItemsByIndex[index]

		if !exists {
			// CREATE: new item, include all fields
			newItem := handler.CreateNew(planItem)
			changedItems = append(changedItems, newItem)
			hasChanges = true
			continue
		}

		// UPDATE: existing item, check which fields changed
		updatedItem, fieldChanged := handler.UpdateExisting(planItem, stateItem)
		if fieldChanged {
			changedItems = append(changedItems, updatedItem)
			hasChanges = true
		}
	}

	// DELETE: Check for deleted items
	for stateIndex := range stateItemsByIndex {
		found := false
		for _, planItem := range planItems {
			if !planItem.GetIndex().IsNull() && planItem.GetIndex().ValueInt64() == stateIndex {
				found = true
				break
			}
		}

		if !found {
			// Item removed - create deletion marker
			deletedItem := handler.CreateDeleted(stateIndex)
			changedItems = append(changedItems, deletedItem)
			hasChanges = true
		}
	}

	return changedItems, hasChanges
}
