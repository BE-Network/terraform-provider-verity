package genericresource

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-verity/internal/bulkops"
	"terraform-provider-verity/internal/transport"
	"terraform-provider-verity/internal/utils"
)

var errNotOnSystem = errors.New("not on the system")

type invalidResponseError struct{ message string }

func (e invalidResponseError) Error() string { return e.message }

func (r *Resource) adopt(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse, plan map[string]attr.Value) {
	if err := r.runtime.EnsureAuthenticated(ctx); err != nil {
		resp.Diagnostics.AddError("Failed to Authenticate", fmt.Sprintf("Error authenticating with API: %s", err))
		return
	}
	name := r.identityOf(plan)
	if !r.validateIdentity(plan, &resp.Diagnostics) {
		return
	}
	label := fmt.Sprintf("%s %s", r.spec.TerraformType, name)

	r.runtime.ClearCache(ctx, r.spec.API.CacheKey)
	current, err := r.fetchAdopted(ctx, name, plan)
	if err != nil {
		r.reportAdoptReadError(err, label, &resp.Diagnostics)
		return
	}
	nullables := r.nullableSource(ctx, req.Config, name, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	object, changed, err := buildUpdate(r.spec.Fields, plan, current, nullables, &resp.Diagnostics)
	if err != nil {
		resp.Diagnostics.AddError("Invalid Configuration", fmt.Sprintf("%s: %s", label, err))
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}
	if !changed {
		tflog.Info(ctx, fmt.Sprintf("Adopted existing %s; its settings already match the configuration", label))
		resp.Diagnostics.Append(r.setState(ctx, &resp.State, current)...)
		return
	}

	guardOutboundSecrets(r.spec.Fields, object, false, label, &resp.Diagnostics)
	if resp.Diagnostics.HasError() {
		return
	}
	value, err := r.adapter.ResourceValue(object)
	if err != nil {
		resp.Diagnostics.AddError("Failed to Build Request", fmt.Sprintf("%s: %s", label, err))
		return
	}
	var writeDiagnostics diag.Diagnostics
	if !bulkops.ExecuteResourceOperationWithOptions(ctx, r.runtime.BulkManager(), r.runtime.NotifyOperationAdded,
		"update", r.spec.API.BulkKey, name, value, &writeDiagnostics, r.operationOptions()) {
		resp.Diagnostics.Append(writeDiagnostics...)
		return
	}
	resp.Diagnostics.Append(writeDiagnostics...)

	confirmed, err := r.confirmAdopted(ctx, name, plan, nullables)
	if err != nil {
		resp.Diagnostics.AddError("Failed to Confirm Settings",
			fmt.Sprintf("The settings of %s were sent, but %s. Nothing was saved to Terraform state; the next apply reads the object again and sends whatever still differs.", label, err))
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Adopted existing %s and applied the configured settings", label))
	resp.Diagnostics.Append(r.setState(ctx, &resp.State, confirmed)...)
}

func (r *Resource) fetchAdopted(ctx context.Context, name string, prior map[string]attr.Value) (map[string]attr.Value, error) {
	collection, err := r.runtime.FetchCollectionStrict(ctx, r.spec)
	if errors.Is(err, transport.ErrMalformedCollection) {
		return nil, invalidResponseError{err.Error()}
	}
	if err != nil {
		return nil, err
	}
	data, exists := collection[name]
	if !exists {
		return nil, errNotOnSystem
	}
	object, ok := data.(map[string]interface{})
	if !ok {
		return nil, invalidResponseError{fmt.Sprintf("%s %s is not an object in the API response", r.spec.TerraformType, name)}
	}
	values, err := stateFromAPI(r.spec.Fields, preserveSecrets(r.spec.Fields, object, prior), r.runtime.Mode(), prior)
	if err != nil {
		return nil, invalidResponseError{fmt.Sprintf("%s %s: %s", r.spec.TerraformType, name, err)}
	}
	values[r.spec.IdentityPath] = types.StringValue(name)
	return values, nil
}

func (r *Resource) confirmAdopted(ctx context.Context, name string, plan map[string]attr.Value, nullables nullableSource) (map[string]attr.Value, error) {
	var lastErr error
	for attempt := 0; attempt <= bulkops.PostOperationVerificationRetries; attempt++ {
		if attempt > 0 {
			time.Sleep(bulkops.PostOperationVerificationBackoff)
		}
		r.runtime.ClearCache(ctx, r.spec.API.CacheKey)
		confirmed, err := r.fetchAdopted(ctx, name, plan)
		switch {
		case errors.Is(err, errNotOnSystem):
			lastErr = errors.New("the object is missing from the API response")
			continue
		case err != nil:
			lastErr = fmt.Errorf("reading it back failed: %s", err)
			continue
		}
		var scratch diag.Diagnostics
		remaining, changed, err := buildUpdate(r.spec.Fields, plan, confirmed, nullables, &scratch)
		if err != nil {
			lastErr = fmt.Errorf("its response could not be compared: %s", err)
			continue
		}
		if !changed && !scratch.HasError() {
			return confirmed, nil
		}
		fields := make([]string, 0, len(remaining))
		for field := range remaining {
			fields = append(fields, field)
		}
		sort.Strings(fields)
		lastErr = fmt.Errorf("the API still shows different values for %v", fields)
	}
	return nil, lastErr
}

func (r *Resource) reportAdoptReadError(err error, label string, diagnostics *diag.Diagnostics) {
	var invalid invalidResponseError
	switch {
	case errors.Is(err, errNotOnSystem):
		diagnostics.AddError("Settings Object Not Found",
			fmt.Sprintf("%s doesn't exist on the system. It's a settings object that can't be created through the API; check the name.", label))
	case errors.As(err, &invalid):
		diagnostics.AddError("Invalid Response", invalid.message)
	default:
		diagnostics.Append(utils.FormatOpenAPIError(err, fmt.Sprintf("Failed to Read %s", label))...)
	}
}
