# verity_acl_v6 (Resource)

Manages a Verity IPv6 IP Filter (ACL).

Supported modes: Campus, Datacenter.

## Example Usage

```hcl
resource "verity_acl_v6" "example" {
  name = "example"
  bidirectional = false
  destination_ip = ""
  destination_port_1 = null
  destination_port_2 = null
  destination_port_operator = ""
  enable = false
  protocol = ""
  source_ip = ""
  source_port_1 = null
  source_port_2 = null
  source_port_operator = ""

  object_properties {
    notes = ""
  }
}
```

## Argument Reference

### Required

* `name` (String) - Object Name. Must be unique. Changing it replaces the resource.

### Optional

* `bidirectional` (Boolean) - If bidirectional is selected, packets will be selected that match the source filters in either the source or destination fields of the packet.
* `destination_ip` (String) - This field matches the destination IP address of an IPv6 packet.
* `destination_port_1` (Integer) - This field is used for equal, greater-than or less-than TCP/UDP port value in match operation. This field is also used for the lower value in the range port match operation. Set it to `null` to clear it.
* `destination_port_2` (Integer) - This field will only be used in the range TCP/UDP port value match operation to define the top value in the range. Set it to `null` to clear it.
* `destination_port_operator` (String) - This field determines which match operation will be applied to TCP/UDP ports. The choices are equal, greater-than, less-than or range.
* `enable` (Boolean) - Enable object.
* `object_properties` (Block) - Additional properties for this object. At most one block.
  * `notes` (String) - User notes.
* `protocol` (String) - Value must be ip/tcp/udp/icmp or a number between 0 and 255 to match packets. Value IP will match all IP protocols.
* `source_ip` (String) - This field matches the source IP address of an IPv6 packet.
* `source_port_1` (Integer) - This field is used for equal, greater-than or less-than TCP/UDP port value in match operation. This field is also used for the lower value in the range port match operation. Set it to `null` to clear it.
* `source_port_2` (Integer) - This field will only be used in the range TCP/UDP port value match operation to define the top value in the range. Set it to `null` to clear it.
* `source_port_operator` (String) - This field determines which match operation will be applied to TCP/UDP ports. The choices are equal, greater-than, less-than or range.

## Import

Import an existing object by its `name`:

```sh
terraform import verity_acl_v6.<resource_name> <name>
```
