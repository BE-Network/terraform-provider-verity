package bulkops

import (
	"context"
	"fmt"
	"terraform-provider-verity/internal/utils"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// ResourceOperationOptions holds optional configuration for resource operations
type ResourceOperationOptions struct {
	// HeaderParams holds custom header parameters for API requests.
	// For example, ACL resources use "ip_version": "4" or "6"
	HeaderParams map[string]string
}

// ExecuteResourceOperation handles the common pattern of bulk operations (Put/Patch/Delete)
// and waiting for completion with proper error handling
func ExecuteResourceOperation(
	ctx context.Context,
	bulkOpsMgr *Manager,
	notifyFunc func(),
	operationType, resourceType, resourceName string,
	resourceData interface{},
	diagnostics *diag.Diagnostics,
) bool {
	return ExecuteResourceOperationWithOptions(ctx, bulkOpsMgr, notifyFunc, operationType, resourceType, resourceName, resourceData, diagnostics, nil)
}

// ExecuteResourceOperationWithOptions handles resource operations with optional configuration
// such as IP version for ACL resources
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

	// Extract header parameters if provided
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
