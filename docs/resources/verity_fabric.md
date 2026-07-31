# Fabric Resource

`verity_fabric` manages Fabric resources in Verity.

## Example Usage

```hcl
resource "verity_fabric" "example" {
  name = "example"
  enable = true
  plane_count = "1"
  su_size = "32"
  server_management = true
  site_type = "enterprise"
  service_for_site = "Management"
  service_for_site_ref_type_ = "service"
  spanning_tree_type = "pvst"
  region_name = ""
  revision = 0
  force_spanning_tree_on_fabric_ports = false
  read_only_mode = false
  domain_for_site = ""
  domain_for_site_ref_type_ = "site_collection"
  enable_dscp = true
  dscp_to_p_bit_map = "0000000011111111222222223333333344444444555555556666666677777777"
  anycast_mac_address_auto_assigned_ = true
  mac_address_aging_time = 600
  mlag_delay_restore_timer = 300
  bgp_keepalive_timer = 60
  bgp_hold_down_timer = 180
  spine_bgp_advertisement_interval = 1
  spine_bgp_connect_timer = 120
  spine_as_number = null
  leaf_bgp_keep_alive_timer = 60
  leaf_bgp_hold_down_timer = 180
  leaf_bgp_advertisement_interval = 1
  leaf_bgp_connect_timer = 120
  link_state_timeout_value = 60
  evpn_multihoming_startup_delay = 300
  evpn_mac_holdtime = 1080
  aggressive_reporting = true
  switch_ip_base = ""
  controller_ip_base = ""
  switch_username = "admin"
  switch_password = "switch-password"
  switch_password_encrypted = ""
  hgx_username = "hgx-admin"
  hgx_password = "hgx-password"
  hgx_password_encrypted = ""
  switch_gateway = "192.168.1.1"
  controller_gateway = "192.168.2.1"
  hgx_gateway = "192.168.3.1"
  gpu_architecture = "hgx"
  multi_tenant = true
  base_bgp_as_number = "61000"
  router_id_base_prefix = "172.16.0.0"
  vtep_id_base_prefix = "172.16.10.0"
  paired_ip_subnet = "192.168.254.0/24"
  max_switches = "2000"
  pause_validation_alarms = false
  starting_octet = null
  max_sus = null
  max_pods = null
  su_support = false
  allow_all_underlay_connections = false
  port_admin_polling_interval = 0
  port_status_polling_interval = 0
  enable_dhcp_snooping = false
  ip_source_guard = false
  duplicate_address_detection_max_number_of_moves = 5
  duplicate_address_detection_time = 180

  object_properties {
    system_graphs {
      index = 1
    }
  }
}
```

## Argument Reference

