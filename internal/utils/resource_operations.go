package utils

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// ExecuteResourceOperation handles the common pattern of bulk operations (Put/Patch/Delete)
// and waiting for completion with proper error handling
func ExecuteResourceOperation(
	ctx context.Context,
	bulkOpsMgr *BulkOperationManager,
	notifyFunc func(),
	operationType, resourceType, resourceName string,
	resourceData interface{},
	diagnostics *diag.Diagnostics,
) bool {
	var operationID string

	switch operationType {
	case "create":
		operationID = bulkOpsMgr.AddPut(ctx, resourceType, resourceName, resourceData)
	case "update":
		operationID = bulkOpsMgr.AddPatch(ctx, resourceType, resourceName, resourceData)
	case "delete":
		operationID = bulkOpsMgr.AddDelete(ctx, resourceType, resourceName)
	default:
		diagnostics.AddError("Invalid Operation Type", fmt.Sprintf("Unknown operation type: %s", operationType))
		return false
	}

	notifyFunc()

	if err := bulkOpsMgr.WaitForOperation(ctx, operationID, OperationTimeout); err != nil {
		diagnostics.Append(
			FormatOpenAPIError(err, fmt.Sprintf("Failed to %s %s %s", operationType, resourceType, resourceName))...,
		)
		return false
	}

	return true
}
