package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

func silk_resampler_private_down_FIR_INTERPOL(out *opus_int16, buf *opus_int32, FIR_Coefs *opus_int16, FIR_Order int64, FIR_Fracs int64, max_index_Q16 opus_int32, index_increment_Q16 opus_int32) *opus_int16 {
	var (
		index_Q16    opus_int32
		res_Q6       opus_int32
		buf_ptr      *opus_int32
		interpol_ind opus_int32
		interpol_ptr *opus_int16
	)
	switch FIR_Order {
	case RESAMPLER_DOWN_ORDER_FIR0:
		for index_Q16 = 0; index_Q16 < max_index_Q16; index_Q16 += index_increment_Q16 {
			buf_ptr = (*opus_int32)(unsafe.Add(unsafe.Pointer(buf), unsafe.Sizeof(opus_int32(0))*uintptr(index_Q16>>16)))
			interpol_ind = ((index_Q16 & math.MaxUint16) * opus_int32(int64(opus_int16(FIR_Fracs)))) >> 16
			interpol_ptr = (*opus_int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(opus_int16(0))*uintptr(opus_int32(RESAMPLER_DOWN_ORDER_FIR0/2)*interpol_ind)))
			res_Q6 = ((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*0))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(opus_int16(0))*0))))) >> 16
			res_Q6 = res_Q6 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*1))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(opus_int16(0))*1))))) >> 16)
			res_Q6 = res_Q6 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*2))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(opus_int16(0))*2))))) >> 16)
			res_Q6 = res_Q6 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*3))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(opus_int16(0))*3))))) >> 16)
			res_Q6 = res_Q6 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*4))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(opus_int16(0))*4))))) >> 16)
			res_Q6 = res_Q6 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*5))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(opus_int16(0))*5))))) >> 16)
			res_Q6 = res_Q6 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*6))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(opus_int16(0))*6))))) >> 16)
			res_Q6 = res_Q6 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*7))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(opus_int16(0))*7))))) >> 16)
			res_Q6 = res_Q6 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*8))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(opus_int16(0))*8))))) >> 16)
			interpol_ptr = (*opus_int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(opus_int16(0))*uintptr(RESAMPLER_DOWN_ORDER_FIR0/2*(FIR_Fracs-1-int64(interpol_ind)))))
			res_Q6 = res_Q6 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*17))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(opus_int16(0))*0))))) >> 16)
			res_Q6 = res_Q6 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*16))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(opus_int16(0))*1))))) >> 16)
			res_Q6 = res_Q6 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*15))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(opus_int16(0))*2))))) >> 16)
			res_Q6 = res_Q6 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*14))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(opus_int16(0))*3))))) >> 16)
			res_Q6 = res_Q6 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*13))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(opus_int16(0))*4))))) >> 16)
			res_Q6 = res_Q6 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*12))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(opus_int16(0))*5))))) >> 16)
			res_Q6 = res_Q6 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*11))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(opus_int16(0))*6))))) >> 16)
			res_Q6 = res_Q6 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*10))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(opus_int16(0))*7))))) >> 16)
			res_Q6 = res_Q6 + (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*9))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(opus_int16(0))*8))))) >> 16)
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
		}
	case RESAMPLER_DOWN_ORDER_FIR1:
		for index_Q16 = 0; index_Q16 < max_index_Q16; index_Q16 += index_increment_Q16 {
			buf_ptr = (*opus_int32)(unsafe.Add(unsafe.Pointer(buf), unsafe.Sizeof(opus_int32(0))*uintptr(index_Q16>>16)))
			res_Q6 = (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*0))) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*23)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(opus_int16(0))*0))))) >> 16
			res_Q6 = res_Q6 + ((((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*1))) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*22)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(opus_int16(0))*1))))) >> 16)
			res_Q6 = res_Q6 + ((((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*2))) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*21)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(opus_int16(0))*2))))) >> 16)
			res_Q6 = res_Q6 + ((((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*3))) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*20)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(opus_int16(0))*3))))) >> 16)
			res_Q6 = res_Q6 + ((((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*4))) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*19)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(opus_int16(0))*4))))) >> 16)
			res_Q6 = res_Q6 + ((((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*5))) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*18)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(opus_int16(0))*5))))) >> 16)
			res_Q6 = res_Q6 + ((((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*6))) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*17)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(opus_int16(0))*6))))) >> 16)
			res_Q6 = res_Q6 + ((((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*7))) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*16)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(opus_int16(0))*7))))) >> 16)
			res_Q6 = res_Q6 + ((((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*8))) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*15)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(opus_int16(0))*8))))) >> 16)
			res_Q6 = res_Q6 + ((((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*9))) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*14)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(opus_int16(0))*9))))) >> 16)
			res_Q6 = res_Q6 + ((((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*10))) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*13)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(opus_int16(0))*10))))) >> 16)
			res_Q6 = res_Q6 + ((((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*11))) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*12)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(opus_int16(0))*11))))) >> 16)
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
		}
	case RESAMPLER_DOWN_ORDER_FIR2:
		for index_Q16 = 0; index_Q16 < max_index_Q16; index_Q16 += index_increment_Q16 {
			buf_ptr = (*opus_int32)(unsafe.Add(unsafe.Pointer(buf), unsafe.Sizeof(opus_int32(0))*uintptr(index_Q16>>16)))
			res_Q6 = (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*0))) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*35)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(opus_int16(0))*0))))) >> 16
			res_Q6 = res_Q6 + ((((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*1))) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*34)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(opus_int16(0))*1))))) >> 16)
			res_Q6 = res_Q6 + ((((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*2))) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*33)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(opus_int16(0))*2))))) >> 16)
			res_Q6 = res_Q6 + ((((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*3))) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*32)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(opus_int16(0))*3))))) >> 16)
			res_Q6 = res_Q6 + ((((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*4))) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*31)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(opus_int16(0))*4))))) >> 16)
			res_Q6 = res_Q6 + ((((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*5))) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*30)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(opus_int16(0))*5))))) >> 16)
			res_Q6 = res_Q6 + ((((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*6))) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*29)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(opus_int16(0))*6))))) >> 16)
			res_Q6 = res_Q6 + ((((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*7))) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*28)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(opus_int16(0))*7))))) >> 16)
			res_Q6 = res_Q6 + ((((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*8))) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*27)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(opus_int16(0))*8))))) >> 16)
			res_Q6 = res_Q6 + ((((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*9))) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*26)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(opus_int16(0))*9))))) >> 16)
			res_Q6 = res_Q6 + ((((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*10))) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*25)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(opus_int16(0))*10))))) >> 16)
			res_Q6 = res_Q6 + ((((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*11))) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*24)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(opus_int16(0))*11))))) >> 16)
			res_Q6 = res_Q6 + ((((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*12))) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*23)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(opus_int16(0))*12))))) >> 16)
			res_Q6 = res_Q6 + ((((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*13))) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*22)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(opus_int16(0))*13))))) >> 16)
			res_Q6 = res_Q6 + ((((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*14))) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*21)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(opus_int16(0))*14))))) >> 16)
			res_Q6 = res_Q6 + ((((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*15))) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*20)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(opus_int16(0))*15))))) >> 16)
			res_Q6 = res_Q6 + ((((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*16))) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*19)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(opus_int16(0))*16))))) >> 16)
			res_Q6 = res_Q6 + ((((*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*17))) + (*(*opus_int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(opus_int32(0))*18)))) * opus_int32(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(opus_int16(0))*17))))) >> 16)
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
		}
	default:
	}
	return out
}
func silk_resampler_private_down_FIR(SS unsafe.Pointer, out [0]opus_int16, in [0]opus_int16, inLen opus_int32) {
	var (
		S                   *silk_resampler_state_struct = (*silk_resampler_state_struct)(SS)
		nSamplesIn          opus_int32
		max_index_Q16       opus_int32
		index_increment_Q16 opus_int32
		buf                 *opus_int32
		FIR_Coefs           *opus_int16
	)
	buf = (*opus_int32)(libc.Malloc(int((S.BatchSize + S.FIR_Order) * int64(unsafe.Sizeof(opus_int32(0))))))
	libc.MemCpy(unsafe.Pointer(buf), unsafe.Pointer(&S.SFIR.I32[0]), int(S.FIR_Order*int64(unsafe.Sizeof(opus_int32(0)))))
	FIR_Coefs = (*opus_int16)(unsafe.Add(unsafe.Pointer(S.Coefs), unsafe.Sizeof(opus_int16(0))*2))
	index_increment_Q16 = S.InvRatio_Q16
	for {
		if inLen < opus_int32(S.BatchSize) {
			nSamplesIn = inLen
		} else {
			nSamplesIn = opus_int32(S.BatchSize)
		}
		silk_resampler_private_AR2(S.SIIR[:], [0]opus_int32((*opus_int32)(unsafe.Add(unsafe.Pointer(buf), unsafe.Sizeof(opus_int32(0))*uintptr(S.FIR_Order)))), in, [0]opus_int16(S.Coefs), nSamplesIn)
		max_index_Q16 = opus_int32(opus_uint32(nSamplesIn) << 16)
		out = [0]opus_int16(silk_resampler_private_down_FIR_INTERPOL(&out[0], buf, FIR_Coefs, S.FIR_Order, S.FIR_Fracs, max_index_Q16, index_increment_Q16))
		in += [0]opus_int16(nSamplesIn)
		inLen -= nSamplesIn
		if inLen > 1 {
			libc.MemCpy(unsafe.Pointer(buf), unsafe.Pointer((*opus_int32)(unsafe.Add(unsafe.Pointer(buf), unsafe.Sizeof(opus_int32(0))*uintptr(nSamplesIn)))), int(S.FIR_Order*int64(unsafe.Sizeof(opus_int32(0)))))
		} else {
			break
		}
	}
	libc.MemCpy(unsafe.Pointer(&S.SFIR.I32[0]), unsafe.Pointer((*opus_int32)(unsafe.Add(unsafe.Pointer(buf), unsafe.Sizeof(opus_int32(0))*uintptr(nSamplesIn)))), int(S.FIR_Order*int64(unsafe.Sizeof(opus_int32(0)))))
}
