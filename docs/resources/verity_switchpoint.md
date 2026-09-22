# verity_switchpoint (Resource)

Manages a Verity Switchpoint.

Supported modes: Campus, Datacenter.

## Example Usage

### Campus mode

```hcl
resource "verity_switchpoint" "example" {
  name = "example"
  authentication_protocol = ""
  bb_switch = false
  cli_access_mode = ""
  comm_type = ""
  communication_mode = ""
  connected_bundle = ""
  connected_bundle_ref_type_ = "endpoint_bundle"
  connection_service = ""
  connection_service_ref_type_ = "service"
  controller_ip_and_mask_auto_assigned_ = true
  device_managed_as = ""
  device_serial_number = ""
  enable = false
  enable_password = ""
  enable_password_encrypted = ""
  expected_breakout = ""
  expected_breakout_uplink_port = ""
  expected_fabric = ""
  expected_fabric_ref_type_ = "fabric"
  expected_uplink_port = null
  gateway_auto_assigned_ = true
  ip_source = ""
  is_fabric = false
  is_top_of_island = false
  lldp_search_string_auto_assigned_ = true
  located_by = ""
  locked = false
  managed_on_native_vlan = false
  out_of_band_management = false
  passphrase = ""
  passphrase_encrypted = ""
  password = ""
  password_encrypted = ""
  port = ""
  power_state = ""
  private_password = ""
  private_password_encrypted = ""
  private_protocol = ""
  read_only_mode = false
  sdlc = ""
  security_type = ""
  snmp_community_string = ""
  snmpv3_username = ""
  ssh_key_or_password = ""
  ssh_key_or_password_encrypted_auto_assigned_ = true
  switch = ""
  switch_gateway_auto_assigned_ = true
  switch_ip_and_mask_auto_assigned_ = true
  switch_ref_type_ = "switchpoint"
  tenant = ""
  tenant_ref_type_ = "tenant"
  uplink_port = ""
  upstream_is_lag = false
  username_auto_assigned_ = true
  uses_tagged_packets = false
  ztp_identification = ""

  badges {
    index = 1
    badge = ""
    badge_ref_type_ = "badge"
  }

  children {
    index = 1
    child_num_device = ""
    child_num_endpoint = ""
    child_num_endpoint_ref_type_ = "switchpoint"
  }

  eths {
    index = 1
    breakout = ""
    customer_vlan = ""
    enable = false
    eth_num_icon = ""
    eth_num_label = ""
    port_name = ""
  }

  object_properties {
    aggregate = false
    draw_as_edge_device = false
    emulate_rf_video_port = false
    expected_parent_endpoint = ""
    expected_parent_endpoint_ref_type_ = "switchpoint"
    is_host = false
    number_of_multipoints = null
    user_notes = ""
  }

  pots {
    index = 1
    pots_num_caller_id = ""
    pots_num_enable = false
    pots_num_hot_line = ""
    pots_num_password = ""
    pots_num_password_encrypted = ""
    pots_num_uri = ""
    pots_num_username = ""
  }

  traffic_mirrors {
    index = 1
    traffic_mirror_num_destination_port = ""
    traffic_mirror_num_enable = false
    traffic_mirror_num_inbound_traffic = false
    traffic_mirror_num_outbound_traffic = false
    traffic_mirror_num_source_lag_indicator = false
    traffic_mirror_num_source_port = ""
  }
}
```

### Datacenter mode

