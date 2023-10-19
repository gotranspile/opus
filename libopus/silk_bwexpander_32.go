package libopus

import "unsafe"

func silk_bwexpander_32(ar *opus_int32, d int64, chirp_Q16 opus_int32) {
	var (
		i                   int64
		chirp_minus_one_Q16 opus_int32 = chirp_Q16 - 65536
	)
	for i = 0; i < d-1; i++ {
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(opus_int32(0))*uintptr(i))) = opus_int32((int64(chirp_Q16) * int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(opus_int32(0))*uintptr(i))))) >> 16)
		if 16 == 1 {
			chirp_Q16 += ((chirp_Q16 * chirp_minus_one_Q16) >> 1) + ((chirp_Q16 * chirp_minus_one_Q16) & 1)
		} else {
			chirp_Q16 += (((chirp_Q16 * chirp_minus_one_Q16) >> (16 - 1)) + 1) >> 1
		}
	}
	*(*opus_int32)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(opus_int32(0))*uintptr(d-1))) = opus_int32((int64(chirp_Q16) * int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(opus_int32(0))*uintptr(d-1))))) >> 16)
}
