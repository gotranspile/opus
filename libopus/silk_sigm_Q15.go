package libopus

import "math"

var sigm_LUT_slope_Q10 [6]opus_int32 = [6]opus_int32{237, 153, 73, 30, 12, 7}
var sigm_LUT_pos_Q15 [6]opus_int32 = [6]opus_int32{16384, 23955, 28861, 31213, 32178, 32548}
var sigm_LUT_neg_Q15 [6]opus_int32 = [6]opus_int32{16384, 8812, 3906, 1554, 589, 219}

func silk_sigm_Q15(in_Q5 int64) int64 {
	var ind int64
	if in_Q5 < 0 {
		in_Q5 = -in_Q5
		if in_Q5 >= 6*32 {
			return 0
		} else {
			ind = in_Q5 >> 5
			return int64(sigm_LUT_neg_Q15[ind] - opus_int32(opus_int16(sigm_LUT_slope_Q10[ind]))*opus_int32(opus_int16(in_Q5&0x1F)))
		}
	} else {
		if in_Q5 >= 6*32 {
			return math.MaxInt16
		} else {
			ind = in_Q5 >> 5
			return int64(sigm_LUT_pos_Q15[ind] + opus_int32(opus_int16(sigm_LUT_slope_Q10[ind]))*opus_int32(opus_int16(in_Q5&0x1F)))
		}
	}
}
