package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

const BWE_COEF = 0.99
const V_PITCH_GAIN_START_MIN_Q14 = 11469
const NB_ATT = 2
const V_PITCH_GAIN_START_MAX_Q14 = 15565
const MAX_PITCH_LAG_MS = 18
const RAND_BUF_SIZE = 128
const RAND_BUF_MASK uint8 = math.MaxInt8
const LOG2_INV_LPC_GAIN_HIGH_THRES = 3
const LOG2_INV_LPC_GAIN_LOW_THRES = 8
const PITCH_DRIFT_FAC_Q16 = 655

var HARM_ATT_Q15 [2]opus_int16 = [2]opus_int16{32440, 31130}
var PLC_RAND_ATTENUATE_V_Q15 [2]opus_int16 = [2]opus_int16{31130, 26214}
var PLC_RAND_ATTENUATE_UV_Q15 [2]opus_int16 = [2]opus_int16{32440, 29491}

func silk_PLC_Reset(psDec *silk_decoder_state) {
	psDec.SPLC.PitchL_Q8 = opus_int32(opus_uint32(psDec.Frame_length) << (8 - 1))
	psDec.SPLC.PrevGain_Q16[0] = opus_int32(1*(1<<16) + 0.5)
	psDec.SPLC.PrevGain_Q16[1] = opus_int32(1*(1<<16) + 0.5)
	psDec.SPLC.Subfr_length = 20
	psDec.SPLC.Nb_subfr = 2
}
func silk_PLC(psDec *silk_decoder_state, psDecCtrl *silk_decoder_control, frame [0]opus_int16, lost int64, arch int64) {
	if psDec.Fs_kHz != psDec.SPLC.Fs_kHz {
		silk_PLC_Reset(psDec)
		psDec.SPLC.Fs_kHz = psDec.Fs_kHz
	}
	if lost != 0 {
		silk_PLC_conceal(psDec, psDecCtrl, frame, arch)
		psDec.LossCnt++
	} else {
		silk_PLC_update(psDec, psDecCtrl)
	}
}
func silk_PLC_update(psDec *silk_decoder_state, psDecCtrl *silk_decoder_control) {
	var (
		LTP_Gain_Q14      opus_int32
		temp_LTP_Gain_Q14 opus_int32
		i                 int64
		j                 int64
		psPLC             *silk_PLC_struct
	)
	psPLC = &psDec.SPLC
	psDec.PrevSignalType = int64(psDec.Indices.SignalType)
	LTP_Gain_Q14 = 0
	if int64(psDec.Indices.SignalType) == TYPE_VOICED {
		for j = 0; j*psDec.Subfr_length < psDecCtrl.PitchL[psDec.Nb_subfr-1]; j++ {
			if j == psDec.Nb_subfr {
				break
			}
			temp_LTP_Gain_Q14 = 0
			for i = 0; i < LTP_ORDER; i++ {
				temp_LTP_Gain_Q14 += opus_int32(psDecCtrl.LTPCoef_Q14[(psDec.Nb_subfr-1-j)*LTP_ORDER+i])
			}
			if temp_LTP_Gain_Q14 > LTP_Gain_Q14 {
				LTP_Gain_Q14 = temp_LTP_Gain_Q14
				libc.MemCpy(unsafe.Pointer(&psPLC.LTPCoef_Q14[0]), unsafe.Pointer(&psDecCtrl.LTPCoef_Q14[opus_int32(opus_int16(psDec.Nb_subfr-1-j))*LTP_ORDER]), int(LTP_ORDER*unsafe.Sizeof(opus_int16(0))))
				psPLC.PitchL_Q8 = opus_int32(opus_uint32(psDecCtrl.PitchL[psDec.Nb_subfr-1-j]) << 8)
			}
		}
		libc.MemSet(unsafe.Pointer(&psPLC.LTPCoef_Q14[0]), 0, int(LTP_ORDER*unsafe.Sizeof(opus_int16(0))))
		psPLC.LTPCoef_Q14[LTP_ORDER/2] = opus_int16(LTP_Gain_Q14)
		if LTP_Gain_Q14 < V_PITCH_GAIN_START_MIN_Q14 {
			var (
				scale_Q10 int64
				tmp       opus_int32
			)
			tmp = opus_int32(V_PITCH_GAIN_START_MIN_Q14 << 10)
			scale_Q10 = int64(tmp / (func() opus_int32 {
				if LTP_Gain_Q14 > 1 {
					return LTP_Gain_Q14
				}
				return 1
			}()))
			for i = 0; i < LTP_ORDER; i++ {
				psPLC.LTPCoef_Q14[i] = opus_int16((opus_int32(psPLC.LTPCoef_Q14[i]) * opus_int32(opus_int16(scale_Q10))) >> 10)
			}
		} else if LTP_Gain_Q14 > V_PITCH_GAIN_START_MAX_Q14 {
			var (
				scale_Q14 int64
				tmp       opus_int32
			)
			tmp = opus_int32(V_PITCH_GAIN_START_MAX_Q14 << 14)
			scale_Q14 = int64(tmp / (func() opus_int32 {
				if LTP_Gain_Q14 > 1 {
					return LTP_Gain_Q14
				}
				return 1
			}()))
			for i = 0; i < LTP_ORDER; i++ {
				psPLC.LTPCoef_Q14[i] = opus_int16((opus_int32(psPLC.LTPCoef_Q14[i]) * opus_int32(opus_int16(scale_Q14))) >> 14)
			}
		}
	} else {
		psPLC.PitchL_Q8 = opus_int32(opus_uint32(opus_int32(opus_int16(psDec.Fs_kHz))*18) << 8)
		libc.MemSet(unsafe.Pointer(&psPLC.LTPCoef_Q14[0]), 0, int(LTP_ORDER*unsafe.Sizeof(opus_int16(0))))
	}
	libc.MemCpy(unsafe.Pointer(&psPLC.PrevLPC_Q12[0]), unsafe.Pointer(&(psDecCtrl.PredCoef_Q12[1])[0]), int(psDec.LPC_order*int64(unsafe.Sizeof(opus_int16(0)))))
	psPLC.PrevLTP_scale_Q14 = opus_int16(psDecCtrl.LTP_scale_Q14)
	libc.MemCpy(unsafe.Pointer(&psPLC.PrevGain_Q16[0]), unsafe.Pointer(&psDecCtrl.Gains_Q16[psDec.Nb_subfr-2]), int(2*unsafe.Sizeof(opus_int32(0))))
	psPLC.Subfr_length = psDec.Subfr_length
	psPLC.Nb_subfr = psDec.Nb_subfr
}
func silk_PLC_energy(energy1 *opus_int32, shift1 *int64, energy2 *opus_int32, shift2 *int64, exc_Q14 *opus_int32, prevGain_Q10 *opus_int32, subfr_length int64, nb_subfr int64) {
	var (
		i           int64
		k           int64
		exc_buf     *opus_int16
		exc_buf_ptr *opus_int16
	)
	exc_buf = (*opus_int16)(libc.Malloc(int((subfr_length * 2) * int64(unsafe.Sizeof(opus_int16(0))))))
	exc_buf_ptr = exc_buf
	for k = 0; k < 2; k++ {
		for i = 0; i < subfr_length; i++ {
			if ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(exc_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(i+(k+nb_subfr-2)*subfr_length)))) * int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(prevGain_Q10), unsafe.Sizeof(opus_int32(0))*uintptr(k))))) >> 16)) >> 8) > silk_int16_MAX {
				*(*opus_int16)(unsafe.Add(unsafe.Pointer(exc_buf_ptr), unsafe.Sizeof(opus_int16(0))*uintptr(i))) = silk_int16_MAX
			} else if ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(exc_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(i+(k+nb_subfr-2)*subfr_length)))) * int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(prevGain_Q10), unsafe.Sizeof(opus_int32(0))*uintptr(k))))) >> 16)) >> 8) < opus_int32(math.MinInt16) {
				*(*opus_int16)(unsafe.Add(unsafe.Pointer(exc_buf_ptr), unsafe.Sizeof(opus_int16(0))*uintptr(i))) = math.MinInt16
			} else {
				*(*opus_int16)(unsafe.Add(unsafe.Pointer(exc_buf_ptr), unsafe.Sizeof(opus_int16(0))*uintptr(i))) = opus_int16((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(exc_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(i+(k+nb_subfr-2)*subfr_length)))) * int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(prevGain_Q10), unsafe.Sizeof(opus_int32(0))*uintptr(k))))) >> 16)) >> 8)
			}
		}
		exc_buf_ptr = (*opus_int16)(unsafe.Add(unsafe.Pointer(exc_buf_ptr), unsafe.Sizeof(opus_int16(0))*uintptr(subfr_length)))
	}
	silk_sum_sqr_shift(energy1, shift1, exc_buf, subfr_length)
	silk_sum_sqr_shift(energy2, shift2, (*opus_int16)(unsafe.Add(unsafe.Pointer(exc_buf), unsafe.Sizeof(opus_int16(0))*uintptr(subfr_length))), subfr_length)
}
func silk_PLC_conceal(psDec *silk_decoder_state, psDecCtrl *silk_decoder_control, frame [0]opus_int16, arch int64) {
	var (
		i              int64
		j              int64
		k              int64
		lag            int64
		idx            int64
		sLTP_buf_idx   int64
		shift1         int64
		shift2         int64
		rand_seed      opus_int32
		harm_Gain_Q15  opus_int32
		rand_Gain_Q15  opus_int32
		inv_gain_Q30   opus_int32
		energy1        opus_int32
		energy2        opus_int32
		rand_ptr       *opus_int32
		pred_lag_ptr   *opus_int32
		LPC_pred_Q10   opus_int32
		LTP_pred_Q12   opus_int32
		rand_scale_Q14 opus_int16
		B_Q14          *opus_int16
		sLPC_Q14_ptr   *opus_int32
		A_Q12          [16]opus_int16
		sLTP           *opus_int16
		sLTP_Q14       *opus_int32
		psPLC          *silk_PLC_struct = &psDec.SPLC
		prevGain_Q10   [2]opus_int32
	)
	sLTP_Q14 = (*opus_int32)(libc.Malloc(int((psDec.Ltp_mem_length + psDec.Frame_length) * int64(unsafe.Sizeof(opus_int32(0))))))
	sLTP = (*opus_int16)(libc.Malloc(int(psDec.Ltp_mem_length * int64(unsafe.Sizeof(opus_int16(0))))))
	prevGain_Q10[0] = (psPLC.PrevGain_Q16[0]) >> 6
	prevGain_Q10[1] = (psPLC.PrevGain_Q16[1]) >> 6
	if psDec.First_frame_after_reset != 0 {
		*(*[16]opus_int16)(unsafe.Pointer(&psPLC.PrevLPC_Q12[0])) = [16]opus_int16{}
	}
	silk_PLC_energy(&energy1, &shift1, &energy2, &shift2, &psDec.Exc_Q14[0], &prevGain_Q10[0], psDec.Subfr_length, psDec.Nb_subfr)
	if (energy1 >> opus_int32(shift2)) < (energy2 >> opus_int32(shift1)) {
		rand_ptr = &psDec.Exc_Q14[silk_max_int(0, (psPLC.Nb_subfr-1)*psPLC.Subfr_length-RAND_BUF_SIZE)]
	} else {
		rand_ptr = &psDec.Exc_Q14[silk_max_int(0, psPLC.Nb_subfr*psPLC.Subfr_length-RAND_BUF_SIZE)]
	}
	B_Q14 = &psPLC.LTPCoef_Q14[0]
	rand_scale_Q14 = psPLC.RandScale_Q14
	harm_Gain_Q15 = opus_int32(HARM_ATT_Q15[silk_min_int(NB_ATT-1, psDec.LossCnt)])
	if psDec.PrevSignalType == TYPE_VOICED {
		rand_Gain_Q15 = opus_int32(PLC_RAND_ATTENUATE_V_Q15[silk_min_int(NB_ATT-1, psDec.LossCnt)])
	} else {
		rand_Gain_Q15 = opus_int32(PLC_RAND_ATTENUATE_UV_Q15[silk_min_int(NB_ATT-1, psDec.LossCnt)])
	}
	silk_bwexpander(&psPLC.PrevLPC_Q12[0], psDec.LPC_order, opus_int32(BWE_COEF*(1<<16)+0.5))
	libc.MemCpy(unsafe.Pointer(&A_Q12[0]), unsafe.Pointer(&psPLC.PrevLPC_Q12[0]), int(psDec.LPC_order*int64(unsafe.Sizeof(opus_int16(0)))))
	if psDec.LossCnt == 0 {
		rand_scale_Q14 = 1 << 14
		if psDec.PrevSignalType == TYPE_VOICED {
			for i = 0; i < LTP_ORDER; i++ {
				rand_scale_Q14 -= *(*opus_int16)(unsafe.Add(unsafe.Pointer(B_Q14), unsafe.Sizeof(opus_int16(0))*uintptr(i)))
			}
			rand_scale_Q14 = silk_max_16(3277, rand_scale_Q14)
			rand_scale_Q14 = opus_int16((opus_int32(rand_scale_Q14) * opus_int32(psPLC.PrevLTP_scale_Q14)) >> 14)
		} else {
			var (
				invGain_Q30    opus_int32
				down_scale_Q30 opus_int32
			)
			invGain_Q30 = func() opus_int32 {
				_ = arch
				return silk_LPC_inverse_pred_gain_c(&psPLC.PrevLPC_Q12[0], psDec.LPC_order)
			}()
			down_scale_Q30 = silk_min_32(opus_int32((1<<30)>>LOG2_INV_LPC_GAIN_HIGH_THRES), invGain_Q30)
			down_scale_Q30 = silk_max_32(opus_int32((1<<30)>>LOG2_INV_LPC_GAIN_LOW_THRES), down_scale_Q30)
			down_scale_Q30 = opus_int32(opus_uint32(down_scale_Q30) << LOG2_INV_LPC_GAIN_HIGH_THRES)
			rand_Gain_Q15 = ((down_scale_Q30 * opus_int32(int64(opus_int16(rand_Gain_Q15)))) >> 16) >> 14
		}
	}
	rand_seed = psPLC.Rand_seed
	if 8 == 1 {
		lag = int64((psPLC.PitchL_Q8 >> 1) + (psPLC.PitchL_Q8 & 1))
	} else {
		lag = int64(((psPLC.PitchL_Q8 >> (8 - 1)) + 1) >> 1)
	}
	sLTP_buf_idx = psDec.Ltp_mem_length
	idx = psDec.Ltp_mem_length - lag - psDec.LPC_order - LTP_ORDER/2
	silk_LPC_analysis_filter((*opus_int16)(unsafe.Add(unsafe.Pointer(sLTP), unsafe.Sizeof(opus_int16(0))*uintptr(idx))), &psDec.OutBuf[idx], &A_Q12[0], opus_int32(psDec.Ltp_mem_length-idx), opus_int32(psDec.LPC_order), arch)
	inv_gain_Q30 = silk_INVERSE32_varQ(psPLC.PrevGain_Q16[1], 46)
	if inv_gain_Q30 < opus_int32(silk_int32_MAX>>1) {
		inv_gain_Q30 = inv_gain_Q30
	} else {
		inv_gain_Q30 = opus_int32(silk_int32_MAX >> 1)
	}
	for i = idx + psDec.LPC_order; i < psDec.Ltp_mem_length; i++ {
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLTP_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(i))) = (inv_gain_Q30 * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(sLTP), unsafe.Sizeof(opus_int16(0))*uintptr(i)))))) >> 16
	}
	for k = 0; k < psDec.Nb_subfr; k++ {
		pred_lag_ptr = (*opus_int32)(unsafe.Add(unsafe.Pointer(sLTP_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(sLTP_buf_idx-lag+LTP_ORDER/2)))
		for i = 0; i < psDec.Subfr_length; i++ {
			LTP_pred_Q12 = 2
			LTP_pred_Q12 = LTP_pred_Q12 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), unsafe.Sizeof(opus_int32(0))*0))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(B_Q14), unsafe.Sizeof(opus_int16(0))*0))))) >> 16)
			LTP_pred_Q12 = LTP_pred_Q12 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), -int(unsafe.Sizeof(opus_int32(0))*1)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(B_Q14), unsafe.Sizeof(opus_int16(0))*1))))) >> 16)
			LTP_pred_Q12 = LTP_pred_Q12 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), -int(unsafe.Sizeof(opus_int32(0))*2)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(B_Q14), unsafe.Sizeof(opus_int16(0))*2))))) >> 16)
			LTP_pred_Q12 = LTP_pred_Q12 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), -int(unsafe.Sizeof(opus_int32(0))*3)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(B_Q14), unsafe.Sizeof(opus_int16(0))*3))))) >> 16)
			LTP_pred_Q12 = LTP_pred_Q12 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), -int(unsafe.Sizeof(opus_int32(0))*4)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(B_Q14), unsafe.Sizeof(opus_int16(0))*4))))) >> 16)
			pred_lag_ptr = (*opus_int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), unsafe.Sizeof(opus_int32(0))*1))
			rand_seed = opus_int32(RAND_INCREMENT + opus_uint32(rand_seed)*RAND_MULTIPLIER)
			idx = int64((rand_seed >> 25) & opus_int32(RAND_BUF_SIZE-1))
			*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLTP_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(sLTP_buf_idx))) = opus_int32(opus_uint32(LTP_pred_Q12+(((*(*opus_int32)(unsafe.Add(unsafe.Pointer(rand_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(idx))))*opus_int32(int64(rand_scale_Q14)))>>16)) << 2)
			sLTP_buf_idx++
		}
		for j = 0; j < LTP_ORDER; j++ {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(B_Q14), unsafe.Sizeof(opus_int16(0))*uintptr(j))) = opus_int16((opus_int32(opus_int16(harm_Gain_Q15)) * opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(B_Q14), unsafe.Sizeof(opus_int16(0))*uintptr(j))))) >> 15)
		}
		rand_scale_Q14 = opus_int16((opus_int32(rand_scale_Q14) * opus_int32(opus_int16(rand_Gain_Q15))) >> 15)
		psPLC.PitchL_Q8 = psPLC.PitchL_Q8 + ((psPLC.PitchL_Q8 * PITCH_DRIFT_FAC_Q16) >> 16)
		psPLC.PitchL_Q8 = silk_min_32(psPLC.PitchL_Q8, opus_int32(opus_uint32(MAX_PITCH_LAG_MS*opus_int32(opus_int16(psDec.Fs_kHz)))<<8))
		if 8 == 1 {
			lag = int64((psPLC.PitchL_Q8 >> 1) + (psPLC.PitchL_Q8 & 1))
		} else {
			lag = int64(((psPLC.PitchL_Q8 >> (8 - 1)) + 1) >> 1)
		}
	}
	sLPC_Q14_ptr = (*opus_int32)(unsafe.Add(unsafe.Pointer(sLTP_Q14), unsafe.Sizeof(opus_int32(0))*uintptr(psDec.Ltp_mem_length-MAX_LPC_ORDER)))
	libc.MemCpy(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Pointer(&psDec.SLPC_Q14_buf[0]), int(MAX_LPC_ORDER*unsafe.Sizeof(opus_int32(0))))
	for i = 0; i < psDec.Frame_length; i++ {
		LPC_pred_Q10 = opus_int32(psDec.LPC_order >> 1)
		LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-1)))) * opus_int32(int64(A_Q12[0]))) >> 16)
		LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-2)))) * opus_int32(int64(A_Q12[1]))) >> 16)
		LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-3)))) * opus_int32(int64(A_Q12[2]))) >> 16)
		LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-4)))) * opus_int32(int64(A_Q12[3]))) >> 16)
		LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-5)))) * opus_int32(int64(A_Q12[4]))) >> 16)
		LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-6)))) * opus_int32(int64(A_Q12[5]))) >> 16)
		LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-7)))) * opus_int32(int64(A_Q12[6]))) >> 16)
		LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-8)))) * opus_int32(int64(A_Q12[7]))) >> 16)
		LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-9)))) * opus_int32(int64(A_Q12[8]))) >> 16)
		LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-10)))) * opus_int32(int64(A_Q12[9]))) >> 16)
		for j = 10; j < psDec.LPC_order; j++ {
			LPC_pred_Q10 = LPC_pred_Q10 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i-j-1)))) * opus_int32(int64(A_Q12[j]))) >> 16)
		}
		if ((opus_uint32(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) + opus_uint32(opus_int32(opus_uint32(func() opus_int32 {
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
			if (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) & (opus_int32(opus_uint32(func() opus_int32 {
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
				*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i))) = 0x80000000
			} else {
				*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i))) = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) + (opus_int32(opus_uint32(func() opus_int32 {
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
		} else if (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) | (opus_int32(opus_uint32(func() opus_int32 {
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
			*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i))) = silk_int32_MAX
		} else {
			*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i))) = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) + (opus_int32(opus_uint32(func() opus_int32 {
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
			if (func() opus_int32 {
				if 8 == 1 {
					return ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(prevGain_Q10[1])) >> 16)) >> 1) + ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(prevGain_Q10[1])) >> 16)) & 1)
				}
				return (((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(prevGain_Q10[1])) >> 16)) >> (8 - 1)) + 1) >> 1
			}()) > silk_int16_MAX {
				return silk_int16_MAX
			}
			if (func() opus_int32 {
				if 8 == 1 {
					return ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(prevGain_Q10[1])) >> 16)) >> 1) + ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(prevGain_Q10[1])) >> 16)) & 1)
				}
				return (((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(prevGain_Q10[1])) >> 16)) >> (8 - 1)) + 1) >> 1
			}()) < opus_int32(math.MinInt16) {
				return math.MinInt16
			}
			if 8 == 1 {
				return ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(prevGain_Q10[1])) >> 16)) >> 1) + ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(prevGain_Q10[1])) >> 16)) & 1)
			}
			return (((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(prevGain_Q10[1])) >> 16)) >> (8 - 1)) + 1) >> 1
		}()) > silk_int16_MAX {
			frame[i] = silk_int16_MAX
		} else if (func() opus_int32 {
			if (func() opus_int32 {
				if 8 == 1 {
					return ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(prevGain_Q10[1])) >> 16)) >> 1) + ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(prevGain_Q10[1])) >> 16)) & 1)
				}
				return (((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(prevGain_Q10[1])) >> 16)) >> (8 - 1)) + 1) >> 1
			}()) > silk_int16_MAX {
				return silk_int16_MAX
			}
			if (func() opus_int32 {
				if 8 == 1 {
					return ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(prevGain_Q10[1])) >> 16)) >> 1) + ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(prevGain_Q10[1])) >> 16)) & 1)
				}
				return (((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(prevGain_Q10[1])) >> 16)) >> (8 - 1)) + 1) >> 1
			}()) < opus_int32(math.MinInt16) {
				return math.MinInt16
			}
			if 8 == 1 {
				return ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(prevGain_Q10[1])) >> 16)) >> 1) + ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(prevGain_Q10[1])) >> 16)) & 1)
			}
			return (((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(prevGain_Q10[1])) >> 16)) >> (8 - 1)) + 1) >> 1
		}()) < opus_int32(math.MinInt16) {
			frame[i] = math.MinInt16
		} else if (func() opus_int32 {
			if 8 == 1 {
				return ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(prevGain_Q10[1])) >> 16)) >> 1) + ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(prevGain_Q10[1])) >> 16)) & 1)
			}
			return (((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(prevGain_Q10[1])) >> 16)) >> (8 - 1)) + 1) >> 1
		}()) > silk_int16_MAX {
			frame[i] = silk_int16_MAX
		} else if (func() opus_int32 {
			if 8 == 1 {
				return ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(prevGain_Q10[1])) >> 16)) >> 1) + ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(prevGain_Q10[1])) >> 16)) & 1)
			}
			return (((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(prevGain_Q10[1])) >> 16)) >> (8 - 1)) + 1) >> 1
		}()) < opus_int32(math.MinInt16) {
			frame[i] = math.MinInt16
		} else if 8 == 1 {
			frame[i] = opus_int16(((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(prevGain_Q10[1])) >> 16)) >> 1) + ((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(prevGain_Q10[1])) >> 16)) & 1))
		} else {
			frame[i] = opus_int16((((opus_int32((int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(MAX_LPC_ORDER+i)))) * int64(prevGain_Q10[1])) >> 16)) >> (8 - 1)) + 1) >> 1)
		}
	}
	libc.MemCpy(unsafe.Pointer(&psDec.SLPC_Q14_buf[0]), unsafe.Pointer((*opus_int32)(unsafe.Add(unsafe.Pointer(sLPC_Q14_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(psDec.Frame_length)))), int(MAX_LPC_ORDER*unsafe.Sizeof(opus_int32(0))))
	psPLC.Rand_seed = rand_seed
	psPLC.RandScale_Q14 = rand_scale_Q14
	for i = 0; i < MAX_NB_SUBFR; i++ {
		psDecCtrl.PitchL[i] = lag
	}
}
func silk_PLC_glue_frames(psDec *silk_decoder_state, frame [0]opus_int16, length int64) {
	var (
		i            int64
		energy_shift int64
		energy       opus_int32
		psPLC        *silk_PLC_struct
	)
	psPLC = &psDec.SPLC
	if psDec.LossCnt != 0 {
		silk_sum_sqr_shift(&psPLC.Conc_energy, &psPLC.Conc_energy_shift, &frame[0], length)
		psPLC.Last_frame_lost = 1
	} else {
		if psDec.SPLC.Last_frame_lost != 0 {
			silk_sum_sqr_shift(&energy, &energy_shift, &frame[0], length)
			if energy_shift > psPLC.Conc_energy_shift {
				psPLC.Conc_energy = psPLC.Conc_energy >> opus_int32(energy_shift-psPLC.Conc_energy_shift)
			} else if energy_shift < psPLC.Conc_energy_shift {
				energy = energy >> opus_int32(psPLC.Conc_energy_shift-energy_shift)
			}
			if energy > psPLC.Conc_energy {
				var (
					frac_Q24  opus_int32
					LZ        opus_int32
					gain_Q16  opus_int32
					slope_Q16 opus_int32
				)
				LZ = silk_CLZ32(psPLC.Conc_energy)
				LZ = LZ - 1
				psPLC.Conc_energy = opus_int32(opus_uint32(psPLC.Conc_energy) << opus_uint32(LZ))
				energy = energy >> silk_max_32(24-LZ, 0)
				frac_Q24 = psPLC.Conc_energy / (func() opus_int32 {
					if energy > 1 {
						return energy
					}
					return 1
				}())
				gain_Q16 = opus_int32(opus_uint32(silk_SQRT_APPROX(frac_Q24)) << 4)
				slope_Q16 = ((1 << 16) - gain_Q16) / opus_int32(length)
				slope_Q16 = opus_int32(opus_uint32(slope_Q16) << 2)
				for i = 0; i < length; i++ {
					frame[i] = opus_int16((gain_Q16 * opus_int32(int64(frame[i]))) >> 16)
					gain_Q16 += slope_Q16
					if gain_Q16 > 1<<16 {
						break
					}
				}
			}
		}
		psPLC.Last_frame_lost = 0
	}
}
