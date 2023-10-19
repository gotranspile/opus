package libopus

import (
	"math"
	"unsafe"
)

const MAX_LOOPS = 20

func silk_NLSF_stabilize(NLSF_Q15 *opus_int16, NDeltaMin_Q15 *opus_int16, L int64) {
	var (
		i               int64
		I               int64 = 0
		k               int64
		loops           int64
		center_freq_Q15 opus_int16
		diff_Q15        opus_int32
		min_diff_Q15    opus_int32
		min_center_Q15  opus_int32
		max_center_Q15  opus_int32
	)
	for loops = 0; loops < MAX_LOOPS; loops++ {
		min_diff_Q15 = opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*0)) - *(*opus_int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(opus_int16(0))*0)))
		I = 0
		for i = 1; i <= L-1; i++ {
			diff_Q15 = opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(i))) - (*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(i-1))) + *(*opus_int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(i)))))
			if diff_Q15 < min_diff_Q15 {
				min_diff_Q15 = diff_Q15
				I = i
			}
		}
		diff_Q15 = opus_int32((1 << 15) - (*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(L-1))) + *(*opus_int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(L)))))
		if diff_Q15 < min_diff_Q15 {
			min_diff_Q15 = diff_Q15
			I = L
		}
		if min_diff_Q15 >= 0 {
			return
		}
		if I == 0 {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*0)) = *(*opus_int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(opus_int16(0))*0))
		} else if I == L {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(L-1))) = (1 << 15) - *(*opus_int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(L)))
		} else {
			min_center_Q15 = 0
			for k = 0; k < I; k++ {
				min_center_Q15 += opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(k))))
			}
			min_center_Q15 += opus_int32((*(*opus_int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I)))) >> 1)
			max_center_Q15 = 1 << 15
			for k = L; k > I; k-- {
				max_center_Q15 -= opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(k))))
			}
			max_center_Q15 -= opus_int32((*(*opus_int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I)))) >> 1)
			if min_center_Q15 > max_center_Q15 {
				if (func() opus_int32 {
					if 1 == 1 {
						return ((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I-1)))) + opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I))))) >> 1) + ((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I-1)))) + opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I))))) & 1)
					}
					return (((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I-1)))) + opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I))))) >> (1 - 1)) + 1) >> 1
				}()) > min_center_Q15 {
					center_freq_Q15 = opus_int16(min_center_Q15)
				} else if (func() opus_int32 {
					if 1 == 1 {
						return ((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I-1)))) + opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I))))) >> 1) + ((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I-1)))) + opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I))))) & 1)
					}
					return (((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I-1)))) + opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I))))) >> (1 - 1)) + 1) >> 1
				}()) < max_center_Q15 {
					center_freq_Q15 = opus_int16(max_center_Q15)
				} else if 1 == 1 {
					center_freq_Q15 = opus_int16(((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I-1)))) + opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I))))) >> 1) + ((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I-1)))) + opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I))))) & 1))
				} else {
					center_freq_Q15 = opus_int16((((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I-1)))) + opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I))))) >> (1 - 1)) + 1) >> 1)
				}
			} else if (func() opus_int32 {
				if 1 == 1 {
					return ((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I-1)))) + opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I))))) >> 1) + ((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I-1)))) + opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I))))) & 1)
				}
				return (((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I-1)))) + opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I))))) >> (1 - 1)) + 1) >> 1
			}()) > max_center_Q15 {
				center_freq_Q15 = opus_int16(max_center_Q15)
			} else if (func() opus_int32 {
				if 1 == 1 {
					return ((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I-1)))) + opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I))))) >> 1) + ((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I-1)))) + opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I))))) & 1)
				}
				return (((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I-1)))) + opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I))))) >> (1 - 1)) + 1) >> 1
			}()) < min_center_Q15 {
				center_freq_Q15 = opus_int16(min_center_Q15)
			} else if 1 == 1 {
				center_freq_Q15 = opus_int16(((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I-1)))) + opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I))))) >> 1) + ((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I-1)))) + opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I))))) & 1))
			} else {
				center_freq_Q15 = opus_int16((((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I-1)))) + opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I))))) >> (1 - 1)) + 1) >> 1)
			}
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I-1))) = center_freq_Q15 - ((*(*opus_int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I)))) >> 1)
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I))) = *(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I-1))) + *(*opus_int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(I)))
		}
	}
	if loops == MAX_LOOPS {
		silk_insertion_sort_increasing_all_values_int16((*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*0)), L)
		*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*0)) = opus_int16(silk_max_int(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*0))), int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(opus_int16(0))*0)))))
		for i = 1; i < L; i++ {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(i))) = opus_int16(silk_max_int(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(i)))), int64(opus_int16(func() opus_int32 {
				if ((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(i-1))))) + opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(i))))) > silk_int16_MAX {
					return silk_int16_MAX
				}
				if ((opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(i-1))))) + opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(i))))) < opus_int32(math.MinInt16) {
					return math.MinInt16
				}
				return (opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(i-1))))) + opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(i))))
			}()))))
		}
		*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(L-1))) = opus_int16(silk_min_int(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(L-1)))), int64((1<<15)-*(*opus_int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(L))))))
		for i = L - 2; i >= 0; i-- {
			*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(i))) = opus_int16(silk_min_int(int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(i)))), int64(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(i+1)))-*(*opus_int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(opus_int16(0))*uintptr(i+1))))))
		}
	}
}
