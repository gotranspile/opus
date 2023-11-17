package libopus

import "github.com/gotranspile/opus/silk"

func silk_gains_quant(ind [4]int8, gain_Q16 [4]int32, prev_ind *int8, conditional int, nb_subfr int) {
	silk.GainsQuant(ind, gain_Q16, prev_ind, conditional, nb_subfr)
}
func silk_gains_dequant(gain_Q16 [4]int32, ind [4]int8, prev_ind *int8, conditional int, nb_subfr int) {
	silk.GainsDequant(gain_Q16, ind, prev_ind, conditional, nb_subfr)
}
