package libopus

func silk_stereo_quant_pred(pred_Q13 [0]opus_int32, ix [2][3]int8) {
	var (
		i              int64
		j              int64
		n              int64
		low_Q13        opus_int32
		step_Q13       opus_int32
		lvl_Q13        opus_int32
		err_min_Q13    opus_int32
		err_Q13        opus_int32
		quant_pred_Q13 opus_int32 = 0
	)
	for n = 0; n < 2; n++ {
		err_min_Q13 = silk_int32_MAX
		for i = 0; i < STEREO_QUANT_TAB_SIZE-1; i++ {
			low_Q13 = opus_int32(silk_stereo_pred_quant_Q13[i])
			step_Q13 = ((opus_int32(silk_stereo_pred_quant_Q13[i+1]) - low_Q13) * opus_int32(int64(opus_int16(opus_int32((0.5/STEREO_QUANT_SUB_STEPS)*(1<<16)+0.5))))) >> 16
			for j = 0; j < STEREO_QUANT_SUB_STEPS; j++ {
				lvl_Q13 = low_Q13 + (opus_int32(opus_int16(step_Q13)))*opus_int32(opus_int16(j*2+1))
				if (pred_Q13[n] - lvl_Q13) > 0 {
					err_Q13 = pred_Q13[n] - lvl_Q13
				} else {
					err_Q13 = -(pred_Q13[n] - lvl_Q13)
				}
				if err_Q13 < err_min_Q13 {
					err_min_Q13 = err_Q13
					quant_pred_Q13 = lvl_Q13
					ix[n][0] = int8(i)
					ix[n][1] = int8(j)
				} else {
					goto done
				}
			}
		}
	done:
		ix[n][2] = int8(opus_int32(int64(ix[n][0]) / 3))
		ix[n][0] -= int8(int64(ix[n][2]) * 3)
		pred_Q13[n] = quant_pred_Q13
	}
	pred_Q13[0] -= pred_Q13[1]
}
