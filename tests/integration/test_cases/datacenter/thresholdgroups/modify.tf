
resource "verity_threshold_group" "threshold_group_test_script1" {
	enable = false
	targets {
		index = 1
		enable = false
	}
	thresholds {
		index = 2
	}
}

resource "verity_threshold_group" "threshold_group_test_script2" {
	enable = false
	targets {
		index = 1
		enable = false
	}
	thresholds {
		index = 1
		enable = false
	}
}