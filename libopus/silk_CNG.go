package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

func silk_CNG_exc(exc_Q14 [0]opus_int32, exc_buf_Q14 [0]opus_int32, length int64, rand_seed *opus_int32) {
	var (
		seed     opus_int32
		i        int64
		idx      int64
		exc_mask int64
	)
	exc_mask = CNG_BUF_MASK_MAX
	for exc_mask > length {
		exc_mask = exc_mask >> 1
	}
	seed = *rand_seed
	for i = 0; i < length; i++ {
		seed = opus_int32(RAND_INCREMENT + opus_uint32(seed)*RAND_MULTIPLIER)
		idx = int64((seed >> 24) & opus_int32(exc_mask))
		exc_Q14[i] = exc_buf_Q14[idx]
	}
	*rand_seed = seed
}
func silk_CNG_Reset(psDec *silk_decoder_state) {
	var (
		i             int64
		NLSF_step_Q15 int64
		NLSF_acc_Q15  int64
	)
	NLSF_step_Q15 = int64(opus_int32(silk_int16_MAX / (psDec.LPC_order + 1)))
	NLSF_acc_Q15 = 0
	for i = 0; i < psDec.LPC_order; i++ {
		NLSF_acc_Q15 += NLSF_step_Q15
		psDec.SCNG.CNG_smth_NLSF_Q15[i] = opus_int16(NLSF_acc_Q15)
	}
	psDec.SCNG.CNG_smth_Gain_Q16 = 0
	psDec.SCNG.Rand_seed = 3176576
}
func silk_CNG(psDec *silk_decoder_state, psDecCtrl *silk_decoder_control, frame [0]opus_int16, length int64) {
	var (
		i            int64
		subfr        int64
		LPC_pred_Q10 opus_int32
		max_Gain_Q16 opus_int32
		gain_Q16     opus_int32
		gain_Q10     opus_int32
		A_Q12        [16]opus_int16
		psCNG        *silk_CNG_struct = &psDec.SCNG
	)
	if psDec.Fs_kHz != psCNG.Fs_kHz {
		silk_CNG_Reset(psDec)
		psCNG.Fs_kHz = psDec.Fs_kHz
	}
	if psDec.LossCnt == 0 && psDec.PrevSignalType == TYPE_NO_VOICE_ACTIVITY {
		for i = 0; i < psDec.LPC_order; i++ {
			psCNG.CNG_smth_NLSF_Q15[i] += opus_int16(((opus_int32(psDec.PrevNLSF_Q15[i]) - opus_int32(psCNG.CNG_smth_NLSF_Q15[i])) * CNG_NLSF_SMTH_Q16) >> 16)
		}
		max_Gain_Q16 = 0
		subfr = 0
		for i = 0; i < psDec.Nb_subfr; i++ {
			if psDecCtrl.Gains_Q16[i] > max_Gain_Q16 {
				max_Gain_Q16 = psDecCtrl.Gains_Q16[i]
				subfr = i
			}
		}
		libc.MemMove(unsafe.Pointer(&psCNG.CNG_exc_buf_Q14[psDec.Subfr_length]), unsafe.Pointer(&psCNG.CNG_exc_buf_Q14[0]), int((psDec.Nb_subfr-1)*psDec.Subfr_length*int64(unsafe.Sizeof(opus_int32(0)))))
		libc.MemCpy(unsafe.Pointer(&psCNG.CNG_exc_buf_Q14[0]), unsafe.Pointer(&psDec.Exc_Q14[subfr*psDec.Subfr_length]), int(psDec.Subfr_length*int64(unsafe.Sizeof(opus_int32(0)))))
		for i = 0; i < psDec.Nb_subfr; i++ {
			psCNG.CNG_smth_Gain_Q16 += ((psDecCtrl.Gains_Q16[i] - psCNG.CNG_smth_Gain_Q16) * CNG_GAIN_SMTH_Q16) >> 16
			if (opus_int32((int64(psCNG.CNG_smth_Gain_Q16) * CNG_GAIN_SMTH_THRESHOLD_Q16) >> 16)) > psDecCtrl.Gains_Q16[i] {
				psCNG.CNG_smth_Gain_Q16 = psDecCtrl.Gains_Q16[i]
			}
		}
	}
	if psDec.LossCnt != 0 {
		var CNG_sig_Q14 *opus_int32
		CNG_sig_Q14 = (*opus_int32)(libc.Malloc(int((length + MAX_LPC_ORDER) * int64(unsafe.Sizeof(opus_int32(0))))))
		gain_Q16 = opus_int32((int64(psDec.SPLC.RandScale_Q14) * int64(psDec.SPLC.PrevGain_Q16[1])) >> 16)
		if gain_Q16 >= (1<<21) || psCNG.CNG_smth_Gain_Q16 > (1<<23) {
			gain_Q16 = (gain_Q16 >> 16) * (gain_Q16 >> 16)
			gain_Q16 = ((psCNG.CNG_smth_Gain_Q16 >> 16) * (psCNG.CNG_smth_Gain_Q16 >> 16)) - (opus_int32(opus_uint32(gain_Q16) << 5))
			gain_Q16 = opus_int32(opus_uint32(silk_SQRT_APPROX(gain_Q16)) << 16)
		} else {
			gain_Q16 = opus_int32((int64(gain_Q16) * int64(gain_Q16)) >> 16)
			gain_Q16 = (opus_int32((int64(psCNG.CNG_smth_Gain_Q16) * int64(psCNG.CNG_smth_Gain_Q16)) >> 16)) - (opus_int32(opus_uint32(gain_Q16) << 5))
			gain_Q16 = opus_int32(opus_uint32(silk_SQRT_APPROX(gain_Q16)) << 8)
		}
		gain_Q10 = gain_Q16 >> 6
		silk_CNG_exc([0]opus_int32((*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER)))), psCNG.CNG_exc_buf_Q14[:], length, &psCNG.Rand_seed)
		silk_NLSF2A(&A_Q12[0], &psCNG.CNG_smth_NLSF_Q15[0], psDec.LPC_order, psDec.Arch)
		libc.MemCpy(unsafe.Pointer(CNG_sig_Q14), unsafe.Pointer(&psCNG.CNG_synth_state[0]), int(MAX_LPC_ORDER*unsafe.Sizeof(opus_int32(0))))
		for i = 0; i < length; i++ {
			LPC_pred_Q10 = opus_int32(psDec.LPC_order >> 1)
			LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-1)))) * opus_int32(int64(A_Q12[0]))) >> 16)
			LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-2)))) * opus_int32(int64(A_Q12[1]))) >> 16)
			LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-3)))) * opus_int32(int64(A_Q12[2]))) >> 16)
			LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-4)))) * opus_int32(int64(A_Q12[3]))) >> 16)
			LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-5)))) * opus_int32(int64(A_Q12[4]))) >> 16)
			LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-6)))) * opus_int32(int64(A_Q12[5]))) >> 16)
			LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-7)))) * opus_int32(int64(A_Q12[6]))) >> 16)
			LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-8)))) * opus_int32(int64(A_Q12[7]))) >> 16)
			LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-9)))) * opus_int32(int64(A_Q12[8]))) >> 16)
			LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-10)))) * opus_int32(int64(A_Q12[9]))) >> 16)
			if psDec.LPC_order == 16 {
				LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-11)))) * opus_int32(int64(A_Q12[10]))) >> 16)
				LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-12)))) * opus_int32(int64(A_Q12[11]))) >> 16)
				LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-13)))) * opus_int32(int64(A_Q12[12]))) >> 16)
				LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-14)))) * opus_int32(int64(A_Q12[13]))) >> 16)
				LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-15)))) * opus_int32(int64(A_Q12[14]))) >> 16)
				LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-16)))) * opus_int32(int64(A_Q12[15]))) >> 16)
			}
			if ((opus_uint32(*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) + opus_uint32(opus_int32(opus_uint32(func() opus_int32 {
				if (0x80000000 >> 4) > (silk_int32_MAX >> 4) {
					if LPC_pred_Q10 > (0x80000000 >> 4) {
						return 0x80000000 >> 4
					}
					if LPC_pred_Q10 < opus_int32(silk_int32_MAX>>4) {
						return opus_int32(silk_int32_MAX >> 4)
					}
					return LPC_pred_Q10
				}
				if LPC_pred_Q10 > opus_int32(silk_int32_MAX>>4) {
					return opus_int32(silk_int32_MAX >> 4)
				}
				if LPC_pred_Q10 < (0x80000000 >> 4) {
					return 0x80000000 >> 4
				}
				return LPC_pred_Q10
			}())<<4))) & 0x80000000) == 0 {
				if (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) & (opus_int32(opus_uint32(func() opus_int32 {
					if (0x80000000 >> 4) > (silk_int32_MAX >> 4) {
						if LPC_pred_Q10 > (0x80000000 >> 4) {
							return 0x80000000 >> 4
						}
						if LPC_pred_Q10 < opus_int32(silk_int32_MAX>>4) {
							return opus_int32(silk_int32_MAX >> 4)
						}
						return LPC_pred_Q10
					}
					if LPC_pred_Q10 > opus_int32(silk_int32_MAX>>4) {
						return opus_int32(silk_int32_MAX >> 4)
					}
					if LPC_pred_Q10 < (0x80000000 >> 4) {
						return 0x80000000 >> 4
					}
					return LPC_pred_Q10
				}()) << 4))) & 0x80000000) != 0 {
					*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i))) = 0x80000000
				} else {
					*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i))) = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) + (opus_int32(opus_uint32(func() opus_int32 {
						if (0x80000000 >> 4) > (silk_int32_MAX >> 4) {
							if LPC_pred_Q10 > (0x80000000 >> 4) {
								return 0x80000000 >> 4
							}
							if LPC_pred_Q10 < opus_int32(silk_int32_MAX>>4) {
								return opus_int32(silk_int32_MAX >> 4)
							}
							return LPC_pred_Q10
						}
						if LPC_pred_Q10 > opus_int32(silk_int32_MAX>>4) {
							return opus_int32(silk_int32_MAX >> 4)
						}
						if LPC_pred_Q10 < (0x80000000 >> 4) {
							return 0x80000000 >> 4
						}
						return LPC_pred_Q10
					}()) << 4))
				}
			} else if (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) | (opus_int32(opus_uint32(func() opus_int32 {
				if (0x80000000 >> 4) > (silk_int32_MAX >> 4) {
					if LPC_pred_Q10 > (0x80000000 >> 4) {
						return 0x80000000 >> 4
					}
					if LPC_pred_Q10 < opus_int32(silk_int32_MAX>>4) {
						return opus_int32(silk_int32_MAX >> 4)
					}
					return LPC_pred_Q10
				}
				if LPC_pred_Q10 > opus_int32(silk_int32_MAX>>4) {
					return opus_int32(silk_int32_MAX >> 4)
				}
				if LPC_pred_Q10 < (0x80000000 >> 4) {
					return 0x80000000 >> 4
				}
				return LPC_pred_Q10
			}()) << 4))) & 0x80000000) == 0 {
				*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i))) = silk_int32_MAX
			} else {
				*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i))) = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) + (opus_int32(opus_uint32(func() opus_int32 {
					if (0x80000000 >> 4) > (silk_int32_MAX >> 4) {
						if LPC_pred_Q10 > (0x80000000 >> 4) {
							return 0x80000000 >> 4
						}
						if LPC_pred_Q10 < opus_int32(silk_int32_MAX>>4) {
							return opus_int32(silk_int32_MAX >> 4)
						}
						return LPC_pred_Q10
					}
					if LPC_pred_Q10 > opus_int32(silk_int32_MAX>>4) {
						return opus_int32(silk_int32_MAX >> 4)
					}
					if LPC_pred_Q10 < (0x80000000 >> 4) {
						return 0x80000000 >> 4
					}
					return LPC_pred_Q10
				}()) << 4))
			}
			if ((opus_int32(frame[i])) + (func() opus_int32 {
				if (func() opus_int32 {
					if 8 == 1 {
						return ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(gain_Q10)) >> 16)) >> 1) + ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(gain_Q10)) >> 16)) & 1)
					}
					return (((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(gain_Q10)) >> 16)) >> (8 - 1)) + 1) >> 1
				}()) > silk_int16_MAX {
					return silk_int16_MAX
				}
				if (func() opus_int32 {
					if 8 == 1 {
						return ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(gain_Q10)) >> 16)) >> 1) + ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(gain_Q10)) >> 16)) & 1)
					}
					return (((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(gain_Q10)) >> 16)) >> (8 - 1)) + 1) >> 1
				}()) < opus_int32(math.MinInt16) {
					return math.MinInt16
				}
				if 8 == 1 {
					return ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(gain_Q10)) >> 16)) >> 1) + ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(gain_Q10)) >> 16)) & 1)
				}
				return (((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(gain_Q10)) >> 16)) >> (8 - 1)) + 1) >> 1
			}())) > silk_int16_MAX {
				frame[i] = silk_int16_MAX
			} else if ((opus_int32(frame[i])) + (func() opus_int32 {
				if (func() opus_int32 {
					if 8 == 1 {
						return ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(gain_Q10)) >> 16)) >> 1) + ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(gain_Q10)) >> 16)) & 1)
					}
					return (((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(gain_Q10)) >> 16)) >> (8 - 1)) + 1) >> 1
				}()) > silk_int16_MAX {
					return silk_int16_MAX
				}
				if (func() opus_int32 {
					if 8 == 1 {
						return ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(gain_Q10)) >> 16)) >> 1) + ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(gain_Q10)) >> 16)) & 1)
					}
					return (((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(gain_Q10)) >> 16)) >> (8 - 1)) + 1) >> 1
				}()) < opus_int32(math.MinInt16) {
					return math.MinInt16
				}
				if 8 == 1 {
					return ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(gain_Q10)) >> 16)) >> 1) + ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(gain_Q10)) >> 16)) & 1)
				}
				return (((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(gain_Q10)) >> 16)) >> (8 - 1)) + 1) >> 1
			}())) < opus_int32(math.MinInt16) {
				frame[i] = math.MinInt16
			} else {
				frame[i] = opus_int16((opus_int32(frame[i])) + (func() opus_int32 {
					if (func() opus_int32 {
						if 8 == 1 {
							return ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(gain_Q10)) >> 16)) >> 1) + ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(gain_Q10)) >> 16)) & 1)
						}
						return (((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(gain_Q10)) >> 16)) >> (8 - 1)) + 1) >> 1
					}()) > silk_int16_MAX {
						return silk_int16_MAX
					}
					if (func() opus_int32 {
						if 8 == 1 {
							return ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(gain_Q10)) >> 16)) >> 1) + ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(gain_Q10)) >> 16)) & 1)
						}
						return (((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(gain_Q10)) >> 16)) >> (8 - 1)) + 1) >> 1
					}()) < opus_int32(math.MinInt16) {
						return math.MinInt16
					}
					if 8 == 1 {
						return ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(gain_Q10)) >> 16)) >> 1) + ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(gain_Q10)) >> 16)) & 1)
					}
					return (((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(gain_Q10)) >> 16)) >> (8 - 1)) + 1) >> 1
				}()))
			}
		}
		libc.MemCpy(unsafe.Pointer(&psCNG.CNG_synth_state[0]), unsafe.Pointer((*opus_int32)(unsafe.Add(unsafe.Pointer(CNG_sig_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(length)))), int(MAX_LPC_ORDER*unsafe.Sizeof(opus_int32(0))))
	} else {
		libc.MemSet(unsafe.Pointer(&psCNG.CNG_synth_state[0]), 0, int(psDec.LPC_order*int64(unsafe.Sizeof(opus_int32(0)))))
	}
}
