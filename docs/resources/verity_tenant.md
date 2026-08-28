# Tenant Resource

`verity_tenant` manages tenant resources in Verity, which define isolated network environments.

## Example Usage

```hcl
resource "verity_tenant" "example" {
  name = "example"
  dhcp_relay_source_ipv4s_subnet = ""
  dhcp_relay_source_ipv6s_subnet = ""
  route_distinguisher = ""
  route_aggregation = "manual"
  import_route_map = ""
  vrf_name = "Vrf_9"
  layer_3_vni = 102009
  route_target_import = ""
  layer_3_vni_auto_assigned_ = true
  enable = true
  route_target_export = ""
  layer_3_vlan_auto_assigned_ = true
  vrf_name_auto_assigned_ = true
  layer_3_vlan = 2009
  export_route_map = ""
  default_originate = false

  route_tenants {
    index = 1
    enable = true
    tenant = ""
  }

  route_aggregators {
    index = 1
    route_aggregation_num_enable = true
    route_aggregation_num_ip_and_mask = "10.0.0.0/8"
  }
}
```

## Argument Reference

* `name` (String) - Unique identifier for the tenant
* `enable` (Boolean) - Enable this tenant. Default is `false`
* `layer_3_vni` (Integer) - Layer 3 VNI value
* `layer_3_vni_auto_assigned_` (Boolean) - Whether Layer 3 VNI is auto-assigned
* `layer_3_vlan` (Integer) - Layer 3 VLAN ID
* `layer_3_vlan_auto_assigned_` (Boolean) - Whether Layer 3 VLAN ID is auto-assigned
* `dhcp_relay_source_ipv4s_subnet` (String) - IPv4 subnet used to allocate Tenant-specific loopback addresses on each leaf switch for DHCP relay and troubleshooting.
* `dhcp_relay_source_ipv6s_subnet` (String) - IPv6 subnet used to allocate Tenant-specific loopback addresses on each leaf switch for DHCP relay and troubleshooting.
* `route_distinguisher` (String) - Route distinguisher for BGP
* `route_aggregation` (String) - Route aggregation configuration. Valid values are `manual` and `automated`; omit or set to an empty string to disable it.
* `route_aggregators` (Array) - Route aggregation entries
  * `index` (Integer) - The index identifying the object. Set to `0` to add an entry
  * `route_aggregation_num_enable` (Boolean) - Enable this route aggregation entry
  * `route_aggregation_num_ip_and_mask` (String) - IP address and mask for route aggregation
* `route_target_import` (String) - Route target import value for BGP
* `route_target_export` (String) - Route target export value for BGP
* `import_route_map` (String) - Import route map
* `import_route_map_ref_type_` (String) - Object type for import route map
* `export_route_map` (String) - Export route map
* `export_route_map_ref_type_` (String) - Object type for export route map
* `vrf_name` (String) - VRF name
* `vrf_name_auto_assigned_` (Boolean) - Whether VRF name is auto-assigned
* `tenant_type` (String) - Tenant type.
* `route_tenants` (Array) - List of route tenant blocks
  * `enable` (Boolean) - Enable this route tenant
  * `tenant` (String) - Reference to another tenant
  * `index` (Integer) - Index value for ordering
* `default_originate` (Boolean) - Enables a leaf switch to originate IPv4 default type-5 EVPN routes across the switching fabric

## Import

Tenant resources can be imported using the `name` attribute:

```sh
terraform import verity_tenant.<resource_name> <name>
```
