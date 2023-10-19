package libopus

import "unsafe"

const QA = 16

func silk_NLSF2A_find_poly(out *opus_int32, cLSF *opus_int32, dd int64) {
	var (
		k    int64
		n    int64
		ftmp opus_int32
	)
	*(*opus_int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int32(0))*0)) = opus_int32(1 << QA)
	*(*opus_int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int32(0))*1)) = -*(*opus_int32)(unsafe.Add(unsafe.Pointer(cLSF), unsafe.Sizeof(opus_int32(0))*0))
	for k = 1; k < dd; k++ {
		ftmp = *(*opus_int32)(unsafe.Add(unsafe.Pointer(cLSF), unsafe.Sizeof(opus_int32(0))*uintptr(k*2)))
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int32(0))*uintptr(k+1))) = (opus_int32(opus_uint32(*(*opus_int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int32(0))*uintptr(k-1)))) << 1)) - opus_int32(func() int64 {
			if QA == 1 {
				return ((int64(ftmp) * int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int32(0))*uintptr(k))))) >> 1) + ((int64(ftmp) * int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int32(0))*uintptr(k))))) & 1)
			}
			return (((int64(ftmp) * int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int32(0))*uintptr(k))))) >> (QA - 1)) + 1) >> 1
		}())
		for n = k; n > 1; n-- {
			*(*opus_int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int32(0))*uintptr(n))) += *(*opus_int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int32(0))*uintptr(n-2))) - opus_int32(func() int64 {
				if QA == 1 {
					return ((int64(ftmp) * int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int32(0))*uintptr(n-1))))) >> 1) + ((int64(ftmp) * int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int32(0))*uintptr(n-1))))) & 1)
				}
				return (((int64(ftmp) * int64(*(*opus_int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int32(0))*uintptr(n-1))))) >> (QA - 1)) + 1) >> 1
			}())
		}
		*(*opus_int32)(unsafe.Add(unsafe.Pointer(out), unsafe.Sizeof(opus_int32(0))*1)) -= ftmp
	}
}
func silk_NLSF2A(a_Q12 *opus_int16, NLSF *opus_int16, d int64, arch int64) {
	var ordering16 [16]uint8 = [16]uint8{0, 15, 8, 7, 4, 11, 12, 3, 2, 13, 10, 5, 6, 9, 14, 1}
	_ = ordering16
	var ordering10 [10]uint8 = [10]uint8{0, 9, 6, 3, 4, 5, 8, 1, 2, 7}
	_ = ordering10
	var ordering *uint8
	var k int64
	var i int64
	var dd int64
	var cos_LSF_QA [24]opus_int32
	var P [13]opus_int32
	var Q [13]opus_int32
	var Ptmp opus_int32
	var Qtmp opus_int32
	var f_int opus_int32
	var f_frac opus_int32
	var cos_val opus_int32
	var delta opus_int32
	var a32_QA1 [24]opus_int32
	ordering = &func() [16]uint8 {
		if d == 16 {
			return ordering16
		}
		return [16]uint8(ordering10)
	}()[0]
	for k = 0; k < d; k++ {
		f_int = opus_int32((*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF), unsafe.Sizeof(opus_int16(0))*uintptr(k)))) >> (15 - 7))
		f_frac = opus_int32(*(*opus_int16)(unsafe.Add(unsafe.Pointer(NLSF), unsafe.Sizeof(opus_int16(0))*uintptr(k)))) - (opus_int32(opus_uint32(f_int) << (15 - 7)))
		cos_val = opus_int32(silk_LSFCosTab_FIX_Q12[f_int])
		delta = opus_int32(silk_LSFCosTab_FIX_Q12[f_int+1]) - cos_val
		if (20 - QA) == 1 {
			cos_LSF_QA[*(*uint8)(unsafe.Add(unsafe.Pointer(ordering), k))] = (((opus_int32(opus_uint32(cos_val) << 8)) + delta*f_frac) >> 1) + (((opus_int32(opus_uint32(cos_val) << 8)) + delta*f_frac) & 1)
		} else {
			cos_LSF_QA[*(*uint8)(unsafe.Add(unsafe.Pointer(ordering), k))] = ((((opus_int32(opus_uint32(cos_val) << 8)) + delta*f_frac) >> opus_int32((20-QA)-1)) + 1) >> 1
		}
	}
	dd = d >> 1
	silk_NLSF2A_find_poly(&P[0], &cos_LSF_QA[0], dd)
	silk_NLSF2A_find_poly(&Q[0], &cos_LSF_QA[1], dd)
	for k = 0; k < dd; k++ {
		Ptmp = P[k+1] + P[k]
		Qtmp = Q[k+1] - Q[k]
		a32_QA1[k] = -Qtmp - Ptmp
		a32_QA1[d-k-1] = Qtmp - Ptmp
	}
	silk_LPC_fit(a_Q12, &a32_QA1[0], 12, QA+1, d)
	for i = 0; (func() opus_int32 {
		_ = arch
		return silk_LPC_inverse_pred_gain_c(a_Q12, d)
	}()) == 0 && i < MAX_LPC_STABILIZE_ITERATIONS; i++ {
		silk_bwexpander_32(&a32_QA1[0], d, 65536-(opus_int32(2<<i)))
		for k = 0; k < d; k++ {
			if (QA + 1 - 12) == 1 {
				*(*opus_int16)(unsafe.Add(unsafe.Pointer(a_Q12), unsafe.Sizeof(opus_int16(0))*uintptr(k))) = opus_int16(((a32_QA1[k]) >> 1) + ((a32_QA1[k]) & 1))
			} else {
				*(*opus_int16)(unsafe.Add(unsafe.Pointer(a_Q12), unsafe.Sizeof(opus_int16(0))*uintptr(k))) = opus_int16((((a32_QA1[k]) >> opus_int32((QA+1-12)-1)) + 1) >> 1)
			}
		}
	}
}
