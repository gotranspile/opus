package libopus

import "unsafe"

func silk_insertion_sort_increasing(a *int32, idx *int, L int, K int) {
	var (
		value int32
		i     int
		j     int
	)
	for i = 0; i < K; i++ {
		*(*int)(unsafe.Add(unsafe.Pointer(idx), unsafe.Sizeof(int(0))*uintptr(i))) = i
	}
	for i = 1; i < K; i++ {
		value = *(*int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int32(0))*uintptr(i)))
		for j = i - 1; j >= 0 && int(value) < int(*(*int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int32(0))*uintptr(j)))); j-- {
			*(*int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int32(0))*uintptr(j+1))) = *(*int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int32(0))*uintptr(j)))
			*(*int)(unsafe.Add(unsafe.Pointer(idx), unsafe.Sizeof(int(0))*uintptr(j+1))) = *(*int)(unsafe.Add(unsafe.Pointer(idx), unsafe.Sizeof(int(0))*uintptr(j)))
		}
		*(*int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int32(0))*uintptr(j+1))) = value
		*(*int)(unsafe.Add(unsafe.Pointer(idx), unsafe.Sizeof(int(0))*uintptr(j+1))) = i
	}
	for i = K; i < L; i++ {
		value = *(*int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int32(0))*uintptr(i)))
		if int(value) < int(*(*int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int32(0))*uintptr(K-1)))) {
			for j = K - 2; j >= 0 && int(value) < int(*(*int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int32(0))*uintptr(j)))); j-- {
				*(*int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int32(0))*uintptr(j+1))) = *(*int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int32(0))*uintptr(j)))
				*(*int)(unsafe.Add(unsafe.Pointer(idx), unsafe.Sizeof(int(0))*uintptr(j+1))) = *(*int)(unsafe.Add(unsafe.Pointer(idx), unsafe.Sizeof(int(0))*uintptr(j)))
			}
			*(*int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int32(0))*uintptr(j+1))) = value
			*(*int)(unsafe.Add(unsafe.Pointer(idx), unsafe.Sizeof(int(0))*uintptr(j+1))) = i
		}
	}
}
func silk_insertion_sort_increasing_all_values_int16(a *int16, L int) {
	var (
		value int
		i     int
		j     int
	)
	for i = 1; i < L; i++ {
		value = int(*(*int16)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int16(0))*uintptr(i))))
		for j = i - 1; j >= 0 && value < int(*(*int16)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int16(0))*uintptr(j)))); j-- {
			*(*int16)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int16(0))*uintptr(j+1))) = *(*int16)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int16(0))*uintptr(j)))
		}
		*(*int16)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(int16(0))*uintptr(j+1))) = int16(value)
	}
}
