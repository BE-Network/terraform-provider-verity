# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_device_voice_settings" "device_voice_settings_test_script1" {
	object_properties {
		group = ""
	}
	codecs {
		index = 3
		codec_num_enable = false
		codec_num_silence_suppression = true
	}
	codecs {
		index = 4
		codec_num_enable = false
	}
	codecs {
		index = 5
	}
	enable = false
	fax_t38 = true
	local_port_max = 30202
	local_port_min = 30002
	register_expires = 3660
	registration_period = 3245
	rtcp = false
	voicemail_server_expires = 3660
}


resource "verity_device_voice_settings" "device_voice_settings_test_script2" {
	object_properties {
		group = "test1"
	}
	codecs {
		index = 2
		codec_num_enable = false
		codec_num_silence_suppression = true
	}
	codecs {
		index = 3
		codec_num_enable = false
	}
	codecs {
		index = 4
    }
	enable = false
	fax_t38 = true
	local_port_max = 30202
	local_port_min = 30002
	register_expires = 3660
	registration_period = 3242
	rtcp = false
	voicemail_server_expires = 3660
}
