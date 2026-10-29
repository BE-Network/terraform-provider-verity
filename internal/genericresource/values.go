package genericresource

import (
	"context"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-verity/internal/spec"
	"terraform-provider-verity/internal/transport"
	"terraform-provider-verity/internal/utils"
)

type attributeSource interface {
	GetAttribute(ctx context.Context, p path.Path, target interface{}) diag.Diagnostics
}

func readScalars(ctx context.Context, source attributeSource, fields []spec.FieldSpec) (map[string]attr.Value, diag.Diagnostics) {
	var diagnostics diag.Diagnostics
	values := make(map[string]attr.Value, len(fields))
	for _, field := range fields {
		if field.Unmanaged {
			continue
		}
		attributePath := path.Root(field.TerraformName)
		switch field.Kind {
		case spec.FieldKindString:
			var value types.String
			diagnostics.Append(source.GetAttribute(ctx, attributePath, &value)...)
			values[field.TerraformName] = value
		case spec.FieldKindBool:
			var value types.Bool
			diagnostics.Append(source.GetAttribute(ctx, attributePath, &value)...)
			values[field.TerraformName] = value
		case spec.FieldKindInt64:
			var value types.Int64
			diagnostics.Append(source.GetAttribute(ctx, attributePath, &value)...)
			values[field.TerraformName] = value
		case spec.FieldKindNumber:
			var value types.Number
			diagnostics.Append(source.GetAttribute(ctx, attributePath, &value)...)
			values[field.TerraformName] = value
		case spec.FieldKindObject, spec.FieldKindList:

			var value types.List
			diagnostics.Append(source.GetAttribute(ctx, attributePath, &value)...)
			values[field.TerraformName] = value
		default:
			diagnostics.AddError(
				"Unsupported Field Kind",
				fmt.Sprintf("field %q has kind %q, which the engine cannot read", field.TerraformName, field.Kind),
			)
		}
	}
	return values, diagnostics
}

func buildCreate(fields []spec.FieldSpec, plan map[string]attr.Value, nullables nullableSource) (transport.WireObject, error) {
	object := make(transport.WireObject, len(fields))

	pairs, autoPaired, err := autoAssignmentPairs(fields)
	if err != nil {
		return nil, err
	}
	for _, pair := range pairs {
		if err := createAutoAssigned(pair, plan, nullables, object); err != nil {
			return nil, err
		}
	}

	for _, field := range fields {
		if field.Unmanaged || autoPaired[field.TerraformName] {
			continue
		}

		var value attr.Value
		if field.Nullable {
			written, isWritten := nullables.known(field)
			if !isWritten {

				continue
			}
			value = written
		} else {
			planned, present := plan[field.TerraformName]
			if !present {
				continue
			}
			value = planned
		}

		var (
			wire transport.WireValue
			send bool
			err  error
		)
		if field.Kind == spec.FieldKindObject {
			wire, send, err = createSingleton(field, value, nullables)
		} else if field.Kind == spec.FieldKindList {
			wire, send, err = createList(field, value, nullables)
		} else {
			wire, send, err = createWire(field, value)
		}
		if err != nil {
			return nil, err
		}
		if send {
			object[field.APIName] = wire
		}
	}
	return object, nil
}

func createWire(field spec.FieldSpec, value attr.Value) (transport.WireValue, bool, error) {
	if value.IsUnknown() {
		switch field.UnknownPlan {
		case spec.UnknownPlanOmitAndRead, spec.UnknownPlanPreserve:
			return transport.WireValue{}, false, nil
		case spec.UnknownPlanReject:
			return transport.WireValue{}, false, fmt.Errorf("%s is not known at apply time, and this field cannot be sent unknown", field.TerraformName)
		default:
			return transport.WireValue{}, false, fmt.Errorf("%s declares no unknown_plan policy", field.TerraformName)
		}
	}
	if value.IsNull() {
		switch field.CreateNull {
		case spec.CreateNullOmit:
			return transport.WireValue{}, false, nil
		case spec.CreateNullAPINull:
			return transport.Null(), true, nil
		case spec.CreateNullReject:
			return transport.WireValue{}, false, fmt.Errorf("%s is null, and this field must have a value", field.TerraformName)
		case spec.CreateNullDefault:
			literal, err := literalWire(field)
			if err != nil {
				return transport.WireValue{}, false, err
			}
			return literal, true, nil
		default:
			return transport.WireValue{}, false, fmt.Errorf("%s declares no create_null policy", field.TerraformName)
		}
	}
	wire, err := toWire(field, value)
	if err != nil {
		return transport.WireValue{}, false, err
	}
	return wire, true, nil
}

