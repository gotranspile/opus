package libopus

import (
	"github.com/gotranspile/cxgo/runtime/libc"
	"math"
	"unsafe"
)

func silk_resampler_private_down_FIR_INTERPOL(out *int16, buf *int32, FIR_Coefs *int16, FIR_Order int, FIR_Fracs int, max_index_Q16 int32, index_increment_Q16 int32) *int16 {
	var (
		index_Q16    int32
		res_Q6       int32
		buf_ptr      *int32
		interpol_ind int32
		interpol_ptr *int16
	)
	switch FIR_Order {
	case RESAMPLER_DOWN_ORDER_FIR0:
		for index_Q16 = 0; int(index_Q16) < int(max_index_Q16); index_Q16 += index_increment_Q16 {
			buf_ptr = (*int32)(unsafe.Add(unsafe.Pointer(buf), unsafe.Sizeof(int32(0))*uintptr(int(index_Q16)>>16)))
			interpol_ind = int32(((int(index_Q16) & math.MaxUint16) * int(int64(int16(FIR_Fracs)))) >> 16)
			interpol_ptr = (*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*uintptr(int(RESAMPLER_DOWN_ORDER_FIR0/2)*int(interpol_ind))))
			res_Q6 = int32((int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*0))) * int64(*(*int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(int16(0))*0)))) >> 16)
			res_Q6 = int32(int64(res_Q6) + ((int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*1))) * int64(*(*int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(int16(0))*1)))) >> 16))
			res_Q6 = int32(int64(res_Q6) + ((int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*2))) * int64(*(*int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(int16(0))*2)))) >> 16))
			res_Q6 = int32(int64(res_Q6) + ((int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*3))) * int64(*(*int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(int16(0))*3)))) >> 16))
			res_Q6 = int32(int64(res_Q6) + ((int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*4))) * int64(*(*int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(int16(0))*4)))) >> 16))
			res_Q6 = int32(int64(res_Q6) + ((int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*5))) * int64(*(*int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(int16(0))*5)))) >> 16))
			res_Q6 = int32(int64(res_Q6) + ((int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*6))) * int64(*(*int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(int16(0))*6)))) >> 16))
			res_Q6 = int32(int64(res_Q6) + ((int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*7))) * int64(*(*int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(int16(0))*7)))) >> 16))
			res_Q6 = int32(int64(res_Q6) + ((int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*8))) * int64(*(*int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(int16(0))*8)))) >> 16))
			interpol_ptr = (*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*uintptr(int(RESAMPLER_DOWN_ORDER_FIR0/2)*(FIR_Fracs-1-int(interpol_ind)))))
			res_Q6 = int32(int64(res_Q6) + ((int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*17))) * int64(*(*int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(int16(0))*0)))) >> 16))
			res_Q6 = int32(int64(res_Q6) + ((int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*16))) * int64(*(*int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(int16(0))*1)))) >> 16))
			res_Q6 = int32(int64(res_Q6) + ((int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*15))) * int64(*(*int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(int16(0))*2)))) >> 16))
			res_Q6 = int32(int64(res_Q6) + ((int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*14))) * int64(*(*int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(int16(0))*3)))) >> 16))
			res_Q6 = int32(int64(res_Q6) + ((int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*13))) * int64(*(*int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(int16(0))*4)))) >> 16))
			res_Q6 = int32(int64(res_Q6) + ((int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*12))) * int64(*(*int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(int16(0))*5)))) >> 16))
			res_Q6 = int32(int64(res_Q6) + ((int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*11))) * int64(*(*int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(int16(0))*6)))) >> 16))
			res_Q6 = int32(int64(res_Q6) + ((int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*10))) * int64(*(*int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(int16(0))*7)))) >> 16))
			res_Q6 = int32(int64(res_Q6) + ((int64(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*9))) * int64(*(*int16)(unsafe.Add(unsafe.Pointer(interpol_ptr), unsafe.Sizeof(int16(0))*8)))) >> 16))
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
		}
	case RESAMPLER_DOWN_ORDER_FIR1:
		for index_Q16 = 0; int(index_Q16) < int(max_index_Q16); index_Q16 += index_increment_Q16 {
			buf_ptr = (*int32)(unsafe.Add(unsafe.Pointer(buf), unsafe.Sizeof(int32(0))*uintptr(int(index_Q16)>>16)))
			res_Q6 = int32(((int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*0))) + int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*23)))) * int(int64(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*0))))) >> 16)
			res_Q6 = int32(int(res_Q6) + (((int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*1))) + int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*22)))) * int(int64(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*1))))) >> 16))
			res_Q6 = int32(int(res_Q6) + (((int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*2))) + int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*21)))) * int(int64(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*2))))) >> 16))
			res_Q6 = int32(int(res_Q6) + (((int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*3))) + int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*20)))) * int(int64(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*3))))) >> 16))
			res_Q6 = int32(int(res_Q6) + (((int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*4))) + int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*19)))) * int(int64(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*4))))) >> 16))
			res_Q6 = int32(int(res_Q6) + (((int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*5))) + int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*18)))) * int(int64(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*5))))) >> 16))
			res_Q6 = int32(int(res_Q6) + (((int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*6))) + int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*17)))) * int(int64(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*6))))) >> 16))
			res_Q6 = int32(int(res_Q6) + (((int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*7))) + int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*16)))) * int(int64(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*7))))) >> 16))
			res_Q6 = int32(int(res_Q6) + (((int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*8))) + int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*15)))) * int(int64(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*8))))) >> 16))
			res_Q6 = int32(int(res_Q6) + (((int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*9))) + int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*14)))) * int(int64(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*9))))) >> 16))
			res_Q6 = int32(int(res_Q6) + (((int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*10))) + int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*13)))) * int(int64(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*10))))) >> 16))
			res_Q6 = int32(int(res_Q6) + (((int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*11))) + int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*12)))) * int(int64(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*11))))) >> 16))
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
		}
	case RESAMPLER_DOWN_ORDER_FIR2:
		for index_Q16 = 0; int(index_Q16) < int(max_index_Q16); index_Q16 += index_increment_Q16 {
			buf_ptr = (*int32)(unsafe.Add(unsafe.Pointer(buf), unsafe.Sizeof(int32(0))*uintptr(int(index_Q16)>>16)))
			res_Q6 = int32(((int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*0))) + int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*35)))) * int(int64(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*0))))) >> 16)
			res_Q6 = int32(int(res_Q6) + (((int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*1))) + int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*34)))) * int(int64(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*1))))) >> 16))
			res_Q6 = int32(int(res_Q6) + (((int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*2))) + int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*33)))) * int(int64(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*2))))) >> 16))
			res_Q6 = int32(int(res_Q6) + (((int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*3))) + int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*32)))) * int(int64(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*3))))) >> 16))
			res_Q6 = int32(int(res_Q6) + (((int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*4))) + int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*31)))) * int(int64(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*4))))) >> 16))
			res_Q6 = int32(int(res_Q6) + (((int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*5))) + int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*30)))) * int(int64(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*5))))) >> 16))
			res_Q6 = int32(int(res_Q6) + (((int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*6))) + int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*29)))) * int(int64(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*6))))) >> 16))
			res_Q6 = int32(int(res_Q6) + (((int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*7))) + int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*28)))) * int(int64(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*7))))) >> 16))
			res_Q6 = int32(int(res_Q6) + (((int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*8))) + int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*27)))) * int(int64(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*8))))) >> 16))
			res_Q6 = int32(int(res_Q6) + (((int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*9))) + int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*26)))) * int(int64(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*9))))) >> 16))
			res_Q6 = int32(int(res_Q6) + (((int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*10))) + int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*25)))) * int(int64(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*10))))) >> 16))
			res_Q6 = int32(int(res_Q6) + (((int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*11))) + int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*24)))) * int(int64(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*11))))) >> 16))
			res_Q6 = int32(int(res_Q6) + (((int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*12))) + int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*23)))) * int(int64(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*12))))) >> 16))
			res_Q6 = int32(int(res_Q6) + (((int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*13))) + int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*22)))) * int(int64(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*13))))) >> 16))
			res_Q6 = int32(int(res_Q6) + (((int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*14))) + int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*21)))) * int(int64(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*14))))) >> 16))
			res_Q6 = int32(int(res_Q6) + (((int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*15))) + int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*20)))) * int(int64(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*15))))) >> 16))
			res_Q6 = int32(int(res_Q6) + (((int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*16))) + int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*19)))) * int(int64(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*16))))) >> 16))
			res_Q6 = int32(int(res_Q6) + (((int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*17))) + int(*(*int32)(unsafe.Add(unsafe.Pointer(buf_ptr), unsafe.Sizeof(int32(0))*18)))) * int(int64(*(*int16)(unsafe.Add(unsafe.Pointer(FIR_Coefs), unsafe.Sizeof(int16(0))*17))))) >> 16))
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
		}
	default:
	}
	return out
}
func silk_resampler_private_down_FIR(SS unsafe.Pointer, out []int16, in []int16, inLen int32) {
	var (
		S                   *silk_resampler_state_struct = (*silk_resampler_state_struct)(SS)
		nSamplesIn          int32
		max_index_Q16       int32
		index_increment_Q16 int32
		buf                 *int32
		FIR_Coefs           *int16
	)
	buf = (*int32)(libc.Malloc((S.BatchSize + S.FIR_Order) * int(unsafe.Sizeof(int32(0)))))
	libc.MemCpy(unsafe.Pointer(buf), unsafe.Pointer(&S.SFIR.I32[0]), S.FIR_Order*int(unsafe.Sizeof(int32(0))))
	FIR_Coefs = &S.Coefs[2]
	index_increment_Q16 = S.InvRatio_Q16
	for {
		if int(inLen) < S.BatchSize {
			nSamplesIn = inLen
		} else {
			nSamplesIn = int32(S.BatchSize)
		}
		silk_resampler_private_AR2(S.SIIR[:], []int32((*int32)(unsafe.Add(unsafe.Pointer(buf), unsafe.Sizeof(int32(0))*uintptr(S.FIR_Order)))), in, S.Coefs, nSamplesIn)
		max_index_Q16 = int32(int(uint32(nSamplesIn)) << 16)
		out = []int16(silk_resampler_private_down_FIR_INTERPOL(&out[0], buf, FIR_Coefs, S.FIR_Order, S.FIR_Fracs, max_index_Q16, index_increment_Q16))
		in += []int16(nSamplesIn)
		inLen -= nSamplesIn
		if int(inLen) > 1 {
			libc.MemCpy(unsafe.Pointer(buf), unsafe.Pointer((*int32)(unsafe.Add(unsafe.Pointer(buf), unsafe.Sizeof(int32(0))*uintptr(nSamplesIn)))), S.FIR_Order*int(unsafe.Sizeof(int32(0))))
		} else {
			break
		}
	}
	libc.MemCpy(unsafe.Pointer(&S.SFIR.I32[0]), unsafe.Pointer((*int32)(unsafe.Add(unsafe.Pointer(buf), unsafe.Sizeof(int32(0))*uintptr(nSamplesIn)))), S.FIR_Order*int(unsafe.Sizeof(int32(0))))
}
