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

func TestGeneratedStagesMatchExecutorOrder(t *testing.T) {
	tests := []struct {
		name   string
		mode   string
		stages []string
	}{
		{
			name: "campus",
			mode: "campus",
			stages: []string{
				"sfp_breakout_stage", "acl_v6_stage", "acl_v4_stage", "mac_filter_stage", "service_stage",
				"port_acl_stage", "tacacs_profile_stage", "ldap_profile_stage", "sflow_collector_stage",
				"eth_port_profile_stage", "packet_queue_stage", "device_aaa_profile_stage", "fabric_stage",
				"service_port_profile_stage", "diagnostics_profile_stage", "authenticated_eth_port_stage",
				"device_settings_stage", "voice_port_profile_stage", "lag_stage", "device_voice_setting_stage",
				"eth_port_settings_stage", "diagnostics_port_profile_stage", "bundle_stage", "badge_stage",
				"grouping_rule_stage", "switchpoint_stage", "threshold_stage", "threshold_group_stage", "pair_stage",
			},
		},
		{
			name: "datacenter",
			mode: "datacenter",
			stages: []string{
				"sfp_breakout_stage", "community_list_stage", "as_path_access_list_stage", "ipv6_prefix_list_stage",
				"ipv4_prefix_list_stage", "extended_community_list_stage", "acl_v6_stage", "acl_v4_stage",
				"route_map_clause_stage", "pb_routing_acl_stage", "route_map_stage", "pb_routing_stage", "tenant_stage",
				"service_stage", "fabric_stage", "tacacs_profile_stage", "ldap_profile_stage", "port_acl_stage",
				"ipv6_list_stage", "ipv4_list_stage", "pod_stage", "packet_queue_stage", "device_aaa_profile_stage",
				"eth_port_profile_stage", "packet_broker_stage", "sflow_collector_stage", "gateway_stage", "su_stage",
				"diagnostics_port_profile_stage", "device_settings_stage", "lag_stage", "diagnostics_profile_stage",
				"gateway_profile_stage", "eth_port_settings_stage", "badge_stage", "plane_stage", "spine_plane_stage",
				"rack_stage", "bundle_stage", "ssp_group_stage", "grouping_rule_stage", "switchpoint_stage", "threshold_stage",
				"threshold_group_stage", "pair_stage",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := (&Importer{ctx: context.Background(), Mode: tt.mode}).generateStagesTF()
			if err != nil {
				t.Fatalf("generateStagesTF returned error: %v", err)
			}

			lastPosition := -1
			for _, stage := range tt.stages {
				position := strings.Index(config, `resource "verity_operation_stage" "`+stage+`"`)
				if position == -1 {
					t.Fatalf("missing stage %s", stage)
				}
				if position <= lastPosition {
					t.Fatalf("stage %s is not after its predecessor", stage)
				}
				lastPosition = position
			}
		})
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
