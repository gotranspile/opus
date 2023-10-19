package libopus

import "math"

func silk_log2lin(inLog_Q7 opus_int32) opus_int32 {
	var (
		out     opus_int32
		frac_Q7 opus_int32
	)
	if inLog_Q7 < 0 {
		return 0
	} else if inLog_Q7 >= 3967 {
		return silk_int32_MAX
	}
	out = 1 << (inLog_Q7 >> 7)
	frac_Q7 = inLog_Q7 & math.MaxInt8
	if inLog_Q7 < 2048 {
		out = out + ((out * (frac_Q7 + (((opus_int32(opus_int16(frac_Q7)) * opus_int32(opus_int16(128-frac_Q7))) * opus_int32(int64(opus_int16(-174)))) >> 16))) >> 7)
	} else {
		out = out + (out>>7)*(frac_Q7+(((opus_int32(opus_int16(frac_Q7))*opus_int32(opus_int16(128-frac_Q7)))*opus_int32(int64(opus_int16(-174))))>>16))
	}
	return out
}
