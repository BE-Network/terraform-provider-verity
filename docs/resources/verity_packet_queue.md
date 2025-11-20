# Packet Queue Resource

`verity_packet_queue` manages packet queue resources in Verity, which define quality of service settings for network traffic.

## Example Usage

```hcl
resource "verity_packet_queue" "example" {
  name = "example"
  enable = true
  
  pbit {
    packet_queue_for_p_bit = 1
    index = 1
  }
  
  queue {
    bandwidth_for_queue = 100
    scheduler_type = "SP"
    scheduler_weight = 10
    index = 1
  }
  
  object_properties {
    group = "queue-group"
  }
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique.
* `enable` (Boolean) - Enable object.
* `pbit` (Array) - 
  * `packet_queue_for_p_bit` (Integer) - Flag indicating this Traffic Class' Queue.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `queue` (Array) - 
  * `bandwidth_for_queue` (Integer) - Percentage bandwidth allocated to Queue. 0 is no limit.
  * `scheduler_type` (String) - Scheduler Type for Queue.
  * `scheduler_weight` (Integer) - Weight associated with WRR or DWRR scheduler.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `object_properties` (Object) - 
  * `group` (String) - Group.

## Import

Packet Queue resources can be imported using the `name` attribute:

```sh
terraform import verity_packet_queue.<resource_name> <name>
```
