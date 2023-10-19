package libopus

import (
	"math"
	"unsafe"
)

var A_fb1_20 opus_int16 = 5394 << 1
var A_fb1_21 opus_int16 = -24290

func silk_ana_filt_bank_1(in *opus_int16, S *opus_int32, outL *opus_int16, outH *opus_int16, N opus_int32) {
	var (
		k     int64
		N2    int64 = int64(N >> 1)
		in32  opus_int32
		X     opus_int32
		Y     opus_int32
		out_1 opus_int32
		out_2 opus_int32
	)
	for k = 0; k < N2; k++ {
		in32 = opus_int32(opus_uint32(opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(opus_int16(0))*uintptr(k*2))))) << 10)
		Y = in32 - (*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*0)))
		X = Y + ((Y * opus_int32(int64(A_fb1_21))) >> 16)
		out_1 = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*0))) + X
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*0)) = in32 + X
		in32 = opus_int32(opus_uint32(opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(opus_int16(0))*uintptr(k*2+1))))) << 10)
		Y = in32 - (*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*1)))
		X = (Y * opus_int32(int64(A_fb1_20))) >> 16
		out_2 = (*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*1))) + X
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*1)) = in32 + X
		if (func() opus_int32 {
			if 11 == 1 {
				return ((out_2 + out_1) >> 1) + ((out_2 + out_1) & 1)
			}
			return (((out_2 + out_1) >> (11 - 1)) + 1) >> 1
		}()) > silk_int16_MAX {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(outL), unsafe.Sizeof(opus_int16(0))*uintptr(k))) = silk_int16_MAX
		} else if (func() opus_int32 {
			if 11 == 1 {
				return ((out_2 + out_1) >> 1) + ((out_2 + out_1) & 1)
			}
			return (((out_2 + out_1) >> (11 - 1)) + 1) >> 1
		}()) < opus_int32(math.MinInt16) {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(outL), unsafe.Sizeof(opus_int16(0))*uintptr(k))) = math.MinInt16
		} else if 11 == 1 {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(outL), unsafe.Sizeof(opus_int16(0))*uintptr(k))) = opus_int16(((out_2 + out_1) >> 1) + ((out_2 + out_1) & 1))
		} else {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(outL), unsafe.Sizeof(opus_int16(0))*uintptr(k))) = opus_int16((((out_2 + out_1) >> (11 - 1)) + 1) >> 1)
		}
		if (func() opus_int32 {
			if 11 == 1 {
				return ((out_2 - out_1) >> 1) + ((out_2 - out_1) & 1)
			}
			return (((out_2 - out_1) >> (11 - 1)) + 1) >> 1
		}()) > silk_int16_MAX {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(outH), unsafe.Sizeof(opus_int16(0))*uintptr(k))) = silk_int16_MAX
		} else if (func() opus_int32 {
			if 11 == 1 {
				return ((out_2 - out_1) >> 1) + ((out_2 - out_1) & 1)
			}
			return (((out_2 - out_1) >> (11 - 1)) + 1) >> 1
		}()) < opus_int32(math.MinInt16) {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(outH), unsafe.Sizeof(opus_int16(0))*uintptr(k))) = math.MinInt16
		} else if 11 == 1 {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(outH), unsafe.Sizeof(opus_int16(0))*uintptr(k))) = opus_int16(((out_2 - out_1) >> 1) + ((out_2 - out_1) & 1))
		} else {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(outH), unsafe.Sizeof(opus_int16(0))*uintptr(k))) = opus_int16((((out_2 - out_1) >> (11 - 1)) + 1) >> 1)
		}
	}
}
