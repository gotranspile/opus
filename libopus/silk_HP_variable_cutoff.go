package libopus

func silk_HP_variable_cutoff(state_Fxx [0]silk_encoder_state_FLP) {
	var (
		quality_Q15       int64
		pitch_freq_Hz_Q16 opus_int32
		pitch_freq_log_Q7 opus_int32
		delta_freq_Q7     opus_int32
		psEncC1           *silk_encoder_state = &state_Fxx[0].SCmn
	)
	if int64(psEncC1.PrevSignalType) == TYPE_VOICED {
		pitch_freq_Hz_Q16 = (opus_int32(opus_uint32(psEncC1.Fs_kHz*1000) << 16)) / opus_int32(psEncC1.PrevLag)
		pitch_freq_log_Q7 = silk_lin2log(pitch_freq_Hz_Q16) - (16 << 7)
		quality_Q15 = psEncC1.Input_quality_bands_Q15[0]
		pitch_freq_log_Q7 = pitch_freq_log_Q7 + (((((opus_int32(opus_uint32(-quality_Q15) << 2)) * opus_int32(int64(opus_int16(quality_Q15)))) >> 16) * opus_int32(int64(opus_int16(pitch_freq_log_Q7-(silk_lin2log(opus_int32(VARIABLE_HP_MIN_CUTOFF_HZ*(1<<16)+0.5))-(16<<7)))))) >> 16)
		delta_freq_Q7 = pitch_freq_log_Q7 - (psEncC1.Variable_HP_smth1_Q15 >> 8)
		if delta_freq_Q7 < 0 {
			delta_freq_Q7 = delta_freq_Q7 * 3
		}
		if (-opus_int32(VARIABLE_HP_MAX_DELTA_FREQ*(1<<7) + 0.5)) > (opus_int32(VARIABLE_HP_MAX_DELTA_FREQ*(1<<7) + 0.5)) {
			if delta_freq_Q7 > (-opus_int32(VARIABLE_HP_MAX_DELTA_FREQ*(1<<7) + 0.5)) {
				delta_freq_Q7 = -opus_int32(VARIABLE_HP_MAX_DELTA_FREQ*(1<<7) + 0.5)
			} else if delta_freq_Q7 < (opus_int32(VARIABLE_HP_MAX_DELTA_FREQ*(1<<7) + 0.5)) {
				delta_freq_Q7 = opus_int32(VARIABLE_HP_MAX_DELTA_FREQ*(1<<7) + 0.5)
			} else {
				delta_freq_Q7 = delta_freq_Q7
			}
		} else if delta_freq_Q7 > (opus_int32(VARIABLE_HP_MAX_DELTA_FREQ*(1<<7) + 0.5)) {
			delta_freq_Q7 = opus_int32(VARIABLE_HP_MAX_DELTA_FREQ*(1<<7) + 0.5)
		} else if delta_freq_Q7 < (-opus_int32(VARIABLE_HP_MAX_DELTA_FREQ*(1<<7) + 0.5)) {
			delta_freq_Q7 = -opus_int32(VARIABLE_HP_MAX_DELTA_FREQ*(1<<7) + 0.5)
		} else {
			delta_freq_Q7 = delta_freq_Q7
		}
		psEncC1.Variable_HP_smth1_Q15 = psEncC1.Variable_HP_smth1_Q15 + (((opus_int32(opus_int16(psEncC1.Speech_activity_Q8)) * opus_int32(opus_int16(delta_freq_Q7))) * opus_int32(int64(opus_int16(opus_int32(VARIABLE_HP_SMTH_COEF1*(1<<16)+0.5))))) >> 16)
		if (opus_int32(opus_uint32(silk_lin2log(VARIABLE_HP_MIN_CUTOFF_HZ)) << 8)) > (opus_int32(opus_uint32(silk_lin2log(VARIABLE_HP_MAX_CUTOFF_HZ)) << 8)) {
			if psEncC1.Variable_HP_smth1_Q15 > (opus_int32(opus_uint32(silk_lin2log(VARIABLE_HP_MIN_CUTOFF_HZ)) << 8)) {
				psEncC1.Variable_HP_smth1_Q15 = opus_int32(opus_uint32(silk_lin2log(VARIABLE_HP_MIN_CUTOFF_HZ)) << 8)
			} else if psEncC1.Variable_HP_smth1_Q15 < (opus_int32(opus_uint32(silk_lin2log(VARIABLE_HP_MAX_CUTOFF_HZ)) << 8)) {
				psEncC1.Variable_HP_smth1_Q15 = opus_int32(opus_uint32(silk_lin2log(VARIABLE_HP_MAX_CUTOFF_HZ)) << 8)
			} else {
				psEncC1.Variable_HP_smth1_Q15 = psEncC1.Variable_HP_smth1_Q15
			}
		} else if psEncC1.Variable_HP_smth1_Q15 > (opus_int32(opus_uint32(silk_lin2log(VARIABLE_HP_MAX_CUTOFF_HZ)) << 8)) {
			psEncC1.Variable_HP_smth1_Q15 = opus_int32(opus_uint32(silk_lin2log(VARIABLE_HP_MAX_CUTOFF_HZ)) << 8)
		} else if psEncC1.Variable_HP_smth1_Q15 < (opus_int32(opus_uint32(silk_lin2log(VARIABLE_HP_MIN_CUTOFF_HZ)) << 8)) {
			psEncC1.Variable_HP_smth1_Q15 = opus_int32(opus_uint32(silk_lin2log(VARIABLE_HP_MIN_CUTOFF_HZ)) << 8)
		} else {
			psEncC1.Variable_HP_smth1_Q15 = psEncC1.Variable_HP_smth1_Q15
		}
	}
}
