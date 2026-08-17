package importer

import (
	"context"
	"strings"
	"testing"
)

func TestGeneratedDatacenterStagesPlaceFabricBeforeGateway(t *testing.T) {
	config, err := (&Importer{ctx: context.Background(), Mode: "datacenter"}).generateStagesTF()
	if err != nil {
		t.Fatalf("generateStagesTF returned error: %v", err)
	}

	fabric := strings.Index(config, `resource "verity_operation_stage" "fabric_stage"`)
	gateway := strings.Index(config, `resource "verity_operation_stage" "gateway_stage"`)
	if fabric == -1 || gateway == -1 || fabric >= gateway {
		t.Fatalf("expected Fabric stage before Gateway stage, got:\n%s", config)
	}
	if !strings.Contains(config, "depends_on = [verity_operation_stage.fabric_stage]") {
		t.Fatal("expected the stage after Fabric to depend on fabric_stage")
	}
}

func TestGeneratedCampusStagesExcludeDatacenterOnlyResources(t *testing.T) {
	config, err := (&Importer{ctx: context.Background(), Mode: "campus"}).generateStagesTF()
	if err != nil {
		t.Fatalf("generateStagesTF returned error: %v", err)
	}

	for _, stage := range []string{"ipv4_list_stage", "ipv6_list_stage", "plane_stage", "rack_stage"} {
		if strings.Contains(config, `"`+stage+`"`) {
			t.Errorf("Campus stages must not contain %s", stage)
		}
	}
	if !strings.Contains(config, `resource "verity_operation_stage" "sfp_breakout_stage"`) {
		t.Fatal("Campus stages must contain the SFP Breakout stage")
	}
	if !strings.Contains(config, "depends_on = [verity_operation_stage.sfp_breakout_stage]") {
		t.Fatal("expected the first compatible stage after SFP Breakout to depend on sfp_breakout_stage")
	}
}
