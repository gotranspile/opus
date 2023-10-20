package libopus

import (
	"math"
	"unsafe"
)

func silk_resampler_private_up2_HQ(S *int32, out *int16, in *int16, len_ int32) {
	var (
		k       int32
		in32    int32
		out32_1 int32
		out32_2 int32
		Y       int32
		X       int32
	)
	for k = 0; int(k) < int(len_); k++ {
		in32 = int32(int(uint32(int32(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(k)))))) << 10)
		Y = int32(int(in32) - int(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0))))
		X = int32((int64(Y) * int64(silk_resampler_up2_hq_0[0])) >> 16)
		out32_1 = int32(int(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0))) + int(X))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0)) = int32(int(in32) + int(X))
		Y = int32(int(out32_1) - int(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1))))
		X = int32((int64(Y) * int64(silk_resampler_up2_hq_0[1])) >> 16)
		out32_2 = int32(int(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1))) + int(X))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1)) = int32(int(out32_1) + int(X))
		Y = int32(int(out32_2) - int(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*2))))
		X = int32(int64(Y) + ((int64(Y) * int64(silk_resampler_up2_hq_0[2])) >> 16))
		out32_1 = int32(int(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*2))) + int(X))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*2)) = int32(int(out32_2) + int(X))
		if (func() int {
			if 10 == 1 {
				return (int(out32_1) >> 1) + (int(out32_1) & 1)
			}
			return ((int(out32_1) >> (10 - 1)) + 1) >> 1
		}()) > silk_int16_MAX {
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(int(k)*2))) = silk_int16_MAX
		} else if (func() int {
			if 10 == 1 {
				return (int(out32_1) >> 1) + (int(out32_1) & 1)
			}
			return ((int(out32_1) >> (10 - 1)) + 1) >> 1
		}()) < int(math.MinInt16) {
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(int(k)*2))) = math.MinInt16
		} else if 10 == 1 {
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(int(k)*2))) = int16((int(out32_1) >> 1) + (int(out32_1) & 1))
		} else {
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(int(k)*2))) = int16(((int(out32_1) >> (10 - 1)) + 1) >> 1)
		}
		Y = int32(int(in32) - int(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*3))))
		X = int32((int64(Y) * int64(silk_resampler_up2_hq_1[0])) >> 16)
		out32_1 = int32(int(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*3))) + int(X))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*3)) = int32(int(in32) + int(X))
		Y = int32(int(out32_1) - int(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*4))))
		X = int32((int64(Y) * int64(silk_resampler_up2_hq_1[1])) >> 16)
		out32_2 = int32(int(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*4))) + int(X))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*4)) = int32(int(out32_1) + int(X))
		Y = int32(int(out32_2) - int(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*5))))
		X = int32(int64(Y) + ((int64(Y) * int64(silk_resampler_up2_hq_1[2])) >> 16))
		out32_1 = int32(int(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*5))) + int(X))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*5)) = int32(int(out32_2) + int(X))
		if (func() int {
			if 10 == 1 {
				return (int(out32_1) >> 1) + (int(out32_1) & 1)
			}
			return ((int(out32_1) >> (10 - 1)) + 1) >> 1
		}()) > silk_int16_MAX {
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(int(k)*2+1))) = silk_int16_MAX
		} else if (func() int {
			if 10 == 1 {
				return (int(out32_1) >> 1) + (int(out32_1) & 1)
			}
			return ((int(out32_1) >> (10 - 1)) + 1) >> 1
		}()) < int(math.MinInt16) {
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(int(k)*2+1))) = math.MinInt16
		} else if 10 == 1 {
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(int(k)*2+1))) = int16((int(out32_1) >> 1) + (int(out32_1) & 1))
		} else {
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(int(k)*2+1))) = int16(((int(out32_1) >> (10 - 1)) + 1) >> 1)
		}
	}
}
func silk_resampler_private_up2_HQ_wrapper(SS unsafe.Pointer, out *int16, in *int16, len_ int32) {
	var S *silk_resampler_state_struct = (*silk_resampler_state_struct)(SS)
	silk_resampler_private_up2_HQ(&S.SIIR[0], out, in, len_)
}
