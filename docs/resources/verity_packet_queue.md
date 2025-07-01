# Packet Queue Resource

`verity_packet_queue` manages packet queue resources in Verity, which define quality of service settings for network traffic.

## Example Usage

```hcl
resource "verity_packet_queue" "example" {
  name = "example_queue"
  enable = true
  
  pbit {
    packet_queue_for_p_bit = 1
    index = 0
  }
  
  pbit {
    packet_queue_for_p_bit = 2
    index = 1
  }
  
  queue {
    bandwidth_for_queue = 1000
    scheduler_type = "strict"
    scheduler_weight = 10
    index = 0
  }
  
  queue {
    bandwidth_for_queue = 2000
    scheduler_type = "wrr"
    scheduler_weight = 20
    index = 1
  }
  
  object_properties {
    isdefault = false
    group = "queue-group"
  }
}
```

## Argument Reference

* `name` - (Required) Unique identifier for the packet queue.
* `enable` - (Optional) Enable this packet queue. Default is `false`.
* `pbit` - (Optional) List of P-bit configurations:
  * `packet_queue_for_p_bit` - (Optional) Queue assignment for this P-bit.
  * `index` - (Optional) Index identifying this P-bit configuration.
* `queue` - (Optional) List of queue configurations:
  * `bandwidth_for_queue` - (Optional) Bandwidth allocated to this queue in Kbps.
  * `scheduler_type` - (Optional) Scheduler type (e.g., `strict`, `wrr`).
  * `scheduler_weight` - (Optional) Scheduler weight for weighted scheduling algorithms.
  * `index` - (Optional) Index identifying this queue configuration.
* `object_properties` - (Optional) Object properties configuration:
  * `isdefault` - (Optional) Whether this is the default queue configuration. Default is `false`.
  * `group` - (Optional) Group name.

## Import

Packet Queue resources can be imported using the `name` attribute:

```
$ terraform import verity_packet_queue.example example_queue
```
