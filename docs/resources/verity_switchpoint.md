# Switchpoint Resource

`verity_switchpoint` manages switchpoint resources in Verity, which represent network devices like switches or routers.

## Example Usage

```hcl
resource "verity_switchpoint" "example" {
  name = "example_switch"
  device_serial_number = "ABC123456"
  connected_bundle = "main_bundle"
  connected_bundle_ref_type_ = "bundle"
  read_only_mode = false
  locked = false
  disabled_ports = "1-4"
  out_of_band_management = true
  type = "access"
  super_pod = "pod1"
  pod = "rack1"
  rack = "A1"
  switch_router_id_ip_mask = "192.168.1.1"
  switch_router_id_ip_mask_auto_assigned_ = false
  switch_vtep_id_ip_mask = "192.168.2.1"
  switch_vtep_id_ip_mask_auto_assigned_ = false
  bgp_as_number = 65001
  bgp_as_number_auto_assigned_ = false
  
  badges {
    badge = "badge1"
    badge_ref_type_ = "badge"
    index = 0
  }
  
  children {
    child_num_endpoint = "child_endpoint"
    child_num_endpoint_ref_type_ = "endpoint"
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
  }
  
  eths {
    breakout = "1x100G"
    index = 0
  }
  
  object_properties {
    user_notes = "Important device"
    expected_parent_endpoint = "parent_endpoint"
    expected_parent_endpoint_ref_type_ = "endpoint"
    number_of_multipoints = 4
    aggregate = true
    is_host = false
    
    eths {
      eth_num_icon = "server"
      eth_num_label = "Server Port"
      index = 0
    }
  }
}
```

## Argument Reference

* `name` - (Required) Unique identifier for the switchpoint.
* `device_serial_number` - (Optional) Serial number of the physical device.
* `connected_bundle` - (Optional) Reference to a bundle resource.
* `connected_bundle_ref_type_` - (Optional) Object type for bundle reference.
* `read_only_mode` - (Optional) Whether the device is in read-only mode.
* `locked` - (Optional) Whether the device is locked.
* `disabled_ports` - (Optional) Range of disabled ports (e.g., "1-4,6,8-10").
* `out_of_band_management` - (Optional) Whether out-of-band management is enabled.
* `type` - (Optional) Type of switchpoint (e.g., "access", "core").
* `super_pod` - (Optional) Super pod identifier.
* `pod` - (Optional) Pod identifier.
* `rack` - (Optional) Rack identifier.
* `switch_router_id_ip_mask` - (Optional) IP address for router ID.
* `switch_router_id_ip_mask_auto_assigned_` - (Optional) Whether router ID is auto-assigned.
* `switch_vtep_id_ip_mask` - (Optional) IP address for VTEP ID.
* `switch_vtep_id_ip_mask_auto_assigned_` - (Optional) Whether VTEP ID is auto-assigned.
* `bgp_as_number` - (Optional) BGP AS number.
* `bgp_as_number_auto_assigned_` - (Optional) Whether BGP AS number is auto-assigned.
* `badges` - (Optional) List of badge configurations:
  * `badge` - (Optional) Badge reference.
  * `badge_ref_type_` - (Optional) Object type for badge reference.
  * `index` - (Optional) Index identifying this badge.
* `children` - (Optional) List of child device configurations:
  * `child_num_endpoint` - (Optional) Child endpoint reference.
  * `child_num_endpoint_ref_type_` - (Optional) Object type for endpoint reference.
  * `child_num_device` - (Optional) Child device reference.
  * `index` - (Optional) Index identifying this child.
* `traffic_mirrors` - (Optional) List of traffic mirroring configurations:
  * `traffic_mirror_num_enable` - (Optional) Enable this traffic mirror.
  * `traffic_mirror_num_source_port` - (Optional) Source port for mirrored traffic.
  * `traffic_mirror_num_source_lag_indicator` - (Optional) Whether source is a LAG.
  * `traffic_mirror_num_destination_port` - (Optional) Destination port for mirrored traffic.
  * `traffic_mirror_num_inbound_traffic` - (Optional) Mirror inbound traffic.
  * `traffic_mirror_num_outbound_traffic` - (Optional) Mirror outbound traffic.
* `eths` - (Optional) List of Ethernet port configurations:
  * `breakout` - (Optional) Breakout configuration (e.g., "1x100G", "4x25G").
  * `index` - (Optional) Index identifying this port.
* `object_properties` - (Optional) Object properties configuration:
  * `user_notes` - (Optional) User notes for this device.
  * `expected_parent_endpoint` - (Optional) Parent endpoint reference.
  * `expected_parent_endpoint_ref_type_` - (Optional) Object type for parent endpoint reference.
  * `number_of_multipoints` - (Optional) Number of multipoints.
  * `aggregate` - (Optional) Whether this is an aggregate device.
  * `is_host` - (Optional) Whether this device is a host.
  * `eths` - (Optional) List of Ethernet port properties:
    * `eth_num_icon` - (Optional) Icon for this port.
    * `eth_num_label` - (Optional) Label for this port.
    * `index` - (Optional) Index identifying this port.

## Import

Switchpoint resources can be imported using the `name` attribute:

```
$ terraform import verity_switchpoint.example example_switch
```
