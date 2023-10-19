package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

func silk_stereo_LR_to_MS(state *stereo_enc_state, x1 [0]opus_int16, x2 [0]opus_int16, ix [2][3]int8, mid_only_flag *int8, mid_side_rates_bps [0]opus_int32, total_rate_bps opus_int32, prev_speech_act_Q8 int64, toMono int64, fs_kHz int64, frame_length int64) {
	var (
		n                int64
		is10msFrame      int64
		denom_Q16        int64
		delta0_Q13       int64
		delta1_Q13       int64
		sum              opus_int32
		diff             opus_int32
		smooth_coef_Q16  opus_int32
		pred_Q13         [2]opus_int32
		pred0_Q13        opus_int32
		pred1_Q13        opus_int32
		LP_ratio_Q14     opus_int32
		HP_ratio_Q14     opus_int32
		frac_Q16         opus_int32
		frac_3_Q16       opus_int32
		min_mid_rate_bps opus_int32
		width_Q14        opus_int32
		w_Q24            opus_int32
		deltaw_Q24       opus_int32
		side             *opus_int16
		LP_mid           *opus_int16
		HP_mid           *opus_int16
		LP_side          *opus_int16
		HP_side          *opus_int16
		mid              *opus_int16 = &x1[-2]
	)
	side = (*opus_int16)(libc.Malloc(int((frame_length + 2) * int64(unsafe.Sizeof(opus_int16(0))))))
	for n = 0; n < frame_length+2; n++ {
		sum = opus_int32(x1[n-2]) + opus_int32(x2[n-2])
		diff = opus_int32(x1[n-2]) - opus_int32(x2[n-2])
		if 1 == 1 {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(mid), unsafe.Sizeof(opus_int16(0))*uintptr(n))) = opus_int16((sum >> 1) + (sum & 1))
		} else {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(mid), unsafe.Sizeof(opus_int16(0))*uintptr(n))) = opus_int16(((sum >> (1 - 1)) + 1) >> 1)
		}
		if (func() opus_int32 {
			if 1 == 1 {
				return (diff >> 1) + (diff & 1)
			}
			return ((diff >> (1 - 1)) + 1) >> 1
		}()) > silk_int16_MAX {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(side), unsafe.Sizeof(opus_int16(0))*uintptr(n))) = silk_int16_MAX
		} else if (func() opus_int32 {
			if 1 == 1 {
				return (diff >> 1) + (diff & 1)
			}
			return ((diff >> (1 - 1)) + 1) >> 1
		}()) < opus_int32(math.MinInt16) {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(side), unsafe.Sizeof(opus_int16(0))*uintptr(n))) = math.MinInt16
		} else if 1 == 1 {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(side), unsafe.Sizeof(opus_int16(0))*uintptr(n))) = opus_int16((diff >> 1) + (diff & 1))
		} else {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(side), unsafe.Sizeof(opus_int16(0))*uintptr(n))) = opus_int16(((diff >> (1 - 1)) + 1) >> 1)
		}
	}
	libc.MemCpy(unsafe.Pointer(mid), unsafe.Pointer(&state.SMid[0]), int(2*unsafe.Sizeof(opus_int16(0))))
	libc.MemCpy(unsafe.Pointer(side), unsafe.Pointer(&state.SSide[0]), int(2*unsafe.Sizeof(opus_int16(0))))
	libc.MemCpy(unsafe.Pointer(&state.SMid[0]), unsafe.Pointer((*opus_int16)(unsafe.Add(unsafe.Pointer(mid), unsafe.Sizeof(opus_int16(0))*uintptr(frame_length)))), int(2*unsafe.Sizeof(opus_int16(0))))
	libc.MemCpy(unsafe.Pointer(&state.SSide[0]), unsafe.Pointer((*opus_int16)(unsafe.Add(unsafe.Pointer(side), unsafe.Sizeof(opus_int16(0))*uintptr(frame_length)))), int(2*unsafe.Sizeof(opus_int16(0))))
	LP_mid = (*opus_int16)(libc.Malloc(int(frame_length * int64(unsafe.Sizeof(opus_int16(0))))))
	HP_mid = (*opus_int16)(libc.Malloc(int(frame_length * int64(unsafe.Sizeof(opus_int16(0))))))
	for n = 0; n < frame_length; n++ {
		if 2 == 1 {
			sum = (((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(mid), unsafe.Sizeof(opus_int16(0))*uintptr(n)))) + opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(mid), unsafe.Sizeof(opus_int16(0))*uintptr(n+2))))) + (opus_int32(opus_uint32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(mid), unsafe.Sizeof(opus_int16(0))*uintptr(n+1)))) << 1))) >> 1) + (((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(mid), unsafe.Sizeof(opus_int16(0))*uintptr(n)))) + opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(mid), unsafe.Sizeof(opus_int16(0))*uintptr(n+2))))) + (opus_int32(opus_uint32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(mid), unsafe.Sizeof(opus_int16(0))*uintptr(n+1)))) << 1))) & 1)
		} else {
			sum = ((((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(mid), unsafe.Sizeof(opus_int16(0))*uintptr(n)))) + opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(mid), unsafe.Sizeof(opus_int16(0))*uintptr(n+2))))) + (opus_int32(opus_uint32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(mid), unsafe.Sizeof(opus_int16(0))*uintptr(n+1)))) << 1))) >> (2 - 1)) + 1) >> 1
		}
		*(*opus_int16)(unsafe.Add(unsafe.Pointer(LP_mid), unsafe.Sizeof(opus_int16(0))*uintptr(n))) = opus_int16(sum)
		*(*opus_int16)(unsafe.Add(unsafe.Pointer(HP_mid), unsafe.Sizeof(opus_int16(0))*uintptr(n))) = opus_int16(opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(mid), unsafe.Sizeof(opus_int16(0))*uintptr(n+1)))) - sum)
	}
	LP_side = (*opus_int16)(libc.Malloc(int(frame_length * int64(unsafe.Sizeof(opus_int16(0))))))
	HP_side = (*opus_int16)(libc.Malloc(int(frame_length * int64(unsafe.Sizeof(opus_int16(0))))))
	for n = 0; n < frame_length; n++ {
		if 2 == 1 {
			sum = (((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(side), unsafe.Sizeof(opus_int16(0))*uintptr(n)))) + opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(side), unsafe.Sizeof(opus_int16(0))*uintptr(n+2))))) + (opus_int32(opus_uint32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(side), unsafe.Sizeof(opus_int16(0))*uintptr(n+1)))) << 1))) >> 1) + (((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(side), unsafe.Sizeof(opus_int16(0))*uintptr(n)))) + opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(side), unsafe.Sizeof(opus_int16(0))*uintptr(n+2))))) + (opus_int32(opus_uint32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(side), unsafe.Sizeof(opus_int16(0))*uintptr(n+1)))) << 1))) & 1)
		} else {
			sum = ((((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(side), unsafe.Sizeof(opus_int16(0))*uintptr(n)))) + opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(side), unsafe.Sizeof(opus_int16(0))*uintptr(n+2))))) + (opus_int32(opus_uint32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(side), unsafe.Sizeof(opus_int16(0))*uintptr(n+1)))) << 1))) >> (2 - 1)) + 1) >> 1
		}
		*(*opus_int16)(unsafe.Add(unsafe.Pointer(LP_side), unsafe.Sizeof(opus_int16(0))*uintptr(n))) = opus_int16(sum)
		*(*opus_int16)(unsafe.Add(unsafe.Pointer(HP_side), unsafe.Sizeof(opus_int16(0))*uintptr(n))) = opus_int16(opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(side), unsafe.Sizeof(opus_int16(0))*uintptr(n+1)))) - sum)
	}
	is10msFrame = int64(libc.BoolToInt(frame_length == fs_kHz*10))
	if is10msFrame != 0 {
		smooth_coef_Q16 = opus_int32((STEREO_RATIO_SMOOTH_COEF/2)*(1<<16) + 0.5)
	} else {
		smooth_coef_Q16 = opus_int32(STEREO_RATIO_SMOOTH_COEF*(1<<16) + 0.5)
	}
	smooth_coef_Q16 = ((opus_int32(opus_int16(prev_speech_act_Q8)) * opus_int32(opus_int16(prev_speech_act_Q8))) * opus_int32(int64(opus_int16(smooth_coef_Q16)))) >> 16
	pred_Q13[0] = silk_stereo_find_predictor(&LP_ratio_Q14, [0]opus_int16(LP_mid), [0]opus_int16(LP_side), [0]opus_int32(&state.Mid_side_amp_Q0[0]), frame_length, int64(smooth_coef_Q16))
	pred_Q13[1] = silk_stereo_find_predictor(&HP_ratio_Q14, [0]opus_int16(HP_mid), [0]opus_int16(HP_side), [0]opus_int32(&state.Mid_side_amp_Q0[2]), frame_length, int64(smooth_coef_Q16))
	frac_Q16 = HP_ratio_Q14 + (opus_int32(opus_int16(LP_ratio_Q14)))*3
	if frac_Q16 < (opus_int32(1*(1<<16) + 0.5)) {
		frac_Q16 = frac_Q16
	} else {
		frac_Q16 = opus_int32(1*(1<<16) + 0.5)
	}
	if is10msFrame != 0 {
		total_rate_bps -= 1200
	} else {
		total_rate_bps -= 600
	}
	if total_rate_bps < 1 {
		total_rate_bps = 1
	}
	min_mid_rate_bps = (opus_int32(opus_int16(fs_kHz)))*600 + 2000
	frac_3_Q16 = frac_Q16 * 3
	mid_side_rates_bps[0] = silk_DIV32_varQ(total_rate_bps, frac_3_Q16+(opus_int32((8+5)*(1<<16)+0.5)), 16+3)
	if mid_side_rates_bps[0] < min_mid_rate_bps {
		mid_side_rates_bps[0] = min_mid_rate_bps
		mid_side_rates_bps[1] = total_rate_bps - mid_side_rates_bps[0]
		width_Q14 = silk_DIV32_varQ((opus_int32(opus_uint32(mid_side_rates_bps[1])<<1))-min_mid_rate_bps, ((frac_3_Q16+(opus_int32(1*(1<<16)+0.5)))*opus_int32(int64(opus_int16(min_mid_rate_bps))))>>16, 14+2)
		if 0 > (opus_int32(1*(1<<14) + 0.5)) {
			if width_Q14 > 0 {
				width_Q14 = 0
			} else if width_Q14 < (opus_int32(1*(1<<14) + 0.5)) {
				width_Q14 = opus_int32(1*(1<<14) + 0.5)
			} else {
				width_Q14 = width_Q14
			}
		} else if width_Q14 > (opus_int32(1*(1<<14) + 0.5)) {
			width_Q14 = opus_int32(1*(1<<14) + 0.5)
		} else if width_Q14 < 0 {
			width_Q14 = 0
		} else {
			width_Q14 = width_Q14
		}
	} else {
		mid_side_rates_bps[1] = total_rate_bps - mid_side_rates_bps[0]
		width_Q14 = opus_int32(1*(1<<14) + 0.5)
	}
	state.Smth_width_Q14 = opus_int16(opus_int32(state.Smth_width_Q14) + (((width_Q14 - opus_int32(state.Smth_width_Q14)) * opus_int32(int64(opus_int16(smooth_coef_Q16)))) >> 16))
	*mid_only_flag = 0
	if toMono != 0 {
		width_Q14 = 0
		pred_Q13[0] = 0
		pred_Q13[1] = 0
		silk_stereo_quant_pred(pred_Q13[:], ix)
	} else if state.Width_prev_Q14 == 0 && (total_rate_bps*8 < min_mid_rate_bps*13 || ((frac_Q16*opus_int32(int64(state.Smth_width_Q14)))>>16) < (opus_int32(0.05*(1<<14)+0.5))) {
		pred_Q13[0] = (opus_int32(state.Smth_width_Q14) * opus_int32(opus_int16(pred_Q13[0]))) >> 14
		pred_Q13[1] = (opus_int32(state.Smth_width_Q14) * opus_int32(opus_int16(pred_Q13[1]))) >> 14
		silk_stereo_quant_pred(pred_Q13[:], ix)
		width_Q14 = 0
		pred_Q13[0] = 0
		pred_Q13[1] = 0
		mid_side_rates_bps[0] = total_rate_bps
		mid_side_rates_bps[1] = 0
		*mid_only_flag = 1
	} else if state.Width_prev_Q14 != 0 && (total_rate_bps*8 < min_mid_rate_bps*11 || ((frac_Q16*opus_int32(int64(state.Smth_width_Q14)))>>16) < (opus_int32(0.02*(1<<14)+0.5))) {
		pred_Q13[0] = (opus_int32(state.Smth_width_Q14) * opus_int32(opus_int16(pred_Q13[0]))) >> 14
		pred_Q13[1] = (opus_int32(state.Smth_width_Q14) * opus_int32(opus_int16(pred_Q13[1]))) >> 14
		silk_stereo_quant_pred(pred_Q13[:], ix)
		width_Q14 = 0
		pred_Q13[0] = 0
		pred_Q13[1] = 0
	} else if opus_int32(state.Smth_width_Q14) > (opus_int32(0.95*(1<<14) + 0.5)) {
		silk_stereo_quant_pred(pred_Q13[:], ix)
		width_Q14 = opus_int32(1*(1<<14) + 0.5)
	} else {
		pred_Q13[0] = (opus_int32(state.Smth_width_Q14) * opus_int32(opus_int16(pred_Q13[0]))) >> 14
		pred_Q13[1] = (opus_int32(state.Smth_width_Q14) * opus_int32(opus_int16(pred_Q13[1]))) >> 14
		silk_stereo_quant_pred(pred_Q13[:], ix)
		width_Q14 = opus_int32(state.Smth_width_Q14)
	}
	if int64(*mid_only_flag) == 1 {
		state.Silent_side_len += opus_int16(frame_length - STEREO_INTERP_LEN_MS*fs_kHz)
		if int64(state.Silent_side_len) < LA_SHAPE_MS*fs_kHz {
			*mid_only_flag = 0
		} else {
			state.Silent_side_len = 10000
		}
	} else {
		state.Silent_side_len = 0
	}
	if int64(*mid_only_flag) == 0 && mid_side_rates_bps[1] < 1 {
		mid_side_rates_bps[1] = 1
		mid_side_rates_bps[0] = opus_int32(silk_max_int(1, int64(total_rate_bps-mid_side_rates_bps[1])))
	}
	pred0_Q13 = opus_int32(int64(-state.Pred_prev_Q13[0]))
	pred1_Q13 = opus_int32(int64(-state.Pred_prev_Q13[1]))
	w_Q24 = opus_int32(opus_uint32(state.Width_prev_Q14) << 10)
	denom_Q16 = int64(opus_int32((1 << 16) / (STEREO_INTERP_LEN_MS * fs_kHz)))
	if 16 == 1 {
		delta0_Q13 = int64(-(((opus_int32(opus_int16(pred_Q13[0]-opus_int32(state.Pred_prev_Q13[0]))) * opus_int32(opus_int16(denom_Q16))) >> 1) + ((opus_int32(opus_int16(pred_Q13[0]-opus_int32(state.Pred_prev_Q13[0]))) * opus_int32(opus_int16(denom_Q16))) & 1)))
	} else {
		delta0_Q13 = int64(-((((opus_int32(opus_int16(pred_Q13[0]-opus_int32(state.Pred_prev_Q13[0]))) * opus_int32(opus_int16(denom_Q16))) >> (16 - 1)) + 1) >> 1))
	}
	if 16 == 1 {
		delta1_Q13 = int64(-(((opus_int32(opus_int16(pred_Q13[1]-opus_int32(state.Pred_prev_Q13[1]))) * opus_int32(opus_int16(denom_Q16))) >> 1) + ((opus_int32(opus_int16(pred_Q13[1]-opus_int32(state.Pred_prev_Q13[1]))) * opus_int32(opus_int16(denom_Q16))) & 1)))
	} else {
		delta1_Q13 = int64(-((((opus_int32(opus_int16(pred_Q13[1]-opus_int32(state.Pred_prev_Q13[1]))) * opus_int32(opus_int16(denom_Q16))) >> (16 - 1)) + 1) >> 1))
	}
	deltaw_Q24 = opus_int32(opus_uint32(((width_Q14-opus_int32(state.Width_prev_Q14))*opus_int32(int64(opus_int16(denom_Q16))))>>16) << 10)
	for n = 0; n < STEREO_INTERP_LEN_MS*fs_kHz; n++ {
		pred0_Q13 += opus_int32(delta0_Q13)
		pred1_Q13 += opus_int32(delta1_Q13)
		w_Q24 += deltaw_Q24
		sum = opus_int32(opus_uint32((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(mid), unsafe.Sizeof(opus_int16(0))*uintptr(n))))+opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(mid), unsafe.Sizeof(opus_int16(0))*uintptr(n+2)))))+(opus_int32(opus_uint32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(mid), unsafe.Sizeof(opus_int16(0))*uintptr(n+1))))<<1))) << 9)
		sum = ((w_Q24 * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(side), unsafe.Sizeof(opus_int16(0))*uintptr(n+1)))))) >> 16) + ((sum * opus_int32(int64(opus_int16(pred0_Q13)))) >> 16)
		sum = sum + (((opus_int32(opus_uint32(opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(mid), unsafe.Sizeof(opus_int16(0))*uintptr(n+1))))) << 11)) * opus_int32(int64(opus_int16(pred1_Q13)))) >> 16)
		if (func() opus_int32 {
			if 8 == 1 {
				return (sum >> 1) + (sum & 1)
			}
			return ((sum >> (8 - 1)) + 1) >> 1
		}()) > silk_int16_MAX {
			x2[n-1] = silk_int16_MAX
		} else if (func() opus_int32 {
			if 8 == 1 {
				return (sum >> 1) + (sum & 1)
			}
			return ((sum >> (8 - 1)) + 1) >> 1
		}()) < opus_int32(math.MinInt16) {
			x2[n-1] = math.MinInt16
		} else if 8 == 1 {
			x2[n-1] = opus_int16((sum >> 1) + (sum & 1))
		} else {
			x2[n-1] = opus_int16(((sum >> (8 - 1)) + 1) >> 1)
		}
	}
	pred0_Q13 = -pred_Q13[0]
	pred1_Q13 = -pred_Q13[1]
	w_Q24 = opus_int32(opus_uint32(width_Q14) << 10)
	for n = STEREO_INTERP_LEN_MS * fs_kHz; n < frame_length; n++ {
		sum = opus_int32(opus_uint32((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(mid), unsafe.Sizeof(opus_int16(0))*uintptr(n))))+opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(mid), unsafe.Sizeof(opus_int16(0))*uintptr(n+2)))))+(opus_int32(opus_uint32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(mid), unsafe.Sizeof(opus_int16(0))*uintptr(n+1))))<<1))) << 9)
		sum = ((w_Q24 * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(side), unsafe.Sizeof(opus_int16(0))*uintptr(n+1)))))) >> 16) + ((sum * opus_int32(int64(opus_int16(pred0_Q13)))) >> 16)
		sum = sum + (((opus_int32(opus_uint32(opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(mid), unsafe.Sizeof(opus_int16(0))*uintptr(n+1))))) << 11)) * opus_int32(int64(opus_int16(pred1_Q13)))) >> 16)
		if (func() opus_int32 {
			if 8 == 1 {
				return (sum >> 1) + (sum & 1)
			}
			return ((sum >> (8 - 1)) + 1) >> 1
		}()) > silk_int16_MAX {
			x2[n-1] = silk_int16_MAX
		} else if (func() opus_int32 {
			if 8 == 1 {
				return (sum >> 1) + (sum & 1)
			}
			return ((sum >> (8 - 1)) + 1) >> 1
		}()) < opus_int32(math.MinInt16) {
			x2[n-1] = math.MinInt16
		} else if 8 == 1 {
			x2[n-1] = opus_int16((sum >> 1) + (sum & 1))
		} else {
			x2[n-1] = opus_int16(((sum >> (8 - 1)) + 1) >> 1)
		}
	}
	state.Pred_prev_Q13[0] = opus_int16(pred_Q13[0])
	state.Pred_prev_Q13[1] = opus_int16(pred_Q13[1])
	state.Width_prev_Q14 = opus_int16(width_Q14)
}
