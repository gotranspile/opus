package libopus

func silk_find_LPC_FLP(psEncC *silk_encoder_state, NLSF_Q15 []int16, x []float32, minInvGain float32) {
	var (
		k              int
		subfr_length   int
		a              [16]float32
		res_nrg        float32
		res_nrg_2nd    float32
		res_nrg_interp float32
		NLSF0_Q15      [16]int16
		a_tmp          [16]float32
		LPC_res        [384]float32
	)
	subfr_length = psEncC.Subfr_length + psEncC.PredictLPCOrder
	psEncC.Indices.NLSFInterpCoef_Q2 = 4
	res_nrg = silk_burg_modified_FLP(a[:], x, minInvGain, subfr_length, psEncC.Nb_subfr, psEncC.PredictLPCOrder)
	if psEncC.UseInterpolatedNLSFs != 0 && psEncC.First_frame_after_reset == 0 && psEncC.Nb_subfr == MAX_NB_SUBFR {
		res_nrg -= silk_burg_modified_FLP(a_tmp[:], []float32(&x[(int(MAX_NB_SUBFR/2))*subfr_length]), minInvGain, subfr_length, int(MAX_NB_SUBFR/2), psEncC.PredictLPCOrder)
		silk_A2NLSF_FLP(&NLSF_Q15[0], &a_tmp[0], psEncC.PredictLPCOrder)
		res_nrg_2nd = FLT_MAX
		for k = 3; k >= 0; k-- {
			silk_interpolate(NLSF0_Q15, psEncC.Prev_NLSFq_Q15, [16]int16(NLSF_Q15), k, psEncC.PredictLPCOrder)
			silk_NLSF2A_FLP(&a_tmp[0], &NLSF0_Q15[0], psEncC.PredictLPCOrder, psEncC.Arch)
			silk_LPC_analysis_filter_FLP(LPC_res[:], a_tmp[:], x, subfr_length*2, psEncC.PredictLPCOrder)
			res_nrg_interp = float32(silk_energy_FLP(&LPC_res[psEncC.PredictLPCOrder], subfr_length-psEncC.PredictLPCOrder) + silk_energy_FLP(&LPC_res[psEncC.PredictLPCOrder+subfr_length], subfr_length-psEncC.PredictLPCOrder))
			if res_nrg_interp < res_nrg {
				res_nrg = res_nrg_interp
				psEncC.Indices.NLSFInterpCoef_Q2 = int8(k)
			} else if res_nrg_interp > res_nrg_2nd {
				break
			}
			res_nrg_2nd = res_nrg_interp
		}
	}
	if int(psEncC.Indices.NLSFInterpCoef_Q2) == 4 {
		silk_A2NLSF_FLP(&NLSF_Q15[0], &a[0], psEncC.PredictLPCOrder)
	}
}
