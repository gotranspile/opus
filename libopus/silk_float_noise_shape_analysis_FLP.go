package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

func warped_gain(coefs *float32, lambda float32, order int) float32 {
	var (
		i    int
		gain float32
	)
	lambda = -lambda
	gain = *(*float32)(unsafe.Add(unsafe.Pointer(coefs), unsafe.Sizeof(float32(0))*uintptr(order-1)))
	for i = order - 2; i >= 0; i-- {
		gain = lambda*gain + *(*float32)(unsafe.Add(unsafe.Pointer(coefs), unsafe.Sizeof(float32(0))*uintptr(i)))
	}
	return 1.0 / (1.0 - lambda*gain)
}
func warped_true2monic_coefs(coefs *float32, lambda float32, limit float32, order int) {
	var (
		i      int
		iter   int
		ind    int = 0
		tmp    float32
		maxabs float32
		chirp  float32
		gain   float32
	)
	for i = order - 1; i > 0; i-- {
		*(*float32)(unsafe.Add(unsafe.Pointer(coefs), unsafe.Sizeof(float32(0))*uintptr(i-1))) -= lambda * *(*float32)(unsafe.Add(unsafe.Pointer(coefs), unsafe.Sizeof(float32(0))*uintptr(i)))
	}
	gain = (1.0 - lambda*lambda) / (lambda**(*float32)(unsafe.Add(unsafe.Pointer(coefs), unsafe.Sizeof(float32(0))*0)) + 1.0)
	for i = 0; i < order; i++ {
		*(*float32)(unsafe.Add(unsafe.Pointer(coefs), unsafe.Sizeof(float32(0))*uintptr(i))) *= gain
	}
	for iter = 0; iter < 10; iter++ {
		maxabs = -1.0
		for i = 0; i < order; i++ {
			tmp = float32(math.Abs(float64(*(*float32)(unsafe.Add(unsafe.Pointer(coefs), unsafe.Sizeof(float32(0))*uintptr(i))))))
			if tmp > maxabs {
				maxabs = tmp
				ind = i
			}
		}
		if maxabs <= limit {
			return
		}
		for i = 1; i < order; i++ {
			*(*float32)(unsafe.Add(unsafe.Pointer(coefs), unsafe.Sizeof(float32(0))*uintptr(i-1))) += lambda * *(*float32)(unsafe.Add(unsafe.Pointer(coefs), unsafe.Sizeof(float32(0))*uintptr(i)))
		}
		gain = 1.0 / gain
		for i = 0; i < order; i++ {
			*(*float32)(unsafe.Add(unsafe.Pointer(coefs), unsafe.Sizeof(float32(0))*uintptr(i))) *= gain
		}
		chirp = float32(0.99 - (float64(iter)*0.1+0.8)*float64(maxabs-limit)/float64(maxabs*float32(ind+1)))
		silk_bwexpander_FLP([]float32(coefs), order, chirp)
		for i = order - 1; i > 0; i-- {
			*(*float32)(unsafe.Add(unsafe.Pointer(coefs), unsafe.Sizeof(float32(0))*uintptr(i-1))) -= lambda * *(*float32)(unsafe.Add(unsafe.Pointer(coefs), unsafe.Sizeof(float32(0))*uintptr(i)))
		}
		gain = (1.0 - lambda*lambda) / (lambda**(*float32)(unsafe.Add(unsafe.Pointer(coefs), unsafe.Sizeof(float32(0))*0)) + 1.0)
		for i = 0; i < order; i++ {
			*(*float32)(unsafe.Add(unsafe.Pointer(coefs), unsafe.Sizeof(float32(0))*uintptr(i))) *= gain
		}
	}
}
func limit_coefs(coefs *float32, limit float32, order int) {
	var (
		i      int
		iter   int
		ind    int = 0
		tmp    float32
		maxabs float32
		chirp  float32
	)
	for iter = 0; iter < 10; iter++ {
		maxabs = -1.0
		for i = 0; i < order; i++ {
			tmp = float32(math.Abs(float64(*(*float32)(unsafe.Add(unsafe.Pointer(coefs), unsafe.Sizeof(float32(0))*uintptr(i))))))
			if tmp > maxabs {
				maxabs = tmp
				ind = i
			}
		}
		if maxabs <= limit {
			return
		}
		chirp = float32(0.99 - (float64(iter)*0.1+0.8)*float64(maxabs-limit)/float64(maxabs*float32(ind+1)))
		silk_bwexpander_FLP([]float32(coefs), order, chirp)
	}
}
func silk_noise_shape_analysis_FLP(psEnc *silk_encoder_state_FLP, psEncCtrl *silk_encoder_control_FLP, pitch_res *float32, x *float32) {
	var (
		psShapeSt        *silk_shape_state_FLP = &psEnc.SShape
		k                int
		nSamples         int
		nSegs            int
		SNR_adj_dB       float32
		HarmShapeGain    float32
		Tilt             float32
		nrg              float32
		log_energy       float32
		log_energy_prev  float32
		energy_variation float32
		BWExp            float32
		gain_mult        float32
		gain_add         float32
		strength         float32
		b                float32
		warping          float32
		x_windowed       [240]float32
		auto_corr        [25]float32
		rc               [25]float32
		x_ptr            *float32
		pitch_res_ptr    *float32
	)
	x_ptr = (*float32)(unsafe.Add(unsafe.Pointer(x), -int(unsafe.Sizeof(float32(0))*uintptr(psEnc.SCmn.La_shape))))
	SNR_adj_dB = float32(float64(psEnc.SCmn.SNR_dB_Q7) * (1 / 128.0))
	psEncCtrl.Input_quality = float32(float64(psEnc.SCmn.Input_quality_bands_Q15[0]+psEnc.SCmn.Input_quality_bands_Q15[1]) * 0.5 * (1.0 / 32768.0))
	psEncCtrl.Coding_quality = silk_sigmoid((SNR_adj_dB - 20.0) * 0.25)
	if psEnc.SCmn.UseCBR == 0 {
		b = float32(1.0 - float64(psEnc.SCmn.Speech_activity_Q8)*(1.0/256.0))
		SNR_adj_dB -= BG_SNR_DECR_dB * psEncCtrl.Coding_quality * (psEncCtrl.Input_quality*0.5 + 0.5) * b * b
	}
	if int(psEnc.SCmn.Indices.SignalType) == TYPE_VOICED {
		SNR_adj_dB += HARM_SNR_INCR_dB * psEnc.LTPCorr
	} else {
		SNR_adj_dB += float32((float64(psEnc.SCmn.SNR_dB_Q7)*(-0.4)*(1/128.0) + 6.0) * float64(1.0-psEncCtrl.Input_quality))
	}
	if int(psEnc.SCmn.Indices.SignalType) == TYPE_VOICED {
		psEnc.SCmn.Indices.QuantOffsetType = 0
	} else {
		nSamples = psEnc.SCmn.Fs_kHz * 2
		energy_variation = 0.0
		log_energy_prev = 0.0
		pitch_res_ptr = pitch_res
		nSegs = (SUB_FRAME_LENGTH_MS * int(int32(int16(psEnc.SCmn.Nb_subfr)))) / 2
		for k = 0; k < nSegs; k++ {
			nrg = float32(nSamples) + float32(silk_energy_FLP([]float32(pitch_res_ptr), nSamples))
			log_energy = silk_log2(float64(nrg))
			if k > 0 {
				energy_variation += float32(math.Abs(float64(log_energy - log_energy_prev)))
			}
			log_energy_prev = log_energy
			pitch_res_ptr = (*float32)(unsafe.Add(unsafe.Pointer(pitch_res_ptr), unsafe.Sizeof(float32(0))*uintptr(nSamples)))
		}
		if float64(energy_variation) > ENERGY_VARIATION_THRESHOLD_QNT_OFFSET*float64(nSegs-1) {
			psEnc.SCmn.Indices.QuantOffsetType = 0
		} else {
			psEnc.SCmn.Indices.QuantOffsetType = 1
		}
	}
	strength = FIND_PITCH_WHITE_NOISE_FRACTION * psEncCtrl.PredGain
	BWExp = BANDWIDTH_EXPANSION / (strength*strength + 1.0)
	warping = float32(psEnc.SCmn.Warping_Q16)/65536.0 + psEncCtrl.Coding_quality*0.01
	for k = 0; k < psEnc.SCmn.Nb_subfr; k++ {
		var (
			shift      int
			slope_part int
			flat_part  int
		)
		flat_part = psEnc.SCmn.Fs_kHz * 3
		slope_part = (psEnc.SCmn.ShapeWinLength - flat_part) / 2
		silk_apply_sine_window_FLP(x_windowed[:], []float32(x_ptr), 1, slope_part)
		shift = slope_part
		libc.MemCpy(unsafe.Pointer(&x_windowed[shift]), unsafe.Pointer((*float32)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(float32(0))*uintptr(shift)))), flat_part*int(unsafe.Sizeof(float32(0))))
		shift += flat_part
		silk_apply_sine_window_FLP([]float32(&x_windowed[shift]), []float32((*float32)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(float32(0))*uintptr(shift)))), 2, slope_part)
		x_ptr = (*float32)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(float32(0))*uintptr(psEnc.SCmn.Subfr_length)))
		if psEnc.SCmn.Warping_Q16 > 0 {
			silk_warped_autocorrelation_FLP(&auto_corr[0], &x_windowed[0], warping, psEnc.SCmn.ShapeWinLength, psEnc.SCmn.ShapingLPCOrder)
		} else {
			silk_autocorrelation_FLP(&auto_corr[0], &x_windowed[0], psEnc.SCmn.ShapeWinLength, psEnc.SCmn.ShapingLPCOrder+1)
		}
		auto_corr[0] += auto_corr[0]*SHAPE_WHITE_NOISE_FRACTION + 1.0
		nrg = silk_schur_FLP(rc[:], auto_corr[:], psEnc.SCmn.ShapingLPCOrder)
		silk_k2a_FLP([]float32(&psEncCtrl.AR[k*MAX_SHAPE_LPC_ORDER]), rc[:], int32(psEnc.SCmn.ShapingLPCOrder))
		psEncCtrl.Gains[k] = float32(math.Sqrt(float64(nrg)))
		if psEnc.SCmn.Warping_Q16 > 0 {
			psEncCtrl.Gains[k] *= warped_gain(&psEncCtrl.AR[k*MAX_SHAPE_LPC_ORDER], warping, psEnc.SCmn.ShapingLPCOrder)
		}
		silk_bwexpander_FLP([]float32(&psEncCtrl.AR[k*MAX_SHAPE_LPC_ORDER]), psEnc.SCmn.ShapingLPCOrder, BWExp)
		if psEnc.SCmn.Warping_Q16 > 0 {
			warped_true2monic_coefs(&psEncCtrl.AR[k*MAX_SHAPE_LPC_ORDER], warping, 3.999, psEnc.SCmn.ShapingLPCOrder)
		} else {
			limit_coefs(&psEncCtrl.AR[k*MAX_SHAPE_LPC_ORDER], 3.999, psEnc.SCmn.ShapingLPCOrder)
		}
	}
	gain_mult = float32(math.Pow(2.0, float64(SNR_adj_dB*(-0.16))))
	gain_add = float32(math.Pow(2.0, MIN_QGAIN_DB*0.16))
	for k = 0; k < psEnc.SCmn.Nb_subfr; k++ {
		psEncCtrl.Gains[k] *= gain_mult
		psEncCtrl.Gains[k] += gain_add
	}
	strength = float32(LOW_FREQ_SHAPING * (LOW_QUALITY_LOW_FREQ_SHAPING_DECR*(float64(psEnc.SCmn.Input_quality_bands_Q15[0])*(1.0/32768.0)-1.0) + 1.0))
	strength *= float32(float64(psEnc.SCmn.Speech_activity_Q8) * (1.0 / 256.0))
	if int(psEnc.SCmn.Indices.SignalType) == TYPE_VOICED {
		for k = 0; k < psEnc.SCmn.Nb_subfr; k++ {
			b = float32(0.2/float64(psEnc.SCmn.Fs_kHz) + 3.0/float64(psEncCtrl.PitchL[k]))
			psEncCtrl.LF_MA_shp[k] = b + (-1.0)
			psEncCtrl.LF_AR_shp[k] = 1.0 - b - b*strength
		}
		Tilt = float32(-HP_NOISE_COEF - (1-HP_NOISE_COEF)*HARM_HP_NOISE_COEF*float64(psEnc.SCmn.Speech_activity_Q8)*(1.0/256.0))
	} else {
		b = float32(1.3 / float64(psEnc.SCmn.Fs_kHz))
		psEncCtrl.LF_MA_shp[0] = b + (-1.0)
		psEncCtrl.LF_AR_shp[0] = 1.0 - b - b*strength*0.6
		for k = 1; k < psEnc.SCmn.Nb_subfr; k++ {
			psEncCtrl.LF_MA_shp[k] = psEncCtrl.LF_MA_shp[0]
			psEncCtrl.LF_AR_shp[k] = psEncCtrl.LF_AR_shp[0]
		}
		Tilt = -HP_NOISE_COEF
	}
	if USE_HARM_SHAPING != 0 && int(psEnc.SCmn.Indices.SignalType) == TYPE_VOICED {
		HarmShapeGain = HARMONIC_SHAPING
		HarmShapeGain += HIGH_RATE_OR_LOW_QUALITY_HARMONIC_SHAPING * (1.0 - (1.0-psEncCtrl.Coding_quality)*psEncCtrl.Input_quality)
		HarmShapeGain *= float32(math.Sqrt(float64(psEnc.LTPCorr)))
	} else {
		HarmShapeGain = 0.0
	}
	for k = 0; k < psEnc.SCmn.Nb_subfr; k++ {
		psShapeSt.HarmShapeGain_smth += SUBFR_SMTH_COEF * (HarmShapeGain - psShapeSt.HarmShapeGain_smth)
		psEncCtrl.HarmShapeGain[k] = psShapeSt.HarmShapeGain_smth
		psShapeSt.Tilt_smth += SUBFR_SMTH_COEF * (Tilt - psShapeSt.Tilt_smth)
		psEncCtrl.Tilt[k] = psShapeSt.Tilt_smth
	}
}
