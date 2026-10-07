# verity_nac_port_profile (Resource)

Manages a Verity NAC Port Profile.

Supported modes: Campus.

## Example Usage

```hcl
resource "verity_nac_port_profile" "example" {
  name = "example"
  enable = false

  eth_ports {
    index = 1
    eth_port_profile_num_enable = false
    eth_port_profile_num_eth_port = ""
    eth_port_profile_num_eth_port_ref_type_ = "eth_port_profile_"
  }

  object_properties {
    port_monitoring = ""
  }
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `enable` (Boolean) - Enable object.
* `eth_ports` (Block List) - Ethernet port profiles used for network access control. Entries are matched by `index`.
  * `eth_port_profile_num_enable` (Boolean) - Enable row.
  * `eth_port_profile_num_eth_port` (String) - Choose an Eth Port Profile. Set together with `eth_port_profile_num_eth_port_ref_type_`.
  * `eth_port_profile_num_eth_port_ref_type_` (String) - Object type for eth_port_profile_num_eth_port field.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `object_properties` (Block) - Object properties for the NAC port profile. At most one block.
  * `port_monitoring` (String) - Defines importance of Link Down on this port.

## Reference Fields

A reference names another Verity object. Set the reference and its type field together. When your configuration sets neither, Terraform leaves the reference to the server, including one set in the Verity UI.

| Field | Type field | Allowed types |
| --- | --- | --- |
| `eth_ports.eth_port_profile_num_eth_port` | `eth_port_profile_num_eth_port_ref_type_` | `eth_port_profile_` |

## Import

Import an existing object by its `name`:

```sh
terraform import verity_nac_port_profile.<resource_name> <name>
```
