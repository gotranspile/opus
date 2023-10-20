package libopus

import "unsafe"

func silk_NLSF_VQ(err_Q24 []int32, in_Q15 []int16, pCB_Q8 []uint8, pWght_Q9 []int16, K int, LPC_order int) {
	var (
		i             int
		m             int
		diff_Q15      int32
		diffw_Q24     int32
		sum_error_Q24 int32
		pred_Q24      int32
		w_Q9_ptr      *int16
		cb_Q8_ptr     *uint8
	)
	cb_Q8_ptr = &pCB_Q8[0]
	w_Q9_ptr = &pWght_Q9[0]
	for i = 0; i < K; i++ {
		sum_error_Q24 = 0
		pred_Q24 = 0
		for m = LPC_order - 2; m >= 0; m -= 2 {
			diff_Q15 = int32(int(in_Q15[m+1]) - int(int32(int(uint32(int32(*(*uint8)(unsafe.Add(unsafe.Pointer(cb_Q8_ptr), m+1)))))<<7)))
			diffw_Q24 = int32(int(int32(int16(diff_Q15))) * int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(w_Q9_ptr), unsafe.Sizeof(int16(0))*uintptr(m+1))))))
			sum_error_Q24 = int32(int(sum_error_Q24) + (func() int {
				if (int(diffw_Q24) - (int(pred_Q24) >> 1)) > 0 {
					return int(diffw_Q24) - (int(pred_Q24) >> 1)
				}
				return int(int32(-(int(diffw_Q24) - (int(pred_Q24) >> 1))))
			}()))
			pred_Q24 = diffw_Q24
			diff_Q15 = int32(int(in_Q15[m]) - int(int32(int(uint32(int32(*(*uint8)(unsafe.Add(unsafe.Pointer(cb_Q8_ptr), m)))))<<7)))
			diffw_Q24 = int32(int(int32(int16(diff_Q15))) * int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(w_Q9_ptr), unsafe.Sizeof(int16(0))*uintptr(m))))))
			sum_error_Q24 = int32(int(sum_error_Q24) + (func() int {
				if (int(diffw_Q24) - (int(pred_Q24) >> 1)) > 0 {
					return int(diffw_Q24) - (int(pred_Q24) >> 1)
				}
				return int(int32(-(int(diffw_Q24) - (int(pred_Q24) >> 1))))
			}()))
			pred_Q24 = diffw_Q24
		}
		err_Q24[i] = sum_error_Q24
		cb_Q8_ptr = (*uint8)(unsafe.Add(unsafe.Pointer(cb_Q8_ptr), LPC_order))
		w_Q9_ptr = (*int16)(unsafe.Add(unsafe.Pointer(w_Q9_ptr), unsafe.Sizeof(int16(0))*uintptr(LPC_order)))
	}
}
