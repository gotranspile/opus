package libopus

import (
	"math"
	"unsafe"
)

func silk_resampler_private_up2_HQ(S *opus_int32, out *opus_int16, in *opus_int16, len_ opus_int32) {
	var (
		k       opus_int32
		in32    opus_int32
		out32_1 opus_int32
		out32_2 opus_int32
		Y       opus_int32
		X       opus_int32
	)
	for k = 0; k < len_; k++ {
		in32 = opus_int32(opus_uint32(opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(opus_int16(0))*uintptr(k))))) << 10)
		Y = in32 - (*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*0)))
		X = (Y * opus_int32(int64(silk_resampler_up2_hq_0[0]))) >> 16
		out32_1 = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*0))) + X
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*0)) = in32 + X
		Y = out32_1 - (*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*1)))
		X = (Y * opus_int32(int64(silk_resampler_up2_hq_0[1]))) >> 16
		out32_2 = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*1))) + X
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*1)) = out32_1 + X
		Y = out32_2 - (*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*2)))
		X = Y + ((Y * opus_int32(int64(silk_resampler_up2_hq_0[2]))) >> 16)
		out32_1 = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*2))) + X
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*2)) = out32_2 + X
		if (func() opus_int32 {
			if 10 == 1 {
				return (out32_1 >> 1) + (out32_1 & 1)
			}
			return ((out32_1 >> (10 - 1)) + 1) >> 1
		}()) > silk_int16_MAX {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int16(0))*uintptr(k*2))) = silk_int16_MAX
		} else if (func() opus_int32 {
			if 10 == 1 {
				return (out32_1 >> 1) + (out32_1 & 1)
			}
			return ((out32_1 >> (10 - 1)) + 1) >> 1
		}()) < opus_int32(math.MinInt16) {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int16(0))*uintptr(k*2))) = math.MinInt16
		} else if 10 == 1 {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int16(0))*uintptr(k*2))) = opus_int16((out32_1 >> 1) + (out32_1 & 1))
		} else {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int16(0))*uintptr(k*2))) = opus_int16(((out32_1 >> (10 - 1)) + 1) >> 1)
		}
		Y = in32 - (*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*3)))
		X = (Y * opus_int32(int64(silk_resampler_up2_hq_1[0]))) >> 16
		out32_1 = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*3))) + X
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*3)) = in32 + X
		Y = out32_1 - (*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*4)))
		X = (Y * opus_int32(int64(silk_resampler_up2_hq_1[1]))) >> 16
		out32_2 = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*4))) + X
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*4)) = out32_1 + X
		Y = out32_2 - (*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*5)))
		X = Y + ((Y * opus_int32(int64(silk_resampler_up2_hq_1[2]))) >> 16)
		out32_1 = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*5))) + X
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*5)) = out32_2 + X
		if (func() opus_int32 {
			if 10 == 1 {
				return (out32_1 >> 1) + (out32_1 & 1)
			}
			return ((out32_1 >> (10 - 1)) + 1) >> 1
		}()) > silk_int16_MAX {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int16(0))*uintptr(k*2+1))) = silk_int16_MAX
		} else if (func() opus_int32 {
			if 10 == 1 {
				return (out32_1 >> 1) + (out32_1 & 1)
			}
			return ((out32_1 >> (10 - 1)) + 1) >> 1
		}()) < opus_int32(math.MinInt16) {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int16(0))*uintptr(k*2+1))) = math.MinInt16
		} else if 10 == 1 {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int16(0))*uintptr(k*2+1))) = opus_int16((out32_1 >> 1) + (out32_1 & 1))
		} else {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int16(0))*uintptr(k*2+1))) = opus_int16(((out32_1 >> (10 - 1)) + 1) >> 1)
		}
	}
}
func silk_resampler_private_up2_HQ_wrapper(SS unsafe.Pointer, out *opus_int16, in *opus_int16, len_ opus_int32) {
	var S *silk_resampler_state_struct = (*silk_resampler_state_struct)(SS)
	silk_resampler_private_up2_HQ(&S.SIIR[0], out, in, len_)
}
