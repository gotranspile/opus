package libopus

import (
	"math"
	"unsafe"
)

func silk_resampler_down2(S *opus_int32, out *opus_int16, in *opus_int16, inLen opus_int32) {
	var (
		k     opus_int32
		len2  opus_int32 = (inLen >> 1)
		in32  opus_int32
		out32 opus_int32
		Y     opus_int32
		X     opus_int32
	)
	for k = 0; k < len2; k++ {
		in32 = opus_int32(opus_uint32(opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(opus_int16(0))*uintptr(k*2))))) << 10)
		Y = in32 - (*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*0)))
		X = Y + ((Y * opus_int32(int64(silk_resampler_down2_1))) >> 16)
		out32 = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*0))) + X
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*0)) = in32 + X
		in32 = opus_int32(opus_uint32(opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(opus_int16(0))*uintptr(k*2+1))))) << 10)
		Y = in32 - (*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*1)))
		X = (Y * opus_int32(int64(silk_resampler_down2_0))) >> 16
		out32 = out32 + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*1)))
		out32 = out32 + X
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*1)) = in32 + X
		if (func() opus_int32 {
			if 11 == 1 {
				return (out32 >> 1) + (out32 & 1)
			}
			return ((out32 >> (11 - 1)) + 1) >> 1
		}()) > silk_int16_MAX {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int16(0))*uintptr(k))) = silk_int16_MAX
		} else if (func() opus_int32 {
			if 11 == 1 {
				return (out32 >> 1) + (out32 & 1)
			}
			return ((out32 >> (11 - 1)) + 1) >> 1
		}()) < opus_int32(math.MinInt16) {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int16(0))*uintptr(k))) = math.MinInt16
		} else if 11 == 1 {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int16(0))*uintptr(k))) = opus_int16((out32 >> 1) + (out32 & 1))
		} else {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int16(0))*uintptr(k))) = opus_int16(((out32 >> (11 - 1)) + 1) >> 1)
		}
	}
}
