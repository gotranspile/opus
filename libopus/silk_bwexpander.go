package libopus

import "unsafe"

func silk_bwexpander(ar *int16, d int, chirp_Q16 int32) {
	var (
		i                   int
		chirp_minus_one_Q16 int32 = int32(int(chirp_Q16) - 65536)
	)
	for i = 0; i < d-1; i++ {
		if 16 == 1 {
			*(*int16)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(int16(0))*uintptr(i))) = int16(((int(chirp_Q16) * int(*(*int16)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(int16(0))*uintptr(i))))) >> 1) + ((int(chirp_Q16) * int(*(*int16)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(int16(0))*uintptr(i))))) & 1))
		} else {
			*(*int16)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(int16(0))*uintptr(i))) = int16((((int(chirp_Q16) * int(*(*int16)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(int16(0))*uintptr(i))))) >> (16 - 1)) + 1) >> 1)
		}
		if 16 == 1 {
			chirp_Q16 += int32(((int(chirp_Q16) * int(chirp_minus_one_Q16)) >> 1) + ((int(chirp_Q16) * int(chirp_minus_one_Q16)) & 1))
		} else {
			chirp_Q16 += int32((((int(chirp_Q16) * int(chirp_minus_one_Q16)) >> (16 - 1)) + 1) >> 1)
		}
	}
	if 16 == 1 {
		*(*int16)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(int16(0))*uintptr(d-1))) = int16(((int(chirp_Q16) * int(*(*int16)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(int16(0))*uintptr(d-1))))) >> 1) + ((int(chirp_Q16) * int(*(*int16)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(int16(0))*uintptr(d-1))))) & 1))
	} else {
		*(*int16)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(int16(0))*uintptr(d-1))) = int16((((int(chirp_Q16) * int(*(*int16)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(int16(0))*uintptr(d-1))))) >> (16 - 1)) + 1) >> 1)
	}
}
