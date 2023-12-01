package libopus

import "github.com/gotranspile/opus/silk"

func silk_A2NLSF_FLP(NLSF_Q15 []int16, pAR []float32, LPC_order int) {
	silk.A2NLSF_FLP(NLSF_Q15, pAR, LPC_order)
}
func silk_NLSF2A_FLP(pAR []float32, NLSF_Q15 []int16, LPC_order int, arch int) {
	silk.NLSF2A_FLP(pAR, NLSF_Q15, LPC_order, arch)
}
func silk_process_NLSFs_FLP(psEncC *silk_encoder_state, PredCoef [2][16]float32, NLSF_Q15 [16]int16, prev_NLSF_Q15 [16]int16) {
	silk.ProcessNLSFs_FLP(psEncC, PredCoef, NLSF_Q15, prev_NLSF_Q15)
}
func silk_NSQ_wrapper_FLP(psEnc *silk_encoder_state_FLP, psEncCtrl *silk_encoder_control_FLP, psIndices *SideInfoIndices, psNSQ *silk_nsq_state, pulses []int8, x []float32) {
	silk.NSQ_wrapper_FLP(psEnc, psEncCtrl, psIndices, psNSQ, pulses, x)
}
func silk_quant_LTP_gains_FLP(B [20]float32, cbk_index [4]int8, periodicity_index *int8, sum_log_gain_Q7 *int32, pred_gain_dB *float32, XX [100]float32, xX [20]float32, subfr_len int, nb_subfr int, arch int) {
	silk.QuantLTPGains_FLP(B, cbk_index, periodicity_index, sum_log_gain_Q7, pred_gain_dB, XX, xX, subfr_len, nb_subfr, arch)
}
