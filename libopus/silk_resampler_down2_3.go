package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

const ORDER_FIR = 4

func silk_resampler_down2_3(S *int32, out *int16, in *int16, inLen int32) {
	var (
		nSamplesIn int32
		counter    int32
		res_Q6     int32
		buf        *int32
		buf_ptr    *int32
	)
	buf = (*int32)(libc.Malloc(((int(RESAMPLER_MAX_BATCH_SIZE_MS * RESAMPLER_MAX_FS_KHZ)) + ORDER_FIR) * int(unsafe.Sizeof(int32(0)))))
	libc.MemCpy(unsafe.Pointer(buf), unsafe.Pointer(S), int(ORDER_FIR*unsafe.Sizeof(int32(0))))
	for {
		if int(inLen) < (int(RESAMPLER_MAX_BATCH_SIZE_MS * RESAMPLER_MAX_FS_KHZ)) {
			nSamplesIn = inLen
		} else {
			nSamplesIn = int32(int(RESAMPLER_MAX_BATCH_SIZE_MS * RESAMPLER_MAX_FS_KHZ))
		}
		silk_resampler_private_AR2([]int32((*int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(int32(0))*uintptr(ORDER_FIR)))), []int32((*int32)(unsafe.Add(unsafe.Pointer(buf), unsafe.Sizeof(int32(0))*uintptr(ORDER_FIR)))), []int16(in), silk_Resampler_2_3_COEFS_LQ[:], nSamplesIn)
		buf_ptr = buf
		counter = nSamplesIn
		for int(counter) > 2 {
			res_Q6 = int32((int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*0))) * int64(silk_Resampler_2_3_COEFS_LQ[2])) >> 16)
			res_Q6 = int32(int64(res_Q6) + ((int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*1))) * int64(silk_Resampler_2_3_COEFS_LQ[3])) >> 16))
			res_Q6 = int32(int64(res_Q6) + ((int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*2))) * int64(silk_Resampler_2_3_COEFS_LQ[5])) >> 16))
			res_Q6 = int32(int64(res_Q6) + ((int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*3))) * int64(silk_Resampler_2_3_COEFS_LQ[4])) >> 16))
			*func() *int16 {
				p := &out
				x := *p
				*p = (*int16)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(int16(0))*1))
				return x
			}() = int16(func() int {
				if (func() int {
					if 6 == 1 {
						return (int(res_Q6) >> 1) + (int(res_Q6) & 1)
					}
					return ((int(res_Q6) >> (6 - 1)) + 1) >> 1
				}()) > silk_int16_MAX {
					return silk_int16_MAX
				}
				if (func() int {
					if 6 == 1 {
						return (int(res_Q6) >> 1) + (int(res_Q6) & 1)
					}
					return ((int(res_Q6) >> (6 - 1)) + 1) >> 1
				}()) < int(math.MinInt16) {
					return math.MinInt16
				}
				if 6 == 1 {
					return (int(res_Q6) >> 1) + (int(res_Q6) & 1)
				}
				return ((int(res_Q6) >> (6 - 1)) + 1) >> 1
			}())
			res_Q6 = int32((int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*1))) * int64(silk_Resampler_2_3_COEFS_LQ[4])) >> 16)
			res_Q6 = int32(int64(res_Q6) + ((int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*2))) * int64(silk_Resampler_2_3_COEFS_LQ[5])) >> 16))
			res_Q6 = int32(int64(res_Q6) + ((int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*3))) * int64(silk_Resampler_2_3_COEFS_LQ[3])) >> 16))
			res_Q6 = int32(int64(res_Q6) + ((int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*4))) * int64(silk_Resampler_2_3_COEFS_LQ[2])) >> 16))
			*func() *int16 {
				p := &out
				x := *p
				*p = (*int16)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(int16(0))*1))
				return x
			}() = int16(func() int {
				if (func() int {
					if 6 == 1 {
						return (int(res_Q6) >> 1) + (int(res_Q6) & 1)
					}
					return ((int(res_Q6) >> (6 - 1)) + 1) >> 1
				}()) > silk_int16_MAX {
					return silk_int16_MAX
				}
				if (func() int {
					if 6 == 1 {
						return (int(res_Q6) >> 1) + (int(res_Q6) & 1)
					}
					return ((int(res_Q6) >> (6 - 1)) + 1) >> 1
				}()) < int(math.MinInt16) {
					return math.MinInt16
				}
				if 6 == 1 {
					return (int(res_Q6) >> 1) + (int(res_Q6) & 1)
				}
				return ((int(res_Q6) >> (6 - 1)) + 1) >> 1
			}())
			buf_ptr = (*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*3))
			counter -= 3
		}
		in = (*int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(int16(0))*uintptr(nSamplesIn)))
		inLen -= nSamplesIn
		if int(inLen) > 0 {
			libc.MemCpy(unsafe.Pointer(buf), unsafe.Pointer((*int32)(unsafe.Add(unsafe.Pointer(buf), unsafe.Sizeof(int32(0))*uintptr(nSamplesIn)))), int(ORDER_FIR*unsafe.Sizeof(int32(0))))
		} else {
			break
		}
	}
	libc.MemCpy(unsafe.Pointer(S), unsafe.Pointer((*int32)(unsafe.Add(unsafe.Pointer(buf), unsafe.Sizeof(int32(0))*uintptr(nSamplesIn)))), int(ORDER_FIR*unsafe.Sizeof(int32(0))))
}
