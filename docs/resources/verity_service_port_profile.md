# Service Port Profile Resource

`verity_service_port_profile` manages service port profile resources in Verity, which define port configurations and associated services.

## Example Usage

```hcl
resource "verity_service_port_profile" "example" {
  name = "example"
  enable = true
  port_type = "access"
  tls_limit_in = 1000
  tls_service = "example_service"
  tls_service_ref_type_ = "service"
  trusted_port = true
  ip_mask = "255.255.255.0"
  
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

* `name` (String) - Object Name. Must be unique
* `enable` (Boolean) - Enable object
* `port_type` (String) - Determines what Service are provisioned on the port and if those Services are propagated upstream
* `tls_limit_in` (Integer) - Speed of ingress (Mbps) for TLS (Transparent LAN Service)
* `tls_service` (String) - Service used for TLS (Transparent LAN Service)
* `tls_service_ref_type_` (String) - Object type for tls_service field
* `trusted_port` (Boolean) - Trusted Ports do not participate in IP Source Guard, Dynamic ARP Inspection, nor DHCP Snooping, meaning all packets are forwarded without any checks
* `ip_mask` (String) - IP/Mask
* `services` (Array) - List of service configurations
  * `row_num_enable` (Boolean) - Enable row
  * `row_num_service` (String) - Connect a Service
  * `row_num_service_ref_type_` (String) - Object type for row_num_service field
  * `row_num_external_vlan` (Integer) - Choose an external vlan
  * `row_num_limit_in` (Integer) - Speed of ingress (Mbps)
  * `row_num_limit_out` (Integer) - Speed of egress (Mbps)
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list
* `object_properties` (Object) - Additional object properties
  * `on_summary` (Boolean) - Show on the summary view
  * `port_monitoring` (String) - Defines importance of Link Down on this port
  * `group` (String) - Group

## Import

Service Port Profile resources can be imported using the `name` attribute:

```sh
terraform import verity_service_port_profile.<resource_name> <name>
```