* `name` (String) - Template Name. Must be unique within type.
* `enable` (Boolean) - Enable object.
* `plane_count` (String) - Number of planes in this Fabric.
* `su_size` (String) - Number of HGXs per SU. Valid values are `32` and `64`; the API default is `32`.
* `su_support` (Boolean) - Support grouping leaf switches in SUs.
* `server_management` (Boolean) - Support managing servers.
* `allow_all_underlay_connections` (Boolean) - Allows underlay connections between PODs.
* `site_type` (String) - Type of Fabric.
* `service_for_site` (String) - Service for Fabric.
* `service_for_site_ref_type_` (String) - Object type for `service_for_site` field.
* `port_admin_polling_interval` (Integer) - Polling interval value in seconds used when aggressive reporting is disabled.
* `port_status_polling_interval` (Integer) - Polling interval value in seconds used when aggressive reporting is disabled.
* `spanning_tree_type` (String) - Sets the spanning tree type for all Ports in this Fabric with Spanning Tree enabled.
* `region_name` (String) - Defines the logical boundary of the network.
* `revision` (Integer) - Revision for the MSTP configuration.
* `force_spanning_tree_on_fabric_ports` (Boolean) - Enable spanning tree on all fabric connections.
* `read_only_mode` (Boolean) - When enabled, vNetC will perform all functions except writing database updates to target hardware.
* `domain_for_site` (String) - Fabric Collection for Fabric.
* `domain_for_site_ref_type_` (String) - Object type for `domain_for_site` field.
* `enable_dscp` (Boolean) - Enable DSCP to p-bit/TC configuration.
* `dscp_to_p_bit_map` (String) - DSCP to TC map string of length 64.
* `anycast_mac_address` (String) - Anycast MAC address to use.
* `anycast_mac_address_auto_assigned_` (Boolean) - Whether the anycast MAC address should be automatically assigned by the API.
* `mac_address_aging_time` (Integer) - MAC Address Aging Time.
* `mlag_delay_restore_timer` (Integer) - MLAG Delay Restore Timer.
* `bgp_keepalive_timer` (Integer) - Spine BGP Keepalive Timer.
* `bgp_hold_down_timer` (Integer) - Spine BGP Hold Down Timer.
* `spine_bgp_advertisement_interval` (Integer) - BGP Advertisement Interval for spines/superspines.
* `spine_bgp_connect_timer` (Integer) - BGP Connect Timer.
* `spine_as_number` (Integer) - BGP AS number applied uniformly to all spine endpoints in this fabric.
* `leaf_bgp_keep_alive_timer` (Integer) - Leaf BGP Keep Alive Timer.
* `leaf_bgp_hold_down_timer` (Integer) - Leaf BGP Hold Down Timer.
* `leaf_bgp_advertisement_interval` (Integer) - BGP Advertisement Interval for leafs.
* `leaf_bgp_connect_timer` (Integer) - BGP Connect Timer.
* `link_state_timeout_value` (Integer) - Link State Timeout Value.
* `evpn_multihoming_startup_delay` (Integer) - Startup Delay.
* `evpn_mac_holdtime` (Integer) - MAC Holdtime.
* `aggressive_reporting` (Boolean) - Fast reporting of switch communications, link up/down, and BGP status.
* `switch_ip_base` (String) - Base IPv4 address for switch IPs in this Fabric.
* `controller_ip_base` (String) - Base IPv4 address for controller IPs in this Fabric.
* `switch_username` (String) - Default username for managed switches in this Fabric.
* `switch_password` (String) - Default password for managed switches in this Fabric.
* `switch_password_encrypted` (String) - Default password for managed switches in this Fabric.
* `hgx_username` (String) - Default username for HGX devices in this Fabric.
* `hgx_password` (String) - Default password for HGX devices in this Fabric.
* `hgx_password_encrypted` (String) - Default password for HGX devices in this Fabric.
* `switch_gateway` (String) - Default switch management gateway IP for devices in this Fabric.
* `controller_gateway` (String) - Default Device Management VM gateway IP for devices in this Fabric.
* `hgx_gateway` (String) - Default HGX management gateway IP for devices in this Fabric.
* `multi_tenant` (Boolean) - Allow multiple tenants to HGX endpoints on this fabric.
* `base_bgp_as_number` (String) - Base BGP Autonomous System Number used for switches in the fabric.
* `router_id_base_prefix` (String) - Router ID starting IP address.
* `vtep_id_base_prefix` (String) - VTEP ID starting IP address.
* `paired_ip_subnet` (String) - IP address range reserved for communication between paired switches.
* `max_switches` (String) - Maximum number of switches to support in this site.
* `pause_validation_alarms` (Boolean) - Validation still runs, but validation alarms are not raised while enabled.
* `starting_octet` (Integer) - Starting Octet for HGX Port IPs.
* `max_sus` (Integer) - Maximum number of SUs allowed per POD.
* `max_pods` (Integer) - Maximum number of PODs allowed in the Fabric.
* `gpu_architecture` (String) - GPU Architecture used within this Fabric.
* `enable_dhcp_snooping` (Boolean) - Enables DHCP snooping.
* `ip_source_guard` (Boolean) - On untrusted ports, only allow known traffic from known IP addresses.
* `duplicate_address_detection_max_number_of_moves` (Integer) - Duplicate Address Detection Max Number of Moves.
* `duplicate_address_detection_time` (Integer) - Duplicate Address Detection Time.
* `object_properties` (Object) - Additional object properties.
  * `system_graphs` (Array) - System graphs.
    * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.

## Import

Fabric resources can be imported using the `name` attribute:

```sh
terraform import verity_fabric.<resource_name> <name>
```
