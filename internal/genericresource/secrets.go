package genericresource

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-verity/internal/spec"
	"terraform-provider-verity/internal/transport"
)

func preserveSecrets(fields []spec.FieldSpec, data map[string]interface{}, source map[string]attr.Value) map[string]interface{} {
	var result map[string]interface{}
	for _, field := range fields {
		if field.Secret == "" || field.Unmanaged {
			continue
		}
		known, ok := source[field.TerraformName].(types.String)
		if !ok || known.IsNull() || known.IsUnknown() {
			continue
		}
		raw, present := data[field.APIName]
		text, isString := raw.(string)
		replace := false
		switch field.Secret {
		case spec.SecretWriteOnly:
			replace = !present || raw == nil || (isString && text == "")
		case spec.SecretMasked:
			replace = isString && spec.IsMaskedSecretPlaceholder(text)
		}
		if !replace {
			continue
		}
		if result == nil {
			result = make(map[string]interface{}, len(data))
			for key, value := range data {
				result[key] = value
			}
		}
		result[field.APIName] = known.ValueString()
	}
	if result == nil {
		return data
	}
	return result
}

func isRealSecret(value attr.Value) bool {
	text, ok := value.(types.String)
	if !ok || text.IsNull() || text.IsUnknown() {
		return false
	}
	return text.ValueString() != "" && !spec.IsMaskedSecretPlaceholder(text.ValueString())
}

func isPlaceholder(value attr.Value) bool {
	text, ok := value.(types.String)
	return ok && !text.IsNull() && !text.IsUnknown() && spec.IsMaskedSecretPlaceholder(text.ValueString())
}

func (r *Resource) planSecrets(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	creating := req.State.Raw.IsNull()
	for _, field := range r.spec.Fields {
		if field.Secret != spec.SecretMasked || field.Unmanaged {
			continue
		}
		attributePath := path.Root(field.TerraformName)
		var planned types.String
		resp.Diagnostics.Append(resp.Plan.GetAttribute(ctx, attributePath, &planned)...)
		if resp.Diagnostics.HasError() {
			return
		}
		var prior types.String
		if !creating {
			resp.Diagnostics.Append(req.State.GetAttribute(ctx, attributePath, &prior)...)
			if resp.Diagnostics.HasError() {
				return
			}
		}

		if isPlaceholder(planned) {
			switch {
			case creating:
				resp.Diagnostics.AddAttributeError(attributePath, "Password placeholder in configuration",
					fmt.Sprintf("%s is set to %q, which is the API's placeholder for an existing password, not a password. Remove the attribute or set the real password.", field.TerraformName, planned.ValueString()))
			case !isPlaceholder(prior):
				resp.Diagnostics.AddAttributeError(attributePath, "Password placeholder in configuration",
					fmt.Sprintf("%s would change to %q, which is the API's placeholder, not a password. Remove the attribute or set the real password.", field.TerraformName, planned.ValueString()))
			default:
				resp.Diagnostics.AddAttributeWarning(attributePath, "Password placeholder in configuration",
					fmt.Sprintf("%s is set to the API's placeholder %q, probably from an earlier import. It is not sent to the API. Remove the attribute, or set the real password if Terraform should manage it.", field.TerraformName, planned.ValueString()))
			}
		}

		if creating || field.AutoAssignment == nil {
			continue
		}
		flagPath := path.Root(field.AutoAssignment.FlagField)
		var priorFlag, plannedFlag types.Bool
		resp.Diagnostics.Append(req.State.GetAttribute(ctx, flagPath, &priorFlag)...)
		resp.Diagnostics.Append(resp.Plan.GetAttribute(ctx, flagPath, &plannedFlag)...)
		if resp.Diagnostics.HasError() {
			return
		}
		turningOff := !priorFlag.IsNull() && !priorFlag.IsUnknown() && priorFlag.ValueBool() &&
			!plannedFlag.IsNull() && !plannedFlag.IsUnknown() && !plannedFlag.ValueBool()
		sent := attr.Value(planned)
		if planned.IsNull() || planned.IsUnknown() {
			sent = prior
		}
		if turningOff && !isRealSecret(sent) && !r.plaintextCompanionIsReal(ctx, field, resp) {
			resp.Diagnostics.AddAttributeError(flagPath, "Password required",
				fmt.Sprintf("Turning off %s requires setting the password (%s or %s), because the API doesn't return the current one.", field.AutoAssignment.FlagField, plaintextName(field), field.TerraformName))
		}
	}
}

