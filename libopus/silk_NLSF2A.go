package libopus

import "unsafe"

const NLSF2A_QA = 16

func silk_NLSF2A_find_poly(out *int32, cLSF *int32, dd int) {
	var (
		k    int
		n    int
		ftmp int32
	)
	*(*int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int32(0))*0)) = int32(int(1 << NLSF2A_QA))
	*(*int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int32(0))*1)) = -*(*int32)(unsafe.Add(unsafe.Pointer(cLSF), unsafe.Sizeof(int32(0))*0))
	for k = 1; k < dd; k++ {
		ftmp = *(*int32)(unsafe.Add(unsafe.Pointer(cLSF), unsafe.Sizeof(int32(0))*uintptr(k*2)))
		*(*int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int32(0))*uintptr(k+1))) = int32(int(int32(int(uint32(*(*int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int32(0))*uintptr(k-1)))))<<1)) - int(int32(func() int64 {
			if NLSF2A_QA == 1 {
				return ((int64(ftmp) * int64(*(*int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int32(0))*uintptr(k))))) >> 1) + ((int64(ftmp) * int64(*(*int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int32(0))*uintptr(k))))) & 1)
			}
			return (((int64(ftmp) * int64(*(*int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int32(0))*uintptr(k))))) >> int64(int(NLSF2A_QA-1))) + 1) >> 1
		}())))
		for n = k; n > 1; n-- {
			*(*int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int32(0))*uintptr(n))) += int32(int(*(*int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int32(0))*uintptr(n-2)))) - int(int32(func() int64 {
				if NLSF2A_QA == 1 {
					return ((int64(ftmp) * int64(*(*int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int32(0))*uintptr(n-1))))) >> 1) + ((int64(ftmp) * int64(*(*int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int32(0))*uintptr(n-1))))) & 1)
				}
				return (((int64(ftmp) * int64(*(*int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int32(0))*uintptr(n-1))))) >> int64(int(NLSF2A_QA-1))) + 1) >> 1
			}())))
		}
		*(*int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(int32(0))*1)) -= ftmp
	}
}
func silk_NLSF2A(a_Q12 *int16, NLSF *int16, d int, arch int) {
	var ordering16 [16]uint8 = [16]uint8{0, 15, 8, 7, 4, 11, 12, 3, 2, 13, 10, 5, 6, 9, 14, 1}
	_ = ordering16
	var ordering10 [10]uint8 = [10]uint8{0, 9, 6, 3, 4, 5, 8, 1, 2, 7}
	_ = ordering10
	var ordering *uint8
	var k int
	var i int
	var dd int
	var cos_LSF_NLSF2A_QA [24]int32
	var P [13]int32
	var Q [13]int32
	var Ptmp int32
	var Qtmp int32
	var f_int int32
	var f_frac int32
	var cos_val int32
	var delta int32
	var a32_NLSF2A_QA1 [24]int32
	ordering = &func() [16]uint8 {
		if d == 16 {
			return ordering16
		}
		return [16]uint8(ordering10)
	}()[0]
	for k = 0; k < d; k++ {
		f_int = int32(int(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF), unsafe.Sizeof(int16(0))*uintptr(k)))) >> (15 - 7))
		f_frac = int32(int(*(*int16)(unsafe.Add(unsafe.Pointer(NLSF), unsafe.Sizeof(int16(0))*uintptr(k)))) - int(int32(int(uint32(f_int))<<(15-7))))
		cos_val = int32(silk_LSFCosTab_FIX_Q12[f_int])
		delta = int32(int(silk_LSFCosTab_FIX_Q12[int(f_int)+1]) - int(cos_val))
		if (int(20 - NLSF2A_QA)) == 1 {
			cos_LSF_NLSF2A_QA[*(*uint8)(unsafe.Add(unsafe.Pointer(ordering), k))] = int32(((int(int32(int(uint32(cos_val))<<8)) + int(delta)*int(f_frac)) >> 1) + ((int(int32(int(uint32(cos_val))<<8)) + int(delta)*int(f_frac)) & 1))
		} else {
			cos_LSF_NLSF2A_QA[*(*uint8)(unsafe.Add(unsafe.Pointer(ordering), k))] = int32((((int(int32(int(uint32(cos_val))<<8)) + int(delta)*int(f_frac)) >> ((int(20 - NLSF2A_QA)) - 1)) + 1) >> 1)
		}
	}
	dd = d >> 1
	silk_NLSF2A_find_poly(&P[0], &cos_LSF_NLSF2A_QA[0], dd)
	silk_NLSF2A_find_poly(&Q[0], &cos_LSF_NLSF2A_QA[1], dd)
	for k = 0; k < dd; k++ {
		Ptmp = int32(int(P[k+1]) + int(P[k]))
		Qtmp = int32(int(Q[k+1]) - int(Q[k]))
		a32_NLSF2A_QA1[k] = int32(int(-Qtmp) - int(Ptmp))
		a32_NLSF2A_QA1[d-k-1] = int32(int(Qtmp) - int(Ptmp))
	}
	silk_LPC_fit(a_Q12, &a32_NLSF2A_QA1[0], 12, int(NLSF2A_QA+1), d)
	for i = 0; int(func() int32 {
		_ = arch
		return silk_LPC_inverse_pred_gain_c(a_Q12, d)
	}()) == 0 && i < MAX_LPC_STABILIZE_ITERATIONS; i++ {
		silk_bwexpander_32(a32_NLSF2A_QA1[:], d, int32(65536-int(int32(2<<i))))
		for k = 0; k < d; k++ {
			if (int(NLSF2A_QA+1) - 12) == 1 {
				*(*int16)(unsafe.Add(unsafe.Pointer(a_Q12), unsafe.Sizeof(int16(0))*uintptr(k))) = int16((int(a32_NLSF2A_QA1[k]) >> 1) + (int(a32_NLSF2A_QA1[k]) & 1))
			} else {
				*(*int16)(unsafe.Add(unsafe.Pointer(a_Q12), unsafe.Sizeof(int16(0))*uintptr(k))) = int16(((int(a32_NLSF2A_QA1[k]) >> ((int(NLSF2A_QA+1) - 12) - 1)) + 1) >> 1)
			}
		}
	}
}
