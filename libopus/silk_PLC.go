package libopus

import (
	"github.com/gotranspile/opus/silk"
)

func silk_PLC_Reset(psDec *silk_decoder_state) {
	silk.PLC_Reset(psDec)
}
func silk_PLC(psDec *silk_decoder_state, psDecCtrl *silk_decoder_control, frame []int16, lost int, arch int) {
	silk.PLC(psDec, psDecCtrl, frame, lost, arch)
}
func silk_PLC_glue_frames(psDec *silk_decoder_state, frame []int16, length int) {
	silk.PLC_glue_frames(psDec, frame, length)
}
