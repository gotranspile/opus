package libopus

func silk_stereo_decode_pred(psRangeDec *ec_dec, pred_Q13 [0]opus_int32) {
	var (
		n        int64
		ix       [2][3]int64
		low_Q13  opus_int32
		step_Q13 opus_int32
	)
	n = ec_dec_icdf(psRangeDec, &silk_stereo_pred_joint_iCDF[0], 8)
	ix[0][2] = int64(opus_int32(n / 5))
	ix[1][2] = n - ix[0][2]*5
	for n = 0; n < 2; n++ {
		ix[n][0] = ec_dec_icdf(psRangeDec, &silk_uniform3_iCDF[0], 8)
		ix[n][1] = ec_dec_icdf(psRangeDec, &silk_uniform5_iCDF[0], 8)
	}
	for n = 0; n < 2; n++ {
		ix[n][0] += ix[n][2] * 3
		low_Q13 = opus_int32(silk_stereo_pred_quant_Q13[ix[n][0]])
		step_Q13 = ((opus_int32(silk_stereo_pred_quant_Q13[ix[n][0]+1]) - low_Q13) * opus_int32(int64(opus_int16(opus_int32((0.5/STEREO_QUANT_SUB_STEPS)*(1<<16)+0.5))))) >> 16
		pred_Q13[n] = low_Q13 + (opus_int32(opus_int16(step_Q13)))*opus_int32(opus_int16(ix[n][1]*2+1))
	}
	pred_Q13[0] -= pred_Q13[1]
}
func silk_stereo_decode_mid_only(psRangeDec *ec_dec, decode_only_mid *int64) {
	*decode_only_mid = ec_dec_icdf(psRangeDec, &silk_stereo_only_code_mid_iCDF[0], 8)
}
