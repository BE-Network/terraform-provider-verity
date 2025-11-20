# Switchpoint Resource

`verity_switchpoint` manages switchpoint resources in Verity, which represent network devices like switches or routers.

## Example Usage

```hcl
resource "verity_switchpoint" "example" {
  name = "example"
  enable = true
  device_serial_number = "ABC123456"
  connected_bundle = "main_bundle"
  connected_bundle_ref_type_ = "endpoint_bundle"
  read_only_mode = false
  locked = false
  out_of_band_management = true
  type = "leaf"
  spine_plane = "plane1"
  spine_plane_ref_type_ = "spine_plane"
  pod = "rack1"
  pod_ref_type_ = "pod"
  rack = "A1"
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
    index = 0
    eth_num_icon = "server"
    eth_num_label = "Server Port"
    enable = true
    port_name = "eth0"
  }
  
  object_properties {
    user_notes = "Important device"
    expected_parent_endpoint = "parent_endpoint"
    expected_parent_endpoint_ref_type_ = "switchpoint"
    number_of_multipoints = 4
    aggregate = true
    is_host = false
    draw_as_edge_device = false
  }
}
```

## Argument Reference

* `name` (String) - Unique identifier for the switchpoint
* `enable` (Boolean) - Enable object
* `device_serial_number` (String) - Serial number of the physical device
* `connected_bundle` (String) - Reference to a bundle resource
* `connected_bundle_ref_type_` (String) - Object type for bundle reference
* `read_only_mode` (Boolean) - Whether the device is in read-only mode
* `locked` (Boolean) - Whether the device is locked
* `out_of_band_management` (Boolean) - Whether out-of-band management is enabled
* `type` (String) - Type of switchpoint
* `spine_plane` (String) - Spine Plane - subgrouping of super spine and spine
* `spine_plane_ref_type_` (String) - Object type for spine_plane field
* `pod` (String) - Pod - subgrouping of spine and leaf switches
* `pod_ref_type_` (String) - Object type for pod field
* `rack` (String) - Physical Rack location of the Switch
* `switch_router_id_ip_mask` (String) - Switch BGP Router Identifier
* `switch_router_id_ip_mask_auto_assigned_` (Boolean) - Whether router ID is auto-assigned
* `switch_vtep_id_ip_mask` (String) - Switch VETP Identifier
* `switch_vtep_id_ip_mask_auto_assigned_` (Boolean) - Whether VTEP ID is auto-assigned
* `bgp_as_number` (Integer) - BGP Autonomous System Number for the site underlay
* `bgp_as_number_auto_assigned_` (Boolean) - Whether BGP AS number is auto-assigned
* `is_fabric` (Boolean) - Whether this is a fabric switch
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
  * `index` (Integer) - Index identifying this port
  * `eth_num_icon` (String) - Icon of this Eth Port
  * `eth_num_label` (String) - Label of this Eth Port
  * `enable` (Boolean) - Enable port
  * `port_name` (String) - The name identifying the port
* `object_properties` (Object) - Object properties configuration
  * `user_notes` (String) - Notes written by User about the site
  * `expected_parent_endpoint` (String) - Expected Parent Endpoint
  * `expected_parent_endpoint_ref_type_` (String) - Object type for expected_parent_endpoint field
  * `number_of_multipoints` (Integer) - Number of Multipoints
  * `aggregate` (Boolean) - Whether this is an aggregate device
  * `is_host` (Boolean) - Whether this device is a host
  * `draw_as_edge_device` (Boolean) - Turn on to display the switch as an edge device instead of as a switch

## Import

Switchpoint resources can be imported using the `name` attribute:

```sh
terraform import verity_switchpoint.<resource_name> <name>
```
