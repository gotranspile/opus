package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

func silk_resampler_private_IIR_FIR_INTERPOL(out *int16, buf *int16, max_index_Q16 int32, index_increment_Q16 int32) *int16 {
	var (
		index_Q16   int32
		res_Q15     int32
		buf_ptr     *int16
		table_index int32
	)
	for index_Q16 = 0; int(index_Q16) < int(max_index_Q16); index_Q16 += index_increment_Q16 {
		table_index = int32(((int(index_Q16) & math.MaxUint16) * 12) >> 16)
		buf_ptr = (*int16)(unsafe.Add(unsafe.Pointer(buf), unsafe.Sizeof(int16(0))*uintptr(int(index_Q16)>>16)))
		res_Q15 = int32(int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int16(0))*0)))) * int(int32(silk_resampler_frac_FIR_12[table_index][0])))
		res_Q15 = int32(int(res_Q15) + int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int16(0))*1))))*int(int32(silk_resampler_frac_FIR_12[table_index][1])))
		res_Q15 = int32(int(res_Q15) + int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int16(0))*2))))*int(int32(silk_resampler_frac_FIR_12[table_index][2])))
		res_Q15 = int32(int(res_Q15) + int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int16(0))*3))))*int(int32(silk_resampler_frac_FIR_12[table_index][3])))
		res_Q15 = int32(int(res_Q15) + int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int16(0))*4))))*int(int32(silk_resampler_frac_FIR_12[11-int(table_index)][3])))
		res_Q15 = int32(int(res_Q15) + int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int16(0))*5))))*int(int32(silk_resampler_frac_FIR_12[11-int(table_index)][2])))
		res_Q15 = int32(int(res_Q15) + int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int16(0))*6))))*int(int32(silk_resampler_frac_FIR_12[11-int(table_index)][1])))
		res_Q15 = int32(int(res_Q15) + int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int16(0))*7))))*int(int32(silk_resampler_frac_FIR_12[11-int(table_index)][0])))
		*func() *int16 {
			p := &out
			x := *p
			*p = (*int16)(unsafe.Add(unsafe.Pointer(*p), unsafe.Sizeof(int16(0))*1))
			return x
		}() = int16(func() int {
			if (func() int {
				if 15 == 1 {
					return (int(res_Q15) >> 1) + (int(res_Q15) & 1)
				}
				return ((int(res_Q15) >> (15 - 1)) + 1) >> 1
			}()) > silk_int16_MAX {
				return silk_int16_MAX
			}
			if (func() int {
				if 15 == 1 {
					return (int(res_Q15) >> 1) + (int(res_Q15) & 1)
				}
				return ((int(res_Q15) >> (15 - 1)) + 1) >> 1
			}()) < int(math.MinInt16) {
				return math.MinInt16
			}
			if 15 == 1 {
				return (int(res_Q15) >> 1) + (int(res_Q15) & 1)
			}
			return ((int(res_Q15) >> (15 - 1)) + 1) >> 1
		}())
	}
	return out
}
func silk_resampler_private_IIR_FIR(SS unsafe.Pointer, out []int16, in []int16, inLen int32) {
	var (
		S                   *silk_resampler_state_struct = (*silk_resampler_state_struct)(SS)
		nSamplesIn          int32
		max_index_Q16       int32
		index_increment_Q16 int32
		buf                 *int16
	)
	buf = (*int16)(libc.Malloc((S.BatchSize*2 + RESAMPLER_ORDER_FIR_12) * int(unsafe.Sizeof(int16(0)))))
	libc.MemCpy(unsafe.Pointer(buf), unsafe.Pointer(&S.SFIR.I16[0]), int(RESAMPLER_ORDER_FIR_12*unsafe.Sizeof(int16(0))))
	index_increment_Q16 = S.InvRatio_Q16
	for {
		if int(inLen) < S.BatchSize {
			nSamplesIn = inLen
		} else {
			nSamplesIn = int32(S.BatchSize)
		}
		silk_resampler_private_up2_HQ(&S.SIIR[0], (*int16)(unsafe.Add(unsafe.Pointer(buf), unsafe.Sizeof(int16(0))*uintptr(RESAMPLER_ORDER_FIR_12))), &in[0], nSamplesIn)
		max_index_Q16 = int32(int(uint32(nSamplesIn)) << (16 + 1))
		out = []int16(silk_resampler_private_IIR_FIR_INTERPOL(&out[0], buf, max_index_Q16, index_increment_Q16))
		in += []int16(nSamplesIn)
		inLen -= nSamplesIn
		if int(inLen) > 0 {
			libc.MemCpy(unsafe.Pointer(buf), unsafe.Pointer((*int16)(unsafe.Add(unsafe.Pointer(buf), unsafe.Sizeof(int16(0))*uintptr(int(nSamplesIn)<<1)))), int(RESAMPLER_ORDER_FIR_12*unsafe.Sizeof(int16(0))))
		} else {
			break
		}
	}
	libc.MemCpy(unsafe.Pointer(&S.SFIR.I16[0]), unsafe.Pointer((*int16)(unsafe.Add(unsafe.Pointer(buf), unsafe.Sizeof(int16(0))*uintptr(int(nSamplesIn)<<1)))), int(RESAMPLER_ORDER_FIR_12*unsafe.Sizeof(int16(0))))
}
