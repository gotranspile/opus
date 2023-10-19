package libopus

import "math"

func silk_stereo_find_predictor(ratio_Q14 *opus_int32, x [0]opus_int16, y [0]opus_int16, mid_res_amp_Q0 [0]opus_int32, length int64, smooth_coef_Q16 int64) opus_int32 {
	var (
		scale     int64
		scale1    int64
		scale2    int64
		nrgx      opus_int32
		nrgy      opus_int32
		corr      opus_int32
		pred_Q13  opus_int32
		pred2_Q10 opus_int32
	)
	silk_sum_sqr_shift(&nrgx, &scale1, &x[0], length)
	silk_sum_sqr_shift(&nrgy, &scale2, &y[0], length)
	scale = silk_max_int(scale1, scale2)
	scale = scale + (scale & 1)
	nrgy = nrgy >> opus_int32(scale-scale2)
	nrgx = nrgx >> opus_int32(scale-scale1)
	nrgx = opus_int32(silk_max_int(int64(nrgx), 1))
	corr = silk_inner_prod_aligned_scale(&x[0], &y[0], scale, length)
	pred_Q13 = silk_DIV32_varQ(corr, nrgx, 13)
	if (-(1 << 14)) > (1 << 14) {
		if pred_Q13 > (-(1 << 14)) {
			pred_Q13 = -(1 << 14)
		} else if pred_Q13 < (1 << 14) {
			pred_Q13 = 1 << 14
		} else {
			pred_Q13 = pred_Q13
		}
	} else if pred_Q13 > (1 << 14) {
		pred_Q13 = 1 << 14
	} else if pred_Q13 < (-(1 << 14)) {
		pred_Q13 = -(1 << 14)
	} else {
		pred_Q13 = pred_Q13
	}
	pred2_Q10 = (pred_Q13 * opus_int32(int64(opus_int16(pred_Q13)))) >> 16
	smooth_coef_Q16 = silk_max_int(smooth_coef_Q16, int64(func() opus_int32 {
		if pred2_Q10 > 0 {
			return pred2_Q10
		}
		return -pred2_Q10
	}()))
	scale = scale >> 1
	mid_res_amp_Q0[0] = (mid_res_amp_Q0[0]) + ((((opus_int32(opus_uint32(silk_SQRT_APPROX(nrgx)) << opus_uint32(scale))) - mid_res_amp_Q0[0]) * opus_int32(int64(opus_int16(smooth_coef_Q16)))) >> 16)
	nrgy = nrgy - (opus_int32(opus_uint32((corr*opus_int32(int64(opus_int16(pred_Q13))))>>16) << (3 + 1)))
	nrgy = nrgy + (opus_int32(opus_uint32((nrgx*opus_int32(int64(opus_int16(pred2_Q10))))>>16) << 6))
	mid_res_amp_Q0[1] = (mid_res_amp_Q0[1]) + ((((opus_int32(opus_uint32(silk_SQRT_APPROX(nrgy)) << opus_uint32(scale))) - mid_res_amp_Q0[1]) * opus_int32(int64(opus_int16(smooth_coef_Q16)))) >> 16)
	*ratio_Q14 = silk_DIV32_varQ(mid_res_amp_Q0[1], func() opus_int32 {
		if (mid_res_amp_Q0[0]) > 1 {
			return mid_res_amp_Q0[0]
		}
		return 1
	}(), 14)
	if 0 > math.MaxInt16 {
		if (*ratio_Q14) > 0 {
			*ratio_Q14 = 0
		} else if (*ratio_Q14) < math.MaxInt16 {
			*ratio_Q14 = math.MaxInt16
		} else {
			*ratio_Q14 = *ratio_Q14
		}
	} else if (*ratio_Q14) > math.MaxInt16 {
		*ratio_Q14 = math.MaxInt16
	} else if (*ratio_Q14) < 0 {
		*ratio_Q14 = 0
	} else {
		*ratio_Q14 = *ratio_Q14
	}
	return pred_Q13
}
