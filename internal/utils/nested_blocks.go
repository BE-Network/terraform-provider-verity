package utils

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type IndexedItem interface {
	GetIndex() types.Int64
}

type IndexedItemHandler[PlanType IndexedItem, APIType any] struct {
	CreateNew func(planItem PlanType) APIType

	UpdateExisting func(planItem PlanType, stateItem PlanType) (APIType, bool)

	CreateDeleted func(index int64) APIType
}

func ProcessIndexedArrayUpdates[PlanType IndexedItem, APIType any](
	planItems []PlanType,
	stateItems []PlanType,
	handler IndexedItemHandler[PlanType, APIType],
) ([]APIType, bool) {

	stateItemsByIndex := make(map[int64]PlanType)
	for _, item := range stateItems {
		if !item.GetIndex().IsNull() {
			stateItemsByIndex[item.GetIndex().ValueInt64()] = item
		}
	}

	var changedItems []APIType
	hasChanges := false

	for _, planItem := range planItems {
		if planItem.GetIndex().IsNull() {
			continue
		}

		index := planItem.GetIndex().ValueInt64()
		stateItem, exists := stateItemsByIndex[index]

		if !exists {

			newItem := handler.CreateNew(planItem)
			changedItems = append(changedItems, newItem)
			hasChanges = true
			continue
		}

		updatedItem, fieldChanged := handler.UpdateExisting(planItem, stateItem)
		if fieldChanged {
			changedItems = append(changedItems, updatedItem)
			hasChanges = true
		}
	}

	for stateIndex := range stateItemsByIndex {
		found := false
		for _, planItem := range planItems {
			if !planItem.GetIndex().IsNull() && planItem.GetIndex().ValueInt64() == stateIndex {
				found = true
				break
			}
		}

		if !found {

			deletedItem := handler.CreateDeleted(stateIndex)
			changedItems = append(changedItems, deletedItem)
			hasChanges = true
		}
	}

	return changedItems, hasChanges
}

func FilterIndexedEntries[T IndexedItem](stateItems []T, refItems []T) []T {
	if len(refItems) == 0 || len(stateItems) == 0 {
		return stateItems
	}

	refIndices := make(map[int64]bool, len(refItems))
	for _, item := range refItems {
		idx := item.GetIndex()
		if !idx.IsNull() && !idx.IsUnknown() {
			refIndices[idx.ValueInt64()] = true
		}
	}

	var filtered []T
	for _, item := range stateItems {
		idx := item.GetIndex()
		if !idx.IsNull() && !idx.IsUnknown() && refIndices[idx.ValueInt64()] {
			filtered = append(filtered, item)
		}
	}

	if len(filtered) > 0 {
		return filtered
	}
	return stateItems
}
