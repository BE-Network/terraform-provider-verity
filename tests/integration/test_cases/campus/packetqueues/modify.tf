# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_packet_queue" "packet_queue_test_script1" {
	object_properties {
		group = ""
	}
	enable = false
	pbit {
		index = 5
		packet_queue_for_p_bit = 0
	}
	pbit {
		index = 8
	}
	queue {
		index = 6
		bandwidth_for_queue = 1
		scheduler_weight = 1
	}
	queue {
		index = 8
	}
}


resource "verity_packet_queue" "packet_queue_test_script2" {
	object_properties {
		group = "test"
	}
	enable = false
	pbit {
		index = 6
		packet_queue_for_p_bit = 2
	}
	pbit {
		index = 8
	}
	queue {
		index = 2
		bandwidth_for_queue = 30
		scheduler_weight = 1
	}
	queue {
		index = 8
	}
}

