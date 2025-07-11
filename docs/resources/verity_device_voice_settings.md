# Device Voice Settings Resource

`verity_device_voice_settings` manages voice settings for devices in Verity, which define the configuration for voice communications.

## Version Compatibility

**This resource requires Verity API version 6.5 or higher.**

## Example Usage

```hcl
resource "verity_device_voice_settings" "example" {
  name = "example_voice_settings"
  enable = true
  dtmf_method = "inband"
  region = "US"
  protocol = "sip"
  proxy_server = "proxy.example.com"
  proxy_server_port = 5060
  registrar_server = "registrar.example.com"
  registrar_server_port = 5060
  user_agent_domain = "example.com"
  user_agent_transport = "udp"
  user_agent_port = 5060
  registration_period = 3600
  
  codecs {
    codec_num_name = "g711ulaw"
    codec_num_enable = true
    codec_num_packetization_period = "20msec"
    codec_num_silence_suppression = false
    index = 1
  }
  
  object_properties {
    isdefault = false
    group = "voice-group"
  }
}
```

## Argument Reference

* `name` - (Required) Unique identifier for the device voice settings.
* `enable` - (Optional) Enable these voice settings. Default is `false`.
* `dtmf_method` - (Optional) DTMF method (e.g., `inband`, `rfc2833`, `sipinfo`).
* `region` - (Optional) Region code for the voice settings.
* `protocol` - (Optional) Voice protocol (e.g., `sip`).
* `proxy_server` - (Optional) SIP proxy server hostname or IP address.
* `proxy_server_port` - (Optional) SIP proxy server port number.
* `proxy_server_secondary` - (Optional) Secondary SIP proxy server hostname or IP address.
* `proxy_server_secondary_port` - (Optional) Secondary SIP proxy server port number.
* `registrar_server` - (Optional) SIP registrar server hostname or IP address.
* `registrar_server_port` - (Optional) SIP registrar server port number.
* `registrar_server_secondary` - (Optional) Secondary SIP registrar server hostname or IP address.
* `registrar_server_secondary_port` - (Optional) Secondary SIP registrar server port number.
* `user_agent_domain` - (Optional) SIP user agent domain.
* `user_agent_transport` - (Optional) SIP user agent transport protocol (e.g., `udp`, `tcp`).
* `user_agent_port` - (Optional) SIP user agent port number.
* `outbound_proxy` - (Optional) SIP outbound proxy hostname or IP address.
* `outbound_proxy_port` - (Optional) SIP outbound proxy port number.
* `outbound_proxy_secondary` - (Optional) Secondary SIP outbound proxy hostname or IP address.
* `outbound_proxy_secondary_port` - (Optional) Secondary SIP outbound proxy port number.
* `registration_period` - (Optional) Registration period in seconds.
* `register_expires` - (Optional) Register expiration time in seconds.
* `voicemail_server` - (Optional) Voicemail server hostname or IP address.
* `voicemail_server_port` - (Optional) Voicemail server port number.
* `voicemail_server_expires` - (Optional) Voicemail server expiration time in seconds.
* `codecs` - (Optional) List of codec configurations:
  * `codec_num_name` - (Optional) Codec name (e.g., `g711ulaw`, `g711alaw`, `g729`).
  * `codec_num_enable` - (Optional) Enable this codec. Default is `false`.
  * `codec_num_packetization_period` - (Optional) Packetization period (e.g., `20msec`, `30msec`).
  * `codec_num_silence_suppression` - (Optional) Enable silence suppression. Default is `false`.
  * `index` - (Optional) Index identifying this codec configuration.
* `object_properties` - (Optional) Object properties configuration:
  * `isdefault` - (Optional) Whether this is the default voice setting. Default is `false`.
  * `group` - (Optional) Group name.

## Import

Device Voice Settings resources can be imported using the `name` attribute:

```
$ terraform import verity_device_voice_settings.example example
```
