package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func silk_LP_interpolate_filter_taps(B_Q28 [3]opus_int32, A_Q28 [2]opus_int32, ind int64, fac_Q16 opus_int32) {
	var (
		nb int64
		na int64
	)
	if ind < TRANSITION_INT_NUM-1 {
		if fac_Q16 > 0 {
			if fac_Q16 < 32768 {
				for nb = 0; nb < TRANSITION_NB; nb++ {
					B_Q28[nb] = (silk_Transition_LP_B_Q28[ind][nb]) + (((silk_Transition_LP_B_Q28[ind+1][nb] - silk_Transition_LP_B_Q28[ind][nb]) * opus_int32(int64(opus_int16(fac_Q16)))) >> 16)
				}
				for na = 0; na < TRANSITION_NA; na++ {
					A_Q28[na] = (silk_Transition_LP_A_Q28[ind][na]) + (((silk_Transition_LP_A_Q28[ind+1][na] - silk_Transition_LP_A_Q28[ind][na]) * opus_int32(int64(opus_int16(fac_Q16)))) >> 16)
				}
			} else {
				for nb = 0; nb < TRANSITION_NB; nb++ {
					B_Q28[nb] = (silk_Transition_LP_B_Q28[ind+1][nb]) + (((silk_Transition_LP_B_Q28[ind+1][nb] - silk_Transition_LP_B_Q28[ind][nb]) * opus_int32(int64(opus_int16(fac_Q16-(1<<16))))) >> 16)
				}
				for na = 0; na < TRANSITION_NA; na++ {
					A_Q28[na] = (silk_Transition_LP_A_Q28[ind+1][na]) + (((silk_Transition_LP_A_Q28[ind+1][na] - silk_Transition_LP_A_Q28[ind][na]) * opus_int32(int64(opus_int16(fac_Q16-(1<<16))))) >> 16)
				}
			}
		} else {
			libc.MemCpy(unsafe.Pointer(&B_Q28[0]), unsafe.Pointer(&(silk_Transition_LP_B_Q28[ind])[0]), int(TRANSITION_NB*unsafe.Sizeof(opus_int32(0))))
			libc.MemCpy(unsafe.Pointer(&A_Q28[0]), unsafe.Pointer(&(silk_Transition_LP_A_Q28[ind])[0]), int(TRANSITION_NA*unsafe.Sizeof(opus_int32(0))))
		}
	} else {
		libc.MemCpy(unsafe.Pointer(&B_Q28[0]), unsafe.Pointer(&(silk_Transition_LP_B_Q28[TRANSITION_INT_NUM-1])[0]), int(TRANSITION_NB*unsafe.Sizeof(opus_int32(0))))
		libc.MemCpy(unsafe.Pointer(&A_Q28[0]), unsafe.Pointer(&(silk_Transition_LP_A_Q28[TRANSITION_INT_NUM-1])[0]), int(TRANSITION_NA*unsafe.Sizeof(opus_int32(0))))
	}
}
func silk_LP_variable_cutoff(psLP *silk_LP_state, frame *opus_int16, frame_length int64) {
	var (
		B_Q28   [3]opus_int32
		A_Q28   [2]opus_int32
		fac_Q16 opus_int32 = 0
		ind     int64      = 0
	)
	if psLP.Mode != 0 {
		fac_Q16 = opus_int32(opus_uint32(opus_int32(TRANSITION_TIME_MS/(SUB_FRAME_LENGTH_MS*MAX_NB_SUBFR))-psLP.Transition_frame_no) << (16 - 6))
		ind = int64(fac_Q16 >> 16)
		fac_Q16 -= opus_int32(opus_uint32(ind) << 16)
		silk_LP_interpolate_filter_taps(B_Q28, A_Q28, ind, fac_Q16)
		if 0 > (TRANSITION_TIME_MS / (SUB_FRAME_LENGTH_MS * MAX_NB_SUBFR)) {
			if (psLP.Transition_frame_no + opus_int32(psLP.Mode)) > 0 {
				psLP.Transition_frame_no = 0
			} else if (psLP.Transition_frame_no + opus_int32(psLP.Mode)) < opus_int32(TRANSITION_TIME_MS/(SUB_FRAME_LENGTH_MS*MAX_NB_SUBFR)) {
				psLP.Transition_frame_no = opus_int32(TRANSITION_TIME_MS / (SUB_FRAME_LENGTH_MS * MAX_NB_SUBFR))
			} else {
				psLP.Transition_frame_no = psLP.Transition_frame_no + opus_int32(psLP.Mode)
			}
		} else if (psLP.Transition_frame_no + opus_int32(psLP.Mode)) > opus_int32(TRANSITION_TIME_MS/(SUB_FRAME_LENGTH_MS*MAX_NB_SUBFR)) {
			psLP.Transition_frame_no = opus_int32(TRANSITION_TIME_MS / (SUB_FRAME_LENGTH_MS * MAX_NB_SUBFR))
		} else if (psLP.Transition_frame_no + opus_int32(psLP.Mode)) < 0 {
			psLP.Transition_frame_no = 0
		} else {
			psLP.Transition_frame_no = psLP.Transition_frame_no + opus_int32(psLP.Mode)
		}
		silk_biquad_alt_stride1(frame, &B_Q28[0], &A_Q28[0], &psLP.In_LP_State[0], frame, opus_int32(frame_length))
	}
}