type nullableSource struct {
	config     map[string]attr.Value
	configured func(terraformName string) bool

	indexed func(block string, index int64, member string) bool

	block func(path string) bool
}

func (n nullableSource) singletonMember(field, member spec.FieldSpec, planMembers map[string]attr.Value) (attr.Value, bool) {
	if n.block == nil || !n.block(field.TerraformName+"."+member.TerraformName) {
		return nil, false
	}
	source := planMembers
	if members, present, err := singletonMembers(field, n.config[field.TerraformName]); err == nil && present {
		source = members
	}
	value, held := source[member.TerraformName]
	return value, held
}

func (n nullableSource) entryMember(field, member spec.FieldSpec, planEntry map[string]attr.Value) (attr.Value, bool) {
	var index int64
	if value, ok := planEntry[field.Collection.IdentityField].(types.Int64); ok {
		index = value.ValueInt64()
	}
	if n.indexed == nil || !n.indexed(field.TerraformName, index, member.TerraformName) {
		return nil, false
	}
	source := planEntry
	if entries, present, err := listEntries(field, n.config[field.TerraformName]); err == nil && present {
		for _, candidate := range entries {
			if value, ok := candidate[field.Collection.IdentityField].(types.Int64); ok && !value.IsNull() && !value.IsUnknown() && value.ValueInt64() == index {
				source = candidate
				break
			}
		}
	}
	value, held := source[member.TerraformName]
	return value, held
}

func (n nullableSource) known(field spec.FieldSpec) (attr.Value, bool) {
	if n.configured == nil || !n.configured(field.TerraformName) {
		return nil, false
	}
	value, held := n.config[field.TerraformName]
	if !held {
		return nil, false
	}
	return value, true
}

func buildUpdate(fields []spec.FieldSpec, plan, state map[string]attr.Value, nullables nullableSource, diagnostics *diag.Diagnostics) (transport.WireObject, bool, error) {
	object := make(transport.WireObject, len(fields))
	changed := false

	companions, paired, err := referencePairs(fields)
	if err != nil {
		return nil, false, err
	}
	for _, field := range fields {
		if field.Unmanaged || field.Reference == nil {
			continue
		}
		pairChanged, err := applyReferencePair(field, companions[field.TerraformName], plan, state, object, diagnostics)
		if err != nil {
			return nil, false, err
		}
		if diagnostics.HasError() {
			return nil, false, nil
		}
		changed = changed || pairChanged
	}

	autoPairs, autoPaired, err := autoAssignmentPairs(fields)
	if err != nil {
		return nil, false, err
	}
	for _, pair := range autoPairs {
		pairChanged, err := updateAutoAssigned(pair, plan, state, nullables.config, object)
		if err != nil {
			return nil, false, err
		}
		changed = changed || pairChanged
	}

	for _, field := range fields {
		if field.Unmanaged || paired[field.TerraformName] || autoPaired[field.TerraformName] {
			continue
		}

		var value attr.Value
		if field.Nullable {
			written, isWritten := nullables.known(field)
			if !isWritten {
				continue
			}
			value = written
		} else {
			planned, present := plan[field.TerraformName]
			if !present {
				continue
			}
			value = planned
		}

		previous, hadPrevious := state[field.TerraformName]
		if hadPrevious && value.Equal(previous) {
			continue
		}
		if field.Kind == spec.FieldKindObject || field.Kind == spec.FieldKindList {
			var (
				wire          transport.WireValue
				objectChanged bool
				err           error
			)
			send := true
			switch {
			case field.Kind == spec.FieldKindList:
				wire, objectChanged, err = updateList(field, value, previous, nullables, diagnostics)
			case !hasManagedMembers(field):
				wire, send, objectChanged, err = updateEmptySingleton(field, value, previous)
			default:
				wire, objectChanged, err = updateSingleton(field, value, previous, nullables, diagnostics)
			}
			if err != nil {
				return nil, false, err
			}
			if diagnostics.HasError() {
				return nil, false, nil
			}
			if objectChanged {
				if send {
					object[field.APIName] = wire
				}
				changed = true
			}
			continue
		}
		wire, send, fieldChanged, err := updateWire(field, value)
		if err != nil {
			return nil, false, err
		}
		changed = changed || fieldChanged
		if send {
			object[field.APIName] = wire
		}
	}
	return object, changed, nil
}

