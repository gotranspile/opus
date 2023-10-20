package libopus

import (
	"math"
	"unsafe"
)

func silk_VQ_WMat_EC_c(ind *int8, res_nrg_Q15 *int32, rate_dist_Q8 *int32, gain_Q7 *int, XX_Q17 *int32, xX_Q17 *int32, cb_Q7 *int8, cb_gain_Q7 *uint8, cl_Q5 *uint8, subfr_len int, max_gain_Q7 int32, L int) {
	var (
		k           int
		gain_tmp_Q7 int
		cb_row_Q7   *int8
		neg_xX_Q24  [5]int32
		sum1_Q15    int32
		sum2_Q24    int32
		bits_res_Q8 int32
		bits_tot_Q8 int32
	)
	neg_xX_Q24[0] = -(int32(int(uint32(*(*int32)(unsafe.Add(unsafe.Pointer(xX_Q17), unsafe.Sizeof(int32(0))*0)))) << 7))
	neg_xX_Q24[1] = -(int32(int(uint32(*(*int32)(unsafe.Add(unsafe.Pointer(xX_Q17), unsafe.Sizeof(int32(0))*1)))) << 7))
	neg_xX_Q24[2] = -(int32(int(uint32(*(*int32)(unsafe.Add(unsafe.Pointer(xX_Q17), unsafe.Sizeof(int32(0))*2)))) << 7))
	neg_xX_Q24[3] = -(int32(int(uint32(*(*int32)(unsafe.Add(unsafe.Pointer(xX_Q17), unsafe.Sizeof(int32(0))*3)))) << 7))
	neg_xX_Q24[4] = -(int32(int(uint32(*(*int32)(unsafe.Add(unsafe.Pointer(xX_Q17), unsafe.Sizeof(int32(0))*4)))) << 7))
	*rate_dist_Q8 = silk_int32_MAX
	*res_nrg_Q15 = silk_int32_MAX
	cb_row_Q7 = cb_Q7
	*ind = 0
	for k = 0; k < L; k++ {
		var penalty int32
		gain_tmp_Q7 = int(*(*uint8)(unsafe.Add(unsafe.Pointer(cb_gain_Q7), k)))
		sum1_Q15 = int32(math.Floor(1.001*(1<<15) + 0.5))
		penalty = int32(int(uint32(int32(func() int {
			if (gain_tmp_Q7 - int(max_gain_Q7)) > 0 {
				return gain_tmp_Q7 - int(max_gain_Q7)
			}
			return 0
		}()))) << 11)
		sum2_Q24 = int32(int(neg_xX_Q24[0]) + int(*(*int32)(unsafe.Add(unsafe.Pointer(XX_Q17), unsafe.Sizeof(int32(0))*1)))*int(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 1))))
		sum2_Q24 = int32(int(sum2_Q24) + int(*(*int32)(unsafe.Add(unsafe.Pointer(XX_Q17), unsafe.Sizeof(int32(0))*2)))*int(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 2))))
		sum2_Q24 = int32(int(sum2_Q24) + int(*(*int32)(unsafe.Add(unsafe.Pointer(XX_Q17), unsafe.Sizeof(int32(0))*3)))*int(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 3))))
		sum2_Q24 = int32(int(sum2_Q24) + int(*(*int32)(unsafe.Add(unsafe.Pointer(XX_Q17), unsafe.Sizeof(int32(0))*4)))*int(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 4))))
		sum2_Q24 = int32(int(uint32(sum2_Q24)) << 1)
		sum2_Q24 = int32(int(sum2_Q24) + int(*(*int32)(unsafe.Add(unsafe.Pointer(XX_Q17), unsafe.Sizeof(int32(0))*0)))*int(*cb_row_Q7))
		sum1_Q15 = int32(int64(sum1_Q15) + ((int64(sum2_Q24) * int64(int16(*cb_row_Q7))) >> 16))
		sum2_Q24 = int32(int(neg_xX_Q24[1]) + int(*(*int32)(unsafe.Add(unsafe.Pointer(XX_Q17), unsafe.Sizeof(int32(0))*7)))*int(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 2))))
		sum2_Q24 = int32(int(sum2_Q24) + int(*(*int32)(unsafe.Add(unsafe.Pointer(XX_Q17), unsafe.Sizeof(int32(0))*8)))*int(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 3))))
		sum2_Q24 = int32(int(sum2_Q24) + int(*(*int32)(unsafe.Add(unsafe.Pointer(XX_Q17), unsafe.Sizeof(int32(0))*9)))*int(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 4))))
		sum2_Q24 = int32(int(uint32(sum2_Q24)) << 1)
		sum2_Q24 = int32(int(sum2_Q24) + int(*(*int32)(unsafe.Add(unsafe.Pointer(XX_Q17), unsafe.Sizeof(int32(0))*6)))*int(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 1))))
		sum1_Q15 = int32(int64(sum1_Q15) + ((int64(sum2_Q24) * int64(int16(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 1))))) >> 16))
		sum2_Q24 = int32(int(neg_xX_Q24[2]) + int(*(*int32)(unsafe.Add(unsafe.Pointer(XX_Q17), unsafe.Sizeof(int32(0))*13)))*int(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 3))))
		sum2_Q24 = int32(int(sum2_Q24) + int(*(*int32)(unsafe.Add(unsafe.Pointer(XX_Q17), unsafe.Sizeof(int32(0))*14)))*int(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 4))))
		sum2_Q24 = int32(int(uint32(sum2_Q24)) << 1)
		sum2_Q24 = int32(int(sum2_Q24) + int(*(*int32)(unsafe.Add(unsafe.Pointer(XX_Q17), unsafe.Sizeof(int32(0))*12)))*int(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 2))))
		sum1_Q15 = int32(int64(sum1_Q15) + ((int64(sum2_Q24) * int64(int16(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 2))))) >> 16))
		sum2_Q24 = int32(int(neg_xX_Q24[3]) + int(*(*int32)(unsafe.Add(unsafe.Pointer(XX_Q17), unsafe.Sizeof(int32(0))*19)))*int(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 4))))
		sum2_Q24 = int32(int(uint32(sum2_Q24)) << 1)
		sum2_Q24 = int32(int(sum2_Q24) + int(*(*int32)(unsafe.Add(unsafe.Pointer(XX_Q17), unsafe.Sizeof(int32(0))*18)))*int(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 3))))
		sum1_Q15 = int32(int64(sum1_Q15) + ((int64(sum2_Q24) * int64(int16(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 3))))) >> 16))
		sum2_Q24 = int32(int(uint32(neg_xX_Q24[4])) << 1)
		sum2_Q24 = int32(int(sum2_Q24) + int(*(*int32)(unsafe.Add(unsafe.Pointer(XX_Q17), unsafe.Sizeof(int32(0))*24)))*int(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 4))))
		sum1_Q15 = int32(int64(sum1_Q15) + ((int64(sum2_Q24) * int64(int16(*(*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), 4))))) >> 16))
		if int(sum1_Q15) >= 0 {
			bits_res_Q8 = int32(int(int32(int16(subfr_len))) * int(int32(int16(int(silk_lin2log(int32(int(sum1_Q15)+int(penalty))))-(15<<7)))))
			bits_tot_Q8 = int32(int(bits_res_Q8) + int(int32(int(uint32(*(*uint8)(unsafe.Add(unsafe.Pointer(cl_Q5), k))))<<(3-1))))
			if int(bits_tot_Q8) <= int(*rate_dist_Q8) {
				*rate_dist_Q8 = bits_tot_Q8
				*res_nrg_Q15 = int32(int(sum1_Q15) + int(penalty))
				*ind = int8(k)
				*gain_Q7 = gain_tmp_Q7
			}
		}
		cb_row_Q7 = (*int8)(unsafe.Add(unsafe.Pointer(cb_row_Q7), LTP_ORDER))
	}
}
