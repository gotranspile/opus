package libopus

import "unsafe"

func silk_find_LTP_FLP(XX [100]float32, xX [20]float32, r_ptr []float32, lag [4]int, subfr_length int, nb_subfr int) {
	var (
		k       int
		xX_ptr  *float32
		XX_ptr  *float32
		lag_ptr *float32
		xx      float32
		temp    float32
	)
	xX_ptr = &xX[0]
	XX_ptr = &XX[0]
	for k = 0; k < nb_subfr; k++ {
		lag_ptr = (*float32)(unsafe.Add(unsafe.Pointer(&r_ptr[0]), -int(unsafe.Sizeof(float32(0))*uintptr(lag[k]+int(LTP_ORDER/2)))))
		silk_corrMatrix_FLP(lag_ptr, subfr_length, LTP_ORDER, XX_ptr)
		silk_corrVector_FLP(lag_ptr, &r_ptr[0], subfr_length, LTP_ORDER, xX_ptr)
		xx = float32(silk_energy_FLP(r_ptr, subfr_length+LTP_ORDER))
		temp = 1.0 / (func() float32 {
			if xx > (LTP_CORR_INV_MAX*0.5*(*(*float32)(unsafe.Add(unsafe.Pointer(XX_ptr), unsafe.Sizeof(float32(0))*0))+*(*float32)(unsafe.Add(unsafe.Pointer(XX_ptr), unsafe.Sizeof(float32(0))*24))) + 1.0) {
				return xx
			}
			return LTP_CORR_INV_MAX*0.5*(*(*float32)(unsafe.Add(unsafe.Pointer(XX_ptr), unsafe.Sizeof(float32(0))*0))+*(*float32)(unsafe.Add(unsafe.Pointer(XX_ptr), unsafe.Sizeof(float32(0))*24))) + 1.0
		}())
		silk_scale_vector_FLP(XX_ptr, temp, int(LTP_ORDER*LTP_ORDER))
		silk_scale_vector_FLP(xX_ptr, temp, LTP_ORDER)
		r_ptr += []float32(subfr_length)
		XX_ptr = (*float32)(unsafe.Add(unsafe.Pointer(XX_ptr), unsafe.Sizeof(float32(0))*uintptr(int(LTP_ORDER*LTP_ORDER))))
		xX_ptr = (*float32)(unsafe.Add(unsafe.Pointer(xX_ptr), unsafe.Sizeof(float32(0))*uintptr(LTP_ORDER)))
	}
}
