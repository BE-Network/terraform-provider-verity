package transport

import (
	"encoding/json"
	"fmt"
)

type ValueKind string

const (
	ValueKindNull    ValueKind = "null"
	ValueKindString  ValueKind = "string"
	ValueKindBool    ValueKind = "bool"
	ValueKindInt64   ValueKind = "int64"
	ValueKindDecimal ValueKind = "decimal"
	ValueKindObject  ValueKind = "object"
	ValueKindList    ValueKind = "list"
)

type WireValue struct {
	kind    ValueKind
	string  string
	boolean bool
	int64   int64
	decimal json.Number
	object  WireObject
	list    []WireValue
}

type WireObject map[string]WireValue

func Null() WireValue               { return WireValue{kind: ValueKindNull} }
func String(value string) WireValue { return WireValue{kind: ValueKindString, string: value} }
func Bool(value bool) WireValue     { return WireValue{kind: ValueKindBool, boolean: value} }
func Int64(value int64) WireValue   { return WireValue{kind: ValueKindInt64, int64: value} }
func Decimal(value string) WireValue {
	return WireValue{kind: ValueKindDecimal, decimal: json.Number(value)}
}
func Object(value WireObject) WireValue     { return WireValue{kind: ValueKindObject, object: value} }
func List(value []WireValue) WireValue      { return WireValue{kind: ValueKindList, list: value} }
func (v WireValue) Kind() ValueKind         { return v.kind }
func (v WireValue) StringValue() string     { return v.string }
func (v WireValue) BoolValue() bool         { return v.boolean }
func (v WireValue) Int64Value() int64       { return v.int64 }
func (v WireValue) DecimalValue() string    { return v.decimal.String() }
func (v WireValue) ObjectValue() WireObject { return v.object }
func (v WireValue) ListValue() []WireValue  { return append([]WireValue(nil), v.list...) }

func (v WireValue) MarshalJSON() ([]byte, error) {
	switch v.kind {
	case ValueKindNull:
		return []byte("null"), nil
	case ValueKindString:
		return json.Marshal(v.string)
	case ValueKindBool:
		return json.Marshal(v.boolean)
	case ValueKindInt64:
		return json.Marshal(v.int64)
	case ValueKindDecimal:
		return json.Marshal(v.decimal)
	case ValueKindObject:
		return json.Marshal(v.object)
	case ValueKindList:
		return json.Marshal(v.list)
	default:
		return nil, fmt.Errorf("unsupported wire value kind %q", v.kind)
	}
}
