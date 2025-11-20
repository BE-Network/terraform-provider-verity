# Voice Port Profile Resource

`verity_voice_port_profile` manages voice port profile resources in Verity, which define configurations for voice ports including protocol settings, call features, and quality parameters.

## Example Usage

```hcl
resource "verity_voice_port_profile" "example" {
  name = "example"
  enable = true
  protocol = "sip"
  digit_map = "(x.)"
  call_three_way_enable = true
  caller_id_enable = true
  caller_id_name_enable = true
  call_waiting_enable = true
  call_forward_unconditional_enable = false
  call_forward_on_busy_enable = true
  call_forward_on_no_answer_ring_count = 4
  call_transfer_enable = true
  audio_mwi_enable = false
  anonymous_call_block_enable = true
  do_not_disturb_enable = true
  cid_blocking_enable = false
  cid_num_presentation_status = "allowed"
  cid_name_presentation_status = "allowed"
  call_waiting_caller_id_enable = true
  call_hold_enable = true
  visual_mwi_enable = false
  mwi_refresh_timer = 3600
  hotline_enable = false
  dial_tone_feature_delay = 500
  intercom_enable = true
  intercom_transfer_enable = true
  jitter_target = 60
  jitter_buffer_max = 200
  release_timer = 20
  roh_timer = 30

  object_properties {
    isdefault = false
    port_monitoring = "enabled"
    group = "voice-group"
    format_dial_plan = true
  }
}
```

## Argument Reference

* `name` (String) - Unique identifier for the voice port profile.
* `enable` (Boolean) - Enable this voice port profile. Default is `false`.
* `protocol` (String) - Voice protocol (e.g., `sip`, `mgcp`).
* `digit_map` (String) - Digit map pattern for call routing.
* `call_three_way_enable` (Boolean) - Enable three-way calling. Default is `false`.
* `caller_id_enable` (Boolean) - Enable caller ID. Default is `false`.
* `caller_id_name_enable` (Boolean) - Enable caller ID name. Default is `false`.
* `call_waiting_enable` (Boolean) - Enable call waiting. Default is `false`.
* `call_forward_unconditional_enable` (Boolean) - Enable unconditional call forwarding. Default is `false`.
* `call_forward_on_busy_enable` (Boolean) - Enable call forwarding on busy. Default is `false`.
* `call_forward_on_no_answer_ring_count` (Integer) - Number of rings before forwarding on no answer.
* `call_transfer_enable` (Boolean) - Enable call transfer. Default is `false`.
* `audio_mwi_enable` (Boolean) - Enable audio message waiting indicator. Default is `false`.
* `anonymous_call_block_enable` (Boolean) - Block anonymous calls. Default is `false`.
* `do_not_disturb_enable` (Boolean) - Enable do not disturb. Default is `false`.
* `cid_blocking_enable` (Boolean) - Enable caller ID blocking. Default is `false`.
* `cid_num_presentation_status` (String) - Caller ID number presentation status (e.g., `allowed`, `restricted`).
* `cid_name_presentation_status` (String) - Caller ID name presentation status.
* `call_waiting_caller_id_enable` (Boolean) - Enable caller ID during call waiting. Default is `false`.
* `call_hold_enable` (Boolean) - Enable call hold. Default is `false`.
* `visual_mwi_enable` (Boolean) - Enable visual message waiting indicator. Default is `false`.
* `mwi_refresh_timer` (Integer) - Message waiting indicator refresh time in seconds.
* `hotline_enable` (Boolean) - Enable hotline. Default is `false`.
* `dial_tone_feature_delay` (Integer) - Dial tone feature delay in milliseconds.
* `intercom_enable` (Boolean) - Enable intercom. Default is `false`.
* `intercom_transfer_enable` (Boolean) - Enable intercom transfer. Default is `false`.
* `transmit_gain` (Integer) - Transmit gain in tenths of a dB. Example -30 would equal -3.0dB.
* `receive_gain` (Integer) - Receive gain in tenths of a dB. Example -30 would equal -3.0dB.
* `echo_cancellation_enable` (Boolean) - Enable echo cancellation. Default is `false`.
* `signaling_code` (String) - Signaling code for voice communication.
* `jitter_target` (Integer) - Target jitter buffer size in milliseconds.
* `jitter_buffer_max` (Integer) - Maximum jitter buffer size in milliseconds.
* `release_timer` (Integer) - Release timer in seconds.
* `roh_timer` (Integer) - Receiver off-hook timer in seconds.
* `object_properties` (Object) - Object properties configuration:
  * `port_monitoring` (String) - Port monitoring status.
  * `group` (String) - Group name.
  * `format_dial_plan` (Boolean) - Format dial plan. Default is `false`.

## Import

Voice Port Profile resources can be imported using the `name` attribute:

```sh
terraform import verity_voice_port_profile.<resource_name> <name>
```
