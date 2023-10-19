package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

const ORDER_FIR = 4

func silk_resampler_down2_3(S *opus_int32, out *opus_int16, in *opus_int16, inLen opus_int32) {
	var (
		nSamplesIn opus_int32
		counter    opus_int32
		res_Q6     opus_int32
		buf        *opus_int32
		buf_ptr    *opus_int32
	)
	buf = (*opus_int32)(libc.Malloc(int(uintptr((RESAMPLER_MAX_BATCH_SIZE_MS*RESAMPLER_MAX_FS_KHZ)+ORDER_FIR) * unsafe.Sizeof(opus_int32(0)))))
	libc.MemCpy(unsafe.Pointer(buf), unsafe.Pointer(S), int(ORDER_FIR*unsafe.Sizeof(opus_int32(0))))
	for {
		if inLen < opus_int32(RESAMPLER_MAX_BATCH_SIZE_MS*RESAMPLER_MAX_FS_KHZ) {
			nSamplesIn = inLen
		} else {
			nSamplesIn = opus_int32(RESAMPLER_MAX_BATCH_SIZE_MS * RESAMPLER_MAX_FS_KHZ)
		}
		silk_resampler_private_AR2([0]opus_int32((*opus_int32)(unsafe.Add(unsafe.Pointer(S), unsafe.Sizeof(opus_int32(0))*uintptr(ORDER_FIR)))), [0]opus_int32((*opus_int32)(unsafe.Add(unsafe.Pointer(buf), unsafe.Sizeof(opus_int32(0))*uintptr(ORDER_FIR)))), [0]opus_int16(in), silk_Resampler_2_3_COEFS_LQ[:], nSamplesIn)
		buf_ptr = buf
		counter = nSamplesIn
		for counter > 2 {
			res_Q6 = ((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*0))) * opus_int32(int64(silk_Resampler_2_3_COEFS_LQ[2]))) >> 16
			res_Q6 = res_Q6 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*1))) * opus_int32(int64(silk_Resampler_2_3_COEFS_LQ[3]))) >> 16)
			res_Q6 = res_Q6 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*2))) * opus_int32(int64(silk_Resampler_2_3_COEFS_LQ[5]))) >> 16)
			res_Q6 = res_Q6 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*3))) * opus_int32(int64(silk_Resampler_2_3_COEFS_LQ[4]))) >> 16)
			*func() *opus_int16 {
				p := &out
				x := *p
				*p = (*opus_int16)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(opus_int16(0))*1))
				return x
			}() = opus_int16(func() opus_int32 {
				if (func() opus_int32 {
					if 6 == 1 {
						return (res_Q6 >> 1) + (res_Q6 & 1)
					}
					return ((res_Q6 >> (6 - 1)) + 1) >> 1
				}()) > silk_int16_MAX {
					return silk_int16_MAX
				}
				if (func() opus_int32 {
					if 6 == 1 {
						return (res_Q6 >> 1) + (res_Q6 & 1)
					}
					return ((res_Q6 >> (6 - 1)) + 1) >> 1
				}()) < opus_int32(math.MinInt16) {
					return math.MinInt16
				}
				if 6 == 1 {
					return (res_Q6 >> 1) + (res_Q6 & 1)
				}
				return ((res_Q6 >> (6 - 1)) + 1) >> 1
			}())
			res_Q6 = ((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*1))) * opus_int32(int64(silk_Resampler_2_3_COEFS_LQ[4]))) >> 16
			res_Q6 = res_Q6 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*2))) * opus_int32(int64(silk_Resampler_2_3_COEFS_LQ[5]))) >> 16)
			res_Q6 = res_Q6 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*3))) * opus_int32(int64(silk_Resampler_2_3_COEFS_LQ[3]))) >> 16)
			res_Q6 = res_Q6 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*4))) * opus_int32(int64(silk_Resampler_2_3_COEFS_LQ[2]))) >> 16)
			*func() *opus_int16 {
				p := &out
				x := *p
				*p = (*opus_int16)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(opus_int16(0))*1))
				return x
			}() = opus_int16(func() opus_int32 {
				if (func() opus_int32 {
					if 6 == 1 {
						return (res_Q6 >> 1) + (res_Q6 & 1)
					}
					return ((res_Q6 >> (6 - 1)) + 1) >> 1
				}()) > silk_int16_MAX {
					return silk_int16_MAX
				}
				if (func() opus_int32 {
					if 6 == 1 {
						return (res_Q6 >> 1) + (res_Q6 & 1)
					}
					return ((res_Q6 >> (6 - 1)) + 1) >> 1
				}()) < opus_int32(math.MinInt16) {
					return math.MinInt16
				}
				if 6 == 1 {
					return (res_Q6 >> 1) + (res_Q6 & 1)
				}
				return ((res_Q6 >> (6 - 1)) + 1) >> 1
			}())
			buf_ptr = (*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*3))
			counter -= 3
		}
		in = (*opus_int16)(unsafe.Add(unsafe.Pointer(in), unsafe.Sizeof(opus_int16(0))*uintptr(nSamplesIn)))
		inLen -= nSamplesIn
		if inLen > 0 {
			libc.MemCpy(unsafe.Pointer(buf), unsafe.Pointer((*opus_int32)(unsafe.Add(unsafe.Pointer(buf), unsafe.Sizeof(opus_int32(0))*uintptr(nSamplesIn)))), int(ORDER_FIR*unsafe.Sizeof(opus_int32(0))))
		} else {
			break
		}
	}
	libc.MemCpy(unsafe.Pointer(S), unsafe.Pointer((*opus_int32)(unsafe.Add(unsafe.Pointer(buf), unsafe.Sizeof(opus_int32(0))*uintptr(nSamplesIn)))), int(ORDER_FIR*unsafe.Sizeof(opus_int32(0))))
}