func updateWire(field spec.FieldSpec, value attr.Value) (wire transport.WireValue, send, changed bool, err error) {
	if value.IsUnknown() {

		switch field.UnknownPlan {
		case spec.UnknownPlanOmitAndRead, spec.UnknownPlanPreserve:
			return transport.WireValue{}, false, false, nil
		case spec.UnknownPlanReject:
			return transport.WireValue{}, false, false, fmt.Errorf("%s is not known at apply time, and this field cannot be sent unknown", field.TerraformName)
		default:
			return transport.WireValue{}, false, false, fmt.Errorf("%s declares no unknown_plan policy", field.TerraformName)
		}
	}
	if value.IsNull() {
		cleared, omit, err := clearedWire(field)
		if err != nil {
			return transport.WireValue{}, false, false, err
		}
		return cleared, !omit, true, nil
	}
	converted, err := toWire(field, value)
	if err != nil {
		return transport.WireValue{}, false, false, err
	}
	return converted, true, true, nil
}

func clearedWire(field spec.FieldSpec) (value transport.WireValue, omit bool, err error) {
	switch field.UpdateClear {
	case spec.UpdateClearEmptyString:
		return transport.String(""), false, nil
	case spec.UpdateClearFalse:
		return transport.Bool(false), false, nil
	case spec.UpdateClearZero:
		if field.Kind == spec.FieldKindNumber {
			return transport.Decimal("0"), false, nil
		}
		return transport.Int64(0), false, nil
	case spec.UpdateClearAPINull:
		return transport.Null(), false, nil
	case spec.UpdateClearOmit, spec.UpdateClearOmitUnmanaged:
		return transport.WireValue{}, true, nil
	case spec.UpdateClearDefault:
		literal, err := literalWire(field)
		return literal, false, err
	case spec.UpdateClearReject:
		return transport.WireValue{}, false, fmt.Errorf("%s cannot be cleared once set", field.TerraformName)
	default:
		return transport.WireValue{}, false, fmt.Errorf("%s declares no update_clear policy", field.TerraformName)
	}
}

func literalWire(field spec.FieldSpec) (transport.WireValue, error) {
	if field.Default == nil {
		return transport.WireValue{}, fmt.Errorf("%s asks for its default but declares none", field.TerraformName)
	}
	switch field.Default.Kind {
	case spec.LiteralNull:
		return transport.Null(), nil
	case spec.LiteralString:
		return transport.String(field.Default.Value), nil
	case spec.LiteralBool:
		return transport.Bool(field.Default.Value == "true"), nil
	case spec.LiteralInt64:
		var parsed int64
		if _, err := fmt.Sscanf(field.Default.Value, "%d", &parsed); err != nil {
			return transport.WireValue{}, fmt.Errorf("%s has default %q, which is not an integer", field.TerraformName, field.Default.Value)
		}
		return transport.Int64(parsed), nil
	case spec.LiteralDecimal:
		return transport.Decimal(field.Default.Value), nil
	default:
		return transport.WireValue{}, fmt.Errorf("%s has a default of unknown kind %q", field.TerraformName, field.Default.Kind)
	}
}

func toWire(field spec.FieldSpec, value attr.Value) (transport.WireValue, error) {
	switch field.Kind {
	case spec.FieldKindString:
		typed, ok := value.(types.String)
		if !ok {
			return transport.WireValue{}, fmt.Errorf("%s is declared string but holds %T", field.TerraformName, value)
		}
		return transport.String(typed.ValueString()), nil
	case spec.FieldKindBool:
		typed, ok := value.(types.Bool)
		if !ok {
			return transport.WireValue{}, fmt.Errorf("%s is declared bool but holds %T", field.TerraformName, value)
		}
		return transport.Bool(typed.ValueBool()), nil
	case spec.FieldKindInt64:
		typed, ok := value.(types.Int64)
		if !ok {
			return transport.WireValue{}, fmt.Errorf("%s is declared int64 but holds %T", field.TerraformName, value)
		}
		return transport.Int64(typed.ValueInt64()), nil
	case spec.FieldKindNumber:
		typed, ok := value.(types.Number)
		if !ok {
			return transport.WireValue{}, fmt.Errorf("%s is declared number but holds %T", field.TerraformName, value)
		}
		return transport.Decimal(typed.ValueBigFloat().Text('f', -1)), nil
	default:
		return transport.WireValue{}, fmt.Errorf("%s has kind %q, which the engine cannot send", field.TerraformName, field.Kind)
	}
}

