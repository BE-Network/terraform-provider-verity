
resource "verity_threshold" "threshold_test_script1" {
	enable = false
	rules {
		index = 1
		enable = false
	}
	rules {
		index = 2
	}
}

resource "verity_threshold" "threshold_test_script2" {
	enable = false
	rules {
		index = 1
		enable = false
	}
}