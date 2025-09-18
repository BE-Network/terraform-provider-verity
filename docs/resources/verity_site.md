# Site Resource

Manages a Site object in the Verity system.

**Note:** Only PATCH operations are supported for this resource. Create and delete operations are prohibited.

## Example Usage

```hcl
resource "verity_site" "TOR_Complex" {
    name = "TOR Complex"
    depends_on = [verity_operation_stage.site_stage]
    object_properties {
        system_graphs {
            graph_num_data = "graph_endpoint_path=::graph_group_name=::graph_lagg_path=::graph_line_set=::graph_processes=::graph_source_mod=::graph_source_path=::graph_sub_path=::graph_time_frame="
            index = 1
        }
    }
    aggressive_reporting = true
    anycast_mac_address_auto_assigned_ = true
    bgp_hold_down_timer = 180
    bgp_keepalive_timer = 60
    crc_failure_threshold = 5
    dscp_to_p_bit_map = "0000000011111111222222223333333344444444555555556666666677777777"
    enable = true
    evpn_mac_holdtime = 1080
    evpn_multihoming_startup_delay = 300
    force_spanning_tree_on_fabric_ports = false
    islands {
        index = 1
        toi_switchpoint = ""
        toi_switchpoint_ref_type_ = ""
    }
    leaf_bgp_advertisement_interval = 1
    leaf_bgp_connect_timer = 120
    leaf_bgp_hold_down_timer = 180
    leaf_bgp_keep_alive_timer = 60
    link_state_timeout_value = 60
    mac_address_aging_time = 600
    mlag_delay_restore_timer = 300
    read_only_mode = false
    region_name = ""
    revision = 0
    service_for_site = "Management"
    service_for_site_ref_type_ = "service"
    spanning_tree_type = "pvst"
    spine_bgp_advertisement_interval = 1
    spine_bgp_connect_timer = 120
}
```

## Argument Reference

The following arguments are supported:

- `name` (String, Required) — Object Name. Must be unique.
- `enable` (Boolean, Optional) — Enable object. Default: `true`.
- `service_for_site` (String, Optional) — Service for Site. Default: `service|Management|`.
- `service_for_site_ref_type_` (String, Optional) — Object type for service_for_site field. Allowed value: `service`.
- `spanning_tree_type` (String, Optional) — Sets the spanning tree type for all Ports in this Site with Spanning Tree enabled. Default: `pvst`. Allowed values: `"", "port", "pvst", "mstp"`.
- `region_name` (String, Optional) — Defines the logical boundary of the network. All switches in an MSTP region must have the same configured region name.
- `revision` (Integer, Optional) — Revision for MSTP configuration. Default: `0`. Maximum: `65535`.
- `force_spanning_tree_on_fabric_ports` (Boolean, Optional) — Enable spanning tree on all fabric connections. Default: `false`.
- `read_only_mode` (Boolean, Optional) — When enabled, vNetC performs all functions except writing database updates to hardware. Default: `false`.
- `dscp_to_p_bit_map` (String, Optional) — DSCP to p-bit map packet prioritization. Default: `0000000011111111222222223333333344444444555555556666666677777777`.
- `anycast_mac_address` (String, Optional) — Site Level MAC Address for Anycast. Default: `(auto)`.
- `anycast_mac_address_auto_assigned_` (Boolean, Optional) — Whether the anycast_mac_address is auto-assigned.
- `mac_address_aging_time` (Integer, Optional) — MAC Address Aging Time. Default: `600`. Range: `1-100000`.
- `mlag_delay_restore_timer` (Integer, Optional) — MLAG Delay Restore Timer. Default: `300`. Range: `1-3600`.
- `bgp_keepalive_timer` (Integer, Optional) — Spine BGP Keepalive Timer. Default: `60`. Range: `1-3600`.
- `bgp_hold_down_timer` (Integer, Optional) — Spine BGP Hold Down Timer. Default: `180`. Range: `1-3600`.
- `spine_bgp_advertisement_interval` (Integer, Optional) — BGP Advertisement Interval for spines/superspines. Default: `1`. Maximum: `3600`.
- `spine_bgp_connect_timer` (Integer, Optional) — BGP Connect Timer. Default: `120`. Range: `1-3600`.
- `leaf_bgp_keep_alive_timer` (Integer, Optional) — Leaf BGP Keep Alive Timer. Default: `60`. Range: `1-3600`.
- `leaf_bgp_hold_down_timer` (Integer, Optional) — Leaf BGP Hold Down Timer. Default: `180`. Range: `1-3600`.
- `leaf_bgp_advertisement_interval` (Integer, Optional) — BGP Advertisement Interval for leafs. Default: `1`. Maximum: `3600`.
- `leaf_bgp_connect_timer` (Integer, Optional) — BGP Connect Timer. Default: `120`. Range: `1-3600`.
- `link_state_timeout_value` (Integer, Optional) — Link State Timeout Value. Default: `60`.
- `evpn_multihoming_startup_delay` (Integer, Optional) — Startup Delay. Default: `300`.
- `evpn_mac_holdtime` (Integer, Optional) — MAC Holdtime. Default: `1080`.
- `aggressive_reporting` (Boolean, Optional) — Fast Reporting of Switch Communications, Link Up/Down, and BGP Status. Default: `true`.
- `crc_failure_threshold` (Integer, Optional) — Threshold in Errors per second to disable links as part of LAGs. Default: `5`.
- `islands` (Block, Optional, repeatable) — Defines islands for the site. Each block supports:
    - `toi_switchpoint` (String, Optional) — TOI Switchpoint.
    - `toi_switchpoint_ref_type_` (String, Optional) — Object type for toi_switchpoint field. Allowed value: `switchpoint`.
    - `index` (Integer, Optional) — Index identifying the object.
- `pairs` (Block, Optional, repeatable) — Defines pairs for the site. Each block supports:
    - `name` (String, Required) — Object Name. Must be unique.
    - `switchpoint_1` (String, Optional) — Switchpoint.
    - `switchpoint_1_ref_type_` (String, Optional) — Object type for switchpoint_1 field. Allowed value: `switchpoint`.
    - `switchpoint_2` (String, Optional) — Switchpoint.
    - `switchpoint_2_ref_type_` (String, Optional) — Object type for switchpoint_2 field. Allowed value: `switchpoint`.
    - `lag_group` (String, Optional) — LAG Group.
    - `lag_group_ref_type_` (String, Optional) — Object type for lag_group field. Allowed value: `lag`.
    - `is_whitebox_pair` (Boolean, Optional) — LAG Pair. Default: `false`.
    - `index` (Integer, Optional) — Index identifying the object.
- `object_properties` (Block, Optional) — Additional object properties. Supports:
    - `system_graphs` (Block, Optional, repeatable) — Graph data for the site. Each block supports:
        - `graph_num_data` (String, Optional) — The graph data detailing this graph choice.
        - `index` (Integer, Optional) — Index identifying the object.
- `enable_dhcp_snooping` (Boolean, Optional) — Enables DHCP snooping. Default: `false`.
- `ip_source_guard` (Boolean, Optional) — Enables IP source guard. Default: `false`.

## Import

This resource can be imported using the object name:

```sh
terraform import verity_site.<resource_name> <name>
```
