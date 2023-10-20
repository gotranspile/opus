package libopus

import (
	"math"
	"unsafe"
)

func silk_resampler_down2(S *int32, out *int16, in *int16, inLen int32) {
	var (
		k     int32
		len2  int32 = int32(int(inLen) >> 1)
		in32  int32
		out32 int32
		Y     int32
		X     int32
	)
	for k = 0; int(k) < int(len2); k++ {
		in32 = int32(int(uint32(int32(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(int(k)*2)))))) << 10)
		Y = int32(int(in32) - int(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0))))
		X = int32(int64(Y) + ((int64(Y) * int64(silk_resampler_down2_1)) >> 16))
		out32 = int32(int(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0))) + int(X))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0)) = int32(int(in32) + int(X))
		in32 = int32(int(uint32(int32(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(int(k)*2+1)))))) << 10)
		Y = int32(int(in32) - int(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1))))
		X = int32((int64(Y) * int64(silk_resampler_down2_0)) >> 16)
		out32 = int32(int(out32) + int(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1))))
		out32 = int32(int(out32) + int(X))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1)) = int32(int(in32) + int(X))
		if (func() int {
			if 11 == 1 {
				return (int(out32) >> 1) + (int(out32) & 1)
			}
			return ((int(out32) >> (11 - 1)) + 1) >> 1
		}()) > silk_int16_MAX {
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(k))) = silk_int16_MAX
		} else if (func() int {
			if 11 == 1 {
				return (int(out32) >> 1) + (int(out32) & 1)
			}
			return ((int(out32) >> (11 - 1)) + 1) >> 1
		}()) < int(math.MinInt16) {
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(k))) = math.MinInt16
		} else if 11 == 1 {
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(k))) = int16((int(out32) >> 1) + (int(out32) & 1))
		} else {
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(k))) = int16(((int(out32) >> (11 - 1)) + 1) >> 1)
		}
	}
}
