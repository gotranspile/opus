package libopus

import (
	"math"
	"unsafe"
)

func silk_LPC_fit(a_QOUT *int16, a_QIN *int32, QOUT int, QIN int, d int) {
	var (
		i         int
		k         int
		idx       int = 0
		maxabs    int32
		absval    int32
		chirp_Q16 int32
	)
	for i = 0; i < 10; i++ {
		maxabs = 0
		for k = 0; k < d; k++ {
			if int(*(*int32)(unsafe.Add(unsafe.Pointer(a_QIN), unsafe.Sizeof(int32(0))*uintptr(k)))) > 0 {
				absval = *(*int32)(unsafe.Add(unsafe.Pointer(a_QIN), unsafe.Sizeof(int32(0))*uintptr(k)))
			} else {
				absval = -(*(*int32)(unsafe.Add(unsafe.Pointer(a_QIN), unsafe.Sizeof(int32(0))*uintptr(k))))
			}
			if int(absval) > int(maxabs) {
				maxabs = absval
				idx = k
			}
		}
		if (QIN - QOUT) == 1 {
			maxabs = int32((int(maxabs) >> 1) + (int(maxabs) & 1))
		} else {
			maxabs = int32(((int(maxabs) >> ((QIN - QOUT) - 1)) + 1) >> 1)
		}
		if int(maxabs) > silk_int16_MAX {
			if int(maxabs) < 163838 {
				maxabs = maxabs
			} else {
				maxabs = 163838
			}
			chirp_Q16 = int32(int(int32(math.Floor(0.999*(1<<16)+0.5))) - int(int32(int(int32(int(uint32(int32(int(maxabs)-silk_int16_MAX)))<<14))/((int(maxabs)*(idx+1))>>2))))
			silk_bwexpander_32([]int32(a_QIN), d, chirp_Q16)
		} else {
			break
		}
	}
	if i == 10 {
		for k = 0; k < d; k++ {
			if (func() int {
				if (QIN - QOUT) == 1 {
					return (int(*(*int32)(unsafe.Add(unsafe.Pointer(a_QIN), unsafe.Sizeof(int32(0))*uintptr(k)))) >> 1) + (int(*(*int32)(unsafe.Add(unsafe.Pointer(a_QIN), unsafe.Sizeof(int32(0))*uintptr(k)))) & 1)
				}
				return ((int(*(*int32)(unsafe.Add(unsafe.Pointer(a_QIN), unsafe.Sizeof(int32(0))*uintptr(k)))) >> ((QIN - QOUT) - 1)) + 1) >> 1
			}()) > silk_int16_MAX {
				*(*int16)(unsafe.Add(unsafe.Pointer(a_QOUT), unsafe.Sizeof(int16(0))*uintptr(k))) = silk_int16_MAX
			} else if (func() int {
				if (QIN - QOUT) == 1 {
					return (int(*(*int32)(unsafe.Add(unsafe.Pointer(a_QIN), unsafe.Sizeof(int32(0))*uintptr(k)))) >> 1) + (int(*(*int32)(unsafe.Add(unsafe.Pointer(a_QIN), unsafe.Sizeof(int32(0))*uintptr(k)))) & 1)
				}
				return ((int(*(*int32)(unsafe.Add(unsafe.Pointer(a_QIN), unsafe.Sizeof(int32(0))*uintptr(k)))) >> ((QIN - QOUT) - 1)) + 1) >> 1
			}()) < int(math.MinInt16) {
				*(*int16)(unsafe.Add(unsafe.Pointer(a_QOUT), unsafe.Sizeof(int16(0))*uintptr(k))) = math.MinInt16
			} else if (QIN - QOUT) == 1 {
				*(*int16)(unsafe.Add(unsafe.Pointer(a_QOUT), unsafe.Sizeof(int16(0))*uintptr(k))) = int16((int(*(*int32)(unsafe.Add(unsafe.Pointer(a_QIN), unsafe.Sizeof(int32(0))*uintptr(k)))) >> 1) + (int(*(*int32)(unsafe.Add(unsafe.Pointer(a_QIN), unsafe.Sizeof(int32(0))*uintptr(k)))) & 1))
			} else {
				*(*int16)(unsafe.Add(unsafe.Pointer(a_QOUT), unsafe.Sizeof(int16(0))*uintptr(k))) = int16(((int(*(*int32)(unsafe.Add(unsafe.Pointer(a_QIN), unsafe.Sizeof(int32(0))*uintptr(k)))) >> ((QIN - QOUT) - 1)) + 1) >> 1)
			}
			*(*int32)(unsafe.Add(unsafe.Pointer(a_QIN), unsafe.Sizeof(int32(0))*uintptr(k))) = int32(int(uint32(int32(*(*int16)(unsafe.Add(unsafe.Pointer(a_QOUT), unsafe.Sizeof(int16(0))*uintptr(k)))))) << (QIN - QOUT))
		}
	} else {
		for k = 0; k < d; k++ {
			if (QIN - QOUT) == 1 {
				*(*int16)(unsafe.Add(unsafe.Pointer(a_QOUT), unsafe.Sizeof(int16(0))*uintptr(k))) = int16((int(*(*int32)(unsafe.Add(unsafe.Pointer(a_QIN), unsafe.Sizeof(int32(0))*uintptr(k)))) >> 1) + (int(*(*int32)(unsafe.Add(unsafe.Pointer(a_QIN), unsafe.Sizeof(int32(0))*uintptr(k)))) & 1))
			} else {
				*(*int16)(unsafe.Add(unsafe.Pointer(a_QOUT), unsafe.Sizeof(int16(0))*uintptr(k))) = int16(((int(*(*int32)(unsafe.Add(unsafe.Pointer(a_QIN), unsafe.Sizeof(int32(0))*uintptr(k)))) >> ((QIN - QOUT) - 1)) + 1) >> 1)
			}
		}
	}
}
