# Device Voice Settings Resource

`verity_device_voice_settings` manages voice settings for devices in Verity, which define the configuration for voice communications.

## Example Usage

```hcl
resource "verity_device_voice_settings" "example" {
  name = "example"
  enable = false
  dtmf_method = "Inband"
  region = "US"
  protocol = "SIP"
  proxy_server = "example.com"
  proxy_server_port = 5060
  proxy_server_secondary = ""
  proxy_server_secondary_port = 0
  registrar_server = "example.com"
  registrar_server_port = 5060
  registrar_server_secondary = ""
  registrar_server_secondary_port = 0
  user_agent_domain = "example.com"
  user_agent_transport = "UDP"
  user_agent_port = 5060
  outbound_proxy = ""
  outbound_proxy_port = 0
  outbound_proxy_secondary = ""
  outbound_proxy_secondary_port = 0
  registration_period = 3240
  register_expires = 3600
  voicemail_server = ""
  voicemail_server_port = 0
  voicemail_server_expires = 3600
  sip_dscp_mark = 0
  call_agent_1 = ""
  call_agent_port_1 = 0
  call_agent_2 = ""
  call_agent_port_2 = 0
  domain = ""
  mgcp_dscp_mark = 0
  termination_base = "aaln/"
  local_port_min = 30000
  local_port_max = 30200
  event_payload_type = 101
  cas_events = 0
  dscp_mark = 0
  rtcp = true
  fax_t38 = false
  bit_rate = "14400"
  cancel_call_waiting = "*70"
  call_hold = "*9"
  cids_activate = "*67"
  cids_deactivate = "*82"
  do_not_disturb_activate = "*78"
  do_not_disturb_deactivate = "*79"
  do_not_disturb_pin_change = "*10"
  emergency_service_number = "911"
  anon_cid_block_activate = "*77"
  anon_cid_block_deactivate = "*87"
  call_forward_unconditional_activate = "*72"
  call_forward_unconditional_deactivate = "*73"
  call_forward_on_busy_activate = "*90"
  call_forward_on_busy_deactivate = "*91"
  call_forward_on_no_answer_activate = "*92"
  call_forward_on_no_answer_deactivate = "*93"
  intercom_1 = "*53"
  intercom_2 = "*54"
  intercom_3 = "*55"

  codecs {
    codec_num_name = "a"
    codec_num_enable = true
    codec_num_packetization_period = "20"
    codec_num_silence_suppression = false
    index = 1
  }

  object_properties {
    group = "voice-group"
  }
}
```

## Argument Reference

