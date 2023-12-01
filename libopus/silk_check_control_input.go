package libopus

import "github.com/gotranspile/opus/silk"

func check_control_input(encControl *silk_EncControlStruct) int {
	return silk.CheckControlInput(encControl)
}
