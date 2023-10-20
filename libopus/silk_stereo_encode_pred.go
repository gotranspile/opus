package libopus

func silk_stereo_encode_pred(psRangeEnc *ec_enc, ix [2][3]int8) {
	var n int
	n = int(ix[0][2])*5 + int(ix[1][2])
	ec_enc_icdf(psRangeEnc, n, &silk_stereo_pred_joint_iCDF[0], 8)
	for n = 0; n < 2; n++ {
		ec_enc_icdf(psRangeEnc, int(ix[n][0]), &silk_uniform3_iCDF[0], 8)
		ec_enc_icdf(psRangeEnc, int(ix[n][1]), &silk_uniform5_iCDF[0], 8)
	}
}
func silk_stereo_encode_mid_only(psRangeEnc *ec_enc, mid_only_flag int8) {
	ec_enc_icdf(psRangeEnc, int(mid_only_flag), &silk_stereo_only_code_mid_iCDF[0], 8)
}