* `name` (String) - Object Name. Must be unique.
* `enable` (Boolean) - Enable object.
* `dtmf_method` (String) - Specifies how DTMF signals are carried.
* `region` (String) - Region.
* `protocol` (String) - Voice Protocol: MGCP or SIP.
* `proxy_server` (String) - IP address or URI of the SIP proxy server for SIP signalling messages.
* `proxy_server_port` (Integer) - Proxy Server Port.
* `proxy_server_secondary` (String) - IP address or URI of the secondary SIP proxy server for SIP signalling messages.
* `proxy_server_secondary_port` (Integer) - Secondary Proxy Server Port.
* `registrar_server` (String) - Name or IP address or resolved name of the registrar server for SIP signalling messages.
* `registrar_server_port` (Integer) - Registrar Server Port.
* `registrar_server_secondary` (String) - Name or IP address or resolved name of the secondary registrar server for SIP signalling messages.
* `registrar_server_secondary_port` (Integer) - Secondary Registrar Server Port.
* `user_agent_domain` (String) - User Agent Domain.
* `user_agent_transport` (String) - User Agent Transport.
* `user_agent_port` (Integer) - User Agent Port.
* `outbound_proxy` (String) - IP address or URI of the outbound proxy server for SIP signalling messages. An outbound SIP proxy may or may not be required within a given network.
* `outbound_proxy_port` (Integer) - Outbound Proxy Port.
* `outbound_proxy_secondary` (String) - IP address or URI of the secondary outbound proxy server for SIP signalling messages. An outbound SIP proxy may or may not be required within a given network.
* `outbound_proxy_secondary_port` (Integer) - Secondary Outbound Proxy Port.
* `registration_period` (Integer) - Specifies the time in seconds to start the re-registration process.
* `register_expires` (Integer) - SIP registration expiration time in seconds. If value is 0, the SIP agent does not add an expiration time to the registration requests and does not perform re-registration.
* `voicemail_server` (String) - Name or IP address or resolved name of the external voicemail server if not provided by SIP server for MWI control.
* `voicemail_server_port` (Integer) - Voicemail Server Port.
* `voicemail_server_expires` (Integer) - Voicemail server expiration time in seconds. If value is 0, the Register Expires time is used instead.
* `sip_dscp_mark` (Integer) - Sip Differentiated Services Code point (DSCP).
* `call_agent_1` (String) - Call Agent 1.
* `call_agent_port_1` (Integer) - Call Agent Port 1.
* `call_agent_2` (String) - Call Agent 2.
* `call_agent_port_2` (Integer) - Call Agent Port 2.
* `domain` (String) - Domain.
* `mgcp_dscp_mark` (Integer) - MGCP Differentiated Services Code point (DSCP).
* `termination_base` (String) - Base string for the MGCP physical termination id(s).
* `local_port_min` (Integer) - Defines the base RTP port that should be used for voice traffic.
* `local_port_max` (Integer) - Defines the highest RTP port used for voice traffic, must be greater than local Local Port Min.
* `event_payload_type` (Integer) - Telephone Event Payload Type.
* `cas_events` (Integer) - Enables or disables handling of CAS via RTP CAS events. Valid values are 0 = off and 1 = on.
* `dscp_mark` (Integer) - Differentiated Services Code Point (DSCP) to be used for outgoing RTP packets.
* `rtcp` (Boolean) - RTCP Enable.
* `fax_t38` (Boolean) - Fax T.38 Enable.
* `bit_rate` (String) - T.38 Bit Rate in bps. Most available fax machines support up to 14,400bps.
* `cancel_call_waiting` (String) - Cancel Call waiting.
* `call_hold` (String) - Call hold.
* `cids_activate` (String) - Caller ID Delivery Blocking (single call) Activate.
* `cids_deactivate` (String) - Caller ID Delivery Blocking (single call) Deactivate.
* `do_not_disturb_activate` (String) - Do not Disturb Activate.
* `do_not_disturb_deactivate` (String) - Do not Disturb Deactivate.
* `do_not_disturb_pin_change` (String) - Do not Disturb PIN Change.
* `emergency_service_number` (String) - Emergency Service Number.
* `anon_cid_block_activate` (String) - Anonymoes Caller ID Block Activate.
* `anon_cid_block_deactivate` (String) - Anonymous Caller ID Block Deactivate.
* `call_forward_unconditional_activate` (String) - Call Forward Unconditional Activate.
* `call_forward_unconditional_deactivate` (String) - Call Forward Unconditional Deactivate.
* `call_forward_on_busy_activate` (String) - Call Forward On Busy Activate.
* `call_forward_on_busy_deactivate` (String) - Call Forward On Busy Deactivate.
* `call_forward_on_no_answer_activate` (String) - Call Forward On No Answer Activate.
* `call_forward_on_no_answer_deactivate` (String) - Call Forward On No Answer Deactivate.
* `intercom_1` (String) - Intercom 1.
* `intercom_2` (String) - Intercom 2.
* `intercom_3` (String) - Intercom 3.
* `codecs` (Array) - 
  * `codec_num_name` (String) - Name of this Codec.
  * `codec_num_enable` (Boolean) - Enable Codec.
  * `codec_num_packetization_period` (String) - Packet period selection interval in milliseconds.
  * `codec_num_silence_suppression` (Boolean) - Specifies whether silence suppression is on or off. Valid values are 0 = off and 1 = on.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `object_properties` (Object) - 
  * `group` (String) - Group.

## Import

Device Voice Settings resources can be imported using the `name` attribute:

```sh
terraform import verity_device_voice_settings.<resource_name> <name>
```
