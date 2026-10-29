
resource "verity_sfp_breakout" "SFP_Breakouts" {
    name = "SFP Breakouts"
    depends_on = [verity_operation_stage.sfp_breakout_stage]
	breakout {
		index = 1
		breakout = "1x100G"
		enable = false
		part_number = ""
		vendor = "aaa"
	}
}

