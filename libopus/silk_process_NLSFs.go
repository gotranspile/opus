package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func silk_process_NLSFs(psEncC *silk_encoder_state, PredCoef_Q12 [2][16]opus_int16, pNLSF_Q15 [16]opus_int16, prev_NLSFq_Q15 [16]opus_int16) {
	var (
		i               int64
		doInterpolate   int64
		NLSF_mu_Q20     int64
		i_sqr_Q15       opus_int16
		pNLSF0_temp_Q15 [16]opus_int16
		pNLSFW_QW       [16]opus_int16
		pNLSFW0_temp_QW [16]opus_int16
	)
	NLSF_mu_Q20 = int64(opus_int32(((int64(opus_int16(psEncC.Speech_activity_Q8)) * int64(opus_int32((-0.001)*(1<<28)+0.5))) >> 16) + int64(opus_int32(0.003*(1<<20)+0.5))))
	if psEncC.Nb_subfr == 2 {
		NLSF_mu_Q20 = NLSF_mu_Q20 + (NLSF_mu_Q20 >> 1)
	}
	silk_NLSF_VQ_weights_laroia(&pNLSFW_QW[0], &pNLSF_Q15[0], psEncC.PredictLPCOrder)
	doInterpolate = int64(libc.BoolToInt(psEncC.UseInterpolatedNLSFs == 1 && int64(psEncC.Indices.NLSFInterpCoef_Q2) < 4))
	if doInterpolate != 0 {
		silk_interpolate(pNLSF0_temp_Q15, prev_NLSFq_Q15, pNLSF_Q15, int64(psEncC.Indices.NLSFInterpCoef_Q2), psEncC.PredictLPCOrder)
		silk_NLSF_VQ_weights_laroia(&pNLSFW0_temp_QW[0], &pNLSF0_temp_Q15[0], psEncC.PredictLPCOrder)
		i_sqr_Q15 = opus_int16(opus_int32(opus_uint32(opus_int32(opus_int16(psEncC.Indices.NLSFInterpCoef_Q2))*opus_int32(opus_int16(psEncC.Indices.NLSFInterpCoef_Q2))) << 11))
		for i = 0; i < psEncC.PredictLPCOrder; i++ {
			pNLSFW_QW[i] = opus_int16(opus_int32((pNLSFW_QW[i])>>1) + ((opus_int32(pNLSFW0_temp_QW[i]) * opus_int32(i_sqr_Q15)) >> 16))
		}
	}
	silk_NLSF_encode(&psEncC.Indices.NLSFIndices[0], &pNLSF_Q15[0], psEncC.PsNLSF_CB, &pNLSFW_QW[0], NLSF_mu_Q20, psEncC.NLSF_MSVQ_Survivors, int64(psEncC.Indices.SignalType))
	silk_NLSF2A(&PredCoef_Q12[1][0], &pNLSF_Q15[0], psEncC.PredictLPCOrder, psEncC.Arch)
	if doInterpolate != 0 {
		silk_interpolate(pNLSF0_temp_Q15, prev_NLSFq_Q15, pNLSF_Q15, int64(psEncC.Indices.NLSFInterpCoef_Q2), psEncC.PredictLPCOrder)
		silk_NLSF2A(&PredCoef_Q12[0][0], &pNLSF0_temp_Q15[0], psEncC.PredictLPCOrder, psEncC.Arch)
	} else {
		libc.MemCpy(unsafe.Pointer(&(PredCoef_Q12[0])[0]), unsafe.Pointer(&(PredCoef_Q12[1])[0]), int(psEncC.PredictLPCOrder*int64(unsafe.Sizeof(opus_int16(0)))))
	}
}
