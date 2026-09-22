# verity_service_port_profile (Resource)

Manages a Verity Service Port Profile.

Supported modes: Campus.

## Example Usage

```hcl
resource "verity_service_port_profile" "example" {
  name = "example"
  enable = false
  ip_mask = ""
  port_type = ""
  tls_limit_in = null
  tls_service = ""
  tls_service_ref_type_ = "service"
  trusted_port = false

  object_properties {
    on_summary = false
    port_monitoring = ""
  }

  services {
    index = 1
    row_num_enable = false
    row_num_external_vlan = null
    row_num_limit_in = null
    row_num_limit_out = null
    row_num_service = ""
    row_num_service_ref_type_ = "service"
  }
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `enable` (Boolean) - Enable object.
* `ip_mask` (String) - IP/Mask.
* `object_properties` (Block) - Object properties for the service port profile. At most one block.
  * `on_summary` (Boolean) - Show on the summary view.
  * `port_monitoring` (String) - Defines importance of Link Down on this port.
* `port_type` (String) - Determines what Service are provisioned on the port and if those Services are propagated upstream "Upstream Switchport" Services specified below. Services are not propagated.; "Downstream Switchport" Services specified below. Services are propagated.; "Crosslink Switchport" Services is union of all Services on each switch. Services are not propagated.; "Upstream L3 (L2/L3 Switches Only" No Services.
* `services` (Block List) - Service configurations. Entries are matched by `index`.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `row_num_enable` (Boolean) - Enable row.
  * `row_num_external_vlan` (Integer) - Choose an external vlan. Set it to `null` to clear it.
  * `row_num_limit_in` (Integer) - Speed of ingress (Mbps). Set it to `null` to clear it.
  * `row_num_limit_out` (Integer) - Speed of egress (Mbps). Set it to `null` to clear it.
  * `row_num_service` (String) - Connect a Service. Set together with `row_num_service_ref_type_`.
  * `row_num_service_ref_type_` (String) - Object type for row_num_service field.
* `tls_limit_in` (Integer) - Speed of ingress (Mbps) for TLS (Transparent LAN Service). Set it to `null` to clear it.
* `tls_service` (String) - Service used for TLS (Transparent LAN Service). Set together with `tls_service_ref_type_`.
* `tls_service_ref_type_` (String) - Object type for tls_service field.
* `trusted_port` (Boolean) - Trusted Ports do not participate in IP Source Guard, Dynamic ARP Inspection, nor DHCP Snooping, meaning all packets are forwarded without any checks.

## Reference Fields

A reference names another Verity object. Set the reference and its type field together. When your configuration sets neither, Terraform leaves the reference to the server, including one set in the Verity UI.

| Field | Type field | Allowed types |
| --- | --- | --- |
| `services.row_num_service` | `row_num_service_ref_type_` | `service` |
| `tls_service` | `tls_service_ref_type_` | `service` |

## Import

Import an existing object by its `name`:

```sh
terraform import verity_service_port_profile.<resource_name> <name>
```