func plaintextName(field spec.FieldSpec) string {
	return strings.TrimSuffix(field.TerraformName, "_encrypted")
}

func (r *Resource) plaintextCompanionIsReal(ctx context.Context, field spec.FieldSpec, resp *resource.ModifyPlanResponse) bool {
	name := plaintextName(field)
	if name == field.TerraformName {
		return false
	}
	for _, candidate := range r.spec.Fields {
		if candidate.TerraformName != name || candidate.Secret != spec.SecretWriteOnly {
			continue
		}
		var planned types.String
		if diags := resp.Plan.GetAttribute(ctx, path.Root(name), &planned); diags.HasError() {
			return false
		}
		return isRealSecret(planned)
	}
	return false
}

func hasRealPlaintext(fields []spec.FieldSpec, field spec.FieldSpec, object transport.WireObject) bool {
	name := plaintextName(field)
	for _, candidate := range fields {
		if candidate.TerraformName != name || candidate.Secret != spec.SecretWriteOnly || name == field.TerraformName {
			continue
		}
		value, present := object[candidate.APIName]
		return present && value.Kind() == transport.ValueKindString && value.StringValue() != "" && !spec.IsMaskedSecretPlaceholder(value.StringValue())
	}
	return false
}

func resendPlaintextOnFlagOff(fields []spec.FieldSpec, plan, state map[string]attr.Value, object transport.WireObject) {
	for _, field := range fields {
		if field.Secret != spec.SecretMasked || field.Unmanaged || field.AutoAssignment == nil {
			continue
		}
		priorFlag, priorOK := state[field.AutoAssignment.FlagField].(types.Bool)
		plannedFlag, plannedOK := plan[field.AutoAssignment.FlagField].(types.Bool)
		turningOff := priorOK && plannedOK && isKnown(priorFlag) && priorFlag.ValueBool() &&
			isKnown(plannedFlag) && !plannedFlag.ValueBool()
		if !turningOff {
			continue
		}
		name := plaintextName(field)
		for _, candidate := range fields {
			if candidate.TerraformName != name || candidate.Secret != spec.SecretWriteOnly || name == field.TerraformName {
				continue
			}
			if _, present := object[candidate.APIName]; present {
				break
			}
			if planned := plan[name]; isRealSecret(planned) {
				object[candidate.APIName] = transport.String(planned.(types.String).ValueString())
			}
		}
	}
}

func guardOutboundSecrets(fields []spec.FieldSpec, object transport.WireObject, creating bool, label string, diagnostics *diag.Diagnostics) {
	for _, field := range fields {
		if field.Secret != spec.SecretMasked || field.Unmanaged {
			continue
		}
		value, present := object[field.APIName]
		if !present || value.Kind() != transport.ValueKindString || !spec.IsMaskedSecretPlaceholder(value.StringValue()) {
			continue
		}
		delete(object, field.APIName)
		if hasRealPlaintext(fields, field, object) {
			continue
		}
		if creating {
			diagnostics.AddError("Password placeholder not sent",
				fmt.Sprintf("%s: %s contains the API's placeholder, not a password. Set the real password or remove the attribute.", label, field.TerraformName))
			continue
		}
		diagnostics.AddWarning("Password placeholder not sent",
			fmt.Sprintf("%s: %s contains the API's placeholder, not a password, so it was left out of the update.", label, field.TerraformName))
	}
}
