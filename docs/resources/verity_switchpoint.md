# Switchpoint Resource

`verity_switchpoint` manages switchpoint resources in Verity, which represent network devices like switches or routers.

## Example Usage

```hcl
resource "verity_switchpoint" "example" {
  name = "example"
  enable = true
  tenant = "tenant-a"
  tenant_ref_type_ = "tenant"
  device_serial_number = "ABC123456"
  connected_bundle = "main_bundle"
  connected_bundle_ref_type_ = "endpoint_bundle"
  is_top_of_island = false
  read_only_mode = false
  locked = false
  expected_site = "dc-fabric"
  expected_site_ref_type_ = "site"
  out_of_band_management = true
  type = "leaf"
  spine_plane = "plane1"
  spine_plane_ref_type_ = "spine_plane"
  pod = "rack1"
  pod_ref_type_ = "pod"
  su = "su-a"
  su_ref_type_ = "su"
  ssp_group = "ssp-a"
  ssp_group_ref_type_ = "superspine_group"
  rack = "A1"
  position = 0
  rail_group = 0
  switch_router_id_ip_mask = "192.168.1.1"
  switch_router_id_ip_mask_auto_assigned_ = false
  switch_vtep_id_ip_mask = "192.168.2.1"
  switch_vtep_id_ip_mask_auto_assigned_ = false
  bgp_as_number = 65001
  bgp_as_number_auto_assigned_ = false
  is_fabric = false
  
  badges {
    badge = "badge1"
    badge_ref_type_ = "badge"
    index = 0
  }
  
  children {
    child_num_endpoint = "child_endpoint"
    child_num_endpoint_ref_type_ = "switchpoint"
    child_num_device = "child_device"
    index = 0
  }
  
  traffic_mirrors {
    traffic_mirror_num_enable = true
    traffic_mirror_num_source_port = "eth1"
    traffic_mirror_num_source_lag_indicator = false
    traffic_mirror_num_destination_port = "eth2"
    traffic_mirror_num_inbound_traffic = true
    traffic_mirror_num_outbound_traffic = true
    index = 0
  }
  
  eths {
    breakout = "1x100G"
    customer_vlan = ""
    index = 0
    eth_num_icon = "server"
    eth_num_label = "Server Port"
    enable = true
    port_name = "eth0"
  }

  pots {
    index = 0
    pots_num_enable = false
    pots_num_uri = ""
    pots_num_username = ""
    pots_num_password = ""
    pots_num_caller_id = ""
    pots_num_hot_line = ""
    pots_num_password_encrypted = ""
  }
  
  object_properties {
    user_notes = "Important device"
    expected_parent_endpoint = "parent_endpoint"
    expected_parent_endpoint_ref_type_ = "switchpoint"
    number_of_multipoints = 4
    aggregate = true
    is_host = false
    emulate_rf_video_port = false
    draw_as_edge_device = false
  }
}
```

## Argument Reference

