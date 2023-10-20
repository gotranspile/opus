package libopus

import (
	"math"
	"unsafe"
)

const MAX_LOOPS = 20

func silk_NLSF_stabilize(NLSF_Q15 *int16, NDeltaMin_Q15 *int16, L int) {
	var (
		i               int
		I               int = 0
		k               int
		loops           int
		center_freq_Q15 int16
		diff_Q15        int32
		min_diff_Q15    int32
		min_center_Q15  int32
		max_center_Q15  int32
	)
	for loops = 0; loops < MAX_LOOPS; loops++ {
		min_diff_Q15 = int32(int(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*0))) - int(*(*int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int16(0))*0))))
		I = 0
		for i = 1; i <= L-1; i++ {
			diff_Q15 = int32(int(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(i)))) - (int(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(i-1)))) + int(*(*int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int16(0))*uintptr(i))))))
			if int(diff_Q15) < int(min_diff_Q15) {
				min_diff_Q15 = diff_Q15
				I = i
			}
		}
		diff_Q15 = int32((1 << 15) - (int(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(L-1)))) + int(*(*int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int16(0))*uintptr(L))))))
		if int(diff_Q15) < int(min_diff_Q15) {
			min_diff_Q15 = diff_Q15
			I = L
		}
		if int(min_diff_Q15) >= 0 {
			return
		}
		if I == 0 {
			*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*0)) = *(*int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int16(0))*0))
		} else if I == L {
			*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(L-1))) = int16((1 << 15) - int(*(*int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int16(0))*uintptr(L)))))
		} else {
			min_center_Q15 = 0
			for k = 0; k < I; k++ {
				min_center_Q15 += int32(*(*int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int16(0))*uintptr(k))))
			}
			min_center_Q15 += int32(int(*(*int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int16(0))*uintptr(I)))) >> 1)
			max_center_Q15 = 1 << 15
			for k = L; k > I; k-- {
				max_center_Q15 -= int32(*(*int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int16(0))*uintptr(k))))
			}
			max_center_Q15 -= int32(int(*(*int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int16(0))*uintptr(I)))) >> 1)
			if int(min_center_Q15) > int(max_center_Q15) {
				if (func() int {
					if 1 == 1 {
						return ((int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I-1))))) + int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I)))))) >> 1) + ((int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I-1))))) + int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I)))))) & 1)
					}
					return (((int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I-1))))) + int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I)))))) >> (1 - 1)) + 1) >> 1
				}()) > int(min_center_Q15) {
					center_freq_Q15 = int16(min_center_Q15)
				} else if (func() int {
					if 1 == 1 {
						return ((int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I-1))))) + int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I)))))) >> 1) + ((int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I-1))))) + int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I)))))) & 1)
					}
					return (((int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I-1))))) + int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I)))))) >> (1 - 1)) + 1) >> 1
				}()) < int(max_center_Q15) {
					center_freq_Q15 = int16(max_center_Q15)
				} else if 1 == 1 {
					center_freq_Q15 = int16(((int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I-1))))) + int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I)))))) >> 1) + ((int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I-1))))) + int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I)))))) & 1))
				} else {
					center_freq_Q15 = int16((((int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I-1))))) + int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I)))))) >> (1 - 1)) + 1) >> 1)
				}
			} else if (func() int {
				if 1 == 1 {
					return ((int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I-1))))) + int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I)))))) >> 1) + ((int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I-1))))) + int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I)))))) & 1)
				}
				return (((int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I-1))))) + int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I)))))) >> (1 - 1)) + 1) >> 1
			}()) > int(max_center_Q15) {
				center_freq_Q15 = int16(max_center_Q15)
			} else if (func() int {
				if 1 == 1 {
					return ((int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I-1))))) + int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I)))))) >> 1) + ((int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I-1))))) + int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I)))))) & 1)
				}
				return (((int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I-1))))) + int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I)))))) >> (1 - 1)) + 1) >> 1
			}()) < int(min_center_Q15) {
				center_freq_Q15 = int16(min_center_Q15)
			} else if 1 == 1 {
				center_freq_Q15 = int16(((int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I-1))))) + int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I)))))) >> 1) + ((int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I-1))))) + int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I)))))) & 1))
			} else {
				center_freq_Q15 = int16((((int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I-1))))) + int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I)))))) >> (1 - 1)) + 1) >> 1)
			}
			*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I-1))) = int16(int(center_freq_Q15) - (int(*(*int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int16(0))*uintptr(I)))) >> 1))
			*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I))) = int16(int(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(I-1)))) + int(*(*int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int16(0))*uintptr(I)))))
		}
	}
	if loops == MAX_LOOPS {
		silk_insertion_sort_increasing_all_values_int16((*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*0)), L)
		*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*0)) = int16(silk_max_int(int(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*0))), int(*(*int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int16(0))*0)))))
		for i = 1; i < L; i++ {
			*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(i))) = int16(silk_max_int(int(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(i)))), int(int16(func() int {
				if (int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(i-1))))) + int(*(*int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int16(0))*uintptr(i))))) > silk_int16_MAX {
					return silk_int16_MAX
				}
				if (int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(i-1))))) + int(*(*int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int16(0))*uintptr(i))))) < int(math.MinInt16) {
					return math.MinInt16
				}
				return int(int32(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(i-1))))) + int(*(*int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int16(0))*uintptr(i))))
			}()))))
		}
		*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(L-1))) = int16(silk_min_int(int(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(L-1)))), (1<<15)-int(*(*int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int16(0))*uintptr(L))))))
		for i = L - 2; i >= 0; i-- {
			*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(i))) = int16(silk_min_int(int(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(i)))), int(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF_Q15), unsafe.Sizeof(int16(0))*uintptr(i+1))))-int(*(*int16)(unsafe.Add(unsafe.Pointer(NDeltaMin_Q15), unsafe.Sizeof(int16(0))*uintptr(i+1))))))
		}
	}
}