```hcl
resource "verity_switchpoint" "example" {
  name = "example"
  authentication_protocol = ""
  bb_switch = false
  bgp_as_number_auto_assigned_ = true
  cli_access_mode = ""
  comm_type = ""
  communication_mode = ""
  connected_bundle = ""
  connected_bundle_ref_type_ = "endpoint_bundle"
  controller_ip_and_mask_auto_assigned_ = true
  device_serial_number = ""
  enable = false
  enable_password = ""
  enable_password_encrypted = ""
  expected_breakout = ""
  expected_breakout_uplink_port = ""
  expected_fabric = ""
  expected_fabric_ref_type_ = "fabric"
  expected_uplink_port = null
  gateway_auto_assigned_ = true
  ip_source = ""
  is_top_of_island = false
  lldp_search_string_auto_assigned_ = true
  located_by = ""
  locked = false
  managed_on_native_vlan = false
  out_of_band_management = false
  passphrase = ""
  passphrase_encrypted = ""
  password = ""
  password_encrypted = ""
  plane = ""
  plane_ref_type_ = "plane"
  pod = ""
  pod_ref_type_ = "pod"
  position = null
  power_state = ""
  private_password = ""
  private_password_encrypted = ""
  private_protocol = ""
  rack = ""
  rack_info = ""
  rack_ref_type_ = "rack"
  rail_group = null
  read_only_mode = false
  sdlc = ""
  security_type = ""
  snmp_community_string = ""
  snmpv3_username = ""
  spine_plane = ""
  spine_plane_ref_type_ = "spine_plane"
  ssh_key_or_password = ""
  ssh_key_or_password_encrypted_auto_assigned_ = true
  ssp_group = ""
  ssp_group_ref_type_ = "superspine_group"
  su = ""
  su_ref_type_ = "su"
  switch_gateway_auto_assigned_ = true
  switch_ip_and_mask_auto_assigned_ = true
  switch_router_id_ip_mask_auto_assigned_ = true
  switch_vtep_id_ip_mask_auto_assigned_ = true
  tenant = ""
  tenant_ref_type_ = "tenant"
  type = ""
  uplink_port = ""
  upstream_is_lag = false
  username_auto_assigned_ = true
  ztp_identification = ""

  badges {
    index = 1
    badge = ""
    badge_ref_type_ = "badge"
  }

  children {
    index = 1
    child_num_device = ""
    child_num_endpoint = ""
    child_num_endpoint_ref_type_ = "switchpoint"
  }

  eths {
    index = 1
    breakout = ""
    customer_vlan = ""
    enable = false
    eth_num_icon = ""
    eth_num_label = ""
    port_name = ""
  }

  object_properties {
    aggregate = false
    emulate_rf_video_port = false
    expected_parent_endpoint = ""
    expected_parent_endpoint_ref_type_ = "switchpoint"
    is_host = false
    number_of_multipoints = null
    user_notes = ""
  }

  traffic_mirrors {
    index = 1
    traffic_mirror_num_destination_port = ""
    traffic_mirror_num_enable = false
    traffic_mirror_num_inbound_traffic = false
    traffic_mirror_num_outbound_traffic = false
    traffic_mirror_num_source_lag_indicator = false
    traffic_mirror_num_source_port = ""
  }
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `authentication_protocol` (String) - Protocol.
* `badges` (Block List) - Badge configurations. Entries are matched by `index`.
  * `badge` (String) - Enable of this POTS port. Set together with `badge_ref_type_`.
  * `badge_ref_type_` (String) - Object type for badge field.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `bb_switch` (Boolean) - Expose fields for Device Management.
* `bgp_as_number` (Integer) - BGP Autonomous System Number for the Fabric Underlay. Datacenter mode only. Set it to `null` to clear it. Assigned by the server while `bgp_as_number_auto_assigned_` is `true`; it cannot be set then.
* `bgp_as_number_auto_assigned_` (Boolean) - Whether or not the value in bgp_as_number field has been automatically assigned or not. Set to false and change bgp_as_number value to edit. Datacenter mode only.
* `children` (Block List) - Child configurations. Entries are matched by `index`.
  * `child_num_device` (String) - Device associated with the Child.
  * `child_num_endpoint` (String) - Switchpoint associated with the Child. Set together with `child_num_endpoint_ref_type_`.
  * `child_num_endpoint_ref_type_` (String) - Object type for child_num_endpoint field.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `cli_access_mode` (String) - CLI Access Mode.
* `comm_type` (String) - Comm Type.
* `communication_mode` (String) - Select the network operating system (NOS) type for this endpoint.
* `connected_bundle` (String) - Connected Bundle. Set together with `connected_bundle_ref_type_`.
* `connected_bundle_ref_type_` (String) - Object type for connected_bundle field.
* `connection_service` (String) - Connect a Service. Campus mode only. Set together with `connection_service_ref_type_`.
* `connection_service_ref_type_` (String) - Object type for connection_service field. Campus mode only.
* `controller_ip_and_mask` (String) - Controller IP and Mask. Assigned by the server while `controller_ip_and_mask_auto_assigned_` is `true`; it cannot be set then.
* `controller_ip_and_mask_auto_assigned_` (Boolean) - Whether or not the value in controller_ip_and_mask field has been automatically assigned or not. Set to false and change controller_ip_and_mask value to edit.
* `device_managed_as` (String) - Device managed as. Campus mode only.
* `device_serial_number` (String) - Device Serial Number.
* `enable` (Boolean) - Enable object. It's highly recommended to set this value to true so that validation on the object will be ran.
* `enable_password` (String) - Enable Password - to enable privileged CLI operations.
* `enable_password_encrypted` (String) - Enable Password - to enable privileged CLI operations.
* `eths` (Block List) - Ethernet port configurations. Entries are matched by `index`.
  * `breakout` (String) - Breakout Port Override. Available options determined by Switch capability, Installed SFP and the capacity of the pipeline.
  * `customer_vlan` (String) - A Value between 1 and 4096.
  * `enable` (Boolean) - Enable port.
  * `eth_num_icon` (String) - Icon of this Eth Port.
  * `eth_num_label` (String) - Label of this Eth Port.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `port_name` (String) - The name identifying the port. Used for reference only, it won't actually change the port name.
* `expected_breakout` (String) - Full breakout configuration for the SFP being used as the uplink.
* `expected_breakout_uplink_port` (String) - Uplink Ethernet Port identifier in the format 1/# or 1/#/#.
* `expected_fabric` (String) - Expected Fabric. Set together with `expected_fabric_ref_type_`.
* `expected_fabric_ref_type_` (String) - Object type for expected_fabric field.
* `expected_uplink_port` (Integer) - Specify the uplink port for ZTP when using an SFP-based port. The Uplink Port is 1-indexed relative to the configured Port Group, where 1 represents the first port in the group. Port Group and breakout configurations are switch-model dependent; consult the switch vendor documentation to determine the correct Port Group and Uplink Port values. If an SFP-based uplink is not specified, ZTP programs the first 32 copper ports for use as uplinks. Set it to `null` to clear it.
* `gateway` (String) - Gateway. Assigned by the server while `gateway_auto_assigned_` is `true`; it cannot be set then.
* `gateway_auto_assigned_` (Boolean) - Whether or not the value in gateway field has been automatically assigned or not. Set to false and change gateway value to edit.
* `ip_source` (String) - IP Source.
* `is_fabric` (Boolean) - For Switch Endpoints. Denotes a Switch that is Fabric rather than an Edge Device. Campus mode only.
* `is_top_of_island` (Boolean) - Mark this Switchpoint as Top of Island.
* `lldp_search_string` (String) - Optional unless Located By is "LLDP" or Device managed as "Active SFP". Must be either the chassis-id or the hostname of the LLDP from the managed device. Used to detect connections between managed devices. If blank, the chassis-id detected by the Device Controller via SNMP/CLI is used. Assigned by the server while `lldp_search_string_auto_assigned_` is `true`; it cannot be set then.
* `lldp_search_string_auto_assigned_` (Boolean) - Whether or not the value in lldp_search_string field has been automatically assigned or not. Set to false and change lldp_search_string value to edit.
* `located_by` (String) - Controls how the system locates this Device within its LAN.
* `locked` (Boolean) - Permission lock.
* `managed_on_native_vlan` (Boolean) - Managed on native VLAN.
* `object_properties` (Block) - Object properties for the switchpoint. At most one block.
  * `aggregate` (Boolean) - For Switch Endpoints. Denotes switch aggregated with all of its sub switches.
  * `draw_as_edge_device` (Boolean) - Turn on to display the switch as an edge device instead of as a switch. Campus mode only.
  * `emulate_rf_video_port` (Boolean) - Emulate RF Video Port.
  * `expected_parent_endpoint` (String) - Expected Parent Endpoint. Set together with `expected_parent_endpoint_ref_type_`.
  * `expected_parent_endpoint_ref_type_` (String) - Object type for expected_parent_endpoint field.
  * `is_host` (Boolean) - For Switch Endpoints. Denotes the Host Switch.
  * `number_of_multipoints` (Integer) - Number of Multipoints. Set it to `null` to clear it.
  * `user_notes` (String) - Notes writen by User about the fabric.
* `out_of_band_management` (Boolean) - For Switch Endpoints. Denotes a Switch is managed out of band via the management port.
* `passphrase` (String) - Passphrase.
* `passphrase_encrypted` (String) - Passphrase.
* `password` (String) - Password.
* `password_encrypted` (String) - Password.
* `plane` (String) - Plane. Datacenter mode only. Set together with `plane_ref_type_`.
* `plane_ref_type_` (String) - Object type for plane field. Datacenter mode only.
* `pod` (String) - Pod - subgrouping of spine and leaf switches. Datacenter mode only. Set together with `pod_ref_type_`.
* `pod_ref_type_` (String) - Object type for pod field. Datacenter mode only.
* `port` (String) - Port locating the Switch to be controlled. Campus mode only.
* `position` (Number) - Position of the Switch. Datacenter mode only. Set it to `null` to clear it.
* `pots` (Block List) - POTS configurations. Entries are matched by `index`. Campus mode only.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `pots_num_caller_id` (String) - ASCII string defining the user for the Caller ID display for POTS port.
  * `pots_num_enable` (Boolean) - Enable POTS port.
  * `pots_num_hot_line` (String) - URI of line to autodial upon off-hook for POTS port.
  * `pots_num_password` (String) - SIP password used for authentication for POTS port.
  * `pots_num_password_encrypted` (String) - SIP password used for authentication for POTS port.
  * `pots_num_uri` (String) - Specific telephone extension for SIP for POTS port.
  * `pots_num_username` (String) - SIP username used for authentication for POTS port.
* `power_state` (String) - Power state of Switch Controller.
* `private_password` (String) - Password.
* `private_password_encrypted` (String) - Password.
* `private_protocol` (String) - Protocol.
* `rack` (String) - Rack. Datacenter mode only. Set together with `rack_ref_type_`.
* `rack_info` (String) - Physical Rack location of the Switch. Datacenter mode only.
* `rack_ref_type_` (String) - Object type for rack field. Datacenter mode only.
* `rail_group` (Number) - Rail Group the Switch is part of. Datacenter mode only. Set it to `null` to clear it.
* `read_only_mode` (Boolean) - When Read Only Mode is checked, vNetC will perform all functions except writing database updates to the target hardware.
* `sdlc` (String) - SDLC that Device Controller belongs to.
* `security_type` (String) - Security level.
* `snmp_community_string` (String) - Comm Credentials.
* `snmpv3_username` (String) - Username.
* `spine_plane` (String) - Spine Plane - subgrouping of super spine and spine. Datacenter mode only. Set together with `spine_plane_ref_type_`.
* `spine_plane_ref_type_` (String) - Object type for spine_plane field. Datacenter mode only.
* `ssh_key_or_password` (String) - SSH Key or Password.
* `ssh_key_or_password_encrypted` (String) - SSH Key or Password. Assigned by the server while `ssh_key_or_password_encrypted_auto_assigned_` is `true`; it cannot be set then.
* `ssh_key_or_password_encrypted_auto_assigned_` (Boolean) - Whether or not the value in ssh_key_or_password_encrypted field has been automatically assigned or not. Set to false and change ssh_key_or_password_encrypted value to edit.
* `ssp_group` (String) - SuperSpine Group - grouping of superspines in 3-tier config. Datacenter mode only. Set together with `ssp_group_ref_type_`.
* `ssp_group_ref_type_` (String) - Object type for ssp_group field. Datacenter mode only.
* `su` (String) - SU. Datacenter mode only. Set together with `su_ref_type_`.
* `su_ref_type_` (String) - Object type for su field. Datacenter mode only.
* `switch` (String) - Switchpoint locating the Switch to be controlled. Campus mode only. Set together with `switch_ref_type_`.
* `switch_gateway` (String) - Gateway of Managed Device. Assigned by the server while `switch_gateway_auto_assigned_` is `true`; it cannot be set then.
* `switch_gateway_auto_assigned_` (Boolean) - Whether or not the value in switch_gateway field has been automatically assigned or not. Set to false and change switch_gateway value to edit.
* `switch_ip_and_mask` (String) - Switch IP and Mask. Assigned by the server while `switch_ip_and_mask_auto_assigned_` is `true`; it cannot be set then.
* `switch_ip_and_mask_auto_assigned_` (Boolean) - Whether or not the value in switch_ip_and_mask field has been automatically assigned or not. Set to false and change switch_ip_and_mask value to edit.
* `switch_ref_type_` (String) - Object type for switch field. Campus mode only.
* `switch_router_id_ip_mask` (String) - Switch BGP Router Identifier. Datacenter mode only. Assigned by the server while `switch_router_id_ip_mask_auto_assigned_` is `true`; it cannot be set then.
* `switch_router_id_ip_mask_auto_assigned_` (Boolean) - Whether or not the value in switch_router_id_ip_mask field has been automatically assigned or not. Set to false and change switch_router_id_ip_mask value to edit. Datacenter mode only.
* `switch_vtep_id_ip_mask` (String) - Switch VETP Identifier. Datacenter mode only. Assigned by the server while `switch_vtep_id_ip_mask_auto_assigned_` is `true`; it cannot be set then.
* `switch_vtep_id_ip_mask_auto_assigned_` (Boolean) - Whether or not the value in switch_vtep_id_ip_mask field has been automatically assigned or not. Set to false and change switch_vtep_id_ip_mask value to edit. Datacenter mode only.
* `tenant` (String) - The Tenant of this Device. Set together with `tenant_ref_type_`.
* `tenant_ref_type_` (String) - Object type for tenant field.
* `traffic_mirrors` (Block List) - Traffic mirror configurations. Entries are matched by `index`.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `traffic_mirror_num_destination_port` (String) - Destination Port for Traffic Mirror.
  * `traffic_mirror_num_enable` (Boolean) - Enable Traffic Mirror.
  * `traffic_mirror_num_inbound_traffic` (Boolean) - Boolean value indicating if the mirror is for inbound traffic.
  * `traffic_mirror_num_outbound_traffic` (Boolean) - Boolean value indicating if the mirror is for outbound traffic.
  * `traffic_mirror_num_source_lag_indicator` (Boolean) - Source LAG Indicator for Traffic Mirror.
  * `traffic_mirror_num_source_port` (String) - Source Port for Traffic Mirror.
* `type` (String) - Type of Switchpoint. Datacenter mode only.
* `uplink_port` (String) - Uplink Port of Managed Device.
* `upstream_is_lag` (Boolean) - If checked, then ZTP will provision the TOR switch with the first 32 ports as a lag to facilitate plug-n-play.
* `username` (String) - Username. Assigned by the server while `username_auto_assigned_` is `true`; it cannot be set then.
* `username_auto_assigned_` (Boolean) - Whether or not the value in username field has been automatically assigned or not. Set to false and change username value to edit.
* `uses_tagged_packets` (Boolean) - Indicates if the direct interface expects tagged or untagged packets. Campus mode only.
* `ztp_identification` (String) - Service Tag or Serial Number to identify device for Zero Touch Provisioning.

## Reference Fields

A reference names another Verity object. Set the reference and its type field together. When your configuration sets neither, Terraform leaves the reference to the server, including one set in the Verity UI.

| Field | Type field | Allowed types |
| --- | --- | --- |
| `badges.badge` | `badge_ref_type_` | `badge` |
| `children.child_num_endpoint` | `child_num_endpoint_ref_type_` | `switchpoint` |
| `connected_bundle` | `connected_bundle_ref_type_` | `endpoint_bundle` |
| `connection_service` | `connection_service_ref_type_` | `service` |
| `expected_fabric` | `expected_fabric_ref_type_` | `fabric` |
| `object_properties.expected_parent_endpoint` | `expected_parent_endpoint_ref_type_` | `switchpoint` |
| `plane` | `plane_ref_type_` | `plane` |
| `pod` | `pod_ref_type_` | `pod` |
| `rack` | `rack_ref_type_` | `rack` |
| `spine_plane` | `spine_plane_ref_type_` | `spine_plane` |
| `ssp_group` | `ssp_group_ref_type_` | `superspine_group` |
| `su` | `su_ref_type_` | `su` |
| `switch` | `switch_ref_type_` | `switchpoint` |
| `tenant` | `tenant_ref_type_` | `tenant` |

## Auto-Assigned Fields

While a flag is `true`, the server assigns the paired value and the value cannot be set in the configuration. Set the flag to `false` to set the value yourself; if you leave the value unset, the value the server assigned is kept.

| Field | Flag | Reassigned when |
| --- | --- | --- |
| `bgp_as_number` | `bgp_as_number_auto_assigned_` | — |
| `controller_ip_and_mask` | `controller_ip_and_mask_auto_assigned_` | — |
| `gateway` | `gateway_auto_assigned_` | — |
| `lldp_search_string` | `lldp_search_string_auto_assigned_` | — |
| `ssh_key_or_password_encrypted` | `ssh_key_or_password_encrypted_auto_assigned_` | — |
| `switch_gateway` | `switch_gateway_auto_assigned_` | — |
| `switch_ip_and_mask` | `switch_ip_and_mask_auto_assigned_` | — |
| `switch_router_id_ip_mask` | `switch_router_id_ip_mask_auto_assigned_` | — |
| `switch_vtep_id_ip_mask` | `switch_vtep_id_ip_mask_auto_assigned_` | — |
| `username` | `username_auto_assigned_` | — |

## Import

Import an existing object by its `name`:

```sh
terraform import verity_switchpoint.<resource_name> <name>
```
