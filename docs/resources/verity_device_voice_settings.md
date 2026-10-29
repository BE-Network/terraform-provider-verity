# verity_device_voice_settings (Resource)

Manages a Verity Device Voice Settings.

Supported modes: Campus.

## Example Usage

```hcl
resource "verity_device_voice_settings" "example" {
  name = "example"
  anon_cid_block_activate = ""
  anon_cid_block_deactivate = ""
  bit_rate = ""
  call_agent_1 = ""
  call_agent_2 = ""
  call_agent_port_1 = null
  call_agent_port_2 = null
  call_forward_on_busy_activate = ""
  call_forward_on_busy_deactivate = ""
  call_forward_on_no_answer_activate = ""
  call_forward_on_no_answer_deactivate = ""
  call_forward_unconditional_activate = ""
  call_forward_unconditional_deactivate = ""
  call_hold = ""
  cancel_call_waiting = ""
  cas_events = null
  cids_activate = ""
  cids_deactivate = ""
  do_not_disturb_activate = ""
  do_not_disturb_deactivate = ""
  do_not_disturb_pin_change = ""
  domain = ""
  dscp_mark = null
  dtmf_method = ""
  emergency_service_number = ""
  enable = false
  event_payload_type = null
  fax_t38 = false
  intercom_1 = ""
  intercom_2 = ""
  intercom_3 = ""
  local_port_max = null
  local_port_min = null
  mgcp_dscp_mark = null
  outbound_proxy = ""
  outbound_proxy_port = null
  outbound_proxy_secondary = ""
  outbound_proxy_secondary_port = null
  protocol = ""
  proxy_server = ""
  proxy_server_port = null
  proxy_server_secondary = ""
  proxy_server_secondary_port = null
  region = ""
  register_expires = null
  registrar_server = ""
  registrar_server_port = null
  registrar_server_secondary = ""
  registrar_server_secondary_port = null
  registration_period = null
  rtcp = false
  sip_dscp_mark = null
  termination_base = ""
  user_agent_domain = ""
  user_agent_port = null
  user_agent_transport = ""
  voicemail_server = ""
  voicemail_server_expires = null
  voicemail_server_port = null

  codecs {
    index = 1
    codec_num_enable = false
    codec_num_name = ""
    codec_num_packetization_period = ""
    codec_num_silence_suppression = false
  }
}
```

## Argument Reference

### Required

* `name` (String) - Template Name. Must be unique within type. Changing it replaces the resource.

### Optional

* `anon_cid_block_activate` (String) - Anonymoes Caller ID Block Activate.
* `anon_cid_block_deactivate` (String) - Anonymous Caller ID Block Deactivate.
* `bit_rate` (String) - T.38 Bit Rate in bps. Most available fax machines support up to 14,400bps.
* `call_agent_1` (String) - Call Agent 1.
* `call_agent_2` (String) - Call Agent 2.
* `call_agent_port_1` (Integer) - Call Agent Port 1. Set it to `null` to clear it.
* `call_agent_port_2` (Integer) - Call Agent Port 2. Set it to `null` to clear it.
* `call_forward_on_busy_activate` (String) - Call Forward On Busy Activate.
* `call_forward_on_busy_deactivate` (String) - Call Forward On Busy Deactivate.
* `call_forward_on_no_answer_activate` (String) - Call Forward On No Answer Activate.
* `call_forward_on_no_answer_deactivate` (String) - Call Forward On No Answer Deactivate.
* `call_forward_unconditional_activate` (String) - Call Forward Unconditional Activate.
* `call_forward_unconditional_deactivate` (String) - Call Forward Unconditional Deactivate.
* `call_hold` (String) - Call hold.
* `cancel_call_waiting` (String) - Cancel Call waiting.
* `cas_events` (Integer) - Enables or disables handling of CAS via RTP CAS events. Valid values are 0 = off and 1 = on. Set it to `null` to clear it.
* `cids_activate` (String) - Caller ID Delivery Blocking (single call) Activate.
* `cids_deactivate` (String) - Caller ID Delivery Blocking (single call) Deactivate.
* `codecs` (Block List) - Codec configurations. Entries are matched by `index`.
  * `codec_num_enable` (Boolean) - Enable Codec.
  * `codec_num_name` (String) - Name of this Codec.
  * `codec_num_packetization_period` (String) - Packet period selection interval in milliseconds.
  * `codec_num_silence_suppression` (Boolean) - Specifies whether silence suppression is on or off. Valid values are 0 = off and 1 = on.
  * `index` (Integer) - The index identifying the object. Zero if you want to add an object to the list.
