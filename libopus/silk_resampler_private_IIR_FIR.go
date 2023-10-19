package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

func silk_resampler_private_IIR_FIR_INTERPOL(out *opus_int16, buf *opus_int16, max_index_Q16 opus_int32, index_increment_Q16 opus_int32) *opus_int16 {
	var (
		index_Q16   opus_int32
		res_Q15     opus_int32
		buf_ptr     *opus_int16
		table_index opus_int32
	)
	for index_Q16 = 0; index_Q16 < max_index_Q16; index_Q16 += index_increment_Q16 {
		table_index = ((index_Q16 & math.MaxUint16) * 12) >> 16
		buf_ptr = (*opus_int16)(unsafe.Add(unsafe.Pointer(buf), unsafe.Sizeof(opus_int16(0))*uintptr(index_Q16>>16)))
		res_Q15 = opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int16(0))*0))) * opus_int32(silk_resampler_frac_FIR_12[table_index][0])
		res_Q15 = res_Q15 + (opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int16(0))*1))))*opus_int32(silk_resampler_frac_FIR_12[table_index][1])
		res_Q15 = res_Q15 + (opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int16(0))*2))))*opus_int32(silk_resampler_frac_FIR_12[table_index][2])
		res_Q15 = res_Q15 + (opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int16(0))*3))))*opus_int32(silk_resampler_frac_FIR_12[table_index][3])
		res_Q15 = res_Q15 + (opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int16(0))*4))))*opus_int32(silk_resampler_frac_FIR_12[11-table_index][3])
		res_Q15 = res_Q15 + (opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int16(0))*5))))*opus_int32(silk_resampler_frac_FIR_12[11-table_index][2])
		res_Q15 = res_Q15 + (opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int16(0))*6))))*opus_int32(silk_resampler_frac_FIR_12[11-table_index][1])
		res_Q15 = res_Q15 + (opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int16(0))*7))))*opus_int32(silk_resampler_frac_FIR_12[11-table_index][0])
		*func() *opus_int16 {
			p := &out
			x := *p
			*p = (*opus_int16)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(opus_int16(0))*1))
			return x
		}() = opus_int16(func() opus_int32 {
			if (func() opus_int32 {
				if 15 == 1 {
					return (res_Q15 >> 1) + (res_Q15 & 1)
				}
				return ((res_Q15 >> (15 - 1)) + 1) >> 1
			}()) > silk_int16_MAX {
				return silk_int16_MAX
			}
			if (func() opus_int32 {
				if 15 == 1 {
					return (res_Q15 >> 1) + (res_Q15 & 1)
				}
				return ((res_Q15 >> (15 - 1)) + 1) >> 1
			}()) < opus_int32(math.MinInt16) {
				return math.MinInt16
			}
			if 15 == 1 {
				return (res_Q15 >> 1) + (res_Q15 & 1)
			}
			return ((res_Q15 >> (15 - 1)) + 1) >> 1
		}())
	}
	return out
}
func silk_resampler_private_IIR_FIR(SS unsafe.Pointer, out [0]opus_int16, in [0]opus_int16, inLen opus_int32) {
	var (
		S                   *silk_resampler_state_struct = (*silk_resampler_state_struct)(SS)
		nSamplesIn          opus_int32
		max_index_Q16       opus_int32
		index_increment_Q16 opus_int32
		buf                 *opus_int16
	)
	buf = (*opus_int16)(libc.Malloc(int((S.BatchSize*2 + RESAMPLER_ORDER_FIR_12) * int64(unsafe.Sizeof(opus_int16(0))))))
	libc.MemCpy(unsafe.Pointer(buf), unsafe.Pointer(&S.SFIR.I16[0]), int(RESAMPLER_ORDER_FIR_12*unsafe.Sizeof(opus_int16(0))))
	index_increment_Q16 = S.InvRatio_Q16
	for {
		if inLen < opus_int32(S.BatchSize) {
			nSamplesIn = inLen
		} else {
			nSamplesIn = opus_int32(S.BatchSize)
		}
		silk_resampler_private_up2_HQ(&S.SIIR[0], (*opus_int16)(unsafe.Add(unsafe.Pointer(buf), unsafe.Sizeof(opus_int16(0))*uintptr(RESAMPLER_ORDER_FIR_12))), &in[0], nSamplesIn)
		max_index_Q16 = opus_int32(opus_uint32(nSamplesIn) << (16 + 1))
		out = [0]opus_int16(silk_resampler_private_IIR_FIR_INTERPOL(&out[0], buf, max_index_Q16, index_increment_Q16))
		in += [0]opus_int16(nSamplesIn)
		inLen -= nSamplesIn
		if inLen > 0 {
			libc.MemCpy(unsafe.Pointer(buf), unsafe.Pointer((*opus_int16)(unsafe.Add(unsafe.Pointer(buf), unsafe.Sizeof(opus_int16(0))*uintptr(nSamplesIn<<1)))), int(RESAMPLER_ORDER_FIR_12*unsafe.Sizeof(opus_int16(0))))
		} else {
			break
		}
	}
	libc.MemCpy(unsafe.Pointer(&S.SFIR.I16[0]), unsafe.Pointer((*opus_int16)(unsafe.Add(unsafe.Pointer(buf), unsafe.Sizeof(opus_int16(0))*uintptr(nSamplesIn<<1)))), int(RESAMPLER_ORDER_FIR_12*unsafe.Sizeof(opus_int16(0))))
}
