package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

type NSQ_del_dec_struct struct {
	SLPC_Q14  [96]opus_int32
	RandState [40]opus_int32
	Q_Q10     [40]opus_int32
	Xq_Q14    [40]opus_int32
	Pred_Q15  [40]opus_int32
	Shape_Q14 [40]opus_int32
	SAR2_Q14  [24]opus_int32
	LF_AR_Q14 opus_int32
	Diff_Q14  opus_int32
	Seed      opus_int32
	SeedInit  opus_int32
	RD_Q10    opus_int32
}
type NSQ_sample_struct struct {
	Q_Q10        opus_int32
	RD_Q10       opus_int32
	Xq_Q14       opus_int32
	LF_AR_Q14    opus_int32
	Diff_Q14     opus_int32
	SLTP_shp_Q14 opus_int32
	LPC_exc_Q14  opus_int32
}
type NSQ_sample_pair [2]NSQ_sample_struct

func silk_NSQ_del_dec_c(psEncC *silk_encoder_state, NSQ *silk_nsq_state, psIndices *SideInfoIndices, x16 [0]opus_int16, pulses [0]int8, PredCoef_Q12 [32]opus_int16, LTPCoef_Q14 [20]opus_int16, AR_Q13 [96]opus_int16, HarmShapeGain_Q14 [4]int64, Tilt_Q14 [4]int64, LF_shp_Q14 [4]opus_int32, Gains_Q16 [4]opus_int32, pitchL [4]int64, Lambda_Q10 int64, LTP_scale_Q14 int64) {
	var (
		i                      int64
		k                      int64
		lag                    int64
		start_idx              int64
		LSF_interpolation_flag int64
		Winner_ind             int64
		subfr                  int64
		last_smple_idx         int64
		smpl_buf_idx           int64
		decisionDelay          int64
		A_Q12                  *opus_int16
		B_Q14                  *opus_int16
		AR_shp_Q13             *opus_int16
		pxq                    *opus_int16
		sLTP_Q15               *opus_int32
		sLTP                   *opus_int16
		HarmShapeFIRPacked_Q14 opus_int32
		offset_Q10             int64
		RDmin_Q10              opus_int32
		Gain_Q10               opus_int32
		x_sc_Q10               *opus_int32
		delayedGain_Q10        *opus_int32
		psDelDec               *NSQ_del_dec_struct
		psDD                   *NSQ_del_dec_struct
	)
	lag = NSQ.LagPrev
	psDelDec = (*NSQ_del_dec_struct)(libc.Malloc(int(psEncC.NStatesDelayedDecision * int64(unsafe.Sizeof(NSQ_del_dec_struct{})))))
	libc.MemSet(unsafe.Pointer(psDelDec), 0, int(psEncC.NStatesDelayedDecision*int64(unsafe.Sizeof(NSQ_del_dec_struct{}))))
	for k = 0; k < psEncC.NStatesDelayedDecision; k++ {
		psDD = (*NSQ_del_dec_struct)(unsafe.Add(unsafe.Pointer(psDelDec), unsafe.Sizeof(NSQ_del_dec_struct{})*uintptr(k)))
		psDD.Seed = opus_int32((k + int64(psIndices.Seed)) & 3)
		psDD.SeedInit = psDD.Seed
		psDD.RD_Q10 = 0
		psDD.LF_AR_Q14 = NSQ.SLF_AR_shp_Q14
		psDD.Diff_Q14 = NSQ.SDiff_shp_Q14
		psDD.Shape_Q14[0] = NSQ.SLTP_shp_Q14[psEncC.Ltp_mem_length-1]
		libc.MemCpy(unsafe.Pointer(&psDD.SLPC_Q14[0]), unsafe.Pointer(&NSQ.SLPC_Q14[0]), int(MAX_LPC_ORDER*unsafe.Sizeof(opus_int32(0))))
		libc.MemCpy(unsafe.Pointer(&psDD.SAR2_Q14[0]), unsafe.Pointer(&NSQ.SAR2_Q14[0]), int(unsafe.Sizeof([24]opus_int32{})))
	}
	offset_Q10 = int64(silk_Quantization_Offsets_Q10[int64(psIndices.SignalType)>>1][psIndices.QuantOffsetType])
	smpl_buf_idx = 0
	decisionDelay = silk_min_int(DECISION_DELAY, psEncC.Subfr_length)
	if int64(psIndices.SignalType) == TYPE_VOICED {
		for k = 0; k < psEncC.Nb_subfr; k++ {
			decisionDelay = silk_min_int(decisionDelay, pitchL[k]-LTP_ORDER/2-1)
		}
	} else {
		if lag > 0 {
			decisionDelay = silk_min_int(decisionDelay, lag-LTP_ORDER/2-1)
		}
	}
	if int64(psIndices.NLSFInterpCoef_Q2) == 4 {
		LSF_interpolation_flag = 0
	} else {
		LSF_interpolation_flag = 1
	}
	sLTP_Q15 = (*opus_int32)(libc.Malloc(int((psEncC.Ltp_mem_length + psEncC.Frame_length) * int64(unsafe.Sizeof(opus_int32(0))))))
	sLTP = (*opus_int16)(libc.Malloc(int((psEncC.Ltp_mem_length + psEncC.Frame_length) * int64(unsafe.Sizeof(opus_int16(0))))))
	x_sc_Q10 = (*opus_int32)(libc.Malloc(int(psEncC.Subfr_length * int64(unsafe.Sizeof(opus_int32(0))))))
	delayedGain_Q10 = (*opus_int32)(libc.Malloc(int(DECISION_DELAY * unsafe.Sizeof(opus_int32(0)))))
	pxq = &NSQ.Xq[psEncC.Ltp_mem_length]
	NSQ.SLTP_shp_buf_idx = psEncC.Ltp_mem_length
	NSQ.SLTP_buf_idx = psEncC.Ltp_mem_length
	subfr = 0
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
				if k == 2 {
					RDmin_Q10 = (*(*NSQ_del_dec_struct)(unsafe.Add(unsafe.Pointer(psDelDec), unsafe.Sizeof(NSQ_del_dec_struct{})*0))).RD_Q10
					Winner_ind = 0
					for i = 1; i < psEncC.NStatesDelayedDecision; i++ {
						if (*(*NSQ_del_dec_struct)(unsafe.Add(unsafe.Pointer(psDelDec), unsafe.Sizeof(NSQ_del_dec_struct{})*uintptr(i)))).RD_Q10 < RDmin_Q10 {
							RDmin_Q10 = (*(*NSQ_del_dec_struct)(unsafe.Add(unsafe.Pointer(psDelDec), unsafe.Sizeof(NSQ_del_dec_struct{})*uintptr(i)))).RD_Q10
							Winner_ind = i
						}
					}
					for i = 0; i < psEncC.NStatesDelayedDecision; i++ {
						if i != Winner_ind {
							(*(*NSQ_del_dec_struct)(unsafe.Add(unsafe.Pointer(psDelDec), unsafe.Sizeof(NSQ_del_dec_struct{})*uintptr(i)))).RD_Q10 += opus_int32(silk_int32_MAX >> 4)
						}
					}
					psDD = (*NSQ_del_dec_struct)(unsafe.Add(unsafe.Pointer(psDelDec), unsafe.Sizeof(NSQ_del_dec_struct{})*uintptr(Winner_ind)))
					last_smple_idx = smpl_buf_idx + decisionDelay
					for i = 0; i < decisionDelay; i++ {
						last_smple_idx = (last_smple_idx - 1) % DECISION_DELAY
						if last_smple_idx < 0 {
							last_smple_idx += DECISION_DELAY
						}
						if 10 == 1 {
							pulses[i-decisionDelay] = int8(((psDD.Q_Q10[last_smple_idx]) >> 1) + ((psDD.Q_Q10[last_smple_idx]) & 1))
						} else {
							pulses[i-decisionDelay] = int8((((psDD.Q_Q10[last_smple_idx]) >> (10 - 1)) + 1) >> 1)
						}
						if (func() opus_int32 {
							if 14 == 1 {
								return ((opus_int32((int64(psDD.Xq_Q14[last_smple_idx]) * int64(Gains_Q16[1])) >> 16)) >> 1) + ((opus_int32((int64(psDD.Xq_Q14[last_smple_idx]) * int64(Gains_Q16[1])) >> 16)) & 1)
							}
							return (((opus_int32((int64(psDD.Xq_Q14[last_smple_idx]) * int64(Gains_Q16[1])) >> 16)) >> (14 - 1)) + 1) >> 1
						}()) > silk_int16_MAX {
							*(*opus_int16)(unsafe.Add(unsafe.Pointer(pxq), unsafe.Sizeof(opus_int16(0))*uintptr(i-decisionDelay))) = silk_int16_MAX
						} else if (func() opus_int32 {
							if 14 == 1 {
								return ((opus_int32((int64(psDD.Xq_Q14[last_smple_idx]) * int64(Gains_Q16[1])) >> 16)) >> 1) + ((opus_int32((int64(psDD.Xq_Q14[last_smple_idx]) * int64(Gains_Q16[1])) >> 16)) & 1)
							}
							return (((opus_int32((int64(psDD.Xq_Q14[last_smple_idx]) * int64(Gains_Q16[1])) >> 16)) >> (14 - 1)) + 1) >> 1
						}()) < opus_int32(math.MinInt16) {
							*(*opus_int16)(unsafe.Add(unsafe.Pointer(pxq), unsafe.Sizeof(opus_int16(0))*uintptr(i-decisionDelay))) = math.MinInt16
						} else if 14 == 1 {
							*(*opus_int16)(unsafe.Add(unsafe.Pointer(pxq), unsafe.Sizeof(opus_int16(0))*uintptr(i-decisionDelay))) = opus_int16(((opus_int32((int64(psDD.Xq_Q14[last_smple_idx]) * int64(Gains_Q16[1])) >> 16)) >> 1) + ((opus_int32((int64(psDD.Xq_Q14[last_smple_idx]) * int64(Gains_Q16[1])) >> 16)) & 1))
						} else {
							*(*opus_int16)(unsafe.Add(unsafe.Pointer(pxq), unsafe.Sizeof(opus_int16(0))*uintptr(i-decisionDelay))) = opus_int16((((opus_int32((int64(psDD.Xq_Q14[last_smple_idx]) * int64(Gains_Q16[1])) >> 16)) >> (14 - 1)) + 1) >> 1)
						}
						NSQ.SLTP_shp_Q14[NSQ.SLTP_shp_buf_idx-decisionDelay+i] = psDD.Shape_Q14[last_smple_idx]
					}
					subfr = 0
				}
				start_idx = psEncC.Ltp_mem_length - lag - psEncC.PredictLPCOrder - LTP_ORDER/2
				silk_LPC_analysis_filter((*opus_int16)(unsafe.Add(unsafe.Pointer(sLTP), unsafe.Sizeof(opus_int16(0))*uintptr(start_idx))), &NSQ.Xq[start_idx+k*psEncC.Subfr_length], A_Q12, opus_int32(psEncC.Ltp_mem_length-start_idx), opus_int32(psEncC.PredictLPCOrder), psEncC.Arch)
				NSQ.SLTP_buf_idx = psEncC.Ltp_mem_length
				NSQ.Rewhite_flag = 1
			}
		}
		silk_nsq_del_dec_scale_states(psEncC, NSQ, [0]NSQ_del_dec_struct(psDelDec), x16, [0]opus_int32(x_sc_Q10), [0]opus_int16(sLTP), [0]opus_int32(sLTP_Q15), k, psEncC.NStatesDelayedDecision, LTP_scale_Q14, Gains_Q16, pitchL, int64(psIndices.SignalType), decisionDelay)
		silk_noise_shape_quantizer_del_dec(NSQ, [0]NSQ_del_dec_struct(psDelDec), int64(psIndices.SignalType), [0]opus_int32(x_sc_Q10), pulses, [0]opus_int16(pxq), [0]opus_int32(sLTP_Q15), [0]opus_int32(delayedGain_Q10), [0]opus_int16(A_Q12), [0]opus_int16(B_Q14), [0]opus_int16(AR_shp_Q13), lag, HarmShapeFIRPacked_Q14, Tilt_Q14[k], LF_shp_Q14[k], Gains_Q16[k], Lambda_Q10, offset_Q10, psEncC.Subfr_length, func() int64 {
			p := &subfr
			x := *p
			*p++
			return x
		}(), psEncC.ShapingLPCOrder, psEncC.PredictLPCOrder, psEncC.Warping_Q16, psEncC.NStatesDelayedDecision, &smpl_buf_idx, decisionDelay, psEncC.Arch)
		x16 += [0]opus_int16(psEncC.Subfr_length)
		pulses += [0]int8(psEncC.Subfr_length)
		pxq = (*opus_int16)(unsafe.Add(unsafe.Pointer(pxq), unsafe.Sizeof(opus_int16(0))*uintptr(psEncC.Subfr_length)))
	}
	RDmin_Q10 = (*(*NSQ_del_dec_struct)(unsafe.Add(unsafe.Pointer(psDelDec), unsafe.Sizeof(NSQ_del_dec_struct{})*0))).RD_Q10
	Winner_ind = 0
	for k = 1; k < psEncC.NStatesDelayedDecision; k++ {
		if (*(*NSQ_del_dec_struct)(unsafe.Add(unsafe.Pointer(psDelDec), unsafe.Sizeof(NSQ_del_dec_struct{})*uintptr(k)))).RD_Q10 < RDmin_Q10 {
			RDmin_Q10 = (*(*NSQ_del_dec_struct)(unsafe.Add(unsafe.Pointer(psDelDec), unsafe.Sizeof(NSQ_del_dec_struct{})*uintptr(k)))).RD_Q10
			Winner_ind = k
		}
	}
	psDD = (*NSQ_del_dec_struct)(unsafe.Add(unsafe.Pointer(psDelDec), unsafe.Sizeof(NSQ_del_dec_struct{})*uintptr(Winner_ind)))
	psIndices.Seed = int8(psDD.SeedInit)
	last_smple_idx = smpl_buf_idx + decisionDelay
	Gain_Q10 = (Gains_Q16[psEncC.Nb_subfr-1]) >> 6
	for i = 0; i < decisionDelay; i++ {
		last_smple_idx = (last_smple_idx - 1) % DECISION_DELAY
		if last_smple_idx < 0 {
			last_smple_idx += DECISION_DELAY
		}
		if 10 == 1 {
			pulses[i-decisionDelay] = int8(((psDD.Q_Q10[last_smple_idx]) >> 1) + ((psDD.Q_Q10[last_smple_idx]) & 1))
		} else {
			pulses[i-decisionDelay] = int8((((psDD.Q_Q10[last_smple_idx]) >> (10 - 1)) + 1) >> 1)
		}
		if (func() opus_int32 {
			if 8 == 1 {
				return ((opus_int32((int64(psDD.Xq_Q14[last_smple_idx]) * int64(Gain_Q10)) >> 16)) >> 1) + ((opus_int32((int64(psDD.Xq_Q14[last_smple_idx]) * int64(Gain_Q10)) >> 16)) & 1)
			}
			return (((opus_int32((int64(psDD.Xq_Q14[last_smple_idx]) * int64(Gain_Q10)) >> 16)) >> (8 - 1)) + 1) >> 1
		}()) > silk_int16_MAX {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(pxq), unsafe.Sizeof(opus_int16(0))*uintptr(i-decisionDelay))) = silk_int16_MAX
		} else if (func() opus_int32 {
			if 8 == 1 {
				return ((opus_int32((int64(psDD.Xq_Q14[last_smple_idx]) * int64(Gain_Q10)) >> 16)) >> 1) + ((opus_int32((int64(psDD.Xq_Q14[last_smple_idx]) * int64(Gain_Q10)) >> 16)) & 1)
			}
			return (((opus_int32((int64(psDD.Xq_Q14[last_smple_idx]) * int64(Gain_Q10)) >> 16)) >> (8 - 1)) + 1) >> 1
		}()) < opus_int32(math.MinInt16) {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(pxq), unsafe.Sizeof(opus_int16(0))*uintptr(i-decisionDelay))) = math.MinInt16
		} else if 8 == 1 {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(pxq), unsafe.Sizeof(opus_int16(0))*uintptr(i-decisionDelay))) = opus_int16(((opus_int32((int64(psDD.Xq_Q14[last_smple_idx]) * int64(Gain_Q10)) >> 16)) >> 1) + ((opus_int32((int64(psDD.Xq_Q14[last_smple_idx]) * int64(Gain_Q10)) >> 16)) & 1))
		} else {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(pxq), unsafe.Sizeof(opus_int16(0))*uintptr(i-decisionDelay))) = opus_int16((((opus_int32((int64(psDD.Xq_Q14[last_smple_idx]) * int64(Gain_Q10)) >> 16)) >> (8 - 1)) + 1) >> 1)
		}
		NSQ.SLTP_shp_Q14[NSQ.SLTP_shp_buf_idx-decisionDelay+i] = psDD.Shape_Q14[last_smple_idx]
	}
	libc.MemCpy(unsafe.Pointer(&NSQ.SLPC_Q14[0]), unsafe.Pointer(&psDD.SLPC_Q14[psEncC.Subfr_length]), int(MAX_LPC_ORDER*unsafe.Sizeof(opus_int32(0))))
	libc.MemCpy(unsafe.Pointer(&NSQ.SAR2_Q14[0]), unsafe.Pointer(&psDD.SAR2_Q14[0]), int(unsafe.Sizeof([24]opus_int32{})))
	NSQ.SLF_AR_shp_Q14 = psDD.LF_AR_Q14
	NSQ.SDiff_shp_Q14 = psDD.Diff_Q14
	NSQ.LagPrev = pitchL[psEncC.Nb_subfr-1]
	libc.MemMove(unsafe.Pointer(&NSQ.Xq[0]), unsafe.Pointer(&NSQ.Xq[psEncC.Frame_length]), int(psEncC.Ltp_mem_length*int64(unsafe.Sizeof(opus_int16(0)))))
	libc.MemMove(unsafe.Pointer(&NSQ.SLTP_shp_Q14[0]), unsafe.Pointer(&NSQ.SLTP_shp_Q14[psEncC.Frame_length]), int(psEncC.Ltp_mem_length*int64(unsafe.Sizeof(opus_int32(0)))))
}
func silk_noise_shape_quantizer_del_dec(NSQ *silk_nsq_state, psDelDec [0]NSQ_del_dec_struct, signalType int64, x_Q10 [0]opus_int32, pulses [0]int8, xq [0]opus_int16, sLTP_Q15 [0]opus_int32, delayedGain_Q10 [0]opus_int32, a_Q12 [0]opus_int16, b_Q14 [0]opus_int16, AR_shp_Q13 [0]opus_int16, lag int64, HarmShapeFIRPacked_Q14 opus_int32, Tilt_Q14 int64, LF_shp_Q14 opus_int32, Gain_Q16 opus_int32, Lambda_Q10 int64, offset_Q10 int64, length int64, subfr int64, shapingLPCOrder int64, predictLPCOrder int64, warping_Q16 int64, nStatesDelayedDecision int64, smpl_buf_idx *int64, decisionDelay int64, arch int64) {
	var (
		i                 int64
		j                 int64
		k                 int64
		Winner_ind        int64
		RDmin_ind         int64
		RDmax_ind         int64
		last_smple_idx    int64
		Winner_rand_state opus_int32
		LTP_pred_Q14      opus_int32
		LPC_pred_Q14      opus_int32
		n_AR_Q14          opus_int32
		n_LTP_Q14         opus_int32
		n_LF_Q14          opus_int32
		r_Q10             opus_int32
		rr_Q10            opus_int32
		rd1_Q10           opus_int32
		rd2_Q10           opus_int32
		RDmin_Q10         opus_int32
		RDmax_Q10         opus_int32
		q1_Q0             opus_int32
		q1_Q10            opus_int32
		q2_Q10            opus_int32
		exc_Q14           opus_int32
		LPC_exc_Q14       opus_int32
		xq_Q14            opus_int32
		Gain_Q10          opus_int32
		tmp1              opus_int32
		tmp2              opus_int32
		sLF_AR_shp_Q14    opus_int32
		pred_lag_ptr      *opus_int32
		shp_lag_ptr       *opus_int32
		psLPC_Q14         *opus_int32
		psSampleState     *NSQ_sample_pair
		psDD              *NSQ_del_dec_struct
		psSS              *NSQ_sample_struct
	)
	psSampleState = (*NSQ_sample_pair)(libc.Malloc(int(nStatesDelayedDecision * int64(unsafe.Sizeof(NSQ_sample_pair{})))))
	shp_lag_ptr = &NSQ.SLTP_shp_Q14[NSQ.SLTP_shp_buf_idx-lag+HARM_SHAPE_FIR_TAPS/2]
	pred_lag_ptr = &sLTP_Q15[NSQ.SLTP_buf_idx-lag+LTP_ORDER/2]
	Gain_Q10 = Gain_Q16 >> 6
	for i = 0; i < length; i++ {
		if signalType == TYPE_VOICED {
			LTP_pred_Q14 = 2
			LTP_pred_Q14 = LTP_pred_Q14 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), unsafe.Sizeof(opus_int32(0))*0))) * opus_int32(int64(b_Q14[0]))) >> 16)
			LTP_pred_Q14 = LTP_pred_Q14 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), -int(unsafe.Sizeof(opus_int32(0))*1)))) * opus_int32(int64(b_Q14[1]))) >> 16)
			LTP_pred_Q14 = LTP_pred_Q14 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), -int(unsafe.Sizeof(opus_int32(0))*2)))) * opus_int32(int64(b_Q14[2]))) >> 16)
			LTP_pred_Q14 = LTP_pred_Q14 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), -int(unsafe.Sizeof(opus_int32(0))*3)))) * opus_int32(int64(b_Q14[3]))) >> 16)
			LTP_pred_Q14 = LTP_pred_Q14 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), -int(unsafe.Sizeof(opus_int32(0))*4)))) * opus_int32(int64(b_Q14[4]))) >> 16)
			LTP_pred_Q14 = opus_int32(opus_uint32(LTP_pred_Q14) << 1)
			pred_lag_ptr = (*opus_int32)(unsafe.Add(unsafe.Pointer(pred_lag_ptr), unsafe.Sizeof(opus_int32(0))*1))
		} else {
			LTP_pred_Q14 = 0
		}
		if lag > 0 {
			n_LTP_Q14 = ((func() opus_int32 {
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
			n_LTP_Q14 = n_LTP_Q14 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(shp_lag_ptr), -int(unsafe.Sizeof(opus_int32(0))*1)))) * opus_int32(int64(HarmShapeFIRPacked_Q14)>>16)) >> 16)
			n_LTP_Q14 = LTP_pred_Q14 - (opus_int32(opus_uint32(n_LTP_Q14) << 2))
			shp_lag_ptr = (*opus_int32)(unsafe.Add(unsafe.Pointer(shp_lag_ptr), unsafe.Sizeof(opus_int32(0))*1))
		} else {
			n_LTP_Q14 = 0
		}
		for k = 0; k < nStatesDelayedDecision; k++ {
			psDD = &psDelDec[k]
			psSS = (*NSQ_sample_struct)(unsafe.Pointer((*NSQ_sample_pair)(unsafe.Add(unsafe.Pointer(psSampleState), unsafe.Sizeof(NSQ_sample_pair{})*uintptr(k)))))
			psDD.Seed = opus_int32(RAND_INCREMENT + opus_uint32(psDD.Seed)*RAND_MULTIPLIER)
			psLPC_Q14 = &psDD.SLPC_Q14[MAX_LPC_ORDER-1+i]
			LPC_pred_Q14 = func() opus_int32 {
				_ = arch
				return silk_noise_shape_quantizer_short_prediction_c(psLPC_Q14, &a_Q12[0], predictLPCOrder)
			}()
			LPC_pred_Q14 = opus_int32(opus_uint32(LPC_pred_Q14) << 4)
			tmp2 = psDD.Diff_Q14 + (((psDD.SAR2_Q14[0]) * opus_int32(int64(opus_int16(warping_Q16)))) >> 16)
			tmp1 = (psDD.SAR2_Q14[0]) + (((psDD.SAR2_Q14[1] - tmp2) * opus_int32(int64(opus_int16(warping_Q16)))) >> 16)
			psDD.SAR2_Q14[0] = tmp2
			n_AR_Q14 = opus_int32(shapingLPCOrder >> 1)
			n_AR_Q14 = n_AR_Q14 + ((tmp2 * opus_int32(int64(AR_shp_Q13[0]))) >> 16)
			for j = 2; j < shapingLPCOrder; j += 2 {
				tmp2 = (psDD.SAR2_Q14[j-1]) + (((psDD.SAR2_Q14[j+0] - tmp1) * opus_int32(int64(opus_int16(warping_Q16)))) >> 16)
				psDD.SAR2_Q14[j-1] = tmp1
				n_AR_Q14 = n_AR_Q14 + ((tmp1 * opus_int32(int64(AR_shp_Q13[j-1]))) >> 16)
				tmp1 = (psDD.SAR2_Q14[j+0]) + (((psDD.SAR2_Q14[j+1] - tmp2) * opus_int32(int64(opus_int16(warping_Q16)))) >> 16)
				psDD.SAR2_Q14[j+0] = tmp2
				n_AR_Q14 = n_AR_Q14 + ((tmp2 * opus_int32(int64(AR_shp_Q13[j]))) >> 16)
			}
			psDD.SAR2_Q14[shapingLPCOrder-1] = tmp1
			n_AR_Q14 = n_AR_Q14 + ((tmp1 * opus_int32(int64(AR_shp_Q13[shapingLPCOrder-1]))) >> 16)
			n_AR_Q14 = opus_int32(opus_uint32(n_AR_Q14) << 1)
			n_AR_Q14 = n_AR_Q14 + ((psDD.LF_AR_Q14 * opus_int32(int64(opus_int16(Tilt_Q14)))) >> 16)
			n_AR_Q14 = opus_int32(opus_uint32(n_AR_Q14) << 2)
			n_LF_Q14 = ((psDD.Shape_Q14[*smpl_buf_idx]) * opus_int32(int64(opus_int16(LF_shp_Q14)))) >> 16
			n_LF_Q14 = n_LF_Q14 + ((psDD.LF_AR_Q14 * opus_int32(int64(LF_shp_Q14)>>16)) >> 16)
			n_LF_Q14 = opus_int32(opus_uint32(n_LF_Q14) << 2)
			if ((opus_uint32(n_AR_Q14) + opus_uint32(n_LF_Q14)) & 0x80000000) == 0 {
				if ((n_AR_Q14 & n_LF_Q14) & 0x80000000) != 0 {
					tmp1 = 0x80000000
				} else {
					tmp1 = n_AR_Q14 + n_LF_Q14
				}
			} else if ((n_AR_Q14 | n_LF_Q14) & 0x80000000) == 0 {
				tmp1 = silk_int32_MAX
			} else {
				tmp1 = n_AR_Q14 + n_LF_Q14
			}
			tmp2 = n_LTP_Q14 + LPC_pred_Q14
			if ((opus_uint32(tmp2) - opus_uint32(tmp1)) & 0x80000000) == 0 {
				if (tmp2 & (tmp1 ^ 0x80000000) & 0x80000000) != 0 {
					tmp1 = 0x80000000
				} else {
					tmp1 = tmp2 - tmp1
				}
			} else if ((tmp2 ^ 0x80000000) & tmp1 & 0x80000000) != 0 {
				tmp1 = silk_int32_MAX
			} else {
				tmp1 = tmp2 - tmp1
			}
			if 4 == 1 {
				tmp1 = (tmp1 >> 1) + (tmp1 & 1)
			} else {
				tmp1 = ((tmp1 >> (4 - 1)) + 1) >> 1
			}
			r_Q10 = (x_Q10[i]) - tmp1
			if psDD.Seed < 0 {
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
				rd1_Q10 = opus_int32(opus_int16(q1_Q10)) * opus_int32(opus_int16(Lambda_Q10))
				rd2_Q10 = opus_int32(opus_int16(q2_Q10)) * opus_int32(opus_int16(Lambda_Q10))
			} else if q1_Q0 == 0 {
				q1_Q10 = opus_int32(offset_Q10)
				q2_Q10 = q1_Q10 + opus_int32(1024-QUANT_LEVEL_ADJUST_Q10)
				rd1_Q10 = opus_int32(opus_int16(q1_Q10)) * opus_int32(opus_int16(Lambda_Q10))
				rd2_Q10 = opus_int32(opus_int16(q2_Q10)) * opus_int32(opus_int16(Lambda_Q10))
			} else if q1_Q0 == -1 {
				q2_Q10 = opus_int32(offset_Q10)
				q1_Q10 = q2_Q10 - opus_int32(1024-QUANT_LEVEL_ADJUST_Q10)
				rd1_Q10 = opus_int32(opus_int16(-q1_Q10)) * opus_int32(opus_int16(Lambda_Q10))
				rd2_Q10 = opus_int32(opus_int16(q2_Q10)) * opus_int32(opus_int16(Lambda_Q10))
			} else {
				q1_Q10 = (opus_int32(opus_uint32(q1_Q0) << 10)) + QUANT_LEVEL_ADJUST_Q10
				q1_Q10 = q1_Q10 + opus_int32(offset_Q10)
				q2_Q10 = q1_Q10 + 1024
				rd1_Q10 = opus_int32(opus_int16(-q1_Q10)) * opus_int32(opus_int16(Lambda_Q10))
				rd2_Q10 = opus_int32(opus_int16(-q2_Q10)) * opus_int32(opus_int16(Lambda_Q10))
			}
			rr_Q10 = r_Q10 - q1_Q10
			rd1_Q10 = (rd1_Q10 + (opus_int32(opus_int16(rr_Q10)))*opus_int32(opus_int16(rr_Q10))) >> 10
			rr_Q10 = r_Q10 - q2_Q10
			rd2_Q10 = (rd2_Q10 + (opus_int32(opus_int16(rr_Q10)))*opus_int32(opus_int16(rr_Q10))) >> 10
			if rd1_Q10 < rd2_Q10 {
				(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*0))).RD_Q10 = psDD.RD_Q10 + rd1_Q10
				(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*1))).RD_Q10 = psDD.RD_Q10 + rd2_Q10
				(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*0))).Q_Q10 = q1_Q10
				(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*1))).Q_Q10 = q2_Q10
			} else {
				(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*0))).RD_Q10 = psDD.RD_Q10 + rd2_Q10
				(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*1))).RD_Q10 = psDD.RD_Q10 + rd1_Q10
				(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*0))).Q_Q10 = q2_Q10
				(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*1))).Q_Q10 = q1_Q10
			}
			exc_Q14 = opus_int32(opus_uint32((*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*0))).Q_Q10) << 4)
			if psDD.Seed < 0 {
				exc_Q14 = -exc_Q14
			}
			LPC_exc_Q14 = exc_Q14 + LTP_pred_Q14
			xq_Q14 = LPC_exc_Q14 + LPC_pred_Q14
			(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*0))).Diff_Q14 = xq_Q14 - (opus_int32(opus_uint32(x_Q10[i]) << 4))
			sLF_AR_shp_Q14 = (*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*0))).Diff_Q14 - n_AR_Q14
			if ((opus_uint32(sLF_AR_shp_Q14) - opus_uint32(n_LF_Q14)) & 0x80000000) == 0 {
				if (sLF_AR_shp_Q14 & (n_LF_Q14 ^ 0x80000000) & 0x80000000) != 0 {
					(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*0))).SLTP_shp_Q14 = 0x80000000
				} else {
					(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*0))).SLTP_shp_Q14 = sLF_AR_shp_Q14 - n_LF_Q14
				}
			} else if ((sLF_AR_shp_Q14 ^ 0x80000000) & n_LF_Q14 & 0x80000000) != 0 {
				(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*0))).SLTP_shp_Q14 = silk_int32_MAX
			} else {
				(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*0))).SLTP_shp_Q14 = sLF_AR_shp_Q14 - n_LF_Q14
			}
			(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*0))).LF_AR_Q14 = sLF_AR_shp_Q14
			(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*0))).LPC_exc_Q14 = LPC_exc_Q14
			(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*0))).Xq_Q14 = xq_Q14
			exc_Q14 = opus_int32(opus_uint32((*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*1))).Q_Q10) << 4)
			if psDD.Seed < 0 {
				exc_Q14 = -exc_Q14
			}
			LPC_exc_Q14 = exc_Q14 + LTP_pred_Q14
			xq_Q14 = LPC_exc_Q14 + LPC_pred_Q14
			(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*1))).Diff_Q14 = xq_Q14 - (opus_int32(opus_uint32(x_Q10[i]) << 4))
			sLF_AR_shp_Q14 = (*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*1))).Diff_Q14 - n_AR_Q14
			if ((opus_uint32(sLF_AR_shp_Q14) - opus_uint32(n_LF_Q14)) & 0x80000000) == 0 {
				if (sLF_AR_shp_Q14 & (n_LF_Q14 ^ 0x80000000) & 0x80000000) != 0 {
					(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*1))).SLTP_shp_Q14 = 0x80000000
				} else {
					(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*1))).SLTP_shp_Q14 = sLF_AR_shp_Q14 - n_LF_Q14
				}
			} else if ((sLF_AR_shp_Q14 ^ 0x80000000) & n_LF_Q14 & 0x80000000) != 0 {
				(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*1))).SLTP_shp_Q14 = silk_int32_MAX
			} else {
				(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*1))).SLTP_shp_Q14 = sLF_AR_shp_Q14 - n_LF_Q14
			}
			(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*1))).LF_AR_Q14 = sLF_AR_shp_Q14
			(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*1))).LPC_exc_Q14 = LPC_exc_Q14
			(*(*NSQ_sample_struct)(unsafe.Add(unsafe.Pointer(psSS), unsafe.Sizeof(NSQ_sample_struct{})*1))).Xq_Q14 = xq_Q14
		}
		*smpl_buf_idx = (*smpl_buf_idx - 1) % DECISION_DELAY
		if *smpl_buf_idx < 0 {
			*smpl_buf_idx += DECISION_DELAY
		}
		last_smple_idx = (*smpl_buf_idx + decisionDelay) % DECISION_DELAY
		RDmin_Q10 = (*(*NSQ_sample_pair)(unsafe.Add(unsafe.Pointer(psSampleState), unsafe.Sizeof(NSQ_sample_pair{})*0)))[0].RD_Q10
		Winner_ind = 0
		for k = 1; k < nStatesDelayedDecision; k++ {
			if (*(*NSQ_sample_pair)(unsafe.Add(unsafe.Pointer(psSampleState), unsafe.Sizeof(NSQ_sample_pair{})*uintptr(k))))[0].RD_Q10 < RDmin_Q10 {
				RDmin_Q10 = (*(*NSQ_sample_pair)(unsafe.Add(unsafe.Pointer(psSampleState), unsafe.Sizeof(NSQ_sample_pair{})*uintptr(k))))[0].RD_Q10
				Winner_ind = k
			}
		}
		Winner_rand_state = psDelDec[Winner_ind].RandState[last_smple_idx]
		for k = 0; k < nStatesDelayedDecision; k++ {
			if psDelDec[k].RandState[last_smple_idx] != Winner_rand_state {
				(*(*NSQ_sample_pair)(unsafe.Add(unsafe.Pointer(psSampleState), unsafe.Sizeof(NSQ_sample_pair{})*uintptr(k))))[0].RD_Q10 = (*(*NSQ_sample_pair)(unsafe.Add(unsafe.Pointer(psSampleState), unsafe.Sizeof(NSQ_sample_pair{})*uintptr(k))))[0].RD_Q10 + opus_int32(silk_int32_MAX>>4)
				(*(*NSQ_sample_pair)(unsafe.Add(unsafe.Pointer(psSampleState), unsafe.Sizeof(NSQ_sample_pair{})*uintptr(k))))[1].RD_Q10 = (*(*NSQ_sample_pair)(unsafe.Add(unsafe.Pointer(psSampleState), unsafe.Sizeof(NSQ_sample_pair{})*uintptr(k))))[1].RD_Q10 + opus_int32(silk_int32_MAX>>4)
			}
		}
		RDmax_Q10 = (*(*NSQ_sample_pair)(unsafe.Add(unsafe.Pointer(psSampleState), unsafe.Sizeof(NSQ_sample_pair{})*0)))[0].RD_Q10
		RDmin_Q10 = (*(*NSQ_sample_pair)(unsafe.Add(unsafe.Pointer(psSampleState), unsafe.Sizeof(NSQ_sample_pair{})*0)))[1].RD_Q10
		RDmax_ind = 0
		RDmin_ind = 0
		for k = 1; k < nStatesDelayedDecision; k++ {
			if (*(*NSQ_sample_pair)(unsafe.Add(unsafe.Pointer(psSampleState), unsafe.Sizeof(NSQ_sample_pair{})*uintptr(k))))[0].RD_Q10 > RDmax_Q10 {
				RDmax_Q10 = (*(*NSQ_sample_pair)(unsafe.Add(unsafe.Pointer(psSampleState), unsafe.Sizeof(NSQ_sample_pair{})*uintptr(k))))[0].RD_Q10
				RDmax_ind = k
			}
			if (*(*NSQ_sample_pair)(unsafe.Add(unsafe.Pointer(psSampleState), unsafe.Sizeof(NSQ_sample_pair{})*uintptr(k))))[1].RD_Q10 < RDmin_Q10 {
				RDmin_Q10 = (*(*NSQ_sample_pair)(unsafe.Add(unsafe.Pointer(psSampleState), unsafe.Sizeof(NSQ_sample_pair{})*uintptr(k))))[1].RD_Q10
				RDmin_ind = k
			}
		}
		if RDmin_Q10 < RDmax_Q10 {
			libc.MemCpy(unsafe.Pointer((*opus_int32)(unsafe.Add(unsafe.Pointer((*opus_int32)(unsafe.Pointer(&psDelDec[RDmax_ind]))), unsafe.Sizeof(opus_int32(0))*uintptr(i)))), unsafe.Pointer((*opus_int32)(unsafe.Add(unsafe.Pointer((*opus_int32)(unsafe.Pointer(&psDelDec[RDmin_ind]))), unsafe.Sizeof(opus_int32(0))*uintptr(i)))), int(unsafe.Sizeof(NSQ_del_dec_struct{})-uintptr(i*int64(unsafe.Sizeof(opus_int32(0))))))
			libc.MemCpy(unsafe.Pointer(&(*(*NSQ_sample_pair)(unsafe.Add(unsafe.Pointer(psSampleState), unsafe.Sizeof(NSQ_sample_pair{})*uintptr(RDmax_ind))))[0]), unsafe.Pointer(&(*(*NSQ_sample_pair)(unsafe.Add(unsafe.Pointer(psSampleState), unsafe.Sizeof(NSQ_sample_pair{})*uintptr(RDmin_ind))))[1]), int(unsafe.Sizeof(NSQ_sample_struct{})))
		}
		psDD = &psDelDec[Winner_ind]
		if subfr > 0 || i >= decisionDelay {
			if 10 == 1 {
				pulses[i-decisionDelay] = int8(((psDD.Q_Q10[last_smple_idx]) >> 1) + ((psDD.Q_Q10[last_smple_idx]) & 1))
			} else {
				pulses[i-decisionDelay] = int8((((psDD.Q_Q10[last_smple_idx]) >> (10 - 1)) + 1) >> 1)
			}
			if (func() opus_int32 {
				if 8 == 1 {
					return ((opus_int32((int64(psDD.Xq_Q14[last_smple_idx]) * int64(delayedGain_Q10[last_smple_idx])) >> 16)) >> 1) + ((opus_int32((int64(psDD.Xq_Q14[last_smple_idx]) * int64(delayedGain_Q10[last_smple_idx])) >> 16)) & 1)
				}
				return (((opus_int32((int64(psDD.Xq_Q14[last_smple_idx]) * int64(delayedGain_Q10[last_smple_idx])) >> 16)) >> (8 - 1)) + 1) >> 1
			}()) > silk_int16_MAX {
				xq[i-decisionDelay] = silk_int16_MAX
			} else if (func() opus_int32 {
				if 8 == 1 {
					return ((opus_int32((int64(psDD.Xq_Q14[last_smple_idx]) * int64(delayedGain_Q10[last_smple_idx])) >> 16)) >> 1) + ((opus_int32((int64(psDD.Xq_Q14[last_smple_idx]) * int64(delayedGain_Q10[last_smple_idx])) >> 16)) & 1)
				}
				return (((opus_int32((int64(psDD.Xq_Q14[last_smple_idx]) * int64(delayedGain_Q10[last_smple_idx])) >> 16)) >> (8 - 1)) + 1) >> 1
			}()) < opus_int32(math.MinInt16) {
				xq[i-decisionDelay] = math.MinInt16
			} else if 8 == 1 {
				xq[i-decisionDelay] = opus_int16(((opus_int32((int64(psDD.Xq_Q14[last_smple_idx]) * int64(delayedGain_Q10[last_smple_idx])) >> 16)) >> 1) + ((opus_int32((int64(psDD.Xq_Q14[last_smple_idx]) * int64(delayedGain_Q10[last_smple_idx])) >> 16)) & 1))
			} else {
				xq[i-decisionDelay] = opus_int16((((opus_int32((int64(psDD.Xq_Q14[last_smple_idx]) * int64(delayedGain_Q10[last_smple_idx])) >> 16)) >> (8 - 1)) + 1) >> 1)
			}
			NSQ.SLTP_shp_Q14[NSQ.SLTP_shp_buf_idx-decisionDelay] = psDD.Shape_Q14[last_smple_idx]
			sLTP_Q15[NSQ.SLTP_buf_idx-decisionDelay] = psDD.Pred_Q15[last_smple_idx]
		}
		NSQ.SLTP_shp_buf_idx++
		NSQ.SLTP_buf_idx++
		for k = 0; k < nStatesDelayedDecision; k++ {
			psDD = &psDelDec[k]
			psSS = &(*(*NSQ_sample_pair)(unsafe.Add(unsafe.Pointer(psSampleState), unsafe.Sizeof(NSQ_sample_pair{})*uintptr(k))))[0]
			psDD.LF_AR_Q14 = psSS.LF_AR_Q14
			psDD.Diff_Q14 = psSS.Diff_Q14
			psDD.SLPC_Q14[MAX_LPC_ORDER+i] = psSS.Xq_Q14
			psDD.Xq_Q14[*smpl_buf_idx] = psSS.Xq_Q14
			psDD.Q_Q10[*smpl_buf_idx] = psSS.Q_Q10
			psDD.Pred_Q15[*smpl_buf_idx] = opus_int32(opus_uint32(psSS.LPC_exc_Q14) << 1)
			psDD.Shape_Q14[*smpl_buf_idx] = psSS.SLTP_shp_Q14
			psDD.Seed = opus_int32(opus_uint32(psDD.Seed) + opus_uint32(func() opus_int32 {
				if 10 == 1 {
					return (psSS.Q_Q10 >> 1) + (psSS.Q_Q10 & 1)
				}
				return ((psSS.Q_Q10 >> (10 - 1)) + 1) >> 1
			}()))
			psDD.RandState[*smpl_buf_idx] = psDD.Seed
			psDD.RD_Q10 = psSS.RD_Q10
		}
		delayedGain_Q10[*smpl_buf_idx] = Gain_Q10
	}
	for k = 0; k < nStatesDelayedDecision; k++ {
		psDD = &psDelDec[k]
		libc.MemCpy(unsafe.Pointer(&psDD.SLPC_Q14[0]), unsafe.Pointer(&psDD.SLPC_Q14[length]), int(MAX_LPC_ORDER*unsafe.Sizeof(opus_int32(0))))
	}
}
func silk_nsq_del_dec_scale_states(psEncC *silk_encoder_state, NSQ *silk_nsq_state, psDelDec [0]NSQ_del_dec_struct, x16 [0]opus_int16, x_sc_Q10 [0]opus_int32, sLTP [0]opus_int16, sLTP_Q15 [0]opus_int32, subfr int64, nStatesDelayedDecision int64, LTP_scale_Q14 int64, Gains_Q16 [4]opus_int32, pitchL [4]int64, signal_type int64, decisionDelay int64) {
	var (
		i            int64
		k            int64
		lag          int64
		gain_adj_Q16 opus_int32
		inv_gain_Q31 opus_int32
		inv_gain_Q26 opus_int32
		psDD         *NSQ_del_dec_struct
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
			for i = NSQ.SLTP_buf_idx - lag - LTP_ORDER/2; i < NSQ.SLTP_buf_idx-decisionDelay; i++ {
				sLTP_Q15[i] = opus_int32((int64(gain_adj_Q16) * int64(sLTP_Q15[i])) >> 16)
			}
		}
		for k = 0; k < nStatesDelayedDecision; k++ {
			psDD = &psDelDec[k]
			psDD.LF_AR_Q14 = opus_int32((int64(gain_adj_Q16) * int64(psDD.LF_AR_Q14)) >> 16)
			psDD.Diff_Q14 = opus_int32((int64(gain_adj_Q16) * int64(psDD.Diff_Q14)) >> 16)
			for i = 0; i < MAX_LPC_ORDER; i++ {
				psDD.SLPC_Q14[i] = opus_int32((int64(gain_adj_Q16) * int64(psDD.SLPC_Q14[i])) >> 16)
			}
			for i = 0; i < MAX_SHAPE_LPC_ORDER; i++ {
				psDD.SAR2_Q14[i] = opus_int32((int64(gain_adj_Q16) * int64(psDD.SAR2_Q14[i])) >> 16)
			}
			for i = 0; i < DECISION_DELAY; i++ {
				psDD.Pred_Q15[i] = opus_int32((int64(gain_adj_Q16) * int64(psDD.Pred_Q15[i])) >> 16)
				psDD.Shape_Q14[i] = opus_int32((int64(gain_adj_Q16) * int64(psDD.Shape_Q14[i])) >> 16)
			}
		}
		NSQ.Prev_gain_Q16 = Gains_Q16[subfr]
	}
}
