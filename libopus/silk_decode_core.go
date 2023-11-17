package libopus

import (
	"github.com/gotranspile/opus/silk"
)

func silk_decode_core(psDec *silk_decoder_state, psDecCtrl *silk_decoder_control, xq []int16, pulses [320]int16, arch int) {
	silk.DecodeCore(psDec, psDecCtrl, xq, pulses, arch)
}
