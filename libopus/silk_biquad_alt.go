package libopus

import (
	"math"
	"unsafe"
)

func silk_biquad_alt_stride1(in *int16, B_Q28 *int32, A_Q28 *int32, S *int32, out *int16, len_ int32) {
	var (
		k         int
		inval     int32
		A0_U_Q28  int32
		A0_L_Q28  int32
		A1_U_Q28  int32
		A1_L_Q28  int32
		out32_Q14 int32
	)
	A0_L_Q28 = int32(int(-*(*int32)(unsafe.Add(unsafe.Pointer(A_Q28), unsafe.Sizeof(int32(0))*0))) & 0x3FFF)
	A0_U_Q28 = int32(int(-*(*int32)(unsafe.Add(unsafe.Pointer(A_Q28), unsafe.Sizeof(int32(0))*0))) >> 14)
	A1_L_Q28 = int32(int(-*(*int32)(unsafe.Add(unsafe.Pointer(A_Q28), unsafe.Sizeof(int32(0))*1))) & 0x3FFF)
	A1_U_Q28 = int32(int(-*(*int32)(unsafe.Add(unsafe.Pointer(A_Q28), unsafe.Sizeof(int32(0))*1))) >> 14)
	for k = 0; k < int(len_); k++ {
		inval = int32(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(k))))
		out32_Q14 = int32(int(uint32(int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0)))+((int64(*(*int32)(unsafe.Add(unsafe.Pointer(B_Q28), unsafe.Sizeof(int32(0))*0)))*int64(int16(inval)))>>16)))) << 2)
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0)) = int32(int(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1))) + (func() int {
			if 14 == 1 {
				return (int(int32((int64(out32_Q14)*int64(int16(A0_L_Q28)))>>16)) >> 1) + (int(int32((int64(out32_Q14)*int64(int16(A0_L_Q28)))>>16)) & 1)
			}
			return ((int(int32((int64(out32_Q14)*int64(int16(A0_L_Q28)))>>16)) >> (14 - 1)) + 1) >> 1
		}()))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0)) = int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0))) + ((int64(out32_Q14) * int64(int16(A0_U_Q28))) >> 16))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0)) = int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0))) + ((int64(*(*int32)(unsafe.Add(unsafe.Pointer(B_Q28), unsafe.Sizeof(int32(0))*1))) * int64(int16(inval))) >> 16))
		if 14 == 1 {
			*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1)) = int32((int(int32((int64(out32_Q14)*int64(int16(A1_L_Q28)))>>16)) >> 1) + (int(int32((int64(out32_Q14)*int64(int16(A1_L_Q28)))>>16)) & 1))
		} else {
			*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1)) = int32(((int(int32((int64(out32_Q14)*int64(int16(A1_L_Q28)))>>16)) >> (14 - 1)) + 1) >> 1)
		}
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1)) = int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1))) + ((int64(out32_Q14) * int64(int16(A1_U_Q28))) >> 16))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1)) = int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1))) + ((int64(*(*int32)(unsafe.Add(unsafe.Pointer(B_Q28), unsafe.Sizeof(int32(0))*2))) * int64(int16(inval))) >> 16))
		if ((int(out32_Q14) + (1 << 14) - 1) >> 14) > silk_int16_MAX {
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(k))) = silk_int16_MAX
		} else if ((int(out32_Q14) + (1 << 14) - 1) >> 14) < int(math.MinInt16) {
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(k))) = math.MinInt16
		} else {
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(k))) = int16((int(out32_Q14) + (1 << 14) - 1) >> 14)
		}
	}
}
func silk_biquad_alt_stride2_c(in *int16, B_Q28 *int32, A_Q28 *int32, S *int32, out *int16, len_ int32) {
	var (
		k         int
		A0_U_Q28  int32
		A0_L_Q28  int32
		A1_U_Q28  int32
		A1_L_Q28  int32
		out32_Q14 [2]int32
	)
	A0_L_Q28 = int32(int(-*(*int32)(unsafe.Add(unsafe.Pointer(A_Q28), unsafe.Sizeof(int32(0))*0))) & 0x3FFF)
	A0_U_Q28 = int32(int(-*(*int32)(unsafe.Add(unsafe.Pointer(A_Q28), unsafe.Sizeof(int32(0))*0))) >> 14)
	A1_L_Q28 = int32(int(-*(*int32)(unsafe.Add(unsafe.Pointer(A_Q28), unsafe.Sizeof(int32(0))*1))) & 0x3FFF)
	A1_U_Q28 = int32(int(-*(*int32)(unsafe.Add(unsafe.Pointer(A_Q28), unsafe.Sizeof(int32(0))*1))) >> 14)
	for k = 0; k < int(len_); k++ {
		out32_Q14[0] = int32(int(uint32(int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0)))+((int64(*(*int32)(unsafe.Add(unsafe.Pointer(B_Q28), unsafe.Sizeof(int32(0))*0)))*int64(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(k*2+0)))))>>16)))) << 2)
		out32_Q14[1] = int32(int(uint32(int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*2)))+((int64(*(*int32)(unsafe.Add(unsafe.Pointer(B_Q28), unsafe.Sizeof(int32(0))*0)))*int64(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(k*2+1)))))>>16)))) << 2)
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0)) = int32(int(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1))) + (func() int {
			if 14 == 1 {
				return (int(int32((int64(out32_Q14[0])*int64(int16(A0_L_Q28)))>>16)) >> 1) + (int(int32((int64(out32_Q14[0])*int64(int16(A0_L_Q28)))>>16)) & 1)
			}
			return ((int(int32((int64(out32_Q14[0])*int64(int16(A0_L_Q28)))>>16)) >> (14 - 1)) + 1) >> 1
		}()))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*2)) = int32(int(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*3))) + (func() int {
			if 14 == 1 {
				return (int(int32((int64(out32_Q14[1])*int64(int16(A0_L_Q28)))>>16)) >> 1) + (int(int32((int64(out32_Q14[1])*int64(int16(A0_L_Q28)))>>16)) & 1)
			}
			return ((int(int32((int64(out32_Q14[1])*int64(int16(A0_L_Q28)))>>16)) >> (14 - 1)) + 1) >> 1
		}()))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0)) = int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0))) + ((int64(out32_Q14[0]) * int64(int16(A0_U_Q28))) >> 16))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*2)) = int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*2))) + ((int64(out32_Q14[1]) * int64(int16(A0_U_Q28))) >> 16))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0)) = int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0))) + ((int64(*(*int32)(unsafe.Add(unsafe.Pointer(B_Q28), unsafe.Sizeof(int32(0))*1))) * int64(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(k*2+0))))) >> 16))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*2)) = int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*2))) + ((int64(*(*int32)(unsafe.Add(unsafe.Pointer(B_Q28), unsafe.Sizeof(int32(0))*1))) * int64(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(k*2+1))))) >> 16))
		if 14 == 1 {
			*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1)) = int32((int(int32((int64(out32_Q14[0])*int64(int16(A1_L_Q28)))>>16)) >> 1) + (int(int32((int64(out32_Q14[0])*int64(int16(A1_L_Q28)))>>16)) & 1))
		} else {
			*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1)) = int32(((int(int32((int64(out32_Q14[0])*int64(int16(A1_L_Q28)))>>16)) >> (14 - 1)) + 1) >> 1)
		}
		if 14 == 1 {
			*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*3)) = int32((int(int32((int64(out32_Q14[1])*int64(int16(A1_L_Q28)))>>16)) >> 1) + (int(int32((int64(out32_Q14[1])*int64(int16(A1_L_Q28)))>>16)) & 1))
		} else {
			*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*3)) = int32(((int(int32((int64(out32_Q14[1])*int64(int16(A1_L_Q28)))>>16)) >> (14 - 1)) + 1) >> 1)
		}
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1)) = int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1))) + ((int64(out32_Q14[0]) * int64(int16(A1_U_Q28))) >> 16))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*3)) = int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*3))) + ((int64(out32_Q14[1]) * int64(int16(A1_U_Q28))) >> 16))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1)) = int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1))) + ((int64(*(*int32)(unsafe.Add(unsafe.Pointer(B_Q28), unsafe.Sizeof(int32(0))*2))) * int64(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(k*2+0))))) >> 16))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*3)) = int32(int64(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*3))) + ((int64(*(*int32)(unsafe.Add(unsafe.Pointer(B_Q28), unsafe.Sizeof(int32(0))*2))) * int64(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(k*2+1))))) >> 16))
		if ((int(out32_Q14[0]) + (1 << 14) - 1) >> 14) > silk_int16_MAX {
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(k*2+0))) = silk_int16_MAX
		} else if ((int(out32_Q14[0]) + (1 << 14) - 1) >> 14) < int(math.MinInt16) {
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(k*2+0))) = math.MinInt16
		} else {
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(k*2+0))) = int16((int(out32_Q14[0]) + (1 << 14) - 1) >> 14)
		}
		if ((int(out32_Q14[1]) + (1 << 14) - 1) >> 14) > silk_int16_MAX {
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(k*2+1))) = silk_int16_MAX
		} else if ((int(out32_Q14[1]) + (1 << 14) - 1) >> 14) < int(math.MinInt16) {
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(k*2+1))) = math.MinInt16
		} else {
			*(*int16)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int16(0))*uintptr(k*2+1))) = int16((int(out32_Q14[1]) + (1 << 14) - 1) >> 14)
		}
	}
}