func stateFromAPI(fields []spec.FieldSpec, data map[string]interface{}, mode string, prior map[string]attr.Value) (map[string]attr.Value, error) {
	values := make(map[string]attr.Value, len(fields))
	for _, field := range fields {
		if field.Unmanaged {
			continue
		}
		if !appliesToMode(field, mode) {
			values[field.TerraformName] = nullFor(field)
			continue
		}
		raw, present := data[field.APIName]
		if !present || raw == nil {
			decoded, err := absentValue(field, prior)
			if err != nil {
				return nil, err
			}
			values[field.TerraformName] = decoded
			continue
		}
		if field.Kind == spec.FieldKindObject || field.Kind == spec.FieldKindList {
			decode := singletonFromAPI
			if field.Kind == spec.FieldKindList {
				decode = listFromAPI
			}
			decoded, err := decode(field, raw, mode)
			if err != nil {
				return nil, err
			}
			values[field.TerraformName] = decoded
			continue
		}
		values[field.TerraformName] = fromAPI(field.Kind, raw)
	}
	return values, nil
}

func appliesToMode(field spec.FieldSpec, mode string) bool {
	for _, candidate := range field.Modes {
		if string(candidate) == mode {
			return true
		}
	}
	return false
}

func absentValue(field spec.FieldSpec, prior map[string]attr.Value) (attr.Value, error) {
	switch field.ResponseAbsence {
	case spec.ResponseAbsenceTerraformNull:
		return nullFor(field), nil
	case spec.ResponseAbsencePreserve:
		if previous, held := prior[field.TerraformName]; held && previous != nil {
			return previous, nil
		}
		return nullFor(field), nil
	case spec.ResponseAbsenceDefault:
		if field.Default == nil {
			return nullFor(field), nil
		}
		return literalAttr(field)
	case spec.ResponseAbsenceError:
		return nil, fmt.Errorf("the response carries no %q, and this field must be present", field.APIName)
	default:
		return nil, fmt.Errorf("%s declares no response_absence policy", field.TerraformName)
	}
}

func literalAttr(field spec.FieldSpec) (attr.Value, error) {
	switch field.Default.Kind {
	case spec.LiteralNull:
		return nullFor(field), nil
	case spec.LiteralString:
		return types.StringValue(field.Default.Value), nil
	case spec.LiteralBool:
		return types.BoolValue(field.Default.Value == "true"), nil
	case spec.LiteralInt64:
		var parsed int64
		if _, err := fmt.Sscanf(field.Default.Value, "%d", &parsed); err != nil {
			return nil, fmt.Errorf("%s has default %q, which is not an integer", field.TerraformName, field.Default.Value)
		}
		return types.Int64Value(parsed), nil
	case spec.LiteralDecimal:
		parsed, _, err := big.ParseFloat(field.Default.Value, 10, 512, big.ToNearestEven)
		if err != nil {
			return nil, fmt.Errorf("%s has default %q, which is not a number", field.TerraformName, field.Default.Value)
		}
		return types.NumberValue(parsed), nil
	default:
		return nil, fmt.Errorf("%s has a default of unknown kind %q", field.TerraformName, field.Default.Kind)
	}
}

func fromAPI(kind spec.FieldKind, raw interface{}) attr.Value {
	switch kind {
	case spec.FieldKindString:
		return utils.MapStringFromAPI(raw)
	case spec.FieldKindBool:
		return utils.MapBoolFromAPI(raw)
	case spec.FieldKindInt64:
		return utils.MapInt64FromAPI(raw)
	case spec.FieldKindNumber:
		return utils.MapNumberFromAPI(raw)
	default:
		return nil
	}
}

func nullOf(kind spec.FieldKind) attr.Value {
	switch kind {
	case spec.FieldKindString:
		return types.StringNull()
	case spec.FieldKindBool:
		return types.BoolNull()
	case spec.FieldKindInt64:
		return types.Int64Null()
	case spec.FieldKindNumber:
		return types.NumberNull()
	default:
		return nil
	}
}

func attributeTypes(fields []spec.FieldSpec) map[string]attr.Type {
	types_ := make(map[string]attr.Type, len(fields))
	for _, field := range fields {
		if field.Unmanaged {
			continue
		}
		if null := nullFor(field); null != nil {
			types_[field.TerraformName] = null.Type(context.Background())
		}
	}
	return types_
}
