# verity_tenant (Resource)

Manages a Tenant resource.

Supported modes: Datacenter.

## Example Usage

```hcl
resource "verity_tenant" "example" {
  name = "example"
  default_originate = false
  dhcp_relay_source_ipv4s_subnet = ""
  dhcp_relay_source_ipv6s_subnet = ""
  enable = false
  export_route_map = ""
  export_route_map_ref_type_ = "route_map"
  import_route_map = ""
  import_route_map_ref_type_ = "route_map"
  layer_3_vlan_auto_assigned_ = true
  layer_3_vni_auto_assigned_ = true
  route_aggregation = ""
  route_distinguisher = ""
  route_target_export = ""
  route_target_import = ""
  tenant_type = ""
  vrf_name_auto_assigned_ = true

  route_aggregators {
    index = 1
    route_aggregation_num_enable = false
    route_aggregation_num_ip_and_mask = ""
  }

  route_tenants {
    index = 1
    enable = false
    tenant = ""
  }
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `default_originate` (Boolean) - When enabled, provision an underlay gateway on the switch for this tenant.
* `dhcp_relay_source_ipv4s_subnet` (String) - IPv4 subnet used to allocate Tenant-specific loopback addresses on each leaf switch for DHCP relay and troubleshooting.
* `dhcp_relay_source_ipv6s_subnet` (String) - IPv6 subnet used to allocate Tenant-specific loopback addresses on each leaf switch for DHCP relay and troubleshooting.
* `enable` (Boolean) - Enable object. It's highly recommended to set this value to true so that validation on the object will be ran.
* `export_route_map` (String) - A route-map applied to routes exported from the current tenant to other tenants with the purpose of filtering or modifying the routes. Set together with `export_route_map_ref_type_`.
* `export_route_map_ref_type_` (String) - Object type for export_route_map field.
* `import_route_map` (String) - A route-map applied to routes imported into the current tenant from other tenants with the purpose of filtering or modifying the routes. Set together with `import_route_map_ref_type_`.
* `import_route_map_ref_type_` (String) - Object type for import_route_map field.
* `layer_3_vlan` (Integer) - VLAN value used to transport traffic between services of a Tenant. Set it to `null` to clear it. Assigned by the server while `layer_3_vlan_auto_assigned_` is `true`; it cannot be set then.
* `layer_3_vlan_auto_assigned_` (Boolean) - Whether or not the value in layer_3_vlan field has been automatically assigned or not. Set to false and change layer_3_vlan value to edit.
* `layer_3_vni` (Integer) - VNI value used to transport traffic between services of a Tenant. Set it to `null` to clear it. Assigned by the server while `layer_3_vni_auto_assigned_` is `true`; it cannot be set then.
* `layer_3_vni_auto_assigned_` (Boolean) - Whether or not the value in layer_3_vni field has been automatically assigned or not. Set to false and change layer_3_vni value to edit.
* `route_aggregation` (String) - Route Aggregation configuration for this tenant.
* `route_aggregators` (Block List) - Route aggregation entries. Entries are matched by `index`.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `route_aggregation_num_enable` (Boolean) - Enable.
  * `route_aggregation_num_ip_and_mask` (String) - IP address and mask for route aggregation.
* `route_distinguisher` (String) - Route Distinguishers are used to maintain uniqueness among identical routes from different routers. If set, then routes from this Tenant will be identified with this Route Distinguisher (BGP Community). It should be two numbers separated by a colon.
* `route_target_export` (String) - A route-target (BGP Community) to attach while exporting routes from the current tenant. It should be a comma-separated list of BGP Communities: each Community being two numbers separated by a colon.
* `route_target_import` (String) - A route-target (BGP Community) to attach while importing routes into the current tenant. It should be a comma-separated list of BGP Communities: each Community being two numbers separated by a colon.
* `route_tenants` (Block List) - Route tenants configuration. Entries are matched by `index`.
  * `enable` (Boolean) - Enable.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `tenant` (String) - Tenant.
* `tenant_type` (String) - Type of Tenant. To Provision on Spectrum-X sites, select East-West.
* `vrf_name` (String) - Virtual Routing and Forwarding instance name associated to tenants. Assigned by the server while `vrf_name_auto_assigned_` is `true`; it cannot be set then.
* `vrf_name_auto_assigned_` (Boolean) - Whether or not the value in vrf_name field has been automatically assigned or not. Set to false and change vrf_name value to edit.

## Reference Fields

A reference names another Verity object. Set the reference and its type field together. When your configuration sets neither, Terraform leaves the reference to the server, including one set in the Verity UI.

| Field | Type field | Allowed types |
| --- | --- | --- |
| `export_route_map` | `export_route_map_ref_type_` | `route_map` |
| `import_route_map` | `import_route_map_ref_type_` | `route_map` |

## Auto-Assigned Fields

While a flag is `true`, the server assigns the paired value and the value cannot be set in the configuration. Set the flag to `false` to set the value yourself; if you leave the value unset, the value the server assigned is kept.

| Field | Flag | Reassigned when |
| --- | --- | --- |
| `layer_3_vlan` | `layer_3_vlan_auto_assigned_` | — |
| `layer_3_vni` | `layer_3_vni_auto_assigned_` | — |
| `vrf_name` | `vrf_name_auto_assigned_` | — |

## Import

Import an existing object by its `name`:

```sh
terraform import verity_tenant.<resource_name> <name>
```
