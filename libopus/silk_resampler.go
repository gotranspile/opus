package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

const USE_silk_resampler_copy = 0
const USE_silk_resampler_private_up2_HQ_wrapper = 1
const USE_silk_resampler_private_IIR_FIR = 2
const USE_silk_resampler_private_down_FIR = 3

var delay_matrix_enc [5][3]int8 = [5][3]int8{{6, 0, 3}, {0, 7, 3}, {0, 1, 10}, {0, 2, 6}, {18, 10, 12}}
var delay_matrix_dec [3][5]int8 = [3][5]int8{{4, 0, 2, 0, 0}, {0, 9, 4, 7, 4}, {0, 3, 12, 7, 7}}

func silk_resampler_init(S *silk_resampler_state_struct, Fs_Hz_in opus_int32, Fs_Hz_out opus_int32, forEnc int64) int64 {
	var up2x int64
	*S = silk_resampler_state_struct{}
	if forEnc != 0 {
		if Fs_Hz_in != 8000 && Fs_Hz_in != 12000 && Fs_Hz_in != 16000 && Fs_Hz_in != 24000 && Fs_Hz_in != 48000 || Fs_Hz_out != 8000 && Fs_Hz_out != 12000 && Fs_Hz_out != 16000 {
			return -1
		}
		S.InputDelay = int64(delay_matrix_enc[(((Fs_Hz_in>>12)-opus_int32(libc.BoolToInt(Fs_Hz_in > 16000)))>>opus_int32(libc.BoolToInt(Fs_Hz_in > 24000)))-1][(((Fs_Hz_out>>12)-opus_int32(libc.BoolToInt(Fs_Hz_out > 16000)))>>opus_int32(libc.BoolToInt(Fs_Hz_out > 24000)))-1])
	} else {
		if Fs_Hz_in != 8000 && Fs_Hz_in != 12000 && Fs_Hz_in != 16000 || Fs_Hz_out != 8000 && Fs_Hz_out != 12000 && Fs_Hz_out != 16000 && Fs_Hz_out != 24000 && Fs_Hz_out != 48000 {
			return -1
		}
		S.InputDelay = int64(delay_matrix_dec[(((Fs_Hz_in>>12)-opus_int32(libc.BoolToInt(Fs_Hz_in > 16000)))>>opus_int32(libc.BoolToInt(Fs_Hz_in > 24000)))-1][(((Fs_Hz_out>>12)-opus_int32(libc.BoolToInt(Fs_Hz_out > 16000)))>>opus_int32(libc.BoolToInt(Fs_Hz_out > 24000)))-1])
	}
	S.Fs_in_kHz = int64(Fs_Hz_in / 1000)
	S.Fs_out_kHz = int64(Fs_Hz_out / 1000)
	S.BatchSize = S.Fs_in_kHz * RESAMPLER_MAX_BATCH_SIZE_MS
	up2x = 0
	if Fs_Hz_out > Fs_Hz_in {
		if Fs_Hz_out == (Fs_Hz_in * 2) {
			S.Resampler_function = 1
		} else {
			S.Resampler_function = 2
			up2x = 1
		}
	} else if Fs_Hz_out < Fs_Hz_in {
		S.Resampler_function = 3
		if (Fs_Hz_out * 4) == (Fs_Hz_in * 3) {
			S.FIR_Fracs = 3
			S.FIR_Order = RESAMPLER_DOWN_ORDER_FIR0
			S.Coefs = &silk_Resampler_3_4_COEFS[0]
		} else if (Fs_Hz_out * 3) == (Fs_Hz_in * 2) {
			S.FIR_Fracs = 2
			S.FIR_Order = RESAMPLER_DOWN_ORDER_FIR0
			S.Coefs = &silk_Resampler_2_3_COEFS[0]
		} else if (Fs_Hz_out * 2) == Fs_Hz_in {
			S.FIR_Fracs = 1
			S.FIR_Order = RESAMPLER_DOWN_ORDER_FIR1
			S.Coefs = &silk_Resampler_1_2_COEFS[0]
		} else if (Fs_Hz_out * 3) == Fs_Hz_in {
			S.FIR_Fracs = 1
			S.FIR_Order = RESAMPLER_DOWN_ORDER_FIR2
			S.Coefs = &silk_Resampler_1_3_COEFS[0]
		} else if (Fs_Hz_out * 4) == Fs_Hz_in {
			S.FIR_Fracs = 1
			S.FIR_Order = RESAMPLER_DOWN_ORDER_FIR2
			S.Coefs = &silk_Resampler_1_4_COEFS[0]
		} else if (Fs_Hz_out * 6) == Fs_Hz_in {
			S.FIR_Fracs = 1
			S.FIR_Order = RESAMPLER_DOWN_ORDER_FIR2
			S.Coefs = &silk_Resampler_1_6_COEFS[0]
		} else {
			return -1
		}
	} else {
		S.Resampler_function = 0
	}
	S.InvRatio_Q16 = opus_int32(opus_uint32((opus_int32(opus_uint32(Fs_Hz_in)<<opus_uint32(up2x+14)))/Fs_Hz_out) << 2)
	for (opus_int32((int64(S.InvRatio_Q16) * int64(Fs_Hz_out)) >> 16)) < (opus_int32(opus_uint32(Fs_Hz_in) << opus_uint32(up2x))) {
		S.InvRatio_Q16++
	}
	return 0
}
func silk_resampler(S *silk_resampler_state_struct, out [0]opus_int16, in [0]opus_int16, inLen opus_int32) int64 {
	var nSamples int64
	nSamples = S.Fs_in_kHz - S.InputDelay
	libc.MemCpy(unsafe.Pointer(&S.DelayBuf[S.InputDelay]), unsafe.Pointer(&in[0]), int(nSamples*int64(unsafe.Sizeof(opus_int16(0)))))
	switch S.Resampler_function {
	case 1:
		silk_resampler_private_up2_HQ_wrapper(unsafe.Pointer(S), &out[0], &S.DelayBuf[0], opus_int32(S.Fs_in_kHz))
		silk_resampler_private_up2_HQ_wrapper(unsafe.Pointer(S), &out[S.Fs_out_kHz], &in[nSamples], inLen-opus_int32(S.Fs_in_kHz))
	case 2:
		silk_resampler_private_IIR_FIR(unsafe.Pointer(S), out, S.DelayBuf[:], opus_int32(S.Fs_in_kHz))
		silk_resampler_private_IIR_FIR(unsafe.Pointer(S), [0]opus_int16(&out[S.Fs_out_kHz]), [0]opus_int16(&in[nSamples]), inLen-opus_int32(S.Fs_in_kHz))
	case 3:
		silk_resampler_private_down_FIR(unsafe.Pointer(S), out, S.DelayBuf[:], opus_int32(S.Fs_in_kHz))
		silk_resampler_private_down_FIR(unsafe.Pointer(S), [0]opus_int16(&out[S.Fs_out_kHz]), [0]opus_int16(&in[nSamples]), inLen-opus_int32(S.Fs_in_kHz))
	default:
		libc.MemCpy(unsafe.Pointer(&out[0]), unsafe.Pointer(&S.DelayBuf[0]), int(S.Fs_in_kHz*int64(unsafe.Sizeof(opus_int16(0)))))
		libc.MemCpy(unsafe.Pointer(&out[S.Fs_out_kHz]), unsafe.Pointer(&in[nSamples]), int((inLen-opus_int32(S.Fs_in_kHz))*opus_int32(unsafe.Sizeof(opus_int16(0)))))
	}
	libc.MemCpy(unsafe.Pointer(&S.DelayBuf[0]), unsafe.Pointer(&in[inLen-opus_int32(S.InputDelay)]), int(S.InputDelay*int64(unsafe.Sizeof(opus_int16(0)))))
	return 0
}
