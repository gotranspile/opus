package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func silk_VAD_Init(psSilk_VAD *silk_VAD_state) int {
	var (
		b   int
		ret int = 0
	)
	*psSilk_VAD = silk_VAD_state{}
	for b = 0; b < VAD_N_BANDS; b++ {
		psSilk_VAD.NoiseLevelBias[b] = silk_max_32(int32(VAD_NOISE_LEVELS_BIAS/(b+1)), 1)
	}
	for b = 0; b < VAD_N_BANDS; b++ {
		psSilk_VAD.NL[b] = int32(int(psSilk_VAD.NoiseLevelBias[b]) * 100)
		psSilk_VAD.Inv_NL[b] = int32(silk_int32_MAX / int(psSilk_VAD.NL[b]))
	}
	psSilk_VAD.Counter = 15
	for b = 0; b < VAD_N_BANDS; b++ {
		psSilk_VAD.NrgRatioSmth_Q8[b] = 100 * 256
	}
	return ret
}

var tiltWeights [4]int32 = [4]int32{30000, 6000, -12000, -12000}

func silk_VAD_GetSA_Q8_c(psEncC *silk_encoder_state, pIn []int16) int {
	var (
		SA_Q15                 int
		pSNR_dB_Q7             int
		input_tilt             int
		decimated_framelength1 int
		decimated_framelength2 int
		decimated_framelength  int
		dec_subframe_length    int
		dec_subframe_offset    int
		SNR_Q7                 int
		i                      int
		b                      int
		s                      int
		sumSquared             int32
		smooth_coef_Q16        int32
		HPstateTmp             int16
		X                      *int16
		Xnrg                   [4]int32
		NrgToNoiseRatio_Q8     [4]int32
		speech_nrg             int32
		x_tmp                  int32
		X_offset               [4]int
		ret                    int             = 0
		psSilk_VAD             *silk_VAD_state = &psEncC.SVAD
	)
	decimated_framelength1 = psEncC.Frame_length >> 1
	decimated_framelength2 = psEncC.Frame_length >> 2
	decimated_framelength = psEncC.Frame_length >> 3
	X_offset[0] = 0
	X_offset[1] = decimated_framelength + decimated_framelength2
	X_offset[2] = X_offset[1] + decimated_framelength
	X_offset[3] = X_offset[2] + decimated_framelength2
	X = (*int16)(libc.Malloc((X_offset[3] + decimated_framelength1) * int(unsafe.Sizeof(int16(0)))))
	silk_ana_filt_bank_1(&pIn[0], &psSilk_VAD.AnaState[0], X, (*int16)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(int16(0))*uintptr(X_offset[3]))), int32(psEncC.Frame_length))
	silk_ana_filt_bank_1(X, &psSilk_VAD.AnaState1[0], X, (*int16)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(int16(0))*uintptr(X_offset[2]))), int32(decimated_framelength1))
	silk_ana_filt_bank_1(X, &psSilk_VAD.AnaState2[0], X, (*int16)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(int16(0))*uintptr(X_offset[1]))), int32(decimated_framelength2))
	*(*int16)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(int16(0))*uintptr(decimated_framelength-1))) = int16(int(*(*int16)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(int16(0))*uintptr(decimated_framelength-1)))) >> 1)
	HPstateTmp = *(*int16)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(int16(0))*uintptr(decimated_framelength-1)))
	for i = decimated_framelength - 1; i > 0; i-- {
		*(*int16)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(int16(0))*uintptr(i-1))) = int16(int(*(*int16)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(int16(0))*uintptr(i-1)))) >> 1)
		*(*int16)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(int16(0))*uintptr(i))) -= *(*int16)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(int16(0))*uintptr(i-1)))
	}
	*(*int16)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(int16(0))*0)) -= psSilk_VAD.HPstate
	psSilk_VAD.HPstate = HPstateTmp
	for b = 0; b < VAD_N_BANDS; b++ {
		decimated_framelength = psEncC.Frame_length >> silk_min_int(VAD_N_BANDS-b, int(VAD_N_BANDS-1))
		dec_subframe_length = decimated_framelength >> VAD_INTERNAL_SUBFRAMES_LOG2
		dec_subframe_offset = 0
		Xnrg[b] = psSilk_VAD.XnrgSubfr[b]
		for s = 0; s < (int(1 << VAD_INTERNAL_SUBFRAMES_LOG2)); s++ {
			sumSquared = 0
			for i = 0; i < dec_subframe_length; i++ {
				x_tmp = int32(int(*(*int16)(unsafe.Add(unsafe.Pointer(X), unsafe.Sizeof(int16(0))*uintptr(X_offset[b]+i+dec_subframe_offset)))) >> 3)
				sumSquared = int32(int(sumSquared) + int(int32(int16(x_tmp)))*int(int32(int16(x_tmp))))
			}
			if s < (int(1<<VAD_INTERNAL_SUBFRAMES_LOG2))-1 {
				if ((int(uint32(Xnrg[b])) + int(uint32(sumSquared))) & 0x80000000) != 0 {
					Xnrg[b] = silk_int32_MAX
				} else {
					Xnrg[b] = int32(int(Xnrg[b]) + int(sumSquared))
				}
			} else {
				if ((int(uint32(Xnrg[b])) + int(uint32(int32(int(sumSquared)>>1)))) & 0x80000000) != 0 {
					Xnrg[b] = silk_int32_MAX
				} else {
					Xnrg[b] = int32(int(Xnrg[b]) + (int(sumSquared) >> 1))
				}
			}
			dec_subframe_offset += dec_subframe_length
		}
		psSilk_VAD.XnrgSubfr[b] = sumSquared
	}
	silk_VAD_GetNoiseLevels([4]int32(&Xnrg[0]), psSilk_VAD)
	sumSquared = 0
	input_tilt = 0
	for b = 0; b < VAD_N_BANDS; b++ {
		speech_nrg = int32(int(Xnrg[b]) - int(psSilk_VAD.NL[b]))
		if int(speech_nrg) > 0 {
			if (int(Xnrg[b]) & 0xFF800000) == 0 {
				NrgToNoiseRatio_Q8[b] = int32(int(int32(int(uint32(Xnrg[b]))<<8)) / (int(psSilk_VAD.NL[b]) + 1))
			} else {
				NrgToNoiseRatio_Q8[b] = int32(int(Xnrg[b]) / ((int(psSilk_VAD.NL[b]) >> 8) + 1))
			}
			SNR_Q7 = int(silk_lin2log(NrgToNoiseRatio_Q8[b])) - 8*128
			sumSquared = int32(int(sumSquared) + int(int32(int16(SNR_Q7)))*int(int32(int16(SNR_Q7))))
			if int(speech_nrg) < (1 << 20) {
				SNR_Q7 = int(int32((int64(int32(int(uint32(silk_SQRT_APPROX(speech_nrg)))<<6)) * int64(int16(SNR_Q7))) >> 16))
			}
			input_tilt = int(int32(input_tilt + int((int64(tiltWeights[b])*int64(int16(SNR_Q7)))>>16)))
		} else {
			NrgToNoiseRatio_Q8[b] = 256
		}
	}
	sumSquared = int32(int(sumSquared) / VAD_N_BANDS)
	pSNR_dB_Q7 = int(int16(int(silk_SQRT_APPROX(sumSquared)) * 3))
	SA_Q15 = silk_sigm_Q15(int(int32((VAD_SNR_FACTOR_Q16*int64(int16(pSNR_dB_Q7)))>>16)) - VAD_NEGATIVE_OFFSET_Q5)
	psEncC.Input_tilt_Q15 = int(int32(int(uint32(int32(silk_sigm_Q15(input_tilt)-16384))) << 1))
	speech_nrg = 0
	for b = 0; b < VAD_N_BANDS; b++ {
		speech_nrg += int32((b + 1) * ((int(Xnrg[b]) - int(psSilk_VAD.NL[b])) >> 4))
	}
	if psEncC.Frame_length == psEncC.Fs_kHz*20 {
		speech_nrg = int32(int(speech_nrg) >> 1)
	}
	if int(speech_nrg) <= 0 {
		SA_Q15 = SA_Q15 >> 1
	} else if int(speech_nrg) < 16384 {
		speech_nrg = int32(int(uint32(speech_nrg)) << 16)
		speech_nrg = silk_SQRT_APPROX(speech_nrg)
		SA_Q15 = int(int32(((int(speech_nrg) + 32768) * int(int64(int16(SA_Q15)))) >> 16))
	}
	psEncC.Speech_activity_Q8 = silk_min_int(SA_Q15>>7, silk_uint8_MAX)
	smooth_coef_Q16 = int32((VAD_SNR_SMOOTH_COEF_Q18 * int64(int16(int32((int64(int32(SA_Q15))*int64(int16(SA_Q15)))>>16)))) >> 16)
	if psEncC.Frame_length == psEncC.Fs_kHz*10 {
		smooth_coef_Q16 >>= 1
	}
	for b = 0; b < VAD_N_BANDS; b++ {
		psSilk_VAD.NrgRatioSmth_Q8[b] = int32(int(psSilk_VAD.NrgRatioSmth_Q8[b]) + (((int(NrgToNoiseRatio_Q8[b]) - int(psSilk_VAD.NrgRatioSmth_Q8[b])) * int(int64(int16(smooth_coef_Q16)))) >> 16))
		SNR_Q7 = (int(silk_lin2log(psSilk_VAD.NrgRatioSmth_Q8[b])) - 8*128) * 3
		psEncC.Input_quality_bands_Q15[b] = silk_sigm_Q15((SNR_Q7 - 16*128) >> 4)
	}
	return ret
}
func silk_VAD_GetNoiseLevels(pX [4]int32, psSilk_VAD *silk_VAD_state) {
	var (
		k        int
		nl       int32
		nrg      int32
		inv_nrg  int32
		coef     int
		min_coef int
	)
	if int(psSilk_VAD.Counter) < 1000 {
		min_coef = int(int32(silk_int16_MAX / ((int(psSilk_VAD.Counter) >> 4) + 1)))
		psSilk_VAD.Counter++
	} else {
		min_coef = 0
	}
	for k = 0; k < VAD_N_BANDS; k++ {
		nl = psSilk_VAD.NL[k]
		if ((int(uint32(pX[k])) + int(uint32(psSilk_VAD.NoiseLevelBias[k]))) & 0x80000000) != 0 {
			nrg = silk_int32_MAX
		} else {
			nrg = int32(int(pX[k]) + int(psSilk_VAD.NoiseLevelBias[k]))
		}
		inv_nrg = int32(silk_int32_MAX / int(nrg))
		if int(nrg) > int(int32(int(uint32(nl))<<3)) {
			coef = int(VAD_NOISE_LEVEL_SMOOTH_COEF_Q16 >> 3)
		} else if int(nrg) < int(nl) {
			coef = VAD_NOISE_LEVEL_SMOOTH_COEF_Q16
		} else {
			coef = int(int32((int64(int32((int64(inv_nrg)*int64(nl))>>16)) * int64(int16(int(VAD_NOISE_LEVEL_SMOOTH_COEF_Q16<<1)))) >> 16))
		}
		coef = silk_max_int(coef, min_coef)
		psSilk_VAD.Inv_NL[k] = int32(int(psSilk_VAD.Inv_NL[k]) + (((int(inv_nrg) - int(psSilk_VAD.Inv_NL[k])) * int(int64(int16(coef)))) >> 16))
		nl = int32(silk_int32_MAX / int(psSilk_VAD.Inv_NL[k]))
		if int(nl) < 0xFFFFFF {
			nl = nl
		} else {
			nl = 0xFFFFFF
		}
		psSilk_VAD.NL[k] = nl
	}
}
