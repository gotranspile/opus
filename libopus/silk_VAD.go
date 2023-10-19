package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func silk_VAD_Init(psSilk_VAD *silk_VAD_state) int64 {
	var (
		b   int64
		ret int64 = 0
	)
	*psSilk_VAD = silk_VAD_state{}
	for b = 0; b < VAD_N_BANDS; b++ {
		psSilk_VAD.NoiseLevelBias[b] = silk_max_32(opus_int32(VAD_NOISE_LEVELS_BIAS/(b+1)), 1)
	}
	for b = 0; b < VAD_N_BANDS; b++ {
		psSilk_VAD.NL[b] = (psSilk_VAD.NoiseLevelBias[b]) * 100
		psSilk_VAD.Inv_NL[b] = silk_int32_MAX / (psSilk_VAD.NL[b])
	}
	psSilk_VAD.Counter = 15
	for b = 0; b < VAD_N_BANDS; b++ {
		psSilk_VAD.NrgRatioSmth_Q8[b] = 100 * 256
	}
	return ret
}

var tiltWeights [4]opus_int32 = [4]opus_int32{30000, 6000, -12000, -12000}

func silk_VAD_GetSA_Q8_c(psEncC *silk_encoder_state, pIn [0]opus_int16) int64 {
	var (
		SA_Q15                 int64
		pSNR_dB_Q7             int64
		input_tilt             int64
		decimated_framelength1 int64
		decimated_framelength2 int64
		decimated_framelength  int64
		dec_subframe_length    int64
		dec_subframe_offset    int64
		SNR_Q7                 int64
		i                      int64
		b                      int64
		s                      int64
		sumSquared             opus_int32
		smooth_coef_Q16        opus_int32
		HPstateTmp             opus_int16
		X                      *opus_int16
		Xnrg                   [4]opus_int32
		NrgToNoiseRatio_Q8     [4]opus_int32
		speech_nrg             opus_int32
		x_tmp                  opus_int32
		X_offset               [4]int64
		ret                    int64           = 0
		psSilk_VAD             *silk_VAD_state = &psEncC.SVAD
	)
	decimated_framelength1 = psEncC.Frame_length >> 1
	decimated_framelength2 = psEncC.Frame_length >> 2
	decimated_framelength = psEncC.Frame_length >> 3
	X_offset[0] = 0
	X_offset[1] = decimated_framelength + decimated_framelength2
	X_offset[2] = X_offset[1] + decimated_framelength
	X_offset[3] = X_offset[2] + decimated_framelength2
	X = (*opus_int16)(libc.Malloc(int((X_offset[3] + decimated_framelength1) * int64(unsafe.Sizeof(opus_int16(0))))))
	silk_ana_filt_bank_1(&pIn[0], &psSilk_VAD.AnaState[0], X, (*opus_int16)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(opus_int16(0))*uintptr(X_offset[3]))), opus_int32(psEncC.Frame_length))
	silk_ana_filt_bank_1(X, &psSilk_VAD.AnaState1[0], X, (*opus_int16)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(opus_int16(0))*uintptr(X_offset[2]))), opus_int32(decimated_framelength1))
	silk_ana_filt_bank_1(X, &psSilk_VAD.AnaState2[0], X, (*opus_int16)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(opus_int16(0))*uintptr(X_offset[1]))), opus_int32(decimated_framelength2))
	*(*opus_int16)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(opus_int16(0))*uintptr(decimated_framelength-1))) = (*(*opus_int16)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(opus_int16(0))*uintptr(decimated_framelength-1)))) >> 1
	HPstateTmp = *(*opus_int16)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(opus_int16(0))*uintptr(decimated_framelength-1)))
	for i = decimated_framelength - 1; i > 0; i-- {
		*(*opus_int16)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(opus_int16(0))*uintptr(i-1))) = (*(*opus_int16)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(opus_int16(0))*uintptr(i-1)))) >> 1
		*(*opus_int16)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(opus_int16(0))*uintptr(i))) -= *(*opus_int16)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(opus_int16(0))*uintptr(i-1)))
	}
	*(*opus_int16)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(opus_int16(0))*0)) -= psSilk_VAD.HPstate
	psSilk_VAD.HPstate = HPstateTmp
	for b = 0; b < VAD_N_BANDS; b++ {
		decimated_framelength = psEncC.Frame_length >> silk_min_int(VAD_N_BANDS-b, VAD_N_BANDS-1)
		dec_subframe_length = decimated_framelength >> VAD_INTERNAL_SUBFRAMES_LOG2
		dec_subframe_offset = 0
		Xnrg[b] = psSilk_VAD.XnrgSubfr[b]
		for s = 0; s < (1 << VAD_INTERNAL_SUBFRAMES_LOG2); s++ {
			sumSquared = 0
			for i = 0; i < dec_subframe_length; i++ {
				x_tmp = opus_int32((*(*opus_int16)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(opus_int16(0))*uintptr(X_offset[b]+i+dec_subframe_offset)))) >> 3)
				sumSquared = sumSquared + (opus_int32(opus_int16(x_tmp)))*opus_int32(opus_int16(x_tmp))
			}
			if s < (1<<VAD_INTERNAL_SUBFRAMES_LOG2)-1 {
				if ((opus_uint32(Xnrg[b]) + opus_uint32(sumSquared)) & 0x80000000) != 0 {
					Xnrg[b] = silk_int32_MAX
				} else {
					Xnrg[b] = (Xnrg[b]) + sumSquared
				}
			} else {
				if ((opus_uint32(Xnrg[b]) + opus_uint32(sumSquared>>1)) & 0x80000000) != 0 {
					Xnrg[b] = silk_int32_MAX
				} else {
					Xnrg[b] = (Xnrg[b]) + (sumSquared >> 1)
				}
			}
			dec_subframe_offset += dec_subframe_length
		}
		psSilk_VAD.XnrgSubfr[b] = sumSquared
	}
	silk_VAD_GetNoiseLevels([4]opus_int32(&Xnrg[0]), psSilk_VAD)
	sumSquared = 0
	input_tilt = 0
	for b = 0; b < VAD_N_BANDS; b++ {
		speech_nrg = Xnrg[b] - psSilk_VAD.NL[b]
		if speech_nrg > 0 {
			if (Xnrg[b] & 0xFF800000) == 0 {
				NrgToNoiseRatio_Q8[b] = (opus_int32(opus_uint32(Xnrg[b]) << 8)) / (psSilk_VAD.NL[b] + 1)
			} else {
				NrgToNoiseRatio_Q8[b] = (Xnrg[b]) / (((psSilk_VAD.NL[b]) >> 8) + 1)
			}
			SNR_Q7 = int64(silk_lin2log(NrgToNoiseRatio_Q8[b]) - 8*128)
			sumSquared = sumSquared + (opus_int32(opus_int16(SNR_Q7)))*opus_int32(opus_int16(SNR_Q7))
			if speech_nrg < (1 << 20) {
				SNR_Q7 = int64(((opus_int32(opus_uint32(silk_SQRT_APPROX(speech_nrg)) << 6)) * opus_int32(int64(opus_int16(SNR_Q7)))) >> 16)
			}
			input_tilt = int64(opus_int32(input_tilt + int64(((tiltWeights[b])*opus_int32(int64(opus_int16(SNR_Q7))))>>16)))
		} else {
			NrgToNoiseRatio_Q8[b] = 256
		}
	}
	sumSquared = sumSquared / VAD_N_BANDS
	pSNR_dB_Q7 = int64(opus_int16(silk_SQRT_APPROX(sumSquared) * 3))
	SA_Q15 = silk_sigm_Q15(int64((opus_int32((VAD_SNR_FACTOR_Q16 * int64(opus_int16(pSNR_dB_Q7))) >> 16)) - VAD_NEGATIVE_OFFSET_Q5))
	psEncC.Input_tilt_Q15 = int64(opus_int32(opus_uint32(silk_sigm_Q15(input_tilt)-16384) << 1))
	speech_nrg = 0
	for b = 0; b < VAD_N_BANDS; b++ {
		speech_nrg += opus_int32((b + 1) * int64((Xnrg[b]-psSilk_VAD.NL[b])>>4))
	}
	if psEncC.Frame_length == psEncC.Fs_kHz*20 {
		speech_nrg = speech_nrg >> 1
	}
	if speech_nrg <= 0 {
		SA_Q15 = SA_Q15 >> 1
	} else if speech_nrg < 16384 {
		speech_nrg = opus_int32(opus_uint32(speech_nrg) << 16)
		speech_nrg = silk_SQRT_APPROX(speech_nrg)
		SA_Q15 = int64(((speech_nrg + 32768) * opus_int32(int64(opus_int16(SA_Q15)))) >> 16)
	}
	psEncC.Speech_activity_Q8 = silk_min_int(SA_Q15>>7, silk_uint8_MAX)
	smooth_coef_Q16 = opus_int32((VAD_SNR_SMOOTH_COEF_Q18 * int64(opus_int16(((opus_int32(SA_Q15))*opus_int32(int64(opus_int16(SA_Q15))))>>16))) >> 16)
	if psEncC.Frame_length == psEncC.Fs_kHz*10 {
		smooth_coef_Q16 >>= 1
	}
	for b = 0; b < VAD_N_BANDS; b++ {
		psSilk_VAD.NrgRatioSmth_Q8[b] = (psSilk_VAD.NrgRatioSmth_Q8[b]) + (((NrgToNoiseRatio_Q8[b] - psSilk_VAD.NrgRatioSmth_Q8[b]) * opus_int32(int64(opus_int16(smooth_coef_Q16)))) >> 16)
		SNR_Q7 = int64((silk_lin2log(psSilk_VAD.NrgRatioSmth_Q8[b]) - 8*128) * 3)
		psEncC.Input_quality_bands_Q15[b] = silk_sigm_Q15((SNR_Q7 - 16*128) >> 4)
	}
	return ret
}
func silk_VAD_GetNoiseLevels(pX [4]opus_int32, psSilk_VAD *silk_VAD_state) {
	var (
		k        int64
		nl       opus_int32
		nrg      opus_int32
		inv_nrg  opus_int32
		coef     int64
		min_coef int64
	)
	if psSilk_VAD.Counter < 1000 {
		min_coef = int64(silk_int16_MAX / ((psSilk_VAD.Counter >> 4) + 1))
		psSilk_VAD.Counter++
	} else {
		min_coef = 0
	}
	for k = 0; k < VAD_N_BANDS; k++ {
		nl = psSilk_VAD.NL[k]
		if ((opus_uint32(pX[k]) + opus_uint32(psSilk_VAD.NoiseLevelBias[k])) & 0x80000000) != 0 {
			nrg = silk_int32_MAX
		} else {
			nrg = (pX[k]) + (psSilk_VAD.NoiseLevelBias[k])
		}
		inv_nrg = silk_int32_MAX / nrg
		if nrg > (opus_int32(opus_uint32(nl) << 3)) {
			coef = VAD_NOISE_LEVEL_SMOOTH_COEF_Q16 >> 3
		} else if nrg < nl {
			coef = VAD_NOISE_LEVEL_SMOOTH_COEF_Q16
		} else {
			coef = int64(((opus_int32((int64(inv_nrg) * int64(nl)) >> 16)) * opus_int32(int64(opus_int16(VAD_NOISE_LEVEL_SMOOTH_COEF_Q16<<1)))) >> 16)
		}
		coef = silk_max_int(coef, min_coef)
		psSilk_VAD.Inv_NL[k] = (psSilk_VAD.Inv_NL[k]) + (((inv_nrg - psSilk_VAD.Inv_NL[k]) * opus_int32(int64(opus_int16(coef)))) >> 16)
		nl = silk_int32_MAX / (psSilk_VAD.Inv_NL[k])
		if nl < 0xFFFFFF {
			nl = nl
		} else {
			nl = 0xFFFFFF
		}
		psSilk_VAD.NL[k] = nl
	}
}
