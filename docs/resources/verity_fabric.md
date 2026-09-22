# verity_fabric (Resource)

Manages a Verity Fabric.

Supported modes: Campus, Datacenter.

## Example Usage

### Campus mode

```hcl
resource "verity_fabric" "example" {
  name = "example"
  aggressive_reporting = false
  allow_all_underlay_connections = false
  base_bgp_as_number = ""
  controller_gateway = ""
  controller_ip_base = ""
  domain_for_fabric = ""
  domain_for_fabric_ref_type_ = "site_collection"
  dscp_to_p_bit_map = ""
  duplicate_address_detection_max_number_of_moves = null
  duplicate_address_detection_time = null
  enable = false
  enable_dhcp_snooping = false
  enable_dscp = false
  evpn_mac_holdtime = null
  evpn_multihoming_startup_delay = null
  fabric_type = ""
  force_spanning_tree_on_fabric_ports = false
  gpu_architecture = ""
  hgx_password = ""
  hgx_password_encrypted = ""
  hgx_username = ""
  ip_source_guard = false
  link_state_timeout_value = null
  max_pods = null
  max_sus = null
  max_switches = ""
  multi_tenant = false
  paired_ip_subnet = ""
  pause_validation_alarms = false
  plane_count = ""
  port_admin_polling_interval = null
  port_status_polling_interval = null
  read_only_mode = false
  region_name = ""
  revision = null
  route_aggregation = ""
  router_id_base_prefix = ""
  server_management = false
  service_for_fabric = ""
  service_for_fabric_ref_type_ = "service"
  set_leaf_router_id_on_bgp = false
  spanning_tree_type = ""
  starting_octet = null
  su_size = ""
  su_support = false
  switch_gateway = ""
  switch_ip_base = ""
  switch_password = ""
  switch_password_encrypted = ""
  switch_username = ""
  vtep_id_base_prefix = ""

  object_properties {

    system_graphs {
      index = 1
    }
  }

  route_aggregators {
    index = 1
    route_aggregation_num_enable = false
    route_aggregation_num_ip_and_mask = ""
  }
}
```

### Datacenter mode

```hcl
resource "verity_fabric" "example" {
  name = "example"
  aggressive_reporting = false
  allow_all_underlay_connections = false
  anycast_mac_address_auto_assigned_ = true
  base_bgp_as_number = ""
  bgp_hold_down_timer = null
  bgp_keepalive_timer = null
  controller_gateway = ""
  controller_ip_base = ""
  domain_for_fabric = ""
  domain_for_fabric_ref_type_ = "site_collection"
  dscp_to_p_bit_map = ""
  duplicate_address_detection_max_number_of_moves = null
  duplicate_address_detection_time = null
  enable = false
  enable_dscp = false
  evpn_mac_holdtime = null
  evpn_multihoming_startup_delay = null
  fabric_type = ""
  force_spanning_tree_on_fabric_ports = false
  gpu_architecture = ""
  hgx_password = ""
  hgx_password_encrypted = ""
  hgx_username = ""
  leaf_bgp_advertisement_interval = null
  leaf_bgp_connect_timer = null
  leaf_bgp_hold_down_timer = null
  leaf_bgp_keep_alive_timer = null
  link_state_timeout_value = null
  mac_address_aging_time = null
  max_pods = null
  max_sus = null
  max_switches = ""
  mlag_delay_restore_timer = null
  multi_tenant = false
  paired_ip_subnet = ""
  pause_validation_alarms = false
  plane_count = ""
  port_admin_polling_interval = null
  port_status_polling_interval = null
  read_only_mode = false
  region_name = ""
  revision = null
  route_aggregation = ""
  router_id_base_prefix = ""
  server_management = false
  service_for_fabric = ""
  service_for_fabric_ref_type_ = "service"
  set_leaf_router_id_on_bgp = false
  spanning_tree_type = ""
  spine_as_number = null
  spine_bgp_advertisement_interval = null
  spine_bgp_connect_timer = null
  starting_octet = null
  su_size = ""
  su_support = false
  switch_gateway = ""
  switch_ip_base = ""
  switch_password = ""
  switch_password_encrypted = ""
  switch_username = ""
  vtep_id_base_prefix = ""

  object_properties {

    system_graphs {
      index = 1
    }
  }

  route_aggregators {
    index = 1
    route_aggregation_num_enable = false
    route_aggregation_num_ip_and_mask = ""
  }
}
```

## Argument Reference

### Required

* `name` (String) - Object Name. Must be unique. Changing it replaces the resource.

### Optional

