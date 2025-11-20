# ACL Resources

`verity_acl_v4` and `verity_acl_v6` manage ACL (Access Control List) IP filter resources in Verity, which define network filtering rules for IPv4 and IPv6 traffic respectively.

## Example Usage

```hcl
resource "verity_acl_v4" "example" {
  name = "example"
  enable = false
  protocol = ""
  bidirectional = false
  source_ip = ""
  source_port_operator = ""
  source_port_1 = null
  source_port_2 = null
  destination_ip = ""
  destination_port_operator = ""
  destination_port_1 = null
  destination_port_2 = null

  object_properties {
    notes = ""
  }
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique.
* `enable` (Boolean) - Enable object.
* `protocol` (String) - Value must be ip/tcp/udp/icmp or a number between 0 and 255 to match packets.  Value IP will match all IP protocols.
* `bidirectional` (Boolean) - If bidirectional is selected, packets will be selected that match the source filters in either the source or destination fields of the packet.
* `source_ip` (String) - This field matches the source IP address of an IPv4 packet.
* `source_port_operator` (String) - This field determines which match operation will be applied to TCP/UDP ports. The choices are equal, greater, less or range.
* `source_port_1` (Integer) - This field is used for equal, greater-than or less-than TCP/UDP port value in match operation. This field is also used for the lower value in the range port match operation.
* `source_port_2` (Integer) - This field will only be used in the range TCP/UDP port value match operation to define the top value in the range.
* `destination_ip` (String) - This field matches the destination IP address of an IPv4 packet.
* `destination_port_operator` (String) - This field determines which match operation will be applied to TCP/UDP ports. The choices are equal, greater, less or range.
* `destination_port_1` (Integer) - This field is used for equal, greater-than or less-than TCP/UDP port value in match operation. This field is also used for the lower value in the range port match operation.
* `destination_port_2` (Integer) - This field will only be used in the range TCP/UDP port value match operation to define the top value in the range.
* `object_properties` (Object) - 
  * `notes` (String) - User Notes.

## Import

ACL IP Filter resources can be imported using the `name` attribute:

```sh
terraform import verity_acl_v4.<resource_name> <name>
terraform import verity_acl_v6.<resource_name> <name>
```
