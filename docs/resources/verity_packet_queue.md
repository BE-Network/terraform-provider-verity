# verity_packet_queue (Resource)

Manages a Verity Packet Queue.

Supported modes: Campus, Datacenter.

## Example Usage

```hcl
resource "verity_packet_queue" "example" {
  name = "example"
  enable = false

  pbit {
    index = 1
    packet_queue_for_p_bit = null
  }

  queue {
    index = 1
    bandwidth_for_queue = null
    scheduler_type = ""
    scheduler_weight = null
  }
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `enable` (Boolean) - Enable object.
* `pbit` (Block List) - P-bit configurations. Entries are matched by `index`.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `packet_queue_for_p_bit` (Integer) - Flag indicating this Traffic Class' Queue. Set it to `null` to clear it.
* `queue` (Block List) - Queue configurations. Entries are matched by `index`.
  * `bandwidth_for_queue` (Integer) - Percentage bandwidth allocated to Queue. 0 is no limit. Set it to `null` to clear it.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
  * `scheduler_type` (String) - Scheduler Type for Queue.
  * `scheduler_weight` (Integer) - Weight associated with WRR or DWRR scheduler. Set it to `null` to clear it.

## Import

Import an existing object by its `name`:

```sh
terraform import verity_packet_queue.<resource_name> <name>
```
