package libopus

import "unsafe"

func silk_NLSF_VQ(err_Q24 [0]opus_int32, in_Q15 [0]opus_int16, pCB_Q8 [0]uint8, pWght_Q9 [0]opus_int16, K int64, LPC_order int64) {
	var (
		i             int64
		m             int64
		diff_Q15      opus_int32
		diffw_Q24     opus_int32
		sum_error_Q24 opus_int32
		pred_Q24      opus_int32
		w_Q9_ptr      *opus_int16
		cb_Q8_ptr     *uint8
	)
	cb_Q8_ptr = &pCB_Q8[0]
	w_Q9_ptr = &pWght_Q9[0]
	for i = 0; i < K; i++ {
		sum_error_Q24 = 0
		pred_Q24 = 0
		for m = LPC_order - 2; m >= 0; m -= 2 {
			diff_Q15 = opus_int32(in_Q15[m+1]) - (opus_int32(opus_uint32(opus_int32(*(*uint8)(unsafe.Add(unsafe.Pointer(cb_Q8_ptr), m+1)))) << 7))
			diffw_Q24 = opus_int32(opus_int16(diff_Q15)) * opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(w_Q9_ptr), unsafe.Sizeof(opus_int16(0))*uintptr(m+1))))
			sum_error_Q24 = sum_error_Q24 + (func() opus_int32 {
				if (diffw_Q24 - (pred_Q24 >> 1)) > 0 {
					return diffw_Q24 - (pred_Q24 >> 1)
				}
				return -(diffw_Q24 - (pred_Q24 >> 1))
			}())
			pred_Q24 = diffw_Q24
			diff_Q15 = opus_int32(in_Q15[m]) - (opus_int32(opus_uint32(opus_int32(*(*uint8)(unsafe.Add(unsafe.Pointer(cb_Q8_ptr), m)))) << 7))
			diffw_Q24 = opus_int32(opus_int16(diff_Q15)) * opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(w_Q9_ptr), unsafe.Sizeof(opus_int16(0))*uintptr(m))))
			sum_error_Q24 = sum_error_Q24 + (func() opus_int32 {
				if (diffw_Q24 - (pred_Q24 >> 1)) > 0 {
					return diffw_Q24 - (pred_Q24 >> 1)
				}
				return -(diffw_Q24 - (pred_Q24 >> 1))
			}())
			pred_Q24 = diffw_Q24
		}
		err_Q24[i] = sum_error_Q24
		cb_Q8_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(cb_Q8_ptr), LPC_order))
		w_Q9_ptr = (*opus_int16)(unsafe.Add(unsafe.Pointer(w_Q9_ptr), unsafe.Sizeof(opus_int16(0))*uintptr(LPC_order)))
	}
}
