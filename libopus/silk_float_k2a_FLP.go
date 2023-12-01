package libopus

import "unsafe"

func silk_k2a_FLP(A *float32, rc *float32, order int32) {
	var (
		k    int
		n    int
		rck  float32
		tmp1 float32
		tmp2 float32
	)
	for k = 0; k < int(order); k++ {
		rck = *(*float32)(unsafe.Add(unsafe.Pointer(rc), unsafe.Sizeof(float32(0))*uintptr(k)))
		for n = 0; n < (k+1)>>1; n++ {
			tmp1 = *(*float32)(unsafe.Add(unsafe.Pointer(A), unsafe.Sizeof(float32(0))*uintptr(n)))
			tmp2 = *(*float32)(unsafe.Add(unsafe.Pointer(A), unsafe.Sizeof(float32(0))*uintptr(k-n-1)))
			*(*float32)(unsafe.Add(unsafe.Pointer(A), unsafe.Sizeof(float32(0))*uintptr(n))) = tmp1 + tmp2*rck
			*(*float32)(unsafe.Add(unsafe.Pointer(A), unsafe.Sizeof(float32(0))*uintptr(k-n-1))) = tmp2 + tmp1*rck
		}
		*(*float32)(unsafe.Add(unsafe.Pointer(A), unsafe.Sizeof(float32(0))*uintptr(k))) = -rck
	}
}