* `do_not_disturb_activate` (String) - Do not Disturb Activate.
* `do_not_disturb_deactivate` (String) - Do not Disturb Deactivate.
* `do_not_disturb_pin_change` (String) - Do not Disturb PIN Change.
* `domain` (String) - Domain.
* `dscp_mark` (Integer) - Differentiated Services Code Point (DSCP) to be used for outgoing RTP packets. Set it to `null` to clear it.
* `dtmf_method` (String) - Specifies how DTMF signals are carried.
* `emergency_service_number` (String) - Emergency Service Number.
* `enable` (Boolean) - Enable object.
* `event_payload_type` (Integer) - Telephone Event Payload Type. Set it to `null` to clear it.
* `fax_t38` (Boolean) - Fax T.38 Enable.
* `intercom_1` (String) - Intercom 1.
* `intercom_2` (String) - Intercom 2.
* `intercom_3` (String) - Intercom 3.
* `local_port_max` (Integer) - Defines the highest RTP port used for voice traffic, must be greater than local Local Port Min. Set it to `null` to clear it.
* `local_port_min` (Integer) - Defines the base RTP port that should be used for voice traffic. Set it to `null` to clear it.
* `mgcp_dscp_mark` (Integer) - MGCP Differentiated Services Code point (DSCP). Set it to `null` to clear it.
* `outbound_proxy` (String) - IP address or URI of the outbound proxy server for SIP signalling messages. An outbound SIP proxy may or may not be required within a given network.
* `outbound_proxy_port` (Integer) - Outbound Proxy Port. Set it to `null` to clear it.
* `outbound_proxy_secondary` (String) - IP address or URI of the secondary outbound proxy server for SIP signalling messages. An outbound SIP proxy may or may not be required within a given network.
* `outbound_proxy_secondary_port` (Integer) - Secondary Outbound Proxy Port. Set it to `null` to clear it.
* `protocol` (String) - Voice Protocol: MGCP or SIP.
* `proxy_server` (String) - IP address or URI of the SIP proxy server for SIP signalling messages.
* `proxy_server_port` (Integer) - Proxy Server Port. Set it to `null` to clear it.
* `proxy_server_secondary` (String) - IP address or URI of the secondary SIP proxy server for SIP signalling messages.
* `proxy_server_secondary_port` (Integer) - Secondary Proxy Server Port. Set it to `null` to clear it.
* `region` (String) - Region.
* `register_expires` (Integer) - SIP registration expiration time in seconds. If value is 0, the SIP agent does not add an expiration time to the registration requests and does not perform re-registration. The default value is 3600 seconds. Set it to `null` to clear it.
* `registrar_server` (String) - Name or IP address or resolved name of the registrar server for SIP signalling messages. Examples: 10.10.10.10 and proxy.voip.net.
* `registrar_server_port` (Integer) - Registrar Server Port. Set it to `null` to clear it.
* `registrar_server_secondary` (String) - Name or IP address or resolved name of the secondary registrar server for SIP signalling messages. Examples: 10.10.10.10 and proxy.voip.net.
* `registrar_server_secondary_port` (Integer) - Secondary Registrar Server Port. Set it to `null` to clear it.
* `registration_period` (Integer) - Specifies the time in seconds to start the re-registration process. The default value is 3240 seconds. Set it to `null` to clear it.
* `rtcp` (Boolean) - RTCP Enable.
* `sip_dscp_mark` (Integer) - Sip Differentiated Services Code point (DSCP). Set it to `null` to clear it.
* `termination_base` (String) - Base string for the MGCP physical termination id(s).
* `user_agent_domain` (String) - User Agent Domain.
* `user_agent_port` (Integer) - User Agent Port. Set it to `null` to clear it.
* `user_agent_transport` (String) - User Agent Transport.
* `voicemail_server` (String) - Name or IP address or resolved name of the external voicemail server if not provided by SIP server for MWI control. Examples: 10.10.10.10 and proxy.voip.net.
* `voicemail_server_expires` (Integer) - Voicemail server expiration time in seconds. If value is 0, the Register Expires time is used instead. The default value is 3600 seconds. Set it to `null` to clear it.
* `voicemail_server_port` (Integer) - Voicemail Server Port. Set it to `null` to clear it.

## Import

Import an existing object by its `name`:

```sh
terraform import verity_device_voice_settings.<resource_name> <name>
```
