package libopus

import "unsafe"

func silk_inner_prod_aligned_scale(inVec1 *opus_int16, inVec2 *opus_int16, scale int64, len_ int64) opus_int32 {
	var (
		i   int64
		sum opus_int32 = 0
	)
	for i = 0; i < len_; i++ {
		sum = sum + ((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(inVec1), unsafe.Sizeof(opus_int16(0))*uintptr(i)))) * opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(inVec2), unsafe.Sizeof(opus_int16(0))*uintptr(i))))) >> opus_int32(scale))
	}
	return sum
}
