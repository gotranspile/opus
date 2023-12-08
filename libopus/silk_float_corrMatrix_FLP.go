package libopus

import "unsafe"

func silk_corrVector_FLP(x *float32, t *float32, L int, Order int, Xt *float32) {
	var (
		lag  int
		ptr1 *float32
	)
	ptr1 = (*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr(Order-1)))
	for lag = 0; lag < Order; lag++ {
		*(*float32)(unsafe.Add(unsafe.Pointer(Xt), unsafe.Sizeof(float32(0))*uintptr(lag))) = float32(silk_inner_product_FLP([]float32(ptr1), []float32(t), L))
		ptr1 = (*float32)(unsafe.Add(unsafe.Pointer(ptr1), -int(unsafe.Sizeof(float32(0))*1)))
	}
}
func silk_corrMatrix_FLP(x *float32, L int, Order int, XX *float32) {
	var (
		j      int
		lag    int
		energy float64
		ptr1   *float32
		ptr2   *float32
	)
	ptr1 = (*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr(Order-1)))
	energy = silk_energy_FLP([]float32(ptr1), L)
	*((*float32)(unsafe.Add(unsafe.Pointer(XX), unsafe.Sizeof(float32(0))*uintptr(Order*0+0)))) = float32(energy)
	for j = 1; j < Order; j++ {
		energy += float64(*(*float32)(unsafe.Add(unsafe.Pointer(ptr1), -int(unsafe.Sizeof(float32(0))*uintptr(j))))**(*float32)(unsafe.Add(unsafe.Pointer(ptr1), -int(unsafe.Sizeof(float32(0))*uintptr(j)))) - *(*float32)(unsafe.Add(unsafe.Pointer(ptr1), unsafe.Sizeof(float32(0))*uintptr(L-j)))**(*float32)(unsafe.Add(unsafe.Pointer(ptr1), unsafe.Sizeof(float32(0))*uintptr(L-j))))
		*((*float32)(unsafe.Add(unsafe.Pointer(XX), unsafe.Sizeof(float32(0))*uintptr(j*Order+j)))) = float32(energy)
	}
	ptr2 = (*float32)(unsafe.Add(unsafe.Pointer(x), unsafe.Sizeof(float32(0))*uintptr(Order-2)))
	for lag = 1; lag < Order; lag++ {
		energy = silk_inner_product_FLP([]float32(ptr1), []float32(ptr2), L)
		*((*float32)(unsafe.Add(unsafe.Pointer(XX), unsafe.Sizeof(float32(0))*uintptr(lag*Order+0)))) = float32(energy)
		*((*float32)(unsafe.Add(unsafe.Pointer(XX), unsafe.Sizeof(float32(0))*uintptr(Order*0+lag)))) = float32(energy)
		for j = 1; j < (Order - lag); j++ {
			energy += float64(*(*float32)(unsafe.Add(unsafe.Pointer(ptr1), -int(unsafe.Sizeof(float32(0))*uintptr(j))))**(*float32)(unsafe.Add(unsafe.Pointer(ptr2), -int(unsafe.Sizeof(float32(0))*uintptr(j)))) - *(*float32)(unsafe.Add(unsafe.Pointer(ptr1), unsafe.Sizeof(float32(0))*uintptr(L-j)))**(*float32)(unsafe.Add(unsafe.Pointer(ptr2), unsafe.Sizeof(float32(0))*uintptr(L-j))))
			*((*float32)(unsafe.Add(unsafe.Pointer(XX), unsafe.Sizeof(float32(0))*uintptr((lag+j)*Order+j)))) = float32(energy)
			*((*float32)(unsafe.Add(unsafe.Pointer(XX), unsafe.Sizeof(float32(0))*uintptr(j*Order+(lag+j))))) = float32(energy)
		}
		ptr2 = (*float32)(unsafe.Add(unsafe.Pointer(ptr2), -int(unsafe.Sizeof(float32(0))*1)))
	}
}
