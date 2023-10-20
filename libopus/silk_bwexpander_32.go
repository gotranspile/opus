package libopus

import "unsafe"

func silk_bwexpander_32(ar *int32, d int, chirp_Q16 int32) {
	var (
		i                   int
		chirp_minus_one_Q16 int32 = int32(int(chirp_Q16) - 65536)
	)
	for i = 0; i < d-1; i++ {
		*(*int32)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(int32(0))*uintptr(i))) = int32((int64(chirp_Q16) * int64(*(*int32)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(int32(0))*uintptr(i))))) >> 16)
		if 16 == 1 {
			chirp_Q16 += int32(((int(chirp_Q16) * int(chirp_minus_one_Q16)) >> 1) + ((int(chirp_Q16) * int(chirp_minus_one_Q16)) & 1))
		} else {
			chirp_Q16 += int32((((int(chirp_Q16) * int(chirp_minus_one_Q16)) >> (16 - 1)) + 1) >> 1)
		}
	}
	*(*int32)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(int32(0))*uintptr(d-1))) = int32((int64(chirp_Q16) * int64(*(*int32)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(int32(0))*uintptr(d-1))))) >> 16)
}
