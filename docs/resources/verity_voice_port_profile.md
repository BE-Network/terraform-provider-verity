# verity_voice_port_profile (Resource)

Manages a Verity Voice Port Profile.

Supported modes: Campus.

## Example Usage

```hcl
resource "verity_voice_port_profile" "example" {
  name = "example"
  anonymous_call_block_enable = false
  audio_mwi_enable = false
  call_forward_on_busy_enable = false
  call_forward_on_no_answer_ring_count = null
  call_forward_unconditional_enable = false
  call_hold_enable = false
  call_three_way_enable = false
  call_transfer_enable = false
  call_waiting_caller_id_enable = false
  call_waiting_enable = false
  caller_id_enable = false
  caller_id_name_enable = false
  cid_blocking_enable = false
  cid_name_presentation_status = ""
  cid_num_presentation_status = ""
  dial_tone_feature_delay = null
  digit_map = ""
  do_not_disturb_enable = false
  echo_cancellation_enable = false
  enable = false
  hotline_enable = false
  intercom_enable = false
  intercom_transfer_enable = false
  jitter_buffer_max = null
  jitter_target = null
  mwi_refresh_timer = null
  protocol = ""
  receive_gain = null
  release_timer = null
  roh_timer = null
  signaling_code = ""
  transmit_gain = null
  visual_mwi_enable = false

  object_properties {
    format_dial_plan = false
    port_monitoring = ""
  }
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `anonymous_call_block_enable` (Boolean) - Block all anonymous calls.
* `audio_mwi_enable` (Boolean) - Audio Message Waiting Indicator.
* `call_forward_on_busy_enable` (Boolean) - Call Forward On Busy.
* `call_forward_on_no_answer_ring_count` (Integer) - Call Forward on number of rings. Set it to `null` to clear it.
* `call_forward_unconditional_enable` (Boolean) - Call Forward Unconditional.
* `call_hold_enable` (Boolean) - Call Hold.
* `call_three_way_enable` (Boolean) - Enable three way calling.
* `call_transfer_enable` (Boolean) - Call Transfer.
* `call_waiting_caller_id_enable` (Boolean) - Call Waiting Caller ID.
* `call_waiting_enable` (Boolean) - Call Waiting.
* `caller_id_enable` (Boolean) - Caller ID.
* `caller_id_name_enable` (Boolean) - Caller ID Name.
* `cid_blocking_enable` (Boolean) - CID Blocking.
* `cid_name_presentation_status` (String) - CID Name Presentation.
* `cid_num_presentation_status` (String) - CID Number Presentation.
* `dial_tone_feature_delay` (Integer) - Dial Tone Feature Delay. Set it to `null` to clear it.
* `digit_map` (String) - Dial Plan.
* `do_not_disturb_enable` (Boolean) - Do not disturb.
* `echo_cancellation_enable` (Boolean) - Echo Cancellation Enable.
* `enable` (Boolean) - Enable object.
* `hotline_enable` (Boolean) - Direct Connect.
* `intercom_enable` (Boolean) - Intercom.
* `intercom_transfer_enable` (Boolean) - Intercom Transfer.
* `jitter_buffer_max` (Integer) - The maximum depth of the jitter buffer in milliseconds. Set it to `null` to clear it.
* `jitter_target` (Integer) - The target value of the jitter buffer in milliseconds. Set it to `null` to clear it.
* `mwi_refresh_timer` (Integer) - Message Waiting Indicator Refresh. Set it to `null` to clear it.
* `object_properties` (Block) - Object properties for the voice port profile. At most one block.
  * `format_dial_plan` (Boolean) - Format dial plan for easier viewing.
  * `port_monitoring` (String) - Defines importance of Link Down on this port.
* `protocol` (String) - Voice Protocol: MGCP or SIP.
* `receive_gain` (Integer) - Receive Gainin tenths of a dB. Example -30 would equal -3.0db. Set it to `null` to clear it.
* `release_timer` (Integer) - Release timer defined in seconds. The default value of this attribute is 10 seconds. Set it to `null` to clear it.
* `roh_timer` (Integer) - Time in seconds for the receiver is off-hook before ROH tone is applied. The value 0 disables ROH timing. The default value is 15 seconds. Set it to `null` to clear it.
* `signaling_code` (String) - Signaling Code.
* `transmit_gain` (Integer) - Transmit Gain in tenths of a dB.Example -30 would equal -3.0db. Set it to `null` to clear it.
* `visual_mwi_enable` (Boolean) - Visual Message Waiting Indicator.

## Import

Import an existing object by its `name`:

```sh
terraform import verity_voice_port_profile.<resource_name> <name>
```