* `aggressive_reporting` (Boolean) - Fast Reporting of Switch Communications, Link Up/Down, and BGP Status.
* `allow_all_underlay_connections` (Boolean) - Allows underlay connections between PODs.
* `anycast_mac_address` (String) - Anycast MAC address to use. This field should not be specified when 'anycast_mac_address_auto_assigned_' is set to true, as the API will assign this value automatically. Used for MAC VRRP. Datacenter mode only. Assigned by the server while `anycast_mac_address_auto_assigned_` is `true`; it cannot be set then.
* `anycast_mac_address_auto_assigned_` (Boolean) - Whether the anycast MAC address should be automatically assigned by the API. When set to true, do not specify the 'anycast_mac_address' field in your configuration. Datacenter mode only.
* `base_bgp_as_number` (String) - Base BGP Autonomous System Number used for switches in the fabric.
* `bgp_hold_down_timer` (Integer) - Spine BGP Hold Down Timer (minimum: 1, maximum: 3600). Datacenter mode only. Set it to `null` to clear it.
* `bgp_keepalive_timer` (Integer) - Spine BGP Keepalive Timer (minimum: 1, maximum: 3600). Datacenter mode only. Set it to `null` to clear it.
* `controller_gateway` (String) - Default Device Management VM gateway IP for devices in this Fabric.
* `controller_ip_base` (String) - Base IPv4 address for controller IPs in this Fabric.
* `domain_for_fabric` (String) - Fabric Collection for Fabric. Set together with `domain_for_fabric_ref_type_`.
* `domain_for_fabric_ref_type_` (String) - Object type for domain_for_fabric field.
* `dscp_to_p_bit_map` (String) - For any Service that is using DSCP to TC map packet prioritization. A string of length 64 with a 0-7 in each position (maxLength: 64).
* `duplicate_address_detection_max_number_of_moves` (Integer) - Duplicate Address Detection Max Number of Moves. Set it to `null` to clear it.
* `duplicate_address_detection_time` (Integer) - Duplicate Address Detection Time. Set it to `null` to clear it.
* `enable` (Boolean) - Enable object.
* `enable_dhcp_snooping` (Boolean) - Enables the switches to monitor DHCP traffic and collect assigned IP addresses which are then placed in the DHCP assigned IPs report. Campus mode only.
* `enable_dscp` (Boolean) - Enable DSCP to p-bit/TC configuration. When enabled, DSCP to p-bit/TC mappings are applied.
* `evpn_mac_holdtime` (Integer) - MAC Holdtime. Set it to `null` to clear it.
* `evpn_multihoming_startup_delay` (Integer) - Startup Delay. Set it to `null` to clear it.
* `fabric_type` (String) - Type of Fabric.
* `force_spanning_tree_on_fabric_ports` (Boolean) - Enable spanning tree on all fabric connections. This overrides the Eth Port Settings for Fabric ports.
* `gpu_architecture` (String) - GPU Architecture used within this Fabric.
* `hgx_password` (String) - Default password for HGX devices in this Fabric.
* `hgx_password_encrypted` (String) - Default password for HGX devices in this Fabric.
* `hgx_username` (String) - Default username for HGX devices in this Fabric.
* `ip_source_guard` (Boolean) - On untrusted ports, only allow known traffic from known IP addresses. IP addresses are discovered via DHCP snooping or with static IP settings. Campus mode only.
* `leaf_bgp_advertisement_interval` (Integer) - BGP Advertisement Interval for leafs. Use "0" for immediate updates (maximum: 3600). Datacenter mode only. Set it to `null` to clear it.
* `leaf_bgp_connect_timer` (Integer) - BGP Connect Timer (minimum: 1, maximum: 3600). Datacenter mode only. Set it to `null` to clear it.
* `leaf_bgp_hold_down_timer` (Integer) - Leaf BGP Hold Down Timer (minimum: 1, maximum: 3600). Datacenter mode only. Set it to `null` to clear it.
* `leaf_bgp_keep_alive_timer` (Integer) - Leaf BGP Keep Alive Timer (minimum: 1, maximum: 3600). Datacenter mode only. Set it to `null` to clear it.
* `link_state_timeout_value` (Integer) - Link State Timeout Value. Set it to `null` to clear it.
* `mac_address_aging_time` (Integer) - MAC Address Aging Time (minimum: 1, maximum: 100000). Datacenter mode only. Set it to `null` to clear it.
* `max_pods` (Integer) - Maximum number of PODs allowed in the Fabric. Set it to `null` to clear it.
* `max_sus` (Integer) - Maximum number of SUs allowed per POD. Set it to `null` to clear it.
* `max_switches` (String) - Max number Switches to support in this site.
* `mlag_delay_restore_timer` (Integer) - MLAG Delay Restore Timer (minimum: 1, maximum: 3600). Datacenter mode only. Set it to `null` to clear it.
* `multi_tenant` (Boolean) - Allow multiple tenants to HGX endpoints on this fabric.
* `object_properties` (Block) - Object properties for the Fabric. At most one block.
  * `system_graphs` (Block List) - System graphs. Entries are matched by `index`.
    * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `paired_ip_subnet` (String) - IP address range reserved for communication between paired switches.
