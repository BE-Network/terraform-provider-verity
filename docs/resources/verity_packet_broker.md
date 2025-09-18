# Packet Broker Resource

`verity_packet_broker` manages packet broker (PB Egress Profile) resources in Verity, which define how packets are filtered and processed.

## Version Compatibility

**This resource requires Verity API version 6.5 or higher.**

## Example Usage

```hcl
resource "verity_packet_broker" "example" {
  name = "example_packet_broker"
  enable = true
  
  ipv4_permit {
    enable = true
    filter = "permit_filter_1"
    filter_ref_type_ = "ipv4_prefix_list"
    index = 1
  }
  
  ipv4_deny {
    enable = true
    filter = "deny_filter_1"
    filter_ref_type_ = "ipv4_prefix_list"
    index = 2
  }
  
  ipv6_permit {
    enable = true
    filter = "permit_filter_v6"
    filter_ref_type_ = "ipv6_prefix_list"
    index = 3
  }
  
  ipv6_deny {
    enable = true
    filter = "deny_filter_v6"
    filter_ref_type_ = "ipv6_prefix_list"
    index = 4
  }
}
```

## Argument Reference

* `name` - (Required) Unique identifier for the packet broker.
* `enable` - (Optional) Enable this packet broker. Default is `false`.
* `ipv4_permit` - (Optional) List of IPv4 permit filter configurations:
  * `enable` - (Optional) Enable this filter. Default is `false`.
  * `filter` - (Optional) Reference to a filter resource.
  * `filter_ref_type_` - (Optional) Object type for filter reference.
  * `index` - (Optional) Index identifying this filter configuration.
* `ipv4_deny` - (Optional) List of IPv4 deny filter configurations:
  * `enable` - (Optional) Enable this filter. Default is `false`.
  * `filter` - (Optional) Reference to a filter resource.
  * `filter_ref_type_` - (Optional) Object type for filter reference.
  * `index` - (Optional) Index identifying this filter configuration.
* `ipv6_permit` - (Optional) List of IPv6 permit filter configurations:
  * `enable` - (Optional) Enable this filter. Default is `false`.
  * `filter` - (Optional) Reference to a filter resource.
  * `filter_ref_type_` - (Optional) Object type for filter reference.
  * `index` - (Optional) Index identifying this filter configuration.
* `ipv6_deny` - (Optional) List of IPv6 deny filter configurations:
  * `enable` - (Optional) Enable this filter. Default is `false`.
  * `filter` - (Optional) Reference to a filter resource.
  * `filter_ref_type_` - (Optional) Object type for filter reference.
  * `index` - (Optional) Index identifying this filter configuration.

## Import

Packet Broker resources can be imported using the `name` attribute:

```sh
terraform import verity_packet_broker.<resource_name> <name>
```
