package libopus

import (
	"math"
	"unsafe"
)

func silk_LPC_fit(a_QOUT *opus_int16, a_QIN *opus_int32, QOUT int64, QIN int64, d int64) {
	var (
		i         int64
		k         int64
		idx       int64 = 0
		maxabs    opus_int32
		absval    opus_int32
		chirp_Q16 opus_int32
	)
	for i = 0; i < 10; i++ {
		maxabs = 0
		for k = 0; k < d; k++ {
			if (*(*opus_int32)(unsafe.Add(unsafe.Pointer(a_QIN), unsafe.Sizeof(opus_int32(0))*uintptr(k)))) > 0 {
				absval = *(*opus_int32)(unsafe.Add(unsafe.Pointer(a_QIN), unsafe.Sizeof(opus_int32(0))*uintptr(k)))
			} else {
				absval = -(*(*opus_int32)(unsafe.Add(unsafe.Pointer(a_QIN), unsafe.Sizeof(opus_int32(0))*uintptr(k))))
			}
			if absval > maxabs {
				maxabs = absval
				idx = k
			}
		}
		if (QIN - QOUT) == 1 {
			maxabs = (maxabs >> 1) + (maxabs & 1)
		} else {
			maxabs = ((maxabs >> opus_int32((QIN-QOUT)-1)) + 1) >> 1
		}
		if maxabs > silk_int16_MAX {
			if maxabs < 163838 {
				maxabs = maxabs
			} else {
				maxabs = 163838
			}
			chirp_Q16 = (opus_int32(0.999*(1<<16) + 0.5)) - (opus_int32(opus_uint32(maxabs-silk_int16_MAX)<<14))/((maxabs*opus_int32(idx+1))>>2)
			silk_bwexpander_32(a_QIN, d, chirp_Q16)
		} else {
			break
		}
	}
	if i == 10 {
		for k = 0; k < d; k++ {
			if (func() opus_int32 {
				if (QIN - QOUT) == 1 {
					return ((*(*opus_int32)(unsafe.Add(unsafe.Pointer(a_QIN), unsafe.Sizeof(opus_int32(0))*uintptr(k)))) >> 1) + ((*(*opus_int32)(unsafe.Add(unsafe.Pointer(a_QIN), unsafe.Sizeof(opus_int32(0))*uintptr(k)))) & 1)
				}
				return (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(a_QIN), unsafe.Sizeof(opus_int32(0))*uintptr(k)))) >> opus_int32((QIN-QOUT)-1)) + 1) >> 1
			}()) > silk_int16_MAX {
				*(*opus_int16)(unsafe.Add(unsafe.Pointer(a_QOUT), unsafe.Sizeof(opus_int16(0))*uintptr(k))) = silk_int16_MAX
			} else if (func() opus_int32 {
				if (QIN - QOUT) == 1 {
					return ((*(*opus_int32)(unsafe.Add(unsafe.Pointer(a_QIN), unsafe.Sizeof(opus_int32(0))*uintptr(k)))) >> 1) + ((*(*opus_int32)(unsafe.Add(unsafe.Pointer(a_QIN), unsafe.Sizeof(opus_int32(0))*uintptr(k)))) & 1)
				}
				return (((*(*opus_int32)(unsafe.Add(unsafe.Pointer(a_QIN), unsafe.Sizeof(opus_int32(0))*uintptr(k)))) >> opus_int32((QIN-QOUT)-1)) + 1) >> 1
			}()) < opus_int32(math.MinInt16) {
				*(*opus_int16)(unsafe.Add(unsafe.Pointer(a_QOUT), unsafe.Sizeof(opus_int16(0))*uintptr(k))) = math.MinInt16
			} else if (QIN - QOUT) == 1 {
				*(*opus_int16)(unsafe.Add(unsafe.Pointer(a_QOUT), unsafe.Sizeof(opus_int16(0))*uintptr(k))) = opus_int16(((*(*opus_int32)(unsafe.Add(unsafe.Pointer(a_QIN), unsafe.Sizeof(opus_int32(0))*uintptr(k)))) >> 1) + ((*(*opus_int32)(unsafe.Add(unsafe.Pointer(a_QIN), unsafe.Sizeof(opus_int32(0))*uintptr(k)))) & 1))
			} else {
				*(*opus_int16)(unsafe.Add(unsafe.Pointer(a_QOUT), unsafe.Sizeof(opus_int16(0))*uintptr(k))) = opus_int16((((*(*opus_int32)(unsafe.Add(unsafe.Pointer(a_QIN), unsafe.Sizeof(opus_int32(0))*uintptr(k)))) >> opus_int32((QIN-QOUT)-1)) + 1) >> 1)
			}
			*(*opus_int32)(unsafe.Add(unsafe.Pointer(a_QIN), unsafe.Sizeof(opus_int32(0))*uintptr(k))) = opus_int32(opus_uint32(opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(a_QOUT), unsafe.Sizeof(opus_int16(0))*uintptr(k))))) << opus_uint32(QIN-QOUT))
		}
	} else {
		for k = 0; k < d; k++ {
			if (QIN - QOUT) == 1 {
				*(*opus_int16)(unsafe.Add(unsafe.Pointer(a_QOUT), unsafe.Sizeof(opus_int16(0))*uintptr(k))) = opus_int16(((*(*opus_int32)(unsafe.Add(unsafe.Pointer(a_QIN), unsafe.Sizeof(opus_int32(0))*uintptr(k)))) >> 1) + ((*(*opus_int32)(unsafe.Add(unsafe.Pointer(a_QIN), unsafe.Sizeof(opus_int32(0))*uintptr(k)))) & 1))
			} else {
				*(*opus_int16)(unsafe.Add(unsafe.Pointer(a_QOUT), unsafe.Sizeof(opus_int16(0))*uintptr(k))) = opus_int16((((*(*opus_int32)(unsafe.Add(unsafe.Pointer(a_QIN), unsafe.Sizeof(opus_int32(0))*uintptr(k)))) >> opus_int32((QIN-QOUT)-1)) + 1) >> 1)
			}
		}
	}
}
