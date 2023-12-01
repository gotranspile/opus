package libopus

import "unsafe"

func silk_LTP_analysis_filter_FLP(LTP_res *float32, x *float32, B [20]float32, pitchL [4]int, invGains [4]float32, subfr_length int, nb_subfr int, pre_length int) {
	var (
		x_ptr       *float32
		x_lag_ptr   *float32
		Btmp        [5]float32
		LTP_res_ptr *float32
		inv_gain    float32
		k           int
		i           int
		j           int
	)
	x_ptr = x
	LTP_res_ptr = LTP_res
	for k = 0; k < nb_subfr; k++ {
		x_lag_ptr = (*float32)(unsafe.Add(unsafe.Pointer(x_ptr), -int(unsafe.Sizeof(float32(0))*uintptr(pitchL[k]))))
		inv_gain = invGains[k]
		for i = 0; i < LTP_ORDER; i++ {
			Btmp[i] = B[k*LTP_ORDER+i]
		}
		for i = 0; i < subfr_length+pre_length; i++ {
			*(*float32)(unsafe.Add(unsafe.Pointer(LTP_res_ptr), unsafe.Sizeof(float32(0))*uintptr(i))) = *(*float32)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(float32(0))*uintptr(i)))
			for j = 0; j < LTP_ORDER; j++ {
				*(*float32)(unsafe.Add(unsafe.Pointer(LTP_res_ptr), unsafe.Sizeof(float32(0))*uintptr(i))) -= Btmp[j] * *(*float32)(unsafe.Add(unsafe.Pointer(x_lag_ptr), unsafe.Sizeof(float32(0))*uintptr(int(LTP_ORDER/2)-j)))
			}
			*(*float32)(unsafe.Add(unsafe.Pointer(LTP_res_ptr), unsafe.Sizeof(float32(0))*uintptr(i))) *= inv_gain
			x_lag_ptr = (*float32)(unsafe.Add(unsafe.Pointer(x_lag_ptr), unsafe.Sizeof(float32(0))*1))
		}
		LTP_res_ptr = (*float32)(unsafe.Add(unsafe.Pointer(LTP_res_ptr), unsafe.Sizeof(float32(0))*uintptr(subfr_length+pre_length)))
		x_ptr = (*float32)(unsafe.Add(unsafe.Pointer(x_ptr), unsafe.Sizeof(float32(0))*uintptr(subfr_length)))
	}
}
