
resource "verity_packet_queue" "test" {
    name = "test"
    depends_on = [verity_operation_stage.packet_queue_stage]
	enable = true
	pbit {
		index = 1
		packet_queue_for_p_bit = 0
	}
	pbit {
		index = 2
		packet_queue_for_p_bit = 1
	}
	pbit {
		index = 3
		packet_queue_for_p_bit = 2
	}
	pbit {
		index = 4
		packet_queue_for_p_bit = 3
	}
	pbit {
		index = 5
		packet_queue_for_p_bit = 0
	}
	pbit {
		index = 6
		packet_queue_for_p_bit = 1
	}
	pbit {
		index = 7
		packet_queue_for_p_bit = 2
	}
	pbit {
		index = 8
		packet_queue_for_p_bit = 3
	}
	queue {
		index = 1
		bandwidth_for_queue = 40
		scheduler_type = ""
		scheduler_weight = 0
	}
	queue {
		index = 2
		bandwidth_for_queue = 40
		scheduler_type = ""
		scheduler_weight = 0
	}
	queue {
		index = 3
		bandwidth_for_queue = 10
		scheduler_type = ""
		scheduler_weight = 0
	}
	queue {
		index = 4
		bandwidth_for_queue = 10
		scheduler_type = ""
		scheduler_weight = 0
	}
	queue {
		index = 5
		bandwidth_for_queue = 0
		scheduler_type = ""
		scheduler_weight = 0
	}
	queue {
		index = 6
		bandwidth_for_queue = 0
		scheduler_type = ""
		scheduler_weight = 0
	}
	queue {
		index = 7
		bandwidth_for_queue = 0
		scheduler_type = ""
		scheduler_weight = 0
	}
	queue {
		index = 8
		bandwidth_for_queue = 0
		scheduler_type = ""
		scheduler_weight = 0
	}
}

