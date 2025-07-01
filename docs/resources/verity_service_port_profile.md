# Service Port Profile Resource

`verity_service_port_profile` manages service port profile resources in Verity, which define port configurations and associated services.

## Example Usage

```hcl
resource "verity_service_port_profile" "example" {
  name = "example_service_port_profile"
  enable = true
  port_type = "access"
  tls_limit_in = 1000
  tls_service = "example_service"
  tls_service_ref_type_ = "service"
  trusted_port = true
  ip_mask = "255.255.255.0"
  
  services {
    row_num_enable = true
    row_num_service = "service1"
    row_num_service_ref_type_ = "service"
    row_num_external_vlan = 100
    row_num_limit_in = 2000
    row_num_limit_out = 2000
    index = 0
  }
  
  services {
    row_num_enable = true
    row_num_service = "service2"
    row_num_service_ref_type_ = "service"
    row_num_external_vlan = 200
    row_num_limit_in = 3000
    row_num_limit_out = 3000
    index = 1
  }
  
  object_properties {
    on_summary = true
    port_monitoring = "enabled"
    group = "profile-group"
  }
}
```

## Argument Reference

* `name` - (Required) Unique identifier for the service port profile.
* `enable` - (Optional) Enable this service port profile. Default is `false`.
* `port_type` - (Optional) Type of port (e.g., `access`, `trunk`).
* `tls_limit_in` - (Optional) Inbound bandwidth limit for TLS in Kbps.
* `tls_service` - (Optional) TLS service reference.
* `tls_service_ref_type_` - (Optional) Reference type for TLS service.
* `trusted_port` - (Optional) Whether this is a trusted port. Default is `false`.
* `ip_mask` - (Optional) IP mask.
* `services` - (Optional) List of service configurations:
  * `row_num_enable` - (Optional) Enable this service. Default is `false`.
  * `row_num_service` - (Optional) Service name or identifier.
  * `row_num_service_ref_type_` - (Optional) Reference type for service.
  * `row_num_external_vlan` - (Optional) External VLAN ID.
  * `row_num_limit_in` - (Optional) Inbound bandwidth limit in Kbps.
  * `row_num_limit_out` - (Optional) Outbound bandwidth limit in Kbps.
  * `index` - (Optional) Index identifying this service configuration.
* `object_properties` - (Optional) Object properties configuration:
  * `on_summary` - (Optional) Whether this profile appears on summary. Default is `false`.
  * `port_monitoring` - (Optional) Port monitoring status.
  * `group` - (Optional) Group name.

## Import

Service Port Profile resources can be imported using the `name` attribute:

```
$ terraform import verity_service_port_profile.example example_service_port_profile
```
