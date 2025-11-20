# Site Resource

Manages a Site object in the Verity system.

**Note:** Only PATCH operations are supported for this resource. Create and delete operations are prohibited.

## Example Usage

```hcl
resource "verity_site" "example" {
    name = "example"
    aggressive_reporting = true
    anycast_mac_address_auto_assigned_ = true
    bgp_hold_down_timer = 180
    bgp_keepalive_timer = 60
    crc_failure_threshold = 5
    dscp_to_p_bit_map = "00000000"
    enable = true
    evpn_mac_holdtime = 1080
    evpn_multihoming_startup_delay = 300
    force_spanning_tree_on_fabric_ports = false
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
    duplicate_address_detection_max_number_of_moves = 5
    duplicate_address_detection_time = 180
    enable_dhcp_snooping = false
    ip_source_guard = false

    pairs {
        name = "pair1"
        switchpoint_1 = "switch1"
        switchpoint_1_ref_type_ = "switchpoint"
        switchpoint_2 = "switch2"
        switchpoint_2_ref_type_ = "switchpoint"
        lag_group = "lag1"
        lag_group_ref_type_ = "lag"
        is_whitebox_pair = false
        index = 1
    }

    islands {
        index = 1
        toi_switchpoint = ""
        toi_switchpoint_ref_type_ = ""
    }

    object_properties {
        system_graphs {
            graph_num_data = "graph_endpoint_path"
            index = 1
        }
    }
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique
* `enable` (Boolean) - Enable object
* `service_for_site` (String) - Service for Site
* `service_for_site_ref_type_` (String) - Object type for service_for_site field
* `spanning_tree_type` (String) - Sets the spanning tree type for all Ports in this Site with Spanning Tree enabled
* `region_name` (String) - Defines the logical boundary of the network. All switches in an MSTP region must have the same configured region name
* `revision` (Integer) - A logical number that signifies a revision for the MSTP configuration. All switches in an MSTP region must have the same revision number
* `force_spanning_tree_on_fabric_ports` (Boolean) - Enable spanning tree on all fabric connections. This overrides the Eth Port Settings for Fabric ports
* `read_only_mode` (Boolean) - When Read Only Mode is checked, vNetC will perform all functions except writing database updates to the target hardware
* `dscp_to_p_bit_map` (String) - For any Service that is using DSCP to p-bit map packet prioritization. A string of length 64 with a 0-7 in each position
* `anycast_mac_address` (String) - Site Level MAC Address for Anycast
* `anycast_mac_address_auto_assigned_` (Boolean) - Whether or not the value in anycast_mac_address field has been automatically assigned or not. Set to false and change anycast_mac_address value to edit.
* `mac_address_aging_time` (Integer) - MAC Address Aging Time (between 1-100000)
* `mlag_delay_restore_timer` (Integer) - MLAG Delay Restore Timer
* `bgp_keepalive_timer` (Integer) - Spine BGP Keepalive Timer
* `bgp_hold_down_timer` (Integer) - Spine BGP Hold Down Timer
* `spine_bgp_advertisement_interval` (Integer) - BGP Advertisement Interval for spines/superspines. Use "0" for immediate updates
* `spine_bgp_connect_timer` (Integer) - BGP Connect Timer
* `leaf_bgp_keep_alive_timer` (Integer) - Leaf BGP Keep Alive Timer
* `leaf_bgp_hold_down_timer` (Integer) - Leaf BGP Hold Down Timer
* `leaf_bgp_advertisement_interval` (Integer) - BGP Advertisement Interval for leafs. Use "0" for immediate updates
* `leaf_bgp_connect_timer` (Integer) - BGP Connect Timer
* `link_state_timeout_value` (Integer) - Link State Timeout Value
* `evpn_multihoming_startup_delay` (Integer) - Startup Delay
* `evpn_mac_holdtime` (Integer) - MAC Holdtime
* `aggressive_reporting` (Boolean) - Fast Reporting of Switch Communications, Link Up/Down, and BGP Status
* `crc_failure_threshold` (Integer) - Threshold in Errors per second that when met will disable the links as part of LAGs
* `duplicate_address_detection_max_number_of_moves` (Integer) - Controls duplicate MAC address detection (DAD) Max Number of Moves for EVPN (Ethernet VPN) within the BGP address-family. Number of moves (2 to 1000; default 5 if left blank)
* `duplicate_address_detection_time` (Integer) - Controls duplicate MAC address detection (DAD) time for EVPN (Ethernet VPN) within the BGP address-family. Time in seconds (2 to 1800; default 180 if left blank)
* `enable_dhcp_snooping` (Boolean) - Enables the switches to monitor DHCP traffic and collect assigned IP addresses which are then placed in the DHCP assigned IPs report.
* `ip_source_guard` (Boolean) - On untrusted ports, only allow known traffic from known IP addresses. IP addresses are discovered via DHCP snooping or with static IP settings
* `islands` (Array) - Defines islands for the site
  * `toi_switchpoint` (String) - TOI Switchpoint
  * `toi_switchpoint_ref_type_` (String) - Object type for toi_switchpoint field
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `pairs` (Array) - Defines pairs for the site
  * `name` (String) - Object Name. Must be unique.
  * `switchpoint_1` (String) - Switchpoint
  * `switchpoint_1_ref_type_` (String) - Object type for switchpoint_1 field
  * `switchpoint_2` (String) - Switchpoint
  * `switchpoint_2_ref_type_` (String) - Object type for switchpoint_2 field
  * `lag_group` (String) - LAG Group
  * `lag_group_ref_type_` (String) - Object type for lag_group field
  * `is_whitebox_pair` (Boolean) - LAG Pair
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `object_properties` (Object) - Additional object properties
  * `system_graphs` (Array) - Graph data for the site
    * `graph_num_data` (String) - The graph data detailing this graph choice
    * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.

## Import

This resource can be imported using the object name:

```sh
terraform import verity_site.<resource_name> <name>
```
