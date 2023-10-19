package libopus

import "unsafe"

func silk_bwexpander(ar *opus_int16, d int64, chirp_Q16 opus_int32) {
	var (
		i                   int64
		chirp_minus_one_Q16 opus_int32 = chirp_Q16 - 65536
	)
	for i = 0; i < d-1; i++ {
		if 16 == 1 {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(opus_int16(0))*uintptr(i))) = opus_int16(((chirp_Q16 * opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(opus_int16(0))*uintptr(i))))) >> 1) + ((chirp_Q16 * opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(opus_int16(0))*uintptr(i))))) & 1))
		} else {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(opus_int16(0))*uintptr(i))) = opus_int16((((chirp_Q16 * opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(opus_int16(0))*uintptr(i))))) >> (16 - 1)) + 1) >> 1)
		}
		if 16 == 1 {
			chirp_Q16 += ((chirp_Q16 * chirp_minus_one_Q16) >> 1) + ((chirp_Q16 * chirp_minus_one_Q16) & 1)
		} else {
			chirp_Q16 += (((chirp_Q16 * chirp_minus_one_Q16) >> (16 - 1)) + 1) >> 1
		}
	}
	if 16 == 1 {
		*(*opus_int16)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(opus_int16(0))*uintptr(d-1))) = opus_int16(((chirp_Q16 * opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(opus_int16(0))*uintptr(d-1))))) >> 1) + ((chirp_Q16 * opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(opus_int16(0))*uintptr(d-1))))) & 1))
	} else {
		*(*opus_int16)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(opus_int16(0))*uintptr(d-1))) = opus_int16((((chirp_Q16 * opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(ar), unsafe.Sizeof(opus_int16(0))*uintptr(d-1))))) >> (16 - 1)) + 1) >> 1)
	}
}
