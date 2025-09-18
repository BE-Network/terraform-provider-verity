# Voice Port Profile Resource

`verity_voice_port_profile` manages voice port profile resources in Verity, which define configurations for voice ports including protocol settings, call features, and quality parameters.

## Version Compatibility

**This resource requires Verity API version 6.5 or higher.**

## Example Usage

```hcl
resource "verity_voice_port_profile" "example" {
  name = "example_voice_profile"
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
  modem_transport = "passthrough"
  preferred_codec = "g711ulaw"
  fax_tone_detect_enable = true
  silent_dialing_enable = false
  dtmf_inband = true
  dtmf_rfc2833 = true
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

* `name` - (Required) Unique identifier for the voice port profile.
* `enable` - (Optional) Enable this voice port profile. Default is `false`.
* `protocol` - (Optional) Voice protocol (e.g., `sip`, `mgcp`).
* `digit_map` - (Optional) Digit map pattern for call routing.
* `call_three_way_enable` - (Optional) Enable three-way calling. Default is `false`.
* `caller_id_enable` - (Optional) Enable caller ID. Default is `false`.
* `caller_id_name_enable` - (Optional) Enable caller ID name. Default is `false`.
* `call_waiting_enable` - (Optional) Enable call waiting. Default is `false`.
* `call_forward_unconditional_enable` - (Optional) Enable unconditional call forwarding. Default is `false`.
* `call_forward_on_busy_enable` - (Optional) Enable call forwarding on busy. Default is `false`.
* `call_forward_on_no_answer_ring_count` - (Optional) Number of rings before forwarding on no answer.
* `call_transfer_enable` - (Optional) Enable call transfer. Default is `false`.
* `audio_mwi_enable` - (Optional) Enable audio message waiting indicator. Default is `false`.
* `anonymous_call_block_enable` - (Optional) Block anonymous calls. Default is `false`.
* `do_not_disturb_enable` - (Optional) Enable do not disturb. Default is `false`.
* `cid_blocking_enable` - (Optional) Enable caller ID blocking. Default is `false`.
* `cid_num_presentation_status` - (Optional) Caller ID number presentation status (e.g., `allowed`, `restricted`).
* `cid_name_presentation_status` - (Optional) Caller ID name presentation status.
* `call_waiting_caller_id_enable` - (Optional) Enable caller ID during call waiting. Default is `false`.
* `call_hold_enable` - (Optional) Enable call hold. Default is `false`.
* `visual_mwi_enable` - (Optional) Enable visual message waiting indicator. Default is `false`.
* `mwi_refresh_timer` - (Optional) Message waiting indicator refresh time in seconds.
* `hotline_enable` - (Optional) Enable hotline. Default is `false`.
* `dial_tone_feature_delay` - (Optional) Dial tone feature delay in milliseconds.
* `intercom_enable` - (Optional) Enable intercom. Default is `false`.
* `intercom_transfer_enable` - (Optional) Enable intercom transfer. Default is `false`.
* `transmit_gain` - (Optional) Transmit gain in tenths of a dB. Example -30 would equal -3.0dB (available as of API version 6.5).
* `receive_gain` - (Optional) Receive gain in tenths of a dB. Example -30 would equal -3.0dB (available as of API version 6.5).
* `echo_cancellation_enable` - (Optional) Enable echo cancellation. Default is `false` (available as of API version 6.5).
* `modem_transport` - (Optional) Modem transport mode (e.g., `passthrough`, `relay`).
* `preferred_codec` - (Optional) Preferred audio codec (e.g., `g711ulaw`, `g729`).
* `fax_tone_detect_enable` - (Optional) Enable fax tone detection. Default is `false`.
* `silent_dialing_enable` - (Optional) Enable silent dialing. Default is `false`.
* `dtmf_inband` - (Optional) Enable DTMF in-band. Default is `false`.
* `dtmf_rfc2833` - (Optional) Enable DTMF RFC 2833. Default is `false`.
* `signaling_code` - (Optional) Signaling code for voice communication (available as of API version 6.5).
* `jitter_target` - (Optional) Target jitter buffer size in milliseconds.
* `jitter_buffer_max` - (Optional) Maximum jitter buffer size in milliseconds.
* `release_timer` - (Optional) Release timer in seconds.
* `roh_timer` - (Optional) Receiver off-hook timer in seconds.
* `object_properties` - (Optional) Object properties configuration:
  * `isdefault` - (Optional) Whether this is the default profile. Default is `false`.
  * `port_monitoring` - (Optional) Port monitoring status.
  * `group` - (Optional) Group name.
  * `format_dial_plan` - (Optional) Format dial plan. Default is `false`.

## Import

Voice Port Profile resources can be imported using the `name` attribute:

```sh
terraform import verity_voice_port_profile.<resource_name> <name>
```