* `pause_validation_alarms` (Boolean) - Validation still runs, but validation alarms are not raised for this Fabric while enabled.
* `plane_count` (String) - Number of planes in this Fabric.
* `port_admin_polling_interval` (Integer) - Polling interval value in seconds used when aggressive reporting is disabled. Set it to `null` to clear it.
* `port_status_polling_interval` (Integer) - Polling interval value in seconds used when aggressive reporting is disabled. Set it to `null` to clear it.
* `read_only_mode` (Boolean) - When Read Only Mode is checked, vNetC will perform all functions except writing database updates to the target hardware.
* `region_name` (String) - Defines the logical boundary of the network. All switches in an MSTP region must have the same configured region name.
* `revision` (Integer) - A logical number that signifies a revision for the MSTP configuration. All switches in an MSTP region must have the same revision number (maximum: 65535). Set it to `null` to clear it.
* `route_aggregation` (String) - Route Aggregation configuration for this fabric.
* `route_aggregators` (Block List) - Route aggregation entries. Entries are matched by `index`.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `route_aggregation_num_enable` (Boolean) - Enable.
  * `route_aggregation_num_ip_and_mask` (String) - IP address and mask for route aggregation.
* `router_id_base_prefix` (String) - Router ID starting IP address.
* `server_management` (Boolean) - Support managing servers.
* `service_for_fabric` (String) - Service for Fabric. Set together with `service_for_fabric_ref_type_`.
* `service_for_fabric_ref_type_` (String) - Object type for service_for_fabric field.
* `set_leaf_router_id_on_bgp` (Boolean) - Use the endpoint loopback0 address as the router ID for all BGP sessions on leaf switches.
* `spanning_tree_type` (String) - Sets the spanning tree type for all Ports in this Fabric with Spanning Tree enabled.
* `spine_as_number` (Integer) - BGP AS number applied uniformly to all spine endpoints in this CLOS fabric on save. Leave blank to manage spine AS numbers individually. Datacenter mode only. Set it to `null` to clear it.
* `spine_bgp_advertisement_interval` (Integer) - BGP Advertisement Interval for spines/superspines. Use "0" for immediate updates (maximum: 3600). Datacenter mode only. Set it to `null` to clear it.
* `spine_bgp_connect_timer` (Integer) - BGP Connect Timer (minimum: 1, maximum: 3600). Datacenter mode only. Set it to `null` to clear it.
* `starting_octet` (Integer) - Starting Octet for HGX Port IPs. Set it to `null` to clear it.
* `su_size` (String) - Number of HGXs per SU. Valid values are "32" and "64".
* `su_support` (Boolean) - Support grouping leaf switches in SUs.
* `switch_gateway` (String) - Default switch management gateway IP for devices in this Fabric.
* `switch_ip_base` (String) - Base IPv4 address for switch IPs in this Fabric.
* `switch_password` (String) - Default password for managed switches in this Fabric.
* `switch_password_encrypted` (String) - Default password for managed switches in this Fabric.
* `switch_username` (String) - Default username for managed switches in this Fabric.
* `vtep_id_base_prefix` (String) - Vtep ID starting IP address.

## Reference Fields

A reference names another Verity object. Set the reference and its type field together. When your configuration sets neither, Terraform leaves the reference to the server, including one set in the Verity UI.

| Field | Type field | Allowed types |
| --- | --- | --- |
| `domain_for_fabric` | `domain_for_fabric_ref_type_` | `site_collection` |
| `service_for_fabric` | `service_for_fabric_ref_type_` | `service` |

## Auto-Assigned Fields

While a flag is `true`, the server assigns the paired value and the value cannot be set in the configuration. Set the flag to `false` to set the value yourself; if you leave the value unset, the value the server assigned is kept.

| Field | Flag | Reassigned when |
| --- | --- | --- |
| `anycast_mac_address` | `anycast_mac_address_auto_assigned_` | — |

## Import

Import an existing object by its `name`:

```sh
terraform import verity_fabric.<resource_name> <name>
```