* `name` (String) - Unique identifier for the switchpoint
* `enable` (Boolean) - Enable object
* `tenant` (String) - The Tenant of this Device
* `tenant_ref_type_` (String) - Object type for tenant field
* `device_serial_number` (String) - Serial number of the physical device
* `connected_bundle` (String) - Reference to a bundle resource
* `connected_bundle_ref_type_` (String) - Object type for bundle reference
* `is_top_of_island` (Boolean) - Mark this Switchpoint as Top of Island
* `read_only_mode` (Boolean) - Whether the device is in read-only mode
* `locked` (Boolean) - Whether the device is locked
* `expected_site` (String) - Expected Fabric
* `expected_site_ref_type_` (String) - Object type for expected_site field
* `out_of_band_management` (Boolean) - Whether out-of-band management is enabled
* `type` (String) - Type of switchpoint
* `spine_plane` (String) - Spine Plane - subgrouping of super spine and spine
* `spine_plane_ref_type_` (String) - Object type for spine_plane field
* `pod` (String) - Pod - subgrouping of spine and leaf switches
* `pod_ref_type_` (String) - Object type for pod field
* `su` (String) - SU
* `su_ref_type_` (String) - Object type for su field
* `ssp_group` (String) - SuperSpine Group - grouping of superspines in 3-tier config
* `ssp_group_ref_type_` (String) - Object type for ssp_group field
* `rack` (String) - Physical Rack location of the Switch
* `position` (Number) - Position of the Switch
* `rail_group` (Number) - Rail Group the Switch is part of
* `switch_router_id_ip_mask` (String) - Switch BGP Router Identifier
* `switch_router_id_ip_mask_auto_assigned_` (Boolean) - Whether router ID is auto-assigned
* `switch_vtep_id_ip_mask` (String) - Switch VETP Identifier
* `switch_vtep_id_ip_mask_auto_assigned_` (Boolean) - Whether VTEP ID is auto-assigned
* `bgp_as_number` (Integer) - BGP Autonomous System Number for the site underlay
* `bgp_as_number_auto_assigned_` (Boolean) - Whether BGP AS number is auto-assigned
* `is_fabric` (Boolean) - Whether this is a fabric switch
* `bb_switch` (Boolean) - Expose fields for Device Management
* `ip_source` (String) - IP Source
* `controller_ip_and_mask` (String) - Controller IP and Mask
* `gateway` (String) - Gateway
* `switch_ip_and_mask` (String) - Switch IP and Mask
* `switch_gateway` (String) - Gateway of Managed Device
* `comm_type` (String) - Comm Type
* `snmp_community_string` (String) - Comm Credentials
* `uplink_port` (String) - Uplink Port of Managed Device
* `lldp_search_string` (String) - LLDP search string used to detect managed-device connections
* `ztp_identification` (String) - Service Tag or Serial Number for Zero Touch Provisioning
* `located_by` (String) - Controls how the system locates this Device within its LAN
* `power_state` (String) - Power state of Switch Controller
* `communication_mode` (String) - Communication Mode
* `cli_access_mode` (String) - CLI Access Mode
* `username` (String) - Username
* `password` (String) - Password
* `enable_password` (String) - Enable Password
* `ssh_key_or_password` (String) - SSH Key or Password
* `managed_on_native_vlan` (Boolean) - Managed on native VLAN
* `sdlc` (String) - SDLC that Device Controller belongs to
* `security_type` (String) - Security level
* `snmpv3_username` (String) - SNMPv3 username
* `authentication_protocol` (String) - Authentication protocol
* `passphrase` (String) - Passphrase
* `private_protocol` (String) - Privacy protocol
* `private_password` (String) - Privacy password
* `password_encrypted` (String) - Encrypted password
* `enable_password_encrypted` (String) - Encrypted enable password
* `ssh_key_or_password_encrypted` (String) - Encrypted SSH key or password
* `passphrase_encrypted` (String) - Encrypted passphrase
* `private_password_encrypted` (String) - Encrypted privacy password
* `device_managed_as` (String) - Device managed as
* `switch` (String) - Switchpoint locating the Switch to be controlled
* `switch_ref_type_` (String) - Object type for switch field
* `connection_service` (String) - Connect a Service
* `connection_service_ref_type_` (String) - Object type for connection_service field
* `port` (String) - Port locating the Switch to be controlled
* `uses_tagged_packets` (Boolean) - Indicates if the direct interface expects tagged or untagged packets
* `badges` (Array) - List of badge configurations
  * `badge` (String) - Badge reference
  * `badge_ref_type_` (String) - Object type for badge reference
  * `index` (Integer) - Index identifying this badge
* `children` (Array) - List of child device configurations
  * `child_num_endpoint` (String) - Child endpoint reference
  * `child_num_endpoint_ref_type_` (String) - Object type for endpoint reference
  * `child_num_device` (String) - Child device reference
  * `index` (Integer) - Index identifying this child
* `traffic_mirrors` (Array) - List of traffic mirroring configurations
  * `traffic_mirror_num_enable` (Boolean) - Enable this traffic mirror
  * `traffic_mirror_num_source_port` (String) - Source port for mirrored traffic
  * `traffic_mirror_num_source_lag_indicator` (Boolean) - Whether source is a LAG
  * `traffic_mirror_num_destination_port` (String) - Destination port for mirrored traffic
  * `traffic_mirror_num_inbound_traffic` (Boolean) - Mirror inbound traffic
  * `traffic_mirror_num_outbound_traffic` (Boolean) - Mirror outbound traffic
  * `index` (Integer) - Index identifying this mirror
* `eths` (Array) - List of Ethernet port configurations
  * `breakout` (String) - Breakout Port Override
  * `customer_vlan` (String) - A Value between 1 and 4096
  * `index` (Integer) - Index identifying this port
  * `eth_num_icon` (String) - Icon of this Eth Port
  * `eth_num_label` (String) - Label of this Eth Port
  * `enable` (Boolean) - Enable port
  * `port_name` (String) - The name identifying the port
* `pots` (Array) - List of POTS configurations
  * `pots_num_enable` (Boolean) - Enable POTS port
  * `pots_num_uri` (String) - Specific telephone extension for SIP for POTS port
  * `pots_num_username` (String) - SIP username used for authentication for POTS port
  * `pots_num_password` (String) - SIP password used for authentication for POTS port
  * `pots_num_caller_id` (String) - ASCII string defining the user for Caller ID display for POTS port
  * `pots_num_hot_line` (String) - URI of line to autodial upon off-hook for POTS port
  * `pots_num_password_encrypted` (String) - SIP password used for authentication for POTS port
  * `index` (Integer) - Index identifying this POTS entry
* `object_properties` (Object) - Object properties configuration
  * `user_notes` (String) - Notes written by User about the site
  * `expected_parent_endpoint` (String) - Expected Parent Endpoint
  * `expected_parent_endpoint_ref_type_` (String) - Object type for expected_parent_endpoint field
  * `number_of_multipoints` (Integer) - Number of Multipoints
  * `aggregate` (Boolean) - Whether this is an aggregate device
  * `is_host` (Boolean) - Whether this device is a host
  * `emulate_rf_video_port` (Boolean) - Emulate RF Video Port
  * `draw_as_edge_device` (Boolean) - Turn on to display the switch as an edge device instead of as a switch

## Import

Switchpoint resources can be imported using the `name` attribute:

```sh
terraform import verity_switchpoint.<resource_name> <name>
```
