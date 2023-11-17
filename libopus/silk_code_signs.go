package libopus

import (
	"github.com/gotranspile/opus/silk"
)

func silk_encode_signs(psRangeEnc *ec_enc, pulses []int8, length int, signalType int, quantOffsetType int, sum_pulses [20]int) {
	silk.EncodeSigns(psRangeEnc, pulses, length, signalType, quantOffsetType, sum_pulses)
}
