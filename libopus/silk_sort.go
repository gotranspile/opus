package libopus

import "unsafe"

func silk_insertion_sort_increasing(a *opus_int32, idx *int64, L int64, K int64) {
	var (
		value opus_int32
		i     int64
		j     int64
	)
	for i = 0; i < K; i++ {
		*(*int64)(unsafe.Add(unsafe.Pointer(idx), unsafe.Sizeof(int64(0))*uintptr(i))) = i
	}
	for i = 1; i < K; i++ {
		value = *(*opus_int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(opus_int32(0))*uintptr(i)))
		for j = i - 1; j >= 0 && value < *(*opus_int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(opus_int32(0))*uintptr(j))); j-- {
			*(*opus_int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(opus_int32(0))*uintptr(j+1))) = *(*opus_int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(opus_int32(0))*uintptr(j)))
			*(*int64)(unsafe.Add(unsafe.Pointer(idx), unsafe.Sizeof(int64(0))*uintptr(j+1))) = *(*int64)(unsafe.Add(unsafe.Pointer(idx), unsafe.Sizeof(int64(0))*uintptr(j)))
		}
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(opus_int32(0))*uintptr(j+1))) = value
		*(*int64)(unsafe.Add(unsafe.Pointer(idx), unsafe.Sizeof(int64(0))*uintptr(j+1))) = i
	}
	for i = K; i < L; i++ {
		value = *(*opus_int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(opus_int32(0))*uintptr(i)))
		if value < *(*opus_int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(opus_int32(0))*uintptr(K-1))) {
			for j = K - 2; j >= 0 && value < *(*opus_int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(opus_int32(0))*uintptr(j))); j-- {
				*(*opus_int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(opus_int32(0))*uintptr(j+1))) = *(*opus_int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(opus_int32(0))*uintptr(j)))
				*(*int64)(unsafe.Add(unsafe.Pointer(idx), unsafe.Sizeof(int64(0))*uintptr(j+1))) = *(*int64)(unsafe.Add(unsafe.Pointer(idx), unsafe.Sizeof(int64(0))*uintptr(j)))
			}
			*(*opus_int32)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(opus_int32(0))*uintptr(j+1))) = value
			*(*int64)(unsafe.Add(unsafe.Pointer(idx), unsafe.Sizeof(int64(0))*uintptr(j+1))) = i
		}
	}
}
func silk_insertion_sort_increasing_all_values_int16(a *opus_int16, L int64) {
	var (
		value int64
		i     int64
		j     int64
	)
	for i = 1; i < L; i++ {
		value = int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(opus_int16(0))*uintptr(i))))
		for j = i - 1; j >= 0 && value < int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(opus_int16(0))*uintptr(j)))); j-- {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(opus_int16(0))*uintptr(j+1))) = *(*opus_int16)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(opus_int16(0))*uintptr(j)))
		}
		*(*opus_int16)(unsafe.Add(unsafe.Pointer(a), unsafe.Sizeof(opus_int16(0))*uintptr(j+1))) = opus_int16(value)
	}
}
