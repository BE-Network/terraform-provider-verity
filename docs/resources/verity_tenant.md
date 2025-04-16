# Tenant Resource

`verity_tenant` manages tenant resources in Verity, which define isolated network environments.

## Example Usage

```hcl
resource "verity_tenant" "example" {
  name = "example"
  object_properties {
    group = ""
  }
  dhcp_relay_source_ips_subnet = ""
  route_distinguisher = ""
  import_route_map = ""
  vrf_name = "Vrf_9"
  route_tenants {
    index = 1
    enable = true
    tenant = ""
    tenant_ref_type_ = ""
  }
  layer_3_vni = 102009
  route_target_import = ""
  export_route_map_ref_type_ = ""
  layer_3_vni_auto_assigned_ = true
  enable = true
  route_target_export = ""
  layer_3_vlan_auto_assigned_ = true
  vrf_name_auto_assigned_ = true
  layer_3_vlan = 2009
  export_route_map = ""
  import_route_map_ref_type_ = ""
}
```

## Argument Reference

* `name` - Unique identifier for the tenant
* `enable` - Enable this tenant. Default is `false`
* `object_properties` - Object properties block
  * `group` - Group name
* `layer_3_vni` - Layer 3 VNI value
* `layer_3_vni_auto_assigned_` - Whether Layer 3 VNI is auto-assigned
* `layer_3_vlan` - Layer 3 VLAN ID
* `layer_3_vlan_auto_assigned_` - Whether Layer 3 VLAN ID is auto-assigned
* `dhcp_relay_source_ips_subnet` - DHCP relay source IPs subnet
* `route_distinguisher` - Route distinguisher for BGP
* `route_target_import` - Route target import value for BGP
* `route_target_export` - Route target export value for BGP
* `import_route_map` - Import route map
* `import_route_map_ref_type_` - Object type for import route map
* `export_route_map` - Export route map
* `export_route_map_ref_type_` - Object type for export route map
* `vrf_name` - VRF name
* `vrf_name_auto_assigned_` - Whether VRF name is auto-assigned
* `route_tenants` - List of route tenant blocks
  * `enable` - Enable this route tenant
  * `tenant` - Reference to another tenant
  * `tenant_ref_type_` - Object type for tenant reference
  * `index` - Index value for ordering

## Import

Tenant resources can be imported using the `name` attribute:

