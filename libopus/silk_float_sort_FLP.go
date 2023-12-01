package libopus

import "unsafe"

func silk_insertion_sort_decreasing_FLP(a *float32, idx *int, L int, K int) {
	var (
		value float32
		i     int
		j     int
	)
	for i = 0; i < K; i++ {
		*(*int)(unsafe.Add(unsafe.Pointer(idx), unsafe.Sizeof(int(0))*uintptr(i))) = i
	}
	for i = 1; i < K; i++ {
		value = *(*float32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(float32(0))*uintptr(i)))
		for j = i - 1; j >= 0 && value > *(*float32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(float32(0))*uintptr(j))); j-- {
			*(*float32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(float32(0))*uintptr(j+1))) = *(*float32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(float32(0))*uintptr(j)))
			*(*int)(unsafe.Add(unsafe.Pointer(idx), unsafe.Sizeof(int(0))*uintptr(j+1))) = *(*int)(unsafe.Add(unsafe.Pointer(idx), unsafe.Sizeof(int(0))*uintptr(j)))
		}
		*(*float32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(float32(0))*uintptr(j+1))) = value
		*(*int)(unsafe.Add(unsafe.Pointer(idx), unsafe.Sizeof(int(0))*uintptr(j+1))) = i
	}
	for i = K; i < L; i++ {
		value = *(*float32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(float32(0))*uintptr(i)))
		if value > *(*float32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(float32(0))*uintptr(K-1))) {
			for j = K - 2; j >= 0 && value > *(*float32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(float32(0))*uintptr(j))); j-- {
				*(*float32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(float32(0))*uintptr(j+1))) = *(*float32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(float32(0))*uintptr(j)))
				*(*int)(unsafe.Add(unsafe.Pointer(idx), unsafe.Sizeof(int(0))*uintptr(j+1))) = *(*int)(unsafe.Add(unsafe.Pointer(idx), unsafe.Sizeof(int(0))*uintptr(j)))
			}
			*(*float32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(float32(0))*uintptr(j+1))) = value
			*(*int)(unsafe.Add(unsafe.Pointer(idx), unsafe.Sizeof(int(0))*uintptr(j+1))) = i
		}
	}
}
