package libopus

func silk_init_encoder(psEnc *silk_encoder_state_FLP, arch int) int {
	var ret int = 0
	*psEnc = silk_encoder_state_FLP{}
	psEnc.SCmn.Arch = arch
	psEnc.SCmn.Variable_HP_smth1_Q15 = int32(int(uint32(int32(int(silk_lin2log(int32(float64(int(VARIABLE_HP_MIN_CUTOFF_HZ*(1<<16)))+0.5)))-(16<<7)))) << 8)
	psEnc.SCmn.Variable_HP_smth2_Q15 = psEnc.SCmn.Variable_HP_smth1_Q15
	psEnc.SCmn.First_frame_after_reset = 1
	ret += silk_VAD_Init(&psEnc.SCmn.SVAD)
	return ret
}
