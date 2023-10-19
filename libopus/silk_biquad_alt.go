package libopus

import (
	"math"
	"unsafe"
)

func silk_biquad_alt_stride1(in *opus_int16, B_Q28 *opus_int32, A_Q28 *opus_int32, S *opus_int32, out *opus_int16, len_ opus_int32) {
	var (
		k         int64
		inval     opus_int32
		A0_U_Q28  opus_int32
		A0_L_Q28  opus_int32
		A1_U_Q28  opus_int32
		A1_L_Q28  opus_int32
		out32_Q14 opus_int32
	)
	A0_L_Q28 = (-*(*opus_int32)(unsafe.Add(unsafe.Pointer(A_Q28), unsafe.Sizeof(opus_int32(0))*0))) & 0x3FFF
	A0_U_Q28 = (-*(*opus_int32)(unsafe.Add(unsafe.Pointer(A_Q28), unsafe.Sizeof(opus_int32(0))*0))) >> 14
	A1_L_Q28 = (-*(*opus_int32)(unsafe.Add(unsafe.Pointer(A_Q28), unsafe.Sizeof(opus_int32(0))*1))) & 0x3FFF
	A1_U_Q28 = (-*(*opus_int32)(unsafe.Add(unsafe.Pointer(A_Q28), unsafe.Sizeof(opus_int32(0))*1))) >> 14
	for k = 0; k < int64(len_); k++ {
		inval = opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(opus_int16(0))*uintptr(k))))
		out32_Q14 = opus_int32(opus_uint32((*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*0)))+(((*(*opus_int32)(unsafe.Add(unsafe.Pointer(B_Q28), unsafe.Sizeof(opus_int32(0))*0)))*opus_int32(int64(opus_int16(inval))))>>16)) << 2)
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*0)) = *(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*1)) + (func() opus_int32 {
			if 14 == 1 {
				return (((out32_Q14 * opus_int32(int64(opus_int16(A0_L_Q28)))) >> 16) >> 1) + (((out32_Q14 * opus_int32(int64(opus_int16(A0_L_Q28)))) >> 16) & 1)
			}
			return ((((out32_Q14 * opus_int32(int64(opus_int16(A0_L_Q28)))) >> 16) >> (14 - 1)) + 1) >> 1
		}())
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*0)) = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*0))) + ((out32_Q14 * opus_int32(int64(opus_int16(A0_U_Q28)))) >> 16)
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*0)) = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*0))) + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(B_Q28), unsafe.Sizeof(opus_int32(0))*1))) * opus_int32(int64(opus_int16(inval)))) >> 16)
		if 14 == 1 {
			*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*1)) = (((out32_Q14 * opus_int32(int64(opus_int16(A1_L_Q28)))) >> 16) >> 1) + (((out32_Q14 * opus_int32(int64(opus_int16(A1_L_Q28)))) >> 16) & 1)
		} else {
			*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*1)) = ((((out32_Q14 * opus_int32(int64(opus_int16(A1_L_Q28)))) >> 16) >> (14 - 1)) + 1) >> 1
		}
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*1)) = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*1))) + ((out32_Q14 * opus_int32(int64(opus_int16(A1_U_Q28)))) >> 16)
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*1)) = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*1))) + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(B_Q28), unsafe.Sizeof(opus_int32(0))*2))) * opus_int32(int64(opus_int16(inval)))) >> 16)
		if ((out32_Q14 + (1 << 14) - 1) >> 14) > silk_int16_MAX {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int16(0))*uintptr(k))) = silk_int16_MAX
		} else if ((out32_Q14 + (1 << 14) - 1) >> 14) < opus_int32(math.MinInt16) {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int16(0))*uintptr(k))) = math.MinInt16
		} else {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int16(0))*uintptr(k))) = opus_int16((out32_Q14 + (1 << 14) - 1) >> 14)
		}
	}
}
func silk_biquad_alt_stride2_c(in *opus_int16, B_Q28 *opus_int32, A_Q28 *opus_int32, S *opus_int32, out *opus_int16, len_ opus_int32) {
	var (
		k         int64
		A0_U_Q28  opus_int32
		A0_L_Q28  opus_int32
		A1_U_Q28  opus_int32
		A1_L_Q28  opus_int32
		out32_Q14 [2]opus_int32
	)
	A0_L_Q28 = (-*(*opus_int32)(unsafe.Add(unsafe.Pointer(A_Q28), unsafe.Sizeof(opus_int32(0))*0))) & 0x3FFF
	A0_U_Q28 = (-*(*opus_int32)(unsafe.Add(unsafe.Pointer(A_Q28), unsafe.Sizeof(opus_int32(0))*0))) >> 14
	A1_L_Q28 = (-*(*opus_int32)(unsafe.Add(unsafe.Pointer(A_Q28), unsafe.Sizeof(opus_int32(0))*1))) & 0x3FFF
	A1_U_Q28 = (-*(*opus_int32)(unsafe.Add(unsafe.Pointer(A_Q28), unsafe.Sizeof(opus_int32(0))*1))) >> 14
	for k = 0; k < int64(len_); k++ {
		out32_Q14[0] = opus_int32(opus_uint32((*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*0)))+(((*(*opus_int32)(unsafe.Add(unsafe.Pointer(B_Q28), unsafe.Sizeof(opus_int32(0))*0)))*opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(opus_int16(0))*uintptr(k*2+0))))))>>16)) << 2)
		out32_Q14[1] = opus_int32(opus_uint32((*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*2)))+(((*(*opus_int32)(unsafe.Add(unsafe.Pointer(B_Q28), unsafe.Sizeof(opus_int32(0))*0)))*opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(opus_int16(0))*uintptr(k*2+1))))))>>16)) << 2)
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*0)) = *(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*1)) + (func() opus_int32 {
			if 14 == 1 {
				return ((((out32_Q14[0]) * opus_int32(int64(opus_int16(A0_L_Q28)))) >> 16) >> 1) + ((((out32_Q14[0]) * opus_int32(int64(opus_int16(A0_L_Q28)))) >> 16) & 1)
			}
			return (((((out32_Q14[0]) * opus_int32(int64(opus_int16(A0_L_Q28)))) >> 16) >> (14 - 1)) + 1) >> 1
		}())
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*2)) = *(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*3)) + (func() opus_int32 {
			if 14 == 1 {
				return ((((out32_Q14[1]) * opus_int32(int64(opus_int16(A0_L_Q28)))) >> 16) >> 1) + ((((out32_Q14[1]) * opus_int32(int64(opus_int16(A0_L_Q28)))) >> 16) & 1)
			}
			return (((((out32_Q14[1]) * opus_int32(int64(opus_int16(A0_L_Q28)))) >> 16) >> (14 - 1)) + 1) >> 1
		}())
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*0)) = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*0))) + (((out32_Q14[0]) * opus_int32(int64(opus_int16(A0_U_Q28)))) >> 16)
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*2)) = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*2))) + (((out32_Q14[1]) * opus_int32(int64(opus_int16(A0_U_Q28)))) >> 16)
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*0)) = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*0))) + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(B_Q28), unsafe.Sizeof(opus_int32(0))*1))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(opus_int16(0))*uintptr(k*2+0)))))) >> 16)
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*2)) = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*2))) + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(B_Q28), unsafe.Sizeof(opus_int32(0))*1))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(opus_int16(0))*uintptr(k*2+1)))))) >> 16)
		if 14 == 1 {
			*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*1)) = ((((out32_Q14[0]) * opus_int32(int64(opus_int16(A1_L_Q28)))) >> 16) >> 1) + ((((out32_Q14[0]) * opus_int32(int64(opus_int16(A1_L_Q28)))) >> 16) & 1)
		} else {
			*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*1)) = (((((out32_Q14[0]) * opus_int32(int64(opus_int16(A1_L_Q28)))) >> 16) >> (14 - 1)) + 1) >> 1
		}
		if 14 == 1 {
			*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*3)) = ((((out32_Q14[1]) * opus_int32(int64(opus_int16(A1_L_Q28)))) >> 16) >> 1) + ((((out32_Q14[1]) * opus_int32(int64(opus_int16(A1_L_Q28)))) >> 16) & 1)
		} else {
			*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*3)) = (((((out32_Q14[1]) * opus_int32(int64(opus_int16(A1_L_Q28)))) >> 16) >> (14 - 1)) + 1) >> 1
		}
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*1)) = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*1))) + (((out32_Q14[0]) * opus_int32(int64(opus_int16(A1_U_Q28)))) >> 16)
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*3)) = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*3))) + (((out32_Q14[1]) * opus_int32(int64(opus_int16(A1_U_Q28)))) >> 16)
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*1)) = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*1))) + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(B_Q28), unsafe.Sizeof(opus_int32(0))*2))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(opus_int16(0))*uintptr(k*2+0)))))) >> 16)
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*3)) = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*3))) + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(B_Q28), unsafe.Sizeof(opus_int32(0))*2))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(opus_int16(0))*uintptr(k*2+1)))))) >> 16)
		if ((out32_Q14[0] + (1 << 14) - 1) >> 14) > silk_int16_MAX {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int16(0))*uintptr(k*2+0))) = silk_int16_MAX
		} else if ((out32_Q14[0] + (1 << 14) - 1) >> 14) < opus_int32(math.MinInt16) {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int16(0))*uintptr(k*2+0))) = math.MinInt16
		} else {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int16(0))*uintptr(k*2+0))) = opus_int16((out32_Q14[0] + (1 << 14) - 1) >> 14)
		}
		if ((out32_Q14[1] + (1 << 14) - 1) >> 14) > silk_int16_MAX {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int16(0))*uintptr(k*2+1))) = silk_int16_MAX
		} else if ((out32_Q14[1] + (1 << 14) - 1) >> 14) < opus_int32(math.MinInt16) {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int16(0))*uintptr(k*2+1))) = math.MinInt16
		} else {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int16(0))*uintptr(k*2+1))) = opus_int16((out32_Q14[1] + (1 << 14) - 1) >> 14)
		}
	}
}
