# Tenant Resource

`verity_tenant` manages tenant resources in Verity, which define isolated network environments.

## Example Usage

```hcl
resource "verity_tenant" "example" {
  name = "example"
  enable = true
  layer_3_vni = null
  layer_3_vni_auto_assigned_ = false
  layer_3_vlan = null
  layer_3_vlan_auto_assigned_ = false
  dhcp_relay_source_ips_subnet = ""
  route_distinguisher = ""
  route_target_import = ""
  route_target_export = ""
  import_route_map = ""
  import_route_map_ref_type_ = ""
  export_route_map = ""
  export_route_map_ref_type_ = ""
  vrf_name = "(auto)"
  vrf_name_auto_assigned_ = false
  default_originate = false

  route_tenants {
    index = 1
    enable = false
    tenant = ""
  }

  object_properties {
    group = ""
  }
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique.
* `enable` (Boolean) - Enable object. It's highly recommended to set this value to true so that validation on the object will be ran.
* `layer_3_vni` (Integer) - VNI value used to transport traffic between services of a Tenant.
* `layer_3_vni_auto_assigned_` (Boolean) - Whether or not the value in layer_3_vni field has been automatically assigned or not. Set to false and change layer_3_vni value to edit.
* `layer_3_vlan` (Integer) - VLAN value used to transport traffic between services of a Tenant.
* `layer_3_vlan_auto_assigned_` (Boolean) - Whether or not the value in layer_3_vlan field has been automatically assigned or not. Set to false and change layer_3_vlan value to edit.
* `dhcp_relay_source_ips_subnet` (String) - Range of IP addresses (represented in IP subnet format) used to configure the source IP of each DHCP Relay on each switch that this Tenant is provisioned on.
* `route_distinguisher` (String) - Route Distinguishers are used to maintain uniqueness among identical routes from different routers. If set, then routes from this Tenant will be identified with this Route Distinguisher (BGP Community). It should be two numbers separated by a colon.
* `route_target_import` (String) - A route-target (BGP Community) to attach while importing routes into the current tenant. It should be a comma-separated list of BGP Communities: each Community being two numbers separated by a colon.
* `route_target_export` (String) - A route-target (BGP Community) to attach while exporting routes from the current tenant. It should be a comma-separated list of BGP Communities: each Community being two numbers separated by a colon.
* `import_route_map` (String) - A route-map applied to routes imported into the current tenant from other tenants with the purpose of filtering or modifying the routes.
* `import_route_map_ref_type_` (String) - Object type for import_route_map field.
* `export_route_map` (String) - A route-map applied to routes exported into the current tenant from other tenants with the purpose of filtering or modifying the routes.
* `export_route_map_ref_type_` (String) - Object type for export_route_map field.
* `vrf_name` (String) - Virtual Routing and Forwarding instance name associated to tenants.
* `vrf_name_auto_assigned_` (Boolean) - Whether or not the value in vrf_name field has been automatically assigned or not. Set to false and change vrf_name value to edit.
* `route_tenants` (Array) - 
  * `enable` (Boolean) - Enable.
  * `tenant` (String) - Tenant.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `object_properties` (Object) - 
  * `group` (String) - Group.
* `default_originate` (Boolean) - Enables a leaf switch to originate IPv4 default type-5 EVPN routes across the switching fabric.

## Import

Tenant resources can be imported using the `name` attribute:

```sh
terraform import verity_tenant.<resource_name> <name>
```