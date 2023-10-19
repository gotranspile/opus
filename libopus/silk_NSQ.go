package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

func silk_noise_shape_quantizer_short_prediction_c(buf32 *opus_int32, coef16 *opus_int16, order int64) opus_int32 {
	var out opus_int32
	out = opus_int32(order >> 1)
	out = out + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf32), unsafe.Sizeof(opus_int32(0))*0))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(coef16), unsafe.Sizeof(opus_int16(0))*0))))) >> 16)
	out = out + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf32), -int(unsafe.Sizeof(opus_int32(0))*1)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(coef16), unsafe.Sizeof(opus_int16(0))*1))))) >> 16)
	out = out + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf32), -int(unsafe.Sizeof(opus_int32(0))*2)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(coef16), unsafe.Sizeof(opus_int16(0))*2))))) >> 16)
	out = out + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf32), -int(unsafe.Sizeof(opus_int32(0))*3)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(coef16), unsafe.Sizeof(opus_int16(0))*3))))) >> 16)
	out = out + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf32), -int(unsafe.Sizeof(opus_int32(0))*4)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(coef16), unsafe.Sizeof(opus_int16(0))*4))))) >> 16)
	out = out + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf32), -int(unsafe.Sizeof(opus_int32(0))*5)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(coef16), unsafe.Sizeof(opus_int16(0))*5))))) >> 16)
	out = out + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf32), -int(unsafe.Sizeof(opus_int32(0))*6)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(coef16), unsafe.Sizeof(opus_int16(0))*6))))) >> 16)
	out = out + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf32), -int(unsafe.Sizeof(opus_int32(0))*7)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(coef16), unsafe.Sizeof(opus_int16(0))*7))))) >> 16)
	out = out + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf32), -int(unsafe.Sizeof(opus_int32(0))*8)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(coef16), unsafe.Sizeof(opus_int16(0))*8))))) >> 16)
	out = out + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf32), -int(unsafe.Sizeof(opus_int32(0))*9)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(coef16), unsafe.Sizeof(opus_int16(0))*9))))) >> 16)
	if order == 16 {
		out = out + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf32), -int(unsafe.Sizeof(opus_int32(0))*10)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(coef16), unsafe.Sizeof(opus_int16(0))*10))))) >> 16)
		out = out + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf32), -int(unsafe.Sizeof(opus_int32(0))*11)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(coef16), unsafe.Sizeof(opus_int16(0))*11))))) >> 16)
		out = out + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf32), -int(unsafe.Sizeof(opus_int32(0))*12)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(coef16), unsafe.Sizeof(opus_int16(0))*12))))) >> 16)
		out = out + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf32), -int(unsafe.Sizeof(opus_int32(0))*13)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(coef16), unsafe.Sizeof(opus_int16(0))*13))))) >> 16)
		out = out + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf32), -int(unsafe.Sizeof(opus_int32(0))*14)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(coef16), unsafe.Sizeof(opus_int16(0))*14))))) >> 16)
		out = out + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf32), -int(unsafe.Sizeof(opus_int32(0))*15)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(coef16), unsafe.Sizeof(opus_int16(0))*15))))) >> 16)
	}
	return out
}
func silk_NSQ_noise_shape_feedback_loop_c(data0 *opus_int32, data1 *opus_int32, coef *opus_int16, order int64) opus_int32 {
	var (
		out  opus_int32
		tmp1 opus_int32
		tmp2 opus_int32
		j    int64
	)
	tmp2 = *(*opus_int32)(unsafe.Add(unsafe.Pointer(data0), unsafe.Sizeof(opus_int32(0))*0))
	tmp1 = *(*opus_int32)(unsafe.Add(unsafe.Pointer(data1), unsafe.Sizeof(opus_int32(0))*0))
	*(*opus_int32)(unsafe.Add(unsafe.Pointer(data1), unsafe.Sizeof(opus_int32(0))*0)) = tmp2
	out = opus_int32(order >> 1)
	out = out + ((tmp2 * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(coef), unsafe.Sizeof(opus_int16(0))*0))))) >> 16)
	for j = 2; j < order; j += 2 {
		tmp2 = *(*opus_int32)(unsafe.Add(unsafe.Pointer(data1), unsafe.Sizeof(opus_int32(0))*uintptr(j-1)))
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(data1), unsafe.Sizeof(opus_int32(0))*uintptr(j-1))) = tmp1
		out = out + ((tmp1 * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(coef), unsafe.Sizeof(opus_int16(0))*uintptr(j-1)))))) >> 16)
		tmp1 = *(*opus_int32)(unsafe.Add(unsafe.Pointer(data1), unsafe.Sizeof(opus_int32(0))*uintptr(j+0)))
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(data1), unsafe.Sizeof(opus_int32(0))*uintptr(j+0))) = tmp2
		out = out + ((tmp2 * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(coef), unsafe.Sizeof(opus_int16(0))*uintptr(j)))))) >> 16)
	}
	*(*opus_int32)(unsafe.Add(unsafe.Pointer(data1), unsafe.Sizeof(opus_int32(0))*uintptr(order-1))) = tmp1
	out = out + ((tmp1 * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(coef), unsafe.Sizeof(opus_int16(0))*uintptr(order-1)))))) >> 16)
	out = opus_int32(opus_uint32(out) << 1)
	return out
}
func silk_NSQ_c(psEncC *silk_encoder_state, NSQ *silk_nsq_state, psIndices *SideInfoIndices, x16 [0]opus_int16, pulses [0]int8, PredCoef_Q12 [32]opus_int16, LTPCoef_Q14 [20]opus_int16, AR_Q13 [96]opus_int16, HarmShapeGain_Q14 [4]int64, Tilt_Q14 [4]int64, LF_shp_Q14 [4]opus_int32, Gains_Q16 [4]opus_int32, pitchL [4]int64, Lambda_Q10 int64, LTP_scale_Q14 int64) {
	var (
		k                      int64
		lag                    int64
		start_idx              int64
		LSF_interpolation_flag int64
		A_Q12                  *opus_int16
		B_Q14                  *opus_int16
		AR_shp_Q13             *opus_int16
		pxq                    *opus_int16
		sLTP_Q15               *opus_int32
		sLTP                   *opus_int16
		HarmShapeFIRPacked_Q14 opus_int32
		offset_Q10             int64
		x_sc_Q10               *opus_int32
	)
	NSQ.Rand_seed = opus_int32(psIndices.Seed)
	lag = NSQ.LagPrev
	offset_Q10 = int64(silk_Quantization_Offsets_Q10[int64(psIndices.SignalType)>>1][psIndices.QuantOffsetType])
	if int64(psIndices.NLSFInterpCoef_Q2) == 4 {
		LSF_interpolation_flag = 0
	} else {
		LSF_interpolation_flag = 1
	}
	sLTP_Q15 = (*opus_int32)(libc.Malloc(int((psEncC.Ltp_mem_length + psEncC.Frame_length) * int64(unsafe.Sizeof(opus_int32(0))))))
	sLTP = (*opus_int16)(libc.Malloc(int((psEncC.Ltp_mem_length + psEncC.Frame_length) * int64(unsafe.Sizeof(opus_int16(0))))))
	x_sc_Q10 = (*opus_int32)(libc.Malloc(int(psEncC.Subfr_length * int64(unsafe.Sizeof(opus_int32(0))))))
	NSQ.SLTP_shp_buf_idx = psEncC.Ltp_mem_length
	NSQ.SLTP_buf_idx = psEncC.Ltp_mem_length
	pxq = &NSQ.Xq[psEncC.Ltp_mem_length]
	for k = 0; k < psEncC.Nb_subfr; k++ {
		A_Q12 = &PredCoef_Q12[((k>>1)|(1-LSF_interpolation_flag))*MAX_LPC_ORDER]
		B_Q14 = &LTPCoef_Q14[k*LTP_ORDER]
		AR_shp_Q13 = &AR_Q13[k*MAX_SHAPE_LPC_ORDER]
		HarmShapeFIRPacked_Q14 = opus_int32((HarmShapeGain_Q14[k]) >> 2)
		HarmShapeFIRPacked_Q14 |= opus_int32(opus_uint32(opus_int32((HarmShapeGain_Q14[k])>>1)) << 16)
		NSQ.Rewhite_flag = 0
		if int64(psIndices.SignalType) == TYPE_VOICED {
			lag = pitchL[k]
			if (k & int64(3-(opus_int32(opus_uint32(LSF_interpolation_flag)<<1)))) == 0 {
				start_idx = psEncC.Ltp_mem_length - lag - psEncC.PredictLPCOrder - LTP_ORDER/2
				silk_LPC_analysis_filter((*opus_int16)(unsafe.Add(unsafe.Pointer(sLTP), unsafe.Sizeof(opus_int16(0))*uintptr(start_idx))), &NSQ.Xq[start_idx+k*psEncC.Subfr_length], A_Q12, opus_int32(psEncC.Ltp_mem_length-start_idx), opus_int32(psEncC.PredictLPCOrder), psEncC.Arch)
				NSQ.Rewhite_flag = 1
				NSQ.SLTP_buf_idx = psEncC.Ltp_mem_length
			}
		}
		silk_nsq_scale_states(psEncC, NSQ, x16, [0]opus_int32(x_sc_Q10), [0]opus_int16(sLTP), [0]opus_int32(sLTP_Q15), k, LTP_scale_Q14, Gains_Q16, pitchL, int64(psIndices.SignalType))
		silk_noise_shape_quantizer(NSQ, int64(psIndices.SignalType), [0]opus_int32(x_sc_Q10), pulses, [0]opus_int16(pxq), [0]opus_int32(sLTP_Q15), [0]opus_int16(A_Q12), [0]opus_int16(B_Q14), [0]opus_int16(AR_shp_Q13), lag, HarmShapeFIRPacked_Q14, Tilt_Q14[k], LF_shp_Q14[k], Gains_Q16[k], Lambda_Q10, offset_Q10, psEncC.Subfr_length, psEncC.ShapingLPCOrder, psEncC.PredictLPCOrder, psEncC.Arch)
		x16 += [0]opus_int16(psEncC.Subfr_length)
		pulses += [0]int8(psEncC.Subfr_length)
		pxq = (*opus_int16)(unsafe.Add(unsafe.Pointer(pxq), unsafe.Sizeof(opus_int16(0))*uintptr(psEncC.Subfr_length)))
	}
	NSQ.LagPrev = pitchL[psEncC.Nb_subfr-1]
	libc.MemMove(unsafe.Pointer(&NSQ.Xq[0]), unsafe.Pointer(&NSQ.Xq[psEncC.Frame_length]), int(psEncC.Ltp_mem_length*int64(unsafe.Sizeof(opus_int16(0)))))
	libc.MemMove(unsafe.Pointer(&NSQ.SLTP_shp_Q14[0]), unsafe.Pointer(&NSQ.SLTP_shp_Q14[psEncC.Frame_length]), int(psEncC.Ltp_mem_length*int64(unsafe.Sizeof(opus_int32(0)))))
}
func silk_noise_shape_quantizer(NSQ *silk_nsq_state, signalType int64, x_sc_Q10 [0]opus_int32, pulses [0]int8, xq [0]opus_int16, sLTP_Q15 [0]opus_int32, a_Q12 [0]opus_int16, b_Q14 [0]opus_int16, AR_shp_Q13 [0]opus_int16, lag int64, HarmShapeFIRPacked_Q14 opus_int32, Tilt_Q14 int64, LF_shp_Q14 opus_int32, Gain_Q16 opus_int32, Lambda_Q10 int64, offset_Q10 int64, length int64, shapingLPCOrder int64, predictLPCOrder int64, arch int64) {
	var (
		i              int64
		LTP_pred_Q13   opus_int32
		LPC_pred_Q10   opus_int32
		n_AR_Q12       opus_int32
		n_LTP_Q13      opus_int32
		n_LF_Q12       opus_int32
		r_Q10          opus_int32
		rr_Q10         opus_int32
		q1_Q0          opus_int32
		q1_Q10         opus_int32
		q2_Q10         opus_int32
		rd1_Q20        opus_int32
		rd2_Q20        opus_int32
		exc_Q14        opus_int32
		LPC_exc_Q14    opus_int32
		xq_Q14         opus_int32
		Gain_Q10       opus_int32
		tmp1           opus_int32
		tmp2           opus_int32
		sLF_AR_shp_Q14 opus_int32
		psLPC_Q14      *opus_int32
		shp_lag_ptr    *opus_int32
		pred_lag_ptr   *opus_int32
	)
	shp_lag_ptr = &NSQ.SLTP_shp_Q14[NSQ.SLTP_shp_buf_idx-lag+HARM_SHAPE_FIR_TAPS/2]
	pred_lag_ptr = &sLTP_Q15[NSQ.SLTP_buf_idx-lag+LTP_ORDER/2]
	Gain_Q10 = Gain_Q16 >> 6
	psLPC_Q14 = &NSQ.SLPC_Q14[MAX_LPC_ORDER-1]
	for i = 0; i < length; i++ {
		NSQ.Rand_seed = opus_int32(RAND_INCREMENT + opus_uint32(NSQ.Rand_seed)*RAND_MULTIPLIER)
		LPC_pred_Q10 = func() opus_int32 {
			_ = arch
			return silk_noise_shape_quantizer_short_prediction_c(psLPC_Q14, &a_Q12[0], predictLPCOrder)
		}()
		if signalType == TYPE_VOICED {
			LTP_pred_Q13 = 2
			LTP_pred_Q13 = LTP_pred_Q13 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), unsafe.Sizeof(opus_int32(0))*0))) * opus_int32(int64(b_Q14[0]))) >> 16)
			LTP_pred_Q13 = LTP_pred_Q13 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), -int(unsafe.Sizeof(opus_int32(0))*1)))) * opus_int32(int64(b_Q14[1]))) >> 16)
			LTP_pred_Q13 = LTP_pred_Q13 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), -int(unsafe.Sizeof(opus_int32(0))*2)))) * opus_int32(int64(b_Q14[2]))) >> 16)
			LTP_pred_Q13 = LTP_pred_Q13 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), -int(unsafe.Sizeof(opus_int32(0))*3)))) * opus_int32(int64(b_Q14[3]))) >> 16)
			LTP_pred_Q13 = LTP_pred_Q13 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), -int(unsafe.Sizeof(opus_int32(0))*4)))) * opus_int32(int64(b_Q14[4]))) >> 16)
			pred_lag_ptr = (*opus_int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), unsafe.Sizeof(opus_int32(0))*1))
		} else {
			LTP_pred_Q13 = 0
		}
		n_AR_Q12 = func() opus_int32 {
			_ = arch
			return silk_NSQ_noise_shape_feedback_loop_c(&NSQ.SDiff_shp_Q14, &NSQ.SAR2_Q14[0], &AR_shp_Q13[0], shapingLPCOrder)
		}()
		n_AR_Q12 = n_AR_Q12 + ((NSQ.SLF_AR_shp_Q14 * opus_int32(int64(opus_int16(Tilt_Q14)))) >> 16)
		n_LF_Q12 = ((NSQ.SLTP_shp_Q14[NSQ.SLTP_shp_buf_idx-1]) * opus_int32(int64(opus_int16(LF_shp_Q14)))) >> 16
		n_LF_Q12 = n_LF_Q12 + ((NSQ.SLF_AR_shp_Q14 * opus_int32(int64(LF_shp_Q14)>>16)) >> 16)
		tmp1 = (opus_int32(opus_uint32(LPC_pred_Q10) << 2)) - n_AR_Q12
		tmp1 = tmp1 - n_LF_Q12
		if lag > 0 {
			n_LTP_Q13 = ((func() opus_int32 {
				if ((opus_uint32(*(*opus_int32)(unsafe.Add(unsafe.Pointer(shp_lag_ptr), unsafe.Sizeof(opus_int32(0))*0))) + opus_uint32(*(*opus_int32)(unsafe.Add(unsafe.Pointer(shp_lag_ptr), -int(unsafe.Sizeof(opus_int32(0))*2))))) & 0x80000000) == 0 {
					if (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(shp_lag_ptr), unsafe.Sizeof(opus_int32(0))*0))) & (*(*opus_int32)(unsafe.Add(unsafe.Pointer(shp_lag_ptr), -int(unsafe.Sizeof(opus_int32(0))*2))))) & 0x80000000) != 0 {
						return 0x80000000
					}
					return (*(*opus_int32)(unsafe.Add(unsafe.Pointer(shp_lag_ptr), unsafe.Sizeof(opus_int32(0))*0))) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(shp_lag_ptr), -int(unsafe.Sizeof(opus_int32(0))*2))))
				}
				if (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(shp_lag_ptr), unsafe.Sizeof(opus_int32(0))*0))) | (*(*opus_int32)(unsafe.Add(unsafe.Pointer(shp_lag_ptr), -int(unsafe.Sizeof(opus_int32(0))*2))))) & 0x80000000) == 0 {
					return silk_int32_MAX
				}
				return (*(*opus_int32)(unsafe.Add(unsafe.Pointer(shp_lag_ptr), unsafe.Sizeof(opus_int32(0))*0))) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(shp_lag_ptr), -int(unsafe.Sizeof(opus_int32(0))*2))))
			}()) * opus_int32(int64(opus_int16(HarmShapeFIRPacked_Q14)))) >> 16
			n_LTP_Q13 = n_LTP_Q13 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(shp_lag_ptr), -int(unsafe.Sizeof(opus_int32(0))*1)))) * opus_int32(int64(HarmShapeFIRPacked_Q14)>>16)) >> 16)
			n_LTP_Q13 = opus_int32(opus_uint32(n_LTP_Q13) << 1)
			shp_lag_ptr = (*opus_int32)(unsafe.Add(unsafe.Pointer(shp_lag_ptr), unsafe.Sizeof(opus_int32(0))*1))
			tmp2 = LTP_pred_Q13 - n_LTP_Q13
			tmp1 = tmp2 + (opus_int32(opus_uint32(tmp1) << 1))
			if 3 == 1 {
				tmp1 = (tmp1 >> 1) + (tmp1 & 1)
			} else {
				tmp1 = ((tmp1 >> (3 - 1)) + 1) >> 1
			}
		} else {
			if 2 == 1 {
				tmp1 = (tmp1 >> 1) + (tmp1 & 1)
			} else {
				tmp1 = ((tmp1 >> (2 - 1)) + 1) >> 1
			}
		}
		r_Q10 = (x_sc_Q10[i]) - tmp1
		if NSQ.Rand_seed < 0 {
			r_Q10 = -r_Q10
		}
		if (-(31 << 10)) > (30 << 10) {
			if r_Q10 > (-(31 << 10)) {
				r_Q10 = -(31 << 10)
			} else if r_Q10 < (30 << 10) {
				r_Q10 = 30 << 10
			} else {
				r_Q10 = r_Q10
			}
		} else if r_Q10 > (30 << 10) {
			r_Q10 = 30 << 10
		} else if r_Q10 < (-(31 << 10)) {
			r_Q10 = -(31 << 10)
		} else {
			r_Q10 = r_Q10
		}
		q1_Q10 = r_Q10 - opus_int32(offset_Q10)
		q1_Q0 = q1_Q10 >> 10
		if Lambda_Q10 > 2048 {
			var rdo_offset int64 = Lambda_Q10/2 - 512
			if q1_Q10 > opus_int32(rdo_offset) {
				q1_Q0 = (q1_Q10 - opus_int32(rdo_offset)) >> 10
			} else if q1_Q10 < opus_int32(-rdo_offset) {
				q1_Q0 = (q1_Q10 + opus_int32(rdo_offset)) >> 10
			} else if q1_Q10 < 0 {
				q1_Q0 = -1
			} else {
				q1_Q0 = 0
			}
		}
		if q1_Q0 > 0 {
			q1_Q10 = (opus_int32(opus_uint32(q1_Q0) << 10)) - QUANT_LEVEL_ADJUST_Q10
			q1_Q10 = q1_Q10 + opus_int32(offset_Q10)
			q2_Q10 = q1_Q10 + 1024
			rd1_Q20 = opus_int32(opus_int16(q1_Q10)) * opus_int32(opus_int16(Lambda_Q10))
			rd2_Q20 = opus_int32(opus_int16(q2_Q10)) * opus_int32(opus_int16(Lambda_Q10))
		} else if q1_Q0 == 0 {
			q1_Q10 = opus_int32(offset_Q10)
			q2_Q10 = q1_Q10 + opus_int32(1024-QUANT_LEVEL_ADJUST_Q10)
			rd1_Q20 = opus_int32(opus_int16(q1_Q10)) * opus_int32(opus_int16(Lambda_Q10))
			rd2_Q20 = opus_int32(opus_int16(q2_Q10)) * opus_int32(opus_int16(Lambda_Q10))
		} else if q1_Q0 == -1 {
			q2_Q10 = opus_int32(offset_Q10)
			q1_Q10 = q2_Q10 - opus_int32(1024-QUANT_LEVEL_ADJUST_Q10)
			rd1_Q20 = opus_int32(opus_int16(-q1_Q10)) * opus_int32(opus_int16(Lambda_Q10))
			rd2_Q20 = opus_int32(opus_int16(q2_Q10)) * opus_int32(opus_int16(Lambda_Q10))
		} else {
			q1_Q10 = (opus_int32(opus_uint32(q1_Q0) << 10)) + QUANT_LEVEL_ADJUST_Q10
			q1_Q10 = q1_Q10 + opus_int32(offset_Q10)
			q2_Q10 = q1_Q10 + 1024
			rd1_Q20 = opus_int32(opus_int16(-q1_Q10)) * opus_int32(opus_int16(Lambda_Q10))
			rd2_Q20 = opus_int32(opus_int16(-q2_Q10)) * opus_int32(opus_int16(Lambda_Q10))
		}
		rr_Q10 = r_Q10 - q1_Q10
		rd1_Q20 = rd1_Q20 + (opus_int32(opus_int16(rr_Q10)))*opus_int32(opus_int16(rr_Q10))
		rr_Q10 = r_Q10 - q2_Q10
		rd2_Q20 = rd2_Q20 + (opus_int32(opus_int16(rr_Q10)))*opus_int32(opus_int16(rr_Q10))
		if rd2_Q20 < rd1_Q20 {
			q1_Q10 = q2_Q10
		}
		if 10 == 1 {
			pulses[i] = int8((q1_Q10 >> 1) + (q1_Q10 & 1))
		} else {
			pulses[i] = int8(((q1_Q10 >> (10 - 1)) + 1) >> 1)
		}
		exc_Q14 = opus_int32(opus_uint32(q1_Q10) << 4)
		if NSQ.Rand_seed < 0 {
			exc_Q14 = -exc_Q14
		}
		LPC_exc_Q14 = exc_Q14 + (opus_int32(opus_uint32(LTP_pred_Q13) << 1))
		xq_Q14 = LPC_exc_Q14 + (opus_int32(opus_uint32(LPC_pred_Q10) << 4))
		if (func() opus_int32 {
			if 8 == 1 {
				return ((opus_int32((int64(xq_Q14) * int64(Gain_Q10)) >> 16)) >> 1) + ((opus_int32((int64(xq_Q14) * int64(Gain_Q10)) >> 16)) & 1)
			}
			return (((opus_int32((int64(xq_Q14) * int64(Gain_Q10)) >> 16)) >> (8 - 1)) + 1) >> 1
		}()) > silk_int16_MAX {
			xq[i] = silk_int16_MAX
		} else if (func() opus_int32 {
			if 8 == 1 {
				return ((opus_int32((int64(xq_Q14) * int64(Gain_Q10)) >> 16)) >> 1) + ((opus_int32((int64(xq_Q14) * int64(Gain_Q10)) >> 16)) & 1)
			}
			return (((opus_int32((int64(xq_Q14) * int64(Gain_Q10)) >> 16)) >> (8 - 1)) + 1) >> 1
		}()) < opus_int32(math.MinInt16) {
			xq[i] = math.MinInt16
		} else if 8 == 1 {
			xq[i] = opus_int16(((opus_int32((int64(xq_Q14) * int64(Gain_Q10)) >> 16)) >> 1) + ((opus_int32((int64(xq_Q14) * int64(Gain_Q10)) >> 16)) & 1))
		} else {
			xq[i] = opus_int16((((opus_int32((int64(xq_Q14) * int64(Gain_Q10)) >> 16)) >> (8 - 1)) + 1) >> 1)
		}
		psLPC_Q14 = (*opus_int32)(unsafe.Add(unsafe.Pointer(psLPC_Q14), unsafe.Sizeof(opus_int32(0))*1))
		*psLPC_Q14 = xq_Q14
		NSQ.SDiff_shp_Q14 = xq_Q14 - (opus_int32(opus_uint32(x_sc_Q10[i]) << 4))
		sLF_AR_shp_Q14 = NSQ.SDiff_shp_Q14 - (opus_int32(opus_uint32(n_AR_Q12) << 2))
		NSQ.SLF_AR_shp_Q14 = sLF_AR_shp_Q14
		NSQ.SLTP_shp_Q14[NSQ.SLTP_shp_buf_idx] = sLF_AR_shp_Q14 - (opus_int32(opus_uint32(n_LF_Q12) << 2))
		sLTP_Q15[NSQ.SLTP_buf_idx] = opus_int32(opus_uint32(LPC_exc_Q14) << 1)
		NSQ.SLTP_shp_buf_idx++
		NSQ.SLTP_buf_idx++
		NSQ.Rand_seed = opus_int32(opus_uint32(NSQ.Rand_seed) + opus_uint32(pulses[i]))
	}
	libc.MemCpy(unsafe.Pointer(&NSQ.SLPC_Q14[0]), unsafe.Pointer(&NSQ.SLPC_Q14[length]), int(MAX_LPC_ORDER*unsafe.Sizeof(opus_int32(0))))
}
func silk_nsq_scale_states(psEncC *silk_encoder_state, NSQ *silk_nsq_state, x16 [0]opus_int16, x_sc_Q10 [0]opus_int32, sLTP [0]opus_int16, sLTP_Q15 [0]opus_int32, subfr int64, LTP_scale_Q14 int64, Gains_Q16 [4]opus_int32, pitchL [4]int64, signal_type int64) {
	var (
		i            int64
		lag          int64
		gain_adj_Q16 opus_int32
		inv_gain_Q31 opus_int32
		inv_gain_Q26 opus_int32
	)
	lag = pitchL[subfr]
	inv_gain_Q31 = silk_INVERSE32_varQ(func() opus_int32 {
		if (Gains_Q16[subfr]) > 1 {
			return Gains_Q16[subfr]
		}
		return 1
	}(), 47)
	if 5 == 1 {
		inv_gain_Q26 = (inv_gain_Q31 >> 1) + (inv_gain_Q31 & 1)
	} else {
		inv_gain_Q26 = ((inv_gain_Q31 >> (5 - 1)) + 1) >> 1
	}
	for i = 0; i < psEncC.Subfr_length; i++ {
		x_sc_Q10[i] = opus_int32((int64(x16[i]) * int64(inv_gain_Q26)) >> 16)
	}
	if NSQ.Rewhite_flag != 0 {
		if subfr == 0 {
			inv_gain_Q31 = opus_int32(opus_uint32((inv_gain_Q31*opus_int32(int64(opus_int16(LTP_scale_Q14))))>>16) << 2)
		}
		for i = NSQ.SLTP_buf_idx - lag - LTP_ORDER/2; i < NSQ.SLTP_buf_idx; i++ {
			sLTP_Q15[i] = (inv_gain_Q31 * opus_int32(int64(sLTP[i]))) >> 16
		}
	}
	if Gains_Q16[subfr] != NSQ.Prev_gain_Q16 {
		gain_adj_Q16 = silk_DIV32_varQ(NSQ.Prev_gain_Q16, Gains_Q16[subfr], 16)
		for i = NSQ.SLTP_shp_buf_idx - psEncC.Ltp_mem_length; i < NSQ.SLTP_shp_buf_idx; i++ {
			NSQ.SLTP_shp_Q14[i] = opus_int32((int64(gain_adj_Q16) * int64(NSQ.SLTP_shp_Q14[i])) >> 16)
		}
		if signal_type == TYPE_VOICED && NSQ.Rewhite_flag == 0 {
			for i = NSQ.SLTP_buf_idx - lag - LTP_ORDER/2; i < NSQ.SLTP_buf_idx; i++ {
				sLTP_Q15[i] = opus_int32((int64(gain_adj_Q16) * int64(sLTP_Q15[i])) >> 16)
			}
		}
		NSQ.SLF_AR_shp_Q14 = opus_int32((int64(gain_adj_Q16) * int64(NSQ.SLF_AR_shp_Q14)) >> 16)
		NSQ.SDiff_shp_Q14 = opus_int32((int64(gain_adj_Q16) * int64(NSQ.SDiff_shp_Q14)) >> 16)
		for i = 0; i < MAX_LPC_ORDER; i++ {
			NSQ.SLPC_Q14[i] = opus_int32((int64(gain_adj_Q16) * int64(NSQ.SLPC_Q14[i])) >> 16)
		}
		for i = 0; i < MAX_SHAPE_LPC_ORDER; i++ {
			NSQ.SAR2_Q14[i] = opus_int32((int64(gain_adj_Q16) * int64(NSQ.SAR2_Q14[i])) >> 16)
		}
		NSQ.Prev_gain_Q16 = Gains_Q16[subfr]
	}
}
