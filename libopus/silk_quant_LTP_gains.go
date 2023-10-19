package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"unsafe"
)

func silk_quant_LTP_gains(B_Q14 [20]opus_int16, cbk_index [4]int8, periodicity_index *int8, sum_log_gain_Q7 *opus_int32, pred_gain_dB_Q7 *int64, XX_Q17 [100]opus_int32, xX_Q17 [20]opus_int32, subfr_len int64, nb_subfr int64, arch int64) {
	var (
		j                    int64
		k                    int64
		cbk_size             int64
		temp_idx             [4]int8
		cl_ptr_Q5            *uint8
		cbk_ptr_Q7           *int8
		cbk_gain_ptr_Q7      *uint8
		XX_Q17_ptr           *opus_int32
		xX_Q17_ptr           *opus_int32
		res_nrg_Q15_subfr    opus_int32
		res_nrg_Q15          opus_int32
		rate_dist_Q7_subfr   opus_int32
		rate_dist_Q7         opus_int32
		min_rate_dist_Q7     opus_int32
		sum_log_gain_tmp_Q7  opus_int32
		best_sum_log_gain_Q7 opus_int32
		max_gain_Q7          opus_int32
		gain_Q7              int64
	)
	min_rate_dist_Q7 = silk_int32_MAX
	best_sum_log_gain_Q7 = 0
	for k = 0; k < 3; k++ {
		var gain_safety opus_int32 = (opus_int32(0.4*(1<<7) + 0.5))
		cl_ptr_Q5 = silk_LTP_gain_BITS_Q5_ptrs[k]
		cbk_ptr_Q7 = silk_LTP_vq_ptrs_Q7[k]
		cbk_gain_ptr_Q7 = silk_LTP_vq_gain_ptrs_Q7[k]
		cbk_size = int64(silk_LTP_vq_sizes[k])
		XX_Q17_ptr = &XX_Q17[0]
		xX_Q17_ptr = &xX_Q17[0]
		res_nrg_Q15 = 0
		rate_dist_Q7 = 0
		sum_log_gain_tmp_Q7 = *sum_log_gain_Q7
		for j = 0; j < nb_subfr; j++ {
			max_gain_Q7 = silk_log2lin(((opus_int32((MAX_SUM_LOG_GAIN_DB/6.0)*(1<<7)+0.5))-sum_log_gain_tmp_Q7)+(opus_int32(7*(1<<7)+0.5))) - gain_safety
			_ = arch
			silk_VQ_WMat_EC_c(&temp_idx[j], &res_nrg_Q15_subfr, &rate_dist_Q7_subfr, &gain_Q7, XX_Q17_ptr, xX_Q17_ptr, cbk_ptr_Q7, cbk_gain_ptr_Q7, cl_ptr_Q5, subfr_len, max_gain_Q7, cbk_size)
			if ((opus_uint32(res_nrg_Q15) + opus_uint32(res_nrg_Q15_subfr)) & 0x80000000) != 0 {
				res_nrg_Q15 = silk_int32_MAX
			} else {
				res_nrg_Q15 = res_nrg_Q15 + res_nrg_Q15_subfr
			}
			if ((opus_uint32(rate_dist_Q7) + opus_uint32(rate_dist_Q7_subfr)) & 0x80000000) != 0 {
				rate_dist_Q7 = silk_int32_MAX
			} else {
				rate_dist_Q7 = rate_dist_Q7 + rate_dist_Q7_subfr
			}
			if 0 > (sum_log_gain_tmp_Q7 + silk_lin2log(gain_safety+opus_int32(gain_Q7)) - (opus_int32(7*(1<<7) + 0.5))) {
				sum_log_gain_tmp_Q7 = 0
			} else {
				sum_log_gain_tmp_Q7 = sum_log_gain_tmp_Q7 + silk_lin2log(gain_safety+opus_int32(gain_Q7)) - (opus_int32(7*(1<<7) + 0.5))
			}
			XX_Q17_ptr = (*opus_int32)(unsafe.Add(unsafe.Pointer(XX_Q17_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(LTP_ORDER*LTP_ORDER)))
			xX_Q17_ptr = (*opus_int32)(unsafe.Add(unsafe.Pointer(xX_Q17_ptr), unsafe.Sizeof(opus_int32(0))*uintptr(LTP_ORDER)))
		}
		if rate_dist_Q7 <= min_rate_dist_Q7 {
			min_rate_dist_Q7 = rate_dist_Q7
			*periodicity_index = int8(k)
			libc.MemCpy(unsafe.Pointer(&cbk_index[0]), unsafe.Pointer(&temp_idx[0]), int(nb_subfr*int64(unsafe.Sizeof(int8(0)))))
			best_sum_log_gain_Q7 = sum_log_gain_tmp_Q7
		}
	}
	cbk_ptr_Q7 = silk_LTP_vq_ptrs_Q7[*periodicity_index]
	for j = 0; j < nb_subfr; j++ {
		for k = 0; k < LTP_ORDER; k++ {
			B_Q14[j*LTP_ORDER+k] = opus_int16(opus_int32(opus_uint32(*(*int8)(unsafe.Add(unsafe.Pointer(cbk_ptr_Q7), int64(cbk_index[j])*LTP_ORDER+k))) << 7))
		}
	}
	if nb_subfr == 2 {
		res_nrg_Q15 = res_nrg_Q15 >> 1
	} else {
		res_nrg_Q15 = res_nrg_Q15 >> 2
	}
	*sum_log_gain_Q7 = best_sum_log_gain_Q7
	*pred_gain_dB_Q7 = int64(opus_int32(opus_int16(silk_lin2log(res_nrg_Q15)-(15<<7))) * opus_int32(opus_int16(-3)))
}
