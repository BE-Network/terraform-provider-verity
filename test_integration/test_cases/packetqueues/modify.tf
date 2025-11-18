# Test case: Modify existing resources
# Define modified versions of the resources from add.tf

resource "verity_packet_queue" "packet_queue_test_script1" {
	object_properties {
		group = "test"
	}
	enable = false
	pbit {
		index = 2
		packet_queue_for_p_bit = 0
	}
	queue {
		index = 2
		bandwidth_for_queue = 40
		scheduler_type = ""
		scheduler_weight = 0
	}
}