package libopus

import (
	"github.com/gotranspile/opus/silk"
)

func silk_decode_parameters(psDec *silk_decoder_state, psDecCtrl *silk_decoder_control, condCoding int) {
	silk.DecodeParameters(psDec, psDecCtrl, condCoding)
}
