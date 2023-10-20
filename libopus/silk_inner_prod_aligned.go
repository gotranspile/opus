package libopus

import "unsafe"

func silk_inner_prod_aligned_scale(inVec1 *int16, inVec2 *int16, scale int, len_ int) int32 {
	var (
		i   int
		sum int32 = 0
	)
	for i = 0; i < len_; i++ {
		sum = int32(int(sum) + ((int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(inVec1), unsafe.Sizeof(int16(0))*uintptr(i))))) * int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(inVec2), unsafe.Sizeof(int16(0))*uintptr(i)))))) >> scale))
	}
	return sum
}
