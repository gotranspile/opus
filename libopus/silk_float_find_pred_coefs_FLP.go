package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

func silk_find_pred_coefs_FLP(psEnc *silk_encoder_state_FLP, psEncCtrl *silk_encoder_control_FLP, res_pitch []float32, x []float32, condCoding int) {
	var (
		i          int
		XXLTP      [100]float32
		xXLTP      [20]float32
		invGains   [4]float32
		NLSF_Q15   [16]int16 = [16]int16{}
		x_ptr      *float32
		x_pre_ptr  *float32
		LPC_in_pre [384]float32
		minInvGain float32
	)
	for i = 0; i < psEnc.SCmn.Nb_subfr; i++ {
		invGains[i] = 1.0 / psEncCtrl.Gains[i]
	}
	if int(psEnc.SCmn.Indices.SignalType) == TYPE_VOICED {
		silk_find_LTP_FLP(XXLTP, xXLTP, res_pitch, psEncCtrl.PitchL, psEnc.SCmn.Subfr_length, psEnc.SCmn.Nb_subfr)
		silk_quant_LTP_gains_FLP(psEncCtrl.LTPCoef, psEnc.SCmn.Indices.LTPIndex, &psEnc.SCmn.Indices.PERIndex, &psEnc.SCmn.Sum_log_gain_Q7, &psEncCtrl.LTPredCodGain, XXLTP, xXLTP, psEnc.SCmn.Subfr_length, psEnc.SCmn.Nb_subfr, psEnc.SCmn.Arch)
		silk_LTP_scale_ctrl_FLP(psEnc, psEncCtrl, condCoding)
		silk_LTP_analysis_filter_FLP(&LPC_in_pre[0], (*float32)(unsafe.Add(unsafe.Pointer(&x[0]), -int(unsafe.Sizeof(float32(0))*uintptr(psEnc.SCmn.PredictLPCOrder)))), psEncCtrl.LTPCoef, psEncCtrl.PitchL, invGains, psEnc.SCmn.Subfr_length, psEnc.SCmn.Nb_subfr, psEnc.SCmn.PredictLPCOrder)
	} else {
		x_ptr = (*float32)(unsafe.Add(unsafe.Pointer(&x[0]), -int(unsafe.Sizeof(float32(0))*uintptr(psEnc.SCmn.PredictLPCOrder))))
		x_pre_ptr = &LPC_in_pre[0]
		for i = 0; i < psEnc.SCmn.Nb_subfr; i++ {
			silk_scale_copy_vector_FLP(x_pre_ptr, x_ptr, invGains[i], psEnc.SCmn.Subfr_length+psEnc.SCmn.PredictLPCOrder)
			x_pre_ptr = (*float32)(unsafe.Add(unsafe.Pointer(x_pre_ptr), unsafe.Sizeof(float32(0))*uintptr(psEnc.SCmn.Subfr_length+psEnc.SCmn.PredictLPCOrder)))
			x_ptr = (*float32)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(float32(0))*uintptr(psEnc.SCmn.Subfr_length)))
		}
		libc.MemSet(unsafe.Pointer(&psEncCtrl.LTPCoef[0]), 0, psEnc.SCmn.Nb_subfr*LTP_ORDER*int(unsafe.Sizeof(float32(0))))
		psEncCtrl.LTPredCodGain = 0.0
		psEnc.SCmn.Sum_log_gain_Q7 = 0
	}
	if psEnc.SCmn.First_frame_after_reset != 0 {
		minInvGain = 1.0 / MAX_PREDICTION_POWER_GAIN_AFTER_RESET
	} else {
		minInvGain = float32(math.Pow(2, float64(psEncCtrl.LTPredCodGain/3))) / MAX_PREDICTION_POWER_GAIN
		minInvGain /= psEncCtrl.Coding_quality*0.75 + 0.25
	}
	silk_find_LPC_FLP(&psEnc.SCmn, NLSF_Q15[:], LPC_in_pre[:], minInvGain)
	silk_process_NLSFs_FLP(&psEnc.SCmn, psEncCtrl.PredCoef, NLSF_Q15, psEnc.SCmn.Prev_NLSFq_Q15)
	silk_residual_energy_FLP(psEncCtrl.ResNrg, LPC_in_pre[:], psEncCtrl.PredCoef, psEncCtrl.Gains[:], psEnc.SCmn.Subfr_length, psEnc.SCmn.Nb_subfr, psEnc.SCmn.PredictLPCOrder)
	libc.MemCpy(unsafe.Pointer(&psEnc.SCmn.Prev_NLSFq_Q15[0]), unsafe.Pointer(&NLSF_Q15[0]), int(unsafe.Sizeof([16]int16{})))
}
