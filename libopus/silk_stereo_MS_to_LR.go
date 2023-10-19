package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

func silk_stereo_MS_to_LR(state *stereo_dec_state, x1 [0]opus_int16, x2 [0]opus_int16, pred_Q13 [0]opus_int32, fs_kHz int64, frame_length int64) {
	var (
		n          int64
		denom_Q16  int64
		delta0_Q13 int64
		delta1_Q13 int64
		sum        opus_int32
		diff       opus_int32
		pred0_Q13  opus_int32
		pred1_Q13  opus_int32
	)
	libc.MemCpy(unsafe.Pointer(&x1[0]), unsafe.Pointer(&state.SMid[0]), int(2*unsafe.Sizeof(opus_int16(0))))
	libc.MemCpy(unsafe.Pointer(&x2[0]), unsafe.Pointer(&state.SSide[0]), int(2*unsafe.Sizeof(opus_int16(0))))
	libc.MemCpy(unsafe.Pointer(&state.SMid[0]), unsafe.Pointer(&x1[frame_length]), int(2*unsafe.Sizeof(opus_int16(0))))
	libc.MemCpy(unsafe.Pointer(&state.SSide[0]), unsafe.Pointer(&x2[frame_length]), int(2*unsafe.Sizeof(opus_int16(0))))
	pred0_Q13 = opus_int32(state.Pred_prev_Q13[0])
	pred1_Q13 = opus_int32(state.Pred_prev_Q13[1])
	denom_Q16 = int64(opus_int32((1 << 16) / (STEREO_INTERP_LEN_MS * fs_kHz)))
	if 16 == 1 {
		delta0_Q13 = int64(((opus_int32(opus_int16(pred_Q13[0]-opus_int32(state.Pred_prev_Q13[0]))) * opus_int32(opus_int16(denom_Q16))) >> 1) + ((opus_int32(opus_int16(pred_Q13[0]-opus_int32(state.Pred_prev_Q13[0]))) * opus_int32(opus_int16(denom_Q16))) & 1))
	} else {
		delta0_Q13 = int64((((opus_int32(opus_int16(pred_Q13[0]-opus_int32(state.Pred_prev_Q13[0]))) * opus_int32(opus_int16(denom_Q16))) >> (16 - 1)) + 1) >> 1)
	}
	if 16 == 1 {
		delta1_Q13 = int64(((opus_int32(opus_int16(pred_Q13[1]-opus_int32(state.Pred_prev_Q13[1]))) * opus_int32(opus_int16(denom_Q16))) >> 1) + ((opus_int32(opus_int16(pred_Q13[1]-opus_int32(state.Pred_prev_Q13[1]))) * opus_int32(opus_int16(denom_Q16))) & 1))
	} else {
		delta1_Q13 = int64((((opus_int32(opus_int16(pred_Q13[1]-opus_int32(state.Pred_prev_Q13[1]))) * opus_int32(opus_int16(denom_Q16))) >> (16 - 1)) + 1) >> 1)
	}
	for n = 0; n < STEREO_INTERP_LEN_MS*fs_kHz; n++ {
		pred0_Q13 += opus_int32(delta0_Q13)
		pred1_Q13 += opus_int32(delta1_Q13)
		sum = opus_int32(opus_uint32((opus_int32(x1[n])+opus_int32(x1[n+2]))+(opus_int32(opus_uint32(x1[n+1])<<1))) << 9)
		sum = (opus_int32(opus_uint32(opus_int32(x2[n+1])) << 8)) + ((sum * opus_int32(int64(opus_int16(pred0_Q13)))) >> 16)
		sum = sum + (((opus_int32(opus_uint32(opus_int32(x1[n+1])) << 11)) * opus_int32(int64(opus_int16(pred1_Q13)))) >> 16)
		if (func() opus_int32 {
			if 8 == 1 {
				return (sum >> 1) + (sum & 1)
			}
			return ((sum >> (8 - 1)) + 1) >> 1
		}()) > silk_int16_MAX {
			x2[n+1] = silk_int16_MAX
		} else if (func() opus_int32 {
			if 8 == 1 {
				return (sum >> 1) + (sum & 1)
			}
			return ((sum >> (8 - 1)) + 1) >> 1
		}()) < opus_int32(math.MinInt16) {
			x2[n+1] = math.MinInt16
		} else if 8 == 1 {
			x2[n+1] = opus_int16((sum >> 1) + (sum & 1))
		} else {
			x2[n+1] = opus_int16(((sum >> (8 - 1)) + 1) >> 1)
		}
	}
	pred0_Q13 = pred_Q13[0]
	pred1_Q13 = pred_Q13[1]
	for n = STEREO_INTERP_LEN_MS * fs_kHz; n < frame_length; n++ {
		sum = opus_int32(opus_uint32((opus_int32(x1[n])+opus_int32(x1[n+2]))+(opus_int32(opus_uint32(x1[n+1])<<1))) << 9)
		sum = (opus_int32(opus_uint32(opus_int32(x2[n+1])) << 8)) + ((sum * opus_int32(int64(opus_int16(pred0_Q13)))) >> 16)
		sum = sum + (((opus_int32(opus_uint32(opus_int32(x1[n+1])) << 11)) * opus_int32(int64(opus_int16(pred1_Q13)))) >> 16)
		if (func() opus_int32 {
			if 8 == 1 {
				return (sum >> 1) + (sum & 1)
			}
			return ((sum >> (8 - 1)) + 1) >> 1
		}()) > silk_int16_MAX {
			x2[n+1] = silk_int16_MAX
		} else if (func() opus_int32 {
			if 8 == 1 {
				return (sum >> 1) + (sum & 1)
			}
			return ((sum >> (8 - 1)) + 1) >> 1
		}()) < opus_int32(math.MinInt16) {
			x2[n+1] = math.MinInt16
		} else if 8 == 1 {
			x2[n+1] = opus_int16((sum >> 1) + (sum & 1))
		} else {
			x2[n+1] = opus_int16(((sum >> (8 - 1)) + 1) >> 1)
		}
	}
	state.Pred_prev_Q13[0] = opus_int16(pred_Q13[0])
	state.Pred_prev_Q13[1] = opus_int16(pred_Q13[1])
	for n = 0; n < frame_length; n++ {
		sum = opus_int32(x1[n+1]) + opus_int32(x2[n+1])
		diff = opus_int32(x1[n+1]) - opus_int32(x2[n+1])
		if sum > silk_int16_MAX {
			x1[n+1] = silk_int16_MAX
		} else if sum < opus_int32(math.MinInt16) {
			x1[n+1] = math.MinInt16
		} else {
			x1[n+1] = opus_int16(sum)
		}
		if diff > silk_int16_MAX {
			x2[n+1] = silk_int16_MAX
		} else if diff < opus_int32(math.MinInt16) {
			x2[n+1] = math.MinInt16
		} else {
			x2[n+1] = opus_int16(diff)
		}
	}
}
