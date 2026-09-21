
resource "verity_bundle" "bundle_test_script1" {
	object_properties {
		group = "test"
		is_for_switch = true
		is_public = true
	}
	enable = false
	user_services {
		row_app_enable = true
		index = 1
	}
}

resource "verity_bundle" "bundle_test_script2" {
	object_properties {
		group = "test"
		is_for_switch = true
		is_public = true
	}
	enable = false
	user_services {
		row_app_enable = true
		index = 1
	}
}