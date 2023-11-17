package libopus

import (
	"github.com/gotranspile/opus/silk"
)

func silk_CNG(psDec *silk_decoder_state, psDecCtrl *silk_decoder_control, frame []int16, length int) {
	silk.CNG(psDec, psDecCtrl, frame, length)
}
