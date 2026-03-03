# Test case: Add new resources
# Define resources to be injected into the corresponding .tf file for testing

resource "verity_packet_queue" "packet_queue_test_script1" {
    name = "packet_queue_test_script1"
    depends_on = [verity_operation_stage.packet_queue_stage]
	object_properties {
		group = ""
	}
	enable = true
	pbit {
		index = 1
		packet_queue_for_p_bit = 0
	}
	queue {
		index = 1
		bandwidth_for_queue = 40
		scheduler_type = ""
		scheduler_weight = 0
	}
}
