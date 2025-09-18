# ACL Resources

`verity_acl_v4` and `verity_acl_v6` manage ACL (Access Control List) IP filter resources in Verity, which define network filtering rules for IPv4 and IPv6 traffic respectively.

## Version Compatibility

**These resources require Verity API version 6.5 or higher.**

## Example Usage

```hcl
# IPv4 Filter example
resource "verity_acl_v4" "ipv4_filter_example" {
  name = "ipv4_filter_example"
  enable = true
  protocol = "tcp"
  bidirectional = false
  source_ip = "192.168.1.0/24"
  source_port_operator = "range"
  source_port_1 = 1024
  source_port_2 = 2048
  destination_ip = "10.0.0.0/24"
  destination_port_operator = "equal"
  destination_port_1 = 80
  
  object_properties {
    notes = ""
  }
}

# IPv6 Filter example
resource "verity_acl_v6" "ipv6_filter_example" {
  name = "ipv6_filter_example"
  enable = true
  protocol = "udp"
  bidirectional = true
  source_ip = "2001:db8::/64"
  source_port_operator = "greater-than"
  source_port_1 = 1024
  destination_ip = "2001:db8:1::/64"
  destination_port_operator = "equal"
  destination_port_1 = 53
  
  object_properties {
    notes = ""
  }
}
```

## Argument Reference

* `name` - (Required) Unique identifier for the IP filter.
* `enable` - (Optional) Enable this IP filter. Default is `false`.
* `protocol` - (Optional) Protocol to match. Value must be ip/tcp/udp/icmp or a number between 0 and 255. Value `ip` will match all IP protocols.
* `bidirectional` - (Optional) If set to `true`, packets will be selected that match the source filters in either the source or destination fields of the packet. Default is `false`.
* `source_ip` - (Optional) Source IP address or subnet to match in CIDR notation.
* `source_port_operator` - (Optional) Match operation to apply to TCP/UDP ports. Values are `equal`, `greater-than`, `less-than`, or `range`.
* `source_port_1` - (Optional) Value for equal, greater-than, or less-than TCP/UDP port match operations, or the lower value in a range operation.
* `source_port_2` - (Optional) Upper value in a range TCP/UDP port match operation. Only used when `source_port_operator` is set to `range`.
* `destination_ip` - (Optional) Destination IP address or subnet to match in CIDR notation.
* `destination_port_operator` - (Optional) Match operation to apply to TCP/UDP ports. Values are `equal`, `greater-than`, `less-than`, or `range`.
* `destination_port_1` - (Optional) Value for equal, greater-than, or less-than TCP/UDP port match operations, or the lower value in a range operation.
* `destination_port_2` - (Optional) Upper value in a range TCP/UDP port match operation. Only used when `destination_port_operator` is set to `range`.
* `object_properties` - (Optional) Additional properties for the IP filter:
  * `notes` - (Optional) User notes for the IP filter.

## Import

ACL IP Filter resources can be imported using the `name` attribute:

```sh
terraform import verity_acl_v4.<resource_name> <name>
terraform import verity_acl_v6.<resource_name> <name>
```
