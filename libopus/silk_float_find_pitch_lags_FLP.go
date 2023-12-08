package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func silk_find_pitch_lags_FLP(psEnc *silk_encoder_state_FLP, psEncCtrl *silk_encoder_control_FLP, res []float32, x []float32, arch int) {
	var (
		buf_len   int
		thrhld    float32
		res_nrg   float32
		x_buf_ptr *float32
		x_buf     *float32
		auto_corr [17]float32
		A         [16]float32
		refl_coef [16]float32
		Wsig      [384]float32
		Wsig_ptr  *float32
	)
	buf_len = psEnc.SCmn.La_pitch + psEnc.SCmn.Frame_length + psEnc.SCmn.Ltp_mem_length
	x_buf = (*float32)(unsafe.Add(unsafe.Pointer(&x[0]), -int(unsafe.Sizeof(float32(0))*uintptr(psEnc.SCmn.Ltp_mem_length))))
	x_buf_ptr = (*float32)(unsafe.Add(unsafe.Pointer((*float32)(unsafe.Add(unsafe.Pointer(x_buf), unsafe.Sizeof(float32(0))*uintptr(buf_len)))), -int(unsafe.Sizeof(float32(0))*uintptr(psEnc.SCmn.Pitch_LPC_win_length))))
	Wsig_ptr = &Wsig[0]
	silk_apply_sine_window_FLP([]float32(Wsig_ptr), []float32(x_buf_ptr), 1, psEnc.SCmn.La_pitch)
	Wsig_ptr = (*float32)(unsafe.Add(unsafe.Pointer(Wsig_ptr), unsafe.Sizeof(float32(0))*uintptr(psEnc.SCmn.La_pitch)))
	x_buf_ptr = (*float32)(unsafe.Add(unsafe.Pointer(x_buf_ptr), unsafe.Sizeof(float32(0))*uintptr(psEnc.SCmn.La_pitch)))
	libc.MemCpy(unsafe.Pointer(Wsig_ptr), unsafe.Pointer(x_buf_ptr), (psEnc.SCmn.Pitch_LPC_win_length-(psEnc.SCmn.La_pitch<<1))*int(unsafe.Sizeof(float32(0))))
	Wsig_ptr = (*float32)(unsafe.Add(unsafe.Pointer(Wsig_ptr), unsafe.Sizeof(float32(0))*uintptr(psEnc.SCmn.Pitch_LPC_win_length-(psEnc.SCmn.La_pitch<<1))))
	x_buf_ptr = (*float32)(unsafe.Add(unsafe.Pointer(x_buf_ptr), unsafe.Sizeof(float32(0))*uintptr(psEnc.SCmn.Pitch_LPC_win_length-(psEnc.SCmn.La_pitch<<1))))
	silk_apply_sine_window_FLP([]float32(Wsig_ptr), []float32(x_buf_ptr), 2, psEnc.SCmn.La_pitch)
	silk_autocorrelation_FLP(&auto_corr[0], &Wsig[0], psEnc.SCmn.Pitch_LPC_win_length, psEnc.SCmn.PitchEstimationLPCOrder+1)
	auto_corr[0] += auto_corr[0]*FIND_PITCH_WHITE_NOISE_FRACTION + 1
	res_nrg = silk_schur_FLP(refl_coef[:], auto_corr[:], psEnc.SCmn.PitchEstimationLPCOrder)
	psEncCtrl.PredGain = auto_corr[0] / (func() float32 {
		if res_nrg > 1.0 {
			return res_nrg
		}
		return 1.0
	}())
	silk_k2a_FLP(A[:], refl_coef[:], int32(psEnc.SCmn.PitchEstimationLPCOrder))
	silk_bwexpander_FLP(A[:], psEnc.SCmn.PitchEstimationLPCOrder, FIND_PITCH_BANDWIDTH_EXPANSION)
	silk_LPC_analysis_filter_FLP(res, A[:], []float32(x_buf), buf_len, psEnc.SCmn.PitchEstimationLPCOrder)
	if int(psEnc.SCmn.Indices.SignalType) != TYPE_NO_VOICE_ACTIVITY && psEnc.SCmn.First_frame_after_reset == 0 {
		thrhld = 0.6
		thrhld -= float32(float64(psEnc.SCmn.PitchEstimationLPCOrder) * 0.004)
		thrhld -= float32(float64(psEnc.SCmn.Speech_activity_Q8) * 0.1 * (1.0 / 256.0))
		thrhld -= float32(float64(int(psEnc.SCmn.PrevSignalType)>>1) * 0.15)
		thrhld -= float32(float64(psEnc.SCmn.Input_tilt_Q15) * 0.1 * (1.0 / 32768.0))
		if silk_pitch_analysis_core_FLP(res, psEncCtrl.PitchL[:], &psEnc.SCmn.Indices.LagIndex, &psEnc.SCmn.Indices.ContourIndex, &psEnc.LTPCorr, psEnc.SCmn.PrevLag, float32(float64(psEnc.SCmn.PitchEstimationThreshold_Q16)/65536.0), thrhld, psEnc.SCmn.Fs_kHz, psEnc.SCmn.PitchEstimationComplexity, psEnc.SCmn.Nb_subfr, arch) == 0 {
			psEnc.SCmn.Indices.SignalType = TYPE_VOICED
		} else {
			psEnc.SCmn.Indices.SignalType = TYPE_UNVOICED
		}
	} else {
		*(*[4]int)(unsafe.Pointer(&psEncCtrl.PitchL[0])) = [4]int{}
		psEnc.SCmn.Indices.LagIndex = 0
		psEnc.SCmn.Indices.ContourIndex = 0
		psEnc.LTPCorr = 0
	}
}
