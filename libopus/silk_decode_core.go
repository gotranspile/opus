package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

func silk_decode_core(psDec *silk_decoder_state, psDecCtrl *silk_decoder_control, xq [0]opus_int16, pulses [320]opus_int16, arch int64) {
	var (
		i                       int64
		k                       int64
		lag                     int64 = 0
		start_idx               int64
		sLTP_buf_idx            int64
		NLSF_interpolation_flag int64
		signalType              int64
		A_Q12                   *opus_int16
		B_Q14                   *opus_int16
		pxq                     *opus_int16
		A_Q12_tmp               [16]opus_int16
		sLTP                    *opus_int16
		sLTP_Q15                *opus_int32
		LTP_pred_Q13            opus_int32
		LPC_pred_Q10            opus_int32
		Gain_Q10                opus_int32
		inv_gain_Q31            opus_int32
		gain_adj_Q16            opus_int32
		rand_seed               opus_int32
		offset_Q10              opus_int32
		pred_lag_ptr            *opus_int32
		pexc_Q14                *opus_int32
		pres_Q14                *opus_int32
		res_Q14                 *opus_int32
		sLPC_Q14                *opus_int32
	)
	sLTP = (*opus_int16)(libc.Malloc(int(psDec.Ltp_mem_length * int64(unsafe.Sizeof(opus_int16(0))))))
	sLTP_Q15 = (*opus_int32)(libc.Malloc(int((psDec.Ltp_mem_length + psDec.Frame_length) * int64(unsafe.Sizeof(opus_int32(0))))))
	res_Q14 = (*opus_int32)(libc.Malloc(int(psDec.Subfr_length * int64(unsafe.Sizeof(opus_int32(0))))))
	sLPC_Q14 = (*opus_int32)(libc.Malloc(int((psDec.Subfr_length + MAX_LPC_ORDER) * int64(unsafe.Sizeof(opus_int32(0))))))
	offset_Q10 = opus_int32(silk_Quantization_Offsets_Q10[int64(psDec.Indices.SignalType)>>1][psDec.Indices.QuantOffsetType])
	if int64(psDec.Indices.NLSFInterpCoef_Q2) < 1<<2 {
		NLSF_interpolation_flag = 1
	} else {
		NLSF_interpolation_flag = 0
	}
	rand_seed = opus_int32(psDec.Indices.Seed)
	for i = 0; i < psDec.Frame_length; i++ {
		rand_seed = opus_int32(RAND_INCREMENT + opus_uint32(rand_seed)*RAND_MULTIPLIER)
		psDec.Exc_Q14[i] = opus_int32(opus_uint32(opus_int32(pulses[i])) << 14)
		if psDec.Exc_Q14[i] > 0 {
			psDec.Exc_Q14[i] -= opus_int32(QUANT_LEVEL_ADJUST_Q10 << 4)
		} else if psDec.Exc_Q14[i] < 0 {
			psDec.Exc_Q14[i] += opus_int32(QUANT_LEVEL_ADJUST_Q10 << 4)
		}
		psDec.Exc_Q14[i] += offset_Q10 << 4
		if rand_seed < 0 {
			psDec.Exc_Q14[i] = -psDec.Exc_Q14[i]
		}
		rand_seed = opus_int32(opus_uint32(rand_seed) + opus_uint32(pulses[i]))
	}
	libc.MemCpy(unsafe.Pointer(sLPC_Q14), unsafe.Pointer(&psDec.SLPC_Q14_buf[0]), int(MAX_LPC_ORDER*unsafe.Sizeof(opus_int32(0))))
	pexc_Q14 = &psDec.Exc_Q14[0]
	pxq = &xq[0]
	sLTP_buf_idx = psDec.Ltp_mem_length
	for k = 0; k < psDec.Nb_subfr; k++ {
		pres_Q14 = res_Q14
		A_Q12 = &psDecCtrl.PredCoef_Q12[k>>1][0]
		libc.MemCpy(unsafe.Pointer(&A_Q12_tmp[0]), unsafe.Pointer(A_Q12), int(psDec.LPC_order*int64(unsafe.Sizeof(opus_int16(0)))))
		B_Q14 = &psDecCtrl.LTPCoef_Q14[k*LTP_ORDER]
		signalType = int64(psDec.Indices.SignalType)
		Gain_Q10 = (psDecCtrl.Gains_Q16[k]) >> 6
		inv_gain_Q31 = silk_INVERSE32_varQ(psDecCtrl.Gains_Q16[k], 47)
		if psDecCtrl.Gains_Q16[k] != psDec.Prev_gain_Q16 {
			gain_adj_Q16 = silk_DIV32_varQ(psDec.Prev_gain_Q16, psDecCtrl.Gains_Q16[k], 16)
			for i = 0; i < MAX_LPC_ORDER; i++ {
				*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(i))) = opus_int32((int64(gain_adj_Q16) * int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(i))))) >> 16)
			}
		} else {
			gain_adj_Q16 = 1 << 16
		}
		psDec.Prev_gain_Q16 = psDecCtrl.Gains_Q16[k]
		if psDec.LossCnt != 0 && psDec.PrevSignalType == TYPE_VOICED && int64(psDec.Indices.SignalType) != TYPE_VOICED && k < MAX_NB_SUBFR/2 {
			libc.MemSet(unsafe.Pointer(B_Q14), 0, int(LTP_ORDER*unsafe.Sizeof(opus_int16(0))))
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(B_Q14), unsafe.Sizeof(opus_int16(0))*uintptr(LTP_ORDER/2))) = opus_int16(opus_int32(0.25*(1<<14) + 0.5))
			signalType = TYPE_VOICED
			psDecCtrl.PitchL[k] = psDec.LagPrev
		}
		if signalType == TYPE_VOICED {
			lag = psDecCtrl.PitchL[k]
			if k == 0 || k == 2 && NLSF_interpolation_flag != 0 {
				start_idx = psDec.Ltp_mem_length - lag - psDec.LPC_order - LTP_ORDER/2
				if k == 2 {
					libc.MemCpy(unsafe.Pointer(&psDec.OutBuf[psDec.Ltp_mem_length]), unsafe.Pointer(&xq[0]), int(psDec.Subfr_length*2*int64(unsafe.Sizeof(opus_int16(0)))))
				}
				silk_LPC_analysis_filter((*opus_int16)(unsafe.Add(unsafe.Pointer(sLTP), unsafe.Sizeof(opus_int16(0))*uintptr(start_idx))), &psDec.OutBuf[start_idx+k*psDec.Subfr_length], A_Q12, opus_int32(psDec.Ltp_mem_length-start_idx), opus_int32(psDec.LPC_order), arch)
				if k == 0 {
					inv_gain_Q31 = opus_int32(opus_uint32((inv_gain_Q31*opus_int32(int64(opus_int16(psDecCtrl.LTP_scale_Q14))))>>16) << 2)
				}
				for i = 0; i < lag+LTP_ORDER/2; i++ {
					*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLTP_Q15), unsafe.Sizeof(opus_int32(0))*uintptr(sLTP_buf_idx-i-1))) = (inv_gain_Q31 * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(sLTP), unsafe.Sizeof(opus_int16(0))*uintptr(psDec.Ltp_mem_length-i-1)))))) >> 16
				}
			} else {
				if gain_adj_Q16 != 1<<16 {
					for i = 0; i < lag+LTP_ORDER/2; i++ {
						*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLTP_Q15), unsafe.Sizeof(opus_int32(0))*uintptr(sLTP_buf_idx-i-1))) = opus_int32((int64(gain_adj_Q16) * int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLTP_Q15), unsafe.Sizeof(opus_int32(0))*uintptr(sLTP_buf_idx-i-1))))) >> 16)
					}
				}
			}
		}
		if signalType == TYPE_VOICED {
			pred_lag_ptr = (*opus_int32)(unsafe.Add(unsafe.Pointer(sLTP_Q15), unsafe.Sizeof(opus_int32(0))*uintptr(sLTP_buf_idx-lag+LTP_ORDER/2)))
			for i = 0; i < psDec.Subfr_length; i++ {
				LTP_pred_Q13 = 2
				LTP_pred_Q13 = LTP_pred_Q13 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), unsafe.Sizeof(opus_int32(0))*0))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(B_Q14), unsafe.Sizeof(opus_int16(0))*0))))) >> 16)
				LTP_pred_Q13 = LTP_pred_Q13 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), -int(unsafe.Sizeof(opus_int32(0))*1)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(B_Q14), unsafe.Sizeof(opus_int16(0))*1))))) >> 16)
				LTP_pred_Q13 = LTP_pred_Q13 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), -int(unsafe.Sizeof(opus_int32(0))*2)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(B_Q14), unsafe.Sizeof(opus_int16(0))*2))))) >> 16)
				LTP_pred_Q13 = LTP_pred_Q13 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), -int(unsafe.Sizeof(opus_int32(0))*3)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(B_Q14), unsafe.Sizeof(opus_int16(0))*3))))) >> 16)
				LTP_pred_Q13 = LTP_pred_Q13 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), -int(unsafe.Sizeof(opus_int32(0))*4)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(B_Q14), unsafe.Sizeof(opus_int16(0))*4))))) >> 16)
				pred_lag_ptr = (*opus_int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), unsafe.Sizeof(opus_int32(0))*1))
				*(*opus_int32)(unsafe.Add(unsafe.Pointer(pres_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(i))) = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(pexc_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(i)))) + (opus_int32(opus_uint32(LTP_pred_Q13) << 1))
				*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLTP_Q15), unsafe.Sizeof(opus_int32(0))*uintptr(sLTP_buf_idx))) = opus_int32(opus_uint32(*(*opus_int32)(unsafe.Add(unsafe.Pointer(pres_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(i)))) << 1)
				sLTP_buf_idx++
			}
		} else {
			pres_Q14 = pexc_Q14
		}
		for i = 0; i < psDec.Subfr_length; i++ {
			LPC_pred_Q10 = opus_int32(psDec.LPC_order >> 1)
			LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-1)))) * opus_int32(int64(A_Q12_tmp[0]))) >> 16)
			LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-2)))) * opus_int32(int64(A_Q12_tmp[1]))) >> 16)
			LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-3)))) * opus_int32(int64(A_Q12_tmp[2]))) >> 16)
			LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-4)))) * opus_int32(int64(A_Q12_tmp[3]))) >> 16)
			LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-5)))) * opus_int32(int64(A_Q12_tmp[4]))) >> 16)
			LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-6)))) * opus_int32(int64(A_Q12_tmp[5]))) >> 16)
			LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-7)))) * opus_int32(int64(A_Q12_tmp[6]))) >> 16)
			LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-8)))) * opus_int32(int64(A_Q12_tmp[7]))) >> 16)
			LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-9)))) * opus_int32(int64(A_Q12_tmp[8]))) >> 16)
			LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-10)))) * opus_int32(int64(A_Q12_tmp[9]))) >> 16)
			if psDec.LPC_order == 16 {
				LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-11)))) * opus_int32(int64(A_Q12_tmp[10]))) >> 16)
				LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-12)))) * opus_int32(int64(A_Q12_tmp[11]))) >> 16)
				LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-13)))) * opus_int32(int64(A_Q12_tmp[12]))) >> 16)
				LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-14)))) * opus_int32(int64(A_Q12_tmp[13]))) >> 16)
				LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-15)))) * opus_int32(int64(A_Q12_tmp[14]))) >> 16)
				LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-16)))) * opus_int32(int64(A_Q12_tmp[15]))) >> 16)
			}
			if ((opus_uint32(*(*opus_int32)(unsafe.Add(unsafe.Pointer(pres_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(i)))) + opus_uint32(opus_int32(opus_uint32(func() opus_int32 {
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
				if (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(pres_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(i)))) & (opus_int32(opus_uint32(func() opus_int32 {
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
					*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i))) = 0x80000000
				} else {
					*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i))) = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(pres_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(i)))) + (opus_int32(opus_uint32(func() opus_int32 {
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
			} else if (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(pres_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(i)))) | (opus_int32(opus_uint32(func() opus_int32 {
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
				*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i))) = silk_int32_MAX
			} else {
				*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i))) = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(pres_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(i)))) + (opus_int32(opus_uint32(func() opus_int32 {
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
			if (func() opus_int32 {
				if 8 == 1 {
					return ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(Gain_Q10)) >> 16)) >> 1) + ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(Gain_Q10)) >> 16)) & 1)
				}
				return (((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(Gain_Q10)) >> 16)) >> (8 - 1)) + 1) >> 1
			}()) > silk_int16_MAX {
				*(*opus_int16)(unsafe.Add(unsafe.Pointer(pxq), unsafe.Sizeof(opus_int16(0))*uintptr(i))) = silk_int16_MAX
			} else if (func() opus_int32 {
				if 8 == 1 {
					return ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(Gain_Q10)) >> 16)) >> 1) + ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(Gain_Q10)) >> 16)) & 1)
				}
				return (((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(Gain_Q10)) >> 16)) >> (8 - 1)) + 1) >> 1
			}()) < opus_int32(math.MinInt16) {
				*(*opus_int16)(unsafe.Add(unsafe.Pointer(pxq), unsafe.Sizeof(opus_int16(0))*uintptr(i))) = math.MinInt16
			} else if 8 == 1 {
				*(*opus_int16)(unsafe.Add(unsafe.Pointer(pxq), unsafe.Sizeof(opus_int16(0))*uintptr(i))) = opus_int16(((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(Gain_Q10)) >> 16)) >> 1) + ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(Gain_Q10)) >> 16)) & 1))
			} else {
				*(*opus_int16)(unsafe.Add(unsafe.Pointer(pxq), unsafe.Sizeof(opus_int16(0))*uintptr(i))) = opus_int16((((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(Gain_Q10)) >> 16)) >> (8 - 1)) + 1) >> 1)
			}
		}
		libc.MemCpy(unsafe.Pointer(sLPC_Q14), unsafe.Pointer((*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(psDec.Subfr_length)))), int(MAX_LPC_ORDER*unsafe.Sizeof(opus_int32(0))))
		pexc_Q14 = (*opus_int32)(unsafe.Add(unsafe.Pointer(pexc_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(psDec.Subfr_length)))
		pxq = (*opus_int16)(unsafe.Add(unsafe.Pointer(pxq), unsafe.Sizeof(opus_int16(0))*uintptr(psDec.Subfr_length)))
	}
	libc.MemCpy(unsafe.Pointer(&psDec.SLPC_Q14_buf[0]), unsafe.Pointer(sLPC_Q14), int(MAX_LPC_ORDER*unsafe.Sizeof(opus_int32(0))))
}
