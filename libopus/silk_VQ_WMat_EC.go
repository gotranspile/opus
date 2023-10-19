package libopus

import "unsafe"

func silk_VQ_WMat_EC_c(ind *int8, res_nrg_Q15 *opus_int32, rate_dist_Q8 *opus_int32, gain_Q7 *int64, XX_Q17 *opus_int32, xX_Q17 *opus_int32, cb_Q7 *int8, cb_gain_Q7 *uint8, cl_Q5 *uint8, subfr_len int64, max_gain_Q7 opus_int32, L int64) {
	var (
		k           int64
		gain_tmp_Q7 int64
		cb_row_Q7   *int8
		neg_xX_Q24  [5]opus_int32
		sum1_Q15    opus_int32
		sum2_Q24    opus_int32
		bits_res_Q8 opus_int32
		bits_tot_Q8 opus_int32
	)
	neg_xX_Q24[0] = -(opus_int32(opus_uint32(*(*opus_int32)(unsafe.Add(unsafe.Pointer(xX_Q17), unsafe.Sizeof(opus_int32(0))*0))) << 7))
	neg_xX_Q24[1] = -(opus_int32(opus_uint32(*(*opus_int32)(unsafe.Add(unsafe.Pointer(xX_Q17), unsafe.Sizeof(opus_int32(0))*1))) << 7))
	neg_xX_Q24[2] = -(opus_int32(opus_uint32(*(*opus_int32)(unsafe.Add(unsafe.Pointer(xX_Q17), unsafe.Sizeof(opus_int32(0))*2))) << 7))
	neg_xX_Q24[3] = -(opus_int32(opus_uint32(*(*opus_int32)(unsafe.Add(unsafe.Pointer(xX_Q17), unsafe.Sizeof(opus_int32(0))*3))) << 7))
	neg_xX_Q24[4] = -(opus_int32(opus_uint32(*(*opus_int32)(unsafe.Add(unsafe.Pointer(xX_Q17), unsafe.Sizeof(opus_int32(0))*4))) << 7))
	*rate_dist_Q8 = silk_int32_MAX
	*res_nrg_Q15 = silk_int32_MAX
	cb_row_Q7 = cb_Q7
	*ind = 0
	for k = 0; k < L; k++ {
		var penalty opus_int32
		gain_tmp_Q7 = int64(*(*uint8)(unsafe.Add(unsafe.Pointer(cb_gain_Q7), k)))
		sum1_Q15 = opus_int32(1.001*(1<<15) + 0.5)
		penalty = opus_int32(opus_uint32(func() int64 {
			if (gain_tmp_Q7 - int64(max_gain_Q7)) > 0 {
				return gain_tmp_Q7 - int64(max_gain_Q7)
			}
			return 0
		}()) << 11)
		sum2_Q24 = (neg_xX_Q24[0]) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(XX_Q17), unsafe.Sizeof(opus_int32(0))*1)))*opus_int32(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 1)))
		sum2_Q24 = sum2_Q24 + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(XX_Q17), unsafe.Sizeof(opus_int32(0))*2)))*opus_int32(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 2)))
		sum2_Q24 = sum2_Q24 + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(XX_Q17), unsafe.Sizeof(opus_int32(0))*3)))*opus_int32(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 3)))
		sum2_Q24 = sum2_Q24 + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(XX_Q17), unsafe.Sizeof(opus_int32(0))*4)))*opus_int32(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 4)))
		sum2_Q24 = opus_int32(opus_uint32(sum2_Q24) << 1)
		sum2_Q24 = sum2_Q24 + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(XX_Q17), unsafe.Sizeof(opus_int32(0))*0)))*opus_int32(*cb_row_Q7)
		sum1_Q15 = sum1_Q15 + ((sum2_Q24 * opus_int32(int64(opus_int16(*cb_row_Q7)))) >> 16)
		sum2_Q24 = (neg_xX_Q24[1]) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(XX_Q17), unsafe.Sizeof(opus_int32(0))*7)))*opus_int32(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 2)))
		sum2_Q24 = sum2_Q24 + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(XX_Q17), unsafe.Sizeof(opus_int32(0))*8)))*opus_int32(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 3)))
		sum2_Q24 = sum2_Q24 + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(XX_Q17), unsafe.Sizeof(opus_int32(0))*9)))*opus_int32(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 4)))
		sum2_Q24 = opus_int32(opus_uint32(sum2_Q24) << 1)
		sum2_Q24 = sum2_Q24 + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(XX_Q17), unsafe.Sizeof(opus_int32(0))*6)))*opus_int32(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 1)))
		sum1_Q15 = sum1_Q15 + ((sum2_Q24 * opus_int32(int64(opus_int16(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 1)))))) >> 16)
		sum2_Q24 = (neg_xX_Q24[2]) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(XX_Q17), unsafe.Sizeof(opus_int32(0))*13)))*opus_int32(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 3)))
		sum2_Q24 = sum2_Q24 + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(XX_Q17), unsafe.Sizeof(opus_int32(0))*14)))*opus_int32(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 4)))
		sum2_Q24 = opus_int32(opus_uint32(sum2_Q24) << 1)
		sum2_Q24 = sum2_Q24 + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(XX_Q17), unsafe.Sizeof(opus_int32(0))*12)))*opus_int32(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 2)))
		sum1_Q15 = sum1_Q15 + ((sum2_Q24 * opus_int32(int64(opus_int16(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 2)))))) >> 16)
		sum2_Q24 = (neg_xX_Q24[3]) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(XX_Q17), unsafe.Sizeof(opus_int32(0))*19)))*opus_int32(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 4)))
		sum2_Q24 = opus_int32(opus_uint32(sum2_Q24) << 1)
		sum2_Q24 = sum2_Q24 + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(XX_Q17), unsafe.Sizeof(opus_int32(0))*18)))*opus_int32(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 3)))
		sum1_Q15 = sum1_Q15 + ((sum2_Q24 * opus_int32(int64(opus_int16(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 3)))))) >> 16)
		sum2_Q24 = opus_int32(opus_uint32(neg_xX_Q24[4]) << 1)
		sum2_Q24 = sum2_Q24 + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(XX_Q17), unsafe.Sizeof(opus_int32(0))*24)))*opus_int32(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 4)))
		sum1_Q15 = sum1_Q15 + ((sum2_Q24 * opus_int32(int64(opus_int16(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 4)))))) >> 16)
		if sum1_Q15 >= 0 {
			bits_res_Q8 = opus_int32(opus_int16(subfr_len)) * opus_int32(opus_int16(silk_lin2log(sum1_Q15+penalty)-(15<<7)))
			bits_tot_Q8 = bits_res_Q8 + (opus_int32(opus_uint32(*(*uint8)(unsafe.Add(unsafe.Pointer(cl_Q5), k))) << (3 - 1)))
			if bits_tot_Q8 <= *rate_dist_Q8 {
				*rate_dist_Q8 = bits_tot_Q8
				*res_nrg_Q15 = sum1_Q15 + penalty
				*ind = int8(k)
				*gain_Q7 = gain_tmp_Q7
			}
		}
		cb_row_Q7 = (*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), LTP_ORDER))
	}
}
