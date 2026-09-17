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

// attributeSource is the part of tfsdk.Plan, tfsdk.State and tfsdk.Config this
// engine uses. Taking the narrow interface rather than the concrete types lets
// one reader serve all three.
type attributeSource interface {
	GetAttribute(ctx context.Context, p path.Path, target interface{}) diag.Diagnostics
}

// readScalars pulls every spec field out of a plan or state into a map keyed by
// Terraform name. The value keeps its null and unknown state, because the
// lifecycle policies are decided on exactly that distinction.
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
		case spec.FieldKindObject:
			// A singleton is stored as a list block; its entry is read later.
			var value types.List
			diagnostics.Append(source.GetAttribute(ctx, attributePath, &value)...)
			values[field.TerraformName] = value
		default:
			diagnostics.AddError(
				"Unsupported Field Kind",
				fmt.Sprintf("field %q has kind %q, which the scalar engine cannot read", field.TerraformName, field.Kind),
			)
		}
	}
	return values, diagnostics
}

// buildCreate turns a plan into the object a create sends.
//
// Each field is decided by its own declared policy and nothing else: an unknown
// by UnknownPlan, a null by CreateNull, and anything else is sent as it stands.
// A field the request omits is read back afterwards, which is what makes
// omission safe for a Computed attribute.
func buildCreate(fields []spec.FieldSpec, plan map[string]attr.Value, nullables nullableSource) (transport.WireObject, error) {
	object := make(transport.WireObject, len(fields))

	// An auto-assignment pair is decided together: the flag decides whether the
	// value is sent at all, so neither half follows its own policy alone.
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

		// The source of the value differs for a nullable field; what is done with
		// the value does not. On create a nullable attribute is Computed with no
		// prior state, so `x = null` plans as unknown and only the configuration
		// scan shows what was written. The scan answers whether the field was
		// written and what it holds, and the declared policies decide the rest.
		var value attr.Value
		if field.Nullable {
			written, isWritten := nullables.known(field)
			if !isWritten {
				// Not written: the server keeps its value.
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
			wire, send, err = createSingleton(field, value)
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

// createWire applies one field's declared create policies to its value: an
// unknown by UnknownPlan, a null by CreateNull, and anything else sent as it
// stands. It reports whether the field belongs in the request at all.
//
// Every field goes through here, whatever supplied its value, so a policy the
// registry declares is honored uniformly rather than assumed for a class of
// field.
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

// buildUpdate turns a plan and the state it replaces into the object an update
// sends, along with whether anything changed at all. An update carries only the
// fields that differ, so a field equal to its state is absent from the result.
//
// Clearing is the case that needs the spec: what "no value" looks like on the
// wire is per field, and the engine must not guess it from the Go type.
// nullableSource carries what only the configuration can answer: whether a
// nullable attribute is written at all, and what it was written as.
//
// A nullable numeric is the one field that can be cleared by an explicit null,
// and `x = null` is indistinguishable from an absent x in a plan for an
// Optional and Computed attribute — both arrive as the value already in state.
// The handwritten resources resolve it by reading the .tf files, and so does
// this, through the same parser.
type nullableSource struct {
	config     map[string]attr.Value
	configured func(terraformName string) bool
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

	// Reference pairs are decided together and then skipped by the loop below,
	// because neither half's policy describes what the API requires of the two.
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

		// As on create, a nullable field differs only in where its value comes
		// from: the configuration scan, because only configuration distinguishes
		// "cleared" from "not mentioned". A field that is not written is left
		// alone; one that is written goes through the declared policies.
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
		if field.Kind == spec.FieldKindObject {
			wire, objectChanged, err := updateSingleton(field, value, previous, diagnostics)
			if err != nil {
				return nil, false, err
			}
			if diagnostics.HasError() {
				return nil, false, nil
			}
			if objectChanged {
				object[field.APIName] = wire
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

// updateWire applies one field's declared update policies to a value that
// differs from state. It reports whether the field belongs in the request and
// whether the update counts as a change: a clear by omission is a change with
// nothing to send, and an unknown left to the read is neither.
func updateWire(field spec.FieldSpec, value attr.Value) (wire transport.WireValue, send, changed bool, err error) {
	if value.IsUnknown() {
		// The request leaves it out and the read that follows supplies it.
		// Leaving it out is not a change to send, so it does not on its own make
		// the update worth making.
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

// clearedWire says what clearing one field looks like, and whether clearing it
// means sending nothing at all. Omission is a real strategy here, not an absence
// of one: a member of an object sent whole is cleared by dropping its key.
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

// toWire converts one known, non-null Framework value. The kind comes from the
// spec, so a value whose Go type disagrees with it is a registry error rather
// than something to coerce quietly.
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
		return transport.WireValue{}, fmt.Errorf("%s has kind %q, which the scalar engine cannot send", field.TerraformName, field.Kind)
	}
}

// stateFromAPI decodes a response object into the values state will hold.
//
// A field that does not apply to the running mode reads as null whatever the
// response says, which is what keeps a campus-only field from appearing in a
// datacenter plan. Otherwise ResponseAbsence decides what an absent member
// means.
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
		if field.Kind == spec.FieldKindObject {
			decoded, err := singletonFromAPI(field, raw, mode)
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

// appliesToMode reports whether a field is part of the API surface the running
// mode exposes.
//
// It reads the spec's own Modes rather than utils.FieldAppliesToMode, which
// answers the same question from the generated lookup table. The table is
// derived from this registry, so the data agrees, but it is reached by two
// strings and fails open on any miss: an unknown endpoint, an unknown field, or
// an unrecognised mode all return true. A field this engine holds the spec for
// should not be decided by a lookup that can silently answer "yes" because it
// found nothing.
//
// Reading Modes fails closed instead, which is safe because it cannot be empty:
// validateModes requires at least one mode and rejects any value that is not
// datacenter or campus, and the registry is validated before the engine sees it.
func appliesToMode(field spec.FieldSpec, mode string) bool {
	for _, candidate := range field.Modes {
		if string(candidate) == mode {
			return true
		}
	}
	return false
}

// absentValue applies ResponseAbsence to a member the response did not carry.
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

// fromAPI decodes one present member. It delegates to the same helpers the
// handwritten resources use, so a value that decoded one way before decodes the
// same way now — including the JSON-number handling an int64 needs.
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

// attributeTypes describes the object state holds, which the engine needs to
// build a whole state value rather than setting attributes one at a time.
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
