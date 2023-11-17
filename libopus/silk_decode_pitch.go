package libopus

import (
	"github.com/gotranspile/opus/silk"
)

func silk_decode_pitch(lagIndex int16, contourIndex int8, pitch_lags []int, Fs_kHz int, nb_subfr int) {
	silk.DecodePitch(lagIndex, contourIndex, pitch_lags, Fs_kHz, nb_subfr)
}
