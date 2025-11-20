# Packet Broker Resource

`verity_packet_broker` manages packet broker (PB Egress Profile) resources in Verity, which define how packets are filtered and processed.

## Example Usage

```hcl
resource "verity_packet_broker" "example" {
  name = "example"
  enable = true
  
  ipv4_permit {
    enable = true
    filter = "permit_filter_1"
    filter_ref_type_ = "ipv4_filter"
    index = 1
  }
  
  ipv4_deny {
    enable = true
    filter = "deny_filter_1"
    filter_ref_type_ = "ipv4_filter"
    index = 1
  }
  
  ipv6_permit {
    enable = true
    filter = "permit_filter_v6"
    filter_ref_type_ = "ipv6_filter"
    index = 1
  }
  
  ipv6_deny {
    enable = true
    filter = "deny_filter_v6"
    filter_ref_type_ = "ipv6_filter"
    index = 1
  }
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique.
* `enable` (Boolean) - Enable object.
* `ipv4_permit` (Array) - 
  * `enable` (Boolean) - Enable.
  * `filter` (String) - Filter.
  * `filter_ref_type_` (String) - Object type for filter field.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `ipv4_deny` (Array) - 
  * `enable` (Boolean) - Enable.
  * `filter` (String) - Filter.
  * `filter_ref_type_` (String) - Object type for filter field.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `ipv6_permit` (Array) - 
  * `enable` (Boolean) - Enable.
  * `filter` (String) - Filter.
  * `filter_ref_type_` (String) - Object type for filter field.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `ipv6_deny` (Array) - 
  * `enable` (Boolean) - Enable.
  * `filter` (String) - Filter.
  * `filter_ref_type_` (String) - Object type for filter field.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.

## Import

Packet Broker resources can be imported using the `name` attribute:

```sh
terraform import verity_packet_broker.<resource_name> <name>
```
