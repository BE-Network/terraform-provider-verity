package utils_test

import (
	"encoding/json"
	"strconv"
	"testing"

	providerutils "terraform-provider-verity/internal/utils"
	"terraform-provider-verity/openapi"

	"github.com/hashicorp/terraform-plugin-framework/types"
)

const fourByteASN int64 = 4201200024

func TestSetNullableInt64FieldsPreservesFourByteASN(t *testing.T) {
	var neighborASN openapi.NullableInt64

	providerutils.SetNullableInt64Fields([]providerutils.NullableInt64FieldMapping{
		{
			FieldName:    "NeighborAsNumber",
			APIField:     &neighborASN,
			TFValue:      types.Int64Value(fourByteASN),
			IsConfigured: true,
		},
	})

	assertSerializedNeighborASN(t, neighborASN)
}

func TestCompareAndSetNullableInt64FieldPreservesFourByteASN(t *testing.T) {
	var neighborASN openapi.NullableInt64
	hasChanges := false

	providerutils.CompareAndSetNullableInt64Field(
		types.Int64Value(fourByteASN),
		types.Int64Null(),
		true,
		func(value *openapi.NullableInt64) { neighborASN = *value },
		&hasChanges,
	)

	if !hasChanges {
		t.Fatal("expected the field change to be recorded")
	}

	assertSerializedNeighborASN(t, neighborASN)
}

func assertSerializedNeighborASN(t *testing.T, neighborASN openapi.NullableInt64) {
	t.Helper()

	gateway := openapi.GatewaysPutRequestGatewayValue{
		NeighborAsNumber: neighborASN,
	}
	payload, err := json.Marshal(gateway)
	if err != nil {
		t.Fatalf("marshal gateway payload: %v", err)
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(payload, &fields); err != nil {
		t.Fatalf("unmarshal gateway payload: %v", err)
	}

	want := strconv.FormatInt(fourByteASN, 10)
	if got := string(fields["neighbor_as_number"]); got != want {
		t.Fatalf("neighbor_as_number = %s, want %s; payload: %s", got, want, payload)
	}
}
