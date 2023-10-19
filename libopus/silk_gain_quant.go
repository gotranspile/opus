package libopus

const OFFSET = 2090
const SCALE_Q16 = 2251
const INV_SCALE_Q16 = 0x1D1C71

func silk_gains_quant(ind [4]int8, gain_Q16 [4]opus_int32, prev_ind *int8, conditional int64, nb_subfr int64) {
	var (
		k                          int64
		double_step_size_threshold int64
	)
	for k = 0; k < nb_subfr; k++ {
		ind[k] = int8(opus_int32(((((N_LEVELS_QGAIN - 1) * 65536) / (((MAX_QGAIN_DB - MIN_QGAIN_DB) * 128) / 6)) * int64(opus_int16(silk_lin2log(gain_Q16[k])-opus_int32((MIN_QGAIN_DB*128)/6+16*128)))) >> 16))
		if int64(ind[k]) < int64(*prev_ind) {
			ind[k]++
		}
		if 0 > (N_LEVELS_QGAIN - 1) {
			if int64(ind[k]) > 0 {
				ind[k] = 0
			} else if int64(ind[k]) < (N_LEVELS_QGAIN - 1) {
				ind[k] = int8(N_LEVELS_QGAIN - 1)
			} else {
				ind[k] = ind[k]
			}
		} else if int64(ind[k]) > (N_LEVELS_QGAIN - 1) {
			ind[k] = int8(N_LEVELS_QGAIN - 1)
		} else if int64(ind[k]) < 0 {
			ind[k] = 0
		} else {
			ind[k] = ind[k]
		}
		if k == 0 && conditional == 0 {
			if (int64(*prev_ind) + (-4)) > (N_LEVELS_QGAIN - 1) {
				if int64(ind[k]) > (int64(*prev_ind) + (-4)) {
					ind[k] = int8(int64(*prev_ind) + (-4))
				} else if int64(ind[k]) < (N_LEVELS_QGAIN - 1) {
					ind[k] = int8(N_LEVELS_QGAIN - 1)
				} else {
					ind[k] = ind[k]
				}
			} else if int64(ind[k]) > (N_LEVELS_QGAIN - 1) {
				ind[k] = int8(N_LEVELS_QGAIN - 1)
			} else if int64(ind[k]) < (int64(*prev_ind) + (-4)) {
				ind[k] = int8(int64(*prev_ind) + (-4))
			} else {
				ind[k] = ind[k]
			}
			*prev_ind = ind[k]
		} else {
			ind[k] = int8(int64(ind[k]) - int64(*prev_ind))
			double_step_size_threshold = MAX_DELTA_GAIN_QUANT*2 - N_LEVELS_QGAIN + int64(*prev_ind)
			if int64(ind[k]) > double_step_size_threshold {
				ind[k] = int8(double_step_size_threshold + ((int64(ind[k]) - double_step_size_threshold + 1) >> 1))
			}
			if int64(-4) > MAX_DELTA_GAIN_QUANT {
				if int64(ind[k]) > int64(-4) {
					ind[k] = -4
				} else if int64(ind[k]) < MAX_DELTA_GAIN_QUANT {
					ind[k] = MAX_DELTA_GAIN_QUANT
				} else {
					ind[k] = ind[k]
				}
			} else if int64(ind[k]) > MAX_DELTA_GAIN_QUANT {
				ind[k] = MAX_DELTA_GAIN_QUANT
			} else if int64(ind[k]) < int64(-4) {
				ind[k] = -4
			} else {
				ind[k] = ind[k]
			}
			if int64(ind[k]) > double_step_size_threshold {
				*prev_ind += int8((opus_int32(opus_uint32(ind[k]) << 1)) - opus_int32(double_step_size_threshold))
				*prev_ind = int8(silk_min_int(int64(*prev_ind), N_LEVELS_QGAIN-1))
			} else {
				*prev_ind += ind[k]
			}
			ind[k] -= -4
		}
		gain_Q16[k] = silk_log2lin(silk_min_32((opus_int32(((((((MAX_QGAIN_DB-MIN_QGAIN_DB)*128)/6)*65536)/(N_LEVELS_QGAIN-1))*int64(opus_int16(*prev_ind)))>>16))+opus_int32((MIN_QGAIN_DB*128)/6+16*128), 3967))
	}
}
func silk_gains_dequant(gain_Q16 [4]opus_int32, ind [4]int8, prev_ind *int8, conditional int64, nb_subfr int64) {
	var (
		k                          int64
		ind_tmp                    int64
		double_step_size_threshold int64
	)
	for k = 0; k < nb_subfr; k++ {
		if k == 0 && conditional == 0 {
			*prev_ind = int8(silk_max_int(int64(ind[k]), int64(*prev_ind)-16))
		} else {
			ind_tmp = int64(ind[k]) + (-4)
			double_step_size_threshold = MAX_DELTA_GAIN_QUANT*2 - N_LEVELS_QGAIN + int64(*prev_ind)
			if ind_tmp > double_step_size_threshold {
				*prev_ind += int8((opus_int32(opus_uint32(ind_tmp) << 1)) - opus_int32(double_step_size_threshold))
			} else {
				*prev_ind += int8(ind_tmp)
			}
		}
		if 0 > (N_LEVELS_QGAIN - 1) {
			if int64(*prev_ind) > 0 {
				*prev_ind = 0
			} else if int64(*prev_ind) < (N_LEVELS_QGAIN - 1) {
				*prev_ind = int8(N_LEVELS_QGAIN - 1)
			} else {
				*prev_ind = *prev_ind
			}
		} else if int64(*prev_ind) > (N_LEVELS_QGAIN - 1) {
			*prev_ind = int8(N_LEVELS_QGAIN - 1)
		} else if int64(*prev_ind) < 0 {
			*prev_ind = 0
		} else {
			*prev_ind = *prev_ind
		}
		gain_Q16[k] = silk_log2lin(silk_min_32((opus_int32(((((((MAX_QGAIN_DB-MIN_QGAIN_DB)*128)/6)*65536)/(N_LEVELS_QGAIN-1))*int64(opus_int16(*prev_ind)))>>16))+opus_int32((MIN_QGAIN_DB*128)/6+16*128), 3967))
	}
}
func silk_gains_ID(ind [4]int8, nb_subfr int64) opus_int32 {
	var (
		k       int64
		gainsID opus_int32
	)
	gainsID = 0
	for k = 0; k < nb_subfr; k++ {
		gainsID = opus_int32(ind[k]) + (opus_int32(opus_uint32(gainsID) << 8))
	}
	return gainsID
}
