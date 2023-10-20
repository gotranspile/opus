package libopus

import (
	"math"
	"unsafe"
)

var A_fb1_20 int16 = 5394 << 1
var A_fb1_21 int16 = -24290

func silk_ana_filt_bank_1(in *int16, S *int32, outL *int16, outH *int16, N int32) {
	var (
		k     int
		N2    int = (int(N) >> 1)
		in32  int32
		X     int32
		Y     int32
		out_1 int32
		out_2 int32
	)
	for k = 0; k < N2; k++ {
		in32 = int32(int(uint32(int32(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(k*2)))))) << 10)
		Y = int32(int(in32) - int(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0))))
		X = int32(int64(Y) + ((int64(Y) * int64(A_fb1_21)) >> 16))
		out_1 = int32(int(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0))) + int(X))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*0)) = int32(int(in32) + int(X))
		in32 = int32(int(uint32(int32(*(*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(k*2+1)))))) << 10)
		Y = int32(int(in32) - int(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1))))
		X = int32((int64(Y) * int64(A_fb1_20)) >> 16)
		out_2 = int32(int(*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1))) + int(X))
		*(*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*1)) = int32(int(in32) + int(X))
		if (func() int {
			if 11 == 1 {
				return ((int(out_2) + int(out_1)) >> 1) + ((int(out_2) + int(out_1)) & 1)
			}
			return (((int(out_2) + int(out_1)) >> (11 - 1)) + 1) >> 1
		}()) > silk_int16_MAX {
			*(*int16)(unsafe.Add(unsafe.Pointer(outL), unsafe.Sizeof(int16(0))*uintptr(k))) = silk_int16_MAX
		} else if (func() int {
			if 11 == 1 {
				return ((int(out_2) + int(out_1)) >> 1) + ((int(out_2) + int(out_1)) & 1)
			}
			return (((int(out_2) + int(out_1)) >> (11 - 1)) + 1) >> 1
		}()) < int(math.MinInt16) {
			*(*int16)(unsafe.Add(unsafe.Pointer(outL), unsafe.Sizeof(int16(0))*uintptr(k))) = math.MinInt16
		} else if 11 == 1 {
			*(*int16)(unsafe.Add(unsafe.Pointer(outL), unsafe.Sizeof(int16(0))*uintptr(k))) = int16(((int(out_2) + int(out_1)) >> 1) + ((int(out_2) + int(out_1)) & 1))
		} else {
			*(*int16)(unsafe.Add(unsafe.Pointer(outL), unsafe.Sizeof(int16(0))*uintptr(k))) = int16((((int(out_2) + int(out_1)) >> (11 - 1)) + 1) >> 1)
		}
		if (func() int {
			if 11 == 1 {
				return ((int(out_2) - int(out_1)) >> 1) + ((int(out_2) - int(out_1)) & 1)
			}
			return (((int(out_2) - int(out_1)) >> (11 - 1)) + 1) >> 1
		}()) > silk_int16_MAX {
			*(*int16)(unsafe.Add(unsafe.Pointer(outH), unsafe.Sizeof(int16(0))*uintptr(k))) = silk_int16_MAX
		} else if (func() int {
			if 11 == 1 {
				return ((int(out_2) - int(out_1)) >> 1) + ((int(out_2) - int(out_1)) & 1)
			}
			return (((int(out_2) - int(out_1)) >> (11 - 1)) + 1) >> 1
		}()) < int(math.MinInt16) {
			*(*int16)(unsafe.Add(unsafe.Pointer(outH), unsafe.Sizeof(int16(0))*uintptr(k))) = math.MinInt16
		} else if 11 == 1 {
			*(*int16)(unsafe.Add(unsafe.Pointer(outH), unsafe.Sizeof(int16(0))*uintptr(k))) = int16(((int(out_2) - int(out_1)) >> 1) + ((int(out_2) - int(out_1)) & 1))
		} else {
			*(*int16)(unsafe.Add(unsafe.Pointer(outH), unsafe.Sizeof(int16(0))*uintptr(k))) = int16((((int(out_2) - int(out_1)) >> (11 - 1)) + 1) >> 1)
		}
	}
}
