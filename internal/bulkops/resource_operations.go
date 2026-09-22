package bulkops

import (
	"context"
	"fmt"
	"terraform-provider-verity/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

type ResourceOperationOptions struct {
	HeaderParams map[string]string
}

func ExecuteResourceOperationWithOptions(
	ctx context.Context,
	bulkOpsMgr *Manager,
	notifyFunc func(),
	operationType, resourceType, resourceName string,
	resourceData interface{},
	diagnostics *diag.Diagnostics,
	options *ResourceOperationOptions,
) bool {
	var operationID string

	var headerParams map[string]string
	if options != nil && options.HeaderParams != nil {
		headerParams = options.HeaderParams
	}

	switch operationType {
	case "create":
		operationID = bulkOpsMgr.AddPut(ctx, resourceType, resourceName, resourceData, headerParams)
	case "update":
		operationID = bulkOpsMgr.AddPatch(ctx, resourceType, resourceName, resourceData, headerParams)
	case "delete":
		operationID = bulkOpsMgr.AddDelete(ctx, resourceType, resourceName, headerParams)
	default:
		diagnostics.AddError("Invalid Operation Type", fmt.Sprintf("Unknown operation type: %s", operationType))
		return false
	}

	notifyFunc()

	if err := bulkOpsMgr.WaitForOperation(ctx, operationID, OperationTimeout); err != nil {
		diagnostics.Append(
			utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to %s %s %s", operationType, resourceType, resourceName))...,
		)
		return false
	}

	return true
}
